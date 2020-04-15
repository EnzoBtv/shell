package main

import (
	"reflect"
)

type mapf func(interface{}, int) interface{}

func mapArray(arr interface{}, function mapf) []interface{} {
	in := reflect.ValueOf(arr)
	out := make([]interface{}, in.Len())

	for i := 0; i < in.Len(); i++ {
		functionReturn := function(in.Index(i).Interface(), i)

		out[i] = functionReturn
	}

	return out
}
