package sdk

import (
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	id := resourceMonitorsTestIdAccountObjectIdentifier
	startTimestamp := "IMMIEDIATELY"
	endTimestamp := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).String()
	alterStartTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).String()
	firstUser := NewAccountObjectIdentifier("FIRST_USER")
	secondUser := NewAccountObjectIdentifier("SECOND_USER")
	user1 := NewAccountObjectIdentifier("user1")
	user2 := NewAccountObjectIdentifier("user2")

	resourceMonitorsTests.Create.
		withAdditionalValidationCase(
			"validation_Create_withOnlyTriggers",
			func(opts *CreateResourceMonitorOptions) {
				opts.With = &ResourceMonitorWith{
					Triggers: []TriggerDefinition{
						{Threshold: 50, TriggerAction: TriggerActionNotify},
					},
				}
			},
			NewError("due to Snowflake limitations you cannot create Resource Monitor with only triggers set"),
		).
		withExpectedSqlf(
			case_ResourceMonitors_sql_Create_basic,
			`CREATE RESOURCE MONITOR %s`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ResourceMonitors_sql_Create_all,
			func(opts *CreateResourceMonitorOptions) {
				opts.OrReplace = new(true)
				opts.With = &ResourceMonitorWith{
					CreditQuota:    new(100),
					Frequency:      new(FrequencyMonthly),
					StartTimestamp: &startTimestamp,
					EndTimestamp:   &endTimestamp,
					NotifyUsers:    &NotifyUsers{Users: []NotifiedUser{{Name: firstUser}, {Name: secondUser}}},
					Triggers: []TriggerDefinition{
						{Threshold: 50, TriggerAction: TriggerActionSuspendImmediate},
						{Threshold: 100, TriggerAction: TriggerActionNotify},
					},
				}
			},
			`CREATE OR REPLACE RESOURCE MONITOR %s WITH CREDIT_QUOTA = 100 FREQUENCY = MONTHLY START_TIMESTAMP = 'IMMIEDIATELY' END_TIMESTAMP = '%s' NOTIFY_USERS = (%s, %s) TRIGGERS ON 50 PERCENT DO SUSPEND_IMMEDIATE ON 100 PERCENT DO NOTIFY`,
			id.FullyQualifiedName(),
			endTimestamp,
			firstUser.FullyQualifiedName(),
			secondUser.FullyQualifiedName(),
		)

	resourceMonitorsTests.Alter.
		withAdditionalValidationCase(
			"validation_Alter_Set_frequencyWithoutStartTimestamp",
			func(opts *AlterResourceMonitorOptions) {
				opts.Set = &ResourceMonitorSet{Frequency: new(FrequencyMonthly)}
			},
			NewError("must specify frequency and start time together"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_startTimestampWithoutFrequency",
			func(opts *AlterResourceMonitorOptions) {
				opts.Set = &ResourceMonitorSet{StartTimestamp: &alterStartTimestamp}
			},
			NewError("must specify frequency and start time together"),
		).
		withModifyAndExpectedSqlf(
			case_ResourceMonitors_sql_Alter_basic,
			func(opts *AlterResourceMonitorOptions) {
				opts.Set = &ResourceMonitorSet{CreditQuota: new(50)}
			},
			`ALTER RESOURCE MONITOR %s SET CREDIT_QUOTA = 50`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ResourceMonitors_sql_Alter_all,
			func(opts *AlterResourceMonitorOptions) {
				opts.IfExists = new(true)
				opts.Set = &ResourceMonitorSet{
					CreditQuota:    new(50),
					Frequency:      new(FrequencyYearly),
					StartTimestamp: &alterStartTimestamp,
					EndTimestamp:   &endTimestamp,
					NotifyUsers:    &NotifyUsers{Users: []NotifiedUser{{Name: user1}, {Name: user2}}},
				}
				opts.Unset = &ResourceMonitorUnset{
					CreditQuota:  new(true),
					EndTimestamp: new(true),
					NotifyUsers:  new(true),
				}
				opts.Triggers = []TriggerDefinition{
					{Threshold: 50, TriggerAction: TriggerActionSuspendImmediate},
					{Threshold: 100, TriggerAction: TriggerActionNotify},
				}
			},
			`ALTER RESOURCE MONITOR IF EXISTS %s SET CREDIT_QUOTA = 50 FREQUENCY = YEARLY START_TIMESTAMP = '%s' END_TIMESTAMP = '%s' NOTIFY_USERS = (%s, %s) SET CREDIT_QUOTA = null END_TIMESTAMP = null NOTIFY_USERS = () TRIGGERS ON 50 PERCENT DO SUSPEND_IMMEDIATE ON 100 PERCENT DO NOTIFY`,
			id.FullyQualifiedName(),
			alterStartTimestamp,
			endTimestamp,
			user1.FullyQualifiedName(),
			user2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setNotifyUsers",
			func(opts *AlterResourceMonitorOptions) {
				opts.Set = &ResourceMonitorSet{
					NotifyUsers: &NotifyUsers{Users: []NotifiedUser{{Name: user1}, {Name: user2}}},
				}
			},
			`ALTER RESOURCE MONITOR %s SET NOTIFY_USERS = (%s, %s)`,
			id.FullyQualifiedName(),
			user1.FullyQualifiedName(),
			user2.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_setMultiple",
			func(opts *AlterResourceMonitorOptions) {
				opts.Set = &ResourceMonitorSet{
					CreditQuota:    new(50),
					Frequency:      new(FrequencyYearly),
					StartTimestamp: &alterStartTimestamp,
				}
			},
			`ALTER RESOURCE MONITOR %s SET CREDIT_QUOTA = 50 FREQUENCY = YEARLY START_TIMESTAMP = '%s'`,
			id.FullyQualifiedName(),
			alterStartTimestamp,
		).
		withAdditionalSqlCasef(
			"sql_Alter_unset",
			func(opts *AlterResourceMonitorOptions) {
				opts.Unset = &ResourceMonitorUnset{
					CreditQuota:  new(true),
					EndTimestamp: new(true),
				}
			},
			`ALTER RESOURCE MONITOR %s SET CREDIT_QUOTA = null END_TIMESTAMP = null`,
			id.FullyQualifiedName(),
		)

	resourceMonitorsTests.Drop.
		withExpectedSqlf(
			case_ResourceMonitors_sql_Drop_basic,
			`DROP RESOURCE MONITOR %s`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ResourceMonitors_sql_Drop_all,
			func(opts *DropResourceMonitorOptions) { opts.IfExists = new(true) },
			`DROP RESOURCE MONITOR IF EXISTS %s`,
			id.FullyQualifiedName(),
		)

	resourceMonitorsTests.Show.
		withExpectedSql(case_ResourceMonitors_sql_Show_basic, `SHOW RESOURCE MONITORS`).
		withModifyAndExpectedSqlf(
			case_ResourceMonitors_sql_Show_all,
			func(opts *ShowResourceMonitorOptions) {
				opts.Like = &Like{Pattern: new("pattern")}
			},
			`SHOW RESOURCE MONITORS LIKE 'pattern'`,
		).
		withModifyAndExpectedSqlf(
			case_ResourceMonitors_sql_Show_Like,
			func(opts *ShowResourceMonitorOptions) { opts.Like = &Like{Pattern: new("pattern")} },
			`SHOW RESOURCE MONITORS LIKE 'pattern'`,
		)
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
