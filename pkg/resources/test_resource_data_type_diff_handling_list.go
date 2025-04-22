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

	dataTypes := make([]datatypes.DataType, 0)
	arguments := d.Get("nesting_list").([]any)
	// TODO [next PR or will be extracted]: extract common method to work with collections on the resource side (will be done with the masking policy and row access policy application)
	for _, arg := range arguments {
		v := arg.(map[string]any)
		dataType, err := readNestedDatatypeCommon(v, "nested_datatype")
		if err != nil {
			return diag.FromErr(err)
		}
		dataTypes = append(dataTypes, dataType)
	}
	err := testResourceDataTypeDiffHandlingListSet(envName, dataTypes)
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
		if value != "" {
			// TODO [next PR]: extract splitting
			dataTypes := strings.Split(value, "|")
			currentConfigDatatypes := d.Get("nesting_list").([]any)
			nestedDatatypesSchema := make([]map[string]any, 0)
			for i, dt := range dataTypes {
				externalDataType, err := datatypes.ParseDataType(dt)
				if err != nil {
					return diag.FromErr(err)
				}
				// TODO [next PR]: handle missing data types too
				if i+1 > len(currentConfigDatatypes) {
					log.Printf("[DEBUG] reading %d: external datatype %s outside of original range, adding", i, SqlNew(externalDataType))
					nestedDatatypesSchema = append(nestedDatatypesSchema, map[string]any{
						"nested_datatype": SqlNew(externalDataType),
					})
				} else {
					v := currentConfigDatatypes[i].(map[string]any)
					currentConfigDataType, err := readNestedDatatypeCommon(v, "nested_datatype")
					if err != nil {
						return diag.FromErr(err)
					}
					// current config data type is saved to state with all attributes known
					// external data type is left without changes as all the unknowns should remain as unknowns
					if datatypes.AreDefinitelyDifferent(currentConfigDataType, externalDataType) {
						log.Printf("[DEBUG] reading %d: external datatype %s is definitely different from the current config %s, updating", i, SqlNew(externalDataType), SqlNew(currentConfigDataType))
						nestedDatatypesSchema = append(nestedDatatypesSchema, map[string]any{
							"nested_datatype": SqlNew(externalDataType),
						})
					} else {
						log.Printf("[DEBUG] reading %d: external datatype %s is not definitely different from the current config %s, not updating", i, SqlNew(externalDataType), SqlNew(currentConfigDataType))
						nestedDatatypesSchema = append(nestedDatatypesSchema, map[string]any{
							// TODO [SNOW-2054238]: add test for StateFunc behavior with collections.
							// using toSql() here as StateFunc seems to be not working in this case
							"nested_datatype": currentConfigDataType.ToSql(),
						})
					}
				}
			}

			// TODO [SNOW-2054240]: test saving collections differently than through the whole pre-created map.
			if err := d.Set("nesting_list", nestedDatatypesSchema); err != nil {
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

func testResourceDataTypeDiffHandlingListSet(envName string, dataTypes []datatypes.DataType) error {
	newValue := strings.Join(collections.Map(dataTypes, func(dt datatypes.DataType) string { return dt.ToSql() }), "|")
	log.Printf("[DEBUG] setting env %s to value `%s`", envName, newValue)
	return os.Setenv(envName, newValue)
}

func testResourceDataTypeDiffHandlingListUnset(envName string) error {
	log.Printf("[DEBUG] unsetting env %s", envName)
	return os.Setenv(envName, "")
}
