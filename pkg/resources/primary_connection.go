package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var primaryConnectionSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("String that specifies the identifier (i.e. name) for the connection. Must start with an alphabetic character and may only contain letters, decimal digits (0-9), and underscores (_). For a primary connection, the name must be unique across connection names and account names in the organization. "),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"is_primary": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Indicates if the connection is primary. When Terraform detects that the connection is not primary, the resource is recreated.",
	},
	"enable_failover_to_accounts": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: relatedResourceDescription("Enables failover for given connection to provided accounts. Specifies a list of accounts in your organization where a secondary connection for this primary connection can be promoted to serve as the primary connection. Include your organization name for each account in the list.", resources.Account),
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the connection.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW CONNECTIONS` for the given connection.",
		Elem: &schema.Resource{
			Schema: schemas.ShowConnectionSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func PrimaryConnection() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseAccountObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.AccountObjectIdentifier] {
			return client.Connections.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.PrimaryConnection, CreateContextPrimaryConnection),
		ReadContext:   TrackingReadWrapper(resources.PrimaryConnection, ReadContextPrimaryConnection),
		UpdateContext: TrackingUpdateWrapper(resources.PrimaryConnection, UpdateContextPrimaryConnection),
		DeleteContext: TrackingDeleteWrapper(resources.PrimaryConnection, deleteFunc),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.PrimaryConnection, customdiff.All(
			ComputedIfAnyAttributeChanged(primaryConnectionSchema, ShowOutputAttributeName, "comment", "is_primary", "enable_failover_to_accounts"),
			RecreateWhenResourceBoolFieldChangedExternally("is_primary", true),
		)),

		Description: "Resource used to manage primary connections. For managing replicated connection check resource [snowflake_secondary_connection](./secondary_connection). For more information, check [connection documentation](https://docs.snowflake.com/en/sql-reference/sql/create-connection.html).",
		Schema:      primaryConnectionSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.PrimaryConnection, ImportName[sdk.AccountObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateContextPrimaryConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateConnectionRequest(id)

	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	err = client.Connections.Create(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	if v, ok := d.GetOk("enable_failover_to_accounts"); ok {
		enableFailoverConfig := v.([]any)

		enableFailoverToAccountsList := make([]sdk.AccountIdentifier, 0)
		for _, enableToAccount := range enableFailoverConfig {
			accountInConfig := enableToAccount.(string)
			accountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(accountInConfig)

			enableFailoverToAccountsList = append(enableFailoverToAccountsList, accountIdentifier)
		}

		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
			WithEnableConnectionFailover(*sdk.NewEnableConnectionFailoverRequest(enableFailoverToAccountsList)))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextPrimaryConnection(ctx, d, meta)
}

func ReadContextPrimaryConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	connections, err := client.Connections.Show(ctx, sdk.NewShowConnectionRequest().WithLike(sdk.Like{Pattern: sdk.String(id.Name())}))
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to retrieve connection.",
				Detail:   fmt.Sprintf("Connection name: %s, Err: %s", id.FullyQualifiedName(), err),
			},
		}
	}

	connection, err := collections.FindFirst(connections, func(c sdk.Connection) bool {
		return c.Name == id.Name() && c.AccountLocator == client.GetAccountLocator()
	})
	if errors.Is(err, sdk.ErrObjectNotFound) {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to retrieve connection. Target object not found. Marking the resource as removed.",
				Detail:   fmt.Sprintf("Connection name: %s, Err: %s", id.FullyQualifiedName(), err),
			},
		}
	}

	errs := errors.Join(
		d.Set("is_primary", connection.IsPrimary),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.ConnectionToSchema(connection)}),
		d.Set("comment", connection.Comment),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	sessionDetails, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	currentAccountIdentifier := sdk.NewAccountIdentifier(sessionDetails.OrganizationName, sessionDetails.AccountName)

	enableFailoverToAccounts := make([]string, 0)
	for _, allowedAccount := range connection.FailoverAllowedToAccounts {
		if currentAccountIdentifier.FullyQualifiedName() == allowedAccount.FullyQualifiedName() {
			continue
		}
		enableFailoverToAccounts = append(enableFailoverToAccounts, allowedAccount.Name())
	}

	if len(enableFailoverToAccounts) == 0 {
		err := d.Set("enable_failover_to_accounts", []any{})
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err := d.Set("enable_failover_to_accounts", enableFailoverToAccounts)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func UpdateContextPrimaryConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	connectionSetRequest := new(sdk.ConnectionSetRequest)
	connectionUnsetRequest := new(sdk.ConnectionUnsetRequest)

	if d.HasChange("enable_failover_to_accounts") {
		before, after := d.GetChange("enable_failover_to_accounts")

		getFailoverToAccounts := func(failoverConfig []any) []sdk.AccountIdentifier {
			failoverEnabledToAccounts := make([]sdk.AccountIdentifier, 0)

			for _, allowedAccount := range failoverConfig {
				accountIdentifier := sdk.NewAccountIdentifierFromFullyQualifiedName(allowedAccount.(string))
				failoverEnabledToAccounts = append(failoverEnabledToAccounts, accountIdentifier)
			}
			return failoverEnabledToAccounts
		}

		beforeFailover := getFailoverToAccounts(before.([]any))
		afterFailover := getFailoverToAccounts(after.([]any))

		addedFailovers, removedFailovers := ListDiff(beforeFailover, afterFailover)

		if len(addedFailovers) > 0 {
			err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
				WithEnableConnectionFailover(*sdk.NewEnableConnectionFailoverRequest(addedFailovers)),
			)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(removedFailovers) > 0 {
			err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).
				WithDisableConnectionFailover(*sdk.NewDisableConnectionFailoverRequest().
					WithToAccounts(*sdk.NewToAccountsRequest(removedFailovers)),
				),
			)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		if len(comment) > 0 {
			connectionSetRequest.WithComment(comment)
		} else {
			connectionUnsetRequest.WithComment(true)
		}
	}

	if (*connectionSetRequest != sdk.ConnectionSetRequest{}) {
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).WithSet(*connectionSetRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if (*connectionUnsetRequest != sdk.ConnectionUnsetRequest{}) {
		err := client.Connections.Alter(ctx, sdk.NewAlterConnectionRequest(id).WithUnset(*connectionUnsetRequest))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextPrimaryConnection(ctx, d, meta)
}
