// Package expect provides scripttest commands for interacting with 
// interactive programs using the expect utility.
package expect

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"rsc.io/script"
)

// Commands returns a map of expect-related commands to add to a scripttest engine.
func Commands() map[string]script.Cmd {
	cmds := make(map[string]script.Cmd)
	
	// Register each expect command
	cmds["expect:spawn"] = spawnCmd()
	cmds["expect:send"] = sendCmd()
	cmds["expect:expect"] = expectCmd()
	cmds["expect:interact"] = interactCmd()
	cmds["expect:script"] = scriptCmd()
	
	return cmds
}

// spawnCmd creates a command to spawn a process to interact with via expect
func spawnCmd() script.Cmd {
	return script.Command(
		script.CmdUsage{
			Summary: "Start a new process to interact with",
			Args:    "program [args...]",
			Detail: []string{
				"expect:spawn starts a new interactive process that can be controlled with",
				"subsequent expect commands. The process remains running until explicitly",
				"closed or until the script ends.",
				"",
				"Example:",
				"  expect:spawn ssh localhost",
				"  expect:spawn python",
				"  expect:spawn telnet example.com 23",
			},
		},
		func(s *script.State, args ...string) (script.WaitFunc, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("expect:spawn requires at least a program name")
			}
			
			// Create a temporary expect script
			scriptFile, err := createTempFile(s, "spawn.exp", `
				#!/usr/bin/expect -f
				log_user 1
				spawn {{.Command}}
				interact
			`, map[string]interface{}{
				"Command": strings.Join(args, " "),
			})
			if err != nil {
				return nil, err
			}
			
			// Make it executable
			if err := os.Chmod(scriptFile, 0755); err != nil {
				return nil, fmt.Errorf("failed to make expect script executable: %v", err)
			}
			
			// Set environment variables to track the current expect session
			if err := s.Setenv("EXPECT_SCRIPT", scriptFile); err != nil {
				return nil, fmt.Errorf("failed to set environment variable: %v", err)
			}
			if err := s.Setenv("EXPECT_PID", ""); err != nil {
				return nil, fmt.Errorf("failed to set environment variable: %v", err)
			}
			
			// Run the expect script
			cmd := exec.Command(scriptFile)
			cmd.Stdin = os.Stdin
			
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			
			if err := cmd.Start(); err != nil {
				return nil, fmt.Errorf("failed to start expect script: %v", err)
			}
			
			// Store the PID
			if err := s.Setenv("EXPECT_PID", fmt.Sprintf("%d", cmd.Process.Pid)); err != nil {
				return nil, fmt.Errorf("failed to set environment variable: %v", err)
			}
			
			return func(s *script.State) (string, string, error) {
				err := cmd.Wait()
				return stdout.String(), stderr.String(), err
			}, nil
		},
	)
}

// sendCmd creates a command to send input to the spawned process
func sendCmd() script.Cmd {
	return script.Command(
		script.CmdUsage{
			Summary: "Send input to the spawned process",
			Args:    "input [no_newline]",
			Detail: []string{
				"expect:send sends the specified input to the currently running process.",
				"By default, a newline is appended to the input. If \"no_newline\" is specified",
				"as the second argument, no newline will be appended.",
				"",
				"Example:",
				"  expect:spawn python",
				"  expect:send \"print('hello')\"",
				"  expect:send \"exit()\" no_newline",
			},
		},
		func(s *script.State, args ...string) (script.WaitFunc, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("expect:send requires input to send")
			}
			
			scriptFile, found := s.LookupEnv("EXPECT_SCRIPT")
			if !found || scriptFile == "" {
				return nil, fmt.Errorf("no expect process is currently running (use expect:spawn first)")
			}
			
			input := args[0]
			appendNewline := true
			if len(args) > 1 && args[1] == "no_newline" {
				appendNewline = false
			}
			
			// Create a temporary expect script for sending input
			pidValue, _ := s.LookupEnv("EXPECT_PID")
			sendScript, err := createTempFile(s, "send.exp", `
				#!/usr/bin/expect -f
				set pid {{.Pid}}
				
				# Send the input to the process
				spawn -noecho /bin/sh -c "echo {{.Input}} | expect_sendto $pid {{.Newline}}"
				expect eof
			`, map[string]interface{}{
				"Pid":      pidValue,
				"Input":    escapeTcl(input),
				"Newline":  appendNewline,
			})
			if err != nil {
				return nil, err
			}
			
			// Make it executable
			if err := os.Chmod(sendScript, 0755); err != nil {
				return nil, fmt.Errorf("failed to make send script executable: %v", err)
			}
			
			// Run the send script
			var stdout, stderr bytes.Buffer
			cmd := exec.Command(sendScript)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			
			if err := cmd.Run(); err != nil {
				return nil, fmt.Errorf("failed to send input: %v, stderr: %s", err, stderr.String())
			}
			
			return func(s *script.State) (string, string, error) {
				return stdout.String(), stderr.String(), nil
			}, nil
		},
	)
}

