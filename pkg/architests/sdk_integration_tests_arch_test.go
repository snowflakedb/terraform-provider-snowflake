package architests

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/architest"
)

func TestArchCheck_IntegrationTests_Sdk(t *testing.T) {
	sdkIntegrationTestFiles := architest.Directory("../sdk/testint/").AllFiles()
	integrationTestFiles := sdkIntegrationTestFiles.Filter(architest.FileNameRegexFilterProvider(architest.IntegrationTestFileRegex))

	t.Run("integration tests files have the right package", func(t *testing.T) {
		integrationTestFiles.All(func(file *architest.File) {
			file.AssertHasPackage(t, "testint")
		})
	})

	t.Run("integration tests are named correctly", func(t *testing.T) {
		integrationTestFiles.All(func(file *architest.File) {
			file.ExportedMethods().All(func(method *architest.Method) {
				method.AssertIntegrationTestNamedCorrectly(t)
			})
		})
	})

	t.Run("there are no integration tests in other test files in the directory", func(t *testing.T) {
		otherTestFiles := sdkIntegrationTestFiles.Filter(architest.FileNameFilterWithExclusionsProvider(architest.TestFileRegex, architest.IntegrationTestFileRegex))

		otherTestFiles.All(func(file *architest.File) {
			file.ExportedMethods().All(func(method *architest.Method) {
				method.AssertNameDoesNotMatch(t, architest.IntegrationTestNameRegex)
				method.AssertTestNamedCorrectly(t)
			})
		})
	})

	t.Run("there are only integration tests in package testint", func(t *testing.T) {
		packageFiles := sdkIntegrationTestFiles.Filter(architest.PackageFilterProvider("testint"))

		packageFiles.All(func(file *architest.File) {
			file.ExportedMethods().All(func(method *architest.Method) {
				// our integration tests have TestMain, let's filter it out now (maybe later we can support in in architest)
				if method.Name() != "TestMain" {
					method.AssertIntegrationTestNamedCorrectly(t)
				}
			})
		})
	})
}
