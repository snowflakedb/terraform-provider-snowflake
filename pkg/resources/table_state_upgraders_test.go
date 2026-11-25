package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_migrateTableToV2_24_0(t *testing.T) {
	defaultBlock := func(values map[string]any) []any {
		return []any{values}
	}
	expressionDefault := func(expression string) []any {
		return defaultBlock(map[string]any{"expression": expression})
	}

	testCases := []struct {
		name     string
		column   map[string]any
		expected map[string]any
	}{
		{
			name:     "text constant becomes a quoted literal",
			column:   map[string]any{"type": "VARCHAR(16)", "default": defaultBlock(map[string]any{"constant": "it's", "expression": "", "sequence": ""})},
			expected: map[string]any{"type": "VARCHAR(16)", "default": expressionDefault("'it''s'")},
		},
		{
			name:     "non-text constant stays raw",
			column:   map[string]any{"type": "BOOLEAN", "default": defaultBlock(map[string]any{"constant": "FALSE", "expression": "", "sequence": ""})},
			expected: map[string]any{"type": "BOOLEAN", "default": expressionDefault("FALSE")},
		},
		{
			name:     "constant with an unparsable type stays raw",
			column:   map[string]any{"type": "NOT A TYPE", "default": defaultBlock(map[string]any{"constant": "abc"})},
			expected: map[string]any{"type": "NOT A TYPE", "default": expressionDefault("abc")},
		},
		{
			name:     "constant without a type stays raw",
			column:   map[string]any{"default": defaultBlock(map[string]any{"constant": "abc"})},
			expected: map[string]any{"default": expressionDefault("abc")},
		},
		{
			name:     "sequence becomes NEXTVAL",
			column:   map[string]any{"type": "NUMBER(38,0)", "default": defaultBlock(map[string]any{"constant": "", "expression": "", "sequence": `"DB"."SCH"."SEQ"`})},
			expected: map[string]any{"type": "NUMBER(38,0)", "default": expressionDefault(`"DB"."SCH"."SEQ".NEXTVAL`)},
		},
		{
			name:     "expression is kept",
			column:   map[string]any{"type": "TIMESTAMP_NTZ(9)", "default": defaultBlock(map[string]any{"constant": "", "expression": "CURRENT_TIMESTAMP()", "sequence": ""})},
			expected: map[string]any{"type": "TIMESTAMP_NTZ(9)", "default": expressionDefault("CURRENT_TIMESTAMP()")},
		},
		{
			name:     "already upgraded default is kept",
			column:   map[string]any{"type": "VARCHAR(16)", "default": expressionDefault("'hello'")},
			expected: map[string]any{"type": "VARCHAR(16)", "default": expressionDefault("'hello'")},
		},
		{
			name:     "default block without values is removed",
			column:   map[string]any{"type": "VARCHAR(16)", "default": defaultBlock(map[string]any{"constant": "", "expression": "", "sequence": ""})},
			expected: map[string]any{"type": "VARCHAR(16)", "default": []any{}},
		},
		{
			name:     "nil default block is removed",
			column:   map[string]any{"type": "VARCHAR(16)", "default": []any{nil}},
			expected: map[string]any{"type": "VARCHAR(16)", "default": []any{}},
		},
		{
			name:     "column without default is unchanged",
			column:   map[string]any{"type": "VARCHAR(16)", "default": []any{}},
			expected: map[string]any{"type": "VARCHAR(16)", "default": []any{}},
		},
	}
	withBooleanDefaults := func(state map[string]any) map[string]any {
		state["is_transient"] = BooleanDefault
		state["enable_schema_evolution"] = BooleanDefault
		state["error_logging"] = BooleanDefault
		state["row_timestamp"] = BooleanDefault
		return state
	}

	maskingPolicyTestCases := []struct {
		name     string
		column   map[string]any
		expected map[string]any
	}{
		{
			name:     "masking policy name becomes a block",
			column:   map[string]any{"masking_policy": "DB.SCH.POLICY"},
			expected: map[string]any{"masking_policy": []any{map[string]any{"policy_name": "DB.SCH.POLICY", "using": []any{}}}},
		},
		{
			name:     "empty masking policy becomes an empty block list",
			column:   map[string]any{"masking_policy": ""},
			expected: map[string]any{"masking_policy": []any{}},
		},
		{
			name:     "missing masking policy becomes an empty block list",
			column:   map[string]any{},
			expected: map[string]any{"masking_policy": []any{}},
		},
		{
			name:     "already upgraded masking policy is kept",
			column:   map[string]any{"masking_policy": []any{map[string]any{"policy_name": `"DB"."SCH"."POLICY"`, "using": []any{"A", "B"}}}},
			expected: map[string]any{"masking_policy": []any{map[string]any{"policy_name": `"DB"."SCH"."POLICY"`, "using": []any{"A", "B"}}}},
		},
		{
			name:     "old default and old masking policy are both upgraded",
			column:   map[string]any{"type": "VARCHAR(16)", "default": defaultBlock(map[string]any{"constant": "x"}), "masking_policy": "DB.SCH.POLICY"},
			expected: map[string]any{"type": "VARCHAR(16)", "default": expressionDefault("'x'"), "masking_policy": []any{map[string]any{"policy_name": "DB.SCH.POLICY", "using": []any{}}}},
		},
	}
	for _, tc := range maskingPolicyTestCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.expected["not_null"] = false
			state, err := migrateTableToV2_24_0(context.Background(), map[string]any{"column": []any{tc.column}}, nil)
			require.NoError(t, err)
			assert.Equal(t, withBooleanDefaults(map[string]any{"column": []any{tc.expected}}), state)
		})
	}

	notNullTestCases := []struct {
		name     string
		column   map[string]any
		expected bool
	}{
		{name: "nullable false becomes not_null true", column: map[string]any{"nullable": false}, expected: true},
		{name: "nullable true becomes not_null false", column: map[string]any{"nullable": true}, expected: false},
		{name: "missing nullable becomes not_null false", column: map[string]any{}, expected: false},
		{name: "null nullable becomes not_null false", column: map[string]any{"nullable": nil}, expected: false},
		{name: "already upgraded not_null is kept", column: map[string]any{"not_null": true}, expected: true},
	}
	for _, tc := range notNullTestCases {
		t.Run(tc.name, func(t *testing.T) {
			state, err := migrateTableToV2_24_0(context.Background(), map[string]any{"column": []any{tc.column}}, nil)
			require.NoError(t, err)
			column := state["column"].([]any)[0].(map[string]any)
			assert.Equal(t, tc.expected, column["not_null"])
			assert.NotContains(t, column, "nullable")
		})
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.column["masking_policy"] = ""
			tc.expected["masking_policy"] = []any{}
			tc.expected["not_null"] = false
			state, err := migrateTableToV2_24_0(context.Background(), map[string]any{"column": []any{tc.column}}, nil)
			require.NoError(t, err)
			assert.Equal(t, withBooleanDefaults(map[string]any{"column": []any{tc.expected}}), state)
		})
	}

	t.Run("nil state", func(t *testing.T) {
		state, err := migrateTableToV2_24_0(context.Background(), nil, nil)
		require.NoError(t, err)
		assert.Nil(t, state)
	})

	t.Run("state without columns", func(t *testing.T) {
		state, err := migrateTableToV2_24_0(context.Background(), map[string]any{"name": "T"}, nil)
		require.NoError(t, err)
		assert.Equal(t, withBooleanDefaults(map[string]any{"name": "T"}), state)
	})

	t.Run("unexpected column shape is left unchanged", func(t *testing.T) {
		state, err := migrateTableToV2_24_0(context.Background(), map[string]any{"column": []any{"not a map"}}, nil)
		require.NoError(t, err)
		assert.Equal(t, withBooleanDefaults(map[string]any{"column": []any{"not a map"}}), state)
	})

	t.Run("boolean string fields keep their values", func(t *testing.T) {
		state, err := migrateTableToV2_24_0(context.Background(), map[string]any{
			"is_transient":            BooleanTrue,
			"enable_schema_evolution": BooleanFalse,
			"error_logging":           nil,
			"row_timestamp":           "",
		}, nil)
		require.NoError(t, err)
		assert.Equal(t, map[string]any{
			"is_transient":            BooleanTrue,
			"enable_schema_evolution": BooleanFalse,
			"error_logging":           BooleanDefault,
			"row_timestamp":           BooleanDefault,
		}, state)
	})
}

