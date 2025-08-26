package testacc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	// for PoC using the imports from testfunctional package
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/customplanmodifiers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/customtypes"
)

var _ resource.ResourceWithConfigure = &WarehouseResource{}

func NewWarehousePocResource() resource.Resource {
	return &WarehouseResource{}
}

type WarehouseResource struct {
	SnowflakeClientEmbeddable
}

type warehousePocModelV0 struct {
	Name                            types.String                             `tfsdk:"name"`
	WarehouseType                   customtypes.EnumValue[sdk.WarehouseType] `tfsdk:"warehouse_type"`
	WarehouseSize                   customtypes.EnumValue[sdk.WarehouseSize] `tfsdk:"warehouse_size"`
	MaxClusterCount                 types.Int64                              `tfsdk:"max_cluster_count"`
	MinClusterCount                 types.Int64                              `tfsdk:"min_cluster_count"`
	ScalingPolicy                   customtypes.EnumValue[sdk.ScalingPolicy] `tfsdk:"scaling_policy"`
	AutoSuspend                     types.Int64                              `tfsdk:"auto_suspend"`
	AutoResume                      types.Bool                               `tfsdk:"auto_resume"`
	InitiallySuspended              types.Bool                               `tfsdk:"initially_suspended"`
	ResourceMonitor                 types.String                             `tfsdk:"resource_monitor"` // TODO [SNOW-2296366]: identifier type?
	Comment                         types.String                             `tfsdk:"comment"`
	EnableQueryAcceleration         types.Bool                               `tfsdk:"enable_query_acceleration"`
	QueryAccelerationMaxScaleFactor types.Int64                              `tfsdk:"query_acceleration_max_scale_factor"`

	// embedding to clearly distinct parameters from other attributes
	warehouseParametersModelV0

	Id types.String `tfsdk:"id"`
	fullyQualifiedNameModelEmbeddable
}

// we can't use here the WarehouseParameter type values as struct tags are pure literals
// this is really easy to generate though
type warehouseParametersModelV0 struct {
	MaxConcurrencyLevel             types.Int64 `tfsdk:"max_concurrency_level"`
	StatementQueuedTimeoutInSeconds types.Int64 `tfsdk:"statement_queued_timeout_in_seconds"`
	StatementTimeoutInSeconds       types.Int64 `tfsdk:"statement_timeout_in_seconds"`
}

type WarehousePocPrivateJson struct {
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

	WarehousePocParametersPrivateJson
}

type WarehousePocParametersPrivateJson struct {
	MaxConcurrencyLevel                  int               `json:"max_concurrency_level,omitempty"`
	MaxConcurrencyLevelLevel             sdk.ParameterType `json:"max_concurrency_level_level,omitempty"`
	StatementQueuedTimeoutInSeconds      int               `json:"statement_queued_timeout_in_seconds,omitempty"`
	StatementQueuedTimeoutInSecondsLevel sdk.ParameterType `json:"statement_queued_timeout_in_seconds_level,omitempty"`
	StatementTimeoutInSeconds            int               `json:"statement_timeout_in_seconds,omitempty"`
	StatementTimeoutInSecondsLevel       sdk.ParameterType `json:"statement_timeout_in_seconds_level,omitempty"`
}

func warehousePocPrivateJsonFromWarehouse(warehouse *sdk.Warehouse) *WarehousePocPrivateJson {
	return &WarehousePocPrivateJson{
		WarehouseType:                   warehouse.Type,
		WarehouseSize:                   warehouse.Size,
		MaxClusterCount:                 warehouse.MaxClusterCount,
		MinClusterCount:                 warehouse.MinClusterCount,
		ScalingPolicy:                   warehouse.ScalingPolicy,
		AutoSuspend:                     warehouse.AutoSuspend,
		AutoResume:                      warehouse.AutoResume,
		ResourceMonitor:                 warehouse.ResourceMonitor.Name(),
		EnableQueryAcceleration:         warehouse.EnableQueryAcceleration,
		QueryAccelerationMaxScaleFactor: warehouse.QueryAccelerationMaxScaleFactor,
	}
}

