// Package bridge provides a clean connection between Go's testing framework and scripttest.
// It allows scripttest tests to be run as part of standard Go test suites.
package bridge

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"rsc.io/script"
	"rsc.io/script/scripttest"
)

// TestOptions defines configuration options for running scripttest tests.
type TestOptions struct {
	// Pattern is the glob pattern to match test files (default: "testdata/*.txt")
	Pattern string

	// UseDocker determines whether to run tests in Docker containers
	UseDocker bool

	// DockerImage specifies which Docker image to use when UseDocker is true (default: golang:latest)
	DockerImage string

	// UpdateSnapshots controls whether snapshots should be updated
	UpdateSnapshots bool

	// Verbose enables verbose output
	Verbose bool

	// EnvVars defines additional environment variables to pass to tests
	EnvVars map[string]string

	// SnapshotDir specifies the directory for snapshots (default: testdata/__snapshots__)
	SnapshotDir string
}

// DefaultOptions returns the default test options.
func DefaultOptions() TestOptions {
	return TestOptions{
		Pattern:         "testdata/*.txt",
		UseDocker:       false,
		DockerImage:     "golang:latest",
		UpdateSnapshots: false,
		Verbose:         false,
		EnvVars:         make(map[string]string),
		SnapshotDir:     "testdata/__snapshots__",
	}
}

// Runner manages the execution of scripttest tests using Go's testing package.
type Runner struct {
	Options TestOptions
}

// NewRunner creates a new scripttest runner with the provided options.
func NewRunner(opts TestOptions) *Runner {
	return &Runner{Options: opts}
}

// Run executes scripttest tests matched by the pattern.
func (r *Runner) Run(t *testing.T) {
	// Get the original working directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Create a temporary directory for this test run
	tempDir, err := os.MkdirTemp("", "scripttest-bridge-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Find all matching test files
	matches, err := filepath.Glob(r.Options.Pattern)
	if err != nil {
		t.Fatalf("Invalid pattern %q: %v", r.Options.Pattern, err)
	}
	if len(matches) == 0 {
		t.Fatalf("No files match pattern %q", r.Options.Pattern)
	}

	// Create snapshot directory if needed
	if r.Options.UpdateSnapshots {
		snapshotDir := r.Options.SnapshotDir
		if !filepath.IsAbs(snapshotDir) {
			snapshotDir = filepath.Join(origDir, snapshotDir)
		}
		if err := os.MkdirAll(snapshotDir, 0755); err != nil {
			t.Fatalf("Failed to create snapshot directory: %v", err)
		}
	}

	// Process each test file
	for _, testFile := range matches {
		testName := filepath.Base(testFile)
		t.Run(testName, func(t *testing.T) {
			// Create test directory
			testDir := filepath.Join(tempDir, testName)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			// Run the test
			if err := r.runTest(t, testFile, testDir); err != nil {
				t.Fatalf("Test failed: %v", err)
			}
		})
	}
}

// RunTest runs a single scripttest test file.
func (r *Runner) RunTest(t *testing.T, testFile string) {
	// Create a temporary directory for this test
	tempDir, err := os.MkdirTemp("", "scripttest-single-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Run the test
	if err := r.runTest(t, testFile, tempDir); err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

// runTest handles the actual execution of a scripttest test.
func (r *Runner) runTest(t *testing.T, testFile, testDir string) error {
	// Setup environment
	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"TMPDIR=" + os.TempDir(),
	}

	// Add GOPATH if present
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		env = append(env, "GOPATH="+gopath)
	}

	// Add user-defined environment variables
	for k, v := range r.Options.EnvVars {
		env = append(env, k+"="+v)
	}

	// Set update snapshots environment variable if needed
	if r.Options.UpdateSnapshots {
		env = append(env, "UPDATE_SNAPSHOTS=1")
	}

	// Start with default script commands
	cmds := scripttest.DefaultCmds()

	// Create a snapshot handler
	setupSnapshotCommand(cmds, r.Options.SnapshotDir)

	// Start with default conditions
	conds := scripttest.DefaultConds()

	// Add platform-specific conditions
	setupPlatformConditions(conds)

	// Create engine
	engine := &script.Engine{
		Cmds:  cmds,
		Conds: conds,
		Quiet: !r.Options.Verbose && !testing.Verbose(),
	}

	// Configure test context
	ctx := context.Background()
	if t.Deadline() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, t.Deadline())
		defer cancel()
	}

	// Create testdata directory and copy the test file
	testdataDir := filepath.Join(testDir, "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		return fmt.Errorf("failed to create testdata directory: %v", err)
	}

	// Copy the test file to testdata directory
	testContent, err := os.ReadFile(testFile)
	if err != nil {
		return fmt.Errorf("failed to read test file: %v", err)
	}

	destFile := filepath.Join(testdataDir, filepath.Base(testFile))
	if err := os.WriteFile(destFile, testContent, 0644); err != nil {
		return fmt.Errorf("failed to write test file: %v", err)
	}

	// If using Docker, add the Docker-specific pattern
	testPattern := "testdata/" + filepath.Base(testFile)

	// Capture stdout and stderr for the test output
	var stdout, stderr bytes.Buffer
	oldStdout, oldStderr := os.Stdout, os.Stderr
	if !r.Options.Verbose && !testing.Verbose() {
		os.Stdout, os.Stderr = &stdout, &stderr
		defer func() { os.Stdout, os.Stderr = oldStdout, oldStderr }()
	}

	// Run the test
	if r.Options.UseDocker {
		// TODO: Implement Docker support by calling the appropriate functions
		return fmt.Errorf("Docker support not yet implemented in bridge")
	} else {
		// Change to the test directory
		oldDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %v", err)
		}
		if err := os.Chdir(testDir); err != nil {
			return fmt.Errorf("failed to change to test directory: %v", err)
		}
		defer os.Chdir(oldDir)

		// Run the test
		scripttest.Test(t, ctx, engine, env, testPattern)
	}

	return nil
}

