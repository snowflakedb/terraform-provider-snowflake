package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_sliceAwareDeepEqual(t *testing.T) {
	var nilStringSlice []string
	var nilAnySlice []any

	testCases := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{name: "same strings", a: "foo", b: "foo", expected: true},
		{name: "different strings", a: "foo", b: "bar", expected: false},
		{name: "[]string vs []any same elements", a: []string{`\N`, ""}, b: []any{`\N`, ""}, expected: true},
		{name: "[]any vs []string same elements", a: []any{"a"}, b: []string{"a"}, expected: true},
		{name: "[]string vs []any different elements", a: []string{"a"}, b: []any{"b"}, expected: false},
		{name: "nil []string vs empty []any", a: nilStringSlice, b: []any{}, expected: true},
		{name: "empty []string vs nil []any", a: []string{}, b: nilAnySlice, expected: true},
		{name: "nil []string vs empty []string", a: nilStringSlice, b: []string{}, expected: true},
		{name: "untyped nil vs empty []string", a: nil, b: []string{}, expected: true},
		{name: "empty []any vs untyped nil", a: []any{}, b: nil, expected: true},
		{name: "untyped nil vs nil []string", a: nil, b: nilStringSlice, expected: true},
		{name: "untyped nil vs non-empty slice", a: nil, b: []string{"a"}, expected: false},
		{name: "non-empty []string vs untyped nil", a: []string{"a"}, b: nil, expected: false},
		{name: "slice vs non-slice", a: []string{"a"}, b: "a", expected: false},
		{name: "both untyped nil", a: nil, b: nil, expected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, sliceAwareDeepEqual(tc.a, tc.b))
		})
	}
}

func Test_handleExternalChangesToObjectDeepEqual_nativeStringSliceVsTypeList(t *testing.T) {
	s := map[string]*schema.Schema{
		DescribeOutputAttributeName: {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"null_if": {
						Type: schema.TypeList,
						Elem: &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
		"null_if": {
			Type: schema.TypeList,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
	}

	d := schema.TestResourceDataRaw(t, s, map[string]any{
		DescribeOutputAttributeName: []any{
			map[string]any{"null_if": []any{`\N`}},
		},
		"null_if": []any{`\N`},
	})

	err := handleExternalChangesToObjectDeepEqual(d, DescribeOutputAttributeName,
		outputMapping{"null_if", "null_if", []string{`\N`}, []any{"DRIFT"}, nil},
	)
	require.NoError(t, err)
	assert.Equal(t, []any{`\N`}, d.Get("null_if"))
}
