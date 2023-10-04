// Code generated by dto builder generator; DO NOT EDIT.

package sdk

import ()

func NewCreateTaskRequest(
	name SchemaObjectIdentifier,
	sql string,
) *CreateTaskRequest {
	s := CreateTaskRequest{}
	s.name = name
	s.sql = sql
	return &s
}

func (s *CreateTaskRequest) WithOrReplace(OrReplace *bool) *CreateTaskRequest {
	s.OrReplace = OrReplace
	return s
}

func (s *CreateTaskRequest) WithIfNotExists(IfNotExists *bool) *CreateTaskRequest {
	s.IfNotExists = IfNotExists
	return s
}

func (s *CreateTaskRequest) WithWarehouse(Warehouse *CreateTaskWarehouseRequest) *CreateTaskRequest {
	s.Warehouse = Warehouse
	return s
}

func (s *CreateTaskRequest) WithSchedule(Schedule *string) *CreateTaskRequest {
	s.Schedule = Schedule
	return s
}

func (s *CreateTaskRequest) WithConfig(Config *string) *CreateTaskRequest {
	s.Config = Config
	return s
}

func (s *CreateTaskRequest) WithAllowOverlappingExecution(AllowOverlappingExecution *bool) *CreateTaskRequest {
	s.AllowOverlappingExecution = AllowOverlappingExecution
	return s
}

func (s *CreateTaskRequest) WithSessionParameters(SessionParameters *SessionParameters) *CreateTaskRequest {
	s.SessionParameters = SessionParameters
	return s
}

func (s *CreateTaskRequest) WithUserTaskTimeoutMs(UserTaskTimeoutMs *int) *CreateTaskRequest {
	s.UserTaskTimeoutMs = UserTaskTimeoutMs
	return s
}

func (s *CreateTaskRequest) WithSuspendTaskAfterNumFailures(SuspendTaskAfterNumFailures *int) *CreateTaskRequest {
	s.SuspendTaskAfterNumFailures = SuspendTaskAfterNumFailures
	return s
}

func (s *CreateTaskRequest) WithErrorIntegration(ErrorIntegration *string) *CreateTaskRequest {
	s.ErrorIntegration = ErrorIntegration
	return s
}

func (s *CreateTaskRequest) WithCopyGrants(CopyGrants *bool) *CreateTaskRequest {
	s.CopyGrants = CopyGrants
	return s
}

func (s *CreateTaskRequest) WithComment(Comment *string) *CreateTaskRequest {
	s.Comment = Comment
	return s
}

func (s *CreateTaskRequest) WithAfter(After []SchemaObjectIdentifier) *CreateTaskRequest {
	s.After = After
	return s
}

func (s *CreateTaskRequest) WithTag(Tag []TagAssociation) *CreateTaskRequest {
	s.Tag = Tag
	return s
}

func (s *CreateTaskRequest) WithWhen(When *string) *CreateTaskRequest {
	s.When = When
	return s
}

func NewCreateTaskWarehouseRequest() *CreateTaskWarehouseRequest {
	return &CreateTaskWarehouseRequest{}
}

func (s *CreateTaskWarehouseRequest) WithWarehouse(Warehouse *AccountObjectIdentifier) *CreateTaskWarehouseRequest {
	s.Warehouse = Warehouse
	return s
}

func (s *CreateTaskWarehouseRequest) WithUserTaskManagedInitialWarehouseSize(UserTaskManagedInitialWarehouseSize *string) *CreateTaskWarehouseRequest {
	s.UserTaskManagedInitialWarehouseSize = UserTaskManagedInitialWarehouseSize
	return s
}

func NewAlterTaskRequest(
	name SchemaObjectIdentifier,
) *AlterTaskRequest {
	s := AlterTaskRequest{}
	s.name = name
	return &s
}

func (s *AlterTaskRequest) WithIfExists(IfExists *bool) *AlterTaskRequest {
	s.IfExists = IfExists
	return s
}

func (s *AlterTaskRequest) WithResume(Resume *bool) *AlterTaskRequest {
	s.Resume = Resume
	return s
}

func (s *AlterTaskRequest) WithSuspend(Suspend *bool) *AlterTaskRequest {
	s.Suspend = Suspend
	return s
}

func (s *AlterTaskRequest) WithRemoveAfter(RemoveAfter []SchemaObjectIdentifier) *AlterTaskRequest {
	s.RemoveAfter = RemoveAfter
	return s
}

func (s *AlterTaskRequest) WithAddAfter(AddAfter []SchemaObjectIdentifier) *AlterTaskRequest {
	s.AddAfter = AddAfter
	return s
}

func (s *AlterTaskRequest) WithSet(Set *TaskSetRequest) *AlterTaskRequest {
	s.Set = Set
	return s
}

func (s *AlterTaskRequest) WithUnset(Unset *TaskUnsetRequest) *AlterTaskRequest {
	s.Unset = Unset
	return s
}

