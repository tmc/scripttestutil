# Expect Commands for Scripttest

This package provides a set of commands for interacting with interactive programs using the `expect` utility within scripttestutil tests.

## Commands

### expect:spawn

Starts a new interactive program that can be controlled with subsequent expect commands.

```
expect:spawn program [args...]
```

Example:
```
expect:spawn python3
expect:spawn ssh user@host
expect:spawn telnet example.com 23
```

### expect:expect

Waits for a specific pattern to appear in the output of the spawned program.

```
expect:expect pattern [timeout]
```

Example:
```
expect:expect ">>>" 5
expect:expect "Password:" 10
expect:expect "Connection established"
```

### expect:send

Sends input to the spawned program. By default, a newline is appended.

```
expect:send input [no_newline]
```

Example:
```
expect:send "print('hello')"
expect:send "mypassword" no_newline
```

### expect:interact

Enters interactive mode with the spawned program, allowing direct user interaction.

```
expect:interact [escape_character]
```

Example:
```
expect:interact "^]"  # Use Ctrl-] to exit
```

### expect:script

Runs a complete expect script for more complex interactions.

```
expect:script 'script_content'
```

Example:
```
expect:script '
spawn ssh user@host
expect "Password:"
send "mypassword\r"
expect "$ "
send "ls -la\r"
expect "$ "
send "exit\r"
expect eof
'
```

## Integration with Testscript

To use these commands in your testscript tests, you need to register them with your testscript runner. Here's how to do it:

```go
import (
    "testing"
    
    "github.com/tmc/scripttestutil/commands"
    "github.com/tmc/scripttestutil/testscript"
)

func TestMyExpectTests(t *testing.T) {
    opts := testscript.DefaultOptions()
    
    // Register expect commands
    opts.SetupHook = func(cmds map[string]script.Cmd) {
        commands.RegisterExpect(cmds)
    }
    
    // Run your tests
    testscript.RunDir(t, "testdata/expect_tests", opts)
}
```

## Example Test

Here's a simple test that uses the expect commands to interact with Python:

```
# Start a Python interpreter
expect:spawn python3
expect:expect ">>>" 5

# Run a simple calculation
expect:send "print(40 + 2)"
expect:expect "42" 5
expect:expect ">>>" 5

# Exit the Python interpreter
expect:send "exit()"

# Verify test completion
echo "Test completed successfully"
stdout "Test completed successfully"
```