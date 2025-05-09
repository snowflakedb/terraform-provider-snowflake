package sdk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Grants = (*grants)(nil)

type grants struct {
	client *Client
}

func (v *grants) GrantPrivilegesToAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *GrantPrivilegesToAccountRoleOptions) error {
	if opts == nil {
		opts = &GrantPrivilegesToAccountRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.accountRole = role

	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "grant x on all pipes" operation, we'll go and grant specified privileges
	// to every Pipe one by one.
	if on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.client.Grants.GrantPrivilegesToAccountRole(
					ctx,
					privileges,
					&AccountRoleGrantOn{
						SchemaObject: &GrantOnSchemaObject{
							SchemaObject: &Object{
								ObjectType: ObjectTypePipe,
								Name:       pipe.ID(),
							},
						},
					},
					role,
					opts,
				)
			},
		)
	}

	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) RevokePrivilegesFromAccountRole(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, role AccountObjectIdentifier, opts *RevokePrivilegesFromAccountRoleOptions) error {
	if opts == nil {
		opts = &RevokePrivilegesFromAccountRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.accountRole = role

	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "revoke x on all pipes" operation, we'll go and revoke specified privileges
	// from every Pipe one by one.
	if on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.client.Grants.RevokePrivilegesFromAccountRole(
					ctx,
					privileges,
					&AccountRoleGrantOn{
						SchemaObject: &GrantOnSchemaObject{
							SchemaObject: &Object{
								ObjectType: ObjectTypePipe,
								Name:       pipe.ID(),
							},
						},
					},
					role,
					opts,
				)
			},
		)
	}

	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) GrantPrivilegesToDatabaseRole(ctx context.Context, privileges *DatabaseRoleGrantPrivileges, on *DatabaseRoleGrantOn, role DatabaseObjectIdentifier, opts *GrantPrivilegesToDatabaseRoleOptions) error {
	if opts == nil {
		opts = &GrantPrivilegesToDatabaseRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.databaseRole = role

	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "grant x on all pipes" operation, we'll go and grant specified privileges
	// to every Pipe one by one.
	if on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.client.Grants.GrantPrivilegesToDatabaseRole(
					ctx,
					privileges,
					&DatabaseRoleGrantOn{
						SchemaObject: &GrantOnSchemaObject{
							SchemaObject: &Object{
								ObjectType: ObjectTypePipe,
								Name:       pipe.ID(),
							},
						},
					},
					role,
					opts,
				)
			},
		)
	}

	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) RevokePrivilegesFromDatabaseRole(ctx context.Context, privileges *DatabaseRoleGrantPrivileges, on *DatabaseRoleGrantOn, role DatabaseObjectIdentifier, opts *RevokePrivilegesFromDatabaseRoleOptions) error {
	if opts == nil {
		opts = &RevokePrivilegesFromDatabaseRoleOptions{}
	}
	opts.privileges = privileges
	opts.on = on
	opts.databaseRole = role

	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "revoke x on all pipes" operation, we'll go and revoke specified privileges
	// from every Pipe one by one.
	if on != nil &&
		on.SchemaObject != nil &&
		on.SchemaObject.All != nil &&
		on.SchemaObject.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.SchemaObject.All.InDatabase,
			on.SchemaObject.All.InSchema,
			func(pipe Pipe) error {
				return v.client.Grants.RevokePrivilegesFromDatabaseRole(
					ctx,
					privileges,
					&DatabaseRoleGrantOn{
						SchemaObject: &GrantOnSchemaObject{
							SchemaObject: &Object{
								ObjectType: ObjectTypePipe,
								Name:       pipe.ID(),
							},
						},
					},
					role,
					opts,
				)
			},
		)
	}

	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) GrantPrivilegeToShare(ctx context.Context, privileges []ObjectPrivilege, on *ShareGrantOn, to AccountObjectIdentifier) error {
	opts := &grantPrivilegeToShareOptions{
		privileges: privileges,
		On:         on,
		to:         to,
	}
	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) RevokePrivilegeFromShare(ctx context.Context, privileges []ObjectPrivilege, on *ShareGrantOn, id AccountObjectIdentifier) error {
	opts := &revokePrivilegeFromShareOptions{
		privileges: privileges,
		On:         on,
		from:       id,
	}
	return validateAndExec(v.client, ctx, opts)
}

