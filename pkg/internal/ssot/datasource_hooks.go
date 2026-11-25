package ssot

import "slices"

// DataSourceHook has the states of ResourceHook.
type DataSourceHook = ResourceHook

// DataSourceHooks are the steps of a generated data source. EffectiveDataSourceHooks gives their defaults.
type DataSourceHooks struct {
	// Constructor.
	Description DataSourceHook
	ReadWrapper DataSourceHook

	// Schema.
	FiltersSchema          DataSourceHook
	SwitchesSchema         DataSourceHook
	ShowOutputSchema       DataSourceHook
	DescribeOutputSchema   DataSourceHook
	ParametersOutputSchema DataSourceHook
	AdditionalSchema       DataSourceHook

	// Read.
	ApplyFilters        DataSourceHook
	ShowInSdk           DataSourceHook
	SetId               DataSourceHook
	DescribeInSdk       DataSourceHook
	DescribeOutputSet   DataSourceHook
	ShowParametersInSdk DataSourceHook
	ParametersOutputSet DataSourceHook
	ItemExtra           DataSourceHook
}

type DataSourceHookOption func(*DataSourceHooks)

func WithDataSourceDescription(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.Description = h }
}

func WithDataSourceReadWrapper(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.ReadWrapper = h }
}

func WithDataSourceFiltersSchema(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.FiltersSchema = h }
}

func WithDataSourceSwitchesSchema(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.SwitchesSchema = h }
}

func WithDataSourceShowOutputSchema(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.ShowOutputSchema = h }
}

func WithDataSourceDescribeOutputSchema(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.DescribeOutputSchema = h }
}

func WithDataSourceParametersOutputSchema(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.ParametersOutputSchema = h }
}

func WithDataSourceAdditionalSchema(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.AdditionalSchema = h }
}

func WithDataSourceApplyFilters(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.ApplyFilters = h }
}

func WithDataSourceShowInSdk(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.ShowInSdk = h }
}

func WithDataSourceSetId(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.SetId = h }
}

func WithDataSourceDescribeInSdk(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.DescribeInSdk = h }
}

func WithDataSourceDescribeOutputSet(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.DescribeOutputSet = h }
}

func WithDataSourceShowParametersInSdk(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.ShowParametersInSdk = h }
}

func WithDataSourceParametersOutputSet(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.ParametersOutputSet = h }
}

func WithDataSourceItemExtra(h DataSourceHook) DataSourceHookOption {
	return func(hooks *DataSourceHooks) { hooks.ItemExtra = h }
}

// DataSourceInKind is the IN flavor of the SDK SHOW request.
type DataSourceInKind string

const (
	DataSourceInNone                  DataSourceInKind = ""
	DataSourceInPlain                 DataSourceInKind = "In"
	DataSourceInExtended              DataSourceInKind = "ExtendedIn"
	DataSourceInExtendedWithoutSchema DataSourceInKind = "ExtendedInWithoutSchema"
	DataSourceInAccount               DataSourceInKind = "InAccount"
)

// DataSourceFilters are the filters that the SDK SHOW request offers.
type DataSourceFilters struct {
	Like       bool
	StartsWith bool
	Limit      bool
	On         bool
	Terse      bool
	History    bool
	In         DataSourceInKind
}

func (obj Object) DataSourceFilters() DataSourceFilters {
	cfg := obj.SdkConfig()
	generator := ObjectGenerators(obj.Generators...).DataSource
	filters := DataSourceFilters{
		Like:       cfg.ShowOptionalLike,
		StartsWith: cfg.ShowOptionalStartsWith && !slices.Contains(generator.WithoutFilters, DataSourceFilterStartsWith),
		Limit:      (cfg.ShowOptionalLimitFrom || cfg.ShowOptionalLimit) && !slices.Contains(generator.WithoutFilters, DataSourceFilterLimit),
		On:         cfg.ShowOptionalOn,
		Terse:      cfg.ShowTerse && generator.WithTerse,
		History:    slices.Contains(cfg.ShowOptionalSql, "HISTORY") && generator.WithHistory,
	}
	switch {
	case cfg.ShowOptionalExtendedInWithoutSchema:
		filters.In = DataSourceInExtendedWithoutSchema
	case cfg.ShowOptionalExtendedIn:
		filters.In = DataSourceInExtended
	case cfg.ShowOptionalIn:
		filters.In = DataSourceInPlain
	case cfg.ShowOptionalInAccount:
		filters.In = DataSourceInAccount
	}
	return filters
}

