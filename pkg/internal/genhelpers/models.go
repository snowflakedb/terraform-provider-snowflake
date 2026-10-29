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
	AllowedGenerationParts []GenerationPartNamer
	// EnabledGenerationParts lists optional (disabled-by-default) generation parts that should be enabled for this object.
	EnabledGenerationParts []GenerationPartNamer
}

type HasObjectGenerationSettings interface {
	getObjectGenerationSettings() *ObjectGenerationSettings
}

func (i *ObjectGenerationSettings) getObjectGenerationSettings() *ObjectGenerationSettings { return i }

// GenerationPartNamer is implemented by typed generation part name constants.
// Each generator package defines its own private type satisfying this interface,
// ensuring that only predefined constants can be passed to generation part registration methods.
type GenerationPartNamer interface {
	GenerationPartName() string
}

type defaultGenerationPartName string

func (d defaultGenerationPartName) GenerationPartName() string { return string(d) }

const DefaultGenerationPart defaultGenerationPartName = "default"

// SdkValueType is the richer Go type generated beside the Terraform storage type.
// SdkType is empty for schema-reflected fields. SdkKind is "bool", "enum", or "identifier" when SdkType is set.
type SdkValueType struct {
	SdkType string
	SdkKind string
}

// SdkMethodSuffix names the extra method that takes the real Go type.
// The storage-type method stays. Go cannot overload it, so the real-type method needs its own suffix: Bool for a Go bool, Sdk for an SDK type.
func SdkMethodSuffix(sdkType string) string {
	if sdkType == "bool" {
		return "Bool"
	}
	return "Sdk"
}

// SdkValueMapper selects the mapper that stores a richer value as the Terraform string.
func SdkValueMapper(kind string) Mapper {
	switch kind {
	case "bool":
		return FormatBool
	case "identifier":
		return Name
	default:
		return CastToString
	}
}
