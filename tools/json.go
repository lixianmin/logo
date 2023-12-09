package tools

import (
	"fmt"
	"github.com/lixianmin/got/convert"
	"reflect"
	"strconv"
	"sync"
	"time"
)

/********************************************************************
created:    2021-01-07
author:     lixianmin

this file is derived from go-redis/v8/internal/util.go
*********************************************************************/

var bufferPool = &sync.Pool{
	New: func() any {
		return make([]byte, 0, 256)
	},
}

func FormatJson(args ...any) string {
	var results = bufferPool.Get().([]byte)
	results = append(results, '{')
	{
		var count = len(args)
		var halfCount = (count + 1) >> 1
		for i := 0; i < halfCount; i++ {
			var index = i << 1
			var key, _ = args[index].(string)

			// 如果只有奇数个参数，则输出默认值null
			index++
			var value interface{} = nil
			if index < count {
				value = args[index]
			}

			results = strconv.AppendQuote(results, key)
			results = append(results, ':')
			results = AppendJson(results, value)

			if i+1 < halfCount {
				results = append(results, ',')
			}
		}
	}
	results = append(results, '}')

	// 这里results马上就要还回去了，不要使用unsafe的[]byte转string了
	var text = string(results)
	bufferPool.Put(results[:0])
	return text
}

func AppendJson(b []byte, v any) []byte {
	// v.(type)有值, 不代表v!=nil
	const null = "nil"
	switch v := v.(type) {
	case nil:
		return append(b, null...)
	case string:
		return strconv.AppendQuote(b, v)
	case []byte:
		return strconv.AppendQuote(b, convert.String(v))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return convert.AppendInt(b, v, 10)
	case float32:
		return strconv.AppendFloat(b, float64(v), 'f', -1, 64)
	case float64:
		return strconv.AppendFloat(b, v, 'f', -1, 64)
	case bool:
		return strconv.AppendBool(b, v)
	case error:
		var v2 = null
		if !isNil(v) {
			v2 = v.Error()
		}
		return strconv.AppendQuote(b, v2)
	case fmt.Stringer: // 实现String()方法
		var v2 = null
		if !isNil(v) {
			v2 = v.String()
		}
		return strconv.AppendQuote(b, v2)
	case time.Time:
		b = append(b, '"')
		b = v.AppendFormat(b, time.RFC3339Nano)
		b = append(b, '"')
		return b
	default:
		return append(b, convert.ToJson(v)...)
	}
}

func isNil(input any) bool {
	var value = reflect.ValueOf(input)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return value.IsNil()
	}

	return false
}

//func appendUTF8String(b []byte, s string) []byte {
//	for _, r := range s {
//		b = appendRune(b, r)
//	}
//	return b
//}
//
//func appendRune(b []byte, r rune) []byte {
//	if r < utf8.RuneSelf {
//		switch c := byte(r); c {
//		case '\n':
//			return append(b, "\\n"...)
//		case '\r':
//			return append(b, "\\r"...)
//		default:
//			return append(b, c)
//		}
//	}
//
//	l := len(b)
//	b = append(b, make([]byte, utf8.UTFMax)...)
//	n := utf8.EncodeRune(b[l:l+utf8.UTFMax], r)
//	b = b[:l+n]
//
//	return b
//}
