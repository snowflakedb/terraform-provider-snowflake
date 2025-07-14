package resources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/maps"
)

var organizationAccountSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: externalChangesNotDetectedFieldDescription("Specifies the identifier (i.e. name) for the organization account. It must be unique within an organization, regardless of which Snowflake Region the organization account is in and must start with an alphabetic character and cannot contain spaces or special characters except for underscores (_). Note that if the organization account name includes underscores, features that do not accept account names with underscores (e.g. Okta SSO or SCIM) can reference a version of the account name that substitutes hyphens (-) for the underscores."),
	},
	"admin_name": {
		Type:        schema.TypeString,
		Required:    true,
		Sensitive:   true,
		Description: externalChangesNotDetectedFieldDescription("Login name of the initial administrative user of the organization account. A new user is created in the new organization account with this name and password and granted the GLOBALORGADMIN role in the account. A login name can be any string consisting of letters, numbers, and underscores. Login names are always case-insensitive."),
	},
	"admin_password": {
		Type:         schema.TypeString,
		Optional:     true,
		Sensitive:    true,
		Description:  externalChangesNotDetectedFieldDescription("Password for the initial administrative user of the organization account. Either admin_password or admin_rsa_public_key has to be specified."),
		ExactlyOneOf: []string{"admin_password", "admin_rsa_public_key"},
	},
	"admin_rsa_public_key": {
		Type:         schema.TypeString,
		Optional:     true,
		Description:  externalChangesNotDetectedFieldDescription("Assigns a public key to the initial administrative user of the organization account. Either admin_password or admin_rsa_public_key has to be specified."),
		AtLeastOneOf: []string{"admin_password", "admin_rsa_public_key"},
	},
	"first_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: externalChangesNotDetectedFieldDescription("First name of the initial administrative user of the organization account."),
	},
	"last_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: externalChangesNotDetectedFieldDescription("Last name of the initial administrative user of the organization account."),
	},
	"email": {
		Type:        schema.TypeString,
		Required:    true,
		Sensitive:   true,
		Description: externalChangesNotDetectedFieldDescription("Email address of the initial administrative user of the organization account. This email address is used to send any notifications about the account."),
	},
	"must_change_password": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: externalChangesNotDetectedFieldDescription("Specifies whether the new user created to administer the organization account is forced to change their password upon first login into the account."),
	},
	"edition": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      externalChangesNotDetectedFieldDescription(fmt.Sprintf("Snowflake Edition of the organization account. See more about Snowflake Editions in the [official documentation](https://docs.snowflake.com/en/user-guide/intro-editions). Valid options are: %s", docs.PossibleValuesListed(sdk.AllOrganizationAccountEditions))),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToOrganizationAccountEdition),
		ValidateDiagFunc: sdkValidation(sdk.ToOrganizationAccountEdition),
	},
	"region_group": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: externalChangesNotDetectedFieldDescription("ID of the region group where the organization account is created. To retrieve the region group ID for existing accounts in your organization, execute the [SHOW REGIONS](https://docs.snowflake.com/en/sql-reference/sql/show-regions) command. For information about when you might need to specify region group, see [Region groups](https://docs.snowflake.com/en/user-guide/admin-account-identifier.html#label-region-groups)."),
	},
	"region": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: externalChangesNotDetectedFieldDescription("[Snowflake Region ID](https://docs.snowflake.com/en/user-guide/admin-account-identifier#region-ids) of the region where the organization account is created. If no value is provided, Snowflake creates the organization account in the same Snowflake Region as the current account (i.e. the account in which the CREATE ORGANIZATION ACCOUNT statement is executed.)"),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: externalChangesNotDetectedFieldDescription("Specifies a comment for the organization account."),
	},
}

func OrganizationAccount() *schema.Resource {
	return &schema.Resource{
		Description: "The organization account resource allows you to create organization accounts. For more information, check [organization account documentation](https://docs.snowflake.com/en/user-guide/organizations-manage-accounts).",

		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.OrganizationAccountResource), TrackingCreateWrapper(resources.OrganizationAccount, CreateOrganizationAccount)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.OrganizationAccountResource), TrackingReadWrapper(resources.OrganizationAccount, ReadOrganizationAccount)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.OrganizationAccountResource), TrackingUpdateWrapper(resources.OrganizationAccount, UpdateOrganizationAccount)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.OrganizationAccountResource), TrackingDeleteWrapper(resources.OrganizationAccount, DeleteOrganizationAccount)),

		Schema: organizationAccountSchema,
	}
}

func CreateOrganizationAccount(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	edition, err := sdk.ToOrganizationAccountEdition(d.Get("edition").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateOrganizationAccountRequest(
		sdk.NewAccountObjectIdentifier(d.Get("name").(string)),
		d.Get("admin_name").(string),
		d.Get("email").(string),
		edition,
	)

	if errs := errors.Join(
		stringAttributeCreateBuilder(d, "admin_password", request.WithAdminPassword),
		stringAttributeCreateBuilder(d, "admin_rsa_public_key", request.WithAdminRsaPublicKey),
		stringAttributeCreateBuilder(d, "first_name", request.WithFirstName),
		stringAttributeCreateBuilder(d, "last_name", request.WithLastName),
		booleanAttributeCreateBuilder(d, "must_change_password", request.WithMustChangePassword),
		stringAttributeCreateBuilder(d, "region_group", request.WithRegionGroup),
		stringAttributeCreateBuilder(d, "region", request.WithRegion),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	); errs != nil {
		return diag.FromErr(err)
	}

	err = client.OrganizationAccounts.Create(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sdk.NewAccountObjectIdentifier(d.Get("name").(string)).Name())

	return ReadOrganizationAccount(ctx, d, meta)
}

func ReadOrganizationAccount(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	return nil
}

func UpdateOrganizationAccount(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	if d.HasChanges(maps.Keys(organizationAccountSchema)...) {
		return diag.FromErr(fmt.Errorf("cannot update organization account, only create is supported; please use the `snowflake_organization_account` resource to alter existing organization accounts"))
	}

	return nil
}

func DeleteOrganizationAccount(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	log.Println("[DEBUG] Deleting organization account is not supported. The resource will be only removed from state. Please contact Snowflake support to delete the organization account if that's what you intended to do.")
	d.SetId("")
	return nil
}
