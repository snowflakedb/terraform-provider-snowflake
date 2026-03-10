package sdk

func NewCreateExternalAccessIntegrationRequest(
	name AccountObjectIdentifier,
	allowedNetworkRules []SchemaObjectIdentifier,
	enabled bool,
) *CreateExternalAccessIntegrationRequest {
	s := CreateExternalAccessIntegrationRequest{}
	s.name = name
	s.AllowedNetworkRules = allowedNetworkRules
	s.Enabled = enabled
	return &s
}

func (s *CreateExternalAccessIntegrationRequest) WithOrReplace(orReplace bool) *CreateExternalAccessIntegrationRequest {
	s.OrReplace = &orReplace
	return s
}

func (s *CreateExternalAccessIntegrationRequest) WithIfNotExists(ifNotExists bool) *CreateExternalAccessIntegrationRequest {
	s.IfNotExists = &ifNotExists
	return s
}

func (s *CreateExternalAccessIntegrationRequest) WithAllowedAuthenticationSecrets(secrets []SchemaObjectIdentifier) *CreateExternalAccessIntegrationRequest {
	s.AllowedAuthenticationSecrets = secrets
	return s
}

func (s *CreateExternalAccessIntegrationRequest) WithComment(comment string) *CreateExternalAccessIntegrationRequest {
	s.Comment = &comment
	return s
}

func NewAlterExternalAccessIntegrationRequest(
	name AccountObjectIdentifier,
) *AlterExternalAccessIntegrationRequest {
	s := AlterExternalAccessIntegrationRequest{}
	s.name = name
	return &s
}

func (s *AlterExternalAccessIntegrationRequest) WithIfExists(ifExists bool) *AlterExternalAccessIntegrationRequest {
	s.IfExists = &ifExists
	return s
}

func (s *AlterExternalAccessIntegrationRequest) WithSet(set ExternalAccessIntegrationSetRequest) *AlterExternalAccessIntegrationRequest {
	s.Set = &set
	return s
}

func (s *AlterExternalAccessIntegrationRequest) WithUnset(unset ExternalAccessIntegrationUnsetRequest) *AlterExternalAccessIntegrationRequest {
	s.Unset = &unset
	return s
}

func NewExternalAccessIntegrationSetRequest() *ExternalAccessIntegrationSetRequest {
	s := ExternalAccessIntegrationSetRequest{}
	return &s
}

func (s *ExternalAccessIntegrationSetRequest) WithAllowedNetworkRules(rules []SchemaObjectIdentifier) *ExternalAccessIntegrationSetRequest {
	s.AllowedNetworkRules = rules
	return s
}

func (s *ExternalAccessIntegrationSetRequest) WithAllowedAuthenticationSecrets(secrets []SchemaObjectIdentifier) *ExternalAccessIntegrationSetRequest {
	s.AllowedAuthenticationSecrets = secrets
	return s
}

func (s *ExternalAccessIntegrationSetRequest) WithEnabled(enabled bool) *ExternalAccessIntegrationSetRequest {
	s.Enabled = &enabled
	return s
}

func (s *ExternalAccessIntegrationSetRequest) WithComment(comment string) *ExternalAccessIntegrationSetRequest {
	s.Comment = &comment
	return s
}

func NewExternalAccessIntegrationUnsetRequest() *ExternalAccessIntegrationUnsetRequest {
	s := ExternalAccessIntegrationUnsetRequest{}
	return &s
}

func (s *ExternalAccessIntegrationUnsetRequest) WithAllowedAuthenticationSecrets(v bool) *ExternalAccessIntegrationUnsetRequest {
	s.AllowedAuthenticationSecrets = &v
	return s
}

func (s *ExternalAccessIntegrationUnsetRequest) WithComment(comment bool) *ExternalAccessIntegrationUnsetRequest {
	s.Comment = &comment
	return s
}

func NewDropExternalAccessIntegrationRequest(
	name AccountObjectIdentifier,
) *DropExternalAccessIntegrationRequest {
	s := DropExternalAccessIntegrationRequest{}
	s.name = name
	return &s
}

func (s *DropExternalAccessIntegrationRequest) WithIfExists(ifExists bool) *DropExternalAccessIntegrationRequest {
	s.IfExists = &ifExists
	return s
}

func NewShowExternalAccessIntegrationRequest() *ShowExternalAccessIntegrationRequest {
	s := ShowExternalAccessIntegrationRequest{}
	return &s
}

func (s *ShowExternalAccessIntegrationRequest) WithLike(like Like) *ShowExternalAccessIntegrationRequest {
	s.Like = &like
	return s
}

func NewDescribeExternalAccessIntegrationRequest(
	name AccountObjectIdentifier,
) *DescribeExternalAccessIntegrationRequest {
	s := DescribeExternalAccessIntegrationRequest{}
	s.name = name
	return &s
}
