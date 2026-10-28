package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_CustomDiffEntries(t *testing.T) {
	comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
	methods := ssot.Attribute{Name: "authentication_methods", Type: ssot.AttrSet}
	nested := ssot.Attribute{
		Name: "mfa_policy",
		Type: ssot.AttrList,
		Elem: ssot.ObjectElem([]ssot.Attribute{
			{Name: "allowed_methods", Type: ssot.AttrSet},
		}),
	}
	createdOn := ssot.Attribute{Name: "created_on", Type: ssot.AttrTime, Readonly: true}
	name := ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true}
	owner := ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true}
	databaseName := ssot.Attribute{Name: "database_name", Type: ssot.AttrString, Readonly: true}

	model := func(opts ...func(*ssot.Object)) ResourceModel {
		obj := ssot.Object{
			Name:                      "Policy",
			IdType:                    "SchemaObjectIdentifier",
			IdentityHierarchyForceNew: true,
			Attributes:                []ssot.Attribute{methods, nested, comment},
			Show: ssot.OutputSpec{Columns: []*ssot.Attribute{
				&createdOn, &name, &comment, &databaseName, &nested,
			}},
			Describe: ssot.OutputSpec{Columns: []*ssot.Attribute{
				&name, &owner, &methods, &nested, &comment,
			}},
			ResourceHooks: []ssot.ResourceHookOption{
				ssot.WithCustomizeDiff(ssot.Generated()),
			},
		}
		for _, opt := range opts {
			opt(&obj)
		}
		return ModelFromInputObject(ResourceObject{Object: obj}, nil, nil)
	}

	t.Run("show describe and non-force-new name", func(t *testing.T) {
		require.Equal(t, []CustomDiffEntry{
			{Key: "ShowOutputAttributeName", Attrs: []string{"name", "comment"}},
			{Key: "DescribeOutputAttributeName", Attrs: []string{"name", "comment"}},
			{Key: "FullyQualifiedNameAttributeName", Attrs: []string{"name"}},
		}, model().CustomDiffEntries())
	})

	t.Run("describe-only read-only column follows the attribute of its show name", func(t *testing.T) {
		enrollment := ssot.Attribute{Name: "enrollment", Type: ssot.AttrString}
		describeEnrollment := ssot.Attribute{Name: "describe_enrollment", ShowName: "enrollment", Type: ssot.AttrString, Readonly: true}
		unrelated := ssot.Attribute{Name: "describe_other", ShowName: "other", Type: ssot.AttrString, Readonly: true}
		got := model(func(obj *ssot.Object) {
			obj.Attributes = append(obj.Attributes, enrollment)
			obj.Describe.Columns = []*ssot.Attribute{&name, &describeEnrollment, &unrelated, &comment}
		}).CustomDiffEntries()
		require.Contains(t, got, CustomDiffEntry{Key: "DescribeOutputAttributeName", Attrs: []string{"name", "enrollment", "comment"}})
	})

	t.Run("fully qualified name includes identifiers that are not force new", func(t *testing.T) {
		got := model(func(obj *ssot.Object) {
			obj.IdentityHierarchyForceNew = false
			obj.IdentityNameForceNew = true
		}).CustomDiffEntries()
		require.Equal(t, []string{"database", "schema"}, got[2].Attrs)
	})

	t.Run("fully qualified name omitted when every identifier is force new", func(t *testing.T) {
		got := model(func(obj *ssot.Object) {
			obj.IdType = "AccountObjectIdentifier"
			obj.IdentityNameForceNew = true
		}).CustomDiffEntries()
		require.Len(t, got, 2)
		require.Equal(t, "DescribeOutputAttributeName", got[1].Key)
	})

	t.Run("set and list columns are excluded", func(t *testing.T) {
		got := model().CustomDiffEntries()
		require.Equal(t, []string{"name", "comment"}, got[1].Attrs)
		require.NotContains(t, got[1].Attrs, "authentication_methods")
		require.NotContains(t, got[1].Attrs, "mfa_policy")
	})

	t.Run("parameters computed if comes before fully qualified name", func(t *testing.T) {
		got := model(func(obj *ssot.Object) {
			obj.ResourceHooks = append(obj.ResourceHooks,
				ssot.WithParametersAttributesSchema(ssot.Generated()),
				ssot.WithParametersFieldNames(ssot.ExtCall()),
				ssot.WithParametersCustomDiff(ssot.Generated()),
			)
		}).CustomDiffCalls()
		require.Equal(t, []string{
			`ComputedIfAnyAttributeChanged(policySchema, ShowOutputAttributeName, "name", "comment")`,
			`ComputedIfAnyAttributeChanged(policySchema, DescribeOutputAttributeName, "name", "comment")`,
			"ComputedIfAnyAttributeChanged(policyParametersAttributesSchema, ParametersAttributeName, policyParametersFieldNamesExt...)",
			`ComputedIfAnyAttributeChanged(policySchema, FullyQualifiedNameAttributeName, "name")`,
			"policyParametersCustomDiff",
		}, got)
	})

	t.Run("additional customize diffs follow the generated calls", func(t *testing.T) {
		got := model(func(obj *ssot.Object) {
			obj.AdditionalCustomizeDiffs = []ssot.CustomizeDiffFunc{
				"policyForceNewIfSizeRemoved",
				"policyForceNewIfWasInteractive",
			}
		}).CustomDiffCalls()
		require.Equal(t, []string{
			`ComputedIfAnyAttributeChanged(policySchema, ShowOutputAttributeName, "name", "comment")`,
			`ComputedIfAnyAttributeChanged(policySchema, DescribeOutputAttributeName, "name", "comment")`,
			`ComputedIfAnyAttributeChanged(policySchema, FullyQualifiedNameAttributeName, "name")`,
			"policyForceNewIfSizeRemoved",
			"policyForceNewIfWasInteractive",
		}, got)
	})

	t.Run("omitted output schemas are left out", func(t *testing.T) {
		got := model(func(obj *ssot.Object) {
			obj.ResourceHooks = append(obj.ResourceHooks,
				ssot.WithShowOutputSchema(ssot.Omitted()),
				ssot.WithDescribeOutputSchema(ssot.Omitted()),
				ssot.WithFqnSchema(ssot.Omitted()),
			)
		}).CustomDiffEntries()
		require.Empty(t, got)
	})
}

