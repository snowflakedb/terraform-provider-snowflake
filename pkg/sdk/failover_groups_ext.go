package sdk

import "context"

func (v *FailoverGroup) ExternalID() ExternalObjectIdentifier {
	return NewExternalObjectIdentifier(AccountIdentifier{
		organizationName: v.OrganizationName,
		accountName:      v.AccountName,
		accountLocator:   v.AccountLocator,
	}, v.ID())
}

func (v *failoverGroups) ShowDatabases(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error) {
	databases, err := v.ShowFailoverGroupDatabases(ctx, NewShowFailoverGroupDatabasesRequest(id))
	if err != nil {
		return nil, err
	}
	result := make([]AccountObjectIdentifier, len(databases))
	for i, db := range databases {
		result[i] = NewAccountObjectIdentifier(db.Name)
	}
	return result, nil
}

func (v *failoverGroups) ShowShares(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error) {
	shares, err := v.ShowFailoverGroupShares(ctx, NewShowFailoverGroupSharesRequest(id))
	if err != nil {
		return nil, err
	}
	result := make([]AccountObjectIdentifier, len(shares))
	for i, s := range shares {
		// TODO [SNOW-1348343]: change during failover groups rework; this was not working correctly with identifiers containing `.` character
		result[i] = NewExternalObjectIdentifier(NewAccountIdentifierFromFullyQualifiedName(s.OwnerAccount), NewAccountObjectIdentifier(s.Name)).objectIdentifier.(AccountObjectIdentifier)
	}
	return result, nil
}
