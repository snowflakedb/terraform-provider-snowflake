package testfunctional_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	providerForSdkV2FunctionalTests          *schema.Provider
	providerForSdkV2FunctionalTestsFactories map[string]func() (tfprotov6.ProviderServer, error)

	v5Server tfprotov5.ProviderServer
	v6Server tfprotov6.ProviderServer
)

func setUpProvidersForFunctionalTests() error {
	providerForSdkV2FunctionalTests = sdkV2FunctionalTestsProvider()

	var err error
	v5Server = providerForSdkV2FunctionalTests.GRPCProvider()
	v6Server, err = tf5to6server.UpgradeServer(
		context.Background(),
		func() tfprotov5.ProviderServer {
			return v5Server
		},
	)
	if err != nil {
		return fmt.Errorf("cannot upgrade server from proto v5 to proto v6, failing, err: %w", err)
	}

	providerForSdkV2FunctionalTestsFactories = map[string]func() (tfprotov6.ProviderServer, error){
		SdkV2FunctionalTestsProviderName: func() (tfprotov6.ProviderServer, error) {
			return v6Server, nil
		},
	}

	return nil
}
