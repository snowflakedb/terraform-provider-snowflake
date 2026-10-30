package resources

import (
	"context"
	"errors"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
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
		StateContext: TrackingImportWrapper(resources.ReplicationGroup, importReplicationGroup),
	}
)

type replicationGroupAlterRequestsExt struct {
	objectLists *sdk.ReplicationGroupSetRequest
	properties  *sdk.ReplicationGroupSetRequest
	unset       *sdk.ReplicationGroupUnsetRequest
}

func replicationGroupApplyCreateOptionalsExt(d *schema.ResourceData, request *sdk.CreateReplicationGroupRequest) error {
	// CREATE has no COMMENT clause. Comment is set with ALTER.
	return errors.Join(
		listValueCreate(d, "object_types", &request.ObjectTypes, toReplicationGroupObjectTypes),
		listValueCreate(d, "allowed_databases", &request.AllowedDatabases, toReplicationGroupAllowedDatabases),
		listValueCreate(d, "allowed_external_volumes", &request.AllowedExternalVolumes, toReplicationGroupAllowedExternalVolumes),
		listValueCreate(d, "allowed_shares", &request.AllowedShares, toReplicationGroupAllowedShares),
		listValueCreate(d, "allowed_integration_types", &request.AllowedIntegrationTypes, toReplicationGroupAllowedIntegrationTypes),
		setMappedValueCreateBuilder(d, "allowed_accounts", func(values []sdk.AccountIdentifier) *sdk.CreateReplicationGroupRequest {
			request.AllowedAccounts = values
			return request
		}, toReplicationGroupAllowedAccounts),
		booleanStringAttributeCreateBuilder(d, "ignore_edition_check", request.WithIgnoreEditionCheck),
		stringAttributeCreateBuilder(d, "replication_schedule", request.WithReplicationSchedule),
		accountObjectIdentifierAttributeCreateBuilder(d, "error_integration", request.WithErrorIntegration),
	)
}

func replicationGroupNewCreateRequestExt(_ *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreateReplicationGroupRequest, error) {
	// OBJECT_TYPES and ALLOWED_ACCOUNTS are constructor arguments. ApplyCreateOptionals fills them.
	return sdk.NewCreateReplicationGroupRequest(id, nil, nil), nil
}

func replicationGroupApplySetUnsetExt(d *schema.ResourceData) (*replicationGroupAlterRequestsExt, error) {
	objectLists := sdk.NewReplicationGroupSetRequest()
	properties := sdk.NewReplicationGroupSetRequest()
	unset := sdk.NewReplicationGroupUnsetRequest()
	errs := errors.Join(
		listValueUpdateSetOnlyBuilder(d, "object_types", objectLists.WithObjectTypes, toReplicationGroupObjectTypes),
		listValueUpdateSetOnlyBuilder(d, "allowed_databases", objectLists.WithAllowedDatabases, toReplicationGroupAllowedDatabases),
		listValueUpdateSetOnlyBuilder(d, "allowed_external_volumes", objectLists.WithAllowedExternalVolumes, toReplicationGroupAllowedExternalVolumes),
		listValueUpdateSetOnlyBuilder(d, "allowed_shares", objectLists.WithAllowedShares, toReplicationGroupAllowedShares),
		listValueUpdateSetOnlyBuilder(d, "allowed_integration_types", objectLists.WithAllowedIntegrationTypes, toReplicationGroupAllowedIntegrationTypes),
		stringAttributeUpdateBuilder(d, "comment", properties.WithComment, unset.WithComment),
		stringAttributeUpdateBuilder(d, "replication_schedule", properties.WithReplicationSchedule, unset.WithReplicationSchedule),
		accountObjectIdentifierAttributeUpdate(d, "error_integration", &properties.ErrorIntegration, &unset.ErrorIntegration),
	)
	if errs != nil {
		return nil, errs
	}
	return &replicationGroupAlterRequestsExt{
		objectLists: objectLists,
		properties:  properties,
		unset:       unset,
	}, nil
}

func replicationGroupAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, setUnset *replicationGroupAlterRequestsExt) error {
	client := meta.(*provider.Context).Client
	if !reflect.DeepEqual(*setUnset.objectLists, *sdk.NewReplicationGroupSetRequest()) {
		if err := client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(id).WithSet(*setUnset.objectLists)); err != nil {
			return err
		}
	}
	if !reflect.DeepEqual(*setUnset.properties, *sdk.NewReplicationGroupSetRequest()) {
		if err := client.ReplicationGroups.Alter(ctx, sdk.NewAlterReplicationGroupRequest(id).WithSet(*setUnset.properties)); err != nil {
			return err
		}
	}
	return nil
}

func replicationGroupSetConfigFieldsExt(ctx context.Context, d *schema.ResourceData, meta any, group *sdk.ReplicationGroup) error {
	allowedAccounts, err := replicationGroupAllowedAccountsWithoutCurrent(ctx, meta, group.AllowedAccounts)
	if err != nil {
		return err
	}
	errs := []error{
		d.Set("object_types", replicationGroupEnumStrings(group.ObjectTypes)),
		d.Set("allowed_accounts", accountIdentifierNames(allowedAccounts)),
		d.Set("comment", group.Comment),
	}
	// An empty integration-type list is omitted from config. Writing [] would plan a diff against null.
	if len(group.AllowedIntegrationTypes) > 0 {
		errs = append(errs, d.Set("allowed_integration_types", replicationGroupEnumStrings(group.AllowedIntegrationTypes)))
	}
	return errors.Join(errs...)
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

func importReplicationGroup(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	result, err := ImportReplicationGroup(ctx, d, meta)
	if err != nil {
		return nil, err
	}
	if err := replicationGroupOmitCurrentAccountFromState(ctx, d, meta); err != nil {
		return nil, err
	}
	return result, nil
}

func replicationGroupCurrentAccount(ctx context.Context, meta any) (sdk.AccountIdentifier, error) {
	details, err := meta.(*provider.Context).Client.ContextFunctions.CurrentSessionDetails(ctx)
	if err != nil {
		return sdk.AccountIdentifier{}, err
	}
	return sdk.NewAccountIdentifier(details.OrganizationName, details.AccountName), nil
}

func replicationGroupAllowedAccountsWithoutCurrent(ctx context.Context, meta any, accounts []sdk.AccountIdentifier) ([]sdk.AccountIdentifier, error) {
	current, err := replicationGroupCurrentAccount(ctx, meta)
	if err != nil {
		return nil, err
	}
	return omitCurrentAccount(accounts, current), nil
}

func omitCurrentAccount(accounts []sdk.AccountIdentifier, current sdk.AccountIdentifier) []sdk.AccountIdentifier {
	filtered := make([]sdk.AccountIdentifier, 0, len(accounts))
	for _, account := range accounts {
		if account.FullyQualifiedName() == current.FullyQualifiedName() {
			continue
		}
		filtered = append(filtered, account)
	}
	return filtered
}

func replicationGroupOmitCurrentAccountFromState(ctx context.Context, d *schema.ResourceData, meta any) error {
	current, err := replicationGroupCurrentAccount(ctx, meta)
	if err != nil {
		return err
	}
	raw, ok := d.GetOk("allowed_accounts")
	if !ok {
		return nil
	}
	filtered := make([]string, 0)
	for _, value := range expandStringList(raw.(*schema.Set).List()) {
		account, err := helpers.DecodeSnowflakeAccountIdentifier(value)
		if err != nil {
			return err
		}
		if account.FullyQualifiedName() == current.FullyQualifiedName() {
			continue
		}
		filtered = append(filtered, value)
	}
	return d.Set("allowed_accounts", filtered)
}

func accountIdentifierNames(values []sdk.AccountIdentifier) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = value.FullyQualifiedName()
	}
	return out
}
