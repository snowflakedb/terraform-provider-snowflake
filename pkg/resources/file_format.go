package resources

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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
		"parse_header",
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
	"parse_header": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Boolean that specifies whether to use the first row headers in the data files to determine column names.",
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
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
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

func FileFormat() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		func(resourceId string) (sdk.SchemaObjectIdentifier, error) {
			id, err := fileFormatIDFromString(resourceId)
			if err != nil {
				return sdk.SchemaObjectIdentifier{}, err
			}
			return sdk.NewSchemaObjectIdentifier(id.DatabaseName, id.SchemaName, id.FileFormatName), nil
		},
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.FileFormats.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.FileFormatResource), TrackingCreateWrapper(resources.FileFormat, CreateFileFormat)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.FileFormatResource), TrackingReadWrapper(resources.FileFormat, ReadFileFormat)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.FileFormatResource), TrackingUpdateWrapper(resources.FileFormat, UpdateFileFormat)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.FileFormatResource), TrackingDeleteWrapper(resources.FileFormat, deleteFunc)),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.FileFormat, customdiff.All(
			ComputedIfAnyAttributeChanged(fileFormatSchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Schema: fileFormatSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateFileFormat implements schema.CreateFunc.
func CreateFileFormat(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	dbName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	fileFormatName := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(dbName, schemaName, fileFormatName)

	opts := sdk.CreateFileFormatOptions{
		Type:                  sdk.FileFormatType(d.Get("format_type").(string)),
		FileFormatTypeOptions: sdk.FileFormatTypeOptions{},
	}

	switch opts.Type {
	case sdk.FileFormatTypeCSV:
		if v, ok := d.GetOk("compression"); ok {
			comp := sdk.CSVCompression(v.(string))
			opts.CSVCompression = &comp
		}
		if v, ok := d.GetOk("record_delimiter"); ok {
			opts.CSVRecordDelimiter = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("field_delimiter"); ok {
			opts.CSVFieldDelimiter = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("file_extension"); ok {
			opts.CSVFileExtension = sdk.String(v.(string))
		}
		opts.CSVParseHeader = sdk.Bool(d.Get("parse_header").(bool))
		if v, ok := d.GetOk("skip_header"); ok {
			opts.CSVSkipHeader = sdk.Int(v.(int))
		}
		opts.CSVSkipBlankLines = sdk.Bool(d.Get("skip_blank_lines").(bool))
		if v, ok := d.GetOk("date_format"); ok {
			opts.CSVDateFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("time_format"); ok {
			opts.CSVTimeFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("timestamp_format"); ok {
			opts.CSVTimestampFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("binary_format"); ok {
			bf := sdk.BinaryFormat(v.(string))
			opts.CSVBinaryFormat = &bf
		}
		if v, ok := d.GetOk("escape"); ok {
			opts.CSVEscape = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("escape_unenclosed_field"); ok {
			opts.CSVEscapeUnenclosedField = sdk.String(v.(string))
		}
		opts.CSVTrimSpace = sdk.Bool(d.Get("trim_space").(bool))
		if v, ok := d.GetOk("field_optionally_enclosed_by"); ok {
			opts.CSVFieldOptionallyEnclosedBy = sdk.String(v.(string))
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
			opts.CSVNullIf = &nullIf
		}
		opts.CSVErrorOnColumnCountMismatch = sdk.Bool(d.Get("error_on_column_count_mismatch").(bool))
		opts.CSVReplaceInvalidCharacters = sdk.Bool(d.Get("replace_invalid_characters").(bool))
		opts.CSVEmptyFieldAsNull = sdk.Bool(d.Get("empty_field_as_null").(bool))
		opts.CSVSkipByteOrderMark = sdk.Bool(d.Get("skip_byte_order_mark").(bool))
		if v, ok := d.GetOk("encoding"); ok {
			enc := sdk.CSVEncoding(v.(string))
			opts.CSVEncoding = &enc
		}
	case sdk.FileFormatTypeJSON:
		if v, ok := d.GetOk("compression"); ok {
			comp := sdk.JSONCompression(v.(string))
			opts.JSONCompression = &comp
		}
		if v, ok := d.GetOk("date_format"); ok {
			opts.JSONDateFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("time_format"); ok {
			opts.JSONTimeFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("timestamp_format"); ok {
			opts.JSONTimestampFormat = sdk.String(v.(string))
		}
		if v, ok := d.GetOk("binary_format"); ok {
			bf := sdk.BinaryFormat(v.(string))
			opts.JSONBinaryFormat = &bf
		}
		opts.JSONTrimSpace = sdk.Bool(d.Get("trim_space").(bool))
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
			opts.JSONNullIf = nullIf
		}
		if v, ok := d.GetOk("file_extension"); ok {
			opts.JSONFileExtension = sdk.String(v.(string))
		}
		opts.JSONEnableOctal = sdk.Bool(d.Get("enable_octal").(bool))
		opts.JSONAllowDuplicate = sdk.Bool(d.Get("allow_duplicate").(bool))
		opts.JSONStripOuterArray = sdk.Bool(d.Get("strip_outer_array").(bool))
		opts.JSONStripNullValues = sdk.Bool(d.Get("strip_null_values").(bool))
		opts.JSONReplaceInvalidCharacters = sdk.Bool(d.Get("replace_invalid_characters").(bool))
		opts.JSONIgnoreUTF8Errors = sdk.Bool(d.Get("ignore_utf8_errors").(bool))
		opts.JSONSkipByteOrderMark = sdk.Bool(d.Get("skip_byte_order_mark").(bool))
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
	case sdk.FileFormatTypeORC:
		opts.ORCTrimSpace = sdk.Bool(d.Get("trim_space").(bool))
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
			opts.ORCNullIf = &nullIf
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
	case sdk.FileFormatTypeXML:
		if v, ok := d.GetOk("compression"); ok {
			comp := sdk.XMLCompression(v.(string))
			opts.XMLCompression = &comp
		}
		opts.XMLIgnoreUTF8Errors = sdk.Bool(d.Get("ignore_utf8_errors").(bool))
		opts.XMLPreserveSpace = sdk.Bool(d.Get("preserve_space").(bool))
		opts.XMLStripOuterElement = sdk.Bool(d.Get("strip_outer_element").(bool))
		opts.XMLDisableSnowflakeData = sdk.Bool(d.Get("disable_snowflake_data").(bool))
		opts.XMLDisableAutoConvert = sdk.Bool(d.Get("disable_auto_convert").(bool))
		opts.XMLSkipByteOrderMark = sdk.Bool(d.Get("skip_byte_order_mark").(bool))
	}

	if v, ok := d.GetOk("comment"); ok {
		opts.Comment = sdk.String(v.(string))
	}

	err := client.FileFormats.Create(ctx, id, &opts)
	if err != nil {
		return diag.FromErr(err)
	}

	fileFormatID := &fileFormatID{
		DatabaseName:   dbName,
		SchemaName:     schemaName,
		FileFormatName: fileFormatName,
	}
	dataIDInput, err := fileFormatID.String()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dataIDInput)

	return ReadFileFormat(ctx, d, meta)
}

// ReadFileFormat implements schema.ReadFunc.
func ReadFileFormat(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	fileFormatID, err := fileFormatIDFromString(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	id := sdk.NewSchemaObjectIdentifier(fileFormatID.DatabaseName, fileFormatID.SchemaName, fileFormatID.FileFormatName)

	fileFormat, err := client.FileFormats.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query file format. Marking the resource as removed.",
					Detail:   fmt.Sprintf("File format id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", fileFormat.Name.Name()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database", fileFormat.Name.DatabaseName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("schema", fileFormat.Name.SchemaName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("format_type", fileFormat.Type); err != nil {
		return diag.FromErr(err)
	}

	switch fileFormat.Type {
	case sdk.FileFormatTypeCSV:
		if err := d.Set("compression", fileFormat.Options.CSVCompression); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("record_delimiter", fileFormat.Options.CSVRecordDelimiter); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("field_delimiter", fileFormat.Options.CSVFieldDelimiter); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("file_extension", fileFormat.Options.CSVFileExtension); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("parse_header", fileFormat.Options.CSVParseHeader); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("skip_header", fileFormat.Options.CSVSkipHeader); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("skip_blank_lines", fileFormat.Options.CSVSkipBlankLines); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("date_format", fileFormat.Options.CSVDateFormat); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("time_format", fileFormat.Options.CSVTimeFormat); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("timestamp_format", fileFormat.Options.CSVTimestampFormat); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("binary_format", fileFormat.Options.CSVBinaryFormat); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("escape", fileFormat.Options.CSVEscape); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("escape_unenclosed_field", fileFormat.Options.CSVEscapeUnenclosedField); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("trim_space", fileFormat.Options.CSVTrimSpace); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("field_optionally_enclosed_by", fileFormat.Options.CSVFieldOptionallyEnclosedBy); err != nil {
			return diag.FromErr(err)
		}
		nullIf := []string{}
		for _, s := range *fileFormat.Options.CSVNullIf {
			nullIf = append(nullIf, s.S)
		}
		if err := d.Set("null_if", nullIf); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("error_on_column_count_mismatch", fileFormat.Options.CSVErrorOnColumnCountMismatch); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("replace_invalid_characters", fileFormat.Options.CSVReplaceInvalidCharacters); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("empty_field_as_null", fileFormat.Options.CSVEmptyFieldAsNull); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("skip_byte_order_mark", fileFormat.Options.CSVSkipByteOrderMark); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("encoding", fileFormat.Options.CSVEncoding); err != nil {
			return diag.FromErr(err)
		}
	case sdk.FileFormatTypeJSON:
		if err := d.Set("compression", fileFormat.Options.JSONCompression); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("date_format", fileFormat.Options.JSONDateFormat); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("time_format", fileFormat.Options.JSONTimeFormat); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("timestamp_format", fileFormat.Options.JSONTimestampFormat); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("binary_format", fileFormat.Options.JSONBinaryFormat); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("trim_space", fileFormat.Options.JSONTrimSpace); err != nil {
			return diag.FromErr(err)
		}
		nullIf := []string{}
		for _, s := range fileFormat.Options.JSONNullIf {
			nullIf = append(nullIf, s.S)
		}
		if err := d.Set("null_if", nullIf); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("file_extension", fileFormat.Options.JSONFileExtension); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("enable_octal", fileFormat.Options.JSONEnableOctal); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("allow_duplicate", fileFormat.Options.JSONAllowDuplicate); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("strip_outer_array", fileFormat.Options.JSONStripOuterArray); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("strip_null_values", fileFormat.Options.JSONStripNullValues); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("replace_invalid_characters", fileFormat.Options.JSONReplaceInvalidCharacters); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("ignore_utf8_errors", fileFormat.Options.JSONIgnoreUTF8Errors); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("skip_byte_order_mark", fileFormat.Options.JSONSkipByteOrderMark); err != nil {
			return diag.FromErr(err)
		}
	case sdk.FileFormatTypeAvro:
		if err := d.Set("compression", fileFormat.Options.AvroCompression); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("trim_space", fileFormat.Options.AvroTrimSpace); err != nil {
			return diag.FromErr(err)
		}
		nullIf := []string{}
		for _, s := range *fileFormat.Options.AvroNullIf {
			nullIf = append(nullIf, s.S)
		}
		if err := d.Set("null_if", nullIf); err != nil {
			return diag.FromErr(err)
		}
	case sdk.FileFormatTypeORC:
		if err := d.Set("trim_space", fileFormat.Options.ORCTrimSpace); err != nil {
			return diag.FromErr(err)
		}
		nullIf := []string{}
		for _, s := range *fileFormat.Options.ORCNullIf {
			nullIf = append(nullIf, s.S)
		}
		if err := d.Set("null_if", nullIf); err != nil {
			return diag.FromErr(err)
		}
	case sdk.FileFormatTypeParquet:
		if err := d.Set("compression", fileFormat.Options.ParquetCompression); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("binary_as_text", fileFormat.Options.ParquetBinaryAsText); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("trim_space", fileFormat.Options.ParquetTrimSpace); err != nil {
			return diag.FromErr(err)
		}
		nullIf := []string{}
		for _, s := range *fileFormat.Options.ParquetNullIf {
			nullIf = append(nullIf, s.S)
		}
		if err := d.Set("null_if", nullIf); err != nil {
			return diag.FromErr(err)
		}
	case sdk.FileFormatTypeXML:
		if err := d.Set("compression", fileFormat.Options.XMLCompression); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("ignore_utf8_errors", fileFormat.Options.XMLIgnoreUTF8Errors); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("preserve_space", fileFormat.Options.XMLPreserveSpace); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("strip_outer_element", fileFormat.Options.XMLStripOuterElement); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("disable_snowflake_data", fileFormat.Options.XMLDisableSnowflakeData); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("disable_auto_convert", fileFormat.Options.XMLDisableAutoConvert); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("skip_byte_order_mark", fileFormat.Options.XMLSkipByteOrderMark); err != nil {
			return diag.FromErr(err)
		}
		// Terraform doesn't like it when computed fields aren't set.
		if err := d.Set("null_if", []string{}); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("comment", fileFormat.Comment); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// UpdateFileFormat implements schema.UpdateFunc.
func UpdateFileFormat(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	fileFormatID, err := fileFormatIDFromString(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	id := sdk.NewSchemaObjectIdentifier(fileFormatID.DatabaseName, fileFormatID.SchemaName, fileFormatID.FileFormatName)

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		err := client.FileFormats.Alter(ctx, id, &sdk.AlterFileFormatOptions{
			Rename: &sdk.AlterFileFormatRenameOptions{
				NewName: newId,
			},
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming file format: %w", err))
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	runSet := false
	opts := sdk.AlterFileFormatOptions{Set: &sdk.FileFormatTypeOptions{}}

	switch sdk.FileFormatType(d.Get("format_type").(string)) {
	case sdk.FileFormatTypeCSV:
		if d.HasChange("compression") {
			v := sdk.CSVCompression(d.Get("compression").(string))
			opts.Set.CSVCompression = &v
			runSet = true
		}
		if d.HasChange("record_delimiter") {
			v := d.Get("record_delimiter").(string)
			opts.Set.CSVRecordDelimiter = &v
			runSet = true
		}
		if d.HasChange("field_delimiter") {
			v := d.Get("field_delimiter").(string)
			opts.Set.CSVFieldDelimiter = &v
			runSet = true
		}
		if d.HasChange("file_extension") {
			v := d.Get("file_extension").(string)
			opts.Set.CSVFileExtension = &v
			runSet = true
		}
		if d.HasChange("parse_header") {
			v := d.Get("parse_header").(bool)
			opts.Set.CSVParseHeader = &v
			runSet = true
		}
		if d.HasChange("skip_header") {
			v := d.Get("skip_header").(int)
			opts.Set.CSVSkipHeader = &v
			runSet = true
		}
		if d.HasChange("skip_blank_lines") {
			v := d.Get("skip_blank_lines").(bool)
			opts.Set.CSVSkipBlankLines = &v
			runSet = true
		}
		if d.HasChange("date_format") {
			v := d.Get("date_format").(string)
			opts.Set.CSVDateFormat = &v
			runSet = true
		}
		if d.HasChange("time_format") {
			v := d.Get("time_format").(string)
			opts.Set.CSVTimeFormat = &v
			runSet = true
		}
		if d.HasChange("timestamp_format") {
			v := d.Get("timestamp_format").(string)
			opts.Set.CSVTimestampFormat = &v
			runSet = true
		}
		if d.HasChange("binary_format") {
			v := sdk.BinaryFormat(d.Get("binary_format").(string))
			opts.Set.CSVBinaryFormat = &v
			runSet = true
		}
		if d.HasChange("escape") {
			v := d.Get("escape").(string)
			opts.Set.CSVEscape = &v
			runSet = true
		}
		if d.HasChange("escape_unenclosed_field") {
			v := d.Get("escape_unenclosed_field").(string)
			opts.Set.CSVEscapeUnenclosedField = &v
			runSet = true
		}
		if d.HasChange("trim_space") {
			v := d.Get("trim_space").(bool)
			opts.Set.CSVTrimSpace = &v
			runSet = true
		}
		if d.HasChange("field_optionally_enclosed_by") {
			v := d.Get("field_optionally_enclosed_by").(string)
			opts.Set.CSVFieldOptionallyEnclosedBy = &v
			runSet = true
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
			opts.Set.CSVNullIf = &nullIf
			runSet = true
		}
		if d.HasChange("error_on_column_count_mismatch") {
			v := d.Get("error_on_column_count_mismatch").(bool)
			opts.Set.CSVErrorOnColumnCountMismatch = &v
			runSet = true
		}
		if d.HasChange("replace_invalid_characters") {
			v := d.Get("replace_invalid_characters").(bool)
			opts.Set.CSVReplaceInvalidCharacters = &v
			runSet = true
		}
		if d.HasChange("empty_field_as_null") {
			v := d.Get("empty_field_as_null").(bool)
			opts.Set.CSVEmptyFieldAsNull = &v
			runSet = true
		}
		if d.HasChange("skip_byte_order_mark") {
			v := d.Get("skip_byte_order_mark").(bool)
			opts.Set.CSVSkipByteOrderMark = &v
			runSet = true
		}
		if d.HasChange("encoding") {
			v := sdk.CSVEncoding(d.Get("encoding").(string))
			opts.Set.CSVEncoding = &v
			runSet = true
		}
	case sdk.FileFormatTypeJSON:
		if d.HasChange("compression") {
			comp := sdk.JSONCompression(d.Get("compression").(string))
			opts.Set.JSONCompression = &comp
			runSet = true
		}
		if d.HasChange("date_format") {
			v := d.Get("date_format").(string)
			opts.Set.JSONDateFormat = &v
			runSet = true
		}
		if d.HasChange("time_format") {
			v := d.Get("time_format").(string)
			opts.Set.JSONTimeFormat = &v
			runSet = true
		}
		if d.HasChange("timestamp_format") {
			v := d.Get("timestamp_format").(string)
			opts.Set.JSONTimestampFormat = &v
			runSet = true
		}
		if d.HasChange("binary_format") {
			v := sdk.BinaryFormat(d.Get("binary_format").(string))
			opts.Set.JSONBinaryFormat = &v
			runSet = true
		}
		if d.HasChange("trim_space") {
			v := d.Get("trim_space").(bool)
			opts.Set.JSONTrimSpace = &v
			runSet = true
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
			opts.Set.JSONNullIf = nullIf
			runSet = true
		}
		if d.HasChange("file_extension") {
			v := d.Get("file_extension").(string)
			opts.Set.JSONFileExtension = &v
			runSet = true
		}
		if d.HasChange("enable_octal") {
			v := d.Get("enable_octal").(bool)
			opts.Set.JSONEnableOctal = &v
			runSet = true
		}
		if d.HasChange("allow_duplicate") {
			v := d.Get("allow_duplicate").(bool)
			opts.Set.JSONAllowDuplicate = &v
			runSet = true
		}
		if d.HasChange("strip_outer_array") {
			v := d.Get("strip_outer_array").(bool)
			opts.Set.JSONStripOuterArray = &v
			runSet = true
		}
		if d.HasChange("strip_null_values") {
			v := d.Get("strip_null_values").(bool)
			opts.Set.JSONStripNullValues = &v
			runSet = true
		}
		if d.HasChange("replace_invalid_characters") {
			v := d.Get("replace_invalid_characters").(bool)
			opts.Set.JSONReplaceInvalidCharacters = &v
			runSet = true
		}
		if d.HasChange("ignore_utf8_errors") {
			v := d.Get("ignore_utf8_errors").(bool)
			opts.Set.JSONIgnoreUTF8Errors = &v
			runSet = true
		}
		if d.HasChange("skip_byte_order_mark") {
			v := d.Get("skip_byte_order_mark").(bool)
			opts.Set.JSONSkipByteOrderMark = &v
			runSet = true
		}
	case sdk.FileFormatTypeAvro:
		if d.HasChange("compression") {
			comp := sdk.AvroCompression(d.Get("compression").(string))
			opts.Set.AvroCompression = &comp
			runSet = true
		}
		if d.HasChange("trim_space") {
			v := d.Get("trim_space").(bool)
			opts.Set.AvroTrimSpace = &v
			runSet = true
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
			runSet = true
		}
	case sdk.FileFormatTypeORC:
		if d.HasChange("trim_space") {
			v := d.Get("trim_space").(bool)
			opts.Set.ORCTrimSpace = &v
			runSet = true
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
			opts.Set.ORCNullIf = &nullIf
			runSet = true
		}
	case sdk.FileFormatTypeParquet:
		if d.HasChange("compression") {
			comp := sdk.ParquetCompression(d.Get("compression").(string))
			opts.Set.ParquetCompression = &comp
			runSet = true
		}
		if d.HasChange("binary_as_text") {
			v := d.Get("binary_as_text").(bool)
			opts.Set.ParquetBinaryAsText = &v
			runSet = true
		}
		if d.HasChange("trim_space") {
			v := d.Get("trim_space").(bool)
			opts.Set.ParquetTrimSpace = &v
			runSet = true
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
			runSet = true
		}
	case sdk.FileFormatTypeXML:
		if d.HasChange("compression") {
			comp := sdk.XMLCompression(d.Get("compression").(string))
			opts.Set.XMLCompression = &comp
			runSet = true
		}
		if d.HasChange("ignore_utf8_errors") {
			v := d.Get("ignore_utf8_errors").(bool)
			opts.Set.XMLIgnoreUTF8Errors = &v
			runSet = true
		}
		if d.HasChange("preserve_space") {
			v := d.Get("preserve_space").(bool)
			opts.Set.XMLPreserveSpace = &v
			runSet = true
		}
		if d.HasChange("strip_outer_element") {
			v := d.Get("strip_outer_element").(bool)
			opts.Set.XMLStripOuterElement = &v
			runSet = true
		}
		if d.HasChange("disable_snowflake_data") {
			v := d.Get("disable_snowflake_data").(bool)
			opts.Set.XMLDisableSnowflakeData = &v
			runSet = true
		}
		if d.HasChange("disable_auto_convert") {
			v := d.Get("disable_auto_convert").(bool)
			opts.Set.XMLDisableAutoConvert = &v
			runSet = true
		}
		if d.HasChange("skip_byte_order_mark") {
			v := d.Get("skip_byte_order_mark").(bool)
			opts.Set.XMLSkipByteOrderMark = &v
			runSet = true
		}
	}

	if d.HasChange("comment") {
		v := d.Get("comment").(string)
		opts.Set.Comment = &v
		runSet = true
	}

	if runSet {
		err = client.FileFormats.Alter(ctx, id, &opts)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadFileFormat(ctx, d, meta)
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
