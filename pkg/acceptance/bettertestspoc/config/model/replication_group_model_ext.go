package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (r *ReplicationGroupModel) WithObjectTypes(objectTypes []string) *ReplicationGroupModel {
	return r.WithObjectTypesValue(replicationGroupStringList(objectTypes))
}

func (r *ReplicationGroupModel) WithAllowedAccounts(allowedAccounts []string) *ReplicationGroupModel {
	return r.WithAllowedAccountsValue(replicationGroupStringSet(allowedAccounts))
}

func (r *ReplicationGroupModel) WithAllowedDatabases(allowedDatabases []string) *ReplicationGroupModel {
	return r.WithAllowedDatabasesValue(replicationGroupStringList(allowedDatabases))
}

func (r *ReplicationGroupModel) WithAllowedExternalVolumes(allowedExternalVolumes []string) *ReplicationGroupModel {
	return r.WithAllowedExternalVolumesValue(replicationGroupStringList(allowedExternalVolumes))
}

func (r *ReplicationGroupModel) WithAllowedShares(allowedShares []string) *ReplicationGroupModel {
	return r.WithAllowedSharesValue(replicationGroupStringList(allowedShares))
}

func (r *ReplicationGroupModel) WithAllowedIntegrationTypes(allowedIntegrationTypes []string) *ReplicationGroupModel {
	return r.WithAllowedIntegrationTypesValue(replicationGroupStringList(allowedIntegrationTypes))
}

func replicationGroupStringSet(values []string) tfconfig.Variable {
	return tfconfig.SetVariable(
		collections.Map(values, func(value string) tfconfig.Variable { return tfconfig.StringVariable(value) })...,
	)
}

func replicationGroupStringList(values []string) tfconfig.Variable {
	if len(values) == 0 {
		return config.EmptyListVariable()
	}
	return tfconfig.ListVariable(
		collections.Map(values, func(value string) tfconfig.Variable { return tfconfig.StringVariable(value) })...,
	)
}
