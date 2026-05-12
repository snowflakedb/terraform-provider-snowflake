package testacc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

var testAccProtoV6ProviderFactoriesWithPluginPoc map[string]func() (tfprotov6.ProviderServer, error)

const TerraformPluginFrameworkPocDefaultCacheKey = "TerraformPluginFrameworkPoC"

func init() {
	// based on https://developer.hashicorp.com/terraform/plugin/framework/migrating/mux#protocol-version-6
	testAccProtoV6ProviderFactoriesWithPluginPoc = providerFactoryPluginPocUsingCache(TerraformPluginFrameworkPocDefaultCacheKey)
}

func providerFactoryPluginPocUsingCache(key string) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"snowflake": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()

			// creating a separate cache for all plugin framework tests
			p, err := providerFactoryUsingCache(key)["snowflake"]()
			if err != nil {
				return nil, err
			}

			providers := []func() tfprotov6.ProviderServer{
				providerserver.NewProtocol6(NewWithCacheKey("dev", key)),
				func() tfprotov6.ProviderServer {
					return p
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