func warehousePocParametersPrivateJsonFromParameters(warehouseParameters []*sdk.Parameter) (*WarehousePocParametersPrivateJson, error) {
	privateJson := &WarehousePocParametersPrivateJson{}
	if err := marshalWarehousePocParameters(warehouseParameters, privateJson); err != nil {
		return nil, err
	}
	return privateJson, nil
}

func marshalWarehousePocPrivateJson(warehouse *sdk.Warehouse, warehouseParameters []*sdk.Parameter) ([]byte, error) {
	warehouseJson := warehousePocPrivateJsonFromWarehouse(warehouse)
	if warehouseParametersJson, err := warehousePocParametersPrivateJsonFromParameters(warehouseParameters); err != nil {
		return nil, err
	} else {
		warehouseJson.WarehousePocParametersPrivateJson = *warehouseParametersJson
	}
	bytes, err := json.Marshal(warehouseJson)
	if err != nil {
		return nil, fmt.Errorf("could not marshal json: %w", err)
	}
	return bytes, nil
}

func marshalWarehousePocParameters(warehouseParameters []*sdk.Parameter, privateJson *WarehousePocParametersPrivateJson) error {
	for _, parameter := range warehouseParameters {
		switch parameter.Key {
		case string(sdk.WarehouseParameterMaxConcurrencyLevel):
			if err := marshalWarehousePocParameterInt(parameter, &privateJson.MaxConcurrencyLevel, &privateJson.MaxConcurrencyLevelLevel); err != nil {
				return err
			}
		case string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds):
			if err := marshalWarehousePocParameterInt(parameter, &privateJson.StatementQueuedTimeoutInSeconds, &privateJson.StatementQueuedTimeoutInSecondsLevel); err != nil {
				return err
			}
		case string(sdk.WarehouseParameterStatementTimeoutInSeconds):
			if err := marshalWarehousePocParameterInt(parameter, &privateJson.StatementTimeoutInSeconds, &privateJson.StatementTimeoutInSecondsLevel); err != nil {
				return err
			}
		}
	}
	return nil
}

func marshalWarehousePocParameterInt(parameter *sdk.Parameter, field *int, levelField *sdk.ParameterType) error {
	value, err := strconv.Atoi(parameter.Value)
	if err != nil {
		return err
	}
	*field = value
	*levelField = parameter.Level
	return nil
}

func (r *WarehouseResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_warehouse_poc"
}

