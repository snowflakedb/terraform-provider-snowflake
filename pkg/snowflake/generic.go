package snowflake

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/template"
)

type EntityType string

type Builder struct {
	entityType EntityType
	name       string
}

func (b *Builder) Show() string {
	return fmt.Sprintf(`SHOW %sS LIKE '%s'`, b.entityType, b.name)
}

func (b *Builder) Describe() string {
	return fmt.Sprintf(`DESCRIBE %s "%s"`, b.entityType, b.name)
}

func (b *Builder) Drop() string {
	return fmt.Sprintf(`DROP %s "%s"`, b.entityType, b.name)
}

func (b *Builder) Rename(newName string) string {
	return fmt.Sprintf(`ALTER %s "%s" RENAME TO "%s"`, b.entityType, b.name, newName)
}

// SettingBuilder is an interface for a builder that allows you to set key value pairs.
type SettingBuilder interface {
	SetString(string, string)
	SetStringList(string, []string)
	SetBool(string, bool)
	SetInt(string, int)
	SetFloat(string, float64)
	SetRaw(string)
}

type AlterPropertiesBuilder struct {
	name                 string
	entityType           EntityType
	stringProperties     map[string]string
	stringListProperties map[string][]string
	boolProperties       map[string]bool
	intProperties        map[string]int
	floatProperties      map[string]float64
	rawStatement         string
	tags                 []TagValue
}

func (b *Builder) Alter() *AlterPropertiesBuilder {
	return &AlterPropertiesBuilder{
		name:                 b.name,
		entityType:           b.entityType,
		stringProperties:     make(map[string]string),
		stringListProperties: make(map[string][]string),
		boolProperties:       make(map[string]bool),
		intProperties:        make(map[string]int),
		floatProperties:      make(map[string]float64),
	}
}

func (ab *AlterPropertiesBuilder) SetString(key, value string) {
	ab.stringProperties[key] = value
}

func (ab *AlterPropertiesBuilder) SetStringList(key string, value []string) {
	ab.stringListProperties[key] = value
}

func (ab *AlterPropertiesBuilder) SetBool(key string, value bool) {
	ab.boolProperties[key] = value
}

func (ab *AlterPropertiesBuilder) SetInt(key string, value int) {
	ab.intProperties[key] = value
}

func (ab *AlterPropertiesBuilder) SetFloat(key string, value float64) {
	ab.floatProperties[key] = value
}

func (ab *AlterPropertiesBuilder) SetRaw(rawStatement string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`%s %s`, ab.rawStatement, rawStatement))
	ab.rawStatement = sb.String()
}

func (ab *AlterPropertiesBuilder) SetTags(tags []TagValue) {
	ab.tags = tags
}

func (ab *AlterPropertiesBuilder) GetTagValueString() string {
	var q strings.Builder
	for _, v := range ab.tags {
		if v.Schema != "" {
			if v.Database != "" {
				q.WriteString(fmt.Sprintf(`"%v".`, v.Database))
			}
			q.WriteString(fmt.Sprintf(`"%v".`, v.Schema))
		}
		q.WriteString(fmt.Sprintf(`"%v" = "%v", `, v.Name, v.Value))
	}
	return strings.TrimSuffix(q.String(), ", ")
}

type CreateBuilder struct {
	name                 string
	entityType           EntityType
	stringProperties     map[string]string
	stringListProperties map[string][]string
	boolProperties       map[string]bool
	intProperties        map[string]int
	floatProperties      map[string]float64
	rawStatement         string
	tags                 []TagValue
}

func (b *Builder) Create() *CreateBuilder {
	return &CreateBuilder{
		name:                 b.name,
		entityType:           b.entityType,
		stringProperties:     make(map[string]string),
		stringListProperties: make(map[string][]string),
		boolProperties:       make(map[string]bool),
		intProperties:        make(map[string]int),
		floatProperties:      make(map[string]float64),
	}
}

func (b *CreateBuilder) SetString(key, value string) {
	b.stringProperties[key] = value
}

func (b *CreateBuilder) SetStringList(key string, value []string) {
	b.stringListProperties[key] = value
}

func (b *CreateBuilder) SetBool(key string, value bool) {
	b.boolProperties[key] = value
}

func (b *CreateBuilder) SetInt(key string, value int) {
	b.intProperties[key] = value
}

func (b *CreateBuilder) SetFloat(key string, value float64) {
	b.floatProperties[key] = value
}

func (b *CreateBuilder) SetRaw(rawStatement string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`%s %s`, b.rawStatement, rawStatement))
	b.rawStatement = sb.String()
}

func (b *CreateBuilder) SetTags(tags []TagValue) {
	b.tags = tags
}

func (b *CreateBuilder) GetTagValueString() string {
	var q strings.Builder
	for _, v := range b.tags {
		if v.Schema != "" {
			if v.Database != "" {
				q.WriteString(fmt.Sprintf(`"%v".`, v.Database))
			}
			q.WriteString(fmt.Sprintf(`"%v".`, v.Schema))
		}
		q.WriteString(fmt.Sprintf(`"%v" = "%v", `, v.Name, v.Value))
	}
	return strings.TrimSuffix(q.String(), ", ")
}

func formatStringList(list []string) string {
	t, err := template.New("StringList").Funcs(template.FuncMap{
		"escapeString": EscapeString,
	}).Parse(`({{ range $i, $v := .}}{{ if $i }}, {{ end }}'{{ escapeString $v }}'{{ end }})`)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer

	if err := t.Execute(&buf, list); err != nil {
		return ""
	}

	return buf.String()
}

func Contains(s []string, str string) bool {
	return slices.Contains(s, str)
}
