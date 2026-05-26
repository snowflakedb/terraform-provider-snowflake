package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var sequenceSet = g.NewQueryStruct("SequenceSet").
	PredefinedQueryStructField("ValuesBehavior", "*ValuesBehavior", g.KeywordOptions()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes())

var sequenceConstraint = g.NewQueryStruct("SequenceConstraint").
	OptionalSQL("CASCADE").
	OptionalSQL("RESTRICT").
	WithValidation(g.ExactlyOneValueSet, "Cascade", "Restrict")

var sequencesDef = g.NewInterface(
	"Sequences",
	"Sequence",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
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
		Identifier("RenameTo", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
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
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-sequences",
	g.StructPair("sequenceRow", "Sequence").
		Text("created_on").
		Text("name").
		Text("schema_name").
		Text("database_name").
		Number("next_value").
		Number("interval").
		Text("owner").
		Text("owner_role_type").
		Text("comment").
		Field("ordered", "string", "bool").
		WithConvertGeneration(),
	g.NewQueryStruct("ShowSequences").
		Show().
		SQL("SEQUENCES").
		OptionalLike().
		OptionalIn(),
	g.ShowByIDInFiltering,
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-sequence",
	g.StructPair("sequenceDetailRow", "SequenceDetail").
		Text("created_on").
		Text("name").
		Text("schema_name").
		Text("database_name").
		Number("next_value").
		Number("interval").
		Text("owner").
		Text("owner_role_type").
		Text("comment").
		Field("ordered", "string", "bool").
		WithConvertGeneration(),
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
