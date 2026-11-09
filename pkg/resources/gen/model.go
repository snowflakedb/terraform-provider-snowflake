package gen

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
)

type ResourceModel struct {
	Name                                string
	NamePlural                          string
	DescribeFailRead                    bool
	DescribeUsesShowId                  bool
	HasPartialOnUpdate                  bool
	AdditionalDescribeOutputDescription string
	DescribeDetailsName                 string
	IdType                              string
	LegacyIdEncoding                    bool
	PreviewFeature                      bool
	IdentityNameForceNew                bool
	IdentityHierarchyForceNew           bool
	AdditionalIdentityNameDescription   string
	ShowColumns                         []*ssot.Attribute
	DescribeColumns                     []*ssot.Attribute
	SdkCreateRequestObjectName          string
	VariantOf                           *ssot.Object
	Attributes                          []ResourceAttribute
	// ComputedAttributes are the SHOW-sourced computed-only schema fields. ComputedColumns are their SHOW columns.
	ComputedAttributes []ResourceAttribute
	ComputedColumns    []*ssot.Attribute
	// ShowSchemaOwner is the object whose generated SHOW schema and mapper this resource uses. A variant that reuses its base uses the base.
	ShowSchemaOwner          string
	AttributesModification   bool
	ObjectType               string
	Hooks                    ssot.ResourceHooks
	SchemaStateUpgraders     []ssot.StateUpgraderFunc
	AdditionalCustomizeDiffs []ssot.CustomizeDiffFunc
	BeforeCreateSteps        BeforeCreateSteps
	AfterCreateSteps         AfterCreateSteps
	AlterStatements          AlterStatementsModel
	*genhelpers.PreambleModel
}

// AfterCreateSteps are the data of the generated AfterCreate. Each kind of step has its own list here and its own block
// in after_create.tmpl. To set another kind of attribute after create, add a list and a block.
type AfterCreateSteps struct {
	// AlterKeywords call an ALTER keyword when a bool attribute is true, such as RESUME for enabled.
	AlterKeywords []AfterCreateAlterKeyword
}

// AfterCreateAlterKeyword calls ALTER <Keyword> after create when the bool attribute is true.
type AfterCreateAlterKeyword struct {
	// Attribute is the Terraform name of the bool attribute.
	Attribute string
	// Keyword is the SQL keyword, such as RESUME.
	Keyword string
	// Field is the Go field of the ALTER request, such as Resume.
	Field string
}

// BeforeCreateSteps are the data of the generated BeforeCreate. Each kind of step has its own field here and its own block
// in before_create.tmpl. To run another kind of step before create, add a field and a block.
type BeforeCreateSteps struct {
	// AdoptExisting adopts an object that Snowflake creates itself, such as the PUBLIC schema, instead of creating it.
	AdoptExisting *BeforeCreateAdoptExisting
}

// BeforeCreateAdoptExisting updates an existing object instead of creating it when the identity name is one of Names.
type BeforeCreateAdoptExisting struct {
	// Names are compared to the identity name ignoring case and outer spaces.
	Names []string
	// ObjectWords is the lowercase object type for the log, such as schema.
	ObjectWords string
}

// AlterStatementsModel lists the fields of the resource ext request struct that hold one ALTER statement each.
// Set and Unset are in send order: the statements of the SSOT definition first, the main set and unset last.
// Requests is Set, then Unset, with the SDK request kind of each field.
type AlterStatementsModel struct {
	Set      []string
	Unset    []string
	Requests []AlterStatementRequest
}

type AlterStatementRequest struct {
	Field string
	// Kind is the SDK request stem: Set or Unset.
	Kind string
}

func newAlterStatementsModel(obj ssot.Object) AlterStatementsModel {
	// A group is a field in declaration order. The main group is the bare set (unset). A block without one sends an empty
	// main request last, which carries the parameters.
	fields := func(main string, groups []ssot.AlterGroupSpec) []string {
		names := collections.Map(groups, func(group ssot.AlterGroupSpec) string { return main + group.Name })
		if !slices.ContainsFunc(groups, func(group ssot.AlterGroupSpec) bool { return group.Name == ssot.MainAlterGroup }) {
			names = append(names, main)
		}
		return names
	}
	model := AlterStatementsModel{
		Set:   fields("set", obj.AlterSetGroups()),
		Unset: fields("unset", obj.AlterUnsetGroups()),
	}
	for _, field := range model.Set {
		model.Requests = append(model.Requests, AlterStatementRequest{Field: field, Kind: "Set"})
	}
	for _, field := range model.Unset {
		model.Requests = append(model.Requests, AlterStatementRequest{Field: field, Kind: "Unset"})
	}
	return model
}

