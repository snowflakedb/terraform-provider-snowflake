package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	tableCreateExt = PreviewFeatureCreateContextWrapper(string(previewfeatures.TableResource), TrackingCreateWrapper(resources.Table, CreateTable))
	tableReadExt   = PreviewFeatureReadContextWrapper(string(previewfeatures.TableResource), TrackingReadWrapper(resources.Table, ReadTableFunc(true)))
	tableUpdateExt = PreviewFeatureUpdateContextWrapper(string(previewfeatures.TableResource), TrackingUpdateWrapper(resources.Table, UpdateTable))

	tableDescriptionExt   = ""
	tableCustomizeDiffExt = TrackingCustomDiffWrapper(resources.Table, customdiff.All(
		TemporaryWorkaroundIdentifierForceNewIfHierarchyRenamesExperimentNotEnabled("database"),
		TemporaryWorkaroundIdentifierForceNewIfHierarchyRenamesExperimentNotEnabled("schema"),
		ComputedIfAnyAttributeChanged(tableSchema, FullyQualifiedNameAttributeName, "name", "database", "schema"),
	))

	tableImporterExt = &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

type tableAlterRequestsExt struct {
	set   *sdk.TableSetRequest
	unset *sdk.TableUnsetRequest
}

func tableNewCreateRequestExt(d *schema.ResourceData, id sdk.SchemaObjectIdentifier) (*sdk.CreateTableRequest, error) {
	tableColumnRequests, err := getTableColumnRequests(d.Get("column").([]any))
	if err != nil {
		return nil, err
	}

	request := sdk.NewCreateTableRequest(id, *sdk.NewCreateTableColumnsAndConstraintsRequest().WithColumns(tableColumnRequests))
	// tag is not an attribute, so generated ApplyCreateOptionals does not see it.
	if _, ok := d.GetOk("tag"); ok {
		request.WithTag(getPropertyTags(d, "tag"))
	}
	return request, nil
}

func toTableClusterBy(value []any) ([]string, error) {
	return expandStringList(value), nil
}

func tableSetConfigFieldsExt(d *schema.ResourceData, table *sdk.Table) error {
	// Set the relevant data in the state
	toSet := map[string]any{
		"name":            table.Name,
		"owner":           table.Owner,
		"database":        table.DatabaseName,
		"schema":          table.SchemaName,
		"comment":         table.Comment,
		"cluster_by":      table.GetClusterByKeys(),
		"change_tracking": table.ChangeTracking,
	}

	for key, val := range toSet {
		if err := d.Set(key, val); err != nil { // lintignore:R001
			return err
		}
	}
	return nil
}

func tableSetDescribeOutputExt(d *schema.ResourceData, _ sdk.SchemaObjectIdentifier, tableDescription []sdk.TableColumnDetails, _ error, _ bool) error {
	return d.Set("column", toColumnConfig(tableDescription))
}

func tableHierarchyRenameExt(ctx context.Context, d *schema.ResourceData, meta any, id *sdk.SchemaObjectIdentifier) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client

	if providerCtx.Experiments.IsEnabled(experimentalfeatures.HierarchyRenames) && (d.HasChange("database") || d.HasChange("schema")) {
		tableRenameFn := func(currentId, targetId sdk.SchemaObjectIdentifier) func() error {
			return func() error {
				return client.Tables.Alter(ctx, sdk.NewAlterTableRequest(currentId).WithRenameTo(targetId))
			}
		}

		if diags := handleThreeLevelHierarchyRename(
			ctx, d, client, id,
			tableRenameFn,
			client.Tables.ShowByID,
			func(id sdk.SchemaObjectIdentifier) string { return helpers.EncodeSnowflakeID(id) },
			"table",
		); diags != nil {
			return diags
		}
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(*id).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming table %v err = %w", d.Id(), err))
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		*id = newId
	}
	return nil
}

