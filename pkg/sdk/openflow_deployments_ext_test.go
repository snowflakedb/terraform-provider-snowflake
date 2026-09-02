package sdk

import "testing"

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
