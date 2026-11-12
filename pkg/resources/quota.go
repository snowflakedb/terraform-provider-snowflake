package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
		Description:      blocklistedCharactersFieldDescription("The database in which to create the quota."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
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
		Set:      quotaSharedResourceHash,
		Description: joinWithSpace(
			"Resources (domains) whose spend counts toward the quota.",
			"A quota can monitor either the WAREHOUSE domain or the AI domains, not both, and a domain either as a whole or by its instances.",
		),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"domain": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      joinWithSpace("Monitored domain.", enumValuesDescription(sdk.AllQuotaSharedResourceDomains)),
					ValidateDiagFunc: sdkValidation(sdk.ToQuotaSharedResourceDomain),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToQuotaSharedResourceDomain),
				},
				"instance": {
					Type:     schema.TypeString,
					Optional: true,
					Description: joinWithSpace(
						"A single instance of the domain; if not set, the whole domain is monitored.",
						"For `AI FUNCTION` and `CORTEX CODE`, the instance name (e.g. `AI_CLASSIFY`, `CORTEX_CODE_CLI`).",
						"For `CORTEX AGENT`, the fully qualified agent name; Snowflake returns only the agent's name, so it's read back only when one agent in the account has it.",
						"For `WAREHOUSE`, `SNOWFLAKE INTELLIGENCE`, and `AI GATEWAY`, the object name.",
					),
					DiffSuppressFunc: suppressQuotaSharedResourceInstance,
				},
			},
		},
	},
	"user_tags": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Applies the quota to users with the given tags.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"operator": {
					Type:             schema.TypeString,
					Optional:         true,
					Default:          string(sdk.QuotaUserTagsOperatorUnion),
					Description:      joinWithSpace("How the tags are combined.", enumValuesDescription(sdk.AllQuotaUserTagsOperators)),
					ValidateDiagFunc: sdkValidation(sdk.ToQuotaUserTagsOperator),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToQuotaUserTagsOperator),
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
		s := quotaTagValueSchema("Excludes the users with the given tag and value from the quota.")
		s.Optional = true
		return s
	}(),
	"included_users": {
		Type:             schema.TypeSet,
		Optional:         true,
		Elem:             &schema.Schema{Type: schema.TypeString},
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompareIdentifiersInSet("included_users"), IgnoreCaseInSet("included_users")),
		Description:      relatedResourceDescription("Names of individual users the quota applies to. Names are resolved case-insensitively.", resources.User),
	},
	"excluded_users": {
		Type:             schema.TypeSet,
		Optional:         true,
		Elem:             &schema.Schema{Type: schema.TypeString},
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompareIdentifiersInSet("excluded_users"), IgnoreCaseInSet("excluded_users")),
		Description:      relatedResourceDescription("Names of individual users excluded from the quota. Names are resolved case-insensitively.", resources.User),
	},
	"per_user_limit_monthly": quotaPerUserLimitSchema("monthly"),
	"per_user_limit_weekly":  quotaPerUserLimitSchema("weekly"),
	"per_user_limit_daily":   quotaPerUserLimitSchema("daily"),
	"block_enforcement_enabled": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
		Description: joinWithSpace(
			"Blocks new requests of users who reached their limit.",
			"It applies only to the AI domains; Snowflake accepts it for the WAREHOUSE domain too, but it has no effect there.",
		),
	},
	"notify_blocked_users": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Sends an email to users when they are blocked. Used only when `block_enforcement_enabled` is true.",
	},
	"notification_threshold": {
		Type:     schema.TypeSet,
		Optional: true,
		Set:      quotaNotificationThresholdHash,
		Description: joinWithSpace(
			"Thresholds (percentages of the limit) at which notifications are sent.",
			"Snowflake keeps one threshold per threshold, spend strategy, and cycle, so two of them differing only in `notify_user` produce a permanent diff.",
		),
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"threshold": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Percentage of the limit.",
				},
				"spend_strategy": quotaSpendStrategySchema(sdk.QuotaSpendStrategyProjected),
				"notify_user": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Notifies the user who crossed the threshold.",
				},
				"cycle": {
					Type:             schema.TypeString,
					Optional:         true,
					Default:          string(sdk.QuotaCycleMonthly),
					Description:      joinWithSpace("Cycle the threshold applies to.", enumValuesDescription(sdk.AllQuotaCycles)),
					ValidateDiagFunc: sdkValidation(sdk.ToQuotaCycle),
					DiffSuppressFunc: NormalizeAndCompare(sdk.ToQuotaCycle),
				},
			},
		},
	},
	"notification_integrations": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		DiffSuppressFunc: NormalizeAndCompareIdentifiersInSet("notification_integrations"),
		Description:      relatedResourceDescription("Webhook or queue notification integrations used to send admin summaries; email integrations aren't accepted (use `admin_emails` for email recipients). USAGE on each integration must be granted to the SNOWFLAKE application.", resources.NotificationIntegration),
	},
	"admin_emails": {
		Type:             schema.TypeSet,
		Optional:         true,
		Elem:             &schema.Schema{Type: schema.TypeString},
		DiffSuppressFunc: IgnoreCaseInSet("admin_emails"),
		Description:      "Email addresses that receive admin summary notifications. Each must be the verified email of a user in the account.",
	},
	"custom_action": {
		Type:        schema.TypeSet,
		Optional:    true,
		Set:         quotaCustomActionHash,
		Description: "Stored procedures called when users cross a threshold of their monthly limit.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"procedure": quotaProcedureSchema(joinWithSpace(
					"Fully qualified name of the procedure with its argument types, e.g. `\"<database>\".\"<schema>\".\"<procedure>\"(VARCHAR)`.",
					"The procedure must use EXECUTE AS OWNER and be granted to the SNOWFLAKE application; its first argument receives the IDs of the users who crossed the threshold.",
				)),
				"arguments": {
					Type:        schema.TypeList,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Additional arguments passed to the procedure.",
				},
				"spend_strategy": quotaSpendStrategySchema(sdk.QuotaSpendStrategyActual),
				"threshold": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Percentage of the monthly limit.",
				},
			},
		},
	},
	"cycle_start_action": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Stored procedure called at the start of each monthly cycle.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"procedure": quotaProcedureSchema(joinWithSpace(
					"Fully qualified name of the procedure with its argument types, e.g. `\"<database>\".\"<schema>\".\"<procedure>\"()`.",
					"The procedure must use EXECUTE AS OWNER and be granted to the SNOWFLAKE application.",
				)),
				"arguments": {
					Type:        schema.TypeList,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Arguments passed to the procedure.",
				},
			},
		},
	},
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

func quotaTagValueSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Set:         quotaTagValueHash,
		Description: description,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tag": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Fully qualified name of the tag, e.g. `\"<database>\".\"<schema>\".\"<tag>\"`. The role managing the quota needs the APPLYBUDGET privilege on the tag.",
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
					DiffSuppressFunc: suppressIdentifierQuoting,
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
		Type:        schema.TypeInt,
		Optional:    true,
		Description: fmt.Sprintf("Credit limit (a whole number of credits) applied to each user in scope for every %s cycle (UTC-aligned). If not set, the limit is unset.", cycle),
	}
}

func quotaSpendStrategySchema(defaultValue sdk.QuotaSpendStrategy) *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Default:          string(defaultValue),
		Description:      joinWithSpace("Whether the threshold applies to actual or projected spend.", enumValuesDescription(sdk.AllQuotaSpendStrategies)),
		ValidateDiagFunc: sdkValidation(sdk.ToQuotaSpendStrategy),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToQuotaSpendStrategy),
	}
}

func quotaProcedureSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Description:      description,
		ValidateDiagFunc: sdkValidation(sdk.ParseSchemaObjectIdentifierWithArguments),
		DiffSuppressFunc: NormalizeAndCompareUsingFunc(sdk.ParseSchemaObjectIdentifierWithArguments, func(a, b sdk.SchemaObjectIdentifierWithArguments) bool {
			return a.FullyQualifiedName() == b.FullyQualifiedName()
		}),
	}
}

func Quota() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Quotas.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.Quota, CreateQuota),
		ReadContext:   TrackingReadWrapper(resources.Quota, ReadQuota),
		UpdateContext: TrackingUpdateWrapper(resources.Quota, UpdateQuota),
		DeleteContext: TrackingDeleteWrapper(resources.Quota, deleteFunc),
		Description: joinWithSpace(
			"Resource used to manage per-user quota objects (`SNOWFLAKE.CORE.QUOTA`).",
			"For more information, check [per-user quotas documentation](https://docs.snowflake.com/en/user-guide/budgets/per-user-quotas).",
		),

		Schema: quotaSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Quota, ImportName[sdk.SchemaObjectIdentifier]),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Quota, customdiff.All(
			TemporaryWorkaroundIdentifierForceNewIfHierarchyRenamesExperimentNotEnabled("database"),
			TemporaryWorkaroundIdentifierForceNewIfHierarchyRenamesExperimentNotEnabled("schema"),
			ComputedIfAnyAttributeChanged(quotaSchema, ShowOutputAttributeName, "name", "comment"),
			// Set and list blocks are excluded from describe_output invalidation, because they produce false-positive
			// HasChange results in ComputedIfAnyAttributeChanged. describe_output is still refreshed after every apply.
			// TODO [SNOW-1648997]: address the above comment (same pattern as external_access_integration)
			ComputedIfAnyAttributeChanged(
				quotaSchema, DescribeOutputAttributeName,
				"per_user_limit_monthly",
				"per_user_limit_weekly",
				"per_user_limit_daily",
				"block_enforcement_enabled",
				"notify_blocked_users",
			),
			ComputedIfAnyAttributeChanged(quotaSchema, FullyQualifiedNameAttributeName, "name", "schema", "database"),
		)),
		Timeouts: defaultTimeouts,
	}
}

