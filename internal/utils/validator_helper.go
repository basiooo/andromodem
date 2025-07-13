package utils

import (
	"reflect"
	"strings"
)

func GetJSONFieldName(fld reflect.StructField) string {
	splited := strings.SplitN(fld.Tag.Get("json"), ",", 2)
	if len(splited) == 0 {
		return ""
	}
	name := splited[0]
	if name == "-" || name == "" {
		return ""
	}
	return name
}
