package resources

import (
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// columnSchema builds the rich "column" schema shared by table-like resources that support the
// full per-column capabilities (as opposed to basicColumnSchema, which only supports name + type).
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
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
					Default:     false,
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
								Optional:    true,
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

// parseColumns parses the "column" list from the resource data into IcebergTableColumnRequests.
func parseColumns(d *schema.ResourceData) ([]sdk.IcebergTableColumnRequest, error) {
	raw := d.Get("column").([]any)
	columns := make([]sdk.IcebergTableColumnRequest, len(raw))
	for i := range raw {
		column, err := parseColumn(d, fmt.Sprintf("column.%d.", i))
		if err != nil {
			return nil, err
		}
		columns[i] = column
	}
	return columns, nil
}

// parseColumn parses a single column at the given prefix (e.g. "column.0.") into an IcebergTableColumnRequest.
func parseColumn(d *schema.ResourceData, prefix string) (sdk.IcebergTableColumnRequest, error) {
	name := d.Get(prefix + "name").(string)
	dataType, err := datatypes.ParseDataType(d.Get(prefix + "type").(string))
	if err != nil {
		return sdk.IcebergTableColumnRequest{}, fmt.Errorf("parsing data type of column %q: %w", name, err)
	}
	req := sdk.NewIcebergTableColumnRequest(name, dataType)

	if err := errors.Join(
		boolAttributeCreate(d, prefix+"not_null", &req.NotNull),
		stringAttributeCreateBuilder(d, prefix+"comment", func(v string) *sdk.IcebergTableColumnRequest { return req.WithComment(v) }),
		attributeMappedValueCreateBuilderNested(d, prefix+"default", func(v sdk.ColumnDefaultValue) *sdk.IcebergTableColumnRequest {
			return req.WithDefaultValue(v)
		}, func(d *schema.ResourceData) (sdk.ColumnDefaultValue, error) {
			return parseColumnDefaultValue(d, prefix+"default.0."), nil
		}),
		attributeMappedValueCreateBuilderNested(d, prefix+"masking_policy", func(v sdk.TableColumnMaskingPolicyRequest) *sdk.IcebergTableColumnRequest {
			return req.WithMaskingPolicy(v)
		}, func(d *schema.ResourceData) (sdk.TableColumnMaskingPolicyRequest, error) {
			return parseColumnMaskingPolicy(d, prefix+"masking_policy.0.")
		}),
		attributeMappedValueCreateBuilderNested(d, prefix+"projection_policy", func(v sdk.TableColumnProjectionPolicyRequest) *sdk.IcebergTableColumnRequest {
			return req.WithProjectionPolicy(v)
		}, func(d *schema.ResourceData) (sdk.TableColumnProjectionPolicyRequest, error) {
			return parseColumnProjectionPolicy(d, prefix+"projection_policy.0.")
		}),
	); err != nil {
		return sdk.IcebergTableColumnRequest{}, fmt.Errorf("parsing column %q: %w", name, err)
	}

	return *req, nil
}

func parseColumnDefaultValue(d *schema.ResourceData, prefix string) sdk.ColumnDefaultValue {
	defaultValue := sdk.ColumnDefaultValue{}
	if v, ok := d.GetOk(prefix + "expression"); ok {
		expression := v.(string)
		defaultValue.Expression = &expression
	}
	return defaultValue
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

// constraintEnforcementSchemaFields builds the enforcement-related fields shared by the UNIQUE/PRIMARY KEY
// and FOREIGN KEY out-of-line constraint schemas. Each field squashes a pair of mutually exclusive SQL
// keywords (e.g. ENFORCED / NOT ENFORCED) into a single tri-state boolean string field.
func constraintEnforcementSchemaFields() map[string]*schema.Schema {
	boolStringField := func(description string) *schema.Schema {
		return &schema.Schema{
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			Default:          BooleanDefault,
			ValidateDiagFunc: validateBooleanString,
			Description:      booleanStringFieldDescription(description),
		}
	}
	return map[string]*schema.Schema{
		"enforced":           boolStringField("Whether the constraint is enforced (`true`) or not enforced (`false`)."),
		"deferrable":         boolStringField("Whether the constraint is deferrable (`true`) or not deferrable (`false`)."),
		"initially_deferred": boolStringField("Whether the constraint is initially deferred (`true`) or initially immediate (`false`)."),
		"enable":             boolStringField("Whether the constraint is enabled (`true`) or disabled (`false`)."),
		"validate":           boolStringField("Whether to validate existing data on the table when the constraint is created (`true`) or skip validation (`false`)."),
		"rely":               boolStringField("Whether a constraint in NOVALIDATE mode is taken into account (`true`) or not (`false`) during query rewrite."),
	}
}

// outOfLineUniqueOrPKConstraintSchemaFields builds the fields shared by the primary_key_constraint
// and unique_constraint schemas. The two attributes are distinguished by their own name (there is
// no boolean discriminator field inside the block).
func outOfLineUniqueOrPKConstraintSchemaFields() map[string]*schema.Schema {
	return collections.MergeMaps(constraintEnforcementSchemaFields(), map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Name of the constraint.",
		},
		"column": {
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			MinItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "The column(s) the constraint applies to.",
		},
		"comment": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Constraint comment.",
		},
	})
}

