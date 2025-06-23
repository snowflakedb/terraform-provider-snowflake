package testacc

import (
	"context"
	"os"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ------ provider interface implementation ------

// Ensure the implementation satisfies the provider.Provider interface.
var _ provider.Provider = &pluginFrameworkPocProvider{}

type pluginFrameworkPocProvider struct {
	// TODO [mux-PR]: fill version automatically like tracking
	version string
}

type pluginFrameworkPocProviderModelV0 struct {
	Todo types.String `tfsdk:"todo"`
}

func (p *pluginFrameworkPocProvider) Metadata(_ context.Context, _ provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "snowflake"
	response.Version = p.version
}

func (p *pluginFrameworkPocProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	// TODO [mux-PR]: schema needs to match based on https://developer.hashicorp.com/terraform/plugin/framework/migrating/mux#preparedconfig-response-from-multiple-servers
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"todo": schema.StringAttribute{
				Description: "TODO",
				Optional:    true,
			},
		},
	}
}

func (p *pluginFrameworkPocProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	// TODO [mux-PR]: implement (populate in *gosnowflake.Config)
	// TODO [mux-PR]: us os wrapper
	todoFromEnv := os.Getenv("SNOWFLAKE_TODO")

	var configModel pluginFrameworkPocProviderModelV0

	// Read configuration data into model
	response.Diagnostics.Append(request.Config.Get(ctx, &configModel)...)

	var todo string
	if !configModel.Todo.IsNull() {
		todo = configModel.Todo.ValueString()
	} else {
		todo = todoFromEnv
	}

	if todo == "" {
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
