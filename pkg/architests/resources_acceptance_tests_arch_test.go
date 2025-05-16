package architests

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/architest"
)

func TestArchCheck_AcceptanceTests_Resources(t *testing.T) {
	resourcesFiles := architest.Directory("../resources/").AllFiles()
	acceptanceTestFiles := resourcesFiles.Filter(architest.FileNameRegexFilterProvider(architest.AcceptanceTestFileRegex))

	t.Run("acceptance tests files have the right package", func(t *testing.T) {
		acceptanceTestFiles.All(func(file *architest.File) {
			file.AssertHasPackage(t, "resources_test")
		})
	})

	t.Run("acceptance tests are named correctly", func(t *testing.T) {
		acceptanceTestFiles.All(func(file *architest.File) {
			file.ExportedMethods().All(func(method *architest.Method) {
				method.AssertAcceptanceTestNamedCorrectly(t)
			})
		})
	})

	t.Run("there are no acceptance tests in other test files in the directory", func(t *testing.T) {
		otherTestFiles := resourcesFiles.Filter(architest.FileNameFilterWithExclusionsProvider(
			architest.TestFileRegex,
			architest.AcceptanceTestFileRegex,
			regexp.MustCompile("helpers_test.go"),
			regexp.MustCompile("view_test.go"),
		))

		otherTestFiles.All(func(file *architest.File) {
			file.ExportedMethods().All(func(method *architest.Method) {
				method.AssertNameDoesNotMatch(t, architest.AcceptanceTestNameRegex)
				method.AssertTestNamedCorrectly(t)
			})
		})
	})

	t.Run("there are only acceptance tests in package resources_test", func(t *testing.T) {
		t.Skipf("Currently there are non-acceptance tests in resources_test package")
		packageFiles := resourcesFiles.Filter(architest.PackageFilterProvider("resources_test"))

		packageFiles.All(func(file *architest.File) {
			file.ExportedMethods().All(func(method *architest.Method) {
				method.AssertAcceptanceTestNamedCorrectly(t)
			})
		})
	})
}
