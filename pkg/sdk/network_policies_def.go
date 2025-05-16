package sdk

import (
	"encoding/json"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/poc/generator"
)

//go:generate go run ./poc/main.go

// NetworkRulesSnowflakeDto is needed to unpack the applied network rules from the JSON response from Snowflake
type NetworkRulesSnowflakeDto struct {
	FullyQualifiedRuleName string
}

func ParseNetworkRulesSnowflakeDto(networkRulesStringValue string) ([]NetworkRulesSnowflakeDto, error) {
	var networkRules []NetworkRulesSnowflakeDto
	err := json.Unmarshal([]byte(networkRulesStringValue), &networkRules)
	if err != nil {
		return nil, err
	}
	return networkRules, nil
}

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

	NetworkPoliciesDef = g.NewInterface(
		"NetworkPolicies",
		"NetworkPolicy",
		g.KindOfT[AccountObjectIdentifier](),
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
				Identifier("RenameTo", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
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
		ShowOperation(
			"https://docs.snowflake.com/en/sql-reference/sql/show-network-policies",
			g.DbStruct("showNetworkPolicyDBRow").
				Field("created_on", "string").
				Field("name", "string").
				Field("comment", "string").
				Field("entries_in_allowed_ip_list", "int").
				Field("entries_in_blocked_ip_list", "int").
				Field("entries_in_allowed_network_rules", "int").
				Field("entries_in_blocked_network_rules", "int"),
			g.PlainStruct("NetworkPolicy").
				Field("CreatedOn", "string").
				Field("Name", "string").
				Field("Comment", "string").
				Field("EntriesInAllowedIpList", "int").
				Field("EntriesInBlockedIpList", "int").
				Field("EntriesInAllowedNetworkRules", "int").
				Field("EntriesInBlockedNetworkRules", "int"),
			g.NewQueryStruct("ShowNetworkPolicies").
				Show().
				SQL("NETWORK POLICIES").
				OptionalLike(),
		).
		ShowByIdOperationWithFiltering(
			g.ShowByIDLikeFiltering,
		).
		DescribeOperation(
			g.DescriptionMappingKindSlice,
			"https://docs.snowflake.com/en/sql-reference/sql/desc-network-policy",
			g.DbStruct("describeNetworkPolicyDBRow").
				Field("name", "string").
				Field("value", "string"),
			g.PlainStruct("NetworkPolicyProperty").
				Field("Name", "string").
				Field("Value", "string"),
			g.NewQueryStruct("DescribeNetworkPolicy").
				Describe().
				SQL("NETWORK POLICY").
				Name().
				WithValidation(g.ValidIdentifier, "name"),
		)
)
