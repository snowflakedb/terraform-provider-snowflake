package gen

type TagBuilder struct {
	db  []string
	ddl []string
	sql []string
}

func Tags() *TagBuilder {
	return &TagBuilder{
		db:  make([]string, 0),
		ddl: make([]string, 0),
		sql: make([]string, 0),
	}
}

func (v *TagBuilder) Static() *TagBuilder {
	v.ddl = append(v.ddl, "static")
	return v
}

func (v *TagBuilder) Keyword() *TagBuilder {
	v.ddl = append(v.ddl, "keyword")
	return v
}

func (v *TagBuilder) Parameter() *TagBuilder {
	v.ddl = append(v.ddl, "parameter")
	return v
}

func (v *TagBuilder) Identifier() *TagBuilder {
	v.ddl = append(v.ddl, "identifier")
	return v
}

func (v *TagBuilder) List() *TagBuilder {
	v.ddl = append(v.ddl, "list")
	return v
}

func (v *TagBuilder) Parentheses() *TagBuilder {
	v.ddl = append(v.ddl, "parentheses")
	return v
}

func (v *TagBuilder) NoParentheses() *TagBuilder {
	v.ddl = append(v.ddl, "no_parentheses")
	return v
}

func (v *TagBuilder) NoEquals() *TagBuilder {
	v.ddl = append(v.ddl, "no_equals")
	return v
}

func (v *TagBuilder) SingleQuotes() *TagBuilder {
	v.ddl = append(v.ddl, "single_quotes")
	return v
}

func (v *TagBuilder) DB(db ...string) *TagBuilder {
	v.db = append(v.db, db...)
	return v
}

func (v *TagBuilder) DDL(ddl ...string) *TagBuilder {
	v.ddl = append(v.ddl, ddl...)
	return v
}

func (v *TagBuilder) SQL(sql ...string) *TagBuilder {
	v.sql = append(v.sql, sql...)
	return v
}

func (v *TagBuilder) Build() map[string][]string {
	res := make(map[string][]string)
	if len(v.db) > 0 {
		res["db"] = v.db
	}
	if len(v.ddl) > 0 {
		res["ddl"] = v.ddl
	}
	if len(v.sql) > 0 {
		res["sql"] = v.sql
	}
	return res
}
