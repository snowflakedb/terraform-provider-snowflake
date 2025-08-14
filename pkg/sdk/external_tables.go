package sdk

import (
	"context"
	"database/sql"
	"time"
)

type ExternalTables interface {
	Create(ctx context.Context, req *CreateExternalTableRequest) error
	CreateWithManualPartitioning(ctx context.Context, req *CreateWithManualPartitioningExternalTableRequest) error
	CreateDeltaLake(ctx context.Context, req *CreateDeltaLakeExternalTableRequest) error
	CreateUsingTemplate(ctx context.Context, req *CreateExternalTableUsingTemplateRequest) error
	// TODO: Add alter options from https://docs.snowflake.com/en/sql-reference/sql/alter-table#external-table-column-actions-exttablecolumnaction (for better UX)
	Alter(ctx context.Context, req *AlterExternalTableRequest) error
	AlterPartitions(ctx context.Context, req *AlterExternalTablePartitionRequest) error
	Drop(ctx context.Context, req *DropExternalTableRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, req *ShowExternalTableRequest) ([]ExternalTable, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*ExternalTable, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*ExternalTable, error)
	DescribeColumns(ctx context.Context, req *DescribeExternalTableColumnsRequest) ([]ExternalTableColumnDetails, error)
	DescribeStage(ctx context.Context, req *DescribeExternalTableStageRequest) ([]ExternalTableStageDetails, error)
}

type ExternalTable struct {
	CreatedOn           time.Time
	Name                string
	DatabaseName        string
	SchemaName          string
	Invalid             bool
	InvalidReason       string
	Owner               string
	Comment             string
	Stage               string
	Location            string
	FileFormatName      string
	FileFormatType      string
	Cloud               string
	Region              string
	NotificationChannel string
	LastRefreshedOn     time.Time
	TableFormat         string
	LastRefreshDetails  string
	OwnerRoleType       string
}

func (v *ExternalTable) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *ExternalTable) ObjectType() ObjectType {
	return ObjectTypeExternalTable
}

type RawFileFormat struct {
	Format string `ddl:"keyword"`
}

// CreateExternalTableOptions based on https://docs.snowflake.com/en/sql-reference/sql/create-external-table
type CreateExternalTableOptions struct {
	create              bool                   `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	externalTable       bool                   `ddl:"static" sql:"EXTERNAL TABLE"`
	IfNotExists         *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
	Columns             []ExternalTableColumn  `ddl:"list,parentheses"`
	CloudProviderParams *CloudProviderParams
	PartitionBy         []string                  `ddl:"keyword,parentheses" sql:"PARTITION BY"`
	Location            string                    `ddl:"parameter" sql:"LOCATION"`
	RefreshOnCreate     *bool                     `ddl:"parameter" sql:"REFRESH_ON_CREATE"`
	AutoRefresh         *bool                     `ddl:"parameter" sql:"AUTO_REFRESH"`
	Pattern             *string                   `ddl:"parameter,single_quotes" sql:"PATTERN"`
	FileFormat          []ExternalTableFileFormat `ddl:"parameter,parentheses" sql:"FILE_FORMAT"`
	// RawFileFormat was introduced, because of the decision taken during https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2228
	// that for now the snowflake_external_table resource should continue on using raw file format, which wasn't previously supported by the new SDK.
	// In the future it should most likely be replaced by a more structured version FileFormat
	RawFileFormat   *RawFileFormat        `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	AwsSnsTopic     *string               `ddl:"parameter,single_quotes" sql:"AWS_SNS_TOPIC"`
	CopyGrants      *bool                 `ddl:"keyword" sql:"COPY GRANTS"`
	Comment         *string               `ddl:"parameter,single_quotes" sql:"COMMENT"`
	RowAccessPolicy *TableRowAccessPolicy `ddl:"keyword"`
	Tag             []TagAssociation      `ddl:"keyword,parentheses" sql:"TAG"`
}

type ExternalTableColumn struct {
	Name             string   `ddl:"keyword"`
	Type             DataType `ddl:"keyword"`
	AsExpression     []string `ddl:"keyword,parentheses" sql:"AS"`
	NotNull          *bool    `ddl:"keyword" sql:"NOT NULL"`
	InlineConstraint *ColumnInlineConstraint
}

