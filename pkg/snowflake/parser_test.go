package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestViewSelectStatementExtractor_Extract(t *testing.T) {
	basic := "create view foo as select * from bar;"
	caps := "CREATE VIEW FOO AS SELECT * FROM BAR;"
	commentWithSingleQuotes := "CREATE VIEW FOO COMMENT = 'test''' AS SELECT * FROM BAR;"
	parens := "create view foo as (select * from bar);"
	multiline := `
create view foo as
select *
from bar;`

	multilineComment := `
create view foo as
-- comment
select *
from bar;`

	secure := "create secure view foo as select * from bar;"
	replace := "create or replace view foo as select * from bar;"
	grants := "create or replace view foo copy grants as select * from bar;"
	recursive := "create recursive view foo as select * from bar;"
	ine := "create view if not exists foo as select * from bar;"

	comment := `create view foo comment='asdf' as select * from bar;`
	commentEscape := `create view foo comment='asdf\'s are fun' as select * from bar;`
	identifier := `create view "foo"."bar"."bam" comment='asdf\'s are fun' as select * from bar;`

	full := `CREATE SECURE VIEW "rgdxfmnfhh"."PUBLIC"."rgdxfmnfhh" COMMENT = 'Terraform test resource' AS SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`
	issue2640 := `CREATE OR REPLACE SECURE VIEW "CLASSIFICATION" comment = 'Classification View of the union of classification tables' AS select * from AB1_SUBSCRIPTION.CLASSIFICATION.CLASSIFICATION    union   select * from AB2_SUBSCRIPTION.CLASSIFICATION.CLASSIFICATION`
	withRowAccessAndAggregationPolicy := `CREATE SECURE VIEW "rgdxfmnfhh"."PUBLIC"."rgdxfmnfhh" COMMENT = 'Terraform test resource' ROW ACCESS policy rap on (title, title2) AGGREGATION POLICY rap   AS SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`
	withRowAccessAndAggregationPolicyWithEntityKey := `CREATE SECURE VIEW "rgdxfmnfhh"."PUBLIC"."rgdxfmnfhh" COMMENT = 'Terraform test resource' ROW ACCESS policy rap on (title, title2) AGGREGATION POLICY rap ENTITY KEY (foo, bar)  AS SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`
	columnsListEndingWithMaskingPolicyWithoutUsing := `CREATE OR REPLACE SECURE TEMPORARY VIEW "rgdxfmnfhh"."PUBLIC"."rgdxfmnfhh" (id PROJECTION POLICY pp MASKING POLICY mp COMMENT 'a (s) df', foo MASKING POLICY pp) COMMENT = 'Terraform test resource' ROW ACCESS policy rap on (title, title2) AGGREGATION POLICY rap ENTITY KEY (foo, bar)  AS SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`
	columnsListEndingWithMaskingPolicyWithUsing := `CREATE OR REPLACE SECURE TEMPORARY VIEW "rgdxfmnfhh"."PUBLIC"."rgdxfmnfhh" (id PROJECTION POLICY pp MASKING POLICY mp COMMENT 'a (s) df', foo MASKING POLICY pp USING ("col1")) COMMENT = 'Terraform test resource' ROW ACCESS policy rap on (title, title2) AGGREGATION POLICY rap ENTITY KEY (foo, bar)  AS SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`
	columnsListEndingWithProjectionPolicy := `CREATE OR REPLACE SECURE TEMPORARY VIEW "rgdxfmnfhh"."PUBLIC"."rgdxfmnfhh" (id PROJECTION POLICY pp MASKING POLICY mp COMMENT 'a (s) df' , foo PROJECTION POLICY pp) COMMENT = 'Terraform test resource' ROW ACCESS policy rap on (title, title2) AGGREGATION POLICY rap ENTITY KEY (foo, bar)  AS SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`
	columnsListEndingWithComment := `CREATE OR REPLACE SECURE TEMPORARY VIEW "rgdxfmnfhh"."PUBLIC"."rgdxfmnfhh" (id PROJECTION POLICY pp MASKING POLICY mp COMMENT 'asdf', foo PROJECTION POLICY pp COMMENT 'foo (bar) hoge') COMMENT = 'Terraform test resource' ROW ACCESS policy rap on (title, title2) AGGREGATION POLICY rap ENTITY KEY (foo, bar)  AS SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`
	columnsListEndingWithID := `CREATE OR REPLACE SECURE TEMPORARY VIEW "rgdxfmnfhh"."PUBLIC"."rgdxfmnfhh" ("ID", "FOO") COMMENT = 'Terraform test resource' ROW ACCESS policy rap on (title, title2) AGGREGATION POLICY rap ENTITY KEY (foo, bar)  AS SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`
	allFields := `CREATE OR REPLACE SECURE TEMPORARY VIEW "rgdxfmnfhh"."PUBLIC"."rgdxfmnfhh" (id PROJECTION POLICY pp MASKING POLICY mp USING ("col1", "cond1") COMMENT 'asdf', foo MASKING POLICY mp USING ("col1", "cond1")) COMMENT = 'Terraform test resource' ROW ACCESS policy rap on (title, title2) AGGREGATION POLICY rap ENTITY KEY (foo, bar)  AS SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`
	// Regression test for SNOW-3308280: panic when consumeToken reaches end of input (off-by-one in bounds check).
	// The last column has a masking policy and no trailing options, so the parser hits end-of-input while probing optional tokens.
	snow3308280 := `CREATE OR REPLACE VIEW "DB"."SCHEMA"."VIEW_NAME" (col1 MASKING POLICY "DB"."SCHEMA"."MASK_POLICY") AS SELECT col1 FROM DB.SCHEMA.SOURCE_TABLE`
	testStatement := "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"basic", args{basic}, "select * from bar;", false},
		{"caps", args{caps}, "SELECT * FROM BAR;", false},
		{"comment with single quotes", args{commentWithSingleQuotes}, "SELECT * FROM BAR;", false},
		{"parens", args{parens}, "(select * from bar);", false},
		{"multiline", args{multiline}, "select *\nfrom bar;", false},
		{"multilineComment", args{multilineComment}, "-- comment\nselect *\nfrom bar;", false},
		{"secure", args{secure}, "select * from bar;", false},
		{"replace", args{replace}, "select * from bar;", false},
		{"grants", args{grants}, "select * from bar;", false},
		{"recursive", args{recursive}, "select * from bar;", false},
		{"ine", args{ine}, "select * from bar;", false},
		{"comment", args{comment}, "select * from bar;", false},
		{"commentEscape", args{commentEscape}, "select * from bar;", false},
		{"identifier", args{identifier}, "select * from bar;", false},
		{"full", args{full}, testStatement, false},
		{"issue2640", args{issue2640}, "select * from AB1_SUBSCRIPTION.CLASSIFICATION.CLASSIFICATION    union   select * from AB2_SUBSCRIPTION.CLASSIFICATION.CLASSIFICATION", false},
		{"with row access policy and aggregation policy", args{withRowAccessAndAggregationPolicy}, testStatement, false},
		{"with row access policy and aggregation policy with entity key", args{withRowAccessAndAggregationPolicyWithEntityKey}, testStatement, false},
		{"with column list ending with masking policy without using", args{columnsListEndingWithMaskingPolicyWithoutUsing}, testStatement, false},
		{"with column list ending with masking policy with using", args{columnsListEndingWithMaskingPolicyWithUsing}, testStatement, false},
		{"with column list ending with projection using", args{columnsListEndingWithProjectionPolicy}, testStatement, false},
		{"with column list ending with comment", args{columnsListEndingWithComment}, testStatement, false},
		{"with column list ending with column name", args{columnsListEndingWithID}, testStatement, false},
		{"all fields", args{allFields}, testStatement, false},
		{"SNOW-3308280: masking policy on last column reaching end of input", args{snow3308280}, "SELECT col1 FROM DB.SCHEMA.SOURCE_TABLE", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := NewViewSelectStatementExtractor(tt.args.input)
			got, err := e.Extract()
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewSelectStatementExtractor.Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestViewSelectStatementExtractor_ExtractMaterializedView(t *testing.T) {
	basic := "create materialized view foo as select * from bar;"
	caps := "CREATE MATERIALIZED VIEW FOO AS SELECT * FROM BAR;"
	parens := "create materialized view foo as (select * from bar);"
	multiline := `
create materialized view foo as
select *
from bar;`

	multilineComment := `
create materialized view foo as
-- comment
select *
from bar;`

	secure := "create secure materialized view foo as select * from bar;"
	replace := "create or replace materialized view foo as select * from bar;"
	ine := "create materialized view if not exists foo as select * from bar;"

	comment := `create materialized view foo comment='asdf' as select * from bar;`
	commentEscape := `create materialized view foo comment='asdf\'s are fun' as select * from bar;`
	clusterBy := "create materialized view foo cluster by (c1, c2) as select * from bar;"
	identifier := `create materialized view "foo"."bar"."bam" comment='asdf\'s are fun' as select * from bar;`

	full := `CREATE SECURE MATERIALIZED VIEW "rgdxfmnfhh"."PUBLIC"."rgdxfmnfhh" COMMENT = 'Terraform test resource' CLUSTER BY (C1, C2) AS SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES`

	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"basic", args{basic}, "select * from bar;", false},
		{"caps", args{caps}, "SELECT * FROM BAR;", false},
		{"parens", args{parens}, "(select * from bar);", false},
		{"multiline", args{multiline}, "select *\nfrom bar;", false},
		{"multilineComment", args{multilineComment}, "-- comment\nselect *\nfrom bar;", false},
		{"secure", args{secure}, "select * from bar;", false},
		{"replace", args{replace}, "select * from bar;", false},
		{"ine", args{ine}, "select * from bar;", false},
		{"comment", args{comment}, "select * from bar;", false},
		{"commentEscape", args{commentEscape}, "select * from bar;", false},
		{"clusterBy", args{clusterBy}, "select * from bar;", false},
		{"identifier", args{identifier}, "select * from bar;", false},
		{"full", args{full}, "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := NewViewSelectStatementExtractor(tt.args.input)
			got, err := e.ExtractMaterializedView()
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewSelectStatementExtractor.ExtractMaterializedView() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ViewSelectStatementExtractor.ExtractMaterializedView() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestViewSelectStatementExtractor_ExtractDynamicTable(t *testing.T) {
	basic := "create dynamic table foo lag = 'DOWNSTREAM' refresh_mode = 'AUTO' initialize = 'ON_CREATE' warehouse = COMPUTE_WH as select * from bar;"
	caps := "CREATE DYNAMIC TABLE FOO LAG = 'DOWNSTREAM' REFRESH_MODE = 'AUTO' INITIALIZE = 'ON_CREATE' WAREHOUSE = COMPUTE_WH AS SELECT * FROM BAR;"
	parens := "create dynamic table foo lag = 'DOWNSTREAM' refresh_mode = 'AUTO' initialize = 'ON_CREATE' warehouse = COMPUTE_WH as (select * from bar);"
	multiline := `
create dynamic table foo
lag = 'DOWNSTREAM'
refresh_mode = 'AUTO'
initialize = 'ON_CREATE'
warehouse = COMPUTE_WH
as select *
from bar;`

	multilineComment := `
create dynamic table foo
lag = 'DOWNSTREAM'
refresh_mode = 'AUTO'
initialize = 'ON_CREATE'
warehouse = COMPUTE_WH
as
-- comment
select *
from bar;`

	comment := `create dynamic table foo lag = 'DOWNSTREAM' refresh_mode = 'AUTO' initialize = 'ON_CREATE' warehouse = COMPUTE_WH comment = 'asdf' as select * from bar;`
	commentEscape := `create dynamic table foo lag = 'DOWNSTREAM' refresh_mode = 'AUTO' initialize = 'ON_CREATE' warehouse = COMPUTE_WH comment = 'asdf\'s are fun' as select * from bar;`
	orReplace := `create or replace dynamic table foo lag = 'DOWNSTREAM' refresh_mode = 'AUTO' initialize = 'ON_CREATE' warehouse = COMPUTE_WH comment = 'asdf' as select * from bar;`
	identifier := `create or replace dynamic table "foo"."bar"."bam" lag = 'DOWNSTREAM' refresh_mode = 'AUTO' initialize = 'ON_CREATE' warehouse = COMPUTE_WH comment = 'asdf\'s are fun' as select * from bar;`
	// running SHOW DYNAMIC TABLE in Snowflake actually returns the query with
	// the comment before other parameters, even though this is inconsistent
	// with the order they are specified in CREATE DYNAMIC TABLE
	commentBeforeOtherParams := `create dynamic table foo comment = 'asdf\'s are fun' lag = 'DOWNSTREAM' refresh_mode = 'AUTO' initialize = 'ON_CREATE' warehouse = COMPUTE_WH as select * from bar;`

	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"basic", args{basic}, "select * from bar;", false},
		{"caps", args{caps}, "SELECT * FROM BAR;", false},
		{"parens", args{parens}, "(select * from bar);", false},
		{"multiline", args{multiline}, "select *\nfrom bar;", false},
		{"multilineComment", args{multilineComment}, "-- comment\nselect *\nfrom bar;", false},
		{"comment", args{comment}, "select * from bar;", false},
		{"commentEscape", args{commentEscape}, "select * from bar;", false},
		{"orReplace", args{orReplace}, "select * from bar;", false},
		{"identifier", args{identifier}, "select * from bar;", false},
		{"commentBeforeOtherParams", args{commentBeforeOtherParams}, "select * from bar;", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := NewViewSelectStatementExtractor(tt.args.input)
			got, err := e.ExtractDynamicTable()
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewSelectStatementExtractor.ExtractDynamicTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ViewSelectStatementExtractor.ExtractDynamicTable() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestViewSelectStatementExtractor_consumeToken(t *testing.T) {
	type fields struct {
		input []rune
		pos   int
	}
	type args struct {
		t string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		posAfter int
	}{
		{"basic - found", fields{[]rune("foo"), 0}, args{"foo"}, 3},
		{"basic - not found", fields{[]rune("foo"), 0}, args{"bar"}, 0},
		{"basic - not found", fields{[]rune("fob"), 0}, args{"foo"}, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &ViewSelectStatementExtractor{
				input: tt.fields.input,
				pos:   tt.fields.pos,
			}
			e.consumeToken(tt.args.t)

			if e.pos != tt.posAfter {
				t.Errorf("pos after = %v, want %v", e.pos, tt.posAfter)
			}
		})
	}
}

func TestViewSelectStatementExtractor_consumeSpace(t *testing.T) {
	type fields struct {
		input []rune
		pos   int
	}
	tests := []struct {
		name     string
		fields   fields
		posAfter int
	}{
		{"simple", fields{[]rune("   foo"), 0}, 3},
		{"empty", fields{[]rune(""), 0}, 0},
		{"middle", fields{[]rune("foo \t\n bar"), 3}, 7},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &ViewSelectStatementExtractor{
				input: tt.fields.input,
				pos:   tt.fields.pos,
			}
			e.consumeSpace()

			if e.pos != tt.posAfter {
				t.Errorf("pos after = %v, want %v", e.pos, tt.posAfter)
			}
		})
	}
}

func TestViewSelectStatementExtractor_consumeComment(t *testing.T) {
	type fields struct {
		input []rune
		pos   int
	}
	tests := []struct {
		name     string
		fields   fields
		posAfter int
	}{
		{"basic", fields{[]rune("comment='foo'"), 0}, 13},
		{"escaped", fields{[]rune(`comment='fo\'o'`), 0}, 15},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &ViewSelectStatementExtractor{
				input: tt.fields.input,
				pos:   tt.fields.pos,
			}
			e.consumeComment()

			if e.pos != tt.posAfter {
				t.Errorf("pos after = %v, want %v", e.pos, tt.posAfter)
			}
		})
	}
}

func TestViewSelectStatementExtractor_consumeClusterBy(t *testing.T) {
	type fields struct {
		input []rune
		pos   int
	}
	tests := []struct {
		name     string
		fields   fields
		posAfter int
	}{
		{"none", fields{[]rune("as foo"), 0}, 0},
		{"single", fields{[]rune("(c1)"), 0}, 4},
		{"double", fields{[]rune("(c1, c2)"), 0}, 8},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			e := &ViewSelectStatementExtractor{
				input: tt.fields.input,
				pos:   tt.fields.pos,
			}
			e.consumeClusterBy()

			if e.pos != tt.posAfter {
				t.Errorf("pos after = %v, want %v", e.pos, tt.posAfter)
			}
		})
	}
}
