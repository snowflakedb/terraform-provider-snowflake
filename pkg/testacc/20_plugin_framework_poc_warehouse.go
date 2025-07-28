package testacc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	ResourceMonitor                 types.String                             `tfsdk:"resource_monitor"` // TODO [mux-PR]: identifier type?
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

func marshallWarehousePocPrivateJson(warehouse *sdk.Warehouse) ([]byte, error) {
	warehouseJson := warehousePocPrivateJsonFromWarehouse(warehouse)
	bytes, err := json.Marshal(warehouseJson)
	if err != nil {
		return nil, fmt.Errorf("could not marshal json: %w", err)
	}
	return bytes, nil
}

func (r *WarehouseResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_warehouse_poc"
}

// TODO [mux-PR]: suppress identifier quoting
// TODO [mux-PR]: IgnoreChangeToCurrentSnowflakeValueInShow?
// TODO [mux-PR]: show_output and parameters
// TODO [this PR]: fully_qualified_name
func (r *WarehouseResource) attributes() map[string]schema.Attribute {
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
			// TODO [mux-PR]: identifier validation
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
				// TODO [mux-PR]: how it behaves with renames?
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
		Attributes: r.attributes(),
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

	response.Diagnostics.Append(response.Plan.Set(ctx, &plan)...)
}

// TODO [mux-PR]: from the docs https://developer.hashicorp.com/terraform/plugin/framework/resources/import
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

	// TODO [this PR]: handle parameters
	//warehouseParameters, err := client.Warehouses.ShowParameters(ctx, id)
	// if err != nil {
	//	response.Diagnostics.AddError("Could not read Warehouse PoC parameters", err.Error())
	//	return
	// }

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
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
		// resource_monitor
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

	// TODO [this PR]: added to pass the initial test
	data.FullyQualifiedName = types.StringValue(id.FullyQualifiedName())

	// we can use the existing encoder
	data.Id = types.StringValue(helpers.EncodeResourceIdentifier(id))

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	b, d := r.readAfterCreateOrUpdate(ctx, data, id, &response.State)
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

func (r *WarehouseResource) readAfterCreateOrUpdate(ctx context.Context, data *warehousePocModelV0, id sdk.AccountObjectIdentifier, state *tfsdk.State) ([]byte, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	// TODO [this PR]: merge with read
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

	bytes, err := marshallWarehousePocPrivateJson(warehouse)
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

// TODO [this PR]: add functional test for saving the field always when it is not null in config
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

	_ = warehouseParameters

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

		// TODO [mux-PR]: introduce function like handleExternalChangesToObjectInShow or something similar
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
		// if warehouse.ResourceMonitor != prevValue.ResourceMonitor {
		//	data.ResourceMonitor = types.StringValue(warehouse.ResourceMonitor.Name())
		// }
		if warehouse.EnableQueryAcceleration != prevValue.EnableQueryAcceleration {
			data.EnableQueryAcceleration = types.BoolValue(warehouse.EnableQueryAcceleration)
		}
		if warehouse.QueryAccelerationMaxScaleFactor != prevValue.QueryAccelerationMaxScaleFactor {
			data.QueryAccelerationMaxScaleFactor = types.Int64Value(int64(warehouse.QueryAccelerationMaxScaleFactor))
		}
	}

	bytes, err := marshallWarehousePocPrivateJson(warehouse)
	if err != nil {
		diags.AddError("Could not marshal json", err.Error())
		return diags
	}
	response.Diagnostics.Append(response.Private.SetKey(ctx, privateStateSnowflakeObjectsStateKey, bytes)...)

	// TODO [this PR]: handle warehouse parameters read
	// TODO [this PR]: setStateToValuesFromConfig ?
	// TODO [mux-PR]: show_output and parameters

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
		{
			newId := sdk.NewAccountObjectIdentifier(plan.Name.ValueString())

			err := r.client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{
				NewName: &newId,
			})
			if err != nil {
				response.Diagnostics.AddError("Could not rename warehouse PoC", err.Error())
				return
			}

			plan.Id = types.StringValue(helpers.EncodeResourceIdentifier(id))
			id = newId
		}
	}

	// Batch SET operations and UNSET operations
	set := sdk.WarehouseSet{}
	unset := sdk.WarehouseUnset{}

	errs := errors.Join(
		// name handled in rename
		// TODO [this PR]: unset for warehouse type does not work
		testfunctional.StringEnumAttributeUpdate(plan.WarehouseType, state.WarehouseType, &set.WarehouseType, &unset.WarehouseType),
		// TODO [this PR]: warehouse size unset?
		// TODO [this PR]: WaitForCompletion
		//testfunctional.StringEnumAttributeUpdate(plan.WarehouseSize, state.WarehouseSize, &set.WarehouseSize, &unset.WarehouseSize),
		testfunctional.Int64AttributeUpdate(plan.MaxClusterCount, state.MaxClusterCount, &set.MaxClusterCount, &unset.MaxClusterCount),
		testfunctional.Int64AttributeUpdate(plan.MinClusterCount, state.MinClusterCount, &set.MinClusterCount, &unset.MinClusterCount),
		// TODO [this PR]: unset for scaling policy does not work
		testfunctional.StringEnumAttributeUpdate(plan.ScalingPolicy, state.ScalingPolicy, &set.ScalingPolicy, &unset.ScalingPolicy),
		// TODO [this PR]: unset for auto_suspend does not work
		testfunctional.Int64AttributeUpdate(plan.AutoSuspend, state.AutoSuspend, &set.AutoSuspend, &unset.AutoSuspend),
		// TODO [this PR]: unset for auto_resume does not work
		testfunctional.BooleanAttributeUpdate(plan.AutoResume, state.AutoResume, &set.AutoResume, &unset.AutoResume),
		// resource_monitor
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

	// TODO [this PR]: added to pass the initial test
	plan.FullyQualifiedName = types.StringValue(id.FullyQualifiedName())

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	b, d := r.readAfterCreateOrUpdate(ctx, plan, id, &response.State)
	if d.HasError() {
		response.Diagnostics.Append(d...)
		return
	}
	response.Diagnostics.Append(response.Private.SetKey(ctx, privateStateSnowflakeObjectsStateKey, b)...)
}

// For SDKv2 resources we have a method handling deletion common cases; we can add somethign similar later
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
