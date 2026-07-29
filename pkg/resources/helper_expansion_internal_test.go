package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpandStringList(t *testing.T) {
	r := require.New(t)

	in := []interface{}{"this", "is", "just", "a", "test"}
	out := expandStringList(in)

	r.Equal("this", out[0])
	r.Equal("is", out[1])
	r.Equal("just", out[2])
	r.Equal("a", out[3])
	r.Equal("test", out[4])
}

func TestExpandBlankStringList(t *testing.T) {
	r := require.New(t)
	in := []interface{}{}
	out := expandStringList(in)

	r.Empty(out)
}
