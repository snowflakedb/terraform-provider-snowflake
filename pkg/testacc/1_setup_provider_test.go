package testacc

import (
	"context"
	"fmt"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
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
)

// TODO [next PRs]: rework this when working on terraform plugin framework PoC
func setUpProvider() error {
	TestAccProvider = provider.Provider()
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
		return fmt.Errorf("cannot upgrade server from proto v5 to proto v6, failing, err: %v", err)
	}

	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"snowflake": func() (tfprotov6.ProviderServer, error) {
			return v6Server, nil
		},
	}

	return nil
}

// TODO [next PRs]: it's currently an exact copy of acceptance.ConfigureProviderWithConfigCache; adjust after moving the tests
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

	if providerCtx != nil && accTestEnabled && oswrapper.Getenv("SF_TF_ACC_TEST_ENABLE_ALL_PREVIEW_FEATURES") == "true" {
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

	if clientErrorDiag.HasError() {
		return nil, clientErrorDiag
	}

	return providerCtx, nil
}
