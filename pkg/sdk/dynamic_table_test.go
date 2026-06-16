package sdk

import (
	"testing"
)

func TestDynamicTableCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *createDynamicTableOptions {
		return &createDynamicTableOptions{
			name: id,
			targetLag: TargetLag{
				MaximumDuration: new("1 minutes"),
			},
			warehouse: AccountObjectIdentifier{
				name: "warehouse_name",
			},
			query: "SELECT product_id, product_name FROM staging_table",
		}
	}
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createDynamicTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = new(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE DYNAMIC TABLE %s TARGET_LAG = '1 minutes' WAREHOUSE = "warehouse_name" AS SELECT product_id, product_name FROM staging_table`, id.FullyQualifiedName())
	})

	t.Run("all optional", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = new(true)
		opts.Comment = new("comment")
		opts.RefreshMode = DynamicTableRefreshModeFull.ToPointer()
		opts.Initialize = DynamicTableInitializeOnSchedule.ToPointer()
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE DYNAMIC TABLE %s TARGET_LAG = '1 minutes' INITIALIZE = ON_SCHEDULE REFRESH_MODE = FULL WAREHOUSE = "warehouse_name" COMMENT = 'comment' AS SELECT product_id, product_name FROM staging_table`, id.FullyQualifiedName())
	})
}

func TestDynamicTableAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *alterDynamicTableOptions {
		return &alterDynamicTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *alterDynamicTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("alterDynamicTableOptions", "Suspend", "Resume", "Refresh", "Set"))
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Resume = new(true)
		opts.Suspend = new(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("alterDynamicTableOptions", "Suspend", "Resume", "Refresh", "Set"))
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("alterDynamicTableOptions", "Suspend", "Resume", "Refresh", "Set"))
	})

	t.Run("suspend", func(t *testing.T) {
		opts := defaultOpts()
		opts.Suspend = new(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER DYNAMIC TABLE %s SUSPEND`, id.FullyQualifiedName())
	})

	t.Run("resume", func(t *testing.T) {
		opts := defaultOpts()
		opts.Resume = new(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER DYNAMIC TABLE %s RESUME`, id.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DynamicTableSet{
			TargetLag: &TargetLag{
				MaximumDuration: new("1 minutes"),
			},
			Warehouse: &AccountObjectIdentifier{
				name: "warehouse_name",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DYNAMIC TABLE %s SET TARGET_LAG = '1 minutes' WAREHOUSE = "warehouse_name"`, id.FullyQualifiedName())
	})
}

func TestDynamicTableDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *dropDynamicTableOptions {
		return &dropDynamicTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *dropDynamicTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP DYNAMIC TABLE %s`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = new(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP DYNAMIC TABLE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestDynamicTableShow(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *showDynamicTableOptions {
		return &showDynamicTableOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *showDynamicTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, ErrPatternRequiredForLikeKeyword)
	})

	t.Run("show with in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Database: NewAccountObjectIdentifier("database"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW DYNAMIC TABLES IN DATABASE "database"`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: new(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW DYNAMIC TABLES LIKE '%s'`, id.Name())
	})

	t.Run("show with like and in", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: new(id.Name()),
		}
		opts.In = &In{
			Database: NewAccountObjectIdentifier("database"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW DYNAMIC TABLES LIKE '%s' IN DATABASE "database"`, id.Name())
	})
}

func TestDynamicTableDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *describeDynamicTableOptions {
		return &describeDynamicTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *describeDynamicTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("describe", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE DYNAMIC TABLE %s`, id.FullyQualifiedName())
	})
}
