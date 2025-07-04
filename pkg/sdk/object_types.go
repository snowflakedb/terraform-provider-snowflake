package sdk

import (
	"fmt"
	"slices"
	"strings"
)

// Object bundles together the object type and name. Its used for DDL statements.
type Object struct {
	ObjectType ObjectType       `ddl:"keyword"`
	Name       ObjectIdentifier `ddl:"identifier"`
}

// ObjectType is the type of object.
type ObjectType string

const (
	ObjectTypeAccount              ObjectType = "ACCOUNT"
	ObjectTypeManagedAccount       ObjectType = "MANAGED ACCOUNT"
	ObjectTypeUser                 ObjectType = "USER"
	ObjectTypeDatabaseRole         ObjectType = "DATABASE ROLE"
	ObjectTypeDataset              ObjectType = "DATASET"
	ObjectTypeRole                 ObjectType = "ROLE"
	ObjectTypeIntegration          ObjectType = "INTEGRATION"
	ObjectTypeNetworkPolicy        ObjectType = "NETWORK POLICY"
	ObjectTypePasswordPolicy       ObjectType = "PASSWORD POLICY"
	ObjectTypeSessionPolicy        ObjectType = "SESSION POLICY"
	ObjectTypePrivacyPolicy        ObjectType = "PRIVACY POLICY"
	ObjectTypeReplicationGroup     ObjectType = "REPLICATION GROUP"
	ObjectTypeFailoverGroup        ObjectType = "FAILOVER GROUP"
	ObjectTypeConnection           ObjectType = "CONNECTION"
	ObjectTypeParameter            ObjectType = "PARAMETER"
	ObjectTypeWarehouse            ObjectType = "WAREHOUSE"
	ObjectTypeResourceMonitor      ObjectType = "RESOURCE MONITOR"
	ObjectTypeDatabase             ObjectType = "DATABASE"
	ObjectTypeSchema               ObjectType = "SCHEMA"
	ObjectTypeShare                ObjectType = "SHARE"
	ObjectTypeTable                ObjectType = "TABLE"
	ObjectTypeDynamicTable         ObjectType = "DYNAMIC TABLE"
	ObjectTypeCortexSearchService  ObjectType = "CORTEX SEARCH SERVICE"
	ObjectTypeExternalTable        ObjectType = "EXTERNAL TABLE"
	ObjectTypeEventTable           ObjectType = "EVENT TABLE"
	ObjectTypeView                 ObjectType = "VIEW"
	ObjectTypeMaterializedView     ObjectType = "MATERIALIZED VIEW"
	ObjectTypeSequence             ObjectType = "SEQUENCE"
	ObjectTypeSnapshot             ObjectType = "SNAPSHOT"
	ObjectTypeFunction             ObjectType = "FUNCTION"
	ObjectTypeExternalFunction     ObjectType = "EXTERNAL FUNCTION"
	ObjectTypeProcedure            ObjectType = "PROCEDURE"
	ObjectTypeStream               ObjectType = "STREAM"
	ObjectTypeTask                 ObjectType = "TASK"
	ObjectTypeMaskingPolicy        ObjectType = "MASKING POLICY"
	ObjectTypeRowAccessPolicy      ObjectType = "ROW ACCESS POLICY"
	ObjectTypeTag                  ObjectType = "TAG"
	ObjectTypeSecret               ObjectType = "SECRET"
	ObjectTypeStage                ObjectType = "STAGE"
	ObjectTypeFileFormat           ObjectType = "FILE FORMAT"
	ObjectTypePipe                 ObjectType = "PIPE"
	ObjectTypeAlert                ObjectType = "ALERT"
	ObjectTypeBudget               ObjectType = "SNOWFLAKE.CORE.BUDGET"
	ObjectTypeClassification       ObjectType = "SNOWFLAKE.ML.CLASSIFICATION"
	ObjectTypeApplication          ObjectType = "APPLICATION"
	ObjectTypeApplicationPackage   ObjectType = "APPLICATION PACKAGE"
	ObjectTypeApplicationRole      ObjectType = "APPLICATION ROLE"
	ObjectTypeStreamlit            ObjectType = "STREAMLIT"
	ObjectTypeColumn               ObjectType = "COLUMN"
	ObjectTypeIcebergTable         ObjectType = "ICEBERG TABLE"
	ObjectTypeExternalVolume       ObjectType = "EXTERNAL VOLUME"
	ObjectTypeNetworkRule          ObjectType = "NETWORK RULE"
	ObjectTypeNotebook             ObjectType = "NOTEBOOK"
	ObjectTypePackagesPolicy       ObjectType = "PACKAGES POLICY"
	ObjectTypeComputePool          ObjectType = "COMPUTE POOL"
	ObjectTypeAggregationPolicy    ObjectType = "AGGREGATION POLICY"
	ObjectTypeAuthenticationPolicy ObjectType = "AUTHENTICATION POLICY"
	ObjectTypeHybridTable          ObjectType = "HYBRID TABLE"
	ObjectTypeImageRepository      ObjectType = "IMAGE REPOSITORY"
	ObjectTypeProjectionPolicy     ObjectType = "PROJECTION POLICY"
	ObjectTypeDataMetricFunction   ObjectType = "DATA METRIC FUNCTION"
	ObjectTypeGitRepository        ObjectType = "GIT REPOSITORY"
	ObjectTypeModel                ObjectType = "MODEL"
	ObjectTypeService              ObjectType = "SERVICE"
	// ObjectTypeProgrammaticAccessToken is a pseudo-object, as it does not support the usual operations in Snowflake, but it is handled by user functions.
	// Programmatic access tokens do not have grants and cannot be tagged.
	ObjectTypeProgrammaticAccessToken ObjectType = "PROGRAMMATIC ACCESS TOKEN"
)