func CreateQuota(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := sdk.NewSchemaObjectIdentifier(d.Get("database").(string), d.Get("schema").(string), d.Get("name").(string))

	request := sdk.NewCreateQuotaRequest(id)
	if err := stringAttributeCreateBuilder(d, "comment", request.WithComment); err != nil {
		return diag.FromErr(err)
	}
	if err := client.Quotas.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating quota %v, err = %w", id.FullyQualifiedName(), err))
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	if err := applyQuotaConfiguration(ctx, d, client.Quotas, id); err != nil {
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

	comment := ""
	if quota.Comment != nil {
		comment = *quota.Comment
	}
	sharedResources, unsupportedSharedResources := quotaSharedResourcesToSchema(details.Scope)
	userTags := make([]map[string]any, 0)
	if tags := details.Scope.SelectionTags(); len(tags) > 0 {
		userTags = append(userTags, map[string]any{"operator": details.Scope.UserTags.Operator, "tag": quotaTagValuesToSchema(tags)})
	}

	errs := errors.Join(
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		d.Set("name", id.Name()),
		d.Set("comment", comment),
		d.Set("shared_resource", sharedResources),
		d.Set("user_tags", userTags),
		d.Set("excluded_user_tag", quotaTagValuesToSchema(details.Scope.ExclusionTags())),
		d.Set("included_users", quotaUserNames(details.Scope.IncludedUsers)),
		d.Set("excluded_users", quotaUserNames(details.Scope.ExcludedUsers)),
		d.Set("per_user_limit_monthly", quotaPerUserLimit(details.Config.PerUserLimit)),
		d.Set("per_user_limit_weekly", quotaPerUserLimit(details.Config.PerUserLimitWeekly)),
		d.Set("per_user_limit_daily", quotaPerUserLimit(details.Config.PerUserLimitDaily)),
		d.Set("block_enforcement_enabled", details.Config.BlockEnforcementEnabled),
		d.Set("notify_blocked_users", details.Config.PerUserBlockNotificationsEnabled),
		d.Set("notification_threshold", quotaNotificationThresholdsToSchema(details.NotificationThresholds)),
		d.Set("notification_integrations", quotaNotificationIntegrationNames(details.NotificationIntegrations)),
		d.Set("admin_emails", details.Config.AdminEmails),
		d.Set("custom_action", quotaCustomActionsToSchema(details.CustomActions)),
		d.Set("cycle_start_action", quotaCycleStartActionToSchema(details.CycleStartAction)),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.QuotaToSchema(quota)}),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.QuotaDetailsToSchema(details)}),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	if len(unsupportedSharedResources) > 0 {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Quota monitors shared resources this provider version doesn't support",
				Detail: fmt.Sprintf("Quota %s monitors %s, which can't be managed with this provider version; they're left out of the state and shown in describe_output only.",
					id.FullyQualifiedName(), strings.Join(unsupportedSharedResources, ", ")),
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

	rename := func(currentId, targetId sdk.SchemaObjectIdentifier) error {
		return client.Quotas.Alter(ctx, sdk.NewAlterQuotaRequest(currentId).WithRenameTo(targetId))
	}
	if diags := handleThreeLevelHierarchyRenameIfEnabled(ctx, d, meta, &id, rename, client.Quotas.ShowByID, "quota"); diags != nil {
		return diags
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))
		if err := rename(id, newId); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming quota from %v to %v, err = %w", id.FullyQualifiedName(), newId.FullyQualifiedName(), err))
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set := sdk.NewQuotaSetRequest()
	unset := sdk.NewQuotaUnsetRequest()
	if err := stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment); err != nil {
		return diag.FromErr(err)
	}
	if !reflect.DeepEqual(*set, *sdk.NewQuotaSetRequest()) {
		if err := client.Quotas.Alter(ctx, sdk.NewAlterQuotaRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(fmt.Errorf("error updating the comment of quota %v, err = %w", id.FullyQualifiedName(), err))
		}
	}
	if !reflect.DeepEqual(*unset, *sdk.NewQuotaUnsetRequest()) {
		if err := client.Quotas.Alter(ctx, sdk.NewAlterQuotaRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(fmt.Errorf("error updating the comment of quota %v, err = %w", id.FullyQualifiedName(), err))
		}
	}

	if err := applyQuotaConfiguration(ctx, d, client.Quotas, id); err != nil {
		d.Partial(true)
		return diag.FromErr(fmt.Errorf("error updating quota %v, err = %w", id.FullyQualifiedName(), err))
	}

	return ReadQuota(ctx, d, meta)
}

