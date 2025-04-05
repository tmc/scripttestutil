# Creating Custom Command Sets for Scripttest

This guide explains how to create reusable command sets for scripttest. Command sets allow you to extend scripttest with specialized functionality for different domains, such as interacting with databases, testing APIs, or controlling interactive programs.

## Overview

A command set is a collection of related commands that extend the scripttest engine. Each command can execute specific functionality and interact with the test environment. Commands are registered with the scripttest engine and can be invoked from test files.

## Creating a Command Set Package

### 1. Create a package structure

Create a new package under the `github.com/tmc/scripttestutil/commands/` directory. For example, to create a command set for PostgreSQL:

```
commands/
  └── postgres/
      ├── postgres.go
      ├── postgres_test.go
      ├── README.md
      └── testdata/
          └── postgres_basic_test.txt
```

### 2. Implement the Commands function

In your command package, implement a `Commands` function that returns a map of script commands:

```go
// Package postgres provides scripttest commands for working with PostgreSQL databases.
package postgres

import (
    "fmt"
    "rsc.io/script"
)

// Commands returns a map of PostgreSQL-related commands for scripttest.
func Commands() map[string]script.Cmd {
    cmds := make(map[string]script.Cmd)
    
    // Register each command
    cmds["postgres:init"] = initCmd()
    cmds["postgres:query"] = queryCmd()
    // Add more commands...
    
    return cmds
}

// initCmd creates a command to initialize a PostgreSQL database
func initCmd() script.Cmd {
    return script.Command(
        script.CmdUsage{
            Summary: "Initialize a PostgreSQL database",
            Args:    "dbname",
            Long:    "postgres:init creates a PostgreSQL database for testing",
        },
        func(s *script.State, args ...string) (script.WaitFunc, error) {
            // Implementation goes here
            return func(s *script.State) (string, string, error) {
                return "Database initialized", "", nil
            }, nil
        },
    )
}

// queryCmd creates a command to run a SQL query
func queryCmd() script.Cmd {
    // Implementation goes here
}
```

### 3. Register the command set

Add your command set to the main `commands` package by editing `commands/commands.go`:

```go
package commands

import (
    "github.com/tmc/scripttestutil/commands/expect"
    "github.com/tmc/scripttestutil/commands/postgres" // Add your new package
    "rsc.io/script"
)

// RegisterAll adds all available command sets.
func RegisterAll(cmds map[string]script.Cmd) {
    RegisterExpect(cmds)
    RegisterPostgres(cmds) // Add registration function
}

// RegisterPostgres adds PostgreSQL commands to the command map.
func RegisterPostgres(cmds map[string]script.Cmd) {
    pgCmds := postgres.Commands()
    for name, cmd := range pgCmds {
        cmds[name] = cmd
    }
}
```

## Best Practices

### Naming Conventions

- Use a namespace prefix for your commands (e.g., `postgres:query`) to avoid conflicts
- Use descriptive command names that clearly indicate their purpose
- Keep the namespace consistent across related commands

### Documentation

- Document each command thoroughly, including:
  - Summary of what the command does
  - Required and optional arguments
  - Examples of usage
  - Any side effects or requirements
- Create a README.md for your command set with complete documentation
- Include example tests in your testdata directory

### Command Implementation

- Make commands robust against user errors with clear error messages
- Use the script.State to manage state between commands
- Environment variables are a good way to store session information
- Clean up any resources created by your commands
- Implement proper error handling

### Testing

- Write tests for your command set
- Include example test files that demonstrate usage
- Test both success and failure cases
- Consider edge cases and boundaries

## Using Command Sets

Command sets can be used in scripttest tests by registering them with the testscript options:

```go
import (
    "testing"
    
    "github.com/tmc/scripttestutil/commands"
    "github.com/tmc/scripttestutil/testscript"
    "rsc.io/script"
)

func TestWithCommands(t *testing.T) {
    opts := testscript.DefaultOptions()
    
    // Register specific command sets
    opts.SetupHook = func(cmds map[string]script.Cmd) {
        commands.RegisterPostgres(cmds)
    }
    
    // Or register all command sets
    // opts.SetupHook = func(cmds map[string]script.Cmd) {
    //     commands.RegisterAll(cmds)
    // }
    
    testscript.RunDir(t, "testdata", opts)
}
```

Then in your test files:

```
# Initialize a PostgreSQL database
postgres:init testdb

# Run a query
postgres:query "CREATE TABLE users (id SERIAL, name TEXT)"
postgres:query "INSERT INTO users (name) VALUES ('Alice')"

# Verify results
postgres:query "SELECT name FROM users"
stdout 'Alice'
```

## Examples of Command Sets

The project includes several command sets as examples:

- **Expect Commands**: For interacting with interactive programs
- **PostgreSQL Commands**: For testing with PostgreSQL databases
- **Docker Commands**: For managing Docker containers during tests
- **SSH Commands**: For remote testing via SSH

Study these examples for guidance on implementing your own command sets.

## Contributing

If you create a generally useful command set, consider contributing it back to the scripttestutil project for others to use.