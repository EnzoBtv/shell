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

//@TODO CONVERT THIS FUNCTIONS TO DYNAMIC TYPING

// SliceContains checks if an item is in the slice
func SliceContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

// JoinSlices joins two slices without repeating elements
func JoinSlices(new []string, existing *[]string) {
	for _, item := range new {
		if !SliceContains(*existing, item) {
			*existing = append(*existing, item)
		}
	}
}