func (o ObjectType) String() string {
	return string(o)
}

func (o ObjectType) IsWithArguments() bool {
	return slices.Contains([]ObjectType{ObjectTypeExternalFunction, ObjectTypeFunction, ObjectTypeProcedure}, o)
}

var allObjectTypes = []ObjectType{
	ObjectTypeAccount,
	ObjectTypeManagedAccount,
	ObjectTypeUser,
	ObjectTypeDatabaseRole,
	ObjectTypeDataset,
	ObjectTypeRole,
	ObjectTypeIntegration,
	ObjectTypeNetworkPolicy,
	ObjectTypePasswordPolicy,
	ObjectTypeSessionPolicy,
	ObjectTypePrivacyPolicy,
	ObjectTypeReplicationGroup,
	ObjectTypeFailoverGroup,
	ObjectTypeConnection,
	ObjectTypeParameter,
	ObjectTypeWarehouse,
	ObjectTypeResourceMonitor,
	ObjectTypeDatabase,
	ObjectTypeSchema,
	ObjectTypeShare,
	ObjectTypeTable,
	ObjectTypeDynamicTable,
	ObjectTypeCortexSearchService,
	ObjectTypeExternalTable,
	ObjectTypeEventTable,
	ObjectTypeView,
	ObjectTypeMaterializedView,
	ObjectTypeSequence,
	ObjectTypeSnapshot,
	ObjectTypeFunction,
	ObjectTypeExternalFunction,
	ObjectTypeProcedure,
	ObjectTypeStream,
	ObjectTypeTask,
	ObjectTypeMaskingPolicy,
	ObjectTypeRowAccessPolicy,
	ObjectTypeTag,
	ObjectTypeSecret,
	ObjectTypeStage,
	ObjectTypeFileFormat,
	ObjectTypePipe,
	ObjectTypeAlert,
	ObjectTypeBudget,
	ObjectTypeClassification,
	ObjectTypeApplication,
	ObjectTypeApplicationPackage,
	ObjectTypeApplicationRole,
	ObjectTypeStreamlit,
	ObjectTypeColumn,
	ObjectTypeIcebergTable,
	ObjectTypeExternalVolume,
	ObjectTypeNetworkRule,
	ObjectTypeNotebook,
	ObjectTypePackagesPolicy,
	ObjectTypeComputePool,
	ObjectTypeAggregationPolicy,
	ObjectTypeAuthenticationPolicy,
	ObjectTypeHybridTable,
	ObjectTypeImageRepository,
	ObjectTypeProjectionPolicy,
	ObjectTypeDataMetricFunction,
	ObjectTypeGitRepository,
	ObjectTypeModel,
	ObjectTypeService,
	ObjectTypeProgrammaticAccessToken,
}

// TODO(SNOW-1834370): use ToObjectType in other places with type conversion (instead of sdk.ObjectType)
func ToObjectType(s string) (ObjectType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allObjectTypes, ObjectType(s)) {
		return "", fmt.Errorf("invalid object type: %s", s)
	}
	return ObjectType(s), nil
}

