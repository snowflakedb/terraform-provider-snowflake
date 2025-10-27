package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var taskSchema = map[string]*schema.Schema{
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the task."),
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the task."),
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the task; must be unique for the database and schema in which the task is created."),
	},
	"started": {
		Type:     schema.TypeBool,
		Required: true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShowWithMapping("state", func(state any) any {
			stateEnum, err := sdk.ToTaskState(state.(string))
			if err != nil {
				return false
			}
			return stateEnum == sdk.TaskStateStarted
		}),
		Description: "Specifies if the task should be started or suspended.",
	},
	"warehouse": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      relatedResourceDescription("The warehouse the task will use. Omit this parameter to use Snowflake-managed compute resources for runs of this task. Due to Snowflake limitations warehouse identifier can consist of only upper-cased letters. (Conflicts with user_task_managed_initial_warehouse_size)", resources.Warehouse),
		ConflictsWith:    []string{"user_task_managed_initial_warehouse_size"},
	},
	"schedule": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		Description:   "The schedule for periodically running the task. This can be a cron or interval in minutes. (Conflicts with finalize and after; when set, one of the sub-fields `minutes` or `using_cron` should be set)",
		ConflictsWith: []string{"finalize", "after"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"minutes": {
					Type:             schema.TypeInt,
					Optional:         true,
					Description:      "Specifies an interval (in minutes) of wait time inserted between runs of the task. Accepts positive integers only. (conflicts with `using_cron`)",
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
					ExactlyOneOf:     []string{"schedule.0.minutes", "schedule.0.using_cron"},
				},
				"using_cron": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Specifies a cron expression and time zone for periodically running the task. Supports a subset of standard cron utility syntax. (conflicts with `minutes`)",
					DiffSuppressFunc: ignoreCaseSuppressFunc,
					ExactlyOneOf:     []string{"schedule.0.minutes", "schedule.0.using_cron"},
				},
			},
		},
	},
	"config": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("config"),
		Description:      "Specifies a string representation of key value pairs that can be accessed by all tasks in the task graph. Must be in JSON format.",
	},
	"allow_overlapping_execution": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("allow_overlapping_execution"),
		Description:      booleanStringFieldDescription("By default, Snowflake ensures that only one instance of a particular DAG is allowed to run at a time, setting the parameter value to TRUE permits DAG runs to overlap."),
	},
	"error_integration": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("error_integration")),
		Description:      relatedResourceDescription(blocklistedCharactersFieldDescription("Specifies the name of the notification integration used for error notifications."), resources.NotificationIntegration),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the task.",
	},
	"finalize": {
		Optional:         true,
		Type:             schema.TypeString,
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(
			suppressIdentifierQuoting,
			IgnoreChangeToCurrentSnowflakeValueInShow("task_relations.0.finalized_root_task"),
		),
		Description:   blocklistedCharactersFieldDescription("Specifies the name of a root task that the finalizer task is associated with. Finalizer tasks run after all other tasks in the task graph run to completion. You can define the SQL of a finalizer task to handle notifications and the release and cleanup of resources that a task graph uses. For more information, see [Release and cleanup of task graphs](https://docs.snowflake.com/en/user-guide/tasks-graphs.html#label-finalizer-task)."),
		ConflictsWith: []string{"schedule", "after"},
	},
	"after": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		},
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("after"),
		Description:      blocklistedCharactersFieldDescription("Specifies one or more predecessor tasks for the current task. Use this option to [create a DAG](https://docs.snowflake.com/en/user-guide/tasks-graphs.html#label-task-dag) of tasks or add this task to an existing DAG. A DAG is a series of tasks that starts with a scheduled root task and is linked together by dependencies."),
		ConflictsWith:    []string{"schedule", "finalize"},
	},
	"when": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: SuppressIfAny(DiffSuppressStatement, IgnoreChangeToCurrentSnowflakeValueInShow("condition")),
		Description:      "Specifies a Boolean SQL expression; multiple conditions joined with AND/OR are supported. When a task is triggered (based on its SCHEDULE or AFTER setting), it validates the conditions of the expression to determine whether to execute. If the conditions of the expression are not met, then the task skips the current run. Any tasks that identify this task as a predecessor also don’t run.",
	},
	"sql_statement": {
		Type:             schema.TypeString,
		Required:         true,
		DiffSuppressFunc: SuppressIfAny(DiffSuppressStatement, IgnoreChangeToCurrentSnowflakeValueInShow("definition")),
		Description:      "Any single SQL statement, or a call to a stored procedure, executed when the task runs.",
	},
	"target_completion_interval": {
		Type:             schema.TypeString,
		Optional:         true,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("target_lag"),
		Description:      "Specifies the target completion interval for tasks. Format: '<num> { HOURS | MINUTES | SECONDS }' (e.g., '10 MINUTES'). For more information, check [TARGET_COMPLETION_INTERVAL docs](https://docs.snowflake.com/en/sql-reference/parameters#target-completion-interval).",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW TASKS` for the given task.",
		Elem: &schema.Resource{
			Schema: schemas.ShowTaskSchema,
		},
	},
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN TASK` for the given task.",
		Elem: &schema.Resource{
			Schema: schemas.ShowTaskParametersSchema,
		},
	},
}

