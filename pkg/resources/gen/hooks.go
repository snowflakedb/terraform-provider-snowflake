package gen

type Hook struct {
	Present  bool
	Generate bool
}

type ResourceHooks struct {
	ParseIdFromConfig          Hook
	ParseId                    Hook
	BeforeCreate               Hook
	NewCreateRequest           Hook
	ApplyCreateOptionals       Hook
	CreateInSdk                Hook
	ShowByIdSafelyInSdk        Hook
	DescribeInSdk              Hook
	HandleExternalChanges      Hook
	SetStateToValuesFromConfig Hook
	ShowOutputSet              Hook
	DescribeOutputSet          Hook
	SetDescribeOutput          Hook
	FqnSet                     Hook
	SetConfigFields            Hook
	HierarchyRename            Hook
	BeforeAlter                Hook
	ApplySetUnset              Hook
	AlterSetInSdk              Hook
	AlterUnsetInSdk            Hook
	SetFieldsNotSetByRead      Hook
	ShowByIdInSdk              Hook
	SchemaVersion              Hook
	IdentitySchema             Hook
	AttributesSchema           Hook
	FqnSchema                  Hook
	ShowOutputSchema           Hook
	DescribeOutputSchema       Hook
	ApplyParametersCreate      Hook
	ShowParametersInSdk        Hook
	ParametersDetailsFromRaw   Hook
	SetParametersFields        Hook
	ParametersOutputSet        Hook
	ApplyParametersChanges     Hook
	ParametersAttributesSchema Hook
	ParametersOutputSchema     Hook
	ParametersFieldNames       Hook
}

type HookOption func(*ResourceHooks)

func extCall() Hook {
	return Hook{Present: true, Generate: false}
}

func generated() Hook {
	return Hook{Present: true, Generate: true}
}

func omitted() Hook {
	return Hook{}
}

func defaultExtHooks(opts ...HookOption) ResourceHooks {
	h := ResourceHooks{
		ParseIdFromConfig:          extCall(),
		ParseId:                    extCall(),
		BeforeCreate:               omitted(),
		NewCreateRequest:           extCall(),
		ApplyCreateOptionals:       extCall(),
		CreateInSdk:                extCall(),
		ShowByIdSafelyInSdk:        extCall(),
		DescribeInSdk:              extCall(),
		HandleExternalChanges:      extCall(),
		SetStateToValuesFromConfig: extCall(),
		ShowOutputSet:              extCall(),
		DescribeOutputSet:          extCall(),
		SetDescribeOutput:          omitted(),
		FqnSet:                     extCall(),
		SetConfigFields:            extCall(),
		HierarchyRename:            omitted(),
		BeforeAlter:                omitted(),
		ApplySetUnset:              extCall(),
		AlterSetInSdk:              extCall(),
		AlterUnsetInSdk:            extCall(),
		SetFieldsNotSetByRead:      extCall(),
		ShowByIdInSdk:              extCall(),
		SchemaVersion:              omitted(),
		IdentitySchema:             generated(),
		AttributesSchema:           extCall(),
		FqnSchema:                  generated(),
		ShowOutputSchema:           generated(),
		DescribeOutputSchema:       generated(),
		ApplyParametersCreate:      omitted(),
		ShowParametersInSdk:        omitted(),
		ParametersDetailsFromRaw:   omitted(),
		SetParametersFields:        omitted(),
		ParametersOutputSet:        omitted(),
		ApplyParametersChanges:     omitted(),
		ParametersAttributesSchema: omitted(),
		ParametersOutputSchema:     omitted(),
		ParametersFieldNames:       omitted(),
	}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func WithParseIdFromConfig(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ParseIdFromConfig = h }
}

func WithParseId(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ParseId = h }
}

func WithBeforeCreate(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.BeforeCreate = h }
}

func WithNewCreateRequest(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.NewCreateRequest = h }
}

func WithApplyCreateOptionals(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ApplyCreateOptionals = h }
}

func WithCreateInSdk(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.CreateInSdk = h }
}

func WithShowByIdSafelyInSdk(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ShowByIdSafelyInSdk = h }
}

func WithDescribeInSdk(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.DescribeInSdk = h }
}

func WithHandleExternalChanges(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.HandleExternalChanges = h }
}

func WithSetStateToValuesFromConfig(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.SetStateToValuesFromConfig = h }
}

func WithShowOutputSet(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ShowOutputSet = h }
}

func WithDescribeOutputSet(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.DescribeOutputSet = h }
}

func WithSetDescribeOutput(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.SetDescribeOutput = h }
}

func WithFqnSet(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.FqnSet = h }
}

func WithSetConfigFields(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.SetConfigFields = h }
}

func WithHierarchyRename(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.HierarchyRename = h }
}

func WithBeforeAlter(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.BeforeAlter = h }
}

func WithApplySetUnset(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ApplySetUnset = h }
}

func WithAlterSetInSdk(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.AlterSetInSdk = h }
}

func WithAlterUnsetInSdk(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.AlterUnsetInSdk = h }
}

func WithSetFieldsNotSetByRead(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.SetFieldsNotSetByRead = h }
}

func WithShowByIdInSdk(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ShowByIdInSdk = h }
}

func WithSchemaVersion(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.SchemaVersion = h }
}

func WithIdentitySchema(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.IdentitySchema = h }
}

func WithAttributesSchema(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.AttributesSchema = h }
}

func WithFqnSchema(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.FqnSchema = h }
}

func WithShowOutputSchema(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ShowOutputSchema = h }
}

func WithDescribeOutputSchema(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.DescribeOutputSchema = h }
}

func WithApplyParametersCreate(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ApplyParametersCreate = h }
}

func WithShowParametersInSdk(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ShowParametersInSdk = h }
}

func WithParametersDetailsFromRaw(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ParametersDetailsFromRaw = h }
}

func WithSetParametersFields(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.SetParametersFields = h }
}

func WithParametersOutputSet(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ParametersOutputSet = h }
}

func WithApplyParametersChanges(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ApplyParametersChanges = h }
}

func WithParametersAttributesSchema(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ParametersAttributesSchema = h }
}

func WithParametersOutputSchema(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ParametersOutputSchema = h }
}

func WithParametersFieldNames(h Hook) HookOption {
	return func(hooks *ResourceHooks) { hooks.ParametersFieldNames = h }
}
