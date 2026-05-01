package sdk

var (
	_ optionsProvider[CreateOpenflowRuntimeOptions]   = new(CreateOpenflowRuntimeRequest)
	_ optionsProvider[AlterOpenflowRuntimeOptions]    = new(AlterOpenflowRuntimeRequest)
	_ optionsProvider[DropOpenflowRuntimeOptions]     = new(DropOpenflowRuntimeRequest)
	_ optionsProvider[ShowOpenflowRuntimeOptions]     = new(ShowOpenflowRuntimeRequest)
	_ optionsProvider[DescribeOpenflowRuntimeOptions] = new(DescribeOpenflowRuntimeRequest)
)

type CreateOpenflowRuntimeRequest struct {
	name                       SchemaObjectIdentifier  // required
	InDeployment               AccountObjectIdentifier // required
	ExecuteAsRole              string                  // required
	NodeType                   OpenflowRuntimeNodeType // required
	MinNodes                   int                     // required
	MaxNodes                   int                     // required
	IfNotExists                *bool
	ExternalAccessIntegrations []AccountObjectIdentifier
	DisplayName                *string
	Comment                    *string
}

type AlterOpenflowRuntimeRequest struct {
	name             SchemaObjectIdentifier // required
	Suspend          *bool
	Resume           *bool
	Terminate        *bool
	TerminateCascade *bool
	Set              *OpenflowRuntimeSetRequest
	Unset            *OpenflowRuntimeUnsetRequest
}

type OpenflowRuntimeSetRequest struct {
	MinNodes                   *int
	MaxNodes                   *int
	ExecuteAsRole              *string
	ExternalAccessIntegrations []AccountObjectIdentifier
	DisplayName                *string
	Comment                    *string
}

type OpenflowRuntimeUnsetRequest struct {
	DisplayName *bool
	Comment     *bool
}

type DropOpenflowRuntimeRequest struct {
	name     SchemaObjectIdentifier // required
	IfExists *bool
	Cascade  *bool
}

type ShowOpenflowRuntimeRequest struct {
	Like *Like
}

type DescribeOpenflowRuntimeRequest struct {
	name SchemaObjectIdentifier // required
}
