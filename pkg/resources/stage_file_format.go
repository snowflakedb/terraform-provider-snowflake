package resources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var stageFileFormatSchema = map[string]*schema.Schema{
	"file_format": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Specifies the file format for the stage.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"format_name": {
					Type:             schema.TypeString,
					Optional:         true,
					ExactlyOneOf:     []string{"file_format.0.format_name", "file_format.0.csv"},
					Description:      "Fully qualified name of the file format (e.g., 'database.schema.format_name').",
					DiffSuppressFunc: suppressIdentifierQuoting,
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
				},
				"csv": {
					Type:         schema.TypeList,
					Optional:     true,
					MaxItems:     1,
					ExactlyOneOf: []string{"file_format.0.format_name", "file_format.0.csv"},
					Description:  "CSV file format options.",
					Elem: &schema.Resource{
						Schema: csvFileFormatSchema,
					},
				},
			},
		},
	},
}

var csvFileFormatSchema = map[string]*schema.Schema{
	"compression": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      fmt.Sprintf("Specifies the compression format. Valid values: %s.", possibleValuesListed(sdk.AllCsvCompressions)),
		ValidateDiagFunc: sdkValidation(sdk.ToCsvCompression),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToCsvCompression),
	},
	"record_delimiter": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "One or more singlebyte or multibyte characters that separate records in an input file. Use `NONE` to specify no delimiter.",
	},
	"field_delimiter": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "One or more singlebyte or multibyte characters that separate fields in an input file. Use `NONE` to specify no delimiter.",
	},
	"multi_line": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Boolean that specifies whether to parse CSV files containing multiple records on a single line."),
	},
	"file_extension": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the extension for files unloaded to a stage.",
	},
	"parse_header": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Boolean that specifies whether to use the first row headers in the data files to determine column names."),
		ConflictsWith:    []string{"file_format.0.csv.0.skip_header"},
	},
	"skip_header": {
		Type:          schema.TypeInt,
		Optional:      true,
		ValidateFunc:  validation.IntAtLeast(0),
		Default:       IntDefault,
		Description:   "Number of lines at the start of the file to skip.",
		ConflictsWith: []string{"file_format.0.csv.0.parse_header"},
	},
	"skip_blank_lines": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Boolean that specifies to skip any blank lines encountered in the data files."),
	},
	"date_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Defines the format of date values in the data files. Use `AUTO` to have Snowflake auto-detect the format.",
	},
	"time_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Defines the format of time values in the data files. Use `AUTO` to have Snowflake auto-detect the format.",
	},
	"timestamp_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Defines the format of timestamp values in the data files. Use `AUTO` to have Snowflake auto-detect the format.",
	},
	"binary_format": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      fmt.Sprintf("Defines the encoding format for binary input or output. Valid values: %s.", possibleValuesListed(sdk.AllBinaryFormats)),
		ValidateDiagFunc: sdkValidation(sdk.ToBinaryFormat),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToBinaryFormat),
	},
	"escape": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Single character string used as the escape character for field values. Use `NONE` to specify no escape character. NOTE: This value may be not imported properly from Snowflake. Snowflake returns escaped values.",
	},
	"escape_unenclosed_field": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Single character string used as the escape character for unenclosed field values only. Use `NONE` to specify no escape character. NOTE: This value may be not imported properly from Snowflake. Snowflake returns escaped values.",
	},
	"trim_space": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Boolean that specifies whether to remove white space from fields."),
	},
	"field_optionally_enclosed_by": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Character used to enclose strings. Use `NONE` to specify no enclosure character.",
	},
	"null_if": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "String used to convert to and from SQL NULL.",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"error_on_column_count_mismatch": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Boolean that specifies whether to generate a parsing error if the number of delimited columns in an input file does not match the number of columns in the corresponding table."),
	},
	"replace_invalid_characters": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Boolean that specifies whether to replace invalid UTF-8 characters with the Unicode replacement character."),
	},
	"empty_field_as_null": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Boolean that specifies whether to insert SQL NULL for empty fields in an input file."),
	},
	"skip_byte_order_mark": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		Description:      booleanStringFieldDescription("Boolean that specifies whether to skip the BOM (byte order mark) if present in a data file."),
	},
	"encoding": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      fmt.Sprintf("Specifies the character set of the source data when loading data into a table. Valid values: %s.", possibleValuesListed(sdk.AllCsvEncodings)),
		ValidateDiagFunc: sdkValidation(sdk.ToCsvEncoding),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToCsvEncoding),
	},
}

