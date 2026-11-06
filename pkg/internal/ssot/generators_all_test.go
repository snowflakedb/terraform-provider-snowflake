package ssot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithAllGenerators_equalsEnablingEveryGeneratorByHand(t *testing.T) {
	byHand := ObjectGenerators(
		WithObjectParametersAssert(),
		WithResourceParametersAssert(),
		WithSchema(SchemaShow(), SchemaDescribe()),
		WithObjectAssert(ObjectAssertShow(), ObjectAssertDescribe()),
		WithShowOutputAssert(ShowOutputShow(), ShowOutputDescribe()),
		WithResource(),
		WithModel(),
		WithResourceAssert(),
	)
	require.Equal(t, byHand, ObjectGenerators(WithAllGenerators()))
}

func Test_WithAllGenerators_leavesTheSdkAlone(t *testing.T) {
	require.Equal(t, SdkGenerator{}, ObjectGenerators(WithAllGenerators()).Sdk)

	g := ObjectGenerators(WithAllGenerators(), WithSdk(SdkCreate()))
	require.True(t, g.Sdk.Create)
	require.True(t, g.Model.Enabled)
}

func Test_WithAllGenerators_overridesApplyInOrder(t *testing.T) {
	t.Run("without turns a generator off and the rest stay on", func(t *testing.T) {
		g := ObjectGenerators(WithAllGenerators(WithoutModel(), WithoutResourceAssert()))
		require.False(t, g.Model.Enabled)
		require.False(t, g.ResourceAssert.Enabled)
		require.True(t, g.Resource.Enabled)
		require.True(t, g.Schema.Show)
	})
	t.Run("a later override wins", func(t *testing.T) {
		require.True(t, ObjectGenerators(WithAllGenerators(WithoutModel(), WithModel())).Model.Enabled)
		require.False(t, ObjectGenerators(WithAllGenerators(WithModel(), WithoutModel())).Model.Enabled)
	})
}

func Test_WithAllGenerators_modificationsAddToTheDefaults(t *testing.T) {
	tables := Attribute{Name: "tables"}
	account := Attribute{Name: "account"}
	g := ObjectGenerators(WithAllGenerators(
		WithSchema(SchemaShow(SchemaShowListFields(&tables))),
		WithObjectAssert(ObjectAssertShow(ObjectAssertShowSkipFields(&tables)), ObjectAssertDescribe(ObjectAssertDescribeMethod("DescribeDetails"))),
		WithShowOutputAssert(ShowOutputDescribe(ShowOutputDescribeSkipFields(&account))),
		WithModel(ModelMultiline(&tables)),
	))

	require.Equal(t, []*Attribute{&tables}, g.Schema.ShowListFields)
	require.True(t, g.Schema.Show)
	require.True(t, g.Schema.Describe, "the other side stays on")
	require.Equal(t, []*Attribute{&tables}, g.ObjectAssert.ShowSkipFields)
	require.Equal(t, "DescribeDetails", g.ObjectAssert.DescribeMethod)
	require.Equal(t, []*Attribute{&account}, g.ShowOutputAssert.DescribeSkipFields)
	require.True(t, g.ShowOutputAssert.Show)
	require.Equal(t, []*Attribute{&tables}, g.Model.Multiline)
	require.True(t, g.Model.Enabled)
}

func Test_Without_wholeGeneratorIsZero(t *testing.T) {
	g := ObjectGenerators(WithAllGenerators(
		WithoutSchema(),
		WithoutObjectAssert(),
		WithoutShowOutputAssert(),
		WithoutModel(),
		WithoutResource(),
		WithoutResourceAssert(),
		WithoutObjectParametersAssert(),
		WithoutResourceParametersAssert(),
	))
	require.Equal(t, Generators{}, g)
}