func (v *grants) GrantOwnership(ctx context.Context, on OwnershipGrantOn, to OwnershipGrantTo, opts *GrantOwnershipOptions) error {
	if opts == nil {
		opts = &GrantOwnershipOptions{}
	}

	opts.On = on
	opts.To = to

	if on.Object != nil && on.Object.ObjectType == ObjectTypePipe {
		return v.grantOwnershipOnPipe(ctx, on.Object.Name.(SchemaObjectIdentifier), opts)
	}

	if on.Object != nil && on.Object.ObjectType == ObjectTypeTask {
		return v.grantOwnershipOnTask(ctx, on.Object.Name.(SchemaObjectIdentifier), opts)
	}

	// Snowflake doesn't allow bulk operations on Pipes. Because of that, when SDK user
	// issues "grant x on all pipes" operation, we'll go and grant specified privileges
	// to every Pipe one by one.
	if on.All != nil && on.All.PluralObjectType == PluralObjectTypePipes {
		return v.runOnAllPipes(
			ctx,
			on.All.InDatabase,
			on.All.InSchema,
			func(pipe Pipe) error {
				return v.client.Grants.GrantOwnership(
					ctx,
					OwnershipGrantOn{
						Object: &Object{
							ObjectType: ObjectTypePipe,
							Name:       pipe.ID(),
						},
					},
					to,
					opts,
				)
			},
		)
	}

	return validateAndExec(v.client, ctx, opts)
}

