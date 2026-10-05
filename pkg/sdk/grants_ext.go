package sdk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

type OwnershipCurrentGrantsOutboundPrivileges string

const (
	Revoke OwnershipCurrentGrantsOutboundPrivileges = "REVOKE"
	Copy   OwnershipCurrentGrantsOutboundPrivileges = "COPY"
)

func (v *AccountRoleGrantPrivilegesRequest) ToInheritedAccountRoleGrantPrivileges() InheritedAccountRoleGrantPrivilegesRequest {
	return InheritedAccountRoleGrantPrivilegesRequest{
		AccountObjectPrivileges: v.AccountObjectPrivileges,
		SchemaPrivileges:        v.SchemaPrivileges,
		SchemaObjectPrivileges:  v.SchemaObjectPrivileges,
		AllPrivileges:           v.AllPrivileges,
	}
}

func (v *DatabaseRoleGrantPrivilegesRequest) ToInheritedDatabaseRoleGrantPrivileges() InheritedDatabaseRoleGrantPrivilegesRequest {
	return InheritedDatabaseRoleGrantPrivilegesRequest{
		SchemaPrivileges:       v.SchemaPrivileges,
		SchemaObjectPrivileges: v.SchemaObjectPrivileges,
		AllPrivileges:          v.AllPrivileges,
	}
}

func (v *Grant) ID() ObjectIdentifier {
	return v.Name
}

// ObjectTypeFromShowGrants maps GRANTED_ON / GRANT_ON from SHOW GRANTS to an ObjectType.
// Snowflake sometimes returns a shortened name (VOLUME, POSTGRES) instead of the SQL type
// (EXTERNAL VOLUME, POSTGRES INSTANCE).
func ObjectTypeFromShowGrants(raw string) (ObjectType, error) {
	if raw == "" {
		return "", nil
	}
	switch raw {
	case "VOLUME":
		return ObjectTypeExternalVolume, nil
	case "MODULE":
		return ObjectTypeModel, nil
	case "CORTEX_AGENT":
		return ObjectTypeAgent, nil
	case "CORTEX_AGENT_SERVER":
		return ObjectTypeMcpServer, nil
	case "QUALITY_MONITOR":
		return ObjectTypeModelMonitor, nil
	case "POSTGRES":
		return ObjectTypePostgresInstance, nil
	default:
		return ObjectType(strings.ReplaceAll(raw, "_", " ")), nil
	}
}

func trimQuotes(raw string) (*string, error) {
	if raw == "" {
		return nil, nil
	}
	trimmed := strings.Trim(raw, `"`)
	return &trimmed, nil
}

func (r grantRow) additionalConvert(result *Grant) error {
	var name ObjectIdentifier
	var err error
	// TODO(SNOW-1569535): use a mapper from object type to parsing function
	if ObjectType(r.GrantedOn).IsWithArguments() {
		name, err = ParseSchemaObjectIdentifierWithArgumentsAndReturnType(r.Name)
	} else {
		name, err = ParseObjectIdentifierString(r.Name)
	}
	if err != nil {
		log.Printf("[DEBUG] Failed to parse identifier [%s], err = \"%s\"; falling back to fully qualified name conversion", r.Name, err)
		name = NewObjectIdentifierFromFullyQualifiedName(r.Name)
	}
	result.Name = name
	return nil
}

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
	ObjectTypeInteractiveTable,
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
	ObjectTypeSnowflakeIntelligence,
	ObjectTypeStage,
	ObjectTypeStream,
	ObjectTypeTable,
	ObjectTypeTag,
	ObjectTypeTask,
	ObjectTypeUser,
	ObjectTypeView,
	ObjectTypeWarehouse,
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
	ObjectTypeSnowflakeIntelligence,
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
	ObjectTypeInteractiveTable,
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

var (
	ValidGrantOwnershipObjectTypesString = make([]string, len(validGrantOwnershipObjectTypes))
	ValidGrantToAccountObjectTypesString = make([]string, len(validGrantToAccountObjectTypes))
	ValidGrantToSchemaObjectTypesString  = make([]string, len(validGrantToSchemaObjectTypes))
)

func init() {
	for i, objectType := range validGrantOwnershipObjectTypes {
		ValidGrantOwnershipObjectTypesString[i] = objectType.String()
	}
	for i, objectType := range validGrantToAccountObjectTypes {
		ValidGrantToAccountObjectTypesString[i] = objectType.String()
	}
	for i, objectType := range validGrantToSchemaObjectTypes {
		ValidGrantToSchemaObjectTypesString[i] = objectType.String()
	}
}

