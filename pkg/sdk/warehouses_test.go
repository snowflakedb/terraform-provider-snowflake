package sdk

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWarehouseCreate(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name: AccountObjectIdentifier{
				name: "mywarehouse",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE WAREHOUSE "mywarehouse"`)
	})

	t.Run("with complete options", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		resourceMonitorId := randomAccountObjectIdentifier()
		opts := &CreateWarehouseOptions{
			OrReplace:   Bool(true),
			name:        NewAccountObjectIdentifier("completewarehouse"),
			IfNotExists: Bool(true),

			WarehouseType:                   Pointer(WarehouseTypeStandard),
			WarehouseSize:                   Pointer(WarehouseSizeX4Large),
			MaxClusterCount:                 Int(8),
			MinClusterCount:                 Int(3),
			ScalingPolicy:                   Pointer(ScalingPolicyEconomy),
			AutoSuspend:                     Int(1000),
			AutoResume:                      Bool(true),
			InitiallySuspended:              Bool(false),
			ResourceMonitor:                 Pointer(resourceMonitorId),
			Comment:                         String("hello"),
			EnableQueryAcceleration:         Bool(true),
			QueryAccelerationMaxScaleFactor: Int(62),

			MaxConcurrencyLevel:             Int(7),
			StatementQueuedTimeoutInSeconds: Int(29),
			StatementTimeoutInSeconds:       Int(89),
			Tag: []TagAssociation{
				{
					Name:  tagId1,
					Value: "v1",
				},
				{
					Name:  tagId2,
					Value: "v2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE WAREHOUSE IF NOT EXISTS "completewarehouse" WAREHOUSE_TYPE = 'STANDARD' WAREHOUSE_SIZE = 'X4LARGE' MAX_CLUSTER_COUNT = 8 MIN_CLUSTER_COUNT = 3 SCALING_POLICY = 'ECONOMY' AUTO_SUSPEND = 1000 AUTO_RESUME = true INITIALLY_SUSPENDED = false RESOURCE_MONITOR = %s COMMENT = 'hello' ENABLE_QUERY_ACCELERATION = true QUERY_ACCELERATION_MAX_SCALE_FACTOR = 62 MAX_CONCURRENCY_LEVEL = 7 STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 29 STATEMENT_TIMEOUT_IN_SECONDS = 89 TAG (%s = 'v1', %s = 'v2')`, resourceMonitorId.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})
}

func TestWarehouseSizing(t *testing.T) {
	t.Run("validation: Min bigger than Max", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			MaxClusterCount: Int(1),
			MinClusterCount: Int(2),
		}
		assertOptsInvalidJoinedErrors(t, opts, fmt.Errorf("MinClusterCount must be less than or equal to MaxClusterCount"))
	})

	t.Run("validation: only min set", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			MinClusterCount: Int(2),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE WAREHOUSE \"mywarehouse\" MIN_CLUSTER_COUNT = 2")
	})

	t.Run("Max equal Min", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			MaxClusterCount: Int(2),
			MinClusterCount: Int(2),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE WAREHOUSE \"mywarehouse\" MAX_CLUSTER_COUNT = 2 MIN_CLUSTER_COUNT = 2")
	})

	t.Run("Max greater than Min", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			MaxClusterCount: Int(2),
			MinClusterCount: Int(1),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE WAREHOUSE \"mywarehouse\" MAX_CLUSTER_COUNT = 2 MIN_CLUSTER_COUNT = 1")
	})

	t.Run("Allow large Min Max Values", func(t *testing.T) {
		opts := &CreateWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			MaxClusterCount: Int(100),
			MinClusterCount: Int(11),
		}
		assertOptsValidAndSQLEquals(t, opts, "CREATE WAREHOUSE \"mywarehouse\" MAX_CLUSTER_COUNT = 100 MIN_CLUSTER_COUNT = 11")
	})
}

