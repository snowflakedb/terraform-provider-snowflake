package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var passwordPolicySchema = map[string]*schema.Schema{
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The database this password policy belongs to.",
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The schema this password policy belongs to.",
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Identifier for the password policy; must be unique for your account.",
	},
	"or_replace": {
		Type:                  schema.TypeBool,
		Optional:              true,
		Default:               false,
		Description:           "Whether to override a previous password policy with the same name.",
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
		DiffSuppressOnRefresh: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return old != new
		},
	},
	"min_length": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      8,
		Description:  "Specifies the minimum number of characters the password must contain. Supported range: 8 to 256, inclusive. Default: 8",
		ValidateFunc: validation.IntBetween(8, 256),
	},
	"max_length": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      256,
		Description:  "Specifies the maximum number of characters the password must contain. This number must be greater than or equal to the sum of PASSWORD_MIN_LENGTH, PASSWORD_MIN_UPPER_CASE_CHARS, and PASSWORD_MIN_LOWER_CASE_CHARS. Supported range: 8 to 256, inclusive. Default: 256",
		ValidateFunc: validation.IntBetween(8, 256),
	},
	"min_upper_case_chars": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "Specifies the minimum number of uppercase characters the password must contain. Supported range: 0 to 256, inclusive. Default: 1",
		ValidateFunc: validation.IntBetween(0, 256),
	},
	"min_lower_case_chars": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "Specifies the minimum number of lowercase characters the password must contain. Supported range: 0 to 256, inclusive. Default: 1",
		ValidateFunc: validation.IntBetween(0, 256),
	},
	"min_numeric_chars": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "Specifies the minimum number of numeric characters the password must contain. Supported range: 0 to 256, inclusive. Default: 1",
		ValidateFunc: validation.IntBetween(0, 256),
	},
	"min_special_chars": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "Specifies the minimum number of special characters the password must contain. Supported range: 0 to 256, inclusive. Default: 1",
		ValidateFunc: validation.IntBetween(0, 256),
	},
	"min_age_days": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      0,
		Description:  "Specifies the number of days the user must wait before a recently changed password can be changed again. Supported range: 0 to 999, inclusive. Default: 0",
		ValidateFunc: validation.IntBetween(0, 999),
	},
	"max_age_days": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      90,
		Description:  "Specifies the maximum number of days before the password must be changed. Supported range: 0 to 999, inclusive. A value of zero (i.e. 0) indicates that the password does not need to be changed. Snowflake does not recommend choosing this value for a default account-level password policy or for any user-level policy. Instead, choose a value that meets your internal security guidelines. Default: 90, which means the password must be changed every 90 days.",
		ValidateFunc: validation.IntBetween(0, 999),
	},
	"max_retries": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      5,
		Description:  "Specifies the maximum number of attempts to enter a password before being locked out. Supported range: 1 to 10, inclusive. Default: 5",
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"lockout_time_mins": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      15,
		Description:  "Specifies the number of minutes the user account will be locked after exhausting the designated number of password retries (i.e. PASSWORD_MAX_RETRIES). Supported range: 1 to 999, inclusive. Default: 15",
		ValidateFunc: validation.IntBetween(1, 999),
	},
	"history": {
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      0,
		Description:  "Specifies the number of the most recent passwords that Snowflake stores. These stored passwords cannot be repeated when a user updates their password value. The current password value does not count towards the history. When you increase the history value, Snowflake saves the previous values. When you decrease the value, Snowflake saves the stored values up to that value that is set. For example, if the history value is 8 and you change the history value to 3, Snowflake stores the most recent 3 passwords and deletes the 5 older password values from the history. Default: 0 Max: 24",
		ValidateFunc: validation.IntBetween(0, 24),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Adds a comment or overwrites an existing comment for the password policy.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func PasswordPolicy() *schema.Resource {
	// TODO(SNOW-1818849): unassign policies before dropping
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.SchemaObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.PasswordPolicies.DropSafely
		},
	)

	return &schema.Resource{
		Description:   "A password policy specifies the requirements that must be met to create and reset a password to authenticate to Snowflake.",
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.PasswordPolicyResource), TrackingCreateWrapper(resources.PasswordPolicy, CreatePasswordPolicy)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.PasswordPolicyResource), TrackingReadWrapper(resources.PasswordPolicy, ReadPasswordPolicy)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.PasswordPolicyResource), TrackingUpdateWrapper(resources.PasswordPolicy, UpdatePasswordPolicy)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.PasswordPolicyResource), TrackingDeleteWrapper(resources.PasswordPolicy, deleteFunc)),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.PasswordPolicy, customdiff.All(
			ComputedIfAnyAttributeChanged(passwordPolicySchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Schema: passwordPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreatePasswordPolicy implements schema.CreateFunc.
func CreatePasswordPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	objectIdentifier := sdk.NewSchemaObjectIdentifier(database, schema, name)

	createOptions := &sdk.CreatePasswordPolicyOptions{
		OrReplace:                 sdk.Bool(d.Get("or_replace").(bool)),
		IfNotExists:               sdk.Bool(d.Get("if_not_exists").(bool)),
		PasswordMinLength:         sdk.Int(d.Get("min_length").(int)),
		PasswordMaxLength:         sdk.Int(d.Get("max_length").(int)),
		PasswordMinUpperCaseChars: sdk.Int(d.Get("min_upper_case_chars").(int)),
		PasswordMinLowerCaseChars: sdk.Int(d.Get("min_lower_case_chars").(int)),
		PasswordMinNumericChars:   sdk.Int(d.Get("min_numeric_chars").(int)),
		PasswordMinSpecialChars:   sdk.Int(d.Get("min_special_chars").(int)),
		PasswordMinAgeDays:        sdk.Int(d.Get("min_age_days").(int)),
		PasswordMaxAgeDays:        sdk.Int(d.Get("max_age_days").(int)),
		PasswordMaxRetries:        sdk.Int(d.Get("max_retries").(int)),
		PasswordLockoutTimeMins:   sdk.Int(d.Get("lockout_time_mins").(int)),
		PasswordHistory:           sdk.Int(d.Get("history").(int)),
	}

	if v, ok := d.GetOk("comment"); ok {
		createOptions.Comment = sdk.String(v.(string))
	}

	err := client.PasswordPolicies.Create(ctx, objectIdentifier, createOptions)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeSnowflakeID(objectIdentifier))
	return ReadPasswordPolicy(ctx, d, meta)
}