// allowedUnquotedCharactersRegex matches non-empty strings consisting only of allowed characters
var allowedUnquotedCharactersRegex = regexp.MustCompile(`^[a-zA-Z ._]+$`)

// validateUserInput checks that the passed string contains only allowed characters.
func validateUserInput(s string) error {
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
		if err := validateUserInput(privilege.String()); err != nil {
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
	if err := validateUserInput(s); err != nil {
		return "", fmt.Errorf("invalid privilege: %w", err)
	}
	return s, nil
}

// The additionalValidations methods below hold the custom validation logic that the generator cannot
// express (privilege character checks and object-type checks).

func (v *AccountRoleGrantPrivileges) additionalValidations() error {
	return errors.Join(
		validatePrivileges(v.GlobalPrivileges),
		validatePrivileges(v.AccountObjectPrivileges),
		validatePrivileges(v.SchemaPrivileges),
		validatePrivileges(v.SchemaObjectPrivileges),
	)
}

func (v *DatabaseRoleGrantPrivileges) additionalValidations() error {
	return errors.Join(
		validatePrivileges(v.DatabasePrivileges),
		validatePrivileges(v.SchemaPrivileges),
		validatePrivileges(v.SchemaObjectPrivileges),
	)
}

func (v InheritedAccountRoleGrantPrivileges) additionalValidations() error {
	return errors.Join(
		validatePrivileges(v.AccountObjectPrivileges),
		validatePrivileges(v.SchemaPrivileges),
		validatePrivileges(v.SchemaObjectPrivileges),
	)
}

func (v InheritedDatabaseRoleGrantPrivileges) additionalValidations() error {
	return errors.Join(
		validatePrivileges(v.SchemaPrivileges),
		validatePrivileges(v.SchemaObjectPrivileges),
	)
}

func (v *GrantOnAccountObject) additionalValidations() error {
	if v.Object == nil {
		return nil
	}
	if err := validateUserInput(v.Object.ObjectType.String()); err != nil {
		return fmt.Errorf("invalid object type: %w", err)
	}
	return nil
}

func (opts *GrantPrivilegeToShareOptions) additionalValidations() error {
	return validatePrivileges(opts.Privileges)
}

func (opts *RevokePrivilegeFromShareOptions) additionalValidations() error {
	return validatePrivileges(opts.Privileges)
}

func (v *grants) RevokeInheritedPrivilegesFromAccountRoleSafely(ctx context.Context, request *RevokeInheritedPrivilegesFromAccountRoleRequest) error {
	return SafeRevokePrivileges(func() error {
		return v.RevokeInheritedPrivilegesFromAccountRole(ctx, request)
	})
}

func (v *grants) RevokeInheritedPrivilegesFromDatabaseRoleSafely(ctx context.Context, request *RevokeInheritedPrivilegesFromDatabaseRoleRequest) error {
	return SafeRevokePrivileges(func() error {
		return v.RevokeInheritedPrivilegesFromDatabaseRole(ctx, request)
	})
}

// pipeSchemaObjectOn builds a per-object "ON <schema object>" request used when expanding a
// bulk "ON ALL PIPES" grant/revoke into one statement per pipe.
func pipeSchemaObjectOn(pipe Pipe) *GrantOnSchemaObjectRequest {
	return NewGrantOnSchemaObjectRequest().WithSchemaObject(Object{ObjectType: ObjectTypePipe, Name: pipe.ID()})
}

func (v *grants) GrantPrivilegesToAccountRole(ctx context.Context, request *GrantPrivilegesToAccountRoleRequest) error {
	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "grant x on all pipes" operation, we'll go and grant specified privileges
	// to every Pipe one by one.
	if on := request.On; on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.GrantPrivilegesToAccountRole(ctx, &GrantPrivilegesToAccountRoleRequest{
					Privileges:      request.Privileges,
					On:              NewAccountRoleGrantOnRequest().WithSchemaObject(*pipeSchemaObjectOn(pipe)),
					AccountRole:     request.AccountRole,
					WithGrantOption: request.WithGrantOption,
				})
			},
		)
	}

	return v.grantPrivilegesToAccountRole(ctx, request)
}

