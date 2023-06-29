package resources

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

const (
	fileFormatIDDelimiter = '|'
)

// valid format type options for each File Format Type
// https://docs.snowflake.com/en/sql-reference/sql/create-file-format.html#syntax
var formatTypeOptions = map[string][]string{
	"CSV": {
		"compression",
		"record_delimiter",
		"field_delimiter",
		"file_extension",
		"skip_header",
		"skip_blank_lines",
		"date_format",
		"time_format",
		"timestamp_format",
		"binary_format",
		"escape",
		"escape_unenclosed_field",
		"trim_space",
		"field_optionally_enclosed_by",
		"null_if",
		"error_on_column_count_mismatch",
		"replace_invalid_characters",
		"empty_field_as_null",
		"skip_byte_order_mark",
		"encoding",
	},
	"JSON": {
		"compression",
		"date_format",
		"time_format",
		"timestamp_format",
		"binary_format",
		"trim_space",
		"null_if",
		"file_extension",
		"enable_octal",
		"allow_duplicate",
		"strip_outer_array",
		"strip_null_values",
		"replace_invalid_characters",
		"ignore_utf8_errors",
		"skip_byte_order_mark",
	},
	"AVRO": {
		"compression",
		"trim_space",
		"null_if",
	},
	"ORC": {
		"trim_space",
		"null_if",
	},
	"PARQUET": {
		"compression",
		"binary_as_text",
		"trim_space",
		"null_if",
	},
	"XML": {
		"compression",
		"ignore_utf8_errors",
		"preserve_space",
		"strip_outer_element",
		"disable_snowflake_data",
		"disable_auto_convert",
		"skip_byte_order_mark",
	},
}

var fileFormatSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the file format; must be unique for the database and schema in which the file format is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the file format.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the file format.",
		ForceNew:    true,
	},
	"format_type": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Specifies the format of the input files (for data loading) or output files (for data unloading).",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"CSV", "JSON", "AVRO", "ORC", "PARQUET", "XML"}, true),
	},
	"compression": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Specifies the current compression algorithm for the data file.",
	},
	"record_delimiter": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Specifies one or more singlebyte or multibyte characters that separate records in an input file (data loading) or unloaded file (data unloading).",
	},
	"field_delimiter": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Specifies one or more singlebyte or multibyte characters that separate fields in an input file (data loading) or unloaded file (data unloading).",
	},
	"file_extension": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the extension for files unloaded to a stage.",
	},
	"skip_header": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Number of lines at the start of the file to skip.",
	},
	"skip_blank_lines": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies to skip any blank lines encountered in the data files.",
	},
	"date_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Defines the format of date values in the data files (data loading) or table (data unloading).",
	},
	"time_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Defines the format of time values in the data files (data loading) or table (data unloading).",
	},
	"timestamp_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Defines the format of timestamp values in the data files (data loading) or table (data unloading).",
	},
	"binary_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Defines the encoding format for binary input or output.",
	},
	"escape": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Single character string used as the escape character for field values.",
	},
	"escape_unenclosed_field": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Single character string used as the escape character for unenclosed field values only.",
	},
	"trim_space": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to remove white space from fields.",
	},
	"field_optionally_enclosed_by": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "Character used to enclose strings.",
	},
	"null_if": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Computed:    true,
		Description: "String used to convert to and from SQL NULL.",
	},
	"error_on_column_count_mismatch": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to generate a parsing error if the number of delimited columns (i.e. fields) in an input file does not match the number of columns in the corresponding table.",
	},
	"replace_invalid_characters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to replace invalid UTF-8 characters with the Unicode replacement character (�).",
	},
	"empty_field_as_null": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether to insert SQL NULL for empty fields in an input file, which are represented by two successive delimiters.",
	},
	"skip_byte_order_mark": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to skip the BOM (byte order mark), if present in a data file.",
	},
	"encoding": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "String (constant) that specifies the character set of the source data when loading data into a table.",
	},
	"enable_octal": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that enables parsing of octal numbers.",
	},
	"allow_duplicate": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies to allow duplicate object field names (only the last one will be preserved).",
	},
	"strip_outer_array": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that instructs the JSON parser to remove outer brackets.",
	},
	"strip_null_values": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that instructs the JSON parser to remove object fields or array elements containing null values.",
	},
	"ignore_utf8_errors": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether UTF-8 encoding errors produce error conditions.",
	},
	"binary_as_text": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to interpret columns with no defined logical data type as UTF-8 text.",
	},
	"preserve_space": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether the XML parser preserves leading and trailing spaces in element content.",
	},
	"strip_outer_element": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether the XML parser strips out the outer XML element, exposing 2nd level elements as separate documents.",
	},
	"disable_snowflake_data": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether the XML parser disables recognition of Snowflake semi-structured data tags.",
	},
	"disable_auto_convert": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether the XML parser disables automatic conversion of numeric and Boolean values from text to native representation.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the file format.",
	},
}