// TODO [SNOW-2296366]: suppress identifier quoting
// TODO [SNOW-2296366]: support all identifier types
// TODO [SNOW-2298113]: show_output and parameters
func warehousePocAttributes() map[string]schema.Attribute {
	existingWarehouseSchema := resources.Warehouse().Schema
	attrs := map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: existingWarehouseSchema["name"].Description,
			Required:    true,
		},
		"warehouse_type": schema.StringAttribute{
			Description: existingWarehouseSchema["warehouse_type"].Description,
			Optional:    true,
			CustomType:  customtypes.EnumType[sdk.WarehouseType]{},
			PlanModifiers: []planmodifier.String{
				customplanmodifiers.EnumSuppressor[sdk.WarehouseType](),
			},
		},
		"warehouse_size": schema.StringAttribute{
			Description: existingWarehouseSchema["warehouse_size"].Description,
			Optional:    true,
			CustomType:  customtypes.EnumType[sdk.WarehouseSize]{},
			PlanModifiers: []planmodifier.String{
				customplanmodifiers.EnumSuppressor[sdk.WarehouseSize](),
				customplanmodifiers.RequiresReplaceIfRemovedFromConfig(),
			},
		},
		"max_cluster_count": schema.Int64Attribute{
			Description: existingWarehouseSchema["max_cluster_count"].Description,
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"min_cluster_count": schema.Int64Attribute{
			Description: existingWarehouseSchema["min_cluster_count"].Description,
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		"scaling_policy": schema.StringAttribute{
			Description: existingWarehouseSchema["scaling_policy"].Description,
			Optional:    true,
			CustomType:  customtypes.EnumType[sdk.ScalingPolicy]{},
			PlanModifiers: []planmodifier.String{
				customplanmodifiers.EnumSuppressor[sdk.ScalingPolicy](),
			},
		},
		"auto_suspend": schema.Int64Attribute{
			Description: existingWarehouseSchema["auto_suspend"].Description,
			Optional:    true,
		},
		// boolean vs tri-value string in the SDKv2 implementation
		"auto_resume": schema.BoolAttribute{
			Description: existingWarehouseSchema["auto_resume"].Description,
			Optional:    true,
		},
		"initially_suspended": schema.BoolAttribute{
			Description: existingWarehouseSchema["initially_suspended"].Description,
			Optional:    true,
			// TODO [mux-PR]: IgnoreAfterCreation
		},
		"resource_monitor": schema.StringAttribute{
			Description: existingWarehouseSchema["resource_monitor"].Description,
			Optional:    true,
			// TODO [SNOW-2296366]: identifier validation
		},
		"comment": schema.StringAttribute{
			Description: existingWarehouseSchema["comment"].Description,
			Optional:    true,
		},
		"enable_query_acceleration": schema.BoolAttribute{
			Description: existingWarehouseSchema["enable_query_acceleration"].Description,
			Optional:    true,
		},
		// no SDKv2 IntDefault(-1) workaround needed
		"query_acceleration_max_scale_factor": schema.Int64Attribute{
			Description: existingWarehouseSchema["query_acceleration_max_scale_factor"].Description,
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.Between(0, 100),
			},
		},
		// parameters are not computed because we can't handle them the same way as in SDKv2 implementation
		strings.ToLower(string(sdk.WarehouseParameterMaxConcurrencyLevel)): schema.Int64Attribute{
			Description: existingWarehouseSchema[strings.ToLower(string(sdk.WarehouseParameterMaxConcurrencyLevel))].Description,
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		strings.ToLower(string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds)): schema.Int64Attribute{
			Description: existingWarehouseSchema[strings.ToLower(string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds))].Description,
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		strings.ToLower(string(sdk.WarehouseParameterStatementTimeoutInSeconds)): schema.Int64Attribute{
			Description: existingWarehouseSchema[strings.ToLower(string(sdk.WarehouseParameterStatementTimeoutInSeconds))].Description,
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.Between(0, 604800),
			},
		},
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Warehouse identifier.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		resources.FullyQualifiedNameAttributeName: GetFullyQualifiedNameResourceSchema(),
	}
	return attrs
}

func (r *WarehouseResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version:    0,
		Attributes: warehousePocAttributes(),
	}
}

func (r *WarehouseResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	if request.State.Raw.IsNull() || request.Plan.Raw.IsNull() {
		return
	}

	var plan, state *warehousePocModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}

	// TODO [mux-PR]: we can extract modifiers like earlier we had ComputedIfAnyAttributeChanged)
	if !plan.Name.Equal(state.Name) {
		plan.FullyQualifiedName = types.StringUnknown()
		plan.Id = types.StringUnknown()
	}

	// TODO [mux-PR]: add a functional test documenting that IgnoreChangeToCurrentSnowflakeValueInShow cannot be achieved that way.
	// Commented out on purpose for now.
	// r.simulateOldIgnoreChangeToCurrentSnowflakeValueInShow(ctx, request, response, plan, state)
	// if response.Diagnostics.HasError() {
	//	return
	// }

	response.Diagnostics.Append(response.Plan.Set(ctx, &plan)...)
}

