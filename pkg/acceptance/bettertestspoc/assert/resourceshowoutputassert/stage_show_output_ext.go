package resourceshowoutputassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
)

func StagesDatasourceShowOutput(t *testing.T, datasourceReference string) *StageShowOutputAssert {
	t.Helper()
	s := StageShowOutputAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "show_output", "stages.0."),
	}
	s.AddAssertion(assert.ValueSet("show_output.#", "1"))
	return &s
}

func (s *StageShowOutputAssert) HasCreatedOnNotEmpty() *StageShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValuePresent("created_on"))
	return s
}

func (s *StageShowOutputAssert) HasCommentEmpty() *StageShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("comment", ""))
	return s
}

func (s *StageShowOutputAssert) HasStorageIntegrationEmpty() *StageShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputStringUnderlyingValueSet("storage_integration", ""))
	return s
}

func (s *StageShowOutputAssert) HasRegionEmpty() *StageShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("region", ""))
	return s
}

func (s *StageShowOutputAssert) HasCloudEmpty() *StageShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("cloud", ""))
	return s
}

func (s *StageShowOutputAssert) HasEndpointEmpty() *StageShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("endpoint", ""))
	return s
}

func (s *StageShowOutputAssert) HasUrlEmpty() *StageShowOutputAssert {
	s.AddAssertion(assert.ResourceShowOutputValueSet("url", ""))
	return s
}
