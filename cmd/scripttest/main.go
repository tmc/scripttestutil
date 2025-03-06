package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	_ "rsc.io/script/scripttest" // not strictly necessary but nice for go odc tool
)

func usage() {
	// Extract the content of the /* ... */ comment in doc.go.
	_, after, _ := strings.Cut(doc, "/*")
	doc, _, _ := strings.Cut(after, "*/")
	io.WriteString(flag.CommandLine.Output(), doc)
	flag.PrintDefaults()
	os.Exit(2)
}

//go:embed doc.go
var doc string

var (
	verbose         bool
	stream          bool
	pattern         string
	useDocker       bool
	dockerImage     string
	autoGoToolchain bool

	flagDebug bool
)

func main() {
	log.SetPrefix("scripttest: ")
	log.SetFlags(0)

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&stream, "stream", false, "stream cgpt output")
	flag.BoolVar(&flagDebug, "debug", false, "debug cgpt output")
	flag.StringVar(&pattern, "p", "testdata/*.txt", "test file pattern")
	flag.BoolVar(&useDocker, "docker", false, "run tests in Docker container")
	flag.StringVar(&dockerImage, "docker-image", "", "Docker image to use (defaults to golang:latest)")
	flag.BoolVar(&autoGoToolchain, "auto-go", true, "automatically download Go toolchain if needed")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
	}

	cmd := flag.Arg(0)
	args := flag.Args()[1:]

	switch cmd {
	case "test", "run":
		if err := runTests(args); err != nil {
			log.Fatal(err)
		}
	case "scaffold":
		if err := runScaffold(args); err != nil {
			log.Fatal(err)
		}
	case "infer":
		if err := runInfer(args); err != nil {
			log.Fatal(err)
		}
	case "help":
		if err := runHelp(args); err != nil {
			log.Fatal(err)
		}
	case "playback":
		if err := runPlayback(args); err != nil {
			log.Fatal(err)
		}
	case "record":
		if err := runRecord(args); err != nil {
			log.Fatal(err)
		}
	case "play-cast":
		if err := runPlayCast(args); err != nil {
			log.Fatal(err)
		}
	case "convert-cast":
		if err := runConvertCast(args); err != nil {
			log.Fatal(err)
		}
	default:
		usage()
	}
}

func scaffold(dir string) error {
	if verbose {
		log.Printf("scaffolding in directory: %s", dir)
	}

	info, err := loadOrInferCommandInfo(dir)
	if err != nil {
		return fmt.Errorf("failed to load or infer command info: %v", err)
	}

	if verbose {
		log.Printf("command info: %s", info)
	}

	prompt, err := generateScaffoldPrompt(info)
	if err != nil {
		return fmt.Errorf("failed to load prompt cgpt: %v", err)
	}
	resp, err := queryCgpt(prompt, "")
	if err != nil {
		return fmt.Errorf("failed to query cgpt: %v", err)
	}

	if verbose {
		log.Printf("cgpt response: %s", resp)
	}

	return applyScaffold(dir, resp)
}

func infer(dir string) error {
	if verbose {
		log.Printf("inferring command info in directory: %s", dir)
	}

	info, err := inferCommandInfo(dir)
	if err != nil {
		return fmt.Errorf("failed to infer command info: %v", err)
	}

	file := filepath.Join(dir, ".scripttest_info")
	if err := os.WriteFile(file, []byte(info), 0644); err != nil {
		return fmt.Errorf("failed to write command info: %v", err)
	}

	if verbose {
		log.Printf("command info written to: %s", file)
	}

	return nil
}

func loadOrInferCommandInfo(dir string) (string, error) {
	file := filepath.Join(dir, ".scripttest_info")
	info, err := os.ReadFile(file)
	if err == nil {
		if verbose {
			log.Printf("loaded existing command info from: %s", file)
		}
		return string(info), nil
	}
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to read command info: %v", err)
	}
	if verbose {
		log.Printf("inferring command info")
	}
	return inferCommandInfo(dir)
}

