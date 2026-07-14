package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeMcpServerSpecification(t *testing.T) {
	t.Run("json equals to yaml specification with different key order and version field", func(t *testing.T) {
		jsonInput := `{"version":1,"tools":[{"name":"sql_exec_tool","type":"SYSTEM_EXECUTE_SQL","title":"SQL Execution Tool","description":"For tests."}]}`
		yamlInput := `tools:
  - title: "SQL Execution Tool"
    name: "sql_exec_tool"
    type: "SYSTEM_EXECUTE_SQL"
    description: "For tests."`

		jsonOutput, err := NormalizeMcpServerSpecification(jsonInput)
		require.NoError(t, err)

		yamlOutput, err := NormalizeMcpServerSpecification(yamlInput)
		require.NoError(t, err)

		require.Equal(t, jsonOutput, yamlOutput)
	})

	t.Run("empty input", func(t *testing.T) {
		got, err := NormalizeMcpServerSpecification("")
		require.NoError(t, err)
		require.Equal(t, "null", got)
	})

	t.Run("invalid input", func(t *testing.T) {
		_, err := NormalizeMcpServerSpecification("{broken")
		require.Error(t, err)
	})
}