// setupSnapshotCommand adds the snapshot command to the engine.
func setupSnapshotCommand(cmds map[string]script.Cmd, snapshotDir string) {
	cmds["snapshot"] = script.Command(
		script.CmdUsage{
			Summary: "Record command output",
			Args:    "[filename]",
		},
		func(s *script.State, args ...string) (script.WaitFunc, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("snapshot command requires a filename argument")
			}
			
			// Create the snapshots directory if it doesn't exist
			if err := os.MkdirAll(filepath.Dir(args[0]), 0755); err != nil {
				return nil, fmt.Errorf("failed to create snapshot directory: %v", err)
			}

			// Handle snapshot command
			// This is a simplified implementation
			// You would need to implement proper snapshot handling here
			return func(s *script.State) (string, string, error) {
				return "", "", nil
			}, nil
		},
	)
}

// setupPlatformConditions adds platform-specific conditions to the engine.
func setupPlatformConditions(conds map[string]script.Cond) {
	// Unix condition
	conds["unix"] = script.OnceCondition("unix system", func() (bool, error) {
		return os.Getenv("GOOS") != "windows", nil
	})

	// Windows condition
	conds["windows"] = script.OnceCondition("windows system", func() (bool, error) {
		return os.Getenv("GOOS") == "windows", nil
	})

	// macOS condition
	conds["darwin"] = script.OnceCondition("darwin system", func() (bool, error) {
		return os.Getenv("GOOS") == "darwin", nil
	})

	// Linux condition
	conds["linux"] = script.OnceCondition("linux system", func() (bool, error) {
		return os.Getenv("GOOS") == "linux", nil
	})
}

// RunFile is a convenience function to run a single scripttest file.
func RunFile(t *testing.T, file string, opts TestOptions) {
	runner := NewRunner(opts)
	runner.RunTest(t, file)
}

// RunPattern is a convenience function to run all scripttest files matching a pattern.
func RunPattern(t *testing.T, pattern string, opts TestOptions) {
	opts.Pattern = pattern
	runner := NewRunner(opts)
	runner.Run(t)
}

// RunDir is a convenience function to run all scripttest files in a directory.
func RunDir(t *testing.T, dir string, opts TestOptions) {
	pattern := filepath.Join(dir, "*.txt")
	RunPattern(t, pattern, opts)
}