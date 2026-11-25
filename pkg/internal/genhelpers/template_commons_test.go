package genhelpers

import (
	"strings"
	"testing"
	"testing/fstest"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Unimplemented(t *testing.T) {
	require.PanicsWithValue(t, "generated implementation for ParseIdFromConfig is not implemented", func() {
		Unimplemented("ParseIdFromConfig")
	})
}

func Test_SubTemplateName(t *testing.T) {
	assert.Equal(t, "parseId", SubTemplateName("parse_id.tmpl"))
	assert.Equal(t, "showByIdSafelyInSdk", SubTemplateName("show_by_id_safely_in_sdk.tmpl"))
	assert.Equal(t, "fqnSet", SubTemplateName("fqn_set.tmpl"))
	assert.Equal(t, "item", SubTemplateName("item.tmpl"))
}

func Test_ParseSubTemplates(t *testing.T) {
	fsys := fstest.MapFS{
		"sub/parse_id.tmpl":     {Data: []byte(`id of {{ . }}`)},
		"sub/show_by_id.tmpl":   {Data: []byte(`{{ template "parseId" . }} then show`)},
		"sub/notes.txt":         {Data: []byte(`not a template`)},
		"sub/nested/other.tmpl": {Data: []byte(`not read`)},
	}
	root := template.New("root")
	ParseSubTemplates(root, fsys, "sub")

	var out strings.Builder
	require.NoError(t, root.ExecuteTemplate(&out, "showById", "alert"))
	assert.Equal(t, "id of alert then show", out.String())
	assert.Nil(t, root.Lookup("notes"))
	assert.Nil(t, root.Lookup("other"))

	assert.Panics(t, func() { ParseSubTemplates(template.New("root"), fsys, "missing") })
	assert.Panics(t, func() {
		ParseSubTemplates(template.New("root"), fstest.MapFS{"sub/broken.tmpl": {Data: []byte(`{{ if }}`)}}, "sub")
	})
}
