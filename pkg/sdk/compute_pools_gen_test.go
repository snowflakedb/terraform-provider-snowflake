package sdk

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/stretchr/testify/require"
)

func TestComputePools_Create(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid CreateComputePoolOptions
	defaultOpts := func() *CreateComputePoolOptions {
		return &CreateComputePoolOptions{
			name:           id,
			MinNodes:       1,
			MaxNodes:       3,
			InstanceFamily: ComputePoolInstanceFamilyCPUX64XS,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateComputePoolOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: MinNodes must be greater than 0", func(t *testing.T) {
		opts := defaultOpts()
		opts.MinNodes = 0
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("CreateComputePoolOptions", "MinNodes", IntErrGreater, 0))
	})

	t.Run("validation: MaxNodes must be greater than or equal MinNodes", func(t *testing.T) {
		opts := defaultOpts()
		opts.MinNodes = 2
		opts.MaxNodes = 1
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("CreateComputePoolOptions", "MaxNodes", IntErrGreaterOrEqual, opts.MinNodes))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "CREATE COMPUTE POOL %s MIN_NODES = 1 MAX_NODES = 3 INSTANCE_FAMILY = CPU_X64_XS", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		appId := randomAccountObjectIdentifier()
		comment := random.Comment()
		tagId := NewAccountObjectIdentifier("tag1")
		opts := &CreateComputePoolOptions{
			IfNotExists:        Pointer(true),
			name:               id,
			ForApplication:     &appId,
			MinNodes:           2,
			MaxNodes:           3,
			InstanceFamily:     ComputePoolInstanceFamilyCPUX64XS,
			AutoResume:         Pointer(true),
			InitiallySuspended: Pointer(true),
			AutoSuspendSecs:    Pointer(42),
			Tag: []TagAssociation{
				{
					Name:  tagId,
					Value: "value1",
				},
			},
			Comment: &comment,
		}
		assertOptsValidAndSQLEquals(t, opts, fmt.Sprintf(`CREATE COMPUTE POOL IF NOT EXISTS %s FOR APPLICATION %s MIN_NODES = 2 MAX_NODES = 3`+
			` INSTANCE_FAMILY = CPU_X64_XS AUTO_RESUME = true INITIALLY_SUSPENDED = true AUTO_SUSPEND_SECS = 42 TAG (%s = 'value1')`+
			` COMMENT = '%s'`, id.FullyQualifiedName(), appId.FullyQualifiedName(), tagId.FullyQualifiedName(), comment))
	})
}

