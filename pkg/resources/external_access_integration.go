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

func eaiAllowedAuthenticationSecretsSchema() map[string]*schema.Schema {
	exactlyOneOf := []string{
		"allowed_authentication_secrets.0.none",
		"allowed_authentication_secrets.0.all",
		"allowed_authentication_secrets.0.secrets",
	}
	return map[string]*schema.Schema{
		"none": {
			Type:         schema.TypeBool,
			Optional:     true,
			Description:  "When true, no secrets are allowed for authentication.",
			ExactlyOneOf: exactlyOneOf,
		},
		"all": {
			Type:         schema.TypeBool,
			Optional:     true,
			Description:  "When true, all secrets in the account are allowed for authentication.",
			ExactlyOneOf: exactlyOneOf,
		},
		"secrets": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
			},
			DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("allowed_authentication_secrets.0.secrets"),
			Description:      "Specifies the fully qualified identifiers of secrets allowed for authentication.",
			ExactlyOneOf:     exactlyOneOf,
		},
	}
}

func eaiAllowedApiAuthIntegrationsSchema() map[string]*schema.Schema {
	exactlyOneOf := []string{
		"allowed_api_authentication_integrations.0.none",
		"allowed_api_authentication_integrations.0.integrations",
	}
	return map[string]*schema.Schema{
		"none": {
			Type:         schema.TypeBool,
			Optional:     true,
			Description:  "When true, no API authentication integrations are allowed.",
			ExactlyOneOf: exactlyOneOf,
		},
		"integrations": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
			DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("allowed_api_authentication_integrations.0.integrations"),
			Description:      "Specifies the API authentication integrations allowed for authenticating to external locations.",
			ExactlyOneOf:     exactlyOneOf,
		},
	}
}

var externalAccessIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the external access integration."),
	},
	"enabled": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Specifies whether the integration is enabled.",
	},
	"allowed_network_rules": {
		Type:     schema.TypeSet,
		Required: true,
		MinItems: 1,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		},
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("allowed_network_rules"),
		Description:      relatedResourceDescription("Specifies the network rules for external locations reachable through this integration. At least one is required.", resources.NetworkRule),
	},
	"allowed_authentication_secrets": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Specifies allowed authentication secrets for this integration. Exactly one of `none`, `all`, or `secrets` must be set inside the block.",
		Elem: &schema.Resource{
			Schema: eaiAllowedAuthenticationSecretsSchema(),
		},
	},
	"allowed_api_authentication_integrations": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Specifies allowed API authentication integrations for this integration. Exactly one of `none` or `integrations` must be set inside the block.",
		Elem: &schema.Resource{
			Schema: eaiAllowedApiAuthIntegrationsSchema(),
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the external access integration.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW EXTERNAL ACCESS INTEGRATIONS` for this integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowExternalAccessIntegrationSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE EXTERNAL ACCESS INTEGRATION` for this integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeExternalAccessIntegrationDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func ExternalAccessIntegration() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ExternalAccessIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		// TODO [SNOW-3895663]: The resource is currently registered only in the test provider. When it is moved to the production
		// provider, add ExternalAccessIntegrationResource to previewfeatures and wrap the functions below with
		// PreviewFeature*ContextWrapper.
		CreateContext: TrackingCreateWrapper(resources.ExternalAccessIntegration, CreateExternalAccessIntegration),
		ReadContext:   TrackingReadWrapper(resources.ExternalAccessIntegration, ReadExternalAccessIntegrationFunc(true)),
		UpdateContext: TrackingUpdateWrapper(resources.ExternalAccessIntegration, UpdateExternalAccessIntegration),
		DeleteContext: TrackingDeleteWrapper(resources.ExternalAccessIntegration, deleteFunc),
		Description:   "Resource used to manage external access integration objects. For more information, check [external access integration documentation](https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration).",

		Schema: externalAccessIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ExternalAccessIntegration, ImportName[sdk.AccountObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: TrackingCustomDiffWrapper(resources.ExternalAccessIntegration, customdiff.All(
			ComputedIfAnyAttributeChanged(externalAccessIntegrationSchema, ShowOutputAttributeName, "enabled", "comment"),
			// For now, the block fields (allowed_authentication_secrets and allowed_api_authentication_integrations)
			// are excluded from describe_output invalidation. TypeList/TypeSet blocks produce false-positive HasChange
			// results in ComputedIfAnyAttributeChanged. describe_output is still refreshed correctly after every apply.
			// TODO [SNOW-1648997]: address the above comment (same pattern as session_policy)
			ComputedIfAnyAttributeChanged(
				externalAccessIntegrationSchema, DescribeOutputAttributeName,
				"enabled",
				"comment",
				"allowed_network_rules",
			),
			ComputedIfAnyAttributeChanged(externalAccessIntegrationSchema, FullyQualifiedNameAttributeName, "name"),
		)),
	}
}

func CreateExternalAccessIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	allowedNetworkRules, err := parseSchemaObjectIdentifierSet(d.Get("allowed_network_rules"))
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateExternalAccessIntegrationRequest(id, allowedNetworkRules, d.Get("enabled").(bool))

	errs := errors.Join(
		attributeMappedValueCreateBuilder(d, "allowed_api_authentication_integrations", request.WithAllowedApiAuthenticationIntegrations, buildEAIAllowedApiAuthIntegrationsRequest),
		attributeMappedValueCreateBuilder(d, "allowed_authentication_secrets", request.WithAllowedAuthenticationSecrets, buildEAIAllowedAuthSecretsRequest),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.ExternalAccessIntegrations.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating external access integration %v, err = %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadExternalAccessIntegrationFunc(false)(ctx, d, meta)
}

func ReadExternalAccessIntegrationFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		externalAccessIntegration, err := client.ExternalAccessIntegrations.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query external access integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("External access integration id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		details, err := client.ExternalAccessIntegrations.DescribeDetails(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			normalizeStringList := func(v any) any {
				if list, ok := v.([]any); ok {
					return expandStringList(list)
				}
				return v
			}
			if err = handleExternalChangesToObjectInFlatDescribeDeepEqual(
				d,
				outputMapping{"allowed_authentication_secrets", "allowed_authentication_secrets", details.AllowedAuthenticationSecrets, eaiAuthSecretsResourceStateFromDescribeOutput(details.AllowedAuthenticationSecrets), normalizeStringList},
				outputMapping{"allowed_api_authentication_integrations", "allowed_api_authentication_integrations", details.AllowedApiAuthenticationIntegrations, eaiApiAuthIntegrationsResourceStateFromDescribeOutput(details.AllowedApiAuthenticationIntegrations), normalizeStringList},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			// not reading allowed_authentication_secrets and allowed_api_authentication_integrations on purpose
			// (handled as external change to describe output; on import/first read set directly above)
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("name", details.Id.Name()),
			d.Set("enabled", details.Enabled),
			d.Set("comment", details.Comment),
			d.Set("allowed_network_rules", collections.Map(details.AllowedNetworkRules, sdk.SchemaObjectIdentifier.FullyQualifiedName)),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.ExternalAccessIntegrationToSchema(externalAccessIntegration)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ExternalAccessIntegrationDetailsToSchema(*details)}),
		)

		return diag.FromErr(errs)
	}
}

func eaiAuthSecretsResourceStateFromDescribeOutput(secrets []string) []any {
	if len(secrets) == 0 {
		return []any{map[string]any{"none": true}}
	}
	if len(secrets) == 1 && strings.EqualFold(secrets[0], "ALL") {
		return []any{map[string]any{"all": true}}
	}
	elems := make([]any, len(secrets))
	for i, s := range secrets {
		elems[i] = s
	}
	return []any{map[string]any{
		"secrets": schema.NewSet(schema.HashString, elems),
	}}
}

func eaiApiAuthIntegrationsResourceStateFromDescribeOutput(integrations []string) []any {
	if len(integrations) == 0 {
		return []any{map[string]any{"none": true}}
	}
	elems := make([]any, len(integrations))
	for i, s := range integrations {
		elems[i] = s
	}
	return []any{map[string]any{
		"integrations": schema.NewSet(schema.HashString, elems),
	}}
}

