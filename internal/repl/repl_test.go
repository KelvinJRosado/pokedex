package repl

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestRunCleanExit(t *testing.T) {
	var out bytes.Buffer
	Run(strings.NewReader("exit\n"), &out)

	got := out.String()
	if !strings.Contains(got, "Closing the Pokedex") {
		t.Errorf("expected farewell message, got: %q", got)
	}
	// After clean exit, Run should not print another prompt.
	if strings.HasSuffix(got, "Pokedex > ") {
		t.Errorf("expected no trailing prompt after exit, got: %q", got)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var out bytes.Buffer
	Run(strings.NewReader("nonsense\nexit\n"), &out)

	if !strings.Contains(out.String(), "Unknown command") {
		t.Errorf("expected 'Unknown command' in output, got: %q", out.String())
	}
}

func TestRunSkipsBlankLines(t *testing.T) {
	var out bytes.Buffer
	Run(strings.NewReader("\n   \nexit\n"), &out)

	if strings.Contains(out.String(), "Unknown command") {
		t.Errorf("blank lines should not be treated as unknown commands, got: %q", out.String())
	}
	if !strings.Contains(out.String(), "Closing the Pokedex") {
		t.Errorf("expected exit to still run after blank lines, got: %q", out.String())
	}
}

func TestRunSurfacesCommandError(t *testing.T) {
	// `mapb` on the first page errors before hitting the network.
	var out bytes.Buffer
	Run(strings.NewReader("mapb\nexit\n"), &out)

	if !strings.Contains(out.String(), "Error:") {
		t.Errorf("expected error message in output, got: %q", out.String())
	}
	if !strings.Contains(out.String(), "first page") {
		t.Errorf("expected error text from commandMapb, got: %q", out.String())
	}
}

func TestRunPromptsBetweenCommands(t *testing.T) {
	var out bytes.Buffer
	// help prints a lot of lines but should not affect prompt count.
	Run(strings.NewReader("help\nnonsense\nexit\n"), &out)

	// Initial prompt + one before each of the 2 non-exit lines processed = 3 prompts.
	// (The exit handler returns before printing another.)
	if got := strings.Count(out.String(), "Pokedex > "); got < 3 {
		t.Errorf("expected at least 3 prompts, got %d in output:\n%s", got, out.String())
	}
}

func TestRunReturnsOnEOF(t *testing.T) {
	var out bytes.Buffer
	// No newline, just EOF after a known command.
	Run(strings.NewReader("help"), &out)

	if !strings.Contains(out.String(), "help:") {
		t.Errorf("expected help command to have run before EOF, got: %q", out.String())
	}
}

// errReader returns an error on every Read call so we can exercise
// scanner.Err() handling.
type errReader struct{}

func (errReader) Read(p []byte) (int, error) { return 0, errors.New("synthetic read failure") }

func TestRunReportsScannerError(t *testing.T) {
	var out bytes.Buffer
	Run(errReader{}, &out)

	if !strings.Contains(out.String(), "Invalid input") {
		t.Errorf("expected scanner error to be reported, got: %q", out.String())
	}
}

// Compile-time assertion that errReader satisfies io.Reader.
var _ io.Reader = errReader{}
