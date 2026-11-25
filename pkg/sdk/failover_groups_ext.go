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

// TODO [SNOW-4490849]: a value of one field that requires another field (here ALLOWED_INTEGRATION_TYPES requires INTEGRATIONS in OBJECT_TYPES)
// has no validation in the SSoT definition yet.
func (v *FailoverGroupSet) additionalValidations() error {
	if len(v.AllowedIntegrationTypes) > 0 {
		if !slices.Contains(v.ObjectTypes, FailoverGroupObjectTypeIntegrations) {
			return errors.New("INTEGRATIONS must be set in OBJECT_TYPES when setting allowed integration types")
		}
	}
	return nil
}

// REMOVE is validated by the generated ChildrenRule. ADD and MOVE are not, see the TODOs in the definition.
func (opts *AlterFailoverGroupOptions) additionalValidations() error {
	var errs []error
	if opts.Add != nil {
		// TODO [SNOW-4490849]: exactly one list, and IGNORE EDITION CHECK only with the accounts, in the definition of ADD.
		if !exactlyOneValueSet(opts.Add.AllowedDatabases, opts.Add.AllowedExternalVolumes, opts.Add.AllowedShares, opts.Add.AllowedAccounts) {
			errs = append(errs, NewError("ADD must set exactly one allowed list"))
		}
		if opts.Add.IgnoreEditionCheck != nil && *opts.Add.IgnoreEditionCheck && !valueSet(opts.Add.AllowedAccounts) {
			errs = append(errs, NewError("IGNORE EDITION CHECK requires ADD TO ALLOWED_ACCOUNTS"))
		}
	}
	if opts.Move != nil {
		// Snowflake accepts several lists in one MOVE, although the page shows one list per statement.
		// TODO [SNOW-4490849]: at least one list, in the definition of MOVE.
		if !anyValueSet(opts.Move.Databases, opts.Move.ExternalVolumes, opts.Move.Shares) {
			errs = append(errs, NewError("MOVE must set at least one of DATABASES, EXTERNAL VOLUMES, and SHARES"))
		}
	}
	return JoinErrors(errs...)
}

// TODO [SNOW-4490849]: each OBJECT_TYPES value that goes together with an ALLOWED_* list, and SECURITY INTEGRATIONS requiring ROLES,
// are dependencies between fields that have no validation in the SSoT definition yet.
func (opts *CreateFailoverGroupOptions) additionalValidations() error {
	var errs []error
	errs = append(errs, failoverGroupRequireObjectList(
		opts.ObjectTypes, FailoverGroupObjectTypeDatabases, opts.AllowedDatabases,
		"DATABASES", "ALLOWED_DATABASES",
	)...)
	errs = append(errs, failoverGroupRequireObjectList(
		opts.ObjectTypes, FailoverGroupObjectTypeExternalVolumes, opts.AllowedExternalVolumes,
		"EXTERNAL VOLUMES", "ALLOWED_EXTERNAL_VOLUMES",
	)...)
	errs = append(errs, failoverGroupRequireObjectList(
		opts.ObjectTypes, FailoverGroupObjectTypeShares, opts.AllowedShares,
		"SHARES", "ALLOWED_SHARES",
	)...)
	errs = append(errs, failoverGroupRequireObjectList(
		opts.ObjectTypes, FailoverGroupObjectTypeIntegrations, opts.AllowedIntegrationTypes,
		"INTEGRATIONS", "ALLOWED_INTEGRATION_TYPES",
	)...)
	if slices.Contains(opts.AllowedIntegrationTypes, FailoverGroupIntegrationTypeSecurityIntegrations) &&
		!slices.Contains(opts.ObjectTypes, FailoverGroupObjectTypeRoles) {
		errs = append(errs, NewError("OBJECT_TYPES must include ROLES when ALLOWED_INTEGRATION_TYPES includes SECURITY INTEGRATIONS"))
	}
	return JoinErrors(errs...)
}

func failoverGroupRequireObjectList[T any](objectTypes []FailoverGroupObjectType, objectType FailoverGroupObjectType, list []T, objectTypeSQL, listSQL string) []error {
	hasType := slices.Contains(objectTypes, objectType)
	hasList := len(list) > 0
	switch {
	case hasType && !hasList:
		return []error{NewError(listSQL + " must be set when OBJECT_TYPES includes " + objectTypeSQL)}
	case hasList && !hasType:
		return []error{NewError("OBJECT_TYPES must include " + objectTypeSQL + " when " + listSQL + " is set")}
	default:
		return nil
	}
}

// TODO [SNOW-4490849]: the identifier validation of an ExternalObjectIdentifier field (it can hold a nil identifier) is not generated.
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
		result[i] = NewAccountObjectIdentifier(s.Name)
	}
	return result, nil
}

// The SHOW converters below keep the leniency of the handwritten converter: a value they do not know is not an error,
// because a SHOW row that cannot be read breaks the read of every failover group (see GH #1418).

// SHOW prints PARAMETERS for the ACCOUNT PARAMETERS object type. A token that is not a known object type is kept in upper case.
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
		objectType, err := ToFailoverGroupObjectType(part)
		if err != nil {
			objectType = FailoverGroupObjectType(strings.ToUpper(part))
		}
		result = append(result, objectType)
	}
	return result, nil
}

// SHOW prints short tokens such as API and EXTERNAL_ACCESS. They become the create-page values (API INTEGRATIONS, EXTERNAL ACCESS INTEGRATIONS).
// A token that is not a known integration type is kept in upper case.
func ParseCommaSeparatedFailoverGroupIntegrationTypeArray(value string) ([]FailoverGroupIntegrationType, error) {
	var result []FailoverGroupIntegrationType
	for part := range strings.SplitSeq(value, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		part = strings.ReplaceAll(part, "_", " ") + " INTEGRATIONS"
		integrationType, err := ToFailoverGroupIntegrationType(part)
		if err != nil {
			integrationType = FailoverGroupIntegrationType(strings.ToUpper(part))
		}
		result = append(result, integrationType)
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
