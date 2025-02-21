package util

import (
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// StructToMap 结构体转map[string]string
func StructToMap(obj interface{}) map[string]string {
	result := make(map[string]string)
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		var fieldValue string
		if field.Kind() == reflect.Struct {
			nestedMap := StructToMap(field.Interface())
			for k, v := range nestedMap {
				result[fieldName+"."+k] = v
			}
		}

		switch field.Kind() {
		case reflect.String:
			fieldValue = field.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldValue = strconv.FormatInt(field.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldValue = strconv.FormatUint(field.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			fieldValue = strconv.FormatFloat(field.Float(), 'g', -1, 64)
		case reflect.Bool:
			fieldValue = strconv.FormatBool(field.Bool())
		default:
			fieldValue = ""
		}
		result[fieldName] = fieldValue
	}
	return result
}

// SortAndConcat 排序拼接
func SortAndConcat(params map[string]string) string {
	// 过滤值为空的参数
	filteredParams := make(map[string]string)
	for k, v := range params {
		if v != "" {
			filteredParams[k] = v
		}
	}
	// 获取键并按ASCII排序
	keys := make([]string, 0, len(filteredParams))
	for k := range filteredParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 拼接键值对
	var strBuilder strings.Builder
	for i, k := range keys {
		if i > 0 {
			strBuilder.WriteString("&")
		}
		strBuilder.WriteString(k)
		strBuilder.WriteString("=")
		strBuilder.WriteString(url.QueryEscape(filteredParams[k]))
	}
	return strBuilder.String()
}
