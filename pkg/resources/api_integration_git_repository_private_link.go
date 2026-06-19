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

var apiIntegrationGitRepositoryPrivateLinkSchema = func() map[string]*schema.Schema {
	apiIntegrationGitRepositoryPrivateLink := map[string]*schema.Schema{
		"all_allowed_authentication_secrets": {
			Type:          schema.TypeBool,
			Optional:      true,
			ConflictsWith: []string{"no_allowed_authentication_secrets", "allowed_authentication_secrets"},
			Description:   "When set to true, all authentication secrets are allowed to be used when authenticating to the git repository. Conflicts with `no_allowed_authentication_secrets` and `allowed_authentication_secrets`.",
		},
		"no_allowed_authentication_secrets": {
			Type:          schema.TypeBool,
			Optional:      true,
			ConflictsWith: []string{"all_allowed_authentication_secrets", "allowed_authentication_secrets"},
			Description:   "When set to true, no authentication secrets are allowed to be used when authenticating to the git repository. Conflicts with `all_allowed_authentication_secrets` and `allowed_authentication_secrets`.",
		},
		"allowed_authentication_secrets": {
			Type:             schema.TypeSet,
			Optional:         true,
			ConflictsWith:    []string{"all_allowed_authentication_secrets", "no_allowed_authentication_secrets"},
			DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("allowed_authentication_secrets"),
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			},
			Description: "A list of fully-qualified secret identifiers (database.schema.secret) allowed to be used when authenticating to the git repository. Conflicts with `all_allowed_authentication_secrets` and `no_allowed_authentication_secrets`.",
		},
		"use_privatelink_endpoint": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Specifies whether to use the private link endpoint for the git repository. When set to true, Snowflake uses the VNet-injected endpoint for the git repository.",
		},
		"tls_trusted_certificates": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Specifies secrets containing self-signed certificates to be used when authenticating with a Git repository server over private link. Only needed when the certificate is self-signed rather than signed by a certificate authority. Each entry must be a fully-qualified name of a Snowflake secret of type generic string whose value is Base64-encoded certificate data.",
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE API INTEGRATION` for the given integration.",
			Elem: &schema.Resource{
				Schema: schemas.DescribeGitRepositoryPrivateLinkApiIntegrationSchema,
			},
		},
	}
	return collections.MergeMaps(apiIntegrationCommonSchema, apiIntegrationGitRepositoryPrivateLink)
}()

func ApiIntegrationGitRepositoryPrivateLink() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ApiIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ApiIntegrationGitRepositoryPrivateLink, CreateApiIntegrationGitRepositoryPrivateLink),
		ReadContext:   TrackingReadWrapper(resources.ApiIntegrationGitRepositoryPrivateLink, ReadApiIntegrationGitRepositoryPrivateLink),
		UpdateContext: TrackingUpdateWrapper(resources.ApiIntegrationGitRepositoryPrivateLink, UpdateApiIntegrationGitRepositoryPrivateLink),
		DeleteContext: TrackingDeleteWrapper(resources.ApiIntegrationGitRepositoryPrivateLink, deleteFunc),
		Description:   "Resource used to manage API integration for git HTTPS API via a private link endpoint. For more information, check [api integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-api-integration).",

		Schema: apiIntegrationGitRepositoryPrivateLinkSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ApiIntegrationGitRepositoryPrivateLink, ImportApiIntegrationGitRepositoryPrivateLink),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(apiIntegrationGitRepositoryPrivateLinkSchema, ShowOutputAttributeName, "enabled", "comment"),
			ComputedIfAnyAttributeChanged(apiIntegrationGitRepositoryPrivateLinkSchema, DescribeOutputAttributeName, "enabled", "api_allowed_prefixes", "api_blocked_prefixes", "comment", "all_allowed_authentication_secrets", "no_allowed_authentication_secrets", "allowed_authentication_secrets", "use_privatelink_endpoint", "tls_trusted_certificates"),
		),
	}
}

func buildAllowedAuthSecretsRequestFromStatePrivateLink(d *schema.ResourceData) (*sdk.ApiIntegrationAllowedAuthenticationSecretsRequest, error) {
	req := sdk.NewApiIntegrationAllowedAuthenticationSecretsRequest()
	if v, ok := d.GetOk("all_allowed_authentication_secrets"); ok && v.(bool) {
		return req.WithAllSecrets(true), nil
	}
	if v, ok := d.GetOk("no_allowed_authentication_secrets"); ok && v.(bool) {
		return req.WithNoSecrets(true), nil
	}
	if v, ok := d.GetOk("allowed_authentication_secrets"); ok && v.(*schema.Set).Len() > 0 {
		ids, err := collections.MapErr(expandStringList(v.(*schema.Set).List()), sdk.ParseSchemaObjectIdentifier)
		if err != nil {
			return nil, err
		}
		return req.WithAllowedList(ids), nil
	}
	return nil, nil
}

func setAllowedAuthSecretFieldsFromDescribePrivateLink(d *schema.ResourceData, raw string) error {
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
	case "":
		return errors.Join(
			d.Set("all_allowed_authentication_secrets", false),
			d.Set("no_allowed_authentication_secrets", false),
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
			d.Set("allowed_authentication_secrets", ids),
		)
	}
}

func buildTlsTrustedCertificates(d *schema.ResourceData) ([]sdk.SchemaObjectIdentifier, error) {
	raw := d.Get("tls_trusted_certificates").([]any)
	ids, err := collections.MapErr(raw, func(v any) (sdk.SchemaObjectIdentifier, error) {
		return sdk.ParseSchemaObjectIdentifier(v.(string))
	})
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func CreateApiIntegrationGitRepositoryPrivateLink(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, request, err := handleApiIntegrationCommonCreate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	gitParams := sdk.NewGitHttpsApiPrivateLinkParamsRequest(d.Get("use_privatelink_endpoint").(bool))

	secretsReq, err := buildAllowedAuthSecretsRequestFromStatePrivateLink(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if secretsReq != nil {
		gitParams.WithAllowedAuthenticationSecrets(*secretsReq)
	}

	tlsCerts, err := buildTlsTrustedCertificates(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(tlsCerts) > 0 {
		gitParams.WithTlsTrustedCertificates(tlsCerts)
	}

	if err = client.ApiIntegrations.Create(ctx, request.WithGitHttpsApiPrivateLinkProviderParams(*gitParams)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating git HTTPS API private link API integration: %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadApiIntegrationGitRepositoryPrivateLink(ctx, d, meta)
}

func ImportApiIntegrationGitRepositoryPrivateLink(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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
			"api integration %s has api_provider %q, not compatible with snowflake_api_integration_git_repository_private_link; use the appropriate resource type",
			id.FullyQualifiedName(),
			details.ApiProvider,
		)
	}

	return ImportName[sdk.AccountObjectIdentifier](ctx, d, meta)
}

func ReadApiIntegrationGitRepositoryPrivateLink(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
					Summary:  "Failed to query API integration git HTTPS API private link. Marking the resource as removed.",
					Detail:   fmt.Sprintf("API integration git HTTPS API private link id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	gitDetails, err := client.ApiIntegrations.DescribeGitHttpsApiDetails(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe API integration git HTTPS API private link (%s): %w", d.Id(), err))
	}

	errs := errors.Join(
		handleApiIntegrationCommonRead(d, id, s, gitDetails.AllowedPrefixes, gitDetails.BlockedPrefixes),
		setAllowedAuthSecretFieldsFromDescribePrivateLink(d, gitDetails.AllowedAuthenticationSecrets),
		d.Set("use_privatelink_endpoint", gitDetails.UsePrivatelinkEndpoint),
		d.Set("tls_trusted_certificates", gitDetails.TlsTrustedCertificates),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ApiIntegrationGitRepositoryPrivateLinkDetailsToSchema(gitDetails)}),
	)
	return diag.FromErr(errs)
}

func UpdateApiIntegrationGitRepositoryPrivateLink(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewApiIntegrationSetRequest()
	unset := sdk.NewApiIntegrationUnsetRequest()
	gitSet := sdk.NewSetGitHttpsApiPrivateLinkParamsRequest()
	gitUnset := sdk.NewUnsetGitHttpsApiPrivateLinkParamsRequest()

	if err := handleApiIntegrationCommonUpdate(d, set, unset); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("all_allowed_authentication_secrets", "no_allowed_authentication_secrets", "allowed_authentication_secrets") {
		secretsReq, err := buildAllowedAuthSecretsRequestFromStatePrivateLink(d)
		if err != nil {
			return diag.FromErr(err)
		}
		if secretsReq != nil {
			gitSet.WithAllowedAuthenticationSecrets(*secretsReq)
		} else {
			gitUnset.WithAllowedAuthenticationSecrets(true)
		}
	}

	if d.HasChange("use_privatelink_endpoint") {
		v := d.Get("use_privatelink_endpoint").(bool)
		gitSet.WithUsePrivatelinkEndpoint(v)
	}

	if d.HasChange("tls_trusted_certificates") {
		tlsCerts, err := buildTlsTrustedCertificates(d)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(tlsCerts) > 0 {
			gitSet.WithTlsTrustedCertificates(tlsCerts)
		} else {
			gitUnset.WithTlsTrustedCertificates(true)
		}
	}

	if !reflect.DeepEqual(*gitSet, *sdk.NewSetGitHttpsApiPrivateLinkParamsRequest()) {
		set.WithGitHttpsApiPrivateLinkParams(*gitSet)
	}
	if !reflect.DeepEqual(*set, *sdk.NewApiIntegrationSetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithSet(*set)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating git HTTPS API private link API integration: %w", err))
		}
	}

	if !reflect.DeepEqual(*gitUnset, *sdk.NewUnsetGitHttpsApiPrivateLinkParamsRequest()) {
		unset.WithGitHttpsApiPrivateLinkParams(*gitUnset)
	}
	if !reflect.DeepEqual(*unset, *sdk.NewApiIntegrationUnsetRequest()) {
		req := sdk.NewAlterApiIntegrationRequest(id).WithUnset(*unset)
		if err = client.ApiIntegrations.Alter(ctx, req); err != nil {
			return diag.FromErr(fmt.Errorf("error updating git HTTPS API private link API integration: %w", err))
		}
	}

	return ReadApiIntegrationGitRepositoryPrivateLink(ctx, d, meta)
}
