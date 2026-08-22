package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandNotificationConfiguration(raw interface{}) *sdk.NotificationConfiguration {
	notificationList := raw.([]interface{})
	if len(notificationList) == 0 {
		return nil
	}
	notificationMap := notificationList[0].(map[string]interface{})
	notification := &sdk.NotificationConfiguration{}

	if notifyAdmins, ok := notificationMap["notify_admins"].(bool); ok {
		notification.NotifyAdmins = &notifyAdmins
	}
	if severityThreshold, ok := notificationMap["severity_threshold"].(string); ok && severityThreshold != "" {
		notification.SeverityThreshold = &severityThreshold
	}
	if usersSet, ok := notificationMap["users"].(*schema.Set); ok && usersSet.Len() > 0 {
		users := make([]string, 0, usersSet.Len())
		for _, u := range usersSet.List() {
			users = append(users, u.(string))
		}
		notification.Users = users
	}
	return notification
}

func flattenNotificationConfiguration(notificationJSON string) ([]interface{}, error) {
	if notificationJSON == "" {
		return nil, nil
	}
	notification, err := sdk.ParseNotificationConfiguration(notificationJSON)
	if err != nil {
		return nil, err
	}
	if notification == nil {
		return nil, nil
	}
	m := map[string]interface{}{
		"notify_admins":      notification.NotifyAdmins != nil && *notification.NotifyAdmins,
		"severity_threshold": "",
		"users":              notification.Users,
	}
	if notification.SeverityThreshold != nil {
		m["severity_threshold"] = *notification.SeverityThreshold
	}
	return []interface{}{m}, nil
}
