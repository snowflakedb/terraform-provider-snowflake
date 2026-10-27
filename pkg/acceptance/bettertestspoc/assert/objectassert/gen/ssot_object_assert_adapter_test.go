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
			ssot.WithObjectAssert(
				ssot.ObjectAssertShow(),
				ssot.ObjectAssertDescribe(),
			),
		},
		Show: ssot.OutputSpec{
			Columns: []*ssot.Attribute{
				&comment,
			},
		},
		Describe: ssot.OutputSpec{
			Columns: []*ssot.Attribute{
				&errorCode,
			},
		},
	}

	t.Run("show", func(t *testing.T) {
		details := objectAssertionDetailsFromSsot(obj, false)
		require.Equal(t, genhelpers.SdkObjectDetails{
			IdType: "sdk.AccountObjectIdentifier",
			StructDetails: genhelpers.StructDetails{
				Name: "sdk.Widget",
				Fields: []genhelpers.Field{
					ssot.ColumnField(obj, comment),
				},
			},
		}, details)

		show := SnowflakeObjectAssertionsFromSsot(obj, false)
		require.Equal(t, "Widget", show.Name)
		require.Equal(t, "sdk.Widget", show.SdkType)
		require.Equal(t, details.IdType, show.IdType)
		require.False(t, show.IsDataSourceOutput)
		require.False(t, show.NoIdentifiableObject)
		require.Equal(t, []string{"Comment"}, fieldAssertionNames(show))
	})

	t.Run("describe", func(t *testing.T) {
		details := objectAssertionDetailsFromSsot(obj, true)
		require.Equal(t, genhelpers.SdkObjectDetails{
			IdType:             "sdk.AccountObjectIdentifier",
			IsDataSourceOutput: true,
			StructDetails: genhelpers.StructDetails{
				Name: "sdk.WidgetDetails",
				Fields: []genhelpers.Field{
					ssot.ColumnField(obj, errorCode),
				},
			},
		}, details)

		desc := SnowflakeObjectAssertionsFromSsot(obj, true)
		require.Equal(t, "WidgetDetails", desc.Name)
		require.Equal(t, "sdk.WidgetDetails", desc.SdkType)
		require.True(t, desc.IsDataSourceOutput)
		require.False(t, desc.NoIdentifiableObject)
		require.Equal(t, []string{"ErrorCode"}, fieldAssertionNames(desc))
	})
}

func Test_objectAssertionDetailsFromSsot_passesSkipFields(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	tables := ssot.Attribute{Name: "tables", Type: ssot.AttrString}
	password := ssot.Attribute{Name: "password", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:   "Widget",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithObjectAssert(
				ssot.ObjectAssertShow(
					ssot.ObjectAssertShowSkipFields(&tables),
				),
				ssot.ObjectAssertDescribe(
					ssot.ObjectAssertDescribeSkipFields(&password),
				),
			),
		},
		Show: ssot.OutputSpec{
			Columns: []*ssot.Attribute{
				&comment,
				&tables,
			},
		},
		Describe: ssot.OutputSpec{
			Columns: []*ssot.Attribute{
				&password,
			},
		},
	}

	show := objectAssertionDetailsFromSsot(obj, false)
	require.Equal(t, []string{"Tables"}, show.SkipFields)

	withoutSkips := ssot.Object{
		Name:   "Widget",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithObjectAssert(
				ssot.ObjectAssertShow(),
			),
		},
		Show: ssot.OutputSpec{
			Columns: []*ssot.Attribute{
				&comment,
			},
		},
	}
	unskipped := objectAssertionDetailsFromSsot(withoutSkips, false)
	require.Nil(t, unskipped.SkipFields)

	desc := objectAssertionDetailsFromSsot(obj, true)
	require.Equal(t, []string{"Password"}, desc.SkipFields)

	model := SnowflakeObjectAssertionsFromSsot(obj, false)
	require.False(t, fieldByName(t, model, "Comment").SkipGeneration)
	require.True(t, fieldByName(t, model, "Tables").SkipGeneration)
}

func Test_objectAssertionDetailsFromSsot_passesNoIdentifiableObject(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	obj := ssot.Object{
		Name:   "Widget",
		IdType: "DatabaseObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithObjectAssert(
				ssot.ObjectAssertDescribe(
					ssot.WithNoIdentifiableObject(),
				),
			),
		},
		Describe: ssot.OutputSpec{
			Columns: []*ssot.Attribute{
				&name,
			},
		},
	}

	details := objectAssertionDetailsFromSsot(obj, true)
	require.Equal(t, "sdk.WidgetDetails", details.ObjectName())
	require.True(t, details.IsDataSourceOutput)
	require.True(t, details.NoIdentifiableObject)
	require.Equal(t, []genhelpers.Field{
		ssot.ColumnField(obj, name),
	}, details.Fields)
}

