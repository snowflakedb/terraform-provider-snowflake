package sdk

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *createTagOptions {
		return &createTagOptions{
			name: id,
		}
	}

	t.Run("create with all optional", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.Comment = String("comment")
		opts.AllowedValues = &AllowedValues{
			Values: []StringAllowEmpty{
				{
					Value: "value1",
				},
				{
					Value: "value2",
				},
			},
		}
		opts.Propagate = &TagPropagate{
			PropagationMethod: Pointer(TagPropagationOnDependencyAndDataMovement),
			OnConflict:        &TagOnConflict{CustomValue: String("FAIL")},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE TAG IF NOT EXISTS %s ALLOWED_VALUES 'value1', 'value2' PROPAGATE = ON_DEPENDENCY_AND_DATA_MOVEMENT ON_CONFLICT = 'FAIL' COMMENT = 'comment'`, id.FullyQualifiedName())
	})

	t.Run("create with propagate only", func(t *testing.T) {
		opts := defaultOpts()
		opts.Propagate = &TagPropagate{PropagationMethod: Pointer(TagPropagationOnDependency)}
		assertOptsValidAndSQLEquals(t, opts, `CREATE TAG %s PROPAGATE = ON_DEPENDENCY`, id.FullyQualifiedName())
	})

	t.Run("create with on_conflict allowed_values_sequence", func(t *testing.T) {
		opts := defaultOpts()
		opts.Propagate = &TagPropagate{
			PropagationMethod: Pointer(TagPropagationOnDataMovement),
			OnConflict:        &TagOnConflict{AllowedValuesSequence: Bool(true)},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE TAG %s PROPAGATE = ON_DATA_MOVEMENT ON_CONFLICT = ALLOWED_VALUES_SEQUENCE`, id.FullyQualifiedName())
	})

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*createTagOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: allowed values count", func(t *testing.T) {
		opts := defaultOpts()
		opts.AllowedValues = &AllowedValues{
			Values: []StringAllowEmpty{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntBetween("AllowedValues", "Values", 1, 300))
	})

	t.Run("validation: both ifNotExists and orReplace present", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("createTagOptions", "OrReplace", "IfNotExists"))
	})

	t.Run("validation: multiple errors", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier, errOneOf("createTagOptions", "OrReplace", "IfNotExists"))
	})
}

func TestTagDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *dropTagOptions {
		return &dropTagOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *dropTagOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("drop with name", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP TAG %s`, id.FullyQualifiedName())
	})

	t.Run("drop with if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP TAG IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestTagUndrop(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *undropTagOptions {
		return &undropTagOptions{
			name: id,
		}
	}
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *undropTagOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("undrop with name", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `UNDROP TAG %s`, id.FullyQualifiedName())
	})
}

func TestTagShow(t *testing.T) {
	defaultOpts := func() *showTagOptions {
		return &showTagOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *showTagOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, ErrPatternRequiredForLikeKeyword)
	})

	t.Run("validation: empty in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &ExtendedIn{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("showTagOptions.In", "Account", "Database", "Schema"))
	})

	t.Run("show with empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW TAGS`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{Pattern: String("test")}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TAGS LIKE 'test'`)
	})

	t.Run("show with in", func(t *testing.T) {
		opts := defaultOpts()
		opts.In = &ExtendedIn{
			In: In{
				Account: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TAGS IN ACCOUNT`)
	})
}

func TestTagAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *alterTagOptions {
		return &alterTagOptions{
			name: id,
		}
	}
	defaultAllowedValues := func() *AllowedValues {
		return &AllowedValues{
			Values: []StringAllowEmpty{
				{
					Value: "value1",
				},
				{
					Value: "value2",
				},
			},
		}
	}
	mp1ID := NewSchemaObjectIdentifierInSchema(id.SchemaId(), "policy1")
	mp2ID := NewSchemaObjectIdentifierInSchema(id.SchemaId(), "policy2")
	defaultMaskingPolicies := func() []TagMaskingPolicy {
		return []TagMaskingPolicy{
			{
				Name: mp1ID,
			},
			{
				Name: mp2ID,
			},
		}
	}

	t.Run("alter with rename to and if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.Rename = &TagRename{Name: randomSchemaObjectIdentifierInSchema(id.SchemaId())}
		opts.ifExists = Pointer(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG IF EXISTS %s RENAME TO %s`, id.FullyQualifiedName(), opts.Rename.Name.FullyQualifiedName())
	})

	t.Run("alter with add", func(t *testing.T) {
		opts := defaultOpts()
		opts.Add = &TagAdd{AllowedValues: defaultAllowedValues()}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s ADD ALLOWED_VALUES 'value1', 'value2'`, id.FullyQualifiedName())
	})

	t.Run("alter with drop", func(t *testing.T) {
		opts := defaultOpts()
		opts.Drop = &TagDrop{AllowedValues: defaultAllowedValues()}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s DROP ALLOWED_VALUES 'value1', 'value2'`, id.FullyQualifiedName())
	})

	t.Run("alter with unset allowed values", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TagUnset{AllowedValues: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s UNSET ALLOWED_VALUES`, id.FullyQualifiedName())
	})

	t.Run("alter with set masking policies", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TagSet{
			MaskingPolicies: &TagSetMaskingPolicies{
				MaskingPolicies: defaultMaskingPolicies(),
				Force:           Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s SET MASKING POLICY %s, MASKING POLICY %s FORCE`, id.FullyQualifiedName(), mp1ID.FullyQualifiedName(), mp2ID.FullyQualifiedName())
	})

	t.Run("alter with unset masking policies", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TagUnset{
			MaskingPolicies: &TagUnsetMaskingPolicies{
				MaskingPolicies: defaultMaskingPolicies(),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s UNSET MASKING POLICY %s, MASKING POLICY %s`, id.FullyQualifiedName(), mp1ID.FullyQualifiedName(), mp2ID.FullyQualifiedName())
	})

	t.Run("alter with set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TagSet{Comment: String("comment")}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s SET COMMENT = 'comment'`, id.FullyQualifiedName())
	})

	t.Run("alter with set propagate", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TagSet{Propagate: &TagPropagate{PropagationMethod: Pointer(TagPropagationOnDependency)}}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s SET PROPAGATE = ON_DEPENDENCY`, id.FullyQualifiedName())
	})

	t.Run("alter with set propagate and on_conflict value", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TagSet{
			Propagate: &TagPropagate{
				PropagationMethod: Pointer(TagPropagationOnDependencyAndDataMovement),
				OnConflict:        &TagOnConflict{CustomValue: String("FAIL")},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s SET PROPAGATE = ON_DEPENDENCY_AND_DATA_MOVEMENT ON_CONFLICT = 'FAIL'`, id.FullyQualifiedName())
	})

	t.Run("alter with unset propagate", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TagUnset{Propagate: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s UNSET PROPAGATE`, id.FullyQualifiedName())
	})

	t.Run("alter with unset on_conflict", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TagUnset{OnConflict: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s UNSET ON_CONFLICT`, id.FullyQualifiedName())
	})

	t.Run("alter with unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TagUnset{Comment: Bool(true)}
		assertOptsValidAndSQLEquals(t, opts, `ALTER TAG %s UNSET COMMENT`, id.FullyQualifiedName())
	})

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*alterTagOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no alter action", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("alterTagOptions", "Add", "Drop", "Set", "Unset", "Rename"))
	})

	t.Run("validation: multiple alter actions", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TagSet{
			Comment: String("comment"),
		}
		opts.Unset = &TagUnset{
			AllowedValues: Bool(true),
		}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("alterTagOptions", "Add", "Drop", "Set", "Unset", "Rename"))
	})

	t.Run("validation: multiple fields in set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TagSet{
			Comment:         String("comment"),
			MaskingPolicies: &TagSetMaskingPolicies{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("TagSet", "MaskingPolicies", "AllowedValues", "Propagate", "Comment"))
	})

	t.Run("validation: empty masking policies in set", func(t *testing.T) {
		opts := defaultOpts()
		opts.Set = &TagSet{
			MaskingPolicies: &TagSetMaskingPolicies{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("TagSet.MaskingPolicies", "MaskingPolicies", IntErrGreater, 0))
	})

	t.Run("validation: empty masking policies in unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TagUnset{
			MaskingPolicies: &TagUnsetMaskingPolicies{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntValue("TagUnset.MaskingPolicies", "MaskingPolicies", IntErrGreater, 0))
	})

	t.Run("validation: invalid new name", func(t *testing.T) {
		opts := defaultOpts()
		opts.Rename = &TagRename{
			Name: emptySchemaObjectIdentifier,
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: new name from different db", func(t *testing.T) {
		newId := randomSchemaObjectIdentifier()

		opts := defaultOpts()
		opts.Rename = &TagRename{
			Name: newId,
		}
		assertOptsValid(t, opts)
	})

	t.Run("validation: no property to unset", func(t *testing.T) {
		opts := defaultOpts()
		opts.Unset = &TagUnset{}
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("TagUnset", "MaskingPolicies", "AllowedValues", "Propagate", "OnConflict", "Comment"))
	})

	t.Run("validation: add allowed values count", func(t *testing.T) {
		opts := defaultOpts()
		opts.Add = &TagAdd{
			AllowedValues: &AllowedValues{
				Values: []StringAllowEmpty{},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntBetween("AllowedValues", "Values", 1, 300))
	})

	t.Run("validation: drop allowed values count", func(t *testing.T) {
		opts := defaultOpts()
		opts.Drop = &TagDrop{
			AllowedValues: &AllowedValues{
				Values: []StringAllowEmpty{},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errIntBetween("AllowedValues", "Values", 1, 300))
	})
}

func TestTagSet(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *setTagOptions {
		return &setTagOptions{
			objectType: ObjectTypeStage,
			objectName: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*setTagOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.objectName = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: unsupported object type", func(t *testing.T) {
		opts := defaultOpts()
		opts.objectType = ObjectTypeSequence
		assertOptsInvalidJoinedErrors(t, opts, errors.New("tagging for object type SEQUENCE is not supported"))
	})

	t.Run("set with all optional", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetTags = []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER %s %s SET TAG "tag1" = 'value1'`, opts.objectType, id.FullyQualifiedName())
	})

	t.Run("set on account", func(t *testing.T) {
		accountId := randomAccountIdentifier()
		opts := &setTagOptions{
			objectType: ObjectTypeStage,
			objectName: accountId,
			SetTags: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag1"),
					Value: "value1",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER %s %s SET TAG "tag1" = 'value1'`, opts.objectType, accountId.FullyQualifiedName())
	})

	t.Run("set with column", func(t *testing.T) {
		objectId := randomTableColumnIdentifierInSchemaObject(id)
		tagId := randomSchemaObjectIdentifier()
		request := NewSetTagRequest(ObjectTypeColumn, objectId).WithSetTags([]TagAssociation{
			{
				Name:  tagId,
				Value: "value1",
			},
		})
		opts := request.toOpts()
		assertOptsValidAndSQLEquals(t, opts, `ALTER TABLE %s MODIFY COLUMN "%s" SET TAG %s = 'value1'`, id.FullyQualifiedName(), objectId.columnName, tagId.FullyQualifiedName())
	})
}

