package gen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShowProjection_emitsMethod(t *testing.T) {
	source := PlainStruct("Widget").Text("Name").OptionalNumber("Running").Text("Comment")
	target := PlainStruct("WidgetMini").Text("Name").OptionalNumber("Running")
	iface := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").
		ShowProjection(source, "AsMini", target)

	var b strings.Builder
	require.NoError(t, OperationStructIterateTemplate.Execute(&b, iface))
	got := b.String()
	require.Contains(t, got, "func (w *Widget) AsMini() *WidgetMini {")
	require.Contains(t, got, "if w == nil {")
	require.Contains(t, got, "return nil")
	require.Contains(t, got, "return &WidgetMini{")
	require.Contains(t, got, "Name: w.Name,")
	require.Contains(t, got, "Running: w.Running,")
	require.NotContains(t, got, "Comment:")
	require.NotContains(t, got, "type WidgetMini struct")
}

func Test_ShowByIDProjection_emitsMethods(t *testing.T) {
	iface := NewInterface("Widgets", "Widget", "AccountObjectIdentifier").
		ShowByIdOperationWithNoFiltering()
	iface.ShowObjectName = "Widget"
	iface.Operations[0].ObjectInterface = iface
	iface.ShowByIDProjection("Mini")

	var ifaceOut strings.Builder
	require.NoError(t, InterfaceTemplate.Execute(&ifaceOut, iface))
	gotIface := ifaceOut.String()
	require.Contains(t, gotIface, "ShowByIDMini(ctx context.Context, id AccountObjectIdentifier) (*WidgetMini, error)")
	require.Contains(t, gotIface, "ShowByIDSafelyMini(ctx context.Context, id AccountObjectIdentifier) (*WidgetMini, error)")

	var implOut strings.Builder
	require.NoError(t, ImplementationTemplate.Execute(&implOut, iface))
	gotImpl := implOut.String()
	require.Contains(t, gotImpl, "func (v *widgets) ShowByIDMini(ctx context.Context, id AccountObjectIdentifier) (*WidgetMini, error) {")
	require.Contains(t, gotImpl, "row, err := v.ShowByID(ctx, id)")
	require.Contains(t, gotImpl, "return nil, err")
	require.Contains(t, gotImpl, "return row.AsMini(), nil")
	require.Contains(t, gotImpl, "func (v *widgets) ShowByIDSafelyMini(ctx context.Context, id AccountObjectIdentifier) (*WidgetMini, error) {")
	require.Contains(t, gotImpl, "return SafeShowById(v.client, v.ShowByIDMini, ctx, id)")
}

func Test_ShowByIDProjection_panics(t *testing.T) {
	withShowByID := func() *Interface {
		return NewInterface("Widgets", "Widget", "AccountObjectIdentifier").ShowByIdOperationWithNoFiltering()
	}
	require.Panics(t, func() { withShowByID().ShowByIDProjection("") })
	require.Panics(t, func() { withShowByID().ShowByIDProjection("Mini") })
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").ShowByIDProjection("Mini")
	})
}

func Test_ShowProjection_panicsMissingField(t *testing.T) {
	source := PlainStruct("Widget").Text("Name")
	target := PlainStruct("WidgetMini").Text("Name").Text("Comment")
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").ShowProjection(source, "AsMini", target)
	})
}
