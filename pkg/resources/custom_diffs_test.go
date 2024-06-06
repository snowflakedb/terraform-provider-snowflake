package resources_test

import (
	"context"
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNestedValueComputedIf(t *testing.T) {
	customDiff := resources.NestedValueComputedIf(
		"nested_value",
		func(client *sdk.Client) (*sdk.Parameter, error) {
			return &sdk.Parameter{
				Key:   "Parameter",
				Value: "snow-value",
			}, nil
		},
		func(v any) string { return v.(string) },
	)
	providerConfig := createProviderWithNestedValueAndCustomDiff(t, schema.TypeString, customDiff)

	t.Run("value set in the configuration and state", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"nested_value": cty.ListVal([]cty.Value{
				cty.MapVal(map[string]cty.Value{
					"value": cty.NumberIntVal(123),
				}),
			}),
		}), map[string]any{
			"nested_value": []any{
				map[string]any{
					"value": 123,
				},
			},
		})
		assert.False(t, diff.Attributes["nested_value.#"].NewComputed)
	})

	t.Run("value set only in the configuration", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"nested_value": cty.ListVal([]cty.Value{
				cty.MapVal(map[string]cty.Value{
					"value": cty.NumberIntVal(123),
				}),
			}),
		}), map[string]any{})
		assert.True(t, diff.Attributes["nested_value.#"].NewComputed)
	})

	t.Run("value set in the state and not equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"nested_value": []any{
				map[string]any{
					"value": "value-to-change",
				},
			},
		})
		assert.True(t, diff.Attributes["nested_value.#"].NewComputed)
	})

	t.Run("value set in the state and equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"nested_value": []any{
				map[string]any{
					"value": "snow-value",
				},
			},
		})
		assert.False(t, diff.Attributes["nested_value.#"].NewComputed)
	})
}

func TestNestedIntValueAccountObjectComputedIf(t *testing.T) {
	providerConfig := createProviderWithNestedValueAndCustomDiff(t, schema.TypeInt, resources.NestedIntValueAccountObjectComputedIf("nested_value", sdk.AccountParameterDataRetentionTimeInDays))

	t.Run("different value than on the Snowflake side", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"nested_value": []any{
				map[string]any{
					"value": 999, // value outside of valid range
				},
			},
		})
		assert.True(t, diff.Attributes["nested_value.#"].NewComputed)
	})

	t.Run("same value as in Snowflake", func(t *testing.T) {
		dataRetentionTimeInDays, err := acc.Client(t).Parameters.ShowAccountParameter(context.Background(), sdk.AccountParameterDataRetentionTimeInDays)
		require.NoError(t, err)

		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"nested_value": []any{
				map[string]any{
					"value": dataRetentionTimeInDays.Value,
				},
			},
		})
		assert.False(t, diff.Attributes["nested_value.#"].NewComputed)
	})
}

func createProviderWithNestedValueAndCustomDiff(t *testing.T, valueType schema.ValueType, customDiffFunc schema.CustomizeDiffFunc) *schema.Provider {
	t.Helper()
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test": {
				Schema: map[string]*schema.Schema{
					"nested_value": {
						Type:     schema.TypeList,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"value": {
									Type:     valueType,
									Required: true,
								},
							},
						},
						Computed: true,
						Optional: true,
					},
				},
				CustomizeDiff: customDiffFunc,
			},
		},
	}
}

func calculateDiff(t *testing.T, providerConfig *schema.Provider, rawConfigValue cty.Value, stateValue map[string]any) *terraform.InstanceDiff {
	t.Helper()
	diff, err := providerConfig.ResourcesMap["test"].Diff(
		context.Background(),
		&terraform.InstanceState{
			RawConfig: rawConfigValue,
		},
		&terraform.ResourceConfig{
			Config: stateValue,
		},
		&provider.Context{Client: acc.Client(t)},
	)
	require.NoError(t, err)
	return diff
}

func Test_NormalizeAndCompare(t *testing.T) {
	genericNormalize := func(value string) (any, error) {
		switch value {
		case "ok", "ok1":
			return "ok", nil
		default:
			return nil, fmt.Errorf("incorrect value %s", value)
		}
	}

	t.Run("generic normalize", func(t *testing.T) {
		result := resources.NormalizeAndCompare(genericNormalize)("", "ok", "ok", nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(genericNormalize)("", "ok", "ok1", nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(genericNormalize)("", "ok", "nok", nil)
		assert.False(t, result)
	})

	t.Run("warehouse size", func(t *testing.T) {
		result := resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), string(sdk.WarehouseSizeX4Large), nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), "4X-LARGE", nil)
		assert.True(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), string(sdk.WarehouseSizeX5Large), nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), "invalid", nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", string(sdk.WarehouseSizeX4Large), "", nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", "invalid", string(sdk.WarehouseSizeX4Large), nil)
		assert.False(t, result)

		result = resources.NormalizeAndCompare(sdk.ToWarehouseSize)("", "", string(sdk.WarehouseSizeX4Large), nil)
		assert.False(t, result)
	})
}
