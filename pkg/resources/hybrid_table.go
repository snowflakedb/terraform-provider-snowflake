package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var hybridTableSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies the identifier for the hybrid table; must be unique for the database and schema in which the table is created.",
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database in which to create the hybrid table.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema in which to create the hybrid table.",
	},
	"column": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		Description: "Definitions of columns to create in the hybrid table. Minimum one required. Column structure changes (add/remove/rename columns) require table recreation, but comment changes can be updated in-place.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Column name",
					DiffSuppressFunc: suppressColumnNameDiff,
				},
				"type": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Column type, e.g. NUMBER, VARCHAR. For a full list of column types, see [Summary of Data Types](https://docs.snowflake.com/en/sql-reference/intro-summary-data-types).",
					ValidateDiagFunc: IsDataTypeValid,
					DiffSuppressFunc: DiffSuppressDataTypes,
				},
				"nullable": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Whether this column can contain null values. Set to false for columns used in primary key constraints.",
				},
				"comment": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Column comment",
				},
			},
		},
	},
	"constraint": {
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		ForceNew:    true,
		Description: "Definitions of constraints for the hybrid table. At least one PRIMARY KEY constraint is required. Constraints cannot be modified after table creation.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the constraint",
				},
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "Type of constraint: PRIMARY KEY, FOREIGN KEY, or UNIQUE. All constraints must be ENFORCED.",
					ValidateFunc: validation.StringInSlice([]string{"PRIMARY KEY", "FOREIGN KEY", "UNIQUE"}, false),
				},
				"columns": {
					Type:        schema.TypeList,
					Required:    true,
					MinItems:    1,
					Description: "Columns to use in the constraint",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"foreign_key": {
					Type:        schema.TypeList,
					Optional:    true,
					MaxItems:    1,
					Description: "Foreign key reference details. Required when type is FOREIGN KEY.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"table_id": {
								Type:             schema.TypeString,
								Required:         true,
								Description:      "Identifier of the referenced table in the format database.schema.table",
								ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
							},
							"columns": {
								Type:        schema.TypeList,
								Required:    true,
								MinItems:    1,
								Description: "Columns in the referenced table",
								Elem:        &schema.Schema{Type: schema.TypeString},
							},
							"match": {
								Type:         schema.TypeString,
								Optional:     true,
								Description:  "The match type for the foreign key: FULL, SIMPLE, or PARTIAL. Note: MATCH is not supported for hybrid tables and will be ignored.",
								ValidateFunc: validation.StringInSlice([]string{"FULL", "SIMPLE", "PARTIAL"}, true),
							},
							"on_update": {
								Type:         schema.TypeString,
								Optional:     true,
								Description:  "Action to perform when the primary/unique key is updated: CASCADE, SET NULL, SET DEFAULT, RESTRICT, or NO ACTION. Note: not supported for hybrid tables and will be ignored.",
								ValidateFunc: validation.StringInSlice([]string{"CASCADE", "SET NULL", "SET DEFAULT", "RESTRICT", "NO ACTION"}, true),
							},
							"on_delete": {
								Type:         schema.TypeString,
								Optional:     true,
								Description:  "Action to perform when the primary/unique key is deleted: CASCADE, SET NULL, SET DEFAULT, RESTRICT, or NO ACTION. Note: not supported for hybrid tables and will be ignored.",
								ValidateFunc: validation.StringInSlice([]string{"CASCADE", "SET NULL", "SET DEFAULT", "RESTRICT", "NO ACTION"}, true),
							},
						},
					},
				},
			},
		},
	},
	"index": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Definitions of secondary indexes for the hybrid table",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Name of the index",
					DiffSuppressFunc: suppressIdentifierQuoting,
				},
				"columns": {
					Type:        schema.TypeList,
					Required:    true,
					MinItems:    1,
					Description: "Columns to include in the index",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the hybrid table",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW TABLES` for the given hybrid table.",
		Elem: &schema.Resource{
			Schema: schemas.ShowTableSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE TABLE COLUMNS` for the given hybrid table.",
		Elem: &schema.Resource{
			Schema: schemas.TableDescribeSchema,
		},
	},
}

func HybridTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.HybridTable, CreateHybridTable),
		ReadContext:   TrackingReadWrapper(resources.HybridTable, ReadHybridTable),
		UpdateContext: TrackingUpdateWrapper(resources.HybridTable, UpdateHybridTable),
		DeleteContext: TrackingDeleteWrapper(resources.HybridTable, DeleteHybridTable),

		Schema: hybridTableSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: defaultTimeouts,

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(hybridTableSchema, ShowOutputAttributeName, "comment"),
			ComputedIfAnyAttributeChanged(hybridTableSchema, FullyQualifiedNameAttributeName, "name"),
			validateHybridTableConstraintColumns,
			validateHybridTableIndexColumns,
			forceNewOnColumnStructureChange,
		),
	}
}

func CreateHybridTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	// Parse columns
	columnRequests, err := getHybridTableColumnRequests(d.Get("column").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	// Parse constraints
	constraintRequests, err := getHybridTableConstraintRequests(d.Get("constraint").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	// Create hybrid table request
	createRequest := sdk.NewCreateHybridTableRequest(id, columnRequests, constraintRequests)

	if v, ok := d.GetOk("comment"); ok {
		createRequest.WithComment(sdk.String(v.(string)))
	}

	// Create hybrid table
	err = client.Tables.CreateHybridTable(ctx, createRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating hybrid table %v: %w", name, err))
	}

	// Create indexes (retry logic for "Another index is being built" is handled in SDK)
	if v, ok := d.GetOk("index"); ok {
		indexes := v.([]interface{})
		for _, idx := range indexes {
			indexMap, ok := idx.(map[string]interface{})
			if !ok {
				return diag.Errorf("unexpected type for index configuration")
			}
			indexName, ok := indexMap["name"].(string)
			if !ok {
				return diag.Errorf("index name must be a string")
			}
			rawIndexColumns := expandStringList(indexMap["columns"].([]interface{}))

			// Use selective quoting - must match how columns were defined
			indexColumns := quoteColumnNamesSelectively(rawIndexColumns)

			// Create index identifier (SDK will use only the name part in SQL)
			indexId := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, indexName)
			indexRequest := sdk.NewCreateIndexRequest(indexId, id, indexColumns)

			err := client.Tables.CreateIndex(ctx, indexRequest)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error creating index %v on hybrid table %v: %w", indexName, name, err))
			}
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadHybridTable(ctx, d, meta)
}

func getHybridTableColumnRequests(from interface{}) ([]sdk.TableColumnRequest, error) {
	cols, ok := from.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for columns configuration")
	}
	requests := make([]sdk.TableColumnRequest, len(cols))

	for i, c := range cols {
		colMap, ok := c.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("column[%d]: unexpected type for column configuration", i)
		}
		colName, ok := colMap["name"].(string)
		if !ok {
			return nil, fmt.Errorf("column[%d]: name must be a string", i)
		}
		// Use selective quoting - only quote when necessary
		columnName := quoteIdentifierIfNeeded(colName)
		columnType, ok := colMap["type"].(string)
		if !ok {
			return nil, fmt.Errorf("column[%d]: type must be a string", i)
		}

		request := sdk.NewTableColumnRequest(columnName, sdk.DataType(columnType))

		// Set nullable
		if nullable, ok := colMap["nullable"].(bool); ok {
			request.WithNotNull(sdk.Bool(!nullable))
		}

		// Set comment
		if comment, ok := colMap["comment"].(string); ok && comment != "" {
			request.WithComment(sdk.String(comment))
		}

		requests[i] = *request
	}

	return requests, nil
}

func getHybridTableConstraintRequests(from interface{}) ([]sdk.OutOfLineConstraintRequest, error) {
	constraints, ok := from.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for constraints configuration")
	}
	requests := make([]sdk.OutOfLineConstraintRequest, len(constraints))

	for i, c := range constraints {
		constraintMap, ok := c.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("constraint[%d]: unexpected type for constraint configuration", i)
		}
		constraintName, ok := constraintMap["name"].(string)
		if !ok {
			return nil, fmt.Errorf("constraint[%d]: name must be a string", i)
		}
		constraintType, ok := constraintMap["type"].(string)
		if !ok {
			return nil, fmt.Errorf("constraint[%d]: type must be a string", i)
		}
		rawColumns := expandStringList(constraintMap["columns"].([]interface{}))
		// Use selective quoting - must match how columns were defined
		columns := quoteColumnNamesSelectively(rawColumns)

		var sdkConstraintType sdk.ColumnConstraintType
		switch constraintType {
		case "PRIMARY KEY":
			sdkConstraintType = sdk.ColumnConstraintTypePrimaryKey
		case "FOREIGN KEY":
			sdkConstraintType = sdk.ColumnConstraintTypeForeignKey
		case "UNIQUE":
			sdkConstraintType = sdk.ColumnConstraintTypeUnique
		default:
			return nil, fmt.Errorf("unsupported constraint type: %s", constraintType)
		}

		request := sdk.NewOutOfLineConstraintRequest(sdkConstraintType).
			WithName(sdk.String(constraintName)).
			WithColumns(columns)

		// Handle foreign key
		if constraintType == "FOREIGN KEY" {
			if fk, ok := constraintMap["foreign_key"].([]interface{}); ok && len(fk) > 0 {
				fkMap, ok := fk[0].(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("constraint[%d]: unexpected type for foreign_key configuration", i)
				}
				tableIdStr, ok := fkMap["table_id"].(string)
				if !ok {
					return nil, fmt.Errorf("constraint[%d]: foreign_key.table_id must be a string", i)
				}
				tableId := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(tableIdStr)
				rawFkColumns := expandStringList(fkMap["columns"].([]interface{}))
				// Use selective quoting - must match how columns were defined
				fkColumns := quoteColumnNamesSelectively(rawFkColumns)

				fkRequest := sdk.NewOutOfLineForeignKeyRequest(tableId, fkColumns)

				// Note: MATCH clause and ON UPDATE/ON DELETE actions are not supported for hybrid tables

				request.WithForeignKey(fkRequest)
			}
		}

		requests[i] = *request
	}

	return requests, nil
}

func ReadHybridTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.SchemaObjectIdentifier)

	// Show hybrid table
	table, err := client.Tables.ShowByIDSafely(ctx, id)
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

	// Verify it's a table (Snowflake returns "TABLE" for both regular and hybrid tables)
	// We rely on the fact that only hybrid tables can have certain constraints/properties
	if table.Kind != "TABLE" {
		return diag.FromErr(fmt.Errorf("expected TABLE but got %s", table.Kind))
	}

	// Set basic attributes
	if err := d.Set("name", table.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name for hybrid table %s: %w", id.FullyQualifiedName(), err))
	}
	if err := d.Set("database", table.DatabaseName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set database for hybrid table %s: %w", id.FullyQualifiedName(), err))
	}
	if err := d.Set("schema", table.SchemaName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set schema for hybrid table %s: %w", id.FullyQualifiedName(), err))
	}
	if err := d.Set("comment", table.Comment); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set comment for hybrid table %s: %w", id.FullyQualifiedName(), err))
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set fully_qualified_name for hybrid table %s: %w", id.FullyQualifiedName(), err))
	}

	// Build a map of configured column names (preserve user's casing)
	// This is important because Snowflake may return names in normalized form
	configColumns := d.Get("column").([]interface{})
	configColumnNames := make(map[string]string) // map normalized name to original casing
	for _, c := range configColumns {
		colMap := c.(map[string]interface{})
		colName := colMap["name"].(string)
		// Map the normalized form to the original name
		normalizedName := NormalizeIdentifier(colName)
		configColumnNames[normalizedName] = colName
		// Also map the exact name
		configColumnNames[colName] = colName
	}

	// Get column details
	columnDetails, err := client.Tables.DescribeColumns(ctx, sdk.NewDescribeTableColumnsRequest(id))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to describe columns for hybrid table %s: %w", id.FullyQualifiedName(), err))
	}

	// Set columns
	columns := make([]map[string]interface{}, 0)
	for _, col := range columnDetails {
		if col.Kind != "COLUMN" {
			continue
		}

		// Normalize the data type to match input format
		dataType, err := datatypes.ParseDataType(string(col.Type))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error parsing data type %s: %w", col.Type, err))
		}
		// Remove spaces after commas to match input format (e.g., "NUMBER(38,0)" not "NUMBER(38, 0)")
		normalizedType := strings.ReplaceAll(dataType.ToSql(), ", ", ",")

		// Strip default precision from TIMESTAMP types to match user input
		// e.g., "TIMESTAMP_NTZ(9)" -> "TIMESTAMP_NTZ" if user didn't specify precision
		if strings.HasPrefix(normalizedType, "TIMESTAMP_") && strings.Contains(normalizedType, "(9)") {
			normalizedType = strings.Replace(normalizedType, "(9)", "", 1)
		}

		// Determine the column name to store in state
		// Priority: 1) Use configured name if available, 2) Infer from Snowflake's name
		columnName := col.Name
		if len(configColumnNames) > 0 {
			// During normal read/refresh: use configured name
			if configName, ok := configColumnNames[col.Name]; ok {
				columnName = configName
			} else if configName, ok := configColumnNames[strings.ToUpper(col.Name)]; ok {
				columnName = configName
			}
		} else {
			// During import (config is empty): infer original name
			// If the returned name is all uppercase and contains only alphanumeric + underscore,
			// it was likely a simple unquoted identifier - convert to lowercase for user-friendliness
			isSimpleUppercase := true
			for _, r := range col.Name {
				if r == '_' || (r >= '0' && r <= '9') {
					continue
				}
				if r < 'A' || r > 'Z' {
					isSimpleUppercase = false
					break
				}
			}
			if isSimpleUppercase && len(col.Name) > 0 {
				// Likely was a simple identifier - use lowercase
				columnName = strings.ToLower(col.Name)
			}
			// Otherwise keep the name as Snowflake returned it (mixed case, special chars were quoted)
		}

		column := map[string]interface{}{
			"name":     columnName,
			"type":     normalizedType,
			"nullable": col.IsNullable,
		}

		if col.Comment != nil {
			column["comment"] = *col.Comment
		}

		columns = append(columns, column)
	}
	if err := d.Set("column", columns); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set column for hybrid table %s: %w", id.FullyQualifiedName(), err))
	}

	// Note: Constraints are ForceNew, so we don't need to read them back
	// They are stored in state and cannot be modified

	// Get indexes
	showIndexesRequest := &sdk.ShowIndexesRequest{
		In: &sdk.In{
			Table: id,
		},
	}
	indexes, err := client.Tables.ShowIndexes(ctx, showIndexesRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to list indexes for hybrid table %s: %w", id.FullyQualifiedName(), err))
	}

	// Set indexes - only include explicitly user-defined indexes, in the same order as config
	// Snowflake auto-creates indexes for PRIMARY KEY and UNIQUE constraints, which we should exclude
	configIndexes := d.Get("index").([]interface{})
	configIndexNames := make(map[string]int) // map name to order
	for i, idx := range configIndexes {
		indexMap := idx.(map[string]interface{})
		indexName := indexMap["name"].(string)
		configIndexNames[strings.ToUpper(indexName)] = i
	}

	// Create a map of indexes from Snowflake
	snowflakeIndexes := make(map[string]*sdk.Index)
	for i := range indexes {
		snowflakeIndexes[strings.ToUpper(indexes[i].Name)] = &indexes[i]
	}

	// Build index list in config order
	indexList := make([]map[string]interface{}, len(configIndexes))
	for name, order := range configIndexNames {
		idx, ok := snowflakeIndexes[name]
		if !ok {
			// Index not found in Snowflake - use config values
			indexList[order] = configIndexes[order].(map[string]interface{})
			continue
		}

		// Unquote column names to match input format
		var columns []string
		if len(idx.Columns) > 0 {
			columns = make([]string, len(idx.Columns))
			for i, col := range idx.Columns {
				columns[i] = strings.Trim(col, `"`)
			}
		} else {
			// SHOW INDEXES may not return columns for hybrid tables - use config values
			configIdx := configIndexes[order].(map[string]interface{})
			configColumns := configIdx["columns"].([]interface{})
			columns = make([]string, len(configColumns))
			for i, col := range configColumns {
				columns[i] = col.(string)
			}
		}

		// Use config name to preserve user's specified case
		configIdx := configIndexes[order].(map[string]interface{})
		indexList[order] = map[string]interface{}{
			"name":    configIdx["name"].(string),
			"columns": columns,
		}
	}

	if err := d.Set("index", indexList); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set index for hybrid table %s: %w", id.FullyQualifiedName(), err))
	}

	// Set show_output
	if err := d.Set(ShowOutputAttributeName, []map[string]any{schemas.TableToSchema(table)}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set show_output for hybrid table %s: %w", id.FullyQualifiedName(), err))
	}

	// Set describe_output
	if err := d.Set(DescribeOutputAttributeName, schemas.TableDescriptionToSchema(columnDetails)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set describe_output for hybrid table %s: %w", id.FullyQualifiedName(), err))
	}

	return nil
}

func UpdateHybridTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.SchemaObjectIdentifier)

	// Handle comment changes
	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		alterRequest := sdk.NewAlterTableRequest(id)
		if comment != "" {
			alterRequest.WithSet(sdk.NewTableSetRequest().WithComment(sdk.String(comment)))
		} else {
			alterRequest.WithUnset(sdk.NewTableUnsetRequest().WithComment(true))
		}

		err := client.Tables.Alter(ctx, alterRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating hybrid table comment: %w", err))
		}
	}

	// Handle column comment changes (only comments can be updated, not structure)
	if d.HasChange("column") {
		oldCols, newCols := d.GetChange("column")
		oldColsList := oldCols.([]interface{})
		newColsList := newCols.([]interface{})

		// If the number of columns changed, this should force recreation
		// But comment-only changes can be handled
		if len(oldColsList) == len(newColsList) {
			for i := range newColsList {
				oldColMap := oldColsList[i].(map[string]interface{})
				newColMap := newColsList[i].(map[string]interface{})

				oldName := oldColMap["name"].(string)
				newName := newColMap["name"].(string)

				// Only handle comment changes for the same column
				if oldName == newName {
					oldComment, oldCommentOk := oldColMap["comment"].(string)
					newComment, newCommentOk := newColMap["comment"].(string)

					if oldCommentOk || newCommentOk {
						if oldComment != newComment {
							// Update column comment via ALTER TABLE
							columnName := quoteIdentifierIfNeeded(newName)

							if newComment != "" {
								// Set comment
								alterSQL := fmt.Sprintf("ALTER COLUMN %s COMMENT '%s'", columnName, strings.ReplaceAll(newComment, "'", "''"))
								_, err := client.ExecForTests(ctx, fmt.Sprintf("ALTER TABLE %s %s", id.FullyQualifiedName(), alterSQL))
								if err != nil {
									return diag.FromErr(fmt.Errorf("error updating column comment for %s: %w", columnName, err))
								}
							} else {
								// Unset comment
								alterSQL := fmt.Sprintf("ALTER COLUMN %s UNSET COMMENT", columnName)
								_, err := client.ExecForTests(ctx, fmt.Sprintf("ALTER TABLE %s %s", id.FullyQualifiedName(), alterSQL))
								if err != nil {
									return diag.FromErr(fmt.Errorf("error unsetting column comment for %s: %w", columnName, err))
								}
							}
						}
					}
				}
			}
		}
	}

	// Handle index changes
	if d.HasChange("index") {
		oldIndexes, newIndexes := d.GetChange("index")
		oldIndexList := oldIndexes.([]interface{})
		newIndexList := newIndexes.([]interface{})

		// Build maps for comparison
		oldIndexMap := make(map[string][]string)
		for _, idx := range oldIndexList {
			indexMap := idx.(map[string]interface{})
			indexName := indexMap["name"].(string)
			oldIndexMap[indexName] = expandStringList(indexMap["columns"].([]interface{}))
		}

		newIndexMap := make(map[string][]string)
		for _, idx := range newIndexList {
			indexMap := idx.(map[string]interface{})
			indexName := indexMap["name"].(string)
			newIndexMap[indexName] = expandStringList(indexMap["columns"].([]interface{}))
		}

		// Drop removed indexes
		for indexName := range oldIndexMap {
			if _, exists := newIndexMap[indexName]; !exists {
				indexId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), indexName)
				dropRequest := sdk.NewDropIndexRequest(indexId).WithIfExists(sdk.Bool(true))
				err := client.Tables.DropIndex(ctx, dropRequest)
				if err != nil {
					// If the error is about table not existing, it might be during table recreation
					// or an SDK issue with identifier formatting - skip if IF EXISTS is set
					if strings.Contains(err.Error(), "does not exist") {
						// Index already gone, continue
						continue
					}
					return diag.FromErr(fmt.Errorf("error dropping index %v: %w", indexName, err))
				}
			}
		}

		// Create new indexes
		for indexName, columns := range newIndexMap {
			if _, exists := oldIndexMap[indexName]; !exists {
				// Apply selective quoting to match column definitions
				quotedColumns := quoteColumnNamesSelectively(columns)
				indexId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), indexName)
				createRequest := sdk.NewCreateIndexRequest(indexId, id, quotedColumns)
				err := client.Tables.CreateIndex(ctx, createRequest)
				if err != nil {
					return diag.FromErr(fmt.Errorf("error creating index %v: %w", indexName, err))
				}
			}
		}
	}

	return ReadHybridTable(ctx, d, meta)
}

