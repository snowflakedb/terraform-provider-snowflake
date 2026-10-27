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
