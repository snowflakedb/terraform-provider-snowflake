package resources

import (
	"context"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// columnSchemaForCustomDiffTest mirrors the subset of the `column` schema that
// parseHybridColumns / forceNewIfColumnFieldChanged read. It deliberately omits
// schema-level defaults so that test cases can pin every field explicitly and
// the assertions don't depend on default-fill behavior.
func columnSchemaForCustomDiffTest() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"column": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name":     {Type: schema.TypeString, Required: true},
					"type":     {Type: schema.TypeString, Required: true},
					"nullable": {Type: schema.TypeBool, Optional: true},
					"collate":  {Type: schema.TypeString, Optional: true},
					"comment":  {Type: schema.TypeString, Optional: true},
					"default": {
						Type:     schema.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"constant":   {Type: schema.TypeString, Optional: true},
								"expression": {Type: schema.TypeString, Optional: true},
							},
						},
					},
				},
			},
		},
	}
}

// runColumnCustomDiff runs the given CustomizeDiffFunc against a synthetic
// resource whose schema mirrors the column block, and returns the resolved
// InstanceDiff so individual attributes can be inspected.
func runColumnCustomDiff(t *testing.T, customDiff schema.CustomizeDiffFunc, oldState map[string]string, newConfig map[string]any) *terraform.InstanceDiff {
	t.Helper()
	p := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"test": {
				Schema:        columnSchemaForCustomDiffTest(),
				CustomizeDiff: customdiff.All(customDiff),
			},
		},
	}
	diff, err := p.ResourcesMap["test"].Diff(
		context.Background(),
		&terraform.InstanceState{Attributes: oldState},
		&terraform.ResourceConfig{Config: newConfig},
		&provider.Context{Client: &sdk.Client{}},
	)
	require.NoError(t, err)
	return diff
}

// columnAttrs serializes a list of column maps into the dotted-key attribute
// shape the SDK expects for InstanceState.Attributes.
func columnAttrs(cols []map[string]string) map[string]string {
	attrs := map[string]string{
		"column.#": strconv.Itoa(len(cols)),
	}
	for i, c := range cols {
		attrs["column."+strconv.Itoa(i)+".name"] = c["name"]
		attrs["column."+strconv.Itoa(i)+".type"] = c["type"]
		attrs["column."+strconv.Itoa(i)+".nullable"] = c["nullable"]
		attrs["column."+strconv.Itoa(i)+".collate"] = c["collate"]
		attrs["column."+strconv.Itoa(i)+".comment"] = c["comment"]
		attrs["column."+strconv.Itoa(i)+".default.#"] = "0"
	}
	return attrs
}

func col(name, ty, nullable, collate, comment string) map[string]any {
	return map[string]any{
		"name":     name,
		"type":     ty,
		"nullable": nullable == "true",
		"collate":  collate,
		"comment":  comment,
	}
}

func Test_forceNewIfColumnNullableChanged(t *testing.T) {
	t.Run("nullable flipped on existing column forces new on that nullable path", func(t *testing.T) {
		old := columnAttrs([]map[string]string{
			{"name": "ID", "type": "NUMBER(38,0)", "nullable": "true", "collate": "", "comment": ""},
		})
		newCfg := map[string]any{
			"column": []any{col("ID", "NUMBER(38,0)", "false", "", "")},
		}
		diff := runColumnCustomDiff(t, forceNewIfColumnNullableChanged(), old, newCfg)
		assert.True(t, diff.RequiresNew(), "expected ForceNew when nullable flips")
		require.NotNil(t, diff.Attributes["column.0.nullable"])
		assert.True(t, diff.Attributes["column.0.nullable"].RequiresNew, "column.0.nullable should carry RequiresNew")
	})

	t.Run("nullable unchanged on existing column does not force new", func(t *testing.T) {
		old := columnAttrs([]map[string]string{
			{"name": "ID", "type": "NUMBER(38,0)", "nullable": "true", "collate": "", "comment": ""},
		})
		newCfg := map[string]any{
			"column": []any{col("ID", "NUMBER(38,0)", "true", "", "")},
		}
		diff := runColumnCustomDiff(t, forceNewIfColumnNullableChanged(), old, newCfg)
		assert.False(t, diff.RequiresNew())
	})

	t.Run("brand-new column added does not force new", func(t *testing.T) {
		old := columnAttrs([]map[string]string{
			{"name": "ID", "type": "NUMBER(38,0)", "nullable": "true", "collate": "", "comment": ""},
		})
		newCfg := map[string]any{
			"column": []any{
				col("ID", "NUMBER(38,0)", "true", "", ""),
				col("AGE", "NUMBER(38,0)", "false", "", ""),
			},
		}
		diff := runColumnCustomDiff(t, forceNewIfColumnNullableChanged(), old, newCfg)
		assert.False(t, diff.RequiresNew(), "adding a column with nullable=false should not trigger ForceNew on existing columns")
	})

	t.Run("column dropped does not force new", func(t *testing.T) {
		old := columnAttrs([]map[string]string{
			{"name": "ID", "type": "NUMBER(38,0)", "nullable": "true", "collate": "", "comment": ""},
			{"name": "AGE", "type": "NUMBER(38,0)", "nullable": "false", "collate": "", "comment": ""},
		})
		newCfg := map[string]any{
			"column": []any{col("ID", "NUMBER(38,0)", "true", "", "")},
		}
		diff := runColumnCustomDiff(t, forceNewIfColumnNullableChanged(), old, newCfg)
		assert.False(t, diff.RequiresNew())
	})

	t.Run("reorder + nullable change forces new on the NEW index", func(t *testing.T) {
		old := columnAttrs([]map[string]string{
			{"name": "ID", "type": "NUMBER(38,0)", "nullable": "true", "collate": "", "comment": ""},
			{"name": "AGE", "type": "NUMBER(38,0)", "nullable": "true", "collate": "", "comment": ""},
		})
		// AGE moves from index 1 to index 0 with nullable flipped to false.
		newCfg := map[string]any{
			"column": []any{
				col("AGE", "NUMBER(38,0)", "false", "", ""),
				col("ID", "NUMBER(38,0)", "true", "", ""),
			},
		}
		diff := runColumnCustomDiff(t, forceNewIfColumnNullableChanged(), old, newCfg)
		assert.True(t, diff.RequiresNew())
		require.NotNil(t, diff.Attributes["column.0.nullable"], "ForceNew must target the NEW index of the renamed-position column")
		assert.True(t, diff.Attributes["column.0.nullable"].RequiresNew)
	})

	t.Run("create scenario (no prior state) does not force new", func(t *testing.T) {
		newCfg := map[string]any{
			"column": []any{col("ID", "NUMBER(38,0)", "false", "", "")},
		}
		diff := runColumnCustomDiff(t, forceNewIfColumnNullableChanged(), map[string]string{}, newCfg)
		assert.False(t, diff.RequiresNew(), "create-time diffs (empty old state) must never trip ForceNew on the new column")
	})
}