type fileFormatID struct {
	DatabaseName   string
	SchemaName     string
	FileFormatName string
}

func (ffi *fileFormatID) String() (string, error) {
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = fileFormatIDDelimiter
	if err := csvWriter.WriteAll([][]string{{ffi.DatabaseName, ffi.SchemaName, ffi.FileFormatName}}); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

// FileFormat returns a pointer to the resource representing a file format.
func FileFormat() *schema.Resource {
	return &schema.Resource{
		Create: CreateFileFormat,
		Read:   ReadFileFormat,
		Update: UpdateFileFormat,
		Delete: DeleteFileFormat,

		Schema: fileFormatSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateFileFormat implements schema.CreateFunc.
func CreateFileFormat(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	dbName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	fileFormatName := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(dbName, schemaName, fileFormatName)

	opts := sdk.CreateFileFormatOptions{
		Type:                  sdk.FileFormatType(d.Get("format_type").(string)),
		FileFormatTypeOptions: sdk.FileFormatTypeOptions{},
	}

	switch opts.Type {
	case sdk.FileFormatTypeCsv:
		if v, ok := d.GetOk("compression"); ok {
			comp := sdk.CsvCompression(v.(string))
			opts.CsvCompression = &comp
		}
		if v, ok := d.GetOk("record_delimiter"); ok {
			opts.CsvRecordDelimiter = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("field_delimiter"); ok {
			opts.CsvFieldDelimiter = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("file_extension"); ok {
			opts.CsvFileExtension = sdk.String(v.(string))
		}
		opts.CsvSkipHeader = sdk.Int(d.Get("skip_header").(int))
		opts.CsvSkipBlankLines = sdk.Bool(d.Get("skip_blank_lines").(bool))
		if v, ok := d.GetOk("date_format"); ok {
			opts.CsvDateFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("time_format"); ok {
			opts.CsvTimeFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("timestamp_format"); ok {
			opts.CsvTimestampFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("binary_format"); ok {
			bf := sdk.BinaryFormat(v.(string))
			opts.CsvBinaryFormat = &bf
		}
		if v, ok := d.GetOk("escape"); ok {
			opts.CsvEscape = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("escape_unenclosed_field"); ok {
			opts.CsvEscapeUnenclosedField = sdk.String(v.(string))
		}
		opts.CsvTrimSpace = sdk.Bool(d.Get("trim_space").(bool))
		if v, ok := d.GetOk("field_optionally_enclosed_by"); ok {
			opts.CsvFieldOptionallyEnclosedBy = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("null_if"); ok {
			nullIf := []sdk.NullString{}
			for _, s := range v.([]interface{}) {
				if s == nil {
					s = ""
				} else {
					s = s.(string)
				}
				nullIf = append(nullIf, sdk.NullString{S: s.(string)})
			}
			opts.CsvNullIf = &nullIf
		}
		opts.CsvErrorOnColumnCountMismatch = sdk.Bool(d.Get("error_on_column_count_mismatch").(bool))
		opts.CsvReplaceInvalidCharacters = sdk.Bool(d.Get("replace_invalid_characters").(bool))
		opts.CsvEmptyFieldAsNull = sdk.Bool(d.Get("empty_field_as_null").(bool))
		opts.CsvSkipByteOrderMark = sdk.Bool(d.Get("skip_byte_order_mark").(bool))
		if v, ok := d.GetOk("encoding"); ok {
			enc := sdk.CsvEncoding(v.(string))
			opts.CsvEncoding = &enc
		}
	case sdk.FileFormatTypeJson:
		if v, ok := d.GetOk("compression"); ok {
			comp := sdk.JsonCompression(v.(string))
			opts.JsonCompression = &comp
		}
		if v, ok := d.GetOk("date_format"); ok {
			opts.JsonDateFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("time_format"); ok {
			opts.JsonTimeFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("timestamp_format"); ok {
			opts.JsonTimestampFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("binary_format"); ok {
			bf := sdk.BinaryFormat(v.(string))
			opts.JsonBinaryFormat = &bf
		}
		opts.JsonTrimSpace = sdk.Bool(d.Get("trim_space").(bool))
		if v, ok := d.GetOk("null_if"); ok {
			nullIf := []sdk.NullString{}
			for _, s := range v.([]interface{}) {
				if s == nil {
					s = ""
				} else {
					s = s.(string)
				}
				nullIf = append(nullIf, sdk.NullString{S: s.(string)})
			}
			opts.JsonNullIf = &nullIf
		}
		if v, ok := d.GetOk("file_extension"); ok {
			opts.JsonFileExtension = sdk.String(v.(string))
		}
		opts.JsonEnableOctal = sdk.Bool(d.Get("enable_octal").(bool))
		opts.JsonAllowDuplicate = sdk.Bool(d.Get("allow_duplicate").(bool))
		opts.JsonStripOuterArray = sdk.Bool(d.Get("strip_outer_array").(bool))
		opts.JsonStripNullValues = sdk.Bool(d.Get("strip_null_values").(bool))
		opts.JsonReplaceInvalidCharacters = sdk.Bool(d.Get("replace_invalid_characters").(bool))
		opts.JsonIgnoreUtf8Errors = sdk.Bool(d.Get("ignore_utf8_errors").(bool))
		opts.JsonSkipByteOrderMark = sdk.Bool(d.Get("skip_byte_order_mark").(bool))
	case sdk.FileFormatTypeAvro:
		if v, ok := d.GetOk("compression"); ok {
			comp := sdk.AvroCompression(v.(string))
			opts.AvroCompression = &comp
		}
		opts.AvroTrimSpace = sdk.Bool(d.Get("trim_space").(bool))
		if v, ok := d.GetOk("null_if"); ok {
			nullIf := []sdk.NullString{}
			for _, s := range v.([]interface{}) {
				if s == nil {
					s = ""
				} else {
					s = s.(string)
				}
				nullIf = append(nullIf, sdk.NullString{S: s.(string)})
			}
			opts.AvroNullIf = &nullIf
		}
	case sdk.FileFormatTypeOrc:
		opts.OrcTrimSpace = sdk.Bool(d.Get("trim_space").(bool))
		if v, ok := d.GetOk("null_if"); ok {
			nullIf := []sdk.NullString{}
			for _, s := range v.([]interface{}) {
				if s == nil {
					s = ""
				} else {
					s = s.(string)
				}
				nullIf = append(nullIf, sdk.NullString{S: s.(string)})
			}
			opts.OrcNullIf = &nullIf
		}
	case sdk.FileFormatTypeParquet:
		if v, ok := d.GetOk("compression"); ok {
			comp := sdk.ParquetCompression(v.(string))
			opts.ParquetCompression = &comp
		}
		opts.ParquetBinaryAsText = sdk.Bool(d.Get("binary_as_text").(bool))
		opts.ParquetTrimSpace = sdk.Bool(d.Get("trim_space").(bool))
		if v, ok := d.GetOk("null_if"); ok {
			nullIf := []sdk.NullString{}
			for _, s := range v.([]interface{}) {
				if s == nil {
					s = ""
				} else {
					s = s.(string)
				}
				nullIf = append(nullIf, sdk.NullString{S: s.(string)})
			}
			opts.ParquetNullIf = &nullIf
		}
	case sdk.FileFormatTypeXml:
		if v, ok := d.GetOk("compression"); ok {
			comp := sdk.XmlCompression(v.(string))
			opts.XmlCompression = &comp
		}
		opts.XmlIgnoreUtf8Errors = sdk.Bool(d.Get("ignore_utf8_errors").(bool))
		opts.XmlPreserveSpace = sdk.Bool(d.Get("preserve_space").(bool))
		opts.XmlStripOuterElement = sdk.Bool(d.Get("strip_outer_element").(bool))
		opts.XmlDisableSnowflakeData = sdk.Bool(d.Get("disable_snowflake_data").(bool))
		opts.XmlDisableAutoConvert = sdk.Bool(d.Get("disable_auto_convert").(bool))
		opts.XmlSkipByteOrderMark = sdk.Bool(d.Get("skip_byte_order_mark").(bool))
	}

	if v, ok := d.GetOk("comment"); ok {
		opts.Comment = sdk.String(v.(string))
	}

	err := client.FileFormats.Create(ctx, id, &opts)
	if err != nil {
		return err
	}

	fileFormatID := &fileFormatID{
		DatabaseName:   dbName,
		SchemaName:     schemaName,
		FileFormatName: fileFormatName,
	}
	dataIDInput, err := fileFormatID.String()
	if err != nil {
		return err
	}
	d.SetId(dataIDInput)

	return ReadFileFormat(d, meta)
}

// ReadFileFormat implements schema.ReadFunc.
func ReadFileFormat(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	fileFormatID, err := fileFormatIDFromString(d.Id())
	if err != nil {
		return err
	}
	id := sdk.NewSchemaObjectIdentifier(fileFormatID.DatabaseName, fileFormatID.SchemaName, fileFormatID.FileFormatName)

	fileFormat, err := client.FileFormats.ShowByID(ctx, id)
	if err != nil {
		return fmt.Errorf("cannot read file format: %w", err)
	}

	if err := d.Set("name", fileFormat.Name.Name()); err != nil {
		return err
	}

	if err := d.Set("database", fileFormat.Name.DatabaseName()); err != nil {
		return err
	}

	if err := d.Set("schema", fileFormat.Name.SchemaName()); err != nil {
		return err
	}

	if err := d.Set("format_type", fileFormat.Type); err != nil {
		return err
	}

	switch fileFormat.Type {
	case sdk.FileFormatTypeCsv:
		if err := d.Set("compression", fileFormat.Options.CsvCompression); err != nil {
			return err
		}
		if err := d.Set("record_delimiter", fileFormat.Options.CsvRecordDelimiter); err != nil {
			return err
		}
		if err := d.Set("field_delimiter", fileFormat.Options.CsvFieldDelimiter); err != nil {
			return err
		}
		if err := d.Set("file_extension", fileFormat.Options.CsvFileExtension); err != nil {
			return err
		}
		if err := d.Set("skip_header", fileFormat.Options.CsvSkipHeader); err != nil {
			return err
		}
		if err := d.Set("skip_blank_lines", fileFormat.Options.CsvSkipBlankLines); err != nil {
			return err
		}
		if err := d.Set("date_format", fileFormat.Options.CsvDateFormat); err != nil {
			return err
		}
		if err := d.Set("time_format", fileFormat.Options.CsvTimeFormat); err != nil {
			return err
		}
		if err := d.Set("timestamp_format", fileFormat.Options.CsvTimestampFormat); err != nil {
			return err
		}
		if err := d.Set("binary_format", fileFormat.Options.CsvBinaryFormat); err != nil {
			return err
		}
		if err := d.Set("escape", fileFormat.Options.CsvEscape); err != nil {
			return err
		}
		if err := d.Set("escape_unenclosed_field", fileFormat.Options.CsvEscapeUnenclosedField); err != nil {
			return err
		}
		if err := d.Set("trim_space", fileFormat.Options.CsvTrimSpace); err != nil {
			return err
		}
		if err := d.Set("field_optionally_enclosed_by", fileFormat.Options.CsvFieldOptionallyEnclosedBy); err != nil {
			return err
		}
		nullIf := []string{}
		for _, s := range *fileFormat.Options.CsvNullIf {
			nullIf = append(nullIf, s.S)
		}
		if err := d.Set("null_if", nullIf); err != nil {
			return err
		}
		if err := d.Set("error_on_column_count_mismatch", fileFormat.Options.CsvErrorOnColumnCountMismatch); err != nil {
			return err
		}
		if err := d.Set("replace_invalid_characters", fileFormat.Options.CsvReplaceInvalidCharacters); err != nil {
			return err
		}
		if err := d.Set("empty_field_as_null", fileFormat.Options.CsvEmptyFieldAsNull); err != nil {
			return err
		}
		if err := d.Set("skip_byte_order_mark", fileFormat.Options.CsvSkipByteOrderMark); err != nil {
			return err
		}
		if err := d.Set("encoding", fileFormat.Options.CsvEncoding); err != nil {
			return err
		}
	case sdk.FileFormatTypeJson:
		if err := d.Set("compression", fileFormat.Options.JsonCompression); err != nil {
			return err
		}
		if err := d.Set("date_format", fileFormat.Options.JsonDateFormat); err != nil {
			return err
		}
		if err := d.Set("time_format", fileFormat.Options.JsonTimeFormat); err != nil {
			return err
		}
		if err := d.Set("timestamp_format", fileFormat.Options.JsonTimestampFormat); err != nil {
			return err
		}
		if err := d.Set("binary_format", fileFormat.Options.JsonBinaryFormat); err != nil {
			return err
		}
		if err := d.Set("trim_space", fileFormat.Options.JsonTrimSpace); err != nil {
			return err
		}
		nullIf := []string{}
		for _, s := range *fileFormat.Options.JsonNullIf {
			nullIf = append(nullIf, s.S)
		}
		if err := d.Set("null_if", nullIf); err != nil {
			return err
		}
		if err := d.Set("file_extension", fileFormat.Options.JsonFileExtension); err != nil {
			return err
		}
		if err := d.Set("enable_octal", fileFormat.Options.JsonEnableOctal); err != nil {
			return err
		}
		if err := d.Set("allow_duplicate", fileFormat.Options.JsonAllowDuplicate); err != nil {
			return err
		}
		if err := d.Set("strip_outer_array", fileFormat.Options.JsonStripOuterArray); err != nil {
			return err
		}
		if err := d.Set("strip_null_values", fileFormat.Options.JsonStripNullValues); err != nil {
			return err
		}
		if err := d.Set("replace_invalid_characters", fileFormat.Options.JsonReplaceInvalidCharacters); err != nil {
			return err
		}
		if err := d.Set("ignore_utf8_errors", fileFormat.Options.JsonIgnoreUtf8Errors); err != nil {
			return err
		}
		if err := d.Set("skip_byte_order_mark", fileFormat.Options.JsonSkipByteOrderMark); err != nil {
			return err
		}
	case sdk.FileFormatTypeAvro:
		if err := d.Set("compression", fileFormat.Options.AvroCompression); err != nil {
			return err
		}
		if err := d.Set("trim_space", fileFormat.Options.AvroTrimSpace); err != nil {
			return err
		}
		nullIf := []string{}
		for _, s := range *fileFormat.Options.AvroNullIf {
			nullIf = append(nullIf, s.S)
		}
		if err := d.Set("null_if", nullIf); err != nil {
			return err
		}
	case sdk.FileFormatTypeOrc:
		if err := d.Set("trim_space", fileFormat.Options.OrcTrimSpace); err != nil {
			return err
		}
		nullIf := []string{}
		for _, s := range *fileFormat.Options.OrcNullIf {
			nullIf = append(nullIf, s.S)
		}
		if err := d.Set("null_if", nullIf); err != nil {
			return err
		}
	case sdk.FileFormatTypeParquet:
		if err := d.Set("compression", fileFormat.Options.ParquetCompression); err != nil {
			return err
		}
		if err := d.Set("binary_as_text", fileFormat.Options.ParquetBinaryAsText); err != nil {
			return err
		}
		if err := d.Set("trim_space", fileFormat.Options.ParquetTrimSpace); err != nil {
			return err
		}
		nullIf := []string{}
		for _, s := range *fileFormat.Options.ParquetNullIf {
			nullIf = append(nullIf, s.S)
		}
		if err := d.Set("null_if", nullIf); err != nil {
			return err
		}
	case sdk.FileFormatTypeXml:
		if err := d.Set("compression", fileFormat.Options.XmlCompression); err != nil {
			return err
		}
		if err := d.Set("ignore_utf8_errors", fileFormat.Options.XmlIgnoreUtf8Errors); err != nil {
			return err
		}
		if err := d.Set("preserve_space", fileFormat.Options.XmlPreserveSpace); err != nil {
			return err
		}
		if err := d.Set("strip_outer_element", fileFormat.Options.XmlStripOuterElement); err != nil {
			return err
		}
		if err := d.Set("disable_snowflake_data", fileFormat.Options.XmlDisableSnowflakeData); err != nil {
			return err
		}
		if err := d.Set("disable_auto_convert", fileFormat.Options.XmlDisableAutoConvert); err != nil {
			return err
		}
		if err := d.Set("skip_byte_order_mark", fileFormat.Options.XmlSkipByteOrderMark); err != nil {
			return err
		}
		// Terraform doesn't like it when computed fields aren't set.
		if err := d.Set("null_if", []string{}); err != nil {
			return err
		}
	}

	if err := d.Set("comment", fileFormat.Comment); err != nil {
		return err
	}
	return nil
}

// UpdateFileFormat implements schema.UpdateFunc.
func UpdateFileFormat(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	fileFormatID, err := fileFormatIDFromString(d.Id())
	if err != nil {
		return err
	}
	id := sdk.NewSchemaObjectIdentifier(fileFormatID.DatabaseName, fileFormatID.SchemaName, fileFormatID.FileFormatName)

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), d.Get("name").(string))
		err := client.FileFormats.Alter(ctx, id, &sdk.AlterFileFormatOptions{
			Rename: &sdk.AlterFileFormatRenameOptions{
				NewName: newId,
			},
		})
		if err != nil {
			return fmt.Errorf("error renaming file format: %w", err)
		}
		id = newId
	}

	opts := sdk.AlterFileFormatOptions{}

	switch d.Get("format_type") {
	case sdk.FileFormatTypeCsv:
		if d.HasChange("compression") {
			v := sdk.CsvCompression(d.Get("compression").(string))
			opts.Set.CsvCompression = &v
		}
		if d.HasChange("record_delimiter") {
			v := d.Get("record_delimiter").(string)
			opts.Set.CsvRecordDelimiter = &v
		}
		if d.HasChange("field_delimiter") {
			v := d.Get("field_delimiter").(string)
			opts.Set.CsvFieldDelimiter = &v
		}
		if d.HasChange("file_extension") {
			v := d.Get("file_extension").(string)
			opts.Set.CsvFileExtension = &v
		}
		if d.HasChange("skip_header") {
			v := d.Get("skip_header").(int)
			opts.Set.CsvSkipHeader = &v
		}
		if d.HasChange("skip_blank_lines") {
			v := d.Get("skip_blank_lines").(bool)
			opts.Set.CsvSkipBlankLines = &v
		}
		if d.HasChange("date_format") {
			v := d.Get("date_format").(string)
			opts.Set.CsvDateFormat = &v
		}
		if d.HasChange("time_format") {
			v := d.Get("time_format").(string)
			opts.Set.CsvTimeFormat = &v
		}
		if d.HasChange("timestamp_format") {
			v := d.Get("timestamp_format").(string)
			opts.Set.CsvTimestampFormat = &v
		}
		if d.HasChange("binary_format") {
			v := sdk.BinaryFormat(d.Get("binary_format").(string))
			opts.Set.CsvBinaryFormat = &v
		}
		if d.HasChange("escape") {
			v := d.Get("escape").(string)
			opts.Set.CsvEscape = &v
		}
		if d.HasChange("escape_unenclosed_field") {
			v := d.Get("escape_unenclosed_field").(string)
			opts.Set.CsvEscapeUnenclosedField = &v
		}
		if d.HasChange("trim_space") {
			v := d.Get("trim_space").(bool)
			opts.Set.CsvTrimSpace = &v
		}
		if d.HasChange("field_optionally_enclosed_by") {
			v := d.Get("field_optionally_enclosed_by").(string)
			opts.Set.CsvFieldOptionallyEnclosedBy = &v
		}
		if d.HasChange("null_if") {
			nullIf := []sdk.NullString{}
			for _, s := range d.Get("null_if").([]interface{}) {
				if s == nil {
					s = ""
				} else {
					s = s.(string)
				}
				nullIf = append(nullIf, sdk.NullString{S: s.(string)})
			}
			opts.Set.CsvNullIf = &nullIf
		}
		if d.HasChange("error_on_column_count_mismatch") {
			v := d.Get("error_on_column_count_mismatch").(bool)
			opts.Set.CsvErrorOnColumnCountMismatch = &v
		}
		if d.HasChange("replace_invalid_characters") {
			v := d.Get("replace_invalid_characters").(bool)
			opts.Set.CsvReplaceInvalidCharacters = &v
		}
		if d.HasChange("empty_field_as_null") {
			v := d.Get("empty_field_as_null").(bool)
			opts.Set.CsvEmptyFieldAsNull = &v
		}
		if d.HasChange("skip_byte_order_mark") {
			v := d.Get("skip_byte_order_mark").(bool)
			opts.Set.CsvSkipByteOrderMark = &v
		}
		if d.HasChange("encoding") {
			v := sdk.CsvEncoding(d.Get("encoding").(string))
			opts.Set.CsvEncoding = &v
		}
	case sdk.FileFormatTypeJson:
		if d.HasChange("compression") {
			comp := sdk.JsonCompression(d.Get("compression").(string))
			opts.Set.JsonCompression = &comp
		}
		if d.HasChange("date_format") {
			v := d.Get("date_format").(string)
			opts.Set.JsonDateFormat = &v
		}
		if d.HasChange("time_format") {
			v := d.Get("time_format").(string)
			opts.Set.JsonTimeFormat = &v
		}
		if d.HasChange("timestamp_format") {
			v := d.Get("timestamp_format").(string)
			opts.Set.JsonTimestampFormat = &v
		}
		if d.HasChange("binary_format") {
			v := sdk.BinaryFormat(d.Get("binary_format").(string))
			opts.Set.JsonBinaryFormat = &v
		}
		if d.HasChange("trim_space") {
			v := d.Get("trim_space").(bool)
			opts.Set.JsonTrimSpace = &v
		}
		if d.HasChange("null_if") {
			nullIf := []sdk.NullString{}
			for _, s := range d.Get("null_if").([]interface{}) {
				if s == nil {
					s = ""
				} else {
					s = s.(string)
				}
				nullIf = append(nullIf, sdk.NullString{S: s.(string)})
			}
			opts.Set.JsonNullIf = &nullIf
		}
		if d.HasChange("file_extension") {
			v := d.Get("file_extension").(string)
			opts.Set.JsonFileExtension = &v
		}
		if d.HasChange("enable_octal") {
			v := d.Get("enable_octal").(bool)
			opts.Set.JsonEnableOctal = &v
		}
		if d.HasChange("allow_duplicate") {
			v := d.Get("allow_duplicate").(bool)
			opts.Set.JsonAllowDuplicate = &v
		}
		if d.HasChange("strip_outer_array") {
			v := d.Get("strip_outer_array").(bool)
			opts.Set.JsonStripOuterArray = &v
		}
		if d.HasChange("strip_null_values") {
			v := d.Get("strip_null_values").(bool)
			opts.Set.JsonStripNullValues = &v
		}
		if d.HasChange("replace_invalid_characters") {
			v := d.Get("replace_invalid_characters").(bool)
			opts.Set.JsonReplaceInvalidCharacters = &v
		}
		if d.HasChange("ignore_utf8_errors") {
			v := d.Get("ignore_utf8_errors").(bool)
			opts.Set.JsonIgnoreUtf8Errors = &v
		}
		if d.HasChange("skip_byte_order_mark") {
			v := d.Get("skip_byte_order_mark").(bool)
			opts.Set.JsonSkipByteOrderMark = &v
		}
	case sdk.FileFormatTypeAvro:
		if d.HasChange("compression") {
			comp := sdk.AvroCompression(d.Get("compression").(string))
			opts.Set.AvroCompression = &comp
		}
		if d.HasChange("trim_space") {
			v := d.Get("trim_space").(bool)
			opts.Set.AvroTrimSpace = &v
		}
		if d.HasChange("null_if") {
			nullIf := []sdk.NullString{}
			for _, s := range d.Get("null_if").([]interface{}) {
				if s == nil {
					s = ""
				} else {
					s = s.(string)
				}
				nullIf = append(nullIf, sdk.NullString{S: s.(string)})
			}
			opts.Set.AvroNullIf = &nullIf
		}
	case sdk.FileFormatTypeOrc:
		if d.HasChange("trim_space") {
			v := d.Get("trim_space").(bool)
			opts.Set.OrcTrimSpace = &v
		}
		if d.HasChange("null_if") {
			nullIf := []sdk.NullString{}
			for _, s := range d.Get("null_if").([]interface{}) {
				if s == nil {
					s = ""
				} else {
					s = s.(string)
				}
				nullIf = append(nullIf, sdk.NullString{S: s.(string)})
			}
			opts.Set.OrcNullIf = &nullIf
		}
	case sdk.FileFormatTypeParquet:
		if d.HasChange("compression") {
			comp := sdk.ParquetCompression(d.Get("compression").(string))
			opts.Set.ParquetCompression = &comp
		}
		if d.HasChange("binary_as_text") {
			v := d.Get("binary_as_text").(bool)
			opts.Set.ParquetBinaryAsText = &v
		}
		if d.HasChange("trim_space") {
			v := d.Get("trim_space").(bool)
			opts.Set.ParquetTrimSpace = &v
		}
		if d.HasChange("null_if") {
			nullIf := []sdk.NullString{}
			for _, s := range d.Get("null_if").([]interface{}) {
				if s == nil {
					s = ""
				} else {
					s = s.(string)
				}
				nullIf = append(nullIf, sdk.NullString{S: s.(string)})
			}
			opts.Set.ParquetNullIf = &nullIf
		}
	case sdk.FileFormatTypeXml:
		if d.HasChange("compression") {
			comp := sdk.XmlCompression(d.Get("compression").(string))
			opts.Set.XmlCompression = &comp
		}
		if d.HasChange("ignore_utf8_errors") {
			v := d.Get("ignore_utf8_errors").(bool)
			opts.Set.XmlIgnoreUtf8Errors = &v
		}
		if d.HasChange("preserve_space") {
			v := d.Get("preserve_space").(bool)
			opts.Set.XmlPreserveSpace = &v
		}
		if d.HasChange("strip_outer_element") {
			v := d.Get("strip_outer_element").(bool)
			opts.Set.XmlStripOuterElement = &v
		}
		if d.HasChange("disable_snowflake_data") {
			v := d.Get("disable_snowflake_data").(bool)
			opts.Set.XmlDisableSnowflakeData = &v
		}
		if d.HasChange("disable_auto_convert") {
			v := d.Get("disable_auto_convert").(bool)
			opts.Set.XmlDisableAutoConvert = &v
		}
		if d.HasChange("skip_byte_order_mark") {
			v := d.Get("skip_byte_order_mark").(bool)
			opts.Set.XmlSkipByteOrderMark = &v
		}
	}

	err = client.FileFormats.Alter(ctx, id, &opts)

	if err != nil {
		return err
	}

	return ReadFileFormat(d, meta)
}

// DeleteFileFormat implements schema.DeleteFunc.
func DeleteFileFormat(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	fileFormatID, err := fileFormatIDFromString(d.Id())
	if err != nil {
		return err
	}
	id := sdk.NewSchemaObjectIdentifier(fileFormatID.DatabaseName, fileFormatID.SchemaName, fileFormatID.FileFormatName)

	err = client.FileFormats.Drop(ctx, id, nil)
	if err != nil {
		return fmt.Errorf("error while deleting file format: %w", err)
	}

	d.SetId("")

	return nil
}

func fileFormatIDFromString(stringID string) (*fileFormatID, error) {
	reader := csv.NewReader(strings.NewReader(stringID))
	reader.Comma = fileFormatIDDelimiter
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not CSV compatible")
	}

	if len(lines) != 1 {
		return nil, fmt.Errorf("1 line at a time")
	}
	if len(lines[0]) != 3 {
		return nil, fmt.Errorf("4 fields allowed")
	}

	return &fileFormatID{
		DatabaseName:   lines[0][0],
		SchemaName:     lines[0][1],
		FileFormatName: lines[0][2],
	}, nil
}

func getFormatTypeOption(d *schema.ResourceData, formatType, formatTypeOption string) (interface{}, bool, error) {
	validFormatTypeOptions := formatTypeOptions[formatType]
	if v, ok := d.GetOk(formatTypeOption); ok {
		if err := validateFormatTypeOptions(formatType, formatTypeOption, validFormatTypeOptions); err != nil {
			return nil, true, err
		}
		return v, true, nil
	}
	return nil, false, nil
}

func validateFormatTypeOptions(formatType, formatTypeOption string, validFormatTypeOptions []string) error {
	for _, f := range validFormatTypeOptions {
		if f == formatTypeOption {
			return nil
		}
	}
	return fmt.Errorf("%v is an invalid format type option for format type %v", formatTypeOption, formatType)
}
