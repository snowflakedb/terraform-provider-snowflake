//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func Test_IgnoreAfterCreation(t *testing.T) {
	testSchema := map[string]*schema.Schema{
		"value": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	t.Run("without id", func(t *testing.T) {
		in := map[string]any{}
		d := schema.TestResourceDataRaw(t, testSchema, in)

		result := resources.IgnoreAfterCreation("", "", "", d)
		assert.False(t, result)
	})

	t.Run("with id", func(t *testing.T) {
		in := map[string]any{}
		d := schema.TestResourceDataRaw(t, testSchema, in)
		d.SetId("something")

		result := resources.IgnoreAfterCreation("", "", "", d)
		assert.True(t, result)
	})
}

func Test_NormalizeAndCompareIdentifiersSet(t *testing.T) {
	rawDataWithValues := func(values []any) *schema.ResourceData {
		return schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"value": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}, map[string]any{
			"value": values,
		})
	}
	emptyResourceData := rawDataWithValues([]any{})

	t.Run("validation: size key", func(t *testing.T) {
		assert.False(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.#", "1", "2", emptyResourceData))
	})

	t.Run("validation: case mismatch", func(t *testing.T) {
		resourceData := rawDataWithValues([]any{"SCHEMA.OBJECT.IDENTIFIER"})
		assert.False(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "schema.object.identifier", "", resourceData))
		// TODO(SNOW-1511594): Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.False(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", `"schema"."object"."identifier"`, resourceData))
	})

	t.Run("validation: case mismatch quoted identifier in the state", func(t *testing.T) {
		resourceData := rawDataWithValues([]any{`"SCHEMA"."OBJECT"."IDENTIFIER"`})
		assert.False(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "schema.object.identifier", "", resourceData))
		// TODO(SNOW-1511594): Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.False(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", `"schema"."object"."identifier"`, resourceData))
	})

	t.Run(`change suppressed from schema.object.identifier to "schema"."object"."identifier" with schema.object.identifier in state`, func(t *testing.T) {
		resourceData := rawDataWithValues([]any{"schema.object.identifier"})
		assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "schema.object.identifier", "", resourceData))
		// TODO(SNOW-1511594): Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", `"schema"."object"."identifier"`, resourceData))
	})

	t.Run(`change suppressed from schema.object.identifier to "schema"."object"."identifier" with "schema"."object"."identifier" in state`, func(t *testing.T) {
		resourceData := rawDataWithValues([]any{`"schema"."object"."identifier"`})
		assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "schema.object.identifier", "", resourceData))
		// TODO: Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", `"schema"."object"."identifier"`, resourceData))
	})

	t.Run(`change suppressed from "schema"."object"."identifier" to schema.object.identifier with schema.object.identifier in state`, func(t *testing.T) {
		resourceData := rawDataWithValues([]any{"schema.object.identifier"})
		assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", `"schema"."object"."identifier"`, "", resourceData))
		// TODO(SNOW-1511594): Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", "schema.object.identifier", resourceData))
	})

	t.Run(`change suppressed from "schema"."object"."identifier" to schema.object.identifier with "schema"."object"."identifier" in state`, func(t *testing.T) {
		resourceData := rawDataWithValues([]any{`"schema"."object"."identifier"`})
		assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", `"schema"."object"."identifier"`, "", resourceData))
		// TODO(SNOW-1511594): Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", "schema.object.identifier`, resourceData))
	})

	t.Run(`change suppressed from "SCHEMA"."OBJECT"."IDENTIFIER" to SCHEMA.OBJECT.IDENTIFIER with "SCHEMA"."OBJECT"."IDENTIFIER" in state`, func(t *testing.T) {
		resourceData := rawDataWithValues([]any{`"SCHEMA"."OBJECT"."IDENTIFIER"`})
		assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", `"SCHEMA"."OBJECT"."IDENTIFIER"`, "", resourceData))
		// TODO(SNOW-1511594): Cannot be tested with schema.TestResourceDataRaw because it doesn't populate raw state which is used in the cases below
		// assert.True(t, resources.NormalizeAndCompareIdentifiersInSet("value")("value.doesnt_matter", "", "SCHEMA.OBJECT.IDENTIFIER", resourceData))
	})
}

func Test_ignoreNewEmptyList(t *testing.T) {
	tests := []struct {
		name      string
		subfields []string
		key       string
		old       string
		new       string
		suppress  bool
	}{
		{
			name:     "suppress on zero count",
			key:      "a.#",
			old:      "5",
			new:      "0",
			suppress: true,
		},
		{
			name:      "suppress on ignored field",
			key:       "a.0.b",
			subfields: []string{"b"},
			suppress:  true,
		},
		{
			name:      "suppress on nested ignored field",
			key:       "a.0.b.c.d",
			subfields: []string{"b.c.d"},
			suppress:  true,
		},
		{
			name:     "do not suppress on non-zero count",
			key:      "a.#",
			new:      "5",
			suppress: false,
		},
		{
			name:      "do not suppress on non-ignored field",
			key:       "a.0.b",
			subfields: []string{"c"},
			suppress:  false,
		},
		{
			name:     "do not suppress on invalid key",
			key:      "a",
			suppress: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.suppress, resources.IgnoreNewEmptyListOrSubfields(tt.subfields...)(tt.key, tt.old, tt.new, nil))
		})
	}
}

func Test_IgnoreMatchingColumnNameAndMaskingPolicyUsingFirstElem(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"column": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"column_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"masking_policy": {
						Type: schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"using": {
									Type:     schema.TypeList,
									Optional: true,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	resourceData := func(using ...any) map[string]any {
		return map[string]any{
			"column": []any{
				map[string]any{
					"column_name": "foo",
					"masking_policy": []any{
						map[string]any{
							"using": using,
						},
					},
				},
			},
		}
	}
	tests := []struct {
		name         string
		key          string
		old          string
		new          string
		resourceData *schema.ResourceData
		wantSuppress bool
	}{
		{
			name:         "suppress when USING is not specified in the config, but is in the state - check count",
			key:          "column.0.masking_policy.0.using.#",
			old:          "1",
			new:          "0",
			resourceData: schema.TestResourceDataRaw(t, resourceSchema, resourceData("foo")),
			wantSuppress: true,
		},
		{
			name:         "suppress when USING is not specified in the config, but is in the state - check elem",
			key:          "column.0.masking_policy.0.using.0",
			old:          "foo",
			new:          "",
			resourceData: schema.TestResourceDataRaw(t, resourceSchema, resourceData("foo")),
			wantSuppress: true,
		},
		{
			name:         "do not suppress when there is column name mismatch",
			key:          "column.0.masking_policy.0.using.0",
			old:          "foo",
			new:          "bar",
			resourceData: schema.TestResourceDataRaw(t, resourceSchema, resourceData("foo")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantSuppress, resources.IgnoreMatchingColumnNameAndMaskingPolicyUsingFirstElem()(tt.key, tt.old, tt.new, tt.resourceData))
		})
	}
}
