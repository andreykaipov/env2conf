package main

import (
	"testing"
)

func TestSet_Happy_InOrder(t *testing.T) {
	a := &VariableArray{}
	a.Set(0, "q")
	a.Set(1, "w")
	a.Set(2, "e")
	if a.Len() != 3 {
		t.Errorf("Expected length of 3; got: %v", a.Len())
	}
	if a.LenNonEmpty() != 3 {
		t.Errorf("Expected non empty length of 3; got: %v", a.LenNonEmpty())
	}
}

func TestSet_Happy_OutOfOrder(t *testing.T) {
	a := &VariableArray{}
	a.Set(30, "q")
	a.Set(11, "w")
	a.Set(18, "e")
	a.Set(1, "r")
	if a.Len() != 31 {
		t.Errorf("Expected length of 30; got %v", a.Len())
	}
	if a.Get(30) != "q" {
		t.Errorf("Expected q; got %v", a.Get(30))
	}
	if a.Get(11) != "w" {
		t.Errorf("Expected e; got %v", a.Get(11))
	}
	if a.Get(18) != "e" {
		t.Errorf("Expected e; got %v", a.Get(18))
	}
	if a.Get(1) != "r" {
		t.Errorf("Expected e; got %v", a.Get(1))
	}
	if a.LenNonEmpty() != 4 {
		t.Errorf("Expected non empty length of 4; got: %v", a.LenNonEmpty())
	}
}

func TestSet_Happy_Nested(t *testing.T) {
	a := &VariableArray{}
	a.Set(2, &VariableArray{})
	a.Get(2).(*VariableArray).Set(1, "d")
	a.Get(2).(*VariableArray).Set(0, "c")
	a.Set(0, &VariableArray{})
	a.Get(0).(*VariableArray).Set(1, "b")
	a.Get(0).(*VariableArray).Set(0, "a")
	if a.Len() != 3 {
		t.Errorf("Expected length of 3; got: %v", a.Len())
	}
	if a.LenNonEmpty() != 2 {
		t.Errorf("Expected non empty length of 3; got: %v", a.LenNonEmpty())
	}
	if a.Get(0).(*VariableArray).Len() != 2 {
		t.Errorf("Expected length of 2 for nested array")
	}
	if a.Get(2).(*VariableArray).Len() != 2 {
		t.Errorf("Expected length of 2 for nested array")
	}
}

func TestJSONMarshal_Happy(t *testing.T) {
	a := &VariableArray{}
	a.Set(0, "a")
	a.Set(2, "b")
	a.Set(3, &VariableArray{})
	a.Get(3).(*VariableArray).Set(1, "c")
	a.Get(3).(*VariableArray).Set(2, "d")
	a.Get(3).(*VariableArray).Set(4, map[string]interface{}{"a": 1, "b": "a"})
	json, err := a.MarshalJSON()
	if err != nil {
		t.Errorf("Expected JSON marshalling to not fail")
	}
	if string(json) != `["a",null,"b",[null,"c","d",null,{"a":1,"b":"a"}]]` {
		t.Errorf("Expected JSON marshalling to produce the correct result")
	}
}
