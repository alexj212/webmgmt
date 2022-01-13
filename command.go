package webmgmt

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"
)

// CommandFunc function definition for a command
type CommandFunc func(client Client, args *CommandArgs) error

// ExecLevel type for admin levels
type ExecLevel int32

// ALL level for all users
const ALL = ExecLevel(0)

// USER level for default users
const USER = ExecLevel(1)

// ADMIN level for admin users
const ADMIN = ExecLevel(2)

// Command is just that, a command for your application.
// E.g.  'go run ...' - 'run' is the command. Cobra requires
// you to define the usage and description as part of your command
// definition to ensure usability.
type Command struct {
	Exec CommandFunc

	// Use is the word to execute the command
	Use string

	// Short is the short description shown in the 'help' output.
	Short string

	// Long is the long message shown in the 'help <this-command>' output.
	Long string

	// Example is examples of how to use the command.
	Example string

	// Hidden defines, if this command is hidden and should NOT show up in the list of available commands.
	Hidden bool
	// Version defines the version for this command. If this value is non-empty and the command does not
	// define a "version" flag, a "version" boolean flag will be added to the command and, if specified,
	// will print content of the "Version" variable.
	Version string

	// Deprecated defines, if this command is deprecated and should print this string when used.
	Deprecated string

	ExecLevel ExecLevel

	// commands is the list of commands supported by this program.
	commands []*Command
	// parent is a parent command for this command.
	parent *Command
	// helpTemplate is help template defined by user.
	helpTemplate string
	// helpFunc is help func defined by user.
	helpFunc func(*Command, []string, Client)
	// helpCommand is command with usage 'help'. If it's not defined by user,
	// cobra uses default help command.
	helpCommand *Command
	// usageFunc is usage func defined by user.
	usageFunc func(*Command, io.Writer) error
	// usageTemplate is usage template defined by user.
	usageTemplate string
	// versionTemplate is the version template defined by user.
	versionTemplate string

	// Max lengths of commands' string lengths for use in padding.
	commandsMaxUseLen         int
	commandsMaxCommandPathLen int
	commandsMaxNameLen        int
	// commandsAreSorted defines, if command slice are sorted or not.
	commandsAreSorted bool

	// flags is full set of flags.
	//flags *flag.FlagSet

	// DisableSuggestions disables the suggestions based on Levenshtein distance
	// that go along with 'unknown command' messages.
	DisableSuggestions bool
	// SuggestionsMinimumDistance defines minimum levenshtein distance to display suggestions.
	// Must be > 0.
	SuggestionsMinimumDistance int
	// SuggestFor is an array of command names for which this command will be suggested -
	// similar to aliases but only suggests.
	SuggestFor []string

	FlagSet  *flag.FlagSet
	HasFlags bool
}

// SetHelpFunc sets help function. Can be defined by Application.
func (c *Command) SetHelpFunc(f func(*Command, []string, Client)) {
	c.helpFunc = f
}

// SetHelpCommand sets help command.
func (c *Command) SetHelpCommand(cmd *Command) {
	c.helpCommand = cmd
}

// SetHelpTemplate sets help template to be used. Application can use it to set custom template.
func (c *Command) SetHelpTemplate(s string) {
	c.helpTemplate = s
}

// HasParent determines if the command is a child command.
func (c *Command) HasParent() bool {
	return c.parent != nil
}

// Parent returns a commands parent command.
func (c *Command) Parent() *Command {
	return c.parent
}

// Runnable determines if the command is itself runnable.
func (c *Command) Runnable() bool {
	return c.Exec != nil
}

// HasSubCommands determines if the command has children commands.
func (c *Command) HasSubCommands() bool {
	return len(c.commands) > 0
}

// UseLine puts out the full usage for a given command (including parents).
func (c *Command) UseLine() string {
	var useline string
	if c.HasParent() && c.Runnable() && c.parent.Runnable() {
		useline = c.parent.Use + " " + c.Use
	} else {
		useline = c.Use
	}

	if c.HasFlags && !strings.Contains(useline, "[flags]") {
		useline += " [flags]"
	}
	return useline
}

