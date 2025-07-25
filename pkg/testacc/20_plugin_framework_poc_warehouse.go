package testacc

import (
	"context"
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
			Description: existingWarehouseSchema[string(sdk.WarehouseParameterMaxConcurrencyLevel)].Description,
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},
		strings.ToLower(string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds)): schema.Int64Attribute{
			Description: existingWarehouseSchema[string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds)].Description,
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		strings.ToLower(string(sdk.WarehouseParameterStatementTimeoutInSeconds)): schema.Int64Attribute{
			Description: existingWarehouseSchema[string(sdk.WarehouseParameterStatementTimeoutInSeconds)].Description,
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

	response.Diagnostics.Append(r.readAfterCreateOrUpdate(data)...)
	if response.Diagnostics.HasError() {
		return
	}

	// we can use the existing encoder
	data.Id = types.StringValue(helpers.EncodeResourceIdentifier(id))
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *WarehouseResource) create(ctx context.Context, id sdk.AccountObjectIdentifier, opts *sdk.CreateWarehouseOptions) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.client.Warehouses.Create(ctx, id, opts)
	if err != nil {
		diags.AddError("Could not create warehouse PoC", err.Error())
	}

	return diags
}

func (r *WarehouseResource) readAfterCreateOrUpdate(data *warehousePocModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	// TODO [this PR]: read

	return diags
}

func (r *WarehouseResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *warehousePocModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	id, err := sdk.ParseAccountObjectIdentifier(data.Id.String())
	if err != nil {
		response.Diagnostics.AddError("Could not read ID in warehouse PoC", err.Error())
		return
	}
	response.Diagnostics.Append(r.read(ctx, data, id)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

// TODO [this PR]: add functional test for saving the field always when it is not null in config
func (r *WarehouseResource) read(ctx context.Context, data *warehousePocModelV0, id sdk.AccountObjectIdentifier) diag.Diagnostics {
	diags := diag.Diagnostics{}

	client := r.client
	w, err := client.Warehouses.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			// TODO [this PR]: mark as removed (old d.SetId(""))
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

	_ = w
	_ = warehouseParameters

	data.FullyQualifiedName = types.StringValue(id.FullyQualifiedName())

	// TODO [this PR]: handle external changes
	// TODO [this PR]: setStateToValuesFromConfig ?
	// TODO [this PR]: handle warehouse parameters read
	// TODO [mux-PR]: show_output and parameters

	return diags
}

func (r *WarehouseResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state *warehousePocModelV0

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	id, err := sdk.ParseAccountObjectIdentifier(state.Id.String())
	if err != nil {
		response.Diagnostics.AddError("Could not read ID in warehouse PoC", err.Error())
		return
	}

	// TODO [this PR]: handle rename

	// Batch SET operations and UNSET operations
	set := sdk.WarehouseSet{}
	unset := sdk.WarehouseUnset{}

	errs := errors.Join(
		// name handled in rename
		// TODO [this PR]: unset for warehouse type does not work
		testfunctional.StringEnumAttributeUpdateWithUnset(plan.WarehouseType, state.WarehouseType, &set.WarehouseType, &unset.WarehouseType),
		// TODO [this PR]: warehouse size unset?
		// TODO [this PR]: WaitForCompletion
		//testfunctional.StringEnumAttributeUpdateWithUnset(plan.WarehouseSize, state.WarehouseSize, &set.WarehouseSize, &unset.WarehouseSize),
		// max_cluster_count
		// min_cluster_count
		// TODO [this PR]: unset for scaling policy does not work
		testfunctional.StringEnumAttributeUpdateWithUnset(plan.ScalingPolicy, state.ScalingPolicy, &set.ScalingPolicy, &unset.ScalingPolicy),
		// auto_suspend
		// auto_resume
		// resource_monitor
		// comment
		// enable_query_acceleration
		// query_acceleration_max_scale_factor

		// max_concurrency_level
		// statement_queued_timeout_in_seconds
		// statement_timeout_in_seconds
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

	response.Diagnostics.Append(r.readAfterCreateOrUpdate(plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *WarehouseResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// TODO [this PR]: implement
}
