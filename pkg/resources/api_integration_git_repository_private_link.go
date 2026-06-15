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
		"allowed_authentication_secrets": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies which authentication secrets are allowed to be used when authenticating to the git repository over a private link endpoint. Valid values are: ALL, NONE, or a comma-separated list of fully-qualified secret identifiers (e.g. \"mydb.myschema.mysecret\").",
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
			Description: "Specifies a list of TLS certificates that are trusted for secure communication with the git repository. Each entry must be a fully-qualified Snowflake schema object identifier (e.g. db.schema.cert_name).",
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
			ComputedIfAnyAttributeChanged(apiIntegrationGitRepositoryPrivateLinkSchema, DescribeOutputAttributeName, "enabled", "api_allowed_prefixes", "api_blocked_prefixes", "comment", "allowed_authentication_secrets", "use_privatelink_endpoint", "tls_trusted_certificates"),
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

func buildTlsTrustedCertificates(d *schema.ResourceData) ([]sdk.SchemaObjectIdentifier, error) {
	raw := d.Get("tls_trusted_certificates").([]interface{})
	ids := make([]sdk.SchemaObjectIdentifier, 0, len(raw))
	for _, v := range raw {
		id, err := sdk.ParseSchemaObjectIdentifier(v.(string))
		if err != nil {
			return nil, fmt.Errorf("invalid tls_trusted_certificates identifier %q: %w", v.(string), err)
		}
		ids = append(ids, id)
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

	if v := d.Get("allowed_authentication_secrets").(string); v != "" {
		secretsReq, err := buildAllowedAuthenticationSecretsRequest(v)
		if err != nil {
			return diag.FromErr(err)
		}
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

	// Validate this is a private-link sub-type: UsePrivatelinkEndpoint=true, no UserAuthType.
	if !details.UsePrivatelinkEndpoint {
		return nil, fmt.Errorf(
			"api integration %s has use_privatelink_endpoint=false, not compatible with snowflake_api_integration_git_repository_private_link; use the appropriate resource type",
			id.FullyQualifiedName(),
		)
	}
	if details.UserAuthType != "" {
		return nil, fmt.Errorf(
			"api integration %s has user_auth_type %q, not compatible with snowflake_api_integration_git_repository_private_link; use the appropriate resource type",
			id.FullyQualifiedName(),
			details.UserAuthType,
		)
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
		d.Set("allowed_authentication_secrets", gitDetails.AllowedAuthenticationSecrets),
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
