package defs

import (
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/stretchr/testify/require"
)

func Test_DefinitionsOmitDefaultAttributeValues(t *testing.T) {
	var findings []string
	for _, obj := range All() {
		findings = append(findings, ssot.RedundantExplicitDefaults(obj)...)
	}
	slices.Sort(findings)
	findings = slices.Compact(findings)
	require.Empty(t, findings)
}
