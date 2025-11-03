package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var authenticationPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the authentication policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the authentication policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the authentication policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"authentication_methods": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToAuthenticationMethodsOption),
		},
		Optional:         true,
		DiffSuppressFunc: NormalizeAndCompareEnumsInSet("authentication_methods", sdk.ToAuthenticationMethodsOption),
		Description:      fmt.Sprintf("A list of authentication methods that are allowed during login. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllAuthenticationMethods)),
	},
	"mfa_authentication_methods": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToMfaAuthenticationMethodsOption),
		},
		Optional:         true,
		DiffSuppressFunc: NormalizeAndCompareEnumsInSet("mfa_authentication_methods", sdk.ToMfaAuthenticationMethodsOption),
		Description:      fmt.Sprintf("A list of authentication methods that enforce multi-factor authentication (MFA) during login. Authentication methods not listed in this parameter do not prompt for multi-factor authentication. Allowed values are %s.", possibleValuesListed(sdk.AllMfaAuthenticationMethods)),
		Deprecated:       "This field is deprecated and will be removed in the future. The new field `ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION` will be added in the next versions of the provider. Read our [BCR Migration Guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/SNOWFLAKE_BCR_MIGRATION_GUIDE.md#changes-in-authentication-policies) for more migration steps and more details.",
	},
	"mfa_enrollment": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      fmt.Sprintf("Determines whether a user must enroll in multi-factor authentication. Valid values are (case-insensitive): %s. When REQUIRED is specified, Enforces users to enroll in MFA. If this value is used, then the `client_types` parameter must include `snowflake_ui`, because Snowsight is the only place users can enroll in multi-factor authentication (MFA).", possibleValuesListed(sdk.AllMfaEnrollmentOptions)),
		ValidateDiagFunc: sdkValidation(sdk.ToMfaEnrollmentOption),
		DiffSuppressFunc: SuppressIfAny(
			NormalizeAndCompare(sdk.ToMfaEnrollmentOption),
			func(_, oldRaw, newRaw string, _ *schema.ResourceData) bool {
				old, err := sdk.ToMfaEnrollmentReadOption(oldRaw)
				if err != nil {
					return false
				}
				new, err := sdk.ToMfaEnrollmentOption(newRaw)
				if err != nil {
					return false
				}
				return old == sdk.MfaEnrollmentReadRequiredSnowflakeUiPasswordOnly && new == sdk.MfaEnrollmentOptional
			}),
	},
	"client_types": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToClientTypesOption),
		},
		Optional:         true,
		DiffSuppressFunc: NormalizeAndCompareEnumsInSet("client_types", sdk.ToClientTypesOption),
		Description:      fmt.Sprintf("A list of clients that can authenticate with Snowflake. If a client tries to connect, and the client is not one of the valid `client_types`, then the login attempt fails. Valid values are (case-insensitive): %s. The `client_types` property of an authentication policy is a best effort method to block user logins based on specific clients. It should not be used as the sole control to establish a security boundary.", possibleValuesListed(sdk.AllClientTypes)),
	},
	"security_integrations": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("security_integrations"),
		Optional:         true,
		Description:      "A list of security integrations the authentication policy is associated with. This parameter has no effect when `saml` or `oauth` are not in the `authentication_methods` list. All values in the `security_integrations` list must be compatible with the values in the `authentication_methods` list. For example, if `security_integrations` contains a SAML security integration, and `authentication_methods` contains OAUTH, then you cannot create the authentication policy. To allow all security integrations use `ALL` as parameter.",
	},
	"mfa_policy": {
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enforce_mfa_on_external_authentication": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      fmt.Sprintf("Determines whether multi-factor authentication (MFA) is enforced on external authentication. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllEnforceMfaOnExternalAuthenticationOptions)),
					ValidateDiagFunc: sdkValidation(sdk.ToEnforceMfaOnExternalAuthenticationOption),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToEnforceMfaOnExternalAuthenticationOption),
					AtLeastOneOf:     []string{"mfa_policy.0.enforce_mfa_on_external_authentication", "mfa_policy.0.allowed_methods"},
				},
				"allowed_methods": {
					Type:        schema.TypeSet,
					Optional:    true,
					Description: fmt.Sprintf("Specifies the allowed methods for the MFA policy. Valid values are: %s. These values are case-sensitive due to Terraform limitations (it's a nested field). Prefer using uppercased values.", possibleValuesListed(sdk.AllMfaPolicyOptions)),
					Elem: &schema.Schema{
						Type:             schema.TypeString,
						ValidateDiagFunc: sdkValidation(sdk.ToMfaPolicyAllowedMethodsOption),
					},
					AtLeastOneOf: []string{"mfa_policy.0.enforce_mfa_on_external_authentication", "mfa_policy.0.allowed_methods"},
				},
			},
		},
		Optional:    true,
		Description: "Specifies the multi-factor authentication (MFA) methods that users can use as a second factor of authentication.",
	},
	"pat_policy": {
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default_expiry_in_days": {
					Type:             schema.TypeInt,
					Optional:         true,
					Description:      "Specifies the default expiration time (in days) for a programmatic access token.",
					AtLeastOneOf:     []string{"pat_policy.0.default_expiry_in_days", "pat_policy.0.max_expiry_in_days", "pat_policy.0.network_policy_evaluation"},
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				},
				"max_expiry_in_days": {
					Type:             schema.TypeInt,
					Optional:         true,
					Description:      "Specifies the maximum number of days that can be set for the expiration time for a programmatic access token.",
					AtLeastOneOf:     []string{"pat_policy.0.default_expiry_in_days", "pat_policy.0.max_expiry_in_days", "pat_policy.0.network_policy_evaluation"},
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				},
				"network_policy_evaluation": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: sdkValidation(sdk.ToNetworkPolicyEvaluationOption),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToNetworkPolicyEvaluationOption),
					Description:      "Specifies the network policy evaluation for the PAT.",
					AtLeastOneOf:     []string{"pat_policy.0.default_expiry_in_days", "pat_policy.0.max_expiry_in_days", "pat_policy.0.network_policy_evaluation"},
				},
			},
		},
		Optional:    true,
		Description: "Specifies the policy for programmatic access tokens.",
	},
	"workload_identity_policy": {
		Type:     schema.TypeList,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allowed_providers": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type:             schema.TypeString,
						ValidateDiagFunc: sdkValidation(sdk.ToAllowedProviderOption),
					},
					Description: fmt.Sprintf("Specifies the allowed providers for the workload identity policy. Valid values are: %s. These values are case-sensitive due to Terraform limitations (it's a nested field). Prefer using uppercased values.", possibleValuesListed(sdk.AllAllowedProviderOptions)),
				},
				"allowed_aws_accounts": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Specifies the list of AWS account IDs allowed by the authentication policy during workload identity authentication of type `AWS`.",
				},
				"allowed_azure_issuers": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Specifies the list of Azure Entra ID issuers allowed by the authentication policy during workload identity authentication of type `AZURE`.",
				},
				"allowed_oidc_issuers": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Specifies the list of OIDC issuers allowed by the authentication policy during workload identity authentication of type `OIDC`.",
				},
			},
		},
		Optional:    true,
		Description: "Specifies the policy for workload identity federation.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the authentication policy.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW AUTHENTICATION POLICIES` for the given policy.",
		Elem: &schema.Resource{
			Schema: schemas.ShowAuthenticationPolicySchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE AUTHENTICATION POLICY` for the given policy.",
		Elem: &schema.Resource{
			Schema: schemas.AuthenticationPolicyDescribeSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func AuthenticationPolicy() *schema.Resource {
	// TODO(SNOW-1818849): unassign policies before dropping
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.AuthenticationPolicies.DropSafely
		},
	)

	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.AuthenticationPolicyResource), TrackingCreateWrapper(resources.AuthenticationPolicy, CreateContextAuthenticationPolicy)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.AuthenticationPolicyResource), TrackingReadWrapper(resources.AuthenticationPolicy, ReadContextAuthenticationPolicy(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.AuthenticationPolicyResource), TrackingUpdateWrapper(resources.AuthenticationPolicy, UpdateContextAuthenticationPolicy)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.AuthenticationPolicyResource), TrackingDeleteWrapper(resources.AuthenticationPolicy, deleteFunc)),
		Description:   "Resource used to manage authentication policy objects. For more information, check [authentication policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-authentication-policy).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.AuthenticationPolicy, customdiff.All(
			// For now, the set/list fields have to be excluded.
			// TODO [SNOW-1648997]: address the above comment
			ComputedIfAnyAttributeChanged(authenticationPolicySchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(authenticationPolicySchema, DescribeOutputAttributeName, "name", "mfa_enrollment", "comment"),
			ComputedIfAnyAttributeChanged(authenticationPolicySchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Schema: authenticationPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.AuthenticationPolicy, ImportAuthenticationPolicy),
		},
		Timeouts: defaultTimeouts,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    cty.EmptyObject,
				Upgrade: v2_9_0_AuthenticationPolicyStateUpgrader,
			},
		},
	}
}

func ImportAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	if _, err := ImportName[sdk.SchemaObjectIdentifier](context.Background(), d, nil); err != nil {
		return nil, err
	}

	// needed as otherwise the resource will be incorrectly imported when a list-parameter value equals a default value
	authenticationPolicyDescriptions, err := client.AuthenticationPolicies.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	authenticationPolicyDetails := sdk.AuthenticationPolicyDetails(authenticationPolicyDescriptions)
	authenticationMethods, err := authenticationPolicyDetails.GetAuthenticationMethods()
	if err != nil {
		return nil, err
	}
	clientTypes, err := authenticationPolicyDetails.GetClientTypes()
	if err != nil {
		return nil, err
	}
	securityIntegrations, err := authenticationPolicyDetails.GetSecurityIntegrations()
	if err != nil {
		return nil, err
	}
	securityIntegrationStrings, err := collections.MapErr(securityIntegrations, func(r sdk.AccountObjectIdentifier) (string, error) { return r.Name(), nil })
	if err != nil {
		return nil, err
	}
	mfaEnrollment, err := authenticationPolicyDetails.GetMfaEnrollment()
	if err != nil {
		return nil, err
	}

	if err := errors.Join(
		d.Set("authentication_methods", authenticationMethods),
		d.Set("mfa_enrollment", mfaEnrollment),
		d.Set("client_types", clientTypes),
		d.Set("security_integrations", securityIntegrationStrings),
	); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	name := d.Get("name").(string)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	req := sdk.NewCreateAuthenticationPolicyRequest(id)

	// Set optionals
	if v, ok := d.GetOk("authentication_methods"); ok {
		authenticationMethodsRawList := expandStringList(v.(*schema.Set).List())
		authenticationMethods := make([]sdk.AuthenticationMethods, len(authenticationMethodsRawList))
		for i, v := range authenticationMethodsRawList {
			option, err := sdk.ToAuthenticationMethodsOption(v)
			if err != nil {
				return diag.FromErr(err)
			}
			authenticationMethods[i] = sdk.AuthenticationMethods{Method: option}
		}
		req.WithAuthenticationMethods(authenticationMethods)
	}

	// TODO(SNOW-2454947): Remove this once the 2025_06 is generally enabled.
	if v, ok := d.GetOk("mfa_authentication_methods"); ok {
		mfaAuthenticationMethodsRawList := expandStringList(v.(*schema.Set).List())
		mfaAuthenticationMethods := make([]sdk.MfaAuthenticationMethods, len(mfaAuthenticationMethodsRawList))
		for i, v := range mfaAuthenticationMethodsRawList {
			option, err := sdk.ToMfaAuthenticationMethodsOption(v)
			if err != nil {
				return diag.FromErr(err)
			}
			mfaAuthenticationMethods[i] = sdk.MfaAuthenticationMethods{Method: option}
		}
		req.WithMfaAuthenticationMethods(mfaAuthenticationMethods)
	}

	if v, ok := d.GetOk("client_types"); ok {
		clientTypesRawList := expandStringList(v.(*schema.Set).List())
		clientTypes := make([]sdk.ClientTypes, len(clientTypesRawList))
		for i, v := range clientTypesRawList {
			option, err := sdk.ToClientTypesOption(v)
			if err != nil {
				return diag.FromErr(err)
			}
			clientTypes[i] = sdk.ClientTypes{ClientType: option}
		}
		req.WithClientTypes(clientTypes)
	}

	if err := errors.Join(
		attributeMappedValueCreateBuilder(d, "security_integrations", req.WithSecurityIntegrations, ToSecurityIntegrationsRequest),
		attributeMappedValueCreateBuilder(d, "mfa_enrollment", req.WithMfaEnrollment, sdk.ToMfaEnrollmentOption),
		stringAttributeCreateBuilder(d, "comment", req.WithComment),
		attributeMappedValueCreateBuilder(d, "mfa_policy", req.WithMfaPolicy, ToMfaPolicyRequest),
		attributeMappedValueCreateBuilder(d, "pat_policy", req.WithPatPolicy, ToPatPolicyRequest),
		attributeMappedValueCreateBuilder(d, "workload_identity_policy", req.WithWorkloadIdentityPolicy, ToWorkloadIdentityPolicyRequest),
	); err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*provider.Context).Client
	if err := client.AuthenticationPolicies.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextAuthenticationPolicy(false)(ctx, d, meta)
}

func ToSecurityIntegrationsRequest(value any) (sdk.SecurityIntegrationsOptionRequest, error) {
	raw := expandStringList(value.(*schema.Set).List())
	securityIntegrations := make([]sdk.AccountObjectIdentifier, len(raw))
	for i, v := range raw {
		securityIntegrations[i] = sdk.NewAccountObjectIdentifier(v)
	}
	// ALL is handled as a single-elem list - a pseudoidentifier.
	return *sdk.NewSecurityIntegrationsOptionRequest().WithSecurityIntegrations(securityIntegrations), nil
}

func ToMfaPolicyRequest(value any) (sdk.AuthenticationPolicyMfaPolicyRequest, error) {
	mfaPolicyRaw, ok := value.([]any)
	if !ok || len(mfaPolicyRaw) != 1 {
		return sdk.AuthenticationPolicyMfaPolicyRequest{}, fmt.Errorf("unable to extract mfa policy, input is either nil or non expected type (%T): %v", value, value)
	}
	mfaConfig := mfaPolicyRaw[0].(map[string]any)
	mfaPolicy := sdk.NewAuthenticationPolicyMfaPolicyRequest()
	if v := mfaConfig["enforce_mfa_on_external_authentication"]; v.(string) != "" {
		enforceMfaOnExternalAuthentication, err := sdk.ToEnforceMfaOnExternalAuthenticationOption(v.(string))
		if err != nil {
			return sdk.AuthenticationPolicyMfaPolicyRequest{}, fmt.Errorf("unable to extract enforce MFA on external authentication, input is either nil or non expected type (%T): %v", v, v)
		}
		mfaPolicy.WithEnforceMfaOnExternalAuthentication(enforceMfaOnExternalAuthentication)
	}
	if v := mfaConfig["allowed_methods"]; v != nil {
		allowedMethods := v.(*schema.Set).List()
		values, err := collections.MapErr(allowedMethods, func(v any) (sdk.AuthenticationPolicyMfaPolicyListItem, error) {
			enum, err := sdk.ToMfaPolicyAllowedMethodsOption(v.(string))
			if err != nil {
				return sdk.AuthenticationPolicyMfaPolicyListItem{}, err
			}
			return sdk.AuthenticationPolicyMfaPolicyListItem{Method: enum}, nil
		})
		if err != nil {
			return sdk.AuthenticationPolicyMfaPolicyRequest{}, err
		}
		mfaPolicy.WithAllowedMethods(values)
	}
	return *mfaPolicy, nil
}

func ToPatPolicyRequest(value any) (sdk.AuthenticationPolicyPatPolicyRequest, error) {
	patPolicyRaw, ok := value.([]any)
	if !ok || len(patPolicyRaw) != 1 {
		return sdk.AuthenticationPolicyPatPolicyRequest{}, fmt.Errorf("unable to extract pat policy, input is either nil or non expected type (%T): %v", value, value)
	}
	patConfig := patPolicyRaw[0].(map[string]any)
	patPolicy := sdk.NewAuthenticationPolicyPatPolicyRequest()
	if v := patConfig["default_expiry_in_days"]; v.(int) != 0 {
		patPolicy.WithDefaultExpiryInDays(v.(int))
	}
	if v := patConfig["max_expiry_in_days"]; v.(int) != 0 {
		patPolicy.WithMaxExpiryInDays(v.(int))
	}
	if v := patConfig["network_policy_evaluation"]; v.(string) != "" {
		networkPolicyEvaluation, err := sdk.ToNetworkPolicyEvaluationOption(v.(string))
		if err != nil {
			return sdk.AuthenticationPolicyPatPolicyRequest{}, err
		}
		patPolicy.WithNetworkPolicyEvaluation(networkPolicyEvaluation)
	}

	return *patPolicy, nil
}

func ToWorkloadIdentityPolicyRequest(value any) (sdk.AuthenticationPolicyWorkloadIdentityPolicyRequest, error) {
	workloadIdentityPolicyRaw, ok := value.([]any)
	if !ok || len(workloadIdentityPolicyRaw) != 1 {
		return sdk.AuthenticationPolicyWorkloadIdentityPolicyRequest{}, fmt.Errorf("unable to extract workload identity policy, input is either nil or non expected type (%T): %v", value, value)
	}
	workloadIdentityPolicyConfig := workloadIdentityPolicyRaw[0].(map[string]any)
	workloadIdentityPolicy := sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest()
	if v := workloadIdentityPolicyConfig["allowed_providers"]; v != nil {
		allowedProviders := v.(*schema.Set).List()
		values, err := collections.MapErr(allowedProviders, func(v any) (sdk.AuthenticationPolicyAllowedProviderListItem, error) {
			enum, err := sdk.ToAllowedProviderOption(v.(string))
			if err != nil {
				return sdk.AuthenticationPolicyAllowedProviderListItem{}, err
			}
			return sdk.AuthenticationPolicyAllowedProviderListItem{Provider: enum}, nil
		})
		if err != nil {
			return sdk.AuthenticationPolicyWorkloadIdentityPolicyRequest{}, err
		}
		workloadIdentityPolicy.WithAllowedProviders(values)
	}
	if v := workloadIdentityPolicyConfig["allowed_aws_accounts"]; v != nil {
		allowedAwsAccounts := v.(*schema.Set).List()
		values, err := collections.MapErr(allowedAwsAccounts, func(v any) (sdk.StringListItemWrapper, error) {
			return sdk.StringListItemWrapper{Value: v.(string)}, nil
		})
		if err != nil {
			return sdk.AuthenticationPolicyWorkloadIdentityPolicyRequest{}, err
		}
		workloadIdentityPolicy.WithAllowedAwsAccounts(values)
	}
	if v := workloadIdentityPolicyConfig["allowed_azure_issuers"]; v != nil {
		allowedAzureIssuers := v.(*schema.Set).List()
		values, err := collections.MapErr(allowedAzureIssuers, func(v any) (sdk.StringListItemWrapper, error) {
			return sdk.StringListItemWrapper{Value: v.(string)}, nil
		})
		if err != nil {
			return sdk.AuthenticationPolicyWorkloadIdentityPolicyRequest{}, err
		}
		workloadIdentityPolicy.WithAllowedAzureIssuers(values)
	}
	if v := workloadIdentityPolicyConfig["allowed_oidc_issuers"]; v != nil {
		allowedOidcIssuers := v.(*schema.Set).List()
		values, err := collections.MapErr(allowedOidcIssuers, func(v any) (sdk.StringListItemWrapper, error) {
			return sdk.StringListItemWrapper{Value: v.(string)}, nil
		})
		if err != nil {
			return sdk.AuthenticationPolicyWorkloadIdentityPolicyRequest{}, err
		}
		workloadIdentityPolicy.WithAllowedOidcIssuers(values)
	}
	return *workloadIdentityPolicy, nil
}

func ReadContextAuthenticationPolicy(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		diags := diag.Diagnostics{}
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		authenticationPolicy, err := client.AuthenticationPolicies.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to retrieve authentication policy. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Authentication policy id: %s, Err: %s", d.Id(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		authenticationPolicyDescriptionsRaw, err := client.AuthenticationPolicies.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}
		authenticationPolicyDescriptions := sdk.AuthenticationPolicyDetails(authenticationPolicyDescriptionsRaw)
		if withExternalChangesMarking {
			authenticationMethods, err := authenticationPolicyDescriptions.GetAuthenticationMethods()
			if err != nil {
				return diag.FromErr(err)
			}
			mfaEnrollment, err := authenticationPolicyDescriptions.GetMfaEnrollment()
			if err != nil {
				return diag.FromErr(err)
			}
			clientTypes, err := authenticationPolicyDescriptions.GetClientTypes()
			if err != nil {
				return diag.FromErr(err)
			}
			securityIntegrations, err := authenticationPolicyDescriptions.GetSecurityIntegrations()
			if err != nil {
				return diag.FromErr(err)
			}
			mfaAuthenticationMethods, err := authenticationPolicyDescriptions.GetMfaAuthenticationMethods()
			if err != nil {
				return diag.FromErr(err)
			}
			var securityIntegrationsStrings []string
			if securityIntegrations != nil {
				securityIntegrationsStrings = make([]string, len(securityIntegrations))
				for i, v := range securityIntegrations {
					securityIntegrationsStrings[i] = v.Name()
				}
			}
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"authentication_methods", "authentication_methods", authenticationPolicyDescriptions.Raw("AUTHENTICATION_METHODS"), authenticationMethods, nil},
				outputMapping{"mfa_enrollment", "mfa_enrollment", authenticationPolicyDescriptions.Raw("MFA_ENROLLMENT"), mfaEnrollment, nil},
				outputMapping{"client_types", "client_types", authenticationPolicyDescriptions.Raw("CLIENT_TYPES"), clientTypes, nil},
				outputMapping{"security_integrations", "security_integrations", authenticationPolicyDescriptions.Raw("SECURITY_INTEGRATIONS"), securityIntegrationsStrings, nil},
				outputMapping{"mfa_authentication_methods", "mfa_authentication_methods", authenticationPolicyDescriptions.Raw("MFA_AUTHENTICATION_METHODS"), mfaAuthenticationMethods, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, authenticationPolicySchema, []string{
			"authentication_methods",
			"mfa_enrollment",
			"client_types",
			"security_integrations",
			"mfa_authentication_methods",
		}); err != nil {
			return diag.FromErr(err)
		}

		if err := errors.Join(
			d.Set("comment", authenticationPolicy.Comment),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.AuthenticationPolicyToSchema(authenticationPolicy)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.AuthenticationPolicyDescriptionToSchema(authenticationPolicyDescriptions)}),
		); err != nil {
			return diag.FromErr(err)
		}

		return diags
	}
}

func UpdateContextAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewAuthenticationPolicySetRequest(), sdk.NewAuthenticationPolicyUnsetRequest()

	// change to name
	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		err = client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(id).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	// change to authentication methods
	if d.HasChange("authentication_methods") {
		if v, ok := d.GetOk("authentication_methods"); ok {
			authenticationMethods := expandStringList(v.(*schema.Set).List())
			authenticationMethodsValues := make([]sdk.AuthenticationMethods, len(authenticationMethods))
			for i, v := range authenticationMethods {
				option, err := sdk.ToAuthenticationMethodsOption(v)
				if err != nil {
					return diag.FromErr(err)
				}
				authenticationMethodsValues[i] = sdk.AuthenticationMethods{Method: option}
			}

			set.WithAuthenticationMethods(authenticationMethodsValues)
		} else {
			unset.WithAuthenticationMethods(true)
		}
	}

	// TODO(SNOW-2454947): Remove this once the 2025_06 is generally enabled.
	// change to mfa authentication methods
	if d.HasChange("mfa_authentication_methods") {
		if v, ok := d.GetOk("mfa_authentication_methods"); ok {
			mfaAuthenticationMethods := expandStringList(v.(*schema.Set).List())
			mfaAuthenticationMethodsValues := make([]sdk.MfaAuthenticationMethods, len(mfaAuthenticationMethods))
			for i, v := range mfaAuthenticationMethods {
				option, err := sdk.ToMfaAuthenticationMethodsOption(v)
				if err != nil {
					return diag.FromErr(err)
				}
				mfaAuthenticationMethodsValues[i] = sdk.MfaAuthenticationMethods{Method: option}
			}

			set.WithMfaAuthenticationMethods(mfaAuthenticationMethodsValues)
		} else {
			unset.WithMfaAuthenticationMethods(true)
		}
	}

	// change to mfa enrollment
	if d.HasChange("mfa_enrollment") {
		if mfaEnrollmentOption, err := sdk.ToMfaEnrollmentOption(d.Get("mfa_enrollment").(string)); err == nil {
			set.WithMfaEnrollment(mfaEnrollmentOption)
		} else {
			unset.WithMfaEnrollment(true)
		}
	}

	// change to client types
	if d.HasChange("client_types") {
		if v, ok := d.GetOk("client_types"); ok {
			clientTypes := expandStringList(v.(*schema.Set).List())
			clientTypesValues := make([]sdk.ClientTypes, len(clientTypes))
			for i, v := range clientTypes {
				option, err := sdk.ToClientTypesOption(v)
				if err != nil {
					return diag.FromErr(err)
				}
				clientTypesValues[i] = sdk.ClientTypes{ClientType: option}
			}

			set.WithClientTypes(clientTypesValues)
		} else {
			unset.WithClientTypes(true)
		}
	}

	if err := errors.Join(
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		attributeMappedValueUpdate(d, "security_integrations", &set.SecurityIntegrations, &unset.SecurityIntegrations, ToSecurityIntegrationsRequest),
		ToMfaPolicyRequestUpdate(d, &set.MfaPolicy, &unset.MfaPolicy),
		ToPatPolicyRequestUpdate(d, &set.PatPolicy, &unset.PatPolicy),
		ToWorkloadIdentityPolicyRequestUpdate(d, &set.WorkloadIdentityPolicy, &unset.WorkloadIdentityPolicy),
	); err != nil {
		return diag.FromErr(err)
	}

	if !reflect.DeepEqual(*set, *sdk.NewAuthenticationPolicySetRequest()) {
		req := sdk.NewAlterAuthenticationPolicyRequest(id).WithSet(*set)
		if err := client.AuthenticationPolicies.Alter(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	if !reflect.DeepEqual(*unset, *sdk.NewAuthenticationPolicyUnsetRequest()) {
		req := sdk.NewAlterAuthenticationPolicyRequest(id).WithUnset(*unset)
		if err := client.AuthenticationPolicies.Alter(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextAuthenticationPolicy(false)(ctx, d, meta)
}

func ToMfaPolicyRequestUpdate(d *schema.ResourceData, set **sdk.AuthenticationPolicyMfaPolicyRequest, unset **bool) error {
	if !d.HasChange("mfa_policy") {
		return nil
	}
	_, mfaConfigRaw := d.GetChange("mfa_policy")
	mfaConfigList := mfaConfigRaw.([]any)
	if len(mfaConfigList) == 0 {
		*unset = sdk.Bool(true)
		return nil
	}
	mfaConfig := mfaConfigList[0].(map[string]any)
	req := sdk.NewAuthenticationPolicyMfaPolicyRequest()
	if d.HasChange("mfa_policy.0.allowed_methods") {
		var allowedMethods []sdk.AuthenticationPolicyMfaPolicyListItem
		if v, ok := mfaConfig["allowed_methods"]; ok && len(v.(*schema.Set).List()) > 0 {
			allowedMethodsRaw := v.(*schema.Set).List()
			values, err := collections.MapErr(allowedMethodsRaw, func(v any) (sdk.AuthenticationPolicyMfaPolicyListItem, error) {
				enum, err := sdk.ToMfaPolicyAllowedMethodsOption(v.(string))
				if err != nil {
					return sdk.AuthenticationPolicyMfaPolicyListItem{}, err
				}
				return sdk.AuthenticationPolicyMfaPolicyListItem{Method: enum}, nil
			})
			if err != nil {
				return err
			}
			allowedMethods = values
		} else {
			allowedMethods = []sdk.AuthenticationPolicyMfaPolicyListItem{
				{Method: sdk.MfaPolicyAllowedMethodAll},
			}
		}
		req.WithAllowedMethods(allowedMethods)
	}
	if d.HasChange("mfa_policy.0.enforce_mfa_on_external_authentication") {
		var enforceMfaOnExternalAuthentication sdk.EnforceMfaOnExternalAuthenticationOption
		if v, ok := mfaConfig["enforce_mfa_on_external_authentication"]; ok {
			v, err := sdk.ToEnforceMfaOnExternalAuthenticationOption(v.(string))
			if err != nil {
				return err
			}
			enforceMfaOnExternalAuthentication = v
		} else {
			enforceMfaOnExternalAuthentication = sdk.EnforceMfaOnExternalAuthenticationNone
		}
		req.WithEnforceMfaOnExternalAuthentication(enforceMfaOnExternalAuthentication)
	}
	if !reflect.DeepEqual(*req, *sdk.NewAuthenticationPolicyMfaPolicyRequest()) {
		*set = req
	}

	return nil
}

func ToPatPolicyRequestUpdate(d *schema.ResourceData, set **sdk.AuthenticationPolicyPatPolicyRequest, unset **bool) error {
	if !d.HasChange("pat_policy") {
		return nil
	}
	_, patPolicyRaw := d.GetChange("pat_policy")
	patConfigList := patPolicyRaw.([]any)
	if len(patConfigList) == 0 {
		*unset = sdk.Bool(true)
		return nil
	}
	patConfig := patConfigList[0].(map[string]any)
	req := sdk.NewAuthenticationPolicyPatPolicyRequest()
	if d.HasChange("pat_policy.0.default_expiry_in_days") {
		if v, ok := patConfig["default_expiry_in_days"]; ok {
			req.WithDefaultExpiryInDays(v.(int))
		} else {
			req.WithDefaultExpiryInDays(15)
		}
	}
	if d.HasChange("pat_policy.0.max_expiry_in_days") {
		if v, ok := patConfig["max_expiry_in_days"]; ok {
			req.WithMaxExpiryInDays(v.(int))
		} else {
			req.WithMaxExpiryInDays(365)
		}
	}
	if d.HasChange("pat_policy.0.network_policy_evaluation") {
		if v, ok := patConfig["network_policy_evaluation"]; ok {
			networkPolicyEvaluation, err := sdk.ToNetworkPolicyEvaluationOption(v.(string))
			if err != nil {
				return err
			}
			req.WithNetworkPolicyEvaluation(networkPolicyEvaluation)
		} else {
			req.WithNetworkPolicyEvaluation(sdk.NetworkPolicyEvaluationEnforcedRequired)
		}
	}
	if !reflect.DeepEqual(*req, *sdk.NewAuthenticationPolicyPatPolicyRequest()) {
		*set = req
	}
	return nil
}

func ToWorkloadIdentityPolicyRequestUpdate(d *schema.ResourceData, set **sdk.AuthenticationPolicyWorkloadIdentityPolicyRequest, unset **bool) error {
	if !d.HasChange("workload_identity_policy") {
		return nil
	}
	_, workloadIdentityPolicyRaw := d.GetChange("workload_identity_policy")
	workloadIdentityPolicyConfigList := workloadIdentityPolicyRaw.([]any)
	if len(workloadIdentityPolicyConfigList) == 0 {
		*unset = sdk.Bool(true)
		return nil
	}
	workloadIdentityPolicyConfig := workloadIdentityPolicyConfigList[0].(map[string]any)
	req := sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest()
	if d.HasChange("workload_identity_policy.0.allowed_providers") {
		var allowedProviders []sdk.AuthenticationPolicyAllowedProviderListItem
		if v, ok := workloadIdentityPolicyConfig["allowed_providers"]; ok && len(v.(*schema.Set).List()) > 0 {
			allowedProvidersRaw := v.(*schema.Set).List()
			values, err := collections.MapErr(allowedProvidersRaw, func(v any) (sdk.AuthenticationPolicyAllowedProviderListItem, error) {
				enum, err := sdk.ToAllowedProviderOption(v.(string))
				if err != nil {
					return sdk.AuthenticationPolicyAllowedProviderListItem{}, err
				}
				return sdk.AuthenticationPolicyAllowedProviderListItem{Provider: enum}, nil
			})
			if err != nil {
				return err
			}
			allowedProviders = values
		} else {
			allowedProviders = []sdk.AuthenticationPolicyAllowedProviderListItem{
				{Provider: sdk.AllowedProviderAll},
			}
		}
		req.WithAllowedProviders(allowedProviders)
	}
	if d.HasChange("workload_identity_policy.0.allowed_aws_accounts") {
		var allowedAwsAccounts []sdk.StringListItemWrapper
		if v, ok := workloadIdentityPolicyConfig["allowed_aws_accounts"]; ok && len(v.(*schema.Set).List()) > 0 {
			allowedAwsAccountsRaw := v.(*schema.Set).List()
			values, err := collections.MapErr(allowedAwsAccountsRaw, func(v any) (sdk.StringListItemWrapper, error) {
				return sdk.StringListItemWrapper{Value: v.(string)}, nil
			})
			if err != nil {
				return err
			}
			allowedAwsAccounts = values
		} else {
			allowedAwsAccounts = []sdk.StringListItemWrapper{
				{Value: "ALL"},
			}
		}
		req.WithAllowedAwsAccounts(allowedAwsAccounts)
	}
	if d.HasChange("workload_identity_policy.0.allowed_azure_issuers") {
		var allowedAzureIssuers []sdk.StringListItemWrapper
		if v, ok := workloadIdentityPolicyConfig["allowed_azure_issuers"]; ok && len(v.(*schema.Set).List()) > 0 {
			allowedAzureIssuersRaw := v.(*schema.Set).List()
			values, err := collections.MapErr(allowedAzureIssuersRaw, func(v any) (sdk.StringListItemWrapper, error) {
				return sdk.StringListItemWrapper{Value: v.(string)}, nil
			})
			if err != nil {
				return err
			}
			allowedAzureIssuers = values
		} else {
			allowedAzureIssuers = []sdk.StringListItemWrapper{
				{Value: "ALL"},
			}
		}
		req.WithAllowedAzureIssuers(allowedAzureIssuers)
	}
	if d.HasChange("workload_identity_policy.0.allowed_oidc_issuers") {
		var allowedOidcIssuers []sdk.StringListItemWrapper
		if v, ok := workloadIdentityPolicyConfig["allowed_oidc_issuers"]; ok && len(v.(*schema.Set).List()) > 0 {
			allowedOidcIssuersRaw := v.(*schema.Set).List()
			values, err := collections.MapErr(allowedOidcIssuersRaw, func(v any) (sdk.StringListItemWrapper, error) {
				return sdk.StringListItemWrapper{Value: v.(string)}, nil
			})
			if err != nil {
				return err
			}
			allowedOidcIssuers = values
		} else {
			allowedOidcIssuers = []sdk.StringListItemWrapper{
				{Value: "ALL"},
			}
		}
		req.WithAllowedOidcIssuers(allowedOidcIssuers)
	}
	if !reflect.DeepEqual(*req, *sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest()) {
		*set = req
	}
	return nil
}
