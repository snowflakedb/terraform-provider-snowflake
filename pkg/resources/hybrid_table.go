package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
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
		Description:      blocklistedCharactersFieldDescription("The database in which to create the hybrid table."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
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
								ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
								DiffSuppressFunc: suppressIdentifierQuoting,
								Description:      "The default sequence for the column (uses NEXTVAL).",
							},
						},
					},
				},
				"collate": {
					Type:             schema.TypeString,
					Optional:         true,
					Default:          "",
					DiffSuppressFunc: ignoreCaseSuppressFunc,
					Description:      "Column collation specification, e.g. en-ci. Case-insensitive (en-ci and EN-CI are treated as equal).",
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
		Description: "Defines the primary key constraint for the hybrid table. Snowflake requires every hybrid table to have a primary key — this block is mandatory and cannot be omitted or removed. Snowflake does not support altering the primary key in place, so any change to `keys` (including reordering, adding, or removing columns) or to `name` forces recreation of the hybrid table.",
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
		Type:        schema.TypeSet,
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
		Type:        schema.TypeSet,
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
			StateContext: TrackingImportWrapper(resources.HybridTable, ImportName[sdk.SchemaObjectIdentifier]),
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

// hybridColumnSpec is the parsed form of a single hybrid-table column declaration,
// shared between the Create path (CREATE TABLE) and the alter-time ADD COLUMN
// path. The two paths hit the same parsing/validation routines but populate
// distinct SDK request types, so the spec carries only the fields that both
// requests can consume directly. Per-request-only fields (e.g. NotNull, which
// HybridTableAddColumnActionRequest does not support) live in their respective
// builders.
type hybridColumnSpec struct {
	dataType     datatypes.DataType
	defaultValue *sdk.ColumnDefaultValue // nil when no default block was provided
	collate      string                  // empty when not set
	comment      string                  // empty when not set
}

func buildHybridColumnSpec(col hybridTableColumn) (hybridColumnSpec, error) {
	dataType, err := datatypes.ParseDataType(col.dataType)
	if err != nil {
		return hybridColumnSpec{}, fmt.Errorf("invalid data type for column %s: %w", col.name, err)
	}

	spec := hybridColumnSpec{
		dataType: dataType,
		collate:  col.collate,
		comment:  col.comment,
	}

	if col._default != nil {
		defaultValue, err := buildHybridColumnDefaultFromParsed(col._default, dataType)
		if err != nil {
			return hybridColumnSpec{}, fmt.Errorf("column %q: %w", col.name, err)
		}
		spec.defaultValue = defaultValue
	}

	return spec, nil
}

func buildHybridTableColumnRequests(cols []any) ([]sdk.HybridTableColumnRequest, error) {
	parsed := parseHybridColumns(cols)
	requests := make([]sdk.HybridTableColumnRequest, len(parsed))
	for i, col := range parsed {
		spec, err := buildHybridColumnSpec(col)
		if err != nil {
			return nil, err
		}

		req := sdk.NewHybridTableColumnRequest(col.name, sdk.LegacyDataTypeWithAttrs(spec.dataType))
		if !col.nullable {
			req.WithNotNull(true)
		}
		if spec.defaultValue != nil {
			req.WithDefaultValue(*spec.defaultValue)
		}
		if spec.collate != "" {
			req.WithCollate(spec.collate)
		}
		if spec.comment != "" {
			req.WithComment(spec.comment)
		}
		requests[i] = *req
	}
	return requests, nil
}

// buildHybridAddColumnAction builds the alter-time add-column action from a
// parsed column. Delegates the shared parsing/validation work to
// buildHybridColumnSpec and only maps the spec onto the
// HybridTableAddColumnActionRequest type-specific fields. Note that this
// request has no NotNull — non-nullable columns can only be added via an
// InlineConstraint, but the resource forces recreation on nullable changes
// (see forceNewIfColumnNullableChanged), so a nullable=false branch here
// would be unreachable.
func buildHybridAddColumnAction(col hybridTableColumn) (*sdk.HybridTableAddColumnActionRequest, error) {
	spec, err := buildHybridColumnSpec(col)
	if err != nil {
		return nil, err
	}
	req := sdk.NewHybridTableAddColumnActionRequest(col.name, sdk.LegacyDataTypeWithAttrs(spec.dataType))
	if spec.defaultValue != nil {
		req.WithDefaultValue(*spec.defaultValue)
	}
	if spec.collate != "" {
		req.WithCollate(spec.collate)
	}
	if spec.comment != "" {
		req.WithComment(spec.comment)
	}
	return req, nil
}

// buildHybridAlterColumnTypeAction constructs a SET DATA TYPE alter-column
// action from a parsed column. Mirrors the Create and ADD COLUMN paths by
// running the raw type string through datatypes.ParseDataType (so that an
// invalid type fails client-side with a clear error) before converting via
// sdk.LegacyDataTypeWithAttrs, which preserves precision/scale/length on the
// way out (e.g. NUMBER(20,5) and VARCHAR(100) survive the round-trip).
//
// Using LegacyDataTypeFrom here would be a regression: it strips attributes,
// silently degrading NUMBER(20,5) to NUMBER (which Snowflake interprets as
// the default NUMBER(38,0)).
func buildHybridAlterColumnTypeAction(col hybridTableColumn) (*sdk.HybridTableAlterColumnActionRequest, error) {
	dataType, err := datatypes.ParseDataType(col.dataType)
	if err != nil {
		return nil, fmt.Errorf("invalid data type for column %s: %w", col.name, err)
	}
	return sdk.NewHybridTableAlterColumnActionRequest(col.name).
		WithType(sdk.LegacyDataTypeWithAttrs(dataType)), nil
}

// buildHybridColumnDefaultFromParsed converts the resource-level columnDefault
// into the map shape expected by buildHybridColumnDefaultValue, then delegates.
// It exists so the create-time and alter-time paths share the same default-value
// formatting and exclusivity validation.
func buildHybridColumnDefaultFromParsed(cd *columnDefault, dataType datatypes.DataType) (*sdk.ColumnDefaultValue, error) {
	defMap := map[string]any{}
	if cd.constant != nil {
		defMap["constant"] = *cd.constant
	}
	if cd.expression != nil {
		defMap["expression"] = *cd.expression
	}
	if cd.sequence != nil {
		defMap["sequence"] = *cd.sequence
	}
	return buildHybridColumnDefaultValue(defMap, dataType)
}

// buildHybridColumnDefaultValue builds a ColumnDefaultValue from a default block map
// and validates that exactly one of constant/expression/sequence is set. This
// validation lives here, rather than at the schema level, because Terraform
// plugin SDK v2 cannot express ExactlyOneOf across non-zero indices inside
// multi-element TypeLists (the parent column list).
func buildHybridColumnDefaultValue(defMap map[string]any, dataType datatypes.DataType) (*sdk.ColumnDefaultValue, error) {
	constant, hasConstant := defMap["constant"].(string)
	hasConstant = hasConstant && len(constant) > 0

	expr, hasExpression := defMap["expression"].(string)
	hasExpression = hasExpression && len(expr) > 0

	seq, hasSequence := defMap["sequence"].(string)
	hasSequence = hasSequence && len(seq) > 0

	set := 0
	for _, has := range []bool{hasConstant, hasExpression, hasSequence} {
		if has {
			set++
		}
	}
	if set != 1 {
		return nil, fmt.Errorf("default block must have exactly one of %q, %q, or %q set", "constant", "expression", "sequence")
	}

	var expression string
	switch {
	case hasConstant:
		if datatypes.IsTextDataType(dataType) {
			expression = "'" + strings.ReplaceAll(constant, "'", "''") + "'"
		} else {
			expression = constant
		}
	case hasExpression:
		expression = expr
	case hasSequence:
		expression = fmt.Sprintf("%v.NEXTVAL", seq)
	}

	return &sdk.ColumnDefaultValue{Expression: &expression}, nil
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
		for _, ucRaw := range v.(*schema.Set).List() {
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
		for _, fkRaw := range v.(*schema.Set).List() {
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
	if !reflect.DeepEqual(*set, *sdk.NewHybridTableSetPropertiesRequest()) {
		if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithSet(*set)); err != nil {
			d.Partial(true)
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

		columnState, err := buildHybridColumnStateFromDescribe(details, d)
		if err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.HybridTableToSchema(hybridTable)}),
			d.Set(DescribeOutputAttributeName, schemas.HybridTableDetailsListToSchema(details)),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("comment", hybridTable.Comment),
			d.Set("column", columnState),
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

	// Handle rename (name, database, or schema change). RENAME TO accepts a
	// fully-qualified identifier, so a database or schema change is realized
	// as a server-side move via the same statement.
	if d.HasChange("name") || d.HasChange("database") || d.HasChange("schema") {
		newId := sdk.NewSchemaObjectIdentifier(
			d.Get("database").(string),
			d.Get("schema").(string),
			d.Get("name").(string),
		)

		if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithNewName(newId)); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error renaming hybrid table from %v to %v: %w", id.FullyQualifiedName(), newId.FullyQualifiedName(), err))
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

	if !reflect.DeepEqual(*set, *sdk.NewHybridTableSetPropertiesRequest()) {
		if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithSet(*set)); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error setting hybrid table properties %v: %w", id.FullyQualifiedName(), err))
		}
	}
	if !reflect.DeepEqual(*unset, *sdk.NewHybridTableUnsetPropertiesRequest()) {
		if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithUnset(*unset)); err != nil {
			d.Partial(true)
			return diag.FromErr(fmt.Errorf("error unsetting hybrid table properties %v: %w", id.FullyQualifiedName(), err))
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
				d.Partial(true)
				return diag.FromErr(fmt.Errorf("error dropping columns from hybrid table %v: %w", id.FullyQualifiedName(), err))
			}
		}

		// Add new columns (one at a time)
		for _, col := range added {
			addReq, err := buildHybridAddColumnAction(col)
			if err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
			if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).WithAddColumnAction(*addReq)); err != nil {
				d.Partial(true)
				return diag.FromErr(fmt.Errorf("error adding column %s to hybrid table %v: %w", col.name, id.FullyQualifiedName(), err))
			}
		}

		// Alter changed columns
		for _, ch := range changed {
			if ch.changedDataType {
				alterReq, err := buildHybridAlterColumnTypeAction(ch.newColumn)
				if err != nil {
					d.Partial(true)
					return diag.FromErr(err)
				}
				if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
					WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{*alterReq})); err != nil {
					d.Partial(true)
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
					d.Partial(true)
					return diag.FromErr(fmt.Errorf("error altering column comment %s on hybrid table %v: %w",
						ch.newColumn.name, id.FullyQualifiedName(), err))
				}
			}

			if ch.droppedDefault {
				alterReq := sdk.NewHybridTableAlterColumnActionRequest(ch.newColumn.name).
					WithDropDefault(true)
				if err := client.HybridTables.Alter(ctx, sdk.NewAlterHybridTableRequest(id).
					WithAlterColumnAction([]sdk.HybridTableAlterColumnActionRequest{*alterReq})); err != nil {
					d.Partial(true)
					return diag.FromErr(fmt.Errorf("error dropping column default %s on hybrid table %v: %w",
						ch.newColumn.name, id.FullyQualifiedName(), err))
				}
			}
		}
	}

	return GetReadHybridTableFunc(false)(ctx, d, meta)
}

