package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	id := catalogIntegrationsTestIdAccountObjectIdentifier
	tagId := NewAccountObjectIdentifier("tag1")
	glueAwsRoleArn := "arn:aws:iam::123456789012:role/sqsAccess"
	glueCatalogId := "123456789012"
	glueRegion := "us-east-2"
	polarisCatalogUri := "https://testorg-testacc.snowflakecomputing.com/polaris/api/catalog"
	restCatalogUri := "https://api.tabular.io/ws"
	oAuthClientId := "my_client_id"
	oAuthClientSecret := "my_client_secret"
	oAuthAllowedScope := "PRINCIPAL_ROLE:ALL"
	oAuthTokenUri := "https://api.tabular.io/ws/v1/oauth/tokens" //nolint:gosec // dummy value, not a production address
	sigV4IamRole := "arn:aws:iam::123456789012:role/my-role"
	sigV4SigningRegion := "us-west-2"
	sigV4ExternalId := "external_id"

	catalogIntegrationsTests.Create.
		withDefaultOpts(func() *CreateCatalogIntegrationOptions {
			return &CreateCatalogIntegrationOptions{
				name: id,
				AwsGlueCatalogSourceParams: &AwsGlueParams{
					GlueAwsRoleArn: glueAwsRoleArn,
					GlueCatalogId:  glueCatalogId,
				},
				Enabled: true,
			}
		}).
		withModify(
			case_CatalogIntegrations_validation_Create_opts_IcebergRestCatalogSourceParams_ExactlyOneValueSet_NoneSet,
			func(opts *CreateCatalogIntegrationOptions) {
				opts.AwsGlueCatalogSourceParams = nil
				opts.IcebergRestCatalogSourceParams = &IcebergRestParams{}
			},
		).
		withModify(
			case_CatalogIntegrations_validation_Create_opts_IcebergRestCatalogSourceParams_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateCatalogIntegrationOptions) {
				opts.AwsGlueCatalogSourceParams = nil
				opts.IcebergRestCatalogSourceParams = &IcebergRestParams{
					OAuthRestAuthentication:  &OAuthRestAuthentication{},
					BearerRestAuthentication: &BearerRestAuthentication{},
				}
			},
		).
		withExpectedSqlf(
			case_CatalogIntegrations_sql_Create_basic,
			"CREATE CATALOG INTEGRATION %s CATALOG_SOURCE = GLUE TABLE_FORMAT = ICEBERG GLUE_AWS_ROLE_ARN = '%s' GLUE_CATALOG_ID = '%s' ENABLED = true",
			id.FullyQualifiedName(), glueAwsRoleArn, glueCatalogId,
		).
		withModifyAndExpectedSqlf(
			case_CatalogIntegrations_sql_Create_all,
			func(opts *CreateCatalogIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.AwsGlueCatalogSourceParams.GlueRegion = new(glueRegion)
				opts.AwsGlueCatalogSourceParams.CatalogNamespace = new("myNamespace")
				opts.Enabled = false
				opts.RefreshIntervalSeconds = new(60)
				opts.Comment = new("test comment")
			},
			"CREATE CATALOG INTEGRATION IF NOT EXISTS %s "+
				"CATALOG_SOURCE = GLUE "+
				"TABLE_FORMAT = ICEBERG "+
				"GLUE_AWS_ROLE_ARN = '%s' "+
				"GLUE_CATALOG_ID = '%s' "+
				"GLUE_REGION = '%s' "+
				"CATALOG_NAMESPACE = 'myNamespace' "+
				"ENABLED = false "+
				"REFRESH_INTERVAL_SECONDS = 60 "+
				"COMMENT = 'test comment'",
			id.FullyQualifiedName(), glueAwsRoleArn, glueCatalogId, glueRegion,
		).
		withAdditionalSqlCasef(
			"sql_Create_objectStorage",
			func(opts *CreateCatalogIntegrationOptions) {
				opts.AwsGlueCatalogSourceParams = nil
				opts.ObjectStorageCatalogSourceParams = &ObjectStorageParams{
					TableFormat: CatalogIntegrationTableFormatDelta,
				}
			},
			"CREATE CATALOG INTEGRATION %s CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = DELTA ENABLED = true",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_openCatalog",
			func(opts *CreateCatalogIntegrationOptions) {
				opts.AwsGlueCatalogSourceParams = nil
				opts.OpenCatalogCatalogSourceParams = &OpenCatalogParams{
					RestConfig: OpenCatalogRestConfig{
						CatalogUri:  polarisCatalogUri,
						CatalogName: "my_catalog_name",
					},
					RestAuthentication: OAuthRestAuthentication{
						OauthClientId:      oAuthClientId,
						OauthClientSecret:  oAuthClientSecret,
						OauthAllowedScopes: []StringListItemWrapper{{Value: oAuthAllowedScope}},
					},
				}
			},
			"CREATE CATALOG INTEGRATION %s "+
				"CATALOG_SOURCE = POLARIS "+
				"TABLE_FORMAT = ICEBERG "+
				"REST_CONFIG = (CATALOG_URI = '%s' CATALOG_NAME = 'my_catalog_name') "+
				"REST_AUTHENTICATION = (TYPE = OAUTH OAUTH_CLIENT_ID = '%s' OAUTH_CLIENT_SECRET = '%s' OAUTH_ALLOWED_SCOPES = ('%s')) "+
				"ENABLED = true",
			id.FullyQualifiedName(), polarisCatalogUri, oAuthClientId, oAuthClientSecret, oAuthAllowedScope,
		).
		withAdditionalSqlCasef(
			"sql_Create_icebergRest",
			func(opts *CreateCatalogIntegrationOptions) {
				opts.AwsGlueCatalogSourceParams = nil
				opts.IcebergRestCatalogSourceParams = &IcebergRestParams{
					RestConfig: IcebergRestRestConfig{
						CatalogUri: restCatalogUri,
					},
					SigV4RestAuthentication: &SigV4RestAuthentication{
						Sigv4IamRole: sigV4IamRole,
					},
				}
			},
			"CREATE CATALOG INTEGRATION %s "+
				"CATALOG_SOURCE = ICEBERG_REST "+
				"TABLE_FORMAT = ICEBERG "+
				"REST_CONFIG = (CATALOG_URI = '%s') "+
				"REST_AUTHENTICATION = (TYPE = SIGV4 SIGV4_IAM_ROLE = '%s') "+
				"ENABLED = true",
			id.FullyQualifiedName(), restCatalogUri, sigV4IamRole,
		).
		withAdditionalSqlCasef(
			"sql_Create_objectStorage_all",
			func(opts *CreateCatalogIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.AwsGlueCatalogSourceParams = nil
				opts.ObjectStorageCatalogSourceParams = &ObjectStorageParams{
					TableFormat: CatalogIntegrationTableFormatDelta,
				}
				opts.Enabled = false
				opts.RefreshIntervalSeconds = new(60)
				opts.Comment = new("test comment")
			},
			"CREATE CATALOG INTEGRATION IF NOT EXISTS %s "+
				"CATALOG_SOURCE = OBJECT_STORE "+
				"TABLE_FORMAT = DELTA "+
				"ENABLED = false "+
				"REFRESH_INTERVAL_SECONDS = 60 "+
				"COMMENT = 'test comment'",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_openCatalog_all",
			func(opts *CreateCatalogIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.AwsGlueCatalogSourceParams = nil
				opts.OpenCatalogCatalogSourceParams = &OpenCatalogParams{
					CatalogNamespace: new("myNamespace"),
					RestConfig: OpenCatalogRestConfig{
						CatalogUri:           polarisCatalogUri,
						CatalogApiType:       new(CatalogIntegrationCatalogApiTypePublic),
						CatalogName:          "my_catalog_name",
						AccessDelegationMode: new(CatalogIntegrationAccessDelegationModeVendedCredentials),
					},
					RestAuthentication: OAuthRestAuthentication{
						OauthTokenUri:      new(oAuthTokenUri),
						OauthClientId:      oAuthClientId,
						OauthClientSecret:  oAuthClientSecret,
						OauthAllowedScopes: []StringListItemWrapper{{Value: oAuthAllowedScope}},
					},
				}
				opts.Enabled = false
				opts.RefreshIntervalSeconds = new(60)
				opts.Comment = new("test comment")
			},
			"CREATE CATALOG INTEGRATION IF NOT EXISTS %s "+
				"CATALOG_SOURCE = POLARIS "+
				"TABLE_FORMAT = ICEBERG "+
				"CATALOG_NAMESPACE = 'myNamespace' "+
				"REST_CONFIG = (CATALOG_URI = '%s' CATALOG_API_TYPE = %s CATALOG_NAME = 'my_catalog_name' ACCESS_DELEGATION_MODE = %s) "+
				"REST_AUTHENTICATION = (TYPE = OAUTH OAUTH_TOKEN_URI = '%s' OAUTH_CLIENT_ID = '%s' OAUTH_CLIENT_SECRET = '%s' OAUTH_ALLOWED_SCOPES = ('%s')) "+
				"ENABLED = false "+
				"REFRESH_INTERVAL_SECONDS = 60 "+
				"COMMENT = 'test comment'",
			id.FullyQualifiedName(), polarisCatalogUri, CatalogIntegrationCatalogApiTypePublic, CatalogIntegrationAccessDelegationModeVendedCredentials, oAuthTokenUri, oAuthClientId, oAuthClientSecret, oAuthAllowedScope,
		).
		withAdditionalSqlCasef(
			"sql_Create_icebergRest_sigV4_all",
			func(opts *CreateCatalogIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.AwsGlueCatalogSourceParams = nil
				opts.IcebergRestCatalogSourceParams = &IcebergRestParams{
					CatalogNamespace: new("myNamespace"),
					RestConfig: IcebergRestRestConfig{
						CatalogUri:           restCatalogUri,
						Prefix:               new("prefix"),
						CatalogName:          new("my_catalog_name"),
						CatalogApiType:       new(CatalogIntegrationCatalogApiTypeAwsApiGateway),
						AccessDelegationMode: new(CatalogIntegrationAccessDelegationModeVendedCredentials),
					},
					SigV4RestAuthentication: &SigV4RestAuthentication{
						Sigv4IamRole:       sigV4IamRole,
						Sigv4ExternalId:    new(sigV4ExternalId),
						Sigv4SigningRegion: new(sigV4SigningRegion),
					},
				}
				opts.Enabled = false
				opts.RefreshIntervalSeconds = new(60)
				opts.Comment = new("test comment")
			},
			"CREATE CATALOG INTEGRATION IF NOT EXISTS %s "+
				"CATALOG_SOURCE = ICEBERG_REST "+
				"TABLE_FORMAT = ICEBERG "+
				"CATALOG_NAMESPACE = 'myNamespace' "+
				"REST_CONFIG = (CATALOG_URI = '%s' PREFIX = 'prefix' CATALOG_NAME = 'my_catalog_name' CATALOG_API_TYPE = %s ACCESS_DELEGATION_MODE = %s) "+
				"REST_AUTHENTICATION = (TYPE = SIGV4 SIGV4_IAM_ROLE = '%s' SIGV4_SIGNING_REGION = '%s' SIGV4_EXTERNAL_ID = '%s') "+
				"ENABLED = false "+
				"REFRESH_INTERVAL_SECONDS = 60 "+
				"COMMENT = 'test comment'",
			id.FullyQualifiedName(), restCatalogUri, CatalogIntegrationCatalogApiTypeAwsApiGateway, CatalogIntegrationAccessDelegationModeVendedCredentials, sigV4IamRole, sigV4SigningRegion, sigV4ExternalId,
		).
		withAdditionalSqlCasef(
			"sql_Create_icebergRest_oauth_all",
			func(opts *CreateCatalogIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.AwsGlueCatalogSourceParams = nil
				opts.IcebergRestCatalogSourceParams = &IcebergRestParams{
					CatalogNamespace: new("myNamespace"),
					RestConfig: IcebergRestRestConfig{
						CatalogUri:           restCatalogUri,
						Prefix:               new("prefix"),
						CatalogName:          new("my_catalog_name"),
						CatalogApiType:       new(CatalogIntegrationCatalogApiTypeAwsApiGateway),
						AccessDelegationMode: new(CatalogIntegrationAccessDelegationModeVendedCredentials),
					},
					OAuthRestAuthentication: &OAuthRestAuthentication{
						OauthClientId:      oAuthClientId,
						OauthClientSecret:  oAuthClientSecret,
						OauthAllowedScopes: []StringListItemWrapper{{Value: oAuthAllowedScope}},
					},
				}
				opts.Enabled = false
				opts.RefreshIntervalSeconds = new(60)
				opts.Comment = new("test comment")
			},
			"CREATE CATALOG INTEGRATION IF NOT EXISTS %s "+
				"CATALOG_SOURCE = ICEBERG_REST "+
				"TABLE_FORMAT = ICEBERG "+
				"CATALOG_NAMESPACE = 'myNamespace' "+
				"REST_CONFIG = (CATALOG_URI = '%s' PREFIX = 'prefix' CATALOG_NAME = 'my_catalog_name' CATALOG_API_TYPE = %s ACCESS_DELEGATION_MODE = %s) "+
				"REST_AUTHENTICATION = (TYPE = OAUTH OAUTH_CLIENT_ID = '%s' OAUTH_CLIENT_SECRET = '%s' OAUTH_ALLOWED_SCOPES = ('%s')) "+
				"ENABLED = false "+
				"REFRESH_INTERVAL_SECONDS = 60 "+
				"COMMENT = 'test comment'",
			id.FullyQualifiedName(), restCatalogUri, CatalogIntegrationCatalogApiTypeAwsApiGateway, CatalogIntegrationAccessDelegationModeVendedCredentials, oAuthClientId, oAuthClientSecret, oAuthAllowedScope,
		).
		withAdditionalSqlCasef(
			"sql_Create_icebergRest_bearer_all",
			func(opts *CreateCatalogIntegrationOptions) {
				opts.IfNotExists = new(true)
				opts.AwsGlueCatalogSourceParams = nil
				opts.IcebergRestCatalogSourceParams = &IcebergRestParams{
					CatalogNamespace: new("myNamespace"),
					RestConfig: IcebergRestRestConfig{
						CatalogUri:           restCatalogUri,
						Prefix:               new("prefix"),
						CatalogName:          new("my_catalog_name"),
						CatalogApiType:       new(CatalogIntegrationCatalogApiTypeAwsApiGateway),
						AccessDelegationMode: new(CatalogIntegrationAccessDelegationModeVendedCredentials),
					},
					BearerRestAuthentication: &BearerRestAuthentication{
						BearerToken: "test-token",
					},
				}
				opts.Enabled = false
				opts.RefreshIntervalSeconds = new(60)
				opts.Comment = new("test comment")
			},
			"CREATE CATALOG INTEGRATION IF NOT EXISTS %s "+
				"CATALOG_SOURCE = ICEBERG_REST "+
				"TABLE_FORMAT = ICEBERG "+
				"CATALOG_NAMESPACE = 'myNamespace' "+
				"REST_CONFIG = (CATALOG_URI = '%s' PREFIX = 'prefix' CATALOG_NAME = 'my_catalog_name' CATALOG_API_TYPE = %s ACCESS_DELEGATION_MODE = %s) "+
				"REST_AUTHENTICATION = (TYPE = BEARER BEARER_TOKEN = 'test-token') "+
				"ENABLED = false "+
				"REFRESH_INTERVAL_SECONDS = 60 "+
				"COMMENT = 'test comment'",
			id.FullyQualifiedName(), restCatalogUri, CatalogIntegrationCatalogApiTypeAwsApiGateway, CatalogIntegrationAccessDelegationModeVendedCredentials,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateCatalogIntegrationOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE CATALOG INTEGRATION %s CATALOG_SOURCE = GLUE TABLE_FORMAT = ICEBERG GLUE_AWS_ROLE_ARN = '%s' GLUE_CATALOG_ID = '%s' ENABLED = true",
			id.FullyQualifiedName(), glueAwsRoleArn, glueCatalogId,
		)

	catalogIntegrationsTests.Alter.
		withModify(
			case_CatalogIntegrations_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterCatalogIntegrationOptions) {
				opts.Set = &CatalogIntegrationSet{Enabled: new(true)}
				opts.SetTags = []TagAssociation{{Name: tagId, Value: "tag_value"}}
			},
		).
		withModifyAndExpectedSqlf(
			case_CatalogIntegrations_sql_Alter_Set,
			func(opts *AlterCatalogIntegrationOptions) {
				opts.Set = &CatalogIntegrationSet{Enabled: new(true)}
			},
			"ALTER CATALOG INTEGRATION %s SET ENABLED = true",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CatalogIntegrations_sql_Alter_SetTags,
			func(opts *AlterCatalogIntegrationOptions) {
				opts.SetTags = []TagAssociation{{Name: tagId, Value: "tag_value"}}
			},
			`ALTER CATALOG INTEGRATION %s SET TAG %s = 'tag_value'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CatalogIntegrations_sql_Alter_UnsetTags,
			func(opts *AlterCatalogIntegrationOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId}
			},
			`ALTER CATALOG INTEGRATION %s UNSET TAG %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_bearer",
			func(opts *AlterCatalogIntegrationOptions) {
				opts.IfExists = new(true)
				opts.Set = &CatalogIntegrationSet{
					SetBearerRestAuthentication: &SetBearerRestAuthentication{
						BearerToken: "test-token",
					},
					Enabled:                new(true),
					RefreshIntervalSeconds: new(60),
					Comment:                &StringAllowEmpty{Value: "test comment"},
				}
			},
			"ALTER CATALOG INTEGRATION IF EXISTS %s SET "+
				"REST_AUTHENTICATION = (BEARER_TOKEN = 'test-token') "+
				"ENABLED = true "+
				"REFRESH_INTERVAL_SECONDS = %d "+
				"COMMENT = '%s'",
			id.FullyQualifiedName(), 60, "test comment",
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_oauth",
			func(opts *AlterCatalogIntegrationOptions) {
				opts.IfExists = new(true)
				opts.Set = &CatalogIntegrationSet{
					SetOAuthRestAuthentication: &SetOAuthRestAuthentication{
						OauthClientSecret: oAuthClientSecret,
					},
					Enabled:                new(true),
					RefreshIntervalSeconds: new(60),
					Comment:                &StringAllowEmpty{Value: "test comment"},
				}
			},
			"ALTER CATALOG INTEGRATION IF EXISTS %s SET "+
				"REST_AUTHENTICATION = (OAUTH_CLIENT_SECRET = '%s') "+
				"ENABLED = true "+
				"REFRESH_INTERVAL_SECONDS = %d "+
				"COMMENT = '%s'",
			id.FullyQualifiedName(), oAuthClientSecret, 60, "test comment",
		)

	catalogIntegrationsTests.Drop.
		withExpectedSqlf(
			case_CatalogIntegrations_sql_Drop_basic,
			"DROP CATALOG INTEGRATION %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CatalogIntegrations_sql_Drop_all,
			func(opts *DropCatalogIntegrationOptions) { opts.IfExists = new(true) },
			"DROP CATALOG INTEGRATION IF EXISTS %s", id.FullyQualifiedName(),
		)

	catalogIntegrationsTests.Show.
		withExpectedSql(case_CatalogIntegrations_sql_Show_basic, "SHOW CATALOG INTEGRATIONS").
		withModifyAndExpectedSqlf(
			case_CatalogIntegrations_sql_Show_all,
			func(opts *ShowCatalogIntegrationOptions) {
				opts.Like = &Like{Pattern: new("like-pattern")}
			},
			"SHOW CATALOG INTEGRATIONS LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_CatalogIntegrations_sql_Show_Like,
			func(opts *ShowCatalogIntegrationOptions) {
				opts.Like = &Like{Pattern: new("like-pattern")}
			},
			"SHOW CATALOG INTEGRATIONS LIKE 'like-pattern'",
		)

	catalogIntegrationsTests.Describe.
		withExpectedSqlf(
			case_CatalogIntegrations_sql_Describe_basic,
			"DESCRIBE CATALOG INTEGRATION %s", id.FullyQualifiedName(),
		)
}

func TestParseCommaSeparatedEnumMap(t *testing.T) {
	testCases := []struct {
		Name   string
		Value  string
		Result []string
	}{
		{
			Name:   "empty enum map",
			Value:  "{}",
			Result: []string{},
		},
		{
			Name:   "empty string",
			Value:  "",
			Result: []string{},
		},
		{
			Name:   "multiple elements",
			Value:  "{KEY=value, KEY2=value2}",
			Result: []string{"KEY=value", "KEY2=value2"},
		},
		{
			Name:   "multiple elements without curly braces",
			Value:  "KEY=value, KEY2=value2",
			Result: []string{"KEY=value", "KEY2=value2"},
		},
		{
			Name:   "multiple elements with nested arrays",
			Value:  "{KEY=value, KEY2=[INNER_KEY=value2, INNER_KEY2=value3]}",
			Result: []string{"KEY=value", "KEY2=[INNER_KEY=value2, INNER_KEY2=value3]"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := parseCommaSeparatedEnumMap(CatalogIntegrationProperty{Value: tc.Value})
			require.Equal(t, tc.Result, actual)
		})
	}
}
