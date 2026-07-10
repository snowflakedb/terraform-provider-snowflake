package sdk

// IntegrationType is the type of integration.
type IntegrationType string

const (
	IntegrationTypeSecurityIntegrations       IntegrationType = "SECURITY INTEGRATIONS"
	IntegrationTypeAPIIntegrations            IntegrationType = "API INTEGRATIONS"
	IntegrationTypeStorageIntegrations        IntegrationType = "STORAGE INTEGRATIONS"
	IntegrationTypeExternalAccessIntegrations IntegrationType = "EXTERNAL ACCESS INTEGRATIONS"
	IntegrationTypeNotificationIntegrations   IntegrationType = "NOTIFICATION INTEGRATIONS"
)

type FailoverGroupSecondaryState string

const (
	FailoverGroupSecondaryStateSuspended FailoverGroupSecondaryState = "SUSPENDED"
	FailoverGroupSecondaryStateStarted   FailoverGroupSecondaryState = "STARTED"
	FailoverGroupSecondaryStateNull      FailoverGroupSecondaryState = "NULL"
)
