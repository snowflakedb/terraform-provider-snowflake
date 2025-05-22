package provider

import "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"

type Context struct {
	Client          *sdk.Client
	EnabledFeatures []string
}