// primaryKeyConstraintSchema builds the primary_key_constraint schema shared by table-like resources,
// covering a table-level PRIMARY KEY constraint.
func primaryKeyConstraintSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Defines a table-level PRIMARY KEY constraint.",
		Elem: &schema.Resource{
			Schema: outOfLineUniqueOrPKConstraintSchemaFields(),
		},
	}
}

// uniqueConstraintSchema builds the unique_constraint schema shared by table-like resources,
// covering a table-level UNIQUE constraint.
func uniqueConstraintSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Defines a table-level UNIQUE constraint.",
		Elem: &schema.Resource{
			Schema: outOfLineUniqueOrPKConstraintSchemaFields(),
		},
	}
}

// foreignKeyConstraintSchema builds the foreign_key_constraint schema shared by table-like resources,
// covering a table-level FOREIGN KEY constraint.
func foreignKeyConstraintSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Defines a table-level FOREIGN KEY constraint.",
		Elem: &schema.Resource{
			Schema: collections.MergeMaps(constraintEnforcementSchemaFields(), map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the constraint.",
				},
				"column": {
					Type:        schema.TypeList,
					Required:    true,
					ForceNew:    true,
					MinItems:    1,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "The local column(s) the foreign key is defined on.",
				},
				"table_name": {
					Type:             schema.TypeString,
					Required:         true,
					ForceNew:         true,
					DiffSuppressFunc: suppressIdentifierQuoting,
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
					Description:      "The table that the foreign key references.",
				},
				"ref_column": {
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "The column(s) in the referenced table that the foreign key references.",
				},
				"match": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					ValidateDiagFunc: sdkValidation(sdk.ToMatchType),
					Description:      fmt.Sprintf("The match type for the foreign key. Valid values are: %v.", sdk.AllMatchTypes),
				},
				"on_update": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					ValidateDiagFunc: sdkValidation(sdk.ToForeignKeyAction),
					Description:      fmt.Sprintf("Specifies the action to perform when the referenced primary/unique key is updated. Valid values are: %v.", sdk.AllForeignKeyActions),
				},
				"on_delete": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					ValidateDiagFunc: sdkValidation(sdk.ToForeignKeyAction),
					Description:      fmt.Sprintf("Specifies the action to perform when the referenced primary/unique key is deleted. Valid values are: %v.", sdk.AllForeignKeyActions),
				},
				"comment": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Constraint comment.",
				},
			}),
		},
	}
}

// checkConstraintSchema builds the check_constraint schema shared by table-like resources, covering
// a table-level CHECK constraint.
func checkConstraintSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Defines a table-level CHECK constraint.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Name of the constraint.",
				},
				"expression": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "The CHECK constraint expression.",
				},
				"validate": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					Default:          BooleanDefault,
					ValidateDiagFunc: validateBooleanString,
					Description:      booleanStringFieldDescription("Whether existing data is validated against the constraint (`true`, `ENABLE VALIDATE`) or not (`false`, `ENABLE NOVALIDATE`)."),
				},
			},
		},
	}
}

