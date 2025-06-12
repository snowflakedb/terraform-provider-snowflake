package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func JobServiceWithSpec(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
	spec string,
) *JobServiceModel {
	s := &JobServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.JobService)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecification(spec)
	return s
}

func JobServiceWithSpecOnStage(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
	stageId sdk.SchemaObjectIdentifier,
	fileName string,
) *JobServiceModel {
	s := &JobServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.JobService)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecificationOnStage(stageId, fileName)
	return s
}

func JobServiceWithSpecTemplate(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
	specTemplate string,
	using []helpers.ServiceSpecUsing,
) *JobServiceModel {
	s := &JobServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.JobService)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecificationTemplate(specTemplate, using...)
	return s
}

func JobServiceWithSpecTemplateOnStage(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
	stageId sdk.SchemaObjectIdentifier,
	fileName string,
	using []helpers.ServiceSpecUsing,
) *JobServiceModel {
	s := &JobServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.JobService)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecificationTemplateOnStage(stageId, fileName, using)
	return s
}

func (s *JobServiceModel) WithFromSpecification(spec string) *JobServiceModel {
	s.WithFromSpecificationValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"text": config.MultilineWrapperVariable(spec),
	}))
	return s
}

func (s *JobServiceModel) WithFromSpecificationOnStage(stageId sdk.SchemaObjectIdentifier, fileName string) *JobServiceModel {
	s.WithFromSpecificationValue(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"stage": tfconfig.StringVariable(stageId.FullyQualifiedName()),
		"file":  tfconfig.StringVariable(fileName),
	}))
	return s
}

func (s *JobServiceModel) WithFromSpecificationTemplate(spec string, using ...helpers.ServiceSpecUsing) *JobServiceModel {
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

func (s *JobServiceModel) WithFromSpecificationTemplateOnStage(stageId sdk.SchemaObjectIdentifier, fileName string, using []helpers.ServiceSpecUsing) *JobServiceModel {
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

func (f *JobServiceModel) WithExternalAccessIntegrations(ids ...sdk.AccountObjectIdentifier) *JobServiceModel {
	return f.WithExternalAccessIntegrationsValue(
		tfconfig.SetVariable(
			collections.Map(ids, func(id sdk.AccountObjectIdentifier) tfconfig.Variable {
				return tfconfig.StringVariable(id.FullyQualifiedName())
			})...,
		),
	)
}
