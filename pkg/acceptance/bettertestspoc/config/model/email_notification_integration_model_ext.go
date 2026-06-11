package model

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
)

func (e *EmailNotificationIntegrationModel) WithAllowedRecipients(emails ...string) *EmailNotificationIntegrationModel {
	return e.WithAllowedRecipientsValue(
		tfconfig.SetVariable(
			collections.Map(emails, func(email string) tfconfig.Variable { return tfconfig.StringVariable(email) })...,
		),
	)
}
