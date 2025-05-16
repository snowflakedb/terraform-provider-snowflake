package sdk

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
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

var _ ResourceMonitors = (*resourceMonitors)(nil)

type resourceMonitors struct {
	client *Client
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

func (row *resourceMonitorRow) convert() (*ResourceMonitor, error) {
	resourceMonitor := &ResourceMonitor{
		Name:      row.Name,
		CreatedOn: row.CreatedOn,
		Owner:     row.Owner,
	}
	if row.CreditQuota.Valid {
		creditQuota, err := strconv.ParseFloat(row.CreditQuota.String, 64)
		if err != nil {
			return nil, err
		}
		resourceMonitor.CreditQuota = creditQuota
	}

	if row.UsedCredits.Valid {
		usedCredits, err := strconv.ParseFloat(row.UsedCredits.String, 64)
		if err != nil {
			return nil, err
		}
		resourceMonitor.UsedCredits = usedCredits
	}

	if row.RemainingCredits.Valid {
		remainingCredits, err := strconv.ParseFloat(row.RemainingCredits.String, 64)
		if err != nil {
			return nil, err
		}
		resourceMonitor.RemainingCredits = remainingCredits
	}

	if row.Level.Valid {
		level, err := ToResourceMonitorLevel(row.Level.String)
		if err != nil {
			log.Printf("[DEBUG] unable to parse resource monitor level: %v", err)
		} else {
			resourceMonitor.Level = &level
		}
	}

	if row.Frequency.Valid {
		frequency, err := ToResourceMonitorFrequency(row.Frequency.String)
		if err != nil {
			return nil, err
		}
		resourceMonitor.Frequency = *frequency
	}

	if row.StartTime.Valid {
		resourceMonitor.StartTime = row.StartTime.String
	}

	if row.EndTime.Valid {
		resourceMonitor.EndTime = row.EndTime.String
	}

	notifyTriggers, err := extractTriggerInts(row.NotifyAt)
	if err != nil {
		return nil, err
	}
	resourceMonitor.NotifyAt = notifyTriggers

	suspendTriggers, err := extractTriggerInts(row.SuspendAt)
	if err != nil {
		return nil, err
	}
	if len(suspendTriggers) > 0 {
		resourceMonitor.SuspendAt = &suspendTriggers[0]
	}

	suspendImmediateTriggers, err := extractTriggerInts(row.SuspendImmediateAt)
	if err != nil {
		return nil, err
	}
	if len(suspendImmediateTriggers) > 0 {
		resourceMonitor.SuspendImmediateAt = &suspendImmediateTriggers[0]
	}

	if row.Comment.Valid {
		resourceMonitor.Comment = row.Comment.String
	}

	if row.NotifyUsers.Valid && row.NotifyUsers.String != "" {
		resourceMonitor.NotifyUsers = strings.Split(row.NotifyUsers.String, ", ")
	}

	return resourceMonitor, nil
}

// extractTriggerInts converts the triggers in the DB (stored as a comma separated string with trailing `%` signs) into a slice of ints.
func extractTriggerInts(s sql.NullString) ([]int, error) {
	// Check if this is NULL
	if !s.Valid || s.String == "" {
		return []int{}, nil
	}
	ints := strings.Split(s.String, ",")
	out := make([]int, 0, len(ints))
	for _, i := range ints {
		numberToParse := strings.TrimRight(i, "%")
		myInt, err := strconv.Atoi(numberToParse)
		if err != nil {
			return out, fmt.Errorf("failed to convert %v to integer err = %w", numberToParse, err)
		}
		out = append(out, myInt)
	}
	return out, nil
}

func (v *ResourceMonitor) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *ResourceMonitor) ObjectType() ObjectType {
	return ObjectTypeResourceMonitor
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

func (opts *CreateResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if everyValueSet(opts.OrReplace, opts.IfNotExists) {
		errs = append(errs, errOneOf("CreateResourceMonitorOptions", "OrReplace", "IfNotExists"))
	}
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, errors.Join(ErrInvalidObjectIdentifier))
	}
	if valueSet(opts.With) && everyValueNil(opts.With.CreditQuota, opts.With.Frequency, opts.With.StartTimestamp, opts.With.EndTimestamp, opts.With.NotifyUsers) && valueSet(opts.With.Triggers) {
		errs = append(errs, fmt.Errorf("due to Snowflake limitations you cannot create Resource Monitor with only triggers set"))
	}
	return errors.Join(errs...)
}

