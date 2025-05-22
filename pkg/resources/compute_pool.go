package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var computePoolSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the compute pool; must be unique for the account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"for_application": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "Specifies the Snowflake Native App name.",
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"min_nodes": {
		Type:             schema.TypeInt,
		Required:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the minimum number of nodes for the compute pool.",
	},
	"max_nodes": {
		Type:             schema.TypeInt,
		Required:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
		Description:      "Specifies the maximum number of nodes for the compute pool.",
	},
	"instance_family": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToComputePoolInstanceFamily),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToComputePoolInstanceFamily)),
		Description:      fmt.Sprintf("Identifies the type of machine you want to provision for the nodes in the compute pool. Valid values are (case-insensitive): %s.", possibleValuesListed(sdk.AllComputePoolInstanceFamilies)),
	},
	"auto_resume": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_resume"),
		Description:      booleanStringFieldDescription("Specifies whether to automatically resume a compute pool when a service or job is submitted to it."),
		Default:          BooleanDefault,
	},
	"initially_suspended": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateBooleanString,
		DiffSuppressFunc: IgnoreAfterCreation,
		Description:      "Specifies whether the compute pool is created initially in the suspended state. This field is used only when creating a compute pool. Changes on this field are ignored after creation.",
		Default:          BooleanDefault,
	},
	"auto_suspend_secs": {
		Type:             schema.TypeInt,
		Optional:         true,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInShow("auto_suspend_secs"),
		Description:      "Number of seconds of inactivity after which you want Snowflake to automatically suspend the compute pool.",
		Default:          IntDefault,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the compute pool.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW COMPUTE POOLS` for the given compute pool.",
		Elem: &schema.Resource{
			Schema: schemas.ShowComputePoolSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE COMPUTE POOL` for the given compute pool.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeComputePoolSchema,
		},
	},
}

func ComputePool() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.ComputePools.DropSafely
		},
	)
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.ComputePoolResource), TrackingCreateWrapper(resources.ComputePool, CreateComputePool)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.ComputePoolResource), TrackingReadWrapper(resources.ComputePool, ReadComputePoolFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.ComputePoolResource), TrackingUpdateWrapper(resources.ComputePool, UpdateComputePool)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.ComputePoolResource), TrackingDeleteWrapper(resources.ComputePool, deleteFunc)),
		Description:   "Resource used to manage compute pools. For more information, check [compute pools documentation](https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.ComputePool, customdiff.All(
			ComputedIfAnyAttributeChanged(computePoolSchema, ShowOutputAttributeName, "auto_suspend_secs", "auto_resume", "min_nodes", "max_nodes", "comment"),
			ComputedIfAnyAttributeChanged(computePoolSchema, DescribeOutputAttributeName, "auto_suspend_secs", "auto_resume", "min_nodes", "max_nodes", "comment"),
		)),

		Schema: computePoolSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.ComputePool, ImportComputePool),
		},

		Timeouts: defaultTimeouts,
	}
}

func ImportComputePool(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeID(d.Id()).(sdk.AccountObjectIdentifier)

	computePool, err := client.ComputePools.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if computePool.Application != nil {
		if err := d.Set("for_application", computePool.Application.FullyQualifiedName()); err != nil {
			return nil, err
		}
	}

	errs := errors.Join(
		d.Set("name", computePool.Name),
		d.Set("auto_resume", booleanStringFromBool(computePool.AutoResume)),
		d.Set("auto_suspend_secs", computePool.AutoSuspendSecs),
	)
	if errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

func CreateComputePool(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	id := sdk.NewAccountObjectIdentifier(name)
	minNodes := d.Get("min_nodes").(int)
	maxNodes := d.Get("max_nodes").(int)
	instanceFamilyRaw := d.Get("instance_family").(string)
	instanceFamily, err := sdk.ToComputePoolInstanceFamily(instanceFamilyRaw)
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateComputePoolRequest(id, minNodes, maxNodes, instanceFamily)
	errs := errors.Join(
		accountObjectIdentifierAttributeCreate(d, "for_application", &request.ForApplication),
		booleanStringAttributeCreateBuilder(d, "auto_resume", request.WithAutoResume),
		booleanStringAttributeCreateBuilder(d, "initially_suspended", request.WithInitiallySuspended),
		intAttributeWithSpecialDefaultCreateBuilder(d, "auto_suspend_secs", request.WithAutoSuspendSecs),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	if err := client.ComputePools.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadComputePoolFunc(false)(ctx, d, meta)
}

func ReadComputePoolFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseAccountObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		computePool, err := client.ComputePools.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query compute pool. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Compute pool id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}
		computePoolDetails, err := client.ComputePools.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}
		if withExternalChangesMarking {
			var applicationFullyQualifiedName string
			if computePool.Application != nil {
				applicationFullyQualifiedName = computePool.Application.FullyQualifiedName()
			}
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"application", "for_application", applicationFullyQualifiedName, applicationFullyQualifiedName, nil},
				outputMapping{"auto_resume", "auto_resume", computePool.AutoResume, booleanStringFromBool(computePool.AutoResume), nil},
				outputMapping{"auto_suspend_secs", "auto_suspend_secs", computePool.AutoSuspendSecs, computePool.AutoSuspendSecs, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, computePoolSchema, []string{
			"for_application",
			"auto_resume",
			"auto_suspend_secs",
		}); err != nil {
			return diag.FromErr(err)
		}
		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.ComputePoolToSchema(computePool)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.ComputePoolDetailsToSchema(*computePoolDetails)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("min_nodes", computePool.MinNodes),
			d.Set("max_nodes", computePool.MaxNodes),
			d.Set("instance_family", computePool.InstanceFamily),
			d.Set("comment", computePool.Comment),
		)
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateComputePool(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	set, unset := sdk.NewComputePoolSetRequest(), sdk.NewComputePoolUnsetRequest()
	errs := errors.Join(
		// name, application, instance_family are handled by ForceNew.
		// initially_suspended is ignored after creation.
		intAttributeUpdateSetOnly(d, "min_nodes", &set.MinNodes),
		intAttributeUpdateSetOnly(d, "max_nodes", &set.MaxNodes),
		intAttributeWithSpecialDefaultUpdate(d, "auto_suspend_secs", &set.AutoSuspendSecs, &unset.AutoSuspendSecs),
		booleanStringAttributeUpdate(d, "auto_resume", &set.AutoResume, &unset.AutoResume),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}
	if (*set != sdk.ComputePoolSetRequest{}) {
		if err := client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if (*unset != sdk.ComputePoolUnsetRequest{}) {
		if err := client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}
	return ReadComputePoolFunc(false)(ctx, d, meta)
}
