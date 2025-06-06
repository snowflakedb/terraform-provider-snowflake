package acceptance

import (
	"context"
	"log"

	internalprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
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
	TestAccProvider *schema.Provider
	v5Server        tfprotov5.ProviderServer
	v6Server        tfprotov6.ProviderServer

	configureClientErrorDiag diag.Diagnostics
	configureProviderCtx     *internalprovider.Context
)

func init() {
	log.Println("[DEBUG] Running init from old acceptance tests setup")

	TestAccProvider = provider.Provider()
	// TODO [SNOW-2054208]: improve during the package extraction
	TestAccProvider.ResourcesMap["snowflake_object_renaming"] = resources.ObjectRenamingListsAndSets()
	TestAccProvider.ResourcesMap["snowflake_test_resource_data_type_diff_handling"] = resources.TestResourceDataTypeDiffHandling()
	TestAccProvider.ResourcesMap["snowflake_test_resource_data_type_diff_handling_list"] = resources.TestResourceDataTypeDiffHandlingList()
	TestAccProvider.ConfigureContextFunc = ConfigureProviderWithConfigCache

	v5Server = TestAccProvider.GRPCProvider()
	var err error
	v6Server, err = tf5to6server.UpgradeServer(
		context.Background(),
		func() tfprotov5.ProviderServer {
			return v5Server
		},
	)
	if err != nil {
		log.Panicf("Cannot upgrade server from proto v5 to proto v6, failing, err: %v", err)
	}
	_ = testAccProtoV6ProviderFactoriesNew
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"snowflake": func() (tfprotov6.ProviderServer, error) {
		return v6Server, nil
	},
}

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

func ConfigureProviderWithConfigCache(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
	accTestEnabled, err := oswrapper.GetenvBool("TF_ACC")
	if err != nil {
		accTestEnabled = false
		log.Printf("TF_ACC environmental variable has incorrect format: %v, using %v as a default value", err, accTestEnabled)
	}
	configureClientOnceEnabled, err := oswrapper.GetenvBool("SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE")
	if err != nil {
		configureClientOnceEnabled = false
		log.Printf("SF_TF_ACC_TEST_CONFIGURE_CLIENT_ONCE environmental variable has incorrect format: %v, using %v as a default value", err, configureClientOnceEnabled)
	}
	// hacky way to speed up our acceptance tests
	if accTestEnabled && configureClientOnceEnabled {
		log.Printf("[DEBUG] Returning cached provider configuration result")
		if configureProviderCtx != nil {
			log.Printf("[DEBUG] Returning cached provider configuration context")
			return configureProviderCtx, nil
		}
		if configureClientErrorDiag.HasError() {
			log.Printf("[DEBUG] Returning cached provider configuration error")
			return nil, configureClientErrorDiag
		}
	}
	log.Printf("[DEBUG] No cached provider configuration found or caching is not enabled; configuring a new provider")

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
