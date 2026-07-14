package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (opts *DummyOperationFileFormatOptions) additionalValidations() error {
	if valueSet(opts.FileFormat) {
		return opts.FileFormat.validate()
	}
	return nil
}

func (opts *FileFormatObjectOptions) fieldsByType() map[FileFormatType][]any {
	return map[FileFormatType][]any{
		FileFormatTypeCsv: {
			opts.CsvCompression,
			opts.CsvRecordDelimiter,
			opts.CsvFieldDelimiter,
			opts.CsvMultiLine,
			opts.CsvFileExtension,
			opts.CsvParseHeader,
			opts.CsvSkipHeader,
			opts.CsvSkipBlankLines,
			opts.CsvDateFormat,
			opts.CsvTimeFormat,
			opts.CsvTimestampFormat,
			opts.CsvBinaryFormat,
			opts.CsvEscape,
			opts.CsvEscapeUnenclosedField,
			opts.CsvTrimSpace,
			opts.CsvFieldOptionallyEnclosedBy,
			opts.CsvNullIf,
			opts.CsvErrorOnColumnCountMismatch,
			opts.CsvReplaceInvalidCharacters,
			opts.CsvEmptyFieldAsNull,
			opts.CsvSkipByteOrderMark,
			opts.CsvEncoding,
		},
		FileFormatTypeJson: {
			opts.JsonCompression,
			opts.JsonDateFormat,
			opts.JsonTimeFormat,
			opts.JsonTimestampFormat,
			opts.JsonBinaryFormat,
			opts.JsonTrimSpace,
			opts.JsonMultiLine,
			opts.JsonNullIf,
			opts.JsonFileExtension,
			opts.JsonEnableOctal,
			opts.JsonAllowDuplicate,
			opts.JsonStripOuterArray,
			opts.JsonStripNullValues,
			opts.JsonReplaceInvalidCharacters,
			opts.JsonIgnoreUtf8Errors,
			opts.JsonSkipByteOrderMark,
		},
		FileFormatTypeAvro: {
			opts.AvroCompression,
			opts.AvroTrimSpace,
			opts.AvroReplaceInvalidCharacters,
			opts.AvroNullIf,
		},
		FileFormatTypeOrc: {
			opts.OrcTrimSpace,
			opts.OrcReplaceInvalidCharacters,
			opts.OrcNullIf,
		},
		FileFormatTypeParquet: {
			opts.ParquetCompression,
			opts.ParquetSnappyCompression,
			opts.ParquetBinaryAsText,
			opts.ParquetUseLogicalType,
			opts.ParquetTrimSpace,
			opts.ParquetUseVectorizedScanner,
			opts.ParquetReplaceInvalidCharacters,
			opts.ParquetNullIf,
		},
		FileFormatTypeXml: {
			opts.XmlCompression,
			opts.XmlIgnoreUtf8Errors,
			opts.XmlPreserveSpace,
			opts.XmlStripOuterElement,
			opts.XmlDisableAutoConvert,
			opts.XmlReplaceInvalidCharacters,
			opts.XmlSkipByteOrderMark,
		},
	}
}

// additionalValidations ensures fields belonging to more than one file format type are never set at once,
// e.g. setting both a Csv* and a Json* field in the same FileFormatObjectOptions (used by both Create and Alter Set).
func (opts *FileFormatObjectOptions) additionalValidations() error {
	typesWithFieldsSet := 0
	for _, fields := range opts.fieldsByType() {
		if anyValueSet(fields...) {
			typesWithFieldsSet++
		}
	}
	if typesWithFieldsSet > 1 {
		return fmt.Errorf("cannot set options for more than one file format type at once")
	}
	return nil
}

// additionalValidations ensures only the fields matching FileFormatType are set on Create.
func (opts *CreateFileFormatOptions) additionalValidations() error {
	fields := opts.FileFormatObjectOptions.fieldsByType()
	for formatType, typeFields := range fields {
		if formatType == opts.FileFormatType {
			continue
		}
		if anyValueSet(typeFields...) {
			return fmt.Errorf("cannot set %s fields when TYPE = %s", formatType, opts.FileFormatType)
		}
	}
	return nil
}

