package gen

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func modelOf(object ssot.Object) DataSourceModel {
	return ModelFromInputObject(DataSourceObject{Object: object}, nil, nil)
}

func renderTemplate(t *testing.T, name string, model DataSourceModel) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, DataSourceTemplate.ExecuteTemplate(&b, name, model))
	return b.String()
}

func syntheticObject(sdk []ssot.SdkOption, hooks ...ssot.DataSourceHookOption) ssot.Object {
	return ssot.Object{
		Name:            "Widget",
		NamePlural:      "Widgets",
		Show:            ssot.OutputSpec{Columns: []*ssot.Attribute{{Name: "name", Readonly: true}}},
		Describe:        ssot.OutputSpec{Columns: []*ssot.Attribute{{Name: "created_on", Readonly: true}}},
		Generators:      []ssot.GeneratorOption{ssot.WithSdk(sdk...)},
		DataSourceHooks: hooks,
	}
}

// widgetWith is a Widget with SHOW options and a DESCRIBE.
func widgetWith(show ...ssot.SdkShowOption) ssot.Object {
	return syntheticObject([]ssot.SdkOption{ssot.SdkShow(show...), ssot.SdkDescribe()})
}

// terseWidgetWith is a widgetWith that opts in to the terse and history filters.
func terseWidgetWith(show ...ssot.SdkShowOption) ssot.Object {
	object := widgetWith(show...)
	object.Generators = append(object.Generators, ssot.WithDataSource(ssot.DataSourceWithTerse(), ssot.DataSourceWithHistory()))
	return object
}

func TestSwitchesSchema(t *testing.T) {
	t.Run("describe only", func(t *testing.T) {
		got := renderTemplate(t, "switchesSchema", modelOf(widgetWith()))
		require.Contains(t, got, `"with_describe": withDescribeSchema("WIDGET", "widget", "WIDGETS"),`)
		require.NotContains(t, got, "with_parameters")
	})
	t.Run("describe and parameters", func(t *testing.T) {
		got := renderTemplate(t, "switchesSchema", modelOf(syntheticObject(
			[]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe(), ssot.SdkShowParameters()},
		)))
		require.Contains(t, got, `"with_describe": withDescribeSchema("WIDGET", "widget", "WIDGETS"),`)
		require.Contains(t, got, `"with_parameters": withParametersSchema("WIDGET", "widget", "WIDGETS"),`)
	})
	t.Run("nothing to switch", func(t *testing.T) {
		got := renderTemplate(t, "switchesSchema", modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()})))
		require.Empty(t, strings.TrimSpace(got))
	})
	t.Run("opting out of describe keeps the parameters switch", func(t *testing.T) {
		model := modelOf(syntheticObject(
			[]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe(), ssot.SdkShowParameters()},
			ssot.WithDataSourceDescribeInSdk(ssot.Omitted()),
		))
		got := renderTemplate(t, "switchesSchema", model)
		require.NotContains(t, got, "with_describe")
		require.Contains(t, got, "with_parameters")
	})
	t.Run("hand-written", func(t *testing.T) {
		model := modelOf(syntheticObject(
			[]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe()},
			ssot.WithDataSourceSwitchesSchema(ssot.ExtCall()),
		))
		require.Contains(t, renderTemplate(t, "switchesSchema", model), "widgetsSwitchesSchemaExt,")
	})
}

