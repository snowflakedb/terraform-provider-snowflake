package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func UserDefaultWorkloadIdentityAwsVariable(arn string) tfconfig.Variable {
	return tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"aws": tfconfig.ListVariable(tfconfig.ObjectVariable(
				map[string]tfconfig.Variable{
					"arn": tfconfig.StringVariable(arn),
				},
			)),
		},
	)
}

func UserDefaultWorkloadIdentityGcpVariable(subject string) tfconfig.Variable {
	return tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"gcp": tfconfig.ListVariable(tfconfig.ObjectVariable(
				map[string]tfconfig.Variable{
					"subject": tfconfig.StringVariable(subject),
				},
			)),
		},
	)
}

func UserDefaultWorkloadIdentityAzureVariable(issuer, subject string) tfconfig.Variable {
	return tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"azure": tfconfig.ListVariable(tfconfig.ObjectVariable(
				map[string]tfconfig.Variable{
					"issuer":  tfconfig.StringVariable(issuer),
					"subject": tfconfig.StringVariable(subject),
				},
			)),
		},
	)
}

func UserDefaultWorkloadIdentityOidcVariable(issuer, subject string, audienceList ...string) tfconfig.Variable {
	m := map[string]tfconfig.Variable{
		"issuer":  tfconfig.StringVariable(issuer),
		"subject": tfconfig.StringVariable(subject),
	}
	if len(audienceList) > 0 {
		audiences := make([]tfconfig.Variable, len(audienceList))
		for i, a := range audienceList {
			audiences[i] = tfconfig.StringVariable(a)
		}
		m["oidc_audience_list"] = tfconfig.ListVariable(audiences...)
	}

	return tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"oidc": tfconfig.ListVariable(tfconfig.ObjectVariable(m)),
		},
	)
}

func UserDefaultWorkloadIdentityAwsEmpty() tfconfig.Variable {
	return tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"aws": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			})),
		},
	)
}

func UserDefaultWorkloadIdentityGcpEmpty() tfconfig.Variable {
	return tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"gcp": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			})),
		},
	)
}

func UserDefaultWorkloadIdentityAzureEmpty() tfconfig.Variable {
	return tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"azure": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			})),
		},
	)
}

func UserDefaultWorkloadIdentityOidcEmpty() tfconfig.Variable {
	return tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"oidc": tfconfig.ListVariable(tfconfig.ObjectVariable(map[string]tfconfig.Variable{
				"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
			})),
		},
	)
}

func UserDefaultWorkloadIdentityMultipleProvidersVariable() tfconfig.Variable {
	return tfconfig.ObjectVariable(
		map[string]tfconfig.Variable{
			"aws": tfconfig.ListVariable(tfconfig.ObjectVariable(
				map[string]tfconfig.Variable{
					"arn": tfconfig.StringVariable("foo"),
				},
			)),
			"gcp": tfconfig.ListVariable(tfconfig.ObjectVariable(
				map[string]tfconfig.Variable{
					"subject": tfconfig.StringVariable("bar"),
				},
			)),
		},
	)
}

func UserDefaultWorkloadIdentityEmpty() tfconfig.Variable {
	return tfconfig.ObjectVariable(map[string]tfconfig.Variable{
		"any": tfconfig.StringVariable(string(config.SnowflakeProviderConfigSingleAttributeWorkaround)),
	})
}
