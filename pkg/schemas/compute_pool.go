package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeComputePoolSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"state": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"min_nodes": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"max_nodes": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"instance_family": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"num_services": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"num_jobs": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"auto_suspend_secs": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"auto_resume": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"active_nodes": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"idle_nodes": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"target_nodes": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"resumed_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"updated_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"owner": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_exclusive": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"application": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"error_code": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"status_message": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func ComputePoolDetailsToSchema(computePool sdk.ComputePoolDetails) map[string]any {
	computePoolSchema := make(map[string]any)
	computePoolSchema["name"] = computePool.Name
	computePoolSchema["state"] = string(computePool.State)
	computePoolSchema["min_nodes"] = computePool.MinNodes
	computePoolSchema["max_nodes"] = computePool.MaxNodes
	computePoolSchema["instance_family"] = string(computePool.InstanceFamily)
	computePoolSchema["num_services"] = computePool.NumServices
	computePoolSchema["num_jobs"] = computePool.NumJobs
	computePoolSchema["auto_suspend_secs"] = computePool.AutoSuspendSecs
	computePoolSchema["auto_resume"] = computePool.AutoResume
	computePoolSchema["active_nodes"] = computePool.ActiveNodes
	computePoolSchema["idle_nodes"] = computePool.IdleNodes
	computePoolSchema["target_nodes"] = computePool.TargetNodes
	computePoolSchema["created_on"] = computePool.CreatedOn.String()
	computePoolSchema["resumed_on"] = computePool.ResumedOn.String()
	computePoolSchema["updated_on"] = computePool.UpdatedOn.String()
	computePoolSchema["owner"] = computePool.Owner
	if computePool.Comment != nil {
		computePoolSchema["comment"] = computePool.Comment
	}
	computePoolSchema["is_exclusive"] = computePool.IsExclusive
	if computePool.Application != nil {
		computePoolSchema["application"] = computePool.Application.Name()
	}
	computePoolSchema["error_code"] = computePool.ErrorCode
	computePoolSchema["status_message"] = computePool.StatusMessage
	return computePoolSchema
}
