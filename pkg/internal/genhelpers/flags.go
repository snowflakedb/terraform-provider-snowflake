package genhelpers

import (
	"errors"
	"fmt"
	"strings"
)

// TODO [this PR]: add filter name to this type
type filters []string

func (f *filters) String() string {
	return fmt.Sprint(*f)
}

func (f *filters) Set(value string) error {
	if len(*f) > 0 {
		return errors.New("filters already")
	}
	for _, fil := range strings.Split(value, ",") {
		*f = append(*f, fil)
	}
	return nil
}
