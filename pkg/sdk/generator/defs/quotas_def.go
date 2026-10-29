package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// Per-user quotas are documented only in the user guide (https://docs.snowflake.com/en/user-guide/budgets/per-user-quotas);
// there is no class reference yet, so where the guide is silent the definitions below follow Snowflake's behavior.

const quotasDocs = "https://docs.snowflake.com/en/user-guide/budgets/per-user-quotas"

var QuotaSharedResourceDomainEnumDef = g.NewEnum(
	"QuotaSharedResourceDomain", "QuotaSharedResourceDomains",
	"WAREHOUSE", "AI FUNCTION", "CORTEX AGENT", "CORTEX CODE", "SNOWFLAKE INTELLIGENCE", "AI GATEWAY",
)

var QuotaCycleEnumDef = g.NewEnum(
	"QuotaCycle", "QuotaCycles",
	"MONTHLY", "WEEKLY", "DAILY",
)

var QuotaSpendStrategyEnumDef = g.NewEnum(
	"QuotaSpendStrategy", "QuotaSpendStrategies",
	"ACTUAL", "PROJECTED",
)

var QuotaUserTagsOperatorEnumDef = g.NewEnum(
	"QuotaUserTagsOperator", "QuotaUserTagsOperators",
	"UNION", "INTERSECTION",
)

// AI FUNCTION and CORTEX CODE instances are plain names; the other domains take object references.
// Without an instance, the whole domain is used.
func quotaSharedResourceArgs(name string) *g.QueryStruct {
	reference := func(objectType string) *g.IdentifierTransformer {
		return g.IdentifierOptions().SystemReference(objectType).InSubquery()
	}
	return g.NewQueryStruct(name).
		Enum("Domain", QuotaSharedResourceDomainEnumDef, g.KeywordOptions().SingleQuotes().Required()).
		OptionalText("InstanceName", g.KeywordOptions().SingleQuotes()).
		OptionalIdentifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), reference("WAREHOUSE")).
		OptionalIdentifier("CortexAgent", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), reference("CORTEX AGENT")).
		OptionalIdentifier("SnowflakeIntelligence", g.KindOfT[sdkcommons.AccountObjectIdentifier](), reference("SNOWFLAKE INTELLIGENCE")).
		// the AI Gateway reference follows the guide; only its whole domain is tested, as an instance needs the AI Gateway preview
		OptionalIdentifier("AiGateway", g.KindOfT[sdkcommons.AccountObjectIdentifier](), reference("AI GATEWAY")).
		WithValidation(g.MoreThanOneValueSet, "InstanceName", "Warehouse", "CortexAgent", "SnowflakeIntelligence", "AiGateway").
		WithAdditionalValidations()
}

func quotaTagValue(name string) *g.QueryStruct {
	return g.NewQueryStruct(name).
		Identifier("Tag", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required().SystemReference("TAG", "SESSION", "APPLYBUDGET").InSubquery()).
		Text("Value", g.KeywordOptions().SingleQuotes().Required()).
		WithValidation(g.ValidIdentifier, "Tag")
}

// SET_USER_TAGS replaces the tags that select the users (and any tags excluded with EXCLUDE_USERS('TAG', ...)).
// The operator is required; the one-argument form is rejected.
var quotaSetUserTagsArgs = g.NewQueryStruct("SetUserTagsArgs").
	ListQueryStructField("Tags", quotaTagValue("QuotaUserTag"), g.ListOptions().Required().MustBrackets().ArrayElements()).
	Enum("Operator", QuotaUserTagsOperatorEnumDef, g.KeywordOptions().SingleQuotes().Required())

// The 'USER' form accumulates like INCLUDE_USERS, and an empty array clears the list. The 'TAG' form replaces the user tags
// (the quota then applies to all users without the given tags) and also clears the included and excluded users.
var quotaExcludeUsersArgs = g.NewQueryStruct("ExcludeUsersArgs").
	OptionalQueryStructField(
		"ByTag",
		g.NewQueryStruct("QuotaExcludeUsersByTag").
			SQLWithCustomFieldName("kind", "'TAG'").
			ListQueryStructField("Tags", quotaTagValue("QuotaExcludedUserTag"), g.ListOptions().Required().MustBrackets().ArrayElements()),
		g.ListOptions(),
	).
	OptionalQueryStructField(
		"ByUser",
		g.NewQueryStruct("QuotaExcludeUsersByUser").
			SQLWithCustomFieldName("kind", "'USER'").
			List("Users", "string", g.ListOptions().Required().MustBrackets().SingleQuotes()),
		g.ListOptions(),
	).
	WithValidation(g.ExactlyOneValueSet, "ByTag", "ByUser")

