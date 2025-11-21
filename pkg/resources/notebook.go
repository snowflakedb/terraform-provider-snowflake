package resources

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var notebookSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the notebook; must be unique for the schema in which the notebook is created."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the notebook."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the notebook."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"from": {
		Type:         schema.TypeList,
		Optional:     true,
		ForceNew:     true,
		Description:  "Specifies the location in a stage of an .ipynb file from which the notebook should be created. MAIN_FILE parameter a user-specified identifier for the notebook file name must also be set alongside it.",
		RequiredWith: []string{"main_file"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"stage": {
					Type:             schema.TypeString,
					Required:         true,
					ForceNew:         true,
					Description:      "Identifier of the stage where the .ipynb file is located.",
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
					DiffSuppressFunc: suppressIdentifierQuoting,
				},
				"path": {
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
					Description: "Location of the .ipynb file in the stage.",
				},
			},
		},
	},
	"main_file": {
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
		Description: "Specifies a user-specified identifier for the notebook file name.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the notebook.",
	},
	"query_warehouse": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies the warehouse where SQL queries in the notebook are run.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"idle_auto_shutdown_time_seconds": {
		Type:             schema.TypeInt,
		Optional:         true,
		Description:      "Specifies the number of seconds of idle time before the notebook is shut down automatically.",
		ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
	},
	"warehouse": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies the warehouse that runs the notebook kernel and python code.",
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"secrets": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specifies secret variables for the notebook.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"secret_variable_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The name of the secret variable.",
				},
				"secret_id": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Fully qualified name of the allowed [secret](https://docs.snowflake.com/en/sql-reference/sql/create-secret).",
					ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
					DiffSuppressFunc: suppressIdentifierQuoting,
				},
			},
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW NOTEBOOKS` for the given notebook",
		Elem: &schema.Resource{
			Schema: schemas.ShowNotebookSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE NOTEBOOK` for the given notebook.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeNotebookSchema,
		},
	},
}

func Notebook() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.Notebooks.DropSafely
		},
	)
	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.NotebookResource), TrackingCreateWrapper(resources.Notebook, CreateNotebook)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.NotebookResource), TrackingReadWrapper(resources.Notebook, GetReadNotebookFunc(true))),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.NotebookResource), TrackingUpdateWrapper(resources.Notebook, UpdateNotebook)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.NotebookResource), TrackingDeleteWrapper(resources.Notebook, deleteFunc)),
		Description:   "Resource used to manage notebooks. For more information, check [notebooks documentation](https://docs.snowflake.com/en/sql-reference/sql/create-notebook).",

		CustomizeDiff: TrackingCustomDiffWrapper(resources.Notebook, customdiff.All(
			ComputedIfAnyAttributeChanged(notebookSchema, ShowOutputAttributeName, "name", "comment", "query_warehouse", "code_warehouse"),
			ComputedIfAnyAttributeChanged(notebookSchema, DescribeOutputAttributeName, "name", "comment", "query_warehouse", "code_warehouse", "idle_auto_shutdown_time_seconds", "main_file"),
			ComputedIfAnyAttributeChanged(notebookSchema, FullyQualifiedNameAttributeName, "name"),
		)),

		Schema: notebookSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.Notebook, ImportNotebook),
		},

		Timeouts: defaultTimeouts,
	}
}

func ImportNotebook(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}

	notebook, err := client.Notebooks.ShowByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := ImportName[sdk.SchemaObjectIdentifier](context.Background(), d, nil); err != nil {
		return nil, err
	}

	errs := errors.Join(
		d.Set("warehouse", notebook.CodeWarehouse.FullyQualifiedName()),
	)
	if errs != nil {
		return nil, errs
	}
	return []*schema.ResourceData{d}, nil
}

