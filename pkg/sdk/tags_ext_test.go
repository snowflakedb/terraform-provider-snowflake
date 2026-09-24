package sdk

import (
	"fmt"
	"testing"
)

func init() {
	id := tagsTestIdSchemaObjectIdentifier
	allowedValues := &AllowedValues{
		Values: []StringAllowEmpty{{Value: "value1"}, {Value: "value2"}},
	}
	mp1ID := NewSchemaObjectIdentifierInSchema(id.SchemaId(), "policy1")
	mp2ID := NewSchemaObjectIdentifierInSchema(id.SchemaId(), "policy2")
	maskingPolicies := []TagMaskingPolicy{{Name: mp1ID}, {Name: mp2ID}}
	renameTarget := randomSchemaObjectIdentifierInSchema(id.SchemaId())
	renameTargetDifferentDb := randomSchemaObjectIdentifier()

	tagsTests.Create.
		withAdditionalValidationCase(
			"validation_Create_AllowedValues_count",
			func(opts *CreateTagOptions) {
				opts.AllowedValues = &AllowedValues{Values: []StringAllowEmpty{}}
			},
			errIntValue("AllowedValues", "Values", IntErrGreater, 0),
		).
		withExpectedSqlf(
			case_Tags_sql_Create_basic,
			`CREATE TAG %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tags_sql_Create_all,
			func(opts *CreateTagOptions) {
				opts.IfNotExists = new(true)
				opts.AllowedValues = allowedValues
				opts.Propagate = &TagPropagate{
					PropagationMethod: new(TagPropagationOnDependencyAndDataMovement),
					OnConflict:        &TagOnConflict{CustomValue: new("FAIL")},
				}
				opts.Comment = new("comment")
			},
			`CREATE TAG IF NOT EXISTS %s ALLOWED_VALUES 'value1', 'value2' PROPAGATE = ON_DEPENDENCY_AND_DATA_MOVEMENT ON_CONFLICT = 'FAIL' COMMENT = 'comment'`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_propagateOnly",
			func(opts *CreateTagOptions) {
				opts.Propagate = &TagPropagate{PropagationMethod: new(TagPropagationOnDependency)}
			},
			`CREATE TAG %s PROPAGATE = ON_DEPENDENCY`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_onConflictAllowedValuesSequence",
			func(opts *CreateTagOptions) {
				opts.Propagate = &TagPropagate{
					PropagationMethod: new(TagPropagationOnDataMovement),
					OnConflict:        &TagOnConflict{AllowedValuesSequence: new(true)},
				}
			},
			`CREATE TAG %s PROPAGATE = ON_DATA_MOVEMENT ON_CONFLICT = ALLOWED_VALUES_SEQUENCE`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateTagOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE TAG %s`, id.FullyQualifiedName(),
		)

	tagsTests.Alter.
		withAdditionalValidationCase(
			"validation_Alter_Set_multipleFields",
			func(opts *AlterTagOptions) {
				opts.Set = &TagSet{
					Comment: new("comment"),
					MaskingPolicies: &TagSetMaskingPolicies{
						MaskingPolicies: []TagMaskingPolicy{{Name: mp1ID}},
					},
				}
			},
			errOneOf("TagSet", "MaskingPolicies", "AllowedValues", "Propagate", "Comment"),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_emptyMaskingPolicies",
			func(opts *AlterTagOptions) {
				opts.Set = &TagSet{MaskingPolicies: &TagSetMaskingPolicies{}}
			},
			errIntValue("TagSet.MaskingPolicies", "MaskingPolicies", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Alter_Unset_emptyMaskingPolicies",
			func(opts *AlterTagOptions) {
				opts.Unset = &TagUnset{MaskingPolicies: &TagUnsetMaskingPolicies{}}
			},
			errIntValue("TagUnset.MaskingPolicies", "MaskingPolicies", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Alter_Add_AllowedValues_count",
			func(opts *AlterTagOptions) {
				opts.Add = &TagAdd{AllowedValues: &AllowedValues{Values: []StringAllowEmpty{}}}
			},
			errIntValue("AllowedValues", "Values", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Alter_Drop_AllowedValues_count",
			func(opts *AlterTagOptions) {
				opts.Drop = &TagDrop{AllowedValues: &AllowedValues{Values: []StringAllowEmpty{}}}
			},
			errIntValue("AllowedValues", "Values", IntErrGreater, 0),
		).
		withAdditionalValidationCase(
			"validation_Alter_Set_AllowedValues_count",
			func(opts *AlterTagOptions) {
				opts.Set = &TagSet{AllowedValues: &AllowedValues{Values: []StringAllowEmpty{}}}
			},
			errIntValue("AllowedValues", "Values", IntErrGreater, 0),
		).
		withModifyAndExpectedSqlf(
			case_Tags_sql_Alter_Add,
			func(opts *AlterTagOptions) {
				opts.Add = &TagAdd{AllowedValues: allowedValues}
			},
			`ALTER TAG %s ADD ALLOWED_VALUES 'value1', 'value2'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tags_sql_Alter_Drop,
			func(opts *AlterTagOptions) {
				opts.Drop = &TagDrop{AllowedValues: allowedValues}
			},
			`ALTER TAG %s DROP ALLOWED_VALUES 'value1', 'value2'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tags_sql_Alter_Set,
			func(opts *AlterTagOptions) { opts.Set = &TagSet{Comment: new("comment")} },
			`ALTER TAG %s SET COMMENT = 'comment'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tags_sql_Alter_Unset,
			func(opts *AlterTagOptions) { opts.Unset = &TagUnset{AllowedValues: new(true)} },
			`ALTER TAG %s UNSET ALLOWED_VALUES`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tags_sql_Alter_RenameTo,
			func(opts *AlterTagOptions) {
				opts.IfExists = new(true)
				opts.RenameTo = &renameTarget
			},
			`ALTER TAG IF EXISTS %s RENAME TO %s`,
			id.FullyQualifiedName(), renameTarget.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_RenameTo_differentDatabase",
			func(opts *AlterTagOptions) { opts.RenameTo = &renameTargetDifferentDb },
			`ALTER TAG %s RENAME TO %s`,
			id.FullyQualifiedName(), renameTargetDifferentDb.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_maskingPolicies",
			func(opts *AlterTagOptions) {
				opts.Set = &TagSet{
					MaskingPolicies: &TagSetMaskingPolicies{
						MaskingPolicies: maskingPolicies,
						Force:           new(true),
					},
				}
			},
			`ALTER TAG %s SET MASKING POLICY %s, MASKING POLICY %s FORCE`,
			id.FullyQualifiedName(), mp1ID.FullyQualifiedName(), mp2ID.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_propagate",
			func(opts *AlterTagOptions) {
				opts.Set = &TagSet{Propagate: &TagPropagate{PropagationMethod: new(TagPropagationOnDependency)}}
			},
			`ALTER TAG %s SET PROPAGATE = ON_DEPENDENCY`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Set_propagateOnConflict",
			func(opts *AlterTagOptions) {
				opts.Set = &TagSet{
					Propagate: &TagPropagate{
						PropagationMethod: new(TagPropagationOnDependencyAndDataMovement),
						OnConflict:        &TagOnConflict{CustomValue: new("FAIL")},
					},
				}
			},
			`ALTER TAG %s SET PROPAGATE = ON_DEPENDENCY_AND_DATA_MOVEMENT ON_CONFLICT = 'FAIL'`,
			id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_maskingPolicies",
			func(opts *AlterTagOptions) {
				opts.Unset = &TagUnset{
					MaskingPolicies: &TagUnsetMaskingPolicies{MaskingPolicies: maskingPolicies},
				}
			},
			`ALTER TAG %s UNSET MASKING POLICY %s, MASKING POLICY %s`,
			id.FullyQualifiedName(), mp1ID.FullyQualifiedName(), mp2ID.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_propagate",
			func(opts *AlterTagOptions) { opts.Unset = &TagUnset{Propagate: new(true)} },
			`ALTER TAG %s UNSET PROPAGATE`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_onConflict",
			func(opts *AlterTagOptions) { opts.Unset = &TagUnset{OnConflict: new(true)} },
			`ALTER TAG %s UNSET ON_CONFLICT`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_Unset_comment",
			func(opts *AlterTagOptions) { opts.Unset = &TagUnset{Comment: new(true)} },
			`ALTER TAG %s UNSET COMMENT`, id.FullyQualifiedName(),
		)

	tagsTests.Show.
		withAdditionalValidationCase(
			"validation_Show_Like_patternRequired",
			func(opts *ShowTagOptions) { opts.Like = &Like{} },
			ErrPatternRequiredForLikeKeyword,
		).
		withAdditionalValidationCase(
			"validation_Show_In_noScope",
			func(opts *ShowTagOptions) { opts.In = &ExtendedIn{} },
			errExactlyOneOf("ShowTagOptions.In", "Account", "Database", "Schema"),
		).
		withAdditionalValidationCase(
			"validation_Show_In_moreThanOneScope",
			func(opts *ShowTagOptions) {
				opts.In = &ExtendedIn{In: In{Account: new(true), Database: id.DatabaseId()}}
			},
			errExactlyOneOf("ShowTagOptions.In", "Account", "Database", "Schema"),
		).
		withExpectedSql(case_Tags_sql_Show_basic, `SHOW TAGS`).
		withModifyAndExpectedSqlf(
			case_Tags_sql_Show_Like,
			func(opts *ShowTagOptions) { opts.Like = &Like{Pattern: new("test")} },
			`SHOW TAGS LIKE 'test'`,
		).
		withModifyAndExpectedSqlf(
			case_Tags_sql_Show_In,
			func(opts *ShowTagOptions) { opts.In = &ExtendedIn{In: In{Account: new(true)}} },
			`SHOW TAGS IN ACCOUNT`,
		).
		withModifyAndExpectedSqlf(
			case_Tags_sql_Show_all,
			func(opts *ShowTagOptions) {
				opts.Like = &Like{Pattern: new("test")}
				opts.In = &ExtendedIn{In: In{Account: new(true)}}
			},
			`SHOW TAGS LIKE 'test' IN ACCOUNT`,
		)

	tagsTests.Drop.
		withExpectedSqlf(
			case_Tags_sql_Drop_basic,
			`DROP TAG %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_Tags_sql_Drop_all,
			func(opts *DropTagOptions) { opts.IfExists = new(true) },
			`DROP TAG IF EXISTS %s`, id.FullyQualifiedName(),
		)

	tagsTests.Undrop.
		withExpectedSqlf(
			case_Tags_sql_Undrop_basic,
			`UNDROP TAG %s`, id.FullyQualifiedName(),
		)

	tagId := NewAccountObjectIdentifier("tag1")
	tagId2 := NewAccountObjectIdentifier("tag2")
	accountId := randomAccountIdentifier()

	tagsTests.Set.
		withDefaultOpts(func() *SetTagOptions {
			return &SetTagOptions{
				objectType: ObjectTypeStage,
				objectName: id,
				SetTags:    []TagAssociation{{Name: tagId, Value: "value1"}},
			}
		}).
		withAdditionalValidationCase(
			"validation_Set_invalidObjectType",
			func(opts *SetTagOptions) { opts.objectType = ObjectType("SEQUENCE;") },
			fmt.Errorf("invalid object type: %s contains disallowed characters; it must follow this regex: %s", "SEQUENCE;", allowedUnquotedCharactersRegex.String()),
		).
		withAdditionalSqlCasef(
			"sql_Set_unknownObjectType",
			func(opts *SetTagOptions) { opts.objectType = ObjectTypePostgresInstance },
			`ALTER %s %s SET TAG "tag1" = 'value1'`, ObjectTypePostgresInstance, id.FullyQualifiedName(),
		).
		withExpectedSqlf(
			case_Tags_sql_Set_basic,
			`ALTER %s %s SET TAG "tag1" = 'value1'`, ObjectTypeStage, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Set_onAccount",
			func(opts *SetTagOptions) { opts.objectName = accountId },
			`ALTER %s %s SET TAG "tag1" = 'value1'`, ObjectTypeStage, accountId.FullyQualifiedName(),
		)

	tagsTests.Unset.
		withDefaultOpts(func() *UnsetTagOptions {
			return &UnsetTagOptions{
				objectType: ObjectTypeStage,
				IfExists:   new(true),
				objectName: id,
				UnsetTags:  []ObjectIdentifier{tagId, tagId2},
			}
		}).
		withAdditionalValidationCase(
			"validation_Unset_invalidObjectType",
			func(opts *UnsetTagOptions) { opts.objectType = ObjectType("SEQUENCE;") },
			fmt.Errorf("invalid object type: %s contains disallowed characters; it must follow this regex: %s", "SEQUENCE;", allowedUnquotedCharactersRegex.String()),
		).
		withExpectedSqlf(
			case_Tags_sql_Unset_basic,
			`ALTER %s IF EXISTS %s UNSET TAG "tag1", "tag2"`, ObjectTypeStage, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Unset_onAccount",
			func(opts *UnsetTagOptions) {
				opts.IfExists = nil
				opts.objectName = accountId
			},
			`ALTER %s %s UNSET TAG "tag1", "tag2"`, ObjectTypeStage, accountId.FullyQualifiedName(),
		)
}

func TestTags_Set_withColumn(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	objectId := randomTableColumnIdentifierInSchemaObject(id)
	tagId := randomSchemaObjectIdentifier()
	request := NewSetTagRequest(ObjectTypeColumn, objectId).WithSetTags([]TagAssociation{
		{
			Name:  tagId,
			Value: "value1",
		},
	})
	request.adjust()
	opts := request.toOpts()
	assertOptsValidAndSqlEqualsf(t, opts, `ALTER TABLE %s MODIFY COLUMN "%s" SET TAG %s = 'value1'`, id.FullyQualifiedName(), objectId.columnName, tagId.FullyQualifiedName())
}

func TestTags_Unset_withColumn(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	objectId := randomTableColumnIdentifierInSchemaObject(id)
	tagId1 := randomSchemaObjectIdentifier()
	tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
	request := NewUnsetTagRequest(ObjectTypeColumn, objectId).
		WithUnsetTags([]ObjectIdentifier{tagId1, tagId2}).
		WithIfExists(true)
	request.adjust()
	opts := request.toOpts()
	assertOptsValidAndSqlEqualsf(t, opts, `ALTER %s IF EXISTS %s MODIFY COLUMN "%s" UNSET TAG %s, %s`, opts.objectType, id.FullyQualifiedName(), objectId.Name(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
}

// The SDK enforces only the lower bound on AllowedValues; the maximum is enforced by Snowflake
// server-side, so a large list — even above the current server default — must pass client-side validation.
func TestTags_AllowedValues_noClientSideUpperBound(t *testing.T) {
	id := randomSchemaObjectIdentifier()

	// 5001 sits above the current server-side default (5000): the SDK must not impose a ceiling of its own.
	values := make([]string, 5001)
	for i := range values {
		values[i] = fmt.Sprintf("value%d", i)
	}
	allowedValues := NewAllowedValuesRequestFromStrings(values)

	assertOptsValid(t, NewCreateTagRequest(id).WithAllowedValues(*allowedValues).toOpts())
	assertOptsValid(t, NewAlterTagRequest(id).WithAdd(*NewTagAddRequest().WithAllowedValues(*allowedValues)).toOpts())
}
