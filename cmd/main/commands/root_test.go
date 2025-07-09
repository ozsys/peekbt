package commands

import (
	"bytes"
	"strings"
	"testing"
)

// TestExecute_NoArgs ensures that running without arguments returns the expected error.
func TestExecute_NoArgs(t *testing.T) {
	// Clear any previous args and set no args
	rootCommand.SetArgs([]string{})
	// Capture output
	buf := &bytes.Buffer{}
	rootCommand.SetOut(buf)
	// Execute
	err := Execute()
	if err == nil || err.Error() != "no command specified" {
		t.Fatalf("expected error 'no command specified', got %v, output: %s", err, buf.String())
	}
}

// TestExecute_HelpFlag verifies that the --help flag prints usage and does not error.
func TestExecute_HelpFlag(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCommand.SetOut(buf)
	rootCommand.SetErr(buf)

	rootCommand.SetArgs([]string{"--help"})

	_ = rootCommand.Execute()

	got := buf.String()
	if !strings.Contains(got, "Usage") {
		t.Errorf("help output should contain 'Usage', got:\n%s", got)
	}
}
