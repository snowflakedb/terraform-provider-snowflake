package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

type TaskRelationsRepresentation struct {
	Predecessors      []string `json:"Predecessors"`
	FinalizerTask     string   `json:"FinalizerTask"`
	FinalizedRootTask string   `json:"FinalizedRootTask"`
}

func (r *TaskRelationsRepresentation) ToTaskRelations() (TaskRelations, error) {
	predecessors := make([]SchemaObjectIdentifier, len(r.Predecessors))
	for i, predecessor := range r.Predecessors {
		id, err := ParseSchemaObjectIdentifier(predecessor)
		if err != nil {
			return TaskRelations{}, err
		}
		predecessors[i] = id
	}

	taskRelations := TaskRelations{
		Predecessors: predecessors,
	}

	if len(r.FinalizerTask) > 0 {
		finalizerTask, err := ParseSchemaObjectIdentifier(r.FinalizerTask)
		if err != nil {
			return TaskRelations{}, err
		}
		taskRelations.FinalizerTask = &finalizerTask
	}

	if len(r.FinalizedRootTask) > 0 {
		finalizedRootTask, err := ParseSchemaObjectIdentifier(r.FinalizedRootTask)
		if err != nil {
			return TaskRelations{}, err
		}
		taskRelations.FinalizedRootTask = &finalizedRootTask
	}

	return taskRelations, nil
}

type TaskRelations struct {
	Predecessors      []SchemaObjectIdentifier
	FinalizerTask     *SchemaObjectIdentifier
	FinalizedRootTask *SchemaObjectIdentifier
}

func ToTaskRelations(s string) (TaskRelations, error) {
	var taskRelationsRepresentation TaskRelationsRepresentation
	if err := json.Unmarshal([]byte(s), &taskRelationsRepresentation); err != nil {
		return TaskRelations{}, err
	}
	taskRelations, err := taskRelationsRepresentation.ToTaskRelations()
	if err != nil {
		return TaskRelations{}, err
	}
	return taskRelations, nil
}

func (v *Task) IsStarted() bool {
	return v.State == TaskStateStarted
}

type TaskSchedule struct {
	Minutes int
	Seconds int
	Hours   int
	Cron    string
}

func ParseTaskSchedule(schedule string) (*TaskSchedule, error) {
	upperSchedule := strings.ToUpper(schedule)

	// Handle cron schedules - preserve original casing for timezone (e.g., America/Los_Angeles)
	if strings.HasPrefix(upperSchedule, "USING CRON ") {
		cron := schedule[len("USING CRON "):]
		return &TaskSchedule{Cron: cron}, nil
	}

	parts := strings.Split(strings.TrimSpace(upperSchedule), " ")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid schedule format: %s", schedule)
	}

	unit := parts[1]
	value, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	switch {
	case slices.Contains([]string{"HOURS", "HOUR", "H"}, unit):
		return &TaskSchedule{Hours: value}, nil
	case slices.Contains([]string{"MINUTES", "MINUTE", "M"}, unit):
		return &TaskSchedule{Minutes: value}, nil
	case slices.Contains([]string{"SECONDS", "SECOND", "S"}, unit):
		return &TaskSchedule{Seconds: value}, nil
	default:
		return nil, fmt.Errorf("invalid schedule format: %s", schedule)
	}
}

type CreateTaskWarehouseRequest struct {
	Warehouse                           *AccountObjectIdentifier
	UserTaskManagedInitialWarehouseSize *WarehouseSize
}

func NewCreateTaskWarehouseRequest() *CreateTaskWarehouseRequest {
	return &CreateTaskWarehouseRequest{}
}

func (s *CreateTaskWarehouseRequest) WithWarehouse(warehouse AccountObjectIdentifier) *CreateTaskWarehouseRequest {
	s.Warehouse = &warehouse
	return s
}

func (s *CreateTaskWarehouseRequest) WithUserTaskManagedInitialWarehouseSize(userTaskManagedInitialWarehouseSize WarehouseSize) *CreateTaskWarehouseRequest {
	s.UserTaskManagedInitialWarehouseSize = &userTaskManagedInitialWarehouseSize
	return s
}

func (r *CreateTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (r *CreateOrAlterTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (r *AlterTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (v *tasks) ShowParameters(ctx context.Context, id SchemaObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Task: id,
		},
	})
}

