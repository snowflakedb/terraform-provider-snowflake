//go:build sdk_generation

package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var openflowDeploymentsDef = g.NewInterface(
	"OpenflowDeployments",
	"OpenflowDeployment",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"TODO: add link when public docs are available",
	g.NewQueryStruct("CreateOpenflowDeployment").
		Create().
		SQL("OPENFLOW DEPLOYMENT").
		IfNotExists().
		Name().
		OptionalAssignment("DEPLOYMENT_TYPE", g.KindOfT[sdkcommons.OpenflowDeploymentType](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("VPC_TYPE", g.KindOfT[sdkcommons.OpenflowVpcType](), g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("CUSTOM_INGRESS_HOSTNAME", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("USE_PRIVATE_LINK", g.ParameterOptions()).
		OptionalBooleanAssignment("USE_USER_AUTH_OVER_PRIVATELINK", g.ParameterOptions()).
		OptionalTextAssignment("EVENT_TABLE", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"TODO: add link when public docs are available",
	g.NewQueryStruct("AlterOpenflowDeployment").
		Alter().
		SQL("OPENFLOW DEPLOYMENT").
		Name().
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("OpenflowDeploymentSet").
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("DISPLAY_NAME", g.ParameterOptions().SingleQuotes()).
				OptionalTextAssignment("EVENT_TABLE", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.AtLeastOneValueSet, "Comment", "DisplayName", "EventTable"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("OpenflowDeploymentUnset").
				OptionalSQL("COMMENT").
				OptionalSQL("DISPLAY_NAME").
				OptionalSQL("EVENT_TABLE").
				WithValidation(g.AtLeastOneValueSet, "Comment", "DisplayName", "EventTable"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset"),
).DropOperation(
	"TODO: add link when public docs are available",
	g.NewQueryStruct("DropOpenflowDeployment").
		Drop().
		SQL("OPENFLOW DEPLOYMENT").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"TODO: add link when public docs are available",
	g.DbStruct("openflowDeploymentRow").
		Text("name").
		Text("type").
		Text("status").
		OptionalText("vpc_type").
		OptionalText("display_name").
		Bool("use_private_link").
		Bool("use_user_auth_over_private_link").
		OptionalText("custom_ingress_hostname").
		Text("owner").
		OptionalText("comment").
		Time("created_on").
		Time("updated_on"),
	g.PlainStruct("OpenflowDeployment").
		Field("DeploymentType", "OpenflowDeploymentType").
		Field("Status", "OpenflowDeploymentStatus").
		Field("VpcType", "*OpenflowVpcType").
		OptionalText("DisplayName").
		Bool("UsePrivateLink").
		Bool("UseUserAuthOverPrivatelink").
		OptionalText("CustomIngressHostname").
		Text("Owner").
		OptionalText("Comment").
		Time("CreatedOn").
		Time("UpdatedOn"),
	g.NewQueryStruct("ShowOpenflowDeployments").
		Show().
		SQL("OPENFLOW DEPLOYMENTS").
		OptionalLike(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"TODO: add link when public docs are available",
	g.DbStruct("openflowDeploymentDetailsRow").
		Text("name").
		Text("type").
		Text("status").
		OptionalText("vpc_type").
		OptionalText("display_name").
		Bool("use_private_link").
		Bool("use_user_auth_over_private_link").
		OptionalText("custom_ingress_hostname").
		Text("owner").
		OptionalText("comment").
		Time("created_on").
		Time("updated_on").
		OptionalText("error_code").
		OptionalText("status_message"),
	g.PlainStruct("OpenflowDeploymentDetails").
		Field("DeploymentType", "OpenflowDeploymentType").
		Field("Status", "OpenflowDeploymentStatus").
		Field("VpcType", "*OpenflowVpcType").
		OptionalText("DisplayName").
		Bool("UsePrivateLink").
		Bool("UseUserAuthOverPrivatelink").
		OptionalText("CustomIngressHostname").
		Text("Owner").
		OptionalText("Comment").
		Time("CreatedOn").
		Time("UpdatedOn").
		OptionalText("ErrorCode").
		OptionalText("StatusMessage"),
	g.NewQueryStruct("DescribeOpenflowDeployment").
		Describe().
		SQL("OPENFLOW DEPLOYMENT").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
