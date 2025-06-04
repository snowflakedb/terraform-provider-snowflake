package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Services interface {
	Create(ctx context.Context, request *CreateServiceRequest) error
	Alter(ctx context.Context, request *AlterServiceRequest) error
	ExecuteJobService(ctx context.Context, request *ExecuteJobServiceServiceRequest) error
	Drop(ctx context.Context, request *DropServiceRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, request *ShowServiceRequest) ([]Service, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Service, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Service, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*ServiceDetails, error)
}

// CreateServiceOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-service.
type CreateServiceOptions struct {
	create                     bool                               `ddl:"static" sql:"CREATE"`
	service                    bool                               `ddl:"static" sql:"SERVICE"`
	IfNotExists                *bool                              `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                       SchemaObjectIdentifier             `ddl:"identifier"`
	InComputePool              AccountObjectIdentifier            `ddl:"identifier" sql:"IN COMPUTE POOL"`
	FromSpecification          *ServiceFromSpecification          `ddl:"keyword"`
	FromSpecificationTemplate  *ServiceFromSpecificationTemplate  `ddl:"keyword"`
	AutoSuspendSecs            *int                               `ddl:"parameter" sql:"AUTO_SUSPEND_SECS"`
	ExternalAccessIntegrations *ServiceExternalAccessIntegrations `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	AutoResume                 *bool                              `ddl:"parameter" sql:"AUTO_RESUME"`
	MinInstances               *int                               `ddl:"parameter" sql:"MIN_INSTANCES"`
	MinReadyInstances          *int                               `ddl:"parameter" sql:"MIN_READY_INSTANCES"`
	MaxInstances               *int                               `ddl:"parameter" sql:"MAX_INSTANCES"`
	QueryWarehouse             *AccountObjectIdentifier           `ddl:"identifier,equals" sql:"QUERY_WAREHOUSE"`
	Tag                        []TagAssociation                   `ddl:"keyword,parentheses" sql:"TAG"`
	Comment                    *string                            `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ServiceExternalAccessIntegrations struct {
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"list,must_parentheses"`
}

type ListItem struct {
	Key         string `ddl:"keyword,double_quotes"`
	arrowEquals bool   `ddl:"static" sql:"=>"`
	Value       any    `ddl:"keyword"`
}

type ServiceFromSpecification struct {
	from              bool     `ddl:"static" sql:"FROM"`
	Location          Location `ddl:"parameter,no_quotes,no_equals"`
	SpecificationFile *string  `ddl:"parameter,single_quotes" sql:"SPECIFICATION_FILE"`
	Specification     *string  `ddl:"parameter,double_dollar_quotes,no_equals" sql:"SPECIFICATION"`
}

type ServiceFromSpecificationTemplate struct {
	from                      bool       `ddl:"static" sql:"FROM"`
	Location                  Location   `ddl:"parameter,no_quotes,no_equals"`
	SpecificationTemplateFile *string    `ddl:"parameter,single_quotes" sql:"SPECIFICATION_TEMPLATE_FILE"`
	SpecificationTemplate     *string    `ddl:"parameter,double_dollar_quotes,no_equals" sql:"SPECIFICATION_TEMPLATE"`
	Using                     []ListItem `ddl:"parameter,parentheses,no_equals" sql:"USING"`
}

// AlterServiceOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-service.
type AlterServiceOptions struct {
	alter                     bool                              `ddl:"static" sql:"ALTER"`
	service                   bool                              `ddl:"static" sql:"SERVICE"`
	IfExists                  *bool                             `ddl:"keyword" sql:"IF EXISTS"`
	name                      SchemaObjectIdentifier            `ddl:"identifier"`
	Resume                    *bool                             `ddl:"keyword" sql:"RESUME"`
	Suspend                   *bool                             `ddl:"keyword" sql:"SUSPEND"`
	FromSpecification         *ServiceFromSpecification         `ddl:"keyword"`
	FromSpecificationTemplate *ServiceFromSpecificationTemplate `ddl:"keyword"`
	Restore                   *Restore                          `ddl:"keyword" sql:"RESTORE"`
	Set                       *ServiceSet                       `ddl:"keyword" sql:"SET"`
	Unset                     *ServiceUnset                     `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTags                   []TagAssociation                  `ddl:"keyword" sql:"SET TAG"`
	UnsetTags                 []ObjectIdentifier                `ddl:"keyword" sql:"UNSET TAG"`
}

type Restore struct {
	Volume       string                 `ddl:"parameter,double_quotes,no_equals" sql:"VOLUME"`
	Instances    []int                  `ddl:"keyword" sql:"INSTANCES"`
	FromSnapshot SchemaObjectIdentifier `ddl:"identifier" sql:"FROM SNAPSHOT"`
}

type ServiceSet struct {
	MinInstances               *int                               `ddl:"parameter" sql:"MIN_INSTANCES"`
	MaxInstances               *int                               `ddl:"parameter" sql:"MAX_INSTANCES"`
	AutoSuspendSecs            *int                               `ddl:"parameter" sql:"AUTO_SUSPEND_SECS"`
	MinReadyInstances          *int                               `ddl:"parameter" sql:"MIN_READY_INSTANCES"`
	QueryWarehouse             *AccountObjectIdentifier           `ddl:"identifier,equals" sql:"QUERY_WAREHOUSE"`
	AutoResume                 *bool                              `ddl:"parameter" sql:"AUTO_RESUME"`
	ExternalAccessIntegrations *ServiceExternalAccessIntegrations `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Comment                    *string                            `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ServiceUnset struct {
	MinInstances               *bool `ddl:"keyword" sql:"MIN_INSTANCES"`
	AutoSuspendSecs            *bool `ddl:"keyword" sql:"AUTO_SUSPEND_SECS"`
	MaxInstances               *bool `ddl:"keyword" sql:"MAX_INSTANCES"`
	MinReadyInstances          *bool `ddl:"keyword" sql:"MIN_READY_INSTANCES"`
	QueryWarehouse             *bool `ddl:"keyword" sql:"QUERY_WAREHOUSE"`
	AutoResume                 *bool `ddl:"keyword" sql:"AUTO_RESUME"`
	ExternalAccessIntegrations *bool `ddl:"keyword" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	Comment                    *bool `ddl:"keyword" sql:"COMMENT"`
}

