package app

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
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
		"  recall     Mock source recall results\n",
		"  logs       Stream mock logs for a duration\n",
		"  long-run   Run a long task with mixed persistent and overwritten output\n",
		"  qr-login   Simulate a QR-code login\n",
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

func TestRecallHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "recall", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock recall\n" +
		"\n" +
		"Description:\n" +
		"  Group mock source recall commands that emit recall result bodies.\n" +
		"\n" +
		"Available Commands:\n" +
		"  knowledge  Publish mock knowledge-base recall sources\n" +
		"  web-search Publish mock web search recall sources\n" +
		"  help       Help about any command\n" +
		"\n" +
		"Flags:\n" +
		"  -h, --help         help for this command\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestRecallKnowledgeHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "recall", "knowledge", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock recall knowledge [flags]\n" +
		"\n" +
		"Description:\n" +
		"  Return a fixed knowledge-base recall result body.\n" +
		"\n" +
		"Flags:\n" +
		"  --output string    Response format: text|json\n" +
		"  --query string     Recall query text\n" +
		"  -h, --help         help for this command\n" +
		"\n" +
		"Params fields:\n" +
		"  name     type     required   default                    description\n" +
		"  query    string   no         办公用品申请流程                   Recall query text\n" +
		"  output   string   no         text                       Response format: text|json\n" +
		"\n" +
		"Examples:\n" +
		"  mock recall knowledge\n" +
		"  mock recall knowledge --query 办公用品申请流程 --output json\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestRecallWebSearchHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "recall", "web-search", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock recall web-search [flags]\n" +
		"\n" +
		"Description:\n" +
		"  Return a fixed web search recall result body.\n" +
		"\n" +
		"Flags:\n" +
		"  --output string    Response format: text|json\n" +
		"  --query string     Recall query text\n" +
		"  -h, --help         help for this command\n" +
		"\n" +
		"Params fields:\n" +
		"  name     type     required   default                            description\n" +
		"  query    string   no         release checklist best practices   Recall query text\n" +
		"  output   string   no         text                               Response format: text|json\n" +
		"\n" +
		"Examples:\n" +
		"  mock recall web-search\n" +
		"  mock recall web-search --query \"release checklist best practices\" --output json\n"

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

func TestQRLoginHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "qr-login", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock qr-login [flags]\n" +
		"\n" +
		"Description:\n" +
		"  Print a mock QR code, wait for simulated confirmation, and complete the login successfully.\n" +
		"\n" +
		"Flags:\n" +
		"  --wait string      Time before automatic login approval\n" +
		"  -h, --help         help for this command\n" +
		"\n" +
		"Params fields:\n" +
		"  name   type     required   default   description\n" +
		"  wait   string   no         20s       Time to wait before the login is approved automatically\n" +
		"\n" +
		"Examples:\n" +
		"  mock qr-login\n" +
		"  mock qr-login --wait 1s\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestMockLoginQRCodeEncodesPayload(t *testing.T) {
	t.Parallel()

	rows := strings.Split(mockLoginQRMatrix, "\n")
	if len(rows) != 25 {
		t.Fatalf("expected a version 2 QR code with 25 rows, got %d", len(rows))
	}
	matrix := make([][]bool, len(rows))
	for rowIndex, row := range rows {
		if len(row) != len(rows) {
			t.Fatalf("row %d has %d modules, want %d", rowIndex, len(row), len(rows))
		}
		matrix[rowIndex] = make([]bool, len(row))
		for columnIndex, module := range row {
			switch module {
			case '0':
			case '1':
				matrix[rowIndex][columnIndex] = true
			default:
				t.Fatalf("row %d column %d has invalid module %q", rowIndex, columnIndex, module)
			}
		}
	}

	const formatLMask1 = 0x72f3
	primaryFormat, secondaryFormat := readQRFormatBits(matrix)
	if primaryFormat != formatLMask1 || secondaryFormat != formatLMask1 {
		t.Fatalf("unexpected format information: primary=%#x secondary=%#x", primaryFormat, secondaryFormat)
	}

	dataBits := readVersion2LMask1DataBits(matrix)
	if len(dataBits) != 359 {
		t.Fatalf("expected 359 data and remainder bits, got %d", len(dataBits))
	}
	for index, bit := range dataBits[352:] {
		if bit {
			t.Fatalf("remainder bit %d is set", index)
		}
	}
	codewords := qrBitsToBytes(dataBits[:352])
	if len(codewords) != 44 {
		t.Fatalf("expected 44 codewords, got %d", len(codewords))
	}
	for exponent := range 10 {
		value := byte(0)
		root := qrGFPower(exponent)
		for _, codeword := range codewords {
			value = qrGFMultiply(value, root) ^ codeword
		}
		if value != 0 {
			t.Fatalf("Reed-Solomon check failed at exponent %d: %#x", exponent, value)
		}
	}

	if mode := qrReadBits(dataBits, 0, 4); mode != 4 {
		t.Fatalf("expected byte mode, got %d", mode)
	}
	payloadLength := qrReadBits(dataBits, 4, 8)
	payload := make([]byte, payloadLength)
	for index := range payload {
		payload[index] = byte(qrReadBits(dataBits, 12+index*8, 8))
	}
	if string(payload) != mockLoginQRPayload {
		t.Fatalf("decoded payload %q, want %q", payload, mockLoginQRPayload)
	}

	renderedRows := strings.Split(renderMockLoginQRCode(), "\n")
	wantRenderedSize := len(rows) + 2*mockLoginQRQuietZone
	if len(renderedRows) != wantRenderedSize {
		t.Fatalf("rendered QR has %d rows, want %d", len(renderedRows), wantRenderedSize)
	}
	for rowIndex, row := range renderedRows {
		if modules := utf8.RuneCountInString(row) / 2; modules != wantRenderedSize {
			t.Fatalf("rendered row %d has %d modules, want %d", rowIndex, modules, wantRenderedSize)
		}
	}
}

func readQRFormatBits(matrix [][]bool) (int, int) {
	size := len(matrix)
	primary := 0
	secondary := 0
	set := func(target *int, bit int, value bool) {
		if value {
			*target |= 1 << bit
		}
	}
	for bit := range 6 {
		set(&primary, bit, matrix[bit][8])
	}
	set(&primary, 6, matrix[7][8])
	set(&primary, 7, matrix[8][8])
	set(&primary, 8, matrix[8][7])
	for bit := 9; bit < 15; bit++ {
		set(&primary, bit, matrix[8][14-bit])
	}
	for bit := range 8 {
		set(&secondary, bit, matrix[8][size-1-bit])
	}
	for bit := 8; bit < 15; bit++ {
		set(&secondary, bit, matrix[size-15+bit][8])
	}
	return primary, secondary
}

func readVersion2LMask1DataBits(matrix [][]bool) []bool {
	size := len(matrix)
	functionModule := make([][]bool, size)
	for row := range functionModule {
		functionModule[row] = make([]bool, size)
	}
	markRectangle := func(left, top, width, height int) {
		for row := top; row < top+height; row++ {
			for column := left; column < left+width; column++ {
				functionModule[row][column] = true
			}
		}
	}
	markRectangle(0, 0, 9, 9)
	markRectangle(size-8, 0, 8, 9)
	markRectangle(0, size-8, 9, 8)
	for index := 8; index < size-8; index++ {
		functionModule[6][index] = true
		functionModule[index][6] = true
	}
	markRectangle(16, 16, 5, 5)

	bits := make([]bool, 0, 359)
	upward := true
	for right := size - 1; right >= 1; right -= 2 {
		if right == 6 {
			right--
		}
		for verticalIndex := range size {
			row := verticalIndex
			if upward {
				row = size - 1 - verticalIndex
			}
			for _, column := range []int{right, right - 1} {
				if functionModule[row][column] {
					continue
				}
				module := matrix[row][column]
				if row%2 == 0 {
					module = !module
				}
				bits = append(bits, module)
			}
		}
		upward = !upward
	}
	return bits
}

