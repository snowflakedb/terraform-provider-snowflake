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

	t.Run("validation: invalid clone", func(t *testing.T) {
		opts := defaultOpts()
		opts.Clone = &Clone{
			SourceObject: emptyAccountObjectIdentifier,
			At: &TimeTravel{
				Timestamp: new(time.Now()),
				Offset:    new(123),
			},
			Before: new(TimeTravel),
		}
		assertOptsInvalidJoinedErrors(t, opts,
			errors.New("only one of AT or BEFORE can be set"),
			errors.New("exactly one of TIMESTAMP, OFFSET or STATEMENT can be set"),
		)
	})

	t.Run("validation: or replace and if not exists set at once", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = new(true)
		opts.IfNotExists = new(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateDatabaseOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalVolume = new(emptyAccountObjectIdentifier)
		opts.Catalog = new(emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts,
			errInvalidIdentifier("CreateDatabaseOptions", "ExternalVolume"),
			errInvalidIdentifier("CreateDatabaseOptions", "Catalog"),
		)
	})

	t.Run("clone", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = new(true)
		opts.Clone = &Clone{
			SourceObject: NewAccountObjectIdentifier("db1"),
			At: &TimeTravel{
				Timestamp: new(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE DATABASE %s CLONE "db1" AT (TIMESTAMP => '2021-01-01 00:00:00 +0000 UTC')`, opts.name.FullyQualifiedName())
	})

	t.Run("complete", func(t *testing.T) {
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.IfNotExists = new(true)
		opts.Transient = new(true)

		opts.DataRetentionTimeInDays = new(1)
		opts.MaxDataExtensionTimeInDays = new(1)
		opts.ExternalVolume = &externalVolumeId
		opts.Catalog = &catalogId
		opts.ReplaceInvalidCharacters = new(true)
		opts.DefaultDDLCollation = new("en_US")
		opts.StorageSerializationPolicy = new(StorageSerializationPolicyCompatible)
		opts.LogLevel = new(LogLevelInfo)
		opts.TraceLevel = new(TraceLevelPropagate)
		opts.SuspendTaskAfterNumFailures = new(10)
		opts.TaskAutoRetryAttempts = new(10)
		opts.UserTaskManagedInitialWarehouseSize = new(WarehouseSizeMedium)
		opts.UserTaskTimeoutMs = new(12000)
		opts.UserTaskMinimumTriggerIntervalInSeconds = new(30)
		opts.QuotedIdentifiersIgnoreCase = new(true)
		opts.EnableConsoleOutput = new(true)

		opts.Comment = new("comment")
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

func TestDatabasesCreateShared(t *testing.T) {
	defaultOpts := func() *CreateSharedDatabaseOptions {
		return &CreateSharedDatabaseOptions{
			name:      randomAccountObjectIdentifier(),
			fromShare: NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account"), randomAccountObjectIdentifier()),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid from share name", func(t *testing.T) {
		opts := defaultOpts()
		opts.fromShare = NewExternalObjectIdentifier(NewAccountIdentifier("", ""), emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("CreateSharedDatabaseOptions", "fromShare"))
	})

	t.Run("validation: or replace and if not exists set at once", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = randomAccountObjectIdentifier()
		opts.OrReplace = new(true)
		opts.IfNotExists = new(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSharedDatabaseOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewAccountObjectIdentifier("db")
		opts.ExternalVolume = new(emptyAccountObjectIdentifier)
		opts.Catalog = new(emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts,
			errInvalidIdentifier("CreateSharedDatabaseOptions", "ExternalVolume"),
			errInvalidIdentifier("CreateSharedDatabaseOptions", "Catalog"),
		)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.Transient = new(true)
		opts.IfNotExists = new(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE TRANSIENT DATABASE IF NOT EXISTS %s FROM SHARE %s`, opts.name.FullyQualifiedName(), opts.fromShare.FullyQualifiedName())
	})

	t.Run("complete", func(t *testing.T) {
		opts := defaultOpts()
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts.OrReplace = new(true)

		opts.ExternalVolume = &externalVolumeId
		opts.Catalog = &catalogId
		opts.ReplaceInvalidCharacters = new(true)
		opts.DefaultDDLCollation = new("en_US")
		opts.StorageSerializationPolicy = new(StorageSerializationPolicyCompatible)
		opts.LogLevel = new(LogLevelInfo)
		opts.TraceLevel = new(TraceLevelPropagate)
		opts.SuspendTaskAfterNumFailures = new(10)
		opts.TaskAutoRetryAttempts = new(10)
		opts.UserTaskManagedInitialWarehouseSize = new(WarehouseSizeMedium)
		opts.UserTaskTimeoutMs = new(12000)
		opts.UserTaskMinimumTriggerIntervalInSeconds = new(30)
		opts.QuotedIdentifiersIgnoreCase = new(true)
		opts.EnableConsoleOutput = new(true)

		opts.Comment = new("comment")
		tagId := randomAccountObjectIdentifier()
		opts.Tag = []TagAssociation{
			{
				Name:  tagId,
				Value: "v1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE DATABASE %s FROM SHARE %s EXTERNAL_VOLUME = %s CATALOG = %s REPLACE_INVALID_CHARACTERS = true DEFAULT_DDL_COLLATION = 'en_US' STORAGE_SERIALIZATION_POLICY = COMPATIBLE LOG_LEVEL = 'INFO' TRACE_LEVEL = 'PROPAGATE' SUSPEND_TASK_AFTER_NUM_FAILURES = 10 TASK_AUTO_RETRY_ATTEMPTS = 10 USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM USER_TASK_TIMEOUT_MS = 12000 USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30 QUOTED_IDENTIFIERS_IGNORE_CASE = true ENABLE_CONSOLE_OUTPUT = true COMMENT = 'comment' TAG (%s = 'v1')`, opts.name.FullyQualifiedName(), opts.fromShare.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName(), tagId.FullyQualifiedName())
	})
}

func TestDatabasesCreateSecondary(t *testing.T) {
	defaultOpts := func() *CreateSecondaryDatabaseOptions {
		return &CreateSecondaryDatabaseOptions{
			name:            randomAccountObjectIdentifier(),
			primaryDatabase: NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account"), randomAccountObjectIdentifier()),
		}
	}

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid primary database", func(t *testing.T) {
		opts := defaultOpts()
		opts.primaryDatabase = NewExternalObjectIdentifier(NewAccountIdentifier("", ""), emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("CreateSecondaryDatabaseOptions", "primaryDatabase"))
	})

	t.Run("validation: or replace and if not exists set at once", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = new(true)
		opts.IfNotExists = new(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSecondaryDatabaseOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: invalid external volume and catalog", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalVolume = new(emptyAccountObjectIdentifier)
		opts.Catalog = new(emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts,
			errInvalidIdentifier("CreateSecondaryDatabaseOptions", "ExternalVolume"),
			errInvalidIdentifier("CreateSecondaryDatabaseOptions", "Catalog"),
		)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = new(true)
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE IF NOT EXISTS %s AS REPLICA OF %s`, opts.name.FullyQualifiedName(), opts.primaryDatabase.FullyQualifiedName())
	})

	t.Run("complete", func(t *testing.T) {
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		primaryDatabaseId := NewExternalObjectIdentifier(NewAccountIdentifierFromAccountLocator("account"), randomAccountObjectIdentifier())
		opts := defaultOpts()
		opts.OrReplace = new(true)
		opts.Transient = new(true)
		opts.primaryDatabase = primaryDatabaseId

		opts.DataRetentionTimeInDays = new(1)
		opts.MaxDataExtensionTimeInDays = new(1)
		opts.ExternalVolume = &externalVolumeId
		opts.Catalog = &catalogId
		opts.ReplaceInvalidCharacters = new(true)
		opts.DefaultDDLCollation = new("en_US")
		opts.StorageSerializationPolicy = new(StorageSerializationPolicyCompatible)
		opts.LogLevel = new(LogLevelInfo)
		opts.TraceLevel = new(TraceLevelPropagate)
		opts.SuspendTaskAfterNumFailures = new(10)
		opts.TaskAutoRetryAttempts = new(10)
		opts.UserTaskManagedInitialWarehouseSize = new(WarehouseSizeMedium)
		opts.UserTaskTimeoutMs = new(12000)
		opts.UserTaskMinimumTriggerIntervalInSeconds = new(30)
		opts.QuotedIdentifiersIgnoreCase = new(true)
		opts.EnableConsoleOutput = new(true)

		opts.Comment = new("comment")
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE TRANSIENT DATABASE %s AS REPLICA OF %s DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 EXTERNAL_VOLUME = %s CATALOG = %s REPLACE_INVALID_CHARACTERS = true DEFAULT_DDL_COLLATION = 'en_US' STORAGE_SERIALIZATION_POLICY = COMPATIBLE LOG_LEVEL = 'INFO' TRACE_LEVEL = 'PROPAGATE' SUSPEND_TASK_AFTER_NUM_FAILURES = 10 TASK_AUTO_RETRY_ATTEMPTS = 10 USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM USER_TASK_TIMEOUT_MS = 12000 USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30 QUOTED_IDENTIFIERS_IGNORE_CASE = true ENABLE_CONSOLE_OUTPUT = true COMMENT = 'comment'`, opts.name.FullyQualifiedName(), primaryDatabaseId.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName())
	})
}