// ParameterHooks returns the catalog parameters for this object's parameter level, in catalog order.
func (m ResourceModel) ParameterHooks() []ParameterHookModel {
	parameters := ssotdefs.ParameterDefsForLevel(ssot.Object{Name: m.Name}.ParameterLevel())
	hooks := make([]ParameterHookModel, 0, len(parameters))
	for _, parameter := range parameters {
		hooks = append(hooks, newParameterHookModel(parameter))
	}
	return hooks
}

func (m ResourceModel) Prefix() string {
	return genhelpers.FirstLetterLowercase(m.Name)
}

// AlterRequestsType is the struct that holds the SET and UNSET requests of the update. The generated ApplySetUnset
// declares it in the generated file, with one field per group of the ALTER spec. An object that keeps ApplySetUnset
// in the ext file declares the struct there, because it may carry more than the groups.
func (m ResourceModel) AlterRequestsType() string {
	if m.Hooks.ApplySetUnset.Generate {
		return m.Prefix() + "AlterRequests"
	}
	return m.Prefix() + "AlterRequestsExt"
}

// HierarchyRenameParents reports whether the generated Update moves a schema object to another database or schema.
func (m ResourceModel) HierarchyRenameParents() bool {
	return m.Hooks.HierarchyRename.Generate && m.IdType == "SchemaObjectIdentifier" && !m.IdentityHierarchyForceNew
}

// HierarchyRenameName reports whether the generated Update renames the object name.
func (m ResourceModel) HierarchyRenameName() bool {
	return m.Hooks.HierarchyRename.Generate && !m.IdentityNameForceNew
}

// ShowObjectName is the SDK type of the SHOW row passed to SetConfigFields.
// A variant's row is the parent type.
func (m ResourceModel) ShowObjectName() string {
	if m.VariantOf != nil && m.VariantOf.Name != "" {
		return m.VariantOf.Name
	}
	return m.Name
}

func (m ResourceModel) EffectiveSdkCreateRequestObjectName() string {
	if m.SdkCreateRequestObjectName != "" {
		return m.SdkCreateRequestObjectName
	}
	return m.Name
}

func (m ResourceModel) EffectiveSdkParametersObjectName() string {
	if m.VariantOf != nil {
		return m.VariantOf.Name
	}
	return m.Name
}

// EffectiveSdkAlterObjectName is the SET and UNSET request stem.
// A variant uses its own name, so WarehouseAdaptive is WarehouseAdaptiveSetRequest.
func (m ResourceModel) EffectiveSdkAlterObjectName() string {
	return m.Name
}

func (m ResourceModel) ParseIdFunc() string {
	if m.Hooks.ParseId.Generate {
		return m.Prefix() + "ParseId"
	}
	return m.Prefix() + "ParseIdExt"
}

func (m ResourceModel) HasUpdateCalls() bool {
	for _, attr := range m.Attributes {
		if attr.UpdateCall != nil {
			return true
		}
	}
	return false
}

// RequiresUpdate reports whether the schema has a field Terraform can change in place.
// An Update function is rejected when every field is ForceNew or computed without Optional.
func (m ResourceModel) RequiresUpdate() bool {
	if m.Hooks.IdentitySchema.Present && !m.Hooks.IdentitySchema.Generate {
		return true
	}
	if m.Hooks.AttributesSchema.Present && !m.Hooks.AttributesSchema.Generate {
		return true
	}
	if m.AttributesModification {
		return true
	}
	if m.Hooks.ParametersAttributesSchema.Present {
		return true
	}
	if m.Hooks.IdentitySchema.Generate {
		if !m.IdentityNameForceNew {
			return true
		}
		if m.IdType == "DatabaseObjectIdentifier" || m.IdType == "SchemaObjectIdentifier" {
			if !m.IdentityHierarchyForceNew {
				return true
			}
		}
	}
	for _, attr := range m.Attributes {
		if !attr.ForceNew {
			return true
		}
	}
	return false
}

