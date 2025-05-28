package generator

func (v *QueryStruct) OptionalSQL(sql string) *QueryStruct {
	v.fields = append(v.fields, NewField(sqlToFieldName(sql, true), "*bool", Tags().Keyword().SQL(sql), nil))
	return v
}

func (v *QueryStruct) NamedList(sql string, itemKind string, transformer *KeywordTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(sqlToFieldName(sql, true), KindOfSlice(itemKind), Tags().Keyword().SQL(sql), transformer))
	return v
}

func (v *QueryStruct) OrReplace() *QueryStruct {
	return v.OptionalSQL("OR REPLACE")
}

func (v *QueryStruct) IfNotExists() *QueryStruct {
	return v.OptionalSQL("IF NOT EXISTS")
}

func (v *QueryStruct) IfExists() *QueryStruct {
	return v.OptionalSQL("IF EXISTS")
}

func (v *QueryStruct) Terse() *QueryStruct {
	return v.OptionalSQL("TERSE")
}

func (v *QueryStruct) Text(name string, transformer *KeywordTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, "string", Tags().Keyword(), transformer))
	return v
}

func (v *QueryStruct) Any(name string, transformer *KeywordTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, "any", Tags().Keyword(), transformer))
	return v
}

func (v *QueryStruct) Number(name string, transformer *KeywordTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, "int", Tags().Keyword(), transformer))
	return v
}

func (v *QueryStruct) OptionalText(name string, transformer *KeywordTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, "*string", Tags().Keyword(), transformer))
	return v
}

func (v *QueryStruct) OptionalNumber(name string, transformer *KeywordTransformer) *QueryStruct {
	v.fields = append(v.fields, NewField(name, "*int", Tags().Keyword(), transformer))
	return v
}

func (v *QueryStruct) OptionalLimitFrom() *QueryStruct {
	return v.PredefinedQueryStructField("Limit", "*LimitFrom", KeywordOptions().SQL("LIMIT"))
}

func (v *QueryStruct) OptionalSessionParameters() *QueryStruct {
	return v.PredefinedQueryStructField("SessionParameters", "*SessionParameters", ListOptions().NoParentheses()).
		WithValidation(ValidateValue, "SessionParameters")
}

func (v *QueryStruct) OptionalSessionParametersUnset() *QueryStruct {
	return v.PredefinedQueryStructField("SessionParametersUnset", "*SessionParametersUnset", ListOptions().NoParentheses()).
		WithValidation(ValidateValue, "SessionParametersUnset")
}

func (v *QueryStruct) NamedListWithParens(sqlPrefix string, listItemKind string, transformer *KeywordTransformer) *QueryStruct {
	if transformer != nil {
		transformer = transformer.Parentheses().SQL(sqlPrefix)
	} else {
		transformer = KeywordOptions().Parentheses().SQL(sqlPrefix)
	}
	v.fields = append(v.fields, NewField(sqlToFieldName(sqlPrefix, true), KindOfSlice(listItemKind), Tags().Keyword(), transformer))
	return v
}

func (v *QueryStruct) OptionalTags() *QueryStruct {
	return v.NamedListWithParens("TAG", "TagAssociation", nil)
}

func (v *QueryStruct) SetTags() *QueryStruct {
	return v.setTags(KeywordOptions().Required())
}

func (v *QueryStruct) OptionalSetTags() *QueryStruct {
	return v.setTags(KeywordOptions().SQL("SET TAG"))
}

func (v *QueryStruct) setTags(transformer *KeywordTransformer) *QueryStruct {
	return v.PredefinedQueryStructField("SetTags", "[]TagAssociation", transformer)
}

func (v *QueryStruct) UnsetTags() *QueryStruct {
	return v.unsetTags(KeywordOptions().Required())
}

func (v *QueryStruct) OptionalUnsetTags() *QueryStruct {
	return v.unsetTags(KeywordOptions().SQL("UNSET TAG"))
}

func (v *QueryStruct) unsetTags(transformer *KeywordTransformer) *QueryStruct {
	return v.PredefinedQueryStructField("UnsetTags", "[]ObjectIdentifier", transformer)
}

func (v *QueryStruct) OptionalLike() *QueryStruct {
	return v.PredefinedQueryStructField("Like", "*Like", KeywordOptions().SQL("LIKE"))
}

func (v *QueryStruct) OptionalIn() *QueryStruct {
	return v.PredefinedQueryStructField("In", "*In", KeywordOptions().SQL("IN"))
}

func (v *QueryStruct) OptionalExtendedIn() *QueryStruct {
	return v.PredefinedQueryStructField("In", "*ExtendedIn", KeywordOptions().SQL("IN"))
}

func (v *QueryStruct) OptionalServiceIn() *QueryStruct {
	return v.PredefinedQueryStructField("In", "*ServiceIn", KeywordOptions().SQL("IN"))
}

func (v *QueryStruct) OptionalStartsWith() *QueryStruct {
	return v.PredefinedQueryStructField("StartsWith", "*string", ParameterOptions().NoEquals().SingleQuotes().SQL("STARTS WITH"))
}

func (v *QueryStruct) OptionalLimit() *QueryStruct {
	return v.PredefinedQueryStructField("Limit", "*LimitFrom", KeywordOptions().SQL("LIMIT"))
}

func (v *QueryStruct) OptionalCopyGrants() *QueryStruct {
	return v.OptionalSQL("COPY GRANTS")
}

func (v *QueryStruct) BodyWithPrecedingArrow() *QueryStruct {
	return v.PredefinedQueryStructField("body", "string", ParameterOptions().NoEquals().NoQuotes().SQL("->").Required())
}

func (v *QueryStruct) OptionalSetBodyWithPrecedingArrow() *QueryStruct {
	return v.PredefinedQueryStructField("SetBody", "*string", ParameterOptions().NoEquals().NoQuotes().SQL("SET BODY ->"))
}
