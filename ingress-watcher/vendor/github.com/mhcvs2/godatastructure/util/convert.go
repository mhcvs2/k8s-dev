package util

import (
	"reflect"
	"strconv"
	"net/url"
)

func ToString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Int32:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Bool:
		return strconv.FormatBool(value.Bool())
	case reflect.ValueOf(&url.URL{}).Kind():
		return value.Interface().(*url.URL).String()
	}
	return "not parse"
}
