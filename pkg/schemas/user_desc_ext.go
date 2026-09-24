package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [next PRs]: every UserDetails field is *XxxProperty (Value/DefaultValue/Description),
// which the generator cannot map, so this ext owns the whole schema. P2 should choose SDK
// flattening to scalars vs native XxxProperty unwrap, then un-skip and delete this file.

func (userDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":         {Type: schema.TypeString, Computed: true},
		"comment":      {Type: schema.TypeString, Computed: true},
		"display_name": {Type: schema.TypeString, Computed: true},
		"type":         {Type: schema.TypeString, Computed: true},
		"login_name":   {Type: schema.TypeString, Computed: true},
		"first_name":   {Type: schema.TypeString, Computed: true},
		"middle_name":  {Type: schema.TypeString, Computed: true},
		"last_name":    {Type: schema.TypeString, Computed: true},
		"email":        {Type: schema.TypeString, Computed: true},
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
		"default_warehouse":       {Type: schema.TypeString, Computed: true},
		"default_namespace":       {Type: schema.TypeString, Computed: true},
		"default_role":            {Type: schema.TypeString, Computed: true},
		"default_secondary_roles": {Type: schema.TypeString, Computed: true},
		"ext_authn_duo": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"ext_authn_uid": {Type: schema.TypeString, Computed: true},
		"mins_to_bypass_mfa": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"mins_to_bypass_network_policy": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"rsa_public_key":     {Type: schema.TypeString, Computed: true},
		"rsa_public_key_fp":  {Type: schema.TypeString, Computed: true},
		"rsa_public_key2":    {Type: schema.TypeString, Computed: true},
		"rsa_public_key2_fp": {Type: schema.TypeString, Computed: true},
		"password_last_set_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"custom_landing_page_url": {Type: schema.TypeString, Computed: true},
		"custom_landing_page_url_flush_next_ui_load": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"has_mfa": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"has_workload_identity": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func (userDetailsToSchemaMapper) additionalToSchema(src *sdk.UserDetails, dst map[string]any) {
	if src.Name != nil {
		dst["name"] = src.Name.Value
	}
	if src.Comment != nil {
		dst["comment"] = src.Comment.Value
	}
	if src.DisplayName != nil {
		dst["display_name"] = src.DisplayName.Value
	}
	if src.Type != nil {
		dst["type"] = src.Type.Value
	}
	if src.LoginName != nil {
		dst["login_name"] = src.LoginName.Value
	}
	if src.FirstName != nil {
		dst["first_name"] = src.FirstName.Value
	}
	if src.MiddleName != nil {
		dst["middle_name"] = src.MiddleName.Value
	}
	if src.LastName != nil {
		dst["last_name"] = src.LastName.Value
	}
	if src.Email != nil {
		dst["email"] = src.Email.Value
	}
	if src.MustChangePassword != nil {
		dst["must_change_password"] = src.MustChangePassword.Value
	}
	if src.Disabled != nil {
		dst["disabled"] = src.Disabled.Value
	}
	if src.SnowflakeLock != nil {
		dst["snowflake_lock"] = src.SnowflakeLock.Value
	}
	if src.SnowflakeSupport != nil {
		dst["snowflake_support"] = src.SnowflakeSupport.Value
	}
	if src.DaysToExpiry != nil && src.DaysToExpiry.Value != nil {
		dst["days_to_expiry"] = *src.DaysToExpiry.Value
	}
	if src.MinsToUnlock != nil && src.MinsToUnlock.Value != nil {
		dst["mins_to_unlock"] = *src.MinsToUnlock.Value
	}
	if src.DefaultWarehouse != nil {
		dst["default_warehouse"] = src.DefaultWarehouse.Value
	}
	if src.DefaultNamespace != nil {
		dst["default_namespace"] = src.DefaultNamespace.Value
	}
	if src.DefaultRole != nil {
		dst["default_role"] = src.DefaultRole.Value
	}
	if src.DefaultSecondaryRoles != nil {
		dst["default_secondary_roles"] = src.DefaultSecondaryRoles.Value
	}
	if src.ExtAuthnDuo != nil {
		dst["ext_authn_duo"] = src.ExtAuthnDuo.Value
	}
	if src.ExtAuthnUid != nil {
		dst["ext_authn_uid"] = src.ExtAuthnUid.Value
	}
	if src.MinsToBypassMfa != nil && src.MinsToBypassMfa.Value != nil {
		dst["mins_to_bypass_mfa"] = *src.MinsToBypassMfa.Value
	}
	if src.MinsToBypassNetworkPolicy != nil && src.MinsToBypassNetworkPolicy.Value != nil {
		dst["mins_to_bypass_network_policy"] = *src.MinsToBypassNetworkPolicy.Value
	}
	if src.RsaPublicKey != nil {
		dst["rsa_public_key"] = src.RsaPublicKey.Value
	}
	if src.RsaPublicKeyFp != nil {
		dst["rsa_public_key_fp"] = src.RsaPublicKeyFp.Value
	}
	if src.RsaPublicKey2 != nil {
		dst["rsa_public_key2"] = src.RsaPublicKey2.Value
	}
	if src.RsaPublicKey2Fp != nil {
		dst["rsa_public_key2_fp"] = src.RsaPublicKey2Fp.Value
	}
	if src.PasswordLastSetTime != nil {
		dst["password_last_set_time"] = src.PasswordLastSetTime.Value
	}
	if src.CustomLandingPageUrl != nil {
		dst["custom_landing_page_url"] = src.CustomLandingPageUrl.Value
	}
	if src.CustomLandingPageUrlFlushNextUiLoad != nil {
		dst["custom_landing_page_url_flush_next_ui_load"] = src.CustomLandingPageUrlFlushNextUiLoad.Value
	}
	if src.HasMfa != nil {
		dst["has_mfa"] = src.HasMfa.Value
	}
	if src.HasWorkloadIdentity != nil {
		dst["has_workload_identity"] = src.HasWorkloadIdentity.Value
	}
}
