package app

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

const (
	ExitSuccess = 0
	ExitFailure = 1
	ExitUsage   = 2
)

type exitError struct {
	Code int
	Err  error
}

func (e *exitError) Error() string {
	if e == nil || e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *exitError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func Execute(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	root := newRootCommand(stdin, stdout, stderr)
	root.SetArgs(args)

	if err := root.Execute(); err != nil {
		var exitErr *exitError
		if errors.As(err, &exitErr) {
			if exitErr.Err != nil {
				fmt.Fprintln(stderr, exitErr.Err.Error())
			}
			return exitErr.Code
		}

		fmt.Fprintln(stderr, err.Error())
		return ExitUsage
	}

	return ExitSuccess
}

func newRootCommand(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "mock",
		Short:         "Mock CLI for scripts and automation tests",
		Description:   "Mock CLI for scripts and automation tests.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.AddCommand(
		newVersionCommand(),
		newSleepCommand(),
		newEchoCommand(),
		newStderrCommand(),
		newExitCommand(),
		newFailCommand(),
		newJSONCommand(),
		newArgsCommand(),
		newEnvCommand(),
		newStdinCommand(),
		newLinesCommand(),
		newStreamCommand(),
		newCreateLeaveCommand(),
		newGetLeaveCommand(),
		newUpdateLeaveCommand(),
		newDeleteLeaveCommand(),
		newExpenseCommand(),
		newProcurementCommand(),
		newRecallCommand(),
		newXDGCommand(),
	)

	return cmd
}

func failuref(format string, args ...any) error {
	return &exitError{
		Code: ExitFailure,
		Err:  fmt.Errorf(format, args...),
	}
}
