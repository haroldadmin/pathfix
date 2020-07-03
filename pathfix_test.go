package pathfix_test

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/haroldadmin/pathfix"
)

func BenchmarkPathFix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pathfix.Fix()
	}
}

func TestPathFix(t *testing.T) {

	ogPath := os.Getenv("PATH")
	ogShell := os.Getenv("SHELL")

	// Workaround for missing shell on github workflow runners
	if ogShell == "" {
		ogShell = "/bin/bash"
		os.Setenv("SHELL", ogShell)
	}

	t.Run("should fix current process's PATH", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Test not applicable to Windows")
		}

		defer resetEnv(t, ogPath, ogShell)

		t.Logf("Current PATH: %s\n\n", os.Getenv("PATH"))

		os.Unsetenv("PATH")
		if path := os.Getenv("PATH"); path != "" {
			t.Error("Failed to unset PATH for this process")
		}

		err := pathfix.Fix()
		if err != nil {
			t.Errorf("Expected no errors, got %v", err)
		}

		if path := os.Getenv("PATH"); path == "" {
			t.Errorf("Fixing path did not work. Current path: %s", path)
		}

	})

	t.Run("should return an error if SHELL env var is not set", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Test not applicable to Windows")
		}

		defer resetEnv(t, ogPath, ogShell)

		os.Unsetenv("SHELL")
		if shell := os.Getenv("SHELL"); shell != "" {
			t.Error("Failed to unset SHELL env var")
		}

		err := pathfix.Fix()
		if err == nil {
			t.Error("Expected an error, got none")
		}

	})

	t.Run("should not attempt to fix PATH if running on windows", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skipf("Not running on windows, this test is not applicable")
		}

		defer resetEnv(t, ogPath, ogShell)

		os.Unsetenv("PATH")
		if path := os.Getenv("PATH"); path != "" {
			t.Error("Failed to unset PATH")
		}

		if err := pathfix.Fix(); err != nil {
			t.Errorf("Expected no errors, got: %v", err)
		}

		if path := os.Getenv("PATH"); path != "" {
			t.Error("Path fix was performed, expected it to be skipped")
		}

	})

	t.Run("should return an error if shell process fails", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Test not applicable to Windows")
		}

		defer resetEnv(t, ogPath, ogShell)

		// Supply an invalid value for SHELL executable so that invoking it fails
		os.Setenv("SHELL", "/blahblahyadayada")

		if err := pathfix.Fix(); err == nil {
			t.Error("Expected an error, got none")
		}
	})

	t.Run("should append to the old PATH if it is not empty", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Test not applicable to Windows")
		}

		defer resetEnv(t, ogPath, ogShell)

		currentPath := "~/blah"
		t.Logf("Starting with PATH:\n%s\n\n", currentPath)

		os.Setenv("PATH", currentPath)

		if err := pathfix.Fix(); err != nil {
			t.Errorf("Expected no errors, got: %v", err)
		}

		path := os.Getenv("PATH")

		if !strings.HasPrefix(path, "~/blah") {
			t.Error("New PATH was not appended to old path: Could not find old PATH at the start")
		}
	})
}

func resetEnv(t *testing.T, ogPath, ogShell string) {
	t.Helper()
	os.Setenv("PATH", ogPath)
	os.Setenv("SHELL", ogShell)
}