func inferCommandInfo(dir string) (string, error) {
	content, err := getCodebaseContent(dir)
	if err != nil {
		return "", fmt.Errorf("failed to get codebase content: %v", err)
	}

	prompt := fmt.Sprintf("Analyze this codebase and identify key binary entrypoints and commnds:\n\n%s\n\n", content)
	prompt += `output a json representation matching this datatype:
type Commands = CommandInfo[];
type CommandInfo = {
  name: string;    // command name
  summary: string; // usage summary
  args: string;    // argument pattern
}`
	res, err := queryCgpt(prompt, "```json\n[")
	if err != nil {
		return "", fmt.Errorf("failed to run cgpt: %w", err)
	}
	return res, nil
}

func getCodebaseContent(dir string) (string, error) {
	// Check if code-to-gpt.sh exists in PATH
	scriptPath, err := exec.LookPath("code-to-gpt.sh")
	if err != nil {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = filepath.Join(os.Getenv("HOME"), "go")
			if runtime.GOOS == "windows" {
				gopath = filepath.Join(os.Getenv("USERPROFILE"), "go")
			}
		}

		binDir := filepath.Join(gopath, "bin")
		scriptPath = filepath.Join(binDir, "code-to-gpt.sh")

		if err := downloadScript(scriptPath); err != nil {
			return "", fmt.Errorf("failed to download code-to-gpt.sh: %v", err)
		}

		// Make executable
		if err := os.Chmod(scriptPath, 0755); err != nil {
			return "", fmt.Errorf("failed to make script executable: %v", err)
		}
	}

	args := []string{"--", ":!.scripttest_history*"}
	cmd := exec.Command(scriptPath, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run code-to-gpt: %v", err)
	}
	return string(output), nil
}

func downloadScript(destPath string) error {
	// Create bin directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// Download script from a trusted source
	resp, err := http.Get("https://raw.githubusercontent.com/tmc/misc/master/code-to-gpt/code-to-gpt.sh")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func queryCgpt(prompt, prefill string) (string, error) {
	// TODO: add history file support
	args := []string{
		// "-I", historyFile,
		// "-O", historyFile,
	}

	if flagDebug {
		args = append(args, "--debug")
	}

	if prefill != "" {
		args = append(args, "--prefill", prefill)
	} else {
		args = append(args, "--prefill", "```json\n{")
	}

	cmd := exec.Command("cgpt", args...)

	var output strings.Builder
	var stderr io.Writer = new(strings.Builder)
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Stdout = &output
	if stream {
		stderr = os.Stderr
		cmd.Stdout = io.MultiWriter(&output, stderr)
	}
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("cgpt query failed: %v", err)
	}
	return extractJSON(output.String()), nil
}

func runTest(pattern string) error {
	if verbose {
		log.Printf("running tests matching pattern: %s", pattern)
	}

	// Get clean work directory
	dir, err := getWorkDir()
	if err != nil {
		return fmt.Errorf("failed to get work directory: %v", err)
	}

	if verbose {
		log.Printf("using work directory: %s", dir)
	}

	// Create testdata directory
	testdata := filepath.Join(dir, "testdata")
	if err := os.MkdirAll(testdata, 0755); err != nil {
		return fmt.Errorf("failed to create testdata directory: %v", err)
	}

	// Set up test files in work directory
	if err := setupTestDir(dir); err != nil {
		return fmt.Errorf("failed to setup test directory: %v", err)
	}

	// Find matching test files
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %v", err)
	}
	if len(matches) == 0 {
		return fmt.Errorf("no files match pattern: %s", pattern)
	}

	// Create symlinks in testdata directory
	for _, file := range matches {
		abs, err := filepath.Abs(file)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %v", file, err)
		}
		dst := filepath.Join(testdata, filepath.Base(file))
		if err := os.Symlink(abs, dst); err != nil {
			return fmt.Errorf("failed to link test file %s: %v", file, err)
		}
	}

	// link in .scripttest_info if it exists:
	scriptTestInfo := ".scripttest_info"
	if _, err := os.Stat(scriptTestInfo); err == nil {
		abs, err := filepath.Abs(scriptTestInfo)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for .scripttest_info: %v", err)
		}
		dst := filepath.Join(dir, ".scripttest_info")
		if err := os.Symlink(abs, dst); err != nil {
			return fmt.Errorf("failed to link .scripttest_info: %v", err)
		}
	}

	// Initialize go modules
	if err := initModules(dir); err != nil {
		return fmt.Errorf("failed to initialize modules: %v", err)
	}

	buildID := getBuildID()
	if verbose {
		log.Printf("build ID: %s", buildID)
	}

	// Run go test in the directory
	args := []string{"test"}
	if verbose {
		args = append(args, "-v")
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %v", err)
	}

	return nil
}