// applyQuotaConfiguration calls the instance methods for the attributes that changed, which on create are the configured ones.
// Block enforcement is disabled before the other changes and enabled after them, so that users aren't blocked by a partially
// applied configuration. The calls aren't transactional: if one fails, the quota is left partially configured.
func applyQuotaConfiguration(ctx context.Context, d *schema.ResourceData, quotas sdk.Quotas, id sdk.SchemaObjectIdentifier) error {
	blockEnforcementEnabled := d.Get("block_enforcement_enabled").(bool)
	blockEnforcementChanged := d.HasChange("block_enforcement_enabled") || d.HasChange("notify_blocked_users")
	if d.IsNewResource() {
		// a new quota has block enforcement disabled and the user notifications enabled
		blockEnforcementChanged = blockEnforcementEnabled || !d.Get("notify_blocked_users").(bool)
	}
	setBlockEnforcement := func() error {
		args := sdk.NewQuotaSetBlockEnforcementEnabledArgsRequest(blockEnforcementEnabled).WithNotifyUsers(d.Get("notify_blocked_users").(bool))
		_, err := quotas.SetBlockEnforcementEnabled(ctx, sdk.NewSetBlockEnforcementEnabledQuotaRequest(id, *args))
		return err
	}

	steps := []func() error{
		func() error {
			if blockEnforcementChanged && !blockEnforcementEnabled {
				return setBlockEnforcement()
			}
			return nil
		},
		func() error { return updateQuotaSharedResources(ctx, d, quotas, id) },
		func() error { return updateQuotaUserTags(ctx, d, quotas, id) },
		func() error { return updateQuotaExcludedUserTags(ctx, d, quotas, id) },
		func() error {
			return updateQuotaUserList(d, "included_users", func(users []string) error {
				_, err := quotas.IncludeUsers(ctx, sdk.NewIncludeUsersQuotaRequest(id, *sdk.NewQuotaIncludeUsersArgsRequest(users)))
				return err
			})
		},
		func() error {
			return updateQuotaUserList(d, "excluded_users", func(users []string) error {
				args := sdk.NewQuotaExcludeUsersArgsRequest().WithByUser(*sdk.NewQuotaExcludeUsersByUserRequest(users))
				_, err := quotas.ExcludeUsers(ctx, sdk.NewExcludeUsersQuotaRequest(id, *args))
				return err
			})
		},
		func() error { return updateQuotaPerUserLimits(ctx, d, quotas, id) },
		func() error { return updateQuotaNotificationThresholds(ctx, d, quotas, id) },
		func() error { return updateQuotaNotificationIntegrations(ctx, d, quotas, id) },
		func() error {
			if !d.HasChange("admin_emails") {
				return nil
			}
			args := sdk.NewQuotaSetAdminEmailsArgsRequestFromEmails(expandStringList(d.Get("admin_emails").(*schema.Set).List())...)
			_, err := quotas.SetAdminEmails(ctx, sdk.NewSetAdminEmailsQuotaRequest(id, *args))
			return err
		},
		func() error { return updateQuotaCustomActions(ctx, d, quotas, id) },
		func() error { return updateQuotaCycleStartAction(ctx, d, quotas, id) },
		func() error {
			if blockEnforcementChanged && blockEnforcementEnabled {
				return setBlockEnforcement()
			}
			return nil
		},
	}
	for _, step := range steps {
		if err := step(); err != nil {
			return err
		}
	}
	return nil
}

func quotaSetChange(d *schema.ResourceData, key string) (removed []any, added []any) {
	oldRaw, newRaw := d.GetChange(key)
	oldSet, newSet := oldRaw.(*schema.Set), newRaw.(*schema.Set)
	return oldSet.Difference(newSet).List(), newSet.Difference(oldSet).List()
}

func updateQuotaSharedResources(ctx context.Context, d *schema.ResourceData, quotas sdk.Quotas, id sdk.SchemaObjectIdentifier) error {
	if !d.HasChange("shared_resource") {
		return nil
	}
	removed, added := quotaSetChange(d, "shared_resource")
	return removeThenAdd(removed, added,
		func(raw any) error {
			resource := raw.(map[string]any)
			domain, err := sdk.ToQuotaSharedResourceDomain(resource["domain"].(string))
			if err != nil {
				return err
			}
			args, err := sdk.NewQuotaRemoveSharedResourceArgsRequestFromInstance(domain, resource["instance"].(string))
			if err != nil {
				return err
			}
			_, err = quotas.RemoveSharedResource(ctx, sdk.NewRemoveSharedResourceQuotaRequest(id, *args))
			return err
		},
		func(raw any) error {
			resource := raw.(map[string]any)
			domain, err := sdk.ToQuotaSharedResourceDomain(resource["domain"].(string))
			if err != nil {
				return err
			}
			args, err := sdk.NewQuotaAddSharedResourceArgsRequestFromInstance(domain, resource["instance"].(string))
			if err != nil {
				return err
			}
			_, err = quotas.AddSharedResource(ctx, sdk.NewAddSharedResourceQuotaRequest(id, *args))
			return err
		},
	)
}

