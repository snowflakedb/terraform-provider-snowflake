package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var tagAssociationSchema = map[string]*schema.Schema{
	"object_identifiers": {
		Type:        schema.TypeSet,
		MinItems:    1,
		Required:    true,
		Description: "Specifies the object identifiers for the tag association.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("object_identifiers"),
	},
	"object_type": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      fmt.Sprintf("Specifies the type of object to add a tag. Allowed object types: %v.", sdk.TagAssociationAllowedObjectTypesString),
		ValidateFunc:     validation.StringInSlice(sdk.TagAssociationAllowedObjectTypesString, true),
		DiffSuppressFunc: ignoreCaseSuppressFunc,
		ForceNew:         true,
	},
	"tag_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the identifier for the tag.",
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"tag_value": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the value of the tag, (e.g. 'finance' or 'engineering')",
	},
	"skip_validation": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If true, skips validation of the tag association.",
		Default:     true,
	},
}

// TagAssociation returns a pointer to the resource representing a schema.
func TagAssociation() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: TrackingCreateWrapper(resources.TagAssociation, CreateContextTagAssociation),
		ReadContext:   TrackingReadWrapper(resources.TagAssociation, ReadContextTagAssociation),
		UpdateContext: TrackingUpdateWrapper(resources.TagAssociation, UpdateContextTagAssociation),
		DeleteContext: TrackingDeleteWrapper(resources.TagAssociation, DeleteContextTagAssociation),

		Description: "Resource used to manage tag associations. For more information, check [object tagging documentation](https://docs.snowflake.com/en/user-guide/object-tagging).",

		Schema: tagAssociationSchema,
		Importer: &schema.ResourceImporter{
			StateContext: ImportTagAssociation,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Minute),
			Update: schema.DefaultTimeout(defaultUpdateTimeout),
			Delete: schema.DefaultTimeout(defaultDeleteTimeout),
			Read:   schema.DefaultTimeout(defaultReadTimeout),
		},

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v0_98_0_TagAssociationStateUpgrader,
			},
		},
	}
}

func ImportTagAssociation(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] Starting tag association import")
	idParts := helpers.ParseResourceIdentifier(d.Id())
	if len(idParts) != 3 {
		return nil, fmt.Errorf("invalid resource id: expected 3 arguments, but got %d", len(idParts))
	}
	objectType, err := sdk.ToObjectType(idParts[2])
	if err != nil {
		return nil, err
	}

	if err := d.Set("tag_id", idParts[0]); err != nil {
		return nil, err
	}
	if err := d.Set("tag_value", idParts[1]); err != nil {
		return nil, err
	}
	if err := d.Set("object_type", objectType); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func TagIdentifierAndObjectIdentifier(d *schema.ResourceData) (sdk.SchemaObjectIdentifier, []sdk.ObjectIdentifier, sdk.ObjectType, error) {
	tag := d.Get("tag_id").(string)
	tagId, err := sdk.ParseSchemaObjectIdentifier(tag)
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, nil, "", fmt.Errorf("invalid tag id: %w", err)
	}

	objectType, err := sdk.ToObjectType(d.Get("object_type").(string))
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, nil, "", err
	}

	ids, err := ExpandObjectIdentifierSet(d.Get("object_identifiers").(*schema.Set).List(), objectType)
	if err != nil {
		return sdk.SchemaObjectIdentifier{}, nil, "", err
	}

	return tagId, ids, objectType, nil
}

func CreateContextTagAssociation(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	tagValue := d.Get("tag_value").(string)

	tagId, ids, objectType, err := TagIdentifierAndObjectIdentifier(d)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, oid := range ids {
		request := sdk.NewSetTagRequest(objectType, oid).WithSetTags([]sdk.TagAssociation{
			{
				Name:  tagId,
				Value: tagValue,
			},
		})
		if err := client.Tags.Set(ctx, request); err != nil {
			return diag.FromErr(err)
		}
		skipValidate := d.Get("skip_validation").(bool)
		if !skipValidate {
			log.Println("[DEBUG] validating tag creation")
			if err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate)-time.Minute, func() *retry.RetryError {
				tag, err := client.SystemFunctions.GetTag(ctx, tagId, oid, objectType)
				if err != nil {
					return retry.NonRetryableError(fmt.Errorf("error getting tag: %w", err))
				}
				// if length of response is zero, tag association was not found. retry
				if tag == nil {
					return retry.RetryableError(fmt.Errorf("expected tag association to be created but not yet created"))
				}
				return nil
			}); err != nil {
				return diag.FromErr(fmt.Errorf("error validating tag creation: %w", err))
			}
		}
	}
	d.SetId(helpers.EncodeResourceIdentifier(tagId.FullyQualifiedName(), tagValue, string(objectType)))
	return ReadContextTagAssociation(ctx, d, meta)
}

