package resources

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var dbtProjectSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the DBT project; must be unique for the schema in which the DBT project is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the DBT project."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the DBT project."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"git_source": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		ForceNew:    true,
		Description: "Git repository source configuration for the DBT project. The provider will clone the repository and upload files to a stage. Mutually exclusive with 'from'.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"repository_url": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Git repository URL (e.g., https://github.com/user/repo.git)",
				},
				"branch": {
					Type:          schema.TypeString,
					Optional:      true,
					Default:       "main",
					Description:   "Git branch to use (default: main). Mutually exclusive with 'tag'.",
					ConflictsWith: []string{"git_source.0.tag"},
				},
				"tag": {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "Git tag to use. Mutually exclusive with 'branch'.",
					ConflictsWith: []string{"git_source.0.branch"},
				},
				"path": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: "Path within the repository to the DBT project (default: root)",
				},
				"stage": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Snowflake stage where the Git repository files will be uploaded",
				},
				"stage_path": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: "Path within the stage where files will be uploaded (optional)",
				},
			},
		},
		ConflictsWith: []string{"from"},
	},
	"from": {
		Type:          schema.TypeString,
		Optional:      true,
		ForceNew:      true,
		Description:   "Specifies the source location for the DBT project. This can be a Git repository stage, existing DBT project stage, internal named stage, or workspace. Mutually exclusive with 'git_source'.",
		ConflictsWith: []string{"git_source"},
	},
	"default_args": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the default arguments to pass to DBT commands.",
	},
	"default_version": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the default version to use. Can be FIRST, LAST, or VERSION$<num>.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the DBT project.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	"git_repository_fqn": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Fully qualified name of the stage used for Git repository files (when using git_source)",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW DBT PROJECTS` for the given DBT project.",
		Elem: &schema.Resource{
			Schema: schemas.ShowDbtProjectSchema,
		},
	},
}

func DbtProject() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.DbtProjects.DropSafely
		},
	)
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.DbtProject, CreateDbtProject),
		ReadContext:   TrackingReadWrapper(resources.DbtProject, ReadDbtProject),
		UpdateContext: TrackingUpdateWrapper(resources.DbtProject, UpdateDbtProject),
		DeleteContext: TrackingDeleteWrapper(resources.DbtProject, deleteFunc),
		Description: joinWithSpace(
			"Resource used to manage DBT projects. For more information, check [DBT projects documentation](https://docs.snowflake.com/en/sql-reference/sql/create-dbt-project).",
			"DBT projects allow you to manage and execute DBT transformations within Snowflake.",
		),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.DbtProject, customdiff.All(
			ComputedIfAnyAttributeChanged(dbtProjectSchema, ShowOutputAttributeName, "default_args", "default_version", "comment"),
		)),

		Schema: dbtProjectSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.DbtProject, ImportName[sdk.SchemaObjectIdentifier]),
		},

		Timeouts: defaultTimeouts,
	}
}

