//go:build sdk_generation

package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var setSpendingLimitArgs = g.NewQueryStruct("SetSpendingLimitArgs").
	PredefinedQueryStructField("spendingLimit", "int", g.ParameterOptions().Required().NoEquals())

var setEmailNotificationsArgs = g.NewQueryStruct("SetEmailNotificationsArgs").
	PredefinedQueryStructField("notificationIntegration", "*string", g.ParameterOptions().NoEquals().SingleQuotes()).
	PredefinedQueryStructField("emails", "string", g.ParameterOptions().Required().NoEquals().SingleQuotes())

var getNotificationIntegrationsResult = g.StructPair(
	"getNotificationIntegrationsRow",
	"BudgetNotificationIntegration",
).Text("integration_name").
	Number("last_notification_time").
	Time("added_date")

var setCycleStartActionArgs = g.NewQueryStruct("SetCycleStartActionArgs").
	PredefinedQueryStructField("procedure", "string", g.ParameterOptions().Required().NoEquals().SingleQuotes()).
	PredefinedQueryStructField("arguments", "string", g.ParameterOptions().Required().NoEquals())

var getCycleStartActionResult = g.StructPair(
	"getCycleStartActionRow",
	"BudgetCycleStartAction",
).Text("action_uuid").
	Text("procedure_fqn").
	Text("procedure_args").
	Time("added_timestamp").
	Time("last_triggered_timestamp")

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
		PredefinedQueryStructField("constructor", "bool", g.StaticOptions().SQL("()")).
		PredefinedQueryStructField("comment", "*string", g.ParameterOptions().SQL("COMMENT").SingleQuotes()).
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
).InstanceMethodOperation(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/method/set_spending_limit",
	"SET_SPENDING_LIMIT",
	setSpendingLimitArgs,
	nil,
	g.InstanceMethodKind("string"),
).InstanceMethodOperation(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/method/get_spending_limit",
	"GET_SPENDING_LIMIT",
	nil,
	nil,
	g.InstanceMethodKind("int"),
).InstanceMethodOperation(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/method/set_email_notifications",
	"SET_EMAIL_NOTIFICATIONS",
	setEmailNotificationsArgs,
	nil,
	g.InstanceMethodKind("string"),
).InstanceMethodOperation(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/method/get_notification_integrations",
	"GET_NOTIFICATION_INTEGRATIONS",
	nil,
	getNotificationIntegrationsResult,
	g.InstanceMethodKindSlice,
).InstanceMethodOperation(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/methods/set_cycle_start_action",
	"SET_CYCLE_START_ACTION",
	setCycleStartActionArgs,
	nil,
	g.InstanceMethodKind("string"),
).InstanceMethodOperation(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/methods/get_cycle_start_action",
	"GET_CYCLE_START_ACTION",
	nil,
	getCycleStartActionResult,
	g.InstanceMethodKindSingleValue,
)
