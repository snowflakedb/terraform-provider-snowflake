package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func importTestSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"boolean_string": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  BooleanDefault,
		},
		"int_with_int_default": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  IntDefault,
		},
		"int": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"string": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"string_set": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
		},
		"block": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"string": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func Test_importBooleanStringAttribute(t *testing.T) {
	tests := []struct {
		name                 string
		importPermissiveMode bool
		snowflakeValue       bool
		expected             string
	}{
		{name: "enforcing, true", snowflakeValue: true, expected: BooleanTrue},
		{name: "enforcing, false", snowflakeValue: false, expected: BooleanFalse},
		{name: "permissive", importPermissiveMode: true, snowflakeValue: true, expected: BooleanDefault},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, importTestSchema(), map[string]any{})
			require.NoError(t, importBooleanStringAttribute(d, "boolean_string", tt.importPermissiveMode, tt.snowflakeValue))
			assert.Equal(t, tt.expected, d.Get("boolean_string"))
		})
	}
}

func Test_importIntWithIntDefaultAttribute(t *testing.T) {
	tests := []struct {
		name                 string
		importPermissiveMode bool
		snowflakeValue       int
		expected             int
	}{
		{name: "enforcing", snowflakeValue: 600, expected: 600},
		{name: "enforcing, zero", snowflakeValue: 0, expected: 0},
		{name: "permissive", importPermissiveMode: true, snowflakeValue: 600, expected: IntDefault},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, importTestSchema(), map[string]any{})
			require.NoError(t, importIntWithIntDefaultAttribute(d, "int_with_int_default", tt.importPermissiveMode, tt.snowflakeValue))
			assert.Equal(t, tt.expected, d.Get("int_with_int_default"))
		})
	}
}

func Test_importIntAttribute(t *testing.T) {
	tests := []struct {
		name                 string
		importPermissiveMode bool
		snowflakeValue       int
		expected             int
	}{
		{name: "enforcing", snowflakeValue: 600, expected: 600},
		{name: "enforcing, zero", snowflakeValue: 0, expected: 0},
		{name: "permissive", importPermissiveMode: true, snowflakeValue: 600, expected: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, importTestSchema(), map[string]any{})
			require.NoError(t, importIntAttribute(d, "int", tt.importPermissiveMode, tt.snowflakeValue))
			assert.Equal(t, tt.expected, d.Get("int"))
		})
	}
}

func Test_importStringAttribute(t *testing.T) {
	tests := []struct {
		name                 string
		importPermissiveMode bool
		snowflakeValue       string
		expected             string
	}{
		{name: "enforcing", snowflakeValue: "value", expected: "value"},
		{name: "enforcing, empty", snowflakeValue: "", expected: ""},
		{name: "permissive", importPermissiveMode: true, snowflakeValue: "value", expected: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, importTestSchema(), map[string]any{})
			require.NoError(t, importStringAttribute(d, "string", tt.importPermissiveMode, tt.snowflakeValue))
			assert.Equal(t, tt.expected, d.Get("string"))
		})
	}
}

// importedStateValue returns nil when set doesn't write key to the state.
func importedStateValue(t *testing.T, key string, set func(d *schema.ResourceData) error) *string {
	t.Helper()
	d := (&schema.Resource{Schema: importTestSchema()}).Data(nil)
	d.SetId("id")
	require.NoError(t, set(d))
	if v, ok := d.State().Attributes[key]; ok {
		return &v
	}
	return nil
}

func Test_importOptionalBooleanStringAttribute(t *testing.T) {
	tests := []struct {
		name                 string
		importPermissiveMode bool
		snowflakeValue       *bool
		expected             *string
	}{
		{name: "enforcing, true", snowflakeValue: new(true), expected: new(BooleanTrue)},
		{name: "enforcing, false", snowflakeValue: new(false), expected: new(BooleanFalse)},
		{name: "enforcing, nil", snowflakeValue: nil, expected: new(BooleanDefault)},
		{name: "permissive", importPermissiveMode: true, snowflakeValue: new(true), expected: new(BooleanDefault)},
		{name: "permissive, nil", importPermissiveMode: true, snowflakeValue: nil, expected: new(BooleanDefault)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, importedStateValue(t, "boolean_string", func(d *schema.ResourceData) error {
				return importOptionalBooleanStringAttribute(d, "boolean_string", tt.importPermissiveMode, tt.snowflakeValue)
			}))
		})
	}
}

