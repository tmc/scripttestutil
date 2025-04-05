package testscript

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAll runs all the tests in the testdata directory
func TestAll(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	// Set up options
	opts := DefaultOptions()
	opts.Verbose = testing.Verbose()
	
	// Set update snapshots if environment variable is set
	if os.Getenv("UPDATE_SNAPSHOTS") == "1" {
		opts.UpdateSnapshots = true
	}

	// Find all test directories under testdata
	testDirs, err := findTestDirs("testdata")
	if err != nil {
		t.Fatalf("Failed to find test directories: %v", err)
	}

	// Run tests for each directory
	for _, dir := range testDirs {
		dirName := filepath.Base(dir)
		t.Run(dirName, func(t *testing.T) {
			RunDir(t, dir, opts)
		})
	}
}

// findTestDirs recursively finds directories containing *.txt test files
func findTestDirs(root string) ([]string, error) {
	var dirs []string
	
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip the snapshot directory
		if info.IsDir() && info.Name() == "__snapshots__" {
			return filepath.SkipDir
		}
		
		// If it's a file and has a .txt extension, add its directory
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			dir := filepath.Dir(path)
			
			// Check if we already have this directory
			alreadyAdded := false
			for _, d := range dirs {
				if d == dir {
					alreadyAdded = true
					break
				}
			}
			
			if !alreadyAdded {
				dirs = append(dirs, dir)
			}
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}
	
	return dirs, nil
}

// TestWithTag allows you to run a specific test category using the -run flag
// Example: go test -run TestWithTag/gotools
func TestWithTag(t *testing.T) {
	// Define test categories
	categories := map[string]string{
		"gotools": "testdata/gotools",
		"spinner": "testdata/spinner",
	}
	
	// Run each category as a subtest
	for name, dir := range categories {
		t.Run(name, func(t *testing.T) {
			opts := DefaultOptions()
			opts.Verbose = testing.Verbose()
			
			RunDir(t, dir, opts)
		})
	}
}

// TestGoTools runs all the Go tools tests using the testscript package
func TestGoTools(t *testing.T) {
	// Set up options
	opts := DefaultOptions()
	opts.Verbose = testing.Verbose()
	
	// Set update snapshots if environment variable is set
	if os.Getenv("UPDATE_SNAPSHOTS") == "1" {
		opts.UpdateSnapshots = true
	}

	// Run all tests in the gotools directory
	RunDir(t, "testdata/gotools", opts)
}

// TestSpecificGoTool runs a specific Go tool test
func TestSpecificGoTool(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	opts := DefaultOptions()
	opts.Verbose = testing.Verbose()
	
	// Run just the gofmt test
	RunFile(t, "testdata/gotools/simple_fmt_test.txt", opts)
}

// TestSpinner runs all the spinner tests using the testscript package
func TestSpinner(t *testing.T) {
	// Set up options
	opts := DefaultOptions()
	opts.Verbose = testing.Verbose()
	
	// Set update snapshots if environment variable is set
	if os.Getenv("UPDATE_SNAPSHOTS") == "1" {
		opts.UpdateSnapshots = true
	}

	// Run all tests in the spinner directory
	RunDir(t, "testdata/spinner", opts)
}

// TestIndividualSpinnerTests runs specific spinner tests
func TestIndividualSpinnerTests(t *testing.T) {
	// Sub-tests for each spinner test
	t.Run("Basic", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping in short mode")
		}
		
		opts := DefaultOptions()
		opts.Verbose = testing.Verbose()
		RunFile(t, "testdata/spinner/basic_test.txt", opts)
	})
	
	t.Run("Concurrent", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping in short mode")
		}
		
		opts := DefaultOptions()
		opts.Verbose = testing.Verbose()
		RunFile(t, "testdata/spinner/concurrent_test.txt", opts)
	})
	
	t.Run("Custom", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping in short mode")
		}
		
		opts := DefaultOptions()
		opts.Verbose = testing.Verbose()
		RunFile(t, "testdata/spinner/custom_spinner_test.txt", opts)
	})
}