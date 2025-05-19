package sdk

import "time"

//go:generate go run ./dto-builder-generator/main.go

type CreateTableAsSelectRequest struct {
	orReplace *bool
	name      SchemaObjectIdentifier       // required
	columns   []TableAsSelectColumnRequest // required
	query     string                       // required
}

type TableAsSelectColumnRequest struct {
	orReplace         *bool
	name              string // required
	type_             *DataType
	maskingPolicyName *SchemaObjectIdentifier
	clusterBy         []string
	copyGrants        *bool
}

type CreateTableUsingTemplateRequest struct {
	orReplace  *bool
	name       SchemaObjectIdentifier // required
	copyGrants *bool
	Query      string // required
}

type CreateTableLikeRequest struct {
	orReplace   *bool
	name        SchemaObjectIdentifier // required
	sourceTable SchemaObjectIdentifier // required
	clusterBy   []string
	copyGrants  *bool
}

type CreateTableCloneRequest struct {
	orReplace   *bool
	name        SchemaObjectIdentifier // required
	sourceTable SchemaObjectIdentifier // required
	copyGrants  *bool
	ClonePoint  *ClonePointRequest
}

type ClonePointRequest struct {
	Moment CloneMoment
	At     TimeTravelRequest
}

type TimeTravelRequest struct {
	Timestamp *time.Time
	Offset    *int
	Statement *string
}

type CreateTableRequest struct {
	orReplace                  *bool
	ifNotExists                *bool
	scope                      *TableScope
	kind                       *TableKind
	name                       SchemaObjectIdentifier // required
	columns                    []TableColumnRequest   // required
	OutOfLineConstraints       []OutOfLineConstraintRequest
	clusterBy                  []string
	enableSchemaEvolution      *bool
	stageFileFormat            *StageFileFormatRequest
	stageCopyOptions           *StageCopyOptionsRequest
	DataRetentionTimeInDays    *int
	MaxDataExtensionTimeInDays *int
	ChangeTracking             *bool
	DefaultDDLCollation        *string
	CopyGrants                 *bool
	RowAccessPolicy            *RowAccessPolicyRequest
	Tags                       []TagAssociationRequest
	Comment                    *string
}

