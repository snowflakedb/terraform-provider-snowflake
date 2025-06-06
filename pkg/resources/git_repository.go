package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var gitRepositorySchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the git repository; must be unique for the schema in which the git repository is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the git repository."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the git repository."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"origin": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the origin URL of the remote Git repository that this Git repository clone represents. The URL must use HTTPS.",
	},
	"api_integration": {
		Type:             schema.TypeString,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Required:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "Specifies the API INTEGRATION that contains information about the remote Git repository such as allowed credentials and prefixes for target URLs.",
	},
	"git_credentials": {
		Type:             schema.TypeString,
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		Optional:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "Specifies the Snowflake secret containing the credentials to use for authenticating with the remote Git repository. Omit this parameter to use the default secret specified by the API integration or if this integration does not require authentication.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the git repository.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW GIT REPOSITORIES` for the given git repository.",
		Elem: &schema.Resource{
			Schema: schemas.ShowGitRepositorySchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE GIT REPOSITORY` for the given git repository.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeGitRepositorySchema,
		},
	},
}

func GitRepository() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.GitRepositories.DropSafely
		},
	)
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.GitRepositoryResource), TrackingCreateWrapper(resources.GitRepository, CreateGitRepository)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.GitRepositoryResource), TrackingReadWrapper(resources.GitRepository, ReadGitRepository)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.GitRepositoryResource), TrackingUpdateWrapper(resources.GitRepository, UpdateGitRepository)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.GitRepositoryResource), TrackingDeleteWrapper(resources.GitRepository, deleteFunc)),
		Description:   "Resource used to manage git repositories. For more information, check [git repositories documentation](https://docs.snowflake.com/en/sql-reference/sql/create-git-repository).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.GitRepository, customdiff.All(
			ComputedIfAnyAttributeChanged(gitRepositorySchema, ShowOutputAttributeName, "origin", "api_integration", "git_credentials", "comment"),
		)),

		Schema: gitRepositorySchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.GitRepository, ImportName[sdk.SchemaObjectIdentifier]),
		},

		Timeouts: defaultTimeouts,
	}
}

func CreateGitRepository(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schemaName := d.Get("schema").(string)
	databaseName := d.Get("database").(string)
	origin := d.Get("origin").(string)

	apiIntegration := sdk.NewAccountObjectIdentifier(d.Get("api_integration").(string))
	gitRepositoryName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	request := sdk.NewCreateGitRepositoryRequest(gitRepositoryName, origin, apiIntegration)
	errs := errors.Join(
		schemaObjectIdentifierAttributeCreate(d, "git_credentials", &request.GitCredentials),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.GitRepositories.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(gitRepositoryName))
	return ReadGitRepository(ctx, d, meta)
}

func ReadGitRepository(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	gitRepository, err := client.GitRepositories.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query git repository. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Git repository id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	gitRepositoryDetails, err := client.GitRepositories.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	var gitCredentials string
	if gitRepository.GitCredentials != nil {
		gitCredentials = gitRepository.GitCredentials.FullyQualifiedName()
	}

	errs := errors.Join(
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.GitRepositoryToSchema(gitRepository)}),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.GitRepositoryDetailsToSchema(*gitRepositoryDetails)}),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set("origin", gitRepository.Origin),
		d.Set("api_integration", gitRepository.ApiIntegration.FullyQualifiedName()),
		d.Set("git_credentials", gitCredentials),
		d.Set("comment", gitRepository.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	return nil
}

func UpdateGitRepository(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewGitRepositorySetRequest(), sdk.NewGitRepositoryUnsetRequest()
	errs := errors.Join(
		accountObjectIdentifierAttributeSetOnly(d, "api_integration", &set.ApiIntegration),
		schemaObjectIdentifierAttributeUpdate(d, "git_credentials", &set.GitCredentials, &unset.GitCredentials),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if (*set != sdk.GitRepositorySetRequest{}) {
		if err := client.GitRepositories.Alter(ctx, sdk.NewAlterGitRepositoryRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if (*unset != sdk.GitRepositoryUnsetRequest{}) {
		if err := client.GitRepositories.Alter(ctx, sdk.NewAlterGitRepositoryRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadGitRepository(ctx, d, meta)
}