func TestTagUnset(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	defaultOpts := func() *unsetTagOptions {
		return &unsetTagOptions{
			objectType: ObjectTypeStage,
			objectName: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		opts := (*unsetTagOptions)(nil)
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.objectName = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: unsupported object type", func(t *testing.T) {
		opts := defaultOpts()
		opts.objectType = ObjectTypeSequence
		assertOptsInvalidJoinedErrors(t, opts, errors.New("tagging for object type SEQUENCE is not supported"))
	})

	t.Run("unset with all optional", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetTags = []ObjectIdentifier{
			NewAccountObjectIdentifier("tag1"),
			NewAccountObjectIdentifier("tag2"),
		}
		opts.IfExists = Pointer(true)
		assertOptsValidAndSQLEquals(t, opts, `ALTER %s IF EXISTS %s UNSET TAG "tag1", "tag2"`, opts.objectType, id.FullyQualifiedName())
	})

	t.Run("unset on account", func(t *testing.T) {
		accountId := randomAccountIdentifier()
		opts := &unsetTagOptions{
			objectType: ObjectTypeStage,
			objectName: accountId,
			UnsetTags: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag1"),
				NewAccountObjectIdentifier("tag2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER %s %s UNSET TAG "tag1", "tag2"`, opts.objectType, accountId.FullyQualifiedName())
	})

	t.Run("unset with column", func(t *testing.T) {
		objectId := randomTableColumnIdentifierInSchemaObject(id)
		tagId1 := randomSchemaObjectIdentifier()
		tagId2 := randomSchemaObjectIdentifierInSchema(tagId1.SchemaId())
		request := UnsetTagRequest{
			objectType: ObjectTypeColumn,
			objectName: objectId,
			UnsetTags: []ObjectIdentifier{
				tagId1,
				tagId2,
			},
			IfExists: Pointer(true),
		}
		opts := request.toOpts()
		assertOptsValidAndSQLEquals(t, opts, `ALTER %s IF EXISTS %s MODIFY COLUMN "%s" UNSET TAG %s, %s`, opts.objectType, id.FullyQualifiedName(), objectId.Name(), tagId1.FullyQualifiedName(), tagId2.FullyQualifiedName())
	})
}

func Test_Tag_ToTagPropagation(t *testing.T) {
	type test struct {
		input string
		want  TagPropagation
	}

	valid := []test{
		// case insensitive
		{input: "on_dependency", want: TagPropagationOnDependency},

		// Supported Values
		{input: "ON_DEPENDENCY_AND_DATA_MOVEMENT", want: TagPropagationOnDependencyAndDataMovement},
		{input: "ON_DEPENDENCY", want: TagPropagationOnDependency},
		{input: "ON_DATA_MOVEMENT", want: TagPropagationOnDataMovement},
		{input: "NONE", want: TagPropagationNone},
	}

	invalid := []test{
		{input: ""},
		{input: "foo"},
		{input: "'ON_DEPENDENCY'"},
	}

	for _, tc := range valid {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ToTagPropagation(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}

	for _, tc := range invalid {
		t.Run(tc.input, func(t *testing.T) {
			_, err := ToTagPropagation(tc.input)
			require.Error(t, err)
		})
	}
}
