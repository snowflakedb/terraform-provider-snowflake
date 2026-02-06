package resources

import (
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// parseStageFileFormat2 parses the stage file format from the resource data to an SDK object.
func parseStageFileFormat2(d *schema.ResourceData) (sdk.StageFileFormatRequest, error) {
	if len(d.Get("file_format").([]any)) == 0 {
		return sdk.StageFileFormatRequest{}, nil
	}
	fileFormatReq := sdk.NewStageFileFormatRequest()

	if v, ok := d.GetOk("file_format.0.format_name"); ok {
		id, err := sdk.ParseSchemaObjectIdentifier(v.(string))
		if err != nil {
			return sdk.StageFileFormatRequest{}, fmt.Errorf("parsing format_name: %w", err)
		}
		fileFormatReq.WithFormatName(id)
	}

	if _, ok := d.GetOk("file_format.0.csv"); ok {
		csvOptions, err := parseCsvFileFormatOptions2(d)
		if err != nil {
			return sdk.StageFileFormatRequest{}, err
		}
		fileFormatReq.WithFileFormatOptions(sdk.FileFormatOptions{
			CsvOptions: csvOptions,
		})
	}

	return *fileFormatReq, nil
}

// parseCsvFileFormatOptions2 parses the CSV file format options from the resource data to an SDK object.
func parseCsvFileFormatOptions2(d *schema.ResourceData) (*sdk.FileFormatCsvOptions, error) {
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
	)
	if err != nil {
		return nil, err
	}

	if v, ok := d.GetOk(prefix + "null_if"); ok {
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

	return csvOptions, nil
}
