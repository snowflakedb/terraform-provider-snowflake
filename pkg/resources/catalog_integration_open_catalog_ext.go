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
	catalogIntegrationOpenCatalogCreateExt        = TrackingCreateWrapper(resources.CatalogIntegrationOpenCatalog, CreateCatalogIntegrationOpenCatalog)
	catalogIntegrationOpenCatalogReadExt          = TrackingReadWrapper(resources.CatalogIntegrationOpenCatalog, ReadCatalogIntegrationOpenCatalogFunc(true))
	catalogIntegrationOpenCatalogUpdateExt        = TrackingUpdateWrapper(resources.CatalogIntegrationOpenCatalog, UpdateCatalogIntegrationOpenCatalog)
	catalogIntegrationOpenCatalogDescriptionExt   = "Resource used to manage catalog integration objects for Apache Iceberg™ tables that integrate with Snowflake Open Catalog. For more information, check [catalog integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-open-catalog)."
	catalogIntegrationOpenCatalogCustomizeDiffExt = TrackingCustomDiffWrapper(resources.CatalogIntegrationOpenCatalog, customdiff.All(
		ComputedIfAnyAttributeChanged(catalogIntegrationOpenCatalogSchema, ShowOutputAttributeName, "enabled", "comment"),
		ComputedIfAnyAttributeChanged(catalogIntegrationOpenCatalogSchema, DescribeOutputAttributeName, "enabled", "refresh_interval_seconds", "comment"),
		RecreateWhenCatalogSourceChangedExternally(sdk.CatalogIntegrationCatalogSourceTypePolaris),
	))
)

var catalogIntegrationOpenCatalogImporterExt = &schema.ResourceImporter{
	StateContext: TrackingImportWrapper(resources.CatalogIntegrationOpenCatalog, ImportCatalogIntegrationOpenCatalog),
}

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

// No config fields were passed to setStateToValuesFromConfig.
// set is the only alter payload. Catalog integration ALTER has no UNSET request type.
type catalogIntegrationOpenCatalogAlterRequestsExt struct {
	set *sdk.CatalogIntegrationSetRequest
}

func catalogIntegrationOpenCatalogParseIdExt(id string) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(id)
}

func catalogIntegrationOpenCatalogNewCreateRequestExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateOpenCatalogCatalogIntegrationRequest, error) {
	enabled := d.Get("enabled").(bool)

	restConfig, restConfigErr := parseOpenCatalogRestConfig(d, "rest_config.0.")
	restAuthentication, restAuthenticationErr := parseOpenCatalogRestAuthentication(d, "rest_authentication.0.")
	if err := errors.Join(restConfigErr, restAuthenticationErr); err != nil {
		return nil, err
	}

	return sdk.NewCreateOpenCatalogCatalogIntegrationRequest(id, restConfig, restAuthentication, enabled), nil
}

func catalogIntegrationOpenCatalogApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateOpenCatalogCatalogIntegrationRequest) error {
	return errors.Join(
		stringAttributeCreateBuilder(d, "catalog_namespace", request.WithCatalogNamespace),
		intAttributeCreateBuilder(d, "refresh_interval_seconds", request.WithRefreshIntervalSeconds),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
}

func catalogIntegrationOpenCatalogHandleExternalChangesInDescribeExt(d *schema.ResourceData, details *sdk.CatalogIntegrationOpenCatalogDetails, withExternalChangesMarking bool) diag.Diagnostics {
	if withExternalChangesMarking {
		if err := handleExternalChangesToObjectInFlatDescribe(
			d,
			catalogIntegrationOpenCatalogOutputMappingsExt(details)...,
		); err != nil {
			return diag.FromErr(err)
		}
		if err := handleExternalChangesToNestedAttrs(d, details); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := catalogIntegrationOpenCatalogSetConfigFieldsExt(d, details); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func catalogIntegrationOpenCatalogOutputMappingsExt(details *sdk.CatalogIntegrationOpenCatalogDetails) []outputMapping {
	return []outputMapping{
		{"refresh_interval_seconds", "refresh_interval_seconds", details.RefreshIntervalSeconds, details.RefreshIntervalSeconds, nil},
	}
}

func catalogIntegrationOpenCatalogSetConfigFieldsExt(d *schema.ResourceData, details *sdk.CatalogIntegrationOpenCatalogDetails) error {
	return errors.Join(
		d.Set("name", details.Id.Name()),
		d.Set("enabled", details.Enabled),
		// not reading refresh_interval_seconds on purpose (handled as external change to describe output)
		d.Set("comment", details.Comment),
		d.Set("catalog_source", string(details.CatalogSource)),
		d.Set("catalog_namespace", details.CatalogNamespace),
		// not reading rest_config and rest_authentication on purpose (set on import)
	)
}

func catalogIntegrationOpenCatalogApplySetUnsetExt(d *schema.ResourceData) (*catalogIntegrationOpenCatalogAlterRequestsExt, error) {
	set := sdk.NewCatalogIntegrationSetRequest()

	if errs := errors.Join(
		attributeMappedValueUpdateSetOnlyAlways(d, "rest_authentication.0.oauth_client_secret", &set.SetOAuthRestAuthentication, func(secret string) (sdk.CatalogIntegrationSetOauthRestAuthenticationRequest, error) {
			return *sdk.NewCatalogIntegrationSetOauthRestAuthenticationRequest(secret), nil
		}),
		handleCatalogIntegrationCommonPropsUpdate(d, set),
	); errs != nil {
		return nil, errs
	}
	return &catalogIntegrationOpenCatalogAlterRequestsExt{set: set}, nil
}

func catalogIntegrationOpenCatalogSetFieldsNotSetByReadExt(d *schema.ResourceData, _ sdk.AccountObjectIdentifier, details *sdk.CatalogIntegrationOpenCatalogDetails, importPermissiveMode bool) error {
	return errors.Join(
		importIntAttribute(d, "refresh_interval_seconds", importPermissiveMode, details.RefreshIntervalSeconds),
		setOpenCatalogRestConfigInState(d, details, importPermissiveMode),
		setOpenCatalogRestAuthInState(d, details, importPermissiveMode),
	)
}

// =============================================================================
// SDK incision points
// =============================================================================

func catalogIntegrationOpenCatalogCreateInSdkExt(ctx context.Context, meta any, _ sdk.AccountObjectIdentifier, req *sdk.CreateOpenCatalogCatalogIntegrationRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.CatalogIntegrations.CreateOpenCatalog(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Open Catalog catalog integration, err = %w", err))
	}
	return nil
}

func catalogIntegrationOpenCatalogShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegration, diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	s, err := client.CatalogIntegrations.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query Open Catalog catalog integration. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Open Catalog catalog integration id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return s, nil
}

func catalogIntegrationOpenCatalogShowByIdInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationOpenCatalogDetails, error) {
	client := meta.(*provider.Context).Client
	details, err := client.CatalogIntegrations.DescribeOpenCatalogDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if details.CatalogSource != sdk.CatalogIntegrationCatalogSourceTypePolaris {
		return nil, fmt.Errorf("invalid catalog source type, expected %s, got %s", sdk.CatalogIntegrationCatalogSourceTypePolaris, details.CatalogSource)
	}
	return details, nil
}

func catalogIntegrationOpenCatalogDescribeInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationOpenCatalogDetails, error) {
	client := meta.(*provider.Context).Client
	details, err := client.CatalogIntegrations.DescribeOpenCatalogDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not describe Open Catalog catalog integration (%s), err = %w", helpers.EncodeResourceIdentifier(id), err)
	}
	return details, nil
}

func catalogIntegrationOpenCatalogAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, reqs *catalogIntegrationOpenCatalogAlterRequestsExt) error {
	if reflect.DeepEqual(*reqs.set, *sdk.NewCatalogIntegrationSetRequest()) {
		return nil
	}
	client := meta.(*provider.Context).Client
	req := sdk.NewAlterCatalogIntegrationRequest(id).WithSet(*reqs.set)
	if err := client.CatalogIntegrations.Alter(ctx, req); err != nil {
		return fmt.Errorf("error updating Open Catalog catalog integration (%s), err = %w", helpers.EncodeResourceIdentifier(id), err)
	}
	return nil
}

