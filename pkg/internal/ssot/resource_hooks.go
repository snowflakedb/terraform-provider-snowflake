package ssot

type ResourceHook struct {
	Present  bool
	Generate bool
}

type ResourceHooks struct {
	ParseIdFromConfig               ResourceHook
	ParseId                         ResourceHook
	BeforeCreate                    ResourceHook
	NewCreateRequest                ResourceHook
	ApplyCreateOptionals            ResourceHook
	CreateInSdk                     ResourceHook
	AfterCreate                     ResourceHook
	ShowByIdSafelyInSdk             ResourceHook
	DescribeInSdk                   ResourceHook
	HandleExternalChanges           ResourceHook
	HandleExternalChangesInDescribe ResourceHook
	SetStateToValuesFromConfig      ResourceHook
	ShowOutputSet                   ResourceHook
	DescribeOutputSet               ResourceHook
	SetDescribeOutput               ResourceHook
	FqnSet                          ResourceHook
	SetConfigFields                 ResourceHook
	HierarchyRename                 ResourceHook
	BeforeAlter                     ResourceHook
	ApplySetUnset                   ResourceHook
	AlterSetInSdk                   ResourceHook
	AlterUnsetInSdk                 ResourceHook
	AfterAlter                      ResourceHook
	SetFieldsNotSetByRead           ResourceHook
	ShowByIdInSdk                   ResourceHook
	DropSafelyInSdk                 ResourceHook
	StateUpgraders                  ResourceHook
	IdentitySchema                  ResourceHook
	AttributesSchema                ResourceHook
	FqnSchema                       ResourceHook
	ShowOutputSchema                ResourceHook
	DescribeOutputSchema            ResourceHook
	ApplyParametersCreate           ResourceHook
	ShowParametersInSdk             ResourceHook
	ParametersDetailsFromRaw        ResourceHook
	SetParametersFields             ResourceHook
	ParametersOutputSet             ResourceHook
	ApplyParametersChanges          ResourceHook
	ParametersAttributesSchema      ResourceHook
	ParametersOutputSchema          ResourceHook
	ParametersFieldNames            ResourceHook
	ParametersCustomDiff            ResourceHook
	CustomizeDiff                   ResourceHook
}

type ResourceHookOption func(*ResourceHooks)

func ExtCall() ResourceHook {
	return ResourceHook{Present: true, Generate: false}
}

func Generated() ResourceHook {
	return ResourceHook{Present: true, Generate: true}
}

func Omitted() ResourceHook {
	return ResourceHook{}
}

func DefaultResourceExtHooks(opts ...ResourceHookOption) ResourceHooks {
	h := ResourceHooks{
		ParseIdFromConfig:               ExtCall(),
		ParseId:                         ExtCall(),
		BeforeCreate:                    Omitted(),
		NewCreateRequest:                ExtCall(),
		ApplyCreateOptionals:            ExtCall(),
		CreateInSdk:                     ExtCall(),
		AfterCreate:                     Omitted(),
		ShowByIdSafelyInSdk:             ExtCall(),
		DescribeInSdk:                   ExtCall(),
		HandleExternalChanges:           ExtCall(),
		HandleExternalChangesInDescribe: Omitted(),
		SetStateToValuesFromConfig:      ExtCall(),
		ShowOutputSet:                   Generated(),
		DescribeOutputSet:               Generated(),
		SetDescribeOutput:               Omitted(),
		FqnSet:                          Generated(),
		SetConfigFields:                 ExtCall(),
		HierarchyRename:                 Omitted(),
		BeforeAlter:                     Omitted(),
		ApplySetUnset:                   ExtCall(),
		AlterSetInSdk:                   ExtCall(),
		AlterUnsetInSdk:                 ExtCall(),
		AfterAlter:                      Omitted(),
		SetFieldsNotSetByRead:           ExtCall(),
		ShowByIdInSdk:                   ExtCall(),
		DropSafelyInSdk:                 Omitted(),
		StateUpgraders:                  Omitted(),
		IdentitySchema:                  Generated(),
		AttributesSchema:                Generated(),
		FqnSchema:                       Generated(),
		ShowOutputSchema:                Generated(),
		DescribeOutputSchema:            Generated(),
		ApplyParametersCreate:           Omitted(),
		ShowParametersInSdk:             Omitted(),
		ParametersDetailsFromRaw:        Omitted(),
		SetParametersFields:             Omitted(),
		ParametersOutputSet:             Omitted(),
		ApplyParametersChanges:          Omitted(),
		ParametersAttributesSchema:      Omitted(),
		ParametersOutputSchema:          Omitted(),
		ParametersFieldNames:            Omitted(),
		ParametersCustomDiff:            Omitted(),
		CustomizeDiff:                   ExtCall(),
	}
	for _, opt := range opts {
		opt(&h)
	}
	deriveOutputSetHooks(&h)
	return h
}

