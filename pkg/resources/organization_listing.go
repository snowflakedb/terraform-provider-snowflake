package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var organizationListingSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
		ValidateDiagFunc: IsValidListingName,
		Description:      "Specifies the listing identifier (name). It must be unique within the organization, regardless of which Snowflake region the account is located in. Must start with an alphabetic character and cannot contain spaces or special characters except for underscores and hyphens.",
	},
	"manifest": {
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Description: joinWithSpace(
			"Specifies the way manifest is provided for the organization listing. For more information on manifest syntax, see [Listing manifest reference](https://docs.snowflake.com/en/progaccess/listing-manifest-reference).",
			"An organization listing manifest has to specify `title`, `organization_targets`, and `locations`.",
			"Publishing an organization listing also requires `organization_profile`.",
			"For `from_string`, the provider compares the manifest returned by `DESCRIBE LISTING` with the configured one and detects external changes to the configured keys; keys that are only set outside of Terraform are not tracked.",
			"Importing a listing that does not use a stage sets `from_string` to the manifest returned by Snowflake.",
			"External changes to a manifest from a stage are not detected.",
		),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from_string": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      joinWithSpace("Manifest provided as a string. For more information on manifest syntax, see [Listing manifest reference](https://docs.snowflake.com/en/progaccess/listing-manifest-reference). Also, the [multiline string syntax](https://developer.hashicorp.com/terraform/language/expressions/strings#heredoc-strings) is a must here. A proper YAML indentation (2 spaces) is required.", doubleDollarQuotesDescription()),
					ValidateDiagFunc: forbidDoubleDollarQuotes,
					DiffSuppressFunc: suppressEquivalentOrganizationListingManifest,
					ExactlyOneOf:     []string{"manifest.0.from_string", "manifest.0.from_stage"},
				},
				"from_stage": {
					Type:         schema.TypeList,
					Optional:     true,
					MaxItems:     1,
					Description:  "Manifest provided from a given stage. If the manifest file is in the root, only stage needs to be passed. For more information on manifest syntax, see [Listing manifest reference](https://docs.snowflake.com/en/progaccess/listing-manifest-reference). A proper YAML indentation (2 spaces) is required.",
					ExactlyOneOf: []string{"manifest.0.from_string", "manifest.0.from_stage"},
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
	"publish": {
		Type:     schema.TypeString,
		Default:  BooleanDefault,
		Optional: true,
		Description: joinWithSpace(
			"Determines if the listing should be published.",
			"When not set, Snowflake defaults to `true`: the listing is published on creation and again on every manifest change, because `ALTER LISTING ... AS` publishes the listing when `PUBLISH` is omitted.",
			"Has to be set to `false` when the listing is not associated with a share or an application package.",
		),
	},
	"comment": {
		Type: schema.TypeString,
		Description: joinWithSpace(
			"Specifies a comment for the listing.",
			"Snowflake rejects `ALTER LISTING ... SET COMMENT` and `UNSET COMMENT` for organization listings, so the comment is updated together with the inline manifest.",
			"Changing the comment of a listing with a manifest from a stage recreates the listing.",
		),
		Optional: true,
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

func OrganizationListing() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.Listings.DropSafely
		},
	)

	return &schema.Resource{
		Description: "Resource used to manage organization listing objects. For more information, check [CREATE ORGANIZATION LISTING](https://docs.snowflake.com/en/sql-reference/sql/create-organization-listing).",

		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.OrganizationListingResource), TrackingCreateWrapper(resources.OrganizationListing, CreateOrganizationListing)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.OrganizationListingResource), TrackingReadWrapper(resources.OrganizationListing, ReadOrganizationListing)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.OrganizationListingResource), TrackingUpdateWrapper(resources.OrganizationListing, UpdateOrganizationListing)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.OrganizationListingResource), TrackingDeleteWrapper(resources.OrganizationListing, deleteFunc)),

		Schema: organizationListingSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.OrganizationListing, ImportName[sdk.AccountObjectIdentifier]),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.OrganizationListing, forceNewCommentForStagedOrganizationListing),

		Timeouts: defaultTimeouts,
	}
}

func forceNewCommentForStagedOrganizationListing(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	if diff.Id() == "" || !diff.HasChange("comment") {
		return nil
	}
	if fromStage, ok := diff.Get("manifest.0.from_stage").([]any); ok && len(fromStage) > 0 {
		return diff.ForceNew("comment")
	}
	return nil
}

func CreateOrganizationListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	req := sdk.NewCreateOrganizationListingRequest(id)
	withReq := sdk.NewListingWithRequest()

	if errs := errors.Join(
		stringAttributeCreateBuilder(d, "manifest.0.from_string", req.WithAs),
		stringAttributeCreateBuilder(d, "comment", req.WithComment),

		attributeMappedValueCreateBuilder(d, "share", withReq.WithShare, sdk.ParseAccountObjectIdentifier),
		attributeMappedValueCreateBuilder(d, "application_package", withReq.WithApplicationPackage, sdk.ParseAccountObjectIdentifier),
	); errs != nil {
		return diag.FromErr(errs)
	}

	if *withReq != *sdk.NewListingWithRequest() {
		req.WithWith(*withReq)
	}

	// REVIEW defaults to true, and Snowflake requires both PUBLISH and REVIEW to be false for a listing
	// without a share or an application package, so REVIEW follows PUBLISH.
	if publishString := d.Get("publish").(string); publishString != BooleanDefault {
		publish, err := booleanStringToBool(publishString)
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithPublish(publish).WithReview(publish)
	}

	if fromStage, ok := d.GetOk("manifest.0.from_stage"); ok && len(fromStage.([]any)) == 1 {
		fromStageMap := fromStage.([]any)[0].(map[string]any)

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

	if err := client.Listings.CreateOrganization(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadOrganizationListing(ctx, d, meta)
}

func UpdateOrganizationListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
		if err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}

		if err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithRenameTo(newId)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	// The comment can only be changed with ALTER LISTING ... AS, so a comment change resends the inline manifest.
	// Comment changes on a listing with a staged manifest force a new resource instead (see CustomizeDiff).
	manifestChanged := d.HasChange("manifest.0.from_string")
	if manifest := d.Get("manifest.0.from_string").(string); manifest != "" && (manifestChanged || d.HasChange("comment")) {
		req := sdk.NewAlterListingAsRequest(manifest).WithComment(d.Get("comment").(string))

		publish, isSet, err := organizationListingPublishForAlterAs(d, manifestChanged)
		if err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
		// REVIEW follows PUBLISH for the same reason as in create.
		if isSet {
			req.WithPublish(publish).WithReview(publish)
		}

		if err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithAlterListingAs(*req)); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}

	if d.HasChange("manifest") {
		if d.HasChange("manifest.0.from_stage") {
			if fromStage := d.Get("manifest.0.from_stage").([]any); len(fromStage) > 0 {
				fromStageMap := fromStage[0].(map[string]any)

				stage, err := sdk.ParseSchemaObjectIdentifier(fromStageMap["stage"].(string))
				if err != nil {
					d.Partial(true)
					return diag.FromErr(err)
				}

				var location string
				if l, ok := fromStageMap["location"]; ok {
					location = l.(string)
				}

				req := sdk.NewAddListingVersionRequest(sdk.NewStageLocation(stage, location))

				if v, ok := fromStageMap["version_name"]; ok && v.(string) != "" {
					req.WithVersionName(v.(string))
				}

				if v, ok := fromStageMap["version_comment"]; ok && v.(string) != "" {
					req.WithComment(v.(string))
				}

				if err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithAddVersion(*req)); err != nil {
					d.Partial(true)
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChanges("publish") {
		if publishString := d.Get("publish").(string); publishString != BooleanDefault {
			publish, err := booleanStringToBool(publishString)
			if err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}

			if publish {
				if err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithPublish(true)); err != nil {
					d.Partial(true)
					return diag.FromErr(err)
				}
			} else {
				if err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithUnpublish(true)); err != nil {
					d.Partial(true)
					return diag.FromErr(err)
				}
			}
		}
	}

	return ReadOrganizationListing(ctx, d, meta)
}

func ReadOrganizationListing(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	listingDetails, err := client.Listings.Describe(ctx, sdk.NewDescribeListingRequest(id))
	if err != nil {
		return diag.FromErr(err)
	}

	if errs := errors.Join(
		setOrganizationListingManifestFromDescribe(d, listingDetails),
		setOptionalValueWithMapping(d, "share", listingDetails.Share, (*sdk.AccountObjectIdentifier).FullyQualifiedName),
		setOptionalValueWithMapping(d, "application_package", listingDetails.ApplicationPackage, (*sdk.AccountObjectIdentifier).FullyQualifiedName),
		d.Set("publish", booleanStringFromBool(listing.State == sdk.ListingStatePublished)),
		d.Set("comment", listing.Comment),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.ListingToSchema(listing)}),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
	); errs != nil {
		return diag.FromErr(errs)
	}

	return nil
}

// setOrganizationListingManifestFromDescribe keeps the inline manifest in the state while Snowflake still has every
// configured key and value, and replaces it with the manifest returned by DESCRIBE LISTING otherwise, so that the next
// plan shows the external change. Manifests from a stage are left untouched.
func setOrganizationListingManifestFromDescribe(d *schema.ResourceData, listingDetails *sdk.ListingDetails) error {
	describedManifest := listingDetails.ManifestYaml
	// DESCRIBE LISTING returns no manifest when Snowflake fails to render it.
	if strings.TrimSpace(describedManifest) == "" {
		return nil
	}
	if fromStage, ok := d.Get("manifest.0.from_stage").([]any); ok && len(fromStage) > 0 {
		return nil
	}

	describedManifestValue := []any{map[string]any{"from_string": describedManifest}}

	currentManifest := d.Get("manifest.0.from_string").(string)
	if currentManifest == "" {
		// Import: only listings that were never sourced from a stage get the inline manifest.
		if listingDetails.LiveVersionUri != nil || listingDetails.LastCommittedVersionUri != nil {
			return nil
		}
		return d.Set("manifest", describedManifestValue)
	}

	contains, err := organizationListingManifestContains(describedManifest, currentManifest)
	if err != nil || contains {
		return nil
	}
	return d.Set("manifest", describedManifestValue)
}

// organizationListingPublishForAlterAs returns the PUBLISH value for ALTER LISTING ... AS and whether it should be sent.
// An explicitly set publish is always sent. Without it, a manifest change follows the Snowflake default (publish),
// while a comment-only change keeps the publication state from the last read so that it does not publish the listing.
func organizationListingPublishForAlterAs(d *schema.ResourceData, manifestChanged bool) (bool, bool, error) {
	oldPublish, newPublish := d.GetChange("publish")
	if newPublishString := newPublish.(string); newPublishString != BooleanDefault {
		publish, err := booleanStringToBool(newPublishString)
		return publish, true, err
	}
	if manifestChanged {
		return false, false, nil
	}
	if oldPublishString := oldPublish.(string); oldPublishString != BooleanDefault {
		publish, err := booleanStringToBool(oldPublishString)
		return publish, true, err
	}
	return false, false, nil
}
