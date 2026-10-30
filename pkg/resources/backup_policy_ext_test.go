package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func Test_backupPolicyScheduleFromConfig(t *testing.T) {
	withFields := func(hours, minutes int, cron string) any {
		return []any{map[string]any{"hours": hours, "minutes": minutes, "using_cron": cron}}
	}

	tests := []struct {
		name     string
		input    any
		expected *string
		err      string
	}{
		{name: "no block (nil)", input: nil},
		{name: "no block (empty)", input: []any{}},
		{name: "hours", input: withFields(2, 0, ""), expected: sdk.String("2 HOUR")},
		{name: "minutes", input: withFields(0, 90, ""), expected: sdk.String("90 MINUTE")},
		{name: "cron", input: withFields(0, 0, "0 0 * * * UTC"), expected: sdk.String("USING CRON 0 0 * * * UTC")},
		{name: "empty block", input: withFields(0, 0, ""), err: "one of hours, minutes, or using_cron"},
		{name: "block with a nil element", input: []any{nil}, err: "one of hours, minutes, or using_cron"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			schedule, err := backupPolicyScheduleFromConfig(tc.input)
			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, schedule)
		})
	}
}

func Test_backupPolicyScheduleToState(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected []any
		err      string
	}{
		{name: "not set (nil)", input: nil},
		{name: "not set (empty)", input: sdk.String("")},
		{name: "hours", input: sdk.String("2 HOUR"), expected: []any{map[string]any{"hours": 2}}},
		{name: "hours (plural)", input: sdk.String("2 HOURS"), expected: []any{map[string]any{"hours": 2}}},
		{name: "minutes", input: sdk.String("60 MINUTE"), expected: []any{map[string]any{"minutes": 60}}},
		{name: "cron keeps the time zone casing", input: sdk.String("USING CRON 0 0 * * * America/Los_Angeles"), expected: []any{map[string]any{"using_cron": "0 0 * * * America/Los_Angeles"}}},
		{name: "seconds are not supported", input: sdk.String("30 SECOND"), err: "unsupported backup policy schedule"},
		{name: "invalid", input: sdk.String("every day"), err: "EVERY"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state, err := backupPolicyScheduleToState(tc.input)
			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, state)
		})
	}
}

func Test_backupPolicyScheduleRoundTrip(t *testing.T) {
	for _, value := range []string{"2 HOUR", "60 MINUTE", "USING CRON 0 0 * * * UTC"} {
		t.Run(value, func(t *testing.T) {
			state, err := backupPolicyScheduleToState(sdk.String(value))
			require.NoError(t, err)
			schedule, err := backupPolicyScheduleFromConfig(state)
			require.NoError(t, err)
			require.Equal(t, sdk.String(value), schedule)
		})
	}
}

func Test_backupPolicyRequestsExt(t *testing.T) {
	id := sdk.NewSchemaObjectIdentifier("db", "sch", "bp")

	t.Run("create request holds the schedule", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, backupPolicySchema, map[string]any{
			"schedule": []any{map[string]any{"hours": 2}},
		})
		req, err := backupPolicyNewCreateRequestExt(d, id)
		require.NoError(t, err)
		require.Equal(t, sdk.NewCreateBackupPolicyRequest(id).WithSchedule("2 HOUR"), req)
	})

	t.Run("create request without a schedule", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, backupPolicySchema, map[string]any{"expire_after_days": 7})
		req, err := backupPolicyNewCreateRequestExt(d, id)
		require.NoError(t, err)
		require.Equal(t, sdk.NewCreateBackupPolicyRequest(id), req)
	})

	t.Run("set and unset", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, backupPolicySchema, map[string]any{
			"schedule":          []any{map[string]any{"using_cron": "0 0 * * * UTC"}},
			"expire_after_days": 7,
			"comment":           "c",
		})
		requests, err := backupPolicyApplySetUnsetExt(d)
		require.NoError(t, err)
		require.Equal(t, sdk.NewBackupPolicySetRequest().WithSchedule("USING CRON 0 0 * * * UTC").WithExpireAfterDays(7).WithComment("c"), requests.set)
		require.Equal(t, sdk.NewBackupPolicyUnsetRequest(), requests.unset)
	})

	t.Run("invalid schedule fails the update", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, backupPolicySchema, map[string]any{
			"schedule": []any{map[string]any{}},
		})
		_, err := backupPolicyApplySetUnsetExt(d)
		require.ErrorContains(t, err, "one of hours, minutes, or using_cron")
	})
}