func objectTypeSingularToPluralMap() map[ObjectType]PluralObjectType {
	return map[ObjectType]PluralObjectType{
		ObjectTypeAccount:                 PluralObjectTypeAccounts,
		ObjectTypeManagedAccount:          PluralObjectTypeManagedAccounts,
		ObjectTypeUser:                    PluralObjectTypeUsers,
		ObjectTypeDatabaseRole:            PluralObjectTypeDatabaseRoles,
		ObjectTypeDataset:                 PluralObjectTypeDatasets,
		ObjectTypeRole:                    PluralObjectTypeRoles,
		ObjectTypeIntegration:             PluralObjectTypeIntegrations,
		ObjectTypeNetworkPolicy:           PluralObjectTypeNetworkPolicies,
		ObjectTypePasswordPolicy:          PluralObjectTypePasswordPolicies,
		ObjectTypeSessionPolicy:           PluralObjectTypeSessionPolicies,
		ObjectTypePrivacyPolicy:           PluralObjectTypePrivacyPolicies,
		ObjectTypeReplicationGroup:        PluralObjectTypeReplicationGroups,
		ObjectTypeFailoverGroup:           PluralObjectTypeFailoverGroups,
		ObjectTypeConnection:              PluralObjectTypeConnections,
		ObjectTypeParameter:               PluralObjectTypeParameters,
		ObjectTypeWarehouse:               PluralObjectTypeWarehouses,
		ObjectTypeResourceMonitor:         PluralObjectTypeResourceMonitors,
		ObjectTypeDatabase:                PluralObjectTypeDatabases,
		ObjectTypeSchema:                  PluralObjectTypeSchemas,
		ObjectTypeShare:                   PluralObjectTypeShares,
		ObjectTypeTable:                   PluralObjectTypeTables,
		ObjectTypeDynamicTable:            PluralObjectTypeDynamicTables,
		ObjectTypeCortexSearchService:     PluralObjectTypeCortexSearchServices,
		ObjectTypeExternalTable:           PluralObjectTypeExternalTables,
		ObjectTypeEventTable:              PluralObjectTypeEventTables,
		ObjectTypeView:                    PluralObjectTypeViews,
		ObjectTypeMaterializedView:        PluralObjectTypeMaterializedViews,
		ObjectTypeSequence:                PluralObjectTypeSequences,
		ObjectTypeSnapshot:                PluralObjectTypeSnapshots,
		ObjectTypeFunction:                PluralObjectTypeFunctions,
		ObjectTypeExternalFunction:        PluralObjectTypeExternalFunctions,
		ObjectTypeProcedure:               PluralObjectTypeProcedures,
		ObjectTypeStream:                  PluralObjectTypeStreams,
		ObjectTypeTask:                    PluralObjectTypeTasks,
		ObjectTypeMaskingPolicy:           PluralObjectTypeMaskingPolicies,
		ObjectTypeRowAccessPolicy:         PluralObjectTypeRowAccessPolicies,
		ObjectTypeTag:                     PluralObjectTypeTags,
		ObjectTypeSecret:                  PluralObjectTypeSecrets,
		ObjectTypeStage:                   PluralObjectTypeStages,
		ObjectTypeFileFormat:              PluralObjectTypeFileFormats,
		ObjectTypePipe:                    PluralObjectTypePipes,
		ObjectTypeAlert:                   PluralObjectTypeAlerts,
		ObjectTypeBudget:                  PluralObjectTypeBudgets,
		ObjectTypeClassification:          PluralObjectTypeClassifications,
		ObjectTypeApplication:             PluralObjectTypeApplications,
		ObjectTypeApplicationPackage:      PluralObjectTypeApplicationPackages,
		ObjectTypeApplicationRole:         PluralObjectTypeApplicationRoles,
		ObjectTypeStreamlit:               PluralObjectTypeStreamlits,
		ObjectTypeIcebergTable:            PluralObjectTypeIcebergTables,
		ObjectTypeExternalVolume:          PluralObjectTypeExternalVolumes,
		ObjectTypeNetworkRule:             PluralObjectTypeNetworkRules,
		ObjectTypeNotebook:                PluralObjectTypeNotebooks,
		ObjectTypePackagesPolicy:          PluralObjectTypePackagesPolicies,
		ObjectTypeComputePool:             PluralObjectTypeComputePool,
		ObjectTypeAggregationPolicy:       PluralObjectTypeAggregationPolicies,
		ObjectTypeAuthenticationPolicy:    PluralObjectTypeAuthenticationPolicies,
		ObjectTypeHybridTable:             PluralObjectTypeHybridTables,
		ObjectTypeImageRepository:         PluralObjectTypeImageRepositories,
		ObjectTypeProjectionPolicy:        PluralObjectTypeProjectionPolicies,
		ObjectTypeDataMetricFunction:      PluralObjectTypeDataMetricFunctions,
		ObjectTypeGitRepository:           PluralObjectTypeGitRepositories,
		ObjectTypeModel:                   PluralObjectTypeModels,
		ObjectTypeService:                 PluralObjectTypeServices,
		ObjectTypeProgrammaticAccessToken: PluralObjectTypeProgrammaticAccessTokens,
	}
}

