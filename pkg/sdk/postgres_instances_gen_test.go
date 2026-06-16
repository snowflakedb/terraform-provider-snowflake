package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func init() {
	allEnumConversionTests = append(allEnumConversionTests, typedEnumTestProvider[PostgresInstanceState]{"PostgresInstanceState", AllPostgresInstanceStates, ToPostgresInstanceState})
	allEnumConversionTests = append(allEnumConversionTests, typedEnumTestProvider[PostgresInstanceAuthenticationAuthority]{"PostgresInstanceAuthenticationAuthority", AllPostgresInstanceAuthenticationAuthorities, ToPostgresInstanceAuthenticationAuthority})
	allEnumConversionTests = append(allEnumConversionTests, typedEnumTestProvider[PostgresInstanceResetAccessRole]{"PostgresInstanceResetAccessRole", AllPostgresInstanceResetAccessRoles, ToPostgresInstanceResetAccessRole})
}

func TestPostgresInstances_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *CreatePostgresInstanceOptions {
		return &CreatePostgresInstanceOptions{
			name:                    id,
			ComputeFamily:           "STANDARD_S",
			StorageSizeGb:           50,
			AuthenticationAuthority: PostgresInstanceAuthenticationAuthorityPostgres,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*CreatePostgresInstanceOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts,
			`CREATE POSTGRES INSTANCE %s COMPUTE_FAMILY = 'STANDARD_S' STORAGE_SIZE_GB = 50 AUTHENTICATION_AUTHORITY = POSTGRES`,
			id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		comment := random.Comment()
		tagId := NewAccountObjectIdentifier("tag1")
		opts := &CreatePostgresInstanceOptions{
			name:                    id,
			ComputeFamily:           "STANDARD_S",
			StorageSizeGb:           50,
			AuthenticationAuthority: PostgresInstanceAuthenticationAuthorityPostgres,
			PostgresVersion:         new(17),
			NetworkPolicy:           new("my_policy"),
			HighAvailability:        new(true),
			StorageIntegration:      new("my_integration"),
			PostgresSettings:        new("{}"),
			Comment:                 &comment,
			Tag: []TagAssociation{
				{
					Name:  tagId,
					Value: "value1",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts,
			`CREATE POSTGRES INSTANCE %s COMPUTE_FAMILY = 'STANDARD_S' STORAGE_SIZE_GB = 50`+
				` AUTHENTICATION_AUTHORITY = POSTGRES POSTGRES_VERSION = 17 NETWORK_POLICY = 'my_policy'`+
				` HIGH_AVAILABILITY = true STORAGE_INTEGRATION = 'my_integration' POSTGRES_SETTINGS = '{}'`+
				` COMMENT = '%s' TAG (%s = 'value1')`,
			id.FullyQualifiedName(), comment, tagId.FullyQualifiedName())
	})
}

func TestPostgresInstances_Fork(t *testing.T) {
	id := randomAccountObjectIdentifier()
	forkId := randomAccountObjectIdentifier()
	defaultOpts := func() *ForkPostgresInstanceOptions {
		return &ForkPostgresInstanceOptions{
			name: id,
			Fork: forkId,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ForkPostgresInstanceOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: valid identifier for [opts.Fork]", func(t *testing.T) {
		opts := defaultOpts()
		opts.Fork = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.At.Timestamp opts.At.Offset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.At = &PostgresInstanceForkAt{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ForkPostgresInstanceOptions.At", "Timestamp", "Offset"))
	})

	t.Run("validation: exactly one field from [opts.Before.Timestamp opts.Before.Offset] should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Before = &PostgresInstanceForkBefore{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("ForkPostgresInstanceOptions.Before", "Timestamp", "Offset"))
	})

	t.Run("validation: conflicting fields [opts.At opts.Before]", func(t *testing.T) {
		opts := defaultOpts()
		opts.At = &PostgresInstanceForkAt{
			Timestamp: new("2025-01-15 12:00:00"),
		}
		opts.Before = &PostgresInstanceForkBefore{
			Timestamp: new("2025-01-15 12:00:00"),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("ForkPostgresInstanceOptions", "At", "Before"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts,
			`CREATE POSTGRES INSTANCE %s FORK %s`,
			id.FullyQualifiedName(), forkId.FullyQualifiedName())
	})

	t.Run("fork with at timestamp", func(t *testing.T) {
		opts := defaultOpts()
		opts.At = &PostgresInstanceForkAt{
			Timestamp: new("2025-01-15 12:00:00"),
		}
		assertOptsValidAndSQLEquals(t, opts,
			`CREATE POSTGRES INSTANCE %s FORK %s AT (TIMESTAMP => '2025-01-15 12:00:00')`,
			id.FullyQualifiedName(), forkId.FullyQualifiedName())
	})

	t.Run("fork with at offset", func(t *testing.T) {
		opts := defaultOpts()
		opts.At = &PostgresInstanceForkAt{
			Offset: new("-7200"),
		}
		assertOptsValidAndSQLEquals(t, opts,
			`CREATE POSTGRES INSTANCE %s FORK %s AT (OFFSET => -7200)`,
			id.FullyQualifiedName(), forkId.FullyQualifiedName())
	})

	t.Run("fork with before timestamp", func(t *testing.T) {
		opts := defaultOpts()
		opts.Before = &PostgresInstanceForkBefore{
			Timestamp: new("2025-01-15 12:00:00"),
		}
		assertOptsValidAndSQLEquals(t, opts,
			`CREATE POSTGRES INSTANCE %s FORK %s BEFORE (TIMESTAMP => '2025-01-15 12:00:00')`,
			id.FullyQualifiedName(), forkId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		comment := random.Comment()
		tagId := NewAccountObjectIdentifier("tag1")
		opts := &ForkPostgresInstanceOptions{
			name: id,
			Fork: forkId,
			At: &PostgresInstanceForkAt{
				Timestamp: new("2025-01-15 12:00:00"),
			},
			ComputeFamily:    new("STANDARD_M"),
			StorageSizeGb:    new(100),
			HighAvailability: new(true),
			PostgresSettings: new("{}"),
			Comment:          &comment,
			Tag: []TagAssociation{
				{
					Name:  tagId,
					Value: "value1",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts,
			`CREATE POSTGRES INSTANCE %s FORK %s AT (TIMESTAMP => '2025-01-15 12:00:00')`+
				` COMPUTE_FAMILY = 'STANDARD_M' STORAGE_SIZE_GB = 100 HIGH_AVAILABILITY = true`+
				` POSTGRES_SETTINGS = '{}' COMMENT = '%s' TAG (%s = 'value1')`,
			id.FullyQualifiedName(), forkId.FullyQualifiedName(), comment, tagId.FullyQualifiedName())
	})
}

func TestPostgresInstances_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *AlterPostgresInstanceOptions {
		return &AlterPostgresInstanceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*AlterPostgresInstanceOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterPostgresInstanceOptions", "RenameTo", "Set", "Unset", "Suspend", "Resume", "ResetAccess", "SetTags", "UnsetTags"))
	})

	t.Run("validation: exactly one field - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &PostgresInstanceSet{}
		opts.Unset = &PostgresInstanceUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterPostgresInstanceOptions", "RenameTo", "Set", "Unset", "Suspend", "Resume", "ResetAccess", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one of the Set fields should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &PostgresInstanceSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterPostgresInstanceOptions.Set", "NetworkPolicy", "AuthenticationAuthority", "Comment", "HighAvailability", "ComputeFamily", "StorageSizeGb", "StorageIntegration", "PostgresVersion", "MaintenanceWindowStart", "PostgresSettings"))
	})

	t.Run("validation: at least one of the Unset fields should be set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &PostgresInstanceUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterPostgresInstanceOptions.Unset", "Comment", "PostgresSettings", "NetworkPolicy", "MaintenanceWindowStart", "StorageIntegration"))
	})

	t.Run("rename", func(t *testing.T) {
		newId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.RenameTo = &newId
		assertOptsValidAndSQLEquals(t, opts, `ALTER POSTGRES INSTANCE %s RENAME TO %s`, id.FullyQualifiedName(), newId.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		comment := random.Comment()
		auth := PostgresInstanceAuthenticationAuthorityPostgresOrSnowflake
		opts := defaultOpts()
		opts.Set = &PostgresInstanceSet{
			NetworkPolicy:           new("my_policy"),
			AuthenticationAuthority: &auth,
			Comment:                 &comment,
			HighAvailability:        new(true),
			ComputeFamily:           new("STANDARD_M"),
			StorageSizeGb:           new(100),
			PostgresVersion:         new(18),
			PostgresSettings:        new("{}"),
		}
		assertOptsValidAndSQLEquals(t, opts,
			`ALTER POSTGRES INSTANCE %s SET NETWORK_POLICY = 'my_policy' AUTHENTICATION_AUTHORITY = POSTGRES_OR_SNOWFLAKE`+
				` COMMENT = '%s' HIGH_AVAILABILITY = true COMPUTE_FAMILY = 'STANDARD_M' STORAGE_SIZE_GB = 100`+
				` POSTGRES_VERSION = 18 POSTGRES_SETTINGS = '{}'`,
			id.FullyQualifiedName(), comment)
	})

	t.Run("set with apply immediately", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &PostgresInstanceSet{
			ComputeFamily: new("STANDARD_L"),
			Apply: &PostgresInstanceApply{
				Immediately: new(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts,
			`ALTER POSTGRES INSTANCE %s SET COMPUTE_FAMILY = 'STANDARD_L' APPLY IMMEDIATELY`,
			id.FullyQualifiedName())
	})

	t.Run("set with apply on timestamp", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &PostgresInstanceSet{
			PostgresVersion: new(18),
			Apply: &PostgresInstanceApply{
				On: new("2026-03-01 12:00:00"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts,
			`ALTER POSTGRES INSTANCE %s SET POSTGRES_VERSION = 18 APPLY ON '2026-03-01 12:00:00'`,
			id.FullyQualifiedName())
	})

	t.Run("validation: exactly one of Apply.Immediately or Apply.On should be present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &PostgresInstanceSet{
			ComputeFamily: new("STANDARD_L"),
			Apply:         &PostgresInstanceApply{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterPostgresInstanceOptions.Set.Apply", "Immediately", "On"))
	})

	t.Run("unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &PostgresInstanceUnset{
			Comment:          new(true),
			PostgresSettings: new(true),
			NetworkPolicy:    new(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER POSTGRES INSTANCE %s UNSET COMMENT, POSTGRES_SETTINGS, NETWORK_POLICY`, id.FullyQualifiedName())
	})

	t.Run("suspend", func(t *testing.T) {
		opts := defaultOpts()
		opts.Suspend = new(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER POSTGRES INSTANCE %s SUSPEND`, id.FullyQualifiedName())
	})

	t.Run("resume", func(t *testing.T) {
		opts := defaultOpts()
		opts.Resume = new(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER POSTGRES INSTANCE %s RESUME`, id.FullyQualifiedName())
	})

	t.Run("reset access for snowflake_admin", func(t *testing.T) {
		opts := defaultOpts()
		opts.ResetAccess = &PostgresInstanceResetAccess{
			ForRole: PostgresInstanceResetAccessRoleSnowflakeAdmin,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER POSTGRES INSTANCE %s RESET ACCESS FOR 'snowflake_admin'`, id.FullyQualifiedName())
	})

	t.Run("reset access for application", func(t *testing.T) {
		opts := defaultOpts()
		opts.ResetAccess = &PostgresInstanceResetAccess{
			ForRole: PostgresInstanceResetAccessRoleApplication,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER POSTGRES INSTANCE %s RESET ACCESS FOR 'application'`, id.FullyQualifiedName())
	})

	t.Run("if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = new(true)
		opts.Suspend = new(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER POSTGRES INSTANCE IF EXISTS %s SUSPEND`, id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER POSTGRES INSTANCE %s SET TAG "tag1" = 'value1'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER POSTGRES INSTANCE %s UNSET TAG "tag1"`, id.FullyQualifiedName())
	})
}

func TestPostgresInstances_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *DropPostgresInstanceOptions {
		return &DropPostgresInstanceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DropPostgresInstanceOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP POSTGRES INSTANCE %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = new(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP POSTGRES INSTANCE IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestPostgresInstances_Show(t *testing.T) {
	defaultOpts := func() *ShowPostgresInstanceOptions {
		return &ShowPostgresInstanceOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*ShowPostgresInstanceOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW POSTGRES INSTANCES")
	})

	t.Run("like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: new("pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW POSTGRES INSTANCES LIKE 'pattern'")
	})

	t.Run("starts with", func(t *testing.T) {
		opts := defaultOpts()
		opts.StartsWith = new("prefix")
		assertOptsValidAndSQLEquals(t, opts, "SHOW POSTGRES INSTANCES STARTS WITH 'prefix'")
	})

	t.Run("limit from", func(t *testing.T) {
		opts := defaultOpts()
		opts.Limit = &LimitFrom{
			Rows: new(10),
			From: new("from"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW POSTGRES INSTANCES LIMIT 10 FROM 'from'")
	})
}

func TestPostgresInstances_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()
	defaultOpts := func() *DescribePostgresInstanceOptions {
		return &DescribePostgresInstanceOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*DescribePostgresInstanceOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE POSTGRES INSTANCE %s", id.FullyQualifiedName())
	})
}
