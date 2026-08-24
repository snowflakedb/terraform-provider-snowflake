package sdk

func init() {
	id := openflowConnectorsTestIdSchemaObjectIdentifier
	runtimeId := randomSchemaObjectIdentifier()
	stageId := randomSchemaObjectIdentifier()
	var stageLocation Location = &StageLocation{stage: stageId, path: "path/"}

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
			func(opts *AlterOpenflowConnectorOptions) { opts.Commit = Bool(true) },
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
	openflowConnectorsTests.Show.
		withExpectedSql(case_OpenflowConnectors_sql_Show_basic, "SHOW OPENFLOW CONNECTORS").
		withModifyAndExpectedSqlf(
			case_OpenflowConnectors_sql_Show_all,
			func(opts *ShowOpenflowConnectorOptions) {
				opts.Like = &Like{Pattern: String("my-connector%")}
				opts.In = &In{Schema: schemaId}
			},
			"SHOW OPENFLOW CONNECTORS LIKE 'my-connector%%' IN SCHEMA %s", schemaId.FullyQualifiedName(),
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
		)

	openflowConnectorsTests.Describe.
		withExpectedSqlf(
			case_OpenflowConnectors_sql_Describe_basic,
			"DESCRIBE OPENFLOW CONNECTOR %s", id.FullyQualifiedName(),
		)
}
