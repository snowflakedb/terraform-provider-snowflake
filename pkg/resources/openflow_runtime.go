package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var openflowRuntimeSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Identifier for the openflow runtime."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Database in which to create the runtime."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Schema in which to create the runtime."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"deployment": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "Name of the openflow deployment this runtime belongs to.",
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"node_type": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateDiagFunc: sdkValidation(sdk.ToOpenflowRuntimeNodeType),
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompare(sdk.ToOpenflowRuntimeNodeType)),
		Description:      fmt.Sprintf("Node type for the runtime. Valid values: %s.", possibleValuesListed(sdk.AllOpenflowRuntimeNodeTypes)),
	},
	"min_nodes": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Minimum number of nodes for the runtime.",
	},
	"max_nodes": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Maximum number of nodes for the runtime.",
	},
	"execute_as_role": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Role to execute the runtime as.",
	},
	"external_access_integrations": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "List of external access integration names (SNOWFLAKE deployments only).",
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"display_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Human-readable display name for the runtime.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Comment for the runtime.",
	},
	"suspended": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "When true the runtime is suspended. When false it is resumed.",
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Output of SHOW OPENFLOW RUNTIMES for this runtime.",
		Elem:        &schema.Resource{Schema: schemas.ShowOpenflowRuntimeSchema},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Output of DESCRIBE OPENFLOW RUNTIME for this runtime.",
		Elem:        &schema.Resource{Schema: schemas.DescribeOpenflowRuntimeSchema},
	},
}

func OpenflowRuntime() *schema.Resource {
	return &schema.Resource{
		CreateContext: TrackingCreateWrapper(resources.OpenflowRuntime, CreateOpenflowRuntime),
		ReadContext:   TrackingReadWrapper(resources.OpenflowRuntime, ReadOpenflowRuntimeFunc(true)),
		UpdateContext: TrackingUpdateWrapper(resources.OpenflowRuntime, UpdateOpenflowRuntime),
		DeleteContext: TrackingDeleteWrapper(resources.OpenflowRuntime, DeleteOpenflowRuntime),
		Description:   "Manages an Openflow Runtime. Runtimes are schema-scoped objects created within an Openflow Deployment.",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.OpenflowRuntime, customdiff.All(
			ComputedIfAnyAttributeChanged(openflowRuntimeSchema, ShowOutputAttributeName,
				"min_nodes", "max_nodes", "execute_as_role", "external_access_integrations", "comment", "display_name", "suspended"),
			ComputedIfAnyAttributeChanged(openflowRuntimeSchema, DescribeOutputAttributeName,
				"min_nodes", "max_nodes", "execute_as_role", "external_access_integrations", "comment", "display_name", "suspended"),
		)),

		Schema: openflowRuntimeSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.OpenflowRuntime, ImportOpenflowRuntime),
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportOpenflowRuntime(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	runtime, err := client.OpenflowRuntimes.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}
	errs := errors.Join(
		d.Set("name", runtime.Name),
		d.Set("database", runtime.DatabaseName),
		d.Set("schema", runtime.SchemaName),
		d.Set("deployment", runtime.Deployment),
		d.Set("node_type", string(runtime.NodeType)),
		d.Set("min_nodes", runtime.MinNodes),
		d.Set("max_nodes", runtime.MaxNodes),
		d.Set("execute_as_role", runtime.ExecuteAsRole),
		d.Set("suspended", runtime.Status == sdk.OpenflowRuntimeStatusSuspended),
	)
	if errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