func qrBitsToBytes(bits []bool) []byte {
	result := make([]byte, len(bits)/8)
	for index := range result {
		result[index] = byte(qrReadBits(bits, index*8, 8))
	}
	return result
}

func qrReadBits(bits []bool, offset, count int) int {
	value := 0
	for index := range count {
		value <<= 1
		if bits[offset+index] {
			value++
		}
	}
	return value
}

func qrGFPower(exponent int) byte {
	value := byte(1)
	for range exponent {
		value = qrGFMultiply(value, 2)
	}
	return value
}

func qrGFMultiply(left, right byte) byte {
	product := 0
	multiplicand := int(left)
	multiplier := int(right)
	for multiplier > 0 {
		if multiplier&1 != 0 {
			product ^= multiplicand
		}
		multiplier >>= 1
		multiplicand <<= 1
		if multiplicand&0x100 != 0 {
			multiplicand ^= 0x11d
		}
	}
	return byte(product)
}

func TestLogsHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "logs", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock logs [flags]\n" +
		"\n" +
		"Description:\n" +
		"  Continuously emit deterministic mock log lines for the requested duration.\n" +
		"\n" +
		"Flags:\n" +
		"  --duration string  Total time to emit logs\n" +
		"  --interval string  Delay between log lines\n" +
		"  -h, --help         help for this command\n" +
		"\n" +
		"Params fields:\n" +
		"  name       type     required   default   description\n" +
		"  duration   string   no         30s       Total time to emit logs\n" +
		"  interval   string   no         1s        Delay between log lines\n" +
		"\n" +
		"Examples:\n" +
		"  mock logs\n" +
		"  mock logs --duration 5s --interval 500ms\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestLogsCommandStreamsForRequestedDuration(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var waits []time.Duration
	cmd := newLogsCommand(func(duration time.Duration) {
		waits = append(waits, duration)
	})
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{"--duration", "250ms", "--interval", "100ms"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute logs: %v", err)
	}
	wantOutput := "" +
		"level=INFO sequence=1 elapsed=0s message=\"mock task is running\"\n" +
		"level=INFO sequence=2 elapsed=100ms message=\"mock task is running\"\n" +
		"level=INFO sequence=3 elapsed=200ms message=\"mock task is running\"\n" +
		"level=INFO sequence=4 elapsed=250ms message=\"mock task completed\"\n"
	if stdout.String() != wantOutput {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", wantOutput, stdout.String())
	}
	wantWaits := []time.Duration{100 * time.Millisecond, 100 * time.Millisecond, 50 * time.Millisecond}
	if len(waits) != len(wantWaits) {
		t.Fatalf("expected waits %v, got %v", wantWaits, waits)
	}
	for index := range wantWaits {
		if waits[index] != wantWaits[index] {
			t.Fatalf("expected waits %v, got %v", wantWaits, waits)
		}
	}
}

func TestLogsCommandRejectsNonPositiveDurations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "duration", args: []string{"logs", "--duration", "0s"}, want: "--duration must be greater than 0"},
		{name: "interval", args: []string{"logs", "--interval", "0s"}, want: "--interval must be greater than 0"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := runCommand(t, nil, test.args...)
			if result.code != ExitUsage {
				t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
			}
			if !strings.Contains(result.stderr, test.want) {
				t.Fatalf("expected stderr to contain %q, got %q", test.want, result.stderr)
			}
		})
	}
}

