package sdk

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	id := warehousesTestIdAccountObjectIdentifier

	tagId1 := randomSchemaObjectIdentifier()
	tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
	resourceMonitorId := randomAccountObjectIdentifier()
	tableId1 := randomSchemaObjectIdentifier()
	tableId2 := randomSchemaObjectIdentifierInSchema(tableId1.SchemaId())
	renameTarget := randomAccountObjectIdentifier()

	warehousesTests.Create.
		withAdditionalValidationCase(
			"validation_Create_WarehouseType_regularOnly",
			func(opts *CreateWarehouseOptions) { opts.WarehouseType = Pointer(WarehouseTypeAdaptive) },
			fmt.Errorf("only STANDARD, SNOWPARK-OPTIMIZED warehouses are supported, got ADAPTIVE"),
		).
		withAdditionalValidationCase(
			"validation_Create_MinClusterCount_lessOrEqualMaxClusterCount",
			func(opts *CreateWarehouseOptions) { opts.MinClusterCount = Int(2); opts.MaxClusterCount = Int(1) },
			fmt.Errorf("MinClusterCount must be less than or equal to MaxClusterCount"),
		).
		withExpectedSqlf(
			case_Warehouses_sql_Create_basic,
			`CREATE WAREHOUSE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Create_all,
			func(opts *CreateWarehouseOptions) {
				opts.OrReplace = Bool(true)
				opts.WarehouseType = Pointer(WarehouseTypeStandard)
				opts.WarehouseSize = Pointer(WarehouseSizeX4Large)
				opts.MaxClusterCount = Int(8)
				opts.MinClusterCount = Int(3)
				opts.ScalingPolicy = Pointer(ScalingPolicyEconomy)
				opts.AutoSuspend = Int(1000)
				opts.AutoResume = Bool(true)
				opts.InitiallySuspended = Bool(false)
				opts.ResourceMonitor = Pointer(resourceMonitorId)
				opts.Comment = String("hello")
				opts.EnableQueryAcceleration = Bool(true)
				opts.QueryAccelerationMaxScaleFactor = Int(62)
				opts.ResourceConstraint = Pointer(WarehouseResourceConstraintMemory1X)
				opts.Generation = Pointer(WarehouseGenerationStandardGen1)
				opts.MaxConcurrencyLevel = Int(7)
				opts.StatementQueuedTimeoutInSeconds = Int(29)
				opts.StatementTimeoutInSeconds = Int(89)
				opts.Tag = []TagAssociation{
					{Name: tagId1, Value: "v1"},
					{Name: tagId2, Value: "v2"},
				}
			},
			`CREATE OR REPLACE WAREHOUSE %s WAREHOUSE_TYPE = 'STANDARD' WAREHOUSE_SIZE = 'X4LARGE' MAX_CLUSTER_COUNT = 8 MIN_CLUSTER_COUNT = 3 SCALING_POLICY = 'ECONOMY' AUTO_SUSPEND = 1000 AUTO_RESUME = true INITIALLY_SUSPENDED = false RESOURCE_MONITOR = %s COMMENT = 'hello' ENABLE_QUERY_ACCELERATION = true QUERY_ACCELERATION_MAX_SCALE_FACTOR = 62 RESOURCE_CONSTRAINT = 'MEMORY_1X' GENERATION = '1' MAX_CONCURRENCY_LEVEL = 7 STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 29 STATEMENT_TIMEOUT_IN_SECONDS = 89 TAG (%s = 'v1', %s = 'v2')`,
			id.FullyQualifiedName(), resourceMonitorId.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		)

	warehousesTests.Create.
		withAdditionalSqlCasef(
			"sql_Create_onlyMinClusterCount",
			func(opts *CreateWarehouseOptions) { opts.MinClusterCount = Int(2) },
			`CREATE WAREHOUSE %s MIN_CLUSTER_COUNT = 2`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_maxEqualMin",
			func(opts *CreateWarehouseOptions) { opts.MaxClusterCount = Int(2); opts.MinClusterCount = Int(2) },
			`CREATE WAREHOUSE %s MAX_CLUSTER_COUNT = 2 MIN_CLUSTER_COUNT = 2`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_maxGreaterThanMin",
			func(opts *CreateWarehouseOptions) { opts.MaxClusterCount = Int(2); opts.MinClusterCount = Int(1) },
			`CREATE WAREHOUSE %s MAX_CLUSTER_COUNT = 2 MIN_CLUSTER_COUNT = 1`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_largeMinMaxValues",
			func(opts *CreateWarehouseOptions) { opts.MaxClusterCount = Int(100); opts.MinClusterCount = Int(11) },
			`CREATE WAREHOUSE %s MAX_CLUSTER_COUNT = 100 MIN_CLUSTER_COUNT = 11`, id.FullyQualifiedName(),
		)

	warehousesTests.CreateAdaptive.
		withExpectedSqlf(
			case_Warehouses_sql_CreateAdaptive_basic,
			`CREATE WAREHOUSE %s WAREHOUSE_TYPE = 'ADAPTIVE'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_CreateAdaptive_all,
			func(opts *CreateAdaptiveWarehouseOptions) {
				opts.OrReplace = Bool(true)
				opts.Comment = String("adaptive warehouse")
				opts.MaxQueryPerformanceLevel = Pointer(MaxQueryPerformanceLevelMedium)
				opts.QueryThroughputMultiplier = Int(22)
				opts.ResourceMonitor = Pointer(NewAccountObjectIdentifier("resmon"))
				opts.StatementQueuedTimeoutInSeconds = Int(30)
				opts.StatementTimeoutInSeconds = Int(60)
				opts.Tag = []TagAssociation{
					{Name: tagId1, Value: "v1"},
					{Name: tagId2, Value: "v2"},
				}
			},
			`CREATE OR REPLACE WAREHOUSE %s WAREHOUSE_TYPE = 'ADAPTIVE' COMMENT = 'adaptive warehouse' MAX_QUERY_PERFORMANCE_LEVEL = 'MEDIUM' QUERY_THROUGHPUT_MULTIPLIER = 22 RESOURCE_MONITOR = "resmon" TAG (%s = 'v1', %s = 'v2') STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 30 STATEMENT_TIMEOUT_IN_SECONDS = 60`,
			id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		)

	warehousesTests.CreateInteractive.
		withAdditionalValidationCase(
			"validation_CreateInteractive_MinClusterCount_lessOrEqualMaxClusterCount",
			func(opts *CreateInteractiveWarehouseOptions) {
				opts.MinClusterCount = Int(2)
				opts.MaxClusterCount = Int(1)
			},
			fmt.Errorf("MinClusterCount must be less than or equal to MaxClusterCount"),
		).
		withExpectedSqlf(
			case_Warehouses_sql_CreateInteractive_basic,
			`CREATE INTERACTIVE WAREHOUSE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_CreateInteractive_all,
			func(opts *CreateInteractiveWarehouseOptions) {
				opts.OrReplace = Bool(true)
				opts.Tables = []SchemaObjectIdentifier{tableId1, tableId2}
				opts.WarehouseSize = Pointer(WarehouseSizeXSmall)
				opts.MaxClusterCount = Int(2)
				opts.MinClusterCount = Int(1)
				opts.AutoSuspend = Int(86400)
				opts.AutoResume = Bool(true)
				opts.InitiallySuspended = Bool(true)
				opts.ResourceMonitor = Pointer(NewAccountObjectIdentifier("resmon"))
				opts.Comment = String("interactive warehouse")
				opts.MaxConcurrencyLevel = Int(8)
				opts.StatementQueuedTimeoutInSeconds = Int(30)
				opts.StatementTimeoutInSeconds = Int(5)
				opts.FallbackWarehouse = Pointer(NewAccountObjectIdentifier("fallbackwh"))
			},
			`CREATE OR REPLACE INTERACTIVE WAREHOUSE %s TABLES (%s, %s) WAREHOUSE_SIZE = 'XSMALL' MAX_CLUSTER_COUNT = 2 MIN_CLUSTER_COUNT = 1 AUTO_SUSPEND = 86400 AUTO_RESUME = true INITIALLY_SUSPENDED = true RESOURCE_MONITOR = "resmon" COMMENT = 'interactive warehouse' MAX_CONCURRENCY_LEVEL = 8 STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 30 STATEMENT_TIMEOUT_IN_SECONDS = 5 FALLBACK_WAREHOUSE = "fallbackwh"`,
			id.FullyQualifiedName(), tableId1.FullyQualifiedName(), tableId2.FullyQualifiedName(),
		)

	warehousesTests.Alter.
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Alter_Suspend,
			func(opts *AlterWarehouseOptions) { opts.Suspend = Bool(true) },
			`ALTER WAREHOUSE %s SUSPEND`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Alter_Resume,
			func(opts *AlterWarehouseOptions) { opts.Resume = Bool(true); opts.IfSuspended = Bool(true) },
			`ALTER WAREHOUSE %s RESUME IF SUSPENDED`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Alter_AbortAllQueries,
			func(opts *AlterWarehouseOptions) { opts.AbortAllQueries = Bool(true) },
			`ALTER WAREHOUSE %s ABORT ALL QUERIES`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Alter_RenameTo,
			func(opts *AlterWarehouseOptions) { opts.RenameTo = &renameTarget },
			`ALTER WAREHOUSE %s RENAME TO %s`, id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Alter_Set,
			func(opts *AlterWarehouseOptions) {
				opts.Set = &WarehouseSet{
					WarehouseType:                   Pointer(WarehouseTypeSnowparkOptimized),
					WaitForCompletion:               Bool(false),
					MinClusterCount:                 Int(4),
					MaxClusterCount:                 Int(5),
					AutoSuspend:                     Int(200),
					ResourceMonitor:                 Pointer(NewAccountObjectIdentifier("resmon")),
					EnableQueryAcceleration:         Bool(false),
					StatementQueuedTimeoutInSeconds: Int(1200),
					ResourceConstraint:              Pointer(WarehouseResourceConstraintMemory1X),
					Generation:                      Pointer(WarehouseGenerationStandardGen1),
				}
			},
			`ALTER WAREHOUSE %s SET WAREHOUSE_TYPE = 'SNOWPARK-OPTIMIZED' WAIT_FOR_COMPLETION = false MAX_CLUSTER_COUNT = 5 MIN_CLUSTER_COUNT = 4 AUTO_SUSPEND = 200 RESOURCE_MONITOR = "resmon" ENABLE_QUERY_ACCELERATION = false RESOURCE_CONSTRAINT = 'MEMORY_1X' GENERATION = '1' STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 1200`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Alter_Unset,
			func(opts *AlterWarehouseOptions) {
				opts.Unset = &WarehouseUnset{
					MaxClusterCount:    Bool(true),
					AutoResume:         Bool(true),
					ResourceConstraint: Bool(true),
					Generation:         Bool(true),
				}
			},
			`ALTER WAREHOUSE %s UNSET MAX_CLUSTER_COUNT, AUTO_RESUME, RESOURCE_CONSTRAINT, GENERATION`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Alter_AddTables,
			func(opts *AlterWarehouseOptions) { opts.AddTables = []SchemaObjectIdentifier{tableId1, tableId2} },
			`ALTER WAREHOUSE %s ADD TABLES (%s, %s)`, id.FullyQualifiedName(), tableId1.FullyQualifiedName(), tableId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Alter_DropTables,
			func(opts *AlterWarehouseOptions) { opts.DropTables = []SchemaObjectIdentifier{tableId1} },
			`ALTER WAREHOUSE %s DROP TABLES (%s)`, id.FullyQualifiedName(), tableId1.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Alter_SetTags,
			func(opts *AlterWarehouseOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tagId1, Value: "v1"},
					{Name: tagId2, Value: "v2"},
				}
			},
			`ALTER WAREHOUSE %s SET TAG %s = 'v1', %s = 'v2'`, id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Alter_UnsetTags,
			func(opts *AlterWarehouseOptions) { opts.UnsetTags = []ObjectIdentifier{tagId1} },
			`ALTER WAREHOUSE %s UNSET TAG %s`, id.FullyQualifiedName(), tagId1.FullyQualifiedName(),
		)

	warehousesTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_Set_FallbackWarehouse",
			func(opts *AlterWarehouseOptions) {
				opts.Set = &WarehouseSet{FallbackWarehouse: Pointer(NewAccountObjectIdentifier("fallbackwh"))}
			},
			`ALTER WAREHOUSE %s SET FALLBACK_WAREHOUSE = "fallbackwh"`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_FallbackWarehouse",
			func(opts *AlterWarehouseOptions) {
				opts.Unset = &WarehouseUnset{FallbackWarehouse: Bool(true)}
			},
			`ALTER WAREHOUSE %s UNSET FALLBACK_WAREHOUSE`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_adaptiveParams",
			func(opts *AlterWarehouseOptions) {
				opts.Set = &WarehouseSet{
					MaxQueryPerformanceLevel:        Pointer(MaxQueryPerformanceLevelXSmall),
					QueryThroughputMultiplier:       Int(5),
					StatementQueuedTimeoutInSeconds: Int(100),
					StatementTimeoutInSeconds:       Int(200),
				}
			},
			`ALTER WAREHOUSE %s SET QUERY_THROUGHPUT_MULTIPLIER = 5 MAX_QUERY_PERFORMANCE_LEVEL = 'XSMALL' STATEMENT_QUEUED_TIMEOUT_IN_SECONDS = 100 STATEMENT_TIMEOUT_IN_SECONDS = 200`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_adaptiveParams",
			func(opts *AlterWarehouseOptions) {
				opts.Unset = &WarehouseUnset{
					MaxQueryPerformanceLevel:        Bool(true),
					QueryThroughputMultiplier:       Bool(true),
					StatementQueuedTimeoutInSeconds: Bool(true),
					StatementTimeoutInSeconds:       Bool(true),
				}
			},
			`ALTER WAREHOUSE %s UNSET STATEMENT_QUEUED_TIMEOUT_IN_SECONDS, STATEMENT_TIMEOUT_IN_SECONDS, QUERY_THROUGHPUT_MULTIPLIER, MAX_QUERY_PERFORMANCE_LEVEL`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_UnsetTags_multiple",
			func(opts *AlterWarehouseOptions) { opts.UnsetTags = []ObjectIdentifier{tagId1, tagId2} },
			`ALTER WAREHOUSE %s UNSET TAG %s, %s`, id.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		)

	warehousesTests.Drop.
		withExpectedSqlf(
			case_Warehouses_sql_Drop_basic,
			`DROP WAREHOUSE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Drop_all,
			func(opts *DropWarehouseOptions) { opts.IfExists = Bool(true) },
			`DROP WAREHOUSE IF EXISTS %s`, id.FullyQualifiedName(),
		)

	warehousesTests.Show.
		withExpectedSql(case_Warehouses_sql_Show_basic, `SHOW WAREHOUSES`).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Show_all,
			func(opts *ShowWarehouseOptions) {
				opts.Like = &Like{Pattern: String("pattern")}
				opts.StartsWith = String("A")
				opts.Limit = &LimitFrom{Rows: Int(1), From: String("B")}
			},
			`SHOW WAREHOUSES LIKE 'pattern' STARTS WITH 'A' LIMIT 1 FROM 'B'`,
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Show_Like,
			func(opts *ShowWarehouseOptions) { opts.Like = &Like{Pattern: String("mywarehouse")} },
			`SHOW WAREHOUSES LIKE 'mywarehouse'`,
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Show_StartsWith,
			func(opts *ShowWarehouseOptions) { opts.StartsWith = String("A") },
			`SHOW WAREHOUSES STARTS WITH 'A'`,
		).
		withModifyAndExpectedSqlf(
			case_Warehouses_sql_Show_Limit,
			func(opts *ShowWarehouseOptions) { opts.Limit = &LimitFrom{Rows: Int(1), From: String("B")} },
			`SHOW WAREHOUSES LIMIT 1 FROM 'B'`,
		)

	warehousesTests.Describe.
		withExpectedSqlf(
			case_Warehouses_sql_Describe_basic,
			`DESCRIBE WAREHOUSE %s`, id.FullyQualifiedName(),
		)
}

