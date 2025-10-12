package resources

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDbtProject_GitSourceValidation(t *testing.T) {
	testCases := []struct {
		name      string
		gitSource map[string]interface{}
		shouldSet bool
	}{
		{
			name: "valid git source with branch",
			gitSource: map[string]interface{}{
				"repository_url": "https://github.com/user/repo.git",
				"branch":         "main",
				"stage":          "my_stage",
			},
			shouldSet: true,
		},
		{
			name: "valid git source with tag",
			gitSource: map[string]interface{}{
				"repository_url": "https://github.com/user/repo.git",
				"tag":            "v1.0.0",
				"stage":          "my_stage",
			},
			shouldSet: true,
		},
		{
			name: "valid git source with path",
			gitSource: map[string]interface{}{
				"repository_url": "https://github.com/user/repo.git",
				"branch":         "main",
				"path":           "dbt_project",
				"stage":          "my_stage",
				"stage_path":     "my_project",
			},
			shouldSet: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test resource data
			resourceSchema := DbtProject().Schema
			resourceData := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
				"database": "test_db",
				"schema":   "test_schema",
				"name":     "test_project",
				"git_source": []interface{}{
					tc.gitSource,
				},
			})

			// Test that the data was set correctly
			if tc.shouldSet {
				gitSourceList := resourceData.Get("git_source").([]interface{})
				assert.Len(t, gitSourceList, 1, "Should have one git_source entry")

				gitSource := gitSourceList[0].(map[string]interface{})
				assert.Equal(t, tc.gitSource["repository_url"], gitSource["repository_url"])
				assert.Equal(t, tc.gitSource["stage"], gitSource["stage"])
			}
		})
	}
}

func TestDbtProject_ConflictingFields(t *testing.T) {
	resourceSchema := DbtProject().Schema

	// Test that git_source and from have ConflictsWith configured
	gitSourceConflicts := resourceSchema["git_source"].ConflictsWith
	assert.Contains(t, gitSourceConflicts, "from", "git_source should conflict with from")

	fromConflicts := resourceSchema["from"].ConflictsWith
	assert.Contains(t, fromConflicts, "git_source", "from should conflict with git_source")

	// Test that branch and tag conflict within git_source
	gitSourceSchema := resourceSchema["git_source"].Elem.(*schema.Resource).Schema
	branchConflicts := gitSourceSchema["branch"].ConflictsWith
	assert.Contains(t, branchConflicts, "git_source.0.tag", "branch should conflict with tag")

	tagConflicts := gitSourceSchema["tag"].ConflictsWith
	assert.Contains(t, tagConflicts, "git_source.0.branch", "tag should conflict with branch")
}

func TestDbtProject_GitCloneLogic(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "dbt-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock Git repository structure
	repoDir := filepath.Join(tempDir, "mock-repo")
	err = os.MkdirAll(repoDir, 0o755)
	require.NoError(t, err)

	// Create some test files
	testFiles := []string{
		"dbt_project.yml",
		"models/model1.sql",
		"models/model2.sql",
		"macros/macro1.sql",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(repoDir, file)
		err = os.MkdirAll(filepath.Dir(filePath), 0o755)
		require.NoError(t, err)

		err = os.WriteFile(filePath, []byte(fmt.Sprintf("-- Test content for %s", file)), 0o644)
		require.NoError(t, err)
	}

	// Test file walking logic (simulating what happens in the resource)
	var foundFiles []string
	err = filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files (same logic as in the resource)
		if info.IsDir() || info.Name()[0] == '.' {
			return nil
		}

		relPath, err := filepath.Rel(repoDir, path)
		if err != nil {
			return err
		}

		foundFiles = append(foundFiles, relPath)
		return nil
	})

	require.NoError(t, err)
	assert.Len(t, foundFiles, len(testFiles), "Should find all test files")

	// Verify all expected files were found
	for _, expectedFile := range testFiles {
		assert.Contains(t, foundFiles, expectedFile, "Should find file: %s", expectedFile)
	}
}

