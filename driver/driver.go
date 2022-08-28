// Copyright (c) 554949297@qq.com . 2022-2022 . All rights reserved

package driver

import (
	"encoding"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"strconv"
)

var binaryUnmarshalerType = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()

func GetT(data string, t encoding.BinaryUnmarshaler) (any, bool) {
	if t == nil {
		return nil, false
	}
	err := t.UnmarshalBinary([]byte(data))
	if err == nil {
		return t, true
	}
	return nil, false
}

func ParseString(str string, result any) bool {
	if reflect.TypeOf(result).Implements(binaryUnmarshalerType) {
		if _result, o := result.(encoding.BinaryUnmarshaler); o {
			if _result.UnmarshalBinary([]byte(str)) == nil {
				return true
			}
		}
		return false
	}
	if reflect.TypeOf(result).Kind() == reflect.Ptr {
		if reflect.TypeOf(result).Elem().Kind() == reflect.String {
			reflect.ValueOf(result).Elem().Set(reflect.ValueOf(str))
			return true
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Slice ||
			reflect.TypeOf(result).Elem().Kind() == reflect.Array ||
			reflect.TypeOf(result).Elem().Kind() == reflect.Map {
			err := jsoniter.UnmarshalFromString(str, result)
			if err == nil {
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Int {
			if i, e := strconv.Atoi(str); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(i))
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Int8 {
			if i, e := strconv.ParseInt(str, 10, 8); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(int8(i)))
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Int16 {
			if i, e := strconv.ParseInt(str, 10, 16); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(int16(i)))
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Int32 {
			if i, e := strconv.ParseInt(str, 10, 32); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(int32(i)))
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Int64 {
			if i, e := strconv.ParseInt(str, 10, 64); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(i))
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Uint {
			if i, e := strconv.ParseUint(str, 10, 0); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(uint(i)))
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Uint8 {
			if i, e := strconv.ParseUint(str, 10, 8); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(uint8(i)))
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Uint16 {
			if i, e := strconv.ParseUint(str, 10, 16); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(uint16(i)))
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Uint32 {
			if i, e := strconv.ParseUint(str, 10, 32); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(uint32(i)))
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Uint64 {
			if i, e := strconv.ParseUint(str, 10, 64); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(i))
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Float32 {
			if f, e := strconv.ParseFloat(str, 32); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(float32(f)))
				return true
			}
		} else if reflect.TypeOf(result).Elem().Kind() == reflect.Float64 {
			if f, e := strconv.ParseFloat(str, 64); e == nil {
				reflect.ValueOf(result).Elem().Set(reflect.ValueOf(f))
				return true
			}
		}
	}
	return false
}
