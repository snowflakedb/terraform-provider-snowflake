package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"
	providerresources "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"

	"github.com/hashicorp/go-cty/cty"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var databaseSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the database; must be unique for your account. As a best practice for [Database Replication and Failover](https://docs.snowflake.com/en/user-guide/db-replication-intro), it is recommended to give each secondary database the same name as its primary database. This practice supports referencing fully-qualified objects (i.e. '<db>.<schema>.<object>') by other objects in the same database, such as querying a fully-qualified table name in a view. If a secondary database has a different name from the primary database, then these object references would break in the secondary database."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"drop_public_schema_on_creation": {
		Type:             schema.TypeBool,
		Optional:         true,
		Description:      "Specifies whether to drop public schema on creation or not. Modifying the parameter after database is already created won't have any effect.",
		DiffSuppressFunc: IgnoreAfterCreation,
	},
	"is_transient": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies the database as transient. Transient databases do not have a Fail-safe period so they do not incur additional storage costs once they leave Time Travel; however, this means they are also not protected by Fail-safe in the event of a data loss.",
	},
	"replication": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Configures replication for a given database. When specified, this database will be promoted to serve as a primary database for replication. A primary database can be replicated in one or more accounts, allowing users in those accounts to query objects in each secondary (i.e. replica) database.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enable_to_account": {
					Type:        schema.TypeList,
					Required:    true,
					Description: "Entry to enable replication and optionally failover for a given account identifier.",
					MinItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"account_identifier": {
								Type:     schema.TypeString,
								Required: true,
								// TODO(SNOW-1438810): Add account identifier validator
								Description: relatedResourceDescription("Specifies account identifier for which replication should be enabled. The account identifiers should be in the form of `\"<organization_name>\".\"<account_name>\"`.", providerresources.Account),
							},
							"with_failover": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: "Specifies if failover should be enabled for the specified account identifier",
							},
						},
					},
				},
				"ignore_edition_check": {
					Type:     schema.TypeBool,
					Optional: true,
					Description: "Allows replicating data to accounts on lower editions in either of the following scenarios: " +
						"1. The primary database is in a Business Critical (or higher) account but one or more of the accounts approved for replication are on lower editions. Business Critical Edition is intended for Snowflake accounts with extremely sensitive data. " +
						"2. The primary database is in a Business Critical (or higher) account and a signed business associate agreement is in place to store PHI data in the account per HIPAA and HITRUST regulations, but no such agreement is in place for one or more of the accounts approved for replication, regardless if they are Business Critical (or higher) accounts. " +
						"Both scenarios are prohibited by default in an effort to help prevent account administrators for Business Critical (or higher) accounts from inadvertently replicating sensitive data to accounts on lower editions.",
				},
			},
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the database.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func Database() *schema.Resource {
	// TODO(SNOW-1818849): unassign network policies inside the database before dropping
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.Databases.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.Database, CreateDatabase),
		UpdateContext: TrackingUpdateWrapper(resources.Database, UpdateDatabase),
		ReadContext:   TrackingReadWrapper(resources.Database, ReadDatabase),
		DeleteContext: TrackingDeleteWrapper(resources.Database, deleteFunc),
		Description:   "Represents a standard database. If replication configuration is specified, the database is promoted to serve as a primary database for replication.",

		Schema: collections.MergeMaps(databaseSchema, databaseParametersSchema),
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Database, ImportName[sdk.AccountObjectIdentifier]),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Database, customdiff.All(
			ComputedIfAnyAttributeChanged(databaseSchema, FullyQualifiedNameAttributeName, "name"),
			databaseParametersCustomDiff,
		)),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: v092DatabaseStateUpgrader,
			},
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	opts := &sdk.CreateDatabaseOptions{
		Transient: GetConfigPropertyAsPointerAllowingZeroValue[bool](d, "is_transient"),
		Comment:   GetConfigPropertyAsPointerAllowingZeroValue[string](d, "comment"),
	}
	if parametersCreateDiags := handleDatabaseParametersCreate(d, opts); len(parametersCreateDiags) > 0 {
		return parametersCreateDiags
	}

	err = client.Databases.Create(ctx, id, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	var diags diag.Diagnostics

	if d.Get("drop_public_schema_on_creation").(bool) {
		var dropSchemaErrs []error
		err := util.Retry(3, time.Second, func() (error, bool) {
			if err := client.Schemas.Drop(ctx, sdk.NewDatabaseObjectIdentifier(id.Name(), "PUBLIC"), &sdk.DropSchemaOptions{IfExists: sdk.Bool(true)}); err != nil {
				dropSchemaErrs = append(dropSchemaErrs, err)
				return nil, false
			}
			return nil, true
		})
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to drop public schema on creation (failed after 3 attempts)",
				Detail:   fmt.Sprintf("The '%s' database was created successfully, but the provider was not able to remove public schema on creation. Please drop the public schema manually. Original errors: %s", id.Name(), errors.Join(dropSchemaErrs...)),
			})
		}
	}

	if v, ok := d.GetOk("replication"); ok {
		replicationConfiguration := v.([]any)[0].(map[string]any)

		var ignoreEditionCheck *bool
		if v, ok := replicationConfiguration["ignore_edition_check"]; ok {
			ignoreEditionCheck = sdk.Pointer(v.(bool))
		}

		if enableToAccounts, ok := replicationConfiguration["enable_to_account"]; ok {
			enableToAccountList := enableToAccounts.([]any)

			if len(enableToAccountList) > 0 {
				replicationToAccounts := make([]sdk.AccountIdentifier, 0)
				failoverToAccounts := make([]sdk.AccountIdentifier, 0)

				for _, enableToAccount := range enableToAccountList {
					accountConfig := enableToAccount.(map[string]any)
					accountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(accountConfig["account_identifier"].(string))

					replicationToAccounts = append(replicationToAccounts, accountIdentifier)
					if v, ok := accountConfig["with_failover"]; ok && v.(bool) {
						failoverToAccounts = append(failoverToAccounts, accountIdentifier)
					}
				}

				if len(replicationToAccounts) > 0 {
					err := client.Databases.AlterReplication(ctx, id, &sdk.AlterDatabaseReplicationOptions{
						EnableReplication: &sdk.EnableReplication{
							ToAccounts:         replicationToAccounts,
							IgnoreEditionCheck: ignoreEditionCheck,
						},
					})
					if err != nil {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Summary:  err.Error(),
						})
					}
				}

				if len(failoverToAccounts) > 0 {
					err = client.Databases.AlterFailover(ctx, id, &sdk.AlterDatabaseFailoverOptions{
						EnableFailover: &sdk.EnableFailover{
							ToAccounts: failoverToAccounts,
						},
					})
					if err != nil {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Warning,
							Summary:  err.Error(),
						})
					}
				}
			}
		}
	}

	return append(diags, ReadDatabase(ctx, d, meta)...)
}

func UpdateDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		err = client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
			NewName: &newId,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	databaseSetRequest := new(sdk.DatabaseSet)
	databaseUnsetRequest := new(sdk.DatabaseUnset)

	if updateParamDiags := handleDatabaseParametersChanges(d, databaseSetRequest, databaseUnsetRequest); len(updateParamDiags) > 0 {
		return updateParamDiags
	}

	if d.HasChange("replication") {
		before, after := d.GetChange("replication")

		getReplicationConfiguration := func(replicationConfigs []any) (replicationEnabledToAccounts []sdk.AccountIdentifier, failoverEnabledToAccounts []sdk.AccountIdentifier) {
			replicationEnabledToAccounts = make([]sdk.AccountIdentifier, 0)
			failoverEnabledToAccounts = make([]sdk.AccountIdentifier, 0)

			for _, replicationConfigurationMap := range replicationConfigs {
				replicationConfiguration := replicationConfigurationMap.(map[string]any)
				for _, enableToAccountMap := range replicationConfiguration["enable_to_account"].([]any) {
					enableToAccount := enableToAccountMap.(map[string]any)
					accountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(enableToAccount["account_identifier"].(string))

					replicationEnabledToAccounts = append(replicationEnabledToAccounts, accountIdentifier)
					if enableToAccount["with_failover"].(bool) {
						failoverEnabledToAccounts = append(failoverEnabledToAccounts, accountIdentifier)
					}
				}
			}

			return replicationEnabledToAccounts, failoverEnabledToAccounts
		}
		beforeReplicationEnabledToAccounts, beforeFailoverEnabledToAccounts := getReplicationConfiguration(before.([]any))
		afterReplicationEnabledToAccounts, afterFailoverEnabledToAccounts := getReplicationConfiguration(after.([]any))

		addedFailovers, removedFailovers := ListDiff(beforeFailoverEnabledToAccounts, afterFailoverEnabledToAccounts)
		addedReplications, removedReplications := ListDiff(beforeReplicationEnabledToAccounts, afterReplicationEnabledToAccounts)
		// Failovers will be disabled implicitly by disabled replications
		removedFailovers = slices.DeleteFunc(removedFailovers, func(identifier sdk.AccountIdentifier) bool { return slices.Contains(removedReplications, identifier) })

		if len(addedReplications) > 0 {
			err := client.Databases.AlterReplication(ctx, id, &sdk.AlterDatabaseReplicationOptions{
				EnableReplication: &sdk.EnableReplication{
					ToAccounts:         addedReplications,
					IgnoreEditionCheck: sdk.Bool(d.Get("replication.0.ignore_edition_check").(bool)),
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(addedFailovers) > 0 {
			err := client.Databases.AlterFailover(ctx, id, &sdk.AlterDatabaseFailoverOptions{
				EnableFailover: &sdk.EnableFailover{
					ToAccounts: addedFailovers,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(removedReplications) > 0 {
			err := client.Databases.AlterReplication(ctx, id, &sdk.AlterDatabaseReplicationOptions{
				DisableReplication: &sdk.DisableReplication{
					ToAccounts: removedReplications,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(removedFailovers) > 0 {
			err := client.Databases.AlterFailover(ctx, id, &sdk.AlterDatabaseFailoverOptions{
				DisableFailover: &sdk.DisableFailover{
					ToAccounts: removedFailovers,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if len(comment) > 0 {
			databaseSetRequest.Comment = &comment
		} else {
			databaseUnsetRequest.Comment = sdk.Bool(true)
		}
	}

	if (*databaseSetRequest != sdk.DatabaseSet{}) {
		err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
			Set: databaseSetRequest,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if (*databaseUnsetRequest != sdk.DatabaseUnset{}) {
		err := client.Databases.Alter(ctx, id, &sdk.AlterDatabaseOptions{
			Unset: databaseUnsetRequest,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadDatabase(ctx, d, meta)
}

func ReadDatabase(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	database, err := client.Databases.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query database. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Database id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("is_transient", database.Transient); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", database.Comment); err != nil {
		return diag.FromErr(err)
	}

	sessionDetails, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	currentAccountIdentifier := sdk.NewAccountIdentifier(sessionDetails.OrganizationName, sessionDetails.AccountName)
	replicationDatabases, err := client.ReplicationFunctions.ShowReplicationDatabases(ctx, &sdk.ShowReplicationDatabasesOptions{
		WithPrimary: sdk.Pointer(sdk.NewExternalObjectIdentifier(currentAccountIdentifier, id)),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if len(replicationDatabases) == 1 {
		replicationAllowedToAccounts := make([]sdk.AccountIdentifier, 0)
		failoverAllowedToAccounts := make([]sdk.AccountIdentifier, 0)

		for _, allowedAccount := range strings.Split(replicationDatabases[0].ReplicationAllowedToAccounts, ",") {
			allowedAccountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(strings.TrimSpace(allowedAccount))
			if currentAccountIdentifier.FullyQualifiedName() == allowedAccountIdentifier.FullyQualifiedName() {
				continue
			}
			replicationAllowedToAccounts = append(replicationAllowedToAccounts, allowedAccountIdentifier)
		}

		for _, allowedAccount := range strings.Split(replicationDatabases[0].FailoverAllowedToAccounts, ",") {
			allowedAccountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(strings.TrimSpace(allowedAccount))
			if currentAccountIdentifier.FullyQualifiedName() == allowedAccountIdentifier.FullyQualifiedName() {
				continue
			}
			failoverAllowedToAccounts = append(failoverAllowedToAccounts, allowedAccountIdentifier)
		}

		enableToAccount := make([]map[string]any, 0)
		for _, allowedAccount := range replicationAllowedToAccounts {
			enableToAccount = append(enableToAccount, map[string]any{
				"account_identifier": allowedAccount.FullyQualifiedName(),
				"with_failover":      slices.Contains(failoverAllowedToAccounts, allowedAccount),
			})
		}

		var ignoreEditionCheck bool
		if v, ok := d.GetOk("replication.0.ignore_edition_check"); ok {
			ignoreEditionCheck = v.(bool)
		}

		if len(enableToAccount) == 0 {
			err := d.Set("replication", []any{})
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := d.Set("replication", []any{
				map[string]any{
					"enable_to_account":    enableToAccount,
					"ignore_edition_check": ignoreEditionCheck,
				},
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	databaseParameters, err := client.Databases.ShowParameters(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := handleDatabaseParameterRead(d, databaseParameters); diags != nil {
		return diags
	}

	return nil
}
