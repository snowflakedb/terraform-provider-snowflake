package sdk

func GetSecondaryRolesOptionFrom(text string) SecondaryRolesOption {
	if text != "" {
		parsedRoles := ParseCommaSeparatedStringArray(text, true)
		if len(parsedRoles) > 0 {
			return SecondaryRolesOptionAll
		} else {
			return SecondaryRolesOptionNone
		}
	}
	return SecondaryRolesOptionDefault
}

func (v *User) GetSecondaryRolesOption() SecondaryRolesOption {
	return GetSecondaryRolesOptionFrom(v.DefaultSecondaryRoles)
}

func userDetailsFromRows(rows []propertyRow) *UserDetails {
	v := &UserDetails{}
	for _, row := range rows {
		switch row.Property {
		case "NAME":
			v.Name = row.toStringProperty()
		case "COMMENT":
			v.Comment = row.toStringProperty()
		case "DISPLAY_NAME":
			v.DisplayName = row.toStringProperty()
		case "TYPE":
			v.Type = row.toStringProperty()
		case "LOGIN_NAME":
			v.LoginName = row.toStringProperty()
		case "FIRST_NAME":
			v.FirstName = row.toStringProperty()
		case "MIDDLE_NAME":
			v.MiddleName = row.toStringProperty()
		case "LAST_NAME":
			v.LastName = row.toStringProperty()
		case "EMAIL":
			v.Email = row.toStringProperty()
		case "PASSWORD":
			v.Password = row.toStringProperty()
		case "MUST_CHANGE_PASSWORD":
			v.MustChangePassword = row.toBoolProperty()
		case "DISABLED":
			v.Disabled = row.toBoolProperty()
		case "SNOWFLAKE_LOCK":
			v.SnowflakeLock = row.toBoolProperty()
		case "SNOWFLAKE_SUPPORT":
			v.SnowflakeSupport = row.toBoolProperty()
		case "DAYS_TO_EXPIRY":
			v.DaysToExpiry = row.toFloatProperty()
		case "MINS_TO_UNLOCK":
			v.MinsToUnlock = row.toIntProperty()
		case "DEFAULT_WAREHOUSE":
			v.DefaultWarehouse = row.toStringProperty()
		case "DEFAULT_NAMESPACE":
			v.DefaultNamespace = row.toStringProperty()
		case "DEFAULT_ROLE":
			v.DefaultRole = row.toStringProperty()
		case "DEFAULT_SECONDARY_ROLES":
			v.DefaultSecondaryRoles = row.toStringProperty()
		case "EXT_AUTHN_DUO":
			v.ExtAuthnDuo = row.toBoolProperty()
		case "EXT_AUTHN_UID":
			v.ExtAuthnUid = row.toStringProperty()
		case "HAS_MFA":
			v.HasMfa = row.toBoolProperty()
		case "MINS_TO_BYPASS_MFA":
			v.MinsToBypassMfa = row.toIntProperty()
		case "MINS_TO_BYPASS_NETWORK_POLICY":
			v.MinsToBypassNetworkPolicy = row.toIntProperty()
		case "RSA_PUBLIC_KEY":
			v.RsaPublicKey = row.toStringProperty()
		case "RSA_PUBLIC_KEY_FP":
			v.RsaPublicKeyFp = row.toStringProperty()
		case "RSA_PUBLIC_KEY_LAST_SET_TIME":
			v.RsaPublicKeyLastSetTime = row.toStringProperty()
		case "RSA_PUBLIC_KEY_2":
			v.RsaPublicKey2 = row.toStringProperty()
		case "RSA_PUBLIC_KEY_2_FP":
			v.RsaPublicKey2Fp = row.toStringProperty()
		case "RSA_PUBLIC_KEY_2_LAST_SET_TIME":
			v.RsaPublicKey2LastSetTime = row.toStringProperty()
		case "PASSWORD_LAST_SET_TIME":
			v.PasswordLastSetTime = row.toStringProperty()
		case "CUSTOM_LANDING_PAGE_URL":
			v.CustomLandingPageUrl = row.toStringProperty()
		case "CUSTOM_LANDING_PAGE_URL_FLUSH_NEXT_UI_LOAD":
			v.CustomLandingPageUrlFlushNextUiLoad = row.toBoolProperty()
		case "HAS_WORKLOAD_IDENTITY":
			v.HasWorkloadIdentity = row.toBoolProperty()
		}
	}
	return v
}
