package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var addNotificationIntegrationArgs = g.NewQueryStruct("AddNotificationIntegrationArgs").
	Identifier("IntegrationName", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions()).
	WithValidation(g.ValidIdentifier, "IntegrationName")

var addNotificationIntegrationResult = g.StructPair(
	"addNotificationIntegrationRow",
	"AddNotificationIntegrationResult",
).Text("status")

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
)