func pluralObjectTypeToSingularMap() map[PluralObjectType]ObjectType {
	m := make(map[PluralObjectType]ObjectType)
	for k, v := range objectTypeSingularToPluralMap() {
		m[v] = k
	}
	return m
}

func (o ObjectType) Plural() PluralObjectType {
	if plural, ok := objectTypeSingularToPluralMap()[o]; ok {
		return plural
	}
	return PluralObjectType(o + "S")
}

// GetObjectIdentifier returns the ObjectIdentifier for the ObjectType and fully qualified name.
func (o ObjectType) GetObjectIdentifier(fullyQualifiedName string) ObjectIdentifier {
	accountObjectIdentifiers := []ObjectType{
		ObjectTypeParameter,
		ObjectTypeDatabase,
		ObjectTypeFailoverGroup,
		ObjectTypeIntegration,
		ObjectTypeResourceMonitor,
		ObjectTypeRole,
		ObjectTypeShare,
		ObjectTypeUser,
		ObjectTypeWarehouse,
	}
	if slices.Contains(accountObjectIdentifiers, o) {
		return NewAccountObjectIdentifier(fullyQualifiedName)
	}
	parts := strings.Split(fullyQualifiedName, ".")
	dbName := parts[0]
	if o == ObjectTypeSchema {
		schemaName := strings.Join(parts[1:], ".")
		return NewDatabaseObjectIdentifier(dbName, schemaName)
	}
	schemaName := parts[1]
	objectName := strings.Join(parts[2:], ".")
	return NewSchemaObjectIdentifier(dbName, schemaName, objectName)
}

type PluralObjectType string

