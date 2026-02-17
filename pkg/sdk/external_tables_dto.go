package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateExternalTableOptions]                       = new(CreateExternalTableRequest)
	_ optionsProvider[CreateWithManualPartitioningExternalTableOptions] = new(CreateWithManualPartitioningExternalTableRequest)
	_ optionsProvider[CreateDeltaLakeExternalTableOptions]              = new(CreateDeltaLakeExternalTableRequest)
	_ optionsProvider[CreateExternalTableUsingTemplateOptions]          = new(CreateExternalTableUsingTemplateRequest)
	_ optionsProvider[AlterExternalTableOptions]                        = new(AlterExternalTableRequest)
	_ optionsProvider[AlterExternalTablePartitionOptions]               = new(AlterExternalTablePartitionRequest)
	_ optionsProvider[DropExternalTableOptions]                         = new(DropExternalTableRequest)
	_ optionsProvider[ShowExternalTableOptions]                         = new(ShowExternalTableRequest)
	_ optionsProvider[describeExternalTableColumnsOptions]              = new(DescribeExternalTableColumnsRequest)
	_ optionsProvider[describeExternalTableStageOptions]                = new(DescribeExternalTableStageRequest)
)

type CreateExternalTableRequest struct {
	orReplace           *bool
	ifNotExists         *bool
	name                SchemaObjectIdentifier // required
	columns             []*ExternalTableColumnRequest
	cloudProviderParams *CloudProviderParamsRequest
	partitionBy         []string
	location            string // required
	refreshOnCreate     *bool
	autoRefresh         *bool
	pattern             *string
	rawFileFormat       *string
	fileFormat          *ExternalTableFileFormatRequest
	awsSnsTopic         *string
	copyGrants          *bool
	comment             *string
	rowAccessPolicy     *RowAccessPolicyRequest
	tag                 []*TagAssociationRequest
}