func TestFiltersSchema(t *testing.T) {
	t.Run("like, starts with and limit", func(t *testing.T) {
		got := renderTemplate(t, "filtersSchema", modelOf(widgetWith(
			ssot.SdkShowOptionalLike(), ssot.SdkShowOptionalStartsWith(), ssot.SdkShowOptionalLimitFrom(),
		)))
		require.Contains(t, got, `"like": likeSchema,`)
		require.Contains(t, got, `"starts_with": startsWithSchema,`)
		require.Contains(t, got, `"limit": limitFromSchema,`)
		require.NotContains(t, got, `"in"`)
	})
	t.Run("like and extended in", func(t *testing.T) {
		got := renderTemplate(t, "filtersSchema", modelOf(widgetWith(
			ssot.SdkShowOptionalLike(), ssot.SdkShowOptionalExtendedIn(),
		)))
		require.Contains(t, got, `"like": likeSchema,`)
		require.Contains(t, got, `"in": extendedInSchema,`)
		require.NotContains(t, got, "starts_with")
		require.NotContains(t, got, "limit")
	})
	t.Run("plain in", func(t *testing.T) {
		got := renderTemplate(t, "filtersSchema", modelOf(widgetWith(ssot.SdkShowOptionalIn())))
		require.Contains(t, got, `"in": inSchema,`)
		require.NotContains(t, got, "extendedInSchema")
	})
	t.Run("extended in and on", func(t *testing.T) {
		got := renderTemplate(t, "filtersSchema", modelOf(widgetWith(
			ssot.SdkShowOptionalExtendedIn(), ssot.SdkShowOptionalOn(),
		)))
		require.Contains(t, got, `"in": extendedInSchema,`)
		require.Contains(t, got, `"on": onSchema,`)
	})
	t.Run("terse and history", func(t *testing.T) {
		got := renderTemplate(t, "filtersSchema", modelOf(terseWidgetWith(
			ssot.SdkShowTerse(), ssot.SdkShowOptionalSql("HISTORY"),
		)))
		require.Contains(t, got, `"terse": terseSchema,`)
		require.Contains(t, got, `"history": historySchema,`)
	})
	t.Run("terse and history are left out by default", func(t *testing.T) {
		got := renderTemplate(t, "filtersSchema", modelOf(widgetWith(ssot.SdkShowOptionalLike(), ssot.SdkShowTerse(), ssot.SdkShowOptionalSql("HISTORY"))))
		require.NotContains(t, got, "terse")
		require.NotContains(t, got, "history")
	})
	t.Run("in without a schema level", func(t *testing.T) {
		got := renderTemplate(t, "filtersSchema", modelOf(widgetWith(ssot.SdkShowOptionalExtendedInWithoutSchema())))
		require.Contains(t, got, `"in": extendedInWithoutSchemaSchema,`)
		require.NotContains(t, got, "extendedInSchema")
	})
	t.Run("in account", func(t *testing.T) {
		got := renderTemplate(t, "filtersSchema", modelOf(widgetWith(ssot.SdkShowOptionalInAccount())))
		require.Contains(t, got, `"in_account": {`)
		require.Contains(t, got, "Limits results to widgets in this account")
		require.Contains(t, got, "[SHOW WIDGETS](https://docs.snowflake.com/en/sql-reference/sql/show-widgets) `IN ACCOUNT`.")
		require.Contains(t, got, "resources.IsValidAccountIdentifier()")
		require.NotContains(t, got, "likeSchema")
	})
	t.Run("no filters", func(t *testing.T) {
		got := renderTemplate(t, "filtersSchema", modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()})))
		require.Empty(t, strings.TrimSpace(got))
	})
	t.Run("hand-written", func(t *testing.T) {
		model := modelOf(syntheticObject(
			[]ssot.SdkOption{ssot.SdkShow(ssot.SdkShowOptionalLike())},
			ssot.WithDataSourceFiltersSchema(ssot.ExtCall()),
		))
		got := renderTemplate(t, "filtersSchema", model)
		require.Contains(t, got, "widgetsFiltersSchemaExt,")
		require.NotContains(t, got, "likeSchema")
	})
}

func TestItemSchema(t *testing.T) {
	t.Run("show and describe", func(t *testing.T) {
		got := renderTemplate(t, "itemSchema", modelOf(widgetWith()))
		require.Contains(t, got, `"widgets": {`)
		require.Contains(t, got, "Description: \"Holds the output of SHOW WIDGETS.\"")
		require.Contains(t, got, "schemas.ShowWidgetSchema")
		require.Contains(t, got, "Description: \"Holds the output of DESCRIBE WIDGET.\"")
		require.Contains(t, got, "schemas.DescribeWidgetDetailsSchema")
		require.NotContains(t, got, "ParametersAttributeName")
	})
	t.Run("parameters", func(t *testing.T) {
		got := renderTemplate(t, "itemSchema", modelOf(syntheticObject(
			[]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe(), ssot.SdkShowParameters()},
		)))
		require.Contains(t, got, "Description: \"Holds the output of SHOW PARAMETERS FOR WIDGET.\"")
		require.Contains(t, got, "schemas.ShowWidgetParametersSchema")
	})
	t.Run("no describe", func(t *testing.T) {
		got := renderTemplate(t, "itemSchema", modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()})))
		require.Contains(t, got, "schemas.ShowWidgetSchema")
		require.NotContains(t, got, "DescribeOutputAttributeName")
	})
	t.Run("the result of a show-only data source is the show output", func(t *testing.T) {
		got := renderTemplate(t, "itemSchema", modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()})))
		require.Contains(t, got, `Description: "Holds the output of SHOW WIDGETS.",`)
		require.NotContains(t, got, "details queries")
	})
	t.Run("the result with describe or parameters aggregates the details", func(t *testing.T) {
		require.Contains(t, renderTemplate(t, "itemSchema", modelOf(widgetWith())), "Holds the aggregated output of all widget details queries.")
		parameters := syntheticObject([]ssot.SdkOption{ssot.SdkShow(), ssot.SdkShowParameters()})
		require.Contains(t, renderTemplate(t, "itemSchema", modelOf(parameters)), "Holds the aggregated output of all widget details queries.")
	})
	t.Run("a hand-written describe schema still counts as details", func(t *testing.T) {
		model := modelOf(syntheticObject(
			[]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe()},
			ssot.WithDataSourceDescribeOutputSchema(ssot.ExtCall()),
		))
		require.Contains(t, renderTemplate(t, "itemSchema", model), "Holds the aggregated output of all widget details queries.")
	})
	t.Run("describe details name", func(t *testing.T) {
		object := syntheticObject([]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe()})
		object.DescribeDetailsName = "WidgetSettings"
		require.Contains(t, renderTemplate(t, "itemSchema", modelOf(object)), "schemas.DescribeWidgetSettingsSchema")
	})
	t.Run("hand-written", func(t *testing.T) {
		model := modelOf(syntheticObject(
			[]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe()},
			ssot.WithDataSourceShowOutputSchema(ssot.ExtCall()),
			ssot.WithDataSourceDescribeOutputSchema(ssot.ExtCall()),
		))
		got := renderTemplate(t, "itemSchema", model)
		require.Contains(t, got, "widgetsShowOutputSchemaExt,")
		require.Contains(t, got, "widgetsDescribeOutputSchemaExt,")
		require.NotContains(t, got, "schemas.")
	})
}

