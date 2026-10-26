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

func Test_ShowProjection_panicsMissingField(t *testing.T) {
	source := PlainStruct("Widget").Text("Name")
	target := PlainStruct("WidgetMini").Text("Name").Text("Comment")
	require.Panics(t, func() {
		NewInterface("Widgets", "Widget", "AccountObjectIdentifier").ShowProjection(source, "AsMini", target)
	})
}