func TestDatabasesCreateFromListing(t *testing.T) {
	defaultOpts := func() *CreateDatabaseFromListingOptions {
		return &CreateDatabaseFromListingOptions{
			name:        randomAccountObjectIdentifier(),
			fromListing: "GZ1M7Z91WTX",
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateDatabaseFromListingOptions
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid name", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptyAccountObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: empty listing global name", func(t *testing.T) {
		opts := defaultOpts()
		opts.fromListing = ""
		assertOptsInvalidJoinedErrors(t, opts, NewError("CreateDatabaseFromListingOptions: listing global name must not be empty"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `CREATE DATABASE %s FROM LISTING '%s'`, opts.name.FullyQualifiedName(), opts.fromListing)
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
			ExternalVolume: new(emptyAccountObjectIdentifier),
			Catalog:        new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("DatabaseSet", "ExternalVolume"), errInvalidIdentifier("DatabaseSet", "Catalog"))
	})

	t.Run("validation: exactly one of actions", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseOptions", "NewName", "Set", "Unset", "SwapWith", "SetTag", "UnsetTag"))
	})

	t.Run("validation: exactly one of actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{}
		opts.Unset = &DatabaseUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseOptions", "NewName", "Set", "Unset", "SwapWith", "SetTag", "UnsetTag"))
	})

	t.Run("validation: at least one set option", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf(
			"DatabaseSet",
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
			"Comment",
		))
	})

	t.Run("validation: at least one unset option", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &DatabaseUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf(
			"DatabaseUnset",
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
			"Comment",
		))
	})

	t.Run("validation: invalid external volume identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{
			ExternalVolume: new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("DatabaseSet", "ExternalVolume"))
	})

	t.Run("validation: invalid catalog integration identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &DatabaseSet{
			Catalog: new(emptyAccountObjectIdentifier),
		}
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("DatabaseSet", "Catalog"))
	})

	t.Run("validation: invalid NewName identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.NewName = new(emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("AlterDatabaseOptions", "NewName"))
	})

	t.Run("validation: invalid SwapWith identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.SwapWith = new(emptyAccountObjectIdentifier)
		assertOptsInvalidJoinedErrors(t, opts, errInvalidIdentifier("AlterDatabaseOptions", "SwapWith"))
	})

	t.Run("rename", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = new(true)
		opts.NewName = new(randomAccountObjectIdentifier())
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE IF EXISTS %s RENAME TO %s`, opts.name.FullyQualifiedName(), opts.NewName.FullyQualifiedName())
	})

	t.Run("swap with", func(t *testing.T) {
		opts := defaultOpts()
		opts.SwapWith = new(randomAccountObjectIdentifier())
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s SWAP WITH %s`, opts.name.FullyQualifiedName(), opts.SwapWith.FullyQualifiedName())
	})

	t.Run("set", func(t *testing.T) {
		externalVolumeId := randomAccountObjectIdentifier()
		catalogId := randomAccountObjectIdentifier()
		opts := defaultOpts()
		opts.Set = &DatabaseSet{
			DataRetentionTimeInDays:    new(1),
			MaxDataExtensionTimeInDays: new(1),
			ExternalVolume:             &externalVolumeId,
			Catalog:                    &catalogId,
			ReplaceInvalidCharacters:   new(true),
			DefaultDDLCollation:        new("en_US"),
			StorageSerializationPolicy: new(StorageSerializationPolicyCompatible),
			LogLevel:                   new(LogLevelError),
			TraceLevel:                 new(TraceLevelPropagate),
			Comment:                    new("comment"),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s SET DATA_RETENTION_TIME_IN_DAYS = 1, MAX_DATA_EXTENSION_TIME_IN_DAYS = 1, EXTERNAL_VOLUME = %s, CATALOG = %s, REPLACE_INVALID_CHARACTERS = true, DEFAULT_DDL_COLLATION = 'en_US', STORAGE_SERIALIZATION_POLICY = COMPATIBLE, LOG_LEVEL = 'ERROR', TRACE_LEVEL = 'PROPAGATE', COMMENT = 'comment'`, opts.name.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName())
	})

	t.Run("unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &DatabaseUnset{
			DataRetentionTimeInDays:    new(true),
			MaxDataExtensionTimeInDays: new(true),
			ExternalVolume:             new(true),
			Catalog:                    new(true),
			ReplaceInvalidCharacters:   new(true),
			DefaultDDLCollation:        new(true),
			StorageSerializationPolicy: new(true),
			LogLevel:                   new(true),
			TraceLevel:                 new(true),
			Comment:                    new(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, EXTERNAL_VOLUME, CATALOG, REPLACE_INVALID_CHARACTERS, DEFAULT_DDL_COLLATION, STORAGE_SERIALIZATION_POLICY, LOG_LEVEL, TRACE_LEVEL, COMMENT`, opts.name.FullyQualifiedName())
	})

	t.Run("with set tag", func(t *testing.T) {
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		opts := defaultOpts()
		opts.SetTag = []TagAssociation{
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
		opts.UnsetTag = []ObjectIdentifier{
			id,
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s UNSET TAG %s`, opts.name.FullyQualifiedName(), id.FullyQualifiedName())
	})
}

func TestDatabasesAlterReplication(t *testing.T) {
	defaultOpts := func() *AlterDatabaseReplicationOptions {
		return &AlterDatabaseReplicationOptions{
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
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseReplicationOptions", "EnableReplication", "DisableReplication", "Refresh"))
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableReplication = &EnableReplication{}
		opts.DisableReplication = &DisableReplication{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseReplicationOptions", "EnableReplication", "DisableReplication", "Refresh"))
	})

	t.Run("enable replication", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableReplication = &EnableReplication{
			ToAccounts: []AccountIdentifier{
				NewAccountIdentifierFromAccountLocator("account1"),
			},
			IgnoreEditionCheck: new(true),
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
		opts.Refresh = new(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER DATABASE %s REFRESH`, opts.name.FullyQualifiedName())
	})
}

