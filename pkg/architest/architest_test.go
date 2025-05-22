package architest_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/architest"
	"github.com/stretchr/testify/assert"
)

func Test_Directory(t *testing.T) {
	t.Run("fail to use non-existing directory", func(t *testing.T) {
		assert.Panics(t, func() {
			architest.Directory("testdata/non_existing")
		})
	})

	t.Run("use directory", func(t *testing.T) {
		assert.NotPanics(t, func() {
			architest.Directory("testdata/dir1")
		})
	})

	tests1 := []struct {
		directory            string
		expectedFileNames    []string
		expectedPackageNames []string
	}{
		{directory: "testdata/dir1", expectedFileNames: []string{"testdata/dir1/sample1.go", "testdata/dir1/sample2.go", "testdata/dir1/different1.go"}, expectedPackageNames: []string{"dir1"}},
		{directory: "testdata/dir2", expectedFileNames: []string{"testdata/dir2/sample1.go", "testdata/dir2/sample1_test.go"}, expectedPackageNames: []string{"dir2", "dir2_test"}},
		{directory: "testdata/dir3", expectedFileNames: []string{"testdata/dir3/sample1.go", "testdata/dir3/sample1_acceptance_test.go"}, expectedPackageNames: []string{"dir3", "dir3_test"}},
		{directory: "testdata/dir4", expectedFileNames: []string{"testdata/dir4/sample1.go", "testdata/dir4/sample1_integration_test.go"}, expectedPackageNames: []string{"dir4", "dir4_test"}},
	}
	for _, tt := range tests1 {
		t.Run("list all files in the given directory", func(t *testing.T) {
			dir := architest.Directory(tt.directory)

			allFiles := dir.AllFiles()
			assert.Len(t, allFiles, len(tt.expectedFileNames))

			fileNames := make([]string, 0, len(allFiles))
			for _, f := range allFiles {
				fileNames = append(fileNames, f.Name())
			}
			assert.ElementsMatch(t, fileNames, tt.expectedFileNames)

			packageNames := make(map[string]bool)
			for _, f := range allFiles {
				packageNames[f.PackageName()] = true
			}
			assert.Len(t, packageNames, len(tt.expectedPackageNames))
			for _, name := range tt.expectedPackageNames {
				assert.Contains(t, packageNames, name)
			}
		})
	}

	tests2 := []struct {
		directory         string
		filter            architest.FileFilter
		expectedFileNames []string
	}{
		{directory: "testdata/dir1", filter: architest.FileNameFilterProvider("sample"), expectedFileNames: []string{"testdata/dir1/sample1.go", "testdata/dir1/sample2.go"}},
		{directory: "testdata/dir1", filter: architest.FileNameRegexFilterProvider(regexp.MustCompile("sample")), expectedFileNames: []string{"testdata/dir1/sample1.go", "testdata/dir1/sample2.go"}},
		{directory: "testdata/dir1", filter: architest.FileNameFilterWithExclusionsProvider(regexp.MustCompile("sample"), regexp.MustCompile("sample1")), expectedFileNames: []string{"testdata/dir1/sample2.go"}},
		{directory: "testdata/dir2", filter: architest.PackageFilterProvider("dir2"), expectedFileNames: []string{"testdata/dir2/sample1.go"}},
		{directory: "testdata/dir2", filter: architest.PackageFilterProvider("dir2_test"), expectedFileNames: []string{"testdata/dir2/sample1_test.go"}},
		{directory: "testdata/dir2", filter: architest.FileNameRegexFilterProvider(architest.AcceptanceTestFileRegex), expectedFileNames: []string{}},
		{directory: "testdata/dir3", filter: architest.FileNameRegexFilterProvider(architest.AcceptanceTestFileRegex), expectedFileNames: []string{"testdata/dir3/sample1_acceptance_test.go"}},
		{directory: "testdata/dir4", filter: architest.FileNameRegexFilterProvider(architest.AcceptanceTestFileRegex), expectedFileNames: []string{}},
		{directory: "testdata/dir2", filter: architest.FileNameRegexFilterProvider(architest.IntegrationTestFileRegex), expectedFileNames: []string{}},
		{directory: "testdata/dir3", filter: architest.FileNameRegexFilterProvider(architest.IntegrationTestFileRegex), expectedFileNames: []string{}},
		{directory: "testdata/dir4", filter: architest.FileNameRegexFilterProvider(architest.IntegrationTestFileRegex), expectedFileNames: []string{"testdata/dir4/sample1_integration_test.go"}},
		{directory: "testdata/dir2", filter: architest.FileNameRegexFilterProvider(architest.TestFileRegex), expectedFileNames: []string{"testdata/dir2/sample1_test.go"}},
		{directory: "testdata/dir3", filter: architest.FileNameRegexFilterProvider(architest.TestFileRegex), expectedFileNames: []string{"testdata/dir3/sample1_acceptance_test.go"}},
		{directory: "testdata/dir4", filter: architest.FileNameRegexFilterProvider(architest.TestFileRegex), expectedFileNames: []string{"testdata/dir4/sample1_integration_test.go"}},
	}
	for _, tt := range tests2 {
		t.Run("list only files matching filter in the given directory", func(t *testing.T) {
			dir := architest.Directory(tt.directory)

			filteredFiles := dir.Files(tt.filter)
			assert.Len(t, filteredFiles, len(tt.expectedFileNames))

			fileNames := make([]string, 0, len(filteredFiles))
			for _, f := range filteredFiles {
				fileNames = append(fileNames, f.Name())
			}
			assert.ElementsMatch(t, fileNames, tt.expectedFileNames)

			// now exactly the same but indirectly
			filteredFiles = dir.AllFiles().Filter(tt.filter)
			assert.Len(t, filteredFiles, len(tt.expectedFileNames))

			fileNames = make([]string, 0, len(filteredFiles))
			for _, f := range filteredFiles {
				fileNames = append(fileNames, f.Name())
			}
			assert.ElementsMatch(t, fileNames, tt.expectedFileNames)
		})
	}
}

