package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseFunctionAndProcedureArguments(t *testing.T) {
	dtOnly := func(dt string) ParsedArgument {
		return ParsedArgument{
			ArgType: dt,
		}
	}

	dtDefault := func(dt string) ParsedArgument {
		return ParsedArgument{
			ArgType:   dt,
			IsDefault: true,
		}
	}

	dtName := func(dt string, name string) ParsedArgument {
		return ParsedArgument{
			ArgType: dt,
			ArgName: name,
		}
	}

	full := func(dt string, name string) ParsedArgument {
		return ParsedArgument{
			ArgType:   dt,
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
		{Arguments: `FLOAT`, Expected: []ParsedArgument{dtOnly("FLOAT")}},
		{Arguments: `FLOAT, NUMBER, TIME`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER"), dtOnly("TIME")}},
		{Arguments: `FLOAT, NUMBER, VECTOR(FLOAT, 20)`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER"), dtOnly("VECTOR(FLOAT, 20)")}},
		{Arguments: `VECTOR(FLOAT, 10), NUMBER, VECTOR(FLOAT, 20)`, Expected: []ParsedArgument{dtOnly("VECTOR(FLOAT, 10)"), dtOnly("NUMBER"), dtOnly("VECTOR(FLOAT, 20)")}},
		{Arguments: `FLOAT, VARCHAR(200), TIME`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("VARCHAR(200)"), dtOnly("TIME")}},
		{Arguments: `FLOAT, VARCHAR(200)`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("VARCHAR(200)")}},
		{Arguments: `VARCHAR(200), FLOAT`, Expected: []ParsedArgument{dtOnly("VARCHAR(200)"), dtOnly("FLOAT")}},
		{Arguments: `FLOAT, NUMBER(10, 2), TIME`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER(10, 2)"), dtOnly("TIME")}},
		{Arguments: `FLOAT, NUMBER(10, 2)`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER(10, 2)")}},
		{Arguments: `NUMBER(10, 2), FLOAT`, Expected: []ParsedArgument{dtOnly("NUMBER(10, 2)"), dtOnly("FLOAT")}},
		{Arguments: `FLOAT, NUMBER, VECTOR()`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER"), dtOnly("VECTOR()")}},
		{Arguments: `NUMBER, VECTOR)2(, FLOAT`, Expected: []ParsedArgument{dtOnly("NUMBER"), dtOnly("VECTOR)2("), dtOnly("FLOAT")}},

		// with defaults
		{Arguments: `DEFAULT FLOAT, DEFAULT NUMBER, DEFAULT TIME`, Expected: []ParsedArgument{dtDefault("FLOAT"), dtDefault("NUMBER"), dtDefault("TIME")}},
		{Arguments: `DEFAULT FLOAT, NUMBER, DEFAULT TIME`, Expected: []ParsedArgument{dtDefault("FLOAT"), dtOnly("NUMBER"), dtDefault("TIME")}},

		// with names
		{Arguments: `a FLOAT`, Expected: []ParsedArgument{dtName("FLOAT", "a")}},
		{Arguments: `a FLOAT, B NUMBER, c TIME`, Expected: []ParsedArgument{dtName("FLOAT", "a"), dtName("NUMBER", "B"), dtName("TIME", "c")}},
		{Arguments: `a FLOAT, NUMBER, c TIME`, Expected: []ParsedArgument{dtName("FLOAT", "a"), dtOnly("NUMBER"), dtName("TIME", "c")}},
		{Arguments: `ab NUMBER(10, 2), x FLOAT, FLOAT`, Expected: []ParsedArgument{dtName("NUMBER(10, 2)", "ab"), dtName("FLOAT", "x"), dtOnly("FLOAT")}},

		// combos
		{Arguments: `DEFAULT ab NUMBER(10, 2), x FLOAT, DEFAULT FLOAT`, Expected: []ParsedArgument{full("NUMBER(10, 2)", "ab"), dtName("FLOAT", "x"), dtDefault("FLOAT")}},
		{Arguments: `DEFAULT NUMBER(10), DEFAULT x FLOAT, aBc FLOAT`, Expected: []ParsedArgument{dtDefault("NUMBER(10)"), full("FLOAT", "x"), dtName("FLOAT", "aBc")}},

		// various spaces
		{Arguments: `  `, Expected: []ParsedArgument{}},
		{Arguments: `FLOAT `, Expected: []ParsedArgument{dtOnly("FLOAT")}},
		{Arguments: `  FLOAT,    NUMBER   , TIME`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER"), dtOnly("TIME")}},
		{Arguments: `VARCHAR(  200), FLOAT`, Expected: []ParsedArgument{dtOnly("VARCHAR(  200)"), dtOnly("FLOAT")}},
		{Arguments: `FLOAT, NUMBER(   10, 2  )   , TIME`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER(   10, 2  )"), dtOnly("TIME")}},
		{Arguments: `DEFAULT      FLOAT,    DEFAULT NUMBER    , DEFAULT TIME    `, Expected: []ParsedArgument{dtDefault("FLOAT"), dtDefault("NUMBER"), dtDefault("TIME")}},
		{Arguments: `    a     FLOAT`, Expected: []ParsedArgument{dtName("FLOAT", "a")}},
		{Arguments: `   ab NUMBER    (   10, 2),x FLOAT, FLOAT  `, Expected: []ParsedArgument{dtName("NUMBER    (   10, 2)", "ab"), dtName("FLOAT", "x"), dtOnly("FLOAT")}},
		{Arguments: `DEFAULT     ab NUMBER(10, 2), x FLOAT, DEFAULT  FLOAT`, Expected: []ParsedArgument{full("NUMBER(10, 2)", "ab"), dtName("FLOAT", "x"), dtDefault("FLOAT")}},

		// incorrect vectors input but we expect it to correctly tokenize
		{Arguments: `FLOAT, NUMBER, VECTOR(VARCHAR, 20)`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER"), dtOnly("VECTOR(VARCHAR, 20)")}},
		{Arguments: `FLOAT, NUMBER, VECTOR(INT, INT)`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER"), dtOnly("VECTOR(INT, INT)")}},
		{Arguments: `FLOAT, NUMBER, VECTOR(20, FLOAT)`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER"), dtOnly("VECTOR(20, FLOAT)")}},

		// incorrect input, but we expect it to correctly tokenize
		{Arguments: `a`, Expected: []ParsedArgument{dtOnly("a")}},
		{Arguments: `a, INT, b`, Expected: []ParsedArgument{dtOnly("a"), dtOnly("INT"), dtOnly("b")}},
		{Arguments: `VECTOR(FLOAT, 10)| NUMBER, VECTOR(FLOAT, 20)`, Expected: []ParsedArgument{dtOnly("VECTOR(FLOAT, 10)| NUMBER"), dtOnly("VECTOR(FLOAT, 20)")}},
		{Arguments: `FLOAT, NUMBER, VECTORFLOAT, 20`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER"), dtOnly("VECTORFLOAT"), dtOnly("20")}},
		{Arguments: `FLOAT, NUMBER, VECTORFLOAT, 20, VECTOR(INT, 10)`, Expected: []ParsedArgument{dtOnly("FLOAT"), dtOnly("NUMBER"), dtOnly("VECTORFLOAT"), dtOnly("20"), dtOnly("VECTOR(INT, 10)")}},

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

		t.Run(fmt.Sprintf("parsing function and procedure arguments, wrapped in parentheses: `%s`", argumentsWithParentheses), func(t *testing.T) {
			body(argumentsWithParentheses)
		})

		t.Run(fmt.Sprintf("parsing function and procedure arguments, wrapped in parentheses and with additional spacing: `%s`", argumentsWithParenthesesAndSpaces), func(t *testing.T) {
			body(argumentsWithParenthesesAndSpaces)
		})
	}
}
