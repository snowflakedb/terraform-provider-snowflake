package sdk

func NewCreateOpenflowRuntimeRequest(
	name SchemaObjectIdentifier,
	inDeployment AccountObjectIdentifier,
	executeAsRole string,
	nodeType OpenflowRuntimeNodeType,
	minNodes int,
	maxNodes int,
) *CreateOpenflowRuntimeRequest {
	s := CreateOpenflowRuntimeRequest{}
	s.name = name
	s.InDeployment = inDeployment
	s.ExecuteAsRole = executeAsRole
	s.NodeType = nodeType
	s.MinNodes = minNodes
	s.MaxNodes = maxNodes
	return &s
}

func (s *CreateOpenflowRuntimeRequest) WithIfNotExists(ifNotExists bool) *CreateOpenflowRuntimeRequest {
	s.IfNotExists = &ifNotExists
	return s
}

func (s *CreateOpenflowRuntimeRequest) WithExternalAccessIntegrations(eais []AccountObjectIdentifier) *CreateOpenflowRuntimeRequest {
	s.ExternalAccessIntegrations = eais
	return s
}

func (s *CreateOpenflowRuntimeRequest) WithDisplayName(displayName string) *CreateOpenflowRuntimeRequest {
	s.DisplayName = &displayName
	return s
}

func (s *CreateOpenflowRuntimeRequest) WithComment(comment string) *CreateOpenflowRuntimeRequest {
	s.Comment = &comment
	return s
}

func NewAlterOpenflowRuntimeRequest(name SchemaObjectIdentifier) *AlterOpenflowRuntimeRequest {
	s := AlterOpenflowRuntimeRequest{}
	s.name = name
	return &s
}

func (s *AlterOpenflowRuntimeRequest) WithSuspend(suspend bool) *AlterOpenflowRuntimeRequest {
	s.Suspend = &suspend
	return s
}

func (s *AlterOpenflowRuntimeRequest) WithResume(resume bool) *AlterOpenflowRuntimeRequest {
	s.Resume = &resume
	return s
}

func (s *AlterOpenflowRuntimeRequest) WithTerminate(terminate bool) *AlterOpenflowRuntimeRequest {
	s.Terminate = &terminate
	return s
}

func (s *AlterOpenflowRuntimeRequest) WithTerminateCascade(terminateCascade bool) *AlterOpenflowRuntimeRequest {
	s.TerminateCascade = &terminateCascade
	return s
}

func (s *AlterOpenflowRuntimeRequest) WithSet(set OpenflowRuntimeSetRequest) *AlterOpenflowRuntimeRequest {
	s.Set = &set
	return s
}

func (s *AlterOpenflowRuntimeRequest) WithUnset(unset OpenflowRuntimeUnsetRequest) *AlterOpenflowRuntimeRequest {
	s.Unset = &unset
	return s
}

func NewOpenflowRuntimeSetRequest() *OpenflowRuntimeSetRequest {
	return &OpenflowRuntimeSetRequest{}
}

func (s *OpenflowRuntimeSetRequest) WithMinNodes(minNodes int) *OpenflowRuntimeSetRequest {
	s.MinNodes = &minNodes
	return s
}

func (s *OpenflowRuntimeSetRequest) WithMaxNodes(maxNodes int) *OpenflowRuntimeSetRequest {
	s.MaxNodes = &maxNodes
	return s
}

func (s *OpenflowRuntimeSetRequest) WithExecuteAsRole(role string) *OpenflowRuntimeSetRequest {
	s.ExecuteAsRole = &role
	return s
}

func (s *OpenflowRuntimeSetRequest) WithExternalAccessIntegrations(eais []AccountObjectIdentifier) *OpenflowRuntimeSetRequest {
	s.ExternalAccessIntegrations = eais
	return s
}

func (s *OpenflowRuntimeSetRequest) WithDisplayName(displayName string) *OpenflowRuntimeSetRequest {
	s.DisplayName = &displayName
	return s
}

func (s *OpenflowRuntimeSetRequest) WithComment(comment string) *OpenflowRuntimeSetRequest {
	s.Comment = &comment
	return s
}

func NewOpenflowRuntimeUnsetRequest() *OpenflowRuntimeUnsetRequest {
	return &OpenflowRuntimeUnsetRequest{}
}

func (s *OpenflowRuntimeUnsetRequest) WithDisplayName(displayName bool) *OpenflowRuntimeUnsetRequest {
	s.DisplayName = &displayName
	return s
}

func (s *OpenflowRuntimeUnsetRequest) WithComment(comment bool) *OpenflowRuntimeUnsetRequest {
	s.Comment = &comment
	return s
}

func NewDropOpenflowRuntimeRequest(name SchemaObjectIdentifier) *DropOpenflowRuntimeRequest {
	s := DropOpenflowRuntimeRequest{}
	s.name = name
	return &s
}

func (s *DropOpenflowRuntimeRequest) WithIfExists(ifExists bool) *DropOpenflowRuntimeRequest {
	s.IfExists = &ifExists
	return s
}

func (s *DropOpenflowRuntimeRequest) WithCascade(cascade bool) *DropOpenflowRuntimeRequest {
	s.Cascade = &cascade
	return s
}

func NewShowOpenflowRuntimeRequest() *ShowOpenflowRuntimeRequest {
	return &ShowOpenflowRuntimeRequest{}
}

func (s *ShowOpenflowRuntimeRequest) WithLike(like Like) *ShowOpenflowRuntimeRequest {
	s.Like = &like
	return s
}

func NewDescribeOpenflowRuntimeRequest(name SchemaObjectIdentifier) *DescribeOpenflowRuntimeRequest {
	s := DescribeOpenflowRuntimeRequest{}
	s.name = name
	return &s
}
