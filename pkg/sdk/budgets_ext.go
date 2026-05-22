package sdk

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"

func NewBudgetSetEmailNotificationsArgsRequestFromEmails(emails ...string) *BudgetSetEmailNotificationsArgsRequest {
	return NewBudgetSetEmailNotificationsArgsRequest(collections.Map(emails, func(email string) BudgetEmailRequest { return BudgetEmailRequest{email} }))
}
