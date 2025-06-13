package sdk

import (
	"errors"
	"fmt"
	"slices"
)

var (
	_ validatable = new(GrantPrivilegesToAccountRoleOptions)
	_ validatable = new(RevokePrivilegesFromAccountRoleOptions)
	_ validatable = new(GrantPrivilegesToDatabaseRoleOptions)
	_ validatable = new(RevokePrivilegesFromDatabaseRoleOptions)
	_ validatable = new(grantPrivilegeToShareOptions)
	_ validatable = new(revokePrivilegeFromShareOptions)
	_ validatable = new(GrantOwnershipOptions)
	_ validatable = new(ShowGrantOptions)
)

// based on https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#required-parameters
var validGrantOwnershipObjectTypes = []ObjectType{
	ObjectTypeAggregationPolicy,
	ObjectTypeAlert,
	ObjectTypeAuthenticationPolicy,
	ObjectTypeComputePool,
	ObjectTypeDataMetricFunction,
	ObjectTypeDatabase,
	ObjectTypeDatabaseRole,
	ObjectTypeDynamicTable,
	ObjectTypeEventTable,
	ObjectTypeExternalTable,
	ObjectTypeExternalVolume,
	ObjectTypeFailoverGroup,
	ObjectTypeFileFormat,
	ObjectTypeFunction,
	ObjectTypeGitRepository,
	ObjectTypeHybridTable,
	ObjectTypeIcebergTable,
	ObjectTypeImageRepository,
	ObjectTypeIntegration,
	ObjectTypeMaterializedView,
	ObjectTypeNetworkPolicy,
	ObjectTypeNetworkRule,
	ObjectTypePackagesPolicy,
	ObjectTypePipe,
	ObjectTypeProcedure,
	ObjectTypeMaskingPolicy,
	ObjectTypePasswordPolicy,
	ObjectTypeProjectionPolicy,
	ObjectTypeReplicationGroup,
	ObjectTypeResourceMonitor,
	ObjectTypeRole,
	ObjectTypeRowAccessPolicy,
	ObjectTypeSchema,
	ObjectTypeSessionPolicy,
	ObjectTypeSecret,
	ObjectTypeSequence,
	ObjectTypeStage,
	ObjectTypeStream,
	ObjectTypeTable,
	ObjectTypeTag,
	ObjectTypeTask,
	ObjectTypeUser,
	ObjectTypeView,
	ObjectTypeWarehouse,
}

// based on https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#required-parameters
var validGrantToObjectTypes = []ObjectType{
	ObjectTypeAggregationPolicy,
	ObjectTypeAlert,
	ObjectTypeAuthenticationPolicy,
	ObjectTypeCortexSearchService,
	ObjectTypeDataMetricFunction,
	ObjectTypeDynamicTable,
	ObjectTypeEventTable,
	ObjectTypeExternalTable,
	ObjectTypeFileFormat,
	ObjectTypeFunction,
	ObjectTypeGitRepository,
	ObjectTypeHybridTable,
	ObjectTypeImageRepository,
	ObjectTypeIcebergTable,
	ObjectTypeMaskingPolicy,
	ObjectTypeMaterializedView,
	ObjectTypeModel,
	ObjectTypeNetworkRule,
	ObjectTypeNotebook,
	ObjectTypePackagesPolicy,
	ObjectTypePasswordPolicy,
	ObjectTypePipe,
	ObjectTypeProcedure,
	ObjectTypeProjectionPolicy,
	ObjectTypeRowAccessPolicy,
	ObjectTypeSecret,
	ObjectTypeService,
	ObjectTypeSessionPolicy,
	ObjectTypeSequence,
	ObjectTypeSnapshot,
	ObjectTypeStage,
	ObjectTypeStream,
	ObjectTypeTable,
	ObjectTypeTag,
	ObjectTypeTask,
	ObjectTypeView,
	ObjectTypeStreamlit, // added because of https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2656
	ObjectTypeDataset,   // added because of https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2807
}

// based on https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#restrictions-and-limitations
var invalidGrantToFutureObjectTypes = []ObjectType{
	ObjectTypeComputePool,
	ObjectTypeExternalFunction,
	ObjectTypeImageRepository,
	ObjectTypeAggregationPolicy,
	ObjectTypeMaskingPolicy,
	ObjectTypePackagesPolicy,
	ObjectTypeProjectionPolicy,
	ObjectTypeRowAccessPolicy,
	ObjectTypeSessionPolicy,
	ObjectTypeTag,
}

