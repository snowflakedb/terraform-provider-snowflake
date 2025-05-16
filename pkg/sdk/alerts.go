package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
)

// Compile-time proof of interface implementation.
var _ Alerts = (*alerts)(nil)

var (
	_ validatable = new(CreateAlertOptions)
	_ validatable = new(AlterAlertOptions)
	_ validatable = new(DropAlertOptions)
	_ validatable = new(ShowAlertOptions)
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

type alerts struct {
	client *Client
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

func (opts *CreateAlertOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *alerts) Create(ctx context.Context, id SchemaObjectIdentifier, warehouse AccountObjectIdentifier, schedule string, condition string, action string, opts *CreateAlertOptions) error {
	if opts == nil {
		opts = &CreateAlertOptions{}
	}
	opts.name = id
	opts.warehouse = warehouse
	opts.schedule = schedule
	opts.name = id
	opts.condition = []AlertCondition{{Condition: []string{condition}}}
	opts.action = action
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

type AlertAction string

var (
	// AlertActionResume makes a suspended alert active.
	AlertActionResume AlertAction = "RESUME"
	// AlertActionSuspend puts the alert into a “Suspended” state.
	AlertActionSuspend AlertAction = "SUSPEND"
)

type AlertState string

var (
	AlertStateStarted   AlertState = "started"
	AlertStateSuspended AlertState = "suspended"
)

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

func (opts *AlterAlertOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if !exactlyOneValueSet(opts.Action, opts.Set, opts.Unset, opts.ModifyCondition, opts.ModifyAction) {
		errs = append(errs, errExactlyOneOf("AlterAlertOptions", "Action", "Set", "Unset", "ModifyCondition", "ModifyAction"))
	}
	return errors.Join(errs...)
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

func (v *alerts) Alter(ctx context.Context, id SchemaObjectIdentifier, opts *AlterAlertOptions) error {
	if opts == nil {
		return errors.New("alter alert options cannot be empty")
	}

	opts.name = id
	if err := opts.validate(); err != nil {
		return err
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	return err
}

// DropAlertOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-alert.
type DropAlertOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	alert    bool                   `ddl:"static" sql:"ALERT"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *DropAlertOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *alerts) Drop(ctx context.Context, id SchemaObjectIdentifier, opts *DropAlertOptions) error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	opts.name = id
	if err := opts.validate(); err != nil {
		return fmt.Errorf("validate alert options: %w", err)
	}
	sql, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = v.client.exec(ctx, sql)
	if err != nil {
		return err
	}
	return err
}

func (v *alerts) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropAlertOptions{IfExists: Bool(true)}) }, ctx, id)
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

func (v *Alert) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Alert) ObjectType() ObjectType {
	return ObjectTypeAlert
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

func (row alertDBRow) convert() *Alert {
	alert := &Alert{
		CreatedOn:    row.CreatedOn,
		Name:         row.Name,
		DatabaseName: row.DatabaseName,
		SchemaName:   row.SchemaName,
		Owner:        row.Owner,
		Comment:      row.Comment,
		Warehouse:    row.Warehouse,
		Schedule:     row.Schedule,
		State:        AlertState(row.State),
		Condition:    row.Condition,
		Action:       row.Action,
	}
	if row.OwnerRoleType.Valid {
		alert.OwnerRoleType = row.OwnerRoleType.String
	}

	return alert
}

func (opts *ShowAlertOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (v *alerts) Show(ctx context.Context, opts *ShowAlertOptions) ([]Alert, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[alertDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[alertDBRow, Alert](dbRows)
	return resultList, nil
}

func (v *alerts) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Alert, error) {
	alerts, err := v.Show(ctx, &ShowAlertOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
		In: &In{
			Schema: id.SchemaId(),
		},
	})
	if err != nil {
		return nil, err
	}

	return collections.FindFirst(alerts, func(alert Alert) bool {
		return alert.ID().FullyQualifiedName() == id.FullyQualifiedName()
	})
}

func (v *alerts) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Alert, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

// describeAlertOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-alert.
type describeAlertOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	alert    bool                   `ddl:"static" sql:"ALERT"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

func (opts *describeAlertOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
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

func (row alertDBRow) toAlertDetails() (*AlertDetails, error) {
	return &AlertDetails{
		CreatedOn:    row.CreatedOn,
		Name:         row.Name,
		DatabaseName: row.DatabaseName,
		SchemaName:   row.SchemaName,
		Owner:        row.Owner,
		Comment:      row.Comment,
		Warehouse:    row.Warehouse,
		Schedule:     row.Schedule,
		State:        row.State,
		Condition:    row.Condition,
		Action:       row.Action,
	}, nil
}

func (v *alerts) Describe(ctx context.Context, id SchemaObjectIdentifier) (*AlertDetails, error) {
	opts := &describeAlertOptions{
		name: id,
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}

	sql, err := structToSQL(opts)
	if err != nil {
		return nil, err
	}

	// SHOW ALERTS and DESCRIBE ALERT SQL statements return the same output
	dest := alertDBRow{}
	err = v.client.queryOne(ctx, &dest, sql)
	if err != nil {
		return nil, err
	}

	return dest.toAlertDetails()
}
