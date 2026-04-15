package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var sessionPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the session policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the session policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the session policy."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"session_idle_timeout_mins": {
		Type:             schema.TypeInt,
		Optional:         true,
		Description:      "For Snowflake clients and programmatic clients, specifies the number of minutes in which a session can be idle before users must authenticate to Snowflake again.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("session_idle_timeout_mins"),
		ValidateFunc:     validation.IntAtLeast(1),
	},
	"session_ui_idle_timeout_mins": {
		Type:             schema.TypeInt,
		Optional:         true,
		Description:      "For Snowsight, specifies the number of minutes in which a session can be idle before users must authenticate to Snowflake again.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("session_ui_idle_timeout_mins"),
		ValidateFunc:     validation.IntAtLeast(1),
	},
	"allowed_secondary_roles": {
		Type:     schema.TypeSet,
		Optional: true,
		Description: joinWithSpace("Specifies the allowed secondary roles for a session policy, if any.",
			"Use a single-element set whose only value is `all` (case-insensitive) to allow all secondary roles, equivalent to `('ALL')` in Snowflake."),
		DiffSuppressFunc: SuppressIfAny(
			IgnoreChangeToCurrentSnowflakeValueInDescribe("allowed_secondary_roles"),
			NormalizeAndCompareIdentifiersInSet("allowed_secondary_roles"),
		),
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
	},
	"blocked_secondary_roles": {
		Type:     schema.TypeSet,
		Optional: true,
		Description: joinWithSpace("Specifies the blocked secondary roles for a session policy, if any.",
			"Blocked secondary roles take precedence over allowed secondary roles.",
			"Use a single-element set whose only value is `all` (case-insensitive) to disallow all secondary roles, equivalent to `('ALL')` in Snowflake.",
		),
		DiffSuppressFunc: SuppressIfAny(
			IgnoreChangeToCurrentSnowflakeValueInDescribe("blocked_secondary_roles"),
			NormalizeAndCompareIdentifiersInSet("blocked_secondary_roles"),
		),
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the session policy.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW SESSION POLICIES` for this session policy.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSessionPolicySchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE SESSION POLICY` for this session policy.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeSessionPolicyDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func SessionPolicy() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.SessionPolicies.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.SessionPolicyResource), TrackingCreateWrapper(resources.SessionPolicy, CreateSessionPolicy)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.SessionPolicyResource), TrackingReadWrapper(resources.SessionPolicy, ReadSessionPolicyFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.SessionPolicyResource), TrackingUpdateWrapper(resources.SessionPolicy, UpdateSessionPolicy)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.SessionPolicyResource), TrackingDeleteWrapper(resources.SessionPolicy, deleteFunc)),
		Description:   "Resource used to manage session policy objects. For more information, check [session policy documentation](https://docs.snowflake.com/en/sql-reference/sql/create-session-policy).",

		Schema: sessionPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.SessionPolicy, ImportName[sdk.SchemaObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: TrackingCustomDiffWrapper(resources.SessionPolicy, customdiff.All(
			ComputedIfAnyAttributeChanged(sessionPolicySchema, ShowOutputAttributeName, "name", "schema", "database", "comment"),
			ComputedIfAnyAttributeChanged(sessionPolicySchema, DescribeOutputAttributeName, "name", "schema", "database", "session_idle_timeout_mins", "session_ui_idle_timeout_mins", "allowed_secondary_roles", "blocked_secondary_roles", "comment"),
			ComputedIfAnyAttributeChanged(sessionPolicySchema, FullyQualifiedNameAttributeName, "name", "schema", "database"),
		)),
	}
}

func CreateSessionPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schemaName := d.Get("schema").(string)
	databaseName := d.Get("database").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	request := sdk.NewCreateSessionPolicyRequest(id)
	errs := errors.Join(
		intAttributeCreateBuilder(d, "session_idle_timeout_mins", request.WithSessionIdleTimeoutMins),
		intAttributeCreateBuilder(d, "session_ui_idle_timeout_mins", request.WithSessionUiIdleTimeoutMins),
		secondaryRolesAttributeCreateBuilder(d, "allowed_secondary_roles", request.WithAllowedSecondaryRoles),
		secondaryRolesAttributeCreateBuilder(d, "blocked_secondary_roles", request.WithBlockedSecondaryRoles),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.SessionPolicies.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating session policy, err = %w", err))
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadSessionPolicyFunc(false)(ctx, d, meta)
}

func ReadSessionPolicyFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		s, err := client.SessionPolicies.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{{
					Severity: diag.Warning,
					Summary:  "Failed to query session policy. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Session policy id: %s, Err: %s", id.FullyQualifiedName(), err),
				}}
			}
			return diag.FromErr(err)
		}

		details, err := client.SessionPolicies.DescribeDetails(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not describe session policy (%s), err = %w", id.FullyQualifiedName(), err))
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"session_idle_timeout_mins", "session_idle_timeout_mins", details.SessionIdleTimeoutMins, details.SessionIdleTimeoutMins, nil},
				outputMapping{"session_ui_idle_timeout_mins", "session_ui_idle_timeout_mins", details.SessionUiIdleTimeoutMins, details.SessionUiIdleTimeoutMins, nil},
			); err != nil {
				return diag.FromErr(err)
			}
			normalizeStringList := func(v any) any {
				if list, ok := v.([]interface{}); ok {
					return expandStringList(list)
				}
				return v
			}
			if err = handleExternalChangesToObjectInFlatDescribeDeepEqual(d,
				outputMapping{"allowed_secondary_roles", "allowed_secondary_roles", details.AllowedSecondaryRoles, details.AllowedSecondaryRoles, normalizeStringList},
				outputMapping{"blocked_secondary_roles", "blocked_secondary_roles", details.BlockedSecondaryRoles, details.BlockedSecondaryRoles, normalizeStringList},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			// not reading session_idle_timeout_mins, session_ui_idle_timeout_mins, allowed_secondary_roles, and blocked_secondary_roles on purpose
			// (handled as external change to describe output)
			d.Set("comment", details.Comment),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.SessionPolicyToSchema(s)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.SessionPolicyDetailsToSchema(details)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)

		return diag.FromErr(errs)
	}
}

func UpdateSessionPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("database") || d.HasChange("schema") || d.HasChange("name") {
		newID := sdk.NewSchemaObjectIdentifier(d.Get("database").(string), d.Get("schema").(string), d.Get("name").(string))
		if err := client.SessionPolicies.Alter(ctx, sdk.NewAlterSessionPolicyRequest(id).WithRenameTo(newID)); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming session policy from %v to %v, err = %w", id.FullyQualifiedName(), newID.FullyQualifiedName(), err))
		}
		d.SetId(helpers.EncodeResourceIdentifier(newID))
		id = newID
	}

	setRequest := sdk.NewSessionPolicySetRequest()
	unsetRequest := sdk.NewSessionPolicyUnsetRequest()

	if err := errors.Join(
		intAttributeUpdate(d, "session_idle_timeout_mins", &setRequest.SessionIdleTimeoutMins, &unsetRequest.SessionIdleTimeoutMins),
		intAttributeUpdate(d, "session_ui_idle_timeout_mins", &setRequest.SessionUiIdleTimeoutMins, &unsetRequest.SessionUiIdleTimeoutMins),
		secondaryRolesAttributeUpdate(d, "allowed_secondary_roles", &setRequest.AllowedSecondaryRoles, &unsetRequest.AllowedSecondaryRoles),
		secondaryRolesAttributeUpdate(d, "blocked_secondary_roles", &setRequest.BlockedSecondaryRoles, &unsetRequest.BlockedSecondaryRoles),
		stringAttributeUpdate(d, "comment", &setRequest.Comment, &unsetRequest.Comment),
	); err != nil {
		return diag.FromErr(err)
	}

	if !reflect.DeepEqual(*setRequest, *sdk.NewSessionPolicySetRequest()) {
		if err := client.SessionPolicies.Alter(ctx, sdk.NewAlterSessionPolicyRequest(id).WithSet(*setRequest)); err != nil {
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unsetRequest, *sdk.NewSessionPolicyUnsetRequest()) {
		if err := client.SessionPolicies.Alter(ctx, sdk.NewAlterSessionPolicyRequest(id).WithUnset(*unsetRequest)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadSessionPolicyFunc(false)(ctx, d, meta)
}

func secondaryRolesAttributeCreateBuilder[T any](d *schema.ResourceData, key string, setValue func(request sdk.SessionPolicySecondaryRolesRequest) *T) error {
	if v, ok := d.GetOk(key); ok {
		secondaryRolesRaw := v.(*schema.Set).List()
		request, err := buildSessionPolicySecondaryRolesRequest(secondaryRolesRaw)
		if err != nil {
			return err
		}
		setValue(*request)
	}
	return nil
}

func secondaryRolesAttributeUpdate(d *schema.ResourceData, key string, setField **sdk.SessionPolicySecondaryRolesRequest, unsetField **bool) error {
	if d.HasChange(key) {
		if v, ok := d.GetOk(key); ok {
			request, err := buildSessionPolicySecondaryRolesRequest(v.(*schema.Set).List())
			if err != nil {
				return err
			}
			*setField = request
		} else {
			*unsetField = sdk.Bool(true)
		}
	}
	return nil
}

func buildSessionPolicySecondaryRolesRequest(rawRoles []any) (*sdk.SessionPolicySecondaryRolesRequest, error) {
	roles := expandStringList(rawRoles)
	request := sdk.NewSessionPolicySecondaryRolesRequest()
	switch {
	case len(roles) == 0:
		request.WithNone(true)
	case len(roles) == 1 && strings.EqualFold("ALL", roles[0]):
		request.WithAll(true)
	default:
		mapped, err := collections.MapErr(roles, sdk.ParseAccountObjectIdentifier)
		if err != nil {
			return nil, err
		}
		request.WithRoles(mapped)
	}
	return request, nil
}
