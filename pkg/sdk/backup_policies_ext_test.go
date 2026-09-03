package sdk

func init() {
	renameTarget := randomSchemaObjectIdentifier()

	backupPoliciesTests.Create.
		withDefaultOpts(func() *CreateBackupPolicyOptions {
			return &CreateBackupPolicyOptions{
				name:            backupPoliciesTestIdSchemaObjectIdentifier,
				ExpireAfterDays: new(7),
			}
		}).
		withExpectedSqlf(
			case_BackupPolicies_sql_Create_basic,
			"CREATE BACKUP POLICY %s EXPIRE_AFTER_DAYS = 7", backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Create_all,
			func(opts *CreateBackupPolicyOptions) {
				opts.IfNotExists = new(true)
				opts.Tag = []TagAssociation{
					{
						Name:  NewAccountObjectIdentifier("tag1"),
						Value: "value1",
					},
					{
						Name:  NewAccountObjectIdentifier("tag2"),
						Value: "value2",
					},
				}
				opts.WithRetentionLock = new(true)
				opts.Schedule = new("USING CRON 0 2 * * * UTC")
				opts.ExpireAfterDays = new(3653)
				opts.Comment = new("some comment")
			},
			`CREATE BACKUP POLICY IF NOT EXISTS %s TAG ("tag1" = 'value1', "tag2" = 'value2') WITH RETENTION LOCK SCHEDULE = 'USING CRON 0 2 * * * UTC' EXPIRE_AFTER_DAYS = 3653 COMMENT = 'some comment'`, backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_scheduleOnly",
			func(opts *CreateBackupPolicyOptions) {
				opts.ExpireAfterDays = nil
				opts.Schedule = new("60 MINUTE")
			},
			"CREATE BACKUP POLICY %s SCHEDULE = '60 MINUTE'", backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateBackupPolicyOptions) {
				opts.OrReplace = new(true)
			},
			"CREATE OR REPLACE BACKUP POLICY %s EXPIRE_AFTER_DAYS = 7", backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_withRetentionLock",
			func(opts *CreateBackupPolicyOptions) {
				opts.WithRetentionLock = new(true)
			},
			"CREATE BACKUP POLICY %s WITH RETENTION LOCK EXPIRE_AFTER_DAYS = 7", backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	backupPoliciesTests.Alter.
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Alter_RenameTo,
			func(opts *AlterBackupPolicyOptions) { opts.RenameTo = &renameTarget },
			"ALTER BACKUP POLICY %s RENAME TO %s", backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Alter_Set,
			func(opts *AlterBackupPolicyOptions) {
				opts.Set = &BackupPolicySet{
					Schedule:        new("2 HOUR"),
					ExpireAfterDays: new(30),
					Comment:         new("some comment"),
				}
			},
			"ALTER BACKUP POLICY %s SET SCHEDULE = '2 HOUR' EXPIRE_AFTER_DAYS = 30 COMMENT = 'some comment'", backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Alter_SetTags,
			func(opts *AlterBackupPolicyOptions) {
				opts.SetTags = []TagAssociation{
					{
						Name:  NewAccountObjectIdentifier("tag1"),
						Value: "value1",
					},
					{
						Name:  NewAccountObjectIdentifier("tag2"),
						Value: "value2",
					},
				}
			},
			`ALTER BACKUP POLICY %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Alter_Unset,
			func(opts *AlterBackupPolicyOptions) {
				opts.Unset = &BackupPolicyUnset{
					Schedule:        new(true),
					ExpireAfterDays: new(true),
					Comment:         new(true),
				}
			},
			"ALTER BACKUP POLICY %s UNSET SCHEDULE, EXPIRE_AFTER_DAYS, COMMENT", backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Alter_UnsetTags,
			func(opts *AlterBackupPolicyOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("tag1"),
					NewAccountObjectIdentifier("tag2"),
				}
			},
			`ALTER BACKUP POLICY %s UNSET TAG "tag1", "tag2"`, backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	backupPoliciesTests.Drop.
		withExpectedSqlf(
			case_BackupPolicies_sql_Drop_basic,
			"DROP BACKUP POLICY %s", backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Drop_all,
			func(opts *DropBackupPolicyOptions) { opts.IfExists = new(true) },
			"DROP BACKUP POLICY IF EXISTS %s", backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	backupPoliciesTests.Show.
		withExpectedSqlf(
			case_BackupPolicies_sql_Show_basic,
			"SHOW BACKUP POLICIES",
		).
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Show_all,
			func(opts *ShowBackupPolicyOptions) {
				opts.Like = &Like{Pattern: String("like-pattern")}
				opts.In = &In{Account: new(true)}
				opts.StartsWith = String("abc")
				opts.Limit = &LimitFrom{Rows: Int(10)}
			},
			"SHOW BACKUP POLICIES LIKE 'like-pattern' IN ACCOUNT STARTS WITH 'abc' LIMIT 10",
		).
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Show_Like,
			func(opts *ShowBackupPolicyOptions) { opts.Like = &Like{Pattern: String("like-pattern")} },
			"SHOW BACKUP POLICIES LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Show_In,
			func(opts *ShowBackupPolicyOptions) { opts.In = &In{Account: new(true)} },
			"SHOW BACKUP POLICIES IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Show_StartsWith,
			func(opts *ShowBackupPolicyOptions) { opts.StartsWith = String("abc") },
			"SHOW BACKUP POLICIES STARTS WITH 'abc'",
		).
		withModifyAndExpectedSqlf(
			case_BackupPolicies_sql_Show_Limit,
			func(opts *ShowBackupPolicyOptions) { opts.Limit = &LimitFrom{Rows: Int(10)} },
			"SHOW BACKUP POLICIES LIMIT 10",
		)

	backupPoliciesTests.Describe.
		withExpectedSqlf(
			case_BackupPolicies_sql_Describe_basic,
			"DESCRIBE BACKUP POLICY %s", backupPoliciesTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)
}