// SET_USER_TAGS replaces all the tags, and an empty list clears them.
func updateQuotaUserTags(ctx context.Context, d *schema.ResourceData, quotas sdk.Quotas, id sdk.SchemaObjectIdentifier) error {
	if !d.HasChange("user_tags") {
		return nil
	}
	operator := sdk.QuotaUserTagsOperatorUnion
	tags := make([]sdk.QuotaUserTagRequest, 0)
	if userTags := d.Get("user_tags").([]any); len(userTags) > 0 && userTags[0] != nil {
		block := userTags[0].(map[string]any)
		var err error
		if operator, err = sdk.ToQuotaUserTagsOperator(block["operator"].(string)); err != nil {
			return err
		}
		for _, raw := range block["tag"].(*schema.Set).List() {
			tag, value, err := quotaTagValueFromRaw(raw)
			if err != nil {
				return err
			}
			tags = append(tags, *sdk.NewQuotaUserTagRequestFromTag(tag, value))
		}
	}
	_, err := quotas.SetUserTags(ctx, sdk.NewSetUserTagsQuotaRequest(id, *sdk.NewQuotaSetUserTagsArgsRequest(tags, operator)))
	return err
}

// EXCLUDE_USERS('TAG', ...) replaces all the excluded tags, and an empty list clears them.
func updateQuotaExcludedUserTags(ctx context.Context, d *schema.ResourceData, quotas sdk.Quotas, id sdk.SchemaObjectIdentifier) error {
	if !d.HasChange("excluded_user_tag") {
		return nil
	}
	tags := make([]sdk.QuotaExcludedUserTagRequest, 0)
	for _, raw := range d.Get("excluded_user_tag").(*schema.Set).List() {
		tag, value, err := quotaTagValueFromRaw(raw)
		if err != nil {
			return err
		}
		tags = append(tags, *sdk.NewQuotaExcludedUserTagRequestFromTag(tag, value))
	}
	args := sdk.NewQuotaExcludeUsersArgsRequest().WithByTag(*sdk.NewQuotaExcludeUsersByTagRequest(tags))
	_, err := quotas.ExcludeUsers(ctx, sdk.NewExcludeUsersQuotaRequest(id, *args))
	return err
}

// INCLUDE_USERS and EXCLUDE_USERS('USER', ...) add to their lists, and an empty list clears them. Removing a user therefore
// clears the list and sets the remaining users again. Names are passed without identifier quotes.
func updateQuotaUserList(d *schema.ResourceData, key string, set func([]string) error) error {
	if !d.HasChange(key) {
		return nil
	}
	removed, added := quotaSetChange(d, key)
	userNames := func(raw []any) []string {
		names := expandStringList(raw)
		for i, name := range names {
			names[i] = strings.Trim(name, `"`)
		}
		return names
	}
	if len(removed) == 0 {
		return set(userNames(added))
	}
	if err := set([]string{}); err != nil {
		return err
	}
	if desired := d.Get(key).(*schema.Set).List(); len(desired) > 0 {
		return set(userNames(desired))
	}
	return nil
}

func updateQuotaPerUserLimits(ctx context.Context, d *schema.ResourceData, quotas sdk.Quotas, id sdk.SchemaObjectIdentifier) error {
	for key, cycle := range map[string]sdk.QuotaCycle{
		"per_user_limit_monthly": sdk.QuotaCycleMonthly,
		"per_user_limit_weekly":  sdk.QuotaCycleWeekly,
		"per_user_limit_daily":   sdk.QuotaCycleDaily,
	} {
		if !d.HasChange(key) {
			continue
		}
		if limit := d.Get(key).(int); limit != 0 {
			args := sdk.NewQuotaSetPerUserLimitArgsRequest(limit)
			if cycle != sdk.QuotaCycleMonthly {
				args.WithCycle(cycle)
			}
			if _, err := quotas.SetPerUserLimit(ctx, sdk.NewSetPerUserLimitQuotaRequest(id, *args)); err != nil {
				return err
			}
			continue
		}
		if _, err := quotas.UnsetPerUserLimit(ctx, sdk.NewUnsetPerUserLimitQuotaRequest(id, *sdk.NewQuotaUnsetPerUserLimitArgsRequest(cycle))); err != nil {
			return err
		}
	}
	return nil
}

