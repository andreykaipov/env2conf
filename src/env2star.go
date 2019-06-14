package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var prefix string
var output string

func init() {
	flag.StringVar(&prefix, "prefix", "config", "A comma-delimited list of prefixes to parse env vars on")
	flag.StringVar(&output, "output", "json", "The output format, e.g. json, yaml, toml")
	flag.Parse()
}

func main() {
	prefixes := strings.Split(prefix, ",")
	parsed := map[string]interface{}{}
	mapsAsArrays := map[string]map[string]interface{}{}

	for _, line := range os.Environ() {
		kv := strings.SplitN(line, "=", 2)
		k := kv[0]
		v := kv[1]
		if len(prefix) == 0 || !any(prefixes, func(x string) bool { return strings.HasPrefix(k, x) }) {
			continue
		}
		if err := parseLineAsMap(k, parsed, v, mapsAsArrays); err != nil {
			fmt.Println(err)
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
	default:
		if z, err := strconv.ParseFloat(s, bitSize); err == nil {
			return z
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
		// for error messages and to store any maps parsed as arrays so we can delete their references later
		fqkey := strings.Join(parts[0:i+1], ".")

		// check for existence and then the proper type
		if _, ok := m[part]; !ok {
			m[part] = map[string]interface{}{}
		}
		if _, ok := m[part].(map[string]interface{}); !ok {
			return fmt.Errorf("wanted to parse a map out of %s, but it was already parsed as a non-map: %v", fqkey, m[part])
		}

		next := m[part].(map[string]interface{})
		if strings.ContainsAny(part, "[") {
			if err := parsePartAsArray(part, m, next); err != nil {
				return fmt.Errorf("failed to parse array out of %s: %s", fqkey, err)
			}
			mapsAsArrays[fqkey] = m
		}
		m = next
	}

	lastPart := parts[len(parts)-1]
	if _, ok := m[lastPart]; ok {
		return fmt.Errorf("wanted to set %s=%s, but %s was already parsed as a map", line, v, line)
	}
	if strings.ContainsAny(lastPart, "[") {
		if err := parsePartAsArray(lastPart, m, sanitize(v)); err != nil {
			return fmt.Errorf("failed to parse array out of %s: %s", line, err)
		}
	} else {
		m[lastPart] = sanitize(v)
	}
	return nil
}

// Given a part "p[n0][n1]...[nk]", an existing map m, and a value v,
// this function creates entries in m such that m[p][n0][n1]...[nk]=v
func parsePartAsArray(part string, m map[string]interface{}, v interface{}) error {
	part = strings.ReplaceAll(part, "[", ".")
	part = strings.ReplaceAll(part, "]", "")
	subparts := strings.Split(part, ".")

	key, indices := subparts[0], subparts[1:]
	if _, ok := m[key]; !ok {
		m[key] = &VariableArray{}
	}
	if _, ok := m[key].(*VariableArray); !ok {
		return fmt.Errorf("expected %s to be an array, but got: %v", key, m[key])
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

	for j, i := range nums[:len(nums)-1] {
		// similar to the map traversal, check for existence and the proper type until we get to the end
		if val := current.Get(i); val == nil {
			current.Set(i, &VariableArray{})
		}
		if _, ok := current.Get(i).(*VariableArray); !ok {
			return fmt.Errorf("wanted to parse an array out of %s[%s], but it was: %v", key, strings.Join(indices[0:j+1], "]["), current.Get(i))
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
		// Now, if j>k, we can fail immediately since m[p][n0][n1]...[nk] cannot both be an array and
		// a terminating value; e.g. a[0]=1 a[0][0]=2. More interestingly, if j=k, then we'd either have
		// a duplicate env var (can't happen), or m[p][n0][n1]...[nk] was parsed as a map; e.g. a[0].a=1 a[0].b=2.
		// So, fail on non-maps only:
		if _, ok := val.(map[string]interface{}); !ok {
			return fmt.Errorf("wanted to set %s[%s]=%v; but %v is already there", key, strings.Join(indices, "]["), v, val)
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
