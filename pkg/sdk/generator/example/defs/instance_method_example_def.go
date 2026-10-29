package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var addNotificationIntegrationArgs = g.NewQueryStruct("AddNotificationIntegrationArgs").
	Identifier("IntegrationName", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions()).
	WithValidation(g.ValidIdentifier, "IntegrationName")

// renders as [(SELECT SYSTEM$REFERENCE('TAG', '<fqn>', 'SESSION', 'APPLYBUDGET')), '<value>']
var exampleTagValue = g.NewQueryStruct("ExampleTagValue").
	Identifier("Tag", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required().SystemReference("TAG", "SESSION", "APPLYBUDGET").InSubquery()).
	Text("Value", g.KeywordOptions().SingleQuotes().Required()).
	WithValidation(g.ValidIdentifier, "Tag")

// renders as ([[<tag reference>, '<value>'], ...], ['<user>', ...], <bool>)
var setScopeArgs = g.NewQueryStruct("SetScopeArgs").
	ListQueryStructField("Tags", exampleTagValue, g.ListOptions().Required().MustBrackets().ArrayElements()).
	List("Users", "string", g.ListOptions().Required().MustBrackets().SingleQuotes()).
	PredefinedQueryStructField("Enabled", "bool", g.ParameterOptions().NoEquals().Required())

var addNotificationIntegrationResult = g.StructPair(
	"addNotificationIntegrationRow",
	"AddNotificationIntegrationResult",
).
	Text("status")

var InstanceMethodExamplesDef = g.NewInterface(
	"InstanceMethodExamples",
	"InstanceMethodExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).InstanceMethodOperation(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/method/add_notification_integration",
	"ADD_NOTIFICATION_INTEGRATION",
	addNotificationIntegrationArgs,
	addNotificationIntegrationResult,
	g.InstanceMethodKindSingleValue,
).InstanceMethodOperationScalar(
	"https://docs.snowflake.com/en/sql-reference/classes/budget/method/get_spending_limit",
	"GET_SPENDING_LIMIT",
	nil,
	"",
).InstanceMethodOperationScalar(
	"https://docs.snowflake.com/en/user-guide/budgets/per-user-quotas",
	"SET_SCOPE",
	setScopeArgs,
	"string",
)
