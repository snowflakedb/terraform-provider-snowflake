package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SdkV2ProviderSchema map[string]*schema.Schema

func (s SdkV2ProviderSchema) ObjectName() string {
	return "Snowflake"
}

// SdkV2ProviderSchemas contains all provider schemas
var SdkV2ProviderSchemas = []SdkV2ProviderSchema{
	provider.GetProviderSchema(),
}
