package testacc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	// for PoC using the imports from testfunctional package
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/customtypes"
)

var _ resource.ResourceWithConfigure = &WarehouseRestApiPocResource{}

func NewWarehouseRestApiPocResource() resource.Resource {
	return &WarehouseRestApiPocResource{}
}

type WarehouseRestApiPocResource struct {
	SnowflakeRestApiEmbeddable
}

type WarehouseRestApiPocPrivateJson struct {
	WarehouseType                   sdk.WarehouseType `json:"warehouse_type,omitempty"`
	WarehouseSize                   sdk.WarehouseSize `json:"warehouse_size,omitempty"`
	MaxClusterCount                 int               `json:"max_cluster_count,omitempty"`
	MinClusterCount                 int               `json:"min_cluster_count,omitempty"`
	ScalingPolicy                   sdk.ScalingPolicy `json:"scaling_policy,omitempty"`
	AutoSuspend                     int               `json:"auto_suspend,omitempty"`
	AutoResume                      bool              `json:"auto_resume,omitempty"`
	ResourceMonitor                 string            `json:"resource_monitor,omitempty"`
	EnableQueryAcceleration         bool              `json:"enable_query_acceleration,omitempty"`
	QueryAccelerationMaxScaleFactor int               `json:"query_acceleration_max_scale_factor,omitempty"`

	WarehouseRestApiPocParametersPrivateJson
}

type WarehouseRestApiPocParametersPrivateJson struct {
	MaxConcurrencyLevel             int `json:"max_concurrency_level,omitempty"`
	StatementQueuedTimeoutInSeconds int `json:"statement_queued_timeout_in_seconds,omitempty"`
	StatementTimeoutInSeconds       int `json:"statement_timeout_in_seconds,omitempty"`
}

func warehouseRestApiPocPrivateJsonFromWarehouseApiModel(warehouse *WarehouseApiModel) (*WarehouseRestApiPocPrivateJson, error) {
	privateJson := &WarehouseRestApiPocPrivateJson{}
	if warehouse.WarehouseType != nil {
		privateJson.WarehouseType = *warehouse.WarehouseType
	}
	if warehouse.WarehouseSize != nil {
		privateJson.WarehouseSize = *warehouse.WarehouseSize
	}
	if warehouse.MaxClusterCount != nil {
		privateJson.MaxClusterCount = *warehouse.MaxClusterCount
	}
	if warehouse.MinClusterCount != nil {
		privateJson.MinClusterCount = *warehouse.MinClusterCount
	}
	if warehouse.ScalingPolicy != nil {
		privateJson.ScalingPolicy = *warehouse.ScalingPolicy
	}
	if warehouse.AutoSuspend != nil {
		privateJson.AutoSuspend = *warehouse.AutoSuspend
	}
	if warehouse.AutoResume != nil {
		v, err := strconv.ParseBool(*warehouse.AutoResume)
		if err != nil {
			return nil, err
		}
		privateJson.AutoResume = v
	}
	if warehouse.ResourceMonitor != nil {
		privateJson.ResourceMonitor = *warehouse.ResourceMonitor
	}
	if warehouse.EnableQueryAcceleration != nil {
		v, err := strconv.ParseBool(*warehouse.EnableQueryAcceleration)
		if err != nil {
			return nil, err
		}
		privateJson.EnableQueryAcceleration = v
	}
	if warehouse.QueryAccelerationMaxScaleFactor != nil {
		privateJson.QueryAccelerationMaxScaleFactor = *warehouse.QueryAccelerationMaxScaleFactor
	}
	if warehouse.MaxConcurrencyLevel != nil {
		privateJson.MaxConcurrencyLevel = *warehouse.MaxConcurrencyLevel
	}
	if warehouse.StatementQueuedTimeoutInSeconds != nil {
		privateJson.StatementQueuedTimeoutInSeconds = *warehouse.StatementQueuedTimeoutInSeconds
	}
	if warehouse.StatementTimeoutInSeconds != nil {
		privateJson.StatementTimeoutInSeconds = *warehouse.StatementTimeoutInSeconds
	}
	return privateJson, nil
}

