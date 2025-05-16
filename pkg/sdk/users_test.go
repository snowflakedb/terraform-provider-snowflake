package sdk

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/stretchr/testify/require"
)

func TestUserCreate(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &CreateUserOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("with only required attributes", func(t *testing.T) {
		opts := &CreateUserOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE USER %s`, id.FullyQualifiedName())
	})

	t.Run("empty secondary roles", func(t *testing.T) {
		opts := &CreateUserOptions{
			name: id,
			ObjectProperties: &UserObjectProperties{
				DefaultSecondaryRoles: &SecondaryRoles{None: Bool(true)},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE USER %s DEFAULT_SECONDARY_ROLES = ()`, id.FullyQualifiedName())
	})

	t.Run("with empty secondary roles", func(t *testing.T) {
		opts := &CreateUserOptions{
			name: id,
			ObjectProperties: &UserObjectProperties{
				DefaultSecondaryRoles: &SecondaryRoles{},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("SecondaryRoles", "All", "None"))
	})

	t.Run("with both options in secondary roles", func(t *testing.T) {
		opts := &CreateUserOptions{
			name: id,
			ObjectProperties: &UserObjectProperties{
				DefaultSecondaryRoles: &SecondaryRoles{
					All:  Bool(true),
					None: Bool(true),
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("SecondaryRoles", "All", "None"))
	})

	t.Run("with type", func(t *testing.T) {
		opts := &CreateUserOptions{
			name: id,
			ObjectProperties: &UserObjectProperties{
				Type: Pointer(UserTypeLegacyService),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE USER %s TYPE = LEGACY_SERVICE`, id.FullyQualifiedName())
	})

	t.Run("with complete options - no type", func(t *testing.T) {
		tagId := randomSchemaObjectIdentifier()
		tags := []TagAssociation{
			{
				Name:  tagId,
				Value: "v1",
			},
		}
		password := random.Password()
		loginName := random.String()
		defaultRoleId := randomAccountObjectIdentifier()
		defaultWarehouseId := randomAccountObjectIdentifier()
		var defaultNamespaceId ObjectIdentifier = randomDatabaseObjectIdentifier()

		opts := &CreateUserOptions{
			OrReplace:   Bool(true),
			name:        id,
			IfNotExists: Bool(true),
			ObjectProperties: &UserObjectProperties{
				Password:              &password,
				LoginName:             &loginName,
				DefaultRole:           Pointer(defaultRoleId),
				DefaultNamespace:      Pointer(defaultNamespaceId),
				DefaultWarehouse:      Pointer(defaultWarehouseId),
				DefaultSecondaryRoles: &SecondaryRoles{All: Bool(true)},
			},
			ObjectParameters: &UserObjectParameters{
				EnableUnredactedQuerySyntaxError: Bool(true),
			},
			SessionParameters: &SessionParameters{
				Autocommit: Bool(true),
			},
			With: Bool(true),
			Tags: tags,
		}

		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE USER IF NOT EXISTS %s PASSWORD = '%s' LOGIN_NAME = '%s' DEFAULT_WAREHOUSE = %s DEFAULT_NAMESPACE = %s DEFAULT_ROLE = %s DEFAULT_SECONDARY_ROLES = ('ALL') ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true AUTOCOMMIT = true WITH TAG (%s = 'v1')`, id.FullyQualifiedName(), password, loginName, defaultWarehouseId.FullyQualifiedName(), defaultNamespaceId.FullyQualifiedName(), defaultRoleId.FullyQualifiedName(), tagId.FullyQualifiedName())
	})
}

func TestUserAlter(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &AlterUserOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterUserOptions", "NewName", "ResetPassword", "AbortAllQueries", "AddDelegatedAuthorization", "RemoveDelegatedAuthorization", "Set", "Unset", "SetTag", "UnsetTag"))
	})

	t.Run("validation: no set", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Set:  &UserSet{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("UserSet", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ObjectProperties", "ObjectParameters", "SessionParameters"))
	})

	t.Run("validation: set more than one policy", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				AuthenticationPolicy: Pointer(randomSchemaObjectIdentifier()),
				PasswordPolicy:       Pointer(randomSchemaObjectIdentifier()),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("UserSet", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy"))
	})

	t.Run("validation: set policy with user parameters and properties", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				AuthenticationPolicy: Pointer(randomSchemaObjectIdentifier()),
				SessionParameters:    &SessionParameters{AbortDetachedQuery: Bool(true)},
				ObjectParameters:     &UserObjectParameters{EnableUnredactedQuerySyntaxError: Bool(true)},
				ObjectProperties:     &UserAlterObjectProperties{DisableMfa: Bool(true)},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errors.New("policies cannot be set with user properties or parameters at the same time"))
	})

	t.Run("two sets", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				SessionParameters: &SessionParameters{AbortDetachedQuery: Bool(true)},
				ObjectParameters:  &UserObjectParameters{EnableUnredactedQuerySyntaxError: Bool(true)},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true ABORT_DETACHED_QUERY = true", id.FullyQualifiedName())
	})

	t.Run("set type", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				ObjectProperties: &UserAlterObjectProperties{UserObjectProperties: UserObjectProperties{Type: Pointer(UserTypeLegacyService)}},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET TYPE = LEGACY_SERVICE", id.FullyQualifiedName())
	})

	t.Run("validation: no unset", func(t *testing.T) {
		opts := &AlterUserOptions{
			name:  id,
			Unset: &UserUnset{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("UserUnset", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ObjectProperties", "ObjectParameters", "SessionParameters"))
	})

	t.Run("validation: unset property with policy", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Unset: &UserUnset{
				PasswordPolicy:   Bool(true),
				ObjectParameters: &UserObjectParametersUnset{EnableUnredactedQuerySyntaxError: Bool(true)},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errors.New("policies cannot be unset with user properties or parameters at the same time"))
	})

	t.Run("alter: unset type", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Unset: &UserUnset{
				ObjectProperties: &UserObjectPropertiesUnset{Type: Bool(true)},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s UNSET TYPE", id.FullyQualifiedName())
	})

	t.Run("validation: unset two policies", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Unset: &UserUnset{
				PasswordPolicy:       Bool(true),
				AuthenticationPolicy: Bool(true),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("UserUnset", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy"))
	})

	t.Run("two compatible unsets", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Unset: &UserUnset{
				ObjectParameters:  &UserObjectParametersUnset{EnableUnredactedQuerySyntaxError: Bool(true)},
				SessionParameters: &SessionParametersUnset{BinaryOutputFormat: Bool(true)},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s UNSET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR, BINARY_OUTPUT_FORMAT", id.FullyQualifiedName())
	})

	t.Run("with setting a policy", func(t *testing.T) {
		passwordPolicy := randomSchemaObjectIdentifier()
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				PasswordPolicy: &passwordPolicy,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET PASSWORD POLICY %s", id.FullyQualifiedName(), passwordPolicy.FullyQualifiedName())
	})

	t.Run("with setting an authentication policy", func(t *testing.T) {
		authenticationPolicy := randomSchemaObjectIdentifier()
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				AuthenticationPolicy: &authenticationPolicy,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET AUTHENTICATION POLICY %s", id.FullyQualifiedName(), authenticationPolicy.FullyQualifiedName())
	})

	t.Run("with setting tags", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		tags := []TagAssociation{
			{
				Name:  tagId1,
				Value: "v1",
			},
			{
				Name:  tagId2,
				Value: "v2",
			},
		}
		opts := &AlterUserOptions{
			name:   id,
			SetTag: tags,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER %s SET TAG %s = 'v1', %s = 'v2'`, id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("with setting properties and parameters", func(t *testing.T) {
		password := random.Password()
		objectProperties := UserAlterObjectProperties{
			UserObjectProperties: UserObjectProperties{
				Password: &password,
				DefaultSecondaryRoles: &SecondaryRoles{
					All: Bool(true),
				},
			},
		}
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				ObjectProperties: &objectProperties,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET PASSWORD = '%s' DEFAULT_SECONDARY_ROLES = ('ALL')", id.FullyQualifiedName(), password)

		objectParameters := UserObjectParameters{
			EnableUnredactedQuerySyntaxError: Bool(true),
		}

		opts = &AlterUserOptions{
			name: id,
			Set: &UserSet{
				ObjectParameters: &objectParameters,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR = true", id.FullyQualifiedName())

		sessionParameters := SessionParameters{
			Autocommit: Bool(true),
		}
		opts = &AlterUserOptions{
			name: id,
			Set: &UserSet{
				SessionParameters: &sessionParameters,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET AUTOCOMMIT = true", id.FullyQualifiedName())
	})

	t.Run("alter: set object properties", func(t *testing.T) {
		objectProperties := UserAlterObjectProperties{
			UserObjectProperties: UserObjectProperties{
				FirstName: String("name"),
			},
			DisableMfa: Bool(true),
		}
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				ObjectProperties: &objectProperties,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET FIRST_NAME = '%s' DISABLE_MFA = true", id.FullyQualifiedName(), "name")
	})

	t.Run("alter: set disable mfa only", func(t *testing.T) {
		objectProperties := UserAlterObjectProperties{
			DisableMfa: Bool(true),
		}
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				ObjectProperties: &objectProperties,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s SET DISABLE_MFA = true", id.FullyQualifiedName())
	})

	t.Run("reset password", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &AlterUserOptions{
			name:          id,
			ResetPassword: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s RESET PASSWORD", id.FullyQualifiedName())
	})

	t.Run("abort all queries", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		opts := &AlterUserOptions{
			name:            id,
			AbortAllQueries: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s ABORT ALL QUERIES", id.FullyQualifiedName())
	})

	t.Run("rename", func(t *testing.T) {
		newID := randomAccountObjectIdentifier()
		opts := &AlterUserOptions{
			name:    id,
			NewName: newID,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s RENAME TO %s", id.FullyQualifiedName(), newID.FullyQualifiedName())
	})

	t.Run("with adding delegated authorization of role", func(t *testing.T) {
		role := "ROLE1"
		integration := "INTEGRATION1"
		opts := &AlterUserOptions{
			name: id,
			AddDelegatedAuthorization: &AddDelegatedAuthorization{
				Role:        role,
				Integration: integration,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s ADD DELEGATED AUTHORIZATION OF ROLE %s TO SECURITY INTEGRATION %s", id.FullyQualifiedName(), role, integration)
	})

	t.Run("with unsetting tags", func(t *testing.T) {
		tag1 := randomSchemaObjectIdentifier()
		tag2 := randomSchemaObjectIdentifier()
		opts := &AlterUserOptions{
			name:     id,
			UnsetTag: []ObjectIdentifier{tag1, tag2},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s UNSET TAG %s, %s", id.FullyQualifiedName(), tag1.FullyQualifiedName(), tag2.FullyQualifiedName())
	})

	t.Run("with unsetting properties", func(t *testing.T) {
		objectProperties := UserObjectPropertiesUnset{
			Password: Bool(true),
		}
		opts := &AlterUserOptions{
			name: id,
			Unset: &UserUnset{
				ObjectProperties: &objectProperties,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s UNSET PASSWORD", id.FullyQualifiedName())
	})

	t.Run("with unsetting a policy", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Unset: &UserUnset{
				SessionPolicy: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s UNSET SESSION POLICY", id.FullyQualifiedName())
	})

	t.Run("with unsetting an authentication policy", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Unset: &UserUnset{
				AuthenticationPolicy: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s UNSET AUTHENTICATION POLICY", id.FullyQualifiedName())
	})

	t.Run("with removing delegated authorization of role", func(t *testing.T) {
		role := "ROLE1"
		integration := "INTEGRATION1"
		opts := &AlterUserOptions{
			name: id,
			RemoveDelegatedAuthorization: &RemoveDelegatedAuthorization{
				Role:        &role,
				Integration: integration,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER USER %s REMOVE DELEGATED AUTHORIZATION OF ROLE %s FROM SECURITY INTEGRATION %s", id.FullyQualifiedName(), role, integration)
	})

	t.Run("empty secondary roles", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				ObjectProperties: &UserAlterObjectProperties{
					UserObjectProperties: UserObjectProperties{
						DefaultSecondaryRoles: &SecondaryRoles{},
					},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("SecondaryRoles", "All", "None"))
	})

	t.Run("with empty secondary roles", func(t *testing.T) {
		opts := &AlterUserOptions{
			name: id,
			Set: &UserSet{
				ObjectProperties: &UserAlterObjectProperties{
					UserObjectProperties: UserObjectProperties{
						DefaultSecondaryRoles: &SecondaryRoles{
							All:  Bool(true),
							None: Bool(true),
						},
					},
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("SecondaryRoles", "All", "None"))
	})
}

func TestUserDrop(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &DropUserOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &DropUserOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP USER %s", id.FullyQualifiedName())
	})
}

func TestUserShow(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &ShowUserOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW USERS")
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowUserOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW USERS LIKE '%s'", id.Name())
	})

	t.Run("with like and from", func(t *testing.T) {
		fromPatern := random.String()
		opts := &ShowUserOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			From: &fromPatern,
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW USERS LIKE '%s' FROM '%s'", id.Name(), fromPatern)
	})

	t.Run("with like and limit", func(t *testing.T) {
		limit := 5
		opts := &ShowUserOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
			Limit: &limit,
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW USERS LIKE '%s' LIMIT %v", id.Name(), limit)
	})

	t.Run("with starts with and from", func(t *testing.T) {
		fromPattern := random.String()
		startsWithPattern := random.String()

		opts := &ShowUserOptions{
			StartsWith: &startsWithPattern,
			From:       &fromPattern,
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW USERS STARTS WITH '%s' FROM '%s'", startsWithPattern, fromPattern)
	})
}

func TestUserDescribe(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &describeUserOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &describeUserOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE USER %s", id.FullyQualifiedName())
	})
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
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToSecondaryRolesOption(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		tc := tc
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
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			got := GetSecondaryRolesOptionFrom(tc.input)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range valid {
		tc := tc
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
