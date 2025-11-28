package main

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var _ ConvertibleCsvRow[WarehouseRepresentation] = new(WarehouseCsvRow)

type WarehouseCsvRow struct {
	Actives                              string `csv:"actives"`
	AutoResume                           string `csv:"auto_resume"`
	AutoSuspend                          string `csv:"auto_suspend"`
	Available                            string `csv:"available"`
	Comment                              string `csv:"comment"`
	CreatedOn                            string `csv:"created_on"`
	EnableQueryAcceleration              string `csv:"enable_query_acceleration"`
	Failed                               string `csv:"failed"`
	Generation                           string `csv:"generation"`
	IsCurrent                            string `csv:"is_current"`
	IsDefault                            string `csv:"is_default"`
	MaxClusterCount                      string `csv:"max_cluster_count"`
	MinClusterCount                      string `csv:"min_cluster_count"`
	Name                                 string `csv:"name"`
	Other                                string `csv:"other"`
	Owner                                string `csv:"owner"`
	OwnerRoleType                        string `csv:"owner_role_type"`
	Pendings                             string `csv:"pendings"`
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
	Suspended                            string `csv:"suspended"`
	Type                                 string `csv:"type"`
	UUID                                 string `csv:"u_u_i_d"`
	UpdatedOn                            string `csv:"updated_on"`
	AutoResumeLevel                      string `csv:"auto_resume_level"`
	AutoResumeValue                      string `csv:"auto_resume_value"`
	AutoSuspendLevel                     string `csv:"auto_suspend_level"`
	AutoSuspendValue                     string `csv:"auto_suspend_value"`
	CommentLevel                         string `csv:"comment_level"`
	CommentValue                         string `csv:"comment_value"`
	EnableQueryAccelerationLevel         string `csv:"enable_query_acceleration_level"`
	EnableQueryAccelerationValue         string `csv:"enable_query_acceleration_value"`
	GenerationLevel                      string `csv:"generation_level"`
	GenerationValue                      string `csv:"generation_value"`
	InitiallySuspendedLevel              string `csv:"initially_suspended_level"`
	InitiallySuspendedValue              string `csv:"initially_suspended_value"`
	MaxClusterCountLevel                 string `csv:"max_cluster_count_level"`
	MaxClusterCountValue                 string `csv:"max_cluster_count_value"`
	MinClusterCountLevel                 string `csv:"min_cluster_count_level"`
	MinClusterCountValue                 string `csv:"min_cluster_count_value"`
	QueryAccelerationMaxScaleFactorLevel string `csv:"query_acceleration_max_scale_factor_level"`
	QueryAccelerationMaxScaleFactorValue string `csv:"query_acceleration_max_scale_factor_value"`
	ResourceConstraintLevel              string `csv:"resource_constraint_level"`
	ResourceConstraintValue              string `csv:"resource_constraint_value"`
	ResourceMonitorLevel                 string `csv:"resource_monitor_level"`
	ResourceMonitorValue                 string `csv:"resource_monitor_value"`
	ScalingPolicyLevel                   string `csv:"scaling_policy_level"`
	ScalingPolicyValue                   string `csv:"scaling_policy_value"`
	WarehouseSizeLevel                   string `csv:"warehouse_size_level"`
	WarehouseSizeValue                   string `csv:"warehouse_size_value"`
	WarehouseTypeLevel                   string `csv:"warehouse_type_level"`
	WarehouseTypeValue                   string `csv:"warehouse_type_value"`
}

type WarehouseRepresentation struct {
	sdk.Warehouse

	// parameters
	AutoResume                      *bool
	AutoSuspend                     *int
	Comment                         *string
	EnableQueryAcceleration         *bool
	Generation                      *string
	InitiallySuspended              *bool
	MaxClusterCount                 *int
	MinClusterCount                 *int
	QueryAccelerationMaxScaleFactor *int
	ResourceConstraint              *string
	ResourceMonitor                 *string
	ScalingPolicy                   *string
	WarehouseSize                   *string
	WarehouseType                   *string
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
		handler.handleBooleanParameter(sdk.ParameterType(row.AutoResumeLevel), row.AutoResumeValue, &warehouseRepresentation.AutoResume),
		handler.handleIntegerParameter(sdk.ParameterType(row.AutoSuspendLevel), row.AutoSuspendValue, &warehouseRepresentation.AutoSuspend),
		handler.handleStringParameter(sdk.ParameterType(row.CommentLevel), row.CommentValue, &warehouseRepresentation.Comment),
		handler.handleBooleanParameter(sdk.ParameterType(row.EnableQueryAccelerationLevel), row.EnableQueryAccelerationValue, &warehouseRepresentation.EnableQueryAcceleration),
		handler.handleStringParameter(sdk.ParameterType(row.GenerationLevel), row.GenerationValue, &warehouseRepresentation.Generation),
		handler.handleBooleanParameter(sdk.ParameterType(row.InitiallySuspendedLevel), row.InitiallySuspendedValue, &warehouseRepresentation.InitiallySuspended),
		handler.handleIntegerParameter(sdk.ParameterType(row.MaxClusterCountLevel), row.MaxClusterCountValue, &warehouseRepresentation.MaxClusterCount),
		handler.handleIntegerParameter(sdk.ParameterType(row.MinClusterCountLevel), row.MinClusterCountValue, &warehouseRepresentation.MinClusterCount),
		handler.handleIntegerParameter(sdk.ParameterType(row.QueryAccelerationMaxScaleFactorLevel), row.QueryAccelerationMaxScaleFactorValue, &warehouseRepresentation.QueryAccelerationMaxScaleFactor),
		handler.handleStringParameter(sdk.ParameterType(row.ResourceConstraintLevel), row.ResourceConstraintValue, &warehouseRepresentation.ResourceConstraint),
		handler.handleStringParameter(sdk.ParameterType(row.ResourceMonitorLevel), row.ResourceMonitorValue, &warehouseRepresentation.ResourceMonitor),
		handler.handleStringParameter(sdk.ParameterType(row.ScalingPolicyLevel), row.ScalingPolicyValue, &warehouseRepresentation.ScalingPolicy),
		handler.handleStringParameter(sdk.ParameterType(row.WarehouseSizeLevel), row.WarehouseSizeValue, &warehouseRepresentation.WarehouseSize),
		handler.handleStringParameter(sdk.ParameterType(row.WarehouseTypeLevel), row.WarehouseTypeValue, &warehouseRepresentation.WarehouseType),
	)
	if errs != nil {
		return nil, errs
	}

	return warehouseRepresentation, nil
}
