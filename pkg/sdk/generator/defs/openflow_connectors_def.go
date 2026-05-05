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
	"TODO: add link when public docs are available",
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
	"TODO: add link when public docs are available",
	g.NewQueryStruct("AlterOpenflowConnector").
		Alter().
		SQL("OPENFLOW CONNECTOR").
		Name().
		OptionalSQL("START").
		OptionalSQL("STOP").
		OptionalSQL("TERMINATE").
		OptionalSQL("COMMIT").
		OptionalSQL("ABORT").
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
		WithValidation(g.ExactlyOneValueSet, "Start", "Stop", "Terminate", "Commit", "Abort", "Set", "Unset"),
).DropOperation(
	"TODO: add link when public docs are available",
	g.NewQueryStruct("DropOpenflowConnector").
		Drop().
		SQL("OPENFLOW CONNECTOR").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"TODO: add link when public docs are available",
	g.DbStruct("openflowConnectorRow").
		Text("name").
		Text("status").
		Text("runtime").
		OptionalText("connector_definition").
		OptionalText("display_name").
		Text("database_name").
		Text("schema_name").
		Text("owner").
		OptionalText("default_version").
		OptionalText("default_version_name").
		OptionalText("default_version_alias").
		OptionalText("default_version_location_uri").
		OptionalText("default_version_source_location_uri").
		OptionalText("live_version_location_uri").
		OptionalText("comment").
		Time("created_on").
		Time("updated_on"),
	g.PlainStruct("OpenflowConnector").
		Text("Name").
		Field("Status", "OpenflowConnectorStatus").
		Text("Runtime").
		OptionalText("ConnectorDefinition").
		OptionalText("DisplayName").
		Text("DatabaseName").
		Text("SchemaName").
		Text("Owner").
		OptionalText("DefaultVersion").
		OptionalText("DefaultVersionName").
		OptionalText("DefaultVersionAlias").
		OptionalText("DefaultVersionLocationUri").
		OptionalText("DefaultVersionSourceLocationUri").
		OptionalText("LiveVersionLocationUri").
		OptionalText("Comment").
		Time("CreatedOn").
		Time("UpdatedOn"),
	g.NewQueryStruct("ShowOpenflowConnectors").
		Show().
		SQL("OPENFLOW CONNECTORS").
		OptionalLike().
		OptionalIn(),
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"TODO: add link when public docs are available",
	g.DbStruct("openflowConnectorDetailsRow").
		Text("name").
		Text("status").
		Text("runtime").
		OptionalText("connector_definition").
		OptionalText("definition_version_name").
		OptionalText("provider").
		OptionalText("display_name").
		Text("database_name").
		Text("schema_name").
		Text("owner").
		OptionalText("default_version").
		OptionalText("default_version_name").
		OptionalText("default_version_alias").
		OptionalText("default_version_location_uri").
		OptionalText("default_version_source_location_uri").
		OptionalText("default_version_git_commit_hash").
		OptionalText("last_version_name").
		OptionalText("last_version_alias").
		OptionalText("last_version_location_uri").
		OptionalText("last_version_source_location_uri").
		OptionalText("last_version_git_commit_hash").
		OptionalText("live_version_location_uri").
		OptionalText("comment").
		Time("created_on").
		Time("updated_on").
		OptionalText("error_code").
		OptionalText("status_message"),
	g.PlainStruct("OpenflowConnectorDetails").
		Text("Name").
		Field("Status", "OpenflowConnectorStatus").
		Text("Runtime").
		OptionalText("ConnectorDefinition").
		OptionalText("DefinitionVersionName").
		OptionalText("Provider").
		OptionalText("DisplayName").
		Text("DatabaseName").
		Text("SchemaName").
		Text("Owner").
		OptionalText("DefaultVersion").
		OptionalText("DefaultVersionName").
		OptionalText("DefaultVersionAlias").
		OptionalText("DefaultVersionLocationUri").
		OptionalText("DefaultVersionSourceLocationUri").
		OptionalText("DefaultVersionGitCommitHash").
		OptionalText("LastVersionName").
		OptionalText("LastVersionAlias").
		OptionalText("LastVersionLocationUri").
		OptionalText("LastVersionSourceLocationUri").
		OptionalText("LastVersionGitCommitHash").
		OptionalText("LiveVersionLocationUri").
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
