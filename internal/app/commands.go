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
		Use:         "version",
		Short:       "Print version information",
		Description: "Print the current mock CLI version string.",
		Example:     "mock version",
		Args:        cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), buildinfo.Summary())
			return err
		},
	}
}

func newSleepCommand() *cobra.Command {
	return &cobra.Command{
		Use:         "sleep <duration>",
		Short:       "Sleep for the requested duration",
		Description: "Pause execution for the requested duration.",
		Example:     "mock sleep 20ms\nmock sleep 1s",
		ArgFields: []cobra.HelpField{
			requiredField("duration", "string", "Duration in Go format, such as 20ms or 1s"),
		},
		Args: cobra.ExactArgs(1),
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
		Use:         "echo <text...>",
		Short:       "Write text to stdout",
		Description: "Write one or more text fragments to stdout as a single line.",
		Example:     "mock echo hello world",
		ArgFields: []cobra.HelpField{
			requiredField("text", "string[]", "One or more text fragments to join with spaces"),
		},
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(args, " "))
			return err
		},
	}
}

func newStderrCommand() *cobra.Command {
	return &cobra.Command{
		Use:         "stderr <text...>",
		Short:       "Write text to stderr",
		Description: "Write one or more text fragments to stderr as a single line.",
		Example:     "mock stderr warning message",
		ArgFields: []cobra.HelpField{
			requiredField("text", "string[]", "One or more text fragments to join with spaces"),
		},
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.ErrOrStderr(), strings.Join(args, " "))
			return err
		},
	}
}

func newExitCommand() *cobra.Command {
	return &cobra.Command{
		Use:         "exit <code>",
		Short:       "Exit with a specific code",
		Description: "Exit immediately with the provided process exit code.",
		Example:     "mock exit 7",
		ArgFields: []cobra.HelpField{
			requiredField("code", "integer", "Exit code, must be between 0 and 255"),
		},
		Args: cobra.ExactArgs(1),
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
		Use:         "fail [message...]",
		Short:       "Write an error message and exit with code 1",
		Description: "Write a failure message to stderr and exit with code 1.",
		Example:     "mock fail\nmock fail broken state",
		ArgFields: []cobra.HelpField{
			optionalField("message", "string[]", "mock failure", "Optional failure message fragments"),
		},
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
		Use:         "json <raw-json>",
		Short:       "Validate and emit compact JSON",
		Description: "Validate a JSON value and write it back in compact form.",
		Example:     "mock json '{\"ok\":true,\"count\":2}'",
		ArgFields: []cobra.HelpField{
			requiredField("raw-json", "string", "A valid JSON value to validate and compact"),
		},
		Args: cobra.ExactArgs(1),
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
		Use:         "args <arg...>",
		Short:       "Print positional arguments as JSON",
		Description: "Print all positional arguments as a JSON array.",
		Example:     "mock args one two three",
		ArgFields: []cobra.HelpField{
			requiredField("arg", "string[]", "One or more positional arguments to encode as JSON"),
		},
		Args: cobra.MinimumNArgs(1),
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
		Use:         "env <key>",
		Short:       "Print the value of an environment variable",
		Description: "Print the value of a single environment variable.",
		Example:     "mock env HOME",
		ArgFields: []cobra.HelpField{
			requiredField("key", "string", "Environment variable name to read"),
		},
		Args: cobra.ExactArgs(1),
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
		Use:         "stdin",
		Short:       "Echo stdin to stdout",
		Description: "Read all stdin content and write it back to stdout unchanged.",
		Example:     "printf 'first\\nsecond\\n' | mock stdin",
		Args:        cobra.NoArgs,
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
		Use:         "lines <count>",
		Short:       "Print numbered lines",
		Description: "Print a fixed number of numbered lines to stdout.",
		Example:     "mock lines 3",
		ArgFields: []cobra.HelpField{
			requiredField("count", "integer", "Number of lines to print, must be greater than 0"),
		},
		Args: cobra.ExactArgs(1),
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
		Use:         "stream <count> [content...]",
		Short:       "Print lines with a delay between each line",
		Description: "Print numbered lines with a delay between each line, or emit custom content sequentially.",
		Example:     "mock stream 3\nmock stream 3 --interval 100ms\nmock stream 3 hello world done --interval 100ms",
		ArgFields: []cobra.HelpField{
			requiredField("count", "integer", "Number of lines to print, must be greater than 0"),
			optionalField("content", "string[]", "-", "Optional lines to emit in order; when provided, the number of items must equal count"),
		},
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			count, err := parsePositiveCount(args[0])
			if err != nil {
				return err
			}

			interval, err := parseDuration(intervalRaw)
			if err != nil {
				return fmt.Errorf("invalid --interval %q: %w", intervalRaw, err)
			}

			lines := make([]string, 0, count)
			if len(args) > 1 {
				lines = append(lines, args[1:]...)
				if len(lines) != count {
					return fmt.Errorf("count %d does not match content item count %d", count, len(lines))
				}
			} else {
				for i := range count {
					lines = append(lines, fmt.Sprintf("line-%d", i+1))
				}
			}

			for i, line := range lines {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", line); err != nil {
					return err
				}
				if i < len(lines)-1 {
					time.Sleep(interval)
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&intervalRaw, "interval", intervalRaw, "Delay between streamed lines")
	return cmd
}

func requiredField(name, fieldType, description string) cobra.HelpField {
	return cobra.HelpField{
		Name:        name,
		Type:        fieldType,
		Required:    "yes",
		Default:     "-",
		Description: description,
	}
}

func optionalField(name, fieldType, defaultValue, description string) cobra.HelpField {
	return cobra.HelpField{
		Name:        name,
		Type:        fieldType,
		Required:    "no",
		Default:     defaultValue,
		Description: description,
	}
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
