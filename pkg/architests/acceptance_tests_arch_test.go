package architests

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/architest"
)

func TestArchCheck_AcceptanceTests(t *testing.T) {
	testAccPackageFiles := architest.Directory("../testacc/").AllFiles()
	acceptanceTestFiles := testAccPackageFiles.Filter(architest.FileNameRegexFilterProvider(architest.AcceptanceTestFileRegex))
	otherTestFiles := testAccPackageFiles.Filter(architest.FileNameFilterWithExclusionsProvider(
		architest.TestFileRegex,
		architest.AcceptanceTestFileRegex,
	))

	t.Run("acceptance tests files have the right package", func(t *testing.T) {
		acceptanceTestFiles.All(func(file *architest.File) {
			file.AssertHasPackage(t, "testacc")
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
		otherTestFiles.All(func(file *architest.File) {
			file.ExportedMethods().All(func(method *architest.Method) {
				method.AssertNameDoesNotMatch(t, architest.AcceptanceTestNameRegex)
			})
		})
	})

	t.Run("there are no non-acceptance tests in other test files in the directory", func(t *testing.T) {
		otherTestFiles.All(func(file *architest.File) {
			file.ExportedMethods().All(func(method *architest.Method) {
				// our acceptance tests have TestMain, let's filter it out now (maybe later we can support it in architest)
				if method.Name() != "TestMain" {
					method.AssertNameDoesNotMatch(t, architest.TestNameRegex)
				}
			})
		})
	})
}
