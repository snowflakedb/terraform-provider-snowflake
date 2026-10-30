package sdk

func init() {
	table := randomSchemaObjectIdentifier()
	dynamicTable := randomSchemaObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifier()
	databaseId := randomAccountObjectIdentifier()
	policy := randomSchemaObjectIdentifier()
	renameTarget := randomSchemaObjectIdentifier()
	backupID := "790d1ee4-88b2-451f-9ccc-eacd1e93a134"

	backupSetsTests.Create.
		withDefaultOpts(func() *CreateBackupSetOptions {
			return &CreateBackupSetOptions{
				name:     backupSetsTestIdSchemaObjectIdentifier,
				ForTable: &table,
			}
		}).
		withExpectedSqlf(
			case_BackupSets_sql_Create_basic,
			"CREATE BACKUP SET %s FOR TABLE %s",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			table.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Create_all,
			func(opts *CreateBackupSetOptions) {
				opts.IfNotExists = new(true)
				opts.Tag = []TagAssociation{
					{Name: NewAccountObjectIdentifier("tag1"), Value: "value1"},
					{Name: NewAccountObjectIdentifier("tag2"), Value: "value2"},
				}
				opts.WithBackupPolicy = &policy
				opts.Comment = new("some comment")
			},
			`CREATE BACKUP SET IF NOT EXISTS %s FOR TABLE %s TAG ("tag1" = 'value1', "tag2" = 'value2') WITH BACKUP POLICY %s COMMENT = 'some comment'`,
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			table.FullyQualifiedName(),
			policy.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateBackupSetOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE BACKUP SET %s FOR TABLE %s",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			table.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_forDynamicTable",
			func(opts *CreateBackupSetOptions) {
				opts.ForTable = nil
				opts.ForDynamicTable = &dynamicTable
			},
			"CREATE BACKUP SET %s FOR DYNAMIC TABLE %s",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			dynamicTable.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_forSchema",
			func(opts *CreateBackupSetOptions) {
				opts.ForTable = nil
				opts.ForSchema = &schemaId
			},
			"CREATE BACKUP SET %s FOR SCHEMA %s",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			schemaId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_forDatabase",
			func(opts *CreateBackupSetOptions) {
				opts.ForTable = nil
				opts.ForDatabase = &databaseId
			},
			"CREATE BACKUP SET %s FOR DATABASE %s",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			databaseId.FullyQualifiedName(),
		)

	backupSetsTests.Alter.
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_AddBackup,
			func(opts *AlterBackupSetOptions) { opts.AddBackup = new(true) },
			"ALTER BACKUP SET %s ADD BACKUP",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_ApplyBackupPolicy,
			func(opts *AlterBackupSetOptions) {
				opts.ApplyBackupPolicy = &BackupSetApplyBackupPolicy{BackupPolicy: policy}
			},
			"ALTER BACKUP SET %s APPLY BACKUP POLICY %s",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			policy.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ApplyBackupPolicy_force",
			func(opts *AlterBackupSetOptions) {
				opts.ApplyBackupPolicy = &BackupSetApplyBackupPolicy{BackupPolicy: policy, Force: new(true)}
			},
			"ALTER BACKUP SET %s APPLY BACKUP POLICY %s FORCE",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			policy.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_SuspendBackupPolicy,
			func(opts *AlterBackupSetOptions) { opts.SuspendBackupPolicy = new(true) },
			"ALTER BACKUP SET %s SUSPEND BACKUP POLICY",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_SuspendBackupCreationPolicy,
			func(opts *AlterBackupSetOptions) { opts.SuspendBackupCreationPolicy = new(true) },
			"ALTER BACKUP SET %s SUSPEND BACKUP CREATION POLICY",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_SuspendBackupExpirationPolicy,
			func(opts *AlterBackupSetOptions) { opts.SuspendBackupExpirationPolicy = new(true) },
			"ALTER BACKUP SET %s SUSPEND BACKUP EXPIRATION POLICY",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_ResumeBackupPolicy,
			func(opts *AlterBackupSetOptions) { opts.ResumeBackupPolicy = new(true) },
			"ALTER BACKUP SET %s RESUME BACKUP POLICY",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_ResumeBackupCreationPolicy,
			func(opts *AlterBackupSetOptions) { opts.ResumeBackupCreationPolicy = new(true) },
			"ALTER BACKUP SET %s RESUME BACKUP CREATION POLICY",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_ResumeBackupExpirationPolicy,
			func(opts *AlterBackupSetOptions) { opts.ResumeBackupExpirationPolicy = new(true) },
			"ALTER BACKUP SET %s RESUME BACKUP EXPIRATION POLICY",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_DeleteBackupIdentifier,
			func(opts *AlterBackupSetOptions) { opts.DeleteBackupIdentifier = new(backupID) },
			"ALTER BACKUP SET %s DELETE BACKUP IDENTIFIER '%s'",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			backupID,
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_ModifyBackup,
			func(opts *AlterBackupSetOptions) {
				opts.ModifyBackup = &BackupSetModifyBackup{
					BackupIdentifier:               backupID,
					AddLegalHoldModifyBackupParams: &AddLegalHoldModifyBackupParams{},
				}
			},
			"ALTER BACKUP SET %s MODIFY BACKUP IDENTIFIER '%s' ADD LEGAL HOLD",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			backupID,
		).
		withAdditionalSqlCasef(
			"sql_Alter_ModifyBackup_removeLegalHold",
			func(opts *AlterBackupSetOptions) {
				opts.ModifyBackup = &BackupSetModifyBackup{
					BackupIdentifier:                  backupID,
					RemoveLegalHoldModifyBackupParams: &RemoveLegalHoldModifyBackupParams{},
				}
			},
			"ALTER BACKUP SET %s MODIFY BACKUP IDENTIFIER '%s' REMOVE LEGAL HOLD",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			backupID,
		).
		withAdditionalSqlCasef(
			"sql_Alter_ModifyBackup_setComment",
			func(opts *AlterBackupSetOptions) {
				opts.ModifyBackup = &BackupSetModifyBackup{
					BackupIdentifier:             backupID,
					SetCommentModifyBackupParams: &SetCommentModifyBackupParams{SetComment: "backup comment"},
				}
			},
			"ALTER BACKUP SET %s MODIFY BACKUP IDENTIFIER '%s' SET COMMENT = 'backup comment'",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			backupID,
		).
		withAdditionalSqlCasef(
			"sql_Alter_ModifyBackup_unsetComment",
			func(opts *AlterBackupSetOptions) {
				opts.ModifyBackup = &BackupSetModifyBackup{
					BackupIdentifier:               backupID,
					UnsetCommentModifyBackupParams: &UnsetCommentModifyBackupParams{},
				}
			},
			"ALTER BACKUP SET %s MODIFY BACKUP IDENTIFIER '%s' UNSET COMMENT",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			backupID,
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_RenameTo,
			func(opts *AlterBackupSetOptions) { opts.RenameTo = &renameTarget },
			"ALTER BACKUP SET %s RENAME TO %s",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
			renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_Set,
			func(opts *AlterBackupSetOptions) { opts.Set = &BackupSetSet{Comment: new("some comment")} },
			"ALTER BACKUP SET %s SET COMMENT = 'some comment'",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_Unset,
			func(opts *AlterBackupSetOptions) { opts.Unset = &BackupSetUnset{Comment: new(true)} },
			"ALTER BACKUP SET %s UNSET COMMENT",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_SetTags,
			func(opts *AlterBackupSetOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("tag1"), Value: "value1"},
					{Name: NewAccountObjectIdentifier("tag2"), Value: "value2"},
				}
			},
			`ALTER BACKUP SET %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`,
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Alter_UnsetTags,
			func(opts *AlterBackupSetOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("tag1"),
					NewAccountObjectIdentifier("tag2"),
				}
			},
			`ALTER BACKUP SET %s UNSET TAG "tag1", "tag2"`,
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ifExists",
			func(opts *AlterBackupSetOptions) {
				opts.IfExists = new(true)
				opts.AddBackup = new(true)
			},
			"ALTER BACKUP SET IF EXISTS %s ADD BACKUP",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	backupSetsTests.Drop.
		withExpectedSqlf(
			case_BackupSets_sql_Drop_basic,
			"DROP BACKUP SET %s",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Drop_all,
			func(opts *DropBackupSetOptions) { opts.IfExists = new(true) },
			"DROP BACKUP SET IF EXISTS %s",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)

	backupSetsTests.Show.
		withExpectedSqlf(case_BackupSets_sql_Show_basic, "SHOW BACKUP SETS").
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Show_all,
			func(opts *ShowBackupSetOptions) {
				opts.Like = &Like{Pattern: String("like-pattern")}
				opts.In = &In{Account: new(true)}
				opts.StartsWith = String("abc")
				opts.Limit = &LimitFrom{Rows: Int(10)}
			},
			"SHOW BACKUP SETS LIKE 'like-pattern' IN ACCOUNT STARTS WITH 'abc' LIMIT 10",
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Show_Like,
			func(opts *ShowBackupSetOptions) { opts.Like = &Like{Pattern: String("like-pattern")} },
			"SHOW BACKUP SETS LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Show_In,
			func(opts *ShowBackupSetOptions) { opts.In = &In{Account: new(true)} },
			"SHOW BACKUP SETS IN ACCOUNT",
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Show_StartsWith,
			func(opts *ShowBackupSetOptions) { opts.StartsWith = String("abc") },
			"SHOW BACKUP SETS STARTS WITH 'abc'",
		).
		withModifyAndExpectedSqlf(
			case_BackupSets_sql_Show_Limit,
			func(opts *ShowBackupSetOptions) { opts.Limit = &LimitFrom{Rows: Int(10)} },
			"SHOW BACKUP SETS LIMIT 10",
		)

	backupSetsTests.Describe.
		withExpectedSqlf(
			case_BackupSets_sql_Describe_basic,
			"DESCRIBE BACKUP SET %s",
			backupSetsTestIdSchemaObjectIdentifier.FullyQualifiedName(),
		)
}
