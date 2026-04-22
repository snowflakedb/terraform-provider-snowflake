package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenflowRuntimes_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	deploymentID := randomAccountObjectIdentifier()

	defaultOpts := func() *CreateOpenflowRuntimeOptions {
		return &CreateOpenflowRuntimeOptions{
			name:          id,
			InDeployment:  deploymentID,
			ExecuteAsRole: "MY_ROLE",
			NodeType:      OpenflowRuntimeNodeTypeSmall,
			MinNodes:      1,
			MaxNodes:      3,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateOpenflowRuntimeOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: MinNodes > 0", func(t *testing.T) {
		opts := defaultOpts()
		opts.MinNodes = 0
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("CreateOpenflowRuntimeOptions", "MinNodes", IntErrGreater, 0))
	})

	t.Run("validation: MaxNodes >= MinNodes", func(t *testing.T) {
		opts := defaultOpts()
		opts.MinNodes = 3
		opts.MaxNodes = 1
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("CreateOpenflowRuntimeOptions", "MaxNodes", IntErrGreaterOrEqual, 3))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW RUNTIME %s IN DEPLOYMENT %s EXECUTE_AS_ROLE = MY_ROLE NODE_TYPE = 'SMALL' MIN_NODES = 1 MAX_NODES = 3",
			id.FullyQualifiedName(), deploymentID.FullyQualifiedName())
	})

	t.Run("if not exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW RUNTIME IF NOT EXISTS %s IN DEPLOYMENT %s EXECUTE_AS_ROLE = MY_ROLE NODE_TYPE = 'SMALL' MIN_NODES = 1 MAX_NODES = 3",
			id.FullyQualifiedName(), deploymentID.FullyQualifiedName())
	})

	t.Run("with comment and display name", func(t *testing.T) {
		comment := "my runtime"
		displayName := "My Runtime"
		opts := defaultOpts()
		opts.Comment = &comment
		opts.DisplayName = &displayName
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW RUNTIME %s IN DEPLOYMENT %s EXECUTE_AS_ROLE = MY_ROLE NODE_TYPE = 'SMALL' MIN_NODES = 1 MAX_NODES = 3 DISPLAY_NAME = 'My Runtime' COMMENT = 'my runtime'",
			id.FullyQualifiedName(), deploymentID.FullyQualifiedName())
	})
}

func TestOpenflowRuntimes_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *AlterOpenflowRuntimeOptions {
		return &AlterOpenflowRuntimeOptions{name: id}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterOpenflowRuntimeOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOpenflowRuntimeOptions", "Suspend", "Resume", "Terminate", "TerminateCascade", "Set", "Unset"))
	})

	t.Run("suspend", func(t *testing.T) {
		opts := defaultOpts()
		opts.Suspend = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s SUSPEND", id.FullyQualifiedName())
	})

	t.Run("resume", func(t *testing.T) {
		opts := defaultOpts()
		opts.Resume = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s RESUME", id.FullyQualifiedName())
	})

	t.Run("terminate", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terminate = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s TERMINATE", id.FullyQualifiedName())
	})

	t.Run("terminate cascade", func(t *testing.T) {
		opts := defaultOpts()
		opts.TerminateCascade = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s TERMINATE CASCADE", id.FullyQualifiedName())
	})

	t.Run("set min/max nodes", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OpenflowRuntimeSet{MinNodes: Int(2), MaxNodes: Int(5)}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s SET MIN_NODES = 2 MAX_NODES = 5", id.FullyQualifiedName())
	})

	t.Run("set comment", func(t *testing.T) {
		comment := "updated"
		opts := defaultOpts()
		opts.Set = &OpenflowRuntimeSet{Comment: &comment}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s SET COMMENT = 'updated'", id.FullyQualifiedName())
	})

	t.Run("unset display name", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OpenflowRuntimeUnset{DisplayName: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s UNSET DISPLAY_NAME", id.FullyQualifiedName())
	})
}

func TestOpenflowRuntimes_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropOpenflowRuntimeOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := &DropOpenflowRuntimeOptions{name: id}
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW RUNTIME %s", id.FullyQualifiedName())
	})

	t.Run("if exists", func(t *testing.T) {
		opts := &DropOpenflowRuntimeOptions{name: id, IfExists: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW RUNTIME IF EXISTS %s", id.FullyQualifiedName())
	})

	t.Run("cascade", func(t *testing.T) {
		opts := &DropOpenflowRuntimeOptions{name: id, Cascade: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW RUNTIME %s CASCADE", id.FullyQualifiedName())
	})
}

func TestOpenflowRuntimes_Show(t *testing.T) {
	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowOpenflowRuntimeOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := &ShowOpenflowRuntimeOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW RUNTIMES")
	})

	t.Run("like", func(t *testing.T) {
		opts := &ShowOpenflowRuntimeOptions{Like: &Like{Pattern: String("rt%")}}
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW RUNTIMES LIKE 'rt%%'")
	})
}

func Test_ToOpenflowRuntimeNodeType(t *testing.T) {
	valid := []struct {
		input string
		want  OpenflowRuntimeNodeType
	}{
		{"SMALL", OpenflowRuntimeNodeTypeSmall},
		{"small", OpenflowRuntimeNodeTypeSmall},
		{"MEDIUM", OpenflowRuntimeNodeTypeMedium},
		{"LARGE", OpenflowRuntimeNodeTypeLarge},
	}
	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToOpenflowRuntimeNodeType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
	_, err := ToOpenflowRuntimeNodeType("XLARGE")
	require.Error(t, err)
}

func Test_ToOpenflowRuntimeStatus(t *testing.T) {
	valid := []string{"CREATING", "CREATE_FAILED", "UPDATING", "UPDATE_FAILED", "SUSPENDING", "SUSPENDED", "SUSPEND_FAILED", "ACTIVATING", "ACTIVE", "ACTIVATE_FAILED", "DELETING", "DELETED", "DELETE_FAILED", "CANCEL_REQUESTED", "RESTARTING", "RESTART_FAILED", "UPGRADING", "UPGRADE_FAILED", "GENERATING_DIAGNOSTIC_BUNDLE", "CLEANING_UP", "INACTIVE"}
	for _, s := range valid {
		t.Run(s, func(t *testing.T) {
			got, err := ToOpenflowRuntimeStatus(s)
			require.NoError(t, err)
			require.Equal(t, OpenflowRuntimeStatus(s), got)
		})
	}
	_, err := ToOpenflowRuntimeStatus("UNKNOWN")
	require.Error(t, err)
}
