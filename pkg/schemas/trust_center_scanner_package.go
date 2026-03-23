package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ShowTrustCenterScannerPackageSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"description": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default_schedule": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"state": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schedule": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"notification": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"provider_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"last_enabled_timestamp": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"last_disabled_timestamp": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowTrustCenterScannerPackageSchema

func ScannerPackageToSchema(pkg *sdk.ScannerPackage) map[string]any {
	s := map[string]any{
		"name":             pkg.Name,
		"id":               pkg.Id,
		"description":      pkg.Description,
		"default_schedule": pkg.DefaultSchedule,
		"state":            pkg.State,
		"schedule":         pkg.Schedule,
		"notification":     pkg.Notification,
		"provider_name":    pkg.Provider,
	}
	if pkg.LastEnabledTimestamp != nil {
		s["last_enabled_timestamp"] = *pkg.LastEnabledTimestamp
	}
	if pkg.LastDisabledTimestamp != nil {
		s["last_disabled_timestamp"] = *pkg.LastDisabledTimestamp
	}
	return s
}

var _ = ScannerPackageToSchema