func Task() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.Task, CreateTask),
		UpdateContext: TrackingUpdateWrapper(resources.Task, UpdateTask),
		ReadContext:   TrackingReadWrapper(resources.Task, ReadTask(true)),
		DeleteContext: TrackingDeleteWrapper(resources.Task, DeleteTask),
		Description:   "Resource used to manage task objects. For more information, check [task documentation](https://docs.snowflake.com/en/user-guide/tasks-intro).",

		Schema: collections.MergeMaps(taskSchema, taskParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Task, ImportTask),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Task, customdiff.All(
			ComputedIfAnyAttributeChanged(taskSchema, ShowOutputAttributeName, "name", "started", "warehouse", "user_task_managed_initial_warehouse_size", "schedule", "config", "allow_overlapping_execution", "error_integration", "comment", "finalize", "after", "when", "target_completion_interval"),
			ComputedIfAnyAttributeChanged(taskParametersSchema, ParametersAttributeName, collections.Map(sdk.AsStringList(sdk.AllTaskParameters), strings.ToLower)...),
			ComputedIfAnyAttributeChanged(taskSchema, FullyQualifiedNameAttributeName, "name"),
			taskParametersCustomDiff,
		)),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    cty.EmptyObject,
				Upgrade: v098TaskStateUpgrader,
			},
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportTask(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	task, err := client.Tasks.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := ImportName[sdk.SchemaObjectIdentifier](context.Background(), d, nil); err != nil {
		return nil, err
	}

	if err := d.Set("allow_overlapping_execution", booleanStringFromBool(task.AllowOverlappingExecution)); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateTask(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	req := sdk.NewCreateTaskRequest(id, d.Get("sql_statement").(string))
	tasksToResume := make([]sdk.SchemaObjectIdentifier, 0)

	if errs := errors.Join(
		attributeMappedValueCreate(d, "warehouse", &req.Warehouse, func(v any) (*sdk.CreateTaskWarehouseRequest, error) {
			warehouseId, err := sdk.ParseAccountObjectIdentifier(v.(string))
			if err != nil {
				return nil, err
			}
			return sdk.NewCreateTaskWarehouseRequest().WithWarehouse(&warehouseId), nil
		}),
		attributeMappedValueCreate(d, "schedule", &req.Schedule, func(v any) (*string, error) {
			if len(v.([]any)) > 0 {
				if minutes, ok := d.GetOk("schedule.0.minutes"); ok {
					return sdk.String(fmt.Sprintf("%d MINUTE", minutes)), nil
				}
				if cron, ok := d.GetOk("schedule.0.using_cron"); ok {
					return sdk.String(fmt.Sprintf("USING CRON %s", cron)), nil
				}
				return nil, fmt.Errorf("when setting a schedule either minutes or using_cron field should be set")
			}
			return nil, nil
		}),
		stringAttributeCreate(d, "config", &req.Config),
		booleanStringAttributeCreate(d, "allow_overlapping_execution", &req.AllowOverlappingExecution),
		accountObjectIdentifierAttributeCreate(d, "error_integration", &req.ErrorIntegration),
		stringAttributeCreate(d, "comment", &req.Comment),
		stringAttributeCreate(d, "when", &req.When),
		stringAttributeCreate(d, "target_completion_interval", &req.TargetCompletionInterval),
	); errs != nil {
		return diag.FromErr(errs)
	}

	if v, ok := d.GetOk("finalize"); ok {
		rootTaskId, err := sdk.ParseSchemaObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		rootTask, err := client.Tasks.ShowByID(ctx, rootTaskId)
		if err != nil {
			return diag.FromErr(err)
		}

		if rootTask.IsStarted() {
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTaskId).WithSuspend(true)); err != nil {
				return diag.FromErr(sdk.JoinErrors(err))
			}
			tasksToResume = append(tasksToResume, rootTaskId)
		}

		req.WithFinalize(rootTaskId)
	}

	if v, ok := d.GetOk("after"); ok {
		after := expandStringList(v.(*schema.Set).List())
		precedingTasks := make([]sdk.SchemaObjectIdentifier, 0)
		for _, parentTaskIdString := range after {
			parentTaskId, err := sdk.ParseSchemaObjectIdentifier(parentTaskIdString)
			if err != nil {
				return diag.FromErr(err)
			}
			resumeTasks, err := client.Tasks.SuspendRootTasks(ctx, parentTaskId, id)
			tasksToResume = append(tasksToResume, resumeTasks...)
			if err != nil {
				return diag.FromErr(sdk.JoinErrors(err))
			}
			precedingTasks = append(precedingTasks, parentTaskId)
		}
		req.WithAFTER(precedingTasks)
	}

	if parameterCreateDiags := handleTaskParametersCreate(d, req); len(parameterCreateDiags) > 0 {
		return parameterCreateDiags
	}

	if err := client.Tasks.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	if d.Get("started").(bool) {
		if err := waitForTaskStart(ctx, client, id); err != nil {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Failed to start the task",
					Detail:   fmt.Sprintf("Id: %s, err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		// Else case not handled, because tasks are created as suspended (https://docs.snowflake.com/en/sql-reference/sql/create-task; "important" section)
	}

	defer func() {
		if err := client.Tasks.ResumeTasks(ctx, tasksToResume); err != nil {
			diags = append(diags, resumeTaskErrorDiag(id, "create", err))
		}
	}()

	return append(diags, ReadTask(false)(ctx, d, meta)...)
}

func UpdateTask(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	task, err := client.Tasks.ShowByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, id, id)
	if err != nil {
		return diag.FromErr(sdk.JoinErrors(err))
	}

	defer func() {
		if err := client.Tasks.ResumeTasks(ctx, tasksToResume); err != nil {
			diags = append(diags, resumeTaskErrorDiag(id, "create", err))
		}
	}()

	if task.IsStarted() {
		if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithSuspend(true)); err != nil {
			return diag.FromErr(sdk.JoinErrors(err))
		}
	}

	unset := sdk.NewTaskUnsetRequest()
	set := sdk.NewTaskSetRequest()

	err = errors.Join(
		attributeMappedValueUpdate(d, "user_task_managed_initial_warehouse_size", &set.UserTaskManagedInitialWarehouseSize, &unset.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		accountObjectIdentifierAttributeUpdate(d, "warehouse", &set.Warehouse, &unset.Warehouse),
		stringAttributeUpdate(d, "config", &set.Config, &unset.Config),
		booleanStringAttributeUpdate(d, "allow_overlapping_execution", &set.AllowOverlappingExecution, &unset.AllowOverlappingExecution),
		accountObjectIdentifierAttributeUpdate(d, "error_integration", &set.ErrorIntegration, &unset.ErrorIntegration),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		stringAttributeUpdate(d, "target_completion_interval", &set.TargetCompletionInterval, &unset.TargetCompletionInterval),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("schedule") {
		_, newSchedule := d.GetChange("schedule")

		if newSchedule != nil && len(newSchedule.([]any)) == 1 {
			if _, newMinutes := d.GetChange("schedule.0.minutes"); newMinutes.(int) > 0 {
				set.Schedule = sdk.String(fmt.Sprintf("%d MINUTE", newMinutes.(int)))
			}
			if _, newCron := d.GetChange("schedule.0.using_cron"); newCron.(string) != "" {
				set.Schedule = sdk.String(fmt.Sprintf("USING CRON %s", newCron.(string)))
			}
		} else {
			unset.Schedule = sdk.Bool(true)
		}
	}

	if updateDiags := handleTaskParametersUpdate(d, set, unset); len(updateDiags) > 0 {
		return updateDiags
	}

	if *unset != (sdk.TaskUnsetRequest{}) {
		if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("when") {
		if v := d.Get("when"); v != "" {
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithModifyWhen(v.(string))); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithRemoveWhen(true)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("sql_statement") {
		if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithModifyAs(d.Get("sql_statement").(string))); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("finalize") {
		if v, ok := d.GetOk("finalize"); ok {
			rootTaskId, err := sdk.ParseSchemaObjectIdentifier(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}

			rootTask, err := client.Tasks.ShowByID(ctx, rootTaskId)
			if err != nil {
				return diag.FromErr(err)
			}

			if rootTask.IsStarted() {
				if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTaskId).WithSuspend(true)); err != nil {
					return diag.FromErr(sdk.JoinErrors(err))
				}
			}

			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithSetFinalize(rootTaskId)); err != nil {
				return diag.FromErr(err)
			}

			if rootTask.IsStarted() && !slices.ContainsFunc(tasksToResume, func(identifier sdk.SchemaObjectIdentifier) bool {
				return identifier.FullyQualifiedName() == rootTaskId.FullyQualifiedName()
			}) {
				tasksToResume = append(tasksToResume, rootTaskId)
			}
		} else {
			rootTask, err := client.Tasks.ShowByID(ctx, *task.TaskRelations.FinalizedRootTask)
			if err != nil {
				return diag.FromErr(err)
			}

			if rootTask.IsStarted() {
				if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(rootTask.ID()).WithSuspend(true)); err != nil {
					return diag.FromErr(sdk.JoinErrors(err))
				}
			}

			if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithUnsetFinalize(true)); err != nil {
				return diag.FromErr(err)
			}

			if rootTask.IsStarted() && !slices.ContainsFunc(tasksToResume, func(identifier sdk.SchemaObjectIdentifier) bool {
				return identifier.FullyQualifiedName() == rootTask.ID().FullyQualifiedName()
			}) {
				tasksToResume = append(tasksToResume, rootTask.ID())
			}
		}
	}

	if d.HasChange("after") {
		oldAfter, newAfter := d.GetChange("after")
		addedTasks, removedTasks := ListDiff(
			expandStringList(oldAfter.(*schema.Set).List()),
			expandStringList(newAfter.(*schema.Set).List()),
		)

		if len(addedTasks) > 0 {
			addedTaskIds, err := collections.MapErr(addedTasks, sdk.ParseSchemaObjectIdentifier)
			if err != nil {
				return diag.FromErr(err)
			}

			for _, addedTaskId := range addedTaskIds {
				addedTasksToResume, err := client.Tasks.SuspendRootTasks(ctx, addedTaskId, sdk.NewSchemaObjectIdentifier("", "", ""))
				tasksToResume = append(tasksToResume, addedTasksToResume...)
				if err != nil {
					return diag.FromErr(sdk.JoinErrors(err))
				}
			}

			err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithAddAfter(addedTaskIds))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(removedTasks) > 0 {
			removedTaskIds, err := collections.MapErr(removedTasks, sdk.ParseSchemaObjectIdentifier)
			if err != nil {
				return diag.FromErr(err)
			}
			err = client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithRemoveAfter(removedTaskIds))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if *set != (sdk.TaskSetRequest{}) {
		if err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("started").(bool) {
		if err := waitForTaskStart(ctx, client, id); err != nil {
			return diag.FromErr(fmt.Errorf("failed to resume task %s, err = %w", id.FullyQualifiedName(), err))
		}
	}
	// We don't process the else case, because the task was already suspended at the beginning of the Update method.
	tasksToResume = slices.DeleteFunc(tasksToResume, func(identifier sdk.SchemaObjectIdentifier) bool {
		return identifier.FullyQualifiedName() == id.FullyQualifiedName()
	})

	return append(diags, ReadTask(false)(ctx, d, meta)...)
}

func ReadTask(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		task, err := client.Tasks.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query task. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Task id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		taskParameters, err := client.Tasks.ShowParameters(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"allow_overlapping_execution", "allow_overlapping_execution", task.AllowOverlappingExecution, booleanStringFromBool(task.AllowOverlappingExecution), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err = setStateToValuesFromConfig(d, taskSchema, []string{
				"allow_overlapping_execution",
			}); err != nil {
				return diag.FromErr(err)
			}
		}

		if errs := errors.Join(
			attributeMappedValueReadOrDefault(d, "finalize", task.TaskRelations.FinalizedRootTask, func(finalizedRootTask *sdk.SchemaObjectIdentifier) (string, error) {
				return finalizedRootTask.FullyQualifiedName(), nil
			}, nil),
			attributeMappedValueReadOrDefault(d, "error_integration", task.ErrorIntegration, func(errorIntegration *sdk.AccountObjectIdentifier) (string, error) {
				return errorIntegration.Name(), nil
			}, nil),
			attributeMappedValueReadOrDefault(d, "warehouse", task.Warehouse, func(warehouse *sdk.AccountObjectIdentifier) (string, error) {
				return warehouse.Name(), nil
			}, nil),
			func() error {
				if len(task.Schedule) > 0 {
					taskSchedule, err := sdk.ParseTaskSchedule(task.Schedule)
					if err != nil {
						return err
					}
					switch {
					case len(taskSchedule.Cron) > 0:
						if err := d.Set("schedule", []any{map[string]any{
							"using_cron": taskSchedule.Cron,
						}}); err != nil {
							return err
						}
					case taskSchedule.Minutes > 0:
						if err := d.Set("schedule", []any{map[string]any{
							"minutes": taskSchedule.Minutes,
						}}); err != nil {
							return err
						}
					}
					return nil
				}
				return d.Set("schedule", nil)
			}(),
			d.Set("started", task.IsStarted()),
			d.Set("when", task.Condition),
			d.Set("config", task.Config),
			d.Set("comment", task.Comment),
			d.Set("sql_statement", task.Definition),
			d.Set("after", collections.Map(task.TaskRelations.Predecessors, sdk.SchemaObjectIdentifier.FullyQualifiedName)),
			handleTaskParameterRead(d, taskParameters),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.TaskToSchema(task)}),
			d.Set(ParametersAttributeName, []map[string]any{schemas.TaskParametersToSchema(taskParameters)}),
		); errs != nil {
			return diag.FromErr(errs)
		}

		return nil
	}
}

