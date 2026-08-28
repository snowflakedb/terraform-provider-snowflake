package sdk

func init() {
	id := userProgrammaticAccessTokensTestIdAccountObjectIdentifier
	userId := randomAccountObjectIdentifier()
	roleId := randomAccountObjectIdentifier()
	renameTarget := randomAccountObjectIdentifier()

	userProgrammaticAccessTokensTests.Add.
		withDefaultOpts(func() *AddUserProgrammaticAccessTokenOptions {
			return &AddUserProgrammaticAccessTokenOptions{
				name:     id,
				UserName: userId,
			}
		}).
		withAdditionalValidationCase(
			"validation_Add_UserName_ValidIdentifier",
			func(opts *AddUserProgrammaticAccessTokenOptions) { opts.UserName = emptyAccountObjectIdentifier },
			errInvalidIdentifier("AddUserProgrammaticAccessTokenOptions", "UserName"),
		).
		withAdditionalValidationCase(
			"validation_Add_DaysToExpiry_greaterOrEqual1",
			func(opts *AddUserProgrammaticAccessTokenOptions) { opts.DaysToExpiry = new(0) },
			errIntValue("AddUserProgrammaticAccessTokenOptions", "DaysToExpiry", IntErrGreaterOrEqual, 1),
		).
		withAdditionalValidationCase(
			"validation_Add_MinsToBypassNetworkPolicyRequirement_greaterOrEqual1",
			func(opts *AddUserProgrammaticAccessTokenOptions) {
				opts.MinsToBypassNetworkPolicyRequirement = new(0)
			},
			errIntValue("AddUserProgrammaticAccessTokenOptions", "MinsToBypassNetworkPolicyRequirement", IntErrGreaterOrEqual, 1),
		).
		withExpectedSqlf(
			case_UserProgrammaticAccessTokens_sql_Add_basic,
			"ALTER USER %s ADD PROGRAMMATIC ACCESS TOKEN %s",
			userId.FullyQualifiedName(), id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Add_all",
			func(opts *AddUserProgrammaticAccessTokenOptions) {
				opts.RoleRestriction = &roleId
				opts.DaysToExpiry = new(30)
				opts.MinsToBypassNetworkPolicyRequirement = new(10)
				opts.Comment = new("test comment")
			},
			"ALTER USER %s ADD PROGRAMMATIC ACCESS TOKEN %s ROLE_RESTRICTION = %s DAYS_TO_EXPIRY = 30 MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT = 10 COMMENT = 'test comment'",
			userId.FullyQualifiedName(), id.FullyQualifiedName(), roleId.FullyQualifiedName(),
		)

	userProgrammaticAccessTokensTests.Modify.
		withDefaultOpts(func() *ModifyUserProgrammaticAccessTokenOptions {
			return &ModifyUserProgrammaticAccessTokenOptions{
				name:     id,
				UserName: userId,
			}
		}).
		withAdditionalValidationCase(
			"validation_Modify_UserName_ValidIdentifier",
			func(opts *ModifyUserProgrammaticAccessTokenOptions) {
				opts.UserName = emptyAccountObjectIdentifier
				opts.RenameTo = &renameTarget
			},
			errInvalidIdentifier("ModifyUserProgrammaticAccessTokenOptions", "UserName"),
		).
		withAdditionalValidationCase(
			"validation_Modify_Set_MinsToBypassNetworkPolicyRequirement_greaterOrEqual1",
			func(opts *ModifyUserProgrammaticAccessTokenOptions) {
				opts.Set = &ModifyProgrammaticAccessTokenSet{MinsToBypassNetworkPolicyRequirement: new(0)}
			},
			errIntValue("ModifyUserProgrammaticAccessTokenOptions", "Set.MinsToBypassNetworkPolicyRequirement", IntErrGreaterOrEqual, 1),
		).
		withModifyAndExpectedSqlf(
			case_UserProgrammaticAccessTokens_sql_Modify_basic,
			func(opts *ModifyUserProgrammaticAccessTokenOptions) { opts.RenameTo = &renameTarget },
			"ALTER USER %s MODIFY PROGRAMMATIC ACCESS TOKEN %s RENAME TO %s",
			userId.FullyQualifiedName(), id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Modify_Set",
			func(opts *ModifyUserProgrammaticAccessTokenOptions) {
				opts.Set = &ModifyProgrammaticAccessTokenSet{
					Disabled:                             new(true),
					MinsToBypassNetworkPolicyRequirement: new(10),
					Comment:                              new("new comment"),
				}
			},
			"ALTER USER %s MODIFY PROGRAMMATIC ACCESS TOKEN %s SET DISABLED = true MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT = 10 COMMENT = 'new comment'",
			userId.FullyQualifiedName(), id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Modify_Unset",
			func(opts *ModifyUserProgrammaticAccessTokenOptions) {
				opts.Unset = &ModifyProgrammaticAccessTokenUnset{
					Disabled:                             new(true),
					MinsToBypassNetworkPolicyRequirement: new(true),
					Comment:                              new(true),
				}
			},
			"ALTER USER %s MODIFY PROGRAMMATIC ACCESS TOKEN %s UNSET DISABLED, MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT, COMMENT",
			userId.FullyQualifiedName(), id.FullyQualifiedName(),
		)

	userProgrammaticAccessTokensTests.Rotate.
		withDefaultOpts(func() *RotateUserProgrammaticAccessTokenOptions {
			return &RotateUserProgrammaticAccessTokenOptions{
				name:     id,
				UserName: userId,
			}
		}).
		withAdditionalValidationCase(
			"validation_Rotate_UserName_ValidIdentifier",
			func(opts *RotateUserProgrammaticAccessTokenOptions) { opts.UserName = emptyAccountObjectIdentifier },
			errInvalidIdentifier("RotateUserProgrammaticAccessTokenOptions", "UserName"),
		).
		withAdditionalValidationCase(
			"validation_Rotate_ExpireRotatedTokenAfterHours_greaterOrEqual0",
			func(opts *RotateUserProgrammaticAccessTokenOptions) { opts.ExpireRotatedTokenAfterHours = new(-1) },
			errIntValue("RotateUserProgrammaticAccessTokenOptions", "ExpireRotatedTokenAfterHours", IntErrGreaterOrEqual, 0),
		).
		withExpectedSqlf(
			case_UserProgrammaticAccessTokens_sql_Rotate_basic,
			"ALTER USER %s ROTATE PROGRAMMATIC ACCESS TOKEN %s",
			userId.FullyQualifiedName(), id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Rotate_all",
			func(opts *RotateUserProgrammaticAccessTokenOptions) { opts.ExpireRotatedTokenAfterHours = new(1) },
			"ALTER USER %s ROTATE PROGRAMMATIC ACCESS TOKEN %s EXPIRE_ROTATED_TOKEN_AFTER_HOURS = 1",
			userId.FullyQualifiedName(), id.FullyQualifiedName(),
		)

	userProgrammaticAccessTokensTests.Remove.
		withDefaultOpts(func() *RemoveUserProgrammaticAccessTokenOptions {
			return &RemoveUserProgrammaticAccessTokenOptions{
				name:     id,
				UserName: userId,
			}
		}).
		withAdditionalValidationCase(
			"validation_Remove_UserName_ValidIdentifier",
			func(opts *RemoveUserProgrammaticAccessTokenOptions) { opts.UserName = emptyAccountObjectIdentifier },
			errInvalidIdentifier("RemoveUserProgrammaticAccessTokenOptions", "UserName"),
		).
		withExpectedSqlf(
			case_UserProgrammaticAccessTokens_sql_Remove_basic,
			"ALTER USER %s REMOVE PROGRAMMATIC ACCESS TOKEN %s",
			userId.FullyQualifiedName(), id.FullyQualifiedName(),
		)

	userProgrammaticAccessTokensTests.Show.
		withExpectedSql(case_UserProgrammaticAccessTokens_sql_Show_basic, "SHOW USER PROGRAMMATIC ACCESS TOKENS").
		withModifyAndExpectedSqlf(
			case_UserProgrammaticAccessTokens_sql_Show_all,
			func(opts *ShowUserProgrammaticAccessTokenOptions) { opts.UserName = &userId },
			"SHOW USER PROGRAMMATIC ACCESS TOKENS FOR USER %s", userId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_UserProgrammaticAccessTokens_sql_Show_UserName,
			func(opts *ShowUserProgrammaticAccessTokenOptions) { opts.UserName = &userId },
			"SHOW USER PROGRAMMATIC ACCESS TOKENS FOR USER %s", userId.FullyQualifiedName(),
		)
}
