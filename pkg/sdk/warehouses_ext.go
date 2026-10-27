package sdk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/util"
	"github.com/hashicorp/go-uuid"
)

// Backward-compatible enum constant aliases.
// The generated names use strict CamelCase from the hyphenated/numeric values,
// while the old hand-written names used different casing conventions.
const (
	WarehouseSizeXSmall   = WarehouseSizeXsmall
	WarehouseSizeXLarge   = WarehouseSizeXlarge
	WarehouseSizeXXLarge  = WarehouseSizeXxlarge
	WarehouseSizeXXXLarge = WarehouseSizeXxxlarge
	WarehouseSizeX4Large  = WarehouseSizeX4large
	WarehouseSizeX5Large  = WarehouseSizeX5large
	WarehouseSizeX6Large  = WarehouseSizeX6large
)

const (
	WarehouseResourceConstraintMemory1X     = WarehouseResourceConstraintMemory1x
	WarehouseResourceConstraintMemory1Xx86  = WarehouseResourceConstraintMemory1xX86
	WarehouseResourceConstraintMemory16X    = WarehouseResourceConstraintMemory16x
	WarehouseResourceConstraintMemory16Xx86 = WarehouseResourceConstraintMemory16xX86
	WarehouseResourceConstraintMemory64X    = WarehouseResourceConstraintMemory64x
	WarehouseResourceConstraintMemory64Xx86 = WarehouseResourceConstraintMemory64xX86
)

const (
	MaxQueryPerformanceLevelXSmall   = MaxQueryPerformanceLevelXsmall
	MaxQueryPerformanceLevelXLarge   = MaxQueryPerformanceLevelXlarge
	MaxQueryPerformanceLevelXXLarge  = MaxQueryPerformanceLevelXxlarge
	MaxQueryPerformanceLevelXXXLarge = MaxQueryPerformanceLevelXxxlarge
	MaxQueryPerformanceLevelX4Large  = MaxQueryPerformanceLevelX4large
)

// AllWarehouseTypesUserSettable excludes ADAPTIVE and INTERACTIVE. Those types are separate resources.
var AllWarehouseTypesUserSettable = []WarehouseType{
	WarehouseTypeStandard,
	WarehouseTypeSnowparkOptimized,
}

// ToWarehouseTypeUserSettable parses a warehouse type string, excluding ADAPTIVE
// which is not settable through the standard warehouse resource.
func ToWarehouseTypeUserSettable(s string) (WarehouseType, error) {
	switch strings.ToUpper(s) {
	case string(WarehouseTypeStandard):
		return WarehouseTypeStandard, nil
	case string(WarehouseTypeSnowparkOptimized):
		return WarehouseTypeSnowparkOptimized, nil
	default:
		return "", fmt.Errorf("invalid warehouse type: %s", s)
	}
}

func IsWarehouseResourceConstraintForSnowparkOptimized(s WarehouseResourceConstraint) bool {
	return slices.Contains(AllWarehouseResourceConstraintsWithoutGenerations, string(s))
}

// FromString methods to satisfy the EnumCreator interface used by custom types.

func (e WarehouseType) FromString(s string) (WarehouseType, error) {
	return ToWarehouseType(s)
}

func (e WarehouseSize) FromString(s string) (WarehouseSize, error) {
	return ToWarehouseSize(s)
}

func (e ScalingPolicy) FromString(s string) (ScalingPolicy, error) {
	return ToScalingPolicy(s)
}

// Validation vars used by resources/datasources.

// ValidWarehouseSizesString is based on https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties
var ValidWarehouseSizesString = AllWarehouseSizesString

// ValidWarehouseScalingPoliciesString is based on https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties
var ValidWarehouseScalingPoliciesString = []string{
	string(ScalingPolicyStandard),
	string(ScalingPolicyEconomy),
}

// ValidWarehouseTypesString is based on https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties
var ValidWarehouseTypesString = []string{
	string(WarehouseTypeStandard),
	string(WarehouseTypeSnowparkOptimized),
	string(WarehouseTypeAdaptive),
}

