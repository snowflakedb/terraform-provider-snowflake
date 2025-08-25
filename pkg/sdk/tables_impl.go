package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ Tables                             = (*tables)(nil)
	_ convertibleRow[TableColumnDetails] = new(tableColumnDetailsRow)
	_ convertibleRow[TableStageDetails]  = new(tableStageDetailsRow)
	_ convertibleRow[Table]              = new(tableDBRow)
)

var (
	_ optionsProvider[createTableOptions]              = new(CreateTableRequest)
	_ optionsProvider[createTableAsSelectOptions]      = new(CreateTableAsSelectRequest)
	_ optionsProvider[createTableUsingTemplateOptions] = new(CreateTableUsingTemplateRequest)
	_ optionsProvider[createTableLikeOptions]          = new(CreateTableLikeRequest)
	_ optionsProvider[createTableCloneOptions]         = new(CreateTableCloneRequest)
	_ optionsProvider[alterTableOptions]               = new(AlterTableRequest)
	_ optionsProvider[dropTableOptions]                = new(DropTableRequest)
	_ optionsProvider[showTableOptions]                = new(ShowTableRequest)
	_ optionsProvider[describeTableColumnsOptions]     = new(DescribeTableColumnsRequest)
	_ optionsProvider[describeTableStageOptions]       = new(DescribeTableStageRequest)
	_ optionsProvider[TableColumnAction]               = new(TableColumnActionRequest)
	_ optionsProvider[TableConstraintAction]           = new(TableConstraintActionRequest)
	_ optionsProvider[TableExternalTableAction]        = new(TableExternalTableActionRequest)
	_ optionsProvider[TableSearchOptimizationAction]   = new(TableSearchOptimizationActionRequest)
	_ optionsProvider[TableSet]                        = new(TableSetRequest)
)

type tables struct {
	client *Client
}

func (v *tables) Create(ctx context.Context, request *CreateTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) CreateAsSelect(ctx context.Context, request *CreateTableAsSelectRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) CreateUsingTemplate(ctx context.Context, request *CreateTableUsingTemplateRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) CreateLike(ctx context.Context, request *CreateTableLikeRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) CreateClone(ctx context.Context, request *CreateTableCloneRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) Alter(ctx context.Context, request *AlterTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) Drop(ctx context.Context, request *DropTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *tables) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropTableRequest(id).WithIfExists(Bool(true))) }, ctx, id)
}

func (v *tables) Show(ctx context.Context, request *ShowTableRequest) ([]Table, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[tableDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}

	return convertRows[tableDBRow, Table](dbRows)
}

func (v *tables) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Table, error) {
	request := NewShowTableRequest().WithIn(ExtendedIn{In: In{Schema: id.SchemaId()}}).
		WithLike(Like{Pattern: String(id.Name())})
	returnedTables, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(returnedTables, func(r Table) bool { return r.Name == id.Name() })
}

func (v *tables) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Table, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *tables) DescribeColumns(ctx context.Context, req *DescribeTableColumnsRequest) ([]TableColumnDetails, error) {
	rows, err := validateAndQuery[tableColumnDetailsRow](v.client, ctx, req.toOpts())
	if err != nil {
		return nil, err
	}
	return convertRows[tableColumnDetailsRow, TableColumnDetails](rows)
}

func (v *tables) DescribeStage(ctx context.Context, req *DescribeTableStageRequest) ([]TableStageDetails, error) {
	rows, err := validateAndQuery[tableStageDetailsRow](v.client, ctx, req.toOpts())
	if err != nil {
		return nil, err
	}
	return convertRows[tableStageDetailsRow, TableStageDetails](rows)
}

