package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribePostgresInstanceSchema represents output of DESCRIBE query for the single PostgresInstance.
var DescribePostgresInstanceSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"updated_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"host": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"origin": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"privatelink_service_identifier": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"compute_family": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"storage_size_gb": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"postgres_version": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"high_availability": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"authentication_authority": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"state": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"retention_time": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"maintenance_window_start": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"network_policy": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"postgres_settings": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"storage_integration": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func PostgresInstanceDetailsToSchema(details *sdk.PostgresInstanceDetails) map[string]any {
	s := make(map[string]any)
	s["name"] = details.Name
	s["owner"] = details.Owner
	s["owner_role_type"] = details.OwnerRoleType
	s["created_on"] = details.CreatedOn
	s["updated_on"] = details.UpdatedOn
	s["type"] = details.Type
	s["host"] = details.Host
	if details.Origin != nil {
		s["origin"] = *details.Origin
	}
	if details.PrivatelinkServiceIdentifier != nil {
		s["privatelink_service_identifier"] = *details.PrivatelinkServiceIdentifier
	}
	s["compute_family"] = details.ComputeFamily
	s["storage_size_gb"] = details.StorageSizeGb
	s["postgres_version"] = details.PostgresVersion
	s["high_availability"] = details.HighAvailability
	s["authentication_authority"] = details.AuthenticationAuthority
	s["state"] = details.State
	s["retention_time"] = details.RetentionTime
	s["maintenance_window_start"] = details.MaintenanceWindowStart
	if details.Comment != nil {
		s["comment"] = *details.Comment
	}
	if details.NetworkPolicy != nil {
		s["network_policy"] = *details.NetworkPolicy
	}
	if details.PostgresSettings != nil {
		s["postgres_settings"] = *details.PostgresSettings
	}
	if details.StorageIntegration != nil {
		s["storage_integration"] = *details.StorageIntegration
	}
	return s
}
