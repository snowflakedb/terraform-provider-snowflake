package tracking

import (
	"fmt"
	"log"
)

var ProviderVersion = "dev"

func SetProviderVersion(version string) {
	providerVersion := fmt.Sprintf("v%s", version)
	log.Println("[INFO] Setting provider version:", providerVersion)
	ProviderVersion = providerVersion
}