func Test_Warehouse_ToWarehouseSize(t *testing.T) {
	type test struct {
		input string
		want  WarehouseSize
	}

	valid := []test{
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

func Test_Warehouse_ToWarehouseTypeUserSettable(t *testing.T) {
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
		// adaptive is not user-settable
		{input: "ADAPTIVE"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToWarehouseTypeUserSettable(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToWarehouseTypeUserSettable(tc.input)
			require.Error(t, err)
		})
	}
}

func Test_Warehouse_Convert(t *testing.T) {
	correctRow := func() warehouseDBRow {
		return warehouseDBRow{
			Size:      sql.NullString{String: string(WarehouseSizeXSmall), Valid: true},
			Available: "100",
		}
	}

	t.Run("convert error: size invalid", func(t *testing.T) {
		t.Skip("SNOW-3108659 - will return error from all convert mappings")
		row := correctRow()
		row.Size = sql.NullString{String: "INCORRECT_SIZE", Valid: true}

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
		assert.Equal(t, Pointer(WarehouseSizeXSmall), wh.Size)
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

	t.Run("convert correct: adaptive warehouse skips size parsing", func(t *testing.T) {
		row := correctRow()
		row.Type = string(WarehouseTypeAdaptive)
		row.Size = sql.NullString{Valid: false}

		wh, err := row.convert()

		require.NoError(t, err)
		require.NotNil(t, wh)
		assert.Equal(t, WarehouseTypeAdaptive, wh.Type)
		assert.Nil(t, wh.Size)
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
