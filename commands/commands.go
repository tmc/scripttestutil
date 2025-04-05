// Package commands provides reusable command sets for scripttestutil.
// Each subpackage offers a specialized set of commands for a particular domain.
package commands

import (
	"github.com/tmc/scripttestutil/commands/expect"
	"rsc.io/script"
)

// RegisterAll adds all available command sets to the provided command map.
// This is a convenience function for registering all command sets at once.
func RegisterAll(cmds map[string]script.Cmd) {
	// Register each command set
	RegisterExpect(cmds)
	
	// Add more command sets here as they become available
	// RegisterSSH(cmds)
	// RegisterDocker(cmds)
	// etc.
}

// RegisterExpect adds the expect command set to the provided command map.
// These commands provide integration with the expect utility for interacting
// with interactive programs.
func RegisterExpect(cmds map[string]script.Cmd) {
	// Get the expect commands
	expectCmds := expect.Commands()
	
	// Add them to the provided map
	for name, cmd := range expectCmds {
		cmds[name] = cmd
	}
}