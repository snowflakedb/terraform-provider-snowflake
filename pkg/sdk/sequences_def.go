package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

var sequenceSet = g.NewQueryStruct("SequenceSet").
	PredefinedQueryStructField("ValuesBehavior", "*ValuesBehavior", g.KeywordOptions()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes())

var sequenceConstraint = g.NewQueryStruct("SequenceConstraint").
	OptionalSQL("CASCADE").
	OptionalSQL("RESTRICT").
	WithValidation(g.ExactlyOneValueSet, "Cascade", "Restrict")

var SequencesDef = g.NewInterface(
	"Sequences",
	"Sequence",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-sequence",
	g.NewQueryStruct("CreateSequence").
		Create().
		OrReplace().
		SQL("SEQUENCE").
		IfNotExists().
		Name().
		OptionalNumberAssignment("START", g.ParameterOptions().NoQuotes()).
		OptionalNumberAssignment("INCREMENT", g.ParameterOptions().NoQuotes()).
		PredefinedQueryStructField("ValuesBehavior", "*ValuesBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-sequence",
	g.NewQueryStruct("AlterSequence").
		Alter().
		SQL("SEQUENCE").
		IfExists().
		Name().
		Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalNumberAssignment("SET INCREMENT", g.ParameterOptions().NoQuotes()).
		OptionalQueryStructField(
			"Set",
			sequenceSet,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalSQL("UNSET COMMENT").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "SetIncrement", "Set", "UnsetComment"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-sequences",
	g.DbStruct("sequenceRow").
		Field("created_on", "string").
		Field("name", "string").
		Field("schema_name", "string").
		Field("database_name", "string").
		Field("next_value", "int").
		Field("interval", "int").
		Field("owner", "string").
		Field("owner_role_type", "string").
		Field("comment", "string").
		Field("ordered", "string"),
	g.PlainStruct("Sequence").
		Field("CreatedOn", "string").
		Field("Name", "string").
		Field("SchemaName", "string").
		Field("DatabaseName", "string").
		Field("NextValue", "int").
		Field("Interval", "int").
		Field("Owner", "string").
		Field("OwnerRoleType", "string").
		Field("Comment", "string").
		Field("Ordered", "bool"),
	g.NewQueryStruct("ShowSequences").
		Show().
		SQL("SEQUENCES").
		OptionalLike().
		OptionalIn(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDInFiltering,
	g.ShowByIDLikeFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-sequence",
	g.DbStruct("sequenceDetailRow").
		Field("created_on", "string").
		Field("name", "string").
		Field("schema_name", "string").
		Field("database_name", "string").
		Field("next_value", "int").
		Field("interval", "int").
		Field("owner", "string").
		Field("owner_role_type", "string").
		Field("comment", "string").
		Field("ordered", "string"),
	g.PlainStruct("SequenceDetail").
		Field("CreatedOn", "string").
		Field("Name", "string").
		Field("SchemaName", "string").
		Field("DatabaseName", "string").
		Field("NextValue", "int").
		Field("Interval", "int").
		Field("Owner", "string").
		Field("OwnerRoleType", "string").
		Field("Comment", "string").
		Field("Ordered", "bool"),
	g.NewQueryStruct("DescribeSequence").
		Describe().
		SQL("SEQUENCE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-sequence",
	g.NewQueryStruct("DropSequence").
		Drop().
		SQL("SEQUENCE").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Constraint",
			sequenceConstraint,
			g.KeywordOptions(),
		).
		WithValidation(g.ValidIdentifier, "name"),
)
