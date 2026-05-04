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
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Specifies the allowed secondary roles for a session policy, if any.",
		Elem: &schema.Resource{
			Schema: sessionPolicySecondaryRolesSchema("allowed_secondary_roles"),
		},
	},
	"blocked_secondary_roles": {
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Description: joinWithSpace("Specifies the blocked secondary roles for a session policy, if any.",
			"Blocked secondary roles take precedence over allowed secondary roles."),
		Elem: &schema.Resource{
			Schema: sessionPolicySecondaryRolesSchema("blocked_secondary_roles"),
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

func sessionPolicySecondaryRolesSchema(attrName string) map[string]*schema.Schema {
	allDescription, noneDescription := "When true, allows all secondary roles.", "When true, disallows all secondary roles."
	action := "allowed"
	if attrName == "blocked_secondary_roles" {
		allDescription, noneDescription = noneDescription, allDescription
		action = "blocked"
	}

	exactlyOneOf := func(attrName string) []string {
		return collections.Map([]string{"none", "all", "roles"}, func(elem string) string {
			return fmt.Sprintf("%s.0.%s", attrName, elem)
		})
	}

	return map[string]*schema.Schema{
		"none": {
			Type:         schema.TypeBool,
			Optional:     true,
			Description:  noneDescription,
			ExactlyOneOf: exactlyOneOf(attrName),
		},
		"all": {
			Type:         schema.TypeBool,
			Optional:     true,
			Description:  allDescription,
			ExactlyOneOf: exactlyOneOf(attrName),
		},
		"roles": {
			Type:         schema.TypeSet,
			Optional:     true,
			Description:  fmt.Sprintf("Specifies roles to be %s as secondary roles.", action),
			ExactlyOneOf: exactlyOneOf(attrName),
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
			},
		},
	}
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
			// For now, the set/list fields have to be excluded.
			// TODO [SNOW-1648997]: address the above comment
			ComputedIfAnyAttributeChanged(sessionPolicySchema, DescribeOutputAttributeName, "name", "schema", "database", "session_idle_timeout_mins", "session_ui_idle_timeout_mins", "comment"),
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
		attributeMappedValueCreateBuilder(d, "allowed_secondary_roles", request.WithAllowedSecondaryRoles, buildSessionPolicySecondaryRolesRequest),
		attributeMappedValueCreateBuilder(d, "blocked_secondary_roles", request.WithBlockedSecondaryRoles, buildSessionPolicySecondaryRolesRequest),
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
				if list, ok := v.([]any); ok {
					return expandStringList(list)
				}
				return v
			}
			if err = handleExternalChangesToObjectInFlatDescribeDeepEqual(d,
				outputMapping{"allowed_secondary_roles", "allowed_secondary_roles", details.AllowedSecondaryRoles, secondaryRolesResourceStateFromDescribeOutput(details.AllowedSecondaryRoles), normalizeStringList},
				outputMapping{"blocked_secondary_roles", "blocked_secondary_roles", details.BlockedSecondaryRoles, secondaryRolesResourceStateFromDescribeOutput(details.BlockedSecondaryRoles), normalizeStringList},
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
		attributeMappedValueUpdate(d, "allowed_secondary_roles", &setRequest.AllowedSecondaryRoles, &unsetRequest.AllowedSecondaryRoles, buildSessionPolicySecondaryRolesRequest),
		attributeMappedValueUpdate(d, "blocked_secondary_roles", &setRequest.BlockedSecondaryRoles, &unsetRequest.BlockedSecondaryRoles, buildSessionPolicySecondaryRolesRequest),
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

func buildSessionPolicySecondaryRolesRequest(v any) (sdk.SessionPolicySecondaryRolesRequest, error) {
	list := v.([]any)
	if len(list) == 0 {
		return sdk.SessionPolicySecondaryRolesRequest{}, fmt.Errorf("session policy secondary roles block is empty")
	}
	block := list[0].(map[string]any)
	request := sdk.NewSessionPolicySecondaryRolesRequest()
	if block["none"].(bool) {
		return *request.WithNone(true), nil
	}
	if block["all"].(bool) {
		return *request.WithAll(true), nil
	}
	rawRoles := block["roles"].(*schema.Set).List()
	roles := expandStringList(rawRoles)
	ids, err := collections.MapErr(roles, sdk.ParseAccountObjectIdentifier)
	if err != nil {
		return sdk.SessionPolicySecondaryRolesRequest{}, err
	}
	return *request.WithRoles(ids), nil
}

// secondaryRolesResourceStateFromDescribeOutput converts DESCRIBE output (list of role names, empty list,
// or a single ALL sentinel) into the nested block shape used by allowed_secondary_roles / blocked_secondary_roles.
func secondaryRolesResourceStateFromDescribeOutput(roles []string) []any {
	if len(roles) == 0 {
		return []any{
			map[string]any{"none": true},
		}
	}
	if len(roles) == 1 && strings.EqualFold(roles[0], "ALL") {
		return []any{
			map[string]any{"all": true},
		}
	}
	elems := make([]any, len(roles))
	for i, r := range roles {
		elems[i] = r
	}
	return []any{
		map[string]any{
			"roles": schema.NewSet(schema.HashString, elems),
		},
	}
}
