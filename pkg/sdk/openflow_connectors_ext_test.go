package sdk

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	id := openflowConnectorsTestIdSchemaObjectIdentifier
	runtimeId := randomSchemaObjectIdentifier()
	stageId := randomSchemaObjectIdentifier()
	renameTarget := randomSchemaObjectIdentifier()
	var stageLocation Location = &StageLocation{stage: stageId, path: "path/"}
	// Fixed rather than random, so the expected URI below is literal rather than derived from the code it
	// checks. All three parts resolve unquoted.
	sourceConnectorId := NewSchemaObjectIdentifier("MY_DB", "MY_SCHEMA", "SOURCE_CONNECTOR")
	var versionLocation Location = NewOpenflowConnectorVersionLocation(sourceConnectorId, "version$1")

	openflowConnectorsTests.Create.
		withDefaultOpts(func() *CreateOpenflowConnectorOptions {
			return &CreateOpenflowConnectorOptions{
				name:      id,
				InRuntime: runtimeId,
			}
		}).
		withModify(case_OpenflowConnectors_validation_Create_opts_ConflictingFields, func(opts *CreateOpenflowConnectorOptions) {
			opts.FromDefinition = String("mydef")
			opts.From = &stageLocation
		}).
		withExpectedSqlf(
			case_OpenflowConnectors_sql_Create_basic,
			"CREATE OPENFLOW CONNECTOR %s IN RUNTIME %s",
			id.FullyQualifiedName(), runtimeId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Create_all,
			func(opts *CreateOpenflowConnectorOptions) {
				opts.IfNotExists = Bool(true)
				opts.FromDefinition = String("mydef")
				opts.DisplayName = String("My Connector")
				opts.Comment = String("some comment")
			},
			"CREATE OPENFLOW CONNECTOR IF NOT EXISTS %s IN RUNTIME %s FROM DEFINITION mydef DISPLAY_NAME = 'My Connector' COMMENT = 'some comment'",
			id.FullyQualifiedName(), runtimeId.FullyQualifiedName(),
		)

	openflowConnectorsTests.Create.
		withAdditionalSqlCasef(
			"sql_Create_fromDefinition",
			func(opts *CreateOpenflowConnectorOptions) {
				opts.FromDefinition = String("mydef")
			},
			"CREATE OPENFLOW CONNECTOR %s IN RUNTIME %s FROM DEFINITION mydef",
			id.FullyQualifiedName(), runtimeId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_fromStage",
			func(opts *CreateOpenflowConnectorOptions) {
				opts.From = &stageLocation
			},
			`CREATE OPENFLOW CONNECTOR %s IN RUNTIME %s FROM '@\"%s\".\"%s\".\"%s\"/path/'`,
			id.FullyQualifiedName(), runtimeId.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		)

	openflowConnectorsTests.Alter.
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_Start,
			func(opts *AlterOpenflowConnectorOptions) { opts.Start = Bool(true) },
			"ALTER OPENFLOW CONNECTOR %s START", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_Stop,
			func(opts *AlterOpenflowConnectorOptions) { opts.Stop = Bool(true) },
			"ALTER OPENFLOW CONNECTOR %s STOP", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_Terminate,
			func(opts *AlterOpenflowConnectorOptions) { opts.Terminate = Bool(true) },
			"ALTER OPENFLOW CONNECTOR %s TERMINATE", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_Commit,
			// COMMIT carries an optional COMMENT, so it is a struct rather than a bare bool.
			func(opts *AlterOpenflowConnectorOptions) { opts.Commit = &OpenflowConnectorCommit{} },
			"ALTER OPENFLOW CONNECTOR %s COMMIT", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_Abort,
			func(opts *AlterOpenflowConnectorOptions) { opts.Abort = Bool(true) },
			"ALTER OPENFLOW CONNECTOR %s ABORT", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_Set,
			func(opts *AlterOpenflowConnectorOptions) {
				opts.Set = &OpenflowConnectorSet{
					DisplayName: String("Updated Connector"),
					Comment:     String("some comment"),
				}
			},
			"ALTER OPENFLOW CONNECTOR %s SET DISPLAY_NAME = 'Updated Connector' COMMENT = 'some comment'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_Unset,
			func(opts *AlterOpenflowConnectorOptions) {
				opts.Unset = &OpenflowConnectorUnset{
					DisplayName: Bool(true),
					Comment:     Bool(true),
				}
			},
			"ALTER OPENFLOW CONNECTOR %s UNSET DISPLAY_NAME, COMMENT", id.FullyQualifiedName(),
		)

	openflowConnectorsTests.Drop.
		withExpectedSqlf(
			case_OpenflowConnectors_sql_Drop_basic,
			"DROP OPENFLOW CONNECTOR %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Drop_all,
			func(opts *DropOpenflowConnectorOptions) { opts.IfExists = Bool(true) },
			"DROP OPENFLOW CONNECTOR IF EXISTS %s", id.FullyQualifiedName(),
		)

	schemaId := randomDatabaseObjectIdentifier()
	// Actions this PR adds. TERMINATE FORCE is the escape hatch for a connector stuck in TERMINATING.
	// ADD VERSION promotes an imported config straight to the default, where ADD LIVE VERSION opens an
	// editable version that later needs a COMMIT.
	openflowConnectorsTests.Alter.
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_TerminateForce,
			func(opts *AlterOpenflowConnectorOptions) { opts.TerminateForce = Bool(true) },
			"ALTER OPENFLOW CONNECTOR %s TERMINATE FORCE", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_RenameTo,
			func(opts *AlterOpenflowConnectorOptions) { opts.RenameTo = &renameTarget },
			"ALTER OPENFLOW CONNECTOR %s RENAME TO %s",
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_AddLiveVersion,
			func(opts *AlterOpenflowConnectorOptions) {
				opts.AddLiveVersion = &OpenflowConnectorAddLiveVersion{VersionAlias: String("v2")}
			},
			"ALTER OPENFLOW CONNECTOR %s ADD LIVE VERSION v2 FROM LAST", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_AddVersion,
			func(opts *AlterOpenflowConnectorOptions) {
				opts.AddVersion = &OpenflowConnectorAddVersion{From: &stageLocation}
			},
			`ALTER OPENFLOW CONNECTOR %s ADD VERSION FROM '@\"%s\".\"%s\".\"%s\"/path/'`,
			id.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_Push,
			func(opts *AlterOpenflowConnectorOptions) {
				opts.Push = &OpenflowConnectorPush{Username: "u", Password: "p", Name: "n", Email: "e@x.com"}
			},
			"ALTER OPENFLOW CONNECTOR %s PUSH USERNAME = 'u' PASSWORD = 'p' NAME = 'n' EMAIL = 'e@x.com'",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Alter_Pull,
			func(opts *AlterOpenflowConnectorOptions) { opts.Pull = Bool(true) },
			"ALTER OPENFLOW CONNECTOR %s PULL", id.FullyQualifiedName(),
		)

	// COMMIT and ADD LIVE VERSION are structs rather than bare keywords because of their COMMENT, so
	// without these cases the struct shape is unasserted.
	openflowConnectorsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_commitWithComment",
			func(opts *AlterOpenflowConnectorOptions) {
				opts.Commit = &OpenflowConnectorCommit{Comment: String("promoting v2")}
			},
			"ALTER OPENFLOW CONNECTOR %s COMMIT COMMENT = 'promoting v2'", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_addLiveVersionWithAliasAndComment",
			func(opts *AlterOpenflowConnectorOptions) {
				opts.AddLiveVersion = &OpenflowConnectorAddLiveVersion{
					VersionAlias: String("v2"),
					Comment:      String("editing config"),
				}
			},
			"ALTER OPENFLOW CONNECTOR %s ADD LIVE VERSION v2 FROM LAST COMMENT = 'editing config'",
			id.FullyQualifiedName(),
		).
		// Why ADD VERSION takes a *Location: StageLocation only renders @<fqn>[/path], and a connector
		// version lives at a snow:// URI with no stage.
		withAdditionalSqlCasef(
			"sql_Alter_addVersionFromConnectorVersion",
			func(opts *AlterOpenflowConnectorOptions) {
				opts.AddVersion = &OpenflowConnectorAddVersion{From: &versionLocation}
			},
			"ALTER OPENFLOW CONNECTOR %s ADD VERSION FROM 'snow://openflow_connector/MY_DB.MY_SCHEMA.SOURCE_CONNECTOR/versions/version$1/'",
			id.FullyQualifiedName(),
		)

	// DEFAULT_VERSION takes a bare FIRST or LAST keyword, or a single-quoted version name, which is why
	// it is a struct rather than one string field.
	openflowConnectorsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_setDefaultVersionLast",
			func(opts *AlterOpenflowConnectorOptions) {
				opts.Set = &OpenflowConnectorSet{DefaultVersion: &OpenflowConnectorDefaultVersion{Last: Bool(true)}}
			},
			"ALTER OPENFLOW CONNECTOR %s SET DEFAULT_VERSION = LAST", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setDefaultVersionFirst",
			func(opts *AlterOpenflowConnectorOptions) {
				opts.Set = &OpenflowConnectorSet{DefaultVersion: &OpenflowConnectorDefaultVersion{First: Bool(true)}}
			},
			"ALTER OPENFLOW CONNECTOR %s SET DEFAULT_VERSION = FIRST", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setDefaultVersionNamed",
			func(opts *AlterOpenflowConnectorOptions) {
				opts.Set = &OpenflowConnectorSet{DefaultVersion: &OpenflowConnectorDefaultVersion{Version: String("VERSION$2")}}
			},
			"ALTER OPENFLOW CONNECTOR %s SET DEFAULT_VERSION = 'VERSION$2'", id.FullyQualifiedName(),
		)

	// Boundary cases the migration lost by recombining variants. ADD LIVE VERSION's alias and PUSH's TO and
	// COMMENT are all optional, so "none of the optionals" and "all of the optionals" are the two shapes
	// that catch a field rendering when it should not, or being dropped when it should not. Every other case
	// here sets some but not all of them.
	openflowConnectorsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_addLiveVersionBare",
			func(opts *AlterOpenflowConnectorOptions) {
				opts.AddLiveVersion = &OpenflowConnectorAddLiveVersion{}
			},
			"ALTER OPENFLOW CONNECTOR %s ADD LIVE VERSION FROM LAST", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_pushAllFields",
			func(opts *AlterOpenflowConnectorOptions) {
				opts.Push = &OpenflowConnectorPush{
					To:       String("git://example.com/repo/branches/main"),
					Username: "u", Password: "p", Name: "n", Email: "e@x.com",
					Comment: String("sync back"),
				}
			},
			"ALTER OPENFLOW CONNECTOR %s PUSH TO 'git://example.com/repo/branches/main'"+
				" USERNAME = 'u' PASSWORD = 'p' NAME = 'n' EMAIL = 'e@x.com' COMMENT = 'sync back'",
			id.FullyQualifiedName(),
		)

	openflowConnectorsTests.Execute.
		withExpectedSqlf(
			case_OpenflowConnectors_sql_Execute_basic,
			"EXECUTE OPENFLOW CONNECTOR %s VALIDATE CONFIGURATION", id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Execute_fromStage",
			func(opts *ExecuteOpenflowConnectorOptions) { opts.From = &stageLocation },
			`EXECUTE OPENFLOW CONNECTOR %s VALIDATE CONFIGURATION FROM '@\"%s\".\"%s\".\"%s\"/path/'`,
			id.FullyQualifiedName(), stageId.DatabaseName(), stageId.SchemaName(), stageId.Name(),
		).
		withAdditionalSqlCasef(
			"sql_Execute_singleStep",
			func(opts *ExecuteOpenflowConnectorOptions) { opts.Step = String("check_connectivity") },
			"EXECUTE OPENFLOW CONNECTOR %s VALIDATE CONFIGURATION STEP 'check_connectivity'", id.FullyQualifiedName(),
		)

	// IF EXISTS is only rendered by ALTER, and the generated Alter SQL cases never set it, so without this
	// the clause is unasserted: retagging the field would not fail any test. Paired with an action that
	// accepts it.
	openflowConnectorsTests.Alter.
		withAdditionalSqlCasef(
			"sql_Alter_ifExists",
			func(opts *AlterOpenflowConnectorOptions) {
				opts.IfExists = Bool(true)
				opts.Stop = Bool(true)
			},
			"ALTER OPENFLOW CONNECTOR IF EXISTS %s STOP", id.FullyQualifiedName(),
		)

	openflowConnectorsTests.Show.
		withExpectedSql(case_OpenflowConnectors_sql_Show_basic, "SHOW OPENFLOW CONNECTORS").
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Show_all,
			func(opts *ShowOpenflowConnectorOptions) {
				opts.Like = &Like{Pattern: String("my-connector%")}
				opts.In = &In{Schema: schemaId}
				opts.StartsWith = String("PROD_")
				opts.Limit = &LimitFrom{Rows: Int(5), From: String("PROD_A")}
			},
			"SHOW OPENFLOW CONNECTORS LIKE 'my-connector%%' IN SCHEMA %s STARTS WITH 'PROD_' LIMIT 5 FROM 'PROD_A'",
			schemaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Show_Like,
			func(opts *ShowOpenflowConnectorOptions) { opts.Like = &Like{Pattern: String("my-connector%")} },
			"SHOW OPENFLOW CONNECTORS LIKE 'my-connector%%'",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Show_In,
			func(opts *ShowOpenflowConnectorOptions) { opts.In = &In{Schema: schemaId} },
			"SHOW OPENFLOW CONNECTORS IN SCHEMA %s", schemaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Show_StartsWith,
			func(opts *ShowOpenflowConnectorOptions) { opts.StartsWith = String("PROD_") },
			"SHOW OPENFLOW CONNECTORS STARTS WITH 'PROD_'",
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Show_Limit,
			func(opts *ShowOpenflowConnectorOptions) { opts.Limit = &LimitFrom{Rows: Int(10)} },
			"SHOW OPENFLOW CONNECTORS LIMIT 10",
		)

	openflowConnectorsTests.Describe.
		withExpectedSqlf(
			case_OpenflowConnectors_sql_Describe_basic,
			"DESCRIBE OPENFLOW CONNECTOR %s", id.FullyQualifiedName(),
		)

	// LIMIT is the only optional field, so the all and Limit cases render the same statement. Bare
	// LIMIT <n> rather than LimitFrom, because LIKE, STARTS WITH and LIMIT <n> FROM '<s>' are all parse
	// errors on this command.
	openflowConnectorsTests.ShowVersions.
		withExpectedSqlf(
			case_OpenflowConnectors_sql_ShowVersions_basic,
			"SHOW VERSIONS IN OPENFLOW CONNECTOR %s", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_ShowVersions_all,
			func(opts *ShowVersionsOpenflowConnectorOptions) { opts.Limit = Int(5) },
			"SHOW VERSIONS IN OPENFLOW CONNECTOR %s LIMIT 5", id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_ShowVersions_Limit,
			func(opts *ShowVersionsOpenflowConnectorOptions) { opts.Limit = Int(5) },
			"SHOW VERSIONS IN OPENFLOW CONNECTOR %s LIMIT 5", id.FullyQualifiedName(),
		)
}

