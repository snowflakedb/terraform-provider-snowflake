package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var failoverGroupsDef = g.NewInterface(
	"FailoverGroups",
	"FailoverGroup",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-failover-group",
	g.NewQueryStruct("CreateFailoverGroup").
		Create().
		SQL("FAILOVER GROUP").
		IfNotExists().
		Name().
		ListAssignment("OBJECT_TYPES", "PluralObjectType", g.ParameterOptions().Required()).
		ListAssignment("ALLOWED_DATABASES", "AccountObjectIdentifier", g.ParameterOptions()).
		ListAssignment("ALLOWED_SHARES", "AccountObjectIdentifier", g.ParameterOptions()).
		ListAssignment("ALLOWED_INTEGRATION_TYPES", "IntegrationType", g.ParameterOptions()).
		ListAssignment("ALLOWED_ACCOUNTS", "AccountIdentifier", g.ParameterOptions().Required()).
		OptionalSQL("IGNORE EDITION CHECK").
		OptionalTextAssignment("REPLICATION_SCHEDULE", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateSecondaryReplicationGroup",
	"https://docs.snowflake.com/en/sql-reference/sql/create-failover-group",
	g.NewQueryStruct("CreateSecondaryReplicationGroup").
		Create().
		SQL("FAILOVER GROUP").
		IfNotExists().
		Name().
		Identifier("PrimaryFailoverGroup", g.KindOfT[sdkcommons.ExternalObjectIdentifier](),
			g.IdentifierOptions().SQL("AS REPLICA OF").Required()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "PrimaryFailoverGroup"),
).CustomOperation(
	"AlterSource",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-failover-group",
	g.NewQueryStruct("AlterSourceFailoverGroup").
		Alter().
		SQL("FAILOVER GROUP").
		IfExists().
		Name().
		OptionalIdentifier("NewName", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](),
			g.IdentifierOptions().SQL("RENAME TO")).
		OptionalQueryStructField("Set",
			g.NewQueryStruct("FailoverGroupSet").
				ListAssignment("OBJECT_TYPES", "PluralObjectType", g.ParameterOptions()).
				ListAssignment("ALLOWED_INTEGRATION_TYPES", "IntegrationType", g.ParameterOptions()).
				OptionalTextAssignment("REPLICATION_SCHEDULE", g.ParameterOptions().SingleQuotes()).
				WithAdditionalValidations(),
			g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("Unset",
			g.NewQueryStruct("FailoverGroupUnset").
				OptionalSQL("REPLICATION_SCHEDULE").
				WithValidation(g.AtLeastOneValueSet, "ReplicationSchedule"),
			g.KeywordOptions().SQL("UNSET")).
		OptionalQueryStructField("Add",
			g.NewQueryStruct("FailoverGroupAdd").
				ListAssignmentWithFieldName("TO ALLOWED_DATABASES", "AccountObjectIdentifier", g.ParameterOptions().Reverse(), "AllowedDatabases").
				ListAssignmentWithFieldName("TO ALLOWED_SHARES", "AccountObjectIdentifier", g.ParameterOptions().Reverse(), "AllowedShares").
				ListAssignmentWithFieldName("TO ALLOWED_ACCOUNTS", "AccountIdentifier", g.ParameterOptions().Reverse(), "AllowedAccounts").
				OptionalSQL("IGNORE_EDITION_CHECK"),
			g.KeywordOptions().SQL("ADD")).
		OptionalQueryStructField("Move",
			g.NewQueryStruct("FailoverGroupMove").
				ListAssignment("DATABASES", "AccountObjectIdentifier", g.ParameterOptions().NoEquals()).
				ListAssignment("SHARES", "AccountObjectIdentifier", g.ParameterOptions().NoEquals()).
				Identifier("To", g.KindOfT[sdkcommons.AccountObjectIdentifier](),
					g.IdentifierOptions().SQL("TO FAILOVER GROUP").Required()),
			g.KeywordOptions().SQL("MOVE")).
		OptionalQueryStructField("Remove",
			g.NewQueryStruct("FailoverGroupRemove").
				ListAssignmentWithFieldName("FROM ALLOWED_DATABASES", "AccountObjectIdentifier", g.ParameterOptions().Reverse(), "AllowedDatabases").
				ListAssignmentWithFieldName("FROM ALLOWED_SHARES", "AccountObjectIdentifier", g.ParameterOptions().Reverse(), "AllowedShares").
				ListAssignmentWithFieldName("FROM ALLOWED_ACCOUNTS", "AccountIdentifier", g.ParameterOptions().Reverse(), "AllowedAccounts"),
			g.KeywordOptions().SQL("REMOVE")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "NewName", "Set", "Unset", "Add", "Move", "Remove"),
).CustomOperation(
	"AlterTarget",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-failover-group",
	g.NewQueryStruct("AlterTargetFailoverGroup").
		Alter().
		SQL("FAILOVER GROUP").
		IfExists().
		Name().
		OptionalSQL("REFRESH").
		OptionalSQL("PRIMARY").
		OptionalSQL("SUSPEND").
		OptionalSQL("RESUME").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Refresh", "Primary", "Suspend", "Resume"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-failover-group",
	g.NewQueryStruct("DropFailoverGroup").
		Drop().
		SQL("FAILOVER GROUP").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-failover-groups",
	g.StructPair("failoverGroupDBRow", "FailoverGroup").
		Text("region_group").
		Text("snowflake_region").
		Time("created_on").
		Text("account_name").
		Text("name").
		Text("type").
		OptionalText("comment", g.WithRequiredInPlain()).
		Field("is_primary", "bool", "bool").
		// ExternalObjectIdentifier: NewExternalObjectIdentifierFromFullyQualifiedName has no error return
		Field("primary", "string", "ExternalObjectIdentifier", g.WithManualConvert()).
		// CSV with special-case "ACCOUNT PARAMETERS" rewrite
		Field("object_types", "string", "[]PluralObjectType", g.WithManualConvert()).
		// CSV with underscore→space replacement
		Field("allowed_integration_types", "string", "[]IntegrationType", g.WithManualConvert()).
		// CSV of "org.account" pairs
		Field("allowed_accounts", "string", "[]AccountIdentifier", g.WithManualConvert()).
		Text("organization_name").
		Text("account_locator").
		OptionalText("replication_schedule", g.WithRequiredInPlain()).
		// sql.NullString → FailoverGroupSecondaryState with non-zero default
		Field("secondary_state", "sql.NullString", "FailoverGroupSecondaryState", g.WithManualConvert()).
		OptionalText("next_scheduled_refresh", g.WithRequiredInPlain()).
		OptionalText("owner", g.WithRequiredInPlain()),
	g.NewQueryStruct("ShowFailoverGroups").
		Show().
		SQL("FAILOVER GROUPS").
		WithField(g.NewField("InAccount", "AccountIdentifier", g.Tags().Identifier().SQL("IN ACCOUNT"), nil)),
	g.ShowByIDNoFiltering,
).CustomShowOperationWithPairedStructs(
	"ShowFailoverGroupDatabases",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/show-databases-in-failover-group",
	g.StructPair("failoverGroupDatabaseDBRow", "FailoverGroupDatabase").
		Text("name"),
	g.NewQueryStruct("ShowFailoverGroupDatabases").
		Show().
		SQL("DATABASES").
		Identifier("In", g.KindOfT[sdkcommons.AccountObjectIdentifier](),
			g.IdentifierOptions().SQL("IN FAILOVER GROUP").Required()),
).CustomShowOperationWithPairedStructs(
	"ShowFailoverGroupShares",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/show-shares-in-failover-group",
	g.StructPair("failoverGroupShareDBRow", "FailoverGroupShare").
		Text("name").
		Text("owner_account"),
	g.NewQueryStruct("ShowFailoverGroupShares").
		Show().
		SQL("SHARES").
		Identifier("In", g.KindOfT[sdkcommons.AccountObjectIdentifier](),
			g.IdentifierOptions().SQL("IN FAILOVER GROUP").Required()),
).
	WithShowByIDFindPredicateKind(g.ShowByIDFindPredicateNameAndLocator).
	WithCustomInterfaceMethod(
		"ShowDatabases",
		"// ShowDatabases returns the list of databases in the failover group as identifiers.",
		[]*g.MethodParameter{
			g.NewMethodParameter("id", "AccountObjectIdentifier"),
		},
		"[]AccountObjectIdentifier", "error",
	).
	WithCustomInterfaceMethod(
		"ShowShares",
		"// ShowShares returns the list of shares in the failover group as identifiers.",
		[]*g.MethodParameter{
			g.NewMethodParameter("id", "AccountObjectIdentifier"),
		},
		"[]AccountObjectIdentifier", "error",
	)