var ValidWarehouseTypesRegularString = []string{
	string(WarehouseTypeStandard),
	string(WarehouseTypeSnowparkOptimized),
}

var AllWarehouseResourceConstraintsWithoutGenerations = []string{
	string(WarehouseResourceConstraintMemory1x),
	string(WarehouseResourceConstraintMemory1xX86),
	string(WarehouseResourceConstraintMemory16x),
	string(WarehouseResourceConstraintMemory16xX86),
	string(WarehouseResourceConstraintMemory64x),
	string(WarehouseResourceConstraintMemory64xX86),
}

// WarehouseParameters is based on https://docs.snowflake.com/en/sql-reference/parameters#object-parameters
var WarehouseParameters = []ObjectParameter{
	ObjectParameterMaxConcurrencyLevel,
	ObjectParameterStatementQueuedTimeoutInSeconds,
	ObjectParameterStatementTimeoutInSeconds,
}

// ShowByIDExperimental is a show by id function with improved performance (using starts with and limit).
func (v *warehouses) ShowByIDExperimental(ctx context.Context, id AccountObjectIdentifier) (*Warehouse, error) {
	warehouses, err := v.Show(ctx, NewShowWarehouseRequest().
		WithLike(Like{Pattern: String(id.Name())}).
		WithStartsWith(id.Name()).
		WithLimit(LimitFrom{Rows: Int(1)}))
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(warehouses, func(r Warehouse) bool { return r.Name == id.Name() })
}

func (v *warehouses) ShowByIDExperimentalSafely(ctx context.Context, id AccountObjectIdentifier) (*Warehouse, error) {
	return SafeShowById(v.client, v.ShowByIDExperimental, ctx, id)
}

func (v *warehouses) ShowParameters(ctx context.Context, id AccountObjectIdentifier) ([]*Parameter, error) {
	return v.client.Parameters.ShowParameters(ctx, &ShowParametersOptions{
		In: &ParametersIn{
			Warehouse: id,
		},
	})
}

// AlterWithSuspend suspends around changes Snowflake rejects on a running warehouse:
// STANDARD/SNOWPARK-OPTIMIZED type changes and interactive resizes. ADAPTIVE conversion is online.
func (v *warehouses) AlterWithSuspend(ctx context.Context, request *AlterWarehouseRequest) error {
	if request.Set == nil || (request.Set.WarehouseType == nil && request.Set.WarehouseSize == nil) {
		return v.Alter(ctx, request)
	}

	warehouse, err := v.ShowByID(ctx, request.name)
	if err != nil {
		return err
	}

	if alterNeedsSuspendedWarehouse(request.Set, warehouse) {
		cleanup, err := v.suspendStartedWarehouse(ctx, request.name)
		defer cleanup()
		if err != nil {
			return err
		}
	}
	return v.Alter(ctx, request)
}

// alterNeedsSuspendedWarehouse reports whether this SET requires a suspended warehouse.
func alterNeedsSuspendedWarehouse(set *WarehouseSetRequest, warehouse *Warehouse) bool {
	changesType := set.WarehouseType != nil
	becomesAdaptive := changesType && *set.WarehouseType == WarehouseTypeAdaptive
	changesSize := set.WarehouseSize != nil
	return (changesType && !becomesAdaptive) || (changesSize && warehouse.IsInteractiveWarehouse())
}

func (v *warehouses) AlterInteractiveWithSuspend(ctx context.Context, request *AlterInteractiveWarehouseRequest) error {
	changesSize := request.Set != nil && request.Set.WarehouseSize != nil
	if !changesSize {
		return v.AlterInteractive(ctx, request)
	}
	cleanup, err := v.suspendStartedWarehouse(ctx, request.name)
	defer cleanup()
	if err != nil {
		return err
	}
	return v.AlterInteractive(ctx, request)
}

