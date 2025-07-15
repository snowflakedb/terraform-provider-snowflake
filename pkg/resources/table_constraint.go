package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// TODO [SNOW-1348114]: refine this resource during redesign:
// - read (from the existing comment it seems that active warehouse is needed (it should be probably added to the resource as required)
// - drop (in tests it's not dropped correctly, probably also because missing warehouse)
// - do we need it?
// - not null cannot be set as a named constraint but it can be set by alter column statement - should it be added back here?
var tableConstraintSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of constraint",
	},
	"type": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Type of constraint, one of 'UNIQUE', 'PRIMARY KEY', or 'FOREIGN KEY'",
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice(sdk.AsStringList(sdk.AllColumnConstraintTypes), false),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.EqualFold(old, new)
		},
		StateFunc: func(val any) string {
			return strings.ToUpper(val.(string))
		},
	},
	"table_id": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      `Identifier for table to create constraint on. Format must follow: "\"&lt;db_name&gt;\".\"&lt;schema_name&gt;\".\"&lt;table_name&gt;\"" or "&lt;db_name&gt;.&lt;schema_name&gt;.&lt;table_name&gt;" (snowflake_table.my_table.id)`,
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
	},
	"columns": {
		Type:     schema.TypeList,
		MinItems: 1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		ForceNew:    true,
		Required:    true,
		Description: "Columns to use in constraint key",
	},
	"enforced": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: "Whether the constraint is enforced",
	},
	"deferrable": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     true,
		Description: "Whether the constraint is deferrable",
	},
	"initially": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		Default:      "DEFERRED",
		Description:  "Whether the constraint is initially deferred or immediate",
		ValidateFunc: validation.StringInSlice([]string{"DEFERRED", "IMMEDIATE"}, true),
	},
	"enable": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     true,
		Description: "Specifies whether the constraint is enabled or disabled. These properties are provided for compatibility with Oracle.",
	},
	"validate": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: "Specifies whether to validate existing data on the table when a constraint is created. Only used in conjunction with the ENABLE property.",
	},
	"rely": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     true,
		Description: "Specifies whether a constraint in NOVALIDATE mode is taken into account during query rewrite.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Comment for the table constraint",
		Deprecated:  "Not used. Will be removed.",
	},
	"foreign_key_properties": {
		Type:        schema.TypeList,
		Optional:    true,
		ForceNew:    true,
		MaxItems:    1,
		Description: "Additional properties when type is set to foreign key. Not applicable for primary/unique keys",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"references": {
					Type:        schema.TypeList,
					Required:    true,
					ForceNew:    true,
					MaxItems:    1,
					Description: "The table and columns that the foreign key references.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"table_id": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "Name of constraint",
							},
							"columns": {
								Type:     schema.TypeList,
								MinItems: 1,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
								Required:    true,
								Description: "Columns to use in foreign key reference",
							},
						},
					},
				},
				"match": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      "FULL",
					ValidateFunc: validation.StringInSlice(sdk.AsStringList(sdk.AllMatchTypes), true),
					Description:  "The match type for the foreign key. Not applicable for primary/unique keys",
				},
				"on_update": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      "NO ACTION",
					ValidateFunc: validation.StringInSlice(sdk.AsStringList(sdk.AllForeignKeyActions), true),
					Description:  "Specifies the action performed when the primary/unique key for the foreign key is updated. Not applicable for primary/unique keys",
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return strings.EqualFold(old, new)
					},
					StateFunc: func(val any) string {
						return strings.ToUpper(val.(string))
					},
				},
				"on_delete": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      "NO ACTION",
					ValidateFunc: validation.StringInSlice(sdk.AsStringList(sdk.AllForeignKeyActions), true),
					Description:  "Specifies the action performed when the primary/unique key for the foreign key is deleted. Not applicable for primary/unique keys",
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return strings.EqualFold(old, new)
					},
					StateFunc: func(val any) string {
						return strings.ToUpper(val.(string))
					},
				},
			},
		},
	},
}

func TableConstraint() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.TableConstraintResource), TrackingCreateWrapper(resources.TableConstraint, CreateTableConstraint)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.TableConstraintResource), TrackingReadWrapper(resources.TableConstraint, ReadTableConstraint)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.TableConstraintResource), TrackingUpdateWrapper(resources.TableConstraint, UpdateTableConstraint)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.TableConstraintResource), TrackingDeleteWrapper(resources.TableConstraint, DeleteTableConstraint)),

		Schema: tableConstraintSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

