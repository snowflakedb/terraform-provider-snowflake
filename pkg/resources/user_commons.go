package resources

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// defaultWorkloadIdentitySchema defines the schema for the default_workload_identity block
// which is applicable to service_user and legacy_service_user resources.
var defaultWorkloadIdentitySchema = map[string]*schema.Schema{
	"default_workload_identity": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Configures the default workload identity for the user. This is used for workload identity federation to allow third-party services to authenticate as this user. Only applicable for service users and legacy service users.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"aws": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "AWS workload identity configuration.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"arn": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The ARN of the AWS IAM role to use for workload identity federation.",
							},
						},
					},
					ExactlyOneOf: []string{
						"default_workload_identity.0.aws",
						"default_workload_identity.0.gcp",
						"default_workload_identity.0.azure",
						"default_workload_identity.0.oidc",
					},
				},
				"gcp": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "GCP workload identity configuration.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"subject": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The GCP service account subject identifier.",
							},
						},
					},
					ExactlyOneOf: []string{
						"default_workload_identity.0.aws",
						"default_workload_identity.0.gcp",
						"default_workload_identity.0.azure",
						"default_workload_identity.0.oidc",
					},
				},
				"azure": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Azure workload identity configuration.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"issuer": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The Azure issuer URL.",
							},
							"subject": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The Azure subject identifier.",
							},
						},
					},
					ExactlyOneOf: []string{
						"default_workload_identity.0.aws",
						"default_workload_identity.0.gcp",
						"default_workload_identity.0.azure",
						"default_workload_identity.0.oidc",
					},
				},
				"oidc": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Generic OIDC workload identity configuration.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"issuer": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The OIDC issuer URL.",
							},
							"subject": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The OIDC subject identifier.",
							},
							"oidc_audience_list": {
								Type:        schema.TypeList,
								Required:    true,
								MinItems:    1,
								Description: "List of allowed OIDC audiences.",
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
					ExactlyOneOf: []string{
						"default_workload_identity.0.aws",
						"default_workload_identity.0.gcp",
						"default_workload_identity.0.azure",
						"default_workload_identity.0.oidc",
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
	// Add WIF schema to service user and legacy service user schemas
	serviceUserSchema = collections.MergeMaps(serviceUserSchema, defaultWorkloadIdentitySchema)
	legacyServiceUserSchema = collections.MergeMaps(legacyServiceUserSchema, defaultWorkloadIdentitySchema)

	for _, attr := range userExternalChangesAttributes {
		if !slices.Contains(serviceUserNotApplicableAttributes, attr) {
			serviceUserExternalChangesAttributes = append(serviceUserExternalChangesAttributes, attr)
		}
		if !slices.Contains(legacyServiceUserNotApplicableAttributes, attr) {
			legacyServiceUserExternalChangesAttributes = append(legacyServiceUserExternalChangesAttributes, attr)
		}
	}
}
