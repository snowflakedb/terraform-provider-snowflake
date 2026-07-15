package resources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
// covering a table-level PRIMARY KEY constraint. MaxItems is capped at 1 because a table can have at
// most one PRIMARY KEY constraint.
func primaryKeyConstraintSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		MaxItems:    1,
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

// rowAccessPolicyCreateRequest extracts the row_access_policy field, if set, into a ready-to-use
// ViewRowAccessPolicyRequest-shaped id + columns pair for use on create. ok is false when the field is unset.
func rowAccessPolicyCreateRequest(d *schema.ResourceData) (id sdk.SchemaObjectIdentifier, columns []sdk.Column, ok bool, err error) {
	v := d.Get("row_access_policy")
	if len(v.([]any)) == 0 {
		return sdk.SchemaObjectIdentifier{}, nil, false, nil
	}
	id, columns, err = extractPolicyWithColumnsSet(v, "on")
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, nil, false, err
	}
	return id, columns, true, nil
}

// rowAccessPolicyUpdateRequests computes the ready-to-use ViewAddRowAccessPolicyRequest/ViewDropRowAccessPolicyRequest
// for a row_access_policy diff. These request types are shared verbatim between AlterViewRequest and
// AlterIcebergTableRequest, so the diffing logic can be shared even though the two resources wire the result into
// different top-level alter requests.
func rowAccessPolicyUpdateRequests(d *schema.ResourceData) (add *sdk.ViewAddRowAccessPolicyRequest, drop *sdk.ViewDropRowAccessPolicyRequest, err error) {
	oldRaw, newRaw := d.GetChange("row_access_policy")
	if len(oldRaw.([]any)) > 0 {
		oldId, _, err := extractPolicyWithColumnsSet(oldRaw, "on")
		if err != nil {
			return nil, nil, err
		}
		drop = sdk.NewViewDropRowAccessPolicyRequest(oldId)
	}
	if len(newRaw.([]any)) > 0 {
		newId, newColumns, err := extractPolicyWithColumnsSet(newRaw, "on")
		if err != nil {
			return nil, nil, err
		}
		add = sdk.NewViewAddRowAccessPolicyRequest(newId, newColumns)
	}
	return add, drop, nil
}

// aggregationPolicyCreateRequest extracts the aggregation_policy field, if set, into a ready-to-use
// ViewAggregationPolicyRequest-shaped id + columns pair for use on create. ok is false when the field is unset.
func aggregationPolicyCreateRequest(d *schema.ResourceData) (id sdk.SchemaObjectIdentifier, entityKey []sdk.Column, ok bool, err error) {
	v := d.Get("aggregation_policy")
	if len(v.([]any)) == 0 {
		return sdk.SchemaObjectIdentifier{}, nil, false, nil
	}
	id, entityKey, err = extractPolicyWithColumnsSet(v, "entity_key")
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, nil, false, err
	}
	return id, entityKey, true, nil
}

// aggregationPolicyUpdateRequests computes the ready-to-use ViewSetAggregationPolicyRequest/ViewUnsetAggregationPolicyRequest
// for the desired aggregation_policy state. These request types are shared verbatim between AlterViewRequest and
// AlterIcebergTableRequest.
func aggregationPolicyUpdateRequests(d *schema.ResourceData) (set *sdk.ViewSetAggregationPolicyRequest, unset *sdk.ViewUnsetAggregationPolicyRequest, err error) {
	v, ok := d.GetOk("aggregation_policy")
	if !ok {
		return nil, sdk.NewViewUnsetAggregationPolicyRequest(), nil
	}
	id, entityKey, err := extractPolicyWithColumnsSet(v, "entity_key")
	if err != nil {
		return nil, nil, err
	}
	set = sdk.NewViewSetAggregationPolicyRequest(id)
	if len(entityKey) > 0 {
		set.WithEntityKey(entityKey)
	}
	return set.WithForce(true), nil, nil
}

// handlePolicyReferences populates row_access_policy and aggregation_policy from the given policy references.
func handlePolicyReferences(policyRefs []sdk.PolicyReference, d *schema.ResourceData) error {
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

// readRootLevelPolicies fetches the policy references for id in the given domain, populates
// row_access_policy and aggregation_policy in the resource state, and returns the fetched references
// so callers can also derive column-level (masking/projection) policy state from them.
func readRootLevelPolicies(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, domain sdk.PolicyEntityDomain, d *schema.ResourceData) ([]sdk.PolicyReference, error) {
	policyRefs, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(id, domain))
	if err != nil {
		return nil, err
	}
	if err := handlePolicyReferences(policyRefs, d); err != nil {
		return nil, err
	}
	return policyRefs, nil
}

// buildColumnPolicyLookups indexes policyRefs by referenced column name, once, so that
// columnPoliciesToState can look up a column's masking/projection policy in O(1) instead of
// rescanning the full policyRefs slice for every column.
func buildColumnPolicyLookups(policyRefs []sdk.PolicyReference) (maskingPolicies, projectionPolicies map[string]sdk.PolicyReference) {
	maskingPolicies = make(map[string]sdk.PolicyReference)
	projectionPolicies = make(map[string]sdk.PolicyReference)
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
	return maskingPolicies, projectionPolicies
}

// columnPoliciesToState converts columnName's masking/projection policy (if any) into column state
// fields. The two lookups are independent: a conversion error on one policy is reported but does not
// prevent the other, valid, policy from being included in the returned state.
func columnPoliciesToState(columnName string, policyRefs []sdk.PolicyReference) (map[string]any, error) {
	maskingPolicies, projectionPolicies := buildColumnPolicyLookups(policyRefs)
	columnPolicies := make(map[string]any)
	var errs []error
	if p, ok := maskingPolicies[columnName]; ok {
		if maskingPolicyState, err := maskingPolicyToState(p); err != nil {
			errs = append(errs, fmt.Errorf("converting masking policy to state: %w", err))
		} else {
			columnPolicies["masking_policy"] = []map[string]any{maskingPolicyState}
		}
	}
	if p, ok := projectionPolicies[columnName]; ok {
		if projectionPolicyState, err := projectionPolicyToState(p); err != nil {
			errs = append(errs, fmt.Errorf("converting projection policy to state: %w", err))
		} else {
			columnPolicies["projection_policy"] = []map[string]any{projectionPolicyState}
		}
	}
	return columnPolicies, errors.Join(errs...)
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
