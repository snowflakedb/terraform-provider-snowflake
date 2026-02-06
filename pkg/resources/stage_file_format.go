package resources

import (
	"fmt"
	"strings"

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
	},
	"escape": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Single character string used as the escape character for field values. Use `NONE` to specify no escape character.",
	},
	"escape_unenclosed_field": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Single character string used as the escape character for unenclosed field values only. Use `NONE` to specify no escape character.",
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
	},
}

func parseStageFileFormat(v any) (sdk.StageFileFormatRequest, error) {
	fileFormatList := v.([]any)
	if len(fileFormatList) == 0 {
		return sdk.StageFileFormatRequest{}, nil
	}
	fileFormatConfig := fileFormatList[0].(map[string]any)
	fileFormatReq := sdk.NewStageFileFormatRequest()

	if formatName, ok := fileFormatConfig["format_name"]; ok && formatName.(string) != "" {
		formatNameStr := formatName.(string)
		id, err := sdk.ParseSchemaObjectIdentifier(formatNameStr)
		if err != nil {
			return sdk.StageFileFormatRequest{}, fmt.Errorf("parsing format_name: %w", err)
		}
		fileFormatReq.WithFormatName(id)
	}

	if csv, ok := fileFormatConfig["csv"]; ok {
		csvList := csv.([]any)
		if len(csvList) > 0 {
			csvOptions, err := parseCsvFileFormatOptions(csvList[0].(map[string]any))
			if err != nil {
				return sdk.StageFileFormatRequest{}, err
			}
			fileFormatReq.WithFileFormatOptions(sdk.FileFormatOptions{
				CsvOptions: csvOptions,
			})
		}
	}

	return *fileFormatReq, nil
}

