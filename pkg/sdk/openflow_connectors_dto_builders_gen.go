package sdk

func NewCreateOpenflowConnectorRequest(
	name SchemaObjectIdentifier,
	inRuntime SchemaObjectIdentifier,
) *CreateOpenflowConnectorRequest {
	s := CreateOpenflowConnectorRequest{}
	s.name = name
	s.InRuntime = inRuntime
	return &s
}

func (s *CreateOpenflowConnectorRequest) WithIfNotExists(ifNotExists bool) *CreateOpenflowConnectorRequest {
	s.IfNotExists = &ifNotExists
	return s
}

func (s *CreateOpenflowConnectorRequest) WithFromDefinition(fromDefinition string) *CreateOpenflowConnectorRequest {
	s.FromDefinition = &fromDefinition
	return s
}

func (s *CreateOpenflowConnectorRequest) WithFromStage(fromStage string) *CreateOpenflowConnectorRequest {
	s.FromStage = &fromStage
	return s
}

func (s *CreateOpenflowConnectorRequest) WithDisplayName(displayName string) *CreateOpenflowConnectorRequest {
	s.DisplayName = &displayName
	return s
}

func (s *CreateOpenflowConnectorRequest) WithComment(comment string) *CreateOpenflowConnectorRequest {
	s.Comment = &comment
	return s
}

func NewAlterOpenflowConnectorRequest(name SchemaObjectIdentifier) *AlterOpenflowConnectorRequest {
	s := AlterOpenflowConnectorRequest{}
	s.name = name
	return &s
}

func (s *AlterOpenflowConnectorRequest) WithStart(start bool) *AlterOpenflowConnectorRequest {
	s.Start = &start
	return s
}

func (s *AlterOpenflowConnectorRequest) WithStop(stop bool) *AlterOpenflowConnectorRequest {
	s.Stop = &stop
	return s
}

func (s *AlterOpenflowConnectorRequest) WithSet(set OpenflowConnectorSetRequest) *AlterOpenflowConnectorRequest {
	s.Set = &set
	return s
}

func (s *AlterOpenflowConnectorRequest) WithUnset(unset OpenflowConnectorUnsetRequest) *AlterOpenflowConnectorRequest {
	s.Unset = &unset
	return s
}

func NewOpenflowConnectorSetRequest() *OpenflowConnectorSetRequest {
	return &OpenflowConnectorSetRequest{}
}

func (s *OpenflowConnectorSetRequest) WithDisplayName(displayName string) *OpenflowConnectorSetRequest {
	s.DisplayName = &displayName
	return s
}

func (s *OpenflowConnectorSetRequest) WithComment(comment string) *OpenflowConnectorSetRequest {
	s.Comment = &comment
	return s
}

func NewOpenflowConnectorUnsetRequest() *OpenflowConnectorUnsetRequest {
	return &OpenflowConnectorUnsetRequest{}
}

func (s *OpenflowConnectorUnsetRequest) WithDisplayName(displayName bool) *OpenflowConnectorUnsetRequest {
	s.DisplayName = &displayName
	return s
}

func (s *OpenflowConnectorUnsetRequest) WithComment(comment bool) *OpenflowConnectorUnsetRequest {
	s.Comment = &comment
	return s
}

func NewDropOpenflowConnectorRequest(name SchemaObjectIdentifier) *DropOpenflowConnectorRequest {
	s := DropOpenflowConnectorRequest{}
	s.name = name
	return &s
}

func (s *DropOpenflowConnectorRequest) WithIfExists(ifExists bool) *DropOpenflowConnectorRequest {
	s.IfExists = &ifExists
	return s
}

func NewShowOpenflowConnectorRequest() *ShowOpenflowConnectorRequest {
	return &ShowOpenflowConnectorRequest{}
}

func (s *ShowOpenflowConnectorRequest) WithLike(like Like) *ShowOpenflowConnectorRequest {
	s.Like = &like
	return s
}

func NewDescribeOpenflowConnectorRequest(name SchemaObjectIdentifier) *DescribeOpenflowConnectorRequest {
	s := DescribeOpenflowConnectorRequest{}
	s.name = name
	return &s
}
