package testacc

import (
	"context"
	"fmt"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// TODO [SNOW-2661409]: check all the places using TestAccProvider directly
	TestAccProvider                 *schema.Provider
	TestAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)

	acceptanceTestsProviderCache *providerInitializationCache[cacheEntry]
)

type cacheEntry struct {
	clientErrorDiag diag.Diagnostics
	providerCtx     *internalprovider.Context
}

func setUpProvider() error {
	acceptanceTestsProviderCache = newProviderInitializationCache[cacheEntry]()

	TestAccProtoV6ProviderFactories, TestAccProvider = providerFactoryUsingCacheReturningProvider("AcceptanceTestDefault")

	return nil
}

var functionsAndProceduresProviderFactory = providerFactoryUsingCache("FunctionsAndProcedures")
var viewsProviderFactory = providerFactoryUsingCache("Views")
var tagsProviderFactory = providerFactoryUsingCache("Tags")
var servicesProviderFactory = providerFactoryUsingCache("Services")
var userPasswordPoliciesProviderFactory = providerFactoryUsingCache("UserPasswordPolicies")
var explicitAccountAdminRoleProviderFactory = providerFactoryUsingCache("ExplicitAccountAdminRole")

// TODO [SNOW-2661409]: secondary account can have also a different configuration, so for now we need to be careful; let's add some hash check for the config or something else to mitigate
var secondaryAccountProviderFactory = providerFactoryUsingCache("SecondaryAccount")

func acceptanceTestsProvider() *schema.Provider {
	p := provider.Provider()
	p.ResourcesMap["snowflake_semantic_view"] = resources.SemanticView()
	p.DataSourcesMap["snowflake_semantic_views"] = datasources.SemanticViews()
	return p
}

// TODO [SNOW-2661409]: we could keep the cache of provider configuration/provider per cache key
func providerFactoryUsingCache(key string) map[string]func() (tfprotov6.ProviderServer, error) {
	factory, _ := providerFactoryUsingCacheReturningProvider(key)
	return factory
}

func providerFactoryUsingCacheReturningProvider(key string) (map[string]func() (tfprotov6.ProviderServer, error), *schema.Provider) {
	p := acceptanceTestsProvider()
	p.ConfigureContextFunc = configureAcceptanceTestProviderWithCacheFunc(key)

	return map[string]func() (tfprotov6.ProviderServer, error){
		"snowflake": func() (tfprotov6.ProviderServer, error) {
			return tf5to6server.UpgradeServer(
				context.Background(),
				p.GRPCProvider,
			)
		},
	}, p
}

// TODO [SNOW-2661409]: check which of the usages wants to really be without cache and which could utilize a dedicated cache entry
func providerFactoryWithoutCache() map[string]func() (tfprotov6.ProviderServer, error) {
	factory, _ := providerFactoryWithoutCacheReturningProvider()
	return factory
}

func providerFactoryWithoutCacheReturningProvider() (map[string]func() (tfprotov6.ProviderServer, error), *schema.Provider) {
	p := acceptanceTestsProvider()
	p.ConfigureContextFunc = configureAcceptanceTestProvider

	return map[string]func() (tfprotov6.ProviderServer, error){
		"snowflake": func() (tfprotov6.ProviderServer, error) {
			return tf5to6server.UpgradeServer(
				context.Background(),
				p.GRPCProvider,
			)
		},
	}, p
}

func configureAcceptanceTestProviderWithCacheFunc(key string) func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		entry := acceptanceTestsProviderCache.getOrInit(key, func() cacheEntry {
			providerCtx, clientErrorDiag := configureAcceptanceTestProvider(ctx, d)
			return cacheEntry{
				providerCtx:     providerCtx.(*internalprovider.Context),
				clientErrorDiag: clientErrorDiag,
			}
		})
		return entry.providerCtx, entry.clientErrorDiag
	}
}

func configureAcceptanceTestProvider(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	accTestLog.Printf("[DEBUG] Initializing acceptance test provider")

	providerCtx, clientErrorDiag := provider.ConfigureProvider(ctx, d)

	if providerCtx != nil && oswrapper.Getenv(fmt.Sprintf("%v", testenvs.EnableAllPreviewFeatures)) == "true" {
		providerCtx.(*internalprovider.Context).EnabledFeatures = previewfeatures.AllPreviewFeatures
	}

	return providerCtx, clientErrorDiag
}