func DeleteTask(ctx context.Context, d *schema.ResourceData, meta any) (diags diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tasksToResume, err := client.Tasks.SuspendRootTasks(ctx, id, id)
	defer func() {
		if err := client.Tasks.ResumeTasks(ctx, tasksToResume); err != nil {
			diags = append(diags, resumeTaskErrorDiag(id, "delete", err))
		}
	}()
	if err != nil {
		return diag.FromErr(sdk.JoinErrors(err))
	}

	if err = client.Tasks.DropSafely(ctx, id); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting task %s err = %w", id.FullyQualifiedName(), err))
	}

	d.SetId("")
	return diags
}

func resumeTaskErrorDiag(id sdk.SchemaObjectIdentifier, operation string, originalErr error) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  fmt.Sprintf("Failed to resume tasks in %s operation (id=%s)", operation, id.FullyQualifiedName()),
		Detail:   fmt.Sprintf("Failed to resume some of the tasks with the following errors (tasks can be resumed by applying the same configuration again): %v", originalErr),
	}
}

func waitForTaskStart(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier) error {
	err := client.Tasks.Alter(ctx, sdk.NewAlterTaskRequest(id).WithResume(true))
	if err != nil {
		return fmt.Errorf("error starting task %s err = %w", id.FullyQualifiedName(), err)
	}
	return util.Retry(5, 5*time.Second, func() (error, bool) {
		task, err := client.Tasks.ShowByID(ctx, id)
		if err != nil {
			return fmt.Errorf("error starting task %s err = %w", id.FullyQualifiedName(), err), false
		}
		if task.State != sdk.TaskStateStarted {
			return nil, false
		}
		return nil, true
	})
}
