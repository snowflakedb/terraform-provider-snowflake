package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

type NetworkRuleType string

const (
	NetworkRuleTypeIpv4             NetworkRuleType = "IPV4"
	NetworkRuleTypeAwsVpcEndpointId NetworkRuleType = "AWSVPCEID"
	NetworkRuleTypeAzureLinkId      NetworkRuleType = "AZURELINKID"
	NetworkRuleTypeHostPort         NetworkRuleType = "HOST_PORT"
)

type NetworkRuleMode string

const (
	NetworkRuleModeIngress       NetworkRuleMode = "INGRESS"
	NetworkRuleModeInternalStage NetworkRuleMode = "INTERNAL_STAGE"
	NetworkRuleModeEgress        NetworkRuleMode = "EGRESS"
)

var NetworkRuleDef = g.NewInterface(
	"NetworkRules",
	"NetworkRule",
	g.KindOfT[SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-network-rule",
		g.NewQueryStruct("Create").
			Create().
			OrReplace().
			SQL("NETWORK RULE").
			Name().
			AssignmentWithFieldName("TYPE", g.KindOfT[NetworkRuleType](), g.ParameterOptions().Required().NoQuotes(), "NetworkRuleType").
			ListAssignment("VALUE_LIST", "NetworkRuleValue", g.ParameterOptions().Required().Parentheses()).
			Assignment("MODE", g.KindOfT[NetworkRuleMode](), g.ParameterOptions().Required().NoQuotes()).
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
					ListAssignment("VALUE_LIST", "NetworkRuleValue", g.ParameterOptions().Required().Parentheses()).
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
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-network-rules",
		g.DbStruct("ShowNetworkRulesRow").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			Text("comment").
			Text("type").
			Text("mode").
			Number("entries_in_valuelist").
			Text("owner_role_type"),
		g.PlainStruct("NetworkRule").
			Time("CreatedOn").
			Text("Name").
			Text("DatabaseName").
			Text("SchemaName").
			Text("Owner").
			Text("Comment").
			Field("Type", "NetworkRuleType").
			Field("Mode", "NetworkRuleMode").
			Number("EntriesInValueList").
			Text("OwnerRoleType"),
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
	DescribeOperation(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-network-rule",
		g.DbStruct("DescNetworkRulesRow").
			Time("created_on").
			Text("name").
			Text("database_name").
			Text("schema_name").
			Text("owner").
			Text("comment").
			Text("type").
			Text("mode").
			Text("value_list"),
		g.PlainStruct("NetworkRuleDetails").
			Time("CreatedOn").
			Text("Name").
			Text("DatabaseName").
			Text("SchemaName").
			Text("Owner").
			Text("Comment").
			Field("Type", "NetworkRuleType").
			Field("Mode", "NetworkRuleMode").
			Field("ValueList", "[]string"),
		g.NewQueryStruct("ShowNetworkRules").
			Describe().
			SQL("NETWORK RULE").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	)
