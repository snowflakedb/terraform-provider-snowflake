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
}

func Test_UpdateCall(t *testing.T) {
	call := IntAttributeUpdateSetOnlyBuilder.UpdateCall("ComputePool", Attribute{Name: "min_nodes"})
	require.Equal(t, UpdateCall{Helper: "intAttributeUpdateSetOnlyBuilder", Target: "set.WithMinNodes"}, call)

	call = BooleanAttributeUpdateSetOnlyBuilder.UpdateCall("Table", Attribute{Name: "change_tracking"})
	require.Equal(t, UpdateCall{Helper: "booleanAttributeUpdateSetOnlyBuilder", Target: "set.WithChangeTracking"}, call)

	call = BooleanStringAttributeUpdateBuilder.UpdateCall("ComputePool", Attribute{Name: "auto_resume"})
	require.Equal(t, UpdateCall{
		Helper: "booleanStringAttributeUpdateBuilder",
		Target: "set.WithAutoResume, unset.WithAutoResume",
	}, call)

	call = ListValueUpdate("toComputePoolBackupInstanceFamilies").UpdateCall("ComputePool", Attribute{Name: "backup_instance_families"})
	require.Equal(t, UpdateCall{
		Helper: "listValueUpdateBuilder",
		Target: "set.WithBackupInstanceFamilies, unset.WithBackupInstanceFamilies, toComputePoolBackupInstanceFamilies",
	}, call)

	require.PanicsWithValue(t, `object "ComputePool": attribute "items" ListValueUpdate needs a mapper function`, func() {
		ListValueUpdate("").UpdateCall("ComputePool", Attribute{Name: "items"})
	})
}
