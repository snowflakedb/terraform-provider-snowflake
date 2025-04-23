package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var networkRuleSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the network rule; must be unique for the database and schema in which the network rule is created.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the network rule.",
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the network rule.",
	},
	"type": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Specifies the type of network identifiers being allowed or blocked. A network rule can have only one type. Allowed values are IPV4, AWSVPCEID, AZURELINKID and HOST_PORT; allowed values are determined by the mode of the network rule; see https://docs.snowflake.com/en/sql-reference/sql/create-network-rule#required-parameters for details.",
		ValidateFunc: validation.StringInSlice([]string{"IPV4", "AWSVPCEID", "AZURELINKID", "HOST_PORT"}, false),
	},
	"value_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Specifies the network identifiers that will be allowed or blocked. Valid values in the list are determined by the type of network rule, see https://docs.snowflake.com/en/sql-reference/sql/create-network-rule#required-parameters for details.",
	},
	"mode": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		Description:  "Specifies what is restricted by the network rule. Valid values are INGRESS, INTERNAL_STAGE and EGRESS; see https://docs.snowflake.com/en/sql-reference/sql/create-network-rule#required-parameters for details.",
		ValidateFunc: validation.StringInSlice([]string{"INGRESS", "INTERNAL_STAGE", "EGRESS"}, false),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the network rule.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// NetworkRule returns a pointer to the resource representing a network rule.
func NetworkRule() *schema.Resource {
	// TODO(SNOW-1818849): unassign network rules before dropping
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErr[sdk.SchemaObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.NetworkRules.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.NetworkRuleResource), TrackingCreateWrapper(resources.NetworkRule, CreateContextNetworkRule)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.NetworkRuleResource), TrackingReadWrapper(resources.NetworkRule, ReadContextNetworkRule)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.NetworkRuleResource), TrackingUpdateWrapper(resources.NetworkRule, UpdateContextNetworkRule)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.NetworkRuleResource), TrackingDeleteWrapper(resources.NetworkRule, deleteFunc)),

		Schema: networkRuleSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateContextNetworkRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	ruleType := sdk.NetworkRuleType(d.Get("type").(string))
	ruleMode := sdk.NetworkRuleMode(d.Get("mode").(string))

	valueList := expandStringList(d.Get("value_list").(*schema.Set).List())
	networkRuleValues := make([]sdk.NetworkRuleValue, len(valueList))
	for i, v := range valueList {
		networkRuleValues[i] = sdk.NetworkRuleValue{Value: v}
	}

	req := sdk.NewCreateNetworkRuleRequest(
		id,
		ruleType,
		networkRuleValues,
		ruleMode,
	)

	// Set optionals
	if v, ok := d.GetOk("comment"); ok {
		req = req.WithComment(sdk.String(v.(string)))
	}

	client := meta.(*provider.Context).Client
	if err := client.NetworkRules.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadContextNetworkRule(ctx, d, meta)
}

func ReadContextNetworkRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	networkRule, err := client.NetworkRules.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query network rule. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Network rule id: %s, Err: %s", d.Id(), err),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve network rule",
				Detail:   fmt.Sprintf("Id: %s\nError: %s", d.Id(), err),
			},
		}
	}

	networkRuleDescriptions, err := client.NetworkRules.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("name", networkRule.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("database", networkRule.DatabaseName); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("schema", networkRule.SchemaName); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("type", networkRule.Type); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("value_list", networkRuleDescriptions.ValueList); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("mode", networkRule.Mode); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("comment", networkRule.Comment); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func UpdateContextNetworkRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)

	valueList := expandStringList(d.Get("value_list").(*schema.Set).List())
	networkRuleValues := make([]sdk.NetworkRuleValue, len(valueList))
	for i, v := range valueList {
		networkRuleValues[i] = sdk.NetworkRuleValue{Value: v}
	}
	comment := d.Get("comment").(string)

	if d.HasChange("value_list") {
		baseReq := sdk.NewAlterNetworkRuleRequest(id)
		if len(valueList) == 0 {
			unsetReq := sdk.NewNetworkRuleUnsetRequest().WithValueList(sdk.Bool(true))
			baseReq.WithUnset(unsetReq)
		} else {
			setReq := sdk.NewNetworkRuleSetRequest(networkRuleValues)
			baseReq.WithSet(setReq)
		}

		if err := client.NetworkRules.Alter(ctx, baseReq); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("comment") {
		baseReq := sdk.NewAlterNetworkRuleRequest(id)
		if len(comment) == 0 {
			unsetReq := sdk.NewNetworkRuleUnsetRequest().WithComment(sdk.Bool(true))
			baseReq.WithUnset(unsetReq)
		} else {
			setReq := sdk.NewNetworkRuleSetRequest(networkRuleValues).WithComment(sdk.String(comment))
			baseReq.WithSet(setReq)
		}

		if err := client.NetworkRules.Alter(ctx, baseReq); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextNetworkRule(ctx, d, meta)
}
