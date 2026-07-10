package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var OpenflowRuntimeNodeTypeEnumDef = g.NewEnum(
	"OpenflowRuntimeNodeType", "OpenflowRuntimeNodeTypes",
	"SMALL", "MEDIUM", "LARGE",
)

var OpenflowRuntimeStatusEnumDef = g.NewEnum(
	"OpenflowRuntimeStatus", "OpenflowRuntimeStatuses",
	"CREATING", "CREATE_FAILED", "UPDATING", "UPDATE_FAILED",
	"SUSPENDING", "SUSPENDED", "SUSPEND_FAILED",
	"ACTIVATING", "ACTIVE", "ACTIVATE_FAILED",
	"DELETING", "DELETED", "DELETE_FAILED",
	"CANCEL_REQUESTED", "RESTARTING", "RESTART_FAILED",
	"UPGRADING", "UPGRADE_FAILED",
	"GENERATING_DIAGNOSTIC_BUNDLE", "CLEANING_UP", "INACTIVE",
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
		Identifier("ExecuteAsRole", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXECUTE_AS_ROLE").Equals().Required()).
		Assignment("NODE_TYPE", OpenflowRuntimeNodeTypeEnumDef.Kind(), g.ParameterOptions().SingleQuotes().Required()).
		NumberAssignment("MIN_NODES", g.ParameterOptions().Required()).
		NumberAssignment("MAX_NODES", g.ParameterOptions().Required()).
		OptionalQueryStructField("ExternalAccessIntegrations", openflowRuntimesExternalAccessIntegrationsDef, g.ParameterOptions().SQL("EXTERNAL_ACCESS_INTEGRATIONS").Parentheses()).
		OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "InDeployment").
		WithValidation(g.ValidIdentifier, "ExecuteAsRole"),
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
		RenameTo().
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("OpenflowRuntimeSet").
				OptionalNumberAssignment("MIN_NODES", g.ParameterOptions()).
				OptionalNumberAssignment("MAX_NODES", g.ParameterOptions()).
				OptionalIdentifier("ExecuteAsRole", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXECUTE_AS_ROLE").Equals()).
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
).ShowOperationWithPairedStructs(
	"TODO: add link when public docs are available",
	g.StructPair("openflowRuntimeRow", "OpenflowRuntime").
		Text("name").
		Enum("status", OpenflowRuntimeStatusEnumDef).
		Text("deployment").
		Number("min_nodes").
		Number("max_nodes").
		Enum("node_type", OpenflowRuntimeNodeTypeEnumDef).
		OptionalText("display_name").
		Field("external_access_integrations", "sql.NullString", "[]AccountObjectIdentifier").
		Bool("initially_suspended").
		Text("execute_as_role").
		Text("owner").
		OptionalText("comment").
		Time("created_on").
		Time("updated_on"),
	g.NewQueryStruct("ShowOpenflowRuntimes").
		Show().
		SQL("OPENFLOW RUNTIMES").
		OptionalLike().
		OptionalIn(),
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"TODO: add link when public docs are available",
	g.StructPair("openflowRuntimeDetailsRow", "OpenflowRuntimeDetails").
		Text("name").
		Enum("status", OpenflowRuntimeStatusEnumDef).
		Text("deployment").
		Number("min_nodes").
		Number("max_nodes").
		Enum("node_type", OpenflowRuntimeNodeTypeEnumDef).
		OptionalText("display_name").
		Field("external_access_integrations", "sql.NullString", "[]AccountObjectIdentifier").
		Bool("initially_suspended").
		Text("execute_as_role").
		Text("owner").
		OptionalText("comment").
		OptionalText("server_url").
		Time("created_on").
		Time("updated_on").
		OptionalText("error_code").
		OptionalText("status_message"),
	g.NewQueryStruct("DescribeOpenflowRuntime").
		Describe().
		SQL("OPENFLOW RUNTIME").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).WithEnums(
	OpenflowRuntimeNodeTypeEnumDef,
	OpenflowRuntimeStatusEnumDef,
)