func (v *resourceMonitors) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateResourceMonitorOptions) error {
	if opts == nil {
		opts = &CreateResourceMonitorOptions{}
	}
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

type ResourceMonitorLevel string

const (
	ResourceMonitorLevelAccount   ResourceMonitorLevel = "ACCOUNT"
	ResourceMonitorLevelWarehouse ResourceMonitorLevel = "WAREHOUSE"
)

func ToResourceMonitorLevel(s string) (ResourceMonitorLevel, error) {
	switch level := ResourceMonitorLevel(strings.ToUpper(s)); level {
	case ResourceMonitorLevelAccount,
		ResourceMonitorLevelWarehouse:
		return level, nil
	default:
		return "", fmt.Errorf("invalid resource monitor level: %s", s)
	}
}

type TriggerDefinition struct {
	Threshold     int           `ddl:"parameter,no_equals" sql:"ON"`
	TriggerAction TriggerAction `ddl:"parameter,no_equals" sql:"PERCENT DO"`
}

type TriggerAction string

const (
	TriggerActionSuspend          TriggerAction = "SUSPEND"
	TriggerActionSuspendImmediate TriggerAction = "SUSPEND_IMMEDIATE"
	TriggerActionNotify           TriggerAction = "NOTIFY"
)

func ToResourceMonitorTriggerAction(s string) (*TriggerAction, error) {
	switch action := TriggerAction(strings.ToUpper(s)); action {
	case TriggerActionSuspend,
		TriggerActionSuspendImmediate,
		TriggerActionNotify:
		return &action, nil
	default:
		return nil, fmt.Errorf("invalid trigger action type: %s", s)
	}
}

type NotifyUsers struct {
	Users []NotifiedUser `ddl:"list,parentheses,comma"`
}

type NotifiedUser struct {
	Name AccountObjectIdentifier `ddl:"identifier"`
}

type Frequency string

func ToResourceMonitorFrequency(s string) (*Frequency, error) {
	switch frequency := Frequency(strings.ToUpper(s)); frequency {
	case FrequencyDaily,
		FrequencyWeekly,
		FrequencyMonthly,
		FrequencyYearly,
		FrequencyNever:
		return &frequency, nil
	default:
		return nil, fmt.Errorf("invalid frequency type: %s", s)
	}
}

const (
	FrequencyMonthly Frequency = "MONTHLY"
	FrequencyDaily   Frequency = "DAILY"
	FrequencyWeekly  Frequency = "WEEKLY"
	FrequencyYearly  Frequency = "YEARLY"
	FrequencyNever   Frequency = "NEVER"
)

var AllFrequencyValues = []Frequency{
	FrequencyMonthly,
	FrequencyDaily,
	FrequencyWeekly,
	FrequencyYearly,
	FrequencyNever,
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

func (opts *AlterResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	var errs []error
	if !ValidObjectIdentifier(opts.name) {
		errs = append(errs, ErrInvalidObjectIdentifier)
	}
	if everyValueNil(opts.Set, opts.Unset, opts.Triggers) {
		errs = append(errs, errAtLeastOneOf("AlterResourceMonitorOptions", "Set", "Unset", "Triggers"))
	}
	if set := opts.Set; valueSet(set) {
		if everyValueNil(set.CreditQuota, set.Frequency, set.StartTimestamp, set.EndTimestamp, set.NotifyUsers) {
			errs = append(errs, errAtLeastOneOf("ResourceMonitorSet", "CreditQuota", "Frequency", "StartTimestamp", "EndTimestamp", "NotifyUsers"))
		}
		if (set.Frequency != nil && set.StartTimestamp == nil) || (set.Frequency == nil && set.StartTimestamp != nil) {
			errs = append(errs, errors.New("must specify frequency and start time together"))
		}
	}
	return errors.Join(errs...)
}

func (v *resourceMonitors) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterResourceMonitorOptions) error {
	if opts == nil {
		opts = &AlterResourceMonitorOptions{}
	}
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
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

func (opts *DropResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	if !ValidObjectIdentifier(opts.name) {
		return errors.Join(ErrInvalidObjectIdentifier)
	}
	return nil
}

func (v *resourceMonitors) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropResourceMonitorOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *resourceMonitors) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropResourceMonitorOptions{IfExists: Bool(true)}) }, ctx, id)
}

// ShowResourceMonitorOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-resource-monitors.
type ShowResourceMonitorOptions struct {
	show             bool  `ddl:"static" sql:"SHOW"`
	resourceMonitors bool  `ddl:"static" sql:"RESOURCE MONITORS"`
	Like             *Like `ddl:"keyword" sql:"LIKE"`
}

func (opts *ShowResourceMonitorOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (v *resourceMonitors) Show(ctx context.Context, opts *ShowResourceMonitorOptions) ([]ResourceMonitor, error) {
	opts = createIfNil(opts)
	dbRows, err := validateAndQuery[resourceMonitorRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := make([]ResourceMonitor, len(dbRows))
	for i, row := range dbRows {
		resourceMonitor, err := row.convert()
		if err != nil {
			return nil, err
		}
		resultList[i] = *resourceMonitor
	}
	return resultList, nil
}

func (v *resourceMonitors) ShowByID(ctx context.Context, id AccountObjectIdentifier) (*ResourceMonitor, error) {
	resourceMonitors, err := v.Show(ctx, &ShowResourceMonitorOptions{
		Like: &Like{
			Pattern: String(id.Name()),
		},
	})
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(resourceMonitors, func(r ResourceMonitor) bool { return r.ID().Name() == id.Name() })
}

func (v *resourceMonitors) ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*ResourceMonitor, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}
