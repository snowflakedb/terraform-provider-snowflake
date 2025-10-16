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
		Description: fmt.Sprintf("A list of authentication methods that are allowed during login. This parameter accepts one or more of the following values: %s", possibleValuesListed(sdk.AllAuthenticationMethods)),
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
		Description:      "Determines whether a user must enroll in multi-factor authentication. Allowed values are REQUIRED and OPTIONAL. When REQUIRED is specified, Enforces users to enroll in MFA. If this value is used, then the CLIENT_TYPES parameter must include SNOWFLAKE_UI, because Snowsight is the only place users can enroll in multi-factor authentication (MFA).",
		ValidateDiagFunc: sdkValidation(sdk.ToMfaEnrollmentOption),
		Default:          "OPTIONAL",
	},
	"client_types": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: sdkValidation(sdk.ToClientTypesOption),
		},
		Optional:    true,
		Description: fmt.Sprintf("A list of clients that can authenticate with Snowflake. If a client tries to connect, and the client is not one of the valid CLIENT_TYPES, then the login attempt fails. Allowed values are %s. The CLIENT_TYPES property of an authentication policy is a best effort method to block user logins based on specific clients. It should not be used as the sole control to establish a security boundary.", possibleValuesListed(sdk.AllClientTypes)),
	},
	"security_integrations": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		Optional:    true,
		Description: "A list of security integrations the authentication policy is associated with. This parameter has no effect when SAML or OAUTH are not in the AUTHENTICATION_METHODS list. All values in the SECURITY_INTEGRATIONS list must be compatible with the values in the AUTHENTICATION_METHODS list. For example, if SECURITY_INTEGRATIONS contains a SAML security integration, and AUTHENTICATION_METHODS contains OAUTH, then you cannot create the authentication policy. To allow all security integrations use ALL as parameter.",
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
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.AuthenticationPolicyResource), TrackingReadWrapper(resources.AuthenticationPolicy, ReadContextAuthenticationPolicy)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.AuthenticationPolicyResource), TrackingUpdateWrapper(resources.AuthenticationPolicy, UpdateContextAuthenticationPolicy)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.AuthenticationPolicyResource), TrackingDeleteWrapper(resources.AuthenticationPolicy, deleteFunc)),
		Description:   "Resource used to manage authentication policy objects. For more information, check [authentication policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-authentication-policy).",

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

	authenticationPolicy, err := client.AuthenticationPolicies.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = d.Set("name", authenticationPolicy.Name); err != nil {
		return nil, err
	}
	if err = d.Set("database", authenticationPolicy.DatabaseName); err != nil {
		return nil, err
	}
	if err = d.Set("schema", authenticationPolicy.SchemaName); err != nil {
		return nil, err
	}
	if err = d.Set("comment", authenticationPolicy.Comment); err != nil {
		return nil, err
	}

	// needed as otherwise the resource will be incorrectly imported when a list-parameter value equals a default value
	authenticationPolicyDescriptions, err := client.AuthenticationPolicies.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	authenticationMethods := getListParameterFromDescribe(authenticationPolicyDescriptions, "AUTHENTICATION_METHODS")
	if err = d.Set("authentication_methods", authenticationMethods); err != nil {
		return nil, err
	}
	mfaAuthenticationMethods := getListParameterFromDescribe(authenticationPolicyDescriptions, "MFA_AUTHENTICATION_METHODS")
	if err = d.Set("mfa_authentication_methods", mfaAuthenticationMethods); err != nil {
		return nil, err
	}
	clientTypes := getListParameterFromDescribe(authenticationPolicyDescriptions, "CLIENT_TYPES")
	if err = d.Set("client_types", clientTypes); err != nil {
		return nil, err
	}
	securityIntegrations := getListParameterFromDescribe(authenticationPolicyDescriptions, "SECURITY_INTEGRATIONS")
	if err = d.Set("security_integrations", securityIntegrations); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
			authenticationMethods[i] = sdk.AuthenticationMethods{Method: *option}
		}
		req.WithAuthenticationMethods(authenticationMethods)
	}

	if v, ok := d.GetOk("mfa_authentication_methods"); ok {
		mfaAuthenticationMethodsRawList := expandStringList(v.(*schema.Set).List())
		mfaAuthenticationMethods := make([]sdk.MfaAuthenticationMethods, len(mfaAuthenticationMethodsRawList))
		for i, v := range mfaAuthenticationMethodsRawList {
			option, err := sdk.ToMfaAuthenticationMethodsOption(v)
			if err != nil {
				return diag.FromErr(err)
			}
			mfaAuthenticationMethods[i] = sdk.MfaAuthenticationMethods{Method: *option}
		}
		req.WithMfaAuthenticationMethods(mfaAuthenticationMethods)
	}

	if v, ok := d.GetOk("mfa_enrollment"); ok {
		option, err := sdk.ToMfaEnrollmentOption(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req = req.WithMfaEnrollment(*option)
	}

	if v, ok := d.GetOk("client_types"); ok {
		clientTypesRawList := expandStringList(v.(*schema.Set).List())
		clientTypes := make([]sdk.ClientTypes, len(clientTypesRawList))
		for i, v := range clientTypesRawList {
			option, err := sdk.ToClientTypesOption(v)
			if err != nil {
				return diag.FromErr(err)
			}
			clientTypes[i] = sdk.ClientTypes{ClientType: *option}
		}
		req.WithClientTypes(clientTypes)
	}

	if v, ok := d.GetOk("security_integrations"); ok {
		securityIntegrationsRawList := expandStringList(v.(*schema.Set).List())
		securityIntegrations := make([]sdk.SecurityIntegrationsOption, len(securityIntegrationsRawList))
		for i, v := range securityIntegrationsRawList {
			securityIntegrations[i] = sdk.SecurityIntegrationsOption{Name: sdk.NewAccountObjectIdentifier(v)}
		}
		req.WithSecurityIntegrations(securityIntegrations)
	}

	if v, ok := d.GetOk("comment"); ok {
		req = req.WithComment(v.(string))
	}

	client := meta.(*provider.Context).Client
	if err := client.AuthenticationPolicies.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextAuthenticationPolicy(ctx, d, meta)
}

func ReadContextAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	authenticationPolicyDescriptions, err := client.AuthenticationPolicies.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	authenticationMethods := getListArgumentWithDefaults(d, "authentication_methods", getListParameterFromDescribe(authenticationPolicyDescriptions, "AUTHENTICATION_METHODS"), []string{"ALL"})
	if err = d.Set("authentication_methods", authenticationMethods); err != nil {
		return diag.FromErr(err)
	}

	mfaAuthenticationMethods := getListArgumentWithDefaults(d, "mfa_authentication_methods", getListParameterFromDescribe(authenticationPolicyDescriptions, "MFA_AUTHENTICATION_METHODS"), []string{"PASSWORD", "SAML"})
	if err = d.Set("mfa_authentication_methods", mfaAuthenticationMethods); err != nil {
		return diag.FromErr(err)
	}

	mfaEnrollment, err := collections.FindFirst(authenticationPolicyDescriptions, func(prop sdk.AuthenticationPolicyDescription) bool { return prop.Property == "MFA_ENROLLMENT" })
	if err == nil {
		if err = d.Set("mfa_enrollment", mfaEnrollment.Value); err != nil {
			return diag.FromErr(err)
		}
	}

	clientTypes := getListArgumentWithDefaults(d, "client_types", getListParameterFromDescribe(authenticationPolicyDescriptions, "CLIENT_TYPES"), []string{"ALL"})
	if err = d.Set("client_types", clientTypes); err != nil {
		return diag.FromErr(err)
	}

	securityIntegrations := getListArgumentWithDefaults(d, "security_integrations", getListParameterFromDescribe(authenticationPolicyDescriptions, "SECURITY_INTEGRATIONS"), []string{"ALL"})
	if err = d.Set("security_integrations", securityIntegrations); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("comment", authenticationPolicy.Comment); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.AuthenticationPolicyToSchema(authenticationPolicy)}); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.AuthenticationPolicyDescriptionToSchema(authenticationPolicyDescriptions)}); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func getListParameterFromDescribe(authenticationPolicyDescriptions []sdk.AuthenticationPolicyDescription, parameterName string) []string {
	parameterList := make([]string, 0)
	if parameterProperty, err := collections.FindFirst(authenticationPolicyDescriptions, func(prop sdk.AuthenticationPolicyDescription) bool {
		return prop.Property == parameterName
	}); err == nil {
		parameterList = append(parameterList, sdk.ParseCommaSeparatedStringArray(parameterProperty.Value, false)...)
	}
	return parameterList
}

