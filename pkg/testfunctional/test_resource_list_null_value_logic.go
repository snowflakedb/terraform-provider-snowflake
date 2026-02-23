package testfunctional

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testResourceListNullValueLogicSchema = map[string]*schema.Schema{
	"env_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Used to make the tests faster (instead of communicating with SF, we read from environment variable).",
	},
	"nullable_list": {
		Type:     schema.TypeList,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	// Computed observation fields - these record what CRUD functions observed about nullable_list
	// using different SDKv2 methods. This allows tests to assert on the actual SDK behavior.
	"get_ok_result": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Result of d.GetOk('nullable_list') - the ok boolean.",
	},
	"get_result": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "String representation of d.Get('nullable_list').([]any).",
	},
	"raw_config_result": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "String representation of d.GetRawConfig().AsValueMap()['nullable_list'].",
	},
}

const (
	listNullValueLogicEnvNull  = "__NULL__"
	listNullValueLogicEnvEmpty = "__EMPTY__"
)

// TestResourceListNullValueLogic explores the list data to implement the logic for the following states:
// - null (not specified in config)
// - empty list (nullable_list = [])
// - filled list (nullable_list = ["a", "b"])
// Usually, the resources handle the null and empty list together, so this resource checks if the logic between
// mainly those two states can be separated using d.Get, d.GetOk, and d.GetRawConfig.
func TestResourceListNullValueLogic() *schema.Resource {
	return &schema.Resource{
		CreateContext: testResourceListNullValueLogicCreate,
		UpdateContext: testResourceListNullValueLogicUpdate,
		ReadContext:   testResourceListNullValueLogicRead,
		DeleteContext: testResourceListNullValueLogicDelete,

		Schema: testResourceListNullValueLogicSchema,
	}
}

// observeAndStoreList examines the nullable_list field using multiple SDKv2 methods,
// stores observations as computed attributes, and persists the list state to the env var.
func observeAndStoreList(d *schema.ResourceData, envName string) diag.Diagnostics {
	// Observation 1: d.GetOk
	_, getOk := d.GetOk("nullable_list")
	log.Printf("[DEBUG] d.GetOk('nullable_list') ok=%t", getOk)
	if err := d.Set("get_ok_result", fmt.Sprintf("%t", getOk)); err != nil {
		return diag.FromErr(err)
	}

	// Observation 2: d.Get raw value
	getResult := d.Get("nullable_list").([]any)
	log.Printf("[DEBUG] d.Get('nullable_list') = %v", getResult)
	if err := d.Set("get_result", fmt.Sprintf("%v", getResult)); err != nil {
		return diag.FromErr(err)
	}

	// Observation 3: d.GetRawConfig value
	rawConfigValue := d.GetRawConfig().AsValueMap()["nullable_list"]
	isNull := rawConfigValue.IsNull()
	var rawConfigStr string
	if isNull {
		rawConfigStr = "null"
	} else {
		slice := rawConfigValue.AsValueSlice()
		items := make([]string, len(slice))
		for i, v := range slice {
			items[i] = v.AsString()
		}
		rawConfigStr = fmt.Sprintf("%v", items)
	}
	log.Printf("[DEBUG] d.GetRawConfig()['nullable_list'] = %s", rawConfigStr)
	if err := d.Set("raw_config_result", rawConfigStr); err != nil {
		return diag.FromErr(err)
	}

	// Store list state in env var (simulates persisting to Snowflake)
	if isNull {
		log.Printf("[DEBUG] storing %s in env %s (null list)", listNullValueLogicEnvNull, envName)
		if err := os.Setenv(envName, listNullValueLogicEnvNull); err != nil {
			return diag.FromErr(err)
		}
	} else if len(getResult) == 0 {
		log.Printf("[DEBUG] storing %s in env %s (empty list)", listNullValueLogicEnvEmpty, envName)
		if err := os.Setenv(envName, listNullValueLogicEnvEmpty); err != nil {
			return diag.FromErr(err)
		}
	} else {
		items := make([]string, len(getResult))
		for i, item := range getResult {
			items[i] = item.(string)
		}
		value := strings.Join(items, ",")
		log.Printf("[DEBUG] storing items in env %s: %s", envName, value)
		if err := os.Setenv(envName, value); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func testResourceListNullValueLogicCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Get("env_name").(string)
	log.Printf("[DEBUG] handling create for %s", envName)

	d.SetId(envName)

	if diags := observeAndStoreList(d, envName); diags.HasError() {
		return diags
	}

	return testResourceListNullValueLogicRead(ctx, d, meta)
}

func testResourceListNullValueLogicUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling update for %s", envName)

	if diags := observeAndStoreList(d, envName); diags.HasError() {
		return diags
	}

	return testResourceListNullValueLogicRead(ctx, d, meta)
}

func testResourceListNullValueLogicRead(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	envName := d.Id()
	value := oswrapper.Getenv(envName)
	log.Printf("[DEBUG] handling read for %s, env value: %s", envName, value)

	switch {
	case value == "" || value == listNullValueLogicEnvNull:
		// List is null / not set - set to nil to represent absence
		if err := d.Set("nullable_list", nil); err != nil {
			return diag.FromErr(err)
		}
	case value == listNullValueLogicEnvEmpty:
		// List is explicitly empty
		if err := d.Set("nullable_list", []string{}); err != nil {
			return diag.FromErr(err)
		}
	default:
		// List has items
		items := strings.Split(value, ",")
		if err := d.Set("nullable_list", items); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func testResourceListNullValueLogicDelete(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling delete for %s", envName)

	if err := os.Setenv(envName, ""); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