func ModelFromInputObject(input ResourceObject, preamble *genhelpers.PreambleModel, _ *log.Logger) ResourceModel {
	hooks := ssot.DefaultResourceExtHooks(input.ResourceHooks...)
	resourceAttributes := input.ResourceAttributes()
	attributes := applyResourceAttributeRules(input.IdType, resourceAttributes)
	applyExactlyOneOf(input.ExactlyOneOf, attributes)
	if hooks.NewCreateRequest.Generate || hooks.ApplyCreateOptionals.Generate {
		attachCreateMappings(input.Name, resourceAttributes, attributes)
	}
	if hooks.ApplySetUnset.Generate {
		attachUpdateMappings(input.Object, resourceAttributes, attributes)
	}
	model := ResourceModel{
		Name:                                input.Name,
		NamePlural:                          input.NamePlural,
		DescribeFailRead:                    input.DescribeFailRead,
		DescribeUsesShowId:                  input.DescribeUsesShowId,
		HasPartialOnUpdate:                  input.HasPartialOnUpdate,
		AdditionalDescribeOutputDescription: input.AdditionalDescribeOutputDescription,
		DescribeDetailsName:                 input.DescribeDetailsName,
		IdType:                              input.IdType,
		LegacyIdEncoding:                    input.LegacyIdEncoding,
		PreviewFeature:                      input.PreviewFeature,
		IdentityNameForceNew:                input.IdentityNameForceNew,
		IdentityHierarchyForceNew:           input.IdentityHierarchyForceNew,
		AdditionalIdentityNameDescription:   input.AdditionalIdentityNameDescription,
		ShowColumns:                         input.Show.Columns,
		DescribeColumns:                     input.Describe.Columns,
		SdkCreateRequestObjectName:          input.SdkCreateRequestObjectName,
		VariantOf:                           input.VariantOf,
		Attributes:                          attributes,
		ShowSchemaOwner:                     input.SchemaShowOwnerName(),
		AttributesModification:              input.AttributesModification,
		ObjectType:                          input.ObjectType(),
		Hooks:                               hooks,
		SchemaStateUpgraders:                input.SchemaStateUpgraders,
		AdditionalCustomizeDiffs:            input.AdditionalCustomizeDiffs,
		BeforeCreateSteps:                   beforeCreateSteps(input.Object),
		AfterCreateSteps:                    afterCreateSteps(input.Object, resourceAttributes),
		AlterStatements:                     newAlterStatementsModel(input.Object),
		PreambleModel:                       preamble,
	}
	if computed, columns := showComputedAttributes(input.Object); len(computed) > 0 {
		if !hooks.ShowByIdSafelyInSdk.Present {
			panic(fmt.Sprintf("object %q: computed attributes are set from SHOW and need ShowByIdSafelyInSdk", input.Name))
		}
		model.ComputedColumns = columns
		model.ComputedAttributes = collections.Map(computed, func(attr ssot.Attribute) ResourceAttribute {
			return computedResourceAttribute(input.Name, attr)
		})
	}
	model.applyDerivedOutputDiffSuppress(resourceAttributes)
	// Building the mappings panics when a generated hook cannot describe a column.
	model.ShowExternalChanges()
	model.DescribeExternalChanges()
	return model
}

// afterCreateSteps collects the SetAfterCreate attributes. A bool calls the ALTER keyword linked to true, because CREATE
// leaves the object in the false state. Validation rejects any other SetAfterCreate attribute for now.
func afterCreateSteps(obj ssot.Object, attributes []ssot.Attribute) AfterCreateSteps {
	var steps AfterCreateSteps
	for _, attribute := range attributes {
		if !attribute.SetAfterCreate {
			continue
		}
		if clause, ok := obj.LinkedAlterKeyword(attribute.Name, true); ok {
			steps.AlterKeywords = append(steps.AlterKeywords, AfterCreateAlterKeyword{
				Attribute: attribute.Name,
				Keyword:   clause.SQL,
				Field:     genhelpers.SqlToGoName(clause.SQL, true),
			})
		}
	}
	return steps
}

func beforeCreateSteps(obj ssot.Object) BeforeCreateSteps {
	var steps BeforeCreateSteps
	if len(obj.IdentityNamesCreatedBySnowflake) > 0 {
		steps.AdoptExisting = &BeforeCreateAdoptExisting{
			Names:       obj.IdentityNamesCreatedBySnowflake,
			ObjectWords: strings.ToLower(obj.ObjectType()),
		}
	}
	return steps
}
