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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	catalogIntegrationIcebergRestCreateExt        = TrackingCreateWrapper(resources.CatalogIntegrationIcebergRest, CreateCatalogIntegrationIcebergRest)
	catalogIntegrationIcebergRestReadExt          = TrackingReadWrapper(resources.CatalogIntegrationIcebergRest, ReadCatalogIntegrationIcebergRestFunc(true))
	catalogIntegrationIcebergRestUpdateExt        = TrackingUpdateWrapper(resources.CatalogIntegrationIcebergRest, UpdateCatalogIntegrationIcebergRest)
	catalogIntegrationIcebergRestDescriptionExt   = "Resource used to manage catalog integration objects for Apache Iceberg™ tables managed in a remote catalog that complies with the open source Apache Iceberg™ REST OpenAPI specification. For more information, check [catalog integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-rest)."
	catalogIntegrationIcebergRestCustomizeDiffExt = TrackingCustomDiffWrapper(resources.CatalogIntegrationIcebergRest, customdiff.All(
		ComputedIfAnyAttributeChanged(catalogIntegrationIcebergRestSchema, ShowOutputAttributeName, "enabled", "comment"),
		ComputedIfAnyAttributeChanged(catalogIntegrationIcebergRestSchema, DescribeOutputAttributeName, "enabled", "refresh_interval_seconds", "comment"),
		RecreateWhenCatalogSourceChangedExternally(sdk.CatalogIntegrationCatalogSourceTypeIcebergRest),
	))
	catalogIntegrationIcebergRestImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.CatalogIntegrationIcebergRest, ImportCatalogIntegrationIcebergRest),
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

type catalogIntegrationIcebergRestAlterRequestsExt struct {
	set *sdk.CatalogIntegrationSetRequest
}

func catalogIntegrationIcebergRestParseIdExt(id string) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(id)
}

func catalogIntegrationIcebergRestNewCreateRequestExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateIcebergRestCatalogIntegrationRequest, error) {
	enabled := d.Get("enabled").(bool)

	restConfig, err := parseIcebergRestConfig(d, "rest_config.0.")
	if err != nil {
		return nil, err
	}
	request := sdk.NewCreateIcebergRestCatalogIntegrationRequest(id, restConfig, enabled)
	if err := errors.Join(
		singleBlockAttributeCreate(d, "oauth_rest_authentication", request.WithOAuthRestAuthentication, parseIcebergRestOAuthRestAuthentication),
		singleBlockAttributeCreate(d, "bearer_rest_authentication", request.WithBearerRestAuthentication, parseBearerRestAuthentication),
		singleBlockAttributeCreate(d, "sigv4_rest_authentication", request.WithSigV4RestAuthentication, parseSigV4RestAuthentication),
	); err != nil {
		return nil, err
	}
	if request.OAuthRestAuthentication == nil && request.BearerRestAuthentication == nil && request.SigV4RestAuthentication == nil {
		return nil, fmt.Errorf("one of oauth_rest_authentication, bearer_rest_authentication, or sigv4_rest_authentication must be set")
	}

	return request, nil
}

func catalogIntegrationIcebergRestApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateIcebergRestCatalogIntegrationRequest) error {
	return errors.Join(
		stringAttributeCreateBuilder(d, "catalog_namespace", request.WithCatalogNamespace),
		intAttributeCreateBuilder(d, "refresh_interval_seconds", request.WithRefreshIntervalSeconds),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
}

func catalogIntegrationIcebergRestHandleExternalChangesInDescribeExt(d *schema.ResourceData, details *sdk.CatalogIntegrationIcebergRestDetails, withExternalChangesMarking bool) diag.Diagnostics {
	if withExternalChangesMarking {
		if err := handleExternalChangesToObjectInFlatDescribe(
			d,
			catalogIntegrationIcebergRestOutputMappingsExt(details)...,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := handleExternalChangesToIcebergRestNestedAttrs(d, details); err != nil {
			return diag.FromErr(err)
		}
	}
	// refresh_interval_seconds is not copied into config; drift is marked on describe output.
	if err := errors.Join(
		d.Set("catalog_source", string(details.CatalogSource)),
		d.Set("catalog_namespace", details.CatalogNamespace),
	); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func catalogIntegrationIcebergRestOutputMappingsExt(details *sdk.CatalogIntegrationIcebergRestDetails) []outputMapping {
	return []outputMapping{
		{"refresh_interval_seconds", "refresh_interval_seconds", details.RefreshIntervalSeconds, details.RefreshIntervalSeconds, nil},
	}
}

func catalogIntegrationIcebergRestSetConfigFieldsExt(d *schema.ResourceData, catalogIntegration *sdk.CatalogIntegration) error {
	return errors.Join(
		d.Set("name", catalogIntegration.Name),
		d.Set("enabled", catalogIntegration.Enabled),
		d.Set("comment", catalogIntegration.Comment),
	)
}

func catalogIntegrationIcebergRestApplySetUnsetExt(d *schema.ResourceData) (*catalogIntegrationIcebergRestAlterRequestsExt, error) {
	set := sdk.NewCatalogIntegrationSetRequest()

	if errs := errors.Join(
		attributeMappedValueUpdateSetOnlyAlways(d, "oauth_rest_authentication.0.oauth_client_secret", &set.SetOAuthRestAuthentication, func(secret string) (sdk.CatalogIntegrationSetOauthRestAuthenticationRequest, error) {
			return *sdk.NewCatalogIntegrationSetOauthRestAuthenticationRequest(secret), nil
		}),
		attributeMappedValueUpdateSetOnlyAlways(d, "bearer_rest_authentication.0.bearer_token", &set.SetBearerRestAuthentication, func(token string) (sdk.CatalogIntegrationSetBearerRestAuthenticationRequest, error) {
			return *sdk.NewCatalogIntegrationSetBearerRestAuthenticationRequest(token), nil
		}),
		handleCatalogIntegrationCommonPropsUpdate(d, set),
	); errs != nil {
		return nil, errs
	}
	return &catalogIntegrationIcebergRestAlterRequestsExt{set: set}, nil
}

func catalogIntegrationIcebergRestSetFieldsNotSetByReadExt(d *schema.ResourceData, _ sdk.AccountObjectIdentifier, details *sdk.CatalogIntegrationIcebergRestDetails, importPermissiveMode bool) error {
	return errors.Join(
		importIntAttribute(d, "refresh_interval_seconds", importPermissiveMode, details.RefreshIntervalSeconds),
		setIcebergRestConfigInState(d, details, importPermissiveMode),
		setIcebergRestAuthInState(d, details, importPermissiveMode),
	)
}

// =============================================================================
// SDK incision points
// =============================================================================

func catalogIntegrationIcebergRestCreateInSdkExt(ctx context.Context, meta any, _ sdk.AccountObjectIdentifier, req *sdk.CreateIcebergRestCatalogIntegrationRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.CatalogIntegrations.CreateIcebergRest(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Iceberg REST catalog integration, err = %w", err))
	}
	return nil
}

func catalogIntegrationIcebergRestShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegration, diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	s, err := client.CatalogIntegrations.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query Iceberg REST catalog integration. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Iceberg REST catalog integration id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return s, nil
}

// Describe plus the ICEBERG_REST catalog-source guard. Import does not call ShowByID.
func catalogIntegrationIcebergRestShowByIdInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationIcebergRestDetails, error) {
	client := meta.(*provider.Context).Client
	details, err := client.CatalogIntegrations.DescribeIcebergRestDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.CatalogSource != sdk.CatalogIntegrationCatalogSourceTypeIcebergRest {
		return nil, fmt.Errorf("invalid catalog source type, expected %s, got %s", sdk.CatalogIntegrationCatalogSourceTypeIcebergRest, details.CatalogSource)
	}
	return details, nil
}

func catalogIntegrationIcebergRestDescribeInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationIcebergRestDetails, error) {
	client := meta.(*provider.Context).Client
	details, err := client.CatalogIntegrations.DescribeIcebergRestDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not describe Iceberg REST catalog integration (%s), err = %w", helpers.EncodeResourceIdentifier(id), err)
	}
	return details, nil
}

func catalogIntegrationIcebergRestAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, reqs *catalogIntegrationIcebergRestAlterRequestsExt) error {
	client := meta.(*provider.Context).Client
	if !reflect.DeepEqual(*reqs.set, *sdk.NewCatalogIntegrationSetRequest()) {
		req := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*reqs.set)
		if err := client.CatalogIntegrations.Alter(ctx, req); err != nil {
			return fmt.Errorf("error updating Iceberg REST catalog integration (%s), err = %w", helpers.EncodeResourceIdentifier(id), err)
		}
	}
	return nil
}

func catalogIntegrationIcebergRestAlterUnsetInSdkExt(context.Context, any, sdk.AccountObjectIdentifier, *catalogIntegrationIcebergRestAlterRequestsExt) error {
	// Catalog integration ALTER has no UNSET request. Update still calls this hook after SET.
	return nil
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

func parseIcebergRestConfig(d *schema.ResourceData, prefix string) (sdk.CatalogIntegrationIcebergRestRestConfigRequest, error) {
	req := sdk.NewCatalogIntegrationIcebergRestRestConfigRequest(d.Get(prefix + "catalog_uri").(string))
	err := errors.Join(
		stringAttributeCreateBuilder(d, prefix+"prefix", req.WithPrefix),
		stringAttributeCreateBuilder(d, prefix+"catalog_name", req.WithCatalogName),
		attributeMappedValueCreateBuilder(d, prefix+"catalog_api_type", req.WithCatalogApiType, sdk.ToCatalogIntegrationCatalogApiType),
		attributeMappedValueCreateBuilder(d, prefix+"access_delegation_mode", req.WithAccessDelegationMode, sdk.ToCatalogIntegrationAccessDelegationMode),
	)
	return *req, err
}

func parseBearerRestAuthentication(d *schema.ResourceData, prefix string) (sdk.CatalogIntegrationIcebergRestBearerRestAuthenticationRequest, error) {
	return *sdk.NewCatalogIntegrationIcebergRestBearerRestAuthenticationRequest(d.Get(prefix + "bearer_token").(string)), nil
}

func parseSigV4RestAuthentication(d *schema.ResourceData, prefix string) (sdk.CatalogIntegrationIcebergRestSigv4RestAuthenticationRequest, error) {
	req := sdk.NewCatalogIntegrationIcebergRestSigv4RestAuthenticationRequest(d.Get(prefix + "sigv4_iam_role").(string))
	err := errors.Join(
		stringAttributeCreateBuilder(d, prefix+"sigv4_signing_region", req.WithSigv4SigningRegion),
		stringAttributeCreateBuilder(d, prefix+"sigv4_external_id", req.WithSigv4ExternalId),
	)
	return *req, err
}

func handleExternalChangesToIcebergRestNestedAttrs(d *schema.ResourceData, details *sdk.CatalogIntegrationIcebergRestDetails) error {
	if err := handleExternalChangesToIcebergRestConfig(d, details); err != nil {
		return err
	}
	if details.OauthRestAuthentication != nil {
		return handleExternalChangesToOAuthRestAuthenticationIcebergRest(d, details)
	}
	if details.Sigv4RestAuthentication != nil {
		return handleExternalChangesToSigV4RestAuthentication(d, details)
	}
	return handleExternalChangesToBearerRestAuthentication(d)
}

func handleExternalChangesToIcebergRestConfig(d *schema.ResourceData, details *sdk.CatalogIntegrationIcebergRestDetails) error {
	restConfig := []any{
		map[string]any{
			"catalog_uri":            details.RestConfig.CatalogUri,
			"prefix":                 details.RestConfig.Prefix,
			"catalog_name":           details.RestConfig.CatalogName,
			"catalog_api_type":       string(details.RestConfig.CatalogApiType),
			"access_delegation_mode": string(details.RestConfig.AccessDelegationMode),
		},
	}
	return handleExternalChangesToObjectInFlatDescribeDeepEqual(
		d,
		outputMapping{"rest_config", "rest_config", restConfig, restConfig, nil},
	)
}

func handleExternalChangesToOAuthRestAuthenticationIcebergRest(d *schema.ResourceData, details *sdk.CatalogIntegrationIcebergRestDetails) error {
	return errors.Join(
		handleExternalChangesToOAuthRestAuthentication(
			d,
			"oauth_rest_authentication",
			details.OauthRestAuthentication.OauthTokenUri,
			details.OauthRestAuthentication.OauthClientId,
			details.OauthRestAuthentication.OauthAllowedScopes,
		),
		d.Set("bearer_rest_authentication", nil),
		d.Set("sigv4_rest_authentication", nil),
	)
}

func handleExternalChangesToSigV4RestAuthentication(d *schema.ResourceData, details *sdk.CatalogIntegrationIcebergRestDetails) error {
	sigV4RestAuthorization := []any{
		map[string]any{
			"sigv4_iam_role":       details.Sigv4RestAuthentication.Sigv4IamRole,
			"sigv4_signing_region": details.Sigv4RestAuthentication.Sigv4SigningRegion,
		},
	}
	err := handleExternalChangesToObjectInFlatDescribeDeepEqual(
		d,
		outputMapping{"sigv4_rest_authentication", "sigv4_rest_authentication", sigV4RestAuthorization, sigV4RestAuthorization, nil},
	)
	return errors.Join(
		err,
		d.Set("oauth_rest_authentication", nil),
		d.Set("bearer_rest_authentication", nil),
	)
}

func handleExternalChangesToBearerRestAuthentication(d *schema.ResourceData) error {
	return errors.Join(d.Set("oauth_rest_authentication", nil), d.Set("sigv4_rest_authentication", nil))
}

func setIcebergRestConfigInState(d *schema.ResourceData, details *sdk.CatalogIntegrationIcebergRestDetails, importPermissiveMode bool) error {
	var catalogApiType, accessDelegationMode string
	if !importPermissiveMode {
		catalogApiType = string(details.RestConfig.CatalogApiType)
		accessDelegationMode = string(details.RestConfig.AccessDelegationMode)
	}
	return d.Set("rest_config", []any{map[string]any{
		"catalog_uri":            details.RestConfig.CatalogUri,
		"prefix":                 details.RestConfig.Prefix,
		"catalog_name":           details.RestConfig.CatalogName,
		"catalog_api_type":       catalogApiType,
		"access_delegation_mode": accessDelegationMode,
	}})
}

func setIcebergRestAuthInState(d *schema.ResourceData, details *sdk.CatalogIntegrationIcebergRestDetails, importPermissiveMode bool) error {
	if details.OauthRestAuthentication != nil {
		var oAuthTokenUri string
		if !importPermissiveMode {
			oAuthTokenUri = details.OauthRestAuthentication.OauthTokenUri
		}
		return errors.Join(
			d.Set("oauth_rest_authentication", []any{map[string]any{
				"oauth_token_uri": oAuthTokenUri,
				"oauth_client_id": details.OauthRestAuthentication.OauthClientId,
				// oauth_client_secret not returned from Snowflake
				"oauth_allowed_scopes": details.OauthRestAuthentication.OauthAllowedScopes,
			}}),
		)
	}
	if details.Sigv4RestAuthentication != nil {
		var sigV4SigningRegion string
		if !importPermissiveMode {
			sigV4SigningRegion = details.Sigv4RestAuthentication.Sigv4SigningRegion
		}
		return errors.Join(
			d.Set("sigv4_rest_authentication", []any{map[string]any{
				"sigv4_iam_role":       details.Sigv4RestAuthentication.Sigv4IamRole,
				"sigv4_signing_region": sigV4SigningRegion,
				// sigv4_external_id not returned from Snowflake
			}}),
		)
	}

	// bearer_rest_authentication.bearer_token not returned from Snowflake

	return nil
}
