package sdk

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	id := openflowDeploymentsTestIdAccountObjectIdentifier
	renameTarget := randomAccountObjectIdentifier()
	eventTableId := randomSchemaObjectIdentifier()

	openflowDeploymentsTests.Create.
		withDefaultOpts(func() *CreateOpenflowDeploymentOptions {
			return &CreateOpenflowDeploymentOptions{
				name:           id,
				DeploymentType: OpenflowDeploymentTypeSnowflake,
			}
		}).
		withExpectedSqlf(
			case_OpenflowDeployments_sql_Create_basic,
			"CREATE OPENFLOW DEPLOYMENT %s DEPLOYMENT_TYPE = 'SNOWFLAKE'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Create_all,
			func(opts *CreateOpenflowDeploymentOptions) {
				opts.IfNotExists = new(true)
				opts.DeploymentType = OpenflowDeploymentTypeByoc
				opts.VpcType = new(OpenflowVpcTypeManaged)
				opts.CustomIngressHostname = new("ingress.example.com")
				opts.UsePrivateLink = new(true)
				opts.UseUserAuthOverPrivatelink = new(false)
				opts.EventTable = &OpenflowDeploymentEventTable{EventTable: &eventTableId}
				opts.DisplayName = new("My Deployment")
				opts.Comment = new("set-comment")
			},
			"CREATE OPENFLOW DEPLOYMENT IF NOT EXISTS %s DEPLOYMENT_TYPE = 'BYOC' VPC_TYPE = 'MANAGED'"+
				" USE_PRIVATE_LINK = true USE_USER_AUTH_OVER_PRIVATELINK = false CUSTOM_INGRESS_HOSTNAME = 'ingress.example.com'"+
				` DISPLAY_NAME = 'My Deployment' COMMENT = 'set-comment' EVENT_TABLE = '\"%s\".\"%s\".\"%s\"'`,
			id.FullyQualifiedName(), eventTableId.DatabaseName(), eventTableId.SchemaName(), eventTableId.Name(),
		)

	openflowDeploymentsTests.Alter.
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Alter_Upgrade,
			func(opts *AlterOpenflowDeploymentOptions) { opts.Upgrade = new(true) },
			"ALTER OPENFLOW DEPLOYMENT %s UPGRADE",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Alter_Terminate,
			func(opts *AlterOpenflowDeploymentOptions) { opts.Terminate = new(true) },
			"ALTER OPENFLOW DEPLOYMENT %s TERMINATE",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Alter_RenameTo,
			func(opts *AlterOpenflowDeploymentOptions) { opts.RenameTo = &renameTarget },
			"ALTER OPENFLOW DEPLOYMENT %s RENAME TO %s",
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Alter_Set,
			func(opts *AlterOpenflowDeploymentOptions) {
				opts.Set = &OpenflowDeploymentSet{
					DisplayName: new("My Deployment"),
					Comment:     new("set-comment"),
					EventTable:  &OpenflowDeploymentEventTable{EventTable: &eventTableId},
				}
			},
			`ALTER OPENFLOW DEPLOYMENT %s SET DISPLAY_NAME = 'My Deployment' COMMENT = 'set-comment' EVENT_TABLE = '\"%s\".\"%s\".\"%s\"'`,
			id.FullyQualifiedName(), eventTableId.DatabaseName(), eventTableId.SchemaName(), eventTableId.Name(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Alter_Unset,
			func(opts *AlterOpenflowDeploymentOptions) {
				opts.Unset = &OpenflowDeploymentUnset{
					DisplayName: new(true),
					Comment:     new(true),
					EventTable:  new(true),
				}
			},
			"ALTER OPENFLOW DEPLOYMENT %s UNSET DISPLAY_NAME, COMMENT, EVENT_TABLE",
			id.FullyQualifiedName(),
		)

	// EVENT_TABLE = NONE is a bare keyword, so it must not be quoted the way a table name is. It also is
	// not the same as UNSET: NONE drops all events for the deployment, UNSET falls back to the account's
	// default event table.
	openflowDeploymentsTests.Create.
		withAdditionalSqlCasef(
			"sql_Create_eventTableNone",
			func(opts *CreateOpenflowDeploymentOptions) {
				opts.EventTable = &OpenflowDeploymentEventTable{None: new(true)}
			},
			"CREATE OPENFLOW DEPLOYMENT %s DEPLOYMENT_TYPE = 'SNOWFLAKE' EVENT_TABLE = NONE",
			id.FullyQualifiedName(),
		)

	openflowDeploymentsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_setEventTableNone",
			func(opts *AlterOpenflowDeploymentOptions) {
				opts.Set = &OpenflowDeploymentSet{EventTable: &OpenflowDeploymentEventTable{None: new(true)}}
			},
			"ALTER OPENFLOW DEPLOYMENT %s SET EVENT_TABLE = NONE",
			id.FullyQualifiedName(),
		)

	// IF EXISTS is only rendered by ALTER, and the generated Alter SQL cases never set it, so without this
	// the clause is unasserted: retagging the field would not fail any test. Paired with an action that
	// accepts it.
	openflowDeploymentsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_ifExists",
			func(opts *AlterOpenflowDeploymentOptions) {
				opts.IfExists = new(true)
				opts.Upgrade = new(true)
			},
			"ALTER OPENFLOW DEPLOYMENT IF EXISTS %s UPGRADE", id.FullyQualifiedName(),
		)

	openflowDeploymentsTests.Drop.
		withExpectedSqlf(
			case_OpenflowDeployments_sql_Drop_basic,
			"DROP OPENFLOW DEPLOYMENT %s",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Drop_all,
			func(opts *DropOpenflowDeploymentOptions) { opts.IfExists = new(true) },
			"DROP OPENFLOW DEPLOYMENT IF EXISTS %s",
			id.FullyQualifiedName(),
		)

	openflowDeploymentsTests.Show.
		withExpectedSql(
			case_OpenflowDeployments_sql_Show_basic,
			"SHOW OPENFLOW DEPLOYMENTS",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Show_all,
			func(opts *ShowOpenflowDeploymentOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
				opts.StartsWith = new("PROD_")
				opts.Limit = &LimitFrom{Rows: new(5), From: new("PROD_A")}
			},
			"SHOW OPENFLOW DEPLOYMENTS LIKE 'pattern' STARTS WITH 'PROD_' LIMIT 5 FROM 'PROD_A'",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Show_Like,
			func(opts *ShowOpenflowDeploymentOptions) {
				opts.Like = &Like{Pattern: new("my-deployment%")}
			},
			"SHOW OPENFLOW DEPLOYMENTS LIKE 'my-deployment%%'",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Show_StartsWith,
			func(opts *ShowOpenflowDeploymentOptions) {
				opts.StartsWith = new("PROD_")
			},
			"SHOW OPENFLOW DEPLOYMENTS STARTS WITH 'PROD_'",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Show_Limit,
			func(opts *ShowOpenflowDeploymentOptions) {
				opts.Limit = &LimitFrom{Rows: new(10)}
			},
			"SHOW OPENFLOW DEPLOYMENTS LIMIT 10",
		)

	openflowDeploymentsTests.Describe.
		withExpectedSqlf(
			case_OpenflowDeployments_sql_Describe_basic,
			"DESCRIBE OPENFLOW DEPLOYMENT %s",
			id.FullyQualifiedName(),
		)
}