func (r *CreateTableRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

type RowAccessPolicyRequest struct {
	Name SchemaObjectIdentifier // required
	On   []string               // required
}

type TableColumnRequest struct {
	name             string   // required
	type_            DataType // required
	collate          *string
	comment          *string
	defaultValue     *ColumnDefaultValueRequest
	notNull          *bool
	maskingPolicy    *ColumnMaskingPolicyRequest
	with             *bool
	tags             []TagAssociation
	inlineConstraint *ColumnInlineConstraintRequest
}

type ColumnDefaultValueRequest struct {
	// One of
	expression *string
	identity   *ColumnIdentityRequest
}

type ColumnIdentityRequest struct {
	Start     int // required
	Increment int // required
	Order     *bool
	Noorder   *bool
}

type ColumnMaskingPolicyRequest struct {
	with  *bool
	name  SchemaObjectIdentifier // required
	using []string
}

type ColumnInlineConstraintRequest struct {
	Name               string               // required
	type_              ColumnConstraintType // required
	foreignKey         *InlineForeignKeyRequest
	enforced           *bool
	notEnforced        *bool
	deferrable         *bool
	notDeferrable      *bool
	initiallyDeferred  *bool
	initiallyImmediate *bool
	enable             *bool
	disable            *bool
	validate           *bool
	noValidate         *bool
	rely               *bool
	noRely             *bool
}

type OutOfLineConstraintRequest struct {
	Name       *string
	Type       ColumnConstraintType // required
	Columns    []string
	ForeignKey *OutOfLineForeignKeyRequest

	// Optional
	Enforced           *bool
	NotEnforced        *bool
	Deferrable         *bool
	NotDeferrable      *bool
	InitiallyDeferred  *bool
	InitiallyImmediate *bool
	Enable             *bool
	Disable            *bool
	Validate           *bool
	NoValidate         *bool
	Rely               *bool
	NoRely             *bool
}

type InlineForeignKeyRequest struct {
	TableName  string // required
	ColumnName []string
	Match      *MatchType
	On         *ForeignKeyOnAction
}

type OutOfLineForeignKeyRequest struct {
	TableName   SchemaObjectIdentifier // required
	ColumnNames []string               // required
	Match       *MatchType
	On          *ForeignKeyOnAction
}

type AlterTableRequest struct {
	IfExists                  *bool
	name                      SchemaObjectIdentifier // required
	NewName                   *SchemaObjectIdentifier
	SwapWith                  *SchemaObjectIdentifier
	ClusteringAction          *TableClusteringActionRequest
	ColumnAction              *TableColumnActionRequest
	ConstraintAction          *TableConstraintActionRequest
	ExternalTableAction       *TableExternalTableActionRequest
	SearchOptimizationAction  *TableSearchOptimizationActionRequest
	Set                       *TableSetRequest
	SetTags                   []TagAssociationRequest
	UnsetTags                 []ObjectIdentifier
	Unset                     *TableUnsetRequest
	AddRowAccessPolicy        *TableAddRowAccessPolicyRequest
	DropRowAccessPolicy       *TableDropRowAccessPolicyRequest
	DropAndAddRowAccessPolicy *TableDropAndAddRowAccessPolicy
	DropAllAccessRowPolicies  *bool
}

type DropTableRequest struct {
	IfExists *bool
	Name     SchemaObjectIdentifier // required
	// One of
	Cascade  *bool
	Restrict *bool
}

func (s *DropTableRequest) toOpts() *dropTableOptions {
	return &dropTableOptions{
		IfExists: s.IfExists,
		name:     s.Name,
		Cascade:  s.Cascade,
		Restrict: s.Restrict,
	}
}

func (s *ShowTableRequest) toOpts() *showTableOptions {
	var like *Like
	if s.Like != nil {
		like = &Like{
			Pattern: s.Like.Pattern,
		}
	}
	var limitFrom *LimitFrom
	if s.Limit != nil {
		limitFrom = &LimitFrom{
			Rows: s.Limit.Rows,
			From: s.Limit.From,
		}
	}
	return &showTableOptions{
		Terse:      s.Terse,
		History:    s.History,
		Like:       like,
		StartsWith: s.StartsWith,
		LimitFrom:  limitFrom,
		In:         s.In,
	}
}

type TableAddRowAccessPolicyRequest struct {
	RowAccessPolicy SchemaObjectIdentifier // required
	On              []string               // required
}

type TableDropRowAccessPolicyRequest struct {
	RowAccessPolicy SchemaObjectIdentifier // required
}

type TableDropAndAddRowAccessPolicyRequest struct {
	Drop TableDropRowAccessPolicyRequest // required
	Add  TableAddRowAccessPolicyRequest  // required
}

type TableUnsetRequest struct {
	DataRetentionTimeInDays    bool
	MaxDataExtensionTimeInDays bool
	ChangeTracking             bool
	DefaultDDLCollation        bool
	EnableSchemaEvolution      bool
	Comment                    bool
}

type AddRowAccessPolicyRequest struct {
	PolicyName string   // required
	ColumnName []string // required
}

type TagAssociationRequest struct {
	Name  ObjectIdentifier // required
	Value string           // required
}

type FileFormatTypeOptionsRequest struct {
	CSVCompression                *CSVCompression
	CSVRecordDelimiter            *string
	CSVFieldDelimiter             *string
	CSVFileExtension              *string
	CSVParseHeader                *bool
	CSVSkipHeader                 *int
	CSVSkipBlankLines             *bool
	CSVDateFormat                 *string
	CSVTimeFormat                 *string
	CSVTimestampFormat            *string
	CSVBinaryFormat               *BinaryFormat
	CSVEscape                     *string
	CSVEscapeUnenclosedField      *string
	CSVTrimSpace                  *bool
	CSVFieldOptionallyEnclosedBy  *string
	CSVNullIf                     *[]NullString
	CSVErrorOnColumnCountMismatch *bool
	CSVReplaceInvalidCharacters   *bool
	CSVEmptyFieldAsNull           *bool
	CSVSkipByteOrderMark          *bool
	CSVEncoding                   *CSVEncoding

	// JSON type options
	JSONCompression              *JSONCompression
	JSONDateFormat               *string
	JSONTimeFormat               *string
	JSONTimestampFormat          *string
	JSONBinaryFormat             *BinaryFormat
	JSONTrimSpace                *bool
	JSONNullIf                   []NullString
	JSONFileExtension            *string
	JSONEnableOctal              *bool
	JSONAllowDuplicate           *bool
	JSONStripOuterArray          *bool
	JSONStripNullValues          *bool
	JSONReplaceInvalidCharacters *bool
	JSONIgnoreUTF8Errors         *bool
	JSONSkipByteOrderMark        *bool

	// AVRO type options
	AvroCompression              *AvroCompression
	AvroTrimSpace                *bool
	AvroReplaceInvalidCharacters *bool
	AvroNullIf                   *[]NullString

	// ORC type options
	ORCTrimSpace                *bool
	ORCReplaceInvalidCharacters *bool
	ORCNullIf                   *[]NullString

	// PARQUET type options
	ParquetCompression              *ParquetCompression
	ParquetSnappyCompression        *bool
	ParquetBinaryAsText             *bool
	ParquetTrimSpace                *bool
	ParquetReplaceInvalidCharacters *bool
	ParquetNullIf                   *[]NullString

	// XML type options
	XMLCompression              *XMLCompression
	XMLIgnoreUTF8Errors         *bool
	XMLPreserveSpace            *bool
	XMLStripOuterElement        *bool
	XMLDisableSnowflakeData     *bool
	XMLDisableAutoConvert       *bool
	XMLReplaceInvalidCharacters *bool
	XMLSkipByteOrderMark        *bool

	Comment *string
}

type TableClusteringActionRequest struct {
	// One of
	ClusterBy            []string
	Recluster            *TableReclusterActionRequest
	ChangeReclusterState *ReclusterState
	DropClusteringKey    *bool
}

type TableReclusterActionRequest struct {
	MaxSize   *int
	Condition *string
}

type TableReclusterChangeStateRequest struct {
	State ReclusterState
}

type TableColumnActionRequest struct {
	Add                 *TableColumnAddActionRequest
	Rename              *TableColumnRenameActionRequest
	Alter               []TableColumnAlterActionRequest
	SetMaskingPolicy    *TableColumnAlterSetMaskingPolicyActionRequest
	UnsetMaskingPolicy  *TableColumnAlterUnsetMaskingPolicyActionRequest
	SetTags             *TableColumnAlterSetTagsActionRequest
	UnsetTags           *TableColumnAlterUnsetTagsActionRequest
	DropColumnsIfExists *bool
	DropColumns         []string
}

type TableColumnAddActionRequest struct {
	IfNotExists      *bool
	Name             string   // required
	Type             DataType // required
	DefaultValue     *ColumnDefaultValueRequest
	InlineConstraint *TableColumnAddInlineConstraintRequest
	MaskingPolicy    *ColumnMaskingPolicyRequest
	With             *bool
	Tags             []TagAssociation
	Comment          *string
	Collate          *string
}

type TableColumnAddInlineConstraintRequest struct {
	NotNull    *bool
	Name       *string
	Type       ColumnConstraintType
	ForeignKey *ColumnAddForeignKey
}

type ColumnAddForeignKeyRequest struct {
	TableName  string
	ColumnName string
}

type TableColumnRenameActionRequest struct {
	OldName string // required
	NewName string // required
}

type TableColumnAlterActionRequest struct {
	Name string // required

	// One of
	DropDefault       *bool
	SetDefault        *SequenceName
	NotNullConstraint *TableColumnNotNullConstraintRequest
	Type              *DataType
	Comment           *string
	UnsetComment      *bool
	Collate           *string
}

type TableColumnAlterSetMaskingPolicyActionRequest struct {
	ColumnName        string                 // required
	MaskingPolicyName SchemaObjectIdentifier // required
	Using             []string               // required
	Force             *bool
}

type TableColumnAlterUnsetMaskingPolicyActionRequest struct {
	ColumnName string // required
}

type TableColumnAlterSetTagsActionRequest struct {
	ColumnName string           // required
	Tags       []TagAssociation // required
}

type TableColumnAlterUnsetTagsActionRequest struct {
	ColumnName string             // required
	Tags       []ObjectIdentifier // required
}

type TableColumnNotNullConstraintRequest struct {
	Set  *bool
	Drop *bool
}

type TableConstraintActionRequest struct {
	Add    *OutOfLineConstraintRequest
	Rename *TableConstraintRenameActionRequest
	Alter  *TableConstraintAlterActionRequest
	Drop   *TableConstraintDropActionRequest
}

type TableConstraintRenameActionRequest struct {
	OldName string
	NewName string
}

type TableConstraintAlterActionRequest struct {
	// One of
	ConstraintName *string
	PrimaryKey     *bool
	Unique         *bool
	ForeignKey     *bool

	// Optional
	Columns     []string
	Enforced    *bool
	NotEnforced *bool
	Validate    *bool
	NoValidate  *bool
	Rely        *bool
	NoRely      *bool
}

type TableConstraintDropActionRequest struct {
	// One of
	ConstraintName *string
	PrimaryKey     *bool
	Unique         *bool
	ForeignKey     *bool

	// Optional
	Columns  []string
	Cascade  *bool
	Restrict *bool
}

type TableExternalTableActionRequest struct {
	// One of
	Add    *TableExternalTableColumnAddActionRequest
	Rename *TableExternalTableColumnRenameActionRequest
	Drop   *TableExternalTableColumnDropActionRequest
}

type TableSearchOptimizationActionRequest struct {
	// One of
	AddSearchOptimizationOn  []string
	DropSearchOptimizationOn []string
}

type TableSetRequest struct {
	EnableSchemaEvolution      *bool
	StageFileFormat            *StageFileFormatRequest
	StageCopyOptions           *StageCopyOptionsRequest
	DataRetentionTimeInDays    *int
	MaxDataExtensionTimeInDays *int
	ChangeTracking             *bool
	DefaultDDLCollation        *string
	Comment                    *string
}

type TableExternalTableColumnAddActionRequest struct {
	IfNotExists *bool
	Name        string
	Type        DataType
	Expression  string
	Comment     *string
}

type TableExternalTableColumnRenameActionRequest struct {
	OldName string
	NewName string
}

type TableExternalTableColumnDropActionRequest struct {
	Columns  []string // required
	IfExists *bool
}

type ShowTableRequest struct {
	Terse      *bool
	History    *bool
	Like       *Like
	In         *ExtendedIn
	StartsWith *string
	Limit      *LimitFrom
}

type ShowTableInRequest struct {
	Account  *bool
	Database AccountObjectIdentifier
	Schema   DatabaseObjectIdentifier
}

type LimitFromRequest struct {
	rows *int
	from *string
}

type DescribeTableColumnsRequest struct {
	id SchemaObjectIdentifier // required
}

type DescribeTableStageRequest struct {
	id SchemaObjectIdentifier // required
}
