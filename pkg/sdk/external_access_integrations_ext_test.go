package sdk

func init() {
	id := externalAccessIntegrationsTestIdAccountObjectIdentifier
	networkRuleId := randomSchemaObjectIdentifier()
	networkRuleId2 := randomSchemaObjectIdentifier()
	apiAuthId := randomAccountObjectIdentifier()
	apiAuthId2 := randomAccountObjectIdentifier()
	secretId := randomSchemaObjectIdentifier()
	secretId2 := randomSchemaObjectIdentifier()
	tagId := randomSchemaObjectIdentifier()

	externalAccessIntegrationsTests.Create.
		withDefaultOpts(func() *CreateExternalAccessIntegrationOptions {
			return &CreateExternalAccessIntegrationOptions{
				name:                id,
				AllowedNetworkRules: []SchemaObjectIdentifier{networkRuleId},
				Enabled:             true,
			}
		}).
		withModify(
			case_ExternalAccessIntegrations_validation_Create_opts_AllowedApiAuthenticationIntegrations_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateExternalAccessIntegrationOptions) {
				opts.AllowedApiAuthenticationIntegrations = &ExternalAccessIntegrationAllowedApiAuthenticationIntegrations{
					None:         Bool(true),
					Integrations: []AccountObjectIdentifier{apiAuthId},
				}
			},
		).
		withExpectedSqlf(
			case_ExternalAccessIntegrations_sql_Create_basic,
			`CREATE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ENABLED = true`,
			id.FullyQualifiedName(), networkRuleId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalAccessIntegrations_sql_Create_all,
			func(opts *CreateExternalAccessIntegrationOptions) {
				opts.OrReplace = Bool(true)
				opts.AllowedNetworkRules = []SchemaObjectIdentifier{networkRuleId, networkRuleId2}
				opts.AllowedApiAuthenticationIntegrations = &ExternalAccessIntegrationAllowedApiAuthenticationIntegrations{
					Integrations: []AccountObjectIdentifier{apiAuthId, apiAuthId2},
				}
				opts.AllowedAuthenticationSecrets = &ExternalAccessIntegrationAllowedAuthenticationSecrets{
					Secrets: []SchemaObjectIdentifier{secretId, secretId2},
				}
				opts.Comment = String("test")
			},
			`CREATE OR REPLACE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s, %s) ALLOWED_API_AUTHENTICATION_INTEGRATIONS = (%s, %s) ALLOWED_AUTHENTICATION_SECRETS = (%s, %s) ENABLED = true COMMENT = 'test'`,
			id.FullyQualifiedName(), networkRuleId.FullyQualifiedName(), networkRuleId2.FullyQualifiedName(),
			apiAuthId.FullyQualifiedName(), apiAuthId2.FullyQualifiedName(),
			secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_allowedApiAuthenticationIntegrations_none",
			func(opts *CreateExternalAccessIntegrationOptions) {
				opts.AllowedApiAuthenticationIntegrations = &ExternalAccessIntegrationAllowedApiAuthenticationIntegrations{
					None: Bool(true),
				}
			},
			`CREATE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ALLOWED_API_AUTHENTICATION_INTEGRATIONS = none ENABLED = true`,
			id.FullyQualifiedName(), networkRuleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_allowedAuthenticationSecrets_all",
			func(opts *CreateExternalAccessIntegrationOptions) {
				opts.AllowedAuthenticationSecrets = &ExternalAccessIntegrationAllowedAuthenticationSecrets{
					All: Bool(true),
				}
			},
			`CREATE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ALLOWED_AUTHENTICATION_SECRETS = all ENABLED = true`,
			id.FullyQualifiedName(), networkRuleId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_allowedAuthenticationSecrets_none",
			func(opts *CreateExternalAccessIntegrationOptions) {
				opts.AllowedAuthenticationSecrets = &ExternalAccessIntegrationAllowedAuthenticationSecrets{
					None: Bool(true),
				}
			},
			`CREATE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ALLOWED_AUTHENTICATION_SECRETS = none ENABLED = true`,
			id.FullyQualifiedName(), networkRuleId.FullyQualifiedName(),
		)

	externalAccessIntegrationsTests.Alter.
		withModify(
			case_ExternalAccessIntegrations_validation_Alter_opts_Set_AllowedApiAuthenticationIntegrations_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterExternalAccessIntegrationOptions) {
				opts.Set = &ExternalAccessIntegrationSet{
					AllowedApiAuthenticationIntegrations: &ExternalAccessIntegrationAllowedApiAuthenticationIntegrations{
						None:         Bool(true),
						Integrations: []AccountObjectIdentifier{apiAuthId},
					},
				}
			},
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_AllowedNetworkRules_notEmptyWhenProvided",
			func(opts *AlterExternalAccessIntegrationOptions) {
				opts.Set = &ExternalAccessIntegrationSet{
					AllowedNetworkRules: []SchemaObjectIdentifier{},
				}
			},
			NewError("AllowedNetworkRules must not be empty when provided"),
		).
		withModifyAndExpectedSqlf(
			case_ExternalAccessIntegrations_sql_Alter_Set,
			func(opts *AlterExternalAccessIntegrationOptions) {
				opts.Set = &ExternalAccessIntegrationSet{
					AllowedNetworkRules: []SchemaObjectIdentifier{networkRuleId},
				}
			},
			`ALTER EXTERNAL ACCESS INTEGRATION %s SET ALLOWED_NETWORK_RULES = (%s)`,
			id.FullyQualifiedName(), networkRuleId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalAccessIntegrations_sql_Alter_Unset,
			func(opts *AlterExternalAccessIntegrationOptions) {
				opts.Unset = &ExternalAccessIntegrationUnset{
					Comment: Bool(true),
				}
			},
			`ALTER EXTERNAL ACCESS INTEGRATION %s UNSET COMMENT`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalAccessIntegrations_sql_Alter_SetTags,
			func(opts *AlterExternalAccessIntegrationOptions) {
				opts.SetTags = []TagAssociation{{Name: tagId, Value: "v"}}
			},
			`ALTER EXTERNAL ACCESS INTEGRATION %s SET TAG %s = 'v'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalAccessIntegrations_sql_Alter_UnsetTags,
			func(opts *AlterExternalAccessIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId}
			},
			`ALTER EXTERNAL ACCESS INTEGRATION %s UNSET TAG %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_all",
			func(opts *AlterExternalAccessIntegrationOptions) {
				opts.Set = &ExternalAccessIntegrationSet{
					AllowedNetworkRules: []SchemaObjectIdentifier{networkRuleId, networkRuleId2},
					AllowedApiAuthenticationIntegrations: &ExternalAccessIntegrationAllowedApiAuthenticationIntegrations{
						Integrations: []AccountObjectIdentifier{apiAuthId, apiAuthId2},
					},
					AllowedAuthenticationSecrets: &ExternalAccessIntegrationAllowedAuthenticationSecrets{
						Secrets: []SchemaObjectIdentifier{secretId, secretId2},
					},
					Enabled: Bool(true),
					Comment: String("test"),
				}
			},
			`ALTER EXTERNAL ACCESS INTEGRATION %s SET ALLOWED_NETWORK_RULES = (%s, %s) ALLOWED_API_AUTHENTICATION_INTEGRATIONS = (%s, %s) ALLOWED_AUTHENTICATION_SECRETS = (%s, %s) ENABLED = true COMMENT = 'test'`,
			id.FullyQualifiedName(), networkRuleId.FullyQualifiedName(), networkRuleId2.FullyQualifiedName(),
			apiAuthId.FullyQualifiedName(), apiAuthId2.FullyQualifiedName(),
			secretId.FullyQualifiedName(), secretId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_allowedApiAuthenticationIntegrations_none",
			func(opts *AlterExternalAccessIntegrationOptions) {
				opts.Set = &ExternalAccessIntegrationSet{
					AllowedApiAuthenticationIntegrations: &ExternalAccessIntegrationAllowedApiAuthenticationIntegrations{
						None: Bool(true),
					},
				}
			},
			`ALTER EXTERNAL ACCESS INTEGRATION %s SET ALLOWED_API_AUTHENTICATION_INTEGRATIONS = none`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_allowedAuthenticationSecrets_all",
			func(opts *AlterExternalAccessIntegrationOptions) {
				opts.Set = &ExternalAccessIntegrationSet{
					AllowedAuthenticationSecrets: &ExternalAccessIntegrationAllowedAuthenticationSecrets{
						All: Bool(true),
					},
				}
			},
			`ALTER EXTERNAL ACCESS INTEGRATION %s SET ALLOWED_AUTHENTICATION_SECRETS = all`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_allowedAuthenticationSecrets_none",
			func(opts *AlterExternalAccessIntegrationOptions) {
				opts.Set = &ExternalAccessIntegrationSet{
					AllowedAuthenticationSecrets: &ExternalAccessIntegrationAllowedAuthenticationSecrets{
						None: Bool(true),
					},
				}
			},
			`ALTER EXTERNAL ACCESS INTEGRATION %s SET ALLOWED_AUTHENTICATION_SECRETS = none`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_all",
			func(opts *AlterExternalAccessIntegrationOptions) {
				opts.Unset = &ExternalAccessIntegrationUnset{
					AllowedNetworkRules:                  Bool(true),
					AllowedApiAuthenticationIntegrations: Bool(true),
					AllowedAuthenticationSecrets:         Bool(true),
					Enabled:                              Bool(true),
					Comment:                              Bool(true),
				}
			},
			`ALTER EXTERNAL ACCESS INTEGRATION %s UNSET ALLOWED_NETWORK_RULES, ALLOWED_API_AUTHENTICATION_INTEGRATIONS, ALLOWED_AUTHENTICATION_SECRETS, ENABLED, COMMENT`,
			id.FullyQualifiedName(),
		)

	externalAccessIntegrationsTests.Drop.
		withExpectedSqlf(
			case_ExternalAccessIntegrations_sql_Drop_basic,
			`DROP EXTERNAL ACCESS INTEGRATION %s`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalAccessIntegrations_sql_Drop_all,
			func(opts *DropExternalAccessIntegrationOptions) {
				opts.IfExists = Bool(true)
			},
			`DROP EXTERNAL ACCESS INTEGRATION IF EXISTS %s`,
			id.FullyQualifiedName(),
		)

	externalAccessIntegrationsTests.Show.
		withExpectedSql(
			case_ExternalAccessIntegrations_sql_Show_basic,
			"SHOW EXTERNAL ACCESS INTEGRATIONS",
		).
		withModifyAndExpectedSqlf(
			case_ExternalAccessIntegrations_sql_Show_all,
			func(opts *ShowExternalAccessIntegrationOptions) {
				opts.Like = &Like{Pattern: String("test")}
			},
			"SHOW EXTERNAL ACCESS INTEGRATIONS LIKE 'test'",
		).
		withModifyAndExpectedSqlf(
			case_ExternalAccessIntegrations_sql_Show_Like,
			func(opts *ShowExternalAccessIntegrationOptions) {
				opts.Like = &Like{Pattern: String("test")}
			},
			"SHOW EXTERNAL ACCESS INTEGRATIONS LIKE 'test'",
		)

	externalAccessIntegrationsTests.Describe.
		withExpectedSqlf(
			case_ExternalAccessIntegrations_sql_Describe_basic,
			`DESCRIBE EXTERNAL ACCESS INTEGRATION %s`,
			id.FullyQualifiedName(),
		)
}
