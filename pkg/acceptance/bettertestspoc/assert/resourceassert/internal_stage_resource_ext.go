package resourceassert

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (i *InternalStageResourceAssert) HasDirectoryEnableString(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("directory.#", "1"))
	i.AddAssertion(assert.ValueSet("directory.0.enable", expected))
	return i
}

func (i *InternalStageResourceAssert) HasDirectoryAutoRefreshString(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("directory.#", "1"))
	i.AddAssertion(assert.ValueSet("directory.0.auto_refresh", expected))
	return i
}

func (i *InternalStageResourceAssert) HasDirectory(enable bool, autoRefresh bool) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("directory.#", "1"))
	i.AddAssertion(assert.ValueSet("directory.0.enable", strconv.FormatBool(enable)))
	i.AddAssertion(assert.ValueSet("directory.0.auto_refresh", strconv.FormatBool(autoRefresh)))
	return i
}

func (i *InternalStageResourceAssert) HasEncryptionSnowflakeFull() *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("encryption.#", "1"))
	i.AddAssertion(assert.ValueSet("encryption.0.snowflake_full.#", "1"))
	i.AddAssertion(assert.ValueSet("encryption.0.snowflake_sse.#", "0"))
	return i
}

func (i *InternalStageResourceAssert) HasEncryptionSnowflakeSse() *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("encryption.#", "1"))
	i.AddAssertion(assert.ValueSet("encryption.0.snowflake_full.#", "0"))
	i.AddAssertion(assert.ValueSet("encryption.0.snowflake_sse.#", "1"))
	return i
}

func (i *InternalStageResourceAssert) HasStageTypeEnum(expected sdk.StageType) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("stage_type", string(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatEmpty() *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.#", "0"))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatFormatName(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.#", "1"))
	i.AddAssertion(assert.ValueSet("file_format.0.format_name", expected))
	i.AddAssertion(assert.ValueSet("file_format.0.csv.#", "0"))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsv() *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.format_name", ""))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvCompression(expected sdk.CsvCompression) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.compression", string(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvRecordDelimiter(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.record_delimiter", expected))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvFieldDelimiter(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.field_delimiter", expected))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvMultiLine(expected bool) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.multi_line", strconv.FormatBool(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvFileExtension(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.file_extension", expected))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvParseHeader(expected bool) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.parse_header", strconv.FormatBool(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvSkipHeader(expected int) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.skip_header", strconv.Itoa(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvSkipBlankLines(expected bool) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.skip_blank_lines", strconv.FormatBool(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvDateFormat(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.date_format", expected))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvTimeFormat(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.time_format", expected))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvTimestampFormat(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.timestamp_format", expected))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvBinaryFormat(expected sdk.BinaryFormat) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.binary_format", string(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvEscape(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.escape", expected))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvEscapeUnenclosedField(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.escape_unenclosed_field", expected))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvTrimSpace(expected bool) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.trim_space", strconv.FormatBool(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvFieldOptionallyEnclosedBy(expected string) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.field_optionally_enclosed_by", expected))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvNullIfCount(expected int) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.null_if.#", strconv.Itoa(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvErrorOnColumnCountMismatch(expected bool) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.error_on_column_count_mismatch", strconv.FormatBool(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvReplaceInvalidCharacters(expected bool) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.replace_invalid_characters", strconv.FormatBool(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvEmptyFieldAsNull(expected bool) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.empty_field_as_null", strconv.FormatBool(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvSkipByteOrderMark(expected bool) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.skip_byte_order_mark", strconv.FormatBool(expected)))
	return i
}

func (i *InternalStageResourceAssert) HasFileFormatCsvEncoding(expected sdk.CsvEncoding) *InternalStageResourceAssert {
	i.AddAssertion(assert.ValueSet("file_format.0.csv.0.encoding", string(expected)))
	return i
}