func marshallWarehouseRestApiPocPrivateJson(warehouseApiModel *WarehouseApiModel) ([]byte, error) {
	warehouseJson, err := warehouseRestApiPocPrivateJsonFromWarehouseApiModel(warehouseApiModel)
	if err != nil {
		return nil, fmt.Errorf("could not create private json: %w", err)
	}
	bytes, err := json.Marshal(warehouseJson)
	if err != nil {
		return nil, fmt.Errorf("could not marshal json: %w", err)
	}
	return bytes, nil
}

func (r *WarehouseRestApiPocResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_warehouse_rest_api_poc"
}

func (r *WarehouseRestApiPocResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version:    0,
		Attributes: warehousePocAttributes(),
	}
}

func (r *WarehouseRestApiPocResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	if request.State.Raw.IsNull() || request.Plan.Raw.IsNull() {
		return
	}

	var plan, state *warehousePocModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	if !plan.Name.Equal(state.Name) {
		plan.FullyQualifiedName = types.StringUnknown()
		plan.Id = types.StringUnknown()
	}

	response.Diagnostics.Append(response.Plan.Set(ctx, &plan)...)
}

func (r *WarehouseRestApiPocResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	id, err := sdk.ParseAccountObjectIdentifier(request.ID)
	if err != nil {
		response.Diagnostics.AddError("Could not read ID in warehouse PoC", err.Error())
		return
	}

	client := r.client
	warehouse, err := client.Warehouses.GetByID(ctx, id)
	if err != nil {
		response.Diagnostics.AddError("Could not read Warehouse PoC", err.Error())
		return
	}
	data := &warehousePocModelV0{
		Id:   types.StringValue(helpers.EncodeResourceIdentifier(id)),
		Name: types.StringValue(id.Name()),
	}
	if warehouse.WarehouseType != nil {
		data.WarehouseType = customtypes.NewEnumValue(*warehouse.WarehouseType)
	}
	if warehouse.WarehouseSize != nil {
		data.WarehouseSize = customtypes.NewEnumValue(*warehouse.WarehouseSize)
	}
	if warehouse.MaxClusterCount != nil {
		data.MaxClusterCount = types.Int64Value(int64(*warehouse.MaxClusterCount))
	}
	if warehouse.MinClusterCount != nil {
		data.MinClusterCount = types.Int64Value(int64(*warehouse.MinClusterCount))
	}
	if warehouse.ScalingPolicy != nil {
		data.ScalingPolicy = customtypes.NewEnumValue(*warehouse.ScalingPolicy)
	}
	if warehouse.AutoSuspend != nil {
		data.AutoSuspend = types.Int64Value(int64(*warehouse.AutoSuspend))
	}
	if warehouse.AutoResume != nil {
		v, err := strconv.ParseBool(*warehouse.AutoResume)
		if err != nil {
			response.Diagnostics.AddError("Could not read AutoResume for Warehouse PoC", err.Error())
		} else {
			data.AutoResume = types.BoolValue(v)
		}
	}
	if warehouse.ResourceMonitor != nil {
		data.ResourceMonitor = types.StringValue(*warehouse.ResourceMonitor)
	}
	if warehouse.Comment != nil {
		data.Comment = types.StringValue(*warehouse.Comment)
	}
	if warehouse.EnableQueryAcceleration != nil {
		v, err := strconv.ParseBool(*warehouse.EnableQueryAcceleration)
		if err != nil {
			response.Diagnostics.AddError("Could not read EnableQueryAcceleration for Warehouse PoC", err.Error())
		} else {
			data.EnableQueryAcceleration = types.BoolValue(v)
		}
	}
	if warehouse.QueryAccelerationMaxScaleFactor != nil {
		data.QueryAccelerationMaxScaleFactor = types.Int64Value(int64(*warehouse.QueryAccelerationMaxScaleFactor))
	}
	if warehouse.MaxConcurrencyLevel != nil {
		data.MaxConcurrencyLevel = types.Int64Value(int64(*warehouse.QueryAccelerationMaxScaleFactor))
	}
	if warehouse.StatementQueuedTimeoutInSeconds != nil {
		data.StatementQueuedTimeoutInSeconds = types.Int64Value(int64(*warehouse.QueryAccelerationMaxScaleFactor))
	}
	if warehouse.StatementTimeoutInSeconds != nil {
		data.StatementTimeoutInSeconds = types.Int64Value(int64(*warehouse.QueryAccelerationMaxScaleFactor))
	}
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *WarehouseRestApiPocResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *warehousePocModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	id := sdk.NewAccountObjectIdentifier(data.Name.ValueString())

	warehouseModel := r.planToApiModel(data)

	response.Diagnostics.Append(r.create(ctx, warehouseModel)...)
	if response.Diagnostics.HasError() {
		return
	}

	data.FullyQualifiedName = types.StringValue(id.FullyQualifiedName())

	// we can use the existing encoder
	data.Id = types.StringValue(helpers.EncodeResourceIdentifier(id))

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	b, d := r.readAfterCreateOrUpdate(ctx, id, &response.State)
	if d.HasError() {
		response.Diagnostics.Append(d...)
		return
	}
	response.Diagnostics.Append(response.Private.SetKey(ctx, privateStateSnowflakeObjectsStateKey, b)...)
}

