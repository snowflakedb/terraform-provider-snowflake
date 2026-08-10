package sdk

import (
	"context"
	"errors"
	"strconv"
)

func (r *CreateExternalAccessIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}

func (d *ExternalAccessIntegrationDetails) ID() AccountObjectIdentifier {
	return d.Id
}

func (v *externalAccessIntegrations) DescribeDetails(ctx context.Context, id AccountObjectIdentifier) (*ExternalAccessIntegrationDetails, error) {
	properties, err := v.Describe(ctx, id)
	if err != nil {
		return nil, err
	}
	return parseExternalAccessIntegrationDetails(properties, id)
}

func parseExternalAccessIntegrationDetails(properties []ExternalAccessIntegrationProperty, id AccountObjectIdentifier) (*ExternalAccessIntegrationDetails, error) {
	details := &ExternalAccessIntegrationDetails{Id: id}
	var errs []error
	for _, prop := range properties {
		switch prop.Name {
		case "ENABLED":
			if v, err := strconv.ParseBool(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.Enabled = v
			}
		case "ALLOWED_NETWORK_RULES":
			if ids, err := ParseCommaSeparatedSchemaObjectIdentifierArray(prop.Value); err != nil {
				errs = append(errs, err)
			} else {
				details.AllowedNetworkRules = ids
			}
		case "ALLOWED_API_AUTHENTICATION_INTEGRATIONS":
			details.AllowedApiAuthenticationIntegrations = ParseCommaSeparatedStringArray(prop.Value, false)
		case "ALLOWED_AUTHENTICATION_SECRETS":
			raw := ParseCommaSeparatedStringArray(prop.Value, false)
			normalized := make([]string, 0, len(raw))
			for _, s := range raw {
				if id, err := ParseSchemaObjectIdentifier(s); err == nil {
					normalized = append(normalized, id.FullyQualifiedName())
				} else {
					normalized = append(normalized, s)
				}
			}
			details.AllowedAuthenticationSecrets = normalized
		case "COMMENT":
			details.Comment = prop.Value
		}
	}
	return details, errors.Join(errs...)
}