// A non-opts test: ShowParameters is a hand-written interface method rather than a generated operation, so
// the generator emits no cases for it. This branch adds ParametersIn.OpenflowDeployment, and without these
// two cases that scope has no coverage at all: dropping it from ParametersIn.validate()'s anyValueSet guard
// leaves the whole suite green while every ShowParameters call starts failing validation.
//
// The LIKE variant is the one the deployment resource actually issues, since EVENT_TABLE appears in neither
// SHOW nor DESCRIBE output and SHOW PARAMETERS is the only way to read it back.
func TestOpenflowDeployments_ShowParameters(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("in openflow deployment", func(t *testing.T) {
		opts := &ShowParametersOptions{
			In: &ParametersIn{OpenflowDeployment: id},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW PARAMETERS IN OPENFLOW DEPLOYMENT %s", id.FullyQualifiedName())
	})

	t.Run("like event table in openflow deployment", func(t *testing.T) {
		opts := &ShowParametersOptions{
			Like: &Like{Pattern: String("EVENT_TABLE")},
			In:   &ParametersIn{OpenflowDeployment: id},
		}
		assertOptsValidAndSqlEqualsf(t, opts, "SHOW PARAMETERS LIKE 'EVENT_TABLE' IN OPENFLOW DEPLOYMENT %s", id.FullyQualifiedName())
	})
}

func Test_quoteOpenflowEventTableValue(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected string
		Error    string
	}{
		{Input: `snowflake.db.schema`, Expected: `"SNOWFLAKE"."DB"."SCHEMA"`},
		{Input: `snowflake.db."schema"`, Expected: `"SNOWFLAKE"."DB"."schema"`},
		{Input: `snowflake.db."sch.ema"`, Expected: `"SNOWFLAKE"."DB"."sch.ema"`},
		// FullyQualifiedName does not double quotes inside a part, so this cannot round-trip through ParseSchemaObjectIdentifier.
		{Input: `snowflake.db."Sch""ema"`, Expected: `"SNOWFLAKE"."DB"."Sch"ema"`},
		{Input: `"snowflake".db."schema"`, Expected: `"snowflake"."DB"."schema"`},
		{Input: `"SNOWFLAKE"."DB"."SCHEMA"`, Expected: `"SNOWFLAKE"."DB"."SCHEMA"`},
		{Input: `SNOWFLAKE.TELEMETRY.EVENTS`, Expected: `"SNOWFLAKE"."TELEMETRY"."EVENTS"`},
		{Input: ``, Expected: ``},
		{Input: `NONE`, Expected: `NONE`},
		{Input: `none`, Expected: `none`},
		{Input: `snowflake.db`, Error: `unexpected number of parts 2 in identifier snowflake.db, expected 3`},
		{Input: `snowflake.db.schema.extra`, Error: `unexpected number of parts 4 in identifier snowflake.db.schema.extra, expected 3`},
		{Input: `snowflake..schema`, Error: "empty part"},
		{Input: `snowflake.db.`, Error: "empty part"},
		{Input: `snowflake.db.""`, Error: "empty part"},
		{Input: `.db.schema`, Error: "empty part"},
		{Input: `snowflake.db."schema`, Error: `extraneous or missing " in quoted-field`},
		{Input: `snowflake.d"b.schema`, Error: `bare " in non-quoted-field`},
		{Input: `snowflake.db."schema"extra`, Error: `extraneous or missing " in quoted-field`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Input, func(t *testing.T) {
			got, err := quoteOpenflowEventTableValue(testCase.Input)
			if testCase.Error != "" {
				require.ErrorContains(t, err, testCase.Error)
				return
			}
			require.NoError(t, err)
			require.Equal(t, testCase.Expected, got)
			if got == "" || strings.EqualFold(got, "NONE") || strings.Contains(testCase.Input, `""`) {
				return
			}
			id, err := ParseSchemaObjectIdentifier(got)
			require.NoError(t, err)
			require.Equal(t, got, id.FullyQualifiedName())
		})
	}
}

