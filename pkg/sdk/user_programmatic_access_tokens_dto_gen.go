package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[AddUserProgrammaticAccessTokenOptions]    = new(AddUserProgrammaticAccessTokenRequest)
	_ optionsProvider[ModifyUserProgrammaticAccessTokenOptions] = new(ModifyUserProgrammaticAccessTokenRequest)
	_ optionsProvider[RotateUserProgrammaticAccessTokenOptions] = new(RotateUserProgrammaticAccessTokenRequest)
	_ optionsProvider[RemoveUserProgrammaticAccessTokenOptions] = new(RemoveUserProgrammaticAccessTokenRequest)
	_ optionsProvider[ShowUserProgrammaticAccessTokenOptions]   = new(ShowUserProgrammaticAccessTokenRequest)
)

type AddUserProgrammaticAccessTokenRequest struct {
	IfExists                             *bool
	User                                 *AccountObjectIdentifier
	name                                 string // required
	RoleRestriction                      *AccountObjectIdentifier
	DaysToExpiry                         *int
	MinsToBypassNetworkPolicyRequirement *int
	Comment                              *string
}

type ModifyUserProgrammaticAccessTokenRequest struct {
	IfExists *bool
	User     *AccountObjectIdentifier
	name     string // required
	Set      *ModifyProgrammaticAccessTokenSetRequest
	Unset    *ModifyProgrammaticAccessTokenUnsetRequest
	RenameTo *string
}

type ModifyProgrammaticAccessTokenSetRequest struct {
	Disabled                             *bool
	MinsToBypassNetworkPolicyRequirement *int
	Comment                              *string
}

type ModifyProgrammaticAccessTokenUnsetRequest struct {
	Disabled                             *bool
	MinsToBypassNetworkPolicyRequirement *bool
	Comment                              *bool
}

type RotateUserProgrammaticAccessTokenRequest struct {
	IfExists                     *bool
	User                         *AccountObjectIdentifier
	name                         string // required
	ExpireRotatedTokenAfterHours *int
}

type RemoveUserProgrammaticAccessTokenRequest struct {
	IfExists *bool
	User     *AccountObjectIdentifier
	name     string // required
}

type ShowUserProgrammaticAccessTokenRequest struct {
	User *AccountObjectIdentifier
}
