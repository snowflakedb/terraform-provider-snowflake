package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type TagPropagation string

const (
	TagPropagationNone                        TagPropagation = "NONE"
	TagPropagationOnDependency                TagPropagation = "ON_DEPENDENCY"
	TagPropagationOnDataMovement              TagPropagation = "ON_DATA_MOVEMENT"
	TagPropagationOnDependencyAndDataMovement TagPropagation = "ON_DEPENDENCY_AND_DATA_MOVEMENT"
)

var AllTagPropagationValues = []TagPropagation{
	TagPropagationNone,
	TagPropagationOnDependency,
	TagPropagationOnDataMovement,
	TagPropagationOnDependencyAndDataMovement,
}

func ToTagPropagation(s string) (TagPropagation, error) {
	tp := TagPropagation(strings.ToUpper(s))
	if !slices.Contains(AllTagPropagationValues, tp) {
		return "", fmt.Errorf("invalid tag propagation value: %s", tp)
	}
	return tp, nil
}
