package genhelpers

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// TODO [SNOW-1501905]: describe
type ObjectNameProvider interface {
	ObjectName() string
}

type Generator[T ObjectNameProvider, M HasPreambleModel] struct {
	objectsProvider func() []T
	modelProvider   func(T, *PreambleModel) M
	// TODO [SNOW-1501905]: consider adding filename to model?
	filenameProvider func(T, M) string
	templates        []*template.Template

	generationParts []GenerationPart[T, M]

	additionalObjectDebugLogProviders []func([]T)
	objectFilters                     []func(T) bool
	generationPartFilters             []func(GenerationPart[T, M]) bool

	preamble *PreambleModel
}

type GenerationPart[T ObjectNameProvider, M HasPreambleModel] struct {
	name             string
	filenameProvider func(T, M) string
	templates        []*template.Template
}

func (g *GenerationPart[_, _]) GetName() string {
	return g.name
}

func NewGenerationPart[T ObjectNameProvider, M HasPreambleModel](name string, filenameProvider func(T, M) string, templates []*template.Template) GenerationPart[T, M] {
	return GenerationPart[T, M]{
		name:             name,
		filenameProvider: filenameProvider,
		templates:        templates,
	}
}

func NewGenerator[T ObjectNameProvider, M HasPreambleModel](preamble *PreambleModel, objectsProvider func() []T, modelProvider func(T, *PreambleModel) M, filenameProvider func(T, M) string, templates []*template.Template) *Generator[T, M] {
	// TODO [SNOW-2324252]: handle vararg input
	parts := []GenerationPart[T, M]{
		// TODO [SNOW-2324252]: change default to name when changing to vararg
		NewGenerationPart("default", filenameProvider, templates),
	}
	return &Generator[T, M]{
		objectsProvider:  objectsProvider,
		modelProvider:    modelProvider,
		filenameProvider: filenameProvider,
		templates:        templates,

		generationParts: parts,

		additionalObjectDebugLogProviders: make([]func([]T), 0),
		objectFilters:                     make([]func(T) bool, 0),
		generationPartFilters:             make([]func(GenerationPart[T, M]) bool, 0),

		preamble: preamble,
	}
}

// TODO [SNOW-2324252]: Probably remove later when we have vararg support in the NewGenerator constructor
func (g *Generator[T, M]) WithGenerationPart(partName string, filenameProvider func(T, M) string, templates []*template.Template) *Generator[T, M] {
	g.generationParts = append(g.generationParts, NewGenerationPart(partName, filenameProvider, templates))
	return g
}

func (g *Generator[T, M]) WithAdditionalObjectsDebugLogs(objectLogsProvider func([]T)) *Generator[T, M] {
	g.additionalObjectDebugLogProviders = append(g.additionalObjectDebugLogProviders, objectLogsProvider)
	return g
}

func (g *Generator[T, M]) WithAdditionalObjectFilter(objectFilter func(T) bool) *Generator[T, M] {
	g.objectFilters = append(g.objectFilters, objectFilter)
	return g
}

func (g *Generator[T, M]) WithAdditionalGenerationPartFilter(generationPartFilter func(GenerationPart[T, M]) bool) *Generator[T, M] {
	g.generationPartFilters = append(g.generationPartFilters, generationPartFilter)
	return g
}

