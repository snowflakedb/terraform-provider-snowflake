package provider

import (
	"context"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &ResourceMonitorResource{}
	_ resource.ResourceWithImportState = &ResourceMonitorResource{}
)

func NewResourceMonitorResource() resource.Resource {
	return &ResourceMonitorResource{}
}

type ResourceMonitorResource struct {
	client *sdk.Client
}

type resourceMonitorModelV0 struct {
	Name                     types.String  `tfsdk:"name"`
	NotifyUsers              types.Set     `tfsdk:"notify_users"`
	CreditQuota              types.Float64 `tfsdk:"credit_quota"`
	Frequency                types.String  `tfsdk:"frequency"`
	StartTimestamp           types.String  `tfsdk:"start_timestamp"`
	EndTimestamp             types.String  `tfsdk:"end_timestamp"`
	SuspendTrigger           types.Int64   `tfsdk:"suspend_trigger"`
	SuspendTriggers          types.Set     `tfsdk:"suspend_triggers"`
	SuspendImmediateTrigger  types.Int64   `tfsdk:"suspend_immediate_trigger"`
	SuspendImmediateTriggers types.Set     `tfsdk:"suspend_immediate_triggers"`
	NotifyTriggers           types.Set     `tfsdk:"notify_triggers"`
	SetForAccount            types.Bool    `tfsdk:"set_for_account"`
	Warehouses               types.Set     `tfsdk:"warehouses"`
	Id                       types.String  `tfsdk:"id"`
}

func resourceMonitorSchemaV0() schema.Schema {
	return schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Identifier for the resource monitor; must be unique for your account.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"notify_users": schema.SetAttribute{
				Description: "Specifies the list of users to receive email notifications on resource monitors.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"credit_quota": schema.Int64Attribute{
				Description: "The number of credits allocated monthly to the resource monitor.",
				Optional:    true,
				Computed:    true,
			},
			"frequency": schema.StringAttribute{
				Description: "The frequency interval at which the credit usage resets to 0. If you set a frequency for a resource monitor, you must also set START_TIMESTAMP.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive([]string{"MONTHLY", "DAILY", "WEEKLY", "YEARLY", "NEVER"}...),
				},
			},
			"start_timestamp": schema.StringAttribute{
				Description: "The date and time when the resource monitor starts monitoring credit usage for the assigned warehouses.",
				Optional:    true,
				Computed:    true,
			},
			"end_timestamp": schema.StringAttribute{
				Description: "The date and time when the resource monitor suspends the assigned warehouses.",
				Optional:    true,
				Computed:    true,
			},
			"suspend_trigger": schema.Int64Attribute{
				Description: "The number that represents the percentage threshold at which to suspend all warehouses.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("suspend_triggers")),
				},
			},
			"suspend_triggers": schema.SetAttribute{
				Description: "A list of percentage thresholds at which to suspend all warehouses.",
				Optional:    true,
				ElementType: types.Int64Type,
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.MatchRoot("suspend_trigger")),
				},
				DeprecationMessage: "Use suspend_trigger instead",
			},
			"suspend_immediate_trigger": schema.Int64Attribute{
				Description: "The number that represents the percentage threshold at which to immediately suspend all warehouses.", Optional: true,
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("suspend_immediate_triggers")),
				},
			},
			"suspend_immediate_triggers": schema.SetAttribute{
				Description: "A list of percentage thresholds at which to suspend all warehouses.",
				Optional:    true,
				ElementType: types.Int64Type,
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.MatchRoot("suspend_immediate_trigger")),
				},
				DeprecationMessage: "Use suspend_immediate_trigger instead",
			},
			"notify_triggers": schema.SetAttribute{
				Description: "A list of percentage thresholds at which to send an alert to subscribed users.",
				Optional:    true,
				ElementType: types.Int64Type,
			},
			"set_for_account": schema.BoolAttribute{
				Description: "Specifies whether the resource monitor should be applied globally to your Snowflake account (defaults to false).",
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				// todo: create a snowflake_resource_monitor_association resource
				// DeprecationMessage: "Use snowflake_resource_monitor_association instead",
			},
			"warehouses": schema.SetAttribute{
				Description: "A list of warehouses to apply the resource monitor to.",
				Optional:    true,
				ElementType: types.StringType,
				// todo: add the `resource_monitor` attribute to the `snowflake_warehouse` resource
				// DeprecationMessage: "Set the `resource_monitor` attribute on the `snowflake_warehouse` resource instead",
			},
		},
	}
}

