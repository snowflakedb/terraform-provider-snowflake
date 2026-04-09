package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const tempTagAllowedValue = "SNOWFLAKE_TERRAFORM_TEMP_TAG_ALLOWED_VALUE"

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
		Type:          schema.TypeSet,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		Description:   "Set of allowed values for the tag (unordered). When specified, only these values can be assigned. When the `TAGS_ALLOW_EMPTY_ALLOWED_VALUES` experiment is enabled, removing this field from the configuration reverts the tag to accepting any value. Conflicts with `no_allowed_values` and `ordered_allowed_values`.",
		Deprecated:    "This field is deprecated and will be removed in the next major version. Use `ordered_allowed_values` instead.",
		ConflictsWith: []string{"no_allowed_values", "ordered_allowed_values"},
	},
	"ordered_allowed_values": {
		Type:          schema.TypeList,
		Elem:          &schema.Schema{Type: schema.TypeString},
		Optional:      true,
		Description:   "Ordered list of allowed values for the tag. The order is preserved in Snowflake and is significant when `on_conflict.allowed_values_sequence` is used — the first matching value in the sequence wins. Use this instead of `allowed_values` when order matters. Conflicts with `allowed_values` and `no_allowed_values`.",
		ConflictsWith: []string{"allowed_values", "no_allowed_values"},
	},
	"no_allowed_values": {
		Type:          schema.TypeBool,
		Optional:      true,
		Description:   "When set to true, the tag explicitly disallows any value from being assigned. This is different from omitting `allowed_values`, which means any value is accepted. Available only when the `TAGS_ALLOW_EMPTY_ALLOWED_VALUES` experiment is enabled. Conflicts with `allowed_values` and `ordered_allowed_values`.",
		ConflictsWith: []string{"allowed_values", "ordered_allowed_values"},
	},
	"propagate": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: fmt.Sprintf("Specifies that the tag will be automatically propagated from source objects to target objects. See more about tag propagation in the [official documentation](https://docs.snowflake.com/en/user-guide/object-tagging/propagation). Valid options are: %s", docs.PossibleValuesListed(sdk.AllTagPropagationValues)),
		DiffSuppressFunc: NormalizeAndCompare(func(s string) (sdk.TagPropagation, error) {
			// Removed from config is the same as NONE.
			if s == "" {
				return sdk.TagPropagationNone, nil
			}
			return sdk.ToTagPropagation(s)
		}),
		ValidateDiagFunc: sdkValidation(sdk.ToTagPropagation),
	},
	"on_conflict": {
		Type:         schema.TypeList,
		Optional:     true,
		Description:  externalChangesNotDetectedFieldDescription("Specifies what happens when there is a conflict between the values of [propagated tags](https://docs.snowflake.com/en/user-guide/object-tagging/propagation)."),
		MaxItems:     1,
		RequiredWith: []string{"propagate"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allowed_values_sequence": {
					Type:         schema.TypeBool,
					Optional:     true,
					Description:  externalChangesNotDetectedFieldDescription("The order of the values in the ALLOWED_VALUES property of the tag determines which value is used when there is a conflict."),
					RequiredWith: []string{"ordered_allowed_values"},
					ExactlyOneOf: []string{
						"on_conflict.0.allowed_values_sequence",
						"on_conflict.0.custom_value",
					},
				},
				"custom_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: externalChangesNotDetectedFieldDescription("Whenever there is a conflict, the value of tag is set to custom_value. If `allowed_values` are set, the value set in this field should be one of the values in the `allowed_values` list."),
					ExactlyOneOf: []string{
						"on_conflict.0.allowed_values_sequence",
						"on_conflict.0.custom_value",
					},
				},
			},
		},
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
			ComputedIfAnyAttributeChanged(tagSchema, ShowOutputAttributeName, "name", "comment", "allowed_values", "ordered_allowed_values", "no_allowed_values"),
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
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client

	if d.Get("no_allowed_values").(bool) && !experimentalfeatures.IsExperimentEnabled(experimentalfeatures.TagsAllowEmptyAllowedValues, providerCtx.EnabledExperiments) {
		return diag.FromErr(fmt.Errorf("no_allowed_values is not supported when the %s experiment is disabled", experimentalfeatures.TagsAllowEmptyAllowedValues))
	}

	name := d.Get("name").(string)
	schemaName := d.Get("schema").(string)
	database := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	request := sdk.NewCreateTagRequest(id)

	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	if v, ok := d.GetOk("ordered_allowed_values"); ok {
		request.WithAllowedValues(expandStringListAllowEmpty(v.([]any)))
	} else if v, ok := d.GetOk("allowed_values"); ok {
		request.WithAllowedValues(expandStringListAllowEmpty(v.(*schema.Set).List()))
	}

	if v, ok := d.GetOk("propagate"); ok {
		propagate, err := buildTagPropagateRequest(v.(string), d)
		if err != nil {
			return diag.FromErr(err)
		}
		request.WithPropagate(*propagate)
	}

	if err := client.Tags.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	var updateAfterCreationDiags diag.Diagnostics

	if v, ok := d.GetOk("masking_policies"); ok {
		ids, err := parseSchemaObjectIdentifierSet(v)
		if err != nil {
			updateAfterCreationDiags = append(updateAfterCreationDiags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to parse masking_policies",
				Detail:   fmt.Sprintf("Unable to parse masking policy identifiers for tag %s, err = %s", id.FullyQualifiedName(), err),
			})
		}

		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(*sdk.NewTagSetRequest().WithMaskingPolicies(ids)))
		if err != nil {
			updateAfterCreationDiags = append(updateAfterCreationDiags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to set masking policies on the tag",
				Detail:   fmt.Sprintf("Unable to alter tag %s, err = %s", id.FullyQualifiedName(), err),
			})
		}
	}

	if experimentalfeatures.IsExperimentEnabled(experimentalfeatures.TagsAllowEmptyAllowedValues, providerCtx.EnabledExperiments) {
		if d.Get("no_allowed_values").(bool) {
			// We have to temporarily add and remove allowed value for Snowflake to make the tag block any value.
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithAdd([]string{tempTagAllowedValue})); err != nil {
				updateAfterCreationDiags = append(updateAfterCreationDiags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to add temporary allowed value on the tag",
					Detail:   fmt.Sprintf("Unable to add temporary allowed value to tag %s, err = %s", id.FullyQualifiedName(), err),
				})
			}

			// Drop can run without checking above command's status as Snowflake doesn't fail on dropped values that are not there.
			// The value is also documented to be "reserved" by the provider in case customer would like to use it.
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithDrop([]string{tempTagAllowedValue})); err != nil {
				updateAfterCreationDiags = append(updateAfterCreationDiags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to remove temporary allowed value on the tag",
					Detail:   fmt.Sprintf("Unable to drop temporary allowed value from tag %s, err = %s", id.FullyQualifiedName(), err),
				})
			}
		}
	}

	return append(updateAfterCreationDiags, ReadContextTag(ctx, d, meta)...)
}