func Test_importOptionalIntWithIntDefaultAttribute(t *testing.T) {
	tests := []struct {
		name                 string
		importPermissiveMode bool
		snowflakeValue       *int
		expected             *string
	}{
		{name: "enforcing", snowflakeValue: new(600), expected: new("600")},
		{name: "enforcing, zero", snowflakeValue: new(0), expected: new("0")},
		{name: "enforcing, nil", snowflakeValue: nil, expected: new(IntDefaultString)},
		{name: "permissive", importPermissiveMode: true, snowflakeValue: new(600), expected: new(IntDefaultString)},
		{name: "permissive, nil", importPermissiveMode: true, snowflakeValue: nil, expected: new(IntDefaultString)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, importedStateValue(t, "int_with_int_default", func(d *schema.ResourceData) error {
				return importOptionalIntWithIntDefaultAttribute(d, "int_with_int_default", tt.importPermissiveMode, tt.snowflakeValue)
			}))
		})
	}
}

func Test_importOptionalIntAttribute(t *testing.T) {
	tests := []struct {
		name                 string
		importPermissiveMode bool
		snowflakeValue       *int
		expected             *string
	}{
		{name: "enforcing", snowflakeValue: new(600), expected: new("600")},
		{name: "enforcing, zero", snowflakeValue: new(0), expected: new("0")},
		{name: "enforcing, nil", snowflakeValue: nil, expected: nil},
		{name: "permissive", importPermissiveMode: true, snowflakeValue: new(600), expected: new("0")},
		{name: "permissive, nil", importPermissiveMode: true, snowflakeValue: nil, expected: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, importedStateValue(t, "int", func(d *schema.ResourceData) error {
				return importOptionalIntAttribute(d, "int", tt.importPermissiveMode, tt.snowflakeValue)
			}))
		})
	}
}

func Test_importOptionalStringAttribute(t *testing.T) {
	type stringEnum string

	tests := []struct {
		name                 string
		importPermissiveMode bool
		snowflakeValue       *stringEnum
		expected             *string
	}{
		{name: "enforcing", snowflakeValue: new(stringEnum("value")), expected: new("value")},
		{name: "enforcing, nil", snowflakeValue: nil, expected: nil},
		{name: "permissive", importPermissiveMode: true, snowflakeValue: new(stringEnum("value")), expected: new("")},
		{name: "permissive, nil", importPermissiveMode: true, snowflakeValue: nil, expected: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, importedStateValue(t, "string", func(d *schema.ResourceData) error {
				return importOptionalStringAttribute(d, "string", tt.importPermissiveMode, tt.snowflakeValue)
			}))
		})
	}
}

func Test_importSetAttribute(t *testing.T) {
	type stringEnum string

	tests := []struct {
		name                 string
		importPermissiveMode bool
		snowflakeValue       []stringEnum
		expected             []string
	}{
		{name: "enforcing", snowflakeValue: []stringEnum{"A", "B"}, expected: []string{"A", "B"}},
		{name: "enforcing, empty", snowflakeValue: []stringEnum{}, expected: []string{}},
		{name: "permissive", importPermissiveMode: true, snowflakeValue: []stringEnum{"A", "B"}, expected: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, importTestSchema(), map[string]any{})
			require.NoError(t, importSetAttribute(d, "string_set", tt.importPermissiveMode, tt.snowflakeValue))
			assert.ElementsMatch(t, tt.expected, expandStringList(d.Get("string_set").(*schema.Set).List()))
		})
	}
}

func Test_importBlockAttribute(t *testing.T) {
	tests := []struct {
		name                 string
		importPermissiveMode bool
		snowflakeValue       []map[string]any
		expected             []any
	}{
		{name: "enforcing", snowflakeValue: []map[string]any{{"string": "value"}}, expected: []any{map[string]any{"string": "value"}}},
		{name: "enforcing, empty", snowflakeValue: []map[string]any{}, expected: []any{}},
		{name: "permissive", importPermissiveMode: true, snowflakeValue: []map[string]any{{"string": "value"}}, expected: []any{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, importTestSchema(), map[string]any{})
			require.NoError(t, importBlockAttribute(d, "block", tt.importPermissiveMode, tt.snowflakeValue))
			assert.Equal(t, tt.expected, d.Get("block"))
		})
	}
}
