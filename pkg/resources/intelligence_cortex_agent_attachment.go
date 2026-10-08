package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var intelligenceCortexAgentAttachmentSchema = map[string]*schema.Schema{
	"snowflake_intelligence_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedPipesFieldDescription("Fully qualified name of the Snowflake Intelligence object to attach the agent to. The account default object is SNOWFLAKE_INTELLIGENCE_OBJECT_DEFAULT."),
		ValidateDiagFunc: IsValidIdentifier[sdk.AccountObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
	"cortex_agent_name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      blocklistedPipesFieldDescription("Fully qualified name of the Cortex agent to attach."),
		ValidateDiagFunc: IsValidIdentifier[sdk.SchemaObjectIdentifier](),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

func IntelligenceCortexAgentAttachment() *schema.Resource {
	return &schema.Resource{
		Description: "Attaches a Cortex agent to a Snowflake Intelligence object so the agent is available in Snowflake CoWork.",

		CreateContext: PreviewFeatureCreateContextWrapper(string(previewfeatures.IntelligenceCortexAgentAttachmentResource), TrackingCreateWrapper(resources.IntelligenceCortexAgentAttachment, CreateIntelligenceCortexAgentAttachment)),
		ReadContext:   PreviewFeatureReadContextWrapper(string(previewfeatures.IntelligenceCortexAgentAttachmentResource), TrackingReadWrapper(resources.IntelligenceCortexAgentAttachment, ReadIntelligenceCortexAgentAttachment)),
		DeleteContext: PreviewFeatureDeleteContextWrapper(string(previewfeatures.IntelligenceCortexAgentAttachmentResource), TrackingDeleteWrapper(resources.IntelligenceCortexAgentAttachment, DeleteIntelligenceCortexAgentAttachment)),

		Schema: intelligenceCortexAgentAttachmentSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateIntelligenceCortexAgentAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	intelligenceName, err := sdk.ParseAccountObjectIdentifier(d.Get("snowflake_intelligence_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	agentName, err := sdk.ParseSchemaObjectIdentifier(d.Get("cortex_agent_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.SnowflakeIntelligences.Alter(ctx, sdk.NewAlterSnowflakeIntelligenceRequest(intelligenceName).WithAddAgent(agentName)); err != nil {
		return diag.FromErr(fmt.Errorf("error while attaching agent %s to Snowflake Intelligence %s, err = %w", agentName.FullyQualifiedName(), intelligenceName.FullyQualifiedName(), err))
	}

	d.SetId(helpers.EncodeResourceIdentifier(intelligenceName.FullyQualifiedName(), agentName.FullyQualifiedName()))

	return ReadIntelligenceCortexAgentAttachment(ctx, d, meta)
}

func ReadIntelligenceCortexAgentAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	intelligenceName, agentName, err := parseIntelligenceCortexAgentAttachmentId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	agents, err := client.SnowflakeIntelligences.ShowAgents(ctx, sdk.NewShowAgentsSnowflakeIntelligenceRequest(intelligenceName))
	if err != nil {
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Failed to query agents in the Snowflake Intelligence object. Marking the resource as removed.",
					Detail:   fmt.Sprintf("Snowflake Intelligence id: %s, Err: %s", intelligenceName.FullyQualifiedName(), err),
				},
			}
		}
		return diag.FromErr(err)
	}

	agent, err := collections.FindFirst(agents, func(agent sdk.SnowflakeIntelligenceAgent) bool {
		return agent.ID().FullyQualifiedName() == agentName.FullyQualifiedName()
	})
	if err != nil {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to find the attached Cortex agent. Marking the resource as removed.",
				Detail:   fmt.Sprintf("Snowflake Intelligence id: %s, Agent id: %s", intelligenceName.FullyQualifiedName(), agentName.FullyQualifiedName()),
			},
		}
	}

	if err := errors.Join(
		d.Set("snowflake_intelligence_name", intelligenceName.FullyQualifiedName()),
		d.Set("cortex_agent_name", agent.ID().FullyQualifiedName()),
	); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func DeleteIntelligenceCortexAgentAttachment(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	intelligenceName, agentName, err := parseIntelligenceCortexAgentAttachmentId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.SnowflakeIntelligences.DropAgentSafely(ctx, intelligenceName, agentName); err != nil {
		return diag.FromErr(fmt.Errorf("error while detaching agent %s from Snowflake Intelligence %s, err = %w", agentName.FullyQualifiedName(), intelligenceName.FullyQualifiedName(), err))
	}

	d.SetId("")

	return nil
}

func parseIntelligenceCortexAgentAttachmentId(id string) (sdk.AccountObjectIdentifier, sdk.SchemaObjectIdentifier, error) {
	parts := helpers.ParseResourceIdentifier(id)
	if len(parts) != 2 {
		return sdk.AccountObjectIdentifier{}, sdk.SchemaObjectIdentifier{}, fmt.Errorf("required id format '<snowflake_intelligence_name>|<cortex_agent_fqn>', but got: '%s'", id)
	}

	intelligenceName, err := sdk.ParseAccountObjectIdentifier(parts[0])
	if err != nil {
		return sdk.AccountObjectIdentifier{}, sdk.SchemaObjectIdentifier{}, err
	}
	agentName, err := sdk.ParseSchemaObjectIdentifier(parts[1])
	if err != nil {
		return sdk.AccountObjectIdentifier{}, sdk.SchemaObjectIdentifier{}, err
	}
	return intelligenceName, agentName, nil
}
