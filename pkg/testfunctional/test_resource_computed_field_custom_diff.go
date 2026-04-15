package testfunctional

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testfunctional/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestResourceComputedFieldCustomDiff explores whether calling SetNewComputed on a purely
// computed field from a CustomizeDiff can trigger an Update when no config-driven fields change.
// This mimics the tag resource's allowed_values_order pattern where order changes must trigger
// an Update even though allowed_values (TypeSet) reports no diff.

// ComputedFieldState is the JSON model stored in / fetched from the HTTP server.
type ComputedFieldState struct {
	Values        []string `json:"values"`
	UpdateCount   int      `json:"update_count"`
	TriggerUpdate bool     `json:"trigger_update"`
}

const ComputedFieldCustomDiffPath = "computed_field_custom_diff"

var testResourceComputedFieldCustomDiffSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"values": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Description: "A TypeSet field — reordering in config does not produce a diff " +
			"because sets are unordered.",
	},
	"computed_order": {
		Type:     schema.TypeList,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Description: "Computed field that stores the ordered list. " +
			"CustomizeDiff calls SetNewComputed on this to trigger Update.",
	},
	"update_count": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Tracks how many times Update was called.",
	},
}

func TestResourceComputedFieldCustomDiff() *schema.Resource {
	return &schema.Resource{
		CreateContext: testResourceComputedFieldCustomDiffCreate,
		UpdateContext: testResourceComputedFieldCustomDiffUpdate,
		ReadContext:   testResourceComputedFieldCustomDiffRead,
		DeleteContext: testResourceComputedFieldCustomDiffDelete,

		CustomizeDiff: func(_ context.Context, diff *schema.ResourceDiff, meta any) error {
			if diff.Id() == "" {
				return nil
			}

			providerCtx := meta.(*common.TestProviderContext)
			var state ComputedFieldState
			if err := common.Get(providerCtx.ServerUrl(), ComputedFieldCustomDiffPath, &state); err != nil {
				return err
			}

			if state.TriggerUpdate {
				log.Printf("[DEBUG] CustomizeDiff: trigger_update is set, calling SetNewComputed on computed_order")
				return diff.SetNewComputed("computed_order")
			}
			log.Printf("[DEBUG] CustomizeDiff: trigger_update is not set, skipping")
			return nil
		},

		Schema: testResourceComputedFieldCustomDiffSchema,
	}
}

func testResourceComputedFieldCustomDiffCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*common.TestProviderContext)
	name := d.Get("name").(string)
	log.Printf("[DEBUG] computed_field_custom_diff: CREATE for %s", name)

	d.SetId(name)

	values := expandStringSet(d.Get("values").(*schema.Set))
	state := ComputedFieldState{
		Values:      values,
		UpdateCount: 0,
	}
	if err := common.Post(providerCtx.ServerUrl(), ComputedFieldCustomDiffPath, state); err != nil {
		return diag.FromErr(err)
	}

	return testResourceComputedFieldCustomDiffRead(ctx, d, meta)
}

func testResourceComputedFieldCustomDiffUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*common.TestProviderContext)
	log.Printf("[DEBUG] computed_field_custom_diff: UPDATE for %s", d.Id())

	// Read current state to get update_count
	var current ComputedFieldState
	if err := common.Get(providerCtx.ServerUrl(), ComputedFieldCustomDiffPath, &current); err != nil {
		return diag.FromErr(err)
	}

	values := expandStringSet(d.Get("values").(*schema.Set))
	state := ComputedFieldState{
		Values:        values,
		UpdateCount:   current.UpdateCount + 1,
		TriggerUpdate: false, // Clear the trigger after Update handles it
	}
	if err := common.Post(providerCtx.ServerUrl(), ComputedFieldCustomDiffPath, state); err != nil {
		return diag.FromErr(err)
	}

	return testResourceComputedFieldCustomDiffRead(ctx, d, meta)
}

func testResourceComputedFieldCustomDiffRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*common.TestProviderContext)
	log.Printf("[DEBUG] computed_field_custom_diff: READ for %s", d.Id())

	var state ComputedFieldState
	if err := common.Get(providerCtx.ServerUrl(), ComputedFieldCustomDiffPath, &state); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("computed_order", state.Values); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("update_count", state.UpdateCount); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func testResourceComputedFieldCustomDiffDelete(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
	log.Printf("[DEBUG] computed_field_custom_diff: DELETE for %s", d.Id())
	d.SetId("")
	return nil
}

func expandStringSet(set *schema.Set) []string {
	items := set.List()
	result := make([]string, len(items))
	for i, v := range items {
		result[i] = v.(string)
	}
	return result
}