const (
	PluralObjectTypeAccounts                 PluralObjectType = "ACCOUNTS"
	PluralObjectTypeManagedAccounts          PluralObjectType = "MANAGED ACCOUNTS"
	PluralObjectTypeUsers                    PluralObjectType = "USERS"
	PluralObjectTypeDatabaseRoles            PluralObjectType = "DATABASE ROLES"
	PluralObjectTypeDatasets                 PluralObjectType = "DATASETS"
	PluralObjectTypeRoles                    PluralObjectType = "ROLES"
	PluralObjectTypeIntegrations             PluralObjectType = "INTEGRATIONS"
	PluralObjectTypeNetworkPolicies          PluralObjectType = "NETWORK POLICIES"
	PluralObjectTypePasswordPolicies         PluralObjectType = "PASSWORD POLICIES"
	PluralObjectTypeSessionPolicies          PluralObjectType = "SESSION POLICIES"
	PluralObjectTypePrivacyPolicies          PluralObjectType = "PRIVACY POLICIES"
	PluralObjectTypeReplicationGroups        PluralObjectType = "REPLICATION GROUPS"
	PluralObjectTypeFailoverGroups           PluralObjectType = "FAILOVER GROUPS"
	PluralObjectTypeConnections              PluralObjectType = "CONNECTIONS"
	PluralObjectTypeParameters               PluralObjectType = "PARAMETERS"
	PluralObjectTypeWarehouses               PluralObjectType = "WAREHOUSES"
	PluralObjectTypeResourceMonitors         PluralObjectType = "RESOURCE MONITORS"
	PluralObjectTypeDatabases                PluralObjectType = "DATABASES"
	PluralObjectTypeSchemas                  PluralObjectType = "SCHEMAS"
	PluralObjectTypeShares                   PluralObjectType = "SHARES"
	PluralObjectTypeTables                   PluralObjectType = "TABLES"
	PluralObjectTypeDynamicTables            PluralObjectType = "DYNAMIC TABLES"
	PluralObjectTypeCortexSearchServices     PluralObjectType = "CORTEX SEARCH SERVICES"
	PluralObjectTypeExternalTables           PluralObjectType = "EXTERNAL TABLES"
	PluralObjectTypeEventTables              PluralObjectType = "EVENT TABLES"
	PluralObjectTypeViews                    PluralObjectType = "VIEWS"
	PluralObjectTypeMaterializedViews        PluralObjectType = "MATERIALIZED VIEWS"
	PluralObjectTypeSequences                PluralObjectType = "SEQUENCES"
	PluralObjectTypeSnapshots                PluralObjectType = "SNAPSHOTS"
	PluralObjectTypeFunctions                PluralObjectType = "FUNCTIONS"
	PluralObjectTypeExternalFunctions        PluralObjectType = "EXTERNAL FUNCTIONS"
	PluralObjectTypeProcedures               PluralObjectType = "PROCEDURES"
	PluralObjectTypeStreams                  PluralObjectType = "STREAMS"
	PluralObjectTypeTasks                    PluralObjectType = "TASKS"
	PluralObjectTypeMaskingPolicies          PluralObjectType = "MASKING POLICIES"
	PluralObjectTypeRowAccessPolicies        PluralObjectType = "ROW ACCESS POLICIES"
	PluralObjectTypeTags                     PluralObjectType = "TAGS"
	PluralObjectTypeSecrets                  PluralObjectType = "SECRETS"
	PluralObjectTypeStages                   PluralObjectType = "STAGES"
	PluralObjectTypeFileFormats              PluralObjectType = "FILE FORMATS"
	PluralObjectTypePipes                    PluralObjectType = "PIPES"
	PluralObjectTypeAlerts                   PluralObjectType = "ALERTS"
	PluralObjectTypeBudgets                  PluralObjectType = "SNOWFLAKE.CORE.BUDGET"
	PluralObjectTypeClassifications          PluralObjectType = "SNOWFLAKE.ML.CLASSIFICATION"
	PluralObjectTypeApplications             PluralObjectType = "APPLICATIONS"
	PluralObjectTypeApplicationPackages      PluralObjectType = "APPLICATION PACKAGES"
	PluralObjectTypeApplicationRoles         PluralObjectType = "APPLICATION ROLES"
	PluralObjectTypeStreamlits               PluralObjectType = "STREAMLITS"
	PluralObjectTypeIcebergTables            PluralObjectType = "ICEBERG TABLES"
	PluralObjectTypeExternalVolumes          PluralObjectType = "EXTERNAL VOLUMES"
	PluralObjectTypeNetworkRules             PluralObjectType = "NETWORK RULES"
	PluralObjectTypeNotebooks                PluralObjectType = "NOTEBOOKS"
	PluralObjectTypePackagesPolicies         PluralObjectType = "PACKAGES POLICIES"
	PluralObjectTypeComputePool              PluralObjectType = "COMPUTE POOLS"
	PluralObjectTypeAggregationPolicies      PluralObjectType = "AGGREGATION POLICIES"
	PluralObjectTypeAuthenticationPolicies   PluralObjectType = "AUTHENTICATION POLICIES"
	PluralObjectTypeHybridTables             PluralObjectType = "HYBRID TABLES"
	PluralObjectTypeImageRepositories        PluralObjectType = "IMAGE REPOSITORIES"
	PluralObjectTypeProjectionPolicies       PluralObjectType = "PROJECTION POLICIES"
	PluralObjectTypeDataMetricFunctions      PluralObjectType = "DATA METRIC FUNCTIONS"
	PluralObjectTypeGitRepositories          PluralObjectType = "GIT REPOSITORIES"
	PluralObjectTypeModels                   PluralObjectType = "MODELS"
	PluralObjectTypeServices                 PluralObjectType = "SERVICES"
	PluralObjectTypeProgrammaticAccessTokens PluralObjectType = "PROGRAMMATIC ACCESS TOKENS"
)

func (p PluralObjectType) String() string {
	return string(p)
}

func (p PluralObjectType) Singular() ObjectType {
	if singular, ok := pluralObjectTypeToSingularMap()[p]; ok {
		return singular
	}
	return ObjectType(strings.TrimSuffix(string(p), "S"))
}
