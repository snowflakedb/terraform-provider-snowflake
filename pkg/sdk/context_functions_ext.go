package sdk

import (
	"encoding/json"
	"fmt"
	"strings"
)

func prefixHTTPS(host string) string {
	return fmt.Sprintf("https://%s", host)
}

func (r currentSecondaryRolesDBRow) additionalConvert(result *CurrentSecondaryRoles) error {
	jsonRoles := &struct {
		Roles string `json:"roles"`
		Value string `json:"value"`
	}{}
	if err := json.Unmarshal([]byte(r.CurrentRoles), jsonRoles); err != nil {
		return err
	}

	var roles []AccountObjectIdentifier
	if len(jsonRoles.Roles) > 0 {
		for role := range strings.SplitSeq(jsonRoles.Roles, ",") {
			roles = append(roles, NewAccountObjectIdentifier(role))
		}
	}

	var value SecondaryRoleOption
	switch jsonRoles.Value {
	case "ALL":
		value = SecondaryRoleOptionAll
	default:
		value = SecondaryRoleOptionNone
	}
	result.Roles = roles
	result.Value = value
	return nil
}

func (acc *CurrentSessionDetails) AccountURL() (string, error) {
	region := acc.Region
	// CURRENT_REGION() may return a prefixed form like "PUBLIC.AWS_US_EAST_1"; use the last segment.
	if idx := strings.LastIndex(region, "."); idx >= 0 {
		region = region[idx+1:]
	}
	if regionID, ok := regionMapping[strings.ToLower(region)]; ok {
		accountID := acc.Account
		if len(regionID) > 0 {
			accountID = fmt.Sprintf("%s.%s", accountID, regionID)
		}
		return fmt.Sprintf("https://%s.%s", accountID, getDomainBasedOnRegion(regionID)), nil
	}
	return "", fmt.Errorf("failed to map Snowflake account region %s to a region_id", acc.Region)
}

const (
	defaultDomain = "snowflakecomputing.com"
	cnDomain      = "snowflakecomputing.cn"
)

func getDomainBasedOnRegion(region string) string {
	if strings.HasPrefix(strings.ToLower(region), "cn-") {
		return cnDomain
	}
	return defaultDomain
}
