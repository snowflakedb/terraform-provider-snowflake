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

	call = ListValueCreate("toComputePoolBackupInstanceFamily").CreateCall("ComputePool", Attribute{Name: "backup_instance_families"})
	require.Equal(t, CreateCall{
		Helper: "listValueCreate",
		Target: "&request.BackupInstanceFamilies, toComputePoolBackupInstanceFamily",
	}, call)

	require.PanicsWithValue(t, `object "ComputePool": attribute "items" ListValueCreate needs a mapper function`, func() {
		ListValueCreate("").CreateCall("ComputePool", Attribute{Name: "items"})
	})
}

func Test_ConstructorGet(t *testing.T) {
	get, ok := Attribute{Name: "min_nodes", Type: AttrInt, Required: true}.ConstructorGet()
	require.True(t, ok)
	require.Equal(t, `d.Get("min_nodes").(int)`, get)

	family := NewEnum("ComputePoolInstanceFamily", "ComputePoolInstanceFamilies", "CPU_X64_XS")
	get, ok = Attribute{Name: "instance_family", Type: AttrEnum, Required: true, Enum: family}.ConstructorGet()
	require.True(t, ok)
	require.Equal(t, `d.Get("instance_family").(string)`, get)

	_, ok = Attribute{Name: "comment", Type: AttrString}.ConstructorGet()
	require.False(t, ok)

	_, ok = Attribute{Name: "comment", Type: AttrString, CreateMapping: StringAttributeCreateBuilder}.ConstructorGet()
	require.False(t, ok)

	require.PanicsWithValue(t, `required attribute "items" has no CreateMapping`, func() {
		Attribute{Name: "items", Type: AttrList, Required: true}.ConstructorGet()
	})
	require.PanicsWithValue(t, `constructor attribute "instance_family" has no Enum`, func() {
		Attribute{Name: "instance_family", Type: AttrEnum, Required: true}.ConstructorGet()
	})
}