func (opts FileFormatOptions) validate() error {
	var errs []error
	if !exactlyOneValueSet(opts.CsvOptions, opts.JsonOptions, opts.AvroOptions, opts.OrcOptions, opts.ParquetOptions, opts.XmlOptions) {
		errs = append(errs, errExactlyOneOf("FileFormat", "CsvOptions", "JsonOptions", "AvroOptions", "OrcOptions", "ParquetOptions", "XmlOptions"))
	}
	if valueSet(opts.CsvOptions) {
		if everyValueSet(opts.CsvOptions.SkipHeader, opts.CsvOptions.ParseHeader) {
			errs = append(errs, errOneOf("FileFormat.CsvOptions", "SkipHeader", "ParseHeader"))
		}
		if valueSet(opts.CsvOptions.RecordDelimiter) {
			if !exactlyOneValueSet(opts.CsvOptions.RecordDelimiter.Value, opts.CsvOptions.RecordDelimiter.None) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.RecordDelimiter", "Value", "None"))
			}
		}
		if valueSet(opts.CsvOptions.FieldDelimiter) {
			if !exactlyOneValueSet(opts.CsvOptions.FieldDelimiter.Value, opts.CsvOptions.FieldDelimiter.None) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.FieldDelimiter", "Value", "None"))
			}
		}
		if valueSet(opts.CsvOptions.DateFormat) {
			if !exactlyOneValueSet(opts.CsvOptions.DateFormat.Value, opts.CsvOptions.DateFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.DateFormat", "Value", "Auto"))
			}
		}
		if valueSet(opts.CsvOptions.TimeFormat) {
			if !exactlyOneValueSet(opts.CsvOptions.TimeFormat.Value, opts.CsvOptions.TimeFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.TimeFormat", "Value", "Auto"))
			}
		}
		if valueSet(opts.CsvOptions.TimestampFormat) {
			if !exactlyOneValueSet(opts.CsvOptions.TimestampFormat.Value, opts.CsvOptions.TimestampFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.TimestampFormat", "Value", "Auto"))
			}
		}
		if valueSet(opts.CsvOptions.Escape) {
			if !exactlyOneValueSet(opts.CsvOptions.Escape.Value, opts.CsvOptions.Escape.None) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.Escape", "Value", "None"))
			}
		}
		if valueSet(opts.CsvOptions.EscapeUnenclosedField) {
			if !exactlyOneValueSet(opts.CsvOptions.EscapeUnenclosedField.Value, opts.CsvOptions.EscapeUnenclosedField.None) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.EscapeUnenclosedField", "Value", "None"))
			}
		}
		if valueSet(opts.CsvOptions.FieldOptionallyEnclosedBy) {
			if !exactlyOneValueSet(opts.CsvOptions.FieldOptionallyEnclosedBy.Value, opts.CsvOptions.FieldOptionallyEnclosedBy.None) {
				errs = append(errs, errExactlyOneOf("FileFormat.CsvOptions.FieldOptionallyEnclosedBy", "Value", "None"))
			}
		}
	}
	if valueSet(opts.JsonOptions) {
		if everyValueSet(opts.JsonOptions.IgnoreUtf8Errors, opts.JsonOptions.ReplaceInvalidCharacters) {
			errs = append(errs, errOneOf("FileFormat.JsonOptions", "IgnoreUtf8Errors", "ReplaceInvalidCharacters"))
		}
		if valueSet(opts.JsonOptions.DateFormat) {
			if !exactlyOneValueSet(opts.JsonOptions.DateFormat.Value, opts.JsonOptions.DateFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.JsonOptions.DateFormat", "Value", "Auto"))
			}
		}
		if valueSet(opts.JsonOptions.TimeFormat) {
			if !exactlyOneValueSet(opts.JsonOptions.TimeFormat.Value, opts.JsonOptions.TimeFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.JsonOptions.TimeFormat", "Value", "Auto"))
			}
		}
		if valueSet(opts.JsonOptions.TimestampFormat) {
			if !exactlyOneValueSet(opts.JsonOptions.TimestampFormat.Value, opts.JsonOptions.TimestampFormat.Auto) {
				errs = append(errs, errExactlyOneOf("FileFormat.JsonOptions.TimestampFormat", "Value", "Auto"))
			}
		}
	}
	if valueSet(opts.ParquetOptions) {
		if everyValueSet(opts.ParquetOptions.Compression, opts.ParquetOptions.SnappyCompression) {
			errs = append(errs, errOneOf("FileFormat.ParquetOptions", "Compression", "SnappyCompression"))
		}
	}
	if valueSet(opts.XmlOptions) {
		if everyValueSet(opts.XmlOptions.IgnoreUtf8Errors, opts.XmlOptions.ReplaceInvalidCharacters) {
			errs = append(errs, errOneOf("FileFormat.XmlOptions", "IgnoreUtf8Errors", "ReplaceInvalidCharacters"))
		}
	}
	return JoinErrors(errs...)
}

