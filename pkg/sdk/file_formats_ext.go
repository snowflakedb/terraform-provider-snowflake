package sdk

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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

func (d *FileFormatCsvDetails) ID() SchemaObjectIdentifier {
	return d.Id
}

func (d *FileFormatJsonDetails) ID() SchemaObjectIdentifier {
	return d.Id
}

func (d *FileFormatAvroDetails) ID() SchemaObjectIdentifier {
	return d.Id
}

func (d *FileFormatOrcDetails) ID() SchemaObjectIdentifier {
	return d.Id
}

func (d *FileFormatParquetDetails) ID() SchemaObjectIdentifier {
	return d.Id
}

func (d *FileFormatXmlDetails) ID() SchemaObjectIdentifier {
	return d.Id
}

// DescribeCsvDetails fetches and parses describe output for a CSV file format.
func (v *fileFormats) DescribeCsvDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatCsvDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatCsvDetails(properties, id)
}

// DescribeJsonDetails fetches and parses describe output for a JSON file format.
func (v *fileFormats) DescribeJsonDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatJsonDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatJsonDetails(properties, id)
}

// DescribeAvroDetails fetches and parses describe output for an Avro file format.
func (v *fileFormats) DescribeAvroDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatAvroDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatAvroDetails(properties, id)
}

// DescribeOrcDetails fetches and parses describe output for an ORC file format.
func (v *fileFormats) DescribeOrcDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatOrcDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatOrcDetails(properties, id)
}

// DescribeParquetDetails fetches and parses describe output for a Parquet file format.
func (v *fileFormats) DescribeParquetDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatParquetDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatParquetDetails(properties, id)
}

// DescribeXmlDetails fetches and parses describe output for an XML file format.
func (v *fileFormats) DescribeXmlDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatXmlDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatXmlDetails(properties, id)
}

