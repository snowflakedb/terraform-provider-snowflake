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
	return r.WithAllowedAccountsValue(replicationGroupStringList(allowedAccounts))
}

func replicationGroupStringList(values []string) tfconfig.Variable {
	if len(values) == 0 {
		return config.EmptyListVariable()
	}
	return tfconfig.ListVariable(
		collections.Map(values, func(value string) tfconfig.Variable { return tfconfig.StringVariable(value) })...,
	)
}
