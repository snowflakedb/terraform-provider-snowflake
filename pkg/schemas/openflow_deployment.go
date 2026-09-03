package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeOpenflowDeploymentSchema represents output of DESCRIBE query for the single OpenflowDeployment.
// DESCRIBE returns the same columns as SHOW minus created_on and updated_on.
var DescribeOpenflowDeploymentSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"status": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"vpc_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"display_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"use_private_link": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"use_user_auth_over_private_link": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"custom_ingress_hostname": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"key": {
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
}

var _ = DescribeOpenflowDeploymentSchema

func OpenflowDeploymentDetailsToSchema(deployment sdk.OpenflowDeploymentDetails) map[string]any {
	deploymentSchema := make(map[string]any)
	deploymentSchema["name"] = deployment.Name
	deploymentSchema["type"] = string(deployment.Type)
	deploymentSchema["status"] = string(deployment.Status)
	if deployment.VpcType != nil {
		deploymentSchema["vpc_type"] = string(*deployment.VpcType)
	}
	if deployment.DisplayName != nil {
		deploymentSchema["display_name"] = deployment.DisplayName
	}
	deploymentSchema["use_private_link"] = deployment.UsePrivateLink
	deploymentSchema["use_user_auth_over_private_link"] = deployment.UseUserAuthOverPrivateLink
	if deployment.CustomIngressHostname != nil {
		deploymentSchema["custom_ingress_hostname"] = deployment.CustomIngressHostname
	}
	if deployment.Key != nil {
		deploymentSchema["key"] = deployment.Key
	}
	deploymentSchema["owner"] = deployment.Owner
	if deployment.Comment != nil {
		deploymentSchema["comment"] = deployment.Comment
	}
	return deploymentSchema
}
