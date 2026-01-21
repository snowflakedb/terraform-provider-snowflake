package sdk

import (
	"fmt"
	"time"
)

func NewCreateTableAsSelectRequest(
	name SchemaObjectIdentifier,
	columns []TableAsSelectColumnRequest,
	query string,
) *CreateTableAsSelectRequest {
	s := CreateTableAsSelectRequest{}
	s.name = name
	s.columns = columns
	s.query = query
	return &s
}

func (s *CreateTableAsSelectRequest) WithOrReplace(orReplace *bool) *CreateTableAsSelectRequest {
	s.orReplace = orReplace
	return s
}

func NewTableAsSelectColumnRequest(
	name string,
) *TableAsSelectColumnRequest {
	s := TableAsSelectColumnRequest{}
	s.name = name
	return &s
}

func (s *TableAsSelectColumnRequest) WithOrReplace(orReplace *bool) *TableAsSelectColumnRequest {
	s.orReplace = orReplace
	return s
}

func (s *TableAsSelectColumnRequest) WithType_(type_ *DataType) *TableAsSelectColumnRequest {
	s.type_ = type_
	return s
}

func (s *TableAsSelectColumnRequest) WithMaskingPolicyName(maskingPolicyName *SchemaObjectIdentifier) *TableAsSelectColumnRequest {
	s.maskingPolicyName = maskingPolicyName
	return s
}

func (s *TableAsSelectColumnRequest) WithClusterBy(clusterBy []string) *TableAsSelectColumnRequest {
	s.clusterBy = clusterBy
	return s
}

func (s *TableAsSelectColumnRequest) WithCopyGrants(copyGrants *bool) *TableAsSelectColumnRequest {
	s.copyGrants = copyGrants
	return s
}

func NewCreateTableUsingTemplateRequest(
	name SchemaObjectIdentifier,
	query string,
) *CreateTableUsingTemplateRequest {
	s := CreateTableUsingTemplateRequest{}
	s.name = name
	s.Query = query
	return &s
}

func (s *CreateTableUsingTemplateRequest) WithOrReplace(orReplace *bool) *CreateTableUsingTemplateRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateTableUsingTemplateRequest) WithCopyGrants(copyGrants *bool) *CreateTableUsingTemplateRequest {
	s.copyGrants = copyGrants
	return s
}

func NewCreateTableLikeRequest(
	name SchemaObjectIdentifier,
	sourceTable SchemaObjectIdentifier,
) *CreateTableLikeRequest {
	s := CreateTableLikeRequest{}
	s.name = name
	s.sourceTable = sourceTable
	return &s
}

func (s *CreateTableLikeRequest) WithOrReplace(orReplace *bool) *CreateTableLikeRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateTableLikeRequest) WithClusterBy(clusterBy []string) *CreateTableLikeRequest {
	s.clusterBy = clusterBy
	return s
}

func (s *CreateTableLikeRequest) WithCopyGrants(copyGrants *bool) *CreateTableLikeRequest {
	s.copyGrants = copyGrants
	return s
}

func NewCreateTableCloneRequest(
	name SchemaObjectIdentifier,
	sourceTable SchemaObjectIdentifier,
) *CreateTableCloneRequest {
	s := CreateTableCloneRequest{}
	s.name = name
	s.sourceTable = sourceTable
	return &s
}

func (s *CreateTableCloneRequest) WithOrReplace(orReplace *bool) *CreateTableCloneRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateTableCloneRequest) WithCopyGrants(copyGrants *bool) *CreateTableCloneRequest {
	s.copyGrants = copyGrants
	return s
}

func (s *CreateTableCloneRequest) WithClonePoint(clonePoint *ClonePointRequest) *CreateTableCloneRequest {
	s.ClonePoint = clonePoint
	return s
}

func NewClonePointRequest() *ClonePointRequest {
	return &ClonePointRequest{}
}

func (s *ClonePointRequest) WithMoment(moment CloneMoment) *ClonePointRequest {
	s.Moment = moment
	return s
}

func (s *ClonePointRequest) WithAt(at TimeTravelRequest) *ClonePointRequest {
	s.At = at
	return s
}

func NewTimeTravelRequest() *TimeTravelRequest {
	return &TimeTravelRequest{}
}

func (s *TimeTravelRequest) WithTimestamp(timestamp *time.Time) *TimeTravelRequest {
	s.Timestamp = timestamp
	return s
}

func (s *TimeTravelRequest) WithOffset(offset *int) *TimeTravelRequest {
	s.Offset = offset
	return s
}

func (s *TimeTravelRequest) WithStatement(statement *string) *TimeTravelRequest {
	s.Statement = statement
	return s
}

func NewCreateTableRequest(
	name SchemaObjectIdentifier,
	columns []TableColumnRequest,
) *CreateTableRequest {
	s := CreateTableRequest{}
	s.name = name
	s.columns = columns
	return &s
}

func (s *CreateTableRequest) WithOrReplace(orReplace *bool) *CreateTableRequest {
	s.orReplace = orReplace
	return s
}

func (s *CreateTableRequest) WithIfNotExists(ifNotExists *bool) *CreateTableRequest {
	s.ifNotExists = ifNotExists
	return s
}

func (s *CreateTableRequest) WithScope(scope *TableScope) *CreateTableRequest {
	s.scope = scope
	return s
}

func (s *CreateTableRequest) WithKind(kind *TableKind) *CreateTableRequest {
	s.kind = kind
	return s
}

func (s *CreateTableRequest) WithOutOfLineConstraint(outOfLineConstraint OutOfLineConstraintRequest) *CreateTableRequest {
	s.OutOfLineConstraints = append(s.OutOfLineConstraints, outOfLineConstraint)
	return s
}

func (s *CreateTableRequest) WithClusterBy(clusterBy []string) *CreateTableRequest {
	s.clusterBy = clusterBy
	return s
}

func (s *CreateTableRequest) WithEnableSchemaEvolution(enableSchemaEvolution *bool) *CreateTableRequest {
	s.enableSchemaEvolution = enableSchemaEvolution
	return s
}

func (s *CreateTableRequest) WithStageFileFormat(stageFileFormat LegacyFileFormatRequest) *CreateTableRequest {
	s.stageFileFormat = &stageFileFormat
	return s
}

func (s *CreateTableRequest) WithStageCopyOptions(stageCopyOptions LegacyTableCopyOptionsRequest) *CreateTableRequest {
	s.stageCopyOptions = &stageCopyOptions
	return s
}

func (s *CreateTableRequest) WithDataRetentionTimeInDays(dataRetentionTimeInDays *int) *CreateTableRequest {
	s.DataRetentionTimeInDays = dataRetentionTimeInDays
	return s
}

func (s *CreateTableRequest) WithMaxDataExtensionTimeInDays(maxDataExtensionTimeInDays *int) *CreateTableRequest {
	s.MaxDataExtensionTimeInDays = maxDataExtensionTimeInDays
	return s
}

func (s *CreateTableRequest) WithChangeTracking(changeTracking *bool) *CreateTableRequest {
	s.ChangeTracking = changeTracking
	return s
}

func (s *CreateTableRequest) WithDefaultDDLCollation(defaultDDLCollation *string) *CreateTableRequest {
	s.DefaultDDLCollation = defaultDDLCollation
	return s
}

func (s *CreateTableRequest) WithCopyGrants(copyGrants *bool) *CreateTableRequest {
	s.CopyGrants = copyGrants
	return s
}

func (s *CreateTableRequest) WithRowAccessPolicy(rowAccessPolicy *RowAccessPolicyRequest) *CreateTableRequest {
	s.RowAccessPolicy = rowAccessPolicy
	return s
}

