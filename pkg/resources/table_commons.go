package resources

import (
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func rowAccessPolicyFieldSchema(objectKind string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"policy_name": {
					Type:             schema.TypeString,
					Required:         true,
					DiffSuppressFunc: suppressIdentifierQuoting,
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
					Description:      relatedResourceDescription("Row access policy name.", resources.RowAccessPolicy),
				},
				"on": {
					Type:     schema.TypeSet,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: "Defines which columns are affected by the policy.",
				},
			},
		},
		Description: fmt.Sprintf("Specifies the row access policy to set on a %s.", objectKind),
	}
}

// aggregationPolicySchema builds the aggregation_policy schema.
func aggregationPolicySchema(objectKind string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"policy_name": {
					Type:             schema.TypeString,
					Required:         true,
					DiffSuppressFunc: suppressIdentifierQuoting,
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
					Description:      "Aggregation policy name.",
				},
				"entity_key": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Description: fmt.Sprintf("Defines which columns uniquely identify an entity within the %s.", objectKind),
				},
			},
		},
		Description: fmt.Sprintf("Specifies the aggregation policy to set on a %s.", objectKind),
	}
}

func extractPolicyWithColumnsSet(v any, columnsKey string) (sdk.SchemaObjectIdentifier, []sdk.Column, error) {
	policyConfig := v.([]any)[0].(map[string]any)
	id, err := sdk.ParseSchemaObjectIdentifier(policyConfig["policy_name"].(string))
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, nil, err
	}
	if policyConfig[columnsKey] == nil {
		return id, nil, nil
	}
	columnsRaw := expandStringList(policyConfig[columnsKey].(*schema.Set).List())
	return id, collections.Map(columnsRaw, func(c string) sdk.Column { return sdk.Column{Value: c} }), nil
}

func extractPolicyWithColumnsList(v any, columnsKey string) (sdk.SchemaObjectIdentifier, []sdk.Column, error) {
	policyConfig := v.([]any)[0].(map[string]any)
	id, err := sdk.ParseSchemaObjectIdentifier(policyConfig["policy_name"].(string))
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, nil, err
	}
	if policyConfig[columnsKey] == nil {
		return id, nil, fmt.Errorf("unable to extract policy with column list, unable to find columnsKey: %s", columnsKey)
	}
	columnsRaw := expandStringList(policyConfig[columnsKey].([]any))
	return id, collections.Map(columnsRaw, func(c string) sdk.Column { return sdk.Column{Value: c} }), nil
}

// rowAccessPolicyAlterRequests computes the ViewAddRowAccessPolicyRequest/ViewDropRowAccessPolicyRequest to apply
// for a row_access_policy diff. These request types are shared verbatim between AlterViewRequest and
// AlterIcebergTableRequest, so the diffing logic can be shared even though the two resources wire the result into
// different top-level alter requests.
func rowAccessPolicyAlterRequests(
	d ResourceValueGetter,
	addRowAccessPolicyFunc func(id sdk.SchemaObjectIdentifier, columns []sdk.Column),
	dropRowAccessPolicyFunc func(id sdk.SchemaObjectIdentifier),
) error {
	oldRaw, newRaw := d.GetChange("row_access_policy")
	if len(oldRaw.([]any)) > 0 {
		oldId, _, err := extractPolicyWithColumnsSet(oldRaw, "on")
		if err != nil {
			return err
		}
		dropRowAccessPolicyFunc(oldId)
	}
	if len(newRaw.([]any)) > 0 {
		newId, newColumns, err := extractPolicyWithColumnsSet(newRaw, "on")
		if err != nil {
			return err
		}
		addRowAccessPolicyFunc(newId, newColumns)
	}
	return nil
}

// aggregationPolicyAlterState extracts the desired aggregation policy state (id + entity key) from the
// aggregation_policy field, or reports that the policy should be unset.
func aggregationPolicyAlterState(d ResourceValueGetter) (id sdk.SchemaObjectIdentifier, entityKey []sdk.Column, isSet bool, err error) {
	v, ok := d.GetOk("aggregation_policy")
	if !ok {
		return sdk.SchemaObjectIdentifier{}, nil, false, nil
	}
	id, entityKey, err = extractPolicyWithColumnsSet(v, "entity_key")
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, nil, false, err
	}
	return id, entityKey, true, nil
}

// ResourceValueGetter is the subset of *schema.ResourceData used by the shared policy helpers.
type ResourceValueGetter interface {
	GetChange(key string) (any, any)
	GetOk(key string) (any, bool)
}

// handlePolicyReferences populates row_access_policy and aggregation_policy from the given policy references.
func handlePolicyReferences(policyRefs []sdk.PolicyReference, d ResourceValueSetter) error {
	var aggregationPolicies []map[string]any
	var rowAccessPolicies []map[string]any
	for _, p := range policyRefs {
		policyName := sdk.NewSchemaObjectIdentifier(*p.PolicyDb, *p.PolicySchema, p.PolicyName)
		switch p.PolicyKind {
		case sdk.PolicyKindAggregationPolicy:
			var entityKey []string
			if p.RefArgColumnNames != nil {
				entityKey = sdk.ParseCommaSeparatedStringArray(*p.RefArgColumnNames, true)
			}
			aggregationPolicies = append(aggregationPolicies, map[string]any{
				"policy_name": policyName.FullyQualifiedName(),
				"entity_key":  entityKey,
			})
		case sdk.PolicyKindRowAccessPolicy:
			var on []string
			if p.RefArgColumnNames != nil {
				on = sdk.ParseCommaSeparatedStringArray(*p.RefArgColumnNames, true)
			}
			rowAccessPolicies = append(rowAccessPolicies, map[string]any{
				"policy_name": policyName.FullyQualifiedName(),
				"on":          on,
			})
		default:
			log.Printf("[DEBUG] unexpected policy kind %v in policy references returned from Snowflake", p.PolicyKind)
		}
	}
	if err := d.Set("aggregation_policy", aggregationPolicies); err != nil {
		return err
	}
	if err := d.Set("row_access_policy", rowAccessPolicies); err != nil {
		return err
	}
	return nil
}
