package ssot

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
)

// attributePlace is where an attribute sits. That decides whether an explicit default can be removed.
type attributePlace struct {
	// nested is a child, an element, or a variant field. The clause then also selects the statement projection.
	nested bool
	// stem is the request stem used to deduce a mapper name: the object name plus each enclosing block.
	stem string
	// checkUpdate is set on a top-level attribute that generation would give an update helper.
	checkUpdate bool
	// setOnly is true when the attribute is in the SET block and not in the UNSET block.
	setOnly bool
}

// RedundantExplicitDefaults lists attribute fields set to the value omission already produces.
// A definition leaves those fields unset. Each property is one check in validateExplicitDefaults.
func RedundantExplicitDefaults(obj Object) []string {
	places := map[*Attribute]attributePlace{}
	wheres := map[*Attribute]string{}
	for _, attr := range obj.explicitDefaultRoots() {
		walkExplicitDefaults(obj.Name+"."+attr.Name, attr, obj.mappingPlace(attr), places, wheres)
	}
	findings := obj.redundantObjectDefaults()
	findings = slices.Grow(findings, len(places))
	for attr, place := range places {
		findings = append(findings, attr.validateExplicitDefaults(wheres[attr], place)...)
	}
	return findings
}

// redundantObjectDefaults lists object name overrides set to the value omission already produces.
// DescribeDetailsName empty is Name plus Details. BaseVariant.ShowPlainName, SdkCreateRequestObjectName empty is Name.
// Variant.CreateName and Variant.AlterName empty are the base object's type words.
// Variant.RequestObjectName empty is Name.
func (obj Object) redundantObjectDefaults() []string {
	var findings []string
	if obj.DescribeDetailsName != "" && obj.DescribeDetailsName == obj.Name+"Details" {
		findings = append(findings, explicitDefaultFinding(obj.Name, "DescribeDetailsName"))
	}
	if obj.ShowPlainName() != "" && obj.ShowPlainName() == obj.Name {
		findings = append(findings, explicitDefaultFinding(obj.Name, "BaseVariant.ShowPlainName"))
	}
	if obj.SdkCreateRequestObjectName != "" && obj.SdkCreateRequestObjectName == obj.Name {
		findings = append(findings, explicitDefaultFinding(obj.Name, "SdkCreateRequestObjectName"))
	}
	if obj.Variant == nil || obj.Variant.Of == nil || obj.Variant.Of.Name == "" {
		return findings
	}
	if obj.Variant.RequestObjectName != "" && obj.Variant.RequestObjectName == obj.Name {
		findings = append(findings, explicitDefaultFinding(obj.Name, "Variant.RequestObjectName"))
	}
	base := obj.Variant.Of.ObjectType()
	if obj.Variant.CreateName != "" && obj.Variant.CreateName == base {
		findings = append(findings, explicitDefaultFinding(obj.Name, "Variant.CreateName"))
	}
	if obj.Variant.AlterName != "" && obj.Variant.AlterName == base {
		findings = append(findings, explicitDefaultFinding(obj.Name, "Variant.AlterName"))
	}
	return findings
}

// walkExplicitDefaults records the strictest place for attr, then walks its children, element, and variant fields.
// where is the path already built for attr. An element has no name, so it keeps its parent's path.
// A later visit of the same pointer merges the place and does not walk again.
func walkExplicitDefaults(where string, attr *Attribute, place attributePlace, places map[*Attribute]attributePlace, wheres map[*Attribute]string) {
	if attr == nil {
		return
	}
	if prev, ok := places[attr]; ok {
		prev.nested = prev.nested || place.nested
		places[attr] = prev
		return
	}
	places[attr] = place
	wheres[attr] = where
	nested := attributePlace{nested: true, stem: place.childStem(attr)}
	for _, child := range attributePointers(attr.Children) {
		walkExplicitDefaults(where+"."+child.Name, child, nested, places, wheres)
	}
	walkExplicitDefaults(where, attr.Elem, nested, places, wheres)
	for _, field := range clauseVariantFields(attr.CreateSql, attr.SetSql, attr.UnsetSql) {
		walkExplicitDefaults(where+"."+field.Name, field, nested, places, wheres)
	}
}

