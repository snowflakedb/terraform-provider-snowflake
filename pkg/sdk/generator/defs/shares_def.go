package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var shareKindEnumDef = g.NewEnum("ShareKind", "ShareKinds", "INBOUND", "OUTBOUND")

var sharePairs = g.StructPair("shareRow", "Share").
	Time("created_on").
	Enum("kind", shareKindEnumDef).
	Text("owner_account").
	Text("name").
	AccountObjectIdentifier("database_name", g.WithPlainFieldName("DatabaseName")).
	AccountIdentifierArray("to").
	Text("owner").
	Text("comment")

var shareInfoPairs = g.StructPair("shareDetailsRow", "ShareInfo").
	PlainField("kind", "ObjectType", g.WithManualConvert()).
	Field("name", "string", "ObjectIdentifier", g.WithManualConvert()).
	Time("shared_on")

func shareAdd() *g.QueryStruct {
	return g.NewQueryStruct("ShareAdd").
		ListAssignment("ACCOUNTS", "AccountIdentifier", g.ParameterOptions().Required()).
		OptionalBooleanAssignment("SHARE_RESTRICTIONS", g.ParameterOptions()).
		WithValidation(g.AtLeastOneValueSet, "Accounts")
}

func shareRemove() *g.QueryStruct {
	return g.NewQueryStruct("ShareRemove").
		ListAssignment("ACCOUNTS", "AccountIdentifier", g.ParameterOptions().Required()).
		WithValidation(g.AtLeastOneValueSet, "Accounts")
}

func shareSet() *g.QueryStruct {
	return g.NewQueryStruct("ShareSet").
		ListAssignment("ACCOUNTS", "AccountIdentifier", g.ParameterOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.AtLeastOneValueSet, "Accounts", "Comment")
}

func shareUnset() *g.QueryStruct {
	return g.NewQueryStruct("ShareUnset").
		OptionalSQL("COMMENT").
		WithValidation(g.ExactlyOneValueSet, "Comment")
}

var sharesDef = g.NewInterface(
	"Shares",
	"Share",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-share",
	g.NewQueryStruct("CreateShare").
		Create().
		OrReplace().
		SQL("SHARE").
		Name().
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-share",
	g.NewQueryStruct("AlterShare").
		Alter().
		SQL("SHARE").
		IfExists().
		Name().
		OptionalQueryStructField("Add", shareAdd(), g.KeywordOptions().SQL("ADD")).
		OptionalQueryStructField("Remove", shareRemove(), g.KeywordOptions().SQL("REMOVE")).
		OptionalQueryStructField("Set", shareSet(), g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("Unset", shareUnset(), g.KeywordOptions().SQL("UNSET")).
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Add", "Remove", "Set", "Unset", "SetTags", "UnsetTags"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-share",
	g.NewQueryStruct("DropShare").
		Drop().
		SQL("SHARE").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-shares",
	sharePairs,
	g.NewQueryStruct("ShowShares").
		Show().
		SQL("SHARES").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimit(),
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-share",
	shareInfoPairs,
	g.NewQueryStruct("DescribeShare").
		Describe().
		SQL("SHARE").
		Identifier("name", "ObjectIdentifier", g.IdentifierOptions().Required()).
		WithValidation(g.ValidIdentifier, "name"),
).WithCustomInterfaceMethod(
	"DescribeProvider",
	"DescribeProvider returns the share contents for a provider (outbound) share.",
	[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.AccountObjectIdentifier]())},
	"[]ShareInfo", "error",
).WithCustomInterfaceMethod(
	"DescribeConsumer",
	"DescribeConsumer returns the share contents for a consumer (inbound) share.",
	[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.ExternalObjectIdentifier]())},
	"[]ShareInfo", "error",
).WithEnums(shareKindEnumDef)
