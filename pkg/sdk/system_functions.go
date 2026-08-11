package sdk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

type SystemFunctions interface {
	GetTag(ctx context.Context, tagID ObjectIdentifier, objectID ObjectIdentifier, objectType ObjectType) (*string, error)
	PipeStatus(pipeId SchemaObjectIdentifier) (PipeExecutionState, error)
	// PipeForceResume unpauses a pipe after ownership transfer. Snowflake will throw an error whenever a pipe changes its owner,
	// and someone tries to unpause it. To unpause a pipe after ownership transfer, this system function has to be called instead of ALTER PIPE.
	PipeForceResume(pipeId SchemaObjectIdentifier, options []ForceResumePipeOption) error
	EnableBehaviorChangeBundle(ctx context.Context, bundle string) error
	DisableBehaviorChangeBundle(ctx context.Context, bundle string) error
	ShowActiveBehaviorChangeBundles(ctx context.Context) ([]BehaviorChangeBundleInfo, error)
	BehaviorChangeBundleStatus(ctx context.Context, bundle string) (BehaviorChangeBundleStatus, error)
	GetIcebergTableInformation(ctx context.Context, id SchemaObjectIdentifier) (*IcebergTableInformation, error)
	// GetClusteringInformation returns clustering information for a table. If columns is empty, the table's defined
	// clustering key is used; otherwise the given expressions are used as the clustering key for the calculation.
	// See more in https://docs.snowflake.com/en/sql-reference/functions/system_clustering_information.
	GetClusteringInformation(ctx context.Context, id SchemaObjectIdentifier, columns ...string) (*ClusteringInformation, error)
	// GetCatalogLinkedDatabaseConfig returns the configuration of a catalog-linked database.
	// See more in https://docs.snowflake.com/en/sql-reference/functions/system_get_catalog_linked_database_config.
	GetCatalogLinkedDatabaseConfig(ctx context.Context, id AccountObjectIdentifier) (*CatalogLinkedDatabaseConfig, error)
	// GetCatalogLinkStatus returns the current state of the link between a catalog-linked database and its remote catalog.
	// See more in https://docs.snowflake.com/en/sql-reference/functions/system_catalog_link_status.
	GetCatalogLinkStatus(ctx context.Context, id AccountObjectIdentifier) (*CatalogLinkStatus, error)
}

var _ SystemFunctions = (*systemFunctions)(nil)

type systemFunctions struct {
	client *Client
}

func (c *systemFunctions) GetTag(ctx context.Context, tagID ObjectIdentifier, objectID ObjectIdentifier, objectType ObjectType) (*string, error) {
	objectType, err := normalizeGetTagObjectType(objectType)
	if err != nil {
		return nil, err
	}
	// ObjectTypeIcebergTableColumn is a dedicated type for handling iceberg table column tags.
	// However, Snowflake expects just the column type.
	if objectType == ObjectTypeIcebergTableColumn {
		objectType = ObjectTypeColumn
	}

	s := &struct {
		Tag sql.NullString `db:"TAG"`
	}{}
	sql := fmt.Sprintf(`SELECT SYSTEM$GET_TAG('%s', '%s', '%v') AS "TAG"`, tagID.FullyQualifiedName(), objectID.FullyQualifiedName(), objectType)
	err = c.client.queryOne(ctx, s, sql)
	if err != nil {
		return nil, err
	}
	if !s.Tag.Valid {
		return nil, nil
	}
	return &s.Tag.String, nil
}

// normalize object types for some values because of errors like below
// SQL compilation error: Invalid value VIEW for argument OBJECT_TYPE. Please use object type TABLE for all kinds of table-like objects.
// TODO [SNOW-1022645]: discuss how we handle situation like this in the SDK
func normalizeGetTagObjectType(objectType ObjectType) (ObjectType, error) {
	if !canBeAssociatedWithTag(objectType) {
		return "", fmt.Errorf("tagging for object type %s is not supported", objectType)
	}
	if slices.Contains([]ObjectType{ObjectTypeView, ObjectTypeMaterializedView, ObjectTypeExternalTable, ObjectTypeEventTable}, objectType) {
		return ObjectTypeTable, nil
	}

	if slices.Contains([]ObjectType{ObjectTypeExternalFunction}, objectType) {
		return ObjectTypeFunction, nil
	}

	return objectType, nil
}

