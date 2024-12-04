// Code generated by sdk-to-schema generator; DO NOT EDIT.

package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ShowAccountSchema represents output of SHOW query for the single Account.
var ShowAccountSchema = map[string]*schema.Schema{
	"organization_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"account_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"snowflake_region": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"region_group": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"edition": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"account_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"account_locator": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"account_locator_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"managed_accounts": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"consumption_billing_entity_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"marketplace_consumer_billing_entity_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"marketplace_provider_billing_entity_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"old_account_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_org_admin": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"account_old_url_saved_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"account_old_url_last_used": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"organization_old_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"organization_old_url_saved_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"organization_old_url_last_used": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"is_events_account": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"is_organization_account": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"dropped_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"scheduled_deletion_time": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"restored_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"moved_to_organization": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"moved_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"organization_url_expiration_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

var _ = ShowAccountSchema

func AccountToSchema(account *sdk.Account) map[string]any {
	accountSchema := make(map[string]any)
	accountSchema["organization_name"] = account.OrganizationName
	accountSchema["account_name"] = account.AccountName
	accountSchema["snowflake_region"] = account.SnowflakeRegion
	// TODO: Check if populated or have to deref
	if account.RegionGroup != nil {
		accountSchema["region_group"] = account.RegionGroup
	}
	if account.Edition != nil {
		// Manually modified, please don't re-generate
		accountSchema["edition"] = string(*account.Edition)
	}
	if account.AccountURL != nil {
		accountSchema["account_url"] = account.AccountURL
	}
	if account.CreatedOn != nil {
		accountSchema["created_on"] = account.CreatedOn.String()
	}
	if account.Comment != nil {
		accountSchema["comment"] = account.Comment
	}
	accountSchema["account_locator"] = account.AccountLocator
	if account.AccountLocatorURL != nil {
		accountSchema["account_locator_url"] = account.AccountLocatorURL
	}
	if account.ManagedAccounts != nil {
		accountSchema["managed_accounts"] = account.ManagedAccounts
	}
	if account.ConsumptionBillingEntityName != nil {
		accountSchema["consumption_billing_entity_name"] = account.ConsumptionBillingEntityName
	}
	if account.MarketplaceConsumerBillingEntityName != nil {
		accountSchema["marketplace_consumer_billing_entity_name"] = account.MarketplaceConsumerBillingEntityName
	}
	if account.MarketplaceProviderBillingEntityName != nil {
		accountSchema["marketplace_provider_billing_entity_name"] = account.MarketplaceProviderBillingEntityName
	}
	if account.OldAccountURL != nil {
		accountSchema["old_account_url"] = account.OldAccountURL
	}
	if account.IsOrgAdmin != nil {
		accountSchema["is_org_admin"] = account.IsOrgAdmin
	}
	if account.AccountOldUrlSavedOn != nil {
		accountSchema["account_old_url_saved_on"] = account.AccountOldUrlSavedOn.String()
	}
	if account.AccountOldUrlLastUsed != nil {
		accountSchema["account_old_url_last_used"] = account.AccountOldUrlLastUsed.String()
	}
	if account.OrganizationOldUrl != nil {
		accountSchema["organization_old_url"] = account.OrganizationOldUrl
	}
	if account.OrganizationOldUrlSavedOn != nil {
		accountSchema["organization_old_url_saved_on"] = account.OrganizationOldUrlSavedOn.String()
	}
	if account.OrganizationOldUrlLastUsed != nil {
		accountSchema["organization_old_url_last_used"] = account.OrganizationOldUrlLastUsed.String()
	}
	if account.IsEventsAccount != nil {
		accountSchema["is_events_account"] = account.IsEventsAccount
	}
	accountSchema["is_organization_account"] = account.IsOrganizationAccount
	if account.DroppedOn != nil {
		accountSchema["dropped_on"] = account.DroppedOn.String()
	}
	if account.ScheduledDeletionTime != nil {
		accountSchema["scheduled_deletion_time"] = account.ScheduledDeletionTime.String()
	}
	if account.RestoredOn != nil {
		accountSchema["restored_on"] = account.RestoredOn.String()
	}
	if account.MovedToOrganization != nil {
		accountSchema["moved_to_organization"] = account.MovedToOrganization
	}
	if account.MovedOn != nil {
		accountSchema["moved_on"] = account.MovedOn
	}
	if account.OrganizationUrlExpirationOn != nil {
		accountSchema["organization_url_expiration_on"] = account.OrganizationUrlExpirationOn.String()
	}
	return accountSchema
}

var _ = AccountToSchema