func CreateNotebook(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	name := d.Get("name").(string)
	schemaName := d.Get("schema").(string)
	databaseName := d.Get("database").(string)

	notebookName := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

	request := sdk.NewCreateNotebookRequest(notebookName)

	errs := errors.Join(
		stringAttributeCreateBuilder(d, "main_file", request.WithMainFile),
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		accountObjectIdentifierAttributeCreate(d, "query_warehouse", &request.QueryWarehouse),
		intAttributeCreateBuilder(d, "idle_auto_shutdown_time_seconds", request.WithIdleAutoShutdownTimeSeconds),
		accountObjectIdentifierAttributeCreate(d, "warehouse", &request.Warehouse),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if fromStage, ok := d.GetOk("from"); ok && len(fromStage.([]any)) > 0 {
		fromStageMap := fromStage.([]any)[0].(map[string]any)

		stage, err := sdk.ParseSchemaObjectIdentifier(fromStageMap["stage"].(string))
		if err != nil {
			return diag.FromErr(err)
		}

		var path string
		if l, ok := fromStageMap["path"]; ok {
			path = l.(string)
		}

		request.WithFrom(sdk.NewStageLocation(stage, path))
	}

	if secrets, ok := d.GetOk("secrets"); ok && len(secrets.([]any)) > 0 {
		secretsList := make([]sdk.SecretReference, 0)
		for _, secret := range secrets.([]any) {
			secretMap := secret.(map[string]any)

			secretVariableName := secretMap["secret_variable_name"].(string)
			secretId, err := sdk.ParseSchemaObjectIdentifier(secretMap["secret_id"].(string))
			if err != nil {
				return diag.FromErr(err)
			}

			secretsList = append(secretsList, sdk.SecretReference{
				VariableName: secretVariableName,
				Name:         secretId,
			})
		}
		request.WithSecrets(sdk.SecretsListRequest{SecretsList: secretsList})
	}

	if err := client.Notebooks.Create(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helpers.EncodeResourceIdentifier(notebookName))
	return GetReadNotebookFunc(false)(ctx, d, meta)
}

func GetReadNotebookFunc(withExternalChangesMarking bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
		client := meta.(*provider.Context).Client
		id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		notebook, err := client.Notebooks.ShowByIDSafely(ctx, id)
		if err != nil {
			if errors.Is(err, sdk.ErrObjectNotFound) {
				d.SetId("")
				return diag.Diagnostics{
					diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "Failed to query notebook. Marking the resource as removed.",
						Detail:   fmt.Sprintf("Notebook id: %s, Err: %s", id.FullyQualifiedName(), err),
					},
				}
			}
			return diag.FromErr(err)
		}

		notebookDetails, err := client.Notebooks.Describe(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		if withExternalChangesMarking {
			warehouse := notebook.CodeWarehouse.Name()
			if err = handleExternalChangesToObjectInShow(d,
				outputMapping{"code_warehouse", "warehouse", warehouse, warehouse, nil},
			); err != nil {
				return diag.FromErr(err)
			}

			secrets := notebookDetails.ExternalAccessSecrets
			if err = handleExternalChangesToObjectInFlatDescribe(d,
				outputMapping{"external_access_secrets", "secrets", secrets, secrets, nil},
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if err = setStateToValuesFromConfig(d, notebookSchema, []string{
			"warehouse",
		}); err != nil {
			return diag.FromErr(err)
		}

		errs := errors.Join(
			d.Set(ShowOutputAttributeName, []map[string]any{schemas.NotebookToSchema(notebook)}),
			d.Set(DescribeOutputAttributeName, []map[string]any{schemas.NotebookDetailsToSchema(notebookDetails)}),
			d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
			d.Set("comment", notebook.Comment),
		)
		if notebook.QueryWarehouse != nil {
			errs = errors.Join(errs, d.Set("query_warehouse", notebook.QueryWarehouse.FullyQualifiedName()))
		} else {
			if _, ok := d.GetOk("query_warehouse"); ok {
				errs = errors.Join(errs, d.Set("query_warehouse", ""))
			}
		}

		if errs != nil {
			return diag.FromErr(errs)
		}
		return nil
	}
}

func UpdateNotebook(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	set, unset := sdk.NewNotebookSetRequest(), sdk.NewNotebookUnsetRequest()

	if d.HasChange("name") {
		databaseName := d.Get("database").(string)
		schemaName := d.Get("schema").(string)
		name := d.Get("name").(string)
		newId := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)

		err := client.Notebooks.Alter(ctx, sdk.NewAlterNotebookRequest(id).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(newId.FullyQualifiedName())
		id = newId
	}

	errs := errors.Join(
		stringAttributeUpdate(d, "comment", &set.Comment, &unset.Comment),
		accountObjectIdentifierAttributeUpdate(d, "query_warehouse", &set.QueryWarehouse, &unset.QueryWarehouse),
		accountObjectIdentifierAttributeUpdate(d, "warehouse", &set.Warehouse, &unset.Warehouse),
		intAttributeUnsetFallbackUpdateWithZeroDefault(d, "idle_auto_shutdown_time_seconds", &set.IdleAutoShutdownTimeSeconds, 1800),
		func() error {
			if d.HasChange("secrets") {
				return setSecretsInBuilder(d, func(references []sdk.SecretReference) error {
					set.Secrets = &sdk.SecretsListRequest{SecretsList: references}
					return nil
				})
			}
			return nil
		}(),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if !reflect.DeepEqual(set, &sdk.NotebookSetRequest{}) {
		if err := client.Notebooks.Alter(ctx, sdk.NewAlterNotebookRequest(id).WithSet(*set)); err != nil {
			return diag.FromErr(err)
		}
	}

	if !reflect.DeepEqual(unset, &sdk.NotebookUnsetRequest{}) {
		// Special case with unsetting warehouse and query_warehouse at the same time.
		if (unset.Warehouse != nil && *unset.Warehouse == true) && (unset.QueryWarehouse != nil && *unset.QueryWarehouse == true) {
			unset.Warehouse = sdk.Bool(false)
			if err := client.Notebooks.Alter(ctx, sdk.NewAlterNotebookRequest(id).WithUnset(*unset)); err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
			if err := client.Notebooks.Alter(ctx, sdk.NewAlterNotebookRequest(id).WithUnset(sdk.NotebookUnsetRequest{Warehouse: sdk.Bool(true)})); err != nil {
				d.Partial(true)
				return diag.FromErr(err)
			}
		} else {
			if err := client.Notebooks.Alter(ctx, sdk.NewAlterNotebookRequest(id).WithUnset(*unset)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return GetReadNotebookFunc(false)(ctx, d, meta)
}