func (s *AlterTableRequest) toOpts() *alterTableOptions {
	var clusteringAction *TableClusteringAction
	if s.ClusteringAction != nil {
		var reclusterAction *TableReclusterAction
		if s.ClusteringAction.Recluster != nil {
			reclusterAction = &TableReclusterAction{
				MaxSize:   s.ClusteringAction.Recluster.MaxSize,
				Condition: s.ClusteringAction.Recluster.Condition,
			}
		}
		var changeReclusterChange *TableReclusterChangeState
		if s.ClusteringAction.ChangeReclusterState != nil {
			changeReclusterChange = &TableReclusterChangeState{State: s.ClusteringAction.ChangeReclusterState}
		}
		clusteringAction = &TableClusteringAction{
			ClusterBy:            s.ClusteringAction.ClusterBy,
			Recluster:            reclusterAction,
			ChangeReclusterState: changeReclusterChange,
			DropClusteringKey:    s.ClusteringAction.DropClusteringKey,
		}
	}
	var columnAction *TableColumnAction
	if s.ColumnAction != nil {
		columnAction = s.ColumnAction.toOpts()
	}
	var constraintAction *TableConstraintAction
	if s.ConstraintAction != nil {
		constraintAction = s.ConstraintAction.toOpts()
	}
	var externalTableAction *TableExternalTableAction
	if s.ExternalTableAction != nil {
		externalTableAction = s.ExternalTableAction.toOpts()
	}
	var searchOptimizationAction *TableSearchOptimizationAction
	if s.SearchOptimizationAction != nil {
		searchOptimizationAction = s.SearchOptimizationAction.toOpts()
	}
	var tableSet *TableSet
	if s.Set != nil {
		tableSet = s.Set.toOpts()
	}

	tagAssociations := make([]TagAssociation, 0, len(s.SetTags))
	for _, tagRequest := range s.SetTags {
		tagAssociations = append(tagAssociations, TagAssociation(tagRequest))
	}
	var tableUnset *TableUnset
	if s.Unset != nil {
		tableUnset = &TableUnset{
			DataRetentionTimeInDays:    Bool(s.Unset.DataRetentionTimeInDays),
			MaxDataExtensionTimeInDays: Bool(s.Unset.MaxDataExtensionTimeInDays),
			ChangeTracking:             Bool(s.Unset.ChangeTracking),
			DefaultDDLCollation:        Bool(s.Unset.DefaultDDLCollation),
			EnableSchemaEvolution:      Bool(s.Unset.EnableSchemaEvolution),
			Comment:                    Bool(s.Unset.Comment),
		}
	}
	var addRowAccessPolicy *TableAddRowAccessPolicy
	if s.AddRowAccessPolicy != nil {
		addRowAccessPolicy = &TableAddRowAccessPolicy{
			RowAccessPolicy: s.AddRowAccessPolicy.RowAccessPolicy,
			On:              s.AddRowAccessPolicy.On,
		}
	}
	var dropRowAccessPolicy *TableDropRowAccessPolicy
	if s.DropRowAccessPolicy != nil {
		dropRowAccessPolicy = &TableDropRowAccessPolicy{
			RowAccessPolicy: s.DropRowAccessPolicy.RowAccessPolicy,
		}
	}
	var dropAndAddRowAccessPolicy *TableDropAndAddRowAccessPolicy
	if s.DropAndAddRowAccessPolicy != nil {
		add := TableAddRowAccessPolicy{
			RowAccessPolicy: s.DropAndAddRowAccessPolicy.Add.RowAccessPolicy,
			On:              s.DropAndAddRowAccessPolicy.Add.On,
		}
		drop := TableDropRowAccessPolicy{
			RowAccessPolicy: s.DropAndAddRowAccessPolicy.Drop.RowAccessPolicy,
		}
		dropAndAddRowAccessPolicy = &TableDropAndAddRowAccessPolicy{
			Drop: drop,
			Add:  add,
		}
	}

	return &alterTableOptions{
		IfExists:                  s.IfExists,
		name:                      s.name,
		NewName:                   s.NewName,
		SwapWith:                  s.SwapWith,
		ClusteringAction:          clusteringAction,
		ColumnAction:              columnAction,
		ConstraintAction:          constraintAction,
		ExternalTableAction:       externalTableAction,
		SearchOptimizationAction:  searchOptimizationAction,
		Set:                       tableSet,
		SetTags:                   tagAssociations,
		UnsetTags:                 s.UnsetTags,
		Unset:                     tableUnset,
		AddRowAccessPolicy:        addRowAccessPolicy,
		DropRowAccessPolicy:       dropRowAccessPolicy,
		DropAndAddRowAccessPolicy: dropAndAddRowAccessPolicy,
		DropAllAccessRowPolicies:  s.DropAllAccessRowPolicies,
	}
}

