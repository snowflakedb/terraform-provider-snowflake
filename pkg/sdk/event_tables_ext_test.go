package sdk

func init() {
	id := eventTablesTestIdSchemaObjectIdentifier
	rowAccessPolicyId := randomSchemaObjectIdentifier()
	rowAccessPolicyId2 := randomSchemaObjectIdentifier()
	tagId := randomSchemaObjectIdentifier()
	renameTarget := randomSchemaObjectIdentifier()
	tagName1 := NewAccountObjectIdentifier("tag1")
	tagName2 := NewAccountObjectIdentifier("tag2")

	eventTablesTests.Create.
		withExpectedSqlf(
			case_EventTables_sql_Create_basic,
			`CREATE EVENT TABLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Create_all,
			func(opts *CreateEventTableOptions) {
				opts.IfNotExists = new(true)
				opts.ClusterBy = []string{"a", "b"}
				opts.DataRetentionTimeInDays = new(1)
				opts.MaxDataExtensionTimeInDays = new(2)
				opts.ChangeTracking = new(true)
				opts.DefaultDdlCollation = new("en_US")
				opts.Comment = new("test")
				opts.RowAccessPolicy = &TableRowAccessPolicyLegacy{
					Name: rowAccessPolicyId,
					On:   []string{"c1", "c2"},
				}
				opts.Tag = []TagAssociation{
					{
						Name:  tagId,
						Value: "v1",
					},
				}
			},
			`CREATE EVENT TABLE IF NOT EXISTS %s CLUSTER BY (a, b) DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 2 CHANGE_TRACKING = true DEFAULT_DDL_COLLATION = 'en_US' COMMENT = 'test' ROW ACCESS POLICY %s ON (c1, c2) TAG (%s = 'v1')`,
			id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_emptyClusterBy",
			func(opts *CreateEventTableOptions) { opts.ClusterBy = []string{} },
			`CREATE EVENT TABLE %s`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateEventTableOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE EVENT TABLE %s COPY GRANTS`, id.FullyQualifiedName(),
		)

	databaseId := NewAccountObjectIdentifier("database")
	likePattern := id.Name()

	eventTablesTests.Show.
		withExpectedSql(case_EventTables_sql_Show_basic, `SHOW EVENT TABLES`).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Show_all,
			func(opts *ShowEventTableOptions) {
				opts.Terse = new(true)
				opts.Like = &Like{Pattern: &likePattern}
				opts.In = &In{Database: databaseId}
				opts.StartsWith = new("prefix")
				opts.Limit = &LimitFrom{Rows: new(10)}
			},
			`SHOW TERSE EVENT TABLES LIKE '%s' IN DATABASE "database" STARTS WITH 'prefix' LIMIT 10`, likePattern,
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Show_Like,
			func(opts *ShowEventTableOptions) { opts.Like = &Like{Pattern: &likePattern} },
			`SHOW EVENT TABLES LIKE '%s'`, likePattern,
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Show_In,
			func(opts *ShowEventTableOptions) { opts.In = &In{Database: databaseId} },
			`SHOW EVENT TABLES IN DATABASE "database"`,
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Show_StartsWith,
			func(opts *ShowEventTableOptions) { opts.StartsWith = new("prefix") },
			`SHOW EVENT TABLES STARTS WITH 'prefix'`,
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Show_Limit,
			func(opts *ShowEventTableOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			`SHOW EVENT TABLES LIMIT 10`,
		).
		withAdditionalSqlCasef(
			"sql_Show_terseLikeAndIn",
			func(opts *ShowEventTableOptions) {
				opts.Terse = new(true)
				opts.Like = &Like{Pattern: &likePattern}
				opts.In = &In{Database: databaseId}
			},
			`SHOW TERSE EVENT TABLES LIKE '%s' IN DATABASE "database"`, likePattern,
		)

	eventTablesTests.Describe.
		withExpectedSqlf(
			case_EventTables_sql_Describe_basic,
			`DESCRIBE EVENT TABLE %s`, id.FullyQualifiedName(),
		)

	eventTablesTests.Drop.
		withExpectedSqlf(
			case_EventTables_sql_Drop_basic,
			`DROP TABLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Drop_all,
			func(opts *DropEventTableOptions) {
				opts.IfExists = new(true)
				opts.Restrict = new(true)
			},
			`DROP TABLE IF EXISTS %s RESTRICT`, id.FullyQualifiedName(),
		)

	eventTablesTests.Alter.
		withDefaultOpts(func() *AlterEventTableOptions {
			return &AlterEventTableOptions{
				name:        id,
				IfNotExists: new(true),
			}
		}).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Alter_RenameTo,
			func(opts *AlterEventTableOptions) { opts.RenameTo = &renameTarget },
			`ALTER TABLE IF NOT EXISTS %s RENAME TO %s`, id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Alter_Set,
			func(opts *AlterEventTableOptions) {
				opts.Set = &EventTableSet{DataRetentionTimeInDays: new(1)}
			},
			`ALTER TABLE IF NOT EXISTS %s SET DATA_RETENTION_TIME_IN_DAYS = 1`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_maxDataExtensionTimeInDays",
			func(opts *AlterEventTableOptions) {
				opts.Set = &EventTableSet{MaxDataExtensionTimeInDays: new(1)}
			},
			`ALTER TABLE IF NOT EXISTS %s SET MAX_DATA_EXTENSION_TIME_IN_DAYS = 1`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_changeTracking",
			func(opts *AlterEventTableOptions) {
				opts.Set = &EventTableSet{ChangeTracking: new(true)}
			},
			`ALTER TABLE IF NOT EXISTS %s SET CHANGE_TRACKING = true`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_comment",
			func(opts *AlterEventTableOptions) {
				opts.Set = &EventTableSet{Comment: new("comment")}
			},
			`ALTER TABLE IF NOT EXISTS %s SET COMMENT = 'comment'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Alter_Unset,
			func(opts *AlterEventTableOptions) {
				opts.Unset = &EventTableUnset{DataRetentionTimeInDays: new(true)}
			},
			`ALTER TABLE IF NOT EXISTS %s UNSET DATA_RETENTION_TIME_IN_DAYS`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_maxDataExtensionTimeInDays",
			func(opts *AlterEventTableOptions) {
				opts.Unset = &EventTableUnset{MaxDataExtensionTimeInDays: new(true)}
			},
			`ALTER TABLE IF NOT EXISTS %s UNSET MAX_DATA_EXTENSION_TIME_IN_DAYS`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_changeTracking",
			func(opts *AlterEventTableOptions) {
				opts.Unset = &EventTableUnset{ChangeTracking: new(true)}
			},
			`ALTER TABLE IF NOT EXISTS %s UNSET CHANGE_TRACKING`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_comment",
			func(opts *AlterEventTableOptions) {
				opts.Unset = &EventTableUnset{Comment: new(true)}
			},
			`ALTER TABLE IF NOT EXISTS %s UNSET COMMENT`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Alter_SetTags,
			func(opts *AlterEventTableOptions) {
				opts.SetTags = []TagAssociation{{Name: tagName1, Value: "value1"}}
			},
			`ALTER TABLE IF NOT EXISTS %s SET TAG "tag1" = 'value1'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Alter_UnsetTags,
			func(opts *AlterEventTableOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagName1, tagName2}
			},
			`ALTER TABLE IF NOT EXISTS %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Alter_AddRowAccessPolicy,
			func(opts *AlterEventTableOptions) {
				opts.AddRowAccessPolicy = &EventTableAddRowAccessPolicy{
					RowAccessPolicy: rowAccessPolicyId,
					On:              []string{"a", "b"},
				}
			},
			`ALTER TABLE IF NOT EXISTS %s ADD ROW ACCESS POLICY %s ON (a, b)`, id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Alter_DropRowAccessPolicy,
			func(opts *AlterEventTableOptions) {
				opts.DropRowAccessPolicy = &EventTableDropRowAccessPolicy{
					RowAccessPolicy: rowAccessPolicyId,
				}
			},
			`ALTER TABLE IF NOT EXISTS %s DROP ROW ACCESS POLICY %s`, id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Alter_DropAndAddRowAccessPolicy,
			func(opts *AlterEventTableOptions) {
				opts.DropAndAddRowAccessPolicy = &EventTableDropAndAddRowAccessPolicy{
					Drop: EventTableDropRowAccessPolicy{RowAccessPolicy: rowAccessPolicyId},
					Add: EventTableAddRowAccessPolicy{
						RowAccessPolicy: rowAccessPolicyId2,
						On:              []string{"a", "b"},
					},
				}
			},
			`ALTER TABLE IF NOT EXISTS %s DROP ROW ACCESS POLICY %s, ADD ROW ACCESS POLICY %s ON (a, b)`,
			id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName(), rowAccessPolicyId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Alter_DropAllRowAccessPolicies,
			func(opts *AlterEventTableOptions) { opts.DropAllRowAccessPolicies = new(true) },
			`ALTER TABLE IF NOT EXISTS %s DROP ALL ROW ACCESS POLICIES`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Alter_ClusteringAction,
			func(opts *AlterEventTableOptions) {
				opts.ClusteringAction = &EventTableClusteringAction{ClusterBy: []string{"a", "b"}}
			},
			`ALTER TABLE IF NOT EXISTS %s CLUSTER BY (a, b)`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ClusteringAction_suspendRecluster",
			func(opts *AlterEventTableOptions) {
				opts.ClusteringAction = &EventTableClusteringAction{SuspendRecluster: new(true)}
			},
			`ALTER TABLE IF NOT EXISTS %s SUSPEND RECLUSTER`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ClusteringAction_resumeRecluster",
			func(opts *AlterEventTableOptions) {
				opts.ClusteringAction = &EventTableClusteringAction{ResumeRecluster: new(true)}
			},
			`ALTER TABLE IF NOT EXISTS %s RESUME RECLUSTER`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_ClusteringAction_dropClusteringKey",
			func(opts *AlterEventTableOptions) {
				opts.ClusteringAction = &EventTableClusteringAction{DropClusteringKey: new(true)}
			},
			`ALTER TABLE IF NOT EXISTS %s DROP CLUSTERING KEY`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_EventTables_sql_Alter_SearchOptimizationAction,
			func(opts *AlterEventTableOptions) {
				opts.SearchOptimizationAction = &EventTableSearchOptimizationAction{
					Add: &SearchOptimization{On: []string{"EQUALITY(*)", "SUBSTRING(*)"}},
				}
			},
			`ALTER TABLE IF NOT EXISTS %s ADD SEARCH OPTIMIZATION ON EQUALITY(*), SUBSTRING(*)`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_SearchOptimizationAction_drop",
			func(opts *AlterEventTableOptions) {
				opts.SearchOptimizationAction = &EventTableSearchOptimizationAction{
					Drop: &SearchOptimization{On: []string{"EQUALITY(*)", "SUBSTRING(*)"}},
				}
			},
			`ALTER TABLE IF NOT EXISTS %s DROP SEARCH OPTIMIZATION ON EQUALITY(*), SUBSTRING(*)`, id.FullyQualifiedName(),
		)
}
