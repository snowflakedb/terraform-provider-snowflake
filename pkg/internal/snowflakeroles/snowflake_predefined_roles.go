package snowflakeroles

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

var (
	GlobalOrgAdmin = sdk.NewAccountObjectIdentifier("GLOBALORGADMIN")
	Orgadmin       = sdk.NewAccountObjectIdentifier("ORGADMIN")
	Accountadmin   = sdk.NewAccountObjectIdentifier("ACCOUNTADMIN")
	SecurityAdmin  = sdk.NewAccountObjectIdentifier("SECURITYADMIN")
	SysAdmin       = sdk.NewAccountObjectIdentifier("SYSADMIN")
	UserAdmin      = sdk.NewAccountObjectIdentifier("USERADMIN")
	Public         = sdk.NewAccountObjectIdentifier("PUBLIC")
	PentestingRole = sdk.NewAccountObjectIdentifier("PENTESTING_ROLE")
	// RESTRICTED is a role that has no grants. It can be used instead of PUBLIC role, so that it's not granted to any users by default.
	Restricted = sdk.NewAccountObjectIdentifier("RESTRICTED")

	OktaProvisioner        = sdk.NewAccountObjectIdentifier("OKTA_PROVISIONER")
	AadProvisioner         = sdk.NewAccountObjectIdentifier("AAD_PROVISIONER")
	GenericScimProvisioner = sdk.NewAccountObjectIdentifier("GENERIC_SCIM_PROVISIONER")
)