func updateQuotaNotificationThresholds(ctx context.Context, d *schema.ResourceData, quotas sdk.Quotas, id sdk.SchemaObjectIdentifier) error {
	if !d.HasChange("notification_threshold") {
		return nil
	}
	removed, added := quotaSetChange(d, "notification_threshold")
	return removeThenAdd(removed, added,
		func(raw any) error {
			threshold, err := quotaNotificationThresholdFromRaw(raw)
			if err != nil {
				return err
			}
			args := sdk.NewQuotaRemoveNotificationThresholdArgsRequest(threshold.Threshold, threshold.SpendStrategy).WithCycle(threshold.Cycle)
			_, err = quotas.RemoveNotificationThreshold(ctx, sdk.NewRemoveNotificationThresholdQuotaRequest(id, *args))
			return err
		},
		func(raw any) error {
			threshold, err := quotaNotificationThresholdFromRaw(raw)
			if err != nil {
				return err
			}
			args := sdk.NewQuotaAddNotificationThresholdArgsRequest(threshold.Threshold, threshold.SpendStrategy, threshold.NotifyUser, threshold.Cycle)
			_, err = quotas.AddNotificationThreshold(ctx, sdk.NewAddNotificationThresholdQuotaRequest(id, *args))
			return err
		},
	)
}

func updateQuotaNotificationIntegrations(ctx context.Context, d *schema.ResourceData, quotas sdk.Quotas, id sdk.SchemaObjectIdentifier) error {
	if !d.HasChange("notification_integrations") {
		return nil
	}
	removed, added := quotaSetChange(d, "notification_integrations")
	return removeThenAdd(removed, added,
		func(raw any) error {
			integrationId, err := sdk.ParseAccountObjectIdentifier(raw.(string))
			if err != nil {
				return err
			}
			args := sdk.NewQuotaRemoveNotificationIntegrationArgsRequest(integrationId)
			_, err = quotas.RemoveNotificationIntegration(ctx, sdk.NewRemoveNotificationIntegrationQuotaRequest(id, *args))
			return err
		},
		func(raw any) error {
			integrationId, err := sdk.ParseAccountObjectIdentifier(raw.(string))
			if err != nil {
				return err
			}
			args := sdk.NewQuotaAddNotificationIntegrationArgsRequest(integrationId)
			_, err = quotas.AddNotificationIntegration(ctx, sdk.NewAddNotificationIntegrationQuotaRequest(id, *args))
			return err
		},
	)
}

// REMOVE_CUSTOM_ACTIONS(<threshold>, <procedure>) removes the actions of both spend strategies, so the removal is done once per
// threshold and procedure, with the procedure passed back exactly as GET_CUSTOM_ACTIONS returns it, and the configured actions
// that shared it are added again.
func updateQuotaCustomActions(ctx context.Context, d *schema.ResourceData, quotas sdk.Quotas, id sdk.SchemaObjectIdentifier) error {
	if !d.HasChange("custom_action") {
		return nil
	}
	removed, added := quotaSetChange(d, "custom_action")
	removedGroups := make(map[string]quotaCustomAction)
	for _, raw := range removed {
		action, err := quotaCustomActionFromRaw(raw)
		if err != nil {
			return err
		}
		removedGroups[quotaCustomActionGroup(action.Procedure, action.Threshold)] = action
	}
	if len(removedGroups) > 0 {
		registered, err := quotas.GetCustomActions(ctx, sdk.NewGetCustomActionsQuotaRequest(id))
		if err != nil {
			return err
		}
		for _, registeredAction := range registered {
			procedure, err := registeredAction.ProcedureId()
			if err != nil {
				return err
			}
			group := quotaCustomActionGroup(procedure, registeredAction.Threshold)
			if _, ok := removedGroups[group]; !ok {
				continue
			}
			args := sdk.NewQuotaRemoveCustomActionsArgsRequest().WithThreshold(registeredAction.Threshold).WithProcedureFqn(registeredAction.ProcedureFqn)
			if _, err := quotas.RemoveCustomActions(ctx, sdk.NewRemoveCustomActionsQuotaRequest(id, *args)); err != nil {
				return err
			}
			delete(removedGroups, group)
		}
		for _, raw := range d.Get("custom_action").(*schema.Set).List() {
			action, err := quotaCustomActionFromRaw(raw)
			if err != nil {
				return err
			}
			if _, wasRemoved := removedGroups[quotaCustomActionGroup(action.Procedure, action.Threshold)]; wasRemoved || slices.ContainsFunc(added, func(a any) bool { return reflect.DeepEqual(a, raw) }) {
				continue
			}
			for _, removedRaw := range removed {
				removedAction, err := quotaCustomActionFromRaw(removedRaw)
				if err != nil {
					return err
				}
				if quotaCustomActionGroup(removedAction.Procedure, removedAction.Threshold) == quotaCustomActionGroup(action.Procedure, action.Threshold) {
					added = append(added, raw)
					break
				}
			}
		}
	}
	for _, raw := range added {
		action, err := quotaCustomActionFromRaw(raw)
		if err != nil {
			return err
		}
		args := sdk.NewQuotaAddCustomActionArgsRequest(action.Procedure, action.Arguments, action.SpendStrategy, action.Threshold)
		if _, err := quotas.AddCustomAction(ctx, sdk.NewAddCustomActionQuotaRequest(id, *args)); err != nil {
			return err
		}
	}
	return nil
}

