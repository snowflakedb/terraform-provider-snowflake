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

func ServiceWithSpecTemplateOnStage(
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
	s.WithFromSpecificationTemplateOnStage(stageId, fileName, []map[string]tfconfig.Variable{
		{
			"key":             tfconfig.StringVariable("key"),
			"value_in_quotes": tfconfig.StringVariable("valueinquotes"),
		},
	})
	return s
}

func ServiceWithSpecTemplate(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
	specTemplate string,
	using []helpers.ServiceSpecUsing,
) *ServiceModel {
	s := &ServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.Service)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecificationTemplateRaw(specTemplate, using)
	return s
}

func (s *ServiceModel) WithFromSpecification(spec string) *ServiceModel {
	s.WithFromSpecificationValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"text": config.MultilineWrapperVariable(spec),
	}))
	return s
}

func (s *ServiceModel) WithFromSpecificationTemplate(spec string, using ...sdk.ListItem) *ServiceModel {
	s.WithFromSpecificationTemplateValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"text": config.MultilineWrapperVariable(spec),
		"using": tfconfig.SetVariable(
			collections.Map(using, func(item sdk.ListItem) tfconfig.Variable {
				v := item.Value.(string)
				return tfconfig.ObjectVariable(map[string]tfconfig.Variable{
					"key":   tfconfig.StringVariable(item.Key),
					"value": tfconfig.StringVariable(v),
				})
			})...,
		),
	}))
	return s
}

func (s *ServiceModel) WithFromSpecificationTemplateRaw(spec string, using []helpers.ServiceSpecUsing) *ServiceModel {
	usingRaw := collections.Map(using, func(item helpers.ServiceSpecUsing) map[string]tfconfig.Variable {
		usingItem := map[string]tfconfig.Variable{
			"key": tfconfig.StringVariable(item.Key),
		}
		if item.Value != nil {
			usingItem["value"] = tfconfig.StringVariable(*item.Value)
		}
		if item.ValueInQuotes != nil {
			usingItem["value_in_quotes"] = tfconfig.StringVariable(*item.ValueInQuotes)
		}
		if item.ValueInDoubleDollars != nil {
			usingItem["value_in_double_dollars"] = tfconfig.StringVariable(*item.ValueInDoubleDollars)
		}
		return usingItem
	})
	s.WithFromSpecificationTemplateValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"text": config.MultilineWrapperVariable(spec),
		"using": tfconfig.SetVariable(
			collections.Map(usingRaw, func(item map[string]tfconfig.Variable) tfconfig.Variable {
				return tfconfig.ObjectVariable(item)
			})...,
		),
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

func (s *ServiceModel) WithFromSpecificationTemplateOnStage(stageId sdk.SchemaObjectIdentifier, fileName string, using []map[string]tfconfig.Variable) *ServiceModel {
	s.WithFromSpecificationTemplateValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"stage": tfconfig.StringVariable(stageId.FullyQualifiedName()),
		"file":  tfconfig.StringVariable(fileName),
		"using": tfconfig.SetVariable(
			collections.Map(using, func(item map[string]tfconfig.Variable) tfconfig.Variable {
				return tfconfig.ObjectVariable(item)
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

func ServiceDynamicUsing(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
) *ServiceModel {
	m := &ServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.Service)}
	m.WithDatabase(database)
	m.WithSchema(schema)
	m.WithName(name)
	m.WithComputePool(computePool)
	return m
}

func (s *ServiceModel) WithTextAndUsing(specTemplate string, using map[string]tfconfig.Variable) *ServiceModel {
	s.WithFromSpecificationTemplateValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"text":  config.MultilineWrapperVariable(specTemplate),
		"using": tfconfig.ObjectVariable(using),
	}))
	return s
}
