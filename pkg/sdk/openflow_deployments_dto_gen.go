package sdk

var (
	_ optionsProvider[CreateOpenflowDeploymentOptions]   = new(CreateOpenflowDeploymentRequest)
	_ optionsProvider[AlterOpenflowDeploymentOptions]    = new(AlterOpenflowDeploymentRequest)
	_ optionsProvider[DropOpenflowDeploymentOptions]     = new(DropOpenflowDeploymentRequest)
	_ optionsProvider[ShowOpenflowDeploymentOptions]     = new(ShowOpenflowDeploymentRequest)
	_ optionsProvider[DescribeOpenflowDeploymentOptions] = new(DescribeOpenflowDeploymentRequest)
)

type CreateOpenflowDeploymentRequest struct {
	name                      AccountObjectIdentifier // required
	IfNotExists               *bool
	DeploymentType            *OpenflowDeploymentType
	VpcType                   *OpenflowVpcType
	CustomIngressHostname     *string
	UsePrivateLink            *bool
	UseUserAuthOverPrivatelink *bool
	EventTable                *string
	DisplayName               *string
	Comment                   *string
}

type AlterOpenflowDeploymentRequest struct {
	name  AccountObjectIdentifier // required
	Set   *OpenflowDeploymentSetRequest
	Unset *OpenflowDeploymentUnsetRequest
}

type OpenflowDeploymentSetRequest struct {
	Comment     *string
	DisplayName *string
	EventTable  *string
}

type OpenflowDeploymentUnsetRequest struct {
	Comment     *bool
	DisplayName *bool
	EventTable  *bool
}

type DropOpenflowDeploymentRequest struct {
	name     AccountObjectIdentifier // required
	IfExists *bool
}

type ShowOpenflowDeploymentRequest struct {
	Like *Like
}

type DescribeOpenflowDeploymentRequest struct {
	name AccountObjectIdentifier // required
}
