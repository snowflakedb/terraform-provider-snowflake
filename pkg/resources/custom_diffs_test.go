//go:build !account_level_tests

package resources_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParameterValueComputedIf(t *testing.T) {
	createProviderConfig := func(parameterLevel sdk.ParameterType, parameterValue sdk.LogLevel) *schema.Provider {
		customDiff := resources.ParameterValueComputedIf(
			"value",
			[]*sdk.Parameter{
				{
					Key:   string(sdk.AccountParameterLogLevel),
					Level: parameterLevel,
					Value: string(parameterValue),
				},
			},
			sdk.ParameterTypeDatabase,
			sdk.AccountParameterLogLevel,
			func(v any) string { return v.(string) },
		)
		return createProviderWithValuePropertyAndCustomDiff(t, &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		}, customDiff)
	}

	t.Run("config: true - state: true - level: different - value: same", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelInfo)
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.True(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: true - state: true - level: different - value: different", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelDebug)
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: true - state: true - level: same - value: same", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeDatabase, sdk.LogLevelInfo)
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: true - state: true - level: same - value: different", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeDatabase, sdk.LogLevelDebug)
		diff := calculateDiff(t, providerConfig, cty.MapVal(map[string]cty.Value{
			"value": cty.StringVal(string(sdk.LogLevelInfo)),
		}), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: false - state: true - level: different - value: same", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelInfo)
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.String), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: false - state: true - level: different - value: different", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelDebug)
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.String), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.True(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: false - state: true - level: same - value: same", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelInfo)
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.String), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.False(t, diff.Attributes["value"].NewComputed)
	})

	t.Run("config: false - state: true - level: same - value: different", func(t *testing.T) {
		providerConfig := createProviderConfig(sdk.ParameterTypeAccount, sdk.LogLevelDebug)
		diff := calculateDiff(t, providerConfig, cty.MapValEmpty(cty.String), map[string]any{
			"value": string(sdk.LogLevelInfo),
		})
		assert.True(t, diff.Attributes["value"].NewComputed)
	})

	// Tests for filled config and empty state were not added as the only way
	// of getting into this situation would be in create operation for which custom diffs are skipped.
}

func createProviderWithValuePropertyAndCustomDiff(t *testing.T, valueSchema *schema.Schema, customDiffFunc schema.CustomizeDiffFunc) *schema.Provider {
	t.Helper()
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test": {
				Schema: map[string]*schema.Schema{
					"value": valueSchema,
				},
				CustomizeDiff: customDiffFunc,
			},
		},
	}
}

func createProviderWithNamedPropertyAndCustomDiff(t *testing.T, propertyName string, valueSchema *schema.Schema, customDiffFunc schema.CustomizeDiffFunc) *schema.Provider {
	t.Helper()
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test": {
				Schema: map[string]*schema.Schema{
					propertyName: valueSchema,
				},
				CustomizeDiff: customDiffFunc,
			},
		},
	}
}

func createProviderWithCustomSchemaAndCustomDiff(t *testing.T, customSchema map[string]*schema.Schema, customDiffFunc schema.CustomizeDiffFunc) *schema.Provider {
	t.Helper()
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test": {
				Schema:        customSchema,
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
		&provider.Context{Client: &sdk.Client{}},
	)
	require.NoError(t, err)
	return diff
}

func calculateDiffFromAttributes(t *testing.T, providerConfig *schema.Provider, oldValue map[string]string, newValue map[string]any) *terraform.InstanceDiff {
	t.Helper()
	diff, err := providerConfig.ResourcesMap["test"].Diff(
		context.Background(),
		&terraform.InstanceState{
			Attributes: oldValue,
		},
		&terraform.ResourceConfig{
			Config: newValue,
		},
		&provider.Context{Client: &sdk.Client{}},
	)
	require.NoError(t, err)
	return diff
}

