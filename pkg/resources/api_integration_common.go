package resources

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