// parseStageFileFormat parses the stage file format from the resource data to an SDK object.
func parseStageFileFormat(d *schema.ResourceData) (sdk.StageFileFormatRequest, error) {
	if len(d.Get("file_format").([]any)) == 0 {
		return sdk.StageFileFormatRequest{}, nil
	}
	prefix := "file_format.0."
	fileFormatReq := sdk.NewStageFileFormatRequest()

	err := errors.Join(
		schemaObjectIdentifierAttributeCreate(d, prefix+"format_name", &fileFormatReq.FormatName),
		attributeMappedValueCreateBuilderNested(d, prefix+"csv", func(fileFormatOptions *sdk.FileFormatCsvOptions) *sdk.StageFileFormatRequest {
			return fileFormatReq.WithFileFormatOptions(sdk.FileFormatOptions{
				CsvOptions: fileFormatOptions,
			})
		}, parseCsvFileFormatOptions),
	)
	if err != nil {
		return sdk.StageFileFormatRequest{}, err
	}

	return *fileFormatReq, nil
}

// parseCsvFileFormatOptions parses the CSV file format options from the resource data to an SDK object.
func parseCsvFileFormatOptions(d *schema.ResourceData) (*sdk.FileFormatCsvOptions, error) {
	csvOptions := &sdk.FileFormatCsvOptions{}
	prefix := "file_format.0.csv.0."

	err := errors.Join(
		attributeMappedValueCreate(d, prefix+"compression", &csvOptions.Compression, func(v any) (*sdk.CsvCompression, error) {
			c, err := sdk.ToCsvCompression(v.(string))
			return &c, err
		}),
		attributeMappedValueCreate(d, prefix+"record_delimiter", &csvOptions.RecordDelimiter, func(v any) (*sdk.StageFileFormatStringOrNone, error) {
			return parseStageFileFormatStringOrNone(v.(string)), nil
		}),
		attributeMappedValueCreate(d, prefix+"field_delimiter", &csvOptions.FieldDelimiter, func(v any) (*sdk.StageFileFormatStringOrNone, error) {
			return parseStageFileFormatStringOrNone(v.(string)), nil
		}),
		booleanStringAttributeCreate(d, prefix+"multi_line", &csvOptions.MultiLine),
		stringAttributeCreate(d, prefix+"file_extension", &csvOptions.FileExtension),
		booleanStringAttributeCreate(d, prefix+"parse_header", &csvOptions.ParseHeader),
		intAttributeWithSpecialDefaultCreate(d, prefix+"skip_header", &csvOptions.SkipHeader),
		booleanStringAttributeCreate(d, prefix+"skip_blank_lines", &csvOptions.SkipBlankLines),
		attributeMappedValueCreate(d, prefix+"date_format", &csvOptions.DateFormat, func(v any) (*sdk.StageFileFormatStringOrAuto, error) {
			return parseStageFileFormatStringOrAuto(v.(string)), nil
		}),
		attributeMappedValueCreate(d, prefix+"time_format", &csvOptions.TimeFormat, func(v any) (*sdk.StageFileFormatStringOrAuto, error) {
			return parseStageFileFormatStringOrAuto(v.(string)), nil
		}),
		attributeMappedValueCreate(d, prefix+"timestamp_format", &csvOptions.TimestampFormat, func(v any) (*sdk.StageFileFormatStringOrAuto, error) {
			return parseStageFileFormatStringOrAuto(v.(string)), nil
		}),
		attributeMappedValueCreate(d, prefix+"binary_format", &csvOptions.BinaryFormat, func(v any) (*sdk.BinaryFormat, error) {
			b, err := sdk.ToBinaryFormat(v.(string))
			return &b, err
		}),
		attributeMappedValueCreate(d, prefix+"escape", &csvOptions.Escape, func(v any) (*sdk.StageFileFormatStringOrNone, error) {
			return parseStageFileFormatStringOrNone(v.(string)), nil
		}),
		attributeMappedValueCreate(d, prefix+"escape_unenclosed_field", &csvOptions.EscapeUnenclosedField, func(v any) (*sdk.StageFileFormatStringOrNone, error) {
			return parseStageFileFormatStringOrNone(v.(string)), nil
		}),
		booleanStringAttributeCreate(d, prefix+"trim_space", &csvOptions.TrimSpace),
		attributeMappedValueCreate(d, prefix+"field_optionally_enclosed_by", &csvOptions.FieldOptionallyEnclosedBy, func(v any) (*sdk.StageFileFormatStringOrNone, error) {
			return parseStageFileFormatStringOrNone(v.(string)), nil
		}),
		booleanStringAttributeCreate(d, prefix+"error_on_column_count_mismatch", &csvOptions.ErrorOnColumnCountMismatch),
		booleanStringAttributeCreate(d, prefix+"replace_invalid_characters", &csvOptions.ReplaceInvalidCharacters),
		booleanStringAttributeCreate(d, prefix+"empty_field_as_null", &csvOptions.EmptyFieldAsNull),
		booleanStringAttributeCreate(d, prefix+"skip_byte_order_mark", &csvOptions.SkipByteOrderMark),
		attributeMappedValueCreate(d, prefix+"encoding", &csvOptions.Encoding, func(v any) (*sdk.CsvEncoding, error) {
			e, err := sdk.ToCsvEncoding(v.(string))
			return &e, err
		}),
		attributeMappedValueCreateBuilder(d, prefix+"null_if", func(nullIf []sdk.NullString) *sdk.FileFormatCsvOptions {
			csvOptions.NullIf = nullIf
			return csvOptions
		}, parseNullIf),
	)
	if err != nil {
		return nil, err
	}

	return csvOptions, nil
}

