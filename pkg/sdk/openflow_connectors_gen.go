package sdk

import (
	"context"
	"database/sql"
	"time"
)

type OpenflowConnectors interface {
	Create(ctx context.Context, request *CreateOpenflowConnectorRequest) error
	Alter(ctx context.Context, request *AlterOpenflowConnectorRequest) error
	Drop(ctx context.Context, request *DropOpenflowConnectorRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, request *ShowOpenflowConnectorRequest) ([]OpenflowConnector, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowConnector, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowConnector, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*OpenflowConnectorDetails, error)
}

// CreateOpenflowConnectorOptions is based on CREATE OPENFLOW CONNECTOR.
type CreateOpenflowConnectorOptions struct {
	create             bool                   `ddl:"static" sql:"CREATE"`
	openflowConnector  bool                   `ddl:"static" sql:"OPENFLOW CONNECTOR"`
	IfNotExists        *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name               SchemaObjectIdentifier `ddl:"identifier"`
	InRuntime          SchemaObjectIdentifier `ddl:"identifier" sql:"IN RUNTIME"`
	FromDefinition     *string                `ddl:"parameter,no_quotes,no_equals" sql:"FROM DEFINITION"`
	FromStage          *string                `ddl:"parameter,single_quotes,no_equals" sql:"FROM"`
	DisplayName        *string                `ddl:"parameter,single_quotes" sql:"DISPLAY_NAME"`
	Comment            *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterOpenflowConnectorOptions is based on ALTER OPENFLOW CONNECTOR.
type AlterOpenflowConnectorOptions struct {
	alter             bool                   `ddl:"static" sql:"ALTER"`
	openflowConnector bool                   `ddl:"static" sql:"OPENFLOW CONNECTOR"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
	Start             *bool                  `ddl:"keyword" sql:"START"`
	Stop              *bool                  `ddl:"keyword" sql:"STOP"`
	Set               *OpenflowConnectorSet  `ddl:"keyword" sql:"SET"`
	Unset             *OpenflowConnectorUnset `ddl:"list,no_parentheses" sql:"UNSET"`
}

type OpenflowConnectorSet struct {
	DisplayName *string `ddl:"parameter,single_quotes" sql:"DISPLAY_NAME"`
	Comment     *string `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type OpenflowConnectorUnset struct {
	DisplayName *bool `ddl:"keyword" sql:"DISPLAY_NAME"`
	Comment     *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropOpenflowConnectorOptions is based on DROP OPENFLOW CONNECTOR.
type DropOpenflowConnectorOptions struct {
	drop              bool                   `ddl:"static" sql:"DROP"`
	openflowConnector bool                   `ddl:"static" sql:"OPENFLOW CONNECTOR"`
	IfExists          *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowOpenflowConnectorOptions is based on SHOW OPENFLOW CONNECTORS.
type ShowOpenflowConnectorOptions struct {
	show               bool  `ddl:"static" sql:"SHOW"`
	openflowConnectors bool  `ddl:"static" sql:"OPENFLOW CONNECTORS"`
	Like               *Like `ddl:"keyword" sql:"LIKE"`
}

type openflowConnectorRow struct {
	Name                string         `db:"name"`
	Status              string         `db:"status"`
	Runtime             string         `db:"runtime"`
	ConnectorDefinition sql.NullString `db:"connector_definition"`
	DisplayName         sql.NullString `db:"display_name"`
	DatabaseName        string         `db:"database_name"`
	SchemaName          string         `db:"schema_name"`
	Owner               string         `db:"owner"`
	Comment             sql.NullString `db:"comment"`
	CreatedOn           time.Time      `db:"created_on"`
	UpdatedOn           time.Time      `db:"updated_on"`
}

type OpenflowConnector struct {
	Name                string
	Status              OpenflowConnectorStatus
	DatabaseName        string
	SchemaName          string
	Runtime             string
	ConnectorDefinition *string
	Started             bool // derived from Status: true when ACTIVE or STARTING
	DisplayName         *string
	Comment             *string
	Owner               string
	CreatedOn           time.Time
	UpdatedOn           time.Time
}

func (v *OpenflowConnector) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *OpenflowConnector) ObjectType() ObjectType {
	return ObjectTypeOpenflowConnector
}

// DescribeOpenflowConnectorOptions is based on DESCRIBE OPENFLOW CONNECTOR.
type DescribeOpenflowConnectorOptions struct {
	describe          bool                   `ddl:"static" sql:"DESCRIBE"`
	openflowConnector bool                   `ddl:"static" sql:"OPENFLOW CONNECTOR"`
	name              SchemaObjectIdentifier `ddl:"identifier"`
}

type openflowConnectorDetailsRow struct {
	Name                string         `db:"name"`
	Status              string         `db:"status"`
	Runtime             string         `db:"runtime"`
	ConnectorDefinition sql.NullString `db:"connector_definition"`
	DisplayName         sql.NullString `db:"display_name"`
	DatabaseName        string         `db:"database_name"`
	SchemaName          string         `db:"schema_name"`
	Owner               string         `db:"owner"`
	Comment             sql.NullString `db:"comment"`
	CreatedOn           time.Time      `db:"created_on"`
	UpdatedOn           time.Time      `db:"updated_on"`
	ErrorCode           sql.NullString `db:"error_code"`
	StatusMessage       sql.NullString `db:"status_message"`
}

type OpenflowConnectorDetails struct {
	Name                string
	Status              OpenflowConnectorStatus
	DatabaseName        string
	SchemaName          string
	Runtime             string
	ConnectorDefinition *string
	Started             bool // derived from Status: true when ACTIVE or STARTING
	DisplayName         *string
	Comment             *string
	Owner               string
	CreatedOn           time.Time
	UpdatedOn           time.Time
	ErrorCode           *string
	StatusMessage       *string
}
