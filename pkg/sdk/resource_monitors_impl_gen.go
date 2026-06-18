package sdk

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ ResourceMonitors = (*resourceMonitors)(nil)

type resourceMonitors struct {
	client *Client
}

func (v *resourceMonitors) Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateResourceMonitorOptions) error {
	if opts == nil {
		opts = &CreateResourceMonitorOptions{}
	}
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *resourceMonitors) Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterResourceMonitorOptions) error {
	if opts == nil {
		opts = &AlterResourceMonitorOptions{}
	}
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *resourceMonitors) Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropResourceMonitorOptions) error {
	opts = createIfNil(opts)
	opts.name = id
	return validateAndExec(v.client, ctx, opts)
}

func (v *resourceMonitors) DropSafely(ctx context.Context, id AccountObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, id, &DropResourceMonitorOptions{IfExists: Bool(true)}) }, ctx, id)
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
