package app

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultLogsDuration = 30 * time.Second
	defaultLogsInterval = time.Second
)

func newLogsCommand(wait func(time.Duration)) *cobra.Command {
	durationRaw := defaultLogsDuration.String()
	intervalRaw := defaultLogsInterval.String()

	cmd := &cobra.Command{
		Use:         "logs",
		Short:       "Stream mock logs for a duration",
		Description: "Continuously emit deterministic mock log lines for the requested duration.",
		Example:     "mock logs\nmock logs --duration 5s --interval 500ms",
		ParamFields: []cobra.HelpField{
			optionalField("duration", "string", defaultLogsDuration.String(), "Total time to emit logs"),
			optionalField("interval", "string", defaultLogsInterval.String(), "Delay between log lines"),
		},
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			duration, err := parseDuration(durationRaw)
			if err != nil {
				return fmt.Errorf("invalid --duration %q: %w", durationRaw, err)
			}
			if duration <= 0 {
				return fmt.Errorf("--duration must be greater than 0")
			}

			interval, err := parseDuration(intervalRaw)
			if err != nil {
				return fmt.Errorf("invalid --interval %q: %w", intervalRaw, err)
			}
			if interval <= 0 {
				return fmt.Errorf("--interval must be greater than 0")
			}

			elapsed := time.Duration(0)
			sequence := 1
			for elapsed < duration {
				if _, err := fmt.Fprintf(
					cmd.OutOrStdout(),
					"level=INFO sequence=%d elapsed=%s message=\"mock task is running\"\n",
					sequence,
					elapsed,
				); err != nil {
					return err
				}

				step := interval
				if remaining := duration - elapsed; remaining < step {
					step = remaining
				}
				wait(step)
				elapsed += step
				sequence++
			}

			_, err = fmt.Fprintf(
				cmd.OutOrStdout(),
				"level=INFO sequence=%d elapsed=%s message=\"mock task completed\"\n",
				sequence,
				elapsed,
			)
			return err
		},
	}

	cmd.Flags().StringVar(&durationRaw, "duration", durationRaw, "Total time to emit logs")
	cmd.Flags().StringVar(&intervalRaw, "interval", intervalRaw, "Delay between log lines")
	return cmd
}
