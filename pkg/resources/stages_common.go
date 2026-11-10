package resources

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func stageCommonSchema(
	describeSchema map[string]*schema.Schema,
) map[string]*schema.Schema {
	stageCommonSchema := map[string]*schema.Schema{
		"name": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the stage; must be unique for the database and schema in which the stage is created."),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"database": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      blocklistedCharactersFieldDescription("The database in which to create the stage."),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"schema": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      blocklistedCharactersFieldDescription("The schema in which to create the stage."),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"stage_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Specifies a type for the stage. This field is used for checking external changes and recreating the resources if needed.",
		},
		"comment": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specifies a comment for the stage.",
		},
		FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
		ShowOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `SHOW STAGES` for the given stage.",
			Elem: &schema.Resource{
				Schema: schemas.ShowStageSchema,
			},
		},
		DescribeOutputAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Outputs the result of `DESCRIBE STAGE` for the given stage.",
			Elem: &schema.Resource{
				Schema: describeSchema,
			},
		},
	}
	return collections.MergeMaps(stageCommonSchema, stageFileFormatSchema)
}

func handleStageRename(ctx context.Context, meta any, d *schema.ResourceData, id *sdk.SchemaObjectIdentifier) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	if diags := handleThreeLevelHierarchyRenameIfEnabled(
		ctx, d, meta, id,
		func(currentId, targetId sdk.SchemaObjectIdentifier) error {
			return client.Stages.Alter(ctx, sdk.NewAlterStageRequest(currentId).WithRenameTo(targetId))
		},
		client.Stages.ShowByID,
		"stage",
	); diags != nil {
		return diags
	}

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), newName)

		err := client.Stages.Alter(ctx, sdk.NewAlterStageRequest(*id).WithRenameTo(newId))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error renaming stage %v to %v: %w", id.FullyQualifiedName(), newId.FullyQualifiedName(), err))
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		*id = newId
	}
	return nil
}

func handleStageDirectoryTable(ctx context.Context, client *sdk.Client, d *schema.ResourceData, id sdk.SchemaObjectIdentifier) error {
	setDirectoryTable := sdk.NewAlterDirectoryTableStageRequest(id)
	parseDirectoryTable := func(value any) (sdk.DirectoryTableSetRequest, error) {
		directoryList := value.([]any)
		if len(directoryList) == 0 {
			return sdk.DirectoryTableSetRequest{}, nil
		}
		directoryConfig := directoryList[0].(map[string]any)
		directoryReq := sdk.NewDirectoryTableSetRequest(directoryConfig["enable"].(bool))
		return *directoryReq, nil
	}
	err := attributeMappedValueUpdateSetOnly(d, "directory", &setDirectoryTable.SetDirectory, parseDirectoryTable)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(setDirectoryTable, sdk.NewAlterDirectoryTableStageRequest(id)) {
		if err := client.Stages.AlterDirectoryTable(ctx, setDirectoryTable); err != nil {
			return fmt.Errorf("error updating stage: %w", err)
		}
	}
	return nil
}

func DeleteStage(resource resources.Resource) schema.DeleteContextFunc {
	return TrackingDeleteWrapper(
		resource,
		ResourceDeleteContextFunc(
			sdk.ParseSchemaObjectIdentifier,
			func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] { return client.Stages.DropSafely },
		),
	)
}

func directoryTableToSet(directoryTable sdk.StageDirectoryTable) []any {
	return []any{
		map[string]any{
			"enable":       directoryTable.Enable,
			"auto_refresh": booleanStringFromBool(directoryTable.AutoRefresh),
		},
	}
}

func directoryTableToCompare(directoryTable sdk.StageDirectoryTable) []any {
	return []any{
		map[string]any{
			"enable":       directoryTable.Enable,
			"auto_refresh": directoryTable.AutoRefresh,
		},
	}
}

func normalizeDirectoryTableForExternalChangeCompare(value any) any {
	items, ok := value.([]any)
	if !ok || len(items) != 1 {
		return value
	}
	directoryTable, ok := items[0].(map[string]any)
	if !ok {
		return value
	}
	return []any{
		map[string]any{
			"enable":       directoryTable["enable"],
			"auto_refresh": directoryTable["auto_refresh"],
		},
	}
}

func directoryTableOutputMapping(directoryTable sdk.StageDirectoryTable) outputMapping {
	return outputMapping{
		"directory_table", "directory", directoryTableToCompare(directoryTable), directoryTableToSet(directoryTable), normalizeDirectoryTableForExternalChangeCompare,
	}
}

func importStageCommonFields(d *schema.ResourceData, details *sdk.StageDetails, setDefaults bool) error {
	if details.DirectoryTable != nil {
		autoRefreshValue := booleanStringFromBool(details.DirectoryTable.AutoRefresh)
		if setDefaults {
			autoRefreshValue = BooleanDefault
		}
		if err := d.Set("directory", []map[string]any{
			{
				"enable":       details.DirectoryTable.Enable,
				"auto_refresh": autoRefreshValue,
			},
		}); err != nil {
			return err
		}
	}

	if fileFormat := stageFileFormatToSchema(details, setDefaults); fileFormat != nil {
		if err := d.Set("file_format", fileFormat); err != nil {
			return err
		}
	}

	return nil
}