type tableConstraintID struct {
	name           string
	constraintType string
	tableID        string
}

func (v *tableConstraintID) String() string {
	return fmt.Sprintf("%s❄️%s❄️%s", v.name, v.constraintType, v.tableID)
}

func (v *tableConstraintID) parse(s string) {
	parts := strings.Split(s, "❄️")
	v.name = parts[0]
	v.constraintType = parts[1]
	v.tableID = parts[2]
}

func getTableIdentifier(s string) (*sdk.SchemaObjectIdentifier, error) {
	var objectIdentifier sdk.ObjectIdentifier
	var err error
	// TODO [SNOW-1348114]: Address during table rework; Fallback for old implementations using table.id instead of table.fully_qualified_name - probably will be removed later.
	if strings.Contains(s, "|") {
		objectIdentifier = helpers.DecodeSnowflakeIDLegacy(s)
	} else {
		objectIdentifier, err = helpers.DecodeSnowflakeParameterID(s)
	}

	if err != nil {
		return nil, fmt.Errorf("table id is incorrect: %s, err: %w", objectIdentifier, err)
	}
	tableIdentifier, ok := objectIdentifier.(sdk.SchemaObjectIdentifier)
	if !ok {
		return nil, fmt.Errorf("table id is incorrect: %s", objectIdentifier)
	}
	return &tableIdentifier, nil
}

// CreateTableConstraint implements schema.CreateFunc.
func CreateTableConstraint(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	cType := d.Get("type").(string)
	tableID := d.Get("table_id").(string)

	tableIdentifier, err := getTableIdentifier(tableID)
	if err != nil {
		return diag.FromErr(err)
	}

	constraintType, err := sdk.ToColumnConstraintType(cType)
	if err != nil {
		return diag.FromErr(err)
	}
	constraintRequest := sdk.NewOutOfLineConstraintRequest(constraintType).WithName(&name)

	cc := d.Get("columns").([]interface{})
	columns := make([]string, 0, len(cc))
	for _, c := range cc {
		columns = append(columns, c.(string))
	}
	constraintRequest.WithColumns(snowflake.QuoteStringList(columns))

	if v, ok := d.GetOk("enforced"); ok {
		constraintRequest.WithEnforced(sdk.Bool(v.(bool)))
	}

	if v, ok := d.GetOk("deferrable"); ok {
		constraintRequest.WithDeferrable(sdk.Bool(v.(bool)))
	}

	if v, ok := d.GetOk("initially"); ok {
		if v.(string) == "DEFERRED" {
			constraintRequest.WithInitiallyDeferred(sdk.Bool(true))
		} else {
			constraintRequest.WithInitiallyImmediate(sdk.Bool(true))
		}
	}

	if v, ok := d.GetOk("enable"); ok {
		constraintRequest.WithEnable(sdk.Bool(v.(bool)))
	}

	if v, ok := d.GetOk("validate"); ok {
		constraintRequest.WithValidate(sdk.Bool(v.(bool)))
	}

	if v, ok := d.GetOk("rely"); ok {
		constraintRequest.WithRely(sdk.Bool(v.(bool)))
	}

	// set foreign key specific settings
	if v, ok := d.GetOk("foreign_key_properties"); ok {
		foreignKeyProperties := v.([]interface{})[0].(map[string]interface{})
		references := foreignKeyProperties["references"].([]interface{})[0].(map[string]interface{})
		fkTableID := references["table_id"].(string)
		fkId, err := helpers.DecodeSnowflakeParameterID(fkTableID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("table id is incorrect: %s, err: %w", fkTableID, err))
		}
		referencedTableIdentifier, ok := fkId.(sdk.SchemaObjectIdentifier)
		if !ok {
			return diag.FromErr(fmt.Errorf("table id is incorrect: %s", fkId))
		}

		cols := references["columns"].([]interface{})
		var fkColumns []string
		for _, c := range cols {
			fkColumns = append(fkColumns, c.(string))
		}
		foreignKeyRequest := sdk.NewOutOfLineForeignKeyRequest(referencedTableIdentifier, snowflake.QuoteStringList(fkColumns))

		matchType, err := sdk.ToMatchType(foreignKeyProperties["match"].(string))
		if err != nil {
			return diag.FromErr(err)
		}
		foreignKeyRequest.WithMatch(&matchType)

		onUpdate, err := sdk.ToForeignKeyAction(foreignKeyProperties["on_update"].(string))
		if err != nil {
			return diag.FromErr(err)
		}
		onDelete, err := sdk.ToForeignKeyAction(foreignKeyProperties["on_delete"].(string))
		if err != nil {
			return diag.FromErr(err)
		}
		foreignKeyRequest.WithOn(sdk.NewForeignKeyOnAction().
			WithOnDelete(&onDelete).
			WithOnUpdate(&onUpdate),
		)
		constraintRequest.WithForeignKey(foreignKeyRequest)
	}

	alterStatement := sdk.NewAlterTableRequest(*tableIdentifier).WithConstraintAction(sdk.NewTableConstraintActionRequest().WithAdd(constraintRequest))
	err = client.Tables.Alter(ctx, alterStatement)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating table constraint %v err = %w", name, err))
	}

	tc := tableConstraintID{
		name,
		cType,
		tableID,
	}
	d.SetId(tc.String())

	return ReadTableConstraint(ctx, d, meta)
}