func TestForceNewIfChangeToEmptyString(t *testing.T) {
	tests := []struct {
		name           string
		stateValue     map[string]string
		rawConfigValue map[string]any
		wantForceNew   bool
	}{
		{
			name:       "empty to non-empty",
			stateValue: map[string]string{},
			rawConfigValue: map[string]any{
				"value": "foo",
			},
			wantForceNew: false,
		}, {
			name:           "empty to empty",
			stateValue:     map[string]string{},
			rawConfigValue: map[string]any{},
			wantForceNew:   false,
		}, {
			name: "non-empty to empty",
			stateValue: map[string]string{
				"value": "foo",
			},
			rawConfigValue: map[string]any{},
			wantForceNew:   true,
		}, {
			name: "non-empty to non-empty",
			stateValue: map[string]string{
				"value": "bar",
			},
			rawConfigValue: map[string]any{
				"value": "foo",
			},
			wantForceNew: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customDiff := resources.ForceNewIfChangeToEmptyString(
				"value",
			)
			provider := createProviderWithValuePropertyAndCustomDiff(t, &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}, customDiff)
			diff := calculateDiffFromAttributes(
				t,
				provider,
				tt.stateValue,
				tt.rawConfigValue,
			)
			assert.Equal(t, tt.wantForceNew, diff.RequiresNew())
		})
	}
}

func TestForceNewIfChangeToEmptySlice(t *testing.T) {
	tests := []struct {
		name           string
		stateValue     map[string]string
		rawConfigValue map[string]any
		wantForceNew   bool
	}{
		{
			name:       "empty to non-empty",
			stateValue: map[string]string{},
			rawConfigValue: map[string]any{
				"value": []any{"foo"},
			},
			wantForceNew: false,
		}, {
			name:           "empty to empty",
			stateValue:     map[string]string{},
			rawConfigValue: map[string]any{},
			wantForceNew:   false,
		}, {
			name: "non-empty to empty",
			stateValue: map[string]string{
				"value.#": "1",
				"value.0": "foo",
			},
			rawConfigValue: map[string]any{},
			wantForceNew:   true,
		}, {
			name: "non-empty to non-empty",
			stateValue: map[string]string{
				"value.#": "2",
				"value.0": "foo",
				"value.1": "bar",
			},
			rawConfigValue: map[string]any{
				"value": []any{"foo"},
			},
			wantForceNew: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customDiff := resources.ForceNewIfChangeToEmptySlice[any](
				"value",
			)
			provider := createProviderWithValuePropertyAndCustomDiff(t, &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			}, customDiff)
			diff := calculateDiffFromAttributes(
				t,
				provider,
				tt.stateValue,
				tt.rawConfigValue,
			)
			assert.Equal(t, tt.wantForceNew, diff.RequiresNew())
		})
	}
}

func TestForceNewIfChangeToEmptySet(t *testing.T) {
	tests := []struct {
		name           string
		stateValue     map[string]string
		rawConfigValue map[string]any
		wantForceNew   bool
	}{
		{
			name:       "empty to non-empty",
			stateValue: map[string]string{},
			rawConfigValue: map[string]any{
				"value": []any{"foo"},
			},
			wantForceNew: false,
		}, {
			name:           "empty to empty",
			stateValue:     map[string]string{},
			rawConfigValue: map[string]any{},
			wantForceNew:   false,
		}, {
			name: "non-empty to empty",
			stateValue: map[string]string{
				"value.#": "1",
				// The Sets are using hashes to generate an index for a given value.
				// In this case: 2577344683 == hash("CREATE DATABASE").
				"value.2577344683": "CREATE DATABASE",
			},
			rawConfigValue: map[string]any{},
			wantForceNew:   true,
		}, {
			name: "non-empty to non-empty",
			stateValue: map[string]string{
				"value.#": "2",
				"value.0": "foo",
				"value.1": "bar",
			},
			rawConfigValue: map[string]any{
				"value": []any{"foo"},
			},
			wantForceNew: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff := calculateDiffFromAttributes(t,
				createProviderWithValuePropertyAndCustomDiff(t,
					&schema.Schema{
						Type: schema.TypeSet,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Optional: true,
					},
					resources.ForceNewIfChangeToEmptySet(
						"value",
					),
				),
				tt.stateValue,
				tt.rawConfigValue,
			)
			assert.Equal(t, tt.wantForceNew, diff.RequiresNew())
		})
	}
}