func upgradeResourceMonitorStateV0toV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var resourceMonitorDataV0 resourceMonitorModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &resourceMonitorDataV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := resourceMonitorDataV0.Name
	notifyUsers := resourceMonitorDataV0.NotifyUsers
	creditQuota := resourceMonitorDataV0.CreditQuota
	frequency := resourceMonitorDataV0.Frequency
	startTimestamp := resourceMonitorDataV0.StartTimestamp
	endTimestamp := resourceMonitorDataV0.EndTimestamp
	suspendTrigger := resourceMonitorDataV0.SuspendTrigger
	suspendTriggers := resourceMonitorDataV0.SuspendTriggers
	if !suspendTriggers.IsNull() {
		suspendTriggersElements := make([]types.Int64, 0, len(suspendTriggers.Elements()))
		suspendTriggers.ElementsAs(ctx, &suspendTriggersElements, false)
		if len(suspendTriggersElements) > 0 {
			suspendTrigger = suspendTriggersElements[0]
		}
	}
	suspendImmediateTrigger := resourceMonitorDataV0.SuspendImmediateTrigger
	suspendImmediateTriggers := resourceMonitorDataV0.SuspendImmediateTriggers
	if !suspendImmediateTriggers.IsNull() {
		suspendImmediateTriggersElements := make([]types.Int64, 0, len(suspendImmediateTriggers.Elements()))
		suspendImmediateTriggers.ElementsAs(ctx, &suspendImmediateTriggersElements, false)
		if len(suspendImmediateTriggersElements) > 0 {
			suspendImmediateTrigger = suspendImmediateTriggersElements[0]
		}
	}
	notifyTriggers := resourceMonitorDataV0.NotifyTriggers

	trigggers := make([]resourceMonitorTriggerModel, 0)
	if !suspendTrigger.IsNull() {
		trigggers = append(trigggers, resourceMonitorTriggerModel{
			Threshold:     suspendTrigger,
			TriggerAction: types.StringValue("SUSPEND"),
		})
	}
	if !suspendImmediateTrigger.IsNull() {
		trigggers = append(trigggers, resourceMonitorTriggerModel{
			Threshold:     suspendImmediateTrigger,
			TriggerAction: types.StringValue("SUSPEND_IMMEDIATE"),
		})
	}
	if !notifyTriggers.IsNull() {
		notifyTriggersElements := make([]types.Int64, 0, len(notifyTriggers.Elements()))
		notifyTriggers.ElementsAs(ctx, &notifyTriggersElements, false)
		for _, v := range notifyTriggersElements {
			trigggers = append(trigggers, resourceMonitorTriggerModel{
				Threshold:     v,
				TriggerAction: types.StringValue("NOTIFY"),
			})
		}
	}
	triggersObjectType := types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"threshold":      types.Int64Type,
		"trigger_action": types.StringType,
	})
	triggersSet, diags := types.SetValueFrom(ctx, triggersObjectType, trigggers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceMonitorV1 := &resourceMonitorModelV1{
		Name:           name,
		NotifyUsers:    notifyUsers,
		CreditQuota:    creditQuota,
		Frequency:      frequency,
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
		Triggers:       triggersSet,
		Id:             resourceMonitorDataV0.Id,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, resourceMonitorV1)...)
}

type resourceMonitorModelV1 struct {
	OrReplace        types.Bool    `tfsdk:"or_replace"`
	Name             types.String  `tfsdk:"name"`
	CreditQuota      types.Float64 `tfsdk:"credit_quota"`
	UsedCredits      types.Float64 `tfsdk:"used_credits"`
	RemainingCredits types.Float64 `tfsdk:"remaining_credits"`
	Frequency        types.String  `tfsdk:"frequency"`
	StartTimestamp   types.String  `tfsdk:"start_timestamp"`
	EndTimestamp     types.String  `tfsdk:"end_timestamp"`
	Level            types.String  `tfsdk:"level"`
	NotifyUsers      types.Set     `tfsdk:"notify_users"`
	Triggers         types.Set     `tfsdk:"triggers"`
	Id               types.String  `tfsdk:"id"`
}