func parseCsvFileFormatOptions(csvConfig map[string]any) (*sdk.FileFormatCsvOptions, error) {
	csvOptions := &sdk.FileFormatCsvOptions{}

	if v, ok := csvConfig["compression"]; ok && v.(string) != "" {
		compression, err := sdk.ToCsvCompression(v.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing compression: %w", err)
		}
		csvOptions.Compression = &compression
	}

	if v, ok := csvConfig["record_delimiter"]; ok && v.(string) != "" {
		csvOptions.RecordDelimiter = parseStageFileFormatStringOrNone(v.(string))
	}

	if v, ok := csvConfig["field_delimiter"]; ok && v.(string) != "" {
		csvOptions.FieldDelimiter = parseStageFileFormatStringOrNone(v.(string))
	}

	if v, ok := csvConfig["multi_line"]; ok && v.(string) != BooleanDefault && v.(string) != "" {
		multiLine, err := booleanStringToBool(v.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing multi_line: %w", err)
		}
		csvOptions.MultiLine = &multiLine
	}

	if v, ok := csvConfig["file_extension"]; ok && v.(string) != "" {
		fileExtension := v.(string)
		csvOptions.FileExtension = &fileExtension
	}

	if v, ok := csvConfig["parse_header"]; ok && v.(string) != BooleanDefault && v.(string) != "" {
		parseHeader, err := booleanStringToBool(v.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing parse_header: %w", err)
		}
		csvOptions.ParseHeader = &parseHeader
	}

	if v, ok := csvConfig["skip_header"]; ok && v.(int) > IntDefault {
		skipHeader := v.(int)
		csvOptions.SkipHeader = &skipHeader
	}

	if v, ok := csvConfig["skip_blank_lines"]; ok && v.(string) != BooleanDefault && v.(string) != "" {
		skipBlankLines, err := booleanStringToBool(v.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing skip_blank_lines: %w", err)
		}
		csvOptions.SkipBlankLines = &skipBlankLines
	}

	if v, ok := csvConfig["date_format"]; ok && v.(string) != "" {
		csvOptions.DateFormat = parseStageFileFormatStringOrAuto(v.(string))
	}

	if v, ok := csvConfig["time_format"]; ok && v.(string) != "" {
		csvOptions.TimeFormat = parseStageFileFormatStringOrAuto(v.(string))
	}

	if v, ok := csvConfig["timestamp_format"]; ok && v.(string) != "" {
		csvOptions.TimestampFormat = parseStageFileFormatStringOrAuto(v.(string))
	}

	if v, ok := csvConfig["binary_format"]; ok && v.(string) != "" {
		binaryFormat, err := sdk.ToBinaryFormat(v.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing binary_format: %w", err)
		}
		csvOptions.BinaryFormat = &binaryFormat
	}

	if v, ok := csvConfig["escape"]; ok && v.(string) != "" {
		csvOptions.Escape = parseStageFileFormatStringOrNone(v.(string))
	}

	if v, ok := csvConfig["escape_unenclosed_field"]; ok && v.(string) != "" {
		csvOptions.EscapeUnenclosedField = parseStageFileFormatStringOrNone(v.(string))
	}

	if v, ok := csvConfig["trim_space"]; ok && v.(string) != BooleanDefault && v.(string) != "" {
		trimSpace, err := booleanStringToBool(v.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing trim_space: %w", err)
		}
		csvOptions.TrimSpace = &trimSpace
	}

	if v, ok := csvConfig["field_optionally_enclosed_by"]; ok && v.(string) != "" {
		csvOptions.FieldOptionallyEnclosedBy = parseStageFileFormatStringOrNone(v.(string))
	}

	if v, ok := csvConfig["null_if"]; ok {
		nullIfList := v.([]any)
		if len(nullIfList) > 0 {
			nullIf := make([]sdk.NullString, len(nullIfList))
			for i, s := range nullIfList {
				str := ""
				if s != nil {
					str = s.(string)
				}
				nullIf[i] = sdk.NullString{S: str}
			}
			csvOptions.NullIf = nullIf
		}
	}

	if v, ok := csvConfig["error_on_column_count_mismatch"]; ok && v.(string) != BooleanDefault && v.(string) != "" {
		errorOnColumnCountMismatch, err := booleanStringToBool(v.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing error_on_column_count_mismatch: %w", err)
		}
		csvOptions.ErrorOnColumnCountMismatch = &errorOnColumnCountMismatch
	}

	if v, ok := csvConfig["replace_invalid_characters"]; ok && v.(string) != BooleanDefault && v.(string) != "" {
		replaceInvalidCharacters, err := booleanStringToBool(v.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing replace_invalid_characters: %w", err)
		}
		csvOptions.ReplaceInvalidCharacters = &replaceInvalidCharacters
	}

	if v, ok := csvConfig["empty_field_as_null"]; ok && v.(string) != BooleanDefault && v.(string) != "" {
		emptyFieldAsNull, err := booleanStringToBool(v.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing empty_field_as_null: %w", err)
		}
		csvOptions.EmptyFieldAsNull = &emptyFieldAsNull
	}

	if v, ok := csvConfig["skip_byte_order_mark"]; ok && v.(string) != BooleanDefault && v.(string) != "" {
		skipByteOrderMark, err := booleanStringToBool(v.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing skip_byte_order_mark: %w", err)
		}
		csvOptions.SkipByteOrderMark = &skipByteOrderMark
	}

	if v, ok := csvConfig["encoding"]; ok && v.(string) != "" {
		encoding, err := sdk.ToCsvEncoding(v.(string))
		if err != nil {
			return nil, fmt.Errorf("parsing encoding: %w", err)
		}
		csvOptions.Encoding = &encoding
	}

	return csvOptions, nil
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

func stageCsvFileFormatToSchema(csv *sdk.FileFormatCsv) map[string]any {
	result := map[string]any{}

	if csv.Compression != "" {
		result["compression"] = csv.Compression
	}
	if csv.RecordDelimiter != "" {
		result["record_delimiter"] = csv.RecordDelimiter
	}
	if csv.FieldDelimiter != "" {
		result["field_delimiter"] = csv.FieldDelimiter
	}
	result["multi_line"] = booleanStringFromBool(csv.MultiLine)
	if csv.FileExtension != "" {
		result["file_extension"] = csv.FileExtension
	}
	result["parse_header"] = booleanStringFromBool(csv.ParseHeader)
	result["skip_header"] = csv.SkipHeader
	result["skip_blank_lines"] = booleanStringFromBool(csv.SkipBlankLines)
	if csv.DateFormat != "" {
		result["date_format"] = csv.DateFormat
	}
	if csv.TimeFormat != "" {
		result["time_format"] = csv.TimeFormat
	}
	if csv.TimestampFormat != "" {
		result["timestamp_format"] = csv.TimestampFormat
	}
	if csv.BinaryFormat != "" {
		result["binary_format"] = csv.BinaryFormat
	}
	if csv.Escape != "" {
		result["escape"] = csv.Escape
	}
	if csv.EscapeUnenclosedField != "" {
		result["escape_unenclosed_field"] = csv.EscapeUnenclosedField
	}
	result["trim_space"] = booleanStringFromBool(csv.TrimSpace)
	if csv.FieldOptionallyEnclosedBy != "" {
		result["field_optionally_enclosed_by"] = csv.FieldOptionallyEnclosedBy
	}
	if len(csv.NullIf) > 0 {
		result["null_if"] = csv.NullIf
	}
	result["error_on_column_count_mismatch"] = booleanStringFromBool(csv.ErrorOnColumnCountMismatch)
	result["replace_invalid_characters"] = booleanStringFromBool(csv.ReplaceInvalidCharacters)
	result["empty_field_as_null"] = booleanStringFromBool(csv.EmptyFieldAsNull)
	result["skip_byte_order_mark"] = booleanStringFromBool(csv.SkipByteOrderMark)
	if csv.Encoding != "" {
		result["encoding"] = csv.Encoding
	}

	return result
}
