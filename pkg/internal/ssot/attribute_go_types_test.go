package ssot

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/stretchr/testify/require"
)

func Test_AttributeGoTypes(t *testing.T) {
	family := NewEnum("PoolFamily", "PoolFamilies", "SMALL")
	stateEnum := NewEnum("PoolState", "PoolStates", "IDLE")
	obj := Object{Name: "Pool", IdType: "AccountObjectIdentifier"}
	schemaObj := Object{Name: "Item", IdType: "SchemaObjectIdentifier"}

	concrete, underlying := AttributeGoTypes(obj, Attribute{Name: "name", Type: AttrString, OptionalInShow: true})
	require.Equal(t, "string", concrete)
	require.Equal(t, "string", underlying)

	concrete, underlying = AttributeGoTypes(obj, Attribute{Name: "min_nodes", Type: AttrInt})
	require.Equal(t, "int", concrete)
	require.Equal(t, "int", underlying)

	concrete, underlying = AttributeGoTypes(obj, Attribute{Name: "auto_resume", Type: AttrBool})
	require.Equal(t, "bool", concrete)
	require.Equal(t, "bool", underlying)

	concrete, underlying = AttributeGoTypes(obj, Attribute{Name: "available", Type: AttrFloat})
	require.Equal(t, "float64", concrete)
	require.Equal(t, "float64", underlying)

	concrete, underlying = AttributeGoTypes(obj, Attribute{Name: "is_default", Type: AttrBool, StoredAsText: true})
	require.Equal(t, "bool", concrete)
	require.Equal(t, "bool", underlying)

	concrete, underlying = AttributeGoTypes(obj, Attribute{Name: "created_on", Type: AttrTime})
	require.Equal(t, "time.Time", concrete)
	require.Equal(t, "struct", underlying)

	concrete, underlying = AttributeGoTypes(obj, Attribute{Name: "state", Type: AttrEnum, Enum: stateEnum})
	require.Equal(t, "sdk.PoolState", concrete)
	require.Equal(t, "string", underlying)

	concrete, underlying = AttributeGoTypes(obj, Attribute{Name: "for_application", Type: AttrIdentifier, IdentifierType: "AccountObjectIdentifier"})
	require.Equal(t, "sdk.AccountObjectIdentifier", concrete)
	require.Equal(t, "struct", underlying)

	concrete, underlying = AttributeGoTypes(obj, Attribute{
		Name: "backup_instance_families",
		Type: AttrList,
		Elem: &Attribute{Type: AttrEnum, Enum: family},
	})
	require.Equal(t, "[]string", concrete)
	require.Equal(t, "slice", underlying)

	concrete, underlying = AttributeGoTypes(obj, Attribute{
		Name: "names",
		Type: AttrSet,
		Elem: &Attribute{Type: AttrString},
	})
	require.Equal(t, "[]string", concrete)
	require.Equal(t, "slice", underlying)

	concrete, underlying = AttributeGoTypes(schemaObj, Attribute{
		Name: "refs",
		Type: AttrList,
		Elem: &Attribute{Type: AttrIdentifier, IdentifierType: "SchemaObjectIdentifier"},
	})
	require.Equal(t, "[]sdk.SchemaObjectIdentifier", concrete)
	require.Equal(t, "slice", underlying)

	concrete, underlying = AttributeGoTypes(obj, Attribute{
		Name: "counts",
		Type: AttrList,
		Elem: &Attribute{Type: AttrInt},
	})
	require.Equal(t, "[]int", concrete)
	require.Equal(t, "slice", underlying)

	concrete, underlying = AttributeGoTypes(obj, Attribute{Name: "rest_config", Type: AttrObject})
	require.Equal(t, "sdk.PoolRestConfigDetails", concrete)
	require.Equal(t, "struct", underlying)
}

