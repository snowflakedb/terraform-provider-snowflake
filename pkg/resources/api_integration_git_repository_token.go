package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var apiIntegrationGitRepositoryAllowedValueExactlyOneOf = []string{
	"all_allowed_authentication_secrets",
	"no_allowed_authentication_secrets",
	"allowed_authentication_secrets",
}

var apiIntegrationGitRepositoryTokenSchema = func() map[string]*schema.Schema {
	apiIntegrationGitRepositoryToken := map[string]*schema.Schema{
		"all_allowed_authentication_secrets": {
			Type:         schema.TypeBool,
			Optional:     true,
			ExactlyOneOf: apiIntegrationGitRepositoryAllowedValueExactlyOneOf,
			Description:  "When set to true, all authentication secrets are allowed to be used when authenticating to the git repository. Exactly one of `all_allowed_authentication_secrets`, `no_allowed_authentication_secrets`, or `allowed_authentication_secrets` must be set.",
		},
		"no_allowed_authentication_secrets": {
			Type:         schema.TypeBool,
			Optional:     true,
			ExactlyOneOf: apiIntegrationGitRepositoryAllowedValueExactlyOneOf,
			Description:  "When set to true, no authentication secrets are allowed to be used when authenticating to the git repository. Exactly one of `all_allowed_authentication_secrets`, `no_allowed_authentication_secrets`, or `allowed_authentication_secrets` must be set.",
		},
		"allowed_authentication_secrets": {
			Type:             schema.TypeSet,
			Optional:         true,
			ExactlyOneOf:     apiIntegrationGitRepositoryAllowedValueExactlyOneOf,
			DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("allowed_authentication_secrets"),
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			},
			Description: "A list of fully-qualified secret identifiers (database.schema.secret) allowed to be used when authenticating to the git repository. Exactly one of `all_allowed_authentication_secrets`, `no_allowed_authentication_secrets`, or `allowed_authentication_secrets` must be set.",
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE API INTEGRATION` for the given integration.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeGitRepositoryTokenApiIntegrationSchema,
			},
		},
	}
	return collections.MergeMaps(apiIntegrationCommonSchema, apiIntegrationGitRepositoryToken)
}()

func ApiIntegrationGitRepositoryToken() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ApiIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ApiIntegrationGitRepositoryToken, CreateApiIntegrationGitRepositoryToken),
		ReadContext:   TrackingReadWrapper(resources.ApiIntegrationGitRepositoryToken, ReadApiIntegrationGitRepositoryToken),
		UpdateContext: TrackingUpdateWrapper(resources.ApiIntegrationGitRepositoryToken, UpdateApiIntegrationGitRepositoryToken),
		DeleteContext: TrackingDeleteWrapper(resources.ApiIntegrationGitRepositoryToken, deleteFunc),
		Description:   "Resource used to manage API integration for git HTTPS API with token-based authentication. For more information, check [api integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-api-integration).",

		Schema: apiIntegrationGitRepositoryTokenSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ApiIntegrationGitRepositoryToken, ImportApiIntegrationGitRepositoryToken),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(apiIntegrationGitRepositoryTokenSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(apiIntegrationGitRepositoryTokenSchema, DescribeOutputAttributeName, "enabled", "api_allowed_prefixes", "api_blocked_prefixes", "comment", "all_allowed_authentication_secrets", "no_allowed_authentication_secrets", "allowed_authentication_secrets"),
		),
	}
}

func ImportApiIntegrationGitRepositoryToken(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.ApiIntegrations.DescribeGitHttpsApiDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not describe API integration %s during import: %w", id.FullyQualifiedName(), err)
	}

	if _, err := sdk.ToApiIntegrationGitApiProviderType(details.ApiProvider); err != nil {
		return nil, fmt.Errorf(
			"api integration %s has api_provider %q, not compatible with snowflake_api_integration_git_repository_token; use the appropriate resource type",
			id.FullyQualifiedName(),
			details.ApiProvider,
		)
	}

	return ImportName[sdk.AccountObjectIdentifier](ctx, d, meta)
}

func buildAllowedAuthSecretsRequestFromState(d *schema.ResourceData) (*sdk.ApiIntegrationAllowedAuthenticationSecretsRequest, error) {
	req := sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest()
	if v, ok := d.GetOk("all_allowed_authentication_secrets"); ok && v.(bool) {
		return req.WithAllSecrets(true), nil
	}
	if v, ok := d.GetOk("no_allowed_authentication_secrets"); ok && v.(bool) {
		return req.WithNoSecrets(true), nil
	}
	if v, ok := d.GetOk("allowed_authentication_secrets"); ok {
		rawList := expandStringList(v.(*schema.Set).List())
		ids, err := collections.MapErr(rawList, sdk.ParseSchemaObjectIdentifier)
		if err != nil {
			return nil, err
		}
		return req.WithAllowedList(ids), nil
	}
	return nil, fmt.Errorf("exactly one of all_allowed_authentication_secrets, no_allowed_authentication_secrets, or allowed_authentication_secrets must be set")
}

func setAllowedAuthSecretFieldsFromDescribe(d *schema.ResourceData, raw string) error {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "ALL":
		return errors.Join(
			d.Set("all_allowed_authentication_secrets", true),
			d.Set("no_allowed_authentication_secrets", false),
			d.Set("allowed_authentication_secrets", []any{}),
		)
	case "NONE":
		return errors.Join(
			d.Set("all_allowed_authentication_secrets", false),
			d.Set("no_allowed_authentication_secrets", true),
			d.Set("allowed_authentication_secrets", []any{}),
		)
	default:
		ids, err := collections.MapErr(sdk.ParseCommaSeparatedStringArray(raw, true), sdk.ParseSchemaObjectIdentifier)
		if err != nil {
			return err
		}
		return errors.Join(
			d.Set("all_allowed_authentication_secrets", false),
			d.Set("no_allowed_authentication_secrets", false),
			d.Set("allowed_authentication_secrets", collections.Map(ids, sdk.SchemaObjectIdentifier.FullyQualifiedName)),
		)
	}
}

func CreateApiIntegrationGitRepositoryToken(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, request, err := handleApiIntegrationCommonCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	gitParams := sdk.NewGitHttpsApiTokenBasedParamsRequest()

	secretsReq, err := buildAllowedAuthSecretsRequestFromState(d)
	if err != nil {
		return diag.FromErr(err)
	}
	gitParams.WithAllowedAuthenticationSecrets(*secretsReq)

	if err = client.ApiIntegrations.Create(ctx, request.WithGitHttpsApiTokenBasedProviderParams(*gitParams)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating git HTTPS API token-based API integration: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadApiIntegrationGitRepositoryToken(ctx, d, meta)
}

func ReadApiIntegrationGitRepositoryToken(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	s, err := client.ApiIntegrations.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query API integration git HTTPS API token-based. Marking the resource as removed.",
					Detail:   fmt.Sprintf("API integration git HTTPS API token-based id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	gitDetails, err := client.ApiIntegrations.DescribeGitHttpsApiDetails(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe API integration git HTTPS API token-based (%s): %w", d.Id(), err))
	}

	errs := errors.Join(
		handleApiIntegrationCommonRead(d, id, s, gitDetails.AllowedPrefixes, gitDetails.BlockedPrefixes),
		setAllowedAuthSecretFieldsFromDescribe(d, gitDetails.AllowedAuthenticationSecrets),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ApiIntegrationGitRepositoryTokenDetailsToSchema(gitDetails)}),
	)
	return diag.FromErr(errs)
}

func UpdateApiIntegrationGitRepositoryToken(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewApiIntegrationSetRequest()
	unset := sdk.NewApiIntegrationUnsetRequest()
	gitSet := sdk.NewSetGitHttpsApiTokenBasedParamsRequest()
	gitUnset := sdk.NewUnsetGitHttpsApiTokenBasedParamsRequest()

	if err := handleApiIntegrationCommonUpdate(d, set, unset); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("all_allowed_authentication_secrets", "no_allowed_authentication_secrets", "allowed_authentication_secrets") {
		secretsReq, err := buildAllowedAuthSecretsRequestFromState(d)
		if err != nil {
			return diag.FromErr(err)
		}
		gitSet.WithAllowedAuthenticationSecrets(*secretsReq)
	}

	if !reflect.DeepEqual(*gitSet, *sdk.NewSetGitHttpsApiTokenBasedParamsRequest()) {
		set.WithGitHttpsApiTokenBasedParams(*gitSet)
	}
	if !reflect.DeepEqual(*set, *sdk.NewApiIntegrationSetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithSet(*set)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating git HTTPS API token-based API integration: %w", err))
		}
	}

	if !reflect.DeepEqual(*gitUnset, *sdk.NewUnsetGitHttpsApiTokenBasedParamsRequest()) {
		unset.WithGitHttpsApiTokenBasedParams(*gitUnset)
	}
	if !reflect.DeepEqual(*unset, *sdk.NewApiIntegrationUnsetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithUnset(*unset)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating git HTTPS API token-based API integration: %w", err))
		}
	}

	return ReadApiIntegrationGitRepositoryToken(ctx, d, meta)
}
