package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var apiIntegrationAllowedAuthSecretsSchema = map[string]*schema.Schema{
	"all_allowed_authentication_secrets": {
		Type:          schema.TypeBool,
		Optional:      true,
		ConflictsWith: []string{"no_allowed_authentication_secrets", "allowed_authentication_secrets"},
		Description:   externalChangesNotDetectedFieldDescription("When set to true, all authentication secrets are allowed to be used when authenticating to the git repository. Conflicts with `no_allowed_authentication_secrets` and `allowed_authentication_secrets`."),
	},
	"no_allowed_authentication_secrets": {
		Type:          schema.TypeBool,
		Optional:      true,
		ConflictsWith: []string{"all_allowed_authentication_secrets", "allowed_authentication_secrets"},
		Description:   externalChangesNotDetectedFieldDescription("When set to true, no authentication secrets are allowed to be used when authenticating to the git repository. Conflicts with `all_allowed_authentication_secrets` and `allowed_authentication_secrets`."),
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
		Description: externalChangesNotDetectedFieldDescription("A list of fully-qualified secret identifiers (database.schema.secret) allowed to be used when authenticating to the git repository. Conflicts with `all_allowed_authentication_secrets` and `no_allowed_authentication_secrets`."),
	},
}

var apiIntegrationCommonSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier (i.e. name) for the integration. This value must be unique in your account."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"enabled": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Specifies whether this API integration is enabled or disabled.",
	},
	"api_allowed_prefixes": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		MinItems:    1,
		Description: "Explicitly limits external functions that use the integration to reference one or more HTTPS proxy service and remote service endpoints and resources.",
	},
	"api_blocked_prefixes": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Lists the endpoints and resources in the HTTPS proxy service that are not allowed to be called from Snowflake.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the integration.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW API INTEGRATIONS` for the given integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowApiIntegrationSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func handleApiIntegrationCommonCreate(d *schema.ResourceData) (sdk.AccountObjectIdentifier, *sdk.CreateApiIntegrationRequest, error) {
	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return sdk.AccountObjectIdentifier{}, nil, err
	}

	request := sdk.NewCreateApiIntegrationRequest(
		id,
		toApiIntegrationEndpointPrefix(expandStringList(d.Get("api_allowed_prefixes").([]any))),
		d.Get("enabled").(bool),
	)

	if errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "api_blocked_prefixes", request.WithApiBlockedPrefixes, func(v []any) ([]sdk.ApiIntegrationEndpointPrefix, error) {
			return collections.Map(expandStringList(v), func(item string) sdk.ApiIntegrationEndpointPrefix {
				return sdk.ApiIntegrationEndpointPrefix{Path: item}
			}), nil
		}),
		stringAttributeCreate(d, "comment", &request.Comment),
	); errs != nil {
		return sdk.AccountObjectIdentifier{}, nil, errs
	}

	return id, request, nil
}

func handleApiIntegrationCommonUpdate(d *schema.ResourceData, set *sdk.ApiIntegrationSetRequest, unset *sdk.ApiIntegrationUnsetRequest) error {
	if d.HasChange("api_allowed_prefixes") {
		set.WithApiAllowedPrefixes(toApiIntegrationEndpointPrefix(expandStringList(d.Get("api_allowed_prefixes").([]any))))
	}
	if d.HasChange("api_blocked_prefixes") {
		v := d.Get("api_blocked_prefixes").([]any)
		if len(v) > 0 {
			set.WithApiBlockedPrefixes(toApiIntegrationEndpointPrefix(expandStringList(v)))
		} else {
			unset.WithApiBlockedPrefixes(true)
		}
	}
	return errors.Join(
		booleanAttributeUpdateSetOnly(d, "enabled", &set.Enabled),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	)
}

func handleApiIntegrationCommonRead(
	d *schema.ResourceData,
	id sdk.AccountObjectIdentifier,
	integration *sdk.ApiIntegration,
	allowedPrefixes []string,
	blockedPrefixes []string,
) error {
	return errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set("enabled", integration.Enabled),
		d.Set("api_allowed_prefixes", allowedPrefixes),
		d.Set("api_blocked_prefixes", blockedPrefixes),
		d.Set("comment", integration.Comment),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.ApiIntegrationToSchema(integration)}),
	)
}

func importApiIntegrationWithDetails[D any](
	ctx context.Context,
	d *schema.ResourceData,
	meta any,
	describeFunc func(ctx context.Context, client *sdk.Client, id sdk.AccountObjectIdentifier) (D, error),
	validateFunc func(details D, id sdk.AccountObjectIdentifier) error,
) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	details, err := describeFunc(ctx, client, id)
	if err != nil {
		return nil, fmt.Errorf("could not describe API integration %s during import: %w", id.FullyQualifiedName(), err)
	}
	if err := validateFunc(details, id); err != nil {
		return nil, err
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
	if v, ok := d.GetOk("allowed_authentication_secrets"); ok && v.(*schema.Set).Len() > 0 {
		ids, err := collections.MapErr(expandStringList(v.(*schema.Set).List()), sdk.ParseSchemaObjectIdentifier)
		if err != nil {
			return nil, err
		}
		return req.WithAllowedList(ids), nil
	}
	return nil, nil
}

func setAllowedAuthSecretFieldsFromDescribe(d *schema.ResourceData, raw string) error {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case string(sdk.ApiIntegrationAllowedAuthenticationSecretsValueAll):
		return errors.Join(
			d.Set("all_allowed_authentication_secrets", true),
			d.Set("no_allowed_authentication_secrets", false),
			d.Set("allowed_authentication_secrets", []any{}),
		)
	case string(sdk.ApiIntegrationAllowedAuthenticationSecretsValueNone):
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
			d.Set("allowed_authentication_secrets", collections.Map(ids, sdk.SchemaObjectIdentifier.FullyQualifiedName)),
		)
	}
}
