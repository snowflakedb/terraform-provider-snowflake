package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func init() {
	allEnumConversionTests = append(allEnumConversionTests, typedEnumTestProvider[OpenflowRuntimeNodeType]{"OpenflowRuntimeNodeType", AllOpenflowRuntimeNodeTypes, ToOpenflowRuntimeNodeType})
	allEnumConversionTests = append(allEnumConversionTests, typedEnumTestProvider[OpenflowRuntimeStatus]{"OpenflowRuntimeStatus", AllOpenflowRuntimeStatuses, ToOpenflowRuntimeStatus})
}

func TestOpenflowRuntimes_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	deploymentId := randomAccountObjectIdentifier()
	defaultOpts := func() *CreateOpenflowRuntimeOptions {
		return &CreateOpenflowRuntimeOptions{
			name:          id,
			InDeployment:  deploymentId,
			ExecuteAsRole: "SYSADMIN",
			NodeType:      OpenflowRuntimeNodeTypeSmall,
			MinNodes:      1,
			MaxNodes:      3,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateOpenflowRuntimeOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW RUNTIME %s IN DEPLOYMENT %s EXECUTE_AS_ROLE = SYSADMIN NODE_TYPE = 'SMALL' MIN_NODES = 1 MAX_NODES = 3",
			id.FullyQualifiedName(), deploymentId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		comment := random.Comment()
		eaiId := randomAccountObjectIdentifier()
		opts := &CreateOpenflowRuntimeOptions{
			IfNotExists:   Bool(true),
			name:          id,
			InDeployment:  deploymentId,
			ExecuteAsRole: "MYROLE",
			NodeType:      OpenflowRuntimeNodeTypeLarge,
			MinNodes:      2,
			MaxNodes:      5,
			ExternalAccessIntegrations: &OpenflowRuntimeExternalAccessIntegrations{
				ExternalAccessIntegrations: []AccountObjectIdentifier{eaiId},
			},
			DisplayName: String("My Runtime"),
			Comment:     &comment,
		}
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW RUNTIME IF NOT EXISTS %s IN DEPLOYMENT %s EXECUTE_AS_ROLE = MYROLE NODE_TYPE = 'LARGE' MIN_NODES = 2 MAX_NODES = 5"+
				" EXTERNAL_ACCESS_INTEGRATIONS = (%s) DISPLAY_NAME = 'My Runtime' COMMENT = '%s'",
			id.FullyQualifiedName(), deploymentId.FullyQualifiedName(), eaiId.FullyQualifiedName(), comment)
	})
}

