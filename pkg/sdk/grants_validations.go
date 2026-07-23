package sdk

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var (
	_ validatable = new(GrantPrivilegesToAccountRoleOptions)
	_ validatable = new(grantInheritedPrivilegesToAccountRoleOptions)
	_ validatable = new(RevokePrivilegesFromAccountRoleOptions)
	_ validatable = new(revokeInheritedPrivilegesFromAccountRoleOptions)
	_ validatable = new(GrantPrivilegesToDatabaseRoleOptions)
	_ validatable = new(grantInheritedPrivilegesToDatabaseRoleOptions)
	_ validatable = new(RevokePrivilegesFromDatabaseRoleOptions)
	_ validatable = new(revokeInheritedPrivilegesFromDatabaseRoleOptions)
	_ validatable = new(grantPrivilegeToShareOptions)
	_ validatable = new(revokePrivilegeFromShareOptions)
	_ validatable = new(GrantOwnershipOptions)
	_ validatable = new(RevokeOwnershipOptions)
	_ validatable = new(ShowGrantOptions)
)

// based on https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#required-parameters
var validGrantOwnershipObjectTypes = []ObjectType{
	ObjectTypeAgent,
	ObjectTypeAggregationPolicy,
	ObjectTypeAlert,
	ObjectTypeAuthenticationPolicy,
	ObjectTypeComputePool,
	ObjectTypeCortexSearchService,
	ObjectTypeDataMetricFunction,
	ObjectTypeDatabase,
	ObjectTypeDatabaseRole,
	ObjectTypeDbtProject,
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
	ObjectTypeSemanticView,
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

// Database roles are excluded
var validGrantOwnershipBulkObjectTypes = []ObjectType{
	ObjectTypeAgent,
	ObjectTypeAggregationPolicy,
	ObjectTypeAlert,
	ObjectTypeAuthenticationPolicy,
	ObjectTypeComputePool,
	ObjectTypeCortexSearchService,
	ObjectTypeDataMetricFunction,
	ObjectTypeDatabase,
	ObjectTypeDbtProject,
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
	ObjectTypeSemanticView,
	ObjectTypeSequence,
	ObjectTypeStage,
	ObjectTypeStream,
	ObjectTypeTable,
	ObjectTypeTag,
	ObjectTypeTask,
	ObjectTypeUser,
	ObjectTypeView,
	ObjectTypeWarehouse,
	ObjectTypeWorkspace,
}

var validGrantToAccountObjectTypes = []ObjectType{
	ObjectTypeUser,
	ObjectTypeResourceMonitor,
	ObjectTypeWarehouse,
	ObjectTypeComputePool,
	ObjectTypeDatabase,
	ObjectTypeIntegration,
	ObjectTypeConnection,
	ObjectTypeFailoverGroup,
	ObjectTypeReplicationGroup,
	ObjectTypeExternalVolume,
}

// based on https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#required-parameters
var validGrantToSchemaObjectTypes = []ObjectType{
	ObjectTypeAgent,
	ObjectTypeAggregationPolicy,
	ObjectTypeAlert,
	ObjectTypeAuthenticationPolicy,
	ObjectTypeCortexSearchService,
	ObjectTypeDataMetricFunction,
	ObjectTypeDataset,
	ObjectTypeDbtProject,
	ObjectTypeDynamicTable,
	ObjectTypeEventTable,
	ObjectTypeExperiment,
	ObjectTypeExternalTable,
	ObjectTypeFileFormat,
	ObjectTypeFunction,
	ObjectTypeGateway,
	ObjectTypeGitRepository,
	ObjectTypeHybridTable,
	ObjectTypeImageRepository,
	ObjectTypeIcebergTable,
	ObjectTypeJoinPolicy,
	ObjectTypeMaskingPolicy,
	ObjectTypeMaterializedView,
	ObjectTypeMcpServer,
	ObjectTypeModel,
	ObjectTypeModelMonitor,
	ObjectTypeNetworkRule,
	ObjectTypeNotebook,
	ObjectTypeNotebookProject,
	ObjectTypeOnlineFeatureTable,
	ObjectTypePackagesPolicy,
	ObjectTypePasswordPolicy,
	ObjectTypePipe,
	ObjectTypePrivacyPolicy,
	ObjectTypeProcedure,
	ObjectTypeProjectionPolicy,
	ObjectTypeRowAccessPolicy,
	ObjectTypeSecret,
	ObjectTypeSemanticView,
	ObjectTypeService,
	ObjectTypeSessionPolicy,
	ObjectTypeSequence,
	ObjectTypeSnapshot,
	ObjectTypeSnapshotPolicy,
	ObjectTypeSnapshotSet,
	ObjectTypeStage,
	ObjectTypeStorageLifecyclePolicy,
	ObjectTypeStream,
	ObjectTypeStreamlit,
	ObjectTypeTable,
	ObjectTypeTag,
	ObjectTypeTask,
	ObjectTypeView,
	ObjectTypeWorkspace,
}

// TODO(SNOW-2370066): Adjust after adding tests
// based on https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#restrictions-and-limitations
var invalidGrantToAllObjectTypes = []ObjectType{
	ObjectTypeComputePool,
	ObjectTypeExperiment,
	ObjectTypeExternalFunction,
	ObjectTypeGateway,
	ObjectTypeJoinPolicy,
	ObjectTypeNotebookProject,
	// ObjectTypeAggregationPolicy,
	// ObjectTypeMaskingPolicy,
	// ObjectTypePackagesPolicy,
	// ObjectTypeProjectionPolicy,
	// ObjectTypeRowAccessPolicy,
	// ObjectTypeSessionPolicy,
	ObjectTypeStorageLifecyclePolicy,
	// ObjectTypeTag,
	ObjectTypeWarehouse,
}

// based on https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#restrictions-and-limitations
var invalidGrantToFutureObjectTypes = []ObjectType{
	ObjectTypeComputePool,
	ObjectTypeExperiment,
	ObjectTypeExternalFunction,
	ObjectTypeGateway,
	ObjectTypeAggregationPolicy,
	ObjectTypeJoinPolicy,
	ObjectTypeNotebookProject,
	ObjectTypeMaskingPolicy,
	ObjectTypePackagesPolicy,
	ObjectTypeProjectionPolicy,
	ObjectTypeRowAccessPolicy,
	ObjectTypeSessionPolicy,
	ObjectTypeSnapshot,
	ObjectTypeStorageLifecyclePolicy,
	ObjectTypeTag,
	ObjectTypeWarehouse,
}

var (
	ValidGrantOwnershipObjectTypesString       = make([]string, len(validGrantOwnershipObjectTypes))
	ValidGrantOwnershipPluralObjectTypesString = make([]string, len(validGrantOwnershipBulkObjectTypes))
	ValidGrantToAccountObjectTypesString       = make([]string, len(validGrantToAccountObjectTypes))
	ValidGrantToAccountObjectPluralTypesString = make([]string, len(validGrantToAccountObjectTypes))
	ValidGrantToSchemaObjectTypesString        = make([]string, len(validGrantToSchemaObjectTypes))
	ValidGrantToAllPluralObjectTypesString     = make([]string, 0)
	ValidGrantToFuturePluralObjectTypesString  = make([]string, 0)
)

func init() {
	for i, objectType := range validGrantOwnershipObjectTypes {
		ValidGrantOwnershipObjectTypesString[i] = objectType.String()
	}
	for i, objectType := range validGrantOwnershipBulkObjectTypes {
		ValidGrantOwnershipPluralObjectTypesString[i] = objectType.Plural().String()
	}
	for i, objectType := range validGrantToAccountObjectTypes {
		ValidGrantToAccountObjectTypesString[i] = objectType.String()
		ValidGrantToAccountObjectPluralTypesString[i] = objectType.Plural().String()
	}
	for i, objectType := range validGrantToSchemaObjectTypes {
		ValidGrantToSchemaObjectTypesString[i] = objectType.String()
		if !slices.Contains(invalidGrantToAllObjectTypes, objectType) {
			ValidGrantToAllPluralObjectTypesString = append(ValidGrantToAllPluralObjectTypesString, objectType.Plural().String())
		}
		if !slices.Contains(invalidGrantToFutureObjectTypes, objectType) {
			ValidGrantToFuturePluralObjectTypesString = append(ValidGrantToFuturePluralObjectTypesString, objectType.Plural().String())
		}
	}
}

// allowedUnquotedCharactersRegex matches non-empty strings consisting only of allowed characters
var allowedUnquotedCharactersRegex = regexp.MustCompile(`^[a-zA-Z ._]+$`)

// validateUnquotedInput checks that the passed string contains only allowed characters.
func validateUnquotedInput(s string) error {
	var errs []error
	if !allowedUnquotedCharactersRegex.MatchString(s) {
		errs = append(errs, fmt.Errorf("%s contains disallowed characters; it must follow this regex: %s", s, allowedUnquotedCharactersRegex.String()))
	}
	return errors.Join(errs...)
}

// validatePrivileges checks that every passed privilege contains only allowed characters.
func validatePrivileges[T fmt.Stringer](privileges []T) error {
	var errs []error
	for _, privilege := range privileges {
		if err := validateUnquotedInput(privilege.String()); err != nil {
			errs = append(errs, fmt.Errorf("invalid privilege: %w", err))
		}
	}
	return errors.Join(errs...)
}

// ToPrivilege converts a string to a privilege name.
// It should be used instead of the raw privilege conversion whenever the input is not trusted.
// There is no dedicated privilege type in the SDK, so we use string instead.
func ToPrivilege(s string) (string, error) {
	s = strings.ToUpper(s)
	if err := validateUnquotedInput(s); err != nil {
		return "", fmt.Errorf("invalid privilege: %w", err)
	}
	return s, nil
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
	return errors.Join(
		validatePrivileges(v.GlobalPrivileges),
		validatePrivileges(v.AccountObjectPrivileges),
		validatePrivileges(v.SchemaPrivileges),
		validatePrivileges(v.SchemaObjectPrivileges),
	)
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
	if !exactlyOneValueSet(v.User, v.ResourceMonitor, v.Warehouse, v.ComputePool, v.Database, v.Integration, v.Connection, v.FailoverGroup, v.ReplicationGroup, v.ExternalVolume) {
		return errExactlyOneOf("GrantOnAccountObject", "User", "ResourceMonitor", "Warehouse", "ComputePool", "Database", "Integration", "Connection", "FailoverGroup", "ReplicationGroup", "ExternalVolume")
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

func (opts *grantInheritedPrivilegesToAccountRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if err := opts.privileges.validate(); err != nil {
		errs = append(errs, err)
	}
	if opts.onAll == "" {
		errs = append(errs, errNotSet("grantInheritedPrivilegesToAccountRoleOptions", "onAll"))
	}
	if err := opts.in.validate(); err != nil {
		errs = append(errs, err)
	}
	if !ValidObjectIdentifier(opts.accountRole) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (v InheritedAccountRoleGrantPrivileges) validate() error {
	if !exactlyOneValueSet(v.AllPrivileges, v.AccountObjectPrivileges, v.SchemaPrivileges, v.SchemaObjectPrivileges) {
		return errExactlyOneOf("InheritedAccountRoleGrantPrivileges", "AllPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges")
	}
	return errors.Join(
		validatePrivileges(v.AccountObjectPrivileges),
		validatePrivileges(v.SchemaPrivileges),
		validatePrivileges(v.SchemaObjectPrivileges),
	)
}

func (v InheritedAccountRoleGrantIn) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.Account, v.Database, v.Schema) {
		errs = append(errs, errExactlyOneOf("InheritedAccountRoleGrantIn", "Account", "Database", "Schema"))
	}
	if v.Database != nil && !ValidObjectIdentifier(*v.Database) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if v.Schema != nil && !ValidObjectIdentifier(*v.Schema) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
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

func (opts *revokeInheritedPrivilegesFromAccountRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if err := opts.privileges.validate(); err != nil {
		errs = append(errs, err)
	}
	if opts.onAll == "" {
		errs = append(errs, errNotSet("revokeInheritedPrivilegesFromAccountRoleOptions", "onAll"))
	}
	if err := opts.in.validate(); err != nil {
		errs = append(errs, err)
	}
	if !ValidObjectIdentifier(opts.accountRole) {
		errs = append(errs, ErrInvalidObjectIdentifier)
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
	errs = append(errs, validatePrivileges(v.DatabasePrivileges))
	errs = append(errs, validatePrivileges(v.SchemaPrivileges))
	errs = append(errs, validatePrivileges(v.SchemaObjectPrivileges))
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

func (opts *grantInheritedPrivilegesToDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if err := opts.privileges.validate(); err != nil {
		errs = append(errs, err)
	}
	if opts.onAll == "" {
		errs = append(errs, errNotSet("grantInheritedPrivilegesToDatabaseRoleOptions", "onAll"))
	}
	if err := opts.in.validate(); err != nil {
		errs = append(errs, err)
	}
	if !ValidObjectIdentifier(opts.databaseRole) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	return errors.Join(errs...)
}

func (v InheritedDatabaseRoleGrantPrivileges) validate() error {
	if !exactlyOneValueSet(v.AllPrivileges, v.SchemaPrivileges, v.SchemaObjectPrivileges) {
		return errExactlyOneOf("InheritedDatabaseRoleGrantPrivileges", "AllPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges")
	}
	return errors.Join(
		validatePrivileges(v.SchemaPrivileges),
		validatePrivileges(v.SchemaObjectPrivileges),
	)
}

func (v InheritedDatabaseRoleGrantIn) validate() error {
	var errs []error
	if !exactlyOneValueSet(v.Database, v.Schema) {
		errs = append(errs, errExactlyOneOf("InheritedDatabaseRoleGrantIn", "Database", "Schema"))
	}
	if v.Database != nil && !ValidObjectIdentifier(*v.Database) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if v.Schema != nil && !ValidObjectIdentifier(*v.Schema) {
		errs = append(errs, ErrInvalidObjectIdentifier)
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

func (opts *revokeInheritedPrivilegesFromDatabaseRoleOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if err := opts.privileges.validate(); err != nil {
		errs = append(errs, err)
	}
	if opts.onAll == "" {
		errs = append(errs, errNotSet("revokeInheritedPrivilegesFromDatabaseRoleOptions", "onAll"))
	}
	if err := opts.in.validate(); err != nil {
		errs = append(errs, err)
	}
	if !ValidObjectIdentifier(opts.databaseRole) {
		errs = append(errs, ErrInvalidObjectIdentifier)
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
	errs = append(errs, validatePrivileges(opts.privileges))
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
	errs = append(errs, validatePrivileges(opts.privileges))
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

func (opts *RevokeOwnershipOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if err := opts.On.validate(); err != nil {
		errs = append(errs, err)
	}
	if err := opts.From.validate(); err != nil {
		errs = append(errs, err)
	}
	if everyValueSet(opts.Restrict, opts.Cascade) {
		errs = append(errs, errOneOf("RevokeOwnershipOptions", "Restrict", "Cascade"))
	}
	return errors.Join(errs...)
}

func (v *RevokeOwnershipGrantOn) validate() error {
	var errs []error
	// Snowflake only supports revoking OWNERSHIP for future grants; ownership of existing objects must be
	// transferred with GRANT OWNERSHIP instead, hence Future is the only allowed variant here.
	if !valueSet(v.Future) {
		errs = append(errs, errNotSet("RevokeOwnershipGrantOn", "Future"))
	} else if err := v.Future.validate(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// TODO: add validations for ShowGrantsOn, ShowGrantsTo, ShowGrantsOf and ShowGrantsIn
func (opts *ShowGrantOptions) validate() error {
	var errs []error
	if moreThanOneValueSet(opts.On, opts.To, opts.Of, opts.In) {
		errs = append(errs, errOneOf("ShowGrantOptions", "On", "To", "Of", "In"))
	}
	if everyValueSet(opts.Inherited, opts.Future) {
		errs = append(errs, errOneOf("ShowGrantOptions", "Inherited", "Future"))
	}
	return errors.Join(errs...)
}
