package genhelpers

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

type namesListFlag struct {
	name             string
	prefix           string
	filters          []string
	availableOptions []string
}

func newInclusionFlag(name string, availableOptions []string) *namesListFlag {
	slices.Sort(availableOptions)
	return &namesListFlag{
		name:             name,
		prefix:           "filter",
		filters:          make([]string, 0),
		availableOptions: availableOptions,
	}
}

func newExclusionFlag(name string, availableOptions []string) *namesListFlag {
	slices.Sort(availableOptions)
	return &namesListFlag{
		name:             name,
		prefix:           "exclude",
		filters:          make([]string, 0),
		availableOptions: availableOptions,
	}
}

func newEnablementFlag(name string, availableOptions []string) *namesListFlag {
	slices.Sort(availableOptions)
	return &namesListFlag{
		name:             name,
		prefix:           "enable",
		filters:          make([]string, 0),
		availableOptions: availableOptions,
	}
}

func (f *namesListFlag) hasValues() bool {
	return len(f.filters) > 0
}

func (f *namesListFlag) flagName() string {
	return fmt.Sprintf("%s-%s", f.prefix, strings.ToLower(strings.ReplaceAll(f.name, " ", "-")))
}

func (f *namesListFlag) usage() string {
	if f.prefix == "exclude" {
		return fmt.Sprintf("exclude the given %[1]s from generation like `<name1>,<name2>,...`; available %[1]s:\n%s", f.name, formatAvailableOptions(f.availableOptions))
	}
	if f.prefix == "enable" {
		return fmt.Sprintf("enable the given optional %[1]s for all objects like `<name1>,<name2>,...`; available %[1]s:\n%s", f.name, formatAvailableOptions(f.availableOptions))
	}
	return fmt.Sprintf("generate only for the given %[1]s like `<name1>,<name2>,...`; available %[1]s:\n%s", f.name, formatAvailableOptions(f.availableOptions))
}

func formatAvailableOptions(options []string) string {
	return strings.Join(collections.Map(options, func(s string) string { return " - " + s }), "\n") + "\n"
}

func (f *namesListFlag) String() string {
	if len(f.filters) == 0 {
		// this is considered a default value; the default usage prints: (default ...)
		return fmt.Sprintf("is an empty list: all %s are used", f.name)
	}
	return fmt.Sprintf("%s %s with values %s (len: %d)", f.name, f.prefix, f.filters, len(f.filters))
}

func (f *namesListFlag) Set(value string) error {
	if len(f.filters) > 0 {
		return fmt.Errorf("%s %s has been already initiated; use single attribute with comma-separated list", f.name, f.prefix)
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s %s cannot be initiated with an empty value; use single attribute with comma-separated list", f.name, f.prefix)
	}
	for _, fil := range strings.Split(value, ",") {
		f.filters = append(f.filters, fil)
	}
	return nil
}
