package sdk

import (
	"testing"
)

func TestPipesCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreatePipeOptions {
		return &CreatePipeOptions{
			name:          id,
			copyStatement: "<copy_statement>",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreatePipeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = new(true)
		opts.OrReplace = new(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreatePipeOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: copy statement required", func(t *testing.T) {
		opts := defaultOpts()
		opts.copyStatement = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreatePipeOptions", "copyStatement"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE PIPE %s AS <copy_statement>`, id.FullyQualifiedName())
	})

	t.Run("all optional", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = new(true)
		opts.AutoIngest = new(true)
		opts.ErrorIntegration = new("some_error_integration")
		opts.AwsSnsTopic = new("some aws sns topic")
		opts.Integration = new("some integration")
		opts.Comment = new("some comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE PIPE IF NOT EXISTS %s AUTO_INGEST = true ERROR_INTEGRATION = some_error_integration AWS_SNS_TOPIC = 'some aws sns topic' INTEGRATION = 'some integration' COMMENT = 'some comment' AS <copy_statement>`, id.FullyQualifiedName())
	})
}

func TestPipesAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *AlterPipeOptions {
		return &AlterPipeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterPipeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterPipeOptions", "Set", "Unset", "SetTags", "UnsetTags", "Refresh"))
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &PipeSet{
			ErrorIntegration: new("new_error_integration"),
		}
		opts.Unset = &PipeUnset{
			Comment: new(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterPipeOptions", "Set", "Unset", "SetTags", "UnsetTags", "Refresh"))
	})

	t.Run("validation: no property to set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &PipeSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterPipeOptions.Set", "ErrorIntegration", "PipeExecutionPaused", "Comment"))
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &PipeUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterPipeOptions.Unset", "ErrorIntegration", "PipeExecutionPaused", "Comment"))
	})

	t.Run("set tag: single", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag_name1"),
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s SET TAG "tag_name1" = 'v1'`, id.FullyQualifiedName())
	})

	t.Run("set tag: multiple", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag_name1"),
				Value: "v1",
			},
			{
				Name:  NewAccountObjectIdentifier("tag_name2"),
				Value: "v2",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s SET TAG "tag_name1" = 'v1', "tag_name2" = 'v2'`, id.FullyQualifiedName())
	})

	t.Run("set all", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = new(true)
		opts.Set = &PipeSet{
			ErrorIntegration:    new("new_error_integration"),
			PipeExecutionPaused: new(true),
			Comment:             new("new comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE IF EXISTS %s SET ERROR_INTEGRATION = new_error_integration, PIPE_EXECUTION_PAUSED = true, COMMENT = 'new comment'`, id.FullyQualifiedName())
	})

	t.Run("unset tag: single", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag_name1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s UNSET TAG "tag_name1"`, id.FullyQualifiedName())
	})

	t.Run("unset tag: multi", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag_name1"),
			NewAccountObjectIdentifier("tag_name2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s UNSET TAG "tag_name1", "tag_name2"`, id.FullyQualifiedName())
	})

	t.Run("unset all", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = new(true)
		opts.Unset = &PipeUnset{
			ErrorIntegration:    new(true),
			PipeExecutionPaused: new(true),
			Comment:             new(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE IF EXISTS %s UNSET ERROR_INTEGRATION, PIPE_EXECUTION_PAUSED, COMMENT`, id.FullyQualifiedName())
	})

	t.Run("refresh", func(t *testing.T) {
		opts := defaultOpts()
		opts.Refresh = &PipeRefresh{}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE %s REFRESH`, id.FullyQualifiedName())
	})

	t.Run("refresh with all", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = new(true)
		opts.Refresh = &PipeRefresh{
			Prefix:        new("/d1"),
			ModifiedAfter: new("2018-07-30T13:56:46-07:00"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER PIPE IF EXISTS %s REFRESH PREFIX = '/d1' MODIFIED_AFTER = '2018-07-30T13:56:46-07:00'`, id.FullyQualifiedName())
	})
}

func TestPipesDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DropPipeOptions {
		return &DropPipeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropPipeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP PIPE %s`, id.FullyQualifiedName())
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = new(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP PIPE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestPipesShow(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *ShowPipeOptions {
		return &ShowPipeOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowPipeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, ErrPatternRequiredForLikeKeyword)
	})

	t.Run("validation: empty in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ShowPipeOptions.In", "Account", "Database", "Schema"))
	})

	t.Run("validation: exactly one scope for in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Account:  new(true),
			Database: id.DatabaseId(),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ShowPipeOptions.In", "Account", "Database", "Schema"))
	})

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES`)
	})

	t.Run("with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: new(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES LIKE '%s'`, id.Name())
	})

	t.Run("in account", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Account: new(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES IN ACCOUNT`)
	})

	t.Run("in database", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Database: id.DatabaseId(),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES IN DATABASE %s`, id.DatabaseId().FullyQualifiedName())
	})

	t.Run("in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Schema: id.SchemaId(),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES IN SCHEMA %s`, id.SchemaId().FullyQualifiedName())
	})

	t.Run("with like and in account", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: new(id.Name()),
		}
		opts.In = &In{
			Account: new(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES LIKE '%s' IN ACCOUNT`, id.Name())
	})

	t.Run("with like and in database", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: new(id.Name()),
		}
		opts.In = &In{
			Database: id.DatabaseId(),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES LIKE '%s' IN DATABASE %s`, id.Name(), id.DatabaseId().FullyQualifiedName())
	})

	t.Run("with like and in schema", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: new(id.Name()),
		}
		opts.In = &In{
			Schema: id.SchemaId(),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW PIPES LIKE '%s' IN SCHEMA %s`, id.Name(), id.SchemaId().FullyQualifiedName())
	})
}

func TestPipesDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DescribePipeOptions {
		return &DescribePipeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribePipeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("with name", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE PIPE %s`, id.FullyQualifiedName())
	})
}
