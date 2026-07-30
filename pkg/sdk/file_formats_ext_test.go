package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFileFormatCsv(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "CSV"},
		{Name: "COMPRESSION", Value: "GZIP"},
		{Name: "RECORD_DELIMITER", Value: "\\n"},
		{Name: "FIELD_DELIMITER", Value: ","},
		{Name: "FILE_EXTENSION", Value: "csv"},
		{Name: "SKIP_HEADER", Value: "1"},
		{Name: "PARSE_HEADER", Value: "false"},
		{Name: "SKIP_BLANK_LINES", Value: "true"},
		{Name: "DATE_FORMAT", Value: "YYYY-MM-DD"},
		{Name: "TIME_FORMAT", Value: "HH24:MI:SS"},
		{Name: "TIMESTAMP_FORMAT", Value: "YYYY-MM-DD HH24:MI:SS"},
		{Name: "BINARY_FORMAT", Value: "HEX"},
		{Name: "ESCAPE", Value: "\\"},
		{Name: "ESCAPE_UNENCLOSED_FIELD", Value: "\\"},
		{Name: "TRIM_SPACE", Value: "true"},
		{Name: "FIELD_OPTIONALLY_ENCLOSED_BY", Value: "'"},
		{Name: "NULL_IF", Value: "[NULL, ]"},
		{Name: "ERROR_ON_COLUMN_COUNT_MISMATCH", Value: "true"},
		{Name: "VALIDATE_UTF8", Value: "true"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "true"},
		{Name: "EMPTY_FIELD_AS_NULL", Value: "true"},
		{Name: "SKIP_BYTE_ORDER_MARK", Value: "true"},
		{Name: "ENCODING", Value: "UTF8"},
		{Name: "MULTI_LINE", Value: "true"},
	}

	csv, err := parseFileFormatCsv(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, csv.Id)
	require.Equal(t, "CSV", csv.Type)
	require.Equal(t, "GZIP", csv.Compression)
	require.Equal(t, "\\n", csv.RecordDelimiter)
	require.Equal(t, ",", csv.FieldDelimiter)
	require.Equal(t, "csv", csv.FileExtension)
	require.Equal(t, 1, csv.SkipHeader)
	require.False(t, csv.ParseHeader)
	require.True(t, csv.SkipBlankLines)
	require.Equal(t, "YYYY-MM-DD", csv.DateFormat)
	require.Equal(t, "HH24:MI:SS", csv.TimeFormat)
	require.Equal(t, "YYYY-MM-DD HH24:MI:SS", csv.TimestampFormat)
	require.Equal(t, "HEX", csv.BinaryFormat)
	require.Equal(t, "\\", csv.Escape)
	require.Equal(t, "\\", csv.EscapeUnenclosedField)
	require.True(t, csv.TrimSpace)
	require.Equal(t, "'", csv.FieldOptionallyEnclosedBy)
	require.Equal(t, []string{"NULL", ""}, csv.NullIf)
	require.True(t, csv.ErrorOnColumnCountMismatch)
	require.True(t, csv.ValidateUtf8)
	require.True(t, csv.ReplaceInvalidCharacters)
	require.True(t, csv.EmptyFieldAsNull)
	require.True(t, csv.SkipByteOrderMark)
	require.Equal(t, "UTF8", csv.Encoding)
	require.True(t, csv.MultiLine)
}

func TestParseFileFormatCsv_invalidSkipHeader(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "CSV"},
		{Name: "SKIP_HEADER", Value: "not-a-number"},
	}

	_, err := parseFileFormatCsv(properties, id)
	require.ErrorContains(t, err, `cannot cast SKIP_HEADER value "not-a-number" to int`)
}

