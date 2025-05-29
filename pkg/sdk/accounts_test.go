package sdk

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountCreate(t *testing.T) {
	t.Run("simplest case", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		password := random.Password()
		opts := &CreateAccountOptions{
			name:          id,
			AdminName:     "someadmin",
			AdminPassword: String(password),
			Email:         "admin@example.com",
			Edition:       EditionBusinessCritical,
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE ACCOUNT %s ADMIN_NAME = 'someadmin' ADMIN_PASSWORD = '%s' EMAIL = 'admin@example.com' EDITION = BUSINESS_CRITICAL`, id.FullyQualifiedName(), password)
	})

	t.Run("every option", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		key := random.Password()
		opts := &CreateAccountOptions{
			name:               id,
			AdminName:          "someadmin",
			AdminRSAPublicKey:  String(key),
			AdminUserType:      Pointer(UserTypeService),
			FirstName:          String("Ad"),
			LastName:           String("Min"),
			Email:              "admin@example.com",
			MustChangePassword: Bool(true),
			Edition:            EditionBusinessCritical,
			RegionGroup:        String("groupid"),
			Region:             String("regionid"),
			Comment:            String("Test account"),
			Polaris:            Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE ACCOUNT %s ADMIN_NAME = 'someadmin' ADMIN_RSA_PUBLIC_KEY = '%s' ADMIN_USER_TYPE = SERVICE FIRST_NAME = 'Ad' LAST_NAME = 'Min' EMAIL = 'admin@example.com' MUST_CHANGE_PASSWORD = true EDITION = BUSINESS_CRITICAL REGION_GROUP = groupid REGION = regionid COMMENT = 'Test account' POLARIS = true`, id.FullyQualifiedName(), key)
	})

	t.Run("static password", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		password := random.Password()
		opts := &CreateAccountOptions{
			name:               id,
			AdminName:          "someadmin",
			AdminPassword:      String(password),
			FirstName:          String("Ad"),
			LastName:           String("Min"),
			Email:              "admin@example.com",
			MustChangePassword: Bool(false),
			Edition:            EditionBusinessCritical,
			RegionGroup:        String("groupid"),
			Region:             String("regionid"),
			Comment:            String("Test account"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE ACCOUNT %s ADMIN_NAME = 'someadmin' ADMIN_PASSWORD = '%s' FIRST_NAME = 'Ad' LAST_NAME = 'Min' EMAIL = 'admin@example.com' MUST_CHANGE_PASSWORD = false EDITION = BUSINESS_CRITICAL REGION_GROUP = groupid REGION = regionid COMMENT = 'Test account'`, id.FullyQualifiedName(), password)
	})
}

func TestAccountAlter(t *testing.T) {
	t.Run("validation: exactly one value set in AccountSet - nothing set", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Set: &AccountSet{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AccountSet", "LegacyParameters", "ResourceMonitor", "PackagesPolicy", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy"))
	})

	t.Run("validation: exactly one value set in AccountSet - multiple set", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				PasswordPolicy:       randomSchemaObjectIdentifier(),
				SessionPolicy:        randomSchemaObjectIdentifier(),
				AuthenticationPolicy: randomSchemaObjectIdentifier(),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AccountSet", "LegacyParameters", "ResourceMonitor", "PackagesPolicy", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy"))
	})

	t.Run("validation: exactly one value set in AccountUnset - nothing set", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AccountUnset", "LegacyParameters", "PackagesPolicy", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ResourceMonitor"))
	})

	t.Run("validation: exactly one value set in AccountUnset - multiple set", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				PasswordPolicy:       Bool(true),
				SessionPolicy:        Bool(true),
				AuthenticationPolicy: Bool(true),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AccountUnset", "LegacyParameters", "PackagesPolicy", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ResourceMonitor"))
	})

	t.Run("with set params", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				LegacyParameters: &AccountLevelParameters{
					AccountParameters: &AccountParameters{
						ClientEncryptionKeySize:       Int(128),
						PreventUnloadToInternalStages: Bool(true),
					},
					SessionParameters: &SessionParameters{
						JsonIndent: Int(16),
					},
					ObjectParameters: &ObjectParameters{
						MaxDataExtensionTimeInDays: Int(30),
					},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET CLIENT_ENCRYPTION_KEY_SIZE = 128, PREVENT_UNLOAD_TO_INTERNAL_STAGES = true, JSON_INDENT = 16, MAX_DATA_EXTENSION_TIME_IN_DAYS = 30`)
	})

	t.Run("with unset params", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				Parameters: &AccountLevelParametersUnset{
					AccountParameters: &AccountParametersUnset{
						InitialReplicationSizeLimitInTB: Bool(true),
						SSOLoginPage:                    Bool(true),
					},
					SessionParameters: &SessionParametersUnset{
						SimulatedDataSharingConsumer: Bool(true),
						Timezone:                     Bool(true),
					},
					ObjectParameters: &ObjectParametersUnset{
						DefaultDDLCollation: Bool(true),
					},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET INITIAL_REPLICATION_SIZE_LIMIT_IN_TB, SSO_LOGIN_PAGE, SIMULATED_DATA_SHARING_CONSUMER, TIMEZONE, DEFAULT_DDL_COLLATION`)
	})

	t.Run("with set resource monitor", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				ResourceMonitor: NewAccountObjectIdentifier("mymonitor"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET RESOURCE_MONITOR = "mymonitor"`)
	})

	t.Run("with set packages policy", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				PackagesPolicy: id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET PACKAGES POLICY %s`, id.FullyQualifiedName())
	})

	t.Run("with set packages policy with force", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				PackagesPolicy: id,
				Force:          Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET PACKAGES POLICY %s FORCE`, id.FullyQualifiedName())
	})

	t.Run("validate: force with other policy than packages", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				PasswordPolicy: id,
				Force:          Bool(true),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("force can only be set with PackagesPolicy field"))
	})

	t.Run("with set password policy", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				PasswordPolicy: id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET PASSWORD POLICY %s`, id.FullyQualifiedName())
	})

	t.Run("with set session policy", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				SessionPolicy: id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET SESSION POLICY %s`, id.FullyQualifiedName())
	})

	t.Run("with set authentication policy", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &AlterAccountOptions{
			Set: &AccountSet{
				AuthenticationPolicy: id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET AUTHENTICATION POLICY %s`, id.FullyQualifiedName())
	})

	t.Run("with unset packages policy", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				PackagesPolicy: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET PACKAGES POLICY`)
	})

	t.Run("with unset password policy", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				PasswordPolicy: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET PASSWORD POLICY`)
	})

	t.Run("with unset session policy", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				SessionPolicy: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET SESSION POLICY`)
	})

	t.Run("with unset authentication policy", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				AuthenticationPolicy: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET AUTHENTICATION POLICY`)
	})

	t.Run("with unset resource monitor", func(t *testing.T) {
		opts := &AlterAccountOptions{
			Unset: &AccountUnset{
				ResourceMonitor: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET RESOURCE_MONITOR`)
	})

	t.Run("with set tag", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		opts := &AlterAccountOptions{
			SetTag: []TagAssociation{
				{
					Name:  tagId1,
					Value: "v1",
				},
				{
					Name:  tagId2,
					Value: "v2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT SET TAG %s = 'v1', %s = 'v2'`, tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("with unset tag", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := &AlterAccountOptions{
			UnsetTag: []ObjectIdentifier{
				id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT UNSET TAG %s`, id.FullyQualifiedName())
	})

	t.Run("set is_org_admin", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &AlterAccountOptions{
			SetIsOrgAdmin: &AccountSetIsOrgAdmin{
				Name:     id,
				OrgAdmin: true,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT %s SET IS_ORG_ADMIN = true`, id.FullyQualifiedName())
	})

	t.Run("rename", func(t *testing.T) {
		oldName := randomAccountObjectIdentifier()
		newName := randomAccountObjectIdentifier()
		opts := &AlterAccountOptions{
			Rename: &AccountRename{
				Name:       oldName,
				NewName:    newName,
				SaveOldURL: Bool(false),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT %s RENAME TO %s SAVE_OLD_URL = false`, oldName.FullyQualifiedName(), newName.FullyQualifiedName())
	})

	t.Run("validation: drop no url set", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &AlterAccountOptions{
			Drop: &AccountDrop{
				Name: id,
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AccountDrop", "OldUrl", "OldOrganizationUrl"))
	})

	t.Run("validation: drop all url options set", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &AlterAccountOptions{
			Drop: &AccountDrop{
				Name:               id,
				OldUrl:             Bool(true),
				OldOrganizationUrl: Bool(true),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AccountDrop", "OldUrl", "OldOrganizationUrl"))
	})

	t.Run("drop old url", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &AlterAccountOptions{
			Drop: &AccountDrop{
				Name:   id,
				OldUrl: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT %s DROP OLD URL`, id.FullyQualifiedName())
	})

	t.Run("drop organization old url", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &AlterAccountOptions{
			Drop: &AccountDrop{
				Name:               id,
				OldOrganizationUrl: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER ACCOUNT %s DROP OLD ORGANIZATION URL`, id.FullyQualifiedName())
	})
}

func TestAccountDrop(t *testing.T) {
	t.Run("validate: empty options", func(t *testing.T) {
		opts := &DropAccountOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("minimal", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &DropAccountOptions{
			name:              id,
			gracePeriodInDays: 10,
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP ACCOUNT %s GRACE_PERIOD_IN_DAYS = 10`, id.FullyQualifiedName())
	})

	t.Run("if exists", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &DropAccountOptions{
			name:              id,
			IfExists:          Bool(true),
			gracePeriodInDays: 10,
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP ACCOUNT IF EXISTS %s GRACE_PERIOD_IN_DAYS = 10`, id.FullyQualifiedName())
	})
}

func TestAccountShow(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		opts := &ShowAccountOptions{}
		assertOptsValidAndSQLEquals(t, opts, `SHOW ACCOUNTS`)
	})

	t.Run("with history and like", func(t *testing.T) {
		opts := &ShowAccountOptions{
			History: Bool(true),
			Like: &Like{
				Pattern: String("myaccount"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW ACCOUNTS HISTORY LIKE 'myaccount'`)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowAccountOptions{
			Like: &Like{
				Pattern: String("myaccount"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW ACCOUNTS LIKE 'myaccount'`)
	})
}

func TestToAccountCreateResponse(t *testing.T) {
	testCases := []struct {
		Name           string
		RawInput       string
		Input          AccountCreateResponse
		ExpectedOutput *AccountCreateResponse
		Error          string
	}{
		{
			Name:     "validation: empty input",
			RawInput: "",
			Error:    "unexpected end of JSON input",
		},
		{
			Name: "validation: only a few fields filled",
			Input: AccountCreateResponse{
				AccountName: "acc_name",
				Url:         `https://org_name-acc_name.snowflakecomputing.com`,
				Edition:     EditionStandard,
				RegionGroup: "region_group",
				Cloud:       "cloud",
				Region:      "region",
			},
			ExpectedOutput: &AccountCreateResponse{
				AccountName:      "acc_name",
				Url:              `https://org_name-acc_name.snowflakecomputing.com`,
				OrganizationName: "ORG_NAME",
				Edition:          EditionStandard,
				RegionGroup:      "region_group",
				Cloud:            "cloud",
				Region:           "region",
			},
		},
		{
			Name: "validation: invalid url",
			Input: AccountCreateResponse{
				Url: `https://org_name_acc_name.snowflake.computing.com`,
			},
			ExpectedOutput: &AccountCreateResponse{
				Url: `https://org_name_acc_name.snowflake.computing.com`,
				// OrganizationName is not filled
			},
		},
		{
			Name: "validation: valid url",
			Input: AccountCreateResponse{
				Url: `https://org_name-acc_name.snowflakecomputing.com`,
			},
			ExpectedOutput: &AccountCreateResponse{
				Url:              `https://org_name-acc_name.snowflakecomputing.com`,
				OrganizationName: "ORG_NAME",
			},
		},
		{
			Name: "validation: valid http url",
			Input: AccountCreateResponse{
				Url: `http://org_name-acc_name.snowflakecomputing.com`,
			},
			ExpectedOutput: &AccountCreateResponse{
				Url:              `http://org_name-acc_name.snowflakecomputing.com`,
				OrganizationName: "ORG_NAME",
			},
		},
		{
			Name: "complete",
			Input: AccountCreateResponse{
				AccountLocator:    "locator",
				AccountLocatorUrl: "locator_url",
				AccountName:       "acc_name",
				Url:               `https://org_name-acc_name.snowflakecomputing.com`,
				Edition:           EditionBusinessCritical,
				RegionGroup:       "region_group",
				Cloud:             "cloud",
				Region:            "region",
			},
			ExpectedOutput: &AccountCreateResponse{
				AccountLocator:    "locator",
				AccountLocatorUrl: "locator_url",
				AccountName:       "acc_name",
				Url:               `https://org_name-acc_name.snowflakecomputing.com`,
				OrganizationName:  "ORG_NAME",
				Edition:           EditionBusinessCritical,
				RegionGroup:       "region_group",
				Cloud:             "cloud",
				Region:            "region",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			input := tc.RawInput
			if tc.Input != (AccountCreateResponse{}) {
				bytes, err := json.Marshal(tc.Input)
				if err != nil {
					assert.Fail(t, err.Error())
				}
				input = string(bytes)
			}

			createResponse, err := ToAccountCreateResponse(input)

			if tc.Error != "" {
				assert.EqualError(t, err, tc.Error)
				assert.Nil(t, createResponse)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.ExpectedOutput, createResponse)
			}
		})
	}
}

func TestToAccountEdition(t *testing.T) {
	type test struct {
		input string
		want  AccountEdition
	}

	valid := []test{
		// case insensitive.
		{input: "standard", want: EditionStandard},

		// Supported Values
		{input: "STANDARD", want: EditionStandard},
		{input: "ENTERPRISE", want: EditionEnterprise},
		{input: "BUSINESS_CRITICAL", want: EditionBusinessCritical},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
		{input: "businesscritical"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToAccountEdition(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToAccountEdition(tc.input)
			require.Error(t, err)
		})
	}
}