type CloudProviderParams struct {
	// One of
	GoogleCloudStorageIntegration *string `ddl:"parameter,single_quotes" sql:"INTEGRATION"`
	MicrosoftAzureIntegration     *string `ddl:"parameter,single_quotes" sql:"INTEGRATION"`
}

type ExternalTableFileFormat struct {
	Name    *string                      `ddl:"parameter,single_quotes" sql:"FORMAT_NAME"`
	Type    *ExternalTableFileFormatType `ddl:"parameter" sql:"TYPE"`
	Options *ExternalTableFileFormatTypeOptions
}

type ExternalTableFileFormatType string

var (
	ExternalTableFileFormatTypeCSV     ExternalTableFileFormatType = "CSV"
	ExternalTableFileFormatTypeJSON    ExternalTableFileFormatType = "JSON"
	ExternalTableFileFormatTypeAvro    ExternalTableFileFormatType = "AVRO"
	ExternalTableFileFormatTypeORC     ExternalTableFileFormatType = "ORC"
	ExternalTableFileFormatTypeParquet ExternalTableFileFormatType = "PARQUET"
)

type ExternalTableFileFormatTypeOptions struct {
	// CSV type options
	CSVCompression               *ExternalTableCsvCompression `ddl:"parameter" sql:"COMPRESSION"`
	CSVRecordDelimiter           *string                      `ddl:"parameter,single_quotes" sql:"RECORD_DELIMITER"`
	CSVFieldDelimiter            *string                      `ddl:"parameter,single_quotes" sql:"FIELD_DELIMITER"`
	CSVSkipHeader                *int                         `ddl:"parameter" sql:"SKIP_HEADER"`
	CSVSkipBlankLines            *bool                        `ddl:"parameter" sql:"SKIP_BLANK_LINES"`
	CSVEscapeUnenclosedField     *string                      `ddl:"parameter,single_quotes" sql:"ESCAPE_UNENCLOSED_FIELD"`
	CSVTrimSpace                 *bool                        `ddl:"parameter" sql:"TRIM_SPACE"`
	CSVFieldOptionallyEnclosedBy *string                      `ddl:"parameter,single_quotes" sql:"FIELD_OPTIONALLY_ENCLOSED_BY"`
	CSVNullIf                    *[]NullString                `ddl:"parameter,parentheses" sql:"NULL_IF"`
	CSVEmptyFieldAsNull          *bool                        `ddl:"parameter" sql:"EMPTY_FIELD_AS_NULL"`
	CSVEncoding                  *CSVEncoding                 `ddl:"parameter,single_quotes" sql:"ENCODING"`

	// JSON type options
	JSONCompression              *ExternalTableJsonCompression `ddl:"parameter" sql:"COMPRESSION"`
	JSONAllowDuplicate           *bool                         `ddl:"parameter" sql:"ALLOW_DUPLICATE"`
	JSONStripOuterArray          *bool                         `ddl:"parameter" sql:"STRIP_OUTER_ARRAY"`
	JSONStripNullValues          *bool                         `ddl:"parameter" sql:"STRIP_NULL_VALUES"`
	JSONReplaceInvalidCharacters *bool                         `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`

	// AVRO type options
	AvroCompression              *ExternalTableAvroCompression `ddl:"parameter" sql:"COMPRESSION"`
	AvroReplaceInvalidCharacters *bool                         `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`

	// ORC type options
	ORCTrimSpace                *bool         `ddl:"parameter" sql:"TRIM_SPACE"`
	ORCReplaceInvalidCharacters *bool         `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
	ORCNullIf                   *[]NullString `ddl:"parameter,parentheses" sql:"NULL_IF"`

	// PARQUET type options
	ParquetCompression              *ExternalTableParquetCompression `ddl:"parameter" sql:"COMPRESSION"`
	ParquetBinaryAsText             *bool                            `ddl:"parameter" sql:"BINARY_AS_TEXT"`
	ParquetReplaceInvalidCharacters *bool                            `ddl:"parameter" sql:"REPLACE_INVALID_CHARACTERS"`
}

type ExternalTableCsvCompression string

var (
	ExternalTableCsvCompressionAuto       ExternalTableCsvCompression = "AUTO"
	ExternalTableCsvCompressionGzip       ExternalTableCsvCompression = "GZIP"
	ExternalTableCsvCompressionBz2        ExternalTableCsvCompression = "BZ2"
	ExternalTableCsvCompressionBrotli     ExternalTableCsvCompression = "BROTLI"
	ExternalTableCsvCompressionZstd       ExternalTableCsvCompression = "ZSTD"
	ExternalTableCsvCompressionDeflate    ExternalTableCsvCompression = "DEFLATE"
	ExternalTableCsvCompressionRawDeflate ExternalTableCsvCompression = "RAW_DEFALTE"
	ExternalTableCsvCompressionNone       ExternalTableCsvCompression = "NONE"
)

type ExternalTableJsonCompression string

var (
	ExternalTableJsonCompressionAuto       ExternalTableJsonCompression = "AUTO"
	ExternalTableJsonCompressionGzip       ExternalTableJsonCompression = "GZIP"
	ExternalTableJsonCompressionBz2        ExternalTableJsonCompression = "BZ2"
	ExternalTableJsonCompressionBrotli     ExternalTableJsonCompression = "BROTLI"
	ExternalTableJsonCompressionZstd       ExternalTableJsonCompression = "ZSTD"
	ExternalTableJsonCompressionDeflate    ExternalTableJsonCompression = "DEFLATE"
	ExternalTableJsonCompressionRawDeflate ExternalTableJsonCompression = "RAW_DEFLATE"
	ExternalTableJsonCompressionNone       ExternalTableJsonCompression = "NONE"
)

type ExternalTableAvroCompression string

var (
	ExternalTableAvroCompressionAuto       ExternalTableAvroCompression = "AUTO"
	ExternalTableAvroCompressionGzip       ExternalTableAvroCompression = "GZIP"
	ExternalTableAvroCompressionBz2        ExternalTableAvroCompression = "BZ2"
	ExternalTableAvroCompressionBrotli     ExternalTableAvroCompression = "BROTLI"
	ExternalTableAvroCompressionZstd       ExternalTableAvroCompression = "ZSTD"
	ExternalTableAvroCompressionDeflate    ExternalTableAvroCompression = "DEFLATE"
	ExternalTableAvroCompressionRawDeflate ExternalTableAvroCompression = "RAW_DEFLATE"
	ExternalTableAvroCompressionNone       ExternalTableAvroCompression = "NONE"
)

type ExternalTableParquetCompression string

var (
	ExternalTableParquetCompressionAuto   ExternalTableParquetCompression = "AUTO"
	ExternalTableParquetCompressionSnappy ExternalTableParquetCompression = "SNAPPY"
	ExternalTableParquetCompressionNone   ExternalTableParquetCompression = "NONE"
)

// CreateWithManualPartitioningExternalTableOptions based on https://docs.snowflake.com/en/sql-reference/sql/create-external-table
type CreateWithManualPartitioningExternalTableOptions struct {
	create                     bool                   `ddl:"static" sql:"CREATE"`
	OrReplace                  *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	externalTable              bool                   `ddl:"static" sql:"EXTERNAL TABLE"`
	IfNotExists                *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier `ddl:"identifier"`
	Columns                    []ExternalTableColumn  `ddl:"list,parentheses"`
	CloudProviderParams        *CloudProviderParams
	PartitionBy                []string                  `ddl:"keyword,parentheses" sql:"PARTITION BY"`
	Location                   string                    `ddl:"parameter" sql:"LOCATION"`
	userSpecifiedPartitionType bool                      `ddl:"static" sql:"PARTITION_TYPE = USER_SPECIFIED"`
	FileFormat                 []ExternalTableFileFormat `ddl:"parameter,parentheses" sql:"FILE_FORMAT"`
	// RawFileFormat was introduced, because of the decision taken during https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2228
	// that for now the snowflake_external_table resource should continue on using raw file format, which wasn't previously supported by the new SDK.
	// In the future it should most likely be replaced by a more structured version FileFormat
	RawFileFormat   *RawFileFormat        `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	CopyGrants      *bool                 `ddl:"keyword" sql:"COPY GRANTS"`
	Comment         *string               `ddl:"parameter,single_quotes" sql:"COMMENT"`
	RowAccessPolicy *TableRowAccessPolicy `ddl:"keyword"`
	Tag             []TagAssociation      `ddl:"keyword,parentheses" sql:"TAG"`
}

