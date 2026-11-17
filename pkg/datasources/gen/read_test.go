package gen

import (
	"go/format"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func renderDataSource(t *testing.T, model DataSourceModel) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, DataSourceTemplate.Execute(&b, model))
	return b.String()
}

func TestReadWrapper(t *testing.T) {
	t.Run("stable", func(t *testing.T) {
		got := renderTemplate(t, "readWrapper", modelOf(widgetWith()))
		require.Equal(t, "TrackingReadWrapper(datasources.Widgets, ReadWidgets)", got)
	})
	t.Run("preview", func(t *testing.T) {
		object := widgetWith()
		object.PreviewFeature = true
		got := renderTemplate(t, "readWrapper", modelOf(object))
		require.Equal(t, "PreviewFeatureReadWrapper(string(previewfeatures.WidgetsDatasource), TrackingReadWrapper(datasources.Widgets, ReadWidgets))", got)
	})
	t.Run("hand-written", func(t *testing.T) {
		model := modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()}, ssot.WithDataSourceReadWrapper(ssot.ExtCall())))
		require.Equal(t, "widgetsReadWrapperExt(ReadWidgets)", renderTemplate(t, "readWrapper", model))
	})
	t.Run("omitted", func(t *testing.T) {
		model := modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()}, ssot.WithDataSourceReadWrapper(ssot.Omitted())))
		require.Equal(t, "ReadWidgets", renderTemplate(t, "readWrapper", model))
	})
}

func TestDescription(t *testing.T) {
	t.Run("show and describe", func(t *testing.T) {
		got := renderTemplate(t, "description", modelOf(widgetWith()))
		require.Contains(t, got, "Data source used to get details of filtered widgets.")
		require.Contains(t, got, "[SHOW WIDGETS](https://docs.snowflake.com/en/sql-reference/sql/show-widgets) query.")
		require.Contains(t, got, "The results of SHOW and DESCRIBE are encapsulated in one output collection `widgets`.")
		require.NotContains(t, got, "docs.WithPreviewFeatureMark()")
	})
	t.Run("preview with parameters", func(t *testing.T) {
		object := syntheticObject([]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe(), ssot.SdkShowParameters()})
		object.PreviewFeature = true
		got := renderTemplate(t, "description", modelOf(object))
		require.True(t, strings.HasPrefix(got, "docs.WithPreviewFeatureMark() + "))
		require.Contains(t, got, "The results of SHOW, DESCRIBE, and SHOW PARAMETERS are encapsulated in one output collection `widgets`.")
	})
	t.Run("show only", func(t *testing.T) {
		got := renderTemplate(t, "description", modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()})))
		require.Contains(t, got, "The results of SHOW are encapsulated")
	})
	t.Run("hand-written", func(t *testing.T) {
		model := modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()}, ssot.WithDataSourceDescription(ssot.ExtCall())))
		require.Equal(t, "widgetsDescriptionExt", renderTemplate(t, "description", model))
	})
}