// DataSourceDescribe is how the data source calls DESCRIBE.
type DataSourceDescribe struct {
	Method string
	// Slice is true when the call returns a list.
	Slice bool
	// RequestConstructor is the SDK request constructor when DESCRIBE takes a request, such as NewDescribeColumnsTableRequest.
	// Empty means the call takes item.ID().
	RequestConstructor string
	// ResultName is the SDK type of the DESCRIBE result, such as TableColumnDetails. Empty uses DescribeDetailsName or Name plus Details.
	ResultName string
}

func (obj Object) HasDataSourceDescribe() bool {
	if ObjectGenerators(obj.Generators...).DataSource.DescribeMethod != "" {
		return true
	}
	return obj.SdkConfig().Describe && len(obj.Describe.Columns) > 0
}

func (obj Object) HasDataSourceParameters() bool {
	return obj.SdkConfig().ShowParameters
}

// DataSourceDescribe prefers a custom-operation method, then DescribeDetails, then Describe.
func (obj Object) DataSourceDescribe() DataSourceDescribe {
	method := ObjectGenerators(obj.Generators...).DataSource.DescribeMethod
	if method != "" {
		op := obj.customOperation(method)
		if op == nil {
			panic("object " + obj.Name + ": DataSourceDescribeMethod " + method + " is not a custom operation")
		}
		describe := DataSourceDescribe{Method: method, RequestConstructor: "New" + method + obj.Name + "Request"}
		if op.Paired != nil {
			describe.Slice = true
			describe.ResultName = op.Paired.Plain
		}
		return describe
	}
	cfg := obj.SdkConfig()
	if cfg.DescribeDetailsMethod || cfg.DescribeDetailsFromIdentifier {
		return DataSourceDescribe{Method: "DescribeDetails"}
	}
	return DataSourceDescribe{Method: "Describe", Slice: cfg.DescribeMapping == SdkDescriptionMappingSlice}
}

func (obj Object) customOperation(kind string) *CustomOperationSpec {
	for i := range obj.CustomOperations {
		if obj.CustomOperations[i].Kind == kind {
			return &obj.CustomOperations[i]
		}
	}
	return nil
}

// EffectiveDataSourceHooks are the hooks of the data source of this object. The data source exposes what the object has,
// so a definition lists only the exceptions.
func (obj Object) EffectiveDataSourceHooks() DataSourceHooks {
	h := DataSourceHooks{
		Description:            Generated(),
		ReadWrapper:            Generated(),
		FiltersSchema:          Generated(),
		SwitchesSchema:         Generated(),
		ShowOutputSchema:       Generated(),
		DescribeOutputSchema:   Generated(),
		ParametersOutputSchema: Generated(),
		AdditionalSchema:       Omitted(),
		ApplyFilters:           Generated(),
		ShowInSdk:              Generated(),
		SetId:                  Generated(),
		DescribeInSdk:          Generated(),
		DescribeOutputSet:      Generated(),
		ShowParametersInSdk:    Generated(),
		ParametersOutputSet:    Generated(),
		ItemExtra:              Omitted(),
	}
	if !obj.HasDataSourceDescribe() {
		h.DescribeInSdk = Omitted()
	}
	if !obj.HasDataSourceParameters() {
		h.ShowParametersInSdk = Omitted()
	}
	for _, opt := range obj.DataSourceHooks {
		opt(&h)
	}
	deriveDataSourceHooks(&h)
	return h
}

// deriveDataSourceHooks drops generated steps that lost what they depend on. Explicit choices stay.
func deriveDataSourceHooks(h *DataSourceHooks) {
	if h.DescribeOutputSchema.Generate && !h.DescribeInSdk.Present {
		h.DescribeOutputSchema = Omitted()
	}
	if h.DescribeOutputSet.Generate && (!h.DescribeInSdk.Present || !h.DescribeOutputSchema.Present) {
		h.DescribeOutputSet = Omitted()
	}
	if h.ParametersOutputSchema.Generate && !h.ShowParametersInSdk.Present {
		h.ParametersOutputSchema = Omitted()
	}
	if h.ParametersOutputSet.Generate && (!h.ShowParametersInSdk.Present || !h.ParametersOutputSchema.Present) {
		h.ParametersOutputSet = Omitted()
	}
	if h.SwitchesSchema.Generate && !h.DescribeInSdk.Present && !h.ShowParametersInSdk.Present {
		h.SwitchesSchema = Omitted()
	}
}