func Test_ComputedIfAnyAttributeChanged(t *testing.T) {
	testSuppressFunc := schema.SchemaDiffSuppressFunc(func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		return strings.Trim(oldValue, `"`) == strings.Trim(newValue, `"`)
	})
	testSchema := map[string]*schema.Schema{
		"value_with_diff_suppress": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: testSuppressFunc,
		},
		"value_without_diff_suppress": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"computed_value": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
	testCustomDiff := resources.ComputedIfAnyAttributeChanged(
		testSchema,
		"computed_value",
		"value_with_diff_suppress",
		"value_without_diff_suppress",
	)
	testProvider := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test": {
				Schema:        testSchema,
				CustomizeDiff: testCustomDiff,
			},
		},
	}

	tests := []struct {
		name           string
		stateValue     map[string]string
		rawConfigValue map[string]any
		expectDiff     bool
	}{
		{
			name: "no change on both fields",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
			},
			expectDiff: false,
		},
		{
			name: "change on field with diff suppression - suppressed (quotes in config added)",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "\"foo\"",
				"value_without_diff_suppress": "foo",
			},
			expectDiff: false,
		},
		{
			name: "change on field with diff suppression - suppressed (quotes in config removed)",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "\"foo\"",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
			},
			expectDiff: false,
		},
		{
			name: "change on field with diff suppression - not suppressed (value change)",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "bar",
				"value_without_diff_suppress": "foo",
			},
			expectDiff: true,
		},
		{
			name: "change on field with diff suppression - not suppressed (value and quotes changed)",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "\"foo\"",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "bar",
				"value_without_diff_suppress": "foo",
			},
			expectDiff: true,
		},
		{
			name: "change on field without diff suppression",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "bar",
			},
			expectDiff: true,
		},
		{
			name: "change on field without diff suppression, suppressed change on field with diff suppression",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "\"foo\"",
				"value_without_diff_suppress": "bar",
			},
			expectDiff: true,
		},
		{
			name: "change on field without diff suppression, not suppressed change on field with diff suppression",
			stateValue: map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			rawConfigValue: map[string]any{
				"value_with_diff_suppress":    "\"bar\"",
				"value_without_diff_suppress": "bar",
			},
			expectDiff: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			diff := calculateDiffFromAttributes(
				t,
				testProvider,
				tt.stateValue,
				tt.rawConfigValue,
			)
			if tt.expectDiff {
				require.NotNil(t, diff)
				require.NotNil(t, diff.Attributes["computed_value"])
				assert.True(t, diff.Attributes["computed_value"].NewComputed)
			} else {
				require.Nil(t, diff)
			}
		})
	}

	t.Run("attributes not found in schema, both fields changed", func(t *testing.T) {
		otherTestSchema := map[string]*schema.Schema{
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: testSuppressFunc,
			},
			"computed_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}
		otherTestCustomDiff := resources.ComputedIfAnyAttributeChanged(
			otherTestSchema,
			"computed_value",
			"value_with_diff_suppress",
			"value_without_diff_suppress",
		)
		otherTestProvider := &schema.Provider{
			ResourcesMap: map[string]*schema.Resource{
				"test": {
					Schema:        testSchema,
					CustomizeDiff: otherTestCustomDiff,
				},
			},
		}

		diff := calculateDiffFromAttributes(
			t,
			otherTestProvider,
			map[string]string{
				"value_with_diff_suppress":    "foo",
				"value_without_diff_suppress": "foo",
				"computed_value":              "foo",
			},
			map[string]any{
				"value_with_diff_suppress":    "\"bar\"",
				"value_without_diff_suppress": "bar",
			},
		)

		require.NotNil(t, diff)
		assert.Nil(t, diff.Attributes["computed_value"])
	})
}

