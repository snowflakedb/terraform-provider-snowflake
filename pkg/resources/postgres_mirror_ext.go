package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Constructor
// =============================================================================

var (
	postgresMirrorCreateExt      = PreviewFeatureCreateContextWrapper(string(previewfeatures.PostgresMirrorResource), TrackingCreateWrapper(resources.PostgresMirror, CreatePostgresMirror))
	postgresMirrorReadExt        = PreviewFeatureReadContextWrapper(string(previewfeatures.PostgresMirrorResource), TrackingReadWrapper(resources.PostgresMirror, ReadPostgresMirrorFunc(true)))
	postgresMirrorUpdateExt      = PreviewFeatureUpdateContextWrapper(string(previewfeatures.PostgresMirrorResource), TrackingUpdateWrapper(resources.PostgresMirror, UpdatePostgresMirror))
	postgresMirrorDeleteExt      = PreviewFeatureDeleteContextWrapper(string(previewfeatures.PostgresMirrorResource), TrackingDeleteWrapper(resources.PostgresMirror, DeletePostgresMirror))
	postgresMirrorDescriptionExt = "Resource used to manage Postgres data mirrors, which replicate Postgres tables into a Snowflake database. For more information, check [data mirroring documentation](https://docs.snowflake.com/en/user-guide/snowflake-postgres/postgres-data-mirroring)."

	postgresMirrorCustomizeDiffExt = TrackingCustomDiffWrapper(resources.PostgresMirror, customdiff.All(
		ComputedIfAnyAttributeChanged(postgresMirrorSchema, DescribeOutputAttributeName, "postgres_tables", "postgres_schemas", "refresh_interval", "warehouse"),
	))

	postgresMirrorImporterExt = &schema.ResourceImporter{
		StateContext: TrackingImportWrapper(resources.PostgresMirror, ImportPostgresMirror),
	}

	// target_table_type stays as configured. Snowflake's SNOWFLAKE default is not copied into state,
	// because the field is ForceNew and an omitted value must not replace the mirror.
	// warehouse and external_volume are written from DESCRIBE_MIRROR in Read. They have no Snowflake default.
	postgresMirrorSetStateToValueFromConfigFieldsExt = []string{
		"target_table_type",
	}
)

// =============================================================================
// Schema
// =============================================================================

var postgresMirrorIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: ignoreCaseSuppressFunc,
		Description:      "Specifies the name of the mirror. It must be unique across all Postgres instances in the account. Snowflake stores the name lowercased, and the lowercased name is the resource ID.",
	},
}

