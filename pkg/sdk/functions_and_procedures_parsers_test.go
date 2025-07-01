package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseFunctionAndProcedureArguments(t *testing.T) {
	dtOnly := func(dt DataType) ParsedArgument {
		return ParsedArgument{
			ArgType: string(dt),
		}
	}

	dtDefault := func(dt DataType) ParsedArgument {
		return ParsedArgument{
			ArgType:   string(dt),
			IsDefault: true,
		}
	}

	dtName := func(dt DataType, name string) ParsedArgument {
		return ParsedArgument{
			ArgType: string(dt),
			ArgName: name,
		}
	}

	full := func(dt DataType, name string) ParsedArgument {
		return ParsedArgument{
			ArgType:   string(dt),
			ArgName:   name,
			IsDefault: true,
		}
	}

	testCases := []struct {
		Arguments string
		Expected  []ParsedArgument
		Error     string
	}{
		// empty
		{Arguments: ``, Expected: []ParsedArgument{}},

		// basic
		{Arguments: `FLOAT`, Expected: []ParsedArgument{dtOnly(DataTypeFloat)}},
		{Arguments: `FLOAT, NUMBER, TIME`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly(DataTypeNumber), dtOnly(DataTypeTime)}},
		{Arguments: `FLOAT, NUMBER, VECTOR(FLOAT, 20)`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly(DataTypeNumber), dtOnly("VECTOR(FLOAT, 20)")}},
		{Arguments: `VECTOR(FLOAT, 10), NUMBER, VECTOR(FLOAT, 20)`, Expected: []ParsedArgument{dtOnly("VECTOR(FLOAT, 10)"), dtOnly(DataTypeNumber), dtOnly("VECTOR(FLOAT, 20)")}},
		{Arguments: `FLOAT, VARCHAR(200), TIME`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly("VARCHAR(200)"), dtOnly(DataTypeTime)}},
		{Arguments: `FLOAT, VARCHAR(200)`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly("VARCHAR(200)")}},
		{Arguments: `VARCHAR(200), FLOAT`, Expected: []ParsedArgument{dtOnly("VARCHAR(200)"), dtOnly(DataTypeFloat)}},
		{Arguments: `FLOAT, NUMBER(10, 2), TIME`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly("NUMBER(10, 2)"), dtOnly(DataTypeTime)}},
		{Arguments: `FLOAT, NUMBER(10, 2)`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly("NUMBER(10, 2)")}},
		{Arguments: `NUMBER(10, 2), FLOAT`, Expected: []ParsedArgument{dtOnly("NUMBER(10, 2)"), dtOnly(DataTypeFloat)}},

		// with defaults
		{Arguments: `DEFAULT FLOAT, DEFAULT NUMBER, DEFAULT TIME`, Expected: []ParsedArgument{dtDefault(DataTypeFloat), dtDefault(DataTypeNumber), dtDefault(DataTypeTime)}},
		{Arguments: `DEFAULT FLOAT, NUMBER, DEFAULT TIME`, Expected: []ParsedArgument{dtDefault(DataTypeFloat), dtOnly(DataTypeNumber), dtDefault(DataTypeTime)}},

		// with names
		{Arguments: `a FLOAT`, Expected: []ParsedArgument{dtName(DataTypeFloat, "a")}},
		{Arguments: `a FLOAT, B NUMBER, c TIME`, Expected: []ParsedArgument{dtName(DataTypeFloat, "a"), dtName(DataTypeNumber, "B"), dtName(DataTypeTime, "c")}},
		{Arguments: `a FLOAT, NUMBER, c TIME`, Expected: []ParsedArgument{dtName(DataTypeFloat, "a"), dtOnly(DataTypeNumber), dtName(DataTypeTime, "c")}},
		{Arguments: `ab NUMBER(10, 2), x FLOAT, FLOAT`, Expected: []ParsedArgument{dtName("NUMBER(10, 2)", "ab"), dtName(DataTypeFloat, "x"), dtOnly(DataTypeFloat)}},

		// combos
		{Arguments: `DEFAULT ab NUMBER(10, 2), x FLOAT, DEFAULT FLOAT`, Expected: []ParsedArgument{full("NUMBER(10, 2)", "ab"), dtName(DataTypeFloat, "x"), dtDefault(DataTypeFloat)}},
		{Arguments: `DEFAULT NUMBER(10), DEFAULT x FLOAT, aBc FLOAT`, Expected: []ParsedArgument{dtDefault("NUMBER(10)"), full(DataTypeFloat, "x"), dtName(DataTypeFloat, "aBc")}},

		// various spaces
		{Arguments: `  `, Expected: []ParsedArgument{}},
		{Arguments: `FLOAT `, Expected: []ParsedArgument{dtOnly(DataTypeFloat)}},
		{Arguments: `  FLOAT,    NUMBER   , TIME`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly(DataTypeNumber), dtOnly(DataTypeTime)}},
		{Arguments: `VARCHAR(  200), FLOAT`, Expected: []ParsedArgument{dtOnly("VARCHAR(  200)"), dtOnly(DataTypeFloat)}},
		{Arguments: `FLOAT, NUMBER(   10, 2  )   , TIME`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly("NUMBER(   10, 2  )"), dtOnly(DataTypeTime)}},
		{Arguments: `DEFAULT      FLOAT,    DEFAULT NUMBER    , DEFAULT TIME    `, Expected: []ParsedArgument{dtDefault(DataTypeFloat), dtDefault(DataTypeNumber), dtDefault(DataTypeTime)}},
		{Arguments: `    a     FLOAT`, Expected: []ParsedArgument{dtName(DataTypeFloat, "a")}},
		{Arguments: `   ab NUMBER    (   10, 2),x FLOAT, FLOAT  `, Expected: []ParsedArgument{dtName("NUMBER    (   10, 2)", "ab"), dtName(DataTypeFloat, "x"), dtOnly(DataTypeFloat)}},
		{Arguments: `DEFAULT     ab NUMBER(10, 2), x FLOAT, DEFAULT  FLOAT`, Expected: []ParsedArgument{full("NUMBER(10, 2)", "ab"), dtName(DataTypeFloat, "x"), dtDefault(DataTypeFloat)}},

		// incorrect vectors input but we expect it to correctly tokenize
		{Arguments: `FLOAT, NUMBER, VECTOR(VARCHAR, 20)`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly(DataTypeNumber), dtOnly("VECTOR(VARCHAR, 20)")}},
		{Arguments: `FLOAT, NUMBER, VECTOR(INT, INT)`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly(DataTypeNumber), dtOnly("VECTOR(INT, INT)")}},
		{Arguments: `FLOAT, NUMBER, VECTOR(20, FLOAT)`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly(DataTypeNumber), dtOnly("VECTOR(20, FLOAT)")}},

		// incorrect input, but we expect it to correctly tokenize
		{Arguments: `a`, Expected: []ParsedArgument{dtOnly("a")}},
		{Arguments: `a, INT, b`, Expected: []ParsedArgument{dtOnly("a"), dtOnly(DataTypeInt), dtOnly("b")}},
		{Arguments: `VECTOR(FLOAT, 10)| NUMBER, VECTOR(FLOAT, 20)`, Expected: []ParsedArgument{dtOnly("VECTOR(FLOAT, 10)| NUMBER"), dtOnly("VECTOR(FLOAT, 20)")}},
		{Arguments: `FLOAT, NUMBER, VECTORFLOAT, 20`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly(DataTypeNumber), dtOnly("VECTORFLOAT"), dtOnly("20")}},
		{Arguments: `FLOAT, NUMBER, VECTORFLOAT, 20, VECTOR(INT, 10)`, Expected: []ParsedArgument{dtOnly(DataTypeFloat), dtOnly(DataTypeNumber), dtOnly("VECTORFLOAT"), dtOnly("20"), dtOnly("VECTOR(INT, 10)")}},

		// incorrect input - parentheses
		{Arguments: `FLOAT, NUMBER, VECTOR(FLOAT, 20, VECTOR(INT, 10)`, Error: `opening and closing parentheses do not match`},
		{Arguments: `FLOAT, NUMBER, VECTOR(FLOAT, 20`, Error: `opening and closing parentheses do not match`},
		{Arguments: `FLOAT, NUMBER, VECTORFLOAT, 20), VECTOR(INT, 10)`, Error: `opening and closing parentheses do not match`},

		// incorrect input - ending with comma
		{Arguments: `,`, Error: "can't end arguments list with a comma"},
		{Arguments: `FLOAT,`, Error: "can't end arguments list with a comma"},
		{Arguments: `FLOAT, NUMBER,`, Error: "can't end arguments list with a comma"},
		{Arguments: `FLOAT, NUMBER,,`, Error: "can't end arguments list with a comma"},
	}
	for _, testCase := range testCases {
		argumentsWithParentheses := "(" + testCase.Arguments + ")"
		argumentsWithParenthesesAndSpaces := "  (" + testCase.Arguments + " )   "

		body := func(arguments string) {
			dataTypes, err := ParseFunctionAndProcedureArguments(arguments)
			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.Expected, dataTypes)
			}
		}

		t.Run(fmt.Sprintf("parsing function and procedure arguments: `%s`", testCase.Arguments), func(t *testing.T) {
			body(testCase.Arguments)
		})

		t.Run(fmt.Sprintf("parsing function and procedure arguments, wrapped in paretheses: `%s`", argumentsWithParentheses), func(t *testing.T) {
			body(argumentsWithParentheses)
		})

		t.Run(fmt.Sprintf("parsing function and procedure arguments, wrapped in paretheses and with additional spacing: `%s`", argumentsWithParenthesesAndSpaces), func(t *testing.T) {
			body(argumentsWithParenthesesAndSpaces)
		})
	}
}
