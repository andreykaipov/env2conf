package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// VariableArray represents an array that does not fail
// whenever a get or set is performed on an out-of-bounds index.
// Like JavaScript arrays <3 e.g.: a=[]; a[10]==undefined
type VariableArray struct {
	slice []interface{}
}

// oh no so many new slices
func (a *VariableArray) expandIfNecessary(i int) {
	if i+1 > len(a.slice) {
		expanded := make([]interface{}, i+1)
		copy(expanded, a.slice)
		a.slice = expanded
	}
}

func (a *VariableArray) Set(i int, val interface{}) {
	a.expandIfNecessary(i)
	a.slice[i] = val
}

func (a *VariableArray) Get(i int) interface{} {
	if i+1 > len(a.slice) {
		return nil
	}
	return a.slice[i]
}

func (a *VariableArray) Iterator() []interface{} {
	return a.slice
}

func (a *VariableArray) Len() int {
	return len(a.slice)
}

func (a *VariableArray) LenNonEmpty() int {
	count := 0
	for _, v := range a.slice {
		if v != nil {
			count++
		}
	}
	return count
}

func (a *VariableArray) String() string {
	s := []string{}
	for _, v := range a.slice {
		s = append(s, fmt.Sprintf("%v", v))
	}
	return "[" + strings.Join(s, " ") + "]"
}

func (a *VariableArray) MarshalJSON() ([]byte, error) {
	b := []byte{'['}
	count := len(a.slice)
	for i, v := range a.slice {
		marshalled, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		b = append(b, marshalled...)
		if i < count-1 {
			b = append(b, ',')
		}
	}
	b = append(b, ']')
	return b, nil
}
