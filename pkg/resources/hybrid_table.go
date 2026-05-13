package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var hybridTableSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the hybrid table."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the hybrid table."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the hybrid table."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the hybrid table.",
	},
	"column": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		Description: "Definitions of a column to create in the hybrid table. Minimum one required.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Column name.",
				},
				"type": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Column type. See [Snowflake data types](https://docs.snowflake.com/en/sql-reference-data-types) for supported values. Example: VARCHAR(256), NUMBER(38,0).",
					ValidateDiagFunc: IsDataTypeValid,
					DiffSuppressFunc: DiffSuppressDataTypes,
				},
				"nullable": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Whether this column allows NULLs. Changing this on an existing column forces recreation because hybrid tables do not support ALTER SET/DROP NOT NULL.",
				},
				"default": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					MinItems:    1,
					Description: "Defines the column default value. Only one of constant, expression, or sequence may be set.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"constant": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "A constant default value for the column.",
							},
							"expression": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "A SQL expression default value for the column.",
							},
							"sequence": {
								Type:             schema.TypeString,
								Optional:         true,
								DiffSuppressFunc: suppressIdentifierQuoting,
								Description:      "The default sequence for the column (uses NEXTVAL).",
							},
						},
					},
				},
				"collate": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: "Column collation specification, e.g. en-ci.",
				},
				"comment": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: "Column-level comment.",
				},
			},
		},
	},
	"primary_key": {
		Type:        schema.TypeList,
		Required:    true,
		ForceNew:    true,
		MaxItems:    1,
		MinItems:    1,
		Description: "Defines the primary key constraint for the hybrid table. Required. Any change forces recreation.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Constraint name. If omitted, Snowflake auto-generates one.",
				},
				"keys": {
					Type:        schema.TypeList,
					Required:    true,
					ForceNew:    true,
					MinItems:    1,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Column names that form the primary key.",
				},
			},
		},
	},
	"unique_constraint": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Defines UNIQUE constraints. Can only be set at creation time. Any change forces recreation.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Constraint name.",
				},
				"columns": {
					Type:        schema.TypeList,
					Required:    true,
					ForceNew:    true,
					MinItems:    1,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Column names for the unique constraint.",
				},
			},
		},
	},
	"foreign_key": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		Description: "Defines FOREIGN KEY constraints. Can only be set at creation time. Any change forces recreation.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Constraint name.",
				},
				"columns": {
					Type:        schema.TypeList,
					Required:    true,
					ForceNew:    true,
					MinItems:    1,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Local column names.",
				},
				"references": {
					Type:        schema.TypeList,
					Required:    true,
					ForceNew:    true,
					MaxItems:    1,
					MinItems:    1,
					Description: "Referenced table and columns.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"table_id": {
								Type:             schema.TypeString,
								Required:         true,
								ForceNew:         true,
								Description:      "Fully qualified name of the referenced table.",
								ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
								DiffSuppressFunc: suppressIdentifierQuoting,
							},
							"columns": {
								Type:        schema.TypeList,
								Required:    true,
								ForceNew:    true,
								MinItems:    1,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Description: "Referenced column names.",
							},
						},
					},
				},
			},
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW HYBRID TABLES` for the given hybrid table.",
		Elem: &schema.Resource{
			Schema: schemas.ShowHybridTableSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE TABLE` for the given hybrid table.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeHybridTableSchema,
		},
	},
}

func HybridTable() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.HybridTables.DropSafely
		},
	)
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.HybridTableResource), TrackingCreateWrapper(resources.HybridTable, CreateHybridTable)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.HybridTableResource), TrackingReadWrapper(resources.HybridTable, GetReadHybridTableFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.HybridTableResource), TrackingUpdateWrapper(resources.HybridTable, UpdateHybridTable)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.HybridTableResource), TrackingDeleteWrapper(resources.HybridTable, deleteFunc)),
		Description:   "Resource used to manage hybrid tables. For more information, check [hybrid tables documentation](https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.HybridTable, customdiff.All(
			hybridTableParametersCustomDiff,
			ComputedIfAnyAttributeChanged(hybridTableSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(hybridTableSchema, DescribeOutputAttributeName, "name", "comment", "column"),
			ComputedIfAnyAttributeChanged(hybridTableSchema, FullyQualifiedNameAttributeName, "name"),
			forceNewIfColumnCollateChanged(),
			forceNewIfColumnNullableChanged(),
		)),

		Schema: collections.MergeMaps(hybridTableSchema, hybridTableParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.HybridTable, ImportHybridTable),
		},

		Timeouts: defaultTimeouts,
	}
}

