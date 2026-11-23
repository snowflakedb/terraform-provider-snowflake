package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
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
		StateContext: TrackingImportWrapper(resources.Table, ImportTable),
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
	request := sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(tableColumnRequests))
	// is_transient is not sent by generated ApplyCreateOptionals: TRANSIENT is the table kind of the CREATE statement.
	if err := booleanStringAttributeCreateBuilder(d, "is_transient", func(isTransient bool) *sdk.CreateTableRequest {
		if isTransient {
			return request.WithKind(sdk.TableKindTransient)
		}
		return request
	}); err != nil {
		return nil, err
	}

	if policyId, columns, ok, err := rowAccessPolicyCreateRequest(d); err != nil {
		return nil, err
	} else if ok {
		request.WithRowAccessPolicy(*sdk.NewTableRowAccessPolicyRequest(policyId, columns))
	}
	return request, nil
}

func toTableClusterBy(value []any) ([]string, error) {
	return expandStringList(value), nil
}

func toTableColumn(value any) (sdk.Column, error) {
	s, ok := value.(string)
	if !ok {
		return sdk.Column{}, fmt.Errorf("expected column name string, got %T", value)
	}
	return sdk.Column{Value: s}, nil
}

func tableSetConfigFieldsExt(ctx context.Context, d *schema.ResourceData, meta any, table *sdk.Table) error {
	if err := errors.Join(
		d.Set("comment", table.Comment),
		d.Set("cluster_by", table.GetClusterByKeys()),
	); err != nil {
		return err
	}

	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return err
	}
	policyRefs, err := readRootLevelPolicies(ctx, meta.(*provider.Context).Client, id, sdk.PolicyEntityDomainTable, d,
		sdk.PolicyKindRowAccessPolicy, sdk.PolicyKindAggregationPolicy, sdk.PolicyKindJoinPolicy)
	if err != nil {
		return err
	}
	columns, err := withColumnMaskingPolicies(d.Get("column").([]any), policyRefs)
	if err != nil {
		return err
	}
	return d.Set("column", columns)
}

func tableSetDescribeOutputExt(d *schema.ResourceData, _ sdk.SchemaObjectIdentifier, tableDescription []sdk.TableColumnDetails, describeErr error, _ bool) error {
	if describeErr != nil {
		return describeErr
	}
	// Columns are rebuilt from DESCRIBE TABLE on every Read, like in snowflake_iceberg_table. DESCRIBE is the only source
	// of the columns, so this is how columns added, dropped, or changed outside of Terraform show up in the plan.
	// Comparing describe_output instead would miss changes it does not hold, such as masking policy USING columns.
	// Formatting differences are handled by the diff suppressors of the column fields.
	return errors.Join(
		d.Set("column", toColumnConfig(tableDescription)),
		d.Set(DescribeOutputAttributeName, schemas.TableColumnDetailsListToSchema(tableDescription)),
	)
}

// tableSetFieldsNotSetByReadExt sets the fields Read does not write from the configuration. Import calls it. In enforcing
// mode the optional boolean fields get the values Snowflake reports; in permissive mode they get their unset values.
func tableSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier, table *sdk.Table, importPermissiveMode bool) error {
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		importBooleanAttribute(d, "change_tracking", importPermissiveMode, table.ChangeTracking),
		importBooleanStringAttribute(d, "enable_schema_evolution", importPermissiveMode, table.EnableSchemaEvolution),
		importBooleanStringAttribute(d, "error_logging", importPermissiveMode, table.ErrorLogging),
		importBooleanStringAttribute(d, "row_timestamp", importPermissiveMode, table.RowTimestamp),
		importBooleanStringAttribute(d, "is_transient", importPermissiveMode, table.Kind == string(sdk.TableKindTransient)),
	)
}

func tableAfterAlterExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	if diags := tableAfterAlterClusterByExt(ctx, d, meta, id); diags != nil {
		return diags
	}
	if diags := tableAfterAlterColumnsExt(ctx, d, meta, id); diags != nil {
		return diags
	}
	if diags := tableAfterAlterPrimaryKeyExt(ctx, d, meta, id); diags != nil {
		return diags
	}
	client := meta.(*provider.Context).Client
	if d.HasChange("row_access_policy") {
		addReq, dropReq, err := rowAccessPolicyUpdateRequests(d)
		if err != nil {
			return diag.FromErr(err)
		}
		alterReq := sdk.NewAlterTableRequest(id)
		switch {
		case addReq != nil && dropReq != nil:
			alterReq.WithDropAndAddRowAccessPolicy(*sdk.NewTableDropAndAddRowAccessPolicyRequest(
				*sdk.NewTableDropRowAccessPolicyRequest(dropReq.RowAccessPolicy),
				*sdk.NewTableAddRowAccessPolicyRequest(addReq.RowAccessPolicy, addReq.On),
			))
		case addReq != nil:
			alterReq.WithAddRowAccessPolicy(*sdk.NewTableAddRowAccessPolicyRequest(addReq.RowAccessPolicy, addReq.On))
		case dropReq != nil:
			alterReq.WithDropRowAccessPolicy(*sdk.NewTableDropRowAccessPolicyRequest(dropReq.RowAccessPolicy))
		}
		if err := client.Tables.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error updating row access policy on %v: %w", d.Id(), err))
		}
	}

	if d.HasChange("aggregation_policy") {
		setReq, unsetReq, err := aggregationPolicyTableUpdateRequests(d)
		if err != nil {
			return diag.FromErr(err)
		}
		alterReq := sdk.NewAlterTableRequest(id)
		if setReq != nil {
			alterReq.WithSetAggregationPolicy(*setReq)
		} else {
			alterReq.WithUnsetAggregationPolicy(*unsetReq)
		}
		if err := client.Tables.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error updating aggregation policy on %v: %w", d.Id(), err))
		}
	}

	if d.HasChange("join_policy") {
		setReq, unsetReq, err := joinPolicyTableUpdateRequests(d)
		if err != nil {
			return diag.FromErr(err)
		}
		alterReq := sdk.NewAlterTableRequest(id)
		if setReq != nil {
			alterReq.WithSetJoinPolicy(*setReq)
		} else {
			alterReq.WithUnsetJoinPolicy(*unsetReq)
		}
		if err := client.Tables.Alter(ctx, alterReq); err != nil {
			return diag.FromErr(fmt.Errorf("error updating join policy on %v: %w", d.Id(), err))
		}
	}
	return nil
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

		if cA._default != nil && cA._default.expression != nil {
			addRequest.WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithExpression(*cA._default.expression))
		}

		if cA.identity != nil {
			addRequest.WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithIdentity(*sdk.NewColumnIdentityRequest(cA.identity.startNum, cA.identity.stepNum)))
		}

		if cA.maskingPolicy != nil {
			maskingPolicyId, err := sdk.ParseSchemaObjectIdentifier(cA.maskingPolicy.name)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to add column %v: %w", cA.name, err))
			}
			addRequest.WithMaskingPolicy(*sdk.NewTableColumnMaskingPolicyRequest(maskingPolicyId).WithUsing(toColumns(cA.maskingPolicy.using)))
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
			if cA.newColumn.maskingPolicy == nil {
				columnAction.WithUnsetMaskingPolicy(*sdk.NewTableUnsetColumnMaskingPolicyRequest(cA.newColumn.name))
			} else {
				maskingPolicyId, err := sdk.ParseSchemaObjectIdentifier(cA.newColumn.maskingPolicy.name)
				if err != nil {
					return diag.FromErr(fmt.Errorf("error altering column masking policy on table %v: %w", id.FullyQualifiedName(), err))
				}
				setRequest := sdk.NewTableSetColumnMaskingPolicyRequest(cA.newColumn.name, maskingPolicyId).WithForce(true)
				if len(cA.newColumn.maskingPolicy.using) > 0 {
					setRequest.WithUsing(toColumns(cA.newColumn.maskingPolicy.using))
				}
				columnAction.WithSetMaskingPolicy(*setRequest)
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
	maskingPolicy *columnMaskingPolicy
	collate       string
}

