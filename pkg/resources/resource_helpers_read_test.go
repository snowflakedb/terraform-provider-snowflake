package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_derefOrZero(t *testing.T) {
	t.Run("nil string pointer returns empty string", func(t *testing.T) {
		var value *string
		assert.Equal(t, "", derefOrZero(value))
	})

	t.Run("string pointer returns the value", func(t *testing.T) {
		value := "WH"
		assert.Equal(t, "WH", derefOrZero(&value))
	})

	t.Run("nil int pointer returns zero", func(t *testing.T) {
		var value *int
		assert.Equal(t, 0, derefOrZero(value))
	})

	t.Run("int pointer returns the value", func(t *testing.T) {
		value := 42
		assert.Equal(t, 42, derefOrZero(&value))
	})
}