type PipeExecutionState string

const (
	FailingOverPipeExecutionState                           PipeExecutionState = "FAILING_OVER"
	PausedPipeExecutionState                                PipeExecutionState = "PAUSED"
	ReadOnlyPipeExecutionState                              PipeExecutionState = "READ_ONLY"
	RunningPipeExecutionState                               PipeExecutionState = "RUNNING"
	StoppedBySnowflakeAdminPipeExecutionState               PipeExecutionState = "STOPPED_BY_SNOWFLAKE_ADMIN"
	StoppedClonedPipeExecutionState                         PipeExecutionState = "STOPPED_CLONED"
	StoppedFeatureDisabledPipeExecutionState                PipeExecutionState = "STOPPED_FEATURE_DISABLED"
	StoppedStageAlteredPipeExecutionState                   PipeExecutionState = "STOPPED_STAGE_ALTERED"
	StoppedStageDroppedPipeExecutionState                   PipeExecutionState = "STOPPED_STAGE_DROPPED"
	StoppedFileFormatDroppedPipeExecutionState              PipeExecutionState = "STOPPED_FILE_FORMAT_DROPPED"
	StoppedNotificationIntegrationDroppedPipeExecutionState PipeExecutionState = "STOPPED_NOTIFICATION_INTEGRATION_DROPPED"
	StoppedMissingPipePipeExecutionState                    PipeExecutionState = "STOPPED_MISSING_PIPE"
	StoppedMissingTablePipeExecutionState                   PipeExecutionState = "STOPPED_MISSING_TABLE"
	StalledCompilationErrorPipeExecutionState               PipeExecutionState = "STALLED_COMPILATION_ERROR"
	StalledInitializationErrorPipeExecutionState            PipeExecutionState = "STALLED_INITIALIZATION_ERROR"
	StalledExecutionErrorPipeExecutionState                 PipeExecutionState = "STALLED_EXECUTION_ERROR"
	StalledInternalErrorPipeExecutionState                  PipeExecutionState = "STALLED_INTERNAL_ERROR"
	StalledStagePermissionErrorPipeExecutionState           PipeExecutionState = "STALLED_STAGE_PERMISSION_ERROR"
)

func (c *systemFunctions) PipeStatus(pipeId SchemaObjectIdentifier) (PipeExecutionState, error) {
	row := &struct {
		PipeStatus string `db:"PIPE_STATUS"`
	}{}
	sql := fmt.Sprintf(`SELECT SYSTEM$PIPE_STATUS('%s') AS "PIPE_STATUS"`, pipeId.FullyQualifiedName())
	ctx := context.Background()

	err := c.client.queryOne(ctx, row, sql)
	if err != nil {
		return "", err
	}

	var pipeStatus map[string]any
	err = json.Unmarshal([]byte(row.PipeStatus), &pipeStatus)
	if err != nil {
		return "", err
	}

	if _, ok := pipeStatus["executionState"]; !ok {
		return "", NewError(fmt.Sprintf("executionState key not found in: %s", pipeStatus))
	}

	return PipeExecutionState(pipeStatus["executionState"].(string)), nil
}

type ForceResumePipeOption string

const (
	StalenessCheckOverrideForceResumePipeOption         ForceResumePipeOption = "STALENESS_CHECK_OVERRIDE"
	OwnershipTransferCheckOverrideForceResumePipeOption ForceResumePipeOption = "OWNERSHIP_TRANSFER_CHECK_OVERRIDE"
)

