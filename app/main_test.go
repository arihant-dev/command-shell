package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// Build the package into a temporary binary and return its path and a cleanup func.
func buildBinary(t *testing.T, dir string) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "cmdbin-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(tmpDir, "appbin")
	cmd := exec.Command("go", "build", "-o", binPath)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("go build failed: %v\noutput:\n%s", err, string(out))
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return binPath, cleanup
}

// Determine directory of this test file (assumed package dir).
func testDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine caller info")
	}
	return filepath.Dir(file)
}

// Test that the package builds successfully.
func TestBuild(t *testing.T) {
	dir := testDir(t)
	bin, cleanup := buildBinary(t, dir)
	defer cleanup()
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("built binary not found at %s: %v", bin, err)
	}
}

// Test that running the built binary with "-h" returns (does not hang).
// This does not assert specific output, only that the process starts and exits within the timeout.
func TestBinaryHelpDoesNotHang(t *testing.T) {
	dir := testDir(t)
	bin, cleanup := buildBinary(t, dir)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, "-h")
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatalf("binary hung or timed out after 2s")
	}
	// Build succeeded and the binary returned (exit code may be non-zero depending on implementation).
	t.Logf("binary executed; output (truncated): %s", truncateOutput(out, 1024))
	if err != nil {
		t.Logf("binary exited with error: %v", err)
	}
}

// small helper to keep logs readable
func truncateOutput(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "...(truncated)"
}

func TestPWD(t *testing.T) {

	sh := &Shell{
		Path:   testDir(t),
		Reader: bufio.NewReader(os.Stdin),
	}

	tests := map[string]struct {
		description string
	}{
		"happy path - get current directory": {
			description: "should return current working directory",
		},
	}

	for name := range tests {
		t.Run(name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("failed to create pipe: %v", err)
			}
			os.Stdout = w

			sh.handlePwd()

			w.Close()
			defer func() { os.Stdout = old }()

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			if err != nil {
				t.Fatalf("failed to read from pipe: %v", err)
			}
			result := strings.TrimSpace(buf.String())

			// Verify we got a valid directory path
			if result == "" {
				t.Error("expected non-empty directory path")
			}
			if !filepath.IsAbs(result) {
				t.Errorf("expected absolute path, got %q", result)
			}
		})
	}
}

func TestCD(t *testing.T) {

	sh := &Shell{
		Path:   testDir(t),
		Reader: bufio.NewReader(os.Stdin),
	}

	tests := map[string]struct {
		description string
	}{
		"happy path - change directory": {
			description: "should change the current working directory",
		},
	}

	for name := range tests {
		t.Run(name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("failed to create pipe: %v", err)
			}
			os.Stdout = w

			// determine home dir portably
			homeDir, err := os.UserHomeDir()
			if err != nil {
				homeDir = os.Getenv("HOME")
			}

			// cross-platform non-existent path
			nonExistent := filepath.Join(os.TempDir(), "non-existent-dir-please-delete")

			testCases := map[string]struct {
				dir string
			}{
				"change to parent directory":       {dir: ".."},
				"change to current directory":      {dir: "."},
				"change to non-existent directory": {dir: nonExistent},
			}
			if homeDir != "" {
				testCases["change to home directory"] = struct{ dir string }{dir: homeDir}
			}

			for _, tt := range testCases {
				sh.handleCd([]string{tt.dir})
				sh.handlePwd()
			}

			w.Close()
			defer func() { os.Stdout = old }()

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			if err != nil {
				t.Fatalf("failed to read from pipe: %v", err)
			}

			// Take the last non-empty line as the resulting PWD
			lines := strings.Split(buf.String(), "\n")
			result := ""
			for i := len(lines) - 1; i >= 0; i-- {
				line := strings.TrimSpace(lines[i])
				if line != "" {
					result = line
					break
				}
			}

			// Verify we got a valid directory path
			if result == "" {
				t.Error("expected non-empty directory path")
			}
			if !filepath.IsAbs(result) {
				t.Errorf("expected absolute path, got %q", result)
			}
		})
	}
}
