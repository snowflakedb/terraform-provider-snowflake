package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
								Default:      "SIMPLE",
								Description:  "The match type for the foreign key: FULL, SIMPLE, or PARTIAL",
								ValidateFunc: validation.StringInSlice([]string{"FULL", "SIMPLE", "PARTIAL"}, true),
							},
							"on_update": {
								Type:         schema.TypeString,
								Optional:     true,
								Default:      "NO ACTION",
								Description:  "Action to perform when the primary/unique key is updated: CASCADE, SET NULL, SET DEFAULT, RESTRICT, or NO ACTION",
								ValidateFunc: validation.StringInSlice([]string{"CASCADE", "SET NULL", "SET DEFAULT", "RESTRICT", "NO ACTION"}, true),
							},
							"on_delete": {
								Type:         schema.TypeString,
								Optional:     true,
								Default:      "NO ACTION",
								Description:  "Action to perform when the primary/unique key is deleted: CASCADE, SET NULL, SET DEFAULT, RESTRICT, or NO ACTION",
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
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the index",
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

		CustomizeDiff: customdiff.All(
			ComputedIfAnyAttributeChanged(hybridTableSchema, ShowOutputAttributeName, "comment"),
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

	// Create indexes
	if v, ok := d.GetOk("index"); ok {
		indexes := v.([]interface{})
		for _, idx := range indexes {
			indexMap := idx.(map[string]interface{})
			indexName := indexMap["name"].(string)
			indexColumns := expandStringList(indexMap["columns"].([]interface{}))

			indexId := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, indexName)
			indexRequest := sdk.NewCreateIndexRequest(indexId, id, indexColumns)

			err = client.Tables.CreateIndex(ctx, indexRequest)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error creating index %v on hybrid table %v: %w", indexName, name, err))
			}
		}
	}

	d.SetId(helpers.EncodeSnowflakeID(id))

	return ReadHybridTable(ctx, d, meta)
}

func getHybridTableColumnRequests(from interface{}) ([]sdk.TableColumnRequest, error) {
	cols := from.([]interface{})
	requests := make([]sdk.TableColumnRequest, len(cols))

	for i, c := range cols {
		colMap := c.(map[string]interface{})
		columnName := fmt.Sprintf(`"%v"`, colMap["name"].(string))
		columnType := colMap["type"].(string)

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
	constraints := from.([]interface{})
	requests := make([]sdk.OutOfLineConstraintRequest, len(constraints))

	for i, c := range constraints {
		constraintMap := c.(map[string]interface{})
		constraintName := constraintMap["name"].(string)
		constraintType := constraintMap["type"].(string)
		columns := expandStringList(constraintMap["columns"].([]interface{}))

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
				fkMap := fk[0].(map[string]interface{})
				tableId := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(fkMap["table_id"].(string))
				fkColumns := expandStringList(fkMap["columns"].([]interface{}))

				fkRequest := sdk.NewOutOfLineForeignKeyRequest(tableId, fkColumns)

				if match, ok := fkMap["match"].(string); ok && match != "" {
					matchType, err := sdk.ToMatchType(match)
					if err != nil {
						return nil, fmt.Errorf("invalid match type: %w", err)
					}
					fkRequest.WithMatch(&matchType)
				}

				if onUpdate, ok := fkMap["on_update"].(string); ok && onUpdate != "" {
					onUpdateAction, err := sdk.ToForeignKeyAction(onUpdate)
					if err != nil {
						return nil, fmt.Errorf("invalid on_update action: %w", err)
					}

					onDelete := fkMap["on_delete"].(string)
					onDeleteAction, err := sdk.ToForeignKeyAction(onDelete)
					if err != nil {
						return nil, fmt.Errorf("invalid on_delete action: %w", err)
					}

					fkRequest.WithOn(sdk.NewForeignKeyOnAction().
						WithOnUpdate(&onUpdateAction).
						WithOnDelete(&onDeleteAction))
				}

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

	// Verify it's a hybrid table
	if table.Kind != "HYBRID TABLE" {
		return diag.FromErr(fmt.Errorf("expected HYBRID TABLE but got %s", table.Kind))
	}

	// Set basic attributes
	if err := d.Set("name", table.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("database", table.DatabaseName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", table.SchemaName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("comment", table.Comment); err != nil {
		return diag.FromErr(err)
	}

	// Get column details
	columnDetails, err := client.Tables.DescribeColumns(ctx, sdk.NewDescribeTableColumnsRequest(id))
	if err != nil {
		return diag.FromErr(err)
	}

	// Set columns
	columns := make([]map[string]interface{}, 0)
	for _, col := range columnDetails {
		if col.Kind != "COLUMN" {
			continue
		}

		column := map[string]interface{}{
			"name":     col.Name,
			"type":     string(col.Type),
			"nullable": col.IsNullable,
		}

		if col.Comment != nil {
			column["comment"] = *col.Comment
		}

		columns = append(columns, column)
	}
	if err := d.Set("column", columns); err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}

	// Set indexes
	indexList := make([]map[string]interface{}, 0, len(indexes))
	for _, idx := range indexes {
		index := map[string]interface{}{
			"name":    idx.Name,
			"columns": idx.Columns,
		}
		indexList = append(indexList, index)
	}
	if err := d.Set("index", indexList); err != nil {
		return diag.FromErr(err)
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
			alterRequest.WithUnset(sdk.NewTableUnsetRequest().WithComment(sdk.Bool(true)))
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
				createRequest := sdk.NewCreateIndexRequest(indexId, id, columns)
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
