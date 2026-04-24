package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
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
	for _, want := range []string{
		"Usage:\n  mock\n",
		"Description:\n  Mock CLI for scripts and automation tests.\n",
		"Available Commands:\n",
		"  create-leave Create a mock leave application\n",
		"  get-leave  Get a mock leave application\n",
		"  expense    Mock expense reimbursement commands\n",
		"  procurement Mock procurement request commands\n",
		"  stream     Print lines with a delay between each line\n",
		"Flags:\n  -h, --help         help for this command\n",
	} {
		if !strings.Contains(result.stdout, want) {
			t.Fatalf("expected help output to contain %q, got %q", want, result.stdout)
		}
	}
	if result.stderr != "" {
		t.Fatalf("expected empty stderr, got %q", result.stderr)
	}
}

func TestVersionHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "version", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock version\n" +
		"\n" +
		"Description:\n" +
		"  Print the current mock CLI version string.\n" +
		"\n" +
		"Flags:\n" +
		"  -h, --help         help for this command\n" +
		"\n" +
		"Examples:\n" +
		"  mock version\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestExpenseHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "expense", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock expense\n" +
		"\n" +
		"Description:\n" +
		"  Group mock expense reimbursement commands under a resource-style namespace.\n" +
		"\n" +
		"Available Commands:\n" +
		"  add        Add a mock expense reimbursement\n" +
		"  get        Get a mock expense reimbursement\n" +
		"  update     Update a mock expense reimbursement\n" +
		"  delete     Delete a mock expense reimbursement\n" +
		"  help       Help about any command\n" +
		"\n" +
		"Flags:\n" +
		"  -h, --help         help for this command\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestProcurementHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "procurement", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock procurement\n" +
		"\n" +
		"Description:\n" +
		"  Group mock procurement request commands under a resource-style namespace.\n" +
		"\n" +
		"Available Commands:\n" +
		"  create     Create a mock procurement request\n" +
		"  get        Get a mock procurement request\n" +
		"  update     Update a mock procurement request\n" +
		"  delete     Delete a mock procurement request\n" +
		"  help       Help about any command\n" +
		"\n" +
		"Flags:\n" +
		"  -h, --help         help for this command\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestEnvHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "env", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock env <key>\n" +
		"\n" +
		"Description:\n" +
		"  Print the value of a single environment variable.\n" +
		"\n" +
		"Flags:\n" +
		"  -h, --help         help for this command\n" +
		"\n" +
		"Args fields:\n" +
		"  name   type     required   default   description\n" +
		"  key    string   yes        -         Environment variable name to read\n" +
		"\n" +
		"Examples:\n" +
		"  mock env HOME\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestStreamHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "stream", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock stream <count> [content...] [flags]\n" +
		"\n" +
		"Description:\n" +
		"  Print numbered lines with a delay between each line, or emit custom content sequentially.\n" +
		"\n" +
		"Flags:\n" +
		"  --interval string  Delay between streamed lines\n" +
		"  -h, --help         help for this command\n" +
		"\n" +
		"Args fields:\n" +
		"  name      type       required   default   description\n" +
		"  count     integer    yes        -         Number of lines to print, must be greater than 0\n" +
		"  content   string[]   no         -         Optional lines to emit in order; when provided, the number of items must equal count\n" +
		"\n" +
		"Examples:\n" +
		"  mock stream 3\n" +
		"  mock stream 3 --interval 100ms\n" +
		"  mock stream 3 hello world done --interval 100ms\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestXDGHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "xdg", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock xdg\n" +
		"\n" +
		"Description:\n" +
		"  Create and inspect mock .config and .local trees under an explicit root.\n" +
		"\n" +
		"Available Commands:\n" +
		"  apply      Create a mock XDG tree from a JSON manifest\n" +
		"  inspect    Inspect a mock XDG tree as JSON\n" +
		"  help       Help about any command\n" +
		"\n" +
		"Flags:\n" +
		"  -h, --help         help for this command\n" +
		"\n" +
		"Examples:\n" +
		"  mock xdg apply --root /tmp/mock-home --manifest ./manifest.json\n" +
		"  mock xdg inspect --root /tmp/mock-home\n" +
		"  mock xdg inspect --root /tmp/mock-home --reveal\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestXDGApplyHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "xdg", "apply", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock xdg apply [flags]\n" +
		"\n" +
		"Description:\n" +
		"  Create .config and .local content under the provided root using a JSON manifest.\n" +
		"\n" +
		"Flags:\n" +
		"  --manifest string  JSON manifest path or - for stdin\n" +
		"  --overwrite bool   Overwrite existing files\n" +
		"  --root string      Fake home root to manage\n" +
		"  -h, --help         help for this command\n" +
		"\n" +
		"Examples:\n" +
		"  mock xdg apply --root /tmp/mock-home --manifest ./manifest.json\n" +
		"  mock xdg apply --root /tmp/mock-home --manifest - --overwrite\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestXDGInspectHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "xdg", "inspect", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock xdg inspect [flags]\n" +
		"\n" +
		"Description:\n" +
		"  Inspect .config and .local metadata under the provided root and optionally reveal file content.\n" +
		"\n" +
		"Flags:\n" +
		"  --reveal bool      Include readable text and JSON file content\n" +
		"  --root string      Fake home root to inspect\n" +
		"  -h, --help         help for this command\n" +
		"\n" +
		"Examples:\n" +
		"  mock xdg inspect --root /tmp/mock-home\n" +
		"  mock xdg inspect --root /tmp/mock-home --reveal\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestHelpCommandMatchesFlagHelp(t *testing.T) {
	t.Parallel()

	for _, path := range [][]string{
		{"version"},
		{"env"},
		{"stream"},
		{"create-leave"},
		{"get-leave"},
		{"update-leave"},
		{"delete-leave"},
		{"expense"},
		{"expense", "add"},
		{"expense", "get"},
		{"expense", "update"},
		{"expense", "delete"},
		{"procurement"},
		{"procurement", "create"},
		{"procurement", "get"},
		{"procurement", "update"},
		{"procurement", "delete"},
	} {
		helpArgs := append([]string{"help"}, path...)
		flagArgs := append(append([]string(nil), path...), "--help")
		resultFromHelp := runCommand(t, nil, helpArgs...)
		resultFromFlag := runCommand(t, nil, flagArgs...)
		name := strings.Join(path, " ")

		if resultFromHelp.code != ExitSuccess {
			t.Fatalf("expected help %s to exit %d, got %d", name, ExitSuccess, resultFromHelp.code)
		}
		if resultFromFlag.code != ExitSuccess {
			t.Fatalf("expected %s --help to exit %d, got %d", name, ExitSuccess, resultFromFlag.code)
		}
		if resultFromHelp.stdout != resultFromFlag.stdout {
			t.Fatalf("expected help outputs for %s to match\nhelp:\n%s\nflag:\n%s", name, resultFromHelp.stdout, resultFromFlag.stdout)
		}
	}
}

