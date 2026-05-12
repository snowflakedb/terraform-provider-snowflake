package resourceassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func stageHasFileFormatFormatName(expected string) []assert.ResourceAssertion {
	return []assert.ResourceAssertion{
		assert.ValueSet("file_format.#", "1"),
		assert.ValueSet("file_format.0.format_name", expected),
		assert.ValueSet("file_format.0.csv.#", "0"),
		assert.ValueSet("file_format.0.json.#", "0"),
		assert.ValueSet("file_format.0.avro.#", "0"),
		assert.ValueSet("file_format.0.orc.#", "0"),
		assert.ValueSet("file_format.0.parquet.#", "0"),
		assert.ValueSet("file_format.0.xml.#", "0"),
	}
}

func stageHasFileFormatCsv() []assert.ResourceAssertion {
	return []assert.ResourceAssertion{
		assert.ValueSet("file_format.#", "1"),
		assert.ValueSet("file_format.0.csv.#", "1"),
		assert.ValueSet("file_format.0.format_name", ""),
		assert.ValueSet("file_format.0.json.#", "0"),
		assert.ValueSet("file_format.0.avro.#", "0"),
		assert.ValueSet("file_format.0.orc.#", "0"),
		assert.ValueSet("file_format.0.parquet.#", "0"),
		assert.ValueSet("file_format.0.xml.#", "0"),
	}
}
