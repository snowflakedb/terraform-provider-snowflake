package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenflowDeployments_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *CreateOpenflowDeploymentOptions {
		return &CreateOpenflowDeploymentOptions{name: id}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreateOpenflowDeploymentOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE OPENFLOW DEPLOYMENT %s", id.FullyQualifiedName())
	})

	t.Run("if not exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OPENFLOW DEPLOYMENT IF NOT EXISTS %s", id.FullyQualifiedName())
	})

	t.Run("snowflake type with all options", func(t *testing.T) {
		dt := OpenflowDeploymentTypeSnowflake
		comment := "my deployment"
		displayName := "My Deployment"
		opts := defaultOpts()
		opts.DeploymentType = &dt
		opts.Comment = &comment
		opts.DisplayName = &displayName
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW DEPLOYMENT %s DEPLOYMENT_TYPE = 'SNOWFLAKE' DISPLAY_NAME = 'My Deployment' COMMENT = 'my deployment'",
			id.FullyQualifiedName())
	})

	t.Run("byoc type", func(t *testing.T) {
		dt := OpenflowDeploymentTypeByoc
		vt := OpenflowVpcTypeManaged
		opts := defaultOpts()
		opts.DeploymentType = &dt
		opts.VpcType = &vt
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW DEPLOYMENT %s DEPLOYMENT_TYPE = 'BYOC' VPC_TYPE = 'MANAGED'",
			id.FullyQualifiedName())
	})

	t.Run("with private link", func(t *testing.T) {
		opts := defaultOpts()
		opts.UsePrivateLink = Bool(true)
		assertOptsValidAndSQLEquals(t, opts,
			"CREATE OPENFLOW DEPLOYMENT %s USE_PRIVATE_LINK = true",
			id.FullyQualifiedName())
	})
}

func TestOpenflowDeployments_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *AlterOpenflowDeploymentOptions {
		return &AlterOpenflowDeploymentOptions{name: id}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterOpenflowDeploymentOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one of Set/Unset", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterOpenflowDeploymentOptions", "Set", "Unset"))
	})

	t.Run("validation: Set requires at least one field", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &OpenflowDeploymentSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOpenflowDeploymentOptions.Set", "Comment", "DisplayName", "EventTable"))
	})

	t.Run("validation: Unset requires at least one field", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OpenflowDeploymentUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterOpenflowDeploymentOptions.Unset", "Comment", "DisplayName", "EventTable"))
	})

	t.Run("set comment", func(t *testing.T) {
		comment := "new comment"
		opts := defaultOpts()
		opts.Set = &OpenflowDeploymentSet{Comment: &comment}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW DEPLOYMENT %s SET COMMENT = 'new comment'", id.FullyQualifiedName())
	})

	t.Run("set display name", func(t *testing.T) {
		displayName := "My Deployment"
		opts := defaultOpts()
		opts.Set = &OpenflowDeploymentSet{DisplayName: &displayName}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW DEPLOYMENT %s SET DISPLAY_NAME = 'My Deployment'", id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OpenflowDeploymentUnset{Comment: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW DEPLOYMENT %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("unset multiple", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &OpenflowDeploymentUnset{Comment: Bool(true), DisplayName: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, "ALTER OPENFLOW DEPLOYMENT %s UNSET COMMENT, DISPLAY_NAME", id.FullyQualifiedName())
	})
}

func TestOpenflowDeployments_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()

	defaultOpts := func() *DropOpenflowDeploymentOptions {
		return &DropOpenflowDeploymentOptions{name: id}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropOpenflowDeploymentOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
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
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP OPENFLOW DEPLOYMENT IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestOpenflowDeployments_Show(t *testing.T) {
	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowOpenflowDeploymentOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := &ShowOpenflowDeploymentOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW DEPLOYMENTS")
	})

	t.Run("like", func(t *testing.T) {
		opts := &ShowOpenflowDeploymentOptions{Like: &Like{Pattern: String("my_dep%")}}
		assertOptsValidAndSQLEquals(t, opts, "SHOW OPENFLOW DEPLOYMENTS LIKE 'my_dep%%'")
	})
}

func TestOpenflowDeployments_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribeOpenflowDeploymentOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier", func(t *testing.T) {
		opts := &DescribeOpenflowDeploymentOptions{name: emptyAccountObjectIdentifier}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := &DescribeOpenflowDeploymentOptions{name: id}
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE OPENFLOW DEPLOYMENT %s", id.FullyQualifiedName())
	})
}

func Test_ToOpenflowDeploymentType(t *testing.T) {
	valid := []struct {
		input string
		want  OpenflowDeploymentType
	}{
		{"SNOWFLAKE", OpenflowDeploymentTypeSnowflake},
		{"snowflake", OpenflowDeploymentTypeSnowflake},
		{"BYOC", OpenflowDeploymentTypeByoc},
		{"byoc", OpenflowDeploymentTypeByoc},
	}
	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToOpenflowDeploymentType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
	for _, bad := range []string{"", "foo", "MANAGED"} {
		t.Run("invalid:"+bad, func(t *testing.T) {
			_, err := ToOpenflowDeploymentType(bad)
			require.Error(t, err)
		})
	}
}

func Test_ToOpenflowVpcType(t *testing.T) {
	valid := []struct {
		input string
		want  OpenflowVpcType
	}{
		{"MANAGED", OpenflowVpcTypeManaged},
		{"managed", OpenflowVpcTypeManaged},
		{"PROVIDED", OpenflowVpcTypeProvided},
	}
	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToOpenflowVpcType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
	_, err := ToOpenflowVpcType("SNOWFLAKE")
	require.Error(t, err)
}

func Test_ToOpenflowDeploymentStatus(t *testing.T) {
	valid := []string{"CREATING", "ACTIVE", "INACTIVE", "PROVISIONING", "NOT_REPORTING", "NOT_HEALTHY", "UPGRADING", "UPGRADE_FAILED", "DEACTIVATION_REQUIRED", "DELETING", "DELETED", "CREATE_FAILED", "DELETE_FAILED"}
	for _, s := range valid {
		t.Run(s, func(t *testing.T) {
			got, err := ToOpenflowDeploymentStatus(s)
			require.NoError(t, err)
			require.Equal(t, OpenflowDeploymentStatus(s), got)
		})
	}
	_, err := ToOpenflowDeploymentStatus("UNKNOWN")
	require.Error(t, err)
}