func TestForceNewIfAllKeysAreNotSet(t *testing.T) {
	tests := []struct {
		name           string
		stateValue     map[string]string
		rawConfigValue map[string]any
		wantForceNew   bool
	}{
		{
			name: "all values set to unset",
			stateValue: map[string]string{
				"value":  "123",
				"value2": "string value",
				"value3": "[one two]",
			},
			rawConfigValue: map[string]any{},
			wantForceNew:   true,
		},
		{
			name: "only value set to unset",
			stateValue: map[string]string{
				"value": "123",
			},
			rawConfigValue: map[string]any{},
			wantForceNew:   true,
		},
		{
			name: "only value2 set to unset",
			stateValue: map[string]string{
				"value2": "string value",
			},
			rawConfigValue: map[string]any{},
			wantForceNew:   true,
		},
		{
			name: "only value3 set to unset",
			stateValue: map[string]string{
				"value3": "[one two]",
			},
			rawConfigValue: map[string]any{},
			// We expect here to not re-create because value3 doesn't have a custom diff on it
			// and the rest custom diffs don't work when the values are not set.
			wantForceNew: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &schema.Provider{
				ResourcesMap: map[string]*schema.Resource{
					"test": {
						Schema: map[string]*schema.Schema{
							"value": {
								Type: schema.TypeInt,
							},
							"value2": {
								Type: schema.TypeString,
							},
							"value3": {
								Type: schema.TypeList,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
						CustomizeDiff: customdiff.All(
							resources.ForceNewIfAllKeysAreNotSet("value", "value", "value2", "value3"),
							resources.ForceNewIfAllKeysAreNotSet("value2", "value", "value2", "value3"),
						),
					},
				},
			}
			diff := calculateDiffFromAttributes(
				t,
				p,
				tt.stateValue,
				tt.rawConfigValue,
			)
			assert.Equal(t, tt.wantForceNew, diff.RequiresNew())
		})
	}
}

func Test_RecreateWhenUserTypeChangedExternally(t *testing.T) {
	tests := []struct {
		name         string
		userType     sdk.UserType
		stateValue   map[string]string
		wantForceNew bool
	}{
		{
			name:         "person - nothing in state",
			userType:     sdk.UserTypePerson,
			stateValue:   map[string]string{},
			wantForceNew: false,
		},
		{
			name:     "person - person in state",
			userType: sdk.UserTypePerson,
			stateValue: map[string]string{
				"user_type": "PERSON",
			},
			wantForceNew: false,
		},
		{
			name:     "person - person in state lowercased",
			userType: sdk.UserTypePerson,
			stateValue: map[string]string{
				"user_type": "person",
			},
			wantForceNew: false,
		},
		{
			name:     "person - service in state",
			userType: sdk.UserTypePerson,
			stateValue: map[string]string{
				"user_type": "SERVICE",
			},
			wantForceNew: true,
		},
		{
			name:     "person - service in state lowercased",
			userType: sdk.UserTypePerson,
			stateValue: map[string]string{
				"user_type": "service",
			},
			wantForceNew: true,
		},
		{
			name:     "person - empty value in state",
			userType: sdk.UserTypePerson,
			stateValue: map[string]string{
				"user_type": "",
			},
			wantForceNew: false,
		},
		{
			name:     "person - garbage in state",
			userType: sdk.UserTypePerson,
			stateValue: map[string]string{
				"user_type": "garbage",
			},
			wantForceNew: true,
		},
		{
			name:         "service - nothing in state",
			userType:     sdk.UserTypeService,
			stateValue:   map[string]string{},
			wantForceNew: true,
		},
		{
			name:     "service - service in state",
			userType: sdk.UserTypeService,
			stateValue: map[string]string{
				"user_type": "SERVICE",
			},
			wantForceNew: false,
		},
		{
			name:     "service - service in state lowercased",
			userType: sdk.UserTypeService,
			stateValue: map[string]string{
				"user_type": "service",
			},
			wantForceNew: false,
		},
		{
			name:     "service - person in state",
			userType: sdk.UserTypeService,
			stateValue: map[string]string{
				"user_type": "PERSON",
			},
			wantForceNew: true,
		},
		{
			name:     "service - person in state lowercased",
			userType: sdk.UserTypeService,
			stateValue: map[string]string{
				"user_type": "person",
			},
			wantForceNew: true,
		},
		{
			name:     "service - empty value in state",
			userType: sdk.UserTypeService,
			stateValue: map[string]string{
				"user_type": "",
			},
			wantForceNew: true,
		},
		{
			name:     "service - garbage in state",
			userType: sdk.UserTypeService,
			stateValue: map[string]string{
				"user_type": "garbage",
			},
			wantForceNew: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			customDiff := resources.RecreateWhenUserTypeChangedExternally(tt.userType)
			testProvider := createProviderWithNamedPropertyAndCustomDiff(t, "user_type", &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			}, customDiff)
			diff := calculateDiffFromAttributes(
				t,
				testProvider,
				tt.stateValue,
				map[string]any{},
			)
			assert.Equal(t, tt.wantForceNew, diff.RequiresNew())
		})
	}
}

