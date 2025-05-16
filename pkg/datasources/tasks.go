package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var tasksSchema = map[string]*schema.Schema{
	"with_parameters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs SHOW PARAMETERS FOR TASK for each task returned by SHOW TASK and saves the output to the parameters field as a map. By default this value is set to true.",
	},
	"like":        likeSchema,
	"in":          extendedInSchema,
	"starts_with": startsWithSchema,
	"root_only": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Filters the command output to return only root tasks (tasks with no predecessors).",
	},
	"limit": limitFromSchema,
	"tasks": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all task details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW TASKS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowTaskSchema,
					},
				},
				resources.ParametersAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PARAMETERS FOR TASK.",
					Elem: &schema.Resource{
						Schema: schemas.ShowTaskParametersSchema,
					},
				},
			},
		},
	},
}

func Tasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Tasks, ReadTasks),
		Schema:      tasksSchema,
		Description: "Data source used to get details of filtered tasks. Filtering is aligned with the current possibilities for [SHOW TASKS](https://docs.snowflake.com/en/sql-reference/sql/show-tasks) query. The results of SHOW and SHOW PARAMETERS IN are encapsulated in one output collection `tasks`.",
	}
}

func ReadTasks(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowTaskRequest()

	handleLike(d, &req.Like)
	if err := handleExtendedIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}
	handleStartsWith(d, &req.StartsWith)
	if v, ok := d.GetOk("root_only"); ok && v.(bool) {
		req.WithRootOnly(true)
	}
	handleLimitFrom(d, &req.Limit)

	tasks, err := client.Tasks.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("tasks_read")

	flattenedTasks := make([]map[string]any, len(tasks))
	for i, task := range tasks {
		task := task

		var taskParameters []map[string]any
		if d.Get("with_parameters").(bool) {
			parameters, err := client.Tasks.ShowParameters(ctx, task.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			taskParameters = []map[string]any{schemas.TaskParametersToSchema(parameters)}
		}

		flattenedTasks[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.TaskToSchema(&task)},
			resources.ParametersAttributeName: taskParameters,
		}
	}

	if err := d.Set("tasks", flattenedTasks); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
