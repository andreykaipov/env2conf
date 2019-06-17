package env2star

import (
	"fmt"
	//"strings"
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

func padPrintf(indent int, s ...interface{}) {
	if indent < 0 {
		indent = 0
	}
	fmt.Printf("%*.s", indent, " ")
	fmt.Printf(s[0].(string), s[1:]...)
}

// Any top-level key in a TOML map whose value is not a map nor an array of maps
// must be printed first before any other key/value. Otherwise, it'll be ambiguous
// which map these keys belong to. TOML is great.
func printTOML(data map[string]interface{}) error {
	for k, v := range data {
		switch val := v.(type) {
		case *VariableArray:
			if !val.IsWellDefined() {
				return fmt.Errorf("TOML arrays can't have mixed data types: %v", val)
			}
			if val.GetType() != fmt.Sprintf("%T", map[string]interface{}{}) {
				fmt.Printf("%s = ", k)
				printTOMLInline(v, 0)
				fmt.Println()
				delete(data, k)
			}
		case map[string]interface{}:
			// nothing yet
		default:
			fmt.Printf("%s = ", k)
			printTOMLInline(v, 0)
			fmt.Println()
			delete(data, k)
		}
	}
	return printTOMLHelper(data, "", 0)
}

func printTOMLHelper(data interface{}, previous string, indent int) error {
	switch t := data.(type) {
	case map[string]interface{}:
		for k, v := range t {
			fqkey := k
			if len(previous) != 0 {
				fqkey = fmt.Sprintf("%s.%s", previous, k)
			}
			switch val := v.(type) {
			case *VariableArray:
				if !val.IsWellDefined() {
					return fmt.Errorf("TOML arrays can't have mixed data types: %v", val)
				}
				if val.GetType() == fmt.Sprintf("%T", map[string]interface{}{}) {
					for _, v := range val.Iterator() {
						fmt.Printf("%*.s[[%s]]\n", indent, " ", fqkey)
						if err := printTOMLHelper(v, fqkey, indent+2); err != nil {
							return err
						}
						fmt.Println()
					}
				} else {
					fmt.Printf("%*.s%s = ", indent-2, " ", k)
					//padPrintf(indent-2,/ "%s = ", k)
					if err := printTOMLHelper(v, fqkey, indent-2); err != nil {
						return err
					}
					fmt.Println()
				}
			case map[string]interface{}:
				fmt.Printf("%*.s[%s]\n", indent, " ", fqkey)
				if err := printTOMLHelper(v, fqkey, indent+2); err != nil {
					return err
				}
			default:
				fmt.Printf("%*.s%s = ", indent, " ", k)
				if err := printTOMLHelper(v, fqkey, indent+2); err != nil {
					return err
				}
				fmt.Println()
			}
		}
	case *VariableArray:
		if !t.IsWellDefined() {
			return fmt.Errorf("TOML arrays can't have mixed data types: %v", t)
		}
		fmt.Print("[")
		for _, v := range t.Iterator() {
			fmt.Printf("\n%*.s", indent+2, " ")
			switch val := v.(type) {
			case map[string]interface{}:
				printTOMLInline(val, 0)
			default:
				if err := printTOMLHelper(val, previous, indent+2); err != nil {
					return err
				}
			}
			fmt.Println(",")
		}
		fmt.Printf("%*.s]", indent, " ")
	case nil:
		fmt.Printf("null")
	case string:
		fmt.Printf("%q", t)
	default:
		fmt.Printf("%v", data)
	}
	return nil
}

func printTOMLInline(data interface{}, indent int) {
	switch t := data.(type) {
	case map[string]interface{}:
		fmt.Printf("{")
		i := 0
		for k, v := range t {
			fmt.Printf("%s = ", k)
			printTOMLInline(v, indent)
			if i < len(t)-1 {
				fmt.Print(", ")
			}
			i++
		}
		fmt.Printf("}")
	case *VariableArray:
		// only arrays support multiline
		multiline := false
		switch vt := t.GetType(); vt {
		case fmt.Sprintf("%T", map[string]interface{}{}):
			multiline = true
		case fmt.Sprintf("%T", &VariableArray{}):
			multiline = true
		}
		fmt.Print("[")
		if multiline {
			fmt.Println()
		}
		for i, v := range t.Iterator() {
			if multiline {
				fmt.Printf("%*.s", indent+2, " ")
			}
			printTOMLInline(v, indent+2)
			if i < t.Len()-1 {
				fmt.Printf(",")
			}
			if multiline {
				fmt.Println()
			}
		}
		if multiline {
			fmt.Printf("%*.s", indent, " ")
		}
		fmt.Print("]")
	case nil:
		fmt.Print("null")
	case string:
		fmt.Printf("%q", t)
	default:
		fmt.Print(data)
	}
}
