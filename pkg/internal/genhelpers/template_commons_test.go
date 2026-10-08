package genhelpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Unimplemented(t *testing.T) {
	require.PanicsWithValue(t, "generated implementation for ParseIdFromConfig is not implemented", func() {
		Unimplemented("ParseIdFromConfig")
	})
}
