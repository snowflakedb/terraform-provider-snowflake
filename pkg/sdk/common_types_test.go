package sdk

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ToStringProperty(t *testing.T) {
	t.Run("with empty property row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "null",
			DefaultValue: "",
			Description:  "desc",
		}
		prop := row.toStringProperty()
		assert.Empty(t, prop.Value)
		assert.Empty(t, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with property row containing values", func(t *testing.T) {
		row := &propertyRow{
			Value:        "value",
			DefaultValue: "default value",
			Description:  "desc",
		}
		prop := row.toStringProperty()
		assert.Equal(t, "value", prop.Value)
		assert.Equal(t, "default value", prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})
}

func Test_ToIntProperty(t *testing.T) {
	t.Run("with empty property row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "null",
			DefaultValue: "",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Nil(t, prop.Value)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with property row not containing numbers", func(t *testing.T) {
		row := &propertyRow{
			Value:        "value",
			DefaultValue: "default value",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Nil(t, prop.Value)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with property not containing default value", func(t *testing.T) {
		row := &propertyRow{
			Value:        "10",
			DefaultValue: "null",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Equal(t, 10, *prop.Value)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with property row containing numbers", func(t *testing.T) {
		row := &propertyRow{
			Value:        "10",
			DefaultValue: "0",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		assert.Equal(t, 10, *prop.Value)
		assert.Equal(t, 0, *prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with negative value row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "-1",
			DefaultValue: "0",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		require.NotNil(t, prop.Value)
		assert.Equal(t, -1, *prop.Value)
	})

	t.Run("with decimal part value - not parsed correctly", func(t *testing.T) {
		row := &propertyRow{
			Value:        "0.85",
			DefaultValue: "0",
			Description:  "desc",
		}
		prop := row.toIntProperty()
		require.Nil(t, prop.Value)
	})
}

func Test_ToBoolProperty(t *testing.T) {
	t.Run("with empty property row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "null",
			DefaultValue: "",
			Description:  "desc",
		}
		prop := row.toBoolProperty()
		assert.False(t, prop.Value)
		assert.False(t, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})

	t.Run("with property row containing values", func(t *testing.T) {
		row := &propertyRow{
			Value:        "true",
			DefaultValue: "false",
			Description:  "desc",
		}
		prop := row.toBoolProperty()
		assert.True(t, prop.Value)
		assert.False(t, prop.DefaultValue)
		assert.Equal(t, row.Description, prop.Description)
	})
}

func Test_ToFloatProperty(t *testing.T) {
	t.Run("with empty property row", func(t *testing.T) {
		row := &propertyRow{
			Value:        "null",
			DefaultValue: "null",
			Description:  "desc",
		}
		prop := row.toFloatProperty()
		assert.Nil(t, prop.Value)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with property row not containing numbers", func(t *testing.T) {
		row := &propertyRow{
			Value:        "value",
			DefaultValue: "default value",
			Description:  "desc",
		}
		prop := row.toFloatProperty()
		assert.Nil(t, prop.Value)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with property not containing default value", func(t *testing.T) {
		row := &propertyRow{
			Value:        "10.5",
			DefaultValue: "null",
			Description:  "desc",
		}
		prop := row.toFloatProperty()
		assert.InDelta(t, 10.5, *prop.Value, testvars.FloatEpsilon)
		assert.Nil(t, prop.DefaultValue)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with property row containing numbers", func(t *testing.T) {
		row := &propertyRow{
			Value:        "10.1",
			DefaultValue: "10.5",
			Description:  "desc",
		}
		prop := row.toFloatProperty()
		assert.InDelta(t, 10.1, *prop.Value, testvars.FloatEpsilon)
		assert.InDelta(t, 10.5, *prop.DefaultValue, testvars.FloatEpsilon)
		assert.Equal(t, prop.Description, row.Description)
	})

	t.Run("with negative value row and zero", func(t *testing.T) {
		row := &propertyRow{
			Value:        "-1.0",
			DefaultValue: "0",
			Description:  "desc",
		}
		prop := row.toFloatProperty()
		assert.InDelta(t, float64(-1), *prop.Value, testvars.FloatEpsilon)
		assert.InDelta(t, float64(0), *prop.DefaultValue, testvars.FloatEpsilon)
		assert.Equal(t, prop.Description, row.Description)
	})
}

func TestToStorageSerializationPolicy(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected StorageSerializationPolicy
		Error    string
	}{
		{Input: string(StorageSerializationPolicyOptimized), Expected: StorageSerializationPolicyOptimized},
		{Input: string(StorageSerializationPolicyCompatible), Expected: StorageSerializationPolicyCompatible},
		{Name: "validation: incorrect storage serialization policy", Input: "incorrect", Error: "unknown storage serialization policy: incorrect"},
		{Name: "validation: empty input", Input: "", Error: "unknown storage serialization policy: "},
		{Name: "validation: lower case input", Input: "optimized", Expected: StorageSerializationPolicyOptimized},
	}

	for _, testCase := range testCases {
		name := testCase.Name
		if name == "" {
			name = fmt.Sprintf("%v storage serialization policy", testCase.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToStorageSerializationPolicy(testCase.Input)
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

func TestToLogLevel(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected LogLevel
		Error    string
	}{
		{Input: string(LogLevelTrace), Expected: LogLevelTrace},
		{Input: string(LogLevelDebug), Expected: LogLevelDebug},
		{Input: string(LogLevelInfo), Expected: LogLevelInfo},
		{Input: string(LogLevelWarn), Expected: LogLevelWarn},
		{Input: string(LogLevelError), Expected: LogLevelError},
		{Input: string(LogLevelFatal), Expected: LogLevelFatal},
		{Input: string(LogLevelOff), Expected: LogLevelOff},
		{Name: "validation: incorrect log level", Input: "incorrect", Error: "unknown log level: incorrect"},
		{Name: "validation: empty input", Input: "", Error: "unknown log level: "},
		{Name: "validation: lower case input", Input: "info", Expected: LogLevelInfo},
	}

	for _, testCase := range testCases {
		name := testCase.Name
		if name == "" {
			name = fmt.Sprintf("%v log level", testCase.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToLogLevel(testCase.Input)
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

func Test_ToExecuteAs(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected ExecuteAs
		Error    string
	}{
		{Input: string(ExecuteAsCaller), Expected: ExecuteAsCaller},
		{Input: string(ExecuteAsOwner), Expected: ExecuteAsOwner},
		{Name: "validation: incorrect execute as", Input: "incorrect", Error: "unknown execute as: incorrect"},
		{Name: "validation: empty input", Input: "", Error: "unknown execute as: "},
		{Name: "validation: lower case input", Input: "caller", Expected: ExecuteAsCaller},
	}

	for _, testCase := range testCases {
		name := testCase.Name
		if name == "" {
			name = fmt.Sprintf("%v execute as", testCase.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToExecuteAs(testCase.Input)
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

func Test_ToNullInputBehavior(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected NullInputBehavior
		Error    string
	}{
		{Input: string(NullInputBehaviorCalledOnNullInput), Expected: NullInputBehaviorCalledOnNullInput},
		{Input: string(NullInputBehaviorReturnsNullInput), Expected: NullInputBehaviorReturnsNullInput},
		{Input: string(NullInputBehaviorStrict), Expected: NullInputBehaviorReturnsNullInput},
		{Name: "validation: incorrect null input behavior", Input: "incorrect", Error: "unknown null input behavior: incorrect"},
		{Name: "validation: empty input", Input: "", Error: "unknown null input behavior: "},
		{Name: "validation: lower case input", Input: "called on null input", Expected: NullInputBehaviorCalledOnNullInput},
	}

	for _, testCase := range testCases {
		name := testCase.Name
		if name == "" {
			name = fmt.Sprintf("%v null input behavior", testCase.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToNullInputBehavior(testCase.Input)
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

func Test_ToReturnResultsBehavior(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected ReturnResultsBehavior
		Error    string
	}{
		{Input: string(ReturnResultsBehaviorVolatile), Expected: ReturnResultsBehaviorVolatile},
		{Input: string(ReturnResultsBehaviorImmutable), Expected: ReturnResultsBehaviorImmutable},
		{Name: "validation: incorrect return results behavior", Input: "incorrect", Error: "unknown return results behavior: incorrect"},
		{Name: "validation: empty input", Input: "", Error: "unknown return results behavior: "},
		{Name: "validation: lower case input", Input: "volatile", Expected: ReturnResultsBehaviorVolatile},
	}

	for _, testCase := range testCases {
		name := testCase.Name
		if name == "" {
			name = fmt.Sprintf("%v null input behavior", testCase.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToReturnResultsBehavior(testCase.Input)
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

func TestToTraceLevel(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected TraceLevel
		Error    string
	}{
		{Input: string(TraceLevelAlways), Expected: TraceLevelAlways},
		{Input: string(TraceLevelOnEvent), Expected: TraceLevelOnEvent},
		{Input: string(TraceLevelPropagate), Expected: TraceLevelPropagate},
		{Input: string(TraceLevelOff), Expected: TraceLevelOff},
		{Name: "validation: incorrect trace level", Input: "incorrect", Error: "unknown trace level: incorrect"},
		{Name: "validation: empty input", Input: "", Error: "unknown trace level: "},
		{Name: "validation: lower case input", Input: "always", Expected: TraceLevelAlways},
	}

	for _, testCase := range testCases {
		name := testCase.Name
		if name == "" {
			name = fmt.Sprintf("%v trace level", testCase.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToTraceLevel(testCase.Input)
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

func Test_ToMetricLevel(t *testing.T) {
	testCases := []struct {
		Name          string
		Input         string
		Expected      MetricLevel
		ExpectedError string
	}{
		{Input: string(MetricLevelAll), Expected: MetricLevelAll},
		{Input: string(MetricLevelNone), Expected: MetricLevelNone},
		{Name: "validation: incorrect metric level", Input: "incorrect", ExpectedError: "unknown metric level: incorrect"},
		{Name: "validation: empty input", Input: "", ExpectedError: "unknown metric level: "},
		{Name: "validation: lower case input", Input: "all", Expected: MetricLevelAll},
	}

	for _, tc := range testCases {
		tc := tc
		name := tc.Name
		if name == "" {
			name = fmt.Sprintf("%v metric level", tc.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToMetricLevel(tc.Input)
			if tc.ExpectedError != "" {
				assert.Empty(t, value)
				assert.ErrorContains(t, err, tc.ExpectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.Expected, value)
			}
		})
	}
}

func Test_ToAutoEventLogging(t *testing.T) {
	testCases := []struct {
		Name          string
		Input         string
		Expected      AutoEventLogging
		ExpectedError string
	}{
		{Input: string(AutoEventLoggingLogging), Expected: AutoEventLoggingLogging},
		{Input: string(AutoEventLoggingTracing), Expected: AutoEventLoggingTracing},
		{Input: string(AutoEventLoggingAll), Expected: AutoEventLoggingAll},
		{Input: string(AutoEventLoggingOff), Expected: AutoEventLoggingOff},
		{Name: "validation: incorrect auto event logging", Input: "incorrect", ExpectedError: "unknown auto event logging: incorrect"},
		{Name: "validation: empty input", Input: "", ExpectedError: "unknown auto event logging: "},
		{Name: "validation: lower case input", Input: "all", Expected: AutoEventLoggingAll},
	}

	for _, tc := range testCases {
		tc := tc
		name := tc.Name
		if name == "" {
			name = fmt.Sprintf("%v auto event logging", tc.Input)
		}
		t.Run(name, func(t *testing.T) {
			value, err := ToAutoEventLogging(tc.Input)
			if tc.ExpectedError != "" {
				assert.Empty(t, value)
				assert.ErrorContains(t, err, tc.ExpectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.Expected, value)
			}
		})
	}
}

func TestStageLocation(t *testing.T) {
	t.Run("get stage sql", func(t *testing.T) {
		stage := NewSchemaObjectIdentifier("db", "schema", "stage")
		identifier := NewStageLocation(stage, "path/to/file")

		assert.Equal(t, `@"db"."schema"."stage"/path/to/file`, identifier.ToSql())
	})

	t.Run("get stage and path sql", func(t *testing.T) {
		stage := NewSchemaObjectIdentifier("db", "schema", "stage")
		identifier := NewStageLocation(stage, "")

		assert.Equal(t, `@"db"."schema"."stage"`, identifier.ToSql())
	})

	t.Run("empty stage and path returns empty sql", func(t *testing.T) {
		stage := NewSchemaObjectIdentifier("", "", "")
		identifier := NewStageLocation(stage, "")

		assert.Equal(t, "", identifier.ToSql())
	})
}
