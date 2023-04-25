package resources

import (
	"context"
	"database/sql"
	"errors"

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
		ForceNew:    true,
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
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Adds a comment or overwrites an existing comment for the password policy.",
	},
}

func PasswordPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "A password policy specifies the requirements that must be met to create and reset a password to authenticate to Snowflake.",
		Create:      CreatePasswordPolicy,
		Read:        ReadPasswordPolicy,
		Update:      UpdatePasswordPolicy,
		Delete:      DeletePasswordPolicy,

		Schema: passwordPolicySchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreatePasswordPolicy implements schema.CreateFunc.
func CreatePasswordPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)
	objectIdentifier := sdk.SchemaObjectIdentifier{
		DatabaseName: database,
		SchemaName:   schema,
		Name:         name,
	}
	fqn := objectIdentifier.FullyQualifiedName()
	err := client.PasswordPolicies.Create(ctx, fqn, &sdk.PasswordPolicyCreateOptions{
		OrReplace:                 sdk.Bool(d.Get("or_replace").(bool)),
		IfNotExists:               sdk.Bool(d.Get("if_not_exists").(bool)),
		PasswordMinLength:         sdk.Int(d.Get("min_length").(int)),
		PasswordMaxLength:         sdk.Int(d.Get("max_length").(int)),
		PasswordMinUpperCaseChars: sdk.Int(d.Get("min_upper_case_chars").(int)),
		PasswordMinLowerCaseChars: sdk.Int(d.Get("min_lower_case_chars").(int)),
		PasswordMinNumericChars:   sdk.Int(d.Get("min_numeric_chars").(int)),
		PasswordMinSpecialChars:   sdk.Int(d.Get("min_special_chars").(int)),
		PasswordMaxAgeDays:        sdk.Int(d.Get("max_age_days").(int)),
		PasswordMaxRetries:        sdk.Int(d.Get("max_retries").(int)),
		PasswordLockoutTimeMins:   sdk.Int(d.Get("lockout_time_mins").(int)),
		Comment:                   sdk.String(d.Get("comment").(string)),
	})
	if err != nil {
		return err
	}
	id := helpers.SnowflakeID(database, schema, name)
	d.SetId(id)
	return ReadPasswordPolicy(d, meta)
}

// ReadPasswordPolicy implements schema.ReadFunc.
func ReadPasswordPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	schemaIdentifier := sdk.NewSchemaIdentifier(objectIdentifier.DatabaseName, objectIdentifier.SchemaName)
	fqn := objectIdentifier.FullyQualifiedName()
	passwordPolicyList, err := client.PasswordPolicies.Show(ctx, &sdk.PasswordPolicyShowOptions{
		Pattern: sdk.String(objectIdentifier.Name),
		In: &sdk.PasswordPolicyShowIn{
			Schema: sdk.String(schemaIdentifier.FullyQualifiedName()),
		},
	})
	if err != nil {
		if errors.Is(err, sdk.ErrNoRecord) {
			d.SetId("")
			return nil
		}
		return err
	}
	passwordPolicy := passwordPolicyList[0]

	if err := d.Set("database", passwordPolicy.DatabaseName); err != nil {
		return err
	}
	if err := d.Set("schema", passwordPolicy.SchemaName); err != nil {
		return err
	}
	if err := d.Set("name", passwordPolicy.Name); err != nil {
		return err
	}
	if err := d.Set("comment", passwordPolicy.Comment); err != nil {
		return err
	}
	passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, fqn)
	if err != nil {
		return err
	}
	if err := d.Set("min_length", passwordPolicyDetails.PasswordMinLength); err != nil {
		return err
	}
	if err := d.Set("max_length", passwordPolicyDetails.PasswordMaxLength); err != nil {
		return err
	}
	if err := d.Set("min_upper_case_chars", passwordPolicyDetails.PasswordMinUpperCaseChars); err != nil {
		return err
	}
	if err := d.Set("min_lower_case_chars", passwordPolicyDetails.PasswordMinLowerCaseChars); err != nil {
		return err
	}
	if err := d.Set("min_numeric_chars", passwordPolicyDetails.PasswordMinNumericChars); err != nil {
		return err
	}
	if err := d.Set("min_special_chars", passwordPolicyDetails.PasswordMinSpecialChars); err != nil {
		return err
	}
	if err := d.Set("max_age_days", passwordPolicyDetails.PasswordMaxAgeDays); err != nil {
		return err
	}
	if err := d.Set("max_retries", passwordPolicyDetails.PasswordMaxRetries); err != nil {
		return err
	}
	if err := d.Set("lockout_time_mins", passwordPolicyDetails.PasswordLockoutTimeMins); err != nil {
		return err
	}

	return nil
}

// UpdatePasswordPolicy implements schema.UpdateFunc.
func UpdatePasswordPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()

	alterOptions := &sdk.PasswordPolicyAlterOptions{
		Set: &sdk.PasswordPolicyAlterSet{},
	}

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		alterOptions.Set.Name = sdk.String(n.(string))
	}

	if d.HasChange("min_length") {
		alterOptions.Set.PasswordMinLength = sdk.Int(d.Get("min_length").(int))
	}
	if d.HasChange("max_length") {
		alterOptions.Set.PasswordMaxLength = sdk.Int(d.Get("max_length").(int))
	}
	if d.HasChange("min_upper_case_chars") {
		alterOptions.Set.PasswordMinUpperCaseChars = sdk.Int(d.Get("min_upper_case_chars").(int))
	}
	if d.HasChange("min_lower_case_chars") {
		alterOptions.Set.PasswordMinLowerCaseChars = sdk.Int(d.Get("min_lower_case_chars").(int))
	}
	if d.HasChange("min_numeric_chars") {
		alterOptions.Set.PasswordMinNumericChars = sdk.Int(d.Get("min_numeric_chars").(int))
	}
	if d.HasChange("min_special_chars") {
		alterOptions.Set.PasswordMinSpecialChars = sdk.Int(d.Get("min_special_chars").(int))
	}
	if d.HasChange("max_age_days") {
		alterOptions.Set.PasswordMaxAgeDays = sdk.Int(d.Get("max_age_days").(int))
	}
	if d.HasChange("max_retries") {
		alterOptions.Set.PasswordMaxRetries = sdk.Int(d.Get("max_retries").(int))
	}
	if d.HasChange("lockout_time_mins") {
		alterOptions.Set.PasswordLockoutTimeMins = sdk.Int(d.Get("lockout_time_mins").(int))
	}
	if d.HasChange("comment") {
		alterOptions.Set.Comment = sdk.String(d.Get("comment").(string))
	}
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	name := objectIdentifier.FullyQualifiedName()
	err := client.PasswordPolicies.Alter(ctx, name, alterOptions)
	if err != nil {
		return err
	}
	return nil
}

// DeletePasswordPolicy implements schema.DeleteFunc.
func DeletePasswordPolicy(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	objectIdentifier := helpers.DecodeSnowflakeID(d.Id()).(sdk.SchemaObjectIdentifier)
	name := objectIdentifier.FullyQualifiedName()
	err := client.PasswordPolicies.Drop(ctx, name, nil)
	if err != nil {
		return err
	}

	return nil
}
