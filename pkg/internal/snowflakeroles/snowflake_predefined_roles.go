package snowflakeroles

import "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"

var (
	Orgadmin       = sdk.NewAccountObjectIdentifier("ORGADMIN")
	Accountadmin   = sdk.NewAccountObjectIdentifier("ACCOUNTADMIN")
	SecurityAdmin  = sdk.NewAccountObjectIdentifier("SECURITYADMIN")
	PentestingRole = sdk.NewAccountObjectIdentifier("PENTESTING_ROLE")

	OktaProvisioner        = sdk.NewAccountObjectIdentifier("OKTA_PROVISIONER")
	AadProvisioner         = sdk.NewAccountObjectIdentifier("AAD_PROVISIONER")
	GenericScimProvisioner = sdk.NewAccountObjectIdentifier("GENERIC_SCIM_PROVISIONER")
)
