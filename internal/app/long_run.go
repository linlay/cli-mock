package app

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultLongRunDuration = 30 * time.Second
	defaultLongRunInterval = time.Second
	longRunProgressWidth   = 96
	longRunClearLine       = "\r\x1b[2K"
)

type longRunOutput struct {
	writer         io.Writer
	seeker         io.Seeker
	progressOffset int64
	seekable       bool
	terminal       bool
	progressActive bool
	lastProgress   string
}

func newLongRunCommand(wait func(time.Duration)) *cobra.Command {
	durationRaw := defaultLongRunDuration.String()
	intervalRaw := defaultLongRunInterval.String()

	cmd := &cobra.Command{
		Use:         "long-run",
		Short:       "Run a long task with mixed persistent and overwritten output",
		Description: "Continuously append permanent log lines while updating an overwritten progress line, then return a partially different final capture.",
		Example:     "mock long-run\nmock long-run --duration 5s --interval 500ms",
		ParamFields: []cobra.HelpField{
			optionalField("duration", "string", defaultLongRunDuration.String(), "Total simulated task duration"),
			optionalField("interval", "string", defaultLongRunInterval.String(), "Delay between progress updates"),
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

			output := newLongRunOutput(cmd.OutOrStdout())
			if err := output.WritePermanent(fmt.Sprintf(
				"event=start duration=%s interval=%s message=\"mock long task started\"",
				duration,
				interval,
			)); err != nil {
				return err
			}
			if err := output.WriteProgress(0, 0, duration, false); err != nil {
				return err
			}

			elapsed := time.Duration(0)
			sequence := 1
			for elapsed < duration {
				step := interval
				if remaining := duration - elapsed; remaining < step {
					step = remaining
				}
				wait(step)
				elapsed += step
				percent := longRunPercent(elapsed, duration)

				if err := output.WriteProgress(percent, elapsed, duration, false); err != nil {
					return err
				}
				if err := output.WritePermanent(fmt.Sprintf(
					"event=log sequence=%d elapsed=%s progress=%d%% message=\"mock work unit completed\"",
					sequence,
					elapsed,
					percent,
				)); err != nil {
					return err
				}
				sequence++
			}

			if err := output.WriteProgress(100, duration, duration, true); err != nil {
				return err
			}
			return output.WritePermanent(fmt.Sprintf(
				"result=success duration=%s logCount=%d message=\"mock long task completed\"",
				duration,
				sequence-1,
			))
		},
	}

	cmd.Flags().StringVar(&durationRaw, "duration", durationRaw, "Total simulated task duration")
	cmd.Flags().StringVar(&intervalRaw, "interval", intervalRaw, "Delay between progress updates")
	return cmd
}

func newLongRunOutput(writer io.Writer) *longRunOutput {
	output := &longRunOutput{writer: writer}
	if file, ok := writer.(*os.File); ok {
		if info, err := file.Stat(); err == nil {
			output.terminal = info.Mode()&os.ModeCharDevice != 0
		}
	}
	if seeker, ok := writer.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekCurrent); err == nil {
			output.seeker = seeker
			output.seekable = true
		}
	}
	return output
}

func (o *longRunOutput) WritePermanent(line string) error {
	if o == nil || o.writer == nil {
		return nil
	}
	if o.seekable {
		if _, err := o.seeker.Seek(0, io.SeekEnd); err != nil {
			return err
		}
		_, err := fmt.Fprintln(o.writer, line)
		return err
	}

	if o.progressActive {
		if err := o.clearProgressLine(); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(o.writer, line); err != nil {
		return err
	}
	if o.progressActive {
		_, err := fmt.Fprint(o.writer, o.progressPrefix(), o.lastProgress)
		return err
	}
	return nil
}

func (o *longRunOutput) WriteProgress(percent int, elapsed, duration time.Duration, complete bool) error {
	if o == nil || o.writer == nil {
		return nil
	}
	status := "working"
	if complete {
		status = "success"
	}
	progress := fmt.Sprintf(
		"progress=%03d%% elapsed=%s/%s status=%s",
		percent,
		elapsed,
		duration,
		status,
	)
	if !o.terminal {
		if padding := longRunProgressWidth - len(progress); padding > 0 {
			progress += strings.Repeat(" ", padding)
		}
	}

	if o.seekable {
		if !o.progressActive {
			offset, err := o.seeker.Seek(0, io.SeekEnd)
			if err != nil {
				return err
			}
			o.progressOffset = offset
		} else if _, err := o.seeker.Seek(o.progressOffset, io.SeekStart); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(o.writer, progress); err != nil {
			return err
		}
		o.progressActive = true
		o.lastProgress = progress
		return nil
	}

	if _, err := fmt.Fprint(o.writer, o.progressPrefix(), progress); err != nil {
		return err
	}
	o.progressActive = !complete
	o.lastProgress = progress
	if complete {
		_, err := fmt.Fprintln(o.writer)
		return err
	}
	return nil
}

func (o *longRunOutput) clearProgressLine() error {
	if o.terminal {
		_, err := fmt.Fprint(o.writer, longRunClearLine)
		return err
	}
	_, err := fmt.Fprintf(o.writer, "\r%s\r", strings.Repeat(" ", longRunProgressWidth))
	return err
}

func (o *longRunOutput) progressPrefix() string {
	if o.terminal {
		return longRunClearLine
	}
	return "\r"
}

func longRunPercent(elapsed, duration time.Duration) int {
	if elapsed >= duration {
		return 100
	}
	percent := int(float64(elapsed) / float64(duration) * 100)
	if percent < 0 {
		return 0
	}
	return percent
}