func resourceMonitorSchemaV1() schema.Schema {
	return schema.Schema{
		Description: "Snowflake resource monitor resource",
		Version:     1,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID of the database",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"or_replace": schema.BoolAttribute{
				Description: "Specifies whether to replace the resource monitor if it exists and has the same name as the one being created",
				Optional:    true,
				Computed:    true,
				Sensitive:   isSensitive("snowflake_resource_monitor.*.or_replace"),
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Specifies the object identifier for the database",
				Required:    true,
				Sensitive:   isSensitive("snowflake_resource_monitor.*.name"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"credit_quota": schema.Float64Attribute{
				Description: "The number of credits allocated to the resource monitor per frequency interval.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(0),
				Sensitive:   isSensitive("snowflake_resource_monitor.*.credit_quota"),
			},
			"used_credits": schema.Float64Attribute{
				Description: "The number of credits used by the resource monitor.",
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"remaining_credits": schema.Float64Attribute{
				Description: "The number of credits remaining for the resource monitor.",
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"level": schema.StringAttribute{
				Description: "resource monitor level",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"frequency": schema.StringAttribute{
				Description: "Specifies the maximum number of days to extend the Fail-safe storage retention period for the database",
				Optional:    true,
				Computed:    true,
				Sensitive:   isSensitive("snowflake_resource_monitor.*.frequency"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"MONTHLY", "DAILY", "WEEKLY", "YEARLY", "NEVER"}...),
					stringvalidator.AlsoRequires(path.MatchRoot("start_timestamp")),
				},
			},
			"start_timestamp": schema.StringAttribute{
				Description: "Specifies the start time of the resource monitor",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("frequency")),
				},
			},
			"end_timestamp": schema.StringAttribute{
				Description: "Specifies the end time of the resource monitor",
				Optional:    true,
			},
			"notify_users": schema.SetAttribute{
				Description: "Specifies the list of users to receive email notifications on resource monitors",
				Optional:    true,
				ElementType: types.StringType,
			},
			"triggers": schema.SetNestedAttribute{
				Description: "Specifies the list of triggers to receive email notifications on resource monitors",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"threshold": schema.Int64Attribute{
							Description: "Specifies the percentage of credits used to trigger an email notification",
							Required:    true,
							Validators: []validator.Int64{
								int64validator.AtLeast(0),
							},
						},
						"trigger_action": schema.StringAttribute{
							Description: "Specifies the action to take when the trigger is activated",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"SUSPEND", "SUSPEND_IMMEDIATE", "NOTIFY"}...),
							},
						},
					},
				},
			},
		},
	}
}

type resourceMonitorTriggerModel struct {
	Threshold     types.Int64  `tfsdk:"threshold"`
	TriggerAction types.String `tfsdk:"trigger_action"`
}

func (old *resourceMonitorModelV1) Equals(new *resourceMonitorModelV1, ctx context.Context) bool {
	if old == nil || new == nil {
		return false
	}
	if !old.Id.Equal(new.Id) {
		return false
	}
	if !old.OrReplace.Equal(new.OrReplace) {
		return false
	}
	if !old.Name.Equal(new.Name) {
		return false
	}
	if !old.CreditQuota.Equal(new.CreditQuota) {
		return false
	}
	if !old.Frequency.Equal(new.Frequency) {
		return false
	}
	if !old.StartTimestamp.Equal(new.StartTimestamp) {
		return false
	}
	if !old.EndTimestamp.Equal(new.EndTimestamp) {
		return false
	}
	if !old.Triggers.Equal(new.Triggers) {
		return false
	}
	if !old.NotifyUsers.Equal(new.NotifyUsers) {
		return false
	}

	return true
}

func (r *ResourceMonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_monitor"
}

func (r *ResourceMonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceMonitorSchemaV1()
}

func (r *ResourceMonitorResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	schemaV0 := resourceMonitorSchemaV0()
	return map[int64]resource.StateUpgrader{
		// State upgrade implementation from 0 to 1
		0: {
			PriorSchema:   &schemaV0,
			StateUpgrader: upgradeResourceMonitorStateV0toV1,
		},
	}
}

func (r *ResourceMonitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*ProviderData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sdk.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = providerData.client
}

