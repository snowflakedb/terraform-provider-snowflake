package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (taskToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"predecessors": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"task_relations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"predecessors": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"finalizer": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"finalized_root_task": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"target_completion_interval": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"hours": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"minutes": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"seconds": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func (taskToSchemaMapper) additionalToSchema(src *sdk.Task, dst map[string]any) {
	// Identifier slice: map each predecessor to FullyQualifiedName.
	dst["predecessors"] = collections.Map(src.Predecessors, sdk.SchemaObjectIdentifier.FullyQualifiedName)
	// Nested TaskRelations struct is a single-item list, not the struct itself.
	finalizer := ""
	if src.TaskRelations.FinalizerTask != nil {
		finalizer = src.TaskRelations.FinalizerTask.FullyQualifiedName()
	}
	finalizedRootTask := ""
	if src.TaskRelations.FinalizedRootTask != nil {
		finalizedRootTask = src.TaskRelations.FinalizedRootTask.FullyQualifiedName()
	}
	dst["task_relations"] = []any{
		map[string]any{
			"predecessors":        collections.Map(src.TaskRelations.Predecessors, sdk.SchemaObjectIdentifier.FullyQualifiedName),
			"finalizer":           finalizer,
			"finalized_root_task": finalizedRootTask,
		},
	}
	// Nested hours/minutes/seconds; skip when unset.
	if src.TargetCompletionInterval != nil {
		dst["target_completion_interval"] = []any{
			map[string]any{
				"hours":   src.TargetCompletionInterval.Hours,
				"minutes": src.TargetCompletionInterval.Minutes,
				"seconds": src.TargetCompletionInterval.Seconds,
			},
		}
	}
}