// TODO(SNOW-2097063): Improve SHOW GRANTS implementation
func (v *grants) Show(ctx context.Context, opts *ShowGrantOptions) ([]Grant, error) {
	if opts == nil {
		opts = &ShowGrantOptions{}
	}

	dbRows, err := validateAndQuery[grantRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[grantRow, Grant](dbRows)
	for i, grant := range resultList {
		// SHOW GRANTS of DATABASE ROLE requires a special handling:
		// - it returns no account name, so for other SHOW GRANTS types it needs to be skipped
		// - it returns fully qualified name for database objects
		granteeNameRaw := dbRows[i].GranteeName
		if !(valueSet(opts.Of) && valueSet(opts.Of.DatabaseRole)) { //nolint:gocritic
			granteeName := granteeNameRaw
			if grant.GrantedTo == ObjectTypeShare {
				granteeName = granteeName[strings.IndexRune(granteeName, '.')+1:]
			}
			resultList[i].GranteeName = NewAccountObjectIdentifier(granteeName)
		} else if !slices.Contains([]ObjectType{ObjectTypeRole, ObjectTypeShare, ObjectTypeUser}, grant.GrantedTo) {
			id, err := ParseDatabaseObjectIdentifier(granteeNameRaw)
			if err != nil {
				return nil, err
			}
			resultList[i].GranteeName = id
		} else if grant.GrantedTo == ObjectTypeUser {
			resultList[i].GranteeName = NewAccountObjectIdentifier(strings.TrimLeft(granteeNameRaw, "USER$"))
		} else {
			resultList[i].GranteeName = NewAccountObjectIdentifier(granteeNameRaw)
		}
	}
	return resultList, nil
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
func (v *grants) grantOwnershipOnPipe(ctx context.Context, pipeId SchemaObjectIdentifier, opts *GrantOwnershipOptions) error {
	currentRole, err := v.client.ContextFunctions.CurrentRole(ctx)
	if err != nil {
		return err
	}

	currentGrants, err := v.client.Grants.Show(ctx, &ShowGrantOptions{
		On: &ShowGrantsOn{
			Object: &Object{
				ObjectType: ObjectTypePipe,
				Name:       pipeId,
			},
		},
	})
	if err != nil {
		return err
	}

	isGrantedWithPrivilege := func(privilege string) bool {
		return slices.ContainsFunc(currentGrants, func(grant Grant) bool {
			return grant.GranteeName == currentRole &&
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
		pipeExecutionState, err := v.client.SystemFunctions.PipeStatus(pipeId)
		if err != nil {
			return err
		}
		originalPipeExecutionState = &pipeExecutionState

		if pipeExecutionState == RunningPipeExecutionState {
			if err := v.client.Pipes.Alter(ctx, pipeId, &AlterPipeOptions{
				Set: &PipeSet{
					PipeExecutionPaused: Bool(true),
				},
			}); err != nil {
				return err
			}
		}
	} else {
		fmt.Printf("[DEBUG] Insufficient permissions to check the status of the pipe (MONITOR privilege): %s, and pause it if it's in running state (OPERATE privilege). Trying to proceed with ownership transfer...", pipeId.FullyQualifiedName())
	}

	if err := validateAndExec(v.client, ctx, opts); err != nil {
		return err
	}

	// If:
	// - The current role was granted with OPERATE privilege before ownership transfer.
	// - GRANT OWNERSHIP command was run with COPY CURRENT GRANTS option.
	// - The pipe was previously running.
	// We can safely use the PIPE_FORCE_RESUME system function to resume the pipe after successful ownership transfer.
	if canOperateOnPipe && opts.CurrentGrants != nil && opts.CurrentGrants.OutboundPrivileges == Copy && originalPipeExecutionState != nil && *originalPipeExecutionState == RunningPipeExecutionState {
		if err := v.client.SystemFunctions.PipeForceResume(pipeId, nil); err != nil {
			return err
		}
	} else {
		log.Printf("[WARN] Insufficient privileges to resume the pipe: %s. Resume has to be done manually with the use of SELECT SYSTEM$PIPE_FORCE_RESUME system function.", pipeId.FullyQualifiedName())
	}

	return nil
}

func (v *grants) grantOwnershipOnTask(ctx context.Context, taskId SchemaObjectIdentifier, opts *GrantOwnershipOptions) error {
	currentGrantsOnObject, err := v.client.Grants.Show(ctx, &ShowGrantOptions{
		On: &ShowGrantsOn{
			Object: &Object{
				ObjectType: ObjectTypeTask,
				Name:       taskId,
			},
		},
	})
	if err != nil {
		return err
	}

	currentGrantsOnAccount, err := v.client.Grants.Show(ctx, &ShowGrantOptions{
		On: &ShowGrantsOn{
			Account: Bool(true),
		},
	})
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

	if currentTask.Warehouse == nil {
		return fmt.Errorf("no warehouse found to be attached to the task: %s", taskId.FullyQualifiedName())
	}

	currentGrantsOnTaskWarehouse, err := v.client.Grants.Show(ctx, &ShowGrantOptions{
		On: &ShowGrantsOn{
			Object: &Object{
				ObjectType: ObjectTypeWarehouse,
				Name:       *currentTask.Warehouse,
			},
		},
	})
	if err != nil {
		return err
	}

	isGrantedWithPrivilege := func(collection []Grant, grantedOn ObjectType, privilege string) bool {
		return slices.ContainsFunc(collection, func(grant Grant) bool {
			return grant.GranteeName == currentRole &&
				grant.GrantedOn == grantedOn &&
				grant.Privilege == privilege
		})
	}
	canOperateOnTask := isGrantedWithPrivilege(currentGrantsOnObject, ObjectTypeTask, SchemaObjectPrivilegeOperate.String())
	isGrantedWithWarehouseUsage := isGrantedWithPrivilege(currentGrantsOnTaskWarehouse, ObjectTypeWarehouse, AccountObjectPrivilegeUsage.String())
	canSuspendTask := canOperateOnTask || isGrantedWithPrivilege(currentGrantsOnObject, ObjectTypeTask, "OWNERSHIP")
	canResumeTask := isGrantedWithWarehouseUsage && canOperateOnTask && isGrantedWithPrivilege(currentGrantsOnAccount, ObjectTypeAccount, GlobalPrivilegeExecuteTask.String())
	canResumeTaskAfterOwnershipTransfer := canResumeTask && ((opts.CurrentGrants != nil && opts.CurrentGrants.OutboundPrivileges == Copy) || (opts.To.AccountRoleName != nil && opts.To.AccountRoleName.Name() == currentRole.Name()))

	var tasksToResume []SchemaObjectIdentifier
	if canSuspendTask {
		tasksToResume, err = v.client.Tasks.SuspendRootTasks(ctx, taskId, taskId)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[WARN] Insufficient privileges to operate on task: %s (OPERATE privilege). Trying to proceed with ownership transfer...", taskId.FullyQualifiedName())
	}

	if err := validateAndExec(v.client, ctx, opts); err != nil {
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

func (v *grants) grantTemporarily(ctx context.Context, privileges *AccountRoleGrantPrivileges, on *AccountRoleGrantOn, accountRoleName AccountObjectIdentifier) (func() error, error) {
	return func() error {
			return v.client.Grants.RevokePrivilegesFromAccountRole(
				ctx,
				privileges,
				on,
				accountRoleName,
				new(RevokePrivilegesFromAccountRoleOptions),
			)
		}, v.client.Grants.GrantPrivilegesToAccountRole(
			ctx,
			privileges,
			on,
			accountRoleName,
			new(GrantPrivilegesToAccountRoleOptions),
		)
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

	pipes, err := v.client.Pipes.Show(ctx, &ShowPipeOptions{In: in})
	if err != nil {
		return err
	}

	return runOnAll(pipes, command)
}

func (v *grants) runOnAllTasks(ctx context.Context, inDatabase *AccountObjectIdentifier, inSchema *DatabaseObjectIdentifier, command func(Task) error) error {
	var in In
	switch {
	case inDatabase != nil:
		in = In{
			Database: *inDatabase,
		}
	case inSchema != nil:
		in = In{
			Schema: *inSchema,
		}
	}

	tasks, err := v.client.Tasks.Show(ctx, NewShowTaskRequest().WithIn(ExtendedIn{In: in}))
	if err != nil {
		return err
	}

	return runOnAll(tasks, command)
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
