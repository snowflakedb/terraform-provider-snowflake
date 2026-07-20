package sdk

import (
	"context"
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
	return parseFileFormatCsv(properties, id), nil
}

// DescribeJsonDetails fetches and parses describe output for a JSON file format.
func (v *fileFormats) DescribeJsonDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatJson, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatJson(properties, id), nil
}

// DescribeAvroDetails fetches and parses describe output for an Avro file format.
func (v *fileFormats) DescribeAvroDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatAvro, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatAvro(properties, id), nil
}

// DescribeOrcDetails fetches and parses describe output for an ORC file format.
func (v *fileFormats) DescribeOrcDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatOrc, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatOrc(properties, id), nil
}

// DescribeParquetDetails fetches and parses describe output for a Parquet file format.
func (v *fileFormats) DescribeParquetDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatParquet, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatParquet(properties, id), nil
}

// DescribeXmlDetails fetches and parses describe output for an XML file format.
func (v *fileFormats) DescribeXmlDetails(ctx context.Context, id SchemaObjectIdentifier) (*FileFormatXml, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseFileFormatXml(properties, id), nil
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
func parseFileFormatCsv(properties []FileFormatProperty, id SchemaObjectIdentifier) *FileFormatCsv {
	csv := &FileFormatCsv{Id: id}
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
			val, _ := strconv.Atoi(prop.Value)
			csv.SkipHeader = val
		case "PARSE_HEADER":
			csv.ParseHeader = prop.Value == "true"
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
			csv.TrimSpace = prop.Value == "true"
		case "FIELD_OPTIONALLY_ENCLOSED_BY":
			csv.FieldOptionallyEnclosedBy = prop.Value
		case "NULL_IF":
			csv.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		case "COMPRESSION":
			csv.Compression = prop.Value
		case "ERROR_ON_COLUMN_COUNT_MISMATCH":
			csv.ErrorOnColumnCountMismatch = prop.Value == "true"
		case "VALIDATE_UTF8":
			csv.ValidateUtf8 = prop.Value == "true"
		case "SKIP_BLANK_LINES":
			csv.SkipBlankLines = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			csv.ReplaceInvalidCharacters = prop.Value == "true"
		case "EMPTY_FIELD_AS_NULL":
			csv.EmptyFieldAsNull = prop.Value == "true"
		case "SKIP_BYTE_ORDER_MARK":
			csv.SkipByteOrderMark = prop.Value == "true"
		case "ENCODING":
			csv.Encoding = prop.Value
		case "MULTI_LINE":
			csv.MultiLine = prop.Value == "true"
		}
	}
	return csv
}

func parseFileFormatJson(properties []FileFormatProperty, id SchemaObjectIdentifier) *FileFormatJson {
	json := &FileFormatJson{Id: id}
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
			json.TrimSpace = prop.Value == "true"
		case "MULTI_LINE":
			json.MultiLine = prop.Value == "true"
		case "NULL_IF":
			json.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		case "FILE_EXTENSION":
			json.FileExtension = prop.Value
		case "ENABLE_OCTAL":
			json.EnableOctal = prop.Value == "true"
		case "ALLOW_DUPLICATE":
			json.AllowDuplicate = prop.Value == "true"
		case "STRIP_OUTER_ARRAY":
			json.StripOuterArray = prop.Value == "true"
		case "STRIP_NULL_VALUES":
			json.StripNullValues = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			json.ReplaceInvalidCharacters = prop.Value == "true"
		case "IGNORE_UTF8_ERRORS":
			json.IgnoreUtf8Errors = prop.Value == "true"
		case "SKIP_BYTE_ORDER_MARK":
			json.SkipByteOrderMark = prop.Value == "true"
		}
	}
	return json
}

func parseFileFormatAvro(properties []FileFormatProperty, id SchemaObjectIdentifier) *FileFormatAvro {
	avro := &FileFormatAvro{Id: id}
	for _, prop := range properties {
		switch prop.Name {
		case "TYPE":
			avro.Type = prop.Value
		case "COMPRESSION":
			avro.Compression = prop.Value
		case "TRIM_SPACE":
			avro.TrimSpace = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			avro.ReplaceInvalidCharacters = prop.Value == "true"
		case "NULL_IF":
			avro.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		}
	}
	return avro
}

func parseFileFormatOrc(properties []FileFormatProperty, id SchemaObjectIdentifier) *FileFormatOrc {
	orc := &FileFormatOrc{Id: id}
	for _, prop := range properties {
		switch prop.Name {
		case "TYPE":
			orc.Type = prop.Value
		case "TRIM_SPACE":
			orc.TrimSpace = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			orc.ReplaceInvalidCharacters = prop.Value == "true"
		case "NULL_IF":
			orc.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		}
	}
	return orc
}

func parseFileFormatParquet(properties []FileFormatProperty, id SchemaObjectIdentifier) *FileFormatParquet {
	parquet := &FileFormatParquet{Id: id}
	for _, prop := range properties {
		switch prop.Name {
		case "TYPE":
			parquet.Type = prop.Value
		case "COMPRESSION":
			parquet.Compression = prop.Value
		case "BINARY_AS_TEXT":
			parquet.BinaryAsText = prop.Value == "true"
		case "USE_LOGICAL_TYPE":
			parquet.UseLogicalType = prop.Value == "true"
		case "TRIM_SPACE":
			parquet.TrimSpace = prop.Value == "true"
		case "USE_VECTORIZED_SCANNER":
			parquet.UseVectorizedScanner = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			parquet.ReplaceInvalidCharacters = prop.Value == "true"
		case "NULL_IF":
			parquet.NullIf = ParseCommaSeparatedStringArray(prop.Value, false)
		}
	}
	return parquet
}

func parseFileFormatXml(properties []FileFormatProperty, id SchemaObjectIdentifier) *FileFormatXml {
	xml := &FileFormatXml{Id: id}
	for _, prop := range properties {
		switch prop.Name {
		case "TYPE":
			xml.Type = prop.Value
		case "COMPRESSION":
			xml.Compression = prop.Value
		case "IGNORE_UTF8_ERRORS":
			xml.IgnoreUtf8Errors = prop.Value == "true"
		case "PRESERVE_SPACE":
			xml.PreserveSpace = prop.Value == "true"
		case "STRIP_OUTER_ELEMENT":
			xml.StripOuterElement = prop.Value == "true"
		case "DISABLE_SNOWFLAKE_DATA":
			xml.DisableSnowflakeData = prop.Value == "true"
		case "DISABLE_AUTO_CONVERT":
			xml.DisableAutoConvert = prop.Value == "true"
		case "REPLACE_INVALID_CHARACTERS":
			xml.ReplaceInvalidCharacters = prop.Value == "true"
		case "SKIP_BYTE_ORDER_MARK":
			xml.SkipByteOrderMark = prop.Value == "true"
		}
	}
	return xml
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

	switch details.Type {
	case FileFormatTypeCsv:
		details.Csv = parseFileFormatCsv(properties, id)
	case FileFormatTypeJson:
		details.Json = parseFileFormatJson(properties, id)
	case FileFormatTypeAvro:
		details.Avro = parseFileFormatAvro(properties, id)
	case FileFormatTypeOrc:
		details.Orc = parseFileFormatOrc(properties, id)
	case FileFormatTypeParquet:
		details.Parquet = parseFileFormatParquet(properties, id)
	case FileFormatTypeXml:
		details.Xml = parseFileFormatXml(properties, id)
	}

	return details, nil
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