// getListArgumentWithDefaults returns the list of values for a given argument, with the defaults applied, if necessary. Otherwise, tf plan will always show a diff with a list parameter with defaults when no value is set.
func getListArgumentWithDefaults(d *schema.ResourceData, argumentName string, argumentIs []string, argumentDefaults []string) []string {
	// in case nothing is set in the tf resource and the is equals the default, we set the is to empty
	argumentShould := d.Get(argumentName).(*schema.Set).List()
	if stringSlicesEqual(argumentIs, argumentDefaults) && len(argumentShould) == 0 {
		argumentIs = []string{}
	}
	return argumentIs
}

func UpdateContextAuthenticationPolicy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewAuthenticationPolicySetRequest(), sdk.NewAuthenticationPolicyUnsetRequest()

	// change to name
	if d.HasChange("name") {
		newId, err := sdk.ParseSchemaObjectIdentifier(d.Get("name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

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
				authenticationMethodsValues[i] = sdk.AuthenticationMethods{Method: *option}
			}

			set.WithAuthenticationMethods(authenticationMethodsValues)
		} else {
			unset.WithAuthenticationMethods(true)
		}
	}

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
				mfaAuthenticationMethodsValues[i] = sdk.MfaAuthenticationMethods{Method: *option}
			}

			set.WithMfaAuthenticationMethods(mfaAuthenticationMethodsValues)
		} else {
			unset.WithMfaAuthenticationMethods(true)
		}
	}

	// change to mfa enrollment
	if d.HasChange("mfa_enrollment") {
		if mfaEnrollmentOption, err := sdk.ToMfaEnrollmentOption(d.Get("mfa_enrollment").(string)); err == nil {
			set.WithMfaEnrollment(*mfaEnrollmentOption)
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
				clientTypesValues[i] = sdk.ClientTypes{ClientType: *option}
			}

			set.WithClientTypes(clientTypesValues)
		} else {
			unset.WithClientTypes(true)
		}
	}

	// change to security integrations
	if d.HasChange("security_integrations") {
		if v, ok := d.GetOk("security_integrations"); ok {
			securityIntegrations := expandStringList(v.(*schema.Set).List())
			securityIntegrationsValues := make([]sdk.SecurityIntegrationsOption, len(securityIntegrations))
			for i, v := range securityIntegrations {
				securityIntegrationsValues[i] = sdk.SecurityIntegrationsOption{Name: sdk.NewAccountObjectIdentifier(v)}
			}

			set.WithSecurityIntegrations(securityIntegrationsValues)
		} else {
			unset.WithSecurityIntegrations(true)
		}
	}

	// change to comment
	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.Comment = sdk.String(v.(string))
		} else {
			unset.WithComment(true)
		}
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

	return ReadContextAuthenticationPolicy(ctx, d, meta)
}

func stringSlicesEqual(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	// convert slices to maps for easy comparison
	set1 := make(map[string]bool)
	for _, v := range s1 {
		set1[v] = true
	}

	set2 := make(map[string]bool)
	for _, v := range s2 {
		set2[v] = true
	}

	for k := range set1 {
		if _, ok := set2[k]; !ok {
			return false
		}
	}
	return true
}