// ---------------------------------------------------------------------------
// Column types and parsing helpers
// ---------------------------------------------------------------------------

type hybridTableColumn struct {
	name     string
	dataType string
	nullable bool
	_default *columnDefault // Reuses table.go's columnDefault type (same package)
	collate  string
	comment  string
}

type hybridTableColumns []hybridTableColumn

type changedHybridColumn struct {
	newColumn       hybridTableColumn
	changedDataType bool
	droppedDefault  bool
	changedComment  bool
}

func parseHybridColumn(from any) hybridTableColumn {
	c := from.(map[string]any)
	var cd *columnDefault

	if defaultList, ok := c["default"].([]any); ok && len(defaultList) == 1 {
		cd = getColumnDefault(defaultList[0].(map[string]any))
	}

	return hybridTableColumn{
		name:     c["name"].(string),
		dataType: c["type"].(string),
		nullable: c["nullable"].(bool),
		_default: cd,
		collate:  c["collate"].(string),
		comment:  c["comment"].(string),
	}
}

func parseHybridColumns(from any) hybridTableColumns {
	cols := from.([]any)
	result := make(hybridTableColumns, len(cols))
	for i, c := range cols {
		result[i] = parseHybridColumn(c)
	}
	return result
}

func (c hybridTableColumns) getNewIn(other hybridTableColumns) hybridTableColumns {
	added := hybridTableColumns{}
	for _, oldColumn := range c {
		if _, err := collections.FindFirst(other, func(newColumn hybridTableColumn) bool {
			return oldColumn.name == newColumn.name
		}); err != nil {
			added = append(added, oldColumn)
		}
	}
	return added
}

func (c hybridTableColumns) getChangedColumnProperties(new hybridTableColumns) []changedHybridColumn {
	changed := make([]changedHybridColumn, 0)
	for _, oldColumn := range c {
		for _, newColumn := range new {
			if oldColumn.name != newColumn.name {
				continue
			}
			ch := changedHybridColumn{newColumn: newColumn}
			if oldColumn.dataType != newColumn.dataType {
				ch.changedDataType = true
			}
			if oldColumn._default != nil && newColumn._default == nil {
				ch.droppedDefault = true
			}
			if oldColumn.comment != newColumn.comment {
				ch.changedComment = true
			}
			// NOTE: collate changes are not detected here because the SDK's
			// HybridTableAlterColumnActionRequest has no WithCollate method.
			// Collation can only be set at column creation time (CREATE or ADD COLUMN).
			if ch.changedDataType || ch.droppedDefault || ch.changedComment {
				changed = append(changed, ch)
			}
			break
		}
	}
	return changed
}

func (c hybridTableColumns) diffs(new hybridTableColumns) (removed hybridTableColumns, added hybridTableColumns, changed []changedHybridColumn) {
	return c.getNewIn(new), new.getNewIn(c), c.getChangedColumnProperties(new)
}

// ---------------------------------------------------------------------------
// Create helpers
// ---------------------------------------------------------------------------

