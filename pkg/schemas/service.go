package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeServiceSchema represents output of DESCRIBE query for the single Service.
var DescribeServiceSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"status": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"database_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"schema_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"compute_pool": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"spec": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"dns_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"current_instances": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"target_instances": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"min_ready_instances": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"min_instances": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"max_instances": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"auto_resume": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"external_access_integrations": {
		// Adjusted manually.
		Type:     schema.TypeSet,
		Elem:     &schema.Schema{Type: schema.TypeString},
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
	"resumed_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"suspended_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"auto_suspend_secs": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner_role_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"query_warehouse": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_job": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_async_job": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"spec_digest": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_upgrading": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"managing_object_domain": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"managing_object_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = DescribeServiceSchema

func ServiceDetailsToSchema(service *sdk.ServiceDetails) map[string]any {
	serviceSchema := make(map[string]any)
	serviceSchema["name"] = service.Name
	serviceSchema["status"] = string(service.Status)
	serviceSchema["spec"] = service.Spec
	serviceSchema["database_name"] = service.DatabaseName
	serviceSchema["schema_name"] = service.SchemaName
	serviceSchema["owner"] = service.Owner
	serviceSchema["compute_pool"] = service.ComputePool.Name()
	serviceSchema["dns_name"] = service.DnsName
	serviceSchema["current_instances"] = service.CurrentInstances
	serviceSchema["target_instances"] = service.TargetInstances
	serviceSchema["min_ready_instances"] = service.MinReadyInstances
	serviceSchema["min_instances"] = service.MinInstances
	serviceSchema["max_instances"] = service.MaxInstances
	serviceSchema["auto_resume"] = service.AutoResume
	serviceSchema["external_access_integrations"] = collections.Map(service.ExternalAccessIntegrations, sdk.AccountObjectIdentifier.Name)
	serviceSchema["created_on"] = service.CreatedOn.String()
	serviceSchema["updated_on"] = service.UpdatedOn.String()
	if service.ResumedOn != nil {
		serviceSchema["resumed_on"] = service.ResumedOn.String()
	}
	if service.SuspendedOn != nil {
		serviceSchema["suspended_on"] = service.SuspendedOn.String()
	}
	serviceSchema["auto_suspend_secs"] = service.AutoSuspendSecs
	if service.Comment != nil {
		serviceSchema["comment"] = service.Comment
	}
	serviceSchema["owner_role_type"] = service.OwnerRoleType
	if service.QueryWarehouse != nil {
		serviceSchema["query_warehouse"] = service.QueryWarehouse.Name()
	}
	serviceSchema["is_job"] = service.IsJob
	serviceSchema["is_async_job"] = service.IsAsyncJob
	serviceSchema["spec_digest"] = service.SpecDigest
	serviceSchema["is_upgrading"] = service.IsUpgrading
	if service.ManagingObjectDomain != nil {
		serviceSchema["managing_object_domain"] = service.ManagingObjectDomain
	}
	if service.ManagingObjectName != nil {
		serviceSchema["managing_object_name"] = service.ManagingObjectName
	}
	return serviceSchema
}

var _ = ServiceToSchema
