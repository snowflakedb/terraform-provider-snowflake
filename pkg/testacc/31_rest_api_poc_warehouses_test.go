package testacc

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

type WarehousesPoc interface {
	Create(ctx context.Context, req WarehouseApiModel) error
	CreateOrAlter(ctx context.Context, req WarehouseApiModel) error
	GetByID(ctx context.Context, id sdk.AccountObjectIdentifier) (*WarehouseApiModel, error)
}

// WarehouseApiModel has almost the same fields as sdk.CreateWarehouseOptions and sdk.WarehouseSet.
// For objects where we already have the request builders, like sdk.CreateDatabaseRoleRequest, we could do conversion from the request temporarily.
// All of POST, PUT, and GET have the same attributes (so we are reusing a single struct for now):
//   - https://docs.snowflake.com/developer-guide/snowflake-rest-api/reference/warehouse#post--api-v2-warehouses
//   - https://docs.snowflake.com/developer-guide/snowflake-rest-api/reference/warehouse#put--api-v2-warehouses-name
//   - https://docs.snowflake.com/developer-guide/snowflake-rest-api/reference/warehouse#get--api-v2-warehouses-name
type WarehouseApiModel struct {
	// required
	Name sdk.AccountObjectIdentifier `json:"name"`

	// optional attributes
	WarehouseType                   *sdk.WarehouseType           `json:"warehouse_type,omitempty"`
	WarehouseSize                   *sdk.WarehouseSize           `json:"warehouse_size,omitempty"`
	WaitForCompletion               *bool                        `json:"wait_for_completion,omitempty"`
	MaxClusterCount                 *int                         `json:"max_cluster_count,omitempty"`
	MinClusterCount                 *int                         `json:"min_cluster_count,omitempty"`
	ScalingPolicy                   *sdk.ScalingPolicy           `json:"scaling_policy,omitempty"`
	AutoSuspend                     *int                         `json:"auto_suspend,omitempty"`
	AutoResume                      *bool                        `json:"auto_resume,omitempty"`
	InitiallySuspended              *bool                        `json:"initially_suspended,omitempty"`
	ResourceMonitor                 *sdk.AccountObjectIdentifier `json:"resource_monitor,omitempty"`
	Comment                         *string                      `json:"comment,omitempty"`
	EnableQueryAcceleration         *bool                        `json:"enable_query_acceleration,omitempty"`
	QueryAccelerationMaxScaleFactor *int                         `json:"query_acceleration_max_scale_factor,omitempty"`

	// optional parameters
	MaxConcurrencyLevel             *int `json:"max_concurrency_level,omitempty"`
	StatementQueuedTimeoutInSeconds *int `json:"statement_queued_timeout_in_seconds,omitempty"`
	StatementTimeoutInSeconds       *int `json:"statement_timeout_in_seconds,omitempty"`
}

var _ WarehousesPoc = (*warehousesPoc)(nil)

type warehousesPoc struct {
	client *RestApiPocClient
}

func (w warehousesPoc) Create(ctx context.Context, req WarehouseApiModel) error {
	_, err := post(ctx, w.client.httpClient, w.client.url, "warehouses", req)
	if err != nil {
		return fmt.Errorf("warehousesPoc.Create: %w", err)
	}
	return nil
}

func (w warehousesPoc) CreateOrAlter(ctx context.Context, req WarehouseApiModel) error {
	panic("implement me")
}

func (w warehousesPoc) GetByID(ctx context.Context, id sdk.AccountObjectIdentifier) (*WarehouseApiModel, error) {
	panic("implement me")
}
