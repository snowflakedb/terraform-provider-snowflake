package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
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

func (s *ServiceResourceAssert) HasFromSpecificationOnStage(stageId sdk.SchemaObjectIdentifier, path, fileName string) *ServiceResourceAssert {
	s.HasFromSpecificationTemplateEmpty()
	s.AddAssertion(assert.ValueSet("from_specification.#", "1"))
	s.AddAssertion(assert.ValueSet("from_specification.0.stage", stageId.FullyQualifiedName()))
	s.AddAssertion(assert.ValueSet("from_specification.0.path", path))
	s.AddAssertion(assert.ValueSet("from_specification.0.file", fileName))
	s.AddAssertion(assert.ValueSet("from_specification.0.text", ""))
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateTextNotEmpty(using ...helpers.ServiceSpecUsing) *ServiceResourceAssert {
	s.HasFromSpecificationEmpty()
	s.AddAssertion(assert.ValueSet("from_specification_template.#", "1"))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.stage", ""))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.path", ""))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.file", ""))
	s.AddAssertion(assert.ValuePresent("from_specification_template.0.text"))
	s.HasFromSpecificationTemplateUsing(using...)
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateOnStage(stageId sdk.SchemaObjectIdentifier, path string, fileName string, using ...helpers.ServiceSpecUsing) *ServiceResourceAssert {
	s.HasFromSpecificationEmpty()
	s.AddAssertion(assert.ValueSet("from_specification_template.#", "1"))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.stage", stageId.FullyQualifiedName()))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.path", path))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.file", fileName))
	s.AddAssertion(assert.ValueSet("from_specification_template.0.text", ""))
	s.HasFromSpecificationTemplateUsing(using...)
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateUsing(using ...helpers.ServiceSpecUsing) *ServiceResourceAssert {
	s.AddAssertion(assert.ValueSet("from_specification_template.0.using.#", fmt.Sprintf("%d", len(using))))
	for i, v := range using {
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("from_specification_template.0.using.%d.key", i), v.Key))
		s.AddAssertion(assert.ValueSet(fmt.Sprintf("from_specification_template.0.using.%d.value", i), v.Value))
	}
	return s
}