func runTestInDocker(pattern string) error {
	if verbose {
		log.Printf("running tests in Docker with pattern: %s", pattern)
	}

	// Get clean work directory
	dir, err := getWorkDir()
	if err != nil {
		return fmt.Errorf("failed to get work directory: %v", err)
	}

	if verbose {
		log.Printf("using work directory: %s", dir)
	}

	// Check for Dockerfile in test files
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %v", err)
	}
	if len(matches) == 0 {
		return fmt.Errorf("no files match pattern: %s", pattern)
	}

	// Look for Dockerfile content in test files
	var dockerfileContent string
	for _, file := range matches {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read test file %s: %v", file, err)
		}
		content := string(data)

		// Look for Dockerfile marker in test file
		if idx := strings.Index(content, "-- Dockerfile --"); idx != -1 {
			// Extract Dockerfile content
			content = content[idx+len("-- Dockerfile --"):]
			if end := strings.Index(content, "\n--"); end != -1 {
				dockerfileContent = strings.TrimSpace(content[:end])
			} else {
				dockerfileContent = strings.TrimSpace(content)
			}
			break
		}
	}

	// Create docker-bake.hcl
	dockerBakeFile := filepath.Join(dir, "docker-bake.hcl")
	if dockerfileContent != "" {
		if verbose {
			log.Printf("using Dockerfile from test file")
		}
		bakeContent := fmt.Sprintf(`
			group "default" {
				targets = ["scripttest-runner"]
			}

			target "scripttest-runner" {
				context = "."
				dockerfile = "Dockerfile"
			}
		`)
		if err := os.WriteFile(dockerBakeFile, []byte(bakeContent), 0644); err != nil {
			return fmt.Errorf("failed to write docker-bake.hcl: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte(dockerfileContent), 0644); err != nil {
			return fmt.Errorf("failed to write Dockerfile: %v", err)
		}
	} else {
		// Use default Dockerfile if none found in test files
		if verbose {
			log.Printf("using default Dockerfile")
		}
		image := dockerImage
		if image == "" {
			image = "golang:latest"
		}
		content := fmt.Sprintf(`FROM %s
WORKDIR /app
COPY . .
RUN go mod download
CMD ["go", "test", "-v"]`, image)
		bakeContent := fmt.Sprintf(`
			group "default" {
				targets = ["scripttest-runner"]
			}

			target "scripttest-runner" {
				context = "."
				dockerfile = "Dockerfile"
			}
		`)
		if err := os.WriteFile(dockerBakeFile, []byte(bakeContent), 0644); err != nil {
			return fmt.Errorf("failed to write docker-bake.hcl: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create Dockerfile: %v", err)
		}
	}

	// Build Docker image using buildx bake
	buildCmd := exec.Command("docker", "buildx", "bake")
	buildCmd.Dir = dir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build Docker image: %v", err)
	}

	// Run tests in container
	args := []string{"run", "--rm"}

	// Mount the workspace
	args = append(args, "-v", fmt.Sprintf("%s:/app", dir))

	// Pass through environment variables
	args = append(args, "-e", "SCRIPTTEST_PATTERN="+pattern)
	if verbose {
		args = append(args, "-e", "VERBOSE=1")
	}
	if os.Getenv("UPDATE_SNAPSHOTS") == "1" {
		args = append(args, "-e", "UPDATE_SNAPSHOTS=1")
	}

	args = append(args, "scripttest-runner")

	runCmd := exec.Command("docker", args...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	if err := runCmd.Run(); err != nil {
		return fmt.Errorf("tests failed in Docker: %v", err)
	}

	return nil
}

