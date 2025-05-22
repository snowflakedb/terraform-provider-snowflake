package resources

import "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"

var DeleteSecurityIntegration = ResourceDeleteContextFunc(
	sdk.ParseAccountObjectIdentifier,
	func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
		return client.SecurityIntegrations.DropSafely
	},
)