func (s *CreateTableRequest) WithTags(tags []TagAssociationRequest) *CreateTableRequest {
	s.Tags = tags
	return s
}

func (s *CreateTableRequest) WithComment(comment *string) *CreateTableRequest {
	s.Comment = comment
	return s
}

func NewTableColumnRequest(
	name string,
	type_ DataType,
) *TableColumnRequest {
	s := TableColumnRequest{}
	s.name = name
	s.type_ = type_
	return &s
}

func (s *TableColumnRequest) WithCollate(collate *string) *TableColumnRequest {
	s.collate = collate
	return s
}

func (s *TableColumnRequest) WithComment(comment *string) *TableColumnRequest {
	s.comment = comment
	return s
}

func (s *TableColumnRequest) WithDefaultValue(defaultValue *ColumnDefaultValueRequest) *TableColumnRequest {
	s.defaultValue = defaultValue
	return s
}

func (s *TableColumnRequest) WithNotNull(notNull *bool) *TableColumnRequest {
	s.notNull = notNull
	return s
}

func (s *TableColumnRequest) WithMaskingPolicy(maskingPolicy *ColumnMaskingPolicyRequest) *TableColumnRequest {
	s.maskingPolicy = maskingPolicy
	return s
}

func (s *TableColumnRequest) WithTags(tags []TagAssociation) *TableColumnRequest {
	s.tags = tags
	return s
}

func (s *TableColumnRequest) WithInlineConstraint(inlineConstraint *ColumnInlineConstraintRequest) *TableColumnRequest {
	s.inlineConstraint = inlineConstraint
	return s
}

func NewColumnDefaultValueRequest() *ColumnDefaultValueRequest {
	return &ColumnDefaultValueRequest{}
}

func (s *ColumnDefaultValueRequest) WithExpression(expression *string) *ColumnDefaultValueRequest {
	s.expression = expression
	return s
}

func (s *ColumnDefaultValueRequest) WithIdentity(identity *ColumnIdentityRequest) *ColumnDefaultValueRequest {
	s.identity = identity
	return s
}

func NewColumnIdentityRequest(
	start int,
	increment int,
) *ColumnIdentityRequest {
	s := ColumnIdentityRequest{}
	s.Start = start
	s.Increment = increment
	return &s
}

func (s *ColumnIdentityRequest) WithOrder() *ColumnIdentityRequest {
	s.Order = Bool(true)
	return s
}

func (s *ColumnIdentityRequest) WithNoorder() *ColumnIdentityRequest {
	s.Noorder = Bool(true)
	return s
}

func NewColumnMaskingPolicyRequest(
	name SchemaObjectIdentifier,
) *ColumnMaskingPolicyRequest {
	s := ColumnMaskingPolicyRequest{}
	s.name = name
	return &s
}

func (s *ColumnMaskingPolicyRequest) WithWith(with *bool) *ColumnMaskingPolicyRequest {
	s.with = with
	return s
}

func (s *ColumnMaskingPolicyRequest) WithUsing(using []string) *ColumnMaskingPolicyRequest {
	s.using = using
	return s
}

func NewColumnInlineConstraintRequest(
	name string,
	type_ ColumnConstraintType,
) *ColumnInlineConstraintRequest {
	s := ColumnInlineConstraintRequest{}
	s.Name = name
	s.type_ = type_
	return &s
}

func (s *ColumnInlineConstraintRequest) WithForeignKey(foreignKey *InlineForeignKeyRequest) *ColumnInlineConstraintRequest {
	s.foreignKey = foreignKey
	return s
}

func (s *ColumnInlineConstraintRequest) WithEnforced(enforced *bool) *ColumnInlineConstraintRequest {
	s.enforced = enforced
	return s
}

func (s *ColumnInlineConstraintRequest) WithNotEnforced(notEnforced *bool) *ColumnInlineConstraintRequest {
	s.notEnforced = notEnforced
	return s
}

func (s *ColumnInlineConstraintRequest) WithDeferrable(deferrable *bool) *ColumnInlineConstraintRequest {
	s.deferrable = deferrable
	return s
}

func (s *ColumnInlineConstraintRequest) WithNotDeferrable(notDeferrable *bool) *ColumnInlineConstraintRequest {
	s.notDeferrable = notDeferrable
	return s
}

func (s *ColumnInlineConstraintRequest) WithInitiallyDeferred(initiallyDeferred *bool) *ColumnInlineConstraintRequest {
	s.initiallyDeferred = initiallyDeferred
	return s
}

func (s *ColumnInlineConstraintRequest) WithInitiallyImmediate(initiallyImmediate *bool) *ColumnInlineConstraintRequest {
	s.initiallyImmediate = initiallyImmediate
	return s
}

func (s *ColumnInlineConstraintRequest) WithEnable(enable *bool) *ColumnInlineConstraintRequest {
	s.enable = enable
	return s
}

func (s *ColumnInlineConstraintRequest) WithDisable(disable *bool) *ColumnInlineConstraintRequest {
	s.disable = disable
	return s
}

func (s *ColumnInlineConstraintRequest) WithValidate(validate *bool) *ColumnInlineConstraintRequest {
	s.validate = validate
	return s
}

func (s *ColumnInlineConstraintRequest) WithNoValidate(noValidate *bool) *ColumnInlineConstraintRequest {
	s.noValidate = noValidate
	return s
}

func (s *ColumnInlineConstraintRequest) WithRely(rely *bool) *ColumnInlineConstraintRequest {
	s.rely = rely
	return s
}

func (s *ColumnInlineConstraintRequest) WithNoRely(noRely *bool) *ColumnInlineConstraintRequest {
	s.noRely = noRely
	return s
}

func NewOutOfLineConstraintRequest(
	constraintType ColumnConstraintType,
) *OutOfLineConstraintRequest {
	s := OutOfLineConstraintRequest{}
	s.Type = constraintType
	return &s
}

func (s *OutOfLineConstraintRequest) WithName(name *string) *OutOfLineConstraintRequest {
	s.Name = name
	return s
}

func (s *OutOfLineConstraintRequest) WithColumns(columns []string) *OutOfLineConstraintRequest {
	s.Columns = columns
	return s
}

func (s *OutOfLineConstraintRequest) WithForeignKey(foreignKey *OutOfLineForeignKeyRequest) *OutOfLineConstraintRequest {
	s.ForeignKey = foreignKey
	return s
}

func (s *OutOfLineConstraintRequest) WithEnforced(enforced *bool) *OutOfLineConstraintRequest {
	s.Enforced = enforced
	return s
}

func (s *OutOfLineConstraintRequest) WithNotEnforced(notEnforced *bool) *OutOfLineConstraintRequest {
	s.NotEnforced = notEnforced
	return s
}

func (s *OutOfLineConstraintRequest) WithDeferrable(deferrable *bool) *OutOfLineConstraintRequest {
	s.Deferrable = deferrable
	return s
}

func (s *OutOfLineConstraintRequest) WithNotDeferrable(notDeferrable *bool) *OutOfLineConstraintRequest {
	s.NotDeferrable = notDeferrable
	return s
}

func (s *OutOfLineConstraintRequest) WithInitiallyDeferred(initiallyDeferred *bool) *OutOfLineConstraintRequest {
	s.InitiallyDeferred = initiallyDeferred
	return s
}

func (s *OutOfLineConstraintRequest) WithInitiallyImmediate(initiallyImmediate *bool) *OutOfLineConstraintRequest {
	s.InitiallyImmediate = initiallyImmediate
	return s
}