// suspendStartedWarehouse suspends id if started. The returned cleanup resumes it and is never nil.
func (v *warehouses) suspendStartedWarehouse(ctx context.Context, id AccountObjectIdentifier) (func(), error) {
	noop := func() {}
	warehouse, err := v.ShowByID(ctx, id)
	if err != nil {
		return noop, err
	}
	if warehouse.State != WarehouseStateStarted {
		return noop, nil
	}
	err = v.Alter(ctx, NewAlterWarehouseRequest(id).WithSuspend(true))
	if err != nil {
		return noop, err
	}
	cleanup := func() {
		err := v.Alter(ctx, NewAlterWarehouseRequest(id).WithResume(true).WithIfSuspended(true))
		if err != nil {
			log.Printf("[DEBUG] error occurred during warehouse resumption, err=%v", err)
		}
	}
	var warehouseSuspensionErrs []error
	err = util.Retry(5, 1*time.Second, func() (error, bool) {
		warehouse, err = v.ShowByID(ctx, id)
		if err != nil {
			warehouseSuspensionErrs = append(warehouseSuspensionErrs, err)
			return nil, false
		}
		if warehouse.State != WarehouseStateSuspended {
			return nil, false
		}
		return nil, true
	})
	if err != nil {
		return cleanup, fmt.Errorf("warehouse suspension failed, err: %w, original errors: %w", err, errors.Join(warehouseSuspensionErrs...))
	}
	return cleanup, nil
}

// CreateInteractivePreservingSession wraps CreateInteractive to work around a known Snowflake behavior:
// creating a warehouse implicitly switches the session onto it, and for interactive warehouses (which
// always have a short statement timeout) this can cause any subsequent statement run in the same session
// to time out. After creation succeeds, this restores whichever warehouse was active in the session
// beforehand. If no warehouse was active, there is no "USE WAREHOUSE NONE"/unset syntax, so a temporary
// warehouse is created and immediately dropped - dropping the warehouse currently in use for a session
// clears CURRENT_WAREHOUSE() back to empty.
// Restoring the session warehouse is best-effort: the interactive warehouse has already been created
// successfully by the time this runs, so a failure here is logged rather than returned, to avoid the
// resource being reported as failed (and possibly recreated) when the object itself is fine.
func (v *warehouses) CreateInteractivePreservingSession(ctx context.Context, request *CreateInteractiveWarehouseRequest) error {
	previousWarehouse, err := v.client.ContextFunctions.CurrentWarehouse(ctx)
	if err != nil {
		return err
	}

	if err := v.CreateInteractive(ctx, request); err != nil {
		return err
	}

	if previousWarehouse.Value != "" {
		if err := v.client.Sessions.UseWarehouse(ctx, NewUseWarehouseSessionRequest(NewAccountObjectIdentifier(previousWarehouse.Value))); err != nil {
			log.Printf("[WARN] could not restore session warehouse to %s after creating interactive warehouse %s, err=%v", previousWarehouse.Value, request.ID().FullyQualifiedName(), err)
		}
		return nil
	}

	if err := v.clearSessionWarehouse(ctx); err != nil {
		log.Printf("[WARN] could not clear session warehouse after creating interactive warehouse %s, err=%v", request.ID().FullyQualifiedName(), err)
	}
	return nil
}

// clearSessionWarehouse leaves the session with no current warehouse, mirroring a session that never had
// one selected. There is no direct syntax for this, so it creates a minimal temporary warehouse (which
// implicitly becomes the session's current warehouse) and drops it again, which clears CURRENT_WAREHOUSE().
func (v *warehouses) clearSessionWarehouse(ctx context.Context) error {
	suffix, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	tempId := NewAccountObjectIdentifier(fmt.Sprintf("TF_TEMP_%s", strings.ReplaceAll(suffix, "-", "_")))

	if err := v.Create(ctx, NewCreateWarehouseRequest(tempId).WithWarehouseSize(WarehouseSizeXSmall).WithInitiallySuspended(true)); err != nil {
		return err
	}
	return v.Drop(ctx, NewDropWarehouseRequest(tempId).WithIfExists(true))
}

