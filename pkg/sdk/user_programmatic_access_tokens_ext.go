package sdk

func (r *AddProgrammaticAccessTokenResult) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(r.TokenName)
}
