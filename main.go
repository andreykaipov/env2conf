package main

import (
	"flag"
	"fmt"
	"os"

	env2star "github.com/andreykaipov/env2star/pkg"
)

func main() {
	printVersion := false
	prefix := ""
	output := ""

	flag.BoolVar(&printVersion, "version", false, "Print version and exit")
	flag.StringVar(&prefix, "prefix", "config", "A comma-delimited list of prefixes to parse env vars on")
	flag.StringVar(&output, "output", "json", "The output format, e.g. json, yaml, toml")
	flag.Parse()

	if printVersion {
		fmt.Printf("env2star %s (Git SHA: %s)\n", version, gitsha)
		os.Exit(0)
	}

	env2star.Run(prefix, output)
}