func (v *grants) RevokePrivilegesFromAccountRole(ctx context.Context, request *RevokePrivilegesFromAccountRoleRequest) error {
	return v.revokeAccountRolePrivileges(ctx, request, noopExecWrapper)
}

func (v *grants) RevokePrivilegesFromAccountRoleSafely(ctx context.Context, request *RevokePrivilegesFromAccountRoleRequest) error {
	return v.revokeAccountRolePrivileges(ctx, request, SafeRevokePrivileges)
}

func noopExecWrapper(f func() error) error {
	return f()
}

func (v *grants) revokeAccountRolePrivileges(
	ctx context.Context,
	request *RevokePrivilegesFromAccountRoleRequest,
	execWrapper func(func() error) error,
) error {
	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "revoke x on all pipes" operation, we'll go and revoke specified privileges
	// from every Pipe one by one.
	if on := request.On; on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.revokeAccountRolePrivileges(ctx, &RevokePrivilegesFromAccountRoleRequest{
					GrantOptionFor: request.GrantOptionFor,
					Privileges:     request.Privileges,
					On:             NewAccountRoleGrantOnRequest().WithSchemaObject(*pipeSchemaObjectOn(pipe)),
					AccountRole:    request.AccountRole,
					Restrict:       request.Restrict,
					Cascade:        request.Cascade,
				}, execWrapper)
			},
		)
	}

	return execWrapper(func() error {
		return v.revokePrivilegesFromAccountRole(ctx, request)
	})
}

func (v *grants) GrantPrivilegesToDatabaseRole(ctx context.Context, request *GrantPrivilegesToDatabaseRoleRequest) error {
	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "grant x on all pipes" operation, we'll go and grant specified privileges
	// to every Pipe one by one.
	if on := request.On; on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.GrantPrivilegesToDatabaseRole(ctx, &GrantPrivilegesToDatabaseRoleRequest{
					Privileges:      request.Privileges,
					On:              NewDatabaseRoleGrantOnRequest().WithSchemaObject(*pipeSchemaObjectOn(pipe)),
					DatabaseRole:    request.DatabaseRole,
					WithGrantOption: request.WithGrantOption,
				})
			},
		)
	}

	return v.grantPrivilegesToDatabaseRole(ctx, request)
}

func (v *grants) RevokePrivilegesFromDatabaseRole(ctx context.Context, request *RevokePrivilegesFromDatabaseRoleRequest) error {
	return v.revokeDatabaseRolePrivileges(ctx, request, noopExecWrapper)
}

func (v *grants) RevokePrivilegesFromDatabaseRoleSafely(ctx context.Context, request *RevokePrivilegesFromDatabaseRoleRequest) error {
	return v.revokeDatabaseRolePrivileges(ctx, request, SafeRevokePrivileges)
}

func (v *grants) revokeDatabaseRolePrivileges(
	ctx context.Context,
	request *RevokePrivilegesFromDatabaseRoleRequest,
	execWrapper func(func() error) error,
) error {
	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "revoke x on all pipes" operation, we'll go and revoke specified privileges
	// from every Pipe one by one.
	if on := request.On; on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.revokeDatabaseRolePrivileges(ctx, &RevokePrivilegesFromDatabaseRoleRequest{
					GrantOptionFor: request.GrantOptionFor,
					Privileges:     request.Privileges,
					On:             NewDatabaseRoleGrantOnRequest().WithSchemaObject(*pipeSchemaObjectOn(pipe)),
					DatabaseRole:   request.DatabaseRole,
					Restrict:       request.Restrict,
					Cascade:        request.Cascade,
				}, execWrapper)
			},
		)
	}

	return execWrapper(func() error {
		return v.revokePrivilegesFromDatabaseRole(ctx, request)
	})
}

func (v *grants) RevokePrivilegeFromShareSafely(ctx context.Context, request *RevokePrivilegeFromShareRequest) error {
	return SafeRevokePrivileges(func() error {
		return v.RevokePrivilegeFromShare(ctx, request)
	})
}

