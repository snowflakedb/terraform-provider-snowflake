package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func init() {
	allEnumConversionTests = append(allEnumConversionTests, typedEnumTestProvider[OpenflowConnectorStatus]{"OpenflowConnectorStatus", AllOpenflowConnectorStatuses, ToOpenflowConnectorStatus})
}

func TestOpenflowConnectors_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	runtimeId := randomSchemaObjectIdentifier()
	defaultOpts := func() *CreateOpenflowConnectorOptions {
		return &CreateOpenflowConnectorOptions{
			name:      id,
			InRuntime: &runtimeId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateOpenflowConnectorOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: FromDefinition and From are conflicting fields", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromDefinition = String("mydef")
		opts.From = String("@MY_STAGE/path/")
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateOpenflowConnectorOptions", "FromDefinition", "From"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE OPENFLOW CONNECTOR %s IN RUNTIME %s",
			id.FullyQualifiedName(), runtimeId.FullyQualifiedName())
	})

	t.Run("from definition", func(t *testing.T) {
		opts := &CreateOpenflowConnectorOptions{
			name:           id,
			InRuntime:      &runtimeId,
			FromDefinition: String("mydef"),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE OPENFLOW CONNECTOR %s IN RUNTIME %s FROM DEFINITION mydef",
			id.FullyQualifiedName(), runtimeId.FullyQualifiedName())
	})

	t.Run("from stage", func(t *testing.T) {
		opts := &CreateOpenflowConnectorOptions{
			name:      id,
			InRuntime: &runtimeId,
			From:      String("@MY_STAGE/path/"),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE OPENFLOW CONNECTOR %s IN RUNTIME %s FROM '@MY_STAGE/path/'",
			id.FullyQualifiedName(), runtimeId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		comment := random.Comment()
		opts := &CreateOpenflowConnectorOptions{
			IfNotExists:    Bool(true),
			name:           id,
			InRuntime:      &runtimeId,
			FromDefinition: String("mydef"),
			DisplayName:    String("My Connector"),
			Comment:        &comment,
		}
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW CONNECTOR IF NOT EXISTS %s IN RUNTIME %s FROM DEFINITION mydef DISPLAY_NAME = 'My Connector' COMMENT = '%s'",
			id.FullyQualifiedName(), runtimeId.FullyQualifiedName(), comment)
	})
}

func TestOpenflowConnectors_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *AlterOpenflowConnectorOptions {
		return &AlterOpenflowConnectorOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterOpenflowConnectorOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Start opts.Stop opts.Terminate opts.Commit opts.Abort opts.Set opts.Unset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOpenflowConnectorOptions", "Start", "Stop", "Terminate", "Commit", "Abort", "Set", "Unset"))
	})

	t.Run("validation: at least one of the fields [opts.Set.DisplayName opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OpenflowConnectorSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOpenflowConnectorOptions.Set", "DisplayName", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.DisplayName opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OpenflowConnectorUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOpenflowConnectorOptions.Unset", "DisplayName", "Comment"))
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

	t.Run("terminate", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terminate = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW CONNECTOR %s TERMINATE", id.FullyQualifiedName())
	})

	t.Run("commit", func(t *testing.T) {
		opts := defaultOpts()
		opts.Commit = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW CONNECTOR %s COMMIT", id.FullyQualifiedName())
	})

	t.Run("abort", func(t *testing.T) {
		opts := defaultOpts()
		opts.Abort = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW CONNECTOR %s ABORT", id.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		comment := random.Comment()
		opts := defaultOpts()
		opts.Set = &OpenflowConnectorSet{
			DisplayName: String("Updated Connector"),
			Comment:     &comment,
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW CONNECTOR %s SET DISPLAY_NAME = 'Updated Connector' COMMENT = '%s'",
			id.FullyQualifiedName(), comment)
	})

	t.Run("unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OpenflowConnectorUnset{
			DisplayName: Bool(true),
			Comment:     Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW CONNECTOR %s UNSET DISPLAY_NAME, COMMENT", id.FullyQualifiedName())
	})
}

func TestOpenflowConnectors_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *DropOpenflowConnectorOptions {
		return &DropOpenflowConnectorOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropOpenflowConnectorOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW CONNECTOR %s", id.FullyQualifiedName())
	})

	t.Run("if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW CONNECTOR IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestOpenflowConnectors_Show(t *testing.T) {
	defaultOpts := func() *ShowOpenflowConnectorOptions {
		return &ShowOpenflowConnectorOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowOpenflowConnectorOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW CONNECTORS")
	})

	t.Run("like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{Pattern: String("my-connector%")}
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW CONNECTORS LIKE 'my-connector%%'")
	})

	t.Run("in schema", func(t *testing.T) {
		schemaId := randomDatabaseObjectIdentifier()
		opts := defaultOpts()
		opts.In = &In{Schema: schemaId}
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW CONNECTORS IN SCHEMA %s", schemaId.FullyQualifiedName())
	})
}

func TestOpenflowConnectors_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *DescribeOpenflowConnectorOptions {
		return &DescribeOpenflowConnectorOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeOpenflowConnectorOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE OPENFLOW CONNECTOR %s", id.FullyQualifiedName())
	})
}
