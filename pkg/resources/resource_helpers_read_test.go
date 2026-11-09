package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_setListKeepingNull(t *testing.T) {
	listSchema := map[string]*schema.Schema{
		"items": {Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
	}
	refresh := func(t *testing.T, stateAttrs map[string]string, value []string) map[string]string {
		t.Helper()
		d, err := schema.InternalMap(listSchema).Data(&terraform.InstanceState{ID: "id", Attributes: stateAttrs}, nil)
		require.NoError(t, err)
		require.NoError(t, setListKeepingNull(d, "items", value))
		return d.State().Attributes
	}

	t.Run("null state stays null when Snowflake returns no items", func(t *testing.T) {
		assert.NotContains(t, refresh(t, map[string]string{}, []string{}), "items.#")
	})

	t.Run("empty state stays empty when Snowflake returns no items", func(t *testing.T) {
		assert.Equal(t, "0", refresh(t, map[string]string{"items.#": "0"}, nil)["items.#"])
	})

	t.Run("items removed in Snowflake are cleared", func(t *testing.T) {
		attrs := refresh(t, map[string]string{"items.#": "1", "items.0": "CPU_X64_M"}, []string{})
		assert.Equal(t, "0", attrs["items.#"])
		assert.NotContains(t, attrs, "items.0")
	})

	t.Run("items from Snowflake are set", func(t *testing.T) {
		attrs := refresh(t, map[string]string{}, []string{"CPU_X64_M", "CPU_X64_L"})
		assert.Equal(t, "2", attrs["items.#"])
		assert.Equal(t, "CPU_X64_M", attrs["items.0"])
		assert.Equal(t, "CPU_X64_L", attrs["items.1"])
	})
}

func Test_derefOrZero(t *testing.T) {
	t.Run("nil string pointer returns empty string", func(t *testing.T) {
		var value *string
		assert.Equal(t, "", derefOrZero(value))
	})

	t.Run("string pointer returns the value", func(t *testing.T) {
		value := "WH"
		assert.Equal(t, "WH", derefOrZero(&value))
	})

	t.Run("nil int pointer returns zero", func(t *testing.T) {
		var value *int
		assert.Equal(t, 0, derefOrZero(value))
	})

	t.Run("int pointer returns the value", func(t *testing.T) {
		value := 42
		assert.Equal(t, 42, derefOrZero(&value))
	})
}