func (s *TableSetRequest) toOpts() *TableSet {
	set := &TableSet{
		EnableSchemaEvolution:      s.EnableSchemaEvolution,
		DataRetentionTimeInDays:    s.DataRetentionTimeInDays,
		MaxDataExtensionTimeInDays: s.MaxDataExtensionTimeInDays,
		ChangeTracking:             s.ChangeTracking,
		DefaultDDLCollation:        s.DefaultDDLCollation,
		Comment:                    s.Comment,
	}

	if s.StageCopyOptions != nil {
		set.StageCopyOptions = s.StageCopyOptions.toOpts()
	}
	if s.StageFileFormat != nil {
		set.StageFileFormat = s.StageFileFormat.toOpts()
	}

	return set
}

func (s *TableSearchOptimizationActionRequest) toOpts() *TableSearchOptimizationAction {
	if len(s.AddSearchOptimizationOn) > 0 {
		return &TableSearchOptimizationAction{
			Add: &AddSearchOptimization{
				On: s.AddSearchOptimizationOn,
			},
		}
	}
	if len(s.DropSearchOptimizationOn) > 0 {
		return &TableSearchOptimizationAction{
			Drop: &DropSearchOptimization{
				On: s.DropSearchOptimizationOn,
			},
		}
	}
	return nil
}

func (r *TableExternalTableActionRequest) toOpts() *TableExternalTableAction {
	if r.Add != nil {
		return &TableExternalTableAction{
			Add: &TableExternalTableColumnAddAction{
				IfNotExists: r.Add.IfNotExists,
				Name:        r.Add.Name,
				Type:        r.Add.Type,
				Expression:  []string{r.Add.Expression},
				Comment:     r.Add.Comment,
			},
		}
	}
	if r.Rename != nil {
		return &TableExternalTableAction{
			Rename: &TableExternalTableColumnRenameAction{
				OldName: r.Rename.OldName,
				NewName: r.Rename.NewName,
			},
		}
	}
	if r.Drop != nil {
		return &TableExternalTableAction{
			Drop: &TableExternalTableColumnDropAction{
				Names:    r.Drop.Columns,
				IfExists: r.Drop.IfExists,
			},
		}
	}
	return nil
}

func (r *TableConstraintActionRequest) toOpts() *TableConstraintAction {
	if r.Add != nil {
		var foreignKey *OutOfLineForeignKey
		if r.Add.ForeignKey != nil {
			var foreignKeyOnAction *ForeignKeyOnAction
			if r.Add.ForeignKey.On != nil {
				foreignKeyOnAction = &ForeignKeyOnAction{
					OnUpdate: r.Add.ForeignKey.On.OnUpdate,
					OnDelete: r.Add.ForeignKey.On.OnDelete,
				}
			}
			foreignKey = &OutOfLineForeignKey{
				TableName:   r.Add.ForeignKey.TableName,
				ColumnNames: r.Add.ForeignKey.ColumnNames,
				Match:       r.Add.ForeignKey.Match,
				On:          foreignKeyOnAction,
			}
		}
		outOfLineConstraint := OutOfLineConstraint{
			Name:               r.Add.Name,
			Type:               r.Add.Type,
			Columns:            r.Add.Columns,
			ForeignKey:         foreignKey,
			Enforced:           r.Add.Enforced,
			NotEnforced:        r.Add.NotEnforced,
			Deferrable:         r.Add.Deferrable,
			NotDeferrable:      r.Add.NotDeferrable,
			InitiallyDeferred:  r.Add.InitiallyDeferred,
			InitiallyImmediate: r.Add.InitiallyImmediate,
			Enable:             r.Add.Enable,
			Disable:            r.Add.Disable,
			Validate:           r.Add.Validate,
			NoValidate:         r.Add.NoValidate,
			Rely:               r.Add.Rely,
			NoRely:             r.Add.NoRely,
		}
		return &TableConstraintAction{
			Add: &outOfLineConstraint,
		}
	}
	if r.Rename != nil {
		return &TableConstraintAction{
			Rename: &TableConstraintRenameAction{
				OldName: r.Rename.OldName,
				NewName: r.Rename.NewName,
			},
		}
	}
	if r.Alter != nil {
		return &TableConstraintAction{
			Alter: &TableConstraintAlterAction{
				ConstraintName: r.Alter.ConstraintName,
				PrimaryKey:     r.Alter.PrimaryKey,
				Unique:         r.Alter.Unique,
				ForeignKey:     r.Alter.ForeignKey,
				Columns:        r.Alter.Columns,
				Enforced:       r.Alter.Enforced,
				NotEnforced:    r.Alter.NotEnforced,
				Validate:       r.Alter.Validate,
				NoValidate:     r.Alter.NoValidate,
				Rely:           r.Alter.Rely,
				NoRely:         r.Alter.NoRely,
			},
		}
	}
	if r.Drop != nil {
		return &TableConstraintAction{
			Drop: &TableConstraintDropAction{
				ConstraintName: r.Drop.ConstraintName,
				PrimaryKey:     r.Drop.PrimaryKey,
				Unique:         r.Drop.Unique,
				ForeignKey:     r.Drop.ForeignKey,
				Columns:        r.Drop.Columns,
				Cascade:        r.Drop.Cascade,
				Restrict:       r.Drop.Restrict,
			},
		}
	}
	return nil
}

