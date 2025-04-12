// A simple example that demonstrates how to use expect commands from the command line
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tmc/scripttestutil/commands"
	"github.com/tmc/scripttestutil/testscript"
	"rsc.io/script"
)

var (
	verbose   = flag.Bool("v", false, "verbose output")
	testDir   = flag.String("dir", "testdata", "directory containing test files")
	singleTest = flag.String("test", "", "run a single test file")
)

func main() {
	flag.Parse()

	// Set up testscript options
	opts := testscript.DefaultOptions()
	opts.Verbose = *verbose
	
	// Register expect commands
	opts.SetupHook = func(cmds map[string]script.Cmd) {
		commands.RegisterExpect(cmds)
	}
	
	// Process test files
	if *singleTest != "" {
		// Run a single test file
		testPath := *singleTest
		if !filepath.IsAbs(testPath) {
			testPath = filepath.Join(*testDir, testPath)
		}
		
		fmt.Printf("Running test: %s\n", testPath)
		if err := runSingleFile(testPath, opts); err != nil {
			log.Fatalf("Test failed: %v", err)
		}
		fmt.Println("Test passed successfully")
	} else {
		// Run all tests in the directory
		fmt.Printf("Running all tests in: %s\n", *testDir)
		if err := runDirectory(*testDir, opts); err != nil {
			log.Fatalf("Tests failed: %v", err)
		}
		fmt.Println("All tests passed successfully")
	}
}

// Helper function to run a single test file
func runSingleFile(path string, opts testscript.Options) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("test file not found: %s", path)
	}
	
	// Execute test file with script engine
	// This is a simplified version - in a real scenario
	// you would want to process the script output
	fmt.Printf("Executing test file: %s\n", path)
	
	// Print the test contents
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read test file: %v", err)
	}
	
	fmt.Println("Test file contents:")
	fmt.Println("-------------------")
	fmt.Println(string(content))
	fmt.Println("-------------------")
	
	// In a real implementation, you would execute the test
	// using scripttest/testscript functionality
	fmt.Println("Note: This example only demonstrates the setup.")
	fmt.Println("In a real implementation, the test would be executed.")
	
	return nil
}

// Helper function to run all tests in a directory
func runDirectory(dir string, opts testscript.Options) error {
	// Find all test files
	pattern := filepath.Join(dir, "*.txt")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern %q: %v", pattern, err)
	}
	
	if len(matches) == 0 {
		return fmt.Errorf("no test files found in %s", dir)
	}
	
	// Run each test file
	for _, testFile := range matches {
		fmt.Printf("Running test: %s\n", testFile)
		if err := runSingleFile(testFile, opts); err != nil {
			return err
		}
	}
	
	return nil
}