func CreateDbtProject(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schemaName := d.Get("schema").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	// Handle Git source configuration for automatic repository cloning and file upload
	var stageFQN string
	var fromLocation string

	if gitSourceList, ok := d.GetOk("git_source"); ok && len(gitSourceList.([]interface{})) > 0 {
		gitSource := gitSourceList.([]interface{})[0].(map[string]interface{})

		repositoryUrl := gitSource["repository_url"].(string)
		branch := gitSource["branch"].(string)
		tag := gitSource["tag"].(string)
		path := gitSource["path"].(string)
		stage := gitSource["stage"].(string)
		stagePath := gitSource["stage_path"].(string)

		// Create a temporary directory for cloning the Git repository
		// This directory will be automatically cleaned up when the function exits
		tempDir, err := os.MkdirTemp("", "dbt-git-clone-")
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to create temporary directory: %w", err))
		}
		defer os.RemoveAll(tempDir) // Ensure cleanup even if errors occur

		// Build git clone command with shallow clone for efficiency
		// Use --depth 1 to only clone the latest commit, reducing download time and disk usage
		var gitCmd *exec.Cmd
		if tag != "" {
			// Clone specific tag for stable/versioned deployments
			gitCmd = exec.Command("git", "clone", "--depth", "1", "--branch", tag, repositoryUrl, tempDir)
		} else {
			// Clone specific branch (default: main) for latest development
			gitCmd = exec.Command("git", "clone", "--depth", "1", "--branch", branch, repositoryUrl, tempDir)
		}

		// Execute git clone and capture both stdout and stderr for better error reporting
		if output, err := gitCmd.CombinedOutput(); err != nil {
			return diag.FromErr(fmt.Errorf("failed to clone Git repository %s: %w\nOutput: %s", repositoryUrl, err, string(output)))
		}

		// Determine the source directory within the cloned repository
		// If a path is specified, use that subdirectory; otherwise use the repository root
		sourceDir := tempDir
		if path != "" {
			sourceDir = filepath.Join(tempDir, path)
			// Verify that the specified path exists in the repository
			if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
				return diag.FromErr(fmt.Errorf("specified path '%s' does not exist in the Git repository", path))
			}
		}

		// Build the stage path for uploading files
		// Combine the base stage with optional stage_path for organization
		stageUploadPath := stage
		if stagePath != "" {
			stageUploadPath = fmt.Sprintf("%s/%s", stage, stagePath)
		}

		// Upload files to Snowflake stage using PUT commands
		// Walk through the directory tree and upload each file individually
		// This approach handles directory structures properly and provides granular error handling
		err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories and hidden files (like .git, .gitignore, etc.)
			// Only upload actual files that are part of the DBT project
			if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
				return nil
			}

			// Calculate relative path from source directory to preserve directory structure
			relPath, err := filepath.Rel(sourceDir, path)
			if err != nil {
				return err
			}

			// Build PUT command for this specific file
			// Use file:// protocol for local file paths
			putPattern := fmt.Sprintf("file://%s", path)

			// Determine the destination path in the stage, preserving directory structure
			putDestination := fmt.Sprintf("@%s/%s", stageUploadPath, filepath.Dir(relPath))
			if filepath.Dir(relPath) == "." {
				// File is in the root directory
				putDestination = fmt.Sprintf("@%s/", stageUploadPath)
			} else {
				// File is in a subdirectory, preserve the structure
				putDestination = fmt.Sprintf("@%s/%s/", stageUploadPath, filepath.Dir(relPath))
			}

			// Execute PUT command with options:
			// - AUTO_COMPRESS=FALSE: Keep files uncompressed for DBT compatibility
			// - OVERWRITE=TRUE: Replace existing files to ensure latest version
			putSQL := fmt.Sprintf("PUT '%s' '%s' AUTO_COMPRESS=FALSE OVERWRITE=TRUE", putPattern, putDestination)

			if _, err := client.ExecForTests(ctx, putSQL); err != nil {
				return fmt.Errorf("failed to upload file %s to stage: %w", relPath, err)
			}

			return nil
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to upload files to stage %s: %w", stage, err))
		}

		// Set the stage FQN and FROM location
		stageFQN = stage
		if stagePath != "" {
			fromLocation = fmt.Sprintf("@%s/%s", stage, stagePath)
		} else {
			fromLocation = fmt.Sprintf("@%s", stage)
		}
	}

	request := sdk.NewCreateDbtProjectRequest(id)
	errs := errors.Join(
		stringAttributeCreateBuilder(d, "from", request.WithFrom),
		stringAttributeCreateBuilder(d, "default_args", request.WithDefaultArgs),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	// Use Git-generated FROM location if available
	if fromLocation != "" {
		request.WithFrom(fromLocation)
	}

	if v, ok := d.GetOk("default_version"); ok {
		defaultVersion, err := sdk.ToDbtProjectDefaultVersion(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		request.WithDefaultVersion(defaultVersion)
	}

	if err := client.DbtProjects.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	// Set computed stage FQN (for Git source)
	if stageFQN != "" {
		if err := d.Set("git_repository_fqn", stageFQN); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadDbtProject(ctx, d, meta)
}

func ReadDbtProject(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dbtProject, err := client.DbtProjects.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query DBT project. Marking the resource as removed.",
					Detail:   fmt.Sprintf("DBT project name: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	errs := errors.Join(
		d.Set("name", dbtProject.Name),
		d.Set("database", dbtProject.DatabaseName),
		d.Set("schema", dbtProject.SchemaName),
		d.Set("from", dbtProject.SourceLocation),
		d.Set("default_args", dbtProject.DefaultArgs),
		d.Set("default_version", dbtProject.DefaultVersion),
		d.Set("comment", dbtProject.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	// Extract stage information for Git integration
	// When using git_source, the provider uploads files to a stage and sets the FROM location
	// We need to extract the stage FQN from the source location for the git_repository_fqn field
	if dbtProject.SourceLocation != nil && strings.HasPrefix(*dbtProject.SourceLocation, "@") {
		sourceLocation := *dbtProject.SourceLocation
		// Extract stage FQN from source location like "@STAGE_NAME" or "@STAGE_NAME/path"
		stageFQN := strings.TrimPrefix(sourceLocation, "@")
		// Remove any path components to get just the stage name
		if slashIndex := strings.Index(stageFQN, "/"); slashIndex != -1 {
			stageFQN = stageFQN[:slashIndex]
		}
		if err := d.Set("git_repository_fqn", stageFQN); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(ShowOutputAttributeName, []map[string]any{schemas.DbtProjectToSchema(dbtProject)}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateDbtProject(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("default_args", "default_version", "comment") {
		request := sdk.NewAlterDbtProjectRequest(id)
		set := sdk.NewDbtProjectSetRequest()

		if d.HasChange("default_args") {
			if v := d.Get("default_args").(string); v != "" {
				set.WithDefaultArgs(v)
			}
		}

		if d.HasChange("default_version") {
			if v := d.Get("default_version").(string); v != "" {
				defaultVersion, err := sdk.ToDbtProjectDefaultVersion(v)
				if err != nil {
					return diag.FromErr(err)
				}
				set.WithDefaultVersion(defaultVersion)
			}
		}

		if d.HasChange("comment") {
			if v := d.Get("comment").(string); v != "" {
				set.WithComment(v)
			}
		}

		request.WithSet(*set)

		if err := client.DbtProjects.Alter(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadDbtProject(ctx, d, meta)
}