func parseNullIf(v any) ([]sdk.NullString, error) {
	nullIfList := v.([]any)
	if len(nullIfList) == 0 {
		return nil, nil
	}
	nullIf := make([]sdk.NullString, len(nullIfList))
	for i, s := range nullIfList {
		str := ""
		if s != nil {
			str = s.(string)
		}
		nullIf[i] = sdk.NullString{S: str}
	}
	return nullIf, nil
}

func parseStageFileFormatStringOrNone(v string) *sdk.StageFileFormatStringOrNone {
	if strings.ToUpper(v) == "NONE" {
		return &sdk.StageFileFormatStringOrNone{None: sdk.Bool(true)}
	}
	return &sdk.StageFileFormatStringOrNone{Value: sdk.String(v)}
}

func parseStageFileFormatStringOrAuto(v string) *sdk.StageFileFormatStringOrAuto {
	if strings.ToUpper(v) == "AUTO" {
		return &sdk.StageFileFormatStringOrAuto{Auto: sdk.Bool(true)}
	}
	return &sdk.StageFileFormatStringOrAuto{Value: sdk.String(v)}
}

// stageFileFormatToSchema converts the SDK details to a Terraform schema.
func stageFileFormatToSchema(details *sdk.StageDetails) []map[string]any {
	if details == nil {
		return nil
	}

	if details.FileFormatName != nil {
		return []map[string]any{
			{
				"format_name": details.FileFormatName.FullyQualifiedName(),
			},
		}
	}

	if details.FileFormatCsv != nil {
		csvSchema := stageCsvFileFormatToSchema(details.FileFormatCsv)
		return []map[string]any{
			{
				"csv": []map[string]any{csvSchema},
			},
		}
	}

	return nil
}

// stageCsvFileFormatToSchema converts the SDK details for a CSV file format to a Terraform schema.
func stageCsvFileFormatToSchema(csv *sdk.FileFormatCsv) map[string]any {
	return map[string]any{
		"record_delimiter":               csv.RecordDelimiter,
		"field_delimiter":                csv.FieldDelimiter,
		"file_extension":                 csv.FileExtension,
		"skip_header":                    csv.SkipHeader,
		"parse_header":                   booleanStringFromBool(csv.ParseHeader),
		"date_format":                    csv.DateFormat,
		"time_format":                    csv.TimeFormat,
		"timestamp_format":               csv.TimestampFormat,
		"binary_format":                  csv.BinaryFormat,
		"escape":                         csv.Escape,
		"escape_unenclosed_field":        csv.EscapeUnenclosedField,
		"trim_space":                     booleanStringFromBool(csv.TrimSpace),
		"field_optionally_enclosed_by":   csv.FieldOptionallyEnclosedBy,
		"null_if":                        collections.Map(csv.NullIf, func(v string) any { return v }),
		"compression":                    csv.Compression,
		"error_on_column_count_mismatch": booleanStringFromBool(csv.ErrorOnColumnCountMismatch),
		"skip_blank_lines":               booleanStringFromBool(csv.SkipBlankLines),
		"replace_invalid_characters":     booleanStringFromBool(csv.ReplaceInvalidCharacters),
		"empty_field_as_null":            booleanStringFromBool(csv.EmptyFieldAsNull),
		"skip_byte_order_mark":           booleanStringFromBool(csv.SkipByteOrderMark),
		"encoding":                       csv.Encoding,
		"multi_line":                     booleanStringFromBool(csv.MultiLine),
	}
}

func handleStageFileFormatRead(d *schema.ResourceData, details *sdk.StageDetails) error {
	fileFormatSchema, err := schemas.StageDescribeToSchema(*details)
	if err != nil {
		return err
	}
	fileFormatToCompare := collections.Map(fileFormatSchema["file_format"].([]map[string]any), func(v map[string]any) any {
		return v
	})
	fileFormatToSet := stageFileFormatToSchema(details)
	return handleExternalChangesToObjectInFlatDescribeDeepEqual(d,
		outputMapping{"file_format", "file_format", fileFormatToCompare, fileFormatToSet, nil},
	)
}
