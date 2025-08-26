package sdk

import (
	"context"
	"database/sql"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
)

type DynamicTables interface {
	Create(ctx context.Context, request *CreateDynamicTableRequest) error
	Alter(ctx context.Context, request *AlterDynamicTableRequest) error
	Describe(ctx context.Context, request *DescribeDynamicTableRequest) (*DynamicTableDetails, error)
	Drop(ctx context.Context, request *DropDynamicTableRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, request *ShowDynamicTableRequest) ([]DynamicTable, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*DynamicTable, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*DynamicTable, error)
}

// createDynamicTableOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-dynamic-table
type createDynamicTableOptions struct {
	create       bool                     `ddl:"static" sql:"CREATE"`
	OrReplace    *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	dynamicTable bool                     `ddl:"static" sql:"DYNAMIC TABLE"`
	name         SchemaObjectIdentifier   `ddl:"identifier"`
	targetLag    TargetLag                `ddl:"parameter,no_quotes" sql:"TARGET_LAG"`
	Initialize   *DynamicTableInitialize  `ddl:"parameter,no_quotes" sql:"INITIALIZE"`
	RefreshMode  *DynamicTableRefreshMode `ddl:"parameter,no_quotes" sql:"REFRESH_MODE"`
	warehouse    AccountObjectIdentifier  `ddl:"identifier,equals" sql:"WAREHOUSE"`
	Comment      *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
	query        string                   `ddl:"parameter,no_equals,no_quotes" sql:"AS"`
}

type TargetLag struct {
	MaximumDuration *string `ddl:"keyword,single_quotes"`
	Downstream      *bool   `ddl:"keyword" sql:"DOWNSTREAM"`
}

type DynamicTableSet struct {
	TargetLag *TargetLag               `ddl:"parameter,no_quotes" sql:"TARGET_LAG"`
	Warehouse *AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
}

// alterDynamicTableOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-dynamic-table
type alterDynamicTableOptions struct {
	alter        bool                   `ddl:"static" sql:"ALTER"`
	dynamicTable bool                   `ddl:"static" sql:"DYNAMIC TABLE"`
	name         SchemaObjectIdentifier `ddl:"identifier"`

	Suspend *bool            `ddl:"keyword" sql:"SUSPEND"`
	Resume  *bool            `ddl:"keyword" sql:"RESUME"`
	Refresh *bool            `ddl:"keyword" sql:"REFRESH"`
	Set     *DynamicTableSet `ddl:"keyword" sql:"SET"`
}

// dropDynamicTableOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-dynamic-table
type dropDynamicTableOptions struct {
	drop         bool                   `ddl:"static" sql:"DROP"`
	dynamicTable bool                   `ddl:"static" sql:"DYNAMIC TABLE"`
	IfExists     *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
}

// showDynamicTableOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-dynamic-tables
type showDynamicTableOptions struct {
	show         bool       `ddl:"static" sql:"SHOW"`
	dynamicTable bool       `ddl:"static" sql:"DYNAMIC TABLES"`
	Like         *Like      `ddl:"keyword" sql:"LIKE"`
	In           *In        `ddl:"keyword" sql:"IN"`
	StartsWith   *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit        *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type DynamicTableRefreshMode string

const (
	DynamicTableRefreshModeAuto        DynamicTableRefreshMode = "AUTO"
	DynamicTableRefreshModeIncremental DynamicTableRefreshMode = "INCREMENTAL"
	DynamicTableRefreshModeFull        DynamicTableRefreshMode = "FULL"
)

func (d DynamicTableRefreshMode) ToPointer() *DynamicTableRefreshMode {
	return &d
}

var AllDynamicRefreshModes = []DynamicTableRefreshMode{DynamicTableRefreshModeAuto, DynamicTableRefreshModeIncremental, DynamicTableRefreshModeFull}

type DynamicTableInitialize string

const (
	DynamicTableInitializeOnCreate   DynamicTableInitialize = "ON_CREATE"
	DynamicTableInitializeOnSchedule DynamicTableInitialize = "ON_SCHEDULE"
)

func (d DynamicTableInitialize) ToPointer() *DynamicTableInitialize {
	return &d
}

var AllDynamicTableInitializes = []DynamicTableInitialize{DynamicTableInitializeOnCreate, DynamicTableInitializeOnSchedule}

type DynamicTableSchedulingState string

const (
	DynamicTableSchedulingStateActive    DynamicTableSchedulingState = "ACTIVE"
	DynamicTableSchedulingStateSuspended DynamicTableSchedulingState = "SUSPENDED"
)

type DynamicTable struct {
	CreatedOn           time.Time
	Name                string
	Reserved            string
	DatabaseName        string
	SchemaName          string
	ClusterBy           string
	Rows                int
	Bytes               int
	Owner               string
	TargetLag           string
	RefreshMode         DynamicTableRefreshMode
	RefreshModeReason   string
	Warehouse           string
	Comment             string
	Text                string
	AutomaticClustering bool
	SchedulingState     DynamicTableSchedulingState
	LastSuspendedOn     time.Time
	IsClone             bool
	IsReplica           bool
	DataTimestamp       time.Time
	OwnerRoleType       string
}

func (dt *DynamicTable) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(dt.DatabaseName, dt.SchemaName, dt.Name)
}

