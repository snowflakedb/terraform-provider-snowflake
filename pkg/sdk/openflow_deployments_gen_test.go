package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func init() {
	allEnumConversionTests = append(allEnumConversionTests, typedEnumTestProvider[OpenflowDeploymentType]{"OpenflowDeploymentType", AllOpenflowDeploymentTypes, ToOpenflowDeploymentType})
	allEnumConversionTests = append(allEnumConversionTests, typedEnumTestProvider[OpenflowVpcType]{"OpenflowVpcType", AllOpenflowVpcTypes, ToOpenflowVpcType})
	allEnumConversionTests = append(allEnumConversionTests, typedEnumTestProvider[OpenflowDeploymentStatus]{"OpenflowDeploymentStatus", AllOpenflowDeploymentStatuses, ToOpenflowDeploymentStatus})
}

func TestOpenflowDeployments_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *CreateOpenflowDeploymentOptions {
		return &CreateOpenflowDeploymentOptions{
			name:           id,
			DeploymentType: OpenflowDeploymentTypeSnowflake,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateOpenflowDeploymentOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE OPENFLOW DEPLOYMENT %s DEPLOYMENT_TYPE = 'SNOWFLAKE'", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		comment := random.Comment()
		vpcType := OpenflowVpcTypeManaged
		opts := &CreateOpenflowDeploymentOptions{
			IfNotExists:                new(true),
			name:                       id,
			DeploymentType:             OpenflowDeploymentTypeByoc,
			VpcType:                    &vpcType,
			CustomIngressHostname:      new("ingress.example.com"),
			UsePrivateLink:             new(true),
			UseUserAuthOverPrivatelink: new(false),
			EventTable:                 new("MY_DB.PUBLIC.EVENTS"),
			DisplayName:                new("My Deployment"),
			Comment:                    &comment,
		}
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW DEPLOYMENT IF NOT EXISTS %s DEPLOYMENT_TYPE = 'BYOC' VPC_TYPE = 'MANAGED'"+
				" CUSTOM_INGRESS_HOSTNAME = 'ingress.example.com' USE_PRIVATE_LINK = true USE_USER_AUTH_OVER_PRIVATELINK = false"+
				" EVENT_TABLE = 'MY_DB.PUBLIC.EVENTS' DISPLAY_NAME = 'My Deployment' COMMENT = '%s'",
			id.FullyQualifiedName(), comment)
	})
}

func TestOpenflowDeployments_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *AlterOpenflowDeploymentOptions {
		return &AlterOpenflowDeploymentOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterOpenflowDeploymentOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.RenameTo] if set", func(t *testing.T) {
		opts := defaultOpts()
		empty := emptyAccountObjectIdentifier
		opts.RenameTo = &empty
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Upgrade opts.Terminate opts.RenameTo opts.Set opts.Unset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOpenflowDeploymentOptions", "Upgrade", "Terminate", "RenameTo", "Set", "Unset"))
	})

	t.Run("validation: at least one of the fields [opts.Set.Comment opts.Set.DisplayName opts.Set.EventTable] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OpenflowDeploymentSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOpenflowDeploymentOptions.Set", "Comment", "DisplayName", "EventTable"))
	})

	t.Run("validation: at least one of the fields [opts.Unset.Comment opts.Unset.DisplayName opts.Unset.EventTable] should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OpenflowDeploymentUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOpenflowDeploymentOptions.Unset", "Comment", "DisplayName", "EventTable"))
	})

	t.Run("upgrade", func(t *testing.T) {
		opts := defaultOpts()
		opts.Upgrade = new(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW DEPLOYMENT %s UPGRADE", id.FullyQualifiedName())
	})

	t.Run("terminate", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terminate = new(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW DEPLOYMENT %s TERMINATE", id.FullyQualifiedName())
	})

	t.Run("rename to", func(t *testing.T) {
		newId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.RenameTo = &newId
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW DEPLOYMENT %s RENAME TO %s", id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		comment := random.Comment()
		opts := defaultOpts()
		opts.Set = &OpenflowDeploymentSet{
			Comment:     &comment,
			DisplayName: new("My Deployment"),
			EventTable:  new("MY_DB.PUBLIC.EVENTS"),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW DEPLOYMENT %s SET COMMENT = '%s' DISPLAY_NAME = 'My Deployment' EVENT_TABLE = 'MY_DB.PUBLIC.EVENTS'",
			id.FullyQualifiedName(), comment)
	})

	t.Run("unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OpenflowDeploymentUnset{
			Comment:     new(true),
			DisplayName: new(true),
			EventTable:  new(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW DEPLOYMENT %s UNSET COMMENT, DISPLAY_NAME, EVENT_TABLE", id.FullyQualifiedName())
	})
}

func TestOpenflowDeployments_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *DropOpenflowDeploymentOptions {
		return &DropOpenflowDeploymentOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropOpenflowDeploymentOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW DEPLOYMENT %s", id.FullyQualifiedName())
	})

	t.Run("if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = new(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW DEPLOYMENT IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestOpenflowDeployments_Show(t *testing.T) {
	defaultOpts := func() *ShowOpenflowDeploymentOptions {
		return &ShowOpenflowDeploymentOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowOpenflowDeploymentOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW DEPLOYMENTS")
	})

	t.Run("like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{Pattern: new("my-deployment%")}
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW DEPLOYMENTS LIKE 'my-deployment%%'")
	})
}

func TestOpenflowDeployments_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *DescribeOpenflowDeploymentOptions {
		return &DescribeOpenflowDeploymentOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeOpenflowDeploymentOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE OPENFLOW DEPLOYMENT %s", id.FullyQualifiedName())
	})
}
