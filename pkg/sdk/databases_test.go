package sdk

import (
	"errors"
	"testing"
	"time"
)

func TestDatabasesCreate(t *testing.T) {
	defaultOpts := func() *CreateDatabaseOptions {
		return &CreateDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: or replace and if not exists set at once", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateDatabaseOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalVolume = Pointer(emptyAccountObjectIdentifier)
		opts.Catalog = Pointer(emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			ErrInvalidObjectIdentifier,
		)
	})

	t.Run("complete", func(t *testing.T) {
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.Transient = Bool(true)

		opts.DataRetentionTimeInDays = Int(1)
		opts.MaxDataExtensionTimeInDays = Int(1)
		opts.ExternalVolume = &externalVolumeId
		opts.Catalog = &catalogId
		opts.ReplaceInvalidCharacters = Bool(true)
		opts.DefaultDdlCollation = String("en_US")
		opts.StorageSerializationPolicy = Pointer(StorageSerializationPolicyCompatible)
		opts.LogLevel = Pointer(LogLevelInfo)
		opts.TraceLevel = Pointer(TraceLevelPropagate)
		opts.SuspendTaskAfterNumFailures = Int(10)
		opts.TaskAutoRetryAttempts = Int(10)
		opts.UserTaskManagedInitialWarehouseSize = Pointer(WarehouseSizeMedium)
		opts.UserTaskTimeoutMs = Int(12000)
		opts.UserTaskMinimumTriggerIntervalInSeconds = Int(30)
		opts.QuotedIdentifiersIgnoreCase = Bool(true)
		opts.EnableConsoleOutput = Bool(true)

		opts.Comment = String("comment")
		tagId := randomAccountObjectIdentifier()
		opts.Tag = []TagAssociation{
			{
				Name:  tagId,
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE TRANSIENT DATABASE IF NOT EXISTS %s DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 EXTERNAL_VOLUME = %s CATALOG = %s REPLACE_INVALID_CHARACTERS = true DEFAULT_DDL_COLLATION = 'en_US' STORAGE_SERIALIZATION_POLICY = COMPATIBLE LOG_LEVEL = 'INFO' TRACE_LEVEL = 'PROPAGATE' SUSPEND_TASK_AFTER_NUM_FAILURES = 10 TASK_AUTO_RETRY_ATTEMPTS = 10 USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM USER_TASK_TIMEOUT_MS = 12000 USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30 QUOTED_IDENTIFIERS_IGNORE_CASE = true ENABLE_CONSOLE_OUTPUT = true COMMENT = 'comment' TAG (%s = 'v1')`, opts.name.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName(), tagId.FullyQualifiedName())
	})
}

func TestDatabasesClone(t *testing.T) {
	defaultOpts := func() *CloneDatabaseOptions {
		return &CloneDatabaseOptions{
			name:  randomAccountObjectIdentifier(),
			Clone: Clone{SourceObject: randomAccountObjectIdentifier()},
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid clone", func(t *testing.T) {
		opts := defaultOpts()
		opts.Clone = Clone{
			SourceObject: emptyAccountObjectIdentifier,
			At: &TimeTravel{
				Timestamp: Pointer(time.Now()),
				Offset:    Int(123),
			},
			Before: new(TimeTravel),
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			errors.New("only one of AT or BEFORE can be set"),
			errors.New("exactly one of TIMESTAMP, OFFSET or STATEMENT can be set"),
		)
	})

	t.Run("validation: or replace and if not exists set at once", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CloneDatabaseOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("complete", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Clone = Clone{
			SourceObject: NewAccountObjectIdentifier("db1"),
			At: &TimeTravel{
				Timestamp: Pointer(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE DATABASE %s CLONE "db1" AT (TIMESTAMP => '2021-01-01 00:00:00 +0000 UTC')`, opts.name.FullyQualifiedName())
	})
}

