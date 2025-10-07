package manual_tests

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ManualTestProvider *schema.Provider
	v5Server           tfprotov5.ProviderServer
	v6Server           tfprotov6.ProviderServer
)

func init() {
	log.Println("[DEBUG] Running manual tests provider setup")

	ManualTestProvider = provider.Provider()

	v5Server = ManualTestProvider.GRPCProvider()
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
}

var ManualTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"snowflake": func() (tfprotov6.ProviderServer, error) {
		return v6Server, nil
	},
}