func TestParseFileFormatCsv_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"PARSE_HEADER"},
		{"TRIM_SPACE"},
		{"ERROR_ON_COLUMN_COUNT_MISMATCH"},
		{"VALIDATE_UTF8"},
		{"SKIP_BLANK_LINES"},
		{"REPLACE_INVALID_CHARACTERS"},
		{"EMPTY_FIELD_AS_NULL"},
		{"SKIP_BYTE_ORDER_MARK"},
		{"MULTI_LINE"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "CSV"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatCsv(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatCsv_multipleInvalidValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "CSV"},
		{Name: "SKIP_HEADER", Value: "not-a-number"},
		{Name: "TRIM_SPACE", Value: "not-a-bool"},
	}

	_, err := parseFileFormatCsv(properties, id)
	require.ErrorContains(t, err, `cannot cast SKIP_HEADER value "not-a-number" to int`)
	require.ErrorContains(t, err, `cannot cast TRIM_SPACE value "not-a-bool" to bool`)
}

func TestParseFileFormatCsv_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "CSV"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	csv, err := parseFileFormatCsv(properties, id)
	require.NoError(t, err)
	require.Equal(t, "CSV", csv.Type)
}

func TestParseFileFormatJson(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "JSON"},
		{Name: "COMPRESSION", Value: "AUTO"},
		{Name: "DATE_FORMAT", Value: "YYYY-MM-DD"},
		{Name: "TIME_FORMAT", Value: "HH24:MI:SS"},
		{Name: "TIMESTAMP_FORMAT", Value: "YYYY-MM-DD HH24:MI:SS"},
		{Name: "BINARY_FORMAT", Value: "HEX"},
		{Name: "TRIM_SPACE", Value: "true"},
		{Name: "MULTI_LINE", Value: "true"},
		{Name: "STRIP_OUTER_ARRAY", Value: "true"},
		{Name: "NULL_IF", Value: "[]"},
		{Name: "FILE_EXTENSION", Value: "json"},
		{Name: "ENABLE_OCTAL", Value: "true"},
		{Name: "ALLOW_DUPLICATE", Value: "true"},
		{Name: "STRIP_NULL_VALUES", Value: "true"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "true"},
		{Name: "IGNORE_UTF8_ERRORS", Value: "true"},
		{Name: "SKIP_BYTE_ORDER_MARK", Value: "true"},
	}

	json, err := parseFileFormatJson(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, json.Id)
	require.Equal(t, "JSON", json.Type)
	require.Equal(t, "AUTO", json.Compression)
	require.Equal(t, "YYYY-MM-DD", json.DateFormat)
	require.Equal(t, "HH24:MI:SS", json.TimeFormat)
	require.Equal(t, "YYYY-MM-DD HH24:MI:SS", json.TimestampFormat)
	require.Equal(t, "HEX", json.BinaryFormat)
	require.True(t, json.TrimSpace)
	require.True(t, json.MultiLine)
	require.True(t, json.StripOuterArray)
	require.Equal(t, []string{}, json.NullIf)
	require.Equal(t, "json", json.FileExtension)
	require.True(t, json.EnableOctal)
	require.True(t, json.AllowDuplicate)
	require.True(t, json.StripNullValues)
	require.True(t, json.ReplaceInvalidCharacters)
	require.True(t, json.IgnoreUtf8Errors)
	require.True(t, json.SkipByteOrderMark)
}

func TestParseFileFormatJson_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"TRIM_SPACE"},
		{"MULTI_LINE"},
		{"ENABLE_OCTAL"},
		{"ALLOW_DUPLICATE"},
		{"STRIP_OUTER_ARRAY"},
		{"STRIP_NULL_VALUES"},
		{"REPLACE_INVALID_CHARACTERS"},
		{"IGNORE_UTF8_ERRORS"},
		{"SKIP_BYTE_ORDER_MARK"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "JSON"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatJson(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatJson_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "JSON"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	json, err := parseFileFormatJson(properties, id)
	require.NoError(t, err)
	require.Equal(t, "JSON", json.Type)
}

func TestParseFileFormatAvro(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "AVRO"},
		{Name: "COMPRESSION", Value: "GZIP"},
		{Name: "TRIM_SPACE", Value: "true"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "true"},
		{Name: "NULL_IF", Value: "[NULL, ]"},
	}

	avro, err := parseFileFormatAvro(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, avro.Id)
	require.Equal(t, "AVRO", avro.Type)
	require.Equal(t, "GZIP", avro.Compression)
	require.True(t, avro.TrimSpace)
	require.True(t, avro.ReplaceInvalidCharacters)
	require.Equal(t, []string{"NULL", ""}, avro.NullIf)
}

