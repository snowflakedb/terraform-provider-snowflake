package sdk

import (
	"fmt"
	"strings"
)

type StreamSourceType string

const (
	StreamSourceTypeTable         StreamSourceType = "TABLE"
	StreamSourceTypeExternalTable StreamSourceType = "EXTERNAL TABLE"
	StreamSourceTypeView          StreamSourceType = "VIEW"
	StreamSourceTypeStage         StreamSourceType = "STAGE"
)

func ToStreamSourceType(s string) (StreamSourceType, error) {
	switch streamSourceType := StreamSourceType(strings.ToUpper(s)); streamSourceType {
	case StreamSourceTypeTable,
		StreamSourceTypeExternalTable,
		StreamSourceTypeView,
		StreamSourceTypeStage:
		return streamSourceType, nil
	default:
		return "", fmt.Errorf("invalid stream source type: %s", s)
	}
}

type StreamMode string

const (
	StreamModeDefault    StreamMode = "DEFAULT"
	StreamModeAppendOnly StreamMode = "APPEND_ONLY"
	StreamModeInsertOnly StreamMode = "INSERT_ONLY"
)

func ToStreamMode(s string) (StreamMode, error) {
	switch streamMode := StreamMode(strings.ToUpper(s)); streamMode {
	case StreamModeDefault,
		StreamModeAppendOnly,
		StreamModeInsertOnly:
		return streamMode, nil
	default:
		return "", fmt.Errorf("invalid stream mode: %s", s)
	}
}
