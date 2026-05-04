//go:build sdk_generation

package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var openflowRuntimesExternalAccessIntegrationsDef = g.NewQueryStruct("OpenflowRuntimeExternalAccessIntegrations").
	List("ExternalAccessIntegrations", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.ListOptions().Required().MustParentheses())

var openflowRuntimesDef = g.NewInterface(
	"OpenflowRuntimes",
	"OpenflowRuntime",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"TODO: add link when public docs are available",
	g.NewQueryStruct("CreateOpenflowRuntime").
		Create().
		SQL("OPENFLOW RUNTIME").
		IfNotExists().
		Name().
		Identifier("InDeployment", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("IN DEPLOYMENT").Required()).
		TextAssignment("EXECUTE_AS_ROLE", g.ParameterOptions().NoQuotes().Required()).
		Assignment("NODE_TYPE", g.KindOfT[sdkcommons.OpenflowRuntimeNodeType](), g.ParameterOptions().SingleQuotes().Required()).
		NumberAssignment("MIN_NODES", g.ParameterOptions().Required()).
		NumberAssignment("MAX_NODES", g.ParameterOptions().Required()).
		OptionalQueryStructField("ExternalAccessIntegrations", openflowRuntimesExternalAccessIntegrationsDef, g.ParameterOptions().SQL("EXTERNAL_ACCESS_INTEGRATIONS").Parentheses()).
		OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"TODO: add link when public docs are available",
	g.NewQueryStruct("AlterOpenflowRuntime").
		Alter().
		SQL("OPENFLOW RUNTIME").
		Name().
		OptionalSQL("SUSPEND").
		OptionalSQL("RESUME").
		OptionalSQL("RESUME RECOVERY").
		OptionalSQL("RESTART").
		OptionalSQL("RESTART RECOVERY").
		OptionalSQL("TERMINATE").
		OptionalSQL("TERMINATE CASCADE").
		OptionalSQL("UPGRADE").
		OptionalIdentifier("RenameTo", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("OpenflowRuntimeSet").
				OptionalNumberAssignment("MIN_NODES", g.ParameterOptions()).
				OptionalNumberAssignment("MAX_NODES", g.ParameterOptions()).
				OptionalTextAssignment("EXECUTE_AS_ROLE", g.ParameterOptions().NoQuotes()).
				OptionalQueryStructField("ExternalAccessIntegrations", openflowRuntimesExternalAccessIntegrationsDef, g.ParameterOptions().SQL("EXTERNAL_ACCESS_INTEGRATIONS").Parentheses()).
				OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.AtLeastOneValueSet, "MinNodes", "MaxNodes", "ExecuteAsRole", "ExternalAccessIntegrations", "DisplayName", "Comment"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("OpenflowRuntimeUnset").
				OptionalSQL("EXECUTE_AS_ROLE").
				OptionalSQL("EXTERNAL_ACCESS_INTEGRATIONS").
				OptionalSQL("DISPLAY_NAME").
				OptionalSQL("COMMENT").
				WithValidation(g.AtLeastOneValueSet, "ExecuteAsRole", "ExternalAccessIntegrations", "DisplayName", "Comment"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "Suspend", "Resume", "ResumeRecovery", "Restart", "RestartRecovery", "Terminate", "TerminateCascade", "Upgrade", "RenameTo", "Set", "Unset"),
).DropOperation(
	"TODO: add link when public docs are available",
	g.NewQueryStruct("DropOpenflowRuntime").
		Drop().
		SQL("OPENFLOW RUNTIME").
		IfExists().
		Name().
		OptionalSQL("CASCADE").
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"TODO: add link when public docs are available",
	g.DbStruct("openflowRuntimeRow").
		Text("name").
		Text("status").
		Text("deployment").
		Number("min_nodes").
		Number("max_nodes").
		Text("node_type").
		OptionalText("display_name").
		OptionalText("external_access_integrations").
		Bool("initially_suspended").
		Text("database_name").
		Text("schema_name").
		Text("execute_as_role").
		Text("owner").
		OptionalText("comment").
		OptionalText("server_url").
		Time("created_on").
		Time("updated_on"),
	g.PlainStruct("OpenflowRuntime").
		Field("Status", "OpenflowRuntimeStatus").
		Text("Deployment").
		Number("MinNodes").
		Number("MaxNodes").
		Field("NodeType", "OpenflowRuntimeNodeType").
		OptionalText("DisplayName").
		OptionalText("ExternalAccessIntegrations").
		Bool("InitiallySuspended").
		Text("DatabaseName").
		Text("SchemaName").
		Text("ExecuteAsRole").
		Text("Owner").
		OptionalText("Comment").
		OptionalText("ServerUrl").
		Time("CreatedOn").
		Time("UpdatedOn"),
	g.NewQueryStruct("ShowOpenflowRuntimes").
		Show().
		SQL("OPENFLOW RUNTIMES").
		OptionalLike(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"TODO: add link when public docs are available",
	g.DbStruct("openflowRuntimeDetailsRow").
		Text("name").
		Text("status").
		Text("deployment").
		Number("min_nodes").
		Number("max_nodes").
		Text("node_type").
		OptionalText("display_name").
		OptionalText("external_access_integrations").
		Bool("initially_suspended").
		Text("database_name").
		Text("schema_name").
		Text("execute_as_role").
		Text("owner").
		OptionalText("comment").
		OptionalText("server_url").
		Time("created_on").
		Time("updated_on").
		OptionalText("error_code").
		OptionalText("status_message"),
	g.PlainStruct("OpenflowRuntimeDetails").
		Field("Status", "OpenflowRuntimeStatus").
		Text("Deployment").
		Number("MinNodes").
		Number("MaxNodes").
		Field("NodeType", "OpenflowRuntimeNodeType").
		OptionalText("DisplayName").
		OptionalText("ExternalAccessIntegrations").
		Bool("InitiallySuspended").
		Text("DatabaseName").
		Text("SchemaName").
		Text("ExecuteAsRole").
		Text("Owner").
		OptionalText("Comment").
		OptionalText("ServerUrl").
		Time("CreatedOn").
		Time("UpdatedOn").
		OptionalText("ErrorCode").
		OptionalText("StatusMessage"),
	g.NewQueryStruct("DescribeOpenflowRuntime").
		Describe().
		SQL("OPENFLOW RUNTIME").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
