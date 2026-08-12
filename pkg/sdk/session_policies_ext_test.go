package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	id := sessionPoliciesTestIdSchemaObjectIdentifier

	sessionPoliciesTests.Create.
		withModifyAndExpectedErr(
			case_SessionPolicies_validation_Create_AllowedSecondaryRoles_ValidateValue,
			func(opts *CreateSessionPolicyOptions) {
				opts.AllowedSecondaryRoles = &SessionPolicySecondaryRoles{All: new(false)}
			},
			errInvalidValue("SessionPolicySecondaryRoles", "All", "false"),
		).
		withAdditionalValidationCase(
			"validation_Create_AllowedSecondaryRoles_None_ValidateValue",
			func(opts *CreateSessionPolicyOptions) {
				opts.AllowedSecondaryRoles = &SessionPolicySecondaryRoles{None: new(false)}
			},
			errInvalidValue("SessionPolicySecondaryRoles", "None", "false"),
		).
		withModifyAndExpectedErr(
			case_SessionPolicies_validation_Create_BlockedSecondaryRoles_ValidateValue,
			func(opts *CreateSessionPolicyOptions) {
				opts.BlockedSecondaryRoles = &SessionPolicySecondaryRoles{All: new(false)}
			},
			errInvalidValue("SessionPolicySecondaryRoles", "All", "false"),
		).
		withAdditionalValidationCase(
			"validation_Create_BlockedSecondaryRoles_None_ValidateValue",
			func(opts *CreateSessionPolicyOptions) {
				opts.BlockedSecondaryRoles = &SessionPolicySecondaryRoles{None: new(false)}
			},
			errInvalidValue("SessionPolicySecondaryRoles", "None", "false"),
		).
		withExpectedSqlf(
			case_SessionPolicies_sql_Create_basic,
			`CREATE SESSION POLICY %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Create_all,
			func(opts *CreateSessionPolicyOptions) {
				opts.OrReplace = new(true)
				opts.SessionIdleTimeoutMins = new(5)
				opts.SessionUiIdleTimeoutMins = new(34)
				opts.AllowedSecondaryRoles = &SessionPolicySecondaryRoles{Roles: []AccountObjectIdentifier{NewAccountObjectIdentifier("ROLE1"), NewAccountObjectIdentifier("ROLE2")}}
				opts.BlockedSecondaryRoles = &SessionPolicySecondaryRoles{All: new(true)}
				opts.Comment = new("some comment")
			},
			`CREATE OR REPLACE SESSION POLICY %s SESSION_IDLE_TIMEOUT_MINS = 5 SESSION_UI_IDLE_TIMEOUT_MINS = 34 ALLOWED_SECONDARY_ROLES = ("ROLE1", "ROLE2") BLOCKED_SECONDARY_ROLES = ('ALL') COMMENT = 'some comment'`,
			id.FullyQualifiedName(),
		)

	newId := randomSchemaObjectIdentifier()
	tagId1 := NewAccountObjectIdentifier("tag1")
	tagId2 := NewAccountObjectIdentifier("tag2")

	sessionPoliciesTests.Alter.
		withModifyAndExpectedErr(
			case_SessionPolicies_validation_Alter_Set_AllowedSecondaryRoles_ValidateValue,
			func(opts *AlterSessionPolicyOptions) {
				opts.Set = &SessionPolicySet{AllowedSecondaryRoles: &SessionPolicySecondaryRoles{All: new(false)}}
			},
			errInvalidValue("SessionPolicySecondaryRoles", "All", "false"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_AllowedSecondaryRoles_None_ValidateValue",
			func(opts *AlterSessionPolicyOptions) {
				opts.Set = &SessionPolicySet{AllowedSecondaryRoles: &SessionPolicySecondaryRoles{None: new(false)}}
			},
			errInvalidValue("SessionPolicySecondaryRoles", "None", "false"),
		).
		withModifyAndExpectedErr(
			case_SessionPolicies_validation_Alter_Set_BlockedSecondaryRoles_ValidateValue,
			func(opts *AlterSessionPolicyOptions) {
				opts.Set = &SessionPolicySet{BlockedSecondaryRoles: &SessionPolicySecondaryRoles{All: new(false)}}
			},
			errInvalidValue("SessionPolicySecondaryRoles", "All", "false"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_BlockedSecondaryRoles_None_ValidateValue",
			func(opts *AlterSessionPolicyOptions) {
				opts.Set = &SessionPolicySet{BlockedSecondaryRoles: &SessionPolicySecondaryRoles{None: new(false)}}
			},
			errInvalidValue("SessionPolicySecondaryRoles", "None", "false"),
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Alter_Set,
			func(opts *AlterSessionPolicyOptions) {
				opts.Set = &SessionPolicySet{Comment: new("some comment")}
			},
			`ALTER SESSION POLICY %s SET COMMENT = 'some comment'`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_secondaryRoles",
			func(opts *AlterSessionPolicyOptions) {
				opts.Set = &SessionPolicySet{
					AllowedSecondaryRoles: &SessionPolicySecondaryRoles{None: new(true)},
					BlockedSecondaryRoles: &SessionPolicySecondaryRoles{Roles: []AccountObjectIdentifier{NewAccountObjectIdentifier("ROLE1"), NewAccountObjectIdentifier("ROLE2")}},
				}
			},
			`ALTER SESSION POLICY %s SET ALLOWED_SECONDARY_ROLES = () BLOCKED_SECONDARY_ROLES = ("ROLE1", "ROLE2")`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Alter_Unset,
			func(opts *AlterSessionPolicyOptions) {
				opts.Unset = &SessionPolicyUnset{
					SessionIdleTimeoutMins:   new(true),
					SessionUiIdleTimeoutMins: new(true),
					Comment:                  new(true),
				}
			},
			`ALTER SESSION POLICY %s UNSET SESSION_IDLE_TIMEOUT_MINS, SESSION_UI_IDLE_TIMEOUT_MINS, COMMENT`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Alter_RenameTo,
			func(opts *AlterSessionPolicyOptions) { opts.RenameTo = new(newId) },
			`ALTER SESSION POLICY %s RENAME TO %s`, id.FullyQualifiedName(), newId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Alter_SetTags,
			func(opts *AlterSessionPolicyOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId1, Value: "value1"},
					{Name: tagId2, Value: "value2"},
				}
			},
			`ALTER SESSION POLICY %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Alter_UnsetTags,
			func(opts *AlterSessionPolicyOptions) { opts.UnsetTags = []ObjectIdentifier{tagId1, tagId2} },
			`ALTER SESSION POLICY %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName(),
		)

	sessionPoliciesTests.Drop.
		withExpectedSqlf(
			case_SessionPolicies_sql_Drop_basic,
			`DROP SESSION POLICY %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Drop_all,
			func(opts *DropSessionPolicyOptions) { opts.IfExists = new(true) },
			`DROP SESSION POLICY IF EXISTS %s`, id.FullyQualifiedName(),
		)

	showId := randomSchemaObjectIdentifier()

	sessionPoliciesTests.Show.
		withExpectedSqlf(
			case_SessionPolicies_sql_Show_basic,
			`SHOW SESSION POLICIES`,
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Show_all,
			func(opts *ShowSessionPolicyOptions) {
				opts.Like = &Like{Pattern: new("like-pattern")}
				opts.In = &ExtendedIn{In: In{Schema: showId.SchemaId()}}
				opts.StartsWith = new("starts-with-pattern")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("limit-from")}
			},
			`SHOW SESSION POLICIES LIKE 'like-pattern' IN SCHEMA %s STARTS WITH 'starts-with-pattern' LIMIT 10 FROM 'limit-from'`,
			showId.SchemaId().FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Show_Like,
			func(opts *ShowSessionPolicyOptions) { opts.Like = &Like{Pattern: new("like-pattern")} },
			`SHOW SESSION POLICIES LIKE 'like-pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Show_In,
			func(opts *ShowSessionPolicyOptions) { opts.In = &ExtendedIn{In: In{Schema: showId.SchemaId()}} },
			`SHOW SESSION POLICIES IN SCHEMA %s`, showId.SchemaId().FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Show_On,
			func(opts *ShowSessionPolicyOptions) { opts.On = &On{Account: new(true)} },
			`SHOW SESSION POLICIES ON ACCOUNT`,
		).
		withAdditionalSqlCasef(
			"sql_Show_On_user",
			func(opts *ShowSessionPolicyOptions) { opts.On = &On{User: NewAccountObjectIdentifier("user_name")} },
			`SHOW SESSION POLICIES ON USER "user_name"`,
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Show_StartsWith,
			func(opts *ShowSessionPolicyOptions) { opts.StartsWith = new("starts-with-pattern") },
			`SHOW SESSION POLICIES STARTS WITH 'starts-with-pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_SessionPolicies_sql_Show_Limit,
			func(opts *ShowSessionPolicyOptions) { opts.Limit = &LimitFrom{Rows: new(10), From: new("limit-from")} },
			`SHOW SESSION POLICIES LIMIT 10 FROM 'limit-from'`,
		)

	sessionPoliciesTests.Describe.
		withExpectedSqlf(
			case_SessionPolicies_sql_Describe_basic,
			`DESCRIBE SESSION POLICY %s`, id.FullyQualifiedName(),
		)
}