func Test_Files(t *testing.T) {
	t.Run("fail to use non-existing file", func(t *testing.T) {
		assert.Panics(t, func() {
			architest.NewFileFromPath("testdata/dir1/non_existing.go")
		})
	})

	t.Run("use file", func(t *testing.T) {
		assert.NotPanics(t, func() {
			architest.NewFileFromPath("testdata/dir1/sample1.go")
		})
	})

	tests1 := []struct {
		filePath            string
		expectedMethodNames []string
	}{
		{filePath: "testdata/dir1/sample1.go", expectedMethodNames: []string{}},
		// object methods are skipped
		{filePath: "testdata/dir1/sample2.go", expectedMethodNames: []string{"A"}},
	}
	for _, tt := range tests1 {
		t.Run("list all methods in file", func(t *testing.T) {
			file := architest.NewFileFromPath(tt.filePath)

			exportedMethods := file.ExportedMethods()
			assert.Len(t, exportedMethods, len(tt.expectedMethodNames))

			methodNames := make([]string, 0, len(exportedMethods))
			for _, m := range exportedMethods {
				methodNames = append(methodNames, m.Name())
			}
			assert.ElementsMatch(t, methodNames, tt.expectedMethodNames)
		})
	}

	tests2 := []struct {
		fileNames []string
	}{
		{fileNames: []string{}},
		{fileNames: []string{"a"}},
		{fileNames: []string{"a", "A"}},
		{fileNames: []string{"A", "a"}},
		{fileNames: []string{"A", "B", "C"}},
	}
	for _, tt := range tests2 {
		t.Run("receiver invoked for every file", func(t *testing.T) {
			files := make(architest.Files, 0, len(tt.fileNames))
			for _, f := range tt.fileNames {
				files = append(files, *architest.NewFile("package", f, nil))
			}
			invokedFiles := make([]string, 0)
			receiver := func(f *architest.File) {
				invokedFiles = append(invokedFiles, f.Name())
			}

			files.All(receiver)

			assert.ElementsMatch(t, tt.fileNames, invokedFiles)
		})
	}
}

func Test_Methods(t *testing.T) {
	tests := []struct {
		methodNames []string
	}{
		{methodNames: []string{}},
		{methodNames: []string{"a"}},
		{methodNames: []string{"a", "A"}},
		{methodNames: []string{"A", "a"}},
		{methodNames: []string{"A", "B", "C"}},
	}
	for _, tt := range tests {
		t.Run("receiver invoked for every method", func(t *testing.T) {
			methods := make(architest.Methods, 0, len(tt.methodNames))
			for _, m := range tt.methodNames {
				methods = append(methods, *architest.NewMethod(m, nil))
			}
			invokedMethods := make([]string, 0)
			receiver := func(m *architest.Method) {
				invokedMethods = append(invokedMethods, m.Name())
			}

			methods.All(receiver)

			assert.ElementsMatch(t, tt.methodNames, invokedMethods)
		})
	}
}

