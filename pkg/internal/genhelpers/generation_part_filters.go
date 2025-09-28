package genhelpers

import (
	"os"
	"slices"
	"strings"
)

// TODO [SNOW-2324252]: Consider extracting this as a command line param
func FilterGenerationPartByNameFromEnv[T ObjectNameProvider, M HasPreambleModel](part GenerationPart[T, M]) bool {
	allowedObjectNamesString := os.Getenv("SF_TF_GENERATOR_EXT_ALLOWED_GENERATION_PARTS_NAMES")
	if allowedObjectNamesString == "" {
		return true
	}
	allowedObjectNames := strings.Split(allowedObjectNamesString, ",")
	return slices.Contains(allowedObjectNames, part.GetName())
}

func filterGenerationPartByNameProvider[T ObjectNameProvider, M HasPreambleModel](allowedGenerationParts []string) func(part GenerationPart[T, M]) bool {
	return func(part GenerationPart[T, M]) bool {
		if len(allowedGenerationParts) == 0 {
			return true
		}
		return slices.Contains(allowedGenerationParts, part.GetName())
	}
}