// It results in error connected with https://github.com/hashicorp/terraform/blob/main/docs/resource-instance-change-lifecycle.md#planresourcechange behavior.
// For each value we receive error like:
//
//	| Error: Provider produced invalid plan
//	|
//	| Provider "registry.terraform.io/hashicorp/snowflake" planned an invalid value
//	| for snowflake_warehouse_poc.test.auto_suspend: planned value
//	| cty.NullVal(cty.Number) does not match config value cty.NumberIntVal(600).
//	|
//	| This is a bug in the provider, which should be reported in the provider's own
//	| issue tracker.
func (r *WarehouseResource) simulateOldIgnoreChangeToCurrentSnowflakeValueInShow(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse, plan *warehousePocModelV0, state *warehousePocModelV0) {
	// this serves the logic behind previous IgnoreChangeToCurrentSnowflakeValueInShow
	prevValueBytes, d := request.Private.GetKey(ctx, privateStateSnowflakeObjectsStateKey)
	response.Diagnostics.Append(d...)
	if response.Diagnostics.HasError() {
		return
	}
	if prevValueBytes != nil {
		var prevValue WarehousePocPrivateJson
		err := json.Unmarshal(prevValueBytes, &prevValue)
		if err != nil {
			response.Diagnostics.AddError("Could not unmarshal json", err.Error())
			return
		}

		// simplified checks for now
		if plan.WarehouseType.ValueString() == string(prevValue.WarehouseType) {
			plan.WarehouseType = state.WarehouseType
		}
		if plan.WarehouseSize.ValueString() == string(prevValue.WarehouseSize) {
			plan.WarehouseSize = state.WarehouseSize
		}
		if plan.MaxClusterCount.ValueInt64() == int64(prevValue.MaxClusterCount) {
			plan.MaxClusterCount = state.MaxClusterCount
		}
		if plan.MinClusterCount.ValueInt64() == int64(prevValue.MinClusterCount) {
			plan.MinClusterCount = state.MinClusterCount
		}
		if plan.ScalingPolicy.ValueString() == string(prevValue.ScalingPolicy) {
			plan.ScalingPolicy = state.ScalingPolicy
		}
		if plan.AutoSuspend.ValueInt64() == int64(prevValue.AutoSuspend) {
			plan.AutoSuspend = state.AutoSuspend
		}
		if plan.AutoResume.ValueBool() == prevValue.AutoResume {
			plan.AutoResume = state.AutoResume
		}
		if plan.ResourceMonitor.ValueString() == prevValue.ResourceMonitor {
			plan.ResourceMonitor = state.ResourceMonitor
		}
		if plan.EnableQueryAcceleration.ValueBool() != prevValue.EnableQueryAcceleration {
			plan.EnableQueryAcceleration = state.EnableQueryAcceleration
		}
		if plan.QueryAccelerationMaxScaleFactor.ValueInt64() == int64(prevValue.QueryAccelerationMaxScaleFactor) {
			plan.QueryAccelerationMaxScaleFactor = state.QueryAccelerationMaxScaleFactor
		}
	}
}

