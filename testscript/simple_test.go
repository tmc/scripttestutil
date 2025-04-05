package testscript_test

import (
	"testing"

	"github.com/tmc/scripttestutil/testscript"
)

// TestSimple runs a simple test using the testscript package
func TestSimple(t *testing.T) {
	// Set up options
	opts := testscript.DefaultOptions()
	opts.Verbose = testing.Verbose()

	// Run the test file
	testscript.RunFile(t, "testdata/simple.txt", opts)
}