func (obj Object) mappingPlace(attr *Attribute) attributePlace {
	restricted := obj.HasAlterSetBlock() || obj.HasAlterUnsetBlock()
	inSet := named(obj.SetAttributes(), attr.Name)
	inUnset := named(obj.UnsetAttributes(), attr.Name)
	return attributePlace{
		stem:        obj.Name,
		checkUpdate: !restricted || inSet || inUnset,
		setOnly:     restricted && !inUnset,
	}
}

func named(attrs []*Attribute, name string) bool {
	for _, attr := range attrs {
		if attr != nil && attr.Name == name {
			return true
		}
	}
	return false
}

func (p attributePlace) childStem(attr *Attribute) string {
	if attr.Name == "" || (len(attr.Children) == 0 && attr.Elem == nil) {
		return p.stem
	}
	return p.stem + genhelpers.SnakeCaseToCamel(attr.Name)
}

func (obj Object) explicitDefaultRoots() []*Attribute {
	var suspend []*Attribute
	if obj.SuspendForAlter != nil {
		suspend = append(slices.Clone(obj.SuspendForAlter.When), obj.SuspendForAlter.State)
	}
	return slices.Concat(
		attributePointers(obj.Attributes),
		attributePointers(obj.ComputedAttributes),
		obj.Show.Columns,
		obj.Describe.Columns,
		obj.ShowGenericColumns(),
		obj.DescribeRows.Columns,
		obj.Create.Clauses,
		obj.Create.KeywordsBefore,
		obj.SetAttributes(),
		obj.UnsetAttributes(),
		// These slots are empty on clauses that do not carry an attribute. A variant's suspend state is empty too.
		present(suspend),
		present(collections.Map(obj.Alter.Clauses, func(clause AlterClause) *Attribute { return clause.Attribute })),
		present(collections.Map(obj.Alter.Clauses, func(clause AlterClause) *Attribute { return clause.Block })),
		present(customOperationAttributes(obj.CustomOperations)),
	)
}

func present(attrs []*Attribute) []*Attribute {
	return collections.Filter(attrs, func(attr *Attribute) bool { return attr != nil })
}

func customOperationAttributes(ops []CustomOperationSpec) []*Attribute {
	return slices.Concat(collections.Map(ops, func(op CustomOperationSpec) []*Attribute {
		return collections.Map(op.Clauses, func(clause CustomOperationClause) *Attribute { return clause.Attribute })
	})...)
}

func attributePointers(attrs []Attribute) []*Attribute {
	out := make([]*Attribute, len(attrs))
	for i := range attrs {
		out[i] = &attrs[i]
	}
	return out
}

func clauseVariantFields(clauses ...SqlClause) []*Attribute {
	return slices.Concat(collections.Map(clauses, func(clause SqlClause) []*Attribute {
		return slices.Concat(collections.Map(clause.Variants, func(variant SqlVariant) []*Attribute {
			return variant.Fields
		})...)
	})...)
}

// validateExplicitDefaults runs one check per property. Add a property by adding a check here and a method below.
func (a Attribute) validateExplicitDefaults(where string, place attributePlace) []string {
	var findings []string
	if a.validateCreateSql(place) {
		findings = append(findings, explicitDefaultFinding(where, clauseUseName(ClauseCreate)))
	}
	if a.validateSetSql(place) {
		findings = append(findings, explicitDefaultFinding(where, clauseUseName(ClauseSet)))
	}
	if a.validateUnsetSql(place) {
		findings = append(findings, explicitDefaultFinding(where, clauseUseName(ClauseUnset)))
	}
	for _, use := range []ClauseUse{ClauseCreate, ClauseSet, ClauseUnset} {
		if a.validateClauseSql(use, place) {
			findings = append(findings, explicitDefaultFinding(where, clauseUseName(use)+" SQL"))
		}
		if a.validateClauseField(use) {
			findings = append(findings, explicitDefaultFinding(where, clauseUseName(use)+" Field"))
		}
	}
	if a.validateShowName() {
		findings = append(findings, explicitDefaultFinding(where, "ShowName"))
	}
	if a.validateSchemaName() {
		findings = append(findings, explicitDefaultFinding(where, "SchemaName"))
	}
	if a.validateShowEnum() {
		findings = append(findings, explicitDefaultFinding(where, "ShowEnum"))
	}
	if a.validateShowString() {
		findings = append(findings, explicitDefaultFinding(where, "ShowString"))
	}
	if a.validateDescribeString() {
		findings = append(findings, explicitDefaultFinding(where, "DescribeString"))
	}
	if a.validateDescribeEnum() {
		findings = append(findings, explicitDefaultFinding(where, "DescribeEnum"))
	}
	if a.validateCreateMapping(place) {
		findings = append(findings, explicitDefaultFinding(where, "CreateMapping"))
	}
	if a.validateUpdateMapping(place) {
		findings = append(findings, explicitDefaultFinding(where, "UpdateMapping"))
	}
	if a.validateUpdateField(place) {
		findings = append(findings, explicitDefaultFinding(where, "UpdateMapping Field"))
	}
	return findings
}