func TestOpenflowRuntimes_Alter(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *AlterOpenflowRuntimeOptions {
		return &AlterOpenflowRuntimeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterOpenflowRuntimeOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.RenameTo] if set", func(t *testing.T) {
		opts := defaultOpts()
		empty := emptySchemaObjectIdentifier
		opts.RenameTo = &empty
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Suspend opts.Resume opts.ResumeRecovery opts.Restart opts.RestartRecovery opts.Terminate opts.TerminateCascade opts.Upgrade opts.RenameTo opts.Set opts.Unset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOpenflowRuntimeOptions", "Suspend", "Resume", "ResumeRecovery", "Restart", "RestartRecovery", "Terminate", "TerminateCascade", "Upgrade", "RenameTo", "Set", "Unset"))
	})

	t.Run("validation: at least one of the fields [opts.Set.MinNodes opts.Set.MaxNodes opts.Set.ExecuteAsRole opts.Set.ExternalAccessIntegrations opts.Set.DisplayName opts.Set.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OpenflowRuntimeSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOpenflowRuntimeOptions.Set", "MinNodes", "MaxNodes", "ExecuteAsRole", "ExternalAccessIntegrations", "DisplayName", "Comment"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.ExecuteAsRole opts.Unset.ExternalAccessIntegrations opts.Unset.DisplayName opts.Unset.Comment] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OpenflowRuntimeUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOpenflowRuntimeOptions.Unset", "ExecuteAsRole", "ExternalAccessIntegrations", "DisplayName", "Comment"))
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

	t.Run("resume recovery", func(t *testing.T) {
		opts := defaultOpts()
		opts.ResumeRecovery = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s RESUME RECOVERY", id.FullyQualifiedName())
	})

	t.Run("restart", func(t *testing.T) {
		opts := defaultOpts()
		opts.Restart = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s RESTART", id.FullyQualifiedName())
	})

	t.Run("restart recovery", func(t *testing.T) {
		opts := defaultOpts()
		opts.RestartRecovery = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s RESTART RECOVERY", id.FullyQualifiedName())
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

	t.Run("upgrade", func(t *testing.T) {
		opts := defaultOpts()
		opts.Upgrade = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s UPGRADE", id.FullyQualifiedName())
	})

	t.Run("rename to", func(t *testing.T) {
		newId := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.RenameTo = &newId
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s RENAME TO %s", id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		comment := random.Comment()
		eaiId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &OpenflowRuntimeSet{
			MinNodes:      Int(2),
			MaxNodes:      Int(5),
			ExecuteAsRole: String("MYROLE"),
			ExternalAccessIntegrations: &OpenflowRuntimeExternalAccessIntegrations{
				ExternalAccessIntegrations: []AccountObjectIdentifier{eaiId},
			},
			DisplayName: String("Updated Runtime"),
			Comment:     &comment,
		}
		assertOptsValidAndSQLEquals(t, opts,
			"ALTER OPENFLOW RUNTIME %s SET MIN_NODES = 2 MAX_NODES = 5 EXECUTE_AS_ROLE = MYROLE EXTERNAL_ACCESS_INTEGRATIONS = (%s) DISPLAY_NAME = 'Updated Runtime' COMMENT = '%s'",
			id.FullyQualifiedName(), eaiId.FullyQualifiedName(), comment)
	})

	t.Run("unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OpenflowRuntimeUnset{
			ExecuteAsRole:              Bool(true),
			ExternalAccessIntegrations: Bool(true),
			DisplayName:                Bool(true),
			Comment:                    Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW RUNTIME %s UNSET EXECUTE_AS_ROLE, EXTERNAL_ACCESS_INTEGRATIONS, DISPLAY_NAME, COMMENT",
			id.FullyQualifiedName())
	})
}

func TestOpenflowRuntimes_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *DropOpenflowRuntimeOptions {
		return &DropOpenflowRuntimeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropOpenflowRuntimeOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW RUNTIME %s", id.FullyQualifiedName())
	})

	t.Run("if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW RUNTIME IF EXISTS %s", id.FullyQualifiedName())
	})

	t.Run("cascade", func(t *testing.T) {
		opts := defaultOpts()
		opts.Cascade = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW RUNTIME %s CASCADE", id.FullyQualifiedName())
	})
}

func TestOpenflowRuntimes_Show(t *testing.T) {
	defaultOpts := func() *ShowOpenflowRuntimeOptions {
		return &ShowOpenflowRuntimeOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowOpenflowRuntimeOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW RUNTIMES")
	})

	t.Run("like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{Pattern: String("my-runtime%")}
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW RUNTIMES LIKE 'my-runtime%%'")
	})

	t.Run("in schema", func(t *testing.T) {
		schemaId := randomDatabaseObjectIdentifier()
		opts := defaultOpts()
		opts.In = &In{Schema: schemaId}
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW RUNTIMES IN SCHEMA %s", schemaId.FullyQualifiedName())
	})
}

func TestOpenflowRuntimes_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *DescribeOpenflowRuntimeOptions {
		return &DescribeOpenflowRuntimeOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeOpenflowRuntimeOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE OPENFLOW RUNTIME %s", id.FullyQualifiedName())
	})
}