// SHOW VERSIONS' SQL and validation cases are generated and registered above; this covers the one thing
// the framework does not reach, the row contract. Even name is nullable: an uncommitted version has none.
func TestOpenflowConnectors_ShowVersionsRowConversion(t *testing.T) {
	t.Run("nulls stay nils", func(t *testing.T) {
		row := openflowConnectorVersionRow{
			CreatedOn:   time.Now(),
			IsDefault:   true,
			IsLive:      true,
			LocationUri: `snow://openflow_connector/MY_DB.MY_SCHEMA."my_connector"/versions/live/`,
		}
		converted, err := row.convert()
		require.NoError(t, err)
		assert.Nil(t, converted.Name)
		assert.Nil(t, converted.Alias)
		assert.Nil(t, converted.Comment)
		assert.Nil(t, converted.SourceLocationUri)
		assert.Nil(t, converted.GitCommitHash)
		assert.True(t, converted.IsDefault)
		assert.Equal(t, `snow://openflow_connector/MY_DB.MY_SCHEMA."my_connector"/versions/live/`, converted.LocationUri)
	})

	t.Run("present values map", func(t *testing.T) {
		row := openflowConnectorVersionRow{
			Name:    sql.NullString{String: "version$1", Valid: true},
			Comment: sql.NullString{String: "Initial version", Valid: true},
		}
		converted, err := row.convert()
		require.NoError(t, err)
		require.NotNil(t, converted.Name)
		assert.Equal(t, "version$1", *converted.Name)
		require.NotNil(t, converted.Comment)
		assert.Equal(t, "Initial version", *converted.Comment)
	})
}

