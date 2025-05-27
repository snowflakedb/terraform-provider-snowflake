package sdk

import "testing"

const (
	origin = "https://github.com/user/repo"
)

func TestGitRepositories_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	apiIntegrationId := randomAccountObjectIdentifier()

	// Minimal valid CreateGitRepositoryOptions
	defaultOpts := func() *CreateGitRepositoryOptions {
		return &CreateGitRepositoryOptions{
			name:           id,
			Origin:         origin,
			ApiIntegration: apiIntegrationId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateGitRepositoryOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateGitRepositoryOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE GIT REPOSITORY %s ORIGIN = '%s' API_INTEGRATION = %s", id.FullyQualifiedName(), origin, apiIntegrationId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		gitCredentialsId := randomAccountObjectIdentifier()

		opts.IfNotExists = Bool(true)
		opts.Origin = origin
		opts.GitCredentials = &gitCredentialsId
		opts.Comment = String("comment")
		opts.Tag = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag-name"),
				Value: "tag-value",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE GIT REPOSITORY IF NOT EXISTS %s ORIGIN = '%s' API_INTEGRATION = %s GIT_CREDENTIALS = %s COMMENT = '%s' TAG (\"tag-name\" = 'tag-value')", id.FullyQualifiedName(), origin, apiIntegrationId.FullyQualifiedName(), gitCredentialsId.FullyQualifiedName(), "comment")
	})
}

func TestGitRepositories_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid AlterGitRepositoryOptions
	defaultOpts := func() *AlterGitRepositoryOptions {
		return &AlterGitRepositoryOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterGitRepositoryOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags opts.Fetch] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &GitRepositorySet{
			Comment: String("comment"),
		}

		opts.Unset = &GitRepositoryUnset{
			Comment: Bool(true),
		}

		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterGitRepositoryOptions", "Set", "Unset", "SetTags", "UnsetTags", "Fetch"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		apiIntegrationId := randomAccountObjectIdentifier()
		gitCredentialsId := randomAccountObjectIdentifier()

		opts.Set = &GitRepositorySet{
			ApiIntegration: &apiIntegrationId,
			GitCredentials: &gitCredentialsId,
			Comment:        String("comment"),
		}

		assertOptsValidAndSQLEquals(t, opts, "ALTER GIT REPOSITORY %s SET API_INTEGRATION = %s GIT_CREDENTIALS = %s COMMENT = 'comment'", id.FullyQualifiedName(), apiIntegrationId.FullyQualifiedName(), gitCredentialsId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()

		apiIntegrationId := randomAccountObjectIdentifier()
		gitCredentialsId := randomAccountObjectIdentifier()

		opts.Set = &GitRepositorySet{
			ApiIntegration: &apiIntegrationId,
			GitCredentials: &gitCredentialsId,
			Comment:        String("comment"),
		}

		assertOptsValidAndSQLEquals(t, opts, "ALTER GIT REPOSITORY %s SET API_INTEGRATION = %s GIT_CREDENTIALS = %s COMMENT = 'comment'", id.FullyQualifiedName(), apiIntegrationId.FullyQualifiedName(), gitCredentialsId.FullyQualifiedName())

		opts.Set = nil
		opts.Unset = &GitRepositoryUnset{
			GitCredentials: Bool(true),
			Comment:        Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER GIT REPOSITORY %s UNSET GIT_CREDENTIALS, COMMENT", id.FullyQualifiedName())

		opts.Unset = nil
		tag := []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag-name"),
				Value: "tag-value",
			},
		}
		opts.SetTags = tag
		assertOptsValidAndSQLEquals(t, opts, "ALTER GIT REPOSITORY %s SET TAG \"tag-name\" = 'tag-value'", id.FullyQualifiedName())

		opts.SetTags = nil
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag-name"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER GIT REPOSITORY %s UNSET TAG \"tag-name\"", id.FullyQualifiedName())

		opts.UnsetTags = nil
		opts.Fetch = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER GIT REPOSITORY %s FETCH", id.FullyQualifiedName())
	})
}

func TestGitRepositories_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid DropGitRepositoryOptions
	defaultOpts := func() *DropGitRepositoryOptions {
		return &DropGitRepositoryOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropGitRepositoryOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP GIT REPOSITORY %s`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP GIT REPOSITORY IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestGitRepositories_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid DescribeGitRepositoryOptions
	defaultOpts := func() *DescribeGitRepositoryOptions {
		return &DescribeGitRepositoryOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeGitRepositoryOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE GIT REPOSITORY %s", id.FullyQualifiedName())
	})
}

func TestGitRepositories_Show(t *testing.T) {
	// Minimal valid ShowGitRepositoryOptions
	defaultOpts := func() *ShowGitRepositoryOptions {
		return &ShowGitRepositoryOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowGitRepositoryOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW GIT REPOSITORIES")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("git-repository-name"),
		}
		opts.In = &In{
			Database: NewAccountObjectIdentifier("database-name"),
		}
		opts.Limit = &LimitFrom{
			Rows: Int(10),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW GIT REPOSITORIES LIKE 'git-repository-name' IN DATABASE \"database-name\" LIMIT 10")
	})
}

func TestGitRepositories_ShowGitBranches(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid ShowGitBranchesGitRepositoryOptions
	defaultOpts := func() *ShowGitBranchesGitRepositoryOptions {
		return &ShowGitBranchesGitRepositoryOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowGitBranchesGitRepositoryOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW GIT BRANCHES IN %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("branch-name"),
		}
		opts.GitRepository = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "SHOW GIT BRANCHES LIKE 'branch-name' IN GIT REPOSITORY %s", id.FullyQualifiedName())
	})
}

func TestGitRepositories_ShowGitTags(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid ShowGitTagsGitRepositoryOptions
	defaultOpts := func() *ShowGitTagsGitRepositoryOptions {
		return &ShowGitTagsGitRepositoryOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowGitTagsGitRepositoryOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW GIT TAGS IN %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("tag-name"),
		}
		opts.GitRepository = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "SHOW GIT TAGS LIKE 'tag-name' IN GIT REPOSITORY %s", id.FullyQualifiedName())
	})
}
