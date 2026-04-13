package genhelpers

import "os"

type ExplicitImport struct {
	Name string
	Path string
}

type PreambleModel struct {
	GeneratorName    string
	GeneratorVersion string
	PackageName      string
	ExplicitImports  []ExplicitImport
}

func NewPreambleModel(name string, version string) *PreambleModel {
	return NewPreambleModelWithImports(name, version, []string{})
}

func NewPreambleModelWithImports(name string, version string, imports []string) *PreambleModel {
	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	m := &PreambleModel{
		GeneratorName:    name,
		GeneratorVersion: version,
		PackageName:      packageWithGenerateDirective,
	}
	for _, imp := range imports {
		m.WithImport(imp)
	}
	return m
}

func (m *PreambleModel) WithImport(path string) *PreambleModel {
	return m.WithNamedImport("", path)
}

func (m *PreambleModel) WithNamedImport(name string, path string) *PreambleModel {
	m.ExplicitImports = append(m.ExplicitImports, ExplicitImport{
		Name: name,
		Path: path,
	})
	return m
}

func (m *PreambleModel) getPreambleModel() *PreambleModel { return m }

type HasPreambleModel interface {
	getPreambleModel() *PreambleModel
}

// ObjectGenerationSettings holds per-object generation configuration.
// Embed this in input object types to enable per-object generation part filtering.
type ObjectGenerationSettings struct {
	// AllowedGenerationParts lists which generation parts apply to this object.
	// nil/empty means "use the generator's default" (all globally-filtered parts).
	AllowedGenerationParts []string
}

type HasObjectGenerationSettings interface {
	getObjectGenerationSettings() *ObjectGenerationSettings
}

func (i *ObjectGenerationSettings) getObjectGenerationSettings() *ObjectGenerationSettings { return i }
