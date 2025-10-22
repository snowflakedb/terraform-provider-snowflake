package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

var eventTableSet = g.NewQueryStruct("EventTableSet").
	OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalBooleanAssignment("CHANGE_TRACKING", g.ParameterOptions()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes())

var eventTableUnset = g.NewQueryStruct("EventTableUnset").
	OptionalSQL("DATA_RETENTION_TIME_IN_DAYS").
	OptionalSQL("MAX_DATA_EXTENSION_TIME_IN_DAYS").
	OptionalSQL("CHANGE_TRACKING").
	OptionalSQL("COMMENT")

var eventTableAddRowAccessPolicy = g.NewQueryStruct("EventTableAddRowAccessPolicy").
	SQL("ADD").
	Identifier("RowAccessPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
	NamedListWithParens("ON", g.KindOfT[string](), g.KeywordOptions().Required()). // TODO: double quotes here?
	WithValidation(g.ValidIdentifier, "RowAccessPolicy")

var eventTableDropRowAccessPolicy = g.NewQueryStruct("EventTableDropRowAccessPolicy").
	SQL("DROP").
	Identifier("RowAccessPolicy", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ROW ACCESS POLICY").Required()).
	WithValidation(g.ValidIdentifier, "RowAccessPolicy")

var eventTableDropAndAddRowAccessPolicy = g.NewQueryStruct("EventTableDropAndAddRowAccessPolicy").
	QueryStructField("Drop", eventTableDropRowAccessPolicy, g.KeywordOptions().Required()).
	QueryStructField("Add", eventTableAddRowAccessPolicy, g.KeywordOptions().Required())

var eventTableClusteringAction = g.NewQueryStruct("EventTableClusteringAction").
	PredefinedQueryStructField("ClusterBy", "*[]string", g.KeywordOptions().Parentheses().SQL("CLUSTER BY")).
	OptionalSQL("SUSPEND RECLUSTER").
	OptionalSQL("RESUME RECLUSTER").
	OptionalSQL("DROP CLUSTERING KEY")

var searchOptimization = g.NewQueryStruct("SearchOptimization").
	SQL("SEARCH OPTIMIZATION").
	PredefinedQueryStructField("On", "[]string", g.KeywordOptions().SQL("ON"))

var eventTableSearchOptimizationAction = g.NewQueryStruct("EventTableSearchOptimizationAction").
	OptionalQueryStructField(
		"Add",
		searchOptimization,
		g.KeywordOptions().SQL("ADD"),
	).
	OptionalQueryStructField(
		"Drop",
		searchOptimization,
		g.KeywordOptions().SQL("DROP"),
	)

var EventTablesDef = g.NewInterface(
	"EventTables",
	"EventTable",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-event-table",
	g.NewQueryStruct("CreateEventTable").
		Create().
		OrReplace().
		SQL("EVENT TABLE").
		IfNotExists().
		Name().
		NamedListWithParens("CLUSTER BY", g.KindOfT[string](), g.KeywordOptions()).
		OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalBooleanAssignment("CHANGE_TRACKING", g.ParameterOptions()).
		OptionalTextAssignment("DEFAULT_DDL_COLLATION", g.ParameterOptions().SingleQuotes()).
		OptionalSQL("COPY GRANTS").
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("RowAccessPolicy", "*RowAccessPolicy", g.KeywordOptions()).
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-event-tables",
	g.DbStruct("eventTableRow").
		Field("created_on", "time.Time").
		Field("name", "string").
		Field("database_name", "string").
		Field("schema_name", "string").
		Field("owner", "sql.NullString").
		Field("comment", "sql.NullString").
		Field("owner_role_type", "sql.NullString"),
	g.PlainStruct("EventTable").
		Field("CreatedOn", "time.Time").
		Field("Name", "string").
		Field("DatabaseName", "string").
		Field("SchemaName", "string").
		Field("Owner", "string").
		Field("Comment", "string").
		Field("OwnerRoleType", "string"),
	g.NewQueryStruct("ShowEventTables").
		Show().
		Terse().
		SQL("EVENT TABLES").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalLimit(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDInFiltering,
	g.ShowByIDLikeFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-event-table",
	g.DbStruct("eventTableDetailsRow").
		Field("name", "string").
		Field("kind", "string").
		Field("comment", "string"),
	g.PlainStruct("EventTableDetails").
		Field("Name", "string").
		Field("Kind", "string").
		Field("Comment", "string"),
	g.NewQueryStruct("DescribeEventTable").
		Describe().
		SQL("EVENT TABLE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-table",
	g.NewQueryStruct("DropEventTable").
		Drop().
		SQL("TABLE").
		IfExists().
		Name().
		OptionalSQL("RESTRICT"). // CASCADE or RESTRICT, and CASCADE for Default
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-table-event-table",
	g.NewQueryStruct("AlterEventTable").
		Alter().
		SQL("TABLE").
		IfNotExists().
		Name().
		OptionalQueryStructField(
			"Set",
			eventTableSet,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			eventTableUnset,
			g.KeywordOptions().SQL("UNSET"),
		).
		OptionalQueryStructField("AddRowAccessPolicy", eventTableAddRowAccessPolicy, g.KeywordOptions()).
		OptionalQueryStructField("DropRowAccessPolicy", eventTableDropRowAccessPolicy, g.KeywordOptions()).
		OptionalQueryStructField("DropAndAddRowAccessPolicy", eventTableDropAndAddRowAccessPolicy, g.ListOptions().NoParentheses()).
		OptionalSQL("DROP ALL ROW ACCESS POLICIES").
		OptionalQueryStructField(
			"ClusteringAction",
			eventTableClusteringAction,
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"SearchOptimizationAction",
			eventTableSearchOptimizationAction,
			g.KeywordOptions(),
		).
		OptionalSetTags().
		OptionalUnsetTags().
		Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set", "Unset", "SetTags", "UnsetTags", "AddRowAccessPolicy", "DropRowAccessPolicy", "DropAndAddRowAccessPolicy", "DropAllRowAccessPolicies", "ClusteringAction", "SearchOptimizationAction"),
)