// DescribeDetails returns the DESCRIBE FILE FORMAT output parsed into typed, per-file-format-type options.
func (v *fileFormats) DescribeDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	details := &FileFormatDetails{}
	for _, p := range properties {
		if p.Name == "TYPE" {
			formatType, err := ToFileFormatType(p.Value)
			if err != nil {
				return nil, err
			}
			details.Type = formatType
			break
		}
	}

	switch details.Type {
	case FileFormatTypeCsv:
		for _, p := range properties {
			if p.Value == "" {
				continue
			}
			v := p.Value
			switch p.Name {
			case "RECORD_DELIMITER":
				details.Options.CsvRecordDelimiter = &StageFileFormatStringOrNone{Value: &v}
			case "FIELD_DELIMITER":
				details.Options.CsvFieldDelimiter = &StageFileFormatStringOrNone{Value: &v}
			case "FILE_EXTENSION":
				details.Options.CsvFileExtension = &v
			case "SKIP_HEADER":
				i, err := strconv.Atoi(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_HEADER value "%s" to int: %w`, v, err)
				}
				details.Options.CsvSkipHeader = &i
			case "PARSE_HEADER":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast PARSE_HEADER value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvParseHeader = &b
			case "DATE_FORMAT":
				details.Options.CsvDateFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "TIME_FORMAT":
				details.Options.CsvTimeFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "TIMESTAMP_FORMAT":
				details.Options.CsvTimestampFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "BINARY_FORMAT":
				bf := BinaryFormat(v)
				details.Options.CsvBinaryFormat = &bf
			case "ESCAPE":
				details.Options.CsvEscape = &StageFileFormatStringOrNone{Value: &v}
			case "ESCAPE_UNENCLOSED_FIELD":
				details.Options.CsvEscapeUnenclosedField = &StageFileFormatStringOrNone{Value: &v}
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvTrimSpace = &b
			case "FIELD_OPTIONALLY_ENCLOSED_BY":
				details.Options.CsvFieldOptionallyEnclosedBy = &StageFileFormatStringOrNone{Value: &v}
			case "NULL_IF":
				details.Options.CsvNullIf = parseNullIfProperty(v)
			case "COMPRESSION":
				comp := CsvCompression(v)
				details.Options.CsvCompression = &comp
			case "ERROR_ON_COLUMN_COUNT_MISMATCH":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast ERROR_ON_COLUMN_COUNT_MISMATCH value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvErrorOnColumnCountMismatch = &b
			case "SKIP_BLANK_LINES":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BLANK_LINES value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvSkipBlankLines = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvReplaceInvalidCharacters = &b
			case "EMPTY_FIELD_AS_NULL":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast EMPTY_FIELD_AS_NULL value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvEmptyFieldAsNull = &b
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err)
				}
				details.Options.CsvSkipByteOrderMark = &b
			case "ENCODING":
				enc := CsvEncoding(v)
				details.Options.CsvEncoding = &enc
			}
		}
	case FileFormatTypeJson:
		for _, p := range properties {
			if p.Value == "" {
				continue
			}
			v := p.Value
			switch p.Name {
			case "FILE_EXTENSION":
				details.Options.JsonFileExtension = &v
			case "DATE_FORMAT":
				details.Options.JsonDateFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "TIME_FORMAT":
				details.Options.JsonTimeFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "TIMESTAMP_FORMAT":
				details.Options.JsonTimestampFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "BINARY_FORMAT":
				bf := BinaryFormat(v)
				details.Options.JsonBinaryFormat = &bf
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonTrimSpace = &b
			case "MULTI_LINE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast MULTI_LINE value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonMultiLine = &b
			case "NULL_IF":
				details.Options.JsonNullIf = parseNullIfProperty(v)
			case "COMPRESSION":
				comp := JsonCompression(v)
				details.Options.JsonCompression = &comp
			case "ENABLE_OCTAL":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast ENABLE_OCTAL value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonEnableOctal = &b
			case "ALLOW_DUPLICATE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast ALLOW_DUPLICATE value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonAllowDuplicate = &b
			case "STRIP_OUTER_ARRAY":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast STRIP_OUTER_ARRAY value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonStripOuterArray = &b
			case "STRIP_NULL_VALUES":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast STRIP_NULL_VALUES value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonStripNullValues = &b
			case "IGNORE_UTF8_ERRORS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonIgnoreUtf8Errors = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonReplaceInvalidCharacters = &b
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err)
				}
				details.Options.JsonSkipByteOrderMark = &b
			}
		}
	case FileFormatTypeAvro:
		for _, p := range properties {
			if p.Value == "" {
				continue
			}
			v := p.Value
			switch p.Name {
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.AvroTrimSpace = &b
			case "NULL_IF":
				details.Options.AvroNullIf = parseNullIfProperty(v)
			case "COMPRESSION":
				comp := AvroCompression(v)
				details.Options.AvroCompression = &comp
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.AvroReplaceInvalidCharacters = &b
			}
		}
	case FileFormatTypeOrc:
		for _, p := range properties {
			if p.Value == "" {
				continue
			}
			v := p.Value
			switch p.Name {
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.OrcTrimSpace = &b
			case "NULL_IF":
				details.Options.OrcNullIf = parseNullIfProperty(v)
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.OrcReplaceInvalidCharacters = &b
			}
		}
	case FileFormatTypeParquet:
		for _, p := range properties {
			if p.Value == "" {
				continue
			}
			v := p.Value
			switch p.Name {
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.ParquetTrimSpace = &b
			case "NULL_IF":
				details.Options.ParquetNullIf = parseNullIfProperty(v)
			case "COMPRESSION":
				comp := ParquetCompression(v)
				details.Options.ParquetCompression = &comp
			case "BINARY_AS_TEXT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast BINARY_AS_TEXT value "%s" to bool: %w`, v, err)
				}
				details.Options.ParquetBinaryAsText = &b
			case "USE_LOGICAL_TYPE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast USE_LOGICAL_TYPE value "%s" to bool: %w`, v, err)
				}
				details.Options.ParquetUseLogicalType = &b
			case "USE_VECTORIZED_SCANNER":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast USE_VECTORIZED_SCANNER value "%s" to bool: %w`, v, err)
				}
				details.Options.ParquetUseVectorizedScanner = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.ParquetReplaceInvalidCharacters = &b
			}
		}
	case FileFormatTypeXml:
		for _, p := range properties {
			if p.Value == "" {
				continue
			}
			v := p.Value
			switch p.Name {
			case "COMPRESSION":
				comp := XmlCompression(v)
				details.Options.XmlCompression = &comp
			case "IGNORE_UTF8_ERRORS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlIgnoreUtf8Errors = &b
			case "PRESERVE_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast PRESERVE_SPACE value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlPreserveSpace = &b
			case "STRIP_OUTER_ELEMENT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast STRIP_OUTER_ELEMENT value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlStripOuterElement = &b
			case "DISABLE_AUTO_CONVERT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast DISABLE_AUTO_CONVERT value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlDisableAutoConvert = &b
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlReplaceInvalidCharacters = &b
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err)
				}
				details.Options.XmlSkipByteOrderMark = &b
			}
		}
	default:
		return nil, fmt.Errorf("describe did not return a recognized file format type")
	}

	return details, nil
}

func parseNullIfProperty(v string) []NullString {
	nullIf := []NullString{}
	for _, s := range strings.Split(strings.Trim(v, "[]"), ", ") {
		if s == "" {
			continue
		}
		nullIf = append(nullIf, NullString{s})
	}
	return nullIf
}

// showFileFormatsOptionsResult mirrors the JSON shape of the SHOW FILE FORMATS format_options column.
type showFileFormatsOptionsResult struct {
	Type                       string   `json:"TYPE"`
	RecordDelimiter            string   `json:"RECORD_DELIMITER"`
	FieldDelimiter             string   `json:"FIELD_DELIMITER"`
	FileExtension              string   `json:"FILE_EXTENSION"`
	SkipHeader                 int      `json:"SKIP_HEADER"`
	ParseHeader                bool     `json:"PARSE_HEADER"`
	DateFormat                 string   `json:"DATE_FORMAT"`
	TimeFormat                 string   `json:"TIME_FORMAT"`
	TimestampFormat            string   `json:"TIMESTAMP_FORMAT"`
	BinaryFormat               string   `json:"BINARY_FORMAT"`
	Escape                     string   `json:"ESCAPE"`
	EscapeUnenclosedField      string   `json:"ESCAPE_UNENCLOSED_FIELD"`
	TrimSpace                  bool     `json:"TRIM_SPACE"`
	FieldOptionallyEnclosedBy  string   `json:"FIELD_OPTIONALLY_ENCLOSED_BY"`
	NullIf                     []string `json:"NULL_IF"`
	Compression                string   `json:"COMPRESSION"`
	ErrorOnColumnCountMismatch bool     `json:"ERROR_ON_COLUMN_COUNT_MISMATCH"`
	SkipBlankLines             bool     `json:"SKIP_BLANK_LINES"`
	ReplaceInvalidCharacters   bool     `json:"REPLACE_INVALID_CHARACTERS"`
	EmptyFieldAsNull           bool     `json:"EMPTY_FIELD_AS_NULL"`
	SkipByteOrderMark          bool     `json:"SKIP_BYTE_ORDER_MARK"`
	Encoding                   string   `json:"ENCODING"`
	MultiLine                  bool     `json:"MULTI_LINE"`
	EnableOctal                bool     `json:"ENABLE_OCTAL"`
	AllowDuplicate             bool     `json:"ALLOW_DUPLICATE"`
	StripOuterArray            bool     `json:"STRIP_OUTER_ARRAY"`
	StripNullValues            bool     `json:"STRIP_NULL_VALUES"`
	IgnoreUTF8Errors           bool     `json:"IGNORE_UTF8_ERRORS"`
	BinaryAsText               bool     `json:"BINARY_AS_TEXT"`
	UseLogicalType             bool     `json:"USE_LOGICAL_TYPE"`
	UseVectorizedScanner       bool     `json:"USE_VECTORIZED_SCANNER"`
	PreserveSpace              bool     `json:"PRESERVE_SPACE"`
	StripOuterElement          bool     `json:"STRIP_OUTER_ELEMENT"`
	DisableAutoConvert         bool     `json:"DISABLE_AUTO_CONVERT"`
	SnappyCompression          bool     `json:"SNAPPY_COMPRESSION"`
}

// fileFormatObjectOptionsFromShowResult parses the JSON blob returned in the SHOW FILE FORMATS
// format_options column into the typed, per-file-format-type FileFormatObjectOptions.
func fileFormatObjectOptionsFromShowResult(formatType FileFormatType, raw string) (*FileFormatObjectOptions, error) {
	var input showFileFormatsOptionsResult
	if err := json.Unmarshal([]byte(raw), &input); err != nil {
		return nil, fmt.Errorf("cannot parse format options: %w", err)
	}

	nullIf := make([]NullString, len(input.NullIf))
	for i, s := range input.NullIf {
		nullIf[i] = NullString{s}
	}

	options := &FileFormatObjectOptions{}
	switch formatType {
	case FileFormatTypeCsv:
		compression := CsvCompression(input.Compression)
		binaryFormat := BinaryFormat(input.BinaryFormat)
		encoding := CsvEncoding(input.Encoding)
		options.CsvCompression = &compression
		options.CsvRecordDelimiter = &StageFileFormatStringOrNone{Value: &input.RecordDelimiter}
		options.CsvFieldDelimiter = &StageFileFormatStringOrNone{Value: &input.FieldDelimiter}
		options.CsvFileExtension = &input.FileExtension
		options.CsvParseHeader = &input.ParseHeader
		options.CsvSkipHeader = &input.SkipHeader
		options.CsvSkipBlankLines = &input.SkipBlankLines
		options.CsvDateFormat = &StageFileFormatStringOrAuto{Value: &input.DateFormat}
		options.CsvTimeFormat = &StageFileFormatStringOrAuto{Value: &input.TimeFormat}
		options.CsvTimestampFormat = &StageFileFormatStringOrAuto{Value: &input.TimestampFormat}
		options.CsvBinaryFormat = &binaryFormat
		options.CsvEscape = &StageFileFormatStringOrNone{Value: &input.Escape}
		options.CsvEscapeUnenclosedField = &StageFileFormatStringOrNone{Value: &input.EscapeUnenclosedField}
		options.CsvTrimSpace = &input.TrimSpace
		options.CsvFieldOptionallyEnclosedBy = &StageFileFormatStringOrNone{Value: &input.FieldOptionallyEnclosedBy}
		options.CsvNullIf = nullIf
		options.CsvErrorOnColumnCountMismatch = &input.ErrorOnColumnCountMismatch
		options.CsvReplaceInvalidCharacters = &input.ReplaceInvalidCharacters
		options.CsvEmptyFieldAsNull = &input.EmptyFieldAsNull
		options.CsvSkipByteOrderMark = &input.SkipByteOrderMark
		options.CsvEncoding = &encoding
	case FileFormatTypeJson:
		compression := JsonCompression(input.Compression)
		binaryFormat := BinaryFormat(input.BinaryFormat)
		options.JsonCompression = &compression
		options.JsonDateFormat = &StageFileFormatStringOrAuto{Value: &input.DateFormat}
		options.JsonTimeFormat = &StageFileFormatStringOrAuto{Value: &input.TimeFormat}
		options.JsonTimestampFormat = &StageFileFormatStringOrAuto{Value: &input.TimestampFormat}
		options.JsonBinaryFormat = &binaryFormat
		options.JsonTrimSpace = &input.TrimSpace
		options.JsonMultiLine = &input.MultiLine
		options.JsonNullIf = nullIf
		options.JsonFileExtension = &input.FileExtension
		options.JsonEnableOctal = &input.EnableOctal
		options.JsonAllowDuplicate = &input.AllowDuplicate
		options.JsonStripOuterArray = &input.StripOuterArray
		options.JsonStripNullValues = &input.StripNullValues
		options.JsonReplaceInvalidCharacters = &input.ReplaceInvalidCharacters
		options.JsonIgnoreUtf8Errors = &input.IgnoreUTF8Errors
		options.JsonSkipByteOrderMark = &input.SkipByteOrderMark
	case FileFormatTypeAvro:
		compression := AvroCompression(input.Compression)
		options.AvroCompression = &compression
		options.AvroTrimSpace = &input.TrimSpace
		options.AvroReplaceInvalidCharacters = &input.ReplaceInvalidCharacters
		options.AvroNullIf = nullIf
	case FileFormatTypeOrc:
		options.OrcTrimSpace = &input.TrimSpace
		options.OrcReplaceInvalidCharacters = &input.ReplaceInvalidCharacters
		options.OrcNullIf = nullIf
	case FileFormatTypeParquet:
		compression := ParquetCompression(input.Compression)
		options.ParquetCompression = &compression
		options.ParquetSnappyCompression = &input.SnappyCompression
		options.ParquetBinaryAsText = &input.BinaryAsText
		options.ParquetUseLogicalType = &input.UseLogicalType
		options.ParquetTrimSpace = &input.TrimSpace
		options.ParquetUseVectorizedScanner = &input.UseVectorizedScanner
		options.ParquetReplaceInvalidCharacters = &input.ReplaceInvalidCharacters
		options.ParquetNullIf = nullIf
	case FileFormatTypeXml:
		compression := XmlCompression(input.Compression)
		options.XmlCompression = &compression
		options.XmlIgnoreUtf8Errors = &input.IgnoreUTF8Errors
		options.XmlPreserveSpace = &input.PreserveSpace
		options.XmlStripOuterElement = &input.StripOuterElement
		options.XmlDisableAutoConvert = &input.DisableAutoConvert
		options.XmlReplaceInvalidCharacters = &input.ReplaceInvalidCharacters
		options.XmlSkipByteOrderMark = &input.SkipByteOrderMark
	}
	return options, nil
}