// Name returns the command's name: the first word in the use line.
func (c *Command) Name() string {
	name := c.Use
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// IsAvailableCommand determines if a command is available as a non-help command
// (this includes all non deprecated/hidden commands).
func (c *Command) IsAvailableCommand() bool {
	if len(c.Deprecated) != 0 || c.Hidden {
		return false
	}

	if c.HasParent() && c.Parent().helpCommand == c {
		return false
	}

	if c.Runnable() || c.HasAvailableSubCommands() {
		return true
	}

	return false
}

// IsAdditionalHelpTopicCommand determines if a command is an additional
// help topic command; additional help topic command is determined by the
// fact that it is NOT runnable/hidden/deprecated, and has no sub commands that
// are runnable/hidden/deprecated.
// Concrete example: https://github.com/spf13/cobra/issues/393#issuecomment-282741924.
func (c *Command) IsAdditionalHelpTopicCommand() bool {
	// if a command is runnable, deprecated, or hidden it is not a 'help' command
	if c.Runnable() || len(c.Deprecated) != 0 || c.Hidden {
		return false
	}

	// if any non-help sub commands are found, the command is not a 'help' command
	for _, sub := range c.commands {
		if !sub.IsAdditionalHelpTopicCommand() {
			return false
		}
	}

	// the command either has no sub commands, or no non-help sub commands
	return true
}

// HasHelpSubCommands determines if a command has any available 'help' sub commands
// that need to be shown in the usage/help default template under 'additional help
// topics'.
func (c *Command) HasHelpSubCommands() bool {
	// return true on the first found available 'help' sub command
	for _, sub := range c.commands {
		if sub.IsAdditionalHelpTopicCommand() {
			return true
		}
	}

	// the command either has no sub commands, or no available 'help' sub commands
	return false
}

// HasAvailableSubCommands determines if a command has available sub commands that
// need to be shown in the usage/help default template under 'available commands'.
func (c *Command) HasAvailableSubCommands() bool {
	// return true on the first found available (non deprecated/help/hidden)
	// sub command
	for _, sub := range c.commands {
		if sub.IsAvailableCommand() {
			return true
		}
	}

	// the command either has no sub commands, or no available (non deprecated/help/hidden)
	// sub commands
	return false
}

// UsageFunc returns either the function set by SetUsageFunc for this command
// or a parent, or it returns a default usage function.
func (c *Command) UsageFunc() (f func(*Command, io.Writer) error) {
	if c.usageFunc != nil {
		return c.usageFunc
	}
	if c.HasParent() {
		return c.Parent().UsageFunc()
	}
	return func(c *Command, io io.Writer) error {

		err := tmpl(io, c.UsageTemplate(), c)
		if err != nil {
			_, _ = io.Write([]byte(err.Error()))
		}
		return err
	}
}

// Usage puts out the usage for the command.
// Used when a user provides invalid input.
// Can be defined by user by overriding UsageFunc.
func (c *Command) Usage(io io.Writer) error {
	return c.UsageFunc()(c, io)
}

// HelpFunc returns either the function set by SetHelpFunc for this command
// or a parent, or it returns a function with default help behavior.
func (c *Command) HelpFunc() func(*Command, []string, Client) {
	if c.helpFunc != nil {
		return c.helpFunc
	}
	if c.HasParent() {
		return c.Parent().HelpFunc()
	}
	return func(c *Command, a []string, client Client) {

		err := tmpl(client.StdOut(), c.HelpTemplate(), c)
		if err != nil {
			_, _ = client.StdOut().Write([]byte(err.Error()))
		}
	}
}

// Help puts out the help for the command.
// Used when a user calls help [command].
// Can be defined by user by overriding HelpFunc.
func (c *Command) Help(client Client) error {
	c.HelpFunc()(c, []string{}, client)
	return nil
}

// UsageString returns usage string.
func (c *Command) UsageString() string {
	bb := new(bytes.Buffer)
	err := c.Usage(bb)
	if err != nil {
		return fmt.Sprintf("UsageString error: %v", err)
	}
	return bb.String()
}

var minUsagePadding = 25

// UsagePadding return padding for the usage.
func (c *Command) UsagePadding() int {
	if c.parent == nil || minUsagePadding > c.parent.commandsMaxUseLen {
		return minUsagePadding
	}
	return c.parent.commandsMaxUseLen
}

var minCommandPathPadding = 11

// CommandPathPadding return padding for the command path.
func (c *Command) CommandPathPadding() int {
	if c.parent == nil || minCommandPathPadding > c.parent.commandsMaxCommandPathLen {
		return minCommandPathPadding
	}
	return c.parent.commandsMaxCommandPathLen
}

var minNamePadding = 11

// NamePadding returns padding for the name.
func (c *Command) NamePadding() int {
	if c.parent == nil || minNamePadding > c.parent.commandsMaxNameLen {
		return minNamePadding
	}
	return c.parent.commandsMaxNameLen
}

// UsageTemplate returns usage template for the command.
func (c *Command) UsageTemplate() string {
	if c.usageTemplate != "" {
		return c.usageTemplate
	}

	if c.HasParent() {
		return c.parent.UsageTemplate()
	}
	return `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{padSpaceAfter .CommandPath}}[command] --help" for more information about a command.{{end}}
`
}

// HelpTemplate return help template for the command.
func (c *Command) HelpTemplate() string {
	if c.helpTemplate != "" {
		return c.helpTemplate
	}

	if c.HasParent() {
		return c.parent.HelpTemplate()
	}
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

// VersionTemplate return version template for the command.
func (c *Command) VersionTemplate() string {
	if c.versionTemplate != "" {
		return c.versionTemplate
	}

	if c.HasParent() {
		return c.parent.VersionTemplate()
	}
	return `{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
`
}

// HasExample determines if the command has example.
func (c *Command) HasExample() bool {
	return len(c.Example) > 0
}

func (c *Command) findSuggestions(arg string) string {
	if c.DisableSuggestions {
		return ""
	}
	if c.SuggestionsMinimumDistance <= 0 {
		c.SuggestionsMinimumDistance = 2
	}
	suggestionsString := ""
	if suggestions := c.SuggestionsFor(arg); len(suggestions) > 0 {
		suggestionsString += "\n\nDid you mean this?\n"
		for _, s := range suggestions {
			suggestionsString += fmt.Sprintf("\t%v\n", s)
		}
	}
	return suggestionsString
}

// SuggestionsFor provides suggestions for the typedName.
func (c *Command) SuggestionsFor(typedName string) []string {
	suggestions := []string{}
	for _, cmd := range c.commands {
		if cmd.IsAvailableCommand() {
			levenshteinDistance := ld(typedName, cmd.Name(), true)
			suggestByLevenshtein := levenshteinDistance <= c.SuggestionsMinimumDistance
			suggestByPrefix := strings.HasPrefix(strings.ToLower(cmd.Name()), strings.ToLower(typedName))
			if suggestByLevenshtein || suggestByPrefix {
				suggestions = append(suggestions, cmd.Name())
			}
			for _, explicitSuggestion := range cmd.SuggestFor {
				if strings.EqualFold(typedName, explicitSuggestion) {
					suggestions = append(suggestions, cmd.Name())
				}
			}
		}
	}
	return suggestions
}

// EnableCommandSorting controls sorting of the slice of commands, which is turned on by default.
// To disable sorting, set it to false.
var EnableCommandSorting = true

// Sorts commands by their names.
type commandSorterByName []*Command

func (c commandSorterByName) Len() int           { return len(c) }
func (c commandSorterByName) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c commandSorterByName) Less(i, j int) bool { return c[i].Name() < c[j].Name() }

// Commands returns a sorted slice of child commands.
func (c *Command) Commands() []*Command {
	// do not sort commands if it already sorted or sorting was disabled
	if EnableCommandSorting && !c.commandsAreSorted {
		sort.Sort(commandSorterByName(c.commands))
		c.commandsAreSorted = true
	}
	return c.commands
}

// AddCommand adds one or more commands to this parent command.
func (c *Command) AddCommand(cmds ...*Command) {
	for i, x := range cmds {
		if cmds[i] == c {
			panic("Command can't be a child of itself")
		}
		cmds[i].parent = c
		// update max lengths
		usageLen := len(x.Use)
		if usageLen > c.commandsMaxUseLen {
			c.commandsMaxUseLen = usageLen
		}
		commandPathLen := len(x.CommandPath())
		if commandPathLen > c.commandsMaxCommandPathLen {
			c.commandsMaxCommandPathLen = commandPathLen
		}
		nameLen := len(x.Name())
		if nameLen > c.commandsMaxNameLen {
			c.commandsMaxNameLen = nameLen
		}

		c.commands = append(c.commands, x)
		c.commandsAreSorted = false
	}
}

// RemoveCommand removes one or more commands from a parent command.
func (c *Command) RemoveCommand(cmds ...*Command) {
	commands := []*Command{}
main:
	for _, command := range c.commands {
		for _, cmd := range cmds {
			if command == cmd {
				command.parent = nil
				continue main
			}
		}
		commands = append(commands, command)
	}
	c.commands = commands
	// recompute all lengths
	c.commandsMaxUseLen = 0
	c.commandsMaxCommandPathLen = 0
	c.commandsMaxNameLen = 0
	for _, command := range c.commands {
		usageLen := len(command.Use)
		if usageLen > c.commandsMaxUseLen {
			c.commandsMaxUseLen = usageLen
		}
		commandPathLen := len(command.CommandPath())
		if commandPathLen > c.commandsMaxCommandPathLen {
			c.commandsMaxCommandPathLen = commandPathLen
		}
		nameLen := len(command.Name())
		if nameLen > c.commandsMaxNameLen {
			c.commandsMaxNameLen = nameLen
		}
	}
}

// CommandPath returns the full path to this command.
func (c *Command) CommandPath() string {
	if c.HasParent() && c.parent.Runnable() {
		return c.Parent().CommandPath() + " " + c.Name()
	}

	return c.Name()
}

// IsSubCommandAvailable checks if sub command exists
func (c *Command) IsSubCommandAvailable(client Client, cmd string) bool {

	if cmd == "help" {
		return true
	}

	for _, command := range c.commands {
		if cmd == command.Use && command.Exec != nil {

			if client.ExecLevel() >= command.ExecLevel {
				return true
			}
			return false
		}
	}
	return false
}

// Execute runs a command thru execution
//nolint:gocyclo,goreportcard
func (c *Command) Execute(client Client, cmdLine *CommandArgs) {

	if cmdLine.CmdName == "help" {
		fmt.Printf("exec 0: %v\n", cmdLine.Debug())

		if len(cmdLine.Args) > 0 {
			for _, cmd := range c.commands {
				if cmd.Name() == cmdLine.Args[0] {
					cmd.Help(client)
					return
				}
			}

			val := c.findSuggestions(cmdLine.Args[0])
			if val == "" {
				c.Help(client)
			} else {
				client.StdErr().Write([]byte(val))
			}
			return
		}

		c.Help(client)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("panic occurred: %v\n", err)
			client.StdErr().Write([]byte(msg))
		}
	}()

	for _, cmd := range c.commands {
		if cmd.Name() == cmdLine.CmdName {
			fmt.Printf("Command matched: %v\n", cmdLine.CmdName)
			shifted, err := cmdLine.Shift()
			if err == nil {
				if cmd.IsSubCommandAvailable(client, shifted.CmdName) {
					fmt.Printf("exec A: %v\n", shifted.Debug())
					cmd.Execute(client, shifted)
					return
				}

				fmt.Printf("exec B Command executing against : %v\n", cmdLine.Debug())
				cmd.Exec(client, cmdLine)
				return

			}

			fmt.Printf("exec C Command executing againt : %v\n", cmdLine.Debug())
			if len(cmdLine.Args) > 0 && (cmdLine.Args[0] == "help" || cmdLine.Args[0] == "--help" || cmdLine.Args[0] == "-help") {
				cmd.Help(client)
			} else {
				cmd.Exec(client, cmdLine)
			}

			return

		}
	}

	fmt.Printf("exec D: %v\n", cmdLine.Debug())

	val := c.findSuggestions(cmdLine.CmdName)
	if val == "" {
		c.Help(client)
	} else {
		client.StdErr().Write([]byte(val))
	}
}