func TestParseFileFormatAvro_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"TRIM_SPACE"},
		{"REPLACE_INVALID_CHARACTERS"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "AVRO"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatAvro(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatAvro_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "AVRO"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	avro, err := parseFileFormatAvro(properties, id)
	require.NoError(t, err)
	require.Equal(t, "AVRO", avro.Type)
}

func TestParseFileFormatOrc(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "ORC"},
		{Name: "TRIM_SPACE", Value: "true"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "false"},
		{Name: "NULL_IF", Value: "[NULL]"},
	}

	orc, err := parseFileFormatOrc(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, orc.Id)
	require.Equal(t, "ORC", orc.Type)
	require.True(t, orc.TrimSpace)
	require.False(t, orc.ReplaceInvalidCharacters)
	require.Equal(t, []string{"NULL"}, orc.NullIf)
}

func TestParseFileFormatOrc_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"TRIM_SPACE"},
		{"REPLACE_INVALID_CHARACTERS"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "ORC"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatOrc(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatOrc_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "ORC"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	orc, err := parseFileFormatOrc(properties, id)
	require.NoError(t, err)
	require.Equal(t, "ORC", orc.Type)
}

func TestParseFileFormatParquet(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "PARQUET"},
		{Name: "COMPRESSION", Value: "SNAPPY"},
		{Name: "BINARY_AS_TEXT", Value: "true"},
		{Name: "USE_LOGICAL_TYPE", Value: "true"},
		{Name: "TRIM_SPACE", Value: "true"},
		{Name: "USE_VECTORIZED_SCANNER", Value: "false"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "true"},
		{Name: "NULL_IF", Value: "[NULL]"},
	}

	parquet, err := parseFileFormatParquet(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, parquet.Id)
	require.Equal(t, "PARQUET", parquet.Type)
	require.Equal(t, "SNAPPY", parquet.Compression)
	require.True(t, parquet.BinaryAsText)
	require.True(t, parquet.UseLogicalType)
	require.True(t, parquet.TrimSpace)
	require.False(t, parquet.UseVectorizedScanner)
	require.True(t, parquet.ReplaceInvalidCharacters)
	require.Equal(t, []string{"NULL"}, parquet.NullIf)
}

func TestParseFileFormatParquet_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"BINARY_AS_TEXT"},
		{"USE_LOGICAL_TYPE"},
		{"TRIM_SPACE"},
		{"USE_VECTORIZED_SCANNER"},
		{"REPLACE_INVALID_CHARACTERS"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "PARQUET"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatParquet(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatParquet_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "PARQUET"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	parquet, err := parseFileFormatParquet(properties, id)
	require.NoError(t, err)
	require.Equal(t, "PARQUET", parquet.Type)
}

func TestParseFileFormatXml(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "XML"},
		{Name: "COMPRESSION", Value: "GZIP"},
		{Name: "IGNORE_UTF8_ERRORS", Value: "true"},
		{Name: "PRESERVE_SPACE", Value: "true"},
		{Name: "STRIP_OUTER_ELEMENT", Value: "false"},
		{Name: "DISABLE_SNOWFLAKE_DATA", Value: "true"},
		{Name: "DISABLE_AUTO_CONVERT", Value: "true"},
		{Name: "REPLACE_INVALID_CHARACTERS", Value: "true"},
		{Name: "SKIP_BYTE_ORDER_MARK", Value: "true"},
	}

	xml, err := parseFileFormatXml(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, xml.Id)
	require.Equal(t, "XML", xml.Type)
	require.Equal(t, "GZIP", xml.Compression)
	require.True(t, xml.IgnoreUtf8Errors)
	require.True(t, xml.PreserveSpace)
	require.False(t, xml.StripOuterElement)
	require.True(t, xml.DisableSnowflakeData)
	require.True(t, xml.DisableAutoConvert)
	require.True(t, xml.ReplaceInvalidCharacters)
	require.True(t, xml.SkipByteOrderMark)
}