// TODO(SNOW-1277135): See if depId is necessary or could be removed
func (v *tasks) SuspendRootTasks(ctx context.Context, taskId SchemaObjectIdentifier, id SchemaObjectIdentifier) ([]SchemaObjectIdentifier, error) {
	rootTasks, err := GetRootTasks(v.client.Tasks, ctx, taskId)
	if err != nil {
		return nil, err
	}

	tasksToResume := make([]SchemaObjectIdentifier, 0)
	suspendErrs := make([]error, 0)

	for _, rootTask := range rootTasks {
		// If a root task is started, then it needs to be suspended before the child tasks can be created
		if rootTask.IsStarted() {
			err := v.client.Tasks.Alter(ctx, NewAlterTaskRequest(rootTask.ID()).WithSuspend(true))
			if err != nil {
				log.Printf("[WARN] failed to suspend task %s", rootTask.ID().FullyQualifiedName())
				suspendErrs = append(suspendErrs, err)
			}

			// Resume the task after modifications are complete as long as it is not a standalone task
			// TODO(SNOW-1277135): Document the purpose of this check and why we need different value for GetRootTasks (depId).
			if rootTask.Name != id.Name() {
				tasksToResume = append(tasksToResume, rootTask.ID())
			}
		}
	}

	return tasksToResume, errors.Join(suspendErrs...)
}

func (v *tasks) ResumeTasks(ctx context.Context, ids []SchemaObjectIdentifier) error {
	resumeErrs := make([]error, 0)
	for _, id := range ids {
		err := v.client.Tasks.Alter(ctx, NewAlterTaskRequest(id).WithResume(true))
		if err != nil {
			log.Printf("[WARN] failed to resume task %s", id.FullyQualifiedName())
			resumeErrs = append(resumeErrs, err)
		}
	}
	return errors.Join(resumeErrs...)
}

// GetRootTasks is a way to get all root tasks for the given tasks.
// Snowflake does not have (yet) a method to do it without traversing the task graph manually.
// Task DAG should have a single root but this is checked when the root task is being resumed; that's why we return here multiple roots.
// Cycles should not be possible in a task DAG, but it is checked when the root task is being resumed; that's why this method has to be cycle-proof.
func GetRootTasks(v Tasks, ctx context.Context, id SchemaObjectIdentifier) ([]Task, error) {
	tasksToExamine := collections.NewQueue[SchemaObjectIdentifier]()
	alreadyExaminedTasksNames := make([]string, 0)
	rootTasks := make([]Task, 0)

	tasksToExamine.Push(id)

	for tasksToExamine.Head() != nil {
		current := tasksToExamine.Pop()

		if slices.Contains(alreadyExaminedTasksNames, current.Name()) {
			continue
		}

		task, err := v.ShowByID(ctx, *current)
		if err != nil {
			return nil, err
		}

		if task.TaskRelations.FinalizedRootTask != nil {
			tasksToExamine.Push(*task.TaskRelations.FinalizedRootTask)
			alreadyExaminedTasksNames = append(alreadyExaminedTasksNames, current.Name())
			continue
		}

		predecessors := task.Predecessors
		if len(predecessors) == 0 {
			rootTasks = append(rootTasks, *task)
		} else {
			for _, p := range predecessors {
				tasksToExamine.Push(p)
			}
		}
		alreadyExaminedTasksNames = append(alreadyExaminedTasksNames, current.Name())
	}

	return rootTasks, nil
}

type TaskTargetCompletionInterval struct {
	Hours   *int
	Minutes *int
	Seconds *int
}

func parseTargetCompletionInterval(interval string) (*TaskTargetCompletionInterval, error) {
	upperInterval := strings.ToUpper(interval)
	parts := strings.Split(strings.TrimSpace(upperInterval), " ")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid task target completion interval format: %s", interval)
	}

	value, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid task target completion interval value: %s", interval)
	}

	unit := parts[1]
	switch {
	case slices.Contains([]string{"HOURS", "HOUR", "H"}, unit):
		return &TaskTargetCompletionInterval{Hours: &value}, nil
	case slices.Contains([]string{"MINUTES", "MINUTE", "M"}, unit):
		return &TaskTargetCompletionInterval{Minutes: &value}, nil
	case slices.Contains([]string{"SECONDS", "SECOND", "S"}, unit):
		return &TaskTargetCompletionInterval{Seconds: &value}, nil
	default:
		return nil, fmt.Errorf("invalid task target completion interval unit: %s", unit)
	}
}
