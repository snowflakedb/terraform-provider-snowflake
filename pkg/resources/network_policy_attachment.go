package resources

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var networkPolicyAttachmentSchema = map[string]*schema.Schema{
	"network_policy_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the network policy; must be unique for the account in which the network policy is created.",
		ForceNew:    true,
	},
	"set_for_account": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Specifies whether the network policy should be applied globally to your Snowflake account<br><br>**Note:** The Snowflake user running `terraform apply` must be on an IP address allowed by the network policy to set that policy globally on the Snowflake account.<br><br>Additionally, a Snowflake account can only have one network policy set globally at any given time. This resource does not enforce one-policy-per-account, it is the user's responsibility to enforce this. If multiple network policy resources have `set_for_account: true`, the final policy set on the account will be non-deterministic.",
		Default:     false,
	},
	"users": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "Specifies which users the network policy should be attached to",
	},
}

// NetworkPolicyAttachment returns a pointer to the resource representing a network policy attachment.
func NetworkPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.NetworkPolicyAttachmentResource), TrackingCreateWrapper(resources.NetworkPolicyAttachment, CreateNetworkPolicyAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.NetworkPolicyAttachmentResource), TrackingReadWrapper(resources.NetworkPolicyAttachment, ReadNetworkPolicyAttachment)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.NetworkPolicyAttachmentResource), TrackingUpdateWrapper(resources.NetworkPolicyAttachment, UpdateNetworkPolicyAttachment)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.NetworkPolicyAttachmentResource), TrackingDeleteWrapper(resources.NetworkPolicyAttachment, DeleteNetworkPolicyAttachment)),

		Schema: networkPolicyAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

// CreateNetworkPolicyAttachment implements schema.CreateFunc.
func CreateNetworkPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	policyName := d.Get("network_policy_name").(string)
	d.SetId(policyName + "_attachment")

	if d.Get("set_for_account").(bool) {
		if err := setOnAccount(ctx, d, meta); err != nil {
			return diag.FromErr(fmt.Errorf("error creating attachment for network policy %v err = %w", policyName, err))
		}
	}

	if u, ok := d.GetOk("users"); ok {
		users := expandStringList(u.(*schema.Set).List())

		if err := ensureUserAlterPrivileges(ctx, users, meta); err != nil {
			return diag.FromErr(err)
		}

		if err := setOnUsers(ctx, users, d, meta); err != nil {
			return diag.FromErr(fmt.Errorf("error creating attachment for network policy %v err = %w", policyName, err))
		}
	}

	return ReadNetworkPolicyAttachment(ctx, d, meta)
}

// ReadNetworkPolicyAttachment implements schema.ReadFunc.
func ReadNetworkPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	policyName := strings.Replace(d.Id(), "_attachment", "", 1)

	var currentUsers []string
	if err := d.Set("network_policy_name", policyName); err != nil {
		return diag.FromErr(err)
	}

	if u, ok := d.GetOk("users"); ok {
		users := expandStringList(u.(*schema.Set).List())
		for _, user := range users {
			parameter, err := client.Parameters.ShowUserParameter(ctx, sdk.UserParameterNetworkPolicy, sdk.NewAccountObjectIdentifier(user))
			if err != nil {
				log.Printf("[DEBUG] network policy (%s) not found on user (%s)", d.Id(), user)
				continue
			}

			if parameter.Level == "USER" && parameter.Key == "NETWORK_POLICY" && parameter.Value == policyName {
				currentUsers = append(currentUsers, user)
			}
		}

		if err := d.Set("users", currentUsers); err != nil {
			return diag.FromErr(err)
		}
	}

	isSetOnAccount := false

	parameter, err := client.Parameters.ShowAccountParameter(ctx, sdk.AccountParameterNetworkPolicy)
	if err != nil {
		log.Printf("[DEBUG] network policy (%s) not found on account", d.Id())
		isSetOnAccount = false
	}

	if err == nil && parameter.Level == "ACCOUNT" && parameter.Key == "NETWORK_POLICY" && parameter.Value == policyName {
		isSetOnAccount = true
	}

	if err := d.Set("set_for_account", isSetOnAccount); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// UpdateNetworkPolicyAttachment implements schema.UpdateFunc.
func UpdateNetworkPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if d.HasChange("set_for_account") {
		oldAcctFlag, newAcctFlag := d.GetChange("set_for_account")
		if newAcctFlag.(bool) {
			if err := setOnAccount(ctx, d, meta); err != nil {
				return diag.FromErr(err)
			}
		} else if !newAcctFlag.(bool) && oldAcctFlag == true {
			if err := unsetOnAccount(ctx, d, meta); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("users") {
		o, n := d.GetChange("users")
		oldUsersSet := o.(*schema.Set)
		newUsersSet := n.(*schema.Set)

		removedUsers := expandStringList(oldUsersSet.Difference(newUsersSet).List())
		addedUsers := expandStringList(newUsersSet.Difference(oldUsersSet).List())

		if err := ensureUserAlterPrivileges(ctx, removedUsers, meta); err != nil {
			return diag.FromErr(err)
		}

		if err := ensureUserAlterPrivileges(ctx, addedUsers, meta); err != nil {
			return diag.FromErr(err)
		}

		for _, user := range removedUsers {
			if err := unsetOnUser(ctx, user, d, meta); err != nil {
				return diag.FromErr(err)
			}
		}

		for _, user := range addedUsers {
			if err := setOnUser(ctx, user, d, meta); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return ReadNetworkPolicyAttachment(ctx, d, meta)
}

// DeleteNetworkPolicyAttachment implements schema.DeleteFunc.
func DeleteNetworkPolicyAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	policyName := d.Get("network_policy_name").(string)
	d.SetId(policyName + "_attachment")

	if d.Get("set_for_account").(bool) {
		if err := unsetOnAccount(ctx, d, meta); err != nil {
			return diag.FromErr(fmt.Errorf("error deleting attachment for network policy %v err = %w", policyName, err))
		}
	}

	if u, ok := d.GetOk("users"); ok {
		users := expandStringList(u.(*schema.Set).List())

		if err := ensureUserAlterPrivileges(ctx, users, meta); err != nil {
			return diag.FromErr(err)
		}

		if err := unsetOnUsers(ctx, users, d, meta); err != nil {
			return diag.FromErr(fmt.Errorf("error deleting attachment for network policy %v err = %w", policyName, err))
		}
	}

	return nil
}

// setOnAccount sets the network policy globally for the Snowflake account
// Note: the ip address of the session executing this SQL must be allowed by the network policy being set.
func setOnAccount(ctx context.Context, d *schema.ResourceData, meta any) error {
	client := meta.(*provider.Context).Client

	policyName := d.Get("network_policy_name").(string)

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{LegacyParameters: &sdk.AccountLevelParameters{ObjectParameters: &sdk.ObjectParameters{NetworkPolicy: sdk.String(policyName)}}}})
	if err != nil {
		return fmt.Errorf("error setting network policy %v on account err = %w", policyName, err)
	}

	return nil
}

// setOnAccount unsets the network policy globally for the Snowflake account.
func unsetOnAccount(ctx context.Context, d *schema.ResourceData, meta any) error {
	client := meta.(*provider.Context).Client

	policyName := d.Get("network_policy_name").(string)

	err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: &sdk.AccountUnset{Parameters: &sdk.AccountLevelParametersUnset{ObjectParameters: &sdk.ObjectParametersUnset{NetworkPolicy: sdk.Bool(true)}}}})
	if err != nil {
		return fmt.Errorf("error unsetting network policy %v on account err = %w", policyName, err)
	}

	return nil
}

// setOnUsers sets the network policy for list of users.
func setOnUsers(ctx context.Context, users []string, data *schema.ResourceData, meta interface{}) error {
	policyName := data.Get("network_policy_name").(string)
	for _, user := range users {
		if err := setOnUser(ctx, user, data, meta); err != nil {
			return fmt.Errorf("error setting network policy %v on user %v err = %w", policyName, user, err)
		}
	}

	return nil
}

// setOnUser sets the network policy for a given user.
func setOnUser(ctx context.Context, user string, data *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client

	policyName := data.Get("network_policy_name").(string)

	err := client.Users.Alter(ctx, sdk.NewAccountObjectIdentifier(user), &sdk.AlterUserOptions{Set: &sdk.UserSet{ObjectParameters: &sdk.UserObjectParameters{NetworkPolicy: sdk.Pointer(sdk.NewAccountObjectIdentifier(policyName))}}})
	if err != nil {
		return fmt.Errorf("error setting network policy %v on user %v err = %w", policyName, user, err)
	}

	return nil
}

// unsetOnUsers unsets the network policy for list of users.
func unsetOnUsers(ctx context.Context, users []string, data *schema.ResourceData, meta interface{}) error {
	policyName := data.Get("network_policy_name").(string)
	for _, user := range users {
		if err := unsetOnUser(ctx, user, data, meta); err != nil {
			return fmt.Errorf("error unsetting network policy %v on user %v err = %w", policyName, user, err)
		}
	}

	return nil
}

// unsetOnUser sets the network policy for a given user.
func unsetOnUser(ctx context.Context, user string, data *schema.ResourceData, meta interface{}) error {
	client := meta.(*provider.Context).Client

	policyName := data.Get("network_policy_name").(string)

	err := client.Users.Alter(ctx, sdk.NewAccountObjectIdentifier(user), &sdk.AlterUserOptions{Unset: &sdk.UserUnset{ObjectParameters: &sdk.UserObjectParametersUnset{NetworkPolicy: sdk.Bool(true)}}})
	if err != nil {
		return fmt.Errorf("error unsetting network policy %v on user %v err = %w", policyName, user, err)
	}

	return nil
}

// ensureUserAlterPrivileges ensures the executing Snowflake user can alter each user in the set of users.
func ensureUserAlterPrivileges(ctx context.Context, users []string, meta interface{}) error {
	client := meta.(*provider.Context).Client

	for _, user := range users {
		_, err := client.Users.Describe(ctx, sdk.NewAccountObjectIdentifier(user))
		if err != nil {
			return fmt.Errorf("error describing user %v err = %w", user, err)
		}
	}

	return nil
}
