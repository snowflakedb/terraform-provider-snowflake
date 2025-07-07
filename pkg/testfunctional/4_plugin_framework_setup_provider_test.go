package testfunctional_test

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var providerForPluginFrameworkFunctionalTestsFactories = map[string]func() (tfprotov6.ProviderServer, error){
	PluginFrameworkFunctionalTestsProviderName: providerserver.NewProtocol6WithError(New()),
}
