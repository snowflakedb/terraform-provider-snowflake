package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
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
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies the type of network identifiers being allowed or blocked. A network rule can have only one type. Allowed values are determined by the mode of the network rule; see https://docs.snowflake.com/en/sql-reference/sql/create-network-rule#required-parameters for details. " + enumValuesDescription(sdk.AllNetworkRuleTypes),
		ValidateDiagFunc: sdkValidation(sdk.ToNetworkRuleType),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToNetworkRuleType),
	},
	"value_list": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Required:    true,
		Description: "Specifies the network identifiers that will be allowed or blocked. Valid values in the list are determined by the type of network rule, see https://docs.snowflake.com/en/sql-reference/sql/create-network-rule#required-parameters for details.",
	},
	"mode": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies what is restricted by the network rule, see https://docs.snowflake.com/en/sql-reference/sql/create-network-rule#required-parameters for details. " + enumValuesDescription(sdk.AllNetworkRuleModes),
		ValidateDiagFunc: sdkValidation(sdk.ToNetworkRuleMode),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToNetworkRuleMode),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the network rule.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW NETWORK RULES` for the given network rule.",
		Elem: &schema.Resource{
			Schema: schemas.ShowNetworkRuleSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE NETWORK RULE` for the given network rule.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeNetworkRuleSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// NetworkRule returns a pointer to the resource representing a network rule.
func NetworkRule() *schema.Resource {
	// TODO(SNOW-1818849): unassign network rules before dropping
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
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
			StateContext: TrackingImportWrapper(resources.NetworkRule, ImportName[sdk.SchemaObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    cty.EmptyObject,
				Upgrade: migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName,
			},
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.NetworkRule, customdiff.All(
			ComputedIfAnyAttributeChanged(networkRuleSchema, ShowOutputAttributeName, "comment", "value_list"),
			ComputedIfAnyAttributeChanged(networkRuleSchema, DescribeOutputAttributeName, "comment", "value_list"),
			ComputedIfAnyAttributeChanged(networkRuleSchema, FullyQualifiedNameAttributeName, "name"),
		)),
	}
}

func CreateContextNetworkRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	ruleType, err := sdk.ToNetworkRuleType(d.Get("type").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	ruleMode, err := sdk.ToNetworkRuleMode(d.Get("mode").(string))
	if err != nil {
		return diag.FromErr(err)
	}

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
	if err := stringAttributeCreateBuilder(d, "comment", req.WithComment); err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*provider.Context).Client
	if err := client.NetworkRules.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextNetworkRule(ctx, d, meta)
}

func ReadContextNetworkRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	diags := diag.Diagnostics{}
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

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
	errs := errors.Join(
		d.Set("name", networkRule.Name),
		d.Set("database", networkRule.DatabaseName),
		d.Set("schema", networkRule.SchemaName),
		d.Set("type", networkRule.Type),
		d.Set("value_list", networkRuleDescriptions.ValueList),
		d.Set("mode", networkRule.Mode),
		d.Set("comment", networkRule.Comment),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.NetworkRuleToSchema(networkRule)}),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.NetworkRuleDetailsToSchema(networkRuleDescriptions)}),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	return diags
}

func UpdateContextNetworkRule(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set := sdk.NewNetworkRuleSetRequest()
	unset := sdk.NewNetworkRuleUnsetRequest()

	errs := errors.Join(
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		setValueUpdate(d, "value_list", &set.ValueList, &unset.ValueList, func(v any) (sdk.NetworkRuleValue, error) {
			return sdk.NetworkRuleValue{Value: v.(string)}, nil
		}),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(*set, *sdk.NewNetworkRuleSetRequest()) {
		if err := client.NetworkRules.Alter(ctx, sdk.NewAlterNetworkRuleRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if !reflect.DeepEqual(*unset, *sdk.NewNetworkRuleUnsetRequest()) {
		if err := client.NetworkRules.Alter(ctx, sdk.NewAlterNetworkRuleRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextNetworkRule(ctx, d, meta)
}
