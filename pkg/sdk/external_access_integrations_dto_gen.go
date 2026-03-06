package sdk

var (
	_ optionsProvider[CreateExternalAccessIntegrationOptions]   = new(CreateExternalAccessIntegrationRequest)
	_ optionsProvider[AlterExternalAccessIntegrationOptions]    = new(AlterExternalAccessIntegrationRequest)
	_ optionsProvider[DropExternalAccessIntegrationOptions]     = new(DropExternalAccessIntegrationRequest)
	_ optionsProvider[ShowExternalAccessIntegrationOptions]     = new(ShowExternalAccessIntegrationRequest)
	_ optionsProvider[DescribeExternalAccessIntegrationOptions] = new(DescribeExternalAccessIntegrationRequest)
)

type CreateExternalAccessIntegrationRequest struct {
	OrReplace                    *bool
	IfNotExists                  *bool
	name                         AccountObjectIdentifier  // required
	AllowedNetworkRules          []SchemaObjectIdentifier // required
	AllowedAuthenticationSecrets []SchemaObjectIdentifier
	Enabled                      bool // required
	Comment                      *string
}

type AlterExternalAccessIntegrationRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
	Set      *ExternalAccessIntegrationSetRequest
	Unset    *ExternalAccessIntegrationUnsetRequest
}

type ExternalAccessIntegrationSetRequest struct {
	AllowedNetworkRules          []SchemaObjectIdentifier
	AllowedAuthenticationSecrets []SchemaObjectIdentifier
	Enabled                      *bool
	Comment                      *string
}

type ExternalAccessIntegrationUnsetRequest struct {
	AllowedAuthenticationSecrets *bool
	Comment                      *bool
}

type DropExternalAccessIntegrationRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowExternalAccessIntegrationRequest struct {
	Like *Like
}

type DescribeExternalAccessIntegrationRequest struct {
	name AccountObjectIdentifier // required
}
