package datatypes

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseDataType_Number(t *testing.T) {
	type test struct {
		input                  string
		expectedPrecision      int
		expectedScale          int
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedPrecision:      DefaultNumberPrecision,
			expectedScale:          DefaultNumberScale,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		{input: "NUMBER(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale, expectedUnderlyingType: "NUMBER"},
		{input: "NUMBER(30, 2)", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "NUMBER"},
		{input: "dec(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale, expectedUnderlyingType: "DEC"},
		{input: "dec(30, 2)", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "DEC"},
		{input: "decimal(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale, expectedUnderlyingType: "DECIMAL"},
		{input: "decimal(30, 2)", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "DECIMAL"},
		{input: "NuMeRiC(30)", expectedPrecision: 30, expectedScale: DefaultNumberScale, expectedUnderlyingType: "NUMERIC"},
		{input: "NuMeRiC(30, 2)", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "NUMERIC"},
		{input: "NUMBER(   30   ,  2   )", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "NUMBER"},
		{input: "    NUMBER   (   30   ,  2   )    ", expectedPrecision: 30, expectedScale: 2, expectedUnderlyingType: "NUMBER"},
		{input: fmt.Sprintf("NUMBER(%d)", DefaultNumberPrecision), expectedPrecision: DefaultNumberPrecision, expectedScale: DefaultNumberScale, expectedUnderlyingType: "NUMBER"},
		{input: fmt.Sprintf("NUMBER(%d, %d)", DefaultNumberPrecision, DefaultNumberScale), expectedPrecision: DefaultNumberPrecision, expectedScale: DefaultNumberScale, expectedUnderlyingType: "NUMBER"},

		defaults("NUMBER"),
		defaults("DEC"),
		defaults("DECIMAL"),
		defaults("NUMERIC"),
		defaults("   NUMBER   "),

		defaults("INT"),
		defaults("INTEGER"),
		defaults("BIGINT"),
		defaults("SMALLINT"),
		defaults("TINYINT"),
		defaults("BYTEINT"),
		defaults("int"),
		defaults("integer"),
		defaults("bigint"),
		defaults("smallint"),
		defaults("tinyint"),
		defaults("byteint"),
	}

	negativeTestCases := []test{
		negative("other(1, 2)"),
		negative("other(1)"),
		negative("other"),
		negative("NUMBER()"),
		negative("NUMBER(x)"),
		negative(fmt.Sprintf("NUMBER(%d, x)", DefaultNumberPrecision)),
		negative(fmt.Sprintf("NUMBER(x, %d)", DefaultNumberScale)),
		negative("NUMBER(1, 2, 3)"),
		negative("NUMBER("),
		negative("NUMBER)"),
		negative("NUM BER"),
		negative("INT(30)"),
		negative("INT(30, 2)"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &NumberDataType{}, parsed)

			assert.Equal(t, tc.expectedPrecision, parsed.(*NumberDataType).precision)
			assert.Equal(t, tc.expectedScale, parsed.(*NumberDataType).scale)
			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*NumberDataType).underlyingType)

			assert.Equal(t, NumberLegacyDataType, parsed.ToLegacyDataTypeSql())
			if slices.Contains(NumberDataTypeSubTypes, parsed.(*NumberDataType).underlyingType) {
				assert.Equal(t, parsed.(*NumberDataType).underlyingType, parsed.ToSql())
			} else {
				assert.Equal(t, fmt.Sprintf("%s(%d, %d)", parsed.(*NumberDataType).underlyingType, parsed.(*NumberDataType).precision, parsed.(*NumberDataType).scale), parsed.ToSql())
			}
			assert.Equal(t, fmt.Sprintf("%s(%d,%d)", NumberLegacyDataType, parsed.(*NumberDataType).precision, parsed.(*NumberDataType).scale), parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Float(t *testing.T) {
	type test struct {
		input                  string
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		defaults("   FLOAT   "),
		defaults("FLOAT"),
		defaults("FLOAT4"),
		defaults("FLOAT8"),
		defaults("DOUBLE PRECISION"),
		defaults("DOUBLE"),
		defaults("REAL"),
		defaults("float"),
		defaults("float4"),
		defaults("float8"),
		defaults("double precision"),
		defaults("double"),
		defaults("real"),
	}

	negativeTestCases := []test{
		negative("FLOAT(38, 0)"),
		negative("FLOAT(38, 2)"),
		negative("FLOAT(38)"),
		negative("FLOAT()"),
		negative("F L O A T"),
		negative("other"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &FloatDataType{}, parsed)

			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*FloatDataType).underlyingType)

			assert.Equal(t, FloatLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, tc.expectedUnderlyingType, parsed.ToSql())
			assert.Equal(t, FloatLegacyDataType, parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Text(t *testing.T) {
	type test struct {
		input                  string
		expectedLength         int
		expectedUnderlyingType string
	}
	defaultsVarchar := func(input string) test {
		return test{
			input:                  input,
			expectedLength:         DefaultVarcharLength,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	defaultsChar := func(input string) test {
		return test{
			input:                  input,
			expectedLength:         DefaultCharLength,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		{input: "VARCHAR(30)", expectedLength: 30, expectedUnderlyingType: "VARCHAR"},
		{input: "string(30)", expectedLength: 30, expectedUnderlyingType: "STRING"},
		{input: "VARCHAR(   30   )", expectedLength: 30, expectedUnderlyingType: "VARCHAR"},
		{input: "    VARCHAR   (   30   )    ", expectedLength: 30, expectedUnderlyingType: "VARCHAR"},
		{input: fmt.Sprintf("VARCHAR(%d)", DefaultVarcharLength), expectedLength: DefaultVarcharLength, expectedUnderlyingType: "VARCHAR"},

		{input: "CHAR(30)", expectedLength: 30, expectedUnderlyingType: "CHAR"},
		{input: "character(30)", expectedLength: 30, expectedUnderlyingType: "CHARACTER"},
		{input: "CHAR(   30   )", expectedLength: 30, expectedUnderlyingType: "CHAR"},
		{input: "    CHAR   (   30   )    ", expectedLength: 30, expectedUnderlyingType: "CHAR"},
		{input: fmt.Sprintf("CHAR(%d)", DefaultCharLength), expectedLength: DefaultCharLength, expectedUnderlyingType: "CHAR"},

		defaultsVarchar("   VARCHAR   "),
		defaultsVarchar("VARCHAR"),
		defaultsVarchar("STRING"),
		defaultsVarchar("TEXT"),
		defaultsVarchar("NVARCHAR"),
		defaultsVarchar("NVARCHAR2"),
		defaultsVarchar("CHAR VARYING"),
		defaultsVarchar("NCHAR VARYING"),
		defaultsVarchar("varchar"),
		defaultsVarchar("string"),
		defaultsVarchar("text"),
		defaultsVarchar("nvarchar"),
		defaultsVarchar("nvarchar2"),
		defaultsVarchar("char varying"),
		defaultsVarchar("nchar varying"),

		defaultsChar("   CHAR   "),
		defaultsChar("CHAR"),
		defaultsChar("CHARACTER"),
		defaultsChar("NCHAR"),
		defaultsChar("char"),
		defaultsChar("character"),
		defaultsChar("nchar"),
	}

	negativeTestCases := []test{
		negative("other(1, 2)"),
		negative("other(1)"),
		negative("other"),
		negative("VARCHAR()"),
		negative("VARCHAR(x)"),
		negative("VARCHAR(   )"),
		negative("CHAR()"),
		negative("CHAR(x)"),
		negative("CHAR(   )"),
		negative("VARCHAR(1, 2)"),
		negative("VARCHAR("),
		negative("VARCHAR)"),
		negative("VAR CHAR"),
		negative("CHAR(1, 2)"),
		negative("CHAR("),
		negative("CHAR)"),
		negative("CH AR"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &TextDataType{}, parsed)

			assert.Equal(t, tc.expectedLength, parsed.(*TextDataType).length)
			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*TextDataType).underlyingType)

			assert.Equal(t, VarcharLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", parsed.(*TextDataType).underlyingType, parsed.(*TextDataType).length), parsed.ToSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", VarcharLegacyDataType, parsed.(*TextDataType).length), parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Binary(t *testing.T) {
	type test struct {
		input                  string
		expectedSize           int
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedSize:           DefaultBinarySize,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		{input: "BINARY(30)", expectedSize: 30, expectedUnderlyingType: "BINARY"},
		{input: "varbinary(30)", expectedSize: 30, expectedUnderlyingType: "VARBINARY"},
		{input: "BINARY(   30   )", expectedSize: 30, expectedUnderlyingType: "BINARY"},
		{input: "    BINARY   (   30   )    ", expectedSize: 30, expectedUnderlyingType: "BINARY"},
		{input: fmt.Sprintf("BINARY(%d)", DefaultBinarySize), expectedSize: DefaultBinarySize, expectedUnderlyingType: "BINARY"},

		defaults("   BINARY   "),
		defaults("BINARY"),
		defaults("VARBINARY"),
		defaults("binary"),
		defaults("varbinary"),
	}

	negativeTestCases := []test{
		negative("other(1, 2)"),
		negative("other(1)"),
		negative("other"),
		negative("BINARY()"),
		negative("BINARY(x)"),
		negative("BINARY(   )"),
		negative("BINARY(1, 2)"),
		negative("BINARY("),
		negative("BINARY)"),
		negative("BIN ARY"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &BinaryDataType{}, parsed)

			assert.Equal(t, tc.expectedSize, parsed.(*BinaryDataType).size)
			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*BinaryDataType).underlyingType)

			assert.Equal(t, BinaryLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", parsed.(*BinaryDataType).underlyingType, parsed.(*BinaryDataType).size), parsed.ToSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", BinaryLegacyDataType, parsed.(*BinaryDataType).size), parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Boolean(t *testing.T) {
	type test struct {
		input                  string
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		defaults("   BOOLEAN   "),
		defaults("BOOLEAN"),
		defaults("boolean"),
	}

	negativeTestCases := []test{
		negative("BOOLEAN(38, 0)"),
		negative("BOOLEAN(38, 2)"),
		negative("BOOLEAN(38)"),
		negative("BOOLEAN()"),
		negative("BOOL"),
		negative("bool"),
		negative("B O O L E A N"),
		negative("other"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &BooleanDataType{}, parsed)

			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*BooleanDataType).underlyingType)

			assert.Equal(t, BooleanLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, tc.expectedUnderlyingType, parsed.ToSql())
			assert.Equal(t, BooleanLegacyDataType, parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Date(t *testing.T) {
	type test struct {
		input                  string
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		defaults("   DATE   "),
		defaults("DATE"),
		defaults("date"),
	}

	negativeTestCases := []test{
		negative("DATE(38, 0)"),
		negative("DATE(38, 2)"),
		negative("DATE(38)"),
		negative("DATE()"),
		negative("D A T E"),
		negative("other"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &DateDataType{}, parsed)

			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*DateDataType).underlyingType)

			assert.Equal(t, DateLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, tc.expectedUnderlyingType, parsed.ToSql())
			assert.Equal(t, DateLegacyDataType, parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Time(t *testing.T) {
	type test struct {
		input                  string
		expectedPrecision      int
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedPrecision:      DefaultTimePrecision,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		defaults("   TIME   "),
		defaults("TIME"),
		defaults("time"),
		{input: "TIME(5)", expectedPrecision: 5, expectedUnderlyingType: "TIME"},
		{input: "time(5)", expectedPrecision: 5, expectedUnderlyingType: "TIME"},
	}

	negativeTestCases := []test{
		negative("TIME(38, 0)"),
		negative("TIME(38, 2)"),
		negative("TIME()"),
		negative("T I M E"),
		negative("other"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &TimeDataType{}, parsed)

			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*TimeDataType).underlyingType)
			assert.Equal(t, tc.expectedPrecision, parsed.(*TimeDataType).precision)

			assert.Equal(t, TimeLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", tc.expectedUnderlyingType, tc.expectedPrecision), parsed.ToSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", TimeLegacyDataType, tc.expectedPrecision), parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_TimestampLtz(t *testing.T) {
	type test struct {
		input                  string
		expectedPrecision      int
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedPrecision:      DefaultTimestampPrecision,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		{input: "TIMESTAMP_LTZ(4)", expectedPrecision: 4, expectedUnderlyingType: "TIMESTAMP_LTZ"},
		{input: "timestamp with local time zone(5)", expectedPrecision: 5, expectedUnderlyingType: "TIMESTAMP WITH LOCAL TIME ZONE"},
		{input: "TIMESTAMP_LTZ(   2   )", expectedPrecision: 2, expectedUnderlyingType: "TIMESTAMP_LTZ"},
		{input: "    TIMESTAMP_LTZ   (   7   )    ", expectedPrecision: 7, expectedUnderlyingType: "TIMESTAMP_LTZ"},
		{input: fmt.Sprintf("TIMESTAMP_LTZ(%d)", DefaultTimestampPrecision), expectedPrecision: DefaultTimestampPrecision, expectedUnderlyingType: "TIMESTAMP_LTZ"},

		defaults("   TIMESTAMP_LTZ   "),
		defaults("TIMESTAMP_LTZ"),
		defaults("TIMESTAMPLTZ"),
		defaults("TIMESTAMP WITH LOCAL TIME ZONE"),
		defaults("timestamp_ltz"),
		defaults("timestampltz"),
		defaults("timestamp with local time zone"),
	}

	negativeTestCases := []test{
		negative("TIMESTAMP_LTZ(38, 0)"),
		negative("TIMESTAMP_LTZ(38, 2)"),
		negative("TIMESTAMP_LTZ()"),
		negative("T I M E S T A M P _ L T Z"),
		negative("other"),
		negative("other(3)"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &TimestampLtzDataType{}, parsed)

			assert.Equal(t, tc.expectedPrecision, parsed.(*TimestampLtzDataType).precision)
			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*TimestampLtzDataType).underlyingType)

			assert.Equal(t, TimestampLtzLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", parsed.(*TimestampLtzDataType).underlyingType, parsed.(*TimestampLtzDataType).precision), parsed.ToSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", TimestampLtzLegacyDataType, parsed.(*TimestampLtzDataType).precision), parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_TimestampNtz(t *testing.T) {
	type test struct {
		input                  string
		expectedPrecision      int
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedPrecision:      DefaultTimestampPrecision,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		{input: "TIMESTAMP_NTZ(4)", expectedPrecision: 4, expectedUnderlyingType: "TIMESTAMP_NTZ"},
		{input: "timestamp without time zone(5)", expectedPrecision: 5, expectedUnderlyingType: "TIMESTAMP WITHOUT TIME ZONE"},
		{input: "TIMESTAMP_NTZ(   2   )", expectedPrecision: 2, expectedUnderlyingType: "TIMESTAMP_NTZ"},
		{input: "    TIMESTAMP_NTZ   (   7   )    ", expectedPrecision: 7, expectedUnderlyingType: "TIMESTAMP_NTZ"},
		{input: fmt.Sprintf("TIMESTAMP_NTZ(%d)", DefaultTimestampPrecision), expectedPrecision: DefaultTimestampPrecision, expectedUnderlyingType: "TIMESTAMP_NTZ"},

		defaults("   TIMESTAMP_NTZ   "),
		defaults("TIMESTAMP_NTZ"),
		defaults("TIMESTAMPNTZ"),
		defaults("TIMESTAMP WITHOUT TIME ZONE"),
		defaults("DATETIME"),
		defaults("timestamp_ntz"),
		defaults("timestampntz"),
		defaults("timestamp without time zone"),
		defaults("datetime"),
	}

	negativeTestCases := []test{
		negative("TIMESTAMP_NTZ(38, 0)"),
		negative("TIMESTAMP_NTZ(38, 2)"),
		negative("TIMESTAMP_NTZ()"),
		negative("T I M E S T A M P _ N T Z"),
		negative("other"),
		negative("other(3)"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &TimestampNtzDataType{}, parsed)

			assert.Equal(t, tc.expectedPrecision, parsed.(*TimestampNtzDataType).precision)
			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*TimestampNtzDataType).underlyingType)

			assert.Equal(t, TimestampNtzLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", parsed.(*TimestampNtzDataType).underlyingType, parsed.(*TimestampNtzDataType).precision), parsed.ToSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", TimestampNtzLegacyDataType, parsed.(*TimestampNtzDataType).precision), parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_TimestampTz(t *testing.T) {
	type test struct {
		input                  string
		expectedPrecision      int
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedPrecision:      DefaultTimestampPrecision,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		{input: "TIMESTAMP_TZ(4)", expectedPrecision: 4, expectedUnderlyingType: "TIMESTAMP_TZ"},
		{input: "timestamp with time zone(5)", expectedPrecision: 5, expectedUnderlyingType: "TIMESTAMP WITH TIME ZONE"},
		{input: "TIMESTAMP_TZ(   2   )", expectedPrecision: 2, expectedUnderlyingType: "TIMESTAMP_TZ"},
		{input: "    TIMESTAMP_TZ   (   7   )    ", expectedPrecision: 7, expectedUnderlyingType: "TIMESTAMP_TZ"},
		{input: fmt.Sprintf("TIMESTAMP_TZ(%d)", DefaultTimestampPrecision), expectedPrecision: DefaultTimestampPrecision, expectedUnderlyingType: "TIMESTAMP_TZ"},

		defaults("   TIMESTAMP_TZ   "),
		defaults("TIMESTAMP_TZ"),
		defaults("TIMESTAMPTZ"),
		defaults("TIMESTAMP WITH TIME ZONE"),
		defaults("timestamp_tz"),
		defaults("timestamptz"),
		defaults("timestamp with time zone"),
	}

	negativeTestCases := []test{
		negative("TIMESTAMP_TZ(38, 0)"),
		negative("TIMESTAMP_TZ(38, 2)"),
		negative("TIMESTAMP_TZ()"),
		negative("T I M E S T A M P _ T Z"),
		negative("other"),
		negative("other(3)"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &TimestampTzDataType{}, parsed)

			assert.Equal(t, tc.expectedPrecision, parsed.(*TimestampTzDataType).precision)
			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*TimestampTzDataType).underlyingType)

			assert.Equal(t, TimestampTzLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", parsed.(*TimestampTzDataType).underlyingType, parsed.(*TimestampTzDataType).precision), parsed.ToSql())
			assert.Equal(t, fmt.Sprintf("%s(%d)", TimestampTzLegacyDataType, parsed.(*TimestampTzDataType).precision), parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Variant(t *testing.T) {
	type test struct {
		input                  string
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		defaults("   VARIANT   "),
		defaults("VARIANT"),
		defaults("variant"),
	}

	negativeTestCases := []test{
		negative("VARIANT(38, 0)"),
		negative("VARIANT(38, 2)"),
		negative("VARIANT(38)"),
		negative("VARIANT()"),
		negative("V A R I A N T"),
		negative("other"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &VariantDataType{}, parsed)

			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*VariantDataType).underlyingType)

			assert.Equal(t, VariantLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, tc.expectedUnderlyingType, parsed.ToSql())
			assert.Equal(t, VariantLegacyDataType, parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Object(t *testing.T) {
	type test struct {
		input                  string
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		defaults("   OBJECT   "),
		defaults("OBJECT"),
		defaults("object"),
	}

	negativeTestCases := []test{
		negative("OBJECT(38, 0)"),
		negative("OBJECT(38, 2)"),
		negative("OBJECT(38)"),
		negative("OBJECT()"),
		negative("O B J E C T"),
		negative("other"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &ObjectDataType{}, parsed)

			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*ObjectDataType).underlyingType)

			assert.Equal(t, ObjectLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, tc.expectedUnderlyingType, parsed.ToSql())
			assert.Equal(t, ObjectLegacyDataType, parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Array(t *testing.T) {
	type test struct {
		input                  string
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		defaults("   ARRAY   "),
		defaults("ARRAY"),
		defaults("array"),
	}

	negativeTestCases := []test{
		negative("ARRAY(38, 0)"),
		negative("ARRAY(38, 2)"),
		negative("ARRAY(38)"),
		negative("ARRAY()"),
		negative("A R R A Y"),
		negative("other"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &ArrayDataType{}, parsed)

			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*ArrayDataType).underlyingType)

			assert.Equal(t, ArrayLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, tc.expectedUnderlyingType, parsed.ToSql())
			assert.Equal(t, ArrayLegacyDataType, parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Geography(t *testing.T) {
	type test struct {
		input                  string
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		defaults("   GEOGRAPHY   "),
		defaults("GEOGRAPHY"),
		defaults("geography"),
	}

	negativeTestCases := []test{
		negative("GEOGRAPHY(38, 0)"),
		negative("GEOGRAPHY(38, 2)"),
		negative("GEOGRAPHY(38)"),
		negative("GEOGRAPHY()"),
		negative("G E O G R A P H Y"),
		negative("other"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &GeographyDataType{}, parsed)

			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*GeographyDataType).underlyingType)

			assert.Equal(t, GeographyLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, tc.expectedUnderlyingType, parsed.ToSql())
			assert.Equal(t, GeographyLegacyDataType, parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Geometry(t *testing.T) {
	type test struct {
		input                  string
		expectedUnderlyingType string
	}
	defaults := func(input string) test {
		return test{
			input:                  input,
			expectedUnderlyingType: strings.TrimSpace(strings.ToUpper(input)),
		}
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		defaults("   GEOMETRY   "),
		defaults("GEOMETRY"),
		defaults("geometry"),
	}

	negativeTestCases := []test{
		negative("GEOMETRY(38, 0)"),
		negative("GEOMETRY(38, 2)"),
		negative("GEOMETRY(38)"),
		negative("GEOMETRY()"),
		negative("G E O M E T R Y"),
		negative("other"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &GeometryDataType{}, parsed)

			assert.Equal(t, tc.expectedUnderlyingType, parsed.(*GeometryDataType).underlyingType)

			assert.Equal(t, GeometryLegacyDataType, parsed.ToLegacyDataTypeSql())
			assert.Equal(t, tc.expectedUnderlyingType, parsed.ToSql())
			assert.Equal(t, GeometryLegacyDataType, parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Vector(t *testing.T) {
	type test struct {
		input             string
		expectedInnerType string
		expectedDimension int
	}
	negative := func(input string) test {
		return test{input: input}
	}

	positiveTestCases := []test{
		{input: "VECTOR(INT, 2)", expectedInnerType: "INT", expectedDimension: 2},
		{input: "VECTOR(FLOAT, 2)", expectedInnerType: "FLOAT", expectedDimension: 2},
		{input: "VeCtOr   ( InT    ,     40     )", expectedInnerType: "INT", expectedDimension: 40},
		{input: "      VECTOR   ( INT    ,     40     )", expectedInnerType: "INT", expectedDimension: 40},
	}

	negativeTestCases := []test{
		negative("VECTOR(1, 2)"),
		negative("VECTOR(1)"),
		negative("VECTOR(2, INT)"),
		negative("VECTOR()"),
		negative("VECTOR"),
		negative("VECTOR(INT, 2, 3)"),
		negative("VECTOR(INT)"),
		negative("VECTOR(x, 2)"),
		negative("VECTOR("),
		negative("VECTOR)"),
		negative("VEC TOR"),
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &VectorDataType{}, parsed)

			assert.Equal(t, tc.expectedInnerType, parsed.(*VectorDataType).innerType)
			assert.Equal(t, tc.expectedDimension, parsed.(*VectorDataType).dimension)
			assert.Equal(t, "VECTOR", parsed.(*VectorDataType).underlyingType)

			assert.Equal(t, fmt.Sprintf("%s(%s, %d)", parsed.(*VectorDataType).underlyingType, parsed.(*VectorDataType).innerType, parsed.(*VectorDataType).dimension), parsed.ToLegacyDataTypeSql())
			assert.Equal(t, fmt.Sprintf("%s(%s, %d)", parsed.(*VectorDataType).underlyingType, parsed.(*VectorDataType).innerType, parsed.(*VectorDataType).dimension), parsed.ToSql())
			assert.Equal(t, fmt.Sprintf("%s(%s, %d)", parsed.(*VectorDataType).underlyingType, parsed.(*VectorDataType).innerType, parsed.(*VectorDataType).dimension), parsed.Canonical())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_ParseDataType_Table(t *testing.T) {
	type column struct {
		Name string
		Type string
	}
	type test struct {
		input           string
		expectedColumns []column
	}

	positiveTestCases := []test{
		{input: "TABLE()", expectedColumns: []column{}},
		{input: "TABLE ()", expectedColumns: []column{}},
		{input: "TABLE ( 	 )", expectedColumns: []column{}},
		{input: "TABLE(arg_name NUMBER)", expectedColumns: []column{{"arg_name", "NUMBER"}}},
		{input: "TABLE(arg_name double precision, arg_name_2 NUMBER)", expectedColumns: []column{{"arg_name", "double precision"}, {"arg_name_2", "NUMBER"}}},
		{input: "TABLE(arg_name NUMBER(38))", expectedColumns: []column{{"arg_name", "NUMBER(38)"}}},
		{input: "TABLE(arg_name NUMBER(38), arg_name_2 VARCHAR)", expectedColumns: []column{{"arg_name", "NUMBER(38)"}, {"arg_name_2", "VARCHAR"}}},
		{input: "TABLE(arg_name number, second float, third GEOGRAPHY)", expectedColumns: []column{{"arg_name", "number"}, {"second", "float"}, {"third", "GEOGRAPHY"}}},
		{input: "TABLE  (		arg_name 		varchar, 		second 	date, third TIME 			)", expectedColumns: []column{{"arg_name", "varchar"}, {"second", "date"}, {"third", "time"}}},
		// TODO [SNOW-2054316]: Support types with parameters (for now, only legacy types are supported because Snowflake returns only with this output), e.g. TABLE(ARG NUMBER(38, 0))
		// TODO [SNOW-2054316]: Support nested tables, e.g. TABLE(ARG NUMBER, NESTED TABLE(A VARCHAR, B GEOMETRY))
		// TODO [SNOW-2054316]: Support complex argument names (with quotes / spaces / special characters / etc)
	}

	negativeTestCases := []test{
		{input: "TABLE())"},
		{input: "TABLE(1, 2)"},
		{input: "TABLE(INT, INT)"},
		{input: "TABLE(a b)"},
		{input: "TABLE(1)"},
		{input: "TABLE(2, INT)"},
		{input: "TABLE"},
		{input: "TABLE(INT, 2, 3)"},
		{input: "TABLE(INT)"},
		{input: "TABLE(x, 2)"},
		{input: "TABLE("},
		{input: "TABLE)"},
		{input: "TA BLE"},
	}

	for _, tc := range positiveTestCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.NoError(t, err)
			require.IsType(t, &TableDataType{}, parsed)

			assert.Equal(t, "TABLE", parsed.(*TableDataType).underlyingType)
			assert.Equal(t, len(tc.expectedColumns), len(parsed.(*TableDataType).columns))
			for i, column := range tc.expectedColumns {
				assert.Equal(t, column.Name, parsed.(*TableDataType).columns[i].name)
				parsedType, err := ParseDataType(column.Type)
				require.NoError(t, err)
				assert.Equal(t, parsedType.ToLegacyDataTypeSql(), parsed.(*TableDataType).columns[i].dataType.ToLegacyDataTypeSql())
			}

			legacyColumns := strings.Join(collections.Map(tc.expectedColumns, func(col column) string {
				parsedType, err := ParseDataType(col.Type)
				require.NoError(t, err)
				return fmt.Sprintf("%s %s", col.Name, parsedType.ToLegacyDataTypeSql())
			}), ", ")
			assert.Equal(t, fmt.Sprintf("TABLE(%s)", legacyColumns), parsed.ToLegacyDataTypeSql())

			canonicalColumns := strings.Join(collections.Map(tc.expectedColumns, func(col column) string {
				parsedType, err := ParseDataType(col.Type)
				require.NoError(t, err)
				return fmt.Sprintf("%s %s", col.Name, parsedType.Canonical())
			}), ", ")
			assert.Equal(t, fmt.Sprintf("TABLE(%s)", canonicalColumns), parsed.Canonical())

			columns := strings.Join(collections.Map(tc.expectedColumns, func(col column) string {
				parsedType, err := ParseDataType(col.Type)
				require.NoError(t, err)
				return fmt.Sprintf("%s %s", col.Name, parsedType.ToSql())
			}), ", ")
			assert.Equal(t, fmt.Sprintf("TABLE(%s)", columns), parsed.ToSql())
		})
	}

	for _, tc := range negativeTestCases {
		tc := tc
		t.Run("negative: "+tc.input, func(t *testing.T) {
			parsed, err := ParseDataType(tc.input)

			require.Error(t, err)
			require.Nil(t, parsed)
		})
	}
}

func Test_AreTheSame(t *testing.T) {
	// empty d1/d2 means nil DataType input
	type test struct {
		d1              string
		d2              string
		expectedOutcome bool
	}

	testCases := []test{
		{d1: "", d2: "", expectedOutcome: true},
		{d1: "", d2: "NUMBER", expectedOutcome: false},
		{d1: "NUMBER", d2: "", expectedOutcome: false},

		{d1: "NUMBER(20)", d2: "NUMBER(20, 2)", expectedOutcome: false},
		{d1: "NUMBER(20, 1)", d2: "NUMBER(20, 2)", expectedOutcome: false},
		{d1: "NUMBER", d2: "NUMBER(20, 2)", expectedOutcome: false},
		{d1: "NUMBER", d2: fmt.Sprintf("NUMBER(%d, %d)", DefaultNumberPrecision, DefaultNumberScale), expectedOutcome: true},
		{d1: fmt.Sprintf("NUMBER(%d)", DefaultNumberPrecision), d2: fmt.Sprintf("NUMBER(%d, %d)", DefaultNumberPrecision, DefaultNumberScale), expectedOutcome: true},
		{d1: "NUMBER", d2: "NUMBER", expectedOutcome: true},
		{d1: "NUMBER(20)", d2: "NUMBER(20)", expectedOutcome: true},
		{d1: "NUMBER(20, 2)", d2: "NUMBER(20, 2)", expectedOutcome: true},
		{d1: "INT", d2: "NUMBER", expectedOutcome: true},
		{d1: "INT", d2: fmt.Sprintf("NUMBER(%d, %d)", DefaultNumberPrecision, DefaultNumberScale), expectedOutcome: true},
		{d1: "INT", d2: "NUMBER(20)", expectedOutcome: false},
		{d1: "NUMBER", d2: "VARCHAR", expectedOutcome: false},
		{d1: "NUMBER(20)", d2: "VARCHAR(20)", expectedOutcome: false},
		{d1: "CHAR", d2: "VARCHAR", expectedOutcome: false},
		{d1: "CHAR", d2: fmt.Sprintf("VARCHAR(%d)", DefaultCharLength), expectedOutcome: true},
		{d1: fmt.Sprintf("CHAR(%d)", DefaultVarcharLength), d2: "VARCHAR", expectedOutcome: true},
		{d1: "BINARY", d2: "BINARY", expectedOutcome: true},
		{d1: "BINARY", d2: "VARBINARY", expectedOutcome: true},
		{d1: "BINARY(20)", d2: "BINARY(20)", expectedOutcome: true},
		{d1: "BINARY(20)", d2: "BINARY(30)", expectedOutcome: false},
		{d1: "BINARY", d2: "BINARY(30)", expectedOutcome: false},
		{d1: fmt.Sprintf("BINARY(%d)", DefaultBinarySize), d2: "BINARY", expectedOutcome: true},
		{d1: "FLOAT", d2: "FLOAT4", expectedOutcome: true},
		{d1: "DOUBLE", d2: "FLOAT8", expectedOutcome: true},
		{d1: "DOUBLE PRECISION", d2: "REAL", expectedOutcome: true},
		{d1: "TIMESTAMPLTZ", d2: "TIMESTAMPNTZ", expectedOutcome: false},
		{d1: "TIMESTAMPLTZ", d2: "TIMESTAMPTZ", expectedOutcome: false},
		{d1: "TIMESTAMPLTZ", d2: fmt.Sprintf("TIMESTAMPLTZ(%d)", DefaultTimestampPrecision), expectedOutcome: true},
		{d1: "VECTOR(INT, 20)", d2: "VECTOR(INT, 20)", expectedOutcome: true},
		{d1: "VECTOR(INT, 20)", d2: "VECTOR(INT, 30)", expectedOutcome: false},
		{d1: "VECTOR(FLOAT, 20)", d2: "VECTOR(INT, 30)", expectedOutcome: false},
		{d1: "VECTOR(FLOAT, 20)", d2: "VECTOR(INT, 20)", expectedOutcome: false},
		{d1: "VECTOR(FLOAT, 20)", d2: "VECTOR(FLOAT, 20)", expectedOutcome: true},
		{d1: "VECTOR(FLOAT, 20)", d2: "FLOAT", expectedOutcome: false},
		{d1: "TIME", d2: "TIME", expectedOutcome: true},
		{d1: "TIME", d2: "TIME(5)", expectedOutcome: false},
		{d1: "TIME", d2: fmt.Sprintf("TIME(%d)", DefaultTimePrecision), expectedOutcome: true},
		{d1: "TABLE()", d2: "TABLE()", expectedOutcome: true},
		{d1: "TABLE(A NUMBER)", d2: "TABLE(B NUMBER)", expectedOutcome: false},
		{d1: "TABLE(A NUMBER)", d2: "TABLE(a NUMBER)", expectedOutcome: false},
		{d1: "TABLE(A NUMBER)", d2: "TABLE(A VARCHAR)", expectedOutcome: false},
		{d1: "TABLE(A NUMBER, B VARCHAR)", d2: "TABLE(A NUMBER, B VARCHAR)", expectedOutcome: true},
		{d1: "TABLE(A NUMBER, B NUMBER)", d2: "TABLE(A NUMBER, B VARCHAR)", expectedOutcome: false},
		{d1: "TABLE()", d2: "TABLE(A NUMBER)", expectedOutcome: false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf(`compare "%s" with "%s" expecting %t`, tc.d1, tc.d2, tc.expectedOutcome), func(t *testing.T) {
			var p1, p2 DataType
			var err error

			if tc.d1 != "" {
				p1, err = ParseDataType(tc.d1)
				require.NoError(t, err)
			}

			if tc.d2 != "" {
				p2, err = ParseDataType(tc.d2)
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedOutcome, AreTheSame(p1, p2))
		})
	}
}

func Test_AreDefinitelyDifferent(t *testing.T) {
	// empty d1/d2 means nil DataType input
	type test struct {
		d1              string
		d2              string
		expectedOutcome bool
	}

	testCases := []test{
		{d1: "", d2: "", expectedOutcome: false},
		{d1: "", d2: "NUMBER", expectedOutcome: true},
		{d1: "NUMBER", d2: "", expectedOutcome: true},

		{d1: "NUMBER(20)", d2: "NUMBER(20, 2)", expectedOutcome: false},
		{d1: "NUMBER(20, 1)", d2: "NUMBER(20, 2)", expectedOutcome: true},
		{d1: "NUMBER", d2: "NUMBER(20, 2)", expectedOutcome: false},
		{d1: "NUMBER", d2: fmt.Sprintf("NUMBER(%d, %d)", DefaultNumberPrecision, DefaultNumberScale), expectedOutcome: false},
		{d1: fmt.Sprintf("NUMBER(%d)", DefaultNumberPrecision), d2: fmt.Sprintf("NUMBER(%d, %d)", DefaultNumberPrecision, DefaultNumberScale), expectedOutcome: false},
		{d1: "NUMBER", d2: "NUMBER", expectedOutcome: false},
		{d1: "NUMBER(20)", d2: "NUMBER(20)", expectedOutcome: false},
		{d1: "NUMBER(20, 2)", d2: "NUMBER(20, 2)", expectedOutcome: false},
		{d1: "INT", d2: "NUMBER", expectedOutcome: false},
		{d1: "INT", d2: fmt.Sprintf("NUMBER(%d, %d)", DefaultNumberPrecision, DefaultNumberScale), expectedOutcome: false},
		{d1: "INT", d2: "NUMBER(20)", expectedOutcome: true},
		{d1: "INT", d2: "NUMBER(20, 4)", expectedOutcome: true},
		{d1: "INT", d2: fmt.Sprintf("NUMBER(%d, %d)", DefaultNumberPrecision, 5), expectedOutcome: true},
		{d1: "NUMBER", d2: "VARCHAR", expectedOutcome: true},
		{d1: "NUMBER(20)", d2: "VARCHAR(20)", expectedOutcome: true},
		{d1: "CHAR", d2: "VARCHAR", expectedOutcome: false},
		{d1: "CHAR", d2: fmt.Sprintf("VARCHAR(%d)", DefaultCharLength), expectedOutcome: false},
		{d1: fmt.Sprintf("CHAR(%d)", DefaultVarcharLength), d2: "VARCHAR", expectedOutcome: false},
		{d1: "BINARY", d2: "BINARY", expectedOutcome: false},
		{d1: "BINARY", d2: "VARBINARY", expectedOutcome: false},
		{d1: "BINARY(20)", d2: "BINARY(20)", expectedOutcome: false},
		{d1: "BINARY(20)", d2: "BINARY(30)", expectedOutcome: true},
		{d1: "BINARY", d2: "BINARY(30)", expectedOutcome: false},
		{d1: fmt.Sprintf("BINARY(%d)", DefaultBinarySize), d2: "BINARY", expectedOutcome: false},
		{d1: "FLOAT", d2: "FLOAT4", expectedOutcome: false},
		{d1: "DOUBLE", d2: "FLOAT8", expectedOutcome: false},
		{d1: "DOUBLE PRECISION", d2: "REAL", expectedOutcome: false},
		{d1: "TIMESTAMPLTZ", d2: "TIMESTAMPNTZ", expectedOutcome: true},
		{d1: "TIMESTAMPLTZ", d2: "TIMESTAMPTZ", expectedOutcome: true},
		{d1: "TIMESTAMPLTZ", d2: fmt.Sprintf("TIMESTAMPLTZ(%d)", DefaultTimestampPrecision), expectedOutcome: false},
		{d1: "VECTOR(INT, 20)", d2: "VECTOR(INT, 20)", expectedOutcome: false},
		{d1: "VECTOR(INT, 20)", d2: "VECTOR(INT, 30)", expectedOutcome: true},
		{d1: "VECTOR(FLOAT, 20)", d2: "VECTOR(INT, 30)", expectedOutcome: true},
		{d1: "VECTOR(FLOAT, 20)", d2: "VECTOR(INT, 20)", expectedOutcome: true},
		{d1: "VECTOR(FLOAT, 20)", d2: "VECTOR(FLOAT, 20)", expectedOutcome: false},
		{d1: "VECTOR(FLOAT, 20)", d2: "FLOAT", expectedOutcome: true},
		{d1: "TIME", d2: "TIME", expectedOutcome: false},
		{d1: "TIME", d2: "TIME(5)", expectedOutcome: false},
		{d1: "TIME", d2: fmt.Sprintf("TIME(%d)", DefaultTimePrecision), expectedOutcome: false},
		{d1: "TABLE()", d2: "TABLE()", expectedOutcome: false},
		{d1: "TABLE(A NUMBER)", d2: "TABLE(B NUMBER)", expectedOutcome: true},
		{d1: "TABLE(A NUMBER)", d2: "TABLE(a NUMBER)", expectedOutcome: true},
		{d1: "TABLE(A NUMBER)", d2: "TABLE(A VARCHAR)", expectedOutcome: true},
		{d1: "TABLE(A NUMBER, B VARCHAR)", d2: "TABLE(A NUMBER, B VARCHAR)", expectedOutcome: false},
		{d1: "TABLE(A NUMBER, B NUMBER)", d2: "TABLE(A NUMBER, B VARCHAR)", expectedOutcome: true},
		{d1: "TABLE()", d2: "TABLE(A NUMBER)", expectedOutcome: true},
		{d1: "TABLE(B CHAR)", d2: "TABLE()", expectedOutcome: true},
		{d1: "TABLE(A NUMBER)", d2: "TABLE(A NUMBER, B VARCHAR)", expectedOutcome: true},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf(`check if "%s" is definitely different than "%s", expecting %t`, tc.d1, tc.d2, tc.expectedOutcome), func(t *testing.T) {
			var p1, p2 DataType
			var err error

			if tc.d1 != "" {
				p1, err = ParseDataType(tc.d1)
				require.NoError(t, err)
			}

			if tc.d2 != "" {
				p2, err = ParseDataType(tc.d2)
				require.NoError(t, err)
			}

			require.Equal(t, tc.expectedOutcome, AreDefinitelyDifferent(p1, p2))
		})
	}
}

func Test_ToSqlWithoutUnknowns(t *testing.T) {
	type test struct {
		dt string
	}

	testCases := []test{
		{dt: "NUMBER"},
		{dt: "NUMBER(20)"},
		{dt: "NUMBER(20, 4)"},
		{dt: fmt.Sprintf("NUMBER(%d, %d)", DefaultNumberPrecision, DefaultNumberScale)},
		{dt: "INT"},
		{dt: "VARCHAR"},
		{dt: "VARCHAR(15)"},
		{dt: "CHAR"},
		{dt: "CHAR(5)"},
		{dt: fmt.Sprintf("VARCHAR(%d)", DefaultCharLength)},
		{dt: fmt.Sprintf("CHAR(%d)", DefaultVarcharLength)},
		{dt: "BINARY"},
		{dt: "VARBINARY"},
		{dt: "BINARY(20)"},
		{dt: "VARBINARY(30)"},
		{dt: fmt.Sprintf("BINARY(%d)", DefaultBinarySize)},
		{dt: "FLOAT"},
		{dt: "FLOAT4"},
		{dt: "DOUBLE PRECISION"},
		{dt: "TIME"},
		{dt: fmt.Sprintf("TIME(%d)", DefaultTimePrecision)},
		{dt: "TIMESTAMPLTZ"},
		{dt: fmt.Sprintf("TIMESTAMPLTZ(%d)", DefaultTimestampPrecision)},
		{dt: "TIMESTAMPNTZ"},
		{dt: fmt.Sprintf("TIMESTAMPNTZ(%d)", DefaultTimestampPrecision)},
		{dt: "TIMESTAMPTZ"},
		{dt: fmt.Sprintf("TIMESTAMPTZ(%d)", DefaultTimestampPrecision)},
		{dt: "TIMESTAMP WITHOUT TIME ZONE"},
		{dt: "VECTOR(INT, 20)"},
		{dt: "VECTOR(FLOAT, 30)"},
		{dt: "TABLE()"},
		{dt: "TABLE(A NUMBER)"},
		{dt: "TABLE(A NUMBER, B INT, C VARCHAR)"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf(`check if ToSqlWithoutUnknowns method works correctly for "%s"`, tc.dt), func(t *testing.T) {
			parsed, err := ParseDataType(tc.dt)
			require.NoError(t, err)
			require.Equal(t, tc.dt, parsed.ToSqlWithoutUnknowns())
		})
	}
}
