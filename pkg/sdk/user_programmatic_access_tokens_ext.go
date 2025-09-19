package sdk

func (r *AddProgrammaticAccessTokenResult) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(r.TokenName)
}

func (r *ProgrammaticAccessToken) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(r.Name)
}