func TestComputePools_Alter(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid AlterComputePoolOptions
	defaultOpts := func() *AlterComputePoolOptions {
		return &AlterComputePoolOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterComputePoolOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags] should be present", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterComputePoolOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("validation: exactly one field from [opts.Set opts.Unset opts.SetTags opts.UnsetTags] should be present - more present", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &ComputePoolSet{}
		opts.Unset = &ComputePoolUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterComputePoolOptions", "Set", "Unset", "SetTags", "UnsetTags"))
	})

	t.Run("set", func(t *testing.T) {
		comment := random.Comment()
		opts := defaultOpts()
		opts.Set = &ComputePoolSet{
			MinNodes:        Pointer(2),
			MaxNodes:        Pointer(3),
			AutoResume:      Pointer(true),
			AutoSuspendSecs: Pointer(60),
			Comment:         &comment,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER COMPUTE POOL %s SET MIN_NODES = 2 MAX_NODES = 3 AUTO_RESUME = true AUTO_SUSPEND_SECS = 60`+
			` COMMENT = '%s'`, id.FullyQualifiedName(), comment)
	})

	t.Run("unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &ComputePoolUnset{
			AutoResume:      Pointer(true),
			AutoSuspendSecs: Pointer(true),
			Comment:         Pointer(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER COMPUTE POOL %s UNSET AUTO_RESUME, AUTO_SUSPEND_SECS, COMMENT", id.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
			{
				Name:  NewAccountObjectIdentifier("tag2"),
				Value: "value2",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER COMPUTE POOL %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, id.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER COMPUTE POOL %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})
}

func TestComputePools_Drop(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid DropComputePoolOptions
	defaultOpts := func() *DropComputePoolOptions {
		return &DropComputePoolOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropComputePoolOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP COMPUTE POOL %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP COMPUTE POOL IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestComputePools_Show(t *testing.T) {
	// Minimal valid ShowComputePoolOptions
	defaultOpts := func() *ShowComputePoolOptions {
		return &ShowComputePoolOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowComputePoolOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW COMPUTE POOLS")
	})

	t.Run("like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW COMPUTE POOLS LIKE 'pattern'")
	})

	t.Run("starts with", func(t *testing.T) {
		opts := defaultOpts()
		opts.StartsWith = Pointer("prefix")
		assertOptsValidAndSQLEquals(t, opts, "SHOW COMPUTE POOLS STARTS WITH 'prefix'")
	})

	t.Run("limit from", func(t *testing.T) {
		opts := defaultOpts()
		opts.Limit = &LimitFrom{
			Rows: Int(10),
			From: String("from"),
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW COMPUTE POOLS LIMIT 10 FROM 'from'")
	})
}

func TestComputePools_Describe(t *testing.T) {
	id := randomAccountObjectIdentifier()
	// Minimal valid DescribeComputePoolOptions
	defaultOpts := func() *DescribeComputePoolOptions {
		return &DescribeComputePoolOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeComputePoolOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})
	t.Run("validation: valid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE COMPUTE POOL %s", id.FullyQualifiedName())
	})
}

func Test_ComputePool_ToComputePoolInstanceFamily(t *testing.T) {
	type test struct {
		input string
		want  ComputePoolInstanceFamily
	}

	valid := []test{
		// case insensitive.
		{input: "cpu_x64_xs", want: ComputePoolInstanceFamilyCPUX64XS},

		// Supported Values
		{input: "CPU_X64_XS", want: ComputePoolInstanceFamilyCPUX64XS},
		{input: "CPU_X64_S", want: ComputePoolInstanceFamilyCPUX64S},
		{input: "CPU_X64_M", want: ComputePoolInstanceFamilyCPUX64M},
		{input: "CPU_X64_L", want: ComputePoolInstanceFamilyCPUX64L},
		{input: "HIGHMEM_X64_S", want: ComputePoolInstanceFamilyHIGHMEMX64S},
		{input: "HIGHMEM_X64_M", want: ComputePoolInstanceFamilyHIGHMEMX64M},
		{input: "HIGHMEM_X64_L", want: ComputePoolInstanceFamilyHIGHMEMX64L},
		{input: "HIGHMEM_X64_SL", want: ComputePoolInstanceFamilyHIGHMEMX64SL},
		{input: "GPU_NV_S", want: ComputePoolInstanceFamilyGPUNVS},
		{input: "GPU_NV_M", want: ComputePoolInstanceFamilyGPUNVM},
		{input: "GPU_NV_L", want: ComputePoolInstanceFamilyGPUNVL},
		{input: "GPU_NV_XS", want: ComputePoolInstanceFamilyGPUNVXS},
		{input: "GPU_NV_SM", want: ComputePoolInstanceFamilyGPUNVSM},
		{input: "GPU_NV_2M", want: ComputePoolInstanceFamilyGPUNV2M},
		{input: "GPU_NV_3M", want: ComputePoolInstanceFamilyGPUNV3M},
		{input: "GPU_NV_SL", want: ComputePoolInstanceFamilyGPUNVSL},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
		{input: "cpux64xs"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToComputePoolInstanceFamily(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToComputePoolInstanceFamily(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_ComputePool_ToComputePoolState(t *testing.T) {
	type test struct {
		input string
		want  ComputePoolState
	}

	valid := []test{
		// case insensitive.
		{input: "idle", want: ComputePoolStateIdle},

		// Supported Values
		{input: "IDLE", want: ComputePoolStateIdle},
		{input: "ACTIVE", want: ComputePoolStateActive},
		{input: "SUSPENDED", want: ComputePoolStateSuspended},
		{input: "STARTING", want: ComputePoolStateStarting},
		{input: "STOPPING", want: ComputePoolStateStopping},
		{input: "RESIZING", want: ComputePoolStateResizing},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToComputePoolState(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToComputePoolState(tc.input)
			require.Error(t, err)
		})
	}
}
