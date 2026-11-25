package sdk

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// TODO [SNOW-4490849]: each OBJECT_TYPES value that goes together with an ALLOWED_* list, and SECURITY INTEGRATIONS requiring ROLES,
// are dependencies between fields that have no validation in the SSoT definition yet.
func (opts *CreateReplicationGroupOptions) additionalValidations() error {
	var errs []error
	errs = append(errs, replicationGroupRequireObjectList(
		opts.ObjectTypes, ReplicationGroupObjectTypeDatabases, opts.AllowedDatabases,
		"DATABASES", "ALLOWED_DATABASES",
	)...)
	errs = append(errs, replicationGroupRequireObjectList(
		opts.ObjectTypes, ReplicationGroupObjectTypeExternalVolumes, opts.AllowedExternalVolumes,
		"EXTERNAL VOLUMES", "ALLOWED_EXTERNAL_VOLUMES",
	)...)
	errs = append(errs, replicationGroupRequireObjectList(
		opts.ObjectTypes, ReplicationGroupObjectTypeShares, opts.AllowedShares,
		"SHARES", "ALLOWED_SHARES",
	)...)
	errs = append(errs, replicationGroupRequireObjectList(
		opts.ObjectTypes, ReplicationGroupObjectTypeIntegrations, opts.AllowedIntegrationTypes,
		"INTEGRATIONS", "ALLOWED_INTEGRATION_TYPES",
	)...)
	if slices.Contains(opts.AllowedIntegrationTypes, ReplicationGroupIntegrationTypeSecurityIntegrations) &&
		!slices.Contains(opts.ObjectTypes, ReplicationGroupObjectTypeRoles) {
		errs = append(errs, NewError("OBJECT_TYPES must include ROLES when ALLOWED_INTEGRATION_TYPES includes SECURITY INTEGRATIONS"))
	}
	return JoinErrors(errs...)
}

