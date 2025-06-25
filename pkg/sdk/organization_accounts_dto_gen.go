package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateOrganizationAccountOptions] = new(CreateOrganizationAccountRequest)
	_ optionsProvider[AlterOrganizationAccountOptions]  = new(AlterOrganizationAccountRequest)
	_ optionsProvider[ShowOrganizationAccountOptions]   = new(ShowOrganizationAccountRequest)
)

type CreateOrganizationAccountRequest struct {
	name               AccountObjectIdentifier // required
	AdminName          string                  // required
	AdminPassword      *string
	AdminRsaPublicKey  *string
	FirstName          *string
	LastName           *string
	Email              string // required
	MustChangePassword *bool
	Edition            OrganizationAccountEdition // required
	RegionGroup        *string
	Region             *string
	Comment            *string
}

type AlterOrganizationAccountRequest struct {
	Name       *AccountObjectIdentifier
	Set        *OrganizationAccountSetRequest
	Unset      *OrganizationAccountUnsetRequest
	SetTags    []TagAssociation
	UnsetTags  []ObjectIdentifier
	RenameTo   *OrganizationAccountRenameRequest
	DropOldUrl *bool
}

type OrganizationAccountSetRequest struct {
	Parameters      *AccountParameters
	ResourceMonitor *AccountObjectIdentifier
	PasswordPolicy  *AccountObjectIdentifier
	SessionPolicy   *AccountObjectIdentifier
}

type OrganizationAccountUnsetRequest struct {
	Parameters      *AccountParametersUnset
	ResourceMonitor *bool
	PasswordPolicy  *bool
	SessionPolicy   *bool
}

type OrganizationAccountRenameRequest struct {
	NewName    *AccountObjectIdentifier // required
	SaveOldUrl *bool
}

type ShowOrganizationAccountRequest struct {
	Like *Like
}
