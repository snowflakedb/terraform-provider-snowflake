package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CollectionMapperName(t *testing.T) {
	families := Attribute{Name: "backup_instance_families"}
	require.Equal(t, "toComputePoolBackupInstanceFamilies", CollectionMapperName("ComputePool", families))

	methods := Attribute{Name: "allowed_methods"}
	require.Equal(t, "toAuthenticationPolicyMfaPolicyAllowedMethods", CollectionMapperName("AuthenticationPolicyMfaPolicy", methods))
}

func Test_DeducedCreateCall(t *testing.T) {
	size := NewEnum("WidgetSize", "WidgetSizes", "SMALL")

	t.Run("boolean string", func(t *testing.T) {
		got, ok := Attribute{Name: "enabled", Type: AttrBool}.DeducedCreateCall("Widget")
		require.True(t, ok)
		require.Equal(t, "booleanStringAttributeCreateBuilder", got.Helper)
		require.Equal(t, "request.WithEnabled", got.Target)
	})
	t.Run("native bool", func(t *testing.T) {
		got, ok := Attribute{Name: "enabled", Type: AttrBool, NativeBool: true}.DeducedCreateCall("Widget")
		require.True(t, ok)
		require.Equal(t, "boolAttributeCreateBuilder", got.Helper)
	})
	t.Run("int whose min is at most zero", func(t *testing.T) {
		got, ok := Attribute{Name: "count", Type: AttrInt, Min: new(0)}.DeducedCreateCall("Widget")
		require.True(t, ok)
		require.Equal(t, "intAttributeWithSpecialDefaultCreateBuilder", got.Helper)
	})
	t.Run("int", func(t *testing.T) {
		got, ok := Attribute{Name: "count", Type: AttrInt, Min: new(1)}.DeducedCreateCall("Widget")
		require.True(t, ok)
		require.Equal(t, "intAttributeCreateBuilder", got.Helper)
	})
	t.Run("enum parser", func(t *testing.T) {
		got, ok := Attribute{Name: "size", Type: AttrEnum, Enum: size}.DeducedCreateCall("Widget")
		require.True(t, ok)
		require.Equal(t, "attributeMappedValueCreateBuilder", got.Helper)
		require.Equal(t, "request.WithSize, sdk.ToWidgetSize", got.Target)
	})
	t.Run("list mapper", func(t *testing.T) {
		got, ok := Attribute{
			Name: "kinds",
			Type: AttrList,
			Elem: &Attribute{Type: AttrEnum, Enum: size},
		}.DeducedCreateCall("Widget")
		require.True(t, ok)
		require.Equal(t, "listValueCreate", got.Helper)
		require.Equal(t, "&request.Kinds, toWidgetKinds", got.Target)
	})
	t.Run("required has no helper", func(t *testing.T) {
		_, ok := Attribute{Name: "name", Type: AttrString, Required: true}.DeducedCreateCall("Widget")
		require.False(t, ok)
	})
}

func Test_DeducedUpdateCall(t *testing.T) {
	size := NewEnum("WidgetSize", "WidgetSizes", "SMALL")

	t.Run("enum parser", func(t *testing.T) {
		got, ok := Attribute{Name: "size", Type: AttrEnum, Enum: size}.DeducedUpdateCall("Widget", false)
		require.True(t, ok)
		require.Equal(t, "attributeMappedValueUpdate", got.Helper)
		require.Equal(t, "sdk.ToWidgetSize", got.Mapper)
		require.True(t, got.Unset)
	})
	t.Run("set-only enum has no helper", func(t *testing.T) {
		_, ok := Attribute{Name: "size", Type: AttrEnum, Enum: size}.DeducedUpdateCall("Widget", true)
		require.False(t, ok)
	})
	t.Run("set-only list mapper", func(t *testing.T) {
		got, ok := Attribute{
			Name: "databases",
			Type: AttrList,
			Elem: &Attribute{Type: AttrIdentifier, IdentifierType: "AccountObjectIdentifier"},
		}.DeducedUpdateCall("Widget", true)
		require.True(t, ok)
		require.Equal(t, "listValueUpdateSetOnlyBuilder", got.Helper)
		require.Equal(t, "toWidgetDatabases", got.Mapper)
		require.False(t, got.Unset)
	})
	t.Run("list mapper", func(t *testing.T) {
		got, ok := Attribute{
			Name: "kinds",
			Type: AttrList,
			Elem: &Attribute{Type: AttrEnum, Enum: size},
		}.DeducedUpdateCall("Widget", false)
		require.True(t, ok)
		require.Equal(t, "listValueUpdateBuilder", got.Helper)
		require.Equal(t, "toWidgetKinds", got.Mapper)
		require.True(t, got.Unset)
	})
	t.Run("set-only boolean string", func(t *testing.T) {
		got, ok := Attribute{Name: "enabled", Type: AttrBool}.DeducedUpdateCall("Widget", true)
		require.True(t, ok)
		require.Equal(t, "booleanStringAttributeUpdateSetOnly", got.Helper)
		require.False(t, got.Unset)
	})
	t.Run("set-only int", func(t *testing.T) {
		got, ok := Attribute{Name: "count", Type: AttrInt}.DeducedUpdateCall("Widget", true)
		require.True(t, ok)
		require.Equal(t, "intAttributeUpdateSetOnlyBuilder", got.Helper)
	})
	t.Run("int whose min is at most zero", func(t *testing.T) {
		got, ok := Attribute{Name: "count", Type: AttrInt, Min: new(0)}.DeducedUpdateCall("Widget", false)
		require.True(t, ok)
		require.Equal(t, "intAttributeWithSpecialDefaultUpdateBuilder", got.Helper)
	})
}

func Test_IntAttributeUnsetFallbackUpdateWithZeroDefault(t *testing.T) {
	call := IntAttributeUnsetFallbackUpdateWithZeroDefault(30).UpdateCall("Widget", Attribute{Name: "interval"})
	require.Equal(t, UpdateCall{Helper: "intAttributeUnsetFallbackUpdateWithZeroDefault", Set: true, Fields: true, Mapper: "30"}, call)
}

func Test_Object_AlterOwner(t *testing.T) {
	comment := Attribute{Name: "comment", Type: AttrString}
	base := Object{Name: "Base", Alter: AlterSpec{Clauses: []AlterClause{AlterSet(&comment)}}}
	require.Equal(t, "Base", Object{Name: "Variant", Variant: &ObjectVariant{Of: &base}}.AlterOwner().Name)
	// A variant with an ALTER of its own, a base without one, and an object that is no variant own their ALTER.
	own := Object{Name: "Variant", Variant: &ObjectVariant{Of: &base}, Alter: AlterSpec{Clauses: []AlterClause{AlterSet(&comment)}}}
	require.Equal(t, "Variant", own.AlterOwner().Name)
	require.Equal(t, "Variant", Object{Name: "Variant", Variant: &ObjectVariant{Of: &Object{Name: "Base"}}}.AlterOwner().Name)
	require.Equal(t, "Base", base.AlterOwner().Name)
}
