package resources

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
			ForceNew:         true,
			Description:      blocklistedCharactersFieldDescription("The database in which to create the stage."),
			DiffSuppressFunc: suppressIdentifierQuoting,
		},
		"schema": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
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

func handleStageRename(ctx context.Context, client *sdk.Client, d *schema.ResourceData, id sdk.SchemaObjectIdentifier) (sdk.SchemaObjectIdentifier, error) {
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		newId := sdk.NewSchemaObjectIdentifierInSchema(id.SchemaId(), newName)

		err := client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).WithRenameTo(newId))
		if err != nil {
			return sdk.SchemaObjectIdentifier{}, fmt.Errorf("error renaming stage %v to %v: %w", id.FullyQualifiedName(), newId.FullyQualifiedName(), err)
		}

		d.SetId(helpers.EncodeResourceIdentifier(newId))
		id = newId
	}
	return id, nil
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
		if notificationIntegration, ok := directoryConfig["notification_integration"]; ok && notificationIntegration.(string) != "" {
			directoryReq.WithNotificationIntegration(notificationIntegration.(string))
		}
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

func directoryTableToSet(d *schema.ResourceData, directoryTable sdk.StageDirectoryTable) []any {
	directory := map[string]any{
		"enable":       directoryTable.Enable,
		"auto_refresh": booleanStringFromBool(directoryTable.AutoRefresh),
	}
	if directoryTable.NotificationIntegration != nil {
		directory["notification_integration"] = *directoryTable.NotificationIntegration
	} else if d != nil {
		if current, ok := d.GetOk("directory.0.notification_integration"); ok {
			directory["notification_integration"] = current.(string)
		}
		if current, ok := d.GetOk("directory.0.refresh_on_create"); ok {
			directory["refresh_on_create"] = current.(string)
		}
	}
	return []any{directory}
}

func directoryNotificationIntegration(directoryTable sdk.StageDirectoryTable) string {
	if directoryTable.NotificationIntegration == nil {
		return ""
	}
	return *directoryTable.NotificationIntegration
}

func directoryTableToCompare(directoryTable sdk.StageDirectoryTable) []any {
	return []any{
		map[string]any{
			"enable":                   directoryTable.Enable,
			"auto_refresh":             directoryTable.AutoRefresh,
			"notification_integration": directoryNotificationIntegration(directoryTable),
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
	notificationIntegration, _ := directoryTable["notification_integration"].(string)
	return []any{
		map[string]any{
			"enable":                   directoryTable["enable"],
			"auto_refresh":             directoryTable["auto_refresh"],
			"notification_integration": notificationIntegration,
		},
	}
}

func directoryTableOutputMapping(d *schema.ResourceData, directoryTable sdk.StageDirectoryTable) outputMapping {
	return outputMapping{
		"directory_table", "directory", directoryTableToCompare(directoryTable), directoryTableToSet(d, directoryTable), normalizeDirectoryTableForExternalChangeCompare,
	}
}

func setImportedDirectoryNotificationIntegration(d *schema.ResourceData, details *sdk.StageDetails) error {
	if details == nil || details.DirectoryTable == nil || details.DirectoryTable.NotificationIntegration == nil || *details.DirectoryTable.NotificationIntegration == "" {
		return nil
	}
	directory, ok := d.Get("directory").([]any)
	if !ok || len(directory) != 1 {
		return nil
	}
	current, ok := directory[0].(map[string]any)
	if !ok {
		return nil
	}
	updated := make(map[string]any, len(current)+1)
	for key, value := range current {
		updated[key] = value
	}
	updated["notification_integration"] = *details.DirectoryTable.NotificationIntegration
	return d.Set("directory", []any{updated})
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
