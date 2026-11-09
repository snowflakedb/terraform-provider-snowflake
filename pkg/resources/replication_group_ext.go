package resources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	replicationGroupCreateExt      = PreviewFeatureCreateContextWrapper(string(previewfeatures.ReplicationGroupResource), TrackingCreateWrapper(resources.ReplicationGroup, CreateReplicationGroup))
	replicationGroupReadExt        = PreviewFeatureReadContextWrapper(string(previewfeatures.ReplicationGroupResource), TrackingReadWrapper(resources.ReplicationGroup, ReadReplicationGroupFunc(true)))
	replicationGroupUpdateExt      = PreviewFeatureUpdateContextWrapper(string(previewfeatures.ReplicationGroupResource), TrackingUpdateWrapper(resources.ReplicationGroup, UpdateReplicationGroup))
	replicationGroupDescriptionExt = "Resource used to manage replication groups in the source account. For more information, check [replication group documentation](https://docs.snowflake.com/en/sql-reference/sql/create-replication-group)."
	replicationGroupImporterExt    = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.ReplicationGroup, ImportReplicationGroup),
	}
	// Listed on the replication group SSOT object and appended to the generated customdiff.All.
	replicationGroupRecreateWhenNotPrimary = RecreateWhenResourceBoolFieldChangedExternally("is_primary", true)
)

type replicationGroupAlterRequestsExt struct {
	set            *sdk.ReplicationGroupSetRequest
	setObjectLists *sdk.ReplicationGroupSetRequest
	unset          *sdk.ReplicationGroupUnsetRequest
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

// TODO [next PR]: ApplySetUnset is generated from the definition, so this function goes away.
func replicationGroupApplySetUnsetExt(d *schema.ResourceData) (*replicationGroupAlterRequestsExt, error) {
	setObjectLists := sdk.NewReplicationGroupSetRequest()
	set := sdk.NewReplicationGroupSetRequest()
	unset := sdk.NewReplicationGroupUnsetRequest()
	errs := errors.Join(
		listValueUpdateSetOnlyBuilder(d, "object_types", setObjectLists.WithObjectTypes, toReplicationGroupObjectTypes),
		listValueUpdateSetOnlyBuilder(d, "allowed_databases", setObjectLists.WithAllowedDatabases, toReplicationGroupAllowedDatabases),
		listValueUpdateSetOnlyBuilder(d, "allowed_external_volumes", setObjectLists.WithAllowedExternalVolumes, toReplicationGroupAllowedExternalVolumes),
		listValueUpdateSetOnlyBuilder(d, "allowed_shares", setObjectLists.WithAllowedShares, toReplicationGroupAllowedShares),
		listValueUpdateSetOnlyBuilder(d, "allowed_integration_types", setObjectLists.WithAllowedIntegrationTypes, toReplicationGroupAllowedIntegrationTypes),
		stringAttributeUpdateBuilder(d, "comment", set.WithComment, unset.WithComment),
		stringAttributeUpdateBuilder(d, "replication_schedule", set.WithReplicationSchedule, unset.WithReplicationSchedule),
		accountObjectIdentifierAttributeUpdate(d, "error_integration", &set.ErrorIntegration, &unset.ErrorIntegration),
	)
	if errs != nil {
		return nil, errs
	}
	return &replicationGroupAlterRequestsExt{
		set:            set,
		setObjectLists: setObjectLists,
		unset:          unset,
	}, nil
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
		d.Set("replication_schedule", group.ReplicationSchedule),
		d.Set("allowed_integration_types", replicationGroupEnumStrings(group.AllowedIntegrationTypes)),
	}
	return errors.Join(errs...)
}

func replicationGroupSetFieldsNotSetByReadExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier, _ *sdk.ReplicationGroup, _ bool) error {
	return d.Set("name", id.Name())
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

func accountIdentifierNames(values []sdk.AccountIdentifier) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = value.FullyQualifiedName()
	}
	return out
}
