package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

const postgresMirrorsDoc = "https://docs.snowflake.com/en/user-guide/snowflake-postgres/postgres-data-mirroring-reference"

var PostgresMirrorTargetTableTypeEnumDef = g.NewEnum(
	"PostgresMirrorTargetTableType", "PostgresMirrorTargetTableTypes",
	"SNOWFLAKE", "ICEBERG",
)

// PostgresMirrorStatusEnumDef lists the statuses DESCRIBE_MIRROR returns today. additionalConvert keeps an unrecognized
// value instead of calling ToPostgresMirrorStatus, so a new status does not fail Describe.
var PostgresMirrorStatusEnumDef = g.NewEnum(
	"PostgresMirrorStatus", "PostgresMirrorStatuses",
	"STARTING", "ACTIVE", "SUSPENDED",
)

var postgresMirrorDetailsPairs = g.StructPair("postgresMirrorDetailsRow", "PostgresMirrorDetails").
	Text("MIRROR_NAME", g.WithPlainFieldName("Name")).
	OptionalText("POSTGRES_INSTANCE", g.WithRequiredInPlain()).
	OptionalText("POSTGRES_DATABASE", g.WithRequiredInPlain()).
	OptionalText("TARGET_DATABASE", g.WithRequiredInPlain()).
	Field("STATUS", "sql.NullString", "PostgresMirrorStatus", g.WithManualConvert()).
	OptionalText("REFRESH_INTERVAL", g.WithRequiredInPlain()).
	// Unknown types must fail Describe. Generated enum mapping expects a pointer and drops the error.
	OptionalEnum("TARGET_TABLE_TYPE", PostgresMirrorTargetTableTypeEnumDef, g.WithRequiredInPlain(), g.WithManualConvert()).
	OptionalText("EXTERNAL_VOLUME").
	// Nil when the apply task runs on serverless compute.
	OptionalText("WAREHOUSE").
	// JSON arrays read from the publication on the source instance; see PostgresMirrorDetails.IsMembershipKnown.
	OptionalPlainField("POSTGRES_TABLES", "[]string", g.WithManualConvert()).
	OptionalPlainField("POSTGRES_SCHEMAS", "[]string", g.WithManualConvert()).
	OptionalBool("FEED_ONLY").
	// *int64, not *int: int is 32 bits on 386, which this provider ships.
	Field("LAST_OPERATION_ID", "sql.NullInt64", "*int64").
	OptionalTime("LAST_APPLY_TIME").
	OptionalTime("LAST_OPERATION_TIME").
	OptionalText("ERROR_MESSAGE").
	OptionalText("CATALOG_INTEGRATION_NAME").
	// The JSON array returned by the procedure, compacted.
	OptionalText("RECENT_QUERY_DURATIONS", g.WithManualConvert()).
	OptionalNumber("APPLY_LAG_SECONDS").
	Field("APPLY_LAG_OPERATIONS", "sql.NullInt64", "*int64").
	Field("APPLY_LAG_BYTES", "sql.NullInt64", "*int64").
	OptionalText("LAST_APPLIED_COMMIT_LSN").
	OptionalTime("LAST_APPLIED_COMMIT_TIME").
	// The JSON object returned by the procedure, compacted; its keys depend on the source's snowflake_cdc version.
	OptionalText("POSTGRES_STATUS", g.WithManualConvert())

var postgresMirrorPairs = g.StructPair("postgresMirrorRow", "PostgresMirror").
	Text("MIRROR_NAME", g.WithPlainFieldName("Name")).
	OptionalText("POSTGRES_INSTANCE", g.WithRequiredInPlain()).
	OptionalText("POSTGRES_DATABASE", g.WithRequiredInPlain()).
	OptionalText("TARGET_DATABASE", g.WithRequiredInPlain()).
	// An unrecognized status is kept; see postgresMirrorRow.additionalConvert.
	Field("STATUS", "sql.NullString", "PostgresMirrorStatus", g.WithManualConvert()).
	OptionalNumber("TABLE_COUNT").
	OptionalNumber("INITIAL_SNAPSHOT_TABLES_REMAINING").
	OptionalText("REFRESH_INTERVAL", g.WithRequiredInPlain()).
	OptionalTime("LAST_APPLY_TIME").
	OptionalTime("LAST_OPERATION_TIME").
	OptionalText("ERROR_MESSAGE")

