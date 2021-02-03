package utils

import (
	"reflect"
	"strings"
)

const TagName = "mapKey"

func Map(s interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	FillMap(s, out)
	return out
}

func FillMap(s interface{}, out map[string]interface{}) {
	if out == nil {
		return
	}

	srcValue := reflect.ValueOf(s)
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}
	fields := structFields(srcValue)

	for _, field := range fields {
		name := field.Name
		val := srcValue.FieldByName(name)

		tagName, tagOpts := parseTag(field.Tag.Get(TagName))
		if tagName == "ignore" || tagOpts.Has("ignore") {
			continue
		}
		if tagName != "" {
			name = tagName
		}

		// if the value is a zero value and the field is marked as omitempty do
		// not include
		if tagOpts.Has("omitempty") {
			zero := reflect.Zero(val.Type()).Interface()
			current := val.Interface()

			if reflect.DeepEqual(current, zero) {
				continue
			}
		}

		out[name] = val.Interface()
	}
}

func structFields(value reflect.Value) []reflect.StructField {
	t := value.Type()
	var f []reflect.StructField

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// we can't access the value of unexported fields
		if field.PkgPath != "" {
			continue
		}

		// don't check if it's omitted
		if tag := field.Tag.Get(TagName); tag == "-" {
			continue
		}

		f = append(f, field)
	}

	return f
}

type tagOptions []string

func (t tagOptions) Has(opt string) bool {
	for _, tagOpt := range t {
		if tagOpt == opt {
			return true
		}
	}

	return false
}

func parseTag(tag string) (string, tagOptions) {
	res := strings.Split(tag, ",")
	return res[0], res[1:]
}
