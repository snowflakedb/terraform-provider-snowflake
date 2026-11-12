package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		expected []map[string]any
		err      string
	}{
		{name: "not set (nil)", input: nil},
		{name: "not set (empty)", input: sdk.String("")},
		{name: "hours", input: sdk.String("2 HOUR"), expected: []map[string]any{{"hours": 2}}},
		{name: "hours (plural)", input: sdk.String("2 HOURS"), expected: []map[string]any{{"hours": 2}}},
		{name: "minutes", input: sdk.String("60 MINUTE"), expected: []map[string]any{{"minutes": 60}}},
		{name: "cron keeps the time zone casing", input: sdk.String("USING CRON 0 0 * * * America/Los_Angeles"), expected: []map[string]any{{"using_cron": "0 0 * * * America/Los_Angeles"}}},
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
			schedule, err := backupPolicyScheduleFromConfig([]any{state[0]})
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
		requests, err := backupPolicyApplySetUnset(d)
		require.NoError(t, err)
		require.Equal(t, sdk.NewBackupPolicySetRequest().WithSchedule("USING CRON 0 0 * * * UTC").WithExpireAfterDays(7).WithComment("c"), requests.set)
		require.Equal(t, sdk.NewBackupPolicyUnsetRequest(), requests.unset)
	})

	t.Run("removing the schedule unsets it in its own request", func(t *testing.T) {
		state := &terraform.InstanceState{
			ID:         "id",
			Attributes: map[string]string{"schedule.#": "1", "schedule.0.hours": "2", "schedule.0.minutes": "0", "schedule.0.using_cron": ""},
		}
		diff, err := schema.InternalMap(backupPolicySchema).Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]any{}), nil, nil, true)
		require.NoError(t, err)
		d, err := schema.InternalMap(backupPolicySchema).Data(state, diff)
		require.NoError(t, err)
		requests, err := backupPolicyApplySetUnset(d)
		require.NoError(t, err)
		require.Equal(t, sdk.NewBackupPolicySetRequest(), requests.set)
		require.Equal(t, sdk.NewBackupPolicyUnsetRequest(), requests.unset)
		require.Equal(t, sdk.NewBackupPolicyUnsetRequest().WithSchedule(true), requests.unsetSchedule)
	})

	t.Run("invalid schedule fails the update", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, backupPolicySchema, map[string]any{
			"schedule": []any{map[string]any{}},
		})
		_, err := backupPolicyApplySetUnset(d)
		require.ErrorContains(t, err, "one of hours, minutes, or using_cron")
	})
}

func Test_backupPolicySetFieldsNotSetByReadExt(t *testing.T) {
	id := sdk.NewSchemaObjectIdentifier("db", "sch", "bp")
	stateOf := func(t *testing.T, backupPolicy *sdk.BackupPolicy, importPermissiveMode bool) map[string]string {
		t.Helper()
		d := (&schema.Resource{Schema: backupPolicySchema}).Data(nil)
		d.SetId("id")
		require.NoError(t, backupPolicySetFieldsNotSetByReadExt(d, id, backupPolicy, importPermissiveMode))
		return d.State().Attributes
	}

	for _, permissive := range []bool{false, true} {
		t.Run(fmt.Sprintf("a policy without an expiration, permissive: %t", permissive), func(t *testing.T) {
			attributes := stateOf(t, &sdk.BackupPolicy{Schedule: sdk.String("2 HOUR")}, permissive)
			require.NotContains(t, attributes, "expire_after_days")
			require.Equal(t, "db", attributes["database"])
			require.Equal(t, "sch", attributes["schema"])
			require.Equal(t, "bp", attributes["name"])
		})
	}

	t.Run("enforcing mode writes the expiration", func(t *testing.T) {
		attributes := stateOf(t, &sdk.BackupPolicy{ExpireAfterDays: sdk.Int(7), HasRetentionLock: true}, false)
		require.Equal(t, "7", attributes["expire_after_days"])
		require.Equal(t, "true", attributes["with_retention_lock"])
	})

	t.Run("permissive mode writes the unset values", func(t *testing.T) {
		attributes := stateOf(t, &sdk.BackupPolicy{ExpireAfterDays: sdk.Int(7), HasRetentionLock: true}, true)
		require.Equal(t, "0", attributes["expire_after_days"])
		require.Equal(t, BooleanDefault, attributes["with_retention_lock"])
	})
}

func Test_toBackupPolicySchedule(t *testing.T) {
	block := func(fields map[string]any) []any { return []any{fields} }

	schedule, err := toBackupPolicySchedule(block(map[string]any{"hours": 2, "minutes": 0, "using_cron": ""}))
	require.NoError(t, err)
	require.Equal(t, "2 HOUR", schedule)

	schedule, err = toBackupPolicySchedule(block(map[string]any{"hours": 0, "minutes": 90, "using_cron": ""}))
	require.NoError(t, err)
	require.Equal(t, "90 MINUTE", schedule)

	schedule, err = toBackupPolicySchedule(block(map[string]any{"hours": 0, "minutes": 0, "using_cron": "0 * * * * UTC"}))
	require.NoError(t, err)
	require.Equal(t, "USING CRON 0 * * * * UTC", schedule)

	_, err = toBackupPolicySchedule(block(map[string]any{"hours": 0, "minutes": 0, "using_cron": ""}))
	require.ErrorContains(t, err, "one of hours, minutes, or using_cron field should be set")

	_, err = toBackupPolicySchedule(nil)
	require.ErrorContains(t, err, "one of hours, minutes, or using_cron field should be set")
}
