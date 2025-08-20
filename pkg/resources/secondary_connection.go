package resources

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var secondaryConnectionSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("String that specifies the identifier (i.e. name) for the connection. Must start with an alphabetic character and may only contain letters, decimal digits (0-9), and underscores (_). For a secondary connection, the name must match the name of its primary connection."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"is_primary": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Indicates if the connection primary status has been changed. If change is detected, resource will be recreated.",
	},
	"as_replica_of": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      relatedResourceDescription("Specifies the identifier for a primary connection from which to create a replica (i.e. a secondary connection).", resources.PrimaryConnection),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the secondary connection.",
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

func SecondaryConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.SecondaryConnection, CreateContextSecondaryConnection),
		ReadContext:   TrackingReadWrapper(resources.SecondaryConnection, ReadContextSecondaryConnection),
		UpdateContext: TrackingUpdateWrapper(resources.SecondaryConnection, UpdateContextSecondaryConnection),
		DeleteContext: TrackingDeleteWrapper(resources.SecondaryConnection, DeleteContextSecondaryConnection),
		Description:   "Resource used to manage secondary (replicated) connections. To manage primary connection check resource [snowflake_primary_connection](./primary_connection). For more information, check [connection documentation](https://docs.snowflake.com/en/sql-reference/sql/create-connection.html).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.SecondaryConnection, customdiff.All(
			ComputedIfAnyAttributeChanged(secondaryConnectionSchema, ShowOutputAttributeName, "comment", "is_primary"),
			RecreateWhenResourceBoolFieldChangedExternally("is_primary", false),
		)),

		Schema: secondaryConnectionSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.SecondaryConnection, ImportName[sdk.AccountObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateContextSecondaryConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateConnectionRequest(id)

	if v, ok := d.GetOk("as_replica_of"); ok {
		externalObjectId, err := sdk.ParseExternalObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		request.WithAsReplicaOf(externalObjectId)
	}

	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	err = client.Connections.Create(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextSecondaryConnection(ctx, d, meta)
}

func ReadContextSecondaryConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	connection, err := client.Connections.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query secondary connection. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Secondary connection id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Failed to query secondary connection.",
				Detail:   fmt.Sprintf("Secondary connection id: %s, Err: %s", id.FullyQualifiedName(), err),
			},
		}
	}

	return diag.FromErr(errors.Join(
		d.Set("is_primary", connection.IsPrimary),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.ConnectionToSchema(connection)}),
		d.Set("comment", connection.Comment),
		d.Set("as_replica_of", connection.Primary.FullyQualifiedName()),
	))
}

func UpdateContextSecondaryConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	connectionSetRequest := new(sdk.ConnectionSetRequest)
	connectionUnsetRequest := new(sdk.ConnectionUnsetRequest)

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

	return ReadContextSecondaryConnection(ctx, d, meta)
}

func DeleteContextSecondaryConnection(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Retry is necessary for cases where the changes in Snowflake didn't have enough time to propagate.
	// An example would be having the following setup:
	// 1. Config with primary connection (con1) and secondary connection (con2)
	// 2. Setting secondary connection (con2) as primary
	// 3. Setting previously primary connection (con1) as primary (we will have the connection dependencies as in the first step)
	// 4. Deleting secondary connection (con2; this may fail without waiting a bit for Snowflake to propagate the changes)
	if err := util.Retry(3, time.Second, func() (error, bool) {
		if err := client.Connections.DropSafely(ctx, id); err != nil && strings.Contains(err.Error(), "is currently a primary connection in a replication relationship and cannot be dropped") {
			return nil, false
		} else {
			log.Printf("[DEBUG] Drop secondary connection failed, err = %s", err)
			return err, true
		}
	}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
