package sdk

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			"CREATE OPENFLOW RUNTIME %s IN DEPLOYMENT %s NODE_TYPE = 'SMALL' MIN_NODES = 1 MAX_NODES = 3 EXECUTE_AS_ROLE = %s",
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
			"CREATE OPENFLOW RUNTIME IF NOT EXISTS %s IN DEPLOYMENT %s NODE_TYPE = 'LARGE' MIN_NODES = 2 MAX_NODES = 5 EXECUTE_AS_ROLE = %s"+
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
			func(opts *AlterOpenflowRuntimeOptions) { opts.Upgrade = &OpenflowRuntimeUpgrade{} },
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
			"ALTER OPENFLOW RUNTIME %s SET DISPLAY_NAME = 'Updated Runtime' MIN_NODES = 2 MAX_NODES = 5 EXTERNAL_ACCESS_INTEGRATIONS = (%s) EXECUTE_AS_ROLE = %s COMMENT = '%s'",
			id.FullyQualifiedName(), eaiId.FullyQualifiedName(), roleId.FullyQualifiedName(), comment,
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

	// ADD and REMOVE edit the integration list in place, where SET replaces it wholesale.
	openflowRuntimesTests.Alter.
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_AddExternalAccessIntegrations,
			func(opts *AlterOpenflowRuntimeOptions) {
				opts.AddExternalAccessIntegrations = &OpenflowRuntimeExternalAccessIntegrations{
					ExternalAccessIntegrations: []AccountObjectIdentifier{eaiId},
				}
			},
			"ALTER OPENFLOW RUNTIME %s ADD EXTERNAL_ACCESS_INTEGRATIONS = (%s)",
			id.FullyQualifiedName(), eaiId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Alter_RemoveExternalAccessIntegrations,
			func(opts *AlterOpenflowRuntimeOptions) {
				opts.RemoveExternalAccessIntegrations = &OpenflowRuntimeExternalAccessIntegrations{
					ExternalAccessIntegrations: []AccountObjectIdentifier{eaiId},
				}
			},
			"ALTER OPENFLOW RUNTIME %s REMOVE EXTERNAL_ACCESS_INTEGRATIONS = (%s)",
			id.FullyQualifiedName(), eaiId.FullyQualifiedName(),
		)

	// IF EXISTS is only rendered by ALTER, and the generated Alter SQL cases never set it, so without this
	// the clause is unasserted: retagging the field would not fail any test. Paired with an action that
	// accepts it.
	openflowRuntimesTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_ifExists",
			func(opts *AlterOpenflowRuntimeOptions) {
				opts.IfExists = new(true)
				opts.Suspend = new(true)
			},
			"ALTER OPENFLOW RUNTIME IF EXISTS %s SUSPEND", id.FullyQualifiedName(),
		)

	// RECOVERY and FORCE are independent modifiers, not alternatives: Snowflake accepts all four
	// combinations. The order is fixed, and UPGRADE FORCE RECOVERY is a parse error, so these cases also
	// pin the order the struct fields render in.
	openflowRuntimesTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_upgradeRecovery",
			func(opts *AlterOpenflowRuntimeOptions) {
				opts.Upgrade = &OpenflowRuntimeUpgrade{Recovery: new(true)}
			},
			"ALTER OPENFLOW RUNTIME %s UPGRADE RECOVERY", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_upgradeForce",
			func(opts *AlterOpenflowRuntimeOptions) {
				opts.Upgrade = &OpenflowRuntimeUpgrade{Force: new(true)}
			},
			"ALTER OPENFLOW RUNTIME %s UPGRADE FORCE", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_upgradeRecoveryForce",
			func(opts *AlterOpenflowRuntimeOptions) {
				opts.Upgrade = &OpenflowRuntimeUpgrade{Recovery: new(true), Force: new(true)}
			},
			"ALTER OPENFLOW RUNTIME %s UPGRADE RECOVERY FORCE", id.FullyQualifiedName(),
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
				opts.StartsWith = new("PROD_")
				opts.Limit = &LimitFrom{Rows: new(5), From: new("PROD_A")}
			},
			"SHOW OPENFLOW RUNTIMES LIKE 'my-runtime%%' IN SCHEMA %s STARTS WITH 'PROD_' LIMIT 5 FROM 'PROD_A'",
			schemaId.FullyQualifiedName(),
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
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Show_StartsWith,
			func(opts *ShowOpenflowRuntimeOptions) { opts.StartsWith = new("PROD_") },
			"SHOW OPENFLOW RUNTIMES STARTS WITH 'PROD_'",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowRuntimes_sql_Show_Limit,
			func(opts *ShowOpenflowRuntimeOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			"SHOW OPENFLOW RUNTIMES LIMIT 10",
		)

	openflowRuntimesTests.Describe.
		withExpectedSqlf(
			case_OpenflowRuntimes_sql_Describe_basic,
			"DESCRIBE OPENFLOW RUNTIME %s", id.FullyQualifiedName(),
		)
}

// A runtime with no integrations reports the column as SQL NULL when freshly created, but as the literal
// string `null` after the last integration is removed with ALTER ... REMOVE. Both must read back empty:
// handing `null` to the comma-separated parser produced one integration actually named "null", which
// surfaced as a live failure asserting an emptied list.
func TestParseOpenflowRuntimeExternalAccessIntegrations(t *testing.T) {
	t.Run("shapes that mean no integrations", func(t *testing.T) {
		for _, value := range []string{"", "   ", "null", "[]"} {
			ids, err := ParseOpenflowRuntimeExternalAccessIntegrations(value)
			require.NoError(t, err, "value %q", value)
			assert.Empty(t, ids, "value %q should mean no integrations", value)
		}
	})

	t.Run("json array of quoted names", func(t *testing.T) {
		ids, err := ParseOpenflowRuntimeExternalAccessIntegrations(`["FIRST_EAI","SECOND_EAI"]`)
		require.NoError(t, err)
		assert.Equal(t, []string{"FIRST_EAI", "SECOND_EAI"}, collections.Map(ids, func(i AccountObjectIdentifier) string { return i.Name() }))
	})

	t.Run("single element", func(t *testing.T) {
		ids, err := ParseOpenflowRuntimeExternalAccessIntegrations(`["ONLY_EAI"]`)
		require.NoError(t, err)
		require.Len(t, ids, 1)
		assert.Equal(t, "ONLY_EAI", ids[0].Name())
	})
}
