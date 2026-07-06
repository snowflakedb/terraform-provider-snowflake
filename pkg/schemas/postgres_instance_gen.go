package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowPostgresInstanceSchema represents output of SHOW query for the single PostgresInstance.
var ShowPostgresInstanceSchema = map[string]*schema.Schema{
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
	"origin": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"host": {
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
	"authentication_authority": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"storage_size": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"postgres_version": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"postgres_settings": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_ha": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"retention_time": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"state": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowPostgresInstanceSchema

func PostgresInstanceToSchema(postgresInstance *sdk.PostgresInstance) map[string]any {
	s := make(map[string]any)
	s["name"] = postgresInstance.Name
	s["owner"] = postgresInstance.Owner
	s["owner_role_type"] = postgresInstance.OwnerRoleType
	s["created_on"] = postgresInstance.CreatedOn.String()
	s["updated_on"] = postgresInstance.UpdatedOn.String()
	s["type"] = postgresInstance.Type
	if postgresInstance.Origin != nil {
		s["origin"] = *postgresInstance.Origin
	}
	if postgresInstance.Host != nil {
		s["host"] = *postgresInstance.Host
	}
	if postgresInstance.PrivatelinkServiceIdentifier != nil {
		s["privatelink_service_identifier"] = *postgresInstance.PrivatelinkServiceIdentifier
	}
	s["compute_family"] = postgresInstance.ComputeFamily
	s["authentication_authority"] = postgresInstance.AuthenticationAuthority
	s["storage_size"] = postgresInstance.StorageSize
	s["postgres_version"] = postgresInstance.PostgresVersion
	if postgresInstance.PostgresSettings != nil {
		s["postgres_settings"] = *postgresInstance.PostgresSettings
	}
	s["is_ha"] = postgresInstance.IsHighlyAvailable
	s["retention_time"] = postgresInstance.RetentionTime
	s["state"] = string(postgresInstance.State)
	if postgresInstance.Comment != nil {
		s["comment"] = *postgresInstance.Comment
	}
	return s
}

var _ = PostgresInstanceToSchema
