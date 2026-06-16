package genhelpers

import (
	"slices"
)

func filterGenerationPartByNameProvider[T ObjectNameProvider, M HasPreambleModel](allowedGenerationParts []string) func(part GenerationPart[T, M]) bool {
	return func(part GenerationPart[T, M]) bool {
		if len(allowedGenerationParts) == 0 {
			return true
		}
		return slices.Contains(allowedGenerationParts, part.GetName())
	}
}

func excludeGenerationPartByNameProvider[T ObjectNameProvider, M HasPreambleModel](excludedGenerationParts []string) func(part GenerationPart[T, M]) bool {
	return func(part GenerationPart[T, M]) bool {
		if len(excludedGenerationParts) == 0 {
			return true
		}
		return !slices.Contains(excludedGenerationParts, part.GetName())
	}
}
