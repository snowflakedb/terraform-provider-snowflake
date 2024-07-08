// Code generated by dto builder generator; DO NOT EDIT.

package sdk

import ()

func NewCreateStreamlitRequest(
	name SchemaObjectIdentifier,
	RootLocation string,
	MainFile string,
) *CreateStreamlitRequest {
	s := CreateStreamlitRequest{}
	s.name = name
	s.RootLocation = RootLocation
	s.MainFile = MainFile
	return &s
}

func (s *CreateStreamlitRequest) WithOrReplace(OrReplace bool) *CreateStreamlitRequest {
	s.OrReplace = &OrReplace
	return s
}

func (s *CreateStreamlitRequest) WithIfNotExists(IfNotExists bool) *CreateStreamlitRequest {
	s.IfNotExists = &IfNotExists
	return s
}

func (s *CreateStreamlitRequest) WithWarehouse(Warehouse AccountObjectIdentifier) *CreateStreamlitRequest {
	s.Warehouse = &Warehouse
	return s
}

func (s *CreateStreamlitRequest) WithExternalAccessIntegrations(ExternalAccessIntegrations ExternalAccessIntegrationsRequest) *CreateStreamlitRequest {
	s.ExternalAccessIntegrations = &ExternalAccessIntegrations
	return s
}

func (s *CreateStreamlitRequest) WithTitle(Title string) *CreateStreamlitRequest {
	s.Title = &Title
	return s
}

func (s *CreateStreamlitRequest) WithComment(Comment string) *CreateStreamlitRequest {
	s.Comment = &Comment
	return s
}

func NewExternalAccessIntegrationsRequest(
	ExternalAccessIntegrations []AccountObjectIdentifier,
) *ExternalAccessIntegrationsRequest {
	s := ExternalAccessIntegrationsRequest{}
	s.ExternalAccessIntegrations = ExternalAccessIntegrations
	return &s
}

func NewAlterStreamlitRequest(
	name SchemaObjectIdentifier,
) *AlterStreamlitRequest {
	s := AlterStreamlitRequest{}
	s.name = name
	return &s
}

func (s *AlterStreamlitRequest) WithIfExists(IfExists bool) *AlterStreamlitRequest {
	s.IfExists = &IfExists
	return s
}

func (s *AlterStreamlitRequest) WithSet(Set StreamlitSetRequest) *AlterStreamlitRequest {
	s.Set = &Set
	return s
}

func (s *AlterStreamlitRequest) WithUnset(Unset StreamlitUnsetRequest) *AlterStreamlitRequest {
	s.Unset = &Unset
	return s
}

func (s *AlterStreamlitRequest) WithRenameTo(RenameTo SchemaObjectIdentifier) *AlterStreamlitRequest {
	s.RenameTo = &RenameTo
	return s
}

func NewStreamlitSetRequest(
	RootLocation *string,
	MainFile *string,
) *StreamlitSetRequest {
	s := StreamlitSetRequest{}
	s.RootLocation = RootLocation
	s.MainFile = MainFile
	return &s
}

func (s *StreamlitSetRequest) WithWarehouse(Warehouse AccountObjectIdentifier) *StreamlitSetRequest {
	s.Warehouse = &Warehouse
	return s
}

func (s *StreamlitSetRequest) WithExternalAccessIntegrations(ExternalAccessIntegrations ExternalAccessIntegrationsRequest) *StreamlitSetRequest {
	s.ExternalAccessIntegrations = &ExternalAccessIntegrations
	return s
}

func (s *StreamlitSetRequest) WithComment(Comment string) *StreamlitSetRequest {
	s.Comment = &Comment
	return s
}

func (s *StreamlitSetRequest) WithTitle(Title string) *StreamlitSetRequest {
	s.Title = &Title
	return s
}

func NewStreamlitUnsetRequest() *StreamlitUnsetRequest {
	return &StreamlitUnsetRequest{}
}

func (s *StreamlitUnsetRequest) WithQueryWarehouse(QueryWarehouse bool) *StreamlitUnsetRequest {
	s.QueryWarehouse = &QueryWarehouse
	return s
}

func (s *StreamlitUnsetRequest) WithComment(Comment bool) *StreamlitUnsetRequest {
	s.Comment = &Comment
	return s
}

func (s *StreamlitUnsetRequest) WithTitle(Title bool) *StreamlitUnsetRequest {
	s.Title = &Title
	return s
}

func NewDropStreamlitRequest(
	name SchemaObjectIdentifier,
) *DropStreamlitRequest {
	s := DropStreamlitRequest{}
	s.name = name
	return &s
}

func (s *DropStreamlitRequest) WithIfExists(IfExists bool) *DropStreamlitRequest {
	s.IfExists = &IfExists
	return s
}

func NewShowStreamlitRequest() *ShowStreamlitRequest {
	return &ShowStreamlitRequest{}
}

func (s *ShowStreamlitRequest) WithTerse(Terse bool) *ShowStreamlitRequest {
	s.Terse = &Terse
	return s
}

func (s *ShowStreamlitRequest) WithLike(Like Like) *ShowStreamlitRequest {
	s.Like = &Like
	return s
}

func (s *ShowStreamlitRequest) WithIn(In In) *ShowStreamlitRequest {
	s.In = &In
	return s
}

func (s *ShowStreamlitRequest) WithLimit(Limit LimitFrom) *ShowStreamlitRequest {
	s.Limit = &Limit
	return s
}

func NewDescribeStreamlitRequest(
	name SchemaObjectIdentifier,
) *DescribeStreamlitRequest {
	s := DescribeStreamlitRequest{}
	s.name = name
	return &s
}