// TODO: add validation tests
func TestWarehouseAlter(t *testing.T) {
	t.Run("with set params", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			Set: &WarehouseSet{
				WarehouseType:                   Pointer(WarehouseTypeSnowparkOptimized),
				WaitForCompletion:               Bool(false),
				MinClusterCount:                 Int(4),
				MaxClusterCount:                 Int(5),
				AutoSuspend:                     Int(200),
				ResourceMonitor:                 NewAccountObjectIdentifier("resmon"),
				EnableQueryAcceleration:         Bool(false),
				StatementQueuedTimeoutInSeconds: Int(1200),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" SET WAREHOUSE_TYPE = 'SNOWPARK-OPTIMIZED' WAIT_FOR_COMPLETION = false MAX_CLUSTER_COUNT = 5 MIN_CLUSTER_COUNT = 4 AUTO_SUSPEND = 200 RESOURCE_MONITOR = "resmon" ENABLE_QUERY_ACCELERATION = false STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 1200`)
	})

	t.Run("with set tag", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			SetTag: []TagAssociation{
				{
					Name:  tagId1,
					Value: "v1",
				},
				{
					Name:  tagId2,
					Value: "v2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" SET TAG %s = 'v1', %s = 'v2'`, tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("with unset tag", func(t *testing.T) {
		tagId := randomSchemaObjectIdentifier()
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			UnsetTag: []ObjectIdentifier{
				tagId,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" UNSET TAG %s`, tagId.FullyQualifiedName())
	})

	t.Run("with unset params", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			Unset: &WarehouseUnset{
				MaxClusterCount: Bool(true),
				AutoResume:      Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" UNSET MAX_CLUSTER_COUNT, AUTO_RESUME`)
	})

	t.Run("rename", func(t *testing.T) {
		newName := NewAccountObjectIdentifier("newName")
		opts := &AlterWarehouseOptions{
			name:    NewAccountObjectIdentifier("oldName"),
			NewName: &newName,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "oldName" RENAME TO "newName"`)
	})

	t.Run("suspend", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name:    NewAccountObjectIdentifier("mywarehouse"),
			Suspend: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" SUSPEND`)
	})

	t.Run("resume", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name:        NewAccountObjectIdentifier("mywarehouse"),
			Resume:      Bool(true),
			IfSuspended: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" RESUME IF SUSPENDED`)
	})

	t.Run("abort all queries", func(t *testing.T) {
		opts := &AlterWarehouseOptions{
			name:            NewAccountObjectIdentifier("mywarehouse"),
			AbortAllQueries: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" ABORT ALL QUERIES`)
	})

	t.Run("with set tag", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			SetTag: []TagAssociation{
				{
					Name:  tagId1,
					Value: "v1",
				},
				{
					Name:  tagId2,
					Value: "v2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" SET TAG %s = 'v1', %s = 'v2'`, tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("with unset tag", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		opts := &AlterWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
			UnsetTag: []ObjectIdentifier{
				tagId1,
				tagId2,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER WAREHOUSE "mywarehouse" UNSET TAG %s, %s`, tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})
}

func TestWarehouseDrop(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &DropWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP WAREHOUSE "mywarehouse"`)
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := &DropWarehouseOptions{
			name:     NewAccountObjectIdentifier("mywarehouse"),
			IfExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP WAREHOUSE IF EXISTS "mywarehouse"`)
	})
}

