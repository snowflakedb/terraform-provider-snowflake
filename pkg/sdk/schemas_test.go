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
			OrReplace:   Bool(true),
			IfNotExists: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSchemaOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := &CreateSchemaOptions{
			name:           id,
			ExternalVolume: Pointer(emptyAccountObjectIdentifier),
			Catalog:        Pointer(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			ErrInvalidObjectIdentifier,
		)
	})

	t.Run("complete", func(t *testing.T) {
		tagId := randomSchemaObjectIdentifier()
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts := &CreateSchemaOptions{
			Transient:                               Bool(true),
			IfNotExists:                             Bool(true),
			name:                                    id,
			WithManagedAccess:                       Bool(true),
			DataRetentionTimeInDays:                 Int(1),
			MaxDataExtensionTimeInDays:              Int(1),
			ExternalVolume:                          &externalVolumeId,
			Catalog:                                 &catalogId,
			PipeExecutionPaused:                     Bool(true),
			ReplaceInvalidCharacters:                Bool(true),
			DefaultDdlCollation:                     Pointer(StringAllowEmpty{Value: "en_US-trim"}),
			StorageSerializationPolicy:              Pointer(StorageSerializationPolicyCompatible),
			LogLevel:                                Pointer(LogLevelInfo),
			TraceLevel:                              Pointer(TraceLevelPropagate),
			SuspendTaskAfterNumFailures:             Int(10),
			TaskAutoRetryAttempts:                   Int(10),
			UserTaskManagedInitialWarehouseSize:     Pointer(WarehouseSizeMedium),
			UserTaskTimeoutMs:                       Int(12000),
			UserTaskMinimumTriggerIntervalInSeconds: Int(30),
			QuotedIdentifiersIgnoreCase:             Bool(true),
			EnableConsoleOutput:                     Bool(true),
			Comment:                                 String("comment"),
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

func TestSchemasClone(t *testing.T) {
	id := randomDatabaseObjectIdentifier()

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CloneSchemaOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := &CloneSchemaOptions{
			name: emptyDatabaseObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.OrReplace opts.IfNotExists]", func(t *testing.T) {
		opts := &CloneSchemaOptions{
			name:        id,
			OrReplace:   Bool(true),
			IfNotExists: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CloneSchemaOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("complete", func(t *testing.T) {
		opts := &CloneSchemaOptions{
			name:      id,
			OrReplace: Bool(true),
			Clone: Clone{
				SourceObject: NewAccountObjectIdentifier("sch1"),
				At: &TimeTravel{
					Timestamp: Pointer(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE SCHEMA %s CLONE "sch1" AT (TIMESTAMP => '2021-01-01 00:00:00 +0000 UTC')`, id.FullyQualifiedName())
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
				ExternalVolume: Pointer(emptyAccountObjectIdentifier),
				Catalog:        Pointer(emptyAccountObjectIdentifier),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: at least one of actions", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSchemaOptions", "RenameTo", "SwapWith", "Set", "Unset", "SetTags", "UnsetTags", "EnableManagedAccess", "DisableManagedAccess"))
	})

	t.Run("validation: exactly one of actions", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:  schemaId,
			Set:   &SchemaSet{},
			Unset: &SchemaUnset{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSchemaOptions", "RenameTo", "SwapWith", "Set", "Unset", "SetTags", "UnsetTags", "EnableManagedAccess", "DisableManagedAccess"))
	})

	t.Run("validation: at least one set option", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
			Set:  &SchemaSet{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf(
			"AlterSchemaOptions.Set",
			"DataRetentionTimeInDays",
			"MaxDataExtensionTimeInDays",
			"ExternalVolume",
			"Catalog",
			"PipeExecutionPaused",
			"ReplaceInvalidCharacters",
			"DefaultDdlCollation",
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
			"Comment",
		))
	})

	t.Run("validation: at least one unset option", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:  schemaId,
			Unset: &SchemaUnset{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf(
			"AlterSchemaOptions.Unset",
			"DataRetentionTimeInDays",
			"MaxDataExtensionTimeInDays",
			"ExternalVolume",
			"Catalog",
			"PipeExecutionPaused",
			"ReplaceInvalidCharacters",
			"DefaultDdlCollation",
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
			"Comment",
		))
	})

	t.Run("validation: invalid external volume identifier", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
			Set: &SchemaSet{
				ExternalVolume: Pointer(emptyAccountObjectIdentifier),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid catalog integration identifier", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
			Set: &SchemaSet{
				Catalog: Pointer(emptyAccountObjectIdentifier),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid RenameTo identifier", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:     schemaId,
			RenameTo: Pointer(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid SwapWith identifier", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:     schemaId,
			SwapWith: Pointer(emptyDatabaseObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
	t.Run("rename to", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:     schemaId,
			IfExists: Bool(true),
			RenameTo: Pointer(newSchemaId),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA IF EXISTS %s RENAME TO %s`, schemaId.FullyQualifiedName(), newSchemaId.FullyQualifiedName())
	})

	t.Run("swap with", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:     schemaId,
			IfExists: Bool(false),
			SwapWith: Pointer(newSchemaId),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA %s SWAP WITH %s`, schemaId.FullyQualifiedName(), newSchemaId.FullyQualifiedName())
	})

	t.Run("set options", func(t *testing.T) {
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts := &AlterSchemaOptions{
			name: schemaId,
			Set: &SchemaSet{
				DataRetentionTimeInDays:                 Int(1),
				MaxDataExtensionTimeInDays:              Int(1),
				ExternalVolume:                          &externalVolumeId,
				Catalog:                                 &catalogId,
				PipeExecutionPaused:                     Bool(true),
				ReplaceInvalidCharacters:                Bool(true),
				DefaultDdlCollation:                     Pointer(StringAllowEmpty{Value: "en_US-trim"}),
				StorageSerializationPolicy:              Pointer(StorageSerializationPolicyCompatible),
				LogLevel:                                Pointer(LogLevelInfo),
				TraceLevel:                              Pointer(TraceLevelPropagate),
				SuspendTaskAfterNumFailures:             Int(10),
				TaskAutoRetryAttempts:                   Int(10),
				UserTaskManagedInitialWarehouseSize:     Pointer(WarehouseSizeMedium),
				UserTaskTimeoutMs:                       Int(12000),
				UserTaskMinimumTriggerIntervalInSeconds: Int(30),
				QuotedIdentifiersIgnoreCase:             Bool(true),
				EnableConsoleOutput:                     Bool(true),
				Comment:                                 String("comment"),
			},
		}
		assertOptsValidAndSQLEquals(
			t, opts, `ALTER SCHEMA %s SET DATA_RETENTION_TIME_IN_DAYS = 1, MAX_DATA_EXTENSION_TIME_IN_DAYS = 1, `+
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
				DataRetentionTimeInDays:                 Bool(true),
				MaxDataExtensionTimeInDays:              Bool(true),
				ExternalVolume:                          Bool(true),
				Catalog:                                 Bool(true),
				PipeExecutionPaused:                     Bool(true),
				ReplaceInvalidCharacters:                Bool(true),
				DefaultDdlCollation:                     Bool(true),
				StorageSerializationPolicy:              Bool(true),
				LogLevel:                                Bool(true),
				TraceLevel:                              Bool(true),
				SuspendTaskAfterNumFailures:             Bool(true),
				TaskAutoRetryAttempts:                   Bool(true),
				UserTaskManagedInitialWarehouseSize:     Bool(true),
				UserTaskTimeoutMs:                       Bool(true),
				UserTaskMinimumTriggerIntervalInSeconds: Bool(true),
				QuotedIdentifiersIgnoreCase:             Bool(true),
				EnableConsoleOutput:                     Bool(true),
				Comment:                                 Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA %s UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, EXTERNAL_VOLUME, CATALOG, PIPE_EXECUTION_PAUSED, `+
			`REPLACE_INVALID_CHARACTERS, DEFAULT_DDL_COLLATION, STORAGE_SERIALIZATION_POLICY, LOG_LEVEL, TRACE_LEVEL, SUSPEND_TASK_AFTER_NUM_FAILURES, TASK_AUTO_RETRY_ATTEMPTS, `+
			`USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE, USER_TASK_TIMEOUT_MS, USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS, QUOTED_IDENTIFIERS_IGNORE_CASE, ENABLE_CONSOLE_OUTPUT, COMMENT`, opts.name.FullyQualifiedName())
	})

	t.Run("set tags", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name: schemaId,
			SetTags: []TagAssociation{
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
			UnsetTags: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag1"),
				NewAccountObjectIdentifier("tag2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA %s UNSET TAG "tag1", "tag2"`, schemaId.FullyQualifiedName())
	})

	t.Run("enable managed access", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:                schemaId,
			EnableManagedAccess: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER SCHEMA %s ENABLE MANAGED ACCESS`, schemaId.FullyQualifiedName())
	})

	t.Run("disable managed access", func(t *testing.T) {
		opts := &AlterSchemaOptions{
			name:                 schemaId,
			DisableManagedAccess: Bool(true),
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
			IfExists: Bool(true),
			Cascade:  Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP SCHEMA IF EXISTS %s CASCADE`, opts.name.FullyQualifiedName())
	})

	t.Run("all options - restrict", func(t *testing.T) {
		opts := &DropSchemaOptions{
			name:     schemaId,
			IfExists: Bool(true),
			Restrict: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP SCHEMA IF EXISTS %s RESTRICT`, opts.name.FullyQualifiedName())
	})

	t.Run("validation: cascade and restrict set together", func(t *testing.T) {
		opts := &DropSchemaOptions{
			name:     schemaId,
			IfExists: Bool(true),
			Cascade:  Bool(true),
			Restrict: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("DropSchemaOptions", "Cascade", "Restrict"))
	})
}

func TestSchemasUndrop(t *testing.T) {
	schemaId := randomDatabaseObjectIdentifier()

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := &UndropSchemaOptions{
			name: emptyDatabaseObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("minimal", func(t *testing.T) {
		opts := &UndropSchemaOptions{
			name: schemaId,
		}
		assertOptsValidAndSQLEquals(t, opts, `UNDROP SCHEMA %s`, opts.name.FullyQualifiedName())
	})
}

func TestSchemasDescribe(t *testing.T) {
	schemaId := randomDatabaseObjectIdentifier()

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := &DescribeSchemaOptions{
			name: emptyDatabaseObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("complete", func(t *testing.T) {
		opts := &DescribeSchemaOptions{
			name: schemaId,
		}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE SCHEMA %s`, opts.name.FullyQualifiedName())
	})
}

func TestSchemasShow(t *testing.T) {
	t.Run("like", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   Bool(true),
			History: Bool(true),
			Like: &Like{
				Pattern: String("schema_pattern"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY LIKE 'schema_pattern'`)
	})

	t.Run("in account", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   Bool(true),
			History: Bool(true),
			In: &ExtendedIn{
				In: In{Account: Bool(true)},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY IN ACCOUNT`)
	})

	t.Run("in database", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   Bool(true),
			History: Bool(true),
			In: &ExtendedIn{
				In: In{Database: NewAccountObjectIdentifier("database_name")},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY IN DATABASE "database_name"`)
	})

	t.Run("starts with", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:      Bool(true),
			History:    Bool(true),
			StartsWith: String("schema_pattern"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY STARTS WITH 'schema_pattern'`)
	})

	t.Run("limit", func(t *testing.T) {
		opts := &ShowSchemaOptions{
			Terse:   Bool(true),
			History: Bool(true),
			Limit: &LimitFrom{
				Rows: Int(3),
				From: String("name_string"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE SCHEMAS HISTORY LIMIT 3 FROM 'name_string'`)
	})
}
