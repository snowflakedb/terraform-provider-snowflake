package sdk

func init() {
	id := applicationsTestIdAccountObjectIdentifier
	packageId := randomAccountObjectIdentifier()
	tagId := randomSchemaObjectIdentifier()

	applicationsTests.Create.
		withDefaultOpts(func() *CreateApplicationOptions {
			return &CreateApplicationOptions{
				name:        id,
				PackageName: packageId,
			}
		}).
		withAdditionalValidationCase(
			"validation_Create_DebugMode_requiresVersion",
			func(opts *CreateApplicationOptions) { opts.DebugMode = Bool(true) },
			NewError("CreateApplicationOptions.DebugMode can be set only when CreateApplicationOptions.Version is set"),
		).
		withExpectedSqlf(
			case_Applications_sql_Create_basic,
			`CREATE APPLICATION %s FROM APPLICATION PACKAGE %s`,
			id.FullyQualifiedName(), packageId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Create_all,
			func(opts *CreateApplicationOptions) {
				opts.Version = &ApplicationVersion{VersionDirectory: String("@test")}
				opts.DebugMode = Bool(true)
				opts.Comment = String("test")
				opts.Tag = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			`CREATE APPLICATION %s FROM APPLICATION PACKAGE %s USING '@test' DEBUG_MODE = true COMMENT = 'test' TAG (%s = 'v1')`,
			id.FullyQualifiedName(), packageId.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_commentAndTag",
			func(opts *CreateApplicationOptions) {
				opts.Comment = String("test")
				opts.Tag = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			`CREATE APPLICATION %s FROM APPLICATION PACKAGE %s COMMENT = 'test' TAG (%s = 'v1')`,
			id.FullyQualifiedName(), packageId.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_versionAndPatch",
			func(opts *CreateApplicationOptions) {
				opts.Version = &ApplicationVersion{VersionAndPatch: &VersionAndPatch{Version: "V001", Patch: Int(1)}}
				opts.DebugMode = Bool(true)
				opts.Comment = String("test")
				opts.Tag = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			`CREATE APPLICATION %s FROM APPLICATION PACKAGE %s USING VERSION V001 PATCH 1 DEBUG_MODE = true COMMENT = 'test' TAG (%s = 'v1')`,
			id.FullyQualifiedName(), packageId.FullyQualifiedName(), tagId.FullyQualifiedName(),
		)

	applicationsTests.Drop.
		withExpectedSqlf(
			case_Applications_sql_Drop_basic,
			`DROP APPLICATION %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Drop_all,
			func(opts *DropApplicationOptions) {
				opts.IfExists = Bool(true)
				opts.Cascade = Bool(true)
			},
			`DROP APPLICATION IF EXISTS %s CASCADE`, id.FullyQualifiedName(),
		)

	applicationsTests.Alter.
		withAdditionalValidationCase(
			"validation_Alter_IfExists_requiresSetOrUnset",
			func(opts *AlterApplicationOptions) {
				opts.IfExists = Bool(true)
				opts.Upgrade = Bool(true)
			},
			NewError("AlterApplicationOptions.IfExists can be set only when AlterApplicationOptions.Set or AlterApplicationOptions.Unset is set"),
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Alter_Set,
			func(opts *AlterApplicationOptions) {
				opts.IfExists = Bool(true)
				opts.Set = &ApplicationSet{
					ShareEventsWithProvider: Bool(true),
					DebugMode:               Bool(true),
					Comment:                 String("test"),
				}
			},
			`ALTER APPLICATION IF EXISTS %s SET COMMENT = 'test' SHARE_EVENTS_WITH_PROVIDER = true DEBUG_MODE = true`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Alter_Unset,
			func(opts *AlterApplicationOptions) {
				opts.IfExists = Bool(true)
				opts.Unset = &ApplicationUnset{Comment: Bool(true)}
			},
			`ALTER APPLICATION IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Alter_Upgrade,
			func(opts *AlterApplicationOptions) { opts.Upgrade = Bool(true) },
			`ALTER APPLICATION %s UPGRADE`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Alter_UpgradeVersion,
			func(opts *AlterApplicationOptions) {
				opts.UpgradeVersion = &ApplicationVersion{VersionDirectory: String("@test")}
			},
			`ALTER APPLICATION %s UPGRADE USING '@test'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Alter_UnsetReferences,
			func(opts *AlterApplicationOptions) { opts.UnsetReferences = &ApplicationReferences{} },
			`ALTER APPLICATION %s UNSET REFERENCES`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Alter_SetTags,
			func(opts *AlterApplicationOptions) {
				opts.SetTags = []TagAssociation{{Name: NewAccountObjectIdentifier("tag1"), Value: "value1"}}
			},
			`ALTER APPLICATION %s SET TAG "tag1" = 'value1'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Alter_UnsetTags,
			func(opts *AlterApplicationOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("tag1"),
					NewAccountObjectIdentifier("tag2"),
				}
			},
			`ALTER APPLICATION %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName(),
		)

	applicationsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_UnsetShareEventsWithProvider",
			func(opts *AlterApplicationOptions) {
				opts.IfExists = Bool(true)
				opts.Unset = &ApplicationUnset{ShareEventsWithProvider: Bool(true)}
			},
			`ALTER APPLICATION IF EXISTS %s UNSET SHARE_EVENTS_WITH_PROVIDER`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_UnsetDebugMode",
			func(opts *AlterApplicationOptions) {
				opts.IfExists = Bool(true)
				opts.Unset = &ApplicationUnset{DebugMode: Bool(true)}
			},
			`ALTER APPLICATION IF EXISTS %s UNSET DEBUG_MODE`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_UpgradeVersionAndPatch",
			func(opts *AlterApplicationOptions) {
				opts.UpgradeVersion = &ApplicationVersion{VersionAndPatch: &VersionAndPatch{Version: "V001", Patch: Int(1)}}
			},
			`ALTER APPLICATION %s UPGRADE USING VERSION V001 PATCH 1`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_UnsetReferencesList",
			func(opts *AlterApplicationOptions) {
				opts.UnsetReferences = &ApplicationReferences{
					References: []ApplicationReference{
						{Reference: "ref1"},
						{Reference: "ref2"},
					},
				}
			},
			`ALTER APPLICATION %s UNSET REFERENCES ('ref1', 'ref2')`, id.FullyQualifiedName(),
		)

	applicationsTests.Show.
		withExpectedSql(case_Applications_sql_Show_basic, `SHOW APPLICATIONS`).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Show_all,
			func(opts *ShowApplicationOptions) {
				opts.Like = &Like{Pattern: String("pattern")}
				opts.StartsWith = String("A")
				opts.Limit = &LimitFrom{Rows: Int(1), From: String("B")}
			},
			`SHOW APPLICATIONS LIKE 'pattern' STARTS WITH 'A' LIMIT 1 FROM 'B'`,
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Show_Like,
			func(opts *ShowApplicationOptions) { opts.Like = &Like{Pattern: String("pattern")} },
			`SHOW APPLICATIONS LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Show_StartsWith,
			func(opts *ShowApplicationOptions) { opts.StartsWith = String("A") },
			`SHOW APPLICATIONS STARTS WITH 'A'`,
		).
		withModifyAndExpectedSqlf(
			case_Applications_sql_Show_Limit,
			func(opts *ShowApplicationOptions) {
				opts.Limit = &LimitFrom{Rows: Int(1), From: String("B")}
			},
			`SHOW APPLICATIONS LIMIT 1 FROM 'B'`,
		)

	applicationsTests.Describe.
		withExpectedSqlf(
			case_Applications_sql_Describe_basic,
			`DESCRIBE APPLICATION %s`, id.FullyQualifiedName(),
		)
}
