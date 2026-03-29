package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/linlay/cli-mock/internal/buildinfo"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), buildinfo.Summary())
			return err
		},
	}
}

func newSleepCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sleep <duration>",
		Short: "Sleep for the requested duration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			duration, err := parseDuration(args[0])
			if err != nil {
				return err
			}
			time.Sleep(duration)
			return nil
		},
	}
}

func newEchoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "echo <text...>",
		Short: "Write text to stdout",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(args, " "))
			return err
		},
	}
}

func newStderrCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stderr <text...>",
		Short: "Write text to stderr",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.ErrOrStderr(), strings.Join(args, " "))
			return err
		},
	}
}

func newExitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "exit <code>",
		Short: "Exit with a specific code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			code, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid exit code %q: %w", args[0], err)
			}
			if code < 0 || code > 255 {
				return fmt.Errorf("exit code must be between 0 and 255")
			}
			return &exitError{Code: code}
		},
	}
}

func newFailCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fail [message...]",
		Short: "Write an error message and exit with code 1",
		RunE: func(cmd *cobra.Command, args []string) error {
			message := "mock failure"
			if len(args) > 0 {
				message = strings.Join(args, " ")
			}
			return failuref("%s", message)
		},
	}
}

func newJSONCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "json <raw-json>",
		Short: "Validate and emit compact JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var value any
			if err := json.Unmarshal([]byte(args[0]), &value); err != nil {
				return failuref("invalid JSON: %v", err)
			}
			data, err := json.Marshal(value)
			if err != nil {
				return failuref("marshal JSON: %v", err)
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(data))
			return err
		},
	}
}

func newArgsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "args <arg...>",
		Short: "Print positional arguments as JSON",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := json.Marshal(args)
			if err != nil {
				return failuref("marshal args: %v", err)
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(data))
			return err
		},
	}
}

func newEnvCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "env <key>",
		Short: "Print the value of an environment variable",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			value, ok := os.LookupEnv(args[0])
			if !ok {
				return failuref("environment variable %q is not set", args[0])
			}
			_, err := fmt.Fprintln(cmd.OutOrStdout(), value)
			return err
		},
	}
}

func newStdinCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stdin",
		Short: "Echo stdin to stdout",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return failuref("read stdin: %v", err)
			}
			_, err = io.Copy(cmd.OutOrStdout(), bytes.NewReader(data))
			return err
		},
	}
}

func newLinesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "lines <count>",
		Short: "Print numbered lines",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			count, err := parsePositiveCount(args[0])
			if err != nil {
				return err
			}

			for i := range count {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "line-%d\n", i+1); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newStreamCommand() *cobra.Command {
	intervalRaw := "1s"

	cmd := &cobra.Command{
		Use:   "stream <count>",
		Short: "Print lines with a delay between each line",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			count, err := parsePositiveCount(args[0])
			if err != nil {
				return err
			}

			interval, err := parseDuration(intervalRaw)
			if err != nil {
				return fmt.Errorf("invalid --interval %q: %w", intervalRaw, err)
			}

			for i := range count {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "line-%d\n", i+1); err != nil {
					return err
				}
				if i < count-1 {
					time.Sleep(interval)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&intervalRaw, "interval", intervalRaw, "Delay between streamed lines")
	return cmd
}

func parseDuration(raw string) (time.Duration, error) {
	duration, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid duration %q: %w", raw, err)
	}
	if duration < 0 {
		return 0, fmt.Errorf("duration must be non-negative")
	}
	return duration, nil
}

func parsePositiveCount(raw string) (int, error) {
	count, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid count %q: %w", raw, err)
	}
	if count <= 0 {
		return 0, fmt.Errorf("count must be greater than 0")
	}
	return count, nil
}
