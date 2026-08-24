package sdk

func init() {
	id := networkPoliciesTestIdAccountObjectIdentifier
	allowedNetworkRuleId := randomSchemaObjectIdentifier()
	blockedNetworkRuleId := randomSchemaObjectIdentifier()
	renameTarget := randomAccountObjectIdentifier()

	networkPoliciesTests.Create.
		withExpectedSqlf(
			case_NetworkPolicies_sql_Create_basic,
			"CREATE NETWORK POLICY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NetworkPolicies_sql_Create_all,
			func(opts *CreateNetworkPolicyOptions) {
				opts.OrReplace = Bool(true)
				opts.AllowedNetworkRuleList = []SchemaObjectIdentifier{allowedNetworkRuleId}
				opts.BlockedNetworkRuleList = []SchemaObjectIdentifier{blockedNetworkRuleId}
				opts.AllowedIpList = []IP{{IP: "123.0.0.1"}, {IP: "321.0.0.1"}}
				opts.BlockedIpList = []IP{{IP: "123.0.0.1"}, {IP: "321.0.0.1"}}
				opts.Comment = String("some_comment")
			},
			"CREATE OR REPLACE NETWORK POLICY %s ALLOWED_NETWORK_RULE_LIST = (%s) BLOCKED_NETWORK_RULE_LIST = (%s) ALLOWED_IP_LIST = ('123.0.0.1', '321.0.0.1') BLOCKED_IP_LIST = ('123.0.0.1', '321.0.0.1') COMMENT = 'some_comment'",
			id.FullyQualifiedName(), allowedNetworkRuleId.FullyQualifiedName(), blockedNetworkRuleId.FullyQualifiedName(),
		)

	networkPoliciesTests.Alter.
		withDefaultOpts(func() *AlterNetworkPolicyOptions {
			return &AlterNetworkPolicyOptions{
				name:     id,
				IfExists: Bool(true),
			}
		}).
		withModify(
			case_NetworkPolicies_validation_Alter_opts_Add_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterNetworkPolicyOptions) {
				opts.Add = &AddNetworkRule{
					AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRuleId},
					BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRuleId},
				}
			},
		).
		withModify(
			case_NetworkPolicies_validation_Alter_opts_Remove_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterNetworkPolicyOptions) {
				opts.Remove = &RemoveNetworkRule{
					AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRuleId},
					BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRuleId},
				}
			},
		).
		withModifyAndExpectedSqlf(
			case_NetworkPolicies_sql_Alter_Set,
			func(opts *AlterNetworkPolicyOptions) {
				opts.Set = &NetworkPolicySet{
					AllowedIpList: &AllowedIPList{[]IP{{IP: "123.0.0.1"}}},
				}
			},
			"ALTER NETWORK POLICY IF EXISTS %s SET ALLOWED_IP_LIST = ('123.0.0.1')", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NetworkPolicies_sql_Alter_Unset,
			func(opts *AlterNetworkPolicyOptions) {
				opts.Unset = &NetworkPolicyUnset{
					AllowedNetworkRuleList: Bool(true),
					BlockedNetworkRuleList: Bool(true),
					AllowedIpList:          Bool(true),
					BlockedIpList:          Bool(true),
					Comment:                Bool(true),
				}
			},
			"ALTER NETWORK POLICY IF EXISTS %s UNSET ALLOWED_NETWORK_RULE_LIST, BLOCKED_NETWORK_RULE_LIST, ALLOWED_IP_LIST, BLOCKED_IP_LIST, COMMENT", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NetworkPolicies_sql_Alter_RenameTo,
			func(opts *AlterNetworkPolicyOptions) { opts.RenameTo = &renameTarget },
			"ALTER NETWORK POLICY IF EXISTS %s RENAME TO %s", id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NetworkPolicies_sql_Alter_Add,
			func(opts *AlterNetworkPolicyOptions) {
				opts.Add = &AddNetworkRule{
					AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRuleId},
				}
			},
			"ALTER NETWORK POLICY IF EXISTS %s ADD ALLOWED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), allowedNetworkRuleId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NetworkPolicies_sql_Alter_Remove,
			func(opts *AlterNetworkPolicyOptions) {
				opts.Remove = &RemoveNetworkRule{
					AllowedNetworkRuleList: []SchemaObjectIdentifier{allowedNetworkRuleId},
				}
			},
			"ALTER NETWORK POLICY IF EXISTS %s REMOVE ALLOWED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), allowedNetworkRuleId.FullyQualifiedName(),
		)

	networkPoliciesTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_Set_empty_ip_list",
			func(opts *AlterNetworkPolicyOptions) {
				opts.Set = &NetworkPolicySet{AllowedIpList: &AllowedIPList{}}
			},
			"ALTER NETWORK POLICY IF EXISTS %s SET ALLOWED_IP_LIST = ()", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_blocked_ip_list",
			func(opts *AlterNetworkPolicyOptions) {
				opts.Set = &NetworkPolicySet{BlockedIpList: &BlockedIPList{[]IP{{IP: "123.0.0.1"}}}}
			},
			"ALTER NETWORK POLICY IF EXISTS %s SET BLOCKED_IP_LIST = ('123.0.0.1')", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_empty_blocked_ip_list",
			func(opts *AlterNetworkPolicyOptions) {
				opts.Set = &NetworkPolicySet{BlockedIpList: &BlockedIPList{}}
			},
			"ALTER NETWORK POLICY IF EXISTS %s SET BLOCKED_IP_LIST = ()", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_allowed_network_rule_list",
			func(opts *AlterNetworkPolicyOptions) {
				opts.Set = &NetworkPolicySet{AllowedNetworkRuleList: &AllowedNetworkRuleList{[]SchemaObjectIdentifier{allowedNetworkRuleId}}}
			},
			"ALTER NETWORK POLICY IF EXISTS %s SET ALLOWED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), allowedNetworkRuleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_empty_allowed_network_rule_list",
			func(opts *AlterNetworkPolicyOptions) {
				opts.Set = &NetworkPolicySet{AllowedNetworkRuleList: &AllowedNetworkRuleList{}}
			},
			"ALTER NETWORK POLICY IF EXISTS %s SET ALLOWED_NETWORK_RULE_LIST = ()", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_blocked_network_rule_list",
			func(opts *AlterNetworkPolicyOptions) {
				opts.Set = &NetworkPolicySet{BlockedNetworkRuleList: &BlockedNetworkRuleList{[]SchemaObjectIdentifier{blockedNetworkRuleId}}}
			},
			"ALTER NETWORK POLICY IF EXISTS %s SET BLOCKED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), blockedNetworkRuleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_empty_blocked_network_rule_list",
			func(opts *AlterNetworkPolicyOptions) {
				opts.Set = &NetworkPolicySet{BlockedNetworkRuleList: &BlockedNetworkRuleList{}}
			},
			"ALTER NETWORK POLICY IF EXISTS %s SET BLOCKED_NETWORK_RULE_LIST = ()", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_comment",
			func(opts *AlterNetworkPolicyOptions) {
				opts.Set = &NetworkPolicySet{Comment: String("some_comment")}
			},
			"ALTER NETWORK POLICY IF EXISTS %s SET COMMENT = 'some_comment'", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_single",
			func(opts *AlterNetworkPolicyOptions) {
				opts.Unset = &NetworkPolicyUnset{Comment: Bool(true)}
			},
			"ALTER NETWORK POLICY IF EXISTS %s UNSET COMMENT", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Add_blocked_network_rule",
			func(opts *AlterNetworkPolicyOptions) {
				opts.Add = &AddNetworkRule{BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRuleId}}
			},
			"ALTER NETWORK POLICY IF EXISTS %s ADD BLOCKED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), blockedNetworkRuleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Remove_blocked_network_rule",
			func(opts *AlterNetworkPolicyOptions) {
				opts.Remove = &RemoveNetworkRule{BlockedNetworkRuleList: []SchemaObjectIdentifier{blockedNetworkRuleId}}
			},
			"ALTER NETWORK POLICY IF EXISTS %s REMOVE BLOCKED_NETWORK_RULE_LIST = (%s)", id.FullyQualifiedName(), blockedNetworkRuleId.FullyQualifiedName(),
		)

	networkPoliciesTests.Drop.
		withExpectedSqlf(
			case_NetworkPolicies_sql_Drop_basic,
			"DROP NETWORK POLICY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NetworkPolicies_sql_Drop_all,
			func(opts *DropNetworkPolicyOptions) { opts.IfExists = Bool(true) },
			"DROP NETWORK POLICY IF EXISTS %s", id.FullyQualifiedName(),
		)

	networkPoliciesTests.Show.
		withExpectedSql(case_NetworkPolicies_sql_Show_basic, "SHOW NETWORK POLICIES").
		withModifyAndExpectedSqlf(
			case_NetworkPolicies_sql_Show_all,
			func(opts *ShowNetworkPolicyOptions) { opts.Like = &Like{Pattern: String("some pattern")} },
			"SHOW NETWORK POLICIES LIKE 'some pattern'",
		).
		withModifyAndExpectedSqlf(
			case_NetworkPolicies_sql_Show_Like,
			func(opts *ShowNetworkPolicyOptions) { opts.Like = &Like{Pattern: String("some pattern")} },
			"SHOW NETWORK POLICIES LIKE 'some pattern'",
		)

	networkPoliciesTests.Describe.
		withExpectedSqlf(
			case_NetworkPolicies_sql_Describe_basic,
			"DESCRIBE NETWORK POLICY %s", id.FullyQualifiedName(),
		)
}