// From the docs https://developer.hashicorp.com/terraform/plugin/framework/resources/import
// (...) which must either specify enough Terraform state for the Read method to refresh [resource] or return an error.
func (r *WarehouseResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	id, err := sdk.ParseAccountObjectIdentifier(request.ID)
	if err != nil {
		response.Diagnostics.AddError("Could not read ID in warehouse PoC", err.Error())
		return
	}

	client := r.client
	warehouse, err := client.Warehouses.ShowByID(ctx, id)
	if err != nil {
		response.Diagnostics.AddError("Could not read Warehouse PoC", err.Error())
		return
	}
	data := &warehousePocModelV0{
		Id:                              types.StringValue(helpers.EncodeResourceIdentifier(id)),
		Name:                            types.StringValue(id.Name()),
		WarehouseType:                   customtypes.NewEnumValue(warehouse.Type),
		WarehouseSize:                   customtypes.NewEnumValue(warehouse.Size),
		MaxClusterCount:                 types.Int64Value(int64(warehouse.MaxClusterCount)),
		MinClusterCount:                 types.Int64Value(int64(warehouse.MinClusterCount)),
		ScalingPolicy:                   customtypes.NewEnumValue(warehouse.ScalingPolicy),
		AutoSuspend:                     types.Int64Value(int64(warehouse.AutoSuspend)),
		AutoResume:                      types.BoolValue(warehouse.AutoResume),
		ResourceMonitor:                 types.StringValue(warehouse.ResourceMonitor.Name()),
		Comment:                         types.StringValue(warehouse.Comment),
		EnableQueryAcceleration:         types.BoolValue(warehouse.EnableQueryAcceleration),
		QueryAccelerationMaxScaleFactor: types.Int64Value(int64(warehouse.QueryAccelerationMaxScaleFactor)),
	}

	warehouseParameters, err := client.Warehouses.ShowParameters(ctx, id)
	if err != nil {
		response.Diagnostics.AddError("Could not read Warehouse PoC parameters", err.Error())
		return
	}
	response.Diagnostics.Append(handleWarehousePocParameterImport(warehouseParameters, data)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func handleWarehousePocParameterImport(warehouseParameters []*sdk.Parameter, data *warehousePocModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}
	for _, parameter := range warehouseParameters {
		switch parameter.Key {
		case string(sdk.WarehouseParameterMaxConcurrencyLevel):
			handleWarehousePocParameterImportInt(parameter, &data.MaxConcurrencyLevel, sdk.ParameterTypeWarehouse, &diags)
		case string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds):
			handleWarehousePocParameterImportInt(parameter, &data.StatementQueuedTimeoutInSeconds, sdk.ParameterTypeWarehouse, &diags)
		case string(sdk.WarehouseParameterStatementTimeoutInSeconds):
			handleWarehousePocParameterImportInt(parameter, &data.StatementTimeoutInSeconds, sdk.ParameterTypeWarehouse, &diags)
		}
	}
	return diags
}

func handleWarehousePocParameterImportInt(parameter *sdk.Parameter, field *types.Int64, objectLevel sdk.ParameterType, diags *diag.Diagnostics) {
	value, err := strconv.Atoi(parameter.Value)
	if err != nil {
		diags.AddError(fmt.Sprintf("Handling parameter %s failed", parameter.Key), err.Error())
		return
	}
	if parameter.Level == objectLevel {
		*field = types.Int64Value(int64(value))
	}
}