// CreateDeltaLakeExternalTableOptions based on https://docs.snowflake.com/en/sql-reference/sql/create-external-table
type CreateDeltaLakeExternalTableOptions struct {
	create              bool                   `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	externalTable       bool                   `ddl:"static" sql:"EXTERNAL TABLE"`
	IfNotExists         *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
	Columns             []ExternalTableColumn  `ddl:"list,parentheses"`
	CloudProviderParams *CloudProviderParams
	PartitionBy         []string                  `ddl:"keyword,parentheses" sql:"PARTITION BY"`
	Location            string                    `ddl:"parameter" sql:"LOCATION"`
	RefreshOnCreate     *bool                     `ddl:"parameter" sql:"REFRESH_ON_CREATE"`
	AutoRefresh         *bool                     `ddl:"parameter" sql:"AUTO_REFRESH"`
	FileFormat          []ExternalTableFileFormat `ddl:"parameter,parentheses" sql:"FILE_FORMAT"`
	// RawFileFormat was introduced, because of the decision taken during https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2228
	// that for now the snowflake_external_table resource should continue on using raw file format, which wasn't previously supported by the new SDK.
	// In the future it should most likely be replaced by a more structured version FileFormat
	RawFileFormat    *RawFileFormat        `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	deltaTableFormat bool                  `ddl:"static" sql:"TABLE_FORMAT = DELTA"`
	CopyGrants       *bool                 `ddl:"keyword" sql:"COPY GRANTS"`
	Comment          *string               `ddl:"parameter,single_quotes" sql:"COMMENT"`
	RowAccessPolicy  *TableRowAccessPolicy `ddl:"keyword"`
	Tag              []TagAssociation      `ddl:"keyword,parentheses" sql:"TAG"`
}

