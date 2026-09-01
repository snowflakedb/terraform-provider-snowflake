package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// OpenflowConnectorStatusEnumDef covers both the legacy status names and the newer aliases, since
// Snowflake can return either spelling. The aliases below fold onto the legacy names, so callers only
// ever see one spelling.
var OpenflowConnectorStatusEnumDef = g.NewEnum(
	"OpenflowConnectorStatus", "OpenflowConnectorStatuses",
	"CREATING", "CREATE_FAILED",
	"STARTING", "START_FAILED",
	"RUNNING",
	"STOPPING", "STOP_FAILED", "STOPPED",
	"UPDATING", "UPDATE_FAILED",
	"DELETING", "DELETE_FAILED", "DELETED",
	// The troubleshooting lifecycle. These have no legacy equivalent and are always returned as-is, so
	// they need no alias. SYSTEM$OPENFLOW_CONNECTOR_EXIT_TROUBLESHOOTING requires the connector to
	// already be in one of them, so a customer can observe them.
	"ENTERING_TROUBLESHOOTING", "ENTER_TROUBLESHOOTING_FAILED", "TROUBLESHOOTING",
	"EXITING_TROUBLESHOOTING", "EXIT_TROUBLESHOOTING_FAILED",
).WithAliases("DELETING", "TERMINATING").
	WithAliases("DELETED", "TERMINATED").
	WithAliases("DELETE_FAILED", "TERMINATE_FAILED")

// openflowConnectorAddLiveVersionDef models `ADD LIVE VERSION [ <alias> ] FROM LAST [ COMMENT = '<c>' ]`.
// Creating a live version is the first step of editing a committed connector's configuration: add the
// live version, PUT the edited files onto its stage, then COMMIT. Only `FROM LAST` is currently
// supported, so it is emitted as a static clause.
func openflowConnectorAddLiveVersionDef() *g.QueryStruct {
	return g.NewQueryStruct("OpenflowConnectorAddLiveVersion").
		OptionalText("VersionAlias", g.KeywordOptions().NoQuotes()).
		SQL("FROM LAST").
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes())
}

// openflowConnectorCommitDef models `COMMIT [ COMMENT = '<c>' ]`, which applies the live version and
// removes it.
func openflowConnectorCommitDef() *g.QueryStruct {
	return g.NewQueryStruct("OpenflowConnectorCommit").
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes())
}

// openflowConnectorAddVersionDef models `ADD VERSION FROM '<stage_reference>'`. Unlike ADD LIVE VERSION,
// which opens an editable live version that later needs a COMMIT, this imports a configuration from a
// stage and promotes it straight to the new default in one statement. The connector must be STOPPED or
// in Draft. The location is a *Location so both `@<stage>/<path>` and `snow://openflow_connector/...`
// are expressible.
func openflowConnectorAddVersionDef() *g.QueryStruct {
	return g.NewQueryStruct("OpenflowConnectorAddVersion").
		PredefinedQueryStructField("From", "*Location", g.ParameterOptions().SQL("FROM").SingleQuotes().NoEquals().Required()).
		// Required() makes From a constructor parameter on the request, but the Options struct can still be
		// built directly with a nil From, which would emit `ADD VERSION` with no FROM at all.
		WithValidation(g.ValidateValueSet, "From")
}

// openflowConnectorPushDef models
// `PUSH [ TO '<git_branch_uri>' ] USERNAME = <u> PASSWORD = <p> NAME = <n> EMAIL = <e> [ COMMENT = <c> ]`,
// which commits the connector's configuration back to a git branch. The four credential and authorship
// fields are all required by the published syntax; only TO and COMMENT are optional.
//
// PASSWORD necessarily appears in the statement text, so it reaches Snowflake's query history. The
// grammar has no secret-reference form. Any resource exposing PUSH must mark its credential attributes
// Sensitive so Terraform keeps them out of plan output and state display.
func openflowConnectorPushDef() *g.QueryStruct {
	return g.NewQueryStruct("OpenflowConnectorPush").
		OptionalTextAssignment("TO", g.ParameterOptions().SingleQuotes().NoEquals()).
		TextAssignment("USERNAME", g.ParameterOptions().SingleQuotes().Required()).
		TextAssignment("PASSWORD", g.ParameterOptions().SingleQuotes().Required()).
		TextAssignment("NAME", g.ParameterOptions().SingleQuotes().Required()).
		TextAssignment("EMAIL", g.ParameterOptions().SingleQuotes().Required()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes())
}

