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

	additionalObjectDebugLogProviders []func([]T)
	objectFilters                     []func(T) bool

	preamble *PreambleModel
}

func NewGenerator[T ObjectNameProvider, M HasPreambleModel](preamble *PreambleModel, objectsProvider func() []T, modelProvider func(T, *PreambleModel) M, filenameProvider func(T, M) string, templates []*template.Template) *Generator[T, M] {
	return &Generator[T, M]{
		objectsProvider:  objectsProvider,
		modelProvider:    modelProvider,
		filenameProvider: filenameProvider,
		templates:        templates,

		additionalObjectDebugLogProviders: make([]func([]T), 0),
		objectFilters:                     make([]func(T) bool, 0),

		// TODO: format version in the preamble differently
		preamble: preamble,
	}
}

func (g *Generator[T, M]) WithAdditionalObjectsDebugLogs(objectLogsProvider func([]T)) *Generator[T, M] {
	g.additionalObjectDebugLogProviders = append(g.additionalObjectDebugLogProviders, objectLogsProvider)
	return g
}

func (g *Generator[T, M]) WithObjectFilter(objectFilter func(T) bool) *Generator[T, M] {
	g.objectFilters = append(g.objectFilters, objectFilter)
	return g
}

func (g *Generator[T, _]) Run() error {
	preprocessArgs()

	file := os.Getenv("GOFILE")
	fmt.Printf("Running generator on %s with args %#v\n", file, os.Args[1:])

	additionalLogs := flag.Bool("verbose", false, "print additional object debug logs")
	dryRun := flag.Bool("dry-run", false, "generate to std out instead of saving")
	flag.Parse()

	objects := g.objectsProvider()

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

	if *additionalLogs {
		for _, p := range g.additionalObjectDebugLogProviders {
			p(objects)
		}
	}

	if *dryRun {
		if err := g.generateAndPrintForAllObjects(objects); err != nil {
			return err
		}
	} else {
		if err := g.generateAndSaveForAllObjects(objects); err != nil {
			return err
		}
	}

	return nil
}

// TODO [SNOW-1501905]: temporary hacky solution to allow easy passing multiple args from the make command
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

func (g *Generator[T, _]) generateAndSaveForAllObjects(objects []T) error {
	var errs []error
	for _, s := range objects {
		buffer := bytes.Buffer{}
		model := g.modelProvider(s, g.preamble)
		if err := executeAllTemplates(model, &buffer, g.templates...); err != nil {
			errs = append(errs, fmt.Errorf("generating output for object %s failed with err: %w", s.ObjectName(), err))
			continue
		}
		filename := g.filenameProvider(s, model)
		if err := WriteCodeToFile(&buffer, filename); err != nil {
			errs = append(errs, fmt.Errorf("saving output for object %s to file %s failed with err: %w", s.ObjectName(), filename, err))
			continue
		}
	}
	return errors.Join(errs...)
}

func (g *Generator[T, _]) generateAndPrintForAllObjects(objects []T) error {
	var errs []error
	for _, s := range objects {
		fmt.Println("===========================")
		fmt.Printf("Generating for object %s\n", s.ObjectName())
		fmt.Println("===========================")
		if err := executeAllTemplates(g.modelProvider(s, g.preamble), os.Stdout, g.templates...); err != nil {
			errs = append(errs, fmt.Errorf("generating output for object %s failed with err: %w", s.ObjectName(), err))
			continue
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
