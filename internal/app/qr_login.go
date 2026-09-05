package app

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultQRLoginWait = 20 * time.Second
	mockLoginQRCode    = `    ██████████████  ██  ████  ██████████████
    ██          ██  ██████    ██          ██
    ██  ██████  ██    ██  ██  ██  ██████  ██
    ██  ██████  ██  ████      ██  ██████  ██
    ██  ██████  ██  ██  ████  ██  ██████  ██
    ██          ██    ████    ██          ██
    ██████████████  ██  ██  ██  ██████████████
                    ████████
    ██  ██  ██  ██      ██████  ████  ██  ██
      ██  ██████  ██  ██      ████  ██████
    ████      ██████    ██  ██████████    ██
    ██  ████      ██████  ████  ██  ██
      ████  ██████  ████        ██  ████████
                    ██  ██  ████      ██
    ██████████████  ██████  ████  ██  ██  ██
    ██          ██    ████████        ██
    ██  ██████  ██  ████    ████████████████
    ██  ██████  ██  ██  ████    ██      ██
    ██  ██████  ██    ██████  ██████  ██  ██
    ██          ██  ████      ██      ████
    ██████████████  ██  ██  ██████  ██  ████`
)

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

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Scan this QR code to log in:\n%s\nWaiting for confirmation...\n", mockLoginQRCode); err != nil {
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
