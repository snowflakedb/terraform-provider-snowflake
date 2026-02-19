package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var hybridTableSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the hybrid table; must be unique for the schema in which the hybrid table is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "The database in which to create the hybrid table.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "The schema in which to create the hybrid table.",
		ValidateDiagFunc: IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"or_replace": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies whether to replace the hybrid table if it already exists.",
	},
	"column": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		Description: "Definitions of columns for the hybrid table.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Column name.",
				},
				"type": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Column data type.",
				},
				"nullable": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Specifies whether the column can contain NULL values. Default is true (nullable).",
				},
				"default": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Defines the default value for the column.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"expression": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Default value expression.",
							},
							"sequence": {
								Type:             schema.TypeString,
								Optional:         true,
								Description:      "Fully qualified name of sequence for default value.",
								ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
							},
						},
					},
				},
				"identity": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Defines the identity/autoincrement configuration for the column.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"start_num": {
								Type:        schema.TypeInt,
								Optional:    true,
								Default:     1,
								Description: "Starting value for the identity column.",
							},
							"step_num": {
								Type:        schema.TypeInt,
								Optional:    true,
								Default:     1,
								Description: "Step/increment value for the identity column.",
							},
						},
					},
				},
				"collate": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Collation specification for string column.",
				},
				"comment": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Column comment.",
				},
				"primary_key": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Specifies whether the column is a primary key (inline constraint).",
				},
				"unique": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Specifies whether the column has a unique constraint (inline constraint).",
				},
				"foreign_key": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Inline foreign key constraint for the column.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"table_name": {
								Type:             schema.TypeString,
								Required:         true,
								Description:      "Name of the table being referenced.",
								ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
							},
							"column_name": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Column name in the referenced table.",
							},
						},
					},
				},
			},
		},
	},
	"index": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Definitions of indexes for the hybrid table.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Index name.",
				},
				"columns": {
					Type:        schema.TypeList,
					Required:    true,
					MinItems:    1,
					Description: "List of column names to include in the index.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	},
	"primary_key": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Out-of-line primary key constraint definition.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Name of the primary key constraint.",
				},
				"columns": {
					Type:        schema.TypeList,
					Required:    true,
					MinItems:    1,
					Description: "List of column names forming the primary key.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	},
	"unique_constraint": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Out-of-line unique constraint definitions.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Name of the unique constraint.",
				},
				"columns": {
					Type:        schema.TypeList,
					Required:    true,
					MinItems:    1,
					Description: "List of column names forming the unique constraint.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	},
	"foreign_key": {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "Out-of-line foreign key constraint definitions.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Name of the foreign key constraint.",
				},
				"columns": {
					Type:        schema.TypeList,
					Required:    true,
					MinItems:    1,
					Description: "List of column names forming the foreign key.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"references_table": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Name of the table being referenced.",
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
				},
				"references_columns": {
					Type:        schema.TypeList,
					Required:    true,
					MinItems:    1,
					Description: "List of column names in the referenced table.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	},
	"data_retention_time_in_days": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "Specifies the retention period for the table in days.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the hybrid table.",
	},
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
		Description: "Outputs the result of `DESCRIBE HYBRID TABLE` for the given hybrid table.",
		Elem: &schema.Resource{
			Schema: schemas.HybridTableDescribeSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func HybridTable() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.HybridTables.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(
			string(previewfeatures.HybridTableResource),
			TrackingCreateWrapper(resources.HybridTable, CreateContextHybridTable)),
		ReadContext: PreviewFeatureReadContextWrapper(
			string(previewfeatures.HybridTableResource),
			TrackingReadWrapper(resources.HybridTable, ReadContextHybridTable(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(
			string(previewfeatures.HybridTableResource),
			TrackingUpdateWrapper(resources.HybridTable, UpdateContextHybridTable)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(
			string(previewfeatures.HybridTableResource),
			TrackingDeleteWrapper(resources.HybridTable, deleteFunc)),

		Description: "Resource used to manage hybrid table objects. For more information, check [hybrid table documentation](https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.HybridTable, customdiff.All(
			ComputedIfAnyAttributeChanged(hybridTableSchema, ShowOutputAttributeName, "name", "comment"),
			ComputedIfAnyAttributeChanged(hybridTableSchema, DescribeOutputAttributeName, "name", "column"),
			ComputedIfAnyAttributeChanged(hybridTableSchema, FullyQualifiedNameAttributeName, "name"),
			validatePrimaryKeyDefined,
			validateNoConflictingConstraints,
			validateColumnAttributes,
		)),

		Schema: hybridTableSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.HybridTable, ImportHybridTable),
		},
	}
}

