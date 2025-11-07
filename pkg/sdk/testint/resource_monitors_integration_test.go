//go:build non_account_level_tests

package testint

import (
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ResourceMonitorsShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	resourceMonitorTest, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)

	t.Run("with like", func(t *testing.T) {
		showOptions := &sdk.ShowResourceMonitorOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(resourceMonitorTest.Name),
			},
		}
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, resourceMonitors, *resourceMonitorTest)
		assert.Len(t, resourceMonitors, 1)
	})

	t.Run("when searching a non-existent resource monitor", func(t *testing.T) {
		showOptions := &sdk.ShowResourceMonitorOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		resourceMonitors, err := client.ResourceMonitors.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Empty(t, resourceMonitors)
	})

	t.Run("show by id", func(t *testing.T) {
		resourceMonitor, err := client.ResourceMonitors.ShowByID(ctx, resourceMonitorTest.ID())
		require.NoError(t, err)
		assert.Equal(t, *resourceMonitor, *resourceMonitorTest)
	})

	t.Run("show by id when searching a non-existent resource monitor", func(t *testing.T) {
		resourceMonitor, err := client.ResourceMonitors.ShowByID(ctx, NonExistingAccountObjectIdentifier)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
		assert.Nil(t, resourceMonitor)
	})
}

func TestInt_ResourceMonitorCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("test complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		name := id.Name()
		frequency := sdk.FrequencyMonthly
		startTimeStamp := "IMMEDIATELY"
		creditQuota := 100
		endTimeStamp := time.Now().Add(24 * 10 * time.Hour).Format("2006-01-02 15:04")
		triggers := []sdk.TriggerDefinition{
			{
				Threshold:     30,
				TriggerAction: sdk.TriggerActionSuspend,
			},
			{
				Threshold:     50,
				TriggerAction: sdk.TriggerActionSuspendImmediate,
			},
			{
				Threshold:     100,
				TriggerAction: sdk.TriggerActionNotify,
			},
		}

		err := client.ResourceMonitors.Create(ctx, id, &sdk.CreateResourceMonitorOptions{
			OrReplace: sdk.Bool(true),
			With: &sdk.ResourceMonitorWith{
				Frequency:      &frequency,
				CreditQuota:    &creditQuota,
				StartTimestamp: &startTimeStamp,
				EndTimestamp:   &endTimeStamp,
				// Users' emails need to be verified in order to use them for notification
				NotifyUsers: nil,
				Triggers:    triggers,
			},
		})
		require.NoError(t, err)

		t.Cleanup(testClientHelper().ResourceMonitor.DropResourceMonitorFunc(t, id))

		assertThatObject(t,
			objectassert.ResourceMonitor(t, id).
				HasName(name).
				HasFrequency(frequency).
				HasCreditQuota(float64(creditQuota)).
				HasNonEmptyStartTime().
				HasNonEmptyEndTime().
				HasNotifyAt(100).
				HasSuspendAt(30).
				HasSuspendImmediateAt(50),
		)
	})

	t.Run("validate: only one suspend trigger", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ResourceMonitors.Create(ctx, id, &sdk.CreateResourceMonitorOptions{
			With: &sdk.ResourceMonitorWith{
				CreditQuota: sdk.Int(100),
				Triggers: []sdk.TriggerDefinition{
					{
						Threshold:     30,
						TriggerAction: sdk.TriggerActionSuspend,
					},
					{
						Threshold:     50,
						TriggerAction: sdk.TriggerActionSuspend,
					},
				},
			},
		})
		require.ErrorContains(t, err, "A resource monitor can have at most one suspend trigger.")
	})

	t.Run("validate: only one suspend immediate trigger", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.ResourceMonitors.Create(ctx, id, &sdk.CreateResourceMonitorOptions{
			With: &sdk.ResourceMonitorWith{
				CreditQuota: sdk.Int(100),
				Triggers: []sdk.TriggerDefinition{
					{
						Threshold:     30,
						TriggerAction: sdk.TriggerActionSuspendImmediate,
					},
					{
						Threshold:     50,
						TriggerAction: sdk.TriggerActionSuspendImmediate,
					},
				},
			},
		})
		require.ErrorContains(t, err, "A resource monitor can have at most one suspend_immediate trigger.")
	})

	t.Run("test no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		name := id.Name()

		err := client.ResourceMonitors.Create(ctx, id, nil)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().ResourceMonitor.DropResourceMonitorFunc(t, id))

		assertThatObject(t,
			objectassert.ResourceMonitor(t, id).
				HasName(name).
				HasFrequency(sdk.FrequencyMonthly).
				HasNonEmptyStartTime().
				HasCreditQuota(0).
				HasEndTime("").
				HasNotifyUsers().
				HasNotifyAt().
				HasSuspendAtNil().
				HasSuspendImmediateAtNil(),
		)
	})
}

