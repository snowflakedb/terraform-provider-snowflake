package resources

import (
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var fileFormatCommonSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the file format; must be unique for the database and schema in which the file format is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the file format."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the file format."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the file format.",
	},
	"type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Specifies the type of the file format. This field is used to detect when the file format type was changed outside of Terraform and to recreate the resource when that happens.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW FILE FORMATS` for this file format.",
		Elem: &schema.Resource{
			Schema: schemas.ShowFileFormatSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// isFileFormatAutoSentinel reports whether v is the special "AUTO" value Snowflake uses to
// mean "detect the format automatically", as opposed to a literal format string.
func isFileFormatAutoSentinel(v string) bool {
	return strings.ToUpper(v) == "AUTO"
}

func fileFormatStringOrAutoMapper(v string) (sdk.StageFileFormatStringOrAutoRequest, error) {
	if isFileFormatAutoSentinel(v) {
		return *sdk.NewStageFileFormatStringOrAutoRequest().WithAuto(true), nil
	}
	return *sdk.NewStageFileFormatStringOrAutoRequest().WithValue(v), nil
}

// isFileFormatNoneSentinel reports whether v is the special "NONE" value Snowflake uses to
// mean "no delimiter/escape/enclosing character", as opposed to a literal character.
func isFileFormatNoneSentinel(v string) bool {
	return strings.ToUpper(v) == "NONE"
}

func fileFormatStringOrNoneMapper(v string) (sdk.StageFileFormatStringOrNoneRequest, error) {
	if isFileFormatNoneSentinel(v) {
		return *sdk.NewStageFileFormatStringOrNoneRequest().WithNone(true), nil
	}
	return *sdk.NewStageFileFormatStringOrNoneRequest().WithValue(v), nil
}

func parseNullIfRequest(v any) (sdk.NullIfListRequest, error) {
	nullIf, err := parseNullIf(v)
	if err != nil {
		return sdk.NullIfListRequest{}, err
	}
	return *sdk.NewNullIfListRequest().WithNullIf(nullIf), nil
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

func csvFileFormatSchema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			ConflictsWith:    []string{prefix + "skip_header"},
		},
		"skip_header": {
			Type:          schema.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IntAtLeast(0),
			Default:       IntDefault,
			Description:   "Number of lines at the start of the file to skip.",
			ConflictsWith: []string{prefix + "parse_header"},
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
}

func jsonFileFormatSchema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"compression": {
			Type:             schema.TypeString,
			Optional:         true,
			Description:      fmt.Sprintf("Specifies the compression format. Valid values: %s.", possibleValuesListed(sdk.AllJsonCompressions)),
			ValidateDiagFunc: sdkValidation(sdk.ToJsonCompression),
			DiffSuppressFunc: NormalizeAndCompare(sdk.ToJsonCompression),
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
		"trim_space": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      booleanStringFieldDescription("Boolean that specifies whether to remove white space from fields."),
		},
		"multi_line": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      booleanStringFieldDescription("Boolean that specifies whether to allow multiple records on a single line."),
		},
		"null_if": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "String used to convert to and from SQL NULL.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"file_extension": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies the extension for files unloaded to a stage.",
		},
		"enable_octal": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      booleanStringFieldDescription("Boolean that enables parsing of octal numbers."),
		},
		"allow_duplicate": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      booleanStringFieldDescription("Boolean that specifies whether to allow duplicate object field names (only the last one will be preserved)."),
		},
		"strip_outer_array": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      booleanStringFieldDescription("Boolean that instructs the JSON parser to remove outer brackets."),
		},
		"strip_null_values": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      booleanStringFieldDescription("Boolean that instructs the JSON parser to remove object fields or array elements containing null values."),
		},
		"replace_invalid_characters": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      booleanStringFieldDescription("Boolean that specifies whether to replace invalid UTF-8 characters with the Unicode replacement character."),
			ConflictsWith:    []string{prefix + "ignore_utf8_errors"},
		},
		"ignore_utf8_errors": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      booleanStringFieldDescription("Boolean that specifies whether UTF-8 encoding errors produce error conditions."),
			ConflictsWith:    []string{prefix + "replace_invalid_characters"},
		},
		"skip_byte_order_mark": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      booleanStringFieldDescription("Boolean that specifies whether to skip the BOM (byte order mark) if present in a data file."),
		},
	}
}
