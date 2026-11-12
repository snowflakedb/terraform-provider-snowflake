package resources

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/validators"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var failoverGroupIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the identifier for the failover group. The identifier must start with an alphabetic character and cannot contain spaces or special characters unless the identifier string is enclosed in double quotes (e.g. \"My object\"). Identifiers enclosed in double quotes are also case-sensitive.",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.EqualFold(old, new)
		},
	},
}

var failoverGroupAttributesSchemaExt = map[string]*schema.Schema{
	"object_types": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validators.NormalizeValidation(sdk.ToPluralObjectType),
		},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Type(s) of objects for which you are enabling replication and failover from the source account to the target account. The following object types are supported: \"ACCOUNT PARAMETERS\", \"DATABASES\", \"INTEGRATIONS\", \"NETWORK POLICIES\", \"RESOURCE MONITORS\", \"ROLES\", \"SHARES\", \"USERS\", \"WAREHOUSES\"",
	},
	"allowed_databases": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Specifies the database or list of databases for which you are enabling replication and failover from the source account to the target account. The OBJECT_TYPES list must include DATABASES to set this parameter.",
	},
	"allowed_shares": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Specifies the share or list of shares for which you are enabling replication and failover from the source account to the target account. The OBJECT_TYPES list must include SHARES to set this parameter.",
	},
	"allowed_integration_types": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Type(s) of integrations for which you are enabling replication and failover from the source account to the target account. This property requires that the OBJECT_TYPES list include INTEGRATIONS to set this parameter. The following integration types are supported: \"SECURITY INTEGRATIONS\", \"API INTEGRATIONS\", \"STORAGE INTEGRATIONS\", \"EXTERNAL ACCESS INTEGRATIONS\", \"NOTIFICATION INTEGRATIONS\"",
	},
	"allowed_accounts": {
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		ConflictsWith: []string{"from_replica"},
		Description:   "Specifies the target account or list of target accounts to which replication and failover of specified objects from the source account is enabled. Secondary failover groups in the target accounts in this list can be promoted to serve as the primary failover group in case of failover. Expected in the form `<org_name>.<target_account_name>`. This value is case-sensitive.",
	},
	"ignore_edition_check": {
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ConflictsWith: []string{"from_replica"},
		Description:   "Allows replicating objects to accounts on lower editions.",
	},
	"from_replica": {
		Type:          schema.TypeList,
		Optional:      true,
		ForceNew:      true,
		MaxItems:      1,
		ConflictsWith: []string{"object_types", "allowed_accounts", "allowed_databases", "allowed_shares", "allowed_integration_types", "ignore_edition_check", "replication_schedule", "optimized_refresh", "rpo_assurance"},
		Description:   "Specifies the name of the replica to use as the source for the failover group.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"organization_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of your Snowflake organization.",
				},
				"source_account_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Source account from which you are enabling replication and failover of the specified objects.",
				},
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Identifier for the primary failover group in the source account.",
				},
			},
		},
	},
	"replication_schedule": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		Description:   "Specifies the schedule for refreshing secondary failover groups.",
		ConflictsWith: []string{"from_replica", "rpo_assurance"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cron": {
					Type:          schema.TypeList,
					Optional:      true,
					MaxItems:      1,
					ConflictsWith: []string{"replication_schedule.interval"},
					Description:   "Specifies the cron expression for the replication schedule. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"expression": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the cron expression for the replication schedule. The cron expression must be in the following format: \"minute hour day-of-month month day-of-week\". The following values are supported: minute: 0-59 hour: 0-23 day-of-month: 1-31 month: 1-12 day-of-week: 0-6 (0 is Sunday)",
							},
							"time_zone": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the time zone for secondary group refresh.",
							},
						},
					},
				},
				"interval": {
					Type:          schema.TypeInt,
					Optional:      true,
					ConflictsWith: []string{"replication_schedule.cron"},
					Description:   "Specifies the interval in minutes for the replication schedule. The interval must be greater than 0 and less than 1440 (24 hours).",
				},
			},
		},
	},
	"optimized_refresh": {
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ConflictsWith: []string{"from_replica", "rpo_assurance"},
		Description:   "Specifies whether the failover group uses Optimized Refresh. Can be set on the primary failover group only. When true, a replication_schedule with an interval of 6 hours or less is required (Snowflake recommends 10 MINUTE or less). Do not set when rpo_assurance is true. Default: false.",
	},
	"rpo_assurance": {
		Type:          schema.TypeBool,
		Optional:      true,
		Default:       false,
		ConflictsWith: []string{"from_replica", "replication_schedule", "optimized_refresh"},
		Description:   "Specifies whether the failover group uses RPO Assurance. Can be set on the primary failover group only. Cannot be set with replication_schedule (unset the schedule first). Enables Optimized Refresh automatically; do not also set optimized_refresh. Default: false.",
	},
}
