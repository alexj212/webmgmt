package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"

	"github.com/alexj212/webmgmt"
)

func main() {
	cmd := &webmgmt.Command{
		Use: "cls",
		Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {

			return
		},
		Short: "send cls event to terminal client",
		Long: `
This is the long version of the description to the 
cls command which will
clear the
screen

`,
		ExecLevel: webmgmt.ALL,
	}

	cmd.FlagSet = flag.NewFlagSet(cmd.Use, flag.ContinueOnError)
	_ = cmd.FlagSet.Int("cnt", 5, "number of lines to print")

	bb := new(bytes.Buffer)

	cmd.Help(bb)
	fmt.Printf(bb.String())
}
