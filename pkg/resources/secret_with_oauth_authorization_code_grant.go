package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var secretAuthorizationCodeGrantSchema = func() map[string]*schema.Schema {
	secretAuthorizationCodeGrant := map[string]*schema.Schema{
		"oauth_refresh_token": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: externalChangesNotDetectedFieldDescription("Specifies the token as a string that is used to obtain a new access token from the OAuth authorization server when the access token expires."),
		},
		"oauth_refresh_token_expiry_time": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("oauth_refresh_token_expiry_time"),
			Description:      "Specifies the timestamp as a string when the OAuth refresh token expires. Accepted string formats: YYYY-MM-DD, YYYY-MM-DD HH:MI, YYYY-MM-DD HH:MI:SS, YYYY-MM-DD HH:MI <timezone>",
		},
		"api_authentication": {
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			Required:         true,
			Description:      relatedResourceDescription("Specifies the name value of the Snowflake security integration that connects Snowflake to an external service.", resources.ApiAuthenticationIntegrationWithAuthorizationCodeGrant),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
	}
	return collections.MergeMaps(secretCommonSchema, secretAuthorizationCodeGrant)
}()

func SecretWithAuthorizationCodeGrant() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.SecretWithAuthorizationCodeGrant, CreateContextSecretWithAuthorizationCodeGrant),
		ReadContext:   TrackingReadWrapper(resources.SecretWithAuthorizationCodeGrant, ReadContextSecretWithAuthorizationCodeGrant(true)),
		UpdateContext: TrackingUpdateWrapper(resources.SecretWithAuthorizationCodeGrant, UpdateContextSecretWithAuthorizationCodeGrant),
		DeleteContext: TrackingDeleteWrapper(resources.SecretWithAuthorizationCodeGrant, DeleteContextSecret),
		Description:   "Resource used to manage secret objects with OAuth Authorization Code Grant. For more information, check [secret documentation](https://docs.snowflake.com/en/sql-reference/sql/create-secret).",

		Schema: secretAuthorizationCodeGrantSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.SecretWithAuthorizationCodeGrant, ImportSecretWithAuthorizationCodeGrant),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.SecretWithAuthorizationCodeGrant, customdiff.All(
			ComputedIfAnyAttributeChanged(secretAuthorizationCodeGrantSchema, ShowOutputAttributeName, "comment"),
			ComputedIfAnyAttributeChanged(secretAuthorizationCodeGrantSchema, DescribeOutputAttributeName, "oauth_refresh_token_expiry_time", "api_authentication"),
			RecreateWhenSecretTypeChangedExternally(sdk.SecretTypeOAuth2AuthorizationCodeGrant),
		)),
		Timeouts: defaultTimeouts,
	}
}

func ImportSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	if err = handleSecretImport(d); err != nil {
		return nil, err
	}

	secretDescription, err := client.Secrets.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = d.Set("oauth_refresh_token_expiry_time", secretDescription.OauthRefreshTokenExpiryTime.String()); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName, schemaName, name := d.Get("database").(string), d.Get("schema").(string), d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	apiIntegrationString := d.Get("api_authentication").(string)
	apiIntegration, err := sdk.ParseAccountObjectIdentifier(apiIntegrationString)
	if err != nil {
		return diag.FromErr(err)
	}

	refreshToken := d.Get("oauth_refresh_token").(string)
	refreshTokenExpiryTime := d.Get("oauth_refresh_token_expiry_time").(string)

	request := sdk.NewCreateWithOAuthAuthorizationCodeFlowSecretRequest(id, refreshToken, refreshTokenExpiryTime, apiIntegration)
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	err = client.Secrets.CreateWithOAuthAuthorizationCodeFlow(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSecretWithAuthorizationCodeGrant(false)(ctx, d, meta)
}

func ReadContextSecretWithAuthorizationCodeGrant(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		secret, err := client.Secrets.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query secret with authorization code grant. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Secret with authorization code grant id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to query secret with authorization code grant.",
					Detail:   fmt.Sprintf("Secret with authorization code grant id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}

		secretDescription, err := client.Secrets.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		// if secret type is changed externally, we wont be able to read oauth_refresh_token_expiry_time value (since it will not be provided)
		// in any other case, there should be oauth_refresh_token_expiry_time value since it is required
		if withExternalChangesMarking && secretDescription.OauthRefreshTokenExpiryTime != nil {
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"oauth_refresh_token_expiry_time", "oauth_refresh_token_expiry_time", secretDescription.OauthRefreshTokenExpiryTime.String(), secretDescription.OauthRefreshTokenExpiryTime.String(), nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		return diag.FromErr(errors.Join(
			handleSecretRead(d, id, secret, secretDescription),
			setStateToValuesFromConfig(d, secretAuthorizationCodeGrantSchema, []string{"oauth_refresh_token_expiry_time"}),
			d.Set("api_authentication", secretDescription.IntegrationName),
		))
	}
}

func UpdateContextSecretWithAuthorizationCodeGrant(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := &sdk.SecretSetRequest{}
	unset := &sdk.SecretUnsetRequest{}
	handleSecretUpdate(d, set, unset)
	setForOAuthAuthorization := &sdk.SetForOAuthAuthorizationRequest{}

	if d.HasChange("oauth_refresh_token") {
		refreshToken := d.Get("oauth_refresh_token").(string)
		setForOAuthAuthorization.WithOauthRefreshToken(refreshToken)
	}

	if d.HasChange("oauth_refresh_token_expiry_time") {
		refreshTokenExpiryTime := d.Get("oauth_refresh_token_expiry_time").(string)
		setForOAuthAuthorization.WithOauthRefreshTokenExpiryTime(refreshTokenExpiryTime)
	}
	if !reflect.DeepEqual(setForOAuthAuthorization, sdk.SetForOAuthAuthorizationRequest{}) {
		set.WithSetForFlow(sdk.SetForFlowRequest{SetForOAuthAuthorization: setForOAuthAuthorization})
	}

	if !reflect.DeepEqual(*set, sdk.SecretSetRequest{}) {
		if err := client.Secrets.Alter(ctx, sdk.NewAlterSecretRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if !reflect.DeepEqual(*unset, sdk.SecretUnsetRequest{}) {
		if err := client.Secrets.Alter(ctx, sdk.NewAlterSecretRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSecretWithAuthorizationCodeGrant(false)(ctx, d, meta)
}