func buildHybridTableColumnRequests(cols []any) ([]sdk.HybridTableColumnRequest, error) {
	requests := make([]sdk.HybridTableColumnRequest, len(cols))
	for i, rawCol := range cols {
		c := rawCol.(map[string]any)
		colType := c["type"].(string)
		dataType, err := datatypes.ParseDataType(colType)
		if err != nil {
			return nil, fmt.Errorf("invalid data type for column %s: %w", c["name"].(string), err)
		}

		req := sdk.NewHybridTableColumnRequest(c["name"].(string), sdk.DataType(colType))

		if nullable, ok := c["nullable"].(bool); ok && !nullable {
			req.WithNotNull(true)
		}

		if defaultList, ok := c["default"].([]any); ok && len(defaultList) == 1 {
			defMap := defaultList[0].(map[string]any)
			if defaultValue := buildHybridColumnDefaultValue(defMap, dataType); defaultValue != nil {
				req.WithDefaultValue(*defaultValue)
			}
		}

		if collate, ok := c["collate"].(string); ok && collate != "" {
			req.WithCollate(collate)
		}

		if comment, ok := c["comment"].(string); ok && comment != "" {
			req.WithComment(comment)
		}

		requests[i] = *req
	}
	return requests, nil
}

// buildHybridColumnDefaultValue builds a ColumnDefaultValue from a default block map.
// Mutual exclusivity of constant/expression/sequence is not validated here;
// invalid combinations are rejected by Snowflake (matches table.go approach).
func buildHybridColumnDefaultValue(defMap map[string]any, dataType datatypes.DataType) *sdk.ColumnDefaultValue {
	constant, hasConstant := defMap["constant"].(string)
	hasConstant = hasConstant && len(constant) > 0

	expr, hasExpression := defMap["expression"].(string)
	hasExpression = hasExpression && len(expr) > 0

	seq, hasSequence := defMap["sequence"].(string)
	hasSequence = hasSequence && len(seq) > 0

	var expression string
	switch {
	case hasConstant:
		if datatypes.IsTextDataType(dataType) {
			expression = snowflake.EscapeSnowflakeString(constant)
		} else {
			expression = constant
		}
	case hasExpression:
		expression = expr
	case hasSequence:
		expression = fmt.Sprintf("%v.NEXTVAL", seq)
	default:
		return nil
	}

	return &sdk.ColumnDefaultValue{Expression: &expression}
}

func buildOutOfLineConstraints(d *schema.ResourceData) ([]sdk.HybridTableOutOfLineConstraintRequest, error) {
	constraints := make([]sdk.HybridTableOutOfLineConstraintRequest, 0)

	// Primary key (required)
	pkList := d.Get("primary_key").([]any)
	pkMap := pkList[0].(map[string]any)
	pkConstraint := sdk.NewHybridTableOutOfLineConstraintRequest(sdk.ColumnConstraintTypePrimaryKey).
		WithColumns(expandStringList(pkMap["keys"].([]any)))
	if pkName, ok := pkMap["name"].(string); ok && pkName != "" {
		pkConstraint.WithName(pkName)
	}
	constraints = append(constraints, *pkConstraint)

	// Unique constraints (optional)
	if v, ok := d.GetOk("unique_constraint"); ok {
		for _, ucRaw := range v.([]any) {
			ucMap := ucRaw.(map[string]any)
			ucConstraint := sdk.NewHybridTableOutOfLineConstraintRequest(sdk.ColumnConstraintTypeUnique).
				WithColumns(expandStringList(ucMap["columns"].([]any)))
			if ucName, ok := ucMap["name"].(string); ok && ucName != "" {
				ucConstraint.WithName(ucName)
			}
			constraints = append(constraints, *ucConstraint)
		}
	}

	// Foreign keys (optional)
	if v, ok := d.GetOk("foreign_key"); ok {
		for _, fkRaw := range v.([]any) {
			fkMap := fkRaw.(map[string]any)
			fkConstraint := sdk.NewHybridTableOutOfLineConstraintRequest(sdk.ColumnConstraintTypeForeignKey).
				WithColumns(expandStringList(fkMap["columns"].([]any)))
			if fkName, ok := fkMap["name"].(string); ok && fkName != "" {
				fkConstraint.WithName(fkName)
			}
			refList := fkMap["references"].([]any)
			refMap := refList[0].(map[string]any)
			refTableId, err := sdk.ParseSchemaObjectIdentifier(refMap["table_id"].(string))
			if err != nil {
				return nil, fmt.Errorf("invalid references.table_id identifier: %w", err)
			}
			fkConstraint.WithForeignKey(sdk.OutOfLineForeignKey{
				TableName:   refTableId,
				ColumnNames: expandStringList(refMap["columns"].([]any)),
			})
			constraints = append(constraints, *fkConstraint)
		}
	}

	return constraints, nil
}

