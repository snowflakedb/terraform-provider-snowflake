package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *ServiceResourceAssert) HasExternalAccessIntegrations(expected ...sdk.AccountObjectIdentifier) *ServiceResourceAssert {
	s.AddAssertion(assert.ValueSet("external_access_integrations.#", fmt.Sprintf("%d", len(expected))))
	for i, v := range expected {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("external_access_integrations.%d", i), v.FullyQualifiedName()))
	}
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTextNotEmpty() *ServiceResourceAssert {
	s.HasNoFromSpecificationTemplate()
	s.AddAssertion(assert.ValueSet("from_specification.#", "1"))
	s.AddAssertion(assert.ValueSet("from_specification.0.stage", ""))
	s.AddAssertion(assert.ValueSet("from_specification.0.path", ""))
	s.AddAssertion(assert.ValueSet("from_specification.0.file", ""))
	s.AddAssertion(assert.ValuePresent("from_specification.0.text"))
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationOnStageNotEmpty() *ServiceResourceAssert {
	s.HasNoFromSpecificationTemplate()
	s.AddAssertion(assert.ValueSet("from_specification.#", "1"))
	s.AddAssertion(assert.ValuePresent("from_specification.0.stage"))
	s.AddAssertion(assert.ValueSet("from_specification.0.path", ""))
	s.AddAssertion(assert.ValuePresent("from_specification.0.file"))
	s.AddAssertion(assert.ValueSet("from_specification.0.text", ""))
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateTextNotEmpty() *ServiceResourceAssert {
	s.HasNoFromSpecification()
	s.AddAssertion(assert.ValueSet("from_specification_template.#", "1"))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.stage", ""))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.path", ""))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.file", ""))
	s.AddAssertion(assert.ValuePresent("from_specification_template.0.text"))

	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateText(text string) *ServiceResourceAssert {
	s.HasNoFromSpecification()
	s.AddAssertion(assert.ValueSet("from_specification_template.#", "1"))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.stage", ""))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.path", ""))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.file", ""))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.text", text))
	s.AddAssertion(assert.ValuePresent("from_specification_template.0.using"))
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateUsing(using map[string]string) *ServiceResourceAssert {
	s.AddAssertion(assert.ValueSet("from_specification_template.#", "1"))
	for k, v := range using {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("from_specification_template.0.using.%s", k), v))
	}
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationEmpty() *ServiceResourceAssert {
	s.AddAssertion(assert.ValueNotSet("from_specification.#"))
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateEmpty() *ServiceResourceAssert {
	s.AddAssertion(assert.ValueNotSet("from_specification_template.#"))
	return s
}