func (r *TableColumnActionRequest) toOpts() *TableColumnAction {
	if r.Add != nil {
		var defaultValue *ColumnDefaultValue
		if r.Add.DefaultValue != nil {
			defaultValue = &ColumnDefaultValue{
				r.Add.DefaultValue.expression,
				&ColumnIdentity{
					Start:     r.Add.DefaultValue.identity.Start,
					Increment: r.Add.DefaultValue.identity.Increment,
				},
			}
		}
		var inlineConstraint *TableColumnAddInlineConstraint
		if r.Add.InlineConstraint != nil {
			var foreignKey *ColumnAddForeignKey
			if r.Add.InlineConstraint.ForeignKey != nil {
				foreignKey = &ColumnAddForeignKey{
					TableName:  r.Add.InlineConstraint.ForeignKey.TableName,
					ColumnName: r.Add.InlineConstraint.ForeignKey.ColumnName,
				}
			}
			inlineConstraint = &TableColumnAddInlineConstraint{
				NotNull:    r.Add.InlineConstraint.NotNull,
				Name:       r.Add.InlineConstraint.Name,
				Type:       r.Add.InlineConstraint.Type,
				ForeignKey: foreignKey,
			}
		}
		return &TableColumnAction{
			Add: &TableColumnAddAction{
				IfNotExists:      r.Add.IfNotExists,
				Name:             r.Add.Name,
				Type:             r.Add.Type,
				DefaultValue:     defaultValue,
				InlineConstraint: inlineConstraint,
				Comment:          r.Add.Comment,
				Collate:          r.Add.Collate,
			},
		}
	}
	if r.Rename != nil {
		return &TableColumnAction{
			Rename: &TableColumnRenameAction{
				OldName: r.Rename.OldName,
				NewName: r.Rename.NewName,
			},
		}
	}
	if len(r.Alter) > 0 {
		var alterActions []TableColumnAlterAction
		for _, alterAction := range r.Alter {
			var notNullConstraint *TableColumnNotNullConstraint
			if alterAction.NotNullConstraint != nil {
				notNullConstraint = &TableColumnNotNullConstraint{
					Set:  alterAction.NotNullConstraint.Set,
					Drop: alterAction.NotNullConstraint.Drop,
				}
			}
			alterActions = append(alterActions, TableColumnAlterAction{
				Name:              alterAction.Name,
				DropDefault:       alterAction.DropDefault,
				SetDefault:        alterAction.SetDefault,
				NotNullConstraint: notNullConstraint,
				Type:              alterAction.Type,
				Comment:           alterAction.Comment,
				Collate:           alterAction.Collate,
				UnsetComment:      alterAction.UnsetComment,
			})
		}
		return &TableColumnAction{
			Alter: alterActions,
		}
	}
	if r.SetMaskingPolicy != nil {
		return &TableColumnAction{
			SetMaskingPolicy: &TableColumnAlterSetMaskingPolicyAction{
				ColumnName:        r.SetMaskingPolicy.ColumnName,
				MaskingPolicyName: r.SetMaskingPolicy.MaskingPolicyName,
				Using:             r.SetMaskingPolicy.Using,
				Force:             r.SetMaskingPolicy.Force,
			},
		}
	}
	if r.UnsetMaskingPolicy != nil {
		return &TableColumnAction{
			UnsetMaskingPolicy: &TableColumnAlterUnsetMaskingPolicyAction{
				ColumnName: r.UnsetMaskingPolicy.ColumnName,
			},
		}
	}
	if r.SetTags != nil {
		return &TableColumnAction{
			SetTags: &TableColumnAlterSetTagsAction{
				ColumnName: r.SetTags.ColumnName,
				Tags:       r.SetTags.Tags,
			},
		}
	}
	if r.UnsetTags != nil {
		return &TableColumnAction{
			UnsetTags: &TableColumnAlterUnsetTagsAction{
				ColumnName: r.UnsetTags.ColumnName,
				Tags:       r.UnsetTags.Tags,
			},
		}
	}
	if len(r.DropColumns) > 0 {
		return &TableColumnAction{
			DropColumns: &TableColumnAlterDropColumns{
				IfExists: r.DropColumnsIfExists,
				Columns:  r.DropColumns,
			},
		}
	}
	return nil
}

