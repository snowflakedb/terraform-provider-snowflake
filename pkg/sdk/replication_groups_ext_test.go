package sdk

func init() {
	id := replicationGroupsTestIdAccountObjectIdentifier
	allowedAccount := NewAccountIdentifier("MY_ORG", "MY_ACCOUNT")

	replicationGroupsTests.Create.
		withDefaultOpts(func() *CreateReplicationGroupOptions {
			return &CreateReplicationGroupOptions{
				name:            id,
				ObjectTypes:     []ReplicationGroupObjectType{ReplicationGroupObjectTypeUsers},
				AllowedAccounts: []AccountIdentifier{allowedAccount},
			}
		}).
		withExpectedSqlf(
			case_ReplicationGroups_sql_Create_basic,
			`CREATE REPLICATION GROUP %s OBJECT_TYPES = USERS ALLOWED_ACCOUNTS = "MY_ORG"."MY_ACCOUNT"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_Create_all,
			func(opts *CreateReplicationGroupOptions) {
				opts.IfNotExists = new(true)
				opts.ObjectTypes = []ReplicationGroupObjectType{
					ReplicationGroupObjectTypeAccountParameters,
					ReplicationGroupObjectTypeDatabases,
					ReplicationGroupObjectTypeExternalVolumes,
					ReplicationGroupObjectTypeIntegrations,
					ReplicationGroupObjectTypeNetworkPolicies,
					ReplicationGroupObjectTypeResourceMonitors,
					ReplicationGroupObjectTypeRoles,
					ReplicationGroupObjectTypeShares,
					ReplicationGroupObjectTypeUsers,
					ReplicationGroupObjectTypeWarehouses,
				}
				opts.AllowedDatabases = []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")}
				opts.AllowedExternalVolumes = []AccountObjectIdentifier{NewAccountObjectIdentifier("vol1")}
				opts.AllowedShares = []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")}
				opts.AllowedIntegrationTypes = []ReplicationGroupIntegrationType{
					ReplicationGroupIntegrationTypeSecurityIntegrations,
					ReplicationGroupIntegrationTypeApiIntegrations,
					ReplicationGroupIntegrationTypeStorageIntegrations,
					ReplicationGroupIntegrationTypeExternalAccessIntegrations,
					ReplicationGroupIntegrationTypeNotificationIntegrations,
				}
				opts.IgnoreEditionCheck = new(true)
				opts.ReplicationSchedule = new("10 MINUTE")
				opts.Tag = []TagAssociation{
					{Name: NewAccountObjectIdentifier("tag1"), Value: "value1"},
					{Name: NewAccountObjectIdentifier("tag2"), Value: "value2"},
				}
				opts.ErrorIntegration = new(NewAccountObjectIdentifier("error_int"))
			},
			`CREATE REPLICATION GROUP IF NOT EXISTS %s OBJECT_TYPES = ACCOUNT PARAMETERS, DATABASES, EXTERNAL VOLUMES, INTEGRATIONS, NETWORK POLICIES, RESOURCE MONITORS, ROLES, SHARES, USERS, WAREHOUSES ALLOWED_DATABASES = "db1" ALLOWED_EXTERNAL_VOLUMES = "vol1" ALLOWED_SHARES = "share1" ALLOWED_INTEGRATION_TYPES = SECURITY INTEGRATIONS, API INTEGRATIONS, STORAGE INTEGRATIONS, EXTERNAL ACCESS INTEGRATIONS, NOTIFICATION INTEGRATIONS ALLOWED_ACCOUNTS = "MY_ORG"."MY_ACCOUNT" IGNORE EDITION CHECK REPLICATION_SCHEDULE = '10 MINUTE' TAG ("tag1" = 'value1', "tag2" = 'value2') ERROR_INTEGRATION = "error_int"`,
			id.FullyQualifiedName(),
		).
		withAdditionalValidationCase(
			"validation_Create_ObjectTypes_required",
			func(opts *CreateReplicationGroupOptions) { opts.ObjectTypes = nil },
			NewError("OBJECT_TYPES must be set"),
		).
		withAdditionalValidationCase(
			"validation_Create_Databases_requireAllowedDatabases",
			func(opts *CreateReplicationGroupOptions) {
				opts.ObjectTypes = []ReplicationGroupObjectType{ReplicationGroupObjectTypeDatabases}
			},
			NewError("ALLOWED_DATABASES must be set when OBJECT_TYPES includes DATABASES"),
		).
		withAdditionalValidationCase(
			"validation_Create_AllowedDatabases_requireDatabases",
			func(opts *CreateReplicationGroupOptions) {
				opts.AllowedDatabases = []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")}
			},
			NewError("OBJECT_TYPES must include DATABASES when ALLOWED_DATABASES is set"),
		).
		withAdditionalValidationCase(
			"validation_Create_ExternalVolumes_requireAllowedExternalVolumes",
			func(opts *CreateReplicationGroupOptions) {
				opts.ObjectTypes = []ReplicationGroupObjectType{ReplicationGroupObjectTypeExternalVolumes}
			},
			NewError("ALLOWED_EXTERNAL_VOLUMES must be set when OBJECT_TYPES includes EXTERNAL VOLUMES"),
		).
		withAdditionalValidationCase(
			"validation_Create_AllowedExternalVolumes_requireExternalVolumes",
			func(opts *CreateReplicationGroupOptions) {
				opts.AllowedExternalVolumes = []AccountObjectIdentifier{NewAccountObjectIdentifier("vol1")}
			},
			NewError("OBJECT_TYPES must include EXTERNAL VOLUMES when ALLOWED_EXTERNAL_VOLUMES is set"),
		).
		withAdditionalValidationCase(
			"validation_Create_Shares_requireAllowedShares",
			func(opts *CreateReplicationGroupOptions) {
				opts.ObjectTypes = []ReplicationGroupObjectType{ReplicationGroupObjectTypeShares}
			},
			NewError("ALLOWED_SHARES must be set when OBJECT_TYPES includes SHARES"),
		).
		withAdditionalValidationCase(
			"validation_Create_AllowedShares_requireShares",
			func(opts *CreateReplicationGroupOptions) {
				opts.AllowedShares = []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")}
			},
			NewError("OBJECT_TYPES must include SHARES when ALLOWED_SHARES is set"),
		).
		withAdditionalValidationCase(
			"validation_Create_Integrations_requireAllowedIntegrationTypes",
			func(opts *CreateReplicationGroupOptions) {
				opts.ObjectTypes = []ReplicationGroupObjectType{ReplicationGroupObjectTypeIntegrations}
			},
			NewError("ALLOWED_INTEGRATION_TYPES must be set when OBJECT_TYPES includes INTEGRATIONS"),
		).
		withAdditionalValidationCase(
			"validation_Create_AllowedIntegrationTypes_requireIntegrations",
			func(opts *CreateReplicationGroupOptions) {
				opts.AllowedIntegrationTypes = []ReplicationGroupIntegrationType{ReplicationGroupIntegrationTypeApiIntegrations}
			},
			NewError("OBJECT_TYPES must include INTEGRATIONS when ALLOWED_INTEGRATION_TYPES is set"),
		).
		withAdditionalValidationCase(
			"validation_Create_SecurityIntegrations_requireRoles",
			func(opts *CreateReplicationGroupOptions) {
				opts.ObjectTypes = []ReplicationGroupObjectType{ReplicationGroupObjectTypeIntegrations}
				opts.AllowedIntegrationTypes = []ReplicationGroupIntegrationType{ReplicationGroupIntegrationTypeSecurityIntegrations}
			},
			NewError("OBJECT_TYPES must include ROLES when ALLOWED_INTEGRATION_TYPES includes SECURITY INTEGRATIONS"),
		)

	replicationGroupsTests.Show.
		withExpectedSqlf(case_ReplicationGroups_sql_Show_basic, "SHOW REPLICATION GROUPS").
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_Show_all,
			func(*ShowReplicationGroupOptions) {},
			"SHOW REPLICATION GROUPS",
		)

	replicationGroupsTests.Drop.
		withExpectedSqlf(
			case_ReplicationGroups_sql_Drop_basic,
			"DROP REPLICATION GROUP %s",
			replicationGroupsTestIdAccountObjectIdentifier.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_Drop_all,
			func(opts *DropReplicationGroupOptions) { opts.IfExists = new(true) },
			"DROP REPLICATION GROUP IF EXISTS %s",
			replicationGroupsTestIdAccountObjectIdentifier.FullyQualifiedName(),
		)

	moveTarget := NewAccountObjectIdentifier("myrg2")
	replicationGroupsTests.Alter.
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_Alter_RenameTo,
			func(opts *AlterReplicationGroupOptions) {
				opts.RenameTo = new(NewAccountObjectIdentifier("myrg2"))
			},
			`ALTER REPLICATION GROUP %s RENAME TO "myrg2"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_Alter_Set,
			func(opts *AlterReplicationGroupOptions) {
				opts.Set = &ReplicationGroupSet{
					ObjectTypes: []ReplicationGroupObjectType{
						ReplicationGroupObjectTypeAccountParameters,
						ReplicationGroupObjectTypeDatabases,
						ReplicationGroupObjectTypeExternalVolumes,
						ReplicationGroupObjectTypeIntegrations,
						ReplicationGroupObjectTypeNetworkPolicies,
						ReplicationGroupObjectTypeResourceMonitors,
						ReplicationGroupObjectTypeRoles,
						ReplicationGroupObjectTypeShares,
						ReplicationGroupObjectTypeUsers,
						ReplicationGroupObjectTypeWarehouses,
					},
					AllowedDatabases:       []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")},
					AllowedExternalVolumes: []AccountObjectIdentifier{NewAccountObjectIdentifier("vol1")},
					AllowedShares:          []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")},
					AllowedIntegrationTypes: []ReplicationGroupIntegrationType{
						ReplicationGroupIntegrationTypeSecurityIntegrations,
						ReplicationGroupIntegrationTypeApiIntegrations,
						ReplicationGroupIntegrationTypeStorageIntegrations,
						ReplicationGroupIntegrationTypeExternalAccessIntegrations,
						ReplicationGroupIntegrationTypeNotificationIntegrations,
					},
					Comment:             new("comment"),
					ReplicationSchedule: new("10 MINUTE"),
					ErrorIntegration:    new(NewAccountObjectIdentifier("error_int")),
				}
			},
			`ALTER REPLICATION GROUP %s SET OBJECT_TYPES = ACCOUNT PARAMETERS, DATABASES, EXTERNAL VOLUMES, INTEGRATIONS, NETWORK POLICIES, RESOURCE MONITORS, ROLES, SHARES, USERS, WAREHOUSES ALLOWED_DATABASES = "db1" ALLOWED_EXTERNAL_VOLUMES = "vol1" ALLOWED_SHARES = "share1" ALLOWED_INTEGRATION_TYPES = SECURITY INTEGRATIONS, API INTEGRATIONS, STORAGE INTEGRATIONS, EXTERNAL ACCESS INTEGRATIONS, NOTIFICATION INTEGRATIONS COMMENT = 'comment' REPLICATION_SCHEDULE = '10 MINUTE' ERROR_INTEGRATION = "error_int"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_Alter_Unset,
			func(opts *AlterReplicationGroupOptions) {
				opts.Unset = &ReplicationGroupUnset{
					Comment:             new(true),
					ReplicationSchedule: new(true),
					ErrorIntegration:    new(true),
				}
			},
			"ALTER REPLICATION GROUP %s UNSET COMMENT, REPLICATION_SCHEDULE, ERROR_INTEGRATION",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_Alter_SetTags,
			func(opts *AlterReplicationGroupOptions) {
				opts.SetTags = []TagAssociation{
					{Name: NewAccountObjectIdentifier("tag1"), Value: "value1"},
					{Name: NewAccountObjectIdentifier("tag2"), Value: "value2"},
				}
			},
			`ALTER REPLICATION GROUP %s SET TAG "tag1" = 'value1', "tag2" = 'value2'`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_Alter_UnsetTags,
			func(opts *AlterReplicationGroupOptions) {
				opts.UnsetTags = []ObjectIdentifier{
					NewAccountObjectIdentifier("tag1"),
					NewAccountObjectIdentifier("tag2"),
				}
			},
			`ALTER REPLICATION GROUP %s UNSET TAG "tag1", "tag2"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_Alter_Move,
			func(opts *AlterReplicationGroupOptions) {
				opts.Move = &ReplicationGroupMove{
					Databases:          []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")},
					ToReplicationGroup: moveTarget,
				}
			},
			`ALTER REPLICATION GROUP %s MOVE DATABASES "db1" TO REPLICATION GROUP "myrg2"`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Move_externalVolumes",
			func(opts *AlterReplicationGroupOptions) {
				opts.Move = &ReplicationGroupMove{
					ExternalVolumes:    []AccountObjectIdentifier{NewAccountObjectIdentifier("vol1")},
					ToReplicationGroup: moveTarget,
				}
			},
			`ALTER REPLICATION GROUP %s MOVE EXTERNAL VOLUMES "vol1" TO REPLICATION GROUP "myrg2"`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Move_shares",
			func(opts *AlterReplicationGroupOptions) {
				opts.Move = &ReplicationGroupMove{
					Shares:             []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")},
					ToReplicationGroup: moveTarget,
				}
			},
			`ALTER REPLICATION GROUP %s MOVE SHARES "share1" TO REPLICATION GROUP "myrg2"`,
			id.FullyQualifiedName(),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_AllowedDatabases_requireDatabases",
			func(opts *AlterReplicationGroupOptions) {
				opts.Set = &ReplicationGroupSet{
					AllowedDatabases: []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")},
				}
			},
			NewError("OBJECT_TYPES must include DATABASES when ALLOWED_DATABASES is set"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_AllowedExternalVolumes_requireExternalVolumes",
			func(opts *AlterReplicationGroupOptions) {
				opts.Set = &ReplicationGroupSet{
					AllowedExternalVolumes: []AccountObjectIdentifier{NewAccountObjectIdentifier("vol1")},
				}
			},
			NewError("OBJECT_TYPES must include EXTERNAL VOLUMES when ALLOWED_EXTERNAL_VOLUMES is set"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_AllowedShares_requireShares",
			func(opts *AlterReplicationGroupOptions) {
				opts.Set = &ReplicationGroupSet{
					AllowedShares: []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")},
				}
			},
			NewError("OBJECT_TYPES must include SHARES when ALLOWED_SHARES is set"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_AllowedIntegrationTypes_requireIntegrations",
			func(opts *AlterReplicationGroupOptions) {
				opts.Set = &ReplicationGroupSet{
					AllowedIntegrationTypes: []ReplicationGroupIntegrationType{ReplicationGroupIntegrationTypeApiIntegrations},
				}
			},
			NewError("OBJECT_TYPES must include INTEGRATIONS when ALLOWED_INTEGRATION_TYPES is set"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_SecurityIntegrations_requireRoles",
			func(opts *AlterReplicationGroupOptions) {
				opts.Set = &ReplicationGroupSet{
					ObjectTypes:             []ReplicationGroupObjectType{ReplicationGroupObjectTypeIntegrations},
					AllowedIntegrationTypes: []ReplicationGroupIntegrationType{ReplicationGroupIntegrationTypeSecurityIntegrations},
				}
			},
			NewError("OBJECT_TYPES must include ROLES when ALLOWED_INTEGRATION_TYPES includes SECURITY INTEGRATIONS"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Move_exactlyOneList",
			func(opts *AlterReplicationGroupOptions) {
				opts.Move = &ReplicationGroupMove{ToReplicationGroup: moveTarget}
			},
			NewError("MOVE must set exactly one of DATABASES, EXTERNAL VOLUMES, and SHARES"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Move_To_ValidIdentifier",
			func(opts *AlterReplicationGroupOptions) {
				opts.Move = &ReplicationGroupMove{
					Databases: []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")},
				}
			},
			errInvalidIdentifier("AlterReplicationGroupOptions.Move", "ToReplicationGroup"),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_Alter_Add,
			func(opts *AlterReplicationGroupOptions) {
				opts.Add = &ReplicationGroupAdd{
					ToAllowedDatabases: []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")},
				}
			},
			`ALTER REPLICATION GROUP %s ADD "db1" TO ALLOWED_DATABASES`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Add_externalVolumes",
			func(opts *AlterReplicationGroupOptions) {
				opts.Add = &ReplicationGroupAdd{
					ToAllowedExternalVolumes: []AccountObjectIdentifier{NewAccountObjectIdentifier("vol1")},
				}
			},
			`ALTER REPLICATION GROUP %s ADD "vol1" TO ALLOWED_EXTERNAL_VOLUMES`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Add_shares",
			func(opts *AlterReplicationGroupOptions) {
				opts.Add = &ReplicationGroupAdd{
					ToAllowedShares: []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")},
				}
			},
			`ALTER REPLICATION GROUP %s ADD "share1" TO ALLOWED_SHARES`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Add_accounts",
			func(opts *AlterReplicationGroupOptions) {
				opts.Add = &ReplicationGroupAdd{
					ToAllowedAccounts:  []AccountIdentifier{NewAccountIdentifier("MY_ORG", "MY_ACCOUNT")},
					IgnoreEditionCheck: new(true),
				}
			},
			`ALTER REPLICATION GROUP %s ADD "MY_ORG"."MY_ACCOUNT" TO ALLOWED_ACCOUNTS IGNORE EDITION CHECK`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_Alter_Remove,
			func(opts *AlterReplicationGroupOptions) {
				opts.Remove = &ReplicationGroupRemove{
					FromAllowedDatabases: []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")},
				}
			},
			`ALTER REPLICATION GROUP %s REMOVE "db1" FROM ALLOWED_DATABASES`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Remove_externalVolumes",
			func(opts *AlterReplicationGroupOptions) {
				opts.Remove = &ReplicationGroupRemove{
					FromAllowedExternalVolumes: []AccountObjectIdentifier{NewAccountObjectIdentifier("vol1")},
				}
			},
			`ALTER REPLICATION GROUP %s REMOVE "vol1" FROM ALLOWED_EXTERNAL_VOLUMES`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Remove_shares",
			func(opts *AlterReplicationGroupOptions) {
				opts.Remove = &ReplicationGroupRemove{
					FromAllowedShares: []AccountObjectIdentifier{NewAccountObjectIdentifier("share1")},
				}
			},
			`ALTER REPLICATION GROUP %s REMOVE "share1" FROM ALLOWED_SHARES`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Remove_accounts",
			func(opts *AlterReplicationGroupOptions) {
				opts.Remove = &ReplicationGroupRemove{
					FromAllowedAccounts: []AccountIdentifier{NewAccountIdentifier("MY_ORG", "MY_ACCOUNT")},
				}
			},
			`ALTER REPLICATION GROUP %s REMOVE "MY_ORG"."MY_ACCOUNT" FROM ALLOWED_ACCOUNTS`,
			id.FullyQualifiedName(),
		).
		withAdditionalValidationCase(
			"validation_Alter_Add_exactlyOneList",
			func(opts *AlterReplicationGroupOptions) {
				opts.Add = &ReplicationGroupAdd{}
			},
			NewError("ADD must set exactly one allowed list"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Add_IgnoreEditionCheck_requiresAccounts",
			func(opts *AlterReplicationGroupOptions) {
				opts.Add = &ReplicationGroupAdd{
					ToAllowedDatabases: []AccountObjectIdentifier{NewAccountObjectIdentifier("db1")},
					IgnoreEditionCheck: new(true),
				}
			},
			NewError("IGNORE EDITION CHECK requires ADD TO ALLOWED_ACCOUNTS"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Remove_exactlyOneList",
			func(opts *AlterReplicationGroupOptions) {
				opts.Remove = &ReplicationGroupRemove{}
			},
			NewError("REMOVE must set exactly one allowed list"),
		)

	replicationGroupsTests.AlterSecondary.
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_AlterSecondary_Refresh,
			func(opts *AlterSecondaryReplicationGroupOptions) { opts.Refresh = new(true) },
			"ALTER REPLICATION GROUP %s REFRESH",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_AlterSecondary_Suspend,
			func(opts *AlterSecondaryReplicationGroupOptions) {
				opts.Suspend = &ReplicationGroupSecondarySuspend{}
			},
			"ALTER REPLICATION GROUP %s SUSPEND",
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_AlterSecondary_Suspend_immediate",
			func(opts *AlterSecondaryReplicationGroupOptions) {
				opts.Suspend = &ReplicationGroupSecondarySuspend{Immediate: new(true)}
			},
			"ALTER REPLICATION GROUP %s SUSPEND IMMEDIATE",
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_AlterSecondary_Resume,
			func(opts *AlterSecondaryReplicationGroupOptions) { opts.Resume = new(true) },
			"ALTER REPLICATION GROUP %s RESUME",
			id.FullyQualifiedName(),
		)

	primary := NewExternalObjectIdentifier(NewAccountIdentifier("myorg", "myaccount"), NewAccountObjectIdentifier("myrg"))
	replicationGroupsTests.CreateSecondary.
		withDefaultOpts(func() *CreateSecondaryReplicationGroupOptions {
			return &CreateSecondaryReplicationGroupOptions{
				name:        id,
				AsReplicaOf: primary,
			}
		}).
		withExpectedSqlf(
			case_ReplicationGroups_sql_CreateSecondary_basic,
			`CREATE REPLICATION GROUP %s AS REPLICA OF "myorg"."myaccount"."myrg"`,
			id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_ReplicationGroups_sql_CreateSecondary_all,
			func(opts *CreateSecondaryReplicationGroupOptions) { opts.IfNotExists = new(true) },
			`CREATE REPLICATION GROUP IF NOT EXISTS %s AS REPLICA OF "myorg"."myaccount"."myrg"`,
			id.FullyQualifiedName(),
		).
		withAdditionalValidationCase(
			"validation_CreateSecondary_AsReplicaOf_ValidIdentifier",
			func(opts *CreateSecondaryReplicationGroupOptions) {
				opts.AsReplicaOf = emptyExternalObjectIdentifier
			},
			errInvalidIdentifier("CreateSecondaryReplicationGroupOptions", "AsReplicaOf"),
		)
}
