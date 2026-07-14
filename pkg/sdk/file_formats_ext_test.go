package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileFormatObjectOptionsFromShowResult_Csv(t *testing.T) {
	raw := `{
		"TYPE": "CSV",
		"COMPRESSION": "GZIP",
		"RECORD_DELIMITER": "\n",
		"FIELD_DELIMITER": ",",
		"SKIP_HEADER": 1,
		"PARSE_HEADER": false,
		"NULL_IF": ["NULL", "null"],
		"TRIM_SPACE": true
	}`

	options, err := fileFormatObjectOptionsFromShowResult(FileFormatTypeCsv, raw)
	require.NoError(t, err)
	require.Equal(t, Pointer(CsvCompressionGzip), options.CsvCompression)
	require.Equal(t, "\n", *options.CsvRecordDelimiter.Value)
	require.Equal(t, ",", *options.CsvFieldDelimiter.Value)
	require.Equal(t, 1, *options.CsvSkipHeader)
	require.False(t, *options.CsvParseHeader)
	require.Equal(t, []NullString{{"NULL"}, {"null"}}, options.CsvNullIf)
	require.True(t, *options.CsvTrimSpace)
}

func TestFileFormatObjectOptionsFromShowResult_Json(t *testing.T) {
	raw := `{
		"TYPE": "JSON",
		"COMPRESSION": "AUTO",
		"STRIP_OUTER_ARRAY": true,
		"NULL_IF": []
	}`

	options, err := fileFormatObjectOptionsFromShowResult(FileFormatTypeJson, raw)
	require.NoError(t, err)
	require.Equal(t, Pointer(JsonCompressionAuto), options.JsonCompression)
	require.True(t, *options.JsonStripOuterArray)
	require.Equal(t, []NullString{}, options.JsonNullIf)
}

func TestFileFormatObjectOptionsFromShowResult_InvalidJSON(t *testing.T) {
	_, err := fileFormatObjectOptionsFromShowResult(FileFormatTypeCsv, "not json")
	require.Error(t, err)
}

func TestParseNullIfProperty(t *testing.T) {
	require.Equal(t, []NullString{{"a"}, {"b"}}, parseNullIfProperty("[a, b]"))
	require.Equal(t, []NullString{}, parseNullIfProperty("[]"))
}
