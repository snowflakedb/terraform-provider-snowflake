package sdk

import (
	"context"
	"fmt"
	"log"
	"strconv"
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

func (d *FileFormatCsv) ID() SchemaObjectIdentifier {
	return d.Id
}

func (d *FileFormatJson) ID() SchemaObjectIdentifier {
	return d.Id
}

func (d *FileFormatAvro) ID() SchemaObjectIdentifier {
	return d.Id
}

func (d *FileFormatOrc) ID() SchemaObjectIdentifier {
	return d.Id
}

func (d *FileFormatParquet) ID() SchemaObjectIdentifier {
	return d.Id
}

func (d *FileFormatXml) ID() SchemaObjectIdentifier {
	return d.Id
}

func (d *FileFormatAllDetails) ID() SchemaObjectIdentifier {
	return d.Id
}

// DescribeCsvDetails fetches and parses describe output for a CSV file format.
func (v *fileFormats) DescribeCsvDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatCsv, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatCsv(properties, id)
}

// DescribeJsonDetails fetches and parses describe output for a JSON file format.
func (v *fileFormats) DescribeJsonDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatJson, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatJson(properties, id)
}

// DescribeAvroDetails fetches and parses describe output for an Avro file format.
func (v *fileFormats) DescribeAvroDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatAvro, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatAvro(properties, id)
}

// DescribeOrcDetails fetches and parses describe output for an ORC file format.
func (v *fileFormats) DescribeOrcDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatOrc, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatOrc(properties, id)
}

// DescribeParquetDetails fetches and parses describe output for a Parquet file format.
func (v *fileFormats) DescribeParquetDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatParquet, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatParquet(properties, id)
}

// DescribeXmlDetails fetches and parses describe output for an XML file format.
func (v *fileFormats) DescribeXmlDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatXml, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatXml(properties, id)
}

// DescribeAllDetails fetches and parses describe output for any file format type.
func (v *fileFormats) DescribeAllDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatAllDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatAllDetails(properties, id)
}

