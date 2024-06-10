package resources_test

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestValueComputedIf(t *testing.T) {
	customDiff := resources.ValueComputedIf[string](
		"value",
		[]*sdk.Parameter{
			{
				Key:   string(sdk.AccountParameterLogLevel),
				Value: string(sdk.LogLevelInfo),
			},
		},
		sdk.AccountParameterLogLevel,
		func(v any) string { return v.(string) },
		func(v string) (string, error) { return v, nil },
	)
	providerConfig := createProviderWithValuePropertyAndCustomDiff(t, schema.TypeString, customDiff)

	t.Run("value set in the configuration and state", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("value set only in the configuration", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{})
		assert.True(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("value set in the state and not equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"value": string(sdk.LogLevelDebug),
		})
		assert.Equal(t, string(sdk.LogLevelInfo), diff.Attributes["value"].New)
	})

	t.Run("value set in the state and equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})
}

func TestAccountObjectStringValueComputedIf(t *testing.T) {
	customDiff := resources.AccountObjectStringValueComputedIf(
		"value",
		[]*sdk.Parameter{
			{
				Key:   string(sdk.AccountParameterLogLevel),
				Value: string(sdk.LogLevelInfo),
			},
		},
		sdk.AccountParameterLogLevel,
	)
	providerConfig := createProviderWithValuePropertyAndCustomDiff(t, schema.TypeString, customDiff)

	t.Run("value set in the configuration and state", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("value set only in the configuration", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{})
		assert.True(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("value set in the state and not equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"value": string(sdk.LogLevelDebug),
		})
		assert.Equal(t, string(sdk.LogLevelInfo), diff.Attributes["value"].New)
	})

	t.Run("value set in the state and equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})
}

func TestAccountObjectIntValueComputedIf(t *testing.T) {
	customDiff := resources.AccountObjectIntValueComputedIf(
		"value",
		[]*sdk.Parameter{
			{
				Key:   string(sdk.AccountParameterDataRetentionTimeInDays),
				Value: "10",
			},
		},
		sdk.AccountParameterDataRetentionTimeInDays,
	)
	providerConfig := createProviderWithValuePropertyAndCustomDiff(t, schema.TypeInt, customDiff)

	t.Run("value set in the configuration and state", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.NumberIntVal(10),
		}), map[string]any{
			"value": "10",
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("value set only in the configuration", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.NumberIntVal(10),
		}), map[string]any{})
		assert.True(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("value set in the state and not equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"value": "20",
		})
		assert.Equal(t, "10", diff.Attributes["value"].New)
	})

	t.Run("value set in the state and equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"value": "10",
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})
}

func TestAccountObjectBoolValueComputedIf(t *testing.T) {
	customDiff := resources.AccountObjectBoolValueComputedIf(
		"value",
		[]*sdk.Parameter{
			{
				Key:   string(sdk.AccountParameterEnableConsoleOutput),
				Value: "true",
			},
		},
		sdk.AccountParameterEnableConsoleOutput,
	)
	providerConfig := createProviderWithValuePropertyAndCustomDiff(t, schema.TypeBool, customDiff)

	t.Run("value set in the configuration and state", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.BoolVal(true),
		}), map[string]any{
			"value": "true",
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("value set only in the configuration", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.BoolVal(true),
		}), map[string]any{})
		assert.True(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("value set in the state and not equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"value": "false",
		})
		assert.Equal(t, "true", diff.Attributes["value"].New)
	})

	t.Run("value set in the state and equals with parameter", func(t *testing.T) {
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.Type{}), map[string]any{
			"value": "true",
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})
}

func createProviderWithValuePropertyAndCustomDiff(t *testing.T, valueType schema.ValueType, customDiffFunc schema.CustomizeDiffFunc) *schema.Provider {
	t.Helper()
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test": {
				Schema: map[string]*schema.Schema{
					"value": {
						Type:     valueType,
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
