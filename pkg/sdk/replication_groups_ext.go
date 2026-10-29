package sdk

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func (opts *CreateReplicationGroupOptions) additionalValidations() error {
	var errs []error
	if len(opts.ObjectTypes) == 0 {
		errs = append(errs, NewError("OBJECT_TYPES must be set"))
	}
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

func (opts *AlterReplicationGroupOptions) additionalValidations() error {
	var errs []error
	if opts.Move != nil {
		if countSetLists(len(opts.Move.Databases), len(opts.Move.ExternalVolumes), len(opts.Move.Shares)) != 1 {
			errs = append(errs, NewError("MOVE must set exactly one of DATABASES, EXTERNAL VOLUMES, and SHARES"))
		}
		if !ValidObjectIdentifier(opts.Move.ToReplicationGroup) {
			errs = append(errs, errInvalidIdentifier("AlterReplicationGroupOptions.Move", "ToReplicationGroup"))
		}
	}
	if opts.Add != nil {
		if countSetLists(len(opts.Add.ToAllowedDatabases), len(opts.Add.ToAllowedExternalVolumes), len(opts.Add.ToAllowedShares), len(opts.Add.ToAllowedAccounts)) != 1 {
			errs = append(errs, NewError("ADD must set exactly one allowed list"))
		}
		if opts.Add.IgnoreEditionCheck != nil && *opts.Add.IgnoreEditionCheck && len(opts.Add.ToAllowedAccounts) == 0 {
			errs = append(errs, NewError("IGNORE EDITION CHECK requires ADD TO ALLOWED_ACCOUNTS"))
		}
	}
	if opts.Remove != nil {
		if countSetLists(len(opts.Remove.FromAllowedDatabases), len(opts.Remove.FromAllowedExternalVolumes), len(opts.Remove.FromAllowedShares), len(opts.Remove.FromAllowedAccounts)) != 1 {
			errs = append(errs, NewError("REMOVE must set exactly one allowed list"))
		}
	}
	return JoinErrors(errs...)
}

func countSetLists(lengths ...int) int {
	n := 0
	for _, length := range lengths {
		if length > 0 {
			n++
		}
	}
	return n
}

func (opts *CreateSecondaryReplicationGroupOptions) additionalValidations() error {
	if opts.AsReplicaOf.objectIdentifier == nil || !ValidObjectIdentifier(opts.AsReplicaOf) {
		return errInvalidIdentifier("CreateSecondaryReplicationGroupOptions", "AsReplicaOf")
	}
	return nil
}

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

func ParseCommaSeparatedReplicationGroupObjectTypeArray(value string) ([]ReplicationGroupObjectType, error) {
	return collections.MapErr(ParseCommaSeparatedStringArray(value, false), ToReplicationGroupObjectType)
}

func ParseCommaSeparatedReplicationGroupIntegrationTypeArray(value string) ([]ReplicationGroupIntegrationType, error) {
	return collections.MapErr(ParseCommaSeparatedStringArray(value, false), ToReplicationGroupIntegrationType)
}
