package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

/*
 * 	todo: add definition for `CREATE APPLICATION <name> FROM LISTING <listing_name> [ COMMENT = '<string_literal>' ] [ WITH TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]`
 */

var versionAndPatch = g.NewQueryStruct("VersionAndPatch").
	TextAssignment("VERSION", g.ParameterOptions().NoEquals().NoQuotes().Required()).
	OptionalNumberAssignment("PATCH", g.ParameterOptions().NoEquals().Required())

var applicationVersion = g.NewQueryStruct("ApplicationVersion").
	OptionalText("VersionDirectory", g.KeywordOptions().SingleQuotes()).
	OptionalQueryStructField("VersionAndPatch", versionAndPatch, g.KeywordOptions().NoQuotes()).
	WithValidation(g.ExactlyOneValueSet, "VersionDirectory", "VersionAndPatch")

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

var ApplicationsDef = g.NewInterface(
	"Applications",
	"Application",
	g.KindOfT[AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-application",
	g.NewQueryStruct("CreateApplication").
		Create().
		SQL("APPLICATION").
		Name().
		SQL("FROM APPLICATION PACKAGE").
		Identifier("PackageName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		OptionalQueryStructField(
			"Version",
			applicationVersion,
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
			applicationVersion,
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
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-applications",
	g.DbStruct("applicationRow").
		Field("created_on", "string").
		Field("name", "string").
		Field("is_default", "string").
		Field("is_current", "string").
		Field("source_type", "string").
		Field("source", "string").
		Field("owner", "string").
		Field("comment", "string").
		Field("version", "string").
		Field("label", "string").
		Field("patch", "int").
		Field("options", "string").
		Field("retention_time", "int"),
	g.PlainStruct("Application").
		Field("CreatedOn", "string").
		Field("Name", "string").
		Field("IsDefault", "bool").
		Field("IsCurrent", "bool").
		Field("SourceType", "string").
		Field("Source", "string").
		Field("Owner", "string").
		Field("Comment", "string").
		Field("Version", "string").
		Field("Label", "string").
		Field("Patch", "int").
		Field("Options", "string").
		Field("RetentionTime", "int"),
	g.NewQueryStruct("ShowApplications").
		Show().
		SQL("APPLICATIONS").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimit(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-application",
	g.DbStruct("applicationPropertyRow").
		Field("property", "string").
		Field("value", "sql.NullString"),
	g.PlainStruct("ApplicationProperty").
		Field("Property", "string").
		Field("Value", "string"),
	g.NewQueryStruct("DescribeApplication").
		Describe().
		SQL("APPLICATION").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
