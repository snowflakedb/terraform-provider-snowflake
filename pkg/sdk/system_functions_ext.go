package sdk

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

type BehaviorChangeBundleInfo struct {
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefault"`
	IsEnabled bool   `json:"isEnabled"`
}

type clusteringInformationDbStruct struct {
	ClusterByKeys               string                         `json:"cluster_by_keys"`
	Version                     string                         `json:"version"`
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

type icebergTableInformationDbStruct struct {
	MetadataLocation string `json:"metadataLocation"`
	Status           string `json:"status"`
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
	LastLinkAttemptStartTime      *time.Time                         `json:"lastLinkAttemptStartTime"`
	FailureDetails                []catalogLinkFailureDetailDbStruct `json:"failureDetails"`
}

func (r *GetTagRequest) adjust() {
	r.Arguments.ObjectType = normalizeGetTagObjectType(r.Arguments.ObjectType)
}

// Snowflake SYSTEM$GET_TAG domains that differ from the DDL object type.
// DDL uses AGENT (CREATE/ALTER AGENT); GET_TAG expects CORTEX AGENT.
const getTagDomainCortexAgent ObjectType = "CORTEX AGENT"

// normalize object types for some values because of errors like below
// SQL compilation error: Invalid value VIEW for argument OBJECT_TYPE. Please use object type TABLE for all kinds of table-like objects.
// TODO [SNOW-1022645]: discuss how we handle situation like this in the SDK
func normalizeGetTagObjectType(objectType ObjectType) ObjectType {
	if slices.Contains([]ObjectType{ObjectTypeView, ObjectTypeMaterializedView, ObjectTypeExternalTable, ObjectTypeEventTable}, objectType) {
		return ObjectTypeTable
	}
	if objectType == ObjectTypeExternalFunction {
		return ObjectTypeFunction
	}
	// ObjectTypeIcebergTableColumn is a dedicated type for handling iceberg table column tags.
	// However, Snowflake expects just the column type.
	if objectType == ObjectTypeIcebergTableColumn {
		return ObjectTypeColumn
	}
	if objectType == ObjectTypeAgent {
		return getTagDomainCortexAgent
	}
	return objectType
}

func (opts *GetTagOptions) additionalValidations() error {
	if err := validateUserInput(opts.Arguments.ObjectType.String()); err != nil {
		return fmt.Errorf("invalid object type: %w", err)
	}
	return nil
}

func parsePipeExecutionState(raw string) (PipeExecutionState, error) {
	var pipeStatus map[string]any
	if err := json.Unmarshal([]byte(raw), &pipeStatus); err != nil {
		return "", err
	}
	executionState, ok := pipeStatus["executionState"]
	if !ok {
		return "", NewError(fmt.Sprintf("executionState key not found in: %s", pipeStatus))
	}
	return PipeExecutionState(executionState.(string)), nil
}

func parseIcebergTableMetadataLocation(raw string) (string, error) {
	var info icebergTableInformationDbStruct
	if err := json.Unmarshal([]byte(raw), &info); err != nil {
		return "", err
	}
	if !strings.EqualFold(info.Status, "success") {
		return "", fmt.Errorf("getting iceberg table information failed: %s", info.Status)
	}
	return info.MetadataLocation, nil
}

func (r clusteringInformationRow) additionalConvert(result *ClusteringInformation) error {
	parsed, err := parseClusteringInformation(r.ClusteringInformation)
	if err != nil {
		return err
	}
	*result = *parsed
	return nil
}

func (r catalogLinkedDatabaseConfigRow) additionalConvert(result *CatalogLinkedDatabaseConfig) error {
	parsed, err := parseCatalogLinkedDatabaseConfig(r.Config)
	if err != nil {
		return err
	}
	*result = *parsed
	return nil
}

func (r catalogLinkStatusRow) additionalConvert(result *CatalogLinkStatus) error {
	parsed, err := parseCatalogLinkStatus(r.Status)
	if err != nil {
		return err
	}
	*result = *parsed
	return nil
}

func parseClusteringInformation(raw string) (*ClusteringInformation, error) {
	var info clusteringInformationDbStruct
	if err := json.Unmarshal([]byte(raw), &info); err != nil {
		return nil, err
	}

	return &ClusteringInformation{
		ClusterByKeys:               info.ClusterByKeys,
		Version:                     info.Version,
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

func parseCatalogLinkedDatabaseConfig(raw string) (*CatalogLinkedDatabaseConfig, error) {
	var config catalogLinkedDatabaseConfigDbStruct
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return nil, err
	}

	result := &CatalogLinkedDatabaseConfig{
		CatalogName:               config.CatalogName,
		SyncIntervalSeconds:       config.SyncIntervalSeconds,
		NamespaceFlattenDelimiter: config.NamespaceFlattenDelimiter,
		IsSuspended:               config.IsSuspended,
		AllowedNamespaces:         config.AllowedNamespaces,
		BlockedNamespaces:         config.BlockedNamespaces,
	}
	mapStringWithMapping(&result.CatalogIntegration, config.CatalogIntegration, ParseAccountObjectIdentifier)
	if config.ExternalVolume != nil {
		externalVolume, err := ParseAccountObjectIdentifier(*config.ExternalVolume)
		if err != nil {
			return nil, err
		}
		result.ExternalVolume = new(externalVolume)
	}
	if config.NamespaceMode != nil {
		namespaceMode, err := ToCatalogLinkedDatabaseNamespaceMode(*config.NamespaceMode)
		if err != nil {
			return nil, err
		}
		result.NamespaceMode = new(namespaceMode)
	}
	if config.AllowedWriteOperations != nil {
		allowedWriteOperations, err := ToCatalogLinkedDatabaseAllowedWriteOperations(*config.AllowedWriteOperations)
		if err != nil {
			return nil, err
		}
		result.AllowedWriteOperations = new(allowedWriteOperations)
	}
	if config.CatalogCaseSensitivity != nil {
		catalogCaseSensitivity, err := ToDatabaseCatalogCaseSensitivity(*config.CatalogCaseSensitivity)
		if err != nil {
			return nil, err
		}
		result.CatalogCaseSensitivity = new(catalogCaseSensitivity)
	}
	return result, nil
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
