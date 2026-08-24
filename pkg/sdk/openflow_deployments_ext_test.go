package sdk

func init() {
	id := openflowDeploymentsTestIdAccountObjectIdentifier
	renameTarget := randomAccountObjectIdentifier()

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
				opts.EventTable = new("MY_DB.PUBLIC.EVENTS")
				opts.DisplayName = new("My Deployment")
				opts.Comment = new("set-comment")
			},
			"CREATE OPENFLOW DEPLOYMENT IF NOT EXISTS %s DEPLOYMENT_TYPE = 'BYOC' VPC_TYPE = 'MANAGED'"+
				" CUSTOM_INGRESS_HOSTNAME = 'ingress.example.com' USE_PRIVATE_LINK = true USE_USER_AUTH_OVER_PRIVATELINK = false"+
				" EVENT_TABLE = 'MY_DB.PUBLIC.EVENTS' DISPLAY_NAME = 'My Deployment' COMMENT = 'set-comment'",
			id.FullyQualifiedName(),
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
					Comment:     new("set-comment"),
					DisplayName: new("My Deployment"),
					EventTable:  new("MY_DB.PUBLIC.EVENTS"),
				}
			},
			"ALTER OPENFLOW DEPLOYMENT %s SET COMMENT = 'set-comment' DISPLAY_NAME = 'My Deployment' EVENT_TABLE = 'MY_DB.PUBLIC.EVENTS'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Alter_Unset,
			func(opts *AlterOpenflowDeploymentOptions) {
				opts.Unset = &OpenflowDeploymentUnset{
					Comment:     new(true),
					DisplayName: new(true),
					EventTable:  new(true),
				}
			},
			"ALTER OPENFLOW DEPLOYMENT %s UNSET COMMENT, DISPLAY_NAME, EVENT_TABLE",
			id.FullyQualifiedName(),
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
			},
			"SHOW OPENFLOW DEPLOYMENTS LIKE 'pattern'",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowDeployments_sql_Show_Like,
			func(opts *ShowOpenflowDeploymentOptions) {
				opts.Like = &Like{Pattern: new("my-deployment%")}
			},
			"SHOW OPENFLOW DEPLOYMENTS LIKE 'my-deployment%%'",
		)

	openflowDeploymentsTests.Describe.
		withExpectedSqlf(
			case_OpenflowDeployments_sql_Describe_basic,
			"DESCRIBE OPENFLOW DEPLOYMENT %s",
			id.FullyQualifiedName(),
		)
}
