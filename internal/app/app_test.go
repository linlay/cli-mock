package app

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestExecuteHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if !strings.Contains(result.stdout, "Mock CLI for scripts and automation tests") {
		t.Fatalf("expected help output, got %q", result.stdout)
	}
	if result.stderr != "" {
		t.Fatalf("expected empty stderr, got %q", result.stderr)
	}
}

func TestExecuteVersion(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "version")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stdout != "mock dev\n" {
		t.Fatalf("unexpected stdout: %q", result.stdout)
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "nope")

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, "unknown command") {
		t.Fatalf("expected unknown command error, got %q", result.stderr)
	}
}

func TestSleepCommand(t *testing.T) {
	t.Parallel()

	start := time.Now()
	result := runCommand(t, nil, "sleep", "20ms")
	elapsed := time.Since(start)

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if elapsed < 15*time.Millisecond {
		t.Fatalf("expected sleep duration, got %v", elapsed)
	}
}

func TestSleepCommandRejectsInvalidDuration(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "sleep", "oops")

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, "invalid duration") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestEchoCommand(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "echo", "hello", "world")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stdout != "hello world\n" {
		t.Fatalf("unexpected stdout: %q", result.stdout)
	}
}

func TestStderrCommand(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "stderr", "hello", "err")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stderr != "hello err\n" {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestExitCommand(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "exit", "7")

	if result.code != 7 {
		t.Fatalf("expected exit 7, got %d", result.code)
	}
	if result.stdout != "" || result.stderr != "" {
		t.Fatalf("expected no output, got stdout=%q stderr=%q", result.stdout, result.stderr)
	}
}

func TestExitCommandRejectsInvalidCode(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "exit", "999")

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, "between 0 and 255") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestFailCommand(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "fail", "broken", "state")

	if result.code != ExitFailure {
		t.Fatalf("expected exit %d, got %d", ExitFailure, result.code)
	}
	if result.stderr != "broken state\n" {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestFailCommandDefaultMessage(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "fail")

	if result.code != ExitFailure {
		t.Fatalf("expected exit %d, got %d", ExitFailure, result.code)
	}
	if result.stderr != "mock failure\n" {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestJSONCommand(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "json", "{\"b\":2,\"a\":1}")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stdout != "{\"a\":1,\"b\":2}\n" {
		t.Fatalf("unexpected stdout: %q", result.stdout)
	}
}

func TestJSONCommandRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "json", "{bad}")

	if result.code != ExitFailure {
		t.Fatalf("expected exit %d, got %d", ExitFailure, result.code)
	}
	if !strings.Contains(result.stderr, "invalid JSON") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestArgsCommand(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "args", "one", "two", "three")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stdout != "[\"one\",\"two\",\"three\"]\n" {
		t.Fatalf("unexpected stdout: %q", result.stdout)
	}
}

func TestEnvCommand(t *testing.T) {
	t.Setenv("MOCK_TEST_KEY", "value-123")
	result := runCommand(t, nil, "env", "MOCK_TEST_KEY")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stdout != "value-123\n" {
		t.Fatalf("unexpected stdout: %q", result.stdout)
	}
}

func TestEnvCommandMissingVariable(t *testing.T) {
	result := runCommand(t, nil, "env", "MOCK_TEST_KEY_MISSING")

	if result.code != ExitFailure {
		t.Fatalf("expected exit %d, got %d", ExitFailure, result.code)
	}
	if !strings.Contains(result.stderr, "is not set") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestStdinCommand(t *testing.T) {
	t.Parallel()

	result := runCommand(t, bytes.NewBufferString("first\nsecond\n"), "stdin")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stdout != "first\nsecond\n" {
		t.Fatalf("unexpected stdout: %q", result.stdout)
	}
}

func TestLinesCommand(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "lines", "3")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stdout != "line-1\nline-2\nline-3\n" {
		t.Fatalf("unexpected stdout: %q", result.stdout)
	}
}

func TestLinesCommandRejectsInvalidCount(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "lines", "0")

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, "greater than 0") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

func TestStreamCommand(t *testing.T) {
	t.Parallel()

	start := time.Now()
	result := runCommand(t, nil, "stream", "3", "--interval", "15ms")
	elapsed := time.Since(start)

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stdout != "line-1\nline-2\nline-3\n" {
		t.Fatalf("unexpected stdout: %q", result.stdout)
	}
	if elapsed < 25*time.Millisecond {
		t.Fatalf("expected streaming delay, got %v", elapsed)
	}
}

func TestStreamCommandRejectsInvalidCount(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "stream", "-1")

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, "greater than 0") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}

type commandResult struct {
	code   int
	stdout string
	stderr string
}

func runCommand(t *testing.T, stdin *bytes.Buffer, args ...string) commandResult {
	t.Helper()

	if stdin == nil {
		stdin = &bytes.Buffer{}
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Execute(args, stdin, &stdout, &stderr)
	return commandResult{
		code:   code,
		stdout: stdout.String(),
		stderr: stderr.String(),
	}
}
