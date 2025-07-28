package sdk

import (
	"context"
	"database/sql"
	"time"
)

type CortexSearchServices interface {
	Create(ctx context.Context, request *CreateCortexSearchServiceRequest) error
	Alter(ctx context.Context, request *AlterCortexSearchServiceRequest) error
	Show(ctx context.Context, request *ShowCortexSearchServiceRequest) ([]CortexSearchService, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*CortexSearchService, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*CortexSearchService, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*CortexSearchServiceDetails, error)
	Drop(ctx context.Context, request *DropCortexSearchServiceRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
}

// CreateCortexSearchServiceOptions is based on https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/create-cortex-search.
type CreateCortexSearchServiceOptions struct {
	create              bool                    `ddl:"static" sql:"CREATE"`
	OrReplace           *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	cortexSearchService bool                    `ddl:"static" sql:"CORTEX SEARCH SERVICE"`
	IfNotExists         *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                SchemaObjectIdentifier  `ddl:"identifier"`
	On                  string                  `ddl:"parameter,no_quotes,no_equals" sql:"ON"`
	Attributes          *Attributes             `ddl:"keyword"`
	Warehouse           AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
	TargetLag           string                  `ddl:"parameter,single_quotes" sql:"TARGET_LAG"`
	EmbeddingModel      *string                 `ddl:"parameter,single_quotes" sql:"EMBEDDING_MODEL"`
	Comment             *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
	QueryDefinition     string                  `ddl:"parameter,no_quotes,no_equals" sql:"AS"`
}
type Attributes struct {
	attributes bool     `ddl:"static" sql:"ATTRIBUTES"`
	Columns    []string `ddl:"list,no_parentheses,no_equals"`
}

// AlterCortexSearchServiceOptions is based on https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/alter-cortex-search.
type AlterCortexSearchServiceOptions struct {
	alter               bool                    `ddl:"static" sql:"ALTER"`
	cortexSearchService bool                    `ddl:"static" sql:"CORTEX SEARCH SERVICE"`
	IfExists            *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name                SchemaObjectIdentifier  `ddl:"identifier"`
	Set                 *CortexSearchServiceSet `ddl:"keyword" sql:"SET"`
}
type CortexSearchServiceSet struct {
	TargetLag *string                  `ddl:"parameter,single_quotes" sql:"TARGET_LAG"`
	Warehouse *AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
	Comment   *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// ShowCortexSearchServiceOptions is based on https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/show-cortex-search.
type ShowCortexSearchServiceOptions struct {
	show                 bool       `ddl:"static" sql:"SHOW"`
	cortexSearchServices bool       `ddl:"static" sql:"CORTEX SEARCH SERVICES"`
	Like                 *Like      `ddl:"keyword" sql:"LIKE"`
	In                   *In        `ddl:"keyword" sql:"IN"`
	StartsWith           *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit                *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}
type cortexSearchServiceRow struct {
	CreatedOn    time.Time      `db:"created_on"`
	Name         string         `db:"name"`
	DatabaseName string         `db:"database_name"`
	SchemaName   string         `db:"schema_name"`
	Comment      sql.NullString `db:"comment"`
}
type CortexSearchService struct {
	CreatedOn    time.Time
	Name         string
	DatabaseName string
	SchemaName   string
	Comment      string
}

func (dt *CortexSearchService) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(dt.DatabaseName, dt.SchemaName, dt.Name)
}

// DescribeCortexSearchServiceOptions is based on https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/desc-cortex-search.
type DescribeCortexSearchServiceOptions struct {
	describe            bool                   `ddl:"static" sql:"DESCRIBE"`
	cortexSearchService bool                   `ddl:"static" sql:"CORTEX SEARCH SERVICE"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
}
type cortexSearchServiceDetailsRow struct {
	CreatedOn         string         `db:"created_on"`
	Name              string         `db:"name"`
	DatabaseName      string         `db:"database_name"`
	SchemaName        string         `db:"schema_name"`
	TargetLag         string         `db:"target_lag"`
	Warehouse         string         `db:"warehouse"`
	SearchColumn      sql.NullString `db:"search_column"`
	AttributeColumns  sql.NullString `db:"attribute_columns"`
	Columns           sql.NullString `db:"columns"`
	Definition        sql.NullString `db:"definition"`
	Comment           sql.NullString `db:"comment"`
	ServiceQueryUrl   string         `db:"service_query_url"`
	DataTimestamp     string         `db:"data_timestamp"`
	SourceDataNumRows int            `db:"source_data_num_rows"`
	IndexingState     string         `db:"indexing_state"`
	IndexingError     sql.NullString `db:"indexing_error"`
	EmbeddingModel    sql.NullString `db:"embedding_model"`
}
type CortexSearchServiceDetails struct {
	CreatedOn         string
	Name              string
	DatabaseName      string
	SchemaName        string
	TargetLag         string
	Warehouse         string
	SearchColumn      *string
	AttributeColumns  []string
	Columns           []string
	Definition        *string
	Comment           *string
	ServiceQueryUrl   string
	DataTimestamp     string
	SourceDataNumRows int
	IndexingState     string
	IndexingError     *string
	EmbeddingModel    *string
}

// DropCortexSearchServiceOptions is based on https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/drop-cortex-search.
type DropCortexSearchServiceOptions struct {
	drop                bool                   `ddl:"static" sql:"DROP"`
	cortexSearchService bool                   `ddl:"static" sql:"CORTEX SEARCH SERVICE"`
	IfExists            *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name                SchemaObjectIdentifier `ddl:"identifier"`
}
