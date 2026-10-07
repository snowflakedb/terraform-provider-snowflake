package defs

import (
	"slices"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
)

var taskParameters = ParameterDefsForLevel(parameterdefs.ParameterLevelTask)

var taskParameterFieldNames = collections.Map(taskParameters, g.ParameterSqlToFieldName)

// CREATE TASK sets the initial warehouse size through the WAREHOUSE clause (taskCreateWarehouse), not as a standalone assignment.
var taskCreateParameters = collections.Filter(taskParameters, func(p parameterdefs.ParameterDef) bool {
	return p.SqlName != UserTaskManagedInitialWarehouseSize.SqlName
})

var taskPairs = g.StructPair("taskDBRow", "Task").
	Text("created_on").
	Text("name").
	Text("id").
	Text("database_name").
	Text("schema_name").
	Text("owner").
	OptionalText("comment", g.WithRequiredInPlain()).
	Field("warehouse", "sql.NullString", "*AccountObjectIdentifier", g.WithPlainFieldName("Warehouse"), g.WithManualConvert()).
	OptionalText("schedule", g.WithRequiredInPlain()).
	Field("predecessors", "string", "[]SchemaObjectIdentifier", g.WithManualConvert()).
	PlainField("state", "TaskState", g.WithManualConvert()).
	Text("definition").
	OptionalText("condition", g.WithRequiredInPlain()).
	Field("allow_overlapping_execution", "string", "bool", g.WithBoolTrueValue("true")).
	Field("error_integration", "sql.NullString", "*AccountObjectIdentifier", g.WithPlainFieldName("ErrorIntegration"), g.WithManualConvert()).
	OptionalText("last_committed_on", g.WithRequiredInPlain()).
	OptionalText("last_suspended_on", g.WithRequiredInPlain()).
	Text("owner_role_type").
	OptionalText("config", g.WithRequiredInPlain()).
	OptionalText("budget", g.WithRequiredInPlain()).
	PlainField("task_relations", "TaskRelations", g.WithManualConvert()).
	OptionalText("last_suspended_reason", g.WithRequiredInPlain()).
	Field("target_completion_interval", "sql.NullString", "*TaskTargetCompletionInterval", g.WithPlainFieldName("TargetCompletionInterval"), g.WithManualConvert()).
	Field("execute_as_user", "sql.NullString", "*AccountObjectIdentifier")

var taskCreateWarehouse = g.NewQueryStruct("CreateTaskWarehouse").
	OptionalIdentifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WAREHOUSE")).
	OptionalAssignment("USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE", "WarehouseSize", g.ParameterOptions().SingleQuotes()).
	WithValidation(g.ExactlyOneValueSet, "Warehouse", "UserTaskManagedInitialWarehouseSize")

