package testacc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewSomeResource() resource.Resource {
	return &SomeResource{}
}

type SomeResource struct{}

func (r *SomeResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	// TODO [mux-PR]: add method for this logic
	response.TypeName = request.ProviderTypeName + "_some"
}

func (r *SomeResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"todo": schema.StringAttribute{
				Description: "TODO",
				Optional:    true,
			},
		},
	}
}

func (r *SomeResource) Create(_ context.Context, _ resource.CreateRequest, _ *resource.CreateResponse) {
	// TODO [mux-PR]: implement
}

func (r *SomeResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// TODO [mux-PR]: implement
}

func (r *SomeResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// TODO [mux-PR]: implement
}

func (r *SomeResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// TODO [mux-PR]: implement
}
