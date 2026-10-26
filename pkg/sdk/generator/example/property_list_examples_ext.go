package example

import (
	"errors"
	"strings"
)

// PropertyListExampleItemDetails is the element type of the prefix list in property_list_example_def.go.
type PropertyListExampleItemDetails struct {
	Name string
}

// ParsePropertyListExampleItemDescribed is the hand-written element parser of the ITEM_ prefix mapping.
func ParsePropertyListExampleItemDescribed(value string) (PropertyListExampleItemDetails, error) {
	name := strings.TrimSpace(value)
	if name == "" {
		return PropertyListExampleItemDetails{}, errors.New("an item value is empty")
	}
	return PropertyListExampleItemDetails{Name: name}, nil
}

// validatePropertyListExampleDescribed is the hand-written check that the generated parse runs last.
func validatePropertyListExampleDescribed(details PropertyListExampleDetails) error {
	if len(details.Items) == 0 {
		return errors.New("no items could be parsed")
	}
	return nil
}
