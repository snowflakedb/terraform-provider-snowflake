package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DerefOrZero(t *testing.T) {
	t.Run("nil string pointer returns empty string", func(t *testing.T) {
		var value *string
		assert.Equal(t, "", DerefOrZero(value))
	})

	t.Run("string pointer returns the value", func(t *testing.T) {
		value := "WH"
		assert.Equal(t, "WH", DerefOrZero(&value))
	})

	t.Run("nil int pointer returns zero", func(t *testing.T) {
		var value *int
		assert.Equal(t, 0, DerefOrZero(value))
	})

	t.Run("int pointer returns the value", func(t *testing.T) {
		value := 42
		assert.Equal(t, 42, DerefOrZero(&value))
	})
}
