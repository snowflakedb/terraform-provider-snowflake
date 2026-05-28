package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var cortexAgentsDef = g.NewInterface(
	"CortexAgents",
	"CortexAgent",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-agent",
		g.NewQueryStruct("CreateCortexAgent").
			Create().
			OrReplace().
			SQL("AGENT").
			IfNotExists().
			Name().
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("PROFILE", g.ParameterOptions().SingleQuotes()).
			TextAssignment("FROM SPECIFICATION", g.ParameterOptions().NoEquals().DoubleDollarQuotes()).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-agent",
		g.NewQueryStruct("AlterCortexAgent").
			Alter().
			SQL("AGENT").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("CortexAgentSet").
					OptionalAssignment("COMMENT", "StringAllowEmpty", g.ParameterOptions()).
					OptionalTextAssignment("PROFILE", g.ParameterOptions().SingleQuotes()).
					WithValidation(g.AtLeastOneValueSet, "Comment", "Profile"),
				g.ListOptions().NoParentheses().SQL("SET"),
			).
			OptionalQueryStructField(
				"ModifyLiveVersionSet",
				g.NewQueryStruct("CortexAgentModifyLiveVersionSet").
					TextAssignment("SPECIFICATION", g.ParameterOptions().DoubleDollarQuotes()),
				g.KeywordOptions().SQL("MODIFY LIVE VERSION SET"),
			).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Set", "ModifyLiveVersionSet"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-agent",
		g.NewQueryStruct("DropCortexAgent").
			Drop().
			SQL("AGENT").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-agents",
		g.DbStruct("showCortexAgentDBRow").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			OptionalText("comment").
			OptionalText("profile"),
		g.PlainStruct("CortexAgent").
			Time("CreatedOn").
			Text("Name").
			Text("DatabaseName").
			Text("SchemaName").
			Text("Owner").
			Text("Comment").
			Field("Profile", "CortexAgentProfile"),
		g.NewQueryStruct("ShowCortexAgents").
			Show().
			SQL("AGENTS").
			OptionalLike().
			OptionalExtendedIn().
			OptionalStartsWith().
			OptionalLimit(),
		g.ShowByIDLikeFiltering,
		g.ShowByIDExtendedInFiltering,
	).
	DescribeOperation(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-agent",
		g.DbStruct("cortexAgentDetailsRow").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			OptionalText("comment").
			OptionalText("profile").
			Text("agent_spec").
			Time("created_on").
			OptionalText("default_version_name").
			OptionalText("versions").
			OptionalText("aliases"),
		g.PlainStruct("CortexAgentDetails").
			Text("Name").
			Text("DatabaseName").
			Text("SchemaName").
			Text("Owner").
			Text("Comment").
			Field("Profile", "CortexAgentProfile").
			Text("AgentSpec").
			Time("CreatedOn").
			OptionalText("DefaultVersionName").
			OptionalText("Versions").
			OptionalText("Aliases"),
		g.NewQueryStruct("DescribeCortexAgent").
			Describe().
			SQL("AGENT").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	WithShowObjectType("Agent")
