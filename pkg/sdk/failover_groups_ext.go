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
	opts := &showFailoverGroupDatabasesOptions{
		in: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []struct {
		Name string `db:"name"`
	}{}
	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]AccountObjectIdentifier, len(dest))
	for i, row := range dest {
		resultList[i] = NewAccountObjectIdentifier(row.Name)
	}
	return resultList, nil
}

func (v *failoverGroups) ShowShares(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error) {
	opts := &showFailoverGroupSharesOptions{
		in: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	dest := []struct {
		Name         string `db:"name"`
		OwnerAccount string `db:"owner_account"`
	}{}
	err = v.client.query(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}
	resultList := make([]AccountObjectIdentifier, len(dest))
	for i, r := range dest {
		// TODO [SNOW-1348343]: change during failover groups rework; this was not working correctly with identifiers containing `.` character
		resultList[i] = NewExternalObjectIdentifier(NewAccountIdentifierFromFullyQualifiedName(r.OwnerAccount), NewAccountObjectIdentifier(r.Name)).objectIdentifier.(AccountObjectIdentifier)
	}
	return resultList, nil
}
