package resources

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var streamlitSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("String that specifies the identifier (i.e. name) for the streamlit; must be unique in your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the streamlit"),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the streamlit."),
		ForceNew:         true,
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"stage": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      relatedResourceDescription("The stage in which streamlit files are located.", resources.Stage),
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInDescribe("root_location")),
	},
	"directory_location": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies the full path to the named stage containing the Streamlit Python files, media files, and the environment.yml file.",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("root_location"),
	},
	"main_file": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Specifies the filename of the Streamlit Python application. This filename is relative to the value of `directory_location`",
		DiffSuppressFunc: IgnoreChangeToCurrentSnowflakeValueInDescribe("main_file"),
	},
	"query_warehouse": {
		Type:             schema.TypeString,
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		Optional:         true,
		Description:      relatedResourceDescription("Specifies the warehouse where SQL queries issued by the Streamlit application are run. Due to Snowflake limitations warehouse identifier can consist of only upper-cased letters.", resources.Warehouse),
		DiffSuppressFunc: SuppressIfAny(suppressIdentifierQuoting, IgnoreChangeToCurrentSnowflakeValueInShow("query_warehouse")),
	},
	"external_access_integrations": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		},
		Optional:         true,
		Description:      "External access integrations connected to the Streamlit.",
		DiffSuppressFunc: SuppressIfAny(NormalizeAndCompareIdentifiersInSet("external_access_integrations"), IgnoreChangeToCurrentSnowflakeValueInDescribe("external_access_integrations")),
	},
	"title": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a title for the Streamlit app to display in Snowsight.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the streamlit.",
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW STREAMLIT` for the given streamlit.",
		Elem: &schema.Resource{
			Schema: schemas.ShowStreamlitSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE STREAMLIT` for the given streamlit.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeStreamlitSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func Streamlit() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Streamlits.DropSafely
		},
	)

	return &schema.Resource{
		SchemaVersion: 1,

		CreateContext: TrackingCreateWrapper(resources.Streamlit, CreateContextStreamlit),
		ReadContext:   TrackingReadWrapper(resources.Streamlit, ReadContextStreamlit),
		UpdateContext: TrackingUpdateWrapper(resources.Streamlit, UpdateContextStreamlit),
		DeleteContext: TrackingDeleteWrapper(resources.Streamlit, deleteFunc),
		Description:   "Resource used to manage streamlits objects. For more information, check [streamlit documentation](https://docs.snowflake.com/en/sql-reference/commands-streamlit).",

		Schema: streamlitSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Streamlit, ImportStreamlit),
		},

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Streamlit, customdiff.All(
			ComputedIfAnyAttributeChanged(streamlitSchema, ShowOutputAttributeName, "name", "title", "comment", "query_warehouse"),
			ComputedIfAnyAttributeChanged(streamlitSchema, FullyQualifiedNameAttributeName, "name"),
			ComputedIfAnyAttributeChanged(streamlitSchema, DescribeOutputAttributeName, "title", "comment", "root_location", "main_file", "query_warehouse", "external_access_integrations"),
		)),

		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				// setting type to cty.EmptyObject is a bit hacky here but following https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/state-upgrade#sdkv2-1 would require lots of repetitive code; this should work with cty.EmptyObject
				Type:    cty.EmptyObject,
				Upgrade: migratePipeSeparatedObjectIdentifierResourceIdToFullyQualifiedName,
			},
		},
		Timeouts: defaultTimeouts,
	}
}

