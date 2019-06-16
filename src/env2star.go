package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var prefix string
var output string
var printVersion bool

func init() {
	flag.StringVar(&prefix, "prefix", "config", "A comma-delimited list of prefixes to parse env vars on")
	flag.StringVar(&output, "output", "json", "The output format, e.g. json, yaml, toml")
	flag.BoolVar(&printVersion, "version", false, "Print version and exit")
	flag.Parse()
}

func main() {
	if printVersion {
		fmt.Printf("env2star %s (Git SHA: %s)\n", version, gitsha)
		os.Exit(0)
	}

	prefixes := strings.Split(prefix, ",")
	parsed := map[string]interface{}{}
	mapsAsArrays := map[string]map[string]interface{}{}

	env := os.Environ()
	sort.Sort(sort.Reverse(sort.StringSlice(env)))

	for _, line := range env {
		kv := strings.SplitN(line, "=", 2)
		k := kv[0]
		v := kv[1]
		if !any(prefixes, func(x string) bool { return strings.HasPrefix(k, x) }) {
			continue
		}
		if err := parseLineAsMap(k, parsed, v, mapsAsArrays); err != nil {
			fmt.Printf("Failed parsing %s: %s\n", line, err)
			os.Exit(1)
		}
	}

	// cleanup
	for k, v := range mapsAsArrays {
		parts := strings.Split(k, ".")
		ogPart := parts[len(parts)-1]
		delete(v, ogPart)
	}

	switch output {
	case "toml":
		printTOML(parsed)
	case "yaml":
		fmt.Print("---")
		printYAML(parsed, 0)
	default:
		printJSON(parsed, 0)
	}
}

var bitSize = 32 + int(^uintptr(0)>>63<<5)

// cleans up an env var value for use in our parsed maps and arrays
func sanitize(s string) interface{} {
	switch s {
	case "{}":
		return map[string]interface{}{}
	case "[]":
		return &VariableArray{}
	case "null":
		return nil
	default:
		if n, err := strconv.Atoi(s); err == nil {
			return n
		}
		if z, err := strconv.ParseFloat(s, bitSize); err == nil {
			return z
		}
		if b, err := strconv.ParseBool(s); err == nil {
			return b
		}
		if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
			s = s[1 : len(s)-1]
		}
		return s
	}
}

// Given a line "a1.a2.a3...an", an existing map m, and a value v,
// this function creates entries in m such that m["a1"]["a2"]...["an"]=v.
// If an ai has a '[', it is parsed as an array as continues to be parsed as a map.
// e.g. "a.b[5[3.c" will create m["a"]["b"][5][3]["c"] and m["a"]["b[5[3"]["c"]
// This is done because we use the latter as a quick reference to what "a.b[5[3" might be.
// To avoid this extra key in the output, we keep track of these funky maps in a map to delete them later.
func parseLineAsMap(line string, m map[string]interface{}, v string, mapsAsArrays map[string]map[string]interface{}) error {
	parts := strings.Split(line, ".")

	for i, part := range parts[:len(parts)-1] {
		// to store any maps parsed as arrays so we can delete their references later
		fqkey := strings.Join(parts[0:i+1], ".")

		// check for existence and then the proper type
		if _, ok := m[part]; !ok {
			m[part] = map[string]interface{}{}
		}
		if _, ok := m[part].(map[string]interface{}); !ok {
			// wanted to parse the part as a map, but it was already parsed as something else
			return fmt.Errorf("%s is both map and %s", part, typeOf(m[part]))
		}

		next := m[part].(map[string]interface{})
		if strings.ContainsAny(part, "[") {
			if err := parsePartAsArray(part, m, next); err != nil {
				return err
			}
			mapsAsArrays[fqkey] = m
		}
		m = next
	}

	lastPart := parts[len(parts)-1]
	cleanVal := sanitize(v)

	if val, ok := m[lastPart]; ok {
		// wanted to set line to the terminating value, but it was already parsed as something else
		return fmt.Errorf("%s is both %s and %s", lastPart, typeOf(cleanVal), typeOf(val))
	}
	if strings.ContainsAny(lastPart, "[") {
		if err := parsePartAsArray(lastPart, m, cleanVal); err != nil {
			return err
		}
	} else {
		m[lastPart] = cleanVal
	}
	return nil
}

// Given a part "p[n0][n1]...[nk]", an existing map m, and a value v,
// this function creates entries in m such that m[p][n0][n1]...[nk]=v
func parsePartAsArray(part string, m map[string]interface{}, v interface{}) error {
	trimmed := part
	trimmed = strings.ReplaceAll(trimmed, "[", ".")
	trimmed = strings.ReplaceAll(trimmed, "]", "")
	subparts := strings.Split(trimmed, ".")

	key, indices := subparts[0], subparts[1:]
	if _, ok := m[key]; !ok {
		m[key] = &VariableArray{}
	}
	if _, ok := m[key].(*VariableArray); !ok {
		// wanted to parse this part as an array, but it was already parsed as something else
		return fmt.Errorf("%s is both array and %s", key, typeOf(m[key]))
	}
	current := m[key].(*VariableArray)

	nums := make([]int, len(indices))
	for j, index := range indices {
		i, err := strconv.Atoi(index)
		if err != nil {
			return fmt.Errorf("found non-number %s in array index", index)
		}
		nums[j] = i
	}

	for _, i := range nums[:len(nums)-1] {
		// similar to the map traversal, check for existence and the proper type until we get to the end
		if val := current.Get(i); val == nil {
			current.Set(i, &VariableArray{})
		}
		if _, ok := current.Get(i).(*VariableArray); !ok {
			// wanted to parse the part as an array, but it was already parsed as something else
			return fmt.Errorf("%s is both array and %s", part, typeOf(current.Get(i)))
		}
		current = current.Get(i).(*VariableArray)
	}

	lastIndex := nums[len(nums)-1]
	if val := current.Get(lastIndex); val != nil {
		// If m[p][n0][n1]...[nk] already exists, then there must exist a j>=k such that there exists
		// a part "p[n0][n1]...[nj]" sharing common descendents with "p[n0][n1]...[nk]" in our env vars
		// that was already parsed. Proof: If j<k, then we would have already failed in the loop above
		// expecting m[p][n0][n1]...[nj] to be an array.
		//
		// A nice corollary to that is if our environment is sorted, this should never happen! :-)
		//
		// However, if j>k, we can fail immediately since m[p][n0][n1]...[nk] cannot both be an array and
		// a terminating value; e.g. a[0]=1 a[0][0]=2. More interestingly, if j=k, then we'd either have
		// a duplicate env var (can't happen), or m[p][n0][n1]...[nk] was parsed as a map; e.g. a[0].a=1 a[0].b=2.
		// So, fail on non-maps only:
		if _, ok := val.(map[string]interface{}); !ok {
			// wanted to set the line to the value, but it was already parsed as something else
			return fmt.Errorf("%s is both %s and %s", part, typeOf(v), typeOf(val))
		}
	}
	current.Set(lastIndex, v)
	return nil
}

func any(xs []string, f func(string) bool) bool {
	for _, x := range xs {
		if f(x) {
			return true
		}
	}
	return false
}

// returns the simple type of an object
func typeOf(v interface{}) string {
	switch t := fmt.Sprintf("%T", v); t {
	case fmt.Sprintf("%T", &VariableArray{}):
		return "array"
	case fmt.Sprintf("%T", map[string]interface{}{}):
		return "map"
	default:
		return t
	}
}