func ReadContextTagAssociation(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	tagValue := d.Get("tag_value").(string)

	tagId, ids, objectType, err := TagIdentifierAndObjectIdentifier(d)
	if err != nil {
		return diag.FromErr(err)
	}
	var correctObjectIds []string
	for _, oid := range ids {
		objectTagValue, err := client.SystemFunctions.GetTag(ctx, tagId, oid, objectType)
		if err != nil {
			return diag.FromErr(err)
		}
		if objectTagValue != nil && *objectTagValue == tagValue {
			correctObjectIds = append(correctObjectIds, oid.FullyQualifiedName())
		}
	}
	if err := d.Set("object_identifiers", correctObjectIds); err != nil {
		return diag.FromErr(err)
	}
	// ensure that object_type is upper case in the state
	if err := d.Set("object_type", objectType); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func UpdateContextTagAssociation(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	tagId, _, objectType, err := TagIdentifierAndObjectIdentifier(d)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChanges("object_identifiers", "tag_value") {
		tagValue := d.Get("tag_value").(string)

		o, n := d.GetChange("object_identifiers")

		oldIds, err := ExpandObjectIdentifierSet(o.(*schema.Set).List(), objectType)
		if err != nil {
			return diag.FromErr(err)
		}

		newIds, err := ExpandObjectIdentifierSet(n.(*schema.Set).List(), objectType)
		if err != nil {
			return diag.FromErr(err)
		}

		oldIdsFQN := collections.Map(oldIds, sdk.ObjectIdentifier.FullyQualifiedName)
		newIdsFQN := collections.Map(newIds, sdk.ObjectIdentifier.FullyQualifiedName)
		addedIdsFQN, removedIdsFQN, commonIdsFQN := ListDiffWithCommonItems(oldIdsFQN, newIdsFQN)

		addedIds, err := collections.MapErr(addedIdsFQN, sdk.ParseObjectIdentifierString)
		if err != nil {
			return diag.FromErr(err)
		}

		removedIds, err := collections.MapErr(removedIdsFQN, sdk.ParseObjectIdentifierString)
		if err != nil {
			return diag.FromErr(err)
		}

		commonIds, err := collections.MapErr(commonIdsFQN, sdk.ParseObjectIdentifierString)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, id := range addedIds {
			request := sdk.NewSetTagRequest(objectType, id).WithSetTags([]sdk.TagAssociation{
				{
					Name:  tagId,
					Value: tagValue,
				},
			})
			if err := client.Tags.Set(ctx, request); err != nil {
				return diag.FromErr(err)
			}
		}

		for _, id := range removedIds {
			if objectType == sdk.ObjectTypeColumn {
				skip, err := skipColumnIfDoesNotExist(ctx, client, id)
				if err != nil {
					return diag.FromErr(err)
				}
				if skip {
					continue
				}
			}
			request := sdk.NewUnsetTagRequest(objectType, id).WithUnsetTags([]sdk.ObjectIdentifier{tagId}).WithIfExists(true)
			if err := client.Tags.Unset(ctx, request); err != nil {
				return diag.FromErr(err)
			}
		}

		if d.HasChange("tag_value") {
			for _, id := range commonIds {
				request := sdk.NewSetTagRequest(objectType, id).WithSetTags([]sdk.TagAssociation{
					{
						Name:  tagId,
						Value: tagValue,
					},
				})
				if err := client.Tags.Set(ctx, request); err != nil {
					return diag.FromErr(err)
				}
			}
			d.SetId(helpers.EncodeResourceIdentifier(tagId.FullyQualifiedName(), tagValue, string(objectType)))
		}
	}

	return ReadContextTagAssociation(ctx, d, meta)
}

func DeleteContextTagAssociation(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	tagId, ids, objectType, err := TagIdentifierAndObjectIdentifier(d)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, id := range ids {
		if objectType == sdk.ObjectTypeColumn {
			skip, err := skipColumnIfDoesNotExist(ctx, client, id)
			if err != nil {
				return diag.FromErr(err)
			}
			if skip {
				continue
			}
		}
		request := sdk.NewUnsetTagRequest(objectType, id).WithUnsetTags([]sdk.ObjectIdentifier{tagId}).WithIfExists(true)
		if err := client.Tags.Unset(ctx, request); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return nil
}

// we need to skip the column manually, because ALTER COLUMN lacks IF EXISTS
func skipColumnIfDoesNotExist(ctx context.Context, client *sdk.Client, id sdk.ObjectIdentifier) (bool, error) {
	columnId, ok := id.(sdk.TableColumnIdentifier)
	if !ok {
		return false, errors.New("invalid column identifier")
	}
	// TODO [SNOW-1007542]: use SHOW COLUMNS
	_, err := client.Tables.ShowByID(ctx, columnId.SchemaObjectId())
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			log.Printf("[DEBUG] table %s not found, skipping", columnId.SchemaObjectId())
			return true, nil
		}
		return false, err
	}
	columns, err := client.Tables.DescribeColumns(ctx, sdk.NewDescribeTableColumnsRequest(columnId.SchemaObjectId()))
	if err != nil {
		return false, err
	}
	if _, err := collections.FindFirst(columns, func(c sdk.TableColumnDetails) bool {
		return c.Name == columnId.Name()
	}); err != nil {
		log.Printf("[DEBUG] column %s not found in table %s, skipping", columnId.Name(), columnId.SchemaObjectId())
		return true, nil
	}
	return false, nil
}
