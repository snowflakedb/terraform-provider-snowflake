package datasourcemodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (s *StorageLifecyclePoliciesModel) WithEmptyIn() *StorageLifecyclePoliciesModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{}),
	)
}

func (s *StorageLifecyclePoliciesModel) WithInDatabase(databaseId sdk.AccountObjectIdentifier) *StorageLifecyclePoliciesModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"database": tfconfig.StringVariable(databaseId.Name()),
		}),
	)
}

func (s *StorageLifecyclePoliciesModel) WithInSchema(schemaId sdk.DatabaseObjectIdentifier) *StorageLifecyclePoliciesModel {
	return s.WithInValue(
		tfconfig.ObjectVariable(map[string]tfconfig.Variable{
			"schema": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
		}),
	)
}
