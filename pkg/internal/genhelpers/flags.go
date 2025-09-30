package genhelpers

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

type filtersFlag struct {
	name             string
	filters          []string
	availableOptions []string
}

func newFiltersFlag(name string, availableOptions []string) *filtersFlag {
	slices.Sort(availableOptions)
	return &filtersFlag{
		name:             name,
		filters:          make([]string, 0),
		availableOptions: availableOptions,
	}
}

func (f *filtersFlag) hasValues() bool {
	return len(f.filters) > 0
}

func (f *filtersFlag) flagName() string {
	return fmt.Sprintf("filter-%s", strings.ToLower(strings.ReplaceAll(f.name, " ", "-")))
}

func (f *filtersFlag) usage() string {
	return fmt.Sprintf("generate only for the given %[1]s like `<name1>,<name2>,...`; available %[1]s:\n%s", f.name, formatAvailableOptions(f.availableOptions))
}

func formatAvailableOptions(options []string) string {
	return strings.Join(collections.Map(options, func(s string) string { return " - " + s }), "\n") + "\n"
}

func (f *filtersFlag) String() string {
	if len(f.filters) == 0 {
		// this is considered a default value; the default usage prints: (default ...)
		return fmt.Sprintf("is an empty list: all %s are used", f.name)
	}
	return fmt.Sprintf("%s filter with values %s (len: %d)", f.name, f.filters, len(f.filters))
}

func (f *filtersFlag) Set(value string) error {
	if len(f.filters) > 0 {
		return fmt.Errorf("%s filters has been already initiated; use single attribute with comma-separated list", f.name)
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s filters cannot be initiated with an empty value; use single attribute with comma-separated list", f.name)
	}
	for _, fil := range strings.Split(value, ",") {
		f.filters = append(f.filters, fil)
	}
	return nil
}
