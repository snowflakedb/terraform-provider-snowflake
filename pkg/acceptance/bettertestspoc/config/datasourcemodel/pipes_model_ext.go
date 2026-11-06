package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (p *PipesModel) WithEmptyIn() *PipesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
		}),
	)
}

func (p *PipesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *PipesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (p *PipesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *PipesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}

func (p *PipesModel) WithInApplication(applicationId sdk.AccountObjectIdentifier) *PipesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"application": tfconfig.StringVariable(applicationId.Name()),
		}),
	)
}

func (p *PipesModel) WithInApplicationPackage(applicationPackageId sdk.AccountObjectIdentifier) *PipesModel {
	return p.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"application_package": tfconfig.StringVariable(applicationPackageId.Name()),
		}),
	)
}
