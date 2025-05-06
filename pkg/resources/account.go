package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var accountSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier (i.e. name) for the account. It must be unique within an organization, regardless of which Snowflake Region the account is in and must start with an alphabetic character and cannot contain spaces or special characters except for underscores (_). Note that if the account name includes underscores, features that do not accept account names with underscores (e.g. Okta SSO or SCIM) can reference a version of the account name that substitutes hyphens (-) for the underscores.",
	},
	"admin_name": {
		Type:             schema.TypeString,
		Required:         true,
		Sensitive:        true,
		Description:      externalChangesNotDetectedFieldDescription("Login name of the initial administrative user of the account. A new user is created in the new account with this name and password and granted the ACCOUNTADMIN role in the account. A login name can be any string consisting of letters, numbers, and underscores. Login names are always case-insensitive."),
		DiffSuppressFunc: IgnoreAfterCreation,
	},
	"admin_password": {
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		Description:      externalChangesNotDetectedFieldDescription("Password for the initial administrative user of the account. Either admin_password or admin_rsa_public_key has to be specified. This field cannot be used whenever admin_user_type is set to SERVICE."),
		DiffSuppressFunc: IgnoreAfterCreation,
		AtLeastOneOf:     []string{"admin_password", "admin_rsa_public_key"},
	},
	"admin_rsa_public_key": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      externalChangesNotDetectedFieldDescription("Assigns a public key to the initial administrative user of the account. Either admin_password or admin_rsa_public_key has to be specified."),
		DiffSuppressFunc: IgnoreAfterCreation,
		AtLeastOneOf:     []string{"admin_password", "admin_rsa_public_key"},
	},
	"admin_user_type": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      externalChangesNotDetectedFieldDescription(fmt.Sprintf("Used for setting the type of the first user that is assigned the ACCOUNTADMIN role during account creation. Valid options are: %s", docs.PossibleValuesListed(sdk.AllUserTypes))),
		DiffSuppressFunc: SuppressIfAny(IgnoreAfterCreation, NormalizeAndCompare(sdk.ToUserType)),
		ValidateDiagFunc: sdkValidation(sdk.ToUserType),
	},
	"first_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		Description:      externalChangesNotDetectedFieldDescription("First name of the initial administrative user of the account. This field cannot be used whenever admin_user_type is set to SERVICE."),
		DiffSuppressFunc: IgnoreAfterCreation,
	},
	"last_name": {
		Type:             schema.TypeString,
		Optional:         true,
		Sensitive:        true,
		Description:      externalChangesNotDetectedFieldDescription("Last name of the initial administrative user of the account. This field cannot be used whenever admin_user_type is set to SERVICE."),
		DiffSuppressFunc: IgnoreAfterCreation,
	},
	"email": {
		Type:             schema.TypeString,
		Required:         true,
		Sensitive:        true,
		Description:      externalChangesNotDetectedFieldDescription("Email address of the initial administrative user of the account. This email address is used to send any notifications about the account."),
		DiffSuppressFunc: IgnoreAfterCreation,
	},
	"must_change_password": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		Description:      externalChangesNotDetectedFieldDescription("Specifies whether the new user created to administer the account is forced to change their password upon first login into the account. This field cannot be used whenever admin_user_type is set to SERVICE."),
		DiffSuppressFunc: IgnoreAfterCreation,
		ValidateDiagFunc: validateBooleanString,
	},
	"edition": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      fmt.Sprintf("Snowflake Edition of the account. See more about Snowflake Editions in the [official documentation](https://docs.snowflake.com/en/user-guide/intro-editions). Valid options are: %s", docs.PossibleValuesListed(sdk.AllAccountEditions)),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToAccountEdition),
		ValidateDiagFunc: sdkValidation(sdk.ToAccountEdition),
	},
	"region_group": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "ID of the region group where the account is created. To retrieve the region group ID for existing accounts in your organization, execute the [SHOW REGIONS](https://docs.snowflake.com/en/sql-reference/sql/show-regions) command. For information about when you might need to specify region group, see [Region groups](https://docs.snowflake.com/en/user-guide/admin-account-identifier.html#label-region-groups).",
	},
	"region": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "[Snowflake Region ID](https://docs.snowflake.com/en/user-guide/admin-account-identifier.html#label-snowflake-region-ids) of the region where the account is created. If no value is provided, Snowflake creates the account in the same Snowflake Region as the current account (i.e. the account in which the CREATE ACCOUNT statement is executed.)",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies a comment for the account.",
		DiffSuppressFunc: SuppressIfAny(
			IgnoreChangeToCurrentSnowflakeValueInShow("comment"),
			func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				return oldValue == "SNOWFLAKE" && newValue == ""
			},
		),
	},
	"is_org_admin": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          BooleanDefault,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("is_org_admin"),
		ValidateDiagFunc: validateBooleanString,
		Description:      "Sets an account property that determines whether the ORGADMIN role is enabled in the account. Only an organization administrator (i.e. user with the ORGADMIN role) can set the property.",
	},
	"grace_period_in_days": {
		Type:             schema.TypeInt,
		Required:         true,
		Description:      "Specifies the number of days during which the account can be restored (“undropped”). The minimum is 3 days and the maximum is 90 days.",
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(3)),
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW ACCOUNTS` for the given account.",
		Elem: &schema.Resource{
			Schema: schemas.ShowAccountSchema,
		},
	},
}

func Account() *schema.Resource {
	return &schema.Resource{
		Description:   "The account resource allows you to create and manage Snowflake accounts. For more information, check [account documentation](https://docs.snowflake.com/en/user-guide/organizations-manage-accounts).",
		CreateContext: TrackingCreateWrapper(resources.Account, CreateAccount),
		ReadContext:   TrackingReadWrapper(resources.Account, ReadAccount(true)),
		UpdateContext: TrackingUpdateWrapper(resources.Account, UpdateAccount),
		DeleteContext: TrackingDeleteWrapper(resources.Account, DeleteAccount),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Account, customdiff.All(
			ComputedIfAnyAttributeChanged(accountSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(accountSchema, ShowOutputAttributeName, "name", "is_org_admin"),
		)),

		Schema: accountSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Account, ImportAccount),
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v0_99_0_AccountStateUpgrader,
			},
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportAccount(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client

	isOrgAdmin, err := client.ContextFunctions.IsRoleInSession(ctx, snowflakeroles.Orgadmin)
	if err != nil {
		return nil, err
	}
	if !isOrgAdmin {
		return nil, errors.New("current user does not have the ORGADMIN role in session")
	}

	id, err := sdk.ParseAccountIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	account, err := client.Accounts.ShowByID(ctx, id.AsAccountObjectIdentifier())
	if err != nil {
		return nil, err
	}

	if _, err := ImportName[sdk.AccountIdentifier](context.Background(), d, nil); err != nil {
		return nil, err
	}

	if account.RegionGroup != nil {
		if err = d.Set("region_group", *account.RegionGroup); err != nil {
			return nil, err
		}
	}

	var comment string
	if account.Comment != nil {
		comment = *account.Comment
	}

	if err := errors.Join(
		d.Set("edition", string(*account.Edition)),
		d.Set("region", account.SnowflakeRegion),
		d.Set("comment", comment),
		d.Set("is_org_admin", booleanStringFromBool(*account.IsOrgAdmin)),
	); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	isOrgAdmin, err := client.ContextFunctions.IsRoleInSession(ctx, snowflakeroles.Orgadmin)
	if err != nil {
		return diag.FromErr(err)
	}
	if !isOrgAdmin {
		return diag.FromErr(errors.New("current user does not have the ORGADMIN role in session"))
	}

	id := sdk.NewAccountObjectIdentifier(d.Get("name").(string))

	opts := &sdk.CreateAccountOptions{
		AdminName: d.Get("admin_name").(string),
		Email:     d.Get("email").(string),
		Edition:   sdk.AccountEdition(d.Get("edition").(string)),
	}

	if v, ok := d.GetOk("admin_password"); ok {
		opts.AdminPassword = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("admin_rsa_public_key"); ok {
		opts.AdminRSAPublicKey = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("admin_user_type"); ok {
		userType, err := sdk.ToUserType(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.AdminUserType = &userType
	}
	if v, ok := d.GetOk("first_name"); ok {
		opts.FirstName = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("last_name"); ok {
		opts.LastName = sdk.String(v.(string))
	}
	if v := d.Get("must_change_password"); v != BooleanDefault {
		parsedBool, err := booleanStringToBool(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.MustChangePassword = &parsedBool
	}
	if v, ok := d.GetOk("region_group"); ok {
		opts.RegionGroup = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("region"); ok {
		opts.Region = sdk.String(v.(string))
	}
	if v, ok := d.GetOk("comment"); ok {
		opts.Comment = sdk.String(v.(string))
	}

	createResponse, err := client.Accounts.Create(ctx, id, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := util.Retry(5, 3*time.Second, func() (error, bool) {
		_, err = client.Accounts.ShowByID(ctx, id)
		if err != nil {
			log.Printf("[DEBUG] retryable operation resulted in error: %v", err)
			return nil, false
		}
		return nil, true
	}); err != nil {
		return diag.FromErr(fmt.Errorf("failed to query account (%s) after creation, err: %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(sdk.NewAccountIdentifier(createResponse.OrganizationName, createResponse.AccountName)))

	if v, ok := d.GetOk("is_org_admin"); ok && v == BooleanTrue {
		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			SetIsOrgAdmin: &sdk.AccountSetIsOrgAdmin{
				Name:     id,
				OrgAdmin: true,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadAccount(false)(ctx, d, meta)
}

func ReadAccount(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client

		isOrgAdmin, err := client.ContextFunctions.IsRoleInSession(ctx, snowflakeroles.Orgadmin)
		if err != nil {
			return diag.FromErr(err)
		}
		if !isOrgAdmin {
			return diag.FromErr(errors.New("current user does not have the ORGADMIN role in session"))
		}

		id, err := sdk.ParseAccountIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		account, err := client.Accounts.ShowByIDSafely(ctx, id.AsAccountObjectIdentifier())
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query account. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Account: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			var regionGroup string
			if account.RegionGroup != nil {
				regionGroup = *account.RegionGroup

				// For organizations that have accounts in multiple region groups, returns <region_group>.<region> so we need to split on "."
				parts := strings.Split(regionGroup, ".")
				if len(parts) == 2 {
					regionGroup = parts[0]
				}
			}
			var comment string
			if account.Comment != nil {
				comment = *account.Comment
			}

			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"edition", "edition", *account.Edition, *account.Edition, nil},
				outputMapping{"is_org_admin", "is_org_admin", *account.IsOrgAdmin, booleanStringFromBool(*account.IsOrgAdmin), nil},
				outputMapping{"region_group", "region_group", regionGroup, regionGroup, nil},
				outputMapping{"snowflake_region", "region", account.SnowflakeRegion, account.SnowflakeRegion, nil},
				outputMapping{"comment", "comment", comment, comment, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err = setStateToValuesFromConfig(d, accountSchema, []string{
				"name",
				"admin_name",
				"admin_password",
				"admin_rsa_public_key",
				"admin_user_type",
				"first_name",
				"last_name",
				"email",
				"must_change_password",
				"edition",
				"region_group",
				"region",
				"comment",
				"is_org_admin",
				"grace_period_in_days",
			}); err != nil {
				return diag.FromErr(err)
			}
		}

		if errs := errors.Join(
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.AccountToSchema(account)}),
		); errs != nil {
			return diag.FromErr(errs)
		}

		return nil
	}
}

func UpdateAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	isOrgAdmin, err := client.ContextFunctions.IsRoleInSession(ctx, snowflakeroles.Orgadmin)
	if err != nil {
		return diag.FromErr(err)
	}
	if !isOrgAdmin {
		return diag.FromErr(errors.New("current user does not have the ORGADMIN role in session"))
	}

	id, err := sdk.ParseAccountIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewAccountIdentifier(id.OrganizationName(), d.Get("name").(string))

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Rename: &sdk.AccountRename{
				Name:    id.AsAccountObjectIdentifier(),
				NewName: newId.AsAccountObjectIdentifier(),
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	if d.HasChange("is_org_admin") {
		oldIsOrgAdmin, newIsOrgAdmin := d.GetChange("is_org_admin")

		// Setting from default to false and vice versa is not allowed because Snowflake throws an error on already disabled IsOrgAdmin
		canUpdate := true
		if (oldIsOrgAdmin.(string) == BooleanFalse && newIsOrgAdmin.(string) == BooleanDefault) ||
			(oldIsOrgAdmin.(string) == BooleanDefault && newIsOrgAdmin.(string) == BooleanFalse) {
			canUpdate = false
		}

		if canUpdate {
			if newIsOrgAdmin.(string) != BooleanDefault {
				parsed, err := booleanStringToBool(newIsOrgAdmin.(string))
				if err != nil {
					return diag.FromErr(err)
				}
				if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
					SetIsOrgAdmin: &sdk.AccountSetIsOrgAdmin{
						Name:     id.AsAccountObjectIdentifier(),
						OrgAdmin: parsed,
					},
				}); err != nil {
					return diag.FromErr(err)
				}
			} else {
				// No unset available for this field (setting Snowflake default)
				if err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
					SetIsOrgAdmin: &sdk.AccountSetIsOrgAdmin{
						Name:     id.AsAccountObjectIdentifier(),
						OrgAdmin: false,
					},
				}); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	return ReadAccount(false)(ctx, d, meta)
}

func DeleteAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	isOrgAdmin, err := client.ContextFunctions.IsRoleInSession(ctx, snowflakeroles.Orgadmin)
	if err != nil {
		return diag.FromErr(err)
	}
	if !isOrgAdmin {
		return diag.FromErr(errors.New("current user does not have the ORGADMIN role in session"))
	}

	id, err := sdk.ParseAccountIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Accounts.Drop(ctx, id.AsAccountObjectIdentifier(), d.Get("grace_period_in_days").(int), &sdk.DropAccountOptions{
		IfExists: sdk.Bool(true),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