func ImportHybridTable(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	if _, err := ImportName[sdk.SchemaObjectIdentifier](ctx, d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateContextHybridTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	// Build column definitions
	columns, err := expandHybridTableColumns(d.Get("column").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	// Build out-of-line constraints
	outOfLineConstraints, err := expandHybridTableOutOfLineConstraints(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Build out-of-line indexes
	outOfLineIndexes, err := expandHybridTableOutOfLineIndexes(d.Get("index").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	// Build the ColumnsAndConstraints structure
	columnsAndConstraints := sdk.HybridTableColumnsConstraintsAndIndexes{
		Columns:             columns,
		OutOfLineConstraint: outOfLineConstraints,
		OutOfLineIndex:      outOfLineIndexes,
	}

	// Create the request
	req := sdk.NewCreateHybridTableRequest(id, columnsAndConstraints)

	// Set OR REPLACE if specified
	if v, ok := d.GetOk("or_replace"); ok && v.(bool) {
		req.WithOrReplace(true)
	}

	// Set DATA_RETENTION_TIME_IN_DAYS if specified
	if v, ok := d.GetOk("data_retention_time_in_days"); ok {
		req.WithDataRetentionTimeInDays(v.(int))
	}

	// Set COMMENT if specified
	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	// Execute the CREATE command
	if err := client.HybridTables.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.FullyQualifiedName())

	return ReadContextHybridTable(false)(ctx, d, meta)
}

func ReadContextHybridTable(readFromProject bool) func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		// Fetch the hybrid table using ShowByIDSafely
		hybridTable, err := client.HybridTables.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to retrieve hybrid table. Marking the resource as removed.",
						Detail:   id.FullyQualifiedName(),
					},
				}
			}
			return diag.FromErr(err)
		}

		// Fetch detailed information using Describe
		details, err := client.HybridTables.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		// Build show_output and describe_output
		showOutput := schemas.HybridTableToSchema(hybridTable)
		describeOutput := make([]map[string]any, len(details))
		for i, detail := range details {
			describeOutput[i] = schemas.HybridTableDetailsToSchema(&detail)
		}

		// Set all state attributes using errors.Join
		// Note: We don't set name, database, schema - Terraform copies them from config automatically
		errs := errors.Join(
			d.Set("comment", hybridTable.Comment),
			d.Set(ShowOutputAttributeName, []map[string]any{showOutput}),
			d.Set(DescribeOutputAttributeName, describeOutput),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}

		return nil
	}
}

func UpdateContextHybridTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle updates to DATA_RETENTION_TIME_IN_DAYS or COMMENT
	if d.HasChange("data_retention_time_in_days") || d.HasChange("comment") {
		setReq := sdk.NewAlterHybridTableRequest(id)
		setProps := &sdk.HybridTableSetPropertiesRequest{}
		hasChanges := false

		if d.HasChange("data_retention_time_in_days") {
			if v, ok := d.GetOk("data_retention_time_in_days"); ok {
				setProps.DataRetentionTimeInDays = sdk.Int(v.(int))
				hasChanges = true
			} else {
				// Unset DATA_RETENTION_TIME_IN_DAYS
				unsetReq := sdk.NewAlterHybridTableRequest(id)
				unsetProps := &sdk.HybridTableUnsetPropertiesRequest{
					DataRetentionTimeInDays: sdk.Bool(true),
				}
				unsetReq.WithUnset(*unsetProps)
				if err := client.HybridTables.Alter(ctx, unsetReq); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if d.HasChange("comment") {
			if v, ok := d.GetOk("comment"); ok {
				setProps.Comment = sdk.String(v.(string))
				hasChanges = true
			} else {
				// Unset COMMENT
				unsetReq := sdk.NewAlterHybridTableRequest(id)
				unsetProps := &sdk.HybridTableUnsetPropertiesRequest{
					Comment: sdk.Bool(true),
				}
				unsetReq.WithUnset(*unsetProps)
				if err := client.HybridTables.Alter(ctx, unsetReq); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if hasChanges {
			setReq.WithSet(*setProps)
			if err := client.HybridTables.Alter(ctx, setReq); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Note: Changes to columns, indexes, and constraints would typically require
	// more complex alter operations or recreation. For the initial implementation,
	// these are marked as ForceNew in the schema, so Terraform will handle recreation.
	// Advanced column/index management can be added in future iterations.

	return ReadContextHybridTable(true)(ctx, d, meta)
}

// Helper functions for model conversions

// expandHybridTableColumns converts Terraform column schema to SDK column structures
func expandHybridTableColumns(columns []interface{}) ([]sdk.HybridTableColumn, error) {
	result := make([]sdk.HybridTableColumn, len(columns))
	for i, col := range columns {
		colMap := col.(map[string]interface{})

		name := colMap["name"].(string)
		dataTypeStr := colMap["type"].(string)
		dataType := sdk.DataType(dataTypeStr)

		column := sdk.HybridTableColumn{
			Name: name,
			Type: dataType,
		}

		// Handle nullable
		if nullable, ok := colMap["nullable"].(bool); ok && !nullable {
			column.NotNull = sdk.Bool(true)
		}

		// Handle default value
		if defaultList, ok := colMap["default"].([]interface{}); ok && len(defaultList) > 0 {
			defaultMap := defaultList[0].(map[string]interface{})
			if expr, ok := defaultMap["expression"].(string); ok && expr != "" {
				column.DefaultValue = &sdk.ColumnDefaultValue{
					Expression: sdk.String(expr),
				}
			}
			// Note: Sequence default is not currently supported in hybrid tables
			// based on the SDK implementation
		}

		// Handle identity
		if identityList, ok := colMap["identity"].([]interface{}); ok && len(identityList) > 0 {
			identityMap := identityList[0].(map[string]interface{})
			startNum := identityMap["start_num"].(int)
			stepNum := identityMap["step_num"].(int)
			column.DefaultValue = &sdk.ColumnDefaultValue{
				Identity: &sdk.ColumnIdentity{
					Start:     startNum,
					Increment: stepNum,
				},
			}
		}

		// Handle collate
		if collate, ok := colMap["collate"].(string); ok && collate != "" {
			column.Collate = sdk.String(collate)
		}

		// Handle comment
		if comment, ok := colMap["comment"].(string); ok && comment != "" {
			column.Comment = sdk.String(comment)
		}

		// Handle inline constraints (primary_key, unique, foreign_key)
		if pkFlag, ok := colMap["primary_key"].(bool); ok && pkFlag {
			column.InlineConstraint = &sdk.HybridTableColumnInlineConstraint{
				Type: sdk.ColumnConstraintTypePrimaryKey,
			}
		} else if uniqueFlag, ok := colMap["unique"].(bool); ok && uniqueFlag {
			column.InlineConstraint = &sdk.HybridTableColumnInlineConstraint{
				Type: sdk.ColumnConstraintTypeUnique,
			}
		} else if fkList, ok := colMap["foreign_key"].([]interface{}); ok && len(fkList) > 0 {
			fkMap := fkList[0].(map[string]interface{})
			tableName := fkMap["table_name"].(string)
			columnName := fkMap["column_name"].(string)

			column.InlineConstraint = &sdk.HybridTableColumnInlineConstraint{
				Type: sdk.ColumnConstraintTypeForeignKey,
				ForeignKey: &sdk.InlineForeignKey{
					TableName:  tableName,
					ColumnName: []string{columnName},
				},
			}
		}

		result[i] = column
	}
	return result, nil
}

// expandHybridTableOutOfLineConstraints converts Terraform constraint schemas to SDK constraint structures
func expandHybridTableOutOfLineConstraints(d *schema.ResourceData) ([]sdk.HybridTableOutOfLineConstraint, error) {
	var result []sdk.HybridTableOutOfLineConstraint

	// Handle primary_key out-of-line constraint
	if pkList, ok := d.GetOk("primary_key"); ok {
		pkListTyped := pkList.([]interface{})
		if len(pkListTyped) > 0 {
			pkMap := pkListTyped[0].(map[string]interface{})
			constraint := sdk.HybridTableOutOfLineConstraint{
				Type:    sdk.ColumnConstraintTypePrimaryKey,
				Columns: expandStringList(pkMap["columns"].([]interface{})),
			}
			if name, ok := pkMap["name"].(string); ok && name != "" {
				constraint.Name = sdk.String(name)
			}
			result = append(result, constraint)
		}
	}

	// Handle unique_constraint out-of-line constraints
	if uniqueSet, ok := d.GetOk("unique_constraint"); ok {
		uniqueList := uniqueSet.(*schema.Set).List()
		for _, item := range uniqueList {
			ucMap := item.(map[string]interface{})
			constraint := sdk.HybridTableOutOfLineConstraint{
				Type:    sdk.ColumnConstraintTypeUnique,
				Columns: expandStringList(ucMap["columns"].([]interface{})),
			}
			if name, ok := ucMap["name"].(string); ok && name != "" {
				constraint.Name = sdk.String(name)
			}
			result = append(result, constraint)
		}
	}

	// Handle foreign_key out-of-line constraints
	if fkSet, ok := d.GetOk("foreign_key"); ok {
		fkList := fkSet.(*schema.Set).List()
		for _, item := range fkList {
			fkMap := item.(map[string]interface{})
			refTable := fkMap["references_table"].(string)
			refTableId, err := sdk.ParseSchemaObjectIdentifier(refTable)
			if err != nil {
				return nil, err
			}

			constraint := sdk.HybridTableOutOfLineConstraint{
				Type:    sdk.ColumnConstraintTypeForeignKey,
				Columns: expandStringList(fkMap["columns"].([]interface{})),
				ForeignKey: &sdk.OutOfLineForeignKey{
					TableName:   refTableId,
					ColumnNames: expandStringList(fkMap["references_columns"].([]interface{})),
				},
			}
			if name, ok := fkMap["name"].(string); ok && name != "" {
				constraint.Name = sdk.String(name)
			}
			result = append(result, constraint)
		}
	}

	return result, nil
}

// expandHybridTableOutOfLineIndexes converts Terraform index schema to SDK index structures
func expandHybridTableOutOfLineIndexes(indexes []interface{}) ([]sdk.HybridTableOutOfLineIndex, error) {
	result := make([]sdk.HybridTableOutOfLineIndex, len(indexes))
	for i, idx := range indexes {
		idxMap := idx.(map[string]interface{})
		index := sdk.HybridTableOutOfLineIndex{
			Name:    idxMap["name"].(string),
			Columns: expandStringList(idxMap["columns"].([]interface{})),
		}
		result[i] = index
	}
	return result, nil
}

// validatePrimaryKeyDefined ensures that a primary key is defined either inline or out-of-line.
// Snowflake hybrid tables require a primary key, and it must be defined in exactly one way.
func validatePrimaryKeyDefined(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	columns := d.Get("column").([]interface{})
	primaryKeyBlock := d.Get("primary_key").([]interface{})

	// Count inline primary keys
	inlinePKCount := 0
	var inlinePKColumns []string
	for _, col := range columns {
		colMap := col.(map[string]interface{})
		colName := colMap["name"].(string)
		if pk, ok := colMap["primary_key"].(bool); ok && pk {
			inlinePKCount++
			inlinePKColumns = append(inlinePKColumns, colName)
		}
	}

	hasOutOfLinePK := len(primaryKeyBlock) > 0

	// Error if both inline and out-of-line primary keys are defined
	if inlinePKCount > 0 && hasOutOfLinePK {
		return errors.New("primary key cannot be defined both inline (column.primary_key = true) and out-of-line (primary_key block); use only one method")
	}

	// Error if no primary key is defined
	if inlinePKCount == 0 && !hasOutOfLinePK {
		return errors.New("hybrid table requires a primary key; define either column.primary_key = true for a single column, or use a primary_key block for single or composite keys")
	}

	// Error if multiple columns have inline primary key
	if inlinePKCount > 1 {
		return fmt.Errorf("only one column can be marked as inline primary key; found %d columns with primary_key = true: %s (use primary_key block for composite keys)", inlinePKCount, strings.Join(inlinePKColumns, ", "))
	}

	return nil
}

// validateNoConflictingConstraints ensures that constraints are not defined both inline and out-of-line for the same column.
func validateNoConflictingConstraints(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	columns := d.Get("column").([]interface{})

	// Build maps of column names with inline constraints
	inlineFKColumns := make(map[string]bool)
	inlineUniqueColumns := make(map[string]bool)

	for _, col := range columns {
		colMap := col.(map[string]interface{})
		name := colMap["name"].(string)

		// Track columns with inline foreign keys
		if fkList, ok := colMap["foreign_key"].([]interface{}); ok && len(fkList) > 0 {
			inlineFKColumns[name] = true
		}

		// Track columns with inline unique constraints
		if unique, ok := colMap["unique"].(bool); ok && unique {
			inlineUniqueColumns[name] = true
		}
	}

	// Check out-of-line foreign keys don't conflict with inline ones
	if fkSet, ok := d.GetOk("foreign_key"); ok {
		outOfLineFKs := fkSet.(*schema.Set).List()
		for _, fk := range outOfLineFKs {
			fkMap := fk.(map[string]interface{})
			fkColumns := fkMap["columns"].([]interface{})
			for _, col := range fkColumns {
				colName := col.(string)
				if inlineFKColumns[colName] {
					return fmt.Errorf("column %q has both inline (column.foreign_key) and out-of-line (foreign_key block) foreign key definitions; use only one method", colName)
				}
			}
		}
	}

	// Check out-of-line unique constraints don't conflict with inline ones
	if ucSet, ok := d.GetOk("unique_constraint"); ok {
		outOfLineUCs := ucSet.(*schema.Set).List()
		for _, uc := range outOfLineUCs {
			ucMap := uc.(map[string]interface{})
			ucColumns := ucMap["columns"].([]interface{})
			for _, col := range ucColumns {
				colName := col.(string)
				if inlineUniqueColumns[colName] {
					return fmt.Errorf("column %q has both inline (column.unique = true) and out-of-line (unique_constraint block) unique constraint definitions; use only one method", colName)
				}
			}
		}
	}

	return nil
}

// validateColumnAttributes validates that column attribute combinations are valid.
// For example, a column cannot have both default and identity, as these are mutually exclusive.
func validateColumnAttributes(ctx context.Context, d *schema.ResourceDiff, meta any) error {
	columns := d.Get("column").([]interface{})

	for i, col := range columns {
		colMap := col.(map[string]interface{})
		name := colMap["name"].(string)

		// Check if both default and identity are specified (mutually exclusive)
		hasDefault := false
		if defaultList, ok := colMap["default"].([]interface{}); ok && len(defaultList) > 0 {
			hasDefault = true
		}

		hasIdentity := false
		if identityList, ok := colMap["identity"].([]interface{}); ok && len(identityList) > 0 {
			hasIdentity = true
		}

		if hasDefault && hasIdentity {
			return fmt.Errorf("column %q (index %d) cannot have both default and identity; these attributes are mutually exclusive", name, i)
		}
	}

	return nil
}
