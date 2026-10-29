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
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func quotaEnumValuesDescription[T ~string](values []T) string {
	return fmt.Sprintf("Valid values are (case-sensitive): %s.", possibleValuesListed(values))
}

// quotaEnumValidation is case-sensitive on purpose: most enums below live inside sets, where DiffSuppressFunc doesn't apply,
// so a value normalized on read would otherwise produce a permanent diff against a differently cased configuration.
func quotaEnumValidation[T ~string](values []T) schema.SchemaValidateDiagFunc {
	allowed := make([]string, len(values))
	for i, value := range values {
		allowed[i] = string(value)
	}
	return validation.ToDiagFunc(validation.StringInSlice(allowed, false))
}

func quotaTagValueSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: description,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tag": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Fully qualified name of the tag, e.g. `\"<database>\".\"<schema>\".\"<tag>\"`. The role managing the quota needs the APPLYBUDGET privilege on the tag.",
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
				},
				"value": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Tag value.",
				},
			},
		},
	}
}

func quotaPerUserLimitSchema(cycle string) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      IntDefault,
		ValidateFunc: validation.IntBetween(1, 1_000_000_000),
		Description:  fmt.Sprintf("Credit limit (a whole number of credits) applied to each user in scope for every %s cycle (UTC-aligned). If not set, the limit is unset.", cycle),
	}
}

var quotaSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the quota; must be unique for the database and schema in which the quota is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the quota."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the quota."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the quota.",
	},
	"shared_resource": {
		Type:     schema.TypeSet,
		Optional: true,
		Description: joinWithSpace(
			"Resources (domains) whose spend counts toward the quota.",
			"A quota can monitor either the WAREHOUSE domain or the AI domains, not both,",
			"and a domain is monitored either as a whole or by its instances.",
		),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"domain": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      joinWithSpace("Monitored domain.", quotaEnumValuesDescription(sdk.AllQuotaSharedResourceDomains)),
					ValidateDiagFunc: quotaEnumValidation(sdk.AllQuotaSharedResourceDomains),
				},
				"instance": {
					Type:     schema.TypeString,
					Optional: true,
					Description: joinWithSpace(
						"A single instance of the domain; if not set, the whole domain is monitored.",
						"For `AI FUNCTION` and `CORTEX CODE`, the instance name (e.g. `AI_CLASSIFY`, `CORTEX_CODE_CLI`).",
						"For `CORTEX AGENT`, the fully qualified agent name.",
						"For `WAREHOUSE`, `SNOWFLAKE INTELLIGENCE`, and `AI GATEWAY`, the object name.",
					),
				},
			},
		},
	},
	"user_tags": {
		Type:          schema.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"excluded_user_tag"},
		Description:   "Applies the quota to users with the given tags. Conflicts with `excluded_user_tag`: Snowflake selects users either by tags or by excluding tags.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"operator": {
					Type:             schema.TypeString,
					Optional:         true,
					Default:          string(sdk.QuotaUserTagsOperatorUnion),
					Description:      joinWithSpace("How the tags are combined.", quotaEnumValuesDescription(sdk.AllQuotaUserTagsOperators)),
					ValidateDiagFunc: quotaEnumValidation(sdk.AllQuotaUserTagsOperators),
				},
				"tag": func() *schema.Schema {
					s := quotaTagValueSchema("Tag and value pairs that select the users.")
					s.Required = true
					s.MinItems = 1
					return s
				}(),
			},
		},
	},
	"excluded_user_tag": func() *schema.Schema {
		s := quotaTagValueSchema("Applies the quota to all users except the ones with the given tag and value. Conflicts with `user_tags`.")
		s.Optional = true
		s.ConflictsWith = []string{"user_tags"}
		return s
	}(),
	"included_users": {
		Type:        schema.TypeSet,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString, DiffSuppressFunc: suppressIdentifierQuoting},
		Description: relatedResourceDescription("Names of individual users the quota applies to. Names are resolved case-insensitively.", resources.User),
	},
	"excluded_users": {
		Type:        schema.TypeSet,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString, DiffSuppressFunc: suppressIdentifierQuoting},
		Description: relatedResourceDescription("Names of individual users excluded from the quota. Names are resolved case-insensitively.", resources.User),
	},
	"per_user_limit_monthly": quotaPerUserLimitSchema("monthly"),
	"per_user_limit_weekly":  quotaPerUserLimitSchema("weekly"),
	"per_user_limit_daily":   quotaPerUserLimitSchema("daily"),
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SNOWFLAKE.CORE.QUOTA` for the given quota.",
		Elem: &schema.Resource{
			Schema: schemas.ShowQuotaSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the quota configuration read with its GET_* methods (scope, config, notifications, and actions).",
		Elem: &schema.Resource{
			Schema: schemas.DescribeQuotaDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func Quota() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Quotas.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.QuotaResource), TrackingCreateWrapper(resources.Quota, CreateQuota)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.QuotaResource), TrackingReadWrapper(resources.Quota, ReadQuota)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.QuotaResource), TrackingUpdateWrapper(resources.Quota, UpdateQuota)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.QuotaResource), TrackingDeleteWrapper(resources.Quota, deleteFunc)),
		Description: joinWithSpace(
			"Resource used to manage per-user quota objects (`SNOWFLAKE.CORE.QUOTA`).",
			"For more information, check [per-user quotas documentation](https://docs.snowflake.com/en/user-guide/budgets/per-user-quotas).",
		),

		Schema: quotaSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Quota, ImportName[sdk.SchemaObjectIdentifier]),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Quota, customdiff.All(
			validateQuotaDomains,
			ComputedIfAnyAttributeChanged(quotaSchema, ShowOutputAttributeName, "name", "comment"),
			// Set and list blocks are excluded from describe_output invalidation, because they produce false-positive
			// HasChange results in ComputedIfAnyAttributeChanged. describe_output is still refreshed after every apply.
			// TODO [SNOW-1648997]: address the above comment (same pattern as external_access_integration)
			ComputedIfAnyAttributeChanged(
				quotaSchema, DescribeOutputAttributeName,
				"included_users",
				"excluded_users",
				"per_user_limit_monthly",
				"per_user_limit_weekly",
				"per_user_limit_daily",
			),
			ComputedIfAnyAttributeChanged(quotaSchema, FullyQualifiedNameAttributeName, "name", "schema", "database"),
		)),
		Timeouts: defaultTimeouts,
	}
}

// validateQuotaDomains rejects configurations Snowflake documents as unsupported but doesn't validate itself.
func validateQuotaDomains(_ context.Context, d *schema.ResourceDiff, _ any) error {
	if !d.NewValueKnown("shared_resource") {
		return nil
	}
	var errs []error
	var monitorsWarehouses, monitorsAi bool
	wholeDomains := make(map[sdk.QuotaSharedResourceDomain]bool)
	domainsWithInstances := make(map[sdk.QuotaSharedResourceDomain]bool)
	unknownInstances := quotaUnknownInstancesByDomain(d.GetRawConfig())
	for _, raw := range d.Get("shared_resource").(*schema.Set).List() {
		resource := raw.(map[string]any)
		if resource["domain"].(string) == "" {
			// an unknown domain (e.g. a reference to another resource) is read as an empty string; it is validated at apply time
			continue
		}
		domain, err := sdk.ToQuotaSharedResourceDomain(resource["domain"].(string))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		switch {
		case resource["instance"].(string) == "" && unknownInstances[domain] > 0:
			// an unknown instance (e.g. a reference to an object created in the same apply) reads as an empty string too.
			// It's an instance, not the whole domain; it's validated when the plan is made again with the known values at apply
			// time, which also catches the whole domain listed next to it (both read the same until then).
			unknownInstances[domain]--
			domainsWithInstances[domain] = true
		case resource["instance"].(string) == "":
			wholeDomains[domain] = true
		default:
			domainsWithInstances[domain] = true
		}
		if _, err := sdk.NewQuotaAddSharedResourceArgsRequestFromInstance(domain, resource["instance"].(string)); err != nil {
			errs = append(errs, fmt.Errorf("invalid instance %q of the %s domain: %w", resource["instance"], domain, err))
		}
		if domain == sdk.QuotaSharedResourceDomainWarehouse {
			monitorsWarehouses = true
		} else {
			monitorsAi = true
		}
	}
	for _, domain := range sdk.AllQuotaSharedResourceDomains {
		if wholeDomains[domain] && domainsWithInstances[domain] {
			errs = append(errs, fmt.Errorf("the %s domain is monitored as a whole, so its instances can't be listed too", domain))
		}
	}
	if monitorsWarehouses && monitorsAi {
		errs = append(errs, errors.New("a quota can monitor either the WAREHOUSE domain or the AI domains, not both"))
	}
	return errors.Join(errs...)
}

func quotaUnknownInstancesByDomain(rawConfig cty.Value) map[sdk.QuotaSharedResourceDomain]int {
	counts := make(map[sdk.QuotaSharedResourceDomain]int)
	if rawConfig.IsNull() || !rawConfig.IsKnown() || !rawConfig.Type().IsObjectType() || !rawConfig.Type().HasAttribute("shared_resource") {
		return counts
	}
	sharedResources := rawConfig.GetAttr("shared_resource")
	if sharedResources.IsNull() || !sharedResources.IsKnown() {
		return counts
	}
	for it := sharedResources.ElementIterator(); it.Next(); {
		_, resource := it.Element()
		if !resource.IsKnown() || resource.IsNull() {
			continue
		}
		domain, instance := resource.GetAttr("domain"), resource.GetAttr("instance")
		if domain.IsKnown() && !domain.IsNull() && !instance.IsKnown() {
			counts[sdk.QuotaSharedResourceDomain(domain.AsString())]++
		}
	}
	return counts
}

func CreateQuota(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := sdk.NewSchemaObjectIdentifier(d.Get("database").(string), d.Get("schema").(string), d.Get("name").(string))

	desired, err := quotaConfigurationFromRaw(d.Get)
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateQuotaRequest(id)
	if err := stringAttributeCreateBuilder(d, "comment", request.WithComment); err != nil {
		return diag.FromErr(err)
	}
	if err := client.Quotas.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating quota %v, err = %w", id.FullyQualifiedName(), err))
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	if err := reconcileQuota(ctx, client, id, &quotaConfiguration{}, desired); err != nil {
		return diag.FromErr(fmt.Errorf("error configuring quota %v, err = %w", id.FullyQualifiedName(), err))
	}

	return ReadQuota(ctx, d, meta)
}

func ReadQuota(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	quota, err := client.Quotas.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query quota. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Quota id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	details, err := client.Quotas.DescribeDetails(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	// The configuration currently in the state is used to keep the user's spelling of equivalent values
	// (e.g. quoted vs. unquoted identifiers).
	current, err := quotaConfigurationFromRaw(d.Get)
	if err != nil {
		return diag.FromErr(err)
	}
	remote, err := quotaConfigurationFromDetails(details, current)
	if err != nil {
		return diag.FromErr(err)
	}

	comment := ""
	if quota.Comment != nil {
		comment = *quota.Comment
	}

	errs := errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set("comment", comment),
		d.Set("shared_resource", remote.sharedResourcesToSchema()),
		d.Set("user_tags", remote.userTagsToSchema()),
		d.Set("excluded_user_tag", quotaTagValuesToSchema(remote.excludedUserTags)),
		d.Set("included_users", remote.includedUsers),
		d.Set("excluded_users", remote.excludedUsers),
		d.Set("per_user_limit_monthly", remote.perUserLimit(sdk.QuotaCycleMonthly)),
		d.Set("per_user_limit_weekly", remote.perUserLimit(sdk.QuotaCycleWeekly)),
		d.Set("per_user_limit_daily", remote.perUserLimit(sdk.QuotaCycleDaily)),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.QuotaToSchema(quota)}),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.QuotaDetailsToSchema(details)}),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	if len(remote.unsupportedSharedResources) > 0 {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Quota monitors shared resources this provider version doesn't support",
				Detail: fmt.Sprintf("Quota %s monitors %s, which can't be managed with this provider version; they're left out of the state and shown in describe_output only.",
					id.FullyQualifiedName(), strings.Join(remote.unsupportedSharedResources, ", ")),
			},
		}
	}
	return nil
}

func UpdateQuota(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))

		if err := client.Quotas.Alter(ctx, sdk.NewAlterQuotaRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming quota from %v to %v, err = %w", id.FullyQualifiedName(), newId.FullyQualifiedName(), err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	if d.HasChange("comment") {
		request := sdk.NewAlterQuotaRequest(id)
		if comment := d.Get("comment").(string); comment != "" {
			request.WithSet(*sdk.NewQuotaSetRequest().WithComment(comment))
		} else {
			request.WithUnset(*sdk.NewQuotaUnsetRequest().WithComment(true))
		}
		if err := client.Quotas.Alter(ctx, request); err != nil {
			return diag.FromErr(fmt.Errorf("error updating the comment of quota %v, err = %w", id.FullyQualifiedName(), err))
		}
	}

	previous, err := quotaConfigurationFromRaw(func(key string) any {
		oldValue, _ := d.GetChange(key)
		return oldValue
	})
	if err != nil {
		return diag.FromErr(err)
	}
	desired, err := quotaConfigurationFromRaw(d.Get)
	if err != nil {
		return diag.FromErr(err)
	}

	// On a failed update the state is saved with the planned values; every value is read back, so the next refresh corrects them.
	if err := reconcileQuota(ctx, client, id, previous, desired); err != nil {
		return diag.FromErr(fmt.Errorf("error updating quota %v, err = %w", id.FullyQualifiedName(), err))
	}

	return ReadQuota(ctx, d, meta)
}

// Removals go first, so that e.g. switching between the WAREHOUSE and AI domains works. The calls aren't transactional:
// if one fails, the quota is left partially configured.
func reconcileQuota(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, previous *quotaConfiguration, desired *quotaConfiguration) error {
	quotas := client.Quotas

	removedResources, addedResources := quotaDiffByKey(previous.sharedResources, desired.sharedResources, quotaSharedResource.key)
	for _, resource := range removedResources {
		args, err := sdk.NewQuotaRemoveSharedResourceArgsRequestFromInstance(resource.domain, resource.instance)
		if err != nil {
			return err
		}
		if _, err := quotas.RemoveSharedResource(ctx, sdk.NewRemoveSharedResourceQuotaRequest(id, *args)); err != nil {
			return fmt.Errorf("removing shared resource %s: %w", resource.description(), err)
		}
	}
	for _, resource := range addedResources {
		args, err := sdk.NewQuotaAddSharedResourceArgsRequestFromInstance(resource.domain, resource.instance)
		if err != nil {
			return err
		}
		if _, err := quotas.AddSharedResource(ctx, sdk.NewAddSharedResourceQuotaRequest(id, *args)); err != nil {
			return fmt.Errorf("adding shared resource %s: %w", resource.description(), err)
		}
	}

	// The users are selected either by tags (SET_USER_TAGS) or by excluding tags (EXCLUDE_USERS('TAG', ...)); each call replaces
	// the other's tags. EXCLUDE_USERS('TAG', ...) also clears the included and excluded users, which are then set again below;
	// with an empty list, it resets the quota to all users.
	includedUsersBefore, excludedUsersBefore := previous.includedUsers, previous.excludedUsers
	if !previous.userTagsEqual(desired) || !quotaKeysEqual(previous.excludedUserTags, desired.excludedUserTags, quotaTagValue.key) {
		if len(desired.userTags) > 0 {
			tags := make([]sdk.QuotaUserTagRequest, len(desired.userTags))
			for i, tag := range desired.userTags {
				tags[i] = *sdk.NewQuotaUserTagRequest(tag.tag, tag.value)
			}
			operator := sdk.QuotaUserTagsOperatorUnion
			if desired.userTagsOperator != nil {
				operator = *desired.userTagsOperator
			}
			if _, err := quotas.SetUserTags(ctx, sdk.NewSetUserTagsQuotaRequest(id, *sdk.NewQuotaSetUserTagsArgsRequest(tags, operator))); err != nil {
				return fmt.Errorf("setting user tags: %w", err)
			}
		} else {
			tags := make([]sdk.QuotaExcludedUserTagRequest, len(desired.excludedUserTags))
			for i, tag := range desired.excludedUserTags {
				tags[i] = *sdk.NewQuotaExcludedUserTagRequest(tag.tag, tag.value)
			}
			args := sdk.NewQuotaExcludeUsersArgsRequest().WithByTag(*sdk.NewQuotaExcludeUsersByTagRequest(tags))
			if _, err := quotas.ExcludeUsers(ctx, sdk.NewExcludeUsersQuotaRequest(id, *args)); err != nil {
				return fmt.Errorf("excluding users by tag: %w", err)
			}
			includedUsersBefore, excludedUsersBefore = nil, nil
		}
	}
	includeUsers := func(users []string) error {
		args := sdk.NewQuotaIncludeUsersArgsRequest(quotaUserNames(users))
		if _, err := quotas.IncludeUsers(ctx, sdk.NewIncludeUsersQuotaRequest(id, *args)); err != nil {
			return fmt.Errorf("including users: %w", err)
		}
		return nil
	}
	if err := quotaReconcileAccumulatingList(includedUsersBefore, desired.includedUsers, quotaUserKey, includeUsers); err != nil {
		return err
	}
	excludeUsers := func(users []string) error {
		args := sdk.NewQuotaExcludeUsersArgsRequest().WithByUser(*sdk.NewQuotaExcludeUsersByUserRequest(quotaUserNames(users)))
		if _, err := quotas.ExcludeUsers(ctx, sdk.NewExcludeUsersQuotaRequest(id, *args)); err != nil {
			return fmt.Errorf("excluding users: %w", err)
		}
		return nil
	}
	if err := quotaReconcileAccumulatingList(excludedUsersBefore, desired.excludedUsers, quotaUserKey, excludeUsers); err != nil {
		return err
	}

	for _, cycle := range sdk.AllQuotaCycles {
		previousLimit, desiredLimit := previous.perUserLimit(cycle), desired.perUserLimit(cycle)
		if previousLimit == desiredLimit {
			continue
		}
		if desiredLimit == IntDefault {
			args := sdk.NewQuotaUnsetPerUserLimitArgsRequest(cycle)
			if _, err := quotas.UnsetPerUserLimit(ctx, sdk.NewUnsetPerUserLimitQuotaRequest(id, *args)); err != nil {
				return fmt.Errorf("unsetting the %s per-user limit: %w", cycle, err)
			}
			continue
		}
		args := sdk.NewQuotaSetPerUserLimitArgsRequest(desiredLimit)
		if cycle != sdk.QuotaCycleMonthly {
			args.WithCycle(cycle)
		}
		if _, err := quotas.SetPerUserLimit(ctx, sdk.NewSetPerUserLimitQuotaRequest(id, *args)); err != nil {
			return fmt.Errorf("setting the %s per-user limit: %w", cycle, err)
		}
	}

	return nil
}