func TestLongRunHelp(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "long-run", "--help")

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d", ExitSuccess, result.code)
	}

	want := "" +
		"Usage:\n" +
		"  mock long-run [flags]\n" +
		"\n" +
		"Description:\n" +
		"  Continuously append permanent log lines while updating an overwritten progress line, then return a partially different final capture.\n" +
		"\n" +
		"Flags:\n" +
		"  --duration string  Total simulated task duration\n" +
		"  --interval string  Delay between progress updates\n" +
		"  -h, --help         help for this command\n" +
		"\n" +
		"Params fields:\n" +
		"  name       type     required   default   description\n" +
		"  duration   string   no         30s       Total simulated task duration\n" +
		"  interval   string   no         1s        Delay between progress updates\n" +
		"\n" +
		"Examples:\n" +
		"  mock long-run\n" +
		"  mock long-run --duration 5s --interval 500ms\n"

	if result.stdout != want {
		t.Fatalf("unexpected stdout:\nwant:\n%s\ngot:\n%s", want, result.stdout)
	}
}

func TestLongRunCommandOverwritesProgressButAppendsLogs(t *testing.T) {
	t.Parallel()

	outputPath := filepath.Join(t.TempDir(), "long-run.log")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("create long-run output: %v", err)
	}
	defer outputFile.Close()

	var liveOutput bytes.Buffer
	var liveOffset int
	readAppended := func() {
		data, readErr := os.ReadFile(outputPath)
		if readErr != nil {
			t.Fatalf("read long-run output: %v", readErr)
		}
		if len(data) > liveOffset {
			liveOutput.Write(data[liveOffset:])
			liveOffset = len(data)
		}
	}

	var waits []time.Duration
	cmd := newLongRunCommand(func(duration time.Duration) {
		readAppended()
		waits = append(waits, duration)
	})
	cmd.SetOut(outputFile)
	cmd.SetArgs([]string{"--duration", "250ms", "--interval", "100ms"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute long-run: %v", err)
	}
	readAppended()
	if err := outputFile.Sync(); err != nil {
		t.Fatalf("sync long-run output: %v", err)
	}
	finalData, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read final long-run output: %v", err)
	}
	finalOutput := string(finalData)

	if !strings.Contains(liveOutput.String(), "progress=000%") || !strings.Contains(liveOutput.String(), "status=working") {
		t.Fatalf("expected live output to retain initial progress, got %q", liveOutput.String())
	}
	if strings.Contains(liveOutput.String(), "status=success") {
		t.Fatalf("live append-only output unexpectedly contained overwritten final progress: %q", liveOutput.String())
	}
	if strings.Contains(finalOutput, "progress=000%") || !strings.Contains(finalOutput, "progress=100%") || !strings.Contains(finalOutput, "status=success") {
		t.Fatalf("expected final capture to contain only final progress, got %q", finalOutput)
	}
	for _, want := range []string{
		"event=start duration=250ms interval=100ms",
		"event=log sequence=1 elapsed=100ms progress=40%",
		"event=log sequence=2 elapsed=200ms progress=80%",
		"event=log sequence=3 elapsed=250ms progress=100%",
		"result=success duration=250ms logCount=3",
	} {
		if !strings.Contains(liveOutput.String(), want) || !strings.Contains(finalOutput, want) {
			t.Fatalf("expected permanent output %q in live and final captures\nlive: %q\nfinal: %q", want, liveOutput.String(), finalOutput)
		}
	}
	if liveOutput.String() == finalOutput {
		t.Fatal("expected live output and final capture to differ")
	}

	wantWaits := []time.Duration{100 * time.Millisecond, 100 * time.Millisecond, 50 * time.Millisecond}
	if len(waits) != len(wantWaits) {
		t.Fatalf("expected waits %v, got %v", wantWaits, waits)
	}
	for index := range wantWaits {
		if waits[index] != wantWaits[index] {
			t.Fatalf("expected waits %v, got %v", wantWaits, waits)
		}
	}
}