func Test_parseSessionPolicyTargetScopes(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected []SessionPolicyTargetScope
	}{
		{
			Name:     "empty options returns nil",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "missing target_scopes key returns nil",
			Input:    `{"other_field":"value"}`,
			Expected: nil,
		},
		{
			Name:     "empty target_scopes array returns empty slice",
			Input:    `{"target_scopes":[]}`,
			Expected: []SessionPolicyTargetScope{},
		},
		{
			Name:     "single scope",
			Input:    `{"target_scopes":["ACCOUNT"]}`,
			Expected: []SessionPolicyTargetScope{SessionPolicyTargetScopeAccount},
		},
		{
			Name:  "multiple scopes are sorted alphabetically",
			Input: `{"target_scopes":["SERVICE_USERS","ACCOUNT","PERSON_USERS"]}`,
			Expected: []SessionPolicyTargetScope{
				SessionPolicyTargetScopeAccount,
				SessionPolicyTargetScopePersonUsers,
				SessionPolicyTargetScopeServiceUsers,
			},
		},
		{
			Name:  "extra attributes are ignored",
			Input: `{"target_scopes":["PERSON_USERS","ACCOUNT"],"comment":"ignored","nested":{"a":1},"list":[1,2,3]}`,
			Expected: []SessionPolicyTargetScope{
				SessionPolicyTargetScopeAccount,
				SessionPolicyTargetScopePersonUsers,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := parseSessionPolicyTargetScopes(tc.Input)
			require.NoError(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}

	t.Run("invalid json returns error", func(t *testing.T) {
		result, err := parseSessionPolicyTargetScopes(`{"target_scopes":`)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
