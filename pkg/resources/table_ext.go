package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	tableCreateExt = PreviewFeatureCreateContextWrapper(string(previewfeatures.TableResource), TrackingCreateWrapper(resources.Table, CreateTable))
	tableReadExt   = PreviewFeatureReadContextWrapper(string(previewfeatures.TableResource), TrackingReadWrapper(resources.Table, ReadTableFunc(true)))
	tableUpdateExt = PreviewFeatureUpdateContextWrapper(string(previewfeatures.TableResource), TrackingUpdateWrapper(resources.Table, UpdateTable))

	tableDescriptionExt = docs.WithPreviewFeatureMark() + "Resource used to manage tables. For more information, check [CREATE TABLE documentation](https://docs.snowflake.com/en/sql-reference/sql/create-table)."

	tableImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.Table, importTableExt),
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

func tableNewCreateRequestExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier) (*sdk.CreateTableRequest, error) {
	tableColumnRequests, err := getTableColumnRequests(d.Get("column").([]any))
	if err != nil {
		return nil, err
	}

	return sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(tableColumnRequests)), nil
}

func toTableClusterBy(value []any) ([]string, error) {
	return expandStringList(value), nil
}

func tableSetConfigFieldsExt(_ context.Context, d *schema.ResourceData, _ any, table *sdk.Table) error {
	return errors.Join(
		d.Set("comment", table.Comment),
		d.Set("cluster_by", table.GetClusterByKeys()),
	)
}

func tableHandleExternalChangesExt(d *schema.ResourceData, table *sdk.Table, withExternalChangesMarking bool) error {
	if !withExternalChangesMarking {
		return nil
	}
	return handleExternalChangesToObjectInShow(d, tableOutputMappingsExt(table)...)
}

func tableOutputMappingsExt(table *sdk.Table) []outputMapping {
	return []outputMapping{
		{"change_tracking", "change_tracking", table.ChangeTracking, booleanStringFromBool(table.ChangeTracking), nil},
	}
}

func tableSetDescribeOutputExt(d *schema.ResourceData, _ sdk.SchemaObjectIdentifier, tableDescription []sdk.TableColumnDetails, describeErr error, _ bool) error {
	if describeErr != nil {
		return describeErr
	}
	// TODO: stop copying DESCRIBE onto column on every Read (nested-block drift, like authentication policy). Follow-up.
	return errors.Join(
		d.Set("column", toColumnConfig(tableDescription)),
		d.Set(DescribeOutputAttributeName, schemas.TableColumnDetailsListToSchema(tableDescription)),
	)
}

func importTableExt(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	results, err := ImportTable(ctx, d, meta)
	if err != nil {
		return nil, err
	}
	id, err := tableParseId(d.Id())
	if err != nil {
		return nil, err
	}
	details, err := tableDescribeInSdkExt(ctx, meta, id)
	if err != nil {
		return nil, err
	}
	if err := d.Set("column", toColumnConfig(details)); err != nil {
		return nil, err
	}
	return results, nil
}

func tableSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, table *sdk.Table, _ bool) error {
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		d.Set("comment", table.Comment),
		d.Set("change_tracking", booleanStringFromBool(table.ChangeTracking)),
	)
}

func tableAfterAlterExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if diags := tableAfterAlterClusterByExt(ctx, d, meta, id); diags != nil {
		return diags
	}
	if diags := tableAfterAlterColumnsExt(ctx, d, meta, id); diags != nil {
		return diags
	}
	return tableAfterAlterPrimaryKeyExt(ctx, d, meta, id)
}

func tableAfterAlterClusterByExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if !d.HasChange("cluster_by") {
		return nil
	}
	client := meta.(*provider.Context).Client
	cb := expandStringList(d.Get("cluster_by").([]any))
	if len(cb) != 0 {
		if err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithClusteringAction(*sdk.NewTableClusteringActionRequest().WithClusterBy(cb))); err != nil {
			return diag.FromErr(fmt.Errorf("error clustering table %v: %w", id.FullyQualifiedName(), err))
		}
		return nil
	}
	if err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithClusteringAction(*sdk.NewTableClusteringActionRequest().WithDropClusteringKey(true))); err != nil {
		return diag.FromErr(fmt.Errorf("error dropping clustering key on table %v: %w", id.FullyQualifiedName(), err))
	}
	return nil
}

func tableAfterAlterColumnsExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if !d.HasChange("column") {
		return nil
	}
	client := meta.(*provider.Context).Client
	t, n := d.GetChange("column")
	removed, added, changed := getColumns(t).diffs(getColumns(n))

	if len(removed) > 0 {
		removedColumnNames := make([]string, len(removed))
		for i, r := range removed {
			removedColumnNames[i] = r.name
		}
		err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithDropColumns(*sdk.NewTableDropColumnActionRequest(toColumns(removedColumnNames)))))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error dropping columns on table %v: %w", id.FullyQualifiedName(), err))
		}
	}

	for _, cA := range added {
		dataType, err := datatypes.ParseDataType(cA.dataType)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to add column %v: %w", cA.name, err))
		}
		addRequest := sdk.NewTableColumnAddActionRequest(cA.name, dataType)
		if !cA.nullable {
			addRequest.WithNotNull(true)
		}

		if cA._default != nil {
			if cA._default._type() != "constant" {
				return diag.FromErr(fmt.Errorf("failed to add column %v => Only adding a column as a constant is supported by Snowflake", cA.name))
			}
			var expression string
			if datatypes.IsTextDataType(dataType) {
				expression = snowflake.EscapeSnowflakeString(*cA._default.constant)
			} else {
				expression = *cA._default.constant
			}
			addRequest.WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithExpression(expression))
		}

		if cA.identity != nil {
			addRequest.WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithIdentity(*sdk.NewColumnIdentityRequest(cA.identity.startNum, cA.identity.stepNum)))
		}

		if cA.maskingPolicy != "" {
			addRequest.WithMaskingPolicy(*sdk.NewTableColumnMaskingPolicyRequest(sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(cA.maskingPolicy)))
		}

		if cA.comment != "" {
			addRequest.WithComment(cA.comment)
		}

		if cA.collate != "" && datatypes.IsTextDataType(dataType) {
			addRequest.WithCollate(cA.collate)
		}

		err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithAdd(*addRequest)))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error adding column %v on table %v: %w", cA.name, id.FullyQualifiedName(), err))
		}
	}
	for _, cA := range changed {
		if cA.changedDataType || cA.changedCollate {
			dataType, err := datatypes.ParseDataType(cA.newColumn.dataType)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error changing property on %v: err %w", d.Id(), err))
			}
			dataTypeRequest := sdk.NewTableColumnAlterDataTypeRequest(dataType)
			if datatypes.IsTextDataType(dataType) && cA.newColumn.collate != "" {
				dataTypeRequest.WithCollate(cA.newColumn.collate)
			}
			columnAlterActionRequest := sdk.NewTableColumnAlterActionRequest(cA.newColumn.name).WithDataType(*dataTypeRequest)
			err = client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{*columnAlterActionRequest})))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error altering column type on table %v: %w", id.FullyQualifiedName(), err))
			}
		}
		if cA.changedNullConstraint {
			nullabilityRequest := sdk.NewTableColumnNotNullConstraintRequest()
			if !cA.newColumn.nullable {
				nullabilityRequest.WithSet(true)
			} else {
				nullabilityRequest.WithDrop(true)
			}
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{*sdk.NewTableColumnAlterActionRequest(cA.newColumn.name).WithNotNullConstraint(*nullabilityRequest)})))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error altering column null constraint on table %v: %w", id.FullyQualifiedName(), err))
			}
		}
		if cA.droppedDefault {
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{*sdk.NewTableColumnAlterActionRequest(cA.newColumn.name).WithDropDefault(true)})))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error dropping column default on table %v: %w", id.FullyQualifiedName(), err))
			}
		}
		if cA.changedComment {
			columnAlterActionRequest := sdk.NewTableColumnAlterActionRequest(cA.newColumn.name)
			if cA.newColumn.comment == "" {
				columnAlterActionRequest.WithUnsetComment(true)
			} else {
				columnAlterActionRequest.WithComment(cA.newColumn.comment)
			}

			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{*columnAlterActionRequest})))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error altering column comment on table %v: %w", id.FullyQualifiedName(), err))
			}
		}
		if cA.changedMaskingPolicy {
			columnAction := sdk.NewTableColumnActionRequest()
			if strings.TrimSpace(cA.newColumn.maskingPolicy) == "" {
				columnAction.WithUnsetMaskingPolicy(*sdk.NewTableUnsetColumnMaskingPolicyRequest(cA.newColumn.name))
			} else {
				columnAction.WithSetMaskingPolicy(*sdk.NewTableSetColumnMaskingPolicyRequest(cA.newColumn.name, sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(cA.newColumn.maskingPolicy)).WithForce(true))
			}
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*columnAction))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error altering column masking policy on table %v: %w", id.FullyQualifiedName(), err))
			}
		}
	}
	return nil
}

func tableAfterAlterPrimaryKeyExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if !d.HasChange("primary_key") {
		return nil
	}
	client := meta.(*provider.Context).Client
	o, n := d.GetChange("primary_key")
	newKey := getPrimaryKey(n)
	oldKey := getPrimaryKey(o)

	if len(oldKey.keys) > 0 || len(newKey.keys) == 0 {
		// drop our pk if there was an old primary key, or pk has been removed
		err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithConstraintAction(
			*sdk.NewTableConstraintActionRequest().
				WithDrop(*sdk.NewTableConstraintDropActionRequest().WithPrimaryKey(true)),
		))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error dropping primary key on table %v: %w", id.FullyQualifiedName(), err))
		}
	}

	if len(newKey.keys) > 0 {
		primaryKey := sdk.NewTableOutOfLineUniquePKRequest().WithPrimaryKey(true).WithColumns(toColumns(newKey.keys))
		if newKey.name != "" {
			primaryKey.WithName(newKey.name)
		}
		err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithConstraintAction(
			*sdk.NewTableConstraintActionRequest().WithAdd(*sdk.NewTableOutOfLineConstraintRequest().WithUniquePK(*primaryKey)),
		))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error adding primary key on table %v: %w", id.FullyQualifiedName(), err))
		}
	}
	return nil
}

// =============================================================================
// SDK incision points
// =============================================================================

