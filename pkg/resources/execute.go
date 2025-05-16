package resources

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var executeSchema = map[string]*schema.Schema{
	"execute": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "SQL statement to execute. Forces recreation of resource on change.",
	},
	"revert": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "SQL statement to revert the execute statement. Invoked when resource is being destroyed.",
	},
	"query": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Optional SQL statement to do a read. Invoked on every resource refresh and every time it is changed.",
	},
	"query_results": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of key-value maps (text to text) retrieved after executing read query. Will be empty if the query results in an error.",
		Elem: &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
}

func Execute() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.Execute, CreateExecute),
		ReadContext:   TrackingReadWrapper(resources.Execute, ReadExecute),
		UpdateContext: TrackingUpdateWrapper(resources.Execute, UpdateExecute),
		DeleteContext: TrackingDeleteWrapper(resources.Execute, DeleteExecute),

		Schema: executeSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Description: "Resource allowing execution of ANY SQL statement.",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Execute, customdiff.All(
			customdiff.ForceNewIfChange("execute", func(ctx context.Context, oldValue, newValue, meta any) bool {
				return oldValue != ""
			}),
			func(_ context.Context, diff *schema.ResourceDiff, _ any) error {
				if diff.HasChange("query") {
					err := diff.SetNewComputed("query_results")
					if err != nil {
						return err
					}
				}
				return nil
			}),
		),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},
	}
}

func CreateExecute(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	executeStatement := d.Get("execute").(string)
	_, err = client.ExecUnsafe(ctx, executeStatement)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return ReadExecute(ctx, d, meta)
}

func UpdateExecute(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if d.HasChange("query") {
		return ReadExecute(ctx, d, meta)
	}
	return nil
}

func ReadExecute(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	readStatement := d.Get("query").(string)

	setNilResults := func() diag.Diagnostics {
		log.Printf(`[DEBUG] Clearing query_results`)
		err := d.Set("query_results", nil)
		if err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	if readStatement == "" {
		return setNilResults()
	} else {
		rows, err := client.QueryUnsafe(ctx, readStatement)
		if err != nil {
			return append(setNilResults(), diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "SQL query statement failed",
				Detail:   fmt.Sprintf("Err: %v", err),
			})
		}
		log.Printf(`[INFO] SQL query statement executed successfully, returned rows count: %d`, len(rows))
		rowsTransformed := make([]map[string]any, len(rows))
		for i, row := range rows {
			t := make(map[string]any)
			for k, v := range row {
				if *v == nil {
					t[k] = nil
				} else {
					switch (*v).(type) {
					case fmt.Stringer:
						t[k] = fmt.Sprintf("%v", *v)
					case string:
						t[k] = *v
					default:
						return diag.FromErr(fmt.Errorf("currently only objects convertible to String are supported by query; got %v", *v))
					}
				}
			}
			rowsTransformed[i] = t
		}
		err = d.Set("query_results", rowsTransformed)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func DeleteExecute(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	revertStatement := d.Get("revert").(string)
	_, err := client.ExecUnsafe(ctx, revertStatement)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] Revert SQL applied successfully")

	return nil
}
