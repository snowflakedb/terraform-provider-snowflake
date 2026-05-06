package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

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

var passwordPolicySchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Identifier for the password policy; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The schema this password policy belongs to."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The database this password policy belongs to."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"or_replace": {
		Type:                  schema.TypeBool,
		Optional:              true,
		Default:               false,
		Description:           "Whether to override a previous password policy with the same name.",
		Deprecated:            "This field is a noop and will be removed in a future version of the provider.",
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return old != new
		},
	},
	"if_not_exists": {
		Type:                  schema.TypeBool,
		Optional:              true,
		Default:               false,
		Description:           "Prevent overwriting a previous password policy with the same name.",
		Deprecated:            "This field is a noop and will be removed in a future version of the provider.",
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return old != new
		},
	},
	"min_length": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the minimum number of characters the password must contain.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("password_min_length"),
	},
	"max_length": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the maximum number of characters the password must contain. This number must be greater than or equal to the sum of PASSWORD_MIN_LENGTH, PASSWORD_MIN_UPPER_CASE_CHARS, and PASSWORD_MIN_LOWER_CASE_CHARS.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("password_max_length"),
	},
	"min_upper_case_chars": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Specifies the minimum number of uppercase characters the password must contain.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("password_min_upper_case_chars"),
	},
	"min_lower_case_chars": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Specifies the minimum number of lowercase characters the password must contain.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("password_min_lower_case_chars"),
	},
	"min_numeric_chars": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Specifies the minimum number of numeric characters the password must contain.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("password_min_numeric_chars"),
	},
	"min_special_chars": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Specifies the minimum number of special characters the password must contain.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("password_min_special_chars"),
	},
	"min_age_days": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Specifies the number of days the user must wait before a recently changed password can be changed again.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("password_min_age_days"),
	},
	"max_age_days": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Specifies the maximum number of days before the password must be changed. A value of zero (i.e. 0) indicates that the password does not need to be changed.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("password_max_age_days"),
	},
	"max_retries": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the maximum number of attempts to enter a password before being locked out.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("password_max_retries"),
	},
	"lockout_time_mins": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the number of minutes the user account will be locked after exhausting the designated number of password retries (i.e. PASSWORD_MAX_RETRIES).",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("password_lockout_time_mins"),
	},
	"history": {
		Type:             schema.TypeInt,
		Optional:         true,
		Default:          IntDefault,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Specifies the number of the most recent passwords that Snowflake stores. These stored passwords cannot be repeated when a user updates their password value. The current password value does not count towards the history. When you increase the history value, Snowflake saves the previous values. When you decrease the value, Snowflake saves the stored values up to that value that is set. For example, if the history value is 8 and you change the history value to 3, Snowflake stores the most recent 3 passwords and deletes the 5 older password values from the history.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("password_history"),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Adds a comment or overwrites an existing comment for the password policy.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PASSWORD POLICIES` for the given password policy.",
		Elem: &schema.Resource{
			Schema: schemas.ShowPasswordPolicySchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE PASSWORD POLICY` for the given password policy.",
		Elem: &schema.Resource{
			Schema: schemas.DescribePasswordPolicyDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func PasswordPolicy() *schema.Resource {
	// TODO(SNOW-1818849): unassign policies before dropping
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.PasswordPolicies.DropSafely
		},
	)

	return &schema.Resource{
		Description:   "A password policy specifies the requirements that must be met to create and reset a password to authenticate to Snowflake.",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.PasswordPolicyResource), TrackingCreateWrapper(resources.PasswordPolicy, CreatePasswordPolicy)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.PasswordPolicyResource), TrackingReadWrapper(resources.PasswordPolicy, ReadPasswordPolicyFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.PasswordPolicyResource), TrackingUpdateWrapper(resources.PasswordPolicy, UpdatePasswordPolicy)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.PasswordPolicyResource), TrackingDeleteWrapper(resources.PasswordPolicy, deleteFunc)),

		Schema: passwordPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.PasswordPolicy, ImportName[sdk.SchemaObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    cty.EmptyObject,
				Upgrade: v2_15_0_PasswordPolicyStateUpgrader,
			},
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.PasswordPolicy, customdiff.All(
			ComputedIfAnyAttributeChanged(passwordPolicySchema, ShowOutputAttributeName, "name", "database", "schema", "comment"),
			ComputedIfAnyAttributeChanged(passwordPolicySchema, DescribeOutputAttributeName, "name", "comment",
				"min_length", "max_length", "min_upper_case_chars", "min_lower_case_chars",
				"min_numeric_chars", "min_special_chars", "min_age_days", "max_age_days",
				"max_retries", "lockout_time_mins", "history"),
			ComputedIfAnyAttributeChanged(passwordPolicySchema, FullyQualifiedNameAttributeName, "name", "database", "schema"),
		)),
	}
}

func CreatePasswordPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	req := sdk.NewCreatePasswordPolicyRequest(id)

	errs := errors.Join(
		intAttributeCreateBuilder(d, "min_length", req.WithPasswordMinLength),
		intAttributeCreateBuilder(d, "max_length", req.WithPasswordMaxLength),
		intAttributeCreateBuilder(d, "max_retries", req.WithPasswordMaxRetries),
		intAttributeCreateBuilder(d, "lockout_time_mins", req.WithPasswordLockoutTimeMins),
		intAttributeWithSpecialDefaultCreateBuilder(d, "min_upper_case_chars", req.WithPasswordMinUpperCaseChars),
		intAttributeWithSpecialDefaultCreateBuilder(d, "min_lower_case_chars", req.WithPasswordMinLowerCaseChars),
		intAttributeWithSpecialDefaultCreateBuilder(d, "min_numeric_chars", req.WithPasswordMinNumericChars),
		intAttributeWithSpecialDefaultCreateBuilder(d, "min_special_chars", req.WithPasswordMinSpecialChars),
		intAttributeWithSpecialDefaultCreateBuilder(d, "min_age_days", req.WithPasswordMinAgeDays),
		intAttributeWithSpecialDefaultCreateBuilder(d, "max_age_days", req.WithPasswordMaxAgeDays),
		intAttributeWithSpecialDefaultCreateBuilder(d, "history", req.WithPasswordHistory),
		stringAttributeCreateBuilder(d, "comment", req.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.PasswordPolicies.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadPasswordPolicyFunc(false)(ctx, d, meta)
}

func ReadPasswordPolicyFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		passwordPolicy, err := client.PasswordPolicies.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{{
					Severity: diag.Warning,
					Summary:  "Failed to query password policy. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Password policy id: %s, Err: %s", id.FullyQualifiedName(), err),
				}}
			}
			return diag.FromErr(err)
		}

		details, err := client.PasswordPolicies.DescribeDetails(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("could not describe password policy (%s), err = %w", id.FullyQualifiedName(), err))
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"password_min_length", "min_length", details.PasswordMinLength, details.PasswordMinLength, nil},
				outputMapping{"password_max_length", "max_length", details.PasswordMaxLength, details.PasswordMaxLength, nil},
				outputMapping{"password_min_upper_case_chars", "min_upper_case_chars", details.PasswordMinUpperCaseChars, details.PasswordMinUpperCaseChars, nil},
				outputMapping{"password_min_lower_case_chars", "min_lower_case_chars", details.PasswordMinLowerCaseChars, details.PasswordMinLowerCaseChars, nil},
				outputMapping{"password_min_numeric_chars", "min_numeric_chars", details.PasswordMinNumericChars, details.PasswordMinNumericChars, nil},
				outputMapping{"password_min_special_chars", "min_special_chars", details.PasswordMinSpecialChars, details.PasswordMinSpecialChars, nil},
				outputMapping{"password_min_age_days", "min_age_days", details.PasswordMinAgeDays, details.PasswordMinAgeDays, nil},
				outputMapping{"password_max_age_days", "max_age_days", details.PasswordMaxAgeDays, details.PasswordMaxAgeDays, nil},
				outputMapping{"password_max_retries", "max_retries", details.PasswordMaxRetries, details.PasswordMaxRetries, nil},
				outputMapping{"password_lockout_time_mins", "lockout_time_mins", details.PasswordLockoutTimeMins, details.PasswordLockoutTimeMins, nil},
				outputMapping{"password_history", "history", details.PasswordHistory, details.PasswordHistory, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		errs := errors.Join(
			// Not reading int fields here — they are handled by handleExternalChangesToObjectInFlatDescribe
			d.Set("comment", passwordPolicy.Comment),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.PasswordPolicyToSchema(passwordPolicy)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.PasswordPolicyDetailsToSchema(details)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		)

		return diag.FromErr(errs)
	}
}

func UpdatePasswordPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("database") || d.HasChange("schema") || d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifier(d.Get("database").(string), d.Get("schema").(string), d.Get("name").(string))
		if err := client.PasswordPolicies.Alter(ctx, sdk.NewAlterPasswordPolicyRequest(id).WithNewName(newId)); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming password policy from %v to %v, err = %w", id.FullyQualifiedName(), newId.FullyQualifiedName(), err))
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set := sdk.NewPasswordPolicySetRequest()
	unset := sdk.NewPasswordPolicyUnsetRequest()

	if err := errors.Join(
		intAttributeUpdate(d, "min_length", &set.PasswordMinLength, &unset.PasswordMinLength),
		intAttributeUpdate(d, "max_length", &set.PasswordMaxLength, &unset.PasswordMaxLength),
		intAttributeUpdate(d, "max_retries", &set.PasswordMaxRetries, &unset.PasswordMaxRetries),
		intAttributeUpdate(d, "lockout_time_mins", &set.PasswordLockoutTimeMins, &unset.PasswordLockoutTimeMins),
		intAttributeWithSpecialDefaultUpdate(d, "min_upper_case_chars", &set.PasswordMinUpperCaseChars, &unset.PasswordMinUpperCaseChars),
		intAttributeWithSpecialDefaultUpdate(d, "min_lower_case_chars", &set.PasswordMinLowerCaseChars, &unset.PasswordMinLowerCaseChars),
		intAttributeWithSpecialDefaultUpdate(d, "min_numeric_chars", &set.PasswordMinNumericChars, &unset.PasswordMinNumericChars),
		intAttributeWithSpecialDefaultUpdate(d, "min_special_chars", &set.PasswordMinSpecialChars, &unset.PasswordMinSpecialChars),
		intAttributeWithSpecialDefaultUpdate(d, "min_age_days", &set.PasswordMinAgeDays, &unset.PasswordMinAgeDays),
		intAttributeWithSpecialDefaultUpdate(d, "max_age_days", &set.PasswordMaxAgeDays, &unset.PasswordMaxAgeDays),
		intAttributeWithSpecialDefaultUpdate(d, "history", &set.PasswordHistory, &unset.PasswordHistory),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	); err != nil {
		return diag.FromErr(err)
	}

	if !reflect.DeepEqual(*set, *sdk.NewPasswordPolicySetRequest()) {
		if err := client.PasswordPolicies.Alter(ctx, sdk.NewAlterPasswordPolicyRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if !reflect.DeepEqual(*unset, *sdk.NewPasswordPolicyUnsetRequest()) {
		if err := client.PasswordPolicies.Alter(ctx, sdk.NewAlterPasswordPolicyRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadPasswordPolicyFunc(false)(ctx, d, meta)
}