var postgresMirrorDescribeOutputSchemaExt = map[string]*schema.Schema{
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE_MIRROR` for the given Postgres mirror.",
		Elem: &schema.Resource{
			Schema: schemas.DescribePostgresMirrorDetailsSchema,
		},
	},
}

// The procedures uppercase the target database and the warehouse, so those are compared ignoring case.
// Postgres table and schema names are not Snowflake identifiers, so the generator leaves them unchecked.
func postgresMirrorAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["target_database"].DiffSuppressFunc = SuppressIfAny(s["target_database"].DiffSuppressFunc, ignoreCaseSuppressFunc)
	s["warehouse"].DiffSuppressFunc = SuppressIfAny(s["warehouse"].DiffSuppressFunc, ignoreCaseSuppressFunc)

	// Snowflake fills the interval when the block is omitted, and removing the block keeps the current interval.
	s["refresh_interval"].Computed = true
	s["drop_target_database_on_destroy"].Default = true
	return s
}

// =============================================================================
// Import
// =============================================================================

func ImportPostgresMirror(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	name := strings.ToLower(d.Id())

	details, err := client.PostgresMirrors.Describe(ctx, sdk.NewDescribePostgresMirrorRequest(*sdk.NewDescribePostgresMirrorArgumentsRequest(name)))
	if err != nil {
		return nil, fmt.Errorf("could not describe Postgres mirror %s: %w", name, err)
	}
	if !details.IsMembershipKnown() {
		return nil, fmt.Errorf("DESCRIBE_MIRROR returned no tables or schemas for Postgres mirror %s, either because the source Postgres instance could not be reached or because the SNOWFLAKE application in this account does not report them yet; retry the import later", name)
	}

	// Tables, schemas, and the refresh interval are read from Snowflake in both modes: Read writes them too.
	// drop_target_database_on_destroy is not a Snowflake value. The schema default is true, and schema defaults
	// are not applied to imported state, so both modes write true.
	// warehouse and external_volume have no Snowflake default. Read writes them, including an empty value when unset.
	importPermissiveMode := meta.(*provider.Context).Experiments.IsEnabled(experimentalfeatures.ImportPermissiveMode)
	if errs := errors.Join(
		d.Set("name", details.Name),
		d.Set("postgres_instance", details.PostgresInstance),
		d.Set("postgres_database", details.PostgresDatabase),
		d.Set("target_database", details.TargetDatabase),
		d.Set("postgres_tables", details.PostgresTables),
		d.Set("postgres_schemas", collections.Map(details.PostgresSchemas, unquotePostgresIdentifier)),
		d.Set("drop_target_database_on_destroy", true),
		setPostgresMirrorRefreshInterval(d, details.RefreshInterval),
		importStringAttribute(d, "target_table_type", importPermissiveMode, details.TargetTableType),
	); errs != nil {
		return nil, errs
	}
	d.SetId(details.Name)
	return []*schema.ResourceData{d}, nil
}

// =============================================================================
// Create
// =============================================================================

func postgresMirrorParseIdFromConfigExt(d *schema.ResourceData) (sdk.AccountObjectIdentifier, error) {
	return sdk.NewAccountObjectIdentifier(strings.ToLower(d.Get("name").(string))), nil
}

func postgresMirrorNewCreateRequestExt(d *schema.ResourceData, id sdk.AccountObjectIdentifier) (*sdk.CreatePostgresMirrorArgumentsRequest, error) {
	instanceId, err := sdk.ParseAccountObjectIdentifier(d.Get("postgres_instance").(string))
	if err != nil {
		return nil, err
	}
	targetDatabaseId, err := sdk.ParseAccountObjectIdentifier(d.Get("target_database").(string))
	if err != nil {
		return nil, err
	}

	arguments := sdk.NewCreatePostgresMirrorArgumentsRequest(id.Name(), instanceId.Name(), d.Get("postgres_database").(string), targetDatabaseId.Name())
	if err := attributeMappedValueCreateBuilder(d, "refresh_interval", arguments.WithRefreshInterval, func(v []any) (string, error) {
		interval, _ := postgresMirrorRefreshInterval(v)
		return interval, nil
	}); err != nil {
		return nil, err
	}
	return arguments, nil
}

func toPostgresMirrorName(value any) (sdk.StringListItemWrapper, error) {
	return sdk.StringListItemWrapper{Value: value.(string)}, nil
}

func postgresMirrorCreateInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier, arguments *sdk.CreatePostgresMirrorArgumentsRequest) diag.Diagnostics {
	if err := meta.(*provider.Context).Client.PostgresMirrors.Create(ctx, sdk.NewCreatePostgresMirrorRequest(*arguments)); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Postgres mirror %s: %w", id.Name(), err))
	}
	return nil
}

// =============================================================================
// Read
// =============================================================================

// A mirror has no SHOW, so DESCRIBE_MIRROR is the existence check and the source of the state.
// A missing mirror is removed from state. When membership cannot be read, Read stops and keeps the previous state.
func postgresMirrorShowByIdSafelyInSdkExt(ctx context.Context, d *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier) (*sdk.PostgresMirrorDetails, diag.Diagnostics) {
	details, err := meta.(*provider.Context).Client.PostgresMirrors.Describe(ctx, sdk.NewDescribePostgresMirrorRequest(*sdk.NewDescribePostgresMirrorArgumentsRequest(id.Name())))
	if err != nil {
		// DESCRIBE_MIRROR returns no row for a mirror that does not exist.
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return nil, diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "Failed to query Postgres mirror. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Postgres mirror name: %s, Err: %s", id.Name(), err),
				},
			}
		}
		return nil, diag.FromErr(err)
	}
	if !details.IsMembershipKnown() {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  "Could not read the tables and schemas of the Postgres mirror. Keeping the previously known values.",
				Detail: fmt.Sprintf("DESCRIBE_MIRROR returned no tables or schemas for Postgres mirror %s, either because the source Postgres instance could not be reached or because the SNOWFLAKE application in this account does not report them yet. "+
					"Changes made to them outside of Terraform are not detected until they can be read again.", id.Name()),
			},
		}
	}
	return details, nil
}

func postgresMirrorHandleExternalChangesExt(d *schema.ResourceData, details *sdk.PostgresMirrorDetails, withExternalChangesMarking bool) error {
	if !withExternalChangesMarking {
		return nil
	}
	return handleExternalChangesToObjectInFlatDescribe(d,
		outputMapping{"postgres_instance", "postgres_instance", details.PostgresInstance, details.PostgresInstance, nil},
		outputMapping{"postgres_database", "postgres_database", details.PostgresDatabase, details.PostgresDatabase, nil},
		outputMapping{"target_database", "target_database", details.TargetDatabase, details.TargetDatabase, nil},
		outputMapping{"target_table_type", "target_table_type", string(details.TargetTableType), string(details.TargetTableType), nil},
	)
}

func postgresMirrorSetConfigFieldsExt(_ context.Context, d *schema.ResourceData, _ any, details *sdk.PostgresMirrorDetails) error {
	return errors.Join(
		d.Set("postgres_tables", details.PostgresTables),
		d.Set("postgres_schemas", collections.Map(details.PostgresSchemas, unquotePostgresIdentifier)),
		d.Set("external_volume", util.DerefOrZero(details.ExternalVolume)),
		d.Set("warehouse", util.DerefOrZero(details.Warehouse)),
		setPostgresMirrorRefreshInterval(d, details.RefreshInterval),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.PostgresMirrorDetailsToSchema(details)}),
	)
}

// =============================================================================
// Update
// =============================================================================

type postgresMirrorAlterRequestsExt struct {
	set   *sdk.AlterPostgresMirrorArgumentsRequest
	unset *sdk.AlterPostgresMirrorArgumentsRequest
}

// ALTER_MIRROR takes additions and removals rather than the full list, and fails when asked to add a table the
// mirror already has, so the lists are computed from the state Read refreshed against the source.
func postgresMirrorApplySetUnsetExt(d *schema.ResourceData) (*postgresMirrorAlterRequestsExt, error) {
	reqs := &postgresMirrorAlterRequestsExt{
		set:   sdk.NewAlterPostgresMirrorArgumentsRequest(d.Id()),
		unset: sdk.NewAlterPostgresMirrorArgumentsRequest(d.Id()),
	}

	if d.HasChange("postgres_tables") {
		oldTables, newTables := d.GetChange("postgres_tables")
		if toAdd := expandStringList(newTables.(*schema.Set).Difference(oldTables.(*schema.Set)).List()); len(toAdd) > 0 {
			reqs.set.WithAddTables(postgresMirrorNames(toAdd))
		}
		if toRemove := expandStringList(oldTables.(*schema.Set).Difference(newTables.(*schema.Set)).List()); len(toRemove) > 0 {
			reqs.set.WithRemoveTables(postgresMirrorNames(toRemove))
		}
	}

	if d.HasChange("postgres_schemas") {
		oldSchemas, newSchemas := d.GetChange("postgres_schemas")
		if toAdd := expandStringList(newSchemas.(*schema.Set).Difference(oldSchemas.(*schema.Set)).List()); len(toAdd) > 0 {
			reqs.set.WithAddSchemas(postgresMirrorNames(toAdd))
		}
		if toRemove := expandStringList(oldSchemas.(*schema.Set).Difference(newSchemas.(*schema.Set)).List()); len(toRemove) > 0 {
			reqs.set.WithRemoveSchemas(postgresMirrorNames(toRemove))
		}
	}

	if d.HasChange("refresh_interval") {
		oldBlock, newBlock := d.GetChange("refresh_interval")
		newInterval, isSet := postgresMirrorRefreshInterval(newBlock)
		oldInterval, _ := postgresMirrorRefreshInterval(oldBlock)
		if isSet && newInterval != oldInterval {
			reqs.set.WithRefreshInterval(newInterval)
		}
	}

	if d.HasChange("warehouse") {
		if warehouse := d.Get("warehouse").(string); warehouse == "" {
			reqs.unset.WithUnsetWarehouse(true)
		} else {
			warehouseName, err := bareAccountObjectName(warehouse)
			if err != nil {
				return nil, err
			}
			reqs.set.WithWarehouse(warehouseName)
		}
	}

	return reqs, nil
}

// ALTER_MIRROR rejects a call that changes nothing, so a request is sent only when it differs from an empty one.
func postgresMirrorAlterSetInSdkExt(ctx context.Context, _ *schema.ResourceData, meta any, id sdk.AccountObjectIdentifier, reqs *postgresMirrorAlterRequestsExt) error {
	if reflect.DeepEqual(*reqs.set, *sdk.NewAlterPostgresMirrorArgumentsRequest(id.Name())) {
		return nil
	}
	if err := meta.(*provider.Context).Client.PostgresMirrors.Alter(ctx, sdk.NewAlterPostgresMirrorRequest(*reqs.set)); err != nil {
		return fmt.Errorf("error altering Postgres mirror %s: %w", id.Name(), err)
	}
	return nil
}

func postgresMirrorAlterUnsetInSdkExt(ctx context.Context, meta any, id sdk.AccountObjectIdentifier, reqs *postgresMirrorAlterRequestsExt) error {
	if reflect.DeepEqual(*reqs.unset, *sdk.NewAlterPostgresMirrorArgumentsRequest(id.Name())) {
		return nil
	}
	if err := meta.(*provider.Context).Client.PostgresMirrors.Alter(ctx, sdk.NewAlterPostgresMirrorRequest(*reqs.unset)); err != nil {
		return fmt.Errorf("error unsetting Postgres mirror %s: %w", id.Name(), err)
	}
	return nil
}

// =============================================================================
// Delete
// =============================================================================

func DeletePostgresMirror(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	request := sdk.NewDropPostgresMirrorRequest(*sdk.NewDropPostgresMirrorArgumentsRequest(d.Id()).WithDropTargetDatabase(d.Get("drop_target_database_on_destroy").(bool)))
	if err := client.PostgresMirrors.DropSafely(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error dropping Postgres mirror %s: %w", d.Id(), err))
	}

	d.SetId("")
	return nil
}

// =============================================================================
// Helpers
// =============================================================================

func postgresMirrorNames(names []string) []sdk.StringListItemWrapper {
	return collections.Map(names, func(name string) sdk.StringListItemWrapper { return sdk.StringListItemWrapper{Value: name} })
}

// DESCRIBE_MIRROR returns schema names through quote_ident, but CREATE_MIRROR and ALTER_MIRROR take them as literal
// names, so the quoting is removed before the names are stored and later sent back.
func unquotePostgresIdentifier(name string) string {
	if len(name) >= 2 && strings.HasPrefix(name, `"`) && strings.HasSuffix(name, `"`) {
		return strings.ReplaceAll(name[1:len(name)-1], `""`, `"`)
	}
	return name
}

