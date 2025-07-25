package testacc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewWarehousePocResource() resource.Resource {
	return &WarehouseResource{}
}

type WarehouseResource struct{}

type warehousePocModelV0 struct {
	Name types.String `tfsdk:"name"`
	Id   types.String `tfsdk:"id"`
}

func (r *WarehouseResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_warehouse_poc"
}

func (r *WarehouseResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "",
				PlanModifiers: []planmodifier.String{
					// TODO [mux-PR]: how it behaves with renames?
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *WarehouseResource) Create(_ context.Context, _ resource.CreateRequest, _ *resource.CreateResponse) {
}

func (r *WarehouseResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// TODO [mux-PR]: implement
}

func (r *WarehouseResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// TODO [mux-PR]: implement
}

func (r *WarehouseResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// TODO [mux-PR]: implement
}
