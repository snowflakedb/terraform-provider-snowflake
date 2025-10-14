package sdk

import "testing"

func TestNotebooks_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	stageId := randomSchemaObjectIdentifier()
	var stageLocation Location = &StageLocation{
		stage: stageId,
		path:  "dir/subdir",
	}

	// Minimal valid CreateNotebookOptions
	defaultOpts := func() *CreateNotebookOptions {
		return &CreateNotebookOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateNotebookOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.QueryWarehouse] if set", func(t *testing.T) {
		opts := defaultOpts()
		opts.QueryWarehouse = &emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.Warehouse] if set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Warehouse = &emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateNotebookOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: valid identifier for [opts.ComputePool] if set", func(t *testing.T) {
		opts := defaultOpts()
		opts.ComputePool = &emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE NOTEBOOK %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()

		opts.IfNotExists = Bool(true)
		opts.From = &stageLocation
		opts.MainFile = String("main_file")
		opts.Comment = String("comment")
		opts.QueryWarehouse = &AccountObjectIdentifier{"sample_qwh"}
		opts.IdleAutoShutdownTimeSeconds = Int(3600)
		opts.Warehouse = &AccountObjectIdentifier{"sample_wh"}
		opts.RuntimeName = String("sample")
		opts.ComputePool = &AccountObjectIdentifier{"sample_cp"}
		opts.ExternalAccessIntegrations = []AccountObjectIdentifier{}
		opts.RuntimeEnvironmentVersion = String("WH-RUNTIME-2.0")
		opts.DefaultVersion = String("FIRST")

		assertOptsValidAndSQLEquals(t, opts, "CREATE NOTEBOOK IF NOT EXISTS %s FROM @%s/dir/subdir MAIN_FILE = 'main_file' COMMENT = 'comment' QUERY_WAREHOUSE = \"sample_qwh\" IDLE_AUTO_SHUTDOWN_TIME_SECONDS = 3600 WAREHOUSE = \"sample_wh\" RUNTIME_NAME = 'sample' COMPUTE_POOL = \"sample_cp\" RUNTIME_ENVIRONMENT_VERSION = 'WH-RUNTIME-2.0' DEFAULT_VERSION = FIRST", id.FullyQualifiedName(), stageId.FullyQualifiedName())
	})
}

func TestNotebooks_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid AlterNotebookOptions
	defaultOpts := func() *AlterNotebookOptions {
		return &AlterNotebookOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterNotebookOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterNotebookOptions", "Set", "Unset"))

		opts.Set = &NotebookSet{
			Comment: String("comment"),
		}

		opts.Unset = &NotebookUnset{
			Comment: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterNotebookOptions", "Set", "Unset"))
	})

	t.Run("validation: valid identifier for [opts.Set.QueryWarehouse] if set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NotebookSet{
			QueryWarehouse: &emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.Set.Warehouse] if set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NotebookSet{
			Warehouse: &emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.Set.ComputePool] if set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NotebookSet{
			ComputePool: &emptyAccountObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NotebookSet{
			Comment: String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER NOTEBOOK %s SET COMMENT = 'comment'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &NotebookSet{
			Comment:                     String("comment"),
			QueryWarehouse:              &AccountObjectIdentifier{"sample_qwh"},
			IdleAutoShutdownTimeSeconds: Int(3600),
			SecretsList:                 &SecretsList{[]SecretReference{{"var_name", true, SchemaObjectIdentifier{}}}},
			MainFile:                    String("main_file"),
			Warehouse:                   &AccountObjectIdentifier{"sample_wh"},
			RuntimeName:                 String("runtime_name"),
			ComputePool:                 &AccountObjectIdentifier{"sample_cp"},
			ExternalAccessIntegrations:  []AccountObjectIdentifier{{"test"}},
			RuntimeEnvironmentVersion:   String("WH-RUNTIME-2.0"),
		}

		assertOptsValidAndSQLEquals(t, opts, "ALTER NOTEBOOK %s SET COMMENT = 'comment' QUERY_WAREHOUSE = \"sample_qwh\" IDLE_AUTO_SHUTDOWN_TIME_SECONDS = 3600 SECRETS = ('var_name' =) MAIN_FILE = 'main_file' WAREHOUSE = \"sample_wh\" RUNTIME_NAME = 'runtime_name' COMPUTE_POOL = \"sample_cp\" EXTERNAL_ACCESS_INTEGRATIONS = (\"test\") RUNTIME_ENVIRONMENT_VERSION = 'WH-RUNTIME-2.0'", id.FullyQualifiedName())
	})
}

func TestNotebooks_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid DropNotebookOptions
	defaultOpts := func() *DropNotebookOptions {
		return &DropNotebookOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropNotebookOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP NOTEBOOK %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP NOTEBOOK IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestNotebooks_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid DescribeNotebookOptions
	defaultOpts := func() *DescribeNotebookOptions {
		return &DescribeNotebookOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeNotebookOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE NOTEBOOK %s", id.FullyQualifiedName())
	})
}

func TestNotebooks_Show(t *testing.T) {
	// Minimal valid ShowNotebookOptions
	defaultOpts := func() *ShowNotebookOptions {
		return &ShowNotebookOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowNotebookOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW NOTEBOOKS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("notebook-name"),
		}
		opts.In = &In{
			Database: NewAccountObjectIdentifier("database-name"),
		}
		opts.Limit = &LimitFrom{
			Rows: Int(10),
		}
		opts.StartsWith = String("prefix")
		assertOptsValidAndSQLEquals(t, opts, "SHOW NOTEBOOKS LIKE 'notebook-name' IN DATABASE \"database-name\" LIMIT 10 STARTS WITH 'prefix'")
	})
}