func (s *OutOfLineConstraintRequest) WithEnable(enable *bool) *OutOfLineConstraintRequest {
	s.Enable = enable
	return s
}

func (s *OutOfLineConstraintRequest) WithDisable(disable *bool) *OutOfLineConstraintRequest {
	s.Disable = disable
	return s
}

func (s *OutOfLineConstraintRequest) WithValidate(validate *bool) *OutOfLineConstraintRequest {
	s.Validate = validate
	return s
}

func (s *OutOfLineConstraintRequest) WithNoValidate(noValidate *bool) *OutOfLineConstraintRequest {
	s.NoValidate = noValidate
	return s
}

func (s *OutOfLineConstraintRequest) WithRely(rely *bool) *OutOfLineConstraintRequest {
	s.Rely = rely
	return s
}

func (s *OutOfLineConstraintRequest) WithNoRely(noRely *bool) *OutOfLineConstraintRequest {
	s.NoRely = noRely
	return s
}

func NewInlineForeignKeyRequest(
	tableName string,
) *InlineForeignKeyRequest {
	s := InlineForeignKeyRequest{}
	s.TableName = tableName
	return &s
}

func (s *InlineForeignKeyRequest) WithColumnName(columnName []string) *InlineForeignKeyRequest {
	s.ColumnName = columnName
	return s
}

func (s *InlineForeignKeyRequest) WithMatch(match *MatchType) *InlineForeignKeyRequest {
	s.Match = match
	return s
}

func (s *InlineForeignKeyRequest) WithOn(on *ForeignKeyOnAction) *InlineForeignKeyRequest {
	s.On = on
	return s
}

func NewOutOfLineForeignKeyRequest(
	tableName SchemaObjectIdentifier,
	columnNames []string,
) *OutOfLineForeignKeyRequest {
	s := OutOfLineForeignKeyRequest{}
	s.TableName = tableName
	s.ColumnNames = columnNames
	return &s
}

func (s *OutOfLineForeignKeyRequest) WithMatch(match *MatchType) *OutOfLineForeignKeyRequest {
	s.Match = match
	return s
}

func (s *OutOfLineForeignKeyRequest) WithOn(on *ForeignKeyOnAction) *OutOfLineForeignKeyRequest {
	s.On = on
	return s
}

func NewForeignKeyOnAction() *ForeignKeyOnAction {
	return &ForeignKeyOnAction{}
}

func (s *ForeignKeyOnAction) WithOnUpdate(onUpdate *ForeignKeyAction) *ForeignKeyOnAction {
	s.OnUpdate = onUpdate
	return s
}

func (s *ForeignKeyOnAction) WithOnDelete(onDelete *ForeignKeyAction) *ForeignKeyOnAction {
	s.OnDelete = onDelete
	return s
}

func NewAlterTableRequest(
	name SchemaObjectIdentifier,
) *AlterTableRequest {
	s := AlterTableRequest{}
	s.name = name
	return &s
}

func (s *AlterTableRequest) WithIfExists(ifExists *bool) *AlterTableRequest {
	s.IfExists = ifExists
	return s
}

func (s *AlterTableRequest) WithNewName(newName *SchemaObjectIdentifier) *AlterTableRequest {
	s.NewName = newName
	return s
}

func (s *AlterTableRequest) WithSwapWith(swapWith *SchemaObjectIdentifier) *AlterTableRequest {
	s.SwapWith = swapWith
	return s
}

func (s *AlterTableRequest) WithClusteringAction(clusteringAction *TableClusteringActionRequest) *AlterTableRequest {
	s.ClusteringAction = clusteringAction
	return s
}

func (s *AlterTableRequest) WithColumnAction(columnAction *TableColumnActionRequest) *AlterTableRequest {
	s.ColumnAction = columnAction
	return s
}

func (s *AlterTableRequest) WithConstraintAction(constraintAction *TableConstraintActionRequest) *AlterTableRequest {
	s.ConstraintAction = constraintAction
	return s
}

func (s *AlterTableRequest) WithExternalTableAction(externalTableAction *TableExternalTableActionRequest) *AlterTableRequest {
	s.ExternalTableAction = externalTableAction
	return s
}

func (s *AlterTableRequest) WithSearchOptimizationAction(searchOptimizationAction *TableSearchOptimizationActionRequest) *AlterTableRequest {
	s.SearchOptimizationAction = searchOptimizationAction
	return s
}

func (s *AlterTableRequest) WithSet(set *TableSetRequest) *AlterTableRequest {
	s.Set = set
	return s
}

func (s *AlterTableRequest) WithSetTags(setTags []TagAssociationRequest) *AlterTableRequest {
	s.SetTags = setTags
	return s
}

func (s *AlterTableRequest) WithUnsetTags(unsetTags []ObjectIdentifier) *AlterTableRequest {
	s.UnsetTags = unsetTags
	return s
}

func (s *AlterTableRequest) WithUnset(unset *TableUnsetRequest) *AlterTableRequest {
	s.Unset = unset
	return s
}

func (s *AlterTableRequest) WithAddRowAccessPolicy(addRowAccessPolicy *TableAddRowAccessPolicyRequest) *AlterTableRequest {
	s.AddRowAccessPolicy = addRowAccessPolicy
	return s
}

func (s *AlterTableRequest) WithDropRowAccessPolicy(dropRowAccessPolicy *TableDropRowAccessPolicyRequest) *AlterTableRequest {
	s.DropRowAccessPolicy = dropRowAccessPolicy
	return s
}

func (s *AlterTableRequest) WithDropAndAddRowAccessPolicy(dropAndAddRowAccessPolicy *TableDropAndAddRowAccessPolicy) *AlterTableRequest {
	s.DropAndAddRowAccessPolicy = dropAndAddRowAccessPolicy
	return s
}

func (s *AlterTableRequest) WithDropAllAccessRowPolicies(dropAllAccessRowPolicies *bool) *AlterTableRequest {
	s.DropAllAccessRowPolicies = dropAllAccessRowPolicies
	return s
}

func NewDropTableRequest(
	name SchemaObjectIdentifier,
) *DropTableRequest {
	s := DropTableRequest{}
	s.Name = name
	return &s
}

func (s *DropTableRequest) WithIfExists(ifExists *bool) *DropTableRequest {
	s.IfExists = ifExists
	return s
}

func (s *DropTableRequest) WithCascade(cascade *bool) *DropTableRequest {
	s.Cascade = cascade
	return s
}

func (s *DropTableRequest) WithRestrict(restrict *bool) *DropTableRequest {
	s.Restrict = restrict
	return s
}

func NewTableAddRowAccessPolicyRequest(
	rowAccessPolicy SchemaObjectIdentifier,
	on []string,
) *TableAddRowAccessPolicyRequest {
	s := TableAddRowAccessPolicyRequest{}
	s.RowAccessPolicy = rowAccessPolicy
	s.On = on
	return &s
}

func NewTableDropRowAccessPolicyRequest(
	rowAccessPolicy SchemaObjectIdentifier,
) *TableDropRowAccessPolicyRequest {
	s := TableDropRowAccessPolicyRequest{}
	s.RowAccessPolicy = rowAccessPolicy
	return &s
}

func NewTableDropAndAddRowAccessPolicyRequest(
	drop TableDropRowAccessPolicyRequest,
	add TableAddRowAccessPolicyRequest,
) *TableDropAndAddRowAccessPolicyRequest {
	s := TableDropAndAddRowAccessPolicyRequest{}
	s.Drop = drop
	s.Add = add
	return &s
}

func NewTableUnsetRequest() *TableUnsetRequest {
	return &TableUnsetRequest{}
}