// expectCmd creates a command to wait for a pattern in the output
func expectCmd() script.Cmd {
	return script.Command(
		script.CmdUsage{
			Summary: "Wait for a pattern in the output",
			Args:    "pattern [timeout]",
			Detail: []string{
				"expect:expect waits for the specified pattern to appear in the output of the",
				"spawned process. If a timeout (in seconds) is specified, the command will",
				"fail if the pattern is not found within that time.",
				"",
				"Example:",
				"  expect:spawn python",
				"  expect:expect \">>>\" 5",
				"  expect:send \"print('hello')\"",
				"  expect:expect \"hello\"",
			},
		},
		func(s *script.State, args ...string) (script.WaitFunc, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("expect:expect requires a pattern to wait for")
			}
			
			scriptFile, found := s.LookupEnv("EXPECT_SCRIPT")
			if !found || scriptFile == "" {
				return nil, fmt.Errorf("no expect process is currently running (use expect:spawn first)")
			}
			
			pattern := args[0]
			timeout := "30" // Default timeout of 30 seconds
			if len(args) > 1 {
				timeout = args[1]
			}
			
			// Create a temporary expect script for waiting for output
			pidValue, _ := s.LookupEnv("EXPECT_PID")
			expectScript, err := createTempFile(s, "expect_pattern.exp", `
				#!/usr/bin/expect -f
				set pid {{.Pid}}
				set timeout {{.Timeout}}
				
				# Connect to the existing process
				spawn -noecho /bin/sh -c "expect_console $pid"
				
				# Wait for the pattern
				expect {
					"{{.Pattern}}" {
						puts "Pattern matched"
						exit 0
					}
					timeout {
						puts "Timeout waiting for pattern"
						exit 1
					}
					eof {
						puts "Process ended before pattern was found"
						exit 1
					}
				}
			`, map[string]interface{}{
				"Pid":     pidValue,
				"Pattern": escapeTcl(pattern),
				"Timeout": timeout,
			})
			if err != nil {
				return nil, err
			}
			
			// Make it executable
			if err := os.Chmod(expectScript, 0755); err != nil {
				return nil, fmt.Errorf("failed to make expect script executable: %v", err)
			}
			
			// Run the expect script
			var stdout, stderr bytes.Buffer
			cmd := exec.Command(expectScript)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			
			if err := cmd.Run(); err != nil {
				return nil, fmt.Errorf("pattern not matched: %v, output: %s", err, stdout.String())
			}
			
			return func(s *script.State) (string, string, error) {
				return stdout.String(), stderr.String(), nil
			}, nil
		},
	)
}

