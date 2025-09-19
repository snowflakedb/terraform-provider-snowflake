package testacc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type SnowflakeRestApiEmbeddable struct {
	client *RestApiPocClient
}

func (r *SnowflakeRestApiEmbeddable) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

	r.client = providerContext.RestApiPocClient
}
