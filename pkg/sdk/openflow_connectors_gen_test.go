package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenflowConnectors_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	runtimeID := randomSchemaObjectIdentifier()

	defaultOpts := func() *CreateOpenflowConnectorOptions {
		return &CreateOpenflowConnectorOptions{
			name:           id,
			InRuntime:      runtimeID,
			FromDefinition: String("MY_CONNECTOR_DEF"),
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateOpenflowConnectorOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: requires from_definition or from_stage", func(t *testing.T) {
		opts := &CreateOpenflowConnectorOptions{name: id, InRuntime: runtimeID}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("CreateOpenflowConnectorOptions", "FromDefinition", "FromStage"))
	})

	t.Run("validation: cannot set both from_definition and from_stage", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromStage = String("@MY_STAGE/path/")
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("CreateOpenflowConnectorOptions", "FromDefinition", "FromStage"))
	})

	t.Run("from definition", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW CONNECTOR %s IN RUNTIME %s FROM DEFINITION MY_CONNECTOR_DEF",
			id.FullyQualifiedName(), runtimeID.FullyQualifiedName())
	})

	t.Run("from stage", func(t *testing.T) {
		opts := &CreateOpenflowConnectorOptions{
			name:      id,
			InRuntime: runtimeID,
			FromStage: String("@MY_REPO/branches/main/connector/"),
		}
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW CONNECTOR %s IN RUNTIME %s FROM '@MY_REPO/branches/main/connector/'",
			id.FullyQualifiedName(), runtimeID.FullyQualifiedName())
	})

	t.Run("if not exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW CONNECTOR IF NOT EXISTS %s IN RUNTIME %s FROM DEFINITION MY_CONNECTOR_DEF",
			id.FullyQualifiedName(), runtimeID.FullyQualifiedName())
	})

	t.Run("with comment", func(t *testing.T) {
		comment := "my connector"
		opts := defaultOpts()
		opts.Comment = &comment
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW CONNECTOR %s IN RUNTIME %s FROM DEFINITION MY_CONNECTOR_DEF COMMENT = 'my connector'",
			id.FullyQualifiedName(), runtimeID.FullyQualifiedName())
	})
}

func TestOpenflowConnectors_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	defaultOpts := func() *AlterOpenflowConnectorOptions {
		return &AlterOpenflowConnectorOptions{name: id}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterOpenflowConnectorOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOpenflowConnectorOptions", "Start", "Stop", "Set", "Unset"))
	})

	t.Run("start", func(t *testing.T) {
		opts := defaultOpts()
		opts.Start = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW CONNECTOR %s START", id.FullyQualifiedName())
	})

	t.Run("stop", func(t *testing.T) {
		opts := defaultOpts()
		opts.Stop = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW CONNECTOR %s STOP", id.FullyQualifiedName())
	})

	t.Run("set comment", func(t *testing.T) {
		comment := "updated"
		opts := defaultOpts()
		opts.Set = &OpenflowConnectorSet{Comment: &comment}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW CONNECTOR %s SET COMMENT = 'updated'", id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OpenflowConnectorUnset{Comment: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW CONNECTOR %s UNSET COMMENT", id.FullyQualifiedName())
	})
}

func TestOpenflowConnectors_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropOpenflowConnectorOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := &DropOpenflowConnectorOptions{name: id}
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW CONNECTOR %s", id.FullyQualifiedName())
	})

	t.Run("if exists", func(t *testing.T) {
		opts := &DropOpenflowConnectorOptions{name: id, IfExists: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW CONNECTOR IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestOpenflowConnectors_Show(t *testing.T) {
	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowOpenflowConnectorOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := &ShowOpenflowConnectorOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW CONNECTORS")
	})

	t.Run("like", func(t *testing.T) {
		opts := &ShowOpenflowConnectorOptions{Like: &Like{Pattern: String("cdc%")}}
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW CONNECTORS LIKE 'cdc%%'")
	})
}

func Test_ToOpenflowConnectorStatus(t *testing.T) {
	valid := []string{"CREATING", "CREATE_FAILED", "STARTING", "START_FAILED", "RUNNING", "STOPPING", "STOP_FAILED", "STOPPED", "UPDATING", "UPDATE_FAILED", "DELETING", "DELETE_FAILED", "DELETED"}
	for _, s := range valid {
		t.Run(s, func(t *testing.T) {
			got, err := ToOpenflowConnectorStatus(s)
			require.NoError(t, err)
			require.Equal(t, OpenflowConnectorStatus(s), got)
		})
	}
	_, err := ToOpenflowConnectorStatus("UNKNOWN")
	require.Error(t, err)
}