func updateQuotaCycleStartAction(ctx context.Context, d *schema.ResourceData, quotas sdk.Quotas, id sdk.SchemaObjectIdentifier) error {
	if !d.HasChange("cycle_start_action") {
		return nil
	}
	cycleStartAction := d.Get("cycle_start_action").([]any)
	if len(cycleStartAction) == 0 || cycleStartAction[0] == nil {
		_, err := quotas.RemoveCycleStartAction(ctx, sdk.NewRemoveCycleStartActionQuotaRequest(id))
		return err
	}
	action := cycleStartAction[0].(map[string]any)
	procedure, err := sdk.ParseSchemaObjectIdentifierWithArguments(action["procedure"].(string))
	if err != nil {
		return err
	}
	args := sdk.NewQuotaSetCycleStartActionArgsRequest(procedure, expandStringListAllowEmpty(action["arguments"].([]any)))
	_, err = quotas.SetCycleStartAction(ctx, sdk.NewSetCycleStartActionQuotaRequest(id, *args))
	return err
}

func quotaTagValueFromRaw(raw any) (sdk.SchemaObjectIdentifier, string, error) {
	tagValue := raw.(map[string]any)
	tag, err := sdk.ParseSchemaObjectIdentifier(tagValue["tag"].(string))
	return tag, tagValue["value"].(string), err
}

func quotaNotificationThresholdFromRaw(raw any) (sdk.QuotaNotificationThreshold, error) {
	threshold := raw.(map[string]any)
	spendStrategy, err := sdk.ToQuotaSpendStrategy(threshold["spend_strategy"].(string))
	if err != nil {
		return sdk.QuotaNotificationThreshold{}, err
	}
	cycle, err := sdk.ToQuotaCycle(threshold["cycle"].(string))
	if err != nil {
		return sdk.QuotaNotificationThreshold{}, err
	}
	return sdk.QuotaNotificationThreshold{
		Threshold:     threshold["threshold"].(int),
		SpendStrategy: spendStrategy,
		NotifyUser:    threshold["notify_user"].(bool),
		Cycle:         cycle,
	}, nil
}

type quotaCustomAction struct {
	Procedure     sdk.SchemaObjectIdentifierWithArguments
	Arguments     []string
	SpendStrategy sdk.QuotaSpendStrategy
	Threshold     int
}

func quotaCustomActionFromRaw(raw any) (quotaCustomAction, error) {
	action := raw.(map[string]any)
	procedure, err := sdk.ParseSchemaObjectIdentifierWithArguments(action["procedure"].(string))
	if err != nil {
		return quotaCustomAction{}, err
	}
	spendStrategy, err := sdk.ToQuotaSpendStrategy(action["spend_strategy"].(string))
	if err != nil {
		return quotaCustomAction{}, err
	}
	return quotaCustomAction{
		Procedure:     procedure,
		Arguments:     expandStringListAllowEmpty(action["arguments"].([]any)),
		SpendStrategy: spendStrategy,
		Threshold:     action["threshold"].(int),
	}, nil
}

func quotaCustomActionGroup(procedure sdk.SchemaObjectIdentifierWithArguments, threshold int) string {
	return fmt.Sprintf("%d/%s", threshold, procedure.FullyQualifiedName())
}

func quotaSharedResourcesToSchema(scope sdk.QuotaScope) ([]map[string]any, []string) {
	sharedResources := make([]map[string]any, 0, len(scope.SharedResources))
	var unsupported []string
	for _, resource := range scope.SharedResources {
		domain, err := resource.NormalizedDomain()
		if err != nil {
			description := resource.Domain
			if !resource.IsWholeDomain() {
				description = fmt.Sprintf("%s (%s)", resource.Domain, resource.Name)
			}
			unsupported = append(unsupported, description)
			continue
		}
		instance := ""
		switch {
		case resource.IsWholeDomain():
		case resource.CortexAgentId != nil:
			instance = resource.CortexAgentId.FullyQualifiedName()
		case slices.Contains(quotaDomainsWithPlainInstances, domain), !strings.ContainsAny(resource.Name, `."`):
			instance = resource.Name
		default:
			instance = sdk.NewAccountObjectIdentifier(resource.Name).FullyQualifiedName()
		}
		sharedResources = append(sharedResources, map[string]any{"domain": string(domain), "instance": instance})
	}
	return sharedResources, unsupported
}

func quotaTagValuesToSchema(tags []sdk.QuotaScopeTag) []map[string]any {
	result := make([]map[string]any, 0)
	for _, tag := range tags {
		for _, value := range tag.TagValues {
			result = append(result, map[string]any{"tag": tag.TagId().FullyQualifiedName(), "value": value})
		}
	}
	return result
}