func (s *CreateTableRequest) toOpts() *createTableOptions {
	tagAssociations := make([]TagAssociation, 0, len(s.Tags))
	for _, tagRequest := range s.Tags {
		tagAssociations = append(tagAssociations, TagAssociation(tagRequest))
	}
	var rowAccessPolicy *TableRowAccessPolicy
	if s.RowAccessPolicy != nil {
		rowAccessPolicy = &TableRowAccessPolicy{
			Name: s.RowAccessPolicy.Name,
			On:   s.RowAccessPolicy.On,
		}
	}
	outOfLineConstraints := make([]OutOfLineConstraint, 0)
	for _, outOfLineConstraintRequest := range s.OutOfLineConstraints {
		var foreignKey *OutOfLineForeignKey
		if outOfLineConstraintRequest.ForeignKey != nil {
			var foreignKeyOnAction *ForeignKeyOnAction
			if outOfLineConstraintRequest.ForeignKey.On != nil {
				foreignKeyOnAction = &ForeignKeyOnAction{
					OnUpdate: outOfLineConstraintRequest.ForeignKey.On.OnUpdate,
					OnDelete: outOfLineConstraintRequest.ForeignKey.On.OnDelete,
				}
			}
			foreignKey = &OutOfLineForeignKey{
				TableName:   outOfLineConstraintRequest.ForeignKey.TableName,
				ColumnNames: outOfLineConstraintRequest.ForeignKey.ColumnNames,
				Match:       outOfLineConstraintRequest.ForeignKey.Match,
				On:          foreignKeyOnAction,
			}
		}
		outOfLineConstraint := OutOfLineConstraint{
			Name:               outOfLineConstraintRequest.Name,
			Type:               outOfLineConstraintRequest.Type,
			Columns:            outOfLineConstraintRequest.Columns,
			ForeignKey:         foreignKey,
			Enforced:           outOfLineConstraintRequest.Enforced,
			NotEnforced:        outOfLineConstraintRequest.NotEnforced,
			Deferrable:         outOfLineConstraintRequest.Deferrable,
			NotDeferrable:      outOfLineConstraintRequest.NotDeferrable,
			InitiallyDeferred:  outOfLineConstraintRequest.InitiallyDeferred,
			InitiallyImmediate: outOfLineConstraintRequest.InitiallyImmediate,
			Enable:             outOfLineConstraintRequest.Enable,
			Disable:            outOfLineConstraintRequest.Disable,
			Validate:           outOfLineConstraintRequest.Validate,
			NoValidate:         outOfLineConstraintRequest.NoValidate,
			Rely:               outOfLineConstraintRequest.Rely,
			NoRely:             outOfLineConstraintRequest.NoRely,
		}
		outOfLineConstraints = append(outOfLineConstraints, outOfLineConstraint)
	}

	opts := &createTableOptions{
		OrReplace:                  s.orReplace,
		IfNotExists:                s.ifNotExists,
		Scope:                      s.scope,
		Kind:                       s.kind,
		name:                       s.name,
		ColumnsAndConstraints:      CreateTableColumnsAndConstraints{convertColumns(s.columns), outOfLineConstraints},
		ClusterBy:                  s.clusterBy,
		EnableSchemaEvolution:      s.enableSchemaEvolution,
		DataRetentionTimeInDays:    s.DataRetentionTimeInDays,
		MaxDataExtensionTimeInDays: s.MaxDataExtensionTimeInDays,
		ChangeTracking:             s.ChangeTracking,
		DefaultDDLCollation:        s.DefaultDDLCollation,
		CopyGrants:                 s.CopyGrants,
		Tags:                       tagAssociations,
		Comment:                    s.Comment,
		RowAccessPolicy:            rowAccessPolicy,
	}

	if s.stageCopyOptions != nil {
		opts.StageCopyOptions = s.stageCopyOptions.toOpts()
	}
	if s.stageFileFormat != nil {
		opts.StageFileFormat = s.stageFileFormat.toOpts()
	}

	return opts
}