// CreateExternalTableUsingTemplateOptions based on https://docs.snowflake.com/en/sql-reference/sql/create-external-table#variant-syntax
type CreateExternalTableUsingTemplateOptions struct {
	create              bool                   `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	externalTable       bool                   `ddl:"static" sql:"EXTERNAL TABLE"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
	CopyGrants          *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
	Query               []string               `ddl:"parameter,no_equals,parentheses" sql:"USING TEMPLATE"`
	CloudProviderParams *CloudProviderParams
	PartitionBy         []string                  `ddl:"keyword,parentheses" sql:"PARTITION BY"`
	Location            string                    `ddl:"parameter" sql:"LOCATION"`
	RefreshOnCreate     *bool                     `ddl:"parameter" sql:"REFRESH_ON_CREATE"`
	AutoRefresh         *bool                     `ddl:"parameter" sql:"AUTO_REFRESH"`
	Pattern             *string                   `ddl:"parameter,single_quotes" sql:"PATTERN"`
	FileFormat          []ExternalTableFileFormat `ddl:"parameter,parentheses" sql:"FILE_FORMAT"`
	// RawFileFormat was introduced, because of the decision taken during https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2228
	// that for now the snowflake_external_table resource should continue on using raw file format, which wasn't previously supported by the new SDK.
	// In the future it should most likely be replaced by a more structured version FileFormat
	RawFileFormat   *RawFileFormat        `ddl:"list,parentheses" sql:"FILE_FORMAT ="`
	AwsSnsTopic     *string               `ddl:"parameter,single_quotes" sql:"AWS_SNS_TOPIC"`
	Comment         *string               `ddl:"parameter,single_quotes" sql:"COMMENT"`
	RowAccessPolicy *TableRowAccessPolicy `ddl:"keyword"`
	Tag             []TagAssociation      `ddl:"keyword,parentheses" sql:"TAG"`
}