// postgresMirrorRefreshInterval spells the configured block as the phrase the procedures take. An omitted block is not
// set: ALTER_MIRROR cannot reset the interval, so removing the block keeps the current one.
func postgresMirrorRefreshInterval(v any) (string, bool) {
	blocks := v.([]any)
	if len(blocks) == 0 || blocks[0] == nil {
		return "", false
	}
	block := blocks[0].(map[string]any)
	return sdk.PostgresMirrorRefreshInterval{
		Hours:   block["hours"].(int),
		Minutes: block["minutes"].(int),
		Seconds: block["seconds"].(int),
	}.String(), true
}

func setPostgresMirrorRefreshInterval(d *schema.ResourceData, value string) error {
	interval, err := sdk.ParsePostgresMirrorRefreshInterval(value)
	if err != nil {
		return err
	}
	return d.Set("refresh_interval", []any{map[string]any{
		"hours":   interval.Hours,
		"minutes": interval.Minutes,
		"seconds": interval.Seconds,
	}})
}

func bareAccountObjectName(name string) (string, error) {
	id, err := sdk.ParseAccountObjectIdentifier(name)
	if err != nil {
		return "", err
	}
	return id.Name(), nil
}

func quotedAccountObjectName(name string) (string, error) {
	id, err := sdk.ParseAccountObjectIdentifier(name)
	if err != nil {
		return "", err
	}
	return id.FullyQualifiedName(), nil
}
