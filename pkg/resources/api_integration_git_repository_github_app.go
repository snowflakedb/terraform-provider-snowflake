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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var apiIntegrationGitRepositoryGithubAppSchema = func() map[string]*schema.Schema {
	apiIntegrationGitRepositoryGithubApp := map[string]*schema.Schema{
		// api_provider intentionally omitted: git_https_api is the only valid value (gotcha #9).
		// The value is visible in describe_output.
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE API INTEGRATION` for the given integration.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeGitRepositoryGithubAppApiIntegrationSchema,
			},
		},
	}
	return collections.MergeMaps(apiIntegrationCommonSchema, apiIntegrationGitRepositoryGithubApp)
}()

func ApiIntegrationGitRepositoryGithubApp() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ApiIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ApiIntegrationGitRepositoryGithubApp, CreateApiIntegrationGitRepositoryGithubApp),
		ReadContext:   TrackingReadWrapper(resources.ApiIntegrationGitRepositoryGithubApp, ReadApiIntegrationGitRepositoryGithubApp),
		UpdateContext: TrackingUpdateWrapper(resources.ApiIntegrationGitRepositoryGithubApp, UpdateApiIntegrationGitRepositoryGithubApp),
		DeleteContext: TrackingDeleteWrapper(resources.ApiIntegrationGitRepositoryGithubApp, deleteFunc),
		Description:   "Resource used to manage API integration for git repositories using GitHub App authentication. For more information, check [api integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-api-integration).",

		Schema: apiIntegrationGitRepositoryGithubAppSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ApiIntegrationGitRepositoryGithubApp, ImportApiIntegrationGitRepositoryGithubApp),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(apiIntegrationGitRepositoryGithubAppSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(apiIntegrationGitRepositoryGithubAppSchema, DescribeOutputAttributeName, "enabled", "api_allowed_prefixes", "api_blocked_prefixes", "comment"),
		),
	}
}

func CreateApiIntegrationGitRepositoryGithubApp(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, request, err := handleApiIntegrationCommonCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = client.ApiIntegrations.Create(ctx, request.WithGitHttpsApiGithubAppProviderParams(*sdk.NewGitHttpsApiGithubAppParamsRequest())); err != nil {
		return diag.FromErr(fmt.Errorf("error creating git repository GitHub App API integration: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadApiIntegrationGitRepositoryGithubApp(ctx, d, meta)
}

func ImportApiIntegrationGitRepositoryGithubApp(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	details, err := client.ApiIntegrations.DescribeGitHttpsApiDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not describe API integration %s during import: %w", id.FullyQualifiedName(), err)
	}

	if details.UserAuthType != string(sdk.ApiIntegrationUserAuthTypeSnowflakeGithubApp) {
		return nil, fmt.Errorf(
			"api integration %s has user_auth_type %q, not compatible with snowflake_api_integration_git_repository_github_app (expected %q); use the appropriate resource type",
			id.FullyQualifiedName(),
			details.UserAuthType,
			sdk.ApiIntegrationUserAuthTypeSnowflakeGithubApp,
		)
	}

	return ImportName[sdk.AccountObjectIdentifier](ctx, d, meta)
}

func ReadApiIntegrationGitRepositoryGithubApp(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
					Summary:  "Failed to query API integration git repository GitHub App. Marking the resource as removed.",
					Detail:   fmt.Sprintf("API integration git repository GitHub App id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	gitDetails, err := client.ApiIntegrations.DescribeGitHttpsApiDetails(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe API integration git repository GitHub App (%s): %w", d.Id(), err))
	}

	errs := errors.Join(
		handleApiIntegrationCommonRead(d, id, s, gitDetails.AllowedPrefixes, gitDetails.BlockedPrefixes),
		// api_provider intentionally not set: git_https_api is the only valid value and not exposed in top-level schema.
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ApiIntegrationGitRepositoryGithubAppDetailsToSchema(gitDetails)}),
	)
	return diag.FromErr(errs)
}

func UpdateApiIntegrationGitRepositoryGithubApp(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewApiIntegrationSetRequest()
	unset := sdk.NewApiIntegrationUnsetRequest()

	if err = handleApiIntegrationCommonUpdate(d, set, unset); err != nil {
		return diag.FromErr(err)
	}

	if !reflect.DeepEqual(*set, *sdk.NewApiIntegrationSetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithSet(*set)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating git repository GitHub App API integration: %w", err))
		}
	}

	if !reflect.DeepEqual(*unset, *sdk.NewApiIntegrationUnsetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithUnset(*unset)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating git repository GitHub App API integration: %w", err))
		}
	}

	return ReadApiIntegrationGitRepositoryGithubApp(ctx, d, meta)
}
