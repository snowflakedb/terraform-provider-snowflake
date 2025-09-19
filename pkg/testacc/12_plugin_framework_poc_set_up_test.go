package testacc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

var testAccProtoV6ProviderFactoriesWithPluginPoc map[string]func() (tfprotov6.ProviderServer, error)

func init() {
	// based on https://developer.hashicorp.com/terraform/plugin/framework/migrating/mux#protocol-version-6
	testAccProtoV6ProviderFactoriesWithPluginPoc = map[string]func() (tfprotov6.ProviderServer, error){
		"snowflake": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()

			upgradedSdkServer, err := tf5to6server.UpgradeServer(
				ctx,
				// using the test acc one, as it has the modified configure method (cache)
				TestAccProvider.GRPCProvider,
			)
			if err != nil {
				return nil, err
			}

			providers := []func() tfprotov6.ProviderServer{
				providerserver.NewProtocol6(New("dev")),
				func() tfprotov6.ProviderServer {
					return upgradedSdkServer
				},
			}

			muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}
