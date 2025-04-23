package resources

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/datatypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-2054208]: extract to the dedicated package
var testResourceDataTypeDiffHandlingListSchema = map[string]*schema.Schema{
	"env_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Used to make the tests faster (instead of communicating with SF, we read from environment variable).",
	},
	"nesting_list": {
		Type:     schema.TypeList,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"nested_datatype": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "An example nested field being a data type.",
					DiffSuppressFunc: DiffSuppressDataTypes,
					ValidateDiagFunc: IsDataTypeValid,
					StateFunc:        DataTypeStateFunc,
					// TODO [SNOW-2054235]: all current use cases have force new on nested data types; try without the force new
					ForceNew: true,
				},
			},
		},
		Required:    true,
		Description: "An example list of objects where one of the nested fields is a data type.",
		// TODO [SNOW-2054235]: all current use cases have force new on nested data types; try without the force new
		ForceNew: true,
	},
}

func TestResourceDataTypeDiffHandlingList() *schema.Resource {
	return &schema.Resource{
		CreateContext: TestResourceDataTypeDiffHandlingListCreate,
		UpdateContext: TestResourceDataTypeDiffHandlingListUpdate,
		ReadContext:   TestResourceDataTypeDiffHandlingListRead(true),
		DeleteContext: TestResourceDataTypeDiffHandlingListDelete,

		Schema: testResourceDataTypeDiffHandlingListSchema,
	}
}

func TestResourceDataTypeDiffHandlingListCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Get("env_name").(string)
	log.Printf("[DEBUG] handling create for %s", envName)

	var dataTypes []datatypes.DataType
	var err error
	if dataTypes, err = handleNestedDataTypeCreate(d, "nesting_list", "nested_datatype", func(v map[string]any, dataType datatypes.DataType) (datatypes.DataType, error) {
		return dataType, nil
	}); err != nil {
		return diag.FromErr(err)
	}
	err = testResourceDataTypeDiffHandlingListSet(envName, dataTypes)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(envName)
	return TestResourceDataTypeDiffHandlingListRead(false)(ctx, d, meta)
}

func TestResourceDataTypeDiffHandlingListUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling update for %s", envName)
	// no-op as we use a ForceNew attribute
	return TestResourceDataTypeDiffHandlingListRead(false)(ctx, d, meta)
}

func TestResourceDataTypeDiffHandlingListRead(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		envName := d.Id()
		log.Printf("[DEBUG] handling read for %s, with marking external changes: %t", envName, withExternalChangesMarking)

		value := oswrapper.Getenv(envName)
		log.Printf("[DEBUG] env %s value is `%s`", envName, value)
		externalDataTypes, envExists, err := testResourceDataTypeDiffHandlingListExternalRead(envName)
		if err != nil {
			return diag.FromErr(err)
		}
		if envExists {
			if err := handleNestedDataTypeSet(d, "nesting_list", "nested_datatype", externalDataTypes,
				func(externalItem datatypes.DataType) datatypes.DataType { return externalItem },
				func(externalItem datatypes.DataType, item map[string]any) {},
			); err != nil {
				return diag.FromErr(err)
			}
		}
		return nil
	}
}

func TestResourceDataTypeDiffHandlingListDelete(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	envName := d.Id()
	log.Printf("[DEBUG] handling delete for %s", envName)

	if err := testResourceDataTypeDiffHandlingListUnset(envName); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// testResourceDataTypeDiffHandlingListExternalRead returns:
// - slice of external datatypes or nil;
// - true if the env variable had any value, false otherwise;
// - error or nil.
// On the callers side, it should be first checked for error, then for the bool, and finally the slice.
func testResourceDataTypeDiffHandlingListExternalRead(envName string) ([]datatypes.DataType, bool, error) {
	value := oswrapper.Getenv(envName)
	log.Printf("[DEBUG] env %s value is `%s`", envName, value)
	if value != "" {
		rawDataTypes := strings.Split(value, "|")
		externalDataTypes := make([]datatypes.DataType, len(rawDataTypes))
		for i, dt := range rawDataTypes {
			externalDataType, err := datatypes.ParseDataType(dt)
			if err != nil {
				return externalDataTypes, true, err
			}
			externalDataTypes[i] = externalDataType
		}
		return externalDataTypes, true, nil
	} else {
		return nil, false, nil
	}
}

func testResourceDataTypeDiffHandlingListSet(envName string, dataTypes []datatypes.DataType) error {
	newValue := strings.Join(collections.Map(dataTypes, func(dt datatypes.DataType) string { return dt.ToSql() }), "|")
	log.Printf("[DEBUG] setting env %s to value `%s`", envName, newValue)
	return os.Setenv(envName, newValue)
}

func testResourceDataTypeDiffHandlingListUnset(envName string) error {
	log.Printf("[DEBUG] unsetting env %s", envName)
	return os.Setenv(envName, "")
}
