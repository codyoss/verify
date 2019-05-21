// Package verify uses struct field tags to verify data. There are five tags currently supported:
//
// minSize -- specifies the minimum allowable length of a field. This can only be used on the following types: string,
// slice, array, or map.
//
// maxSize -- specifies the maximum allowable length of a field. This can only be used on the following types: string,
// slice, array, or map.
//
// min -- specifies the minimum allowable value of a field. This should only be used on types that can be parsed into
// an int64 or float64.
//
// max -- specifies the maximum allowable value of a field. This should only be used on types that can be parsed into
// an int64 or float64.
//
// required -- specifies the field may not be set to the zero value for the given type. This may be used on any types
// except arrays and structs.
//
// Here is an example of the usage of each tag:
//
//  type Foo struct {
//		A []string 	`verify:"minSize=5"`
//		B string 	`verify:"maxSize=10"`
//		C int8 		`verify:"min=3"`
//		D float32 	`verify:"max=1.2"`
//		E int64 	`verify:"min=3,max=7"`
//		F *bool 	`verify:"required"`
//  }
//
// There are currently a few limitation with this project. The first is verify only supports working with flat
// structures at the moment; it will not work with inner/embedded structs. Also, because the package makes use of
// reflection the tags may only be used on exported fields.
package verify

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	verifyTagKey = "verify"
	tagMinSize   = "minSize"
	tagMaxSize   = "maxSize"
	tagMin       = "min"
	tagMax       = "max"
	tagRequired  = "required"

	parseBase = 10
	parseBit  = 64
)

var (
	errInvalidKind = errors.New("v provided must be a struct, interface, or pointer to a struct")

	errMissingValueMinSize = errors.New("minSize must specify a size")
	errMissingValueMaxSize = errors.New("maxSize must specify a size")
	errMissingValueMin     = errors.New("min must specify a size")
	errMissingValueMax     = errors.New("max must specify a size")

	errValueTypeMinSize = errors.New("minSize can only be used with types: string, slice, array, or map")
	errValueTypeMaxSize = errors.New("maxSize can only be used with types: string, slice, array, or map")
	errValueTypeMin     = errors.New("min can only be used with types: int, int8, int16, int32, int64, float32, or float64")
	errValueTypeMax     = errors.New("max can only be used with types: int, int8, int16, int32, int64, float32, or float64")

	errConvertToNumberMinSize = errors.New("minSize value must be an int")
	errConvertToNumberMaxSize = errors.New("maxSize value must be an int")
	errConvertToNumberMin     = errors.New("min value must be an int64 or float64")
	errConvertToNumberMax     = errors.New("max value must be an int or float64")
)

// It takes a struct and uses reflection to verify it based on its struct field tags. An error is returned should any of
// the fields fail their validation. The returned error will describe each field that failed validation. Only interfaces
// a struct, or a pointer to struct should be passed to this function.
func It(v interface{}) error {
	rv := reflect.ValueOf(v)

	for rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return errInvalidKind
	}

	rt := rv.Type()
	// TODO: Append errors
	for i := 0; i < rt.NumField(); i++ {
		if tags, ok := rt.Field(i).Tag.Lookup(verifyTagKey); ok {
			err := verifyField(rv.Field(i), rt.Field(i).Name, tags)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func verifyField(f reflect.Value, name string, tag string) error {
	var tagErrs []string
	var tagPrefix string
	st := strings.Split(tag, ",")

	// verify each valid sub-tag found
	for _, v := range st {
		tagPrefix = v
		i := strings.IndexByte(v, '=')
		if i != -1 {
			tagPrefix = v[:i]
		}
		switch tagPrefix {
		case tagMinSize:
			if i == -1 {
				return errMissingValueMinSize
			}
			min, err := strconv.Atoi(v[i+1:])
			if err != nil {
				return errConvertToNumberMinSize
			}

			switch f.Kind() {
			case reflect.String, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
				if f.Len() < min {
					tagErrs = append(tagErrs, fmt.Sprintf("%s has a length less than %d", name, min))
				}
			default:
				return errValueTypeMinSize
			}
		case tagMaxSize:
			if i == -1 {
				return errMissingValueMaxSize
			}
			max, err := strconv.Atoi(v[i+1:])
			if err != nil {
				return errConvertToNumberMaxSize
			}

			switch f.Kind() {
			case reflect.String, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
				if f.Len() > max {
					tagErrs = append(tagErrs, fmt.Sprintf("%s has a length greater than %d", name, max))
				}
			default:
				return errValueTypeMaxSize
			}
		case tagMin:
			var minI int64
			var minF float64
			var isMinFloat bool
			if i == -1 {
				return errMissingValueMin
			}
			minI, err := strconv.ParseInt(v[i+1:], parseBase, parseBit)
			if err != nil {
				minF, err = strconv.ParseFloat(v[i+1:], parseBit)
				if err != nil {
					return errConvertToNumberMin
				}
				isMinFloat = true
			}
			switch f.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if isMinFloat {
					return fmt.Errorf("%s type is int while min is float", name)
				}
				if f.Int() < minI {
					tagErrs = append(tagErrs, fmt.Sprintf("%s has value less than min %d", name, minI))
				}
			case reflect.Float32, reflect.Float64:
				if !isMinFloat {
					return fmt.Errorf("%s type is float while min is int", name)
				}
				if f.Float() < minF {
					tagErrs = append(tagErrs, fmt.Sprintf("%s has value less than min %f", name, minF))
				}
			default:
				return errValueTypeMin
			}
		case tagMax:
			var maxI int64
			var maxF float64
			var isMaxFloat bool
			if i == -1 {
				return errMissingValueMax
			}
			maxI, err := strconv.ParseInt(v[i+1:], parseBase, parseBit)
			if err != nil {
				maxF, err = strconv.ParseFloat(v[i+1:], parseBit)
				if err != nil {
					return errConvertToNumberMax
				}
				isMaxFloat = true
			}
			switch f.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if isMaxFloat {
					return fmt.Errorf("%s type is int while max is float", name)
				}
				if f.Int() > maxI {
					tagErrs = append(tagErrs, fmt.Sprintf("%s has value greater than max %d", name, maxI))
				}
			case reflect.Float32, reflect.Float64:
				if !isMaxFloat {
					return fmt.Errorf("%s type is float while max is int", name)
				}
				if f.Float() > maxF {
					tagErrs = append(tagErrs, fmt.Sprintf("%s has value greater than max %f", name, maxF))
				}
			default:
				return errValueTypeMax
			}
		case tagRequired:
			switch f.Kind() {
			case reflect.Func, reflect.Map, reflect.Slice:
				if f.IsNil() {
					tagErrs = append(tagErrs, fmt.Sprintf("%s is required but is set to zero value", name))
				}
			case reflect.Array, reflect.Struct:
			default:
				if f.Interface() == reflect.Zero(f.Type()).Interface() {
					tagErrs = append(tagErrs, fmt.Sprintf("%s is required but is set to zero value", name))
				}
			}
		}
	}

	// collect all errors to return to user
	if tagErrs != nil {
		var sb strings.Builder
		for i, v := range tagErrs {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(v)
		}
		return fmt.Errorf("verify found the following errors: [%s]", sb.String())
	}

	return nil
}
