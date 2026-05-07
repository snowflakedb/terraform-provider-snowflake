package sdk

import (
	"testing"
)

func TestExternalTablesCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("basic options", func(t *testing.T) {
		opts := &CreateExternalTableOptions{
			IfNotExists: Bool(true),
			name:        id,
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			CloudProviderParams: &CloudProviderParams{
				GoogleCloudStorageIntegration: String("123"),
			},
			Location: "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Type: &ExternalTableFileFormatTypeJson,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL TABLE IF NOT EXISTS %s (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) INTEGRATION = '123' LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON)`, id.FullyQualifiedName())
	})

	t.Run("every optional field", func(t *testing.T) {
		rowAccessPolicyId := randomSchemaObjectIdentifier()
		opts := &CreateExternalTableOptions{
			OrReplace: Bool(true),
			name:      id,
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			CloudProviderParams: &CloudProviderParams{
				GoogleCloudStorageIntegration: String("123"),
			},
			Location: "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Type: &ExternalTableFileFormatTypeJson,
				},
			},
			AwsSnsTopic: String("aws_sns_topic"),
			CopyGrants:  Bool(true),
			RowAccessPolicy: &TableRowAccessPolicy{
				Name: rowAccessPolicyId,
				On:   []string{"value1", "value2"},
			},
			Tag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag1"),
					Value: "value1",
				},
				{
					Name:  NewAccountObjectIdentifier("tag2"),
					Value: "value2",
				},
			},
			Comment: String("some_comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EXTERNAL TABLE %s (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) INTEGRATION = '123' LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON) AWS_SNS_TOPIC = 'aws_sns_topic' COPY GRANTS COMMENT = 'some_comment' ROW ACCESS POLICY %s ON (value1, value2) TAG ("tag1" = 'value1', "tag2" = 'value2')`, id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &CreateExternalTableOptions{
			OrReplace:   Bool(true),
			IfNotExists: Bool(true),
			name:        emptySchemaObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			errOneOf("CreateExternalTableOptions", "OrReplace", "IfNotExists"),
			ErrInvalidObjectIdentifier,
			errNotSet("CreateExternalTableOptions", "Location"),
			errExactlyOneOf("CreateExternalTableOptions", "RawFileFormat", "FileFormat"),
		)
	})

	t.Run("raw file format", func(t *testing.T) {
		opts := &CreateExternalTableOptions{
			name: id,
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location:      "@s1/logs/",
			RawFileFormat: &RawFileFormat{Format: "TYPE = JSON"},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL TABLE %s (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON)`, id.FullyQualifiedName())
	})

	t.Run("validation: neither raw file format is set, nor file format", func(t *testing.T) {
		opts := &CreateExternalTableOptions{
			name: id,
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location: "@s1/logs/",
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("CreateExternalTableOptions", "RawFileFormat", "FileFormat"))
	})
}

func TestExternalTablesCreateWithManualPartitioning(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("valid options", func(t *testing.T) {
		rowAccessPolicyId := randomSchemaObjectIdentifier()
		opts := &CreateWithManualPartitioningExternalTableOptions{
			OrReplace: Bool(true),
			name:      id,
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			CloudProviderParams: &CloudProviderParams{
				GoogleCloudStorageIntegration: String("123"),
			},
			Location: "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Type: &ExternalTableFileFormatTypeJson,
				},
			},
			CopyGrants: Bool(true),
			RowAccessPolicy: &TableRowAccessPolicy{
				Name: rowAccessPolicyId,
				On:   []string{"value1", "value2"},
			},
			Tag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag1"),
					Value: "value1",
				},
				{
					Name:  NewAccountObjectIdentifier("tag2"),
					Value: "value2",
				},
			},
			Comment: String("some_comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EXTERNAL TABLE %s (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) INTEGRATION = '123' LOCATION = @s1/logs/ PARTITION_TYPE = USER_SPECIFIED FILE_FORMAT = (TYPE = JSON) COPY GRANTS COMMENT = 'some_comment' ROW ACCESS POLICY %s ON (value1, value2) TAG ("tag1" = 'value1', "tag2" = 'value2')`, id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &CreateWithManualPartitioningExternalTableOptions{
			OrReplace:   Bool(true),
			IfNotExists: Bool(true),
			name:        emptySchemaObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			errOneOf("CreateWithManualPartitioningExternalTableOptions", "OrReplace", "IfNotExists"),
			ErrInvalidObjectIdentifier,
			errNotSet("CreateWithManualPartitioningExternalTableOptions", "Location"),
			errExactlyOneOf("CreateWithManualPartitioningExternalTableOptions", "RawFileFormat", "FileFormat"),
		)
	})

	t.Run("raw file format", func(t *testing.T) {
		opts := &CreateWithManualPartitioningExternalTableOptions{
			name: id,
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location:      "@s1/logs/",
			RawFileFormat: &RawFileFormat{Format: "TYPE = JSON"},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL TABLE %s (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) LOCATION = @s1/logs/ PARTITION_TYPE = USER_SPECIFIED FILE_FORMAT = (TYPE = JSON)`, id.FullyQualifiedName())
	})

	t.Run("validation: neither raw file format is set, nor file format", func(t *testing.T) {
		opts := &CreateWithManualPartitioningExternalTableOptions{
			name: id,
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location: "@s1/logs/",
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("CreateWithManualPartitioningExternalTableOptions", "RawFileFormat", "FileFormat"))
	})
}

func TestExternalTablesCreateDeltaLake(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("valid options", func(t *testing.T) {
		rowAccessPolicyId := randomSchemaObjectIdentifier()
		opts := &CreateDeltaLakeExternalTableOptions{
			OrReplace: Bool(true),
			name:      id,
			Columns: []ExternalTableColumn{
				{
					Name:             "column",
					Type:             "varchar",
					AsExpression:     []string{"value::column::varchar"},
					InlineConstraint: nil,
				},
			},
			CloudProviderParams: &CloudProviderParams{
				MicrosoftAzureIntegration: String("123"),
			},
			PartitionBy: []string{"column"},
			Location:    "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Name: String("JSON"),
				},
			},
			CopyGrants: Bool(true),
			RowAccessPolicy: &TableRowAccessPolicy{
				Name: rowAccessPolicyId,
				On:   []string{"value1", "value2"},
			},
			Tag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag1"),
					Value: "value1",
				},
				{
					Name:  NewAccountObjectIdentifier("tag2"),
					Value: "value2",
				},
			},
			Comment: String("some_comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EXTERNAL TABLE %s (column varchar AS (value::column::varchar)) INTEGRATION = '123' PARTITION BY (column) LOCATION = @s1/logs/ FILE_FORMAT = (FORMAT_NAME = 'JSON') TABLE_FORMAT = DELTA COPY GRANTS COMMENT = 'some_comment' ROW ACCESS POLICY %s ON (value1, value2) TAG ("tag1" = 'value1', "tag2" = 'value2')`, id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &CreateDeltaLakeExternalTableOptions{
			OrReplace:   Bool(true),
			IfNotExists: Bool(true),
			name:        emptySchemaObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			errOneOf("CreateDeltaLakeExternalTableOptions", "OrReplace", "IfNotExists"),
			ErrInvalidObjectIdentifier,
			errNotSet("CreateDeltaLakeExternalTableOptions", "Location"),
			errExactlyOneOf("CreateDeltaLakeExternalTableOptions", "RawFileFormat", "FileFormat"),
		)
	})

	t.Run("raw file format", func(t *testing.T) {
		opts := &CreateDeltaLakeExternalTableOptions{
			name: id,
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location:      "@s1/logs/",
			RawFileFormat: &RawFileFormat{Format: "TYPE = JSON"},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL TABLE %s (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON) TABLE_FORMAT = DELTA`, id.FullyQualifiedName())
	})

	t.Run("validation: neither raw file format is set, nor file format", func(t *testing.T) {
		opts := &CreateDeltaLakeExternalTableOptions{
			name: id,
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location: "@s1/logs/",
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("CreateDeltaLakeExternalTableOptions", "RawFileFormat", "FileFormat"))
	})
}

func TestExternalTableUsingTemplateOpts(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("valid options", func(t *testing.T) {
		rowAccessPolicyId := randomSchemaObjectIdentifier()
		opts := &CreateExternalTableUsingTemplateOptions{
			OrReplace:  Bool(true),
			name:       id,
			CopyGrants: Bool(true),
			Query:      []string{"query statement"},
			CloudProviderParams: &CloudProviderParams{
				MicrosoftAzureIntegration: String("123"),
			},
			PartitionBy: []string{"column"},
			Location:    "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Name: String("JSON"),
				},
			},
			Comment: String("some_comment"),
			RowAccessPolicy: &TableRowAccessPolicy{
				Name: rowAccessPolicyId,
				On:   []string{"value1", "value2"},
			},
			Tag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag1"),
					Value: "value1",
				},
				{
					Name:  NewAccountObjectIdentifier("tag2"),
					Value: "value2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EXTERNAL TABLE %s COPY GRANTS USING TEMPLATE (query statement) INTEGRATION = '123' PARTITION BY (column) LOCATION = @s1/logs/ FILE_FORMAT = (FORMAT_NAME = 'JSON') COMMENT = 'some_comment' ROW ACCESS POLICY %s ON (value1, value2) TAG ("tag1" = 'value1', "tag2" = 'value2')`, id.FullyQualifiedName(), rowAccessPolicyId.FullyQualifiedName())
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &CreateExternalTableUsingTemplateOptions{
			name: emptySchemaObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			errNotSet("CreateExternalTableUsingTemplateOptions", "Query"),
			errNotSet("CreateExternalTableUsingTemplateOptions", "Location"),
			errExactlyOneOf("CreateExternalTableUsingTemplateOptions", "RawFileFormat", "FileFormat"),
		)
	})

	t.Run("raw file format", func(t *testing.T) {
		opts := &CreateExternalTableUsingTemplateOptions{
			name:     id,
			Location: "@s1/logs/",
			Query: []string{
				"query statement",
			},
			RawFileFormat: &RawFileFormat{Format: "TYPE = JSON"},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL TABLE %s USING TEMPLATE (query statement) LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON)`, id.FullyQualifiedName())
	})

	t.Run("validation: neither raw file format is set, nor file format", func(t *testing.T) {
		opts := &CreateExternalTableUsingTemplateOptions{
			name:     id,
			Location: "@s1/logs/",
			Query: []string{
				"query statement",
			},
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("CreateExternalTableUsingTemplateOptions", "RawFileFormat", "FileFormat"))
	})
}

func TestExternalTablesAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("refresh without path", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			IfExists: Bool(true),
			name:     id,
			Refresh:  &RefreshExternalTable{},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE IF EXISTS %s REFRESH ''`, id.FullyQualifiedName())
	})

	t.Run("refresh with path", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			IfExists: Bool(true),
			name:     id,
			Refresh: &RefreshExternalTable{
				Path: "some/path",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE IF EXISTS %s REFRESH 'some/path'`, id.FullyQualifiedName())
	})

	t.Run("add files", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: id,
			AddFiles: []ExternalTableFile{
				{Name: "one/file.txt"},
				{Name: "second/file.txt"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE %s ADD FILES ('one/file.txt', 'second/file.txt')`, id.FullyQualifiedName())
	})

	t.Run("remove files", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: id,
			RemoveFiles: []ExternalTableFile{
				{Name: "one/file.txt"},
				{Name: "second/file.txt"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE %s REMOVE FILES ('one/file.txt', 'second/file.txt')`, id.FullyQualifiedName())
	})

	t.Run("set auto refresh", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name:        id,
			AutoRefresh: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE %s SET AUTO_REFRESH = true`, id.FullyQualifiedName())
	})

	t.Run("set tag", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: id,
			SetTag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag1"),
					Value: "tag_value1",
				},
				{
					Name:  NewAccountObjectIdentifier("tag2"),
					Value: "tag_value2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE %s SET TAG "tag1" = 'tag_value1', "tag2" = 'tag_value2'`, id.FullyQualifiedName())
	})

	t.Run("unset tag", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: id,
			UnsetTag: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag1"),
				NewAccountObjectIdentifier("tag2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE %s UNSET TAG "tag1", "tag2"`, id.FullyQualifiedName())
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name:        emptySchemaObjectIdentifier,
			AddFiles:    []ExternalTableFile{{Name: "some file"}},
			RemoveFiles: []ExternalTableFile{{Name: "some other file"}},
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			errExactlyOneOf("AlterExternalTableOptions", "Refresh", "AddFiles", "RemoveFiles", "AutoRefresh", "SetTag", "UnsetTag"),
		)
	})
}