func Test_RecreateWhenSecretTypeChangedExternally(t *testing.T) {
	tests := []struct {
		name         string
		secretType   sdk.SecretType
		stateValue   map[string]string
		wantForceNew bool
	}{
		// password type
		{
			name:         "password - nothing in state",
			secretType:   sdk.SecretTypePassword,
			stateValue:   map[string]string{},
			wantForceNew: true,
		},
		{
			name:       "password - empty value in state",
			secretType: sdk.SecretTypePassword,
			stateValue: map[string]string{
				"secret_type": "",
			},
			wantForceNew: true,
		},
		{
			name:       "password - password in state",
			secretType: sdk.SecretTypePassword,
			stateValue: map[string]string{
				"secret_type": "PASSWORD",
			},
			wantForceNew: false,
		},
		{
			name:       "password - password in state lowercased",
			secretType: sdk.SecretTypePassword,
			stateValue: map[string]string{
				"secret_type": "password",
			},
			wantForceNew: false,
		},
		{
			name:       "password - oauth2 in state",
			secretType: sdk.SecretTypePassword,
			stateValue: map[string]string{
				"secret_type": "OAUTH2",
			},
			wantForceNew: true,
		},
		{
			name:       "password - generic_string in state",
			secretType: sdk.SecretTypePassword,
			stateValue: map[string]string{
				"secret_type": "GENERIC_STRING",
			},
			wantForceNew: true,
		},
		{
			name:       "password - oauth2 in state lowercased",
			secretType: sdk.SecretTypePassword,
			stateValue: map[string]string{
				"secret_type": "oauth2",
			},
			wantForceNew: true,
		},
		// generic string type
		{
			name:         "generic_string - nothing in state",
			secretType:   sdk.SecretTypeGenericString,
			stateValue:   map[string]string{},
			wantForceNew: true,
		},
		{
			name:       "generic_string - empty value in state",
			secretType: sdk.SecretTypeGenericString,
			stateValue: map[string]string{
				"secret_type": "",
			},
			wantForceNew: true,
		},
		{
			name:       "generic_string - generic_string in state",
			secretType: sdk.SecretTypeGenericString,
			stateValue: map[string]string{
				"secret_type": "generic_string",
			},
			wantForceNew: false,
		},
		{
			name:       "generic_string - generic_string in state lowercased",
			secretType: sdk.SecretTypeGenericString,
			stateValue: map[string]string{
				"secret_type": "generic_string",
			},
			wantForceNew: false,
		},
		{
			name:       "generic_string - oauth2 in state",
			secretType: sdk.SecretTypeGenericString,
			stateValue: map[string]string{
				"secret_type": "OAUTH2",
			},
			wantForceNew: true,
		},
		{
			name:       "generic_string - password in state",
			secretType: sdk.SecretTypeGenericString,
			stateValue: map[string]string{
				"secret_type": "PASSWORD",
			},
			wantForceNew: true,
		},
		{
			name:       "generic_string - oauth2 in state lowercased",
			secretType: sdk.SecretTypeGenericString,
			stateValue: map[string]string{
				"secret_type": "oauth2",
			},
			wantForceNew: true,
		},
		// oauth2 authorization code grant type
		{
			name:         "oauth2 authorization code grant - nothing in state",
			secretType:   sdk.SecretTypeOAuth2AuthorizationCodeGrant,
			stateValue:   map[string]string{},
			wantForceNew: true,
		},
		{
			name:       "oauth2 authorization code grant - empty value in state",
			secretType: sdk.SecretTypeOAuth2AuthorizationCodeGrant,
			stateValue: map[string]string{
				"secret_type": "",
			},
			wantForceNew: true,
		},
		{
			name:       "oauth2 authorization code grant - password in state",
			secretType: sdk.SecretTypeOAuth2AuthorizationCodeGrant,
			stateValue: map[string]string{
				"secret_type": "PASSWORD",
			},
			wantForceNew: true,
		},
		{
			name:       "oauth2 authorization code grant - generic_string in state",
			secretType: sdk.SecretTypeOAuth2AuthorizationCodeGrant,
			stateValue: map[string]string{
				"secret_type": "GENERIC_STRING",
			},
			wantForceNew: true,
		},
		// oauth2 client credentials type
		{
			name:         "oauth2 client credentials - nothing in state",
			secretType:   sdk.SecretTypeOAuth2ClientCredentials,
			stateValue:   map[string]string{},
			wantForceNew: true,
		},
		{
			name:       "oauth2 client credentials - empty value in state",
			secretType: sdk.SecretTypeOAuth2ClientCredentials,
			stateValue: map[string]string{
				"secret_type": "",
			},
			wantForceNew: true,
		},
		{
			name:       "oauth2 client credentials - password in state",
			secretType: sdk.SecretTypeOAuth2ClientCredentials,
			stateValue: map[string]string{
				"secret_type": "PASSWORD",
			},
			wantForceNew: true,
		},
		{
			name:       "oauth2 client credentials - generic_string in state",
			secretType: sdk.SecretTypeOAuth2ClientCredentials,
			stateValue: map[string]string{
				"secret_type": "GENERIC_STRING",
			},
			wantForceNew: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			customDiff := resources.RecreateWhenSecretTypeChangedExternally(tt.secretType)
			testProvider := createProviderWithNamedPropertyAndCustomDiff(t, "secret_type", &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			}, customDiff)
			diff := calculateDiffFromAttributes(
				t,
				testProvider,
				tt.stateValue,
				map[string]any{},
			)
			assert.Equal(t, tt.wantForceNew, diff.RequiresNew())
		})
	}
}