func (r *WarehouseResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *warehousePocModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)

	opts := &sdk.CreateWarehouseOptions{}
	errs := errors.Join(
		testfunctional.StringEnumAttributeCreate(data.WarehouseType, &opts.WarehouseType),
		testfunctional.StringEnumAttributeCreate(data.WarehouseSize, &opts.WarehouseSize),
		testfunctional.Int64AttributeCreate(data.MaxClusterCount, &opts.MaxClusterCount),
		testfunctional.Int64AttributeCreate(data.MinClusterCount, &opts.MinClusterCount),
		testfunctional.StringEnumAttributeCreate(data.ScalingPolicy, &opts.ScalingPolicy),
		testfunctional.Int64AttributeCreate(data.AutoSuspend, &opts.AutoSuspend),
		testfunctional.BooleanAttributeCreate(data.AutoResume, &opts.AutoResume),
		testfunctional.BooleanAttributeCreate(data.InitiallySuspended, &opts.InitiallySuspended),
		testfunctional.IdAttributeCreate(data.ResourceMonitor, &opts.ResourceMonitor),
		testfunctional.StringAttributeCreate(data.Comment, &opts.Comment),
		testfunctional.BooleanAttributeCreate(data.EnableQueryAcceleration, &opts.EnableQueryAcceleration),
		testfunctional.Int64AttributeCreate(data.QueryAccelerationMaxScaleFactor, &opts.QueryAccelerationMaxScaleFactor),

		testfunctional.Int64AttributeCreate(data.MaxConcurrencyLevel, &opts.MaxConcurrencyLevel),
		testfunctional.Int64AttributeCreate(data.StatementQueuedTimeoutInSeconds, &opts.StatementQueuedTimeoutInSeconds),
		testfunctional.Int64AttributeCreate(data.StatementTimeoutInSeconds, &opts.StatementTimeoutInSeconds),
	)
	if errs != nil {
		response.Diagnostics.AddError("Error creating warehouse PoC", errs.Error())
		return
	}

	response.Diagnostics.Append(r.create(ctx, id, opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	// TODO [SNOW-2298108]: Adjust fully_qualified_name logic
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

func (r *WarehouseResource) create(ctx context.Context, id sdk.AccountObjectIdentifier, opts *sdk.CreateWarehouseOptions) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.client.Warehouses.Create(ctx, id, opts)
	if err != nil {
		diags.AddError("Could not create warehouse PoC", err.Error())
	}

	return diags
}

func (r *WarehouseResource) readAfterCreateOrUpdate(ctx context.Context, id sdk.AccountObjectIdentifier, state *tfsdk.State) ([]byte, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	client := r.client
	warehouse, err := client.Warehouses.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			state.RemoveResource(ctx)
			diags.AddWarning("Failed to query warehouse. Marking the resource as removed.", fmt.Sprintf("Warehouse id: %s, Err: %s", id.FullyQualifiedName(), err))
		} else {
			diags.AddError("Could not read Warehouse PoC", err.Error())
		}
		return nil, diags
	}

	warehouseParameters, err := client.Warehouses.ShowParameters(ctx, id)
	if err != nil {
		diags.AddError("Could not read Warehouse PoC parameters", err.Error())
		return nil, diags
	}

	bytes, err := marshalWarehousePocPrivateJson(warehouse, warehouseParameters)
	if err != nil {
		diags.AddError("Could not marshal json", err.Error())
		return nil, diags
	}

	return bytes, diags
}

