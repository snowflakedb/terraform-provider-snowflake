package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddProgrammaticAccessToken(t *testing.T) {
	name := randomAccountObjectIdentifier()
	roleId := randomAccountObjectIdentifier()
	userId := randomAccountObjectIdentifier()

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AddUserProgrammaticAccessTokenOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid object name", func(t *testing.T) {
		opts := &AddUserProgrammaticAccessTokenOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid user name", func(t *testing.T) {
		opts := &AddUserProgrammaticAccessTokenOptions{
			name:     name,
			UserName: emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("AddUserProgrammaticAccessTokenOptions", "UserName"))
	})

	t.Run("validation: invalid days to expiry", func(t *testing.T) {
		opts := &AddUserProgrammaticAccessTokenOptions{
			name:         name,
			DaysToExpiry: Int(0),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("AddUserProgrammaticAccessTokenOptions", "DaysToExpiry", IntErrGreaterOrEqual, 1))
	})

	t.Run("validation: invalid mins to bypass network policy requirement", func(t *testing.T) {
		opts := &AddUserProgrammaticAccessTokenOptions{
			name:                                 name,
			MinsToBypassNetworkPolicyRequirement: Int(0),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("AddUserProgrammaticAccessTokenOptions", "MinsToBypassNetworkPolicyRequirement", IntErrGreaterOrEqual, 1))
	})

	t.Run("with only required attributes", func(t *testing.T) {
		opts := &AddUserProgrammaticAccessTokenOptions{
			UserName: userId,
			name:     name,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER %s ADD PROGRAMMATIC ACCESS TOKEN %s`, userId.FullyQualifiedName(), name.FullyQualifiedName())
	})

	t.Run("with all attributes", func(t *testing.T) {
		opts := &AddUserProgrammaticAccessTokenOptions{
			UserName:                             userId,
			name:                                 name,
			RoleRestriction:                      &roleId,
			DaysToExpiry:                         Int(30),
			MinsToBypassNetworkPolicyRequirement: Int(10),
			Comment:                              String("test comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER %s ADD PROGRAMMATIC ACCESS TOKEN %s ROLE_RESTRICTION = %s DAYS_TO_EXPIRY = 30 MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT = 10 COMMENT = 'test comment'`, userId.FullyQualifiedName(), name.FullyQualifiedName(), roleId.FullyQualifiedName())
	})
}

