package sdk

import (
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResourceMonitorCreate(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &CreateResourceMonitorOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: OrReplace and IfExists specified", func(t *testing.T) {
		opts := &CreateResourceMonitorOptions{
			name:        id,
			OrReplace:   Bool(true),
			IfNotExists: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateResourceMonitorOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("with complete options", func(t *testing.T) {
		creditQuota := Int(100)
		frequency := FrequencyMonthly
		startTimeStamp := "IMMIEDIATELY"
		endTimeStamp := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).String()
		notifiedUsers := []NotifiedUser{{Name: NewAccountObjectIdentifier("FIRST_USER")}, {Name: NewAccountObjectIdentifier("SECOND_USER")}}
		triggers := []TriggerDefinition{
			{
				Threshold:     50,
				TriggerAction: TriggerActionSuspendImmediate,
			},
			{
				Threshold:     100,
				TriggerAction: TriggerActionNotify,
			},
		}

		opts := &CreateResourceMonitorOptions{
			OrReplace: Bool(true),
			With: &ResourceMonitorWith{
				CreditQuota:    creditQuota,
				Frequency:      &frequency,
				StartTimestamp: &startTimeStamp,
				EndTimestamp:   &endTimeStamp,
				NotifyUsers:    &NotifyUsers{notifiedUsers},
				Triggers:       triggers,
			},
			name: id,
		}

		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE RESOURCE MONITOR %s WITH CREDIT_QUOTA = 100 FREQUENCY = MONTHLY START_TIMESTAMP = 'IMMIEDIATELY' END_TIMESTAMP = '%s' NOTIFY_USERS = ("FIRST_USER", "SECOND_USER") TRIGGERS ON 50 PERCENT DO SUSPEND_IMMEDIATE ON 100 PERCENT DO NOTIFY`,
			id.FullyQualifiedName(),
			endTimeStamp,
		)
	})
}

func TestResourceMonitorAlter(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("validation: empty options", func(t *testing.T) {
		opts := &AlterResourceMonitorOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &AlterResourceMonitorOptions{
			name: id,
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterResourceMonitorOptions", "Set", "Unset", "Triggers"))
	})

	t.Run("validation: no option for set provided", func(t *testing.T) {
		opts := &AlterResourceMonitorOptions{
			name: id,
			Set:  &ResourceMonitorSet{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errAtLeastOneOf("AlterResourceMonitorOptions.Set", "CreditQuota", "Frequency", "StartTimestamp", "EndTimestamp", "NotifyUsers"))
	})

	t.Run("with a single set", func(t *testing.T) {
		newCreditQuota := Int(50)
		opts := &AlterResourceMonitorOptions{
			name: id,
			Set: &ResourceMonitorSet{
				CreditQuota: newCreditQuota,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER RESOURCE MONITOR %s SET CREDIT_QUOTA = %d", id.FullyQualifiedName(), *newCreditQuota)
	})

	t.Run("set notify users", func(t *testing.T) {
		opts := &AlterResourceMonitorOptions{
			name: id,
			Set: &ResourceMonitorSet{
				NotifyUsers: &NotifyUsers{
					Users: []NotifiedUser{
						{Name: NewAccountObjectIdentifier("user1")},
						{Name: NewAccountObjectIdentifier("user2")},
					},
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER RESOURCE MONITOR %s SET NOTIFY_USERS = (\"user1\", \"user2\")", id.FullyQualifiedName())
	})

	t.Run("with a multitple set", func(t *testing.T) {
		newCreditQuota := Int(50)
		newFrequency := FrequencyYearly
		newStartTimeStamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).String()
		opts := &AlterResourceMonitorOptions{
			name: id,
			Set: &ResourceMonitorSet{
				CreditQuota:    newCreditQuota,
				Frequency:      &newFrequency,
				StartTimestamp: &newStartTimeStamp,
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER RESOURCE MONITOR %s SET CREDIT_QUOTA = %d FREQUENCY = %s START_TIMESTAMP = '%s'", id.FullyQualifiedName(), *newCreditQuota, newFrequency, newStartTimeStamp)
	})

	t.Run("with unset", func(t *testing.T) {
		opts := &AlterResourceMonitorOptions{
			name: id,
			Unset: &ResourceMonitorUnset{
				CreditQuota:  Bool(true),
				EndTimestamp: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "ALTER RESOURCE MONITOR %s SET CREDIT_QUOTA = null END_TIMESTAMP = null", id.FullyQualifiedName())
	})
}

func TestResourceMonitorDrop(t *testing.T) {
	id := randomAccountObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &DropResourceMonitorOptions{}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &DropResourceMonitorOptions{
			name: id,
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP RESOURCE MONITOR %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := &DropResourceMonitorOptions{
			name:     id,
			IfExists: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, "DROP RESOURCE MONITOR IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestResourceMonitorShow(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	t.Run("empty options", func(t *testing.T) {
		opts := &ShowResourceMonitorOptions{}
		assertOptsValidAndSQLEquals(t, opts, "SHOW RESOURCE MONITORS")
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowResourceMonitorOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, "SHOW RESOURCE MONITORS LIKE '%s'", id.Name())
	})
}

func TestExtractTriggerInts(t *testing.T) {
	testCases := []struct {
		Input    sql.NullString
		Expected []int
		Error    string
	}{
		{Input: sql.NullString{String: "51%,63%,123%", Valid: true}, Expected: []int{51, 63, 123}},
		{Input: sql.NullString{String: "51%,63%", Valid: true}, Expected: []int{51, 63}},
		{Input: sql.NullString{String: "51%", Valid: true}, Expected: []int{51}},
		{Input: sql.NullString{String: "", Valid: false}, Expected: []int{}},
		{Input: sql.NullString{String: "", Valid: true}, Expected: []int{}},
		{Input: sql.NullString{String: "51,63", Valid: true}, Expected: []int{51, 63}},
		{Input: sql.NullString{String: "1", Valid: true}, Expected: []int{1}},
		{Input: sql.NullString{String: "ab,cd", Valid: true}, Error: "failed to convert ab to integer err = strconv.Atoi"},
		{Input: sql.NullString{String: "12,,34", Valid: true}, Error: "failed to convert  to integer err = strconv.Atoi"},
		{Input: sql.NullString{String: ",", Valid: true}, Error: "failed to convert  to integer err = strconv.Atoi"},
		{Input: sql.NullString{String: "12.34", Valid: true}, Error: "failed to convert 12.34 to integer err = strconv.Atoi"},
	}

	for _, tc := range testCases {
		t.Run("extract trigger ints: "+tc.Input.String+":"+strconv.FormatBool(tc.Input.Valid), func(t *testing.T) {
			result, err := extractTriggerInts(tc.Input)
			if tc.Error != "" {
				require.ErrorContains(t, err, tc.Error)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.Expected, result)
			}
		})
	}
}