func (v *grants) GrantOwnership(ctx context.Context, request *GrantOwnershipRequest) error {
	if on := request.On; on.Object != nil && on.Object.ObjectType == ObjectTypePipe {
		return v.grantOwnershipOnPipe(ctx, on.Object.Name.(SchemaObjectIdentifier), request)
	}

	if on := request.On; on.Object != nil && on.Object.ObjectType == ObjectTypeTask {
		return v.grantOwnershipOnTask(ctx, on.Object.Name.(SchemaObjectIdentifier), request)
	}

	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "grant x on all pipes" operation, we'll go and grant specified privileges
	// to every Pipe one by one.
	if on := request.On; on.All != nil && on.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.All.InDatabase,
			on.All.InSchema,
			func(pipe Pipe) error {
				return v.GrantOwnership(ctx, &GrantOwnershipRequest{
					On:            OwnershipGrantOnRequest{Object: &Object{ObjectType: ObjectTypePipe, Name: pipe.ID()}},
					To:            request.To,
					CurrentGrants: request.CurrentGrants,
				})
			},
		)
	}

	return v.grantOwnership(ctx, request)
}

// TODO(SNOW-2097063): Improve SHOW GRANTS implementation
func (v *grants) Show(ctx context.Context, request *ShowGrantsRequest) ([]Grant, error) {
	if request == nil {
		request = NewShowGrantsRequest()
	}
	opts := request.toOpts()

	dbRows, err := validateAndQuery[grantRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList, err := convertRows[grantRow, Grant](dbRows)
	if err != nil {
		return nil, err
	}
	for i, grant := range resultList {
		// SHOW GRANTS of DATABASE ROLE requires a special handling:
		// - it returns no account name, so for other SHOW GRANTS types it needs to be skipped
		// - it returns fully qualified name for database objects
		granteeNameRaw := dbRows[i].GranteeName
		if !(valueSet(opts.Of) && valueSet(opts.Of.DatabaseRole)) { //nolint:gocritic
			granteeName := granteeNameRaw
			if grant.GrantedTo == ObjectTypeShare {
				granteeName = normalizeShareGranteeName(granteeName, v.client.GetAccountLocator())
			}
			resultList[i].GranteeName = NewAccountObjectIdentifier(granteeName)
		} else if !slices.Contains([]ObjectType{ObjectTypeRole, ObjectTypeShare, ObjectTypeUser, ObjectTypeApplication}, grant.GrantedTo) {
			// TODO(SNOW-3954332): cleanup when 2026_06 bundle is GA (enforced and can no longer be disabled).
			// Before BCR-2371 the grantee (another database role) is fully qualified (<database>.<database_role>); after, only <database_role>.
			// A database role can only be granted within the same database, so reconstruct the missing prefix from the queried role.
			if id, err := ParseDatabaseObjectIdentifier(granteeNameRaw); err == nil {
				resultList[i].GranteeName = id
			} else if id, err := ParseAccountObjectIdentifier(granteeNameRaw); err == nil {
				resultList[i].GranteeName = NewDatabaseObjectIdentifier(opts.Of.DatabaseRole.DatabaseName(), id.Name())
			} else {
				return nil, err
			}
		} else if grant.GrantedTo == ObjectTypeShare {
			resultList[i].GranteeName = NewAccountObjectIdentifier(normalizeShareGranteeName(granteeNameRaw, v.client.GetAccountLocator()))
		} else if grant.GrantedTo == ObjectTypeUser {
			resultList[i].GranteeName = NewAccountObjectIdentifier(strings.TrimPrefix(granteeNameRaw, "USER$"))
		} else {
			resultList[i].GranteeName = NewAccountObjectIdentifier(granteeNameRaw)
		}
	}
	return resultList, nil
}

// SQLKey renders the exact SHOW GRANTS statement that Show would execute for this request, so it can
// be used as a cache key.
func (r *ShowGrantsRequest) SQLKey() (string, error) {
	if r == nil {
		r = NewShowGrantsRequest()
	}
	return structToSQL(r.toOpts())
}

func normalizeShareGranteeName(granteeName string, accountLocator string) string {
	if accountLocator == "" {
		return granteeName
	}
	prefix := accountLocator + "."
	if len(granteeName) >= len(prefix) && strings.EqualFold(granteeName[:len(prefix)], prefix) {
		return granteeName[len(prefix):]
	}
	return granteeName
}

// grantOwnershipOnPipe execution sequence
//  1. Get the current role.
//  2. Show grants on the pipe.
//  3. See if the current role can "operate" on the pipe (has either OPERATE or OWNERSHIP privileges granted).
//  4. If the current role can "operate" on the pipe.
//     4.1. Check the current execution status of the pipe.
//     4.2. Pause the pipe execution if it's running.
//  5. If it cannot, try to proceed with the grant ownership in case the pipe is already paused.
//  6. Grant ownership.
//  7. If the current role could "operate" on the pipe, and the ownership was granted with COPY CURRENT GRANTS option.
//     6.1. Resume with the use of system function.
//  8. If it couldn't, notify the user that the pipe has to be resumed manually with the use of system function.
func (v *grants) grantOwnershipOnPipe(ctx context.Context, pipeId SchemaObjectIdentifier, request *GrantOwnershipRequest) error {
	currentRole, err := v.client.ContextFunctions.CurrentRole(ctx)
	if err != nil {
		return err
	}

	currentGrants, err := v.client.Grants.Show(ctx, NewShowGrantsRequest().WithOn(*NewShowGrantsOnRequest().WithObject(Object{
		ObjectType: ObjectTypePipe,
		Name:       pipeId,
	})))
	if err != nil {
		return err
	}

	isGrantedWithPrivilege := func(privilege string) bool {
		return slices.ContainsFunc(currentGrants, func(grant Grant) bool {
			return grant.GranteeName == currentRole.Value &&
				grant.GrantedOn == ObjectTypePipe &&
				grant.Privilege == privilege
		})
	}
	// To be able to call ALTER on a pipe to stop its execution,
	// the current role has to be either the owner (OWNERSHIP privilege) of this pipe or be granted with OPERATE privilege.
	// MONITOR privilege is also needed to be able to check the current pipe execution state.
	canOperateOnPipe := isGrantedWithPrivilege(SchemaObjectPrivilegeOperate.String())
	canMonitorPipe := isGrantedWithPrivilege(SchemaObjectPrivilegeMonitor.String())
	hasOwnershipOnPipe := isGrantedWithPrivilege("OWNERSHIP")

	var originalPipeExecutionState *PipeExecutionState
	if hasOwnershipOnPipe || (canOperateOnPipe && canMonitorPipe) {
		pipeStatus, err := v.client.SystemFunctions.PipeStatus(ctx, NewPipeStatusRequest(*NewPipeStatusArgumentsRequest(pipeId)))
		if err != nil {
			return err
		}
		originalPipeExecutionState = &pipeStatus.ExecutionState

		if pipeStatus.ExecutionState == PipeExecutionStateRunning {
			if err := v.client.Pipes.Alter(ctx, NewAlterPipeRequest(pipeId).WithSet(*NewPipeSetRequest().WithPipeExecutionPaused(true))); err != nil {
				return err
			}
		}
	} else {
		fmt.Printf("[DEBUG] Insufficient permissions to check the status of the pipe (MONITOR privilege): %s, and pause it if it's in running state (OPERATE privilege). Trying to proceed with ownership transfer...", pipeId.FullyQualifiedName())
	}

	if err := v.grantOwnership(ctx, request); err != nil {
		return err
	}

	// If:
	// - The current role was granted with OPERATE privilege before ownership transfer.
	// - GRANT OWNERSHIP command was run with COPY CURRENT GRANTS option.
	// - The pipe was previously running.
	// We can safely use the PIPE_FORCE_RESUME system function to resume the pipe after successful ownership transfer.
	if canOperateOnPipe && request.CurrentGrants != nil && request.CurrentGrants.OutboundPrivileges == Copy && originalPipeExecutionState != nil && *originalPipeExecutionState == PipeExecutionStateRunning {
		if err := v.client.SystemFunctions.PipeForceResume(ctx, NewPipeForceResumeRequest(*NewPipeForceResumeArgumentsRequest(pipeId))); err != nil {
			return err
		}
	} else {
		log.Printf("[WARN] Insufficient privileges to resume the pipe: %s. Resume has to be done manually with the use of SELECT SYSTEM$PIPE_FORCE_RESUME system function.", pipeId.FullyQualifiedName())
	}

	return nil
}

func (v *grants) grantOwnershipOnTask(ctx context.Context, taskId SchemaObjectIdentifier, request *GrantOwnershipRequest) error {
	currentGrantsOnObject, err := v.client.Grants.Show(ctx, NewShowGrantsRequest().WithOn(*NewShowGrantsOnRequest().WithObject(Object{
		ObjectType: ObjectTypeTask,
		Name:       taskId,
	})))
	if err != nil {
		return err
	}

	currentGrantsOnAccount, err := v.client.Grants.Show(ctx, NewShowGrantsRequest().WithOn(*NewShowGrantsOnRequest().WithAccount(true)))
	if err != nil {
		return err
	}

	currentRole, err := v.client.ContextFunctions.CurrentRole(ctx)
	if err != nil {
		return err
	}

	currentTask, err := v.client.Tasks.ShowByID(ctx, taskId)
	if err != nil {
		return err
	}

	isGrantedWithPrivilege := func(collection []Grant, grantedOn ObjectType, privilege string) bool {
		return slices.ContainsFunc(collection, func(grant Grant) bool {
			return grant.GranteeName == currentRole.Value &&
				grant.GrantedOn == grantedOn &&
				grant.Privilege == privilege
		})
	}

	var isGrantedWithWarehouseUsage bool

	if currentTask.Warehouse == nil {
		// For serverless tasks (tasks that are not associated with any warehouse), we don't need to check for warehouse usage privileges.
		isGrantedWithWarehouseUsage = true
	} else {
		currentGrantsOnTaskWarehouse, err := v.client.Grants.Show(ctx, NewShowGrantsRequest().WithOn(*NewShowGrantsOnRequest().WithObject(Object{
			ObjectType: ObjectTypeWarehouse,
			Name:       *currentTask.Warehouse,
		})))
		if err != nil {
			return err
		}

		isGrantedWithWarehouseUsage = isGrantedWithPrivilege(currentGrantsOnTaskWarehouse, ObjectTypeWarehouse, AccountObjectPrivilegeUsage.String())
	}

	canOperateOnTask := isGrantedWithPrivilege(currentGrantsOnObject, ObjectTypeTask, SchemaObjectPrivilegeOperate.String())
	canSuspendTask := canOperateOnTask || isGrantedWithPrivilege(currentGrantsOnObject, ObjectTypeTask, "OWNERSHIP")
	canResumeTask := isGrantedWithWarehouseUsage && canOperateOnTask && isGrantedWithPrivilege(currentGrantsOnAccount, ObjectTypeAccount, GlobalPrivilegeExecuteTask.String())
	canResumeTaskAfterOwnershipTransfer := canResumeTask && ((request.CurrentGrants != nil && request.CurrentGrants.OutboundPrivileges == Copy) || (request.To.AccountRoleName != nil && request.To.AccountRoleName.Name() == currentRole.Value.Name()))

	var tasksToResume []SchemaObjectIdentifier
	if canSuspendTask {
		tasksToResume, err = v.client.Tasks.SuspendRootTasks(ctx, taskId, taskId)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[WARN] Insufficient privileges to operate on task: %s (OPERATE privilege). Trying to proceed with ownership transfer...", taskId.FullyQualifiedName())
	}

	if err := v.grantOwnership(ctx, request); err != nil {
		return err
	}

	if currentTask.IsStarted() && !slices.ContainsFunc(tasksToResume, func(id SchemaObjectIdentifier) bool {
		return id.FullyQualifiedName() == currentTask.ID().FullyQualifiedName()
	}) {
		tasksToResume = append(tasksToResume, currentTask.ID())
	}

	if len(tasksToResume) > 0 {
		if canResumeTaskAfterOwnershipTransfer {
			err = v.client.Tasks.ResumeTasks(ctx, tasksToResume)
			if err != nil {
				return err
			}
		} else {
			tasksToResumeString := collections.Map(tasksToResume, func(id SchemaObjectIdentifier) string { return id.FullyQualifiedName() })
			log.Printf("[WARN] Insufficient privileges to resume tasks: %v (EXECUTE TASK privilege). Tasks have to be resumed manually.", tasksToResumeString)
		}
	}

	return nil
}

func (v *grants) runOnAllPipes(ctx context.Context, inDatabase *AccountObjectIdentifier, inSchema *DatabaseObjectIdentifier, command func(Pipe) error) error {
	var in *In
	switch {
	case inDatabase != nil:
		in = &In{
			Database: *inDatabase,
		}
	case inSchema != nil:
		in = &In{
			Schema: *inSchema,
		}
	}

	showReq := NewShowPipeRequest()
	if in != nil {
		showReq.WithIn(*in)
	}
	pipes, err := v.client.Pipes.Show(ctx, showReq)
	if err != nil {
		return err
	}

	return runOnAll(pipes, command)
}

func runOnAll[T any](collection []T, command func(T) error) error {
	var errs []error
	for _, element := range collection {
		if err := command(element); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