func TestInt_ResourceMonitorAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when adding a new trigger", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitorCleanup)

		var oldNotifyTriggers []sdk.TriggerDefinition
		for _, threshold := range resourceMonitor.NotifyAt {
			oldNotifyTriggers = append(oldNotifyTriggers, sdk.TriggerDefinition{Threshold: threshold, TriggerAction: sdk.TriggerActionNotify})
		}

		newTriggers := oldNotifyTriggers
		newTriggers = append(newTriggers, sdk.TriggerDefinition{Threshold: *resourceMonitor.SuspendAt, TriggerAction: sdk.TriggerActionSuspend})
		newTriggers = append(newTriggers, sdk.TriggerDefinition{Threshold: *resourceMonitor.SuspendImmediateAt, TriggerAction: sdk.TriggerActionSuspendImmediate})
		newTriggers = append(newTriggers, sdk.TriggerDefinition{Threshold: 30, TriggerAction: sdk.TriggerActionNotify})
		alterOptions := &sdk.AlterResourceMonitorOptions{
			Triggers: newTriggers,
		}
		err := client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), alterOptions)
		require.NoError(t, err)

		resourceMonitor, err = client.ResourceMonitors.ShowByID(ctx, resourceMonitor.ID())
		require.NoError(t, err)

		var newNotifyTriggers []sdk.TriggerDefinition
		for _, threshold := range resourceMonitor.NotifyAt {
			newNotifyTriggers = append(newNotifyTriggers, sdk.TriggerDefinition{Threshold: threshold, TriggerAction: sdk.TriggerActionNotify})
		}

		var allTriggers []sdk.TriggerDefinition
		allTriggers = append(allTriggers, newNotifyTriggers...)
		allTriggers = append(allTriggers, sdk.TriggerDefinition{Threshold: *resourceMonitor.SuspendAt, TriggerAction: sdk.TriggerActionSuspend})
		allTriggers = append(allTriggers, sdk.TriggerDefinition{Threshold: *resourceMonitor.SuspendImmediateAt, TriggerAction: sdk.TriggerActionSuspendImmediate})

		assert.ElementsMatch(t, newTriggers, allTriggers)
	})

	t.Run("when setting and unsetting credit quota", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitorCleanup)

		creditQuota := 100

		err := client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), &sdk.AlterResourceMonitorOptions{
			Set: &sdk.ResourceMonitorSet{
				CreditQuota: &creditQuota,
			},
		})
		require.NoError(t, err)

		resourceMonitor, err = client.ResourceMonitors.ShowByID(ctx, resourceMonitor.ID())
		require.NoError(t, err)
		assert.Equal(t, creditQuota, int(resourceMonitor.CreditQuota))

		err = client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), &sdk.AlterResourceMonitorOptions{
			Unset: &sdk.ResourceMonitorUnset{
				CreditQuota: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		resourceMonitor, err = client.ResourceMonitors.ShowByID(ctx, resourceMonitor.ID())
		require.NoError(t, err)
		assert.InDelta(t, float64(0), resourceMonitor.CreditQuota, testvars.FloatEpsilon)
	})

	t.Run("when changing notify users", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitorCleanup)

		err := client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), &sdk.AlterResourceMonitorOptions{
			Set: &sdk.ResourceMonitorSet{
				NotifyUsers: &sdk.NotifyUsers{
					Users: []sdk.NotifiedUser{{Name: sdk.NewAccountObjectIdentifier("JAN_CIESLAK")}},
				},
			},
		})
		require.NoError(t, err)

		resourceMonitor, err = client.ResourceMonitors.ShowByID(ctx, resourceMonitor.ID())
		require.NoError(t, err)
		assert.Len(t, resourceMonitor.NotifyUsers, 1)
		assert.Equal(t, "JAN_CIESLAK", resourceMonitor.NotifyUsers[0])

		err = client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), &sdk.AlterResourceMonitorOptions{
			Unset: &sdk.ResourceMonitorUnset{
				NotifyUsers: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		resourceMonitor, err = client.ResourceMonitors.ShowByID(ctx, resourceMonitor.ID())
		require.NoError(t, err)
		assert.Empty(t, resourceMonitor.NotifyUsers)
	})

	t.Run("when changing scheduling info", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitorCleanup)

		frequency := sdk.FrequencyNever
		startTimeStamp := "2050-01-01 12:34"
		endTimeStamp := "2051-01-01 12:34"

		err := client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), &sdk.AlterResourceMonitorOptions{
			Set: &sdk.ResourceMonitorSet{
				Frequency:      &frequency,
				StartTimestamp: &startTimeStamp,
				EndTimestamp:   &endTimeStamp,
			},
		})
		require.NoError(t, err)

		resourceMonitor, err = client.ResourceMonitors.ShowByID(ctx, resourceMonitor.ID())
		require.NoError(t, err)

		assert.Equal(t, frequency, resourceMonitor.Frequency)
		assert.NotEmpty(t, resourceMonitor.StartTime)
		assert.NotEmpty(t, resourceMonitor.EndTime)

		err = client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), &sdk.AlterResourceMonitorOptions{
			Unset: &sdk.ResourceMonitorUnset{
				EndTimestamp: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		resourceMonitor, err = client.ResourceMonitors.ShowByID(ctx, resourceMonitor.ID())
		require.NoError(t, err)

		assert.NotEmpty(t, resourceMonitor.StartTime)
		assert.Empty(t, resourceMonitor.EndTime)
	})

	t.Run("all options together", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitorCleanup)

		newTriggers := make([]sdk.TriggerDefinition, 0)
		newTriggers = append(newTriggers, sdk.TriggerDefinition{Threshold: 30, TriggerAction: sdk.TriggerActionNotify})

		creditQuota := 100
		err := client.ResourceMonitors.Alter(ctx, resourceMonitor.ID(), &sdk.AlterResourceMonitorOptions{
			Set: &sdk.ResourceMonitorSet{
				CreditQuota: &creditQuota,
				NotifyUsers: &sdk.NotifyUsers{
					Users: []sdk.NotifiedUser{{Name: sdk.NewAccountObjectIdentifier("JAN_CIESLAK")}},
				},
			},
			Triggers: newTriggers,
		})
		require.NoError(t, err)

		assertThatObject(t,
			objectassert.ResourceMonitor(t, resourceMonitor.ID()).
				HasCreditQuota(float64(creditQuota)).
				HasNotifyUsers("JAN_CIESLAK").
				HasNotifyAt(30),
		)
	})
}

func TestInt_ResourceMonitorDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when resource monitor exists", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitorCleanup)

		err := client.ResourceMonitors.Drop(ctx, resourceMonitor.ID(), nil)
		require.NoError(t, err)

		_, err = client.ResourceMonitors.ShowByID(ctx, resourceMonitor.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("when resource monitor does not exist", func(t *testing.T) {
		err := client.ResourceMonitors.Drop(ctx, NonExistingAccountObjectIdentifier, nil)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
