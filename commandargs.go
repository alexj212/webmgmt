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

type CommandArgs struct {
    CmdLine string
    CmdName string
    Args    []string
    FlagSet *flag.FlagSet
    output  io.Writer
}

func (c *CommandArgs) String() string {
    return c.CmdLine
}

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

func (c *CommandArgs) Debug() string {
    var buffer bytes.Buffer

    buffer.WriteString(fmt.Sprintf("fullcmd: %v  args: [%v]", c.CmdLine, strings.Join(c.Args, ", ")))
    buffer.WriteString(fmt.Sprintf("flagSet.NArg(): %v\n", c.FlagSet.NArg()))
    for i, val := range c.FlagSet.Args() {
        buffer.WriteString(fmt.Sprintf("flagSet.Args()[%d]: %v\n", i, val))
    }
    return buffer.String()
}

func (c *CommandArgs) Parse() error {
    return c.FlagSet.Parse(c.Args)
}
func (c *CommandArgs) Shift() (*CommandArgs, error) {
    if strings.HasPrefix(c.CmdLine, c.CmdName) {
        newCmdLine := c.CmdLine[len(c.CmdName):]
        return NewCommandArgs(newCmdLine, c.output)
    } else {
        return nil, errors.New("Unable to shift command line. ")
    }
}

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
