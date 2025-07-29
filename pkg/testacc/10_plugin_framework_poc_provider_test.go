package testacc

import (
	"context"
	"fmt"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	sdkV2Provider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// our test acc needed variables
var (
	configurePluginFrameworkProviderCtx     *internalprovider.Context
	configureClientErrorPluginFrameworkDiag diag.Diagnostics
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

// The logic for caching is based on the caching we have for the current acceptance tests set.
func (p *pluginFrameworkPocProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	accTestEnabled, err := oswrapper.GetenvBool("TF_ACC")
	if err != nil {
		accTestEnabled = false
		accTestLog.Printf("[ERROR] TF_ACC environmental variable has incorrect format: %v, using %v as a default value", err, accTestEnabled)
	}
	configureClientOnceEnabled, err := oswrapper.GetenvBool("SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE")
	if err != nil {
		configureClientOnceEnabled = false
		accTestLog.Printf("[ERROR] SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE environmental variable has incorrect format: %v, using %v as a default value", err, configureClientOnceEnabled)
	}
	// hacky way to speed up our acceptance tests
	if accTestEnabled && configureClientOnceEnabled {
		accTestLog.Printf("[DEBUG] Returning cached terraform plugin framework PoC provider configuration result")
		if configurePluginFrameworkProviderCtx != nil {
			accTestLog.Printf("[DEBUG] Returning cached terraform plugin framework PoC provider configuration context")
			response.DataSourceData = configurePluginFrameworkProviderCtx
			response.ResourceData = configurePluginFrameworkProviderCtx
			return
		}
		if configureClientErrorPluginFrameworkDiag.HasError() {
			accTestLog.Printf("[DEBUG] Returning cached terraform plugin framework PoC provider configuration error")
			response.Diagnostics.Append(configureClientErrorPluginFrameworkDiag...)
			return
		}
	}
	accTestLog.Printf("[DEBUG] No cached terraform plugin framework PoC provider configuration found or caching is not enabled; configuring a new provider")

	providerCtx, clientErrorDiag := p.configureWithoutCache(ctx, request, response)
	if clientErrorDiag.HasError() {
		response.Diagnostics.Append(clientErrorDiag...)
	}

	if providerCtx != nil && accTestEnabled && oswrapper.Getenv(fmt.Sprintf("%v", testenvs.EnableAllPreviewFeatures)) == "true" {
		providerCtx.EnabledFeatures = previewfeatures.AllPreviewFeatures
	}

	// needed for tests verifying different provider setups
	if accTestEnabled && configureClientOnceEnabled {
		configurePluginFrameworkProviderCtx = providerCtx
		configureClientErrorPluginFrameworkDiag = clientErrorDiag
	} else {
		configurePluginFrameworkProviderCtx = nil
		configureClientErrorPluginFrameworkDiag = make(diag.Diagnostics, 0)
	}
	// no last configured provider
}

func (p *pluginFrameworkPocProvider) configureWithoutCache(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) (*internalprovider.Context, diag.Diagnostics) {
	var configModel pluginFrameworkPocProviderModelV0
	diags := diag.Diagnostics{}

	// Read configuration data into model
	diags.Append(request.Config.Get(ctx, &configModel)...)
	if diags.HasError() {
		return nil, diags
	}

	config, err := p.getDriverConfigFromTerraform(configModel)
	if err != nil {
		diags.AddError("Could not read the Terraform config", err.Error())
		return nil, diags
	}

	// TODO [mux-PR]: handle skip_toml_file_permission_verification and use_legacy_toml_file
	if profile := getProfile(configModel); profile != "" {
		tomlConfig, err := sdkV2Provider.GetDriverConfigFromTOML(profile, false, false)
		if err != nil {
			diags.AddError("Could not read the Toml config", err.Error())
			return nil, diags
		}
		config = sdk.MergeConfig(config, tomlConfig)
	}

	providerCtx := &internalprovider.Context{}
	if client, err := sdk.NewClient(config); err != nil {
		diags.AddError("Could not initialize client", err.Error())
		return nil, diags
	} else {
		providerCtx.Client = client
	}

	// TODO [mux-PR]: set preview_features_enabled
	response.DataSourceData = providerCtx
	response.ResourceData = providerCtx

	return providerCtx, nil
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
