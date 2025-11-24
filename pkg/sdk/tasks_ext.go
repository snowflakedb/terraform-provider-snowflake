package sdk

import "encoding/json"

type TaskRelationsRepresentation struct {
	Predecessors      []string `json:"Predecessors"`
	FinalizerTask     string   `json:"FinalizerTask"`
	FinalizedRootTask string   `json:"FinalizedRootTask"`
}

func (r *TaskRelationsRepresentation) ToTaskRelations() (TaskRelations, error) {
	predecessors := make([]SchemaObjectIdentifier, len(r.Predecessors))
	for i, predecessor := range r.Predecessors {
		id, err := ParseSchemaObjectIdentifier(predecessor)
		if err != nil {
			return TaskRelations{}, err
		}
		predecessors[i] = id
	}

	taskRelations := TaskRelations{
		Predecessors: predecessors,
	}

	if len(r.FinalizerTask) > 0 {
		finalizerTask, err := ParseSchemaObjectIdentifier(r.FinalizerTask)
		if err != nil {
			return TaskRelations{}, err
		}
		taskRelations.FinalizerTask = &finalizerTask
	}

	if len(r.FinalizedRootTask) > 0 {
		finalizedRootTask, err := ParseSchemaObjectIdentifier(r.FinalizedRootTask)
		if err != nil {
			return TaskRelations{}, err
		}
		taskRelations.FinalizedRootTask = &finalizedRootTask
	}

	return taskRelations, nil
}

type TaskRelations struct {
	Predecessors      []SchemaObjectIdentifier
	FinalizerTask     *SchemaObjectIdentifier
	FinalizedRootTask *SchemaObjectIdentifier
}

func ToTaskRelations(s string) (TaskRelations, error) {
	var taskRelationsRepresentation TaskRelationsRepresentation
	if err := json.Unmarshal([]byte(s), &taskRelationsRepresentation); err != nil {
		return TaskRelations{}, err
	}
	taskRelations, err := taskRelationsRepresentation.ToTaskRelations()
	if err != nil {
		return TaskRelations{}, err
	}
	return taskRelations, nil
}
