package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var externalAccessIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Specifies the name of the external access integration. This name follows the rules for Object Identifiers. The name should be unique among integrations in your account.",
	},
	"allowed_network_rules": {
		Type:        schema.TypeSet,
		Required:    true,
		MinItems:    1,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Specifies the network rules that can be used in the integration. Each value is a fully-qualified network rule identifier in the form `database.schema.name`.",
	},
	"allowed_authentication_secrets": {
		Type:        schema.TypeSet,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Specifies the secrets that can be used when making requests to the external network location. Each value is a fully-qualified secret identifier in the form `database.schema.name`.",
	},
	"enabled": {
		Type:        schema.TypeBool,
		Required:    true,
		Description: "Specifies whether this external access integration is enabled or disabled.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the integration.",
	},
	"created_on": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Date and time when the external access integration was created.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func ExternalAccessIntegration() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.AccountObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ExternalAccessIntegrations.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ExternalAccessIntegrationResource), TrackingCreateWrapper(resources.ExternalAccessIntegration, CreateExternalAccessIntegration)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ExternalAccessIntegrationResource), TrackingReadWrapper(resources.ExternalAccessIntegration, ReadExternalAccessIntegration)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ExternalAccessIntegrationResource), TrackingUpdateWrapper(resources.ExternalAccessIntegration, UpdateExternalAccessIntegration)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ExternalAccessIntegrationResource), TrackingDeleteWrapper(resources.ExternalAccessIntegration, deleteFunc)),

		Schema: externalAccessIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// parseSchemaObjectIdentifiers converts a list of fully-qualified identifier strings into SDK SchemaObjectIdentifiers.
func parseSchemaObjectIdentifiers(raw []string) ([]sdk.SchemaObjectIdentifier, error) {
	ids := make([]sdk.SchemaObjectIdentifier, len(raw))
	for i, s := range raw {
		id, err := sdk.ParseSchemaObjectIdentifier(s)
		if err != nil {
			return nil, fmt.Errorf("invalid identifier %q: %w", s, err)
		}
		ids[i] = id
	}
	return ids, nil
}

// schemaObjectIdentifiersToStrings converts a slice of SchemaObjectIdentifiers to their fully-qualified string forms.
func schemaObjectIdentifiersToStrings(ids []sdk.SchemaObjectIdentifier) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = id.FullyQualifiedName()
	}
	return result
}

func CreateExternalAccessIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	enabled := d.Get("enabled").(bool)

	networkRulesRaw := expandStringList(d.Get("allowed_network_rules").(*schema.Set).List())
	networkRules, err := parseSchemaObjectIdentifiers(networkRulesRaw)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error parsing allowed_network_rules: %w", err))
	}

	createRequest := sdk.NewCreateExternalAccessIntegrationRequest(id, networkRules, enabled)

	if v, ok := d.GetOk("allowed_authentication_secrets"); ok {
		secretsRaw := expandStringList(v.(*schema.Set).List())
		secrets, err := parseSchemaObjectIdentifiers(secretsRaw)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error parsing allowed_authentication_secrets: %w", err))
		}
		createRequest.WithAllowedAuthenticationSecrets(secrets)
	}

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(v.(string))
	}

	if err := client.ExternalAccessIntegrations.Create(ctx, createRequest); err != nil {
		return diag.FromErr(fmt.Errorf("error creating external access integration: %w", err))
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadExternalAccessIntegration(ctx, d, meta)
}

func ReadExternalAccessIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.AccountObjectIdentifier)

	integration, err := client.ExternalAccessIntegrations.ShowByIDSafely(ctx, id)
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

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", integration.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", integration.Enabled); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", integration.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_on", integration.CreatedOn.String()); err != nil {
		return diag.FromErr(err)
	}

	// Fetch additional properties from DESCRIBE INTEGRATION
	integrationProperties, err := client.ExternalAccessIntegrations.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not describe external access integration: %w", err))
	}

	for _, property := range integrationProperties {
		name := property.Name
		value := property.Value
		switch name {
		case "ENABLED":
			// Already set from SHOW above.
		case "ALLOWED_NETWORK_RULES":
			if value != "" {
				rules := strings.Split(value, ",")
				for i, r := range rules {
					rules[i] = strings.TrimSpace(r)
				}
				if err := d.Set("allowed_network_rules", rules); err != nil {
					return diag.FromErr(err)
				}
			}
		case "ALLOWED_AUTHENTICATION_SECRETS":
			if value != "" {
				secrets := strings.Split(value, ",")
				for i, s := range secrets {
					secrets[i] = strings.TrimSpace(s)
				}
				if err := d.Set("allowed_authentication_secrets", secrets); err != nil {
					return diag.FromErr(err)
				}
			}
		case "COMMENT":
			if err := d.Set("comment", value); err != nil {
				return diag.FromErr(err)
			}
		default:
			log.Printf("[WARN] unexpected external access integration property %v returned from Snowflake", name)
		}
	}

	return nil
}

func UpdateExternalAccessIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.AccountObjectIdentifier)

	var runSetStatement bool
	setRequest := sdk.NewExternalAccessIntegrationSetRequest()

	if d.HasChange("allowed_network_rules") {
		runSetStatement = true
		networkRulesRaw := expandStringList(d.Get("allowed_network_rules").(*schema.Set).List())
		networkRules, err := parseSchemaObjectIdentifiers(networkRulesRaw)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error parsing allowed_network_rules: %w", err))
		}
		setRequest.WithAllowedNetworkRules(networkRules)
	}

	if d.HasChange("allowed_authentication_secrets") {
		v := d.Get("allowed_authentication_secrets").(*schema.Set).List()
		if len(v) == 0 {
			unsetReq := sdk.NewExternalAccessIntegrationUnsetRequest().WithAllowedAuthenticationSecrets(true)
			if err := client.ExternalAccessIntegrations.Alter(ctx, sdk.NewAlterExternalAccessIntegrationRequest(id).WithUnset(*unsetReq)); err != nil {
				return diag.FromErr(fmt.Errorf("error unsetting allowed_authentication_secrets: %w", err))
			}
		} else {
			runSetStatement = true
			secretsRaw := expandStringList(v)
			secrets, err := parseSchemaObjectIdentifiers(secretsRaw)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error parsing allowed_authentication_secrets: %w", err))
			}
			setRequest.WithAllowedAuthenticationSecrets(secrets)
		}
	}

	if d.HasChange("enabled") {
		runSetStatement = true
		setRequest.WithEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("comment") {
		v := d.Get("comment").(string)
		if v == "" {
			unsetReq := sdk.NewExternalAccessIntegrationUnsetRequest().WithComment(true)
			if err := client.ExternalAccessIntegrations.Alter(ctx, sdk.NewAlterExternalAccessIntegrationRequest(id).WithUnset(*unsetReq)); err != nil {
				return diag.FromErr(fmt.Errorf("error unsetting comment: %w", err))
			}
		} else {
			runSetStatement = true
			setRequest.WithComment(v)
		}
	}

	if runSetStatement {
		if err := client.ExternalAccessIntegrations.Alter(ctx, sdk.NewAlterExternalAccessIntegrationRequest(id).WithSet(*setRequest)); err != nil {
			return diag.FromErr(fmt.Errorf("error updating external access integration: %w", err))
		}
	}

	return ReadExternalAccessIntegration(ctx, d, meta)
}