// parseOutOfLineConstraints parses the "primary_key_constraint", "unique_constraint",
// "foreign_key_constraint", and "check_constraint" lists from the resource data into
// TableOutOfLineConstraintRequests.
func parseOutOfLineConstraints(d *schema.ResourceData) ([]sdk.TableOutOfLineConstraintRequest, error) {
	var constraints []sdk.TableOutOfLineConstraintRequest

	primaryKeyRaw := d.Get("primary_key_constraint").([]any)
	for i := range primaryKeyRaw {
		primaryKey, err := parseOutOfLinePrimaryKey(d, fmt.Sprintf("primary_key_constraint.%d.", i))
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, sdk.TableOutOfLineConstraintRequest{UniquePK: &primaryKey})
	}

	uniqueRaw := d.Get("unique_constraint").([]any)
	for i := range uniqueRaw {
		unique, err := parseOutOfLineUnique(d, fmt.Sprintf("unique_constraint.%d.", i))
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, sdk.TableOutOfLineConstraintRequest{UniquePK: &unique})
	}

	foreignKeyRaw := d.Get("foreign_key_constraint").([]any)
	for i := range foreignKeyRaw {
		fk, err := parseOutOfLineForeignKey(d, fmt.Sprintf("foreign_key_constraint.%d.", i))
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, sdk.TableOutOfLineConstraintRequest{FK: &fk})
	}

	checkRaw := d.Get("check_constraint").([]any)
	for i := range checkRaw {
		prefix := fmt.Sprintf("check_constraint.%d.", i)
		ch := sdk.TableOutOfLineCHRequest{Expression: d.Get(prefix + "expression").(string)}
		if err := errors.Join(
			stringAttributeCreate(d, prefix+"name", &ch.Name),
			booleanStringPairAttributeCreate(d, prefix+"validate", &ch.EnableValidate, &ch.EnableNovalidate),
		); err != nil {
			return nil, err
		}
		constraints = append(constraints, sdk.TableOutOfLineConstraintRequest{CH: &ch})
	}

	return constraints, nil
}

// parseOutOfLineUniquePKCommon parses the fields shared by primary_key_constraint and
// unique_constraint. Callers set the Unique/PrimaryKey discriminator themselves.
func parseOutOfLineUniquePKCommon(d *schema.ResourceData, prefix string) (sdk.TableOutOfLineUniquePKRequest, error) {
	uniquePK := sdk.TableOutOfLineUniquePKRequest{}
	if err := errors.Join(
		stringAttributeCreate(d, prefix+"name", &uniquePK.Name),
		stringAttributeCreate(d, prefix+"comment", &uniquePK.Comment),
		booleanStringPairAttributeCreate(d, prefix+"enforced", &uniquePK.Enforced, &uniquePK.NotEnforced),
		booleanStringPairAttributeCreate(d, prefix+"deferrable", &uniquePK.Deferrable, &uniquePK.NotDeferrable),
		booleanStringPairAttributeCreate(d, prefix+"initially_deferred", &uniquePK.InitiallyDeferred, &uniquePK.InitiallyImmediate),
		booleanStringPairAttributeCreate(d, prefix+"enable", &uniquePK.Enable, &uniquePK.Disable),
		booleanStringPairAttributeCreate(d, prefix+"validate", &uniquePK.Validate, &uniquePK.Novalidate),
		booleanStringPairAttributeCreate(d, prefix+"rely", &uniquePK.Rely, &uniquePK.Norely),
	); err != nil {
		return sdk.TableOutOfLineUniquePKRequest{}, err
	}
	if columnRaw := d.Get(prefix + "column").([]any); len(columnRaw) > 0 {
		uniquePK.Columns = collections.Map(expandStringList(columnRaw), func(c string) sdk.Column { return sdk.Column{Value: c} })
	}
	return uniquePK, nil
}

// parseOutOfLinePrimaryKey parses a single primary_key_constraint entry.
func parseOutOfLinePrimaryKey(d *schema.ResourceData, prefix string) (sdk.TableOutOfLineUniquePKRequest, error) {
	primaryKey, err := parseOutOfLineUniquePKCommon(d, prefix)
	if err != nil {
		return sdk.TableOutOfLineUniquePKRequest{}, err
	}
	primaryKey.PrimaryKey = new(true)
	return primaryKey, nil
}

// parseOutOfLineUnique parses a single unique_constraint entry.
func parseOutOfLineUnique(d *schema.ResourceData, prefix string) (sdk.TableOutOfLineUniquePKRequest, error) {
	unique, err := parseOutOfLineUniquePKCommon(d, prefix)
	if err != nil {
		return sdk.TableOutOfLineUniquePKRequest{}, err
	}
	unique.Unique = new(true)
	return unique, nil
}

