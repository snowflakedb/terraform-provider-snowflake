package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_SnowflakeObjectAssertionsFromSsot_matchesHandwrittenFields(t *testing.T) {
	family := ssot.NewEnum("PoolFamily", "PoolFamilies", "SMALL")
	stateEnum := ssot.NewEnum("PoolState", "PoolStates", "IDLE")
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	autoResume := ssot.Attribute{Name: "auto_resume", Type: ssot.AttrBool}
	isDefault := ssot.Attribute{Name: "is_default", Type: ssot.AttrBool, StoredAsText: true}
	options := ssot.Attribute{Name: "options", Type: ssot.AttrString, OptionalInShow: true}
	droppedOn := ssot.Attribute{Name: "dropped_on", Type: ssot.AttrTime, OptionalInDb: true}
	state := ssot.Attribute{Name: "state", Type: ssot.AttrEnum, Enum: stateEnum}
	forApplication := ssot.Attribute{Name: "for_application", Type: ssot.AttrIdentifier, OptionalInShow: true, ShowName: "application"}
	backupFamilies := ssot.Attribute{
		Name: "backup_instance_families",
		Type: ssot.AttrList,
		Elem: &ssot.Attribute{Type: ssot.AttrEnum, Enum: family},
	}
	minNodes := ssot.Attribute{Name: "min_nodes", Type: ssot.AttrInt}
	errorCode := ssot.Attribute{Name: "error_code", Type: ssot.AttrString}
	showColumns := []*ssot.Attribute{
		&name, &autoResume, &isDefault, &options, &droppedOn, &state, &forApplication, &backupFamilies, &minNodes,
	}
	describeColumns := []*ssot.Attribute{&name, &errorCode}
	obj := ssot.Object{
		Name:   "Pool",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithObjectAssert(ssot.ObjectAssertShow(), ssot.ObjectAssertDescribe()),
		},
		Show:     ssot.OutputSpec{Columns: showColumns},
		Describe: ssot.OutputSpec{Columns: describeColumns},
	}

	showDetails := objectAssertionDetailsFromSsot(obj, false)
	require.Equal(t, "sdk.Pool", showDetails.ObjectName())
	require.False(t, showDetails.IsDataSourceOutput)
	require.Equal(t, "sdk.AccountObjectIdentifier", showDetails.IdType)
	require.Empty(t, showDetails.ObjectTypeName)

	showFields := []genhelpers.Field{
		{Name: "Name", ConcreteType: "string", UnderlyingType: "string"},
		{Name: "AutoResume", ConcreteType: "bool", UnderlyingType: "bool"},
		{Name: "IsDefault", ConcreteType: "bool", UnderlyingType: "bool"},
		{Name: "Options", ConcreteType: "*string", UnderlyingType: "*string"},
		{Name: "DroppedOn", ConcreteType: "time.Time", UnderlyingType: "struct"},
		{Name: "State", ConcreteType: "sdk.PoolState", UnderlyingType: "string"},
		{Name: "Application", ConcreteType: "*sdk.AccountObjectIdentifier", UnderlyingType: "*struct"},
		{Name: "BackupInstanceFamilies", ConcreteType: "[]string", UnderlyingType: "slice"},
		{Name: "MinNodes", ConcreteType: "int", UnderlyingType: "int"},
	}
	show := SnowflakeObjectAssertionsFromSsot(obj, false)
	requireSnowflakeObjectAssertionsModelEqual(t, ModelFromSdkObjectDetails(genhelpers.SdkObjectDetails{
		IdType: "sdk.AccountObjectIdentifier",
		StructDetails: genhelpers.StructDetails{
			Name:   "sdk.Pool",
			Fields: showFields,
		},
	}, nil), show)
	application := fieldByName(t, show, "Application")
	require.Contains(t, application.Mapper("o.Application"), ".FullyQualifiedName()")
	require.Contains(t, application.ExpectedValueMapper("expected"), ".FullyQualifiedName()")

	descDetails := objectAssertionDetailsFromSsot(obj, true)
	require.Equal(t, "sdk.PoolDetails", descDetails.ObjectName())
	require.True(t, descDetails.IsDataSourceOutput)
	require.Empty(t, descDetails.ObjectTypeName)
	desc := SnowflakeObjectAssertionsFromSsot(obj, true)
	requireSnowflakeObjectAssertionsModelEqual(t, ModelFromSdkObjectDetails(genhelpers.SdkObjectDetails{
		IdType:             "sdk.AccountObjectIdentifier",
		IsDataSourceOutput: true,
		StructDetails: genhelpers.StructDetails{
			Name: "sdk.PoolDetails",
			Fields: []genhelpers.Field{
				{Name: "Name", ConcreteType: "string", UnderlyingType: "string"},
				{Name: "ErrorCode", ConcreteType: "string", UnderlyingType: "string"},
			},
		},
	}, nil), desc)
}

