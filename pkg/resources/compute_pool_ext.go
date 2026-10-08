package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	computePoolCreateExt = TrackingCreateWrapper(resources.ComputePool, CreateComputePool)
	computePoolReadExt   = TrackingReadWrapper(resources.ComputePool, ReadComputePoolFunc(true))
	computePoolUpdateExt = TrackingUpdateWrapper(resources.ComputePool, UpdateComputePool)
	computePoolDeleteExt = TrackingDeleteWrapper(resources.ComputePool, ResourceDeleteContextFunc(
		computePoolParseIdExt,
		computePoolDropSafelyInSdkExt,
	))
	computePoolDescriptionExt = joinWithSpace(
		"Resource used to manage compute pools. For more information, check [compute pools documentation](https://docs.snowflake.com/en/sql-reference/sql/create-compute-pool).",
		"A compute pool is a collection of one or more virtual machine (VM) nodes on which Snowflake runs your Snowpark Container Services services (including job services).",
		"See [Working with compute pools](https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-compute-pool) developer guide for more details.",
	)
	computePoolCustomizeDiffExt = TrackingCustomDiffWrapper(resources.ComputePool, customdiff.All(
		// For now, the list fields have to be excluded.
		// TODO [SNOW-1648997]: address the above comment
		ComputedIfAnyAttributeChanged(computePoolSchema, ShowOutputAttributeName, "auto_suspend_secs", "auto_resume", "min_nodes", "max_nodes", "comment"),
		ComputedIfAnyAttributeChanged(computePoolSchema, DescribeOutputAttributeName, "auto_suspend_secs", "auto_resume", "min_nodes", "max_nodes", "comment"),
	))
	computePoolImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.ComputePool, ImportComputePool),
	}
)

// =============================================================================
// Request / mapping / Read-Update-Import helpers
// =============================================================================

var computePoolSetStateToValueFromConfigFieldsExt = []string{
	"for_application",
	"auto_resume",
	"auto_suspend_secs",
}

type computePoolAlterRequestsExt struct {
	set   *sdk.ComputePoolSetRequest
	unset *sdk.ComputePoolUnsetRequest
}

func computePoolParseIdFromConfigExt(d *schema.ResourceData) (sdk.AccountObjectIdentifier, error) {
	return sdk.NewAccountObjectIdentifier(d.Get("name").(string)), nil
}

func computePoolParseIdExt(id string) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(id)
}

func computePoolNewCreateRequestExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateComputePoolRequest, error) {
	minNodes := d.Get("min_nodes").(int)
	maxNodes := d.Get("max_nodes").(int)
	instanceFamily, err := sdk.ToComputePoolInstanceFamily(d.Get("instance_family").(string))
	if err != nil {
		return nil, err
	}
	return sdk.NewCreateComputePoolRequest(id, minNodes, maxNodes, instanceFamily), nil
}

func computePoolApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateComputePoolRequest) error {
	return errors.Join(
		accountObjectIdentifierAttributeCreate(d, "for_application", &request.ForApplication),
		booleanStringAttributeCreateBuilder(d, "auto_resume", request.WithAutoResume),
		booleanStringAttributeCreateBuilder(d, "initially_suspended", request.WithInitiallySuspended),
		intAttributeWithSpecialDefaultCreateBuilder(d, "auto_suspend_secs", request.WithAutoSuspendSecs),
		attributeMappedValueCreateBuilder(d, "backup_instance_families", request.WithBackupInstanceFamilies, toComputePoolBackupInstanceFamilies),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	)
}

func computePoolHandleExternalChangesExt(d *schema.ResourceData, computePool *sdk.ComputePool, withExternalChangesMarking bool) error {
	if !withExternalChangesMarking {
		return nil
	}
	return handleExternalChangesToObjectInShow(d, computePoolOutputMappingsExt(computePool)...)
}

func computePoolOutputMappingsExt(computePool *sdk.ComputePool) []outputMapping {
	var applicationFullyQualifiedName string
	if computePool.Application != nil {
		applicationFullyQualifiedName = computePool.Application.FullyQualifiedName()
	}
	return []outputMapping{
		{"application", "for_application", applicationFullyQualifiedName, applicationFullyQualifiedName, nil},
		{"auto_resume", "auto_resume", computePool.AutoResume, booleanStringFromBool(computePool.AutoResume), nil},
		{"auto_suspend_secs", "auto_suspend_secs", computePool.AutoSuspendSecs, computePool.AutoSuspendSecs, nil},
	}
}