// interactCmd creates a command to interact with the spawned process
func interactCmd() script.Cmd {
	return script.Command(
		script.CmdUsage{
			Summary: "Start interactive mode with the spawned process",
			Args:    "[escape_character]",
			Detail: []string{
				"expect:interact enters interactive mode with the spawned process. Input and output",
				"are passed directly between the user and the process. If an escape character is",
				"specified, typing that character will exit interactive mode.",
				"",
				"Example:",
				"  expect:spawn ssh localhost",
				"  expect:expect \"password:\"",
				"  expect:send \"mypassword\"",
				"  expect:interact \"^]\"  # Ctrl-]",
			},
		},
		func(s *script.State, args ...string) (script.WaitFunc, error) {
			scriptFile, found := s.LookupEnv("EXPECT_SCRIPT")
			if !found || scriptFile == "" {
				return nil, fmt.Errorf("no expect process is currently running (use expect:spawn first)")
			}
			
			escapeChar := "\\035" // Default escape character (Ctrl-])
			if len(args) > 0 {
				escapeChar = args[0]
			}
			
			// Create a temporary expect script for interaction
			pidValue, _ := s.LookupEnv("EXPECT_PID")
			interactScript, err := createTempFile(s, "interact.exp", `
				#!/usr/bin/expect -f
				set pid {{.Pid}}
				
				# Connect to the existing process
				spawn -noecho /bin/sh -c "expect_console $pid"
				
				# Enter interactive mode
				interact {
					{{.EscapeChar}} {
						puts "\nExiting interactive mode"
						return
					}
				}
			`, map[string]interface{}{
				"Pid":        pidValue,
				"EscapeChar": escapeChar,
			})
			if err != nil {
				return nil, err
			}
			
			// Make it executable
			if err := os.Chmod(interactScript, 0755); err != nil {
				return nil, fmt.Errorf("failed to make interact script executable: %v", err)
			}
			
			// Run the interact script
			var stdout, stderr bytes.Buffer
			cmd := exec.Command(interactScript)
			cmd.Stdin = os.Stdin
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			
			if err := cmd.Run(); err != nil {
				return nil, fmt.Errorf("interaction failed: %v", err)
			}
			
			return func(s *script.State) (string, string, error) {
				return stdout.String(), stderr.String(), nil
			}, nil
		},
	)
}

// scriptCmd creates a command to run a complete expect script
func scriptCmd() script.Cmd {
	return script.Command(
		script.CmdUsage{
			Summary: "Run a complete expect script",
			Args:    "script_content",
			Detail: []string{
				"expect:script runs a complete expect script. The script content should be",
				"provided as a single argument. This is useful for complex interactions that",
				"would be cumbersome to express with individual commands.",
				"",
				"Example:",
				"  expect:script '",
				"  spawn ssh localhost",
				"  expect \"password:\"",
				"  send \"mypassword\\r\"",
				"  expect \"$ \"",
				"  send \"ls -la\\r\"",
				"  expect \"$ \"",
				"  send \"exit\\r\"",
				"  expect eof",
				"  '",
			},
		},
		func(s *script.State, args ...string) (script.WaitFunc, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("expect:script requires script content")
			}
			
			scriptContent := args[0]
			
			// Create a temporary expect script
			scriptFile, err := createTempFile(s, "script.exp", `
				#!/usr/bin/expect -f
				log_user 1
				{{.Script}}
			`, map[string]interface{}{
				"Script": scriptContent,
			})
			if err != nil {
				return nil, err
			}
			
			// Make it executable
			if err := os.Chmod(scriptFile, 0755); err != nil {
				return nil, fmt.Errorf("failed to make expect script executable: %v", err)
			}
			
			// Run the expect script
			var stdout, stderr bytes.Buffer
			cmd := exec.Command(scriptFile)
			cmd.Stdin = os.Stdin
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			
			if err := cmd.Start(); err != nil {
				return nil, fmt.Errorf("failed to start expect script: %v", err)
			}
			
			return func(s *script.State) (string, string, error) {
				err := cmd.Wait()
				return stdout.String(), stderr.String(), err
			}, nil
		},
	)
}

// Helper functions

// createTempFile creates a temporary file with content from a template
func createTempFile(s *script.State, name, templateText string, data map[string]interface{}) (string, error) {
	// Create a temporary directory if needed
	tmpDir := filepath.Join(s.Getwd(), ".expect-tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}
	
	// Create a temporary file
	filePath := filepath.Join(tmpDir, name)
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer file.Close()
	
	// Parse and execute the template
	tmpl, err := template.New(name).Parse(templateText)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}
	
	if err := tmpl.Execute(file, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}
	
	return filePath, nil
}

// escapeTcl escapes special characters in TCL strings
func escapeTcl(s string) string {
	// Replace TCL special characters
	replacer := strings.NewReplacer(
		"[", "\\[",
		"]", "\\]",
		"$", "\\$",
		"\"", "\\\"",
		"\\", "\\\\",
	)
	return replacer.Replace(s)
}