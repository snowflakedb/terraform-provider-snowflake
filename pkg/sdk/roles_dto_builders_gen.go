package sdk

func (s *CreateRoleRequest) GetName() AccountObjectIdentifier {
	return s.name
}

func NewCreateRoleRequest(name AccountObjectIdentifier) *CreateRoleRequest {
	s := CreateRoleRequest{}
	s.name = name
	return &s
}

func (s *CreateRoleRequest) WithOrReplace(orReplace bool) *CreateRoleRequest {
	s.OrReplace = Bool(orReplace)
	return s
}

func (s *CreateRoleRequest) WithIfNotExists(ifNotExists bool) *CreateRoleRequest {
	s.IfNotExists = Bool(ifNotExists)
	return s
}

func (s *CreateRoleRequest) WithComment(comment string) *CreateRoleRequest {
	s.Comment = String(comment)
	return s
}

func (s *CreateRoleRequest) WithTag(tag []TagAssociation) *CreateRoleRequest {
	s.Tag = tag
	return s
}

func NewAlterRoleRequest(name AccountObjectIdentifier) *AlterRoleRequest {
	s := AlterRoleRequest{}
	s.name = name
	return &s
}

func (s *AlterRoleRequest) WithIfExists(ifExists bool) *AlterRoleRequest {
	s.IfExists = Bool(ifExists)
	return s
}

func (s *AlterRoleRequest) WithRenameTo(renameTo AccountObjectIdentifier) *AlterRoleRequest {
	s.RenameTo = &renameTo
	return s
}

func (s *AlterRoleRequest) WithSetComment(setComment string) *AlterRoleRequest {
	s.SetComment = String(setComment)
	return s
}

func (s *AlterRoleRequest) WithSetTags(setTags []TagAssociation) *AlterRoleRequest {
	s.SetTags = setTags
	return s
}

func (s *AlterRoleRequest) WithUnsetComment(unsetComment bool) *AlterRoleRequest {
	s.UnsetComment = Bool(unsetComment)
	return s
}

func (s *AlterRoleRequest) WithUnsetTags(unsetTags []ObjectIdentifier) *AlterRoleRequest {
	s.UnsetTags = unsetTags
	return s
}

func NewDropRoleRequest(name AccountObjectIdentifier) *DropRoleRequest {
	s := DropRoleRequest{}
	s.name = name
	return &s
}

func (s *DropRoleRequest) WithIfExists(ifExists bool) *DropRoleRequest {
	s.IfExists = Bool(ifExists)
	return s
}

func NewShowRoleRequest() *ShowRoleRequest {
	return &ShowRoleRequest{}
}

// LikeRequest is a temporary helper used by ShowByID; removed after generator migration.
type LikeRequest struct {
	pattern string // required
}

func NewLikeRequest(pattern string) *LikeRequest {
	return &LikeRequest{
		pattern: pattern,
	}
}

func (s *ShowRoleRequest) WithLike(like *LikeRequest) *ShowRoleRequest {
	s.Like = &Like{
		Pattern: String(like.pattern),
	}
	return s
}

func (s *ShowRoleRequest) WithInClass(inClass RolesInClass) *ShowRoleRequest {
	s.InClass = &inClass
	return s
}

func NewGrantRoleRequest(name AccountObjectIdentifier, grant GrantRole) *GrantRoleRequest {
	s := GrantRoleRequest{}
	s.name = name
	s.Grant = grant
	return &s
}

func NewRevokeRoleRequest(name AccountObjectIdentifier, revoke RevokeRole) *RevokeRoleRequest {
	s := RevokeRoleRequest{}
	s.name = name
	s.Revoke = revoke
	return &s
}