// ReadPasswordPolicy implements schema.ReadFunc.
func ReadPasswordPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.SchemaObjectIdentifier)

	passwordPolicy, err := client.PasswordPolicies.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query password policy. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Password policy id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database", passwordPolicy.DatabaseName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", passwordPolicy.SchemaName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", passwordPolicy.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("comment", passwordPolicy.Comment); err != nil {
		return diag.FromErr(err)
	}
	passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setFromIntProperty(d, "min_length", passwordPolicyDetails.PasswordMinLength); err != nil {
		return diag.FromErr(err)
	}
	if err := setFromIntProperty(d, "max_length", passwordPolicyDetails.PasswordMaxLength); err != nil {
		return diag.FromErr(err)
	}
	if err := setFromIntProperty(d, "min_upper_case_chars", passwordPolicyDetails.PasswordMinUpperCaseChars); err != nil {
		return diag.FromErr(err)
	}
	if err := setFromIntProperty(d, "min_lower_case_chars", passwordPolicyDetails.PasswordMinLowerCaseChars); err != nil {
		return diag.FromErr(err)
	}
	if err := setFromIntProperty(d, "min_numeric_chars", passwordPolicyDetails.PasswordMinNumericChars); err != nil {
		return diag.FromErr(err)
	}
	if err := setFromIntProperty(d, "min_special_chars", passwordPolicyDetails.PasswordMinSpecialChars); err != nil {
		return diag.FromErr(err)
	}
	if err := setFromIntProperty(d, "min_age_days", passwordPolicyDetails.PasswordMinAgeDays); err != nil {
		return diag.FromErr(err)
	}
	if err := setFromIntProperty(d, "max_age_days", passwordPolicyDetails.PasswordMaxAgeDays); err != nil {
		return diag.FromErr(err)
	}
	if err := setFromIntProperty(d, "max_retries", passwordPolicyDetails.PasswordMaxRetries); err != nil {
		return diag.FromErr(err)
	}
	if err := setFromIntProperty(d, "lockout_time_mins", passwordPolicyDetails.PasswordLockoutTimeMins); err != nil {
		return diag.FromErr(err)
	}
	if err := setFromIntProperty(d, "history", passwordPolicyDetails.PasswordHistory); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// UpdatePasswordPolicy implements schema.UpdateFunc.
func UpdatePasswordPolicy(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	objectIdentifier := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.SchemaObjectIdentifier)

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(objectIdentifier.SchemaId(), d.Get("name").(string))

		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, &sdk.AlterPasswordPolicyOptions{
			NewName: &newId,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeSnowflakeID(newId))
		objectIdentifier = newId
	}

	if d.HasChange("min_length") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordMinLength: sdk.Int(d.Get("min_length").(int)),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("max_length") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordMaxLength: sdk.Int(d.Get("max_length").(int)),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("min_upper_case_chars") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordMinUpperCaseChars: sdk.Int(d.Get("min_upper_case_chars").(int)),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("min_lower_case_chars") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordMinLowerCaseChars: sdk.Int(d.Get("min_lower_case_chars").(int)),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("min_numeric_chars") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordMinNumericChars: sdk.Int(d.Get("min_numeric_chars").(int)),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("min_special_chars") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordMinSpecialChars: sdk.Int(d.Get("min_special_chars").(int)),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("min_age_days") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordMinAgeDays: sdk.Int(d.Get("min_age_days").(int)),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("max_age_days") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordMaxAgeDays: sdk.Int(d.Get("max_age_days").(int)),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("max_retries") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordMaxRetries: sdk.Int(d.Get("max_retries").(int)),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("lockout_time_mins") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordLockoutTimeMins: sdk.Int(d.Get("lockout_time_mins").(int)),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("history") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordHistory: sdk.Int(d.Get("history").(int)),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("comment") {
		alterOptions := &sdk.AlterPasswordPolicyOptions{}
		if v, ok := d.GetOk("comment"); ok {
			alterOptions.Set = &sdk.PasswordPolicySet{
				Comment: sdk.String(v.(string)),
			}
		} else {
			alterOptions.Unset = &sdk.PasswordPolicyUnset{
				Comment: sdk.Bool(true),
			}
		}
		err := client.PasswordPolicies.Alter(ctx, objectIdentifier, alterOptions)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadPasswordPolicy(ctx, d, meta)
}