func (r *CreateExternalTableRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

func (s *CreateExternalTableRequest) GetColumns() []*ExternalTableColumnRequest {
	return s.columns
}

type ExternalTableColumnRequest struct {
	name             string   // required
	dataType         DataType // required
	asExpression     string   // required
	notNull          *bool
	inlineConstraint *ColumnInlineConstraintRequest
}

func (v ExternalTableColumnRequest) toOpts() ExternalTableColumn {
	var inlineConstraint *ColumnInlineConstraint
	if v.inlineConstraint != nil {
		inlineConstraint = v.inlineConstraint.toOpts()
	}

	return ExternalTableColumn{
		Name:             v.name,
		Type:             v.dataType,
		AsExpression:     []string{v.asExpression},
		NotNull:          v.notNull,
		InlineConstraint: inlineConstraint,
	}
}

func (v *ColumnInlineConstraintRequest) toOpts() *ColumnInlineConstraint {
	return &ColumnInlineConstraint{
		Name:               &v.Name,
		Type:               v.type_,
		ForeignKey:         v.foreignKey.toOpts(),
		Enforced:           v.enforced,
		NotEnforced:        v.notEnforced,
		Deferrable:         v.deferrable,
		NotDeferrable:      v.notDeferrable,
		InitiallyDeferred:  v.initiallyDeferred,
		InitiallyImmediate: v.initiallyImmediate,
		Enable:             v.enable,
		Disable:            v.disable,
		Validate:           v.validate,
		NoValidate:         v.noValidate,
		Rely:               v.rely,
		NoRely:             v.noRely,
	}
}

func (v *InlineForeignKeyRequest) toOpts() *InlineForeignKey {
	if v == nil {
		return nil
	}
	return &InlineForeignKey{
		TableName:  v.TableName,
		ColumnName: v.ColumnName,
		Match:      v.Match,
		On:         v.On,
	}
}

type CloudProviderParamsRequest struct {
	googleCloudStorageIntegration *string
	microsoftAzureIntegration     *string
}

func (v *CloudProviderParamsRequest) toOpts() *CloudProviderParams {
	return &CloudProviderParams{
		GoogleCloudStorageIntegration: v.googleCloudStorageIntegration,
		MicrosoftAzureIntegration:     v.microsoftAzureIntegration,
	}
}

type ExternalTableFileFormatRequest struct {
	name           *string
	fileFormatType *ExternalTableFileFormatType
	options        *ExternalTableFileFormatTypeOptionsRequest
}

func (v *ExternalTableFileFormatTypeOptionsRequest) toOpts() *ExternalTableFileFormatTypeOptions {
	var csvNullIf []NullString
	if v.csvNullIf != nil {
		for _, n := range *v.csvNullIf {
			csvNullIf = append(csvNullIf, n.toOpts())
		}
	}

	var orcNullIf []NullString
	if v.orcNullIf != nil {
		for _, n := range *v.orcNullIf {
			orcNullIf = append(orcNullIf, n.toOpts())
		}
	}

	return &ExternalTableFileFormatTypeOptions{
		CSVCompression:                  v.csvCompression,
		CSVRecordDelimiter:              v.csvRecordDelimiter,
		CSVFieldDelimiter:               v.csvFieldDelimiter,
		CSVSkipHeader:                   v.csvSkipHeader,
		CSVSkipBlankLines:               v.csvSkipBlankLines,
		CSVEscapeUnenclosedField:        v.csvEscapeUnenclosedField,
		CSVTrimSpace:                    v.csvTrimSpace,
		CSVFieldOptionallyEnclosedBy:    v.csvFieldOptionallyEnclosedBy,
		CSVNullIf:                       &csvNullIf,
		CSVEmptyFieldAsNull:             v.csvEmptyFieldAsNull,
		CSVEncoding:                     v.csvEncoding,
		JSONCompression:                 v.jsonCompression,
		JSONAllowDuplicate:              v.jsonAllowDuplicate,
		JSONStripOuterArray:             v.jsonStripOuterArray,
		JSONStripNullValues:             v.jsonStripNullValues,
		JSONReplaceInvalidCharacters:    v.jsonReplaceInvalidCharacters,
		AvroCompression:                 v.avroCompression,
		AvroReplaceInvalidCharacters:    v.avroReplaceInvalidCharacters,
		ORCTrimSpace:                    v.orcTrimSpace,
		ORCReplaceInvalidCharacters:     v.orcReplaceInvalidCharacters,
		ORCNullIf:                       &orcNullIf,
		ParquetCompression:              v.parquetCompression,
		ParquetBinaryAsText:             v.parquetBinaryAsText,
		ParquetReplaceInvalidCharacters: v.parquetReplaceInvalidCharacters,
	}
}

func (v ExternalTableFileFormatRequest) toOpts() ExternalTableFileFormat {
	var options *ExternalTableFileFormatTypeOptions
	if v.options != nil {
		options = v.options.toOpts()
	}

	return ExternalTableFileFormat{
		Name:    v.name,
		Type:    v.fileFormatType,
		Options: options,
	}
}

type ExternalTableFileFormatTypeOptionsRequest struct {
	// CSV type options
	csvCompression               *ExternalTableCsvCompression
	csvRecordDelimiter           *string
	csvFieldDelimiter            *string
	csvSkipHeader                *int
	csvSkipBlankLines            *bool
	csvEscapeUnenclosedField     *string
	csvTrimSpace                 *bool
	csvFieldOptionallyEnclosedBy *string
	csvNullIf                    *[]NullStringRequest
	csvEmptyFieldAsNull          *bool
	csvEncoding                  *CsvEncoding

	// JSON type options
	jsonCompression              *ExternalTableJsonCompression
	jsonAllowDuplicate           *bool
	jsonStripOuterArray          *bool
	jsonStripNullValues          *bool
	jsonReplaceInvalidCharacters *bool

	// AVRO type options
	avroCompression              *ExternalTableAvroCompression
	avroReplaceInvalidCharacters *bool

	// ORC type options
	orcTrimSpace                *bool
	orcReplaceInvalidCharacters *bool
	orcNullIf                   *[]NullStringRequest

	// PARQUET type options
	parquetCompression              *ExternalTableParquetCompression
	parquetBinaryAsText             *bool
	parquetReplaceInvalidCharacters *bool
}

type NullStringRequest struct {
	str string
}

func (v NullStringRequest) toOpts() NullString {
	return NullString{
		S: v.str,
	}
}

func (v *RowAccessPolicyRequest) toOpts() *TableRowAccessPolicy {
	return &TableRowAccessPolicy{
		Name: v.Name,
		On:   v.On,
	}
}

func (v TagAssociationRequest) toOpts() TagAssociation {
	return TagAssociation(v)
}

func (s *CreateExternalTableRequest) toOpts() *CreateExternalTableOptions {
	columns := make([]ExternalTableColumn, len(s.columns))
	if s.columns != nil {
		for i, c := range s.columns {
			columns[i] = c.toOpts()
		}
	}

	var fileFormat []ExternalTableFileFormat
	if s.fileFormat != nil {
		fileFormat = []ExternalTableFileFormat{s.fileFormat.toOpts()}
	}

	var rawFileFormat *RawFileFormat
	if s.rawFileFormat != nil {
		rawFileFormat = &RawFileFormat{
			Format: *s.rawFileFormat,
		}
	}

	var cloudProviderParams *CloudProviderParams
	if s.cloudProviderParams != nil {
		cloudProviderParams = s.cloudProviderParams.toOpts()
	}

	var rowAccessPolicy *TableRowAccessPolicy
	if s.rowAccessPolicy != nil {
		rowAccessPolicy = s.rowAccessPolicy.toOpts()
	}

	tag := make([]TagAssociation, len(s.tag))
	if s.tag != nil {
		for i, t := range s.tag {
			tag[i] = t.toOpts()
		}
	}

	return &CreateExternalTableOptions{
		OrReplace:           s.orReplace,
		IfNotExists:         s.ifNotExists,
		name:                s.name,
		Columns:             columns,
		CloudProviderParams: cloudProviderParams,
		PartitionBy:         s.partitionBy,
		Location:            s.location,
		RefreshOnCreate:     s.refreshOnCreate,
		AutoRefresh:         s.autoRefresh,
		Pattern:             s.pattern,
		RawFileFormat:       rawFileFormat,
		FileFormat:          fileFormat,
		AwsSnsTopic:         s.awsSnsTopic,
		CopyGrants:          s.copyGrants,
		Comment:             s.comment,
		RowAccessPolicy:     rowAccessPolicy,
		Tag:                 tag,
	}
}

type CreateWithManualPartitioningExternalTableRequest struct {
	orReplace           *bool
	ifNotExists         *bool
	name                SchemaObjectIdentifier // required
	columns             []*ExternalTableColumnRequest
	cloudProviderParams *CloudProviderParamsRequest
	partitionBy         []string
	location            string // required
	rawFileFormat       *string
	fileFormat          *ExternalTableFileFormatRequest
	copyGrants          *bool
	comment             *string
	rowAccessPolicy     *RowAccessPolicyRequest
	tag                 []*TagAssociationRequest
}

func (v *CreateWithManualPartitioningExternalTableRequest) toOpts() *CreateWithManualPartitioningExternalTableOptions {
	columns := make([]ExternalTableColumn, len(v.columns))
	if v.columns != nil {
		for i, c := range v.columns {
			columns[i] = c.toOpts()
		}
	}

	var cloudProviderParams *CloudProviderParams
	if v.cloudProviderParams != nil {
		cloudProviderParams = v.cloudProviderParams.toOpts()
	}

	var fileFormat []ExternalTableFileFormat
	if v.fileFormat != nil {
		fileFormat = []ExternalTableFileFormat{v.fileFormat.toOpts()}
	}

	var rawFileFormat *RawFileFormat
	if v.rawFileFormat != nil {
		rawFileFormat = &RawFileFormat{
			Format: *v.rawFileFormat,
		}
	}

	var rowAccessPolicy *TableRowAccessPolicy
	if v.rowAccessPolicy != nil {
		rowAccessPolicy = v.rowAccessPolicy.toOpts()
	}

	tag := make([]TagAssociation, len(v.tag))
	if v.tag != nil {
		for i, t := range v.tag {
			tag[i] = t.toOpts()
		}
	}

	return &CreateWithManualPartitioningExternalTableOptions{
		OrReplace:           v.orReplace,
		IfNotExists:         v.ifNotExists,
		name:                v.name,
		Columns:             columns,
		CloudProviderParams: cloudProviderParams,
		PartitionBy:         v.partitionBy,
		Location:            v.location,
		RawFileFormat:       rawFileFormat,
		FileFormat:          fileFormat,
		CopyGrants:          v.copyGrants,
		Comment:             v.comment,
		RowAccessPolicy:     rowAccessPolicy,
		Tag:                 tag,
	}
}

type CreateDeltaLakeExternalTableRequest struct {
	orReplace           *bool
	ifNotExists         *bool
	name                SchemaObjectIdentifier // required
	columns             []*ExternalTableColumnRequest
	cloudProviderParams *CloudProviderParamsRequest
	partitionBy         []string
	location            string // required
	refreshOnCreate     *bool
	autoRefresh         *bool
	rawFileFormat       *string
	fileFormat          *ExternalTableFileFormatRequest
	copyGrants          *bool
	comment             *string
	rowAccessPolicy     *RowAccessPolicyRequest
	tag                 []*TagAssociationRequest
}

func (v *CreateDeltaLakeExternalTableRequest) toOpts() *CreateDeltaLakeExternalTableOptions {
	columns := make([]ExternalTableColumn, len(v.columns))
	if v.columns != nil {
		for i, c := range v.columns {
			columns[i] = c.toOpts()
		}
	}

	var cloudProviderParams *CloudProviderParams
	if v.cloudProviderParams != nil {
		cloudProviderParams = v.cloudProviderParams.toOpts()
	}

	var fileFormat []ExternalTableFileFormat
	if v.fileFormat != nil {
		fileFormat = []ExternalTableFileFormat{v.fileFormat.toOpts()}
	}

	var rawFileFormat *RawFileFormat
	if v.rawFileFormat != nil {
		rawFileFormat = &RawFileFormat{
			Format: *v.rawFileFormat,
		}
	}

	var rowAccessPolicy *TableRowAccessPolicy
	if v.rowAccessPolicy != nil {
		rowAccessPolicy = v.rowAccessPolicy.toOpts()
	}

	tag := make([]TagAssociation, len(v.tag))
	if v.tag != nil {
		for i, t := range v.tag {
			tag[i] = t.toOpts()
		}
	}

	return &CreateDeltaLakeExternalTableOptions{
		OrReplace:           v.orReplace,
		IfNotExists:         v.ifNotExists,
		name:                v.name,
		Columns:             columns,
		CloudProviderParams: cloudProviderParams,
		PartitionBy:         v.partitionBy,
		Location:            v.location,
		RefreshOnCreate:     v.refreshOnCreate,
		AutoRefresh:         v.autoRefresh,
		RawFileFormat:       rawFileFormat,
		FileFormat:          fileFormat,
		CopyGrants:          v.copyGrants,
		Comment:             v.comment,
		RowAccessPolicy:     rowAccessPolicy,
		Tag:                 tag,
	}
}

type CreateExternalTableUsingTemplateRequest struct {
	orReplace           *bool
	name                SchemaObjectIdentifier // required
	copyGrants          *bool
	query               string
	cloudProviderParams *CloudProviderParamsRequest
	partitionBy         []string
	location            string // required
	refreshOnCreate     *bool
	autoRefresh         *bool
	pattern             *string
	rawFileFormat       *string
	fileFormat          *ExternalTableFileFormatRequest
	awsSnsTopic         *string
	comment             *string
	rowAccessPolicy     *RowAccessPolicyRequest
	tag                 []*TagAssociationRequest
}

func (v *CreateExternalTableUsingTemplateRequest) toOpts() *CreateExternalTableUsingTemplateOptions {
	var cloudProviderParams *CloudProviderParams
	if v.cloudProviderParams != nil {
		cloudProviderParams = v.cloudProviderParams.toOpts()
	}

	var fileFormat []ExternalTableFileFormat
	if v.fileFormat != nil {
		fileFormat = []ExternalTableFileFormat{v.fileFormat.toOpts()}
	}

	var rawFileFormat *RawFileFormat
	if v.rawFileFormat != nil {
		rawFileFormat = &RawFileFormat{
			Format: *v.rawFileFormat,
		}
	}

	var rowAccessPolicy *TableRowAccessPolicy
	if v.rowAccessPolicy != nil {
		rowAccessPolicy = v.rowAccessPolicy.toOpts()
	}

	tag := make([]TagAssociation, len(v.tag))
	if v.tag != nil {
		for i, t := range v.tag {
			tag[i] = t.toOpts()
		}
	}

	return &CreateExternalTableUsingTemplateOptions{
		OrReplace:           v.orReplace,
		name:                v.name,
		CopyGrants:          v.copyGrants,
		Query:               []string{v.query},
		CloudProviderParams: cloudProviderParams,
		PartitionBy:         v.partitionBy,
		Location:            v.location,
		RefreshOnCreate:     v.refreshOnCreate,
		AutoRefresh:         v.autoRefresh,
		Pattern:             v.pattern,
		RawFileFormat:       rawFileFormat,
		FileFormat:          fileFormat,
		AwsSnsTopic:         v.awsSnsTopic,
		Comment:             v.comment,
		RowAccessPolicy:     rowAccessPolicy,
		Tag:                 tag,
	}
}

type AlterExternalTableRequest struct {
	ifExists    *bool
	name        SchemaObjectIdentifier // required
	refresh     *RefreshExternalTableRequest
	addFiles    []*ExternalTableFileRequest
	removeFiles []*ExternalTableFileRequest
	autoRefresh *bool
	setTag      []*TagAssociationRequest
	unsetTag    []ObjectIdentifier
}

type RefreshExternalTableRequest struct {
	path string // required
}

type ExternalTableFileRequest struct {
	name string // required
}

func (v *AlterExternalTableRequest) toOpts() *AlterExternalTableOptions {
	var refresh *RefreshExternalTable
	if v.refresh != nil {
		refresh = &RefreshExternalTable{
			Path: v.refresh.path,
		}
	}

	addFiles := make([]ExternalTableFile, len(v.addFiles))
	if v.addFiles != nil {
		for i, f := range v.addFiles {
			addFiles[i] = ExternalTableFile{
				Name: f.name,
			}
		}
	}

	removeFiles := make([]ExternalTableFile, len(v.removeFiles))
	if v.removeFiles != nil {
		for i, f := range v.removeFiles {
			removeFiles[i] = ExternalTableFile{
				Name: f.name,
			}
		}
	}

	setTag := make([]TagAssociation, len(v.setTag))
	if v.setTag != nil {
		for i, t := range v.setTag {
			setTag[i] = t.toOpts()
		}
	}

	return &AlterExternalTableOptions{
		IfExists:    v.ifExists,
		name:        v.name,
		Refresh:     refresh,
		AddFiles:    addFiles,
		RemoveFiles: removeFiles,
		AutoRefresh: v.autoRefresh,
		SetTag:      setTag,
		UnsetTag:    v.unsetTag,
	}
}

type AlterExternalTablePartitionRequest struct {
	ifExists      *bool
	name          SchemaObjectIdentifier // required
	addPartitions []*PartitionRequest
	dropPartition *bool
	location      string
}

type PartitionRequest struct {
	columnName string // required
	value      string // required
}

func (v *AlterExternalTablePartitionRequest) toOpts() *AlterExternalTablePartitionOptions {
	addPartitions := make([]Partition, len(v.addPartitions))
	if v.addPartitions != nil {
		for i, p := range v.addPartitions {
			addPartitions[i] = Partition{
				ColumnName: p.columnName,
				Value:      p.value,
			}
		}
	}

	return &AlterExternalTablePartitionOptions{
		IfExists:      v.ifExists,
		name:          v.name,
		AddPartitions: addPartitions,
		DropPartition: v.dropPartition,
		Location:      v.location,
	}
}

type DropExternalTableRequest struct {
	ifExists   *bool
	name       SchemaObjectIdentifier // required
	dropOption *ExternalTableDropOptionRequest
}

type ExternalTableDropOptionRequest struct {
	restrict *bool
	cascade  *bool
}

func (v *ExternalTableDropOptionRequest) toOpts() *ExternalTableDropOption {
	return &ExternalTableDropOption{
		Restrict: v.restrict,
		Cascade:  v.cascade,
	}
}

func (v *DropExternalTableRequest) toOpts() *DropExternalTableOptions {
	var dropOption *ExternalTableDropOption
	if v.dropOption != nil {
		dropOption = v.dropOption.toOpts()
	}

	return &DropExternalTableOptions{
		IfExists:   v.ifExists,
		name:       v.name,
		DropOption: dropOption,
	}
}

type ShowExternalTableRequest struct {
	terse      *bool
	like       *string
	in         *ShowExternalTableInRequest
	startsWith *string
	limitFrom  *LimitFromRequest
}

type ShowExternalTableInRequest struct {
	account  *bool
	database AccountObjectIdentifier
	schema   DatabaseObjectIdentifier
}

func (v *ShowExternalTableInRequest) toOpts() *In {
	return &In{
		Account:  v.account,
		Database: v.database,
		Schema:   v.schema,
	}
}

func (v *LimitFromRequest) toOpts() *LimitFrom {
	return &LimitFrom{
		Rows: v.rows,
		From: v.from,
	}
}

func (v *ShowExternalTableRequest) toOpts() *ShowExternalTableOptions {
	var like *Like
	if v.like != nil {
		like = &Like{
			Pattern: v.like,
		}
	}

	var in *In
	if v.in != nil {
		in = v.in.toOpts()
	}

	var limitFrom *LimitFrom
	if v.limitFrom != nil {
		limitFrom = v.limitFrom.toOpts()
	}

	return &ShowExternalTableOptions{
		Terse:      v.terse,
		Like:       like,
		In:         in,
		StartsWith: v.startsWith,
		LimitFrom:  limitFrom,
	}
}

type DescribeExternalTableColumnsRequest struct {
	id SchemaObjectIdentifier // required
}

type DescribeExternalTableStageRequest struct {
	id SchemaObjectIdentifier // required
}

func (v *DescribeExternalTableColumnsRequest) toOpts() *describeExternalTableColumnsOptions {
	return &describeExternalTableColumnsOptions{
		name: v.id,
	}
}

func (v *DescribeExternalTableStageRequest) toOpts() *describeExternalTableStageOptions {
	return &describeExternalTableStageOptions{
		name: v.id,
	}
}
