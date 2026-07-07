package sdk

import (
	"context"
	"fmt"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

// TagOnConflictAllowedValuesSequence is the value returned by SHOW TAGS in the on_conflict column
// when the tag propagation on-conflict strategy is set to ALLOWED_VALUES_SEQUENCE.
const TagOnConflictAllowedValuesSequence = "ALLOWED_VALUES_SEQUENCE"

func (r *CreateTagRequest) GetName() SchemaObjectIdentifier {
	return r.name
}

// additionalConvert handles fields that require custom parsing not expressible via struct pair options.
// Specifically, allowed_values in SHOW TAGS output uses quoted values that must be stripped.
// TODO [next PRs]: generator limitation — AssignmentKindNullableToStringArray hardcodes trimQuotes=false in the
// ParseCommaSeparatedStringArray call. Snowflake returns allowed_values with double-quoted items
// ("value1","value2"), so trimQuotes=true is needed here. This can be improved in the generator by
// adding a functional option (e.g. WithTrimQuotes()) to OptionalPlainField.
func (r tagRow) additionalConvert(result *Tag) error {
	if r.AllowedValues.Valid {
		result.AllowedValues = ParseCommaSeparatedStringArray(r.AllowedValues.String, true)
	}
	return nil
}

func (opts *AllowedValues) additionalValidations() error {
	if !validateIntInRangeInclusive(len(opts.Values), 1, 300) {
		return errIntBetween("AllowedValues", "Values", 1, 300)
	}
	return nil
}

func (opts *TagSet) additionalValidations() error {
	var errs []error
	if valueSet(opts.MaskingPolicies) && anyValueSet(opts.AllowedValues, opts.Propagate, opts.Comment) {
		errs = append(errs, errOneOf("TagSet", "MaskingPolicies", "AllowedValues", "Propagate", "Comment"))
	}
	if valueSet(opts.MaskingPolicies) {
		if !validateIntGreaterThan(len(opts.MaskingPolicies.MaskingPolicies), 0) {
			errs = append(errs, errIntValue("TagSet.MaskingPolicies", "MaskingPolicies", IntErrGreater, 0))
		}
	}
	return JoinErrors(errs...)
}

func (opts *TagUnset) additionalValidations() error {
	var errs []error
	if valueSet(opts.MaskingPolicies) {
		if !validateIntGreaterThan(len(opts.MaskingPolicies.MaskingPolicies), 0) {
			errs = append(errs, errIntValue("TagUnset.MaskingPolicies", "MaskingPolicies", IntErrGreater, 0))
		}
	}
	return JoinErrors(errs...)
}

func (opts *ShowTagOptions) additionalValidations() error {
	var errs []error
	if valueSet(opts.Like) && !valueSet(opts.Like.Pattern) {
		errs = append(errs, ErrPatternRequiredForLikeKeyword)
	}
	if valueSet(opts.In) && !exactlyOneValueSet(opts.In.Account, opts.In.Database, opts.In.Schema) {
		errs = append(errs, errExactlyOneOf("ShowTagOptions.In", "Account", "Database", "Schema"))
	}
	return JoinErrors(errs...)
}

func (opts *SetTagOptions) additionalValidations() error {
	if !canBeAssociatedWithTag(opts.objectType) {
		return fmt.Errorf("tagging for object type %s is not supported", opts.objectType)
	}
	return nil
}

func (opts *UnsetTagOptions) additionalValidations() error {
	if !canBeAssociatedWithTag(opts.objectType) {
		return fmt.Errorf("tagging for object type %s is not supported", opts.objectType)
	}
	return nil
}

// adjust applies request normalization before validation and execution.
func (r *SetTagRequest) adjust() {
	normalizeTagObjectType(&r.objectType)
	normalizeTagColumnIdentifier(&r.objectType, &r.objectName, &r.Column)
	normalizeTagAccountIdentifier(r.objectType, &r.objectName)
}

// adjust applies request normalization before validation and execution.
func (r *UnsetTagRequest) adjust() {
	normalizeTagObjectType(&r.objectType)
	normalizeTagColumnIdentifier(&r.objectType, &r.objectName, &r.Column)
	normalizeTagAccountIdentifier(r.objectType, &r.objectName)
}

// normalizeTagObjectType maps object types to the types Snowflake expects in ALTER ... SET/UNSET TAG.
func normalizeTagObjectType(objectType *ObjectType) {
	if slices.Contains([]ObjectType{ObjectTypeView, ObjectTypeMaterializedView, ObjectTypeExternalTable, ObjectTypeEventTable}, *objectType) {
		*objectType = ObjectTypeTable
	}
	if slices.Contains([]ObjectType{ObjectTypeExternalFunction}, *objectType) {
		*objectType = ObjectTypeFunction
	}
}

// normalizeTagColumnIdentifier splits a column identifier into its table-level components.
// TODO [SNOW-1022645]: discuss how we handle situation like this in the SDK.
func normalizeTagColumnIdentifier(objectType *ObjectType, objectName *ObjectIdentifier, column **string) {
	if slices.Contains([]ObjectType{ObjectTypeColumn, ObjectTypeIcebergTableColumn}, *objectType) {
		if id, ok := (*objectName).(TableColumnIdentifier); ok {
			if *objectType == ObjectTypeColumn {
				*objectType = ObjectTypeTable
			} else {
				*objectType = ObjectTypeIcebergTable
			}
			*objectName = id.SchemaObjectId()
			*column = Pointer(id.Name())
		}
	}
}

// normalizeTagAccountIdentifier strips the org name prefix from account identifiers.
// TODO(SNOW-1818976): Remove this workaround when Snowflake fixes ALTER "ORGNAME"."ACCOUNTNAME" SET/UNSET TAG.
func normalizeTagAccountIdentifier(objectType ObjectType, objectName *ObjectIdentifier) {
	if objectType == ObjectTypeAccount {
		if id, ok := (*objectName).(AccountIdentifier); ok {
			*objectName = NewAccountIdentifierFromFullyQualifiedName(id.AccountName())
		}
	}
}

// UnsetSafely removes tags from a Snowflake object, ignoring "object does not exist" errors.
func (v *tags) UnsetSafely(ctx context.Context, request *UnsetTagRequest) error {
	return SafeUnsetTag(func() error {
		return v.Unset(ctx, request.WithIfExists(true))
	})
}

// SetOnCurrentAccount applies tags to the current account.
func (v *tags) SetOnCurrentAccount(ctx context.Context, setTags []TagAssociation) error {
	return v.client.Accounts.Alter(ctx, &AlterAccountOptions{
		SetTag: setTags,
	})
}

// UnsetOnCurrentAccount removes tags from the current account.
func (v *tags) UnsetOnCurrentAccount(ctx context.Context, unsetTags []ObjectIdentifier) error {
	return v.client.Accounts.Alter(ctx, &AlterAccountOptions{
		UnsetTag: unsetTags,
	})
}

func NewAllowedValuesRequestFromStrings(values []string) *AllowedValuesRequest {
	return NewAllowedValuesRequest().WithValues(collections.Map(values, func(v string) StringAllowEmpty { return StringAllowEmpty{Value: v} }))
}

func NewTagSetMaskingPoliciesRequestWithIds(ids []SchemaObjectIdentifier) *TagSetMaskingPoliciesRequest {
	return NewTagSetMaskingPoliciesRequest().WithMaskingPolicies(NewTagMaskingPolicyRequestFromIds(ids))
}

func NewTagUnsetMaskingPoliciesRequestWithIds(ids []SchemaObjectIdentifier) *TagUnsetMaskingPoliciesRequest {
	return NewTagUnsetMaskingPoliciesRequest().WithMaskingPolicies(NewTagMaskingPolicyRequestFromIds(ids))
}

func NewTagMaskingPolicyRequestFromIds(ids []SchemaObjectIdentifier) []TagMaskingPolicyRequest {
	return collections.Map(ids, func(id SchemaObjectIdentifier) TagMaskingPolicyRequest { return *NewTagMaskingPolicyRequest(id) })
}