func TestExecuteVersion(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "version")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stdout != "mock dev (commit none, built unknown)\n" {
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

func TestStreamCommandWithCustomContent(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "stream", "3", "hello", "world", "done")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stdout != "hello\nworld\ndone\n" {
		t.Fatalf("unexpected stdout: %q", result.stdout)
	}
}

func TestStreamCommandWithCustomContentDelay(t *testing.T) {
	t.Parallel()

	start := time.Now()
	result := runCommand(t, nil, "stream", "3", "hello", "world", "done", "--interval", "15ms")
	elapsed := time.Since(start)

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}
	if result.stdout != "hello\nworld\ndone\n" {
		t.Fatalf("unexpected stdout: %q", result.stdout)
	}
	if elapsed < 25*time.Millisecond {
		t.Fatalf("expected streaming delay, got %v", elapsed)
	}
}

func TestStreamCommandRejectsMismatchedContentCount(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "stream", "3", "hello", "world")

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, "does not match content item count") {
		t.Fatalf("unexpected stderr: %q", result.stderr)
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

func TestXDGApplyCreatesTree(t *testing.T) {
	t.Parallel()

	root := filepath.Join(t.TempDir(), "mock-home")
	manifestPath := writeManifestFile(t, `{
  "entries": [
    {
      "path": ".config/demo/config.toml",
      "type": "file",
      "format": "text",
      "content": "token = \"demo\"\n"
    },
    {
      "path": ".local/share/demo/secret.json",
      "type": "file",
      "format": "json",
      "content": { "api_key": "demo-key" }
    },
    {
      "path": ".local/state/demo",
      "type": "dir",
      "mode": "0700"
    }
  ]
}`)

	result := runCommand(t, nil, "xdg", "apply", "--root", root, "--manifest", manifestPath)

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d stderr=%q", ExitSuccess, result.code, result.stderr)
	}
	if result.stderr != "" {
		t.Fatalf("expected empty stderr, got %q", result.stderr)
	}
	for _, want := range []string{
		"root: " + root + "\n",
		"created:\n- .config/demo/config.toml\n- .local/share/demo/secret.json\n- .local/state/demo\n",
		"updated:\n-\n",
		"HOME=" + root,
		"XDG_CONFIG_HOME=" + filepath.Join(root, ".config"),
		"XDG_DATA_HOME=" + filepath.Join(root, ".local", "share"),
		"XDG_STATE_HOME=" + filepath.Join(root, ".local", "state"),
	} {
		if !strings.Contains(result.stdout, want) {
			t.Fatalf("expected stdout to contain %q, got %q", want, result.stdout)
		}
	}

	configContent, err := os.ReadFile(filepath.Join(root, ".config", "demo", "config.toml"))
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	if string(configContent) != "token = \"demo\"\n" {
		t.Fatalf("unexpected config content: %q", string(configContent))
	}

	secretContent, err := os.ReadFile(filepath.Join(root, ".local", "share", "demo", "secret.json"))
	if err != nil {
		t.Fatalf("read secret file: %v", err)
	}
	if string(secretContent) != "{\"api_key\":\"demo-key\"}\n" {
		t.Fatalf("unexpected secret content: %q", string(secretContent))
	}

	configInfo, err := os.Stat(filepath.Join(root, ".config", "demo", "config.toml"))
	if err != nil {
		t.Fatalf("stat config file: %v", err)
	}
	if got := configInfo.Mode().Perm(); got != 0o644 {
		t.Fatalf("unexpected config mode: %04o", got)
	}

	secretInfo, err := os.Stat(filepath.Join(root, ".local", "share", "demo", "secret.json"))
	if err != nil {
		t.Fatalf("stat secret file: %v", err)
	}
	if got := secretInfo.Mode().Perm(); got != 0o600 {
		t.Fatalf("unexpected secret mode: %04o", got)
	}

	stateInfo, err := os.Stat(filepath.Join(root, ".local", "state"))
	if err != nil {
		t.Fatalf("stat state dir: %v", err)
	}
	if got := stateInfo.Mode().Perm(); got != 0o700 {
		t.Fatalf("unexpected state mode: %04o", got)
	}

	demoStateInfo, err := os.Stat(filepath.Join(root, ".local", "state", "demo"))
	if err != nil {
		t.Fatalf("stat nested state dir: %v", err)
	}
	if got := demoStateInfo.Mode().Perm(); got != 0o700 {
		t.Fatalf("unexpected nested state mode: %04o", got)
	}
}