func Test_objectAssertionDetailsFromSsot_passesNoIdentifiableObject(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:   "Schema",
		IdType: "DatabaseObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithObjectAssert(ssot.ObjectAssertDescribe(ssot.WithNoIdentifiableObject())),
		},
		Describe: ssot.OutputSpec{Columns: []*ssot.Attribute{&name}},
	}

	details := objectAssertionDetailsFromSsot(obj, true)
	require.Equal(t, "sdk.SchemaDetails", details.ObjectName())
	require.True(t, details.IsDataSourceOutput)
	require.True(t, details.NoIdentifiableObject)
}

func Test_SnowflakeObjectAssertionsFromSsot_panics(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	enabled := []ssot.GeneratorOption{ssot.WithObjectAssert(ssot.ObjectAssertShow())}

	t.Run("show option off", func(t *testing.T) {
		require.Panics(t, func() {
			SnowflakeObjectAssertionsFromSsot(ssot.Object{
				Name:   "Pool",
				IdType: "AccountObjectIdentifier",
				Show:   ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}},
			}, false)
		})
	})
	t.Run("describe option off", func(t *testing.T) {
		require.Panics(t, func() {
			SnowflakeObjectAssertionsFromSsot(ssot.Object{
				Name:       "Pool",
				IdType:     "AccountObjectIdentifier",
				Generators: enabled,
				Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}},
				Describe:   ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}},
			}, true)
		})
	})
	t.Run("empty columns", func(t *testing.T) {
		require.Panics(t, func() {
			SnowflakeObjectAssertionsFromSsot(ssot.Object{
				Name:       "Pool",
				IdType:     "AccountObjectIdentifier",
				Generators: enabled,
			}, false)
		})
	})
	t.Run("nil column", func(t *testing.T) {
		require.Panics(t, func() {
			SnowflakeObjectAssertionsFromSsot(ssot.Object{
				Name:       "Pool",
				IdType:     "AccountObjectIdentifier",
				Generators: enabled,
				Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{nil}},
			}, false)
		})
	})
}

func requireSnowflakeObjectAssertionsModelEqual(t *testing.T, expected, actual SnowflakeObjectAssertionsModel) {
	t.Helper()
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.SdkType, actual.SdkType)
	require.Equal(t, expected.IdType, actual.IdType)
	require.Equal(t, expected.IsDataSourceOutput, actual.IsDataSourceOutput)
	require.Equal(t, expected.IsSubStruct, actual.IsSubStruct)
	require.Equal(t, expected.ObjectTypeName, actual.ObjectTypeName)
	require.Equal(t, expected.NoShowById, actual.NoShowById)
	require.Equal(t, expected.NoIdentifiableObject, actual.NoIdentifiableObject)
	require.Equal(t, expected.ShowByParentId, actual.ShowByParentId)
	require.Equal(t, expected.DescribeOverride, actual.DescribeOverride)
	require.Equal(t, expected.FromObjectIDExpr, actual.FromObjectIDExpr)
	require.Equal(t, expected.NestedAssertFields, actual.NestedAssertFields)
	require.Equal(t, expected.PreambleModel, actual.PreambleModel)
	require.Len(t, actual.Fields, len(expected.Fields))
	for i := range expected.Fields {
		requireFieldAssertionEqual(t, expected.Fields[i], actual.Fields[i])
	}
}

func requireFieldAssertionEqual(t *testing.T, expected, actual SnowflakeObjectFieldAssertion) {
	t.Helper()
	require.Equal(t, expected.Name, actual.Name)
	require.Equal(t, expected.ConcreteType, actual.ConcreteType)
	require.Equal(t, expected.IsOriginalTypePointer, actual.IsOriginalTypePointer)
	require.Equal(t, expected.IsOriginalTypeSlice, actual.IsOriginalTypeSlice)
	require.Equal(t, expected.ComparisonMethod, actual.ComparisonMethod)
	require.Equal(t, expected.GenerateNestedAdapter, actual.GenerateNestedAdapter)
	require.Equal(t, expected.SkipGeneration, actual.SkipGeneration)
	const probe = "field"
	require.Equal(t, expected.Mapper(probe), actual.Mapper(probe))
	require.Equal(t, expected.ExpectedValueMapper(probe), actual.ExpectedValueMapper(probe))
	require.Equal(t, expected.ErrorMapper(probe), actual.ErrorMapper(probe))
	require.Equal(t, expected.ErrorExpectedValueMapper(probe), actual.ErrorExpectedValueMapper(probe))
}

func fieldByName(t *testing.T, model SnowflakeObjectAssertionsModel, name string) SnowflakeObjectFieldAssertion {
	t.Helper()
	for _, field := range model.Fields {
		if field.Name == name {
			return field
		}
	}
	t.Fatalf("field %q not found", name)
	return SnowflakeObjectFieldAssertion{}
}
