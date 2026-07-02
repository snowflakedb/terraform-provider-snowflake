package resourceshowoutputassert

func (u *UserShowOutputAssert) HasCreatedOnNotEmpty() *UserShowOutputAssert {
	u.ValuePresent("created_on")
	return u
}

func (u *UserShowOutputAssert) HasDaysToExpiryNotEmpty() *UserShowOutputAssert {
	u.ValuePresent("days_to_expiry")
	return u
}

func (u *UserShowOutputAssert) HasMinsToUnlockNotEmpty() *UserShowOutputAssert {
	u.ValuePresent("mins_to_unlock")
	return u
}

func (u *UserShowOutputAssert) HasMinsToBypassMfaNotEmpty() *UserShowOutputAssert {
	u.ValuePresent("mins_to_bypass_mfa")
	return u
}

func (u *UserShowOutputAssert) HasTypeEmpty() *UserShowOutputAssert {
	u.StringValueSet("type", "")
	return u
}
