package sdk

func (v *ReplicationAccount) ID() AccountIdentifier {
	return AccountIdentifier{
		organizationName: v.OrganizationName,
		accountName:      v.AccountName,
		accountLocator:   v.AccountLocator,
	}
}