func (r *ResourceMonitorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// we aren't really modifying the plan, just logging what the plan intends to do
	resp.Plan = req.Plan
	var plan, state *resourceMonitorModelV1
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resourceName := "snowflake_database"
	// DELETE
	if req.Plan.Raw.IsNull() {
		_, readLogs, _ := r.read(ctx, state, true)
		_, deleteLogs, _ := r.delete(ctx, state, true)
		deleteLogs = append(deleteLogs, readLogs...)
		tflog.Debug(ctx, formatSQLPreview(DeleteOperation, resourceName, state.Id.ValueString(), deleteLogs))
		return
	}

	// CREATE
	if plan.Id.IsUnknown() {
		_, createLogs, _ := r.create(ctx, plan, true)
		plan.Id = types.StringValue(sdk.NewAccountObjectIdentifier(plan.Name.ValueString()).FullyQualifiedName())
		_, readLogs, _ := r.read(ctx, plan, true)
		createLogs = append(createLogs, readLogs...)
		tflog.Debug(ctx, formatSQLPreview(CreateOperation, resourceName, "", createLogs))
		return
	}

	if plan.Equals(state, ctx) {
		// READ
		_, logs, _ := r.read(ctx, state, true)
		tflog.Debug(ctx, formatSQLPreview(ReadOperation, resourceName, state.Id.ValueString(), logs))
		return
	} else {
		// UPDATE
		_, updateLogs, _ := r.update(ctx, plan, state, true)
		_, readLogs, _ := r.read(ctx, plan, true)
		updateLogs = append(updateLogs, readLogs...)
		tflog.Debug(ctx, formatSQLPreview(UpdateOperation, resourceName, state.Id.ValueString(), updateLogs))
	}
}

func (r *ResourceMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *resourceMonitorModelV1
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	data, _, diags := r.create(ctx, data, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceMonitorResource) create(ctx context.Context, data *resourceMonitorModelV1, dryRun bool) (*resourceMonitorModelV1, []string, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	client := r.client
	if dryRun {
		client = sdk.NewDryRunClient()
	}

	name := data.Name.ValueString()

	id := sdk.NewAccountObjectIdentifier(name)

	opts := &sdk.CreateResourceMonitorOptions{
		OrReplace: data.OrReplace.ValueBoolPointer(),
	}

	with := &sdk.ResourceMonitorWith{}
	setWith := false
	if !data.CreditQuota.IsNull() && !data.CreditQuota.IsUnknown() && data.CreditQuota.ValueFloat64() > 0 {
		setWith = true
		with.CreditQuota = sdk.Int(int(data.CreditQuota.ValueFloat64()))
	}
	if !data.Frequency.IsNull() && data.Frequency.ValueString() != "" {
		setWith = true
		frequency, err := sdk.ToResourceMonitorFrequency(data.Frequency.ValueString())
		if err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Unable to create resource monitor, got error: %s", err))
		}
		with.Frequency = frequency
	}
	if !data.StartTimestamp.IsNull() && data.StartTimestamp.ValueString() != "" {
		setWith = true
		with.StartTimestamp = data.StartTimestamp.ValueStringPointer()
	}

	if !data.EndTimestamp.IsNull() && data.EndTimestamp.ValueString() != "" {
		setWith = true
		with.EndTimestamp = data.EndTimestamp.ValueStringPointer()
	}

	if !data.NotifyUsers.IsNull() && len(data.NotifyUsers.Elements()) > 0 {
		setWith = true
		elements := make([]types.String, 0, len(data.NotifyUsers.Elements()))
		var notifiedUsers []sdk.NotifiedUser
		for _, e := range elements {
			notifiedUsers = append(notifiedUsers, sdk.NotifiedUser{Name: sdk.NewAccountObjectIdentifier(e.ValueString())})
		}
		with.NotifyUsers = &sdk.NotifyUsers{
			Users: notifiedUsers,
		}
	}

	if !data.Triggers.IsNull() && len(data.Triggers.Elements()) > 0 {
		setWith = true
		elements := make([]resourceMonitorTriggerModel, 0, len(data.Triggers.Elements()))
		data.Triggers.ElementsAs(ctx, &elements, false)
		var triggers []sdk.TriggerDefinition
		for _, e := range elements {
			triggers = append(triggers, sdk.TriggerDefinition{
				Threshold:     int(e.Threshold.ValueInt64()),
				TriggerAction: sdk.TriggerAction(e.TriggerAction.ValueString()),
			})
		}
		with.Triggers = triggers
	}

	if setWith {
		opts.With = with
	}
	err := client.ResourceMonitors.Create(ctx, id, opts)

	if dryRun {
		return data, client.TraceLogs(), diags
	}
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Unable to create resource monitor, got error: %s", err))
	}

	data.Id = types.StringValue(id.FullyQualifiedName())
	r.read(ctx, data, false)
	return data, nil, diags
}