func Test_lifecycle_customizeDiff(t *testing.T) {
	t.Run("ext call", func(t *testing.T) {
		got := renderLifecycle(t, lifecycleModel())
		require.Contains(t, got, "CustomizeDiff: poolCustomizeDiffExt,")
		require.NotContains(t, got, "var poolCustomizeDiff")
	})

	t.Run("generated", func(t *testing.T) {
		comment := ssot.Attribute{Name: "comment", Type: ssot.AttrString}
		model := lifecycleModel(func(obj *ssot.Object) {
			obj.IdentityHierarchyForceNew = true
			obj.Attributes = []ssot.Attribute{comment}
			obj.Show = ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}}
			obj.Describe = ssot.OutputSpec{Columns: []*ssot.Attribute{&comment}}
			obj.ResourceHooks = []ssot.ResourceHookOption{
				ssot.WithCustomizeDiff(ssot.Generated()),
				ssot.WithDescribeOutputSchema(ssot.Omitted()),
			}
		})
		got := renderLifecycle(t, model)
		require.Contains(t, got, "CustomizeDiff: poolCustomizeDiff,")
		require.NotContains(t, got, "poolCustomizeDiffExt")
		require.Contains(t, got, "var poolCustomizeDiff = TrackingCustomDiffWrapper(resources.Pool, customdiff.All(")
		require.NotContains(t, got, "SNOW-1648997")
		require.Contains(t, got, `ComputedIfAnyAttributeChanged(poolSchema, ShowOutputAttributeName, "comment"),`)
		require.Contains(t, got, `ComputedIfAnyAttributeChanged(poolSchema, FullyQualifiedNameAttributeName, "name"),`)
		_, diff, ok := strings.Cut(got, "var poolCustomizeDiff")
		require.True(t, ok)
		require.NotContains(t, diff, "DescribeOutputAttributeName")
		require.Equal(t, 1, strings.Count(got, "customdiff.All("))
	})

	t.Run("generated with additional functions", func(t *testing.T) {
		model := lifecycleModel(func(obj *ssot.Object) {
			obj.AdditionalCustomizeDiffs = []ssot.CustomizeDiffFunc{
				"poolForceNewIfSizeRemoved",
			}
			obj.ResourceHooks = []ssot.ResourceHookOption{
				ssot.WithCustomizeDiff(ssot.Generated()),
				ssot.WithShowOutputSchema(ssot.Omitted()),
				ssot.WithDescribeOutputSchema(ssot.Omitted()),
				ssot.WithFqnSchema(ssot.Omitted()),
			}
		})
		got := renderLifecycle(t, model)
		require.Contains(t, got, "CustomizeDiff: poolCustomizeDiff,")
		require.Contains(t, got, "poolForceNewIfSizeRemoved,")
		require.NotContains(t, got, "append([]tfschema.CustomizeDiffFunc{")
		require.NotContains(t, got, "CustomizeDiff: poolCustomizeDiffExt,")
	})
}
