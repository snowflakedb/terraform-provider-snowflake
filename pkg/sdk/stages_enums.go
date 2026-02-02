package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type InternalStageEncryptionOption string

var (
	InternalStageEncryptionFull InternalStageEncryptionOption = "SNOWFLAKE_FULL"
	InternalStageEncryptionSSE  InternalStageEncryptionOption = "SNOWFLAKE_SSE"
)

type ExternalStageS3EncryptionOption string

var (
	ExternalStageS3EncryptionCSE       ExternalStageS3EncryptionOption = "AWS_CSE"
	ExternalStageS3EncryptionSSES3     ExternalStageS3EncryptionOption = "AWS_SSE_S3"
	ExternalStageS3EncryptionSSEKMS    ExternalStageS3EncryptionOption = "AWS_SSE_KMS"
	ExternalStageS3EncryptionNoneValue ExternalStageS3EncryptionOption = "NONE"
)

type ExternalStageGCSEncryptionOption string

var (
	ExternalStageGCSEncryptionSSEKMS    ExternalStageGCSEncryptionOption = "GCS_SSE_KMS"
	ExternalStageGCSEncryptionNoneValue ExternalStageGCSEncryptionOption = "NONE"
)

type ExternalStageAzureEncryptionOption string

var (
	ExternalStageAzureEncryptionCSE       ExternalStageAzureEncryptionOption = "AZURE_CSE"
	ExternalStageAzureEncryptionNoneValue ExternalStageAzureEncryptionOption = "NONE"
)

type StageCopyColumnMapOption string

var (
	StageCopyColumnMapCaseSensitive   StageCopyColumnMapOption = "CASE_SENSITIVE"
	StageCopyColumnMapCaseInsensitive StageCopyColumnMapOption = "CASE_INSENSITIVE"
	StageCopyColumnMapCaseNone        StageCopyColumnMapOption = "NONE"
)

type StageType string

var (
	StageTypeInternal          StageType = "INTERNAL"
	StageTypeInternalNoCse     StageType = "INTERNAL NO CSE"
	StageTypeInternalTemporary StageType = "INTERNAL TEMPORARY"
	StageTypeExternal          StageType = "EXTERNAL"
	StageTypeExternalTemporary StageType = "EXTERNAL TEMPORARY"
)

var allStageTypes = []StageType{
	StageTypeInternal,
	StageTypeInternalNoCse,
	StageTypeInternalTemporary,
	StageTypeExternal,
	StageTypeExternalTemporary,
}

func ToStageType(s string) (StageType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allStageTypes, StageType(s)) {
		return "", fmt.Errorf("invalid stage type: %s", s)
	}
	return StageType(s), nil
}

func (s StageType) Canonical() StageType {
	switch s {
	case StageTypeInternalNoCse, StageTypeInternal, StageTypeInternalTemporary:
		return StageTypeInternal
	case StageTypeExternal, StageTypeExternalTemporary:
		return StageTypeExternal
	default:
		return s
	}
}