// ---------------------------------------------------------------------------
// CRUD
// ---------------------------------------------------------------------------

func CreateHybridTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schemaName := d.Get("schema").(string)
	databaseName := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	columnRequests, err := buildHybridTableColumnRequests(d.Get("column").([]any))
	if err != nil {
		return diag.FromErr(err)
	}

	constraints, err := buildOutOfLineConstraints(d)
	if err != nil {
		return diag.FromErr(err)
	}

	columnsAndConstraints := sdk.HybridTableColumnsConstraintsAndIndexesRequest{
		Columns:             columnRequests,
		OutOfLineConstraint: constraints,
	}
	request := sdk.NewCreateHybridTableRequest(id, columnsAndConstraints)

	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	if err := client.HybridTables.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating hybrid table %v: %w", id.FullyQualifiedName(), err))
	}

	// CRITICAL: Set ID immediately after CREATE, BEFORE any ALTER.
	// If ALTER fails, the table exists in Snowflake and must be in state
	// so the next apply can update it rather than fail on duplicate create.
	d.SetId(helpers.EncodeResourceIdentifier(id))

	// Set properties not available on CREATE (data_retention, max_data_extension).
	// CreateHybridTableOptions only has Comment *string, so apply parameters via a follow-up ALTER.
	set := sdk.NewHybridTableSetPropertiesRequest()
	if diags := handleHybridTableParametersCreate(d, set); diags.HasError() {
		return diags
	}
	if set.DataRetentionTimeInDays != nil || set.MaxDataExtensionTimeInDays != nil {
		if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting hybrid table properties %v: %w", id.FullyQualifiedName(), err))
		}
	}

	return GetReadHybridTableFunc(false)(ctx, d, meta)
}

func GetReadHybridTableFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		hybridTable, err := client.HybridTables.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query hybrid table. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Hybrid table id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		details, err := client.HybridTables.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		parameters, err := client.HybridTables.ShowParameters(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if diags := handleHybridTableParameterRead(d, parameters); diags.HasError() {
			return diags
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"comment", "comment", hybridTable.Comment, hybridTable.Comment, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.HybridTableToSchema(hybridTable)}),
			d.Set(DescribeOutputAttributeName, schemas.HybridTableDetailsListToSchema(details)),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("comment", hybridTable.Comment),
			d.Set("column", buildHybridColumnStateFromDescribe(details, d)),
			d.Set("primary_key", buildPrimaryKeyStateFromDescribe(details, d)),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateHybridTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle rename
	if d.HasChange("name") {
		databaseName := d.Get("database").(string)
		schemaName := d.Get("schema").(string)
		name := d.Get("name").(string)
		newId := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

		err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithNewName(newId))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	// Handle property changes using SET/UNSET pattern via shared helpers
	set := sdk.NewHybridTableSetPropertiesRequest()
	unset := sdk.NewHybridTableUnsetPropertiesRequest()

	if err := stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment); err != nil {
		return diag.FromErr(err)
	}
	if diags := handleHybridTableParametersChanges(d, set, unset); diags.HasError() {
		return diags
	}

	if set.Comment != nil || set.DataRetentionTimeInDays != nil || set.MaxDataExtensionTimeInDays != nil {
		if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting hybrid table properties %v: %w", id.FullyQualifiedName(), err))
		}
	}
	// Snowflake hybrid tables only support UNSET for one property at a time,
	// unlike regular tables which accept space-separated properties.
	// Issue separate ALTER ... UNSET for each property.
	if unset.Comment != nil {
		u := sdk.NewHybridTableUnsetPropertiesRequest()
		u.Comment = unset.Comment
		if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithUnset(*u)); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting hybrid table comment %v: %w", id.FullyQualifiedName(), err))
		}
	}
	if unset.DataRetentionTimeInDays != nil {
		u := sdk.NewHybridTableUnsetPropertiesRequest()
		u.DataRetentionTimeInDays = unset.DataRetentionTimeInDays
		if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithUnset(*u)); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting hybrid table data_retention_time_in_days %v: %w", id.FullyQualifiedName(), err))
		}
	}
	if unset.MaxDataExtensionTimeInDays != nil {
		u := sdk.NewHybridTableUnsetPropertiesRequest()
		u.MaxDataExtensionTimeInDays = unset.MaxDataExtensionTimeInDays
		if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithUnset(*u)); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting hybrid table max_data_extension_time_in_days %v: %w", id.FullyQualifiedName(), err))
		}
	}

	// Handle column changes
	if d.HasChange("column") {
		oldRaw, newRaw := d.GetChange("column")
		removed, added, changed := parseHybridColumns(oldRaw).diffs(parseHybridColumns(newRaw))

		// Drop removed columns
		if len(removed) > 0 {
			dropNames := make([]string, len(removed))
			for i, col := range removed {
				dropNames[i] = col.name
			}
			dropReq := sdk.NewHybridTableDropColumnActionRequest(dropNames).WithIfExists(true)
			if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithDropColumnAction(*dropReq)); err != nil {
				return diag.FromErr(fmt.Errorf("error dropping columns from hybrid table %v: %w", id.FullyQualifiedName(), err))
			}
		}

		// Add new columns (one at a time)
		for _, col := range added {
			dataType, err := datatypes.ParseDataType(col.dataType)
			if err != nil {
				return diag.FromErr(err)
			}
			addReq := sdk.NewHybridTableAddColumnActionRequest(col.name, sdk.DataType(col.dataType))
			if col.collate != "" {
				addReq.WithCollate(col.collate)
			}
			if col.comment != "" {
				addReq.WithComment(col.comment)
			}
			if col._default != nil {
				defMap := map[string]any{}
				if col._default.constant != nil {
					defMap["constant"] = *col._default.constant
				}
				if col._default.expression != nil {
					defMap["expression"] = *col._default.expression
				}
				if col._default.sequence != nil {
					defMap["sequence"] = *col._default.sequence
				}
				if defaultValue := buildHybridColumnDefaultValue(defMap, dataType); defaultValue != nil {
					addReq.WithDefaultValue(*defaultValue)
				}
			}
			if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithAddColumnAction(*addReq)); err != nil {
				return diag.FromErr(fmt.Errorf("error adding column %s to hybrid table %v: %w", col.name, id.FullyQualifiedName(), err))
			}
		}

		// Alter changed columns
		for _, ch := range changed {
			if ch.changedDataType {
				alterReq := sdk.NewHybridTableAlterColumnActionRequest(ch.newColumn.name).
					WithType(sdk.DataType(ch.newColumn.dataType))
				if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
					WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{*alterReq})); err != nil {
					return diag.FromErr(fmt.Errorf("error altering column type %s on hybrid table %v: %w",
						ch.newColumn.name, id.FullyQualifiedName(), err))
				}
			}

			if ch.changedComment {
				alterReq := sdk.NewHybridTableAlterColumnActionRequest(ch.newColumn.name)
				if ch.newColumn.comment == "" {
					alterReq.WithUnsetComment(true)
				} else {
					alterReq.WithComment(ch.newColumn.comment)
				}
				if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
					WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{*alterReq})); err != nil {
					return diag.FromErr(fmt.Errorf("error altering column comment %s on hybrid table %v: %w",
						ch.newColumn.name, id.FullyQualifiedName(), err))
				}
			}

			if ch.droppedDefault {
				alterReq := sdk.NewHybridTableAlterColumnActionRequest(ch.newColumn.name).
					WithDropDefault(true)
				if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
					WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{*alterReq})); err != nil {
					return diag.FromErr(fmt.Errorf("error dropping column default %s on hybrid table %v: %w",
						ch.newColumn.name, id.FullyQualifiedName(), err))
				}
			}
		}
	}

	return GetReadHybridTableFunc(false)(ctx, d, meta)
}

