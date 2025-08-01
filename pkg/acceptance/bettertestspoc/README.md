# Better tests poc

<!-- TOC -->
* [Better tests poc](#better-tests-poc)
  * [How it works](#how-it-works)
    * [Adding new resource assertions](#adding-new-resource-assertions)
    * [Adding new resource show output assertions](#adding-new-resource-show-output-assertions)
    * [Adding new resource parameters assertions](#adding-new-resource-parameters-assertions)
    * [Adding new Snowflake object assertions](#adding-new-snowflake-object-assertions)
    * [Adding new Snowflake object parameters assertions](#adding-new-snowflake-object-parameters-assertions)
    * [Adding new resource config model builders](#adding-new-resource-config-model-builders)
    * [Running the generators](#running-the-generators)
    * [Example usage in practice](#example-usage-in-practice)
  * [Known limitations/planned improvements](#known-limitationsplanned-improvements)
<!-- TOC -->

This package contains a quick implementation of helpers that should allow us a quicker, more pleasant, and more readable implementation of tests, mainly the acceptance ones.
It contains the following packages:
- `assert` - all the assertions reside here. Also, the utilities to build assertions for new objects. All the current assertions are generated. The currently supported assertions are:
  - Snowflake object assertions (generated in subpackage `objectassert`)
  - Snowflake object parameters assertions (generated in subpackage `objectparametersassert`)
  - resource assertions (generated in subpackage `resourceassert`)
  - resource parameters assertions (generated in subpackage `resourceparametersassert`)
  - show output assertions (generated in subpackage `resourceshowoutputassert`)

- `config` - the new model abstractions (`ResourceModel`, `DatasourceModel`, and `ProviderModel`) reside here. They provide models for objects and the builder methods allowing better config preparation in the acceptance tests.
It aims to be more readable than using `Config:` with hardcoded string or `ConfigFile:` for file that is not directly reachable from the test body. Also, it should be easier to reuse the models and prepare convenience extension methods. The models are already generated.

## How it works
### Adding new resource assertions
Resource assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allStructs` slice in the `assert/resourceassert/gen/resource_schema_def.go`
- to add custom (not generated assertions) create file `abc_resource_ext.go` in the `assert/resourceassert` package. Example would be:
```go
func (w *WarehouseResourceAssert) HasDefaultMaxConcurrencyLevel() *WarehouseResourceAssert {
    w.AddAssertion(assert.ValueSet("max_concurrency_level", "8"))
    return w
}
```

### Generated resource assertions
The generated resource assertions are available in the `assert/resourceassert/gen` package.
The generated assertions provide a set of methods that can be used to assert the state of the resource.
Here's a short description of every generated assertion:
- HasX methods check that a field is set with the provided value.
- HasNoX methods check that a field is not set (not present in state or is empty in case of lists/sets).
- HasXEmpty methods are generated only for string type fields and are checking if a given field is set to empty string.
- HasXNotEmpty methods check that a field is set with any value.

### Adding new resource show output assertions
Resource show output assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allResourceSchemaDefs` slice in the `assert/objectassert/gen/sdk_object_def.go`
- to add custom (not generated assertions) create file `abc_show_output_ext.go` in the `assert/resourceshowoutputassert` package. Example would be:
```go
func (u *UserShowOutputAssert) HasNameAndLoginName(expected string) *UserShowOutputAssert {
	return u.
		HasName(expected).
		HasLoginName(expected)
}
```

### Adding new resource parameters assertions
Resource parameters assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allObjectsParameters` slice in the `assert/objectparametersassert/gen/object_parameters_def.go`
- to add custom (not generated assertions) create file `warehouse_resource_parameters_ext.go` in the `assert/resourceparametersassert` package. Example would be:
```go
func (w *WarehouseResourceParametersAssert) HasDefaultMaxConcurrencyLevel() *WarehouseResourceParametersAssert {
	return w.
		HasMaxConcurrencyLevel(8).
		HasMaxConcurrencyLevelLevel("")
}
```

### Adding new Snowflake object assertions
Snowflake object assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allStructs` slice in the `assert/objectassert/gen/sdk_object_def.go`
- to add custom (not generated assertions) create file `abc_snowflake_ext.go` in the `objectassert` package. Example would be:
```go
func (w *WarehouseAssert) HasStateOneOf(expected ...sdk.WarehouseState) *WarehouseAssert {
    w.assertions = append(w.assertions, func(t *testing.T, o *sdk.Warehouse) error {
        t.Helper()
        if !slices.Contains(expected, o.State) {
            return fmt.Errorf("expected state one of: %v; got: %v", expected, string(o.State))
        }
        return nil
    })
    return w
}
```

### Adding new Snowflake object parameters assertions
Snowflake object parameters assertions can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allObjectsParameters` slice in the `assert/objectparametersassert/gen/main/main.go`
- make sure that test helper method `acc.TestClient().Parameter.ShowAbcParameters` exists in `/pkg/acceptance/helpers/parameter_client.go`
- to add custom (not generated) assertions create file `abc_parameters_snowflake_ext.go` in the `objectparametersassert` package. Example would be:
```go
func (w *WarehouseParametersAssert) HasDefaultMaxConcurrencyLevel() *WarehouseParametersAssert {
    return w.
        HasMaxConcurrencyLevel(8).
        HasMaxConcurrencyLevelLevel("")
}
```

### Adding new resource config model builders
Resource config model builders can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allResourceSchemaDefs` slice in the `assert/resourceassert/gen/resource_schema_def.go`
- to add custom (not generated) config builder methods create file `warehouse_model_ext` in the `config/model` package. Example would be:
```go
func BasicAbcModel(
	name string,
	comment string,
) *AbcModel {
	return AbcWithDefaultMeta(name).WithComment(comment)
}

func (w *AbcModel) WithWarehouseSizeEnum(warehouseSize sdk.WarehouseSize) *AbcModel {
	return w.WithWarehouseSize(string(warehouseSize))
}
```

### Adding new datasource config model builders
Data source config model builders can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allDatasourcesSchemaDefs` slice in the `config/datasourcemodel/gen/datasource_schema_def.go`
- to add custom (not generated) config builder methods create file `abc_model_ext` in the `config/datasourcemodel` package. Example would be:
```go
func BasicAbcModel(
	name string,
	comment string,
) *AbcModel {
	return AbcWithDefaultMeta(name).WithComment(comment)
}

func (d *AbcModel) WithLimit(rows int) *AbcModel {
    return d.WithLimitValue(
        tfconfig.ObjectVariable(map[string]tfconfig.Variable{
            "rows": tfconfig.IntegerVariable(rows),
        }),
    )
}
```

### Adding new provider config model builders
Provider config model builders can be generated automatically. For object `abc` do the following:
- add object you want to generate to `allProviderSchemaDefs` slice in the `config/providermodel/gen/provider_schema_def.go`
- to add custom (not generated) config builder methods create file `abc_model_ext` in the `config/providermodel` package. Example would be:
```go
func BasicAbcModel(
	name string,
	comment string,
) *AbcModel {
	return AbcWithDefaultMeta(name).WithComment(comment)
}

func (w *AbcModel) WithWarehouseSizeEnum(warehouseSize sdk.WarehouseSize) *AbcModel {
	return w.WithWarehouseSize(string(warehouseSize))
}
```

*Note*: our provider's config is already generated, so there should not be a need to generate any more providers (the regeneration or adding custom methods are still expected).

### Running the generators
Each of the above assertion types/config models has its own generator and cleanup entry in our Makefile.
You can generate config models with:
```shell
  make clean-resource-model-builders generate-resource-model-builders
```

You can use cli flags:
```shell
  make clean-resource-model-builders generate-resource-model-builders SF_TF_GENERATOR_ARGS='--dry-run --verbose'
```

To clean/generate all from this package run
```shell
  make clean-all-assertions-and-config-models generate-all-assertions-and-config-models
```

### Example usage in practice
You can check the current example usage in `TestAcc_Warehouse_BasicFlows` and the `create: complete` inside `TestInt_Warehouses`. To see the output after invalid assertions:
- add the following to the first step of `TestAcc_Warehouse_BasicFlows`
```go
    // bad checks below
    resourceassert.WarehouseResource(t, "snowflake_warehouse.w").
        HasWarehouseTypeString(string(sdk.WarehouseTypeSnowparkOptimized)).
        HasWarehouseSizeString(string(sdk.WarehouseSizeMedium)),
    resourceshowoutputassert.WarehouseShowOutput(t, "snowflake_warehouse.w").
        HasType(sdk.WarehouseTypeSnowparkOptimized),
    resourceparametersassert.WarehouseResourceParameters(t, "snowflake_warehouse.w").
        HasMaxConcurrencyLevel(16).
        HasMaxConcurrencyLevelLevel(sdk.ParameterTypeWarehouse),
    objectassert.Warehouse(t, warehouseId).
        HasName("bad name").
        HasState(sdk.WarehouseStateSuspended).
        HasType(sdk.WarehouseTypeSnowparkOptimized).
        HasSize(sdk.WarehouseSizeMedium).
        HasMaxClusterCount(12).
        HasMinClusterCount(13).
        HasScalingPolicy(sdk.ScalingPolicyEconomy).
        HasAutoSuspend(123).
        HasAutoResume(false).
        HasResourceMonitor(sdk.NewAccountObjectIdentifier("some-id")).
        HasComment("bad comment").
        HasEnableQueryAcceleration(true).
        HasQueryAccelerationMaxScaleFactor(12),
    objectparametersassert.WarehouseParameters(t, warehouseId).
        HasMaxConcurrencyLevel(16),
    assert.Check(resource.TestCheckResourceAttr("snowflake_warehouse.w", "warehouse_type", string(sdk.WarehouseTypeSnowparkOptimized))),
```
it will result in:
```
    warehouse_acceptance_test.go:51: Step 1/8 error: Check failed: check 7/12 error:
        snowflake_warehouse.w resource assertion [1/2]: failed with error: Attribute 'warehouse_type' not found
        snowflake_warehouse.w resource assertion [2/2]: failed with error: Attribute 'warehouse_size' not found
        check 8/12 error:
        snowflake_warehouse.w show_output assertion [2/2]: failed with error: Attribute 'show_output.0.type' expected "SNOWPARK-OPTIMIZED", got "STANDARD"
        check 9/12 error:
        snowflake_warehouse.w parameters assertion [2/3]: failed with error: Attribute 'parameters.0.max_concurrency_level.0.value' expected "16", got "8"
        snowflake_warehouse.w parameters assertion [3/3]: failed with error: Attribute 'parameters.0.max_concurrency_level.0.level' expected "WAREHOUSE", got ""
        check 10/12 error:
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [1/13]: failed with error: expected name: bad name; got: XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [2/13]: failed with error: expected state: SUSPENDED; got: STARTED
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [3/13]: failed with error: expected type: SNOWPARK-OPTIMIZED; got: STANDARD
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [4/13]: failed with error: expected size: MEDIUM; got: XSMALL
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [5/13]: failed with error: expected max cluster count: 12; got: 1
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [6/13]: failed with error: expected min cluster count: 13; got: 1
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [7/13]: failed with error: expected scaling policy: ECONOMY; got: STANDARD
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [8/13]: failed with error: expected auto suspend: 123; got: 600
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [9/13]: failed with error: expected auto resume: false; got: true
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [10/13]: failed with error: expected resource monitor: some-id; got:
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [11/13]: failed with error: expected comment: bad comment; got: Who does encouraging eagerly annoying dream several their scold straightaway.
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [12/13]: failed with error: expected enable query acceleration: true; got: false
        object WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"] assertion [13/13]: failed with error: expected query acceleration max scale factor: 12; got: 8
        check 11/12 error:
        parameter assertion for WAREHOUSE["XHZJCKAT_35D0BCC1_7797_974E_ACAF_C622C56FA2D2"][MAX_CONCURRENCY_LEVEL][1/1] failed: expected value 16, got 8
        check 12/12 error:
        snowflake_warehouse.w: Attribute 'warehouse_type' not found
```

- add the following to the second step of `TestAcc_Warehouse_BasicFlows`
```go
    // bad checks below
    assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(warehouseId.Name(), "bad name", name)),
    resourceassert.ImportedWarehouseResource(t, warehouseId.Name()).
        HasNameString("bad name").
        HasWarehouseTypeString(string(sdk.WarehouseTypeSnowparkOptimized)).
        HasWarehouseSizeString(string(sdk.WarehouseSizeMedium)).
        HasMaxClusterCountString("2").
        HasMinClusterCountString("3").
        HasScalingPolicyString(string(sdk.ScalingPolicyEconomy)).
        HasAutoSuspendString("123").
        HasAutoResumeString("false").
        HasResourceMonitorString("abc").
        HasCommentString("bad comment").
        HasEnableQueryAccelerationString("true").
        HasQueryAccelerationMaxScaleFactorString("16"),
    resourceparametersassert.ImportedWarehouseResourceParameters(t, warehouseId.Name()).
        HasMaxConcurrencyLevel(1).
        HasMaxConcurrencyLevelLevel(sdk.ParameterTypeWarehouse).
        HasStatementQueuedTimeoutInSeconds(23).
        HasStatementQueuedTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse).
        HasStatementTimeoutInSeconds(1232).
        HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse),
    objectassert.Warehouse(t, warehouseId).
        HasName("bad name").
        HasState(sdk.WarehouseStateSuspended).
        HasType(sdk.WarehouseTypeSnowparkOptimized).
        HasSize(sdk.WarehouseSizeMedium).
        HasMaxClusterCount(12).
        HasMinClusterCount(13).
        HasScalingPolicy(sdk.ScalingPolicyEconomy).
        HasAutoSuspend(123).
        HasAutoResume(false).
        HasResourceMonitor(sdk.NewAccountObjectIdentifier("some-id")).
        HasComment("bad comment").
        HasEnableQueryAcceleration(true).
        HasQueryAccelerationMaxScaleFactor(12),
    objectparametersassert.WarehouseParameters(t, warehouseId).
        HasMaxConcurrencyLevel(1),
```
it will result in:
```
    warehouse_acceptance_test.go:51: check 7/11 error:
        attribute bad name not found in instance state
        check 8/11 error:
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [1/12]: failed with error: expected: bad name, got: WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [2/12]: failed with error: expected: SNOWPARK-OPTIMIZED, got: STANDARD
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [3/12]: failed with error: expected: MEDIUM, got: XSMALL
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [4/12]: failed with error: expected: 2, got: 1
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [5/12]: failed with error: expected: 3, got: 1
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [6/12]: failed with error: expected: ECONOMY, got: STANDARD
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [7/12]: failed with error: expected: 123, got: 600
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [8/12]: failed with error: expected: false, got: true
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [9/12]: failed with error: expected: abc, got:
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [10/12]: failed with error: expected: bad comment, got: Promise my huh off certain you bravery dynasty with Roman.
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [11/12]: failed with error: expected: true, got: false
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported resource assertion [12/12]: failed with error: expected: 16, got: 8
        check 9/11 error:
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [2/7]: failed with error: expected: 1, got: 8
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [3/7]: failed with error: expected: WAREHOUSE, got:
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [4/7]: failed with error: expected: 23, got: 0
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [5/7]: failed with error: expected: WAREHOUSE, got:
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [6/7]: failed with error: expected: 1232, got: 172800
        WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65 imported parameters assertion [7/7]: failed with error: expected: WAREHOUSE, got:
        check 10/11 error:
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [1/13]: failed with error: expected name: bad name; got: WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [2/13]: failed with error: expected state: SUSPENDED; got: STARTED
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [3/13]: failed with error: expected type: SNOWPARK-OPTIMIZED; got: STANDARD
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [4/13]: failed with error: expected size: MEDIUM; got: XSMALL
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [5/13]: failed with error: expected max cluster count: 12; got: 1
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [6/13]: failed with error: expected min cluster count: 13; got: 1
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [7/13]: failed with error: expected scaling policy: ECONOMY; got: STANDARD
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [8/13]: failed with error: expected auto suspend: 123; got: 600
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [9/13]: failed with error: expected auto resume: false; got: true
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [10/13]: failed with error: expected resource monitor: some-id; got:
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [11/13]: failed with error: expected comment: bad comment; got: Promise my huh off certain you bravery dynasty with Roman.
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [12/13]: failed with error: expected enable query acceleration: true; got: false
        object WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"] assertion [13/13]: failed with error: expected query acceleration max scale factor: 12; got: 8
        check 11/11 error:
        parameter assertion for WAREHOUSE["WBJKHLAT_2E52D1E6_D23D_33A0_F568_4EBDBE083B65"][MAX_CONCURRENCY_LEVEL][1/1] failed: expected value 1, got 8
```

- add the following to the `create: complete` in `TestInt_Warehouses`:
```go
    // to show errors
    warehouseAssertionsBad := objectassert.Warehouse(t, id).
        HasName("bad name").
        HasState(sdk.WarehouseStateSuspended).
        HasType(sdk.WarehouseTypeSnowparkOptimized).
        HasSize(sdk.WarehouseSizeMedium).
        HasMaxClusterCount(12).
        HasMinClusterCount(13).
        HasScalingPolicy(sdk.ScalingPolicyStandard).
        HasAutoSuspend(123).
        HasAutoResume(false).
        HasResourceMonitor(sdk.NewAccountObjectIdentifier("some-id")).
        HasComment("bad comment").
        HasEnableQueryAcceleration(false).
        HasQueryAccelerationMaxScaleFactor(12)
    assertions.AssertThatObject(t, warehouseAssertionsBad, testClientHelper())
```
it will result in:
```
    commons.go:101:
        	Error Trace:	/Users/asawicki/Projects/terraform-provider-snowflake/pkg/sdk/testint/warehouses_integration_test.go:149
        	Error:      	Received unexpected error:
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [1/13]: failed with error: expected name: bad name; got: VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [2/13]: failed with error: expected state: SUSPENDED; got: STARTED
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [3/13]: failed with error: expected type: SNOWPARK-OPTIMIZED; got: STANDARD
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [4/13]: failed with error: expected size: MEDIUM; got: SMALL
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [5/13]: failed with error: expected max cluster count: 12; got: 8
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [6/13]: failed with error: expected min cluster count: 13; got: 2
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [7/13]: failed with error: expected type: STANDARD; got: ECONOMY
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [8/13]: failed with error: expected auto suspend: 123; got: 1000
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [9/13]: failed with error: expected auto resume: false; got: true
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [10/13]: failed with error: expected resource monitor: some-id; got: OOUJMDIT_535F314F_6549_348F_370E_AB430EE4BC7B
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [11/13]: failed with error: expected comment: bad comment; got: comment
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [12/13]: failed with error: expected enable query acceleration: false; got: true
        	            	object WAREHOUSE["VKSENEIT_535F314F_6549_348F_370E_AB430EE4BC7B"] assertion [13/13]: failed with error: expected query acceleration max scale factor: 12; got: 90
        	Test:       	TestInt_Warehouses/create:_complete
```

## Planned improvements
- Unify assertions (resource, resource_show, resource_parameters, etc.), so that the assertions, generated methods, and helper methods used inside them are consistent across assertion types
- Test all the utilities for assertion/model construction (public interfaces, methods, functions).
- Verify if all the config types are supported.
- Consider a better implementation for the model conversion to config (TODO left in `config/config.go`).
- Support additional methods for references in models (TODO left in `config/config.go`).
- Generate depends_on for all compatible models. Consider exporting it in meta (discussion: https://github.com/snowflakedb/terraform-provider-snowflake/pull/3207#discussion_r1850053618).
- Generate timeouts for all compatible models, similarly to depends_on.
- Add a convenience function to concatenate multiple models (TODO left in `config/config.go`).
- Add function to support using `ConfigFile:` in the acceptance tests (TODO left in `config/config.go`).
- Replace `acceptance/snowflakechecks` with the new proposed Snowflake objects assertions.
- Support `showOutputValueUnset` and add a second function for each `show_output` attribute.
- Support `resourceAssertionTypeValueNotSet` for import checks (`panic` left currently).
- Add assertions for
  - The `describe_output` and handle describe objects too.
  - Custom objects: functions using `NewSnowflakeObjectAssertWithTestClientObjectProvider` can handle only objects that have proper identifiers. Not all of them have such IDs (check user PAT), so this function can't be used.
- Add support for datasource tests (assertions and config builders).
- Consider overriding the assertions when invoking same check multiple times with different params (e.g. `Warehouse(...).HasType(X).HasType(Y)`; it could use the last-check-wins approach, to more easily reuse complex checks between the test steps).
- Consider not adding the check for `show_output` presence on creation (same with `parameters`). The majority of the use cases need it to be present but there are a few others (like conditional presence in the datasources). Currently, it seems that they should be always present in the resources, so no change is made. Later, with adding the support for the datasource tests, consider simple destructive implementation like:
- Add support for `set` so that assertions like e.g. `oauth_scopes.*` could be done.
```go
func (w *WarehouseDatasourceShowOutputAssert) IsEmpty() {
    w.assertions = make([]resourceAssertion, 0)
    w.assertions = append(w.assertions, valueSet("show_output.#", "0"))
}
```
- add an assertion method or helper method that would be better at checking for missing complex types (like Lists, Sets, Maps, etc.), currently ValueSet with `.#` is used
- support other mappings if needed (TODO left in `assert/objectassert/gen/model.go`)
- consider extracting preamble model to commons (TODOs left in `assert/objectassert/gen/model.go` and in `assert/objectparametersassert/gen/model.go`)
- get a runtime name for the assertion creator (TODOs left in `assert/objectparametersassert/gen/model.go`)
- use a better definition for each objet's snowflake parameters (TODO left in `assert/objectparametersassert/gen/main/main.go`)
- add possibility to have enums generated in config builders (TODO left in `config/model/warehouse_model_ext.go`)
- handle situations where snowflake default behaves inconsistently (TODO left in `assert/objectparametersassert/gen/object_parameters_def.go`)
- handle attribute types in resource assertions (currently strings only; TODO left in `assert/resourceassert/gen/model.go`)
- distinguish between different enum types (TODO left in `assert/resourceshowoutputassert/gen/templates.go`)
- support the rest of attribute types in config model builders (TODO left in `config/model/gen/model.go`)
- Omit computed fields in the model (like FullyQualifiedName), because it doesn't make sense to set them
- There's an error when generating models, steps to reproduce:
  - Go to view resource code and change `data_metric_function` field to `testing` and make it required
  - During the generation, the following error appears: mixed named and unnamed parameters.
    It's a golang error indicating that the parameter has both unnamed and named parameters in function (e.g. `func(abc string, int)`).
    The error is a result of both things:
    1. Lists of objects are partially generated, and only parameter name is generated in some functions (the type has to be added manually).
    2. `testing` is a package name that makes Go think that we want to have unnamed parameter there, but we just didn't generate the type for that field in the function argument.
- generate assertions checking that time is not empty - we often do not compare time fields by value, but check if they are set
- utilize `ContainsExactlyInAnyOrder` function in `pkg/acceptance/bettertestspoc/assert/commons.go` to create asserts on collections that are order independent
- Additional asserts for sets and lists that wouldn't rely on the order of items saved to the state (SNOW-1706544)
  - this should also support nested sets and lists (see `accountRolesDataSourceContainsRole` for example)
- support generating provider config and use generated configs in `pkg/provider/provider_acceptance_test.go`
- add config builders for other block types (Variable, Output, Locals, Module, Terraform)
- add provider to resource/datasource models (use in the grant_ownership_acceptance_test)
- explore HCL v2 in more detail (especially struct tags generation; probably with migration to plugin framework because of schema models); ref: https://github.com/hashicorp/hcl/blob/bee2dc2e75f7528ad85777b7a013c13796426bd6/gohcl/encode_test.go#L48
- introduce some common interface for all three existing models (ResourceModel, DatasourceModel, and ProviderModel)
- rename ResourceSchemaDetails (because it is used for the datasources and provider too)
- consider duplicating the builders template from resource (currently same template used for datasources and provider which limits the customization possibilities for just one block type)
- consider merging ResourceModel with DatasourceModel (currently the implementation is really similar)
- remove schema.TypeMap workaround or make it wiser (e.g. during generation we could programmatically gather all schema.TypeMap and use this workaround only for them)
- support asserting resource id in `assert/resourceassert/*_gen.go`
- add possibility for object parameter assert not take any identifier (currently there's a workaround in `account_parameters_snowflake_gen.go`, because `SHOW PARAMETERS FOR ACCOUNT` don't take any identifiers)
- SNOW-2048330: add possibility to override the default value (and optionally default level) used in HasAllDefaults -> HasDefaultParameterValueOnLevel parameter assertions (it's required for cases like asserting `NETWORK_POLICY` which is predefined for our testing environments and causes test failures)
- `Mapper` is just a `func(string) string`, so there is no easy way inside the template to know what mapper is being applied. Because of that, we have the `Identity` mapper which just returns the input string which leads to easier template logic applicable both to cases needed the mapping and not needing it. It leads to more unnecessary code (like for collections comparison in object asserts), so it would be great to change the logic, e.g. by handling the list of mappers (logic allowing the emptiness check, easier piping).
- Do not generate `Has...String` for complex types, like maps and lists. Instead, generate correct assertions for checking all the values in such type.

## Known limitations
- generating provider config may misbehave when used only with one object/map paramter (like `params`), e.g.:
```json
{
  "provider": {
    "snowflake": {
      "params": {
        "statement_timeout_in_seconds": 31337
      }
    }
  }
}
```
will be converted to HCL:
```hcl
provider "snowflake" "params" {
  statement_timeout_in_seconds = 31337
}
```
