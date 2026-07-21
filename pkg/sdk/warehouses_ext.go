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

// WarehouseGeneration is kept manual because its values ("1", "2") produce non-descriptive generated names.
type WarehouseGeneration string

const (
	WarehouseGenerationStandardGen1 WarehouseGeneration = "1"
	WarehouseGenerationStandardGen2 WarehouseGeneration = "2"
)

var AllWarehouseGenerations = []string{
	string(WarehouseGenerationStandardGen1),
	string(WarehouseGenerationStandardGen2),
}

func ToWarehouseGeneration(s string) (WarehouseGeneration, error) {
	switch s {
	case "1":
		return WarehouseGenerationStandardGen1, nil
	case "2":
		return WarehouseGenerationStandardGen2, nil
	default:
		return "", fmt.Errorf("invalid generation: %s", s)
	}
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

// AlterWithSuspend wraps Alter with automatic suspend/resume for changes that Snowflake refuses to
// apply to a running warehouse. Changing the warehouse type requires a suspended warehouse for any
// warehouse, and an interactive warehouse additionally cannot be resized while running. In both
// cases the warehouse is suspended before the alter and resumed afterwards.
func (v *warehouses) AlterWithSuspend(ctx context.Context, request *AlterWarehouseRequest) error {
	changesType := request.Set != nil && request.Set.WarehouseType != nil
	changesSize := request.Set != nil && request.Set.WarehouseSize != nil
	if !changesType && !changesSize {
		return v.Alter(ctx, request)
	}

	warehouse, err := v.ShowByID(ctx, request.name)
	if err != nil {
		return err
	}

	// A type change always needs a suspended warehouse; a size change needs it only for interactive
	// warehouses (regular warehouses resize live).
	mustSuspend := changesType || (changesSize && warehouse.IsInteractiveWarehouse())
	if mustSuspend && warehouse.State == WarehouseStateStarted {
		err := v.Alter(ctx, NewAlterWarehouseRequest(request.name).WithSuspend(true))
		if err != nil {
			return err
		}
		defer func() {
			err := v.Alter(ctx, NewAlterWarehouseRequest(request.name).WithResume(true).WithIfSuspended(true))
			if err != nil {
				log.Printf("[DEBUG] error occurred during warehouse resumption, err=%v", err)
			}
		}()

		// needed to make sure that warehouse is suspended
		var warehouseSuspensionErrs []error
		err = util.Retry(5, 1*time.Second, func() (error, bool) {
			warehouse, err = v.ShowByID(ctx, request.name)
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
			return fmt.Errorf("warehouse suspension failed, err: %w, original errors: %w", err, errors.Join(warehouseSuspensionErrs...))
		}
	}
	return v.Alter(ctx, request)
}

// additionalConvert handles manual field conversions that the generator cannot express.
func (r warehouseDBRow) additionalConvert(wh *Warehouse) error {
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
	if r.ResourceConstraint.Valid {
		switch wh.Type {
		case WarehouseTypeStandard:
			// After BCR 2026_02, resource_constraint is NULL for Standard warehouses; generation column is used instead.
		case WarehouseTypeSnowparkOptimized:
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
		case WarehouseTypeAdaptive:
			// Adaptive warehouses don't use resource constraints; ignore.
		default:
			return fmt.Errorf("invalid warehouse type: %s", wh.Type)
		}
	}

	// Tables - only present for interactive warehouses; may be NULL. SHOW WAREHOUSES returns the
	// associated tables as a comma-separated list of fully-qualified names. Identifiers can contain
	// commas when quoted, so we split in a quote-aware manner rather than using strings.Split.
	if r.Tables.Valid {
		if tables := strings.TrimSpace(r.Tables.String); tables != "" {
			for _, raw := range splitCommaSeparatedIdentifiers(tables) {
				raw = strings.TrimSpace(raw)
				if raw == "" {
					continue
				}
				id, err := ParseSchemaObjectIdentifier(raw)
				if err != nil {
					return fmt.Errorf("parsing table identifier %q: %w", raw, err)
				}
				wh.Tables = append(wh.Tables, id)
			}
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
	if v.QueryThroughputMultiplier != nil {
		if ok := validateIntGreaterThanOrEqual(*v.QueryThroughputMultiplier, 0); !ok {
			errs = append(errs, fmt.Errorf("QueryThroughputMultiplier must be greater than or equal to 0"))
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
