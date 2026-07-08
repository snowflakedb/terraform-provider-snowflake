package sdk

import (
	"context"
	"errors"
	"slices"
	"strings"
)

func (v *FailoverGroup) ExternalID() ExternalObjectIdentifier {
	return NewExternalObjectIdentifier(AccountIdentifier{
		organizationName: v.OrganizationName,
		accountName:      v.AccountName,
		accountLocator:   v.AccountLocator,
	}, v.ID())
}

func (v *FailoverGroupSet) additionalValidations() error {
	if len(v.AllowedIntegrationTypes) > 0 {
		if !slices.Contains(v.ObjectTypes, PluralObjectTypeIntegrations) {
			return errors.New("INTEGRATIONS must be set in OBJECT_TYPES when setting allowed integration types")
		}
	}
	return nil
}

func (row failoverGroupDBRow) additionalConvert(result *FailoverGroup) error {
	result.Primary = NewExternalObjectIdentifierFromFullyQualifiedName(row.Primary)

	ots := strings.Split(row.ObjectTypes, ",")
	result.ObjectTypes = make([]PluralObjectType, 0, len(ots))
	for _, ot := range ots {
		pot := PluralObjectType(strings.TrimSpace(ot))
		if pot == PluralObjectTypeParameters {
			result.ObjectTypes = append(result.ObjectTypes, PluralObjectType("ACCOUNT PARAMETERS"))
		} else {
			result.ObjectTypes = append(result.ObjectTypes, pot)
		}
	}

	its := strings.Split(row.AllowedIntegrationTypes, ",")
	result.AllowedIntegrationTypes = make([]IntegrationType, 0, len(its))
	for _, it := range its {
		if it == "" {
			continue
		}
		result.AllowedIntegrationTypes = append(result.AllowedIntegrationTypes, IntegrationType(strings.ReplaceAll(strings.TrimSpace(it), "_", " ")+" INTEGRATIONS"))
	}

	aas := strings.Split(row.AllowedAccounts, ",")
	result.AllowedAccounts = make([]AccountIdentifier, 0, len(aas))
	for _, aa := range aas {
		s := strings.TrimSpace(aa)
		p := strings.Split(s, ".")
		if len(p) != 2 {
			continue
		}
		result.AllowedAccounts = append(result.AllowedAccounts, NewAccountIdentifier(p[0], p[1]))
	}

	result.SecondaryState = FailoverGroupSecondaryStateNull
	if row.SecondaryState.Valid {
		result.SecondaryState = FailoverGroupSecondaryState(row.SecondaryState.String)
	}

	return nil
}

func (r *CreateFailoverGroupRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (v *failoverGroups) ShowDatabases(ctx context.Context, id AccountObjectIdentifier) ([]AccountObjectIdentifier, error) {
	databases, err := v.ShowFailoverGroupDatabases(ctx, NewShowFailoverGroupDatabasesFailoverGroupRequest(id))
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
	shares, err := v.ShowFailoverGroupShares(ctx, NewShowFailoverGroupSharesFailoverGroupRequest(id))
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
