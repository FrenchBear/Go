// golangci-lint.go
// Wrapper to golangci-lint.exe, eliminate VSCode hardcoed options --print-issued-lines=false and --out-format=colored-line-number
// not supported anymore by golangci-lint.exe
//
// 2025-06-12	PV		First version, code provoded by Gemini

package main

import (
	"fmt"     // Package for formatted I/O (e.g., printing error messages)
	"os"      // Package for operating system functions (e.g., command-line arguments, environment variables, exit codes)
	"os/exec" // Package for running external commands
	"strings" // Package for string manipulation (e.g., checking prefixes)
)

// golangciLintPath defines the path to the golangci-lint.exe executable.
const golangciLintPath = `C:\Utils\golangci-lint\golangci-lint.exe` // Assumes it's in PATH or current directory if left as is

func main() {
	// 1. Get all command-line arguments passed to this relay launcher.
	// os.Args[0] is the name of the program itself.
	// os.Args[1:] contains all arguments passed after the program name.
	args := os.Args[1:]

	// 2. Filter out unwanted arguments.
	// We'll create a new slice to store only the arguments we want to pass on.
	var filteredArgs []string
	for _, arg := range args {
		// Check if the argument starts with "--print-issued-lines=" or "--out-format=".
		// If it does, we skip it.
		if strings.HasPrefix(arg, "--print-issued-lines=") {
			//fmt.Fprintf(os.Stderr, "INFO: Filtering out argument: %s\n", arg)
			continue // Skip this argument
		}
		if strings.HasPrefix(arg, "--out-format=") {
			//fmt.Fprintf(os.Stderr, "INFO: Filtering out argument: %s\n", arg)
			continue // Skip this argument
		}
		// If it's not one of the unwanted arguments, add it to our filtered list.
		filteredArgs = append(filteredArgs, arg)
	}

	// 3. Prepare the command to execute golangci-lint.exe.
	// exec.Command takes the executable path as the first argument,
	// followed by all the arguments to pass to that executable.
	cmd := exec.Command(golangciLintPath, filteredArgs...)

	// 4. Redirect standard I/O (input, output, error) to the current process.
	// This ensures that:
	// - Any input expected by golangci-lint.exe comes from where this relay launcher is run.
	// - Any output from golangci-lint.exe (e.g., linting results) is printed to the console.
	// - Any errors from golangci-lint.exe are printed to the console.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 5. Run the command.
	err := cmd.Run()

	// 6. Handle errors and exit with the appropriate status code.
	if err != nil {
		// If the command failed to start or exited with a non-zero status, err will not be nil.
		if exitError, ok := err.(*exec.ExitError); ok {
			// If it's an ExitError, it means the command ran but returned a non-zero exit code.
			// We should exit with the same code to reflect the golangci-lint.exe's status.
			fmt.Fprintf(os.Stderr, "ERROR: golangci-lint.exe exited with error: %v\n", exitError)
			os.Exit(exitError.ExitCode()) // Exit with the same exit code as golangci-lint.exe
		} else {
			// Other types of errors (e.g., command not found, permission denied).
			fmt.Fprintf(os.Stderr, "ERROR: Failed to run golangci-lint.exe: %v\n", err)
			os.Exit(1) // Exit with a generic error code
		}
	}

	// If cmd.Run() returns nil, it means the command ran successfully and exited with status 0.
	os.Exit(0) // Exit successfully
}
