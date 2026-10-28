package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	authenticationPolicyCreateExt      = TrackingCreateWrapper(resources.AuthenticationPolicy, CreateAuthenticationPolicy)
	authenticationPolicyReadExt        = TrackingReadWrapper(resources.AuthenticationPolicy, ReadAuthenticationPolicyFunc(true))
	authenticationPolicyUpdateExt      = TrackingUpdateWrapper(resources.AuthenticationPolicy, UpdateAuthenticationPolicy)
	authenticationPolicyDescriptionExt = "Resource used to manage authentication policy objects. For more information, check [authentication policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-authentication-policy)."
	authenticationPolicyImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.AuthenticationPolicy, ImportAuthenticationPolicy),
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

type authenticationPolicyAlterRequestsExt struct {
	set   *sdk.AuthenticationPolicySetRequest
	unset *sdk.AuthenticationPolicyUnsetRequest
}

func authenticationPolicyParseIdExt(id string) (sdk.SchemaObjectIdentifier, error) {
	return sdk.ParseSchemaObjectIdentifier(id)
}

func authenticationPolicyNewCreateRequestExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier) (*sdk.CreateAuthenticationPolicyRequest, error) {
	return sdk.NewCreateAuthenticationPolicyRequest(id), nil
}

func authenticationPolicyHandleExternalChangesInDescribeExt(d *schema.ResourceData, details *sdk.AuthenticationPolicyDetails, withExternalChangesMarking bool) diag.Diagnostics {
	if details == nil {
		return nil
	}
	if err := d.Set("client_policy", authenticationPolicyClientPolicyToSchema(details.ClientPolicy)); err != nil {
		return diag.FromErr(err)
	}
	if withExternalChangesMarking {
		described := details.Describe()
		if described == nil {
			return nil
		}
		// valueToCompare matches the TypeString describe_output values. valueToSet is the config representation.
		if err := handleExternalChangesToObjectInFlatDescribe(d,
			outputMapping{"authentication_methods", "authentication_methods", described.AuthenticationMethods, authenticationPolicyEnumStrings(details.AuthenticationMethods), nil},
			outputMapping{"mfa_enrollment", "mfa_enrollment", described.MfaEnrollment, string(details.MfaEnrollment), nil},
			outputMapping{"client_types", "client_types", described.ClientTypes, authenticationPolicyEnumStrings(details.ClientTypes), nil},
			outputMapping{"security_integrations", "security_integrations", described.SecurityIntegrations, authenticationPolicySecurityIntegrationNames(details.SecurityIntegrations), nil},
			outputMapping{"mfa_policy", "mfa_policy", described.MfaPolicy, authenticationPolicyMfaPolicyToSchema(details.MfaPolicy), nil},
			outputMapping{"pat_policy", "pat_policy", described.PatPolicy, authenticationPolicyPatPolicyToSchema(details.PatPolicy), nil},
			outputMapping{"workload_identity_policy", "workload_identity_policy", described.WorkloadIdentityPolicy, authenticationPolicyWorkloadIdentityPolicyToSchema(details.WorkloadIdentityPolicy), nil},
		); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func authenticationPolicyEnumStrings[T ~string](values []T) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = string(value)
	}
	return out
}

func authenticationPolicyHierarchyRenameExt(ctx context.Context, d *schema.ResourceData, meta any, id *sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		client := meta.(*provider.Context).Client
		err := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(*id).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		*id = newId
	}
	return nil
}