func Test_Assertions(t *testing.T) {
	tests1 := []struct {
		filePath        string
		expectedPackage string
	}{
		{filePath: "testdata/dir1/sample1.go", expectedPackage: "dir1"},
		{filePath: "testdata/dir2/sample1.go", expectedPackage: "dir2"},
		{filePath: "testdata/dir2/sample1_test.go", expectedPackage: "dir2_test"},
	}
	for _, tt := range tests1 {
		t.Run("file package assertions", func(t *testing.T) {
			file := architest.NewFileFromPath(tt.filePath)
			tut1 := &testing.T{}
			tut2 := &testing.T{}

			file.AssertHasPackage(tut1, tt.expectedPackage)
			file.AssertHasPackage(tut2, "some_other_package")

			assert.Equal(t, false, tut1.Failed())
			assert.Equal(t, true, tut2.Failed())
		})
	}

	tests2 := []struct {
		methodName string
		correct    bool
	}{
		{methodName: "TestAcc_abc", correct: true},
		{methodName: "TestAcc_TestAcc_Test", correct: true},
		{methodName: "TestAcc_", correct: false},
		{methodName: "ATestAcc_", correct: false},
		{methodName: "TestAcc", correct: false},
		{methodName: "Test_", correct: false},
		{methodName: "Test", correct: false},
		{methodName: "Test_asdf", correct: false},
		{methodName: "TestAccc_", correct: false},
		{methodName: "TestInt_Abc", correct: false},
	}
	for _, tt := range tests2 {
		t.Run(fmt.Sprintf("acceptance test name assertions for method %s", tt.methodName), func(t *testing.T) {
			file := architest.NewFileFromPath("testdata/dir1/sample1.go")
			method := architest.NewMethod(tt.methodName, file)
			tut := &testing.T{}

			method.AssertAcceptanceTestNamedCorrectly(tut)

			assert.Equal(t, !tt.correct, tut.Failed())
		})
	}

	tests3 := []struct {
		methodName string
		regexRaw   string
		correct    bool
	}{
		{methodName: "sample1", regexRaw: "sample", correct: true},
		{methodName: "Sample1", regexRaw: "sample", correct: false},
		{methodName: "Sample1", regexRaw: "Sample", correct: true},
	}
	for _, tt := range tests3 {
		t.Run(fmt.Sprintf("matching and not matching method name assertions for %s", tt.methodName), func(t *testing.T) {
			file := architest.NewFileFromPath("testdata/dir1/sample1.go")
			method := architest.NewMethod(tt.methodName, file)
			tut1 := &testing.T{}
			tut2 := &testing.T{}

			method.AssertNameMatches(tut1, regexp.MustCompile(tt.regexRaw))
			method.AssertNameDoesNotMatch(tut2, regexp.MustCompile(tt.regexRaw))

			assert.Equal(t, !tt.correct, tut1.Failed())
			assert.Equal(t, tt.correct, tut2.Failed())
		})
	}

	tests4 := []struct {
		methodName string
		correct    bool
	}{
		{methodName: "Test", correct: true},
		{methodName: "aTest", correct: false},
		{methodName: "Test_", correct: true},
		{methodName: "Test_adsfadf", correct: true},
	}
	for _, tt := range tests4 {
		t.Run(fmt.Sprintf("test name assertions for method %s", tt.methodName), func(t *testing.T) {
			file := architest.NewFileFromPath("testdata/dir1/sample1.go")
			method := architest.NewMethod(tt.methodName, file)
			tut := &testing.T{}

			method.AssertTestNamedCorrectly(tut)

			assert.Equal(t, !tt.correct, tut.Failed())
		})
	}

	tests5 := []struct {
		methodName string
		correct    bool
	}{
		{methodName: "TestInt_abc", correct: true},
		{methodName: "TestInt_TestInt_Test", correct: true},
		{methodName: "TestInt_", correct: false},
		{methodName: "ATestInt_", correct: false},
		{methodName: "TestInt", correct: false},
		{methodName: "Test_", correct: false},
		{methodName: "Test", correct: false},
		{methodName: "Test_asdf", correct: false},
		{methodName: "TestIntt_", correct: false},
		{methodName: "TestAcc_Abc", correct: false},
	}
	for _, tt := range tests5 {
		t.Run(fmt.Sprintf("intagration test name assertions for method %s", tt.methodName), func(t *testing.T) {
			file := architest.NewFileFromPath("testdata/dir1/sample1.go")
			method := architest.NewMethod(tt.methodName, file)
			tut := &testing.T{}

			method.AssertIntegrationTestNamedCorrectly(tut)

			assert.Equal(t, !tt.correct, tut.Failed())
		})
	}
}

func Test_SampleArchiTestUsage(t *testing.T) {
	t.Run("acceptance tests", func(t *testing.T) {
		acceptanceTestFiles := architest.Directory("testdata/dir3/").
			AllFiles().
			Filter(architest.FileNameRegexFilterProvider(architest.AcceptanceTestFileRegex))

		acceptanceTestFiles.All(func(file *architest.File) {
			file.AssertHasPackage(t, "dir3_test")
			file.ExportedMethods().All(func(method *architest.Method) {
				method.AssertAcceptanceTestNamedCorrectly(t)
			})
		})
	})

	t.Run("integration tests", func(t *testing.T) {
		integrationTestFiles := architest.Directory("testdata/dir4/").
			AllFiles().
			Filter(architest.FileNameRegexFilterProvider(architest.IntegrationTestFileRegex))

		integrationTestFiles.All(func(file *architest.File) {
			file.AssertHasPackage(t, "dir4_test")
			file.ExportedMethods().All(func(method *architest.Method) {
				method.AssertIntegrationTestNamedCorrectly(t)
			})
		})
	})

	t.Run("tests", func(t *testing.T) {
		testFiles := architest.Directory("testdata/dir2/").
			AllFiles().
			Filter(architest.FileNameRegexFilterProvider(architest.TestNameRegex))

		testFiles.All(func(file *architest.File) {
			file.AssertHasPackage(t, "dir2_test")
			file.ExportedMethods().All(func(method *architest.Method) {
				method.AssertTestNamedCorrectly(t)
			})
		})
	})
}