// ExecuteJobServiceServiceOptions is based on https://docs.snowflake.com/en/sql-reference/sql/execute-job-service.
type ExecuteJobServiceServiceOptions struct {
	executeJobService                   bool                                 `ddl:"static" sql:"EXECUTE JOB SERVICE"`
	InComputePool                       AccountObjectIdentifier              `ddl:"identifier" sql:"IN COMPUTE POOL"`
	Name                                SchemaObjectIdentifier               `ddl:"identifier,equals" sql:"NAME"`
	Async                               *bool                                `ddl:"parameter" sql:"ASYNC"`
	QueryWarehouse                      *AccountObjectIdentifier             `ddl:"identifier,equals" sql:"QUERY_WAREHOUSE"`
	Comment                             *string                              `ddl:"parameter,single_quotes" sql:"COMMENT"`
	ExternalAccessIntegrations          *ServiceExternalAccessIntegrations   `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	JobServiceFromSpecification         *JobServiceFromSpecification         `ddl:"keyword"`
	JobServiceFromSpecificationTemplate *JobServiceFromSpecificationTemplate `ddl:"keyword"`
	Tag                                 []TagAssociation                     `ddl:"keyword,parentheses" sql:"TAG"`
}

type JobServiceFromSpecification struct {
	from              bool     `ddl:"static" sql:"FROM"`
	Location          Location `ddl:"parameter,no_quotes,no_equals"`
	SpecificationFile *string  `ddl:"parameter,single_quotes" sql:"SPECIFICATION_FILE"`
	Specification     *string  `ddl:"parameter,double_dollar_quotes,no_equals" sql:"SPECIFICATION"`
}

type JobServiceFromSpecificationTemplate struct {
	from                      bool       `ddl:"static" sql:"FROM"`
	Location                  Location   `ddl:"parameter,no_quotes,no_equals"`
	SpecificationTemplateFile *string    `ddl:"parameter,single_quotes" sql:"SPECIFICATION_TEMPLATE_FILE"`
	SpecificationTemplate     *string    `ddl:"parameter,double_dollar_quotes,no_equals" sql:"SPECIFICATION_TEMPLATE"`
	Using                     []ListItem `ddl:"parameter,parentheses,no_equals" sql:"USING"`
}

// DropServiceOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-service.
type DropServiceOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	service  bool                   `ddl:"static" sql:"SERVICE"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
	Force    *bool                  `ddl:"keyword" sql:"FORCE"`
}

// ShowServiceOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-services.
type ShowServiceOptions struct {
	show        bool       `ddl:"static" sql:"SHOW"`
	Job         *bool      `ddl:"keyword" sql:"JOB"`
	services    bool       `ddl:"static" sql:"SERVICES"`
	ExcludeJobs *bool      `ddl:"keyword" sql:"EXCLUDE JOBS"`
	Like        *Like      `ddl:"keyword" sql:"LIKE"`
	In          *ServiceIn `ddl:"keyword" sql:"IN"`
	StartsWith  *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit       *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type servicesRow struct {
	Name                       string         `db:"name"`
	Status                     string         `db:"status"`
	DatabaseName               string         `db:"database_name"`
	SchemaName                 string         `db:"schema_name"`
	Owner                      string         `db:"owner"`
	ComputePool                string         `db:"compute_pool"`
	DnsName                    string         `db:"dns_name"`
	CurrentInstances           int            `db:"current_instances"`
	TargetInstances            int            `db:"target_instances"`
	MinReadyInstances          int            `db:"min_ready_instances"`
	MinInstances               int            `db:"min_instances"`
	MaxInstances               int            `db:"max_instances"`
	AutoResume                 bool           `db:"auto_resume"`
	ExternalAccessIntegrations sql.NullString `db:"external_access_integrations"`
	CreatedOn                  time.Time      `db:"created_on"`
	UpdatedOn                  time.Time      `db:"updated_on"`
	ResumedOn                  sql.NullTime   `db:"resumed_on"`
	SuspendedOn                sql.NullTime   `db:"suspended_on"`
	AutoSuspendSecs            int            `db:"auto_suspend_secs"`
	Comment                    sql.NullString `db:"comment"`
	OwnerRoleType              string         `db:"owner_role_type"`
	QueryWarehouse             sql.NullString `db:"query_warehouse"`
	IsJob                      bool           `db:"is_job"`
	IsAsyncJob                 bool           `db:"is_async_job"`
	SpecDigest                 string         `db:"spec_digest"`
	IsUpgrading                bool           `db:"is_upgrading"`
	ManagingObjectDomain       sql.NullString `db:"managing_object_domain"`
	ManagingObjectName         sql.NullString `db:"managing_object_name"`
}

type Service struct {
	Name                       string
	Status                     ServiceStatus
	DatabaseName               string
	SchemaName                 string
	Owner                      string
	ComputePool                AccountObjectIdentifier
	DnsName                    string
	CurrentInstances           int
	TargetInstances            int
	MinReadyInstances          int
	MinInstances               int
	MaxInstances               int
	AutoResume                 bool
	ExternalAccessIntegrations []AccountObjectIdentifier
	CreatedOn                  time.Time
	UpdatedOn                  time.Time
	ResumedOn                  *time.Time
	SuspendedOn                *time.Time
	AutoSuspendSecs            int
	Comment                    *string
	OwnerRoleType              string
	QueryWarehouse             *AccountObjectIdentifier
	IsJob                      bool
	IsAsyncJob                 bool
	SpecDigest                 string
	IsUpgrading                bool
	ManagingObjectDomain       *string
	ManagingObjectName         *string
}

func (v *Service) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Service) ObjectType() ObjectType {
	return ObjectTypeService
}

// DescribeServiceOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-service.
type DescribeServiceOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	service  bool                   `ddl:"static" sql:"SERVICE"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

type serviceDescRow struct {
	Name                       string         `db:"name"`
	Status                     string         `db:"status"`
	DatabaseName               string         `db:"database_name"`
	SchemaName                 string         `db:"schema_name"`
	Owner                      string         `db:"owner"`
	ComputePool                string         `db:"compute_pool"`
	Spec                       string         `db:"spec"`
	DnsName                    string         `db:"dns_name"`
	CurrentInstances           int            `db:"current_instances"`
	TargetInstances            int            `db:"target_instances"`
	MinReadyInstances          int            `db:"min_ready_instances"`
	MinInstances               int            `db:"min_instances"`
	MaxInstances               int            `db:"max_instances"`
	AutoResume                 bool           `db:"auto_resume"`
	ExternalAccessIntegrations sql.NullString `db:"external_access_integrations"`
	CreatedOn                  time.Time      `db:"created_on"`
	UpdatedOn                  time.Time      `db:"updated_on"`
	ResumedOn                  sql.NullTime   `db:"resumed_on"`
	SuspendedOn                sql.NullTime   `db:"suspended_on"`
	AutoSuspendSecs            int            `db:"auto_suspend_secs"`
	Comment                    sql.NullString `db:"comment"`
	OwnerRoleType              string         `db:"owner_role_type"`
	QueryWarehouse             sql.NullString `db:"query_warehouse"`
	IsJob                      bool           `db:"is_job"`
	IsAsyncJob                 bool           `db:"is_async_job"`
	SpecDigest                 string         `db:"spec_digest"`
	IsUpgrading                bool           `db:"is_upgrading"`
	ManagingObjectDomain       sql.NullString `db:"managing_object_domain"`
	ManagingObjectName         sql.NullString `db:"managing_object_name"`
}

type ServiceDetails struct {
	Name                       string
	Status                     ServiceStatus
	DatabaseName               string
	SchemaName                 string
	Owner                      string
	ComputePool                AccountObjectIdentifier
	Spec                       string
	DnsName                    string
	CurrentInstances           int
	TargetInstances            int
	MinReadyInstances          int
	MinInstances               int
	MaxInstances               int
	AutoResume                 bool
	ExternalAccessIntegrations []AccountObjectIdentifier
	CreatedOn                  time.Time
	UpdatedOn                  time.Time
	ResumedOn                  *time.Time
	SuspendedOn                *time.Time
	AutoSuspendSecs            int
	Comment                    *string
	OwnerRoleType              string
	QueryWarehouse             *AccountObjectIdentifier
	IsJob                      bool
	IsAsyncJob                 bool
	SpecDigest                 string
	IsUpgrading                bool
	ManagingObjectDomain       *string
	ManagingObjectName         *string
}