// INCLUDE_USERS calls accumulate; an empty array clears the list. User names are resolved case-insensitively.
var quotaIncludeUsersArgs = g.NewQueryStruct("IncludeUsersArgs").
	List("Users", "string", g.ListOptions().Required().MustBrackets().SingleQuotes())

// The limit must be a positive integer no greater than 1000000000.
var quotaSetPerUserLimitArgs = g.NewQueryStruct("SetPerUserLimitArgs").
	PredefinedQueryStructField("Limit", "int", g.ParameterOptions().NoEquals().Required()).
	OptionalEnum("Cycle", QuotaCycleEnumDef, g.KeywordOptions().SingleQuotes())

var quotaUnsetPerUserLimitArgs = g.NewQueryStruct("UnsetPerUserLimitArgs").
	Enum("Cycle", QuotaCycleEnumDef, g.KeywordOptions().SingleQuotes().Required())

// The one-argument form also disables the per-user block notifications, so NotifyUsers should always be passed.
var quotaSetBlockEnforcementEnabledArgs = g.NewQueryStruct("SetBlockEnforcementEnabledArgs").
	PredefinedQueryStructField("Enabled", "bool", g.ParameterOptions().NoEquals().Required()).
	PredefinedQueryStructField("NotifyUsers", "*bool", g.ParameterOptions().NoEquals())

// The threshold must be greater than 0 and no greater than 1000; all four arguments are required.
var quotaAddNotificationThresholdArgs = g.NewQueryStruct("AddNotificationThresholdArgs").
	PredefinedQueryStructField("Threshold", "int", g.ParameterOptions().NoEquals().Required()).
	Enum("SpendStrategy", QuotaSpendStrategyEnumDef, g.KeywordOptions().SingleQuotes().Required()).
	PredefinedQueryStructField("NotifyUser", "bool", g.ParameterOptions().NoEquals().Required()).
	Enum("Cycle", QuotaCycleEnumDef, g.KeywordOptions().SingleQuotes().Required())

// Thresholds are identified by (threshold, spend strategy, cycle); without the cycle, the monthly threshold is removed.
var quotaRemoveNotificationThresholdArgs = g.NewQueryStruct("RemoveNotificationThresholdArgs").
	PredefinedQueryStructField("Threshold", "int", g.ParameterOptions().NoEquals().Required()).
	Enum("SpendStrategy", QuotaSpendStrategyEnumDef, g.KeywordOptions().SingleQuotes().Required()).
	OptionalEnum("Cycle", QuotaCycleEnumDef, g.KeywordOptions().SingleQuotes())