func UpdateExternalAccessIntegration(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	setRequest := sdk.NewExternalAccessIntegrationSetRequest()
	unsetRequest := sdk.NewExternalAccessIntegrationUnsetRequest()

	// ALLOWED_NETWORK_RULES is required and must not be empty, so it is only ever set, never unset.
	if d.HasChange("allowed_network_rules") {
		allowedNetworkRules, err := parseSchemaObjectIdentifierSet(d.Get("allowed_network_rules"))
		if err != nil {
			return diag.FromErr(err)
		}
		setRequest.WithAllowedNetworkRules(allowedNetworkRules)
	}

	if errs := errors.Join(
		// ENABLED is required, so it is only ever set, never unset.
		booleanAttributeUpdateSetOnly(d, "enabled", &setRequest.Enabled),
		stringAttributeUpdate(d, "comment", &setRequest.Comment, &unsetRequest.Comment),
		attributeMappedValueUpdate(d, "allowed_api_authentication_integrations", &setRequest.AllowedApiAuthenticationIntegrations, &unsetRequest.AllowedApiAuthenticationIntegrations, buildEAIAllowedApiAuthIntegrationsRequest),
		attributeMappedValueUpdate(d, "allowed_authentication_secrets", &setRequest.AllowedAuthenticationSecrets, &unsetRequest.AllowedAuthenticationSecrets, buildEAIAllowedAuthSecretsRequest),
	); errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(*setRequest, *sdk.NewExternalAccessIntegrationSetRequest()) {
		if err := client.ExternalAccessIntegrations.Alter(ctx, sdk.NewAlterExternalAccessIntegrationRequest(id).WithSet(*setRequest)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting properties for external access integration %v, err = %w", id.FullyQualifiedName(), err))
		}
	}

	if !reflect.DeepEqual(*unsetRequest, *sdk.NewExternalAccessIntegrationUnsetRequest()) {
		if err := client.ExternalAccessIntegrations.Alter(ctx, sdk.NewAlterExternalAccessIntegrationRequest(id).WithUnset(*unsetRequest)); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting properties for external access integration %v, err = %w", id.FullyQualifiedName(), err))
		}
	}

	return ReadExternalAccessIntegrationFunc(false)(ctx, d, meta)
}

func buildEAIAllowedApiAuthIntegrationsRequest(v any) (sdk.ExternalAccessIntegrationAllowedApiAuthenticationIntegrationsRequest, error) {
	list := v.([]any)
	if len(list) == 0 {
		return sdk.ExternalAccessIntegrationAllowedApiAuthenticationIntegrationsRequest{}, fmt.Errorf("allowed_api_authentication_integrations block is empty")
	}
	block := list[0].(map[string]any)
	req := sdk.NewExternalAccessIntegrationAllowedApiAuthenticationIntegrationsRequest()
	if none, _ := block["none"].(bool); none {
		return *req.WithNone(true), nil
	}
	if intSet, ok := block["integrations"].(*schema.Set); ok && intSet.Len() > 0 {
		ids, err := collections.MapErr(expandStringList(intSet.List()), sdk.ParseAccountObjectIdentifier)
		if err != nil {
			return sdk.ExternalAccessIntegrationAllowedApiAuthenticationIntegrationsRequest{}, err
		}
		return *req.WithIntegrations(ids), nil
	}
	return sdk.ExternalAccessIntegrationAllowedApiAuthenticationIntegrationsRequest{}, fmt.Errorf("allowed_api_authentication_integrations block has no recognized field set")
}

func buildEAIAllowedAuthSecretsRequest(v any) (sdk.ExternalAccessIntegrationAllowedAuthenticationSecretsRequest, error) {
	list := v.([]any)
	if len(list) == 0 {
		return sdk.ExternalAccessIntegrationAllowedAuthenticationSecretsRequest{}, fmt.Errorf("allowed_authentication_secrets block is empty")
	}
	block := list[0].(map[string]any)
	req := sdk.NewExternalAccessIntegrationAllowedAuthenticationSecretsRequest()
	if all, _ := block["all"].(bool); all {
		return *req.WithAll(true), nil
	}
	if none, _ := block["none"].(bool); none {
		return *req.WithNone(true), nil
	}
	if secretSet, ok := block["secrets"].(*schema.Set); ok && secretSet.Len() > 0 {
		ids, err := parseSchemaObjectIdentifierSet(secretSet)
		if err != nil {
			return sdk.ExternalAccessIntegrationAllowedAuthenticationSecretsRequest{}, err
		}
		return *req.WithSecrets(ids), nil
	}
	return sdk.ExternalAccessIntegrationAllowedAuthenticationSecretsRequest{}, fmt.Errorf("allowed_authentication_secrets block has no recognized field set")
}