func Test_forceNewIfColumnCollateChanged(t *testing.T) {
	t.Run("collate changed on existing column forces new on collate path", func(t *testing.T) {
		old := columnAttrs([]map[string]string{
			{"name": "NAME", "type": "VARCHAR(200)", "nullable": "true", "collate": "en-ci", "comment": ""},
		})
		newCfg := map[string]any{
			"column": []any{col("NAME", "VARCHAR(200)", "true", "en-cs", "")},
		}
		diff := runColumnCustomDiff(t, forceNewIfColumnCollateChanged(), old, newCfg)
		assert.True(t, diff.RequiresNew())
		require.NotNil(t, diff.Attributes["column.0.collate"])
		assert.True(t, diff.Attributes["column.0.collate"].RequiresNew)
	})

	t.Run("collate unchanged does not force new", func(t *testing.T) {
		old := columnAttrs([]map[string]string{
			{"name": "NAME", "type": "VARCHAR(200)", "nullable": "true", "collate": "en-ci", "comment": ""},
		})
		newCfg := map[string]any{
			"column": []any{col("NAME", "VARCHAR(200)", "true", "en-ci", "")},
		}
		diff := runColumnCustomDiff(t, forceNewIfColumnCollateChanged(), old, newCfg)
		assert.False(t, diff.RequiresNew())
	})
}

func Test_forceNewIfColumnFieldChanged(t *testing.T) {
	t.Run("no `column` change at all is a no-op", func(t *testing.T) {
		// Equal old and new state for `column`; the changed predicate must never fire.
		called := false
		predicate := func(_, _ column) bool {
			called = true
			return true
		}
		old := columnAttrs([]map[string]string{
			{"name": "ID", "type": "NUMBER(38,0)", "nullable": "true", "collate": "", "comment": ""},
		})
		newCfg := map[string]any{
			"column": []any{col("ID", "NUMBER(38,0)", "true", "", "")},
		}
		diff := runColumnCustomDiff(t, forceNewIfColumnFieldChanged("nullable", predicate), old, newCfg)
		assert.False(t, diff.RequiresNew())
		assert.False(t, called, "predicate must not be invoked when `column` has no change")
	})

	t.Run("predicate that always returns false never forces new", func(t *testing.T) {
		old := columnAttrs([]map[string]string{
			{"name": "ID", "type": "NUMBER(38,0)", "nullable": "true", "collate": "", "comment": ""},
		})
		newCfg := map[string]any{
			"column": []any{col("ID", "NUMBER(38,0)", "false", "", "")},
		}
		diff := runColumnCustomDiff(t, forceNewIfColumnFieldChanged("nullable", func(_, _ column) bool { return false }), old, newCfg)
		assert.False(t, diff.RequiresNew())
	})

	t.Run("custom field name targets the correct nested path", func(t *testing.T) {
		// Force ForceNew on `column.0.comment` using a synthetic predicate.
		old := columnAttrs([]map[string]string{
			{"name": "ID", "type": "NUMBER(38,0)", "nullable": "true", "collate": "", "comment": "old"},
		})
		newCfg := map[string]any{
			"column": []any{col("ID", "NUMBER(38,0)", "true", "", "new")},
		}
		diff := runColumnCustomDiff(t, forceNewIfColumnFieldChanged("comment", func(o, n column) bool {
			return o.comment != n.comment
		}), old, newCfg)
		assert.True(t, diff.RequiresNew())
		require.NotNil(t, diff.Attributes["column.0.comment"])
		assert.True(t, diff.Attributes["column.0.comment"].RequiresNew)
	})
}
