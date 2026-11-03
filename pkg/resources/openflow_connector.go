package resources

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	openflowConnectorCreateTimeout = 30 * time.Minute
	openflowConnectorUpdateTimeout = 30 * time.Minute
	openflowConnectorDeleteTimeout = 30 * time.Minute
)

var openflowConnectorTimeouts = &schema.ResourceTimeout{
	Create: schema.DefaultTimeout(openflowConnectorCreateTimeout),
	Read:   schema.DefaultTimeout(defaultReadTimeout),
	Update: schema.DefaultTimeout(openflowConnectorUpdateTimeout),
	Delete: schema.DefaultTimeout(openflowConnectorDeleteTimeout),
}

var openflowConnectorSchema = map[string]*schema.Schema{
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the Openflow connector."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the Openflow connector."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the Openflow connector; must be unique for the schema in which the connector is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"runtime": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Specifies the fully qualified name of the Openflow runtime the connector runs in. The connector is created in the runtime's schema, so `database` and `schema` must match it. Snowflake has no ALTER for it, so changing it recreates the connector.",
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"from": {
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		MinItems: 1,
		MaxItems: 1,
		// Accepted configurations:
		// - definition
		// - stage, and optional path
		// Required, because CREATE always needs one of the two sources. ExactlyOneOf rejects a block setting
		// both or neither, but does not fire when the block is absent, so leaving it Optional would let a
		// configuration without it plan cleanly and fail only once CREATE runs.
		Description: "Specifies what the connector is created from. Snowflake has no ALTER for it, so changing it recreates the connector. Note that external changes on this field and nested fields are not detected: Snowflake resolves a definition for a connector created from a stage too, so a value read from SHOW could not be told apart from a configured one. `show_output.connector_definition` reports what Snowflake resolved.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"definition": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					ExactlyOneOf: []string{"from.0.definition", "from.0.stage"},
					Description:  "Catalog definition ID for the connector type, for example `OPENFLOW_POSTGRES_CDC`. List the available IDs with the `snowflake_openflow_connector_definitions` data source. A connector created this way is a draft: it settles on STOPPED and stays there until a configuration version is committed, which this resource does not do.",
				},
				// A configured connector can only be created this way: configuration is not an argument of
				// CREATE OPENFLOW CONNECTOR, it is a config.json in a bundle Snowflake reads from a stage.
				"stage": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					ExactlyOneOf:     []string{"from.0.definition", "from.0.stage"},
					Description:      "Identifier of a stage holding a complete configuration bundle, which is how a connector arrives already configured and able to start without a commit. A git repository stage works here too.",
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
					DiffSuppressFunc: suppressIdentifierQuoting,
				},
				"path": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					RequiredWith: []string{"from.0.stage"},
					Description:  "Path to the bundle within the stage. The bundle's root is used when omitted.",
				},
			},
		},
	},
	"version_trigger": {
		Type:         schema.TypeString,
		Optional:     true,
		RequiredWith: []string{"from.0.stage"},
		Description:  "An arbitrary value, such as a hash of the bundle's files, that, if and only if changed from a non-empty to a different non-empty value (or known after apply), loads the bundle at `from` again as a new configuration version of the same connector, rather than recreating it. When you add this field to the configuration, or remove it, nothing is loaded, so add it before the bundle changes. A running connector is stopped for the load and started again afterwards; a stopped one stays stopped. When a version is loaded, `show_output` and `describe_output` are marked as computed. Only for a connector created from a stage.",
	},
	"display_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A free-text alias for the connector. Shown in the Openflow UI in place of the connector's identifier when set.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the Openflow connector.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW OPENFLOW CONNECTORS` for the given connector.",
		Elem: &schema.Resource{
			Schema: schemas.ShowOpenflowConnectorSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE OPENFLOW CONNECTOR` for the given connector.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeOpenflowConnectorDetailsSchema,
		},
	},
}

func OpenflowConnector() *schema.Resource {
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.OpenflowConnectorResource), TrackingCreateWrapper(resources.OpenflowConnector, CreateOpenflowConnector)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.OpenflowConnectorResource), TrackingReadWrapper(resources.OpenflowConnector, ReadOpenflowConnector(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.OpenflowConnectorResource), TrackingUpdateWrapper(resources.OpenflowConnector, UpdateOpenflowConnector)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.OpenflowConnectorResource), TrackingDeleteWrapper(resources.OpenflowConnector, DeleteOpenflowConnector)),
		Description: joinWithSpace(
			"Resource used to manage Openflow connectors, which run inside an Openflow runtime.",
			"Every mutating statement is asynchronous, so create and update return once the connector settles.",
			"Starting and stopping are operational actions and are not exposed here, although the provider may execute such commands when altering the object; the connector's state is reported in `show_output`.",
			"For more information, check [Openflow connector documentation](https://docs.snowflake.com/en/sql-reference/sql/create-openflow-connector).",
		),

		CustomizeDiff: TrackingCustomDiffWrapper(resources.OpenflowConnector, customdiff.All(
			ComputedIfAnyAttributeChanged(openflowConnectorSchema, ShowOutputAttributeName, "name", "display_name", "comment"),
			ComputedIfAnyAttributeChanged(openflowConnectorSchema, DescribeOutputAttributeName, "name", "display_name", "comment"),
			ComputedIfAnyAttributeChanged(openflowConnectorSchema, FullyQualifiedNameAttributeName, "name"),
			func(_ context.Context, diff *schema.ResourceDiff, _ any) error {
				// A new version changes the version columns of both outputs.
				o, n := diff.GetChange("version_trigger")
				if keeperChanged(o.(string), n.(string), diff.GetRawPlan().AsValueMap()["version_trigger"].IsKnown()) {
					return errors.Join(
						diff.SetNewComputed(ShowOutputAttributeName),
						diff.SetNewComputed(DescribeOutputAttributeName),
					)
				}
				return nil
			},
		)),

		Schema:   openflowConnectorSchema,
		Timeouts: openflowConnectorTimeouts,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.OpenflowConnector, ImportOpenflowConnector),
		},
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
		d.Set("database", id.DatabaseName()),
		d.Set("schema", id.SchemaName()),
		d.Set("name", id.Name()),
		d.Set("runtime", sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), connector.Runtime).FullyQualifiedName()),
		// `from` is deliberately not set. SHOW reports a definition whichever source the connector was created
		// from, so importing one would write a definition for a stage-created connector, which its
		// configuration cannot hold next to `stage`. The first plan after an import therefore asks to replace
		// the connector, since `from` is create-only.
	)
	if connector.DisplayName != nil && *connector.DisplayName != "" {
		errs = errors.Join(errs, d.Set("display_name", *connector.DisplayName))
	}
	if connector.Comment != nil && *connector.Comment != "" {
		errs = errors.Join(errs, d.Set("comment", *connector.Comment))
	}
	if errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

// openflowConnectorBundleLocation is the stage location of the bundle a `from` block with a stage names.
func openflowConnectorBundleLocation(from map[string]any) (sdk.Location, error) {
	stage, err := sdk.ParseSchemaObjectIdentifier(from["stage"].(string))
	if err != nil {
		return nil, err
	}
	// Snowflake resolves the bundle's config.json under the path only when the path ends in a slash, and
	// rejects the create with "No connector configuration file found at location" without one. The
	// trailing slash is added here rather than being asked of the configuration, so that state keeps the
	// path as written.
	path := from["path"].(string)
	if path != "" && !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return sdk.NewStageLocation(stage, path), nil
}

func CreateOpenflowConnector(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := sdk.NewSchemaObjectIdentifier(d.Get("database").(string), d.Get("schema").(string), d.Get("name").(string))

	runtime, err := sdk.ParseSchemaObjectIdentifier(d.Get("runtime").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateOpenflowConnectorRequest(id, runtime)
	if err := errors.Join(
		stringAttributeCreateBuilder(d, "display_name", request.WithDisplayName),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
	); err != nil {
		return diag.FromErr(err)
	}

	// ExactlyOneOf on the two sources means the block is present with one of them filled, so only their
	// emptiness is checked here.
	from := d.Get("from").([]any)[0].(map[string]any)
	if definition := from["definition"].(string); definition != "" {
		request.WithFromDefinition(definition)
	}
	if from["stage"].(string) != "" {
		location, err := openflowConnectorBundleLocation(from)
		if err != nil {
			return diag.FromErr(err)
		}
		request.WithFrom(location)
	}

	if err := client.OpenflowConnectors.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	if err := waitForOpenflowConnectorReady(ctx, client, id, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}
	return ReadOpenflowConnector(false)(ctx, d, meta)
}

func ReadOpenflowConnector(withExternalChangesMarking bool) schema.ReadContextFunc {
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
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query openflow connector. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Openflow connector: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		connectorDetails, err := client.OpenflowConnectors.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		// Flattened up front: the drift comparison is an `any` equality against the string in show_output, so a
		// *string would never match and would mark drift on every read.
		var displayName, comment string
		if connector.DisplayName != nil {
			displayName = *connector.DisplayName
		}
		if connector.Comment != nil {
			comment = *connector.Comment
		}

		if withExternalChangesMarking {
			if err = handleExternalChangesToObjectInShow(
				d,
				outputMapping{"display_name", "display_name", displayName, displayName, nil},
				outputMapping{"comment", "comment", comment, comment, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, openflowConnectorSchema, []string{
			"display_name",
			"comment",
		}); err != nil {
			return diag.FromErr(err)
		}

		if errs := errors.Join(
			// Read back so that an external change to a ForceNew attribute diverges from config and forces
			// replacement on the next plan. SHOW reports the runtime by name only, and a connector can only be
			// created in its runtime's schema, so the qualified name is rebuilt from the connector's own.
			d.Set("runtime", sdk.NewSchemaObjectIdentifier(id.DatabaseName(), id.SchemaName(), connector.Runtime).FullyQualifiedName()),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.OpenflowConnectorToSchema(connector)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.OpenflowConnectorDetailsToSchema(connectorDetails)}),
		); errs != nil {
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
	timeout := d.Timeout(schema.TimeoutUpdate)

	// ALTER is refused while the connector is in a transient status, so an apply that lands while another
	// change is still in flight has to let it settle first.
	if err := waitForOpenflowConnectorSettled(ctx, client, id, timeout); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), d.Get("name").(string))
		if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithRenameTo(newId)); err != nil {
			return diag.FromErr(err)
		}
		if err := waitForOpenflowConnectorReady(ctx, client, newId, timeout); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}

	set, unset := sdk.NewOpenflowConnectorSetRequest(), sdk.NewOpenflowConnectorUnsetRequest()
	if err := errors.Join(
		stringAttributeUpdate(d, "display_name", &set.DisplayName, &unset.DisplayName),
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
	); err != nil {
		return diag.FromErr(err)
	}

	if (*set != sdk.OpenflowConnectorSetRequest{}) {
		if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
		if err := waitForOpenflowConnectorReady(ctx, client, id, timeout); err != nil {
			return diag.FromErr(err)
		}
	}
	if (*unset != sdk.OpenflowConnectorUnsetRequest{}) {
		if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
		if err := waitForOpenflowConnectorReady(ctx, client, id, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	if o, n := d.GetChange("version_trigger"); keeperChanged(o.(string), n.(string), d.GetRawPlan().AsValueMap()["version_trigger"].IsKnown()) {
		if err := addOpenflowConnectorVersionFromBundle(ctx, client, id, d.Get("from").([]any)[0].(map[string]any), timeout); err != nil {
			// Without this the new trigger would be saved even though no version was added, and the next plan
			// would show nothing left to do.
			d.Partial(true)
			return diag.FromErr(err)
		}
	}

	return ReadOpenflowConnector(false)(ctx, d, meta)
}

// addOpenflowConnectorVersionFromBundle loads the bundle again with ADD VERSION, which makes it the new
// default version of the same connector. Snowflake accepts ADD VERSION only from STOPPED, so a running
// connector is stopped around it and started again, leaving it in the status it was found in. A failure after
// the STOP leaves the connector stopped, for an operator to look at rather than to have it started on a
// version that may not have loaded.
func addOpenflowConnectorVersionFromBundle(ctx context.Context, client *sdk.Client, id sdk.SchemaObjectIdentifier, from map[string]any, timeout time.Duration) error {
	location, err := openflowConnectorBundleLocation(from)
	if err != nil {
		return err
	}

	connector, err := client.OpenflowConnectors.ShowByID(ctx, id)
	if err != nil {
		return err
	}
	var wasRunning bool
	switch connector.Status {
	case sdk.OpenflowConnectorStatusRunning:
		wasRunning = true
		if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithStop(true)); err != nil {
			return err
		}
		if err := waitForOpenflowConnectorStopped(ctx, client, id, timeout); err != nil {
			return err
		}
	case sdk.OpenflowConnectorStatusStopped:
	default:
		// A failed or troubleshooting connector needs an operator; starting it again afterwards would hide that.
		return fmt.Errorf("openflow connector %s is %s; a new version can only be added to a RUNNING or STOPPED connector", id.FullyQualifiedName(), connector.Status)
	}

	if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithAddVersion(*sdk.NewOpenflowConnectorAddVersionRequest(&location))); err != nil {
		return err
	}
	if err := waitForOpenflowConnectorReady(ctx, client, id, timeout); err != nil {
		return err
	}

	if wasRunning {
		if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithStart(true)); err != nil {
			return err
		}
		if err := waitForOpenflowConnectorRunning(ctx, client, id, timeout); err != nil {
			return err
		}
	}
	return nil
}

func DeleteOpenflowConnector(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	timeout := d.Timeout(schema.TimeoutDelete)

	connector, err := client.OpenflowConnectors.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	switch connector.Status {
	case sdk.OpenflowConnectorStatusDeleting, sdk.OpenflowConnectorStatusDeleted:
		// Already terminated or on its way; fall through to the wait and drop. DELETE_FAILED is deliberately
		// not here: a terminate that failed has to be issued again, which the default branch does.
	default:
		if err := waitForOpenflowConnectorSettled(ctx, client, id, timeout); err != nil {
			return diag.FromErr(err)
		}
		// The status has to be re-read after the wait, because the wait can settle on a failed status, which
		// is torn down differently from a clean one.
		settled, err := client.OpenflowConnectors.ShowByIDSafely(ctx, id)
		switch {
		case errors.Is(err, sdk.ErrObjectNotFound):
			d.SetId("")
			return nil
		case err != nil:
			return diag.FromErr(err)
		}
		// DELETE_FAILED is the one failure status STOP does not clear: Snowflake refuses STOP from it and
		// accepts TERMINATE, so it goes straight to the terminate below.
		if slices.Contains(sdk.OpenflowConnectorFailureStatuses, settled.Status) &&
			settled.Status != sdk.OpenflowConnectorStatusDeleteFailed {
			if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithIfExists(true).WithStop(true)); err != nil {
				return diag.FromErr(err)
			}
			if err := waitForOpenflowConnectorStopped(ctx, client, id, timeout); err != nil {
				return diag.FromErr(err)
			}
		}
		if err := client.OpenflowConnectors.Alter(ctx, sdk.NewAlterOpenflowConnectorRequest(id).WithIfExists(true).WithTerminate(true)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := waitForOpenflowConnectorTerminated(ctx, client, id, timeout); err != nil {
		return diag.FromErr(err)
	}
	if err := client.OpenflowConnectors.DropSafely(ctx, id); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