// openflowConnectorDefaultVersionDef models the value of `SET DEFAULT_VERSION = { FIRST | LAST | '<v>' }`.
// FIRST and LAST are bare keywords while a specific version is single-quoted, so one text field cannot
// express all three.
//
// Version is the Snowflake-generated version name, VERSION$<n>, not a version alias.
func openflowConnectorDefaultVersionDef() *g.QueryStruct {
	return g.NewQueryStruct("OpenflowConnectorDefaultVersion").
		OptionalSQL("FIRST").
		OptionalSQL("LAST").
		OptionalText("Version", g.KeywordOptions().SingleQuotes()).
		WithValidation(g.ExactlyOneValueSet, "First", "Last", "Version")
}

var openflowConnectorsDef = g.NewInterface(
	"OpenflowConnectors",
	"OpenflowConnector",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-connector#create-openflow-connector",
	g.NewQueryStruct("CreateOpenflowConnector").
		Create().
		SQL("OPENFLOW CONNECTOR").
		IfNotExists().
		Name().
		Identifier("InRuntime", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("IN RUNTIME").Required()).
		OptionalTextAssignment("FROM DEFINITION", g.ParameterOptions().NoQuotes().NoEquals()).
		PredefinedQueryStructField("From", "*Location", g.ParameterOptions().SQL("FROM").SingleQuotes().NoEquals()).
		OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "InRuntime").
		WithValidation(g.ConflictingFields, "FromDefinition", "From"),
).AlterOperation(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-connector#alter-openflow-connector",
	g.NewQueryStruct("AlterOpenflowConnector").
		Alter().
		SQL("OPENFLOW CONNECTOR").
		IfExists().
		Name().
		OptionalSQL("START").
		OptionalSQL("STOP").
		OptionalSQL("TERMINATE").
		// TERMINATE FORCE is the escape hatch for a connector stuck in TERMINATING; plain TERMINATE is
		// rejected from that state.
		OptionalSQL("TERMINATE FORCE").
		OptionalSQL("ABORT").
		// RENAME TO cannot move a connector between schemas.
		RenameTo().
		OptionalQueryStructField("AddVersion", openflowConnectorAddVersionDef(), g.KeywordOptions().SQL("ADD VERSION")).
		OptionalQueryStructField("AddLiveVersion", openflowConnectorAddLiveVersionDef(), g.KeywordOptions().SQL("ADD LIVE VERSION")).
		OptionalQueryStructField("Commit", openflowConnectorCommitDef(), g.KeywordOptions().SQL("COMMIT")).
		OptionalQueryStructField("Push", openflowConnectorPushDef(), g.KeywordOptions().SQL("PUSH")).
		OptionalSQL("PULL").
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("OpenflowConnectorSet").
				OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				OptionalQueryStructField("DefaultVersion", openflowConnectorDefaultVersionDef(), g.ParameterOptions().SQL("DEFAULT_VERSION")).
				WithValidation(g.AtLeastOneValueSet, "DisplayName", "Comment", "DefaultVersion"),
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
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "Start", "Stop", "Terminate", "TerminateForce", "Abort", "RenameTo", "AddVersion", "AddLiveVersion", "Commit", "Push", "Pull", "Set", "Unset").
		// IF EXISTS is only grammatical with a subset of the actions. Every version and git action rejects
		// it with 001003 (42000) syntax error, unexpected '<action>', whereas STOP and TERMINATE are
		// accepted and no-op. Rejecting the bad combinations here turns a
		// server-side syntax error into a client-side one. Pairwise because ConflictingFields errors only
		// when every listed field is set.
		WithValidation(g.ConflictingFields, "IfExists", "Commit").
		WithValidation(g.ConflictingFields, "IfExists", "Abort").
		WithValidation(g.ConflictingFields, "IfExists", "AddLiveVersion").
		WithValidation(g.ConflictingFields, "IfExists", "AddVersion").
		WithValidation(g.ConflictingFields, "IfExists", "Push").
		WithValidation(g.ConflictingFields, "IfExists", "Pull"),
).DropOperation(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-connector#drop-openflow-connector",
	g.NewQueryStruct("DropOpenflowConnector").
		Drop().
		SQL("OPENFLOW CONNECTOR").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-connector#show-openflow-connectors",
	g.StructPair("openflowConnectorRow", "OpenflowConnector").
		Text("name").
		Enum("status", OpenflowConnectorStatusEnumDef).
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
		Time("updated_on").
		OptionalText("connector_url"),
	g.NewQueryStruct("ShowOpenflowConnectors").
		Show().
		SQL("OPENFLOW CONNECTORS").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalLimitFrom(),
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-connector#describe-openflow-connector",
	g.StructPair("openflowConnectorDetailsRow", "OpenflowConnectorDetails").
		// DESCRIBE OPENFLOW CONNECTOR returns neither database_name nor schema_name (unlike SHOW), so the
		// identifier cannot be rebuilt from the row. Id is populated by the caller, as for NotebookDetails.
		PlainOnlyField("Id", "SchemaObjectIdentifier").
		Text("name").
		Enum("status", OpenflowConnectorStatusEnumDef).
		Text("runtime").
		OptionalText("connector_definition").
		OptionalText("display_name").
		Text("owner").
		OptionalText("comment").
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
		OptionalText("connector_url"),
	g.NewQueryStruct("DescribeOpenflowConnector").
		Describe().
		SQL("OPENFLOW CONNECTOR").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	// EXECUTE is its own statement rather than an ALTER action, so it needs a custom operation. It
	// validates a connector's configuration without applying it: against the connector's own current
	// configuration by default, or against the one at FROM. STEP narrows validation to a single named
	// step.
	"Execute",
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-connector#execute-openflow-connector",
	g.NewQueryStruct("ExecuteOpenflowConnector").
		SQL("EXECUTE").
		SQL("OPENFLOW CONNECTOR").
		Name().
		SQL("VALIDATE CONFIGURATION").
		PredefinedQueryStructField("From", "*Location", g.ParameterOptions().SQL("FROM").SingleQuotes().NoEquals()).
		OptionalTextAssignment("STEP", g.ParameterOptions().SingleQuotes().NoEquals()).
		WithValidation(g.ValidIdentifier, "name"),
).CustomShowOperationWithPairedStructs(
	// SHOW VERSIONS IN OPENFLOW CONNECTOR lists a connector's configuration versions. Its row shape differs
	// from SHOW OPENFLOW CONNECTORS, which is what a custom show operation is for - the same arrangement
	// SHOW GIT BRANCHES and SHOW GIT TAGS use on GitRepositories.
	"ShowVersions",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-connector",
	// name is optional because a version that has been created but not yet committed has none.
	g.StructPair("openflowConnectorVersionRow", "OpenflowConnectorVersion").
		OptionalText("name").
		OptionalText("alias").
		OptionalText("comment").
		Time("created_on").
		Bool("is_default").
		Bool("is_first").
		Bool("is_last").
		Bool("is_live").
		Text("location_uri").
		OptionalText("source_location_uri").
		OptionalText("git_commit_hash"),
	g.NewQueryStruct("ShowVersionsOpenflowConnector").
		SQL("SHOW VERSIONS").
		SQL("IN OPENFLOW CONNECTOR").
		Name().
		// Only a bare LIMIT <n>, so this cannot reuse OptionalLimitFrom the way the four Openflow SHOW
		// commands do. LIKE, STARTS WITH and LIMIT <n> FROM '<s>' are all parse errors on this command.
		PredefinedQueryStructField("Limit", "*int", g.ParameterOptions().NoEquals().SQL("LIMIT")).
		WithValidation(g.ValidIdentifier, "name"),
).WithEnums(
	OpenflowConnectorStatusEnumDef,
)