func Test_RecreateWhenSecretTypeChangedExternallyForOAuth2(t *testing.T) {
	tests := []struct {
		name         string
		secretType   sdk.SecretType
		stateValue   map[string]string
		wantForceNew bool
	}{
		// config          - authorization code
		// external change - drop and recreate with the same id but as oauth2 with client credentials
		{
			name:       "oauth2 authorization code - oauth2 client credentials in state",
			secretType: sdk.SecretTypeOAuth2AuthorizationCodeGrant,
			stateValue: map[string]string{
				"secret_type": "OAUTH2",
				"describe_output.0.oauth_refresh_token_expiry_time": "",
			},
			wantForceNew: true,
		},
		// config          - client credentials
		// external change - drop and recreate with the same id but as oauth2 with authorization code grant
		{
			name:       "oauth2 client credentials - oauth2 authorization code in state",
			secretType: sdk.SecretTypeOAuth2ClientCredentials,
			stateValue: map[string]string{
				"secret_type": "OAUTH2",
				"describe_output.0.oauth_refresh_token_expiry_time": "some test date here",
			},
			wantForceNew: true,
		},
		// no external change
		{
			name:       "oauth2 authorization code - oauth2 authorization code in state",
			secretType: sdk.SecretTypeOAuth2AuthorizationCodeGrant,
			stateValue: map[string]string{
				"secret_type": "OAUTH2",
				"describe_output.0.oauth_refresh_token_expiry_time": "some test date here",
			},
			wantForceNew: false,
		},
		// no external change
		{
			name:       "oauth2 client credentials - oauth2 client credentials code in state",
			secretType: sdk.SecretTypeOAuth2ClientCredentials,
			stateValue: map[string]string{
				"secret_type": "OAUTH2",
				"describe_output.0.oauth_refresh_token_expiry_time": "",
			},
			wantForceNew: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			customDiff := resources.RecreateWhenSecretTypeChangedExternally(tt.secretType)
			testProvider := createProviderWithCustomSchemaAndCustomDiff(t,
				map[string]*schema.Schema{
					"secret_type": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"describe_output": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"oauth_refresh_token_expiry_time": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
				customDiff)
			diff := calculateDiffFromAttributes(
				t,
				testProvider,
				tt.stateValue,
				map[string]any{},
			)
			assert.Equal(t, tt.wantForceNew, diff.RequiresNew())
		})
	}
}

func Test_RecreateWhenResourceBoolFieldChangedExternally(t *testing.T) {
	tests := []struct {
		name         string
		isPrimary    bool
		stateValue   map[string]string
		wantForceNew bool
	}{
		{
			name:      "changed is_primary from false to true",
			isPrimary: false,
			stateValue: map[string]string{
				"is_primary": "true",
			},
			wantForceNew: true,
		},
		{
			name:      "changed is_primary from true to false",
			isPrimary: true,
			stateValue: map[string]string{
				"is_primary": "false",
			},
			wantForceNew: true,
		},
		{
			name:      "no change in is_primary - true to true",
			isPrimary: true,
			stateValue: map[string]string{
				"is_primary": "true",
			},
			wantForceNew: false,
		},
		{
			name:      "no change in is_primary - false to false",
			isPrimary: false,
			stateValue: map[string]string{
				"is_primary": "false",
			},
			wantForceNew: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			customDiff := resources.RecreateWhenResourceBoolFieldChangedExternally("is_primary", tt.isPrimary)
			testProvider := createProviderWithCustomSchemaAndCustomDiff(t,
				map[string]*schema.Schema{
					"is_primary": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
				customDiff)
			diff := calculateDiffFromAttributes(
				t,
				testProvider,
				tt.stateValue,
				map[string]any{},
			)
			assert.Equal(t, tt.wantForceNew, diff.RequiresNew())
		})
	}
}