// AlterExternalTableOptions based on https://docs.snowflake.com/en/sql-reference/sql/alter-external-table
type AlterExternalTableOptions struct {
	alterExternalTable bool                   `ddl:"static" sql:"ALTER EXTERNAL TABLE"`
	IfExists           *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name               SchemaObjectIdentifier `ddl:"identifier"`
	// One of
	Refresh     *RefreshExternalTable `ddl:"keyword" sql:"REFRESH"`
	AddFiles    []ExternalTableFile   `ddl:"keyword,no_quotes,parentheses" sql:"ADD FILES"`
	RemoveFiles []ExternalTableFile   `ddl:"keyword,no_quotes,parentheses" sql:"REMOVE FILES"`
	AutoRefresh *bool                 `ddl:"parameter" sql:"SET AUTO_REFRESH"`
	SetTag      []TagAssociation      `ddl:"keyword" sql:"SET TAG"`
	UnsetTag    []ObjectIdentifier    `ddl:"keyword" sql:"UNSET TAG"`
}

type RefreshExternalTable struct {
	Path string `ddl:"parameter,no_equals,single_quotes"`
}

type ExternalTableFile struct {
	Name string `ddl:"keyword,single_quotes"`
}

// AlterExternalTablePartitionOptions based on https://docs.snowflake.com/en/sql-reference/sql/alter-external-table
type AlterExternalTablePartitionOptions struct {
	alterExternalTable bool                   `ddl:"static" sql:"ALTER EXTERNAL TABLE"`
	IfExists           *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name               SchemaObjectIdentifier `ddl:"identifier"`
	AddPartitions      []Partition            `ddl:"keyword,parentheses" sql:"ADD PARTITION"`
	DropPartition      *bool                  `ddl:"keyword" sql:"DROP PARTITION"`
	Location           string                 `ddl:"parameter,no_equals,single_quotes" sql:"LOCATION"`
}

type Partition struct {
	ColumnName string `ddl:"keyword"`
	Value      string `ddl:"parameter,single_quotes"`
}

// DropExternalTableOptions based on https://docs.snowflake.com/en/sql-reference/sql/drop-external-table
type DropExternalTableOptions struct {
	dropExternalTable bool                   `ddl:"static" sql:"DROP EXTERNAL TABLE"`
	IfExists          *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
	DropOption        *ExternalTableDropOption
}

type ExternalTableDropOption struct {
	Restrict *bool `ddl:"keyword" sql:"RESTRICT"`
	Cascade  *bool `ddl:"keyword" sql:"CASCADE"`
}