type columnMaskingPolicy struct {
	name  string
	using []string
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

			if cO.name == cN.name && !cO.maskingPolicy.equal(cN.maskingPolicy) {
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

// columnDefaultEmptyStringLiteral is the SQL literal of an empty string. It is what Snowflake stores and
// DESCRIBE TABLE returns for a column created with an empty string default (two single quotes).
const columnDefaultEmptyStringLiteral = "''"

// getColumnDefaultExpression returns the expression of a column default block. The block is optional, but its
// expression is required, so a block that is present always has a default. An empty expression is sent as the
// empty string literal, because `DEFAULT` without an expression is not valid SQL.
//
// SDKv2 decodes a block whose fields are all zero values (here: `expression = ""`) as a nil element, so the element
// is not type-asserted directly.
func getColumnDefaultExpression(defaults []any) (string, bool) {
	if len(defaults) != 1 {
		return "", false
	}
	def, ok := defaults[0].(map[string]any)
	if !ok {
		return columnDefaultEmptyStringLiteral, true
	}
	expression, _ := def["expression"].(string)
	if expression == "" {
		return columnDefaultEmptyStringLiteral, true
	}
	return expression, true
}

func getColumnMaskingPolicy(maskingPolicies []any) *columnMaskingPolicy {
	if len(maskingPolicies) != 1 {
		return nil
	}
	maskingPolicy, ok := maskingPolicies[0].(map[string]any)
	if !ok {
		return nil
	}
	name, _ := maskingPolicy["policy_name"].(string)
	if name == "" {
		return nil
	}
	var using []string
	if usingRaw, ok := maskingPolicy["using"].([]any); ok {
		using = expandStringList(usingRaw)
	}
	return &columnMaskingPolicy{name: name, using: using}
}

// equal compares the policies as they are in the state. Formatting differences (identifier quoting, the column itself as
// the only USING argument) are handled by the diff suppressors of the masking_policy schema, so they never reach here.
func (p *columnMaskingPolicy) equal(other *columnMaskingPolicy) bool {
	if p == nil || other == nil {
		return p == other
	}
	return p.name == other.name && slices.Equal(p.using, other.using)
}

// withColumnMaskingPolicies sets each column's masking_policy from POLICY_REFERENCES, like in snowflake_iceberg_table.
func withColumnMaskingPolicies(columns []any, policyRefs []sdk.PolicyReference) ([]any, error) {
	var errs []error
	for _, c := range columns {
		column, ok := c.(map[string]any)
		if !ok {
			continue
		}
		name, _ := column["name"].(string)
		columnPolicies, err := columnPoliciesToState(name, policyRefs)
		if err != nil {
			errs = append(errs, fmt.Errorf("column %q: %w", name, err))
		}
		if maskingPolicy, ok := columnPolicies["masking_policy"]; ok {
			column["masking_policy"] = maskingPolicy
		} else {
			column["masking_policy"] = []any{}
		}
	}
	return columns, errors.Join(errs...)
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

	identity := c["identity"].([]any)

	if expression, ok := getColumnDefaultExpression(c["default"].([]any)); ok {
		cd = &columnDefault{expression: &expression}
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
		maskingPolicy: getColumnMaskingPolicy(c["masking_policy"].([]any)),
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

	if expression, ok := getColumnDefaultExpression(c["default"].([]any)); ok {
		request.WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithExpression(expression))
	}

	identity := c["identity"].([]any)
	if len(identity) == 1 {
		identityProp := identity[0].(map[string]any)
		startNum := identityProp["start_num"].(int)
		stepNum := identityProp["step_num"].(int)
		request.WithDefaultValue(*sdk.NewColumnDefaultValueRequest().WithIdentity(*sdk.NewColumnIdentityRequest(startNum, stepNum)))
	}

	if maskingPolicy := getColumnMaskingPolicy(c["masking_policy"].([]any)); maskingPolicy != nil {
		maskingPolicyId, err := sdk.ParseSchemaObjectIdentifier(maskingPolicy.name)
		if err != nil {
			return nil, err
		}
		request.WithMaskingPolicy(*sdk.NewTableColumnMaskingPolicyRequest(maskingPolicyId).WithUsing(toColumns(maskingPolicy.using)))
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

		identity := toColumnIdentityConfig(td)
		if identity != nil {
			flat["identity"] = []any{identity}
		} else {
			def := toColumnDefaultConfig(td)
			if def != nil {
				flat["default"] = []any{def}
			}
		}

		flattened = append(flattened, flat)
	}
	return flattened
}

func toColumnDefaultConfig(td sdk.TableColumnDetails) map[string]any {
	if td.Default == nil {
		return nil
	}

	return map[string]any{"expression": *td.Default}
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
