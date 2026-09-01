package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var OpenflowDeploymentTypeEnumDef = g.NewEnum(
	"OpenflowDeploymentType", "OpenflowDeploymentTypes",
	"SNOWFLAKE", "BYOC",
)

var OpenflowVpcTypeEnumDef = g.NewEnum(
	"OpenflowVpcType", "OpenflowVpcTypes",
	"MANAGED", "PROVIDED",
)

// OpenflowDeploymentStatusEnumDef covers both the legacy status names and the newer aliases, since
// Snowflake can return either spelling. The aliases below fold TERMINATING, TERMINATED,
// TERMINATE_FAILED and PROVISIONING_FAILED onto DELETING, DELETED, DELETE_FAILED and CREATE_FAILED,
// so callers only ever see one spelling.
var OpenflowDeploymentStatusEnumDef = g.NewEnum(
	"OpenflowDeploymentStatus", "OpenflowDeploymentStatuses",
	"CREATING", "ACTIVE", "INACTIVE", "PROVISIONING", "NOT_REPORTING",
	"NOT_HEALTHY", "UPGRADING", "UPGRADE_FAILED", "DEACTIVATION_REQUIRED",
	"DELETING", "DELETED", "CREATE_FAILED", "DELETE_FAILED",
	// These have no legacy equivalent and are always returned as-is, so they need no alias.
	// GENERATING_DIAGNOSTIC_BUNDLE is included for parity with the runtime enum, which already had it.
	"MIGRATING", "MIGRATION_FAILED", "ROLLING_BACK", "ROLLBACK_FAILED",
	"GENERATING_DIAGNOSTIC_BUNDLE",
).WithAliases("DELETING", "TERMINATING").
	WithAliases("DELETED", "TERMINATED").
	WithAliases("DELETE_FAILED", "TERMINATE_FAILED").
	WithAliases("CREATE_FAILED", "PROVISIONING_FAILED")

// openflowDeploymentEventTableDef models `EVENT_TABLE = { '<db>.<schema>.<table>' | NONE }`. The table is
// an identifier rather than free text, and NONE is a bare keyword, so the two alternatives cannot share a
// field.
//
// NONE is not the same as omitting the clause. NONE drops all events for the deployment, whereas leaving
// it unset - or UNSET on ALTER - falls back to the account's default event table.
func openflowDeploymentEventTableDef() *g.QueryStruct {
	return g.NewQueryStruct("OpenflowDeploymentEventTable").
		OptionalIdentifier("EventTable", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SingleQuotes()).
		OptionalSQL("NONE").
		WithValidation(g.ExactlyOneValueSet, "EventTable", "None")
}

var openflowDeploymentsDef = g.NewInterface(
	"OpenflowDeployments",
	"OpenflowDeployment",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-deployment#create-openflow-deployment",
	g.NewQueryStruct("CreateOpenflowDeployment").
		Create().
		SQL("OPENFLOW DEPLOYMENT").
		IfNotExists().
		Name().
		// Clause order follows the published CREATE OPENFLOW DEPLOYMENT syntax exactly, since option order
		// has broken queries here before.
		Assignment("DEPLOYMENT_TYPE", OpenflowDeploymentTypeEnumDef.Kind(), g.ParameterOptions().SingleQuotes().Required()).
		OptionalAssignment("VPC_TYPE", OpenflowVpcTypeEnumDef.KindPtr(), g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("USE_PRIVATE_LINK", g.ParameterOptions()).
		OptionalBooleanAssignment("USE_USER_AUTH_OVER_PRIVATELINK", g.ParameterOptions()).
		OptionalTextAssignment("CUSTOM_INGRESS_HOSTNAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalQueryStructField("EventTable", openflowDeploymentEventTableDef(), g.ParameterOptions().SQL("EVENT_TABLE")).
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-deployment#alter-openflow-deployment",
	g.NewQueryStruct("AlterOpenflowDeployment").
		Alter().
		SQL("OPENFLOW DEPLOYMENT").
		IfExists().
		Name().
		OptionalSQL("UPGRADE").
		OptionalSQL("TERMINATE").
		RenameTo().
		OptionalQueryStructField(
			"Set",
			// Order follows the published ALTER OPENFLOW DEPLOYMENT ... SET syntax.
			g.NewQueryStruct("OpenflowDeploymentSet").
				OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				OptionalQueryStructField("EventTable", openflowDeploymentEventTableDef(), g.ParameterOptions().SQL("EVENT_TABLE")).
				WithValidation(g.AtLeastOneValueSet, "Comment", "DisplayName", "EventTable"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("OpenflowDeploymentUnset").
				OptionalSQL("DISPLAY_NAME").
				OptionalSQL("COMMENT").
				OptionalSQL("EVENT_TABLE").
				WithValidation(g.AtLeastOneValueSet, "Comment", "DisplayName", "EventTable"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "Upgrade", "Terminate", "RenameTo", "Set", "Unset"),
).DropOperation(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-deployment#drop-openflow-deployment",
	g.NewQueryStruct("DropOpenflowDeployment").
		Drop().
		SQL("OPENFLOW DEPLOYMENT").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-deployment#show-openflow-deployments",
	g.StructPair("openflowDeploymentRow", "OpenflowDeployment").
		Text("name").
		Enum("type", OpenflowDeploymentTypeEnumDef, g.WithDbFieldName("DeploymentType")).
		Enum("status", OpenflowDeploymentStatusEnumDef).
		OptionalEnum("vpc_type", OpenflowVpcTypeEnumDef).
		OptionalText("display_name").
		Bool("use_private_link").
		Bool("use_user_auth_over_private_link").
		OptionalText("custom_ingress_hostname").
		OptionalText("key").
		Text("owner").
		OptionalText("comment").
		Time("created_on").
		Time("updated_on"),
	g.NewQueryStruct("ShowOpenflowDeployments").
		Show().
		SQL("OPENFLOW DEPLOYMENTS").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimitFrom(),
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/LIMITEDACCESS/openflow-gen2/sql-reference/openflow-deployment#describe-openflow-deployment",
	g.StructPair("openflowDeploymentDetailsRow", "OpenflowDeploymentDetails").
		Text("name").
		Enum("type", OpenflowDeploymentTypeEnumDef, g.WithDbFieldName("DeploymentType")).
		Enum("status", OpenflowDeploymentStatusEnumDef).
		OptionalEnum("vpc_type", OpenflowVpcTypeEnumDef).
		OptionalText("display_name").
		Bool("use_private_link").
		Bool("use_user_auth_over_private_link").
		OptionalText("custom_ingress_hostname").
		OptionalText("key").
		Text("owner").
		OptionalText("comment"),
	g.NewQueryStruct("DescribeOpenflowDeployment").
		Describe().
		SQL("OPENFLOW DEPLOYMENT").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
	// EVENT_TABLE is settable on CREATE and via ALTER SET/UNSET, but it is not returned by SHOW OPENFLOW
	// DEPLOYMENTS or DESCRIBE OPENFLOW DEPLOYMENT. The only way to read it back is
	// SHOW PARAMETERS LIKE 'EVENT_TABLE' IN OPENFLOW DEPLOYMENT <name>, hence ShowParameters.
).ShowParameters(
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).WithEnums(
	OpenflowDeploymentTypeEnumDef,
	OpenflowVpcTypeEnumDef,
	OpenflowDeploymentStatusEnumDef,
)
