package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var OpenflowRuntimeNodeTypeEnumDef = g.NewEnum(
	"OpenflowRuntimeNodeType", "OpenflowRuntimeNodeTypes",
	"SMALL", "MEDIUM", "LARGE",
)

// OpenflowRuntimeStatusEnumDef covers both the legacy status names and the newer aliases, since
// Snowflake can return either spelling. The aliases below fold onto the legacy names, so callers only
// ever see one spelling.
var OpenflowRuntimeStatusEnumDef = g.NewEnum(
	"OpenflowRuntimeStatus", "OpenflowRuntimeStatuses",
	"CREATING", "CREATE_FAILED", "UPDATING", "UPDATE_FAILED",
	"SUSPENDING", "SUSPENDED", "SUSPEND_FAILED",
	"ACTIVATING", "ACTIVE", "ACTIVATE_FAILED",
	"DELETING", "DELETED", "DELETE_FAILED",
	"CANCEL_REQUESTED", "RESTARTING", "RESTART_FAILED",
	"UPGRADING", "UPGRADE_FAILED",
	"GENERATING_DIAGNOSTIC_BUNDLE", "CLEANING_UP", "INACTIVE",
	// These have no legacy equivalent and are always returned as-is, so they need no alias.
	"MIGRATING", "MIGRATION_FAILED", "ROLLING_BACK", "ROLLBACK_FAILED",
).WithAliases("DELETING", "TERMINATING").
	WithAliases("DELETED", "TERMINATED").
	WithAliases("DELETE_FAILED", "TERMINATE_FAILED").
	WithAliases("ACTIVATING", "RESUMING").
	WithAliases("ACTIVATE_FAILED", "RESUME_FAILED")

var openflowRuntimesExternalAccessIntegrationsDef = g.NewQueryStruct("OpenflowRuntimeExternalAccessIntegrations").
	List("ExternalAccessIntegrations", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.ListOptions().Required().MustParentheses())

// openflowRuntimeUpgradeDef models `UPGRADE [ RECOVERY ] [ FORCE ]`. RECOVERY and FORCE are independent
// optional modifiers rather than alternatives: Snowflake accepts UPGRADE, UPGRADE RECOVERY, UPGRADE FORCE
// and UPGRADE RECOVERY FORCE. The order is fixed, and UPGRADE FORCE RECOVERY is a parse error, which is
// why they are separate fields in this order rather than one enum.
func openflowRuntimeUpgradeDef() *g.QueryStruct {
	return g.NewQueryStruct("OpenflowRuntimeUpgrade").
		OptionalSQL("RECOVERY").
		OptionalSQL("FORCE")
}

