package resources

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var stageCommonSchema = map[string]*schema.Schema{
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
			Schema: schemas.StageDescribeSchema,
		},
	},
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

var DeleteStage = PreviewFeatureDeleteContextWrapper(
	string(previewfeatures.InternalStageResource),
	TrackingDeleteWrapper(
		resources.InternalStage,
		ResourceDeleteContextFunc(
			sdk.ParseSchemaObjectIdentifier,
			func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] { return client.Stages.DropSafely },
		),
	),
)
