package cobra

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type boolFlag interface {
	IsBoolFlag() bool
}

type PositionalArgs func(cmd *Command, args []string) error

type CompletionOptions struct {
	DisableDefaultCmd bool
}

type HelpField struct {
	Name        string
	Type        string
	Required    string
	Default     string
	Description string
}

type HelpSection struct {
	Title string
	Body  string
}

type Command struct {
	Use               string
	Short             string
	Description       string
	Example           string
	ArgFields         []HelpField
	ParamFields       []HelpField
	HelpSections      []HelpSection
	Args              PositionalArgs
	RunE              func(cmd *Command, args []string) error
	SilenceUsage      bool
	SilenceErrors     bool
	CompletionOptions CompletionOptions

	parent   *Command
	children []*Command
	flagSet  *flag.FlagSet
	args     []string
	in       io.Reader
	out      io.Writer
	err      io.Writer
}

func (c *Command) AddCommand(children ...*Command) {
	for _, child := range children {
		if child == nil {
			continue
		}
		child.parent = c
		c.children = append(c.children, child)
	}
}

func (c *Command) SetArgs(args []string) {
	c.args = append([]string(nil), args...)
}

func (c *Command) SetIn(r io.Reader) {
	c.in = r
}

func (c *Command) SetOut(w io.Writer) {
	c.out = w
}

func (c *Command) SetErr(w io.Writer) {
	c.err = w
}

func (c *Command) InOrStdin() io.Reader {
	if c.in != nil {
		return c.in
	}
	if c.parent != nil {
		return c.parent.InOrStdin()
	}
	return os.Stdin
}

func (c *Command) OutOrStdout() io.Writer {
	if c.out != nil {
		return c.out
	}
	if c.parent != nil {
		return c.parent.OutOrStdout()
	}
	return os.Stdout
}

func (c *Command) ErrOrStderr() io.Writer {
	if c.err != nil {
		return c.err
	}
	if c.parent != nil {
		return c.parent.ErrOrStderr()
	}
	return os.Stderr
}

func (c *Command) Flags() *flag.FlagSet {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		c.flagSet.SetOutput(io.Discard)
	}
	return c.flagSet
}

func (c *Command) Name() string {
	name := strings.TrimSpace(c.Use)
	if name == "" {
		return ""
	}
	return strings.Fields(name)[0]
}

func (c *Command) CommandPath() string {
	if c.parent == nil {
		return c.Name()
	}
	parentPath := c.parent.CommandPath()
	if parentPath == "" {
		return c.Name()
	}
	if c.Name() == "" {
		return parentPath
	}
	return parentPath + " " + c.Name()
}

func (c *Command) Execute() error {
	return c.execute(c.args)
}

func (c *Command) execute(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "help":
			return c.helpFor(args[1:])
		case "-h", "--help":
			return c.printHelp()
		}
	}

	if child := c.findChild(args); child != nil {
		return child.execute(args[1:])
	}

	if len(c.children) > 0 && len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		return fmt.Errorf("unknown command %q for %q", args[0], c.CommandPath())
	}

	if len(c.children) > 0 && len(args) == 0 {
		return c.printHelp()
	}

	parsedArgs, err := c.parseFlags(args)
	if err != nil {
		return err
	}

	for _, arg := range parsedArgs {
		if arg == "-h" || arg == "--help" {
			return c.printHelp()
		}
	}

	if c.Args != nil {
		if err := c.Args(c, parsedArgs); err != nil {
			return err
		}
	}
	if c.RunE != nil {
		return c.RunE(c, parsedArgs)
	}
	return c.printHelp()
}

func (c *Command) helpFor(args []string) error {
	target := c
	for _, arg := range args {
		next := target.findChild([]string{arg})
		if next == nil {
			return fmt.Errorf("unknown help topic %q for %q", arg, target.CommandPath())
		}
		target = next
	}
	return target.printHelp()
}

func (c *Command) findChild(args []string) *Command {
	if len(args) == 0 {
		return nil
	}
	name := args[0]
	for _, child := range c.children {
		if child.Name() == name {
			return child
		}
	}
	return nil
}

func (c *Command) parseFlags(args []string) ([]string, error) {
	if c.flagSet == nil {
		return args, nil
	}

	positionals := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			positionals = append(positionals, args[i+1:]...)
			break
		}
		if !strings.HasPrefix(arg, "--") || arg == "--help" {
			positionals = append(positionals, arg)
			continue
		}

		nameValue := strings.TrimPrefix(arg, "--")
		name, value, hasValue := strings.Cut(nameValue, "=")
		f := c.flagSet.Lookup(name)
		if f == nil {
			return nil, fmt.Errorf("unknown flag: --%s", name)
		}
		if !hasValue {
			if bf, ok := f.Value.(boolFlag); ok && bf.IsBoolFlag() {
				value = "true"
			} else {
				i++
				if i >= len(args) {
					return nil, fmt.Errorf("flag needs an argument: --%s", name)
				}
				value = args[i]
			}
		}
		if err := f.Value.Set(value); err != nil {
			return nil, err
		}
	}
	return positionals, nil
}

func (c *Command) printHelp() error {
	_, err := fmt.Fprint(c.OutOrStdout(), c.helpText())
	return err
}

