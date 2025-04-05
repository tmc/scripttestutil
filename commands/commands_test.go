package commands_test

import (
	"fmt"
	"testing"

	"github.com/tmc/scripttestutil/commands"
	"github.com/tmc/scripttestutil/testscript"
	"rsc.io/script"
)

// ExampleRegisterAll demonstrates how to register all command sets
func ExampleRegisterAll() {
	// Create a map for commands (this would typically be done by testscript)
	cmds := make(map[string]script.Cmd)
	
	// Register all command sets at once
	commands.RegisterAll(cmds)
	
	fmt.Println("Command sets registered successfully")
	
	// Output:
	// Command sets registered successfully
}

// ExampleRegisterExpect demonstrates how to register just the expect commands
func ExampleRegisterExpect() {
	// Create options for testscript
	opts := testscript.DefaultOptions()
	
	// Set up a hook to register only expect commands
	opts.SetupHook = func(cmds map[string]script.Cmd) {
		commands.RegisterExpect(cmds)
	}
	
	fmt.Println("Test script file:")
	fmt.Println(`
# Start a Python interpreter
expect:spawn python3
expect:expect ">>>" 5

# Run a simple program
expect:send "print('Hello, World!')"
expect:expect "Hello, World!" 5

# Exit Python
expect:send "exit()"
	`)
	
	// Output:
	// Test script file:
	// 
	// # Start a Python interpreter
	// expect:spawn python3
	// expect:expect ">>>" 5
	// 
	// # Run a simple program
	// expect:send "print('Hello, World!')"
	// expect:expect "Hello, World!" 5
	// 
	// # Exit Python
	// expect:send "exit()"
}

// TestCommandRegistration verifies that commands can be registered successfully
func TestCommandRegistration(t *testing.T) {
	// Create a map to hold commands
	cmds := make(map[string]script.Cmd)
	
	// Register all commands
	commands.RegisterAll(cmds)
	
	// Verify that expect commands were registered
	if _, ok := cmds["expect:spawn"]; !ok {
		t.Error("expect:spawn command was not registered")
	}
	if _, ok := cmds["expect:expect"]; !ok {
		t.Error("expect:expect command was not registered")
	}
	if _, ok := cmds["expect:send"]; !ok {
		t.Error("expect:send command was not registered")
	}
}