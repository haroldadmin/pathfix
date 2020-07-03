package pathfix

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Fix appends the PATH value of the user's terminal shell to current process's PATH environment variable
// It returns an error if the "SHELL" environment variable can not be found, or points to a non executable file.
func Fix() error {
	if runtime.GOOS == "windows" {
		return nil
	}

	defaultShell := os.Getenv("SHELL")
	if defaultShell == "" {
		return errors.New("Failed to retrieve default shell: No SHELL environment variable found")
	}

	buf := &bytes.Buffer{}
	envCommand := exec.Command(defaultShell, "-ilc", "env")
	envCommand.Stdout = buf

	err := envCommand.Start()
	if err != nil {
		return fmt.Errorf("Error starting shell: %v", err)
	}

	// Without explicitly calling Wait, tests hang forever
	// https://github.com/golang/go/issues/24050
	err = envCommand.Wait()

	if err != nil {
		return fmt.Errorf("Failed to run shell for retrieving environment variables: %v", err)
	}

	for _, envVar := range strings.Split(buf.String(), "\n") {
		if strings.HasPrefix(envVar, "PATH=") {
			split := strings.Split(envVar, "=")

			if len(split) < 2 {
				// New PATH is empty so return early
				return nil
			}

			newPath := split[1]
			currentPath := os.Getenv("PATH")
			completePath := buildCompletePath(currentPath, newPath)

			os.Setenv("PATH", completePath)
			break
		}
	}

	return nil
}

func buildCompletePath(currentPath string, newPath string) string {
	if strings.TrimSpace(currentPath) == "" {
		return newPath
	}

	return currentPath + string(os.PathListSeparator) + newPath
}