func parseFileFormatCsvDetails(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatCsvDetails, error) {
	details := &FileFormatCsvDetails{Id: id}
	var errs []error
	for _, p := range properties {
		if p.Value == "" {
			continue
		}
		v := p.Value
		switch p.Name {
		case "RECORD_DELIMITER":
			details.RecordDelimiter = &StageFileFormatStringOrNone{Value: &v}
		case "FIELD_DELIMITER":
			details.FieldDelimiter = &StageFileFormatStringOrNone{Value: &v}
		case "FILE_EXTENSION":
			details.FileExtension = &v
		case "SKIP_HEADER":
			i, err := strconv.Atoi(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast SKIP_HEADER value "%s" to int: %w`, v, err))
			} else {
				details.SkipHeader = &i
			}
		case "PARSE_HEADER":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast PARSE_HEADER value "%s" to bool: %w`, v, err))
			} else {
				details.ParseHeader = &b
			}
		case "DATE_FORMAT":
			details.DateFormat = &StageFileFormatStringOrAuto{Value: &v}
		case "TIME_FORMAT":
			details.TimeFormat = &StageFileFormatStringOrAuto{Value: &v}
		case "TIMESTAMP_FORMAT":
			details.TimestampFormat = &StageFileFormatStringOrAuto{Value: &v}
		case "BINARY_FORMAT":
			bf := BinaryFormat(v)
			details.BinaryFormat = &bf
		case "ESCAPE":
			details.Escape = &StageFileFormatStringOrNone{Value: &v}
		case "ESCAPE_UNENCLOSED_FIELD":
			details.EscapeUnenclosedField = &StageFileFormatStringOrNone{Value: &v}
		case "TRIM_SPACE":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err))
			} else {
				details.TrimSpace = &b
			}
		case "FIELD_OPTIONALLY_ENCLOSED_BY":
			details.FieldOptionallyEnclosedBy = &StageFileFormatStringOrNone{Value: &v}
		case "NULL_IF":
			details.NullIf = parseNullIfProperty(v)
		case "COMPRESSION":
			comp := CsvCompression(v)
			details.Compression = &comp
		case "ERROR_ON_COLUMN_COUNT_MISMATCH":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast ERROR_ON_COLUMN_COUNT_MISMATCH value "%s" to bool: %w`, v, err))
			} else {
				details.ErrorOnColumnCountMismatch = &b
			}
		case "SKIP_BLANK_LINES":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast SKIP_BLANK_LINES value "%s" to bool: %w`, v, err))
			} else {
				details.SkipBlankLines = &b
			}
		case "REPLACE_INVALID_CHARACTERS":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
			} else {
				details.ReplaceInvalidCharacters = &b
			}
		case "EMPTY_FIELD_AS_NULL":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast EMPTY_FIELD_AS_NULL value "%s" to bool: %w`, v, err))
			} else {
				details.EmptyFieldAsNull = &b
			}
		case "SKIP_BYTE_ORDER_MARK":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err))
			} else {
				details.SkipByteOrderMark = &b
			}
		case "ENCODING":
			enc := CsvEncoding(v)
			details.Encoding = &enc
		}
	}
	return details, errors.Join(errs...)
}

func parseFileFormatJsonDetails(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatJsonDetails, error) {
	details := &FileFormatJsonDetails{Id: id}
	var errs []error
	for _, p := range properties {
		if p.Value == "" {
			continue
		}
		v := p.Value
		switch p.Name {
		case "FILE_EXTENSION":
			details.FileExtension = &v
		case "DATE_FORMAT":
			details.DateFormat = &StageFileFormatStringOrAuto{Value: &v}
		case "TIME_FORMAT":
			details.TimeFormat = &StageFileFormatStringOrAuto{Value: &v}
		case "TIMESTAMP_FORMAT":
			details.TimestampFormat = &StageFileFormatStringOrAuto{Value: &v}
		case "BINARY_FORMAT":
			bf := BinaryFormat(v)
			details.BinaryFormat = &bf
		case "TRIM_SPACE":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err))
			} else {
				details.TrimSpace = &b
			}
		case "MULTI_LINE":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast MULTI_LINE value "%s" to bool: %w`, v, err))
			} else {
				details.MultiLine = &b
			}
		case "NULL_IF":
			details.NullIf = parseNullIfProperty(v)
		case "COMPRESSION":
			comp := JsonCompression(v)
			details.Compression = &comp
		case "ENABLE_OCTAL":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast ENABLE_OCTAL value "%s" to bool: %w`, v, err))
			} else {
				details.EnableOctal = &b
			}
		case "ALLOW_DUPLICATE":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast ALLOW_DUPLICATE value "%s" to bool: %w`, v, err))
			} else {
				details.AllowDuplicate = &b
			}
		case "STRIP_OUTER_ARRAY":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast STRIP_OUTER_ARRAY value "%s" to bool: %w`, v, err))
			} else {
				details.StripOuterArray = &b
			}
		case "STRIP_NULL_VALUES":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast STRIP_NULL_VALUES value "%s" to bool: %w`, v, err))
			} else {
				details.StripNullValues = &b
			}
		case "IGNORE_UTF8_ERRORS":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, v, err))
			} else {
				details.IgnoreUtf8Errors = &b
			}
		case "REPLACE_INVALID_CHARACTERS":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
			} else {
				details.ReplaceInvalidCharacters = &b
			}
		case "SKIP_BYTE_ORDER_MARK":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err))
			} else {
				details.SkipByteOrderMark = &b
			}
		}
	}
	return details, errors.Join(errs...)
}

func parseFileFormatAvroDetails(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatAvroDetails, error) {
	details := &FileFormatAvroDetails{Id: id}
	var errs []error
	for _, p := range properties {
		if p.Value == "" {
			continue
		}
		v := p.Value
		switch p.Name {
		case "TRIM_SPACE":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err))
			} else {
				details.TrimSpace = &b
			}
		case "NULL_IF":
			details.NullIf = parseNullIfProperty(v)
		case "COMPRESSION":
			comp := AvroCompression(v)
			details.Compression = &comp
		case "REPLACE_INVALID_CHARACTERS":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
			} else {
				details.ReplaceInvalidCharacters = &b
			}
		}
	}
	return details, errors.Join(errs...)
}

