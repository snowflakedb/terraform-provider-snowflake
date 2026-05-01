package sdk

import (
	"context"
	"database/sql"
	"time"
)

type OpenflowRuntimes interface {
	Create(ctx context.Context, request *CreateOpenflowRuntimeRequest) error
	Alter(ctx context.Context, request *AlterOpenflowRuntimeRequest) error
	Drop(ctx context.Context, request *DropOpenflowRuntimeRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, request *ShowOpenflowRuntimeRequest) ([]OpenflowRuntime, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowRuntime, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowRuntime, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowRuntimeDetails, error)
}

// CreateOpenflowRuntimeOptions is based on CREATE OPENFLOW RUNTIME.
type CreateOpenflowRuntimeOptions struct {
	create          bool                      `ddl:"static" sql:"CREATE"`
	openflowRuntime bool                      `ddl:"static" sql:"OPENFLOW RUNTIME"`
	IfNotExists     *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            SchemaObjectIdentifier    `ddl:"identifier"`
	InDeployment    AccountObjectIdentifier   `ddl:"identifier" sql:"IN DEPLOYMENT"`
	ExecuteAsRole   string                    `ddl:"parameter,no_quotes" sql:"EXECUTE_AS_ROLE"`
	NodeType        OpenflowRuntimeNodeType   `ddl:"parameter,single_quotes" sql:"NODE_TYPE"`
	MinNodes        int                       `ddl:"parameter" sql:"MIN_NODES"`
	MaxNodes        int                       `ddl:"parameter" sql:"MAX_NODES"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	DisplayName     *string                   `ddl:"parameter,single_quotes" sql:"DISPLAY_NAME"`
	Comment         *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterOpenflowRuntimeOptions is based on ALTER OPENFLOW RUNTIME.
type AlterOpenflowRuntimeOptions struct {
	alter           bool                   `ddl:"static" sql:"ALTER"`
	openflowRuntime bool                   `ddl:"static" sql:"OPENFLOW RUNTIME"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
	Suspend         *bool                  `ddl:"keyword" sql:"SUSPEND"`
	Resume          *bool                  `ddl:"keyword" sql:"RESUME"`
	Terminate       *bool                  `ddl:"keyword" sql:"TERMINATE"`
	TerminateCascade *bool                 `ddl:"keyword" sql:"TERMINATE CASCADE"`
	Set             *OpenflowRuntimeSet    `ddl:"keyword" sql:"SET"`
	Unset           *OpenflowRuntimeUnset  `ddl:"list,no_parentheses" sql:"UNSET"`
}

type OpenflowRuntimeSet struct {
	MinNodes                   *int                      `ddl:"parameter" sql:"MIN_NODES"`
	MaxNodes                   *int                      `ddl:"parameter" sql:"MAX_NODES"`
	ExecuteAsRole              *string                   `ddl:"parameter,no_quotes" sql:"EXECUTE_AS_ROLE"`
	ExternalAccessIntegrations []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	DisplayName                *string                   `ddl:"parameter,single_quotes" sql:"DISPLAY_NAME"`
	Comment                    *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type OpenflowRuntimeUnset struct {
	DisplayName *bool `ddl:"keyword" sql:"DISPLAY_NAME"`
	Comment     *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropOpenflowRuntimeOptions is based on DROP OPENFLOW RUNTIME.
type DropOpenflowRuntimeOptions struct {
	drop            bool                   `ddl:"static" sql:"DROP"`
	openflowRuntime bool                   `ddl:"static" sql:"OPENFLOW RUNTIME"`
	IfExists        *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
	Cascade         *bool                  `ddl:"keyword" sql:"CASCADE"`
}

// ShowOpenflowRuntimeOptions is based on SHOW OPENFLOW RUNTIMES.
type ShowOpenflowRuntimeOptions struct {
	show             bool  `ddl:"static" sql:"SHOW"`
	openflowRuntimes bool  `ddl:"static" sql:"OPENFLOW RUNTIMES"`
	Like             *Like `ddl:"keyword" sql:"LIKE"`
}

type openflowRuntimeRow struct {
	Name                       string         `db:"name"`
	Status                     string         `db:"status"`
	Deployment                 string         `db:"deployment"`
	MinNodes                   int            `db:"min_nodes"`
	MaxNodes                   int            `db:"max_nodes"`
	NodeType                   string         `db:"node_type"`
	DisplayName                sql.NullString `db:"display_name"`
	ExternalAccessIntegrations sql.NullString `db:"external_access_integrations"`
	InitiallySuspended         bool           `db:"initially_suspended"`
	DatabaseName               string         `db:"database_name"`
	SchemaName                 string         `db:"schema_name"`
	ExecuteAsRole              string         `db:"execute_as_role"`
	Key                        sql.NullString `db:"key"`
	Owner                      string         `db:"owner"`
	Comment                    sql.NullString `db:"comment"`
	CreatedOn                  time.Time      `db:"created_on"`
	UpdatedOn                  time.Time      `db:"updated_on"`
}

type OpenflowRuntime struct {
	Name                       string
	Status                     OpenflowRuntimeStatus
	DatabaseName               string
	SchemaName                 string
	Deployment                 string
	NodeType                   OpenflowRuntimeNodeType
	MinNodes                   int
	MaxNodes                   int
	ExecuteAsRole              string
	ExternalAccessIntegrations []string
	DisplayName                *string
	Comment                    *string
	Owner                      string
	CreatedOn                  time.Time
}

func (v *OpenflowRuntime) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *OpenflowRuntime) ObjectType() ObjectType {
	return ObjectTypeOpenflowRuntime
}

// DescribeOpenflowRuntimeOptions is based on DESCRIBE OPENFLOW RUNTIME.
type DescribeOpenflowRuntimeOptions struct {
	describe        bool                   `ddl:"static" sql:"DESCRIBE"`
	openflowRuntime bool                   `ddl:"static" sql:"OPENFLOW RUNTIME"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
}

type openflowRuntimeDetailsRow struct {
	Name                       string         `db:"name"`
	Status                     string         `db:"status"`
	Deployment                 string         `db:"deployment"`
	MinNodes                   int            `db:"min_nodes"`
	MaxNodes                   int            `db:"max_nodes"`
	NodeType                   string         `db:"node_type"`
	DisplayName                sql.NullString `db:"display_name"`
	ExternalAccessIntegrations sql.NullString `db:"external_access_integrations"`
	InitiallySuspended         bool           `db:"initially_suspended"`
	DatabaseName               string         `db:"database_name"`
	SchemaName                 string         `db:"schema_name"`
	ExecuteAsRole              string         `db:"execute_as_role"`
	Key                        sql.NullString `db:"key"`
	Owner                      string         `db:"owner"`
	Comment                    sql.NullString `db:"comment"`
	CreatedOn                  time.Time      `db:"created_on"`
	UpdatedOn                  time.Time      `db:"updated_on"`
	ErrorCode                  sql.NullString `db:"error_code"`
	StatusMessage              sql.NullString `db:"status_message"`
}

type OpenflowRuntimeDetails struct {
	Name                       string
	Status                     OpenflowRuntimeStatus
	DatabaseName               string
	SchemaName                 string
	Deployment                 string
	NodeType                   OpenflowRuntimeNodeType
	MinNodes                   int
	MaxNodes                   int
	ExecuteAsRole              string
	ExternalAccessIntegrations []string
	DisplayName                *string
	Comment                    *string
	Owner                      string
	CreatedOn                  time.Time
	ErrorCode                  *string
	StatusMessage              *string
}
