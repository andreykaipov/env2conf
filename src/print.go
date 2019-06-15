package main

import (
	"fmt"
)

func printJSON(data interface{}, indent int) {
	switch t := data.(type) {
	case map[string]interface{}:
		fmt.Println("{")
		i := 0
		for k, v := range t {
			fmt.Printf("%*.s%q: ", indent+2, " ", k)
			printJSON(v, indent+2)
			if i < len(t)-1 {
				fmt.Printf(",\n")
			}
			i++
		}
		fmt.Printf("\n%*.s}", indent, " ")
	case *VariableArray:
		fmt.Println("[")
		for i, v := range t.Iterator() {
			fmt.Printf("%*.s", indent+2, " ")
			printJSON(v, indent+2)
			if i < t.Len()-1 {
				fmt.Printf(",\n")
			}
		}
		fmt.Printf("\n%*.s]", indent, " ")
	case nil:
		fmt.Print("null")
	case string:
		fmt.Printf("%q", data)
	default:
		fmt.Print(data)
	}
}

func printYAML(data interface{}, indent int) {
	switch t := data.(type) {
	case map[string]interface{}:
		if len(t) == 0 {
			fmt.Println("{}")
		} else {
			fmt.Println()
		}
		for k, v := range t {
			fmt.Printf("%*.s%s: ", indent, " ", k)
			printYAML(v, indent+2)
		}
	case *VariableArray:
		if t.Len() == 0 {
			fmt.Println("[]")
		} else {
			fmt.Println()
		}
		for _, v := range t.Iterator() {
			fmt.Printf("%*.s- ", indent, " ")
			printYAML(v, indent+2)
		}
	case nil:
		fmt.Println("null")
	default:
		fmt.Println(data)
	}
}

func printTOML(data map[string]interface{}) {
	for k, v := range data {
		fmt.Printf("%s = ", k)
		printTOMLHelper(v)
		fmt.Println()
	}
}

func printTOMLHelper(data interface{}) {
	switch t := data.(type) {
	case map[string]interface{}:
		fmt.Printf("{")
		i := 0
		for k, v := range t {
			fmt.Printf("%s = ", k)
			printTOMLHelper(v)
			if i < len(t)-1 {
				fmt.Print(", ")
			}
			i++
		}
		fmt.Printf("}")
	case *VariableArray:
		fmt.Printf("[")
		for i, v := range t.Iterator() {
			printTOMLHelper(v)
			if i < t.Len()-1 {
				fmt.Print(", ")
			}
		}
		fmt.Printf("]")
	case nil:
		fmt.Print("null")
	case string:
		fmt.Printf("%q", t)
	default:
		fmt.Print(data)
	}
}
