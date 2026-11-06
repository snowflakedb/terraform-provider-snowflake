package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func computedObject(computed []ssot.Attribute, show ...*ssot.Attribute) ssot.Object {
	return ssot.Object{
		Name:               "Pool",
		NamePlural:         "Pools",
		IdType:             "SchemaObjectIdentifier",
		Attributes:         []ssot.Attribute{{Name: "comment", Type: ssot.AttrString}},
		ComputedAttributes: computed,
		Show:               ssot.OutputSpec{Columns: show},
	}
}

func renderAttributesSchema(t *testing.T, model ResourceModel) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, SchemaTemplate.ExecuteTemplate(&b, "attributesSchemaDefinition", model))
	return b.String()
}

func Test_showComputedAttributes_matchesByName(t *testing.T) {
	owner := ssot.Attribute{Name: "owner", Type: ssot.AttrString, Description: "Owner."}
	showOwner := &ssot.Attribute{Name: "owner", Type: ssot.AttrString, OptionalInShow: true}
	describeOnly := ssot.Attribute{Name: "catalog_source", Type: ssot.AttrString, Description: "Source."}

	t.Run("a SHOW column with the name is selected, and it supplies the column", func(t *testing.T) {
		attrs, columns := showComputedAttributes(computedObject([]ssot.Attribute{owner, describeOnly}, showOwner))
		require.Equal(t, []ssot.Attribute{owner}, attrs)
		require.Equal(t, []*ssot.Attribute{showOwner}, columns)
	})
	t.Run("not in SHOW stays hand-written", func(t *testing.T) {
		attrs, columns := showComputedAttributes(computedObject([]ssot.Attribute{describeOnly}, showOwner))
		require.Empty(t, attrs)
		require.Empty(t, columns)
	})
}

func Test_computedAttributes_schemaAndSetter(t *testing.T) {
	owner := ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true, Description: "Role that owns the pool."}
	primary := ssot.Attribute{Name: "is_primary", Type: ssot.AttrBool, NativeBool: true, Readonly: true, Description: "Whether the pool is primary."}
	hooks := []ssot.ResourceHookOption{ssot.WithShowByIdSafelyInSdk(ssot.Generated())}
	obj := computedObject([]ssot.Attribute{owner, primary}, &owner, &primary)
	obj.ResourceHooks = hooks
	model := ModelFromInputObject(ResourceObject{Object: obj}, nil, nil)

	t.Run("schema entries are computed only", func(t *testing.T) {
		got := renderAttributesSchema(t, model)
		require.Contains(t, got, "\"owner\": {\n\t\tType: tfschema.TypeString,\n\t\tComputed: true,\n\t\tDescription: \"Role that owns the pool.\",")
		require.Contains(t, got, "\"is_primary\": {\n\t\tType: tfschema.TypeBool,\n\t\tComputed: true,")
		require.NotContains(t, got, "Optional: true,\n\t\tComputed")
	})
	t.Run("computed attributes do not need an update", func(t *testing.T) {
		require.False(t, model.HasUpdateCalls())
		require.Len(t, model.Attributes, 1)
	})
	t.Run("read sets them from SHOW", func(t *testing.T) {
		require.Contains(t, renderTemplate(t, "setComputedFields", model), "poolSetComputedFields(d, pool),")
		fn := renderTemplate(t, "generatedSetComputedFields", model)
		require.Contains(t, fn, "func poolSetComputedFields(d *tfschema.ResourceData, obj *sdk.Pool) error {")
		require.Contains(t, fn, "d.Set(\"owner\", obj.Owner),")
		require.Contains(t, fn, "d.Set(\"is_primary\", obj.IsPrimary),")
	})
	t.Run("lifecycle emits the function and the call once", func(t *testing.T) {
		got := renderLifecycle(t, model)
		require.Equal(t, 1, strings.Count(got, "func poolSetComputedFields("))
		require.Equal(t, 1, strings.Count(got, "poolSetComputedFields(d, pool),"))
	})
}

func Test_computedAttributes_noneMeansNoOutput(t *testing.T) {
	model := ModelFromInputObject(ResourceObject{Object: computedObject(nil)}, nil, nil)
	require.Empty(t, model.ComputedAttributes)
	require.Empty(t, renderTemplate(t, "setComputedFields", model))
	require.NotContains(t, renderLifecycle(t, model), "SetComputedFields")
}

func Test_computedAttributes_panics(t *testing.T) {
	hooks := ssot.WithShowByIdSafelyInSdk(ssot.Generated())
	build := func(attr ssot.Attribute, opts ...ssot.ResourceHookOption) {
		obj := computedObject([]ssot.Attribute{attr}, &attr)
		obj.ResourceHooks = opts
		ModelFromInputObject(ResourceObject{Object: obj}, nil, nil)
	}
	t.Run("without a description", func(t *testing.T) {
		require.PanicsWithValue(t, `object "Pool": computed attribute "owner" needs a Description`, func() {
			build(ssot.Attribute{Name: "owner", Type: ssot.AttrString}, hooks)
		})
	})
	t.Run("a list", func(t *testing.T) {
		require.Panics(t, func() {
			build(ssot.Attribute{Name: "tables", Type: ssot.AttrList, Elem: &ssot.Attribute{Type: ssot.AttrString}, Description: "Tables."}, hooks)
		})
	})
	t.Run("without SHOW lookup", func(t *testing.T) {
		require.Panics(t, func() {
			build(ssot.Attribute{Name: "owner", Type: ssot.AttrString, Description: "Owner."}, ssot.WithShowByIdSafelyInSdk(ssot.Omitted()))
		})
	})
}