func tableDescribeInSdkExt(ctx context.Context, meta any, id sdk.SchemaObjectIdentifier) ([]sdk.TableColumnDetails, error) {
	client := meta.(*provider.Context).Client
	return client.Tables.DescribeColumns(ctx, sdk.NewDescribeColumnsTableRequest(id))
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

type columnDefault struct {
	constant   *string
	expression *string
	sequence   *string
}

func (cd *columnDefault) _type() string {
	if cd.constant != nil {
		return "constant"
	}

	if cd.expression != nil {
		return "expression"
	}

	if cd.sequence != nil {
		return "sequence"
	}

	return "unknown"
}

type columnIdentity struct {
	startNum int
	stepNum  int
}

type column struct {
	name          string
	dataType      string
	nullable      bool
	_default      *columnDefault
	identity      *columnIdentity
	comment       string
	maskingPolicy string
	collate       string
}

type columns []column

type changedColumns []changedColumn

type changedColumn struct {
	newColumn             column // our new column
	changedDataType       bool
	changedNullConstraint bool
	droppedDefault        bool
	changedComment        bool
	changedMaskingPolicy  bool
	changedCollate        bool
}

func (c columns) getChangedColumnProperties(new columns) (changed changedColumns) {
	changed = make(changedColumns, 0, len(new)*len(c))
	for _, cO := range c {
		for _, cN := range new {
			changeColumn := changedColumn{cN, false, false, false, false, false, false}
			if cO.name == cN.name && cO.dataType != cN.dataType {
				changeColumn.changedDataType = true
			}
			if cO.name == cN.name && cO.nullable != cN.nullable {
				changeColumn.changedNullConstraint = true
			}
			if cO.name == cN.name && cO._default != nil && cN._default == nil {
				changeColumn.droppedDefault = true
			}

			if cO.name == cN.name && cO.comment != cN.comment {
				changeColumn.changedComment = true
			}

			if cO.name == cN.name && cO.maskingPolicy != cN.maskingPolicy {
				changeColumn.changedMaskingPolicy = true
			}

			if cO.name == cN.name && cO.collate != cN.collate {
				changeColumn.changedCollate = true
			}

			changed = append(changed, changeColumn)
		}
	}
	return
}

func (c columns) diffs(new columns) (removed columns, added columns, changed changedColumns) {
	return c.getNewIn(new), new.getNewIn(c), c.getChangedColumnProperties(new)
}

func getColumnDefault(def map[string]any) *columnDefault {
	if c, ok := def["constant"]; ok {
		if constant, ok := c.(string); ok && len(constant) > 0 {
			return &columnDefault{
				constant: &constant,
			}
		}
	}

	if e, ok := def["expression"]; ok {
		if expr, ok := e.(string); ok && len(expr) > 0 {
			return &columnDefault{
				expression: &expr,
			}
		}
	}

	if s, ok := def["sequence"]; ok {
		if seq := s.(string); ok && len(seq) > 0 {
			return &columnDefault{
				sequence: &seq,
			}
		}
	}

	return nil
}

func getColumnIdentity(identity map[string]any) *columnIdentity {
	if len(identity) > 0 {
		startNum := identity["start_num"].(int)
		stepNum := identity["step_num"].(int)
		return &columnIdentity{startNum, stepNum}
	}

	return nil
}

func getColumn(from any) (to column) {
	c := from.(map[string]any)
	var cd *columnDefault
	var id *columnIdentity

	_default := c["default"].([]any)
	identity := c["identity"].([]any)

	if len(_default) == 1 {
		cd = getColumnDefault(_default[0].(map[string]any))
	}
	if len(identity) == 1 {
		id = getColumnIdentity(identity[0].(map[string]any))
	}

	return column{
		name:          c["name"].(string),
		dataType:      c["type"].(string),
		nullable:      c["nullable"].(bool),
		_default:      cd,
		identity:      id,
		comment:       c["comment"].(string),
		collate:       c["collate"].(string),
		maskingPolicy: c["masking_policy"].(string),
	}
}

func getColumns(from any) (to columns) {
	cols := from.([]any)
	to = make(columns, len(cols))
	for i, c := range cols {
		to[i] = getColumn(c)
	}
	return to
}

func getTableColumnRequest(from any) (*sdk.TableColumnRequest, error) {
	c := from.(map[string]any)
	_type := c["type"].(string)
	dataType, err := datatypes.ParseDataType(_type)
	if err != nil {
		return nil, err
	}

	request := sdk.NewTableColumnRequest(c["name"].(string), dataType)

	_default := c["default"].([]any)
	var expression string
	if len(_default) == 1 {
		if c, ok := _default[0].(map[string]any)["constant"]; ok {
			if constant, ok := c.(string); ok && len(constant) > 0 {
				if datatypes.IsTextDataType(dataType) {
					expression = snowflake.EscapeSnowflakeString(constant)
				} else {
					expression = constant
				}
			}
		}

		if e, ok := _default[0].(map[string]any)["expression"]; ok {
			if expr, ok := e.(string); ok && len(expr) > 0 {
				expression = expr
			}
		}

		if s, ok := _default[0].(map[string]any)["sequence"]; ok {
			if seq, ok2 := s.(string); ok2 && len(seq) > 0 {
				expression = fmt.Sprintf(`%v.NEXTVAL`, seq)
			}
		}
		request.WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithExpression(expression))
	}

	identity := c["identity"].([]any)
	if len(identity) == 1 {
		identityProp := identity[0].(map[string]any)
		startNum := identityProp["start_num"].(int)
		stepNum := identityProp["step_num"].(int)
		request.WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithIdentity(*sdk.NewColumnIdentityRequest(startNum, stepNum)))
	}

	maskingPolicy := c["masking_policy"].(string)
	if maskingPolicy != "" {
		request.WithMaskingPolicy(*sdk.NewTableColumnMaskingPolicyRequest(sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(maskingPolicy)))
	}

	if datatypes.IsTextDataType(dataType) {
		request.WithCollate(c["collate"].(string))
	}

	return request.
		WithNotNull(!c["nullable"].(bool)).
		WithComment(c["comment"].(string)), nil
}

