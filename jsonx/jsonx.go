package jsonx

import (
	"encoding/json"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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

// Valid reports whether data is a valid JSON encoding.
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

// Valid returns true if the input is valid json text.
// func Valid(text string) bool {
// 	return gjson.Valid(text)
// }

// ValidBytes returns true if the input is valid json bytes.
// If working with bytes, this method preferred over ValidBytes(string(bytes))
// func ValidBytes(bytes []byte) bool {
// 	return gjson.ValidBytes(bytes)
// }

// Get searches json text for the specified path.
// A path is in dot syntax, such as "name.last" or "age".
// When the value is found it's returned immediately.
//
// A path is a series of keys separated by a dot.
// A key may contain special wildcard characters '*' and '?'.
// To access an array value use the index as the key.
// To get the number of elements in an array or to access a child path, use
// the '#' character.
// The dot and wildcard character can be escaped with '\'.
//
//	{
//	  "name": {"first": "Tom", "last": "Anderson"},
//	  "age":37,
//	  "children": ["Sara","Alex","Jack"],
//	  "friends": [
//	    {"first": "James", "last": "Murphy"},
//	    {"first": "Roger", "last": "Craig"}
//	  ]
//	}
//	"name.last"          >> "Anderson"
//	"age"                >> 37
//	"children"           >> ["Sara","Alex","Jack"]
//	"children.#"         >> 3
//	"children.1"         >> "Alex"
//	"child*.2"           >> "Jack"
//	"c?ildren.0"         >> "Sara"
//	"friends.#.first"    >> ["James","Roger"]
//
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return unexpected results.
// If you are consuming JSON from an unpredictable source then you may want to
// use the Valid function first.
// func Get(text, path string) Result {
// 	return gjson.Get(text, path)
// }

// GetBytes searches json bytes for the specified path.
// If working with bytes, this method preferred over Get(string(bytes), path)
// func GetBytes(bytes []byte, path string) Result {
// 	return gjson.GetBytes(bytes, path)
// }

// MGet searches json text for the multiple paths.
// The return value is a Result array where the number of items
// will be equal to the number of input paths.
// func MGet(text string, paths []string) []gjson.Result {
// 	return gjson.GetMany(text, paths...)
// }

// MGetBytes searches json bytes for the multiple paths.
// The return value is a Result array where the number of items
// will be equal to the number of input paths.
// If working with bytes, this method preferred over BatchGet(string(bytes), paths)
// func MGetBytes(bytes []byte, paths []string) []gjson.Result {
// 	return gjson.GetManyBytes(bytes, paths...)
// }

// Parse parses the json text and returns a result.
//
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return unexpected results.
// If you are consuming JSON from an unpredictable source then you may want to
// use the Valid function first.
// func Parse(text string) gjson.Result {
// 	return gjson.Parse(text)
// }

// ParseBytes parses the json bytes and returns a result.
// If working with bytes, this method preferred over Parse(string(bytes))
// func ParseBytes(bytes []byte) gjson.Result {
// 	return gjson.ParseBytes(bytes)
// }

// ForEachLine iterates through lines of json as specified by the json Lines
// format (http://jsonlines.org/).
// Each line is returned as a Result.
// func ForEachLine(text string, f func(result gjson.Result) bool) {
// 	gjson.ForEachLine(text, f)
// }

// Option represents additional options for the Set and Delete functions.
type Option func(opts *sjson.Options)

// BeOptimistic sets a hint that the value likely exists which
// allows for the json to perform a fast-track search and replace.
func BeOptimistic() Option {
	return func(opts *sjson.Options) {
		opts.Optimistic = true
	}
}

// DoReplaceInPlace sets a hint to replace the input json rather than
// allocate a new json byte slice. When this field is specified
// the input json will no longer be valid, and it should not be used
// In the case when the destination slice doesn't have enough free
// bytes to replace the data in place, a new bytes slice will be
// created under the hood.
// The BeOptimistic option must be set and the input must be a byte
// slice in order to use this field.
func DoReplaceInPlace() Option {
	return func(opts *sjson.Options) {
		opts.ReplaceInPlace = true
	}
}

// Set sets a json value for the specified path with options.
// A path is in dot syntax, such as "name.last" or "age".
// This function expects that the json is well-formed, and does not validate.
// Invalid json will not panic, but it may return unexpected results.
// An error is returned if the path is not valid.
//
// A path is a series of keys separated by a dot.
//
//	{
//	  "name": {"first": "Tom", "last": "Anderson"},
//	  "age":37,
//	  "children": ["Sara","Alex","Jack"],
//	  "friends": [
//	    {"first": "James", "last": "Murphy"},
//	    {"first": "Roger", "last": "Craig"}
//	  ]
//	}
//	"name.last"          >> "Anderson"
//	"age"                >> 37
//	"children.1"         >> "Alex"
func Set(text, path string, value any, opts ...Option) (string, error) {
	config := &sjson.Options{}
	for _, opt := range opts {
		if opt != nil {
			opt(config)
		}
	}
	return sjson.SetOptions(text, path, value, config)
}

// SetBytes sets a json value for the specified path with options.
// If working with bytes, this method preferred over Set(string(bytes), path, value, opts...)
func SetBytes(bytes []byte, path string, value any, opts ...Option) ([]byte, error) {
	config := &sjson.Options{}
	for _, opt := range opts {
		if opt != nil {
			opt(config)
		}
	}
	return sjson.SetBytesOptions(bytes, path, value, config)
}

// SetRaw sets a raw json value for the specified path with options.
// This function works the same as Set except that the value is set as a
// raw block of json. This allows for setting pre-marshalled json objects.
func SetRaw(text, path, value string, opts ...Option) (string, error) {
	config := &sjson.Options{}
	for _, opt := range opts {
		if opt != nil {
			opt(config)
		}
	}
	return sjson.SetRawOptions(text, path, value, config)
}

// SetRawBytes sets a raw json value for the specified path with options.
// If working with bytes, this method preferred over SetRaw(string(bytes), path, string(value),
// opts...)
func SetRawBytes(bytes []byte, path string, value []byte, opts ...Option) ([]byte, error) {
	config := &sjson.Options{}
	for _, opt := range opts {
		if opt != nil {
			opt(config)
		}
	}
	return sjson.SetRawBytesOptions(bytes, path, value, config)
}

// Delete deletes a value from json text for the specified path.
// func Delete(text, path string) (string, error) {
// 	return sjson.Delete(text, path)
// }

// DeleteBytes deletes a value from json bytes for the specified path.
// If working with bytes, this method preferred over Delete(string(bytes), path)
// func DeleteBytes(bytes []byte, path string) ([]byte, error) {
// 	return sjson.DeleteBytes(bytes, path)
// }