func TestApplyFilters(t *testing.T) {
	t.Run("like, starts with and limit", func(t *testing.T) {
		got := renderTemplate(t, "applyFilters", modelOf(widgetWith(
			ssot.SdkShowOptionalLike(), ssot.SdkShowOptionalStartsWith(), ssot.SdkShowOptionalLimitFrom(),
		)))
		require.Contains(t, got, "handleLike(d, &req.Like)")
		require.Contains(t, got, "handleStartsWith(d, &req.StartsWith)")
		require.Contains(t, got, "handleLimitFrom(d, &req.Limit)")
		require.NotContains(t, got, "handleIn")
	})
	t.Run("extended in", func(t *testing.T) {
		got := renderTemplate(t, "applyFilters", modelOf(widgetWith(
			ssot.SdkShowOptionalLike(), ssot.SdkShowOptionalExtendedIn(),
		)))
		require.Contains(t, got, "handleExtendedIn(d, &req.In)")
		require.NotContains(t, got, "handleStartsWith")
	})
	t.Run("plain in", func(t *testing.T) {
		got := renderTemplate(t, "applyFilters", modelOf(widgetWith(ssot.SdkShowOptionalIn())))
		require.Contains(t, got, "handleIn(d, &req.In)")
		require.NotContains(t, got, "handleExtendedIn")
	})
	t.Run("on", func(t *testing.T) {
		got := renderTemplate(t, "applyFilters", modelOf(widgetWith(ssot.SdkShowOptionalOn())))
		require.Contains(t, got, "handleOn(d, &req.On)")
	})
	t.Run("terse and history", func(t *testing.T) {
		got := renderTemplate(t, "applyFilters", modelOf(terseWidgetWith(
			ssot.SdkShowTerse(), ssot.SdkShowOptionalSql("HISTORY"),
		)))
		require.Contains(t, got, `d.Get("terse").(bool)`)
		require.Contains(t, got, "req.WithTerse(true)")
		require.Contains(t, got, `d.Get("history").(bool)`)
		require.Contains(t, got, "req.WithHistory(true)")
	})
	t.Run("in account", func(t *testing.T) {
		got := renderTemplate(t, "applyFilters", modelOf(widgetWith(ssot.SdkShowOptionalInAccount())))
		require.Contains(t, got, `d.GetOk("in_account")`)
		require.Contains(t, got, "helpers.DecodeSnowflakeAccountIdentifier")
		require.Contains(t, got, "req.WithInAccount(account)")
	})
	t.Run("no filters", func(t *testing.T) {
		got := renderTemplate(t, "applyFilters", modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()})))
		require.Empty(t, strings.TrimSpace(got))
	})
	t.Run("hand-written", func(t *testing.T) {
		model := modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow(ssot.SdkShowOptionalLike())}, ssot.WithDataSourceApplyFilters(ssot.ExtCall())))
		got := renderTemplate(t, "applyFilters", model)
		require.Contains(t, got, "widgetsApplyFiltersExt(d, req)")
		require.NotContains(t, got, "handleLike")
	})
}

func TestShowInSdkAndSetId(t *testing.T) {
	require.Contains(t, renderTemplate(t, "showInSdk", modelOf(widgetWith())), "items, err := client.Widgets.Show(ctx, req)")
	require.Contains(t, renderTemplate(t, "setId", modelOf(widgetWith())), `d.SetId("widgets_read")`)

	model := modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()},
		ssot.WithDataSourceShowInSdk(ssot.ExtCall()), ssot.WithDataSourceSetId(ssot.ExtCall())))
	require.Contains(t, renderTemplate(t, "showInSdk", model), "items, err := widgetsShowInSdkExt(ctx, client, req, &diags)")
	require.Contains(t, renderTemplate(t, "setId", model), "widgetsSetIdExt(d)")
}

func TestDescribe(t *testing.T) {
	t.Run("single result", func(t *testing.T) {
		model := modelOf(widgetWith())
		require.Contains(t, renderTemplate(t, "describeInSdk", model), "describeResult, err := client.Widgets.Describe(ctx, item.ID())")
		require.Contains(t, renderTemplate(t, "describeOutputSet", model), "describeOutput = []map[string]any{schemas.WidgetDetailsToSchema(describeResult)}")
	})
	t.Run("details method is preferred", func(t *testing.T) {
		model := modelOf(syntheticObject([]ssot.SdkOption{
			ssot.SdkShow(), ssot.SdkDescribe(ssot.SdkDescribeDetails(ssot.SdkDescribeDetailsMethod())),
		}))
		require.Contains(t, renderTemplate(t, "describeInSdk", model), "client.Widgets.DescribeDetails(ctx, item.ID())")
	})
	t.Run("list result", func(t *testing.T) {
		model := modelOf(syntheticObject([]ssot.SdkOption{
			ssot.SdkShow(), ssot.SdkDescribe(ssot.SdkDescribeMapping(ssot.SdkDescriptionMappingSlice)),
		}))
		got := renderTemplate(t, "describeOutputSet", model)
		require.Contains(t, got, "collections.Map(describeResult, func(details sdk.WidgetDetails) map[string]any {")
		require.Contains(t, got, "schemas.WidgetDetailsToSchema(&details)")
	})
	t.Run("hand-written", func(t *testing.T) {
		model := modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe()},
			ssot.WithDataSourceDescribeInSdk(ssot.ExtCall()), ssot.WithDataSourceDescribeOutputSet(ssot.ExtCall())))
		require.Contains(t, renderTemplate(t, "describeInSdk", model), "describeResult, err := widgetsDescribeInSdkExt(ctx, client, item, &diags)")
		require.Contains(t, renderTemplate(t, "describeOutputSet", model), "describeOutput = widgetsDescribeOutputSetExt(describeResult)")
	})
	t.Run("nothing without describe", func(t *testing.T) {
		model := modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()}))
		require.Empty(t, strings.TrimSpace(renderTemplate(t, "describeInSdk", model)))
		require.Empty(t, strings.TrimSpace(renderTemplate(t, "describeOutputSet", model)))
	})
}