func (r *ResourceMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resourceMonitorModelV1
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data, _, diags := r.read(ctx, data, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceMonitorResource) read(ctx context.Context, data *resourceMonitorModelV1, dryRun bool) (*resourceMonitorModelV1, []string, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	client := r.client
	if dryRun {
		client = sdk.NewDryRunClient()
	}

	id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(data.Id.ValueString())
	resourceMonitor, err := client.ResourceMonitors.ShowByID(ctx, id)
	if dryRun {
		return data, client.TraceLogs(), diags
	}
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err))
		return data, nil, diags
	}

	data.CreditQuota = types.Float64Value(resourceMonitor.CreditQuota)
	data.Frequency = types.StringValue(string(resourceMonitor.Frequency))
	switch *resourceMonitor.Level {
	case sdk.ResourceMonitorLevelAccount:
		data.Level = types.StringValue("ACCOUNT")
	case sdk.ResourceMonitorLevelWarehouse:
		data.Level = types.StringValue("WAREHOUSE")
	}
	data.UsedCredits = types.Float64Value(resourceMonitor.UsedCredits)
	data.RemainingCredits = types.Float64Value(resourceMonitor.RemainingCredits)

	if resourceMonitor.StartTime != "" {
		if data.StartTimestamp.ValueString() != "IMMEDIATELY" {
			data.StartTimestamp = types.StringValue(resourceMonitor.StartTime)
		}
	} else {
		data.StartTimestamp = types.StringNull()
	}
	if resourceMonitor.EndTime != "" {
		data.EndTimestamp = types.StringValue(resourceMonitor.EndTime)
	}
	if len(resourceMonitor.NotifyUsers) == 0 {
		data.NotifyUsers = types.SetNull(types.StringType)
	} else {
		var notifyUsers []types.String
		for _, e := range resourceMonitor.NotifyUsers {
			notifyUsers = append(notifyUsers, types.StringValue(e))
		}
		var diag diag.Diagnostics
		data.NotifyUsers, diag = types.SetValueFrom(ctx, types.StringType, notifyUsers)
		diags = append(diags, diag...)
	}

	triggersObjectType := types.ObjectType{}.WithAttributeTypes(map[string]attr.Type{
		"threshold":      types.Int64Type,
		"trigger_action": types.StringType,
	})
	if len(resourceMonitor.NotifyAt) == 0 && resourceMonitor.SuspendAt == nil && resourceMonitor.SuspendImmediateAt == nil {
		data.Triggers = types.SetNull(triggersObjectType)
	} else {
		var triggers []resourceMonitorTriggerModel
		for _, e := range resourceMonitor.NotifyAt {
			triggers = append(triggers, resourceMonitorTriggerModel{
				Threshold:     types.Int64Value(int64(e)),
				TriggerAction: types.StringValue(string(sdk.TriggerActionNotify)),
			})
		}
		if resourceMonitor.SuspendAt != nil {
			triggers = append(triggers, resourceMonitorTriggerModel{
				Threshold:     types.Int64Value(int64(*resourceMonitor.SuspendAt)),
				TriggerAction: types.StringValue(string(sdk.TriggerActionSuspend)),
			})
		}

		var diag diag.Diagnostics
		data.Triggers, diag = types.SetValueFrom(ctx, triggersObjectType, triggers)
		diags = append(diags, diag...)
	}

	data.Id = types.StringValue(id.FullyQualifiedName())
	return data, nil, diags
}

