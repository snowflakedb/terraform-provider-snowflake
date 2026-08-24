package sdk

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
)

func init() {
	id := openflowRuntimesTestIdSchemaObjectIdentifier
	deploymentId := randomAccountObjectIdentifier()
	roleId := randomAccountObjectIdentifier()
	eaiId := randomAccountObjectIdentifier()
	comment := random.Comment()
	renameTarget := randomSchemaObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifier()

	openflowRuntimesTests.Create.
		withDefaultOpts(func() *CreateOpenflowRuntimeOptions {
			return &CreateOpenflowRuntimeOptions{
				name:          id,
				InDeployment:  deploymentId,
				ExecuteAsRole: roleId,
				NodeType:      OpenflowRuntimeNodeTypeSmall,
				MinNodes:      1,
				MaxNodes:      3,
			}
		}).
		withExpectedSqlf(
			case_OpenflowRuntimes_sql_Create_basic,
			"CREATE OPENFLOW RUNTIME %s IN DEPLOYMENT %s EXECUTE_AS_ROLE = %s NODE_TYPE = 'SMALL' MIN_NODES = 1 MAX_NODES = 3",
			id.FullyQualifiedName(), deploymentId.FullyQualifiedName(), roleId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Create_all,
			func(opts *CreateOpenflowRuntimeOptions) {
				opts.IfNotExists = new(true)
				opts.NodeType = OpenflowRuntimeNodeTypeLarge
				opts.MinNodes = 2
				opts.MaxNodes = 5
				opts.ExternalAccessIntegrations = &OpenflowRuntimeExternalAccessIntegrations{
					ExternalAccessIntegrations: []AccountObjectIdentifier{eaiId},
				}
				opts.DisplayName = new("My Runtime")
				opts.Comment = &comment
			},
			"CREATE OPENFLOW RUNTIME IF NOT EXISTS %s IN DEPLOYMENT %s EXECUTE_AS_ROLE = %s NODE_TYPE = 'LARGE' MIN_NODES = 2 MAX_NODES = 5"+
				" EXTERNAL_ACCESS_INTEGRATIONS = (%s) DISPLAY_NAME = 'My Runtime' COMMENT = '%s'",
			id.FullyQualifiedName(), deploymentId.FullyQualifiedName(), roleId.FullyQualifiedName(), eaiId.FullyQualifiedName(), comment,
		)

	openflowRuntimesTests.Alter.
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_Suspend,
			func(opts *AlterOpenflowRuntimeOptions) { opts.Suspend = new(true) },
			"ALTER OPENFLOW RUNTIME %s SUSPEND", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_Resume,
			func(opts *AlterOpenflowRuntimeOptions) { opts.Resume = new(true) },
			"ALTER OPENFLOW RUNTIME %s RESUME", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_ResumeRecovery,
			func(opts *AlterOpenflowRuntimeOptions) { opts.ResumeRecovery = new(true) },
			"ALTER OPENFLOW RUNTIME %s RESUME RECOVERY", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_Restart,
			func(opts *AlterOpenflowRuntimeOptions) { opts.Restart = new(true) },
			"ALTER OPENFLOW RUNTIME %s RESTART", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_RestartRecovery,
			func(opts *AlterOpenflowRuntimeOptions) { opts.RestartRecovery = new(true) },
			"ALTER OPENFLOW RUNTIME %s RESTART RECOVERY", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_Terminate,
			func(opts *AlterOpenflowRuntimeOptions) { opts.Terminate = new(true) },
			"ALTER OPENFLOW RUNTIME %s TERMINATE", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_TerminateCascade,
			func(opts *AlterOpenflowRuntimeOptions) { opts.TerminateCascade = new(true) },
			"ALTER OPENFLOW RUNTIME %s TERMINATE CASCADE", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_Upgrade,
			func(opts *AlterOpenflowRuntimeOptions) { opts.Upgrade = new(true) },
			"ALTER OPENFLOW RUNTIME %s UPGRADE", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_RenameTo,
			func(opts *AlterOpenflowRuntimeOptions) { opts.RenameTo = &renameTarget },
			"ALTER OPENFLOW RUNTIME %s RENAME TO %s", id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_Set,
			func(opts *AlterOpenflowRuntimeOptions) {
				opts.Set = &OpenflowRuntimeSet{
					MinNodes:      new(2),
					MaxNodes:      new(5),
					ExecuteAsRole: &roleId,
					ExternalAccessIntegrations: &OpenflowRuntimeExternalAccessIntegrations{
						ExternalAccessIntegrations: []AccountObjectIdentifier{eaiId},
					},
					DisplayName: new("Updated Runtime"),
					Comment:     &comment,
				}
			},
			"ALTER OPENFLOW RUNTIME %s SET MIN_NODES = 2 MAX_NODES = 5 EXECUTE_AS_ROLE = %s EXTERNAL_ACCESS_INTEGRATIONS = (%s) DISPLAY_NAME = 'Updated Runtime' COMMENT = '%s'",
			id.FullyQualifiedName(), roleId.FullyQualifiedName(), eaiId.FullyQualifiedName(), comment,
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_Unset,
			func(opts *AlterOpenflowRuntimeOptions) {
				opts.Unset = &OpenflowRuntimeUnset{
					ExecuteAsRole:              new(true),
					ExternalAccessIntegrations: new(true),
					DisplayName:                new(true),
					Comment:                    new(true),
				}
			},
			"ALTER OPENFLOW RUNTIME %s UNSET EXECUTE_AS_ROLE, EXTERNAL_ACCESS_INTEGRATIONS, DISPLAY_NAME, COMMENT",
			id.FullyQualifiedName(),
		)

	openflowRuntimesTests.Drop.
		withExpectedSqlf(
			case_OpenflowRuntimes_sql_Drop_basic,
			"DROP OPENFLOW RUNTIME %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Drop_all,
			func(opts *DropOpenflowRuntimeOptions) {
				opts.IfExists = new(true)
				opts.Cascade = new(true)
			},
			"DROP OPENFLOW RUNTIME IF EXISTS %s CASCADE", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Drop_ifExists",
			func(opts *DropOpenflowRuntimeOptions) { opts.IfExists = new(true) },
			"DROP OPENFLOW RUNTIME IF EXISTS %s", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Drop_cascade",
			func(opts *DropOpenflowRuntimeOptions) { opts.Cascade = new(true) },
			"DROP OPENFLOW RUNTIME %s CASCADE", id.FullyQualifiedName(),
		)

	openflowRuntimesTests.Show.
		withExpectedSql(case_OpenflowRuntimes_sql_Show_basic, "SHOW OPENFLOW RUNTIMES").
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Show_all,
			func(opts *ShowOpenflowRuntimeOptions) {
				opts.Like = &Like{Pattern: new("my-runtime%")}
				opts.In = &In{Schema: schemaId}
			},
			"SHOW OPENFLOW RUNTIMES LIKE 'my-runtime%%' IN SCHEMA %s", schemaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Show_Like,
			func(opts *ShowOpenflowRuntimeOptions) { opts.Like = &Like{Pattern: new("my-runtime%")} },
			"SHOW OPENFLOW RUNTIMES LIKE 'my-runtime%%'",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Show_In,
			func(opts *ShowOpenflowRuntimeOptions) { opts.In = &In{Schema: schemaId} },
			"SHOW OPENFLOW RUNTIMES IN SCHEMA %s", schemaId.FullyQualifiedName(),
		)

	openflowRuntimesTests.Describe.
		withExpectedSqlf(
			case_OpenflowRuntimes_sql_Describe_basic,
			"DESCRIBE OPENFLOW RUNTIME %s", id.FullyQualifiedName(),
		)
}