// tableAfterAlterExt runs the ALTER TABLE statements issued after SET/UNSET: clustering, columns, primary key, and tags.
func tableAfterAlterExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if d.HasChange("cluster_by") {
		cb := expandStringList(d.Get("cluster_by").([]any))

		if len(cb) != 0 {
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithClusteringAction(*sdk.NewTableClusteringActionRequest().WithClusterBy(cb)))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error updating table: %w", err))
			}
		} else {
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithClusteringAction(*sdk.NewTableClusteringActionRequest().WithDropClusteringKey(true)))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error updating table: %w", err))
			}
		}
	}

	if d.HasChange("column") {
		t, n := d.GetChange("column")
		removed, added, changed := getColumns(t).diffs(getColumns(n))

		if len(removed) > 0 {
			removedColumnNames := make([]string, len(removed))
			for i, r := range removed {
				removedColumnNames[i] = r.name
			}
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithDropColumns(*sdk.NewTableDropColumnActionRequest(toColumns(removedColumnNames)))))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error updating table: %w", err))
			}
		}

		for _, cA := range added {
			addRequest := sdk.NewTableColumnAddActionRequest(cA.name, sdk.DataType(cA.dataType)).
				WithInlineConstraint(*sdk.NewTableColumnAddInlineConstraintRequest("").WithNotNull(!cA.nullable))

			if cA._default != nil {
				if cA._default._type() != "constant" {
					return diag.FromErr(fmt.Errorf("failed to add column %v => Only adding a column as a constant is supported by Snowflake", cA.name))
				}
				var expression string
				if sdk.IsStringType(cA.dataType) {
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

			if cA.collate != "" && sdk.IsStringType(cA.dataType) {
				addRequest.WithCollate(cA.collate)
			}

			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithAdd(*addRequest)))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error adding column: %w", err))
			}
		}
		for _, cA := range changed {
			if cA.changedDataType || cA.changedCollate {
				columnAlterActionRequest := sdk.NewTableColumnAlterActionRequest(cA.newColumn.name).WithDataType(sdk.DataType(cA.newColumn.dataType))
				if sdk.IsStringType(cA.newColumn.dataType) && cA.newColumn.collate != "" {
					columnAlterActionRequest.WithCollate(cA.newColumn.collate)
				}
				err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{*columnAlterActionRequest})))
				if err != nil {
					return diag.FromErr(fmt.Errorf("error changing property on %v: err %w", d.Id(), err))
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
					return diag.FromErr(fmt.Errorf("error changing property on %v: err %w", d.Id(), err))
				}
			}
			if cA.droppedDefault {
				err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithColumnAction(*sdk.NewTableColumnActionRequest().WithAlter([]sdk.TableColumnAlterActionRequest{*sdk.NewTableColumnAlterActionRequest(cA.newColumn.name).WithDropDefault(true)})))
				if err != nil {
					return diag.FromErr(fmt.Errorf("error changing property on %v: err %w", d.Id(), err))
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
					return diag.FromErr(fmt.Errorf("error changing property on %v: err %w", d.Id(), err))
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
					return diag.FromErr(fmt.Errorf("error changing property on %v: err %w", d.Id(), err))
				}
			}
		}
	}

	if d.HasChange("primary_key") {
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
				return diag.FromErr(fmt.Errorf("error updating table: %w", err))
			}
		}

		if len(newKey.keys) > 0 {
			constraint := sdk.NewOutOfLineConstraintRequest(sdk.ColumnConstraintTypePrimaryKey).WithColumns(toColumns(newKey.keys))
			if newKey.name != "" {
				constraint.WithName(newKey.name)
			}
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithConstraintAction(
				*sdk.NewTableConstraintActionRequest().WithAdd(*constraint),
			))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error updating table: %w", err))
			}
		}
	}

	if d.HasChange("tag") {
		unsetTags, setTags := GetTagsDiff(d, "tag")

		if len(unsetTags) > 0 {
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithUnsetTags(unsetTags))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting tags on %v, err = %w", d.Id(), err))
			}
		}

		if len(setTags) > 0 {
			err := client.Tables.Alter(ctx, sdk.NewAlterTableRequest(id).WithSetTags(setTags))
			if err != nil {
				return diag.FromErr(fmt.Errorf("error setting tags on %v, err = %w", d.Id(), err))
			}
		}
	}
	return nil
}

// =============================================================================
// SDK incision points
// =============================================================================

func tableShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.SchemaObjectIdentifier) (*sdk.Table, diag.Diagnostics) {
	client := meta.(*provider.Context).Client

	table, err := client.Tables.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query table. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Table id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}

	s, err := client.Schemas.ShowByID(ctx, sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName()))
	if err != nil {
		log.Printf("[DEBUG] schema (%s) not found", d.Id())
		d.SetId("")
		return nil, diag.Diagnostics{}
	}

	// "retention_time" may sometimes be empty string instead of an integer.
	retentionTime := s.RetentionTime
	if retentionTime == "" {
		retentionTime = "0"
	}
	schemaRetentionTime, err := strconv.ParseInt(retentionTime, 10, 64)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	if v := d.Get("data_retention_time_in_days"); v.(int) != IntDefault || int64(table.RetentionTime) != schemaRetentionTime {
		if err := d.Set("data_retention_time_in_days", table.RetentionTime); err != nil {
			return nil, diag.FromErr(err)
		}
	}
	return table, nil
}

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

	request := sdk.NewTableColumnRequest(c["name"].(string), sdk.DataType(_type))

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
		flat["type"] = string(td.Type)
		flat["nullable"] = td.IsNullable

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

	if sdk.IsStringType(string(td.Type)) {
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