func (r *WarehouseResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
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

func (r *WarehouseResource) read(ctx context.Context, data *warehousePocModelV0, id sdk.AccountObjectIdentifier, request resource.ReadRequest, response *resource.ReadResponse) diag.Diagnostics {
	diags := diag.Diagnostics{}

	client := r.client
	warehouse, err := client.Warehouses.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			response.State.RemoveResource(ctx)
			diags.AddWarning("Failed to query warehouse. Marking the resource as removed.", fmt.Sprintf("Warehouse id: %s, Err: %s", id.FullyQualifiedName(), err))
		} else {
			diags.AddError("Could not read Warehouse PoC", err.Error())
		}
		return diags
	}

	warehouseParameters, err := client.Warehouses.ShowParameters(ctx, id)
	if err != nil {
		diags.AddError("Could not read Warehouse PoC parameters", err.Error())
		return diags
	}

	// TODO [SNOW-2298108]: Adjust fully_qualified_name logic
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

		// TODO [SNOW-2298113]: introduce function like handleExternalChangesToObjectInShow or something similar
		if warehouse.Type != prevValue.WarehouseType {
			data.WarehouseType = customtypes.NewEnumValue(warehouse.Type)
		}
		if warehouse.Size != prevValue.WarehouseSize {
			data.WarehouseSize = customtypes.NewEnumValue(warehouse.Size)
		}
		if warehouse.MaxClusterCount != prevValue.MaxClusterCount {
			data.MaxClusterCount = types.Int64Value(int64(warehouse.MaxClusterCount))
		}
		if warehouse.MinClusterCount != prevValue.MinClusterCount {
			data.MinClusterCount = types.Int64Value(int64(warehouse.MinClusterCount))
		}
		if warehouse.ScalingPolicy != prevValue.ScalingPolicy {
			data.ScalingPolicy = customtypes.NewEnumValue(warehouse.ScalingPolicy)
		}
		if warehouse.AutoSuspend != prevValue.AutoSuspend {
			data.AutoSuspend = types.Int64Value(int64(warehouse.AutoSuspend))
		}
		if warehouse.AutoResume != prevValue.AutoResume {
			data.AutoResume = types.BoolValue(warehouse.AutoResume)
		}
		if warehouse.ResourceMonitor.Name() != prevValue.ResourceMonitor {
			data.ResourceMonitor = types.StringValue(warehouse.ResourceMonitor.Name())
		}
		if warehouse.EnableQueryAcceleration != prevValue.EnableQueryAcceleration {
			data.EnableQueryAcceleration = types.BoolValue(warehouse.EnableQueryAcceleration)
		}
		if warehouse.QueryAccelerationMaxScaleFactor != prevValue.QueryAccelerationMaxScaleFactor {
			data.QueryAccelerationMaxScaleFactor = types.Int64Value(int64(warehouse.QueryAccelerationMaxScaleFactor))
		}

		if parametersJson, err := warehousePocParametersPrivateJsonFromParameters(warehouseParameters); err != nil {
			diags.AddError("Could not read Warehouse PoC parameters", err.Error())
			return diags
		} else {
			if parametersJson.MaxConcurrencyLevel != prevValue.MaxConcurrencyLevel {
				data.MaxConcurrencyLevel = types.Int64Value(int64(parametersJson.MaxConcurrencyLevel))
			}
			if parametersJson.MaxConcurrencyLevelLevel != sdk.ParameterTypeWarehouse {
				data.MaxConcurrencyLevel = types.Int64Null()
			} else {
				data.MaxConcurrencyLevel = types.Int64Value(int64(parametersJson.MaxConcurrencyLevel))
			}
			if parametersJson.StatementQueuedTimeoutInSeconds != prevValue.StatementQueuedTimeoutInSeconds {
				data.StatementQueuedTimeoutInSeconds = types.Int64Value(int64(parametersJson.StatementQueuedTimeoutInSeconds))
			}
			if parametersJson.StatementQueuedTimeoutInSecondsLevel != sdk.ParameterTypeWarehouse {
				data.StatementQueuedTimeoutInSeconds = types.Int64Null()
			} else {
				data.StatementQueuedTimeoutInSeconds = types.Int64Value(int64(parametersJson.StatementQueuedTimeoutInSeconds))
			}
			if parametersJson.StatementTimeoutInSeconds != prevValue.StatementTimeoutInSeconds {
				data.StatementTimeoutInSeconds = types.Int64Value(int64(parametersJson.StatementTimeoutInSeconds))
			}
			if parametersJson.StatementTimeoutInSecondsLevel != sdk.ParameterTypeWarehouse {
				data.StatementTimeoutInSeconds = types.Int64Null()
			} else {
				data.StatementTimeoutInSeconds = types.Int64Value(int64(parametersJson.StatementTimeoutInSeconds))
			}
		}
	}

	bytes, err := marshalWarehousePocPrivateJson(warehouse, warehouseParameters)
	if err != nil {
		diags.AddError("Could not marshal json", err.Error())
		return diags
	}
	response.Diagnostics.Append(response.Private.SetKey(ctx, privateStateSnowflakeObjectsStateKey, bytes)...)

	// TODO [SNOW-2298113]: show_output and parameters

	return diags
}

