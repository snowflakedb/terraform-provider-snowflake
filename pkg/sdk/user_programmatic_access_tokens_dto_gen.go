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
	UserName                             AccountObjectIdentifier // required
	name                                 AccountObjectIdentifier // required
	RoleRestriction                      *AccountObjectIdentifier
	DaysToExpiry                         *int
	MinsToBypassNetworkPolicyRequirement *int
	Comment                              *string
}

type ModifyUserProgrammaticAccessTokenRequest struct {
	IfExists *bool
	UserName AccountObjectIdentifier // required
	name     AccountObjectIdentifier // required
	Set      *ModifyProgrammaticAccessTokenSetRequest
	Unset    *ModifyProgrammaticAccessTokenUnsetRequest
	RenameTo *AccountObjectIdentifier
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
	UserName                     AccountObjectIdentifier // required
	name                         AccountObjectIdentifier // required
	ExpireRotatedTokenAfterHours *int
}

type RemoveUserProgrammaticAccessTokenRequest struct {
	IfExists *bool
	UserName AccountObjectIdentifier // required
	name     AccountObjectIdentifier // required
}

type ShowUserProgrammaticAccessTokenRequest struct {
	UserName *AccountObjectIdentifier
}
