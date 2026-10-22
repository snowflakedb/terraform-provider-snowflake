package sdk

func init() {
	id := applicationPackagesTestIdAccountObjectIdentifier
	tagId := randomSchemaObjectIdentifier()
	alterTagId := NewAccountObjectIdentifier("tag1")
	alterTagId2 := NewAccountObjectIdentifier("tag2")

	applicationPackagesTests.Create.
		withExpectedSqlf(
			case_ApplicationPackages_sql_Create_basic,
			"CREATE APPLICATION PACKAGE %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Create_all,
			func(opts *CreateApplicationPackageOptions) {
				opts.IfNotExists = new(true)
				opts.DataRetentionTimeInDays = new(1)
				opts.MaxDataExtensionTimeInDays = new(1)
				opts.DefaultDdlCollation = new("en_US")
				opts.Comment = new("comment")
				opts.Distribution = new(DistributionInternal)
				opts.Tag = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			"CREATE APPLICATION PACKAGE IF NOT EXISTS %s DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 DEFAULT_DDL_COLLATION = 'en_US' COMMENT = 'comment' DISTRIBUTION = INTERNAL TAG (%s = 'v1')",
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		)

	applicationPackagesTests.Alter.
		withDefaultOpts(func() *AlterApplicationPackageOptions {
			return &AlterApplicationPackageOptions{
				IfExists: new(true),
				name:     id,
			}
		}).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Alter_Set,
			func(opts *AlterApplicationPackageOptions) {
				opts.Set = &ApplicationPackageSet{
					DataRetentionTimeInDays:    new(1),
					MaxDataExtensionTimeInDays: new(1),
					DefaultDdlCollation:        new("en_US"),
					Comment:                    new("comment"),
					Distribution:               new(DistributionInternal),
				}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s SET DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 DEFAULT_DDL_COLLATION = 'en_US' COMMENT = 'comment' DISTRIBUTION = INTERNAL`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Alter_Unset,
			func(opts *AlterApplicationPackageOptions) {
				opts.Unset = &ApplicationPackageUnset{
					DataRetentionTimeInDays:    new(true),
					MaxDataExtensionTimeInDays: new(true),
					DefaultDdlCollation:        new(true),
					Comment:                    new(true),
					Distribution:               new(true),
				}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, DEFAULT_DDL_COLLATION, COMMENT, DISTRIBUTION`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Alter_ModifyReleaseDirective,
			func(opts *AlterApplicationPackageOptions) {
				opts.ModifyReleaseDirective = &ModifyReleaseDirective{
					ReleaseDirective: "DEFAULT",
					Version:          "V1",
					Patch:            1,
				}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s MODIFY RELEASE DIRECTIVE DEFAULT VERSION = V1 PATCH = 1`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Alter_SetDefaultReleaseDirective,
			func(opts *AlterApplicationPackageOptions) {
				opts.SetDefaultReleaseDirective = &SetDefaultReleaseDirective{
					Version: "V1",
					Patch:   1,
				}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s SET DEFAULT RELEASE DIRECTIVE VERSION = V1 PATCH = 1`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Alter_SetReleaseDirective,
			func(opts *AlterApplicationPackageOptions) {
				opts.SetReleaseDirective = &SetReleaseDirective{
					ReleaseDirective: "DEFAULT",
					Accounts: []string{
						"org1.acc1",
						"org2.acc2",
					},
					Version: "V1",
					Patch:   1,
				}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s SET RELEASE DIRECTIVE DEFAULT ACCOUNTS = (org1.acc1, org2.acc2) VERSION = V1 PATCH = 1`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Alter_UnsetReleaseDirective,
			func(opts *AlterApplicationPackageOptions) {
				opts.UnsetReleaseDirective = &UnsetReleaseDirective{
					ReleaseDirective: "DEFAULT",
				}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s UNSET RELEASE DIRECTIVE DEFAULT`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Alter_AddVersion,
			func(opts *AlterApplicationPackageOptions) {
				opts.AddVersion = &AddVersion{
					VersionIdentifier: new("v1_1"),
					Using:             "@hello_snowflake_code.core.hello_snowflake_stage",
					Label:             new("test"),
				}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s ADD VERSION v1_1 USING '@hello_snowflake_code.core.hello_snowflake_stage' LABEL = 'test'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Alter_DropVersion,
			func(opts *AlterApplicationPackageOptions) {
				opts.DropVersion = &DropVersion{
					VersionIdentifier: "v1_1",
				}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s DROP VERSION v1_1`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Alter_AddPatchForVersion,
			func(opts *AlterApplicationPackageOptions) {
				opts.AddPatchForVersion = &AddPatchForVersion{
					VersionIdentifier: new("v1_1"),
					Using:             "@hello_snowflake_code.core.hello_snowflake_stage",
					Label:             new("test"),
				}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s ADD PATCH FOR VERSION v1_1 USING '@hello_snowflake_code.core.hello_snowflake_stage' LABEL = 'test'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Alter_SetTags,
			func(opts *AlterApplicationPackageOptions) {
				opts.SetTags = []TagAssociation{{Name: alterTagId, Value: "value1"}}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s SET TAG %s = 'value1'`,
			id.FullyQualifiedName(), alterTagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Alter_UnsetTags,
			func(opts *AlterApplicationPackageOptions) {
				opts.UnsetTags = []ObjectIdentifier{alterTagId, alterTagId2}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), alterTagId.FullyQualifiedName(), alterTagId2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_SetReleaseDirective_noAccounts",
			func(opts *AlterApplicationPackageOptions) {
				opts.SetReleaseDirective = &SetReleaseDirective{
					ReleaseDirective: "DEFAULT",
					Version:          "V1",
					Patch:            1,
				}
			},
			`ALTER APPLICATION PACKAGE IF EXISTS %s SET RELEASE DIRECTIVE DEFAULT ACCOUNTS = () VERSION = V1 PATCH = 1`,
			id.FullyQualifiedName(),
		)

	applicationPackagesTests.Drop.
		withExpectedSqlf(
			case_ApplicationPackages_sql_Drop_basic,
			`DROP APPLICATION PACKAGE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Drop_all,
			func(opts *DropApplicationPackageOptions) { opts.IfExists = new(true) },
			`DROP APPLICATION PACKAGE IF EXISTS %s`, id.FullyQualifiedName(),
		)

	applicationPackagesTests.Show.
		withExpectedSql(case_ApplicationPackages_sql_Show_basic, `SHOW APPLICATION PACKAGES`).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Show_all,
			func(opts *ShowApplicationPackageOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.StartsWith = new("A")
				opts.Limit = &LimitFrom{Rows: new(1), From: new("B")}
			},
			`SHOW APPLICATION PACKAGES LIKE 'pattern' STARTS WITH 'A' LIMIT 1 FROM 'B'`,
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Show_Like,
			func(opts *ShowApplicationPackageOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			`SHOW APPLICATION PACKAGES LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Show_StartsWith,
			func(opts *ShowApplicationPackageOptions) { opts.StartsWith = new("A") },
			`SHOW APPLICATION PACKAGES STARTS WITH 'A'`,
		).
		withModifyAndExpectedSqlf(
			case_ApplicationPackages_sql_Show_Limit,
			func(opts *ShowApplicationPackageOptions) {
				opts.Limit = &LimitFrom{Rows: new(1), From: new("B")}
			},
			`SHOW APPLICATION PACKAGES LIMIT 1 FROM 'B'`,
		)
}
