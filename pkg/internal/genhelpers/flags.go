package genhelpers

import (
	"fmt"
	"strings"
)

type filtersFlag struct {
	name    string
	filters []string
}

func newFiltersFlag(name string) *filtersFlag {
	return &filtersFlag{
		name:    name,
		filters: make([]string, 0),
	}
}

func (f *filtersFlag) hasValues() bool {
	return len(f.filters) > 0
}

func (f *filtersFlag) String() string {
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
