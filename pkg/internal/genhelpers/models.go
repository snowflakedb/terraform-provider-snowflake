package genhelpers

import "os"

type PreambleModel struct {
	GeneratorName    string
	GeneratorVersion string
	PackageName      string
	ExplicitImports  []string
}

func NewPreambleModel(name string, version string, imports []string) *PreambleModel {
	packageWithGenerateDirective := os.Getenv("GOPACKAGE")
	return &PreambleModel{
		GeneratorName:    name,
		GeneratorVersion: version,
		PackageName:      packageWithGenerateDirective,
		ExplicitImports:  imports,
	}
}

func (m *PreambleModel) getPreambleModel() *PreambleModel { return m }

type HasPreambleModel interface {
	getPreambleModel() *PreambleModel
}
