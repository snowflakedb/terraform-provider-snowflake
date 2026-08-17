package sdk

func init() {
	id := gitRepositoriesTestIdSchemaObjectIdentifier
	apiIntegrationId := randomAccountObjectIdentifier()
	gitCredentialsId := randomSchemaObjectIdentifier()
	databaseId := NewAccountObjectIdentifier("database-name")
	tagId := NewAccountObjectIdentifier("tag-name")
	origin := "https://github.com/user/repo"

	gitRepositoriesTests.Create.
		withDefaultOpts(func() *CreateGitRepositoryOptions {
			return &CreateGitRepositoryOptions{
				name:           id,
				Origin:         origin,
				ApiIntegration: apiIntegrationId,
			}
		}).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Create_basic,
			func(opts *CreateGitRepositoryOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE GIT REPOSITORY %s ORIGIN = '%s' API_INTEGRATION = %s",
			id.FullyQualifiedName(), origin, apiIntegrationId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Create_all,
			func(opts *CreateGitRepositoryOptions) {
				opts.IfNotExists = new(true)
				opts.GitCredentials = &gitCredentialsId
				opts.Comment = new("comment")
				opts.Tag = []TagAssociation{{Name: tagId, Value: "tag-value"}}
			},
			`CREATE GIT REPOSITORY IF NOT EXISTS %s ORIGIN = '%s' API_INTEGRATION = %s GIT_CREDENTIALS = %s COMMENT = 'comment' TAG (%s = 'tag-value')`,
			id.FullyQualifiedName(), origin, apiIntegrationId.FullyQualifiedName(), gitCredentialsId.FullyQualifiedName(), tagId.FullyQualifiedName(),
		)

	gitRepositoriesTests.Alter.
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Alter_Set,
			func(opts *AlterGitRepositoryOptions) {
				opts.Set = &GitRepositorySet{
					ApiIntegration: &apiIntegrationId,
					GitCredentials: &gitCredentialsId,
					Comment:        new("comment"),
				}
			},
			"ALTER GIT REPOSITORY %s SET API_INTEGRATION = %s GIT_CREDENTIALS = %s COMMENT = 'comment'",
			id.FullyQualifiedName(), apiIntegrationId.FullyQualifiedName(), gitCredentialsId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Alter_Unset,
			func(opts *AlterGitRepositoryOptions) {
				opts.Unset = &GitRepositoryUnset{
					GitCredentials: new(true),
					Comment:        new(true),
				}
			},
			"ALTER GIT REPOSITORY %s UNSET GIT_CREDENTIALS, COMMENT", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Alter_SetTags,
			func(opts *AlterGitRepositoryOptions) {
				opts.SetTags = []TagAssociation{{Name: tagId, Value: "tag-value"}}
			},
			`ALTER GIT REPOSITORY %s SET TAG %s = 'tag-value'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Alter_UnsetTags,
			func(opts *AlterGitRepositoryOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId}
			},
			`ALTER GIT REPOSITORY %s UNSET TAG %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Alter_Fetch,
			func(opts *AlterGitRepositoryOptions) { opts.Fetch = new(true) },
			"ALTER GIT REPOSITORY %s FETCH", id.FullyQualifiedName(),
		)

	gitRepositoriesTests.Drop.
		withExpectedSqlf(
			case_GitRepositories_sql_Drop_basic,
			"DROP GIT REPOSITORY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Drop_all,
			func(opts *DropGitRepositoryOptions) { opts.IfExists = new(true) },
			"DROP GIT REPOSITORY IF EXISTS %s", id.FullyQualifiedName(),
		)

	gitRepositoriesTests.Describe.
		withExpectedSqlf(
			case_GitRepositories_sql_Describe_basic,
			"DESCRIBE GIT REPOSITORY %s", id.FullyQualifiedName(),
		)

	gitRepositoriesTests.Show.
		withExpectedSql(case_GitRepositories_sql_Show_basic, "SHOW GIT REPOSITORIES").
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Show_all,
			func(opts *ShowGitRepositoryOptions) {
				opts.Like = &Like{Pattern: new("git-repository-name")}
				opts.In = &In{Database: databaseId}
				opts.Limit = &LimitFrom{Rows: new(10)}
			},
			`SHOW GIT REPOSITORIES LIKE 'git-repository-name' IN DATABASE %s LIMIT 10`,
			databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Show_Like,
			func(opts *ShowGitRepositoryOptions) {
				opts.Like = &Like{Pattern: new("git-repository-name")}
			},
			`SHOW GIT REPOSITORIES LIKE 'git-repository-name'`,
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Show_In,
			func(opts *ShowGitRepositoryOptions) { opts.In = &In{Database: databaseId} },
			`SHOW GIT REPOSITORIES IN DATABASE %s`, databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_Show_Limit,
			func(opts *ShowGitRepositoryOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			`SHOW GIT REPOSITORIES LIMIT 10`,
		)

	gitRepositoriesTests.ShowGitBranches.
		withExpectedSqlf(
			case_GitRepositories_sql_ShowGitBranches_basic,
			"SHOW GIT BRANCHES IN %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_ShowGitBranches_all,
			func(opts *ShowGitBranchesGitRepositoryOptions) {
				opts.Like = &Like{Pattern: new("branch-name")}
				opts.GitRepository = new(true)
			},
			"SHOW GIT BRANCHES LIKE 'branch-name' IN GIT REPOSITORY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_ShowGitBranches_Like,
			func(opts *ShowGitBranchesGitRepositoryOptions) {
				opts.Like = &Like{Pattern: new("branch-name")}
			},
			"SHOW GIT BRANCHES LIKE 'branch-name' IN %s", id.FullyQualifiedName(),
		)

	gitRepositoriesTests.ShowGitTags.
		withExpectedSqlf(
			case_GitRepositories_sql_ShowGitTags_basic,
			"SHOW GIT TAGS IN %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_ShowGitTags_all,
			func(opts *ShowGitTagsGitRepositoryOptions) {
				opts.Like = &Like{Pattern: new("tag-name")}
				opts.GitRepository = new(true)
			},
			"SHOW GIT TAGS LIKE 'tag-name' IN GIT REPOSITORY %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_GitRepositories_sql_ShowGitTags_Like,
			func(opts *ShowGitTagsGitRepositoryOptions) {
				opts.Like = &Like{Pattern: new("tag-name")}
			},
			"SHOW GIT TAGS LIKE 'tag-name' IN %s", id.FullyQualifiedName(),
		)
}
