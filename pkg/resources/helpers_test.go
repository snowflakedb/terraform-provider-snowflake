package resources_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func Test_GetPropertyAsPointer(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"integer": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"second_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"third_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"second_string": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"third_string": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"boolean": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"second_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"third_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}, map[string]any{
		"integer":        123,
		"second_integer": 0,
		"string":         "some string",
		"second_string":  "",
		"boolean":        true,
		"second_boolean": false,
		"invalid":        true,
	})

	assert.Equal(t, 123, *resources.GetPropertyAsPointer[int](d, "integer"))
	assert.Equal(t, "some string", *resources.GetPropertyAsPointer[string](d, "string"))
	assert.True(t, *resources.GetPropertyAsPointer[bool](d, "boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "invalid"))

	assert.Equal(t, 123, *resources.GetPropertyAsPointer[int](d, "integer"))
	assert.Nil(t, resources.GetPropertyAsPointer[int](d, "second_integer"))
	assert.Nil(t, resources.GetPropertyAsPointer[int](d, "third_integer"))
	assert.Equal(t, "some string", *resources.GetPropertyAsPointer[string](d, "string"))
	assert.Nil(t, resources.GetPropertyAsPointer[string](d, "second_integer"))
	assert.Nil(t, resources.GetPropertyAsPointer[string](d, "third_string"))
	assert.True(t, *resources.GetPropertyAsPointer[bool](d, "boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "second_boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "third_boolean"))
	assert.Nil(t, resources.GetPropertyAsPointer[bool](d, "invalid"))
}

// TODO [SNOW-1511594]: provide TestResourceDataRaw with working GetRawConfig()
func Test_GetConfigPropertyAsPointerAllowingZeroValue(t *testing.T) {
	t.Skip("TestResourceDataRaw does not set up the ResourceData correctly - GetRawConfig is nil")
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"integer": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"second_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"third_integer": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"second_string": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"third_string": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"boolean": {
			Type:     schema.TypeBool,
			Required: true,
		},
		"second_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"third_boolean": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}, map[string]any{
		"integer":        123,
		"second_integer": 0,
		"string":         "some string",
		"second_string":  "",
		"boolean":        true,
		"second_boolean": false,
		"invalid":        true,
	})

	assert.Equal(t, 123, *resources.GetConfigPropertyAsPointerAllowingZeroValue[int](d, "integer"))
	assert.Equal(t, 0, *resources.GetConfigPropertyAsPointerAllowingZeroValue[int](d, "second_integer"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[int](d, "third_integer"))
	assert.Equal(t, "some string", *resources.GetConfigPropertyAsPointerAllowingZeroValue[string](d, "string"))
	assert.Equal(t, "", *resources.GetConfigPropertyAsPointerAllowingZeroValue[string](d, "second_integer"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[string](d, "third_string"))
	assert.True(t, *resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "boolean"))
	assert.False(t, *resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "second_boolean"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "third_boolean"))
	assert.Nil(t, resources.GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "invalid"))
}

func TestListDiff(t *testing.T) {
	testCases := []struct {
		Name    string
		Before  []any
		After   []any
		Added   []any
		Removed []any
	}{
		{
			Name:    "no changes",
			Before:  []any{1, 2, 3, 4},
			After:   []any{1, 2, 3, 4},
			Removed: []any{},
			Added:   []any{},
		},
		{
			Name:    "only removed",
			Before:  []any{1, 2, 3, 4},
			After:   []any{},
			Removed: []any{1, 2, 3, 4},
			Added:   []any{},
		},
		{
			Name:    "only added",
			Before:  []any{},
			After:   []any{1, 2, 3, 4},
			Removed: []any{},
			Added:   []any{1, 2, 3, 4},
		},
		{
			Name:    "added repeated items",
			Before:  []any{2},
			After:   []any{1, 2, 1},
			Removed: []any{},
			Added:   []any{1, 1},
		},
		{
			Name:    "removed repeated items",
			Before:  []any{1, 2, 1},
			After:   []any{2},
			Removed: []any{1, 1},
			Added:   []any{},
		},
		{
			Name:    "simple diff: ints",
			Before:  []any{1, 2, 3, 4, 5, 6, 7, 8, 9},
			After:   []any{1, 3, 5, 7, 9, 12, 13, 14},
			Removed: []any{2, 4, 6, 8},
			Added:   []any{12, 13, 14},
		},
		{
			Name:    "simple diff: strings",
			Before:  []any{"one", "two", "three", "four"},
			After:   []any{"five", "two", "four", "six"},
			Removed: []any{"one", "three"},
			Added:   []any{"five", "six"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			added, removed := resources.ListDiff(tc.Before, tc.After)
			assert.Equal(t, tc.Added, added)
			assert.Equal(t, tc.Removed, removed)
		})
	}
}

func TestListDiffWithCommonItems(t *testing.T) {
	testCases := []struct {
		Name    string
		Before  []any
		After   []any
		Added   []any
		Removed []any
		Common  []any
	}{
		{
			Name:    "no changes",
			Before:  []any{1, 2, 3, 4},
			After:   []any{1, 2, 3, 4},
			Removed: []any{},
			Added:   []any{},
			Common:  []any{1, 2, 3, 4},
		},
		{
			Name:    "only removed",
			Before:  []any{1, 2, 3, 4},
			After:   []any{},
			Removed: []any{1, 2, 3, 4},
			Added:   []any{},
			Common:  []any{},
		},
		{
			Name:    "only added",
			Before:  []any{},
			After:   []any{1, 2, 3, 4},
			Removed: []any{},
			Added:   []any{1, 2, 3, 4},
			Common:  []any{},
		},
		{
			Name:    "added repeated items",
			Before:  []any{2},
			After:   []any{1, 2, 1},
			Removed: []any{},
			Added:   []any{1, 1},
			Common:  []any{2},
		},
		{
			Name:    "removed repeated items",
			Before:  []any{1, 2, 1},
			After:   []any{2},
			Removed: []any{1, 1},
			Added:   []any{},
			Common:  []any{2},
		},
		{
			Name:    "simple diff: ints",
			Before:  []any{1, 2, 3, 4, 5, 6, 7, 8, 9},
			After:   []any{1, 3, 5, 7, 9, 12, 13, 14},
			Removed: []any{2, 4, 6, 8},
			Added:   []any{12, 13, 14},
			Common:  []any{1, 3, 5, 7, 9},
		},
		{
			Name:    "simple diff: strings",
			Before:  []any{"one", "two", "three", "four"},
			After:   []any{"five", "two", "four", "six"},
			Removed: []any{"one", "three"},
			Added:   []any{"five", "six"},
			Common:  []any{"two", "four"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			added, removed, common := resources.ListDiffWithCommonItems(tc.Before, tc.After)
			assert.Equal(t, tc.Added, added)
			assert.Equal(t, tc.Removed, removed)
			assert.Equal(t, tc.Common, common)
		})
	}
}