func TestDataSourceSchema(t *testing.T) {
	t.Run("parts keep their order", func(t *testing.T) {
		model := modelOf(syntheticObject(
			[]ssot.SdkOption{ssot.SdkShow(ssot.SdkShowOptionalLike()), ssot.SdkDescribe()},
			ssot.WithDataSourceAdditionalSchema(ssot.ExtCall()),
		))
		var b strings.Builder
		require.NoError(t, DataSourceTemplate.Execute(&b, model))
		got := b.String()

		require.Contains(t, got, "var widgetsSchema = collections.MergeMaps(")
		order := []string{"with_describe", "likeSchema", `"widgets": {`, "widgetsAdditionalSchemaExt,"}
		last := -1
		for _, part := range order {
			index := strings.Index(got, part)
			require.Greater(t, index, last, part)
			last = index
		}
	})
}

func TestAdditionalOutputDescriptions(t *testing.T) {
	object := func(opts ...ssot.DataSourceOption) ssot.Object {
		obj := syntheticObject([]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe(), ssot.SdkShowParameters()})
		obj.Generators = append(obj.Generators, ssot.WithDataSource(opts...))
		return obj
	}

	t.Run("none by default", func(t *testing.T) {
		model := modelOf(object())
		require.Contains(t, renderTemplate(t, "showOutputSchema", model), `Description: "Holds the output of SHOW WIDGETS.",`)
		require.Contains(t, renderTemplate(t, "describeOutputSchema", model), `Description: "Holds the output of DESCRIBE WIDGET.",`)
		require.Contains(t, renderTemplate(t, "parametersOutputSchema", model), `Description: "Holds the output of SHOW PARAMETERS FOR WIDGET.",`)
	})
	t.Run("appended to the matching block only", func(t *testing.T) {
		model := modelOf(object(
			ssot.DataSourceAdditionalShowOutputDescription("The full list is read from `VIEW`."),
			ssot.DataSourceAdditionalDescribeOutputDescription("Needs a privilege."),
			ssot.DataSourceAdditionalParametersOutputDescription("Read per object."),
		))
		require.Contains(t, renderTemplate(t, "showOutputSchema", model), "Holds the output of SHOW WIDGETS. The full list is read from `VIEW`.")
		require.Contains(t, renderTemplate(t, "describeOutputSchema", model), "Holds the output of DESCRIBE WIDGET. Needs a privilege.")
		require.Contains(t, renderTemplate(t, "parametersOutputSchema", model), "Holds the output of SHOW PARAMETERS FOR WIDGET. Read per object.")
		require.NotContains(t, renderTemplate(t, "describeOutputSchema", model), "VIEW")
	})
	t.Run("quotes are escaped", func(t *testing.T) {
		model := modelOf(object(ssot.DataSourceAdditionalShowOutputDescription(`Says "hello".`)))
		require.Contains(t, renderTemplate(t, "showOutputSchema", model), `Holds the output of SHOW WIDGETS. Says \"hello\".`)
	})
}