func parseFileFormatOrcDetails(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatOrcDetails, error) {
	details := &FileFormatOrcDetails{Id: id}
	var errs []error
	for _, p := range properties {
		if p.Value == "" {
			continue
		}
		v := p.Value
		switch p.Name {
		case "TRIM_SPACE":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err))
			} else {
				details.TrimSpace = &b
			}
		case "NULL_IF":
			details.NullIf = parseNullIfProperty(v)
		case "REPLACE_INVALID_CHARACTERS":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
			} else {
				details.ReplaceInvalidCharacters = &b
			}
		}
	}
	return details, errors.Join(errs...)
}

func parseFileFormatParquetDetails(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatParquetDetails, error) {
	details := &FileFormatParquetDetails{Id: id}
	var errs []error
	for _, p := range properties {
		if p.Value == "" {
			continue
		}
		v := p.Value
		switch p.Name {
		case "TRIM_SPACE":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err))
			} else {
				details.TrimSpace = &b
			}
		case "NULL_IF":
			details.NullIf = parseNullIfProperty(v)
		case "COMPRESSION":
			comp := ParquetCompression(v)
			details.Compression = &comp
		case "BINARY_AS_TEXT":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast BINARY_AS_TEXT value "%s" to bool: %w`, v, err))
			} else {
				details.BinaryAsText = &b
			}
		case "USE_LOGICAL_TYPE":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast USE_LOGICAL_TYPE value "%s" to bool: %w`, v, err))
			} else {
				details.UseLogicalType = &b
			}
		case "USE_VECTORIZED_SCANNER":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast USE_VECTORIZED_SCANNER value "%s" to bool: %w`, v, err))
			} else {
				details.UseVectorizedScanner = &b
			}
		case "REPLACE_INVALID_CHARACTERS":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
			} else {
				details.ReplaceInvalidCharacters = &b
			}
		}
	}
	return details, errors.Join(errs...)
}

func parseFileFormatXmlDetails(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatXmlDetails, error) {
	details := &FileFormatXmlDetails{Id: id}
	var errs []error
	for _, p := range properties {
		if p.Value == "" {
			continue
		}
		v := p.Value
		switch p.Name {
		case "COMPRESSION":
			comp := XmlCompression(v)
			details.Compression = &comp
		case "IGNORE_UTF8_ERRORS":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, v, err))
			} else {
				details.IgnoreUtf8Errors = &b
			}
		case "PRESERVE_SPACE":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast PRESERVE_SPACE value "%s" to bool: %w`, v, err))
			} else {
				details.PreserveSpace = &b
			}
		case "STRIP_OUTER_ELEMENT":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast STRIP_OUTER_ELEMENT value "%s" to bool: %w`, v, err))
			} else {
				details.StripOuterElement = &b
			}
		case "DISABLE_AUTO_CONVERT":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast DISABLE_AUTO_CONVERT value "%s" to bool: %w`, v, err))
			} else {
				details.DisableAutoConvert = &b
			}
		case "REPLACE_INVALID_CHARACTERS":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
			} else {
				details.ReplaceInvalidCharacters = &b
			}
		case "SKIP_BYTE_ORDER_MARK":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err))
			} else {
				details.SkipByteOrderMark = &b
			}
		}
	}
	return details, errors.Join(errs...)
}

func (d *FileFormatAllDetails) ID() SchemaObjectIdentifier {
	return d.Id
}

// DescribeAllDetails fetches and parses describe output for any file format type.
func (v *fileFormats) DescribeAllDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatAllDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatAllDetails(properties, id)
}