// additionalConvert handles manual field conversions that the generator cannot express.
func (r warehousesRow) additionalConvert(wh *Warehouse) error {
	// State and Type - simple casts from string
	wh.State = WarehouseState(r.State)
	wh.Type = WarehouseType(r.Type)

	// Float fields from string
	if available := strings.TrimSpace(r.Available); available != "" {
		if val, err := strconv.ParseFloat(available, 64); err != nil {
			return fmt.Errorf(`row 'available' has incorrect value '%s', %w`, available, err)
		} else {
			wh.Available = val
		}
	}
	if provisioning := strings.TrimSpace(r.Provisioning); provisioning != "" {
		if val, err := strconv.ParseFloat(provisioning, 64); err != nil {
			return fmt.Errorf(`row 'provisioning' has incorrect value '%s', %w`, provisioning, err)
		} else {
			wh.Provisioning = val
		}
	}
	if quiescing := strings.TrimSpace(r.Quiescing); quiescing != "" {
		if val, err := strconv.ParseFloat(quiescing, 64); err != nil {
			return fmt.Errorf(`row 'quiescing' has incorrect value '%s', %w`, quiescing, err)
		} else {
			wh.Quiescing = val
		}
	}
	if other := strings.TrimSpace(r.Other); other != "" {
		if val, err := strconv.ParseFloat(other, 64); err != nil {
			return fmt.Errorf(`row 'other' has incorrect value '%s', %w`, other, err)
		} else {
			wh.Other = val
		}
	}

	// Generation
	if r.Generation.Valid {
		generation, err := ToWarehouseGeneration(r.Generation.String)
		if err != nil {
			return err
		}
		wh.Generation = &generation
	}

	// ResourceConstraint - conditional on warehouse type.
	// We use EqualFold instead of the generated ToWarehouseResourceConstraint because
	// the values contain lowercase "x86" and the generated function uses ToUpper which breaks matching.
	if r.ResourceConstraint.Valid && wh.Type == WarehouseTypeSnowparkOptimized {
		var found bool
		for _, rc := range AllWarehouseResourceConstraints {
			if strings.EqualFold(string(rc), r.ResourceConstraint.String) {
				v := rc
				wh.ResourceConstraint = &v
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid resource constraint: %s", r.ResourceConstraint.String)
		}
	}
	return nil
}

// additionalValidations for CreateWarehouseOptions.
func (opts *CreateWarehouseOptions) additionalValidations() error {
	var errs []error
	if valueSet(opts.MinClusterCount) && valueSet(opts.MaxClusterCount) && !validateIntGreaterThanOrEqual(*opts.MaxClusterCount, *opts.MinClusterCount) {
		errs = append(errs, fmt.Errorf("MinClusterCount must be less than or equal to MaxClusterCount"))
	}
	if valueSet(opts.QueryAccelerationMaxScaleFactor) && !validateIntInRangeInclusive(*opts.QueryAccelerationMaxScaleFactor, 0, 100) {
		errs = append(errs, errIntBetween("CreateWarehouseOptions", "QueryAccelerationMaxScaleFactor", 0, 100))
	}
	if valueSet(opts.WarehouseType) && !slices.Contains(ValidWarehouseTypesRegularString, string(*opts.WarehouseType)) {
		errs = append(errs, fmt.Errorf("only %s warehouses are supported, got %s", collections.JoinStrings(ValidWarehouseTypesRegularString, ", "), *opts.WarehouseType))
	}
	return JoinErrors(errs...)
}

// additionalValidations for CreateAdaptiveWarehouseOptions.
func (opts *CreateAdaptiveWarehouseOptions) additionalValidations() error {
	var errs []error
	if valueSet(opts.QueryThroughputMultiplier) && !validateIntGreaterThanOrEqual(*opts.QueryThroughputMultiplier, 0) {
		errs = append(errs, fmt.Errorf("QueryThroughputMultiplier must be greater than or equal to 0"))
	}
	return JoinErrors(errs...)
}

// additionalValidations for CreateInteractiveWarehouseOptions.
func (opts *CreateInteractiveWarehouseOptions) additionalValidations() error {
	var errs []error
	if valueSet(opts.MinClusterCount) && valueSet(opts.MaxClusterCount) && !validateIntGreaterThanOrEqual(*opts.MaxClusterCount, *opts.MinClusterCount) {
		errs = append(errs, fmt.Errorf("MinClusterCount must be less than or equal to MaxClusterCount"))
	}
	return JoinErrors(errs...)
}

// additionalValidations for AlterWarehouseOptions.
func (opts *AlterWarehouseOptions) additionalValidations() error {
	var errs []error
	if everyValueSet(opts.Suspend, opts.Resume) && (*opts.Suspend && *opts.Resume) {
		errs = append(errs, errOneOf("AlterWarehouseOptions", "Suspend", "Resume"))
	}
	if (valueSet(opts.IfSuspended) && *opts.IfSuspended) && (!valueSet(opts.Resume) || !*opts.Resume) {
		errs = append(errs, fmt.Errorf(`"Resume" has to be set when using "IfSuspended"`))
	}
	return JoinErrors(errs...)
}

// additionalValidations for WarehouseSet.
func (v *WarehouseSet) additionalValidations() error {
	var errs []error
	// we validate only the case when both are set together, if only MinClusterCount is set, we leave it for Snowflake to validate
	if v.MinClusterCount != nil && valueSet(v.MaxClusterCount) {
		if ok := validateIntInRangeInclusive(*v.MinClusterCount, 1, *v.MaxClusterCount); !ok {
			errs = append(errs, fmt.Errorf("MinClusterCount must be less than or equal to MaxClusterCount"))
		}
	}
	if v.AutoSuspend != nil {
		if ok := validateIntGreaterThanOrEqual(*v.AutoSuspend, 0); !ok {
			errs = append(errs, fmt.Errorf("AutoSuspend must be greater than or equal to 0"))
		}
	}
	if v.QueryAccelerationMaxScaleFactor != nil {
		if ok := validateIntInRangeInclusive(*v.QueryAccelerationMaxScaleFactor, 0, 100); !ok {
			errs = append(errs, fmt.Errorf("QueryAccelerationMaxScaleFactor must be between 0 and 100"))
		}
	}
	return JoinErrors(errs...)
}

func (v *WarehouseAdaptiveSet) additionalValidations() error {
	var errs []error
	if v.QueryThroughputMultiplier != nil {
		if ok := validateIntGreaterThanOrEqual(*v.QueryThroughputMultiplier, 0); !ok {
			errs = append(errs, fmt.Errorf("QueryThroughputMultiplier must be greater than or equal to 0"))
		}
	}
	return JoinErrors(errs...)
}

func (v *WarehouseInteractiveSet) additionalValidations() error {
	var errs []error
	if v.AutoSuspend != nil {
		if ok := validateIntGreaterThanOrEqual(*v.AutoSuspend, 0); !ok {
			errs = append(errs, fmt.Errorf("AutoSuspend must be greater than or equal to 0"))
		}
	}
	return JoinErrors(errs...)
}

func (s *CreateWarehouseRequest) ID() AccountObjectIdentifier {
	return s.name
}

func (s *CreateAdaptiveWarehouseRequest) ID() AccountObjectIdentifier {
	return s.name
}

func (s *CreateInteractiveWarehouseRequest) ID() AccountObjectIdentifier {
	return s.name
}

// IsInteractiveWarehouse reports whether the warehouse is interactive. Snowflake surfaces this
// through the type column (type = INTERACTIVE); there is no separate is_interactive column.
func (w *Warehouse) IsInteractiveWarehouse() bool {
	return w.Type == WarehouseTypeInteractive
}

func (v *WarehouseRegular) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *WarehouseDetails) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *WarehouseAdaptive) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *WarehouseInteractive) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}
