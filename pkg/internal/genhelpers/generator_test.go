package genhelpers

import (
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

type parallelTestObject struct {
	name string
}

func (o parallelTestObject) ObjectName() string { return o.name }

type parallelTestModel struct {
	Name string
	*PreambleModel
}

func TestGenerateAndSaveParallelMatchesSequential(t *testing.T) {
	objects := []parallelTestObject{{name: "Alpha"}, {name: "Beta"}}
	sequential := generateTestObjects(t, objects, 1)
	parallel := generateTestObjects(t, objects, 4)

	require.Equal(t, sequential, parallel)
	require.Contains(t, sequential["alpha_gen.go"], "Alpha")
	require.Contains(t, sequential["beta_gen.go"], "Beta")
}

func generateTestObjects(t *testing.T, objects []parallelTestObject, parallelism int) map[string]string {
	t.Helper()

	body := template.Must(template.New("body").Parse("package gentest\n\nfunc {{ .Name }}Value() string {\n\treturn \"{{ .Name }}\"\n}\n"))
	preamble := NewPreambleModel("test generator", "0.0.0")
	generator := NewGenerator(
		preamble,
		func() []parallelTestObject { return objects },
		func(object parallelTestObject, preamble *PreambleModel) parallelTestModel {
			return parallelTestModel{Name: object.name, PreambleModel: preamble}
		},
		func(_ parallelTestObject, model parallelTestModel) string {
			return strings.ToLower(model.Name) + "_gen.go"
		},
		[]*template.Template{body},
	)
	generator.parallelism = parallelism

	directory := t.TempDir()
	t.Chdir(directory)
	require.NoError(t, generator.generateAndSave(objects, generator.generationParts))

	entries, err := os.ReadDir(directory)
	require.NoError(t, err)

	root, err := os.OpenRoot(directory)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, root.Close())
	})

	files := make(map[string]string, len(entries))
	for _, entry := range entries {
		content, err := root.ReadFile(entry.Name())
		require.NoError(t, err)
		files[entry.Name()] = string(content)
	}
	return files
}