func TestWarehouseShow(t *testing.T) {
	t.Run("empty options", func(t *testing.T) {
		opts := &ShowWarehouseOptions{}
		assertOptsValidAndSQLEquals(t, opts, `SHOW WAREHOUSES`)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowWarehouseOptions{
			Like: &Like{
				Pattern: String("mywarehouse"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW WAREHOUSES LIKE 'mywarehouse'`)
	})
}

func TestWarehouseDescribe(t *testing.T) {
	t.Run("only name", func(t *testing.T) {
		opts := &describeWarehouseOptions{
			name: NewAccountObjectIdentifier("mywarehouse"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE WAREHOUSE "mywarehouse"`)
	})
}

func TestWarehouseToResourceConstraint(t *testing.T) {
	type test struct {
		input string
		want  WarehouseResourceConstraint
	}

	valid := []test{
		// Case insensitive.
		{input: "standard_gen_1", want: WarehouseResourceConstraintStandardGen1},

		// Supported Values.
		{input: "STANDARD_GEN_1", want: WarehouseResourceConstraintStandardGen1},
		{input: "STANDARD_GEN_2", want: WarehouseResourceConstraintStandardGen2},
		{input: "MEMORY_1X", want: WarehouseResourceConstraintMemory1X},
		{input: "MEMORY_1X_X86", want: WarehouseResourceConstraintMemory1Xx86},
		{input: "MEMORY_16X", want: WarehouseResourceConstraintMemory16X},
		{input: "MEMORY_16X_X86", want: WarehouseResourceConstraintMemory16Xx86},
		{input: "MEMORY_64X", want: WarehouseResourceConstraintMemory64X},
		{input: "MEMORY_64X_X86", want: WarehouseResourceConstraintMemory64Xx86},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToWarehouseResourceConstraint(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToWarehouseResourceConstraint(tc.input)
			require.Error(t, err)
		})
	}
}

func TestWarehouseToResourceConstraintWithoutGeneration(t *testing.T) {
	type test struct {
		input string
		want  WarehouseResourceConstraint
	}

	valid := []test{
		// Case insensitive.
		{input: "memory_1x", want: WarehouseResourceConstraintMemory1X},
		// Supported Values.
		{input: "MEMORY_1X", want: WarehouseResourceConstraintMemory1X},
		{input: "MEMORY_1X_X86", want: WarehouseResourceConstraintMemory1Xx86},
		{input: "MEMORY_16X", want: WarehouseResourceConstraintMemory16X},
		{input: "MEMORY_16X_X86", want: WarehouseResourceConstraintMemory16Xx86},
		{input: "MEMORY_64X", want: WarehouseResourceConstraintMemory64X},
		{input: "MEMORY_64X_X86", want: WarehouseResourceConstraintMemory64Xx86},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
		{input: "STANDARD_GEN_1"},
		{input: "STANDARD_GEN_2"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToWarehouseResourceConstraintWithoutGeneration(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToWarehouseResourceConstraintWithoutGeneration(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_Warehouse_ToWarehouseSize(t *testing.T) {
	type test struct {
		input string
		want  WarehouseSize
	}

	valid := []test{
		// case insensitive.
		{input: "xsmall", want: WarehouseSizeXSmall},

		// Supported Values
		{input: "XSMALL", want: WarehouseSizeXSmall},
		{input: "SMALL", want: WarehouseSizeSmall},
		{input: "MEDIUM", want: WarehouseSizeMedium},
		{input: "LARGE", want: WarehouseSizeLarge},
		{input: "XLARGE", want: WarehouseSizeXLarge},
		{input: "XXLARGE", want: WarehouseSizeXXLarge},
		{input: "XXXLARGE", want: WarehouseSizeXXXLarge},
		{input: "X4LARGE", want: WarehouseSizeX4Large},
		{input: "X5LARGE", want: WarehouseSizeX5Large},
		{input: "X6LARGE", want: WarehouseSizeX6Large},

		// Synonyms
		{input: "X-SMALL", want: WarehouseSizeXSmall},
		{input: "X-LARGE", want: WarehouseSizeXLarge},
		{input: "X2LARGE", want: WarehouseSizeXXLarge},
		{input: "2X-LARGE", want: WarehouseSizeXXLarge},
		{input: "X3LARGE", want: WarehouseSizeXXXLarge},
		{input: "3X-LARGE", want: WarehouseSizeXXXLarge},
		{input: "4X-LARGE", want: WarehouseSizeX4Large},
		{input: "5X-LARGE", want: WarehouseSizeX5Large},
		{input: "6X-LARGE", want: WarehouseSizeX6Large},
	}

	invalid := []test{
		// old values
		{input: "2XLARGE"},
		{input: "3XLARGE"},
		{input: "4XLARGE"},
		{input: "5XLARGE"},
		{input: "6XLARGE"},

		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToWarehouseSize(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToWarehouseSize(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_Warehouse_ToWarehouseType(t *testing.T) {
	type test struct {
		input string
		want  WarehouseType
	}

	valid := []test{
		// case insensitive.
		{input: "standard", want: WarehouseTypeStandard},

		// Supported Values
		{input: "STANDARD", want: WarehouseTypeStandard},
		{input: "SNOWPARK-OPTIMIZED", want: WarehouseTypeSnowparkOptimized},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},

		// not supported values (single-quoted)
		{input: "'STANDARD'"},
		{input: "'SNOWPARK-OPTIMIZED'"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToWarehouseType(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToWarehouseType(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_Warehouse_ToScalingPolicy(t *testing.T) {
	type test struct {
		input string
		want  ScalingPolicy
	}

	valid := []test{
		// case insensitive.
		{input: "standard", want: ScalingPolicyStandard},

		// Supported Values
		{input: "STANDARD", want: ScalingPolicyStandard},
		{input: "ECONOMY", want: ScalingPolicyEconomy},
	}

	invalid := []test{
		// bad values
		{input: ""},
		{input: "foo"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToScalingPolicy(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToScalingPolicy(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_Warehouse_Convert(t *testing.T) {
	correctRow := func() warehouseDBRow {
		return warehouseDBRow{
			Size:      string(WarehouseSizeXSmall),
			Available: "100",
		}
	}

	t.Run("convert error: size invalid", func(t *testing.T) {
		row := correctRow()
		row.Size = "INCORRECT_SIZE"

		wh, err := row.convert()

		require.ErrorContains(t, err, "invalid warehouse size: INCORRECT_SIZE")
		require.Nil(t, wh)
	})

	t.Run("convert error: available NaN", func(t *testing.T) {
		row := correctRow()
		row.Available = "not a number"

		wh, err := row.convert()

		require.ErrorContains(t, err, "row 'available' has incorrect value 'not a number'")
		require.Nil(t, wh)
	})

	t.Run("convert correct", func(t *testing.T) {
		row := correctRow()

		wh, err := row.convert()

		require.NoError(t, err)
		require.NotNil(t, wh)
		assert.Equal(t, WarehouseSizeXSmall, wh.Size)
		assert.InDelta(t, 100.0, wh.Available, testvars.FloatEpsilon)
	})

	t.Run("convert correct: available empty", func(t *testing.T) {
		row := correctRow()
		row.Available = " "

		wh, err := row.convert()

		require.NoError(t, err)
		require.NotNil(t, wh)
		assert.InDelta(t, 0.0, wh.Available, testvars.FloatEpsilon)
	})
}

func TestIsWarehouseResourceConstraintForSnowparkOptimized(t *testing.T) {
	trueCases := []WarehouseResourceConstraint{
		WarehouseResourceConstraintMemory1X,
		WarehouseResourceConstraintMemory1Xx86,
		WarehouseResourceConstraintMemory16X,
		WarehouseResourceConstraintMemory16Xx86,
		WarehouseResourceConstraintMemory64X,
		WarehouseResourceConstraintMemory64Xx86,
	}

	falseCases := []WarehouseResourceConstraint{
		WarehouseResourceConstraintStandardGen1,
		WarehouseResourceConstraintStandardGen2,
		WarehouseResourceConstraint("UNKNOWN"),
	}

	for _, c := range trueCases {
		t.Run(string(c), func(t *testing.T) {
			assert.True(t, IsWarehouseResourceConstraintForSnowparkOptimized(c))
		})
	}

	for _, c := range falseCases {
		t.Run(string(c), func(t *testing.T) {
			assert.False(t, IsWarehouseResourceConstraintForSnowparkOptimized(c))
		})
	}
}

func Test_Warehouse_ToWarehouseGeneration(t *testing.T) {
	type test struct {
		input string
		want  WarehouseGeneration
	}

	valid := []test{
		{input: "1", want: WarehouseGenerationStandardGen1},
		{input: "2", want: WarehouseGenerationStandardGen2},
	}

	invalid := []string{
		"",
		"0",
		"GEN_1",
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToWarehouseGeneration(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, in := range invalid {
		t.Run(in, func(t *testing.T) {
			_, err := ToWarehouseGeneration(in)
			require.Error(t, err)
		})
	}
}

func Test_Warehouse_WarehouseGenerationToWarehouseResourceConstraint(t *testing.T) {
	type test struct {
		input WarehouseGeneration
		want  WarehouseResourceConstraint
	}

	valid := []test{
		{input: "1", want: WarehouseResourceConstraintStandardGen1},
		{input: "2", want: WarehouseResourceConstraintStandardGen2},
	}

	invalid := []string{
		"",
		"0",
		"GEN_1",
	}

	for _, tc := range valid {
		t.Run(string(tc.input), func(t *testing.T) {
			got, err := WarehouseGenerationToWarehouseResourceConstraint(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, in := range invalid {
		t.Run(in, func(t *testing.T) {
			_, err := WarehouseGenerationToWarehouseResourceConstraint(WarehouseGeneration(in))
			require.Error(t, err)
		})
	}
}

func Test_Warehouse_WarehouseResourceConstraintToWarehouseGeneration(t *testing.T) {
	type test struct {
		input WarehouseResourceConstraint
		want  WarehouseGeneration
	}

	valid := []test{
		{input: WarehouseResourceConstraintStandardGen1, want: WarehouseGenerationStandardGen1},
		{input: WarehouseResourceConstraintStandardGen2, want: WarehouseGenerationStandardGen2},
	}

	invalid := []string{
		"",
		"0",
		"1",
	}

	for _, tc := range valid {
		t.Run(string(tc.input), func(t *testing.T) {
			got, err := WarehouseResourceConstraintToWarehouseGeneration(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, in := range invalid {
		t.Run(in, func(t *testing.T) {
			_, err := WarehouseResourceConstraintToWarehouseGeneration(WarehouseResourceConstraint(in))
			require.Error(t, err)
		})
	}
}
