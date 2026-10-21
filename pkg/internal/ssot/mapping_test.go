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
}

func Test_ConstructorGet(t *testing.T) {
	minNodes := Attribute{Name: "min_nodes", Type: AttrInt, Required: true}
	get, ok := minNodes.ConstructorGet()
	require.True(t, ok)
	require.Equal(t, `d.Get("min_nodes").(int)`, get)

	family := NewEnum("ComputePoolInstanceFamily", "ComputePoolInstanceFamilies", "CPU_X64_XS")
	instanceFamily := Attribute{Name: "instance_family", Type: AttrEnum, Required: true, Enum: family}
	get, ok = instanceFamily.ConstructorGet()
	require.True(t, ok)
	require.Equal(t, `d.Get("instance_family").(string)`, get)

	comment := Attribute{Name: "comment", Type: AttrString}
	_, ok = comment.ConstructorGet()
	require.False(t, ok)

	mapped := Attribute{Name: "comment", Type: AttrString, CreateMapping: StringAttributeCreateBuilder}
	_, ok = mapped.ConstructorGet()
	require.False(t, ok)

	require.PanicsWithValue(t, `required attribute "items" has no CreateMapping`, func() {
		items := Attribute{Name: "items", Type: AttrList, Required: true}
		items.ConstructorGet()
	})
	require.PanicsWithValue(t, `constructor attribute "instance_family" has no Enum`, func() {
		missingEnum := Attribute{Name: "instance_family", Type: AttrEnum, Required: true}
		missingEnum.ConstructorGet()
	})
}

func Test_UpdateCall(t *testing.T) {
	call := IntAttributeUpdateSetOnlyBuilder.UpdateCall("ComputePool", Attribute{Name: "min_nodes"})
	require.Equal(t, UpdateCall{Helper: "intAttributeUpdateSetOnlyBuilder", Target: "set.WithMinNodes"}, call)

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