// ---------------------------------------------------------------------------
// Read helpers
// ---------------------------------------------------------------------------

func buildHybridColumnStateFromDescribe(details []sdk.HybridTableDetails, d *schema.ResourceData) ([]any, error) {
	flattened := make([]any, 0)

	// Build lookups from config for fields where DESCRIBE returns a different value
	// than the user wrote and the framework cannot reconcile it via DiffSuppressFunc:
	// - collate: DESCRIBE returns "X COLLATE 'Y'" combined; the SDK splits it but the
	//   server-side spelling can differ from what the user wrote (e.g. case).
	// - nullable: DESCRIBE returns NOT NULL for PK columns even when config says nullable=true.
	//   Preserve the config value since Snowflake silently enforces NOT NULL on PK columns.
	//
	// Data type drift (e.g. INTEGER vs NUMBER(38,0)) is handled by DiffSuppressDataTypes
	// on the column.type field — we do not need to substitute the config value here.
	type configColumnInfo struct {
		collate  string
		nullable bool
	}
	configByName := make(map[string]configColumnInfo)
	if configCols, ok := d.GetOk("column"); ok {
		for _, rawCol := range configCols.([]any) {
			colMap := rawCol.(map[string]any)
			if colName, ok := colMap["name"].(string); ok {
				info := configColumnInfo{nullable: true}
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
		// Prefer the user's config value for collate (it preserves the exact spelling
		// the user wrote, e.g. "en-ci"). Fall back to the SDK-split Collation from
		// DESCRIBE for imported tables where no config exists.
		collate := colInfo.collate
		if collate == "" && td.Collation != nil {
			collate = *td.Collation
		}
		flat := map[string]any{
			"name":     td.Name,
			"type":     td.Type,
			"nullable": colInfo.nullable,
			"comment":  td.Comment,
			"collate":  collate,
		}

		def, err := toHybridColumnDefaultConfig(td)
		if err != nil {
			return nil, err
		}
		if def != nil {
			flat["default"] = def
		}

		flattened = append(flattened, flat)
	}
	return flattened, nil
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
//
// Returns an error if the DESCRIBE output advertises a sequence default whose
// identifier portion cannot be parsed. The schema-level ValidateDiagFunc on the
// `sequence` field guarantees user-supplied values are valid identifiers, so a
// parse failure on round-trip indicates either a Snowflake bug or a previously
// drifted state — both warrant surfacing the error rather than silently falling
// back to a malformed value.
func toHybridColumnDefaultConfig(td sdk.HybridTableDetails) ([]any, error) {
	if td.Default == "" {
		return nil, nil
	}

	defaultRaw := td.Default

	// Sequence detection: ends with .NEXTVAL (case-insensitive)
	const nextvalSuffix = ".NEXTVAL"
	if strings.HasSuffix(strings.ToUpper(defaultRaw), nextvalSuffix) {
		sequenceIdRaw := defaultRaw[:len(defaultRaw)-len(nextvalSuffix)]
		id, err := sdk.ParseSchemaObjectIdentifier(sequenceIdRaw)
		if err != nil {
			return nil, fmt.Errorf("hybrid_table column %q: failed to parse sequence identifier %q from DESCRIBE: %w", td.Name, sequenceIdRaw, err)
		}
		return []any{map[string]any{"sequence": id.FullyQualifiedName()}}, nil
	}

	// Expression detection: contains parentheses
	if strings.Contains(defaultRaw, "(") && strings.Contains(defaultRaw, ")") {
		return []any{map[string]any{
			"expression": defaultRaw,
		}}, nil
	}

	// Constant: unescape for string types
	if sdk.IsStringType(td.Type) {
		unquoted := strings.TrimSuffix(strings.TrimPrefix(defaultRaw, "'"), "'")
		return []any{map[string]any{
			"constant": strings.ReplaceAll(unquoted, "''", "'"),
		}}, nil
	}

	return []any{map[string]any{
		"constant": defaultRaw,
	}}, nil
}

// ---------------------------------------------------------------------------
// CustomizeDiff
// ---------------------------------------------------------------------------

// forceNewIfColumnCollateChanged forces recreation when collation changes on an
// existing column. The SDK's HybridTableAlterColumnActionRequest has no WithCollate
// method, so collation can only be set at column creation time.
func forceNewIfColumnCollateChanged() schema.CustomizeDiffFunc {
	return forceNewIfColumnFieldChanged("collate", func(o, n hybridTableColumn) bool {
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
	return forceNewIfColumnFieldChanged("nullable", func(o, n hybridTableColumn) bool {
		return o.nullable != n.nullable
	})
}

// forceNewIfColumnFieldChanged returns a CustomizeDiffFunc that forces recreation
// when a nested column field changes. It must call diff.ForceNew on the specific
// nested path (e.g. "column.0.nullable"), not on the parent "column" list — the
// terraform-plugin-sdk/v2 ForceNew sets RequiresNew on the resolved leaf schema,
// and a TypeList parent does not propagate that flag down to its diff entries.
func forceNewIfColumnFieldChanged(fieldName string, changed func(old, new hybridTableColumn) bool) schema.CustomizeDiffFunc {
	return func(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
		if !diff.HasChange("column") {
			return nil
		}
		oldRaw, newRaw := diff.GetChange("column")
		oldCols := parseHybridColumns(oldRaw)
		newCols := parseHybridColumns(newRaw)
		for _, o := range oldCols {
			for newIdx, n := range newCols {
				if o.name == n.name && changed(o, n) {
					return diff.ForceNew(fmt.Sprintf("column.%d.%s", newIdx, fieldName))
				}
			}
		}
		return nil
	}
}