func ImportHybridTable(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, nil); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

// ---------------------------------------------------------------------------
// Read helpers
// ---------------------------------------------------------------------------

func buildHybridColumnStateFromDescribe(details []sdk.HybridTableDetails, d *schema.ResourceData) []any {
	flattened := make([]any, 0)

	// Build lookups from config for fields where DESCRIBE returns a normalized/different value
	// that would cause permanent drift against the user's config:
	// - type: DESCRIBE normalizes e.g. INTEGER -> NUMBER(38,0), which triggers HasChange("column")
	//   and cascading ComputedIfAnyAttributeChanged drift on describe_output/show_output.
	//   DiffSuppressDataTypes handles this at plan level but not for the non-refresh plan check.
	// - collate: DESCRIBE does not return collation at all.
	// - nullable: DESCRIBE returns NOT NULL for PK columns even when config says nullable=true.
	//   Preserve the config value since Snowflake silently enforces NOT NULL on PK columns.
	type configColumnInfo struct {
		dataType string
		collate  string
		nullable bool
	}
	configByName := make(map[string]configColumnInfo)
	if configCols, ok := d.GetOk("column"); ok {
		for _, rawCol := range configCols.([]any) {
			colMap := rawCol.(map[string]any)
			if colName, ok := colMap["name"].(string); ok {
				info := configColumnInfo{nullable: true}
				if dt, ok := colMap["type"].(string); ok {
					info.dataType = dt
				}
				if collate, ok := colMap["collate"].(string); ok {
					info.collate = collate
				}
				if nullable, ok := colMap["nullable"].(bool); ok {
					info.nullable = nullable
				}
				configByName[strings.ToUpper(colName)] = info
			}
		}
	}

	for _, td := range details {
		if td.Kind != "COLUMN" {
			continue
		}

		colInfo, found := configByName[strings.ToUpper(td.Name)]
		if !found {
			// Externally added column not in config: schema default is nullable=true.
			// Without this, zero-value would set nullable=false, which is ForceNew and
			// would produce a spurious delete+create plan instead of the desired Update.
			colInfo = configColumnInfo{nullable: true}
		}
		// Use the config type if it's semantically equivalent to the DESCRIBE type,
		// to avoid spurious diffs (e.g. config "INTEGER" vs DESCRIBE "NUMBER(38,0)").
		colType := td.Type
		if colInfo.dataType != "" {
			configDT, configErr := datatypes.ParseDataType(colInfo.dataType)
			describeDT, describeErr := datatypes.ParseDataType(td.Type)
			if configErr != nil {
				log.Printf("[WARN] hybrid_table Read: failed to parse config data type %q for column %q: %v", colInfo.dataType, td.Name, configErr)
			}
			if describeErr != nil {
				log.Printf("[WARN] hybrid_table Read: failed to parse describe data type %q for column %q: %v", td.Type, td.Name, describeErr)
			}
			if configErr == nil && describeErr == nil && datatypes.AreTheSame(configDT, describeDT) {
				colType = colInfo.dataType
			}
		}
		// Prefer the user's config value for collate (it preserves the exact spelling
		// the user wrote, e.g. "en-ci"). Fall back to the SDK-split Collation from
		// DESCRIBE for imported tables where no config exists.
		collate := colInfo.collate
		if collate == "" && td.Collation != nil {
			collate = *td.Collation
		}
		flat := map[string]any{
			"name":     td.Name,
			"type":     colType,
			"nullable": colInfo.nullable,
			"comment":  td.Comment,
			"collate":  collate,
		}

		if def := toHybridColumnDefaultConfig(td); def != nil {
			flat["default"] = def
		}

		flattened = append(flattened, flat)
	}
	return flattened
}