func (c *systemFunctions) PipeForceResume(pipeId SchemaObjectIdentifier, options []ForceResumePipeOption) error {
	ctx := context.Background()
	var functionOpts string
	if len(options) > 0 {
		stringOptions := collections.Map(options, func(opt ForceResumePipeOption) string { return string(opt) })
		functionOpts = fmt.Sprintf(", '%s'", strings.Join(stringOptions, ","))
	}
	_, err := c.client.exec(ctx, fmt.Sprintf("SELECT SYSTEM$PIPE_FORCE_RESUME('%s')%s", pipeId.FullyQualifiedName(), functionOpts))
	return err
}

func (c *systemFunctions) EnableBehaviorChangeBundle(ctx context.Context, bundle string) error {
	_, err := c.client.exec(ctx, fmt.Sprintf("SELECT SYSTEM$ENABLE_BEHAVIOR_CHANGE_BUNDLE('%s')", bundle))
	return err
}

func (c *systemFunctions) DisableBehaviorChangeBundle(ctx context.Context, bundle string) error {
	_, err := c.client.exec(ctx, fmt.Sprintf("SELECT SYSTEM$DISABLE_BEHAVIOR_CHANGE_BUNDLE('%s')", bundle))
	return err
}

type BehaviorChangeBundleInfo struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault"`
	IsEnabled bool   `json:"isEnabled"`
}

func (c *systemFunctions) ShowActiveBehaviorChangeBundles(ctx context.Context) ([]BehaviorChangeBundleInfo, error) {
	row := &struct {
		BundlesRaw string `db:"BUNDLES"`
	}{}
	sql := `SELECT SYSTEM$SHOW_ACTIVE_BEHAVIOR_CHANGE_BUNDLES() AS "BUNDLES"`
	err := c.client.queryOne(ctx, row, sql)
	if err != nil {
		return nil, err
	}
	var bundles []BehaviorChangeBundleInfo
	err = json.Unmarshal([]byte(row.BundlesRaw), &bundles)
	if err != nil {
		return nil, err
	}
	return bundles, nil
}

type BehaviorChangeBundleStatus string

const (
	BehaviorChangeBundleStatusEnabled  BehaviorChangeBundleStatus = "ENABLED"
	BehaviorChangeBundleStatusDisabled BehaviorChangeBundleStatus = "DISABLED"
	BehaviorChangeBundleStatusReleased BehaviorChangeBundleStatus = "RELEASED"
)

var allBehaviorChangeBundleStatuses = []BehaviorChangeBundleStatus{
	BehaviorChangeBundleStatusEnabled,
	BehaviorChangeBundleStatusDisabled,
	BehaviorChangeBundleStatusReleased,
}

func ToBehaviorChangeBundleStatus(s string) (BehaviorChangeBundleStatus, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allBehaviorChangeBundleStatuses, BehaviorChangeBundleStatus(s)) {
		return "", fmt.Errorf("invalid behavior change bundle status: %s", s)
	}
	return BehaviorChangeBundleStatus(s), nil
}

func (c *systemFunctions) BehaviorChangeBundleStatus(ctx context.Context, bundle string) (BehaviorChangeBundleStatus, error) {
	row := &struct {
		StatusRaw string `db:"STATUS"`
	}{}
	sql := fmt.Sprintf(`SELECT SYSTEM$BEHAVIOR_CHANGE_BUNDLE_STATUS('%s') AS "STATUS"`, bundle)
	err := c.client.queryOne(ctx, row, sql)
	if err != nil {
		return "", err
	}
	return ToBehaviorChangeBundleStatus(row.StatusRaw)
}

type icebergTableInformationDbStruct struct {
	MetadataLocation string `json:"metadataLocation"`
	Status           string `json:"status"`
}

type IcebergTableInformation struct {
	MetadataLocation string
}

