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
func Fix() error {
	if runtime.GOOS == "windows" {
		return nil
	}

	existingPath := os.Getenv("PATH")

	env, err := getEnv()
	if err != nil {
		return err
	}

	newPath := extractPath(env)
	if newPath == "" {
		// New PATH is empty, there's nothing to fix
		return nil
	}

	newPathset := createPathset(newPath)

	pathBuilder := strings.Builder{}
	pathBuilder.WriteString(existingPath)

	for p := range newPathset {
		if !strings.Contains(existingPath, p) {
			pathBuilder.WriteRune(os.PathListSeparator)
			pathBuilder.WriteString(p)
		}
	}

	combinedPath := pathBuilder.String()

	err = os.Setenv("PATH", combinedPath)
	if err != nil {
		return err
	}

	return nil
}

// getEnv extracts all the environment variables from the user's default shell by launching a
// interactive, login session and running the `env` command. It returns an error if the SHELL
// environment variable is not set, or if there is an error in invoking the shell process
func getEnv() (string, error) {
	defaultShell := os.Getenv("SHELL")

	if defaultShell == "" {
		return "", errors.New("Failed to retrieve default shell: No SHELL environment variable found")
	}

	buf := &bytes.Buffer{}
	envCommand := exec.Command(defaultShell, "-ilc", "env")
	envCommand.Stdout = buf

	err := envCommand.Start()
	if err != nil {
		return "", fmt.Errorf("Error starting shell: %v", err)
	}

	// Without explicitly calling Wait, tests hang forever
	// https://github.com/golang/go/issues/24050
	err = envCommand.Wait()

	if err != nil {
		return "", fmt.Errorf("Failed to run shell for retrieving environment variables: %v", err)
	}

	return buf.String(), nil
}

// extractPath takes in extracted environment variables from the user's shell, and returns the
// PATH environment variable from it. If PATH is not found, then it returns an empty string.
func extractPath(env string) string {
	for _, envVar := range strings.Split(env, "\n") {
		if strings.HasPrefix(envVar, "PATH=") {
			split := strings.Split(envVar, "=")
			if len(split) < 2 {
				// Discovered PATH is empty
				return ""
			}
			return split[1]
		}
	}
	return ""
}

// pathset is a makeshift Set datastructure using a map of empty structs
type pathset map[string]struct{}

// createPathset takes in the value of the PATH environment variable, and returns a set of
// unique path elements inside it
func createPathset(path string) pathset {
	set := make(pathset)
	for _, p := range strings.Split(path, string(os.PathListSeparator)) {
		if p == "" {
			continue
		}
		set[p] = struct{}{}
	}
	return set
}
