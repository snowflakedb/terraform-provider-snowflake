package model

import (
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

func Saml2SecurityIntegration2(
	resourceName string,
	name string,
	saml2Issuer string,
	saml2Provider string,
	saml2SsoUrl string,
	saml2CertVariableName string,
) *Saml2SecurityIntegrationModel {
	s := &Saml2SecurityIntegrationModel{ResourceModelMeta: config.Meta(resourceName, resources.Saml2SecurityIntegration)}
	s.WithName(name)
	s.WithSaml2Issuer(saml2Issuer)
	s.WithSaml2Provider(saml2Provider)
	s.WithSaml2SsoUrl(saml2SsoUrl)
	s.WithSaml2X509Cert(fmt.Sprintf("%[1]s%[2]s%[1]s", config.SnowflakeProviderConfigUnquoteMarker, saml2CertVariableName))
	return s
}

func (s *Saml2SecurityIntegrationModel) WithAllowedEmailPatterns(values ...string) *Saml2SecurityIntegrationModel {
	s.AllowedEmailPatterns = tfconfig.SetVariable(
		collections.Map(values, func(value string) tfconfig.Variable {
			return tfconfig.StringVariable(value)
		})...,
	)
	return s
}

func (s *Saml2SecurityIntegrationModel) WithAllowedUserDomains(values ...string) *Saml2SecurityIntegrationModel {
	s.AllowedUserDomains = tfconfig.SetVariable(
		collections.Map(values, func(value string) tfconfig.Variable {
			return tfconfig.StringVariable(value)
		})...,
	)
	return s
}