func TestParameters(t *testing.T) {
	model := modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow(), ssot.SdkShowParameters()}))
	require.Contains(t, renderTemplate(t, "showParametersInSdk", model), "parameters, err := client.Widgets.ShowParameters(ctx, item.ID())")
	require.Contains(t, renderTemplate(t, "parametersOutputSet", model), "schemas.WidgetParametersToSchema(parameters, providerCtx)")

	across := syntheticObject([]ssot.SdkOption{ssot.SdkShow(), ssot.SdkShowParameters()})
	across.ShowParametersDetailsAdditionalLevels = []ssot.ParameterLevel{ssot.ParameterLevel("WIDGET_EXTRA")}
	require.Contains(t, renderTemplate(t, "parametersOutputSet", modelOf(across)), "schemas.AllWidgetParametersToSchema(parameters, providerCtx)")
	require.Contains(t, renderTemplate(t, "parametersOutputSchema", modelOf(across)), "schemas.ShowAllWidgetParametersSchema(),")

	custom := modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow(), ssot.SdkShowParameters()},
		ssot.WithDataSourceShowParametersInSdk(ssot.ExtCall()), ssot.WithDataSourceParametersOutputSet(ssot.ExtCall())))
	require.Contains(t, renderTemplate(t, "showParametersInSdk", custom), "widgetsShowParametersInSdkExt(ctx, client, item, &diags)")
	require.Contains(t, renderTemplate(t, "parametersOutputSet", custom), "widgetsParametersOutputSetExt(parameters)")
}

func TestRead(t *testing.T) {
	t.Run("show and describe", func(t *testing.T) {
		got := renderDataSource(t, modelOf(widgetWith()))
		require.Contains(t, got, "func Widgets() *schema.Resource {")
		require.Contains(t, got, "func ReadWidgets(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {")
		require.Contains(t, got, "client := meta.(*provider.Context).Client")
		require.NotContains(t, got, "providerCtx")
		require.Contains(t, got, "req := sdk.NewShowWidgetRequest()")
		require.Contains(t, got, `if d.Get("with_describe").(bool) {`)
		require.Contains(t, got, "resources.ShowOutputAttributeName: []map[string]any{schemas.WidgetToSchema(&item)},")
		require.Contains(t, got, "resources.DescribeOutputAttributeName: describeOutput,")
		require.NotContains(t, got, "ParametersAttributeName")
		require.Contains(t, got, `d.Set("widgets", flattened)`)
	})
	t.Run("parameters need the provider context", func(t *testing.T) {
		got := renderDataSource(t, modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow(), ssot.SdkShowParameters()})))
		require.Contains(t, got, "providerCtx := meta.(*provider.Context)")
		require.Contains(t, got, "client := providerCtx.Client")
		require.Contains(t, got, `if d.Get("with_parameters").(bool) {`)
		require.Contains(t, got, "resources.ParametersAttributeName: parametersOutput,")
	})
	t.Run("show only", func(t *testing.T) {
		got := renderDataSource(t, modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()})))
		require.NotContains(t, got, "with_describe")
		require.NotContains(t, got, "describeOutput")
		require.NotContains(t, got, "DescribeOutputAttributeName")
	})
	t.Run("extra item output is merged", func(t *testing.T) {
		model := modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()}, ssot.WithDataSourceItemExtra(ssot.ExtCall())))
		got := renderDataSource(t, model)
		require.Contains(t, got, "itemExtra, err := widgetsItemExtraExt(ctx, d, client, item, &diags)")
		require.Contains(t, got, "collections.MergeMaps(map[string]any{")
		require.Contains(t, got, "}, itemExtra)")
	})
	t.Run("a hand-written step can report warnings", func(t *testing.T) {
		model := modelOf(syntheticObject([]ssot.SdkOption{ssot.SdkShow()}, ssot.WithDataSourceShowInSdk(ssot.ExtCall())))
		got := renderDataSource(t, model)
		require.Contains(t, got, "var diags diag.Diagnostics")
		require.Contains(t, got, "return append(diags, diag.FromErr(err)...)")
		require.Contains(t, got, "return diags")
	})
	t.Run("no diagnostics without hand-written steps", func(t *testing.T) {
		got := renderDataSource(t, modelOf(widgetWith()))
		require.NotContains(t, got, "diags")
	})
	t.Run("no merge without extras", func(t *testing.T) {
		require.NotContains(t, renderDataSource(t, modelOf(widgetWith())), "itemExtra")
	})
}