func (c *Command) helpText() string {
	var b strings.Builder

	writeUsage(&b, c)
	writeDescription(&b, c)
	writeAvailableCommands(&b, c)
	writeFlags(&b, c)
	writeFieldSection(&b, "Args fields", c.ArgFields)
	writeFieldSection(&b, "Params fields", c.ParamFields)
	writeHelpSections(&b, c)
	writeExamples(&b, c)

	return b.String()
}

func writeUsage(b *strings.Builder, c *Command) {
	if c.CommandPath() == "" {
		return
	}

	b.WriteString("Usage:\n")
	fmt.Fprintf(b, "  %s", c.CommandPath())
	if c.Name() != "" && c.Use != "" && len(strings.Fields(c.Use)) > 1 {
		fmt.Fprintf(b, " %s", strings.Join(strings.Fields(c.Use)[1:], " "))
	}
	if c.flagSet != nil {
		b.WriteString(" [flags]")
	}
	b.WriteString("\n")
}

func writeDescription(b *strings.Builder, c *Command) {
	description := c.Description
	if description == "" {
		description = c.Short
	}
	if description == "" {
		return
	}

	b.WriteString("\nDescription:\n")
	writeIndentedBlock(b, description)
}

func writeAvailableCommands(b *strings.Builder, c *Command) {
	if len(c.children) == 0 {
		return
	}

	b.WriteString("\nAvailable Commands:\n")
	for _, child := range c.children {
		fmt.Fprintf(b, "  %-10s %s\n", child.Name(), child.Short)
	}
	if !c.CompletionOptions.DisableDefaultCmd {
		b.WriteString("  completion Generate the autocompletion script\n")
	}
	b.WriteString("  help       Help about any command\n")
}

func writeFlags(b *strings.Builder, c *Command) {
	b.WriteString("\nFlags:\n")

	if c.flagSet != nil {
		c.flagSet.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(b, "  %-18s %s\n", formatFlagUsage(f), f.Usage)
		})
	}

	b.WriteString("  -h, --help         help for this command\n")
}

func writeFieldSection(b *strings.Builder, title string, fields []HelpField) {
	if len(fields) == 0 {
		return
	}

	nameWidth := len("name")
	typeWidth := len("type")
	requiredWidth := len("required")
	defaultWidth := len("default")

	for _, field := range fields {
		nameWidth = max(nameWidth, len(field.Name))
		typeWidth = max(typeWidth, len(field.Type))
		requiredWidth = max(requiredWidth, len(field.Required))
		defaultWidth = max(defaultWidth, len(field.Default))
	}

	fmt.Fprintf(b, "\n%s:\n", title)
	fmt.Fprintf(b, "  %-*s   %-*s   %-*s   %-*s   %s\n",
		nameWidth, "name",
		typeWidth, "type",
		requiredWidth, "required",
		defaultWidth, "default",
		"description",
	)

	for _, field := range fields {
		fmt.Fprintf(b, "  %-*s   %-*s   %-*s   %-*s   %s\n",
			nameWidth, field.Name,
			typeWidth, field.Type,
			requiredWidth, field.Required,
			defaultWidth, field.Default,
			field.Description,
		)
	}
}

func writeHelpSections(b *strings.Builder, c *Command) {
	for _, section := range c.HelpSections {
		if strings.TrimSpace(section.Title) == "" || strings.TrimSpace(section.Body) == "" {
			continue
		}
		fmt.Fprintf(b, "\n%s:\n", section.Title)
		writeIndentedBlock(b, section.Body)
	}
}

func writeExamples(b *strings.Builder, c *Command) {
	if strings.TrimSpace(c.Example) == "" {
		return
	}

	b.WriteString("\nExamples:\n")
	writeIndentedBlock(b, c.Example)
}

func writeIndentedBlock(b *strings.Builder, body string) {
	for _, line := range strings.Split(strings.TrimRight(body, "\n"), "\n") {
		fmt.Fprintf(b, "  %s\n", line)
	}
}

func formatFlagUsage(f *flag.Flag) string {
	name := "--" + f.Name
	typeName := flagTypeName(f)
	if typeName == "" {
		return name
	}
	return name + " " + typeName
}

func flagTypeName(f *flag.Flag) string {
	if boolFlag, ok := f.Value.(interface{ IsBoolFlag() bool }); ok && boolFlag.IsBoolFlag() {
		return "bool"
	}

	typeName := fmt.Sprintf("%T", f.Value)
	switch {
	case strings.HasSuffix(typeName, ".stringValue"):
		return "string"
	case strings.HasSuffix(typeName, ".intValue"):
		return "int"
	case strings.HasSuffix(typeName, ".int64Value"):
		return "int64"
	case strings.HasSuffix(typeName, ".uintValue"):
		return "uint"
	case strings.HasSuffix(typeName, ".uint64Value"):
		return "uint64"
	case strings.HasSuffix(typeName, ".float64Value"):
		return "float"
	case strings.HasSuffix(typeName, ".durationValue"):
		return "duration"
	default:
		return ""
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func NoArgs(cmd *Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("accepts 0 arg(s), received %d", len(args))
	}
	return nil
}

func ExactArgs(n int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("accepts %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}

func MinimumNArgs(n int) PositionalArgs {
	return func(cmd *Command, args []string) error {
		if len(args) < n {
			return fmt.Errorf("requires at least %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}
