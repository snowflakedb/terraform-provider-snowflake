package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"log"
	"slices"
	"strings"
)

var (
	SecurityIntegrationDescribeSchema = collections.MergeMaps(
		DescribeApiAuthSecurityIntegrationSchema,
		DescribeExternalOauthSecurityIntegrationSchema,
		DescribeOauthIntegrationForCustomClients,
		DescribeOauthIntegrationForPartnerApplications,
		DescribeSaml2IntegrationSchema,
		DescribeScimSecurityIntegrationSchema,
	)
	allSecurityIntegrationPropertiesNames = slices.Concat(
		ApiAuthenticationPropertiesNames,
		ExternalOauthPropertiesNames,
		OauthIntegrationForCustomClientsPropertiesNames,
		OauthIntegrationForPartnerApplicationsPropertiesNames,
		Saml2PropertiesNames,
		ScimPropertiesNames,
	)
	InvalidIntegrationPropertyNames = []string{
		"OAUTH_CLIENT_ID",
		"OAUTH_REDIRECT_URI",
		"SAML2_SNOWFLAKE_X509_CERT",
		"SAML2_X509_CERT",
	}
)

func SecurityIntegrationsDescriptionsToSchema(integrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	securityIntegrationProperties := make(map[string]any)
	for _, desc := range integrationProperties {
		desc := desc
		if !slices.Contains(InvalidIntegrationPropertyNames, desc.Name) {
			if slices.Contains(allSecurityIntegrationPropertiesNames, desc.Name) {
				securityIntegrationProperties[strings.ToLower(desc.Name)] = []map[string]any{SecurityIntegrationPropertyToSchema(&desc)}
			} else {
				log.Printf("[WARN] unexpected property %v in security integration returned from Snowflake", desc.Name)
			}
		}
	}
	return securityIntegrationProperties
}
