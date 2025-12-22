package testfunctional

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
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
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"nullable_list_output": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Saves the state of the nullable list",
	},
}

// TestResourceListNullValueLogic explores the list data to implement the logic for the following states:
// - null
// - empty list
// - filled list
// Usually, the resources handle the null and empty list together, so this resource checks if the logic between
// mainly those two states can be separated.
func TestResourceListNullValueLogic() *schema.Resource {
	return &schema.Resource{
		CreateContext: TestResourceListNullValueLogicCreate,
		UpdateContext: TestResourceListNullValueLogicUpdate,
		ReadContext:   TestResourceListNullValueLogicRead(true),
		DeleteContext: TestResourceListNullValueLogicDelete,

		Schema: testResourceListNullValueLogicSchema,
	}
}

func TestResourceListNullValueLogicCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Get("env_name").(string)
	log.Printf("[DEBUG] handling create for %s", envName)

	if err := resources.HandleDatatypeCreate(d, "top_level_datatype", func(dataType datatypes.DataType) error {
		return testResourceDataTypeDiffHandlingSet(envName, dataType)
	}); err != nil {
		return diag.FromErr(err)
	}

	strings.Join(d.Get("nullable_list_output").([]string), ",")
	log.Printf("[DEBUG] setting env %s to value `%s`", envName)
	if err := os.Setenv(envName); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(envName)
	return TestResourceDataTypeDiffHandlingRead(false)(ctx, d, meta)
}

func TestResourceListNullValueLogicUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling update for %s", envName)

	if err := resources.HandleDatatypeUpdate(d, "top_level_datatype", func(dataType datatypes.DataType) error {
		return testResourceDataTypeDiffHandlingSet(envName, dataType)
	}); err != nil {
		return diag.FromErr(err)
	}

	return TestResourceDataTypeDiffHandlingRead(false)(ctx, d, meta)
}

func TestResourceListNullValueLogicRead(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		envName := d.Id()
		log.Printf("[DEBUG] handling read for %s, with marking external changes: %t", envName, withExternalChangesMarking)

		value := oswrapper.Getenv(envName)
		log.Printf("[DEBUG] env %s value is `%s`", envName, value)
		if value != "" {
			externalDataType, err := datatypes.ParseDataType(value)
			if err != nil {
				return diag.FromErr(err)
			}

			if err := resources.HandleDatatypeSet(d, "top_level_datatype", externalDataType); err != nil {
				return diag.FromErr(err)
			}
		}
		return nil
	}
}

func TestResourceListNullValueLogicDelete(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling delete for %s", envName)

	if err := testResourceDataTypeDiffHandlingUnset(envName); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
