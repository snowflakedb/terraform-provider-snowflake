package sdk

import (
	"fmt"
	"strings"
)

type AccountEdition string

var (
	EditionStandard         AccountEdition = "STANDARD"
	EditionEnterprise       AccountEdition = "ENTERPRISE"
	EditionBusinessCritical AccountEdition = "BUSINESS_CRITICAL"
)

var AllAccountEditions = []AccountEdition{
	EditionStandard,
	EditionEnterprise,
	EditionBusinessCritical,
}

func ToAccountEdition(edition string) (AccountEdition, error) {
	switch typedEdition := AccountEdition(strings.ToUpper(edition)); typedEdition {
	case EditionStandard, EditionEnterprise, EditionBusinessCritical:
		return typedEdition, nil
	default:
		return "", fmt.Errorf("unknown account edition: %s", edition)
	}
}
