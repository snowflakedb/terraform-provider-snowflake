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

var imageRepositorySchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the image repository; must be unique for the schema in which the image repository is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the image repository."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the image repository."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the object.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW IMAGE REPOSITORIES` for the given image repository.",
		Elem: &schema.Resource{
			Schema: schemas.ShowImageRepositorySchema,
		},
	},
}

func ImageRepository() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.ImageRepositories.DropSafely
		},
	)
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ImageRepositoryResource), TrackingCreateWrapper(resources.ImageRepository, CreateImageRepository)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ImageRepositoryResource), TrackingReadWrapper(resources.ImageRepository, ReadImageRepository)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ImageRepositoryResource), TrackingUpdateWrapper(resources.ImageRepository, UpdateImageRepository)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ImageRepositoryResource), TrackingDeleteWrapper(resources.ImageRepository, deleteFunc)),
		Description: longDescription(
			"Resource used to manage image repositories. For more information, check [image repositories documentation](https://docs.snowflake.com/en/sql-reference/sql/create-image-repository).",
			"Snowpark Container Services provides an OCIv2-compliant image registry service and a storage unit call repository to store images.",
			"See [Working with an image registry and repository](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-registry-repository) developer guide for more details.",
		),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ImageRepository, customdiff.All(
			ComputedIfAnyAttributeChanged(imageRepositorySchema, ShowOutputAttributeName, "comment"),
		)),

		Schema: imageRepositorySchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ImageRepository, ImportName[sdk.SchemaObjectIdentifier]),
		},

		Timeouts: defaultTimeouts,
	}
}

func CreateImageRepository(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schemaName := d.Get("schema").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	request := sdk.NewCreateImageRepositoryRequest(id)
	errs := errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	if err := client.ImageRepositories.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadImageRepository(ctx, d, meta)
}

func ReadImageRepository(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	imageRepository, err := client.ImageRepositories.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query image repository. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Image repository id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	errs := errors.Join(
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.ImageRepositoryToSchema(imageRepository)}),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set("comment", imageRepository.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	return nil
}

func UpdateImageRepository(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	set := sdk.NewImageRepositorySetRequest()
	errs := errors.Join(
		stringAttributeUpdateSetOnly(d, "comment", &set.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	if err := client.ImageRepositories.Alter(ctx, sdk.NewAlterImageRepositoryRequest(id).WithSet(*set)); err != nil {
		return diag.FromErr(err)
	}
	return ReadImageRepository(ctx, d, meta)
}
