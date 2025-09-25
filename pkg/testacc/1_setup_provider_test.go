package testacc

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	TestAccProvider                 *schema.Provider
	TestAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)

	v5Server tfprotov5.ProviderServer
	v6Server tfprotov6.ProviderServer

	configureClientErrorDiag diag.Diagnostics
	configureProviderCtx     *internalprovider.Context

	providerInitializationCache map[string]cacheEntry

	// temporary unsafe way to get the last configuration for the provider (to verify in tests);
	// should be used with caution as it is not prepared for the parallel tests
	// should be replaced in the future (e.g. map with test name as key)
	lastConfiguredProviderContext *internalprovider.Context
)

type cacheEntry struct {
	clientErrorDiag diag.Diagnostics
	providerCtx     *internalprovider.Context
}

// TODO [SNOW-2312385]: rework this when improving the caching logic
func setUpProvider() error {
	providerInitializationCache = make(map[string]cacheEntry)

	TestAccProvider = provider.Provider()
	TestAccProvider.ResourcesMap["snowflake_semantic_view"] = resources.SemanticView()
	TestAccProvider.ConfigureContextFunc = configureProviderWithConfigCache

	var err error
	v5Server = TestAccProvider.GRPCProvider()
	v6Server, err = tf5to6server.UpgradeServer(
		context.Background(),
		func() tfprotov5.ProviderServer {
			return v5Server
		},
	)
	if err != nil {
		return fmt.Errorf("cannot upgrade server from proto v5 to proto v6, failing, err: %w", err)
	}

	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"snowflake": func() (tfprotov6.ProviderServer, error) {
			return v6Server, nil
		},
	}
	_ = testAccProtoV6ProviderFactoriesNew

	return nil
}

// TODO [SNOW-2298291]: investigate this (it was moved from the old testing.go file)
// if we do not reuse the created objects there is no `Previously configured provider being re-configured.` warning
// currently left for possible usage after other improvements
var testAccProtoV6ProviderFactoriesNew = map[string]func() (tfprotov6.ProviderServer, error){
	"snowflake": func() (tfprotov6.ProviderServer, error) {
		return tf5to6server.UpgradeServer(
			context.Background(),
			provider.Provider().GRPCProvider,
		)
	},
}

// TODO [SNOW-2312385]: remove this function and use configureProviderWithConfigCacheFunc
func configureProviderWithConfigCache(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
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
		accTestLog.Printf("[DEBUG] Returning cached provider configuration result")
		if configureProviderCtx != nil {
			accTestLog.Printf("[DEBUG] Returning cached provider configuration context")
			return configureProviderCtx, nil
		}
		if configureClientErrorDiag.HasError() {
			accTestLog.Printf("[DEBUG] Returning cached provider configuration error")
			return nil, configureClientErrorDiag
		}
	}
	accTestLog.Printf("[DEBUG] No cached provider configuration found or caching is not enabled; configuring a new provider")

	providerCtx, clientErrorDiag := provider.ConfigureProvider(ctx, d)

	if providerCtx != nil && accTestEnabled && oswrapper.Getenv(fmt.Sprintf("%v", testenvs.EnableAllPreviewFeatures)) == "true" {
		providerCtx.(*internalprovider.Context).EnabledFeatures = previewfeatures.AllPreviewFeatures
	}

	// needed for tests verifying different provider setups
	if accTestEnabled && configureClientOnceEnabled {
		configureProviderCtx = providerCtx.(*internalprovider.Context)
		configureClientErrorDiag = clientErrorDiag
	} else {
		configureProviderCtx = nil
		configureClientErrorDiag = make(diag.Diagnostics, 0)
	}
	if v, ok := providerCtx.(*internalprovider.Context); ok {
		lastConfiguredProviderContext = v
	}

	if clientErrorDiag.HasError() {
		return nil, clientErrorDiag
	}

	return providerCtx, nil
}

var taskDedicatedProviderFactory = providerFactoryUsingCache("task")

// TODO [SNOW-2312385]: we could keep the cache of provider per cache key
func providerFactoryUsingCache(key string) map[string]func() (tfprotov6.ProviderServer, error) {
	p := provider.Provider()
	p.ConfigureContextFunc = configureProviderWithConfigCacheFunc(key)

	return map[string]func() (tfprotov6.ProviderServer, error){
		"snowflake": func() (tfprotov6.ProviderServer, error) {
			return tf5to6server.UpgradeServer(
				context.Background(),
				p.GRPCProvider,
			)
		},
	}
}

func configureProviderWithConfigCacheFunc(key string) func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		// TODO [SNOW-2312385]: do we still keep this guard if this configure func override is only used in the acceptance tests already?
		accTestEnabled, err := oswrapper.GetenvBool("TF_ACC")
		if err != nil {
			accTestEnabled = false
			accTestLog.Printf("[ERROR] TF_ACC environmental variable has incorrect format: %v, using %v as a default value", err, accTestEnabled)
		}

		// check if we cached initialized provider context with the key already
		if cached, ok := providerInitializationCache[key]; accTestEnabled && ok {
			accTestLog.Printf("[DEBUG] Returning cached provider configuration result for key %s", key)
			if cached.providerCtx != nil {
				accTestLog.Printf("[DEBUG] Returning cached provider configuration context")
				return cached.providerCtx, nil
			} else if cached.clientErrorDiag.HasError() {
				accTestLog.Printf("[DEBUG] Returning cached provider configuration error")
				return nil, cached.clientErrorDiag
			}
		}
		accTestLog.Printf("[DEBUG] No cached provider configuration found for key %s or caching is not enabled; configuring a new provider", key)

		providerCtx, clientErrorDiag := provider.ConfigureProvider(ctx, d)

		if providerCtx != nil && accTestEnabled && oswrapper.Getenv(fmt.Sprintf("%v", testenvs.EnableAllPreviewFeatures)) == "true" {
			providerCtx.(*internalprovider.Context).EnabledFeatures = previewfeatures.AllPreviewFeatures
		}

		providerInitializationCache[key] = cacheEntry{
			providerCtx:     providerCtx.(*internalprovider.Context),
			clientErrorDiag: clientErrorDiag,
		}

		// TODO [SNOW-2312385]: what do we want to do with this?
		if v, ok := providerCtx.(*internalprovider.Context); ok {
			lastConfiguredProviderContext = v
		}

		return providerCtx, clientErrorDiag
	}
}
