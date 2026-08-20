package sdk

func init() {
	id := networkRulesTestIdSchemaObjectIdentifier
	comment := "some comment"
	valueList := []NetworkRuleValue{
		{Value: "0.0.0.0"},
		{Value: "1.1.1.1"},
	}

	networkRulesTests.Create.
		withDefaultOpts(func() *CreateNetworkRuleOptions {
			return &CreateNetworkRuleOptions{
				name:            id,
				NetworkRuleType: NetworkRuleTypeIpv4,
				ValueList:       valueList,
				Mode:            NetworkRuleModeIngress,
			}
		}).
		withExpectedSqlf(
			case_NetworkRules_sql_Create_basic,
			`CREATE NETWORK RULE %s TYPE = IPV4 VALUE_LIST = ('0.0.0.0', '1.1.1.1') MODE = INGRESS`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NetworkRules_sql_Create_all,
			func(opts *CreateNetworkRuleOptions) {
				opts.OrReplace = new(true)
				opts.Comment = &comment
			},
			`CREATE OR REPLACE NETWORK RULE %s TYPE = IPV4 VALUE_LIST = ('0.0.0.0', '1.1.1.1') MODE = INGRESS COMMENT = '%s'`,
			id.FullyQualifiedName(), comment,
		).
		withAdditionalSqlCasef(
			"sql_Create_emptyValueList",
			func(opts *CreateNetworkRuleOptions) { opts.ValueList = []NetworkRuleValue{} },
			`CREATE NETWORK RULE %s TYPE = IPV4 MODE = INGRESS`,
			id.FullyQualifiedName(),
		)

	networkRulesTests.Alter.
		withModifyAndExpectedSqlf(
			case_NetworkRules_sql_Alter_basic,
			func(opts *AlterNetworkRuleOptions) {
				opts.Set = &NetworkRuleSet{ValueList: valueList}
			},
			`ALTER NETWORK RULE %s SET VALUE_LIST = ('0.0.0.0', '1.1.1.1')`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NetworkRules_sql_Alter_all,
			func(opts *AlterNetworkRuleOptions) {
				opts.IfExists = new(true)
				opts.Set = &NetworkRuleSet{ValueList: valueList, Comment: &comment}
			},
			`ALTER NETWORK RULE IF EXISTS %s SET VALUE_LIST = ('0.0.0.0', '1.1.1.1'), COMMENT = '%s'`,
			id.FullyQualifiedName(), comment,
		).
		withAdditionalSqlCasef(
			"sql_Alter_unset",
			func(opts *AlterNetworkRuleOptions) {
				opts.Unset = &NetworkRuleUnset{ValueList: new(true), Comment: new(true)}
			},
			`ALTER NETWORK RULE %s UNSET VALUE_LIST, COMMENT`,
			id.FullyQualifiedName(),
		)

	networkRulesTests.Drop.
		withExpectedSqlf(
			case_NetworkRules_sql_Drop_basic,
			`DROP NETWORK RULE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_NetworkRules_sql_Drop_all,
			func(opts *DropNetworkRuleOptions) { opts.IfExists = new(true) },
			`DROP NETWORK RULE IF EXISTS %s`, id.FullyQualifiedName(),
		)

	networkRulesTests.Show.
		withExpectedSql(case_NetworkRules_sql_Show_basic, `SHOW NETWORK RULES`).
		withModifyAndExpectedSqlf(
			case_NetworkRules_sql_Show_all,
			func(opts *ShowNetworkRuleOptions) {
				opts.Like = &Like{Pattern: new("name")}
				opts.In = &In{Database: NewAccountObjectIdentifier("database-name")}
				opts.StartsWith = new("abc")
				opts.Limit = &LimitFrom{Rows: new(10)}
			},
			`SHOW NETWORK RULES LIKE 'name' IN DATABASE "database-name" STARTS WITH 'abc' LIMIT 10`,
		).
		withModifyAndExpectedSqlf(
			case_NetworkRules_sql_Show_Like,
			func(opts *ShowNetworkRuleOptions) { opts.Like = &Like{Pattern: new("name")} },
			`SHOW NETWORK RULES LIKE 'name'`,
		).
		withModifyAndExpectedSqlf(
			case_NetworkRules_sql_Show_In,
			func(opts *ShowNetworkRuleOptions) {
				opts.In = &In{Database: NewAccountObjectIdentifier("database-name")}
			},
			`SHOW NETWORK RULES IN DATABASE "database-name"`,
		).
		withModifyAndExpectedSqlf(
			case_NetworkRules_sql_Show_StartsWith,
			func(opts *ShowNetworkRuleOptions) { opts.StartsWith = new("abc") },
			`SHOW NETWORK RULES STARTS WITH 'abc'`,
		).
		withModifyAndExpectedSqlf(
			case_NetworkRules_sql_Show_Limit,
			func(opts *ShowNetworkRuleOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			`SHOW NETWORK RULES LIMIT 10`,
		)

	networkRulesTests.Describe.
		withExpectedSqlf(
			case_NetworkRules_sql_Describe_basic,
			`DESCRIBE NETWORK RULE %s`, id.FullyQualifiedName(),
		)
}
