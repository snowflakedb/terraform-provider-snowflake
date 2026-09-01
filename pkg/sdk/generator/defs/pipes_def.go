package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

func pipeSet() *g.QueryStruct {
	return g.NewQueryStruct("PipeSet").
		OptionalTextAssignment("ERROR_INTEGRATION", g.ParameterOptions().NoQuotes()).
		OptionalBooleanAssignment("PIPE_EXECUTION_PAUSED", g.ParameterOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.AtLeastOneValueSet, "ErrorIntegration", "PipeExecutionPaused", "Comment")
}

func pipeUnset() *g.QueryStruct {
	return g.NewQueryStruct("PipeUnset").
		OptionalSQL("ERROR_INTEGRATION").
		OptionalSQL("PIPE_EXECUTION_PAUSED").
		OptionalSQL("COMMENT").
		WithValidation(g.AtLeastOneValueSet, "ErrorIntegration", "PipeExecutionPaused", "Comment")
}

func pipeRefresh() *g.QueryStruct {
	return g.NewQueryStruct("PipeRefresh").
		OptionalTextAssignment("PREFIX", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("MODIFIED_AFTER", g.ParameterOptions().SingleQuotes())
}

var pipePairs = g.StructPair("pipeDBRow", "Pipe").
	Text("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Text("definition").
	Text("owner").
	OptionalText("notification_channel", g.WithRequiredInPlain()).
	OptionalText("comment", g.WithRequiredInPlain()).
	OptionalText("integration", g.WithRequiredInPlain()).
	OptionalText("pattern", g.WithRequiredInPlain()).
	OptionalText("error_integration", g.WithRequiredInPlain()).
	OptionalText("owner_role_type", g.WithRequiredInPlain()).
	OptionalText("invalid_reason", g.WithRequiredInPlain())

var pipeDetailPairs = g.StructPair("pipeDetailRow", "PipeDetail").
	Text("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	Text("definition").
	Text("owner").
	OptionalText("notification_channel", g.WithRequiredInPlain()).
	OptionalText("comment", g.WithRequiredInPlain()).
	OptionalText("integration", g.WithRequiredInPlain()).
	OptionalText("pattern", g.WithRequiredInPlain()).
	OptionalText("error_integration", g.WithRequiredInPlain()).
	OptionalText("owner_role_type", g.WithRequiredInPlain()).
	OptionalText("invalid_reason", g.WithRequiredInPlain())

var pipesDef = g.NewInterface(
	"Pipes",
	"Pipe",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-pipe",
	g.NewQueryStruct("CreatePipe").
		Create().
		OrReplace().
		SQL("PIPE").
		IfNotExists().
		Name().
		OptionalBooleanAssignment("AUTO_INGEST", g.ParameterOptions()).
		OptionalTextAssignment("ERROR_INTEGRATION", g.ParameterOptions().NoQuotes()).
		OptionalTextAssignment("AWS_SNS_TOPIC", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("INTEGRATION", g.ParameterOptions().SingleQuotes()).
		OptionalComment().
		SQL("AS").
		Text("copyStatement", g.KeywordOptions().NoQuotes().Required()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithAdditionalValidations(),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-pipe",
	g.NewQueryStruct("AlterPipe").
		Alter().
		SQL("PIPE").
		IfExists().
		Name().
		OptionalQueryStructField("Set", pipeSet(), g.ListOptions().NoParentheses().SQL("SET")).
		OptionalQueryStructField("Unset", pipeUnset(), g.ListOptions().NoParentheses().SQL("UNSET")).
		OptionalSetTags().
		OptionalUnsetTags().
		OptionalQueryStructField("Refresh", pipeRefresh(), g.KeywordOptions().SQL("REFRESH")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags", "Refresh"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-pipe",
	g.NewQueryStruct("DropPipe").
		Drop().
		SQL("PIPE").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-pipes",
	pipePairs,
	g.NewQueryStruct("ShowPipes").
		Show().
		SQL("PIPES").
		OptionalLike().
		OptionalIn().
		WithAdditionalValidations(),
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-pipe",
	pipeDetailPairs,
	g.NewQueryStruct("DescribePipe").
		Describe().
		SQL("PIPE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
