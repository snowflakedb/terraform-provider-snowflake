package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// columnSchema builds the rich "column" schema shared by table-like resources that support the
// full per-column capabilities (name, type, not_null, comment, default, masking_policy, projection_policy).
func columnSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		ForceNew:    true,
		MinItems:    1,
		Description: "Definitions of the columns to create in the table. Minimum one required.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Column name.",
				},
				"type": {
					Type:             schema.TypeString,
					Required:         true,
					ForceNew:         true,
					Description:      "Column type, e.g. VARIANT. For a full list of column types, see [Summary of Data Types](https://docs.snowflake.com/en/sql-reference/intro-summary-data-types).",
					ValidateDiagFunc: IsDataTypeValid,
					DiffSuppressFunc: DiffSuppressDataTypes,
				},
				"not_null": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Default:          BooleanDefault,
					ValidateDiagFunc: validateBooleanString,
					DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
						return (oldValue == BooleanTrue || oldValue == BooleanFalse) && (newValue == BooleanDefault)
					},
					Description: "Whether to restrict the column to NOT NULL values.",
				},
				"default": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					MaxItems:    1,
					Description: "Defines the column default value.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"expression": {
								Type:        schema.TypeString,
								Required:    true,
								ForceNew:    true,
								Description: "The default expression value for the column.",
							},
						},
					},
				},
				"masking_policy":    columnMaskingPolicySchema(true, IgnoreMatchingColumnNameAndMaskingPolicyUsingFirstElem("name")),
				"projection_policy": columnProjectionPolicySchema(),
				"comment": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Column comment.",
				},
			},
		},
	}
}

func columnMaskingPolicySchema(forceNew bool, usingDiffSuppress schema.SchemaDiffSuppressFunc) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    forceNew,
		MaxItems:    1,
		Description: relatedResourceDescription("Specifies the masking policy to set on a column.", resources.MaskingPolicy),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"policy_name": {
					Type:             schema.TypeString,
					Required:         true,
					ForceNew:         forceNew,
					DiffSuppressFunc: suppressIdentifierQuoting,
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
					Description:      relatedResourceDescription("Masking policy name.", resources.MaskingPolicy),
				},
				"using": {
					Type:             schema.TypeList,
					Optional:         true,
					ForceNew:         forceNew,
					Elem:             &schema.Schema{Type: schema.TypeString},
					DiffSuppressFunc: usingDiffSuppress,
					Description:      "Specifies the arguments to pass into the conditional masking policy SQL expression, in order. The first column in the list specifies the column for the policy conditions to mask or tokenize the data and must match the column to which the masking policy is set. The additional columns specify the columns to evaluate to determine whether to mask or tokenize the data in each row of the query result when a query is made on the first column. If the USING clause is omitted, Snowflake treats the conditional masking policy as a normal masking policy.",
				},
			},
		},
	}
}

func columnProjectionPolicySchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		MaxItems:    1,
		Description: "Specifies the projection policy to set on a column.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"policy_name": {
					Type:             schema.TypeString,
					Required:         true,
					ForceNew:         true,
					DiffSuppressFunc: suppressIdentifierQuoting,
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
					Description:      "Projection policy name.",
				},
			},
		},
	}
}

func parseColumnMaskingPolicy(d *schema.ResourceData, prefix string) (sdk.TableColumnMaskingPolicyRequest, error) {
	id, err := sdk.ParseSchemaObjectIdentifier(d.Get(prefix + "policy_name").(string))
	if err != nil {
		return sdk.TableColumnMaskingPolicyRequest{}, err
	}
	req := sdk.TableColumnMaskingPolicyRequest{MaskingPolicy: id}
	if usingRaw := d.Get(prefix + "using").([]any); len(usingRaw) > 0 {
		req.Using = collections.Map(expandStringList(usingRaw), func(v string) sdk.Column { return sdk.Column{Value: v} })
	}
	return req, nil
}

func parseColumnProjectionPolicy(d *schema.ResourceData, prefix string) (sdk.TableColumnProjectionPolicyRequest, error) {
	id, err := sdk.ParseSchemaObjectIdentifier(d.Get(prefix + "policy_name").(string))
	if err != nil {
		return sdk.TableColumnProjectionPolicyRequest{}, err
	}
	return sdk.TableColumnProjectionPolicyRequest{ProjectionPolicy: id}, nil
}