func TestDbtProject_SchemaValidation(t *testing.T) {
	resource := DbtProject()

	// Test that the resource has the expected schema fields
	assert.NotNil(t, resource.Schema["database"], "Should have database field")
	assert.NotNil(t, resource.Schema["schema"], "Should have schema field")
	assert.NotNil(t, resource.Schema["name"], "Should have name field")
	assert.NotNil(t, resource.Schema["git_source"], "Should have git_source field")
	assert.NotNil(t, resource.Schema["from"], "Should have from field")

	// Test git_source schema structure
	gitSourceSchema := resource.Schema["git_source"].Elem.(*schema.Resource).Schema
	assert.NotNil(t, gitSourceSchema["repository_url"], "git_source should have repository_url")
	assert.NotNil(t, gitSourceSchema["branch"], "git_source should have branch")
	assert.NotNil(t, gitSourceSchema["tag"], "git_source should have tag")
	assert.NotNil(t, gitSourceSchema["path"], "git_source should have path")
	assert.NotNil(t, gitSourceSchema["stage"], "git_source should have stage")
	assert.NotNil(t, gitSourceSchema["stage_path"], "git_source should have stage_path")

	// Test that branch and tag are mutually exclusive
	branchConflicts := gitSourceSchema["branch"].ConflictsWith
	assert.Contains(t, branchConflicts, "git_source.0.tag", "branch should conflict with tag")

	tagConflicts := gitSourceSchema["tag"].ConflictsWith
	assert.Contains(t, tagConflicts, "git_source.0.branch", "tag should conflict with branch")

	// Test that git_source and from are mutually exclusive
	gitSourceConflicts := resource.Schema["git_source"].ConflictsWith
	assert.Contains(t, gitSourceConflicts, "from", "git_source should conflict with from")

	fromConflicts := resource.Schema["from"].ConflictsWith
	assert.Contains(t, fromConflicts, "git_source", "from should conflict with git_source")
}

func TestDbtProject_DefaultValues(t *testing.T) {
	resourceSchema := DbtProject().Schema

	// Test default values
	gitSourceSchema := resourceSchema["git_source"].Elem.(*schema.Resource).Schema

	assert.Equal(t, "main", gitSourceSchema["branch"].Default, "branch should default to 'main'")
	assert.Equal(t, "", gitSourceSchema["path"].Default, "path should default to empty string")
	assert.Equal(t, "", gitSourceSchema["stage_path"].Default, "stage_path should default to empty string")
}

func TestDbtProject_RequiredFields(t *testing.T) {
	resourceSchema := DbtProject().Schema

	// Test required fields at resource level
	assert.True(t, resourceSchema["database"].Required, "database should be required")
	assert.True(t, resourceSchema["schema"].Required, "schema should be required")
	assert.True(t, resourceSchema["name"].Required, "name should be required")

	// Test required fields in git_source
	gitSourceSchema := resourceSchema["git_source"].Elem.(*schema.Resource).Schema
	assert.True(t, gitSourceSchema["repository_url"].Required, "repository_url should be required")
	assert.True(t, gitSourceSchema["stage"].Required, "stage should be required")

	// Test optional fields in git_source
	assert.True(t, gitSourceSchema["branch"].Optional, "branch should be optional")
	assert.True(t, gitSourceSchema["tag"].Optional, "tag should be optional")
	assert.True(t, gitSourceSchema["path"].Optional, "path should be optional")
	assert.True(t, gitSourceSchema["stage_path"].Optional, "stage_path should be optional")
}

func TestDbtProject_ForceNew(t *testing.T) {
	resourceSchema := DbtProject().Schema

	// Test that git_source forces new resource
	assert.True(t, resourceSchema["git_source"].ForceNew, "git_source should force new resource")

	// Test that from also forces new resource for consistency
	assert.True(t, resourceSchema["from"].ForceNew, "from should force new resource")
}
