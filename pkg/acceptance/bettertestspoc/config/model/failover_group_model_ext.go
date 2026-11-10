package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (f *FailoverGroupModel) WithObjectTypes(objectTypes []string) *FailoverGroupModel {
	return f.WithObjectTypesValue(failoverGroupStringList(objectTypes))
}

func (f *FailoverGroupModel) WithAllowedAccounts(allowedAccounts []string) *FailoverGroupModel {
	return f.WithAllowedAccountsValue(failoverGroupStringSet(allowedAccounts))
}

func (f *FailoverGroupModel) WithAllowedDatabases(allowedDatabases []string) *FailoverGroupModel {
	return f.WithAllowedDatabasesValue(failoverGroupStringList(allowedDatabases))
}

func (f *FailoverGroupModel) WithAllowedExternalVolumes(allowedExternalVolumes []string) *FailoverGroupModel {
	return f.WithAllowedExternalVolumesValue(failoverGroupStringList(allowedExternalVolumes))
}

func (f *FailoverGroupModel) WithAllowedShares(allowedShares []string) *FailoverGroupModel {
	return f.WithAllowedSharesValue(failoverGroupStringList(allowedShares))
}

func (f *FailoverGroupModel) WithAllowedIntegrationTypes(allowedIntegrationTypes []string) *FailoverGroupModel {
	return f.WithAllowedIntegrationTypesValue(failoverGroupStringList(allowedIntegrationTypes))
}

func failoverGroupStringSet(values []string) tfconfig.Variable {
	return tfconfig.SetVariable(
		collections.Map(values, func(value string) tfconfig.Variable { return tfconfig.StringVariable(value) })...,
	)
}

func failoverGroupStringList(values []string) tfconfig.Variable {
	if len(values) == 0 {
		return config.EmptyListVariable()
	}
	return tfconfig.ListVariable(
		collections.Map(values, func(value string) tfconfig.Variable { return tfconfig.StringVariable(value) })...,
	)
}
