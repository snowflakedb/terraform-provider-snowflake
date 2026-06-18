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

type AlterRoleRequest struct {
	IfExists     *bool
	name         AccountObjectIdentifier // required
	RenameTo     *AccountObjectIdentifier
	SetComment   *string
	SetTags      []TagAssociation
	UnsetComment *bool
	UnsetTags    []ObjectIdentifier
}

type DropRoleRequest struct {
	IfExists *bool
	name     AccountObjectIdentifier // required
}

type ShowRoleRequest struct {
	Like    *Like
	InClass *RolesInClass
}

type GrantRoleRequest struct {
	name  AccountObjectIdentifier // required
	Grant GrantRole               // required
}

type RevokeRoleRequest struct {
	name   AccountObjectIdentifier // required
	Revoke RevokeRole              // required
}