func (s *AlterTaskRequest) WithSetTags(SetTags []TagAssociation) *AlterTaskRequest {
	s.SetTags = SetTags
	return s
}

func (s *AlterTaskRequest) WithUnsetTags(UnsetTags []ObjectIdentifier) *AlterTaskRequest {
	s.UnsetTags = UnsetTags
	return s
}

func (s *AlterTaskRequest) WithModifyAs(ModifyAs *string) *AlterTaskRequest {
	s.ModifyAs = ModifyAs
	return s
}

func (s *AlterTaskRequest) WithModifyWhen(ModifyWhen *string) *AlterTaskRequest {
	s.ModifyWhen = ModifyWhen
	return s
}

func NewTaskSetRequest() *TaskSetRequest {
	return &TaskSetRequest{}
}

func (s *TaskSetRequest) WithWarehouse(Warehouse *AccountObjectIdentifier) *TaskSetRequest {
	s.Warehouse = Warehouse
	return s
}

func (s *TaskSetRequest) WithSchedule(Schedule *string) *TaskSetRequest {
	s.Schedule = Schedule
	return s
}

func (s *TaskSetRequest) WithConfig(Config *string) *TaskSetRequest {
	s.Config = Config
	return s
}

func (s *TaskSetRequest) WithAllowOverlappingExecution(AllowOverlappingExecution *bool) *TaskSetRequest {
	s.AllowOverlappingExecution = AllowOverlappingExecution
	return s
}

func (s *TaskSetRequest) WithUserTaskTimeoutMs(UserTaskTimeoutMs *int) *TaskSetRequest {
	s.UserTaskTimeoutMs = UserTaskTimeoutMs
	return s
}

func (s *TaskSetRequest) WithSuspendTaskAfterNumFailures(SuspendTaskAfterNumFailures *int) *TaskSetRequest {
	s.SuspendTaskAfterNumFailures = SuspendTaskAfterNumFailures
	return s
}

func (s *TaskSetRequest) WithComment(Comment *string) *TaskSetRequest {
	s.Comment = Comment
	return s
}

func (s *TaskSetRequest) WithSessionParameters(SessionParameters *SessionParameters) *TaskSetRequest {
	s.SessionParameters = SessionParameters
	return s
}

func NewTaskUnsetRequest() *TaskUnsetRequest {
	return &TaskUnsetRequest{}
}

func (s *TaskUnsetRequest) WithWarehouse(Warehouse *bool) *TaskUnsetRequest {
	s.Warehouse = Warehouse
	return s
}

func (s *TaskUnsetRequest) WithSchedule(Schedule *bool) *TaskUnsetRequest {
	s.Schedule = Schedule
	return s
}

func (s *TaskUnsetRequest) WithConfig(Config *bool) *TaskUnsetRequest {
	s.Config = Config
	return s
}

func (s *TaskUnsetRequest) WithAllowOverlappingExecution(AllowOverlappingExecution *bool) *TaskUnsetRequest {
	s.AllowOverlappingExecution = AllowOverlappingExecution
	return s
}

func (s *TaskUnsetRequest) WithUserTaskTimeoutMs(UserTaskTimeoutMs *bool) *TaskUnsetRequest {
	s.UserTaskTimeoutMs = UserTaskTimeoutMs
	return s
}

func (s *TaskUnsetRequest) WithSuspendTaskAfterNumFailures(SuspendTaskAfterNumFailures *bool) *TaskUnsetRequest {
	s.SuspendTaskAfterNumFailures = SuspendTaskAfterNumFailures
	return s
}

func (s *TaskUnsetRequest) WithComment(Comment *bool) *TaskUnsetRequest {
	s.Comment = Comment
	return s
}

func (s *TaskUnsetRequest) WithSessionParametersUnset(SessionParametersUnset *SessionParametersUnset) *TaskUnsetRequest {
	s.SessionParametersUnset = SessionParametersUnset
	return s
}

func NewDropTaskRequest(
	name SchemaObjectIdentifier,
) *DropTaskRequest {
	s := DropTaskRequest{}
	s.name = name
	return &s
}

func (s *DropTaskRequest) WithIfExists(IfExists *bool) *DropTaskRequest {
	s.IfExists = IfExists
	return s
}

func NewShowTaskRequest() *ShowTaskRequest {
	return &ShowTaskRequest{}
}

func (s *ShowTaskRequest) WithTerse(Terse *bool) *ShowTaskRequest {
	s.Terse = Terse
	return s
}

func NewDescribeTaskRequest(
	name SchemaObjectIdentifier,
) *DescribeTaskRequest {
	s := DescribeTaskRequest{}
	s.name = name
	return &s
}

func NewExecuteTaskRequest(
	name SchemaObjectIdentifier,
) *ExecuteTaskRequest {
	s := ExecuteTaskRequest{}
	s.name = name
	return &s
}

func (s *ExecuteTaskRequest) WithRetryLast(RetryLast *bool) *ExecuteTaskRequest {
	s.RetryLast = RetryLast
	return s
}