func replicationGroupRequireObjectList[T any](objectTypes []ReplicationGroupObjectType, objectType ReplicationGroupObjectType, list []T, objectTypeSQL, listSQL string) []error {
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

func (r *CreateReplicationGroupRequest) GetName() AccountObjectIdentifier {
	return r.name
}

// REMOVE is validated by the generated ChildrenRule. ADD and MOVE are not, see the TODOs in the definition.
func (opts *AlterReplicationGroupOptions) additionalValidations() error {
	var errs []error
	if opts.Move != nil {
		// TODO [SNOW-4490849]: exactly one list, in the definition of MOVE.
		if !exactlyOneValueSet(opts.Move.Databases, opts.Move.ExternalVolumes, opts.Move.Shares) {
			errs = append(errs, NewError("MOVE must set exactly one of DATABASES, EXTERNAL VOLUMES, and SHARES"))
		}
		if !ValidObjectIdentifier(opts.Move.ToReplicationGroup) {
			errs = append(errs, errInvalidIdentifier("AlterReplicationGroupOptions.Move", "ToReplicationGroup"))
		}
	}
	if opts.Add != nil {
		// TODO [SNOW-4490849]: exactly one list, and IGNORE EDITION CHECK only with the accounts, in the definition of ADD.
		if !exactlyOneValueSet(opts.Add.ToAllowedDatabases, opts.Add.ToAllowedExternalVolumes, opts.Add.ToAllowedShares, opts.Add.ToAllowedAccounts) {
			errs = append(errs, NewError("ADD must set exactly one allowed list"))
		}
		if opts.Add.IgnoreEditionCheck != nil && *opts.Add.IgnoreEditionCheck && !valueSet(opts.Add.ToAllowedAccounts) {
			errs = append(errs, NewError("IGNORE EDITION CHECK requires ADD TO ALLOWED_ACCOUNTS"))
		}
	}
	return JoinErrors(errs...)
}

// TODO [SNOW-4490849]: the identifier validation of an ExternalObjectIdentifier field (it can hold a nil identifier) is not generated.
func (opts *CreateSecondaryReplicationGroupOptions) additionalValidations() error {
	if opts.AsReplicaOf.objectIdentifier == nil || !ValidObjectIdentifier(opts.AsReplicaOf) {
		return errInvalidIdentifier("CreateSecondaryReplicationGroupOptions", "AsReplicaOf")
	}
	return nil
}

// TODO [SNOW-4490849]: a value of one field that requires another field (an ALLOWED_* list requires its OBJECT_TYPES value) has no validation
// in the SSoT definition yet.
func (opts *ReplicationGroupSet) additionalValidations() error {
	var errs []error
	if len(opts.AllowedDatabases) > 0 && !slices.Contains(opts.ObjectTypes, ReplicationGroupObjectTypeDatabases) {
		errs = append(errs, NewError("OBJECT_TYPES must include DATABASES when ALLOWED_DATABASES is set"))
	}
	if len(opts.AllowedExternalVolumes) > 0 && !slices.Contains(opts.ObjectTypes, ReplicationGroupObjectTypeExternalVolumes) {
		errs = append(errs, NewError("OBJECT_TYPES must include EXTERNAL VOLUMES when ALLOWED_EXTERNAL_VOLUMES is set"))
	}
	if len(opts.AllowedShares) > 0 && !slices.Contains(opts.ObjectTypes, ReplicationGroupObjectTypeShares) {
		errs = append(errs, NewError("OBJECT_TYPES must include SHARES when ALLOWED_SHARES is set"))
	}
	if len(opts.AllowedIntegrationTypes) > 0 && !slices.Contains(opts.ObjectTypes, ReplicationGroupObjectTypeIntegrations) {
		errs = append(errs, NewError("OBJECT_TYPES must include INTEGRATIONS when ALLOWED_INTEGRATION_TYPES is set"))
	}
	if slices.Contains(opts.AllowedIntegrationTypes, ReplicationGroupIntegrationTypeSecurityIntegrations) &&
		!slices.Contains(opts.ObjectTypes, ReplicationGroupObjectTypeRoles) {
		errs = append(errs, NewError("OBJECT_TYPES must include ROLES when ALLOWED_INTEGRATION_TYPES includes SECURITY INTEGRATIONS"))
	}
	return JoinErrors(errs...)
}

// SHOW tokens that differ from the create-page values. A token that is not listed is left unchanged.
var replicationGroupObjectTypeShowTokens = map[string]string{
	"PARAMETERS": string(ReplicationGroupObjectTypeAccountParameters),
}

var replicationGroupIntegrationTypeShowTokens = map[string]string{
	"API":             string(ReplicationGroupIntegrationTypeApiIntegrations),
	"EXTERNAL_ACCESS": string(ReplicationGroupIntegrationTypeExternalAccessIntegrations),
	"NOTIFICATION":    string(ReplicationGroupIntegrationTypeNotificationIntegrations),
	"SECURITY":        string(ReplicationGroupIntegrationTypeSecurityIntegrations),
	"STORAGE":         string(ReplicationGroupIntegrationTypeStorageIntegrations),
}

func replicationGroupShowToken(part string, tokens map[string]string) string {
	if mapped, ok := tokens[strings.ToUpper(part)]; ok {
		return mapped
	}
	return part
}

func ParseCommaSeparatedReplicationGroupObjectTypeArray(value string) ([]ReplicationGroupObjectType, error) {
	parts := ParseCommaSeparatedStringArray(value, false)
	rewritten := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		rewritten = append(rewritten, replicationGroupShowToken(part, replicationGroupObjectTypeShowTokens))
	}
	return collections.MapErr(rewritten, ToReplicationGroupObjectType)
}

func ParseCommaSeparatedReplicationGroupIntegrationTypeArray(value string) ([]ReplicationGroupIntegrationType, error) {
	parts := ParseCommaSeparatedStringArray(value, false)
	rewritten := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		rewritten = append(rewritten, replicationGroupShowToken(part, replicationGroupIntegrationTypeShowTokens))
	}
	return collections.MapErr(rewritten, ToReplicationGroupIntegrationType)
}
