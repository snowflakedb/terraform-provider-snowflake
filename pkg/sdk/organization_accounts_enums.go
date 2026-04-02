package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type OrganizationAccountEdition string

var (
	OrganizationAccountEditionEnterprise       OrganizationAccountEdition = "ENTERPRISE"
	OrganizationAccountEditionBusinessCritical OrganizationAccountEdition = "BUSINESS_CRITICAL"
)

var AllOrganizationAccountEditions = []OrganizationAccountEdition{
	OrganizationAccountEditionEnterprise,
	OrganizationAccountEditionBusinessCritical,
}

func ToOrganizationAccountEdition(s string) (OrganizationAccountEdition, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllOrganizationAccountEditions, OrganizationAccountEdition(s)) {
		return "", fmt.Errorf("invalid organization account edition: %s", s)
	}
	return OrganizationAccountEdition(s), nil
}
