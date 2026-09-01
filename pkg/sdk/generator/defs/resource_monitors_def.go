package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var resourceMonitorLevelEnum = g.NewEnum(
	"ResourceMonitorLevel", "ResourceMonitorLevels",
	"ACCOUNT", "WAREHOUSE",
)

var triggerActionEnum = g.NewEnum(
	"TriggerAction", "TriggerActions",
	"SUSPEND", "SUSPEND_IMMEDIATE", "NOTIFY",
)

var frequencyEnum = g.NewEnum(
	"Frequency", "Frequencies",
	"MONTHLY", "DAILY", "WEEKLY", "YEARLY", "NEVER",
)

func notifyUsersStruct() *g.QueryStruct {
	return g.NewQueryStruct("NotifyUsers").
		ListQueryStructField("Users", notifiedUserStruct(), g.ListOptions().Parentheses())
}

func notifiedUserStruct() *g.QueryStruct {
	return g.NewQueryStruct("NotifiedUser").
		Identifier("Name", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Required())
}

func triggerDefinitions() *g.QueryStruct {
	return g.NewQueryStruct("TriggerDefinition").
		WithField(g.NewField("Threshold", "int", nil, g.ParameterOptions().NoEquals().SQL("ON").Required())).
		WithField(g.NewField("TriggerAction", "TriggerAction", nil, g.ParameterOptions().NoEquals().SQL("PERCENT DO").Required()))
}

var resourceMonitorWithStruct = g.NewQueryStruct("ResourceMonitorWith").
	OptionalNumberAssignment("CREDIT_QUOTA", g.ParameterOptions()).
	OptionalEnumAssignment("FREQUENCY", frequencyEnum, g.ParameterOptions()).
	OptionalTextAssignment("START_TIMESTAMP", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("END_TIMESTAMP", g.ParameterOptions().SingleQuotes()).
	OptionalQueryStructField("NotifyUsers", notifyUsersStruct(), g.ParameterOptions().SQL("NOTIFY_USERS")).
	ListQueryStructField("Triggers", triggerDefinitions(), g.KeywordOptions().SQL("TRIGGERS").NoComma())

var resourceMonitorSetStruct = g.NewQueryStruct("ResourceMonitorSet").
	OptionalNumberAssignment("CREDIT_QUOTA", g.ParameterOptions()).
	OptionalEnumAssignment("FREQUENCY", frequencyEnum, g.ParameterOptions()).
	OptionalTextAssignment("START_TIMESTAMP", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("END_TIMESTAMP", g.ParameterOptions().SingleQuotes()).
	OptionalQueryStructField("NotifyUsers", notifyUsersStruct(), g.ParameterOptions().SQL("NOTIFY_USERS")).
	WithValidation(g.AtLeastOneValueSet, "CreditQuota", "Frequency", "StartTimestamp", "EndTimestamp", "NotifyUsers")

var resourceMonitorUnsetStruct = g.NewQueryStruct("ResourceMonitorUnset").
	OptionalSQLWithCustomFieldName("CreditQuota", "CREDIT_QUOTA = null").
	OptionalSQLWithCustomFieldName("EndTimestamp", "END_TIMESTAMP = null").
	OptionalSQLWithCustomFieldName("NotifyUsers", "NOTIFY_USERS = ()")

var resourceMonitorPairs = g.StructPair("resourceMonitorRow", "ResourceMonitor").
	Text("name").
	Field("credit_quota", "sql.NullString", "float64", g.WithManualConvert()).
	Field("used_credits", "sql.NullString", "float64", g.WithManualConvert()).
	Field("remaining_credits", "sql.NullString", "float64", g.WithManualConvert()).
	OptionalEnum("level", resourceMonitorLevelEnum).
	OptionalEnum("frequency", frequencyEnum).
	OptionalText("start_time", g.WithRequiredInPlain()).
	OptionalText("end_time", g.WithRequiredInPlain()).
	Field("notify_at", "sql.NullString", "[]int", g.WithManualConvert()).
	Field("suspend_at", "sql.NullString", "*int", g.WithManualConvert()).
	Field("suspend_immediately_at", "sql.NullString", "*int", g.WithManualConvert()).
	Time("created_on").
	Text("owner").
	OptionalText("comment", g.WithRequiredInPlain()).
	OptionalPlainField("notify_users", "[]string")

var resourceMonitorsDef = g.NewInterface(
	"ResourceMonitors",
	"ResourceMonitor",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-resource-monitor",
	g.NewQueryStruct("CreateResourceMonitor").
		Create().
		OrReplace().
		SQL("RESOURCE MONITOR").
		IfNotExists().
		Name().
		OptionalQueryStructField("With", resourceMonitorWithStruct, g.KeywordOptions().SQL("WITH")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithAdditionalValidations(),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-resource-monitor",
	g.NewQueryStruct("AlterResourceMonitor").
		Alter().
		SQL("RESOURCE MONITOR").
		IfExists().
		Name().
		OptionalQueryStructField("Set", resourceMonitorSetStruct, g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("Unset", resourceMonitorUnsetStruct, g.KeywordOptions().SQL("SET")).
		ListQueryStructField("Triggers", triggerDefinitions(), g.KeywordOptions().SQL("TRIGGERS").NoComma()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.AtLeastOneValueSet, "Set", "Unset", "Triggers").
		WithAdditionalValidations(),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-resource-monitor",
	g.NewQueryStruct("DropResourceMonitor").
		Drop().
		SQL("RESOURCE MONITOR").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-resource-monitors",
	resourceMonitorPairs,
	g.NewQueryStruct("ShowResourceMonitors").
		Show().
		SQL("RESOURCE MONITORS").
		OptionalLike(),
	g.ShowByIDLikeFiltering,
).WithEnums(
	resourceMonitorLevelEnum,
	triggerActionEnum,
	frequencyEnum,
)
