package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	id := mcpServersTestIdSchemaObjectIdentifier
	schemaId := id.SchemaId()
	spec := `tools:
  - title: "SQL Execution Tool"
    name: "sql_exec_tool"
    type: "SYSTEM_EXECUTE_SQL"
    description: "Unit test."
`

	mcpServersTests.Create.
		withDefaultOpts(func() *CreateMcpServerOptions {
			return &CreateMcpServerOptions{
				name:              id,
				FromSpecification: spec,
			}
		}).
		withExpectedSqlf(
			case_McpServers_sql_Create_basic,
			"CREATE MCP SERVER %s FROM SPECIFICATION $$%s$$", id.FullyQualifiedName(), spec,
		).
		withModifyAndExpectedSqlf(
			case_McpServers_sql_Create_all,
			func(opts *CreateMcpServerOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("some comment")
			},
			"CREATE MCP SERVER IF NOT EXISTS %s COMMENT = 'some comment' FROM SPECIFICATION $$%s$$",
			id.FullyQualifiedName(), spec,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateMcpServerOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE MCP SERVER %s FROM SPECIFICATION $$%s$$", id.FullyQualifiedName(), spec,
		)

	mcpServersTests.Drop.
		withExpectedSqlf(
			case_McpServers_sql_Drop_basic,
			"DROP MCP SERVER %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_McpServers_sql_Drop_all,
			func(opts *DropMcpServerOptions) { opts.IfExists = new(true) },
			"DROP MCP SERVER IF EXISTS %s", id.FullyQualifiedName(),
		)

	mcpServersTests.Show.
		withExpectedSql(case_McpServers_sql_Show_basic, "SHOW MCP SERVERS").
		withModifyAndExpectedSqlf(
			case_McpServers_sql_Show_all,
			func(opts *ShowMcpServerOptions) {
				opts.Like = &Like{Pattern: new("like-pattern")}
				opts.In = &In{Schema: schemaId}
			},
			"SHOW MCP SERVERS LIKE 'like-pattern' IN SCHEMA %s",
			schemaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_McpServers_sql_Show_Like,
			func(opts *ShowMcpServerOptions) { opts.Like = &Like{Pattern: new("like-pattern")} },
			"SHOW MCP SERVERS LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_McpServers_sql_Show_In,
			func(opts *ShowMcpServerOptions) { opts.In = &In{Schema: schemaId} },
			"SHOW MCP SERVERS IN SCHEMA %s", schemaId.FullyQualifiedName(),
		)

	mcpServersTests.Describe.
		withExpectedSqlf(
			case_McpServers_sql_Describe_basic,
			"DESCRIBE MCP SERVER %s", id.FullyQualifiedName(),
		)
}

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

		require.NotEmpty(t, jsonOutput)
		require.NotEmpty(t, yamlOutput)
		require.JSONEq(t, jsonOutput, yamlOutput)
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
