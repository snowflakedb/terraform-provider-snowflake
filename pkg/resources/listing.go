package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var listingSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		Description:      "Specifies the listing identifier (name). It must be unique within the organization, regardless of which Snowflake region the account is located in. Must start with an alphabetic character and cannot contain spaces or special characters except for underscores.",
	},
	"manifest": {
		Type:        schema.TypeList,
		Required:    true,
		MaxItems:    1,
		Description: "Specifies the way manifest is provided for the listing. For more information on manifest syntax, see [Listing manifest reference](https://docs.snowflake.com/en/progaccess/listing-manifest-reference).",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from_string": {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "Manifest provided as a string. For more information on manifest syntax, see [Listing manifest reference](https://docs.snowflake.com/en/progaccess/listing-manifest-reference). Also, the [multiline string syntax](https://developer.hashicorp.com/terraform/language/expressions/strings#heredoc-strings) is a must here. A proper YAML indentation (2 spaces) is required.",
					ConflictsWith: []string{"manifest.0.from_stage"},
				},
				"from_stage": {
					Type:          schema.TypeList,
					Optional:      true,
					MaxItems:      1,
					Description:   "Manifest provided as a string. For more information on manifest syntax, see [Listing manifest reference](https://docs.snowflake.com/en/progaccess/listing-manifest-reference). A proper YAML indentation (2 spaces) is required.",
					ConflictsWith: []string{"manifest.0.from_string"},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"stage": {
								Type:             schema.TypeString,
								Required:         true,
								Description:      "Identifier of the stage where the manifest file is located.",
								ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
								DiffSuppressFunc: suppressIdentifierQuoting,
							},
							"location": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Location of the manifest file in the stage. If not specified, the manifest file will be expected to be at the root of the stage.",
							},
							"version_name": {
								Type:     schema.TypeString,
								Optional: true,
								Description: joinWithSpace(
									"Represents manifest version name.",
									"It's case-sensitive and used in manifest versioning.",
									"Version name should be specified or changed whenever any changes in the manifest should be applied to the listing.",
									"Later on the versions of the listing can be analyzed by calling the [SHOW VERSIONS IN LISTING](https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-listing) command. The resource does not track the changes on the specified stage.",
								),
							},
							"version_comment": {
								Type: schema.TypeString,
								Description: joinWithSpace(
									"Specifies a comment for the listing version.",
									"Whenever a new version is created, this comment will be associated with it.",
									"The comment on the version will be visible in the [SHOW VERSIONS IN LISTING](https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-listing) command output.",
								),
								Optional: true,
							},
						},
					},
				},
			},
		},
	},
	"share": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "Specifies the identifier for the share to attach to the listing.",
		ConflictsWith:    []string{"application_package"},
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"application_package": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "Specifies the application package attached to the listing.",
		ConflictsWith:    []string{"share"},
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"review": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Determines if the listing should be submitted for Marketplace Ops review when a new manifest version is specified.",
		Default:     BooleanDefault,
	},
	"publish": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Determines if the listing should be published immediately after it receives approval from Marketplace Ops review and a new manifest version is specified.",
		Default:     BooleanDefault,
	},
	"comment": {
		Type:        schema.TypeString,
		Description: "Specifies a comment for the listing.",
		Optional:    true,
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW LISTINGS` for the given listing.",
		Elem: &schema.Resource{
			Schema: schemas.ShowListingSchema,
		},
	},
}

func Listing() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.Listings.DropSafely
		},
	)

	return &schema.Resource{
		// TODO: Desc
		Description: "",

		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ListingResource), TrackingCreateWrapper(resources.Listing, CreateListing)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ListingResource), TrackingReadWrapper(resources.Listing, ReadListing)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ListingResource), TrackingUpdateWrapper(resources.Listing, UpdateListing)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ListingResource), TrackingDeleteWrapper(resources.Listing, deleteFunc)),

		Schema: listingSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Listing, schema.ImportStatePassthroughContext),
		},

		Timeouts: defaultTimeouts,
	}
}

func CreateListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	req := sdk.NewCreateListingRequest(id)

	manifest := d.Get("manifest").([]any)
	if len(manifest) != 1 {
		return diag.Errorf("manifest must be specified")
	}

	manifestMap := manifest[0].(map[string]any)

	if v, ok := manifestMap["from_string"]; ok {
		fromString := v.(string)
		if fromString == "" {
			return diag.Errorf("manifest.from_string cannot be empty")
		}
		req.WithAs(fromString)
	}

	if v, ok := manifestMap["from_stage"]; ok {
		if len(v.([]any)) != 1 {
			return diag.Errorf("manifest.from_stage cannot be empty when specified")
		}

		fromStageMap := v.([]any)[0].(map[string]any)

		stage, err := sdk.ParseSchemaObjectIdentifier(fromStageMap["stage"].(string))
		if err != nil {
			return diag.FromErr(err)
		}

		var location string
		if l, ok := fromStageMap["location"]; ok {
			location = l.(string)
		}

		req.WithFrom(sdk.NewStageLocation(stage, location))
	}

	if v, ok := d.GetOk("share"); ok {
		share, err := sdk.ParseAccountObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithWith(*sdk.NewListingWithRequest().WithShare(share))
	}

	if v, ok := d.GetOk("application_package"); ok {
		share, err := sdk.ParseAccountObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithWith(*sdk.NewListingWithRequest().WithShare(share))
	}

	if err := client.Listings.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadListing(ctx, d, meta)
}

func UpdateListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		if err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		id = newId
	}

	if d.HasChange("manifest") {
		if d.HasChange("manifest.0.from_string") {
			if manifest := d.Get("manifest.0.from_string").(string); manifest != "" {
				if err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithAlterListingAs(*sdk.NewAlterListingAsRequest(manifest))); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		// TODO: Maybe only version_name should trigger a new version?
		if d.HasChange("manifest.0.from_stage") {
			if v, ok := d.GetOk("manifest.0.from_stage"); ok {
				fromStage := v.([]any)
				if len(fromStage) != 1 {
					return diag.Errorf("manifest.from_stage cannot be empty when specified")
				}
				fromStageMap := fromStage[0].(map[string]any)

				stage, err := sdk.ParseSchemaObjectIdentifier(fromStageMap["stage"].(string))
				if err != nil {
					return diag.FromErr(err)
				}

				var location string
				if l, ok := fromStageMap["location"]; ok {
					location = l.(string)
				}

				req := sdk.NewAddListingVersionRequest("", sdk.NewStageLocation(stage, location))

				//if v, ok := fromStageMap["version_name"]; ok && v.(string) != "" {
				//	req.WithVersionName(v.(string))
				//}

				if v, ok := fromStageMap["version_comment"]; ok && v.(string) != "" {
					req.WithComment(v.(string))
				}

				if err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithAddVersion(*req)); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChanges("review", "publish") {
		// TODO: What do we do here?
	}

	if d.HasChange("comment") {
		if comment := d.Get("comment").(string); comment != "" {
			if err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithSet(*sdk.NewListingSetRequest().WithComment(comment))); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithUnset(*sdk.NewListingUnsetRequest().WithComment(true))); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return ReadListing(ctx, d, meta)
}

func ReadListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	listing, err := client.Listings.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query listing. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Listing id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	_ = listing

	return nil
}
