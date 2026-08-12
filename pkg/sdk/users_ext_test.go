package sdk

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/stretchr/testify/require"
)

func init() {
	id := usersTestIdAccountObjectIdentifier

	tagId := randomSchemaObjectIdentifier()
	password := random.Password()
	loginName := random.String()
	defaultRoleId := randomAccountObjectIdentifier()
	defaultWarehouseId := randomAccountObjectIdentifier()
	var defaultNamespaceId ObjectIdentifier = randomDatabaseObjectIdentifier()

	passwordPolicyId := randomSchemaObjectIdentifier()
	sessionPolicyId := randomSchemaObjectIdentifier()
	authenticationPolicyId := randomSchemaObjectIdentifier()
	renameTargetId := randomAccountObjectIdentifier()
	setTagId1 := randomSchemaObjectIdentifier()
	setTagId2 := randomSchemaObjectIdentifierInSchema(setTagId1.SchemaId())
	unsetTagId1 := randomSchemaObjectIdentifier()
	unsetTagId2 := randomSchemaObjectIdentifier()
	alterPassword := random.Password()

	usersTests.Create.
		withExpectedSqlf(
			case_Users_sql_Create_basic,
			`CREATE USER %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Create_all,
			func(opts *CreateUserOptions) {
				opts.IfNotExists = new(true)
				opts.ObjectProperties = &UserObjectProperties{
					Password:              &password,
					LoginName:             &loginName,
					DefaultRole:           &defaultRoleId,
					DefaultNamespace:      &defaultNamespaceId,
					DefaultWarehouse:      &defaultWarehouseId,
					DefaultSecondaryRoles: &SecondaryRoles{All: new(true)},
				}
				opts.ObjectParameters = &UserObjectParameters{EnableUnredactedQuerySyntaxError: new(true)}
				opts.SessionParameters = &SessionParameters{Autocommit: new(true)}
				opts.With = new(true)
				opts.Tag = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			`CREATE USER IF NOT EXISTS %s PASSWORD = '%s' LOGIN_NAME = '%s' DEFAULT_WAREHOUSE = %s DEFAULT_NAMESPACE = %s DEFAULT_ROLE = %s DEFAULT_SECONDARY_ROLES = ('ALL') ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true AUTOCOMMIT = true WITH TAG (%s = 'v1')`,
			id.FullyQualifiedName(), password, loginName, defaultWarehouseId.FullyQualifiedName(), defaultNamespaceId.FullyQualifiedName(), defaultRoleId.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_withType",
			func(opts *CreateUserOptions) {
				opts.ObjectProperties = &UserObjectProperties{UserType: Pointer(UserTypeLegacyService)}
			},
			`CREATE USER %s TYPE = LEGACY_SERVICE`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_emptySecondaryRoles",
			func(opts *CreateUserOptions) {
				opts.ObjectProperties = &UserObjectProperties{DefaultSecondaryRoles: &SecondaryRoles{None: new(true)}}
			},
			`CREATE USER %s DEFAULT_SECONDARY_ROLES = ()`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_workloadIdentityGcp",
			func(opts *CreateUserOptions) {
				opts.ObjectProperties = &UserObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						GcpType: &UserObjectWorkloadIdentityGcp{Subject: new("system:serviceaccount:service_account_namespace:service_account_name")},
					},
				}
			},
			`CREATE USER %s WORKLOAD_IDENTITY = (TYPE = %s SUBJECT = '%s')`, id.FullyQualifiedName(), WIFTypeGcp, "system:serviceaccount:service_account_namespace:service_account_name",
		).
		withAdditionalSqlCasef(
			"sql_Create_workloadIdentityAzure",
			func(opts *CreateUserOptions) {
				opts.ObjectProperties = &UserObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						AzureType: &UserObjectWorkloadIdentityAzure{
							Issuer:  new("https://accounts.google.com"),
							Subject: new("system:serviceaccount:service_account_namespace:service_account_name"),
						},
					},
				}
			},
			`CREATE USER %s WORKLOAD_IDENTITY = (TYPE = %s ISSUER = '%s' SUBJECT = '%s')`, id.FullyQualifiedName(), WIFTypeAzure, "https://accounts.google.com", "system:serviceaccount:service_account_namespace:service_account_name",
		).
		withAdditionalSqlCasef(
			"sql_Create_workloadIdentityAws",
			func(opts *CreateUserOptions) {
				opts.ObjectProperties = &UserObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						AwsType: &UserObjectWorkloadIdentityAws{Arn: new("arn:aws:iam::123456789012:role/test-role")},
					},
				}
			},
			`CREATE USER %s WORKLOAD_IDENTITY = (TYPE = %s ARN = '%s')`, id.FullyQualifiedName(), WIFTypeAws, "arn:aws:iam::123456789012:role/test-role",
		).
		withAdditionalSqlCasef(
			"sql_Create_workloadIdentityAwsWithIssuer",
			func(opts *CreateUserOptions) {
				opts.ObjectProperties = &UserObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						AwsType: &UserObjectWorkloadIdentityAws{
							Arn:    new("arn:aws:iam::123456789012:role/test-role"),
							Issuer: new("https://sts.amazonaws.com"),
						},
					},
				}
			},
			`CREATE USER %s WORKLOAD_IDENTITY = (TYPE = %s ARN = '%s' ISSUER = '%s')`, id.FullyQualifiedName(), WIFTypeAws, "arn:aws:iam::123456789012:role/test-role", "https://sts.amazonaws.com",
		).
		withAdditionalSqlCasef(
			"sql_Create_workloadIdentityOidcBasic",
			func(opts *CreateUserOptions) {
				opts.ObjectProperties = &UserObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						OidcType: &UserObjectWorkloadIdentityOidc{
							Issuer:  new("https://accounts.google.com"),
							Subject: new("system:serviceaccount:service_account_namespace:service_account_name"),
						},
					},
				}
			},
			`CREATE USER %s WORKLOAD_IDENTITY = (TYPE = %s ISSUER = '%s' SUBJECT = '%s')`, id.FullyQualifiedName(), WIFTypeOidc, "https://accounts.google.com", "system:serviceaccount:service_account_namespace:service_account_name",
		).
		withAdditionalSqlCasef(
			"sql_Create_workloadIdentityOidcComplete",
			func(opts *CreateUserOptions) {
				opts.ObjectProperties = &UserObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						OidcType: &UserObjectWorkloadIdentityOidc{
							Issuer:  new("https://accounts.google.com"),
							Subject: new("system:serviceaccount:service_account_namespace:service_account_name"),
							OidcAudienceList: []StringListItemWrapper{
								{Value: "https://accounts.google.com/o/oauth2/auth"},
							},
						},
					},
				}
			},
			`CREATE USER %s WORKLOAD_IDENTITY = (TYPE = %s ISSUER = '%s' SUBJECT = '%s' OIDC_AUDIENCE_LIST = ('https://accounts.google.com/o/oauth2/auth'))`, id.FullyQualifiedName(), WIFTypeOidc, "https://accounts.google.com", "system:serviceaccount:service_account_namespace:service_account_name",
		)

	usersTests.Alter.
		withAdditionalValidationCase(
			"validation_Alter_Set_policyWithPropertiesOrParameters",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{
					AuthenticationPolicy: &authenticationPolicyId,
					ObjectParameters:     &UserObjectParameters{EnableUnredactedQuerySyntaxError: new(true)},
				}
			},
			NewError("policies cannot be set with user properties or parameters at the same time"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_forceWithoutPolicy",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{
					ObjectProperties: &UserAlterObjectProperties{LoginName: new("some_login")},
					Force:            new(true),
				}
			},
			NewError("force can only be set with PasswordPolicy, SessionPolicy, or AuthenticationPolicy"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Unset_policyWithPropertiesOrParameters",
			func(opts *AlterUserOptions) {
				opts.Unset = &UserUnset{
					PasswordPolicy:   new(true),
					ObjectParameters: &UserObjectParametersUnset{EnableUnredactedQuerySyntaxError: new(true)},
				}
			},
			NewError("policies cannot be unset with user properties or parameters at the same time"),
		).
		withModify(
			case_Users_validation_Alter_RemoveDelegatedAuthorization_Integration_ValidateValueSet,
			func(opts *AlterUserOptions) {
				opts.RemoveDelegatedAuthorization = &RemoveDelegatedAuthorization{
					RemoveDelegatedAuthorizationOfRole: new("ROLE1"),
				}
			},
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Alter_RenameTo,
			func(opts *AlterUserOptions) { opts.RenameTo = &renameTargetId },
			`ALTER USER %s RENAME TO %s`, id.FullyQualifiedName(), renameTargetId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Alter_ResetPassword,
			func(opts *AlterUserOptions) { opts.ResetPassword = new(true) },
			`ALTER USER %s RESET PASSWORD`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Alter_AbortAllQueries,
			func(opts *AlterUserOptions) { opts.AbortAllQueries = new(true) },
			`ALTER USER %s ABORT ALL QUERIES`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Alter_AddDelegatedAuthorization,
			func(opts *AlterUserOptions) {
				opts.AddDelegatedAuthorization = &AddDelegatedAuthorization{Role: "ROLE1", Integration: "INTEGRATION1"}
			},
			`ALTER USER %s ADD DELEGATED AUTHORIZATION OF ROLE %s TO SECURITY INTEGRATION %s`, id.FullyQualifiedName(), "ROLE1", "INTEGRATION1",
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Alter_RemoveDelegatedAuthorization,
			func(opts *AlterUserOptions) {
				opts.RemoveDelegatedAuthorization = &RemoveDelegatedAuthorization{
					RemoveDelegatedAuthorizationOfRole: new("ROLE1"),
					Integration:                        "INTEGRATION1",
				}
			},
			`ALTER USER %s REMOVE DELEGATED AUTHORIZATION OF ROLE %s FROM SECURITY INTEGRATION %s`, id.FullyQualifiedName(), "ROLE1", "INTEGRATION1",
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Alter_SetTags,
			func(opts *AlterUserOptions) {
				opts.SetTags = []TagAssociation{
					{Name: setTagId1, Value: "v1"},
					{Name: setTagId2, Value: "v2"},
				}
			},
			`ALTER USER %s SET TAG %s = 'v1', %s = 'v2'`, id.FullyQualifiedName(), setTagId1.FullyQualifiedName(), setTagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Alter_UnsetTags,
			func(opts *AlterUserOptions) { opts.UnsetTags = []ObjectIdentifier{unsetTagId1, unsetTagId2} },
			`ALTER USER %s UNSET TAG %s, %s`, id.FullyQualifiedName(), unsetTagId1.FullyQualifiedName(), unsetTagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Alter_Set,
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{
					SessionParameters: &SessionParameters{AbortDetachedQuery: new(true)},
					ObjectParameters:  &UserObjectParameters{EnableUnredactedQuerySyntaxError: new(true)},
				}
			},
			`ALTER USER %s SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true ABORT_DETACHED_QUERY = true`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_UserType",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{ObjectProperties: &UserAlterObjectProperties{UserType: Pointer(UserTypeLegacyService)}}
			},
			`ALTER USER %s SET TYPE = LEGACY_SERVICE`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_PasswordPolicy",
			func(opts *AlterUserOptions) { opts.Set = &UserSet{PasswordPolicy: &passwordPolicyId} },
			`ALTER USER %s SET PASSWORD POLICY %s`, id.FullyQualifiedName(), passwordPolicyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_AuthenticationPolicy",
			func(opts *AlterUserOptions) { opts.Set = &UserSet{AuthenticationPolicy: &authenticationPolicyId} },
			`ALTER USER %s SET AUTHENTICATION POLICY %s`, id.FullyQualifiedName(), authenticationPolicyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_SessionPolicy",
			func(opts *AlterUserOptions) { opts.Set = &UserSet{SessionPolicy: &sessionPolicyId} },
			`ALTER USER %s SET SESSION POLICY %s`, id.FullyQualifiedName(), sessionPolicyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_PasswordPolicy_withForce",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{PasswordPolicy: &passwordPolicyId, Force: new(true)}
			},
			`ALTER USER %s SET PASSWORD POLICY %s FORCE`, id.FullyQualifiedName(), passwordPolicyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_SessionPolicy_withForce",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{SessionPolicy: &sessionPolicyId, Force: new(true)}
			},
			`ALTER USER %s SET SESSION POLICY %s FORCE`, id.FullyQualifiedName(), sessionPolicyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_AuthenticationPolicy_withForce",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{AuthenticationPolicy: &authenticationPolicyId, Force: new(true)}
			},
			`ALTER USER %s SET AUTHENTICATION POLICY %s FORCE`, id.FullyQualifiedName(), authenticationPolicyId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_ObjectProperties",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{
					ObjectProperties: &UserAlterObjectProperties{
						Password:              &alterPassword,
						DefaultSecondaryRoles: &SecondaryRoles{All: new(true)},
					},
				}
			},
			`ALTER USER %s SET PASSWORD = '%s' DEFAULT_SECONDARY_ROLES = ('ALL')`, id.FullyQualifiedName(), alterPassword,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_ObjectParameters",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{ObjectParameters: &UserObjectParameters{EnableUnredactedQuerySyntaxError: new(true)}}
			},
			`ALTER USER %s SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_SessionParameters",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{SessionParameters: &SessionParameters{Autocommit: new(true)}}
			},
			`ALTER USER %s SET AUTOCOMMIT = true`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_FirstNameAndDisableMfa",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{ObjectProperties: &UserAlterObjectProperties{FirstName: new("name"), DisableMfa: new(true)}}
			},
			`ALTER USER %s SET FIRST_NAME = '%s' DISABLE_MFA = true`, id.FullyQualifiedName(), "name",
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_DisableMfaOnly",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{ObjectProperties: &UserAlterObjectProperties{DisableMfa: new(true)}}
			},
			`ALTER USER %s SET DISABLE_MFA = true`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_WorkloadIdentityGcp",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{ObjectProperties: &UserAlterObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						GcpType: &UserObjectWorkloadIdentityGcp{Subject: new("system:serviceaccount:service_account_namespace:service_account_name")},
					},
				}}
			},
			`ALTER USER %s SET WORKLOAD_IDENTITY = (TYPE = %s SUBJECT = '%s')`, id.FullyQualifiedName(), WIFTypeGcp, "system:serviceaccount:service_account_namespace:service_account_name",
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_WorkloadIdentityAzure",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{ObjectProperties: &UserAlterObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						AzureType: &UserObjectWorkloadIdentityAzure{
							Issuer:  new("https://accounts.google.com"),
							Subject: new("system:serviceaccount:service_account_namespace:service_account_name"),
						},
					},
				}}
			},
			`ALTER USER %s SET WORKLOAD_IDENTITY = (TYPE = %s ISSUER = '%s' SUBJECT = '%s')`, id.FullyQualifiedName(), WIFTypeAzure, "https://accounts.google.com", "system:serviceaccount:service_account_namespace:service_account_name",
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_WorkloadIdentityAws",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{ObjectProperties: &UserAlterObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						AwsType: &UserObjectWorkloadIdentityAws{Arn: new("arn:aws:iam::123456789012:role/test-role")},
					},
				}}
			},
			`ALTER USER %s SET WORKLOAD_IDENTITY = (TYPE = %s ARN = '%s')`, id.FullyQualifiedName(), WIFTypeAws, "arn:aws:iam::123456789012:role/test-role",
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_WorkloadIdentityAwsWithIssuer",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{ObjectProperties: &UserAlterObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						AwsType: &UserObjectWorkloadIdentityAws{
							Arn:    new("arn:aws:iam::123456789012:role/test-role"),
							Issuer: new("https://sts.amazonaws.com"),
						},
					},
				}}
			},
			`ALTER USER %s SET WORKLOAD_IDENTITY = (TYPE = %s ARN = '%s' ISSUER = '%s')`, id.FullyQualifiedName(), WIFTypeAws, "arn:aws:iam::123456789012:role/test-role", "https://sts.amazonaws.com",
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_WorkloadIdentityOidcBasic",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{ObjectProperties: &UserAlterObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						OidcType: &UserObjectWorkloadIdentityOidc{
							Issuer:  new("https://accounts.google.com"),
							Subject: new("system:serviceaccount:service_account_namespace:service_account_name"),
						},
					},
				}}
			},
			`ALTER USER %s SET WORKLOAD_IDENTITY = (TYPE = %s ISSUER = '%s' SUBJECT = '%s')`, id.FullyQualifiedName(), WIFTypeOidc, "https://accounts.google.com", "system:serviceaccount:service_account_namespace:service_account_name",
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_WorkloadIdentityOidcComplete",
			func(opts *AlterUserOptions) {
				opts.Set = &UserSet{ObjectProperties: &UserAlterObjectProperties{
					WorkloadIdentity: &UserObjectWorkloadIdentityProperties{
						OidcType: &UserObjectWorkloadIdentityOidc{
							Issuer:  new("https://accounts.google.com"),
							Subject: new("system:serviceaccount:service_account_namespace:service_account_name"),
							OidcAudienceList: []StringListItemWrapper{
								{Value: "https://accounts.google.com/o/oauth2/auth"},
							},
						},
					},
				}}
			},
			`ALTER USER %s SET WORKLOAD_IDENTITY = (TYPE = %s ISSUER = '%s' SUBJECT = '%s' OIDC_AUDIENCE_LIST = ('https://accounts.google.com/o/oauth2/auth'))`, id.FullyQualifiedName(), WIFTypeOidc, "https://accounts.google.com", "system:serviceaccount:service_account_namespace:service_account_name",
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Alter_Unset,
			func(opts *AlterUserOptions) {
				opts.Unset = &UserUnset{
					ObjectParameters:  &UserObjectParametersUnset{EnableUnredactedQuerySyntaxError: new(true)},
					SessionParameters: &SessionParametersUnset{BinaryOutputFormat: new(true)},
				}
			},
			`ALTER USER %s UNSET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR, BINARY_OUTPUT_FORMAT`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_UserType",
			func(opts *AlterUserOptions) {
				opts.Unset = &UserUnset{ObjectProperties: &UserObjectPropertiesUnset{UserType: new(true)}}
			},
			`ALTER USER %s UNSET TYPE`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_Password",
			func(opts *AlterUserOptions) {
				opts.Unset = &UserUnset{ObjectProperties: &UserObjectPropertiesUnset{Password: new(true)}}
			},
			`ALTER USER %s UNSET PASSWORD`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_SessionPolicy",
			func(opts *AlterUserOptions) { opts.Unset = &UserUnset{SessionPolicy: new(true)} },
			`ALTER USER %s UNSET SESSION POLICY`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_AuthenticationPolicy",
			func(opts *AlterUserOptions) { opts.Unset = &UserUnset{AuthenticationPolicy: new(true)} },
			`ALTER USER %s UNSET AUTHENTICATION POLICY`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_WorkloadIdentity",
			func(opts *AlterUserOptions) {
				opts.Unset = &UserUnset{ObjectProperties: &UserObjectPropertiesUnset{WorkloadIdentity: new(true)}}
			},
			`ALTER USER %s UNSET WORKLOAD_IDENTITY`, id.FullyQualifiedName(),
		)

	usersTests.Drop.
		withExpectedSqlf(case_Users_sql_Drop_basic, `DROP USER %s`, id.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_Users_sql_Drop_all,
			func(opts *DropUserOptions) { opts.IfExists = new(true) },
			`DROP USER IF EXISTS %s`, id.FullyQualifiedName(),
		)

	usersTests.Show.
		withAdditionalValidationCase(
			"validation_Show_fromWithoutLimit",
			func(opts *ShowUserOptions) {
				opts.Limit = &LimitFrom{From: new("from_pattern")}
			},
			errNotSet("ShowUserOptions.Limit", "Rows"),
		).
		withExpectedSqlf(case_Users_sql_Show_basic, `SHOW USERS`).
		withModifyAndExpectedSqlf(
			case_Users_sql_Show_all,
			func(opts *ShowUserOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.StartsWith = new("prefix")
				opts.Limit = &LimitFrom{Rows: new(5), From: new("from")}
			},
			`SHOW USERS LIKE '%s' STARTS WITH '%s' LIMIT %v FROM '%s'`, "pattern", "prefix", 5, "from",
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Show_Like,
			func(opts *ShowUserOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			`SHOW USERS LIKE '%s'`, "pattern",
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Show_StartsWith,
			func(opts *ShowUserOptions) { opts.StartsWith = new("prefix") },
			`SHOW USERS STARTS WITH '%s'`, "prefix",
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_Show_Limit,
			func(opts *ShowUserOptions) { opts.Limit = &LimitFrom{Rows: new(5), From: new("from")} },
			`SHOW USERS LIMIT %v FROM '%s'`, 5, "from",
		).
		withAdditionalSqlCasef(
			"sql_Show_LikeAndLimit",
			func(opts *ShowUserOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.Limit = &LimitFrom{Rows: new(5)}
			},
			`SHOW USERS LIKE '%s' LIMIT %v`, "pattern", 5,
		).
		withAdditionalSqlCasef(
			"sql_Show_LikeAndLimitAndFrom",
			func(opts *ShowUserOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.Limit = &LimitFrom{Rows: new(5), From: new("from")}
			},
			`SHOW USERS LIKE '%s' LIMIT %v FROM '%s'`, "pattern", 5, "from",
		).
		withAdditionalSqlCasef(
			"sql_Show_StartsWithAndLimitAndFrom",
			func(opts *ShowUserOptions) {
				opts.StartsWith = new("prefix")
				opts.Limit = &LimitFrom{Rows: new(5), From: new("from")}
			},
			`SHOW USERS STARTS WITH '%s' LIMIT %v FROM '%s'`, "prefix", 5, "from",
		)

	usersTests.Describe.
		withExpectedSqlf(case_Users_sql_Describe_basic, `DESCRIBE USER %s`, id.FullyQualifiedName())

	usersTests.ShowUserWorkloadIdentityAuthenticationMethodOptions.
		withDefaultOpts(func() *ShowUserWorkloadIdentityAuthenticationMethodOptionsUserOptions {
			return &ShowUserWorkloadIdentityAuthenticationMethodOptionsUserOptions{
				ForUser: id,
			}
		}).
		withExpectedSqlf(
			case_Users_sql_ShowUserWorkloadIdentityAuthenticationMethodOptions_basic,
			`SHOW USER WORKLOAD IDENTITY AUTHENTICATION METHODS FOR USER %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Users_sql_ShowUserWorkloadIdentityAuthenticationMethodOptions_all,
			func(opts *ShowUserWorkloadIdentityAuthenticationMethodOptionsUserOptions) {},
			`SHOW USER WORKLOAD IDENTITY AUTHENTICATION METHODS FOR USER %s`, id.FullyQualifiedName(),
		)
}

func Test_User_ToGeographyOutputFormat(t *testing.T) {
	type test struct {
		input string
		want  GeographyOutputFormat
	}

	valid := []test{
		// case insensitive.
		{input: "geojson", want: GeographyOutputFormatGeoJSON},

		// Supported Values
		{input: string(GeographyOutputFormatGeoJSON), want: GeographyOutputFormatGeoJSON},
		{input: string(GeographyOutputFormatWKT), want: GeographyOutputFormatWKT},
		{input: string(GeographyOutputFormatWKB), want: GeographyOutputFormatWKB},
		{input: string(GeographyOutputFormatEWKT), want: GeographyOutputFormatEWKT},
		{input: string(GeographyOutputFormatEWKB), want: GeographyOutputFormatEWKB},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},

		// not supported values (single-quoted)
		{input: "'GeoJSON'"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToGeographyOutputFormat(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToGeographyOutputFormat(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_User_ToGeometryOutputFormat(t *testing.T) {
	type test struct {
		input string
		want  GeometryOutputFormat
	}

	valid := []test{
		// case insensitive.
		{input: "geojson", want: GeometryOutputFormatGeoJSON},

		// Supported Values
		{input: string(GeometryOutputFormatGeoJSON), want: GeometryOutputFormatGeoJSON},
		{input: string(GeometryOutputFormatWKT), want: GeometryOutputFormatWKT},
		{input: string(GeometryOutputFormatWKB), want: GeometryOutputFormatWKB},
		{input: string(GeometryOutputFormatEWKT), want: GeometryOutputFormatEWKT},
		{input: string(GeometryOutputFormatEWKB), want: GeometryOutputFormatEWKB},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},

		// not supported values (single-quoted)
		{input: "'GeoJSON'"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToGeometryOutputFormat(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToGeometryOutputFormat(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_User_ToBinaryInputFormat(t *testing.T) {
	type test struct {
		input string
		want  BinaryInputFormat
	}

	valid := []test{
		// case insensitive.
		{input: "hex", want: BinaryInputFormatHex},

		// Supported Values
		{input: string(BinaryInputFormatHex), want: BinaryInputFormatHex},
		{input: string(BinaryInputFormatBase64), want: BinaryInputFormatBase64},
		{input: string(BinaryInputFormatUTF8), want: BinaryInputFormatUTF8},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},

		// not supported values (single-quoted)
		{input: "'HEX'"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToBinaryInputFormat(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToBinaryInputFormat(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_User_ToBinaryOutputFormat(t *testing.T) {
	type test struct {
		input string
		want  BinaryOutputFormat
	}

	valid := []test{
		// case insensitive.
		{input: "hex", want: BinaryOutputFormatHex},

		// Supported Values
		{input: string(BinaryOutputFormatHex), want: BinaryOutputFormatHex},
		{input: string(BinaryOutputFormatBase64), want: BinaryOutputFormatBase64},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},

		// not supported values (single-quoted)
		{input: "'HEX'"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToBinaryOutputFormat(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToBinaryOutputFormat(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_User_ToClientTimestampTypeMapping(t *testing.T) {
	type test struct {
		input string
		want  ClientTimestampTypeMapping
	}

	valid := []test{
		// case insensitive.
		{input: "timestamp_ltz", want: ClientTimestampTypeMappingLtz},

		// Supported Values
		{input: string(ClientTimestampTypeMappingLtz), want: ClientTimestampTypeMappingLtz},
		{input: string(ClientTimestampTypeMappingNtz), want: ClientTimestampTypeMappingNtz},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},

		// not supported values (single-quoted)
		{input: "'TIMESTAMP_LTZ'"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToClientTimestampTypeMapping(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToClientTimestampTypeMapping(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_User_ToTimestampTypeMapping(t *testing.T) {
	type test struct {
		input string
		want  TimestampTypeMapping
	}

	valid := []test{
		// case insensitive.
		{input: "timestamp_ltz", want: TimestampTypeMappingLtz},

		// Supported Values
		{input: string(TimestampTypeMappingLtz), want: TimestampTypeMappingLtz},
		{input: string(TimestampTypeMappingNtz), want: TimestampTypeMappingNtz},
		{input: string(TimestampTypeMappingTz), want: TimestampTypeMappingTz},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},

		// not supported values (single-quoted)
		{input: "'TIMESTAMP_LTZ'"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToTimestampTypeMapping(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToTimestampTypeMapping(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_User_ToTransactionDefaultIsolationLevel(t *testing.T) {
	type test struct {
		input string
		want  TransactionDefaultIsolationLevel
	}

	valid := []test{
		// case insensitive.
		{input: "read committed", want: TransactionDefaultIsolationLevelReadCommitted},

		// Supported Values
		{input: string(TransactionDefaultIsolationLevelReadCommitted), want: TransactionDefaultIsolationLevelReadCommitted},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},

		// not supported values (single-quoted)
		{input: "'READ COMMITTED'"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToTransactionDefaultIsolationLevel(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToTransactionDefaultIsolationLevel(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_User_ToUnsupportedDDLAction(t *testing.T) {
	type test struct {
		input string
		want  UnsupportedDDLAction
	}

	valid := []test{
		// case insensitive.
		{input: "ignore", want: UnsupportedDDLActionIgnore},

		// Supported Values
		{input: string(UnsupportedDDLActionIgnore), want: UnsupportedDDLActionIgnore},
		{input: string(UnsupportedDDLActionFail), want: UnsupportedDDLActionFail},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},

		// not supported values (single-quoted)
		{input: "'IGNORE'"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToUnsupportedDDLAction(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToUnsupportedDDLAction(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_User_ToSecondaryRolesOption(t *testing.T) {
	type test struct {
		input string
		want  SecondaryRolesOption
	}

	valid := []test{
		// case insensitive.
		{input: "none", want: SecondaryRolesOptionNone},

		// Supported Values
		{input: "NONE", want: SecondaryRolesOptionNone},
		{input: "ALL", want: SecondaryRolesOptionAll},
		{input: "DEFAULT", want: SecondaryRolesOptionDefault},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToSecondaryRolesOption(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToSecondaryRolesOption(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_User_GetSecondaryRolesOptionFrom(t *testing.T) {
	type test struct {
		input string
		want  SecondaryRolesOption
	}

	valid := []test{
		{input: "", want: SecondaryRolesOptionDefault},
		{input: "[]", want: SecondaryRolesOptionNone},
		{input: `["ALL"]`, want: SecondaryRolesOptionAll},
		{input: `["any"]`, want: SecondaryRolesOptionAll},
		{input: `["more", "than", "one"]`, want: SecondaryRolesOptionAll},
		{input: `no list`, want: SecondaryRolesOptionAll},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got := GetSecondaryRolesOptionFrom(tc.input)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range valid {
		t.Run(fmt.Sprintf("invoked from user: %s", tc.input), func(t *testing.T) {
			user := User{DefaultSecondaryRoles: tc.input}
			got := user.GetSecondaryRolesOption()
			require.Equal(t, tc.want, got)
		})
	}
}

func Test_User_ToUserType(t *testing.T) {
	type test struct {
		input string
		want  UserType
	}

	valid := []test{
		// case insensitive.
		{input: "person", want: UserTypePerson},

		// Supported Values
		{input: "PERSON", want: UserTypePerson},
		{input: "SERVICE", want: UserTypeService},
		{input: "LEGACY_SERVICE", want: UserTypeLegacyService},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
		{input: "legacyservice"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToUserType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToUserType(tc.input)
			require.Error(t, err)
		})
	}
}
