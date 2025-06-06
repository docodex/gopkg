package jsonx

import (
	"encoding/json"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// MarshalToString returns the JSON-encoded string of v.
func MarshalToString(v any) (string, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// UnmarshalFromString parses the JSON-encoded string and stores the result
// in the value pointed to by v.
func UnmarshalFromString(data string, v any) error {
	return json.Unmarshal([]byte(data), v)
}

// MGet searches json text for the multiple paths.
// The return value is a Result array where the number of items
// will be equal to the number of input paths.
func MGet(text string, paths ...string) []gjson.Result {
	return gjson.GetMany(text, paths...)
}

// MGetBytes searches json bytes for the multiple paths.
// The return value is a Result array where the number of items
// will be equal to the number of input paths.
// If working with bytes, this method preferred over BatchGet(string(bytes), paths)
func MGetBytes(bytes []byte, paths ...string) []gjson.Result {
	return gjson.GetManyBytes(bytes, paths...)
}

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