func TestXDGApplyRejectsInvalidPaths(t *testing.T) {
	t.Parallel()

	root := filepath.Join(t.TempDir(), "mock-home")
	tests := []struct {
		name    string
		entry   string
		wantErr string
	}{
		{
			name:    "absolute",
			entry:   `{"path":"/tmp/nope","type":"file","format":"text","content":"x"}`,
			wantErr: "must be relative",
		},
		{
			name:    "traversal",
			entry:   `{"path":"../.config/nope","type":"file","format":"text","content":"x"}`,
			wantErr: "must not escape the root",
		},
		{
			name:    "outside managed dirs",
			entry:   `{"path":".cache/demo","type":"file","format":"text","content":"x"}`,
			wantErr: "must start with .config/ or .local/",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			manifestPath := writeManifestFile(t, `{"entries":[`+tt.entry+`]}`)
			result := runCommand(t, nil, "xdg", "apply", "--root", root, "--manifest", manifestPath)

			if result.code != ExitUsage {
				t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
			}
			if !strings.Contains(result.stderr, tt.wantErr) {
				t.Fatalf("expected stderr to contain %q, got %q", tt.wantErr, result.stderr)
			}
		})
	}
}

func TestXDGApplyRejectsExistingFileWithoutOverwrite(t *testing.T) {
	t.Parallel()

	root := filepath.Join(t.TempDir(), "mock-home")
	firstManifest := writeManifestFile(t, `{"entries":[{"path":".config/demo/config.toml","type":"file","format":"text","content":"first\n"}]}`)
	secondManifest := writeManifestFile(t, `{"entries":[{"path":".config/demo/config.toml","type":"file","format":"text","content":"second\n"}]}`)

	first := runCommand(t, nil, "xdg", "apply", "--root", root, "--manifest", firstManifest)
	if first.code != ExitSuccess {
		t.Fatalf("first apply failed: code=%d stderr=%q", first.code, first.stderr)
	}

	second := runCommand(t, nil, "xdg", "apply", "--root", root, "--manifest", secondManifest)
	if second.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, second.code)
	}
	if !strings.Contains(second.stderr, "already exists; use --overwrite") {
		t.Fatalf("unexpected stderr: %q", second.stderr)
	}

	content, err := os.ReadFile(filepath.Join(root, ".config", "demo", "config.toml"))
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	if string(content) != "first\n" {
		t.Fatalf("unexpected config content after failed overwrite: %q", string(content))
	}
}