func (s *CreateTableAsSelectRequest) toOpts() *createTableAsSelectOptions {
	columns := make([]TableAsSelectColumn, 0, len(s.columns))
	for _, column := range s.columns {
		var maskingPolicy *TableAsSelectColumnMaskingPolicy
		if column.maskingPolicyName != nil {
			maskingPolicy = &TableAsSelectColumnMaskingPolicy{
				Name: *column.maskingPolicyName,
			}
		}
		columns = append(columns, TableAsSelectColumn{
			Name:          column.name,
			Type:          column.type_,
			MaskingPolicy: maskingPolicy,
		})
	}
	return &createTableAsSelectOptions{
		OrReplace: s.orReplace,
		name:      s.name,
		Columns:   columns,
		Query:     s.query,
	}
}

func (s *CreateTableUsingTemplateRequest) toOpts() *createTableUsingTemplateOptions {
	return &createTableUsingTemplateOptions{
		OrReplace:  s.orReplace,
		name:       s.name,
		CopyGrants: s.copyGrants,
		Query:      []string{s.Query},
	}
}

func (s *CreateTableLikeRequest) toOpts() *createTableLikeOptions {
	return &createTableLikeOptions{
		OrReplace:   s.orReplace,
		name:        s.name,
		CopyGrants:  s.copyGrants,
		SourceTable: s.sourceTable,
		ClusterBy:   s.clusterBy,
	}
}

func (s *CreateTableCloneRequest) toOpts() *createTableCloneOptions {
	var clonePoint *ClonePoint
	if s.ClonePoint != nil {
		clonePoint = &ClonePoint{
			Moment: s.ClonePoint.Moment,
			At: TimeTravel{
				Timestamp: s.ClonePoint.At.Timestamp,
				Offset:    s.ClonePoint.At.Offset,
				Statement: s.ClonePoint.At.Statement,
			},
		}
	}
	return &createTableCloneOptions{
		OrReplace:   s.orReplace,
		name:        s.name,
		CopyGrants:  s.copyGrants,
		SourceTable: s.sourceTable,
		ClonePoint:  clonePoint,
	}
}

func (v *StageFileFormatRequest) toOpts() *StageFileFormat {
	return &StageFileFormat{
		FormatName: v.FormatName,
		Type:       v.Type,
		Options:    v.Options.toOpts(),
	}
}

func (v *StageCopyOptionsRequest) toOpts() *StageCopyOptions {
	return &StageCopyOptions{
		OnError:           v.OnError.toOpts(),
		SizeLimit:         v.SizeLimit,
		Purge:             v.Purge,
		ReturnFailedOnly:  v.ReturnFailedOnly,
		MatchByColumnName: v.MatchByColumnName,
		EnforceLength:     v.EnforceLength,
		Truncatecolumns:   v.Truncatecolumns,
		Force:             v.Force,
	}
}

func (v *StageCopyOnErrorOptionsRequest) toOpts() *StageCopyOnErrorOptions {
	return &StageCopyOnErrorOptions{
		Continue:       v.Continue,
		SkipFile:       v.SkipFile,
		AbortStatement: v.AbortStatement,
	}
}

