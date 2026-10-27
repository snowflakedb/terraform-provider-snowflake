package resources

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	computePoolCreateExt      = TrackingCreateWrapper(resources.ComputePool, CreateComputePool)
	computePoolReadExt        = TrackingReadWrapper(resources.ComputePool, ReadComputePoolFunc(true))
	computePoolUpdateExt      = TrackingUpdateWrapper(resources.ComputePool, UpdateComputePool)
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

type computePoolAlterRequestsExt struct {
	set   *sdk.ComputePoolSetRequest
	unset *sdk.ComputePoolUnsetRequest
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

func computePoolSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, computePool *sdk.ComputePool, importPermissiveMode bool) error {
	var applicationFullyQualifiedName string
	if computePool.Application != nil {
		applicationFullyQualifiedName = computePool.Application.FullyQualifiedName()
	}
	return errors.Join(
		d.Set("name", id.Name()),
		importStringAttribute(d, "for_application", importPermissiveMode, applicationFullyQualifiedName),
		importBooleanStringAttribute(d, "auto_resume", importPermissiveMode, computePool.AutoResume),
		importIntWithIntDefaultAttribute(d, "auto_suspend_secs", importPermissiveMode, computePool.AutoSuspendSecs),
	)
}

// =============================================================================
// Pre-existing mappers
// =============================================================================

func toComputePoolBackupInstanceFamilies(value any) (sdk.ComputePoolInstanceFamilyListItem, error) {
	instanceFamily, err := sdk.ToComputePoolInstanceFamily(value.(string))
	if err != nil {
		return sdk.ComputePoolInstanceFamilyListItem{}, err
	}
	return sdk.ComputePoolInstanceFamilyListItem{Value: instanceFamily}, nil
}
