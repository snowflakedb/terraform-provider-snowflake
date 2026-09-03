package sdk

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func init() {
	id := imageRepositoriesTestIdSchemaObjectIdentifier
	schemaId := randomDatabaseObjectIdentifier()
	tagId := NewAccountObjectIdentifier("tag1")
	tagId2 := NewAccountObjectIdentifier("tag2")
	createComment := random.Comment()

	imageRepositoriesTests.Create.
		withExpectedSqlf(
			case_ImageRepositories_sql_Create_basic,
			"CREATE IMAGE REPOSITORY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ImageRepositories_sql_Create_all,
			func(opts *CreateImageRepositoryOptions) {
				opts.IfNotExists = new(true)
				opts.Encryption = &ImageRepositoryEncryption{
					EncryptionType: ImageRepositoryEncryptionTypeSnowflakeFull,
				}
				opts.Comment = &createComment
				opts.Tag = []TagAssociation{{Name: tagId, Value: "value1"}}
			},
			"CREATE IMAGE REPOSITORY IF NOT EXISTS %s ENCRYPTION = (TYPE = 'SNOWFLAKE_FULL') COMMENT = '%s' TAG (%s = 'value1')",
			id.FullyQualifiedName(), createComment, tagId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateImageRepositoryOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE IMAGE REPOSITORY %s", id.FullyQualifiedName(),
		)

	imageRepositoriesTests.Alter.
		withModify(
			case_ImageRepositories_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterImageRepositoryOptions) {
				opts.Set = &ImageRepositorySet{}
				opts.UnsetTags = []ObjectIdentifier{tagId}
			},
		).
		withModifyAndExpectedSqlf(
			case_ImageRepositories_sql_Alter_Set,
			func(opts *AlterImageRepositoryOptions) {
				opts.Set = &ImageRepositorySet{Comment: &StringAllowEmpty{Value: "test"}}
			},
			"ALTER IMAGE REPOSITORY %s SET COMMENT = 'test'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ImageRepositories_sql_Alter_SetTags,
			func(opts *AlterImageRepositoryOptions) {
				opts.IfExists = new(true)
				opts.SetTags = []TagAssociation{
					{Name: tagId, Value: "value1"},
					{Name: tagId2, Value: "value2"},
				}
			},
			`ALTER IMAGE REPOSITORY IF EXISTS %s SET TAG %s = 'value1', %s = 'value2'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ImageRepositories_sql_Alter_UnsetTags,
			func(opts *AlterImageRepositoryOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			`ALTER IMAGE REPOSITORY %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		)

	imageRepositoriesTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_Set_emptyComment",
			func(opts *AlterImageRepositoryOptions) {
				opts.Set = &ImageRepositorySet{Comment: &StringAllowEmpty{Value: ""}}
			},
			"ALTER IMAGE REPOSITORY %s SET COMMENT = ''", id.FullyQualifiedName(),
		)

	imageRepositoriesTests.Drop.
		withExpectedSqlf(
			case_ImageRepositories_sql_Drop_basic,
			"DROP IMAGE REPOSITORY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ImageRepositories_sql_Drop_all,
			func(opts *DropImageRepositoryOptions) { opts.IfExists = new(true) },
			"DROP IMAGE REPOSITORY IF EXISTS %s", id.FullyQualifiedName(),
		)

	imageRepositoriesTests.Show.
		withExpectedSql(case_ImageRepositories_sql_Show_basic, "SHOW IMAGE REPOSITORIES").
		withModifyAndExpectedSqlf(
			case_ImageRepositories_sql_Show_all,
			func(opts *ShowImageRepositoryOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.In = &In{Schema: schemaId}
			},
			"SHOW IMAGE REPOSITORIES LIKE 'pattern' IN SCHEMA %s", schemaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ImageRepositories_sql_Show_Like,
			func(opts *ShowImageRepositoryOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			"SHOW IMAGE REPOSITORIES LIKE 'pattern'",
		).
		withModifyAndExpectedSqlf(
			case_ImageRepositories_sql_Show_In,
			func(opts *ShowImageRepositoryOptions) { opts.In = &In{Schema: schemaId} },
			"SHOW IMAGE REPOSITORIES IN SCHEMA %s", schemaId.FullyQualifiedName(),
		)
}
