package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var scimIntegrationSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("String that specifies the identifier (i.e. name) for the integration; must be unique in your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"enabled": {
		Type:             schema.TypeBool,
		Required:         true,
		Description:      "Specify whether the security integration is enabled.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("enabled"),
	},
	"scim_client": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      fmt.Sprintf("Specifies the client type for the scim integration. Valid options are: %v.", possibleValuesListed(sdk.AllScimSecurityIntegrationScimClients)),
		ValidateDiagFunc: sdkValidation(sdk.ToScimSecurityIntegrationScimClientOption),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToScimSecurityIntegrationScimClientOption),
	},
	"run_as_role": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		Description: fmt.Sprintf("Specify the SCIM role in Snowflake that owns any users and roles that are imported from the identity provider into Snowflake using SCIM."+
			" Provider assumes that the specified role is already provided. Valid options are: %v.", possibleValuesListed(sdk.AllScimSecurityIntegrationRunAsRoles)),
		ValidateDiagFunc: sdkValidation(sdk.ToScimSecurityIntegrationRunAsRoleOption),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToScimSecurityIntegrationRunAsRoleOption),
	},
	"network_policy": {
		Type:             schema.TypeString,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Optional:         true,
		Description:      relatedResourceDescription("Specifies an existing network policy that controls SCIM network traffic.", resources.NetworkPolicy),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeListValueInDescribe("network_policy")),
	},
	"sync_password": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeListValueInDescribe("sync_password"),
		Description:      booleanStringFieldDescription("Specifies whether to enable or disable the synchronization of a user password from an Okta SCIM client as part of the API request to Snowflake. This property is not supported for Azure SCIM."),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the integration.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SECURITY INTEGRATIONS` for the given security integration.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSecurityIntegrationSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SECURITY INTEGRATIONS` for the given security integration.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeScimSecurityIntegrationSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func SCIMIntegration() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.ScimSecurityIntegration, CreateContextSCIMIntegration),
		ReadContext:   TrackingReadWrapper(resources.ScimSecurityIntegration, ReadContextSCIMIntegration(true)),
		UpdateContext: TrackingUpdateWrapper(resources.ScimSecurityIntegration, UpdateContextSCIMIntegration),
		DeleteContext: TrackingDeleteWrapper(resources.ScimSecurityIntegration, DeleteSecurityIntegration),
		Description:   "Resource used to manage scim security integration objects. For more information, check [security integrations documentation](https://docs.snowflake.com/en/sql-reference/sql/create-security-integration-scim).",

		Schema: scimIntegrationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ScimSecurityIntegration, ImportScimIntegration),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ScimSecurityIntegration, customdiff.All(
			ComputedIfAnyAttributeChanged(scimIntegrationSchema, ShowOutputAttributeName, "enabled", "scim_client", "comment"),
			ComputedIfAnyAttributeChanged(scimIntegrationSchema, DescribeOutputAttributeName, "enabled", "comment", "network_policy", "run_as_role", "sync_password"),
		)),

		SchemaVersion: 2,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v092ScimIntegrationStateUpgrader,
			},
			{
				Version: 1,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v093ScimIntegrationStateUpgrader,
			},
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportScimIntegration(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	integrationProperties, err := client.SecurityIntegrations.Describe(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = d.Set("name", integration.ID().Name()); err != nil {
		return nil, err
	}
	if err = d.Set("enabled", integration.Enabled); err != nil {
		return nil, err
	}
	scimClient, err := integration.SubType()
	if err != nil {
		return nil, err
	}
	if err = d.Set("scim_client", scimClient); err != nil {
		return nil, err
	}
	if runAsRoleProperty, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "RUN_AS_ROLE" }); err == nil {
		if err = d.Set("run_as_role", runAsRoleProperty.Value); err != nil {
			return nil, err
		}
	}
	if networkPolicyProperty, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "NETWORK_POLICY" }); err == nil {
		if err = d.Set("network_policy", networkPolicyProperty.Value); err != nil {
			return nil, err
		}
	}
	if strings.EqualFold(strings.TrimSpace(scimClient), string(sdk.ScimSecurityIntegrationScimClientAzure)) {
		if err = d.Set("sync_password", BooleanDefault); err != nil {
			return nil, err
		}
	} else {
		if syncPasswordProperty, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SYNC_PASSWORD" }); err == nil {
			if err = d.Set("sync_password", syncPasswordProperty.Value); err != nil {
				return nil, err
			}
		}
	}
	if err = d.Set("comment", integration.Comment); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateContextSCIMIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	scimClient, err := sdk.ToScimSecurityIntegrationScimClientOption(d.Get("scim_client").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	runAsRole, err := sdk.ToScimSecurityIntegrationRunAsRoleOption(d.Get("run_as_role").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	req := sdk.NewCreateScimSecurityIntegrationRequest(id, scimClient, runAsRole).WithEnabled(d.Get("enabled").(bool))

	if v, ok := d.GetOk("network_policy"); ok {
		req.WithNetworkPolicy(sdk.NewAccountObjectIdentifier(v.(string)))
	}

	if v := d.Get("sync_password").(string); v != BooleanDefault {
		if scimClient := d.Get("scim_client").(string); strings.EqualFold(strings.TrimSpace(scimClient), string(sdk.ScimSecurityIntegrationScimClientAzure)) {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "field `sync_password` is not supported for scim_client = \"AZURE\"",
					Detail:   "can not CREATE scim integration with field `sync_password` for scim_client = \"AZURE\"",
				},
			}
		}
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return diag.FromErr(err)
		}

		req.WithSyncPassword(parsed)
	}

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if err := client.SecurityIntegrations.CreateScim(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSCIMIntegration(false)(ctx, d, meta)
}

