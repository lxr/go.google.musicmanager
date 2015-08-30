/*

Package convert uses reflection to deep convert between arbitrary value
hierarchies.  The rules for conversion are as follows:

 - Basic types are converted according to the usual Go conversion
   rules.  They panic if they cannot be converted.
 - Pointers and interfaces are indirected.  If a pointer on the
   destination side is nil and the source side has a non-pointer or a
   non-nil pointer, a new zero value is allocated for it before
   traversal.  Nil pointers on the destination side are left as-is
   otherwise.  A nil pointer or interface on the source side causes the
   destination to be set to its zero value.
 - Structs can be converted from other structs or string-keyed maps.
   Missing field names on the destination side are silently ignored.
   (Trying to convert an array, slice, or non-string-keyed map to a
   struct results in an  empty struct.)
 - Maps can be converted from other maps, arrays, slices or structs,
   provided that the key type is convertible to the map's key type; if
   it is not, the conversion panics.  (Structs have "keys" of type
   string and slices and arrays of type int.)
 - Arrays and slices can be converted from each other and signed integer
   -kind-keyed-maps.  (Structs and incorrectly keyed maps cause a
   panic.)  Indices out of range cause a panic, unless the index is
   simply too large for a slice on the destination side, in which case
   it is grown to accommodate.
 - When trying to convert an array, slice, struct, or map to a nil
   interface value, convert tries to store a new map[T]interface{} in
   it, where T is the appropriate key type.  It panics if this fails.

Values in the destination hierarchy that the source hierarchy does not
map to are not touched by the conversion.

As an additional convenience, convert supports the "convert" tag key on
struct fields, which declares a location in the destination/source value
hierarchy where/from which to convert the field.  The first character in
the tag is the key separator.  Keys that are parsable as integers are
treated as such, otherwise all keys are taken as string keys.  Note
that, as the empty string is a legal map key, all of the following are
distinct paths:

   /testing
   /testing/
   //testing

A convert tag with the value "-" means to ignore the field.  If both the
source and destination values are structs, the convert tags on the
destination side take precedence.

*/
package convert

import (
	"reflect"
	"strconv"
	"strings"
)

// Convert converts the source value hierarchy to the destination one.
func Convert(dst, src interface{}) {
	convert(reflect.ValueOf(dst), reflect.ValueOf(src))
}

func convert(dst, src value) {
	switch dst.Kind() {
	case reflect.Invalid:
		return
	case reflect.Ptr:
		if dst.IsNil() && (src.Kind() != reflect.Ptr || !src.IsNil()) {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		convert(dst.Elem(), src)
		return
	case reflect.Struct:
		t := dst.Type()
		for i := 0; i < t.NumField(); i++ {
			if src := index(src, false, getPath(t.Field(i))...); src.IsValid() {
				convert(dst.Field(i), src)
			}
		}
		return
	}

	switch src.Kind() {
	case reflect.Invalid:
		dst.Set(reflect.Zero(dst.Type()))
	case reflect.Ptr, reflect.Interface:
		convert(dst, src.Elem())
	case reflect.Struct:
		t := src.Type()
		for i := 0; i < t.NumField(); i++ {
			convert(index(dst, true, getPath(t.Field(i))...), src.Field(i))
		}
	case reflect.Slice, reflect.Array:
		// We iterate over the values in reverse to ensure that
		// index() creates a correctly-sized array right off the
		// bat.
		for i := src.Len() - 1; i >= 0; i-- {
			convert(index(dst, true, reflect.ValueOf(i)), src.Index(i))
		}
	case reflect.Map:
		for _, key := range src.MapKeys() {
			convert(index(dst, true, key), src.MapIndex(key))
		}
	default:
		dst.Set(src.Convert(dst.Type()))
	}
}

// reflect.Values acquired using reflect.Value.MapIndex aren't settable,
// unlike those acquired with reflect.Value.Index or
// reflect.Value.Field.  In order to make reflect.Setting of values
// uniform, we define a special type for map values and an interface
// so that it can be used transparently in place of reflect.Value.

type mapValue struct {
	reflect.Value
	m value
	k reflect.Value
}

func (v *mapValue) Set(x reflect.Value) {
	v.m.SetMapIndex(v.k, x)
	v.Value = v.m.MapIndex(v.k)
}

type value interface {
	Convert(reflect.Type) reflect.Value
	Elem() reflect.Value
	Field(i int) reflect.Value
	FieldByName(s string) reflect.Value
	Index(i int) reflect.Value
	Interface() interface{}
	IsNil() bool
	IsValid() bool
	Kind() reflect.Kind
	Len() int
	MapIndex(key reflect.Value) reflect.Value
	MapKeys() []reflect.Value
	Set(x reflect.Value)
	SetMapIndex(key, elem reflect.Value)
	Type() reflect.Type
}

// Shortcut to the reflect.Type of the empty interface.
// Is there a smarter way to express this?
var emptyInterface = reflect.TypeOf([]interface{}{}).Elem()

// index walks the value hierarchy rooted v at according to the given
// path, creating any missing intermediate values if create is true.
// An empty path returns v.  An invalid path or a path that can't be
// created return the zero Value.  index panics if an otherwise valid
// intermediate value encountered during the walk isn't indexable.
func index(v value, create bool, path ...reflect.Value) value {
	for _, key := range path {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() && create {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		}
		if v.Kind() == reflect.Interface {
			if v.IsNil() && create {
				mt := reflect.MapOf(key.Type(), emptyInterface)
				v.Set(reflect.MakeMap(mt))
			}
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Invalid:
			return v
		case reflect.Struct:
			v = v.FieldByName(key.String())
		case reflect.Slice, reflect.Array:
			i := int(key.Int())
			if want := i + 1 - v.Len(); want > 0 && v.Kind() == reflect.Slice {
				x := reflect.MakeSlice(v.Type(), want, want)
				if mv, ok := v.(*mapValue); ok {
					v = mv.Value
				}
				v.Set(reflect.AppendSlice(v.(reflect.Value), x))
			}
			v = v.Index(i)
		case reflect.Map:
			if v.IsNil() && create {
				v.Set(reflect.MakeMap(v.Type()))
			}
			key = key.Convert(v.Type().Key())
			x := v.MapIndex(key)
			if !v.IsNil() && !x.IsValid() {
				x = reflect.Zero(v.Type().Elem())
			}
			v = &mapValue{
				Value: x,
				m:     v,
				k:     key,
			}
		default:
			panic("type is not indexable")
		}
	}
	return v
}

// getPath returns the convert path of the given struct field.  If
// the field has a "convert" tag key, its value is parsed as the convert
// path.  If this key is missing or empty, the convert path defaults to
// the field name.  If the key has the value "-", getPath returns a
// path consisting of a single zero Value, which index interprets as
// always missing.
func getPath(sf reflect.StructField) []reflect.Value {
	tag := sf.Tag.Get("convert")
	switch tag {
	case "-":
		return []reflect.Value{{}}
	case "":
		return []reflect.Value{reflect.ValueOf(sf.Name)}
	}
	keys := strings.Split(tag[1:], tag[:1])
	edges := make([]reflect.Value, len(keys))
	for i, key := range keys {
		n, err := strconv.Atoi(key)
		if err == nil {
			edges[i] = reflect.ValueOf(n)
		} else {
			edges[i] = reflect.ValueOf(key)
		}
	}
	return edges
}
