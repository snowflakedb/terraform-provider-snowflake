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

var apiIntegrationGitRepositoryTokenSchema = func() map[string]*schema.Schema {
	apiIntegrationGitRepositoryToken := map[string]*schema.Schema{
		"allowed_authentication_secrets": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Specifies which authentication secrets are allowed to be used when authenticating to the git repository. Valid values are: ALL, NONE, or a comma-separated list of fully-qualified secret identifiers (e.g. \"mydb.myschema.mysecret\").",
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
			ComputedIfAnyAttributeChanged(apiIntegrationGitRepositoryTokenSchema, DescribeOutputAttributeName, "enabled", "api_allowed_prefixes", "api_blocked_prefixes", "comment", "allowed_authentication_secrets"),
		),
	}
}

func buildAllowedAuthenticationSecretsRequest(value string) (*sdk.ApiIntegrationAllowedAuthenticationSecretsRequest, error) {
	req := sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest()
	upper := strings.ToUpper(strings.TrimSpace(value))
	switch upper {
	case "ALL":
		return req.WithAllSecrets(true), nil
	case "NONE":
		return req.WithNoSecrets(true), nil
	default:
		parts := strings.Split(value, ",")
		identifiers := make([]sdk.SchemaObjectIdentifier, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			id, err := sdk.ParseSchemaObjectIdentifier(p)
			if err != nil {
				return nil, fmt.Errorf("invalid secret identifier %q in allowed_authentication_secrets: %w", p, err)
			}
			identifiers = append(identifiers, id)
		}
		if len(identifiers) == 0 {
			return nil, fmt.Errorf("allowed_authentication_secrets must be ALL, NONE, or a non-empty comma-separated list of secret identifiers")
		}
		return req.WithAllowedList(identifiers), nil
	}
}

func CreateApiIntegrationGitRepositoryToken(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, request, err := handleApiIntegrationCommonCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	gitParams := sdk.NewGitHttpsApiTokenBasedParamsRequest()

	if v := d.Get("allowed_authentication_secrets").(string); v != "" {
		secretsReq, err := buildAllowedAuthenticationSecretsRequest(v)
		if err != nil {
			return diag.FromErr(err)
		}
		gitParams.WithAllowedAuthenticationSecrets(*secretsReq)
	}

	if err = client.ApiIntegrations.Create(ctx, request.WithGitHttpsApiTokenBasedProviderParams(*gitParams)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating git HTTPS API token-based API integration: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadApiIntegrationGitRepositoryToken(ctx, d, meta)
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

	// Validate this is a token-based sub-type: no UserAuthType and no UsePrivatelinkEndpoint.
	if details.UserAuthType != "" {
		return nil, fmt.Errorf(
			"api integration %s has user_auth_type %q, not compatible with snowflake_api_integration_git_repository_token; use the appropriate resource type",
			id.FullyQualifiedName(),
			details.UserAuthType,
		)
	}
	if details.UsePrivatelinkEndpoint {
		return nil, fmt.Errorf(
			"api integration %s has use_privatelink_endpoint=true, not compatible with snowflake_api_integration_git_repository_token; use the appropriate resource type",
			id.FullyQualifiedName(),
		)
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
		d.Set("allowed_authentication_secrets", gitDetails.AllowedAuthenticationSecrets),
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

	if d.HasChange("allowed_authentication_secrets") {
		v := d.Get("allowed_authentication_secrets").(string)
		if v == "" {
			gitUnset.WithAllowedAuthenticationSecrets(true)
		} else {
			secretsReq, err := buildAllowedAuthenticationSecretsRequest(v)
			if err != nil {
				return diag.FromErr(err)
			}
			gitSet.WithAllowedAuthenticationSecrets(*secretsReq)
		}
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
