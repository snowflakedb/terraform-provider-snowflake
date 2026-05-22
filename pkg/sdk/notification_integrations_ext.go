package sdk

func (s *NotificationIntegrationSet) additionalValidations() error {
	// TODO [next PR]: this validation type will be added (conflicting fields currently works only for 2 fields)
	if moreThanOneValueSet(s.SetPushParams, s.SetEmailParams, s.SetWebhookParams) {
		return errOneOf("AlterNotificationIntegrationOptions.Set", "SetPushParams", "SetEmailParams", "SetWebhookParams")
	}
	return nil
}

func (r *CreateNotificationIntegrationRequest) GetName() AccountObjectIdentifier {
	return r.name
}