func (r *WarehouseRestApiPocResource) create(ctx context.Context, req WarehouseApiModel) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.client.Warehouses.Create(ctx, req)
	if err != nil {
		diags.AddError("Could not create warehouse PoC", err.Error())
	}

	return diags
}

func (r *WarehouseRestApiPocResource) readAfterCreateOrUpdate(ctx context.Context, id sdk.AccountObjectIdentifier, state *tfsdk.State) ([]byte, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	client := r.client
	warehouse, err := client.Warehouses.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			state.RemoveResource(ctx)
			diags.AddWarning("Failed to query warehouse. Marking the resource as removed.", fmt.Sprintf("Warehouse id: %s, Err: %s", id.FullyQualifiedName(), err))
		} else {
			diags.AddError("Could not read Warehouse PoC", err.Error())
		}
		return nil, diags
	}

	bytes, err := marshallWarehouseRestApiPocPrivateJson(warehouse)
	if err != nil {
		diags.AddError("Could not marshal json", err.Error())
		return nil, diags
	}

	return bytes, diags
}

func (r *WarehouseRestApiPocResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *warehousePocModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	id, err := sdk.ParseAccountObjectIdentifier(data.Id.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Could not read ID in warehouse PoC", err.Error())
		return
	}
	response.Diagnostics.Append(r.read(ctx, data, id, request, response)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *WarehouseRestApiPocResource) read(ctx context.Context, data *warehousePocModelV0, id sdk.AccountObjectIdentifier, request resource.ReadRequest, response *resource.ReadResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}

	client := r.client
	warehouse, err := client.Warehouses.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			response.State.RemoveResource(ctx)
			diags.AddWarning("Failed to query warehouse. Marking the resource as removed.", fmt.Sprintf("Warehouse id: %s, Err: %s", id.FullyQualifiedName(), err))
		} else {
			diags.AddError("Could not read Warehouse PoC", err.Error())
		}
		return diags
	}

	data.FullyQualifiedName = types.StringValue(id.FullyQualifiedName())

	prevValueBytes, d := request.Private.GetKey(ctx, privateStateSnowflakeObjectsStateKey)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	if prevValueBytes != nil {
		var prevValue WarehousePocPrivateJson
		err := json.Unmarshal(prevValueBytes, &prevValue)
		if err != nil {
			diags.AddError("Could not unmarshal json", err.Error())
			return diags
		}

		if warehouse.WarehouseType != nil && *warehouse.WarehouseType != prevValue.WarehouseType {
			data.WarehouseType = customtypes.NewEnumValue(*warehouse.WarehouseType)
		}
		if warehouse.WarehouseSize != nil && *warehouse.WarehouseSize != prevValue.WarehouseSize {
			data.WarehouseSize = customtypes.NewEnumValue(*warehouse.WarehouseSize)
		}
		if warehouse.MaxClusterCount != nil && *warehouse.MaxClusterCount != prevValue.MaxClusterCount {
			data.MaxClusterCount = types.Int64Value(int64(*warehouse.MaxClusterCount))
		}
		if warehouse.MinClusterCount != nil && *warehouse.MinClusterCount != prevValue.MinClusterCount {
			data.MinClusterCount = types.Int64Value(int64(*warehouse.MinClusterCount))
		}
		if warehouse.ScalingPolicy != nil && *warehouse.ScalingPolicy != prevValue.ScalingPolicy {
			data.ScalingPolicy = customtypes.NewEnumValue(*warehouse.ScalingPolicy)
		}
		if warehouse.AutoSuspend != nil && *warehouse.AutoSuspend != prevValue.AutoSuspend {
			data.AutoSuspend = types.Int64Value(int64(*warehouse.AutoSuspend))
		}
		if warehouse.AutoResume != nil && *warehouse.AutoResume != fmt.Sprintf("%t", prevValue.AutoResume) {
			v, err := strconv.ParseBool(*warehouse.AutoResume)
			if err != nil {
				response.Diagnostics.AddError("Could not read AutoResume for Warehouse PoC", err.Error())
			} else {
				data.AutoResume = types.BoolValue(v)
			}
		}
		if warehouse.ResourceMonitor != nil && *warehouse.ResourceMonitor != prevValue.ResourceMonitor {
			data.ResourceMonitor = types.StringValue(*warehouse.ResourceMonitor)
		}
		if warehouse.EnableQueryAcceleration != nil && *warehouse.EnableQueryAcceleration != fmt.Sprintf("%t", prevValue.EnableQueryAcceleration) {
			v, err := strconv.ParseBool(*warehouse.EnableQueryAcceleration)
			if err != nil {
				response.Diagnostics.AddError("Could not read EnableQueryAcceleration for Warehouse PoC", err.Error())
			} else {
				data.EnableQueryAcceleration = types.BoolValue(v)
			}
		}
		if warehouse.QueryAccelerationMaxScaleFactor != nil && *warehouse.QueryAccelerationMaxScaleFactor != prevValue.QueryAccelerationMaxScaleFactor {
			data.QueryAccelerationMaxScaleFactor = types.Int64Value(int64(*warehouse.QueryAccelerationMaxScaleFactor))
		}

		// simplified parameter handling as we don't have level
		if warehouse.MaxConcurrencyLevel != nil && *warehouse.MaxConcurrencyLevel != prevValue.MaxConcurrencyLevel {
			data.MaxConcurrencyLevel = types.Int64Value(int64(*warehouse.MaxConcurrencyLevel))
		}
		if warehouse.StatementQueuedTimeoutInSeconds != nil && *warehouse.StatementQueuedTimeoutInSeconds != prevValue.StatementQueuedTimeoutInSeconds {
			data.StatementQueuedTimeoutInSeconds = types.Int64Value(int64(*warehouse.StatementQueuedTimeoutInSeconds))
		}
		if warehouse.StatementTimeoutInSeconds != nil && *warehouse.StatementTimeoutInSeconds != prevValue.StatementTimeoutInSeconds {
			data.StatementTimeoutInSeconds = types.Int64Value(int64(*warehouse.StatementTimeoutInSeconds))
		}
	}

	if diags.HasError() {
		return diags
	}

	bytes, err := marshallWarehouseRestApiPocPrivateJson(warehouse)
	if err != nil {
		diags.AddError("Could not marshal json", err.Error())
		return diags
	}
	response.Diagnostics.Append(response.Private.SetKey(ctx, privateStateSnowflakeObjectsStateKey, bytes)...)

	return diags
}

