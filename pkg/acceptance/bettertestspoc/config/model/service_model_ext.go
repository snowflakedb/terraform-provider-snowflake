package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func ServiceWithSpec(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
	spec string,
) *ServiceModel {
	s := &ServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.Service)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecification(spec)
	return s
}

func ServiceWithSpecOnStage(
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

func ServiceWithSpecTemplate(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
	specTemplate string,
	using ...helpers.ServiceSpecUsing,
) *ServiceModel {
	s := &ServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.Service)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecificationTemplate(specTemplate, using...)
	return s
}

func ServiceWithSpecTemplateOnStage(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
	stageId sdk.SchemaObjectIdentifier,
	fileName string,
	using []helpers.ServiceSpecUsing,
) *ServiceModel {
	s := &ServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.Service)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecificationTemplateOnStage(stageId, fileName, using)
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

func (s *ServiceModel) WithFromSpecificationTemplate(spec string, using ...helpers.ServiceSpecUsing) *ServiceModel {
	s.WithFromSpecificationTemplateValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"text": config.MultilineWrapperVariable(spec),
		"using": tfconfig.SetVariable(
			collections.Map(using, func(item helpers.ServiceSpecUsing) tfconfig.Variable {
				return item.ToTfVariable()
			})...,
		),
	}))
	return s
}

func (s *ServiceModel) WithFromSpecificationTemplateOnStage(stageId sdk.SchemaObjectIdentifier, fileName string, using []helpers.ServiceSpecUsing) *ServiceModel {
	s.WithFromSpecificationTemplateValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"stage": tfconfig.StringVariable(stageId.FullyQualifiedName()),
		"file":  tfconfig.StringVariable(fileName),
		"using": tfconfig.SetVariable(
			collections.Map(using, func(item helpers.ServiceSpecUsing) tfconfig.Variable {
				return item.ToTfVariable()
			})...,
		),
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
