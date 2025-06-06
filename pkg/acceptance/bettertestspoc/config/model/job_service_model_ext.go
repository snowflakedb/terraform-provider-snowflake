package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func JobServiceWithDefaultSpec(
	resourceName string,
	database string,
	schema string,
	name string,
	computePool string,
) *JobServiceModel {
	spec := `
spec:
  containers:
  - name: example-container
    image: /snowflake/images/snowflake_images/exampleimage:latest
`
	s := &JobServiceModel{ResourceModelMeta: config.Meta(resourceName, resources.JobService)}
	s.WithDatabase(database)
	s.WithSchema(schema)
	s.WithName(name)
	s.WithComputePool(computePool)
	s.WithFromSpecification(spec)
	return s
}

func JobServiceWithDefaultSpecOnStage(
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

func (f *JobServiceModel) WithExternalAccessIntegrations(ids ...sdk.AccountObjectIdentifier) *JobServiceModel {
	return f.WithExternalAccessIntegrationsValue(
		tfconfig.SetVariable(
			collections.Map(ids, func(id sdk.AccountObjectIdentifier) tfconfig.Variable {
				return tfconfig.StringVariable(id.FullyQualifiedName())
			})...,
		),
	)
}