// postgresMirrorsDef defines Snowflake Postgres data mirrors. Mirrors are not catalog objects: there is no DDL, SHOW, or
// DESCRIBE for them, so every operation calls a procedure of the SNOWFLAKE application in SNOWFLAKE.POSTGRES.
// DropSafely recognizes the procedure's missing-mirror message.
var postgresMirrorsDef = g.NewInterface(
	"PostgresMirrors",
	"PostgresMirror",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](), // Mirrors are referenced by a plain name; required placeholder.
).CreateOperation(
	postgresMirrorsDoc,
	g.NewQueryStruct("CreatePostgresMirror").
		SQLWithCustomFieldName("createMirror", "CALL SNOWFLAKE.POSTGRES.CREATE_MIRROR").
		QueryStructField(
			"arguments",
			g.NewQueryStruct("CreatePostgresMirrorArguments").
				TextAssignment("mirror_name", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				TextAssignment("postgres_instance", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				TextAssignment("postgres_database", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				TextAssignment("target_database", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				PredefinedQueryStructField("PostgresTables", "[]StringListItemWrapper", g.ParameterOptions().NoEquals().Parentheses().SQL("postgres_tables => ARRAY_CONSTRUCT")).
				PredefinedQueryStructField("PostgresSchemas", "[]StringListItemWrapper", g.ParameterOptions().NoEquals().Parentheses().SQL("postgres_schemas => ARRAY_CONSTRUCT")).
				OptionalTextAssignment("refresh_interval", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				OptionalEnumAssignment("target_table_type", PostgresMirrorTargetTableTypeEnumDef, g.ParameterOptions().SingleQuotes().ArrowEquals()).
				OptionalTextAssignment("warehouse", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				OptionalTextAssignment("external_volume", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				WithValidation(g.ValidateValueSet, "PostgresInstance").
				WithValidation(g.ValidateValueSet, "PostgresDatabase").
				WithValidation(g.ValidateValueSet, "TargetDatabase").
				WithValidation(g.AtLeastOneValueSet, "PostgresTables", "PostgresSchemas").
				WithAdditionalValidations(),
			g.ListOptions().MustParentheses().Required(),
		),
).CustomOperation(
	"Alter",
	postgresMirrorsDoc,
	g.NewQueryStruct("AlterPostgresMirror").
		SQLWithCustomFieldName("alterMirror", "CALL SNOWFLAKE.POSTGRES.ALTER_MIRROR").
		QueryStructField(
			"arguments",
			g.NewQueryStruct("AlterPostgresMirrorArguments").
				TextAssignment("mirror_name", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				OptionalTextAssignment("refresh_interval", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				PredefinedQueryStructField("AddTables", "[]StringListItemWrapper", g.ParameterOptions().NoEquals().Parentheses().SQL("add_tables => ARRAY_CONSTRUCT")).
				PredefinedQueryStructField("RemoveTables", "[]StringListItemWrapper", g.ParameterOptions().NoEquals().Parentheses().SQL("remove_tables => ARRAY_CONSTRUCT")).
				PredefinedQueryStructField("AddSchemas", "[]StringListItemWrapper", g.ParameterOptions().NoEquals().Parentheses().SQL("add_schemas => ARRAY_CONSTRUCT")).
				PredefinedQueryStructField("RemoveSchemas", "[]StringListItemWrapper", g.ParameterOptions().NoEquals().Parentheses().SQL("remove_schemas => ARRAY_CONSTRUCT")).
				OptionalTextAssignment("warehouse", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				// An omitted warehouse leaves the apply task's compute unchanged; only an explicit NULL switches it to serverless.
				// The builder emits this keyword only for true. Generated valueSet would treat false as a change and as a
				// conflict with Warehouse, so those checks are in additionalValidations.
				OptionalSQLWithCustomFieldName("UnsetWarehouse", "warehouse => NULL").
				WithAdditionalValidations(),
			g.ListOptions().MustParentheses().Required(),
		),
).DropOperation(
	postgresMirrorsDoc,
	g.NewQueryStruct("DropPostgresMirror").
		SQLWithCustomFieldName("dropMirror", "CALL SNOWFLAKE.POSTGRES.DROP_MIRROR").
		QueryStructField(
			"arguments",
			g.NewQueryStruct("DropPostgresMirrorArguments").
				TextAssignment("mirror_name", g.ParameterOptions().SingleQuotes().ArrowEquals()).
				// Defaults to true in the procedure. A kept database is marked as previously mirrored, and its name can never be
				// a mirror target again.
				OptionalBooleanAssignment("drop_target_database", g.ParameterOptions().ArrowEquals()),
			g.ListOptions().MustParentheses().Required(),
		),
	// DROP_MIRROR has no IF EXISTS, and DropSafely needs the whole request to pass drop_target_database.
	g.WithDropSafelySuppressed(),
).CustomShowOperationWithPairedStructs(
	"Describe",
	g.ShowMappingKindSingleValue,
	postgresMirrorsDoc,
	postgresMirrorDetailsPairs,
	g.NewQueryStruct("DescribePostgresMirror").
		SQLWithCustomFieldName("describeMirror", "CALL SNOWFLAKE.POSTGRES.DESCRIBE_MIRROR").
		QueryStructField(
			"arguments",
			g.NewQueryStruct("DescribePostgresMirrorArguments").
				TextAssignment("mirror_name", g.ParameterOptions().SingleQuotes().ArrowEquals()),
			g.ListOptions().MustParentheses().Required(),
		),
).CustomShowOperationWithPairedStructs(
	"List",
	g.ShowMappingKindSlice,
	postgresMirrorsDoc,
	postgresMirrorPairs,
	g.NewQueryStruct("ListPostgresMirror").
		SQLWithCustomFieldName("listMirrors", "CALL SNOWFLAKE.POSTGRES.LIST_MIRRORS").
		QueryStructField(
			"arguments",
			g.NewQueryStruct("ListPostgresMirrorArguments").
				OptionalTextAssignment("postgres_instance", g.ParameterOptions().SingleQuotes().ArrowEquals()),
			g.ListOptions().MustParentheses().Required(),
		),
).WithCustomInterfaceMethod(
	"DropSafely",
	"DropSafely treats a mirror that no longer exists as dropped.",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*DropPostgresMirrorRequest")},
	"error",
).WithCustomInterfaceMethod(
	"WaitForInitialSnapshot",
	"WaitForInitialSnapshot polls LIST_MIRRORS until the mirror is ACTIVE, no table is left in its first copy, and at\n"+
		"least minimumTables tables are registered. minimumTables is the number of tables named explicitly on create and\n"+
		"must be at least 1: LIST_MIRRORS reports 0 remaining before registration. Tables reached only through a schema,\n"+
		"or added by a later alter, are included only when the caller counts them. The caller sets the deadline on ctx.",
	[]*g.MethodParameter{
		g.NewMethodParameter("postgresInstance", "string"),
		g.NewMethodParameter("name", "string"),
		g.NewMethodParameter("minimumTables", "int"),
	},
	"error",
).
	WithEnums(PostgresMirrorTargetTableTypeEnumDef, PostgresMirrorStatusEnumDef)
