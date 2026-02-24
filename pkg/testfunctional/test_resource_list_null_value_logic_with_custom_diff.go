package testfunctional

import (
	"context"
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	listPresenceNull  = "null"
	listPresenceEmpty = "empty"
	listPresenceItems = "items"
)

func listPresenceFromRawConfig(d *schema.ResourceDiff) string {
	rawConfigValue := d.GetRawConfig().AsValueMap()["nullable_list"]
	if rawConfigValue.IsNull() {
		return listPresenceNull
	}
	if len(rawConfigValue.AsValueSlice()) == 0 {
		return listPresenceEmpty
	}
	return listPresenceItems
}

// TestResourceListNullValueLogicWithPresence is the same resource as TestResourceListNullValueLogic,
// but with an additional nullable_list_presence computed field and a CustomizeDiff that forces a plan
// diff on null <-> empty transitions, proving that the helper field approach solves the limitation.
func TestResourceListNullValueLogicWithPresence() *schema.Resource {
	schemaWithPresence := map[string]*schema.Schema{}
	for k, v := range testResourceListNullValueLogicSchema {
		schemaWithPresence[k] = v
	}
	schemaWithPresence["nullable_list_presence"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		CreateContext: testResourceListNullValueLogicWithPresenceCreate,
		UpdateContext: testResourceListNullValueLogicWithPresenceUpdate,
		ReadContext:   testResourceListNullValueLogicWithPresenceRead,
		DeleteContext: testResourceListNullValueLogicDelete,

		Schema: schemaWithPresence,

		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ any) error {
			configPresence := listPresenceFromRawConfig(d)
			statePresence, _ := d.GetOk("nullable_list_presence")
			if statePresence != configPresence {
				log.Printf("[DEBUG] nullable_list_presence diff: state=%q config=%q, forcing new value", statePresence, configPresence)
				return d.SetNew("nullable_list_presence", configPresence)
			}
			return nil
		},
	}
}

func observeAndStoreListWithPresence(d *schema.ResourceData, envName string) diag.Diagnostics {
	if diags := observeAndStoreList(d, envName); diags.HasError() {
		return diags
	}

	rawConfigValue := d.GetRawConfig().AsValueMap()["nullable_list"]
	getResult := d.Get("nullable_list").([]any)

	var presence string
	switch {
	case rawConfigValue.IsNull():
		presence = listPresenceNull
	case len(getResult) == 0:
		presence = listPresenceEmpty
	default:
		presence = listPresenceItems
	}
	if err := d.Set("nullable_list_presence", presence); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func testResourceListNullValueLogicWithPresenceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Get("env_name").(string)
	log.Printf("[DEBUG] handling create (with presence) for %s", envName)

	d.SetId(envName)

	if diags := observeAndStoreListWithPresence(d, envName); diags.HasError() {
		return diags
	}

	return testResourceListNullValueLogicWithPresenceRead(ctx, d, meta)
}

func testResourceListNullValueLogicWithPresenceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling update (with presence) for %s", envName)

	if diags := observeAndStoreListWithPresence(d, envName); diags.HasError() {
		return diags
	}

	return testResourceListNullValueLogicWithPresenceRead(ctx, d, meta)
}

func testResourceListNullValueLogicWithPresenceRead(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	envName := d.Id()
	value := oswrapper.Getenv(envName)
	log.Printf("[DEBUG] handling read (with presence) for %s, env value: %s", envName, value)

	switch {
	case value == "" || value == listNullValueLogicEnvNull:
		if err := d.Set("nullable_list", nil); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("nullable_list_presence", listPresenceNull); err != nil {
			return diag.FromErr(err)
		}
	case value == listNullValueLogicEnvEmpty:
		if err := d.Set("nullable_list", []string{}); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("nullable_list_presence", listPresenceEmpty); err != nil {
			return diag.FromErr(err)
		}
	default:
		items := strings.Split(value, ",")
		if err := d.Set("nullable_list", items); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("nullable_list_presence", listPresenceItems); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
