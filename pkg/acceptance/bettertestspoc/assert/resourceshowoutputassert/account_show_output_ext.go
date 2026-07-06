package resourceshowoutputassert

func (a *AccountShowOutputAssert) HasAccountUrlNotEmpty() *AccountShowOutputAssert {
	a.ValuePresent("account_url")
	return a
}

func (a *AccountShowOutputAssert) HasCreatedOnNotEmpty() *AccountShowOutputAssert {
	a.ValuePresent("created_on")
	return a
}

func (a *AccountShowOutputAssert) HasAccountLocatorNotEmpty() *AccountShowOutputAssert {
	a.ValuePresent("account_locator")
	return a
}

func (a *AccountShowOutputAssert) HasAccountLocatorUrlNotEmpty() *AccountShowOutputAssert {
	a.ValuePresent("account_locator_url")
	return a
}

func (a *AccountShowOutputAssert) HasConsumptionBillingEntityNameNotEmpty() *AccountShowOutputAssert {
	a.ValuePresent("consumption_billing_entity_name")
	return a
}

func (a *AccountShowOutputAssert) HasOrganizationOldUrlEmpty() *AccountShowOutputAssert {
	a.StringValueSet("organization_old_url", "")
	return a
}

func (a *AccountShowOutputAssert) HasMarketplaceProviderBillingEntityNameNotEmpty() *AccountShowOutputAssert {
	a.ValuePresent("marketplace_provider_billing_entity_name")
	return a
}

func (a *AccountShowOutputAssert) HasAccountOldUrlSavedOnEmpty() *AccountShowOutputAssert {
	a.StringValueSet("account_old_url_saved_on", "")
	return a
}

func (a *AccountShowOutputAssert) HasAccountOldUrlLastUsedEmpty() *AccountShowOutputAssert {
	a.StringValueSet("account_old_url_last_used", "")
	return a
}

func (a *AccountShowOutputAssert) HasOrganizationOldUrlSavedOnEmpty() *AccountShowOutputAssert {
	a.StringValueSet("organization_old_url_saved_on", "")
	return a
}

func (a *AccountShowOutputAssert) HasOrganizationOldUrlLastUsedEmpty() *AccountShowOutputAssert {
	a.StringValueSet("organization_old_url_last_used", "")
	return a
}

func (a *AccountShowOutputAssert) HasDroppedOnEmpty() *AccountShowOutputAssert {
	a.StringValueSet("dropped_on", "")
	return a
}

func (a *AccountShowOutputAssert) HasScheduledDeletionTimeEmpty() *AccountShowOutputAssert {
	a.StringValueSet("scheduled_deletion_time", "")
	return a
}

func (a *AccountShowOutputAssert) HasRestoredOnEmpty() *AccountShowOutputAssert {
	a.StringValueSet("restored_on", "")
	return a
}

func (a *AccountShowOutputAssert) HasMovedToOrganizationEmpty() *AccountShowOutputAssert {
	a.StringValueSet("moved_to_organization", "")
	return a
}

func (a *AccountShowOutputAssert) HasMovedOnEmpty() *AccountShowOutputAssert {
	a.StringValueSet("moved_on", "")
	return a
}

func (a *AccountShowOutputAssert) HasOrganizationUrlExpirationOnEmpty() *AccountShowOutputAssert {
	a.StringValueSet("organization_url_expiration_on", "")
	return a
}

func (a *AccountShowOutputAssert) HasIsEventsAccountEmpty() *AccountShowOutputAssert {
	a.StringValueSet("is_events_account", "")
	return a
}

func (a *AccountShowOutputAssert) HasIsOrganizationAccountEmpty() *AccountShowOutputAssert {
	a.StringValueSet("is_organization_account", "")
	return a
}
