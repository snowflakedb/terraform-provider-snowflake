package main

import (
	"context"
	"flag"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/internal/tracking"
	oldprovider "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

var version string = "dev" // goreleaser can pass other information to the main package, such as the specific commit
// https://goreleaser.com/cookbooks/using-main.version/

func main() {
	tracking.SetProviderVersion(version)

	ctx := context.Background()

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	upgradedSdkServer, err := tf5to6server.UpgradeServer(
		ctx,
		oldprovider.Provider().GRPCProvider,
	)
	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		// disabled until ready to start using
		//	providerserver.NewProtocol6(provider.New(version)()),
		func() tfprotov6.ProviderServer {
			return upgradedSdkServer
		},
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/snowflakedb/snowflake",
		muxServer.ProviderServer,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
