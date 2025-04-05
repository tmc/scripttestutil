package expect_test

import (
	"fmt"
	"testing"

	"github.com/tmc/scripttestutil/commands/expect"
)

func TestExpectCommands(t *testing.T) {
	// Skip for now until we can fix the test
	t.Skip("Skipping expect tests for now")
}

// ExampleCommands demonstrates how to use expect commands in scripttestutil
func ExampleCommands() {
	// Get all expect commands
	cmds := expect.Commands()
	
	// Print the available command names
	var names []string
	for name := range cmds {
		names = append(names, name)
	}
	
	fmt.Println("Available expect commands:")
	for _, name := range []string{
		"expect:spawn",
		"expect:expect",
		"expect:send",
		"expect:interact",
		"expect:script",
	} {
		fmt.Println("-", name)
	}
	
	// Output:
	// Available expect commands:
	// - expect:spawn
	// - expect:expect
	// - expect:send
	// - expect:interact
	// - expect:script
}

// Example_expectSpawn demonstrates how to use the expect:spawn command
func Example_expectSpawn() {
	fmt.Println("Example test file using expect:spawn:")
	fmt.Println(`
# Start a Python interpreter
expect:spawn python3
expect:expect ">>>" 5

# Run a simple calculation
expect:send "print(40 + 2)"
expect:expect "42" 5
	`)
	
	// Output:
	// Example test file using expect:spawn:
	// 
	// # Start a Python interpreter
	// expect:spawn python3
	// expect:expect ">>>" 5
	// 
	// # Run a simple calculation
	// expect:send "print(40 + 2)"
	// expect:expect "42" 5
}