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
		Description: "Definitions of columns to create in the hybrid table. Minimum one required.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Column name",
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

			// Quote column names to match the quoted names in column definitions
			indexColumns := quoteColumnNames(rawIndexColumns)

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
		columnName := fmt.Sprintf(`"%v"`, colName)
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
		// Quote column names to match the quoted names in column definitions
		columns := quoteColumnNames(rawColumns)

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
				// Quote column names to match the quoted names in column definitions
				fkColumns := quoteColumnNames(rawFkColumns)

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

		column := map[string]interface{}{
			"name":     col.Name,
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

	// Build index list with only indexes that exist in Snowflake, maintaining config order
	indexList := make([]map[string]interface{}, 0, len(configIndexes))
	for name, order := range configIndexNames {
		idx, ok := snowflakeIndexes[name]
		if !ok {
			// Index not found in Snowflake - skip it to allow drift detection
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
		indexList = append(indexList, map[string]interface{}{
			"name":    configIdx["name"].(string),
			"columns": columns,
		})
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
				dropRequest := sdk.NewDropIndexRequest(indexId)
				err := client.Tables.DropIndex(ctx, dropRequest)
				if err != nil {
					return diag.FromErr(fmt.Errorf("error dropping index %v: %w", indexName, err))
				}
			}
		}

		// Create new indexes
		for indexName, columns := range newIndexMap {
			if _, exists := oldIndexMap[indexName]; !exists {
				indexId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), indexName)
				// Quote column names to match the quoted names in column definitions
				quotedColumns := quoteColumnNames(columns)
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
