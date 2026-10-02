package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	id := gatewaysTestIdSchemaObjectIdentifier
	schemaID := id.SchemaId()
	specification := `tools:
  - name: sql_exec_tool
    type: SYSTEM_EXECUTE_SQL
`

	gatewaysTests.Create.
		withDefaultOpts(func() *CreateGatewayOptions {
			return &CreateGatewayOptions{
				name:              id,
				FromSpecification: specification,
			}
		}).
		withExpectedSqlf(
			case_Gateways_sql_Create_basic,
			"CREATE GATEWAY %s FROM SPECIFICATION $$%s$$", id.FullyQualifiedName(), specification,
		).
		withModifyAndExpectedSqlf(
			case_Gateways_sql_Create_all,
			func(opts *CreateGatewayOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("some comment")
			},
			"CREATE GATEWAY IF NOT EXISTS %s COMMENT = 'some comment' FROM SPECIFICATION $$%s$$",
			id.FullyQualifiedName(), specification,
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateGatewayOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE GATEWAY %s FROM SPECIFICATION $$%s$$", id.FullyQualifiedName(), specification,
		)

	gatewaysTests.Alter.
		withModifyAndExpectedSqlf(
			case_Gateways_sql_Alter_Set,
			func(opts *AlterGatewayOptions) {
				opts.IfExists = new(true)
				opts.Set = &GatewaySet{Comment: &StringAllowEmpty{Value: "some comment"}}
			},
			"ALTER GATEWAY IF EXISTS %s SET COMMENT = 'some comment'", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Gateways_sql_Alter_ModifyLiveVersionSet,
			func(opts *AlterGatewayOptions) {
				opts.ModifyLiveVersionSet = &GatewayModifyLiveVersionSet{Specification: specification}
			},
			"ALTER GATEWAY %s MODIFY LIVE VERSION SET SPECIFICATION = $$%s$$",
			id.FullyQualifiedName(), specification,
		)

	gatewaysTests.Drop.
		withExpectedSqlf(case_Gateways_sql_Drop_basic, "DROP GATEWAY %s", id.FullyQualifiedName()).
		withModifyAndExpectedSqlf(
			case_Gateways_sql_Drop_all,
			func(opts *DropGatewayOptions) { opts.IfExists = new(true) },
			"DROP GATEWAY IF EXISTS %s", id.FullyQualifiedName(),
		)

	gatewaysTests.Show.
		withExpectedSql(case_Gateways_sql_Show_basic, "SHOW GATEWAYS").
		withModifyAndExpectedSqlf(
			case_Gateways_sql_Show_all,
			func(opts *ShowGatewayOptions) {
				opts.Like = &Like{Pattern: new("like-pattern")}
				opts.In = &ExtendedIn{In: In{Schema: schemaID}}
				opts.StartsWith = new("starts-with-pattern")
				opts.Limit = &LimitFrom{Rows: new(10), From: new("limit-from")}
			},
			"SHOW GATEWAYS LIKE 'like-pattern' IN SCHEMA %s STARTS WITH 'starts-with-pattern' LIMIT 10 FROM 'limit-from'",
			schemaID.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Gateways_sql_Show_Like,
			func(opts *ShowGatewayOptions) { opts.Like = &Like{Pattern: new("like-pattern")} },
			"SHOW GATEWAYS LIKE 'like-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_Gateways_sql_Show_In,
			func(opts *ShowGatewayOptions) { opts.In = &ExtendedIn{In: In{Schema: schemaID}} },
			"SHOW GATEWAYS IN SCHEMA %s", schemaID.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Gateways_sql_Show_StartsWith,
			func(opts *ShowGatewayOptions) { opts.StartsWith = new("starts-with-pattern") },
			"SHOW GATEWAYS STARTS WITH 'starts-with-pattern'",
		).
		withModifyAndExpectedSqlf(
			case_Gateways_sql_Show_Limit,
			func(opts *ShowGatewayOptions) { opts.Limit = &LimitFrom{Rows: new(10), From: new("limit-from")} },
			"SHOW GATEWAYS LIMIT 10 FROM 'limit-from'",
		)

	gatewaysTests.Describe.withExpectedSqlf(
		case_Gateways_sql_Describe_basic,
		"DESCRIBE GATEWAY %s", id.FullyQualifiedName(),
	)
}

func TestNormalizeGatewaySpecification(t *testing.T) {
	t.Run("equivalent JSON and YAML specifications", func(t *testing.T) {
		jsonSpec := `{"version":1,"tools":[{"name":"sql_exec_tool","type":"SYSTEM_EXECUTE_SQL"}]}`
		yamlSpec := `version: 1
tools:
  - name: sql_exec_tool
    type: SYSTEM_EXECUTE_SQL
`

		jsonOutput, err := NormalizeGatewaySpecification(jsonSpec)
		require.NoError(t, err)

		yamlOutput, err := NormalizeGatewaySpecification(yamlSpec)
		require.NoError(t, err)

		require.JSONEq(t, jsonOutput, yamlOutput)
	})

	t.Run("empty input", func(t *testing.T) {
		got, err := NormalizeGatewaySpecification(" \n\t ")
		require.NoError(t, err)
		require.Equal(t, "{}", got)
	})

	t.Run("invalid input", func(t *testing.T) {
		_, err := NormalizeGatewaySpecification("{broken")
		require.Error(t, err)
	})
}
