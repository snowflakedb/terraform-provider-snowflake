package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseFunctionOrProcedureShowArguments(t *testing.T) {
	mustParse := func(sql string) datatypes.DataType {
		dataType, err := datatypes.ParseDataType(sql)
		require.NoError(t, err)
		return dataType
	}

	testCases := []struct {
		name                  string
		arguments             string
		expectedArgumentsRaw  []string
		expectedArguments     []datatypes.DataType
		expectedReturnType    datatypes.DataType
		expectedReturnTypeRaw string
	}{
		{
			name:                  "no arguments",
			arguments:             "NAME() RETURN VARCHAR",
			expectedArgumentsRaw:  []string{},
			expectedArguments:     []datatypes.DataType{},
			expectedReturnType:    mustParse("VARCHAR"),
			expectedReturnTypeRaw: "VARCHAR",
		},
		{
			name:                  "multiple arguments",
			arguments:             "NAME(NUMBER, VARCHAR, VECTOR(INT, 16)) RETURN TABLE ()",
			expectedArgumentsRaw:  []string{"NUMBER", "VARCHAR", "VECTOR(INT, 16)"},
			expectedArguments:     []datatypes.DataType{mustParse("NUMBER"), mustParse("VARCHAR"), mustParse("VECTOR(INT, 16)")},
			expectedReturnType:    mustParse("TABLE ()"),
			expectedReturnTypeRaw: "TABLE ()",
		},
		{
			name:                  "unparseable argument and return type fall back to nil",
			arguments:             "NAME(UNKNOWN_TYPE, FLOAT) RETURN UNKNOWN_TYPE",
			expectedArgumentsRaw:  []string{"UNKNOWN_TYPE", "FLOAT"},
			expectedArguments:     []datatypes.DataType{nil, mustParse("FLOAT")},
			expectedReturnType:    nil,
			expectedReturnTypeRaw: "UNKNOWN_TYPE",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			showArguments, err := parseFunctionOrProcedureShowArguments("NAME", tc.arguments)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedArgumentsRaw, showArguments.ArgumentsRaw)
			require.Len(t, showArguments.Arguments, len(tc.expectedArguments))
			for i, expected := range tc.expectedArguments {
				assert.True(t, datatypes.AreTheSame(expected, showArguments.Arguments[i]), "argument %d: expected %v, got %v", i, expected, showArguments.Arguments[i])
			}
			assert.True(t, datatypes.AreTheSame(tc.expectedReturnType, showArguments.ReturnType))
			assert.Equal(t, tc.expectedReturnTypeRaw, showArguments.ReturnTypeRaw)
		})
	}
}
