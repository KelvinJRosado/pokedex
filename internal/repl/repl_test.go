package repl

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/kelvinjrosado/pokedex/internal/logger"
)

// newTestLogger returns a logger whose console handler writes to the returned
// buffer, so tests can assert on what was logged. The file handler discards.
func newTestLogger() (*logger.CustomLogger, *bytes.Buffer) {
	var logBuf bytes.Buffer
	return logger.NewWithWriters(&logBuf, io.Discard), &logBuf
}

func TestRunCleanExit(t *testing.T) {

	var out bytes.Buffer
	lgr, _ := newTestLogger()
	Run(strings.NewReader("exit\n"), &out, lgr)

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
	lgr, logBuf := newTestLogger()
	Run(strings.NewReader("nonsense\nexit\n"), &out, lgr)

	if !strings.Contains(logBuf.String(), "Unknown command") {
		t.Errorf("expected 'Unknown command' in log output, got: %q", logBuf.String())
	}
}

func TestRunSkipsBlankLines(t *testing.T) {

	var out bytes.Buffer
	lgr, logBuf := newTestLogger()
	Run(strings.NewReader("\n   \nexit\n"), &out, lgr)

	if strings.Contains(logBuf.String(), "Unknown command") {
		t.Errorf("blank lines should not be treated as unknown commands, got: %q", logBuf.String())
	}
	if !strings.Contains(out.String(), "Closing the Pokedex") {
		t.Errorf("expected exit to still run after blank lines, got: %q", out.String())
	}
}

func TestRunSurfacesCommandError(t *testing.T) {

	// `mapb` on the first page errors before hitting the network.
	var out bytes.Buffer
	lgr, logBuf := newTestLogger()
	Run(strings.NewReader("mapb\nexit\n"), &out, lgr)

	if !strings.Contains(logBuf.String(), "Error executing command") {
		t.Errorf("expected error log, got: %q", logBuf.String())
	}
	if !strings.Contains(logBuf.String(), "first page") {
		t.Errorf("expected error text from commandMapb, got: %q", logBuf.String())
	}
}

func TestRunPromptsBetweenCommands(t *testing.T) {

	var out bytes.Buffer
	lgr, _ := newTestLogger()
	// help prints a lot of lines but should not affect prompt count.
	Run(strings.NewReader("help\nnonsense\nexit\n"), &out, lgr)

	// Initial prompt + one before each of the 2 non-exit lines processed = 3 prompts.
	// (The exit handler returns before printing another.)
	if got := strings.Count(out.String(), "Pokedex > "); got < 3 {
		t.Errorf("expected at least 3 prompts, got %d in output:\n%s", got, out.String())
	}
}

func TestRunReturnsOnEOF(t *testing.T) {

	var out bytes.Buffer
	lgr, _ := newTestLogger()
	// No newline, just EOF after a known command.
	Run(strings.NewReader("help"), &out, lgr)

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
	lgr, logBuf := newTestLogger()
	Run(errReader{}, &out, lgr)

	if !strings.Contains(logBuf.String(), "Invalid input") {
		t.Errorf("expected scanner error to be reported, got: %q", logBuf.String())
	}
}

// Compile-time assertion that errReader satisfies io.Reader.
var _ io.Reader = errReader{}
