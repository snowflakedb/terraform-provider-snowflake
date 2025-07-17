package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var stageSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Specifies the identifier for the stage; must be unique for the database and schema in which the stage is created.",
		ForceNew:    true,
	},
	"database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database in which to create the stage.",
		ForceNew:    true,
	},
	"schema": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The schema in which to create the stage.",
		ForceNew:    true,
	},
	"url": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the URL for the stage.",
	},
	"credentials": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the credentials for the stage.",
		Sensitive:   true,
	},
	"storage_integration": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the name of the storage integration used to delegate authentication responsibility for external cloud storage to a Snowflake identity and access management (IAM) entity.",
	},
	"file_format": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies the file format for the stage. Specifying the default Snowflake value (e.g. TYPE = CSV) will currently result in a permadiff (check [#2679](https://github.com/snowflakedb/terraform-provider-snowflake/issues/2679)). For now, omit the default values; it will be fixed in the upcoming provider versions. Examples of usage: <b>1. with hardcoding value:</b> `file_format=\"FORMAT_NAME = DB.SCHEMA.FORMATNAME\"` <b>2. from dynamic value:</b> `file_format = \"FORMAT_NAME = ${snowflake_file_format.myfileformat.fully_qualified_name}\"` <b>3. from expression:</b> `file_format = format(\"FORMAT_NAME =%s.%s.MYFILEFORMAT\", var.db_name, each.value.schema_name)`. Reference: [#265](https://github.com/snowflakedb/terraform-provider-snowflake/issues/265)",
		DiffSuppressFunc: suppressQuoting,
	},
	"copy_options": {
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "Specifies the copy options for the stage.",
		DiffSuppressFunc: suppressQuoting,
	},
	"encryption": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies the encryption settings for the stage.",
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the stage.",
	},
	"directory": {
		Type:        schema.TypeString,
		ForceNew:    true,
		Optional:    true,
		Description: "Specifies the directory settings for the stage.",
	},
	"aws_external_id": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		// Description based on https://docs.snowflake.com/en/user-guide/data-load-s3-config-aws-iam-role#step-3-create-an-external-stage
		Description: "A unique ID assigned to the specific stage. The ID has the following format: &lt;snowflakeAccount&gt;_SFCRole=&lt;snowflakeRoleId&gt;_&lt;randomId&gt;",
	},
	"snowflake_iam_user": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
		// Description based on https://docs.snowflake.com/en/user-guide/data-load-s3-config-aws-iam-role#step-3-create-an-external-stage
		Description: "An AWS IAM user created for your Snowflake account. This user is the same for every external S3 stage created in your account.",
	},
	"tag":                           tagReferenceSchema,
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

