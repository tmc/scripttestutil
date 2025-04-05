// A simple example that shows how to use expect commands in scripttest tests
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

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
	
	// Create a minimal testing.T implementation for the runner
	t := &testingT{
		verbose: *verbose,
	}
	
	// Run the tests
	if *singleTest != "" {
		// Run a single test file
		testPath := *singleTest
		if !filepath.IsAbs(testPath) {
			testPath = filepath.Join(*testDir, testPath)
		}
		
		fmt.Printf("Running test: %s\n", testPath)
		testscript.RunFile(t, testPath, opts)
	} else {
		// Run all tests in the directory
		fmt.Printf("Running all tests in: %s\n", *testDir)
		testscript.RunDir(t, *testDir, opts)
	}
	
	// Report results
	if t.failed {
		fmt.Printf("FAIL: %d tests failed\n", t.failCount)
		os.Exit(1)
	} else {
		fmt.Printf("PASS: All tests passed\n")
	}
}

// testingT is a minimal implementation of testing.T for use in the example
type testingT struct {
	verbose   bool
	failed    bool
	failCount int
}

func (t *testingT) Fail() {
	t.failed = true
	t.failCount++
}

func (t *testingT) Failed() bool {
	return t.failed
}

func (t *testingT) FailNow() {
	t.Fail()
	panic("test failed")
}

func (t *testingT) Name() string {
	return "expect-tester"
}

func (t *testingT) Fatalf(format string, args ...interface{}) {
	fmt.Printf("FAIL: "+format+"\n", args...)
	t.FailNow()
}

func (t *testingT) Errorf(format string, args ...interface{}) {
	fmt.Printf("ERROR: "+format+"\n", args...)
	t.Fail()
}

func (t *testingT) Log(args ...interface{}) {
	if t.verbose {
		fmt.Println(args...)
	}
}

func (t *testingT) Logf(format string, args ...interface{}) {
	if t.verbose {
		fmt.Printf(format+"\n", args...)
	}
}

func (t *testingT) Helper() {}

func (t *testingT) Skip(args ...interface{}) {
	fmt.Println("SKIP:", args)
}

func (t *testingT) Skipf(format string, args ...interface{}) {
	fmt.Printf("SKIP: "+format+"\n", args...)
}

func (t *testingT) Skipped() bool {
	return false
}

func (t *testingT) Run(name string, f func(t *testing.T)) bool {
	fmt.Printf("=== RUN %s\n", name)
	
	subT := &testingT{
		verbose: t.verbose,
	}
	
	defer func() {
		if r := recover(); r != nil {
			if r != "test failed" {
				panic(r)
			}
		}
		
		if subT.failed {
			fmt.Printf("--- FAIL: %s\n", name)
			t.failed = true
			t.failCount++
		} else {
			fmt.Printf("--- PASS: %s\n", name)
		}
	}()
	
	f(subT)
	return !subT.failed
}

func (t *testingT) Deadline() (deadline bool, ok bool) {
	return false, false
}