func TestXDGApplyOverwriteUpdatesExistingFile(t *testing.T) {
	t.Parallel()

	root := filepath.Join(t.TempDir(), "mock-home")
	firstManifest := writeManifestFile(t, `{"entries":[{"path":".config/demo/config.toml","type":"file","format":"text","content":"first\n"}]}`)
	secondManifest := writeManifestFile(t, `{"entries":[{"path":".config/demo/config.toml","type":"file","format":"text","content":"second\n"}]}`)

	first := runCommand(t, nil, "xdg", "apply", "--root", root, "--manifest", firstManifest)
	if first.code != ExitSuccess {
		t.Fatalf("first apply failed: code=%d stderr=%q", first.code, first.stderr)
	}

	second := runCommand(t, nil, "xdg", "apply", "--root", root, "--manifest", secondManifest, "--overwrite")
	if second.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d stderr=%q", ExitSuccess, second.code, second.stderr)
	}
	if !strings.Contains(second.stdout, "updated:\n- .config/demo/config.toml\n") {
		t.Fatalf("expected updated section, got %q", second.stdout)
	}

	content, err := os.ReadFile(filepath.Join(root, ".config", "demo", "config.toml"))
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	if string(content) != "second\n" {
		t.Fatalf("unexpected config content after overwrite: %q", string(content))
	}
}