func parseOutOfLineForeignKey(d *schema.ResourceData, prefix string) (sdk.TableOutOfLineFKRequest, error) {
	id, err := sdk.ParseSchemaObjectIdentifier(d.Get(prefix + "table_name").(string))
	if err != nil {
		return sdk.TableOutOfLineFKRequest{}, err
	}
	fk := sdk.TableOutOfLineFKRequest{References: id}

	if err := errors.Join(
		stringAttributeCreate(d, prefix+"name", &fk.Name),
		stringAttributeCreate(d, prefix+"comment", &fk.Comment),
		attributeMappedValueCreate(d, prefix+"match", &fk.Match, func(v any) (*sdk.MatchType, error) {
			matchType, err := sdk.ToMatchType(v.(string))
			return &matchType, err
		}),
		booleanStringPairAttributeCreate(d, prefix+"enforced", &fk.Enforced, &fk.NotEnforced),
		booleanStringPairAttributeCreate(d, prefix+"deferrable", &fk.Deferrable, &fk.NotDeferrable),
		booleanStringPairAttributeCreate(d, prefix+"initially_deferred", &fk.InitiallyDeferred, &fk.InitiallyImmediate),
		booleanStringPairAttributeCreate(d, prefix+"enable", &fk.Enable, &fk.Disable),
		booleanStringPairAttributeCreate(d, prefix+"validate", &fk.Validate, &fk.Novalidate),
		booleanStringPairAttributeCreate(d, prefix+"rely", &fk.Rely, &fk.Norely),
	); err != nil {
		return sdk.TableOutOfLineFKRequest{}, err
	}

	if columnRaw := d.Get(prefix + "column").([]any); len(columnRaw) > 0 {
		fk.Columns = collections.Map(expandStringList(columnRaw), func(c string) sdk.Column { return sdk.Column{Value: c} })
	}
	if refColumnRaw := d.Get(prefix + "ref_column").([]any); len(refColumnRaw) > 0 {
		fk.RefColumns = collections.Map(expandStringList(refColumnRaw), func(c string) sdk.Column { return sdk.Column{Value: c} })
	}

	onUpdate, onUpdateOk := d.GetOk(prefix + "on_update")
	onDelete, onDeleteOk := d.GetOk(prefix + "on_delete")
	if onUpdateOk || onDeleteOk {
		on := &sdk.ForeignKeyOnAction{}
		if onUpdateOk {
			action, err := sdk.ToForeignKeyAction(onUpdate.(string))
			if err != nil {
				return sdk.TableOutOfLineFKRequest{}, err
			}
			on.OnUpdate = &action
		}
		if onDeleteOk {
			action, err := sdk.ToForeignKeyAction(onDelete.(string))
			if err != nil {
				return sdk.TableOutOfLineFKRequest{}, err
			}
			on.OnDelete = &action
		}
		fk.On = on
	}

	return fk, nil
}

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

func columnPoliciesToState(columnName string, policyRefs []sdk.PolicyReference) (map[string]any, error) {
	columnPolicies := make(map[string]any)
	maskingPolicies := make(map[string]sdk.PolicyReference)
	projectionPolicies := make(map[string]sdk.PolicyReference)
	for _, p := range policyRefs {
		if p.RefColumnName == nil {
			continue
		}
		switch p.PolicyKind {
		case sdk.PolicyKindMaskingPolicy:
			maskingPolicies[*p.RefColumnName] = p
		case sdk.PolicyKindProjectionPolicy:
			projectionPolicies[*p.RefColumnName] = p
		}
	}
	if p, ok := maskingPolicies[columnName]; ok {
		maskingPolicyState, err := maskingPolicyToState(p)
		if err != nil {
			return nil, fmt.Errorf("converting masking policy to state: %w", err)
		} else {
			columnPolicies["masking_policy"] = []map[string]any{maskingPolicyState}
		}
	}
	if p, ok := projectionPolicies[columnName]; ok {
		projectionPolicyState, err := projectionPolicyToState(p)
		if err != nil {
			return nil, fmt.Errorf("converting projection policy to state: %w", err)
		} else {
			columnPolicies["projection_policy"] = []map[string]any{projectionPolicyState}
		}
	}
	return columnPolicies, nil
}

func maskingPolicyToState(maskingPolicy sdk.PolicyReference) (map[string]any, error) {
	if maskingPolicy.PolicyDb == nil || maskingPolicy.PolicySchema == nil {
		return nil, fmt.Errorf("policy db and schema can not be empty")
	}
	var usingArgs []string
	if maskingPolicy.RefArgColumnNames != nil {
		usingArgs = sdk.ParseCommaSeparatedStringArray(*maskingPolicy.RefArgColumnNames, true)
	}
	return map[string]any{
		"policy_name": sdk.NewSchemaObjectIdentifier(*maskingPolicy.PolicyDb, *maskingPolicy.PolicySchema, maskingPolicy.PolicyName).FullyQualifiedName(),
		"using":       append([]string{*maskingPolicy.RefColumnName}, usingArgs...),
	}, nil
}

func projectionPolicyToState(projectionPolicy sdk.PolicyReference) (map[string]any, error) {
	if projectionPolicy.PolicyDb == nil || projectionPolicy.PolicySchema == nil {
		return nil, fmt.Errorf("policy db and schema can not be empty")
	}
	return map[string]any{
		"policy_name": sdk.NewSchemaObjectIdentifier(*projectionPolicy.PolicyDb, *projectionPolicy.PolicySchema, projectionPolicy.PolicyName).FullyQualifiedName(),
	}, nil
}