func convertStageFileFormatOptions(stageFileFormatRequests []StageFileFormatRequest) []StageFileFormat {
	fileFormats := make([]StageFileFormat, 0, len(stageFileFormatRequests))
	for _, request := range stageFileFormatRequests {
		var options *FileFormatTypeOptions
		if request.Options != nil {
			options = request.Options.toOpts()
		}
		format := StageFileFormat{
			FormatName: request.FormatName,
			Type:       request.Type,
			Options:    options,
		}
		fileFormats = append(fileFormats, format)
	}
	return fileFormats
}

func (v *FileFormatTypeOptionsRequest) toOpts() *FileFormatTypeOptions {
	if v == nil {
		return nil
	}
	return &FileFormatTypeOptions{
		CSVCompression:                  v.CSVCompression,
		CSVRecordDelimiter:              v.CSVRecordDelimiter,
		CSVFieldDelimiter:               v.CSVFieldDelimiter,
		CSVFileExtension:                v.CSVFileExtension,
		CSVParseHeader:                  v.CSVParseHeader,
		CSVSkipHeader:                   v.CSVSkipHeader,
		CSVSkipBlankLines:               v.CSVSkipBlankLines,
		CSVDateFormat:                   v.CSVDateFormat,
		CSVTimeFormat:                   v.CSVTimeFormat,
		CSVTimestampFormat:              v.CSVTimestampFormat,
		CSVBinaryFormat:                 v.CSVBinaryFormat,
		CSVEscape:                       v.CSVEscape,
		CSVEscapeUnenclosedField:        v.CSVEscapeUnenclosedField,
		CSVTrimSpace:                    v.CSVTrimSpace,
		CSVFieldOptionallyEnclosedBy:    v.CSVFieldOptionallyEnclosedBy,
		CSVNullIf:                       v.CSVNullIf,
		CSVErrorOnColumnCountMismatch:   v.CSVErrorOnColumnCountMismatch,
		CSVReplaceInvalidCharacters:     v.CSVReplaceInvalidCharacters,
		CSVEmptyFieldAsNull:             v.CSVEmptyFieldAsNull,
		CSVSkipByteOrderMark:            v.CSVSkipByteOrderMark,
		CSVEncoding:                     v.CSVEncoding,
		JSONCompression:                 v.JSONCompression,
		JSONDateFormat:                  v.JSONDateFormat,
		JSONTimeFormat:                  v.JSONTimeFormat,
		JSONTimestampFormat:             v.JSONTimestampFormat,
		JSONBinaryFormat:                v.JSONBinaryFormat,
		JSONTrimSpace:                   v.JSONTrimSpace,
		JSONNullIf:                      v.JSONNullIf,
		JSONFileExtension:               v.JSONFileExtension,
		JSONEnableOctal:                 v.JSONEnableOctal,
		JSONAllowDuplicate:              v.JSONAllowDuplicate,
		JSONStripOuterArray:             v.JSONStripOuterArray,
		JSONStripNullValues:             v.JSONStripNullValues,
		JSONReplaceInvalidCharacters:    v.JSONReplaceInvalidCharacters,
		JSONIgnoreUTF8Errors:            v.JSONIgnoreUTF8Errors,
		JSONSkipByteOrderMark:           v.JSONSkipByteOrderMark,
		AvroCompression:                 v.AvroCompression,
		AvroTrimSpace:                   v.AvroTrimSpace,
		AvroReplaceInvalidCharacters:    v.AvroReplaceInvalidCharacters,
		AvroNullIf:                      v.AvroNullIf,
		ORCTrimSpace:                    v.ORCTrimSpace,
		ORCReplaceInvalidCharacters:     v.ORCReplaceInvalidCharacters,
		ORCNullIf:                       v.ORCNullIf,
		ParquetCompression:              v.ParquetCompression,
		ParquetSnappyCompression:        v.ParquetSnappyCompression,
		ParquetBinaryAsText:             v.ParquetBinaryAsText,
		ParquetTrimSpace:                v.ParquetTrimSpace,
		ParquetReplaceInvalidCharacters: v.ParquetReplaceInvalidCharacters,
		ParquetNullIf:                   v.ParquetNullIf,
		XMLCompression:                  v.XMLCompression,
		XMLIgnoreUTF8Errors:             v.XMLIgnoreUTF8Errors,
		XMLPreserveSpace:                v.XMLPreserveSpace,
		XMLStripOuterElement:            v.XMLStripOuterElement,
		XMLDisableSnowflakeData:         v.XMLDisableSnowflakeData,
		XMLDisableAutoConvert:           v.XMLDisableAutoConvert,
		XMLReplaceInvalidCharacters:     v.XMLReplaceInvalidCharacters,
		XMLSkipByteOrderMark:            v.XMLSkipByteOrderMark,
	}
}