func Test_objectAssertionDetailsFromSsot_showProjection(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	base := ssot.Object{Name: "Widget"}
	obj := ssot.Object{
		Name:      "WidgetAdaptive",
		IdType:    "AccountObjectIdentifier",
		VariantOf: &base,
		Generators: []ssot.GeneratorOption{
			ssot.WithObjectAssert(
				ssot.ObjectAssertShow(
					ssot.ObjectAssertShowFromBase("ShowByIDAdaptive"),
				),
			),
		},
		Show: ssot.OutputSpec{
			Columns: []*ssot.Attribute{
				&comment,
			},
		},
	}

	details := objectAssertionDetailsFromSsot(obj, false)
	require.Equal(t, "Widget", details.ObjectTypeName)
	require.Equal(t, &genhelpers.ShowProjectionDef{
		ClientName: "Widget",
		Method:     "ShowByIDAdaptive",
	}, details.ShowProjection)
	model := SnowflakeObjectAssertionsFromSsot(obj, false)
	require.Equal(t, "Widget", model.ObjectTypeName)
	require.Equal(t, details.ShowProjection, model.ShowProjection)
}

func Test_SnowflakeObjectAssertionsFromSsot_panics(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	enabled := []ssot.GeneratorOption{
		ssot.WithObjectAssert(
			ssot.ObjectAssertShow(),
		),
	}

	t.Run("show option off", func(t *testing.T) {
		obj := ssot.Object{
			Name:   "Widget",
			IdType: "AccountObjectIdentifier",
			Show: ssot.OutputSpec{
				Columns: []*ssot.Attribute{
					&comment,
				},
			},
		}
		require.Panics(t, func() {
			SnowflakeObjectAssertionsFromSsot(obj, false)
		})
	})

	t.Run("describe option off", func(t *testing.T) {
		obj := ssot.Object{
			Name:       "Widget",
			IdType:     "AccountObjectIdentifier",
			Generators: enabled,
			Show: ssot.OutputSpec{
				Columns: []*ssot.Attribute{
					&comment,
				},
			},
			Describe: ssot.OutputSpec{
				Columns: []*ssot.Attribute{
					&comment,
				},
			},
		}
		require.Panics(t, func() {
			SnowflakeObjectAssertionsFromSsot(obj, true)
		})
	})

	t.Run("empty columns", func(t *testing.T) {
		obj := ssot.Object{
			Name:       "Widget",
			IdType:     "AccountObjectIdentifier",
			Generators: enabled,
		}
		require.Panics(t, func() {
			SnowflakeObjectAssertionsFromSsot(obj, false)
		})
	})

	t.Run("nil column", func(t *testing.T) {
		obj := ssot.Object{
			Name:       "Widget",
			IdType:     "AccountObjectIdentifier",
			Generators: enabled,
			Show: ssot.OutputSpec{
				Columns: []*ssot.Attribute{
					nil,
				},
			},
		}
		require.Panics(t, func() {
			SnowflakeObjectAssertionsFromSsot(obj, false)
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

func Test_objectAssertionDetailsFromSsot_passesNestedFields(t *testing.T) {
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString}
	locations := ssot.Attribute{Name: "storage_locations", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrString}}
	obj := ssot.Object{
		Name:   "Widget",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithObjectAssert(
				ssot.ObjectAssertDescribe(
					ssot.ObjectAssertDescribeNestedField(&locations),
				),
			),
		},
		Describe: ssot.OutputSpec{
			Columns: []*ssot.Attribute{
				&name,
				&locations,
			},
		},
	}

	desc := objectAssertionDetailsFromSsot(obj, true)
	require.Equal(t, []string{"StorageLocations"}, desc.NestedAssertFields)
	require.Equal(t, "[]sdk.WidgetStorageLocationDetails", desc.Fields[1].ConcreteType)
	require.Equal(t, "slice", desc.Fields[1].UnderlyingType)
	require.Equal(t, "string", desc.Fields[0].ConcreteType)

	model := SnowflakeObjectAssertionsFromSsot(obj, true)
	require.False(t, model.Fields[0].GenerateNestedAdapter)
	require.True(t, model.Fields[1].GenerateNestedAdapter)
}

func Test_objectAssertionDetailsFromSsot_nestedFieldGoTypeOverride(t *testing.T) {
	locations := ssot.Attribute{Name: "storage_locations", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrString}}
	obj := ssot.Object{
		Name:   "Widget",
		IdType: "AccountObjectIdentifier",
		Generators: []ssot.GeneratorOption{
			ssot.WithObjectAssert(
				ssot.ObjectAssertDescribe(
					ssot.ObjectAssertDescribeNestedField(&locations, ssot.ObjectAssertNestedGoType("[]sdk.Custom")),
				),
			),
		},
		Describe: ssot.OutputSpec{
			Columns: []*ssot.Attribute{
				&locations,
			},
		},
	}

	desc := objectAssertionDetailsFromSsot(obj, true)
	require.Equal(t, "[]sdk.Custom", desc.Fields[0].ConcreteType)
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