func authenticationPolicyApplySetUnsetExt(d *schema.ResourceData) (*authenticationPolicyAlterRequestsExt, error) {
	set, unset := sdk.NewAuthenticationPolicySetRequest(), sdk.NewAuthenticationPolicyUnsetRequest()

	if err := errors.Join(
		setValueUpdate(d, "authentication_methods", &set.AuthenticationMethods, &unset.AuthenticationMethods, toAuthenticationPolicyAuthenticationMethods),
		attributeMappedValueUpdate(d, "mfa_enrollment", &set.MfaEnrollment, &unset.MfaEnrollment, sdk.ToMfaEnrollmentOption),
		setValueUpdate(d, "client_types", &set.ClientTypes, &unset.ClientTypes, toAuthenticationPolicyClientTypes),
		setValueUpdate(d, "client_policy", &set.ClientPolicy, &unset.ClientPolicy, toAuthenticationPolicyClientPolicy),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		attributeMappedValueUpdate(d, "security_integrations", &set.SecurityIntegrations, &unset.SecurityIntegrations, ToSecurityIntegrationsRequest),
		singleBlockAttributeUpdate(d, "mfa_policy", &set.MfaPolicy, &unset.MfaPolicy, applyMfaPolicyUpdate),
		singleBlockAttributeUpdate(d, "pat_policy", &set.PatPolicy, &unset.PatPolicy, applyPatPolicyUpdate),
		singleBlockAttributeUpdate(d, "workload_identity_policy", &set.WorkloadIdentityPolicy, &unset.WorkloadIdentityPolicy, applyWorkloadIdentityPolicyUpdate),
	); err != nil {
		return nil, err
	}
	return &authenticationPolicyAlterRequestsExt{set: set, unset: unset}, nil
}

func authenticationPolicySetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, details *sdk.AuthenticationPolicyDetails, importPermissiveMode bool) error {
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		importSetAttribute(d, "authentication_methods", importPermissiveMode, details.AuthenticationMethods),
		importStringAttribute(d, "mfa_enrollment", importPermissiveMode, string(details.MfaEnrollment)),
		importSetAttribute(d, "client_types", importPermissiveMode, details.ClientTypes),
		importSetAttribute(d, "security_integrations", importPermissiveMode, authenticationPolicySecurityIntegrationNames(details.SecurityIntegrations)),
		importBlockAttribute(d, "mfa_policy", importPermissiveMode, authenticationPolicyMfaPolicyToSchema(details.MfaPolicy)),
		importBlockAttribute(d, "pat_policy", importPermissiveMode, authenticationPolicyPatPolicyToSchema(details.PatPolicy)),
		importBlockAttribute(d, "workload_identity_policy", importPermissiveMode, authenticationPolicyWorkloadIdentityPolicyToSchema(details.WorkloadIdentityPolicy)),
	)
}

func authenticationPolicySecurityIntegrationNames(details sdk.AuthenticationPolicySecurityIntegrationsDetails) []string {
	if details.All {
		return []string{"ALL"}
	}
	names := make([]string, len(details.SecurityIntegrations))
	for i, id := range details.SecurityIntegrations {
		names[i] = id.Name()
	}
	return names
}

// =============================================================================
// SDK incision points
// =============================================================================

