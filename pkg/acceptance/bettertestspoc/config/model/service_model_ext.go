package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func ServiceWithDefaultSpec(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
) *ServiceModel {
	spec := `
spec:
  containers:
  - name: example-container
    image: /snowflake/images/snowflake_images/exampleimage:latest
`
	s := &ServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.Service)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecification(spec)
	return s
}

func ServiceWithDefaultSpecOnStage(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
	stageId sdk.SchemaObjectIdentifier,
	fileName string,
) *ServiceModel {
	s := &ServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.Service)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecificationOnStage(stageId, fileName)
	return s
}

func ServiceWithDefaultSpecTemplate(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
) *ServiceModel {
	spec := `
spec:
  containers:
  - name: {{ container_name }}
    image: /snowflake/images/snowflake_images/exampleimage:latest
`
	s := &ServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.Service)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecificationTemplate(spec, map[string]tfconfig.Variable{
		"key":             tfconfig.StringVariable("container_name"),
		"value_in_quotes": tfconfig.StringVariable("example"),
	})
	return s
}

func (s *ServiceModel) WithFromSpecification(spec string) *ServiceModel {
	s.WithFromSpecificationValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"text": config.MultilineWrapperVariable(spec),
	}))
	return s
}

func (s *ServiceModel) WithFromSpecificationOnStage(stageId sdk.SchemaObjectIdentifier, fileName string) *ServiceModel {
	s.WithFromSpecificationValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"stage": tfconfig.StringVariable(stageId.FullyQualifiedName()),
		"file":  tfconfig.StringVariable(fileName),
	}))
	return s
}

func (s *ServiceModel) WithFromSpecificationTemplate(spec string, using map[string]tfconfig.Variable) *ServiceModel {
	s.WithFromSpecificationTemplateValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"text":  config.MultilineWrapperVariable(spec),
		"using": tfconfig.ObjectVariable(using),
	}))
	return s
}

func (f *ServiceModel) WithExternalAccessIntegrations(ids ...sdk.AccountObjectIdentifier) *ServiceModel {
	return f.WithExternalAccessIntegrationsValue(
		tfconfig.SetVariable(
			collections.Map(ids, func(id sdk.AccountObjectIdentifier) tfconfig.Variable {
				return tfconfig.StringVariable(id.FullyQualifiedName())
			})...,
		),
	)
}
