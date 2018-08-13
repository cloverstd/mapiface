package mapiface

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Option .
type Option struct {
	MaxDep int
	Tag    string
}

var mapiface = map[string]interface{}{}
var ifaceType = reflect.TypeOf(mapiface).Elem()
var zero = reflect.Zero(ifaceType)
var sliceIfaceType = reflect.TypeOf([]interface{}{})
var mapIfaceType = reflect.TypeOf(mapiface)

// ErrMaxDepth .
var ErrMaxDepth = errors.New("overflow maxdep")

// Convert a value to mapiface
func Convert(value interface{}) (interface{}, error) {
	return ConvertWithOption(value, Option{})
}

// ConvertWithOption a value to mapiface with option
func ConvertWithOption(value interface{}, op Option) (interface{}, error) {
	if op.MaxDep <= 0 {
		op.MaxDep = 1024
	}
	if op.Tag == "" {
		op.Tag = "json"
	}
	v, err := convert(reflect.ValueOf(value), op, 0)
	if err != nil {
		return nil, err
	}
	return v.Interface(), nil
}

func convert(v reflect.Value, op Option, dep int) (reflect.Value, error) {
	if dep > op.MaxDep {
		return v, ErrMaxDepth
	}
	if !v.IsValid() {
		return v, fmt.Errorf("invalid value")
	}
	kind := v.Kind()
	typ := v.Type()
	switch kind {
	case reflect.Map:
		mapiface := reflect.MakeMapWithSize(reflect.MapOf(typ.Key(), ifaceType), v.Len())
		if !v.IsNil() {
			for _, key := range v.MapKeys() {
				v, err := convert(v.MapIndex(key), op, dep+1)
				if err != nil {
					return v, err
				}
				mapiface.SetMapIndex(key, v)
			}
		}
		return mapiface, nil
	case reflect.Ptr:
		if v.IsNil() {
			return reflect.Zero(v.Type()), nil
		}
		value := reflect.New(typ.Elem())
		v, err := convert(v.Elem(), op, dep+1)
		if err != nil {
			return v, err
		}
		value.Elem().Set(v)
		return value, nil
	case reflect.Array, reflect.Slice:
		sliceiface := reflect.MakeSlice(sliceIfaceType, v.Len(), v.Cap())
		if !v.IsNil() {
			for i := 0; i < v.Len(); i++ {
				v, err := convert(v.Index(i), op, dep+1)
				if err != nil {
					return v, err
				}
				sliceiface.Index(i).Set(v)
			}
		}
		return sliceiface, nil
	case reflect.Struct:
		mapiface := reflect.MakeMap(mapIfaceType)
		for i := 0; i < v.NumField(); i++ {
			field := typ.Field(i)
			if field.PkgPath != "" {
				continue
			}
			fieldName := field.Name
			omitempty := false
			if tag, ok := field.Tag.Lookup(op.Tag); ok {
				tags := strings.Split(tag, ",")
				fieldName = tags[0]
				if len(tags) >= 2 && tags[1] == "omitempty" {
					omitempty = true
				}
			}
			value := v.Field(i)
			if isNil(value) && omitempty {
				continue
			}
			value, err := convert(value, op, dep+1)
			if err != nil {
				return value, err
			}
			mapiface.SetMapIndex(reflect.ValueOf(fieldName), value)
		}
		return mapiface, nil
	default:
		return v, nil
	}
}

func isNil(v reflect.Value) bool {
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return false
}
