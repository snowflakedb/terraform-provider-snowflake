package sdk

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Tasks interface {
	Create(ctx context.Context, request *CreateTaskRequest) error
	CreateOrAlter(ctx context.Context, request *CreateOrAlterTaskRequest) error
	Clone(ctx context.Context, request *CloneTaskRequest) error
	Alter(ctx context.Context, request *AlterTaskRequest) error
	Drop(ctx context.Context, request *DropTaskRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, request *ShowTaskRequest) ([]Task, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Task, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Task, error)
	ShowParameters(ctx context.Context, id SchemaObjectIdentifier) ([]*Parameter, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*Task, error)
	Execute(ctx context.Context, request *ExecuteTaskRequest) error
	SuspendRootTasks(ctx context.Context, taskId SchemaObjectIdentifier, id SchemaObjectIdentifier) ([]SchemaObjectIdentifier, error)
	ResumeTasks(ctx context.Context, ids []SchemaObjectIdentifier) error
}

// CreateTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-task.
type CreateTaskOptions struct {
	create                                  bool                     `ddl:"static" sql:"CREATE"`
	OrReplace                               *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	task                                    bool                     `ddl:"static" sql:"TASK"`
	IfNotExists                             *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                                    SchemaObjectIdentifier   `ddl:"identifier"`
	Warehouse                               *CreateTaskWarehouse     `ddl:"keyword"`
	Schedule                                *string                  `ddl:"parameter,single_quotes" sql:"SCHEDULE"`
	Config                                  *string                  `ddl:"parameter,no_quotes" sql:"CONFIG"`
	AllowOverlappingExecution               *bool                    `ddl:"parameter" sql:"ALLOW_OVERLAPPING_EXECUTION"`
	SessionParameters                       *SessionParameters       `ddl:"list,no_parentheses"`
	UserTaskTimeoutMs                       *int                     `ddl:"parameter" sql:"USER_TASK_TIMEOUT_MS"`
	SuspendTaskAfterNumFailures             *int                     `ddl:"parameter" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	ErrorIntegration                        *AccountObjectIdentifier `ddl:"identifier,equals" sql:"ERROR_INTEGRATION"`
	Comment                                 *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Finalize                                *SchemaObjectIdentifier  `ddl:"identifier,equals" sql:"FINALIZE"`
	TaskAutoRetryAttempts                   *int                     `ddl:"parameter" sql:"TASK_AUTO_RETRY_ATTEMPTS"`
	Tag                                     []TagAssociation         `ddl:"keyword,parentheses" sql:"TAG"`
	UserTaskMinimumTriggerIntervalInSeconds *int                     `ddl:"parameter" sql:"USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"`
	TargetCompletionInterval                *string                  `ddl:"parameter,single_quotes" sql:"TARGET_COMPLETION_INTERVAL"`
	ServerlessTaskMinStatementSize          *WarehouseSize           `ddl:"parameter,single_quotes" sql:"SERVERLESS_TASK_MIN_STATEMENT_SIZE"`
	ServerlessTaskMaxStatementSize          *WarehouseSize           `ddl:"parameter,single_quotes" sql:"SERVERLESS_TASK_MAX_STATEMENT_SIZE"`
	After                                   []SchemaObjectIdentifier `ddl:"parameter,no_equals" sql:"AFTER"`
	When                                    *string                  `ddl:"parameter,no_quotes,no_equals" sql:"WHEN"`
	as                                      bool                     `ddl:"static" sql:"AS"`
	sql                                     string                   `ddl:"keyword,no_quotes"`
}

type CreateTaskWarehouse struct {
	Warehouse                           *AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
	UserTaskManagedInitialWarehouseSize *WarehouseSize           `ddl:"parameter,single_quotes" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
}

// CreateOrAlterTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-task#create-or-alter-task.
type CreateOrAlterTaskOptions struct {
	createOrAlter               bool                     `ddl:"static" sql:"CREATE OR ALTER"`
	task                        bool                     `ddl:"static" sql:"TASK"`
	name                        SchemaObjectIdentifier   `ddl:"identifier"`
	Warehouse                   *CreateTaskWarehouse     `ddl:"keyword"`
	Schedule                    *string                  `ddl:"parameter,single_quotes" sql:"SCHEDULE"`
	Config                      *string                  `ddl:"parameter,no_quotes" sql:"CONFIG"`
	AllowOverlappingExecution   *bool                    `ddl:"parameter" sql:"ALLOW_OVERLAPPING_EXECUTION"`
	UserTaskTimeoutMs           *int                     `ddl:"parameter" sql:"USER_TASK_TIMEOUT_MS"`
	SessionParameters           *SessionParameters       `ddl:"list,no_parentheses"`
	SuspendTaskAfterNumFailures *int                     `ddl:"parameter" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	ErrorIntegration            *AccountObjectIdentifier `ddl:"identifier,equals" sql:"ERROR_INTEGRATION"`
	Comment                     *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Finalize                    *SchemaObjectIdentifier  `ddl:"identifier,equals" sql:"FINALIZE"`
	TaskAutoRetryAttempts       *int                     `ddl:"parameter" sql:"TASK_AUTO_RETRY_ATTEMPTS"`
	After                       []SchemaObjectIdentifier `ddl:"parameter,no_equals" sql:"AFTER"`
	When                        *string                  `ddl:"parameter,no_quotes,no_equals" sql:"WHEN"`
	as                          bool                     `ddl:"static" sql:"AS"`
	sql                         string                   `ddl:"keyword,no_quotes"`
}

// CloneTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-task#create-task-clone.
type CloneTaskOptions struct {
	create     bool                   `ddl:"static" sql:"CREATE"`
	OrReplace  *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	task       bool                   `ddl:"static" sql:"TASK"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
	clone      bool                   `ddl:"static" sql:"CLONE"`
	sourceTask SchemaObjectIdentifier `ddl:"identifier"`
	CopyGrants *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
}

// AlterTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-task.
type AlterTaskOptions struct {
	alter         bool                     `ddl:"static" sql:"ALTER"`
	task          bool                     `ddl:"static" sql:"TASK"`
	IfExists      *bool                    `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier   `ddl:"identifier"`
	Resume        *bool                    `ddl:"keyword" sql:"RESUME"`
	Suspend       *bool                    `ddl:"keyword" sql:"SUSPEND"`
	RemoveAfter   []SchemaObjectIdentifier `ddl:"parameter,no_equals" sql:"REMOVE AFTER"`
	AddAfter      []SchemaObjectIdentifier `ddl:"parameter,no_equals" sql:"ADD AFTER"`
	Set           *TaskSet                 `ddl:"list,no_parentheses" sql:"SET"`
	Unset         *TaskUnset               `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTags       []TagAssociation         `ddl:"keyword" sql:"SET TAG"`
	UnsetTags     []ObjectIdentifier       `ddl:"keyword" sql:"UNSET TAG"`
	SetFinalize   *SchemaObjectIdentifier  `ddl:"identifier,equals" sql:"SET FINALIZE"`
	UnsetFinalize *bool                    `ddl:"keyword" sql:"UNSET FINALIZE"`
	ModifyAs      *string                  `ddl:"parameter,no_quotes,no_equals" sql:"MODIFY AS"`
	ModifyWhen    *string                  `ddl:"parameter,no_quotes,no_equals" sql:"MODIFY WHEN"`
	RemoveWhen    *bool                    `ddl:"keyword" sql:"REMOVE WHEN"`
}