func TestExternalTablesAlterPartitions(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("add partition", func(t *testing.T) {
		opts := &AlterExternalTablePartitionOptions{
			name:     id,
			IfExists: Bool(true),
			AddPartitions: []Partition{
				{
					ColumnName: "one",
					Value:      "one_value",
				},
				{
					ColumnName: "two",
					Value:      "two_value",
				},
			},
			Location: "123",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE IF EXISTS %s ADD PARTITION (one = 'one_value', two = 'two_value') LOCATION '123'`, id.FullyQualifiedName())
	})

	t.Run("remove partition", func(t *testing.T) {
		opts := &AlterExternalTablePartitionOptions{
			name:          id,
			IfExists:      Bool(true),
			DropPartition: Bool(true),
			Location:      "partition_location",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE IF EXISTS %s DROP PARTITION LOCATION 'partition_location'`, id.FullyQualifiedName())
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &AlterExternalTablePartitionOptions{
			name:          emptySchemaObjectIdentifier,
			AddPartitions: []Partition{{ColumnName: "colName", Value: "value"}},
			DropPartition: Bool(true),
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			errOneOf("AlterExternalTablePartitionOptions", "AddPartitions", "DropPartition"),
		)
	})
}

func TestExternalTablesDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("restrict", func(t *testing.T) {
		opts := &DropExternalTableOptions{
			IfExists: Bool(true),
			name:     id,
			DropOption: &ExternalTableDropOption{
				Restrict: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP EXTERNAL TABLE IF EXISTS %s RESTRICT`, id.FullyQualifiedName())
	})

	t.Run("cascade", func(t *testing.T) {
		opts := &DropExternalTableOptions{
			IfExists: Bool(true),
			name:     id,
			DropOption: &ExternalTableDropOption{
				Cascade: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP EXTERNAL TABLE IF EXISTS %s CASCADE`, id.FullyQualifiedName())
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &DropExternalTableOptions{
			name: emptySchemaObjectIdentifier,
			DropOption: &ExternalTableDropOption{
				Restrict: Bool(true),
				Cascade:  Bool(true),
			},
		}

		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			errOneOf("ExternalTableDropOption", "Restrict", "Cascade"),
		)
	})
}