func TestRead_everyShapeIsValidGo(t *testing.T) {
	preview := widgetWith()
	preview.PreviewFeature = true
	sliceDescribe := syntheticObject([]ssot.SdkOption{
		ssot.SdkShow(), ssot.SdkDescribe(ssot.SdkDescribeMapping(ssot.SdkDescriptionMappingSlice)),
	})
	detailsDescribe := syntheticObject([]ssot.SdkOption{
		ssot.SdkShow(), ssot.SdkDescribe(ssot.SdkDescribeDetails(ssot.SdkDescribeDetailsMethod())),
	})
	handWritten := syntheticObject(
		[]ssot.SdkOption{ssot.SdkShow(ssot.SdkShowOptionalLike()), ssot.SdkDescribe(), ssot.SdkShowParameters()},
		ssot.WithDataSourceDescription(ssot.ExtCall()),
		ssot.WithDataSourceReadWrapper(ssot.ExtCall()),
		ssot.WithDataSourceFiltersSchema(ssot.ExtCall()),
		ssot.WithDataSourceSwitchesSchema(ssot.ExtCall()),
		ssot.WithDataSourceShowOutputSchema(ssot.ExtCall()),
		ssot.WithDataSourceDescribeOutputSchema(ssot.ExtCall()),
		ssot.WithDataSourceParametersOutputSchema(ssot.ExtCall()),
		ssot.WithDataSourceAdditionalSchema(ssot.ExtCall()),
		ssot.WithDataSourceApplyFilters(ssot.ExtCall()),
		ssot.WithDataSourceShowInSdk(ssot.ExtCall()),
		ssot.WithDataSourceSetId(ssot.ExtCall()),
		ssot.WithDataSourceDescribeInSdk(ssot.ExtCall()),
		ssot.WithDataSourceDescribeOutputSet(ssot.ExtCall()),
		ssot.WithDataSourceShowParametersInSdk(ssot.ExtCall()),
		ssot.WithDataSourceParametersOutputSet(ssot.ExtCall()),
		ssot.WithDataSourceItemExtra(ssot.ExtCall()),
	)

	objects := map[string]ssot.Object{
		"show only":          syntheticObject([]ssot.SdkOption{ssot.SdkShow()}),
		"show and describe":  widgetWith(),
		"preview":            preview,
		"slice describe":     sliceDescribe,
		"details describe":   detailsDescribe,
		"parameters":         syntheticObject([]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe(), ssot.SdkShowParameters()}),
		"every filter":       terseWidgetWith(ssot.SdkShowOptionalLike(), ssot.SdkShowOptionalStartsWith(), ssot.SdkShowOptionalLimitFrom(), ssot.SdkShowOptionalExtendedIn(), ssot.SdkShowOptionalOn(), ssot.SdkShowTerse(), ssot.SdkShowOptionalSql("HISTORY")),
		"plain in":           widgetWith(ssot.SdkShowOptionalIn()),
		"in account":         widgetWith(ssot.SdkShowOptionalInAccount()),
		"every hook omitted": syntheticObject([]ssot.SdkOption{ssot.SdkShow(), ssot.SdkDescribe(), ssot.SdkShowParameters()}, ssot.WithDataSourceDescribeInSdk(ssot.Omitted()), ssot.WithDataSourceShowParametersInSdk(ssot.Omitted()), ssot.WithDataSourceReadWrapper(ssot.Omitted())),
		"hand-written":       handWritten,
	}
	for name, object := range objects {
		t.Run(name, func(t *testing.T) {
			_, err := format.Source([]byte("package datasources\n\n" + renderDataSource(t, modelOf(object))))
			require.NoError(t, err)
		})
	}
}
