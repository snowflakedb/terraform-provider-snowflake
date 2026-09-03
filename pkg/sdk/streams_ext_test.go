package sdk

func init() {
	id := streamsTestIdSchemaObjectIdentifier
	tableId := randomSchemaObjectIdentifier()
	externalTableId := randomSchemaObjectIdentifier()
	stageId := randomSchemaObjectIdentifier()
	viewId := randomSchemaObjectIdentifier()
	sourceStreamId := randomSchemaObjectIdentifier()
	streamId := randomSchemaObjectIdentifier()
	schemaId := randomDatabaseObjectIdentifier()
	tagId := NewAccountObjectIdentifier("tag1")
	tagId2 := NewAccountObjectIdentifier("tag2")
	timestamp := "2024-09-25 06:16:10.359 -0700"
	queryId := "0111447d-0905-8a5c-0062-f3820281547a"

	streamsTests.CreateOnTable.
		withDefaultOpts(func() *CreateOnTableStreamOptions {
			return &CreateOnTableStreamOptions{
				name:    id,
				TableId: tableId,
			}
		}).
		withExpectedSqlf(
			case_Streams_sql_CreateOnTable_basic,
			"CREATE STREAM %s ON TABLE %s", id.FullyQualifiedName(), tableId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_CreateOnTable_all,
			func(opts *CreateOnTableStreamOptions) {
				opts.IfNotExists = new(true)
				opts.On = &OnStream{At: new(true), Statement: OnStreamStatement{Stream: new("123")}}
				opts.AppendOnly = new(true)
				opts.ShowInitialRows = new(true)
				opts.Comment = new("some comment")
			},
			"CREATE STREAM IF NOT EXISTS %s ON TABLE %s AT (STREAM => '123') APPEND_ONLY = true SHOW_INITIAL_ROWS = true COMMENT = 'some comment'",
			id.FullyQualifiedName(), tableId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnTable_atTimestamp",
			func(opts *CreateOnTableStreamOptions) {
				opts.On = &OnStream{At: new(true), Statement: OnStreamStatement{Timestamp: new(timestamp)}}
			},
			"CREATE STREAM %s ON TABLE %s AT (TIMESTAMP => '%s')",
			id.FullyQualifiedName(), tableId.FullyQualifiedName(), timestamp,
		).
		withAdditionalSqlCasef(
			"sql_CreateOnTable_atOffset",
			func(opts *CreateOnTableStreamOptions) {
				opts.On = &OnStream{At: new(true), Statement: OnStreamStatement{Offset: new("-10")}}
			},
			"CREATE STREAM %s ON TABLE %s AT (OFFSET => -10)",
			id.FullyQualifiedName(), tableId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnTable_atStatement",
			func(opts *CreateOnTableStreamOptions) {
				opts.On = &OnStream{At: new(true), Statement: OnStreamStatement{Statement: new(queryId)}}
			},
			"CREATE STREAM %s ON TABLE %s AT (STATEMENT => '%s')",
			id.FullyQualifiedName(), tableId.FullyQualifiedName(), queryId,
		).
		withAdditionalSqlCasef(
			"sql_CreateOnTable_atStream",
			func(opts *CreateOnTableStreamOptions) {
				opts.On = &OnStream{At: new(true), Statement: OnStreamStatement{Stream: new(streamId.FullyQualifiedName())}}
			},
			"CREATE STREAM %s ON TABLE %s AT (STREAM => '%s')",
			id.FullyQualifiedName(), tableId.FullyQualifiedName(), temporaryReplace(streamId),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnTable_beforeTimestamp",
			func(opts *CreateOnTableStreamOptions) {
				opts.On = &OnStream{Before: new(true), Statement: OnStreamStatement{Timestamp: new(timestamp)}}
			},
			"CREATE STREAM %s ON TABLE %s BEFORE (TIMESTAMP => '%s')",
			id.FullyQualifiedName(), tableId.FullyQualifiedName(), timestamp,
		).
		withAdditionalSqlCasef(
			"sql_CreateOnTable_orReplace",
			func(opts *CreateOnTableStreamOptions) { opts.OrReplace = new(true) },
			"CREATE OR REPLACE STREAM %s ON TABLE %s", id.FullyQualifiedName(), tableId.FullyQualifiedName(),
		)

	streamsTests.CreateOnExternalTable.
		withDefaultOpts(func() *CreateOnExternalTableStreamOptions {
			return &CreateOnExternalTableStreamOptions{
				name:            id,
				ExternalTableId: externalTableId,
			}
		}).
		withExpectedSqlf(
			case_Streams_sql_CreateOnExternalTable_basic,
			"CREATE STREAM %s ON EXTERNAL TABLE %s", id.FullyQualifiedName(), externalTableId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_CreateOnExternalTable_all,
			func(opts *CreateOnExternalTableStreamOptions) {
				opts.IfNotExists = new(true)
				opts.On = &OnStream{At: new(true), Statement: OnStreamStatement{Statement: new("123")}}
				opts.InsertOnly = new(true)
				opts.Comment = new("some comment")
			},
			`CREATE STREAM IF NOT EXISTS %s ON EXTERNAL TABLE %s AT (STATEMENT => '123') INSERT_ONLY = true COMMENT = 'some comment'`,
			id.FullyQualifiedName(), externalTableId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnExternalTable_orReplace",
			func(opts *CreateOnExternalTableStreamOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			"CREATE OR REPLACE STREAM %s COPY GRANTS ON EXTERNAL TABLE %s",
			id.FullyQualifiedName(), externalTableId.FullyQualifiedName(),
		)

	streamsTests.CreateOnDirectoryTable.
		withDefaultOpts(func() *CreateOnDirectoryTableStreamOptions {
			return &CreateOnDirectoryTableStreamOptions{
				name:    id,
				StageId: stageId,
			}
		}).
		withExpectedSqlf(
			case_Streams_sql_CreateOnDirectoryTable_basic,
			"CREATE STREAM %s ON STAGE %s", id.FullyQualifiedName(), stageId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_CreateOnDirectoryTable_all,
			func(opts *CreateOnDirectoryTableStreamOptions) {
				opts.IfNotExists = new(true)
				opts.Comment = new("some comment")
			},
			`CREATE STREAM IF NOT EXISTS %s ON STAGE %s COMMENT = 'some comment'`,
			id.FullyQualifiedName(), stageId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnDirectoryTable_orReplace",
			func(opts *CreateOnDirectoryTableStreamOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			"CREATE OR REPLACE STREAM %s COPY GRANTS ON STAGE %s", id.FullyQualifiedName(), stageId.FullyQualifiedName(),
		)

	streamsTests.CreateOnView.
		withDefaultOpts(func() *CreateOnViewStreamOptions {
			return &CreateOnViewStreamOptions{
				name:   id,
				ViewId: viewId,
			}
		}).
		withExpectedSqlf(
			case_Streams_sql_CreateOnView_basic,
			"CREATE STREAM %s ON VIEW %s", id.FullyQualifiedName(), viewId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_CreateOnView_all,
			func(opts *CreateOnViewStreamOptions) {
				opts.IfNotExists = new(true)
				opts.On = &OnStream{Before: new(true), Statement: OnStreamStatement{Stream: new("123")}}
				opts.AppendOnly = new(true)
				opts.ShowInitialRows = new(true)
				opts.Comment = new("some comment")
			},
			`CREATE STREAM IF NOT EXISTS %s ON VIEW %s BEFORE (STREAM => '123') APPEND_ONLY = true SHOW_INITIAL_ROWS = true COMMENT = 'some comment'`,
			id.FullyQualifiedName(), viewId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateOnView_orReplace",
			func(opts *CreateOnViewStreamOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			"CREATE OR REPLACE STREAM %s COPY GRANTS ON VIEW %s", id.FullyQualifiedName(), viewId.FullyQualifiedName(),
		)

	streamsTests.Clone.
		withDefaultOpts(func() *CloneStreamOptions {
			return &CloneStreamOptions{
				name:         id,
				sourceStream: sourceStreamId,
			}
		}).
		withExpectedSqlf(
			case_Streams_sql_Clone_basic,
			"CREATE STREAM %s CLONE %s", id.FullyQualifiedName(), sourceStreamId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Clone_all",
			func(opts *CloneStreamOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			"CREATE OR REPLACE STREAM %s CLONE %s COPY GRANTS",
			id.FullyQualifiedName(), sourceStreamId.FullyQualifiedName(),
		)

	streamsTests.Alter.
		withModify(
			case_Streams_validation_Alter_opts_ConflictingFields,
			func(opts *AlterStreamOptions) {
				opts.IfExists = new(true)
				opts.UnsetTags = []ObjectIdentifier{tagId}
			},
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_Alter_SetComment,
			func(opts *AlterStreamOptions) {
				opts.IfExists = new(true)
				opts.SetComment = new("some comment")
			},
			`ALTER STREAM IF EXISTS %s SET COMMENT = 'some comment'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_Alter_UnsetComment,
			func(opts *AlterStreamOptions) {
				opts.IfExists = new(true)
				opts.UnsetComment = new(true)
			},
			`ALTER STREAM IF EXISTS %s UNSET COMMENT`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_Alter_SetTags,
			func(opts *AlterStreamOptions) {
				opts.IfExists = new(true)
				opts.SetTags = []TagAssociation{
					{Name: tagId, Value: "value1"},
					{Name: tagId2, Value: "value2"},
				}
			},
			`ALTER STREAM IF EXISTS %s SET TAG %s = 'value1', %s = 'value2'`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_Alter_UnsetTags,
			func(opts *AlterStreamOptions) {
				opts.UnsetTags = []ObjectIdentifier{tagId, tagId2}
			},
			`ALTER STREAM %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tagId.FullyQualifiedName(), tagId2.FullyQualifiedName(),
		)

	streamsTests.Drop.
		withExpectedSqlf(
			case_Streams_sql_Drop_basic,
			`DROP STREAM %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_Drop_all,
			func(opts *DropStreamOptions) { opts.IfExists = new(true) },
			`DROP STREAM IF EXISTS %s`, id.FullyQualifiedName(),
		)

	streamsTests.Show.
		withExpectedSql(case_Streams_sql_Show_basic, "SHOW STREAMS").
		withModifyAndExpectedSqlf(
			case_Streams_sql_Show_all,
			func(opts *ShowStreamOptions) {
				opts.Terse = new(true)
				opts.Like = &Like{Pattern: new("pattern")}
				opts.In = &ExtendedIn{In: In{Schema: schemaId}}
				opts.StartsWith = new("starts with pattern")
				opts.Limit = &LimitFrom{Rows: new(123), From: new("from pattern")}
			},
			`SHOW TERSE STREAMS LIKE 'pattern' IN SCHEMA %s STARTS WITH 'starts with pattern' LIMIT 123 FROM 'from pattern'`,
			schemaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_Show_Like,
			func(opts *ShowStreamOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			`SHOW STREAMS LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_Show_In,
			func(opts *ShowStreamOptions) { opts.In = &ExtendedIn{In: In{Schema: schemaId}} },
			`SHOW STREAMS IN SCHEMA %s`, schemaId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_Show_StartsWith,
			func(opts *ShowStreamOptions) { opts.StartsWith = new("starts with pattern") },
			`SHOW STREAMS STARTS WITH 'starts with pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_Streams_sql_Show_Limit,
			func(opts *ShowStreamOptions) { opts.Limit = &LimitFrom{Rows: new(123)} },
			`SHOW STREAMS LIMIT 123`,
		)

	streamsTests.Describe.
		withExpectedSqlf(
			case_Streams_sql_Describe_basic,
			`DESCRIBE STREAM %s`, id.FullyQualifiedName(),
		)
}
