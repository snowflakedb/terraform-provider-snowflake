package testacc

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Context is extended copy of provider.Context
type Context struct {
	Client           *sdk.Client
	RestApiPocClient *RestApiPocClient
	EnabledFeatures  []string
}

const privateStateSnowflakeObjectsStateKey = "state_in_snowflake"

type SnowflakeClientEmbeddable struct {
	client *sdk.Client
}

func (r *SnowflakeClientEmbeddable) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	providerContext, ok := request.ProviderData.(*Context)
	if !ok {
		response.Diagnostics.AddError("Provider context is broken", "Set up the context correctly in the provider's Configure func.")
		return
	}

	if providerContext.Client == nil {
		response.Diagnostics.AddError("Snowflake client cannot be null", "Set up the context correctly in the provider's Configure func.")
		return
	}

	r.client = providerContext.Client
}

type fullyQualifiedNameModelEmbeddable struct {
	FullyQualifiedName types.String `tfsdk:"fully_qualified_name"`
}

func GetFullyQualifiedNameResourceSchema() schema.Attribute {
	return schema.StringAttribute{
		Computed:    true,
		Description: schemas.FullyQualifiedNameSchema.Description,
		// TODO [SNOW-2298108]: decide what should be the logic behind fully_qualified_name
		// PlanModifiers: []planmodifier.String{
		//	stringplanmodifier.UseStateForUnknown(),
		// },
	}
}