// ShowExternalTableOptions based on https://docs.snowflake.com/en/sql-reference/sql/show-external-tables
type ShowExternalTableOptions struct {
	show           bool       `ddl:"static" sql:"SHOW"`
	Terse          *bool      `ddl:"keyword" sql:"TERSE"`
	externalTables bool       `ddl:"static" sql:"EXTERNAL TABLES"`
	Like           *Like      `ddl:"keyword" sql:"LIKE"`
	In             *In        `ddl:"keyword" sql:"IN"`
	StartsWith     *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	LimitFrom      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type externalTableRow struct {
	CreatedOn           time.Time      `db:"created_on"`
	Name                string         `db:"name"`
	DatabaseName        string         `db:"database_name"`
	SchemaName          string         `db:"schema_name"`
	Invalid             bool           `db:"invalid"`
	InvalidReason       sql.NullString `db:"invalid_reason"`
	Owner               string         `db:"owner"`
	Comment             string         `db:"comment"`
	Stage               string         `db:"stage"`
	Location            string         `db:"location"`
	FileFormatName      string         `db:"file_format_name"`
	FileFormatType      sql.NullString `db:"file_format_type"`
	Cloud               string         `db:"cloud"`
	Region              string         `db:"region"`
	NotificationChannel sql.NullString `db:"notification_channel"`
	LastRefreshedOn     sql.NullTime   `db:"last_refreshed_on"`
	TableFormat         string         `db:"table_format"`
	LastRefreshDetails  sql.NullString `db:"last_refresh_details"`
	OwnerRoleType       string         `db:"owner_role_type"`
}

func (e externalTableRow) convertErr() (*ExternalTable, error) {
	et := &ExternalTable{
		CreatedOn:      e.CreatedOn,
		Name:           e.Name,
		DatabaseName:   e.DatabaseName,
		SchemaName:     e.SchemaName,
		Invalid:        e.Invalid,
		Owner:          e.Owner,
		Stage:          e.Stage,
		Location:       e.Location,
		FileFormatName: e.FileFormatName,
		Cloud:          e.Cloud,
		Region:         e.Region,
		TableFormat:    e.TableFormat,
		OwnerRoleType:  e.OwnerRoleType,
		Comment:        e.Comment,
	}
	if e.InvalidReason.Valid {
		et.InvalidReason = e.InvalidReason.String
	}
	if e.NotificationChannel.Valid {
		et.NotificationChannel = e.NotificationChannel.String
	}
	if e.LastRefreshedOn.Valid {
		et.LastRefreshedOn = e.LastRefreshedOn.Time
	}
	if e.LastRefreshDetails.Valid {
		et.LastRefreshDetails = e.LastRefreshDetails.String
	}
	if e.FileFormatType.Valid {
		et.FileFormatType = e.FileFormatType.String
	}
	return et, nil
}

// describeExternalTableColumnsOptions based on https://docs.snowflake.com/en/sql-reference/sql/desc-external-table
type describeExternalTableColumnsOptions struct {
	describeExternalTable bool                   `ddl:"static" sql:"DESCRIBE EXTERNAL TABLE"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	columnsType           bool                   `ddl:"static" sql:"TYPE = COLUMNS"`
}

type ExternalTableColumnDetails struct {
	Name       string
	Type       DataType
	Kind       string
	IsNullable bool
	Default    *string
	IsPrimary  bool
	IsUnique   bool
	Check      *bool
	Expression *string
	Comment    *string
	PolicyName *string
}

// externalTableColumnDetailsRow based on https://docs.snowflake.com/en/sql-reference/sql/desc-external-table
type externalTableColumnDetailsRow struct {
	Name       string         `db:"name"`
	Type       DataType       `db:"type"`
	Kind       string         `db:"kind"`
	IsNullable string         `db:"null?"`
	Default    sql.NullString `db:"default"`
	IsPrimary  string         `db:"primary key"`
	IsUnique   string         `db:"unique key"`
	Check      sql.NullString `db:"check"`
	Expression sql.NullString `db:"expression"`
	Comment    sql.NullString `db:"comment"`
	PolicyName sql.NullString `db:"policy name"`
}

func (r externalTableColumnDetailsRow) convertErr() (*ExternalTableColumnDetails, error) {
	details := &ExternalTableColumnDetails{
		Name:       r.Name,
		Type:       r.Type,
		Kind:       r.Kind,
		IsNullable: r.IsNullable == "Y",
		IsPrimary:  r.IsPrimary == "Y",
		IsUnique:   r.IsUnique == "Y",
	}
	if r.Default.Valid {
		details.Default = String(r.Default.String)
	}
	if r.Check.Valid {
		details.Check = Bool(r.Check.String == "Y")
	}
	if r.Expression.Valid {
		details.Expression = String(r.Expression.String)
	}
	if r.Comment.Valid {
		details.Comment = String(r.Comment.String)
	}
	if r.PolicyName.Valid {
		details.PolicyName = String(r.PolicyName.String)
	}
	return details, nil
}

type describeExternalTableStageOptions struct {
	describeExternalTable bool                   `ddl:"static" sql:"DESCRIBE EXTERNAL TABLE"`
	name                  SchemaObjectIdentifier `ddl:"identifier"`
	stageType             bool                   `ddl:"static" sql:"TYPE = STAGE"`
}

type ExternalTableStageDetails struct {
	ParentProperty  string
	Property        string
	PropertyType    string
	PropertyValue   string
	PropertyDefault string
}

type externalTableStageDetailsRow struct {
	ParentProperty  string `db:"parent_property"`
	Property        string `db:"property"`
	PropertyType    string `db:"property_type"`
	PropertyValue   string `db:"property_value"`
	PropertyDefault string `db:"property_default"`
}

func (r externalTableStageDetailsRow) convertErr() (*ExternalTableStageDetails, error) {
	return &ExternalTableStageDetails{
		ParentProperty:  r.ParentProperty,
		Property:        r.Property,
		PropertyType:    r.PropertyType,
		PropertyValue:   r.PropertyValue,
		PropertyDefault: r.PropertyDefault,
	}, nil
}