func (r *WarehouseRestApiPocResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state *warehousePocModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	id, err := sdk.ParseAccountObjectIdentifier(state.Id.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Could not read ID in warehouse PoC", err.Error())
		return
	}

	// Change name separately
	if !plan.Name.Equal(state.Name) {
		newId := sdk.NewAccountObjectIdentifier(plan.Name.ValueString())

		err := r.client.Warehouses.Rename(ctx, id, newId)
		if err != nil {
			response.Diagnostics.AddError("Could not rename warehouse PoC", err.Error())
			return
		}

		plan.Id = types.StringValue(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	warehouseModel := r.planToApiModel(plan)

	// workaround for WaitForCompletion
	if warehouseModel.WarehouseSize != nil {
		warehouseModel.WaitForCompletion = sdk.String("true")
	}

	// Put (CoA)
	if err := r.client.Warehouses.CreateOrAlter(ctx, warehouseModel); err != nil {
		response.Diagnostics.AddError("Could not run create or alter in REST API PoC warehouse", err.Error())
		return
	}

	plan.FullyQualifiedName = types.StringValue(id.FullyQualifiedName())

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	b, d := r.readAfterCreateOrUpdate(ctx, id, &response.State)
	if d.HasError() {
		response.Diagnostics.Append(d...)
		return
	}
	response.Diagnostics.Append(response.Private.SetKey(ctx, privateStateSnowflakeObjectsStateKey, b)...)
}

func (r *WarehouseRestApiPocResource) planToApiModel(plan *warehousePocModelV0) WarehouseApiModel {
	warehouseModel := WarehouseApiModel{
		Name: plan.Name.ValueString(),
	}
	if !plan.WarehouseType.IsNull() {
		warehouseModel.WarehouseType = sdk.Pointer(sdk.WarehouseType(plan.WarehouseType.ValueString()))
	}
	if !plan.WarehouseSize.IsNull() {
		warehouseModel.WarehouseSize = sdk.Pointer(sdk.WarehouseSize(plan.WarehouseSize.ValueString()))
	}
	if !plan.MaxClusterCount.IsNull() {
		warehouseModel.MaxClusterCount = sdk.Pointer(int(plan.MaxClusterCount.ValueInt64()))
	}
	if !plan.MinClusterCount.IsNull() {
		warehouseModel.MinClusterCount = sdk.Pointer(int(plan.MinClusterCount.ValueInt64()))
	}
	if !plan.ScalingPolicy.IsNull() {
		warehouseModel.ScalingPolicy = sdk.Pointer(sdk.ScalingPolicy(plan.ScalingPolicy.ValueString()))
	}
	if !plan.AutoSuspend.IsNull() {
		warehouseModel.AutoSuspend = sdk.Pointer(int(plan.AutoSuspend.ValueInt64()))
	}
	if !plan.AutoResume.IsNull() {
		warehouseModel.AutoResume = sdk.Pointer(fmt.Sprintf("%t", plan.AutoResume.ValueBool()))
	}
	if !plan.InitiallySuspended.IsNull() {
		warehouseModel.InitiallySuspended = sdk.Pointer(fmt.Sprintf("%t", plan.InitiallySuspended.ValueBool()))
	}
	if !plan.ResourceMonitor.IsNull() {
		warehouseModel.ResourceMonitor = plan.ResourceMonitor.ValueStringPointer()
	}
	if !plan.Comment.IsNull() {
		warehouseModel.Comment = plan.Comment.ValueStringPointer()
	}
	if !plan.EnableQueryAcceleration.IsNull() {
		warehouseModel.EnableQueryAcceleration = sdk.Pointer(fmt.Sprintf("%t", plan.EnableQueryAcceleration.ValueBool()))
	}
	if !plan.QueryAccelerationMaxScaleFactor.IsNull() {
		warehouseModel.QueryAccelerationMaxScaleFactor = sdk.Pointer(int(plan.QueryAccelerationMaxScaleFactor.ValueInt64()))
	}
	if !plan.MaxConcurrencyLevel.IsNull() {
		warehouseModel.MaxConcurrencyLevel = sdk.Pointer(int(plan.MaxConcurrencyLevel.ValueInt64()))
	}
	if !plan.StatementQueuedTimeoutInSeconds.IsNull() {
		warehouseModel.StatementQueuedTimeoutInSeconds = sdk.Pointer(int(plan.StatementQueuedTimeoutInSeconds.ValueInt64()))
	}
	if !plan.StatementTimeoutInSeconds.IsNull() {
		warehouseModel.StatementTimeoutInSeconds = sdk.Pointer(int(plan.StatementTimeoutInSeconds.ValueInt64()))
	}
	return warehouseModel
}

func (r *WarehouseRestApiPocResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data *warehousePocModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	id, err := sdk.ParseAccountObjectIdentifier(data.Id.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Could not read ID in warehouse PoC", err.Error())
		return
	}

	err = r.client.Warehouses.Drop(ctx, id, sdk.Bool(true))
	if err != nil {
		response.Diagnostics.AddError("Could not delete warehouse PoC", err.Error())
		return
	}
}
