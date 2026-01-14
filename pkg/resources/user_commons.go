package resources

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var defaultWorkloadIdentitySchema = map[string]*schema.Schema{
	"default_workload_identity": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Specifies the default workload identity that the user will use to authenticate. For more information, see [Workload identity federation](https://docs.snowflake.com/en/user-guide/workload-identity-federation).",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"aws": {
					Type:         schema.TypeList,
					Optional:     true,
					MaxItems:     1,
					ExactlyOneOf: []string{"default_workload_identity.0.aws", "default_workload_identity.0.azure", "default_workload_identity.0.gcp", "default_workload_identity.0.oidc"},
					Description:  "AWS workload identity configuration.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"arn": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the Amazon Resource Name (ARN) for the AWS IAM user or role that will be used for authentication.",
							},
						},
					},
				},
				"azure": {
					Type:         schema.TypeList,
					Optional:     true,
					MaxItems:     1,
					ExactlyOneOf: []string{"default_workload_identity.0.aws", "default_workload_identity.0.azure", "default_workload_identity.0.gcp", "default_workload_identity.0.oidc"},
					Description:  "Azure workload identity configuration.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"issuer": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the issuer URL.",
							},
							"subject": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the subject.",
							},
						},
					},
				},
				"gcp": {
					Type:         schema.TypeList,
					Optional:     true,
					MaxItems:     1,
					ExactlyOneOf: []string{"default_workload_identity.0.aws", "default_workload_identity.0.azure", "default_workload_identity.0.gcp", "default_workload_identity.0.oidc"},
					Description:  "GCP workload identity configuration.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"subject": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the subject.",
							},
						},
					},
				},
				"oidc": {
					Type:         schema.TypeList,
					Optional:     true,
					MaxItems:     1,
					ExactlyOneOf: []string{"default_workload_identity.0.aws", "default_workload_identity.0.azure", "default_workload_identity.0.gcp", "default_workload_identity.0.oidc"},
					Description:  "OIDC workload identity configuration.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"issuer": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the issuer URL.",
							},
							"subject": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Specifies the subject.",
							},
							"oidc_audience_list": {
								Type:        schema.TypeList,
								Optional:    true,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Description: "Specifies the custom audience list for OIDC workload identity.",
							},
						},
					},
				},
			},
		},
	},
}

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
	"default_workload_identity",
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

	// Add workload identity schema only to service and legacy service users
	for k, v := range defaultWorkloadIdentitySchema {
		serviceUserSchema[k] = v
		legacyServiceUserSchema[k] = v
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