func ReadContextTag(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client

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

	if experimentalfeatures.IsExperimentEnabled(experimentalfeatures.TagsAllowEmptyAllowedValues, providerCtx.EnabledExperiments) {
		if err := d.Set("no_allowed_values", reflect.DeepEqual(tag.AllowedValues, []string{})); err != nil {
			return diag.FromErr(err)
		}
	}

	errs := errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.TagToSchema(tag)}),
		d.Set("comment", tag.Comment),
		// Use ordered_allowed_values by default (including import where rawConfig is null).
		// Fall back to allowed_values only when it is explicitly set in the config.
		func() error {
			useAllowedValues := false
			if !d.GetRawConfig().IsNull() {
				v, ok := d.GetRawConfig().AsValueMap()["allowed_values"]
				useAllowedValues = ok && !v.IsNull() && v.LengthInt() > 0
			}
			if useAllowedValues {
				return d.Set("allowed_values", tag.AllowedValues)
			}
			return d.Set("ordered_allowed_values", tag.AllowedValues)
		}(),
		d.Set("propagate", tag.Propagate),
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
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client

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
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(*set)); err != nil {
				return diag.FromErr(err)
			}
		} else {
			unset := sdk.NewTagUnsetRequest().WithComment(true)
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(*unset)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("on_conflict") {
		if len(d.Get("on_conflict").([]any)) == 0 {
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(*sdk.NewTagUnsetRequest().WithOnConflict(true))); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	experimentEnabled := experimentalfeatures.IsExperimentEnabled(experimentalfeatures.TagsAllowEmptyAllowedValues, providerCtx.EnabledExperiments)
	noAllowedValues := d.Get("no_allowed_values").(bool)
	if noAllowedValues && !experimentEnabled {
		return diag.FromErr(fmt.Errorf("no_allowed_values is not supported when the %s experiment is disabled", experimentalfeatures.TagsAllowEmptyAllowedValues))
	}

	allowedValuesSequenceActive := false
	if v, ok := d.GetOk("on_conflict.0.allowed_values_sequence"); ok && v.(bool) {
		allowedValuesSequenceActive = true
	}

	// Handle allowed values changes.
	// The three fields (allowed_values, ordered_allowed_values, no_allowed_values) are mutually exclusive,
	// so exactly one target state applies. Determine it from config and apply the matching Snowflake operation.
	if d.HasChange("ordered_allowed_values") || d.HasChange("allowed_values") || d.HasChange("no_allowed_values") {
		newOrdered := expandStringListAllowEmpty(d.Get("ordered_allowed_values").([]any))
		newUnordered := expandStringListAllowEmpty(d.Get("allowed_values").(*schema.Set).List())

		switch {
		case len(newOrdered) > 0:
			// Target: ordered values - SET atomically replaces in the specified order.
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(*sdk.NewTagSetRequest().WithAllowedValues(newOrdered))); err != nil {
				return diag.FromErr(err)
			}

		case len(newUnordered) > 0:
			// Target: unordered values - SET atomically replaces.
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(*sdk.NewTagSetRequest().WithAllowedValues(newUnordered))); err != nil {
				return diag.FromErr(err)
			}

		case noAllowedValues:
			// Target: block all values (empty allowed_values in Snowflake). Requires experiment flag (validated above).
			oldUnorderedRaw, _ := d.GetChange("allowed_values")
			oldValues := expandStringListAllowEmpty(oldUnorderedRaw.(*schema.Set).List())
			if len(oldValues) == 0 {
				oldOrderedRaw, _ := d.GetChange("ordered_allowed_values")
				oldValues = expandStringListAllowEmpty(oldOrderedRaw.([]any))
			}
			if len(oldValues) > 0 {
				if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithDrop(oldValues)); err != nil {
					return diag.FromErr(err)
				}
			} else {
				// No previous values - use ADD temp + DROP temp to enter blocking state.
				if err := errors.Join(
					client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithAdd([]string{tempTagAllowedValue})),
					client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithDrop([]string{tempTagAllowedValue})),
				); err != nil {
					return diag.FromErr(err)
				}
			}

		default:
			// No allowed values field is set - values were removed from config.
			if experimentEnabled {
				// UNSET makes the tag accept any value (null state).
				if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(*sdk.NewTagUnsetRequest().WithAllowedValues(true))); err != nil {
					return diag.FromErr(err)
				}
			} else {
				// Without experiment: removing values = blocking state (legacy behavior).
				// DROP old values; if already blocking (no old values), do nothing.
				oldUnorderedRaw, _ := d.GetChange("allowed_values")
				oldValues := expandStringListAllowEmpty(oldUnorderedRaw.(*schema.Set).List())
				if len(oldValues) == 0 {
					oldOrderedRaw, _ := d.GetChange("ordered_allowed_values")
					oldValues = expandStringListAllowEmpty(oldOrderedRaw.([]any))
				}
				if len(oldValues) > 0 {
					if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithDrop(oldValues)); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
	}

	// Determine if we need to re-apply propagation settings.
	// Covers:
	// - propagate changed
	// - on_conflict changed (set or unset)
	// - allowed_values changed or reordered while on_conflict is configured
	shouldPropagate := d.HasChange("propagate") || d.HasChange("on_conflict")
	if !shouldPropagate && (d.HasChange("allowed_values") || d.HasChange("ordered_allowed_values")) {
		// If on_conflict is configured, re-send propagation when values change
		// to re-evaluate conflict resolution on dependent objects.
		if allowedValuesSequenceActive {
			shouldPropagate = true
		}
	}

	if shouldPropagate {
		if v, ok := d.GetOk("propagate"); ok {
			propagate, err := buildTagPropagateRequest(v.(string), d)
			if err != nil {
				return diag.FromErr(err)
			}

			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(
				*sdk.NewTagSetRequest().WithPropagate(*propagate),
			)); err != nil {
				return diag.FromErr(err)
			}
		} else if d.HasChange("propagate") {
			// Propagate was removed from config.
			// Note: Snowflake's UNSET PROPAGATE does NOT remove already-propagated tags from dependent objects.
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(
				*sdk.NewTagUnsetRequest().WithPropagate(true),
			)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("masking_policies") {
		o, n := d.GetChange("masking_policies")
		oldAllowedValues := expandStringListAllowEmpty(o.(*schema.Set).List())
		newAllowedValues := expandStringListAllowEmpty(n.(*schema.Set).List())

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
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(*sdk.NewTagUnsetRequest().WithMaskingPolicies(removedids))); err != nil {
				return diag.FromErr(err)
			}
		}

		if len(addedItems) > 0 {
			if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithSet(*sdk.NewTagSetRequest().WithMaskingPolicies(addedids))); err != nil {
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
		if err := client.Tags.Alter(ctx, sdk.NewAlterTagRequest(id).WithUnset(*sdk.NewTagUnsetRequest().WithMaskingPolicies(removedPolicies))); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := client.Tags.DropSafely(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func buildTagPropagateRequest(propagate string, d *schema.ResourceData) (*sdk.TagPropagateRequest, error) {
	tagPropagation, err := sdk.ToTagPropagation(propagate)
	if err != nil {
		return nil, err
	}
	propagateReq := sdk.NewTagPropagateRequest(tagPropagation)
	if v, ok := d.GetOk("on_conflict"); ok && len(v.([]any)) > 0 {
		onConflictMap := v.([]any)[0].(map[string]any)
		if v, ok := onConflictMap["allowed_values_sequence"]; ok && v.(bool) {
			propagateReq.WithOnConflict(sdk.TagOnConflict{
				AllowedValuesSequence: sdk.Bool(true),
			})
		}
		if v, ok := onConflictMap["custom_value"]; ok && v.(string) != "" {
			propagateReq.WithOnConflict(sdk.TagOnConflict{
				CustomValue: sdk.String(v.(string)),
			})
		}
	}
	return propagateReq, nil
}