func DeleteHybridTable(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.SchemaObjectIdentifier)

	dropRequest := sdk.NewDropTableRequest(id).WithIfExists(sdk.Bool(true))
	err := client.Tables.Drop(ctx, dropRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error dropping hybrid table %v: %w", id.FullyQualifiedName(), err))
	}

	d.SetId("")
	return nil
}

// validateHybridTableConstraintColumns validates that all columns referenced in constraints exist in the column list.
func validateHybridTableConstraintColumns(_ context.Context, d *schema.ResourceDiff, _ any) error {
	columns := d.Get("column").([]interface{})
	constraints := d.Get("constraint").([]interface{})

	// Build a set of valid column names
	validColumns := make(map[string]bool)
	for _, c := range columns {
		colMap, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := colMap["name"].(string); ok {
			validColumns[strings.ToUpper(name)] = true
		}
	}

	// Check each constraint's columns
	for i, constraint := range constraints {
		constraintMap, ok := constraint.(map[string]interface{})
		if !ok {
			continue
		}

		constraintName := ""
		if name, ok := constraintMap["name"].(string); ok {
			constraintName = name
		}

		constraintColumns, ok := constraintMap["columns"].([]interface{})
		if !ok {
			continue
		}

		for _, col := range constraintColumns {
			colName, ok := col.(string)
			if !ok {
				continue
			}
			if !validColumns[strings.ToUpper(colName)] {
				return fmt.Errorf("constraint[%d] (%s) references non-existent column %q; valid columns are: %v",
					i, constraintName, colName, getColumnNames(columns))
			}
		}
	}

	return nil
}

