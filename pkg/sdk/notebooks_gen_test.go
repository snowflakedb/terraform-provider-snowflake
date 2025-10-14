package sdk

import (
	"testing"
)

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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set] should be present", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterNotebookOptions", "Set"))
	})

	t.Run("validation: valid identifier for [opts.Set.QueryWarehouse] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.Set.Warehouse] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.Set.ComputePool] if set", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
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
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		// TODO: fill me
		assertOptsValidAndSQLEquals(t, opts, "TODO: fill me")
	})
}