func (g *Generator[T, M]) Run() error {
	preprocessArgs()

	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	// getting them early to be able to easily list available options in help
	objects := g.objectsProvider()
	parts := g.generationParts

	filterObjects := newFiltersFlag("object names", collections.Map(objects, func(o T) string { return o.ObjectName() }))
	filterParts := newFiltersFlag("generation part names", collections.Map(parts, func(p GenerationPart[T, M]) string { return p.GetName() }))

	additionalLogs := flag.Bool("verbose", false, "print additional object debug logs")
	dryRun := flag.Bool("dry-run", false, "generate to std out instead of saving")
	flag.Var(filterObjects, filterObjects.flagName(), filterObjects.usage())
	flag.Var(filterParts, filterParts.flagName(), filterParts.usage())

	// TODO [this PR]: generic description
	flag.Usage = func() {
		usage := `Generate SDK objects based on the SQL definitions provided.

usage: make generate-sdk SF_TF_GENERATOR_ARGS='<args>'
`
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", usage)
		flag.PrintDefaults()
	}

	flag.Parse()

	if filterObjects.hasValues() {
		fmt.Printf("Object filters present: %s\n", filterObjects)
		g.objectFilters = append(g.objectFilters, filterObjectByNameProvider[T](filterObjects.filters))
	}
	if filterParts.hasValues() {
		fmt.Printf("Generation part filters present: %s\n", filterParts)
		g.generationPartFilters = append(g.generationPartFilters, filterGenerationPartByNameProvider[T, M](filterParts.filters))
	}

	if len(g.objectFilters) > 0 {
		filteredObjects := make([]T, 0)
		for _, o := range objects {
			matches := true
			for _, f := range g.objectFilters {
				matches = matches && f(o)
			}
			if matches {
				filteredObjects = append(filteredObjects, o)
			}
		}
		objects = filteredObjects
	}

	if len(g.generationPartFilters) > 0 {
		filteredGenerationParts := make([]GenerationPart[T, M], 0)
		for _, p := range g.generationParts {
			matches := true
			for _, f := range g.generationPartFilters {
				matches = matches && f(p)
			}
			if matches {
				filteredGenerationParts = append(filteredGenerationParts, p)
			}
		}
		parts = filteredGenerationParts
	}

	if *additionalLogs {
		for _, p := range g.additionalObjectDebugLogProviders {
			p(objects)
		}
	}

	if *dryRun {
		if err := g.generateAndPrint(objects, parts); err != nil {
			return err
		}
	} else {
		if err := g.generateAndSave(objects, parts); err != nil {
			return err
		}
	}

	return nil
}

// TODO [SNOW-1501905]: temporary hacky solution to allow easy passing multiple args from the make command
// TODO [this PR]: describe the reasoning
func preprocessArgs() {
	rest := os.Args[1:]
	newArgs := []string{os.Args[0]}
	for _, a := range rest {
		newArgs = append(newArgs, strings.Split(a, " ")...)
	}
	os.Args = newArgs
}

func (g *Generator[_, _]) RunAndHandleOsReturn() {
	err := g.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Generator[T, M]) generateAndSave(objects []T, parts []GenerationPart[T, M]) error {
	var errs []error
	for _, s := range objects {
		model := g.modelProvider(s, g.preamble)

		for _, p := range parts {
			buffer := bytes.Buffer{}

			if err := executeAllTemplates(model, &buffer, p.templates...); err != nil {
				errs = append(errs, fmt.Errorf("generating output for object %s failed with err: %w", s.ObjectName(), err))
				continue
			}
			filename := p.filenameProvider(s, model)
			if err := WriteCodeToFile(&buffer, filename); err != nil {
				errs = append(errs, fmt.Errorf("saving output for object %s to file %s failed with err: %w", s.ObjectName(), filename, err))
				continue
			}
		}
	}
	return errors.Join(errs...)
}

func (g *Generator[T, M]) generateAndPrint(objects []T, parts []GenerationPart[T, M]) error {
	var errs []error
	for _, s := range objects {
		fmt.Println("===========================")
		fmt.Printf("Generating for object %s\n", s.ObjectName())
		fmt.Println("===========================")
		for _, p := range parts {
			if err := executeAllTemplates(g.modelProvider(s, g.preamble), os.Stdout, p.templates...); err != nil {
				errs = append(errs, fmt.Errorf("generating output for object %s failed with err: %w", s.ObjectName(), err))
				continue
			}
		}
	}
	return errors.Join(errs...)
}

func executeAllTemplates[M HasPreambleModel](model M, writer io.Writer, templates ...*template.Template) error {
	var errs []error
	for _, t := range templates {
		if err := t.Execute(writer, model); err != nil {
			errs = append(errs, fmt.Errorf("template execution for template %s failed with err: %w", t.Name(), err))
		}
	}
	return errors.Join(errs...)
}