func convertColumns(columnRequests []TableColumnRequest) []TableColumn {
	columns := make([]TableColumn, 0, len(columnRequests))
	for _, columnRequest := range columnRequests {
		columnRequest := columnRequest
		var defaultValue *ColumnDefaultValue
		if columnRequest.defaultValue != nil {
			var columnIdentity *ColumnIdentity
			if columnRequest.defaultValue.identity != nil {
				columnIdentity = &ColumnIdentity{
					Start:     columnRequest.defaultValue.identity.Start,
					Increment: columnRequest.defaultValue.identity.Increment,
					Order:     columnRequest.defaultValue.identity.Order,
					Noorder:   columnRequest.defaultValue.identity.Noorder,
				}
			}
			defaultValue = &ColumnDefaultValue{
				columnRequest.defaultValue.expression,
				columnIdentity,
			}
		}
		var inlineConstraint *ColumnInlineConstraint
		if columnRequest.inlineConstraint != nil {
			var foreignKey *InlineForeignKey
			if columnRequest.inlineConstraint.foreignKey != nil {
				var onActionRequest *ForeignKeyOnAction
				if columnRequest.inlineConstraint.foreignKey.On != nil {
					onActionRequest = &ForeignKeyOnAction{
						OnUpdate: columnRequest.inlineConstraint.foreignKey.On.OnUpdate,
						OnDelete: columnRequest.inlineConstraint.foreignKey.On.OnDelete,
					}
				}
				foreignKey = &InlineForeignKey{
					TableName:  columnRequest.inlineConstraint.foreignKey.TableName,
					ColumnName: columnRequest.inlineConstraint.foreignKey.ColumnName,
					Match:      columnRequest.inlineConstraint.foreignKey.Match,
					On:         onActionRequest,
				}
			}
			inlineConstraint = &ColumnInlineConstraint{
				Name:               &columnRequest.inlineConstraint.Name,
				Type:               columnRequest.inlineConstraint.type_,
				ForeignKey:         foreignKey,
				Enforced:           columnRequest.inlineConstraint.enforced,
				NotEnforced:        columnRequest.inlineConstraint.notEnforced,
				Deferrable:         columnRequest.inlineConstraint.deferrable,
				NotDeferrable:      columnRequest.inlineConstraint.notDeferrable,
				InitiallyDeferred:  columnRequest.inlineConstraint.initiallyDeferred,
				InitiallyImmediate: columnRequest.inlineConstraint.initiallyImmediate,
				Enable:             columnRequest.inlineConstraint.enable,
				Disable:            columnRequest.inlineConstraint.disable,
				Validate:           columnRequest.inlineConstraint.validate,
				NoValidate:         columnRequest.inlineConstraint.noValidate,
				Rely:               columnRequest.inlineConstraint.rely,
				NoRely:             columnRequest.inlineConstraint.noRely,
			}
		}
		var maskingPolicy *ColumnMaskingPolicy
		if columnRequest.maskingPolicy != nil {
			maskingPolicy = &ColumnMaskingPolicy{
				With:  columnRequest.maskingPolicy.with,
				Name:  columnRequest.maskingPolicy.name,
				Using: columnRequest.maskingPolicy.using,
			}
		}
		columns = append(columns, TableColumn{
			Name:             columnRequest.name,
			Type:             columnRequest.type_,
			Collate:          columnRequest.collate,
			Comment:          columnRequest.comment,
			DefaultValue:     defaultValue,
			MaskingPolicy:    maskingPolicy,
			NotNull:          columnRequest.notNull,
			Tags:             columnRequest.tags,
			InlineConstraint: inlineConstraint,
		})
	}
	return columns
}

func (v *DescribeTableColumnsRequest) toOpts() *describeTableColumnsOptions {
	return &describeTableColumnsOptions{
		name: v.id,
	}
}

func (v *DescribeTableStageRequest) toOpts() *describeTableStageOptions {
	return &describeTableStageOptions{
		name: v.id,
	}
}