// validateClauseSql reports a clause SQL text set to the uppercase attribute name. Empty SQL already means that.
// A clause that restates the whole type default is reported there instead.
func (a Attribute) validateClauseSql(use ClauseUse, place attributePlace) bool {
	clause := a.rawSqlClause(use)
	if !clause.Present || clause.SQL == "" || clause.SQL != strings.ToUpper(a.Name) || !a.clauseSqlUsesAttributeName(clause) {
		return false
	}
	return !a.restatesDefaultClause(use, place)
}

// clauseSqlUsesAttributeName reports that an empty SQL is the uppercase attribute name in the text the generator writes.
// A static clause, an object, and a preexisting struct print clause.SQL itself. A list parameter or a value keyword that
// sets FieldName also prints clause.SQL, so clearing it would drop the prefix.
func (a Attribute) clauseSqlUsesAttributeName(clause SqlClause) bool {
	switch clause.Kind {
	case SqlKindParameter:
		if a.IsObjectList() {
			return true
		}
		return clause.ItemStruct == "" || clause.Field == ""
	case SqlKindKeyword:
		if !clause.ValueKeyword {
			return true
		}
		return clause.Field == ""
	case SqlKindIdentifier:
		return true
	default:
		return false
	}
}

// validateClauseField reports a FieldName set to the camel-cased attribute name.
// Only a clause whose omitted field name is that camel case is checked. A parameter, a plain keyword, and an identifier
// derive the name from the SQL text instead.
func (a Attribute) validateClauseField(use ClauseUse) bool {
	clause := a.rawSqlClause(use)
	if !clause.Present || clause.Field == "" || !clause.fieldNameUsesAttributeName() {
		return false
	}
	return clause.Field == genhelpers.SnakeCaseToCamel(a.Name)
}

// fieldNameUsesAttributeName reports that an empty Field is the camel-cased attribute name.
// An object and a preexisting struct do that. A value keyword does that only with SQL: without SQL, an empty Field writes
// the attribute name as the keyword, so a FieldName that leaves the keyword out is not redundant. A parameter, a plain
// keyword, and an identifier use the SQL text.
func (c SqlClause) fieldNameUsesAttributeName() bool {
	switch c.Kind {
	case SqlKindObject, SqlKindPreexistingStruct:
		return true
	case SqlKindKeyword:
		return c.ValueKeyword && c.SQL != ""
	default:
		return false
	}
}

// validateUpdateField reports a set-only update field set to the camel-cased attribute name.
// An empty field already writes that name. A mapping that restates the whole update call is reported there instead.
func (a Attribute) validateUpdateField(place attributePlace) bool {
	if !place.checkUpdate || a.UpdateMapping == nil {
		return false
	}
	explicit := a.UpdateMapping.UpdateCall(place.stem, a)
	if explicit.Field == "" || explicit.Field != genhelpers.SnakeCaseToCamel(a.Name) {
		return false
	}
	return !a.validateUpdateMapping(place)
}

func (a Attribute) validateCreateMapping(place attributePlace) bool {
	if a.CreateMapping == nil || place.stem == "" {
		return false
	}
	explicit := a.CreateMapping.CreateCall(place.stem, a)
	cleared := a
	cleared.CreateMapping = nil
	deduced, ok := cleared.DeducedCreateCall(place.stem)
	return ok && explicit == deduced
}

