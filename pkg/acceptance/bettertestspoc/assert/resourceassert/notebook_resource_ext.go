package resourceassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"

func (n *NotebookResourceAssert) HasFromPathAndStage(expectedPath string, expectedStageId string) *NotebookResourceAssert {
	n.AddAssertion(assert.ValueSet("from.0.path", expectedPath))
	n.AddAssertion(assert.ValueSet("from.0.stage", expectedStageId))
	return n
}