// Only webhook and queue integrations are accepted (email recipients are set with SET_ADMIN_EMAILS), and USAGE on the integration
// must be granted to the SNOWFLAKE application; otherwise ADD_NOTIFICATION_INTEGRATION fails with INVALID_INTEGRATION.
func quotaNotificationIntegrationArgs(name string) *g.QueryStruct {
	return g.NewQueryStruct(name).
		Identifier("Integration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SingleQuotes().Required()).
		WithValidation(g.ValidIdentifier, "Integration")
}

// SET_ADMIN_EMAILS takes a single string with comma-separated emails, which must be verified emails of users in the account;
// an empty string clears them.
var quotaSetAdminEmailsArgs = g.NewQueryStruct("SetAdminEmailsArgs").
	Text("Emails", g.KeywordOptions().SingleQuotes().Required())

// The arguments are passed after the IDs of the users who crossed the threshold, so their number must match the procedure's
// remaining parameters.
var quotaAddCustomActionArgs = g.NewQueryStruct("AddCustomActionArgs").
	Identifier("Procedure", g.KindOfT[sdkcommons.SchemaObjectIdentifierWithArguments](), g.IdentifierOptions().Required().SystemReference("PROCEDURE")).
	List("Arguments", "string", g.ListOptions().Required().MustBrackets().SingleQuotes()).
	Enum("SpendStrategy", QuotaSpendStrategyEnumDef, g.KeywordOptions().SingleQuotes().Required()).
	PredefinedQueryStructField("Threshold", "int", g.ParameterOptions().NoEquals().Required()).
	WithValidation(g.ValidIdentifier, "Procedure")

// The arguments are positional, and removing by procedure removes the actions of both spend strategies. The procedure is passed
// as the unquoted PROCEDURE_FQN returned by GET_CUSTOM_ACTIONS.
var quotaRemoveCustomActionsArgs = g.NewQueryStruct("RemoveCustomActionsArgs").
	OptionalNumber("Threshold", g.KeywordOptions()).
	OptionalText("ProcedureFqn", g.KeywordOptions().SingleQuotes()).
	WithAdditionalValidations()

var quotaSetCycleStartActionArgs = g.NewQueryStruct("SetCycleStartActionArgs").
	Identifier("Procedure", g.KindOfT[sdkcommons.SchemaObjectIdentifierWithArguments](), g.IdentifierOptions().Required().SystemReference("PROCEDURE")).
	List("Arguments", "string", g.ListOptions().Required().MustBrackets().SingleQuotes()).
	WithValidation(g.ValidIdentifier, "Procedure")

// GET_CONFIG returns one row; numbers and booleans ("0"/"1") are returned as strings.
var quotaConfigResult = g.StructPair("quotaConfigRow", "QuotaConfig").
	// uppercase column names used on purpose, here and in the results below
	OptionalText("QUOTA_ID").
	OptionalNumber("PER_USER_LIMIT").
	OptionalNumber("PER_USER_LIMIT_WEEKLY").
	OptionalNumber("PER_USER_LIMIT_DAILY").
	Bool("BLOCK_ENFORCEMENT_ENABLED").
	Bool("PER_USER_BLOCK_NOTIFICATIONS_ENABLED").
	Field("ADMIN_EMAILS", "sql.NullString", "[]string", g.WithCustomParser("ParseQuotaAdminEmails")).
	OptionalText("ADMIN_LAST_SENT_AT").
	OptionalText("REFRESH_TIER")

var quotaNotificationThresholdResult = g.StructPair("quotaNotificationThresholdRow", "QuotaNotificationThreshold").
	Number("THRESHOLD").
	Enum("SPEND_STRATEGY", QuotaSpendStrategyEnumDef).
	Bool("NOTIFY_USER").
	Enum("CYCLE", QuotaCycleEnumDef).
	Time("ADDED_TIMESTAMP")

// LAST_NOTIFICATION_TIME is NULL until a notification is sent; it's read as text, so it scans whatever its type is.
var quotaNotificationIntegrationResult = g.StructPair("quotaNotificationIntegrationRow", "QuotaNotificationIntegration").
	Text("INTEGRATION_NAME").
	OptionalText("LAST_NOTIFICATION_TIME").
	Time("ADDED_TIMESTAMP")

// PROCEDURE_FQN is returned unquoted (e.g. DB.SCH.PROC(VARCHAR)), and PROCEDURE_ARGS is a JSON array.
var quotaCustomActionResult = g.StructPair("quotaCustomActionRow", "QuotaCustomAction").
	Text("ACTION_ID").
	Text("PROCEDURE_FQN").
	Field("PROCEDURE_ARGS", "string", "[]string", g.WithCustomParser("ParseQuotaProcedureArguments")).
	Enum("SPEND_STRATEGY", QuotaSpendStrategyEnumDef).
	Number("THRESHOLD").
	OptionalTime("LAST_TRIGGER_ATTEMPT_TIME").
	Time("ADDED_TIMESTAMP")

// GET_CYCLE_START_ACTION returns no rows when no action is set; the columns are nullable so that a row of NULLs scans too.
var quotaCycleStartActionResult = g.StructPair("quotaCycleStartActionRow", "QuotaCycleStartAction").
	OptionalText("ACTION_ID").
	OptionalText("PROCEDURE_FQN").
	Field("PROCEDURE_ARGS", "sql.NullString", "[]string", g.WithCustomParser("ParseQuotaProcedureArguments")).
	OptionalTime("ADDED_TIMESTAMP").
	OptionalTime("LAST_TRIGGER_ATTEMPT_TIME")

var quotasDef = g.NewInterface(
	"Quotas",
	"Quota",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	quotasDocs,
	g.NewQueryStruct("CreateQuotaOptions").
		Create().
		OrReplace().
		SQLWithCustomFieldName("snowflakeCoreQuota", "SNOWFLAKE.CORE.QUOTA").
		IfNotExists().
		Name().
		SQLWithCustomFieldName("parens", "()").
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	quotasDocs,
	g.NewQueryStruct("AlterQuotaOptions").
		Alter().
		SQLWithCustomFieldName("snowflakeCoreQuota", "SNOWFLAKE.CORE.QUOTA").
		Name().
		RenameTo().
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("QuotaSet").
				OptionalComment().
				WithValidation(g.AtLeastOneValueSet, "Comment"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("QuotaUnset").
				OptionalSQL("COMMENT").
				WithValidation(g.AtLeastOneValueSet, "Comment"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set", "Unset"),
).DropOperation(
	quotasDocs,
	g.NewQueryStruct("DropQuotaOptions").
		Drop().
		SQLWithCustomFieldName("snowflakeCoreQuota", "SNOWFLAKE.CORE.QUOTA").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	quotasDocs,
	g.StructPair("quotaDBRow", "Quota").
		Time("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		OptionalText("current_version").
		OptionalText("comment").
		Text("owner").
		OptionalText("owner_role_type"),
	g.NewQueryStruct("ShowQuotas").
		Show().
		SQLWithCustomFieldName("snowflakeCoreQuota", "SNOWFLAKE.CORE.QUOTA").
		OptionalLike().
		OptionalIn().
		OptionalLimit(),
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).InstanceMethodOperationScalar(
	quotasDocs,
	"ADD_SHARED_RESOURCE",
	quotaSharedResourceArgs("AddSharedResourceArgs"),
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"REMOVE_SHARED_RESOURCE",
	quotaSharedResourceArgs("RemoveSharedResourceArgs"),
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"SET_USER_TAGS",
	quotaSetUserTagsArgs,
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"EXCLUDE_USERS",
	quotaExcludeUsersArgs,
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"INCLUDE_USERS",
	quotaIncludeUsersArgs,
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"GET_QUOTA_SCOPE",
	nil,
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"SET_PER_USER_LIMIT",
	quotaSetPerUserLimitArgs,
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"UNSET_PER_USER_LIMIT",
	quotaUnsetPerUserLimitArgs,
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"SET_BLOCK_ENFORCEMENT_ENABLED",
	quotaSetBlockEnforcementEnabledArgs,
	"string",
).InstanceMethodOperation(
	quotasDocs,
	"GET_CONFIG",
	nil,
	quotaConfigResult,
	g.InstanceMethodKindSingleValue,
).InstanceMethodOperationScalar(
	quotasDocs,
	"ADD_NOTIFICATION_THRESHOLD",
	quotaAddNotificationThresholdArgs,
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"REMOVE_NOTIFICATION_THRESHOLD",
	quotaRemoveNotificationThresholdArgs,
	"string",
).InstanceMethodOperation(
	quotasDocs,
	"GET_NOTIFICATION_THRESHOLDS",
	nil,
	quotaNotificationThresholdResult,
	g.InstanceMethodKindSlice,
).InstanceMethodOperationScalar(
	quotasDocs,
	"ADD_NOTIFICATION_INTEGRATION",
	quotaNotificationIntegrationArgs("AddNotificationIntegrationArgs"),
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"REMOVE_NOTIFICATION_INTEGRATION",
	quotaNotificationIntegrationArgs("RemoveNotificationIntegrationArgs"),
	"string",
).InstanceMethodOperation(
	quotasDocs,
	"GET_NOTIFICATION_INTEGRATIONS",
	nil,
	quotaNotificationIntegrationResult,
	g.InstanceMethodKindSlice,
).InstanceMethodOperationScalar(
	quotasDocs,
	"SET_ADMIN_EMAILS",
	quotaSetAdminEmailsArgs,
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"ADD_CUSTOM_ACTION",
	quotaAddCustomActionArgs,
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"REMOVE_CUSTOM_ACTIONS",
	quotaRemoveCustomActionsArgs,
	"string",
).InstanceMethodOperation(
	quotasDocs,
	"GET_CUSTOM_ACTIONS",
	nil,
	quotaCustomActionResult,
	g.InstanceMethodKindSlice,
).InstanceMethodOperationScalar(
	quotasDocs,
	"SET_CYCLE_START_ACTION",
	quotaSetCycleStartActionArgs,
	"string",
).InstanceMethodOperationScalar(
	quotasDocs,
	"REMOVE_CYCLE_START_ACTION",
	nil,
	"string",
).InstanceMethodOperation(
	quotasDocs,
	"GET_CYCLE_START_ACTION",
	nil,
	quotaCycleStartActionResult,
	g.InstanceMethodKindSingleValue,
).
	WithCustomInterfaceMethod("DescribeDetails", "", []*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())}, "*QuotaDetails", "error").
	WithEnums(QuotaSharedResourceDomainEnumDef, QuotaCycleEnumDef, QuotaSpendStrategyEnumDef, QuotaUserTagsOperatorEnumDef)
