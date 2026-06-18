package sdk

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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

func (r resourceMonitorRow) additionalConvert(resourceMonitor *ResourceMonitor) error {
	if r.CreditQuota.Valid {
		creditQuota, err := strconv.ParseFloat(r.CreditQuota.String, 64)
		if err != nil {
			return err
		}
		resourceMonitor.CreditQuota = creditQuota
	}

	if r.UsedCredits.Valid {
		usedCredits, err := strconv.ParseFloat(r.UsedCredits.String, 64)
		if err != nil {
			return err
		}
		resourceMonitor.UsedCredits = usedCredits
	}

	if r.RemainingCredits.Valid {
		remainingCredits, err := strconv.ParseFloat(r.RemainingCredits.String, 64)
		if err != nil {
			return err
		}
		resourceMonitor.RemainingCredits = remainingCredits
	}

	notifyTriggers, err := extractTriggerInts(r.NotifyAt)
	if err != nil {
		return err
	}
	resourceMonitor.NotifyAt = notifyTriggers

	suspendTriggers, err := extractTriggerInts(r.SuspendAt)
	if err != nil {
		return err
	}
	if len(suspendTriggers) > 0 {
		resourceMonitor.SuspendAt = &suspendTriggers[0]
	}

	suspendImmediateTriggers, err := extractTriggerInts(r.SuspendImmediatelyAt)
	if err != nil {
		return err
	}
	if len(suspendImmediateTriggers) > 0 {
		resourceMonitor.SuspendImmediatelyAt = &suspendImmediateTriggers[0]
	}

	return nil
}

func (opts *CreateResourceMonitorOptions) additionalValidations() error {
	if valueSet(opts.With) && everyValueNil(opts.With.CreditQuota, opts.With.Frequency, opts.With.StartTimestamp, opts.With.EndTimestamp, opts.With.NotifyUsers) && valueSet(opts.With.Triggers) {
		return fmt.Errorf("due to Snowflake limitations you cannot create Resource Monitor with only triggers set")
	}
	return nil
}

func (opts *AlterResourceMonitorOptions) additionalValidations() error {
	var errs []error
	if set := opts.Set; valueSet(set) {
		if (set.Frequency != nil && set.StartTimestamp == nil) || (set.Frequency == nil && set.StartTimestamp != nil) {
			errs = append(errs, errors.New("must specify frequency and start time together"))
		}
	}
	return errors.Join(errs...)
}
