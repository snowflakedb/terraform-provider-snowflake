package sdk

import (
	"errors"
	"testing"
	"time"
)

func init() {
	id := databasesTestIdAccountObjectIdentifier
	externalVolumeId := randomAccountObjectIdentifier()
	catalogId := randomAccountObjectIdentifier()
	tagId := randomAccountObjectIdentifier()
	cloneSourceId := randomAccountObjectIdentifier()
	fromShareId := randomExternalObjectIdentifier()
	primaryDatabaseId := randomExternalObjectIdentifier()
	listingGlobalName := "GZ1M7Z91WTX"
	renameTarget := randomAccountObjectIdentifier()
	swapWithId := randomAccountObjectIdentifier()
	setTagId1 := randomSchemaObjectIdentifier()
	setTagId2 := randomSchemaObjectIdentifierInSchema(setTagId1.SchemaId())
	unsetTagId := randomSchemaObjectIdentifier()
	replicationAccount := NewAccountIdentifierFromAccountLocator("account1")
	cloneCompleteSource := NewAccountObjectIdentifier("db1")
	cloneTimestamp := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

	databasesTests.Create.
		withExpectedSqlf(
			case_Databases_sql_Create_basic,
			`CREATE DATABASE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_Create_all,
			func(opts *CreateDatabaseOptions) {
				opts.IfNotExists = new(true)
				opts.Transient = new(true)
				opts.DataRetentionTimeInDays = new(1)
				opts.MaxDataExtensionTimeInDays = new(1)
				opts.ExternalVolume = &externalVolumeId
				opts.Catalog = &catalogId
				opts.ReplaceInvalidCharacters = new(true)
				opts.DefaultDdlCollation = new("en_US")
				opts.DefaultNotebookComputePoolCpu = new("CPU_X64_S")
				opts.DefaultNotebookComputePoolGpu = new("GPU_NV_S")
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
				opts.Tag = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			`CREATE TRANSIENT DATABASE IF NOT EXISTS %s DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 EXTERNAL_VOLUME = %s CATALOG = %s REPLACE_INVALID_CHARACTERS = true DEFAULT_DDL_COLLATION = 'en_US' DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU = 'CPU_X64_S' DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU = 'GPU_NV_S' STORAGE_SERIALIZATION_POLICY = COMPATIBLE LOG_LEVEL = 'INFO' TRACE_LEVEL = 'PROPAGATE' SUSPEND_TASK_AFTER_NUM_FAILURES = 10 TASK_AUTO_RETRY_ATTEMPTS = 10 USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM USER_TASK_TIMEOUT_MS = 12000 USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30 QUOTED_IDENTIFIERS_IGNORE_CASE = true ENABLE_CONSOLE_OUTPUT = true COMMENT = 'comment' TAG (%s = 'v1')`,
			id.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateDatabaseOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE DATABASE %s`, id.FullyQualifiedName(),
		)

	databasesTests.Clone.
		withDefaultOpts(func() *CloneDatabaseOptions {
			return &CloneDatabaseOptions{
				name:  id,
				Clone: Clone{SourceObject: cloneSourceId},
			}
		}).
		withAdditionalValidationCase(
			"validation_Clone_AtAndBeforeSet",
			func(opts *CloneDatabaseOptions) {
				opts.Clone.At = &TimeTravel{Timestamp: new(cloneTimestamp)}
				opts.Clone.Before = &TimeTravel{Offset: new(123)}
			},
			errors.New("only one of AT or BEFORE can be set"),
		).
		withAdditionalValidationCase(
			"validation_Clone_TimeTravelExactlyOne",
			func(opts *CloneDatabaseOptions) {
				opts.Clone.At = &TimeTravel{
					Timestamp: new(cloneTimestamp),
					Offset:    new(123),
				}
			},
			errors.New("exactly one of TIMESTAMP, OFFSET or STATEMENT can be set"),
		).
		withExpectedSqlf(
			case_Databases_sql_Clone_basic,
			`CREATE DATABASE %s CLONE %s`, id.FullyQualifiedName(), cloneSourceId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Clone_all",
			func(opts *CloneDatabaseOptions) {
				opts.IfNotExists = new(true)
				opts.Clone = Clone{
					SourceObject: cloneCompleteSource,
					At:           &TimeTravel{Timestamp: &cloneTimestamp},
				}
			},
			`CREATE DATABASE IF NOT EXISTS %s CLONE "db1" AT (TIMESTAMP => '2021-01-01 00:00:00 +0000 UTC')`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Clone_orReplace",
			func(opts *CloneDatabaseOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE DATABASE %s CLONE %s`, id.FullyQualifiedName(), cloneSourceId.FullyQualifiedName(),
		)

	databasesTests.CreateShared.
		withDefaultOpts(func() *CreateSharedDatabaseOptions {
			return &CreateSharedDatabaseOptions{
				name:      id,
				FromShare: fromShareId,
			}
		}).
		withModifyAndExpectedSqlf(
			case_Databases_sql_CreateShared_basic,
			func(opts *CreateSharedDatabaseOptions) {
				opts.Transient = new(true)
			},
			`CREATE TRANSIENT DATABASE %s FROM SHARE %s`,
			id.FullyQualifiedName(), fromShareId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_CreateShared_all,
			func(opts *CreateSharedDatabaseOptions) {
				opts.IfNotExists = new(true)
				opts.ExternalVolume = &externalVolumeId
				opts.Catalog = &catalogId
				opts.ReplaceInvalidCharacters = new(true)
				opts.DefaultDdlCollation = new("en_US")
				opts.DefaultNotebookComputePoolCpu = new("CPU_X64_S")
				opts.DefaultNotebookComputePoolGpu = new("GPU_NV_S")
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
				opts.Tag = []TagAssociation{{Name: tagId, Value: "v1"}}
			},
			`CREATE DATABASE IF NOT EXISTS %s FROM SHARE %s EXTERNAL_VOLUME = %s CATALOG = %s REPLACE_INVALID_CHARACTERS = true DEFAULT_DDL_COLLATION = 'en_US' DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU = 'CPU_X64_S' DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU = 'GPU_NV_S' STORAGE_SERIALIZATION_POLICY = COMPATIBLE LOG_LEVEL = 'INFO' TRACE_LEVEL = 'PROPAGATE' SUSPEND_TASK_AFTER_NUM_FAILURES = 10 TASK_AUTO_RETRY_ATTEMPTS = 10 USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM USER_TASK_TIMEOUT_MS = 12000 USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30 QUOTED_IDENTIFIERS_IGNORE_CASE = true ENABLE_CONSOLE_OUTPUT = true COMMENT = 'comment' TAG (%s = 'v1')`,
			id.FullyQualifiedName(), fromShareId.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName(), tagId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateShared_orReplace",
			func(opts *CreateSharedDatabaseOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE DATABASE %s FROM SHARE %s`,
			id.FullyQualifiedName(), fromShareId.FullyQualifiedName(),
		)

	databasesTests.CreateSecondary.
		withDefaultOpts(func() *CreateSecondaryDatabaseOptions {
			return &CreateSecondaryDatabaseOptions{
				name:            id,
				PrimaryDatabase: primaryDatabaseId,
			}
		}).
		withExpectedSqlf(
			case_Databases_sql_CreateSecondary_basic,
			`CREATE DATABASE %s AS REPLICA OF %s`,
			id.FullyQualifiedName(), primaryDatabaseId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_CreateSecondary_all,
			func(opts *CreateSecondaryDatabaseOptions) {
				opts.IfNotExists = new(true)
				opts.Transient = new(true)
				opts.DataRetentionTimeInDays = new(1)
				opts.MaxDataExtensionTimeInDays = new(1)
				opts.ExternalVolume = &externalVolumeId
				opts.Catalog = &catalogId
				opts.ReplaceInvalidCharacters = new(true)
				opts.DefaultDdlCollation = new("en_US")
				opts.DefaultNotebookComputePoolCpu = new("CPU_X64_S")
				opts.DefaultNotebookComputePoolGpu = new("GPU_NV_S")
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
			},
			`CREATE TRANSIENT DATABASE IF NOT EXISTS %s AS REPLICA OF %s DATA_RETENTION_TIME_IN_DAYS = 1 MAX_DATA_EXTENSION_TIME_IN_DAYS = 1 EXTERNAL_VOLUME = %s CATALOG = %s REPLACE_INVALID_CHARACTERS = true DEFAULT_DDL_COLLATION = 'en_US' DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU = 'CPU_X64_S' DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU = 'GPU_NV_S' STORAGE_SERIALIZATION_POLICY = COMPATIBLE LOG_LEVEL = 'INFO' TRACE_LEVEL = 'PROPAGATE' SUSPEND_TASK_AFTER_NUM_FAILURES = 10 TASK_AUTO_RETRY_ATTEMPTS = 10 USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE = MEDIUM USER_TASK_TIMEOUT_MS = 12000 USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS = 30 QUOTED_IDENTIFIERS_IGNORE_CASE = true ENABLE_CONSOLE_OUTPUT = true COMMENT = 'comment'`,
			id.FullyQualifiedName(), primaryDatabaseId.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_CreateSecondary_orReplace",
			func(opts *CreateSecondaryDatabaseOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE DATABASE %s AS REPLICA OF %s`,
			id.FullyQualifiedName(), primaryDatabaseId.FullyQualifiedName(),
		)

	databasesTests.CreateFromListing.
		withDefaultOpts(func() *CreateFromListingDatabaseOptions {
			return &CreateFromListingDatabaseOptions{
				name:        id,
				FromListing: listingGlobalName,
			}
		}).
		withAdditionalValidationCase(
			"validation_CreateFromListing_emptyListingGlobalName",
			func(opts *CreateFromListingDatabaseOptions) { opts.FromListing = "" },
			NewError("CreateFromListingDatabaseOptions: listing global name must not be empty"),
		).
		withExpectedSqlf(
			case_Databases_sql_CreateFromListing_basic,
			`CREATE DATABASE %s FROM LISTING '%s'`, id.FullyQualifiedName(), listingGlobalName,
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_CreateFromListing_all,
			func(opts *CreateFromListingDatabaseOptions) {},
			`CREATE DATABASE %s FROM LISTING '%s'`, id.FullyQualifiedName(), listingGlobalName,
		)

	databasesTests.CreateCatalogLinked.
		withDefaultOpts(func() *CreateCatalogLinkedDatabaseOptions {
			return &CreateCatalogLinkedDatabaseOptions{
				name: id,
				LinkedCatalog: LinkedCatalog{
					Catalog: catalogId,
				},
			}
		}).
		withExpectedSqlf(
			case_Databases_sql_CreateCatalogLinked_basic,
			`CREATE DATABASE %s LINKED_CATALOG = (CATALOG = '\"%s\"')`, id.FullyQualifiedName(), catalogId.Name(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_CreateCatalogLinked_all,
			func(opts *CreateCatalogLinkedDatabaseOptions) {
				opts.LinkedCatalog = LinkedCatalog{
					Catalog:                   catalogId,
					AllowedNamespaces:         []StringListItemWrapper{{Value: "ns1"}, {Value: "ns2"}},
					BlockedNamespaces:         []StringListItemWrapper{{Value: "ns3"}},
					AllowedWriteOperations:    new(CatalogLinkedDatabaseAllowedWriteOperationsAll),
					NamespaceMode:             new(CatalogLinkedDatabaseNamespaceModeFlattenNestedNamespace),
					NamespaceFlattenDelimiter: new("-"),
					SyncIntervalSeconds:       new(60),
				}
				opts.ExternalVolume = &externalVolumeId
				opts.Comment = new("comment")
				opts.Tag = []TagAssociation{{Name: tagId, Value: "v1"}}
				opts.CatalogCaseSensitivity = new(DatabaseCatalogCaseSensitivityCaseInsensitive)
			},
			`CREATE DATABASE %s LINKED_CATALOG = (CATALOG = '\"%s\"', ALLOWED_NAMESPACES = ('ns1', 'ns2'), BLOCKED_NAMESPACES = ('ns3'), ALLOWED_WRITE_OPERATIONS = ALL, NAMESPACE_MODE = FLATTEN_NESTED_NAMESPACE, NAMESPACE_FLATTEN_DELIMITER = '-', SYNC_INTERVAL_SECONDS = 60) EXTERNAL_VOLUME = '\"%s\"' COMMENT = 'comment' TAG (%s = 'v1') CATALOG_CASE_SENSITIVITY = CASE_INSENSITIVE`,
			id.FullyQualifiedName(), catalogId.Name(), externalVolumeId.Name(), tagId.FullyQualifiedName(),
		)

	databasesTests.Alter.
		withModifyAndExpectedSqlf(
			case_Databases_sql_Alter_RenameTo,
			func(opts *AlterDatabaseOptions) {
				opts.IfExists = new(true)
				opts.RenameTo = &renameTarget
			},
			`ALTER DATABASE IF EXISTS %s RENAME TO %s`,
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_Alter_Set,
			func(opts *AlterDatabaseOptions) {
				opts.Set = &DatabaseSet{
					DataRetentionTimeInDays:       new(1),
					MaxDataExtensionTimeInDays:    new(1),
					ExternalVolume:                &externalVolumeId,
					Catalog:                       &catalogId,
					ReplaceInvalidCharacters:      new(true),
					DefaultDdlCollation:           new("en_US"),
					DefaultNotebookComputePoolCpu: new("CPU_X64_S"),
					DefaultNotebookComputePoolGpu: new("GPU_NV_S"),
					StorageSerializationPolicy:    new(StorageSerializationPolicyCompatible),
					LogLevel:                      new(LogLevelError),
					TraceLevel:                    new(TraceLevelPropagate),
					Comment:                       new("comment"),
				}
			},
			`ALTER DATABASE %s SET DATA_RETENTION_TIME_IN_DAYS = 1, MAX_DATA_EXTENSION_TIME_IN_DAYS = 1, EXTERNAL_VOLUME = %s, CATALOG = %s, REPLACE_INVALID_CHARACTERS = true, DEFAULT_DDL_COLLATION = 'en_US', DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU = 'CPU_X64_S', DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU = 'GPU_NV_S', STORAGE_SERIALIZATION_POLICY = COMPATIBLE, LOG_LEVEL = 'ERROR', TRACE_LEVEL = 'PROPAGATE', COMMENT = 'comment'`,
			id.FullyQualifiedName(), externalVolumeId.FullyQualifiedName(), catalogId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_Alter_Unset,
			func(opts *AlterDatabaseOptions) {
				opts.Unset = &DatabaseUnset{
					DataRetentionTimeInDays:       new(true),
					MaxDataExtensionTimeInDays:    new(true),
					ExternalVolume:                new(true),
					Catalog:                       new(true),
					ReplaceInvalidCharacters:      new(true),
					DefaultDdlCollation:           new(true),
					DefaultNotebookComputePoolCpu: new(true),
					DefaultNotebookComputePoolGpu: new(true),
					StorageSerializationPolicy:    new(true),
					LogLevel:                      new(true),
					TraceLevel:                    new(true),
					Comment:                       new(true),
				}
			},
			`ALTER DATABASE %s UNSET DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS, EXTERNAL_VOLUME, CATALOG, REPLACE_INVALID_CHARACTERS, DEFAULT_DDL_COLLATION, DEFAULT_NOTEBOOK_COMPUTE_POOL_CPU, DEFAULT_NOTEBOOK_COMPUTE_POOL_GPU, STORAGE_SERIALIZATION_POLICY, LOG_LEVEL, TRACE_LEVEL, COMMENT`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_Alter_SwapWith,
			func(opts *AlterDatabaseOptions) { opts.SwapWith = &swapWithId },
			`ALTER DATABASE %s SWAP WITH %s`,
			id.FullyQualifiedName(), swapWithId.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_Alter_SetTags,
			func(opts *AlterDatabaseOptions) {
				opts.SetTags = []TagAssociation{
					{Name: setTagId1, Value: "v1"},
					{Name: setTagId2, Value: "v2"},
				}
			},
			`ALTER DATABASE %s SET TAG %s = 'v1', %s = 'v2'`,
			id.FullyQualifiedName(), setTagId1.FullyQualifiedName(), setTagId2.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_Alter_UnsetTags,
			func(opts *AlterDatabaseOptions) {
				opts.UnsetTags = []ObjectIdentifier{unsetTagId}
			},
			`ALTER DATABASE %s UNSET TAG %s`,
			id.FullyQualifiedName(), unsetTagId.FullyQualifiedName(),
		)

	databasesTests.AlterReplication.
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterReplication_EnableReplication,
			func(opts *AlterReplicationDatabaseOptions) {
				opts.EnableReplication = &EnableReplication{
					ToAccounts:         []AccountIdentifier{replicationAccount},
					IgnoreEditionCheck: new(true),
				}
			},
			`ALTER DATABASE %s ENABLE REPLICATION TO ACCOUNTS "account1" IGNORE EDITION CHECK`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterReplication_DisableReplication,
			func(opts *AlterReplicationDatabaseOptions) {
				opts.DisableReplication = &DisableReplication{
					ToAccounts: []AccountIdentifier{replicationAccount},
				}
			},
			`ALTER DATABASE %s DISABLE REPLICATION TO ACCOUNTS "account1"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterReplication_Refresh,
			func(opts *AlterReplicationDatabaseOptions) { opts.Refresh = new(true) },
			`ALTER DATABASE %s REFRESH`,
			id.FullyQualifiedName(),
		)

	databasesTests.AlterFailover.
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterFailover_EnableFailover,
			func(opts *AlterFailoverDatabaseOptions) {
				opts.EnableFailover = &EnableFailover{
					ToAccounts: []AccountIdentifier{replicationAccount},
				}
			},
			`ALTER DATABASE %s ENABLE FAILOVER TO ACCOUNTS "account1"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterFailover_DisableFailover,
			func(opts *AlterFailoverDatabaseOptions) {
				opts.DisableFailover = &DisableFailover{
					ToAccounts: []AccountIdentifier{replicationAccount},
				}
			},
			`ALTER DATABASE %s DISABLE FAILOVER TO ACCOUNTS "account1"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterFailover_Primary,
			func(opts *AlterFailoverDatabaseOptions) { opts.Primary = new(true) },
			`ALTER DATABASE %s PRIMARY`,
			id.FullyQualifiedName(),
		)

	databasesTests.AlterCatalogLinked.
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterCatalogLinked_AddToAllowedNamespaces,
			func(opts *AlterCatalogLinkedDatabaseOptions) {
				opts.IfExists = new(true)
				opts.AddToAllowedNamespaces = &AddToAllowedNamespaces{
					Namespaces: []StringListItemWrapper{{Value: "ns1"}, {Value: "ns2"}},
				}
			},
			`ALTER DATABASE IF EXISTS %s UPDATE LINKED_CATALOG ADD ('ns1', 'ns2') TO ALLOWED_NAMESPACES`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterCatalogLinked_RemoveFromAllowedNamespaces,
			func(opts *AlterCatalogLinkedDatabaseOptions) {
				opts.RemoveFromAllowedNamespaces = &RemoveFromAllowedNamespaces{
					Namespaces: []StringListItemWrapper{{Value: "ns1"}},
				}
			},
			`ALTER DATABASE %s UPDATE LINKED_CATALOG REMOVE ('ns1') FROM ALLOWED_NAMESPACES`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterCatalogLinked_UnsetAllowedNamespaces,
			func(opts *AlterCatalogLinkedDatabaseOptions) { opts.UnsetAllowedNamespaces = new(true) },
			`ALTER DATABASE %s UPDATE LINKED_CATALOG UNSET ALLOWED_NAMESPACES`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterCatalogLinked_AddToBlockedNamespaces,
			func(opts *AlterCatalogLinkedDatabaseOptions) {
				opts.AddToBlockedNamespaces = &AddToBlockedNamespaces{
					Namespaces: []StringListItemWrapper{{Value: "ns3"}},
				}
			},
			`ALTER DATABASE %s UPDATE LINKED_CATALOG ADD ('ns3') TO BLOCKED_NAMESPACES`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterCatalogLinked_RemoveFromBlockedNamespaces,
			func(opts *AlterCatalogLinkedDatabaseOptions) {
				opts.RemoveFromBlockedNamespaces = &RemoveFromBlockedNamespaces{
					Namespaces: []StringListItemWrapper{{Value: "ns3"}},
				}
			},
			`ALTER DATABASE %s UPDATE LINKED_CATALOG REMOVE ('ns3') FROM BLOCKED_NAMESPACES`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterCatalogLinked_UnsetBlockedNamespaces,
			func(opts *AlterCatalogLinkedDatabaseOptions) { opts.UnsetBlockedNamespaces = new(true) },
			`ALTER DATABASE %s UPDATE LINKED_CATALOG UNSET BLOCKED_NAMESPACES`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_AlterCatalogLinked_Set,
			func(opts *AlterCatalogLinkedDatabaseOptions) {
				opts.Set = &CatalogLinkedDatabaseSet{
					SyncIntervalSeconds:    new(120),
					AllowedWriteOperations: new(CatalogLinkedDatabaseAllowedWriteOperationsNone),
				}
			},
			`ALTER DATABASE %s UPDATE LINKED_CATALOG SET SYNC_INTERVAL_SECONDS = 120, ALLOWED_WRITE_OPERATIONS = NONE`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterCatalogLinked_Set_syncIntervalSeconds",
			func(opts *AlterCatalogLinkedDatabaseOptions) {
				opts.Set = &CatalogLinkedDatabaseSet{
					SyncIntervalSeconds: new(120),
				}
			},
			`ALTER DATABASE %s UPDATE LINKED_CATALOG SET SYNC_INTERVAL_SECONDS = 120`,
			id.FullyQualifiedName(),
		)

	databasesTests.Drop.
		withExpectedSqlf(
			case_Databases_sql_Drop_basic,
			`DROP DATABASE %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_Drop_all,
			func(opts *DropDatabaseOptions) {
				opts.IfExists = new(true)
				opts.Cascade = new(true)
			},
			`DROP DATABASE IF EXISTS %s CASCADE`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Drop_all_restrict",
			func(opts *DropDatabaseOptions) {
				opts.IfExists = new(true)
				opts.Restrict = new(true)
			},
			`DROP DATABASE IF EXISTS %s RESTRICT`,
			id.FullyQualifiedName(),
		)

	databasesTests.Undrop.
		withExpectedSqlf(
			case_Databases_sql_Undrop_basic,
			`UNDROP DATABASE %s`, id.FullyQualifiedName(),
		)

	databasesTests.Show.
		withExpectedSql(
			case_Databases_sql_Show_basic,
			`SHOW DATABASES`,
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_Show_all,
			func(opts *ShowDatabaseOptions) {
				opts.Terse = new(true)
				opts.History = new(true)
				opts.Like = &Like{Pattern: new("db2")}
				opts.StartsWith = new("prefix")
				opts.Limit = &LimitFrom{Rows: new(1), From: new("db1")}
			},
			`SHOW TERSE DATABASES HISTORY LIKE 'db2' STARTS WITH 'prefix' LIMIT 1 FROM 'db1'`,
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_Show_Like,
			func(opts *ShowDatabaseOptions) { opts.Like = &Like{Pattern: new("db1")} },
			`SHOW DATABASES LIKE 'db1'`,
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_Show_StartsWith,
			func(opts *ShowDatabaseOptions) { opts.StartsWith = new("prefix") },
			`SHOW DATABASES STARTS WITH 'prefix'`,
		).
		withModifyAndExpectedSqlf(
			case_Databases_sql_Show_Limit,
			func(opts *ShowDatabaseOptions) {
				opts.Limit = &LimitFrom{Rows: new(1), From: new("db1")}
			},
			`SHOW DATABASES LIMIT 1 FROM 'db1'`,
		).
		withAdditionalSqlCasef(
			"sql_Show_Terse",
			func(opts *ShowDatabaseOptions) { opts.Terse = new(true) },
			`SHOW TERSE DATABASES`,
		).
		withAdditionalSqlCasef(
			"sql_Show_History",
			func(opts *ShowDatabaseOptions) { opts.History = new(true) },
			`SHOW DATABASES HISTORY`,
		)
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
		assertOptsValidAndSqlEqualsf(t, opts, `DESCRIBE DATABASE %s`, opts.name.FullyQualifiedName())
	})
}