// deriveOutputSetHooks drops a generated output set call when the state it fills does not exist.
// It only turns Generated into Omitted, so an explicit ExtCall or Omitted stays as it was set.
func deriveOutputSetHooks(h *ResourceHooks) {
	// show_output is only in the schema when the show output schema is present.
	if h.ShowOutputSet.Generate && !h.ShowOutputSchema.Present {
		h.ShowOutputSet = Omitted()
	}
	// describe_output needs the describe step and its schema.
	if h.DescribeOutputSet.Generate && (!h.DescribeInSdk.Present || !h.DescribeOutputSchema.Present) {
		h.DescribeOutputSet = Omitted()
	}
}

func WithParseIdFromConfig(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ParseIdFromConfig = h }
}

func WithParseId(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ParseId = h }
}

func WithBeforeCreate(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.BeforeCreate = h }
}

func WithNewCreateRequest(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.NewCreateRequest = h }
}

func WithApplyCreateOptionals(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ApplyCreateOptionals = h }
}

func WithCreateInSdk(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.CreateInSdk = h }
}

func WithAfterCreate(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.AfterCreate = h }
}

func WithShowByIdSafelyInSdk(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ShowByIdSafelyInSdk = h }
}

func WithDescribeInSdk(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.DescribeInSdk = h }
}

func WithHandleExternalChanges(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.HandleExternalChanges = h }
}

func WithHandleExternalChangesInDescribe(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.HandleExternalChangesInDescribe = h }
}

func WithSetStateToValuesFromConfig(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.SetStateToValuesFromConfig = h }
}

func WithShowOutputSet(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ShowOutputSet = h }
}

func WithDescribeOutputSet(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.DescribeOutputSet = h }
}

func WithSetDescribeOutput(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.SetDescribeOutput = h }
}

func WithFqnSet(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.FqnSet = h }
}

func WithSetConfigFields(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.SetConfigFields = h }
}

func WithHierarchyRename(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.HierarchyRename = h }
}

func WithBeforeAlter(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.BeforeAlter = h }
}

func WithApplySetUnset(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ApplySetUnset = h }
}

func WithAlterSetInSdk(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.AlterSetInSdk = h }
}

func WithAlterUnsetInSdk(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.AlterUnsetInSdk = h }
}

func WithAfterAlter(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.AfterAlter = h }
}

func WithSetFieldsNotSetByRead(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.SetFieldsNotSetByRead = h }
}

func WithShowByIdInSdk(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ShowByIdInSdk = h }
}

func WithDropSafelyInSdk(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.DropSafelyInSdk = h }
}

func WithStateUpgraders(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.StateUpgraders = h }
}

func WithIdentitySchema(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.IdentitySchema = h }
}

func WithAttributesSchema(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.AttributesSchema = h }
}

func WithFqnSchema(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.FqnSchema = h }
}

func WithShowOutputSchema(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ShowOutputSchema = h }
}

func WithDescribeOutputSchema(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.DescribeOutputSchema = h }
}

func WithApplyParametersCreate(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ApplyParametersCreate = h }
}

func WithShowParametersInSdk(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ShowParametersInSdk = h }
}

func WithParametersDetailsFromRaw(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ParametersDetailsFromRaw = h }
}

func WithSetParametersFields(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.SetParametersFields = h }
}

func WithParametersOutputSet(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ParametersOutputSet = h }
}

func WithApplyParametersChanges(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ApplyParametersChanges = h }
}

func WithParametersAttributesSchema(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ParametersAttributesSchema = h }
}

func WithParametersOutputSchema(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ParametersOutputSchema = h }
}

func WithParametersFieldNames(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ParametersFieldNames = h }
}

func WithParametersCustomDiff(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.ParametersCustomDiff = h }
}

func WithCustomizeDiff(h ResourceHook) ResourceHookOption {
	return func(hooks *ResourceHooks) { hooks.CustomizeDiff = h }
}
