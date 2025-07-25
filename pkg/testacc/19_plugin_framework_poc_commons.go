package testacc

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type SnowflakeClientEmbeddable struct {
	client *sdk.Client
}

func (r *SnowflakeClientEmbeddable) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	providerContext, ok := request.ProviderData.(*provider.Context)
	if !ok {
		response.Diagnostics.AddError("Provider context is broken", "Set up the context correctly in the provider's Configure func.")
		return
	}

	r.client = providerContext.Client
}