type dynamicTableRow struct {
	CreatedOn           time.Time      `db:"created_on"`
	Name                string         `db:"name"`
	Reserved            string         `db:"reserved"`
	DatabaseName        string         `db:"database_name"`
	SchemaName          string         `db:"schema_name"`
	ClusterBy           string         `db:"cluster_by"`
	Rows                int            `db:"rows"`
	Bytes               int            `db:"bytes"`
	Owner               string         `db:"owner"`
	TargetLag           string         `db:"target_lag"`
	RefreshMode         string         `db:"refresh_mode"`
	RefreshModeReason   sql.NullString `db:"refresh_mode_reason"`
	Warehouse           string         `db:"warehouse"`
	Comment             string         `db:"comment"`
	Text                string         `db:"text"`
	AutomaticClustering string         `db:"automatic_clustering"`
	SchedulingState     string         `db:"scheduling_state"`
	LastSuspendedOn     sql.NullTime   `db:"last_suspended_on"`
	IsClone             bool           `db:"is_clone"`
	IsReplica           bool           `db:"is_replica"`
	DataTimestamp       sql.NullTime   `db:"data_timestamp"`
	OwnerRoleType       sql.NullString `db:"owner_role_type"`
}

func (dtr dynamicTableRow) convert() (*DynamicTable, error) {
	dt := &DynamicTable{
		CreatedOn:           dtr.CreatedOn,
		Name:                dtr.Name,
		Reserved:            dtr.Reserved,
		DatabaseName:        dtr.DatabaseName,
		SchemaName:          dtr.SchemaName,
		ClusterBy:           dtr.ClusterBy,
		Rows:                dtr.Rows,
		Bytes:               dtr.Bytes,
		Owner:               dtr.Owner,
		TargetLag:           dtr.TargetLag,
		RefreshMode:         DynamicTableRefreshMode(dtr.RefreshMode),
		Warehouse:           dtr.Warehouse,
		Comment:             dtr.Comment,
		Text:                tracking.TrimMetadata(dtr.Text),
		AutomaticClustering: dtr.AutomaticClustering == "ON", // "ON" or "OFF
		SchedulingState:     DynamicTableSchedulingState(dtr.SchedulingState),
		IsClone:             dtr.IsClone,
		IsReplica:           dtr.IsReplica,
	}
	if dtr.RefreshModeReason.Valid {
		dt.RefreshModeReason = dtr.RefreshModeReason.String
	}
	if dtr.DataTimestamp.Valid {
		dt.DataTimestamp = dtr.DataTimestamp.Time
	}
	if dtr.LastSuspendedOn.Valid {
		dt.LastSuspendedOn = dtr.LastSuspendedOn.Time
	}
	if dtr.OwnerRoleType.Valid {
		dt.OwnerRoleType = dtr.OwnerRoleType.String
	}

	return dt, nil
}

// describeDynamicTableOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-dynamic-table
type describeDynamicTableOptions struct {
	describe     bool                   `ddl:"static" sql:"DESCRIBE"`
	dynamicTable bool                   `ddl:"static" sql:"DYNAMIC TABLE"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
}

type DynamicTableDetails struct {
	Name       string
	Type       DataType
	Kind       string
	IsNull     bool
	Default    string
	PrimaryKey string
	UniqueKey  string
	Check      string
	Expression string
	Comment    string
	PolicyName string
}

type dynamicTableDetailsRow struct {
	Name       string         `db:"name"`
	Type       string         `db:"type"`
	Kind       string         `db:"kind"`
	IsNull     string         `db:"null?"`
	Default    sql.NullString `db:"default"`
	PrimaryKey string         `db:"primary key"`
	UniqueKey  string         `db:"unique key"`
	Check      sql.NullString `db:"check"`
	Expression sql.NullString `db:"expression"`
	Comment    sql.NullString `db:"comment"`
	PolicyName sql.NullString `db:"policy name"`
}

func (row dynamicTableDetailsRow) convert() (*DynamicTableDetails, error) {
	typ, _ := datatypes.ParseDataType(row.Type)
	dtd := &DynamicTableDetails{
		Name:       row.Name,
		Type:       LegacyDataTypeFrom(typ),
		Kind:       row.Kind,
		IsNull:     row.IsNull == "Y",
		PrimaryKey: row.PrimaryKey,
		UniqueKey:  row.UniqueKey,
	}
	if row.Default.Valid {
		dtd.Default = row.Default.String
	}
	if row.Check.Valid {
		dtd.Check = row.Check.String
	}
	if row.Expression.Valid {
		dtd.Expression = row.Expression.String
	}
	if row.Comment.Valid {
		dtd.Comment = row.Comment.String
	}
	if row.PolicyName.Valid {
		dtd.PolicyName = row.PolicyName.String
	}
	return dtd, nil
}