// validateHybridTableIndexColumns validates that all columns referenced in indexes exist in the column list.
func validateHybridTableIndexColumns(_ context.Context, d *schema.ResourceDiff, _ any) error {
	columns := d.Get("column").([]interface{})
	indexes := d.Get("index").([]interface{})

	// Build a set of valid column names
	validColumns := make(map[string]bool)
	for _, c := range columns {
		colMap, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := colMap["name"].(string); ok {
			validColumns[strings.ToUpper(name)] = true
		}
	}

	// Check each index's columns
	for i, idx := range indexes {
		indexMap, ok := idx.(map[string]interface{})
		if !ok {
			continue
		}

		indexName := ""
		if name, ok := indexMap["name"].(string); ok {
			indexName = name
		}

		indexColumns, ok := indexMap["columns"].([]interface{})
		if !ok {
			continue
		}

		for _, col := range indexColumns {
			colName, ok := col.(string)
			if !ok {
				continue
			}
			if !validColumns[strings.ToUpper(colName)] {
				return fmt.Errorf("index[%d] (%s) references non-existent column %q; valid columns are: %v",
					i, indexName, colName, getColumnNames(columns))
			}
		}
	}

	return nil
}

// getColumnNames extracts column names from the column list for error messages.
func getColumnNames(columns []interface{}) []string {
	names := make([]string, 0, len(columns))
	for _, c := range columns {
		colMap, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := colMap["name"].(string); ok {
			names = append(names, name)
		}
	}
	return names
}