func (a Attribute) validateUpdateMapping(place attributePlace) bool {
	if !place.checkUpdate || a.UpdateMapping == nil || place.stem == "" {
		return false
	}
	explicit := a.UpdateMapping.UpdateCall(place.stem, a)
	cleared := a
	cleared.UpdateMapping = nil
	deduced, ok := cleared.DeducedUpdateCall(place.stem, place.setOnly)
	return ok && explicit == deduced
}

func explicitDefaultFinding(where, field string) string {
	return fmt.Sprintf("%s %s restates the omitted default", where, field)
}

// validateCreateSql reports a CreateSql set to the clause omission already produces.
func (a Attribute) validateCreateSql(place attributePlace) bool {
	return a.restatesDefaultClause(ClauseCreate, place)
}

// validateSetSql reports a SetSql set to the clause omission already produces.
func (a Attribute) validateSetSql(place attributePlace) bool {
	return a.restatesDefaultClause(ClauseSet, place)
}

// validateUnsetSql reports an UnsetSql set to the clause omission already produces.
func (a Attribute) validateUnsetSql(place attributePlace) bool {
	return a.restatesDefaultClause(ClauseUnset, place)
}

// restatesDefaultClause compares the stored clause with EffectiveSqlClause of a zero clause.
// The type default itself lives in defaultSqlClause.
func (a Attribute) restatesDefaultClause(use ClauseUse, place attributePlace) bool {
	// A SqlOnly attribute has no schema and the validator requires its own CREATE clause, so omitting the clause
	// produces no SQL. The type default is not there to restate.
	if a.SqlOnly || !a.rawSqlClause(use).Present || !a.hasSqlDefault(use) {
		return false
	}
	omitted := a.withoutClause(use)
	if !reflect.DeepEqual(a.EffectiveSqlClause(use), omitted.EffectiveSqlClause(use)) {
		return false
	}
	return !place.nested || !a.clauseChangesProjection(use)
}

func (a Attribute) withoutClause(use ClauseUse) Attribute {
	switch use {
	case ClauseCreate:
		a.CreateSql = SqlClause{}
	case ClauseSet:
		a.SetSql = SqlClause{}
	case ClauseUnset:
		a.UnsetSql = SqlClause{}
	default:
		panic(fmt.Sprintf("attribute %q: unknown sql clause", a.Name))
	}
	return a
}

func (a Attribute) clauseChangesProjection(use ClauseUse) bool {
	omitted := a.withoutClause(use)
	for _, projection := range []ClauseUse{ClauseCreate, ClauseSet, ClauseUnset} {
		if a.inSqlProjection(projection) != omitted.inSqlProjection(projection) {
			return true
		}
	}
	return false
}

func (a Attribute) hasSqlDefault(use ClauseUse) bool {
	if use == ClauseUnset {
		return true
	}
	switch a.SdkType() {
	case AttrInt, AttrBool, AttrString, AttrEnum, AttrList, AttrSet, AttrIdentifier:
		return true
	default:
		return false
	}
}

// validateShowName reports a ShowName set to Name. OutputKey already uses Name when ShowName is empty.
func (a Attribute) validateShowName() bool {
	return a.ShowName != "" && a.ShowName == a.Name
}

// validateSchemaName reports a SchemaName set to the snake_case of the Go field. SchemaKey already produces that when SchemaName is empty.
func (a Attribute) validateSchemaName() bool {
	if a.SchemaName == "" {
		return false
	}
	cleared := a
	cleared.SchemaName = ""
	return a.SchemaName == cleared.SchemaKey()
}

// validateShowEnum reports a ShowEnum set to Enum. ColumnEnum already uses Enum when ShowEnum is empty.
func (a Attribute) validateShowEnum() bool {
	return a.ShowEnum != nil && a.ShowEnum == a.Enum
}

// validateShowString reports ShowString on an attribute the SHOW column already reads as a string.
func (a Attribute) validateShowString() bool {
	return a.ShowString && a.SdkType() == AttrString
}

// validateDescribeString reports DescribeString on an attribute the DESCRIBE column already reads as a string.
func (a Attribute) validateDescribeString() bool {
	return a.DescribeString && a.SdkType() == AttrString
}

// validateDescribeEnum reports a DescribeEnum set to Enum. ForDescribe already uses Enum when DescribeEnum is empty.
func (a Attribute) validateDescribeEnum() bool {
	return a.DescribeEnum != nil && a.DescribeEnum == a.Enum
}