func parseFileFormatAllDetails(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatAllDetails, error) {
	details := &FileFormatAllDetails{Id: id}
	var errs []error

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
				details.CsvRecordDelimiter = &StageFileFormatStringOrNone{Value: &v}
			case "FIELD_DELIMITER":
				details.CsvFieldDelimiter = &StageFileFormatStringOrNone{Value: &v}
			case "FILE_EXTENSION":
				details.CsvFileExtension = &v
			case "SKIP_HEADER":
				i, err := strconv.Atoi(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast SKIP_HEADER value "%s" to int: %w`, v, err))
				} else {
					details.CsvSkipHeader = &i
				}
			case "PARSE_HEADER":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast PARSE_HEADER value "%s" to bool: %w`, v, err))
				} else {
					details.CsvParseHeader = &b
				}
			case "DATE_FORMAT":
				details.CsvDateFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "TIME_FORMAT":
				details.CsvTimeFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "TIMESTAMP_FORMAT":
				details.CsvTimestampFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "BINARY_FORMAT":
				bf := BinaryFormat(v)
				details.CsvBinaryFormat = &bf
			case "ESCAPE":
				details.CsvEscape = &StageFileFormatStringOrNone{Value: &v}
			case "ESCAPE_UNENCLOSED_FIELD":
				details.CsvEscapeUnenclosedField = &StageFileFormatStringOrNone{Value: &v}
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err))
				} else {
					details.CsvTrimSpace = &b
				}
			case "FIELD_OPTIONALLY_ENCLOSED_BY":
				details.CsvFieldOptionallyEnclosedBy = &StageFileFormatStringOrNone{Value: &v}
			case "NULL_IF":
				details.CsvNullIf = parseNullIfProperty(v)
			case "COMPRESSION":
				comp := CsvCompression(v)
				details.CsvCompression = &comp
			case "ERROR_ON_COLUMN_COUNT_MISMATCH":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast ERROR_ON_COLUMN_COUNT_MISMATCH value "%s" to bool: %w`, v, err))
				} else {
					details.CsvErrorOnColumnCountMismatch = &b
				}
			case "SKIP_BLANK_LINES":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast SKIP_BLANK_LINES value "%s" to bool: %w`, v, err))
				} else {
					details.CsvSkipBlankLines = &b
				}
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
				} else {
					details.CsvReplaceInvalidCharacters = &b
				}
			case "EMPTY_FIELD_AS_NULL":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast EMPTY_FIELD_AS_NULL value "%s" to bool: %w`, v, err))
				} else {
					details.CsvEmptyFieldAsNull = &b
				}
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err))
				} else {
					details.CsvSkipByteOrderMark = &b
				}
			case "ENCODING":
				enc := CsvEncoding(v)
				details.CsvEncoding = &enc
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
				details.JsonFileExtension = &v
			case "DATE_FORMAT":
				details.JsonDateFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "TIME_FORMAT":
				details.JsonTimeFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "TIMESTAMP_FORMAT":
				details.JsonTimestampFormat = &StageFileFormatStringOrAuto{Value: &v}
			case "BINARY_FORMAT":
				bf := BinaryFormat(v)
				details.JsonBinaryFormat = &bf
			case "TRIM_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err))
				} else {
					details.JsonTrimSpace = &b
				}
			case "MULTI_LINE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast MULTI_LINE value "%s" to bool: %w`, v, err))
				} else {
					details.JsonMultiLine = &b
				}
			case "NULL_IF":
				details.JsonNullIf = parseNullIfProperty(v)
			case "COMPRESSION":
				comp := JsonCompression(v)
				details.JsonCompression = &comp
			case "ENABLE_OCTAL":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast ENABLE_OCTAL value "%s" to bool: %w`, v, err))
				} else {
					details.JsonEnableOctal = &b
				}
			case "ALLOW_DUPLICATE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast ALLOW_DUPLICATE value "%s" to bool: %w`, v, err))
				} else {
					details.JsonAllowDuplicate = &b
				}
			case "STRIP_OUTER_ARRAY":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast STRIP_OUTER_ARRAY value "%s" to bool: %w`, v, err))
				} else {
					details.JsonStripOuterArray = &b
				}
			case "STRIP_NULL_VALUES":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast STRIP_NULL_VALUES value "%s" to bool: %w`, v, err))
				} else {
					details.JsonStripNullValues = &b
				}
			case "IGNORE_UTF8_ERRORS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, v, err))
				} else {
					details.JsonIgnoreUtf8Errors = &b
				}
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
				} else {
					details.JsonReplaceInvalidCharacters = &b
				}
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err))
				} else {
					details.JsonSkipByteOrderMark = &b
				}
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
					errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err))
				} else {
					details.AvroTrimSpace = &b
				}
			case "NULL_IF":
				details.AvroNullIf = parseNullIfProperty(v)
			case "COMPRESSION":
				comp := AvroCompression(v)
				details.AvroCompression = &comp
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
				} else {
					details.AvroReplaceInvalidCharacters = &b
				}
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
					errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err))
				} else {
					details.OrcTrimSpace = &b
				}
			case "NULL_IF":
				details.OrcNullIf = parseNullIfProperty(v)
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
				} else {
					details.OrcReplaceInvalidCharacters = &b
				}
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
					errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, v, err))
				} else {
					details.ParquetTrimSpace = &b
				}
			case "NULL_IF":
				details.ParquetNullIf = parseNullIfProperty(v)
			case "COMPRESSION":
				comp := ParquetCompression(v)
				details.ParquetCompression = &comp
			case "BINARY_AS_TEXT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast BINARY_AS_TEXT value "%s" to bool: %w`, v, err))
				} else {
					details.ParquetBinaryAsText = &b
				}
			case "USE_LOGICAL_TYPE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast USE_LOGICAL_TYPE value "%s" to bool: %w`, v, err))
				} else {
					details.ParquetUseLogicalType = &b
				}
			case "USE_VECTORIZED_SCANNER":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast USE_VECTORIZED_SCANNER value "%s" to bool: %w`, v, err))
				} else {
					details.ParquetUseVectorizedScanner = &b
				}
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
				} else {
					details.ParquetReplaceInvalidCharacters = &b
				}
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
				details.XmlCompression = &comp
			case "IGNORE_UTF8_ERRORS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, v, err))
				} else {
					details.XmlIgnoreUtf8Errors = &b
				}
			case "PRESERVE_SPACE":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast PRESERVE_SPACE value "%s" to bool: %w`, v, err))
				} else {
					details.XmlPreserveSpace = &b
				}
			case "STRIP_OUTER_ELEMENT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast STRIP_OUTER_ELEMENT value "%s" to bool: %w`, v, err))
				} else {
					details.XmlStripOuterElement = &b
				}
			case "DISABLE_AUTO_CONVERT":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast DISABLE_AUTO_CONVERT value "%s" to bool: %w`, v, err))
				} else {
					details.XmlDisableAutoConvert = &b
				}
			case "REPLACE_INVALID_CHARACTERS":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, v, err))
				} else {
					details.XmlReplaceInvalidCharacters = &b
				}
			case "SKIP_BYTE_ORDER_MARK":
				b, err := strconv.ParseBool(v)
				if err != nil {
					errs = append(errs, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, v, err))
				} else {
					details.XmlSkipByteOrderMark = &b
				}
			}
		}
	default:
		return nil, fmt.Errorf("describe did not return a recognized file format type")
	}

	return details, errors.Join(errs...)
}

func parseNullIfProperty(v string) []NullString {
	nullIf := []NullString{}
	for s := range strings.SplitSeq(strings.Trim(v, "[]"), ", ") {
		if s == "" {
			continue
		}
		nullIf = append(nullIf, NullString{s})
	}
	return nullIf
}
