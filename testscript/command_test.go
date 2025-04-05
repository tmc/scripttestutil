package testscript_test

import (
	"fmt"
	"testing"

	"github.com/tmc/scripttestutil/commands"
	"github.com/tmc/scripttestutil/testscript"
	"rsc.io/script"
)

// ExampleOptions_SetupHook demonstrates how to use the SetupHook to register custom commands
func ExampleOptions_SetupHook() {
	// Create custom options
	opts := testscript.DefaultOptions()
	
	// Register expect commands using the SetupHook
	opts.SetupHook = func(cmds map[string]script.Cmd) {
		commands.RegisterExpect(cmds)
	}
	
	fmt.Println("SetupHook configured to register expect commands")
	
	// Output:
	// SetupHook configured to register expect commands
}

// ExampleRunWithCommands demonstrates how to run tests with custom commands
func ExampleRunWithCommands() {
	// This won't actually run in the example, just showing the pattern
	fmt.Println("Running directory with custom commands:")
	fmt.Println(`
	import (
		"testing"
		"github.com/tmc/scripttestutil/commands"
		"github.com/tmc/scripttestutil/testscript"
		"rsc.io/script"
	)

	func TestWithExpectCommands(t *testing.T) {
		opts := testscript.DefaultOptions()
		opts.Verbose = testing.Verbose()
		
		// Register expect commands
		opts.SetupHook = func(cmds map[string]script.Cmd) {
			commands.RegisterExpect(cmds)
		}
		
		// Run all tests in the expect directory
		testscript.RunDir(t, "testdata/expect", opts)
	}
	`)
	
	// Output:
	// Running directory with custom commands:
	// 
	// 	import (
	// 		"testing"
	// 		"github.com/tmc/scripttestutil/commands"
	// 		"github.com/tmc/scripttestutil/testscript"
	// 		"rsc.io/script"
	// 	)
	// 
	// 	func TestWithExpectCommands(t *testing.T) {
	// 		opts := testscript.DefaultOptions()
	// 		opts.Verbose = testing.Verbose()
	// 		
	// 		// Register expect commands
	// 		opts.SetupHook = func(cmds map[string]script.Cmd) {
	// 			commands.RegisterExpect(cmds)
	// 		}
	// 		
	// 		// Run all tests in the expect directory
	// 		testscript.RunDir(t, "testdata/expect", opts)
	// 	}
	// 
}

// TestExpectCommands runs the expect command examples
func TestExpectCommands(t *testing.T) {
	// Skip this test in normal builds to avoid external dependencies
	// but still allow it to be run explicitly with go test -run=TestExpectCommands
	if testing.Short() {
		t.Skip("Skipping expect command tests in short mode")
	}
	
	// Create custom options with expect commands
	opts := testscript.DefaultOptions()
	opts.Verbose = testing.Verbose()
	
	// Register expect commands using the SetupHook
	opts.SetupHook = func(cmds map[string]script.Cmd) {
		commands.RegisterExpect(cmds)
	}
	
	// Run expect test
	testscript.RunDir(t, "../testdata/expect", opts)
}

// TestMultipleCommandSets demonstrates using multiple command sets
func TestMultipleCommandSets(t *testing.T) {
	// Skip this test in normal builds to avoid external dependencies
	// but still allow it to be run explicitly
	if testing.Short() {
		t.Skip("Skipping multiple command sets test in short mode")
	}
	
	// Create custom options with all command sets
	opts := testscript.DefaultOptions()
	opts.Verbose = testing.Verbose()
	
	// Register all command sets
	opts.SetupHook = func(cmds map[string]script.Cmd) {
		commands.RegisterAll(cmds)
	}
	
	// Run tests
	testscript.RunDir(t, "../testdata", opts)
}