func Test_suppressColumnDefaultSequenceQuoting(t *testing.T) {
	testCases := []struct {
		oldValue string
		newValue string
		expected bool
	}{
		{oldValue: "DB.SCH.SEQ.NEXTVAL", newValue: `"DB"."SCH"."SEQ".NEXTVAL`, expected: true},
		{oldValue: `"DB"."SCH"."SEQ".NEXTVAL`, newValue: `"DB"."SCH"."SEQ".nextval`, expected: true},
		{oldValue: "DB.SCH.SEQ.NEXTVAL", newValue: `"DB"."SCH"."OTHER".NEXTVAL`, expected: false},
		{oldValue: "DB.SCH.SEQ.NEXTVAL", newValue: "CURRENT_TIMESTAMP()", expected: false},
		{oldValue: "'a'", newValue: "'a'", expected: false},
		{oldValue: ".NEXTVAL", newValue: ".NEXTVAL", expected: false},
	}
	for _, tc := range testCases {
		t.Run(tc.oldValue+" vs "+tc.newValue, func(t *testing.T) {
			assert.Equal(t, tc.expected, suppressColumnDefaultSequenceQuoting("", tc.oldValue, tc.newValue, nil))
		})
	}
}

func Test_validateTableColumnDefaultExpression(t *testing.T) {
	// "" is the empty string default, see Test_getColumnDefaultExpression.
	assert.False(t, validateTableColumnDefaultExpression("", nil).HasError())
	assert.True(t, validateTableColumnDefaultExpression("   ", nil).HasError())
	assert.False(t, validateTableColumnDefaultExpression("''", nil).HasError())
	assert.False(t, validateTableColumnDefaultExpression("CURRENT_TIMESTAMP()", nil).HasError())
}