// buildPrimaryKeyStateFromDescribe reconstructs the primary_key block from DESCRIBE output.
// Called unconditionally on every Read (not just import) so that PK column drift is visible.
func buildPrimaryKeyStateFromDescribe(details []sdk.HybridTableDetails, d *schema.ResourceData) []map[string]any {
	pkKeys := make([]string, 0)
	for _, detail := range details {
		if detail.PrimaryKey {
			pkKeys = append(pkKeys, detail.Name)
		}
	}
	if len(pkKeys) == 0 {
		return nil
	}

	// Preserve constraint name from config if available
	pkName := ""
	if configPK, ok := d.GetOk("primary_key"); ok {
		pkList := configPK.([]any)
		if len(pkList) > 0 {
			pkMap := pkList[0].(map[string]any)
			if n, ok := pkMap["name"].(string); ok {
				pkName = n
			}
		}
	}

	return []map[string]any{
		{
			"name": pkName,
			"keys": pkKeys,
		},
	}
}

// toHybridColumnDefaultConfig converts HybridTableDetails.Default (string, not *string)
// into a config-compatible []any (single-element list or nil). Empty string means no default.
// This is NOT the same as table.go's toColumnDefaultConfig which checks == nil.
func toHybridColumnDefaultConfig(td sdk.HybridTableDetails) []any {
	if td.Default == "" {
		return nil
	}

	defaultRaw := td.Default

	// Sequence detection: ends with .NEXTVAL (case-insensitive)
	const nextvalSuffix = ".NEXTVAL"
	if strings.HasSuffix(strings.ToUpper(defaultRaw), nextvalSuffix) {
		sequenceIdRaw := defaultRaw[:len(defaultRaw)-len(nextvalSuffix)]
		id, err := sdk.ParseSchemaObjectIdentifier(sequenceIdRaw)
		if err != nil {
			log.Printf("[WARN] hybrid_table Read: failed to parse sequence identifier %q: %v", sequenceIdRaw, err)
			return []any{map[string]any{"sequence": sequenceIdRaw}}
		}
		return []any{map[string]any{"sequence": id.FullyQualifiedName()}}
	}

	// Expression detection: contains parentheses
	if strings.Contains(defaultRaw, "(") && strings.Contains(defaultRaw, ")") {
		return []any{map[string]any{
			"expression": defaultRaw,
		}}
	}

	// Constant: unescape for string types
	if sdk.IsStringType(td.Type) {
		return []any{map[string]any{
			"constant": snowflake.UnescapeSnowflakeString(defaultRaw),
		}}
	}

	return []any{map[string]any{
		"constant": defaultRaw,
	}}
}

// ---------------------------------------------------------------------------
// CustomizeDiff
// ---------------------------------------------------------------------------

// forceNewIfColumnCollateChanged forces recreation when collation changes on an
// existing column. The SDK's HybridTableAlterColumnActionRequest has no WithCollate
// method, so collation can only be set at column creation time.
func forceNewIfColumnCollateChanged() schema.CustomizeDiffFunc {
	return forceNewIfColumnFieldChanged(func(o, n hybridTableColumn) bool {
		return o.collate != n.collate
	})
}

// forceNewIfColumnNullableChanged forces recreation when nullable changes on an
// existing column. Hybrid tables do not support ALTER COLUMN SET/DROP NOT NULL,
// so toggling nullable requires recreation. Using a custom diff (rather than
// ForceNew on the schema field) ensures that adding a brand-new column does not
// spuriously trigger ForceNew when its nullable field initializes from the
// schema default.
func forceNewIfColumnNullableChanged() schema.CustomizeDiffFunc {
	return forceNewIfColumnFieldChanged(func(o, n hybridTableColumn) bool {
		return o.nullable != n.nullable
	})
}

func forceNewIfColumnFieldChanged(changed func(old, new hybridTableColumn) bool) schema.CustomizeDiffFunc {
	return func(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
		if !diff.HasChange("column") {
			return nil
		}
		oldRaw, newRaw := diff.GetChange("column")
		oldCols := parseHybridColumns(oldRaw)
		newCols := parseHybridColumns(newRaw)
		for _, o := range oldCols {
			for _, n := range newCols {
				if o.name == n.name && changed(o, n) {
					return diff.ForceNew("column")
				}
			}
		}
		return nil
	}
}