// suppressColumnNameDiff suppresses differences in column name casing
// Snowflake uppercases unquoted identifiers but preserves casing for quoted identifiers
func suppressColumnNameDiff(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	if oldValue == "" || newValue == "" {
		return false
	}

	// Use the normalization logic: simple identifiers compare as uppercase, complex ones as-is
	normalizedOld := NormalizeIdentifier(oldValue)
	normalizedNew := NormalizeIdentifier(newValue)

	return normalizedOld == normalizedNew
}

// forceNewOnColumnStructureChange forces recreation when column structure changes
// but allows comment-only changes to be updated in-place
func forceNewOnColumnStructureChange(_ context.Context, d *schema.ResourceDiff, _ any) error {
	if !d.HasChange("column") {
		return nil
	}

	oldCols, newCols := d.GetChange("column")
	oldColsList, ok1 := oldCols.([]interface{})
	newColsList, ok2 := newCols.([]interface{})

	if !ok1 || !ok2 {
		return nil
	}

	// If column count changed, force recreation
	if len(oldColsList) != len(newColsList) {
		return d.ForceNew("column")
	}

	// Check if any column names or types changed
	for i := range newColsList {
		if i >= len(oldColsList) {
			break
		}

		oldColMap, ok1 := oldColsList[i].(map[string]interface{})
		newColMap, ok2 := newColsList[i].(map[string]interface{})

		if !ok1 || !ok2 {
			continue
		}

		// Check if name changed (force recreation)
		oldName, _ := oldColMap["name"].(string)
		newName, _ := newColMap["name"].(string)
		if oldName != newName {
			return d.ForceNew("column")
		}

		// Check if type changed (force recreation)
		oldType, _ := oldColMap["type"].(string)
		newType, _ := newColMap["type"].(string)
		if oldType != newType {
			return d.ForceNew("column")
		}

		// Check if nullable changed (force recreation)
		oldNullable, oldHasNullable := oldColMap["nullable"].(bool)
		newNullable, newHasNullable := newColMap["nullable"].(bool)
		if oldHasNullable && newHasNullable && oldNullable != newNullable {
			return d.ForceNew("column")
		}

		// Comment changes are allowed - don't force recreation
	}

	return nil
}
