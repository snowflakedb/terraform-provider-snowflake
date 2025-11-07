package genhelpers

import (
	"fmt"

	"golang.org/x/tools/imports"
)

func AddImports(outputPath string, input []byte) ([]byte, error) {
	src, err := imports.Process(outputPath, input, &imports.Options{
		Fragment:   false,
		AllErrors:  true,
		Comments:   true,
		TabIndent:  true,
		TabWidth:   4,
		FormatOnly: false,
	})
	if err != nil {
		return nil, fmt.Errorf("adding missing imports to file %s failed with err: %w", outputPath, err)
	}
	return src, nil
}
