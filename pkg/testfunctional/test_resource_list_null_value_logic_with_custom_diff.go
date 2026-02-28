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
	listNullValueLogicHelperFieldNull  = "null"
	listNullValueLogicHelperFieldEmpty = "empty"
	listNullValueLogicHelperFieldItems = "items"
)

func listNullValueLogicHelperFieldFromRawConfig(d *schema.ResourceDiff) string {
	rawConfigValue := d.GetRawConfig().AsValueMap()["nullable_list"]
	if rawConfigValue.IsNull() {
		return listNullValueLogicHelperFieldNull
	}
	if len(rawConfigValue.AsValueSlice()) == 0 {
		return listNullValueLogicHelperFieldEmpty
	}
	return listNullValueLogicHelperFieldItems
}

var testResourceListNullValueLogicWithHelperFieldSchema = func() map[string]*schema.Schema {
	s := map[string]*schema.Schema{}
	for k, v := range testResourceListNullValueLogicSchema {
		s[k] = v
	}
	s["nullable_list_presence"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	return s
}()

// TestResourceListNullValueLogicWithHelperField is the same resource as TestResourceListNullValueLogic,
// but with an additional nullable_list_presence computed helper field and a CustomizeDiff that forces
// a plan diff on null <-> empty transitions, proving that the helper field approach solves the limitation.
func TestResourceListNullValueLogicWithHelperField() *schema.Resource {
	return &schema.Resource{
		CreateContext: testResourceListNullValueLogicWithHelperFieldCreate,
		UpdateContext: testResourceListNullValueLogicWithHelperFieldUpdate,
		ReadContext:   testResourceListNullValueLogicWithHelperFieldRead,
		DeleteContext: testResourceListNullValueLogicDelete,

		Schema: testResourceListNullValueLogicWithHelperFieldSchema,

		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ any) error {
			configPresence := listNullValueLogicHelperFieldFromRawConfig(d)
			statePresence, _ := d.GetOk("nullable_list_presence")
			if statePresence != configPresence {
				log.Printf("[DEBUG] nullable_list_presence diff: state=%q config=%q, forcing new value", statePresence, configPresence)
				return d.SetNew("nullable_list_presence", configPresence)
			}
			return nil
		},
	}
}

func observeAndStoreListWithHelperField(d *schema.ResourceData, envName string) diag.Diagnostics {
	if diags := observeAndStoreList(d, envName); diags.HasError() {
		return diags
	}

	rawConfigValue := d.GetRawConfig().AsValueMap()["nullable_list"]
	getResult := d.Get("nullable_list").([]any)

	var presence string
	switch {
	case rawConfigValue.IsNull():
		presence = listNullValueLogicHelperFieldNull
	case len(getResult) == 0:
		presence = listNullValueLogicHelperFieldEmpty
	default:
		presence = listNullValueLogicHelperFieldItems
	}
	if err := d.Set("nullable_list_presence", presence); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func testResourceListNullValueLogicWithHelperFieldCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Get("env_name").(string)
	log.Printf("[DEBUG] handling create (with helper field) for %s", envName)

	d.SetId(envName)

	if diags := observeAndStoreListWithHelperField(d, envName); diags.HasError() {
		return diags
	}

	return testResourceListNullValueLogicWithHelperFieldRead(ctx, d, meta)
}

func testResourceListNullValueLogicWithHelperFieldUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling update (with helper field) for %s", envName)

	if diags := observeAndStoreListWithHelperField(d, envName); diags.HasError() {
		return diags
	}

	return testResourceListNullValueLogicWithHelperFieldRead(ctx, d, meta)
}

func testResourceListNullValueLogicWithHelperFieldRead(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	envName := d.Id()
	value := oswrapper.Getenv(envName)
	log.Printf("[DEBUG] handling read (with helper field) for %s, env value: %s", envName, value)

	switch {
	case value == "" || value == listNullValueLogicEnvNull:
		if err := d.Set("nullable_list", nil); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("nullable_list_presence", listNullValueLogicHelperFieldNull); err != nil {
			return diag.FromErr(err)
		}
	case value == listNullValueLogicEnvEmpty:
		if err := d.Set("nullable_list", []string{}); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("nullable_list_presence", listNullValueLogicHelperFieldEmpty); err != nil {
			return diag.FromErr(err)
		}
	default:
		items := strings.Split(value, ",")
		if err := d.Set("nullable_list", items); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("nullable_list_presence", listNullValueLogicHelperFieldItems); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
