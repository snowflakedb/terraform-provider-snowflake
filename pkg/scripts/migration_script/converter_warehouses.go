package main

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var _ ConvertibleCsvRow[WarehouseRepresentation] = new(WarehouseCsvRow)

type WarehouseCsvRow struct {
	AutoResume                           string `csv:"auto_resume"`
	AutoSuspend                          string `csv:"auto_suspend"`
	Available                            string `csv:"available"`
	Comment                              string `csv:"comment"`
	CreatedOn                            string `csv:"created_on"`
	EnableQueryAcceleration              string `csv:"enable_query_acceleration"`
	Generation                           string `csv:"generation"`
	IsCurrent                            string `csv:"is_current"`
	IsDefault                            string `csv:"is_default"`
	MaxClusterCount                      string `csv:"max_cluster_count"`
	MinClusterCount                      string `csv:"min_cluster_count"`
	Name                                 string `csv:"name"`
	Other                                string `csv:"other"`
	Owner                                string `csv:"owner"`
	OwnerRoleType                        string `csv:"owner_role_type"`
	Provisioning                         string `csv:"provisioning"`
	QueryAccelerationMaxScaleFactor      string `csv:"query_acceleration_max_scale_factor"`
	Queued                               string `csv:"queued"`
	Quiescing                            string `csv:"quiescing"`
	ResourceConstraint                   string `csv:"resource_constraint"`
	ResourceMonitor                      string `csv:"resource_monitor"`
	ResumedOn                            string `csv:"resumed_on"`
	Running                              string `csv:"running"`
	ScalingPolicy                        string `csv:"scaling_policy"`
	Size                                 string `csv:"size"`
	StartedClusters                      string `csv:"started_clusters"`
	State                                string `csv:"state"`
	Type                                 string `csv:"type"`
	UpdatedOn                            string `csv:"updated_on"`
	MaxConcurrencyLevelLevel             string `csv:"max_concurrency_level_level"`
	MaxConcurrencyLevelValue             string `csv:"max_concurrency_level_value"`
	StatementQueuedTimeoutInSecondsLevel string `csv:"statement_queued_timeout_in_seconds_level"`
	StatementQueuedTimeoutInSecondsValue string `csv:"statement_queued_timeout_in_seconds_value"`
	StatementTimeoutInSecondsLevel       string `csv:"statement_timeout_in_seconds_level"`
	StatementTimeoutInSecondsValue       string `csv:"statement_timeout_in_seconds_value"`
}

type WarehouseRepresentation struct {
	sdk.Warehouse

	// parameters
	MaxConcurrencyLevel             *int
	StatementQueuedTimeoutInSeconds *int
	StatementTimeoutInSeconds       *int
}

func (row WarehouseCsvRow) convert() (*WarehouseRepresentation, error) {
	warehouseRepresentation := &WarehouseRepresentation{
		Warehouse: sdk.Warehouse{
			Name:                    row.Name,
			IsCurrent:               row.IsCurrent == "Y",
			IsDefault:               row.IsDefault == "Y",
			Owner:                   row.Owner,
			Comment:                 row.Comment,
			Type:                    sdk.WarehouseType(row.Type),
			Size:                    sdk.WarehouseSize(row.Size),
			OwnerRoleType:           row.OwnerRoleType,
			State:                   sdk.WarehouseState(row.State),
			ScalingPolicy:           sdk.ScalingPolicy(row.ScalingPolicy),
			EnableQueryAcceleration: row.EnableQueryAcceleration == "true",
		},
	}

	handler := newParameterHandler(sdk.ParameterTypeWarehouse)
	errs := errors.Join(
		handler.handleIntegerParameter(sdk.ParameterType(row.MaxConcurrencyLevelLevel), row.MaxConcurrencyLevelValue, &warehouseRepresentation.MaxConcurrencyLevel),
		handler.handleIntegerParameter(sdk.ParameterType(row.StatementQueuedTimeoutInSecondsLevel), row.StatementQueuedTimeoutInSecondsValue, &warehouseRepresentation.StatementQueuedTimeoutInSeconds),
		handler.handleIntegerParameter(sdk.ParameterType(row.StatementTimeoutInSecondsLevel), row.StatementTimeoutInSecondsValue, &warehouseRepresentation.StatementTimeoutInSeconds),
	)
	if errs != nil {
		return nil, errs
	}

	return warehouseRepresentation, nil
}
