package testacc

import (
	"context"
	"fmt"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

func envNameFieldDescription(description, envName string) string {
	return fmt.Sprintf("%s Can also be sourced from the `%s` environment variable.", description, envName)
}

func (p *pluginFrameworkPocProvider) Schema(_ context.Context, _ provider.SchemaRequest, response *provider.SchemaResponse) {
	// schema needs to match based on https://developer.hashicorp.com/terraform/plugin/framework/migrating/mux#preparedconfig-response-from-multiple-servers
	response.Schema = schema.Schema{
		Attributes: pluginFrameworkPocProviderSchemaV0,
		Blocks: map[string]schema.Block{
			"token_accessor": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"token_endpoint": schema.StringAttribute{
							Description: envNameFieldDescription("The token endpoint for the OAuth provider e.g. https://{yourDomain}/oauth/token when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorTokenEndpoint),
							Required:    true,
							Sensitive:   true,
						},
						"refresh_token": schema.StringAttribute{
							Description: envNameFieldDescription("The refresh token for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorRefreshToken),
							Required:    true,
							Sensitive:   true,
						},
						"client_id": schema.StringAttribute{
							Description: envNameFieldDescription("The client ID for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorClientId),
							Required:    true,
							Sensitive:   true,
						},
						"client_secret": schema.StringAttribute{
							Description: envNameFieldDescription("The client secret for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorClientSecret),
							Required:    true,
							Sensitive:   true,
						},
						"redirect_uri": schema.StringAttribute{
							Description: envNameFieldDescription("The redirect URI for the OAuth provider when using a refresh token to renew access token.", snowflakeenvs.TokenAccessorRedirectUri),
							Required:    true,
							Sensitive:   true,
						},
					},
				},
			},
		},
	}
}

func (p *pluginFrameworkPocProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	// TODO [mux-PR]: implement (populate in *gosnowflake.Config)
	// TODO [mux-PR]: use os wrapper
	// TODO [mux-PR]: handle envs
	// todoFromEnv := os.Getenv("SNOWFLAKE_TODO")

	var configModel pluginFrameworkPocProviderModelV0

	// Read configuration data into model
	response.Diagnostics.Append(request.Config.Get(ctx, &configModel)...)

	// TODO [mux-PR]: configure attributes
	// var authenticator string
	// if !configModel.Authenticator.IsNull() {
	// 	authenticator = configModel.Authenticator.ValueString()
	// } else {
	// 	authenticator = todoFromEnv
	// }
	//
	// if authenticator == "" {
	// 	 response.Diagnostics.AddError(
	//	 	"TODO summary",
	//		"TODO details",
	//	 )
	// }

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
		NewWarehousePocResource,
	}
}

// ------ convenience ------

func New(version string) provider.Provider {
	return &pluginFrameworkPocProvider{
		version: version,
	}
}
