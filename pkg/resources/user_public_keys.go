package resources

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userPublicKeyProperties = []string{
	"rsa_public_key",
	"rsa_public_key_2",
}

// sanitize input to suppress diffs, etc.
func publicKeyStateFunc(v interface{}) string {
	value := v.(string)
	value = strings.TrimSuffix(value, "\n")
	return value
}

var userPublicKeysSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the user.",
		ForceNew:    true,
	},

	"rsa_public_key": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the user’s RSA public key; used for key-pair authentication. Must be on 1 line without header and trailer.",
		StateFunc:   publicKeyStateFunc,
	},
	"rsa_public_key_2": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the user’s second RSA public key; used to rotate the public and Public keys for key-pair authentication based on an expiration schedule set by your organization. Must be on 1 line without header and trailer.",
		StateFunc:   publicKeyStateFunc,
	},
}

func UserPublicKeys() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.UserPublicKeysResource), TrackingCreateWrapper(resources.UserPublicKeys, CreateUserPublicKeys)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.UserPublicKeysResource), TrackingReadWrapper(resources.UserPublicKeys, ReadUserPublicKeys)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.UserPublicKeysResource), TrackingUpdateWrapper(resources.UserPublicKeys, UpdateUserPublicKeys)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.UserPublicKeysResource), TrackingDeleteWrapper(resources.UserPublicKeys, DeleteUserPublicKeys)),

		Schema: userPublicKeysSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func checkUserExists(ctx context.Context, client *sdk.Client, userId sdk.AccountObjectIdentifier) (bool, error) {
	// First check if user exists
	_, err := client.Users.Describe(ctx, userId)
	if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
		log.Printf("[DEBUG] user (%s) not found", userId.Name())
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func ReadUserPublicKeys(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.AccountObjectIdentifier)

	exists, err := checkUserExists(ctx, client, id)
	if err != nil {
		return diag.FromErr(err)
	}
	// If not found, mark resource to be removed from state file during apply or refresh
	if !exists {
		d.SetId("")
		return nil
	}
	// we can't really read the public keys back from Snowflake so assume they haven't changed
	return nil
}

func CreateUserPublicKeys(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	name := d.Get("name").(string)

	for _, prop := range userPublicKeyProperties {
		publicKey, publicKeyOK := d.GetOk(prop)
		if !publicKeyOK {
			continue
		}
		err := updateUserPublicKeys(db, name, prop, publicKey.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(name)
	return ReadUserPublicKeys(ctx, d, meta)
}

func UpdateUserPublicKeys(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	name := d.Id()

	propsToSet := map[string]string{}
	propsToUnset := map[string]string{}

	for _, prop := range userPublicKeyProperties {
		// if key hasn't changed, continue
		if !d.HasChange(prop) {
			continue
		}
		// if it has changed then we should do something about it
		publicKey, publicKeyOK := d.GetOk(prop)
		if publicKeyOK { // if set, then we should update the value
			propsToSet[prop] = publicKey.(string)
		} else { // if now unset, we should unset the key from the user
			propsToUnset[prop] = publicKey.(string)
		}
	}

	// set the keys we decided should be set
	for prop, value := range propsToSet {
		err := updateUserPublicKeys(db, name, prop, value)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// unset the keys we decided should be unset
	for k := range propsToUnset {
		err := unsetUserPublicKeys(db, name, k)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// re-sync
	return ReadUserPublicKeys(ctx, d, meta)
}

func DeleteUserPublicKeys(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	name := d.Id()

	for _, prop := range userPublicKeyProperties {
		err := unsetUserPublicKeys(db, name, prop)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")
	return nil
}

func updateUserPublicKeys(db *sql.DB, name string, prop string, value string) error {
	stmt := fmt.Sprintf(`ALTER USER "%s" SET %s = '%s'`, name, prop, value)
	return snowflake.Exec(db, stmt)
}

func unsetUserPublicKeys(db *sql.DB, name string, prop string) error {
	stmt := fmt.Sprintf(`ALTER USER "%s" UNSET %s`, name, prop)
	return snowflake.Exec(db, stmt)
}