func Test_ColumnField_matchesExtractedStructField(t *testing.T) {
	stateEnum := NewEnum("PoolState", "PoolStates", "IDLE")
	obj := Object{Name: "Pool", IdType: "AccountObjectIdentifier"}
	schemaObj := Object{Name: "Item", IdType: "SchemaObjectIdentifier"}

	require.Equal(t, genhelpers.Field{Name: "Options", ConcreteType: "*string", UnderlyingType: "*string"},
		ColumnField(obj, Attribute{Name: "options", Type: AttrString, OptionalInShow: true}))
	require.Equal(t, genhelpers.Field{Name: "MinNodes", ConcreteType: "*int", UnderlyingType: "*int"},
		ColumnField(obj, Attribute{Name: "min_nodes", Type: AttrInt, OptionalInShow: true}))
	require.Equal(t, genhelpers.Field{Name: "AutoResume", ConcreteType: "*bool", UnderlyingType: "*bool"},
		ColumnField(obj, Attribute{Name: "auto_resume", Type: AttrBool, OptionalInShow: true}))
	require.Equal(t, genhelpers.Field{Name: "UpdatedOn", ConcreteType: "*time.Time", UnderlyingType: "*struct"},
		ColumnField(obj, Attribute{Name: "updated_on", Type: AttrTime, OptionalInShow: true}))
	require.Equal(t, genhelpers.Field{Name: "State", ConcreteType: "*sdk.PoolState", UnderlyingType: "*string"},
		ColumnField(obj, Attribute{Name: "state", Type: AttrEnum, Enum: stateEnum, OptionalInShow: true}))
	require.Equal(t, genhelpers.Field{Name: "Application", ConcreteType: "*sdk.AccountObjectIdentifier", UnderlyingType: "*struct"},
		ColumnField(obj, Attribute{Name: "for_application", Type: AttrIdentifier, IdentifierType: "AccountObjectIdentifier", OptionalInShow: true, ShowName: "application"}))
	require.Equal(t, genhelpers.Field{Name: "Owner", ConcreteType: "sdk.AccountObjectIdentifier", UnderlyingType: "struct"},
		ColumnField(obj, Attribute{Name: "owner", Type: AttrIdentifier, IdentifierType: "AccountObjectIdentifier"}))
	require.Equal(t, genhelpers.Field{Name: "Names", ConcreteType: "[]string", UnderlyingType: "slice"},
		ColumnField(obj, Attribute{Name: "names", Type: AttrSet, Elem: &Attribute{Type: AttrString}}))
	require.Equal(t, genhelpers.Field{Name: "Names", ConcreteType: "[]string", UnderlyingType: "slice"},
		ColumnField(obj, Attribute{Name: "names", Type: AttrList, Elem: &Attribute{Type: AttrString}}))
	require.Equal(t, genhelpers.Field{Name: "Counts", ConcreteType: "[]int", UnderlyingType: "slice"},
		ColumnField(obj, Attribute{Name: "counts", Type: AttrList, Elem: &Attribute{Type: AttrInt}}))
	require.Equal(t, genhelpers.Field{Name: "Refs", ConcreteType: "[]sdk.SchemaObjectIdentifier", UnderlyingType: "slice"},
		ColumnField(schemaObj, Attribute{Name: "refs", Type: AttrList, Elem: &Attribute{Type: AttrIdentifier, IdentifierType: "SchemaObjectIdentifier"}}))
	require.Equal(t, genhelpers.Field{Name: "Names", ConcreteType: "*[]string", UnderlyingType: "*slice"},
		ColumnField(obj, Attribute{Name: "names", Type: AttrList, OptionalInShow: true, Elem: &Attribute{Type: AttrString}}))
	require.Equal(t, genhelpers.Field{Name: "DroppedOn", ConcreteType: "time.Time", UnderlyingType: "struct"},
		ColumnField(obj, Attribute{Name: "dropped_on", Type: AttrTime, OptionalInDb: true}))
	require.Equal(t, genhelpers.Field{Name: "IsDefault", ConcreteType: "bool", UnderlyingType: "bool"},
		ColumnField(obj, Attribute{Name: "is_default", Type: AttrBool, StoredAsText: true}))
	showEnum := NewEnum("WarehouseType", "WarehouseTypes", "STANDARD")
	require.Equal(t, genhelpers.Field{Name: "Type", ConcreteType: "sdk.WarehouseType", UnderlyingType: "string"},
		ColumnField(obj, Attribute{
			Name:     "warehouse_type",
			Type:     AttrEnum,
			Enum:     NewEnum("WarehouseTypeUserSettable", "", "STANDARD"),
			ShowName: "type",
			ShowEnum: showEnum,
		}))
}

func Test_AttributeGoTypes_panics(t *testing.T) {
	obj := Object{Name: "Pool", IdType: "AccountObjectIdentifier"}

	require.Panics(t, func() {
		AttributeGoTypes(obj, Attribute{Name: "names", Type: AttrList})
	})
	require.Panics(t, func() {
		AttributeGoTypes(obj, Attribute{Name: "tag", Type: "unknown"})
	})
}
