package sdk

import (
	"fmt"
	"strings"
)

type TaskState string

const (
	TaskStateStarted   TaskState = "started"
	TaskStateSuspended TaskState = "suspended"
)

func ToTaskState(s string) (TaskState, error) {
	switch taskState := TaskState(strings.ToLower(s)); taskState {
	case TaskStateStarted, TaskStateSuspended:
		return taskState, nil
	default:
		return "", fmt.Errorf("unknown task state: %s", s)
	}
}