// The expected URIs are the forms Snowflake reports in location_uri. Note the asymmetry in the second
// case - bare database and schema, quoted connector - which is why callers do not assemble the string.
func TestOpenflowConnectorVersionLocation_ToSql(t *testing.T) {
	tests := []struct {
		name      string
		connector SchemaObjectIdentifier
		version   string
		expected  string
	}{
		{
			name:      "all parts resolve unquoted",
			connector: NewSchemaObjectIdentifier("OPENFLOW", "OPENFLOW", "MY_CONNECTOR"),
			version:   "version$1",
			expected:  "snow://openflow_connector/OPENFLOW.OPENFLOW.MY_CONNECTOR/versions/version$1/",
		},
		{
			name:      "only the connector needs quoting",
			connector: NewSchemaObjectIdentifier("OPENFLOW", "OPENFLOW", "farhan_test12_v4"),
			version:   "live",
			expected:  `snow://openflow_connector/OPENFLOW.OPENFLOW."farhan_test12_v4"/versions/live/`,
		},
		{
			name:      "hyphenated connector",
			connector: NewSchemaObjectIdentifier("OPENFLOW", "OPENFLOW", "farhan-connector-11"),
			version:   "version$1",
			expected:  `snow://openflow_connector/OPENFLOW.OPENFLOW."farhan-connector-11"/versions/version$1/`,
		},
		{
			name:      "every part needs quoting",
			connector: NewSchemaObjectIdentifier("my_db", "my_schema", "my_connector"),
			version:   "live",
			expected:  `snow://openflow_connector/"my_db"."my_schema"."my_connector"/versions/live/`,
		},
		{
			// SHOW VERSIONS names it VERSION$1, but the URI segment is version$1 and matched case-sensitively.
			name:      "version name as SHOW VERSIONS reports it",
			connector: NewSchemaObjectIdentifier("OPENFLOW", "OPENFLOW", "MY_CONNECTOR"),
			version:   "VERSION$1",
			expected:  "snow://openflow_connector/OPENFLOW.OPENFLOW.MY_CONNECTOR/versions/version$1/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, NewOpenflowConnectorVersionLocation(tt.connector, tt.version).ToSql())
		})
	}
}
