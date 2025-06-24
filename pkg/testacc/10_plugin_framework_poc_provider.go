package testacc

import (
	"context"
	"os"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ------ provider interface implementation ------

// Ensure the implementation satisfies the provider.Provider interface.
var _ provider.Provider = &pluginFrameworkPocProvider{}

type pluginFrameworkPocProvider struct {
	// TODO [mux-PR]: fill version automatically like tracking
	version string
}

func (p *pluginFrameworkPocProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "snowflake"
	response.Version = p.version
}

func (p *pluginFrameworkPocProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	// schema needs to match based on https://developer.hashicorp.com/terraform/plugin/framework/migrating/mux#preparedconfig-response-from-multiple-servers
	response.Schema = schema.Schema{
		Attributes: pluginFrameworkPocProviderSchemaV0,
	}
}

func (p *pluginFrameworkPocProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	// TODO [mux-PR]: implement (populate in *gosnowflake.Config)
	// TODO [mux-PR]: us os wrapper
	todoFromEnv := os.Getenv("SNOWFLAKE_TODO")

	var configModel pluginFrameworkPocProviderModelV0

	// Read configuration data into model
	response.Diagnostics.Append(request.Config.Get(ctx, &configModel)...)

	// TODO [mux-PR]: configure other attributes
	var authenticator string
	if !configModel.Authenticator.IsNull() {
		authenticator = configModel.Authenticator.ValueString()
	} else {
		authenticator = todoFromEnv
	}

	if authenticator == "" {
		response.Diagnostics.AddError(
			"TODO summary",
			"TODO details",
		)
	}

	if response.Diagnostics.HasError() {
		return
	}

	// TODO [mux-PR]: try to initialize the client and set it
	providerCtx := &internalprovider.Context{Client: nil}
	// TODO [mux-PR]: set preview_features_enabled
	response.DataSourceData = providerCtx
	response.ResourceData = providerCtx
}

func (p *pluginFrameworkPocProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// TODO [mux-PR]: implement
	}
}

func (p *pluginFrameworkPocProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSomeResource,
	}
}

// ------ convenience ------

func New(version string) provider.Provider {
	return &pluginFrameworkPocProvider{
		version: version,
	}
}