// TODO (SNOW-1019005): Remove snowflake package that is used in Create and Update operations
func Stage() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		helpers.DecodeSnowflakeIDErrLegacy[sdk.SchemaObjectIdentifier],
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] { return client.Stages.DropSafely },
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.StageResource), TrackingCreateWrapper(resources.Stage, CreateStage)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.StageResource), TrackingReadWrapper(resources.Stage, ReadStage)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.StageResource), TrackingUpdateWrapper(resources.Stage, UpdateStage)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.StageResource), TrackingDeleteWrapper(resources.Stage, deleteFunc)),

		Schema: stageSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	db := client.GetConn().DB
	name := d.Get("name").(string)
	database := d.Get("database").(string)
	schema := d.Get("schema").(string)

	builder := snowflake.NewStageBuilder(name, database, schema)

	if v, ok := d.GetOk("url"); ok {
		builder.WithURL(v.(string))
	}

	if v, ok := d.GetOk("credentials"); ok {
		builder.WithCredentials(v.(string))
	}

	if v, ok := d.GetOk("storage_integration"); ok {
		builder.WithStorageIntegration(v.(string))
	}

	if v, ok := d.GetOk("file_format"); ok {
		builder.WithFileFormat(v.(string))
	}

	if v, ok := d.GetOk("copy_options"); ok {
		builder.WithCopyOptions(v.(string))
	}

	if v, ok := d.GetOk("directory"); ok {
		builder.WithDirectory(v.(string))
	}

	if v, ok := d.GetOk("encryption"); ok {
		builder.WithEncryption(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		builder.WithComment(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		tags := getTags(v)
		builder.WithTags(tags.toSnowflakeTagValues())
	}

	q := builder.Create()

	if err := snowflake.Exec(db, q); err != nil {
		return diag.Errorf("error creating stage %v, err: %v", name, err)
	}

	d.SetId(helpers.EncodeSnowflakeID(database, schema, name))

	return ReadStage(ctx, d, meta)
}

func ReadStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.SchemaObjectIdentifier)

	stage, err := client.Stages.ShowByID(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query stage. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Stage id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	properties, err := client.Stages.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", stage.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("database", stage.DatabaseName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("schema", stage.SchemaName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url", strings.Trim(findStagePropertyValueByName(properties, "URL"), "[\"]")); err != nil {
		return diag.FromErr(err)
	}

	fileFormat := make([]string, 0)
	for _, property := range properties {
		if property.Parent == "STAGE_FILE_FORMAT" && property.Value != property.Default {
			fileFormat = append(fileFormat, fmt.Sprintf("%s = %s", property.Name, property.Value))
		}
	}
	if err := d.Set("file_format", strings.Join(fileFormat, " ")); err != nil {
		return diag.FromErr(err)
	}

	copyOptions := make([]string, 0)
	for _, property := range properties {
		if property.Parent == "STAGE_COPY_OPTIONS" && property.Value != property.Default {
			copyOptions = append(copyOptions, fmt.Sprintf("%s = %s", property.Name, property.Value))
		}
	}
	if err := d.Set("copy_options", strings.Join(copyOptions, " ")); err != nil {
		return diag.FromErr(err)
	}

	directory := make([]string, 0)
	for _, property := range properties {
		if property.Parent == "DIRECTORY" && property.Value != property.Default && property.Name != "LAST_REFRESHED_ON" {
			directory = append(directory, fmt.Sprintf("%s = %s", property.Name, property.Value))
		}
	}
	if err := d.Set("directory", strings.Join(directory, " ")); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("storage_integration", stage.StorageIntegration); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("comment", stage.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("aws_external_id", findStagePropertyValueByName(properties, "AWS_EXTERNAL_ID")); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("snowflake_iam_user", findStagePropertyValueByName(properties, "SNOWFLAKE_IAM_USER")); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateStage(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	id := helpers.DecodeSnowflakeIDLegacy(d.Id()).(sdk.SchemaObjectIdentifier)

	builder := snowflake.NewStageBuilder(id.Name(), id.DatabaseName(), id.SchemaName())

	client := meta.(*provider.Context).Client
	db := client.GetConn().DB

	if d.HasChange("credentials") {
		credentials := d.Get("credentials")
		q := builder.ChangeCredentials(credentials.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage credentials on %v", d.Id())
		}
	}

	if d.HasChange("storage_integration") && d.HasChange("url") {
		si := d.Get("storage_integration")
		url := d.Get("url")
		q := builder.ChangeStorageIntegrationAndUrl(si.(string), url.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage storage integration and url on %v", d.Id())
		}
	} else {
		if d.HasChange("storage_integration") {
			si := d.Get("storage_integration")
			q := builder.ChangeStorageIntegration(si.(string))
			if err := snowflake.Exec(db, q); err != nil {
				return diag.Errorf("error updating stage storage integration on %v", d.Id())
			}
		}

		if d.HasChange("url") {
			url := d.Get("url")
			q := builder.ChangeURL(url.(string))
			if err := snowflake.Exec(db, q); err != nil {
				return diag.Errorf("error updating stage url on %v", d.Id())
			}
		}
	}

	if d.HasChange("encryption") {
		encryption := d.Get("encryption")
		q := builder.ChangeEncryption(encryption.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage encryption on %v", d.Id())
		}
	}

	if d.HasChange("file_format") {
		fileFormat := d.Get("file_format")
		q := builder.ChangeFileFormat(fileFormat.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage file format on %v", d.Id())
		}
	}

	if d.HasChange("copy_options") {
		copyOptions := d.Get("copy_options")
		q := builder.ChangeCopyOptions(copyOptions.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage copy options on %v", d.Id())
		}
	}

	if d.HasChange("comment") {
		comment := d.Get("comment")
		q := builder.ChangeComment(comment.(string))
		if err := snowflake.Exec(db, q); err != nil {
			return diag.Errorf("error updating stage comment on %v", d.Id())
		}
	}

	if d.HasChange("tag") {
		unsetTags, setTags := GetTagsDiff(d, "tag")

		if len(unsetTags) > 0 {
			err := client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).WithUnsetTags(unsetTags))
			if err != nil {
				return diag.Errorf("error occurred when dropping tags on stage with id: %v, err = %s", d.Id(), err)
			}
		}

		if len(setTags) > 0 {
			err := client.Stages.Alter(ctx, sdk.NewAlterStageRequest(id).WithSetTags(setTags))
			if err != nil {
				return diag.Errorf("error occurred when setting tags on stage with id: %v, err = %s", d.Id(), err)
			}
		}
	}

	return ReadStage(ctx, d, meta)
}

func findStagePropertyValueByName(properties []sdk.StageProperty, name string) string {
	for _, property := range properties {
		if property.Name == name {
			return property.Value
		}
	}
	return ""
}
