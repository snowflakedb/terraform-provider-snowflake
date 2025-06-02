package genhelpers

import (
	"reflect"
	"runtime"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [SNOW-1501905]: describe all methods in this file
// TODO [SNOW-1501905]: test all methods in this file

func FirstLetterLowercase(in string) string {
	return strings.ToLower(in[:1]) + in[1:]
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

func SnakeCase(name string) string {
	return ToSnakeCase(name)
}

func CamelToWords(camel string) string {
	return strings.ReplaceAll(ToSnakeCase(camel), "_", " ")
}

func SnakeCaseToCamel(snake string) string {
	var suffix string
	if strings.HasSuffix(snake, "_") {
		suffix = "_"
		snake = strings.TrimSuffix(snake, "_")
	}
	snake = strings.ToLower(snake)
	parts := strings.Split(snake, "_")
	for idx, p := range parts {
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
