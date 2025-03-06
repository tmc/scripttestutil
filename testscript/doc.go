/*
Package testscript provides integration between Go's testing package and scripttest.

The testscript package makes it easy to run scripttest tests as part of standard Go
test suites, combining traditional Go tests and scripttest-based integration tests.

# Basic Usage

Here's a simple example:

	package mypackage_test

	import (
		"testing"

		"github.com/tmc/scripttestutil/testscript"
	)

	func TestCLI(t *testing.T) {
		// Run all scripttest tests in the testdata directory
		opts := testscript.DefaultOptions()
		opts.Verbose = testing.Verbose()
		testscript.RunDir(t, "testdata", opts)
	}

# Features

The testscript package supports:

- Running individual test files or patterns
- Docker-based tests
- Snapshot creation and verification
- Custom environment variables
- Organized sub-tests for each scripttest file

# Running Tests

The package provides several ways to run tests:

	// Run a single test file
	testscript.RunFile(t, "testdata/basic.txt", opts)

	// Run all tests in a directory
	testscript.RunDir(t, "testdata", opts)

	// Run tests matching a pattern
	testscript.Run(t, "testdata/cli/*.txt", opts)

# Options

Configure test behavior with Options:

	opts := testscript.DefaultOptions()
	opts.Verbose = true // Enable verbose output
	opts.UpdateSnapshots = true // Update test snapshots
	opts.EnvVars = map[string]string{ // Set environment variables
		"DEBUG": "true",
		"API_URL": "http://localhost:8080",
	}

# Example with Docker

To run tests in Docker containers:

	func TestInDocker(t *testing.T) {
		opts := testscript.DefaultOptions()
		opts.UseDocker = true
		opts.DockerImage = "golang:1.22"
		testscript.RunDir(t, "testdata/docker", opts)
	}

# Parallel Testing

The package works with Go's parallel testing:

	func TestParallel(t *testing.T) {
		t.Parallel() // Enable parallel testing
		testscript.RunDir(t, "testdata", testscript.DefaultOptions())
	}
*/
package testscript