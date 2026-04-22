package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ShowOpenflowDeploymentSchema = map[string]*schema.Schema{
	"name": {Type: schema.TypeString, Computed: true},
	"status": {Type: schema.TypeString, Computed: true},
	"deployment_type": {Type: schema.TypeString, Computed: true},
	"vpc_type": {Type: schema.TypeString, Computed: true},
	"use_private_link": {Type: schema.TypeBool, Computed: true},
	"use_user_auth_over_privatelink": {Type: schema.TypeBool, Computed: true},
	"custom_ingress_hostname": {Type: schema.TypeString, Computed: true},
	"event_table": {Type: schema.TypeString, Computed: true},
	"display_name": {Type: schema.TypeString, Computed: true},
	"comment": {Type: schema.TypeString, Computed: true},
	"owner": {Type: schema.TypeString, Computed: true},
	"created_on": {Type: schema.TypeString, Computed: true},
}

var _ = ShowOpenflowDeploymentSchema

func OpenflowDeploymentToSchema(d *sdk.OpenflowDeployment) map[string]any {
	m := make(map[string]any)
	m["name"] = d.Name
	m["status"] = string(d.Status)
	m["deployment_type"] = string(d.DeploymentType)
	m["use_private_link"] = d.UsePrivateLink
	m["use_user_auth_over_privatelink"] = d.UseUserAuthOverPrivatelink
	m["owner"] = d.Owner
	m["created_on"] = d.CreatedOn.String()
	if d.VpcType != nil {
		m["vpc_type"] = string(*d.VpcType)
	}
	if d.CustomIngressHostname != nil {
		m["custom_ingress_hostname"] = *d.CustomIngressHostname
	}
	if d.EventTable != nil {
		m["event_table"] = *d.EventTable
	}
	if d.DisplayName != nil {
		m["display_name"] = *d.DisplayName
	}
	if d.Comment != nil {
		m["comment"] = *d.Comment
	}
	return m
}

var DescribeOpenflowDeploymentSchema = map[string]*schema.Schema{
	"name": {Type: schema.TypeString, Computed: true},
	"status": {Type: schema.TypeString, Computed: true},
	"deployment_type": {Type: schema.TypeString, Computed: true},
	"vpc_type": {Type: schema.TypeString, Computed: true},
	"use_private_link": {Type: schema.TypeBool, Computed: true},
	"use_user_auth_over_privatelink": {Type: schema.TypeBool, Computed: true},
	"custom_ingress_hostname": {Type: schema.TypeString, Computed: true},
	"event_table": {Type: schema.TypeString, Computed: true},
	"display_name": {Type: schema.TypeString, Computed: true},
	"comment": {Type: schema.TypeString, Computed: true},
	"owner": {Type: schema.TypeString, Computed: true},
	"created_on": {Type: schema.TypeString, Computed: true},
	"error_code": {Type: schema.TypeString, Computed: true},
	"status_message": {Type: schema.TypeString, Computed: true},
}

func OpenflowDeploymentDetailsToSchema(d sdk.OpenflowDeploymentDetails) map[string]any {
	m := make(map[string]any)
	m["name"] = d.Name
	m["status"] = string(d.Status)
	m["deployment_type"] = string(d.DeploymentType)
	m["use_private_link"] = d.UsePrivateLink
	m["use_user_auth_over_privatelink"] = d.UseUserAuthOverPrivatelink
	m["owner"] = d.Owner
	m["created_on"] = d.CreatedOn.String()
	if d.VpcType != nil {
		m["vpc_type"] = string(*d.VpcType)
	}
	if d.CustomIngressHostname != nil {
		m["custom_ingress_hostname"] = *d.CustomIngressHostname
	}
	if d.EventTable != nil {
		m["event_table"] = *d.EventTable
	}
	if d.DisplayName != nil {
		m["display_name"] = *d.DisplayName
	}
	if d.Comment != nil {
		m["comment"] = *d.Comment
	}
	if d.ErrorCode != nil {
		m["error_code"] = *d.ErrorCode
	}
	if d.StatusMessage != nil {
		m["status_message"] = *d.StatusMessage
	}
	return m
}
