package resources

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

type testResourceValueSetter struct {
	internalMap map[string]any
}

func newTestResourceValueSetter() *testResourceValueSetter {
	return &testResourceValueSetter{
		internalMap: make(map[string]any),
	}
}

func (s *testResourceValueSetter) Set(key string, value any) error {
	s.internalMap[key] = value
	return nil
}

func Test_handleColumns(t *testing.T) {
	testCases := []struct {
		InputColumns          []sdk.ViewDetails
		InputPolicyReferences []sdk.PolicyReference
		Expected              map[string]any
	}{
		{
			InputColumns:          []sdk.ViewDetails{},
			InputPolicyReferences: []sdk.PolicyReference{},
			Expected: map[string]any{
				"column": nil,
			},
		},
		{
			InputColumns: []sdk.ViewDetails{
				{
					Name:    "name",
					Comment: nil,
				},
			},
			InputPolicyReferences: []sdk.PolicyReference{},
			Expected: map[string]any{
				"column": []map[string]any{
					{
						"column_name": "name",
						"comment":     nil,
					},
				},
			},
		},
		{
			InputColumns: []sdk.ViewDetails{
				{
					Name:    "name",
					Comment: sdk.String("comment"),
				},
			},
			InputPolicyReferences: []sdk.PolicyReference{},
			Expected: map[string]any{
				"column": []map[string]any{
					{
						"column_name": "name",
						"comment":     "comment",
					},
				},
			},
		},
		{
			InputColumns: []sdk.ViewDetails{
				{
					Name:    "name",
					Comment: sdk.String("comment"),
				},
				{
					Name:    "name2",
					Comment: sdk.String("comment2"),
				},
			},
			InputPolicyReferences: []sdk.PolicyReference{},
			Expected: map[string]any{
				"column": []map[string]any{
					{
						"column_name": "name",
						"comment":     "comment",
					},
					{
						"column_name": "name2",
						"comment":     "comment2",
					},
				},
			},
		},
		{
			InputColumns: []sdk.ViewDetails{
				{
					Name:    "name",
					Comment: sdk.String("comment"),
				},
				{
					Name:    "name2",
					Comment: sdk.String("comment2"),
				},
			},
			InputPolicyReferences: []sdk.PolicyReference{
				{
					PolicyDb:      sdk.String("db"),
					PolicySchema:  sdk.String("sch"),
					PolicyName:    "policyName",
					PolicyKind:    sdk.PolicyKindProjectionPolicy,
					RefColumnName: sdk.String("name"),
				},
			},
			Expected: map[string]any{
				"column": []map[string]any{
					{
						"column_name": "name",
						"comment":     "comment",
						"projection_policy": []map[string]any{
							{
								"policy_name": sdk.NewSchemaObjectIdentifier("db", "sch", "policyName").FullyQualifiedName(),
							},
						},
					},
					{
						"column_name": "name2",
						"comment":     "comment2",
					},
				},
			},
		},
		{
			InputColumns: []sdk.ViewDetails{
				{
					Name:    "name",
					Comment: sdk.String("comment"),
				},
				{
					Name:    "name2",
					Comment: sdk.String("comment2"),
				},
			},
			InputPolicyReferences: []sdk.PolicyReference{
				{
					PolicyDb:      sdk.String("db"),
					PolicySchema:  sdk.String("sch"),
					PolicyName:    "policyName",
					PolicyKind:    sdk.PolicyKindProjectionPolicy,
					RefColumnName: sdk.String("name"),
				},
				{
					PolicyDb:          sdk.String("db"),
					PolicySchema:      sdk.String("sch"),
					PolicyName:        "policyName2",
					PolicyKind:        sdk.PolicyKindMaskingPolicy,
					RefColumnName:     sdk.String("name"),
					RefArgColumnNames: sdk.String("[one,two]"),
				},
			},
			Expected: map[string]any{
				"column": []map[string]any{
					{
						"column_name": "name",
						"comment":     "comment",
						"projection_policy": []map[string]any{
							{
								"policy_name": sdk.NewSchemaObjectIdentifier("db", "sch", "policyName").FullyQualifiedName(),
							},
						},
						"masking_policy": []map[string]any{
							{
								"policy_name": sdk.NewSchemaObjectIdentifier("db", "sch", "policyName2").FullyQualifiedName(),
								"using":       []string{"name", "one", "two"},
							},
						},
					},
					{
						"column_name": "name2",
						"comment":     "comment2",
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("handle columns(%d): %v - %v", i, tc.InputColumns, tc.InputPolicyReferences), func(t *testing.T) {
			valueSetter := newTestResourceValueSetter()
			err := handleColumns(valueSetter, tc.InputColumns, tc.InputPolicyReferences)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, valueSetter.internalMap)
		})
	}
}

func Test_extractColumns(t *testing.T) {
	testCases := []struct {
		Input    any
		Expected []sdk.ViewColumnRequest
		Error    string
	}{
		{
			Input: "",
			Error: "unable to extract columns, input is either nil or non expected type (string): ",
		},
		{
			Input: nil,
			Error: "unable to extract columns, input is either nil or non expected type (<nil>): <nil>",
		},
		{
			Input: []any{""},
			Error: "unable to extract column, non expected type of string: ",
		},
		{
			Input: []any{
				map[string]any{},
			},
			Error: "unable to extract column, missing column_name key in column",
		},
		{
			Input: []any{
				map[string]any{
					"column_name": "abc",
				},
			},
			Expected: []sdk.ViewColumnRequest{
				*sdk.NewViewColumnRequest("abc"),
			},
		},
		{
			Input: []any{
				map[string]any{
					"column_name": "abc",
				},
				map[string]any{
					"column_name": "cba",
				},
			},
			Expected: []sdk.ViewColumnRequest{
				*sdk.NewViewColumnRequest("abc"),
				*sdk.NewViewColumnRequest("cba"),
			},
		},
		{
			Input: []any{
				map[string]any{
					"column_name": "abc",
					"projection_policy": []any{
						map[string]any{
							"policy_name": "db.sch.proj",
						},
					},
					"masking_policy": []any{
						map[string]any{
							"policy_name": "db.sch.mask",
							"using":       []any{"one", "two"},
						},
					},
				},
				map[string]any{
					"column_name": "cba",
				},
			},
			Expected: []sdk.ViewColumnRequest{
				*sdk.NewViewColumnRequest("abc").
					WithProjectionPolicy(*sdk.NewViewColumnProjectionPolicyRequest(sdk.NewSchemaObjectIdentifier("db", "sch", "proj"))).
					WithMaskingPolicy(*sdk.NewViewColumnMaskingPolicyRequest(sdk.NewSchemaObjectIdentifier("db", "sch", "mask")).WithUsing([]sdk.Column{{Value: "one"}, {Value: "two"}})),
				*sdk.NewViewColumnRequest("cba"),
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d: %s", i, tc.Input), func(t *testing.T) {
			req, err := extractColumns(tc.Input)

			if tc.Error != "" {
				assert.Nil(t, req)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.Error)
			} else {
				assert.True(t, reflect.DeepEqual(tc.Expected, req))
				assert.NoError(t, err)
			}
		})
	}
}

func Test_extractPolicyWithColumnsList(t *testing.T) {
	testCases := []struct {
		Input           any
		ColumnKey       string
		ExpectedId      sdk.SchemaObjectIdentifier
		ExpectedColumns []sdk.Column
		Error           string
	}{
		{
			Input: []any{
				map[string]any{
					"policy_name": "db.sch.pol",
					"using":       []any{"one", "two"},
				},
			},
			ColumnKey: "non-existing",
			Error:     "unable to extract policy with column list, unable to find columnsKey: non-existing",
		},
		{
			Input: []any{
				map[string]any{
					"policy_name": "db.sch.pol",
				},
			},
			ColumnKey: "using",
			Error:     "unable to extract policy with column list, unable to find columnsKey: using",
		},
		{
			Input: []any{
				map[string]any{
					"policy_name": "db.sch.pol",
					"using":       []any{"one", "two"},
				},
			},
			ColumnKey:       "using",
			ExpectedId:      sdk.NewSchemaObjectIdentifier("db", "sch", "pol"),
			ExpectedColumns: []sdk.Column{{Value: "one"}, {Value: "two"}},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d: %s", i, tc.Input), func(t *testing.T) {
			id, cols, err := extractPolicyWithColumnsList(tc.Input, tc.ColumnKey)

			if tc.Error != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.ExpectedId, id)
				assert.Equal(t, tc.ExpectedColumns, cols)
			}
		})
	}
}
