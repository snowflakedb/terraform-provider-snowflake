package sdk

import (
	"testing"
	"time"
)

func TestSchemasCreate(t *testing.T) {
	id := randomDatabaseObjectIdentifier()

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := &CreateSchemaOptions{
			name: emptyDatabaseObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateSchemaOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := &CreateSchemaOptions{
			name:        id,
			OrReplace:   new(true),
			IfNotExists: new(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSchemaOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := &CreateSchemaOptions{
			name:           id,
			ExternalVolume: new(emptyAccountObjectIdentifier),
			Catalog:        new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts,
			errInvalidIdentifier("CreateSchemaOptions", "ExternalVolume"),
			errInvalidIdentifier("CreateSchemaOptions", "Catalog"),
		)
	})

	t.Run("clone", func(t *testing.T) {
		opts := &CreateSchemaOptions{
			name:      id,
			OrReplace: new(true),
			Clone: &Clone{
				SourceObject: NewAccountObjectIdentifier("sch1"),
				At: &TimeTravel{
					Timestamp: new(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SCHEMA %s CLONE "sch1" AT (TIMESTAMP => '2021-01-01 00:00:00 +0000 UTC')`, id.FullyQualifiedName())
	})

	t.Run("complete", func(t *testing.T) {
		tagId := randomSchemaObjectIdentifier()
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts := &CreateSchemaOptions{
			Transient:                               new(true),
			IfNotExists:                             new(true),
			name:                                    id,
			WithManagedAccess:                       new(true),
			DataRetentionTimeInDays:                 new(1),
			MaxDataExtensionTimeInDays:              new(1),
			ExternalVolume:                          &externalVolumeId,
			Catalog:                                 &catalogId,
			PipeExecutionPaused:                     new(true),
			ReplaceInvalidCharacters:                new(true),
			DefaultDDLCollation:                     new(StringAllowEmpty{Value: "en_US-trim"}),
			StorageSerializationPolicy:              new(StorageSerializationPolicyCompatible),
			LogLevel:                                new(LogLevelInfo),
			TraceLevel:                              new(TraceLevelPropagate),
			SuspendTaskAfterNumFailures:             new(10),
			TaskAutoRetryAttempts:                   new(10),
			UserTaskManagedInitialWarehouseSize:     new(WarehouseSizeMedium),
			UserTaskTimeoutMs:                       new(12000),
			UserTaskMinimumTriggerIntervalInSeconds: new(30),
			QuotedIdentifiersIgnoreCase:             new(true),
			EnableConsoleOutput:                     new(true),
			Comment:                                 new("comment"),
			Tag:                                     []TagAssociation{{Name: tagId, Value: "v1"}},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE TRANSIENT SCHEMA IF NOT EXISTS %s WITH MANAGED ACCESS DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 `+
			`EXTERNAL_VOLUME = "%s" CATALOG = "%s" PIPE_EXECUTION_PAUSED = true REPLACE_INVALID_CHARACTERS = true DEFAULT_DDL_COLLATION = 'en_US-trim' STORAGE_SERIALIZATION_POLICY = COMPATIBLE `+
			`LOG_LEVEL = 'INFO' TRACE_LEVEL = 'PROPAGATE' SUSPEND_TASK_AFTER_NUM_FAILURES = 10 TASK_AUTO_RETRY_ATTEMPTS = 10 USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM `+
			`USER_TASK_TIMEOUT_MS = 12000 USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30 QUOTED_IDENTIFIERS_IGNORE_CASE = true ENABLE_CONSOLE_OUTPUT = true `+
			`COMMENT = 'comment' TAG (%s = 'v1')`,
			id.FullyQualifiedName(), externalVolumeId.Name(), catalogId.Name(), tagId.FullyQualifiedName())
	})
}

func TestSchemasAlter(t *testing.T) {
	schemaId := randomDatabaseObjectIdentifier()
	newSchemaId := randomDatabaseObjectIdentifierInDatabase(schemaId.DatabaseId())

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: emptyDatabaseObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: emptyDatabaseObjectIdentifier,
			Set: &SchemaSet{
				ExternalVolume: new(emptyAccountObjectIdentifier),
				Catalog:        new(emptyAccountObjectIdentifier),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("SchemaSet", "ExternalVolume"), errInvalidIdentifier("SchemaSet", "Catalog"))
	})

	t.Run("validation: at least one of actions", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSchemaOptions", "NewName", "SwapWith", "Set", "Unset", "SetTag", "UnsetTag", "EnableManagedAccess", "DisableManagedAccess"))
	})

	t.Run("validation: exactly one of actions", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:  schemaId,
			Set:   &SchemaSet{},
			Unset: &SchemaUnset{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSchemaOptions", "NewName", "SwapWith", "Set", "Unset", "SetTag", "UnsetTag", "EnableManagedAccess", "DisableManagedAccess"))
	})

	t.Run("validation: at least one set option", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
			Set:  &SchemaSet{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf(
			"SchemaSet",
			"DataRetentionTimeInDays",
			"MaxDataExtensionTimeInDays",
			"ExternalVolume",
			"Catalog",
			"ReplaceInvalidCharacters",
			"DefaultDDLCollation",
			"StorageSerializationPolicy",
			"LogLevel",
			"LogEventLevel",
			"TraceLevel",
			"SuspendTaskAfterNumFailures",
			"TaskAutoRetryAttempts",
			"UserTaskManagedInitialWarehouseSize",
			"UserTaskTimeoutMs",
			"UserTaskMinimumTriggerIntervalInSeconds",
			"QuotedIdentifiersIgnoreCase",
			"EnableConsoleOutput",
			"PipeExecutionPaused",
			"Comment",
		))
	})

	t.Run("validation: at least one unset option", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:  schemaId,
			Unset: &SchemaUnset{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf(
			"SchemaUnset",
			"DataRetentionTimeInDays",
			"MaxDataExtensionTimeInDays",
			"ExternalVolume",
			"Catalog",
			"ReplaceInvalidCharacters",
			"DefaultDDLCollation",
			"StorageSerializationPolicy",
			"LogLevel",
			"LogEventLevel",
			"TraceLevel",
			"SuspendTaskAfterNumFailures",
			"TaskAutoRetryAttempts",
			"UserTaskManagedInitialWarehouseSize",
			"UserTaskTimeoutMs",
			"UserTaskMinimumTriggerIntervalInSeconds",
			"QuotedIdentifiersIgnoreCase",
			"EnableConsoleOutput",
			"PipeExecutionPaused",
			"Comment",
		))
	})

	t.Run("validation: invalid external volume identifier", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
			Set: &SchemaSet{
				ExternalVolume: new(emptyAccountObjectIdentifier),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("SchemaSet", "ExternalVolume"))
	})

	t.Run("validation: invalid catalog integration identifier", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
			Set: &SchemaSet{
				Catalog: new(emptyAccountObjectIdentifier),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("SchemaSet", "Catalog"))
	})

	t.Run("validation: invalid NewName identifier", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:    schemaId,
			NewName: new(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("AlterSchemaOptions", "NewName"))
	})

	t.Run("validation: invalid SwapWith identifier", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:     schemaId,
			SwapWith: new(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("AlterSchemaOptions", "SwapWith"))
	})
	t.Run("rename to", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:     schemaId,
			IfExists: new(true),
			NewName:  new(newSchemaId),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA IF EXISTS %s RENAME TO %s`, schemaId.FullyQualifiedName(), newSchemaId.FullyQualifiedName())
	})

	t.Run("swap with", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:     schemaId,
			IfExists: new(false),
			SwapWith: new(newSchemaId),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA %s SWAP WITH %s`, schemaId.FullyQualifiedName(), newSchemaId.FullyQualifiedName())
	})

	t.Run("set options", func(t *testing.T) {
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts := &AlterSchemaOptions{
			name: schemaId,
			Set: &SchemaSet{
				DataRetentionTimeInDays:                 new(1),
				MaxDataExtensionTimeInDays:              new(1),
				ExternalVolume:                          &externalVolumeId,
				Catalog:                                 &catalogId,
				PipeExecutionPaused:                     new(true),
				ReplaceInvalidCharacters:                new(true),
				DefaultDDLCollation:                     new(StringAllowEmpty{Value: "en_US-trim"}),
				StorageSerializationPolicy:              new(StorageSerializationPolicyCompatible),
				LogLevel:                                new(LogLevelInfo),
				TraceLevel:                              new(TraceLevelPropagate),
				SuspendTaskAfterNumFailures:             new(10),
				TaskAutoRetryAttempts:                   new(10),
				UserTaskManagedInitialWarehouseSize:     new(WarehouseSizeMedium),
				UserTaskTimeoutMs:                       new(12000),
				UserTaskMinimumTriggerIntervalInSeconds: new(30),
				QuotedIdentifiersIgnoreCase:             new(true),
				EnableConsoleOutput:                     new(true),
				Comment:                                 new("comment"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA %s SET DATA_RETENTION_TIME_IN_DAYS = 1, MAX_DATA_EXTENSION_TIME_IN_DAYS = 1, `+
			`EXTERNAL_VOLUME = "%s", CATALOG = "%s", PIPE_EXECUTION_PAUSED = true, REPLACE_INVALID_CHARACTERS = true, DEFAULT_DDL_COLLATION = 'en_US-trim', STORAGE_SERIALIZATION_POLICY = COMPATIBLE, `+
			`LOG_LEVEL = 'INFO', TRACE_LEVEL = 'PROPAGATE', SUSPEND_TASK_AFTER_NUM_FAILURES = 10, TASK_AUTO_RETRY_ATTEMPTS = 10, USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM, `+
			`USER_TASK_TIMEOUT_MS = 12000, USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30, QUOTED_IDENTIFIERS_IGNORE_CASE = true, ENABLE_CONSOLE_OUTPUT = true, `+
			`COMMENT = 'comment'`,
			schemaId.FullyQualifiedName(), externalVolumeId.Name(), catalogId.Name(),
		)
	})

	t.Run("unset", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
			Unset: &SchemaUnset{
				DataRetentionTimeInDays:                 new(true),
				MaxDataExtensionTimeInDays:              new(true),
				ExternalVolume:                          new(true),
				Catalog:                                 new(true),
				PipeExecutionPaused:                     new(true),
				ReplaceInvalidCharacters:                new(true),
				DefaultDDLCollation:                     new(true),
				StorageSerializationPolicy:              new(true),
				LogLevel:                                new(true),
				TraceLevel:                              new(true),
				SuspendTaskAfterNumFailures:             new(true),
				TaskAutoRetryAttempts:                   new(true),
				UserTaskManagedInitialWarehouseSize:     new(true),
				UserTaskTimeoutMs:                       new(true),
				UserTaskMinimumTriggerIntervalInSeconds: new(true),
				QuotedIdentifiersIgnoreCase:             new(true),
				EnableConsoleOutput:                     new(true),
				Comment:                                 new(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA %s UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, EXTERNAL_VOLUME, CATALOG, PIPE_EXECUTION_PAUSED, `+
			`REPLACE_INVALID_CHARACTERS, DEFAULT_DDL_COLLATION, STORAGE_SERIALIZATION_POLICY, LOG_LEVEL, TRACE_LEVEL, SUSPEND_TASK_AFTER_NUM_FAILURES, TASK_AUTO_RETRY_ATTEMPTS, `+
			`USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE, USER_TASK_TIMEOUT_MS, USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS, QUOTED_IDENTIFIERS_IGNORE_CASE, ENABLE_CONSOLE_OUTPUT, COMMENT`, opts.name.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
			SetTag: []TagAssociation{
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
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`, schemaId.FullyQualifiedName())
	})

	t.Run("unset tags", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
			UnsetTag: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag1"),
				NewAccountObjectIdentifier("tag2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA %s UNSET TAG "tag1", "tag2"`, schemaId.FullyQualifiedName())
	})

	t.Run("enable managed access", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:                schemaId,
			EnableManagedAccess: new(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA %s ENABLE MANAGED ACCESS`, schemaId.FullyQualifiedName())
	})

	t.Run("disable managed access", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:                 schemaId,
			DisableManagedAccess: new(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA %s DISABLE MANAGED ACCESS`, schemaId.FullyQualifiedName())
	})
}

func TestSchemasDrop(t *testing.T) {
	schemaId := randomDatabaseObjectIdentifier()

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := &DropSchemaOptions{
			name: emptyDatabaseObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("minimal", func(t *testing.T) {
		opts := &DropSchemaOptions{
			name: schemaId,
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP SCHEMA %s`, opts.name.FullyQualifiedName())
	})

	t.Run("all options - cascade", func(t *testing.T) {
		opts := &DropSchemaOptions{
			name:     schemaId,
			IfExists: new(true),
			Cascade:  new(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP SCHEMA IF EXISTS %s CASCADE`, opts.name.FullyQualifiedName())
	})

	t.Run("all options - restrict", func(t *testing.T) {
		opts := &DropSchemaOptions{
			name:     schemaId,
			IfExists: new(true),
			Restrict: new(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP SCHEMA IF EXISTS %s RESTRICT`, opts.name.FullyQualifiedName())
	})

	t.Run("validation: cascade and restrict set together", func(t *testing.T) {
		opts := &DropSchemaOptions{
			name:     schemaId,
			IfExists: new(true),
			Cascade:  new(true),
			Restrict: new(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("DropSchemaOptions", "Cascade", "Restrict"))
	})
}

func TestSchemasUndrop(t *testing.T) {
	schemaId := randomDatabaseObjectIdentifier()

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := &undropSchemaOptions{
			name: emptyDatabaseObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("minimal", func(t *testing.T) {
		opts := &undropSchemaOptions{
			name: schemaId,
		}
		assertOptsValidAndSQLEquals(t, opts, `UNDROP SCHEMA %s`, opts.name.FullyQualifiedName())
	})
}

func TestSchemasDescribe(t *testing.T) {
	schemaId := randomDatabaseObjectIdentifier()

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := &describeSchemaOptions{
			name: emptyDatabaseObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("complete", func(t *testing.T) {
		opts := &describeSchemaOptions{
			name: schemaId,
		}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE SCHEMA %s`, opts.name.FullyQualifiedName())
	})
}

func TestSchemasShow(t *testing.T) {
	t.Run("like", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   new(true),
			History: new(true),
			Like: &Like{
				Pattern: new("schema_pattern"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY LIKE 'schema_pattern'`)
	})

	t.Run("in account", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   new(true),
			History: new(true),
			In: &SchemaIn{
				Account: new(true),
				Name:    NewAccountObjectIdentifier("account_name"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY IN ACCOUNT "account_name"`)
	})

	t.Run("in database", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   new(true),
			History: new(true),
			In: &SchemaIn{
				Database: new(true),
				Name:     NewAccountObjectIdentifier("database_name"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY IN DATABASE "database_name"`)
	})

	t.Run("starts with", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:      new(true),
			History:    new(true),
			StartsWith: new("schema_pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY STARTS WITH 'schema_pattern'`)
	})

	t.Run("limit", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   new(true),
			History: new(true),
			LimitFrom: &LimitFrom{
				Rows: new(3),
				From: new("name_string"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY LIMIT 3 FROM 'name_string'`)
	})
}
