package resources

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var serviceUserNotApplicableAttributes = []string{
	"password",
	"first_name",
	"middle_name",
	"last_name",
	"must_change_password",
	"mins_to_bypass_mfa",
	"disable_mfa",
}

var legacyServiceUserNotApplicableAttributes = []string{
	"first_name",
	"middle_name",
	"last_name",
	"mins_to_bypass_mfa",
	"disable_mfa",
}

var userExternalChangesAttributes = []string{
	"password",
	"login_name",
	"display_name",
	"first_name",
	"last_name",
	"email",
	"must_change_password",
	"disabled",
	"days_to_expiry",
	"mins_to_unlock",
	"default_warehouse",
	"default_namespace",
	"default_role",
	"default_secondary_roles_option",
	"mins_to_bypass_mfa",
	"rsa_public_key",
	"rsa_public_key_2",
	"comment",
	"disable_mfa",
	"workload_identity",
}

var (
	serviceUserSchema       = make(map[string]*schema.Schema)
	legacyServiceUserSchema = make(map[string]*schema.Schema)

	serviceUserExternalChangesAttributes       = make([]string, 0)
	legacyServiceUserExternalChangesAttributes = make([]string, 0)
)

func init() {
	for k, v := range userSchema {
		if !slices.Contains(serviceUserNotApplicableAttributes, k) {
			serviceUserSchema[k] = v
		}
		if !slices.Contains(legacyServiceUserNotApplicableAttributes, k) {
			legacyServiceUserSchema[k] = v
		}
	}
	for _, attr := range userExternalChangesAttributes {
		if !slices.Contains(serviceUserNotApplicableAttributes, attr) {
			serviceUserExternalChangesAttributes = append(serviceUserExternalChangesAttributes, attr)
		}
		if !slices.Contains(legacyServiceUserNotApplicableAttributes, attr) {
			legacyServiceUserExternalChangesAttributes = append(legacyServiceUserExternalChangesAttributes, attr)
		}
	}
}