func TestDatabasesCreateShared(t *testing.T) {
	defaultOpts := func() *CreateSharedDatabaseOptions {
		return &CreateSharedDatabaseOptions{
			name:      randomAccountObjectIdentifier(),
			FromShare: NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account"), randomAccountObjectIdentifier()),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid from share name", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromShare = NewExternalObjectIdentifier(NewAccountIdentifier("", ""), emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: or replace and if not exists set at once", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = randomAccountObjectIdentifier()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSharedDatabaseOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("db")
		opts.ExternalVolume = Pointer(emptyAccountObjectIdentifier)
		opts.Catalog = Pointer(emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			ErrInvalidObjectIdentifier,
		)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.Transient = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE TRANSIENT DATABASE IF NOT EXISTS %s FROM SHARE %s`, opts.name.FullyQualifiedName(), opts.FromShare.FullyQualifiedName())
	})

	t.Run("complete", func(t *testing.T) {
		opts := defaultOpts()
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts.OrReplace = Bool(true)

		opts.ExternalVolume = &externalVolumeId
		opts.Catalog = &catalogId
		opts.ReplaceInvalidCharacters = Bool(true)
		opts.DefaultDdlCollation = String("en_US")
		opts.StorageSerializationPolicy = Pointer(StorageSerializationPolicyCompatible)
		opts.LogLevel = Pointer(LogLevelInfo)
		opts.TraceLevel = Pointer(TraceLevelPropagate)
		opts.SuspendTaskAfterNumFailures = Int(10)
		opts.TaskAutoRetryAttempts = Int(10)
		opts.UserTaskManagedInitialWarehouseSize = Pointer(WarehouseSizeMedium)
		opts.UserTaskTimeoutMs = Int(12000)
		opts.UserTaskMinimumTriggerIntervalInSeconds = Int(30)
		opts.QuotedIdentifiersIgnoreCase = Bool(true)
		opts.EnableConsoleOutput = Bool(true)

		opts.Comment = String("comment")
		tagId := randomAccountObjectIdentifier()
		opts.Tag = []TagAssociation{
			{
				Name:  tagId,
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE DATABASE %s FROM SHARE %s EXTERNAL_VOLUME = %s CATALOG = %s REPLACE_INVALID_CHARACTERS = true DEFAULT_DDL_COLLATION = 'en_US' STORAGE_SERIALIZATION_POLICY = COMPATIBLE LOG_LEVEL = 'INFO' TRACE_LEVEL = 'PROPAGATE' SUSPEND_TASK_AFTER_NUM_FAILURES = 10 TASK_AUTO_RETRY_ATTEMPTS = 10 USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM USER_TASK_TIMEOUT_MS = 12000 USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30 QUOTED_IDENTIFIERS_IGNORE_CASE = true ENABLE_CONSOLE_OUTPUT = true COMMENT = 'comment' TAG (%s = 'v1')`, opts.name.FullyQualifiedName(), opts.FromShare.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName(), tagId.FullyQualifiedName())
	})
}

