package gen

import (
	"log"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/genhelpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
)

type DataSourceModel struct {
	Name       string
	NamePlural string
	// ObjectType is the object words of SHOW PARAMETERS FOR. ShowSql, ShowDoc and DescribeSql come from the SDK config.
	ObjectType          string
	ShowSql             string
	ShowDoc             string
	DescribeSql         string
	DescribeDetailsName string
	PreviewFeature      bool
	// ParametersAcrossLevels is set when the parameters of the data source cover the levels of more than one variant.
	ParametersAcrossLevels bool
	Filters                ssot.DataSourceFilters
	Describe               ssot.DataSourceDescribe
	Hooks                  ssot.DataSourceHooks
	Additional             ssot.DataSourceGenerator
	*genhelpers.PreambleModel
}

func (m DataSourceModel) Prefix() string {
	return genhelpers.FirstLetterLowercase(m.NamePlural)
}

// ResultName is the attribute that holds the list of results.
func (m DataSourceModel) ResultName() string {
	return genhelpers.ToSnakeCase(m.NamePlural)
}

func (m DataSourceModel) DescribeSchemaName() string {
	if m.DescribeDetailsName != "" {
		return m.DescribeDetailsName
	}
	return m.Name + "Details"
}

// ResultSources lists the commands whose output the result holds, for the description.
func (m DataSourceModel) ResultSources() string {
	sources := []string{"SHOW"}
	if m.Hooks.DescribeOutputSchema.Present {
		sources = append(sources, "DESCRIBE")
	}
	if m.Hooks.ParametersOutputSchema.Present {
		sources = append(sources, "SHOW PARAMETERS")
	}
	switch len(sources) {
	case 1:
		return sources[0]
	case 2:
		return sources[0] + " and " + sources[1]
	default:
		return strings.Join(sources[:len(sources)-1], ", ") + ", and " + sources[len(sources)-1]
	}
}

// ResultDescription is the description of the result list. A list that holds only SHOW output has no details to aggregate.
func (m DataSourceModel) ResultDescription() string {
	if !m.Hooks.DescribeOutputSchema.Present && !m.Hooks.ParametersOutputSchema.Present {
		return "Holds the output of SHOW " + m.ShowSql + "."
	}
	return "Holds the aggregated output of all " + genhelpers.CamelToWords(m.Name) + " details queries."
}

func (m DataSourceModel) ShowOutputDescription() string {
	return joinDescription("Holds the output of SHOW "+m.ShowSql+".", m.Additional.AdditionalShowOutputDescription)
}

func (m DataSourceModel) DescribeOutputDescription() string {
	return joinDescription("Holds the output of DESCRIBE "+m.DescribeSql+".", m.Additional.AdditionalDescribeOutputDescription)
}

// ParametersSchemaExpr is the schema of the parameters block. Across levels it is the hand-written union of all of them.
func (m DataSourceModel) ParametersSchemaExpr() string {
	if m.ParametersAcrossLevels {
		return "schemas.ShowAll" + m.Name + "ParametersSchema()"
	}
	return "schemas.Show" + m.Name + "ParametersSchema"
}

func (m DataSourceModel) ParametersMapperName() string {
	if m.ParametersAcrossLevels {
		return "All" + m.Name + "ParametersToSchema"
	}
	return m.Name + "ParametersToSchema"
}

func (m DataSourceModel) ParametersOutputDescription() string {
	return joinDescription("Holds the output of SHOW PARAMETERS FOR "+m.ObjectType+".", m.Additional.AdditionalParametersOutputDescription)
}

func joinDescription(description, additional string) string {
	if additional == "" {
		return description
	}
	return description + " " + additional
}

func (m DataSourceModel) HasFilters() bool {
	f := m.Filters
	return f.Like || f.StartsWith || f.Limit || f.On || f.Terse || f.History || f.In != ssot.DataSourceInNone
}

func ModelFromInputObject(input DataSourceObject, preamble *genhelpers.PreambleModel, _ *log.Logger) DataSourceModel {
	sdkConfig := input.SdkConfig()
	return DataSourceModel{
		Name:                   input.Name,
		NamePlural:             input.NamePlural,
		ObjectType:             input.ObjectType(),
		ShowSql:                sdkConfig.ShowSql,
		ShowDoc:                sdkConfig.ShowDoc,
		DescribeSql:            sdkConfig.DescribeSql,
		DescribeDetailsName:    input.DescribeDetailsName,
		PreviewFeature:         input.PreviewFeature,
		ParametersAcrossLevels: len(input.ShowParametersDetailsAdditionalLevels) > 0,
		Filters:                input.DataSourceFilters(),
		Describe:               input.DataSourceDescribe(),
		Hooks:                  input.EffectiveDataSourceHooks(),
		Additional:             ssot.ObjectGenerators(input.Generators...).DataSource,
		PreambleModel:          preamble,
	}
}
