package sdk

var (
	_ optionsProvider[CreateRoleOptions] = new(CreateRoleRequest)
	_ optionsProvider[AlterRoleOptions]  = new(AlterRoleRequest)
	_ optionsProvider[DropRoleOptions]   = new(DropRoleRequest)
	_ optionsProvider[ShowRoleOptions]   = new(ShowRoleRequest)
	_ optionsProvider[GrantRoleOptions]  = new(GrantRoleRequest)
	_ optionsProvider[RevokeRoleOptions] = new(RevokeRoleRequest)
)

type CreateRoleRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        AccountObjectIdentifier // required
	Comment     *string
	Tag         []TagAssociation
}

func (s *CreateRoleRequest) GetName() AccountObjectIdentifier {
	return s.name
}

func NewCreateRoleRequest(name AccountObjectIdentifier) *CreateRoleRequest {
	return &CreateRoleRequest{
		name: name,
	}
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

func (s *CreateRoleRequest) toOpts() *CreateRoleOptions {
	return &CreateRoleOptions{
		OrReplace:   s.OrReplace,
		IfNotExists: s.IfNotExists,
		name:        s.name,
		Comment:     s.Comment,
		Tag:         s.Tag,
	}
}

type AlterRoleRequest struct {
	IfExists     *bool
	name         AccountObjectIdentifier // required
	RenameTo     *AccountObjectIdentifier
	SetComment   *string
	SetTags      []TagAssociation
	UnsetComment *bool
	UnsetTags    []ObjectIdentifier
}

func NewAlterRoleRequest(name AccountObjectIdentifier) *AlterRoleRequest {
	return &AlterRoleRequest{
		name: name,
	}
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

func (s *AlterRoleRequest) toOpts() *AlterRoleOptions {
	return &AlterRoleOptions{
		IfExists:     s.IfExists,
		name:         s.name,
		RenameTo:     s.RenameTo,
		SetComment:   s.SetComment,
		SetTags:      s.SetTags,
		UnsetComment: s.UnsetComment,
		UnsetTags:    s.UnsetTags,
	}
}

type DropRoleRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

func NewDropRoleRequest(name AccountObjectIdentifier) *DropRoleRequest {
	return &DropRoleRequest{
		name: name,
	}
}

func (s *DropRoleRequest) WithIfExists(ifExists bool) *DropRoleRequest {
	s.IfExists = Bool(ifExists)
	return s
}

func (s *DropRoleRequest) toOpts() *DropRoleOptions {
	return &DropRoleOptions{
		IfExists: s.IfExists,
		name:     s.name,
	}
}

type ShowRoleRequest struct {
	Like    *Like
	InClass *RolesInClass
}

func NewShowRoleRequest() *ShowRoleRequest {
	return &ShowRoleRequest{}
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

type LikeRequest struct {
	pattern string // required
}

func NewLikeRequest(pattern string) *LikeRequest {
	return &LikeRequest{
		pattern: pattern,
	}
}

func (s *ShowRoleRequest) toOpts() *ShowRoleOptions {
	return &ShowRoleOptions{
		Like:    s.Like,
		InClass: s.InClass,
	}
}

type GrantRoleRequest struct {
	name  AccountObjectIdentifier // required
	Grant GrantRole               // required
}

func NewGrantRoleRequest(name AccountObjectIdentifier, grant GrantRole) *GrantRoleRequest {
	return &GrantRoleRequest{
		name:  name,
		Grant: grant,
	}
}

func (s *GrantRoleRequest) toOpts() *GrantRoleOptions {
	return &GrantRoleOptions{
		name:  s.name,
		Grant: s.Grant,
	}
}

type RevokeRoleRequest struct {
	name   AccountObjectIdentifier // required
	Revoke RevokeRole              // required
}

func NewRevokeRoleRequest(name AccountObjectIdentifier, revoke RevokeRole) *RevokeRoleRequest {
	return &RevokeRoleRequest{
		name:   name,
		Revoke: revoke,
	}
}

func (s *RevokeRoleRequest) toOpts() *RevokeRoleOptions {
	return &RevokeRoleOptions{
		name:   s.name,
		Revoke: s.Revoke,
	}
}
