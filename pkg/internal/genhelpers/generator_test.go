package genhelpers

import (
	"bytes"
	"log"
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
		func(object parallelTestObject, preamble *PreambleModel, _ *log.Logger) parallelTestModel {
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

type testPartName string

func (p testPartName) GenerationPartName() string { return string(p) }

func TestGenerateAndSaveParallelLogBatchesStayGrouped(t *testing.T) {
	var captured bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&captured)
	t.Cleanup(func() { log.SetOutput(previous) })

	objects := []parallelTestObject{{name: "Alpha"}, {name: "Beta"}}
	body := template.Must(template.New("body").Parse("package gentest\n\nfunc {{ .Name }}Value() string {\n\treturn \"{{ .Name }}\"\n}\n"))
	preamble := NewPreambleModel("test generator", "0.0.0")
	generator := NewGenerator(
		preamble,
		func() []parallelTestObject { return objects },
		func(object parallelTestObject, preamble *PreambleModel, logger *log.Logger) parallelTestModel {
			logger.Printf("building %s", object.name)
			return parallelTestModel{Name: object.name, PreambleModel: preamble}
		},
		func(_ parallelTestObject, model parallelTestModel) string {
			return strings.ToLower(model.Name) + "_gen.go"
		},
		[]*template.Template{body},
	)
	skipped := func(_ parallelTestObject, _ parallelTestModel) string { return "skipped_gen.go" }
	never := func(parallelTestObject) bool { return false }
	generator.
		WithObjectPrepare(func(object parallelTestObject, logger *log.Logger) {
			logger.Printf("preparing %s", object.name)
		}).
		WithConditionalGenerationPart(testPartName("first"), skipped, []*template.Template{body}, never).
		WithConditionalGenerationPart(testPartName("second"), skipped, []*template.Template{body}, never)
	generator.parallelism = 4

	t.Chdir(t.TempDir())
	require.NoError(t, generator.generateAndSave(objects, generator.generationParts))

	var order []string
	for line := range strings.SplitSeq(captured.String(), "\n") {
		switch {
		case strings.Contains(line, "Alpha:"):
			order = append(order, "Alpha")
		case strings.Contains(line, "Beta:"):
			order = append(order, "Beta")
		}
	}
	require.Equal(t, []string{"Alpha", "Alpha", "Alpha", "Alpha", "Beta", "Beta", "Beta", "Beta"}, order)
	text := captured.String()
	require.Less(t, strings.Index(text, "preparing Alpha"), strings.Index(text, "building Alpha"))
	require.Less(t, strings.Index(text, "building Alpha"), strings.Index(text, "preparing Beta"))
}
