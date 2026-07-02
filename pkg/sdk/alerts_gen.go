package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Alerts interface {
	Create(ctx context.Context, id SchemaObjectIdentifier, warehouse AccountObjectIdentifier, schedule string, condition string, action string, opts *CreateAlertOptions) error
	Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterAlertOptions) error
	Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropAlertOptions) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, opts *ShowAlertOptions) ([]Alert, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Alert, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Alert, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*AlertDetails, error)
}

// CreateAlertOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-alert.
type CreateAlertOptions struct {
	create      bool                   `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	alert       bool                   `ddl:"static" sql:"ALERT"`
	IfNotExists *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        SchemaObjectIdentifier `ddl:"identifier"`

	// required
	warehouse AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
	schedule  string                  `ddl:"parameter,single_quotes" sql:"SCHEDULE"`

	// optional
	Comment *string `ddl:"parameter,single_quotes" sql:"COMMENT"`

	// required
	condition []AlertCondition `ddl:"keyword,parentheses,no_comma"   sql:"IF"`
	action    string           `ddl:"parameter,no_equals" sql:"THEN"`
}

type AlertCondition struct {
	Condition []string `ddl:"keyword,parentheses,no_comma" sql:"EXISTS"`
}

// AlterAlertOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-alert.
type AlterAlertOptions struct {
	alter    bool                   `ddl:"static" sql:"ALTER"`
	alert    bool                   `ddl:"static" sql:"ALERT"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`

	// One of
	Action          *AlertAction `ddl:"keyword"`
	Set             *AlertSet    `ddl:"keyword" sql:"SET"`
	Unset           *AlertUnset  `ddl:"keyword" sql:"UNSET"`
	ModifyCondition *[]string    `ddl:"keyword,parentheses,no_comma" sql:"MODIFY CONDITION EXISTS"`
	ModifyAction    *string      `ddl:"parameter,no_equals" sql:"MODIFY ACTION"`
}

type AlertSet struct {
	Warehouse *AccountObjectIdentifier `ddl:"identifier,equals" sql:"WAREHOUSE"`
	Schedule  *string                  `ddl:"parameter,single_quotes" sql:"SCHEDULE"`
	Comment   *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type AlertUnset struct {
	Warehouse *bool `ddl:"keyword" sql:"WAREHOUSE"`
	Schedule  *bool `ddl:"keyword" sql:"SCHEDULE"`
	Comment   *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropAlertOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-alert.
type DropAlertOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	alert    bool                   `ddl:"static" sql:"ALERT"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowAlertOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-alerts.
type ShowAlertOptions struct {
	show   bool  `ddl:"static" sql:"SHOW"`
	Terse  *bool `ddl:"keyword" sql:"TERSE"`
	alerts bool  `ddl:"static" sql:"ALERTS"`

	// optional
	Like       *Like   `ddl:"keyword" sql:"LIKE"`
	In         *In     `ddl:"keyword" sql:"IN"`
	StartsWith *string `ddl:"parameter,no_equals,single_quotes" sql:"STARTS WITH"`
	Limit      *int    `ddl:"parameter,no_equals" sql:"LIMIT"`
}

// describeAlertOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-alert.
type describeAlertOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	alert    bool                   `ddl:"static" sql:"ALERT"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

type alertDBRow struct {
	CreatedOn     time.Time      `db:"created_on"`
	Name          string         `db:"name"`
	DatabaseName  string         `db:"database_name"`
	SchemaName    string         `db:"schema_name"`
	Owner         string         `db:"owner"`
	Comment       *string        `db:"comment"`
	Warehouse     string         `db:"warehouse"`
	Schedule      string         `db:"schedule"`
	State         string         `db:"state"` // suspended, started
	Condition     string         `db:"condition"`
	Action        string         `db:"action"`
	OwnerRoleType sql.NullString `db:"owner_role_type"`
}

type Alert struct {
	CreatedOn     time.Time
	Name          string
	DatabaseName  string
	SchemaName    string
	Owner         string
	Comment       *string
	Warehouse     string
	Schedule      string
	State         AlertState
	Condition     string
	Action        string
	OwnerRoleType string
}

func (v *Alert) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Alert) ObjectType() ObjectType {
	return ObjectTypeAlert
}

type AlertDetails struct {
	CreatedOn    time.Time
	Name         string
	DatabaseName string
	SchemaName   string
	Owner        string
	Comment      *string
	Warehouse    string
	Schedule     string
	State        string
	Condition    string
	Action       string
}