var (
	ValidGrantOwnershipObjectTypesString       = make([]string, len(validGrantOwnershipObjectTypes))
	ValidGrantOwnershipPluralObjectTypesString = make([]string, len(validGrantOwnershipObjectTypes))
	ValidGrantToObjectTypesString              = make([]string, len(validGrantToObjectTypes))
	ValidGrantToPluralObjectTypesString        = make([]string, len(validGrantToObjectTypes))
	ValidGrantToFuturePluralObjectTypesString  = make([]string, 0)
)

func init() {
	for i, objectType := range validGrantOwnershipObjectTypes {
		ValidGrantOwnershipObjectTypesString[i] = objectType.String()
		ValidGrantOwnershipPluralObjectTypesString[i] = objectType.Plural().String()
	}
	for i, objectType := range validGrantToObjectTypes {
		ValidGrantToObjectTypesString[i] = objectType.String()
		ValidGrantToPluralObjectTypesString[i] = objectType.Plural().String()
		if !slices.Contains(invalidGrantToFutureObjectTypes, objectType) {
			ValidGrantToFuturePluralObjectTypesString = append(ValidGrantToFuturePluralObjectTypesString, objectType.Plural().String())
		}
	}
}

func (opts *GrantPrivilegesToAccountRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !valueSet(opts.privileges) {
		errs = append(errs, errNotSet("GrantPrivilegesToAccountRoleOptions", "privileges"))
	} else {
		if err := opts.privileges.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !valueSet(opts.on) {
		errs = append(errs, errNotSet("GrantPrivilegesToAccountRoleOptions", "on"))
	} else {
		if err := opts.on.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *AccountRoleGrantPrivileges) validate() error {
	if !exactlyOneValueSet(v.AllPrivileges, v.GlobalPrivileges, v.AccountObjectPrivileges, v.SchemaPrivileges, v.SchemaObjectPrivileges) {
		return errExactlyOneOf("AccountRoleGrantPrivileges", "AllPrivileges", "GlobalPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges")
	}
	return nil
}

func (v *AccountRoleGrantOn) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.Account, v.AccountObject, v.Schema, v.SchemaObject) {
		errs = append(errs, errExactlyOneOf("AccountRoleGrantOn", "Account", "AccountObject", "Schema", "SchemaObject"))
	}
	if valueSet(v.AccountObject) {
		if err := v.AccountObject.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.Schema) {
		if err := v.Schema.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.SchemaObject) {
		if err := v.SchemaObject.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *GrantOnAccountObject) validate() error {
	if !exactlyOneValueSet(v.User, v.ResourceMonitor, v.Warehouse, v.ComputePool, v.Database, v.Integration, v.FailoverGroup, v.ReplicationGroup, v.ExternalVolume) {
		return errExactlyOneOf("GrantOnAccountObject", "User", "ResourceMonitor", "Warehouse", "ComputePool", "Database", "Integration", "FailoverGroup", "ReplicationGroup", "ExternalVolume")
	}
	return nil
}

func (v *GrantOnSchema) validate() error {
	if !exactlyOneValueSet(v.Schema, v.AllSchemasInDatabase, v.FutureSchemasInDatabase) {
		return errExactlyOneOf("GrantOnSchema", "Schema", "AllSchemasInDatabase", "FutureSchemasInDatabase")
	}
	return nil
}

func (v *GrantOnSchemaObject) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.SchemaObject, v.All, v.Future) {
		errs = append(errs, errExactlyOneOf("GrantOnSchemaObject", "SchemaObject", "All", "Future"))
	}
	if valueSet(v.All) {
		if err := v.All.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.Future) {
		if err := v.Future.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *GrantOnSchemaObjectIn) validate() error {
	if !exactlyOneValueSet(v.InDatabase, v.InSchema) {
		return errExactlyOneOf("GrantOnSchemaObjectIn", "InDatabase", "InSchema")
	}
	return nil
}

func (opts *RevokePrivilegesFromAccountRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !valueSet(opts.privileges) {
		errs = append(errs, errNotSet("RevokePrivilegesFromAccountRoleOptions", "privileges"))
	} else {
		if err := opts.privileges.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !valueSet(opts.on) {
		errs = append(errs, errNotSet("RevokePrivilegesFromAccountRoleOptions", "on"))
	} else {
		if err := opts.on.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !ValidObjectIdentifier(opts.accountRole) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.Restrict, opts.Cascade) {
		errs = append(errs, errOneOf("RevokePrivilegesFromAccountRoleOptions", "Restrict", "Cascade"))
	}
	return errors.Join(errs...)
}

func (opts *GrantPrivilegesToDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !valueSet(opts.privileges) {
		errs = append(errs, errNotSet("GrantPrivilegesToDatabaseRoleOptions", "privileges"))
	} else {
		if err := opts.privileges.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !valueSet(opts.on) {
		errs = append(errs, errNotSet("GrantPrivilegesToDatabaseRoleOptions", "on"))
	} else {
		if err := opts.on.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *DatabaseRoleGrantPrivileges) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.DatabasePrivileges, v.SchemaPrivileges, v.SchemaObjectPrivileges, v.AllPrivileges) {
		errs = append(errs, errExactlyOneOf("DatabaseRoleGrantPrivileges", "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges", "AllPrivileges"))
	}
	return errors.Join(errs...)
}

func (v *DatabaseRoleGrantOn) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.Database, v.Schema, v.SchemaObject) {
		errs = append(errs, errExactlyOneOf("DatabaseRoleGrantOn", "Database", "Schema", "SchemaObject"))
	}
	if valueSet(v.Schema) {
		if err := v.Schema.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.SchemaObject) {
		if err := v.SchemaObject.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (opts *RevokePrivilegesFromDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !valueSet(opts.privileges) {
		errs = append(errs, errNotSet("RevokePrivilegesFromDatabaseRoleOptions", "privileges"))
	} else {
		if err := opts.privileges.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !valueSet(opts.on) {
		errs = append(errs, errNotSet("RevokePrivilegesFromDatabaseRoleOptions", "on"))
	} else {
		if err := opts.on.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if !ValidObjectIdentifier(opts.databaseRole) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueSet(opts.Restrict, opts.Cascade) {
		errs = append(errs, errOneOf("RevokePrivilegesFromDatabaseRoleOptions", "Restrict", "Cascade"))
	}
	return errors.Join(errs...)
}

func (opts *grantPrivilegeToShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.to) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.On) || len(opts.privileges) == 0 {
		errs = append(errs, fmt.Errorf("on and privilege are required"))
	}
	if valueSet(opts.On) {
		if err := opts.On.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *ShareGrantOn) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.Database, v.Schema, v.Function, v.Table, v.Tag, v.View) {
		errs = append(errs, errExactlyOneOf("ShareGrantOn", "Database", "Schema", "Function", "Table", "Tag", "View"))
	}
	if valueSet(v.Table) {
		if err := v.Table.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *OnTable) validate() error {
	if !exactlyOneValueSet(v.Name, v.AllInSchema) {
		return errExactlyOneOf("OnTable", "Name", "AllInSchema")
	}
	return nil
}

func (opts *revokePrivilegeFromShareOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.from) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !valueSet(opts.On) || len(opts.privileges) == 0 {
		errs = append(errs, errNotSet("revokePrivilegeFromShareOptions", "On", "privileges"))
	}
	if valueSet(opts.On) {
		if err := opts.On.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *OnView) validate() error {
	if !exactlyOneValueSet(v.Name, v.AllInSchema) {
		return errExactlyOneOf("OnView", "Name", "AllInSchema")
	}
	return nil
}

func (opts *GrantOwnershipOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if valueSet(opts.On) {
		if err := opts.On.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(opts.To) {
		if err := opts.To.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *OwnershipGrantOn) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.Object, v.All, v.Future) {
		errs = append(errs, errExactlyOneOf("OwnershipGrantOn", "Object", "AllIn", "Future"))
	}
	if valueSet(v.All) {
		if err := v.All.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	if valueSet(v.Future) {
		if err := v.Future.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (v *OwnershipGrantTo) validate() error {
	if !exactlyOneValueSet(v.DatabaseRoleName, v.AccountRoleName) {
		return errExactlyOneOf("OwnershipGrantTo", "databaseRoleName", "accountRoleName")
	}
	return nil
}

// TODO: add validations for ShowGrantsOn, ShowGrantsTo, ShowGrantsOf and ShowGrantsIn
func (opts *ShowGrantOptions) validate() error {
	if moreThanOneValueSet(opts.On, opts.To, opts.Of, opts.In) {
		return errOneOf("ShowGrantOptions", "On", "To", "Of", "In")
	}
	return nil
}