// Implementation of SYSTEM$GET_ICEBERG_TABLE_INFORMATION - see https://docs.snowflake.com/en/sql-reference/functions/system_get_iceberg_table_information.
func (c *systemFunctions) GetIcebergTableInformation(ctx context.Context, id SchemaObjectIdentifier) (*IcebergTableInformation, error) {
	row := &struct {
		Info string `db:"ICEBERG_TABLE_INFORMATION"`
	}{}
	sql := fmt.Sprintf(`SELECT SYSTEM$GET_ICEBERG_TABLE_INFORMATION('%s') AS "ICEBERG_TABLE_INFORMATION"`, id.FullyQualifiedName())
	err := c.client.queryOne(ctx, row, sql)
	if err != nil {
		return nil, err
	}
	var info icebergTableInformationDbStruct
	if err = json.Unmarshal([]byte(row.Info), &info); err != nil {
		return nil, err
	}
	if !strings.EqualFold(info.Status, "success") {
		return nil, fmt.Errorf("getting iceberg table information failed: %s", info.Status)
	}

	return &IcebergTableInformation{
		MetadataLocation: info.MetadataLocation,
	}, nil
}

type clusteringInformationDbStruct struct {
	ClusterByKeys               string                         `json:"cluster_by_keys"`
	TotalPartitionCount         int                            `json:"total_partition_count"`
	TotalConstantPartitionCount int                            `json:"total_constant_partition_count"`
	AverageOverlaps             float64                        `json:"average_overlaps"`
	AverageDepth                float64                        `json:"average_depth"`
	PartitionDepthHistogram     map[string]int                 `json:"partition_depth_histogram"`
	ClusteringErrors            []clusteringErrorEntryDbStruct `json:"clustering_errors"`
}

type clusteringErrorEntryDbStruct struct {
	Timestamp string `json:"timestamp"`
	Error     string `json:"error"`
}

type ClusteringInformation struct {
	ClusterByKeys               string
	TotalPartitionCount         int
	TotalConstantPartitionCount int
	AverageOverlaps             float64
	AverageDepth                float64
	PartitionDepthHistogram     map[string]int
	ClusteringErrors            []ClusteringError
}

type ClusteringError struct {
	Timestamp string
	Error     string
}