func TestLongRunOutputUsesANSIClearLineForTerminal(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	output := &longRunOutput{writer: &stdout, terminal: true}
	if err := output.WriteProgress(0, 0, time.Second, false); err != nil {
		t.Fatalf("write initial progress: %v", err)
	}
	if err := output.WriteProgress(50, 500*time.Millisecond, time.Second, false); err != nil {
		t.Fatalf("update progress: %v", err)
	}
	if err := output.WritePermanent("event=log sequence=1"); err != nil {
		t.Fatalf("write permanent log: %v", err)
	}
	if err := output.WriteProgress(100, time.Second, time.Second, true); err != nil {
		t.Fatalf("write completed progress: %v", err)
	}

	got := stdout.String()
	if count := strings.Count(got, longRunClearLine); count != 5 {
		t.Fatalf("expected 5 terminal clear-line sequences, got %d in %q", count, got)
	}
	wantUpdates := longRunClearLine + "progress=000% elapsed=0s/1s status=working" +
		longRunClearLine + "progress=050% elapsed=500ms/1s status=working"
	if !strings.HasPrefix(got, wantUpdates) {
		t.Fatalf("expected terminal progress updates without fixed-width padding, got %q", got)
	}
	if !strings.Contains(got, longRunClearLine+"event=log sequence=1\n"+longRunClearLine+"progress=050%") {
		t.Fatalf("expected the progress line to be cleared before the permanent log and redrawn afterward, got %q", got)
	}
	if !strings.HasSuffix(got, "\n") || !strings.HasSuffix(strings.TrimRight(got, " \n"), "status=success") {
		t.Fatalf("expected completed terminal progress to end with a newline, got %q", got)
	}
}

func TestLongRunCommandUsesCarriageReturnsForNonSeekableOutput(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	cmd := newLongRunCommand(func(time.Duration) {})
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{"--duration", "100ms", "--interval", "100ms"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute long-run: %v", err)
	}
	if !strings.Contains(stdout.String(), "\rprogress=000%") || !strings.Contains(stdout.String(), "status=success") {
		t.Fatalf("expected carriage-return progress fallback, got %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "event=log sequence=1") || !strings.Contains(stdout.String(), "result=success") {
		t.Fatalf("expected permanent fallback output, got %q", stdout.String())
	}
}

func TestLongRunCommandRejectsNonPositiveDurations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "duration", args: []string{"long-run", "--duration", "0s"}, want: "--duration must be greater than 0"},
		{name: "interval", args: []string{"long-run", "--interval", "0s"}, want: "--interval must be greater than 0"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := runCommand(t, nil, test.args...)
			if result.code != ExitUsage {
				t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
			}
			if !strings.Contains(result.stderr, test.want) {
				t.Fatalf("expected stderr to contain %q, got %q", test.want, result.stderr)
			}
		})
	}
}

func TestQRLoginCommandWaitsThenSucceeds(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var waited time.Duration
	cmd := newQRLoginCommand(func(duration time.Duration) {
		outputAtWait := stdout.String()
		if !strings.Contains(outputAtWait, "Waiting for confirmation...\n") {
			t.Errorf("expected QR prompt before wait, got %q", outputAtWait)
		}
		if strings.Contains(outputAtWait, "Login successful.\n") {
			t.Errorf("success was emitted before wait: %q", outputAtWait)
		}
		waited = duration
	})
	cmd.SetOut(&stdout)
	cmd.SetArgs(nil)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute qr-login: %v", err)
	}
	if waited != 20*time.Second {
		t.Fatalf("expected 20s wait, got %v", waited)
	}
	output := stdout.String()
	for _, want := range []string{
		"Scan this QR code to log in:\n",
		"██████████████",
		"Waiting for confirmation...\n",
		"Login successful.\n",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got %q", want, output)
		}
	}
	if strings.Index(output, "Waiting for confirmation...") > strings.Index(output, "Login successful.") {
		t.Fatalf("expected confirmation wait before success, got %q", output)
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
