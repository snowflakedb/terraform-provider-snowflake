package genhelpers

import (
	"fmt"
	"io/fs"
	"path"
	"reflect"
	"runtime"
	"strings"
	"text/template"

	_ "embed"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-1501905]: describe all methods in this file
// TODO [SNOW-1501905]: test all methods in this file

// common templates should go here
var (
	//go:embed templates/preamble.tmpl
	preambleTemplateContent string
	PreambleTemplate, _     = template.New("preambleTemplate").Parse(preambleTemplateContent)
)

func FirstLetterLowercase(in string) string {
	return strings.ToLower(in[:1]) + in[1:]
}

// FirstLetterLowercaseParameterName lowercases the first letter like FirstLetterLowercase and appends the forbidden
// attribute name suffix when the result is a Go keyword that cannot be a parameter name (e.g. Type -> type_).
func FirstLetterLowercaseParameterName(in string) string {
	return SanitizeAttributeName(FirstLetterLowercase(in))
}

func FirstLetterUppercase(in string) string {
	return strings.ToUpper(in[:1]) + in[1:]
}

func FirstLetter(in string) string {
	return in[:1]
}

func RunMapper(mapper Mapper, in ...string) string {
	return mapper(strings.Join(in, ""))
}

func TypeWithoutPointer(t string) string {
	without, _ := strings.CutPrefix(t, "*")
	return without
}

func TypeWithoutPointerAndBrackets(t string) string {
	without := strings.TrimLeft(t, "*[]")
	return without
}

func IsTypeSlice(t string) bool {
	return strings.HasPrefix(t, "[]")
}

func SnakeCase(name string) string {
	return ToSnakeCase(name)
}

func CamelToWords(camel string) string {
	return strings.ReplaceAll(ToSnakeCase(camel), "_", " ")
}

// SnakeCaseToCamel converts a snake_case string to a camelCase string.
func SnakeCaseToCamel(snake string) string {
	var suffix string
	if strings.HasSuffix(snake, "_") {
		suffix = "_"
		snake = strings.TrimSuffix(snake, "_")
	}
	snake = strings.ToLower(snake)
	parts := strings.Split(snake, "_")
	for idx, p := range parts {
		if p == "" {
			continue
		}
		parts[idx] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "") + suffix
}

func RemoveForbiddenAttributeNameSuffix(input string) string {
	return strings.TrimRight(input, forbiddenAttributeNameSuffix)
}

func ShouldGenerateWithForAttributeType(valueType schema.ValueType) bool {
	switch valueType {
	case schema.TypeBool, schema.TypeInt, schema.TypeFloat, schema.TypeString:
		return true
	default:
		return false
	}
}

func IsLastItem(itemIdx int, collectionLength int) bool {
	return itemIdx+1 == collectionLength
}

// Unimplemented panics at generate time.
// Call it from a template branch that is not implemented yet, e.g. {{ Unimplemented "ParseIdFromConfig" }}.
func Unimplemented(name string) string {
	panic(fmt.Sprintf("generated implementation for %s is not implemented", name))
}

func BuildTemplateFuncMap(funcs ...any) template.FuncMap {
	allFuncs := make(map[string]any)
	for _, f := range funcs {
		allFuncs[getFunctionName(f)] = f
	}
	return allFuncs
}

func getFunctionName(f any) string {
	fullFuncName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	parts := strings.Split(fullFuncName, ".")
	return parts[len(parts)-1]
}

// ParseSubTemplates adds every .tmpl file of dir in fsys to root. A file is named after its snake_case name in
// lowerCamelCase: parse_id.tmpl is the template "parseId", which a parent template calls with {{ template "parseId" . }}.
// It panics when a file cannot be read or parsed.
func ParseSubTemplates(root *template.Template, fsys fs.FS, dir string) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		panic(fmt.Sprintf("read sub templates of %s: %v", dir, err))
	}
	for _, entry := range entries {
		file := entry.Name()
		if entry.IsDir() || path.Ext(file) != ".tmpl" {
			continue
		}
		content, err := fs.ReadFile(fsys, path.Join(dir, file))
		if err != nil {
			panic(fmt.Sprintf("read sub template %s: %v", file, err))
		}
		template.Must(root.New(SubTemplateName(file)).Parse(string(content)))
	}
}

// SubTemplateName is the name of the sub template in the file: its snake_case name in lowerCamelCase, without the extension.
func SubTemplateName(file string) string {
	return FirstLetterLowercase(SnakeCaseToCamel(strings.TrimSuffix(file, path.Ext(file))))
}