type TaskSet struct {
	Warehouse                               *AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
	UserTaskManagedInitialWarehouseSize     *WarehouseSize           `ddl:"parameter,single_quotes" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
	Schedule                                *string                  `ddl:"parameter,single_quotes" sql:"SCHEDULE"`
	Config                                  *string                  `ddl:"parameter,no_quotes" sql:"CONFIG"`
	AllowOverlappingExecution               *bool                    `ddl:"parameter" sql:"ALLOW_OVERLAPPING_EXECUTION"`
	UserTaskTimeoutMs                       *int                     `ddl:"parameter" sql:"USER_TASK_TIMEOUT_MS"`
	SuspendTaskAfterNumFailures             *int                     `ddl:"parameter" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	ErrorIntegration                        *AccountObjectIdentifier `ddl:"identifier,equals" sql:"ERROR_INTEGRATION"`
	Comment                                 *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
	SessionParameters                       *SessionParameters       `ddl:"list,no_parentheses"`
	TaskAutoRetryAttempts                   *int                     `ddl:"parameter" sql:"TASK_AUTO_RETRY_ATTEMPTS"`
	UserTaskMinimumTriggerIntervalInSeconds *int                     `ddl:"parameter" sql:"USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"`
	TargetCompletionInterval                *string                  `ddl:"parameter,single_quotes" sql:"TARGET_COMPLETION_INTERVAL"`
	ServerlessTaskMinStatementSize          *WarehouseSize           `ddl:"parameter,single_quotes" sql:"SERVERLESS_TASK_MIN_STATEMENT_SIZE"`
	ServerlessTaskMaxStatementSize          *WarehouseSize           `ddl:"parameter,single_quotes" sql:"SERVERLESS_TASK_MAX_STATEMENT_SIZE"`
}

type TaskUnset struct {
	Warehouse                               *bool                   `ddl:"keyword" sql:"WAREHOUSE"`
	UserTaskManagedInitialWarehouseSize     *bool                   `ddl:"keyword" sql:"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE"`
	Schedule                                *bool                   `ddl:"keyword" sql:"SCHEDULE"`
	Config                                  *bool                   `ddl:"keyword" sql:"CONFIG"`
	AllowOverlappingExecution               *bool                   `ddl:"keyword" sql:"ALLOW_OVERLAPPING_EXECUTION"`
	UserTaskTimeoutMs                       *bool                   `ddl:"keyword" sql:"USER_TASK_TIMEOUT_MS"`
	SuspendTaskAfterNumFailures             *bool                   `ddl:"keyword" sql:"SUSPEND_TASK_AFTER_NUM_FAILURES"`
	ErrorIntegration                        *bool                   `ddl:"keyword" sql:"ERROR_INTEGRATION"`
	Comment                                 *bool                   `ddl:"keyword" sql:"COMMENT"`
	TaskAutoRetryAttempts                   *bool                   `ddl:"keyword" sql:"TASK_AUTO_RETRY_ATTEMPTS"`
	UserTaskMinimumTriggerIntervalInSeconds *bool                   `ddl:"keyword" sql:"USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS"`
	TargetCompletionInterval                *bool                   `ddl:"keyword" sql:"TARGET_COMPLETION_INTERVAL"`
	ServerlessTaskMinStatementSize          *bool                   `ddl:"keyword" sql:"SERVERLESS_TASK_MIN_STATEMENT_SIZE"`
	ServerlessTaskMaxStatementSize          *bool                   `ddl:"keyword" sql:"SERVERLESS_TASK_MAX_STATEMENT_SIZE"`
	SessionParametersUnset                  *SessionParametersUnset `ddl:"list,no_parentheses"`
}

// DropTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-task.
type DropTaskOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	task     bool                   `ddl:"static" sql:"TASK"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-tasks.
type ShowTaskOptions struct {
	show       bool        `ddl:"static" sql:"SHOW"`
	Terse      *bool       `ddl:"keyword" sql:"TERSE"`
	tasks      bool        `ddl:"static" sql:"TASKS"`
	Like       *Like       `ddl:"keyword" sql:"LIKE"`
	In         *ExtendedIn `ddl:"keyword" sql:"IN"`
	StartsWith *string     `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	RootOnly   *bool       `ddl:"keyword" sql:"ROOT ONLY"`
	Limit      *LimitFrom  `ddl:"keyword" sql:"LIMIT"`
}

