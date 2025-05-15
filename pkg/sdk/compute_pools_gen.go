package sdk

import (
	"context"
	"database/sql"
	"time"
)

type ComputePools interface {
	Create(ctx context.Context, request *CreateComputePoolRequest) error
	Alter(ctx context.Context, request *AlterComputePoolRequest) error
	Drop(ctx context.Context, request *DropComputePoolRequest) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, request *ShowComputePoolRequest) ([]ComputePool, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ComputePool, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*ComputePool, error)
	Describe(ctx context.Context, id AccountObjectIdentifier) (*ComputePoolDetails, error)
}

// CreateComputePoolOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool.
type CreateComputePoolOptions struct {
	create             bool                      `ddl:"static" sql:"CREATE"`
	computePool        bool                      `ddl:"static" sql:"COMPUTE POOL"`
	IfNotExists        *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name               AccountObjectIdentifier   `ddl:"identifier"`
	ForApplication     *AccountObjectIdentifier  `ddl:"identifier" sql:"FOR APPLICATION"`
	MinNodes           int                       `ddl:"parameter" sql:"MIN_NODES"`
	MaxNodes           int                       `ddl:"parameter" sql:"MAX_NODES"`
	InstanceFamily     ComputePoolInstanceFamily `ddl:"parameter,no_quotes" sql:"INSTANCE_FAMILY"`
	AutoResume         *bool                     `ddl:"parameter" sql:"AUTO_RESUME"`
	InitiallySuspended *bool                     `ddl:"parameter" sql:"INITIALLY_SUSPENDED"`
	AutoSuspendSecs    *int                      `ddl:"parameter" sql:"AUTO_SUSPEND_SECS"`
	Tag                []TagAssociation          `ddl:"keyword,parentheses" sql:"TAG"`
	Comment            *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterComputePoolOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-compute-pool.
type AlterComputePoolOptions struct {
	alter       bool                    `ddl:"static" sql:"ALTER"`
	computePool bool                    `ddl:"static" sql:"COMPUTE POOL"`
	IfExists    *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name        AccountObjectIdentifier `ddl:"identifier"`
	Set         *ComputePoolSet         `ddl:"keyword" sql:"SET"`
	Unset       *ComputePoolUnset       `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTags     []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags   []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
}

type ComputePoolSet struct {
	MinNodes        *int    `ddl:"parameter" sql:"MIN_NODES"`
	MaxNodes        *int    `ddl:"parameter" sql:"MAX_NODES"`
	AutoResume      *bool   `ddl:"parameter" sql:"AUTO_RESUME"`
	AutoSuspendSecs *int    `ddl:"parameter" sql:"AUTO_SUSPEND_SECS"`
	Comment         *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ComputePoolUnset struct {
	AutoResume      *bool `ddl:"keyword" sql:"AUTO_RESUME"`
	AutoSuspendSecs *bool `ddl:"keyword" sql:"AUTO_SUSPEND_SECS"`
	Comment         *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropComputePoolOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-compute-pool.
type DropComputePoolOptions struct {
	drop        bool                    `ddl:"static" sql:"DROP"`
	computePool bool                    `ddl:"static" sql:"COMPUTE POOL"`
	IfExists    *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name        AccountObjectIdentifier `ddl:"identifier"`
}

// ShowComputePoolOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-compute-pools.
type ShowComputePoolOptions struct {
	show         bool       `ddl:"static" sql:"SHOW"`
	computePools bool       `ddl:"static" sql:"COMPUTE POOLS"`
	Like         *Like      `ddl:"keyword" sql:"LIKE"`
	StartsWith   *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit        *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type computePoolsRow struct {
	Name            string         `db:"name"`
	State           string         `db:"state"`
	MinNodes        int            `db:"min_nodes"`
	MaxNodes        int            `db:"max_nodes"`
	InstanceFamily  string         `db:"instance_family"`
	NumServices     int            `db:"num_services"`
	NumJobs         int            `db:"num_jobs"`
	AutoSuspendSecs int            `db:"auto_suspend_secs"`
	AutoResume      bool           `db:"auto_resume"`
	ActiveNodes     int            `db:"active_nodes"`
	IdleNodes       int            `db:"idle_nodes"`
	TargetNodes     int            `db:"target_nodes"`
	CreatedOn       time.Time      `db:"created_on"`
	ResumedOn       time.Time      `db:"resumed_on"`
	UpdatedOn       time.Time      `db:"updated_on"`
	Owner           string         `db:"owner"`
	Comment         sql.NullString `db:"comment"`
	IsExclusive     bool           `db:"is_exclusive"`
	Application     sql.NullString `db:"application"`
	Budget          sql.NullString `db:"budget"`
}

type ComputePool struct {
	Name            string
	State           ComputePoolState
	MinNodes        int
	MaxNodes        int
	InstanceFamily  ComputePoolInstanceFamily
	NumServices     int
	NumJobs         int
	AutoSuspendSecs int
	AutoResume      bool
	ActiveNodes     int
	IdleNodes       int
	TargetNodes     int
	CreatedOn       time.Time
	ResumedOn       time.Time
	UpdatedOn       time.Time
	Owner           string
	Comment         *string
	IsExclusive     bool
	Application     *string
	Budget          *string
}

func (v *ComputePool) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *ComputePool) ObjectType() ObjectType {
	return ObjectTypeComputePool
}

// DescribeComputePoolOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-compute-pool.
type DescribeComputePoolOptions struct {
	describe    bool                    `ddl:"static" sql:"DESCRIBE"`
	computePool bool                    `ddl:"static" sql:"COMPUTE POOL"`
	name        AccountObjectIdentifier `ddl:"identifier"`
}

type computePoolDescRow struct {
	Name            string         `db:"name"`
	State           string         `db:"state"`
	MinNodes        int            `db:"min_nodes"`
	MaxNodes        int            `db:"max_nodes"`
	InstanceFamily  string         `db:"instance_family"`
	NumServices     int            `db:"num_services"`
	NumJobs         int            `db:"num_jobs"`
	AutoSuspendSecs int            `db:"auto_suspend_secs"`
	AutoResume      bool           `db:"auto_resume"`
	ActiveNodes     int            `db:"active_nodes"`
	IdleNodes       int            `db:"idle_nodes"`
	TargetNodes     int            `db:"target_nodes"`
	CreatedOn       time.Time      `db:"created_on"`
	ResumedOn       time.Time      `db:"resumed_on"`
	UpdatedOn       time.Time      `db:"updated_on"`
	Owner           string         `db:"owner"`
	Comment         sql.NullString `db:"comment"`
	IsExclusive     bool           `db:"is_exclusive"`
	Application     sql.NullString `db:"application"`
	Budget          sql.NullString `db:"budget"`
	ErrorCode       string         `db:"error_code"`
	StatusMessage   string         `db:"status_message"`
}

type ComputePoolDetails struct {
	Name            string
	State           ComputePoolState
	MinNodes        int
	MaxNodes        int
	InstanceFamily  ComputePoolInstanceFamily
	NumServices     int
	NumJobs         int
	AutoSuspendSecs int
	AutoResume      bool
	ActiveNodes     int
	IdleNodes       int
	TargetNodes     int
	CreatedOn       time.Time
	ResumedOn       time.Time
	UpdatedOn       time.Time
	Owner           string
	Comment         *string
	IsExclusive     bool
	Application     *AccountObjectIdentifier
	Budget          *string
	ErrorCode       string
	StatusMessage   string
}