func TestParseFileFormatXml_invalidBoolValues(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	testCases := []struct {
		propertyName string
	}{
		{"IGNORE_UTF8_ERRORS"},
		{"PRESERVE_SPACE"},
		{"STRIP_OUTER_ELEMENT"},
		{"DISABLE_SNOWFLAKE_DATA"},
		{"DISABLE_AUTO_CONVERT"},
		{"REPLACE_INVALID_CHARACTERS"},
		{"SKIP_BYTE_ORDER_MARK"},
	}
	for _, tc := range testCases {
		t.Run(tc.propertyName, func(t *testing.T) {
			properties := []FileFormatProperty{
				{Name: "TYPE", Value: "XML"},
				{Name: tc.propertyName, Value: "not-a-bool"},
			}

			_, err := parseFileFormatXml(properties, id)
			require.ErrorContains(t, err, `cannot cast `+tc.propertyName+` value "not-a-bool" to bool`)
		})
	}
}

func TestParseFileFormatXml_unknownProperty(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "XML"},
		{Name: "SOME_UNKNOWN_PROPERTY", Value: "whatever"},
	}

	xml, err := parseFileFormatXml(properties, id)
	require.NoError(t, err)
	require.Equal(t, "XML", xml.Type)
}

func TestParseFileFormatAllDetails(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("csv", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "CSV"},
			{Name: "COMPRESSION", Value: "GZIP"},
			{Name: "TRIM_SPACE", Value: "true"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, id, details.Id)
		require.Equal(t, FileFormatTypeCsv, details.Type)
		require.NotNil(t, details.Csv)
		require.Equal(t, "GZIP", details.Csv.Compression)
		require.True(t, details.Csv.TrimSpace)
		require.Nil(t, details.Json)
		require.Nil(t, details.Avro)
		require.Nil(t, details.Orc)
		require.Nil(t, details.Parquet)
		require.Nil(t, details.Xml)
	})

	t.Run("json", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "JSON"},
			{Name: "COMPRESSION", Value: "AUTO"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeJson, details.Type)
		require.NotNil(t, details.Json)
		require.Equal(t, "AUTO", details.Json.Compression)
		require.Nil(t, details.Csv)
	})

	t.Run("avro", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "AVRO"},
			{Name: "TRIM_SPACE", Value: "true"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeAvro, details.Type)
		require.NotNil(t, details.Avro)
		require.True(t, details.Avro.TrimSpace)
	})

	t.Run("orc", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "ORC"},
			{Name: "TRIM_SPACE", Value: "true"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeOrc, details.Type)
		require.NotNil(t, details.Orc)
		require.True(t, details.Orc.TrimSpace)
	})

	t.Run("parquet", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "PARQUET"},
			{Name: "COMPRESSION", Value: "SNAPPY"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeParquet, details.Type)
		require.NotNil(t, details.Parquet)
		require.Equal(t, "SNAPPY", details.Parquet.Compression)
	})

	t.Run("xml", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "XML"},
			{Name: "COMPRESSION", Value: "GZIP"},
		}
		details, err := parseFileFormatAllDetails(properties, id)
		require.NoError(t, err)
		require.Equal(t, FileFormatTypeXml, details.Type)
		require.NotNil(t, details.Xml)
		require.Equal(t, "GZIP", details.Xml.Compression)
	})

	t.Run("invalid type", func(t *testing.T) {
		_, err := parseFileFormatAllDetails([]FileFormatProperty{{Name: "TYPE", Value: "NOT_A_TYPE"}}, id)
		require.Error(t, err)
	})

	t.Run("missing type", func(t *testing.T) {
		_, err := parseFileFormatAllDetails([]FileFormatProperty{{Name: "COMPRESSION", Value: "GZIP"}}, id)
		require.ErrorContains(t, err, "describe did not return a recognized file format type")
	})

	t.Run("propagates csv parsing error", func(t *testing.T) {
		properties := []FileFormatProperty{
			{Name: "TYPE", Value: "CSV"},
			{Name: "SKIP_HEADER", Value: "not-a-number"},
		}
		_, err := parseFileFormatAllDetails(properties, id)
		require.ErrorContains(t, err, `cannot cast SKIP_HEADER value "not-a-number" to int`)
	})
}
