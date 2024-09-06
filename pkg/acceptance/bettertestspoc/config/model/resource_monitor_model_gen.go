// Code generated by config model builder generator; DO NOT EDIT.

package model

import (
	"reflect"
	"strings"

	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
)

type ResourceMonitorModel struct {
	CreditQuota        tfconfig.Variable `json:"credit_quota,omitempty"`
	EndTimestamp       tfconfig.Variable `json:"end_timestamp,omitempty"`
	Frequency          tfconfig.Variable `json:"frequency,omitempty"`
	FullyQualifiedName tfconfig.Variable `json:"fully_qualified_name,omitempty"`
	Name               tfconfig.Variable `json:"name,omitempty"`
	NotifyUsers        tfconfig.Variable `json:"notify_users,omitempty"`
	StartTimestamp     tfconfig.Variable `json:"start_timestamp,omitempty"`
	Trigger            tfconfig.Variable `json:"trigger,omitempty"`

	*config.ResourceModelMeta
}

/////////////////////////////////////////////////
// Basic builders (resource name and required) //
/////////////////////////////////////////////////

func ResourceMonitor(
	resourceName string,
	name string,
) *ResourceMonitorModel {
	r := &ResourceMonitorModel{ResourceModelMeta: config.Meta(resourceName, resources.ResourceMonitor)}
	r.WithName(name)
	return r
}

func ResourceMonitorWithDefaultMeta(
	name string,
) *ResourceMonitorModel {
	r := &ResourceMonitorModel{ResourceModelMeta: config.DefaultMeta(resources.ResourceMonitor)}
	r.WithName(name)
	return r
}

func (r *ResourceMonitorModel) ToConfigVariables() tfconfig.Variables {
	variables := make(tfconfig.Variables)
	rType := reflect.TypeOf(r).Elem()
	rValue := reflect.ValueOf(r).Elem()
	for i := 0; i < rType.NumField(); i++ {
		field := rType.Field(i)
		if jsonTag, ok := field.Tag.Lookup("json"); ok {
			name := strings.Split(jsonTag, ",")[0]
			if fieldValue, ok := rValue.Field(i).Interface().(tfconfig.Variable); ok {
				variables[name] = fieldValue
			}
		}
	}
	return variables
}

/////////////////////////////////
// below all the proper values //
/////////////////////////////////

func (r *ResourceMonitorModel) WithCreditQuota(creditQuota int) *ResourceMonitorModel {
	r.CreditQuota = tfconfig.IntegerVariable(creditQuota)
	return r
}

func (r *ResourceMonitorModel) WithEndTimestamp(endTimestamp string) *ResourceMonitorModel {
	r.EndTimestamp = tfconfig.StringVariable(endTimestamp)
	return r
}

func (r *ResourceMonitorModel) WithFrequency(frequency string) *ResourceMonitorModel {
	r.Frequency = tfconfig.StringVariable(frequency)
	return r
}

func (r *ResourceMonitorModel) WithFullyQualifiedName(fullyQualifiedName string) *ResourceMonitorModel {
	r.FullyQualifiedName = tfconfig.StringVariable(fullyQualifiedName)
	return r
}

func (r *ResourceMonitorModel) WithName(name string) *ResourceMonitorModel {
	r.Name = tfconfig.StringVariable(name)
	return r
}

// notify_users attribute type is not yet supported, so WithNotifyUsers can't be generated

func (r *ResourceMonitorModel) WithStartTimestamp(startTimestamp string) *ResourceMonitorModel {
	r.StartTimestamp = tfconfig.StringVariable(startTimestamp)
	return r
}

// trigger attribute type is not yet supported, so WithTrigger can't be generated

//////////////////////////////////////////
// below it's possible to set any value //
//////////////////////////////////////////

func (r *ResourceMonitorModel) WithCreditQuotaValue(value tfconfig.Variable) *ResourceMonitorModel {
	r.CreditQuota = value
	return r
}

func (r *ResourceMonitorModel) WithEndTimestampValue(value tfconfig.Variable) *ResourceMonitorModel {
	r.EndTimestamp = value
	return r
}

func (r *ResourceMonitorModel) WithFrequencyValue(value tfconfig.Variable) *ResourceMonitorModel {
	r.Frequency = value
	return r
}

func (r *ResourceMonitorModel) WithFullyQualifiedNameValue(value tfconfig.Variable) *ResourceMonitorModel {
	r.FullyQualifiedName = value
	return r
}

func (r *ResourceMonitorModel) WithNameValue(value tfconfig.Variable) *ResourceMonitorModel {
	r.Name = value
	return r
}

func (r *ResourceMonitorModel) WithNotifyUsersValue(value tfconfig.Variable) *ResourceMonitorModel {
	r.NotifyUsers = value
	return r
}

func (r *ResourceMonitorModel) WithStartTimestampValue(value tfconfig.Variable) *ResourceMonitorModel {
	r.StartTimestamp = value
	return r
}

func (r *ResourceMonitorModel) WithTriggerValue(value tfconfig.Variable) *ResourceMonitorModel {
	r.Trigger = value
	return r
}
