package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
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
					// Matches the data-type field pattern. StateFunc does not reliably fire for
					// nested TypeList fields, so buildHybridColumnStateFromDescribe is the actual
					// Read-path normalizer — keep it.
					StateFunc: DataTypeStateFunc,
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
					Computed:    true,
					ForceNew:    true,
					Description: "Constraint name. If omitted, Snowflake auto-generates one (visible in state after the first apply).",
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
		Set:         uniqueConstraintHash,
		Description: "Defines UNIQUE constraints. Can only be set at creation time. Any change forces recreation.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: "Constraint name. If omitted, Snowflake auto-generates one (visible in state after the first apply).",
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
		Set:         foreignKeyHash,
		Description: "Defines FOREIGN KEY constraints. Can only be set at creation time. Any change forces recreation.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					ForceNew:    true,
					Description: "Constraint name. If omitted, Snowflake auto-generates one (visible in state after the first apply).",
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
	"index": {
		Type:        schema.TypeSet,
		Optional:    true,
		ForceNew:    true,
		Set:         indexHash,
		Description: "Defines secondary indexes on the hybrid table. Can only be set at creation time (declared inline in CREATE HYBRID TABLE). Any change to an index forces recreation of the table.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					Description: "Name of the secondary index. Snowflake requires an explicit name for inline indexes.",
				},
				"columns": {
					Type:        schema.TypeList,
					Required:    true,
					ForceNew:    true,
					MinItems:    1,
					Elem:        &schema.Schema{Type: schema.TypeString, DiffSuppressFunc: ignoreCaseSuppressFunc},
					Description: "Index key columns, in order. Order is semantically meaningful.",
				},
				"include_columns": {
					Type:        schema.TypeSet,
					Optional:    true,
					ForceNew:    true,
					Set:         indexIncludeColumnsHash,
					Elem:        &schema.Schema{Type: schema.TypeString, DiffSuppressFunc: ignoreCaseSuppressFunc},
					Description: "Columns included in the index payload via INCLUDE (...). Order carries no meaning.",
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
		// TODO(SNOW-3637521): Add PreviewFeatureCreateContextWrapper when the hybrid table resource is moved to the production provider.
		CreateContext: TrackingCreateWrapper(resources.HybridTable, CreateHybridTable),
		ReadContext:   TrackingReadWrapper(resources.HybridTable, GetReadHybridTableFunc(true)),
		UpdateContext: TrackingUpdateWrapper(resources.HybridTable, UpdateHybridTable),
		DeleteContext: TrackingDeleteWrapper(resources.HybridTable, deleteFunc),
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
// Set hash functions for TypeSet constraint blocks
// ---------------------------------------------------------------------------

// uniqueConstraintHash hashes a unique_constraint set element on its columns only,
// excluding name so an auto-generated server name does not churn the set element.
func uniqueConstraintHash(v any) int {
	m := v.(map[string]any)
	var b strings.Builder
	for _, col := range m["columns"].([]any) {
		b.WriteString(col.(string))
		b.WriteByte(',')
	}
	return schema.HashString(b.String())
}

// foreignKeyHash hashes a foreign_key set element on its stable identity: local
// columns + referenced table + referenced columns. name is excluded (same reason as
// uniqueConstraintHash). table_id is normalized so a quoted read-back value and an
// unquoted config value hash identically; a parse failure falls back to the raw string.
func foreignKeyHash(v any) int {
	m := v.(map[string]any)
	var b strings.Builder
	for _, col := range m["columns"].([]any) {
		b.WriteString(col.(string))
		b.WriteByte(',')
	}
	if refList, ok := m["references"].([]any); ok && len(refList) > 0 {
		ref := refList[0].(map[string]any)
		tableId := ref["table_id"].(string)
		if parsed, err := sdk.ParseSchemaObjectIdentifier(tableId); err == nil {
			tableId = parsed.FullyQualifiedName()
		}
		b.WriteByte('|')
		b.WriteString(tableId)
		b.WriteByte('|')
		for _, col := range ref["columns"].([]any) {
			b.WriteString(col.(string))
			b.WriteByte(',')
		}
	}
	return schema.HashString(b.String())
}

// indexIncludeColumnsHash hashes a single include_column string (uppercased) so
// that a lowercase config value and its uppercase SHOW INDEXES read-back land in
// the same set bucket. Used as the Set function for the include_columns TypeSet.
func indexIncludeColumnsHash(v any) int {
	return schema.HashString(strings.ToUpper(v.(string)))
}

// indexHash hashes an index set element on its stable identity: the user-supplied
// name plus the index columns and include columns, with column names uppercased.
// Uppercasing is required because SHOW INDEXES returns uppercase column names while
// a config may use any case; the TypeSet element identity is resolved via this hash
// before DiffSuppressFunc runs, so without normalization a lowercase config element
// and its uppercase read-back would hash differently and churn the set (spurious
// ForceNew). name is included verbatim — it is user-supplied and round-trips exactly
// (no server auto-naming, unlike the constraint blocks).
func indexHash(v any) int {
	m := v.(map[string]any)
	var b strings.Builder
	b.WriteString(m["name"].(string))
	b.WriteByte('|')
	for _, col := range m["columns"].([]any) {
		b.WriteString(strings.ToUpper(col.(string)))
		b.WriteByte(',')
	}
	b.WriteByte('|')
	if inc, ok := m["include_columns"]; ok && inc != nil {
		incSet, ok := inc.(*schema.Set)
		if ok {
			incCols := make([]string, 0, incSet.Len())
			for _, c := range incSet.List() {
				incCols = append(incCols, strings.ToUpper(c.(string)))
			}
			sort.Strings(incCols)
			for _, c := range incCols {
				b.WriteString(c)
				b.WriteByte(',')
			}
		}
	}
	return schema.HashString(b.String())
}

// ---------------------------------------------------------------------------
// Column types and parsing helpers
// ---------------------------------------------------------------------------

func parseHybridColumn(from any) column {
	c := from.(map[string]any)
	var cd *columnDefault

	if defaultList, ok := c["default"].([]any); ok && len(defaultList) == 1 {
		cd = getHybridColumnDefault(defaultList[0].(map[string]any))
	}

	return column{
		name:     c["name"].(string),
		dataType: c["type"].(string),
		nullable: c["nullable"].(bool),
		_default: cd,
		collate:  c["collate"].(string),
		comment:  c["comment"].(string),
	}
}

// getHybridColumnDefault fills every non-empty default sub-field present in
// the schema, unlike table.go's getColumnDefault which short-circuits on the
// first match. The build helper (buildHybridColumnDefaultValue) requires
// visibility into all set fields to enforce mutual exclusivity — short-
// circuiting at parse time would silently drop the conflict and let apply
// proceed with whichever field was checked first.
func getHybridColumnDefault(def map[string]any) *columnDefault {
	cd := &columnDefault{}
	hasAny := false

	if v, ok := def["constant"].(string); ok && len(v) > 0 {
		cd.constant = &v
		hasAny = true
	}
	if v, ok := def["expression"].(string); ok && len(v) > 0 {
		cd.expression = &v
		hasAny = true
	}
	if v, ok := def["sequence"].(string); ok && len(v) > 0 {
		cd.sequence = &v
		hasAny = true
	}

	if !hasAny {
		return nil
	}
	return cd
}

func parseHybridColumns(from any) columns {
	cols := from.([]any)
	result := make(columns, len(cols))
	for i, c := range cols {
		result[i] = parseHybridColumn(c)
	}
	return result
}

func (c columns) getChangedHybridColumnProperties(new columns) []changedColumn {
	changed := make([]changedColumn, 0)
	for _, oldColumn := range c {
		for _, newColumn := range new {
			if oldColumn.name != newColumn.name {
				continue
			}
			ch := changedColumn{newColumn: newColumn}
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

func hybridTableColumnDiffs(oldCols, newCols columns) (removed columns, added columns, changed []changedColumn) {
	return oldCols.getNewIn(newCols), newCols.getNewIn(oldCols), oldCols.getChangedHybridColumnProperties(newCols)
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

func buildHybridColumnSpec(col column) (hybridColumnSpec, error) {
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
func buildHybridAddColumnAction(col column) (*sdk.HybridTableAddColumnActionRequest, error) {
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
func buildHybridAlterColumnTypeAction(col column) (*sdk.HybridTableAlterColumnActionRequest, error) {
	dataType, err := datatypes.ParseDataType(col.dataType)
	if err != nil {
		return nil, fmt.Errorf("invalid data type for column %s: %w", col.name, err)
	}
	return sdk.NewHybridTableAlterColumnActionRequest(col.name).
		WithDataType(sdk.LegacyDataTypeWithAttrs(dataType)), nil
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
// validation lives here, rather than at the schema level. Schema-level
// ExactlyOneOf would require paths like "column.0.default.0.constant", but
// terraform-plugin-sdk/v2 (helper/schema/schema.go: checkKeysAgainstSchemaFlags)
// rejects multi-element TypeList parents in path references at provider boot:
// "configuration block reference (...) can only be used with TypeList and
// MaxItems: 1 configuration blocks". The `column` list has no MaxItems.
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

	switch {
	case hasConstant:
		if datatypes.IsTextDataType(dataType) {
			val := "'" + strings.ReplaceAll(constant, "'", "''") + "'"
			return &sdk.ColumnDefaultValue{Expression: &val}, nil
		}
		c := constant
		return &sdk.ColumnDefaultValue{Expression: &c}, nil
	case hasExpression:
		return &sdk.ColumnDefaultValue{Expression: &expr}, nil
	case hasSequence:
		val := fmt.Sprintf("%v.NEXTVAL", seq)
		return &sdk.ColumnDefaultValue{Expression: &val}, nil
	}
	return nil, fmt.Errorf("unreachable: default block passed validation but no type matched")
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

// buildOutOfLineIndexes builds the inline secondary-index requests from the
// `index` TypeSet. Mirrors the optional-constraint loops in buildOutOfLineConstraints.
// The block is create-only and ForceNew, so there is no corresponding update path.
func buildOutOfLineIndexes(d *schema.ResourceData) []sdk.HybridTableOutOfLineIndexRequest {
	indexes := make([]sdk.HybridTableOutOfLineIndexRequest, 0)
	if v, ok := d.GetOk("index"); ok {
		for _, idxRaw := range v.(*schema.Set).List() {
			idxMap := idxRaw.(map[string]any)
			req := sdk.NewHybridTableOutOfLineIndexRequest(
				idxMap["name"].(string),
				expandStringList(idxMap["columns"].([]any)),
			)
			if incRaw, ok := idxMap["include_columns"].(*schema.Set); ok && incRaw.Len() > 0 {
				req.WithIncludeColumns(expandStringList(incRaw.List()))
			}
			indexes = append(indexes, *req)
		}
	}
	return indexes
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
		OutOfLineIndex:      buildOutOfLineIndexes(d),
	}
	request := sdk.NewCreateHybridTableRequest(id, columnsAndConstraints)

	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	if diags := handleHybridTableParametersCreate(d, request); diags.HasError() {
		return diags
	}

	if err := client.HybridTables.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating hybrid table %v: %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

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

		constraints, err := client.HybridTables.GetConstraints(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("reading hybrid table constraints: %w", err))
		}

		// Index read-back is best-effort — failure must not fail Read of an otherwise-healthy table.
		indexes, indexErr := client.HybridTables.ShowIndexes(ctx,
			sdk.NewShowIndexesHybridTableRequest().WithIn(sdk.TableIn{Table: id}))

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.HybridTableToSchema(hybridTable)}),
			d.Set(DescribeOutputAttributeName, schemas.HybridTableDetailsListToSchema(details)),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("comment", hybridTable.Comment),
			d.Set("column", columnState),
			d.Set("primary_key", buildPrimaryKeyStateFromConstraints(constraints)),
			d.Set("unique_constraint", buildUniqueConstraintsStateFromConstraints(constraints)),
			d.Set("foreign_key", buildForeignKeysStateFromConstraints(constraints)),
			d.Set("index", readIndexState(indexes, indexErr, constraints, id)),
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
		d.Partial(true)
		return diag.FromErr(err)
	}
	if diags := handleHybridTableParametersChanges(d, set, unset); diags.HasError() {
		d.Partial(true)
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
		removed, added, changed := hybridTableColumnDiffs(parseHybridColumns(oldRaw), parseHybridColumns(newRaw))

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

	// Reconciliation strategy: preserve the user's config spelling on Read
	// whenever the DESCRIBE value is equivalent (modulo the relevant
	// suppression). Field-level DiffSuppressFunc hides drift in the displayed
	// plan but does NOT prevent HasChange("column") from returning true at the
	// block level (TypeList HasChange does a DeepEqual on raw values). That
	// flips ComputedIfAnyAttributeChanged → describe_output Computed → plan
	// shows describe_output going to (known after apply) → "non-empty plan
	// after apply".
	//
	// - column.<idx>.type: substitute config spelling when canonically equal
	//   (uses datatypes.AreTheSame, the same comparator behind
	//   DiffSuppressDataTypes). Real type changes still surface as drift.
	// - column.<idx>.collate: substitute config spelling when case-equal
	//   (mirrors ignoreCaseSuppressFunc on the field).
	// - column.<idx>.nullable: PK columns silently come back as NOT NULL.
	//   Substitute the config value. Hybrid tables do not support ALTER
	//   SET/DROP NOT NULL, so external drift on this attribute cannot occur
	//   — the substitution is safe.
	type configColumnInfo struct {
		typeStr  string
		collate  string
		nullable bool
		found    bool
	}
	configByName := make(map[string]configColumnInfo)
	if configCols, ok := d.GetOk("column"); ok {
		for _, rawCol := range configCols.([]any) {
			colMap, ok := rawCol.(map[string]any)
			if !ok {
				continue
			}
			colName, ok := colMap["name"].(string)
			if !ok {
				continue
			}
			info := configColumnInfo{nullable: true, found: true}
			if t, ok := colMap["type"].(string); ok {
				info.typeStr = t
			}
			if c, ok := colMap["collate"].(string); ok {
				info.collate = c
			}
			if n, ok := colMap["nullable"].(bool); ok {
				info.nullable = n
			}
			configByName[strings.ToUpper(colName)] = info
		}
	}
	pkKeys := make(map[string]struct{})
	if pkRaw, ok := d.GetOk("primary_key"); ok {
		if pkList, ok := pkRaw.([]any); ok && len(pkList) > 0 {
			if pkMap, ok := pkList[0].(map[string]any); ok {
				if keysRaw, ok := pkMap["keys"].([]any); ok {
					for _, k := range keysRaw {
						if s, ok := k.(string); ok {
							pkKeys[strings.ToUpper(s)] = struct{}{}
						}
					}
				}
			}
		}
	}

	for _, td := range details {
		if td.Kind != "COLUMN" {
			continue
		}

		cfg := configByName[strings.ToUpper(td.Name)]

		typeOut := td.Type
		if cfg.found && cfg.typeStr != "" {
			cfgParsed, errCfg := datatypes.ParseDataType(cfg.typeStr)
			descParsed, errDesc := datatypes.ParseDataType(td.Type)
			if errCfg == nil && errDesc == nil && datatypes.AreTheSame(cfgParsed, descParsed) {
				typeOut = cfg.typeStr
			}
		}

		collate := ""
		if td.Collation != nil {
			collate = *td.Collation
		}
		if cfg.found && cfg.collate != "" && strings.EqualFold(cfg.collate, collate) {
			collate = cfg.collate
		}

		nullable := td.IsNullable
		if _, isPK := pkKeys[strings.ToUpper(td.Name)]; isPK {
			if cfg.found {
				nullable = cfg.nullable
			} else {
				// Externally added PK column not in config: schema default is true.
				nullable = true
			}
		}
		flat := map[string]any{
			"name":     td.Name,
			"type":     typeOut,
			"nullable": nullable,
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

// buildPrimaryKeyStateFromConstraints returns the primary_key block (at most one) from
// the merged constraint list, including the server-side name (possibly auto-generated).
func buildPrimaryKeyStateFromConstraints(constraints []sdk.HybridTableConstraint) []map[string]any {
	for _, c := range constraints {
		if c.Kind == sdk.ColumnConstraintTypePrimaryKey {
			return []map[string]any{{"name": c.Name, "keys": c.Columns}}
		}
	}
	return nil
}

// buildUniqueConstraintsStateFromConstraints returns the unique_constraint blocks from
// the merged constraint list, including the server-side name (possibly auto-generated).
func buildUniqueConstraintsStateFromConstraints(constraints []sdk.HybridTableConstraint) []map[string]any {
	var result []map[string]any
	for _, c := range constraints {
		if c.Kind == sdk.ColumnConstraintTypeUnique {
			result = append(result, map[string]any{"name": c.Name, "columns": c.Columns})
		}
	}
	return result
}

// buildForeignKeysStateFromConstraints returns the foreign_key blocks from the merged
// constraint list, including the server-side name (possibly auto-generated). The inner
// []map[string]any matches the references TypeList (MaxItems 1) schema block.
func buildForeignKeysStateFromConstraints(constraints []sdk.HybridTableConstraint) []map[string]any {
	var result []map[string]any
	for _, c := range constraints {
		if c.Kind == sdk.ColumnConstraintTypeForeignKey {
			// DeleteRule/UpdateRule are intentionally not mapped — the schema does not expose FK rules.
			result = append(result, map[string]any{
				"name":    c.Name,
				"columns": c.Columns,
				"references": []map[string]any{{
					"table_id": c.ReferencedTable.FullyQualifiedName(),
					"columns":  c.ReferencedColumns,
				}},
			})
		}
	}
	return result
}

// readIndexState builds index state from a SHOW INDEXES result, excluding FK-backing
// indexes. Returns nil on error (best-effort: SHOW INDEXES failure must not fail Read
// of an otherwise-healthy table).
func readIndexState(indexes []sdk.HybridTableIndex, indexErr error, constraints []sdk.HybridTableConstraint, id sdk.SchemaObjectIdentifier) []map[string]any {
	if indexErr != nil {
		log.Printf("[WARN] SHOW INDEXES failed for %s; skipping index read-back: %v", id.FullyQualifiedName(), indexErr)
		return nil
	}
	// FK constraints produce a system-managed backing index in SHOW INDEXES that must
	// be excluded: named FKs share the constraint name; anonymous FKs are named SYS_INDEX_..._FOREIGN_KEY_...
	var userIndexes []sdk.HybridTableIndex
outer:
	for _, idx := range indexes {
		upper := strings.ToUpper(idx.Name)
		if strings.HasPrefix(upper, "SYS_INDEX_") {
			continue
		}
		for _, c := range constraints {
			if c.Kind == sdk.ColumnConstraintTypeForeignKey && strings.EqualFold(c.Name, idx.Name) {
				continue outer
			}
		}
		userIndexes = append(userIndexes, idx)
	}
	return buildIndexesStateFromShowIndexes(userIndexes)
}

// buildIndexesStateFromShowIndexes maps SHOW INDEXES rows to index block state, keeping only user secondary indexes (IsUnique=false; nil excluded).
func buildIndexesStateFromShowIndexes(indexes []sdk.HybridTableIndex) []map[string]any {
	var result []map[string]any
	for _, idx := range indexes {
		if idx.IsUnique == nil || *idx.IsUnique {
			continue
		}
		var columns []string
		if idx.Columns != nil {
			columns = sdk.ParseCommaSeparatedStringArray(*idx.Columns, false)
		}
		result = append(result, map[string]any{
			"name":            idx.Name,
			"columns":         columns,
			"include_columns": sdk.ParseCommaSeparatedStringArray(idx.IncludedColumns, false),
		})
	}
	return result
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
	return forceNewIfColumnFieldChanged("collate", func(o, n column) bool {
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
	return forceNewIfColumnFieldChanged("nullable", func(o, n column) bool {
		return o.nullable != n.nullable
	})
}

// forceNewIfColumnFieldChanged returns a CustomizeDiffFunc that forces recreation
// when a nested column field changes. It must call diff.ForceNew on the specific
// nested path (e.g. "column.0.nullable"), not on the parent "column" list — the
// terraform-plugin-sdk/v2 ForceNew sets RequiresNew on the resolved leaf schema,
// and a TypeList parent does not propagate that flag down to its diff entries.
func forceNewIfColumnFieldChanged(fieldName string, changed func(old, new column) bool) schema.CustomizeDiffFunc {
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
