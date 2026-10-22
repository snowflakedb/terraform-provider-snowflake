package genhelpers

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
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

	description         string
	makefileCommandPart string
	cliEnabledParts     []string

	preamble *PreambleModel
}

type GenerationPart[T ObjectNameProvider, M HasPreambleModel] struct {
	name             string
	filenameProvider func(T, M) string
	templates        []*template.Template
	condition        func(T) bool
	enabledByDefault bool
}

func (g *GenerationPart[_, _]) GetName() string {
	return g.name
}

func (g *GenerationPart[_, _]) IsEnabledByDefault() bool {
	return g.enabledByDefault
}

func NewGenerationPart[T ObjectNameProvider, M HasPreambleModel](name GenerationPartNamer, filenameProvider func(T, M) string, templates []*template.Template) GenerationPart[T, M] {
	return GenerationPart[T, M]{
		name:             name.GenerationPartName(),
		filenameProvider: filenameProvider,
		templates:        templates,
		enabledByDefault: true,
	}
}

func NewGenerator[T ObjectNameProvider, M HasPreambleModel](preamble *PreambleModel, objectsProvider func() []T, modelProvider func(T, *PreambleModel) M, filenameProvider func(T, M) string, templates []*template.Template) *Generator[T, M] {
	// TODO [SNOW-2324252]: handle vararg input
	parts := []GenerationPart[T, M]{
		// TODO [SNOW-2324252]: change default to name when changing to vararg
		NewGenerationPart(DefaultGenerationPart, filenameProvider, templates),
	}
	// TODO [SNOW-2324252]: Probably remove later; it should be a part of the constructor
	defaultDescription := "Generator's description missing."
	defaultMakefileCommandPart := "<makefile-command-part>"
	return &Generator[T, M]{
		objectsProvider:  objectsProvider,
		modelProvider:    modelProvider,
		filenameProvider: filenameProvider,
		templates:        templates,

		generationParts: parts,

		additionalObjectDebugLogProviders: make([]func([]T), 0),
		objectFilters:                     make([]func(T) bool, 0),
		generationPartFilters:             make([]func(GenerationPart[T, M]) bool, 0),

		description:         defaultDescription,
		makefileCommandPart: defaultMakefileCommandPart,

		preamble: preamble,
	}
}

// TODO [SNOW-2324252]: Probably remove later when we have vararg support in the NewGenerator constructor
func (g *Generator[T, M]) WithGenerationPart(partName GenerationPartNamer, filenameProvider func(T, M) string, templates []*template.Template) *Generator[T, M] {
	g.generationParts = append(g.generationParts, NewGenerationPart(partName, filenameProvider, templates))
	return g
}

// WithConditionalGenerationPart registers a generation part that is only executed for an object when condition returns true.
func (g *Generator[T, M]) WithConditionalGenerationPart(partName GenerationPartNamer, filenameProvider func(T, M) string, templates []*template.Template, condition func(T) bool) *Generator[T, M] {
	part := NewGenerationPart(partName, filenameProvider, templates)
	part.condition = condition
	g.generationParts = append(g.generationParts, part)
	return g
}

// WithOptionalGenerationPart registers a generation part that is disabled by default and must be explicitly enabled per object.
func (g *Generator[T, M]) WithOptionalGenerationPart(partName GenerationPartNamer, filenameProvider func(T, M) string, templates []*template.Template) *Generator[T, M] {
	part := NewGenerationPart(partName, filenameProvider, templates)
	part.enabledByDefault = false
	g.generationParts = append(g.generationParts, part)
	return g
}

