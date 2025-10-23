package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateTaskOptions]        = new(CreateTaskRequest)
	_ optionsProvider[CreateOrAlterTaskOptions] = new(CreateOrAlterTaskRequest)
	_ optionsProvider[CloneTaskOptions]         = new(CloneTaskRequest)
	_ optionsProvider[AlterTaskOptions]         = new(AlterTaskRequest)
	_ optionsProvider[DropTaskOptions]          = new(DropTaskRequest)
	_ optionsProvider[ShowTaskOptions]          = new(ShowTaskRequest)
	_ optionsProvider[DescribeTaskOptions]      = new(DescribeTaskRequest)
	_ optionsProvider[ExecuteTaskOptions]       = new(ExecuteTaskRequest)
)

type CreateTaskRequest struct {
	OrReplace                               *bool
	IfNotExists                             *bool
	name                                    SchemaObjectIdentifier // required
	Warehouse                               *CreateTaskWarehouseRequest
	Schedule                                *string
	Config                                  *string
	AllowOverlappingExecution               *bool
	SessionParameters                       *SessionParameters
	UserTaskTimeoutMs                       *int
	SuspendTaskAfterNumFailures             *int
	ErrorIntegration                        *AccountObjectIdentifier
	Comment                                 *string
	Finalize                                *SchemaObjectIdentifier
	TaskAutoRetryAttempts                   *int
	Tag                                     []TagAssociation
	UserTaskMinimumTriggerIntervalInSeconds *int
	TargetCompletionInterval                *string
	ServerlessTaskMinStatementSize          *WarehouseSize
	ServerlessTaskMaxStatementSize          *WarehouseSize
	After                                   []SchemaObjectIdentifier
	When                                    *string
	sql                                     string // required
}

type CreateTaskWarehouseRequest struct {
	Warehouse                           *AccountObjectIdentifier
	UserTaskManagedInitialWarehouseSize *WarehouseSize
}

func (r *CreateTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type CreateOrAlterTaskRequest struct {
	name                        SchemaObjectIdentifier // required
	Warehouse                   *CreateTaskWarehouseRequest
	Schedule                    *string
	Config                      *string
	AllowOverlappingExecution   *bool
	UserTaskTimeoutMs           *int
	SessionParameters           *SessionParameters
	SuspendTaskAfterNumFailures *int
	ErrorIntegration            *AccountObjectIdentifier
	Comment                     *string
	Finalize                    *SchemaObjectIdentifier
	TaskAutoRetryAttempts       *int
	After                       []SchemaObjectIdentifier
	When                        *string
	sql                         string // required
}

func (r *CreateOrAlterTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type CloneTaskRequest struct {
	OrReplace  *bool
	name       SchemaObjectIdentifier // required
	sourceTask SchemaObjectIdentifier // required
	CopyGrants *bool
}

type AlterTaskRequest struct {
	IfExists      *bool
	name          SchemaObjectIdentifier // required
	Resume        *bool
	Suspend       *bool
	RemoveAfter   []SchemaObjectIdentifier
	AddAfter      []SchemaObjectIdentifier
	Set           *TaskSetRequest
	Unset         *TaskUnsetRequest
	SetTags       []TagAssociation
	UnsetTags     []ObjectIdentifier
	SetFinalize   *SchemaObjectIdentifier
	UnsetFinalize *bool
	ModifyAs      *string
	ModifyWhen    *string
	RemoveWhen    *bool
}

func (r *AlterTaskRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type TaskSetRequest struct {
	Warehouse                               *AccountObjectIdentifier
	UserTaskManagedInitialWarehouseSize     *WarehouseSize
	Schedule                                *string
	Config                                  *string
	AllowOverlappingExecution               *bool
	UserTaskTimeoutMs                       *int
	SuspendTaskAfterNumFailures             *int
	ErrorIntegration                        *AccountObjectIdentifier
	Comment                                 *string
	SessionParameters                       *SessionParameters
	TaskAutoRetryAttempts                   *int
	UserTaskMinimumTriggerIntervalInSeconds *int
	TargetCompletionInterval                *string
	ServerlessTaskMinStatementSize          *WarehouseSize
	ServerlessTaskMaxStatementSize          *WarehouseSize
}

type TaskUnsetRequest struct {
	Warehouse                               *bool
	UserTaskManagedInitialWarehouseSize     *bool
	Schedule                                *bool
	Config                                  *bool
	AllowOverlappingExecution               *bool
	UserTaskTimeoutMs                       *bool
	SuspendTaskAfterNumFailures             *bool
	ErrorIntegration                        *bool
	Comment                                 *bool
	TaskAutoRetryAttempts                   *bool
	UserTaskMinimumTriggerIntervalInSeconds *bool
	TargetCompletionInterval                *bool
	ServerlessTaskMinStatementSize          *bool
	ServerlessTaskMaxStatementSize          *bool
	SessionParametersUnset                  *SessionParametersUnset
}

type DropTaskRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type ShowTaskRequest struct {
	Terse      *bool
	Like       *Like
	In         *ExtendedIn
	StartsWith *string
	RootOnly   *bool
	Limit      *LimitFrom
}

type DescribeTaskRequest struct {
	name SchemaObjectIdentifier // required
}

type ExecuteTaskRequest struct {
	name      SchemaObjectIdentifier // required
	RetryLast *bool
}
