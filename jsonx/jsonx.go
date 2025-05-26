package jsonx

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

func MarshalToString(v any) (string, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// func Marshal(v any) ([]byte, error) {
// 	return json.Marshal(v)
// }

// func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
// 	return json.MarshalIndent(v, prefix, indent)
// }

func UnmarshalFromString(data string, v any) error {
	return json.Unmarshal([]byte(data), v)
}

// func Unmarshal(data []byte, v any) error {
// 	return json.Unmarshal(data, v)
// }

// func Valid(data []byte) bool {
// 	return json.Valid(data)
// }

// Type is Result type
type Type = gjson.Type

const (
	Null   = gjson.Null   // Null is a null json value
	False  = gjson.False  // False is a json false boolean
	Number = gjson.Number // Number is json number
	String = gjson.String // String is a json string
	True   = gjson.True   // True is a json true boolean
	JSON   = gjson.JSON   // JSON is a raw block of JSON
)

// Result represents a json value that is returned from Get().
type Result = gjson.Result