func Test_getColumnDefaultExpression(t *testing.T) {
	tests := map[string]struct {
		defaults           []any
		expectedExpression string
		expectedOk         bool
	}{
		"no default block":             {defaults: []any{}, expectedOk: false},
		"expression":                   {defaults: []any{map[string]any{"expression": "'pending'"}}, expectedExpression: "'pending'", expectedOk: true},
		"function":                     {defaults: []any{map[string]any{"expression": "CURRENT_TIMESTAMP()"}}, expectedExpression: "CURRENT_TIMESTAMP()", expectedOk: true},
		"empty string literal":         {defaults: []any{map[string]any{"expression": "''"}}, expectedExpression: "''", expectedOk: true},
		"empty expression":             {defaults: []any{map[string]any{"expression": ""}}, expectedExpression: "''", expectedOk: true},
		"block decoded as nil element": {defaults: []any{nil}, expectedExpression: "''", expectedOk: true},
		"block without expression key": {defaults: []any{map[string]any{}}, expectedExpression: "''", expectedOk: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			expression, ok := getColumnDefaultExpression(tc.defaults)
			assert.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedExpression, expression)
		})
	}
}

func Test_suppressColumnDefaultExpressionDifferences(t *testing.T) {
	tests := []struct {
		oldValue string
		newValue string
		expected bool
	}{
		{oldValue: "''", newValue: "", expected: true},
		{oldValue: "", newValue: "''", expected: true},
		{oldValue: "''", newValue: "''", expected: true},
		{oldValue: "''", newValue: "'a'", expected: false},
		{oldValue: "'a'", newValue: "", expected: false},
		{oldValue: "CURRENT_TIMESTAMP()", newValue: "", expected: false},
		{oldValue: `DB.SCHEMA.SEQ.NEXTVAL`, newValue: `"DB"."SCHEMA"."SEQ".NEXTVAL`, expected: true},
	}
	for _, tc := range tests {
		t.Run(tc.oldValue+" vs "+tc.newValue, func(t *testing.T) {
			assert.Equal(t, tc.expected, suppressColumnDefaultExpressionDifferences("", tc.oldValue, tc.newValue, nil))
		})
	}
}