var tasksDef = g.NewInterface(
	"Tasks",
	"Task",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-task",
		g.NewQueryStruct("CreateTask").
			Create().
			OrReplace().
			SQL("TASK").
			IfNotExists().
			Name().
			OptionalQueryStructField("Warehouse", taskCreateWarehouse, g.KeywordOptions()).
			OptionalTextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("CONFIG", g.ParameterOptions().DoubleDollarQuotes()).
			OptionalBooleanAssignment("ALLOW_OVERLAPPING_EXECUTION", nil).
			OptionalIdentifier("ErrorIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("ERROR_INTEGRATION")).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			OptionalIdentifier("Finalize", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Equals().SQL("FINALIZE")).
			OptionalTags().
			OptionalTextAssignment("TARGET_COMPLETION_INTERVAL", g.ParameterOptions().SingleQuotes()).
			WithParameters(taskCreateParameters...).
			ListAssignment("AFTER", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().NoEquals()).
			OptionalIdentifier("ExecuteAsUser", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXECUTE AS USER")).
			OptionalTextAssignment("WHEN", g.ParameterOptions().NoQuotes().NoEquals()).
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifierIfSet, "ErrorIntegration").
			WithValidation(g.ValidIdentifierIfSet, "ExecuteAsUser").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
			WithValidation(g.NoDoubleDollarQuotesIfSet, "Config").
			WithValidation(g.ValidIdentifier, "After"),
	).
	CustomOperation(
		"CreateOrAlter",
		"https://docs.snowflake.com/en/sql-reference/sql/create-task#create-or-alter-task",
		g.NewQueryStruct("CloneTask").
			CreateOrAlter().
			SQL("TASK").
			Name().
			OptionalQueryStructField("Warehouse", taskCreateWarehouse, g.KeywordOptions()).
			OptionalTextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("CONFIG", g.ParameterOptions().DoubleDollarQuotes()).
			OptionalBooleanAssignment("ALLOW_OVERLAPPING_EXECUTION", nil).
			OptionalIdentifier("ErrorIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("ERROR_INTEGRATION")).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			OptionalIdentifier("Finalize", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Equals().SQL("FINALIZE")).
			WithParameters(taskCreateParameters...).
			ListAssignment("AFTER", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().NoEquals()).
			OptionalIdentifier("ExecuteAsUser", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXECUTE AS USER")).
			OptionalTextAssignment("WHEN", g.ParameterOptions().NoQuotes().NoEquals()).
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifierIfSet, "ErrorIntegration").
			WithValidation(g.ValidIdentifierIfSet, "ExecuteAsUser").
			WithValidation(g.NoDoubleDollarQuotesIfSet, "Config").
			WithValidation(g.ValidIdentifier, "After"),
	).
	CustomOperation(
		"Clone",
		"https://docs.snowflake.com/en/sql-reference/sql/create-task#create-task-clone",
		g.NewQueryStruct("CloneTask").
			Create().
			OrReplace().
			SQL("TASK").
			Name().
			SQL("CLONE").
			Identifier("sourceTask", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
			OptionalSQL("COPY GRANTS").
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifier, "sourceTask"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-task",
		g.NewQueryStruct("AlterTask").
			Alter().
			SQL("TASK").
			IfExists().
			Name().
			OptionalSQL("RESUME").
			OptionalSQL("SUSPEND").
			ListAssignment("REMOVE AFTER", "SchemaObjectIdentifier", g.ParameterOptions().NoEquals()).
			ListAssignment("ADD AFTER", "SchemaObjectIdentifier", g.ParameterOptions().NoEquals()).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("TaskSet").
					OptionalIdentifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WAREHOUSE")).
					OptionalTextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("CONFIG", g.ParameterOptions().DoubleDollarQuotes()).
					OptionalBooleanAssignment("ALLOW_OVERLAPPING_EXECUTION", nil).
					OptionalIdentifier("ErrorIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("ERROR_INTEGRATION")).
					OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("TARGET_COMPLETION_INTERVAL", g.ParameterOptions().SingleQuotes()).
					WithParameters(taskParameters...).
					WithValidation(g.AtLeastOneValueSet, append(slices.Clone(taskParameterFieldNames), "Warehouse", "Schedule", "Config", "AllowOverlappingExecution", "ErrorIntegration", "Comment", "TargetCompletionInterval")...).
					WithValidation(g.ConflictingFields, "Warehouse", "UserTaskManagedInitialWarehouseSize").
					WithValidation(g.ValidIdentifierIfSet, "ErrorIntegration").
					WithValidation(g.NoDoubleDollarQuotesIfSet, "Config"),
				g.ListOptions().SQL("SET").NoParentheses(),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("TaskUnset").
					OptionalSQL("WAREHOUSE").
					OptionalSQL("SCHEDULE").
					OptionalSQL("CONFIG").
					OptionalSQL("ALLOW_OVERLAPPING_EXECUTION").
					OptionalSQL("ERROR_INTEGRATION").
					OptionalSQL("COMMENT").
					OptionalSQL("TARGET_COMPLETION_INTERVAL").
					WithParametersUnset(taskParameters...).
					WithValidation(g.AtLeastOneValueSet, append(slices.Clone(taskParameterFieldNames), "Warehouse", "Schedule", "Config", "AllowOverlappingExecution", "ErrorIntegration", "Comment", "TargetCompletionInterval")...),
				g.ListOptions().SQL("UNSET").NoParentheses(),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalIdentifier("SetFinalize", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Equals().SQL("SET FINALIZE")).
			OptionalSQL("UNSET FINALIZE").
			OptionalIdentifier("SetExecuteAsUser", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("SET EXECUTE AS USER")).
			OptionalSQL("UNSET EXECUTE AS USER").
			OptionalTextAssignment("MODIFY AS", g.ParameterOptions().NoQuotes().NoEquals()).
			OptionalTextAssignment("MODIFY WHEN", g.ParameterOptions().NoQuotes().NoEquals()).
			OptionalSQL("REMOVE WHEN").
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifierIfSet, "SetExecuteAsUser").
			WithValidation(g.ExactlyOneValueSet, "Resume", "Suspend", "RemoveAfter", "AddAfter", "Set", "Unset", "SetTags", "UnsetTags", "SetFinalize", "UnsetFinalize", "SetExecuteAsUser", "UnsetExecuteAsUser", "ModifyAs", "ModifyWhen", "RemoveWhen").
			WithValidation(g.ValidIdentifier, "RemoveAfter").
			WithValidation(g.ValidIdentifier, "AddAfter"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-task",
		g.NewQueryStruct("DropTask").
			Drop().
			SQL("TASK").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-tasks",
		taskPairs,
		g.NewQueryStruct("ShowTasks").
			Show().
			Terse().
			SQL("TASKS").
			OptionalLike().
			OptionalExtendedIn().
			OptionalStartsWith().
			OptionalSQL("ROOT ONLY").
			OptionalLimit(),
		g.ShowByIDExtendedInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-task",
		taskPairs,
		g.NewQueryStruct("DescribeTask").
			Describe().
			SQL("TASK").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	CustomOperation(
		"Execute",
		"https://docs.snowflake.com/en/sql-reference/sql/execute-task",
		g.NewQueryStruct("ExecuteTask").
			SQL("EXECUTE").
			SQL("TASK").
			Name().
			OptionalSQL("RETRY LAST").
			WithValidation(g.ValidIdentifier, "name"),
	).
	WithCustomInterfaceMethod(
		"ShowParameters",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"[]*Parameter", "error",
	).
	ShowParametersDetails(taskParameters...).
	WithCustomInterfaceMethod(
		"SuspendRootTasks",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("taskId", g.KindOfT[sdkcommons.SchemaObjectIdentifier]()),
			g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]()),
		},
		"[]SchemaObjectIdentifier", "error",
	).
	WithCustomInterfaceMethod(
		"ResumeTasks",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("ids", "[]SchemaObjectIdentifier")},
		"error",
	)