func getTableColumnRequests(from any) ([]sdk.TableColumnRequest, error) {
	cols := from.([]any)
	to := make([]sdk.TableColumnRequest, len(cols))
	for i, c := range cols {
		cReq, err := getTableColumnRequest(c)
		if err != nil {
			return nil, err
		}
		to[i] = *cReq
	}
	return to, nil
}

type primarykey struct {
	name string
	keys []string
}

func getPrimaryKey(from any) (to primarykey) {
	pk := from.([]any)
	to = primarykey{}
	if len(pk) > 0 {
		pkDetails := pk[0].(map[string]any)
		to.name = pkDetails["name"].(string)
		to.keys = expandStringList(pkDetails["keys"].([]any))
		return to
	}
	return to
}

func toColumnConfig(descriptions []sdk.TableColumnDetails) []any {
	flattened := make([]any, 0)
	for _, td := range descriptions {
		if td.Kind != "COLUMN" {
			continue
		}

		flat := map[string]any{}
		flat["name"] = td.Name
		flat["type"] = td.TypeRaw
		flat["nullable"] = td.Null

		if td.Comment != nil {
			flat["comment"] = *td.Comment
		}

		if td.Collation != nil {
			flat["collate"] = *td.Collation
		}

		if td.PolicyName != nil {
			// TODO [SNOW-867240]: SHOW TABLE returns last part of id without double quotes... we have to quote it again. Move it to SDK.
			flat["masking_policy"] = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(*td.PolicyName).FullyQualifiedName()
		}

		identity := toColumnIdentityConfig(td)
		if identity != nil {
			flat["identity"] = []any{identity}
		} else {
			def := toColumnDefaultConfig(td)
			if def != nil {
				flat["default"] = []any{def}
			}
		}

		if td.SchemaEvolutionRecord != nil {
			flat["schema_evolution_record"] = *td.SchemaEvolutionRecord
		}

		flattened = append(flattened, flat)
	}
	return flattened
}

func toColumnDefaultConfig(td sdk.TableColumnDetails) map[string]any {
	if td.Default == nil {
		return nil
	}

	defaultRaw := *td.Default
	def := map[string]any{}
	if before, ok := strings.CutSuffix(defaultRaw, ".NEXTVAL"); ok {
		// TODO [SNOW-867240]: SHOW TABLE returns last part of id without double quotes... we have to quote it again. Move it to SDK.
		sequenceIdRaw := before
		def["sequence"] = sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(sequenceIdRaw).FullyQualifiedName()
		return def
	}

	if strings.Contains(defaultRaw, "(") && strings.Contains(defaultRaw, ")") {
		def["expression"] = defaultRaw
		return def
	}

	if td.Type != nil && datatypes.IsTextDataType(td.Type) {
		def["constant"] = snowflake.UnescapeSnowflakeString(defaultRaw)
		return def
	}

	def["constant"] = defaultRaw
	return def
}

func toColumnIdentityConfig(td sdk.TableColumnDetails) map[string]any {
	// if autoincrement is used this is reflected back IDENTITY START 1 INCREMENT 1
	if td.Default == nil {
		return nil
	}

	defaultRaw := *td.Default

	if strings.Contains(defaultRaw, "IDENTITY") {
		identity := map[string]any{}

		split := strings.Split(defaultRaw, " ")
		start, err := strconv.Atoi(split[2])
		if err == nil {
			identity["start_num"] = start
		}
		step, err := strconv.Atoi(split[4])
		if err == nil {
			identity["step_num"] = step
		}

		return identity
	}
	return nil
}
