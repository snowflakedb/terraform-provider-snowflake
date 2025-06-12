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
	s.HasFromSpecificationTemplateEmpty()
	s.AddAssertion(assert.ValueSet("from_specification.#", "1"))
	s.AddAssertion(assert.ValueSet("from_specification.0.stage", ""))
	s.AddAssertion(assert.ValueSet("from_specification.0.path", ""))
	s.AddAssertion(assert.ValueSet("from_specification.0.file", ""))
	s.AddAssertion(assert.ValuePresent("from_specification.0.text"))
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationOnStageNotEmpty() *ServiceResourceAssert {
	s.HasFromSpecificationTemplateEmpty()
	s.AddAssertion(assert.ValueSet("from_specification.#", "1"))
	s.AddAssertion(assert.ValuePresent("from_specification.0.stage"))
	s.AddAssertion(assert.ValueSet("from_specification.0.path", ""))
	s.AddAssertion(assert.ValuePresent("from_specification.0.file"))
	s.AddAssertion(assert.ValueSet("from_specification.0.text", ""))
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateText(using []map[string]string) *ServiceResourceAssert {
	s.HasFromSpecificationEmpty()
	s.AddAssertion(assert.ValueSet("from_specification_template.#", "1"))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.stage", ""))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.path", ""))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.file", ""))
	s.AddAssertion(assert.ValuePresent("from_specification_template.0.text"))
	s.HasFromSpecificationTemplateUsing(using)
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateOnStage(stageId sdk.SchemaObjectIdentifier, path string, fileName string, using []map[string]string) *ServiceResourceAssert {
	s.HasFromSpecificationEmpty()
	s.AddAssertion(assert.ValueSet("from_specification_template.#", "1"))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.stage", stageId.FullyQualifiedName()))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.path", path))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.file", fileName))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.text", ""))
	s.HasFromSpecificationTemplateUsing(using)
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateUsing(using []map[string]string) *ServiceResourceAssert {
	len := len(using)
	s.AddAssertion(assert.ValueSet("from_specification_template.#", fmt.Sprintf("%d", len)))
	for i, v := range using {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("from_specification_template.0.using.%d.key", i), v["key"]))
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("from_specification_template.0.using.%d.value", i), v["value"]))
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("from_specification_template.0.using.%d.value_in_quotes", i), v["value_in_quotes"]))
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("from_specification_template.0.using.%d.value_in_double_dollars", i), v["value_in_double_dollars"]))
	}
	return s
}
