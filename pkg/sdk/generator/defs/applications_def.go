package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

/*
 * 	todo: add definition for `CREATE APPLICATION <name> FROM LISTING <listing_name> [ COMMENT = '<string_literal>' ] [ WITH TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]`
 */

var versionAndPatch = g.NewQueryStruct("VersionAndPatch").
	TextAssignment("VERSION", g.ParameterOptions().NoEquals().NoQuotes().Required()).
	OptionalNumberAssignment("PATCH", g.ParameterOptions().NoEquals().Required())

var applicationVersion = func() *g.QueryStruct {
	return g.NewQueryStruct("ApplicationVersion").
		OptionalText("VersionDirectory", g.KeywordOptions().SingleQuotes()).
		OptionalQueryStructField("VersionAndPatch", versionAndPatch, g.KeywordOptions().NoQuotes()).
		WithValidation(g.ExactlyOneValueSet, "VersionDirectory", "VersionAndPatch")
}

var applicationSet = g.NewQueryStruct("ApplicationSet").
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("SHARE_EVENTS_WITH_PROVIDER", g.ParameterOptions()).
	OptionalBooleanAssignment("DEBUG_MODE", g.ParameterOptions())

var applicationUnset = g.NewQueryStruct("ApplicationUnset").
	OptionalSQL("COMMENT").
	OptionalSQL("SHARE_EVENTS_WITH_PROVIDER").
	OptionalSQL("DEBUG_MODE")

var applicationReferences = g.NewQueryStruct("ApplicationReferences").ListQueryStructField(
	"References",
	g.NewQueryStruct("ApplicationReference").Text("Reference", g.KeywordOptions().SingleQuotes()),
	g.ParameterOptions().Parentheses().NoEquals(),
)

var applicationsDef = g.NewInterface(
	"Applications",
	"Application",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-application",
	g.NewQueryStruct("CreateApplication").
		Create().
		SQL("APPLICATION").
		Name().
		SQL("FROM APPLICATION PACKAGE").
		Identifier("PackageName", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		OptionalQueryStructField(
			"Version",
			applicationVersion(),
			g.KeywordOptions().SQL("USING"),
		).
		OptionalBooleanAssignment("DEBUG_MODE", g.ParameterOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "PackageName"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-application",
	g.NewQueryStruct("DropApplication").
		Drop().
		SQL("APPLICATION").
		IfExists().
		Name().
		OptionalSQL("CASCADE").
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-application",
	g.NewQueryStruct("AlterApplication").
		Alter().
		SQL("APPLICATION").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Set",
			applicationSet,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			applicationUnset,
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		OptionalSQL("UPGRADE").
		OptionalQueryStructField(
			"UpgradeVersion",
			applicationVersion(),
			g.KeywordOptions().SQL("UPGRADE USING"),
		).
		OptionalQueryStructField(
			"UnsetReferences",
			applicationReferences,
			g.KeywordOptions().SQL("UNSET REFERENCES"),
		).
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "Upgrade", "UpgradeVersion", "UnsetReferences", "SetTags", "UnsetTags"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-applications",
	g.StructPair("applicationRow", "Application").
		Text("created_on").
		Text("name").
		Field("is_default", "string", "bool").
		Field("is_current", "string", "bool").
		Text("source_type").
		Text("source").
		Text("owner").
		Text("comment").
		Text("version").
		Text("label").
		Number("patch").
		Text("options").
		Number("retention_time"),
	g.NewQueryStruct("ShowApplications").
		Show().
		SQL("APPLICATIONS").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimit(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-application",
	g.StructPair("applicationPropertyRow", "ApplicationProperty").
		Text("property").
		OptionalText("value", g.WithRequiredInPlain()),
	g.NewQueryStruct("DescribeApplication").
		Describe().
		SQL("APPLICATION").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
