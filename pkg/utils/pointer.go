package utils

import "reflect"

func PAdd[T int | int64 | int32 | float64](a, b *T) *T {
	if a == nil && b == nil {
		return nil
	} else if a == nil {
		return b
	} else if b == nil {
		return a
	} else {
		v := *a + *b
		return &v
	}
}

func PSub[T int | int64 | int32 | float64](a, b *T) *T {
	if a == nil && b == nil {
		return nil
	} else if a == nil {
		v := -*b
		return &v
	} else if b == nil {
		return a
	} else {
		v := *a - *b
		return &v
	}
}

func GetPointer[T any](a T) *T {
	v := a
	return &v
}

func GetPointerOrNil[T int | int64 | int32 | string](a T) *T {
	var v T
	if a == v {
		return nil
	}
	return &a
}

func getNestedZeroReflectValue(obj reflect.Type) reflect.Value {
	// recursively get the zero value of the nested objects
	// in case of pointer, get the zero value of the type it points to and populate the pointer with it
	// in case of struct, recursively call this function to get the zero value of the nested objects
	// in case of interface put nil
	// in case of other types, put the zero value
	switch obj.Kind() {
	case reflect.Ptr:
		zeroValue := getNestedZeroReflectValue(obj.Elem())
		result := reflect.New(obj.Elem())
		result.Elem().Set(zeroValue)
		return result
	case reflect.Struct:
		zeroValue := reflect.New(obj).Elem()
		for i := 0; i < obj.NumField(); i++ {
			field := obj.Field(i)
			if !field.IsExported() {
				continue
			}
			zeroValue.Field(i).Set(getNestedZeroReflectValue(obj.Field(i).Type))
		}
		return zeroValue
	case reflect.Interface:
		return reflect.Zero(obj)
	default:
		return reflect.Zero(obj)
	}
}

func GetNestedZeroValue(obj any) any {
	return getNestedZeroReflectValue(reflect.TypeOf(obj)).Interface()
}
