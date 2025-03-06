// Package main demonstrates testing a simple CLI app with scripttest bridge
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	name    string
	verbose bool
	count   int
)

func init() {
	flag.StringVar(&name, "name", "World", "Name to greet")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.IntVar(&count, "count", 1, "Number of times to print the greeting")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	greeting := fmt.Sprintf("Hello, %s!", name)

	if verbose {
		fmt.Fprintf(os.Stderr, "About to print greeting %d times\n", count)
	}

	for i := 0; i < count; i++ {
		fmt.Println(greeting)
	}

	if flag.NArg() > 0 {
		fmt.Printf("Extra arguments: %s\n", strings.Join(flag.Args(), ", "))
	}

	return nil
}