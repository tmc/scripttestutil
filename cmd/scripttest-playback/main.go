package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <snapshot-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nPlays back a scripttest snapshot file\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	snapshotPath := flag.Arg(0)

	// Read snapshot file
	data, err := ioutil.ReadFile(snapshotPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading snapshot: %v\n", err)
		os.Exit(1)
	}

	// Parse snapshot content
	var content map[string]string
	if err := json.Unmarshal(data, &content); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid snapshot format: %v\n", err)
		os.Exit(1)
	}

	// Output the snapshot content
	if stdout, ok := content["stdout"]; ok {
		fmt.Print(stdout)
	}
	if stderr, ok := content["stderr"]; ok {
		fmt.Fprintf(os.Stderr, "%s", stderr)
	}
}
