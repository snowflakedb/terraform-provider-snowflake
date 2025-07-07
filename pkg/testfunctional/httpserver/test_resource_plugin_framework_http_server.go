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
	return &httpServerResource{}
}

type httpServerResource struct {
	serverUrl string
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

func (r *httpServerResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	providerContext, ok := request.ProviderData.(*common.TestProviderContext)
	if !ok {
		response.Diagnostics.AddError("Provider context is broken", "Set up the context correctly in the provider's Configure func.")
		return
	}

	r.serverUrl = providerContext.ServerUrl()
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
	err := common.Post(r.serverUrl, "http_server_example", &exampleRead)
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

	exampleRead := Read{}
	err := common.Get(r.serverUrl, "http_server_example", &exampleRead)
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
