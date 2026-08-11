package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_ToDatabaseKind_Aliases covers the aliased values. The generic conversion test (registered in Test_AllEnumConversions)
// iterates AllDatabaseKinds, which holds only the canonical values, so it does not exercise the aliases at all.
// TODO [SNOW-2324252]: Remove this test after generator improvements.
func Test_ToDatabaseKind_Aliases(t *testing.T) {
	type test struct {
		input string
		want  DatabaseKind
	}

	valid := []test{
		{input: "APPLICATION_PACKAGE", want: DatabaseKindApplicationPackage},
		{input: "application_package", want: DatabaseKindApplicationPackage},
	}

	// APPLICATION PACKAGE is the only value with an underscored alias, because it's the only one DATABASES.TYPE
	// is known to report differently than SHOW DATABASES.
	invalid := []test{
		{input: "IMPORTED_DATABASE"},
		{input: "PERSONAL_DATABASE"},
		{input: "CATALOG_LINKED_DATABASE"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToDatabaseKind(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToDatabaseKind(tc.input)
			require.Error(t, err)
		})
	}
}

// Test_ToDatabaseKind_MultiWordValuesCaseInsensitivity covers the values containing spaces and hyphens. The
// generic conversion test only lowercases the first enum value (STANDARD), which is a single word, so it
// proves nothing about the values that strings.ToUpper has to match across a space or a hyphen.
func Test_ToDatabaseKind_MultiWordValuesCaseInsensitivity(t *testing.T) {
	type test struct {
		input string
		want  DatabaseKind
	}

	valid := []test{
		{input: "imported database", want: DatabaseKindImportedDatabase},
		{input: "Application Package", want: DatabaseKindApplicationPackage},
		{input: "personal database", want: DatabaseKindPersonalDatabase},
		{input: "catalog-linked database", want: DatabaseKindCatalogLinkedDatabase},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToDatabaseKind(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
