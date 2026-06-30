package resourceassert

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *ServiceResourceAssert) HasExternalAccessIntegrationsIdentifier(expected ...sdk.AccountObjectIdentifier) *ServiceResourceAssert {
	return s.HasExternalAccessIntegrations(collections.Map(expected, func(v sdk.AccountObjectIdentifier) string {
		return v.FullyQualifiedName()
	})...)
}

func (s *ServiceResourceAssert) HasFromSpecificationTextNotEmpty() *ServiceResourceAssert {
	s.HasFromSpecificationTemplateEmpty()
	s.ValueSet("from_specification.#", "1")
	s.ValueSet("from_specification.0.stage", "")
	s.ValueSet("from_specification.0.path", "")
	s.ValueSet("from_specification.0.file", "")
	s.ValuePresent("from_specification.0.text")
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationOnStage(stageId sdk.SchemaObjectIdentifier, path, fileName string) *ServiceResourceAssert {
	s.HasFromSpecificationTemplateEmpty()
	s.ValueSet("from_specification.#", "1")
	s.ValueSet("from_specification.0.stage", stageId.FullyQualifiedName())
	s.ValueSet("from_specification.0.path", path)
	s.ValueSet("from_specification.0.file", fileName)
	s.ValueSet("from_specification.0.text", "")
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateTextNotEmpty(using ...helpers.ServiceSpecUsing) *ServiceResourceAssert {
	s.HasFromSpecificationEmpty()
	s.ValueSet("from_specification_template.#", "1")
	s.ValueSet("from_specification_template.0.stage", "")
	s.ValueSet("from_specification_template.0.path", "")
	s.ValueSet("from_specification_template.0.file", "")
	s.ValuePresent("from_specification_template.0.text")
	s.HasFromSpecificationTemplateUsing(using...)
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateOnStage(stageId sdk.SchemaObjectIdentifier, path string, fileName string, using ...helpers.ServiceSpecUsing) *ServiceResourceAssert {
	s.HasFromSpecificationEmpty()
	s.ValueSet("from_specification_template.#", "1")
	s.ValueSet("from_specification_template.0.stage", stageId.FullyQualifiedName())
	s.ValueSet("from_specification_template.0.path", path)
	s.ValueSet("from_specification_template.0.file", fileName)
	s.ValueSet("from_specification_template.0.text", "")
	s.HasFromSpecificationTemplateUsing(using...)
	return s
}

func (s *ServiceResourceAssert) HasFromSpecificationTemplateUsing(using ...helpers.ServiceSpecUsing) *ServiceResourceAssert {
	s.ValueSet("from_specification_template.0.using.#", fmt.Sprintf("%d", len(using)))
	for i, v := range using {
		s.ValueSet(fmt.Sprintf("from_specification_template.0.using.%d.key", i), v.Key)
		s.ValueSet(fmt.Sprintf("from_specification_template.0.using.%d.value", i), v.Value)
	}
	return s
}
