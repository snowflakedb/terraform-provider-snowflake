package sdk

import "testing"

func TestSemanticViews_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateSemanticViewOptions {
		return &CreateSemanticViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateSemanticViewOptions = nil
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
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSemanticViewOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SEMANTIC VIEW %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE SEMANTIC VIEW IF NOT EXISTS %s COMMENT = '%s'`, id.FullyQualifiedName(), "comment")
	})
}

func TestSemanticViews_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DropSemanticViewOptions {
		return &DropSemanticViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropSemanticViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP SEMANTIC VIEW %s`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP SEMANTIC VIEW IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestSemanticViews_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *DescribeSemanticViewOptions {
		return &DescribeSemanticViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeSemanticViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE SEMANTIC VIEW %s", id.FullyQualifiedName())
	})
}

func TestSemanticViews_Show(t *testing.T) {
	defaultOpts := func() *ShowSemanticViewsOptions {
		return &ShowSemanticViewsOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowSemanticViewsOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW SEMANTIC VIEWS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		opts.Like = &Like{
			Pattern: String("semantic-view-name"),
		}
		opts.In = &In{
			Database: NewAccountObjectIdentifier("database-name"),
		}
		opts.StartsWith = String("semantic-view-name-start")
		opts.Limit = &LimitFrom{
			Rows: Int(10),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SEMANTIC VIEWS LIKE 'semantic-view-name' IN DATABASE "database-name" STARTS WITH 'semantic-view-name-start' LIMIT 10`)
	})
}
