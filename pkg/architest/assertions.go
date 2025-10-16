package architest

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (f *File) AssertHasPackage(t *testing.T, expectedPackage string) {
	t.Helper()
	assert.Equalf(t, expectedPackage, f.packageName, "file %s has package %s, expected package %s", f.Name(), f.PackageName(), expectedPackage)
}

func (method *Method) AssertAcceptanceTestNamedCorrectly(t *testing.T) {
	t.Helper()
	method.AssertNameMatches(t, AcceptanceTestNameRegex)
}

func (method *Method) AssertIntegrationTestNamedCorrectly(t *testing.T) {
	t.Helper()
	method.AssertNameMatches(t, IntegrationTestNameRegex)
}

func (method *Method) AssertTestNamedCorrectly(t *testing.T) {
	t.Helper()
	method.AssertNameMatches(t, TestNameRegex)
}

func (method *Method) AssertNameMatches(t *testing.T, regex *regexp.Regexp) {
	t.Helper()
	assert.Truef(t, regex.MatchString(method.Name()), "file %s contains exported method %s which does not match %s", method.FileName(), method.Name(), regex.String())
}

func (method *Method) AssertNameDoesNotMatch(t *testing.T, regex *regexp.Regexp) {
	t.Helper()
	assert.Falsef(t, regex.MatchString(method.Name()), "file %s contains exported method %s which matches %s", method.FileName(), method.Name(), regex.String())
}

// TODO [SNOW-2036808]: add build directive assertion on file (to check the presence of //go:build non_account_level_tests and //go:build account_level_tests directives)