// WithOptionalConditionalGenerationPart registers a generation part that is disabled by default and has an additional runtime condition.
func (g *Generator[T, M]) WithOptionalConditionalGenerationPart(partName GenerationPartNamer, filenameProvider func(T, M) string, templates []*template.Template, condition func(T) bool) *Generator[T, M] {
	part := NewGenerationPart(partName, filenameProvider, templates)
	part.enabledByDefault = false
	part.condition = condition
	g.generationParts = append(g.generationParts, part)
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

// TODO [SNOW-2324252]: Probably remove later; it should be a part of the constructor
func (g *Generator[T, M]) WithDescription(description string) *Generator[T, M] {
	g.description = description
	return g
}

func (g *Generator[T, M]) WithMakefileCommandPart(part string) *Generator[T, M] {
	g.makefileCommandPart = part
	return g
}

// effectivePartsForObject returns the generation parts that should be used for the given object.
// It applies two layers of filtering:
// 1. Restrict to allowed parts if the object defines AllowedGenerationParts.
// 2. Filter out optional (disabled-by-default) parts unless explicitly enabled for the object or via CLI.
func (g *Generator[T, M]) effectivePartsForObject(object T, globalParts []GenerationPart[T, M]) []GenerationPart[T, M] {
	parts := globalParts

	// Restrict to allowed parts if configured on the object
	if settings, ok := any(object).(HasObjectGenerationSettings); ok {
		if s := settings.getObjectGenerationSettings(); s != nil && len(s.AllowedGenerationParts) > 0 {
			filtered := make([]GenerationPart[T, M], 0)
			for _, p := range parts {
				if slices.ContainsFunc(s.AllowedGenerationParts, func(n GenerationPartNamer) bool { return n.GenerationPartName() == p.name }) {
					filtered = append(filtered, p)
				}
			}
			parts = filtered
		}
	}

	// Filter out optional parts not explicitly enabled for this object or via CLI
	var enabledParts []GenerationPartNamer
	if settings, ok := any(object).(HasObjectGenerationSettings); ok {
		if s := settings.getObjectGenerationSettings(); s != nil {
			enabledParts = s.EnabledGenerationParts
		}
	}
	result := make([]GenerationPart[T, M], 0, len(parts))
	for _, p := range parts {
		if p.enabledByDefault {
			result = append(result, p)
		} else if slices.ContainsFunc(enabledParts, func(n GenerationPartNamer) bool { return n.GenerationPartName() == p.name }) || slices.Contains(g.cliEnabledParts, p.name) {
			result = append(result, p)
		}
	}
	return result
}

func (g *Generator[T, M]) Run() error {
	preprocessArgs()

	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	// getting them early to be able to easily list available options in help
	objects := g.objectsProvider()
	parts := g.generationParts

	filterObjects := newInclusionFlag("object names", collections.Map(objects, func(o T) string { return o.ObjectName() }))
	filterParts := newInclusionFlag("generation part names", collections.Map(parts, func(p GenerationPart[T, M]) string { return p.GetName() }))
	excludeObjects := newExclusionFlag("object names", collections.Map(objects, func(o T) string { return o.ObjectName() }))
	excludeParts := newExclusionFlag("generation part names", collections.Map(parts, func(p GenerationPart[T, M]) string { return p.GetName() }))
	optionalPartNames := collections.Map(
		collections.Filter(parts, func(p GenerationPart[T, M]) bool { return !p.enabledByDefault }),
		func(p GenerationPart[T, M]) string { return p.GetName() },
	)
	enableParts := newEnablementFlag("generation part names", optionalPartNames)

	additionalLogs := flag.Bool("verbose", false, "print additional object debug logs")
	dryRun := flag.Bool("dry-run", false, "generate to std out instead of saving")
	flag.Var(filterObjects, filterObjects.flagName(), filterObjects.usage())
	flag.Var(filterParts, filterParts.flagName(), filterParts.usage())
	flag.Var(excludeObjects, excludeObjects.flagName(), excludeObjects.usage())
	flag.Var(excludeParts, excludeParts.flagName(), excludeParts.usage())
	flag.Var(enableParts, enableParts.flagName(), enableParts.usage())

	flag.Usage = func() {
		usage := fmt.Sprintf(`%[1]s

usage: make [clean-%[2]s] generate-%[2]s SF_TF_GENERATOR_ARGS='<args>'
`, g.description, g.makefileCommandPart)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", usage)
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(os.Stderr, "\nGeneration parts enabled by default:\n")
		for _, p := range parts {
			if p.enabledByDefault {
				_, _ = fmt.Fprintf(os.Stderr, "  - %s\n", p.name)
			}
		}
		if len(optionalPartNames) > 0 {
			_, _ = fmt.Fprintf(os.Stderr, "\nGeneration parts disabled by default (opt-in per object):\n")
			for _, p := range parts {
				if !p.enabledByDefault {
					_, _ = fmt.Fprintf(os.Stderr, "  - %s\n", p.name)
				}
			}
		}
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
	if excludeObjects.hasValues() {
		fmt.Printf("Object exclusions present: %s\n", excludeObjects)
		g.objectFilters = append(g.objectFilters, excludeObjectByNameProvider[T](excludeObjects.filters))
	}
	if excludeParts.hasValues() {
		fmt.Printf("Generation part exclusions present: %s\n", excludeParts)
		g.generationPartFilters = append(g.generationPartFilters, excludeGenerationPartByNameProvider[T, M](excludeParts.filters))
	}
	if enableParts.hasValues() {
		fmt.Printf("CLI-enabled generation parts: %s\n", enableParts)
		g.cliEnabledParts = enableParts.filters
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
		if len(filteredObjects) == 0 {
			return fmt.Errorf("no objects found for the given filters: %s; exclusions: %s", filterObjects, excludeObjects)
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
		if len(filteredGenerationParts) == 0 {
			return fmt.Errorf("no generation parts found for the given filters: %s; exclusions: %s", filterParts, excludeParts)
		}
		parts = filteredGenerationParts
	}

	// Validate per-object generation settings
	for _, o := range objects {
		if settings, ok := any(o).(HasObjectGenerationSettings); ok {
			if s := settings.getObjectGenerationSettings(); s != nil {
				for _, name := range s.AllowedGenerationParts {
					if !slices.ContainsFunc(parts, func(p GenerationPart[T, M]) bool { return p.name == name.GenerationPartName() }) {
						return fmt.Errorf("object %s: allowed generation part %q does not exist", o.ObjectName(), name.GenerationPartName())
					}
				}
				for _, name := range s.EnabledGenerationParts {
					idx := slices.IndexFunc(parts, func(p GenerationPart[T, M]) bool { return p.name == name.GenerationPartName() })
					if idx == -1 {
						return fmt.Errorf("object %s: enabled generation part %q does not exist", o.ObjectName(), name.GenerationPartName())
					}
					if parts[idx].enabledByDefault {
						return fmt.Errorf("object %s: enabled generation part %q is already enabled by default", o.ObjectName(), name.GenerationPartName())
					}
				}
			}
		}
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

// We would like to be able to alter the generator behavior based on the command line flags.
// The easiest way to do this is to use a dedicated environment variable and pass it to every generator invocation, e.g.:
//
//	//go:generate go run ./gen/main/main.go $SF_TF_GENERATOR_ARGS
//
// The go:generate directive does only a simple string replacement without retokenization, so:
//
//	//go:generate go run .cmd/mygen ${MYFLAGS}
//
// is tokenized to:
//
//	[go, run, .cmd/mygen, ${MYFLAGS}]
//
// Let's say MYFLAGS="-a=42 b=somevalue" after replacement we'll get:
//
//	[go, run, .cmd/mygen, -a=42 b=somevalue].
//
// Because of that, we do the retokenization ourselves in this method.
//
// One of the potential workarounds is to wrap the invocation in sh -c '...' but for compatibility issues we will stick with the direct go run for now.
//
// References:
// - https://pkg.go.dev/cmd/go/internal/generate
func preprocessArgs() {
	newArgs := make([]string, 0, len(os.Args))
	newArgs = append(newArgs, os.Args[0])
	rest := os.Args[1:]
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

		for _, p := range g.effectivePartsForObject(s, parts) {
			if p.condition != nil && !p.condition(s) {
				log.Printf("[DEBUG] Condition for generation part %s in object %s not satisfied, skipping", p.name, s.ObjectName())
				continue
			}
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
		for _, p := range g.effectivePartsForObject(s, parts) {
			if p.condition != nil && !p.condition(s) {
				log.Printf("[DEBUG] Condition for generation part %s in object %s not satisfied, skipping", p.name, s.ObjectName())
				continue
			}
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