// ReadTableConstraint implements schema.ReadFunc.
func ReadTableConstraint(ctx context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// TODO(issue-2683): Implement read operation
	// commenting this out since it requires an active warehouse to be set which may not be intuitive.
	// also it takes a while for the database to reflect changes. Would likely need to add a validation
	// step like in tag association. People don't like waiting 40 minutes for Terraform to run.

	/*providerContext := meta.(*provider.Context)
	tc := tableConstraintID{}
	tc.parse(d.Id())
	databaseName, schemaName, tableName := snowflakeValidation.ParseFullyQualifiedObjectID(tc.tableID)

	// just need to check to make sure it exists
	_, err := snowflake.ShowTableConstraint(tc.name, databaseName, schemaName, tableName, db)
	if err != nil {
		return diag.FromErr(fmt.Errorf(fmt.Sprintf("error reading table constraint %v", tc.String()))
	}*/

	return nil
}

// UpdateTableConstraint implements schema.UpdateFunc.
func UpdateTableConstraint(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	/* TODO(issue-2683): Update isn't be possible with non-existing Read operation. The Update logic is ready to be uncommented once the Read operation is ready.
	client := meta.(*provider.Context).Client


	tc := tableConstraintID{}
	tc.parse(d.Id())

	tableIdentifier, err := getTableIdentifier(tc.tableID)
	if err != nil {
		return diag.FromErr(err)
	}

		if d.HasChange("name") {
			newName := d.Get("name").(string)
			constraintRequest := sdk.NewTableConstraintRenameActionRequest().WithOldName(tc.name).WithNewName(newName)
			alterStatement := sdk.NewAlterTableRequest(*tableIdentifier).WithConstraintAction(sdk.NewTableConstraintActionRequest().WithRename(constraintRequest))

			err = client.Tables.Alter(ctx, alterStatement)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error renaming table constraint %s err = %w", tc.name, err)
			}

			tc.name = newName
			d.SetId(tc.String())
		}
	*/

	return ReadTableConstraint(ctx, d, meta)
}

// DeleteTableConstraint implements schema.DeleteFunc.
func DeleteTableConstraint(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	tc := tableConstraintID{}
	tc.parse(d.Id())

	tableIdentifier, err := getTableIdentifier(tc.tableID)
	if err != nil {
		return diag.FromErr(err)
	}

	dropRequest := sdk.NewTableConstraintDropActionRequest().WithConstraintName(&tc.name)
	alterStatement := sdk.NewAlterTableRequest(*tableIdentifier).WithConstraintAction(sdk.NewTableConstraintActionRequest().WithDrop(dropRequest))
	err = client.Tables.Alter(ctx, alterStatement)
	if err != nil {
		// if the table constraint does not exist, then remove from state file
		if strings.Contains(err.Error(), "does not exist") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error dropping table constraint %v err = %w", tc.name, err))
	}

	d.SetId("")

	return nil
}