func TestXDGInspectMetadataAndReveal(t *testing.T) {
	t.Parallel()

	root := filepath.Join(t.TempDir(), "mock-home")
	manifestPath := writeManifestFile(t, `{
  "entries": [
    {
      "path": ".config/demo/config.toml",
      "type": "file",
      "format": "text",
      "content": "token = \"demo\"\n"
    },
    {
      "path": ".local/share/demo/secret.json",
      "type": "file",
      "format": "json",
      "content": { "api_key": "demo-key" }
    }
  ]
}`)

	applied := runCommand(t, nil, "xdg", "apply", "--root", root, "--manifest", manifestPath)
	if applied.code != ExitSuccess {
		t.Fatalf("apply failed: code=%d stderr=%q", applied.code, applied.stderr)
	}

	inspect := runCommand(t, nil, "xdg", "inspect", "--root", root)
	if inspect.code != ExitSuccess {
		t.Fatalf("inspect failed: code=%d stderr=%q", inspect.code, inspect.stderr)
	}

	var metadata xdgInspectResponse
	if err := json.Unmarshal([]byte(inspect.stdout), &metadata); err != nil {
		t.Fatalf("unmarshal metadata inspect: %v", err)
	}
	if metadata.Root != root {
		t.Fatalf("unexpected root: %q", metadata.Root)
	}
	if metadata.Reveal {
		t.Fatalf("expected reveal=false, got true")
	}
	if len(metadata.Entries) == 0 {
		t.Fatalf("expected inspect entries")
	}
	for _, entry := range metadata.Entries {
		if entry.Content != nil {
			t.Fatalf("expected hidden content in metadata inspect: %#v", entry)
		}
	}

	revealed := runCommand(t, nil, "xdg", "inspect", "--root", root, "--reveal")
	if revealed.code != ExitSuccess {
		t.Fatalf("revealed inspect failed: code=%d stderr=%q", revealed.code, revealed.stderr)
	}

	var full xdgInspectResponse
	if err := json.Unmarshal([]byte(revealed.stdout), &full); err != nil {
		t.Fatalf("unmarshal revealed inspect: %v", err)
	}

	foundText := false
	foundJSON := false
	for _, entry := range full.Entries {
		switch entry.Path {
		case ".config/demo/config.toml":
			foundText = true
			if entry.Format != "text" {
				t.Fatalf("expected text format, got %#v", entry)
			}
			content, ok := entry.Content.(string)
			if !ok || content != "token = \"demo\"\n" {
				t.Fatalf("unexpected text content: %#v", entry)
			}
		case ".local/share/demo/secret.json":
			foundJSON = true
			if entry.Format != "json" {
				t.Fatalf("expected json format, got %#v", entry)
			}
			content, ok := entry.Content.(map[string]any)
			if !ok {
				t.Fatalf("unexpected json content type: %#v", entry)
			}
			if got, ok := content["api_key"].(string); !ok || got != "demo-key" {
				t.Fatalf("unexpected json content: %#v", entry)
			}
		}
	}
	if !foundText || !foundJSON {
		t.Fatalf("expected both revealed files, got %#v", full.Entries)
	}
}

func TestXDGInspectHandlesMissingManagedDirs(t *testing.T) {
	t.Parallel()

	root := filepath.Join(t.TempDir(), "missing-home")
	result := runCommand(t, nil, "xdg", "inspect", "--root", root)

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d stderr=%q", ExitSuccess, result.code, result.stderr)
	}

	var inspect xdgInspectResponse
	if err := json.Unmarshal([]byte(result.stdout), &inspect); err != nil {
		t.Fatalf("unmarshal inspect output: %v", err)
	}
	if len(inspect.Entries) != 0 {
		t.Fatalf("expected empty entries, got %#v", inspect.Entries)
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

func writeManifestFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	return path
}

func mustJSONArg(t *testing.T, value any) string {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal JSON arg: %v", err)
	}
	return string(data)
}

func mustJSONLine(t *testing.T, value any) string {
	t.Helper()

	return mustJSONArg(t, value) + "\n"
}
