package webmgmt

import (
    "bytes"
    "errors"
    "flag"
    "fmt"
    "io"
    "strings"

    "github.com/kballard/go-shellquote"
)

// CommandArgs is a struct that is used to store the contents of a parsed command line string.
type CommandArgs struct {
    CmdLine string
    CmdName string
    Args    []string
    FlagSet *flag.FlagSet
    output  io.Writer
}

// String will return the CmdLine the original one that is parsed.
func (c *CommandArgs) String() string {
    return c.CmdLine
}

// PealOff will return a command line string after N commands have been pealed off from the front of the command line.
func (c *CommandArgs) PealOff(pos int) string {
    var buffer bytes.Buffer

    for i := pos; i < len(c.Args); i++ {

        if strings.Contains(c.Args[i], " ") || strings.Contains(c.Args[i], "\t") {
            buffer.WriteString("\"")
            buffer.WriteString(c.Args[i])
            buffer.WriteString("\"")
        } else {
            buffer.WriteString(c.Args[i])
        }

        if i < len(c.Args) {
            buffer.WriteString(" ")
        }
    }

    return buffer.String()
}

// Debug will return a string listing the original command line and the parsed arguments and flags
func (c *CommandArgs) Debug() string {
    var buffer bytes.Buffer

    buffer.WriteString(fmt.Sprintf("fullcmd: %v  args: [%v]", c.CmdLine, strings.Join(c.Args, ", ")))
    buffer.WriteString(fmt.Sprintf("flagSet.NArg(): %v\n", c.FlagSet.NArg()))
    for i, val := range c.FlagSet.Args() {
        buffer.WriteString(fmt.Sprintf("flagSet.Args()[%d]: %v\n", i, val))
    }
    return buffer.String()
}

// Parse will use the defined FlagSet parsed and return an error if help is invoked or invalid flags
func (c *CommandArgs) Parse() error {
    return c.FlagSet.Parse(c.Args)
}

// Shift will return a new CommandArgs after shifting the first cmd in the string
func (c *CommandArgs) Shift() (*CommandArgs, error) {
    if strings.HasPrefix(c.CmdLine, c.CmdName) {
        newCmdLine := c.CmdLine[len(c.CmdName):]
        return NewCommandArgs(newCmdLine, c.output)
    } else {
        return nil, errors.New("Unable to shift command line. ")
    }
}

// NewCommandArgs will take a raw string and output. The Raw string will be parsed into a CommandArgs structure.
// If an error is encountered such as empty command string it will be returned and the CommandArgs will be nil.
func NewCommandArgs(cmdLine string, output io.Writer) (*CommandArgs, error) {

    if cmdLine == "" {
        return nil, errors.New("empty command string")
    }

    words, err := shellquote.Split(cmdLine)

    if err != nil {
        return nil, err
    } else if len(words) == 0 {
        return nil, errors.New("no words parsed")

    } else {

        invoke := &CommandArgs{}
        invoke.CmdLine = cmdLine
        invoke.CmdName = words[0]
        invoke.Args = words[1:]
        invoke.output = output
        invoke.FlagSet = flag.NewFlagSet(invoke.CmdName, flag.ContinueOnError)

        if output != nil {
            invoke.FlagSet.SetOutput(output)
        }
        return invoke, nil
    }
    // topic Hello World Stinky
    // topic "Hello World" Stinky
    // group 156 topic Hello World Stinky

}