func (r *ResourceMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceMonitorModelV1
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, _, diags := r.update(ctx, &plan, &state, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags.Append(resp.State.Set(ctx, &data)...)
}

func (r *ResourceMonitorResource) update(ctx context.Context, plan *resourceMonitorModelV1, state *resourceMonitorModelV1, dryRun bool) (*resourceMonitorModelV1, []string, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	client := r.client
	if dryRun {
		client = sdk.NewDryRunClient()
	}
	id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(state.Id.ValueString())
	opts := &sdk.AlterResourceMonitorOptions{}
	runUpdate := false
	if !plan.CreditQuota.Equal(state.CreditQuota) {
		runUpdate = true
		if opts.Set == nil {
			opts.Set = &sdk.ResourceMonitorSet{}
		}
		opts.Set.CreditQuota = sdk.Int(int(plan.CreditQuota.ValueFloat64()))
	}
	if !plan.Frequency.Equal(state.Frequency) {
		runUpdate = true
		if opts.Set == nil {
			opts.Set = &sdk.ResourceMonitorSet{}
		}
		frequency, err := sdk.ToResourceMonitorFrequency(plan.Frequency.ValueString())
		if err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Unable to update resource monitor, got error: %s", err))
			return plan, nil, diags
		}
		opts.Set.Frequency = frequency
		opts.Set.StartTimestamp = plan.StartTimestamp.ValueStringPointer()
	}
	if !plan.StartTimestamp.Equal(state.StartTimestamp) {
		runUpdate = true
		if opts.Set == nil {
			opts.Set = &sdk.ResourceMonitorSet{}
		}
		frequency, err := sdk.ToResourceMonitorFrequency(plan.Frequency.ValueString())
		if err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Unable to update resource monitor, got error: %s", err))
			return plan, nil, diags
		}
		opts.Set.Frequency = frequency
		opts.Set.StartTimestamp = plan.StartTimestamp.ValueStringPointer()
	}
	if !plan.EndTimestamp.Equal(state.EndTimestamp) && plan.EndTimestamp.ValueString() != "" {
		runUpdate = true
		if opts.Set == nil {
			opts.Set = &sdk.ResourceMonitorSet{}
		}
		opts.Set.EndTimestamp = plan.EndTimestamp.ValueStringPointer()
	}
	if !plan.NotifyUsers.Equal(state.NotifyUsers) {
		runUpdate = true
		var notifiedUsers []sdk.NotifiedUser
		elements := make([]types.String, 0, len(plan.NotifyUsers.Elements()))
		plan.NotifyUsers.ElementsAs(ctx, &elements, false)
		for _, e := range elements {
			notifiedUsers = append(notifiedUsers, sdk.NotifiedUser{Name: sdk.NewAccountObjectIdentifier(e.ValueString())})
		}
		opts.Set.NotifyUsers = &sdk.NotifyUsers{
			Users: notifiedUsers,
		}
	}

	if !plan.Triggers.Equal(state.Triggers) {
		runUpdate = true
		var triggers []sdk.TriggerDefinition
		elements := make([]resourceMonitorTriggerModel, 0, len(plan.Triggers.Elements()))
		plan.Triggers.ElementsAs(ctx, &elements, false)
		for _, e := range elements {
			triggers = append(triggers, sdk.TriggerDefinition{
				Threshold:     int(e.Threshold.ValueInt64()),
				TriggerAction: sdk.TriggerAction(e.TriggerAction.ValueString()),
			})
		}
		opts.Triggers = triggers
	}

	if runUpdate {
		err := client.ResourceMonitors.Alter(ctx, id, opts)
		if dryRun {
			return plan, client.TraceLogs(), diags
		}
		if err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Unable to update resource monitor, got error: %s", err))
			return plan, nil, diags
		}
	}
	data, _, readDiags := r.read(ctx, plan, false)
	diags.Append(readDiags...)
	return data, nil, diags
}

func (r *ResourceMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *resourceMonitorModelV1
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, _, diags := r.delete(ctx, data, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceMonitorResource) delete(ctx context.Context, data *resourceMonitorModelV1, dryRun bool) (*resourceMonitorModelV1, []string, diag.Diagnostics) {
	client := r.client
	if dryRun {
		client = sdk.NewDryRunClient()
	}

	diags := diag.Diagnostics{}
	id := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(data.Id.ValueString())
	err := client.ResourceMonitors.Drop(ctx, id, &sdk.DropResourceMonitorOptions{IfExists: sdk.Bool(true)})
	if dryRun {
		return data, client.TraceLogs(), diags
	}
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Unable to delete database, got error: %s", err))
		return data, nil, diags
	}
	return data, nil, diags
}

func (r *ResourceMonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
