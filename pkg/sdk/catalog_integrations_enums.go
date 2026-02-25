package sdk

type CatalogIntegrationCatalogSourceType string

const (
	CatalogIntegrationCatalogSourceTypeAWSGlue              CatalogIntegrationCatalogSourceType = "GLUE"
	CatalogIntegrationCatalogSourceTypeObjectStorage        CatalogIntegrationCatalogSourceType = "OBJECT_STORE"
	CatalogIntegrationCatalogSourceTypePolaris              CatalogIntegrationCatalogSourceType = "POLARIS"
	CatalogIntegrationCatalogSourceTypeIcebergREST          CatalogIntegrationCatalogSourceType = "ICEBERG_REST"
	CatalogIntegrationCatalogSourceTypeSAPBusinessDataCloud CatalogIntegrationCatalogSourceType = "SAP_BDC"
)

type CatalogIntegrationTableFormat string

const (
	CatalogIntegrationTableFormatIceberg CatalogIntegrationTableFormat = "ICEBERG"
	CatalogIntegrationTableFormatDelta   CatalogIntegrationTableFormat = "DELTA"
)

type CatalogIntegrationRestAuthenticationType string

const (
	CatalogIntegrationRestAuthenticationTypeOAuth  CatalogIntegrationRestAuthenticationType = "OAUTH"
	CatalogIntegrationRestAuthenticationTypeBearer CatalogIntegrationRestAuthenticationType = "BEARER"
	CatalogIntegrationRestAuthenticationTypeSigV4  CatalogIntegrationRestAuthenticationType = "SIGV4"
)

type CatalogIntegrationAccessDelegationMode string

const (
	CatalogIntegrationAccessDelegationModeVendedCredentials         CatalogIntegrationAccessDelegationMode = "VENDED_CREDENTIALS"
	CatalogIntegrationAccessDelegationModeExternalVolumeCredentials CatalogIntegrationAccessDelegationMode = "EXTERNAL_VOLUME_CREDENTIALS"
)

type CatalogIntegrationCatalogApiType string

const (
	CatalogIntegrationCatalogApiTypePublic               CatalogIntegrationCatalogApiType = "PUBLIC"
	CatalogIntegrationCatalogApiTypePrivate              CatalogIntegrationCatalogApiType = "PRIVATE"
	CatalogIntegrationCatalogApiTypeAwsApiGateway        CatalogIntegrationCatalogApiType = "AWS_API_GATEWAY"
	CatalogIntegrationCatalogApiTypeAwsPrivateApiGateway CatalogIntegrationCatalogApiType = "AWS_PRIVATE_API_GATEWAY"
	CatalogIntegrationCatalogApiTypeAwsGlue              CatalogIntegrationCatalogApiType = "AWS_GLUE"
	CatalogIntegrationCatalogApiTypeAwsPrivateGlue       CatalogIntegrationCatalogApiType = "AWS_PRIVATE_GLAUE"
)
