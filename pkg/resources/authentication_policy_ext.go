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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	authenticationPolicyCreateExt = TrackingCreateWrapper(resources.AuthenticationPolicy, CreateAuthenticationPolicy)
	authenticationPolicyReadExt   = TrackingReadWrapper(resources.AuthenticationPolicy, ReadAuthenticationPolicyFunc(true))
	authenticationPolicyUpdateExt = TrackingUpdateWrapper(resources.AuthenticationPolicy, UpdateAuthenticationPolicy)
	// TODO(SNOW-1818849): unassign policies before dropping
	authenticationPolicyDeleteExt = TrackingDeleteWrapper(resources.AuthenticationPolicy, ResourceDeleteContextFunc(
		authenticationPolicyParseIdExt,
		authenticationPolicyDropSafelyInSdkExt,
	))
	authenticationPolicyDescriptionExt   = "Resource used to manage authentication policy objects. For more information, check [authentication policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-authentication-policy)."
	authenticationPolicyCustomizeDiffExt = TrackingCustomDiffWrapper(resources.AuthenticationPolicy, customdiff.All(
		// For now, the set/list fields have to be excluded.
		// TODO [SNOW-1648997]: address the above comment
		ComputedIfAnyAttributeChanged(authenticationPolicySchema, ShowOutputAttributeName, "name", "comment"),
		ComputedIfAnyAttributeChanged(authenticationPolicySchema, DescribeOutputAttributeName, "name", "mfa_enrollment", "comment"),
		ComputedIfAnyAttributeChanged(authenticationPolicySchema, FullyQualifiedNameAttributeName, "name"),
	))
	authenticationPolicyImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.AuthenticationPolicy, ImportAuthenticationPolicy),
	}
	authenticationPolicySchemaVersionExt  = 1
	authenticationPolicyStateUpgradersExt = []schema.StateUpgrader{
		{
			Version: 0,
			Type:    cty.EmptyObject,
			Upgrade: v2_9_0_AuthenticationPolicyStateUpgrader,
		},
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

var authenticationPolicySetStateToValueFromConfigFieldsExt = []string{
	"authentication_methods",
	"mfa_enrollment",
	"client_types",
	"security_integrations",
}

type authenticationPolicyAlterRequestsExt struct {
	set   *sdk.AuthenticationPolicySetRequest
	unset *sdk.AuthenticationPolicyUnsetRequest
}

func authenticationPolicyParseIdFromConfigExt(d *schema.ResourceData) (sdk.SchemaObjectIdentifier, error) {
	name := d.Get("name").(string)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	return id, nil
}

func authenticationPolicyParseIdExt(id string) (sdk.SchemaObjectIdentifier, error) {
	return sdk.ParseSchemaObjectIdentifier(id)
}

func authenticationPolicyNewCreateRequestExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier) (*sdk.CreateAuthenticationPolicyRequest, error) {
	return sdk.NewCreateAuthenticationPolicyRequest(id), nil
}

func authenticationPolicyApplyCreateOptionalsExt(d *schema.ResourceData, req *sdk.CreateAuthenticationPolicyRequest) error {
	return errors.Join(
		setMappedValueCreateBuilder(d, "authentication_methods", req.WithAuthenticationMethods, authenticationMethodItem),
		setMappedValueCreateBuilder(d, "client_types", req.WithClientTypes, clientTypeItem),
		setMappedValueCreateBuilder(d, "client_policy", req.WithClientPolicy, toClientPolicyEntries),
		attributeMappedValueCreateBuilder(d, "security_integrations", req.WithSecurityIntegrations, ToSecurityIntegrationsRequest),
		attributeMappedValueCreateBuilder(d, "mfa_enrollment", req.WithMfaEnrollment, sdk.ToMfaEnrollmentOption),
		stringAttributeCreateBuilder(d, "comment", req.WithComment),
		singleBlockAttributeCreate(d, "mfa_policy", req.WithMfaPolicy, parseMfaPolicy),
		singleBlockAttributeCreate(d, "pat_policy", req.WithPatPolicy, parsePatPolicy),
		singleBlockAttributeCreate(d, "workload_identity_policy", req.WithWorkloadIdentityPolicy, parseWorkloadIdentityPolicy),
	)
}

func authenticationPolicySetConfigFieldsExt(d *schema.ResourceData, authenticationPolicy *sdk.AuthenticationPolicy) error {
	return d.Set("comment", authenticationPolicy.Comment)
}

func authenticationPolicyHandleExternalChangesInDescribeExt(d *schema.ResourceData, authenticationPolicyDescriptions sdk.AuthenticationPolicyDetails, withExternalChangesMarking bool) diag.Diagnostics {
	if !withExternalChangesMarking {
		return nil
	}
	diags := make(diag.Diagnostics, 0)
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
	return diags
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
		return nil, err
	}
	return &authenticationPolicyAlterRequestsExt{set: set, unset: unset}, nil
}

func authenticationPolicySetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, authenticationPolicyDetails sdk.AuthenticationPolicyDetails) error {
	if err := errors.Join(
		d.Set("name", id.Name()),
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
	); err != nil {
		return err
	}

	// needed as otherwise the resource will be incorrectly imported when a list-parameter value equals a default value
	authenticationMethods, err := authenticationPolicyDetails.GetAuthenticationMethods()
	if err != nil {
		return err
	}
	clientTypes, err := authenticationPolicyDetails.GetClientTypes()
	if err != nil {
		return err
	}
	securityIntegrations, err := authenticationPolicyDetails.GetSecurityIntegrations()
	if err != nil {
		return err
	}
	securityIntegrationStrings, err := collections.MapErr(securityIntegrations, func(r sdk.AccountObjectIdentifier) (string, error) { return r.Name(), nil })
	if err != nil {
		return err
	}
	mfaEnrollment, err := authenticationPolicyDetails.GetMfaEnrollment()
	if err != nil {
		return err
	}

	return errors.Join(
		d.Set("authentication_methods", authenticationMethods),
		d.Set("mfa_enrollment", mfaEnrollment),
		d.Set("client_types", clientTypes),
		d.Set("security_integrations", securityIntegrationStrings),
	)
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

func authenticationPolicyShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) (*sdk.AuthenticationPolicy, diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	authenticationPolicy, err := client.AuthenticationPolicies.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to retrieve authentication policy. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Authentication policy id: %s, Err: %s", d.Id(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return authenticationPolicy, nil
}

// Import does not SHOW. This returns the same DESCRIBE details Read already uses, in one query.
func authenticationPolicyShowByIdInSdkExt(ctx context.Context, meta any, id sdk.SchemaObjectIdentifier) (sdk.AuthenticationPolicyDetails, error) {
	return authenticationPolicyDescribeInSdkExt(ctx, meta, id)
}

func authenticationPolicyDescribeInSdkExt(ctx context.Context, meta any, id sdk.SchemaObjectIdentifier) (sdk.AuthenticationPolicyDetails, error) {
	client := meta.(*provider.Context).Client
	authenticationPolicyDescriptionsRaw, err := client.AuthenticationPolicies.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return sdk.AuthenticationPolicyDetails(authenticationPolicyDescriptionsRaw), nil
}

func authenticationPolicyAlterSetInSdkExt(ctx context.Context, meta any, id sdk.SchemaObjectIdentifier, reqs *authenticationPolicyAlterRequestsExt) error {
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

func authenticationPolicyDropSafelyInSdkExt(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
	return client.AuthenticationPolicies.DropSafely
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

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
