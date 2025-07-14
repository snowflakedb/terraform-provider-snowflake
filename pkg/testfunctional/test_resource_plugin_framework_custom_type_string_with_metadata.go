package testfunctional

import (
	"context"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/customtypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigure = &StringWithMetadataResource{}

func NewStringWithMetadataResource() resource.Resource {
	return &StringWithMetadataResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[StringWithMetadataOpts]("string_with_metadata"),
	}
}

type StringWithMetadataResource struct {
	common.HttpServerEmbeddable[StringWithMetadataOpts]
}

type stringWithMetadataResourceModelV0 struct {
	Name        types.String                        `tfsdk:"name"`
	StringValue customtypes.StringWithMetadataValue `tfsdk:"string_value"`
	Id          types.String                        `tfsdk:"id"`
}

type StringWithMetadataOpts struct {
	StringValue *string
}

func (r *StringWithMetadataResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_string_with_metadata"
}

func (r *StringWithMetadataResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name for this resource.",
				Required:    true,
			},
			"string_value": schema.StringAttribute{
				CustomType:  customtypes.StringWithMetadataType{},
				Description: "String value.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *StringWithMetadataResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *StringWithMetadataResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *stringWithMetadataResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	opts := &StringWithMetadataOpts{}
	customtypes.StringWithMetadataAttributeCreate(data.StringValue, &opts.StringValue)

	response.Diagnostics.Append(r.create(opts)...)
	if response.Diagnostics.HasError() {
		return
	}

	meta := customtypes.Metadata{
		FieldA: time.Now().Format(time.RFC3339),
	}
	data.StringValue = customtypes.StringWithMetadataValue{
		StringValue: types.StringValue(data.StringValue.ValueString()),
		Metadata:    meta,
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *StringWithMetadataResource) create(opts *StringWithMetadataOpts) diag.Diagnostics {
	diags := diag.Diagnostics{}

	err := r.HttpServerEmbeddable.Post(*opts)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *StringWithMetadataResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *stringWithMetadataResourceModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(r.readStringWithMetadataResource(data)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *StringWithMetadataResource) readStringWithMetadataResource(data *stringWithMetadataResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	opts, err := r.HttpServerEmbeddable.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else if opts.StringValue != nil {
		meta := customtypes.Metadata{
			FieldA: time.Now().Format(time.RFC3339),
		}
		data.StringValue = customtypes.StringWithMetadataValue{
			StringValue: types.StringValue(data.StringValue.ValueString()),
			Metadata:    meta,
		}
	}
	return diags
}

func (r *StringWithMetadataResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *StringWithMetadataResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