func TestDatabasesCreateSecondary(t *testing.T) {
	defaultOpts := func() *CreateSecondaryDatabaseOptions {
		return &CreateSecondaryDatabaseOptions{
			name:            randomAccountObjectIdentifier(),
			PrimaryDatabase: NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account"), randomAccountObjectIdentifier()),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid primary database", func(t *testing.T) {
		opts := defaultOpts()
		opts.PrimaryDatabase = NewExternalObjectIdentifier(NewAccountIdentifier("", ""), emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: or replace and if not exists set at once", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.IfNotExists = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSecondaryDatabaseOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalVolume = Pointer(emptyAccountObjectIdentifier)
		opts.Catalog = Pointer(emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			ErrInvalidObjectIdentifier,
		)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE IF NOT EXISTS %s AS REPLICA OF %s`, opts.name.FullyQualifiedName(), opts.PrimaryDatabase.FullyQualifiedName())
	})

	t.Run("complete", func(t *testing.T) {
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		primaryDatabaseId := NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account"), randomAccountObjectIdentifier())
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		opts.Transient = Bool(true)
		opts.PrimaryDatabase = primaryDatabaseId

		opts.DataRetentionTimeInDays = Int(1)
		opts.MaxDataExtensionTimeInDays = Int(1)
		opts.ExternalVolume = &externalVolumeId
		opts.Catalog = &catalogId
		opts.ReplaceInvalidCharacters = Bool(true)
		opts.DefaultDdlCollation = String("en_US")
		opts.StorageSerializationPolicy = Pointer(StorageSerializationPolicyCompatible)
		opts.LogLevel = Pointer(LogLevelInfo)
		opts.TraceLevel = Pointer(TraceLevelPropagate)
		opts.SuspendTaskAfterNumFailures = Int(10)
		opts.TaskAutoRetryAttempts = Int(10)
		opts.UserTaskManagedInitialWarehouseSize = Pointer(WarehouseSizeMedium)
		opts.UserTaskTimeoutMs = Int(12000)
		opts.UserTaskMinimumTriggerIntervalInSeconds = Int(30)
		opts.QuotedIdentifiersIgnoreCase = Bool(true)
		opts.EnableConsoleOutput = Bool(true)

		opts.Comment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TRANSIENT DATABASE %s AS REPLICA OF %s DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 EXTERNAL_VOLUME = %s CATALOG = %s REPLACE_INVALID_CHARACTERS = true DEFAULT_DDL_COLLATION = 'en_US' STORAGE_SERIALIZATION_POLICY = COMPATIBLE LOG_LEVEL = 'INFO' TRACE_LEVEL = 'PROPAGATE' SUSPEND_TASK_AFTER_NUM_FAILURES = 10 TASK_AUTO_RETRY_ATTEMPTS = 10 USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM USER_TASK_TIMEOUT_MS = 12000 USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30 QUOTED_IDENTIFIERS_IGNORE_CASE = true ENABLE_CONSOLE_OUTPUT = true COMMENT = 'comment'`, opts.name.FullyQualifiedName(), primaryDatabaseId.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName())
	})
}

func TestDatabasesCreateFromListing(t *testing.T) {
	defaultOpts := func() *CreateFromListingDatabaseOptions {
		return &CreateFromListingDatabaseOptions{
			name:        randomAccountObjectIdentifier(),
			FromListing: "GZ1M7Z91WTX",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateFromListingDatabaseOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: empty listing global name", func(t *testing.T) {
		opts := defaultOpts()
		opts.FromListing = ""
		assertOptsInvalidJoinedErrors(t, opts, NewError("CreateFromListingDatabaseOptions: listing global name must not be empty"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE %s FROM LISTING '%s'`, opts.name.FullyQualifiedName(), opts.FromListing)
	})
}

func TestDatabasesAlter(t *testing.T) {
	defaultOpts := func() *AlterDatabaseOptions {
		return &AlterDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{
			ExternalVolume: Pointer(emptyAccountObjectIdentifier),
			Catalog:        Pointer(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one of actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{}
		opts.Unset = &DatabaseUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseOptions", "RenameTo", "Set", "Unset", "SwapWith", "SetTags", "UnsetTags"))
	})

	t.Run("validation: at least one set option", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf(
			"AlterDatabaseOptions.Set",
			"DataRetentionTimeInDays",
			"MaxDataExtensionTimeInDays",
			"ExternalVolume",
			"Catalog",
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
		opts := defaultOpts()
		opts.Unset = &DatabaseUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf(
			"AlterDatabaseOptions.Unset",
			"DataRetentionTimeInDays",
			"MaxDataExtensionTimeInDays",
			"ExternalVolume",
			"Catalog",
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
		opts := defaultOpts()
		opts.Set = &DatabaseSet{
			ExternalVolume: Pointer(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid catalog integration identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{
			Catalog: Pointer(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid RenameTo identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.RenameTo = Pointer(emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid SwapWith identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.SwapWith = Pointer(emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("rename", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.RenameTo = Pointer(randomAccountObjectIdentifier())
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE IF EXISTS %s RENAME TO %s`, opts.name.FullyQualifiedName(), opts.RenameTo.FullyQualifiedName())
	})

	t.Run("swap with", func(t *testing.T) {
		opts := defaultOpts()
		opts.SwapWith = Pointer(randomAccountObjectIdentifier())
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s SWAP WITH %s`, opts.name.FullyQualifiedName(), opts.SwapWith.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &DatabaseSet{
			DataRetentionTimeInDays:    Int(1),
			MaxDataExtensionTimeInDays: Int(1),
			ExternalVolume:             &externalVolumeId,
			Catalog:                    &catalogId,
			ReplaceInvalidCharacters:   Bool(true),
			DefaultDdlCollation:        String("en_US"),
			StorageSerializationPolicy: Pointer(StorageSerializationPolicyCompatible),
			LogLevel:                   Pointer(LogLevelError),
			TraceLevel:                 Pointer(TraceLevelPropagate),
			Comment:                    String("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s SET DATA_RETENTION_TIME_IN_DAYS = 1, MAX_DATA_EXTENSION_TIME_IN_DAYS = 1, EXTERNAL_VOLUME = %s, CATALOG = %s, REPLACE_INVALID_CHARACTERS = true, DEFAULT_DDL_COLLATION = 'en_US', STORAGE_SERIALIZATION_POLICY = COMPATIBLE, LOG_LEVEL = 'ERROR', TRACE_LEVEL = 'PROPAGATE', COMMENT = 'comment'`, opts.name.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName())
	})

	t.Run("unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &DatabaseUnset{
			DataRetentionTimeInDays:    Bool(true),
			MaxDataExtensionTimeInDays: Bool(true),
			ExternalVolume:             Bool(true),
			Catalog:                    Bool(true),
			ReplaceInvalidCharacters:   Bool(true),
			DefaultDdlCollation:        Bool(true),
			StorageSerializationPolicy: Bool(true),
			LogLevel:                   Bool(true),
			TraceLevel:                 Bool(true),
			Comment:                    Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, EXTERNAL_VOLUME, CATALOG, REPLACE_INVALID_CHARACTERS, DEFAULT_DDL_COLLATION, STORAGE_SERIALIZATION_POLICY, LOG_LEVEL, TRACE_LEVEL, COMMENT`, opts.name.FullyQualifiedName())
	})

	t.Run("with set tag", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  tagId1,
				Value: "v1",
			},
			{
				Name:  tagId2,
				Value: "v2",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s SET TAG %s = 'v1', %s = 'v2'`, opts.name.FullyQualifiedName(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})

	t.Run("with unset tag", func(t *testing.T) {
		id := randomSchemaObjectIdentifier()
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			id,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s UNSET TAG %s`, opts.name.FullyQualifiedName(), id.FullyQualifiedName())
	})
}

func TestDatabasesAlterReplication(t *testing.T) {
	defaultOpts := func() *AlterReplicationDatabaseOptions {
		return &AlterReplicationDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterReplicationDatabaseOptions", "EnableReplication", "DisableReplication", "Refresh"))
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableReplication = &EnableReplication{}
		opts.DisableReplication = &DisableReplication{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterReplicationDatabaseOptions", "EnableReplication", "DisableReplication", "Refresh"))
	})

	t.Run("enable replication", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableReplication = &EnableReplication{
			ToAccounts: []AccountIdentifier{
				NewAccountIdentifierFromAccountLocator("account1"),
			},
			IgnoreEditionCheck: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s ENABLE REPLICATION TO ACCOUNTS "account1" IGNORE EDITION CHECK`, opts.name.FullyQualifiedName())
	})

	t.Run("disable replication", func(t *testing.T) {
		opts := defaultOpts()
		opts.DisableReplication = &DisableReplication{
			ToAccounts: []AccountIdentifier{
				NewAccountIdentifierFromAccountLocator("account1"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s DISABLE REPLICATION TO ACCOUNTS "account1"`, opts.name.FullyQualifiedName())
	})

	t.Run("refresh", func(t *testing.T) {
		opts := defaultOpts()
		opts.Refresh = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s REFRESH`, opts.name.FullyQualifiedName())
	})
}

func TestDatabasesAlterFailover(t *testing.T) {
	defaultOpts := func() *AlterFailoverDatabaseOptions {
		return &AlterFailoverDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterFailoverDatabaseOptions", "EnableFailover", "DisableFailover", "Primary"))
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableFailover = &EnableFailover{}
		opts.DisableFailover = &DisableFailover{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterFailoverDatabaseOptions", "EnableFailover", "DisableFailover", "Primary"))
	})

	t.Run("enable failover", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableFailover = &EnableFailover{
			ToAccounts: []AccountIdentifier{
				NewAccountIdentifierFromAccountLocator("account1"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s ENABLE FAILOVER TO ACCOUNTS "account1"`, opts.name.FullyQualifiedName())
	})

	t.Run("disable failover", func(t *testing.T) {
		opts := defaultOpts()
		opts.DisableFailover = &DisableFailover{
			ToAccounts: []AccountIdentifier{
				NewAccountIdentifierFromAccountLocator("account1"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s DISABLE FAILOVER TO ACCOUNTS "account1"`, opts.name.FullyQualifiedName())
	})

	t.Run("primary", func(t *testing.T) {
		opts := defaultOpts()
		opts.Primary = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s PRIMARY`, opts.name.FullyQualifiedName())
	})
}

func TestDatabasesDrop(t *testing.T) {
	defaultOpts := func() *DropDatabaseOptions {
		return &DropDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("minimal", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE %s`, opts.name.FullyQualifiedName())
	})

	t.Run("all options - cascade", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Cascade = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE IF EXISTS %s CASCADE`, opts.name.FullyQualifiedName())
	})

	t.Run("all options - restrict", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		opts.Restrict = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE IF EXISTS %s RESTRICT`, opts.name.FullyQualifiedName())
	})

	t.Run("validation: cascade and restrict set together", func(t *testing.T) {
		opts := defaultOpts()
		opts.Cascade = Bool(true)
		opts.Restrict = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("DropDatabaseOptions", "Cascade", "Restrict"))
	})
}

