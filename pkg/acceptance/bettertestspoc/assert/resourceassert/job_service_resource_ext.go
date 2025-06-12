package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *JobServiceResourceAssert) HasExternalAccessIntegrations(expected ...sdk.AccountObjectIdentifier) *JobServiceResourceAssert {
	s.AddAssertion(assert.ValueSet("external_access_integrations.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("external_access_integrations.%d", i), v.FullyQualifiedName()))
	}
	return s
}

func (s *JobServiceResourceAssert) HasFromSpecificationTextNotEmpty() *JobServiceResourceAssert {
	s.HasFromSpecificationTemplateEmpty()
	s.AddAssertion(assert.ValueSet("from_specification.#", "1"))
	s.AddAssertion(assert.ValueSet("from_specification.0.stage", ""))
	s.AddAssertion(assert.ValueSet("from_specification.0.path", ""))
	s.AddAssertion(assert.ValueSet("from_specification.0.file", ""))
	s.AddAssertion(assert.ValuePresent("from_specification.0.text"))
	return s
}

func (s *JobServiceResourceAssert) HasFromSpecificationOnStageNotEmpty() *JobServiceResourceAssert {
	s.HasFromSpecificationTemplateEmpty()
	s.AddAssertion(assert.ValueSet("from_specification.#", "1"))
	s.AddAssertion(assert.ValuePresent("from_specification.0.stage"))
	s.AddAssertion(assert.ValueSet("from_specification.0.path", ""))
	s.AddAssertion(assert.ValuePresent("from_specification.0.file"))
	s.AddAssertion(assert.ValueSet("from_specification.0.text", ""))
	return s
}