func ImportStreamlit(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	if err = d.Set("name", id.Name()); err != nil {
		return nil, err
	}
	if err = d.Set("database", id.DatabaseName()); err != nil {
		return nil, err
	}
	if err = d.Set("schema", id.SchemaName()); err != nil {
		return nil, err
	}

	streamlit, err := client.Streamlits.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	streamlitDetails, err := client.Streamlits.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	stageId, location, err := helpers.ParseRootLocation(streamlitDetails.RootLocation)
	if err != nil {
		return nil, err
	}
	if err := d.Set("stage", stageId.FullyQualifiedName()); err != nil {
		return nil, err
	}
	if err := d.Set("directory_location", location); err != nil {
		return nil, err
	}
	if err = d.Set("main_file", streamlitDetails.MainFile); err != nil {
		return nil, err
	}
	if err = d.Set("query_warehouse", streamlit.QueryWarehouse); err != nil {
		return nil, err
	}
	if err = d.Set("external_access_integrations", streamlitDetails.ExternalAccessIntegrations); err != nil {
		return nil, err
	}
	if err = d.Set("title", streamlit.Title); err != nil {
		return nil, err
	}
	if err = d.Set("comment", streamlit.Comment); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func CreateContextStreamlit(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	stageId, err := sdk.ParseSchemaObjectIdentifier(d.Get("stage").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	rootLocation := fmt.Sprintf("@%s", stageId.FullyQualifiedName())
	if v, ok := d.GetOk("directory_location"); ok {
		rootLocation = path.Join(rootLocation, v.(string))
	}

	req := sdk.NewCreateStreamlitRequest(id, rootLocation, d.Get("main_file").(string))

	if v, ok := d.GetOk("query_warehouse"); ok {
		warehouseId, err := sdk.ParseAccountObjectIdentifier(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		req.WithQueryWarehouse(warehouseId)
	}

	if v, ok := d.GetOk("title"); ok {
		req.WithTitle(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		req.WithComment(v.(string))
	}

	if v, ok := d.GetOk("external_access_integrations"); ok {
		raw := expandStringList(v.(*schema.Set).List())
		integrations := make([]sdk.AccountObjectIdentifier, len(raw))
		for i, v := range raw {
			integrations[i] = sdk.NewAccountObjectIdentifier(v)
		}
		req.WithExternalAccessIntegrations(sdk.ExternalAccessIntegrationsRequest{
			ExternalAccessIntegrations: integrations,
		})
	}

	if err := client.Streamlits.Create(ctx, req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))

	return ReadContextStreamlit(ctx, d, meta)
}

func ReadContextStreamlit(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	streamlit, err := client.Streamlits.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query streamlit. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Streamlit id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}
	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	streamlitDetails, err := client.Streamlits.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	stageId, location, err := helpers.ParseRootLocation(streamlitDetails.RootLocation)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("stage", stageId.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("directory_location", location); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("main_file", streamlitDetails.MainFile); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("query_warehouse", streamlit.QueryWarehouse); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("external_access_integrations", streamlitDetails.ExternalAccessIntegrations); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("title", streamlit.Title); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("comment", streamlit.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set(ShowOutputAttributeName, []map[string]any{schemas.StreamlitToSchema(streamlit)}); err != nil {
		return diag.FromErr(err)
	}
	schemaDetails, err := schemas.StreamlitPropertiesToSchema(*streamlitDetails)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set(DescribeOutputAttributeName, []map[string]any{schemaDetails}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateContextStreamlit(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewStreamlitSetRequest(), sdk.NewStreamlitUnsetRequest()

	if d.HasChange("name") {
		databaseName := d.Get("database").(string)
		schemaName := d.Get("schema").(string)
		name := d.Get("name").(string)
		newId := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

		err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(newId.FullyQualifiedName())
		id = newId
	}

	if d.HasChange("stage") || d.HasChange("directory_location") {
		stageId := sdk.NewSchemaObjectIdentifierFromFullyQualifiedName(d.Get("stage").(string))
		rootLocation := fmt.Sprintf("@%s", stageId.FullyQualifiedName())
		if v, ok := d.GetOk("directory_location"); ok {
			rootLocation = path.Join(rootLocation, v.(string))
		}
		set.WithRootLocation(rootLocation)
	}

	if d.HasChange("main_file") {
		// required field
		set.WithMainFile(d.Get("main_file").(string))
	}

	if d.HasChange("title") {
		if v, ok := d.GetOk("title"); ok {
			set.WithTitle(v.(string))
		} else {
			unset.WithTitle(true)
		}
	}

	if d.HasChange("query_warehouse") {
		if v, ok := d.GetOk("query_warehouse"); ok {
			warehouseId, err := sdk.ParseAccountObjectIdentifier(v.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			set.WithQueryWarehouse(warehouseId)
		} else {
			unset.WithQueryWarehouse(true)
		}
	}

	if d.HasChange("title") {
		if v, ok := d.GetOk("title"); ok {
			set.WithTitle(v.(string))
		} else {
			unset.WithTitle(true)
		}
	}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			set.WithComment(v.(string))
		} else {
			unset.WithComment(true)
		}
	}

	if d.HasChange("external_access_integrations") {
		raw := expandStringList(d.Get("external_access_integrations").(*schema.Set).List())
		integrations := make([]sdk.AccountObjectIdentifier, len(raw))
		for i, v := range raw {
			integrationId, err := sdk.ParseAccountObjectIdentifier(v)
			if err != nil {
				return diag.FromErr(err)
			}
			integrations[i] = integrationId
		}
		set.WithExternalAccessIntegrations(sdk.ExternalAccessIntegrationsRequest{
			ExternalAccessIntegrations: integrations,
		})
	}

	if (*set != sdk.StreamlitSetRequest{}) {
		if err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if (*unset != sdk.StreamlitUnsetRequest{}) {
		if err := client.Streamlits.Alter(ctx, sdk.NewAlterStreamlitRequest(id).WithUnset(*unset)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadContextStreamlit(ctx, d, meta)
}