func quotaUserNames(users []sdk.QuotaScopeUser) []string {
	names := make([]string, len(users))
	for i, user := range users {
		names[i] = user.UserName
	}
	return names
}

func quotaPerUserLimit(limit *int) int {
	if limit == nil {
		return 0
	}
	return *limit
}

func quotaNotificationThresholdsToSchema(thresholds []sdk.QuotaNotificationThreshold) []map[string]any {
	result := make([]map[string]any, len(thresholds))
	for i, threshold := range thresholds {
		result[i] = map[string]any{
			"threshold":      threshold.Threshold,
			"spend_strategy": string(threshold.SpendStrategy),
			"notify_user":    threshold.NotifyUser,
			"cycle":          string(threshold.Cycle),
		}
	}
	return result
}

func quotaNotificationIntegrationNames(integrations []sdk.QuotaNotificationIntegration) []string {
	names := make([]string, len(integrations))
	for i, integration := range integrations {
		names[i] = integration.IntegrationName
	}
	return names
}

func quotaCustomActionsToSchema(actions []sdk.QuotaCustomAction) []map[string]any {
	result := make([]map[string]any, 0, len(actions))
	for _, action := range actions {
		procedure := action.ProcedureFqn
		if id, err := action.ProcedureId(); err == nil {
			procedure = id.FullyQualifiedName()
		}
		result = append(result, map[string]any{
			"procedure":      procedure,
			"arguments":      action.ProcedureArgs,
			"spend_strategy": string(action.SpendStrategy),
			"threshold":      action.Threshold,
		})
	}
	return result
}

func quotaCycleStartActionToSchema(action *sdk.QuotaCycleStartAction) []map[string]any {
	if action == nil || !action.IsSet() {
		return []map[string]any{}
	}
	procedure := *action.ProcedureFqn
	if id, err := action.ProcedureId(); err == nil {
		procedure = id.FullyQualifiedName()
	}
	return []map[string]any{{"procedure": procedure, "arguments": action.ProcedureArgs}}
}

var quotaDomainsWithPlainInstances = []sdk.QuotaSharedResourceDomain{sdk.QuotaSharedResourceDomainAiFunction, sdk.QuotaSharedResourceDomainCortexCode}

// quotaInstanceKey normalizes an instance of a shared resource domain, so that its equivalent spellings hash and compare the same.
func quotaInstanceKey(domainRaw string, instance string) string {
	domain, err := sdk.ToQuotaSharedResourceDomain(domainRaw)
	switch {
	case instance == "", err != nil:
		return instance
	case slices.Contains(quotaDomainsWithPlainInstances, domain):
		return strings.ToUpper(instance)
	case domain == sdk.QuotaSharedResourceDomainCortexAgent:
		if id, err := sdk.ParseSchemaObjectIdentifier(instance); err == nil {
			return id.FullyQualifiedName()
		}
	default:
		if id, err := sdk.ParseAccountObjectIdentifier(instance); err == nil {
			return id.FullyQualifiedName()
		}
	}
	return instance
}

func quotaSharedResourceHash(v any) int {
	resource := v.(map[string]any)
	domain := strings.ToUpper(resource["domain"].(string))
	return schema.HashString(domain + "/" + quotaInstanceKey(domain, resource["instance"].(string)))
}

func suppressQuotaSharedResourceInstance(k, oldValue, newValue string, d *schema.ResourceData) bool {
	domain, _ := d.Get(strings.TrimSuffix(k, "instance") + "domain").(string)
	return quotaInstanceKey(domain, oldValue) == quotaInstanceKey(domain, newValue)
}

func quotaTagValueHash(v any) int {
	tagValue := v.(map[string]any)
	tag := tagValue["tag"].(string)
	if id, err := sdk.ParseSchemaObjectIdentifier(tag); err == nil {
		tag = id.FullyQualifiedName()
	}
	return schema.HashString(tag + "=" + tagValue["value"].(string))
}

func quotaNotificationThresholdHash(v any) int {
	threshold := v.(map[string]any)
	return schema.HashString(fmt.Sprintf("%d/%s/%t/%s",
		threshold["threshold"].(int),
		strings.ToUpper(threshold["spend_strategy"].(string)),
		threshold["notify_user"].(bool),
		strings.ToUpper(threshold["cycle"].(string)),
	))
}

func quotaCustomActionHash(v any) int {
	action := v.(map[string]any)
	procedure := action["procedure"].(string)
	if id, err := sdk.ParseSchemaObjectIdentifierWithArguments(procedure); err == nil {
		procedure = id.FullyQualifiedName()
	}
	arguments, _ := action["arguments"].([]any)
	return schema.HashString(fmt.Sprintf("%s/%q/%s/%d", procedure, expandStringListAllowEmpty(arguments), strings.ToUpper(action["spend_strategy"].(string)), action["threshold"].(int)))
}
