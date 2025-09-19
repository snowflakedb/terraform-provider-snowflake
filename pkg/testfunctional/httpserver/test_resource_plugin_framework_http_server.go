package httpserver

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigure = &httpServerResource{}

func NewHttpServerResource() resource.Resource {
	return &httpServerResource{
		HttpServerEmbeddable: *common.NewHttpServerEmbeddable[Read]("http_server_example"),
	}
}

type httpServerResource struct {
	common.HttpServerEmbeddable[Read]
}

type Read struct {
	Msg string `json:"msg,omitempty"`
}

type httpServerResourceModelV0 struct {
	Name    types.String `tfsdk:"name"`
	Id      types.String `tfsdk:"id"`
	Message types.String `tfsdk:"message"`
}

func (r *httpServerResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_http_server"
}

func (r *httpServerResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name for this resource.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"message": schema.StringAttribute{
				Computed:    true,
				Description: "Externally settable value.",
			},
		},
	}
}

func (r *httpServerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data *httpServerResourceModelV0
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	name := data.Name.ValueString()
	id := sdk.NewAccountObjectIdentifier(name)
	data.Id = types.StringValue(id.FullyQualifiedName())

	response.Diagnostics.Append(r.create()...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(r.read(data)...)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *httpServerResource) create() diag.Diagnostics {
	diags := diag.Diagnostics{}

	exampleRead := Read{
		Msg: "set through resource",
	}
	err := r.Post(exampleRead)
	if err != nil {
		diags.AddError("Could not create resource", err.Error())
	}
	return diags
}

func (r *httpServerResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data *httpServerResourceModelV0
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)

	response.Diagnostics.Append(r.read(data)...)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *httpServerResource) read(data *httpServerResourceModelV0) diag.Diagnostics {
	diags := diag.Diagnostics{}

	exampleRead, err := r.Get()
	if err != nil {
		diags.AddError("Could not read resources state", err.Error())
	} else {
		data.Message = types.StringValue(exampleRead.Msg)
	}
	return diags
}

func (r *httpServerResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *httpServerResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
