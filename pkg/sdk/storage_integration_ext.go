package sdk

import (
	"fmt"
	"strings"
)

type S3Protocol string

const (
	RegularS3Protocol S3Protocol = "S3"
	GovS3Protocol     S3Protocol = "S3GOV"
	ChinaS3Protocol   S3Protocol = "S3CHINA"
)

var (
	AllS3Protocols      = []S3Protocol{RegularS3Protocol, GovS3Protocol, ChinaS3Protocol}
	AllStorageProviders = append(AsStringList(AllS3Protocols), "GCS", "AZURE")
)

func ToS3Protocol(s string) (S3Protocol, error) {
	switch protocol := S3Protocol(strings.ToUpper(s)); protocol {
	case RegularS3Protocol, GovS3Protocol, ChinaS3Protocol:
		return protocol, nil
	default:
		return "", fmt.Errorf("invalid S3 protocol: %s", s)
	}
}