func applyScaffold(dir string, resp string) error {
	var files map[string]string
	if err := json.Unmarshal([]byte(resp), &files); err != nil {
		// TODO
	}

	for path, content := range files {
		fullPath := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %v", path, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", path, err)
		}
		log.Printf("created %s", path)
	}
	return nil
}

func extractJSON(output string) string {
	// Try to find JSON between markdown code fences
	prefix := "```json"
	suffix := "```"
	start := strings.Index(output, prefix)
	if start == -1 {
		// Try alternate code fence
		prefix = "~~~json"
		start = strings.Index(output, prefix)
	}
	if start != -1 {
		start += len(prefix)
		// Find closing fence after the start position
		end := strings.Index(output[start:], suffix)
		if end != -1 {
			// Trim whitespace and validate JSON
			jsonStr := strings.TrimSpace(output[start : start+end])
			if json.Valid([]byte(jsonStr)) {
				return jsonStr
			}
		}
	}
	// If no valid JSON found between fences, try to find and validate any JSON in the string
	if json.Valid([]byte(output)) {
		return output
	}
	return "" // Return empty if no valid JSON found
}

// getCacheDir returns the scripttest cache directory, creating it if needed
func getCacheDir() (string, error) {
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %v", err)
		}
		cacheDir = filepath.Join(home, ".cache")
	}

	scripttestCache := filepath.Join(cacheDir, "scripttest")
	if err := os.MkdirAll(scripttestCache, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %v", err)
	}

	workDir := filepath.Join(scripttestCache, "workdir")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create work directory: %v", err)
	}

	return workDir, nil
}

// getWorkDir returns a clean temporary directory for the test run
func getWorkDir() (string, error) {
	// Create temp directory that will be automatically cleaned up
	tempDir, err := os.MkdirTemp("", "scripttest-workdir-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	return tempDir, nil
}

func initModules(dir string) error {
	// Check if Go is installed and install it if needed
	if err := ensureGoToolchain(); err != nil {
		return fmt.Errorf("failed to ensure Go toolchain: %v", err)
	}

	// Run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	cmd.Env = os.Environ() // Ensure we pass through GO111MODULE, GOPATH etc.

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %v\n%s", err, stderr.String())
	}

	return nil
}

func runPlayback(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("playback requires snapshot file argument")
	}
	snapshotPath := args[0]

	// Use scriptreplay for playback
	cmd := exec.Command("scriptreplay", snapshotPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runTests(args []string) error {
	// If pattern provided as argument, override flag
	if len(args) > 0 {
		pattern = args[0]
	}
	if useDocker {
		return runTestInDocker(pattern)
	}
	return runTest(pattern)
}

func runScaffold(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	return scaffold(dir)
}

func runInfer(args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	return infer(dir)
}

func runHelp(args []string) error {
	// Print usage first
	_, after, _ := strings.Cut(doc, "/*")
	doc, _, _ := strings.Cut(after, "*/")
	io.WriteString(os.Stdout, doc)
	return nil
}


func runRecord(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("record requires test file and output file arguments")
	}
	testFile := args[0]
	outputFile := args[1]
	
	return recordAsciicast(testFile, outputFile)
}

func runPlayCast(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("play-cast requires asciicast file argument")
	}
	asciicastFile := args[0]
	
	return playAsciicast(asciicastFile)
}

func runConvertCast(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("convert-cast requires snapshot file and output file arguments")
	}
	snapshotFile := args[0]
	outputFile := args[1]
	
	return convertSnapshotToAsciicast(snapshotFile, outputFile)
}
