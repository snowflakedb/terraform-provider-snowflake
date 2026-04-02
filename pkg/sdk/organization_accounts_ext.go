package sdk

func (v *OrganizationAccount) ID() AccountIdentifier {
	return NewAccountIdentifier(v.OrganizationName, v.AccountName)
}
