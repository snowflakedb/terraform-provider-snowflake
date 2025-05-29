package providermodel

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type SnowflakeTomlModel struct {
	sdk.ConfigDTO
	*config.ProviderModelMeta
}

func SnowflakeProviderToml(profile string) *SnowflakeTomlModel {
	s := &SnowflakeTomlModel{ProviderModelMeta: config.DefaultProviderMeta(profile)}
	return s
}
