package resources

import (
	"context"
	"log"
	"os"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-2054208]: extract to the dedicated package
var testResourceDataTypeDiffHandlingSchema = map[string]*schema.Schema{
	"env_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Used to make the tests faster (instead of communicating with SF, we read from environment variable).",
	},
	"top_level_datatype": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "An example field being a data type.",
		DiffSuppressFunc: DiffSuppressDataTypes,
		ValidateDiagFunc: IsDataTypeValid,
		StateFunc:        DataTypeStateFunc,
	},
}

func TestResourceDataTypeDiffHandling() *schema.Resource {
	return &schema.Resource{
		CreateContext: TestResourceDataTypeDiffHandlingCreate,
		UpdateContext: TestResourceDataTypeDiffHandlingUpdate,
		ReadContext:   TestResourceDataTypeDiffHandlingRead(true),
		DeleteContext: TestResourceDataTypeDiffHandlingDelete,

		Schema: testResourceDataTypeDiffHandlingSchema,
	}
}

func TestResourceDataTypeDiffHandlingCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Get("env_name").(string)
	log.Printf("[DEBUG] handling create for %s", envName)

	if err := handleDatatypeCreate(d, "top_level_datatype", func(dataType datatypes.DataType) error {
		return testResourceDataTypeDiffHandlingSet(envName, dataType)
	}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(envName)
	return TestResourceDataTypeDiffHandlingRead(false)(ctx, d, meta)
}

func TestResourceDataTypeDiffHandlingUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling update for %s", envName)

	if err := handleDatatypeUpdate(d, "top_level_datatype", func(dataType datatypes.DataType) error {
		return testResourceDataTypeDiffHandlingSet(envName, dataType)
	}); err != nil {
		return diag.FromErr(err)
	}

	return TestResourceDataTypeDiffHandlingRead(false)(ctx, d, meta)
}

func TestResourceDataTypeDiffHandlingRead(withExternalChangesMarking bool) schema.ReadContextFunc {
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

			if err := handleDatatypeSet(d, "top_level_datatype", externalDataType); err != nil {
				return diag.FromErr(err)
			}
		}
		return nil
	}
}

func TestResourceDataTypeDiffHandlingDelete(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling delete for %s", envName)

	if err := testResourceDataTypeDiffHandlingUnset(envName); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func testResourceDataTypeDiffHandlingSet(envName string, dataType datatypes.DataType) error {
	log.Printf("[DEBUG] setting env %s to value `%s`", envName, dataType.ToSql())
	return os.Setenv(envName, dataType.ToSql())
}

func testResourceDataTypeDiffHandlingUnset(envName string) error {
	log.Printf("[DEBUG] unsetting env %s", envName)
	return os.Setenv(envName, "")
}