func (r *WarehouseResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
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

		err := r.client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
			NewName: &newId,
		})
		if err != nil {
			response.Diagnostics.AddError("Could not rename warehouse PoC", err.Error())
			return
		}

		plan.Id = types.StringValue(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	// Batch SET operations and UNSET operations
	set := sdk.WarehouseSet{}
	unset := sdk.WarehouseUnset{}

	errs := errors.Join(
		// name handled in rename
		// unset for warehouse type does not work, setting the default instead
		testfunctional.StringEnumAttributeUpdateSetDefaultInsteadOfUnset(plan.WarehouseType, state.WarehouseType, &set.WarehouseType, sdk.WarehouseTypeStandard),
		// removing from config is handled with resource recreation
		testfunctional.StringEnumAttributeUpdateSetOnly(plan.WarehouseSize, state.WarehouseSize, &set.WarehouseSize),
		testfunctional.Int64AttributeUpdate(plan.MaxClusterCount, state.MaxClusterCount, &set.MaxClusterCount, &unset.MaxClusterCount),
		testfunctional.Int64AttributeUpdate(plan.MinClusterCount, state.MinClusterCount, &set.MinClusterCount, &unset.MinClusterCount),
		// unset for scaling policy does not work, setting the default instead
		testfunctional.StringEnumAttributeUpdateSetDefaultInsteadOfUnset(plan.ScalingPolicy, state.ScalingPolicy, &set.ScalingPolicy, sdk.ScalingPolicyStandard),
		// unset for auto_suspend does not work, setting the default instead
		testfunctional.Int64AttributeUpdateSetDefaultInsteadOfUnset(plan.AutoSuspend, state.AutoSuspend, &set.AutoSuspend, 600),
		// unset for auto_resume works incorrectly, setting the default instead
		testfunctional.BooleanAttributeUpdateSetDefaultInsteadOfUnset(plan.AutoResume, state.AutoResume, &set.AutoResume, true),
		testfunctional.IdAttributeUpdate(plan.ResourceMonitor, state.ResourceMonitor, &set.ResourceMonitor, &unset.ResourceMonitor),
		testfunctional.StringAttributeUpdate(plan.Comment, state.Comment, &set.Comment, &unset.Comment),
		testfunctional.BooleanAttributeUpdate(plan.EnableQueryAcceleration, state.EnableQueryAcceleration, &set.EnableQueryAcceleration, &unset.EnableQueryAcceleration),
		testfunctional.Int64AttributeUpdate(plan.QueryAccelerationMaxScaleFactor, state.QueryAccelerationMaxScaleFactor, &set.QueryAccelerationMaxScaleFactor, &unset.QueryAccelerationMaxScaleFactor),

		// in the SDK implementation we have the parameters handling separated; for now, here it was not needed
		testfunctional.Int64AttributeUpdate(plan.MaxConcurrencyLevel, state.MaxConcurrencyLevel, &set.MaxConcurrencyLevel, &unset.MaxConcurrencyLevel),
		testfunctional.Int64AttributeUpdate(plan.StatementQueuedTimeoutInSeconds, state.StatementQueuedTimeoutInSeconds, &set.StatementQueuedTimeoutInSeconds, &unset.StatementQueuedTimeoutInSeconds),
		testfunctional.Int64AttributeUpdate(plan.StatementTimeoutInSeconds, state.StatementTimeoutInSeconds, &set.StatementTimeoutInSeconds, &unset.StatementTimeoutInSeconds),
	)
	if errs != nil {
		response.Diagnostics.AddError("Error updating warehouse PoC", errs.Error())
		return
	}
	// workaround for WaitForCompletion
	if set.WarehouseSize != nil {
		set.WaitForCompletion = sdk.Bool(true)
	}

	// Apply SET and UNSET changes
	if (set != sdk.WarehouseSet{}) {
		err := r.client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
			Set: &set,
		})
		if err != nil {
			response.Diagnostics.AddError("Could not update (alter set) warehouse PoC", err.Error())
			return
		}
	}
	if (unset != sdk.WarehouseUnset{}) {
		err := r.client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
			Unset: &unset,
		})
		if err != nil {
			response.Diagnostics.AddError("Could not update (alter unset) warehouse PoC", err.Error())
			return
		}
	}

	// TODO [SNOW-2298108]: Adjust fully_qualified_name logic
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

// For SDKv2 resources we have a method handling deletion common cases; we can add something similar later
func (r *WarehouseResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data *warehousePocModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	id, err := sdk.ParseAccountObjectIdentifier(data.Id.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Could not read ID in warehouse PoC", err.Error())
		return
	}

	err = r.client.Warehouses.DropSafely(ctx, id)
	if err != nil {
		response.Diagnostics.AddError("Could not delete warehouse PoC", err.Error())
		return
	}
}
