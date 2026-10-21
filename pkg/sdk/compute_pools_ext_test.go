package sdk

func init() {
	id := computePoolsTestIdAccountObjectIdentifier
	appId := randomAccountObjectIdentifier()
	tagId := NewAccountObjectIdentifier("tag1")
	tagId2 := NewAccountObjectIdentifier("tag2")
	setComment := "set-comment"
	backupInstanceFamilies := []ComputePoolInstanceFamilyListItem{
		{Value: ComputePoolInstanceFamilyCpuX64M},
		{Value: ComputePoolInstanceFamilyCpuX64S},
	}

	computePoolsTests.Create.
		withDefaultOpts(func() *CreateComputePoolOptions {
			return &CreateComputePoolOptions{
				name:           id,
				MinNodes:       1,
				MaxNodes:       3,
				InstanceFamily: ComputePoolInstanceFamilyCpuX64Xs,
			}
		}).
		withAdditionalValidationCase(
			"validation_Create_MinNodes_greaterThan0",
			func(opts *CreateComputePoolOptions) { opts.MinNodes = 0 },
			errIntValue("CreateComputePoolOptions", "MinNodes", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Create_MaxNodes_greaterOrEqualMinNodes",
			func(opts *CreateComputePoolOptions) { opts.MinNodes = 2; opts.MaxNodes = 1 },
			errIntValue("CreateComputePoolOptions", "MaxNodes", IntErrGreaterOrEqual, 2),
		).
		withExpectedSqlf(
			case_ComputePools_sql_Create_basic,
			"CREATE COMPUTE POOL %s MIN_NODES = 1 MAX_NODES = 3 INSTANCE_FAMILY = CPU_X64_XS",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Create_all,
			func(opts *CreateComputePoolOptions) {
				opts.IfNotExists = new(true)
				opts.ForApplication = &appId
				opts.MinNodes = 2
				opts.MaxNodes = 3
				opts.AutoResume = new(true)
				opts.InitiallySuspended = new(true)
				opts.AutoSuspendSecs = new(42)
				opts.Tag = []TagAssociation{{Name: tagId, Value: "value1"}}
				opts.Comment = &setComment
				opts.BackupInstanceFamilies = backupInstanceFamilies
			},
			"CREATE COMPUTE POOL IF NOT EXISTS %s FOR APPLICATION %s MIN_NODES = 2 MAX_NODES = 3 INSTANCE_FAMILY = CPU_X64_XS AUTO_RESUME = true INITIALLY_SUSPENDED = true AUTO_SUSPEND_SECS = 42 TAG (%s = 'value1') COMMENT = '%s' BACKUP_INSTANCE_FAMILIES = ('CPU_X64_M', 'CPU_X64_S')",
			id.FullyQualifiedName(), appId.FullyQualifiedName(), tagId.FullyQualifiedName(), setComment,
		).
		withAdditionalSqlCasef(
			"sql_Create_backupInstanceFamilies",
			func(opts *CreateComputePoolOptions) { opts.BackupInstanceFamilies = backupInstanceFamilies },
			"CREATE COMPUTE POOL %s MIN_NODES = 1 MAX_NODES = 3 INSTANCE_FAMILY = CPU_X64_XS BACKUP_INSTANCE_FAMILIES = ('CPU_X64_M', 'CPU_X64_S')",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_backupInstanceFamilies_singleEntry",
			func(opts *CreateComputePoolOptions) {
				opts.BackupInstanceFamilies = []ComputePoolInstanceFamilyListItem{{Value: ComputePoolInstanceFamilyCpuX64S}}
			},
			"CREATE COMPUTE POOL %s MIN_NODES = 1 MAX_NODES = 3 INSTANCE_FAMILY = CPU_X64_XS BACKUP_INSTANCE_FAMILIES = ('CPU_X64_S')",
			id.FullyQualifiedName(),
		)

	computePoolsTests.Alter.
		withAdditionalValidationCase(
			"validation_Alter_Set_MinNodes_greaterThan0",
			func(opts *AlterComputePoolOptions) { opts.Set = &ComputePoolSet{MinNodes: new(0)} },
			errIntValue("AlterComputePoolOptions", "Set.MinNodes", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_MaxNodes_greaterThan0",
			func(opts *AlterComputePoolOptions) { opts.Set = &ComputePoolSet{MaxNodes: new(0)} },
			errIntValue("AlterComputePoolOptions", "Set.MaxNodes", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_MaxNodes_greaterOrEqualMinNodes",
			func(opts *AlterComputePoolOptions) {
				opts.Set = &ComputePoolSet{MinNodes: new(2), MaxNodes: new(1)}
			},
			errIntValue("AlterComputePoolOptions", "Set.MaxNodes", IntErrGreaterOrEqual, 2),
		).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Alter_Resume,
			func(opts *AlterComputePoolOptions) { opts.Resume = new(true) },
			"ALTER COMPUTE POOL %s RESUME", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Alter_Suspend,
			func(opts *AlterComputePoolOptions) { opts.Suspend = new(true) },
			"ALTER COMPUTE POOL %s SUSPEND", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Alter_StopAll,
			func(opts *AlterComputePoolOptions) { opts.StopAll = new(true) },
			"ALTER COMPUTE POOL %s STOP ALL", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Alter_Set,
			func(opts *AlterComputePoolOptions) {
				opts.Set = &ComputePoolSet{
					MinNodes:               new(2),
					MaxNodes:               new(3),
					AutoResume:             new(true),
					AutoSuspendSecs:        new(60),
					BackupInstanceFamilies: backupInstanceFamilies,
					Comment:                &setComment,
				}
			},
			"ALTER COMPUTE POOL %s SET MIN_NODES = 2 MAX_NODES = 3 AUTO_RESUME = true AUTO_SUSPEND_SECS = 60 BACKUP_INSTANCE_FAMILIES = ('CPU_X64_M', 'CPU_X64_S') COMMENT = '%s'",
			id.FullyQualifiedName(), setComment,
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_backupInstanceFamilies",
			func(opts *AlterComputePoolOptions) {
				opts.Set = &ComputePoolSet{BackupInstanceFamilies: backupInstanceFamilies}
			},
			"ALTER COMPUTE POOL %s SET BACKUP_INSTANCE_FAMILIES = ('CPU_X64_M', 'CPU_X64_S')",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Alter_Unset,
			func(opts *AlterComputePoolOptions) {
				opts.Unset = &ComputePoolUnset{
					AutoResume:             new(true),
					AutoSuspendSecs:        new(true),
					BackupInstanceFamilies: new(true),
					Comment:                new(true),
				}
			},
			"ALTER COMPUTE POOL %s UNSET AUTO_RESUME, AUTO_SUSPEND_SECS, BACKUP_INSTANCE_FAMILIES, COMMENT", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_backupInstanceFamilies",
			func(opts *AlterComputePoolOptions) {
				opts.Unset = &ComputePoolUnset{BackupInstanceFamilies: new(true)}
			},
			"ALTER COMPUTE POOL %s UNSET BACKUP_INSTANCE_FAMILIES", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Alter_SetTags,
			func(opts *AlterComputePoolOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId, Value: "value1"},
					{Name: tagId2, Value: "value2"},
				}
			},
			`ALTER COMPUTE POOL %s SET TAG %s = 'value1', %s = 'value2'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Alter_UnsetTags,
			func(opts *AlterComputePoolOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			`ALTER COMPUTE POOL %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		)

	computePoolsTests.Drop.
		withExpectedSqlf(case_ComputePools_sql_Drop_basic,
			"DROP COMPUTE POOL %s", id.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Drop_all,
			func(opts *DropComputePoolOptions) { opts.IfExists = new(true) },
			"DROP COMPUTE POOL IF EXISTS %s", id.FullyQualifiedName(),
		)

	computePoolsTests.Show.
		withExpectedSql(case_ComputePools_sql_Show_basic, "SHOW COMPUTE POOLS").
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Show_all,
			func(opts *ShowComputePoolOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.StartsWith = new("prefix")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("from")}
			},
			"SHOW COMPUTE POOLS LIKE 'pattern' STARTS WITH 'prefix' LIMIT 10 FROM 'from'",
		).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Show_Like,
			func(opts *ShowComputePoolOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			"SHOW COMPUTE POOLS LIKE 'pattern'",
		).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Show_StartsWith,
			func(opts *ShowComputePoolOptions) { opts.StartsWith = new("prefix") },
			"SHOW COMPUTE POOLS STARTS WITH 'prefix'",
		).
		withModifyAndExpectedSqlf(
			case_ComputePools_sql_Show_Limit,
			func(opts *ShowComputePoolOptions) {
				opts.Limit = &LimitFrom{Rows: new(10), From: new("from")}
			},
			"SHOW COMPUTE POOLS LIMIT 10 FROM 'from'",
		)

	computePoolsTests.Describe.
		withExpectedSqlf(case_ComputePools_sql_Describe_basic,
			"DESCRIBE COMPUTE POOL %s", id.FullyQualifiedName())
}
