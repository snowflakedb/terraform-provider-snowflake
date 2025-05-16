package sdk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
)

type SystemFunctions interface {
	GetTag(ctx context.Context, tagID ObjectIdentifier, objectID ObjectIdentifier, objectType ObjectType) (*string, error)
	PipeStatus(pipeId SchemaObjectIdentifier) (PipeExecutionState, error)
	// PipeForceResume unpauses a pipe after ownership transfer. Snowflake will throw an error whenever a pipe changes its owner,
	// and someone tries to unpause it. To unpause a pipe after ownership transfer, this system function has to be called instead of ALTER PIPE.
	PipeForceResume(pipeId SchemaObjectIdentifier, options []ForceResumePipeOption) error
	EnableBehaviorChangeBundle(ctx context.Context, bundle string) error
	DisableBehaviorChangeBundle(ctx context.Context, bundle string) error
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
