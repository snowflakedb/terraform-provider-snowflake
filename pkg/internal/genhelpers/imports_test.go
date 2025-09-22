package genhelpers_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/imports"
)

// TODO [next PR - SNOW-2324252]: extract common setup, template for easier test writing
func Test_ExperimentWithImportsProcess(t *testing.T) {
	t.Run("add standard imports", func(t *testing.T) {
		src := []byte(`package somepackagename

func hello() {
	fmt.Println(time.Now())
}
`)
		expected := `package somepackagename

import (
	"fmt"
	"time"
)

func hello() {
	fmt.Println(time.Now())
}
`

		out, err := imports.Process("", src, nil)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})

	t.Run("add sdk import", func(t *testing.T) {
		src := []byte(`package somepackagename

func hello() {
	id := sdk.NewAccountObjectIdentifier("test")
	fmt.Println(id.FullyQualifiedName())
}
`)
		expected := `package somepackagename

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func hello() {
	id := sdk.NewAccountObjectIdentifier("test")
	fmt.Println(id.FullyQualifiedName())
}
`

		out, err := imports.Process("", src, nil)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})

	t.Run("add ambiguous import", func(t *testing.T) {
		src := []byte(`package somepackagename

func hello() {
	r := resources.Database
	fmt.Println(r)
}
`)
		// the second one possible is "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
		expected := `package somepackagename

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
)

func hello() {
	r := resources.Database
	fmt.Println(r)
}
`

		out, err := imports.Process("", src, nil)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})

	t.Run("add ambiguous import - with explicit import", func(t *testing.T) {
		src := []byte(`package somepackagename

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

func hello() {
	r := resources.Database
	fmt.Println(r)
}
`)
		expected := `package somepackagename

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

func hello() {
	r := resources.Database
	fmt.Println(r)
}
`

		out, err := imports.Process("", src, nil)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})

	t.Run("add unused explicit import", func(t *testing.T) {
		src := []byte(`package somepackagename

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

func hello() {
	fmt.Println("not used")
}
`)
		expected := `package somepackagename

import "fmt"

func hello() {
	fmt.Println("not used")
}
`

		out, err := imports.Process("", src, nil)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})

	t.Run("add explicit named import", func(t *testing.T) {
		src := []byte(`package somepackagename

import re "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

func hello() {
	r := re.Database
	fmt.Println(r)
}
`)
		expected := `package somepackagename

import (
	"fmt"

	re "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

func hello() {
	r := re.Database
	fmt.Println(r)
}
`

		out, err := imports.Process("", src, nil)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})
}

func Test_AddImports(t *testing.T) {
	t.Run("add standard imports", func(t *testing.T) {
		src := []byte(`package somepackagename

func hello() {
	fmt.Println(time.Now())
}
`)
		expected := `package somepackagename

import (
	"fmt"
	"time"
)

func hello() {
	fmt.Println(time.Now())
}
`

		out, err := genhelpers.AddImports("", src)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})

	t.Run("add sdk import", func(t *testing.T) {
		src := []byte(`package somepackagename

func hello() {
	id := sdk.NewAccountObjectIdentifier("test")
	fmt.Println(id.FullyQualifiedName())
}
`)
		expected := `package somepackagename

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func hello() {
	id := sdk.NewAccountObjectIdentifier("test")
	fmt.Println(id.FullyQualifiedName())
}
`

		out, err := genhelpers.AddImports("", src)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})

	t.Run("add ambiguous import", func(t *testing.T) {
		src := []byte(`package somepackagename

func hello() {
	r := resources.Database
	fmt.Println(r)
}
`)
		// the second one possible is "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
		expected := `package somepackagename

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
)

func hello() {
	r := resources.Database
	fmt.Println(r)
}
`

		out, err := genhelpers.AddImports("", src)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})

	t.Run("add ambiguous import - with explicit import", func(t *testing.T) {
		src := []byte(`package somepackagename

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

func hello() {
	r := resources.Database
	fmt.Println(r)
}
`)
		expected := `package somepackagename

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

func hello() {
	r := resources.Database
	fmt.Println(r)
}
`

		out, err := genhelpers.AddImports("", src)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})

	t.Run("add unused explicit import", func(t *testing.T) {
		src := []byte(`package somepackagename

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

func hello() {
	fmt.Println("not used")
}
`)
		expected := `package somepackagename

import "fmt"

func hello() {
	fmt.Println("not used")
}
`

		out, err := genhelpers.AddImports("", src)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})

	t.Run("add explicit named import", func(t *testing.T) {
		src := []byte(`package somepackagename

import re "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

func hello() {
	r := re.Database
	fmt.Println(r)
}
`)
		expected := `package somepackagename

import (
	"fmt"

	re "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

func hello() {
	r := re.Database
	fmt.Println(r)
}
`

		out, err := genhelpers.AddImports("", src)

		require.NoError(t, err)
		require.Equal(t, expected, string(out))
	})
}