func authenticationPolicyCreateInSdkExt(ctx context.Context, meta any, _ sdk.SchemaObjectIdentifier, req *sdk.CreateAuthenticationPolicyRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.AuthenticationPolicies.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// Import does not SHOW. This returns the same DESCRIBE details Read already uses, in one query.
func authenticationPolicyShowByIdInSdkExt(ctx context.Context, meta any, id sdk.SchemaObjectIdentifier) (*sdk.AuthenticationPolicyDetails, error) {
	return authenticationPolicyDescribeInSdkExt(ctx, meta, id)
}

func authenticationPolicyDescribeInSdkExt(ctx context.Context, meta any, id sdk.SchemaObjectIdentifier) (*sdk.AuthenticationPolicyDetails, error) {
	client := meta.(*provider.Context).Client
	return client.AuthenticationPolicies.DescribeDetails(ctx, id)
}

func authenticationPolicyAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier, reqs *authenticationPolicyAlterRequestsExt) error {
	set := reqs.set
	if !reflect.DeepEqual(*set, *sdk.NewAuthenticationPolicySetRequest()) {
		client := meta.(*provider.Context).Client
		req := sdk.NewAlterAuthenticationPolicyRequest(id).WithSet(*set)
		if err := client.AuthenticationPolicies.Alter(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

func authenticationPolicyAlterUnsetInSdkExt(ctx context.Context, meta any, id sdk.SchemaObjectIdentifier, reqs *authenticationPolicyAlterRequestsExt) error {
	unset := reqs.unset
	if !reflect.DeepEqual(*unset, *sdk.NewAuthenticationPolicyUnsetRequest()) {
		client := meta.(*provider.Context).Client
		req := sdk.NewAlterAuthenticationPolicyRequest(id).WithUnset(*unset)
		if err := client.AuthenticationPolicies.Alter(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

func toAuthenticationPolicyAuthenticationMethods(v any) (sdk.AuthenticationMethods, error) {
	option, err := sdk.ToAuthenticationMethodsOption(v.(string))
	if err != nil {
		return sdk.AuthenticationMethods{}, err
	}
	return sdk.AuthenticationMethods{Method: option}, nil
}

func toAuthenticationPolicyClientTypes(v any) (sdk.ClientTypes, error) {
	option, err := sdk.ToClientTypesOption(v.(string))
	if err != nil {
		return sdk.ClientTypes{}, err
	}
	return sdk.ClientTypes{ClientType: option}, nil
}

func applyMfaPolicyUpdate(d *schema.ResourceData, prefix string, req *sdk.AuthenticationPolicyMfaPolicyRequest) error {
	return errors.Join(
		attributeMappedValueUpdateSetOnlyFallback(d, prefix+"enforce_mfa_on_external_authentication", &req.EnforceMfaOnExternalAuthentication, sdk.ToEnforceMfaOnExternalAuthenticationOption, sdk.EnforceMfaOnExternalAuthenticationOptionNone),
		setValueUpdateFallback(d, prefix+"allowed_methods", &req.AllowedMethods, toAuthenticationPolicyMfaPolicyAllowedMethods, []sdk.AuthenticationPolicyMfaPolicyListItem{{
			Method: sdk.MfaPolicyAllowedMethodsOptionAll,
		}}),
	)
}

func toAuthenticationPolicyMfaPolicyAllowedMethods(v any) (sdk.AuthenticationPolicyMfaPolicyListItem, error) {
	option, err := sdk.ToMfaPolicyAllowedMethodsOption(v.(string))
	if err != nil {
		return sdk.AuthenticationPolicyMfaPolicyListItem{}, err
	}
	return sdk.AuthenticationPolicyMfaPolicyListItem{Method: option}, nil
}

func applyPatPolicyUpdate(d *schema.ResourceData, prefix string, req *sdk.AuthenticationPolicyPatPolicyRequest) error {
	return errors.Join(
		intAttributeUnsetFallbackUpdateWithZeroDefault(d, prefix+"default_expiry_in_days", &req.DefaultExpiryInDays, 15),
		intAttributeUnsetFallbackUpdateWithZeroDefault(d, prefix+"max_expiry_in_days", &req.MaxExpiryInDays, 365),
		attributeMappedValueUpdateSetOnlyFallback(d, prefix+"network_policy_evaluation", &req.NetworkPolicyEvaluation, sdk.ToNetworkPolicyEvaluationOption, sdk.NetworkPolicyEvaluationOptionEnforcedRequired),
		booleanStringAttributeUnsetFallbackUpdate(d, prefix+"require_role_restriction_for_service_users", &req.RequireRoleRestrictionForServiceUsers, true),
	)
}

func applyWorkloadIdentityPolicyUpdate(d *schema.ResourceData, prefix string, req *sdk.AuthenticationPolicyWorkloadIdentityPolicyRequest) error {
	return errors.Join(
		setValueUpdateFallback(d, prefix+"allowed_providers", &req.AllowedProviders, toAuthenticationPolicyWorkloadIdentityPolicyAllowedProviders, []sdk.AuthenticationPolicyAllowedProviderListItem{{
			Provider: sdk.AllowedProviderOptionAll,
		}}),
		setValueUpdateFallback(d, prefix+"allowed_aws_accounts", &req.AllowedAwsAccounts, toAuthenticationPolicyWorkloadIdentityPolicyAllowedAwsAccounts, allStringListItems()),
		setValueUpdateFallback(d, prefix+"allowed_azure_issuers", &req.AllowedAzureIssuers, toAuthenticationPolicyWorkloadIdentityPolicyAllowedAzureIssuers, allStringListItems()),
		setValueUpdateFallback(d, prefix+"allowed_oidc_issuers", &req.AllowedOidcIssuers, toAuthenticationPolicyWorkloadIdentityPolicyAllowedOidcIssuers, allStringListItems()),
	)
}

func allStringListItems() []sdk.StringListItemWrapper {
	return []sdk.StringListItemWrapper{{Value: "ALL"}}
}

func toAuthenticationPolicyWorkloadIdentityPolicyAllowedProviders(v any) (sdk.AuthenticationPolicyAllowedProviderListItem, error) {
	option, err := sdk.ToAllowedProviderOption(v.(string))
	if err != nil {
		return sdk.AuthenticationPolicyAllowedProviderListItem{}, err
	}
	return sdk.AuthenticationPolicyAllowedProviderListItem{Provider: option}, nil
}

func stringListItem(v any) (sdk.StringListItemWrapper, error) {
	return sdk.StringListItemWrapper{Value: v.(string)}, nil
}

var (
	toAuthenticationPolicyWorkloadIdentityPolicyAllowedAwsAccounts  = stringListItem
	toAuthenticationPolicyWorkloadIdentityPolicyAllowedAzureIssuers = stringListItem
	toAuthenticationPolicyWorkloadIdentityPolicyAllowedOidcIssuers  = stringListItem
)

func ToSecurityIntegrationsRequest(value any) (sdk.SecurityIntegrationsOptionRequest, error) {
	raw := expandStringList(value.(*schema.Set).List())
	securityIntegrations := make([]sdk.AccountObjectIdentifier, len(raw))
	for i, v := range raw {
		securityIntegrations[i] = sdk.NewAccountObjectIdentifier(v)
	}
	// ALL is handled as a single-elem list - a pseudoidentifier.
	return *sdk.NewSecurityIntegrationsOptionRequest().WithSecurityIntegrations(securityIntegrations), nil
}

func toAuthenticationPolicyClientPolicy(v any) (sdk.AuthenticationPolicyClientPolicyEntry, error) {
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

func authenticationPolicyClientPolicyToSchema(entries []sdk.AuthenticationPolicyClientPolicyEntryDetails) []map[string]any {
	out := make([]map[string]any, len(entries))
	for i, entry := range entries {
		out[i] = map[string]any{
			"client_type":     string(entry.ClientType),
			"minimum_version": entry.MinimumVersion,
		}
	}
	return out
}

func authenticationPolicyMfaPolicyToSchema(details sdk.AuthenticationPolicyMfaPolicyDetails) []map[string]any {
	return []map[string]any{{
		"allowed_methods":                        authenticationPolicyEnumStrings(details.AllowedMethods),
		"enforce_mfa_on_external_authentication": string(details.EnforceMfaOnExternalAuthentication),
	}}
}

func authenticationPolicyPatPolicyToSchema(details sdk.AuthenticationPolicyPatPolicyDetails) []map[string]any {
	return []map[string]any{{
		"default_expiry_in_days":                     details.DefaultExpiryInDays,
		"max_expiry_in_days":                         details.MaxExpiryInDays,
		"network_policy_evaluation":                  string(details.NetworkPolicyEvaluation),
		"require_role_restriction_for_service_users": booleanStringFromBool(details.RequireRoleRestrictionForServiceUsers),
	}}
}

func authenticationPolicyWorkloadIdentityPolicyToSchema(details sdk.AuthenticationPolicyWorkloadIdentityPolicyDetails) []map[string]any {
	return []map[string]any{{
		"allowed_providers":     authenticationPolicyEnumStrings(details.AllowedProviders),
		"allowed_aws_accounts":  details.AllowedAwsAccounts,
		"allowed_azure_issuers": details.AllowedAzureIssuers,
		"allowed_oidc_issuers":  details.AllowedOidcIssuers,
	}}
}
