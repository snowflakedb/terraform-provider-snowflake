package resources

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	replicationGroupCreateExt      = TrackingCreateWrapper(resources.ReplicationGroup, CreateReplicationGroup)
	replicationGroupReadExt        = TrackingReadWrapper(resources.ReplicationGroup, ReadReplicationGroupFunc(true))
	replicationGroupUpdateExt      = TrackingUpdateWrapper(resources.ReplicationGroup, UpdateReplicationGroup)
	replicationGroupDescriptionExt = "Resource used to manage replication groups in the source account. For more information, check [replication group documentation](https://docs.snowflake.com/en/sql-reference/sql/create-replication-group)."
	replicationGroupImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.ReplicationGroup, ImportReplicationGroup),
	}
)

type replicationGroupAlterRequestsExt struct {
	set   *sdk.ReplicationGroupSetRequest
	unset *sdk.ReplicationGroupUnsetRequest
}

func replicationGroupApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateReplicationGroupRequest) error {
	// CREATE has no COMMENT clause. Comment is set with ALTER.
	return errors.Join(
		listValueCreate(d, "object_types", &request.ObjectTypes, toReplicationGroupObjectTypes),
		listValueCreate(d, "allowed_databases", &request.AllowedDatabases, toReplicationGroupAllowedDatabases),
		listValueCreate(d, "allowed_external_volumes", &request.AllowedExternalVolumes, toReplicationGroupAllowedExternalVolumes),
		listValueCreate(d, "allowed_shares", &request.AllowedShares, toReplicationGroupAllowedShares),
		listValueCreate(d, "allowed_integration_types", &request.AllowedIntegrationTypes, toReplicationGroupAllowedIntegrationTypes),
		listValueCreate(d, "allowed_accounts", &request.AllowedAccounts, toReplicationGroupAllowedAccounts),
		booleanStringAttributeCreateBuilder(d, "ignore_edition_check", request.WithIgnoreEditionCheck),
		stringAttributeCreateBuilder(d, "replication_schedule", request.WithReplicationSchedule),
		accountObjectIdentifierAttributeCreateBuilder(d, "error_integration", request.WithErrorIntegration),
	)
}

func replicationGroupNewCreateRequestExt(_ *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateReplicationGroupRequest, error) {
	// OBJECT_TYPES and ALLOWED_ACCOUNTS are constructor arguments. ApplyCreateOptionals fills them.
	return sdk.NewCreateReplicationGroupRequest(id, nil, nil), nil
}

func replicationGroupSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, group *sdk.ReplicationGroup, importPermissiveMode bool) error {
	comment := ""
	if group.Comment != nil {
		comment = *group.Comment
	}
	schedule := ""
	if group.ReplicationSchedule != nil {
		schedule = *group.ReplicationSchedule
	}
	return errors.Join(
		d.Set("name", id.Name()),
		d.Set("object_types", replicationGroupEnumStrings(group.ObjectTypes)),
		d.Set("allowed_integration_types", replicationGroupEnumStrings(group.AllowedIntegrationTypes)),
		d.Set("allowed_accounts", accountIdentifierNames(group.AllowedAccounts)),
		importStringAttribute(d, "comment", importPermissiveMode, comment),
		importStringAttribute(d, "replication_schedule", importPermissiveMode, schedule),
	)
}

func toReplicationGroupObjectTypes(value any) (sdk.ReplicationGroupObjectType, error) {
	return sdk.ToReplicationGroupObjectType(value.(string))
}

func toReplicationGroupAllowedIntegrationTypes(value any) (sdk.ReplicationGroupIntegrationType, error) {
	return sdk.ToReplicationGroupIntegrationType(value.(string))
}

func toReplicationGroupAllowedDatabases(value any) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(value.(string))
}

func toReplicationGroupAllowedExternalVolumes(value any) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(value.(string))
}

func toReplicationGroupAllowedShares(value any) (sdk.AccountObjectIdentifier, error) {
	return sdk.ParseAccountObjectIdentifier(value.(string))
}

func toReplicationGroupAllowedAccounts(value any) (sdk.AccountIdentifier, error) {
	return sdk.ParseAccountIdentifier(value.(string))
}

func replicationGroupEnumStrings[T ~string](values []T) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = string(value)
	}
	return out
}

func accountIdentifierNames(values []sdk.AccountIdentifier) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = value.FullyQualifiedName()
	}
	return out
}
