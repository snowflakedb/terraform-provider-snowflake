package model

import (
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (s *StorageIntegrationGcsModel) WithStorageAllowedLocations(allowedLocations []sdk.StorageLocation) *StorageIntegrationGcsModel {
	allowedLocationsStringVariables := collections.Map(allowedLocations, func(location sdk.StorageLocation) tfconfig.Variable { return tfconfig.StringVariable(location.Path) })
	s.WithStorageAllowedLocationsValue(tfconfig.ListVariable(allowedLocationsStringVariables...))
	return s
}

func (s *StorageIntegrationGcsModel) WithStorageBlockedLocations(blockedLocations []sdk.StorageLocation) *StorageIntegrationGcsModel {
	blockedLocationsStringVariables := collections.Map(blockedLocations, func(location sdk.StorageLocation) tfconfig.Variable { return tfconfig.StringVariable(location.Path) })
	s.WithStorageAllowedLocationsValue(tfconfig.ListVariable(blockedLocationsStringVariables...))
	return s
}
