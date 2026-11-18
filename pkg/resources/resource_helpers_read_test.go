package resources

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

func Test_showByIdErrorDiagnostics(t *testing.T) {
	id := sdk.NewSchemaObjectIdentifier("db", "sch", "policy")
	newData := func(t *testing.T) *schema.ResourceData {
		t.Helper()
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, nil)
		d.SetId("db.sch.policy")
		return d
	}

	t.Run("an object that does not exist is removed from the state with a warning", func(t *testing.T) {
		d := newData(t)
		diags := showByIdErrorDiagnostics(d, sdk.ObjectTypeBackupPolicy, id, fmt.Errorf("wrapped: %w", sdk.ErrObjectNotFound))

		require.Len(t, diags, 1)
		assert.Equal(t, diag.Warning, diags[0].Severity)
		assert.Equal(t, "Failed to query backup policy. Marking the resource as removed.", diags[0].Summary)
		assert.Equal(t, fmt.Sprintf("Resource backup policy id: %s, Err: wrapped: %s", id.FullyQualifiedName(), sdk.ErrObjectNotFound), diags[0].Detail)
		assert.Empty(t, d.Id())
	})

	t.Run("another error is returned and keeps the state", func(t *testing.T) {
		d := newData(t)
		diags := showByIdErrorDiagnostics(d, sdk.ObjectTypeBackupPolicy, id, errors.New("boom"))

		require.Len(t, diags, 1)
		assert.Equal(t, diag.Error, diags[0].Severity)
		assert.Equal(t, "boom", diags[0].Summary)
		assert.Equal(t, "db.sch.policy", d.Id())
	})
}

func Test_omitCurrentAccount(t *testing.T) {
	t.Parallel()

	current := sdk.NewAccountIdentifier("ORG", "SOURCE")
	target := sdk.NewAccountIdentifier("ORG", "TARGET")

	require.Equal(t, []sdk.AccountIdentifier{target}, omitCurrentAccount([]sdk.AccountIdentifier{current, target}, current))
	require.Empty(t, omitCurrentAccount([]sdk.AccountIdentifier{current}, current))
}
