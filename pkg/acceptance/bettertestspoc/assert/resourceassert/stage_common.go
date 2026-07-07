package resourceassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func stageApplyFileFormatFormatNameChecks(e *assert.ResourceAssert, expected string) {
	e.ValueSet("file_format.#", "1")
	e.ValueSet("file_format.0.format_name", expected)
	e.ValueSet("file_format.0.csv.#", "0")
	e.ValueSet("file_format.0.json.#", "0")
	e.ValueSet("file_format.0.avro.#", "0")
	e.ValueSet("file_format.0.orc.#", "0")
	e.ValueSet("file_format.0.parquet.#", "0")
	e.ValueSet("file_format.0.xml.#", "0")
}

func stageApplyFileFormatCsvChecks(e *assert.ResourceAssert) {
	e.ValueSet("file_format.#", "1")
	e.ValueSet("file_format.0.csv.#", "1")
	e.ValueSet("file_format.0.format_name", "")
	e.ValueSet("file_format.0.json.#", "0")
	e.ValueSet("file_format.0.avro.#", "0")
	e.ValueSet("file_format.0.orc.#", "0")
	e.ValueSet("file_format.0.parquet.#", "0")
	e.ValueSet("file_format.0.xml.#", "0")
}