// parseFileFormatCsv parses DESCRIBE FILE FORMAT output for a CSV file format. It is also reused
// by stages_ext.go to parse the file format properties embedded in DESCRIBE STAGE output.
func parseFileFormatCsv(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatCsv, error) {
	csv := &FileFormatCsv{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "TYPE":
			csv.Type = prop.Value
		case "RECORD_DELIMITER":
			csv.RecordDelimiter = prop.Value
		case "FIELD_DELIMITER":
			csv.FieldDelimiter = prop.Value
		case "FILE_EXTENSION":
			csv.FileExtension = prop.Value
		case "SKIP_HEADER":
			val, err := strconv.Atoi(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast SKIP_HEADER value "%s" to int: %w`, prop.Value, err))
			} else {
				csv.SkipHeader = val
			}
		case "PARSE_HEADER":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast PARSE_HEADER value "%s" to bool: %w`, prop.Value, err))
			} else {
				csv.ParseHeader = val
			}
		case "DATE_FORMAT":
			csv.DateFormat = prop.Value
		case "TIME_FORMAT":
			csv.TimeFormat = prop.Value
		case "TIMESTAMP_FORMAT":
			csv.TimestampFormat = prop.Value
		case "BINARY_FORMAT":
			csv.BinaryFormat = prop.Value
		case "ESCAPE":
			csv.Escape = prop.Value
		case "ESCAPE_UNENCLOSED_FIELD":
			csv.EscapeUnenclosedField = prop.Value
		case "TRIM_SPACE":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, prop.Value, err))
			} else {
				csv.TrimSpace = val
			}
		case "FIELD_OPTIONALLY_ENCLOSED_BY":
			csv.FieldOptionallyEnclosedBy = prop.Value
		case "NULL_IF":
			csv.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMPRESSION":
			csv.Compression = prop.Value
		case "ERROR_ON_COLUMN_COUNT_MISMATCH":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast ERROR_ON_COLUMN_COUNT_MISMATCH value "%s" to bool: %w`, prop.Value, err))
			} else {
				csv.ErrorOnColumnCountMismatch = val
			}
		case "VALIDATE_UTF8":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast VALIDATE_UTF8 value "%s" to bool: %w`, prop.Value, err))
			} else {
				csv.ValidateUtf8 = val
			}
		case "SKIP_BLANK_LINES":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast SKIP_BLANK_LINES value "%s" to bool: %w`, prop.Value, err))
			} else {
				csv.SkipBlankLines = val
			}
		case "REPLACE_INVALID_CHARACTERS":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, prop.Value, err))
			} else {
				csv.ReplaceInvalidCharacters = val
			}
		case "EMPTY_FIELD_AS_NULL":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast EMPTY_FIELD_AS_NULL value "%s" to bool: %w`, prop.Value, err))
			} else {
				csv.EmptyFieldAsNull = val
			}
		case "SKIP_BYTE_ORDER_MARK":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, prop.Value, err))
			} else {
				csv.SkipByteOrderMark = val
			}
		case "ENCODING":
			csv.Encoding = prop.Value
		case "MULTI_LINE":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast MULTI_LINE value "%s" to bool: %w`, prop.Value, err))
			} else {
				csv.MultiLine = val
			}
		default:
			log.Printf("[DEBUG] unknown CSV file format property: %s", prop.Name)
		}
	}
	return csv, JoinErrors(errs...)
}

func parseFileFormatJson(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatJson, error) {
	json := &FileFormatJson{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "TYPE":
			json.Type = prop.Value
		case "COMPRESSION":
			json.Compression = prop.Value
		case "DATE_FORMAT":
			json.DateFormat = prop.Value
		case "TIME_FORMAT":
			json.TimeFormat = prop.Value
		case "TIMESTAMP_FORMAT":
			json.TimestampFormat = prop.Value
		case "BINARY_FORMAT":
			json.BinaryFormat = prop.Value
		case "TRIM_SPACE":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, prop.Value, err))
			} else {
				json.TrimSpace = val
			}
		case "MULTI_LINE":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast MULTI_LINE value "%s" to bool: %w`, prop.Value, err))
			} else {
				json.MultiLine = val
			}
		case "NULL_IF":
			json.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		case "FILE_EXTENSION":
			json.FileExtension = prop.Value
		case "ENABLE_OCTAL":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast ENABLE_OCTAL value "%s" to bool: %w`, prop.Value, err))
			} else {
				json.EnableOctal = val
			}
		case "ALLOW_DUPLICATE":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast ALLOW_DUPLICATE value "%s" to bool: %w`, prop.Value, err))
			} else {
				json.AllowDuplicate = val
			}
		case "STRIP_OUTER_ARRAY":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast STRIP_OUTER_ARRAY value "%s" to bool: %w`, prop.Value, err))
			} else {
				json.StripOuterArray = val
			}
		case "STRIP_NULL_VALUES":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast STRIP_NULL_VALUES value "%s" to bool: %w`, prop.Value, err))
			} else {
				json.StripNullValues = val
			}
		case "REPLACE_INVALID_CHARACTERS":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, prop.Value, err))
			} else {
				json.ReplaceInvalidCharacters = val
			}
		case "IGNORE_UTF8_ERRORS":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, prop.Value, err))
			} else {
				json.IgnoreUtf8Errors = val
			}
		case "SKIP_BYTE_ORDER_MARK":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, prop.Value, err))
			} else {
				json.SkipByteOrderMark = val
			}
		default:
			log.Printf("[DEBUG] unknown JSON file format property: %s", prop.Name)
		}
	}
	return json, JoinErrors(errs...)
}

func parseFileFormatAvro(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatAvro, error) {
	avro := &FileFormatAvro{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "TYPE":
			avro.Type = prop.Value
		case "COMPRESSION":
			avro.Compression = prop.Value
		case "TRIM_SPACE":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, prop.Value, err))
			} else {
				avro.TrimSpace = val
			}
		case "REPLACE_INVALID_CHARACTERS":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, prop.Value, err))
			} else {
				avro.ReplaceInvalidCharacters = val
			}
		case "NULL_IF":
			avro.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		default:
			log.Printf("[DEBUG] unknown Avro file format property: %s", prop.Name)
		}
	}
	return avro, JoinErrors(errs...)
}

func parseFileFormatOrc(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatOrc, error) {
	orc := &FileFormatOrc{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "TYPE":
			orc.Type = prop.Value
		case "TRIM_SPACE":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, prop.Value, err))
			} else {
				orc.TrimSpace = val
			}
		case "REPLACE_INVALID_CHARACTERS":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, prop.Value, err))
			} else {
				orc.ReplaceInvalidCharacters = val
			}
		case "NULL_IF":
			orc.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		default:
			log.Printf("[DEBUG] unknown ORC file format property: %s", prop.Name)
		}
	}
	return orc, JoinErrors(errs...)
}