func (s *TableUnsetRequest) WithDataRetentionTimeInDays(dataRetentionTimeInDays bool) *TableUnsetRequest {
	s.DataRetentionTimeInDays = dataRetentionTimeInDays
	return s
}

func (s *TableUnsetRequest) WithMaxDataExtensionTimeInDays(maxDataExtensionTimeInDays bool) *TableUnsetRequest {
	s.MaxDataExtensionTimeInDays = maxDataExtensionTimeInDays
	return s
}

func (s *TableUnsetRequest) WithChangeTracking(changeTracking bool) *TableUnsetRequest {
	s.ChangeTracking = changeTracking
	return s
}

func (s *TableUnsetRequest) WithDefaultDDLCollation(defaultDDLCollation bool) *TableUnsetRequest {
	s.DefaultDDLCollation = defaultDDLCollation
	return s
}

func (s *TableUnsetRequest) WithEnableSchemaEvolution(enableSchemaEvolution bool) *TableUnsetRequest {
	s.EnableSchemaEvolution = enableSchemaEvolution
	return s
}

func (s *TableUnsetRequest) WithComment(comment bool) *TableUnsetRequest {
	s.Comment = comment
	return s
}

func NewAddRowAccessPolicyRequest(
	policyName string,
	columnName []string,
) *AddRowAccessPolicyRequest {
	s := AddRowAccessPolicyRequest{}
	s.PolicyName = policyName
	s.ColumnName = columnName
	return &s
}

func NewTagAssociationRequest(
	name ObjectIdentifier,
	value string,
) *TagAssociationRequest {
	s := TagAssociationRequest{}
	s.Name = name
	s.Value = value
	return &s
}

