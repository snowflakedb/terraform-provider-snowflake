package testacc

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
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
	TestAccProvider.DataSourcesMap["snowflake_semantic_views"] = datasources.SemanticViews()
	TestAccProvider.ConfigureContextFunc = configureProviderWithConfigCacheFunc("AcceptanceTestDefault")

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

// TODO [SNOW-2312385]: add dedicated factories (authentication, views tests, functions tests, procedure tests, secondary tests), address all SNOW-2324320
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
		// TODO [SNOW-2312385]: lock access to cache map
		// check if we cached initialized provider context with the key already
		if cached, ok := providerInitializationCache[key]; ok {
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

		if providerCtx != nil && oswrapper.Getenv(fmt.Sprintf("%v", testenvs.EnableAllPreviewFeatures)) == "true" {
			providerCtx.(*internalprovider.Context).EnabledFeatures = previewfeatures.AllPreviewFeatures
		}

		providerInitializationCache[key] = cacheEntry{
			providerCtx:     providerCtx.(*internalprovider.Context),
			clientErrorDiag: clientErrorDiag,
		}

		// TODO [SNOW-2312385]: what do we want to do with this? - get from cache?
		if v, ok := providerCtx.(*internalprovider.Context); ok {
			lastConfiguredProviderContext = v
		}

		return providerCtx, clientErrorDiag
	}
}

func providerFactoryWithoutCache() map[string]func() (tfprotov6.ProviderServer, error) {
	p := provider.Provider()
	p.ConfigureContextFunc = configureProviderWithoutCache

	return map[string]func() (tfprotov6.ProviderServer, error){
		"snowflake": func() (tfprotov6.ProviderServer, error) {
			return tf5to6server.UpgradeServer(
				context.Background(),
				p.GRPCProvider,
			)
		},
	}
}

func configureProviderWithoutCache(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	accTestLog.Printf("[DEBUG] TODO")

	providerCtx, clientErrorDiag := provider.ConfigureProvider(ctx, d)

	if providerCtx != nil && oswrapper.Getenv(fmt.Sprintf("%v", testenvs.EnableAllPreviewFeatures)) == "true" {
		providerCtx.(*internalprovider.Context).EnabledFeatures = previewfeatures.AllPreviewFeatures
	}

	// TODO [SNOW-2312385]: what do we want to do with this when used without cache?
	if v, ok := providerCtx.(*internalprovider.Context); ok {
		lastConfiguredProviderContext = v
	}

	return providerCtx, clientErrorDiag
}
