package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ShowTrustCenterScannerSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"short_description": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"description": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"scanner_package_id": {
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
	"last_scan_timestamp": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowTrustCenterScannerSchema

func ScannerToSchema(scanner *sdk.Scanner) map[string]any {
	s := map[string]any{
		"name":               scanner.Name,
		"id":                 scanner.Id,
		"short_description":  scanner.ShortDescription,
		"description":        scanner.Description,
		"scanner_package_id": scanner.ScannerPackageId,
		"state":              scanner.State,
		"schedule":           scanner.Schedule,
		"notification":       scanner.Notification,
	}
	if scanner.LastScanTimestamp != nil {
		s["last_scan_timestamp"] = *scanner.LastScanTimestamp
	}
	return s
}

var _ = ScannerToSchema
