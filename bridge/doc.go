/*
Package bridge provides a clean integration between Go's testing package and scripttest.

The bridge package allows scripttest tests to be seamlessly integrated into standard Go tests,
making it easy to combine both traditional Go unit tests and scripttest-based integration tests
in the same test suite.

# Basic Usage

Here's a simple example of using the bridge package:

	package mypackage_test

	import (
		"testing"

		"github.com/tmc/scripttestutil/bridge"
	)

	func TestWithScripttest(t *testing.T) {
		// Run all scripttest tests in the testdata directory
		opts := bridge.DefaultOptions()
		opts.Verbose = testing.Verbose()
		bridge.RunDir(t, "testdata", opts)
	}

# Key Features

1. Seamless integration with Go's testing package
2. Support for running individual test files or patterns of files
3. Configuration options for Docker, snapshots, environments, etc.
4. Preserves test output and error messages in proper Go test format
5. Sub-tests for each scripttest file for better organization and parallel execution

# Configuration Options

The TestOptions struct provides several configuration options:

- Pattern: Glob pattern to match test files
- UseDocker: Whether to run tests in Docker containers
- DockerImage: Docker image to use when UseDocker is true
- UpdateSnapshots: Whether to update snapshots
- Verbose: Enable verbose output
- EnvVars: Additional environment variables to pass to tests
- SnapshotDir: Directory for test snapshots

# Running Methods

The package provides several ways to run scripttest tests:

1. RunFile: Run a single scripttest file
2. RunPattern: Run all scripttest files matching a pattern
3. RunDir: Run all scripttest files in a directory
4. Runner.Run: Run tests using a configured Runner instance
5. Runner.RunTest: Run a single test using a configured Runner instance

# Example with Custom Options

	func TestWithCustomOptions(t *testing.T) {
		opts := bridge.DefaultOptions()
		opts.Verbose = true
		opts.UpdateSnapshots = true
		opts.EnvVars = map[string]string{
			"TEST_ENV": "custom-value",
			"DEBUG": "true",
		}
		
		// Run tests matching a specific pattern
		bridge.RunPattern(t, "testdata/advanced/*.txt", opts)
	}

# Docker Support

To run tests in Docker containers:

	func TestInDocker(t *testing.T) {
		opts := bridge.DefaultOptions()
		opts.UseDocker = true
		opts.DockerImage = "ubuntu:latest"
		
		bridge.RunDir(t, "testdata/docker", opts)
	}

# Parallel Testing

The bridge package supports parallel test execution:

	func TestParallel(t *testing.T) {
		t.Parallel() // Enable parallel testing
		
		opts := bridge.DefaultOptions()
		bridge.RunDir(t, "testdata", opts)
	}
*/
package bridge