func NewFileFormatTypeOptionsRequest() *LegacyFileFormatTypeOptionsRequest {
	return &LegacyFileFormatTypeOptionsRequest{}
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVCompression(csvCompression *CSVCompression) *LegacyFileFormatTypeOptionsRequest {
	s.CSVCompression = csvCompression
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVRecordDelimiter(csvRecordDelimiter *string) *LegacyFileFormatTypeOptionsRequest {
	s.CSVRecordDelimiter = csvRecordDelimiter
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVFieldDelimiter(csvFieldDelimiter *string) *LegacyFileFormatTypeOptionsRequest {
	s.CSVFieldDelimiter = csvFieldDelimiter
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVFileExtension(csvFileExtension *string) *LegacyFileFormatTypeOptionsRequest {
	s.CSVFileExtension = csvFileExtension
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVParseHeader(csvParseHeader *bool) *LegacyFileFormatTypeOptionsRequest {
	s.CSVParseHeader = csvParseHeader
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVSkipHeader(csvSkipHeader *int) *LegacyFileFormatTypeOptionsRequest {
	s.CSVSkipHeader = csvSkipHeader
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVSkipBlankLines(csvSkipBlankLines *bool) *LegacyFileFormatTypeOptionsRequest {
	s.CSVSkipBlankLines = csvSkipBlankLines
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVDateFormat(csvDateFormat *string) *LegacyFileFormatTypeOptionsRequest {
	s.CSVDateFormat = csvDateFormat
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVTimeFormat(csvTimeFormat *string) *LegacyFileFormatTypeOptionsRequest {
	s.CSVTimeFormat = csvTimeFormat
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVTimestampFormat(csvTimestampFormat *string) *LegacyFileFormatTypeOptionsRequest {
	s.CSVTimestampFormat = csvTimestampFormat
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVBinaryFormat(csvBinaryFormat *BinaryFormat) *LegacyFileFormatTypeOptionsRequest {
	s.CSVBinaryFormat = csvBinaryFormat
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVEscape(csvEscape *string) *LegacyFileFormatTypeOptionsRequest {
	s.CSVEscape = csvEscape
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVEscapeUnenclosedField(csvEscapeUnenclosedField *string) *LegacyFileFormatTypeOptionsRequest {
	s.CSVEscapeUnenclosedField = csvEscapeUnenclosedField
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVTrimSpace(csvTrimSpace *bool) *LegacyFileFormatTypeOptionsRequest {
	s.CSVTrimSpace = csvTrimSpace
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVFieldOptionallyEnclosedBy(csvFieldOptionallyEnclosedBy *string) *LegacyFileFormatTypeOptionsRequest {
	s.CSVFieldOptionallyEnclosedBy = csvFieldOptionallyEnclosedBy
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVNullIf(csvNullIf *[]NullString) *LegacyFileFormatTypeOptionsRequest {
	s.CSVNullIf = csvNullIf
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVErrorOnColumnCountMismatch(csvErrorOnColumnCountMismatch *bool) *LegacyFileFormatTypeOptionsRequest {
	s.CSVErrorOnColumnCountMismatch = csvErrorOnColumnCountMismatch
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVReplaceInvalidCharacters(csvReplaceInvalidCharacters *bool) *LegacyFileFormatTypeOptionsRequest {
	s.CSVReplaceInvalidCharacters = csvReplaceInvalidCharacters
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVEmptyFieldAsNull(csvEmptyFieldAsNull *bool) *LegacyFileFormatTypeOptionsRequest {
	s.CSVEmptyFieldAsNull = csvEmptyFieldAsNull
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVSkipByteOrderMark(csvSkipByteOrderMark *bool) *LegacyFileFormatTypeOptionsRequest {
	s.CSVSkipByteOrderMark = csvSkipByteOrderMark
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithCSVEncoding(csvEncoding *CSVEncoding) *LegacyFileFormatTypeOptionsRequest {
	s.CSVEncoding = csvEncoding
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONCompression(jsonCompression *JSONCompression) *LegacyFileFormatTypeOptionsRequest {
	s.JSONCompression = jsonCompression
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONDateFormat(jsonDateFormat *string) *LegacyFileFormatTypeOptionsRequest {
	s.JSONDateFormat = jsonDateFormat
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONTimeFormat(jsonTimeFormat *string) *LegacyFileFormatTypeOptionsRequest {
	s.JSONTimeFormat = jsonTimeFormat
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONTimestampFormat(jsonTimestampFormat *string) *LegacyFileFormatTypeOptionsRequest {
	s.JSONTimestampFormat = jsonTimestampFormat
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONBinaryFormat(jsonBinaryFormat *BinaryFormat) *LegacyFileFormatTypeOptionsRequest {
	s.JSONBinaryFormat = jsonBinaryFormat
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONTrimSpace(jsonTrimSpace *bool) *LegacyFileFormatTypeOptionsRequest {
	s.JSONTrimSpace = jsonTrimSpace
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONNullIf(jsonNullIf []NullString) *LegacyFileFormatTypeOptionsRequest {
	s.JSONNullIf = jsonNullIf
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONFileExtension(jsonFileExtension *string) *LegacyFileFormatTypeOptionsRequest {
	s.JSONFileExtension = jsonFileExtension
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONEnableOctal(jsonEnableOctal *bool) *LegacyFileFormatTypeOptionsRequest {
	s.JSONEnableOctal = jsonEnableOctal
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONAllowDuplicate(jsonAllowDuplicate *bool) *LegacyFileFormatTypeOptionsRequest {
	s.JSONAllowDuplicate = jsonAllowDuplicate
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONStripOuterArray(jsonStripOuterArray *bool) *LegacyFileFormatTypeOptionsRequest {
	s.JSONStripOuterArray = jsonStripOuterArray
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONStripNullValues(jsonStripNullValues *bool) *LegacyFileFormatTypeOptionsRequest {
	s.JSONStripNullValues = jsonStripNullValues
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONReplaceInvalidCharacters(jsonReplaceInvalidCharacters *bool) *LegacyFileFormatTypeOptionsRequest {
	s.JSONReplaceInvalidCharacters = jsonReplaceInvalidCharacters
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONIgnoreUTF8Errors(jsonIgnoreUTF8Errors *bool) *LegacyFileFormatTypeOptionsRequest {
	s.JSONIgnoreUTF8Errors = jsonIgnoreUTF8Errors
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithJSONSkipByteOrderMark(jsonSkipByteOrderMark *bool) *LegacyFileFormatTypeOptionsRequest {
	s.JSONSkipByteOrderMark = jsonSkipByteOrderMark
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithAvroCompression(avroCompression *AvroCompression) *LegacyFileFormatTypeOptionsRequest {
	s.AvroCompression = avroCompression
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithAvroTrimSpace(avroTrimSpace *bool) *LegacyFileFormatTypeOptionsRequest {
	s.AvroTrimSpace = avroTrimSpace
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithAvroReplaceInvalidCharacters(avroReplaceInvalidCharacters *bool) *LegacyFileFormatTypeOptionsRequest {
	s.AvroReplaceInvalidCharacters = avroReplaceInvalidCharacters
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithAvroNullIf(avroNullIf *[]NullString) *LegacyFileFormatTypeOptionsRequest {
	s.AvroNullIf = avroNullIf
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithORCTrimSpace(orcTrimSpace *bool) *LegacyFileFormatTypeOptionsRequest {
	s.ORCTrimSpace = orcTrimSpace
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithORCReplaceInvalidCharacters(orcReplaceInvalidCharacters *bool) *LegacyFileFormatTypeOptionsRequest {
	s.ORCReplaceInvalidCharacters = orcReplaceInvalidCharacters
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithORCNullIf(orcNullIf *[]NullString) *LegacyFileFormatTypeOptionsRequest {
	s.ORCNullIf = orcNullIf
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithParquetCompression(parquetCompression *ParquetCompression) *LegacyFileFormatTypeOptionsRequest {
	s.ParquetCompression = parquetCompression
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithParquetSnappyCompression(parquetSnappyCompression *bool) *LegacyFileFormatTypeOptionsRequest {
	s.ParquetSnappyCompression = parquetSnappyCompression
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithParquetBinaryAsText(parquetBinaryAsText *bool) *LegacyFileFormatTypeOptionsRequest {
	s.ParquetBinaryAsText = parquetBinaryAsText
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithParquetTrimSpace(parquetTrimSpace *bool) *LegacyFileFormatTypeOptionsRequest {
	s.ParquetTrimSpace = parquetTrimSpace
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithParquetReplaceInvalidCharacters(parquetReplaceInvalidCharacters *bool) *LegacyFileFormatTypeOptionsRequest {
	s.ParquetReplaceInvalidCharacters = parquetReplaceInvalidCharacters
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithParquetNullIf(parquetNullIf *[]NullString) *LegacyFileFormatTypeOptionsRequest {
	s.ParquetNullIf = parquetNullIf
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithXMLCompression(xmlCompression *XMLCompression) *LegacyFileFormatTypeOptionsRequest {
	s.XMLCompression = xmlCompression
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithXMLIgnoreUTF8Errors(xmlIgnoreUTF8Errors *bool) *LegacyFileFormatTypeOptionsRequest {
	s.XMLIgnoreUTF8Errors = xmlIgnoreUTF8Errors
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithXMLPreserveSpace(xmlPreserveSpace *bool) *LegacyFileFormatTypeOptionsRequest {
	s.XMLPreserveSpace = xmlPreserveSpace
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithXMLStripOuterElement(xmlStripOuterElement *bool) *LegacyFileFormatTypeOptionsRequest {
	s.XMLStripOuterElement = xmlStripOuterElement
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithXMLDisableSnowflakeData(xmlDisableSnowflakeData *bool) *LegacyFileFormatTypeOptionsRequest {
	s.XMLDisableSnowflakeData = xmlDisableSnowflakeData
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithXMLDisableAutoConvert(xmlDisableAutoConvert *bool) *LegacyFileFormatTypeOptionsRequest {
	s.XMLDisableAutoConvert = xmlDisableAutoConvert
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithXMLReplaceInvalidCharacters(xmlReplaceInvalidCharacters *bool) *LegacyFileFormatTypeOptionsRequest {
	s.XMLReplaceInvalidCharacters = xmlReplaceInvalidCharacters
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithXMLSkipByteOrderMark(xmlSkipByteOrderMark *bool) *LegacyFileFormatTypeOptionsRequest {
	s.XMLSkipByteOrderMark = xmlSkipByteOrderMark
	return s
}

func (s *LegacyFileFormatTypeOptionsRequest) WithComment(comment *string) *LegacyFileFormatTypeOptionsRequest {
	s.Comment = comment
	return s
}

func NewTableClusteringActionRequest() *TableClusteringActionRequest {
	return &TableClusteringActionRequest{}
}

func (s *TableClusteringActionRequest) WithClusterBy(clusterBy []string) *TableClusteringActionRequest {
	s.ClusterBy = clusterBy
	return s
}

func (s *TableClusteringActionRequest) WithRecluster(recluster *TableReclusterActionRequest) *TableClusteringActionRequest {
	s.Recluster = recluster
	return s
}

func (s *TableClusteringActionRequest) WithChangeReclusterState(changeReclusterState *ReclusterState) *TableClusteringActionRequest {
	s.ChangeReclusterState = changeReclusterState
	return s
}

func (s *TableClusteringActionRequest) WithDropClusteringKey(dropClusteringKey *bool) *TableClusteringActionRequest {
	s.DropClusteringKey = dropClusteringKey
	return s
}

func NewTableReclusterActionRequest() *TableReclusterActionRequest {
	return &TableReclusterActionRequest{}
}

func (s *TableReclusterActionRequest) WithMaxSize(maxSize *int) *TableReclusterActionRequest {
	s.MaxSize = maxSize
	return s
}

func (s *TableReclusterActionRequest) WithCondition(condition *string) *TableReclusterActionRequest {
	s.Condition = condition
	return s
}

func NewTableReclusterChangeStateRequest() *TableReclusterChangeStateRequest {
	return &TableReclusterChangeStateRequest{}
}

func (s *TableReclusterChangeStateRequest) WithState(state ReclusterState) *TableReclusterChangeStateRequest {
	s.State = state
	return s
}

func NewTableColumnActionRequest() *TableColumnActionRequest {
	return &TableColumnActionRequest{}
}

func (s *TableColumnActionRequest) WithAdd(add *TableColumnAddActionRequest) *TableColumnActionRequest {
	s.Add = add
	return s
}

func (s *TableColumnActionRequest) WithRename(rename *TableColumnRenameActionRequest) *TableColumnActionRequest {
	s.Rename = rename
	return s
}

func (s *TableColumnActionRequest) WithAlter(alter []TableColumnAlterActionRequest) *TableColumnActionRequest {
	s.Alter = alter
	return s
}

func (s *TableColumnActionRequest) WithSetMaskingPolicy(setMaskingPolicy *TableColumnAlterSetMaskingPolicyActionRequest) *TableColumnActionRequest {
	s.SetMaskingPolicy = setMaskingPolicy
	return s
}

func (s *TableColumnActionRequest) WithUnsetMaskingPolicy(unsetMaskingPolicy *TableColumnAlterUnsetMaskingPolicyActionRequest) *TableColumnActionRequest {
	s.UnsetMaskingPolicy = unsetMaskingPolicy
	return s
}

func (s *TableColumnActionRequest) WithSetTags(setTags *TableColumnAlterSetTagsActionRequest) *TableColumnActionRequest {
	s.SetTags = setTags
	return s
}

func (s *TableColumnActionRequest) WithUnsetTags(unsetTags *TableColumnAlterUnsetTagsActionRequest) *TableColumnActionRequest {
	s.UnsetTags = unsetTags
	return s
}

func (s *TableColumnActionRequest) WithDropColumnsIfExists() *TableColumnActionRequest {
	s.DropColumnsIfExists = Bool(true)
	return s
}

func (s *TableColumnActionRequest) WithDropColumns(dropColumns []string) *TableColumnActionRequest {
	s.DropColumns = dropColumns
	return s
}

func NewTableColumnAddActionRequest(
	name string,
	dataType DataType,
) *TableColumnAddActionRequest {
	s := TableColumnAddActionRequest{}
	s.Name = name
	s.Type = dataType
	return &s
}

func (s *TableColumnAddActionRequest) WithIfNotExists() *TableColumnAddActionRequest {
	s.IfNotExists = Bool(true)
	return s
}

func (s *TableColumnAddActionRequest) WithDefaultValue(defaultValue *ColumnDefaultValueRequest) *TableColumnAddActionRequest {
	s.DefaultValue = defaultValue
	return s
}

func (s *TableColumnAddActionRequest) WithInlineConstraint(inlineConstraint *TableColumnAddInlineConstraintRequest) *TableColumnAddActionRequest {
	s.InlineConstraint = inlineConstraint
	return s
}

func (s *TableColumnAddActionRequest) WithMaskingPolicy(maskingPolicy *ColumnMaskingPolicyRequest) *TableColumnAddActionRequest {
	s.MaskingPolicy = maskingPolicy
	return s
}

func (s *TableColumnAddActionRequest) WithWith(with *bool) *TableColumnAddActionRequest {
	s.With = with
	return s
}

func (s *TableColumnAddActionRequest) WithTags(tags []TagAssociation) *TableColumnAddActionRequest {
	s.Tags = tags
	return s
}

func (s *TableColumnAddActionRequest) WithComment(comment *string) *TableColumnAddActionRequest {
	s.Comment = comment
	return s
}

func (s *TableColumnAddActionRequest) WithCollate(collate *string) *TableColumnAddActionRequest {
	s.Collate = collate
	return s
}

func NewTableColumnAddInlineConstraintRequest() *TableColumnAddInlineConstraintRequest {
	return &TableColumnAddInlineConstraintRequest{}
}

func (s *TableColumnAddInlineConstraintRequest) WithNotNull(notNull *bool) *TableColumnAddInlineConstraintRequest {
	s.NotNull = notNull
	return s
}

func (s *TableColumnAddInlineConstraintRequest) WithName(name *string) *TableColumnAddInlineConstraintRequest {
	s.Name = name
	return s
}

func (s *TableColumnAddInlineConstraintRequest) WithType(constraintType ColumnConstraintType) *TableColumnAddInlineConstraintRequest {
	s.Type = constraintType
	return s
}

func (s *TableColumnAddInlineConstraintRequest) WithForeignKey(foreignKey *ColumnAddForeignKey) *TableColumnAddInlineConstraintRequest {
	s.ForeignKey = foreignKey
	return s
}

func NewColumnAddForeignKeyRequest() *ColumnAddForeignKeyRequest {
	return &ColumnAddForeignKeyRequest{}
}

func (s *ColumnAddForeignKeyRequest) WithTableName(tableName string) *ColumnAddForeignKeyRequest {
	s.TableName = tableName
	return s
}

func (s *ColumnAddForeignKeyRequest) WithColumnName(columnName string) *ColumnAddForeignKeyRequest {
	s.ColumnName = columnName
	return s
}

func NewTableColumnRenameActionRequest(
	oldName string,
	newName string,
) *TableColumnRenameActionRequest {
	s := TableColumnRenameActionRequest{}
	s.OldName = oldName
	s.NewName = newName
	return &s
}

func NewTableColumnAlterActionRequest(
	name string,
) *TableColumnAlterActionRequest {
	s := TableColumnAlterActionRequest{}
	s.Name = name
	return &s
}

func (s *TableColumnAlterActionRequest) WithDropDefault(dropDefault *bool) *TableColumnAlterActionRequest {
	s.DropDefault = dropDefault
	return s
}

func (s *TableColumnAlterActionRequest) WithSetDefault(setDefault *SequenceName) *TableColumnAlterActionRequest {
	s.SetDefault = setDefault
	return s
}

func (s *TableColumnAlterActionRequest) WithNotNullConstraint(notNullConstraint *TableColumnNotNullConstraintRequest) *TableColumnAlterActionRequest {
	s.NotNullConstraint = notNullConstraint
	return s
}

func (s *TableColumnAlterActionRequest) WithType(dataType *DataType) *TableColumnAlterActionRequest {
	s.Type = dataType
	return s
}

func (s *TableColumnAlterActionRequest) WithComment(comment *string) *TableColumnAlterActionRequest {
	s.Comment = comment
	return s
}

func (s *TableColumnAlterActionRequest) WithCollate(collate *string) *TableColumnAlterActionRequest {
	s.Collate = collate
	return s
}

func (s *TableColumnAlterActionRequest) WithUnsetComment(unsetComment *bool) *TableColumnAlterActionRequest {
	s.UnsetComment = unsetComment
	return s
}

func NewTableColumnAlterSetMaskingPolicyActionRequest(
	columnName string,
	maskingPolicyName SchemaObjectIdentifier,
	using []string,
) *TableColumnAlterSetMaskingPolicyActionRequest {
	s := TableColumnAlterSetMaskingPolicyActionRequest{}
	s.ColumnName = columnName
	s.MaskingPolicyName = maskingPolicyName
	s.Using = using
	return &s
}

func (s *TableColumnAlterSetMaskingPolicyActionRequest) WithForce(force *bool) *TableColumnAlterSetMaskingPolicyActionRequest {
	s.Force = force
	return s
}

func NewTableColumnAlterUnsetMaskingPolicyActionRequest(
	columnName string,
) *TableColumnAlterUnsetMaskingPolicyActionRequest {
	s := TableColumnAlterUnsetMaskingPolicyActionRequest{}
	s.ColumnName = columnName
	return &s
}

func NewTableColumnAlterSetTagsActionRequest(
	columnName string,
	tags []TagAssociation,
) *TableColumnAlterSetTagsActionRequest {
	s := TableColumnAlterSetTagsActionRequest{}
	s.ColumnName = columnName
	s.Tags = tags
	return &s
}

func NewTableColumnAlterUnsetTagsActionRequest(
	columnName string,
	tags []ObjectIdentifier,
) *TableColumnAlterUnsetTagsActionRequest {
	s := TableColumnAlterUnsetTagsActionRequest{}
	s.ColumnName = columnName
	s.Tags = tags
	return &s
}

func NewTableColumnNotNullConstraintRequest() *TableColumnNotNullConstraintRequest {
	return &TableColumnNotNullConstraintRequest{}
}

func (s *TableColumnNotNullConstraintRequest) WithSet(set *bool) *TableColumnNotNullConstraintRequest {
	s.Set = set
	return s
}

func (s *TableColumnNotNullConstraintRequest) WithDrop(drop *bool) *TableColumnNotNullConstraintRequest {
	s.Drop = drop
	return s
}

func NewTableConstraintActionRequest() *TableConstraintActionRequest {
	return &TableConstraintActionRequest{}
}

func (s *TableConstraintActionRequest) WithAdd(add *OutOfLineConstraintRequest) *TableConstraintActionRequest {
	s.Add = add
	return s
}

func (s *TableConstraintActionRequest) WithRename(rename *TableConstraintRenameActionRequest) *TableConstraintActionRequest {
	s.Rename = rename
	return s
}

func (s *TableConstraintActionRequest) WithAlter(alter *TableConstraintAlterActionRequest) *TableConstraintActionRequest {
	s.Alter = alter
	return s
}

func (s *TableConstraintActionRequest) WithDrop(drop *TableConstraintDropActionRequest) *TableConstraintActionRequest {
	s.Drop = drop
	return s
}

func NewTableConstraintRenameActionRequest() *TableConstraintRenameActionRequest {
	return &TableConstraintRenameActionRequest{}
}

func (s *TableConstraintRenameActionRequest) WithOldName(oldName string) *TableConstraintRenameActionRequest {
	s.OldName = oldName
	return s
}

func (s *TableConstraintRenameActionRequest) WithNewName(newName string) *TableConstraintRenameActionRequest {
	s.NewName = newName
	return s
}

func NewTableConstraintAlterActionRequest() *TableConstraintAlterActionRequest {
	return &TableConstraintAlterActionRequest{}
}

func (s *TableConstraintAlterActionRequest) WithConstraintName(constraintName *string) *TableConstraintAlterActionRequest {
	s.ConstraintName = constraintName
	return s
}

func (s *TableConstraintAlterActionRequest) WithPrimaryKey(primaryKey *bool) *TableConstraintAlterActionRequest {
	s.PrimaryKey = primaryKey
	return s
}

func (s *TableConstraintAlterActionRequest) WithUnique(unique *bool) *TableConstraintAlterActionRequest {
	s.Unique = unique
	return s
}

func (s *TableConstraintAlterActionRequest) WithForeignKey(foreignKey *bool) *TableConstraintAlterActionRequest {
	s.ForeignKey = foreignKey
	return s
}

func (s *TableConstraintAlterActionRequest) WithColumns(columns []string) *TableConstraintAlterActionRequest {
	s.Columns = columns
	return s
}

func (s *TableConstraintAlterActionRequest) WithEnforced(enforced *bool) *TableConstraintAlterActionRequest {
	s.Enforced = enforced
	return s
}

func (s *TableConstraintAlterActionRequest) WithNotEnforced(notEnforced *bool) *TableConstraintAlterActionRequest {
	s.NotEnforced = notEnforced
	return s
}

func (s *TableConstraintAlterActionRequest) WithValidate(validate *bool) *TableConstraintAlterActionRequest {
	s.Validate = validate
	return s
}

func (s *TableConstraintAlterActionRequest) WithNoValidate(noValidate *bool) *TableConstraintAlterActionRequest {
	s.NoValidate = noValidate
	return s
}

func (s *TableConstraintAlterActionRequest) WithRely(rely *bool) *TableConstraintAlterActionRequest {
	s.Rely = rely
	return s
}

func (s *TableConstraintAlterActionRequest) WithNoRely(noRely *bool) *TableConstraintAlterActionRequest {
	s.NoRely = noRely
	return s
}

func NewTableConstraintDropActionRequest() *TableConstraintDropActionRequest {
	return &TableConstraintDropActionRequest{}
}

func (s *TableConstraintDropActionRequest) WithConstraintName(constraintName *string) *TableConstraintDropActionRequest {
	s.ConstraintName = constraintName
	return s
}

func (s *TableConstraintDropActionRequest) WithPrimaryKey(primaryKey *bool) *TableConstraintDropActionRequest {
	s.PrimaryKey = primaryKey
	return s
}

func (s *TableConstraintDropActionRequest) WithUnique(unique *bool) *TableConstraintDropActionRequest {
	s.Unique = unique
	return s
}

func (s *TableConstraintDropActionRequest) WithForeignKey(foreignKey *bool) *TableConstraintDropActionRequest {
	s.ForeignKey = foreignKey
	return s
}

func (s *TableConstraintDropActionRequest) WithColumns(columns []string) *TableConstraintDropActionRequest {
	s.Columns = columns
	return s
}

func (s *TableConstraintDropActionRequest) WithCascade(cascade *bool) *TableConstraintDropActionRequest {
	s.Cascade = cascade
	return s
}

func (s *TableConstraintDropActionRequest) WithRestrict(restrict *bool) *TableConstraintDropActionRequest {
	s.Restrict = restrict
	return s
}

func NewTableExternalTableActionRequest() *TableExternalTableActionRequest {
	return &TableExternalTableActionRequest{}
}

func (s *TableExternalTableActionRequest) WithAdd(add *TableExternalTableColumnAddActionRequest) *TableExternalTableActionRequest {
	s.Add = add
	return s
}

func (s *TableExternalTableActionRequest) WithRename(rename *TableExternalTableColumnRenameActionRequest) *TableExternalTableActionRequest {
	s.Rename = rename
	return s
}

func (s *TableExternalTableActionRequest) WithDrop(drop *TableExternalTableColumnDropActionRequest) *TableExternalTableActionRequest {
	s.Drop = drop
	return s
}

func NewTableSearchOptimizationActionRequest() *TableSearchOptimizationActionRequest {
	return &TableSearchOptimizationActionRequest{}
}

func (s *TableSearchOptimizationActionRequest) WithAddSearchOptimizationOn(addSearchOptimizationOn []string) *TableSearchOptimizationActionRequest {
	s.AddSearchOptimizationOn = addSearchOptimizationOn
	return s
}

func (s *TableSearchOptimizationActionRequest) WithDropSearchOptimizationOn(dropSearchOptimizationOn []string) *TableSearchOptimizationActionRequest {
	s.DropSearchOptimizationOn = dropSearchOptimizationOn
	return s
}

func NewTableSetRequest() *TableSetRequest {
	return &TableSetRequest{}
}

func (s *TableSetRequest) WithEnableSchemaEvolution(enableSchemaEvolution *bool) *TableSetRequest {
	s.EnableSchemaEvolution = enableSchemaEvolution
	return s
}

func (s *TableSetRequest) WithStageFileFormat(stageFileFormat LegacyFileFormatRequest) *TableSetRequest {
	s.StageFileFormat = &stageFileFormat
	return s
}

func (s *TableSetRequest) WithLegacyTableCopyOptions(LegacyTableCopyOptions LegacyTableCopyOptionsRequest) *TableSetRequest {
	s.LegacyTableCopyOptions = &LegacyTableCopyOptions
	return s
}

func (s *TableSetRequest) WithDataRetentionTimeInDays(dataRetentionTimeInDays *int) *TableSetRequest {
	s.DataRetentionTimeInDays = dataRetentionTimeInDays
	return s
}

func (s *TableSetRequest) WithMaxDataExtensionTimeInDays(maxDataExtensionTimeInDays *int) *TableSetRequest {
	s.MaxDataExtensionTimeInDays = maxDataExtensionTimeInDays
	return s
}

func (s *TableSetRequest) WithChangeTracking(changeTracking *bool) *TableSetRequest {
	s.ChangeTracking = changeTracking
	return s
}

func (s *TableSetRequest) WithDefaultDDLCollation(defaultDDLCollation *string) *TableSetRequest {
	s.DefaultDDLCollation = defaultDDLCollation
	return s
}

func (s *TableSetRequest) WithComment(comment *string) *TableSetRequest {
	s.Comment = comment
	return s
}

func NewTableExternalTableColumnAddActionRequest() *TableExternalTableColumnAddActionRequest {
	return &TableExternalTableColumnAddActionRequest{}
}

func (s *TableExternalTableColumnAddActionRequest) WithIfNotExists() *TableExternalTableColumnAddActionRequest {
	s.IfNotExists = Bool(true)
	return s
}

func (s *TableExternalTableColumnAddActionRequest) WithName(name string) *TableExternalTableColumnAddActionRequest {
	s.Name = name
	return s
}

func (s *TableExternalTableColumnAddActionRequest) WithType(dataType DataType) *TableExternalTableColumnAddActionRequest {
	s.Type = dataType
	return s
}

func (s *TableExternalTableColumnAddActionRequest) WithExpression(expression string) *TableExternalTableColumnAddActionRequest {
	s.Expression = expression
	return s
}

func (s *TableExternalTableColumnAddActionRequest) WithComment(comment *string) *TableExternalTableColumnAddActionRequest {
	s.Comment = comment
	return s
}

func NewTableExternalTableColumnRenameActionRequest() *TableExternalTableColumnRenameActionRequest {
	return &TableExternalTableColumnRenameActionRequest{}
}

func (s *TableExternalTableColumnRenameActionRequest) WithOldName(oldName string) *TableExternalTableColumnRenameActionRequest {
	s.OldName = oldName
	return s
}

func (s *TableExternalTableColumnRenameActionRequest) WithNewName(newName string) *TableExternalTableColumnRenameActionRequest {
	s.NewName = newName
	return s
}

func NewTableExternalTableColumnDropActionRequest(columns []string) *TableExternalTableColumnDropActionRequest {
	return &TableExternalTableColumnDropActionRequest{
		Columns: columns,
	}
}

func (s *TableExternalTableColumnDropActionRequest) WithIfExists() *TableExternalTableColumnDropActionRequest {
	s.IfExists = Bool(true)
	return s
}

func NewShowTableRequest() *ShowTableRequest {
	return &ShowTableRequest{}
}

func (s *ShowTableRequest) WithTerse(terse bool) *ShowTableRequest {
	s.Terse = &terse
	return s
}

func (s *ShowTableRequest) WithHistory(history *bool) *ShowTableRequest {
	s.History = history
	return s
}

func (s *ShowTableRequest) WithLike(like Like) *ShowTableRequest {
	s.Like = &like
	return s
}

func (s *ShowTableRequest) WithIn(in ExtendedIn) *ShowTableRequest {
	s.In = &in
	return s
}

func (s *ShowTableRequest) WithStartsWith(startsWith string) *ShowTableRequest {
	s.StartsWith = &startsWith
	return s
}

func (s *ShowTableRequest) WithLimitFrom(limit LimitFrom) *ShowTableRequest {
	s.Limit = &limit
	return s
}

func NewLimitFromRequest() *LimitFromRequest {
	return &LimitFromRequest{}
}

func (s *LimitFromRequest) WithRows(rows *int) *LimitFromRequest {
	s.rows = rows
	return s
}

func (s *LimitFromRequest) WithFrom(from *string) *LimitFromRequest {
	s.from = from
	return s
}

func NewDescribeTableColumnsRequest(
	id SchemaObjectIdentifier,
) *DescribeTableColumnsRequest {
	s := DescribeTableColumnsRequest{}
	s.id = id
	return &s
}

func NewDescribeTableStageRequest(
	id SchemaObjectIdentifier,
) *DescribeTableStageRequest {
	s := DescribeTableStageRequest{}
	s.id = id
	return &s
}

func NewLegacyTableCopyOptionsRequest() *LegacyTableCopyOptionsRequest {
	return &LegacyTableCopyOptionsRequest{}
}

func (s *LegacyTableCopyOptionsRequest) WithOnError(onError LegacyTableCopyOnErrorOptionsRequest) *LegacyTableCopyOptionsRequest {
	s.OnError = &onError
	return s
}

func (s *LegacyTableCopyOptionsRequest) WithSizeLimit(sizeLimit int) *LegacyTableCopyOptionsRequest {
	s.SizeLimit = &sizeLimit
	return s
}

func (s *LegacyTableCopyOptionsRequest) WithPurge(purge bool) *LegacyTableCopyOptionsRequest {
	s.Purge = &purge
	return s
}

func (s *LegacyTableCopyOptionsRequest) WithReturnFailedOnly(returnFailedOnly bool) *LegacyTableCopyOptionsRequest {
	s.ReturnFailedOnly = &returnFailedOnly
	return s
}

func (s *LegacyTableCopyOptionsRequest) WithMatchByColumnName(matchByColumnName StageCopyColumnMapOption) *LegacyTableCopyOptionsRequest {
	s.MatchByColumnName = &matchByColumnName
	return s
}

func (s *LegacyTableCopyOptionsRequest) WithEnforceLength(enforceLength bool) *LegacyTableCopyOptionsRequest {
	s.EnforceLength = &enforceLength
	return s
}

func (s *LegacyTableCopyOptionsRequest) WithTruncatecolumns(truncatecolumns bool) *LegacyTableCopyOptionsRequest {
	s.Truncatecolumns = &truncatecolumns
	return s
}

func (s *LegacyTableCopyOptionsRequest) WithForce(force bool) *LegacyTableCopyOptionsRequest {
	s.Force = &force
	return s
}

// Builder functions for LegacyTableCopyOnErrorOptionsRequest

func NewLegacyTableCopyOnErrorOptionsRequest() *LegacyTableCopyOnErrorOptionsRequest {
	return &LegacyTableCopyOnErrorOptionsRequest{}
}

func (s *LegacyTableCopyOnErrorOptionsRequest) WithContinue_(continue_ bool) *LegacyTableCopyOnErrorOptionsRequest {
	s.Continue_ = &continue_
	return s
}

func (s *LegacyTableCopyOnErrorOptionsRequest) WithAbortStatement(abortStatement bool) *LegacyTableCopyOnErrorOptionsRequest {
	s.AbortStatement = &abortStatement
	return s
}

// WithSkipFile sets SkipFile to "SKIP_FILE"
func (s *LegacyTableCopyOnErrorOptionsRequest) WithSkipFile() *LegacyTableCopyOnErrorOptionsRequest {
	s.SkipFile = String("SKIP_FILE")
	return s
}

// WithSkipFileX sets SkipFile to "SKIP_FILE_n" where n is the provided integer
func (s *LegacyTableCopyOnErrorOptionsRequest) WithSkipFileX(x int) *LegacyTableCopyOnErrorOptionsRequest {
	s.SkipFile = String(fmt.Sprintf("SKIP_FILE_%d", x))
	return s
}

// WithSkipFileXPercent sets SkipFile to "'SKIP_FILE_n%'" where n is the provided integer
func (s *LegacyTableCopyOnErrorOptionsRequest) WithSkipFileXPercent(x int) *LegacyTableCopyOnErrorOptionsRequest {
	s.SkipFile = String(fmt.Sprintf("'SKIP_FILE_%d%%'", x))
	return s
}

func NewLegacyFileFormatRequest() *LegacyFileFormatRequest {
	s := LegacyFileFormatRequest{}
	return &s
}

func (s *LegacyFileFormatRequest) WithFormatName(formatName string) *LegacyFileFormatRequest {
	s.FormatName = &formatName
	return s
}

func (s *LegacyFileFormatRequest) WithFileFormatType(fileFormatType FileFormatType) *LegacyFileFormatRequest {
	s.FileFormatType = &fileFormatType
	return s
}

// adjusted manually
func (s *LegacyFileFormatRequest) WithOptions(options LegacyFileFormatTypeOptionsRequest) *LegacyFileFormatRequest {
	s.Options = &options
	return s
}
