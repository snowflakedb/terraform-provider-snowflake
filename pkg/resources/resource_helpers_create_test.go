package resources

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stringSetSchema is a TypeSet of strings named "items" for the set create and update tests.
func stringSetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"items": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func Test_booleanStringAttributeCreateFallbackBuilder(t *testing.T) {
	flagSchema := map[string]*schema.Schema{
		"flag": {Type: schema.TypeString, Optional: true, Default: BooleanDefault},
	}

	d := schema.TestResourceDataRaw(t, flagSchema, map[string]any{"flag": BooleanDefault})
	var got *bool
	require.NoError(t, booleanStringAttributeCreateFallbackBuilder(d, "flag", func(v bool) any {
		got = &v
		return nil
	}, true))
	require.NotNil(t, got)
	assert.True(t, *got)

	d = schema.TestResourceDataRaw(t, flagSchema, map[string]any{"flag": BooleanFalse})
	got = nil
	require.NoError(t, booleanStringAttributeCreateFallbackBuilder(d, "flag", func(v bool) any {
		got = &v
		return nil
	}, true))
	require.NotNil(t, got)
	assert.False(t, *got)
}

func Test_setValueCreate(t *testing.T) {
	d := schema.TestResourceDataRaw(t, stringSetSchema(), map[string]any{
		"items": []any{"b", "a"},
	})

	var got []string
	require.NoError(t, setValueCreate(d, "items", &got, func(v any) (string, error) {
		return v.(string), nil
	}))
	assert.ElementsMatch(t, []string{"a", "b"}, got)

	absent := schema.TestResourceDataRaw(t, stringSetSchema(), map[string]any{})
	var untouched []string
	require.NoError(t, setValueCreate(absent, "items", &untouched, func(v any) (string, error) {
		return v.(string), nil
	}))
	assert.Nil(t, untouched)
}

type exampleBlock struct {
	Name *string
}

func exampleBlockSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"block": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {Type: schema.TypeString, Optional: true},
				},
			},
		},
	}
}

func parseExampleBlock(d *schema.ResourceData, prefix string) (exampleBlock, error) {
	var value exampleBlock
	err := stringAttributeCreate(d, prefix+"name", &value.Name)
	return value, err
}

func Test_singleBlockAttributeCreate(t *testing.T) {
	d := schema.TestResourceDataRaw(t, exampleBlockSchema(), map[string]any{
		"block": []any{map[string]any{"name": "a"}},
	})

	var got *exampleBlock
	require.NoError(t, singleBlockAttributeCreate(d, "block", func(value exampleBlock) any {
		got = &value
		return nil
	}, parseExampleBlock))
	require.NotNil(t, got)
	require.NotNil(t, got.Name)
	assert.Equal(t, "a", *got.Name)

	absent := schema.TestResourceDataRaw(t, exampleBlockSchema(), map[string]any{})
	var untouched *exampleBlock
	require.NoError(t, singleBlockAttributeCreate(absent, "block", func(value exampleBlock) any {
		untouched = &value
		return nil
	}, parseExampleBlock))
	assert.Nil(t, untouched)
}

func Test_singleBlockAttributeCreate_parseError(t *testing.T) {
	d := schema.TestResourceDataRaw(t, exampleBlockSchema(), map[string]any{
		"block": []any{map[string]any{"name": "a"}},
	})
	err := singleBlockAttributeCreate(d, "block", func(exampleBlock) any { return nil }, func(*schema.ResourceData, string) (exampleBlock, error) {
		return exampleBlock{}, errors.New("parse failed")
	})
	require.EqualError(t, err, "parse failed")
}
