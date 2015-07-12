package json

// Package mappath provides functions that simplify dealing with
// deeply nested map[string]interface{}/[]interface{} hierarchies, such
// as those produced by encoding/json.  A 'map path' is a string that
// indicates a location in such a nested map, consisting of a sequence
// of map keys or slice indices delimited by a separator character,
// which is taken to be the first character in the path.  Note that as
// an empty string is a legal map key, a map path ending in the
// separator character refers to the empty string key of its last map,
// and a sequence of separator characters is semantically distinct from
// a single one.

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Possible Error messages.
const (
	mappathErrIndexOutOfRange = "slice index out of range"
	mappathErrIndexNotInt     = "slice received non-int index"
	mappathErrNotIndexable    = "don't know how to index value"
)

// An Error is returned by Set or Encode if one tries to set an index
// on an object that does not support it.
type mappathError struct {
	Mappath string      // path to the object where the error happened
	Message string      // error message
	Value   interface{} // key or value that triggered the error
}

func (e *mappathError) Error() string {
	return fmt.Sprintf("mappath: %s: %s: %#v", e.Mappath, e.Message, e.Value)
}

// Join reconstructs a map path from an array of its keys and a
// separator string.  No sanity checks are performed on the input, so
// unlike Split, a separator can be longer than one character.
func mappathJoin(keys []string, sep string) string {
	return sep + strings.Join(keys, sep)
}

// Split splits a map path into its constituent keys and returns the
// separator string.  Returns an empty array and string if given an
// empty string.
func mappathSplit(mappath string) ([]string, string) {
	if mappath == "" {
		return []string{}, ""
	}
	sep := mappath[:1]
	return strings.Split(mappath[1:], sep), sep
}

// Get returns the value at the given path in the object and true if it
// exists, or nil and false if it does not.  Returns the object itself
// for an empty path.
func mappathGet(obj interface{}, mappath string) (interface{}, bool) {
	val := obj
	keys, _ := mappathSplit(mappath)
	for _, key := range keys {
		var ok bool
		switch obj := val.(type) {
		case map[string]interface{}:
			val, ok = obj[key]
			if !ok {
				return nil, false
			}
		case []interface{}:
			i, err := strconv.Atoi(key)
			if err != nil || i < 0 || i >= len(obj) {
				return nil, false
			}
			val = obj[i]
		default:
			return nil, false
		}
	}
	return val, true
}

// Set sets the value at the given path in the object, creating any
// missing intermediate objects in the process.  Only maps will be
// created; decimal keys in the map path don't create slices.
// Intermediate objects are never overwritten: an error is instead
// returned if such an object is not a string map or a slice, or if a
// slice is indexed with a non-integer or out-of-range index.  However,
// the object referenced by the path is always overwritten by the given
// value, even if it contains a deep hierarchy.
//
// An empty map path is a no-op.
func mappathSet(obj interface{}, mappath string, v interface{}) error {
	keys, sep := mappathSplit(mappath)
	for i, key := range keys {
		imappath := mappathJoin(keys[:i], sep)
		switch obj1 := obj.(type) {
		case map[string]interface{}:
			if i == len(keys)-1 {
				obj1[key] = v
				return nil
			}
			if _, ok := obj1[key]; !ok {
				obj1[key] = make(map[string]interface{})
			}
			obj = obj1[key]
		case []interface{}:
			j, err := strconv.Atoi(key)
			if err != nil {
				return &mappathError{
					imappath,
					mappathErrIndexNotInt,
					key,
				}
			}
			if j < 0 || j >= len(obj1) {
				return &mappathError{
					imappath,
					mappathErrIndexOutOfRange,
					j,
				}
			}
			if i == len(keys)-1 {
				obj1[j] = v
				return nil
			}
			obj = obj1[j]
		default:
			return &mappathError{
				imappath,
				mappathErrNotIndexable,
				obj1,
			}
		}
	}
	return nil
}

// Decode decodes an object hierarchy into a struct, using the mappath
// tag key of each field as the map path to read its value from.
// Fields with an empty mappath and missing values in the object
// hierarchy will be ignored.
func mappathDecode(obj interface{}, dst interface{}) {
	v := reflect.ValueOf(dst).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		mappath, _ := parseTag(t.Field(i).Tag.Get("mappath"))
		if mappath == "" {
			continue
		}
		val, ok := mappathGet(obj, mappath)
		if !ok {
			continue
		}
		v.Field(i).Set(reflect.ValueOf(val))
	}
}

// Encode encodes a struct into an object hierarchy, using the mappath
// tag key of each field as the map path to write its value to.
// Fields with an empty mappath will be ignored.  If an error
// occurs, the in-progress hierarchy and the error are returned.
func mappathEncode(src interface{}) (map[string]interface{}, error) {
	msg := make(map[string]interface{})
	v := reflect.ValueOf(src).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		mappath, opts := parseTag(t.Field(i).Tag.Get("mappath"))
		if mappath == "" {
			continue
		}
		x := v.Field(i).Interface()
		if opts["omitempty"] && isZero(x) {
			continue
		}
		if err := mappathSet(msg, mappath, x); err != nil {
			return msg, err
		}
	}
	return msg, nil
}

// parseTag breaks a tag string into a map path and an options map.
// it is used by Decode and Encode.
func parseTag(tag string) (string, map[string]bool) {
	parts := strings.Split(tag, ",")
	mappath, sopts := parts[0], parts[1:]
	mopts := make(map[string]bool)
	for _, opt := range sopts {
		mopts[opt] = true
	}
	return mappath, mopts
}

// isZero returns true if its argument is the zero value of its type.
func isZero(x interface{}) bool {
	return reflect.Zero(reflect.TypeOf(x)).Interface() == x
}
