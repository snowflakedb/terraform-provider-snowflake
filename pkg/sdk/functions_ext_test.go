package sdk

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/datatypes"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1850370]: test parsing single
func Test_parseFunctionOrProcedureImports(t *testing.T) {
	inputs := []struct {
		rawInput string
		expected []NormalizedPath
	}{
		{"", []NormalizedPath{}},
		{`[]`, []NormalizedPath{}},
		{`[@~/abc]`, []NormalizedPath{{"~", "abc"}}},
		{`[@~/abc/def]`, []NormalizedPath{{"~", "abc/def"}}},
		{`[@"db"."sc"."st"/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@db.sc.st/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[db.sc.st/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@"db"."sc".st/abc/def]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}}},
		{`[@"db"."sc".st/abc/def, db."sc".st/abc]`, []NormalizedPath{{`"db"."sc"."st"`, "abc/def"}, {`"db"."sc"."st"`, "abc"}}},
	}

	badInputs := []struct {
		rawInput          string
		expectedErrorPart string
	}{
		{"[", "wrapping brackets not found"},
		{"]", "wrapping brackets not found"},
		{`[@~/]`, "contains empty path"},
		{`[@~]`, "cannot be split into stage and path"},
		{`[@"db"."sc"/abc]`, "contains incorrect stage location"},
		{`[@"db"/abc]`, "contains incorrect stage location"},
		{`[@"db"."sc"."st"."smth"/abc]`, "contains incorrect stage location"},
		{`[@"db/a"."sc"."st"/abc]`, "contains incorrect stage location"},
		{`[@"db"."sc"."st"/abc], @"db"."sc"/abc]`, "contains incorrect stage location"},
	}

	for _, tc := range inputs {
		tc := tc
		t.Run(fmt.Sprintf("Snowflake raw imports: %s", tc.rawInput), func(t *testing.T) {
			results, err := parseFunctionOrProcedureImports(&tc.rawInput)
			require.NoError(t, err)
			require.Equal(t, tc.expected, results)
		})
	}

	for _, tc := range badInputs {
		tc := tc
		t.Run(fmt.Sprintf("incorrect Snowflake input: %s, expecting error with: %s", tc.rawInput, tc.expectedErrorPart), func(t *testing.T) {
			_, err := parseFunctionOrProcedureImports(&tc.rawInput)
			require.Error(t, err)
			require.ErrorContains(t, err, "could not parse imports from Snowflake")
			require.ErrorContains(t, err, tc.expectedErrorPart)
		})
	}

	t.Run("Snowflake raw imports nil", func(t *testing.T) {
		results, err := parseFunctionOrProcedureImports(nil)
		require.NoError(t, err)
		require.Equal(t, []NormalizedPath{}, results)
	})
}

func Test_parseFunctionOrProcedureReturns(t *testing.T) {
	inputs := []struct {
		rawInput              string
		expectedRawDataType   string
		expectedReturnNotNull bool
	}{
		{"CHAR", "CHAR(1)", false},
		{"CHAR(1)", "CHAR(1)", false},
		{"NUMBER(30, 2)", "NUMBER(30, 2)", false},
		{"NUMBER(30,2)", "NUMBER(30, 2)", false},
		{"NUMBER(30,2) NOT NULL", "NUMBER(30, 2)", true},
		{"CHAR NOT NULL", "CHAR(1)", true},
		{"  CHAR   NOT NULL  ", "CHAR(1)", true},
		{"OBJECT", "OBJECT", false},
		{"OBJECT NOT NULL", "OBJECT", true},
	}

	badInputs := []struct {
		rawInput          string
		expectedErrorPart string
	}{
		{"", "invalid data type"},
		{"NOT NULL", "invalid data type"},
		{"CHA NOT NULL", "invalid data type"},
		{"CHA NOT NULLS", "invalid data type"},
	}

	for _, tc := range inputs {
		tc := tc
		t.Run(fmt.Sprintf("return data type raw: %s", tc.rawInput), func(t *testing.T) {
			dt, returnNotNull, err := parseFunctionOrProcedureReturns(tc.rawInput)
			require.NoError(t, err)
			require.Equal(t, tc.expectedRawDataType, dt.ToSql())
			require.Equal(t, tc.expectedReturnNotNull, returnNotNull)
		})
	}

	for _, tc := range badInputs {
		tc := tc
		t.Run(fmt.Sprintf("incorrect return data type raw: %s, expecting error with: %s", tc.rawInput, tc.expectedErrorPart), func(t *testing.T) {
			_, _, err := parseFunctionOrProcedureReturns(tc.rawInput)
			require.Error(t, err)
			require.ErrorContains(t, err, tc.expectedErrorPart)
		})
	}
}

func Test_parseFunctionOrProcedureSignature(t *testing.T) {
	inputs := []struct {
		rawInput     string
		expectedArgs []NormalizedArgument
	}{
		{"()", []NormalizedArgument{}},
		{"(abc CHAR)", []NormalizedArgument{{"abc", dataTypeChar}}},
		{"(abc CHAR(1))", []NormalizedArgument{{"abc", dataTypeChar}}},
		{"(abc CHAR(100))", []NormalizedArgument{{"abc", dataTypeChar_100}}},
		{"  (   abc CHAR(100  )  )", []NormalizedArgument{{"abc", dataTypeChar_100}}},
		{"(  abc   CHAR  )", []NormalizedArgument{{"abc", dataTypeChar}}},
		{"(abc DOUBLE PRECISION)", []NormalizedArgument{{"abc", dataTypeDoublePrecision}}},
		{"(abc double precision)", []NormalizedArgument{{"abc", dataTypeDoublePrecision}}},
		{"(abc TIMESTAMP WITHOUT TIME ZONE(5))", []NormalizedArgument{{"abc", dataTypeTimestampWithoutTimeZone_5}}},
	}

	badInputs := []struct {
		rawInput          string
		expectedErrorPart string
	}{
		{"", "can't be empty"},
		{"(abc CHAR", "wrapping parentheses not found"},
		{"abc CHAR)", "wrapping parentheses not found"},
		{"(abc)", "cannot be split into arg name, data type, and default"},
		{"(CHAR)", "cannot be split into arg name, data type, and default"},
		{"(abc CHA)", "invalid data type"},
		{"(abc CHA(123))", "invalid data type"},
		{"(abc CHAR(1) DEFAULT)", "cannot be parsed"},
		{"(abc CHAR(1) DEFAULT 'a')", "cannot be parsed"},
		// TODO [SNOW-1850370]: Snowflake currently does not return concrete data types so we can fail on them currently but it should be improved in the future
		{"(abc NUMBER(30,2))", "cannot be parsed"},
		{"(abc NUMBER(30, 2))", "cannot be parsed"},
	}

	for _, tc := range inputs {
		tc := tc
		t.Run(fmt.Sprintf("return data type raw: %s", tc.rawInput), func(t *testing.T) {
			args, err := parseFunctionOrProcedureSignature(tc.rawInput)

			require.NoError(t, err)
			require.Len(t, args, len(tc.expectedArgs))
			for i, arg := range args {
				require.Equal(t, tc.expectedArgs[i].Name, arg.Name)
				require.True(t, datatypes.AreTheSame(tc.expectedArgs[i].DataType, arg.DataType))
			}
		})
	}

	for _, tc := range badInputs {
		tc := tc
		t.Run(fmt.Sprintf("incorrect signature raw: %s, expecting error with: %s", tc.rawInput, tc.expectedErrorPart), func(t *testing.T) {
			_, err := parseFunctionOrProcedureSignature(tc.rawInput)
			require.Error(t, err)
			require.ErrorContains(t, err, "could not parse signature from Snowflake")
			require.ErrorContains(t, err, tc.expectedErrorPart)
		})
	}
}
