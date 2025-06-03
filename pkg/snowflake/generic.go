package snowflake

import (
	"bytes"
	"fmt"
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

func (ab *AlterPropertiesBuilder) Statement() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`ALTER %s "%s" SET`, ab.entityType, ab.name)) // TODO handle error

	sb.WriteString(ab.rawStatement)

	for _, k := range sortStrings(ab.stringProperties) {
		sb.WriteString(fmt.Sprintf(" %s='%s'", strings.ToUpper(k), EscapeString(ab.stringProperties[k])))
	}

	for _, k := range sortStringList(ab.stringListProperties) {
		sb.WriteString(fmt.Sprintf(" %s=%s", strings.ToUpper(k), formatStringList(ab.stringListProperties[k])))
	}

	for _, k := range sortStringsBool(ab.boolProperties) {
		sb.WriteString(fmt.Sprintf(" %s=%t", strings.ToUpper(k), ab.boolProperties[k]))
	}

	for _, k := range sortStringsInt(ab.intProperties) {
		sb.WriteString(fmt.Sprintf(" %s=%d", strings.ToUpper(k), ab.intProperties[k]))
	}

	for _, k := range sortStringsFloat(ab.floatProperties) {
		sb.WriteString(fmt.Sprintf(" %s=%.2f", strings.ToUpper(k), ab.floatProperties[k]))
	}

	if len(ab.tags) > 0 {
		sb.WriteString(fmt.Sprintf(` TAG %s`, ab.GetTagValueString()))
	}
	return sb.String()
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

func (b *CreateBuilder) Statement() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`CREATE %s "%s"`, b.entityType, b.name)) // TODO handle error

	sb.WriteString(b.rawStatement)

	for _, k := range sortStrings(b.stringProperties) {
		sb.WriteString(fmt.Sprintf(" %s='%s'", strings.ToUpper(k), EscapeString(b.stringProperties[k])))
	}

	for _, k := range sortStringList(b.stringListProperties) {
		sb.WriteString(fmt.Sprintf(" %s=%s", strings.ToUpper(k), formatStringList(b.stringListProperties[k])))
	}

	for _, k := range sortStringsBool(b.boolProperties) {
		sb.WriteString(fmt.Sprintf(" %s=%t", strings.ToUpper(k), b.boolProperties[k]))
	}

	for _, k := range sortStringsInt(b.intProperties) {
		sb.WriteString(fmt.Sprintf(" %s=%d", strings.ToUpper(k), b.intProperties[k]))
	}

	for _, k := range sortStringsFloat(b.floatProperties) {
		sb.WriteString(fmt.Sprintf(" %s=%.2f", strings.ToUpper(k), b.floatProperties[k]))
	}

	if len(b.tags) > 0 {
		sb.WriteString(fmt.Sprintf(` WITH TAG (%s)`, b.GetTagValueString()))
	}

	return sb.String()
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
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
