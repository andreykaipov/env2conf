package env2star

import (
	"fmt"
	"strings"
)

// VariableArray represents an array that does not fail
// whenever a get or set is performed on an out-of-bounds index.
// Like JavaScript arrays <3.
// It's represented as a pointer to a slice of `interface{}`.
type VariableArray []interface{}

// oh no so many new slices
func (a *VariableArray) expandIfNecessary(i int) {
	if i+1 > len(*a) {
		expanded := make([]interface{}, i+1)
		copy(expanded, *a)
		*a = expanded
	}
}

func (a *VariableArray) Set(i int, val interface{}) {
	a.expandIfNecessary(i)
	(*a)[i] = val
}

func (a *VariableArray) Get(i int) interface{} {
	if i+1 > len(*a) {
		return nil
	}
	return (*a)[i]
}

// Iterator returns the underlying slice
func (a *VariableArray) Iterator() []interface{} {
	return *a
}

// IsWellDefined returns true if every element in the array
// is of the same type.
func (a *VariableArray) IsWellDefined() bool {
	t := fmt.Sprintf("%T", a.Get(0))
	for _, v := range *a {
		if fmt.Sprintf("%T", v) != t {
			return false
		}
	}
	return true
}

func (a *VariableArray) GetType() string {
	if a.IsWellDefined() {
		return fmt.Sprintf("%T", a.Get(0))
	}
	return "mixed"
}

// Len returns the length of the underlying slice
func (a *VariableArray) Len() int {
	return len(*a)
}

// LenNonEmpty returns the number of non-nil elements in the underlying slice
func (a *VariableArray) LenNonEmpty() int {
	count := 0
	for _, v := range a.Iterator() {
		if v != nil {
			count++
		}
	}
	return count
}

func (a *VariableArray) String() string {
	s := []string{}
	for _, v := range *a {
		s = append(s, fmt.Sprintf("%v", v))
	}
	return "[" + strings.Join(s, " ") + "]"
}
