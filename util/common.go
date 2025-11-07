package util

import (
	"encoding/binary"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"ppt/config"
	"ppt/log"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
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

// StructToJsonMap 结构体转json字段名Map
func StructToJsonMap(obj interface{}) map[string]string {
	result := make(map[string]string)
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	re := regexp.MustCompile(`[,;]+`)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		jsonTag := val.Type().Field(i).Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		jsonTags := re.Split(jsonTag, -1)
		if len(jsonTags) == 0 {
			continue
		}
		jsonName := jsonTags[0]

		var fieldValue string
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
		result[jsonName] = fieldValue
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

// FormatESLogIndex 获取es索引
func FormatESLogIndex(specTimeSec int64) string {
	specTime := time.Unix(specTimeSec, 0)
	y, m, d := specTime.Date()
	return fmt.Sprintf("%s-%s-%04d-%02d-%02d", config.AppName, config.Env, y, m, d)
}

// Int32ToBytes Int32转字节-小端序
func Int32ToBytes(v int32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, uint32(v))
	return res
}

// ByteToInt32 Byte转Int32-小端序
func ByteToInt32(v []byte) int32 {
	return int32(binary.LittleEndian.Uint32(v))
}

func UnixSecToDateStr(sec int64) string {
	t := time.Unix(sec, 0)
	return t.Format(time.DateOnly)
}

func UnixMilliToDateStr(milli int64) string {
	sec := milli / 1000
	nsec := (milli % 1000) * 1000000
	t := time.Unix(sec, nsec)
	return t.Format(time.DateOnly)
}

func TimeToDateStr(t time.Time) string {
	return t.Format(time.DateOnly)
}

// GetTz 获取时区
func GetTz(timeZone string) *time.Location {
	loc, err := time.LoadLocation("")
	if err != nil {
		log.Error("GetTz LoadLocation error", zap.Error(err))
		return nil
	}
	if timeZone != "" {
		loc, err = time.LoadLocation(timeZone)
		if err != nil {
			log.Error("GetTz spec LoadLocation error", zap.String("time_zone", timeZone), zap.Error(err))
			return nil
		}
	}
	return loc
}

// CalcDiffSetGeneric 求差集
func CalcDiffSetGeneric[T comparable](s1, s2 []T) []T {
	set2 := make(map[T]struct{}, len(s2))
	for _, v := range s2 {
		set2[v] = struct{}{}
	}
	diff := make([]T, 0, len(s1))
	for _, v := range s1 {
		if _, ok := set2[v]; !ok {
			diff = append(diff, v)
		}
	}
	return diff
}
