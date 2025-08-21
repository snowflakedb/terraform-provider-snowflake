package common

import (
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
)

func GetAccessToken() string {
	token := oswrapper.Getenv("SF_TF_SCRIPT_GH_ACCESS_TOKEN")
	if token == "" {
		panic(errors.New("GitHub access token missing"))
	}
	return token
}
