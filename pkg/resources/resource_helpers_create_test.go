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

func Test_listValueCreate(t *testing.T) {
	listSchema := map[string]*schema.Schema{
		"items": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
	d := schema.TestResourceDataRaw(t, listSchema, map[string]any{
		"items": []any{"a", "b"},
	})

	var got []string
	require.NoError(t, listValueCreate(d, "items", &got, func(v any) (string, error) {
		return v.(string), nil
	}))
	assert.Equal(t, []string{"a", "b"}, got)

	absent := schema.TestResourceDataRaw(t, listSchema, map[string]any{})
	var untouched []string
	require.NoError(t, listValueCreate(absent, "items", &untouched, func(v any) (string, error) {
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

func exampleBlockListSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"block": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {Type: schema.TypeString, Optional: true},
				},
			},
		},
	}
}

func Test_listBlockAttributeCreate(t *testing.T) {
	d := schema.TestResourceDataRaw(t, exampleBlockListSchema(), map[string]any{
		"block": []any{
			map[string]any{"name": "a"},
			map[string]any{"name": "b"},
		},
	})

	var got []exampleBlock
	require.NoError(t, listBlockAttributeCreate(d, "block", &got, parseExampleBlock))
	require.Len(t, got, 2)
	require.NotNil(t, got[0].Name)
	assert.Equal(t, "a", *got[0].Name)
	require.NotNil(t, got[1].Name)
	assert.Equal(t, "b", *got[1].Name)

	absent := schema.TestResourceDataRaw(t, exampleBlockListSchema(), map[string]any{})
	var untouched []exampleBlock
	require.NoError(t, listBlockAttributeCreate(absent, "block", &untouched, parseExampleBlock))
	assert.Nil(t, untouched)
}

func Test_listBlockAttributeCreate_parseError(t *testing.T) {
	d := schema.TestResourceDataRaw(t, exampleBlockListSchema(), map[string]any{
		"block": []any{map[string]any{"name": "a"}},
	})
	var got []exampleBlock
	err := listBlockAttributeCreate(d, "block", &got, func(*schema.ResourceData, string) (exampleBlock, error) {
		return exampleBlock{}, errors.New("parse failed")
	})
	require.EqualError(t, err, "parse failed")
	assert.Nil(t, got)
}
