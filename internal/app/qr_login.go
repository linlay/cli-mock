package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultQRLoginWait   = 20 * time.Second
	mockLoginQRPayload   = "https://example.com/?login=mock"
	mockLoginQRQuietZone = 4
	mockLoginQRMatrix    = `1111111011100010001111111
1000001011001110101000001
1011101001000100101011101
1011101001110110001011101
1011101010100010001011101
1000001011001010001000001
1111111010101010101111111
0000000010101011100000000
1110011011101010111110011
0011000111011101011101011
0011011110010001010011101
0111100101001010100101000
1011111011000011101100001
0111010010010111101100011
1110011001011001111001101
0010100011111010110111000
1100101101101000111110010
0000000011011100100010001
1111111000110100101010001
1000001011000001100010011
1011101001001001111110000
1011101001011100010010110
1011101011010001110111011
1000001010111011111110000
1111111011111000101001001`
)

func renderMockLoginQRCode() string {
	rows := strings.Split(mockLoginQRMatrix, "\n")
	quietRow := strings.Repeat(" ", (len(rows)+2*mockLoginQRQuietZone)*2)
	sideQuietZone := strings.Repeat(" ", mockLoginQRQuietZone*2)
	rendered := make([]string, 0, len(rows)+2*mockLoginQRQuietZone)

	for range mockLoginQRQuietZone {
		rendered = append(rendered, quietRow)
	}
	for _, row := range rows {
		var line strings.Builder
		line.Grow(len(quietRow))
		line.WriteString(sideQuietZone)
		for _, module := range row {
			if module == '1' {
				line.WriteString("██")
			} else {
				line.WriteString("  ")
			}
		}
		line.WriteString(sideQuietZone)
		rendered = append(rendered, line.String())
	}
	for range mockLoginQRQuietZone {
		rendered = append(rendered, quietRow)
	}

	return strings.Join(rendered, "\n")
}

func newQRLoginCommand(wait func(time.Duration)) *cobra.Command {
	waitRaw := defaultQRLoginWait.String()

	cmd := &cobra.Command{
		Use:         "qr-login",
		Short:       "Simulate a QR-code login",
		Description: "Print a mock QR code, wait for simulated confirmation, and complete the login successfully.",
		Example:     "mock qr-login\nmock qr-login --wait 1s",
		ParamFields: []cobra.HelpField{
			optionalField("wait", "string", defaultQRLoginWait.String(), "Time to wait before the login is approved automatically"),
		},
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			duration, err := parseDuration(waitRaw)
			if err != nil {
				return fmt.Errorf("invalid --wait %q: %w", waitRaw, err)
			}

			if _, err := fmt.Fprintf(
				cmd.OutOrStdout(),
				"Scan this QR code to log in:\n%s\nPayload: %s\nWaiting for confirmation...\n",
				renderMockLoginQRCode(),
				mockLoginQRPayload,
			); err != nil {
				return err
			}

			wait(duration)

			_, err = fmt.Fprintln(cmd.OutOrStdout(), "Login successful.")
			return err
		},
	}

	cmd.Flags().StringVar(&waitRaw, "wait", waitRaw, "Time before automatic login approval")
	return cmd
}
