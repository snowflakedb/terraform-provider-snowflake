package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var openflowConnectorSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Identifier for the openflow connector."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Database in which to create the connector."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Schema in which to create the connector (must match the runtime's schema)."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"runtime": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Name of the openflow runtime this connector belongs to.",
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"from_definition": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		ExactlyOneOf: []string{"from_definition", "from_stage"},
		Description:  "Name of a known connector definition to install from.",
	},
	"from_stage": {
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		ExactlyOneOf: []string{"from_definition", "from_stage"},
		Description:  "Stage path to install the connector from (e.g. @MY_REPO/branches/main/path/).",
	},
	"display_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Human-readable display name for the connector.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Comment for the connector.",
	},
	"started": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "When true the connector is started. When false it is stopped.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Output of SHOW OPENFLOW CONNECTORS for this connector.",
		Elem:        &schema.Resource{Schema: schemas.ShowOpenflowConnectorSchema},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Output of DESCRIBE OPENFLOW CONNECTOR for this connector.",
		Elem:        &schema.Resource{Schema: schemas.DescribeOpenflowConnectorSchema},
	},
}

func OpenflowConnector() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.OpenflowConnector, CreateOpenflowConnector),
		ReadContext:   TrackingReadWrapper(resources.OpenflowConnector, ReadOpenflowConnectorFunc(true)),
		UpdateContext: TrackingUpdateWrapper(resources.OpenflowConnector, UpdateOpenflowConnector),
		DeleteContext: TrackingDeleteWrapper(resources.OpenflowConnector, DeleteOpenflowConnector),
		Description:   "Manages an Openflow Connector. Connectors are schema-scoped objects created within an Openflow Runtime.",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.OpenflowConnector, customdiff.All(
			ComputedIfAnyAttributeChanged(openflowConnectorSchema, ShowOutputAttributeName, "comment", "display_name", "started"),
			ComputedIfAnyAttributeChanged(openflowConnectorSchema, DescribeOutputAttributeName, "comment", "display_name", "started"),
		)),

		Schema: openflowConnectorSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.OpenflowConnector, ImportOpenflowConnector),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportOpenflowConnector(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	connector, err := client.OpenflowConnectors.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	errs := errors.Join(
		d.Set("name", connector.Name),
		d.Set("database", connector.DatabaseName),
		d.Set("schema", connector.SchemaName),
		d.Set("runtime", connector.Runtime),
		d.Set("started", connector.Started),
	)
	if connector.ConnectorDefinition != nil {
		errs = errors.Join(errs, d.Set("from_definition", *connector.ConnectorDefinition))
	}
	if errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

func CreateOpenflowConnector(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	database := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	runtimeName := d.Get("runtime").(string)
	runtimeID := sdk.NewSchemaObjectIdentifier(database, schemaName, runtimeName)

	request := sdk.NewCreateOpenflowConnectorRequest(id, runtimeID).WithIfNotExists(true)

	if v, ok := d.GetOk("from_definition"); ok {
		request.WithFromDefinition(v.(string))
	}
	if v, ok := d.GetOk("from_stage"); ok {
		request.WithFromStage(v.(string))
	}
	if v, ok := d.GetOk("display_name"); ok {
		request.WithDisplayName(v.(string))
	}
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	if err := client.OpenflowConnectors.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	if err := waitForOpenflowConnectorActive(ctx, client, id); err != nil {
		return diag.FromErr(err)
	}

	// Handle started=true: start the connector after creation.
	if d.Get("started").(bool) {
		if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithStart(true)); err != nil {
			return diag.FromErr(err)
		}
		if err := waitForOpenflowConnectorActive(ctx, client, id); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadOpenflowConnectorFunc(false)(ctx, d, meta)
}

func ReadOpenflowConnectorFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		connector, err := client.OpenflowConnectors.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{{
					Severity: diag.Warning,
					Summary:  "Failed to query openflow connector. Marking as removed.",
					Detail:   fmt.Sprintf("Connector id: %s, Err: %s", id.FullyQualifiedName(), err),
				}}
			}
			return diag.FromErr(err)
		}

		details, err := client.OpenflowConnectors.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			var displayName, comment string
			if connector.DisplayName != nil {
				displayName = *connector.DisplayName
			}
			if connector.Comment != nil {
				comment = *connector.Comment
			}
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"started", "started", connector.Started, connector.Started, nil},
				outputMapping{"display_name", "display_name", displayName, displayName, nil},
				outputMapping{"comment", "comment", comment, comment, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, openflowConnectorSchema, []string{
			"started", "display_name", "comment",
		}); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.OpenflowConnectorToSchema(connector)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.OpenflowConnectorDetailsToSchema(*details)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("started", connector.Started),
		)
		if connector.DisplayName != nil {
			errs = errors.Join(errs, d.Set("display_name", *connector.DisplayName))
		}
		if connector.Comment != nil {
			errs = errors.Join(errs, d.Set("comment", *connector.Comment))
		}
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateOpenflowConnector(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle start/stop.
	if d.HasChange("started") {
		started := d.Get("started").(bool)
		if started {
			if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithStart(true)); err != nil {
				return diag.FromErr(err)
			}
			if err := waitForOpenflowConnectorActive(ctx, client, id); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithStop(true)); err != nil {
				return diag.FromErr(err)
			}
			if err := waitForOpenflowConnectorStopped(ctx, client, id); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	set := sdk.NewOpenflowConnectorSetRequest()
	unset := sdk.NewOpenflowConnectorUnsetRequest()
	hasSet, hasUnset := false, false

	if d.HasChange("display_name") {
		if v, ok := d.GetOk("display_name"); ok {
			set.WithDisplayName(v.(string))
			hasSet = true
		} else {
			unset.WithDisplayName(true)
			hasUnset = true
		}
	}
	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.WithComment(v.(string))
			hasSet = true
		} else {
			unset.WithComment(true)
			hasUnset = true
		}
	}
	if hasSet {
		if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}
	if hasUnset {
		if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadOpenflowConnectorFunc(false)(ctx, d, meta)
}

func DeleteOpenflowConnector(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Stop connector before dropping if it's running.
	connector, err := client.OpenflowConnectors.ShowByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if connector.Started {
		if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithStop(true)); err != nil {
			return diag.FromErr(err)
		}
		if err := waitForOpenflowConnectorStopped(ctx, client, id); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := client.OpenflowConnectors.DropSafely(ctx, id); err != nil {
		return diag.FromErr(err)
	}
	if err := waitForOpenflowConnectorDeleted(ctx, client, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
