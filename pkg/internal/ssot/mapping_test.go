package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CreateCall(t *testing.T) {
	call := AccountObjectIdentifierAttributeCreateBuilder.CreateCall("ComputePool", Attribute{Name: "for_application"})
	require.Equal(t, CreateCall{
		Helper: "accountObjectIdentifierAttributeCreateBuilder",
		Target: "request.WithForApplication",
	}, call)

	call = ListValueCreate("toComputePoolBackupInstanceFamilies").CreateCall("ComputePool", Attribute{Name: "backup_instance_families"})
	require.Equal(t, CreateCall{
		Helper: "listValueCreate",
		Target: "&request.BackupInstanceFamilies, toComputePoolBackupInstanceFamilies",
	}, call)

	require.PanicsWithValue(t, `object "ComputePool": attribute "items" ListValueCreate needs a mapper function`, func() {
		ListValueCreate("").CreateCall("ComputePool", Attribute{Name: "items"})
	})

	call = AttributeMappedValueCreate("ToSecurityIntegrationsRequest").CreateCall("AuthenticationPolicy", Attribute{Name: "security_integrations"})
	require.Equal(t, CreateCall{
		Helper: "attributeMappedValueCreateBuilder",
		Target: "request.WithSecurityIntegrations, ToSecurityIntegrationsRequest",
	}, call)

	require.PanicsWithValue(t, `object "AuthenticationPolicy": attribute "security_integrations" attributeMappedValueCreateBuilder needs a mapper function`, func() {
		AttributeMappedValueCreate("").CreateCall("AuthenticationPolicy", Attribute{Name: "security_integrations"})
	})

	call = SetMappedValueCreate("toPostgresMirrorName").CreateCall("PostgresMirror", Attribute{Name: "postgres_tables"})
	require.Equal(t, CreateCall{
		Helper: "setMappedValueCreateBuilder",
		Target: "request.WithPostgresTables, toPostgresMirrorName",
	}, call)

	require.PanicsWithValue(t, `object "PostgresMirror": attribute "postgres_tables" SetMappedValueCreate needs a mapper function`, func() {
		SetMappedValueCreate("").CreateCall("PostgresMirror", Attribute{Name: "postgres_tables"})
	})
}

func Test_UpdateCall(t *testing.T) {
	call := IntAttributeUpdateSetOnlyBuilder.UpdateCall("ComputePool", Attribute{Name: "min_nodes"})
	require.Equal(t, UpdateCall{Helper: "intAttributeUpdateSetOnlyBuilder", Set: true}, call)

	call = BooleanAttributeUpdateSetOnlyBuilder.UpdateCall("Table", Attribute{Name: "change_tracking"})
	require.Equal(t, UpdateCall{Helper: "booleanAttributeUpdateSetOnlyBuilder", Set: true}, call)

	call = BooleanStringAttributeUpdateBuilder.UpdateCall("ComputePool", Attribute{Name: "auto_resume"})
	require.Equal(t, UpdateCall{
		Helper: "booleanStringAttributeUpdateBuilder",
		Set:    true,
		Unset:  true,
	}, call)

	call = ListValueUpdate("toComputePoolBackupInstanceFamilies").UpdateCall("ComputePool", Attribute{Name: "backup_instance_families"})
	require.Equal(t, UpdateCall{
		Helper: "listValueUpdateBuilder",
		Set:    true,
		Unset:  true,
		Mapper: "toComputePoolBackupInstanceFamilies",
	}, call)

	require.PanicsWithValue(t, `object "ComputePool": attribute "items" ListValueUpdate needs a mapper function`, func() {
		ListValueUpdate("").UpdateCall("ComputePool", Attribute{Name: "items"})
	})
}

func Test_UpdateCall_listValueUpdateSetOnly(t *testing.T) {
	call := ListValueUpdateSetOnly("toWidgetObjectTypes").UpdateCall("Widget", Attribute{Name: "object_types"})
	require.Equal(t, UpdateCall{Helper: "listValueUpdateSetOnlyBuilder", Set: true, Mapper: "toWidgetObjectTypes"}, call)
}

func Test_UpdateCall_fieldsStyle(t *testing.T) {
	call := AccountObjectIdentifierAttributeUpdate.UpdateCall("ReplicationGroup", Attribute{Name: "error_integration"})
	require.Equal(t, UpdateCall{Helper: "accountObjectIdentifierAttributeUpdate", Set: true, Unset: true, Fields: true}, call)

	call = AttributeMappedValueUpdate("toBackupPolicySchedule").UpdateCall("BackupPolicy", Attribute{Name: "schedule"})
	require.Equal(t, UpdateCall{
		Helper: "attributeMappedValueUpdate",
		Set:    true,
		Unset:  true,
		Fields: true,
		Mapper: "toBackupPolicySchedule",
	}, call)

	require.PanicsWithValue(t, `object "BackupPolicy": attribute "schedule" AttributeMappedValueUpdate needs a mapper function`, func() {
		AttributeMappedValueUpdate("").UpdateCall("BackupPolicy", Attribute{Name: "schedule"})
	})
}

func Test_UpdatesWholeValue(t *testing.T) {
	require.True(t, UpdatesWholeValue(AttributeMappedValueUpdate("toBackupPolicySchedule")))
	require.False(t, UpdatesWholeValue(StringAttributeUpdateBuilder))
	require.False(t, UpdatesWholeValue(AccountObjectIdentifierAttributeUpdate))
	require.False(t, UpdatesWholeValue(nil))
}

func Test_UpdateCall_booleanStringSetOnly(t *testing.T) {
	require.Equal(t,
		UpdateCall{Helper: "booleanStringAttributeUpdateSetOnly", Set: true, Fields: true},
		BooleanStringAttributeUpdateSetOnly.UpdateCall("Widget", Attribute{Name: "enabled"}))
}
