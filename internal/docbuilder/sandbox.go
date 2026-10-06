package docbuilder

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SandboxClient struct {
	InputDir  string
	OutputDir string
	Timeout   time.Duration
}

type runDirs struct {
	input  string
	output string
}

type RunResult struct {
	OutputPath string
	StdOut     string
	StdErr     string
}

const defaultTimeout = 30 * time.Second

var allowedFormat = map[string]bool{
	"docx": true,
	"pdf":  true,
	"xlsx": true,
	"md":   true,
}

func NewSandboxClient(inputDir, outputDir string) *SandboxClient {
	return &SandboxClient{
		InputDir:  inputDir,
		OutputDir: outputDir,
	}
}

func ParseFormat(format string) (string, error) {
	ext := strings.ToLower(strings.TrimSpace(format))
	if !allowedFormat[ext] {
		return "", fmt.Errorf("Unsupported output format: %q", format)
	}
	return ext, nil
}

func SafeFileBase(title string) string {
	var b strings.Builder
	prevDash := false
	for _, r := range strings.ToLower(title) {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			prevDash = false
		case !prevDash && b.Len() > 0:
			b.WriteByte('-')
			prevDash = true
		}
	}
	s := strings.Trim(b.String(), "-")
	if len(s) > 60 {
		s = strings.Trim(s[:60], "-")
	}
	if s == "" {
		s = "dokumen"
	}
	return s
}

func RenameOutput(oldPath, title, ext string) (string, error) {
	newPath := filepath.Join(filepath.Dir(oldPath), SafeFileBase(title)+"."+ext)
	if newPath == oldPath {
		return oldPath, nil
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		return "", err
	}
	return newPath, nil
}

func (s *SandboxClient) prepareRunDirs() (*runDirs, error) {
	runID := uuid.New().String()

	inputDir, err := filepath.Abs(filepath.Join(s.InputDir, runID))
	if err != nil {
		return nil, fmt.Errorf("Failed to resolve the absolute path: %w", err)
	}

	outputDir, err := filepath.Abs(filepath.Join(s.OutputDir, runID))
	if err != nil {
		return nil, fmt.Errorf("Failed to resolve the absolute path: %w", err)
	}

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		return nil, fmt.Errorf("Failed to build input directory: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("Failed to build output directory: %w", err)
	}

	return &runDirs{
		input:  inputDir,
		output: outputDir,
	}, nil
}

func writeScript(inputDir, code string) error {
	scriptPath := filepath.Join(inputDir, "script.py")
	if err := os.WriteFile(scriptPath, []byte(code), 0644); err != nil {
		return fmt.Errorf("Failed to build the script: %w", err)
	}
	return nil
}

func (s *SandboxClient) dockerExec(ctx context.Context, dirs *runDirs) (string, string, error) {
	timeout := s.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run",
		"--rm",
		"--network", "none",
		"--memory", "256m",
		"--cpus", "1",
		"--pids-limit", "64",
		"-v", dirs.input+":/workspace/input:ro",
		"-v", dirs.output+":/workspace/output",
		"docbuilder-sandbox",
		"python", "/workspace/input/script.py",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return stdout.String(), stderr.String(), fmt.Errorf("Execution timeout after %s", timeout)
	}
	if err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("Failed to run the script: %w", err)
	}

	return stdout.String(), stderr.String(), nil
}

func verifyOutput(outputDir, outputFile string) (string, error) {
	outputPath := filepath.Join(outputDir, outputFile)

	info, err := os.Stat(outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("Output file %q is not found after execution", outputFile)
		}
		return "", fmt.Errorf("Failed to check the output file: %w", err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("%q is a folder", outputFile)
	}

	return outputPath, nil
}

func (s *SandboxClient) RunCode(ctx context.Context, code, outputFilename string) (*RunResult, error) {
	dirs, err := s.prepareRunDirs()
	if err != nil {
		return nil, fmt.Errorf("prepareRunDirs: %w", err)
	}
	defer os.RemoveAll(dirs.input)

	if err := writeScript(dirs.input, code); err != nil {
		return nil, err
	}

	stdout, stderr, err := s.dockerExec(ctx, dirs)
	if err != nil {
		return nil, fmt.Errorf("Execution failed: %w (stdout: %s, stderr: %s)", err, stdout, stderr)
	}

	outputPath, err := verifyOutput(dirs.output, outputFilename)
	if err != nil {
		return nil, err
	}

	return &RunResult{
		StdOut:     stdout,
		StdErr:     stderr,
		OutputPath: outputPath,
	}, nil
}