var openflowRuntimesDef = g.NewInterface(
	"OpenflowRuntimes",
	"OpenflowRuntime",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-runtime#create-openflow-runtime",
	g.NewQueryStruct("CreateOpenflowRuntime").
		Create().
		SQL("OPENFLOW RUNTIME").
		IfNotExists().
		Name().
		// Clause order follows the published CREATE OPENFLOW RUNTIME syntax exactly - IN DEPLOYMENT,
		// NODE_TYPE, MIN_NODES, MAX_NODES, EXECUTE_AS_ROLE - because option order has broken queries here
		// before.
		Identifier("InDeployment", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("IN DEPLOYMENT").Required()).
		Assignment("NODE_TYPE", OpenflowRuntimeNodeTypeEnumDef.Kind(), g.ParameterOptions().SingleQuotes().Required()).
		NumberAssignment("MIN_NODES", g.ParameterOptions().Required()).
		NumberAssignment("MAX_NODES", g.ParameterOptions().Required()).
		Identifier("ExecuteAsRole", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXECUTE_AS_ROLE").Equals().Required()).
		OptionalQueryStructField("ExternalAccessIntegrations", openflowRuntimesExternalAccessIntegrationsDef, g.ParameterOptions().SQL("EXTERNAL_ACCESS_INTEGRATIONS").Parentheses()).
		OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "InDeployment").
		WithValidation(g.ValidIdentifier, "ExecuteAsRole"),
).AlterOperation(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-runtime#alter-openflow-runtime",
	g.NewQueryStruct("AlterOpenflowRuntime").
		Alter().
		SQL("OPENFLOW RUNTIME").
		IfExists().
		Name().
		OptionalSQL("SUSPEND").
		OptionalSQL("RESUME").
		OptionalSQL("RESUME RECOVERY").
		OptionalSQL("RESTART").
		OptionalSQL("RESTART RECOVERY").
		OptionalSQL("TERMINATE").
		OptionalSQL("TERMINATE CASCADE").
		OptionalQueryStructField("Upgrade", openflowRuntimeUpgradeDef(), g.KeywordOptions().SQL("UPGRADE")).
		RenameTo().
		OptionalQueryStructField(
			"Set",
			// Order follows the published ALTER OPENFLOW RUNTIME ... SET syntax.
			g.NewQueryStruct("OpenflowRuntimeSet").
				OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
				OptionalNumberAssignment("MIN_NODES", g.ParameterOptions()).
				OptionalNumberAssignment("MAX_NODES", g.ParameterOptions()).
				OptionalQueryStructField("ExternalAccessIntegrations", openflowRuntimesExternalAccessIntegrationsDef, g.ParameterOptions().SQL("EXTERNAL_ACCESS_INTEGRATIONS").Parentheses()).
				OptionalIdentifier("ExecuteAsRole", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXECUTE_AS_ROLE").Equals()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.AtLeastOneValueSet, "MinNodes", "MaxNodes", "ExecuteAsRole", "ExternalAccessIntegrations", "DisplayName", "Comment"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			// EXECUTE_AS_ROLE has no UNSET: a runtime always runs as some role, and the server rejects the
			// attempt. It can only be moved with SET.
			g.NewQueryStruct("OpenflowRuntimeUnset").
				OptionalSQL("EXTERNAL_ACCESS_INTEGRATIONS").
				OptionalSQL("DISPLAY_NAME").
				OptionalSQL("COMMENT").
				WithValidation(g.AtLeastOneValueSet, "ExternalAccessIntegrations", "DisplayName", "Comment"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		// ADD and REMOVE edit the external access integration list in place, where SET replaces it
		// wholesale. Both take the same parenthesised, `=`-assigned list as SET.
		OptionalQueryStructField("AddExternalAccessIntegrations", openflowRuntimesExternalAccessIntegrationsDef, g.ParameterOptions().SQL("ADD EXTERNAL_ACCESS_INTEGRATIONS").Parentheses()).
		OptionalQueryStructField("RemoveExternalAccessIntegrations", openflowRuntimesExternalAccessIntegrationsDef, g.ParameterOptions().SQL("REMOVE EXTERNAL_ACCESS_INTEGRATIONS").Parentheses()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "Suspend", "Resume", "ResumeRecovery", "Restart", "RestartRecovery", "Terminate", "TerminateCascade", "Upgrade", "RenameTo", "Set", "Unset", "AddExternalAccessIntegrations", "RemoveExternalAccessIntegrations"),
).DropOperation(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-runtime#drop-openflow-runtime",
	g.NewQueryStruct("DropOpenflowRuntime").
		Drop().
		SQL("OPENFLOW RUNTIME").
		IfExists().
		Name().
		OptionalSQL("CASCADE").
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-runtime#show-openflow-runtimes",
	g.StructPair("openflowRuntimeRow", "OpenflowRuntime").
		Text("name").
		Enum("status", OpenflowRuntimeStatusEnumDef).
		Text("deployment").
		Number("min_nodes").
		Number("max_nodes").
		Enum("node_type", OpenflowRuntimeNodeTypeEnumDef).
		OptionalText("display_name").
		Field("external_access_integrations", "sql.NullString", "[]AccountObjectIdentifier", g.WithCustomParser("ParseOpenflowRuntimeExternalAccessIntegrations")).
		Bool("initially_suspended").
		Text("database_name").
		Text("schema_name").
		OptionalText("execute_as_role").
		OptionalText("key").
		Text("owner").
		OptionalText("comment").
		Time("created_on").
		Time("updated_on"),
	g.NewQueryStruct("ShowOpenflowRuntimes").
		Show().
		SQL("OPENFLOW RUNTIMES").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalLimitFrom(),
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-runtime#describe-openflow-runtime",
	g.StructPair("openflowRuntimeDetailsRow", "OpenflowRuntimeDetails").
		// DESCRIBE OPENFLOW RUNTIME returns neither database_name nor schema_name (unlike SHOW), so the
		// identifier cannot be rebuilt from the row. Id is populated by the caller, as for NotebookDetails.
		PlainOnlyField("Id", "SchemaObjectIdentifier").
		Text("name").
		Enum("status", OpenflowRuntimeStatusEnumDef).
		Text("deployment").
		Number("min_nodes").
		Number("max_nodes").
		Enum("node_type", OpenflowRuntimeNodeTypeEnumDef).
		OptionalText("display_name").
		Field("external_access_integrations", "sql.NullString", "[]AccountObjectIdentifier", g.WithCustomParser("ParseOpenflowRuntimeExternalAccessIntegrations")).
		Bool("initially_suspended").
		OptionalText("execute_as_role").
		OptionalText("key").
		Text("owner").
		OptionalText("comment").
		OptionalText("server_url").
		// Optional defensively, not because the column is gated: a pointer costs nothing and avoids the
		// whole DESCRIBE failing if node_type_tier is ever NULL.
		OptionalText("node_type_tier"),
	g.NewQueryStruct("DescribeOpenflowRuntime").
		Describe().
		SQL("OPENFLOW RUNTIME").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).WithEnums(
	OpenflowRuntimeNodeTypeEnumDef,
	OpenflowRuntimeStatusEnumDef,
)