func Test_Without_oneSideDropsItsSettingsAndKeepsTheOtherSide(t *testing.T) {
	a := Attribute{Name: "a"}
	nested := []ObjectAssertNestedOption{}

	t.Run("schema", func(t *testing.T) {
		g := ObjectGenerators(
			WithSchema(
				SchemaShow(SchemaShowFromBase(), SchemaShowManualFields(&a), SchemaShowListFields(&a)),
				SchemaDescribe(SchemaDescribeManualFields(&a), SchemaDescribeListFields(&a)),
			),
			WithoutSchemaShow(),
		)
		require.Equal(t, SchemaGenerator{
			Describe:             true,
			DescribeManualFields: []*Attribute{&a},
			DescribeListFields:   []*Attribute{&a},
		}, g.Schema)

		g = ObjectGenerators(
			WithSchema(SchemaShow(SchemaShowListFields(&a)), SchemaDescribe(SchemaDescribeListFields(&a))),
			WithoutSchemaDescribe(),
		)
		require.Equal(t, SchemaGenerator{Show: true, ShowListFields: []*Attribute{&a}}, g.Schema)
	})
	t.Run("object assertion", func(t *testing.T) {
		g := ObjectGenerators(
			WithObjectAssert(
				ObjectAssertShow(ObjectAssertShowSkipFields(&a), ObjectAssertShowFromBase("ShowByIDAdaptive"), ObjectAssertShowNestedField(&a, nested...)),
				ObjectAssertDescribe(ObjectAssertDescribeSkipFields(&a), ObjectAssertDescribeMethod("DescribeDetails"), WithNoIdentifiableObject()),
			),
			WithoutObjectAssertShow(),
		)
		require.Equal(t, ObjectAssertGenerator{
			Describe:             true,
			NoIdentifiableObject: true,
			DescribeSkipFields:   []*Attribute{&a},
			DescribeMethod:       "DescribeDetails",
		}, g.ObjectAssert)

		g = ObjectGenerators(
			WithObjectAssert(
				ObjectAssertShow(ObjectAssertShowRowFromBase(), ObjectAssertShowSkipFields(&a)),
				ObjectAssertDescribe(ObjectAssertDescribeFromBase("DescribeAwsGlue"), ObjectAssertDescribeSkipFields(&a), WithNoIdentifiableObject(), ObjectAssertDescribeNestedField(&a)),
			),
			WithoutObjectAssertDescribe(),
		)
		require.Equal(t, ObjectAssertGenerator{Show: true, ShowRowFromBase: true, ShowSkipFields: []*Attribute{&a}}, g.ObjectAssert)
	})
	t.Run("show output assertion", func(t *testing.T) {
		g := ObjectGenerators(
			WithShowOutputAssert(
				ShowOutputShow(ShowOutputShowRowFromBase(), ShowOutputShowSkipFields(&a)),
				ShowOutputDescribe(ShowOutputDescribeSkipFields(&a)),
			),
			WithoutShowOutputAssertShow(),
		)
		require.Equal(t, ShowOutputAssertGenerator{Describe: true, DescribeSkipFields: []*Attribute{&a}}, g.ShowOutputAssert)

		g = ObjectGenerators(
			WithShowOutputAssert(ShowOutputShow(ShowOutputShowSkipFields(&a)), ShowOutputDescribe(ShowOutputDescribeSkipFields(&a))),
			WithoutShowOutputAssertDescribe(),
		)
		require.Equal(t, ShowOutputAssertGenerator{Show: true, ShowSkipFields: []*Attribute{&a}}, g.ShowOutputAssert)
	})
}

func Test_Without_leavesAnObjectThatValidates(t *testing.T) {
	// Skip fields set while the side is off fail validation, so turning a side off must drop them.
	tables := Attribute{Name: "tables", Type: AttrString}
	obj := Object{
		Name:       "Widget",
		Attributes: []Attribute{tables},
		Show:       OutputSpec{Columns: []*Attribute{&tables}},
		Generators: []GeneratorOption{
			WithObjectAssert(ObjectAssertShow(ObjectAssertShowSkipFields(&tables))),
			WithoutObjectAssertShow(),
		},
	}
	require.NotPanics(t, func() { ValidateObject(obj) })
}
