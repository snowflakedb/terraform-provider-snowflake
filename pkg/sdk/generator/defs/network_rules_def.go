package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var networkRulesDef = g.NewInterface(
	"NetworkRules",
	"NetworkRule",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-network-rule",
		g.NewQueryStruct("Create").
			Create().
			OrReplace().
			SQL("NETWORK RULE").
			Name().
			AssignmentWithFieldName("TYPE", g.KindOfT[sdkcommons.NetworkRuleType](), g.ParameterOptions().Required().NoQuotes(), "NetworkRuleType").
			ListAssignment("VALUE_LIST", "NetworkRuleValue", g.ParameterOptions().Required().Parentheses()).
			Assignment("MODE", g.KindOfT[sdkcommons.NetworkRuleMode](), g.ParameterOptions().Required().NoQuotes()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name"),
		g.NewQueryStruct("NetworkRuleValue").
			Text("Value", g.KeywordOptions().SingleQuotes().Required()),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-network-rule",
		g.NewQueryStruct("AlterNetworkRule").
			Alter().
			SQL("NETWORK RULE").
			IfExists().
			Name().
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("NetworkRuleSet").
					ListAssignment("VALUE_LIST", "NetworkRuleValue", g.ParameterOptions().Parentheses()).
					OptionalComment().
					WithValidation(g.AtLeastOneValueSet, "ValueList", "Comment"),
				g.ListOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("NetworkRuleUnset").
					OptionalSQL("VALUE_LIST").
					OptionalSQL("COMMENT").
					WithValidation(g.AtLeastOneValueSet, "ValueList", "Comment"),
				g.ListOptions().SQL("UNSET"),
			).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.AtLeastOneValueSet, "Set", "Unset"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-network-rule",
		g.NewQueryStruct("DropNetworkRule").
			Drop().
			SQL("NETWORK RULE").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-network-rules",
		g.StructPair("ShowNetworkRulesRow", "NetworkRule").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			Text("comment").
			PlainField("type", "NetworkRuleType").
			PlainField("mode", "NetworkRuleMode").
			Number("entries_in_valuelist", g.WithPlainFieldName("EntriesInValueList")).
			Text("owner_role_type"),
		g.NewQueryStruct("ShowNetworkRules").
			Show().
			SQL("NETWORK RULES").
			OptionalLike().
			OptionalIn().
			OptionalStartsWith().
			OptionalLimitFrom(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-network-rule",
		g.StructPair("DescNetworkRulesRow", "NetworkRuleDetails").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			Text("comment").
			PlainField("type", "NetworkRuleType").
			PlainField("mode", "NetworkRuleMode").
			StringList("value_list"),
		g.NewQueryStruct("ShowNetworkRules").
			Describe().
			SQL("NETWORK RULE").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
