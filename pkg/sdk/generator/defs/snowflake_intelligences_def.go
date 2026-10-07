package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// There is no SQL reference page for this object yet. The user guide documents the commands that are stable.
const snowflakeIntelligenceDoc = "https://docs.snowflake.com/en/user-guide/snowflake-cortex/snowflake-cowork/deploy-agents"

var snowflakeIntelligencesDef = g.NewInterface(
	"SnowflakeIntelligences",
	"SnowflakeIntelligence",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	CreateOperation(
		snowflakeIntelligenceDoc,
		g.NewQueryStruct("CreateSnowflakeIntelligence").
			Create().
			OrReplace().
			SQL("SNOWFLAKE INTELLIGENCE").
			IfNotExists().
			Name().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
	).
	AlterOperation(
		snowflakeIntelligenceDoc,
		g.NewQueryStruct("AlterSnowflakeIntelligence").
			Alter().
			SQL("SNOWFLAKE INTELLIGENCE").
			IfExists().
			Name().
			OptionalIdentifier("AddAgent", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("ADD AGENT")).
			OptionalIdentifier("DropAgent", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("DROP AGENT")).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifierIfSet, "AddAgent").
			WithValidation(g.ValidIdentifierIfSet, "DropAgent").
			WithValidation(g.ExactlyOneValueSet, "AddAgent", "DropAgent"),
	).
	DropOperation(
		snowflakeIntelligenceDoc,
		g.NewQueryStruct("DropSnowflakeIntelligence").
			Drop().
			SQL("SNOWFLAKE INTELLIGENCE").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		snowflakeIntelligenceDoc,
		g.StructPair("snowflakeIntelligenceDBRow", "SnowflakeIntelligence").
			Time("created_on").
			Text("name").
			Text("owner").
			OptionalText("comment", g.WithRequiredInPlain()),
		g.NewQueryStruct("ShowSnowflakeIntelligences").
			Show().
			SQL("SNOWFLAKE INTELLIGENCES"),
		g.ShowByIDNoFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		snowflakeIntelligenceDoc,
		g.StructPair("snowflakeIntelligenceDetailsRow", "SnowflakeIntelligenceDetails").
			Text("name").
			Text("owner").
			OptionalText("brand_name").
			OptionalText("vanity_url").
			OptionalText("welcome_message").
			OptionalText("icon_light_path").
			OptionalText("icon_dark_path").
			OptionalText("logo_light_path").
			OptionalText("logo_dark_path").
			OptionalText("favicon_path").
			OptionalText("accent_color_light").
			OptionalText("accent_color_dark").
			OptionalText("comment").
			Time("created_on"),
		g.NewQueryStruct("DescribeSnowflakeIntelligence").
			Describe().
			SQL("SNOWFLAKE INTELLIGENCE").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	CustomShowOperationWithPairedStructs(
		"ShowAgents",
		g.ShowMappingKindSlice,
		snowflakeIntelligenceDoc,
		g.StructPair("snowflakeIntelligenceAgentDBRow", "SnowflakeIntelligenceAgent").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			OptionalText("comment", g.WithRequiredInPlain()).
			OptionalPlainField("profile", "CortexAgentProfile", g.WithCustomParser("UnmarshalCortexAgentProfile")).
			OptionalBoolFromText("is_secure", g.WithBoolParsed(), g.WithRequiredInPlain()),
		g.NewQueryStruct("ShowAgentsInSnowflakeIntelligence").
			Show().
			SQL("AGENTS IN SNOWFLAKE INTELLIGENCE").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
