package sdk

func NewCreateOpenflowDeploymentRequest(name AccountObjectIdentifier) *CreateOpenflowDeploymentRequest {
	s := CreateOpenflowDeploymentRequest{}
	s.name = name
	return &s
}

func (s *CreateOpenflowDeploymentRequest) WithIfNotExists(ifNotExists bool) *CreateOpenflowDeploymentRequest {
	s.IfNotExists = &ifNotExists
	return s
}

func (s *CreateOpenflowDeploymentRequest) WithDeploymentType(deploymentType OpenflowDeploymentType) *CreateOpenflowDeploymentRequest {
	s.DeploymentType = &deploymentType
	return s
}

func (s *CreateOpenflowDeploymentRequest) WithVpcType(vpcType OpenflowVpcType) *CreateOpenflowDeploymentRequest {
	s.VpcType = &vpcType
	return s
}

func (s *CreateOpenflowDeploymentRequest) WithCustomIngressHostname(customIngressHostname string) *CreateOpenflowDeploymentRequest {
	s.CustomIngressHostname = &customIngressHostname
	return s
}

func (s *CreateOpenflowDeploymentRequest) WithUsePrivateLink(usePrivateLink bool) *CreateOpenflowDeploymentRequest {
	s.UsePrivateLink = &usePrivateLink
	return s
}

func (s *CreateOpenflowDeploymentRequest) WithUseUserAuthOverPrivatelink(useUserAuthOverPrivatelink bool) *CreateOpenflowDeploymentRequest {
	s.UseUserAuthOverPrivatelink = &useUserAuthOverPrivatelink
	return s
}

func (s *CreateOpenflowDeploymentRequest) WithEventTable(eventTable string) *CreateOpenflowDeploymentRequest {
	s.EventTable = &eventTable
	return s
}

func (s *CreateOpenflowDeploymentRequest) WithDisplayName(displayName string) *CreateOpenflowDeploymentRequest {
	s.DisplayName = &displayName
	return s
}

func (s *CreateOpenflowDeploymentRequest) WithComment(comment string) *CreateOpenflowDeploymentRequest {
	s.Comment = &comment
	return s
}

func NewAlterOpenflowDeploymentRequest(name AccountObjectIdentifier) *AlterOpenflowDeploymentRequest {
	s := AlterOpenflowDeploymentRequest{}
	s.name = name
	return &s
}

func (s *AlterOpenflowDeploymentRequest) WithSet(set OpenflowDeploymentSetRequest) *AlterOpenflowDeploymentRequest {
	s.Set = &set
	return s
}

func (s *AlterOpenflowDeploymentRequest) WithUnset(unset OpenflowDeploymentUnsetRequest) *AlterOpenflowDeploymentRequest {
	s.Unset = &unset
	return s
}

func NewOpenflowDeploymentSetRequest() *OpenflowDeploymentSetRequest {
	return &OpenflowDeploymentSetRequest{}
}

func (s *OpenflowDeploymentSetRequest) WithComment(comment string) *OpenflowDeploymentSetRequest {
	s.Comment = &comment
	return s
}

func (s *OpenflowDeploymentSetRequest) WithDisplayName(displayName string) *OpenflowDeploymentSetRequest {
	s.DisplayName = &displayName
	return s
}

func (s *OpenflowDeploymentSetRequest) WithEventTable(eventTable string) *OpenflowDeploymentSetRequest {
	s.EventTable = &eventTable
	return s
}

func NewOpenflowDeploymentUnsetRequest() *OpenflowDeploymentUnsetRequest {
	return &OpenflowDeploymentUnsetRequest{}
}

func (s *OpenflowDeploymentUnsetRequest) WithComment(comment bool) *OpenflowDeploymentUnsetRequest {
	s.Comment = &comment
	return s
}

func (s *OpenflowDeploymentUnsetRequest) WithDisplayName(displayName bool) *OpenflowDeploymentUnsetRequest {
	s.DisplayName = &displayName
	return s
}

func (s *OpenflowDeploymentUnsetRequest) WithEventTable(eventTable bool) *OpenflowDeploymentUnsetRequest {
	s.EventTable = &eventTable
	return s
}

func NewDropOpenflowDeploymentRequest(name AccountObjectIdentifier) *DropOpenflowDeploymentRequest {
	s := DropOpenflowDeploymentRequest{}
	s.name = name
	return &s
}

func (s *DropOpenflowDeploymentRequest) WithIfExists(ifExists bool) *DropOpenflowDeploymentRequest {
	s.IfExists = &ifExists
	return s
}

func NewShowOpenflowDeploymentRequest() *ShowOpenflowDeploymentRequest {
	return &ShowOpenflowDeploymentRequest{}
}

func (s *ShowOpenflowDeploymentRequest) WithLike(like Like) *ShowOpenflowDeploymentRequest {
	s.Like = &like
	return s
}

func NewDescribeOpenflowDeploymentRequest(name AccountObjectIdentifier) *DescribeOpenflowDeploymentRequest {
	s := DescribeOpenflowDeploymentRequest{}
	s.name = name
	return &s
}