func TestDatabasesUndrop(t *testing.T) {
	defaultOpts := func() *UndropDatabaseOptions {
		return &UndropDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("minimal", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `UNDROP DATABASE %s`, opts.name.FullyQualifiedName())
	})
}

func TestDatabasesShow(t *testing.T) {
	defaultOpts := func() *ShowDatabaseOptions {
		return &ShowDatabaseOptions{}
	}

	t.Run("without show options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASES`)
	})

	t.Run("terse", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE DATABASES`)
	})

	t.Run("history", func(t *testing.T) {
		opts := defaultOpts()
		opts.History = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASES HISTORY`)
	})

	t.Run("like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String("db1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASES LIKE 'db1'`)
	})

	t.Run("complete", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		opts.History = Bool(true)
		opts.Like = &Like{
			Pattern: String("db2"),
		}
		opts.Limit = &LimitFrom{
			Rows: Int(1),
			From: String("db1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE DATABASES HISTORY LIKE 'db2' LIMIT 1 FROM 'db1'`)
	})
}

func TestDatabasesDescribe(t *testing.T) {
	defaultOpts := func() *describeDatabaseOptions {
		return &describeDatabaseOptions{
			name: randomAccountObjectIdentifier(),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("complete", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE DATABASE %s`, opts.name.FullyQualifiedName())
	})
}
