package sdk

func sampleListingManifest() string {
	return `
title: "MyListing"
subtitle: "Subtitle for MyListing"
description: "Description for MyListing"
listing_terms:
   type: "STANDARD"
targets:
    accounts: ["Org1.Account1"]
usage_examples:
    - title: "this is a test sql"
      description: "Simple example"
      query: "select *"
`
}

func init() {
	id := listingsTestIdAccountObjectIdentifier
	stageId := randomSchemaObjectIdentifier()
	stageLocation := NewStageLocation(stageId, "dir/subdir")
	var from Location = stageLocation
	manifest := sampleListingManifest()
	shareId := randomAccountObjectIdentifier()
	applicationPackageId := randomAccountObjectIdentifier()
	renameTarget := randomAccountObjectIdentifier()

	listingsTests.Create.
		withModify(
			case_Listings_validation_Create_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateListingOptions) {
				opts.As = &manifest
				opts.From = &from
			},
		).
		withModify(
			case_Listings_validation_Create_opts_With_ExactlyOneValueSet_NoneSet,
			func(opts *CreateListingOptions) {
				opts.As = &manifest
				opts.With = &ListingWith{}
			},
		).
		withModify(
			case_Listings_validation_Create_opts_With_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateListingOptions) {
				opts.As = &manifest
				opts.With = &ListingWith{
					Share:              &shareId,
					ApplicationPackage: &applicationPackageId,
				}
			},
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Create_basic,
			func(opts *CreateListingOptions) { opts.As = &manifest },
			"CREATE EXTERNAL LISTING %s AS $$%s$$", id.FullyQualifiedName(), manifest,
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Create_all,
			func(opts *CreateListingOptions) {
				opts.IfNotExists = new(true)
				opts.With = &ListingWith{Share: &shareId}
				opts.As = &manifest
				opts.Publish = new(true)
				opts.Review = new(true)
				opts.Comment = new("comment")
			},
			"CREATE EXTERNAL LISTING IF NOT EXISTS %s SHARE %s AS $$%s$$ PUBLISH = true REVIEW = true COMMENT = 'comment'",
			id.FullyQualifiedName(), shareId.FullyQualifiedName(), manifest,
		).
		withAdditionalSqlCasef(
			"sql_Create_From",
			func(opts *CreateListingOptions) { opts.From = &from },
			`CREATE EXTERNAL LISTING %s FROM '@\"%s\".\"%s\".\"%s\"/dir/subdir'`,
			id.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		).
		withAdditionalSqlCasef(
			"sql_Create_all_As_applicationPackage",
			func(opts *CreateListingOptions) {
				opts.IfNotExists = new(true)
				opts.With = &ListingWith{ApplicationPackage: &applicationPackageId}
				opts.As = &manifest
				opts.Publish = new(true)
				opts.Review = new(true)
				opts.Comment = new("comment")
			},
			"CREATE EXTERNAL LISTING IF NOT EXISTS %s APPLICATION PACKAGE %s AS $$%s$$ PUBLISH = true REVIEW = true COMMENT = 'comment'",
			id.FullyQualifiedName(), applicationPackageId.FullyQualifiedName(), manifest,
		).
		withAdditionalSqlCasef(
			"sql_Create_all_From_share",
			func(opts *CreateListingOptions) {
				opts.IfNotExists = new(true)
				opts.With = &ListingWith{Share: &shareId}
				opts.From = &from
				opts.Publish = new(true)
				opts.Review = new(true)
				opts.Comment = new("comment")
			},
			`CREATE EXTERNAL LISTING IF NOT EXISTS %s SHARE %s FROM '@\"%s\".\"%s\".\"%s\"/dir/subdir' PUBLISH = true REVIEW = true COMMENT = 'comment'`,
			id.FullyQualifiedName(), shareId.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		).
		withAdditionalSqlCasef(
			"sql_Create_all_From_applicationPackage",
			func(opts *CreateListingOptions) {
				opts.IfNotExists = new(true)
				opts.With = &ListingWith{ApplicationPackage: &applicationPackageId}
				opts.From = &from
				opts.Publish = new(true)
				opts.Review = new(true)
				opts.Comment = new("comment")
			},
			`CREATE EXTERNAL LISTING IF NOT EXISTS %s APPLICATION PACKAGE %s FROM '@\"%s\".\"%s\".\"%s\"/dir/subdir' PUBLISH = true REVIEW = true COMMENT = 'comment'`,
			id.FullyQualifiedName(), applicationPackageId.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		)

	listingsTests.Alter.
		withModifyAndExpectedSqlf(
			case_Listings_sql_Alter_Publish,
			func(opts *AlterListingOptions) {
				opts.IfExists = new(true)
				opts.Publish = new(true)
			},
			"ALTER LISTING IF EXISTS %s PUBLISH", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Alter_Unpublish,
			func(opts *AlterListingOptions) {
				opts.IfExists = new(true)
				opts.Unpublish = new(true)
			},
			"ALTER LISTING IF EXISTS %s UNPUBLISH", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Alter_Review,
			func(opts *AlterListingOptions) {
				opts.IfExists = new(true)
				opts.Review = new(true)
			},
			"ALTER LISTING IF EXISTS %s REVIEW", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Alter_AlterListingAs,
			func(opts *AlterListingOptions) {
				opts.IfExists = new(true)
				opts.AlterListingAs = &AlterListingAs{As: manifest}
			},
			"ALTER LISTING IF EXISTS %s AS $$%s$$", id.FullyQualifiedName(), manifest,
		).
		withAdditionalSqlCasef(
			"sql_Alter_AlterListingAs_complete",
			func(opts *AlterListingOptions) {
				opts.IfExists = new(true)
				opts.AlterListingAs = &AlterListingAs{
					As:      manifest,
					Publish: new(true),
					Review:  new(true),
					Comment: new("comment"),
				}
			},
			"ALTER LISTING IF EXISTS %s AS $$%s$$ PUBLISH = true REVIEW = true COMMENT = 'comment'",
			id.FullyQualifiedName(), manifest,
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Alter_AddVersion,
			func(opts *AlterListingOptions) {
				opts.AddVersion = &AddListingVersion{
					IfNotExists: new(true),
					VersionName: "version-name",
					From:        stageLocation,
					Comment:     new("comment"),
				}
			},
			`ALTER LISTING %s ADD VERSION IF NOT EXISTS "version-name" FROM '@\"%s\".\"%s\".\"%s\"/dir/subdir' COMMENT = 'comment'`,
			id.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Alter_RenameTo,
			func(opts *AlterListingOptions) { opts.RenameTo = &renameTarget },
			"ALTER LISTING %s RENAME TO %s", id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Alter_Set,
			func(opts *AlterListingOptions) { opts.Set = &ListingSet{Comment: new("comment")} },
			"ALTER LISTING %s SET COMMENT = 'comment'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Alter_Unset,
			func(opts *AlterListingOptions) { opts.Unset = &ListingUnset{Comment: new(true)} },
			"ALTER LISTING %s UNSET COMMENT", id.FullyQualifiedName(),
		)

	listingsTests.Drop.
		withExpectedSqlf(
			case_Listings_sql_Drop_basic,
			"DROP LISTING %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Drop_all,
			func(opts *DropListingOptions) { opts.IfExists = new(true) },
			"DROP LISTING IF EXISTS %s", id.FullyQualifiedName(),
		)

	listingsTests.Show.
		withExpectedSql(case_Listings_sql_Show_basic, "SHOW LISTINGS").
		withModifyAndExpectedSqlf(
			case_Listings_sql_Show_all,
			func(opts *ShowListingOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.StartsWith = new("startsWith")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("from")}
			},
			"SHOW LISTINGS LIKE 'pattern' STARTS WITH 'startsWith' LIMIT 10 FROM 'from'",
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Show_Like,
			func(opts *ShowListingOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			"SHOW LISTINGS LIKE 'pattern'",
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Show_StartsWith,
			func(opts *ShowListingOptions) { opts.StartsWith = new("startsWith") },
			"SHOW LISTINGS STARTS WITH 'startsWith'",
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_Show_Limit,
			func(opts *ShowListingOptions) {
				opts.Limit = &LimitFrom{Rows: new(10), From: new("from")}
			},
			"SHOW LISTINGS LIMIT 10 FROM 'from'",
		)

	listingsTests.Describe.
		withExpectedSqlf(
			case_Listings_sql_Describe_basic,
			"DESCRIBE LISTING %s", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Describe_all",
			func(opts *DescribeListingOptions) { opts.Revision = new(ListingRevisionDraft) },
			"DESCRIBE LISTING %s REVISION = DRAFT", id.FullyQualifiedName(),
		)

	listingsTests.ShowVersions.
		withExpectedSqlf(
			case_Listings_sql_ShowVersions_basic,
			"SHOW VERSIONS IN LISTING %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_ShowVersions_all,
			func(opts *ShowVersionsListingOptions) { opts.Limit = &LimitFrom{Rows: new(5)} },
			"SHOW VERSIONS IN LISTING %s LIMIT 5", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Listings_sql_ShowVersions_Limit,
			func(opts *ShowVersionsListingOptions) { opts.Limit = &LimitFrom{Rows: new(5)} },
			"SHOW VERSIONS IN LISTING %s LIMIT 5", id.FullyQualifiedName(),
		)
}
