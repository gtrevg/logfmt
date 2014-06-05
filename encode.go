// Package implements the encoding of logfmt key-value pairs.
//
// Example data:
//
//	{ "foo", "bar", "a", 14, "baz", "hello kitty", "cool%story": "bro", "f": true, "%^asdf": true }
//
// Example result logfmt message:
//
//      foo=bar a=14 baz="hello kitty" cool%story=bro f=1 %^asdf=1
//
// This code is completely taken from Inconshreveable's log15 project:
//
//   https://github.com/inconshreveable/log15/blob/master/format.go
//
// The "Marshal" function is almost as the original logfmt function.  This file
// should benefit everyone who's trying to generate this format in go.
//

package logfmt

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	timeLayout  = "2006-01-02T15:04:05-0700"
	floatFormat = 'f'
	errorKey    = "logfmtError"
)

func Marshal(items ...interface{}) ([]byte, error) {
	pieces := make([]string, 0)

	for i := 0; i < len(items); i += 2 {
		k, ok := items[i].(string)
		var s string
		if !ok {
			return nil, nil //s = fmt.Sprintf(`%s="%+v is not a string key"`, errorKey, items[i])
		} else {
			// XXX: we should probably check that all of your key bytes aren't invalid`
			s = fmt.Sprintf(`%s=%s`, k, encodeValue(items[i+1]))
		}

		pieces = append(pieces, s)
	}

	return []byte(strings.Join(pieces, " ")), nil // No newline line in the original
}

func formatShared(value interface{}) interface{} {
	switch v := value.(type) {
	case time.Time:
		return v.Format(timeLayout)

	case error:
		return v.Error()

	case fmt.Stringer:
		return v.String()

	default:
		return v
	}
}

// formatValue formats a value for serialization
func encodeValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	value = formatShared(value)
	switch v := value.(type) {
	case string:
		return escapeString(v)

	case bool:
		return strconv.FormatBool(v)

	case int:
		return strconv.FormatInt(int64(v), 10)

	case int8:
		return strconv.FormatInt(int64(v), 10)

	case int16:
		return strconv.FormatInt(int64(v), 10)

	case int32:
		return strconv.FormatInt(int64(v), 10)

	case int64:
		return strconv.FormatInt(v, 10)

	case float32:
		return strconv.FormatFloat(float64(v), floatFormat, -1, 64) // Don't round. Was 3 insteadof -1

	case float64:
		return strconv.FormatFloat(v, floatFormat, -1, 64) // Don't round. was 3 instead of -1

	case uint:
		return strconv.FormatUint(uint64(v), 10)

	case uint8:
		return strconv.FormatUint(uint64(v), 10)

	case uint16:
		return strconv.FormatUint(uint64(v), 10)

	case uint32:
		return strconv.FormatUint(uint64(v), 10)

	case uint64:
		return strconv.FormatUint(v, 10)

	default:
		return escapeString(fmt.Sprintf("%+v", value))
	}
}

func escapeString(s string) string {
	needQuotes := false
	e := new(bytes.Buffer)
	e.WriteByte('"')
	for _, r := range s {
		if r <= ' ' || r == '=' || r == '"' {
			needQuotes = true
		}

		switch r {
		case '"':
			e.WriteByte('\\')
			e.WriteByte(byte(r))
		default:
			e.WriteRune(r)
		}
	}
	e.WriteByte('"')
	start, stop := 0, e.Len()
	if !needQuotes {
		start, stop = 1, stop-1
	}
	return string(e.Bytes()[start:stop])
}
