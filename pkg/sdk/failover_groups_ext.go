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
		if !slices.Contains(v.ObjectTypes, FailoverGroupObjectTypeIntegrations) {
			return errors.New("INTEGRATIONS must be set in OBJECT_TYPES when setting allowed integration types")
		}
	}
	return nil
}

func (opts *CreateSecondaryFailoverGroupOptions) additionalValidations() error {
	if opts.AsReplicaOf.objectIdentifier == nil || !ValidObjectIdentifier(opts.AsReplicaOf) {
		return errInvalidIdentifier("CreateSecondaryFailoverGroupOptions", "AsReplicaOf")
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

// The SHOW converters below keep the leniency of the handwritten converter: a value they do not know is not an error,
// because a SHOW row that cannot be read breaks the read of every failover group (see GH #1418).

// SHOW prints PARAMETERS for the ACCOUNT PARAMETERS object type. Any other token is left unchanged.
func ParseCommaSeparatedFailoverGroupObjectTypeArray(value string) ([]FailoverGroupObjectType, error) {
	var result []FailoverGroupObjectType
	for part := range strings.SplitSeq(value, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.EqualFold(part, "PARAMETERS") {
			part = string(FailoverGroupObjectTypeAccountParameters)
		}
		result = append(result, FailoverGroupObjectType(part))
	}
	return result, nil
}

// SHOW prints short tokens such as API and EXTERNAL_ACCESS. They become the create-page values (API INTEGRATIONS, EXTERNAL ACCESS INTEGRATIONS).
func ParseCommaSeparatedFailoverGroupIntegrationTypeArray(value string) ([]FailoverGroupIntegrationType, error) {
	var result []FailoverGroupIntegrationType
	for part := range strings.SplitSeq(value, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		result = append(result, FailoverGroupIntegrationType(strings.ReplaceAll(part, "_", " ")+" INTEGRATIONS"))
	}
	return result, nil
}

// SHOW prints org.account. A value that has not exactly two parts is skipped.
func ParseCommaSeparatedFailoverGroupAllowedAccountArray(value string) ([]AccountIdentifier, error) {
	var result []AccountIdentifier
	for part := range strings.SplitSeq(value, ",") {
		parts := strings.Split(strings.TrimSpace(part), ".")
		if len(parts) != 2 {
			continue
		}
		result = append(result, NewAccountIdentifier(parts[0], parts[1]))
	}
	return result, nil
}
