package sdk

import "context"

func (r *AddProgrammaticAccessTokenResult) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(r.TokenName)
}

func (v *ProgrammaticAccessToken) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *userProgrammaticAccessTokens) RemoveByIDSafely(ctx context.Context, request *RemoveUserProgrammaticAccessTokenRequest) error {
	return SafeRemoveProgrammaticAccessToken(v.client, ctx, request)
}
