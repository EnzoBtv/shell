package utils

import (
	"reflect"
)

type mapf func(interface{}, int) interface{}

// MapArray receives an array of any type and a function and returns another array with the modifications returned in the function
func MapArray(arr interface{}, function mapf) []interface{} {
	in := reflect.ValueOf(arr)
	out := make([]interface{}, in.Len())

	for i := 0; i < in.Len(); i++ {
		functionReturn := function(in.Index(i).Interface(), i)

		out[i] = functionReturn
	}

	return out
}
