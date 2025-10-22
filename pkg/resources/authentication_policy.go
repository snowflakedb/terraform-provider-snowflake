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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Optional:    true,
		Description: fmt.Sprintf("A list of authentication methods that are allowed during login. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllAuthenticationMethods)),
	},
	"mfa_authentication_methods": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToMfaAuthenticationMethodsOption),
		},
		Optional:    true,
		Description: fmt.Sprintf("A list of authentication methods that enforce multi-factor authentication (MFA) during login. Authentication methods not listed in this parameter do not prompt for multi-factor authentication. Allowed values are %s.", possibleValuesListed(sdk.AllMfaAuthenticationMethods)),
		Deprecated:  "This field is deprecated and will be removed in the future. The new field `ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION` will be added in the next versions of the provider. Read our [BCR Migration Guide](https://github.com/snowflakedb/terraform-provider-snowflake/blob/main/SNOWFLAKE_BCR_MIGRATION_GUIDE.md#changes-in-authentication-policies) for more migration steps and more details.",
	},
	"mfa_enrollment": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      fmt.Sprintf("Determines whether a user must enroll in multi-factor authentication. Valid values are (case-insensitive): %s. When REQUIRED is specified, Enforces users to enroll in MFA. If this value is used, then the `client_types` parameter must include `snowflake_ui`, because Snowsight is the only place users can enroll in multi-factor authentication (MFA).", possibleValuesListed(sdk.AllMfaEnrollmentOptions)),
		ValidateDiagFunc: sdkValidation(sdk.ToMfaEnrollmentOption),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToMfaEnrollmentOption),
	},
	"client_types": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToClientTypesOption),
		},
		// TODO(next PR): add diff suppression for enum sets.
		Optional:    true,
		Description: fmt.Sprintf("A list of clients that can authenticate with Snowflake. If a client tries to connect, and the client is not one of the valid `client_types`, then the login attempt fails. Valid values are (case-insensitive): %s. The `client_types` property of an authentication policy is a best effort method to block user logins based on specific clients. It should not be used as the sole control to establish a security boundary.", possibleValuesListed(sdk.AllClientTypes)),
	},
	"security_integrations": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("security_integrations"),
		Optional:         true,
		Description:      "A list of security integrations the authentication policy is associated with. This parameter has no effect when `saml` or `oauth` are not in the `authentication_methods` list. All values in the `security_integrations` list must be compatible with the values in the `authentication_methods` list. For example, if `security_integrations` contains a SAML security integration, and `authentication_methods` contains OAUTH, then you cannot create the authentication policy. To allow all security integrations use `all` as parameter.",
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
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.AuthenticationPolicyResource), TrackingCreateWrapper(resources.AuthenticationPolicy, CreateContextAuthenticationPolicy)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.AuthenticationPolicyResource), TrackingReadWrapper(resources.AuthenticationPolicy, ReadContextAuthenticationPolicy(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.AuthenticationPolicyResource), TrackingUpdateWrapper(resources.AuthenticationPolicy, UpdateContextAuthenticationPolicy)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.AuthenticationPolicyResource), TrackingDeleteWrapper(resources.AuthenticationPolicy, deleteFunc)),
		Description:   "Resource used to manage authentication policy objects. For more information, check [authentication policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-authentication-policy).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ComputePool, customdiff.All(
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

	// TODO(XXX): Remove this once the 2025_06 is generally enabled.
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
	// TODO: handle ALL
	return *sdk.NewSecurityIntegrationsOptionRequest().WithSecurityIntegrations(securityIntegrations), nil
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
			var securityIntegrationsStrings []string
			if securityIntegrations != nil {
				securityIntegrationsStrings = make([]string, len(securityIntegrations))
				for i, v := range securityIntegrations {
					securityIntegrationsStrings[i] = v.Name()
				}
			}
			if err = handleExternalChangesToObjectInDescribe(d,
				describeMapping{"authentication_methods", "authentication_methods", authenticationPolicyDescriptions.Raw("AUTHENTICATION_METHODS"), authenticationMethods, nil},
				describeMapping{"mfa_enrollment", "mfa_enrollment", authenticationPolicyDescriptions.Raw("MFA_ENROLLMENT"), mfaEnrollment, nil},
				describeMapping{"client_types", "client_types", authenticationPolicyDescriptions.Raw("CLIENT_TYPES"), clientTypes, nil},
				describeMapping{"security_integrations", "security_integrations", authenticationPolicyDescriptions.Raw("SECURITY_INTEGRATIONS"), securityIntegrationsStrings, nil},
			); err != nil {
				return diag.FromErr(err)
			}

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

	// TODO(XXX): Remove this once the 2025_06 is generally enabled.
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
