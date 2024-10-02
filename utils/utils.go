package utils

import (
	"reflect"
	"strings"
)

func FieldToColumnName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return field.Name
	}
	return strings.Split(jsonTag, ",")[0]
}