func TestDatabasesAlterFailover(t *testing.T) {
	defaultOpts := func() *AlterDatabaseFailoverOptions {
		return &AlterDatabaseFailoverOptions{
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
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseFailoverOptions", "EnableFailover", "DisableFailover", "Primary"))
	})

	t.Run("validation: exactly one action", func(t *testing.T) {
		opts := defaultOpts()
		opts.EnableFailover = &EnableFailover{}
		opts.DisableFailover = &DisableFailover{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterDatabaseFailoverOptions", "EnableFailover", "DisableFailover", "Primary"))
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
		opts.Primary = new(true)
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
		opts.IfExists = new(true)
		opts.Cascade = new(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE IF EXISTS %s CASCADE`, opts.name.FullyQualifiedName())
	})

	t.Run("all options - restrict", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = new(true)
		opts.Restrict = new(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP DATABASE IF EXISTS %s RESTRICT`, opts.name.FullyQualifiedName())
	})

	t.Run("validation: cascade and restrict set together", func(t *testing.T) {
		opts := defaultOpts()
		opts.Cascade = new(true)
		opts.Restrict = new(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("DropDatabaseOptions", "Cascade", "Restrict"))
	})
}

func TestDatabasesUndrop(t *testing.T) {
	defaultOpts := func() *undropDatabaseOptions {
		return &undropDatabaseOptions{
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
	defaultOpts := func() *ShowDatabasesOptions {
		return &ShowDatabasesOptions{}
	}

	t.Run("without show options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASES`)
	})

	t.Run("terse", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = new(true)
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE DATABASES`)
	})

	t.Run("history", func(t *testing.T) {
		opts := defaultOpts()
		opts.History = new(true)
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASES HISTORY`)
	})

	t.Run("like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: new("db1"),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW DATABASES LIKE 'db1'`)
	})

	t.Run("complete", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = new(true)
		opts.History = new(true)
		opts.Like = &Like{
			Pattern: new("db2"),
		}
		opts.LimitFrom = &LimitFrom{
			Rows: new(1),
			From: new("db1"),
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
