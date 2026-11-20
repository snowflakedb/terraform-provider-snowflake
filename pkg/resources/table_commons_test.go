package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_handlePolicyReferences(t *testing.T) {
	tablePolicyKinds := []sdk.PolicyKind{
		sdk.PolicyKindRowAccessPolicy,
		sdk.PolicyKindAggregationPolicy,
		sdk.PolicyKindJoinPolicy,
	}
	rapAggKinds := []sdk.PolicyKind{
		sdk.PolicyKindRowAccessPolicy,
		sdk.PolicyKindAggregationPolicy,
	}

	t.Run("populates allowed policies and skips others", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, tableSchema, map[string]any{})
		err := handlePolicyReferences([]sdk.PolicyReference{
			rootLevelPolicyRef(sdk.PolicyKindRowAccessPolicy, "rap", "[id,name]"),
			rootLevelPolicyRef(sdk.PolicyKindAggregationPolicy, "ap", "[id]"),
			rootLevelPolicyRef(sdk.PolicyKindJoinPolicy, "jp", "[id]"),
			rootLevelPolicyRef(sdk.PolicyKindMaskingPolicy, "mp", "[id]"),
		}, d, tablePolicyKinds...)
		require.NoError(t, err)

		assertRootLevelPolicy(t, d, "row_access_policy", "on", "rap", "id", "name")
		assertRootLevelPolicy(t, d, "aggregation_policy", "entity_key", "ap", "id")
		assertRootLevelPolicy(t, d, "join_policy", "allowed_join_keys", "jp", "id")
	})

	t.Run("does not set disallowed kinds", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, tableSchema, map[string]any{})
		err := handlePolicyReferences([]sdk.PolicyReference{
			rootLevelPolicyRef(sdk.PolicyKindJoinPolicy, "jp", "[id]"),
			rootLevelPolicyRef(sdk.PolicyKindRowAccessPolicy, "rap", "[id]"),
		}, d, rapAggKinds...)
		require.NoError(t, err)

		assertRootLevelPolicy(t, d, "row_access_policy", "on", "rap", "id")
		require.Empty(t, d.Get("aggregation_policy"))
		require.Empty(t, d.Get("join_policy"))
	})

	t.Run("clears allowed kinds when there are no matching references", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, tableSchema, map[string]any{
			"join_policy": []any{map[string]any{
				"policy_name":       `"db"."sch"."jp"`,
				"allowed_join_keys": []any{"id"},
			}},
		})
		err := handlePolicyReferences(nil, d, tablePolicyKinds...)
		require.NoError(t, err)
		require.Empty(t, d.Get("row_access_policy"))
		require.Empty(t, d.Get("aggregation_policy"))
		require.Empty(t, d.Get("join_policy"))
	})

	t.Run("treats missing column names as empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, tableSchema, map[string]any{})
		err := handlePolicyReferences([]sdk.PolicyReference{{
			PolicyDb:     sdk.String("db"),
			PolicySchema: sdk.String("sch"),
			PolicyName:   "ap",
			PolicyKind:   sdk.PolicyKindAggregationPolicy,
		}}, d, sdk.PolicyKindAggregationPolicy)
		require.NoError(t, err)
		assertRootLevelPolicy(t, d, "aggregation_policy", "entity_key", "ap")
	})

	t.Run("rejects an unsupported allowed kind", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, tableSchema, map[string]any{})
		err := handlePolicyReferences(nil, d, sdk.PolicyKindMaskingPolicy)
		require.ErrorContains(t, err, "unsupported root-level policy kind")
	})
}

func rootLevelPolicyRef(kind sdk.PolicyKind, name, columns string) sdk.PolicyReference {
	return sdk.PolicyReference{
		PolicyDb:          sdk.String("db"),
		PolicySchema:      sdk.String("sch"),
		PolicyName:        name,
		PolicyKind:        kind,
		RefArgColumnNames: sdk.String(columns),
	}
}

func assertRootLevelPolicy(t *testing.T, d *schema.ResourceData, field, columnsKey, policyName string, columns ...string) {
	t.Helper()
	got := d.Get(field).([]any)
	require.Len(t, got, 1)
	block := got[0].(map[string]any)
	require.Equal(t, sdk.NewSchemaObjectIdentifier("db", "sch", policyName).FullyQualifiedName(), block["policy_name"])
	require.ElementsMatch(t, columns, block[columnsKey].(*schema.Set).List())
}