func TestModifyProgrammaticAccessToken(t *testing.T) {
	name := randomAccountObjectIdentifier()
	userId := randomAccountObjectIdentifier()

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ModifyUserProgrammaticAccessTokenOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid object name", func(t *testing.T) {
		opts := &ModifyUserProgrammaticAccessTokenOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid user name", func(t *testing.T) {
		opts := &ModifyUserProgrammaticAccessTokenOptions{
			name:     name,
			UserName: emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("ModifyUserProgrammaticAccessTokenOptions", "UserName"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.RenameTo] should be present: none set", func(t *testing.T) {
		opts := &ModifyUserProgrammaticAccessTokenOptions{
			name: name,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ModifyUserProgrammaticAccessTokenOptions", "Set", "Unset", "RenameTo"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.RenameTo] should be present: all set", func(t *testing.T) {
		newId := randomAccountObjectIdentifier()
		opts := &ModifyUserProgrammaticAccessTokenOptions{
			name:     name,
			RenameTo: &newId,
			Set: &ModifyProgrammaticAccessTokenSet{
				Disabled: Bool(true),
			},
			Unset: &ModifyProgrammaticAccessTokenUnset{
				Disabled: Bool(true),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ModifyUserProgrammaticAccessTokenOptions", "Set", "Unset", "RenameTo"))
	})

	t.Run("validation: invalid mins to bypass network policy requirement", func(t *testing.T) {
		opts := &ModifyUserProgrammaticAccessTokenOptions{
			name: name,
			Set: &ModifyProgrammaticAccessTokenSet{
				MinsToBypassNetworkPolicyRequirement: Int(0),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("ModifyUserProgrammaticAccessTokenOptions", "Set.MinsToBypassNetworkPolicyRequirement", IntErrGreaterOrEqual, 1))
	})

	t.Run("with rename to", func(t *testing.T) {
		newId := randomAccountObjectIdentifier()
		opts := &ModifyUserProgrammaticAccessTokenOptions{
			UserName: userId,
			name:     name,
			RenameTo: &newId,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER %s MODIFY PROGRAMMATIC ACCESS TOKEN %s RENAME TO %s`, userId.FullyQualifiedName(), name.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("with set: all attributes", func(t *testing.T) {
		opts := &ModifyUserProgrammaticAccessTokenOptions{
			UserName: userId,
			name:     name,
			Set: &ModifyProgrammaticAccessTokenSet{
				Disabled:                             Bool(true),
				MinsToBypassNetworkPolicyRequirement: Int(10),
				Comment:                              String("new comment"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER %s MODIFY PROGRAMMATIC ACCESS TOKEN %s SET DISABLED = true MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT = 10 COMMENT = 'new comment'`, userId.FullyQualifiedName(), name.FullyQualifiedName())
	})

	t.Run("with unset: all attributes", func(t *testing.T) {
		opts := &ModifyUserProgrammaticAccessTokenOptions{
			UserName: userId,
			name:     name,
			Unset: &ModifyProgrammaticAccessTokenUnset{
				Disabled:                             Bool(true),
				MinsToBypassNetworkPolicyRequirement: Bool(true),
				Comment:                              Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER %s MODIFY PROGRAMMATIC ACCESS TOKEN %s UNSET DISABLED, MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT, COMMENT`, userId.FullyQualifiedName(), name.FullyQualifiedName())
	})
}

func TestRotateProgrammaticAccessToken(t *testing.T) {
	name := randomAccountObjectIdentifier()
	userId := randomAccountObjectIdentifier()
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *RotateUserProgrammaticAccessTokenOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid object name", func(t *testing.T) {
		opts := &RotateUserProgrammaticAccessTokenOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid user name", func(t *testing.T) {
		opts := &RotateUserProgrammaticAccessTokenOptions{
			name:     name,
			UserName: emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("RotateUserProgrammaticAccessTokenOptions", "UserName"))
	})

	t.Run("validation: invalid expire rotated token after hours", func(t *testing.T) {
		opts := &RotateUserProgrammaticAccessTokenOptions{
			name:                         name,
			ExpireRotatedTokenAfterHours: Int(-1),
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("RotateUserProgrammaticAccessTokenOptions", "ExpireRotatedTokenAfterHours", IntErrGreaterOrEqual, 0))
	})

	t.Run("with required attributes", func(t *testing.T) {
		opts := &RotateUserProgrammaticAccessTokenOptions{
			UserName: userId,
			name:     name,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER %s ROTATE PROGRAMMATIC ACCESS TOKEN %s`, userId.FullyQualifiedName(), name.FullyQualifiedName())
	})

	t.Run("with all attributes", func(t *testing.T) {
		opts := &RotateUserProgrammaticAccessTokenOptions{
			UserName:                     userId,
			name:                         name,
			ExpireRotatedTokenAfterHours: Int(1),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER %s ROTATE PROGRAMMATIC ACCESS TOKEN %s EXPIRE_ROTATED_TOKEN_AFTER_HOURS = 1`, userId.FullyQualifiedName(), name.FullyQualifiedName())
	})
}

func TestRemoveProgrammaticAccessToken(t *testing.T) {
	name := randomAccountObjectIdentifier()
	userId := randomAccountObjectIdentifier()

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *RemoveUserProgrammaticAccessTokenOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid object name", func(t *testing.T) {
		opts := &RemoveUserProgrammaticAccessTokenOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid user name", func(t *testing.T) {
		opts := &RemoveUserProgrammaticAccessTokenOptions{
			name:     name,
			UserName: emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("RemoveUserProgrammaticAccessTokenOptions", "UserName"))
	})

	t.Run("with all attributes", func(t *testing.T) {
		opts := &RemoveUserProgrammaticAccessTokenOptions{
			UserName: userId,
			name:     name,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER USER %s REMOVE PROGRAMMATIC ACCESS TOKEN %s`, userId.FullyQualifiedName(), name.FullyQualifiedName())
	})
}

func TestShowProgrammaticAccessTokens(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("with basic attributes", func(t *testing.T) {
		opts := &ShowUserProgrammaticAccessTokenOptions{}
		assertOptsValidAndSQLEquals(t, opts, `SHOW USER PROGRAMMATIC ACCESS TOKENS`)
	})

	t.Run("with optional attributes", func(t *testing.T) {
		opts := &ShowUserProgrammaticAccessTokenOptions{
			UserName: &id,
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW USER PROGRAMMATIC ACCESS TOKENS FOR USER %s`, id.FullyQualifiedName())
	})
}

func Test_ProgrammaticAccessTokenStatus(t *testing.T) {
	type test struct {
		input string
		want  ProgrammaticAccessTokenStatus
	}

	valid := []test{
		// case insensitive.
		{input: "active", want: ProgrammaticAccessTokenStatusActive},

		// Supported Values
		{input: "ACTIVE", want: ProgrammaticAccessTokenStatusActive},
		{input: "EXPIRED", want: ProgrammaticAccessTokenStatusExpired},
		{input: "DISABLED", want: ProgrammaticAccessTokenStatusDisabled},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := toProgrammaticAccessTokenStatus(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := toProgrammaticAccessTokenStatus(tc.input)
			require.Error(t, err)
		})
	}
}
