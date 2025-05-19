package helpers

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
)

// AssertErrorContainsPartsFunc returns a function asserting error message contains each string in parts
func AssertErrorContainsPartsFunc(t *testing.T, parts []string) resource.ErrorCheckFunc {
	t.Helper()
	return func(err error) error {
		for _, part := range parts {
			assert.Contains(t, err.Error(), part)
		}
		return nil
	}
}