func ReadContextSCIMIntegration(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		integration, err := client.SecurityIntegrations.ShowByID(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query security integration. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Security integration name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		integrationProperties, err := client.SecurityIntegrations.Describe(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query security integration properties. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Security integration name: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}
		if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
			return diag.FromErr(err)
		}

		if c := integration.Category; c != sdk.SecurityIntegrationCategory {
			return diag.FromErr(fmt.Errorf("expected %v to be a SECURITY integration, got %v", id, c))
		}

		if err := d.Set("enabled", integration.Enabled); err != nil {
			return diag.FromErr(err)
		}

		scimClient, err := integration.SubType()
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("scim_client", scimClient); err != nil {
			return diag.FromErr(err)
		}

		runAsRoleProperty, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "RUN_AS_ROLE" })
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("run_as_role", runAsRoleProperty.Value); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set("comment", integration.Comment); err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			networkPolicyProperty, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "NETWORK_POLICY" })
			if err != nil {
				return diag.FromErr(err)
			}

			if err = handleExternalChangesToObjectInDescribe(d,
				describeMapping{"network_policy", "network_policy", networkPolicyProperty.Value, networkPolicyProperty.Value, nil},
			); err != nil {
				return diag.FromErr(err)
			}

			if !strings.EqualFold(strings.TrimSpace(scimClient), string(sdk.ScimSecurityIntegrationScimClientAzure)) {
				syncPasswordProperty, err := collections.FindFirst(integrationProperties, func(property sdk.SecurityIntegrationProperty) bool { return property.Name == "SYNC_PASSWORD" })
				if err != nil {
					return diag.FromErr(err)
				}
				if err = handleExternalChangesToObjectInDescribe(d,
					describeMapping{"sync_password", "sync_password", syncPasswordProperty.Value, syncPasswordProperty.Value, nil},
				); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if err = setStateToValuesFromConfig(d, scimIntegrationSchema, []string{
			"network_policy",
			"sync_password",
		}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.SecurityIntegrationToSchema(integration)}); err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ScimSecurityIntegrationPropertiesToSchema(integrationProperties)}); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func UpdateContextSCIMIntegration(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewScimIntegrationSetRequest(), sdk.NewScimIntegrationUnsetRequest()

	if d.HasChange("enabled") {
		set.WithEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("network_policy") {
		if v := d.Get("network_policy").(string); v != "" {
			set.WithNetworkPolicy(sdk.NewAccountObjectIdentifier(v))
		} else {
			unset.WithNetworkPolicy(true)
		}
	}

	if d.HasChange("sync_password") {
		if scimClient := d.Get("scim_client").(string); strings.EqualFold(strings.TrimSpace(scimClient), string(sdk.ScimSecurityIntegrationScimClientAzure)) {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "field `sync_password` is not supported for scim_client = \"AZURE\"",
					Detail:   "can not SET and UNSET field `sync_password` for scim_client = \"AZURE\"",
				},
			}
		}
		if v := d.Get("sync_password").(string); v != BooleanDefault {
			parsed, err := strconv.ParseBool(v)
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithSyncPassword(parsed)
		} else {
			unset.WithSyncPassword(true)
		}
	}

	errs := errors.Join(
		stringAttributeUpdateSetOnly(d, "comment", &set.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if (*set != sdk.ScimIntegrationSetRequest{}) {
		if err := client.SecurityIntegrations.AlterScim(ctx, sdk.NewAlterScimSecurityIntegrationRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if (*unset != sdk.ScimIntegrationUnsetRequest{}) {
		if err := client.SecurityIntegrations.AlterScim(ctx, sdk.NewAlterScimSecurityIntegrationRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextSCIMIntegration(false)(ctx, d, meta)
}
