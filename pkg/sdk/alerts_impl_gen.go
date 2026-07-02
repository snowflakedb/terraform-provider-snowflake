package sdk

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// Compile-time proof of interface implementation.
var _ Alerts = (*alerts)(nil)

var _ convertibleRow[Alert] = new(alertDBRow)

type alerts struct {
	client *Client
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

func (v *alerts) Show(ctx context.Context, opts *ShowAlertOptions) ([]Alert, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[alertDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[alertDBRow, Alert](dbRows)
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

func (row alertDBRow) convert() (*Alert, error) {
	alert := &Alert{
		CreatedOn:    row.CreatedOn,
		Name:         row.Name,
		DatabaseName: row.DatabaseName,
		SchemaName:   row.SchemaName,
		Owner:        row.Owner,
		Comment:      row.Comment,
		Warehouse:    row.Warehouse,
		Schedule:     row.Schedule,
		// TODO [SNOW-3108659]: use enum mapping instead
		State:     AlertState(row.State),
		Condition: row.Condition,
		Action:    row.Action,
	}
	if row.OwnerRoleType.Valid {
		alert.OwnerRoleType = row.OwnerRoleType.String
	}

	return alert, nil
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
