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

var cortexAgentSchema = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("Specifies the identifier for the Cortex agent."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"schema": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The schema in which to create the Cortex agent."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"database": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedCharactersFieldDescription("The database in which to create the Cortex agent."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"specification": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      joinWithSpace("Specifies a YAML object containing the settings for the Cortex agent.", doubleDollarQuotesDescription()),
		DiffSuppressFunc: NormalizeAndCompare(sdk.NormalizeCortexAgentSpecification),
		ValidateDiagFunc: validation.AllDiag(validation.ToDiagFunc(validation.StringIsNotEmpty), forbidDoubleDollarQuotes),
	},
	"comment": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Specifies a comment for the Cortex agent.",
	},
	"profile": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Specifies agent profile information, such as display_name, avatar, and color.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Specifies a display name for the Cortex agent.",
					AtLeastOneOf: []string{"profile.0.display_name", "profile.0.avatar", "profile.0.color"},
				},
				"avatar": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Specifies an avatar image file name or identifier.",
					AtLeastOneOf: []string{"profile.0.display_name", "profile.0.avatar", "profile.0.color"},
				},
				"color": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Specifies a color theme for the Cortex agent.",
					AtLeastOneOf: []string{"profile.0.display_name", "profile.0.avatar", "profile.0.color"},
				},
			},
		},
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW AGENTS` for this Cortex agent.",
		Elem: &schema.Resource{
			Schema: schemas.ShowCortexAgentSchema,
		},
	},
	DescribeOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `DESCRIBE AGENT` for this Cortex agent.",
		Elem: &schema.Resource{
			Schema: schemas.DescribeCortexAgentDetailsSchema,
		},
	},
	FullyQualifiedNameAttributeName: schemas.FullyQualifiedNameSchema,
}

func CortexAgent() *schema.Resource {
	deleteFunc := ResourceDeleteContextFunc(
		sdk.ParseSchemaObjectIdentifier,
		func(client *sdk.Client) DropSafelyFunc[sdk.SchemaObjectIdentifier] {
			return client.CortexAgents.DropSafely
		},
	)

	return &schema.Resource{
		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.CortexAgentResource), TrackingCreateWrapper(resources.CortexAgent, CreateCortexAgent)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.CortexAgentResource), TrackingReadWrapper(resources.CortexAgent, ReadCortexAgent)),
		UpdateContext: PreviewFeatureUpdateContextWrapper(string(previewfeatures.CortexAgentResource), TrackingUpdateWrapper(resources.CortexAgent, UpdateCortexAgent)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.CortexAgentResource), TrackingDeleteWrapper(resources.CortexAgent, deleteFunc)),
		Description:   "Resource used to manage Cortex agent objects. For more information, check [Cortex agent documentation](https://docs.snowflake.com/en/sql-reference/sql/create-agent).",

		Schema: cortexAgentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.CortexAgent, ImportName[sdk.SchemaObjectIdentifier]),
		},
		Timeouts: defaultTimeouts,
		CustomizeDiff: TrackingCustomDiffWrapper(resources.CortexAgent, customdiff.All(
			ComputedIfAnyAttributeChanged(cortexAgentSchema, ShowOutputAttributeName, "comment", "profile"),
			ComputedIfAnyAttributeChanged(cortexAgentSchema, DescribeOutputAttributeName, "specification", "comment", "profile"),
		)),
	}
}

func CreateCortexAgent(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	databaseName := d.Get("database").(string)
	schemaName := d.Get("schema").(string)
	name := d.Get("name").(string)
	id := sdk.NewSchemaObjectIdentifier(databaseName, schemaName, name)
	spec := d.Get("specification").(string)

	request := sdk.NewCreateCortexAgentRequest(id, spec)
	errs := errors.Join(
		stringAttributeCreateBuilder(d, "comment", request.WithComment),
		attributeMappedValueCreateBuilder(d, "profile", request.WithProfile, cortexAgentProfileToJsonString),
	)
	if errs != nil {
		return diag.FromErr(errs)
	}

	if err := client.CortexAgents.Create(ctx, request); err != nil {
		return diag.FromErr(fmt.Errorf("error creating Cortex agent, err = %w", err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(id))
	return ReadCortexAgent(ctx, d, meta)
}

func ReadCortexAgent(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	agent, err := client.CortexAgents.ShowByIDSafely(ctx, id)
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotFound) {
			d.SetId("")
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "Failed to query Cortex agent. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Cortex agent id: %s, Err: %s", id.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	details, err := client.CortexAgents.Describe(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	profileList := make([]any, 0)
	if !reflect.DeepEqual(details.Profile, sdk.CortexAgentProfile{}) {
		block := map[string]any{}
		if details.Profile.DisplayName != nil {
			block["display_name"] = *details.Profile.DisplayName
		}
		if details.Profile.Avatar != nil {
			block["avatar"] = *details.Profile.Avatar
		}
		if details.Profile.Color != nil {
			block["color"] = *details.Profile.Color
		}
		profileList = append(profileList, block)
	}

	errs := errors.Join(
		d.Set("specification", details.AgentSpec),
		d.Set("comment", details.Comment),
		d.Set("profile", profileList),
		d.Set(ShowOutputAttributeName, []map[string]any{schemas.CortexAgentToSchema(agent)}),
		d.Set(DescribeOutputAttributeName, []map[string]any{schemas.CortexAgentDetailsToSchema(details)}),
		d.Set(FullyQualifiedNameAttributeName, id.FullyQualifiedName()),
	)
	return diag.FromErr(errs)
}

func UpdateCortexAgent(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var modifyLiveVersionSetRequest *sdk.CortexAgentModifyLiveVersionSetRequest
	setRequest := sdk.NewCortexAgentSetRequest()

	if d.HasChange("specification") {
		modifyLiveVersionSetRequest = sdk.NewCortexAgentModifyLiveVersionSetRequest(d.Get("specification").(string))
	}

	if err := errors.Join(
		// TODO [SNOW-3530838]: UNSET not implemented
		stringAttributeUpdateSetOnly(d, "comment", &setRequest.Comment),
		attributeMappedValueUpdateSetOnlyFallback(d, "profile", &setRequest.Profile, cortexAgentProfileToJsonString, "{}"),
	); err != nil {
		return diag.FromErr(err)
	}

	if modifyLiveVersionSetRequest != nil {
		if err := client.CortexAgents.Alter(ctx, sdk.NewAlterCortexAgentRequest(id).WithModifyLiveVersionSet(*modifyLiveVersionSetRequest)); err != nil {
			return diag.FromErr(err)
		}
	}

	if !reflect.DeepEqual(*setRequest, *sdk.NewCortexAgentSetRequest()) {
		if err := client.CortexAgents.Alter(ctx, sdk.NewAlterCortexAgentRequest(id).WithSet(*setRequest)); err != nil {
			return diag.FromErr(err)
		}
	}

	return ReadCortexAgent(ctx, d, meta)
}

func cortexAgentProfileToJsonString(v any) (string, error) {
	list := v.([]any)
	if len(list) == 0 {
		return "", fmt.Errorf("Cortex agent profile block is empty")
	}
	block := list[0].(map[string]any)

	profile := sdk.CortexAgentProfile{}
	displayName := block["display_name"].(string)
	if displayName != "" {
		profile.DisplayName = &displayName
	}
	avatar := block["avatar"].(string)
	if avatar != "" {
		profile.Avatar = &avatar
	}
	color := block["color"].(string)
	if color != "" {
		profile.Color = &color
	}

	profileAsJson, err := sdk.MarshalCortexAgentProfile(profile)
	if err != nil {
		return "", err
	}
	return profileAsJson, nil
}
