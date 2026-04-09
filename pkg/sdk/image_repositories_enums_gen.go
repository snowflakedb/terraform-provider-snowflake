package sdk

import (
	"fmt"
	"slices"
	"strings"
)

type ImageRepositoryEncryptionType string

const (
	ImageRepositoryEncryptionTypeSnowflakeFull ImageRepositoryEncryptionType = "SNOWFLAKE_FULL"
	ImageRepositoryEncryptionTypeSnowflakeSse  ImageRepositoryEncryptionType = "SNOWFLAKE_SSE"
)

var AllImageRepositoryEncryptionTypes = []ImageRepositoryEncryptionType{
	ImageRepositoryEncryptionTypeSnowflakeFull,
	ImageRepositoryEncryptionTypeSnowflakeSse,
}

func ToImageRepositoryEncryptionType(s string) (ImageRepositoryEncryptionType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllImageRepositoryEncryptionTypes, ImageRepositoryEncryptionType(s)) {
		return "", fmt.Errorf("invalid image repository encryption type: %s", s)
	}
	return ImageRepositoryEncryptionType(s), nil
}
