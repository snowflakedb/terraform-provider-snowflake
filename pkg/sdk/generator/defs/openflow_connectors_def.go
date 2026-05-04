//go:build sdk_generation

package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var openflowConnectorsDef = g.NewInterface(
	"OpenflowConnectors",
	"OpenflowConnector",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-openflow-connector",
	g.NewQueryStruct("CreateOpenflowConnector").
		Create().
		SQL("OPENFLOW CONNECTOR").
		IfNotExists().
		Name().
		OptionalIdentifier("InRuntime", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("IN RUNTIME")).
		OptionalTextAssignment("FROM DEFINITION", g.ParameterOptions().NoQuotes().NoEquals()).
		OptionalTextAssignment("FROM", g.ParameterOptions().SingleQuotes().NoEquals()).
		OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-openflow-connector",
	g.NewQueryStruct("AlterOpenflowConnector").
		Alter().
		SQL("OPENFLOW CONNECTOR").
		Name().
		OptionalSQL("START").
		OptionalSQL("STOP").
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("OpenflowConnectorSet").
				OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.AtLeastOneValueSet, "DisplayName", "Comment"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("OpenflowConnectorUnset").
				OptionalSQL("DISPLAY_NAME").
				OptionalSQL("COMMENT").
				WithValidation(g.AtLeastOneValueSet, "DisplayName", "Comment"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Start", "Stop", "Set", "Unset"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-openflow-connector",
	g.NewQueryStruct("DropOpenflowConnector").
		Drop().
		SQL("OPENFLOW CONNECTOR").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-openflow-connectors",
	g.DbStruct("openflowConnectorRow").
		Text("name").
		Text("status").
		Text("runtime").
		OptionalText("connector_definition").
		OptionalText("display_name").
		Text("database_name").
		Text("schema_name").
		Text("owner").
		OptionalText("comment").
		Time("created_on").
		Time("updated_on"),
	g.PlainStruct("OpenflowConnector").
		Field("Status", "OpenflowConnectorStatus").
		Text("Runtime").
		OptionalText("ConnectorDefinition").
		OptionalText("DisplayName").
		Text("DatabaseName").
		Text("SchemaName").
		Text("Owner").
		OptionalText("Comment").
		Time("CreatedOn").
		Time("UpdatedOn"),
	g.NewQueryStruct("ShowOpenflowConnectors").
		Show().
		SQL("OPENFLOW CONNECTORS").
		OptionalLike(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-openflow-connector",
	g.DbStruct("openflowConnectorDetailsRow").
		Text("name").
		Text("status").
		Text("runtime").
		OptionalText("connector_definition").
		OptionalText("display_name").
		Text("database_name").
		Text("schema_name").
		Text("owner").
		OptionalText("comment").
		Time("created_on").
		Time("updated_on").
		OptionalText("error_code").
		OptionalText("status_message"),
	g.PlainStruct("OpenflowConnectorDetails").
		Field("Status", "OpenflowConnectorStatus").
		Text("Runtime").
		OptionalText("ConnectorDefinition").
		OptionalText("DisplayName").
		Text("DatabaseName").
		Text("SchemaName").
		Text("Owner").
		OptionalText("Comment").
		Time("CreatedOn").
		Time("UpdatedOn").
		OptionalText("ErrorCode").
		OptionalText("StatusMessage"),
	g.NewQueryStruct("DescribeOpenflowConnector").
		Describe().
		SQL("OPENFLOW CONNECTOR").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