func parseFileFormatParquet(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatParquet, error) {
	parquet := &FileFormatParquet{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "TYPE":
			parquet.Type = prop.Value
		case "COMPRESSION":
			parquet.Compression = prop.Value
		case "BINARY_AS_TEXT":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast BINARY_AS_TEXT value "%s" to bool: %w`, prop.Value, err))
			} else {
				parquet.BinaryAsText = val
			}
		case "USE_LOGICAL_TYPE":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast USE_LOGICAL_TYPE value "%s" to bool: %w`, prop.Value, err))
			} else {
				parquet.UseLogicalType = val
			}
		case "TRIM_SPACE":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast TRIM_SPACE value "%s" to bool: %w`, prop.Value, err))
			} else {
				parquet.TrimSpace = val
			}
		case "USE_VECTORIZED_SCANNER":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast USE_VECTORIZED_SCANNER value "%s" to bool: %w`, prop.Value, err))
			} else {
				parquet.UseVectorizedScanner = val
			}
		case "REPLACE_INVALID_CHARACTERS":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, prop.Value, err))
			} else {
				parquet.ReplaceInvalidCharacters = val
			}
		case "NULL_IF":
			parquet.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		default:
			log.Printf("[DEBUG] unknown Parquet file format property: %s", prop.Name)
		}
	}
	return parquet, JoinErrors(errs...)
}

func parseFileFormatXml(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatXml, error) {
	xml := &FileFormatXml{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "TYPE":
			xml.Type = prop.Value
		case "COMPRESSION":
			xml.Compression = prop.Value
		case "IGNORE_UTF8_ERRORS":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast IGNORE_UTF8_ERRORS value "%s" to bool: %w`, prop.Value, err))
			} else {
				xml.IgnoreUtf8Errors = val
			}
		case "PRESERVE_SPACE":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast PRESERVE_SPACE value "%s" to bool: %w`, prop.Value, err))
			} else {
				xml.PreserveSpace = val
			}
		case "STRIP_OUTER_ELEMENT":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast STRIP_OUTER_ELEMENT value "%s" to bool: %w`, prop.Value, err))
			} else {
				xml.StripOuterElement = val
			}
		case "DISABLE_SNOWFLAKE_DATA":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast DISABLE_SNOWFLAKE_DATA value "%s" to bool: %w`, prop.Value, err))
			} else {
				xml.DisableSnowflakeData = val
			}
		case "DISABLE_AUTO_CONVERT":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast DISABLE_AUTO_CONVERT value "%s" to bool: %w`, prop.Value, err))
			} else {
				xml.DisableAutoConvert = val
			}
		case "REPLACE_INVALID_CHARACTERS":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast REPLACE_INVALID_CHARACTERS value "%s" to bool: %w`, prop.Value, err))
			} else {
				xml.ReplaceInvalidCharacters = val
			}
		case "SKIP_BYTE_ORDER_MARK":
			val, err := strconv.ParseBool(prop.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf(`cannot cast SKIP_BYTE_ORDER_MARK value "%s" to bool: %w`, prop.Value, err))
			} else {
				xml.SkipByteOrderMark = val
			}
		default:
			log.Printf("[DEBUG] unknown XML file format property: %s", prop.Name)
		}
	}
	return xml, JoinErrors(errs...)
}

// parseFileFormatAllDetails reuses the per-type parsers instead of re-implementing the property
// parsing for each type.
func parseFileFormatAllDetails(properties []FileFormatProperty, id SchemaObjectIdentifier) (*FileFormatAllDetails, error) {
	details := &FileFormatAllDetails{Id: id}
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

	var err error
	switch details.Type {
	case FileFormatTypeCsv:
		details.Csv, err = parseFileFormatCsv(properties, id)
	case FileFormatTypeJson:
		details.Json, err = parseFileFormatJson(properties, id)
	case FileFormatTypeAvro:
		details.Avro, err = parseFileFormatAvro(properties, id)
	case FileFormatTypeOrc:
		details.Orc, err = parseFileFormatOrc(properties, id)
	case FileFormatTypeParquet:
		details.Parquet, err = parseFileFormatParquet(properties, id)
	case FileFormatTypeXml:
		details.Xml, err = parseFileFormatXml(properties, id)
	}
	if err != nil {
		return nil, err
	}

	return details, nil
}
