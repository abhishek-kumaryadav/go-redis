package service

import (
	"fmt"
	"reflect"
)

func CastToType[T any](repository map[string]interface{}, key string, create bool) (*T, error) {
	data, ok := repository[key]
	if !ok {
		if !create {
			return nil, fmt.Errorf("Value not present for key %s", key)
		} else {
			temp := allocateMemory[T]()
			data = &temp
			repository[key] = data
		}
	}
	if data == nil {
		return nil, fmt.Errorf("Initialization failure")
	}

	castedValue, ok := data.(*T)
	if !ok {
		return nil, fmt.Errorf("failed to cast to required type")
	}
	return castedValue, nil
}

func allocateMemory[T any]() T {
	var temp T
	switch any(temp).(type) {
	case *struct{}:
		// For structs, create a new pointer
		temp = any(new(T)).(T)
	case []byte:
		// For byte slices, create an empty slice
		temp = any(make([]byte, 0)).(T)
	case map[string]string:
		// For maps, create an empty map
		temp = any(make(map[string]string)).(T)
	case map[string]int:
		// For maps, create an empty map
		temp = any(make(map[string]int)).(T)
	case []string:
		// For string slices, create an empty slice
		temp = any(make([]string, 0)).(T)
	case *int, *string, *bool:
		// For pointer types, create a new pointer
		temp = any(new(T)).(T)
	default:
		// For other types, try to use reflection if needed
		v := reflect.New(reflect.TypeOf(temp)).Elem()
		temp = v.Interface().(T)
	}
	return temp
}