func CreateOpenflowRuntime(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	database := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(database, schemaName, name)

	deploymentName := d.Get("deployment").(string)
	deploymentID := sdk.NewAccountObjectIdentifier(deploymentName)

	nodeType, err := sdk.ToOpenflowRuntimeNodeType(d.Get("node_type").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	request := sdk.NewCreateOpenflowRuntimeRequest(
		id,
		deploymentID,
		d.Get("execute_as_role").(string),
		nodeType,
		d.Get("min_nodes").(int),
		d.Get("max_nodes").(int),
	).WithIfNotExists(true)

	if v, ok := d.GetOk("external_access_integrations"); ok {
		eaiList := v.([]interface{})
		eais := make([]sdk.AccountObjectIdentifier, len(eaiList))
		for i, e := range eaiList {
			eais[i] = sdk.NewAccountObjectIdentifier(e.(string))
		}
		request.WithExternalAccessIntegrations(eais)
	}
	if v, ok := d.GetOk("display_name"); ok {
		request.WithDisplayName(v.(string))
	}
	if v, ok := d.GetOk("comment"); ok {
		request.WithComment(v.(string))
	}

	if err := client.OpenflowRuntimes.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(helpers.EncodeResourceIdentifier(id))

	if err := waitForOpenflowRuntimeActive(ctx, client, id); err != nil {
		return diag.FromErr(err)
	}

	// Handle initially_suspended equivalent: if suspended=true, suspend after creation.
	if d.Get("suspended").(bool) {
		if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithSuspend(true)); err != nil {
			return diag.FromErr(err)
		}
		if err := waitForOpenflowRuntimeSuspended(ctx, client, id); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadOpenflowRuntimeFunc(false)(ctx, d, meta)
}

func ReadOpenflowRuntimeFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		runtime, err := client.OpenflowRuntimes.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{{
					Severity: diag.Warning,
					Summary:  "Failed to query openflow runtime. Marking as removed.",
					Detail:   fmt.Sprintf("Runtime id: %s, Err: %s", id.FullyQualifiedName(), err),
				}}
			}
			return diag.FromErr(err)
		}

		details, err := client.OpenflowRuntimes.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		isSuspended := runtime.Status == sdk.OpenflowRuntimeStatusSuspended

		if withExternalChangesMarking {
			var displayName, comment string
			if runtime.DisplayName != nil {
				displayName = *runtime.DisplayName
			}
			if runtime.Comment != nil {
				comment = *runtime.Comment
			}
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"min_nodes", "min_nodes", runtime.MinNodes, runtime.MinNodes, nil},
				outputMapping{"max_nodes", "max_nodes", runtime.MaxNodes, runtime.MaxNodes, nil},
				outputMapping{"execute_as_role", "execute_as_role", runtime.ExecuteAsRole, runtime.ExecuteAsRole, nil},
				outputMapping{"display_name", "display_name", displayName, displayName, nil},
				outputMapping{"comment", "comment", comment, comment, nil},
				outputMapping{"suspended", "suspended", isSuspended, isSuspended, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, openflowRuntimeSchema, []string{
			"min_nodes", "max_nodes", "execute_as_role", "external_access_integrations",
			"display_name", "comment", "suspended",
		}); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.OpenflowRuntimeToSchema(runtime)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.OpenflowRuntimeDetailsToSchema(*details)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("min_nodes", runtime.MinNodes),
			d.Set("max_nodes", runtime.MaxNodes),
			d.Set("execute_as_role", runtime.ExecuteAsRole),
			d.Set("external_access_integrations", runtime.ExternalAccessIntegrations),
			d.Set("suspended", isSuspended),
		)
		if runtime.DisplayName != nil {
			errs = errors.Join(errs, d.Set("display_name", *runtime.DisplayName))
		}
		if runtime.Comment != nil {
			errs = errors.Join(errs, d.Set("comment", *runtime.Comment))
		}
		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateOpenflowRuntime(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle suspend/resume first.
	if d.HasChange("suspended") {
		suspended := d.Get("suspended").(bool)
		if suspended {
			if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithSuspend(true)); err != nil {
				return diag.FromErr(err)
			}
			if err := waitForOpenflowRuntimeSuspended(ctx, client, id); err != nil {
				return diag.FromErr(err)
			}
		} else {
			if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithResume(true)); err != nil {
				return diag.FromErr(err)
			}
			if err := waitForOpenflowRuntimeActive(ctx, client, id); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Handle SET for async-capable properties (triggers UPDATING state in CP).
	asyncSet := sdk.NewOpenflowRuntimeSetRequest()
	hasAsyncChanges := false
	if d.HasChange("min_nodes") {
		asyncSet.WithMinNodes(d.Get("min_nodes").(int))
		hasAsyncChanges = true
	}
	if d.HasChange("max_nodes") {
		asyncSet.WithMaxNodes(d.Get("max_nodes").(int))
		hasAsyncChanges = true
	}
	if d.HasChange("execute_as_role") {
		asyncSet.WithExecuteAsRole(d.Get("execute_as_role").(string))
		hasAsyncChanges = true
	}
	if d.HasChange("external_access_integrations") {
		eaiList := d.Get("external_access_integrations").([]interface{})
		eais := make([]sdk.AccountObjectIdentifier, len(eaiList))
		for i, e := range eaiList {
			eais[i] = sdk.NewAccountObjectIdentifier(e.(string))
		}
		asyncSet.WithExternalAccessIntegrations(eais)
		hasAsyncChanges = true
	}
	if hasAsyncChanges {
		if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithSet(*asyncSet)); err != nil {
			return diag.FromErr(err)
		}
		if err := waitForOpenflowRuntimeActive(ctx, client, id); err != nil {
			return diag.FromErr(err)
		}
	}

	// Metadata-only SET (no CP task).
	metaSet := sdk.NewOpenflowRuntimeSetRequest()
	metaUnset := sdk.NewOpenflowRuntimeUnsetRequest()
	hasMetaSet, hasMetaUnset := false, false
	if d.HasChange("display_name") {
		if v, ok := d.GetOk("display_name"); ok {
			metaSet.WithDisplayName(v.(string))
			hasMetaSet = true
		} else {
			metaUnset.WithDisplayName(true)
			hasMetaUnset = true
		}
	}
	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			metaSet.WithComment(v.(string))
			hasMetaSet = true
		} else {
			metaUnset.WithComment(true)
			hasMetaUnset = true
		}
	}
	if hasMetaSet {
		if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithSet(*metaSet)); err != nil {
			return diag.FromErr(err)
		}
	}
	if hasMetaUnset {
		if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithUnset(*metaUnset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadOpenflowRuntimeFunc(false)(ctx, d, meta)
}

func DeleteOpenflowRuntime(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Multi-step delete: suspend → wait SUSPENDED → terminate → drop.
	runtime, err := client.OpenflowRuntimes.ShowByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if runtime.Status != sdk.OpenflowRuntimeStatusSuspended {
		if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithSuspend(true)); err != nil {
			return diag.FromErr(err)
		}
		if err := waitForOpenflowRuntimeSuspended(ctx, client, id); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithTerminate(true)); err != nil {
		// If terminate fails, try cascade.
		if strings.Contains(err.Error(), "connector") {
			if err2 := client.OpenflowRuntimes.Alter(ctx, sdk.NewAlterOpenflowRuntimeRequest(id).WithTerminateCascade(true)); err2 != nil {
				return diag.FromErr(fmt.Errorf("terminate failed (%s); terminate cascade also failed: %w", err, err2))
			}
		} else {
			return diag.FromErr(err)
		}
	}

	if err := waitForOpenflowRuntimeDeleted(ctx, client, id); err != nil {
		return diag.FromErr(err)
	}

	if err := client.OpenflowRuntimes.DropSafely(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
