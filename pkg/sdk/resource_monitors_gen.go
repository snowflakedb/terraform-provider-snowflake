package sdk

import (
	"context"
	"database/sql"
	"time"
)

var (
	_ validatable = new(CreateResourceMonitorOptions)
	_ validatable = new(AlterResourceMonitorOptions)
	_ validatable = new(DropResourceMonitorOptions)
	_ validatable = new(ShowResourceMonitorOptions)
)

type ResourceMonitors interface {
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateResourceMonitorOptions) error
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterResourceMonitorOptions) error
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropResourceMonitorOptions) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, opts *ShowResourceMonitorOptions) ([]ResourceMonitor, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ResourceMonitor, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*ResourceMonitor, error)
}

// CreateResourceMonitorOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-resource-monitor.
type CreateResourceMonitorOptions struct {
	create          bool                    `ddl:"static" sql:"CREATE"`
	OrReplace       *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"`
	IfNotExists     *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	With            *ResourceMonitorWith    `ddl:"keyword" sql:"WITH"`
}

type ResourceMonitorWith struct {
	CreditQuota    *int                `ddl:"parameter,equals" sql:"CREDIT_QUOTA"`
	Frequency      *Frequency          `ddl:"parameter,equals" sql:"FREQUENCY"`
	StartTimestamp *string             `ddl:"parameter,equals,single_quotes" sql:"START_TIMESTAMP"`
	EndTimestamp   *string             `ddl:"parameter,equals,single_quotes" sql:"END_TIMESTAMP"`
	NotifyUsers    *NotifyUsers        `ddl:"parameter,equals" sql:"NOTIFY_USERS"`
	Triggers       []TriggerDefinition `ddl:"keyword,no_comma" sql:"TRIGGERS"`
}

type TriggerDefinition struct {
	Threshold     int           `ddl:"parameter,no_equals" sql:"ON"`
	TriggerAction TriggerAction `ddl:"parameter,no_equals" sql:"PERCENT DO"`
}

type NotifyUsers struct {
	Users []NotifiedUser `ddl:"list,parentheses,comma"`
}

type NotifiedUser struct {
	Name AccountObjectIdentifier `ddl:"identifier"`
}

// AlterResourceMonitorOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-resource-monitor.
type AlterResourceMonitorOptions struct {
	alter           bool                    `ddl:"static" sql:"ALTER"`
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"`
	IfExists        *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name            AccountObjectIdentifier `ddl:"identifier"`
	Set             *ResourceMonitorSet     `ddl:"keyword" sql:"SET"`
	Unset           *ResourceMonitorUnset   `ddl:"keyword" sql:"SET"`
	Triggers        []TriggerDefinition     `ddl:"keyword,no_comma" sql:"TRIGGERS"`
}

type ResourceMonitorSet struct {
	// at least one
	CreditQuota    *int         `ddl:"parameter,equals" sql:"CREDIT_QUOTA"`
	Frequency      *Frequency   `ddl:"parameter,equals" sql:"FREQUENCY"`
	StartTimestamp *string      `ddl:"parameter,equals,single_quotes" sql:"START_TIMESTAMP"`
	EndTimestamp   *string      `ddl:"parameter,equals,single_quotes" sql:"END_TIMESTAMP"`
	NotifyUsers    *NotifyUsers `ddl:"parameter,equals" sql:"NOTIFY_USERS"`
}

type ResourceMonitorUnset struct {
	CreditQuota  *bool `ddl:"keyword" sql:"CREDIT_QUOTA = null"`
	EndTimestamp *bool `ddl:"keyword" sql:"END_TIMESTAMP = null"`
	NotifyUsers  *bool `ddl:"keyword" sql:"NOTIFY_USERS = ()"`
}

// DropResourceMonitorOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-resource-monitor.
type DropResourceMonitorOptions struct {
	drop            bool                    `ddl:"static" sql:"DROP"`
	resourceMonitor bool                    `ddl:"static" sql:"RESOURCE MONITOR"`
	IfExists        *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name            AccountObjectIdentifier `ddl:"identifier"`
}

// ShowResourceMonitorOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-resource-monitors.
type ShowResourceMonitorOptions struct {
	show             bool  `ddl:"static" sql:"SHOW"`
	resourceMonitors bool  `ddl:"static" sql:"RESOURCE MONITORS"`
	Like             *Like `ddl:"keyword" sql:"LIKE"`
}

type resourceMonitorRow struct {
	Name               string         `db:"name"`
	CreditQuota        sql.NullString `db:"credit_quota"`
	UsedCredits        sql.NullString `db:"used_credits"`
	RemainingCredits   sql.NullString `db:"remaining_credits"`
	Level              sql.NullString `db:"level"`
	Frequency          sql.NullString `db:"frequency"`
	StartTime          sql.NullString `db:"start_time"`
	EndTime            sql.NullString `db:"end_time"`
	NotifyAt           sql.NullString `db:"notify_at"`
	SuspendAt          sql.NullString `db:"suspend_at"`
	SuspendImmediateAt sql.NullString `db:"suspend_immediately_at"`
	CreatedOn          time.Time      `db:"created_on"`
	Owner              string         `db:"owner"`
	Comment            sql.NullString `db:"comment"`
	NotifyUsers        sql.NullString `db:"notify_users"`
}

type ResourceMonitor struct {
	Name               string
	CreditQuota        float64
	UsedCredits        float64
	RemainingCredits   float64
	Level              *ResourceMonitorLevel
	Frequency          Frequency
	StartTime          string
	EndTime            string
	NotifyAt           []int
	SuspendAt          *int
	SuspendImmediateAt *int
	CreatedOn          time.Time
	Owner              string
	Comment            string
	NotifyUsers        []string
}

func (v *ResourceMonitor) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *ResourceMonitor) ObjectType() ObjectType {
	return ObjectTypeResourceMonitor
}
