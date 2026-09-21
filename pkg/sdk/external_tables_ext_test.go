package sdk

import "fmt"

func init() {
	id := externalTablesTestIdSchemaObjectIdentifier
	rowAccessPolicyId := randomSchemaObjectIdentifier()
	tag1 := NewAccountObjectIdentifier("tag1")
	tag2 := NewAccountObjectIdentifier("tag2")
	location := "@s1/logs/"
	databaseId := NewAccountObjectIdentifier("database_name")
	schemaId := randomDatabaseObjectIdentifier()

	externalTablesTests.Create.
		withDefaultOpts(func() *CreateExternalTableOptions {
			return &CreateExternalTableOptions{
				name:       id,
				Location:   location,
				FileFormat: &ExternalTableFileFormat{FileFormatType: new(ExternalTableFileFormatTypeJson)},
			}
		}).
		withModify(
			case_ExternalTables_validation_Create_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateExternalTableOptions) {
				opts.FileFormat.Name = new("JSON")
			},
		).
		withAdditionalValidationCase(
			"validation_Create_FileFormat_optionsMustMatchType",
			func(opts *CreateExternalTableOptions) {
				opts.FileFormat.Options = &ExternalTableFileFormatTypeOptions{
					CSVCompression: new(ExternalTableCsvCompressionGzip),
				}
			},
			fmt.Errorf("cannot set %s fields when TYPE = %s", ExternalTableFileFormatTypeCsv, ExternalTableFileFormatTypeJson),
		).
		withExpectedSqlf(
			case_ExternalTables_sql_Create_basic,
			`CREATE EXTERNAL TABLE %s LOCATION = %s FILE_FORMAT = (TYPE = JSON)`,
			id.FullyQualifiedName(), location,
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Create_all,
			func(opts *CreateExternalTableOptions) {
				opts.IfNotExists = new(true)
				opts.Columns = []ExternalTableColumn{externalTableTestColumn()}
				opts.CloudProviderParams = &CloudProviderParams{GoogleCloudStorageIntegration: new("123")}
				opts.PartitionBy = []string{"column"}
				opts.RefreshOnCreate = new(true)
				opts.AutoRefresh = new(true)
				opts.Pattern = new(".*")
				opts.FileFormat = &ExternalTableFileFormat{
					FileFormatType: new(ExternalTableFileFormatTypeCsv),
					Options: &ExternalTableFileFormatTypeOptions{
						CSVCompression: new(ExternalTableCsvCompressionGzip),
					},
				}
				opts.AwsSnsTopic = new("aws_sns_topic")
				opts.CopyGrants = new(true)
				opts.Comment = new("some_comment")
				opts.RowAccessPolicy = &TableRowAccessPolicyLegacy{Name: rowAccessPolicyId, On: []string{"value1", "value2"}}
				opts.Tag = []TagAssociation{
					{Name: tag1, Value: "value1"},
					{Name: tag2, Value: "value2"},
				}
			},
			`CREATE EXTERNAL TABLE IF NOT EXISTS %s (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) INTEGRATION = '123' PARTITION BY (column) LOCATION = %s REFRESH_ON_CREATE = true AUTO_REFRESH = true PATTERN = '.*' FILE_FORMAT = (TYPE = CSV COMPRESSION = GZIP) AWS_SNS_TOPIC = 'aws_sns_topic' COPY GRANTS COMMENT = 'some_comment' ROW ACCESS POLICY %s ON (value1, value2) TAG (%s = 'value1', %s = 'value2')`,
			id.FullyQualifiedName(), location, rowAccessPolicyId.FullyQualifiedName(), tag1.FullyQualifiedName(), tag2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateExternalTableOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE EXTERNAL TABLE %s LOCATION = %s FILE_FORMAT = (TYPE = JSON) COPY GRANTS`,
			id.FullyQualifiedName(), location,
		).
		withAdditionalSqlCasef(
			"sql_Create_rawFileFormat",
			func(opts *CreateExternalTableOptions) {
				opts.FileFormat = nil
				opts.RawFileFormat = &RawFileFormat{Format: "TYPE = JSON"}
			},
			`CREATE EXTERNAL TABLE %s LOCATION = %s FILE_FORMAT = (TYPE = JSON)`,
			id.FullyQualifiedName(), location,
		)

	externalTablesTests.CreateWithManualPartitioning.
		withDefaultOpts(func() *CreateWithManualPartitioningExternalTableOptions {
			return &CreateWithManualPartitioningExternalTableOptions{
				name:       id,
				Location:   location,
				FileFormat: &ExternalTableFileFormat{FileFormatType: new(ExternalTableFileFormatTypeJson)},
			}
		}).
		withModify(
			case_ExternalTables_validation_CreateWithManualPartitioning_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateWithManualPartitioningExternalTableOptions) {
				opts.FileFormat.Name = new("JSON")
			},
		).
		withExpectedSqlf(
			case_ExternalTables_sql_CreateWithManualPartitioning_basic,
			`CREATE EXTERNAL TABLE %s LOCATION = %s PARTITION_TYPE = USER_SPECIFIED FILE_FORMAT = (TYPE = JSON)`,
			id.FullyQualifiedName(), location,
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_CreateWithManualPartitioning_all,
			func(opts *CreateWithManualPartitioningExternalTableOptions) {
				opts.IfNotExists = new(true)
				opts.Columns = []ExternalTableColumn{externalTableTestColumn()}
				opts.CloudProviderParams = &CloudProviderParams{GoogleCloudStorageIntegration: new("123")}
				opts.PartitionBy = []string{"column"}
				opts.CopyGrants = new(true)
				opts.Comment = new("some_comment")
				opts.RowAccessPolicy = &TableRowAccessPolicyLegacy{Name: rowAccessPolicyId, On: []string{"value1", "value2"}}
				opts.Tag = []TagAssociation{
					{Name: tag1, Value: "value1"},
					{Name: tag2, Value: "value2"},
				}
			},
			`CREATE EXTERNAL TABLE IF NOT EXISTS %s (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) INTEGRATION = '123' PARTITION BY (column) LOCATION = %s PARTITION_TYPE = USER_SPECIFIED FILE_FORMAT = (TYPE = JSON) COPY GRANTS COMMENT = 'some_comment' ROW ACCESS POLICY %s ON (value1, value2) TAG (%s = 'value1', %s = 'value2')`,
			id.FullyQualifiedName(), location, rowAccessPolicyId.FullyQualifiedName(), tag1.FullyQualifiedName(), tag2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateWithManualPartitioning_orReplace",
			func(opts *CreateWithManualPartitioningExternalTableOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE EXTERNAL TABLE %s LOCATION = %s PARTITION_TYPE = USER_SPECIFIED FILE_FORMAT = (TYPE = JSON) COPY GRANTS`,
			id.FullyQualifiedName(), location,
		).
		withAdditionalSqlCasef(
			"sql_CreateWithManualPartitioning_rawFileFormat",
			func(opts *CreateWithManualPartitioningExternalTableOptions) {
				opts.FileFormat = nil
				opts.RawFileFormat = &RawFileFormat{Format: "TYPE = JSON"}
			},
			`CREATE EXTERNAL TABLE %s LOCATION = %s PARTITION_TYPE = USER_SPECIFIED FILE_FORMAT = (TYPE = JSON)`,
			id.FullyQualifiedName(), location,
		)

	externalTablesTests.CreateDeltaLake.
		withDefaultOpts(func() *CreateDeltaLakeExternalTableOptions {
			return &CreateDeltaLakeExternalTableOptions{
				name:       id,
				Location:   location,
				FileFormat: &ExternalTableFileFormat{FileFormatType: new(ExternalTableFileFormatTypeJson)},
			}
		}).
		withModify(
			case_ExternalTables_validation_CreateDeltaLake_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateDeltaLakeExternalTableOptions) {
				opts.FileFormat.Name = new("JSON")
			},
		).
		withExpectedSqlf(
			case_ExternalTables_sql_CreateDeltaLake_basic,
			`CREATE EXTERNAL TABLE %s LOCATION = %s FILE_FORMAT = (TYPE = JSON) TABLE_FORMAT = DELTA`,
			id.FullyQualifiedName(), location,
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_CreateDeltaLake_all,
			func(opts *CreateDeltaLakeExternalTableOptions) {
				opts.IfNotExists = new(true)
				opts.Columns = []ExternalTableColumn{{
					Name:         "column",
					DataType:     "varchar",
					AsExpression: AsExpression{Expression: "value::column::varchar"},
				}}
				opts.CloudProviderParams = &CloudProviderParams{MicrosoftAzureIntegration: new("123")}
				opts.PartitionBy = []string{"column"}
				opts.RefreshOnCreate = new(true)
				opts.AutoRefresh = new(true)
				opts.FileFormat = &ExternalTableFileFormat{Name: new("JSON")}
				opts.CopyGrants = new(true)
				opts.Comment = new("some_comment")
				opts.RowAccessPolicy = &TableRowAccessPolicyLegacy{Name: rowAccessPolicyId, On: []string{"value1", "value2"}}
				opts.Tag = []TagAssociation{
					{Name: tag1, Value: "value1"},
					{Name: tag2, Value: "value2"},
				}
			},
			`CREATE EXTERNAL TABLE IF NOT EXISTS %s (column varchar AS (value::column::varchar)) INTEGRATION = '123' PARTITION BY (column) LOCATION = %s REFRESH_ON_CREATE = true AUTO_REFRESH = true FILE_FORMAT = (FORMAT_NAME = 'JSON') TABLE_FORMAT = DELTA COPY GRANTS COMMENT = 'some_comment' ROW ACCESS POLICY %s ON (value1, value2) TAG (%s = 'value1', %s = 'value2')`,
			id.FullyQualifiedName(), location, rowAccessPolicyId.FullyQualifiedName(), tag1.FullyQualifiedName(), tag2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateDeltaLake_orReplace",
			func(opts *CreateDeltaLakeExternalTableOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
			},
			`CREATE OR REPLACE EXTERNAL TABLE %s LOCATION = %s FILE_FORMAT = (TYPE = JSON) TABLE_FORMAT = DELTA COPY GRANTS`,
			id.FullyQualifiedName(), location,
		).
		withAdditionalSqlCasef(
			"sql_CreateDeltaLake_rawFileFormat",
			func(opts *CreateDeltaLakeExternalTableOptions) {
				opts.FileFormat = nil
				opts.RawFileFormat = &RawFileFormat{Format: "TYPE = JSON"}
			},
			`CREATE EXTERNAL TABLE %s LOCATION = %s FILE_FORMAT = (TYPE = JSON) TABLE_FORMAT = DELTA`,
			id.FullyQualifiedName(), location,
		)

	externalTablesTests.CreateUsingTemplate.
		withDefaultOpts(func() *CreateUsingTemplateExternalTableOptions {
			return &CreateUsingTemplateExternalTableOptions{
				name:       id,
				Query:      []string{"query statement"},
				Location:   location,
				FileFormat: &ExternalTableFileFormat{FileFormatType: new(ExternalTableFileFormatTypeJson)},
			}
		}).
		withModify(
			case_ExternalTables_validation_CreateUsingTemplate_opts_FileFormat_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateUsingTemplateExternalTableOptions) {
				opts.FileFormat.Name = new("JSON")
			},
		).
		withExpectedSqlf(
			case_ExternalTables_sql_CreateUsingTemplate_basic,
			`CREATE EXTERNAL TABLE %s USING TEMPLATE (query statement) LOCATION = %s FILE_FORMAT = (TYPE = JSON)`,
			id.FullyQualifiedName(), location,
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_CreateUsingTemplate_all,
			func(opts *CreateUsingTemplateExternalTableOptions) {
				opts.OrReplace = new(true)
				opts.CopyGrants = new(true)
				opts.CloudProviderParams = &CloudProviderParams{MicrosoftAzureIntegration: new("123")}
				opts.PartitionBy = []string{"column"}
				opts.RefreshOnCreate = new(true)
				opts.AutoRefresh = new(true)
				opts.Pattern = new(".*")
				opts.FileFormat = &ExternalTableFileFormat{Name: new("JSON")}
				opts.AwsSnsTopic = new("aws_sns_topic")
				opts.Comment = new("some_comment")
				opts.RowAccessPolicy = &TableRowAccessPolicyLegacy{Name: rowAccessPolicyId, On: []string{"value1", "value2"}}
				opts.Tag = []TagAssociation{
					{Name: tag1, Value: "value1"},
					{Name: tag2, Value: "value2"},
				}
			},
			`CREATE OR REPLACE EXTERNAL TABLE %s COPY GRANTS USING TEMPLATE (query statement) INTEGRATION = '123' PARTITION BY (column) LOCATION = %s REFRESH_ON_CREATE = true AUTO_REFRESH = true PATTERN = '.*' FILE_FORMAT = (FORMAT_NAME = 'JSON') AWS_SNS_TOPIC = 'aws_sns_topic' COMMENT = 'some_comment' ROW ACCESS POLICY %s ON (value1, value2) TAG (%s = 'value1', %s = 'value2')`,
			id.FullyQualifiedName(), location, rowAccessPolicyId.FullyQualifiedName(), tag1.FullyQualifiedName(), tag2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateUsingTemplate_rawFileFormat",
			func(opts *CreateUsingTemplateExternalTableOptions) {
				opts.FileFormat = nil
				opts.RawFileFormat = &RawFileFormat{Format: "TYPE = JSON"}
			},
			`CREATE EXTERNAL TABLE %s USING TEMPLATE (query statement) LOCATION = %s FILE_FORMAT = (TYPE = JSON)`,
			id.FullyQualifiedName(), location,
		)

	externalTablesTests.Alter.
		withModify(
			case_ExternalTables_validation_Alter_opts_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *AlterExternalTableOptions) {
				opts.Refresh = &RefreshExternalTable{}
				opts.AddFiles = []ExternalTableFile{{Name: "one/file.txt"}}
			},
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Alter_Refresh,
			func(opts *AlterExternalTableOptions) {
				opts.IfExists = new(true)
				opts.Refresh = &RefreshExternalTable{}
			},
			`ALTER EXTERNAL TABLE IF EXISTS %s REFRESH ''`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Refresh_withPath",
			func(opts *AlterExternalTableOptions) {
				opts.IfExists = new(true)
				opts.Refresh = &RefreshExternalTable{Path: "some/path"}
			},
			`ALTER EXTERNAL TABLE IF EXISTS %s REFRESH 'some/path'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Alter_AddFiles,
			func(opts *AlterExternalTableOptions) {
				opts.AddFiles = []ExternalTableFile{{Name: "one/file.txt"}, {Name: "second/file.txt"}}
			},
			`ALTER EXTERNAL TABLE %s ADD FILES ('one/file.txt', 'second/file.txt')`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Alter_RemoveFiles,
			func(opts *AlterExternalTableOptions) {
				opts.RemoveFiles = []ExternalTableFile{{Name: "one/file.txt"}, {Name: "second/file.txt"}}
			},
			`ALTER EXTERNAL TABLE %s REMOVE FILES ('one/file.txt', 'second/file.txt')`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Alter_AutoRefresh,
			func(opts *AlterExternalTableOptions) { opts.AutoRefresh = new(true) },
			`ALTER EXTERNAL TABLE %s SET AUTO_REFRESH = true`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Alter_SetTags,
			func(opts *AlterExternalTableOptions) {
				opts.SetTags = []TagAssociation{
					{Name: tag1, Value: "tag_value1"},
					{Name: tag2, Value: "tag_value2"},
				}
			},
			`ALTER EXTERNAL TABLE %s SET TAG %s = 'tag_value1', %s = 'tag_value2'`,
			id.FullyQualifiedName(), tag1.FullyQualifiedName(), tag2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Alter_UnsetTags,
			func(opts *AlterExternalTableOptions) {
				opts.UnsetTags = []ObjectIdentifier{tag1, tag2}
			},
			`ALTER EXTERNAL TABLE %s UNSET TAG %s, %s`,
			id.FullyQualifiedName(), tag1.FullyQualifiedName(), tag2.FullyQualifiedName(),
		)

	externalTablesTests.AlterPartitions.
		withDefaultOpts(func() *AlterPartitionsExternalTableOptions {
			return &AlterPartitionsExternalTableOptions{
				name:     id,
				Location: "123",
			}
		}).
		withModify(
			case_ExternalTables_validation_AlterPartitions_opts_ConflictingFields,
			func(opts *AlterPartitionsExternalTableOptions) {
				opts.AddPartitions = []Partition{{ColumnName: "colName", Value: "value"}}
				opts.DropPartition = new(true)
			},
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_AlterPartitions_basic,
			func(opts *AlterPartitionsExternalTableOptions) {
				opts.IfExists = new(true)
				opts.DropPartition = new(true)
			},
			`ALTER EXTERNAL TABLE IF EXISTS %s DROP PARTITION LOCATION '123'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_AlterPartitions_all,
			func(opts *AlterPartitionsExternalTableOptions) {
				opts.IfExists = new(true)
				opts.AddPartitions = []Partition{
					{ColumnName: "one", Value: "one_value"},
					{ColumnName: "two", Value: "two_value"},
				}
			},
			`ALTER EXTERNAL TABLE IF EXISTS %s ADD PARTITION (one = 'one_value', two = 'two_value') LOCATION '123'`,
			id.FullyQualifiedName(),
		)

	externalTablesTests.Drop.
		withExpectedSqlf(
			case_ExternalTables_sql_Drop_basic,
			`DROP EXTERNAL TABLE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Drop_all,
			func(opts *DropExternalTableOptions) {
				opts.IfExists = new(true)
				opts.DropOption = &ExternalTableDropOption{Restrict: new(true)}
			},
			`DROP EXTERNAL TABLE IF EXISTS %s RESTRICT`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Drop_cascade",
			func(opts *DropExternalTableOptions) {
				opts.IfExists = new(true)
				opts.DropOption = &ExternalTableDropOption{Cascade: new(true)}
			},
			`DROP EXTERNAL TABLE IF EXISTS %s CASCADE`, id.FullyQualifiedName(),
		)

	likePattern := "some_pattern"

	externalTablesTests.Show.
		withExpectedSql(case_ExternalTables_sql_Show_basic, `SHOW EXTERNAL TABLES`).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Show_all,
			func(opts *ShowExternalTableOptions) {
				opts.Terse = new(true)
				opts.Like = &Like{Pattern: &likePattern}
				opts.In = &In{Account: new(true)}
				opts.StartsWith = new("some_external_table")
				opts.LimitFrom = &LimitFrom{Rows: new(123), From: new("some_string")}
			},
			`SHOW TERSE EXTERNAL TABLES LIKE '%s' IN ACCOUNT STARTS WITH 'some_external_table' LIMIT 123 FROM 'some_string'`,
			likePattern,
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Show_Like,
			func(opts *ShowExternalTableOptions) { opts.Like = &Like{Pattern: &likePattern} },
			`SHOW EXTERNAL TABLES LIKE '%s'`, likePattern,
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Show_In,
			func(opts *ShowExternalTableOptions) { opts.In = &In{Database: databaseId} },
			`SHOW EXTERNAL TABLES IN DATABASE %s`, databaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Show_StartsWith,
			func(opts *ShowExternalTableOptions) { opts.StartsWith = new("some_external_table") },
			`SHOW EXTERNAL TABLES STARTS WITH 'some_external_table'`,
		).
		withModifyAndExpectedSqlf(
			case_ExternalTables_sql_Show_LimitFrom,
			func(opts *ShowExternalTableOptions) { opts.LimitFrom = &LimitFrom{Rows: new(123)} },
			`SHOW EXTERNAL TABLES LIMIT 123`,
		).
		withAdditionalSqlCasef(
			"sql_Show_inSchema",
			func(opts *ShowExternalTableOptions) {
				opts.Terse = new(true)
				opts.In = &In{Schema: schemaId}
			},
			`SHOW TERSE EXTERNAL TABLES IN SCHEMA %s`, schemaId.FullyQualifiedName(),
		)

	externalTablesTests.DescribeColumns.
		withExpectedSqlf(
			case_ExternalTables_sql_DescribeColumns_basic,
			`DESCRIBE EXTERNAL TABLE %s TYPE = COLUMNS`, id.FullyQualifiedName(),
		)

	externalTablesTests.DescribeStage.
		withExpectedSqlf(
			case_ExternalTables_sql_DescribeStage_basic,
			`DESCRIBE EXTERNAL TABLE %s TYPE = STAGE`, id.FullyQualifiedName(),
		)
}

func externalTableTestColumn() ExternalTableColumn {
	return ExternalTableColumn{
		Name:         "column",
		DataType:     "varchar",
		AsExpression: AsExpression{Expression: "value::column::varchar"},
		NotNull:      new(true),
		InlineConstraint: &ColumnInlineConstraint{
			Name: new("my_constraint"),
			Type: ColumnConstraintTypeUnique,
		},
	}
}
