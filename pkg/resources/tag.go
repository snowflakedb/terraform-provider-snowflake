package resources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

var tagSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the tag; must be unique for the database in which the tag is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the tag."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the tag."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the tag.",
	},
	"allowed_values": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Set of allowed values for the tag.",
	},
	"masking_policies": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		},
		Optional:         true,
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("masking_policies"),
		Description:      relatedResourceDescription("Set of masking policies for the tag. A tag can support one masking policy for each data type. If masking policies are assigned to the tag, before dropping the tag, the provider automatically unassigns them.", resources.MaskingPolicy),
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW TAGS` for the given tag.",
		Elem: &schema.Resource{
			Schema: schemas.ShowTagSchema,
		},
	},
}

// TODO(SNOW-1348114, SNOW-1348110, SNOW-1348355, SNOW-1348353): remove after rework of external table, materialized view, stage and table
var tagReferenceSchema = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	MinItems:    0,
	Description: "Definitions of a tag to associate with the resource.",
	Deprecated:  "Use the 'snowflake_tag_association' resource instead.",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag name, e.g. department.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag value, e.g. marketing_info.",
			},
			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the database that the tag was created in.",
			},
			"schema": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the schema that the tag was created in.",
			},
		},
	},
}

func Tag() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: TrackingCreateWrapper(resources.Tag, CreateContextTag),
		ReadContext:   TrackingReadWrapper(resources.Tag, ReadContextTag),
		UpdateContext: TrackingUpdateWrapper(resources.Tag, UpdateContextTag),
		DeleteContext: TrackingDeleteWrapper(resources.Tag, DeleteContextTag),
		Description:   "Resource used to manage tags. For more information, check [tag documentation](https://docs.snowflake.com/en/sql-reference/sql/create-tag). For assigning tags to Snowflake objects, see [tag_association resource](./tag_association).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Tag, customdiff.All(
			ComputedIfAnyAttributeChanged(tagSchema, ShowOutputAttributeName, "name", "comment", "allowed_values"),
			ComputedIfAnyAttributeChanged(tagSchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Schema: tagSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Tag, ImportName[sdk.SchemaObjectIdentifier]),
		},

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName,
			},
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateContextTag(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schemaName := d.Get("schema").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	request := sdk.NewCreateTagRequest(id)
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(sdk.String(v.(string)))
	}
	if v, ok := d.GetOk("allowed_values"); ok {
		request.WithAllowedValues(expandStringListAllowEmpty(v.(*schema.Set).List()))
	}
	if err := client.Tags.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))
	if v, ok := d.GetOk("masking_policies"); ok {
		ids, err := parseSchemaObjectIdentifierSet(v)
		if err != nil {
			return diag.FromErr(err)
		}
		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(sdk.NewTagSetRequest().WithMaskingPolicies(ids)))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error setting masking policies in tag %v err = %w", id.Name(), err))
		}
	}
	return ReadContextTag(ctx, d, meta)
}

func ReadContextTag(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tag, err := client.Tags.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query tag. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Tag id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	errs := errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.TagToSchema(tag)}),
		d.Set("comment", tag.Comment),
		d.Set("allowed_values", tag.AllowedValues),
		func() error {
			policyRefs, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(id, sdk.PolicyEntityDomainTag))
			if err != nil {
				return (fmt.Errorf("getting policy references for view: %w", err))
			}
			policyIds := make([]string, 0, len(policyRefs))
			for _, p := range policyRefs {
				if p.PolicyKind == sdk.PolicyKindMaskingPolicy {
					policyId := sdk.NewSchemaObjectIdentifier(*p.PolicyDb, *p.PolicySchema, p.PolicyName)
					policyIds = append(policyIds, policyId.FullyQualifiedName())
				}
			}
			return d.Set("masking_policies", policyIds)
		}(),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	return nil
}

func UpdateContextTag(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithRename(newId))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming tag %v err = %w", d.Id(), err))
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}
	if d.HasChange("comment") {
		comment, ok := d.GetOk("comment")
		if ok {
			set := sdk.NewTagSetRequest().WithComment(comment.(string))
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(set)); err != nil {
				return diag.FromErr(err)
			}
		} else {
			unset := sdk.NewTagUnsetRequest().WithComment(true)
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(unset)); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("allowed_values") {
		o, n := d.GetChange("allowed_values")
		oldAllowedValues := expandStringListAllowEmpty(o.(*schema.Set).List())
		newAllowedValues := expandStringListAllowEmpty(n.(*schema.Set).List())

		addedItems, removedItems := ListDiff(oldAllowedValues, newAllowedValues)

		if len(addedItems) > 0 {
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithAdd(addedItems)); err != nil {
				return diag.FromErr(err)
			}
		}

		if len(removedItems) > 0 {
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithDrop(removedItems)); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("masking_policies") {
		o, n := d.GetChange("masking_policies")
		oldAllowedValues := expandStringList(o.(*schema.Set).List())
		newAllowedValues := expandStringList(n.(*schema.Set).List())

		addedItems, removedItems := ListDiff(oldAllowedValues, newAllowedValues)

		removedids := make([]sdk.SchemaObjectIdentifier, len(removedItems))
		for i, idRaw := range removedItems {
			id, err := sdk.ParseSchemaObjectIdentifier(idRaw)
			if err != nil {
				return diag.FromErr(err)
			}
			removedids[i] = id
		}

		addedids := make([]sdk.SchemaObjectIdentifier, len(addedItems))
		for i, idRaw := range addedItems {
			id, err := sdk.ParseSchemaObjectIdentifier(idRaw)
			if err != nil {
				return diag.FromErr(err)
			}
			addedids[i] = id
		}

		if len(removedItems) > 0 {
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(sdk.NewTagUnsetRequest().WithMaskingPolicies(removedids))); err != nil {
				return diag.FromErr(err)
			}
		}

		if len(addedItems) > 0 {
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(sdk.NewTagSetRequest().WithMaskingPolicies(addedids))); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return ReadContextTag(ctx, d, meta)
}

func DeleteContextTag(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// before dropping the resource, all policies must be unset
	policyRefs, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(id, sdk.PolicyEntityDomainTag))
	if err != nil {
		return diag.FromErr(fmt.Errorf("getting policy references for view: %w", err))
	}
	removedPolicies := make([]sdk.SchemaObjectIdentifier, 0, len(policyRefs))
	for _, p := range policyRefs {
		if p.PolicyKind == sdk.PolicyKindMaskingPolicy {
			policyName := sdk.NewSchemaObjectIdentifier(*p.PolicyDb, *p.PolicySchema, p.PolicyName)
			removedPolicies = append(removedPolicies, policyName)
		}
	}

	if len(removedPolicies) > 0 {
		log.Printf("[DEBUG] unsetting masking policies before dropping tag: %s", id.FullyQualifiedName())
		if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(sdk.NewTagUnsetRequest().WithMaskingPolicies(removedPolicies))); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := client.Tags.DropSafely(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
