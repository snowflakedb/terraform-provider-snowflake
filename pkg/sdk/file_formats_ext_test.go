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
		{Name: "SKIP_HEADER", Value: "1"},
		{Name: "PARSE_HEADER", Value: "false"},
		{Name: "NULL_IF", Value: "[NULL, ]"},
		{Name: "TRIM_SPACE", Value: "true"},
	}

	csv, err := parseFileFormatCsv(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, csv.Id)
	require.Equal(t, "CSV", csv.Type)
	require.Equal(t, "GZIP", csv.Compression)
	require.Equal(t, "\\n", csv.RecordDelimiter)
	require.Equal(t, ",", csv.FieldDelimiter)
	require.Equal(t, 1, csv.SkipHeader)
	require.False(t, csv.ParseHeader)
	require.Equal(t, []string{"NULL", ""}, csv.NullIf)
	require.True(t, csv.TrimSpace)
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

func TestParseFileFormatJson(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "JSON"},
		{Name: "COMPRESSION", Value: "AUTO"},
		{Name: "STRIP_OUTER_ARRAY", Value: "true"},
		{Name: "NULL_IF", Value: "[]"},
	}

	json, err := parseFileFormatJson(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, json.Id)
	require.Equal(t, "JSON", json.Type)
	require.Equal(t, "AUTO", json.Compression)
	require.True(t, json.StripOuterArray)
	require.Equal(t, []string{}, json.NullIf)
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

func TestParseFileFormatParquet(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	properties := []FileFormatProperty{
		{Name: "TYPE", Value: "PARQUET"},
		{Name: "COMPRESSION", Value: "SNAPPY"},
		{Name: "BINARY_AS_TEXT", Value: "true"},
		{Name: "USE_LOGICAL_TYPE", Value: "true"},
		{Name: "USE_VECTORIZED_SCANNER", Value: "false"},
		{Name: "NULL_IF", Value: "[NULL]"},
	}

	parquet, err := parseFileFormatParquet(properties, id)
	require.NoError(t, err)

	require.Equal(t, id, parquet.Id)
	require.Equal(t, "PARQUET", parquet.Type)
	require.Equal(t, "SNAPPY", parquet.Compression)
	require.True(t, parquet.BinaryAsText)
	require.True(t, parquet.UseLogicalType)
	require.False(t, parquet.UseVectorizedScanner)
	require.Equal(t, []string{"NULL"}, parquet.NullIf)
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
	require.True(t, xml.SkipByteOrderMark)
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
}