func catalogIntegrationOpenCatalogAlterUnsetInSdkExt(context.Context, any, sdk.AccountObjectIdentifier, *catalogIntegrationOpenCatalogAlterRequestsExt) error {
	// Catalog integration ALTER has no UNSET clause.
	return nil
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

func parseOpenCatalogRestConfig(d *schema.ResourceData, prefix string) (sdk.CatalogIntegrationOpenCatalogRestConfigRequest, error) {
	req := sdk.NewCatalogIntegrationOpenCatalogRestConfigRequest(d.Get(prefix+"catalog_uri").(string), d.Get(prefix+"catalog_name").(string))
	err := errors.Join(
		attributeMappedValueCreateBuilder(d, prefix+"catalog_api_type", req.WithCatalogApiType, sdk.ToCatalogIntegrationCatalogApiType),
		attributeMappedValueCreateBuilder(d, prefix+"access_delegation_mode", req.WithAccessDelegationMode, sdk.ToCatalogIntegrationAccessDelegationMode),
	)
	return *req, err
}

// handleExternalChangesToNestedAttrs marks drift on rest_config and rest_authentication when
// the prior describe_output differs from the current Snowflake response.
func handleExternalChangesToNestedAttrs(d *schema.ResourceData, details *sdk.CatalogIntegrationOpenCatalogDetails) error {
	restConfig := []any{
		map[string]any{
			"catalog_uri":            details.RestConfig.CatalogUri,
			"catalog_name":           details.RestConfig.CatalogName,
			"catalog_api_type":       string(details.RestConfig.CatalogApiType),
			"access_delegation_mode": string(details.RestConfig.AccessDelegationMode),
		},
	}
	err := handleExternalChangesToObjectInFlatDescribeDeepEqual(
		d,
		outputMapping{"rest_config", "rest_config", restConfig, restConfig, nil},
	)
	if err != nil {
		return err
	}
	return handleExternalChangesToRestAuthentication(d, details)
}

func handleExternalChangesToRestAuthentication(d *schema.ResourceData, details *sdk.CatalogIntegrationOpenCatalogDetails) error {
	return handleExternalChangesToOAuthRestAuthentication(
		d,
		"rest_authentication",
		details.RestAuthentication.OauthTokenUri,
		details.RestAuthentication.OauthClientId,
		details.RestAuthentication.OauthAllowedScopes,
	)
}

func setOpenCatalogRestConfigInState(d *schema.ResourceData, details *sdk.CatalogIntegrationOpenCatalogDetails, importPermissiveMode bool) error {
	var catalogApiType, accessDelegationMode string
	if !importPermissiveMode {
		catalogApiType = string(details.RestConfig.CatalogApiType)
		accessDelegationMode = string(details.RestConfig.AccessDelegationMode)
	}
	return d.Set("rest_config", []any{map[string]any{
		"catalog_uri":            details.RestConfig.CatalogUri,
		"catalog_name":           details.RestConfig.CatalogName,
		"catalog_api_type":       catalogApiType,
		"access_delegation_mode": accessDelegationMode,
	}})
}

func setOpenCatalogRestAuthInState(d *schema.ResourceData, details *sdk.CatalogIntegrationOpenCatalogDetails, importPermissiveMode bool) error {
	var oAuthTokenUri string
	if !importPermissiveMode {
		oAuthTokenUri = details.RestAuthentication.OauthTokenUri
	}
	return d.Set("rest_authentication", []any{map[string]any{
		"oauth_token_uri": oAuthTokenUri,
		"oauth_client_id": details.RestAuthentication.OauthClientId,
		// oauth_client_secret not returned from Snowflake
		"oauth_allowed_scopes": details.RestAuthentication.OauthAllowedScopes,
	}})
}
