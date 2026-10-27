package sdk

func init() {
	id := codeBundlesTestIdSchemaObjectIdentifier
	from := `snow://workspace/USER$.PUBLIC.ws/versions/last`
	comment := "comment"

	codeBundlesTests.Create.
		withDefaultOpts(func() *CreateCodeBundleOptions {
			return &CreateCodeBundleOptions{
				name: id,
				From: from,
			}
		}).
		withExpectedSqlf(
			case_CodeBundles_sql_Create_basic,
			`CREATE CODE BUNDLE %s FROM '%s'`,
			id.FullyQualifiedName(), from,
		).
		withModifyAndExpectedSqlf(
			case_CodeBundles_sql_Create_all,
			func(opts *CreateCodeBundleOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new(comment)
			},
			`CREATE CODE BUNDLE IF NOT EXISTS %s FROM '%s' COMMENT = '%s'`,
			id.FullyQualifiedName(), from, comment,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateCodeBundleOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE CODE BUNDLE %s FROM '%s'`,
			id.FullyQualifiedName(), from,
		)

	codeBundlesTests.Show.
		withExpectedSql(case_CodeBundles_sql_Show_basic, `SHOW CODE BUNDLES`).
		withModifyAndExpectedSqlf(
			case_CodeBundles_sql_Show_In,
			func(opts *ShowCodeBundleOptions) {
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
			},
			`SHOW CODE BUNDLES IN ACCOUNT`,
		).
		withModifyAndExpectedSqlf(
			case_CodeBundles_sql_Show_all,
			func(opts *ShowCodeBundleOptions) {
				opts.In = &ExtendedIn{In: In{Schema: id.SchemaId()}}
			},
			`SHOW CODE BUNDLES IN SCHEMA %s`,
			id.SchemaId().FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Show_In_Database",
			func(opts *ShowCodeBundleOptions) {
				opts.In = &ExtendedIn{In: In{Database: id.DatabaseId()}}
			},
			`SHOW CODE BUNDLES IN DATABASE %s`,
			id.DatabaseId().FullyQualifiedName(),
		)

	codeBundlesTests.Alter.
		withModifyAndExpectedSqlf(
			case_CodeBundles_sql_Alter_AddVersionFrom,
			func(opts *AlterCodeBundleOptions) { opts.AddVersionFrom = new(from) },
			`ALTER CODE BUNDLE %s ADD VERSION FROM '%s'`,
			id.FullyQualifiedName(), from,
		).
		withModifyAndExpectedSqlf(
			case_CodeBundles_sql_Alter_Set,
			func(opts *AlterCodeBundleOptions) {
				opts.Set = &CodeBundleSet{
					Comment: new(comment),
				}
			},
			`ALTER CODE BUNDLE %s SET COMMENT = '%s'`,
			id.FullyQualifiedName(), comment,
		).
		withModifyAndExpectedSqlf(
			case_CodeBundles_sql_Alter_Unset,
			func(opts *AlterCodeBundleOptions) {
				opts.Unset = &CodeBundleUnset{
					Comment: new(true),
				}
			},
			`ALTER CODE BUNDLE %s UNSET COMMENT`,
			id.FullyQualifiedName(),
		)

	codeBundlesTests.Drop.
		withExpectedSqlf(
			case_CodeBundles_sql_Drop_basic,
			`DROP CODE BUNDLE %s`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_CodeBundles_sql_Drop_all,
			func(opts *DropCodeBundleOptions) { opts.IfExists = new(true) },
			`DROP CODE BUNDLE IF EXISTS %s`,
			id.FullyQualifiedName(),
		)

	codeBundlesTests.Describe.
		withExpectedSqlf(
			case_CodeBundles_sql_Describe_basic,
			`DESCRIBE CODE BUNDLE %s`,
			id.FullyQualifiedName(),
		)
}
