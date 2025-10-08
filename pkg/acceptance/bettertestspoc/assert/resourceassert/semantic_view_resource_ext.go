package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *SemanticViewResourceAssert) HasTablesString(expected ...sdk.LogicalTable) *SemanticViewResourceAssert {
	s.AddAssertion(assert.ValueSet("tables.#", fmt.Sprintf("%d", len(expected))))
	return s
}

func (s *SemanticViewResourceAssert) HasMetricsString(expected ...sdk.MetricDefinition) *SemanticViewResourceAssert {
	s.AddAssertion(assert.ValueSet("metrics.#", fmt.Sprintf("%d", len(expected))))

	return s
}
