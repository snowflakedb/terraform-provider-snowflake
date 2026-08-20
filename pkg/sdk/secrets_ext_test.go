package sdk

func init() {
	id := secretsTestIdSchemaObjectIdentifier
	apiIntegrationId := randomAccountObjectIdentifier()

	secretsTests.CreateWithOAuthClientCredentialsFlow.
		withDefaultOpts(func() *CreateWithOAuthClientCredentialsFlowSecretOptions {
			return &CreateWithOAuthClientCredentialsFlowSecretOptions{
				name:           id,
				ApiIntegration: apiIntegrationId,
			}
		}).
		withExpectedSqlf(
			case_Secrets_sql_CreateWithOAuthClientCredentialsFlow_basic,
			"CREATE SECRET %s TYPE = OAUTH2 API_AUTHENTICATION = %s",
			id.FullyQualifiedName(), apiIntegrationId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Secrets_sql_CreateWithOAuthClientCredentialsFlow_all,
			func(opts *CreateWithOAuthClientCredentialsFlowSecretOptions) {
				opts.IfNotExists = new(true)
				opts.OauthScopes = &OauthScopesList{[]ApiIntegrationScope{{Scope: "test"}}}
				opts.Comment = new("foo")
			},
			"CREATE SECRET IF NOT EXISTS %s TYPE = OAUTH2 API_AUTHENTICATION = %s OAUTH_SCOPES = ('test') COMMENT = 'foo'",
			id.FullyQualifiedName(), apiIntegrationId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateWithOAuthClientCredentialsFlow_emptyOauthScopesList",
			func(opts *CreateWithOAuthClientCredentialsFlowSecretOptions) {
				opts.IfNotExists = new(true)
				opts.OauthScopes = &OauthScopesList{[]ApiIntegrationScope{}}
				opts.Comment = new("foo")
			},
			"CREATE SECRET IF NOT EXISTS %s TYPE = OAUTH2 API_AUTHENTICATION = %s OAUTH_SCOPES = () COMMENT = 'foo'",
			id.FullyQualifiedName(), apiIntegrationId.FullyQualifiedName(),
		)

	secretsTests.CreateWithOAuthAuthorizationCodeFlow.
		withDefaultOpts(func() *CreateWithOAuthAuthorizationCodeFlowSecretOptions {
			return &CreateWithOAuthAuthorizationCodeFlowSecretOptions{
				name:                        id,
				OauthRefreshToken:           "foo",
				OauthRefreshTokenExpiryTime: "bar",
				ApiIntegration:              apiIntegrationId,
			}
		}).
		withExpectedSqlf(
			case_Secrets_sql_CreateWithOAuthAuthorizationCodeFlow_basic,
			"CREATE SECRET %s TYPE = OAUTH2 OAUTH_REFRESH_TOKEN = 'foo' OAUTH_REFRESH_TOKEN_EXPIRY_TIME = 'bar' API_AUTHENTICATION = %s",
			id.FullyQualifiedName(), apiIntegrationId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Secrets_sql_CreateWithOAuthAuthorizationCodeFlow_all,
			func(opts *CreateWithOAuthAuthorizationCodeFlowSecretOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("test")
			},
			"CREATE SECRET IF NOT EXISTS %s TYPE = OAUTH2 OAUTH_REFRESH_TOKEN = 'foo' OAUTH_REFRESH_TOKEN_EXPIRY_TIME = 'bar' API_AUTHENTICATION = %s COMMENT = 'test'",
			id.FullyQualifiedName(), apiIntegrationId.FullyQualifiedName(),
		)

	secretsTests.CreateWithBasicAuthentication.
		withDefaultOpts(func() *CreateWithBasicAuthenticationSecretOptions {
			return &CreateWithBasicAuthenticationSecretOptions{
				name:     id,
				Username: "foo",
				Password: "bar",
			}
		}).
		withExpectedSqlf(
			case_Secrets_sql_CreateWithBasicAuthentication_basic,
			"CREATE SECRET %s TYPE = PASSWORD USERNAME = 'foo' PASSWORD = 'bar'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Secrets_sql_CreateWithBasicAuthentication_all,
			func(opts *CreateWithBasicAuthenticationSecretOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("test")
			},
			"CREATE SECRET IF NOT EXISTS %s TYPE = PASSWORD USERNAME = 'foo' PASSWORD = 'bar' COMMENT = 'test'",
			id.FullyQualifiedName(),
		)

	secretsTests.CreateWithGenericString.
		withDefaultOpts(func() *CreateWithGenericStringSecretOptions {
			return &CreateWithGenericStringSecretOptions{
				name:         id,
				SecretString: "test",
			}
		}).
		withExpectedSqlf(
			case_Secrets_sql_CreateWithGenericString_basic,
			"CREATE SECRET %s TYPE = GENERIC_STRING SECRET_STRING = 'test'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Secrets_sql_CreateWithGenericString_all,
			func(opts *CreateWithGenericStringSecretOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("test")
			},
			"CREATE SECRET IF NOT EXISTS %s TYPE = GENERIC_STRING SECRET_STRING = 'test' COMMENT = 'test'",
			id.FullyQualifiedName(),
		)

	secretsTests.Alter.
		withDefaultOpts(func() *AlterSecretOptions {
			return &AlterSecretOptions{
				name:     id,
				IfExists: new(true),
			}
		}).
		withModifyAndExpectedSqlf(
			case_Secrets_sql_Alter_Set,
			func(opts *AlterSecretOptions) {
				opts.Set = &SecretSet{Comment: new("test")}
			},
			"ALTER SECRET IF EXISTS %s SET COMMENT = 'test'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Secrets_sql_Alter_Unset,
			func(opts *AlterSecretOptions) {
				opts.Unset = &SecretUnset{Comment: new(true)}
			},
			"ALTER SECRET IF EXISTS %s SET COMMENT = NULL",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_SetForOAuthClientCredentials",
			func(opts *AlterSecretOptions) {
				opts.Set = &SecretSet{
					Comment: new("test"),
					SetForFlow: &SetForFlow{SetForOAuthClientCredentials: &SetForOAuthClientCredentials{
						OauthScopes: &OauthScopesList{[]ApiIntegrationScope{{Scope: "sample_scope"}}},
					}},
				}
			},
			"ALTER SECRET IF EXISTS %s SET COMMENT = 'test' OAUTH_SCOPES = ('sample_scope')",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_SetForOAuthClientCredentials_emptyOauthScopesList",
			func(opts *AlterSecretOptions) {
				opts.Set = &SecretSet{
					Comment: new("test"),
					SetForFlow: &SetForFlow{SetForOAuthClientCredentials: &SetForOAuthClientCredentials{
						OauthScopes: &OauthScopesList{},
					}},
				}
			},
			"ALTER SECRET IF EXISTS %s SET COMMENT = 'test' OAUTH_SCOPES = ()",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_SetForOAuthAuthorization",
			func(opts *AlterSecretOptions) {
				opts.Set = &SecretSet{
					Comment: new("test"),
					SetForFlow: &SetForFlow{SetForOAuthAuthorization: &SetForOAuthAuthorization{
						OauthRefreshToken:           new("test_token"),
						OauthRefreshTokenExpiryTime: new("2024-11-11"),
					}},
				}
			},
			"ALTER SECRET IF EXISTS %s SET COMMENT = 'test' OAUTH_REFRESH_TOKEN = 'test_token' OAUTH_REFRESH_TOKEN_EXPIRY_TIME = '2024-11-11'",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_SetForBasicAuthentication",
			func(opts *AlterSecretOptions) {
				opts.Set = &SecretSet{
					Comment: new("test"),
					SetForFlow: &SetForFlow{SetForBasicAuthentication: &SetForBasicAuthentication{
						Username: new("foo"),
						Password: new("bar"),
					}},
				}
			},
			"ALTER SECRET IF EXISTS %s SET COMMENT = 'test' USERNAME = 'foo' PASSWORD = 'bar'",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_SetForGenericString",
			func(opts *AlterSecretOptions) {
				opts.Set = &SecretSet{
					Comment:    new("test"),
					SetForFlow: &SetForFlow{SetForGenericString: &SetForGenericString{SecretString: new("test")}},
				}
			},
			"ALTER SECRET IF EXISTS %s SET COMMENT = 'test' SECRET_STRING = 'test'",
			id.FullyQualifiedName(),
		)

	secretsTests.Drop.
		withExpectedSqlf(
			case_Secrets_sql_Drop_basic,
			"DROP SECRET %s",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Secrets_sql_Drop_all,
			func(opts *DropSecretOptions) { opts.IfExists = new(true) },
			"DROP SECRET IF EXISTS %s",
			id.FullyQualifiedName(),
		)

	secretsTests.Show.
		withExpectedSql(
			case_Secrets_sql_Show_basic,
			"SHOW SECRETS",
		).
		withModifyAndExpectedSqlf(
			case_Secrets_sql_Show_all,
			func(opts *ShowSecretOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
			},
			"SHOW SECRETS LIKE 'pattern' IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_Secrets_sql_Show_Like,
			func(opts *ShowSecretOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			"SHOW SECRETS LIKE 'pattern'",
		).
		withModifyAndExpectedSqlf(
			case_Secrets_sql_Show_In,
			func(opts *ShowSecretOptions) { opts.In = &ExtendedIn{In: In{Account: new(true)}} },
			"SHOW SECRETS IN ACCOUNT",
		)

	secretsTests.Describe.
		withExpectedSqlf(
			case_Secrets_sql_Describe_basic,
			"DESCRIBE SECRET %s",
			id.FullyQualifiedName(),
		)
}
