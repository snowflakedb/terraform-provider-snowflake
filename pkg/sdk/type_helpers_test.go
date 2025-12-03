package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToParameterType(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected ParameterType
		Error    string
	}{
		{Input: string(ParameterTypeSnowflakeDefault), Expected: ParameterTypeSnowflakeDefault},
		{Input: string(ParameterTypeAccount), Expected: ParameterTypeAccount},
		{Input: string(ParameterTypeUser), Expected: ParameterTypeUser},
		{Input: string(ParameterTypeSession), Expected: ParameterTypeSession},
		{Input: string(ParameterTypeObject), Expected: ParameterTypeObject},
		{Input: string(ParameterTypeWarehouse), Expected: ParameterTypeWarehouse},
		{Input: string(ParameterTypeDatabase), Expected: ParameterTypeDatabase},
		{Input: string(ParameterTypeSchema), Expected: ParameterTypeSchema},
		{Input: string(ParameterTypeTask), Expected: ParameterTypeTask},
		{Input: string(ParameterTypeFunction), Expected: ParameterTypeFunction},
		{Input: string(ParameterTypeProcedure), Expected: ParameterTypeProcedure},
		{Name: "validation: incorrect parameter type", Input: "incorrect", Error: "invalid parameter type: incorrect"},
		{Name: "validation: lower case input", Input: "account", Expected: ParameterType("account")},
	}

	for _, testCase := range testCases {
		name := testCase.Name
		if name == "" {
			name = fmt.Sprintf("%v parameter type", testCase.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToParameterType(testCase.Input)
			if testCase.Error != "" {
				assert.Empty(t, value)
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.Expected, value)
			}
		})
	}
}