func computePoolSetConfigFieldsExt(d *schema.ResourceData, computePool *sdk.ComputePool) error {
	return errors.Join(
		d.Set("min_nodes", computePool.MinNodes),
		d.Set("max_nodes", computePool.MaxNodes),
		d.Set("instance_family", computePool.InstanceFamily),
		d.Set("backup_instance_families", computePool.BackupInstanceFamilies),
		d.Set("comment", computePool.Comment),
	)
}

func computePoolApplySetUnsetExt(d *schema.ResourceData) (*computePoolAlterRequestsExt, error) {
	set, unset := sdk.NewComputePoolSetRequest(), sdk.NewComputePoolUnsetRequest()
	errs := errors.Join(
		// name, application, instance_family are handled by ForceNew.
		// initially_suspended is ignored after creation.
		intAttributeUpdateSetOnly(d, "min_nodes", &set.MinNodes),
		intAttributeUpdateSetOnly(d, "max_nodes", &set.MaxNodes),
		intAttributeWithSpecialDefaultUpdate(d, "auto_suspend_secs", &set.AutoSuspendSecs, &unset.AutoSuspendSecs),
		booleanStringAttributeUpdate(d, "auto_resume", &set.AutoResume, &unset.AutoResume),
		listValueUpdate(d, "backup_instance_families", &set.BackupInstanceFamilies, &unset.BackupInstanceFamilies, toComputePoolBackupInstanceFamily),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	)
	if errs != nil {
		return nil, errs
	}
	return &computePoolAlterRequestsExt{set: set, unset: unset}, nil
}

func computePoolSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, computePool *sdk.ComputePool) error {
	if computePool.Application != nil {
		if err := d.Set("for_application", computePool.Application.FullyQualifiedName()); err != nil {
			return err
		}
	}
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("auto_resume", booleanStringFromBool(computePool.AutoResume)),
		d.Set("auto_suspend_secs", computePool.AutoSuspendSecs),
	)
}

// =============================================================================
// SDK incision points
// =============================================================================

func computePoolCreateInSdkExt(ctx context.Context, meta any, _ sdk.AccountObjectIdentifier, req *sdk.CreateComputePoolRequest) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	if err := client.ComputePools.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func computePoolShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.ComputePool, diag.Diagnostics) {
	client := meta.(*provider.Context).Client
	computePool, err := client.ComputePools.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query compute pool. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Compute pool id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	return computePool, nil
}

func computePoolShowByIdInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.ComputePool, error) {
	client := meta.(*provider.Context).Client
	return client.ComputePools.ShowByID(ctx, id)
}

func computePoolDescribeInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier) (*sdk.ComputePoolDetails, error) {
	client := meta.(*provider.Context).Client
	return client.ComputePools.Describe(ctx, id)
}

func computePoolAlterSetInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier, reqs *computePoolAlterRequestsExt) error {
	if reflect.DeepEqual(*reqs.set, *sdk.NewComputePoolSetRequest()) {
		return nil
	}
	client := meta.(*provider.Context).Client
	return client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithSet(*reqs.set))
}

func computePoolAlterUnsetInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier, reqs *computePoolAlterRequestsExt) error {
	if reflect.DeepEqual(*reqs.unset, *sdk.NewComputePoolUnsetRequest()) {
		return nil
	}
	client := meta.(*provider.Context).Client
	return client.ComputePools.Alter(ctx, sdk.NewAlterComputePoolRequest(id).WithUnset(*reqs.unset))
}

func computePoolDropSafelyInSdkExt(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
	return client.ComputePools.DropSafely
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

func toComputePoolBackupInstanceFamily(value any) (sdk.ComputePoolBackupInstanceFamilyListItem, error) {
	instanceFamily, err := sdk.ToComputePoolInstanceFamily(value.(string))
	if err != nil {
		return sdk.ComputePoolBackupInstanceFamilyListItem{}, err
	}
	return sdk.ComputePoolBackupInstanceFamilyListItem{Value: instanceFamily}, nil
}

func toComputePoolBackupInstanceFamilies(values []any) ([]sdk.ComputePoolBackupInstanceFamilyListItem, error) {
	return collections.MapErr(values, toComputePoolBackupInstanceFamily)
}
