package gen

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_SnowflakeObjectAssertionsFromSsot_usesShowAndDescribeColumns(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	errorCode := ssot.Attribute{Name: "error_code", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:   "Widget",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithObjectAssert(ssot.ObjectAssertShow(), ssot.ObjectAssertDescribe()),
		},
		Show:     ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}},
		Describe: ssot.OutputSpec{Columns: []*ssot.Attribute{&errorCode}},
	}

	showDetails := objectAssertionDetailsFromSsot(obj, false)
	require.Equal(t, genhelpers.SdkObjectDetails{
		IdType: "sdk.AccountObjectIdentifier",
		StructDetails: genhelpers.StructDetails{
			Name:   "sdk.Widget",
			Fields: []genhelpers.Field{ssot.ColumnField(obj, comment)},
		},
	}, showDetails)
	show := SnowflakeObjectAssertionsFromSsot(obj, false)
	require.Equal(t, "Widget", show.Name)
	require.Equal(t, "sdk.Widget", show.SdkType)
	require.Equal(t, showDetails.IdType, show.IdType)
	require.False(t, show.IsDataSourceOutput)
	require.False(t, show.NoIdentifiableObject)
	require.Equal(t, []string{"Comment"}, fieldAssertionNames(show))

	descDetails := objectAssertionDetailsFromSsot(obj, true)
	require.Equal(t, genhelpers.SdkObjectDetails{
		IdType:             "sdk.AccountObjectIdentifier",
		IsDataSourceOutput: true,
		StructDetails: genhelpers.StructDetails{
			Name:   "sdk.WidgetDetails",
			Fields: []genhelpers.Field{ssot.ColumnField(obj, errorCode)},
		},
	}, descDetails)
	desc := SnowflakeObjectAssertionsFromSsot(obj, true)
	require.Equal(t, "WidgetDetails", desc.Name)
	require.Equal(t, "sdk.WidgetDetails", desc.SdkType)
	require.True(t, desc.IsDataSourceOutput)
	require.False(t, desc.NoIdentifiableObject)
	require.Equal(t, []string{"ErrorCode"}, fieldAssertionNames(desc))
}

func Test_objectAssertionDetailsFromSsot_passesNoIdentifiableObject(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:   "Widget",
		IdType: "DatabaseObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithObjectAssert(ssot.ObjectAssertDescribe(ssot.WithNoIdentifiableObject())),
		},
		Describe: ssot.OutputSpec{Columns: []*ssot.Attribute{&name}},
	}

	details := objectAssertionDetailsFromSsot(obj, true)
	require.Equal(t, "sdk.WidgetDetails", details.ObjectName())
	require.True(t, details.IsDataSourceOutput)
	require.True(t, details.NoIdentifiableObject)
	require.Equal(t, []genhelpers.Field{ssot.ColumnField(obj, name)}, details.Fields)
}

func Test_SnowflakeObjectAssertionsFromSsot_panics(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	enabled := []ssot.GeneratorOption{ssot.WithObjectAssert(ssot.ObjectAssertShow())}

	t.Run("show option off", func(t *testing.T) {
		require.Panics(t, func() {
			SnowflakeObjectAssertionsFromSsot(ssot.Object{
				Name:   "Widget",
				IdType: "AccountObjectIdentifier",
				Show:   ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}},
			}, false)
		})
	})
	t.Run("describe option off", func(t *testing.T) {
		require.Panics(t, func() {
			SnowflakeObjectAssertionsFromSsot(ssot.Object{
				Name:       "Widget",
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
				Name:       "Widget",
				IdType:     "AccountObjectIdentifier",
				Generators: enabled,
			}, false)
		})
	})
	t.Run("nil column", func(t *testing.T) {
		require.Panics(t, func() {
			SnowflakeObjectAssertionsFromSsot(ssot.Object{
				Name:       "Widget",
				IdType:     "AccountObjectIdentifier",
				Generators: enabled,
				Show:       ssot.OutputSpec{Columns: []*ssot.Attribute{nil}},
			}, false)
		})
	})
}

func fieldAssertionNames(model SnowflakeObjectAssertionsModel) []string {
	names := make([]string, len(model.Fields))
	for i, field := range model.Fields {
		names[i] = field.Name
	}
	return names
}
