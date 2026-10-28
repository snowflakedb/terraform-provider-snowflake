package ssot

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/stretchr/testify/require"
)

func Test_IdentityFields(t *testing.T) {
	account := Object{Name: "Pool", IdType: "AccountObjectIdentifier"}
	require.Equal(t, []string{"name"}, classifiedNames(IdentityFields(account)))

	database := Object{Name: "Child", IdType: "DatabaseObjectIdentifier"}
	require.Equal(t, []string{"database", "name"}, classifiedNames(IdentityFields(database)))

	nested := Object{Name: "Nested", IdType: "SchemaObjectIdentifier", IdentityNameForceNew: true, IdentityHierarchyForceNew: true}
	fields := IdentityFields(nested)
	require.Equal(t, []string{"database", "schema", "name"}, classifiedNames(fields))
	for _, field := range fields {
		require.True(t, field.ForceNew)
	}

	renameable := Object{Name: "Policy", IdType: "SchemaObjectIdentifier", IdentityHierarchyForceNew: true}
	fields = IdentityFields(renameable)
	require.False(t, fields[2].ForceNew)
	require.Equal(t, "name", fields[2].Name)
	require.True(t, fields[0].ForceNew)
	require.True(t, fields[1].ForceNew)
}

func Test_ClassifyAttribute_storageAndSdkType(t *testing.T) {
	family := NewEnum("PoolFamily", "PoolFamilies", "SMALL")
	obj := Object{Name: "Pool", IdType: "AccountObjectIdentifier"}

	optionalBool := ClassifyAttribute(obj, Attribute{Name: "auto_resume", Type: AttrBool})
	require.Equal(t, AttrString, optionalBool.Storage)
	require.Equal(t, "bool", optionalBool.SdkType)
	require.Equal(t, "bool", optionalBool.SdkKind)

	nativeBool := ClassifyAttribute(obj, Attribute{Name: "initially_suspended", Type: AttrBool, NativeBool: true})
	require.Equal(t, AttrBool, nativeBool.Storage)
	require.Empty(t, nativeBool.SdkType)
	require.False(t, nativeBool.Required)

	requiredBool := ClassifyAttribute(obj, Attribute{Name: "enabled", Type: AttrBool, Required: true})
	require.Equal(t, AttrBool, requiredBool.Storage)
	require.Empty(t, requiredBool.SdkType)

	enumField := ClassifyAttribute(obj, Attribute{Name: "instance_family", Type: AttrEnum, Enum: family, Required: true})
	require.Equal(t, AttrString, enumField.Storage)
	require.Equal(t, "sdk.PoolFamily", enumField.SdkType)
	require.Equal(t, "enum", enumField.SdkKind)

	narrow := NewEnum("WarehouseTypeUserSettable", "", "STANDARD").WithSdkType("WarehouseType")
	narrowField := ClassifyAttribute(obj, Attribute{Name: "warehouse_type", Type: AttrEnum, Enum: narrow})
	require.Equal(t, "sdk.WarehouseType", narrowField.SdkType)

	identifier := ClassifyAttribute(obj, Attribute{Name: "for_application", Type: AttrIdentifier, IdentifierType: "AccountObjectIdentifier"})
	require.Equal(t, "sdk.AccountObjectIdentifier", identifier.SdkType)
	require.Equal(t, "identifier", identifier.SdkKind)

	number := ClassifyAttribute(obj, Attribute{Name: "min_nodes", Type: AttrInt, Required: true})
	require.Equal(t, AttrInt, number.Storage)
	require.Empty(t, number.SdkType)
}

func Test_ClassifyParameter_primitiveStaysOnStorageType(t *testing.T) {
	boolParameter := ClassifyParameter("enable_console_output", "bool")
	require.Equal(t, AttrBool, boolParameter.Storage)
	require.Empty(t, boolParameter.SdkType)

	intParameter := ClassifyParameter("data_retention_time_in_days", "int")
	require.Equal(t, AttrInt, intParameter.Storage)
	require.Empty(t, intParameter.SdkType)

	enumParameter := ClassifyParameter("log_level", "LogLevel")
	require.Equal(t, AttrString, enumParameter.Storage)
	require.Equal(t, "sdk.LogLevel", enumParameter.SdkType)
	require.Equal(t, "enum", enumParameter.SdkKind)

	identifierParameter := ClassifyParameter("external_volume", "AccountObjectIdentifier")
	require.Equal(t, "sdk.AccountObjectIdentifier", identifierParameter.SdkType)
	require.Equal(t, "identifier", identifierParameter.SdkKind)
}

func classifiedNames(fields []ClassifiedField) []string {
	return collections.Map(fields, func(field ClassifiedField) string { return field.Name })
}

func Test_ClassifyAttribute_objectList(t *testing.T) {
	list := Attribute{Name: "location", Type: AttrList, Required: true, Elem: ObjectElem([]Attribute{{Name: "url", Type: AttrString}})}
	require.Equal(t, ClassifiedField{Name: "location", Storage: AttrList, Required: true}, ClassifyAttribute(Object{Name: "Widget"}, list))
}

func Test_ClassifyAttribute_singleBlock(t *testing.T) {
	got := ClassifyAttribute(Object{Name: "Widget"}, Attribute{Name: "block", Type: AttrObject, Required: true, Children: []Attribute{{Name: "leaf", Type: AttrString}}})
	require.Equal(t, ClassifiedField{Name: "block", Storage: AttrList, Required: true}, got)
}
