package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// UserDescribeSchema represents output of DESCRIBE query for the single UserDetails.
var UserDescribeSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"display_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"login_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"first_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"middle_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"last_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"email": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"must_change_password": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"disabled": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"snowflake_lock": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"snowflake_support": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"days_to_expiry": {
		Type:     schema.TypeFloat,
		Computed: true,
	},
	"mins_to_unlock": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"default_warehouse": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default_namespace": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default_role": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default_secondary_roles": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"ext_authn_duo": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"ext_authn_uid": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"mins_to_bypass_mfa": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"mins_to_bypass_network_policy": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"rsa_public_key": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"rsa_public_key_fp": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"rsa_public_key2": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"rsa_public_key2_fp": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"password_last_set_time": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"custom_landing_page_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"custom_landing_page_url_flush_next_ui_load": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"has_mfa": {
		Type:     schema.TypeBool,
		Computed: true,
	},
}

var _ = UserDescribeSchema

func UserDescriptionToSchema(userDetails sdk.UserDetails) []map[string]any {
	userDetailsSchema := make(map[string]any)
	if userDetails.Name != nil {
		userDetailsSchema["name"] = userDetails.Name.Value
	}
	if userDetails.Comment != nil {
		userDetailsSchema["comment"] = userDetails.Comment.Value
	}
	if userDetails.DisplayName != nil {
		userDetailsSchema["display_name"] = userDetails.DisplayName.Value
	}
	if userDetails.Type != nil {
		userDetailsSchema["type"] = userDetails.Type.Value
	}
	if userDetails.LoginName != nil {
		userDetailsSchema["login_name"] = userDetails.LoginName.Value
	}
	if userDetails.FirstName != nil {
		userDetailsSchema["first_name"] = userDetails.FirstName.Value
	}
	if userDetails.MiddleName != nil {
		userDetailsSchema["middle_name"] = userDetails.MiddleName.Value
	}
	if userDetails.LastName != nil {
		userDetailsSchema["last_name"] = userDetails.LastName.Value
	}
	if userDetails.Email != nil {
		userDetailsSchema["email"] = userDetails.Email.Value
	}
	if userDetails.MustChangePassword != nil {
		userDetailsSchema["must_change_password"] = userDetails.MustChangePassword.Value
	}
	if userDetails.Disabled != nil {
		userDetailsSchema["disabled"] = userDetails.Disabled.Value
	}
	if userDetails.SnowflakeLock != nil {
		userDetailsSchema["snowflake_lock"] = userDetails.SnowflakeLock.Value
	}
	if userDetails.SnowflakeSupport != nil {
		userDetailsSchema["snowflake_support"] = userDetails.SnowflakeSupport.Value
	}
	if userDetails.DaysToExpiry != nil && userDetails.DaysToExpiry.Value != nil {
		userDetailsSchema["days_to_expiry"] = *userDetails.DaysToExpiry.Value
	}
	if userDetails.MinsToUnlock != nil && userDetails.MinsToUnlock.Value != nil {
		userDetailsSchema["mins_to_unlock"] = *userDetails.MinsToUnlock.Value
	}
	if userDetails.DefaultWarehouse != nil {
		userDetailsSchema["default_warehouse"] = userDetails.DefaultWarehouse.Value
	}
	if userDetails.DefaultNamespace != nil {
		userDetailsSchema["default_namespace"] = userDetails.DefaultNamespace.Value
	}
	if userDetails.DefaultRole != nil {
		userDetailsSchema["default_role"] = userDetails.DefaultRole.Value
	}
	if userDetails.DefaultSecondaryRoles != nil {
		userDetailsSchema["default_secondary_roles"] = userDetails.DefaultSecondaryRoles.Value
	}
	if userDetails.ExtAuthnDuo != nil {
		userDetailsSchema["ext_authn_duo"] = userDetails.ExtAuthnDuo.Value
	}
	if userDetails.ExtAuthnUid != nil {
		userDetailsSchema["ext_authn_uid"] = userDetails.ExtAuthnUid.Value
	}
	if userDetails.MinsToBypassMfa != nil && userDetails.MinsToBypassMfa.Value != nil {
		userDetailsSchema["mins_to_bypass_mfa"] = userDetails.MinsToBypassMfa.Value
	}
	if userDetails.MinsToBypassNetworkPolicy != nil && userDetails.MinsToBypassNetworkPolicy.Value != nil {
		userDetailsSchema["mins_to_bypass_network_policy"] = userDetails.MinsToBypassNetworkPolicy.Value
	}
	if userDetails.RsaPublicKey != nil {
		userDetailsSchema["rsa_public_key"] = userDetails.RsaPublicKey.Value
	}
	if userDetails.RsaPublicKeyFp != nil {
		userDetailsSchema["rsa_public_key_fp"] = userDetails.RsaPublicKeyFp.Value
	}
	if userDetails.RsaPublicKey2 != nil {
		userDetailsSchema["rsa_public_key2"] = userDetails.RsaPublicKey2.Value
	}
	if userDetails.RsaPublicKey2Fp != nil {
		userDetailsSchema["rsa_public_key2_fp"] = userDetails.RsaPublicKey2Fp.Value
	}
	if userDetails.PasswordLastSetTime != nil {
		userDetailsSchema["password_last_set_time"] = userDetails.PasswordLastSetTime.Value
	}
	if userDetails.CustomLandingPageUrl != nil {
		userDetailsSchema["custom_landing_page_url"] = userDetails.CustomLandingPageUrl.Value
	}
	if userDetails.CustomLandingPageUrlFlushNextUiLoad != nil {
		userDetailsSchema["custom_landing_page_url_flush_next_ui_load"] = userDetails.CustomLandingPageUrlFlushNextUiLoad.Value
	}
	if userDetails.HasMfa != nil {
		userDetailsSchema["has_mfa"] = userDetails.HasMfa.Value
	}
	return []map[string]any{
		userDetailsSchema,
	}
}

var _ = UserDescriptionToSchema
