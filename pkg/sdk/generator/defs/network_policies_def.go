package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var (
	ip = g.NewQueryStruct("IP").
		Text("IP", g.KeywordOptions().SingleQuotes().Required())

	allowedNetworkRuleList = g.NewQueryStruct("AllowedNetworkRuleList").
				List("AllowedNetworkRuleList", "SchemaObjectIdentifier", g.ListOptions().MustParentheses())

	blockedNetworkRuleList = g.NewQueryStruct("BlockedNetworkRuleList").
				List("BlockedNetworkRuleList", "SchemaObjectIdentifier", g.ListOptions().MustParentheses())

	allowedIPList = g.NewQueryStruct("AllowedIPList").
			ListQueryStructField("AllowedIPList", ip, g.ListOptions().MustParentheses())

	blockedIPList = g.NewQueryStruct("BlockedIPList").
			ListQueryStructField("BlockedIPList", ip, g.ListOptions().MustParentheses())

	networkPoliciesAddNetworkRule = g.NewQueryStruct("AddNetworkRule").
					ListAssignment("ALLOWED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
					ListAssignment("BLOCKED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
					WithValidation(g.ExactlyOneValueSet, "AllowedNetworkRuleList", "BlockedNetworkRuleList")

	networkPoliciesRemoveNetworkRule = g.NewQueryStruct("RemoveNetworkRule").
						ListAssignment("ALLOWED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
						ListAssignment("BLOCKED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
						WithValidation(g.ExactlyOneValueSet, "AllowedNetworkRuleList", "BlockedNetworkRuleList")

	networkPoliciesDef = g.NewInterface(
		"NetworkPolicies",
		"NetworkPolicy",
		g.KindOfT[sdkcommons.AccountObjectIdentifier](),
	).
		CreateOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/create-network-policy",
			g.NewQueryStruct("CreateNetworkPolicies").
				Create().
				OrReplace().
				SQL("NETWORK POLICY").
				Name().
				ListAssignment("ALLOWED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
				ListAssignment("BLOCKED_NETWORK_RULE_LIST", "SchemaObjectIdentifier", g.ParameterOptions().Parentheses()).
				ListQueryStructField("AllowedIpList", ip, g.ParameterOptions().SQL("ALLOWED_IP_LIST").Parentheses()).
				ListQueryStructField("BlockedIpList", ip, g.ParameterOptions().SQL("BLOCKED_IP_LIST").Parentheses()).
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				WithValidation(g.ValidIdentifier, "name"),
		).
		AlterOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/alter-network-policy",
			g.NewQueryStruct("AlterNetworkPolicy").
				Alter().
				SQL("NETWORK POLICY").
				IfExists().
				Name().
				OptionalQueryStructField(
					"Set",
					g.NewQueryStruct("NetworkPolicySet").
						OptionalQueryStructField("AllowedNetworkRuleList", allowedNetworkRuleList, g.ParameterOptions().SQL("ALLOWED_NETWORK_RULE_LIST").Parentheses()).
						OptionalQueryStructField("BlockedNetworkRuleList", blockedNetworkRuleList, g.ParameterOptions().SQL("BLOCKED_NETWORK_RULE_LIST").Parentheses()).
						OptionalQueryStructField("AllowedIpList", allowedIPList, g.ParameterOptions().SQL("ALLOWED_IP_LIST").Parentheses()).
						OptionalQueryStructField("BlockedIpList", blockedIPList, g.ParameterOptions().SQL("BLOCKED_IP_LIST").Parentheses()).
						OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
						WithValidation(g.AtLeastOneValueSet, "AllowedIpList", "BlockedIpList", "Comment", "AllowedNetworkRuleList", "BlockedNetworkRuleList"),
					g.KeywordOptions().SQL("SET"),
				).
				OptionalQueryStructField(
					"Unset",
					g.NewQueryStruct("NetworkPolicyUnset").
						OptionalSQL("ALLOWED_NETWORK_RULE_LIST").
						OptionalSQL("BLOCKED_NETWORK_RULE_LIST").
						OptionalSQL("ALLOWED_IP_LIST").
						OptionalSQL("BLOCKED_IP_LIST").
						OptionalSQL("COMMENT").
						WithValidation(g.AtLeastOneValueSet, "AllowedIpList", "BlockedIpList", "Comment", "AllowedNetworkRuleList", "BlockedNetworkRuleList"),
					g.ListOptions().NoParentheses().SQL("UNSET"),
				).
				OptionalQueryStructField(
					"Add",
					networkPoliciesAddNetworkRule,
					g.KeywordOptions().SQL("ADD"),
				).
				OptionalQueryStructField(
					"Remove",
					networkPoliciesRemoveNetworkRule,
					g.KeywordOptions().SQL("REMOVE"),
				).
				Identifier("RenameTo", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
				WithValidation(g.ValidIdentifier, "name").
				WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "RenameTo", "Add", "Remove").
				WithValidation(g.ValidIdentifierIfSet, "RenameTo"),
		).
		DropOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/drop-network-policy",
			g.NewQueryStruct("DropNetworkPolicy").
				Drop().
				SQL("NETWORK POLICY").
				IfExists().
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		).
		ShowOperationWithPairedStructs(
			"https://docs.snowflake.com/en/sql-reference/sql/show-network-policies",
			g.StructPair("showNetworkPolicyDBRow", "NetworkPolicy").
				Text("created_on").
				Text("name").
				Text("comment").
				Number("entries_in_allowed_ip_list").
				Number("entries_in_blocked_ip_list").
				Number("entries_in_allowed_network_rules").
				Number("entries_in_blocked_network_rules").
				WithConvertGeneration(),
			g.NewQueryStruct("ShowNetworkPolicies").
				Show().
				SQL("NETWORK POLICIES").
				OptionalLike(),
		).
		DescribeOperationWithPairedStructs(
			g.DescriptionMappingKindSlice,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-network-policy",
			g.StructPair("describeNetworkPolicyDBRow", "NetworkPolicyProperty").
				Text("name").
				Text("value").
				WithConvertGeneration(),
			g.NewQueryStruct("DescribeNetworkPolicy").
				Describe().
				SQL("NETWORK POLICY").
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		)
)
