package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var warehouseAdaptiveSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Identifier for the adaptive warehouse; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies a comment for the adaptive warehouse.",
	},
	"max_query_performance_level": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToMaxQueryPerformanceLevel),
		DiffSuppressFunc: NormalizeAndCompare(sdk.ToMaxQueryPerformanceLevel),
		Description:      fmt.Sprintf("Specifies the maximum query performance level for the adaptive warehouse. Determines the initial compute capacity. Can only be set at creation time. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllMaxQueryPerformanceLevels)),
	},
	"query_throughput_multiplier": {
		Type:             schema.TypeInt,
		Optional:         true,
		ForceNew:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Specifies the query throughput multiplier for the adaptive warehouse.",
	},
	strings.ToLower(string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds)): {
		Type:             schema.TypeInt,
		Optional:         true,
		Computed:         true,
		ForceNew:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		Description:      "Object parameter that specifies the time, in seconds, a SQL statement (query, DDL, DML, etc.) can be queued on a warehouse before it is canceled by the system.",
	},
	strings.ToLower(string(sdk.WarehouseParameterStatementTimeoutInSeconds)): {
		Type:             schema.TypeInt,
		Optional:         true,
		Computed:         true,
		ForceNew:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 604800)),
		Description:      "Specifies the time, in seconds, after which a running SQL statement (query, DDL, DML, etc.) is canceled by the system.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW WAREHOUSES` for the given adaptive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowAdaptiveWarehouseSchema,
		},
	},
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN WAREHOUSE` for the given adaptive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseAdaptiveParametersSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func WarehouseAdaptive() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.Warehouses.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.WarehouseAdaptive, CreateAdaptiveWarehouse),
		ReadContext:   TrackingReadWrapper(resources.WarehouseAdaptive, ReadAdaptiveWarehouseFunc(true)),
		// TODO: uncomment after ALTER is cleared.
		// UpdateContext: TrackingUpdateWrapper(resources.WarehouseAdaptive, UpdateAdaptiveWarehouse),
		DeleteContext: TrackingDeleteWrapper(resources.WarehouseAdaptive, deleteFunc),
		Description:   "Resource used to manage adaptive warehouse objects. Adaptive warehouses automatically scale compute resources based on workload. For more information, check [adaptive warehouse documentation](https://docs.snowflake.com/en/LIMITEDACCESS/adaptive-warehouses).",

		Schema: warehouseAdaptiveSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.WarehouseAdaptive, ImportAdaptiveWarehouse),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.WarehouseAdaptive, customdiff.All(
			// TODO: uncomment after ALTER is cleared.
			// ComputedIfAnyAttributeChanged(warehouseAdaptiveSchema, ShowOutputAttributeName, "name", "comment", "max_query_performance_level", "query_throughput_multiplier"),
			// ComputedIfAnyAttributeChanged(warehouseAdaptiveSchema, ParametersAttributeName,
			// 	strings.ToLower(string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds)),
			// 	strings.ToLower(string(sdk.WarehouseParameterStatementTimeoutInSeconds)),
			// ),
			// ComputedIfAnyAttributeChanged(warehouseAdaptiveSchema, FullyQualifiedNameAttributeName, "name"),
			ParametersCustomDiff(
				warehouseParametersProvider,
				parameter[sdk.AccountParameter]{sdk.AccountParameterStatementQueuedTimeoutInSeconds, valueTypeInt, sdk.ParameterTypeWarehouse},
				parameter[sdk.AccountParameter]{sdk.AccountParameterStatementTimeoutInSeconds, valueTypeInt, sdk.ParameterTypeWarehouse},
			),
		)),
		Timeouts: defaultTimeouts,
	}
}

func ImportAdaptiveWarehouse(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	w, err := client.Warehouses.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if w.Type != sdk.WarehouseTypeAdaptive {
		return nil, fmt.Errorf("warehouse %s is not of type ADAPTIVE, got %s; use snowflake_warehouse instead", id.FullyQualifiedName(), w.Type)
	}

	errs := errors.Join(
		d.Set("name", id.Name()),
		d.Set("comment", w.Comment),
		setOptionalFromPtr(d, "max_query_performance_level", w.MaxQueryPerformanceLevel),
		setOptionalFromPtr(d, "query_throughput_multiplier", w.QueryThroughputMultiplier),
	)
	if err = errs; err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func CreateAdaptiveWarehouse(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	name := d.Get("name").(string)
	id, err := sdk.ParseAccountObjectIdentifier(name)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := &sdk.CreateAdaptiveWarehouseOptions{}
	errs := errors.Join(
		stringAttributeCreate(d, "comment", &opts.Comment),
		attributeMappedValueCreate(d, "max_query_performance_level", &opts.MaxQueryPerformanceLevel, func(v any) (*sdk.MaxQueryPerformanceLevel, error) {
			s, err := sdk.ToMaxQueryPerformanceLevel(v.(string))
			if err != nil {
				return nil, err
			}
			return &s, nil
		}),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	opts.QueryThroughputMultiplier = GetConfigPropertyAsPointerAllowingZeroValue[int](d, "query_throughput_multiplier")

	opts.StatementQueuedTimeoutInSeconds = GetConfigPropertyAsPointerAllowingZeroValue[int](d, "statement_queued_timeout_in_seconds")
	opts.StatementTimeoutInSeconds = GetConfigPropertyAsPointerAllowingZeroValue[int](d, "statement_timeout_in_seconds")

	if err := client.Warehouses.CreateAdaptive(ctx, id, opts); err != nil {
		return diag.FromErr(fmt.Errorf("error creating adaptive warehouse %s: %w", id.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadAdaptiveWarehouseFunc(false)(ctx, d, meta)
}

func ReadAdaptiveWarehouseFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		w, err := client.Warehouses.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					{
						Severity: diag.Warning,
						Summary:  "Failed to query adaptive warehouse. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Adaptive warehouse id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		warehouseParameters, err := client.Warehouses.ShowParameters(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			var maxQueryPerformanceLevel string
			if w.MaxQueryPerformanceLevel != nil {
				maxQueryPerformanceLevel = string(*w.MaxQueryPerformanceLevel)
			}
			var queryThroughputMultiplier int
			if w.QueryThroughputMultiplier != nil {
				queryThroughputMultiplier = *w.QueryThroughputMultiplier
			}
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"max_query_performance_level", "max_query_performance_level", maxQueryPerformanceLevel, maxQueryPerformanceLevel, nil},
				outputMapping{"query_throughput_multiplier", "query_throughput_multiplier", queryThroughputMultiplier, queryThroughputMultiplier, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		providerCtx := meta.(*provider.Context)
		errs := errors.Join(
			d.Set("name", w.Name),
			d.Set("comment", w.Comment),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.WarehouseToSchema(w)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ParametersAttributeName, []map[string]any{schemas.WarehouseAdaptiveParametersToSchema(warehouseParameters, providerCtx)}),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}

		for _, parameter := range warehouseParameters {
			switch parameter.Key {
			case
				string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds),
				string(sdk.WarehouseParameterStatementTimeoutInSeconds):
				value, err := strconv.Atoi(parameter.Value)
				if err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set(strings.ToLower(parameter.Key), value); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		return nil
	}
}

func UpdateAdaptiveWarehouse(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewAccountObjectIdentifier(d.Get("name").(string))
		if err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{NewName: &newId}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set := sdk.WarehouseSet{}
	unset := sdk.WarehouseUnset{}

	if err := stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := intAttributeUpdate(d, "query_throughput_multiplier", &set.QueryThroughputMultiplier, &unset.QueryThroughputMultiplier); err != nil {
		return diag.FromErr(err)
	}

	if diags := JoinDiags(
		handleParameterUpdate(d, sdk.WarehouseParameterStatementQueuedTimeoutInSeconds, &set.StatementQueuedTimeoutInSeconds, &unset.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, sdk.WarehouseParameterStatementTimeoutInSeconds, &set.StatementTimeoutInSeconds, &unset.StatementTimeoutInSeconds),
	); diags != nil {
		return diags
	}

	if (set != sdk.WarehouseSet{}) {
		if err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{Set: &set}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting adaptive warehouse properties: %w", err))
		}
	}

	if (unset != sdk.WarehouseUnset{}) {
		if err := client.Warehouses.Alter(ctx, id, &sdk.AlterWarehouseOptions{Unset: &unset}); err != nil {
			return diag.FromErr(fmt.Errorf("error unsetting adaptive warehouse properties: %w", err))
		}
	}

	return ReadAdaptiveWarehouseFunc(false)(ctx, d, meta)
}
