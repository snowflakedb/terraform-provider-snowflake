package sdk

import "testing"

func TestCortexSearchServices_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid CreateCortexSearchServiceOptions
	defaultOpts := func() *CreateCortexSearchServiceOptions {
		return &CreateCortexSearchServiceOptions{
			name:      id,
			On:        "searchable_text",
			TargetLag: "1 minutes",
			Warehouse: AccountObjectIdentifier{
				name: "warehouse_name",
			},
			QueryDefinition: "SELECT product_id, product_name, searchable_text FROM staging_table",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateCortexSearchServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: [opts.On] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.On = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateCortexSearchServiceOptions", "On"))
	})

	t.Run("validation: [opts.TargetLag] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.TargetLag = ""
		assertOptsInvalidJoinedErrors(t, opts, errNotSet("CreateCortexSearchServiceOptions", "TargetLag"))
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateCortexSearchServiceOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE CORTEX SEARCH SERVICE %s ON searchable_text WAREHOUSE = "warehouse_name" TARGET_LAG = '1 minutes' AS SELECT product_id, product_name, searchable_text FROM staging_table`, id.FullyQualifiedName())
	})

	t.Run("with embedding model", func(t *testing.T) {
		opts := defaultOpts()
		opts.EmbeddingModel = String("snowflake-arctic-embed-m-v1.5")
		assertOptsValidAndSQLEquals(t, opts, `CREATE CORTEX SEARCH SERVICE %s ON searchable_text WAREHOUSE = "warehouse_name" TARGET_LAG = '1 minutes' EMBEDDING_MODEL = 'snowflake-arctic-embed-m-v1.5' AS SELECT product_id, product_name, searchable_text FROM staging_table`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.Attributes = &Attributes{
			Columns: []string{"product_id", "product_name"},
		}
		opts.EmbeddingModel = String("snowflake-arctic-embed-l-v2.0")
		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE CORTEX SEARCH SERVICE IF NOT EXISTS %s ON searchable_text ATTRIBUTES product_id, product_name WAREHOUSE = "warehouse_name" TARGET_LAG = '1 minutes' EMBEDDING_MODEL = 'snowflake-arctic-embed-l-v2.0' COMMENT = 'comment' AS SELECT product_id, product_name, searchable_text FROM staging_table`, id.FullyQualifiedName())
	})
}

func TestCortexSearchServices_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid AlterCortexSearchServiceOptions
	defaultOpts := func() *AlterCortexSearchServiceOptions {
		return &AlterCortexSearchServiceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterCortexSearchServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterCortexSearchServiceOptions", "Set"))
	})

	t.Run("validation: at least one of the fields [opts.Set.TargetLag opts.Set.Warehouse opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &CortexSearchServiceSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterCortexSearchServiceOptions.Set", "TargetLag", "Warehouse", "Comment"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &CortexSearchServiceSet{
			TargetLag: String("1 minutes"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER CORTEX SEARCH SERVICE %s SET TARGET_LAG = '1 minutes'`, id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &CortexSearchServiceSet{
			TargetLag: String("1 minutes"),
			Warehouse: &AccountObjectIdentifier{
				name: "warehouse_name",
			},
			Comment: String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER CORTEX SEARCH SERVICE %s SET TARGET_LAG = '1 minutes' WAREHOUSE = "warehouse_name" COMMENT = 'comment'`, id.FullyQualifiedName())
	})
}

func TestCortexSearchServices_Show(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid ShowCortexSearchServiceOptions
	defaultOpts := func() *ShowCortexSearchServiceOptions {
		return &ShowCortexSearchServiceOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowCortexSearchServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW CORTEX SEARCH SERVICES")
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW CORTEX SEARCH SERVICES LIKE '%s'`, id.Name())
	})

	t.Run("show with in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &In{
			Database: NewAccountObjectIdentifier("database"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW CORTEX SEARCH SERVICES IN DATABASE "database"`)
	})

	t.Run("show with starts with", func(t *testing.T) {
		opts := defaultOpts()
		opts.StartsWith = String("foo")
		assertOptsValidAndSQLEquals(t, opts, `SHOW CORTEX SEARCH SERVICES STARTS WITH 'foo'`)
	})

	t.Run("show with limit", func(t *testing.T) {
		opts := defaultOpts()
		opts.Limit = &LimitFrom{
			Rows: Int(1),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW CORTEX SEARCH SERVICES LIMIT 1`)
	})

	t.Run("show with limit from", func(t *testing.T) {
		opts := defaultOpts()
		opts.Limit = &LimitFrom{
			Rows: Int(1),
			From: String("foo"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW CORTEX SEARCH SERVICES LIMIT 1 FROM 'foo'`)
	})

	t.Run("show with all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		opts.In = &In{
			Account: Bool(true),
		}
		opts.StartsWith = String("foo")
		opts.Limit = &LimitFrom{
			Rows: Int(1),
			From: String("bar"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW CORTEX SEARCH SERVICES LIKE 'pattern' IN ACCOUNT STARTS WITH 'foo' LIMIT 1 FROM 'bar'`)
	})
}

func TestCortexSearchServices_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid DescribeCortexSearchServiceOptions
	defaultOpts := func() *DescribeCortexSearchServiceOptions {
		return &DescribeCortexSearchServiceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeCortexSearchServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE CORTEX SEARCH SERVICE %s`, id.FullyQualifiedName())
	})
}

func TestCortexSearchServices_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// Minimal valid DropCortexSearchServiceOptions
	defaultOpts := func() *DropCortexSearchServiceOptions {
		return &DropCortexSearchServiceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropCortexSearchServiceOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP CORTEX SEARCH SERVICE %s`, id.FullyQualifiedName())
	})
}