func TestExternalTablesShow(t *testing.T) {
	t.Run("all options", func(t *testing.T) {
		opts := &ShowExternalTableOptions{
			Terse: Bool(true),
			Like: &Like{
				Pattern: String("some_pattern"),
			},
			In: &In{
				Account: Bool(true),
			},
			StartsWith: String("some_external_table"),
			LimitFrom: &LimitFrom{
				Rows: Int(123),
				From: String("some_string"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE EXTERNAL TABLES LIKE 'some_pattern' IN ACCOUNT STARTS WITH 'some_external_table' LIMIT 123 FROM 'some_string'`)
	})

	t.Run("in database", func(t *testing.T) {
		opts := &ShowExternalTableOptions{
			Terse: Bool(true),
			In: &In{
				Database: NewAccountObjectIdentifier("database_name"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE EXTERNAL TABLES IN DATABASE "database_name"`)
	})

	t.Run("in schema", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		opts := &ShowExternalTableOptions{
			Terse: Bool(true),
			In: &In{
				Schema: id,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE EXTERNAL TABLES IN SCHEMA %s`, id.FullyQualifiedName())
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &DropExternalTableOptions{
			name: emptySchemaObjectIdentifier,
			DropOption: &ExternalTableDropOption{
				Restrict: Bool(true),
				Cascade:  Bool(true),
			},
		}

		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			errOneOf("ExternalTableDropOption", "Restrict", "Cascade"),
		)
	})
}

func TestExternalTablesDescribe(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("type columns", func(t *testing.T) {
		opts := &describeExternalTableColumnsOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE EXTERNAL TABLE %s TYPE = COLUMNS`, id.FullyQualifiedName())
	})

	t.Run("type stage", func(t *testing.T) {
		opts := &describeExternalTableStageOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE EXTERNAL TABLE %s TYPE = STAGE`, id.FullyQualifiedName())
	})
}
