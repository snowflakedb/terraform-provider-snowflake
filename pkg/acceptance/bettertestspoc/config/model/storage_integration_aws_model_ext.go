package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationAwsModel) WithStorageAllowedLocations(allowedLocations []sdk.StorageLocation) *StorageIntegrationAwsModel {
	allowedLocationsStringVariables := collections.Map(allowedLocations, func(location sdk.StorageLocation) tfconfig.Variable { return tfconfig.StringVariable(location.Path) })
	s.WithStorageAllowedLocationsValue(tfconfig.ListVariable(allowedLocationsStringVariables...))
	return s
}

func (s *StorageIntegrationAwsModel) WithStorageBlockedLocations(blockedLocations []sdk.StorageLocation) *StorageIntegrationAwsModel {
	blockedLocationsStringVariables := collections.Map(blockedLocations, func(location sdk.StorageLocation) tfconfig.Variable { return tfconfig.StringVariable(location.Path) })
	s.WithStorageAllowedLocationsValue(tfconfig.ListVariable(blockedLocationsStringVariables...))
	return s
}
