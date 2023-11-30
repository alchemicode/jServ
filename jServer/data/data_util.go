package data

import (
	"cmp"
	"reflect"
	"time"
)

func Intersect[T comparable](slice1 []T, slice2 []T) []T {
	final := make([]T, 0)
	m := make(map[T]struct{})
	for _, v := range slice1 {
		m[v] = struct{}{}
	}
	for _, v := range slice2 {
		if _, ok := m[v]; ok {
			final = append(final, v)
		}
	}
	return final
}

func TypeSwitchComparison(val1 interface{}, val2 interface{}) (int, bool) {
	type1 := reflect.TypeOf(val1)
	type2 := reflect.TypeOf(val2)
	if type1 == type2 {
		switch v1 := val1.(type) {
		case string:
			return cmp.Compare[string](v1, val2.(string)), true
		case int:
			return cmp.Compare[int](v1, val2.(int)), true
		case float64:
			return cmp.Compare[float64](v1, val2.(float64)), true
		case time.Time:
			v2 := val2.(time.Time)
			if v1.After(v2) {
				return 1, true
			} else if v1.Before(v2) {
				return -1, true
			} else if v1.Equal(v2) {
				return 0, true
			} else {
				return 0, false
			}
		case []interface{}:
			v2 := val2.([]interface{})
			if reflect.DeepEqual(v1, v2) {
				return 0, true
			} else {
				return 0, false
			}
		case map[string]interface{}:
			v2 := val2.(map[string]interface{})
			if reflect.DeepEqual(v1, v2) {
				return 0, true
			} else {
				return 0, false
			}
		default:
			return 0, false
		}
	} else {
		return 0, false
	}
}

func Remove[T comparable](elem T, slice []T) []T {
	var new_slice []T = slice
	if ok, index := Contains(new_slice, elem); ok {
		new_slice = append(new_slice[:index], slice[index+1:]...)
	}
	return new_slice
}

func Contains[T comparable](slice []T, elem T) (bool, int) {
	for i, a := range slice {
		if a == elem {
			return true, i
		}
	}
	return false, -1
}
