package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var setSpendingLimitArgs = g.NewQueryStruct("SetSpendingLimitArgs").
	PredefinedQueryStructField("SpendingLimit", "int", g.ParameterOptions().Required().NoEquals())

var budgetEmail = g.NewQueryStruct("BudgetEmail").
	Text("Email", g.KeywordOptions().SingleQuotes().Required())

var setEmailNotificationsArgs = g.NewQueryStruct("SetEmailNotificationsArgs").
	OptionalIdentifier("NotificationIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SingleQuotes()).
	ListQueryStructField("Emails", budgetEmail, g.ListOptions().Required())

var getNotificationIntegrationsResult = g.StructPair(
	"getNotificationIntegrationsRow",
	"BudgetNotificationIntegration",
).Text("integration_name").
	Number("last_notification_time").
	Time("added_date").
	WithConvertGeneration()

// stored_procedure_reference and array_construct_statement are omitted as they are duplicated ways of setting the same input
var setCycleStartActionArgs = g.NewQueryStruct("SetCycleStartActionArgs").
	Identifier("Procedure", g.KindOfT[sdkcommons.SchemaObjectIdentifierWithArguments](), g.IdentifierOptions().Required().SystemReference("PROCEDURE")).
	List("Arguments", "string", g.ListOptions().Required())

var getCycleStartActionResult = g.StructPair(
	"getCycleStartActionRow",
	"BudgetCycleStartAction",
).
	// uppercase column names used on purpose
	Text("ACTION_UUID").
	SchemaObjectIdentifierWithArguments("PROCEDURE_FQN", g.WithPlainFieldName("ProcedureId")).
	StringList("PROCEDURE_ARGS").
	Time("ADDED_TIMESTAMP").
	Time("LAST_TRIGGERED_TIMESTAMP").
	WithConvertGeneration()

var budgetsDef = g.NewInterface(
	"Budgets",
	"Budget",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/commands/create-budget",
	g.NewQueryStruct("CreateBudgetOptions").
		Create().
		OrReplace().
		SQLWithCustomFieldName("snowflakeCoreBudget", "SNOWFLAKE.CORE.BUDGET").
		IfNotExists().
		Name().
		SQLWithCustomFieldName("parens", "()").
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/commands/drop-budget",
	g.NewQueryStruct("DropBudgetOptions").
		Drop().
		SQLWithCustomFieldName("snowflakeCoreBudget", "SNOWFLAKE.CORE.BUDGET").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).InstanceMethodOperationScalar(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/method/set_spending_limit",
	"SET_SPENDING_LIMIT",
	setSpendingLimitArgs,
	"string",
).InstanceMethodOperationScalar(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/method/get_spending_limit",
	"GET_SPENDING_LIMIT",
	nil,
	"int",
).InstanceMethodOperationScalar(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/method/set_email_notifications",
	"SET_EMAIL_NOTIFICATIONS",
	setEmailNotificationsArgs,
	"string",
).InstanceMethodOperation(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/method/get_notification_integrations",
	"GET_NOTIFICATION_INTEGRATIONS",
	nil,
	getNotificationIntegrationsResult,
	g.InstanceMethodKindSlice,
).InstanceMethodOperationScalar(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/methods/get_notification_email",
	"GET_NOTIFICATION_EMAIL",
	nil,
	"string",
).InstanceMethodOperationScalar(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/methods/get_notification_integration_name",
	"GET_NOTIFICATION_INTEGRATION_NAME",
	nil,
	"string",
).InstanceMethodOperationScalar(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/methods/set_cycle_start_action",
	"SET_CYCLE_START_ACTION",
	setCycleStartActionArgs,
	"string",
).InstanceMethodOperation(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/methods/get_cycle_start_action",
	"GET_CYCLE_START_ACTION",
	nil,
	getCycleStartActionResult,
	g.InstanceMethodKindSingleValue,
)