type taskDBRow struct {
	CreatedOn                 string         `db:"created_on"`
	Name                      string         `db:"name"`
	Id                        string         `db:"id"`
	DatabaseName              string         `db:"database_name"`
	SchemaName                string         `db:"schema_name"`
	Owner                     string         `db:"owner"`
	Comment                   sql.NullString `db:"comment"`
	Warehouse                 sql.NullString `db:"warehouse"`
	Schedule                  sql.NullString `db:"schedule"`
	Predecessors              string         `db:"predecessors"`
	State                     string         `db:"state"`
	Definition                string         `db:"definition"`
	Condition                 sql.NullString `db:"condition"`
	AllowOverlappingExecution string         `db:"allow_overlapping_execution"`
	ErrorIntegration          sql.NullString `db:"error_integration"`
	LastCommittedOn           sql.NullString `db:"last_committed_on"`
	LastSuspendedOn           sql.NullString `db:"last_suspended_on"`
	OwnerRoleType             string         `db:"owner_role_type"`
	Config                    sql.NullString `db:"config"`
	Budget                    sql.NullString `db:"budget"`
	TaskRelations             string         `db:"task_relations"`
	LastSuspendedReason       sql.NullString `db:"last_suspended_reason"`
}

type Task struct {
	CreatedOn                 string
	Name                      string
	Id                        string
	DatabaseName              string
	SchemaName                string
	Owner                     string
	Comment                   string
	Warehouse                 *AccountObjectIdentifier
	Schedule                  string
	Predecessors              []SchemaObjectIdentifier
	State                     TaskState
	Definition                string
	Condition                 string
	AllowOverlappingExecution bool
	ErrorIntegration          *AccountObjectIdentifier
	LastCommittedOn           string
	LastSuspendedOn           string
	OwnerRoleType             string
	Config                    string
	Budget                    string
	TaskRelations             TaskRelations
	LastSuspendedReason       string
}

func (v *Task) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Task) IsStarted() bool {
	return v.State == TaskStateStarted
}

type TaskSchedule struct {
	Minutes int
	Cron    string
}

func ParseTaskSchedule(schedule string) (*TaskSchedule, error) {
	upperSchedule := strings.ToUpper(schedule)
	switch {
	case strings.Contains(upperSchedule, "USING CRON"):
		// We have to do it this was because we want to get rid of the prefix and leave the casing as is (mostly because timezones like America/Los_Angeles are case-sensitive).
		// That why the prefix trimming has to be done by slicing rather than using strings.TrimPrefix.
		cron := schedule[len("USING CRON "):]
		return &TaskSchedule{Cron: cron}, nil
	case strings.HasSuffix(upperSchedule, "M") ||
		strings.HasSuffix(upperSchedule, "MINUTE") ||
		strings.HasSuffix(upperSchedule, "MINUTES"):
		minuteParts := strings.Split(upperSchedule, " ")
		minutes, err := strconv.Atoi(minuteParts[0])
		if err != nil {
			return nil, err
		}

		return &TaskSchedule{Minutes: minutes}, nil
	default:
		return nil, fmt.Errorf("invalid schedule format: %s", schedule)
	}
}

// DescribeTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-task.
type DescribeTaskOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	task     bool                   `ddl:"static" sql:"TASK"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// ExecuteTaskOptions is based on https://docs.snowflake.com/en/sql-reference/sql/execute-task.
type ExecuteTaskOptions struct {
	execute   bool                   `ddl:"static" sql:"EXECUTE"`
	task      bool                   `ddl:"static" sql:"TASK"`
	name      SchemaObjectIdentifier `ddl:"identifier"`
	RetryLast *bool                  `ddl:"keyword" sql:"RETRY LAST"`
}
