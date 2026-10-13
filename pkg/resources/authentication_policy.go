package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
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
		Description:      fmt.Sprintf("A list of authentication methods that are allowed during login. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllAuthenticationMethodsOptions)),
	},
	"mfa_authentication_methods": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToMfaAuthenticationMethodsOption),
		},
		Optional:         true,
		DiffSuppressFunc: ignoreAlwaysSuppressFunc,
		Description:      fmt.Sprintf("A list of authentication methods that enforce multi-factor authentication (MFA) during login. Authentication methods not listed in this parameter do not prompt for multi-factor authentication. Allowed values are %s.", possibleValuesListed(sdk.AllMfaAuthenticationMethodsOptions)),
		Deprecated:       "This field is deprecated and will be removed in the future. Currently, it has no effect. Use the new `enforce_mfa_on_external_authentication` field instead. Read our [BCR Migration Guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/SNOWFLAKE_BCR_MIGRATION_GUIDE.md#changes-in-authentication-policies) for more migration steps and more details.",
	},
	"mfa_enrollment": {
		Type:     schema.TypeString,
		Optional: true,
		Description: fmt.Sprintf("Determines whether a user must enroll in multi-factor authentication. Valid values are (case-insensitive): %s. When REQUIRED is specified, Enforces users to enroll in MFA."+
			" If this value is used, then the `client_types` parameter must include `snowflake_ui`, because Snowsight is the only place users can enroll in multi-factor authentication (MFA)."+
			" Note that when you set this value to OPTIONAL, and your account setup forces users to enroll in MFA, then Snowflake may set quietly this value to `REQUIRED_PASSWORD_ONLY`, which may cause permadiff."+
			" In this case, you may want to adjust this field value.", possibleValuesListed(sdk.AllMfaEnrollmentOptions)),
		ValidateDiagFunc: sdkValidation(sdk.ToMfaEnrollmentOption),
		DiffSuppressFunc: SuppressIfAny(
			NormalizeAndCompare(sdk.ToMfaEnrollmentOption),
			// We need a custom diff because when this field is set to OPTIONAL, the value is returned as REQUIRED_PASSWORD_ONLY.
			func(_, oldRaw, newRaw string, _ *schema.ResourceData) bool {
				old, err := sdk.ToMfaEnrollmentReadOption(oldRaw)
				if err != nil {
					return false
				}
				new, err := sdk.ToMfaEnrollmentOption(newRaw)
				if err != nil {
					return false
				}
				return old == sdk.MfaEnrollmentReadOptionRequiredPasswordOnly && new == sdk.MfaEnrollmentOptionOptional
			},
		),
	},
	"client_types": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToClientTypesOption),
		},
		Optional:         true,
		DiffSuppressFunc: NormalizeAndCompareEnumsInSet("client_types", sdk.ToClientTypesOption),
		Description:      fmt.Sprintf("A list of clients that can authenticate with Snowflake. If a client tries to connect, and the client is not one of the valid `client_types`, then the login attempt fails. Valid values are (case-insensitive): %s. The `client_types` property of an authentication policy is a best effort method to block user logins based on specific clients. It should not be used as the sole control to establish a security boundary.", possibleValuesListed(sdk.AllClientTypesOptions)),
	},
	"client_policy": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Allows to set policies per-client type. Only valid when `client_types` is empty, contains ALL, or contains DRIVERS.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"client_type": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      fmt.Sprintf("The client or driver type. Valid values (case-insensitive): %s.", possibleValuesListed(sdk.AllClientPolicyDriverTypes)),
					ValidateDiagFunc: sdkValidation(sdk.ToClientPolicyDriverType),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToClientPolicyDriverType),
				},
				"minimum_version": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Minimum allowed version for this client/driver type (e.g. '1.14.1').",
				},
			},
		},
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
					Description: fmt.Sprintf("Specifies the allowed methods for the MFA policy. Valid values are: %s. These values are case-sensitive due to Terraform limitations (it's a nested field). Prefer using uppercased values.", possibleValuesListed(sdk.AllMfaPolicyAllowedMethodsOptions)),
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
					AtLeastOneOf:     []string{"pat_policy.0.default_expiry_in_days", "pat_policy.0.max_expiry_in_days", "pat_policy.0.require_role_restriction_for_service_users", "pat_policy.0.network_policy_evaluation"},
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				},
				"max_expiry_in_days": {
					Type:             schema.TypeInt,
					Optional:         true,
					Description:      "Specifies the maximum number of days that can be set for the expiration time for a programmatic access token.",
					AtLeastOneOf:     []string{"pat_policy.0.default_expiry_in_days", "pat_policy.0.max_expiry_in_days", "pat_policy.0.require_role_restriction_for_service_users", "pat_policy.0.network_policy_evaluation"},
					ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
				},
				"require_role_restriction_for_service_users": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      booleanStringFieldDescription("If true, when you generate a programmatic access token for a service user, you must restrict the use of that token to a specific role."),
					AtLeastOneOf:     []string{"pat_policy.0.default_expiry_in_days", "pat_policy.0.max_expiry_in_days", "pat_policy.0.require_role_restriction_for_service_users", "pat_policy.0.network_policy_evaluation"},
					Default:          BooleanDefault,
					ValidateDiagFunc: validateBooleanString,
				},
				"network_policy_evaluation": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: sdkValidation(sdk.ToNetworkPolicyEvaluationOption),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToNetworkPolicyEvaluationOption),
					Description:      fmt.Sprintf("Specifies the network policy evaluation for the PAT. Valid values are: %s.", possibleValuesListed(sdk.AllNetworkPolicyEvaluationOptions)),
					AtLeastOneOf:     []string{"pat_policy.0.default_expiry_in_days", "pat_policy.0.max_expiry_in_days", "pat_policy.0.require_role_restriction_for_service_users", "pat_policy.0.network_policy_evaluation"},
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
			Schema: schemas.DescribeAuthenticationPolicyDescribeDetailsSchema,
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

		CreateContext: TrackingCreateWrapper(resources.AuthenticationPolicy, CreateContextAuthenticationPolicy),
		ReadContext:   TrackingReadWrapper(resources.AuthenticationPolicy, ReadContextAuthenticationPolicy(true)),
		UpdateContext: TrackingUpdateWrapper(resources.AuthenticationPolicy, UpdateContextAuthenticationPolicy),
		DeleteContext: TrackingDeleteWrapper(resources.AuthenticationPolicy, deleteFunc),
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

	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
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

	if err := errors.Join(
		setMappedValueCreateBuilder(d, "authentication_methods", req.WithAuthenticationMethods, authenticationMethodItem),
		setMappedValueCreateBuilder(d, "client_types", req.WithClientTypes, clientTypeItem),
		setMappedValueCreateBuilder(d, "client_policy", req.WithClientPolicy, toClientPolicyEntries),
		attributeMappedValueCreateBuilder(d, "security_integrations", req.WithSecurityIntegrations, ToSecurityIntegrationsRequest),
		attributeMappedValueCreateBuilder(d, "mfa_enrollment", req.WithMfaEnrollment, sdk.ToMfaEnrollmentOption),
		stringAttributeCreateBuilder(d, "comment", req.WithComment),
		singleBlockAttributeCreate(d, "mfa_policy", req.WithMfaPolicy, parseMfaPolicy),
		singleBlockAttributeCreate(d, "pat_policy", req.WithPatPolicy, parsePatPolicy),
		singleBlockAttributeCreate(d, "workload_identity_policy", req.WithWorkloadIdentityPolicy, parseWorkloadIdentityPolicy),
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

func ReadContextAuthenticationPolicy(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		diags := make(diag.Diagnostics, 0)

		if withExternalChangesMarking {
			describeMappings := make([]outputMapping, 0)

			authenticationMethods, err := authenticationPolicyDescriptions.GetAuthenticationMethods()
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve authentication from authentication policy describe details.",
					Detail:   fmt.Sprintf("Authentication policy id: %s, Err: %s", d.Id(), err),
				})
			} else {
				describeMappings = append(
					describeMappings,
					outputMapping{"authentication_methods", "authentication_methods", authenticationPolicyDescriptions.Raw("AUTHENTICATION_METHODS"), authenticationMethods, nil},
				)
			}

			mfaEnrollment, err := authenticationPolicyDescriptions.GetMfaEnrollment()
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve mfa enrollment from authentication policy describe details.",
					Detail:   fmt.Sprintf("Authentication policy id: %s, Err: %s", d.Id(), err),
				})
			} else {
				describeMappings = append(
					describeMappings,
					outputMapping{"mfa_enrollment", "mfa_enrollment", authenticationPolicyDescriptions.Raw("MFA_ENROLLMENT"), mfaEnrollment, nil},
				)
			}

			clientTypes, err := authenticationPolicyDescriptions.GetClientTypes()
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve client types from authentication policy describe details.",
					Detail:   fmt.Sprintf("Authentication policy id: %s, Err: %s", d.Id(), err),
				})
			} else {
				describeMappings = append(
					describeMappings,
					outputMapping{"client_types", "client_types", authenticationPolicyDescriptions.Raw("CLIENT_TYPES"), clientTypes, nil},
				)
			}

			securityIntegrations, err := authenticationPolicyDescriptions.GetSecurityIntegrations()
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve security integrations from authentication policy describe details.",
					Detail:   fmt.Sprintf("Authentication policy id: %s, Err: %s", d.Id(), err),
				})
			} else {
				var securityIntegrationsStrings []string
				if securityIntegrations != nil {
					securityIntegrationsStrings = make([]string, len(securityIntegrations))
					for i, v := range securityIntegrations {
						securityIntegrationsStrings[i] = v.Name()
					}
				}

				describeMappings = append(
					describeMappings,
					outputMapping{"security_integrations", "security_integrations", authenticationPolicyDescriptions.Raw("SECURITY_INTEGRATIONS"), securityIntegrationsStrings, nil},
				)
			}

			if len(describeMappings) > 0 {
				if err = handleExternalChangesToObjectInFlatDescribe(d, describeMappings...); err != nil {
					return append(diags, diag.FromErr(err)...)
				}
			}
		}

		if err = setStateToValuesFromConfig(d, authenticationPolicySchema, []string{
			"authentication_methods",
			"mfa_enrollment",
			"client_types",
			"security_integrations",
		}); err != nil {
			return append(diags, diag.FromErr(err)...)
		}

		if err := errors.Join(
			d.Set("comment", authenticationPolicy.Comment),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.AuthenticationPolicyToSchema(authenticationPolicy)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.AuthenticationPolicyDescriptionsToSchema(authenticationPolicyDescriptions)}),
		); err != nil {
			return append(diags, diag.FromErr(err)...)
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

	if err := errors.Join(
		setValueUpdate(d, "authentication_methods", &set.AuthenticationMethods, &unset.AuthenticationMethods, authenticationMethodItem),
		attributeMappedValueUpdate(d, "mfa_enrollment", &set.MfaEnrollment, &unset.MfaEnrollment, sdk.ToMfaEnrollmentOption),
		setValueUpdate(d, "client_types", &set.ClientTypes, &unset.ClientTypes, clientTypeItem),
		setValueUpdate(d, "client_policy", &set.ClientPolicy, &unset.ClientPolicy, toClientPolicyEntries),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		attributeMappedValueUpdate(d, "security_integrations", &set.SecurityIntegrations, &unset.SecurityIntegrations, ToSecurityIntegrationsRequest),
		singleBlockAttributeUpdate(d, "mfa_policy", &set.MfaPolicy, &unset.MfaPolicy, applyMfaPolicyUpdate),
		singleBlockAttributeUpdate(d, "pat_policy", &set.PatPolicy, &unset.PatPolicy, applyPatPolicyUpdate),
		singleBlockAttributeUpdate(d, "workload_identity_policy", &set.WorkloadIdentityPolicy, &unset.WorkloadIdentityPolicy, applyWorkloadIdentityPolicyUpdate),
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

func authenticationMethodItem(v any) (sdk.AuthenticationMethods, error) {
	option, err := sdk.ToAuthenticationMethodsOption(v.(string))
	if err != nil {
		return sdk.AuthenticationMethods{}, err
	}
	return sdk.AuthenticationMethods{Method: option}, nil
}

func clientTypeItem(v any) (sdk.ClientTypes, error) {
	option, err := sdk.ToClientTypesOption(v.(string))
	if err != nil {
		return sdk.ClientTypes{}, err
	}
	return sdk.ClientTypes{ClientType: option}, nil
}

func parseMfaPolicy(d *schema.ResourceData, prefix string) (sdk.AuthenticationPolicyMfaPolicyRequest, error) {
	req := sdk.NewAuthenticationPolicyMfaPolicyRequest()
	err := errors.Join(
		attributeMappedValueCreateBuilder(d, prefix+"enforce_mfa_on_external_authentication", req.WithEnforceMfaOnExternalAuthentication, sdk.ToEnforceMfaOnExternalAuthenticationOption),
		setValueCreate(d, prefix+"allowed_methods", &req.AllowedMethods, mfaPolicyAllowedMethodItem),
	)
	return *req, err
}

func applyMfaPolicyUpdate(d *schema.ResourceData, prefix string, req *sdk.AuthenticationPolicyMfaPolicyRequest) error {
	return errors.Join(
		attributeMappedValueUpdateSetOnlyFallback(d, prefix+"enforce_mfa_on_external_authentication", &req.EnforceMfaOnExternalAuthentication, sdk.ToEnforceMfaOnExternalAuthenticationOption, sdk.EnforceMfaOnExternalAuthenticationOptionNone),
		setValueUpdateFallback(d, prefix+"allowed_methods", &req.AllowedMethods, mfaPolicyAllowedMethodItem, []sdk.AuthenticationPolicyMfaPolicyListItem{{
			Method: sdk.MfaPolicyAllowedMethodsOptionAll,
		}}),
	)
}

func mfaPolicyAllowedMethodItem(v any) (sdk.AuthenticationPolicyMfaPolicyListItem, error) {
	option, err := sdk.ToMfaPolicyAllowedMethodsOption(v.(string))
	if err != nil {
		return sdk.AuthenticationPolicyMfaPolicyListItem{}, err
	}
	return sdk.AuthenticationPolicyMfaPolicyListItem{Method: option}, nil
}

func parsePatPolicy(d *schema.ResourceData, prefix string) (sdk.AuthenticationPolicyPatPolicyRequest, error) {
	req := sdk.NewAuthenticationPolicyPatPolicyRequest()
	err := errors.Join(
		intAttributeCreateBuilder(d, prefix+"default_expiry_in_days", req.WithDefaultExpiryInDays),
		intAttributeCreateBuilder(d, prefix+"max_expiry_in_days", req.WithMaxExpiryInDays),
		attributeMappedValueCreateBuilder(d, prefix+"network_policy_evaluation", req.WithNetworkPolicyEvaluation, sdk.ToNetworkPolicyEvaluationOption),
		booleanStringAttributeCreateFallbackBuilder(d, prefix+"require_role_restriction_for_service_users", req.WithRequireRoleRestrictionForServiceUsers, true),
	)
	return *req, err
}

func applyPatPolicyUpdate(d *schema.ResourceData, prefix string, req *sdk.AuthenticationPolicyPatPolicyRequest) error {
	return errors.Join(
		intAttributeUnsetFallbackUpdateWithZeroDefault(d, prefix+"default_expiry_in_days", &req.DefaultExpiryInDays, 15),
		intAttributeUnsetFallbackUpdateWithZeroDefault(d, prefix+"max_expiry_in_days", &req.MaxExpiryInDays, 365),
		attributeMappedValueUpdateSetOnlyFallback(d, prefix+"network_policy_evaluation", &req.NetworkPolicyEvaluation, sdk.ToNetworkPolicyEvaluationOption, sdk.NetworkPolicyEvaluationOptionEnforcedRequired),
		booleanStringAttributeUnsetFallbackUpdate(d, prefix+"require_role_restriction_for_service_users", &req.RequireRoleRestrictionForServiceUsers, true),
	)
}

func parseWorkloadIdentityPolicy(d *schema.ResourceData, prefix string) (sdk.AuthenticationPolicyWorkloadIdentityPolicyRequest, error) {
	req := sdk.NewAuthenticationPolicyWorkloadIdentityPolicyRequest()
	err := errors.Join(
		setValueCreate(d, prefix+"allowed_providers", &req.AllowedProviders, allowedProviderItem),
		setValueCreate(d, prefix+"allowed_aws_accounts", &req.AllowedAwsAccounts, stringListItem),
		setValueCreate(d, prefix+"allowed_azure_issuers", &req.AllowedAzureIssuers, stringListItem),
		setValueCreate(d, prefix+"allowed_oidc_issuers", &req.AllowedOidcIssuers, stringListItem),
	)
	return *req, err
}

func applyWorkloadIdentityPolicyUpdate(d *schema.ResourceData, prefix string, req *sdk.AuthenticationPolicyWorkloadIdentityPolicyRequest) error {
	return errors.Join(
		setValueUpdateFallback(d, prefix+"allowed_providers", &req.AllowedProviders, allowedProviderItem, []sdk.AuthenticationPolicyAllowedProviderListItem{{
			Provider: sdk.AllowedProviderOptionAll,
		}}),
		setValueUpdateFallback(d, prefix+"allowed_aws_accounts", &req.AllowedAwsAccounts, stringListItem, allStringListItems()),
		setValueUpdateFallback(d, prefix+"allowed_azure_issuers", &req.AllowedAzureIssuers, stringListItem, allStringListItems()),
		setValueUpdateFallback(d, prefix+"allowed_oidc_issuers", &req.AllowedOidcIssuers, stringListItem, allStringListItems()),
	)
}

func allStringListItems() []sdk.StringListItemWrapper {
	return []sdk.StringListItemWrapper{{Value: "ALL"}}
}

func allowedProviderItem(v any) (sdk.AuthenticationPolicyAllowedProviderListItem, error) {
	option, err := sdk.ToAllowedProviderOption(v.(string))
	if err != nil {
		return sdk.AuthenticationPolicyAllowedProviderListItem{}, err
	}
	return sdk.AuthenticationPolicyAllowedProviderListItem{Provider: option}, nil
}

func stringListItem(v any) (sdk.StringListItemWrapper, error) {
	return sdk.StringListItemWrapper{Value: v.(string)}, nil
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

func toClientPolicyEntries(v any) (sdk.AuthenticationPolicyClientPolicyEntry, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return sdk.AuthenticationPolicyClientPolicyEntry{}, fmt.Errorf("client_policy entry: expected map, got %T", v)
	}
	clientType, err := sdk.ToClientPolicyDriverType(m["client_type"].(string))
	if err != nil {
		return sdk.AuthenticationPolicyClientPolicyEntry{}, err
	}
	return sdk.AuthenticationPolicyClientPolicyEntry{
		ClientType: clientType,
		Params: &sdk.AuthenticationPolicyClientPolicyEntryParams{
			MinimumVersion: sdk.Pointer(m["minimum_version"].(string)),
		},
	}, nil
}