func Test_foldUnquotedIdentifierParts(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{Name: "empty", Input: ``, Expected: ``},
		{Name: "lowercase unquoted", Input: `snowflake.db.schema`, Expected: `SNOWFLAKE.DB.SCHEMA`},
		{Name: "uppercase unquoted stays", Input: `SNOWFLAKE.DB.SCHEMA`, Expected: `SNOWFLAKE.DB.SCHEMA`},
		{Name: "digits and underscores", Input: `my_db1.sch_2`, Expected: `MY_DB1.SCH_2`},
		{Name: "quoted part keeps case", Input: `snowflake."schema"`, Expected: `SNOWFLAKE."schema"`},
		{Name: "quoted part with dot", Input: `db."sch.ema"`, Expected: `DB."sch.ema"`},
		{Name: "escaped quote inside quoted part", Input: `db."Sch""ema"`, Expected: `DB."Sch""ema"`},
		{Name: "escaped quote followed by unquoted", Input: `ab"c""d"e`, Expected: `AB"c""d"E`},
		{Name: "already quoted parts", Input: `"snowflake"."db"."schema"`, Expected: `"snowflake"."db"."schema"`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			require.Equal(t, testCase.Expected, foldUnquotedIdentifierParts(testCase.Input))
		})
	}
}

func Test_quoteOpenflowEventTableParameter(t *testing.T) {
	param := &Parameter{
		Key:     "EVENT_TABLE",
		Value:   `snowflake.db."schema"`,
		Default: `snowflake.db.schema`,
		Level:   ParameterTypeSnowflakeDefault,
	}
	require.NoError(t, quoteOpenflowEventTableParameter(param))
	require.Equal(t, `"SNOWFLAKE"."DB"."schema"`, param.Value)
	require.Equal(t, `"SNOWFLAKE"."DB"."SCHEMA"`, param.Default)

	details, err := ToOpenflowDeploymentParametersDetails([]*Parameter{param})
	require.NoError(t, err)
	require.Equal(t, `"SNOWFLAKE"."DB"."schema"`, details.EventTable.Value.FullyQualifiedName())
	require.Equal(t, `"SNOWFLAKE"."DB"."SCHEMA"`, details.EventTable.Default.FullyQualifiedName())

	t.Run("rewrites value and default independently", func(t *testing.T) {
		mixed := &Parameter{
			Key:     "EVENT_TABLE",
			Value:   `"snowflake".telemetry.events`,
			Default: `snowflake.db."schema"`,
		}
		require.NoError(t, quoteOpenflowEventTableParameter(mixed))
		require.Equal(t, `"snowflake"."TELEMETRY"."EVENTS"`, mixed.Value)
		require.Equal(t, `"SNOWFLAKE"."DB"."schema"`, mixed.Default)
	})

	t.Run("empty and NONE stay unchanged", func(t *testing.T) {
		unchanged := &Parameter{Key: "EVENT_TABLE", Value: "", Default: "NONE"}
		require.NoError(t, quoteOpenflowEventTableParameter(unchanged))
		require.Empty(t, unchanged.Value)
		require.Equal(t, "NONE", unchanged.Default)
	})

	t.Run("invalid value keeps default unparsed", func(t *testing.T) {
		invalid := &Parameter{Key: "EVENT_TABLE", Value: "not-an-id", Default: `snowflake.db.schema`}
		require.Error(t, quoteOpenflowEventTableParameter(invalid))
		require.Equal(t, "not-an-id", invalid.Value)
		require.Equal(t, `snowflake.db.schema`, invalid.Default)
	})

	t.Run("invalid default keeps value unparsed", func(t *testing.T) {
		invalid := &Parameter{Key: "EVENT_TABLE", Value: `snowflake.db.schema`, Default: "not-an-id"}
		require.Error(t, quoteOpenflowEventTableParameter(invalid))
		require.Equal(t, `snowflake.db.schema`, invalid.Value)
		require.Equal(t, "not-an-id", invalid.Default)
	})
}