// GetClusteringInformation is an implementation of SYSTEM$CLUSTERING_INFORMATION - see https://docs.snowflake.com/en/sql-reference/functions/system_clustering_information.
// The clustering key is not returned by SHOW/DESCRIBE for Iceberg tables, so this is the way to verify it.
func (c *systemFunctions) GetClusteringInformation(ctx context.Context, id SchemaObjectIdentifier, columns ...string) (*ClusteringInformation, error) {
	row := &struct {
		Info string `db:"CLUSTERING_INFORMATION"`
	}{}
	opts := &getClusteringInformationOptions{
		arguments: &clusteringInformationArgs{
			Name: id,
			Columns: collections.Map(columns, func(col string) clusteringInformationColumn {
				return clusteringInformationColumn{Name: col}
			}),
		},
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	if err := c.client.queryOne(ctx, row, sql); err != nil {
		return nil, err
	}
	return parseClusteringInformation(row.Info)
}

type getClusteringInformationOptions struct {
	selectSystemClusteringInformation bool                       `ddl:"static" sql:"SELECT SYSTEM$CLUSTERING_INFORMATION"`
	arguments                         *clusteringInformationArgs `ddl:"list,parentheses,must_parentheses"`
	as                                bool                       `ddl:"static" sql:"AS \"CLUSTERING_INFORMATION\""`
}

type clusteringInformationArgs struct {
	Name    SchemaObjectIdentifier        `ddl:"identifier,single_quotes"`
	Columns []clusteringInformationColumn `ddl:"parameter,no_equals,single_quotes,parentheses"`
}

type clusteringInformationColumn struct {
	Name string `ddl:"keyword,double_quotes"`
}

func parseClusteringInformation(raw string) (*ClusteringInformation, error) {
	var info clusteringInformationDbStruct
	if err := json.Unmarshal([]byte(raw), &info); err != nil {
		return nil, err
	}

	return &ClusteringInformation{
		ClusterByKeys:               info.ClusterByKeys,
		TotalPartitionCount:         info.TotalPartitionCount,
		TotalConstantPartitionCount: info.TotalConstantPartitionCount,
		AverageOverlaps:             info.AverageOverlaps,
		AverageDepth:                info.AverageDepth,
		PartitionDepthHistogram:     info.PartitionDepthHistogram,
		ClusteringErrors: collections.Map(info.ClusteringErrors, func(entry clusteringErrorEntryDbStruct) ClusteringError {
			return ClusteringError{
				Timestamp: entry.Timestamp,
				Error:     entry.Error,
			}
		}),
	}, nil
}

type catalogLinkedDatabaseConfigDbStruct struct {
	CatalogIntegration        string   `json:"catalog_integration"`
	CatalogName               *string  `json:"catalog_name"`
	ExternalVolume            *string  `json:"external_volume"`
	SyncIntervalSeconds       *int     `json:"sync_interval_seconds"`
	NamespaceMode             *string  `json:"namespace_mode"`
	NamespaceFlattenDelimiter *string  `json:"namespace_flatten_delimiter"`
	AllowedWriteOperations    *string  `json:"allowed_write_operations"`
	CatalogCaseSensitivity    *string  `json:"catalog_case_sensitivity"`
	IsSuspended               *bool    `json:"is_suspended"`
	AllowedNamespaces         []string `json:"allowed_namespaces"`
	BlockedNamespaces         []string `json:"blocked_namespaces"`
}

type CatalogLinkedDatabaseConfig struct {
	CatalogIntegration        string
	CatalogName               *string
	ExternalVolume            *string
	SyncIntervalSeconds       *int
	NamespaceMode             *string
	NamespaceFlattenDelimiter *string
	AllowedWriteOperations    *string
	CatalogCaseSensitivity    *string
	IsSuspended               *bool
	AllowedNamespaces         []string
	BlockedNamespaces         []string
}

// GetCatalogLinkedDatabaseConfig is an implementation of SYSTEM$GET_CATALOG_LINKED_DATABASE_CONFIG - see https://docs.snowflake.com/en/sql-reference/functions/system_get_catalog_linked_database_config.
// The LINKED_CATALOG parameters are not returned by SHOW/DESCRIBE for databases, so this is the way to read them back.
func (c *systemFunctions) GetCatalogLinkedDatabaseConfig(ctx context.Context, id AccountObjectIdentifier) (*CatalogLinkedDatabaseConfig, error) {
	row := &struct {
		Config string `db:"CONFIG"`
	}{}
	opts := &getCatalogLinkedDatabaseConfigOptions{
		arguments: &catalogLinkedDatabaseArgs{
			Name: id,
		},
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	if err := c.client.queryOne(ctx, row, sql); err != nil {
		return nil, err
	}
	return parseCatalogLinkedDatabaseConfig(row.Config)
}

type getCatalogLinkedDatabaseConfigOptions struct {
	selectSystemGetCatalogLinkedDatabaseConfig bool                       `ddl:"static" sql:"SELECT SYSTEM$GET_CATALOG_LINKED_DATABASE_CONFIG"`
	arguments                                  *catalogLinkedDatabaseArgs `ddl:"list,parentheses,must_parentheses"`
	as                                         bool                       `ddl:"static" sql:"AS \"CONFIG\""`
}

type catalogLinkedDatabaseArgs struct {
	Name AccountObjectIdentifier `ddl:"identifier,single_quotes"`
}

func parseCatalogLinkedDatabaseConfig(raw string) (*CatalogLinkedDatabaseConfig, error) {
	var config catalogLinkedDatabaseConfigDbStruct
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return nil, err
	}

	return &CatalogLinkedDatabaseConfig{
		CatalogIntegration:        config.CatalogIntegration,
		CatalogName:               config.CatalogName,
		ExternalVolume:            config.ExternalVolume,
		SyncIntervalSeconds:       config.SyncIntervalSeconds,
		NamespaceMode:             config.NamespaceMode,
		NamespaceFlattenDelimiter: config.NamespaceFlattenDelimiter,
		AllowedWriteOperations:    config.AllowedWriteOperations,
		CatalogCaseSensitivity:    config.CatalogCaseSensitivity,
		IsSuspended:               config.IsSuspended,
		AllowedNamespaces:         config.AllowedNamespaces,
		BlockedNamespaces:         config.BlockedNamespaces,
	}, nil
}

type catalogLinkFailureDetailDbStruct struct {
	QualifiedEntityName string `json:"qualifiedEntityName"`
	EntityDomain        string `json:"entityDomain"`
	Operation           string `json:"operation"`
	ErrorCode           string `json:"errorCode"`
	ErrorMessage        string `json:"errorMessage"`
}

type catalogLinkStatusDbStruct struct {
	ExecutionState                string                             `json:"executionState"`
	FailedExecutionStateReason    *string                            `json:"failedExecutionStateReason"`
	FailedExecutionStateErrorCode *string                            `json:"failedExecutionStateErrorCode"`
	LastLinkAttemptStartTime      *string                            `json:"lastLinkAttemptStartTime"`
	FailureDetails                []catalogLinkFailureDetailDbStruct `json:"failureDetails"`
}

type CatalogLinkStatus struct {
	ExecutionState                string
	FailedExecutionStateReason    *string
	FailedExecutionStateErrorCode *string
	LastLinkAttemptStartTime      *string
	FailureDetails                []CatalogLinkFailureDetail
}

type CatalogLinkFailureDetail struct {
	QualifiedEntityName string
	EntityDomain        string
	Operation           string
	ErrorCode           string
	ErrorMessage        string
}

// GetCatalogLinkStatus is an implementation of SYSTEM$CATALOG_LINK_STATUS - see https://docs.snowflake.com/en/sql-reference/functions/system_catalog_link_status.
func (c *systemFunctions) GetCatalogLinkStatus(ctx context.Context, id AccountObjectIdentifier) (*CatalogLinkStatus, error) {
	row := &struct {
		Status string `db:"STATUS"`
	}{}
	opts := &getCatalogLinkStatusOptions{
		arguments: &catalogLinkedDatabaseArgs{
			Name: id,
		},
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}
	if err := c.client.queryOne(ctx, row, sql); err != nil {
		return nil, err
	}
	return parseCatalogLinkStatus(row.Status)
}

type getCatalogLinkStatusOptions struct {
	selectSystemCatalogLinkStatus bool                       `ddl:"static" sql:"SELECT SYSTEM$CATALOG_LINK_STATUS"`
	arguments                     *catalogLinkedDatabaseArgs `ddl:"list,parentheses,must_parentheses"`
	as                            bool                       `ddl:"static" sql:"AS \"STATUS\""`
}

func parseCatalogLinkStatus(raw string) (*CatalogLinkStatus, error) {
	var status catalogLinkStatusDbStruct
	if err := json.Unmarshal([]byte(raw), &status); err != nil {
		return nil, err
	}

	return &CatalogLinkStatus{
		ExecutionState:                status.ExecutionState,
		FailedExecutionStateReason:    status.FailedExecutionStateReason,
		FailedExecutionStateErrorCode: status.FailedExecutionStateErrorCode,
		LastLinkAttemptStartTime:      status.LastLinkAttemptStartTime,
		FailureDetails: collections.Map(status.FailureDetails, func(entry catalogLinkFailureDetailDbStruct) CatalogLinkFailureDetail {
			return CatalogLinkFailureDetail{
				QualifiedEntityName: entry.QualifiedEntityName,
				EntityDomain:        entry.EntityDomain,
				Operation:           entry.Operation,
				ErrorCode:           entry.ErrorCode,
				ErrorMessage:        entry.ErrorMessage,
			}
		}),
	}, nil
}
